from __future__ import annotations

import json
from dataclasses import dataclass
from pathlib import Path
from typing import Any


@dataclass(frozen=True)
class ContractEntry:
    type_version: str
    schema: dict[str, Any]
    schema_path: str


class ContractRegistry:

    def __init__(self) -> None:
        self._entries: dict[str, ContractEntry] = {}

    def register(self, type_version: str, schema: dict[str, Any], schema_path: str) -> None:
        self._entries[type_version] = ContractEntry(
            type_version=type_version,
            schema=schema,
            schema_path=schema_path,
        )

    def get(self, type_version: str) -> ContractEntry:
        try:
            return self._entries[type_version]
        except KeyError:
            raise KeyError(f"no contract registered for {type_version!r}") from None

    def list(self) -> list[ContractEntry]:
        return list(self._entries.values())

    def __contains__(self, type_version: str) -> bool:
        return type_version in self._entries

    def __len__(self) -> int:
        return len(self._entries)


def _auto_register(reg: ContractRegistry) -> None:
    contracts_dir = Path(__file__).parent / "contracts"
    if not contracts_dir.is_dir():
        return
    for type_dir in sorted(contracts_dir.iterdir()):
        if not type_dir.is_dir():
            continue
        base_type = type_dir.name
        for schema_file in sorted(type_dir.glob("v*.schema.json")):
            stem = schema_file.stem.replace(".schema", "")
            type_version = f"{base_type}.{stem}"
            with open(schema_file) as f:
                schema = json.load(f)
            reg.register(type_version, schema, str(schema_file))


registry = ContractRegistry()
_auto_register(registry)
