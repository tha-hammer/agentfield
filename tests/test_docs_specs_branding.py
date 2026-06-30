from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
import fnmatch
import re

import pytest


ROOT = Path(__file__).resolve().parents[1]
MANIFEST_PATH = ROOT / "docs" / "silmari-rebrand-manifest.md"

SCOPED_PATHS = (
    "docs",
    "specs",
    "control-plane/README.md",
    "control-plane/scripts/README.md",
    "control-plane/tools/perf/README.md",
    "control-plane/migrations/README.md",
    "tests/functional/README.md",
    "tests/functional/docker/LOG_DEMO.md",
)

REQUIRED_HEADINGS = (
    "# Silmari Rebrand Manifest",
    "## Summary",
    "## Audited Files",
    "## Preserved Non-Silmari Identifiers",
    "## CodeCleanup Passes",
    "## Verification Commands",
    "## Deferred Or Excluded",
)

ALLOWED_ACTIONS = {
    "rebranded",
    "rebranded-with-preserved-identifiers",
    "audited-no-change",
    "excluded-historical",
    "excluded-runtime-compatibility",
}

ALLOWED_CATEGORIES = {
    "package-name",
    "import-module-path",
    "go-module-path",
    "api-path",
    "json-field",
    "env-var",
    "yaml-config-path",
    "cli-command",
    "docker-image-or-volume",
    "helm-chart-or-k8s-name",
    "skill-or-repo-slug",
    "published-link-target",
    "historical-record",
    "test-fixture",
}

OLD_BRAND_RE = re.compile(
    r"agentfield\.ai|Agent-Field|AgentField|AGENTFIELD|agentfield|AgentPlane|agentplane|Agent Plane|agent plane"
)
IDENTIFIER_CHARS = frozenset("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789._/@:-~")


@dataclass(frozen=True)
class AuditedFile:
    path: str
    action: str
    verification: str


@dataclass(frozen=True)
class PreservedIdentifier:
    path: str
    identifier: str
    category: str
    reason: str


@dataclass(frozen=True)
class ManifestData:
    audited: dict[str, AuditedFile]
    preserved: list[PreservedIdentifier]


def _extract_section(text: str, heading: str) -> str:
    start = text.find(heading)
    if start == -1:
        raise ValueError(f"missing required heading: {heading}")
    end = len(text)
    for candidate in REQUIRED_HEADINGS:
        if candidate == heading:
            continue
        candidate_start = text.find(candidate, start + len(heading))
        if candidate_start != -1:
            end = min(end, candidate_start)
    return text[start:end]


def _parse_table(section: str) -> list[list[str]]:
    rows = [line.strip() for line in section.splitlines() if line.strip().startswith("|")]
    if len(rows) < 2:
        raise ValueError("manifest table is missing header rows")
    data_rows = []
    for row in rows[2:]:
        parts = [part.strip() for part in row.strip("|").split("|")]
        data_rows.append(parts)
    return data_rows


def parse_manifest(text: str) -> ManifestData:
    for heading in REQUIRED_HEADINGS:
        if heading not in text:
            raise ValueError(f"missing required heading: {heading}")

    audited_rows = _parse_table(_extract_section(text, "## Audited Files"))
    preserved_rows = _parse_table(_extract_section(text, "## Preserved Non-Silmari Identifiers"))

    audited: dict[str, AuditedFile] = {}
    for row in audited_rows:
        if len(row) != 3:
            raise ValueError(f"audited files row must have 3 columns: {row}")
        path, action, verification = row
        if action not in ALLOWED_ACTIONS:
            raise ValueError(f"invalid audited action for {path}: {action}")
        audited[path] = AuditedFile(path=path, action=action, verification=verification)

    preserved: list[PreservedIdentifier] = []
    for row in preserved_rows:
        if len(row) != 4:
            raise ValueError(f"preserved identifiers row must have 4 columns: {row}")
        path, identifier, category, reason = row
        if category not in ALLOWED_CATEGORIES:
            raise ValueError(f"invalid preserved category for {path}: {category}")
        if not reason:
            raise ValueError(f"missing preserved identifier reason for {path}: {identifier}")
        preserved.append(
            PreservedIdentifier(
                path=path,
                identifier=identifier.strip("`"),
                category=category,
                reason=reason,
            )
        )

    return ManifestData(audited=audited, preserved=preserved)


