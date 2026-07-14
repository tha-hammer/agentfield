"""MW Phase 2 — HandoffMiddleware (B1–B5).

Tests the three-op surface: subscribe, fetch_body, announce.
TDD: each behavior has a red-green-refactor cycle.
"""

from __future__ import annotations

import json
import time
from pathlib import Path
from typing import Any
from unittest.mock import MagicMock

import pytest

from agentfield.handoff import (
    ContractRegistry,
    HandoffDTO,
    HandoffMiddleware,
    ValidationError,
    announce,
    registry,
)
from agentfield.handoff.consumer import ConsumerHandle


# ─────────────── test helpers ───────────────


class FakeCursorStore:
    """In-memory cursor store for tests."""

    def __init__(self, initial: int = 0) -> None:
        self._cursors: dict[str, int] = {}
        self._initial = initial

    def get(self, consumer: str) -> int:
        return self._cursors.get(consumer, self._initial)

    def advance(self, consumer: str, seq: int) -> None:
        current = self._cursors.get(consumer, 0)
        self._cursors[consumer] = max(current, seq)


class FakeReader:
    """Injectable event reader that returns pre-loaded records."""

    def __init__(self, records: list[dict[str, Any]]) -> None:
        self._records = records
        self._delivered = False

    def read_since(self, cursor: int, limit: int) -> list[dict[str, Any]]:
        if self._delivered:
            return []
        self._delivered = True
        return [r for r in self._records if r.get("sequence", 0) > cursor]


GOLDEN_DATA = {
    "run_id": "3f6d1a90-0000-4000-8000-000000000abc",
    "status": "succeeded",
    "title": "How short-form reels convert viewers",
    "result_ref": "cp-execution://exec_abc123/result",
    "research_prompt": "How do short-form reels convert viewers into subscribers?",
    "research_document_id": "exec_abc123",
}

EVENT_TYPE = "com.silmari.research.completed.v1"


def _make_record(
    seq: int = 1,
    event_id: str = "ce-001",
    subject: str = "exec_abc123",
    event_type: str = EVENT_TYPE,
    data: dict | None = None,
) -> dict[str, Any]:
    return {
        "sequence": seq,
        "id": event_id,
        "type": event_type,
        "subject": subject,
        "data": data or GOLDEN_DATA,
    }


def _subscribe_and_wait(
    mw: HandoffMiddleware,
    handler,
    reader: FakeReader,
    *,
    consumer_name: str = "test",
    wait: float = 0.5,
) -> ConsumerHandle:
    handle = mw.subscribe(
        EVENT_TYPE,
        handler,
        consumer_name=consumer_name,
        poll_interval=0.05,
        reader=reader,
    )
    time.sleep(wait)
    handle.stop()
    return handle


# ════════════════════════════════════════════════════════════════
# B1a — HandoffDTO
# ════════════════════════════════════════════════════════════════


class TestB1a_HandoffDTO:
    def test_dto_has_required_fields(self):
        dto = HandoffDTO(
            event_id="ce-001",
            execution_id="exec-001",
            event_type=EVENT_TYPE,
            sequence=1,
            data=GOLDEN_DATA,
        )
        assert dto.event_id == "ce-001"
        assert dto.execution_id == "exec-001"
        assert dto.event_type == EVENT_TYPE
        assert dto.sequence == 1
        assert dto.data == GOLDEN_DATA

    def test_dto_is_frozen(self):
        dto = HandoffDTO(
            event_id="ce-001",
            execution_id="exec-001",
            event_type=EVENT_TYPE,
            sequence=1,
            data=GOLDEN_DATA,
        )
        with pytest.raises(AttributeError):
            dto.event_id = "other"  # type: ignore[misc]


# ════════════════════════════════════════════════════════════════
# B1b — HandoffMiddleware construction
# ════════════════════════════════════════════════════════════════


class TestB1b_Construction:
    def test_construct_with_required_args(self):
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=FakeCursorStore(),
            registry=registry,
        )
        assert mw is not None

    def test_subscribe_rejects_unregistered_type(self):
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=FakeCursorStore(),
            registry=registry,
        )
        with pytest.raises(ValueError, match="not in the contract registry"):
            mw.subscribe("com.silmari.nonexistent.v99", lambda dto, fb: None)


# ════════════════════════════════════════════════════════════════
# B1 — subscribe: durable, in-memory-dedup'd, validated
# ════════════════════════════════════════════════════════════════


