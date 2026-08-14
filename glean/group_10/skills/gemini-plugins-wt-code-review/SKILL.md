---
name: wt-code-review
description: Use whenever the user asks for a code review, audit, or sanity-check of a WiredTiger change — a PR, a branch, a diff, or unstaged working-tree edits under `src/`. Triggers on phrasing like "review this PR", "review my WT changes", "audit this diff", "check this before I push", "look this over", "is this WT change OK", "do a deep review of WT-XXXXX", as well as on bare pasted diffs that touch WiredTiger source. Orchestrates the review in three steps — acquiring the diff, dispatching parallel specialised reviewers plus git history, then summarising. Each reviewer loads the relevant topic files from `references/` before inspecting the diff, so findings stay grounded in the team's accumulated reviewer wisdom. Do not invoke for code outside the WiredTiger repository.
source: 10gen/agent-skills
when_to_use: orchestrated parallel review of a WiredTiger diff by specialised reviewers
license: Internal
mongodb:
  team: storage-engines
  owner:
    - sean.watt@mongodb.com
    - ayesha.ahmed@mongodb.com
  internal: true
---

# Orchestrating a WiredTiger code review

This skill turns a single "review this" request into a coordinated, parallel review of a WiredTiger change. The orchestrator stays high level and diff-free: acquire the diff, dispatch reviewer subagents, summarise. Each reviewer's domain knowledge lives in `references/`.

Work through the three steps below in order. Create a todo list before starting, one item per step, and mark each completed as it finishes.

## Acquiring the diff

Resolve the user's request to a concrete diff *once*, here. The reviewer subagents trust the path you give them; capturing it centrally guarantees every reviewer sees the same snapshot.

Map the user's scope to a command:

- PR number → `gh pr diff <num>`
- Branch → `git diff develop...HEAD` (or `main...HEAD` if the repo uses `main`)
- Working tree → `git diff` (staged + unstaged), or `git diff --staged` for "what I'm about to commit"
- Pasted diff → write it to disk verbatim

Save the diff to `/tmp/wt-review-<short-id>.diff`. Hold onto that path — every reviewer prompt must include it. If `git diff` returns nothing, stop and ask the user what to review.

Capture a one-line **context** string: the PR number or branch, the base ref, and any scope hint the user gave ("just style", "skip tests", "focus on disagg"). Pass it verbatim to every reviewer.

If the diff is not WiredTiger (different repo, no `src/` paths matching the WT layout), stop and tell the user — the reference files assume the WiredTiger codebase and will produce noise on anything else.

Do not read the diff into your own context. You route paths and findings; you do not review.

## Dispatching reviewers in parallel

Issue **one message containing one `Agent` call per group below plus the history agent**, in a single turn. The runtime fans these out concurrently. There is no triage and no confirmation step — every group runs on every diff. Respect an explicit user scope hint ("just style", "only the disagg parts") by dispatching only the matching group(s); otherwise dispatch all of them.

Each group reviewer reads its reference file(s), then the diff. Groups are non-overlapping by design. Pass the reviewer the absolute path `${CLAUDE_SKILL_DIR}/references/<file>.md` for each file listed.

| Group | Reference files |
|---|---|
| core-c | `references/wt-correctness-reviewer.md`, `references/wt-concurrency-reviewer.md`, `references/wt-assert-reviewer.md`, `references/wt-error-cleanup-reviewer.md` |
| interfaces | `references/wt-api-reviewer.md`, `references/wt-disagg-reviewer.md` |
| hygiene | `references/wt-style-reviewer.md`, `references/wt-perf-reviewer.md`, `references/wt-logging-reviewer.md` |

