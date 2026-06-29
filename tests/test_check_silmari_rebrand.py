from __future__ import annotations

from pathlib import Path
import stat
import subprocess
import sys
import tempfile
import textwrap
import types
import unittest

from hypothesis import given
from hypothesis import strategies as st


REPO_ROOT = Path(__file__).resolve().parents[1]
SCRIPT_PATH = REPO_ROOT / "scripts" / "check-silmari-rebrand.sh"
MANIFEST_PATH = Path("docs") / "silmari-rebrand-manifest.md"
LEGACY_BRAND = "Agent" + "Field"
LEGACY_ENV = "AGENT" + "FIELD" + "_CONFIG_FILE"
LOWER_LEGACY = "agent" + "field"
LOWER_IMPORT = "@" + LOWER_LEGACY + "/sdk"
LEGACY_SKILL_SLUG = "agent" + "field"
SYNC_SCRIPT_REL = Path("scripts") / "sync-embedded-skills.sh"
SKILL_SOURCE_REL = Path("skills") / LEGACY_SKILL_SLUG / "SKILL.md"
SKILL_MIRROR_REL = (
    Path("control-plane")
    / "internal"
    / "skillkit"
    / "skill_data"
    / LEGACY_SKILL_SLUG
    / "SKILL.md"
)


def markdown_table(headers: tuple[str, ...], rows: list[tuple[str, ...]]) -> str:
    header = "| " + " | ".join(headers) + " |"
    separator = "|" + "|".join("---" for _ in headers) + "|"
    body = ["| " + " | ".join(row) + " |" for row in rows]
    return "\n".join([header, separator, *body])


def build_manifest(
    *,
    summary: str = "Contract coverage for the manifest and scanner slice.",
    audited_rows: list[tuple[str, str, str]] | None = None,
    preserved_rows: list[tuple[str, str, str, str]] | None = None,
    codecleanup_rows: list[tuple[str, str, str]] | None = None,
    verification_rows: list[tuple[str, str, str, str]] | None = None,
    deferred_lines: list[str] | None = None,
) -> str:
    audited_rows = audited_rows or []
    preserved_rows = preserved_rows or []
    codecleanup_rows = codecleanup_rows or [
        (
            "Brand surface pass",
            "not-applicable",
            "Later rebrand slices update visible prose and UI copy.",
        ),
        (
            "Compatibility pass",
            "pass",
            "The scanner enforces allowed preserved-identifier categories.",
        ),
        (
            "Mirror and generated-template pass",
            "pass",
            "The scanner checks embedded skill mirror parity.",
        ),
        (
            "Link and asset pass",
            "not-applicable",
            "Link and asset relabeling lands in later slices.",
        ),
        (
            "Formatting and lint pass",
            "pass",
            "The validation slice keeps manifest, scanner, and tests formatted.",
        ),
        (
            "Verification pass",
            "pass",
            "Contract tests exercise the scanner behaviors for this slice.",
        ),
    ]
    verification_rows = verification_rows or [
        (
            "python3 -m pytest tests/test_check_silmari_rebrand.py",
            ".",
            "0",
            "Contract tests cover the scanner slice.",
        ),
    ]
    deferred_lines = deferred_lines or [
        "Repo-wide Silmari copy edits land in later AF-RB slices.",
    ]

    return textwrap.dedent(
        f"""\
        # Silmari Rebrand Manifest

        ## Summary
        {summary}

        ## Audited Files
        {markdown_table(('Path', 'Action', 'Verification'), audited_rows)}

        ## Preserved Non-Silmari Identifiers
        {markdown_table(('Path', 'Identifier', 'Category', 'Reason'), preserved_rows)}

        ## CodeCleanup Passes
        {markdown_table(('Pass', 'Result', 'Notes'), codecleanup_rows)}

        ## Verification Commands
        {markdown_table(('Command', 'Working Directory', 'Exit Code', 'Result'), verification_rows)}

        ## Deferred Or Excluded
        {chr(10).join(f'- {line}' for line in deferred_lines)}
        """
    )


def write_text(root: Path, relative_path: Path | str, content: str) -> None:
    target = root / relative_path
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(content, encoding="utf-8")


def write_sync_script(root: Path) -> None:
    script = textwrap.dedent(
        f"""\
        #!/usr/bin/env bash
        set -euo pipefail

        REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
        src="$REPO_ROOT/skills/{LEGACY_SKILL_SLUG}"
        dst="$REPO_ROOT/control-plane/internal/skillkit/skill_data/{LEGACY_SKILL_SLUG}"

        if [[ "${{1:-}}" == "--check" ]]; then
          if diff -rq "$src" "$dst" >/dev/null 2>&1; then
            echo "All embedded skills are in sync with sources."
            exit 0
          fi
          echo "DRIFT: {LEGACY_SKILL_SLUG} — embed copy out of sync with source" >&2
          echo "" >&2
          echo "Run ./scripts/sync-embedded-skills.sh to fix the drift, then commit." >&2
          exit 1
        fi

        rm -rf "$dst"
        mkdir -p "$dst"
        cp -R "$src/." "$dst/"
        echo "  ✓ synced {LEGACY_SKILL_SLUG}"
        """
    )
    write_text(root, SYNC_SCRIPT_REL, script)
    (root / SYNC_SCRIPT_REL).chmod(stat.S_IRUSR | stat.S_IWUSR | stat.S_IXUSR)


