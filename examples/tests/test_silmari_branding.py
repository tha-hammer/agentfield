from __future__ import annotations

import asyncio
import importlib.util
import json
import re
import textwrap
import unittest
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable


REPO_ROOT = Path(__file__).resolve().parents[2]
MANIFEST_PATH = REPO_ROOT / "docs" / "silmari-rebrand-manifest.md"
RAG_UI_PUBLIC = (
    REPO_ROOT / "examples" / "python_agent_nodes" / "rag_evaluation" / "ui" / "public"
)
OLD_BRAND_RE = re.compile(
    r"agentfield\.ai|Agent-Field|\bAgentField\b|AGENTFIELD|\bagentfield\b|"
    r"AgentPlane|agentplane|Agent Plane|agent plane"
)
SVG_SRC_RE = re.compile(r'src="/([^"]+\.svg)"')
CURATED_EXAMPLE_SURFACE_PATTERNS = (
    "assets/utm-links.csv",
    "examples/**/README.md",
    "examples/**/main.py",
    "examples/**/main.ts",
    "examples/**/main.go",
    "examples/**/agent.py",
    "examples/**/benchmark.py",
    "examples/**/benchmark.ts",
    "examples/**/analyze.py",
    "examples/**/run_benchmarks.sh",
    "examples/**/page.tsx",
    "examples/**/PoweredBy.tsx",
    "examples/**/*.svg",
)
ALLOWED_CATEGORIES = {
    "published-link-target",
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
    "runtime-compatibility",
    "historical-record",
    "test-fixture",
}


def read_text(rel_path: str) -> str:
    return (REPO_ROOT / rel_path).read_text(encoding="utf-8")


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
class BrandOccurrence:
    source: str
    identifier: str
    start: int
    end: int


def _parse_markdown_table(text: str, heading: str) -> list[list[str]]:
    marker = f"{heading}\n"
    if marker not in text:
        raise ValueError(f"Missing heading: {heading}")

    table_started = False
    rows: list[list[str]] = []
    for line in text.splitlines()[text.splitlines().index(heading) + 1 :]:
        if not line.startswith("|"):
            if table_started:
                break
            continue
        table_started = True
        cells = [cell.strip() for cell in line.strip().strip("|").split("|")]
        if set("".join(cells)) <= {"-"}:
            continue
        rows.append(cells)

    if len(rows) < 2:
        raise ValueError(f"Missing rows for {heading}")

    return rows[1:]


def parse_manifest() -> tuple[dict[str, AuditedFile], list[PreservedIdentifier]]:
    text = MANIFEST_PATH.read_text(encoding="utf-8")

    audited_rows = _parse_markdown_table(text, "## Audited Files")
    audited = {
        row[0]: AuditedFile(path=row[0], action=row[1], verification=row[2])
        for row in audited_rows
    }

    preserved_rows = _parse_markdown_table(text, "## Preserved Non-Silmari Identifiers")
    preserved = [
        PreservedIdentifier(path=row[0], identifier=row[1], category=row[2], reason=row[3])
        for row in preserved_rows
    ]

    return audited, preserved


def load_rag_eval_client_module():
    module_path = (
        REPO_ROOT
        / "examples"
        / "python_agent_nodes"
        / "rag_evaluation"
        / "rag_eval_client.py"
    )
    spec = importlib.util.spec_from_file_location("rag_eval_client", module_path)
    if spec is None or spec.loader is None:
        raise AssertionError("Unable to load rag_eval_client module")

    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def identifier_ranges(line_text: str | None, identifier: str | None) -> list[tuple[int, int]]:
    if not line_text or not identifier:
        return []

    ranges: list[tuple[int, int]] = []
    start = 0
    while True:
        index = line_text.find(identifier, start)
        if index == -1:
            return ranges
        ranges.append((index, index + len(identifier)))
        start = index + 1


def iter_old_brand_matches(text: str | None) -> Iterable[re.Match[str]]:
    if not text:
        return ()
    return OLD_BRAND_RE.finditer(text)


