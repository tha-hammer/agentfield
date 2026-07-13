"""Behavior 2: CloudEvents envelope — the standard <domain>.completed shape (R2)."""
from agentfield.provenance_handoff import (
    CompletionEnvelope,
    EnvelopeError,
    build_envelope,
)


def test_builds_valid_cloudevents():
    env = build_envelope(
        node_id="deep-research",
        domain="research",
        correlation_id="exec_123",
        result_ref="pkg://run/abc",
        snapshot={"prompt": "why FDO", "document_id": "doc_9", "title": "T"},
    )
    assert env.specversion == "1.0"
    assert env.type == "com.silmari.research.completed"
    assert env.source == "deep-research"
    assert env.subject == "exec_123"
    assert env.data["result_ref"] == "pkg://run/abc"
    assert env.data["prompt"] == "why FDO"
    assert env.data["document_id"] == "doc_9"
    assert isinstance(env.id, str) and len(env.id) > 0


def test_round_trip_json_lossless():
    env = build_envelope(
        node_id="deep-research",
        domain="research",
        correlation_id="exec_123",
        result_ref="pkg://run/abc",
        snapshot={"prompt": "why FDO", "document_id": "doc_9"},
    )
    restored = CompletionEnvelope.from_json(env.to_json())
    assert restored == env


def test_rejects_mutable_aggregate_in_data():
    class _Agg:
        pass

    try:
        build_envelope(
            node_id="x",
            domain="research",
            correlation_id="e1",
            result_ref="r",
            snapshot={"body": _Agg()},
        )
        assert False, "should have raised EnvelopeError"
    except EnvelopeError:
        pass


def test_rejects_empty_correlation_id():
    try:
        build_envelope(
            node_id="x",
            domain="research",
            correlation_id="",
            result_ref="r",
            snapshot={},
        )
        assert False, "should have raised EnvelopeError"
    except EnvelopeError:
        pass


def test_rejects_snapshot_deeper_than_max_depth():
    too_deep = {"a": {"b": {"c": 1}}}  # depth 3 > MAX_SNAPSHOT_DEPTH (2)
    try:
        build_envelope(
            node_id="x",
            domain="research",
            correlation_id="e1",
            result_ref="r",
            snapshot=too_deep,
        )
        assert False, "should have raised EnvelopeError"
    except EnvelopeError:
        pass


def test_accepts_snapshot_at_max_depth():
    at_max = {"a": {"b": 1}}  # depth 2 == MAX_SNAPSHOT_DEPTH
    env = build_envelope(
        node_id="x",
        domain="research",
        correlation_id="e1",
        result_ref="r",
        snapshot=at_max,
    )
    assert env.data["a"] == {"b": 1}


def test_rejects_snapshot_over_byte_cap():
    huge = {"blob": "x" * (9 * 1024)}  # > MAX_SNAPSHOT_BYTES (8 KiB)
    try:
        build_envelope(
            node_id="x",
            domain="research",
            correlation_id="e1",
            result_ref="r",
            snapshot=huge,
        )
        assert False, "should have raised EnvelopeError"
    except EnvelopeError:
        pass


def test_accepts_snapshot_under_byte_cap():
    small = {"note": "x" * 100}
    env = build_envelope(
        node_id="x",
        domain="test",
        correlation_id="e1",
        result_ref="r",
        snapshot=small,
    )
    assert env.data["note"] == "x" * 100


def test_list_in_snapshot_at_valid_depth():
    with_list = {"tags": ["a", "b", "c"]}  # depth 2 (dict -> list of primitives)
    env = build_envelope(
        node_id="x",
        domain="test",
        correlation_id="e1",
        result_ref="r",
        snapshot=with_list,
    )
    assert env.data["tags"] == ["a", "b", "c"]


def test_dataschema_included_when_provided():
    env = build_envelope(
        node_id="x",
        domain="research",
        correlation_id="e1",
        result_ref="r",
        snapshot={},
        dataschema="com.silmari.research.completed/v1",
    )
    assert env.dataschema == "com.silmari.research.completed/v1"
    assert "dataschema" in env.to_json()


def test_dataschema_absent_when_not_provided():
    env = build_envelope(
        node_id="x",
        domain="research",
        correlation_id="e1",
        result_ref="r",
        snapshot={},
    )
    assert env.dataschema is None
    assert "dataschema" not in env.to_json()
