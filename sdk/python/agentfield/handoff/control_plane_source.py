"""fetch_body — reads ONLY the control plane (GET /executions/:id).

Source-of-truth invariant (C3-Notification): there is NO code path that reads
a peer app's storage.  A test asserts this by grepping the module source.
"""

from __future__ import annotations

from typing import Any

import requests

from agentfield.handoff.types import ExecutionRecord


class ControlPlaneSource:
    """Thin HTTP client for the CP execution read surface."""

    def __init__(self, base_url: str, api_key: str) -> None:
        self._base_url = base_url.rstrip("/")
        self._api_key = api_key

    def fetch_body(self, execution_id: str) -> ExecutionRecord:
        resp = requests.get(
            f"{self._base_url}/api/v1/executions/{execution_id}",
            headers={"X-API-Key": self._api_key},
            timeout=30,
        )
        resp.raise_for_status()
        body: dict[str, Any] = resp.json()
        return ExecutionRecord(
            execution_id=body.get("execution_id", execution_id),
            status=body.get("status", ""),
            result=body.get("result"),
            run_id=body.get("run_id"),
            started_at=body.get("started_at"),
            completed_at=body.get("completed_at"),
        )