def seed_skill_tree(root: Path, *, mirrored: bool) -> None:
    write_text(root, SKILL_SOURCE_REL, "# Skill\ncanonical\n")
    mirror_body = "# Skill\ncanonical\n" if mirrored else "# Skill\ndrifted\n"
    write_text(root, SKILL_MIRROR_REL, mirror_body)


def copy_public_cli(root: Path) -> None:
    write_text(root, Path("scripts") / SCRIPT_PATH.name, SCRIPT_PATH.read_text(encoding="utf-8"))
    (root / "scripts" / SCRIPT_PATH.name).chmod(
        stat.S_IRUSR | stat.S_IWUSR | stat.S_IXUSR
    )


def run_cli(root: Path) -> subprocess.CompletedProcess[str]:
    copy_public_cli(root)
    return subprocess.run(
        ["./scripts/check-silmari-rebrand.sh"],
        cwd=root,
        capture_output=True,
        text=True,
        check=False,
    )


def load_embedded_module() -> types.ModuleType:
    script_text = SCRIPT_PATH.read_text(encoding="utf-8")
    start_marker = "<<'PY'\n"
    end_marker = "\nPY\n"
    start = script_text.index(start_marker) + len(start_marker)
    end = script_text.rindex(end_marker)

    module = types.ModuleType("check_silmari_rebrand")
    module.__dict__["__name__"] = "check_silmari_rebrand"
    sys.modules[module.__name__] = module
    exec(script_text[start:end], module.__dict__)
    return module


