from __future__ import annotations

import argparse
from collections import Counter
from contextlib import redirect_stdout
import importlib.util
import io
import json
from pathlib import Path
import tempfile
import unittest
from unittest import mock


SCRIPT_PATH = (
    Path(__file__).resolve().parents[1]
    / "scripts"
    / "collect_silmari_rebrand_inventory.py"
)
BASELINE_ARTIFACT_PATH = (
    Path(__file__).resolve().parents[1]
    / "thoughts"
    / "searchable"
    / "shared"
    / "plans"
    / "2026-06-29-silmari-rebrand-baseline-inventory.json"
)
ALLOWED_CLASSIFICATIONS = {
    "edit",
    "preserve",
    "historical",
    "runtime-compatible",
}
EXPECTED_UNTRACKED_ROOT_ARTIFACTS = {
    ".python-version",
    "pyproject.toml",
    "uv.lock",
}


def load_module():
    spec = importlib.util.spec_from_file_location(
        "collect_silmari_rebrand_inventory", SCRIPT_PATH
    )
    if spec is None or spec.loader is None:
        raise RuntimeError("Unable to load inventory collector module")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


inventory = load_module()
SDK_PY_PACKAGE = "sdk/python/" + "agent" + "field"
HELM_CHART_DIR = "deployments/helm/" + "agent" + "field"
BENCHMARK_RESULT = (
    "examples/benchmarks/100k-scale/results/" + "Agent" + "Field_Python.json"
)
UTM_LINKS_PATH = "assets/utm-links.csv"


def load_baseline_artifact() -> dict[str, object]:
    return json.loads(BASELINE_ARTIFACT_PATH.read_text(encoding="utf-8"))


class ClassifyPathTests(unittest.TestCase):
    def test_root_readme_is_edit(self) -> None:
        classification, reason = inventory.classify_path("README.md")
        self.assertEqual(classification, "edit")
        self.assertIn("Silmari", reason)

    def test_assets_utm_links_is_edit(self) -> None:
        classification, reason = inventory.classify_path(UTM_LINKS_PATH)
        self.assertEqual(classification, "edit")
        self.assertIn("link labels", reason)

    def test_runtime_sdk_source_is_runtime_compatible(self) -> None:
        classification, reason = inventory.classify_path(f"{SDK_PY_PACKAGE}/agent.py")
        self.assertEqual(classification, "runtime-compatible")
        self.assertIn("runtime", reason)

    def test_sdk_test_file_is_preserve(self) -> None:
        classification, reason = inventory.classify_path(
            "sdk/python/tests/test_agent_server.py"
        )
        self.assertEqual(classification, "preserve")
        self.assertIn("tests", reason)

    def test_historical_benchmark_result_is_historical(self) -> None:
        classification, reason = inventory.classify_path(BENCHMARK_RESULT)
        self.assertEqual(classification, "historical")
        self.assertIn("historical", reason)

    def test_mode_context_is_runtime_compatible(self) -> None:
        classification, reason = inventory.classify_path(
            "control-plane/web/client/src/contexts/ModeContext.tsx"
        )
        self.assertEqual(classification, "runtime-compatible")
        self.assertIn("localStorage", reason)

    def test_helm_notes_is_edit(self) -> None:
        classification, reason = inventory.classify_path(
            f"{HELM_CHART_DIR}/templates/NOTES.txt"
        )
        self.assertEqual(classification, "edit")
        self.assertIn("user-visible", reason)


