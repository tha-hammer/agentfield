"""MW Phase 1 — Contract Registry + Freeze Governance (B1–B4).

Governing spec: specs/cross-app-handoff.pattern.md (6 contracts).
Decisions: D1 (shared package), D2 (version in type string), D3 (validate data only),
D4 (nullable title/research_prompt), D5 (jsonschema dep), D6 (Go producer alignment),
D7 (directory convention).
"""

import json
from dataclasses import is_dataclass
from pathlib import Path

import pytest

# ──────────────────────── Golden fixture ────────────────────────

GOLDEN_FIXTURE = {
    "id": "ce-0f9b2c7a-4d1e-4a2b-9c3d-000000000001",
    "source": "silmari-af-deep-research",
    "type": "com.silmari.research.completed.v1",
    "subject": "exec_abc123",
    "time": "2026-07-12T18:00:00Z",
    "data": {
        "run_id": "3f6d1a90-0000-4000-8000-000000000abc",
        "status": "succeeded",
        "title": "How short-form reels convert viewers",
        "result_ref": "cp-execution://exec_abc123/result",
        "research_prompt": "How do short-form reels convert viewers into subscribers?",
        "research_document_id": "exec_abc123",
    },
}

GOLDEN_DTO = GOLDEN_FIXTURE["data"]


# ════════════════════════════════════════════════════════════════
# B1 — Registry declares each contract as a versioned, schema'd type
# ════════════════════════════════════════════════════════════════


class TestB1_RegistryLookup:
    def test_get_returns_entry_for_research_completed_v1(self):
        from agentfield.handoff import registry

        entry = registry.get("com.silmari.research.completed.v1")
        assert entry.type_version == "com.silmari.research.completed.v1"
        assert isinstance(entry.schema, dict)
        assert entry.schema.get("type") == "object"
        assert "run_id" in entry.schema.get("properties", {})

    def test_get_raises_keyerror_for_unknown_type(self):
        from agentfield.handoff import registry

        with pytest.raises(KeyError, match="no contract registered"):
            registry.get("com.silmari.nonexistent.v1")

    def test_schema_has_all_dto_fields(self):
        from agentfield.handoff import registry

        entry = registry.get("com.silmari.research.completed.v1")
        props = set(entry.schema["properties"].keys())
        expected = {"run_id", "status", "title", "result_ref", "research_prompt", "research_document_id"}
        assert props == expected

    def test_schema_forbids_additional_properties(self):
        from agentfield.handoff import registry

        entry = registry.get("com.silmari.research.completed.v1")
        assert entry.schema.get("additionalProperties") is False

    def test_list_returns_all_registered_contracts(self):
        from agentfield.handoff import registry

        entries = registry.list()
        assert len(entries) >= 1
        type_versions = {e.type_version for e in entries}
        assert "com.silmari.research.completed.v1" in type_versions

    def test_contains_check(self):
        from agentfield.handoff import registry

        assert "com.silmari.research.completed.v1" in registry
        assert "com.silmari.nonexistent.v99" not in registry

    def test_schema_path_points_to_real_file(self):
        from agentfield.handoff import registry

        entry = registry.get("com.silmari.research.completed.v1")
        assert Path(entry.schema_path).exists()

    def test_schema_declares_nullable_title(self):
        from agentfield.handoff import registry

        entry = registry.get("com.silmari.research.completed.v1")
        title_type = entry.schema["properties"]["title"]["type"]
        assert title_type == ["string", "null"]

    def test_schema_declares_nullable_research_prompt(self):
        from agentfield.handoff import registry

        entry = registry.get("com.silmari.research.completed.v1")
        rp_type = entry.schema["properties"]["research_prompt"]["type"]
        assert rp_type == ["string", "null"]


# ════════════════════════════════════════════════════════════════
# B2 — Boundary validation: fail-closed
# ════════════════════════════════════════════════════════════════


