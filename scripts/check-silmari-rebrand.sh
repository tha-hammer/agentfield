#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if ! command -v python3 >/dev/null 2>&1; then
  echo "ERROR: python3 is required" >&2
  exit 2
fi

python3 - "$REPO_ROOT" <<'PY'
from dataclasses import dataclass
from fnmatch import fnmatch
from pathlib import Path
import re
import subprocess
import sys


OLD_BRAND_RE = re.compile(
    r"agentfield\.ai|Agent-Field|\bAgentField\b|AGENTFIELD|\bagentfield\b|AgentPlane|agentplane|Agent Plane|agent plane"
)

MANIFEST_REL_PATH = "docs/silmari-rebrand-manifest.md"

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

SCAN_ROOTS = (
    "README.md",
    "CODE_OF_CONDUCT.md",
    "SUPPORT.md",
    "SECURITY.md",
    "CLAUDE.md",
    ".github",
    "docs",
    "specs",
    "examples",
    "deployments",
    "skills",
    "control-plane/README.md",
    "control-plane/scripts/README.md",
    "control-plane/tools/perf/README.md",
    "control-plane/migrations/README.md",
    "control-plane/internal/templates",
    "control-plane/internal/skillkit/skill_data/agentfield",
    "control-plane/web/client",
    "sdk",
    "tests",
)

EXCLUDED_DIR_NAMES = {
    ".git",
    ".hypothesis",
    "node_modules",
    "dist",
    "build",
    ".pytest_cache",
    ".mypy_cache",
    ".ruff_cache",
    ".venv",
    "venv",
    "__pycache__",
    "coverage",
    ".next",
}

EXCLUDED_FILE_GLOBS = (
    MANIFEST_REL_PATH,
    "CHANGELOG.md",
    "sdk/python/CHANGELOG.md",
    "*.png",
    "*.jpg",
    "*.jpeg",
    "*.webp",
    "*.gif",
    "*.ico",
    "*.woff",
    "*.woff2",
    "*.ttf",
    "*.eot",
    "*.map",
)

REQUIRED_CODECLEANUP_PASSES = (
    "Brand surface pass",
    "Compatibility pass",
    "Mirror and generated-template pass",
    "Link and asset pass",
    "Formatting and lint pass",
    "Verification pass",
)

ALLOWED_PASS_RESULTS = {"pass", "fail", "not-applicable"}


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
class BrandMatch:
    path: str
    line_number: int
    column_start: int
    column_end: int
    match_text: str
    line_text: str


@dataclass(frozen=True)
class Manifest:
    audited_files: dict[str, AuditedFile]
    preserved_identifiers: list[PreservedIdentifier]


@dataclass(frozen=True)
class VerificationCommandLogView:
    command: str
    working_directory: str
    exit_code: str
    result: str


@dataclass(frozen=True)
class ManifestHealthView:
    missing_headings: tuple[str, ...]
    audited_file_count: int
    preserved_identifier_count: int
    verification_command_count: int
    errors: tuple[str, ...]


@dataclass(frozen=True)
class PreservedIdentifierCoverageView:
    scanned_file_count: int
    total_matches: int
    covered_matches: int
    missing_audited_files: tuple[str, ...]
    unmanifested_matches: tuple[BrandMatch, ...]


def read_manifest_text(manifest_path: Path) -> str:
    try:
        return manifest_path.read_text(encoding="utf-8")
    except FileNotFoundError:
        return ""


def section_lines_by_heading(text: str) -> dict[str, list[str]]:
    lines = text.splitlines()
    heading_indexes: dict[str, int] = {}
    for index, line in enumerate(lines):
        stripped = line.strip()
        if stripped in REQUIRED_HEADINGS:
            heading_indexes[stripped] = index

    sections: dict[str, list[str]] = {}
    for position, heading in enumerate(REQUIRED_HEADINGS):
        if heading not in heading_indexes:
            continue
        start = heading_indexes[heading] + 1
        end = len(lines)
        for next_heading in REQUIRED_HEADINGS[position + 1 :]:
            if next_heading in heading_indexes:
                end = heading_indexes[next_heading]
                break
        sections[heading] = lines[start:end]
    return sections


def split_markdown_row(line: str) -> list[str]:
    stripped = line.strip()
    if not stripped.startswith("|") or not stripped.endswith("|"):
        return []
    return [cell.strip() for cell in stripped[1:-1].split("|")]


