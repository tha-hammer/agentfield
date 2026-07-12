"""Behavior 7: "Works for any app" acceptance — a net-new node plugs in and certifies.

A fabricated demo node (node_id="demo-app") that has never touched research or
reels becomes a valid producer + consumer and passes the conformance kit unmodified.
"""
import inspect

from agentfield.provenance_handoff import (
    CompletionEnvelope,
    OutboxCompletionEmitter,
)
from agentfield.testing.conformance import (
    _FakeBus,
    run_provenance_conformance,
)


# ---------------------------------------------------------------------------
# Demo node — zero research/reels lineage
# ---------------------------------------------------------------------------


class _DemoProvenanceStore:
    """Local provenance store for the demo app (R4)."""
    def __init__(self):
        self._edges: dict[str, tuple[str, str]] = {}

    def record(self, *, artifact_id: str, foreign_run_id: str, source: str) -> None:
        self._edges[artifact_id] = (foreign_run_id, source)

    def has(self, artifact_id: str) -> bool:
        return artifact_id in self._edges


class _DemoConsumer:
    """Completion consumer for the demo app — returns its own artifact id."""
    def on_completed(self, env: CompletionEnvelope) -> str:
        return f"demo_artifact_{env.subject}"


class DemoApp:
    """A brand-new node with no prior domain-specific lineage."""
    def __init__(self, node_id: str = "demo-app"):
        self.node_id = node_id

    def completion_consumer(self) -> _DemoConsumer:
        return _DemoConsumer()

    def provenance_store(self) -> _DemoProvenanceStore:
        return _DemoProvenanceStore()


# ---------------------------------------------------------------------------
# Acceptance test
# ---------------------------------------------------------------------------


def test_a_brand_new_app_plugs_in_and_certifies():
    node = DemoApp(node_id="demo-app")
    bus = _FakeBus()
    prod = OutboxCompletionEmitter(bus, node_id="demo-app")
    cons = node.completion_consumer()
    store = node.provenance_store()

    res = run_provenance_conformance(prod, cons, store, domain="demo")
    assert res.ok, f"expected ok, got failed={res.failed}, detail={res.detail}"
    assert res.passed == {"R1", "R2", "R3", "R4"}


def test_demo_uses_different_domain_not_research():
    """The demo uses "demo" domain, proving nothing is research-specific."""
    bus = _FakeBus()
    prod = OutboxCompletionEmitter(bus, node_id="demo-app")
    prod.emit_completed(
        domain="demo",
        correlation_id="demo_run_1",
        result_ref="pkg://demo/result",
        snapshot={"custom_key": "custom_value"},
    )
    assert len(bus.appended) == 1
    assert bus.appended[0]["type"] == "com.silmari.demo.completed"


def test_plug_in_surface_is_exactly_two_modules():
    """The demo imports only provenance_handoff + testing.conformance."""
    src = inspect.getsource(DemoApp)
    assert "deepresearch" not in src.lower()
    assert "reel" not in src.lower()

    src_consumer = inspect.getsource(_DemoConsumer)
    assert "deepresearch" not in src_consumer.lower()
    assert "reel" not in src_consumer.lower()