class BranchAssertionTests(unittest.TestCase):
    def test_assert_branch_succeeds_for_exact_branch(self) -> None:
        with mock.patch.object(
            inventory, "current_branch", return_value=inventory.EXPECTED_BRANCH
        ):
            stream = io.StringIO()
            with redirect_stdout(stream):
                exit_code = inventory.assert_branch(
                    inventory.EXPECTED_BRANCH, Path("/tmp/repo")
                )

        payload = json.loads(stream.getvalue())
        self.assertEqual(exit_code, 0)
        self.assertEqual(payload["current_branch"], inventory.EXPECTED_BRANCH)
        self.assertEqual(payload["expected_branch"], inventory.EXPECTED_BRANCH)
        self.assertTrue(payload["matches"])

    def test_assert_branch_fails_for_non_matching_branch(self) -> None:
        with mock.patch.object(inventory, "current_branch", return_value="main"):
            stream = io.StringIO()
            with redirect_stdout(stream):
                exit_code = inventory.assert_branch(
                    inventory.EXPECTED_BRANCH, Path("/tmp/repo")
                )

        payload = json.loads(stream.getvalue())
        self.assertEqual(exit_code, 1)
        self.assertEqual(payload["current_branch"], "main")
        self.assertFalse(payload["matches"])


class CollectMatchesTests(unittest.TestCase):
    def test_collect_matches_returns_empty_when_rg_finds_no_matches(self) -> None:
        proc = mock.Mock(returncode=1, stdout="", stderr="")
        with mock.patch.object(inventory.subprocess, "run", return_value=proc):
            matches = inventory.collect_matches(["LICENSE"], Path("/tmp/repo"))

        self.assertEqual(matches, {})

    def test_collect_matches_aggregates_multiple_events_for_same_file(self) -> None:
        stdout = "\n".join(
            [
                json.dumps({"type": "begin"}),
                json.dumps(
                    {
                        "type": "match",
                        "data": {
                            "path": {"text": "README.md"},
                            "submatches": [{}, {}],
                        },
                    }
                ),
                json.dumps(
                    {
                        "type": "match",
                        "data": {
                            "path": {"text": "README.md"},
                            "submatches": [{}],
                        },
                    }
                ),
                json.dumps(
                    {
                        "type": "match",
                        "data": {
                            "path": {"text": f"{SDK_PY_PACKAGE}/agent.py"},
                            "submatches": [{}, {}, {}],
                        },
                    }
                ),
            ]
        )
        proc = mock.Mock(returncode=0, stdout=stdout, stderr="")
        with mock.patch.object(inventory.subprocess, "run", return_value=proc):
            matches = inventory.collect_matches(
                ["README.md", SDK_PY_PACKAGE], Path("/tmp/repo")
            )

        self.assertEqual(
            matches,
            {"README.md": 3, f"{SDK_PY_PACKAGE}/agent.py": 3},
        )

    def test_collect_matches_raises_runtime_error_on_rg_failure(self) -> None:
        proc = mock.Mock(returncode=2, stdout="", stderr="rg failed")
        with mock.patch.object(inventory.subprocess, "run", return_value=proc):
            with self.assertRaisesRegex(RuntimeError, "rg failed"):
                inventory.collect_matches(["README.md"], Path("/tmp/repo"))


class BuildInventoryTests(unittest.TestCase):
    def test_build_inventory_handles_empty_scope(self) -> None:
        with mock.patch.object(inventory, "collect_matches", return_value={}):
            payload = inventory.build_inventory(["LICENSE"], Path("/tmp/repo"))

        self.assertEqual(payload["scanRoots"], ["LICENSE"])
        self.assertEqual(payload["classifiedFiles"], [])
        self.assertEqual(payload["classificationCounts"], {})
        self.assertEqual(payload["totalMatchedFiles"], 0)
        self.assertEqual(payload["totalMatches"], 0)

    def test_build_inventory_classifies_mixed_matches(self) -> None:
        mocked_matches = {
            "README.md": 2,
            "sdk/python/tests/test_agent_server.py": 1,
            BENCHMARK_RESULT: 4,
            f"{SDK_PY_PACKAGE}/agent.py": 3,
        }
        with mock.patch.object(inventory, "collect_matches", return_value=mocked_matches):
            payload = inventory.build_inventory(["README.md", SDK_PY_PACKAGE], Path("."))

        self.assertEqual(payload["totalMatchedFiles"], 4)
        self.assertEqual(payload["totalMatches"], 10)
        self.assertEqual(
            payload["classificationCounts"],
            {
                "edit": 1,
                "historical": 1,
                "preserve": 1,
                "runtime-compatible": 1,
            },
        )
        self.assertEqual(
            {
                entry["path"]: entry["classification"]
                for entry in payload["classifiedFiles"]
            },
            {
                "README.md": "edit",
                "sdk/python/tests/test_agent_server.py": "preserve",
                BENCHMARK_RESULT: "historical",
                f"{SDK_PY_PACKAGE}/agent.py": "runtime-compatible",
            },
        )