def parse_table(section_lines: list[str]) -> tuple[list[str], list[list[str]]] | None:
    table_lines = [line for line in section_lines if line.strip().startswith("|")]
    if len(table_lines) < 2:
        return None
    header = split_markdown_row(table_lines[0])
    rows = [split_markdown_row(line) for line in table_lines[2:]]
    return header, rows


def parse_manifest(manifest_path: Path) -> Manifest:
    text = read_manifest_text(manifest_path)
    sections = section_lines_by_heading(text)

    audited_files: dict[str, AuditedFile] = {}
    audited_table = parse_table(sections.get("## Audited Files", []))
    if audited_table is not None:
        _, rows = audited_table
        for row in rows:
            if len(row) != 3:
                continue
            path, action, verification = row
            if not path:
                continue
            audited_files[path] = AuditedFile(path=path, action=action, verification=verification)

    preserved_identifiers: list[PreservedIdentifier] = []
    preserved_table = parse_table(sections.get("## Preserved Non-Silmari Identifiers", []))
    if preserved_table is not None:
        _, rows = preserved_table
        for row in rows:
            if len(row) != 4:
                continue
            path, identifier, category, reason = row
            if not path:
                continue
            preserved_identifiers.append(
                PreservedIdentifier(
                    path=path,
                    identifier=identifier,
                    category=category,
                    reason=reason,
                )
            )

    return Manifest(
        audited_files=audited_files,
        preserved_identifiers=preserved_identifiers,
    )


def looks_like_separator(cells: list[str]) -> bool:
    if not cells:
        return False
    for cell in cells:
        candidate = cell.replace(":", "").replace("-", "").strip()
        if candidate:
            return False
    return True


def is_concrete_reason(reason: str) -> bool:
    normalized = " ".join(reason.strip().lower().split()).strip(" .")
    if not normalized:
        return False
    if normalized in {
        "compatibility",
        "legacy",
        "preserved",
        "historical",
        "n/a",
        "na",
        "none",
        "tbd",
        "todo",
        "same as above",
    }:
        return False
    tokens = [token for token in re.split(r"[^a-z0-9]+", normalized) if token]
    return len(tokens) >= 3


def reason_covers_all_occurrences(reason: str) -> bool:
    normalized = " ".join(reason.strip().lower().split())
    return bool(
        re.search(
            r"\bcovers all occurrences in (?:this|that) file\b", normalized
        )
    )


def validate_table_shape(
    sections: dict[str, list[str]],
    heading: str,
    expected_header: tuple[str, ...],
) -> tuple[list[str], list[list[str]], list[str]]:
    errors: list[str] = []
    section_lines = sections.get(heading)
    if section_lines is None:
        return [], [], errors

    table = parse_table(section_lines)
    if table is None:
        errors.append(f"Missing required markdown table under {heading}.")
        return [], [], errors

    header, rows = table
    if header != list(expected_header):
        errors.append(
            f"Unexpected columns under {heading}: expected {' | '.join(expected_header)}."
        )

    table_lines = [line for line in section_lines if line.strip().startswith("|")]
    separator = split_markdown_row(table_lines[1])
    if not looks_like_separator(separator):
        errors.append(f"Missing markdown separator row under {heading}.")

    return header, rows, errors


def parse_verification_command_rows(text: str) -> list[VerificationCommandLogView]:
    sections = section_lines_by_heading(text)
    table = parse_table(sections.get("## Verification Commands", []))
    if table is None:
        return []
    _, rows = table
    result: list[VerificationCommandLogView] = []
    for row in rows:
        if len(row) != 4:
            continue
        command, working_directory, exit_code, verification_result = row
        result.append(
            VerificationCommandLogView(
                command=command,
                working_directory=working_directory,
                exit_code=exit_code,
                result=verification_result,
            )
        )
    return result


def build_manifest_health_view(
    text: str, manifest: Manifest, errors: list[str]
) -> ManifestHealthView:
    lines = {line.strip() for line in text.splitlines()}
    missing_headings = tuple(
        heading for heading in REQUIRED_HEADINGS if heading not in lines
    )
    verification_commands = parse_verification_command_rows(text)
    return ManifestHealthView(
        missing_headings=missing_headings,
        audited_file_count=len(manifest.audited_files),
        preserved_identifier_count=len(manifest.preserved_identifiers),
        verification_command_count=len(verification_commands),
        errors=tuple(errors),
    )


