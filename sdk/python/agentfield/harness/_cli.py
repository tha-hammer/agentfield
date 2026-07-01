"""Shared async subprocess utilities for CLI-based harness providers."""

from __future__ import annotations

import asyncio
import json
import os
import re
import time
from typing import Any, Dict, List, Optional, Tuple

_ANSI_RE = re.compile(r"\x1B\[[0-?]*[ -/]*[@-~]")


def strip_ansi(text: str) -> str:
    return _ANSI_RE.sub("", text)


async def _drain(
    stream: Optional[asyncio.StreamReader],
    chunks: List[bytes],
    activity: List[float],
) -> None:
    """Read a pipe to EOF, buffering chunks and stamping last-output time."""
    if stream is None:
        return
    while True:
        chunk = await stream.read(65536)
        if not chunk:
            return
        chunks.append(chunk)
        activity[0] = time.monotonic()


async def run_cli(
    cmd: List[str],
    *,
    env: Optional[Dict[str, str]] = None,
    cwd: Optional[str] = None,
    timeout: Optional[float] = None,
    idle_timeout: Optional[float] = None,
) -> Tuple[str, str, int]:
    """Run a CLI command async. Returns (stdout, stderr, returncode).

    ``timeout`` is a total wall-clock cap. ``idle_timeout``, when set,
    additionally kills the process if it produces NO stdout/stderr output for
    that many seconds — catching silent stalls (e.g. an upstream LLM stream
    that hangs mid-response) far sooner than the total cap. A progressing CLI
    streams events continuously, so idle detection won't cut a working call as
    long as ``idle_timeout`` exceeds the longest expected gap between output
    (e.g. a long-running tool/test invocation). Both raise ``TimeoutError``.
    When ``idle_timeout`` is None the behavior is identical to the prior
    total-cap-only implementation.
    """
    merged_env = {**os.environ}
    if env:
        merged_env.update(env)

    proc = await asyncio.create_subprocess_exec(
        *cmd,
        stdin=asyncio.subprocess.DEVNULL,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
        env=merged_env,
        cwd=cwd,
    )

    if idle_timeout is None:
        # Total-cap-only path (unchanged behavior).
        try:
            stdout_bytes, stderr_bytes = await asyncio.wait_for(
                proc.communicate(), timeout=timeout
            )
        except asyncio.TimeoutError:
            proc.kill()
            await proc.wait()
            raise TimeoutError(
                f"CLI command timed out after {timeout}s: {' '.join(cmd)}"
            )
        return (
            stdout_bytes.decode("utf-8", errors="replace"),
            stderr_bytes.decode("utf-8", errors="replace"),
            proc.returncode if proc.returncode is not None else -1,
        )

    # Idle-aware path: stream both pipes concurrently, tracking the last time
    # any output arrived, and enforce both the idle cap and the total cap.
    stdout_chunks: List[bytes] = []
    stderr_chunks: List[bytes] = []
    activity = [time.monotonic()]
    start = time.monotonic()
    readers = asyncio.gather(
        _drain(proc.stdout, stdout_chunks, activity),
        _drain(proc.stderr, stderr_chunks, activity),
    )

    poll = min(max(idle_timeout / 4.0, 0.05), 15.0)
    killed_reason: Optional[str] = None
    try:
        while True:
            try:
                # shield so an idle/total-cap timeout cancels only the wait,
                # not the underlying reader tasks (which keep draining).
                await asyncio.wait_for(asyncio.shield(readers), timeout=poll)
                break  # both pipes hit EOF -> process is exiting
            except asyncio.TimeoutError:
                now = time.monotonic()
                if now - activity[0] >= idle_timeout:
                    killed_reason = f"no output for {idle_timeout}s (idle stall)"
                    break
                if timeout is not None and now - start >= timeout:
                    killed_reason = f"exceeded total timeout {timeout}s"
                    break
    finally:
        if killed_reason is not None:
            proc.kill()
        # Let the readers drain post-kill EOF, then reap the process.
        try:
            await readers
        except Exception:
            pass
        await proc.wait()

    if killed_reason is not None:
        raise TimeoutError(f"CLI command killed: {killed_reason}: {' '.join(cmd)}")

    return (
        b"".join(stdout_chunks).decode("utf-8", errors="replace"),
        b"".join(stderr_chunks).decode("utf-8", errors="replace"),
        proc.returncode if proc.returncode is not None else -1,
    )


def parse_jsonl(text: str) -> List[Dict[str, Any]]:
    """Parse JSONL (newline-delimited JSON) output. Skips invalid lines."""
    events = []
    for line in text.strip().splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            events.append(json.loads(line))
        except json.JSONDecodeError:
            continue
    return events


def extract_final_text(events: List[Dict[str, Any]]) -> Optional[str]:
    """Extract the final result text from a list of JSONL events.

    Looks for common patterns across different CLI tools:
    - type: "result" with text/result field
    - type: "item.completed" with item.text field (Codex)
    - type: "text" with part.text field (OpenCode JSON stream)
    - Last assistant message text
    """
    result_text = None
    current_text_parts: List[str] = []

    for event in events:
        event_type = event.get("type", "")

        if event_type == "step_start":
            current_text_parts = []
        elif event_type == "item.completed":
            item = event.get("item", {})
            if item.get("type") == "agent_message":
                text = item.get("text", "")
                if text:
                    result_text = text
        elif event_type == "result":
            result_text = event.get("result", event.get("text", result_text))
        elif event_type == "turn.completed":
            text = event.get("text", "")
            if text:
                result_text = text
        elif event_type in ("message", "assistant"):
            content = event.get("content", event.get("text", ""))
            if isinstance(content, str) and content:
                result_text = content
        elif event_type == "text":
            content = event.get("text", event.get("content", ""))
            part = event.get("part")
            if not content and isinstance(part, dict):
                content = part.get("text", "")
            if isinstance(content, str) and content:
                current_text_parts.append(content)
                result_text = "".join(current_text_parts)

    return result_text


def estimate_cli_cost(
    model: str,
    prompt: str,
    result_text: str | None,
) -> float | None:
    """Estimate LLM cost from prompt/completion text using litellm.

    Returns None if the model isn't in litellm's pricing DB or litellm
    is not available — callers should treat None as "unknown", not "free".
    """
    if not model:
        return None
    try:
        import litellm

        cost = litellm.completion_cost(
            model=model,
            prompt=prompt,
            completion=result_text or "",
        )
        return cost if cost and cost > 0 else None
    except Exception:
        return None
