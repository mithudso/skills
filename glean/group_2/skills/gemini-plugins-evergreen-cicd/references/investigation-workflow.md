# CI/CD Failure Investigation Workflow

Systematic workflow for investigating and fixing CI/CD test failures in Evergreen.

**Tools**: Evergreen MCP (`list_user_recent_patches`, `get_patch_failed_jobs`, `get_task_logs`, `get_task_test_results`), `codebase-retrieval`, repo build tools
**Critical**: NEVER investigate or commit without user approval, NEVER fix master branch failures

## Phase 1: Retrieve Failing Tasks

Use `list_user_recent_patches_evergreen` (limit: 5-10), identify failed patch, get failed jobs with `get_patch_failed_jobs_evergreen`, extract task names/IDs, variants, test counts, errors.

## Phase 2: Analyze and Prioritize Failures

**CRITICAL**: Do NOT investigate yet. Only analyze and categorize.

**Categorize**: PRIORITY 1: Compilation failures -> PRIORITY 2: Lint/format failures -> PRIORITY 3: Unit test failures -> PRIORITY 4: Integration test failures -> LOWER: Other tasks

**Determine**: Code relationship (compare modified files), flaky tests (no related changes), group related failures (same module or suite).

**Analysis Format**:

```markdown
## Evergreen Failure Analysis - Patch #{patch_id}

### PRIORITY 1: Compilation Failures

- Task: <compilation_task> | Status: Failed
  Error: "error: no matching function for call to 'foo'"
  Related: YES (modified src/path/to/foo.cpp)

### PRIORITY 2: Lint/Format

- Task: <lint_task> | Status: Failed
  Error: "Formatting differences found in 3 files"
  Related: YES (modified files need formatting)

### PRIORITY 3: Unit Test Failures

- Task: <unit_test_task> | Failed Tests: 2
  Tests: FooTest.TestBar, BazTest.TestQux
  Related: YES (modified src/path/to/foo.cpp)

### PRIORITY 4: Integration Test Failures

- Task: <integration_test_task> | Failed Tests: 1
  Tests: tests/integration/some_test
  Related: MAYBE (changed related behavior)

RECOMMENDED ORDER: Fix compilation -> lint/format -> unit tests -> integration tests
```

## Phase 3: Get User Approval on Investigation Scope

Present analysis, recommend compilation/lint first, note any suspected flaky tests. **Wait for explicit approval** before proceeding.

## Phase 4: Investigate Failures Using Evergreen Logs

For approved failures: Fetch logs (`get_task_logs_evergreen`, filter_errors=true, max_lines=500-1000), get test results (`get_task_test_results_evergreen`, failed_only=true), analyze errors, use `codebase-retrieval` for context.

**Important**: Integration test logs are large — fetch from bottom, stop after identifying issue.

## Phase 5: Run Tests Locally if Needed

If Evergreen logs are insufficient, run specific failing tests locally.

Check your repo-specific references file (e.g., `references/server.md`, `references/mms.md`) for local test commands, build commands, and formatting commands.

**Important**: Do NOT attempt to compile the entire codebase locally unless necessary — use Evergreen logs for compilation errors when possible.

## Phase 6: Implement Fixes

Use `codebase-retrieval` for context, make minimal changes, update tests if behavior changed, fix bugs, handle downstream impacts (callers, types, interfaces).

## Phase 7: Verify Fixes Locally

**CRITICAL**: Confirm fixes work before committing. Run previously failing tests, related tests (no regressions), and formatting.

Check your repo-specific references file (e.g., `references/server.md`, `references/mms.md`) for local test and formatting commands.

## Phase 8: Commit and Push

**CRITICAL**: NEVER commit/push without approval. Request approval with summary (fixes, files, test results), then commit with descriptive message (include ticket/patch), push to branch, confirm success.

## Key Principles

1. **NEVER** investigate or commit/push without user approval, **NEVER** fix master branch failures unless that's the primary focus of the session
2. **Always** use Evergreen MCP tools first, prioritize compilation, note suspected flaky tests
3. **Follow** TDD (fail -> fix -> pass), make minimal focused changes
4. **Quality**: All tests pass locally, follow style guidelines, no unrelated changes

## Error Handling

**API**: Rate limiting (wait/retry), auth (verify MCP), patch not found (check ID)
**Investigation**: Insufficient logs (run locally), persistent failures (escalate), master failures (notify, don't fix)
**Local Tests**: Build failures (check deps with Bazel), test environment (verify mongod is built), timeouts (check for hangs)

## Quality Checklist

- [ ] Failures analyzed, user approved scope, prioritized correctly (compilation -> lint -> unit -> integration)
- [ ] Flaky tests identified, Evergreen logs analyzed, local tests only if needed
- [ ] Fixes minimal/focused, failing tests now pass, related tests pass (no regressions)
- [ ] `bazel run format` passes
- [ ] User approved commit/push