def validate_manifest_shape(text: str, manifest: Manifest) -> list[str]:
    errors: list[str] = []
    lines = {line.strip() for line in text.splitlines()}

    for heading in REQUIRED_HEADINGS:
        if heading not in lines:
            errors.append(f"Missing required heading: {heading}")

    sections = section_lines_by_heading(text)

    _, audited_rows, table_errors = validate_table_shape(
        sections,
        "## Audited Files",
        ("Path", "Action", "Verification"),
    )
    errors.extend(table_errors)
    for row in audited_rows:
        if len(row) != 3:
            errors.append("Malformed audited file row: expected 3 columns.")
            continue
        path, action, verification = row
        if not path:
            errors.append("Audited file row is missing Path.")
        if path.startswith("./"):
            errors.append(f"Audited file path must be repo-relative without ./ prefix: {path}")
        if action not in ALLOWED_ACTIONS:
            errors.append(f"Unknown audited file action for {path}: {action}")
        if not verification:
            errors.append(f"Audited file row is missing Verification for {path}.")

    _, preserved_rows, table_errors = validate_table_shape(
        sections,
        "## Preserved Non-Silmari Identifiers",
        ("Path", "Identifier", "Category", "Reason"),
    )
    errors.extend(table_errors)
    for row in preserved_rows:
        if len(row) != 4:
            errors.append("Malformed preserved identifier row: expected 4 columns.")
            continue
        path, identifier, category, reason = row
        if not path:
            errors.append("Preserved identifier row is missing Path.")
        if path.startswith("./"):
            errors.append(
                f"Preserved identifier path must be repo-relative without ./ prefix: {path}"
            )
        if not identifier:
            errors.append(f"Preserved identifier row is missing Identifier for {path}.")
        if category not in ALLOWED_CATEGORIES:
            errors.append(f"Unknown preserved identifier category for {path}: {category}")
        if not is_concrete_reason(reason):
            errors.append(f"Reason must be concrete for preserved identifier {path}: {identifier}")

    _, codecleanup_rows, table_errors = validate_table_shape(
        sections,
        "## CodeCleanup Passes",
        ("Pass", "Result", "Notes"),
    )
    errors.extend(table_errors)
    seen_passes: list[str] = []
    for row in codecleanup_rows:
        if len(row) != 3:
            errors.append("Malformed CodeCleanup row: expected 3 columns.")
            continue
        pass_name, result, notes = row
        seen_passes.append(pass_name)
        if result not in ALLOWED_PASS_RESULTS:
            errors.append(f"Unknown CodeCleanup result for {pass_name}: {result}")
        if not notes:
            errors.append(f"CodeCleanup row is missing Notes for {pass_name}.")
    if codecleanup_rows and tuple(seen_passes) != REQUIRED_CODECLEANUP_PASSES:
        errors.append("CodeCleanup Passes must list the six required pass names in order.")

    _, verification_rows, table_errors = validate_table_shape(
        sections,
        "## Verification Commands",
        ("Command", "Working Directory", "Exit Code", "Result"),
    )
    errors.extend(table_errors)
    if not verification_rows:
        errors.append("Verification Commands table must contain at least one row.")
    for row in verification_rows:
        if len(row) != 4:
            errors.append("Malformed verification command row: expected 4 columns.")
            continue
        command, working_directory, exit_code, result = row
        if not command:
            errors.append("Verification command row is missing Command.")
        if not working_directory:
            errors.append(f"Verification command row is missing Working Directory for {command}.")
        if not exit_code:
            errors.append(f"Verification command row is missing Exit Code for {command}.")
        if not result:
            errors.append(f"Verification command row is missing Result for {command}.")

    return errors


def is_excluded_file(relative_path: str) -> bool:
    return any(fnmatch(relative_path, pattern) for pattern in EXCLUDED_FILE_GLOBS)


def is_binary_file(path: Path) -> bool:
    try:
        with path.open("rb") as handle:
            return b"\x00" in handle.read(4096)
    except OSError:
        return False


def has_excluded_relative_directory(relative_path: Path) -> bool:
    return any(part in EXCLUDED_DIR_NAMES for part in relative_path.parent.parts)