def _identifier_candidates(line: str, match_start: int, match_end: int) -> tuple[str, ...]:
    left = match_start
    while left > 0 and line[left - 1] in IDENTIFIER_CHARS:
        left -= 1

    right = match_end
    while right < len(line) and line[right] in IDENTIFIER_CHARS:
        right += 1

    raw = line[left:right]
    if not raw:
        return ()

    candidates: list[str] = []
    rel_start = match_start - left
    rel_end = match_end - left

    slash_boundaries = [0]
    slash_boundaries.extend(index + 1 for index, char in enumerate(raw[:rel_start]) if char == "/")
    slash_boundaries = sorted(set(slash_boundaries))

    for boundary in slash_boundaries:
        next_separator = raw.find("/", rel_end)
        candidate_values = [raw[boundary:]]
        if next_separator != -1:
            candidate_values.append(raw[boundary:next_separator])

        for candidate in candidate_values:
            if not candidate:
                continue
            if candidate not in candidates:
                candidates.append(candidate)
            if candidate.startswith("./"):
                trimmed = candidate[2:]
                if trimmed not in candidates:
                    candidates.append(trimmed)
            if candidate.startswith("/."):
                trimmed = candidate[1:]
                if trimmed not in candidates:
                    candidates.append(trimmed)
            stripped = candidate.rstrip(".,;:")
            if stripped and stripped != candidate and stripped not in candidates:
                candidates.append(stripped)
            if stripped.startswith("./"):
                trimmed = stripped[2:]
                if trimmed not in candidates:
                    candidates.append(trimmed)
            if stripped.startswith("/."):
                trimmed = stripped[1:]
                if trimmed not in candidates:
                    candidates.append(trimmed)

    segment_end = raw.find("/", rel_end)
    if segment_end == -1:
        segment_end = len(raw)
    segment_start = raw.rfind("/", 0, rel_start) + 1
    for candidate in (raw[segment_start:segment_end],):
        if not candidate:
            continue
        if candidate not in candidates:
            candidates.append(candidate)
        stripped = candidate.rstrip(".,;:")
        if stripped and stripped != candidate and stripped not in candidates:
            candidates.append(stripped)

    return tuple(candidates)


def _identifier_covers_match(identifier: str, line: str, match_start: int, match_end: int) -> bool:
    candidates = _identifier_candidates(line, match_start, match_end)
    if any(wildcard in identifier for wildcard in "*?[]"):
        return any(fnmatch.fnmatch(candidate, identifier) for candidate in candidates)
    return identifier in candidates


def _collect_scoped_markdown_files(root: Path) -> list[Path]:
    files: list[Path] = []
    for scoped_path in SCOPED_PATHS:
        path = root / scoped_path
        if path.is_dir():
            files.extend(sorted(path.rglob("*.md")))
        elif path.exists():
            files.append(path)
    return [path for path in files if path != MANIFEST_PATH]


def validate_scope(root: Path, files: list[Path], manifest_text: str) -> list[str]:
    try:
        manifest = parse_manifest(manifest_text)
    except ValueError as exc:
        return [str(exc)]

    preserved_by_path: dict[str, list[PreservedIdentifier]] = {}
    for row in manifest.preserved:
        preserved_by_path.setdefault(row.path, []).append(row)

    errors: list[str] = []
    for path in files:
        rel_path = path.relative_to(root).as_posix()
        text = path.read_text(encoding="utf-8")
        found_brand = False

        for line_number, line in enumerate(text.splitlines(), start=1):
            for match in OLD_BRAND_RE.finditer(line):
                found_brand = True
                rows = preserved_by_path.get(rel_path, [])
                if any(_identifier_covers_match(row.identifier, line, match.start(), match.end()) for row in rows):
                    continue
                errors.append(
                    f"{rel_path}:{line_number}: unmanifested identifier '{match.group(0)}' in {line.strip()}"
                )

        if found_brand and rel_path not in manifest.audited:
            errors.append(f"{rel_path}: file contains preserved identifiers but is missing from Audited Files")

    return errors


