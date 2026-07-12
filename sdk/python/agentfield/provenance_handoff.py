"""ProvenanceHandoff capability — reusable Node-Contract ports for cross-app handoff (R1–R4).

Ports (developer surface):
  CompletionEmitter  — R2: emit a past-tense <domain>.completed event
  CompletionConsumer — R1+R3: consume idempotently, return YOUR artifact id
  ProvenanceStore    — R1+R4: record a local lineage edge (artifact_id → foreign_run_id)

Value objects:
  CompletionEnvelope — CloudEvents 1.0 envelope for <domain>.completed

Adapters (reference implementations):
  OutboxCompletionEmitter — emitter over the shipped DurableExecutionBus
  IdempotentConsumer      — dedup on correlation_id + at-most-once provenance edge

The SDK owns NO persistent data (specs/agent-runtime-sdk.domain.md §7).
Each consuming node supplies its own ProvenanceStore adapter backed by its own table.
"""

from __future__ import annotations

import json
import uuid
from dataclasses import dataclass, field
from typing import Any, Mapping, Protocol, runtime_checkable

# ---------------------------------------------------------------------------
# Constants
# ---------------------------------------------------------------------------

CE_SPECVERSION = "1.0"
CE_TYPE_PREFIX = "com.silmari."
CE_TYPE_SUFFIX = ".completed"
MAX_SNAPSHOT_DEPTH = 2
MAX_SNAPSHOT_BYTES = 8 * 1024

# ---------------------------------------------------------------------------
# Errors
# ---------------------------------------------------------------------------


class EnvelopeError(Exception):
    """Raised when a CompletionEnvelope cannot be constructed (R2 violation)."""


class MalformedEnvelopeError(Exception):
    """Raised when raw data cannot be parsed into a CompletionEnvelope."""


# ---------------------------------------------------------------------------
# Value objects
# ---------------------------------------------------------------------------


@dataclass(frozen=True)
class CompletionEnvelope:
    """CloudEvents 1.0 envelope for a <domain>.completed handoff event.

    subject == correlation_id == execution_id — the single cross-app join key.
    """

    specversion: str
    id: str
    source: str
    type: str
    subject: str
    data: dict[str, Any]
    dataschema: str | None = None

    def to_json(self) -> str:
        d: dict[str, Any] = {
            "specversion": self.specversion,
            "id": self.id,
            "source": self.source,
            "type": self.type,
            "subject": self.subject,
            "data": self.data,
        }
        if self.dataschema is not None:
            d["dataschema"] = self.dataschema
        return json.dumps(d, sort_keys=True, separators=(",", ":"))

    @classmethod
    def from_json(cls, raw: str) -> CompletionEnvelope:
        try:
            d = json.loads(raw)
        except (json.JSONDecodeError, TypeError) as exc:
            raise MalformedEnvelopeError(f"invalid JSON: {exc}") from exc
        required = {"specversion", "id", "source", "type", "subject", "data"}
        missing = required - set(d)
        if missing:
            raise MalformedEnvelopeError(f"missing required fields: {missing}")
        return cls(
            specversion=d["specversion"],
            id=d["id"],
            source=d["source"],
            type=d["type"],
            subject=d["subject"],
            data=d["data"],
            dataschema=d.get("dataschema"),
        )


# ---------------------------------------------------------------------------
# Envelope builder
# ---------------------------------------------------------------------------


def _check_snapshot_depth(value: Any, current_depth: int, max_depth: int) -> None:
    if isinstance(value, dict):
        if current_depth > max_depth:
            raise EnvelopeError(
                f"snapshot depth {current_depth} exceeds MAX_SNAPSHOT_DEPTH ({max_depth})"
            )
        for v in value.values():
            _check_snapshot_depth(v, current_depth + 1, max_depth)
    elif isinstance(value, list):
        if current_depth > max_depth:
            raise EnvelopeError(
                f"snapshot depth {current_depth} exceeds MAX_SNAPSHOT_DEPTH ({max_depth})"
            )
        for item in value:
            _check_snapshot_depth(item, current_depth + 1, max_depth)
    elif not isinstance(value, (str, int, float, bool, type(None))):
        raise EnvelopeError(
            f"snapshot contains non-primitive value of type {type(value).__name__}"
        )


