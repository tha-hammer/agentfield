"""announce — validated producer-side append to the outbox.

The Python announce is the reference implementation / test surface; the
Go control-plane's BuildResearchCompletedOutboxRecord is the production
producer.  A contract test asserts announce's output matches the frozen
golden fixture.
"""

from __future__ import annotations

from typing import Any, Mapping

from agentfield.handoff._registry import ContractRegistry
from agentfield.handoff._validate import ValidationError, validate
from agentfield.provenance_handoff import build_envelope


def announce(
    event_type: str,
    dto: Mapping[str, Any],
    *,
    execution_id: str,
    node_id: str,
    publisher: Any,
    registry: ContractRegistry | None = None,
) -> str:
    """Validate DTO against the frozen schema and append to the outbox.

    Returns the CloudEvents id.  Raises ValidationError on schema mismatch
    (reject on produce — never emit an invalid event).
    """
    event_for_validation = {"type": event_type, "data": dict(dto)}
    validate(event_for_validation, registry=registry)

    domain = _extract_domain(event_type)
    result_ref = dto.get("result_ref", "")

    snapshot = {k: v for k, v in dto.items() if k != "result_ref"}

    env = build_envelope(
        node_id=node_id,
        domain=domain,
        correlation_id=execution_id,
        result_ref=result_ref,
        snapshot=snapshot,
    )

    publisher.publish({
        "type": env.type,
        "execution_id": env.subject,
        "agent_node_id": node_id,
        "data": env.to_json(),
    })
    return env.id


def _extract_domain(event_type: str) -> str:
    """Extract the domain from a versioned event type string.

    'com.silmari.research.completed.v1' → 'research'
    """
    parts = event_type.split(".")
    if len(parts) >= 4 and parts[0] == "com" and parts[1] == "silmari":
        return parts[2]
    return event_type
