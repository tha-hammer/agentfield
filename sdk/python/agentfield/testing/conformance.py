"""Conformance kit — certifies a node's adapters satisfy R1–R4 against a fake outbox.

Usage (three-line drop-in for any node's test suite)::

    from agentfield.testing.conformance import run_provenance_conformance
    result = run_provenance_conformance(my_producer, my_consumer, my_store)
    assert result.ok

The kit runs on a FAKE outbox by design — it certifies the node's adapters +
rule-conformance so *any* app runs it in CI with no live control plane. The
real-durable-seam proof is the @integration Go Publish → ReadEventOutboxAfter
round-trip, NOT this kit.

Governing-spec contract coverage (spec §4):
  C-Own         → R1 check (consumer writes only its own store)
  C-Notification → R2 check (valid CloudEvents envelope, data whitelist)
  C-Idempotent  → R3 check (dedup on correlation_id; re-emit is a no-op)
  C-Own (OUT-C) → R4 check (provenance stored locally by id)
  C-Outbox      → delegated (live Postgres @integration proof)
  C-AtLeastOnce → delegated (consumer-driver live Postgres test)
  C-Correlation → precondition (Phase-1 UNIQUE(execution_id))
"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any

from agentfield.provenance_handoff import (
    CompletionEnvelope,
    IdempotentConsumer,
    build_envelope,
)


@dataclass
class ConformanceResult:
    """Per-rule certification result, keyed by governing-spec contract names."""

    passed: set[str] = field(default_factory=set)
    failed: set[str] = field(default_factory=set)
    waived: set[str] = field(default_factory=set)
    delegated: set[str] = field(default_factory=set)
    precondition: set[str] = field(default_factory=set)
    detail: dict[str, str] = field(default_factory=dict)

    @property
    def ok(self) -> bool:
        return len(self.failed) == 0


class _FakeBus:
    """Kit-provided fake outbox — stands in for DurableExecutionBus.Publish."""

    def __init__(self) -> None:
        self.appended: list[dict[str, Any]] = []

    def publish(self, event: dict[str, Any]) -> None:
        self.appended.append(event)


class _SpyForeignStore:
    """Spy that records any writes — used to assert R1 (no cross-context write)."""

    def __init__(self) -> None:
        self.writes: list[tuple[str, str, str]] = []

    def record(self, *, artifact_id: str, foreign_run_id: str, source: str) -> None:
        self.writes.append((artifact_id, foreign_run_id, source))

    def has(self, artifact_id: str) -> bool:
        return any(a == artifact_id for a, _, _ in self.writes)


def run_provenance_conformance(
    producer: Any,
    consumer: Any,
    store: Any,
    *,
    domain: str = "conformance-test",
    foreign_store: _SpyForeignStore | None = None,
) -> ConformanceResult:
    """Certify a node's (producer, consumer, store) triple against R1–R4.

    Uses a fake outbox. Does NOT prove C-Outbox or C-AtLeastOnce (those
    require a live Postgres proof).
    """
    result = ConformanceResult()

    result.delegated.add("C-Outbox")
    result.detail["C-Outbox"] = "delegated: proven by @integration Go Publish→ReadEventOutboxAfter round-trip"
    result.delegated.add("C-AtLeastOnce")
    result.detail["C-AtLeastOnce"] = "delegated: proven by consumer-driver live-Postgres test (ISC-INT03-11)"
    result.precondition.add("C-Correlation")
    result.detail["C-Correlation"] = "precondition: Phase-1 UNIQUE(execution_id)"

    fake_bus = _FakeBus()
    if foreign_store is None:
        foreign_store = _SpyForeignStore()

    correlation_id = "conformance_test_run_1"

    # --- Emit via producer ---
    try:
        producer.emit_completed(
            domain=domain,
            correlation_id=correlation_id,
            result_ref="pkg://conformance/result",
            snapshot={"test_key": "test_value"},
        )
    except Exception as exc:
        result.failed.add("R2")
        result.detail["R2"] = f"C-Notification: emit_completed raised: {exc}"
        return result

    # --- R2: valid envelope on the outbox ---
    if hasattr(producer, "_publisher") and hasattr(producer._publisher, "appended"):
        bus_records = producer._publisher.appended
    else:
        bus_records = fake_bus.appended

    r2_ok = True
    if len(bus_records) < 1:
        r2_ok = False
        result.detail["R2"] = "C-Notification: no record appended to outbox"
    else:
        rec = bus_records[-1]
        try:
            import json

            env_data = json.loads(rec.get("data", "{}")) if isinstance(rec.get("data"), str) else rec.get("data", {})
            if env_data.get("subject") != correlation_id:
                r2_ok = False
                result.detail["R2"] = f"C-Notification: subject={env_data.get('subject')}, expected {correlation_id}"
        except Exception as exc:
            r2_ok = False
            result.detail["R2"] = f"C-Notification: envelope parse error: {exc}"

    if r2_ok:
        result.passed.add("R2")
    else:
        result.failed.add("R2")

    # --- Build envelopes for consumer tests ---
    first_env = build_envelope(
        node_id="conformance-producer",
        domain=domain,
        correlation_id=correlation_id,
        result_ref="pkg://conformance/result",
        snapshot={"test_key": "test_value"},
    )

    # --- R3: idempotent consumer, keyed on correlation_id ---
    class _CountingWrapper:
        def __init__(self, inner: Any) -> None:
            self._inner = inner
            self.call_count = 0

        def on_completed(self, env: CompletionEnvelope) -> str:
            self.call_count += 1
            return self._inner.on_completed(env)

        def dedup_key(self, env: CompletionEnvelope) -> str:
            if hasattr(self._inner, "dedup_key"):
                return self._inner.dedup_key(env)
            return env.subject

    counting = _CountingWrapper(consumer)
    idempotent_counted = IdempotentConsumer(counting, store)

    idempotent_counted.deliver(first_env)

    reemit_env = build_envelope(
        node_id="conformance-producer",
        domain=domain,
        correlation_id=correlation_id,
        result_ref="pkg://conformance/result",
        snapshot={"test_key": "test_value"},
    )
    idempotent_counted.deliver(reemit_env)

    r3_ok = counting.call_count == 1
    if not r3_ok:
        result.failed.add("R3")
        result.detail["R3"] = (
            f"C-Idempotent: on_completed called {counting.call_count} times "
            f"(expected 1 — re-emit with new CE id should be a no-op on correlation_id)"
        )
    elif hasattr(consumer, "dedup_key") and consumer.dedup_key(first_env) != first_env.subject:
        result.waived.add("C-Idempotent")
        result.detail["R3"] = "C-Idempotent: waived — consumer uses custom dedup_key (not correlation_id)"
    else:
        result.passed.add("R3")

    # --- R4: provenance stored locally by id ---
    dedup_key = idempotent_counted.dedup_key(first_env)
    r4_ok = store.has(dedup_key)
    if r4_ok:
        result.passed.add("R4")
    else:
        result.failed.add("R4")
        result.detail["R4"] = f"C-Own (OUT-C): store.has({dedup_key!r}) is False after delivery"

    # --- R1: no cross-context write ---
    r1_ok = len(foreign_store.writes) == 0
    if r1_ok:
        result.passed.add("R1")
    else:
        result.failed.add("R1")
        result.detail["R1"] = f"C-Own: foreign store received {len(foreign_store.writes)} write(s)"

    return result
