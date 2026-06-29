#!/usr/bin/env python3
"""Collect the baseline Silmari rebrand inventory for downstream issues."""

from __future__ import annotations

import argparse
import json
import subprocess
import sys
from collections import Counter
from pathlib import Path
from typing import Iterable

EXPECTED_BRANCH = "integration/silmari-rebrand-agentfield"
BRAND_PATTERN = (
    r"\bAgentField\b|\bagentfield\b|AGENTFIELD|Agent-Field|agentfield\.ai|"
    r"AgentPlane|agentplane|Agent Plane|agent plane"
)
SCAN_ROOTS = [
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
]
SCAN_GLOBS = [
    "!**/node_modules/**",
    "!**/dist/**",
    "!**/.pytest_cache/**",
]
PROTECTED_RUNTIME_AREAS = [
    "control-plane/internal/server/**",
    "control-plane/internal/handlers/**",
    "control-plane/internal/storage/**",
    "control-plane/internal/services/**",
    "control-plane/migrations/*.sql",
    "control-plane/pkg/types/**",
    "sdk/python/agentfield/**",
    "sdk/go/**",
    "sdk/typescript/src/**",
    "runtime route names",
    "storage schema and JSON field names",
    "package names and module paths",
    "Docker image names and Kubernetes resource names",
    "Helm chart names and helper names",
    "env vars and YAML config paths",
    "skill slugs and embedded skill paths",
    "web UI storage keys and API path identifiers",
]
EDITABLE_BRAND_SURFACES = [
    "root/community docs and templates",
    "docs and specs",
    "example READMEs, UI copy, and visible benchmark labels",
    "deployment READMEs, chart descriptions, and NOTES output",
    "generated template README/docstring surfaces",
    "skills/agentfield source prose and embedded mirror prose",
    "admin UI labels, navigation, empty states, and About copy",
    "SDK README and package metadata descriptions",
]
PLANNED_RENAMES = {
    "runtime_routes": False,
    "storage_schema": False,
    "sdk_api": False,
    "package_name": False,
    "module_path": False,
    "image_name": False,
    "helm_chart_name": False,
    "env_var": False,
    "yaml_config_path": False,
    "skill_slug": False,
}


def repo_root() -> Path:
    return Path(
        subprocess.check_output(
            ["git", "rev-parse", "--show-toplevel"], text=True
        ).strip()
    )


def current_branch(cwd: Path) -> str:
    return (
        subprocess.check_output(
            ["git", "rev-parse", "--abbrev-ref", "HEAD"], cwd=cwd, text=True
        )
        .strip()
    )


def assert_branch(expected: str, cwd: Path) -> int:
    branch = current_branch(cwd)
    payload = {
        "expected_branch": expected,
        "current_branch": branch,
        "matches": branch == expected,
    }
    print(json.dumps(payload, indent=2, sort_keys=True))
    return 0 if branch == expected else 1


def rg_command(paths: Iterable[str]) -> list[str]:
    cmd = ["rg", "--json", "--ignore-case", "-o", BRAND_PATTERN, *paths]
    for glob in SCAN_GLOBS:
        cmd.extend(["--glob", glob])
    return cmd


def collect_matches(paths: Iterable[str], cwd: Path) -> dict[str, int]:
    proc = subprocess.run(
        rg_command(paths),
        cwd=cwd,
        check=False,
        capture_output=True,
        text=True,
    )
    if proc.returncode not in (0, 1):
        raise RuntimeError(proc.stderr.strip() or proc.stdout.strip())

    matches: dict[str, int] = {}
    for raw_line in proc.stdout.splitlines():
        event = json.loads(raw_line)
        if event.get("type") != "match":
            continue
        data = event["data"]
        path = data["path"]["text"]
        count = len(data["submatches"])
        matches[path] = matches.get(path, 0) + count
    return dict(sorted(matches.items()))


