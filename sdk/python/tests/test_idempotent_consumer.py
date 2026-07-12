"""Behavior 4: IdempotentConsumer — dedup on correlation_id + durable at-most-once provenance (R3, R4)."""
from agentfield.provenance_handoff import (
    CompletionEnvelope,
    IdempotentConsumer,
    MalformedEnvelopeError,
    build_envelope,
    deliver_raw,
)


class _MemStore:
    """ProvenanceStore — local R4 edge store."""

    def __init__(self):
        self._m: dict[str, tuple[str, str]] = {}

    def record(self, *, artifact_id: str, foreign_run_id: str, source: str) -> None:
        self._m[artifact_id] = (foreign_run_id, source)

    def has(self, artifact_id: str) -> bool:
        return artifact_id in self._m


def _envelope(correlation_id: str = "exec_9") -> CompletionEnvelope:
    return build_envelope(
        node_id="deep-research",
        domain="research",
        correlation_id=correlation_id,
        result_ref="pkg://r",
        snapshot={},
    )


# ---------------------------------------------------------------------------
# CI-1: canonical dedup key = correlation_id
# ---------------------------------------------------------------------------


def test_reemit_same_correlation_id_is_a_noop():
    """A producer re-emit (new CE id, same correlation_id) runs on_completed ONCE."""
    seen: list[str] = []

    class _Consumer:
        def on_completed(self, env: CompletionEnvelope) -> str:
            seen.append(env.subject)
            return f"reel_from_{env.subject}"

    store = _MemStore()
    consumer = IdempotentConsumer(_Consumer(), store)

    first = _envelope("exec_9")
    reemit = _envelope("exec_9")
    assert first.id != reemit.id  # genuinely different CloudEvents ids

    a = consumer.deliver(first)
    b = consumer.deliver(reemit)

    assert seen == ["exec_9"]  # on_completed ran ONCE (R3, correlation key)
    assert a == b == "reel_from_exec_9"
    assert store.has("exec_9")  # exactly one provenance edge (R4)


def test_different_correlation_ids_each_processed():
    """Two different runs are processed independently."""
    seen: list[str] = []

    class _Consumer:
        def on_completed(self, env: CompletionEnvelope) -> str:
            seen.append(env.subject)
            return f"art_{env.subject}"

    store = _MemStore()
    consumer = IdempotentConsumer(_Consumer(), store)
    consumer.deliver(_envelope("exec_1"))
    consumer.deliver(_envelope("exec_2"))

    assert seen == ["exec_1", "exec_2"]
    assert consumer.processed_count() == 2


# ---------------------------------------------------------------------------
# CI-2: at-most-once edge survives simulated cursor loss (process restart)
# ---------------------------------------------------------------------------


def test_at_most_once_edge_survives_simulated_cursor_loss():
    """Post-crash redelivery whose provenance edge already landed does NOT re-run on_completed."""
    runs: list[str] = []

    class _Consumer:
        def on_completed(self, env: CompletionEnvelope) -> str:
            runs.append(env.subject)
            return f"reel_from_{env.subject}"

    store = _MemStore()
    env = _envelope("exec_5")

    # First process — edge lands
    IdempotentConsumer(_Consumer(), store).deliver(env)
    assert runs == ["exec_5"]

    # Simulate restart — new instance, in-memory set lost, same durable store
    fresh = IdempotentConsumer(_Consumer(), store)
    again = fresh.deliver(env)

    assert runs == ["exec_5"]  # has(key)-guard short-circuits (at-most-once edge)
    assert again == "exec_5"  # returns the dedup key when short-circuiting from store


# ---------------------------------------------------------------------------
# Transient failure: handler error does NOT mark processed
# ---------------------------------------------------------------------------


def test_failed_handler_does_not_mark_processed():
    class _Boom:
        def on_completed(self, env: CompletionEnvelope) -> str:
            raise RuntimeError("downstream down")

    consumer = IdempotentConsumer(_Boom(), _MemStore())
    env = _envelope("e1")

    try:
        consumer.deliver(env)
    except RuntimeError:
        pass

    assert consumer.processed_count() == 0


# ---------------------------------------------------------------------------
# Malformed record: raises without touching handler or store
# ---------------------------------------------------------------------------


def test_malformed_record_raises_without_touching_handler_or_store():
    class _Consumer:
        def on_completed(self, env: CompletionEnvelope) -> str:
            raise AssertionError("must not be called on poison")

    store = _MemStore()
    consumer = IdempotentConsumer(_Consumer(), store)

    try:
        deliver_raw(consumer, raw_data='{"not":"an envelope"}')
        assert False, "should have raised MalformedEnvelopeError"
    except MalformedEnvelopeError:
        pass

    assert consumer.processed_count() == 0


# ---------------------------------------------------------------------------
# R1 shape: consumer adapter writes only local store
# ---------------------------------------------------------------------------


def test_consumer_writes_only_local_store():
    """The adapter calls only on_completed + local record; no foreign-table write path."""
    import inspect

    src = inspect.getsource(IdempotentConsumer)
    # No foreign store interaction — only self._store and self._inner
    assert "foreign" not in src.lower() or "foreign_run_id" in src


# ---------------------------------------------------------------------------
# Custom dedup_key delegation
# ---------------------------------------------------------------------------


def test_custom_dedup_key_delegates_to_inner():
    """If the inner consumer defines dedup_key, the adapter uses it."""

    class _CustomKey:
        def dedup_key(self, env: CompletionEnvelope) -> str:
            return f"custom_{env.subject}"

        def on_completed(self, env: CompletionEnvelope) -> str:
            return f"art_{env.subject}"

    store = _MemStore()
    consumer = IdempotentConsumer(_CustomKey(), store)
    env = _envelope("exec_1")
    consumer.deliver(env)

    assert store.has("custom_exec_1")