class ParserTests(unittest.TestCase):
    def test_default_assert_branch_target_is_exact_integration_branch(self) -> None:
        parser = inventory.build_parser()
        args = parser.parse_args(["assert-branch"])
        self.assertEqual(args.expected, inventory.EXPECTED_BRANCH)

    def test_collect_uses_repeated_path_overrides(self) -> None:
        parser = inventory.build_parser()
        args = parser.parse_args(
            ["collect", "--path", "README.md", "--path", SDK_PY_PACKAGE]
        )
        self.assertEqual(args.paths, ["README.md", SDK_PY_PACKAGE])

    def test_rg_command_includes_paths_and_exclusion_globs(self) -> None:
        cmd = inventory.rg_command(["README.md"])
        self.assertEqual(
            cmd[:5],
            ["rg", "--json", "--ignore-case", "-o", inventory.BRAND_PATTERN],
        )
        self.assertIn("README.md", cmd)
        for glob in inventory.SCAN_GLOBS:
            self.assertIn(glob, cmd)

    def test_default_scan_roots_include_assets_utm_links(self) -> None:
        self.assertIn(UTM_LINKS_PATH, inventory.SCAN_ROOTS)


class CollectCommandTests(unittest.TestCase):
    def test_cmd_collect_uses_default_scan_roots_and_writes_output(self) -> None:
        fake_inventory = {
            "classifiedFiles": [],
            "classificationCounts": {},
            "command": "rg README.md",
            "commandRef": "Brand Surface Inventory",
            "scanRoots": inventory.SCAN_ROOTS,
            "totalMatchedFiles": 0,
            "totalMatches": 0,
        }
        repo_root = Path("/tmp/repo")
        with tempfile.TemporaryDirectory() as temp_dir:
            output_path = Path(temp_dir) / "inventory.json"
            args = argparse.Namespace(paths=None, output=str(output_path))
            with (
                mock.patch.object(inventory, "repo_root", return_value=repo_root),
                mock.patch.object(
                    inventory,
                    "current_branch",
                    return_value=inventory.EXPECTED_BRANCH,
                ),
                mock.patch.object(
                    inventory, "build_inventory", return_value=fake_inventory
                ) as build_inventory,
            ):
                exit_code = inventory.cmd_collect(args)

            written = json.loads(output_path.read_text(encoding="utf-8"))

        self.assertEqual(exit_code, 0)
        build_inventory.assert_called_once_with(inventory.SCAN_ROOTS, repo_root)
        self.assertEqual(written["branch"], inventory.EXPECTED_BRANCH)
        self.assertEqual(written["brandInventoryCollected"], fake_inventory)
        self.assertTrue(
            all(
                rename_planned is False
                for rename_planned in written["changeSetBoundaryVerified"][
                    "plannedRenames"
                ].values()
            )
        )

    def test_cmd_collect_prints_boundary_notes_when_output_is_none(self) -> None:
        fake_inventory = {
            "classifiedFiles": [],
            "classificationCounts": {},
            "command": "rg LICENSE",
            "commandRef": "Brand Surface Inventory",
            "scanRoots": ["LICENSE"],
            "totalMatchedFiles": 0,
            "totalMatches": 0,
        }
        args = argparse.Namespace(paths=["LICENSE"], output=None)
        with (
            mock.patch.object(inventory, "repo_root", return_value=Path("/tmp/repo")),
            mock.patch.object(
                inventory, "current_branch", return_value=inventory.EXPECTED_BRANCH
            ),
            mock.patch.object(inventory, "build_inventory", return_value=fake_inventory),
        ):
            stream = io.StringIO()
            with redirect_stdout(stream):
                exit_code = inventory.cmd_collect(args)

        payload = json.loads(stream.getvalue())
        self.assertEqual(exit_code, 0)
        self.assertEqual(payload["brandInventoryCollected"]["scanRoots"], ["LICENSE"])
        self.assertIn(
            "No runtime route, storage schema, SDK API, package name",
            payload["changeSetBoundaryVerified"]["notes"],
        )


class BaselineArtifactTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.artifact = load_baseline_artifact()

    def test_baseline_artifact_records_branch_switch_and_untracked_files(self) -> None:
        branch_data = self.artifact["branchPrepared"]
        self.assertEqual(branch_data["expectedBranch"], inventory.EXPECTED_BRANCH)
        self.assertEqual(branch_data["branchAssertionBefore"]["exitCode"], 1)
        self.assertEqual(branch_data["branchAssertionAfter"]["exitCode"], 0)
        self.assertEqual(
            branch_data["branchAssertionAfter"]["currentBranch"],
            inventory.EXPECTED_BRANCH,
        )
        for artifact_name in EXPECTED_UNTRACKED_ROOT_ARTIFACTS:
            self.assertTrue(
                any(
                    entry.endswith(artifact_name)
                    for entry in branch_data["rootRepoUnrelatedUntrackedStillUnstaged"]
                )
            )

    def test_full_inventory_classifies_every_match_into_allowed_categories(self) -> None:
        brand_inventory = self.artifact["brandInventoryCollected"]
        classified_files = brand_inventory["classifiedFiles"]
        computed_counts = Counter(
            entry["classification"] for entry in classified_files
        )
        self.assertEqual(brand_inventory["scanRoots"], inventory.SCAN_ROOTS)
        self.assertTrue(classified_files)
        self.assertTrue(
            all(
                entry["classification"] in ALLOWED_CLASSIFICATIONS
                for entry in classified_files
            )
        )
        self.assertEqual(
            brand_inventory["classificationCounts"], dict(sorted(computed_counts.items()))
        )
        self.assertEqual(
            brand_inventory["totalMatchedFiles"], len(classified_files)
        )
        self.assertEqual(
            brand_inventory["totalMatches"],
            sum(entry["matches"] for entry in classified_files),
        )
        self.assertEqual(
            next(
                entry["classification"]
                for entry in classified_files
                if entry["path"] == UTM_LINKS_PATH
            ),
            "edit",
        )

    def test_verification_samples_cover_empty_single_file_and_runtime_boundary(self) -> None:
        verification_samples = self.artifact["verificationSamples"]
        self.assertEqual(
            verification_samples["emptyScope"]["classificationCounts"], {}
        )
        self.assertEqual(
            verification_samples["emptyScope"]["classifiedFiles"], []
        )
        self.assertEqual(
            verification_samples["singleFile"]["scanRoots"], ["README.md"]
        )
        self.assertEqual(
            verification_samples["singleFile"]["classificationCounts"], {"edit": 1}
        )
        self.assertTrue(
            verification_samples["protectedRuntime"]["classifiedFiles"]
        )
        self.assertTrue(
            all(
                entry["classification"] == "runtime-compatible"
                for entry in verification_samples["protectedRuntime"][
                    "classifiedFiles"
                ]
            )
        )

    def test_change_boundary_keeps_all_planned_runtime_renames_disabled(self) -> None:
        boundary = self.artifact["changeSetBoundaryVerified"]
        self.assertIn("env vars and YAML config paths", boundary["protectedRuntimeAreas"])
        self.assertIn("sdk/go/**/*.go", boundary["protectedRuntimeAreas"])
        self.assertNotIn("sdk/go/**", boundary["protectedRuntimeAreas"])
        self.assertIn(
            "web UI storage keys and API path identifiers",
            boundary["protectedRuntimeAreas"],
        )
        self.assertTrue(
            all(rename_planned is False for rename_planned in boundary["plannedRenames"].values())
        )
        self.assertIn("No runtime route", boundary["notes"])


if __name__ == "__main__":
    unittest.main()