class TestB1_Subscribe:
    def test_handler_receives_valid_dto(self):
        received = []

        def handler(dto: HandoffDTO, fetch_body):
            received.append(dto)

        store = FakeCursorStore()
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=store,
            registry=registry,
        )
        reader = FakeReader([_make_record(seq=1)])
        _subscribe_and_wait(mw, handler, reader)

        assert len(received) == 1
        dto = received[0]
        assert dto.event_id == "ce-001"
        assert dto.execution_id == "exec_abc123"
        assert dto.event_type == EVENT_TYPE
        assert dto.sequence == 1
        assert dto.data == GOLDEN_DATA

    def test_cursor_advances_after_handler_success(self):
        store = FakeCursorStore()
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=store,
            registry=registry,
        )
        reader = FakeReader([_make_record(seq=5)])
        _subscribe_and_wait(mw, lambda dto, fb: None, reader)

        assert store.get("test") == 5

    def test_cursor_not_advanced_when_handler_raises(self):
        def bad_handler(dto, fb):
            raise RuntimeError("handler failure")

        store = FakeCursorStore()
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=store,
            registry=registry,
        )
        reader = FakeReader([_make_record(seq=3)])
        _subscribe_and_wait(mw, bad_handler, reader)

        assert store.get("test") == 0

    def test_dedup_on_subject(self):
        received = []

        def handler(dto: HandoffDTO, fetch_body):
            received.append(dto)

        store = FakeCursorStore()
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=store,
            registry=registry,
        )
        records = [
            _make_record(seq=1, event_id="ce-001", subject="exec-same"),
            _make_record(seq=2, event_id="ce-002", subject="exec-same"),
        ]
        reader = FakeReader(records)
        _subscribe_and_wait(mw, handler, reader)

        assert len(received) == 1
        assert store.get("test") == 2

    def test_filters_by_event_type(self):
        received = []

        def handler(dto: HandoffDTO, fetch_body):
            received.append(dto)

        store = FakeCursorStore()
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=store,
            registry=registry,
        )
        records = [
            _make_record(seq=1, event_type="com.silmari.other.completed.v1"),
            _make_record(seq=2, event_type=EVENT_TYPE),
        ]
        reader = FakeReader(records)
        _subscribe_and_wait(mw, handler, reader)

        assert len(received) == 1
        assert received[0].sequence == 2
        assert store.get("test") == 2

    def test_schema_drift_skips_event(self):
        received = []

        def handler(dto: HandoffDTO, fetch_body):
            received.append(dto)

        store = FakeCursorStore()
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=store,
            registry=registry,
        )
        bad_data = {**GOLDEN_DATA, "extra_field": "should_fail"}
        records = [_make_record(seq=1, data=bad_data)]
        reader = FakeReader(records)
        _subscribe_and_wait(mw, handler, reader)

        assert len(received) == 0
        assert store.get("test") == 1

    def test_data_as_json_string_is_parsed(self):
        received = []

        def handler(dto: HandoffDTO, fetch_body):
            received.append(dto)

        store = FakeCursorStore()
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=store,
            registry=registry,
        )
        records = [_make_record(seq=1, data=json.dumps(GOLDEN_DATA))]
        reader = FakeReader(records)
        _subscribe_and_wait(mw, handler, reader)

        assert len(received) == 1
        assert received[0].data == GOLDEN_DATA


# ════════════════════════════════════════════════════════════════
# B1c — Malformed envelope handling
# ════════════════════════════════════════════════════════════════


class TestB1c_MalformedEnvelope:
    def test_missing_id_skips_and_advances(self):
        received = []
        store = FakeCursorStore()
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=store,
            registry=registry,
        )
        record = {"sequence": 1, "type": EVENT_TYPE, "subject": "exec-001", "data": GOLDEN_DATA}
        reader = FakeReader([record])
        _subscribe_and_wait(mw, lambda dto, fb: received.append(dto), reader)

        assert len(received) == 0
        assert store.get("test") == 1

    def test_missing_subject_skips_and_advances(self):
        received = []
        store = FakeCursorStore()
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=store,
            registry=registry,
        )
        record = {"sequence": 1, "id": "ce-001", "type": EVENT_TYPE, "data": GOLDEN_DATA}
        reader = FakeReader([record])
        _subscribe_and_wait(mw, lambda dto, fb: received.append(dto), reader)

        assert len(received) == 0
        assert store.get("test") == 1

    def test_missing_sequence_skips(self):
        received = []
        store = FakeCursorStore()
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=store,
            registry=registry,
        )
        record = {"id": "ce-001", "type": EVENT_TYPE, "subject": "exec-001", "data": GOLDEN_DATA}
        reader = FakeReader([record])
        _subscribe_and_wait(mw, lambda dto, fb: received.append(dto), reader)

        assert len(received) == 0