def _make_manifest(audited_rows: list[str], preserved_rows: list[str]) -> str:
    return "\n".join(
        [
            "# Silmari Rebrand Manifest",
            "",
            "## Summary",
            "Scoped test manifest.",
            "",
            "## Audited Files",
            "| Path | Action | Verification |",
            "|------|--------|--------------|",
            *audited_rows,
            "",
            "## Preserved Non-Silmari Identifiers",
            "| Path | Identifier | Category | Reason |",
            "|------|------------|----------|--------|",
            *preserved_rows,
            "",
            "## CodeCleanup Passes",
            "| Pass | Result | Notes |",
            "|------|--------|-------|",
            "| Brand surface pass | pass | Scoped test manifest. |",
            "",
            "## Verification Commands",
            "| Command | Working Directory | Exit Code | Result |",
            "|---------|-------------------|-----------|--------|",
            "| `python3 -m pytest tests/test_docs_specs_branding.py` | `.` | 0 | Scoped tests. |",
            "",
            "## Deferred Or Excluded",
            "- None.",
        ]
    )


def test_validate_scope_accepts_empty_docs(tmp_path: Path) -> None:
    doc = tmp_path / "docs" / "empty.md"
    doc.parent.mkdir(parents=True)
    doc.write_text("", encoding="utf-8")

    manifest = _make_manifest([], [])
    assert validate_scope(tmp_path, [doc], manifest) == []


def test_validate_scope_accepts_empty_file_list() -> None:
    manifest = _make_manifest([], [])
    assert validate_scope(ROOT, [], manifest) == []


def test_validate_scope_accepts_single_preserved_env_var(tmp_path: Path) -> None:
    doc = tmp_path / "docs" / "env.md"
    doc.parent.mkdir(parents=True)
    doc.write_text("Set AGENTFIELD_PORT before starting the control plane.\n", encoding="utf-8")

    manifest = _make_manifest(
        ["| docs/env.md | audited-no-change | `pytest` |"],
        ["| docs/env.md | `AGENTFIELD*` | env-var | Env vars remain stable compatibility surfaces in this file. |"],
    )

    assert validate_scope(tmp_path, [doc], manifest) == []


def test_validate_scope_accepts_many_preserved_snippets(tmp_path: Path) -> None:
    doc = tmp_path / "specs" / "compat.md"
    doc.parent.mkdir(parents=True)
    doc.write_text(
        "\n".join(
            [
                "Use config/agentfield.yaml by default.",
                "Legacy install URL: https://agentfield.ai/install.sh",
                "Keep /api/v1/did/agentfield-server available for compatibility.",
            ]
        )
        + "\n",
        encoding="utf-8",
    )

    manifest = _make_manifest(
        ["| specs/compat.md | audited-no-change | `pytest` |"],
        [
            "| specs/compat.md | `config/agentfield.yaml` | yaml-config-path | The default YAML config path remains stable in this file. |",
            "| specs/compat.md | `agentfield.ai*` | published-link-target | Published URLs remain on the legacy domain in this file. |",
            "| specs/compat.md | `/api/v1/did/agentfield-server` | api-path | The DID endpoint path remains a runtime compatibility surface. |",
        ],
    )

    assert validate_scope(tmp_path, [doc], manifest) == []


def test_validate_scope_flags_agentplane_prose_without_historical_allowance(tmp_path: Path) -> None:
    doc = tmp_path / "specs" / "ui.md"
    doc.parent.mkdir(parents=True)
    doc.write_text("AgentPlane UI API worklist\n", encoding="utf-8")

    manifest = _make_manifest(["| specs/ui.md | audited-no-change | `pytest` |"], [])
    errors = validate_scope(tmp_path, [doc], manifest)

    assert errors
    assert "AgentPlane" in errors[0]


def test_validate_scope_accepts_manifested_agentplane_historical_reference(tmp_path: Path) -> None:
    doc = tmp_path / "specs" / "ui.md"
    doc.parent.mkdir(parents=True)
    doc.write_text("See agentplane-ui-api-worklist.md for the historical planning note.\n", encoding="utf-8")

    manifest = _make_manifest(
        ["| specs/ui.md | audited-no-change | `pytest` |"],
        [
            "| specs/ui.md | `agentplane-ui-api-worklist.md` | historical-record | The planning note keeps its historical filename in this file. |"
        ],
    )

    assert validate_scope(tmp_path, [doc], manifest) == []


