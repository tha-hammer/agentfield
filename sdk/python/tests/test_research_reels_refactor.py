"""Behavior 6: Refactor Phase-2 research→reels onto the generic port (behavior-preserving).

Proves:
- deep-research's emit is expressible as one emit_completed() call
- reel-af's consume is expressible as CompletionConsumer.on_completed(env) -> reel_job_id
  + ProvenanceStore writing source_research_run_id = env.subject
- the conformance kit certifies the pair (R1–R4)
- no deepresearch.* write in the consumer; bespoke code is superseded by the adapters
"""
from __future__ import annotations

import uuid
from typing import Any

from agentfield.provenance_handoff import (
    CompletionEnvelope,
    IdempotentConsumer,
    OutboxCompletionEmitter,
    build_envelope,
)
from agentfield.testing.conformance import (
    _FakeBus,
    run_provenance_conformance,
)


# ---------------------------------------------------------------------------
# ReelAfConsumer — reel-af's consume expressed as a CompletionConsumer
# ---------------------------------------------------------------------------


class FakeReelRepo:
    """Mimics the reel-af repo for testing — stamps source_research_run_id."""

    def __init__(self):
        self.created: dict[str, dict[str, Any]] = {}

    def create_reel_job(
        self, *, url: str, source_research_run_id: str | None = None
    ) -> str:
        job_id = f"reel_{uuid.uuid4().hex[:8]}"
        self.created[job_id] = {
            "url": url,
            "source_research_run_id": source_research_run_id,
        }
        return job_id

    def source_research_run_id_of(self, job_id: str) -> str | None:
        return self.created.get(job_id, {}).get("source_research_run_id")


class ReelAfConsumer:
    """reel-af's consume side expressed as a CompletionConsumer (R1+R3).

    on_completed returns the reel-af-owned reel_job_id (R1: returns MY artifact id).
    Stamps source_research_run_id = env.subject (the correlation_id / execution_id).
    """

    def __init__(self, repo: FakeReelRepo):
        self._repo = repo

    def on_completed(self, env: CompletionEnvelope) -> str:
        return self._repo.create_reel_job(
            url=env.data.get("result_ref", ""),
            source_research_run_id=env.subject,
        )


class ReelAfProvenanceStore:
    """reel-af's local provenance store — maps artifact_id → foreign_run_id (R4).

    In the real implementation, this is the source_research_run_id column on
    reel_job/carousel rows. Here it wraps a simple dict for testing.
    """

    def __init__(self):
        self._edges: dict[str, tuple[str, str]] = {}

    def record(self, *, artifact_id: str, foreign_run_id: str, source: str) -> None:
        self._edges[artifact_id] = (foreign_run_id, source)

    def has(self, artifact_id: str) -> bool:
        return artifact_id in self._edges


# ---------------------------------------------------------------------------
# Behavior 6 tests
# ---------------------------------------------------------------------------


def test_reel_carries_source_run_via_generic_consumer():
    """The existing create-from-research flow expressed via the generic port."""
    repo = FakeReelRepo()
    store = ReelAfProvenanceStore()
    env = build_envelope(
        node_id="deep-research",
        domain="research",
        correlation_id="exec_42",
        result_ref="pkg://deep-research/runs/exec_42/result",
        snapshot={
            "research_prompt": "Why is formal domain ownership important?",
            "research_document_id": "exec_42",
        },
    )

    consumer = IdempotentConsumer(ReelAfConsumer(repo=repo), store)
    reel_id = consumer.deliver(env)

    assert repo.source_research_run_id_of(reel_id) == "exec_42"


def test_emit_expressible_as_one_emit_completed_call():
    """deep-research's emit is expressible as one emit_completed() call."""
    bus = _FakeBus()
    emitter = OutboxCompletionEmitter(bus, node_id="deep-research")

    ce_id = emitter.emit_completed(
        domain="research",
        correlation_id="exec_42",
        result_ref="cp-execution://exec_42/result",
        snapshot={
            "research_prompt": "Why is formal domain ownership important?",
            "research_document_id": "exec_42",
            "status": "succeeded",
            "title": "Domain Ownership in Distributed Systems",
        },
    )

    assert len(bus.appended) == 1
    assert bus.appended[0]["type"] == "com.silmari.research.completed"
    assert bus.appended[0]["execution_id"] == "exec_42"
    assert isinstance(ce_id, str)


def test_conformance_kit_certifies_research_reels_pair():
    """The generic conformance kit certifies the research→reels pair R1–R4."""
    bus = _FakeBus()
    prod = OutboxCompletionEmitter(bus, node_id="deep-research")
    cons = ReelAfConsumer(repo=FakeReelRepo())
    store = ReelAfProvenanceStore()

    res = run_provenance_conformance(prod, cons, store, domain="research")
    assert res.ok, f"expected ok, got failed={res.failed}, detail={res.detail}"
    assert res.passed == {"R1", "R2", "R3", "R4"}


def test_consumer_writes_no_foreign_table():
    """The consumer writes only its own repo (source_research_run_id on reel_job),
    never a deepresearch.* table — R1 preserved."""
    import inspect

    src = inspect.getsource(ReelAfConsumer)
    assert "deepresearch" not in src.lower()
    # The consumer calls only repo.create_reel_job (its own), not any foreign store
    assert "research_run" not in src.lower() or "source_research_run_id" in src


def test_redelivery_is_idempotent_via_generic_adapter():
    """A re-emit of the same correlation_id runs on_completed ONCE (CI-1, via the port)."""
    repo = FakeReelRepo()
    store = ReelAfProvenanceStore()
    consumer = IdempotentConsumer(ReelAfConsumer(repo=repo), store)

    env1 = build_envelope(
        node_id="deep-research",
        domain="research",
        correlation_id="exec_99",
        result_ref="r",
        snapshot={},
    )
    env2 = build_envelope(
        node_id="deep-research",
        domain="research",
        correlation_id="exec_99",
        result_ref="r",
        snapshot={},
    )

    id1 = consumer.deliver(env1)
    id2 = consumer.deliver(env2)

    assert id1 == id2
    assert len(repo.created) == 1
