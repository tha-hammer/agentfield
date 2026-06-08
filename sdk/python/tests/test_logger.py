import logging
import threading

import pytest

import agentfield.logger as logger_module
from agentfield.logger import (
    get_logger,
    log_info,
    set_log_level,
)


@pytest.fixture(autouse=True)
def reset_logger_state():
    """Isolate each test from global logger state.

    Two distinct kinds of global state have to be reset:

    1. The SDK-level caches (``_logger_cache`` / ``_global_log_level``). These
       are referenced via ``logger_module`` so the reset hits the *real* module
       globals — assigning a bare module-level name here would only rebind this
       test module's copy and silently fail to reset anything.
    2. The stdlib ``logging`` registry. ``AgentFieldLogger`` attaches a handler
       and sets ``propagate = False`` on the underlying ``logging.Logger``.
       Those mutations outlive the test and leak across the whole session — e.g.
       creating the ``"agentfield"`` logger here would stop ``"agentfield.cancel"``
       records from reaching root handlers (pytest's ``caplog``), breaking
       unrelated tests like ``test_cancel.py``. Snapshot and restore so this
       test file is order-independent.
    """
    manager = logging.root.manager
    saved = {
        name: (lgr.propagate, list(lgr.handlers))
        for name, lgr in manager.loggerDict.items()
        if isinstance(lgr, logging.Logger)
    }

    logger_module._logger_cache.clear()
    logger_module._global_log_level = None
    try:
        yield
    finally:
        logger_module._logger_cache.clear()
        logger_module._global_log_level = None
        for name, lgr in list(manager.loggerDict.items()):
            if not isinstance(lgr, logging.Logger):
                continue
            if name in saved:
                propagate, handlers = saved[name]
                lgr.propagate = propagate
                lgr.handlers[:] = handlers
            else:
                # Logger created during the test — return it to stdlib defaults.
                lgr.propagate = True
                lgr.handlers.clear()


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


@pytest.mark.unit
def test_concurrent_access_is_threadsafe():
    """Concurrent get_logger()/set_log_level() must not raise or corrupt the cache.

    Without locking, iterating the cache in set_log_level() while get_logger()
    inserts into it can raise ``RuntimeError: dictionary changed size during
    iteration``. Hammer both paths from several threads and assert clean results.
    """
    errors: list[BaseException] = []
    barrier = threading.Barrier(8)

    def worker(i: int) -> None:
        try:
            barrier.wait()
            for j in range(50):
                get_logger(f"concurrent.{i}.{j}")
                set_log_level("DEBUG" if j % 2 else "INFO")
        except BaseException as exc:  # noqa: BLE001 - surface any thread failure
            errors.append(exc)

    threads = [threading.Thread(target=worker, args=(i,)) for i in range(8)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert not errors, f"thread-safety violation: {errors[:3]}"
    # Every distinct name must have produced its own correctly-named logger.
    assert get_logger("concurrent.0.0").logger.name == "concurrent.0.0"
