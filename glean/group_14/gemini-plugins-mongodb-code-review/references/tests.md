# Sub-Agent: Test Review

Tag: `tests`

Pass the test file diffs to the sub-agent with this prompt:

> Review the tests in this MongoDB PR.
>
> Reference docs/unit_test.md if this PR is targeting the 10gen/mongo codebase and includes unit
> tests (in a <filename>_test.cpp file).
>
> Focus on:
>
> **Coverage (blocking if gaps are significant)**
>
> - Failure/negative case tested, not just happy path
> - Edge cases: empty input, single element, max size, concurrent access
> - New config knobs / feature flags have tests for both on and off
> - Test exercises the actual code path added
> - Deployment styles: sharded vs unsharded, mongos presence, replica set configurations
>
> **Assertion Quality (non-blocking or blocking)**
>
> - Assertions verify behavior, not just "no crash" — ASSERT_EQ > ASSERT_DOES_NOT_THROW
> - Assertions are tight enough to catch regressions (eq vs gte)
> - JS tests use named arguments for readability
> - Assertion messages include enough context to understand the failure without re-running
>
> **Diagnosability & Flakiness (non-blocking)**
>
> - Test failure message is self-explanatory without reading the code
> - Prefer doTest(scenario1); doTest(scenario2); over looping test cases (stack trace clarity)
> - No hardcoded sleep / timing dependencies — use assert.soon with a message
> - Test doesn't over-couple to implementation details that will change
> - Passthrough suite compatibility — feature flag tag present where needed
> - Replication / sharding state machines: waiting for the right state?
>
> **Readability & Structure (optional or style)**
>
> - Test name follows SuiteName_State_ExpectedBehavior convention
> - Setup/teardown clearly separated from assertions
> - Common setup extracted to a helper rather than repeated
> - Use st.shardColl() / standard helpers rather than reimplementing
> - Readers should understand coverage quickly; don't bury assertions in setup boilerplate