def classify_path(path: str) -> tuple[str, str]:
    name = Path(path).name

    if path == "sdk/python/CHANGELOG.md":
        return (
            "historical",
            "The changelog is preserved as a historical record rather than "
            "rewritten during the rebrand.",
        )
    if path.startswith("examples/benchmarks/100k-scale/results/"):
        return (
            "historical",
            "Recorded benchmark outputs are preserved as historical evidence "
            "for later manifest rows.",
        )

    if path == "control-plane/internal/templates/templates.go":
        return (
            "runtime-compatible",
            "Template embedding code is runtime implementation and keeps "
            "stable image and skill identifiers.",
        )
    if path == "control-plane/web/client/ui_embed.go":
        return (
            "runtime-compatible",
            "The embedded UI registration code is a runtime surface and "
            "should not be renamed for branding.",
        )
    if path == "control-plane/web/client/src/contexts/ModeContext.tsx":
        return (
            "runtime-compatible",
            "This file carries the stable web UI localStorage key "
            "`agentfield-app-mode` and stays compatibility-safe.",
        )
    if path.startswith("sdk/python/agentfield/"):
        return (
            "runtime-compatible",
            "Python SDK runtime and API source keeps stable package names, "
            "method names, and compatibility identifiers.",
        )
    if path.startswith("sdk/typescript/src/"):
        return (
            "runtime-compatible",
            "TypeScript SDK runtime and API source keeps stable module paths, "
            "client field names, and request identifiers.",
        )
    if (
        path.startswith("sdk/go/")
        and name != "go.mod"
        and not path.endswith(".md")
        and not path.endswith("_test.go")
    ):
        return (
            "runtime-compatible",
            "Go SDK runtime and API source keeps stable module paths and "
            "public compatibility identifiers.",
        )
    if path.startswith("deployments/kubernetes/"):
        if path.endswith("README.md"):
            return (
                "edit",
                "Deployment README prose is user-facing and should use "
                "Silmari while preserving config identifiers.",
            )
        return (
            "runtime-compatible",
            "Kubernetes manifests carry stable resource names, image names, "
            "env vars, and config paths.",
        )
    if path.startswith("deployments/docker/"):
        if path.endswith("README.md"):
            return (
                "edit",
                "Docker deployment documentation is user-facing and should "
                "use Silmari wording.",
            )
        return (
            "runtime-compatible",
            "Dockerfiles and compose manifests carry stable image names, "
            "binary names, volumes, and env vars.",
        )
    if path.startswith("deployments/helm/agentfield/"):
        if path.endswith("README.md"):
            return (
                "edit",
                "Helm documentation is user-facing and should use Silmari "
                "wording.",
            )
        if path.endswith("Chart.yaml"):
            return (
                "edit",
                "The chart description needs Silmari copy even though the "
                "chart name `agentfield` stays stable.",
            )
        if path.endswith("templates/NOTES.txt"):
            return (
                "edit",
                "Post-install NOTES output is user-visible and should use "
                "Silmari copy while helper names stay unchanged.",
            )
        return (
            "runtime-compatible",
            "Helm templates and values keep stable chart, helper, resource, "
            "and env-var identifiers.",
        )

    if path in {
        "README.md",
        "CODE_OF_CONDUCT.md",
        "SUPPORT.md",
        "SECURITY.md",
        "CLAUDE.md",
    }:
        return (
            "edit",
            "Root-facing product prose is user-visible and should use "
            "Silmari branding.",
        )
    if path.startswith(".github/"):
        return (
            "edit",
            "Issue templates, PR templates, and workflow-facing prose need "
            "Silmari branding review.",
        )
    if path.startswith("docs/") or path.startswith("specs/"):
        return (
            "edit",
            "Documentation and spec prose are in-scope brand surfaces for the "
            "Silmari copy pass.",
        )
    if path in {
        "control-plane/README.md",
        "control-plane/scripts/README.md",
        "control-plane/tools/perf/README.md",
        "control-plane/migrations/README.md",
    }:
        return (
            "edit",
            "Control-plane documentation is user-facing and should use "
            "Silmari wording.",
        )
    if path.startswith("skills/agentfield/") or path.startswith(
        "control-plane/internal/skillkit/skill_data/agentfield/"
    ):
        return (
            "edit",
            "Skill prose and reference content need Silmari wording while "
            "skill slugs stay stable.",
        )
    if path.startswith("control-plane/internal/templates/"):
        if name in {"requirements.txt.tmpl", "package.json.tmpl", "go.mod.tmpl"}:
            return (
                "preserve",
                "Generated dependency and module template identifiers stay "
                "stable for compatibility.",
            )
        if name == "docker-compose.yml.tmpl":
            return (
                "runtime-compatible",
                "Generated deployment config keeps stable image names, env "
                "vars, and config paths.",
            )
        if path.endswith("_test.go"):
            return (
                "preserve",
                "Template tests intentionally assert stable package and import "
                "identifiers.",
            )
        return (
            "edit",
            "Generated README or source-template comments shown to users need "
            "Silmari copy.",
        )
    if path.startswith("control-plane/web/client/src/test/"):
        if "/pages/" in path or "/components/" in path:
            return (
                "edit",
                "UI tests that assert visible copy should move with the "
                "Silmari UI wording pass.",
            )
        return (
            "preserve",
            "Service and utility tests intentionally keep stable JSON fields, "
            "workspace identifiers, or API fixtures.",
        )
    if path.startswith("control-plane/web/client/src/"):
        runtime_client_paths = (
            "/services/",
            "/types/",
            "/hooks/queries/",
            "/utils/",
        )
        if any(segment in path for segment in runtime_client_paths):
            return (
                "preserve",
                "This admin UI source carries compatibility-sensitive service "
                "or schema identifiers rather than visible brand copy.",
            )
        return (
            "edit",
            "Admin UI labels, navigation, or visible component copy are "
            "rebrand surfaces.",
        )

    if path in {
        "sdk/python/README.md",
        "sdk/typescript/README.md",
        "sdk/go/README.md",
        "sdk/go/ai/README.md",
    }:
        return (
            "edit",
            "SDK README prose is user-facing and should use Silmari wording.",
        )
    if path in {"sdk/python/pyproject.toml", "sdk/typescript/package.json"}:
        return (
            "edit",
            "Published SDK metadata descriptions and keywords need Silmari "
            "branding while package names stay stable.",
        )
    if path in {
        "sdk/python/uv.lock",
        "sdk/python/requirements-dev.txt",
        "sdk/python/MANIFEST.in",
        "sdk/typescript/package-lock.json",
        "sdk/go/go.mod",
    }:
        return (
            "preserve",
            "Lockfile, manifest, or module identifiers stay stable for "
            "package compatibility.",
        )
    if path.startswith("sdk/python/tests/") or path.startswith(
        "sdk/typescript/tests/"
    ):
        return (
            "preserve",
            "SDK tests intentionally keep stable AgentField fixtures, import "
            "paths, and compatibility assertions.",
        )
    if path.startswith("sdk/go/") and path.endswith("_test.go"):
        return (
            "preserve",
            "Go SDK tests intentionally keep stable module paths and "
            "compatibility assertions.",
        )
    if path.startswith("sdk/python/scripts/") or path.startswith(
        "sdk/typescript/scripts/"
    ):
        return (
            "preserve",
            "SDK helper scripts keep stable package and runtime identifiers "
            "for compatibility.",
        )

    if path.startswith("tests/functional/"):
        if path.endswith("README.md") or path.endswith("LOG_DEMO.md"):
            return (
                "edit",
                "Functional test documentation is user-facing and should use "
                "Silmari wording.",
            )
        if path.endswith(".yaml") or path.endswith(".yml"):
            return (
                "runtime-compatible",
                "Functional test configs keep stable config keys, env vars, "
                "and file paths.",
            )
        return (
            "preserve",
            "Functional test fixtures keep stable package, route, and "
            "compatibility identifiers.",
        )

    if path.startswith("examples/"):
        if name in {
            "package.json",
            "package-lock.json",
            "go.mod",
            "requirements.txt",
            "Dockerfile",
        }:
            return (
                "preserve",
                "Example dependency and module identifiers stay stable so the "
                "sample remains runnable.",
            )
        if name.startswith("test_") or name.endswith("_test.py"):
            return (
                "preserve",
                "Example-local tests keep stable compatibility fixtures and "
                "package imports.",
            )
        if "/ui/" in path or path.endswith(".svg"):
            return (
                "edit",
                "Example UI copy and visible asset labels are user-facing "
                "rebrand surfaces.",
            )
        if "benchmarks/100k-scale/" in path:
            return (
                "edit",
                "Benchmark labels, chart text, and comparison copy are visible "
                "brand surfaces for the rebrand pass.",
            )
        return (
            "edit",
            "Example README, product context, script output, or visible sample "
            "copy should use Silmari wording.",
        )

    if path.startswith("deployments/"):
        return (
            "edit",
            "Deployment-facing documentation is user-facing and should use "
            "Silmari wording.",
        )

    return (
        "preserve",
        "This match is kept in the baseline inventory for compatibility review "
        "instead of a broad rename.",
    )


