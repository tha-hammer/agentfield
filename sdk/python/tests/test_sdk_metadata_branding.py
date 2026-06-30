"""
Regression coverage for AF-RB-09 SDK branding and metadata compatibility.
"""
from __future__ import annotations

import json
import re
from pathlib import Path

import pytest

try:
    import tomllib
except ModuleNotFoundError:  # pragma: no cover
    import tomli as tomllib  # type: ignore[import-not-found]


REPO_ROOT = Path(__file__).resolve().parents[3]
PYTHON_README = REPO_ROOT / "sdk/python/README.md"
PYTHON_PYPROJECT = REPO_ROOT / "sdk/python/pyproject.toml"
TYPESCRIPT_README = REPO_ROOT / "sdk/typescript/README.md"
TYPESCRIPT_PACKAGE_JSON = REPO_ROOT / "sdk/typescript/package.json"
GO_MOD = REPO_ROOT / "sdk/go/go.mod"
GO_README = REPO_ROOT / "sdk/go/README.md"
MANIFEST = REPO_ROOT / "docs/silmari-rebrand-manifest.md"

CODE_FENCE_RE = re.compile(r"```[^\n]*\n(.*?)```", re.DOTALL)
EXPECTED_AUDITED_FILES = [
    "sdk/python/README.md",
    "sdk/python/pyproject.toml",
    "sdk/typescript/README.md",
    "sdk/typescript/package.json",
    "sdk/go/go.mod",
    "sdk/go/README.md",
    "sdk/python/tests/test_sdk_metadata_branding.py",
    "docs/silmari-rebrand-manifest.md",
]
EXPECTED_MANIFEST_AUDIT_ROWS = [
    (
        "sdk/go/go.mod",
        "audited-no-change",
        "direct Go module-path grep; targeted SDK metadata branding pytest",
    ),
    (
        "sdk/python/tests/test_sdk_metadata_branding.py",
        "rebranded-with-preserved-identifiers",
        "targeted SDK metadata branding pytest",
    ),
    (
        "docs/silmari-rebrand-manifest.md",
        "rebranded-with-preserved-identifiers",
        "manual review; targeted SDK metadata branding pytest",
    ),
]
EXPECTED_TEST_FIXTURE_PRESERVATION_ROWS = [
    (
        "sdk/python/tests/test_sdk_metadata_branding.py",
        "AgentField",
        "test-fixture",
    ),
    (
        "sdk/python/tests/test_sdk_metadata_branding.py",
        "agentfield",
        "test-fixture",
    ),
    (
        "sdk/python/tests/test_sdk_metadata_branding.py",
        "@agentfield/sdk",
        "test-fixture",
    ),
    (
        "sdk/python/tests/test_sdk_metadata_branding.py",
        "https://github.com/Agent-Field/agentfield",
        "test-fixture",
    ),
    (
        "sdk/python/tests/test_sdk_metadata_branding.py",
        "github.com/Agent-Field/agentfield/sdk/go",
        "test-fixture",
    ),
]


def _read_text(path: Path) -> str:
    return path.read_text(encoding="utf-8")


def _load_json(path: Path) -> dict[str, object]:
    return json.loads(_read_text(path))


def _load_toml(path: Path) -> dict[str, object]:
    return tomllib.loads(_read_text(path))


def _extract_fenced_code_blocks(markdown: str | None) -> list[str]:
    if not markdown:
        return []
    return [match.group(1) for match in CODE_FENCE_RE.finditer(markdown)]


def _strip_fenced_code_blocks(markdown: str | None) -> str:
    if not markdown:
        return ""
    return CODE_FENCE_RE.sub("", markdown)


def _extract_markdown_table_rows(markdown: str, heading: str) -> list[tuple[str, ...]]:
    section_match = re.search(
        rf"^{re.escape(heading)}\n(.*?)(?=^## |\Z)",
        markdown,
        re.MULTILINE | re.DOTALL,
    )
    assert section_match is not None, f"Missing section: {heading}"

    rows: list[tuple[str, ...]] = []
    for line in section_match.group(1).splitlines():
        stripped = line.strip()
        if not stripped.startswith("|"):
            continue

        cells = tuple(cell.strip() for cell in stripped.strip("|").split("|"))
        if cells and cells[0] == "Path":
            continue
        if cells and all(re.fullmatch(r":?-{3,}:?", cell) for cell in cells):
            continue
        rows.append(cells)

    return rows


