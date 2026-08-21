<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `testing-and-vitest-expert` skill.
> Sibling topics in this family are now reference files under the hubs (`software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: testing-and-vitest-expert
description: >
  Vitest testing expert: writing tests, reviewing suites, diagnosing flaky behavior, choosing mocking strategies, tuning vitest.config.ts, and interpreting coverage and snapshot output. Covers vi.fn/vi.spyOn/vi.mock, async test patterns, V8 vs Istanbul coverage, browser mode providers (Playwright recommended), snapshot discipline, and configuration override behavior.
  TRIGGER: user writes, reviews, or debugs Vitest tests; asks about vi mocks, test fixtures, flaky tests, coverage setup, snapshot updates, or vitest.config; requests help migrating from Jest to Vitest.
  SKIP: Jest-only projects with no Vitest; Python pytest (use python-patterns); Playwright end-to-end tests without Vitest integration.
version: "1.1"
category: developer
updated: "2026-05-29"
tags:
  - vitest
  - testing
  - javascript
  - mocking
  - coverage
  - snapshots
  - browser-mode
related_skills:
  - programming-languages
---

# Testing and Vitest Expert

Practical Vitest reference for writing tests, reviewing suites, diagnosing flaky behavior, mocking, coverage, snapshots, and browser mode. Treat the official [Vitest Guide](https://vitest.dev/guide/), [API](https://vitest.dev/api/), and [Config](https://vitest.dev/config/) as the source of truth.

**Version note:** based on Vitest rolling docs as accessed 2026-05-10. Some features (object-form `retry`, `tags`, `meta`) are labeled Vitest 4.1+.

## When to use this skill

- Writing, reviewing, or debugging Vitest test files
- Choosing between `vi.fn`, `vi.spyOn`, and `vi.mock`
- Diagnosing flaky tests with `retry` and `repeats`
- Setting up or tuning coverage (V8 vs Istanbul)
- Understanding snapshot files and inline snapshots
- Configuring `vitest.config.ts` vs `vite.config.ts` override behavior
- Running tests in Browser Mode with Playwright or WebdriverIO
- Migrating from Jest's `done` callback to async/promise patterns

## When NOT to use this skill

- Jest-only projects without Vitest
- Python testing (use `python-patterns` skill)
- Playwright E2E tests outside Vitest (use `extension-e2e-testing`)
- General JavaScript debugging (use `javascript-node-html-css-debugging-expert`)

---

## Quick rules

1. Keep test files named with `.test.` or `.spec.` — Vitest's default pattern.
2. Use `async` tests and returned promises; avoid the legacy `done` callback.
3. Reset, clear, or restore mocks between tests to prevent state leakage.
4. Commit snapshot files alongside code; treat them like code changes in review.
5. Prefer Playwright over the default `preview` provider for CI-realistic browser tests.
6. Start with V8 coverage — it is the recommended default and matches Istanbul accuracy since v3.2.0.
7. `vitest.config.ts` overrides `vite.config.ts` — use `mergeConfig` for intentional extension.
8. Use `retry`/`repeats` to diagnose flakiness, not to mask it.

---

## Core Vitest model

- Framework: next-generation testing powered by Vite. Requires Vite ≥ 6.0.0 and Node ≥ 20.0.0.
- `test()` and `it()` are aliases. Omitting the body marks the test as `todo`.
- A promise-returning or `async` test body is awaited; rejection fails the test.
- Default timeout: 5 seconds. Override globally via `testTimeout` or per-test via options.

---

## Mocking

- Import `vi` from `vitest`, or enable `globals: true` in config to use it globally.
- **`vi.fn()`** — create a mock function. Must be cleared/restored between tests.
- **`vi.spyOn(target, key)`** — wrap an existing method. Some export-spy patterns do not work in Browser Mode.
- **`vi.mock(modulePath, factory?)`** — replace a module. **Hoisted to top of file.** Only affects external imports, not internal same-module calls.
- **`vi.setSystemTime(date)`** — mock the current date. Does not auto-reset; call `vi.useRealTimers()` in `afterEach`.

### Mock reset strategy

| Method | What it resets |
|--------|---------------|
| `vi.clearAllMocks()` | Call history and return values (keeps implementation) |
| `vi.resetAllMocks()` | Call history, return values, and implementation |
| `vi.restoreAllMocks()` | Restores `spyOn` targets to their original implementation |

---

## Async testing

- Prefer `async`/`await` or returning a promise. Avoid `done` callbacks.
- For concurrent snapshot tests, use the **test context `expect`** (`ctx.expect`), not the global one, so the correct test is associated with the snapshot.

---

## Coverage

| Provider | Default? | Notes |
|----------|----------|-------|
| `v8` | Yes | Recommended; AST-based remapping since v3.2.0; faster and lower memory |
| `istanbul` | No | May suit complex module-heavy projects |

Configure via `test.coverage.provider` in `vitest.config.ts`.

---

## Snapshots

- External snapshots: `toMatchSnapshot()` — written to `__snapshots__/` directory.
- Inline snapshots: `toMatchInlineSnapshot()` — stored in the test file itself. Use for small, local assertions.
- Update snapshots: `vitest --update-snapshots` (or `-u`).
- Commit snapshot artifacts. Review them in PRs like code changes.
- Formatting controlled by `snapshotFormat` config or custom serializers.

---

## Browser Mode

| Provider | Notes |
|----------|-------|
| `preview` | Default; simulates events (not CDP). Lower CI confidence. |
| Playwright | Recommended for fresh setups; real CDP, parallel execution. |
| WebdriverIO | Alternative real-browser provider. |

Some `vi.spyOn` export patterns do not work in Browser Mode. Test them in Node mode first, then adapt.

---

## Configuration

- `vitest.config.ts` has **higher priority** than `vite.config.ts` and overrides it rather than merging.
- To share Vite config, use `mergeConfig` explicitly:
  ```ts
  import { mergeConfig } from 'vitest/config';
  import viteConfig from './vite.config';
  export default mergeConfig(viteConfig, defineConfig({ test: { ... } }));
  ```
- Use `configDefaults` to extend defaults without retyping them:
  ```ts
  import { configDefaults, defineConfig } from 'vitest/config';
  export default defineConfig({
    test: { exclude: [...configDefaults.exclude, 'e2e/**'] }
  });
  ```
- Vite options belong at the top level; only Vitest `test` options go under `test`.

---

## API / pattern inventory

| API / pattern | Purpose | Key caveats |
|---|---|---|
| `test()` / `it()` | Define a test | Omitted body = `todo` |
| `test.skip()` / `{ skip: true }` | Skip a test | Easy to leave permanently |
| `testTimeout` / per-test `timeout` | Control timeout | Old timeout-as-last-arg cannot combine with options object |
| `retry` | Retry on failure | Object form is v4.1+; functions not allowed in serialized config |
| `repeats` | Re-run multiple times | Diagnostic only, not a correctness fix |
| `tags` | Group/filter tests | Fails if tags not declared with strict tag behavior enabled |
| `meta` | Reporter-visible metadata | Top-level keys merge; nested objects are not deep-merged |
| `vi.fn()` | Mock function | Reset/clear/restore between tests |
| `vi.spyOn()` | Spy on method/getter | Some export patterns break in Browser Mode |
| `vi.mock()` | Replace module | Hoisted; only affects external module access |
| `vi.setSystemTime()` | Mock current time | Does not auto-reset |
| `toMatchSnapshot()` | External snapshot | Commit artifacts; review in PRs |
| `toMatchInlineSnapshot()` | Inline snapshot | Clutter risk if overused |
| `coverage.provider` | Select V8 or Istanbul | V8 default and recommended since v3.2.0 |

---

## Testing best practices

### Structure and naming

- Name tests to describe expected behavior — the runner surfaces names in output and snapshot files.
- Prefer explicit string names over function-name inference.
- Keep test files alongside source or in a `__tests__` directory; follow project convention.

### Isolation

- Clear or restore mocks in `beforeEach` / `afterEach`.
- Reset time mocks and cached objects — they persist across test boundaries if not explicitly reset.
- Avoid shared mutable state between test cases.

### Mocking strategy

Escalate only as needed: function mock → spy → partial module mock → full module mock.

### Async discipline

- Default to `async`/`await`. Use `retry`/`repeats` for flaky diagnosis, not as a substitute for deterministic setup.

### Snapshot discipline

- Use snapshots for output-stability regression tests.
- Prefer inline snapshots for small assertions; external for large serialized output.
- Review snapshot diffs in PR review as carefully as code diffs.

### Coverage

- Use V8 first. Switch to Istanbul only if you have a concrete reason (complex instrumentation needs).
- Coverage changes are engineering decisions with tradeoffs — not cosmetic config switches.

### Flaky test prevention

- Keep mock state, timer state, and shared object state under control.
- Diagnose with `repeats` and `retry`; then fix the root cause.

---

## Version-sensitive notes

- Object-form `retry`, `tags`, and `meta` are Vitest 4.1+.
- V8 coverage accuracy reached Istanbul equivalence at v3.2.0.
- Browser Mode has explicit mocking limitations not present in Node mode.
- See the [Vitest changelog](https://github.com/vitest-dev/vitest/releases) for version-specific behavior.

---

## References

- [Vitest Guide](https://vitest.dev/guide/)
- [Vitest API](https://vitest.dev/api/)
- [Vitest Config](https://vitest.dev/config/)
- [Vitest Mocking](https://vitest.dev/guide/mocking.html)
- [Vitest Coverage](https://vitest.dev/guide/coverage.html)
- [Vitest Snapshots](https://vitest.dev/guide/snapshot.html)
- [Vitest Browser Mode](https://vitest.dev/guide/browser/)