def build_inventory(paths: Iterable[str], cwd: Path) -> dict[str, object]:
    matches = collect_matches(paths, cwd)
    files = []
    counts = Counter()
    for path, match_count in matches.items():
        classification, reason = classify_path(path)
        counts[classification] += 1
        files.append(
            {
                "path": path,
                "matches": match_count,
                "classification": classification,
                "reason": reason,
            }
        )

    return {
        "commandRef": "Brand Surface Inventory",
        "command": " ".join(rg_command(paths)),
        "scanRoots": list(paths),
        "classifiedFiles": files,
        "classificationCounts": dict(sorted(counts.items())),
        "totalMatchedFiles": len(files),
        "totalMatches": sum(item["matches"] for item in files),
    }


def cmd_collect(args: argparse.Namespace) -> int:
    cwd = repo_root()
    paths = args.paths or SCAN_ROOTS
    payload = {
        "branch": current_branch(cwd),
        "brandInventoryCollected": build_inventory(paths, cwd),
        "changeSetBoundaryVerified": {
            "protectedRuntimeAreas": PROTECTED_RUNTIME_AREAS,
            "editableBrandSurfaces": EDITABLE_BRAND_SURFACES,
            "plannedRenames": PLANNED_RENAMES,
            "notes": (
                "No runtime route, storage schema, SDK API, package name, "
                "module path, image name, Helm chart name, env var, YAML "
                "config path, or skill slug rename is planned."
            ),
        },
    }
    rendered = json.dumps(payload, indent=2, sort_keys=True)
    if args.output:
        Path(args.output).write_text(rendered + "\n", encoding="utf-8")
    else:
        print(rendered)
    return 0


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description="Collect the baseline Silmari rebrand inventory."
    )
    subparsers = parser.add_subparsers(dest="command", required=True)

    assert_parser = subparsers.add_parser(
        "assert-branch", help="Assert the exact rebrand integration branch."
    )
    assert_parser.add_argument(
        "--expected",
        default=EXPECTED_BRANCH,
        help="Expected branch name.",
    )

    collect_parser = subparsers.add_parser(
        "collect", help="Collect and classify brand matches."
    )
    collect_parser.add_argument(
        "--path",
        dest="paths",
        action="append",
        help="Optional scan path override. Repeat to scan multiple paths.",
    )
    collect_parser.add_argument(
        "--output",
        help="Write the JSON report to this file instead of stdout.",
    )
    return parser


def main(argv: list[str] | None = None) -> int:
    parser = build_parser()
    args = parser.parse_args(argv)
    cwd = repo_root()

    if args.command == "assert-branch":
        return assert_branch(args.expected, cwd)
    if args.command == "collect":
        return cmd_collect(args)

    parser.error(f"unknown command: {args.command}")
    return 2


if __name__ == "__main__":
    sys.exit(main())
