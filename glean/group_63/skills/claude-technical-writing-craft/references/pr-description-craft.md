<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `pr-description-craft` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: pr-description-craft
description: Author pull-request descriptions that get merged fast — the "what / why / how / test" template, GitHub/GitLab/Bitbucket PR templates, the "TL;DR for skimmers + details for nit-pickers" two-tier pattern, before/after screenshots and recordings discipline, when to link out vs inline detail, reviewer-shopping (assigning the right CODEOWNERS), checklists that actually get verified, draft-vs-ready signaling, and Conventional Comments for review threads. TRIGGER when user asks "write a PR description", "draft a pull request", "review my PR body", "what should I put in this PR", "PR template", "how do I structure this PR", "before/after screenshots", "request review from", "the PR body is too long", "merge this PR", "split this PR", "stacked PRs", or pastes a Git diff/branch and asks for the description. SKIP when authoring the commit messages inside the PR (use commit-message-craft); reviewing the PR diff for correctness or running a PR review loop (use pr-review-loop, code-reviewer, or pull-request-guidelines); rolling up merged PRs into customer-facing release notes (use changelog-and-release-notes or changelogs-for-humans); writing an RFC or design doc that precedes the implementation work (use rfc-and-design-docs); editing the prose inside the description for sentence-level style (use writing-expert or technical-writing-craft).
---

# Pull Request Description Craft

## Overview

**What this skill covers.** Author the *body* of a pull request — the description, checklists, screenshots, "how to test" section, and reviewer assignment — so it gets approved on the first review pass instead of the third. Covers the four-part WWHT template (What / Why / How / Test), the two-tier reader pattern (skimmers above the fold, nit-pickers below), the before/after evidence discipline (screenshots, recordings, perf numbers), the link-out-vs-inline trade-off, and the social-protocol layer (CODEOWNERS, draft signaling, the "this PR is blocked on X" callout).

**When to use.**
- User is about to open a PR and asks for the description drafted from the diff or branch.
- User pastes an existing PR body and asks for critique or rewrite.
- User asks how to structure a PR for fast review.
- User asks about screenshots, recordings, or before/after evidence in a PR.
- User asks how a PR should signal it's blocked, stacked on another PR, or in draft.
- User asks how to handle a "this PR is too big" situation — splitting strategies, stacked PRs.

**When to skip.**
- Authoring commit messages inside the PR (use `commit-message-craft`).
- Customer-facing release notes that aggregate many merged PRs (use `changelog-and-release-notes`).
- Pre-implementation RFCs and design documents (use `rfc-and-design-docs`).
- Sentence-level prose editing of the PR body (use `writing-expert`).
- Writing PR review *comments* themselves rather than the PR body (use Conventional Comments — referenced here but a separate concern).

## Core Concepts

### 1. The WWHT template — What, Why, How, Test

The four sections every PR body needs, in this order:

```markdown
## What
One or two sentences naming the change at a high level. Read it
once and a senior engineer should know what shipped.

## Why
The problem this PR solves. Link to the issue, the incident, the
customer ticket. If there's no preexisting problem statement,
write it here — three sentences max.

## How
The approach. Name the modules touched, the design decision made,
and the alternative considered-and-rejected (one sentence).

## Test
What you ran. What a reviewer should run to verify. Screenshots,
recordings, or perf numbers attached below.
```

The order matters: reviewers read top-down and decide whether to keep reading after each section. **What** earns 5 seconds. **Why** earns the next 30. **How** earns the careful read. **Test** earns the merge.

**Common WWHT variants:**
- GitHub default template: Description / Motivation / Type of change / Checklist.
- Atlassian variant: Context / Summary of changes / Testing instructions / Reviewer notes.
- Google internal: Problem / Solution / Risk / Rollback plan.

All are isomorphic to WWHT. Pick the variant your repo's `.github/PULL_REQUEST_TEMPLATE.md` enforces and stick to it.

### 2. The two-tier reader pattern — skimmers above, nit-pickers below

A PR has two audiences:
- **Skimmers** — the manager, the on-call, the eventual archaeologist running `git log` six months from now. They read the first paragraph and the screenshots.
- **Nit-pickers** — the assigned reviewer. They read every word, click every link, run the test script.

Optimize the top of the body for skimmers; put depth below.

