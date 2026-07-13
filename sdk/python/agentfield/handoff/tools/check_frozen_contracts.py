"""CI guard: published contract schemas are frozen — never modified, only new versions added.

Usage (CI):
    git diff --name-status HEAD~1 | python -m agentfield.handoff.tools.check_frozen_contracts

Usage (programmatic):
    from agentfield.handoff.tools.check_frozen_contracts import check_frozen_contracts
    verdict, messages = check_frozen_contracts(diff_lines)
"""

from __future__ import annotations

import re
import sys

SCHEMA_PATTERN = re.compile(r"contracts/.+/v\d+\.schema\.json$")


def check_frozen_contracts(diff_lines: list[str]) -> tuple[str, list[str]]:
    messages: list[str] = []
    failed = False

    for line in diff_lines:
        line = line.strip()
        if not line:
            continue

        parts = line.split("\t", 1)
        if len(parts) != 2:
            continue

        status, filepath = parts[0].strip(), parts[1].strip()

        if not SCHEMA_PATTERN.search(filepath):
            continue

        if status == "A":
            messages.append(f"PASS: new schema added: {filepath}")
        elif status == "M":
            version_match = re.search(r"(v\d+)\.schema\.json", filepath)
            version_label = version_match.group(1) if version_match else "vN"
            messages.append(
                f"FAIL: published contract {version_label} is frozen — "
                f"add {_next_version(version_label)} instead: {filepath}"
            )
            failed = True
        elif status == "D":
            messages.append(f"FAIL: published contract schema deleted: {filepath}")
            failed = True

    verdict = "fail" if failed else "pass"
    return verdict, messages


def _next_version(version: str) -> str:
    match = re.match(r"v(\d+)", version)
    if match:
        return f"v{int(match.group(1)) + 1}"
    return "vN+1"


def main() -> None:
    diff_lines = sys.stdin.read().splitlines()
    verdict, messages = check_frozen_contracts(diff_lines)
    for msg in messages:
        print(msg)
    if verdict == "fail":
        sys.exit(1)


if __name__ == "__main__":
    main()