def collect_scan_files(repo_root: Path) -> list[Path]:
    collected: list[Path] = []
    seen: set[str] = set()

    for scan_root in SCAN_ROOTS:
        candidate = repo_root / scan_root
        if not candidate.exists():
            continue

        if candidate.is_file():
            relative = candidate.relative_to(repo_root).as_posix()
            if relative == MANIFEST_REL_PATH or is_excluded_file(relative):
                continue
            if is_binary_file(candidate):
                continue
            if relative not in seen:
                seen.add(relative)
                collected.append(candidate)
            continue

        for path in candidate.rglob("*"):
            if not path.is_file():
                continue
            relative_path = path.relative_to(repo_root)
            if has_excluded_relative_directory(relative_path):
                continue
            relative = relative_path.as_posix()
            if relative == MANIFEST_REL_PATH or is_excluded_file(relative):
                continue
            if is_binary_file(path):
                continue
            if relative not in seen:
                seen.add(relative)
                collected.append(path)

    return sorted(collected, key=lambda path: path.relative_to(repo_root).as_posix())


def find_brand_matches(repo_root: Path, files: list[Path]) -> list[BrandMatch]:
    matches: list[BrandMatch] = []
    for path in files:
        relative = path.relative_to(repo_root).as_posix()
        try:
            text = path.read_text(encoding="utf-8", errors="replace")
        except OSError:
            continue
        for line_number, line_text in enumerate(text.splitlines(), start=1):
            for match in OLD_BRAND_RE.finditer(line_text):
                matches.append(
                    BrandMatch(
                        path=relative,
                        line_number=line_number,
                        column_start=match.start(),
                        column_end=match.end(),
                        match_text=match.group(0),
                        line_text=line_text,
                    )
                )
    return matches


def identifier_ranges(line_text: str, identifier: str) -> list[tuple[int, int]]:
    ranges: list[tuple[int, int]] = []
    if not identifier:
        return ranges
    start = 0
    while True:
        index = line_text.find(identifier, start)
        if index == -1:
            return ranges
        ranges.append((index, index + len(identifier)))
        start = index + 1


def is_exact_token_coverage(match: BrandMatch, row: PreservedIdentifier) -> bool:
    return match.match_text == row.identifier


def is_enclosing_identifier_coverage(
    match: BrandMatch, row: PreservedIdentifier
) -> bool:
    if match.match_text not in row.identifier:
        return False
    for start, end in identifier_ranges(match.line_text, row.identifier):
        if start <= match.column_start and match.column_end <= end:
            return True
    return False


def row_covers_match(match: BrandMatch, row: PreservedIdentifier) -> bool:
    return is_exact_token_coverage(match, row) or is_enclosing_identifier_coverage(
        match, row
    )


def valid_preserved_rows_by_path(
    manifest: Manifest,
) -> dict[str, list[PreservedIdentifier]]:
    rows_by_path: dict[str, list[PreservedIdentifier]] = {}
    for row in manifest.preserved_identifiers:
        if row.category not in ALLOWED_CATEGORIES:
            continue
        if not is_concrete_reason(row.reason):
            continue
        rows_by_path.setdefault(row.path, []).append(row)
    return rows_by_path


def build_exact_token_row_counts(
    rows_by_path: dict[str, list[PreservedIdentifier]],
) -> tuple[dict[tuple[str, str], int], set[tuple[str, str]]]:
    counts: dict[tuple[str, str], int] = {}
    all_occurrence_keys: set[tuple[str, str]] = set()
    for path, rows in rows_by_path.items():
        for row in rows:
            key = (path, row.identifier)
            if reason_covers_all_occurrences(row.reason):
                all_occurrence_keys.add(key)
                continue
            counts[key] = counts.get(key, 0) + 1
    return counts, all_occurrence_keys


def is_manifested(
    match: BrandMatch,
    manifest: Manifest,
    rows_by_path: dict[str, list[PreservedIdentifier]],
    available_exact_rows: dict[tuple[str, str], int],
    all_occurrence_exact_rows: set[tuple[str, str]],
) -> bool:
    if match.path not in manifest.audited_files:
        return False

    for row in rows_by_path.get(match.path, []):
        if row.identifier == match.match_text:
            continue
        if row_covers_match(match, row):
            return True

    exact_key = (match.path, match.match_text)
    if exact_key in all_occurrence_exact_rows:
        return True
    if available_exact_rows.get(exact_key, 0) > 0:
        available_exact_rows[exact_key] -= 1
        return True

    return False