def test_validate_scope_flags_old_brand_tokens_in_heading_table_diagram_and_prose(tmp_path: Path) -> None:
    doc = tmp_path / "docs" / "surfaces.md"
    doc.parent.mkdir(parents=True)
    doc.write_text(
        "\n".join(
            [
                "# AgentField Overview",
                "AgentPlane remains in this introduction.",
                "",
                "| Surface | AgentField |",
                "|---------|------------|",
                "",
                "```mermaid",
                "flowchart LR",
                '  A["agentplane dashboard"] --> B["ok"]',
                "```",
            ]
        )
        + "\n",
        encoding="utf-8",
    )

    manifest = _make_manifest(["| docs/surfaces.md | audited-no-change | `pytest` |"], [])
    errors = validate_scope(tmp_path, [doc], manifest)

    assert errors == [
        "docs/surfaces.md:1: unmanifested identifier 'AgentField' in # AgentField Overview",
        "docs/surfaces.md:2: unmanifested identifier 'AgentPlane' in AgentPlane remains in this introduction.",
        "docs/surfaces.md:4: unmanifested identifier 'AgentField' in | Surface | AgentField |",
        'docs/surfaces.md:9: unmanifested identifier \'agentplane\' in A["agentplane dashboard"] --> B["ok"]',
    ]


def test_validate_scope_rejects_malformed_manifest(tmp_path: Path) -> None:
    doc = tmp_path / "docs" / "env.md"
    doc.parent.mkdir(parents=True)
    doc.write_text("AGENTFIELD_PORT\n", encoding="utf-8")

    bad_manifest = _make_manifest(
        ["| docs/env.md | audited-no-change | `pytest` |"],
        ["| docs/env.md | `AGENTFIELD*` | wrong-category | Missing category enum. |"],
    )

    errors = validate_scope(tmp_path, [doc], bad_manifest)
    assert errors == ["invalid preserved category for docs/env.md: wrong-category"]


def test_validate_scope_rejects_invalid_audited_action(tmp_path: Path) -> None:
    doc = tmp_path / "docs" / "env.md"
    doc.parent.mkdir(parents=True)
    doc.write_text("AGENTFIELD_PORT\n", encoding="utf-8")

    bad_manifest = _make_manifest(
        ["| docs/env.md | rebranded-for-real | `pytest` |"],
        ["| docs/env.md | `AGENTFIELD*` | env-var | Env vars remain stable compatibility surfaces in this file. |"],
    )

    errors = validate_scope(tmp_path, [doc], bad_manifest)
    assert errors == ["invalid audited action for docs/env.md: rebranded-for-real"]


def test_validate_scope_reports_missing_preserved_identifier_reason(tmp_path: Path) -> None:
    doc = tmp_path / "docs" / "env.md"
    doc.parent.mkdir(parents=True)
    doc.write_text("AGENTFIELD_PORT\n", encoding="utf-8")

    bad_manifest = _make_manifest(
        ["| docs/env.md | audited-no-change | `pytest` |"],
        ["| docs/env.md | `AGENTFIELD*` | env-var | |"],
    )

    errors = validate_scope(tmp_path, [doc], bad_manifest)
    assert errors == ["missing preserved identifier reason for docs/env.md: `AGENTFIELD*`"]


def test_validate_scope_reports_missing_required_heading(tmp_path: Path) -> None:
    doc = tmp_path / "docs" / "env.md"
    doc.parent.mkdir(parents=True)
    doc.write_text("AGENTFIELD_PORT\n", encoding="utf-8")

    bad_manifest = _make_manifest([], []).replace("## Verification Commands\n", "")

    errors = validate_scope(tmp_path, [doc], bad_manifest)
    assert errors == ["missing required heading: ## Verification Commands"]


def test_validate_scope_flags_preserved_identifier_without_audited_file_entry(tmp_path: Path) -> None:
    doc = tmp_path / "docs" / "env.md"
    doc.parent.mkdir(parents=True)
    doc.write_text("Set AGENTFIELD_PORT before starting the control plane.\n", encoding="utf-8")

    manifest = _make_manifest(
        [],
        ["| docs/env.md | `AGENTFIELD*` | env-var | Env vars remain stable compatibility surfaces in this file. |"],
    )

    errors = validate_scope(tmp_path, [doc], manifest)
    assert errors == ["docs/env.md: file contains preserved identifiers but is missing from Audited Files"]


