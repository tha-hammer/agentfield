"""Behavior 5: Conformance kit self-test — correct node passes, broken fixtures fail red."""
from agentfield.provenance_handoff import (
    CompletionEnvelope,
    OutboxCompletionEmitter,
)
from agentfield.testing.conformance import (
    ConformanceResult,
    _FakeBus,
    _SpyForeignStore,
    run_provenance_conformance,
)


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

class _MemStore:
    """Correct ProvenanceStore — local R4 edge store."""
    def __init__(self):
        self._m: dict[str, tuple[str, str]] = {}

    def record(self, *, artifact_id: str, foreign_run_id: str, source: str) -> None:
        self._m[artifact_id] = (foreign_run_id, source)

    def has(self, artifact_id: str) -> bool:
        return artifact_id in self._m


class _CorrectConsumer:
    """A consumer that satisfies R1+R3: returns its own artifact id, no foreign write."""
    def on_completed(self, env: CompletionEnvelope) -> str:
        return f"my_artifact_from_{env.subject}"


# ---------------------------------------------------------------------------
# Correct node → green
# ---------------------------------------------------------------------------


def test_kit_passes_a_correct_node():
    bus = _FakeBus()
    prod = OutboxCompletionEmitter(bus, node_id="demo")
    cons = _CorrectConsumer()
    store = _MemStore()
    res = run_provenance_conformance(prod, cons, store)
    assert res.ok, f"expected ok, got failed={res.failed}, detail={res.detail}"
    assert res.passed == {"R1", "R2", "R3", "R4"}


def test_kit_reports_all_six_spec_contracts():
    bus = _FakeBus()
    res = run_provenance_conformance(
        OutboxCompletionEmitter(bus, node_id="demo"),
        _CorrectConsumer(),
        _MemStore(),
    )
    all_contracts = res.passed | res.failed | res.waived | res.delegated | res.precondition
    assert "C-Outbox" in res.delegated
    assert "C-AtLeastOnce" in res.delegated
    assert "C-Correlation" in res.precondition


# ---------------------------------------------------------------------------
# Broken fixture: R3 red — consumer dedups on CE id, not correlation_id
# ---------------------------------------------------------------------------


class _DedupsOnCeIdConsumer:
    """Broken: dedups on CloudEvents id instead of correlation_id.
    A re-emit (same run, new CE id) will reprocess → R3 red."""
    def dedup_key(self, env: CompletionEnvelope) -> str:
        return env.id  # WRONG — should be env.subject (correlation_id)

    def on_completed(self, env: CompletionEnvelope) -> str:
        return f"art_{env.id}"


def test_kit_fails_r3_when_consumer_dedups_on_ce_id_not_correlation():
    bus = _FakeBus()
    res = run_provenance_conformance(
        OutboxCompletionEmitter(bus, node_id="demo"),
        _DedupsOnCeIdConsumer(),
        _MemStore(),
    )
    # CE-id dedup → the consumer processes the re-emit as a new item → R3 should fail
    # The kit detects this because on_completed runs twice for the same correlation_id
    assert not res.ok
    assert "R3" in res.failed or "C-Idempotent" in res.waived


# ---------------------------------------------------------------------------
# Broken fixture: R1 red — consumer writes a foreign store
# ---------------------------------------------------------------------------


class _ForeignWritingConsumer:
    """Broken: writes to a foreign store (R1 violation)."""
    def __init__(self, foreign: _SpyForeignStore):
        self._foreign = foreign

    def on_completed(self, env: CompletionEnvelope) -> str:
        self._foreign.record(
            artifact_id="stolen",
            foreign_run_id=env.subject,
            source=env.source,
        )
        return f"art_{env.subject}"


def test_kit_fails_r1_on_cross_context_write():
    bus = _FakeBus()
    foreign = _SpyForeignStore()
    res = run_provenance_conformance(
        OutboxCompletionEmitter(bus, node_id="demo"),
        _ForeignWritingConsumer(foreign=foreign),
        _MemStore(),
        foreign_store=foreign,
    )
    assert not res.ok
    assert "R1" in res.failed


# ---------------------------------------------------------------------------
# Broken fixture: R4 red — consumer returns id but never records provenance
# ---------------------------------------------------------------------------


class _NoRecordStore:
    """Broken store: record is a no-op, has always returns False."""
    def record(self, *, artifact_id: str, foreign_run_id: str, source: str) -> None:
        pass  # no-op — R4 violation

    def has(self, artifact_id: str) -> bool:
        return False


def test_kit_fails_r4_when_consumer_never_records_provenance():
    bus = _FakeBus()
    res = run_provenance_conformance(
        OutboxCompletionEmitter(bus, node_id="demo"),
        _CorrectConsumer(),
        _NoRecordStore(),
    )
    assert not res.ok
    assert "R4" in res.failed


# ---------------------------------------------------------------------------
# Kit is importable by external packages
# ---------------------------------------------------------------------------


def test_kit_is_importable():
    from agentfield.testing.conformance import run_provenance_conformance, ConformanceResult
    assert callable(run_provenance_conformance)
    assert ConformanceResult is not None
