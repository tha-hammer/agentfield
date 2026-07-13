"""agentfield.handoff — Cross-App Handoff Middleware.

Phase 1: Contract Registry + Freeze Governance.
Phase 2: HandoffMiddleware (subscribe/announce/fetch_body) + HandoffDTO.

Provides:
  - ContractRegistry: versioned, frozen JSON-Schema contracts for handoff events
  - validate(): fail-closed boundary validation of CloudEvent data against frozen schemas
  - HandoffMiddleware: durable-cursor consumer with schema validation + dedup
  - HandoffDTO: typed projection of a validated CloudEvent for handler callbacks
  - CursorStore: protocol for per-consumer durable cursors (app provides adapter)
  - announce(): validated producer-side append to the outbox
"""

from agentfield.handoff._registry import ContractEntry, ContractRegistry, registry
from agentfield.handoff._validate import ValidationError, validate
from agentfield.handoff.consumer import ConsumerHandle, HandoffMiddleware
from agentfield.handoff.producer import announce
from agentfield.handoff.types import CursorStore, ExecutionRecord, HandoffDTO

__all__ = [
    "ConsumerHandle",
    "ContractEntry",
    "ContractRegistry",
    "CursorStore",
    "ExecutionRecord",
    "HandoffDTO",
    "HandoffMiddleware",
    "ValidationError",
    "announce",
    "registry",
    "validate",
]