```markdown
## TL;DR
Adds retry-with-backoff to the S3 uploader so 3% of uploads that
fail on NLB resets now succeed. No API change. Behind no flag.

<details>
<summary>Details</summary>

## Why
Incident #4821 (May 12) showed 3% of uploads failing during the
14:00 UTC spike, all on NLB-reset traces ...

## How
Wraps the existing `putObject` call in a `withRetry` helper that
retries on `EPIPE`, `ECONNRESET`, and TLS-handshake errors with
250ms / 750ms backoff. Caps at 2 retries to stay inside the
existing 30s deadline.

## Test
- `npm run test:integration -- s3-uploader.retry`
- Replayed the captured NLB-reset trace; all 200 sessions succeed.

</details>
```

The `<details>` collapsible is supported on GitHub, GitLab, Bitbucket, and Azure DevOps. Skimmers see the TL;DR; nit-pickers expand the rest.

### 3. Before/After evidence — screenshots, recordings, numbers

Every UI change needs a screenshot. Every behavior change needs a number or a recording.

**Screenshot discipline:**
- Two screenshots side by side: `Before` and `After`. Use a markdown table for layout:
  ```markdown
  | Before | After |
  |---|---|
  | ![before](url) | ![after](url) |
  ```
- Crop tight to the changed region. Full-window screenshots waste reviewer attention.
- Annotate with arrows or boxes if the change is subtle. Tools: macOS Preview, Skitch, CleanShot, Excalidraw.
- For dark-mode features, include both light- and dark-mode screenshots.
- For responsive features, include mobile + desktop.

**Recording discipline:**
- Use recordings for interaction flows, animation, or anything that can't be captured in one frame.
- Trim to ≤ 30 seconds. Reviewers don't scrub through 4-minute videos.
- macOS: `Cmd+Shift+5`. Linux: `peek`, `kooha`. Cross-platform: Loom, CleanShot, OBS.
- For terminal demos, prefer asciinema (text-searchable) over a video.
- Embed via drag-drop into the GitHub PR body so it uploads to GitHub's CDN; external links rot.

**Numerical evidence:**
- For performance PRs: before/after benchmark numbers, with units and confidence intervals.
- For bundle-size PRs: before/after KB.
- For test-coverage PRs: before/after % from the coverage tool.
- Cite the tool that produced the number ("`hyperfine --warmup 3 ./bench`", "`vitest --coverage`").

### 4. The "how to test" section — concrete, copy-pasteable

The most-skipped, most-valuable section. Every step must be executable.

**Bad:**
```
Run the tests and make sure they pass.
```

**Good:**
```markdown
## How to test
1. `git checkout this-branch && npm install`
2. `npm run test:unit -- s3-uploader.retry`
3. In Chrome with the extension loaded, click an attachment in any
   case page. Confirm the upload succeeds without a console error.
4. Open DevTools → Network. Throttle to "Slow 3G". Re-upload.
   Confirm two retry attempts in the network log before success.
```

Rules:
- Number the steps. Reviewers check off as they go.
- Include the commands verbatim. No "run the appropriate test".
- Include the expected observation. "Confirm X" not just "test it".
- If the test requires setup (env vars, fixtures, test accounts), say so explicitly.

### 5. Link-out vs inline detail

The decision rule: if the detail is necessary to *decide* the review, inline it. If it's necessary only to *verify* a specific claim, link out.

**Inline:**
- The rationale for the design decision (one paragraph).
- The error message or log line that motivated the fix.
- The before/after screenshots.
- The 3-line code snippet that's central to the change.

**Link out:**
- The full incident report.
- The original RFC.
- The customer ticket (Jira, Zendesk).
- The benchmark methodology.
- The migration guide.
- Stack traces longer than 10 lines (use a gist).

Always use markdown link syntax `[text](url)` — bare URLs are noise and may not be clickable in some viewers.

### 6. Reviewer shopping — picking the right reviewer, not the available one

A PR sits open because the wrong person was asked to review.

**Rules:**
- Use the repo's `CODEOWNERS` file as the default — those reviewers are auto-assigned and tracked.
- For cross-cutting changes, request *one* primary reviewer from each affected area; don't ping the whole team.
- For domain-specific changes (security, perf, accessibility, i18n), tag the domain owner even if they're not CODEOWNERS.
- "I'll just merge it myself" is okay for typos and version bumps; never okay for behavior changes in shared code.