def iter_old_brand_identifiers(text: str | None) -> Iterable[str]:
    if not text:
        return ()
    return (match.group(0) for match in iter_old_brand_matches(text))


def iter_old_brand_occurrences(rel_path: str, text: str) -> Iterable[BrandOccurrence]:
    for match in iter_old_brand_matches(rel_path):
        start, end = match.span()
        yield BrandOccurrence("path", match.group(0), start, end)

    for match in iter_old_brand_matches(text):
        start, end = match.span()
        yield BrandOccurrence("text", match.group(0), start, end)


def extract_svg_sources(text: str | None) -> list[str]:
    if not text:
        return []
    return SVG_SRC_RE.findall(text)


def iter_curated_example_surface_paths() -> Iterable[Path]:
    seen: set[Path] = set()
    for pattern in CURATED_EXAMPLE_SURFACE_PATTERNS:
        for path in REPO_ROOT.glob(pattern):
            if not path.is_file() or path in seen:
                continue
            seen.add(path)
            yield path


def manifested_identifier_covers_occurrence(
    manifest_identifier: str,
    occurrence: BrandOccurrence,
    *,
    rel_path: str,
    text: str,
) -> bool:
    if manifest_identifier == occurrence.identifier:
        return True

    if occurrence.identifier not in manifest_identifier:
        return False

    haystack = rel_path if occurrence.source == "path" else text
    for start, end in identifier_ranges(haystack, manifest_identifier):
        if start <= occurrence.start and occurrence.end <= end:
            return True

    return False


