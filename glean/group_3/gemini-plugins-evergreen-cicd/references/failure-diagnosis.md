# Evergreen CI/CD Failure Diagnosis — Detailed Reference

## MCP Tool Usage Strategy

**ALWAYS use MCP tools first** for Evergreen failure analysis:

1. Use `list_user_recent_patches_evergreen` to get build patch status
2. Use `get_patch_failed_jobs_evergreen` to analyze failures
3. Use `get_task_test_results_evergreen` to get specific failed test names and errors
4. Use `get_task_logs_evergreen` only when job name and general failure are insufficient

IMPORTANT: If a job is failing on master, do not attempt to fix it.

## Compilation Failures

**Issue**: Code fails to compile due to syntax errors, missing includes, or dependency issues

**MCP Strategy**:

1. Use `get_patch_failed_jobs_evergreen` to get failure details
2. Check logs for the first compilation error — fixing it often resolves cascading errors
3. Do not attempt to compile the entire codebase locally as this takes significant time

**Common Patterns**:

- Missing includes or imports after moving or adding code
- Type errors or template instantiation errors
- Linker errors from missing or renamed symbols
- Build file issues (missing dependencies, wrong targets)

**Local Verification**:

Check your repo-specific references file (e.g., `references/server.md`, `references/mms.md`) for local build commands.

## Unit Test Failures

**Issue**: Unit tests failing due to code changes
**Priority**: Medium — fix after compilation issues

**MCP Strategy**:

1. Use `get_patch_failed_jobs_evergreen` to identify the failing task
2. Use `get_task_test_results_evergreen` to get specific failed test names
3. Use `get_task_logs_evergreen` to get assertion failure details and stack traces

**Common Patterns**:

- Assertion failures
- Segfaults or memory errors (check for use-after-free, null dereference)
- Timeout failures (test took too long)
- Missing test fixtures or setup failures

**Local Verification**:

Check your repo-specific references file (e.g., `references/server.md`, `references/mms.md`) for local test commands.

**Debugging Steps**:

1. Read test failure messages from MCP logs carefully
2. Check if code changes broke test assumptions
3. Update tests if behavior change is intentional
4. Fix code if tests reveal actual bugs
5. For segfaults, look at the stack trace to identify the crash location

## Integration Test Failures

**Issue**: Integration tests failing due to behavior changes or environmental issues
**Priority**: Medium — fix after compilation and unit test issues

**MCP Strategy**:

1. Use `get_patch_failed_jobs_evergreen` to identify failed integration tasks
2. Use `get_task_test_results_evergreen` to get the specific failing test files
3. Use `get_task_logs_evergreen` to check for errors, assertion failures, and test output

**Important**: Log files can be very large. Fetch from the bottom and stop after identifying the first issue.

**Common Patterns**:

- Test assertion failures
- Service crashes during test (look for stack traces in logs)
- Timeout failures (test or service took too long to respond)
- Configuration issues
- Changed behavior or error codes
- External dependency issues

**Local Verification**:

Check your repo-specific references file (e.g., `references/server.md`, `references/mms.md`) for local test commands.

## Lint and Format Failures

**Issue**: Code quality checks failing
**Priority**: High — fix early to prevent blocking commits

**Common Patterns**:

- Linting violations (code style, static analysis)
- Formatting violations (indentation, spacing)

**Local Verification**:

Check your repo-specific references file (e.g., `references/server.md`, `references/mms.md`) for local lint and format commands.

## Application-Specific Debugging

### Core Dumps and Crashes

When a process crashes during testing, look for:

1. "core dump" or "signal" messages in task logs
2. Stack traces or backtrace information
3. The crash location to identify the root cause

### Timeout Analysis

For timeout failures, look for:

- Stack traces of threads or processes
- Lock or resource contention information
- Information about operations in progress

### Application Log Files

Tests may produce application log files. Key things to look for:

- Error or fatal log messages
- Stack traces
- Assertion or invariant failure messages
- Connection errors or timeout messages

## Exclusions — DO NOT FIX

### Known Flaky Tests

- Tests that fail intermittently and are unrelated to your changes
- Check if the test is on the known flaky list or has failed on master recently

### Infrastructure Issues

- **Host allocation failures** — Evergreen could not provision a host
- **Archive failures** — Artifact upload/download issues
- **Network timeouts** — Transient CI infrastructure problems
- **Resource exhaustion** — Escalate to the Build team

## Deduplication Strategy

### Fix Compilation First

- Compilation failures often cause cascading test failures
- Fix `archive_dist_test` tasks before investigating test failures
- Re-evaluate other failures after compilation fixes

### Group Related Failures

- Multiple test failures in the same suite likely share a single root cause
- Unit test failures + integration test failures in the same area suggest a code bug
- Lint failures are independent and can be fixed in parallel with other issues

### Wait for Master Fixes

- If failures appear unrelated to patch changes, check if they also fail on master
- Use MCP tools to compare your patch results with the master waterfall
- Avoid fixing issues that are already being addressed upstream

## Emergency Procedures

### Getting Help

- **Build team**: For infrastructure and CI configuration issues
- **Server team leads**: For persistent test failures in specific areas
- **Escalation path**: Check task logs for owner information, then escalate to team leads
