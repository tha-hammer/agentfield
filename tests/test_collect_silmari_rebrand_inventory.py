from __future__ import annotations

import importlib.util
from pathlib import Path
import unittest


SCRIPT_PATH = (
    Path(__file__).resolve().parents[1]
    / "scripts"
    / "collect_silmari_rebrand_inventory.py"
)


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
BENCHMARK_RESULT = "examples/benchmarks/100k-scale/results/" + "Agent" + "Field_Python.json"


class ClassifyPathTests(unittest.TestCase):
    def test_root_readme_is_edit(self) -> None:
        classification, reason = inventory.classify_path("README.md")
        self.assertEqual(classification, "edit")
        self.assertIn("Silmari", reason)

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


if __name__ == "__main__":
    unittest.main()