For each group, dispatch with `subagent_type: "general-purpose"`, model inherited, `description` like `"Correctness review of <PR/branch>"`, and this prompt (substitute the group's reference list):

```
You are reviewing a WiredTiger code change for one group of related surfaces.
The checklists for what to flag and how to calibrate severity live in reference
files in this repository. Read them before anything else.

References: <space-separated ${CLAUDE_SKILL_DIR}/references/<file>.md paths for this group>
Diff:       <diff-path captured while acquiring the diff>
Context:    <context string captured while acquiring the diff>
User hint:  <verbatim scope hint, or "none">

Steps:
  1. Read every listed reference file end-to-end.
  2. Read the diff. Infer the author's intent from the diff itself.
  3. For each modified function a reference flags as in-scope, read source as
     needed (Read / Grep / Glob / Bash). The references include starter greps.
  4. Form findings. Before keeping each one, re-read the cited file at the cited
     line, confirm the EVIDENCE snippet matches current source, and drop it if it
     is pre-existing, clearly intentional, or you cannot point to it in the diff.
  5. Return the surviving findings, in the output format below.

Output format — one block per surviving finding, separated by a blank line,
grouped BUG → WARNING → NIT, no prose between blocks:

  SEVERITY: BUG | WARNING | NIT
  FILE: <relative-path:line>
  CATEGORY: <short kebab-case label from the reference>
  ISSUE: <one-sentence statement of the defect>
  EVIDENCE: <1–3 lines copied verbatim from the diff or source>
  WHY: <one or two sentences — what breaks, when, who notices>
  FIX: <concrete suggested change, ideally a short replacement snippet>

Severity calibration:
  BUG     — would fail tests, corrupt data, or break ABI/contract once exercised.
  WARNING — likely wrong, needs author's eyes; may be intentional.
  NIT     — minor; not load-bearing, but reviewers will ask.

If no findings survive, emit exactly one line instead of blocks:
  NO FINDINGS — exercised: <comma-separated checklist items>

End with one line:
  SUMMARY: <N> BUG, <N> WARNING, <N> NIT — exercised: <checklist items>

Apply throughout:
  - Flag only what you can point to in the diff.
  - Investigate, don't report — verify each candidate against current source
    before keeping it; drop anything pre-existing, intentional, or unprovable.
  - Stay within your references' scope; the other groups cover the surfaces
    yours doesn't. Do not duplicate them.
```

**History agent** — dispatch alongside the groups. Use `subagent_type: "general-purpose"`, `model: "sonnet"`, `effort: "medium"`, `description: "Git history review of <PR/branch>"`. Prompt:

```
You are reviewing a WiredTiger change for historical regressions. Do NOT review
style or correctness — only whether this change conflicts with prior decisions
in git history.

Diff:    <diff-path captured while acquiring the diff>
Context: <context string captured while acquiring the diff>

Steps:
  1. Extract the modified file paths from the diff.
  2. For each file: git log --oneline -20 -- <file>
  3. For each changed symbol: git log --oneline -10 -S <symbol> -- <file>
  4. For non-obvious deletions or guard removals:
       git blame -L <start>,<end> <file>  on the original lines.
  5. If a commit message, blame, or adjacent comment shows this change undoes a
     deliberate prior decision, flag it.

Output: same finding blocks as the group reviewers. CATEGORY one of:
regression, undoes-fix, invariant-removal, prior-comment-conflict.
If nothing: NO FINDINGS — exercised: git-log, git-blame
End: SUMMARY: <N> BUG, <N> WARNING, <N> NIT — exercised: git-log, git-blame
```

Wait for all reviewers and the history agent. Collect their returned findings verbatim. Do not re-review or rewrite them. Preserve `file:line` citations exactly.

## Summarising for the user

The reviewers return findings as tagged blocks (`SEVERITY` / `FILE` / `CATEGORY` / `ISSUE` / `EVIDENCE` / `WHY` / `FIX`). Those tags are a transport format — never paste the raw uppercase blocks at the user. Render them into a clean markdown report.

Open with a one-line headline: the review target and the totals, in plain sentence case — e.g. `Review of PR #1234 — 1 bug, 2 warnings, 1 nit.`

Then one section per severity that has findings, ordered **Bugs → Warnings → Nits** (a `###` heading in Title Case; omit a severity with no findings). Within a section, order by file. Render each finding as a short markdown entry:

- A bold header line: **`file:line`** — <issue, as a short sentence>.
- An italic metadata line: the category and the reviewer that found it, e.g. *disagg-api-usability · found by interfaces reviewer*. You know each finding's source because you collected it from that group's (or the history agent's) returned output — attribute every finding to its reviewer.
- One or two sentences of plain prose explaining why it matters.
- The evidence as a fenced `c` code block, copied verbatim.
- A final line beginning `Fix:` with the concrete suggested change.

Preserve `file:line` and evidence verbatim — do not paraphrase. Keep the issue line short and let the prose and fix lines carry the detail.

After the findings, a short **Coverage** list: one line per group and the history agent naming the surface and whether it had surviving findings (`correctness: clean`, `concurrency: 1 warning`, …). This tells the user what was checked.

Close with a one-line **Next steps** — typically `cd dist && ./s_fast` if hygiene flagged style, plus any reviewer-recommended follow-up.

## Guardrails

- **Stay diff-free.** You acquire the diff and route its path; you never read it into your own context or form opinions about C source. If you are reading a barrier or a lock, you are doing a reviewer's job — stop.
- **Reviewers load knowledge from the reference files.** Do not paste macro semantics into dispatch prompts; the prompt already names the files to read.
- **Preserve citations.** Quote findings verbatim between steps; never paraphrase.
- **Stop on empty diff.** If `git diff` returns nothing, ask before dispatching.
- **Respect user scope hints.** "just style" or "only the disagg parts" → dispatch only the matching group(s).
- **One diff snapshot per review.** Capture it while acquiring the diff, reuse the path everywhere.