class TestB2_Validation:
    def test_conformant_event_passes(self):
        from agentfield.handoff import validate

        validate(GOLDEN_FIXTURE)

    def test_missing_data_field_raises(self):
        from agentfield.handoff import ValidationError, validate

        event = {"type": "com.silmari.research.completed.v1"}
        with pytest.raises(ValidationError, match="no 'data' field"):
            validate(event)

    def test_missing_type_field_raises(self):
        from agentfield.handoff import ValidationError, validate

        with pytest.raises(ValidationError, match="no 'type' field"):
            validate({"data": GOLDEN_DTO})

    def test_unregistered_type_raises(self):
        from agentfield.handoff import ValidationError, validate

        event = {"type": "com.silmari.unknown.v1", "data": GOLDEN_DTO}
        with pytest.raises(ValidationError, match="no contract registered"):
            validate(event)

    def test_missing_required_field_raises(self):
        from agentfield.handoff import ValidationError, validate

        bad_dto = {k: v for k, v in GOLDEN_DTO.items() if k != "run_id"}
        event = {**GOLDEN_FIXTURE, "data": bad_dto}
        with pytest.raises(ValidationError):
            validate(event)

    def test_extra_field_raises(self):
        from agentfield.handoff import ValidationError, validate

        bad_dto = {**GOLDEN_DTO, "body": "should not be here"}
        event = {**GOLDEN_FIXTURE, "data": bad_dto}
        with pytest.raises(ValidationError):
            validate(event)

    def test_wrong_type_field_raises(self):
        from agentfield.handoff import ValidationError, validate

        bad_dto = {**GOLDEN_DTO, "run_id": 12345}
        event = {**GOLDEN_FIXTURE, "data": bad_dto}
        with pytest.raises(ValidationError):
            validate(event)

    def test_invalid_status_enum_raises(self):
        from agentfield.handoff import ValidationError, validate

        bad_dto = {**GOLDEN_DTO, "status": "running"}
        event = {**GOLDEN_FIXTURE, "data": bad_dto}
        with pytest.raises(ValidationError):
            validate(event)

    def test_null_title_passes(self):
        from agentfield.handoff import validate

        dto = {**GOLDEN_DTO, "title": None}
        event = {**GOLDEN_FIXTURE, "data": dto}
        validate(event)

    def test_null_research_prompt_passes(self):
        from agentfield.handoff import validate

        dto = {**GOLDEN_DTO, "research_prompt": None}
        event = {**GOLDEN_FIXTURE, "data": dto}
        validate(event)

    def test_null_run_id_raises(self):
        from agentfield.handoff import ValidationError, validate

        bad_dto = {**GOLDEN_DTO, "run_id": None}
        event = {**GOLDEN_FIXTURE, "data": bad_dto}
        with pytest.raises(ValidationError):
            validate(event)


# ════════════════════════════════════════════════════════════════
# B3 — Golden fixture drift check
# ════════════════════════════════════════════════════════════════


PRODUCER_FIXTURE = Path("/home/maceo/ntm_Dev/silmari-agentfield-system/silmari-af-deep-research/tests/producer/fixtures/research_completed.cloudevent.json")
CONSUMER_FIXTURE = Path("/home/maceo/ntm_Dev/carousel-impl/tests/web/fixtures/research_completed.cloudevent.json")

# Cross-repo consistency check: requires the deep-research and reels-af
# sibling repos checked out side-by-side (the local multi-repo dev layout).
# CI checks out only this repo, so the sibling paths never exist there —
# skip rather than fail, same idiom as tests/integration/conftest.py's
# "not available in this checkout" guard.
_siblings_available = PRODUCER_FIXTURE.parent.exists() and CONSUMER_FIXTURE.parent.exists()


