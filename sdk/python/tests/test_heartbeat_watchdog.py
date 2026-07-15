"""Tests for the idle/heartbeat watchdog in ``_execute_async_with_callback``.

The watchdog used to cap *total* active time, which would cancel a
legitimately long but steadily-progressing reasoner the moment it crossed the
``default_execution_timeout`` budget. It is now an *idle* watchdog: it cancels
only when no progress (``Agent.heartbeat``) has been observed for the budget
window. A reasoner that never heartbeats keeps the original behavior.

These tests drive ``_execute_async_with_callback`` directly with a patched
``_post_execution_status`` so they need no control plane and run sub-second.
"""

import asyncio

import pytest

from agentfield.agent import Agent


def _make_agent(timeout: float) -> Agent:
    agent = Agent(
        node_id="hb-agent",
        agentfield_server="http://control",
        auto_register=False,
    )
    agent.async_config.default_execution_timeout = timeout
    return agent


def _patch_callback(agent: Agent, captured: list) -> None:
    async def fake_post(callback_url, payload, execution_id, max_retries=5):
        captured.append(payload)

    agent._post_execution_status = fake_post  # type: ignore[assignment]
    # ``_build_execution_callback_url`` must return a truthy URL or the method
    # returns early before running the watchdog.
    agent._build_execution_callback_url = (  # type: ignore[assignment]
        lambda execution_id: f"http://control/api/v1/executions/{execution_id}/status"
    )


@pytest.mark.asyncio
async def test_idle_watchdog_cancels_when_no_progress():
    """A reasoner that sleeps past the budget and never heartbeats is killed."""
    agent = _make_agent(timeout=1.0)
    captured: list = []
    _patch_callback(agent, captured)

    async def slow_reasoner():
        await asyncio.sleep(30.0)  # far past the 1.0s budget
        return {"done": True}

    await asyncio.wait_for(
        agent._execute_async_with_callback(
            reasoner_coro=slow_reasoner,
            execution_id="exec-idle",
            reasoner_name="slow",
        ),
        timeout=10.0,
    )

    assert captured, "expected a status callback"
    payload = captured[-1]
    assert payload["status"] == "failed"
    assert "no progress" in payload["error"]
    assert payload["error_details"]["reason"] == "reasoner_timeout"
    # Marker cleaned up.
    assert "exec-idle" not in agent._progress_markers


@pytest.mark.asyncio
async def test_heartbeat_resets_idle_deadline():
    """The core regression: heartbeating while progressing survives well past
    the budget and completes successfully."""
    agent = _make_agent(timeout=1.0)
    captured: list = []
    _patch_callback(agent, captured)

    execution_id = "exec-hb"

    async def progressing_reasoner():
        # Run for ~3x the budget, heartbeating every 0.4s (< 1.0s budget).
        for _ in range(8):
            await asyncio.sleep(0.4)
            agent.heartbeat(execution_id)
        return {"levels": 5}

    await asyncio.wait_for(
        agent._execute_async_with_callback(
            reasoner_coro=progressing_reasoner,
            execution_id=execution_id,
            reasoner_name="progressing",
        ),
        timeout=20.0,
    )

    assert captured, "expected a status callback"
    payload = captured[-1]
    assert payload["status"] == "succeeded", payload
    assert payload["result"]["levels"] == 5


@pytest.mark.asyncio
async def test_no_heartbeat_preserves_original_behavior():
    """A non-heartbeating reasoner with a small budget still times out,
    proving equivalence with the old total-active cap."""
    agent = _make_agent(timeout=1.0)
    captured: list = []
    _patch_callback(agent, captured)

    async def busy_no_hb():
        # Steady work but never heartbeats -> idle_active grows unbounded.
        for _ in range(20):
            await asyncio.sleep(0.3)
        return {"done": True}

    await asyncio.wait_for(
        agent._execute_async_with_callback(
            reasoner_coro=busy_no_hb,
            execution_id="exec-nohb",
            reasoner_name="busy",
        ),
        timeout=10.0,
    )

    payload = captured[-1]
    assert payload["status"] == "failed"
    assert "no progress" in payload["error"]


@pytest.mark.asyncio
async def test_heartbeat_unknown_execution_is_noop():
    """Heartbeating an unknown execution_id must not raise."""
    agent = _make_agent(timeout=1.0)
    # No marker registered -> silent no-op.
    agent.heartbeat("does-not-exist")