def _assert_python_metadata_contract(pyproject: dict[str, object]) -> None:
    project = pyproject["project"]
    urls = project["urls"]

    assert project["name"] == "agentfield"
    assert project["description"].startswith("Silmari Python SDK")
    assert project["authors"] == [{"name": "Silmari Maintainers"}]
    assert "silmari" in project["keywords"]
    assert "agentfield" in project["keywords"]
    assert urls["Homepage"] == "https://github.com/Agent-Field/agentfield"
    assert (
        urls["Documentation"]
        == "https://github.com/Agent-Field/agentfield/tree/main/docs"
    )
    assert urls["Issues"] == "https://github.com/Agent-Field/agentfield/issues"


def _assert_typescript_metadata_contract(package_json: dict[str, object]) -> None:
    assert package_json["name"] == "@agentfield/sdk"
    assert package_json["description"] == "Silmari TypeScript SDK"
    assert package_json["author"] == "Silmari Team"
    assert "silmari" in package_json["keywords"]
    assert "agentfield" in package_json["keywords"]
    assert package_json["repository"] == {
        "type": "git",
        "url": "https://github.com/Agent-Field/agentfield.git",
        "directory": "sdk/typescript",
    }
    assert (
        package_json["homepage"]
        == "https://github.com/Agent-Field/agentfield/tree/main/sdk/typescript"
    )
    assert package_json["bugs"] == {
        "url": "https://github.com/Agent-Field/agentfield/issues"
    }


class TestBrandingHelpers:
    @pytest.mark.parametrize("markdown", [None, ""])
    def test_extract_fenced_code_blocks_handles_none_and_empty(self, markdown):
        assert _extract_fenced_code_blocks(markdown) == []

    @pytest.mark.parametrize("markdown", [None, ""])
    def test_strip_fenced_code_blocks_handles_none_and_empty(self, markdown):
        assert _strip_fenced_code_blocks(markdown) == ""

    def test_strip_fenced_code_blocks_keeps_unclosed_boundary_case(self):
        markdown = "Intro\n```python\nfrom agentfield import Agent\n"
        assert _strip_fenced_code_blocks(markdown) == markdown

    def test_extract_markdown_table_rows_raises_for_missing_heading(self):
        with pytest.raises(AssertionError, match="Missing section: ## Audited Files"):
            _extract_markdown_table_rows(
                "## Summary\nNo table here.\n", "## Audited Files"
            )

    def test_extract_markdown_table_rows_ignores_headers_separators_and_prose(self):
        markdown = """
## Audited Files
Introductory prose that should be ignored.
| Path | Action | Verification |
|---|---|---|
| sdk/python/README.md | rebranded | visual review |
Trailing prose that should also be ignored.
## Preserved Non-Silmari Identifiers
| Path | Identifier | Category | Reason |
|---|---|---|---|
"""

        assert _extract_markdown_table_rows(markdown, "## Audited Files") == [
            ("sdk/python/README.md", "rebranded", "visual review")
        ]

    def test_load_json_raises_for_invalid_package_metadata(self, tmp_path: Path):
        package_json = tmp_path / "package.json"
        package_json.write_text('{"name": ', encoding="utf-8")

        with pytest.raises(json.JSONDecodeError):
            _load_json(package_json)

    def test_load_toml_raises_for_invalid_package_metadata(self, tmp_path: Path):
        pyproject = tmp_path / "pyproject.toml"
        pyproject.write_text("[project\nname = 'agentfield'", encoding="utf-8")

        with pytest.raises(tomllib.TOMLDecodeError):
            _load_toml(pyproject)