def build_envelope(
    *,
    node_id: str,
    domain: str,
    correlation_id: str,
    result_ref: str,
    snapshot: Mapping[str, Any],
    dataschema: str | None = None,
) -> CompletionEnvelope:
    """Build a valid CloudEvents CompletionEnvelope. Raises EnvelopeError on R2 violation."""
    if not correlation_id:
        raise EnvelopeError("correlation_id must be non-empty (it is the cross-app join key)")

    data = {"result_ref": result_ref, **snapshot}

    try:
        serialized = json.dumps(data, sort_keys=True)
    except (TypeError, ValueError) as exc:
        raise EnvelopeError(f"snapshot is not JSON-serializable: {exc}") from exc

    _check_snapshot_depth(data, 1, MAX_SNAPSHOT_DEPTH)

    if len(serialized) > MAX_SNAPSHOT_BYTES:
        raise EnvelopeError(
            f"snapshot serialized size ({len(serialized)} bytes) exceeds "
            f"MAX_SNAPSHOT_BYTES ({MAX_SNAPSHOT_BYTES})"
        )

    return CompletionEnvelope(
        specversion=CE_SPECVERSION,
        id=str(uuid.uuid4()),
        source=node_id,
        type=f"{CE_TYPE_PREFIX}{domain}{CE_TYPE_SUFFIX}",
        subject=correlation_id,
        data=data,
        dataschema=dataschema,
    )


# ---------------------------------------------------------------------------
# Ports (R1–R4 as @runtime_checkable Protocols)
# ---------------------------------------------------------------------------


@runtime_checkable
class ProvenanceStore(Protocol):
    """R1+R4: local lineage edge store. Each consuming node owns its own table."""

    def record(self, *, artifact_id: str, foreign_run_id: str, source: str) -> None: ...
    def has(self, artifact_id: str) -> bool: ...


@runtime_checkable
class CompletionEmitter(Protocol):
    """R2: the single entrypoint for emitting a past-tense <domain>.completed event."""

    def emit_completed(
        self,
        *,
        domain: str,
        correlation_id: str,
        result_ref: str,
        snapshot: Mapping[str, Any],
    ) -> str: ...


@runtime_checkable
class CompletionConsumer(Protocol):
    """R1+R3: consume a completion event idempotently. Returns the consumer's OWN artifact id."""

    def dedup_key(self, envelope: CompletionEnvelope) -> str: ...
    def on_completed(self, envelope: CompletionEnvelope) -> str: ...


# ---------------------------------------------------------------------------
# Reference adapter: OutboxCompletionEmitter (R2, over the shipped DurableExecutionBus)
# ---------------------------------------------------------------------------


class OutboxCompletionEmitter:
    """Reference emitter adapter wrapping an OutboxStore-shaped publisher.

    One emit_completed → one ExecutionEvent appended via the existing
    DurableExecutionBus.Publish path (append-durable-first). Creates no new bus.
    """

    def __init__(self, publisher: Any, *, node_id: str) -> None:
        self._publisher = publisher
        self._node_id = node_id

    def emit_completed(
        self,
        *,
        domain: str,
        correlation_id: str,
        result_ref: str,
        snapshot: Mapping[str, Any],
    ) -> str:
        env = build_envelope(
            node_id=self._node_id,
            domain=domain,
            correlation_id=correlation_id,
            result_ref=result_ref,
            snapshot=snapshot,
        )
        self._publisher.publish({
            "type": env.type,
            "execution_id": env.subject,
            "agent_node_id": self._node_id,
            "data": env.to_json(),
        })
        return env.id


# ---------------------------------------------------------------------------
# Reference adapter: IdempotentConsumer (R3+R4, dedup on correlation_id)
# ---------------------------------------------------------------------------


class IdempotentConsumer:
    """Dedup on correlation_id (CI-1) + at-most-once provenance edge (CI-2).

    The canonical dedup key is env.subject (correlation_id == execution_id).
    A producer re-emit (new CE id, same correlation_id) is a no-op.
    The in-memory set is only a fast front; the durable guard is store.has(key).
    """

    def __init__(self, inner: Any, store: ProvenanceStore) -> None:
        self._inner = inner
        self._store = store
        self._processed: dict[str, str] = {}

    def dedup_key(self, envelope: CompletionEnvelope) -> str:
        if hasattr(self._inner, "dedup_key"):
            return self._inner.dedup_key(envelope)
        return envelope.subject

    def deliver(self, envelope: CompletionEnvelope) -> str:
        key = self.dedup_key(envelope)

        if self._store.has(key):
            return self._processed.get(key, key)

        if key in self._processed:
            return self._processed[key]

        artifact_id = self._inner.on_completed(envelope)
        self._processed[key] = artifact_id
        self._store.record(
            artifact_id=key,
            foreign_run_id=envelope.subject,
            source=envelope.source,
        )
        return artifact_id

    def processed_count(self) -> int:
        return len(self._processed)


def deliver_raw(consumer: IdempotentConsumer, raw_data: str) -> str:
    """Parse raw JSON data into an envelope and deliver. Raises MalformedEnvelopeError on poison."""
    envelope = CompletionEnvelope.from_json(raw_data)
    return consumer.deliver(envelope)