class SilmariBrandingTests(unittest.TestCase):
    def test_visible_copy_uses_silmari(self) -> None:
        checks = {
            "assets/utm-links.csv": ("Built with Silmari", "Built with AF"),
            "examples/README.md": ("# Silmari Examples", "AgentField Examples"),
            "examples/triggers-demo/README.md": ("Run Silmari", "Run AgentField"),
            "examples/triggers-demo/agent.py": (
                "Silmari triggers demo",
                "AgentField triggers demo",
            ),
            "examples/triggers-demo/scripts/fire-events.sh": (
                "local Silmari control plane",
                "local AgentField control plane",
            ),
            "examples/python_agent_nodes/hello_world/main.py": (
                "Welcome to Silmari.",
                "Welcome to Agentfield.",
            ),
            "examples/python_agent_nodes/docker_hello_world/main.py": (
                "full Silmari execution path",
                "full AgentField execution path",
            ),
            "examples/python_agent_nodes/hello_world_rag/README.md": (
                "Silmari CLI",
                "AgentField CLI",
            ),
            "examples/python_agent_nodes/hello_world_rag/main.py": (
                "Using Silmari memory vectors",
                "Using AgentField memory vectors",
            ),
            "examples/python_agent_nodes/image_generation_hello_world/README.md": (
                "Silmari Control Plane",
                "AgentField Control Plane",
            ),
            "examples/python_agent_nodes/agentic_rag/main.py": (
                "Silmari's unified vector memory",
                "AgentField's unified vector memory",
            ),
            "examples/python_agent_nodes/deep_research/README.md": (
                '"research_question": "What is Silmari?"',
                "What is AgentField.ai?",
            ),
            "examples/python_agent_nodes/deep_research/main.py": (
                "Elegant and simple Silmari primitives",
                "Elegant and simple AgentField primitives",
            ),
            "examples/python_agent_nodes/documentation_chatbot/README.md": (
                "stores vectors in Silmari's global memory scope",
                "stores vectors in AgentField's global memory scope",
            ),
            "examples/python_agent_nodes/documentation_chatbot/product_context.py": (
                "Silmari is a Kubernetes-style control plane",
                "AgentField is a Kubernetes-style control plane",
            ),
            "examples/python_agent_nodes/documentation_chatbot/routers/qa.py": (
                "To get started with Silmari:",
                "To get started with AgentField:",
            ),
            "examples/python_agent_nodes/rag_evaluation/rag_eval_client.py": (
                "RAG Evaluation Silmari node",
                "RAG Evaluation AgentField node",
            ),
            "examples/python_agent_nodes/serverless_hello/main.py": (
                'name: str = "Silmari"',
                'name: str = "AgentField"',
            ),
            "examples/ts-node-examples/serverless-hello/main.ts": (
                "Silmari'}!",
                "AgentField'}!",
            ),
            "examples/ts-node-examples/verifiable-credentials/README.md": (
                "in Silmari to create cryptographically verifiable audit trails",
                "in AgentField to create cryptographically verifiable audit trails",
            ),
            "examples/ts-node-examples/verifiable-credentials/main.ts": (
                "Verifiable Credentials (VCs) in Silmari",
                "Verifiable Credentials (VCs) in AgentField",
            ),
            "examples/benchmarks/100k-scale/README.md": (
                "# Silmari Scale Benchmark",
                "# AgentField Scale Benchmark",
            ),
            "examples/benchmarks/100k-scale/go-bench/main.go": (
                "Silmari Go SDK Benchmark",
                "AgentField Go SDK Benchmark",
            ),
            "examples/benchmarks/100k-scale/python-bench/benchmark.py": (
                "Silmari Python SDK Benchmark",
                "AgentField Python SDK Benchmark",
            ),
            "examples/benchmarks/100k-scale/ts-bench/benchmark.ts": (
                "Silmari TypeScript SDK Benchmark",
                "AgentField TypeScript SDK Benchmark",
            ),
            "examples/benchmarks/100k-scale/langchain-bench/benchmark.py": (
                "equivalent operations to Silmari",
                "equivalent operations to AgentField",
            ),
            "examples/benchmarks/100k-scale/crewai-bench/benchmark.py": (
                "equivalent operations to Silmari",
                "equivalent operations to AgentField",
            ),
            "examples/benchmarks/100k-scale/run_benchmarks.sh": (
                "Silmari Scale Benchmark Suite",
                "AgentField Scale Benchmark Suite",
            ),
        }

        for rel_path, (required, forbidden) in checks.items():
            with self.subTest(path=rel_path):
                text = read_text(rel_path)
                self.assertIn(required, text)
                self.assertNotIn(forbidden, text)

    def test_helper_edge_cases_cover_empty_and_boundary_inputs(self) -> None:
        self.assertEqual([], list(iter_old_brand_matches(None)))
        self.assertEqual([], list(iter_old_brand_matches("")))
        self.assertEqual([], list(iter_old_brand_identifiers(None)))
        self.assertEqual([], list(iter_old_brand_identifiers("")))
        self.assertEqual([], extract_svg_sources(None))
        self.assertEqual([], extract_svg_sources(""))
        self.assertEqual([], identifier_ranges(None, "agentfield"))
        self.assertEqual([], identifier_ranges("", "agentfield"))
        self.assertEqual([], identifier_ranges("AgentField", ""))
        self.assertEqual([(5, 15)], identifier_ranges("link agentfield docs", "agentfield"))
        self.assertEqual([(0, 10), (11, 21)], identifier_ranges("agentfield agentfield", "agentfield"))
        self.assertEqual([], list(iter_old_brand_identifiers("AgentField_Python.json")))
        self.assertEqual([], list(iter_old_brand_identifiers("AgentField_TypeScript.json")))
        self.assertEqual([], list(iter_old_brand_identifiers("agentfield_server")))

    def test_rag_evaluation_client_accepts_legacy_server_keyword(self) -> None:
        module = load_rag_eval_client_module()
        evaluator = module.RAGEvaluator(**{"agent" "field_server": "http://legacy.example"})
        self.addCleanup(evaluator._client.close)
        self.assertEqual("http://legacy.example", evaluator.base_url)

    def test_rag_evaluation_client_normalizes_urls_and_rejects_unknown_kwargs(self) -> None:
        module = load_rag_eval_client_module()

        evaluator = module.RAGEvaluator(
            silmari_server="http://explicit.example/",
            **{"agent" "field_server": "http://legacy.example///"},
        )
        self.addCleanup(evaluator._client.close)
        self.assertEqual("http://legacy.example", evaluator.base_url)

        with self.assertRaisesRegex(TypeError, "Unexpected keyword arguments: unexpected"):
            module.RAGEvaluator(unexpected="value")

    def test_async_rag_evaluation_client_accepts_legacy_server_keyword(self) -> None:
        module = load_rag_eval_client_module()
        evaluator = module.AsyncRAGEvaluator(
            silmari_server="http://explicit.example/",
            **{"agent" "field_server": "http://legacy.example///"},
        )
        self.addCleanup(lambda: asyncio.run(evaluator.close()))
        self.assertEqual("http://legacy.example", evaluator.base_url)

        with self.assertRaisesRegex(TypeError, "Unexpected keyword arguments: unexpected"):
            module.AsyncRAGEvaluator(unexpected="value")

    def test_manifest_parser_errors_when_required_heading_is_missing(self) -> None:
        broken_manifest = textwrap.dedent(
            """
            # Silmari Rebrand Manifest

            ## Summary
            """
        ).strip()

        with self.assertRaisesRegex(ValueError, "Missing heading: ## Audited Files"):
            _parse_markdown_table(broken_manifest, "## Audited Files")

    def test_manifest_parser_errors_when_required_rows_are_missing(self) -> None:
        incomplete_table = textwrap.dedent(
            """
            # Silmari Rebrand Manifest

            ## Audited Files
            | Path | Action | Verification |
            |------|--------|--------------|
            """
        ).strip()

        with self.assertRaisesRegex(ValueError, "Missing rows for ## Audited Files"):
            _parse_markdown_table(incomplete_table, "## Audited Files")

    def test_rag_evaluation_assets_are_renamed_and_consistent(self) -> None:
        page_text = read_text("examples/python_agent_nodes/rag_evaluation/ui/app/page.tsx")
        powered_by_text = read_text(
            "examples/python_agent_nodes/rag_evaluation/ui/components/PoweredBy.tsx"
        )

        self.assertNotIn("/agentfield-", page_text)
        self.assertNotIn("/agentfield-", powered_by_text)
        self.assertIn('/silmari-icon-dark.svg', page_text)
        self.assertIn('/silmari-logo-dark.svg', powered_by_text)
        self.assertIn('alt="Silmari"', page_text)
        self.assertIn('alt="Silmari"', powered_by_text)

        expected_assets = {
            "silmari-icon-dark.svg",
            "silmari-icon-light.svg",
            "silmari-logo-dark.svg",
            "silmari-logo-light.svg",
        }
        actual_assets = {path.name for path in RAG_UI_PUBLIC.glob("*.svg")}
        self.assertEqual(expected_assets, actual_assets)
        self.assertFalse(any(name.startswith("agentfield-") for name in actual_assets))

        for rel_path in sorted(
            "examples/python_agent_nodes/rag_evaluation/ui/public/" + name
            for name in actual_assets
        ):
            with self.subTest(path=rel_path):
                text = read_text(rel_path)
                self.assertFalse(list(iter_old_brand_matches(text)))
                if "logo" in rel_path:
                    self.assertIn("Silmari", text)

        referenced_assets = extract_svg_sources(page_text + powered_by_text)
        for asset_name in referenced_assets:
            with self.subTest(asset=asset_name):
                self.assertTrue((RAG_UI_PUBLIC / asset_name).exists())

    def test_rag_evaluation_ui_lint_script_runs_typescript_check(self) -> None:
        package_json = json.loads(
            read_text("examples/python_agent_nodes/rag_evaluation/ui/package.json")
        )

        self.assertEqual("tsc --noEmit", package_json["scripts"]["lint"])

    def test_benchmark_assets_keep_historical_filenames_but_render_silmari_labels(self) -> None:
        analyze_text = read_text("examples/benchmarks/100k-scale/analyze.py")
        readme_text = read_text("examples/benchmarks/100k-scale/README.md")
        self.assertIn("Silmari Benchmark Comparison", analyze_text)
        self.assertIn('Silmari (Python)', analyze_text)
        self.assertIn("| Silmari Python |", readme_text)

        results_dir = REPO_ROOT / "examples" / "benchmarks" / "100k-scale" / "results"
        historical_files = {
            "AgentField_Go.json",
            "AgentField_Python.json",
            "AgentField_TypeScript.json",
        }
        for file_name in historical_files:
            with self.subTest(file=file_name):
                self.assertTrue((results_dir / file_name).exists())
        self.assertTrue((results_dir / "benchmark_summary.png").exists())

    def test_curated_example_surfaces_with_old_brand_tokens_are_audited(self) -> None:
        audited_files, _ = parse_manifest()
        missing_audits: list[str] = []

        for path in iter_curated_example_surface_paths():
            rel_path = path.relative_to(REPO_ROOT).as_posix()
            text = path.read_text(encoding="utf-8")
            has_old_brand = bool(
                list(iter_old_brand_identifiers(rel_path))
                or list(iter_old_brand_identifiers(text))
            )
            if not has_old_brand:
                continue
            if rel_path not in audited_files:
                missing_audits.append(rel_path)

        self.assertEqual(
            [],
            missing_audits,
            "curated user-facing example surfaces that still contain legacy-brand"
            " identifiers are expected to be listed in docs/silmari-rebrand-manifest.md"
            " under ## Audited Files",
        )

    def test_manifest_rows_exist_for_audited_files_with_preserved_identifiers(self) -> None:
        audited_files, preserved_rows = parse_manifest()
        preserved_by_path: dict[str, list[PreservedIdentifier]] = {}
        for row in preserved_rows:
            preserved_by_path.setdefault(row.path, []).append(row)

        missing_rows: list[str] = []

        for rel_path, audited in audited_files.items():
            if not rel_path.startswith(("examples/", "assets/")):
                continue
            if rel_path.endswith((".png", ".jpg", ".jpeg", ".gif", ".webp")):
                continue
            text = read_text(rel_path)
            has_old_brand = bool(
                list(iter_old_brand_identifiers(rel_path))
                or list(iter_old_brand_identifiers(text))
            )
            if not has_old_brand:
                continue
            if not preserved_by_path.get(rel_path):
                missing_rows.append(rel_path)

        self.assertEqual([], missing_rows)

    def test_manifested_identifiers_cover_old_brand_tokens_in_audited_example_files(self) -> None:
        audited_files, preserved_rows = parse_manifest()
        preserved_by_path: dict[str, list[PreservedIdentifier]] = {}
        for row in preserved_rows:
            preserved_by_path.setdefault(row.path, []).append(row)

        uncovered_matches: list[str] = []
        for rel_path in audited_files:
            if not rel_path.startswith(("examples/", "assets/utm-links.csv")):
                continue
            if rel_path.endswith((".png", ".jpg", ".jpeg", ".gif", ".webp")):
                continue

            text = read_text(rel_path)
            rows = preserved_by_path.get(rel_path, [])
            for occurrence in iter_old_brand_occurrences(rel_path, text):
                if any(
                    manifested_identifier_covers_occurrence(
                        row.identifier,
                        occurrence,
                        rel_path=rel_path,
                        text=text,
                    )
                    for row in rows
                ):
                    continue
                uncovered_matches.append(
                    f"{rel_path} [{occurrence.source}@{occurrence.start}]:"
                    f" {occurrence.identifier}"
                )

        self.assertEqual([], uncovered_matches)

    def test_manifest_uses_allowed_categories_and_concrete_reasons(self) -> None:
        _, preserved_rows = parse_manifest()
        example_rows = [
            row
            for row in preserved_rows
            if row.path.startswith("examples/") or row.path == "assets/utm-links.csv"
        ]

        self.assertTrue(example_rows)

        for row in example_rows:
            with self.subTest(path=row.path, identifier=row.identifier):
                self.assertIn(row.category, ALLOWED_CATEGORIES)
                self.assertGreaterEqual(len(row.reason.strip()), 20)

    def test_historical_benchmark_identifiers_are_manifested_as_historical_records(self) -> None:
        _, preserved_rows = parse_manifest()
        preserved_lookup = {
            (row.path, row.identifier): row
            for row in preserved_rows
        }

        expected_rows = [
            (
                "examples/benchmarks/100k-scale/analyze.py",
                "AgentField_Go",
            ),
            (
                "examples/benchmarks/100k-scale/analyze.py",
                "AgentField_TypeScript",
            ),
            (
                "examples/benchmarks/100k-scale/analyze.py",
                "AgentField_Python",
            ),
            (
                "examples/benchmarks/100k-scale/run_benchmarks.sh",
                "AgentField_Go.json",
            ),
            (
                "examples/benchmarks/100k-scale/run_benchmarks.sh",
                "AgentField_Python.json",
            ),
            (
                "examples/benchmarks/100k-scale/results/AgentField_Go.json",
                "AgentField_Go.json",
            ),
            (
                "examples/benchmarks/100k-scale/results/AgentField_Go.json",
                "AgentField",
            ),
            (
                "examples/benchmarks/100k-scale/results/AgentField_Python.json",
                "AgentField_Python.json",
            ),
            (
                "examples/benchmarks/100k-scale/results/AgentField_Python.json",
                "AgentField",
            ),
            (
                "examples/benchmarks/100k-scale/results/AgentField_TypeScript.json",
                "AgentField_TypeScript.json",
            ),
            (
                "examples/benchmarks/100k-scale/results/AgentField_TypeScript.json",
                "AgentField",
            ),
        ]

        for key in expected_rows:
            with self.subTest(path=key[0], identifier=key[1]):
                row = preserved_lookup.get(key)
                self.assertIsNotNone(row)
                self.assertEqual("historical-record", row.category)
                self.assertGreaterEqual(len(row.reason.strip()), 20)

    def test_runtime_compatibility_identifiers_are_manifested(self) -> None:
        _, preserved_rows = parse_manifest()
        runtime_rows = [
            row
            for row in preserved_rows
            if row.identifier == "agentfield_server"
            and row.path.startswith("examples/")
            and row.path != "examples/tests/test_silmari_branding.py"
        ]

        self.assertTrue(runtime_rows)

        for row in runtime_rows:
            with self.subTest(path=row.path, identifier=row.identifier):
                self.assertEqual("json-field", row.category)
                self.assertGreaterEqual(len(row.reason.strip()), 20)

    def test_historical_benchmark_filename_reasons_are_concrete(self) -> None:
        _, preserved_rows = parse_manifest()
        filename_rows = [
            row
            for row in preserved_rows
            if row.path.startswith("examples/benchmarks/100k-scale/")
            and row.identifier.endswith(".json")
        ]

        self.assertTrue(filename_rows)

        for row in filename_rows:
            with self.subTest(path=row.path, identifier=row.identifier):
                self.assertEqual("historical-record", row.category)
                self.assertRegex(
                    row.reason.lower(),
                    r"filename|artifact|result file",
                )

    def test_test_fixture_old_brand_identifiers_are_manifested(self) -> None:
        _, preserved_rows = parse_manifest()
        fixture_rows = [
            row
            for row in preserved_rows
            if row.path == "examples/tests/test_silmari_branding.py" and row.category == "test-fixture"
        ]

        self.assertTrue(
            fixture_rows,
            "examples/tests/test_silmari_branding.py keeps old-brand fixture strings and should be"
            " manifested as test-fixture data",
        )


if __name__ == "__main__":
    unittest.main()
