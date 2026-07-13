"""Phase 2 types: HandoffDTO, CursorStore protocol, ExecutionRecord."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Protocol, runtime_checkable


@dataclass(frozen=True)
class HandoffDTO:
    """The handler callback's first argument — a validated, typed projection of
    the CloudEvent that the middleware already schema-checked."""

    event_id: str
    execution_id: str
    event_type: str
    sequence: int
    data: dict[str, Any]


@runtime_checkable
class CursorStore(Protocol):
    """Per-consumer durable cursor. Each consumer app supplies its own adapter
    backed by its own table. The middleware owns NO persistent tables."""

    def get(self, consumer: str) -> int: ...
    def advance(self, consumer: str, seq: int) -> None: ...


@dataclass(frozen=True)
class ExecutionRecord:
    """Projection of the CP's GET /executions/:id response."""

    execution_id: str
    status: str
    result: dict[str, Any] | None = None
    run_id: str | None = None
    started_at: str | None = None
    completed_at: str | None = None