class TestSdkPackageMetadata:
    def test_go_module_declares_published_compatibility_path(self):
        assert _read_text(GO_MOD).splitlines()[0] == (
            "module github.com/Agent-Field/agentfield/sdk/go"
        )

    def test_python_metadata_rebrands_safe_fields(self):
        _assert_python_metadata_contract(_load_toml(PYTHON_PYPROJECT))

    def test_python_metadata_preserves_repository_urls(self):
        _assert_python_metadata_contract(_load_toml(PYTHON_PYPROJECT))

    def test_python_metadata_rejects_empty_project_table(self):
        with pytest.raises(KeyError):
            _assert_python_metadata_contract({})

    def test_python_metadata_rejects_rebranded_package_name(self):
        pyproject = _load_toml(PYTHON_PYPROJECT)
        pyproject["project"]["name"] = "silmari"

        with pytest.raises(AssertionError):
            _assert_python_metadata_contract(pyproject)

    def test_python_metadata_rejects_homepage_url_drift(self):
        pyproject = _load_toml(PYTHON_PYPROJECT)
        pyproject["project"]["urls"]["Homepage"] = "https://github.com/Silmari/silmari"

        with pytest.raises(AssertionError):
            _assert_python_metadata_contract(pyproject)

    def test_python_metadata_rejects_missing_legacy_keyword(self):
        pyproject = _load_toml(PYTHON_PYPROJECT)
        pyproject["project"]["keywords"] = ["silmari", "sdk", "agents"]

        with pytest.raises(AssertionError):
            _assert_python_metadata_contract(pyproject)

    def test_python_metadata_rejects_missing_silmari_keyword(self):
        pyproject = _load_toml(PYTHON_PYPROJECT)
        pyproject["project"]["keywords"] = ["agentfield", "sdk", "agents"]

        with pytest.raises(AssertionError):
            _assert_python_metadata_contract(pyproject)

    def test_python_metadata_rejects_author_drift(self):
        pyproject = _load_toml(PYTHON_PYPROJECT)
        pyproject["project"]["authors"] = [{"name": "AgentField Maintainers"}]

        with pytest.raises(AssertionError):
            _assert_python_metadata_contract(pyproject)

    def test_typescript_metadata_rebrands_safe_fields(self):
        _assert_typescript_metadata_contract(_load_json(TYPESCRIPT_PACKAGE_JSON))

    def test_typescript_metadata_preserves_repository_urls(self):
        package_json = _load_json(TYPESCRIPT_PACKAGE_JSON)

        _assert_typescript_metadata_contract(package_json)

    def test_typescript_metadata_rejects_empty_metadata(self):
        with pytest.raises(KeyError):
            _assert_typescript_metadata_contract({})

    def test_typescript_metadata_rejects_repository_url_drift(self):
        package_json = _load_json(TYPESCRIPT_PACKAGE_JSON)
        package_json["repository"] = {
            "type": "git",
            "url": "https://github.com/Silmari/silmari.git",
            "directory": "sdk/typescript",
        }

        with pytest.raises(AssertionError):
            _assert_typescript_metadata_contract(package_json)

    def test_typescript_metadata_rejects_rebranded_package_name(self):
        package_json = _load_json(TYPESCRIPT_PACKAGE_JSON)
        package_json["name"] = "@silmari/sdk"

        with pytest.raises(AssertionError):
            _assert_typescript_metadata_contract(package_json)

    def test_typescript_metadata_rejects_missing_silmari_keyword(self):
        package_json = _load_json(TYPESCRIPT_PACKAGE_JSON)
        package_json["keywords"] = [
            "agentfield",
            "ai-agents",
            "llm",
            "multi-agent",
            "orchestration",
        ]

        with pytest.raises(AssertionError):
            _assert_typescript_metadata_contract(package_json)

    def test_typescript_metadata_rejects_missing_legacy_keyword(self):
        package_json = _load_json(TYPESCRIPT_PACKAGE_JSON)
        package_json["keywords"] = [
            "silmari",
            "ai-agents",
            "llm",
            "multi-agent",
            "orchestration",
        ]

        with pytest.raises(AssertionError):
            _assert_typescript_metadata_contract(package_json)

    def test_typescript_metadata_rejects_author_drift(self):
        package_json = _load_json(TYPESCRIPT_PACKAGE_JSON)
        package_json["author"] = "AgentField Team"

        with pytest.raises(AssertionError):
            _assert_typescript_metadata_contract(package_json)


