"""MW Phase 3 B4 (Part 1) — reel.completed contract registration + freeze.

Mirrors the research.completed registry/validation/freeze tests
(test_handoff_registry.py) for the newly frozen contract
``com.silmari.reel.completed.v1``. Registered purely by presence of
``contracts/com.silmari.reel.completed/v1.schema.json`` (auto-register in
_registry.py), exactly like research.completed.

DTO: {run_id, status, reel_ref, source_execution_id, duration_s, beat_count}.
"""

import pytest

REEL_TYPE = "com.silmari.reel.completed.v1"

GOLDEN_REEL_DTO = {
    "run_id": "3f6d1a90-0000-4000-8000-000000000abc",
    "status": "succeeded",
    "reel_ref": "cp-execution://exec_abc123/result",
    "source_execution_id": "exec_abc123",
    "duration_s": 12.5,
    "beat_count": 5,
}

GOLDEN_REEL_EVENT = {"type": REEL_TYPE, "data": GOLDEN_REEL_DTO}

_EXPECTED_KEYS = {
    "run_id",
    "status",
    "reel_ref",
    "source_execution_id",
    "duration_s",
    "beat_count",
}


# ─────────────────────────── registry lookup ───────────────────────────


class TestReelRegistryLookup:
    def test_get_returns_entry_for_reel_completed_v1(self):
        from agentfield.handoff import registry

        entry = registry.get(REEL_TYPE)
        assert entry.type_version == REEL_TYPE
        assert isinstance(entry.schema, dict)
        assert entry.schema.get("type") == "object"
        assert "reel_ref" in entry.schema.get("properties", {})

    def test_schema_has_all_dto_fields(self):
        from agentfield.handoff import registry

        entry = registry.get(REEL_TYPE)
        assert set(entry.schema["properties"].keys()) == _EXPECTED_KEYS

    def test_schema_forbids_additional_properties(self):
        from agentfield.handoff import registry

        assert registry.get(REEL_TYPE).schema.get("additionalProperties") is False

    def test_registered_alongside_research(self):
        from agentfield.handoff import registry

        type_versions = {e.type_version for e in registry.list()}
        assert REEL_TYPE in type_versions
        assert "com.silmari.research.completed.v1" in type_versions  # unchanged

    def test_contains_check(self):
        from agentfield.handoff import registry

        assert REEL_TYPE in registry
        assert "com.silmari.reel.completed.v99" not in registry

    def test_schema_path_points_to_real_file(self):
        from pathlib import Path

        from agentfield.handoff import registry

        assert Path(registry.get(REEL_TYPE).schema_path).exists()

    def test_status_enum_matches_research(self):
        from agentfield.handoff import registry

        assert registry.get(REEL_TYPE).schema["properties"]["status"]["enum"] == [
            "succeeded",
            "failed",
        ]

    def test_beat_count_is_integer_and_duration_is_number(self):
        from agentfield.handoff import registry

        props = registry.get(REEL_TYPE).schema["properties"]
        assert props["beat_count"]["type"] == "integer"
        assert props["duration_s"]["type"] == "number"


# ─────────────────────────── boundary validation (fail-closed) ───────────────────────────


class TestReelValidation:
    def test_conformant_event_passes(self):
        from agentfield.handoff import validate

        validate(GOLDEN_REEL_EVENT)

    def test_missing_required_field_raises(self):
        from agentfield.handoff import ValidationError, validate

        bad = {k: v for k, v in GOLDEN_REEL_DTO.items() if k != "reel_ref"}
        with pytest.raises(ValidationError):
            validate({"type": REEL_TYPE, "data": bad})

    def test_extra_field_raises(self):
        from agentfield.handoff import ValidationError, validate

        bad = {**GOLDEN_REEL_DTO, "video_path": "/leak/reel.mp4"}
        with pytest.raises(ValidationError):
            validate({"type": REEL_TYPE, "data": bad})

    def test_wrong_type_field_raises(self):
        from agentfield.handoff import ValidationError, validate

        bad = {**GOLDEN_REEL_DTO, "beat_count": "five"}
        with pytest.raises(ValidationError):
            validate({"type": REEL_TYPE, "data": bad})

    def test_invalid_status_enum_raises(self):
        from agentfield.handoff import ValidationError, validate

        bad = {**GOLDEN_REEL_DTO, "status": "running"}
        with pytest.raises(ValidationError):
            validate({"type": REEL_TYPE, "data": bad})

    def test_unregistered_before_registration_would_raise(self):
        # Sanity: an unknown reel version is not registered.
        from agentfield.handoff import ValidationError, validate

        with pytest.raises(ValidationError, match="no contract registered"):
            validate({"type": "com.silmari.reel.completed.v99", "data": GOLDEN_REEL_DTO})


# ─────────────────────────── freeze guard (immutability) ───────────────────────────


class TestReelFreezeGuard:
    def test_modifying_reel_v1_schema_fails(self):
        from agentfield.handoff.tools.check_frozen_contracts import check_frozen_contracts

        diff = ["M\tcontracts/com.silmari.reel.completed/v1.schema.json"]
        verdict, messages = check_frozen_contracts(diff)
        assert verdict == "fail"
        assert any("frozen" in m for m in messages)

    def test_adding_reel_v2_schema_passes(self):
        from agentfield.handoff.tools.check_frozen_contracts import check_frozen_contracts

        diff = ["A\tcontracts/com.silmari.reel.completed/v2.schema.json"]
        verdict, _ = check_frozen_contracts(diff)
        assert verdict == "pass"

    def test_deleting_reel_v1_schema_fails(self):
        from agentfield.handoff.tools.check_frozen_contracts import check_frozen_contracts

        diff = ["D\tcontracts/com.silmari.reel.completed/v1.schema.json"]
        verdict, _ = check_frozen_contracts(diff)
        assert verdict == "fail"
