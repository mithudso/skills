# Verification-patch failure taxonomy

Classify every failed task/test in the bump PR's DSC verification patch into one
of these categories before rendering a verdict on the pin. Categories are
ordered roughly by how their signatures present. Historical instances are cited
so you can compare shapes, not so you can assume recurrence.

## 1. Uniform wall (all disagg exec tasks fail identically)

**Signature:** every disagg_storage exec task (~65) fails the same way, usually
a 600s `waitForService` timeout during fixture setup while an SLS container
fatally exits at startup. Generator tasks (no SLS launch) still pass.

**Meaning:** a mongod-harness vs SLS contract break at the new pin. Not flaky,
not test-specific. Read the failing service's container log for the fatal
error.

**Historical instance (flags rollout, SERVER-128073 era):** newer SLS builds
made `flags-state.json` mandatory per service, service by service across
successive pins (a moving target: first CMS+PAGE, later LOG). Each
newly-mandatory service fatally exited with
`Failed to initialize new flag system ... NotFound path:"etc/k8s/flags-state.json"`,
re-walling all tasks even after earlier services were wired.

**Two-harness gotcha:** there are TWO disagg SLS harnesses and a wall can clear
in one while persisting in the other:

- `src/mongo/db/modules/atlas/jstests/disagg_storage/libs/slstest.js`
  (SLSTest / SLSMinimalThreeCellTest / SLSBackupRestoreTest) starts individual
  containers via dockerRun.
- `libs/sls_fixture.js` (SLSMultiCellFixture) drives
  `buildscripts/modules/atlas/sls-multicell-docker-compose.yml` via docker
  compose (service names like `logd-cell3-0-1` give it away).

A harness config fix generally must be applied to both (service config blocks
in slstest.js AND the per-service blocks + top-level config defs in the compose
file). If a patch clears some tests but restart/reconfig-style tests still
fail, suspect the unfixed second harness.

## 2. Consistent partial failures (contract drift)

**Signature:** a stable subset of tests fails deterministically across
retries, with assertion or parse errors rather than startup walls.

**Meaning:** the new pin changed an observable contract the tests or harness
depend on. Usually a small mechanical fix. Known shapes:

- **Error-string rename:** SLS renamed "log segment not found in log store" to
  "... in log server" (SLS-5682); three log_server_startup_deleted_segment
  tests asserting the old string failed until updated.
- **Proto3 empty repeated field:** an empty list arrives as `{}`, not
  `{field: []}`. A harness `hasOwnProperty` assertion broke on
  GetLogSegmentList; fix shape is `response.field ?? []`.
- **Config schema drift:** SLS changed the CRS `remote_cells` config from a
  map to a sequence of `{cell, scheduler_uri}` (SLS-3590); the hardcoded old
  shape in the compose file and slstest.js had to move with it
  (SERVER-131421). Note only some fixtures exercise a given config path
  (SLSBackupRestoreTest starts CRS; base SLSTest does not), so coverage of the
  fix can be narrower than it looks.

Fixes in this category are the prime candidates for the Variant A
old-pin-safety test: a pure test-string or tolerant-parse fix is often
old-pin-safe; a config-shape change usually is not (the old binary may
hard-reject unknown fields; see category 6).

## 3. Known flakes (pre-existing, pin-independent)

**Signature:** a test that also fails on recent mainline runs at the OLD pin,
with no signature connecting it to the SLS commit range.

**Verify live, never from a remembered list.** Method: check the same test on
recent master waterfall / mainline disagg runs. If it fails there too, classify
as pre-existing and do not chase it in this workflow.

Example shape (2026-07 era): `catchup_step_up_checkpoint_exists_not_installed_kek.js`
failing with connection-refused + `assert.soon ... Could not find log entries
... 12850101`. It was a known mainline flake with no pin signature. Do not
assume it still is; re-verify.

Flake handling: once every real failure is explained, restart the flaky tasks
to get a clean run rather than merging on top of red.

## 4. Infra noise

**Signature:** failures in machinery around the tests, not the tests:

- S3 binary download failures / curl exit 56 during task setup.
- Hook-level noise (e.g. internode validation "Missing hashes" while the
  collections validate clean).
- Expired ECR docker login on locally-run repros (pull fails, exit 18).

Classify, restart, move on. Not a pin verdict input.

## 5. ECR image rot (old pin dead)

**Signature:** `ImageNotFoundException` / `manifest unknown` pulling
`disagg-storage/*` images tagged with the OLD pin's githash, typically on
mainline runs rather than the bump patch (the new pin's images are fresh).

**Meaning:** the SLS image retention in ECR (account 664315256653) is
count-based (keep newest 20000 per repo, roughly 3 weeks of hard guarantee at
observed push rates, sometimes more via prune backlog). A pin left stale long
enough loses its images, at which point master disagg goes red on its own and
the old-pin regression gate for any fix becomes impossible.

**Implication for shepherding:** rot raises urgency to land the bump and can
force a glass-break merge (bypassing the required gate). That trade-off is a
human decision: document the state and the option, never execute a glass-break
yourself. Note also that a bump merged via glass-break may have had reduced
verification coverage; read what its patches actually ran.

## 6. SLS-side regression (no mongod-side lever)

**Signature:** the failing mechanism lives inside the SLS binaries and no
harness/compose/config change on the mongod side reaches it.

**Case study (store_id, SLS-8204):** the new pin's scheduler placement began
hard-requiring a scraped `store_id`; log servers registered with
`store_id: None`; the multi-cell fixture's scrape didn't reliably win the
race, so `GetLogServers` returned `NotEnoughHealthyLogServersFound` and tests
hung in `startLog`. No static harness lever existed. Resolution came from
newer SLS commits improving the scrape; the correct move was escalating to SLS
with a precise question, and a later pin cleared it.

**Config placebo warning (from the same case):** a plausible-looking config
lever (`enforce_store_id: true`) appeared to fix a local run. It was a placebo
three ways: (1) the key was indented at the wrong level and serde silently
drops unknown YAML keys (no `deny_unknown_fields`), so the field was never
set; (2) even set, the field was on a different code path than the failure;
(3) the "passing" run had won the underlying race by luck. Lessons:

- A single local pass is weak evidence for a race-shaped failure.
- Verify a config lever is actually PARSED (right nesting, field exists in the
  pinned source) and is on the causal code path before shipping it.
- When you cannot show the lever's mechanism, treat it as a gamble and say so.

**Escalation form:** one concrete question stating observed mechanism plus the
ask, e.g. "In pin X, scheduler placement excludes log servers with
store_id:None; servers self-init a store_id but register with None. What is
the intended way for store_id to reach placement in the multi-cell compose
topology: config C, a bootstrap RPC, or is None-on-registration a regression?"

## 7. Post-merge fallout

**Signature:** the bump merged green, then the disagg variant breaks on
mainline afterwards (a BF ticket appears).

**Historical instance (BF-43221):** the first production bump committed a
stale `flags-state.json` because the bot's import allowlist was missing two
output paths; the SLS scheduler needed the fresh copy at runtime and the arm64
disagg variant hit 300s startup timeouts post-merge.

**Handling:** the ticket-per-bump design exists precisely so release managers
can revert an individual bump by ticket. Hand off to BF triage with the bump
ticket, pin range, and this revertability context. If the root cause is a bot
gap (as in BF-43221), file a bot fix ticket as well.
