# Team-specific BF Triage Knowledge

This file holds knowledge that is specific to the **team currently using
the skill**, not to MongoDB BF triage in general. The skill's main files
(`SKILL.md`, `reference/workflow_overview.md`, `reference/log_patterns.md`,
`reference/severity_types.md`) reference this file at decision points and
inherit its rules.

If you are a team other than the default below, **clear this file's
content and replace it with your own** — priority queues, front-line
ownership patterns, error-class examples, and disposition rules.
The shape of the file (the section headings used as anchors by the
main files) should stay roughly the same so the cross-references keep
working; the contents under each heading are yours to rewrite.

---

## Owning team

**Workload Resilience** (default content shipped with the skill).

Other teams: replace the section above and rewrite the rest of the file
with your own routing, prioritisation, and pattern rules.

## Triage priority queue

When this team's triager picks the next BF to work, the queue order is:

1. `rapid-response` labelled tickets (P0, daily comment required)
2. Hot master BFs counting toward `block-on-red` thresholds
3. Hot non-master BFs
4. AF tickets
5. Cold master BFs counting toward `block-on-red` thresholds
6. Cold non-master BFs

Master BF thresholds: > 5 cold master BFs at any one time triggers a
"drop what you're doing" mode for everyone on the team that has a master
BF.

## Keep-and-link routing pattern

This team does not always re-route a BF when the actual fix lives outside
its source code. The common pattern is:

1. Keep the BF assigned to this team as the front-line owner.
2. Open or link a downstream fix ticket in the appropriate project
   (e.g. `SERVER-XXXX` for server-side bugs, `PERF-XXXX` for performance
   /backpressure, `TUNE-XXXX` for tuning, dsi PR for workload code).
3. Move the BF to "Waiting for bug fix" until the linked ticket lands.

When the skill recommends "re-route", it should consider this pattern
first: if the workload itself or admission-control code is in the BF's
description, the recommendation should usually be **"Keep with the
owning team, file downstream fix ticket in `<project>`"** rather than
"Re-route to `<engineering team>`".

## Team front-line routing rules

These rows extend the generic re-routing decision matrix in
[`workflow_overview.md`](workflow_overview.md). When the symptom in the
left column matches, the front-line owner is this team. Other teams
should rewrite this section with the symptoms they own front-line.

| Symptom in logs | Disposition |
| --------------- | ----------- |
| Concurrency-suite hangs, FSM workload deadlock | Front-line: this team. |
| `assert.soon` from background hooks (CheckReplDBHash, CheckMetadataConsistency) | Front-line: this team for the cascade; downstream fix ticket against the true-owner team. |
| Genny / Locust / DSI infrastructure failure | Front-line: this team (DSI shared with this team). |
| Admission-control / Flow Control source code (`src/mongo/db/admission/`) | Front-line: this team. |
| Workload-client overload (the workload's runner reports a CPU / latency / queue-saturation error class) | Front-line: this team. Downstream fix ticket against PERF / TUNE for the resource pressure. |
| Server-side overload errors propagating to workloads (rate-limit-exceeded or admission-rejected error classes returned to the client) | Front-line: this team — the workload should catch the retryable error. Downstream fix only if the server side is mis-classifying retryable as fatal. |

## Workload-path-based routing override

When the failing locust workload lives under
`workloads/<team>/<workload-name>/` in the DSI repo, the team prefix in
the path is the canonical owner of the **workload code itself**. The
routing decision should:

- Identify the team prefix from the sliced description / log path
  (e.g. `workloads/<some-team>/<some-workload>/...` → `<some-team>`).
- Prefer **re-routing to the path-owner team** over the default
  "Keep with this team + downstream fix" pattern when
  (a) the failure is in the workload's Python code (parse error,
  ValidationError, undefined attribute, etc.) and
  (b) the team prefix is not the directory this team specifically
  owns availability of (for the default WR config, that directory is
  `workloads/availability/`).
- Cite the workload path as the routing justification.

## Variant-incompatibility patterns

### Forward direction — new variant added, workload not compatible

