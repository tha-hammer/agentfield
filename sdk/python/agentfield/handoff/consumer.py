"""HandoffMiddleware — durable-cursor consumer with schema validation + dedup.

subscribe() starts a background poll loop that:
  1. Reads events from the CP outbox (via HTTP) after the consumer's cursor
  2. Client-side filters by event_type (the CP read has no server-side type filter)
  3. Dedups on subject (correlation_id) in memory per R3
  4. Validates data against the frozen schema (drift → skip+log, never crash)
  5. Calls the handler with (HandoffDTO, fetch_body)
  6. Advances the cursor ONLY after the handler returns without exception (D3)

Malformed envelopes (missing CE fields) are logged, cursor advanced past them,
handler never called (§5 poison safety).
"""

from __future__ import annotations

import json
import logging
import threading
from typing import Any, Callable

from agentfield.handoff._registry import ContractRegistry
from agentfield.handoff._validate import ValidationError, validate
from agentfield.handoff.control_plane_source import ControlPlaneSource
from agentfield.handoff.types import CursorStore, ExecutionRecord, HandoffDTO
from agentfield.provenance_handoff import MalformedEnvelopeError

logger = logging.getLogger(__name__)

HandlerFn = Callable[[HandoffDTO, Callable[[str], ExecutionRecord]], None]

_REQUIRED_CE_FIELDS = {"id", "type", "subject"}


class ConsumerHandle:
    """Handle to a running background consumer thread."""

    def __init__(self, stop_event: threading.Event, thread: threading.Thread) -> None:
        self._stop_event = stop_event
        self._thread = thread

    def stop(self, timeout: float = 5.0) -> None:
        self._stop_event.set()
        self._thread.join(timeout=timeout)

    def is_alive(self) -> bool:
        return self._thread.is_alive()


class HandoffMiddleware:
    """Phase 2 middleware: subscribe to handoff events through the control plane."""

    def __init__(
        self,
        cp_base_url: str,
        cp_api_key: str,
        cursor_store: CursorStore,
        registry: ContractRegistry,
    ) -> None:
        self._cp = ControlPlaneSource(cp_base_url, cp_api_key)
        self._cursor_store = cursor_store
        self._registry = registry

    def subscribe(
        self,
        event_type: str,
        handler: HandlerFn,
        *,
        consumer_name: str = "default",
        batch_limit: int = 50,
        poll_interval: float = 5.0,
        reader: Any | None = None,
    ) -> ConsumerHandle:
        if event_type not in self._registry:
            raise ValueError(
                f"cannot subscribe to {event_type!r}: not in the contract registry"
            )

        stop_event = threading.Event()
        thread = threading.Thread(
            target=self._poll_loop,
            args=(event_type, handler, consumer_name, batch_limit, poll_interval, stop_event, reader),
            name=f"handoff-{consumer_name}",
            daemon=True,
        )
        thread.start()
        return ConsumerHandle(stop_event, thread)

    def _poll_loop(
        self,
        event_type: str,
        handler: HandlerFn,
        consumer_name: str,
        batch_limit: int,
        poll_interval: float,
        stop_event: threading.Event,
        reader: Any | None,
    ) -> None:
        seen: set[str] = set()
        while not stop_event.is_set():
            try:
                cursor = self._cursor_store.get(consumer_name)
                records = self._read_events(cursor, batch_limit, reader)
                for record in records:
                    if stop_event.is_set():
                        break
                    self._process_record(
                        record, event_type, handler, consumer_name, seen,
                    )
            except Exception:
                logger.exception("handoff poll error")
            if stop_event.is_set():
                break
            stop_event.wait(poll_interval)

    def _read_events(self, cursor: int, limit: int, reader: Any | None) -> list[dict]:
        if reader is not None:
            return reader.read_since(cursor, limit)
        import requests as _requests

        resp = _requests.get(
            f"{self._cp._base_url}/api/v1/events",
            headers={"X-API-Key": self._cp._api_key},
            params={"since": cursor, "limit": limit},
            timeout=30,
        )
        resp.raise_for_status()
        return resp.json()

    def _process_record(
        self,
        record: dict[str, Any],
        event_type: str,
        handler: HandlerFn,
        consumer_name: str,
        seen: set[str],
    ) -> None:
        seq = record.get("sequence")
        if seq is None:
            logger.warning("handoff: record missing sequence, skipping")
            return

        missing = _REQUIRED_CE_FIELDS - set(record)
        if missing:
            logger.warning("handoff: malformed envelope missing %s at seq=%s", missing, seq)
            self._cursor_store.advance(consumer_name, seq)
            return

        record_type = record.get("type", "")
        if record_type != event_type:
            self._cursor_store.advance(consumer_name, seq)
            return

        subject = record["subject"]
        if subject in seen:
            self._cursor_store.advance(consumer_name, seq)
            return

        data = record.get("data")
        if isinstance(data, str):
            try:
                data = json.loads(data)
            except (json.JSONDecodeError, TypeError):
                logger.warning("handoff: unparseable data at seq=%s", seq)
                self._cursor_store.advance(consumer_name, seq)
                return

        try:
            validate({"type": event_type, "data": data}, registry=self._registry)
        except ValidationError as exc:
            logger.warning("handoff: schema drift at seq=%s: %s", seq, exc)
            self._cursor_store.advance(consumer_name, seq)
            return

        dto = HandoffDTO(
            event_id=record["id"],
            execution_id=subject,
            event_type=event_type,
            sequence=seq,
            data=data,
        )

        def _fetch(eid: str) -> ExecutionRecord:
            return self._cp.fetch_body(eid)

        try:
            handler(dto, _fetch)
        except Exception:
            logger.exception("handoff: handler failed for seq=%s, cursor NOT advanced", seq)
            return

        seen.add(subject)
        self._cursor_store.advance(consumer_name, seq)