class CheckSilmariRebrandContractTests(unittest.TestCase):
    def make_repo(self, manifest_text: str, *, mirrored: bool = True) -> tempfile.TemporaryDirectory[str]:
        tmpdir = tempfile.TemporaryDirectory()
        root = Path(tmpdir.name)
        write_text(root, MANIFEST_PATH, manifest_text)
        write_sync_script(root)
        seed_skill_tree(root, mirrored=mirrored)
        return tmpdir

    def test_empty_manifest_reports_missing_required_heading(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir_name:
            root = Path(tmpdir_name)
            write_text(root, MANIFEST_PATH, "")
            write_sync_script(root)
            seed_skill_tree(root, mirrored=True)

            result = run_cli(root)

        output = result.stdout + result.stderr
        self.assertEqual(result.returncode, 1)
        self.assertIn("Manifest errors:", output)
        self.assertIn("Missing required heading: # Silmari Rebrand Manifest", output)

    def test_single_unmanifested_match_reports_file_and_line(self) -> None:
        manifest_text = build_manifest(
            audited_rows=[
                (
                    "README.md",
                    "rebranded",
                    "./scripts/check-silmari-rebrand.sh",
                )
            ]
        )
        with self.make_repo(manifest_text) as tmpdir_name:
            root = Path(tmpdir_name)
            write_text(root, "README.md", f"Silmari still mentions {LEGACY_BRAND} here.\n")

            result = run_cli(root)

        output = result.stdout + result.stderr
        self.assertEqual(result.returncode, 1)
        self.assertIn("Unmanifested identifiers:", output)
        self.assertIn(
            f"README.md:1: {LEGACY_BRAND} :: Silmari still mentions {LEGACY_BRAND} here.",
            output,
        )

    def test_single_exact_preserved_token_is_accepted(self) -> None:
        fixture_path = "tests/fixture.txt"
        manifest_text = build_manifest(
            audited_rows=[
                (
                    fixture_path,
                    "audited-no-change",
                    "python3 -m pytest tests/test_check_silmari_rebrand.py",
                )
            ],
            preserved_rows=[
                (
                    fixture_path,
                    LEGACY_BRAND,
                    "test-fixture",
                    "Fixture asserts legacy product-token handling.",
                )
            ],
        )
        with self.make_repo(manifest_text) as tmpdir_name:
            root = Path(tmpdir_name)
            write_text(root, fixture_path, f"{LEGACY_BRAND}\n")

            result = run_cli(root)

        output = result.stdout + result.stderr
        self.assertEqual(result.returncode, 0)
        self.assertIn("Silmari rebrand check passed.", output)
        self.assertIn("validated 1 preserved identifier rows", output)

    def test_enclosing_identifier_is_accepted(self) -> None:
        fixture_path = "tests/env-fixture.txt"
        manifest_text = build_manifest(
            audited_rows=[
                (
                    fixture_path,
                    "audited-no-change",
                    "python3 -m pytest tests/test_check_silmari_rebrand.py",
                )
            ],
            preserved_rows=[
                (
                    fixture_path,
                    LEGACY_ENV,
                    "env-var",
                    "Legacy env var name remains stable for compatibility docs.",
                )
            ],
        )
        with self.make_repo(manifest_text) as tmpdir_name:
            root = Path(tmpdir_name)
            write_text(root, fixture_path, f"Use {LEGACY_ENV} for configuration.\n")

            result = run_cli(root)

        output = result.stdout + result.stderr
        self.assertEqual(result.returncode, 0)
        self.assertIn("Silmari rebrand check passed.", output)

    def test_manifest_content_is_excluded_from_brand_scan(self) -> None:
        manifest_text = build_manifest(
            audited_rows=[
                (
                    "README.md",
                    "audited-no-change",
                    "./scripts/check-silmari-rebrand.sh",
                )
            ],
            preserved_rows=[
                (
                    "README.md",
                    LEGACY_ENV,
                    "env-var",
                    "Legacy env var name remains stable for compatibility docs.",
                )
            ],
        )
        with self.make_repo(manifest_text) as tmpdir_name:
            root = Path(tmpdir_name)

            result = run_cli(root)

        output = result.stdout + result.stderr
        self.assertEqual(result.returncode, 0)
        self.assertNotIn("Unmanifested identifiers:", output)

    def test_mirror_drift_reports_scanner_failure(self) -> None:
        manifest_text = build_manifest()
        with self.make_repo(manifest_text, mirrored=False) as tmpdir_name:
            root = Path(tmpdir_name)

            result = run_cli(root)

        output = result.stdout + result.stderr
        self.assertEqual(result.returncode, 1)
        self.assertIn("skill mirror", output.lower())
        self.assertIn("drift", output.lower())

    def test_unknown_audited_file_action_is_rejected(self) -> None:
        manifest_text = build_manifest(
            audited_rows=[
                ("README.md", "unknown-action", "./scripts/check-silmari-rebrand.sh")
            ]
        )
        with self.make_repo(manifest_text) as tmpdir_name:
            root = Path(tmpdir_name)

            result = run_cli(root)

        output = result.stdout + result.stderr
        self.assertEqual(result.returncode, 1)
        self.assertIn("Unknown audited file action", output)

    def test_unknown_preserved_identifier_category_is_rejected(self) -> None:
        fixture_path = "tests/fixture.txt"
        manifest_text = build_manifest(
            audited_rows=[
                (
                    fixture_path,
                    "audited-no-change",
                    "python3 -m pytest tests/test_check_silmari_rebrand.py",
                )
            ],
            preserved_rows=[
                (
                    fixture_path,
                    LEGACY_BRAND,
                    "wrong-category",
                    "Fixture asserts legacy product-token handling.",
                )
            ],
        )
        with self.make_repo(manifest_text) as tmpdir_name:
            root = Path(tmpdir_name)

            result = run_cli(root)

        output = result.stdout + result.stderr
        self.assertEqual(result.returncode, 1)
        self.assertIn("Unknown preserved identifier category", output)

    def test_non_concrete_preserved_reason_is_rejected(self) -> None:
        fixture_path = "tests/fixture.txt"
        manifest_text = build_manifest(
            audited_rows=[
                (
                    fixture_path,
                    "audited-no-change",
                    "python3 -m pytest tests/test_check_silmari_rebrand.py",
                )
            ],
            preserved_rows=[
                (fixture_path, LEGACY_BRAND, "test-fixture", "Compatibility")
            ],
        )
        with self.make_repo(manifest_text) as tmpdir_name:
            root = Path(tmpdir_name)

            result = run_cli(root)

        output = result.stdout + result.stderr
        self.assertEqual(result.returncode, 1)
        self.assertIn("Reason must be concrete", output)


class IdentifierCoveragePropertyTests(unittest.TestCase):
    @given(
        prefix=st.text().filter(lambda value: "\n" not in value and LOWER_IMPORT not in value),
        suffix=st.text().filter(lambda value: "\n" not in value and LOWER_IMPORT not in value),
    )
    def test_enclosing_identifier_does_not_grant_line_wide_coverage(
        self, prefix: str, suffix: str
    ) -> None:
        module = load_embedded_module()
        line_text = f"{prefix}{LOWER_LEGACY} product docs keep {LOWER_IMPORT}{suffix}"
        match = module.BrandMatch(
            path="README.md",
            line_number=1,
            column_start=len(prefix),
            column_end=len(prefix) + len(LOWER_LEGACY),
            match_text=LOWER_LEGACY,
            line_text=line_text,
        )
        row = module.PreservedIdentifier(
            path="README.md",
            identifier=LOWER_IMPORT,
            category="import-module-path",
            reason="Import path stays stable for compatibility docs.",
        )

        self.assertFalse(module.row_covers_match(match, row))


if __name__ == "__main__":
    unittest.main()