When the BF's first failing revision has a commit message of the form
"Add `<name>` variant for `<test>`", "Enable `<test>` on `<variant>`",
or "DEVPROD-XXXXX: Add `<...>` variants" — and the symptom is a
workload-side error class (workload runner exits non-zero, retryable
errors reported from the workload's client library) — the BF is often a
sign that the workload is **fundamentally incompatible** with the new
variant rather than a bug. Common shapes:

- The new variant lacks a setting the workload expects (e.g. a
  no-mongotune variant lacks a safety knob a mongotune variant
  provides).
- The new variant cannot exercise the failure the workload was
  designed for (e.g. a storage-offloaded variant where the disk never
  actually exhausts).
- The new variant collects artefacts the runner cannot store (e.g.
  oversized core-dump uploads from a crashed cluster overflow the DSI
  runner's disk).

In these cases the recommendation set is, in addition to "wrap
retryable errors in the workload":

1. **Disable the workload on the incompatible variant** in
   `configurations/test_control/*.yml` (DSI repo). This is a one-line
   change and is typical when the workload's failure mode cannot occur
   on the new variant.
2. **Suppress core-dump uploads** for variants where the workload is
   *expected* to crash the cluster.
3. **Coordinate with the variant author** named in the first-failing
   revision's commit message before proposing a fix — they may have a
   reason to keep the workload running there.

Treat "disable on variant" as a legitimate "Recommended next steps"
item alongside the more common "wrap in retry" recommendation when the
commit message points at variant addition.

### Inverse direction — required workload-config field added, test_control not updated

The mirror image of the variant-incompatibility pattern is also common
and is **NOT** variant-specific:

- A recent commit to `workloads/<team>/<workload-name>/src/*.py`
  added a new **required** Pydantic field to the workload's `Config`
  class (`<field>: <type> = Field(...)` with no default).
- The corresponding `configurations/test_control/<task>.yml` for one
  or more tasks was not updated to pass the new flag.
- The locust primary fails immediately at startup with `Error parsing
  option`, `pydantic.ValidationError`, or "missing required field" —
  workers then fail to connect because the master never came up.

The BF surfaces on the variants that run the un-updated `test_control`
file and looks variant-specific, but the actual fault is task-wide.

When triaging this shape:

1. Run `git_log` on `10gen/dsi` for the failing workload's
   `workloads/<team>/<workload-name>/src/` directory in the
   cutoff − 7d to cutoff window, looking for added `Field(...)` lines
   without default values in the workload's `Config` subclass.
2. Read the failing `configurations/test_control/<task>.yml` and check
   whether the new `locust_<field_name>` argument is present.
3. Recommend updating `test_control` to add the missing flag (Type of
   fix: Test change), or routing the BF to the team that owns the
   `workloads/<team>/<workload-name>/` directory rather than keeping
   it at this team.

### Sideways direction — workload-level config removed, per-variant migration incomplete

A common cousin of the inverse direction: a DSI commit removes a
setting that previously applied to *every* variant running the
workload (e.g. a workload-level `locust_poetry_override` /
`locust_poetry_dependencies` / a baseline `setParameter` line in
`configurations/test_control/<task>.yml`), with the stated intent of
moving control to the variant level. The PR then re-adds the setting
only inside the `expansions:` blocks of the variants the author
explicitly tested. Variants not in the PR's test matrix silently
switch to the default and regress — typically with very large
ErrorRate / ErrorsTotal jumps because the workload's open-loop
spike-user calibration was tuned against the now-removed setting.

Diagnostic shape:

- Perf-change-detector BF on a `*_locust` workload (or any workload
  with a meaningful `locust_poetry_override` / driver pin).
- DSI `git_log -- configurations/test_control/<failing-task>.yml`
  in the bisect window shows a commit removing one of:
  `locust_poetry_override`, `locust_poetry_dependencies`,
  `locust_*` driver pins, or a workload-level setParameter.
- The same commit's diff to
  `evergreen/system_perf/master/variants.yml` re-adds the setting
  only inside *some* `expansions:` blocks.
- The failing variant's `expansions:` block does not contain the
  setting (`git show <ref>:evergreen/system_perf/master/variants.yml |
  rg <setting-name>` returns hits only for the variants the PR
  added/migrated, not the failing one).
- A sibling variant on the same task is unaffected by the
  perf-change-detector's report — that sibling's `expansions:`
  contains the migrated setting.

Recommendation set:

1. **Restore the workload-level setting** in
   `configurations/test_control/<failing-task>.yml` (smallest change,
   preserves the original behaviour while the per-variant migration
   is completed). This is usually the right immediate fix.
2. OR **add the setting to every variant block that runs the failing
   task and currently lacks it.** Run an audit of all variants whose
   `tasks:` list includes the failing task and whose `expansions:`
   lack the setting. Cite the specific lines / variants in the audit.
3. File the fix ticket in the same project as the offending PR
   (typically `DEVPROD-XXXXX` or `TUNE-XXXXX`). The PR author is the
   right assignee — they have the original migration context.

Routing: keep the BF with the workload's owning team (this team for
WR's default config); re-route the *fix ticket* to the PR author's
team. The BF stays linked to the fix ticket via "Caused by" / "Is
caused by" until the variant audit lands.

## Intentional-stress / expected-failure workload pattern

Some perf workloads are *designed* to push the cluster into a known-bad
state — disk fills, oplog overflows, CPU saturates, change streams stall.
On those workloads a "System Failure" BF is sometimes the workload doing
its job, not a regression.

Diagnostic shape:

- Task name implies intentional resource exhaustion (regex match on
  `*_exhaustion_*`, `*_overflow_*`, `*_stress_*`), or the workload runs
  only on variant families whose name advertises "Dev Policies" /
  "Release Policies" / similar tuning-experiment markers.
- The bisect points at a workload-tuning commit (in the workload-tuning
  repo or in `10gen/dsi/workloads/` or a `test_control/*_stress*.yml`
  change), NOT at a server-side code change.
- The downstream symptom is fragility in the analysis path (parse
  crash on a truncated `mongod.log`, FTDC parse error, unbounded log
  volume) rather than a real correctness failure.

When this shape fires, the skill should treat **"declare expected,
remove or guard the failing test"** as a first-class recommendation
alongside (and often ahead of) the upstream-fix-ticket path. Concrete
recommendation set:

1. Recommend `Type of fix: Test change` with `Justification: This
   failure is expected. Removed failing test` (or "Disabled the
   failing workload on the failing variant"). This is the typical
   resolution for resolved BFs of this shape.
2. Open a separate hardening ticket against the **DSI / analysis
   side** (e.g. "harden analysis-side parser against truncated
   `mongod.log`") only if the analysis-side crash blocks other
   unrelated tasks; do NOT couple the BF resolution to that hardening
   ticket.
3. Do NOT recommend "Waiting for Bug Fix" against the upstream
   tuning-commit ticket unless the workload owner has confirmed the
   tuning change is itself a regression. The default for an
   intentional-stress workload is that the tuning change merely
   *exposed* the existing fragility.

Routing: keep with this team for the test-change disposition; re-route
to the workload-path owner only if the workload code itself must
change (rare for intentional-stress workloads).

## Performance-change BF disposition rules

Performance-change BFs (auto-generated by the change-point detector;
see [`log_patterns.md`](log_patterns.md) §
"Performance-change BF interpretation" for the detection rules) get
the following dispositions in this team:

- When the commit message of the first failing revision contains a
  `TUNE-` / `DEVPROD-` ticket reference, the fix is almost always in
  DSI / workload-tuning code, not in `mongo/`. Recommend checking
  `10gen/dsi` and the workload-tuning repo's `git_log` first, before
  bisecting `10gen/mongo`.
- **Even when `First Failing Revision` carries a `SERVER-` prefix**,
  the perf-change-detector's bisect locates a *mongo* SHA but the
  DSI module bumped in the same Evergreen run is an independent
  variable and can be the actual cause. Always run a parallel
  `git_log` on `10gen/dsi` for the failure window, especially when
  the failing variant is sys-perf-only / Mongotune-only / disagg-only
  / non-Atlas (these tend to share workloads with Atlas variants
  and so are exposed to per-variant config divergence). The
  DSI-side checklist is in [`SKILL.md`](../SKILL.md) § "Step 6 —
  DSI sub-checklist for sys-perf BFs".
- **When the candidate mongo SHA's diff is "inert in this
  configuration"** — gated by a default-`false` server parameter not
  enabled by the variant's Mongotune policy / setParameter list, or
  in a code path the workload doesn't exercise — redirect to DSI
  investigation immediately. Do NOT wait for the suspect commit's
  author to confirm the inertness in a Jira comment; the gating
  evidence (parameter default, setParameter list, workload call
  graph) is sufficient on its own. Cite the gating parameter and the
  reasoning in the report. The full rule lives in
  [`SKILL.md`](../SKILL.md) § "Step 6 — Inert-mongo-diff redirect
  rule".
- When the diagnostic shape matches the **Sideways direction —
  workload-level config removed, per-variant migration incomplete**
  pattern (see § "Variant-incompatibility patterns" above), the
  default recommendation pair is (a) restore the workload-level
  setting OR add it to the failing variant's `expansions:`, and (b)
  file the fix ticket against the migration PR's project (typically
  `DEVPROD` or `TUNE`). Keep the BF with this team; re-route only
  the fix ticket.
- When the suspected mechanism is a controller-tuning artefact (e.g.
  an older controller hard-codes a value that a new policy expects to
  be configurable), the perf-change BF will often be closed as
  `Type of fix = Non-test code change, Resolution = Gone away` even
  though no obvious revert happened. Mention this as a plausible
  outcome rather than insisting the change-point be reverted.
- When the BF was auto-generated and the summary already names a
  fix-by commit ("fixed by `<sha>`, `<date>`"), default recommendation
  is **Accept-as-known, close as `Gone away`** after `git_show`
  confirms the caused-by and fixed-by commits cancel each other.

---

## When to add new content to this file (Mode C learning)

Mode C (verification / replay) is where the skill discovers new
patterns by grading its own triage against the real resolution. When a
graded BF reveals a new pattern, the **classification rule** below
decides whether the new content goes into a shared file or into this
team file:

- **Shared file** (`SKILL.md`, `workflow_overview.md`,
  `log_patterns.md`, `severity_types.md`) — when the pattern is true
  for **any** team using the skill: a new MongoDB log marker, a new
  severity-type identifier, a generic workflow step, a new tool-failure
  recovery procedure, a new gateway / CLI integration.
- **Team file** (`team_knowledge.md`, this file) — when the pattern
  encodes **who** owns a symptom, **how this team specifically routes
  it**, or **which projects this team links downstream fix tickets in**.
  Concretely: new front-line routing rules, new error-class examples
  the team sees, new variant-family names the team handles, new
  workload-path conventions, new disposition defaults.

Heuristic: if the new rule contains a team name, a workload-tuning
repo / project prefix, or a per-team priority weight, it's team-file
content. If it could be written without naming any team, it's
shared-file content.

See [`verification_mode.md`](verification_mode.md) §
"Promoting findings to skill content" for the full Mode C workflow.