class TestSdkReadmeBranding:
    @pytest.mark.parametrize(
        ("path", "heading"),
        [
            (PYTHON_README, "# Silmari Python SDK"),
            (TYPESCRIPT_README, "# Silmari TypeScript SDK"),
            (GO_README, "# Silmari Go SDK"),
        ],
    )
    def test_sdk_readmes_rebrand_visible_prose(self, path: Path, heading: str):
        markdown = _read_text(path)
        prose = _strip_fenced_code_blocks(markdown)

        assert heading in markdown
        assert "Silmari" in prose
        assert re.search(r"\bAgentField\b", prose) is None

    @pytest.mark.parametrize(
        ("path", "snippet"),
        [
            (PYTHON_README, "pip install agentfield"),
            (PYTHON_README, "from agentfield import Agent"),
            (TYPESCRIPT_README, "npm install @agentfield/sdk"),
            (TYPESCRIPT_README, "import { Agent } from '@agentfield/sdk';"),
            (GO_README, "go get github.com/Agent-Field/agentfield/sdk/go"),
            (
                GO_README,
                'agentfieldagent "github.com/Agent-Field/agentfield/sdk/go/agent"',
            ),
        ],
    )
    def test_sdk_readmes_preserve_install_and_import_examples(
        self, path: Path, snippet: str
    ):
        markdown = _read_text(path)
        code_blocks = "\n".join(_extract_fenced_code_blocks(markdown))

        assert snippet in code_blocks

    @pytest.mark.parametrize(
        ("path", "forbidden_snippet"),
        [
            (PYTHON_README, "pip install silmari"),
            (PYTHON_README, "from silmari import Agent"),
            (TYPESCRIPT_README, "npm install @silmari/sdk"),
            (TYPESCRIPT_README, "import { Agent } from '@silmari/sdk';"),
            (GO_README, "go get github.com/Silmari/silmari/sdk/go"),
        ],
    )
    def test_sdk_readmes_do_not_rebrand_package_or_module_identifiers(
        self, path: Path, forbidden_snippet: str
    ):
        markdown = _read_text(path)

        assert forbidden_snippet not in markdown


class TestCompatibilityManifest:
    def test_manifest_audited_files_rows_have_expected_shape(self):
        audited_rows = _extract_markdown_table_rows(
            _read_text(MANIFEST), "## Audited Files"
        )

        assert audited_rows
        assert all(len(row) == 3 for row in audited_rows)

    @pytest.mark.parametrize("path", EXPECTED_AUDITED_FILES)
    def test_manifest_audits_each_sdk_slice_file(self, path: str):
        audited_rows = _extract_markdown_table_rows(
            _read_text(MANIFEST), "## Audited Files"
        )

        assert any(len(row) == 3 and row[0] == path for row in audited_rows)

    @pytest.mark.parametrize(
        ("path", "action", "verification"), EXPECTED_MANIFEST_AUDIT_ROWS
    )
    def test_manifest_preserves_expected_audit_actions(
        self, path: str, action: str, verification: str
    ):
        manifest_text = _read_text(MANIFEST)

        assert f"| {path} | {action} | {verification} |" in manifest_text

    @pytest.mark.parametrize(
        ("source_path", "identifier", "category"),
        [
            ("sdk/python/pyproject.toml", "agentfield", "package-name"),
            ("sdk/typescript/package.json", "@agentfield/sdk", "package-name"),
            (
                "sdk/go/README.md",
                "github.com/Agent-Field/agentfield/sdk/go",
                "go-module-path",
            ),
            (
                "sdk/go/go.mod",
                "module github.com/Agent-Field/agentfield/sdk/go",
                "go-module-path",
            ),
            (
                "sdk/python/tests/test_sdk_metadata_branding.py",
                "AgentField",
                "test-fixture",
            ),
            (
                "sdk/python/tests/test_sdk_metadata_branding.py",
                "agentfield",
                "test-fixture",
            ),
            (
                "sdk/python/tests/test_sdk_metadata_branding.py",
                "@agentfield/sdk",
                "test-fixture",
            ),
            (
                "sdk/python/tests/test_sdk_metadata_branding.py",
                "https://github.com/Agent-Field/agentfield",
                "test-fixture",
            ),
            (
                "sdk/python/tests/test_sdk_metadata_branding.py",
                "github.com/Agent-Field/agentfield/sdk/go",
                "test-fixture",
            ),
        ],
    )
    def test_manifest_records_preserved_sdk_identifiers(
        self, source_path: str, identifier: str, category: str
    ):
        manifest_text = _read_text(MANIFEST)
        source_text = _read_text(REPO_ROOT / source_path)

        assert identifier in source_text
        assert f"| {source_path} | {identifier} | {category} |" in manifest_text

    @pytest.mark.parametrize(
        ("source_path", "identifier", "category"),
        EXPECTED_TEST_FIXTURE_PRESERVATION_ROWS,
    )
    def test_manifest_records_all_python_test_fixture_preservations(
        self, source_path: str, identifier: str, category: str
    ):
        manifest_text = _read_text(MANIFEST)
        source_text = _read_text(REPO_ROOT / source_path)

        assert identifier in source_text
        assert f"| {source_path} | {identifier} | {category} |" in manifest_text

    def test_manifest_preserved_identifier_rows_have_expected_shape(self):
        preserved_rows = _extract_markdown_table_rows(
            _read_text(MANIFEST), "## Preserved Non-Silmari Identifiers"
        )

        assert preserved_rows
        assert all(len(row) == 4 for row in preserved_rows)