def build_preserved_identifier_coverage_view(
    matches: list[BrandMatch],
    manifest: Manifest,
    scanned_file_count: int,
) -> PreservedIdentifierCoverageView:
    missing_audited_files: list[str] = []
    missing_seen: set[str] = set()
    unmanifested_matches: list[BrandMatch] = []
    covered_matches = 0
    rows_by_path = valid_preserved_rows_by_path(manifest)
    available_exact_rows, all_occurrence_exact_rows = build_exact_token_row_counts(
        rows_by_path
    )

    for match in matches:
        if match.path not in manifest.audited_files:
            if match.path not in missing_seen:
                missing_seen.add(match.path)
                missing_audited_files.append(match.path)
            continue
        if is_manifested(
            match,
            manifest,
            rows_by_path,
            available_exact_rows,
            all_occurrence_exact_rows,
        ):
            covered_matches += 1
            continue
        unmanifested_matches.append(match)

    return PreservedIdentifierCoverageView(
        scanned_file_count=scanned_file_count,
        total_matches=len(matches),
        covered_matches=covered_matches,
        missing_audited_files=tuple(sorted(missing_audited_files)),
        unmanifested_matches=tuple(unmanifested_matches),
    )


def validate_matches(matches: list[BrandMatch], manifest: Manifest) -> list[str]:
    errors: list[str] = []
    coverage_view = build_preserved_identifier_coverage_view(matches, manifest, 0)
    for path in coverage_view.missing_audited_files:
        errors.append(f"Missing audited file row for {path}")
    for match in coverage_view.unmanifested_matches:
        errors.append(
            f"Unmanifested identifier {match.path}:{match.line_number}: {match.match_text}"
        )
    return errors


def verify_skill_mirror(repo_root: Path) -> tuple[bool, str]:
    try:
        proc = subprocess.run(
            ["./scripts/sync-embedded-skills.sh", "--check"],
            cwd=repo_root,
            capture_output=True,
            text=True,
            check=False,
        )
    except OSError as exc:
        return False, f"unable to run ./scripts/sync-embedded-skills.sh --check: {exc}"

    message = proc.stdout.strip() or proc.stderr.strip()
    if proc.returncode == 0:
        return True, message or "All embedded skills are in sync."
    return False, message or "skill mirror drift detected"


def format_match(match: BrandMatch) -> str:
    excerpt = match.line_text.strip() or match.line_text
    return f"{match.path}:{match.line_number}: {match.match_text} :: {excerpt}"


def print_failure(
    coverage_view: PreservedIdentifierCoverageView,
    manifest_errors: list[str],
    manifest: Manifest,
    mirror_status: str,
) -> None:
    print("Silmari rebrand check failed.")
    print()
    print("Unmanifested identifiers:")
    if coverage_view.unmanifested_matches:
        for match in coverage_view.unmanifested_matches:
            print(f"  {format_match(match)}")
    else:
        print("  (none)")

    print()
    print("Files with identifiers missing from Audited Files:")
    if coverage_view.missing_audited_files:
        for path in coverage_view.missing_audited_files:
            print(f"  {path}")
    else:
        print("  (none)")

    print()
    print("Manifest errors:")
    if manifest_errors:
        for error in manifest_errors:
            print(f"  {error}")
    else:
        print("  (none)")

    print()
    print(
        f"Scanned {coverage_view.scanned_file_count} files; "
        f"validated {len(manifest.preserved_identifiers)} preserved identifier rows; "
        f"skill mirror status: {mirror_status}"
    )


def main(repo_root: Path) -> int:
    manifest_path = repo_root / MANIFEST_REL_PATH
    manifest_text = read_manifest_text(manifest_path)
    manifest = parse_manifest(manifest_path)
    manifest_errors = validate_manifest_shape(manifest_text, manifest)

    scan_files = collect_scan_files(repo_root)
    matches = find_brand_matches(repo_root, scan_files)
    coverage_view = build_preserved_identifier_coverage_view(
        matches, manifest, len(scan_files)
    )

    mirror_ok, mirror_message = verify_skill_mirror(repo_root)
    mirror_status = mirror_message or "skill mirror check unavailable"
    if not mirror_ok:
        manifest_errors.append(f"Skill mirror drift detected: {mirror_status}")

    manifest_view = build_manifest_health_view(manifest_text, manifest, manifest_errors)
    if manifest_view.errors or coverage_view.unmanifested_matches or coverage_view.missing_audited_files:
        print_failure(coverage_view, list(manifest_view.errors), manifest, mirror_status)
        return 1

    print("Silmari rebrand check passed.")
    print(
        f"Scanned {coverage_view.scanned_file_count} files; "
        f"validated {manifest_view.preserved_identifier_count} preserved identifier rows; "
        f"skill mirror is in sync."
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main(Path(sys.argv[1])))
PY
