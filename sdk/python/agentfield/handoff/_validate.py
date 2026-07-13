from __future__ import annotations

from typing import Any

import jsonschema

from agentfield.handoff._registry import ContractRegistry


class ValidationError(Exception):
    pass


def validate(event: dict[str, Any], *, registry: ContractRegistry | None = None) -> None:
    if registry is None:
        from agentfield.handoff._registry import registry as _default
        registry = _default

    event_type = event.get("type")
    if not event_type:
        raise ValidationError("event has no 'type' field")

    if event_type not in registry:
        raise ValidationError(f"no contract registered for type {event_type!r}")

    entry = registry.get(event_type)

    data = event.get("data")
    if data is None:
        raise ValidationError("event has no 'data' field")

    try:
        jsonschema.validate(instance=data, schema=entry.schema)
    except jsonschema.ValidationError as exc:
        raise ValidationError(str(exc.message)) from None
