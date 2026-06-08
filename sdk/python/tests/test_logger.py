import pytest

from agentfield.logger import (
    AgentFieldLogger,
    get_logger,
    log_info,
    set_log_level,
    _logger_cache,
    _global_log_level,
)


@pytest.fixture(autouse=True)
def reset_logger_cache():
    """Reset the logger cache before and after each test."""
    _logger_cache.clear()
    global _global_log_level
    _global_log_level = None
    yield
    _logger_cache.clear()
    _global_log_level = None


@pytest.mark.unit
def test_get_logger_returns_correct_name():
    """Test that get_logger returns logger with correct name."""
    logger = get_logger("agentfield.client")
    assert logger.logger.name == "agentfield.client"


@pytest.mark.unit
def test_different_names_produce_different_loggers():
    """Test that different names produce different logger instances."""
    a = get_logger("module_a")
    b = get_logger("module_b")
    assert a.logger.name == "module_a"
    assert b.logger.name == "module_b"
    assert a is not b


@pytest.mark.unit
def test_same_name_returns_same_logger():
    """Test that requesting the same name returns the cached logger."""
    a = get_logger("module_a")
    b = get_logger("module_a")
    assert a is b


@pytest.mark.unit
def test_set_log_level_affects_all_loggers():
    """Test that set_log_level affects all existing logger instances."""
    a = get_logger("a")
    b = get_logger("b")

    set_log_level("DEBUG")

    assert a.log_level == "DEBUG"
    assert b.log_level == "DEBUG"


@pytest.mark.unit
def test_set_log_level_applies_to_new_loggers():
    """Test that set_log_level also applies to loggers created after it's called."""
    set_log_level("DEBUG")
    a = get_logger("a")
    b = get_logger("b")

    assert a.log_level == "DEBUG"
    assert b.log_level == "DEBUG"


@pytest.mark.unit
def test_log_info_works_without_arguments():
    """Test that log_info works without explicit logger arguments."""
    log_info("test message")


@pytest.mark.unit
def test_default_name_returns_agentfield():
    """Test that get_logger() without arguments returns logger with name 'agentfield'."""
    logger = get_logger()
    assert logger.logger.name == "agentfield"


@pytest.mark.unit
def test_backward_compatibility_default_logger():
    """Test backward compatibility: get_logger() returns proper default logger."""
    logger1 = get_logger()
    logger2 = get_logger("agentfield")
    assert logger1 is logger2
    assert logger1.logger.name == "agentfield"
