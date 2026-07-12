"""Behavior 3: OutboxCompletionEmitter reference adapter over the shipped durable outbox (R2)."""
import json
import pathlib

from agentfield.provenance_handoff import (
    OutboxCompletionEmitter,
    build_envelope,
)

GOLDEN_PATH = pathlib.Path(__file__).resolve().parents[3] / "testdata" / "completion_envelope.golden.json"


class _FakeBus:
    """OutboxStore-shaped double — stands in for DurableExecutionBus.Publish."""

    def __init__(self):
        self.appended: list[dict] = []

    def publish(self, event: dict) -> None:
        self.appended.append(event)


def test_emit_appends_one_completed_event():
    bus = _FakeBus()
    emitter = OutboxCompletionEmitter(bus, node_id="deep-research")
    ce_id = emitter.emit_completed(
        domain="research",
        correlation_id="exec_9",
        result_ref="pkg://x",
        snapshot={"prompt": "p", "document_id": "d"},
    )
    assert len(bus.appended) == 1
    ev = bus.appended[0]
    assert ev["type"] == "com.silmari.research.completed"
    assert ev["execution_id"] == "exec_9"
    assert '"result_ref":"pkg://x"' in ev["data"] or '"result_ref": "pkg://x"' in ev["data"]
    assert isinstance(ce_id, str) and len(ce_id) > 0


def test_emit_propagates_publish_error():
    class _BrokenBus:
        def publish(self, event: dict) -> None:
            raise IOError("disk full")

    emitter = OutboxCompletionEmitter(_BrokenBus(), node_id="x")
    try:
        emitter.emit_completed(
            domain="test", correlation_id="e1", result_ref="r", snapshot={}
        )
        assert False, "should propagate the publish error"
    except IOError:
        pass


def test_two_emits_same_correlation_appends_two():
    bus = _FakeBus()
    emitter = OutboxCompletionEmitter(bus, node_id="x")
    emitter.emit_completed(domain="t", correlation_id="e1", result_ref="r", snapshot={})
    emitter.emit_completed(domain="t", correlation_id="e1", result_ref="r", snapshot={})
    assert len(bus.appended) == 2


def test_envelope_json_matches_shared_golden():
    golden = json.loads(GOLDEN_PATH.read_text())
    snapshot = {k: v for k, v in golden["data"].items() if k != "result_ref"}
    env = build_envelope(
        node_id="deep-research",
        domain="research",
        correlation_id=golden["subject"],
        result_ref=golden["data"]["result_ref"],
        snapshot=snapshot,
        dataschema=golden.get("dataschema"),
    )
    got = json.loads(env.to_json())
    got.pop("id")
    expected = dict(golden)
    expected.pop("id")
    assert got == expected, f"Python envelope drifted from golden:\ngot={got}\nexpected={expected}"


def test_emitter_creates_no_new_bus():
    """Grep-equivalent: the adapter takes a publisher by injection, never constructs one."""
    import inspect
    src = inspect.getsource(OutboxCompletionEmitter)
    assert "EventBus(" not in src
    assert "GlobalExecutionEventBus" not in src