@pytest.mark.skipif(
    not _siblings_available,
    reason="deep-research/reels-af sibling repos not available in this checkout",
)
class TestB3_GoldenFixtureDrift:
    def test_producer_fixture_exists(self):
        assert PRODUCER_FIXTURE.exists(), f"producer fixture missing: {PRODUCER_FIXTURE}"

    def test_consumer_fixture_exists(self):
        assert CONSUMER_FIXTURE.exists(), f"consumer fixture missing: {CONSUMER_FIXTURE}"

    def test_fixtures_are_byte_identical(self):
        assert PRODUCER_FIXTURE.read_bytes() == CONSUMER_FIXTURE.read_bytes()

    def test_producer_fixture_data_validates_against_schema(self):
        from agentfield.handoff import validate

        fixture = json.loads(PRODUCER_FIXTURE.read_text())
        event = {**fixture, "type": "com.silmari.research.completed.v1"}
        validate(event)

    def test_consumer_fixture_data_validates_against_schema(self):
        from agentfield.handoff import validate

        fixture = json.loads(CONSUMER_FIXTURE.read_text())
        event = {**fixture, "type": "com.silmari.research.completed.v1"}
        validate(event)

    def test_fixture_dto_keys_match_schema_properties(self):
        from agentfield.handoff import registry

        entry = registry.get("com.silmari.research.completed.v1")
        schema_keys = set(entry.schema["properties"].keys())

        fixture = json.loads(PRODUCER_FIXTURE.read_text())
        fixture_keys = set(fixture["data"].keys())

        assert fixture_keys == schema_keys


# ════════════════════════════════════════════════════════════════
# B4 — Published-schema immutability CI guard (the freeze)
# ════════════════════════════════════════════════════════════════


class TestB4_FreezeGuard:
    def test_modified_schema_fails(self):
        from agentfield.handoff.tools.check_frozen_contracts import check_frozen_contracts

        diff = ["M\tcontracts/com.silmari.research.completed/v1.schema.json"]
        verdict, messages = check_frozen_contracts(diff)
        assert verdict == "fail"
        assert any("frozen" in m for m in messages)

    def test_added_schema_passes(self):
        from agentfield.handoff.tools.check_frozen_contracts import check_frozen_contracts

        diff = ["A\tcontracts/com.silmari.research.completed/v2.schema.json"]
        verdict, messages = check_frozen_contracts(diff)
        assert verdict == "pass"

    def test_deleted_schema_fails(self):
        from agentfield.handoff.tools.check_frozen_contracts import check_frozen_contracts

        diff = ["D\tcontracts/com.silmari.research.completed/v1.schema.json"]
        verdict, messages = check_frozen_contracts(diff)
        assert verdict == "fail"

    def test_non_schema_file_ignored(self):
        from agentfield.handoff.tools.check_frozen_contracts import check_frozen_contracts

        diff = ["M\tsome/other/file.py"]
        verdict, messages = check_frozen_contracts(diff)
        assert verdict == "pass"
        assert len(messages) == 0

    def test_mixed_diff_fails_on_any_modification(self):
        from agentfield.handoff.tools.check_frozen_contracts import check_frozen_contracts

        diff = [
            "A\tcontracts/com.silmari.research.completed/v2.schema.json",
            "M\tcontracts/com.silmari.research.completed/v1.schema.json",
        ]
        verdict, messages = check_frozen_contracts(diff)
        assert verdict == "fail"

    def test_message_suggests_next_version(self):
        from agentfield.handoff.tools.check_frozen_contracts import check_frozen_contracts

        diff = ["M\tcontracts/com.silmari.research.completed/v1.schema.json"]
        _, messages = check_frozen_contracts(diff)
        assert any("v2" in m for m in messages)

    def test_empty_diff_passes(self):
        from agentfield.handoff.tools.check_frozen_contracts import check_frozen_contracts

        verdict, messages = check_frozen_contracts([])
        assert verdict == "pass"

    def test_deep_path_still_detected(self):
        from agentfield.handoff.tools.check_frozen_contracts import check_frozen_contracts

        diff = ["M\tagentfield/handoff/contracts/com.silmari.reel.completed/v1.schema.json"]
        verdict, _ = check_frozen_contracts(diff)
        assert verdict == "fail"


# ════════════════════════════════════════════════════════════════
# Import surface
# ════════════════════════════════════════════════════════════════


class TestImportSurface:
    def test_registry_importable(self):
        from agentfield.handoff import ContractEntry, ContractRegistry, registry
        assert isinstance(registry, ContractRegistry)
        assert is_dataclass(ContractEntry)

    def test_validate_importable(self):
        from agentfield.handoff import ValidationError, validate
        assert callable(validate)
        assert issubclass(ValidationError, Exception)
