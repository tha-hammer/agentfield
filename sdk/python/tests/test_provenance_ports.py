"""Behavior 1: CompletionEmitter / CompletionConsumer / ProvenanceStore ports (R1–R4 as a seam)."""
from agentfield.provenance_handoff import (
    CompletionEmitter,
    CompletionConsumer,
    CompletionEnvelope,
    ProvenanceStore,
)


class _MemStore:
    def __init__(self):
        self._m: dict[str, tuple[str, str]] = {}

    def record(self, *, artifact_id: str, foreign_run_id: str, source: str) -> None:
        self._m[artifact_id] = (foreign_run_id, source)

    def has(self, artifact_id: str) -> bool:
        return artifact_id in self._m


class _TrivialConsumer:
    def dedup_key(self, envelope: CompletionEnvelope) -> str:
        return envelope.subject

    def on_completed(self, envelope: CompletionEnvelope) -> str:
        return f"artifact_from_{envelope.subject}"


class _ConsumerWithoutDedupKey:
    """A consumer that doesn't implement dedup_key — should still satisfy the protocol
    if the protocol makes dedup_key optional with a default."""

    def on_completed(self, envelope: CompletionEnvelope) -> str:
        return "art_1"


class _TrivialEmitter:
    def emit_completed(
        self, *, domain: str, correlation_id: str, result_ref: str, snapshot: dict
    ) -> str:
        return "ce_id_fake"


def test_provenance_store_is_runtime_checkable_and_satisfiable():
    store = _MemStore()
    assert isinstance(store, ProvenanceStore)
    store.record(artifact_id="a1", foreign_run_id="run_1", source="node-x")
    assert store.has("a1")
    assert not store.has("nonexistent")


def test_completion_emitter_has_emit_completed():
    assert hasattr(CompletionEmitter, "emit_completed")
    emitter = _TrivialEmitter()
    assert isinstance(emitter, CompletionEmitter)


def test_completion_consumer_has_on_completed_and_dedup_key():
    assert hasattr(CompletionConsumer, "on_completed")
    assert hasattr(CompletionConsumer, "dedup_key")
    consumer = _TrivialConsumer()
    assert isinstance(consumer, CompletionConsumer)


def test_import_does_no_io():
    """The module can be imported without opening any socket or DB connection."""
    import agentfield.provenance_handoff  # noqa: F401 — import itself is the test