**Avoid these review-shopping anti-patterns:**
- Tagging 8 reviewers to "see who responds first". This creates a bystander effect; nobody owns it.
- Tagging only your friendly reviewer who'll rubber-stamp. The review is a quality gate, not a social ritual.
- Re-requesting review without doing anything (the GitHub "Re-request review" button). Push at least a comment explaining what changed.

### 7. Draft, blocked, and stacked signaling

A PR isn't always "ready to merge as-is". Signal the state in the title and body.

**Title prefixes (de facto standard):**
- `[Draft]` or use GitHub's native "Draft PR" toggle — reviewers know not to spend cycles yet.
- `[WIP]` — older convention; "Draft" PR is the modern replacement.
- `[Do not merge]` — for PRs that exist for discussion, demonstration, or CI-validation only.
- `[RFC]` — proposal stage, even though code is written.
- `[Stacked on #1234]` — depends on another PR; merging this one first will produce a wrong diff.

**Body callouts:**
```markdown
> [!IMPORTANT]
> Stacked on #1234. Review #1234 first; this PR's diff will look
> wrong against `main` until #1234 merges.

> [!WARNING]
> Do not merge until the security review in #1235 is approved.
```

GitHub renders `[!NOTE]`, `[!TIP]`, `[!IMPORTANT]`, `[!WARNING]`, `[!CAUTION]` as styled callouts.

### 8. Stacked PRs — when one PR is too big

When a feature can't reasonably fit in one reviewable PR (~400 lines is the soft ceiling per most-cited research), stack:

1. **PR #1: refactor / scaffolding** — no behavior change, sets up types and module structure.
2. **PR #2: the substantive change** — depends on #1.
3. **PR #3: tests / docs / migration** — depends on #2.

Each PR is independently reviewable and mergeable. Tools that help: Graphite, Sapling, `git-branchless`, `git-spice`. Plain Git works too — use `git rebase --onto` to re-parent later PRs after earlier ones merge.

Mark each stacked PR's body with its position:
```markdown
This is PR 2 of 3 in the auth-overhaul stack:
- #1234 — refactor auth module (merged)
- **#1235 — add OIDC flow** ← you are here
- #1236 — migration & docs
```

### 9. Checklists that get verified, not skipped

The Graphite study found that vague checklist items are checked unconditionally; specific items are actually checked. Rule: every checkbox should be verifiable in under 30 seconds.

**Bad:**
- [ ] Code follows best practices
- [ ] Tests are good
- [ ] Documentation is updated

**Good:**
- [ ] `npm run lint` passes locally
- [ ] `npm run test` passes locally
- [ ] Manifest version bumped in `manifest.json` and `package.json`
- [ ] Screenshot added for any UI change
- [ ] CHANGELOG.md entry added under `## Unreleased`

Put the checklist *at the bottom*, after the WWHT body. Reviewers scan the body for context, then run the checklist as the gate.

### 10. The PR title — Conventional Commits applies here too

PR titles get squash-merged into the default branch's history in many repos. Treat them as commit subjects:

- Imperative mood: `Add retry logic` not `Added retry logic` or `Adding retry logic`.
- Conventional-Commits prefix if your repo uses it: `feat(api): add /v3/users endpoint`.
- ≤ 72 chars (GitHub's PR list truncates ~70).
- No issue numbers in the title — they go in the body or in trailers; GitHub PR titles auto-link `#1234` only inconsistently across UI surfaces.

If you're using squash-merge, the PR title *is* the resulting commit's subject. If you're using merge-commits, the commit messages inside the PR carry the story (see `commit-message-craft`).

## Templates and Examples

### Template — GitHub PR template (drop into `.github/PULL_REQUEST_TEMPLATE.md`)

```markdown
## TL;DR
<!-- One sentence. What changed and what now works that didn't. -->

## Why
<!-- Link the issue / incident / ticket. Restate the problem in 1-3 sentences. -->
Closes #

## How
<!-- The approach. Modules touched. One sentence on the rejected alternative. -->

## How to test
<!-- Numbered, copy-pasteable steps with expected observations. -->
1.
2.
3.

## Screenshots / recordings
<!-- Before / After for UI. Trim recordings to ≤ 30s. -->

## Checklist
- [ ] Tests pass locally (`npm test`)
- [ ] Lint passes (`npm run lint`)
- [ ] Manifest/package version bumped
- [ ] CHANGELOG entry added (if user-facing)
- [ ] Screenshots attached (if UI changed)
- [ ] Breaking change noted in body (if any)
```

### Template — minimal PR body for a bug fix

```markdown
## TL;DR
Fixes a 3% upload failure rate caused by NLB resets during peak load.

## Why
Incident #4821: 14:00 UTC May 12 spike caused ~3% of S3 uploads to
fail with `ECONNRESET` before any bytes left the client. Captured
NLB traces in the incident postmortem.

## How
Wraps `putObject` in `withRetry` (250ms, 750ms backoff, 2 retries
max). Retries only on `EPIPE`, `ECONNRESET`, and TLS-handshake
errors — *not* on 5xx (those should fail loud).

Considered moving to multipart upload to dodge the issue entirely;
rejected because the change surface is bigger and the symptom is
specifically pre-flight.

## How to test
1. `npm run test:unit -- s3-uploader.retry`
2. With the NLB-reset replay fixture: `npm run test:integration -- s3-uploader.replay`. Expect 200/200 successes; previously 6/200 failed.
3. In the live extension against the test bucket, upload a file; confirm one network call, success.

## Screenshots
N/A — no UI change.

Closes #4821
```

### Template — feature PR with UI change

```markdown
## TL;DR
Adds an "ownerless case" acknowledgement popup that lets the on-call
TAM claim a case in one click from any Hub page.

## Why
Ownerless cases currently sit in the queue until someone manually
filters for them. Two SLA breaches in April traced to this gap
(tickets #4001, #4055).

## How
- New module: `src/alerts/ownerless-case-alert-view.js` (pure helper, unit-tested).
- New popup wired through `MCA_OWNERLESS_ALERT` in the service worker.
- Renders in the existing overlay shell — no new Chrome permissions needed.

Considered making it a chrome notification; rejected because dismissed notifications don't re-appear and TAMs reported missing them.

## How to test
1. `git checkout this-branch && npm install`
2. Load unpacked at `chrome://extensions`. Open any Hub case page.
3. From a separate TAM account, leave a case ownerless. Within 60s,
   the popup should appear in your Hub view.
4. Click "Claim". The case should be assigned to you and the popup
   should dismiss.

## Screenshots
| Before | After |
|---|---|
| ![queue without alert](.../before.png) | ![popup visible on Hub](.../after.png) |

## Recording
[ownerless-alert.mp4](https://github.com/user-attachments/.../recording.mp4) (22s)

## Checklist
- [x] Tests pass (`npm test` — 47/47)
- [x] Lint clean
- [x] Version bumped: manifest.json 1.0.178 → 1.0.179
- [x] CHANGELOG entry under `## Unreleased`
- [x] Screenshots attached
- [x] Not a breaking change

Closes #4055
```

### Template — stacked PR

```markdown
This is PR 2 of 3 in the auth-overhaul stack:
- #1234 — refactor auth module (merged)
- **#1235 — add OIDC flow** ← you are here
- #1236 — migration & docs (depends on #1235)

## TL;DR
Adds the OIDC sign-in flow on top of the refactored auth module from #1234.

## Why
...
```

### Template — breaking-change PR

```markdown
## TL;DR
**BREAKING**: Removes the deprecated `/v1/users` endpoint. Six-month
deprecation window closed 2026-04-01.

## Why
/v1 was sunset on 2026-04-01 (announced in v3.2.0 release notes,
Deprecation: header shipped 2025-10-01). Two clients still on /v1
have been emailed and confirmed migration.

## How
- Removes the `/v1/*` route handlers and their tests.
- Adds 410 Gone responses for any residual /v1 traffic.
- Updates the migration guide with a "you should already be off /v1"
  callout.

## Migration impact
- Clients on `/v2` (95% of traffic in May): no impact.
- Clients on `/v1` (last 5% — TenantA, TenantB): pre-coordinated.
- Internal services: confirmed all on `/v2` per service-registry audit.

## Rollback plan
Revert this commit. The /v1 handlers are deleted, not feature-flagged,
so a revert is the only rollback.

## How to test
1. `curl -i https://api.example.com/v1/users` → expect `410 Gone`.
2. `curl -i https://api.example.com/v2/users` → unchanged.
3. Integration suite: `npm run test:integration -- api/legacy-removal`.

Closes #2901
```

## Anti-Patterns

- **Empty PR descriptions** ("see commits" or worse, nothing) — guarantees a slow review. The reviewer has to reverse-engineer intent from the diff.
- **Bullet list of every changed file** — the diff already shows this. Bodies are for the *why*, not the file inventory.
- **Screenshots without before/after** — the reviewer can't tell if the screenshot shows the bug or the fix.
- **"Tested locally"** with no steps — unverifiable. Either spell out the steps or say "no manual testing — covered by added unit tests".
- **40-item checklist of every conceivable concern** — gets check-the-box-checked unconditionally. Trim to 5-7 items that are actually verifiable.
- **Tagging 8 reviewers** — bystander effect. Tag one primary plus CODEOWNERS.
- **Massive PRs (1,000+ lines)** — review quality drops sharply past ~400 lines. Split into a stack.
- **PR title that doesn't match the squash-commit convention** — if you use squash-merge, the title becomes the commit message. `update stuff` becomes `update stuff` in `main`'s history forever.
- **Hiding the breaking-change disclosure in the body** — put `**BREAKING**` in the TL;DR. Reviewers should not have to scroll to discover the contract change.
- **"Will fix in follow-up PR"** with no link — name the follow-up issue or PR number. Otherwise it disappears.

## Decision Heuristics

- **Inline detail or link out?** If the detail is necessary to *decide* whether the change is correct → inline. If it's necessary only to *verify* a specific claim → link out.
- **One PR or split into a stack?** Soft ceiling ~400 lines of meaningful diff (excluding generated code, snapshots, lockfile). Past that, the review degrades. Split.
- **Draft or ready for review?** If a reviewer's first comment is likely to be "did you mean to push this yet?", it's a draft. If you'd be embarrassed by some specific thing, fix it before requesting review.
- **Screenshot or recording?** Static change → screenshot. Interaction, animation, or sequence → recording. If unsure, recording (you can always frame-grab a still later).
- **Re-request review or just push?** If the change is substantive, re-request explicitly with a comment summarizing what changed since the last review pass. Re-requesting without a summary signals "this is the same PR" and risks an unchanged ⌘+enter approval.
- **Squash-merge or merge-commit?** Squash for PRs whose commits are scaffolding (3 commits of "fix typo", "more tests", "fix the fix"). Merge-commit (or rebase-merge) for PRs whose commits are deliberately staged (refactor → feature → tests).
- **Use the repo template or write fresh?** Always use the repo template if one exists. If it omits a section your change needs, add the section *below* the template, don't replace the template.

## References

- [Graphite — Best practices for GitHub pull request descriptions](https://graphite.com/guides/github-pr-description-best-practices) — WWHT structure, screenshot discipline, link-out rules.
- [Graphite — Comprehensive Checklist: GitHub PR Template](https://graphite.com/guides/comprehensive-checklist-github-pr-template) — checklist anti-patterns and verifiability rules.
- [GitHub Docs — Creating a pull request template for your repository](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/creating-a-pull-request-template-for-your-repository) — `.github/PULL_REQUEST_TEMPLATE.md` setup.
- [Conventional Comments specification](https://conventionalcomments.org/) — companion convention for *review-comment* labels (nit, suggestion, issue, praise, question, thought) — referenced for review-thread discipline.
- [Everhour — GitHub PR Template Examples & Setup Guide](https://everhour.com/blog/github-pr-template/) — template variants across stacks.
- [DeployHQ — Pull Request Best Practices: A Complete Guide](https://www.deployhq.com/blog/the-perfect-pull-request-best-practices-for-collaborative-development) — review-shopping anti-patterns and stacked-PR guidance.
- [Dev.to — GitHub PR Review: Best Practices and Tools (2026)](https://dev.to/rahulxsingh/github-pr-review-best-practices-and-tools-2026-1p90) — 2026-current review-flow norms.
- [GitHub Docs — Alerts in markdown (`[!NOTE]`, `[!WARNING]`)](https://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/basic-writing-and-formatting-syntax#alerts) — callout syntax for blocked/important signals.