def test_validate_scope_accepts_embedded_preserved_identifiers(tmp_path: Path) -> None:
    doc = tmp_path / "docs" / "embedded.md"
    doc.parent.mkdir(parents=True)
    doc.write_text(
        "\n".join(
            [
                "Harness outputs live in .agentfield_output.json and .agentfield_schema.json.",
                "SDK adapters still expose AgentFieldHandler and AgentFieldClient for compatibility.",
            ]
        )
        + "\n",
        encoding="utf-8",
    )

    manifest = _make_manifest(
        ["| docs/embedded.md | audited-no-change | `pytest` |"],
        [
            "| docs/embedded.md | `.agentfield_output.json` | test-fixture | The harness output filename remains stable for downstream tooling in this file. |",
            "| docs/embedded.md | `.agentfield_schema.json` | test-fixture | The harness schema filename remains stable for downstream tooling in this file. |",
            "| docs/embedded.md | `AgentFieldHandler` | import-module-path | The public Python SDK handler class name remains stable in this file. |",
            "| docs/embedded.md | `AgentFieldClient` | import-module-path | The public Python SDK client class name remains stable in this file. |",
        ],
    )

    assert validate_scope(tmp_path, [doc], manifest) == []


def test_validate_scope_flags_unmanifested_embedded_identifier(tmp_path: Path) -> None:
    doc = tmp_path / "docs" / "embedded.md"
    doc.parent.mkdir(parents=True)
    doc.write_text("Harness outputs live in .agentfield_output.json.\n", encoding="utf-8")

    manifest = _make_manifest(
        ["| docs/embedded.md | audited-no-change | `pytest` |"],
        [],
    )

    errors = validate_scope(tmp_path, [doc], manifest)
    assert errors == [
        "docs/embedded.md:1: unmanifested identifier 'agentfield' in Harness outputs live in .agentfield_output.json."
    ]


@pytest.mark.parametrize(
    ("identifier", "line", "expected"),
    [
        (".agentfield_output.json", "Harness outputs live in .agentfield_output.json.", True),
        ("agentfield.ai*", "Legacy install URL: https://agentfield.ai/install.sh", True),
        ("AGENTFIELD*", "Set AGENTFIELD_CONFIG_FILE before booting the control plane.", True),
        ("agentfield", "Use agentfield-server for compatibility.", False),
    ],
)
def test_identifier_covers_match_respects_boundaries(identifier: str, line: str, expected: bool) -> None:
    match = OLD_BRAND_RE.search(line)

    assert match is not None
    assert _identifier_covers_match(identifier, line, match.start(), match.end()) is expected


def test_repository_configuration_docs_keep_yaml_and_env_var_compatibility_surfaces() -> None:
    compatibility_targets = {
        ROOT / "docs" / "ENVIRONMENT_VARIABLES.md": "first-class legacy-compatible configuration surfaces",
        ROOT / "control-plane" / "README.md": "peer configuration surfaces",
        ROOT / "specs" / "architecture-overview.md": "legacy-compatible `AGENTFIELD_CONFIG_FILE`",
        ROOT / "specs" / "control-plane.md": "legacy-compatible `AGENTFIELD_CONFIG_FILE`",
        ROOT / "specs" / "deployment.md": "legacy env var name stable",
    }

    for path, phrase in compatibility_targets.items():
        text = path.read_text(encoding="utf-8")

        assert "config/agentfield.yaml" in text, path.as_posix()
        assert "AGENTFIELD_CONFIG_FILE" in text, path.as_posix()
        assert phrase in text, path.as_posix()


def test_repository_docs_and_specs_use_silmari_visible_brand() -> None:
    visible_brand_targets = {
        ROOT / "docs" / "ARCHITECTURE.md": "# Silmari Architecture",
        ROOT / "specs" / "README.md": "# Silmari — Specifications & Architecture",
        ROOT / "specs" / "agentplane-ui-api-spec.md": "# Silmari UI API Specification",
        ROOT / "control-plane" / "README.md": "# Silmari Control Plane",
        ROOT / "control-plane" / "tools" / "perf" / "README.md": "containerised Silmari server",
        ROOT / "tests" / "functional" / "README.md": "# Silmari Functional Tests",
    }

    for path, phrase in visible_brand_targets.items():
        text = path.read_text(encoding="utf-8")
        assert phrase in text, path.as_posix()