# ════════════════════════════════════════════════════════════════
# B2 — fetch_body reads ONLY the control plane
# ════════════════════════════════════════════════════════════════


class TestB2_FetchBody:
    def test_fetch_body_closure_is_provided_to_handler(self):
        fetch_body_calls = []

        def handler(dto: HandoffDTO, fetch_body):
            fetch_body_calls.append(fetch_body)

        store = FakeCursorStore()
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=store,
            registry=registry,
        )
        reader = FakeReader([_make_record(seq=1)])
        _subscribe_and_wait(mw, handler, reader)

        assert len(fetch_body_calls) == 1
        assert callable(fetch_body_calls[0])

    def test_no_peer_db_read_path_in_source(self):
        """Source-of-truth invariant: control_plane_source.py has no peer-DB read."""
        source = Path(__file__).resolve().parent.parent / "agentfield" / "handoff" / "control_plane_source.py"
        code = source.read_text()
        assert "psycopg" not in code
        assert "psycopg2" not in code
        assert "sqlalchemy" not in code
        assert "deepresearch." not in code
        assert "reel_job" not in code
        assert "carousel" not in code


# ════════════════════════════════════════════════════════════════
# B4 — announce: validated producer append
# ════════════════════════════════════════════════════════════════


class TestB4_Announce:
    def test_announce_validates_and_publishes(self):
        publisher = MagicMock()
        ce_id = announce(
            EVENT_TYPE,
            GOLDEN_DATA,
            execution_id="exec_abc123",
            node_id="silmari-af-deep-research",
            publisher=publisher,
            registry=registry,
        )
        assert isinstance(ce_id, str)
        assert len(ce_id) > 0
        publisher.publish.assert_called_once()
        payload = publisher.publish.call_args[0][0]
        assert payload["type"] == "com.silmari.research.completed"
        assert payload["execution_id"] == "exec_abc123"

    def test_announce_rejects_invalid_dto(self):
        publisher = MagicMock()
        bad_dto = {**GOLDEN_DATA, "extra_field": "nope"}
        with pytest.raises(ValidationError):
            announce(
                EVENT_TYPE,
                bad_dto,
                execution_id="exec_abc123",
                node_id="silmari-af-deep-research",
                publisher=publisher,
                registry=registry,
            )
        publisher.publish.assert_not_called()

    def test_announce_rejects_missing_required_field(self):
        publisher = MagicMock()
        bad_dto = {k: v for k, v in GOLDEN_DATA.items() if k != "run_id"}
        with pytest.raises(ValidationError):
            announce(
                EVENT_TYPE,
                bad_dto,
                execution_id="exec_abc123",
                node_id="silmari-af-deep-research",
                publisher=publisher,
                registry=registry,
            )


# ════════════════════════════════════════════════════════════════
# B5 — Enable in prod (fail-closed checks)
# ════════════════════════════════════════════════════════════════


class TestB5_FailClosed:
    def test_subscribe_fails_with_empty_registry(self):
        empty_reg = ContractRegistry()
        mw = HandoffMiddleware(
            cp_base_url="http://localhost:8080",
            cp_api_key="test-key",
            cursor_store=FakeCursorStore(),
            registry=empty_reg,
        )
        with pytest.raises(ValueError, match="not in the contract registry"):
            mw.subscribe(EVENT_TYPE, lambda dto, fb: None)


# ════════════════════════════════════════════════════════════════
# Import surface — Phase 2
# ════════════════════════════════════════════════════════════════


class TestPhase2ImportSurface:
    def test_handoff_middleware_importable(self):
        from agentfield.handoff import HandoffMiddleware
        assert HandoffMiddleware is not None

    def test_handoff_dto_importable(self):
        from agentfield.handoff import HandoffDTO
        assert HandoffDTO is not None

    def test_cursor_store_importable(self):
        from agentfield.handoff import CursorStore
        assert CursorStore is not None

    def test_announce_importable(self):
        from agentfield.handoff import announce
        assert callable(announce)

    def test_consumer_handle_importable(self):
        from agentfield.handoff import ConsumerHandle
        assert ConsumerHandle is not None

    def test_execution_record_importable(self):
        from agentfield.handoff import ExecutionRecord
        assert ExecutionRecord is not None
