"""MW Phase 4 B2 — the conformance surface, worked for the SECOND frozen contract.

Phase 3 B4 registered ``com.silmari.reel.completed.v1`` in the frozen registry
(see test_handoff_registry_reel.py: registry + validation + freeze). This module
proves the SHIPPED conformance surface now generalizes to that second contract
end-to-end, so a reel-producing app (e.g. reel-af) can run it as its onboarding
gate with ZERO per-pair code:

  * the R1-R4 provenance conformance kit certifies a reel node (the 6 pattern
    contracts: R1/R4 = C-Own, R2 = C-Notification, R3 = C-Idempotent,
    C-Outbox + C-AtLeastOnce delegated, C-Correlation precondition), and
  * the frozen registry validates a ``reel.completed.v1`` event, pinned to the
    frozen schema version, and rejects any off-contract event (fail-closed).

Mirrors test_conformance_kit.py (which works a generic/demo node) for the reel
contract; complements test_handoff_registry_reel.py (registry + freeze only).
"""

import pytest

from agentfield.handoff import ValidationError, registry, validate
from agentfield.provenance_handoff import CompletionEnvelope, OutboxCompletionEmitter
from agentfield.testing.conformance import (
    _FakeBus,
    _SpyForeignStore,
    run_provenance_conformance,
)

REEL_TYPE = "com.silmari.reel.completed.v1"
RESEARCH_TYPE = "com.silmari.research.completed.v1"

# Frozen reel.completed.v1 DTO — matches contracts/com.silmari.reel.completed/v1.schema.json.
GOLDEN_REEL_DTO = {
    "run_id": "3f6d1a90-0000-4000-8000-000000000abc",
    "status": "succeeded",
    "reel_ref": "cp-execution://exec_abc123/result",
    "source_execution_id": "exec_abc123",
    "duration_s": 12.5,
    "beat_count": 5,
}


class _ReelConsumer:
    """A reel app's consumer: returns its OWN artifact id (R1), never a foreign write."""

    def on_completed(self, env: CompletionEnvelope) -> str:
        return f"reel_{env.subject}"


class _MemStore:
    """Correct local provenance store (R4 edge)."""

    def __init__(self) -> None:
        self._m: dict[str, tuple[str, str]] = {}

    def record(self, *, artifact_id: str, foreign_run_id: str, source: str) -> None:
        self._m[artifact_id] = (foreign_run_id, source)

    def has(self, artifact_id: str) -> bool:
        return artifact_id in self._m


# ─────────── the second contract is live alongside the first ───────────


def test_both_frozen_contracts_registered():
    type_versions = {e.type_version for e in registry.list()}
    assert RESEARCH_TYPE in type_versions
    assert REEL_TYPE in type_versions  # the worked SECOND contract


# ─────────── onboarding gate, part 1: R1-R4 kit certifies a reel node ───────────


def test_conformance_kit_certifies_a_reel_node():
    bus = _FakeBus()
    res = run_provenance_conformance(
        OutboxCompletionEmitter(bus, node_id="reel-af"),
        _ReelConsumer(),
        _MemStore(),
        domain="reel",
    )
    assert res.ok, f"expected ok, got failed={res.failed}, detail={res.detail}"
    assert res.passed == {"R1", "R2", "R3", "R4"}


def test_kit_reports_all_six_contracts_for_reel_node():
    bus = _FakeBus()
    res = run_provenance_conformance(
        OutboxCompletionEmitter(bus, node_id="reel-af"),
        _ReelConsumer(),
        _MemStore(),
        domain="reel",
    )
    # C-Own (R1/R4), C-Notification (R2), C-Idempotent (R3) proved above; the
    # remaining three of the six pattern contracts are delegated/precondition.
    assert "C-Outbox" in res.delegated
    assert "C-AtLeastOnce" in res.delegated
    assert "C-Correlation" in res.precondition


def test_reel_node_writes_no_foreign_table():
    spy = _SpyForeignStore()
    bus = _FakeBus()
    res = run_provenance_conformance(
        OutboxCompletionEmitter(bus, node_id="reel-af"),
        _ReelConsumer(),
        _MemStore(),
        domain="reel",
        foreign_store=spy,
    )
    assert res.ok
    assert len(spy.writes) == 0  # C-Own: only the reel app's own store is written


# ─────── onboarding gate, part 2: reel.completed.v1 pinned to the frozen version ───────


def test_reel_completed_event_validates_against_frozen_version():
    # Pins to the frozen schema resolved from the registry (no per-pair code).
    validate({"type": REEL_TYPE, "data": GOLDEN_REEL_DTO})


def test_off_contract_reel_event_extra_field_rejected():
    # An extra field would leak the mutable body — fail closed on the frozen contract.
    bad = {**GOLDEN_REEL_DTO, "video_path": "/leak/reel.mp4"}
    with pytest.raises(ValidationError):
        validate({"type": REEL_TYPE, "data": bad})


def test_off_contract_reel_event_bad_enum_rejected():
    bad = {**GOLDEN_REEL_DTO, "status": "running"}
    with pytest.raises(ValidationError):
        validate({"type": REEL_TYPE, "data": bad})