def test_agentplane_ui_api_spec_only_keeps_manifested_historical_reference() -> None:
    spec_path = ROOT / "specs" / "agentplane-ui-api-spec.md"
    spec_text = spec_path.read_text(encoding="utf-8")
    manifest = parse_manifest(MANIFEST_PATH.read_text(encoding="utf-8"))

    assert spec_text.startswith("# Silmari UI API Specification")
    assert "AgentPlane" not in spec_text
    assert spec_text.count("agentplane") == 1
    assert "historical planning note" in spec_text
    assert "agentplane-ui-api-worklist.md" in spec_text
    assert any(
        row.path == "specs/agentplane-ui-api-spec.md"
        and row.identifier == "agentplane-ui-api-worklist.md"
        and row.category == "historical-record"
        and "historical filename" in row.reason
        for row in manifest.preserved
    )


def test_repository_manifest_covers_embedded_legacy_identifiers() -> None:
    manifest = parse_manifest(MANIFEST_PATH.read_text(encoding="utf-8"))
    expected_rows = {
        ("docs/design/harness-v2-design.md", ".agentfield_output.json", "test-fixture"),
        ("docs/design/harness-v2-design.md", ".agentfield_schema.json", "test-fixture"),
        ("specs/sdk-python.md", "AgentFieldHandler", "import-module-path"),
        ("specs/sdk-python.md", "AgentFieldClient", "import-module-path"),
    }

    actual_rows = {(row.path, row.identifier, row.category) for row in manifest.preserved}
    assert expected_rows.issubset(actual_rows)


def test_repository_manifest_logs_broad_legacy_token_verification_scan() -> None:
    verification_section = _extract_section(MANIFEST_PATH.read_text(encoding="utf-8"), "## Verification Commands")
    expected_row = (
        "| `rg -n -e 'AgentField' -e 'agentfield' -e 'AgentPlane' -e 'agentplane' -e 'Agent Plane' -e 'agent plane' "
        "docs specs control-plane/README.md control-plane/scripts/README.md "
        "control-plane/tools/perf/README.md control-plane/migrations/README.md "
        "tests/functional/README.md tests/functional/docker/LOG_DEMO.md --glob '*.md'` | `.` | 0 | "
        "Scoped Markdown review shows only manifested compatibility tokens remain, including "
        "`.agentfield_output.json`, `.agentfield_schema.json`, `AgentFieldHandler`, `AgentFieldClient`, "
        "and the historical `agentplane-ui-api-worklist.md` reference. |"
    )

    assert expected_row in verification_section.splitlines()


def test_repository_manifest_logs_embedded_identifier_verification_scan() -> None:
    verification_section = _extract_section(MANIFEST_PATH.read_text(encoding="utf-8"), "## Verification Commands")
    expected_row = (
        "| `rg -n -e '\\\\.agentfield_output\\\\.json' -e '\\\\.agentfield_schema\\\\.json' "
        "-e 'AgentFieldHandler' -e 'AgentFieldClient' "
        'docs/design/harness-v2-design.md specs/sdk-python.md` | `.` | 0 | Passed: embedded '
        "compatibility identifiers are limited to the manifested harness filenames and public "
        "Python SDK class names. |"
    )

    assert expected_row in verification_section.splitlines()


def test_repository_manifest_logs_yaml_config_verification_scan() -> None:
    verification_section = _extract_section(MANIFEST_PATH.read_text(encoding="utf-8"), "## Verification Commands")

    assert "rg -n -e 'config/agentfield.yaml' -e 'AGENTFIELD_CONFIG_FILE'" in verification_section
    assert "first-class surfaces" in verification_section


def test_log_demo_doc_uses_host_compatibility_data_dir() -> None:
    log_demo_text = (ROOT / "tests" / "functional" / "docker" / "LOG_DEMO.md").read_text(encoding="utf-8")
    manifest = parse_manifest(MANIFEST_PATH.read_text(encoding="utf-8"))

    assert "/tmp/agentfield-log-demo" in log_demo_text
    assert "/tmp/silmari-log-demo" not in log_demo_text
    assert any(
        row.path == "tests/functional/docker/LOG_DEMO.md"
        and row.identifier == "/tmp/agentfield-log-demo"
        and row.category == "test-fixture"
        for row in manifest.preserved
    )


def test_repository_scope_has_no_unmanifested_old_brand_tokens() -> None:
    files = _collect_scoped_markdown_files(ROOT)
    manifest_text = MANIFEST_PATH.read_text(encoding="utf-8")

    errors = validate_scope(ROOT, files, manifest_text)
    assert errors == []
