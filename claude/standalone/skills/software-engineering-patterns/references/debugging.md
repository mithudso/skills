<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `debugging` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: debugging
description: >
  Systematic 7-phase debugging workflow for diagnosing and fixing bugs: symptom collection, environment analysis, unit test execution, integration testing, root cause analysis, fix implementation, and verification. Includes common bug patterns, deep trace logging templates, and a prevention checklist.
  TRIGGER: "bug", "issue", "not working", "broken", "error", "crash", "exception", "failed", "unexpected behavior", "wrong output", "race condition", "balance wrong", "cache issue", "data mismatch", "investigate this", "why is X happening".
  SKIP: performance profiling or optimization without a bug (use performance-profiling-expert or debugging-strategies); general architecture questions unrelated to a bug; production incident response requiring runbooks (use incident-response).
version: 1.1.0
updated: "2026-05-29"
category: developer
tags: [debugging, root-cause-analysis, testing, bug-fix, investigation]
related_skills:
  - debugging-strategies
  - performance-profiling-expert
  - coding-standards
  - testing-and-vitest-expert
whenToUse:
  - diagnosing a bug or unexpected behavior
  - tracking down a crash, exception, or error
  - investigating data integrity issues (wrong values, mismatches, stale cache)
  - finding and fixing race conditions or concurrency bugs
  - verifying a fix and preventing regression
whenNotToUse:
  - performance optimization without a specific bug (use performance-profiling-expert)
  - production incident response requiring runbook coordination (use incident-response)
  - general code quality review without a reported bug (use code-reviewer)
---

# Systematic Debugging

A 7-phase workflow for diagnosing and fixing bugs. Follow the phases in order — skipping steps leads to misdiagnoses.

## Key Principles

- **Systematic over intuitive** — follow phases; don't jump to fixes
- **Evidence-based** — use test results, logs, and code evidence, not assumptions
- **Minimal changes** — fix the root cause, not symptoms
- **Test-first** — write a test for the bug before implementing the fix
- **Prevent recurrence** — every fix ships with a regression test

---

## Phase 1: Symptom Collection & Reproduction

**Objective:** Understand what's broken and reproduce it consistently.

1. **Gather information**
   - What is the user experiencing?
   - When did it start? After what change?
   - Is it reproducible or intermittent?
   - Any error messages or log entries?
   - Which components are affected?

2. **Ask clarifying questions**
   - What steps lead to the issue?
   - Expected vs actual behavior?
   - Does it happen in all contexts or specific scenarios?
   - Any recent code, config, or dependency changes?

3. **Document minimal reproduction steps**
   - Create the smallest set of steps that reliably triggers the issue
   - Identify scope: single feature, multiple features, or system-wide
   - Note prerequisites and required state

---

## Phase 2: Environment & Context Analysis

**Objective:** Understand the system state around the failure.

```bash
# Review recent changes
git log --oneline -20
git diff HEAD~5 -- relevant/files

# Search for related past fixes
git log --all --oneline | grep -i "fix\|bug"
git show <commit-hash>
```

- Check CONFIG settings and defaults
- Review relevant data/sheet structure
- Look for related cached state
- Review CLAUDE.md "Known Issues" section
- Check architectural constraints

---

## Phase 3: Unit Test Execution & Analysis

**Objective:** Identify the scope of the failure and validate assumptions.

1. Run the test suite for affected modules
2. Analyze which tests pass/fail and whether failures match the reported issue
3. Check error messages in Logger/console output
4. **Document test gaps** — is the bug scenario covered? Create a test case for it.

---

## Phase 4: Integration Testing

**Objective:** Understand how modules interact around the failure.

1. Run relevant integration test scenarios
2. Trace data flow through modules: entry point → processing → output
3. Verify assumptions at each boundary:
   - Data created/saved correctly?
   - Cache synchronized after mutations?
   - Balances calculated after updates?
   - Audit logs recorded?

---

## Phase 5: Root Cause Analysis & Deep Trace

**Objective:** Find the exact cause of the issue.

### Code inspection
- Read the primary module handling the failing scenario
- Look for unhandled edge cases and missing error handling
- Check error propagation (try/catch blocks)

### Deep trace logging

Add strategic logging to trace execution:

```javascript
/**
 * Deep trace for [issue description]
 * Related to: [component/feature]
 * Test case: [reproduction steps]
 */
function debugFunction(data) {
  Logger.log(`[MODULE] Entry: data=${JSON.stringify(data)}`);
  try {
    const step1 = calculateSomething(data);
    Logger.log(`[MODULE] After step1: ${JSON.stringify(step1)}`);

    const step2 = updateCache(step1);
    Logger.log(`[MODULE] After step2 (cache): ${JSON.stringify(step2)}`);

    return { success: true, data: step2 };
  } catch (error) {
    Logger.log(`[MODULE] ERROR: ${error.message}`);
    Logger.log(`[MODULE] Stack: ${error.stack}`);
    throw error;
  }
}
```

### Common bug patterns (check these first)

| Pattern | Symptom | Investigation |
|---------|---------|---------------|
| **Cache invalidation timing** | Stale data, balance wrong | Is cache cleared before or after the write? Always write first, then clear cache. |
| **Formula vs values** | Unexpected string values | Storing formula text instead of evaluated result? |
| **Lock scope** | Intermittent failures | Lock released too early or not released? |
| **User resolution** | Wrong user recorded | `Session.getActiveUser()` vs `UserResolver.getCurrentUser()`? |
| **Trigger context** | Works manually, fails in trigger | Simple trigger accessing restricted APIs? Use installable trigger. |
| **Null/undefined** | TypeError at runtime | Unhandled null in data flow? |
| **Type mismatch** | Comparison fails | String vs number, date formatting? |
| **Concurrent access** | Intermittent duplicates | Race condition without proper locking? |

### Hypothesis testing

1. Form a specific hypothesis: "Bug occurs because X happens before Y"
2. Write a test case to prove or disprove it
3. Add logging to validate the hypothesis
4. Trace execution with real data

---

## Phase 6: Fix Implementation

**Objective:** Implement the minimal, targeted fix.

1. **Design the fix**
   - What is the smallest change that resolves the root cause?
   - Does it align with module responsibilities?
   - Will it break existing functionality?
   - Are there other similar issues in the codebase?

2. **Implement**
   - Apply a minimal change to the source file
   - Keep changes focused on the root cause
   - Add defensive code if appropriate
   - Include explanatory comments

3. **Fix pattern example — timing issue:**

```javascript
// BAD: Cache cleared before data is written
cache.clear();
log.appendRow(data);

// GOOD: Write data first, then invalidate cache
log.appendRow(data);
cache.clear();
```

4. **Add logging for future debugging** using consistent prefixes: `[MODULE] message`

---

## Phase 7: Verification & Regression Testing

**Objective:** Prove the fix works and doesn't break anything else.

1. **Test the specific bug scenario** — run reproduction steps; confirm issue is gone
2. **Run the full test suite** — all tests must pass; no new failures; no performance regression
3. **Integration testing** — test with batch operations, concurrent access, edge cases
4. **Create a prevention test case:**
   - Confirms test fails without fix
   - Confirms test passes with fix
   - Documents what regression it prevents
5. **Review for side effects** — other modules, caching, audit logs, performance

---

## Investigation Checklist

- [ ] Issue is reproducible
- [ ] Reproduction steps documented
- [ ] Recent code changes reviewed (`git log`)
- [ ] Relevant test files identified and run
- [ ] Unit tests analyzed
- [ ] Integration tests run and analyzed
- [ ] Known gotchas reviewed (CLAUDE.md)
- [ ] Similar past issues checked (`git log | grep fix`)
- [ ] Root cause identified with evidence
- [ ] Minimal fix designed and implemented
- [ ] Fix tested and verified
- [ ] Full test suite passes
- [ ] Prevention test case added
- [ ] No side effects identified

---

## Output Format

For each investigation, produce:

1. **Symptom Analysis** — what's broken and under what conditions
2. **Environment Context** — recent changes, affected modules, relevant state
3. **Test Results** — unit and integration test findings
4. **Root Cause** — why the bug occurs, with specific evidence
5. **Fix Implementation** — minimal code change with explanation
6. **Verification** — testing approach and results
7. **Prevention** — test case or safeguard added to prevent recurrence
