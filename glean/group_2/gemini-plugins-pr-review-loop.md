# pr-review-loop

**Category:** Frontend & Web Development
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/pr-review-loop/skills/pr-review-loop

## Description
Use when monitoring a PR for CI failures or reviewer feedback, responding to review comments, or babysitting a build.

---

# PR Review Loop

Single-cycle skill that monitors Evergreen CI and PR reviews for the current branch. Each invocation performs one pass and exits — use with `/loop 10m /pr-review-loop` for continuous monitoring.

> Before running, check `references/` for your repo-specific configuration. Load the file matching your repo (e.g., `references/mms.md`, `references/server.md`) to get the correct GitHub repo, build variants, test commands, and ticket prefix. To add support for a new repo, copy [references/_template.md](references/_template.md).

## Modes

| Mode | Trigger | Behavior |
|---|---|---|
| **Interactive** (default) | `/pr-review-loop` or `/pr-review-loop interactive` | Pauses for user approval before **pushing code** or **commenting on the PR** |
| **Autonomous** | `/pr-review-loop autonomous` | Fully autonomous — no approval gates |

Interactive mode performs all analysis, diagnosis, and local code fixes automatically. It only pauses at two external-action gates:

1. **Before implementing review fixes** (Phase 5) — presents proposed action plan for approval
2. **Before pushing code and posting replies** (Phase 6) — presents commit summary and proposed replies for approval

Both modes work with `/loop`. In interactive mode, the loop pauses at approval gates and waits for user input — the loop timer resets after the user responds, so no iterations are lost.

**Tip**: For fully hands-off monitoring, use `/loop 10m /pr-review-loop autonomous`. This will monitor CI, fix failures, address review comments, and push — all without interruption.

### Mode detection

| User says | Mode |
|---|---|
| "auto", "full auto", "hands off", "autonomous", "don't ask" | **Autonomous** |
| Everything else (default) | **Interactive** |

---

## Execution Model

> **Workers skip this section.** If you arrived here from `references/worker.md`, go directly to [Phase 0](#phase-0-state-check).

**Determine mode first** (using Mode detection table above), then:

### Autonomous mode — delegate to Worker agent

The skill argument ("autonomous", "full auto", etc.) is **only for mode detection** — do NOT use it as the branch name.

1. Run `git branch --show-current` to get the current branch name (e.g. "CLOUDP-402360"). Store the result.
2. Build a **context briefing** for the Worker — a short block of text (target: under 500 tokens) that gives the Worker enough judgment context to handle review comments well. Assemble it from two sources:

   **a. User-provided instructions** — anything the user said when invoking the skill or during the conversation that constrains how review feedback should be handled. Examples: "reject scope expansion", "only fix bugs, don't refactor", "this PR is just a migration, push back on feature requests." Include these verbatim or lightly paraphrased.

   **b. PR context you infer from the conversation** — the PR's purpose, key constraints, what's in/out of scope. Distill this from the conversation history. If the user has been discussing the PR with you, you already know why it exists and what matters. Extract the bits that would help a reviewer-agent make good judgment calls.

   If neither source yields anything (e.g., the user just said `/loop 10m /pr-review-loop autonomous` with no prior conversation), omit the briefing block entirely — the Worker will rely on SKILL.md's default categorization rules.

3. Spawn a Worker agent whose prompt includes the branch name and the context briefing:
   > "Read the pr-review-loop skill's `references/worker.md` and execute all phases for branch ACTUAL_BRANCH_NAME. State file: `tmp/pr-review-loop-state.json`.
   >
   > Context for handling review comments:
   > CONTEXT_BRIEFING"
   >
   > (Replace ACTUAL_BRANCH_NAME with the real branch from step 1. Replace CONTEXT_BRIEFING with the briefing from step 2, or remove the "Context for handling review comments:" section if no briefing was assembled.)

   Allowed tools for the Worker: Read, Edit, Write, Glob, Grep, Bash, Agent, mcp__Evergreen__get_my_patches, mcp__Evergreen__get_parsely_logs_for_patch, mcp__Evergreen__debug_failing_task, mcp__Evergreen__summarize_patch_tasks
4. Print the Worker's output (the iteration report) and exit.

**Do not execute Phases 0–6 yourself in autonomous mode** — the Worker handles everything, including state file management and lock lifecycle.

### Handling user decisions after a FLAGGED iteration (autonomous mode)

If the Worker returned a FLAGGED summary and the user responds with a decision before the next loop iteration fires, record it in the state file immediately so the next Worker picks it up:

1. Parse the comment ID(s) from the FLAGGED summary line (format: `[id:XXXXXXX]`).
2. Read `tmp/pr-review-loop-state.json`.
3. Locate the `flagged_items` field per the schema in [state-file.md](references/state-file.md): if `last_push` is non-null, write into `last_push.flagged_items`; otherwise write into the top-level `flagged_items` field. Do **not** synthesize a placeholder `last_push`.
4. For each comment the user made a decision about, add or update an entry as a JSON object (use a real JSON encoder; do not concatenate strings):
   ```json
   { "comment_id": 3184327788, "decision": "<user's decision verbatim>" }
   ```
5. Write the updated state file.
6. Confirm: "Decision recorded. The next Worker iteration will act on it."

The user does not need to visit GitHub — their reply in this conversation is sufficient. The next Worker reads `flagged_items.decision` and treats it as the instruction for that comment.

**Why this matters**: Each Worker starts with a fresh context (~10–15K tokens) regardless of how large the parent conversation has grown. This prevents context accumulation from compounding across dozens of loop iterations.

### Interactive mode — run directly

Continue to Phases 0–6 below.

---

## State file

This skill tracks progress across invocations via `tmp/pr-review-loop-state.json`. See [state-file.md](references/state-file.md) for the full schema and lifecycle rules.

---

## Phase execution checklist

Every invocation MUST execute these phases in order. Do not skip phases for efficiency.

| Phase | When to run | What it does | Approval gate |
|---|---|---|---|
| **0. State Check** | Always | Read state file, check for concurrent run, write lock | — |
| **1. Identify PR and Patch** | Always | Get branch, PR, HEAD, patch, stale check | — |
| **2. Trigger AI Code Review** | If `BUILD_STALE = false` AND bot is configured | Check if AI review bot needs to review current HEAD | — |
| **3. Check Evergreen Failures** | If `BUILD_STALE = false` | Check for new build failures | — |
| **4. Fix Build Failures** | If Phase 3 found failures | Diagnose and fix failures locally | — |
| **5. Check PR Reviews** | **Always — never skip** | Analyze and address review comments | **Interactive: approve plan before implementing** |
| **6. Commit, Push, Report** | Always | Commit changes (if any), post replies, print status report | **Interactive: approve before push + replies** |

If `BUILD_STALE = true`: skip Phases 2–4, but still run Phases 5 and 6.

---

## Phase 0: State Check

1. Get the current branch name (`git branch --show-current`).
2. Ensure `tmp/` directory exists (`mkdir -p tmp`). Read `tmp/pr-review-loop-state.json`. If the file does not exist, create it with `branch: <current branch>`, `running_since: null`, and `last_push: null`. If the file exists but contains malformed JSON, check the file's mtime (`stat -f %m` on macOS, `stat -c %Y` on Linux):
   - **mtime is more than 30 minutes old**: treat as a crashed write — delete the file and recreate as a fresh state file, then proceed.
   - **mtime is within the last 30 minutes**: assume an in-flight write from a concurrent invocation. Print "State file malformed (possible concurrent write); will retry next invocation" and **STOP** without modifying the file.
3. If the file exists but `branch` does not match the current branch, the state is from a different PR. Delete the file and recreate it as in step 2 (this resets `last_push` so stale data from another branch is not applied).
4. If `running_since` is set and less than 30 minutes old, print "Another invocation is still running" and **STOP** (do not clear the lock — the other invocation owns it).
5. If `running_since` is set but older than 30 minutes, treat it as stale (previous run crashed) and proceed.
6. Write `running_since` with the current UTC timestamp.
7. On **every** exit path from this point forward — including early **STOP** in Phase 1 (e.g., on `master`, no PR found) — clear `running_since` before exiting.

---

## Phase 1: Identify PR and Patch

1. Get the current branch:
   ```bash
   git branch --show-current
   ```
2. If on `master` or `main`, report error and **STOP**.
3. Find the PR for this branch. Try `gh pr view` first; if it fails (common with draft PRs or missing upstream tracking), fall back to `gh pr list`:
   ```bash
   gh pr view --json number,url,headRefName 2>/dev/null \
     || gh pr list --head "$(git branch --show-current)" --state open --json number,url,headRefName --limit 1 --jq '.[0]'
   ```
4. If no PR exists from either method, report error and **STOP**.
5. Get `HEAD_SHA` from the local git repo (do NOT use `headRefOid` from `gh pr view` — it can lag after a push):
   ```bash
   git rev-parse HEAD
   ```
6. Store `PR_NUMBER`, `PR_URL`, and `HEAD_SHA`.
7. Use Evergreen MCP `get_my_patches` (limit: 5) to find the most recent patch for this branch.
   Filter results by matching the PR number in the patch description (e.g., `pull request #160683`).
   Do not assume the first result is for this branch — patches from other branches may appear.
8. Store `PATCH_ID` and `PATCH_CREATE_TIME`.

### Stale build check

Determine if the Evergreen patch covers the current HEAD. The patch's `Git Hash` field is the **master base revision**, not the branch head — do not compare it against `HEAD_SHA`.

Instead, check whether a new commit was pushed after the patch was created:
```bash
git log --format="%H %ci" HEAD -1
```
Compare the latest commit timestamp against `PATCH_CREATE_TIME`. If the commit is newer than the patch, the build is stale.

If the build is stale:
- Set `BUILD_STALE = true` — skip Phases 2–4 (AI review trigger and Evergreen checks).
- **Still proceed to Phase 5** (review comments) — reviews can be addressed while waiting for a new build.
- Note the stale status in the iteration report.

---

## Phase 2: Trigger AI Code Review (if configured)

**Skip this phase if** your repo's `references/` file does not configure an AI review bot (i.e., "Enabled: No" or the section is absent).

> Check `references/` for your repo's AI review bot configuration (bot name and trigger phrase).

1. Fetch PR reviews and issue comments (replace `{REPO}` with your repo from `references/`):
   ```bash
   gh api repos/{REPO}/pulls/{PR_NUMBER}/reviews --paginate
   gh api repos/{REPO}/issues/{PR_NUMBER}/comments --paginate
   ```
2. Find the most recent bot review (any `commit_id`) and the most recent trigger comment.
3. Decide whether to post a new trigger:

   | Condition | Action |
   |---|---|
   | Bot review exists for current `HEAD_SHA` | **Skip** — already reviewed this commit |
   | A trigger comment exists that is newer than the last bot review (or no bot review exists yet) | **Skip** — bot hasn't responded to the last trigger yet; wait |
   | Bot's last review is for an older commit AND no trigger comment is newer than that review | **Post** — new commit needs review |

   Only post if there are actually new commits or new review comments since the last bot review. If the bot already reviewed this HEAD and the only new activity is the loop's own replies, do not re-trigger.

4. If posting (use trigger phrase from `references/`):
   ```bash
   gh pr comment {PR_NUMBER} --body "<trigger phrase>"
   ```

This is a fire-and-forget step — do not wait for the bot to complete.

---

## Phase 3: Check Evergreen for New Failures

Use Evergreen MCP `get_parsely_logs_for_patch` with the `PATCH_ID`.

This tool is fast because it:
- Only checks required build variants (see `references/` for your repo's variant list)
- Returns only failed tasks with Parsely log URLs
- Returns immediately if nothing has failed

### Filtering already-handled failures

Check `last_push.failures_fixed` in the state file. Skip any task IDs that were already fixed in a previous invocation. This is more reliable than time-based filtering and works across all platforms.

### Quick status check

Also check overall patch status:
```bash
evergreen list-patches -n 1 -j --id PATCH_ID | python3 -c "
import json,sys
try:
    d = json.load(sys.stdin)
    p = d[0] if isinstance(d, list) else d
    print(f'status:{p.get(\"status\",\"unknown\")}')
    print(f'finished:{p.get(\"finish_time\") or \"\"}')
except (json.JSONDecodeError, KeyError, IndexError):
    print('status:unknown')
    print('finished:')
"
```

- `BUILD_COMPLETE` = `finish_time` is not None/empty
- `BUILD_SUCCESS` = `status` is `"success"`

### Evaluate state

| Condition | Action |
|---|---|
| `BUILD_SUCCESS` AND no failures | → Go to Phase 5 (check reviews), then Phase 6 (report success — candidate for termination) |
| `BUILD_COMPLETE` AND no *new* failures (but `status != "success"`, e.g. unrelated failures were filtered out) | → Go to Phase 5, then Phase 6 (report complete but not green — do not terminate) |
| New failures exist | → Go to Phase 4, then Phase 5 |
| No new failures AND build still running | → Go to Phase 5 (check reviews while waiting) |

---

## Phase 4: Analyze and Fix Build Failures

Read and follow the instructions in the `evergreen-cicd` skill's SKILL.md and its references for detailed guidance on failure analysis, priority ordering, fix strategies, and local test commands.

### Pre-check: filter unrelated failures and detect regressions

Before deep-analyzing any failure:

1. **Already handled?** Check `last_push.failures_fixed` in the state file. If the task ID is listed, skip it.
2. **Regression from previous fix?** Check `last_push.files_changed` in the state file. If the failure's stack trace or error points to a file this skill modified in the previous push, do NOT auto-fix. Report it as `[POSSIBLE REGRESSION]` in the iteration report for human review.
3. **Unrelated?** Check if the same task is also failing on master. Use the Evergreen MCP tools or check the project's recent patches. If a task fails on master too, skip it.

### Diagnosis and fixing

For each new, related failure (in priority order per the `evergreen-cicd` skill):

1. Follow the `evergreen-cicd` investigation workflow to diagnose the failure.
2. If the failure is caused by the branch being behind the base branch, merge the base branch. Determine the base branch from `gh pr view --json baseRefName` (usually `master`, but could be a release branch). Run `git fetch origin <base> && git merge origin/<base>`.
3. Make minimal, focused changes following repository coding standards.
4. **Validate config changes**: If any Evergreen config file was modified, validate before committing. The config path is repo-specific — check `references/` for the correct path (e.g., `references/mms.md`):
   ```bash
   evergreen evaluate <config-path>   # e.g. ./etc/evergreen.yml for Server
   ```
5. After all fixes are applied, proceed to Phase 5 to check for review comments before committing (batch everything into one push).

---

## Phase 5: Check and Address PR Reviews

This phase **always runs** — even when the build is stale or pending. Reviews can be addressed independently of build status.

### Fetch comments

Fetch inline comments, general comments, and review summaries using the endpoints in the [GitHub API Reference](#github-api-reference). Extract comment ID, body, file/line, author. Group into threads.

### Scoping to new comments

Only address comments that have not yet been replied to by the branch author. Check both:
- `reviews_addressed` in the state file — read from `last_push.reviews_addressed` if `last_push` is non-null, otherwise from the top-level `reviews_addressed` field
- GitHub API for existing replies (belt and suspenders)

This prevents re-processing comments already handled in a previous invocation.

### Resuming flagged items

Check `flagged_items` in the state file — read from `last_push.flagged_items` if `last_push` is non-null, otherwise from the top-level `flagged_items` field. These are comments from a previous iteration that required a human decision. Each entry is either a bare comment ID (no decision yet) or an object with a `decision` field:

```json
3184327788
{ "comment_id": 3184327788, "decision": "go ahead with the refactor" }
```

For each flagged item:

1. If a `decision` field is present, use it as the instruction — re-categorize as **Fix** or **Reply** and act on it. Remove from `flagged_items`.
2. If no `decision`, fetch the comment thread from GitHub. If the PR author has replied since flagged, treat the reply as the decision and act on it. Remove from `flagged_items`. If 404 (deleted), remove silently.
3. If neither decision nor reply exists, leave it flagged — re-surface in the iteration report.

### Categorize each comment

| Comment form | Action |
|---|---|
| Reviewer provides explicit replacement code (code block, GitHub suggestion, or inline diff) | **Fix** — apply their code |
| Points out a concrete bug (NPE, wrong return value, missing null check, off-by-one) | **Fix** |
| Nits: rename, formatting, import order, typo fix | **Fix** — group similar nits together, reply "Fixed" |
| Questions a design choice or proposes an alternative approach | **Reply** explaining current rationale, flag for user attention |
| Requests a refactor, new abstraction, or significant restructuring | **Reply** suggesting a follow-up ticket, do not change code |
| Ambiguous, or you cannot determine what change is being requested | **Reply** asking for clarification, do not change code |
| Two reviewers disagree, or comment requires an architectural decision you cannot make confidently | **Interactive**: ask the user to decide. **Autonomous**: **FLAGGED** — do not touch code, do not reply; include `[id:XXXXXXX]` in the iteration report |

### Analyze and present feedback

For each comment, build a structured analysis:

```
### Comment N: @reviewer on File.java:45
**Feedback**: [quoted comment]
**Category**: Fix / Reply / Clarification needed
**Proposed action**: [what will be changed or replied]
**Files affected**: [list]
```

**Interactive mode — Gate 1 (approve the plan)**: Present the full analysis and wait for approval before implementing fixes. The user may edit, skip, or approve individual items.

**Autonomous mode**: Proceed directly to implementation based on the categorization above.

### Implement fixes

For actionable comments: make minimal, focused code changes following AGENTS.md standards. Update tests if the fix affects behavior. Group nit comments (typos, formatting, import order) and batch-fix them together.

### Validate changes

After implementing fixes, run relevant tests to verify nothing is broken. Check `references/` for your repo's local test commands and code quality checks.

If tests fail, analyze and fix before proceeding.

### Post replies

Prefix every reply with `*[AI-assisted response]*` to disclose that the response was generated by an agent posting under the user's account.

Reply format:
- Fixes: `*[AI-assisted response]* Addressed. [description of change]`
- Design questions: `*[AI-assisted response]* [explanation of rationale]`
- Unclear: `*[AI-assisted response]* Could you clarify [specific question]?`

Do **NOT** mark conversations as resolved — that is the comment author's responsibility.

Replies are bundled into the Phase 6 approval gate — do not post them separately.

### Conflicting or complex feedback

- **Conflicting reviewers**: **Interactive**: present pros/cons, recommend an approach, ask the user to decide. **Autonomous**: FLAGGED — record comment IDs, do not touch code or reply.
- **Large refactoring requests**: Reply suggesting a follow-up ticket. Do not implement in this PR.
- **Architectural concerns**: **Interactive**: explain current approach and propose alternatives, flag for user attention. **Autonomous**: if the concern implies a decision the user must make, FLAGGED instead.

---

## Phase 6: Commit, Push, and Report

### If any changes were made (build fixes or review responses):

1. Stage changes carefully — do NOT use `git add -A` (it can stage unrelated/generated files):
   ```bash
   git status --short
   ```
   Review the output and stage only the files you intentionally modified:
   ```bash
   git add <file1> <file2> ...
   ```

2. **Interactive mode — Gate 2 (approve push + replies)**: Present a summary before committing:
   ```
   ### Ready to push
   **Files changed**: [list]
   **Build failures fixed**: [list or "none"]
   **Review comments addressed**: [list or "none"]
   **Commit message**: "<TICKET-PREFIX>: [description]"

   **Proposed replies**:
   1. @reviewer on File.java:45: "*[AI-assisted response]* Addressed. [description]"
   2. @reviewer on Bar.java:100: "*[AI-assisted response]* [explanation]"

   Approve commit, push, and replies?
   ```
   Wait for user approval. The user may edit replies, change the commit message, or exclude files.

   **Autonomous mode**: Proceed directly.

3. Commit with a descriptive message following the repo convention (see `references/` for ticket prefix format):
   ```bash
   git commit -m "<TICKET-PREFIX>: Fix CI failures and address review feedback

   - [list of what was fixed]"
   ```
   Use the ticket number from the branch name if available.

4. Push. If the user uses Graphite (`gt` CLI is available on `$PATH`), use `gt stack submit` instead of `git push`:
   ```bash
   if command -v gt &>/dev/null && gt branch info 2>/dev/null; then
     gt stack submit
   else
     git push
   fi
   ```

5. **Optional: Submit a targeted Evergreen patch for faster iteration.** If build failures were fixed in this cycle, submitting a scoped patch avoids waiting for a full rebuild:

   a. Infer the project ID from the patch's `project_identifier` field (returned by `get_my_patches`). If unknown, run `evergreen list --projects` or `evergreen list --projects | grep <keyword>`.

   b. Submit for the specific failing variant(s). See `evergreen-cicd` skill for flag guidance, especially the `-v`/`-t` vs `--rv`/`--rt` distinction:
      ```bash
      evergreen patch -p <project-id> -v <variant-name> -t all -d "<description>" -f -y
      ```

   c. Comment the new patch URL on the PR (infer repo from `git remote get-url origin`):
      ```bash
      gh pr comment PR_NUMBER --body "*[AI-assisted response]* Submitted targeted Evergreen patch for failing variants: <patch-url>"
      ```

   d. Update `PATCH_ID` in the state file with the new patch ID.

6. Update the state file (`tmp/pr-review-loop-state.json`) with:
   - `last_push.sha`: the new HEAD after push
   - `last_push.files_changed`: list of files modified in this commit
   - `last_push.failures_fixed`: list of Evergreen task IDs fixed in this invocation
   - `last_push.reviews_addressed`: list of GitHub comment IDs addressed this invocation (merge with prior list)
   - `last_push.flagged_items`: list of unresolved flagged entries (merge with prior, remove resolved)
   - `last_push.timestamp`: current time

### If no changes were pushed

The `last_push` object describes the most recent commit pushed by this skill. When no push occurs, do not clear or overwrite it:

- **If `last_push` is null** (no push has ever happened): leave it null. Track `reviews_addressed`, `flagged_items`, and `timestamp` at the top level of the state file instead — do not synthesize a placeholder `last_push`.
- **If `last_push` already exists** (from a previous iteration's push): leave `sha`, `files_changed`, and `failures_fixed` unchanged. Only update `reviews_addressed`, `flagged_items`, and `timestamp` within `last_push`.

In both cases, merge `reviews_addressed` and `flagged_items` with prior values rather than overwriting.

### Iteration report

Always print a status report before exiting, whether or not changes were made:

```
## PR Review Loop — Iteration (TIME)

**PR**: #PR_NUMBER (PR_URL) | **Branch**: BRANCH | **Mode**: interactive/autonomous
**Patch**: PATCH_ID | **Build Status**: started|failed|created | **Complete**: yes/no

### Build Failures:
- [FIXED]: TASK_NAME - description of fix
- [UNRELATED]: TASK_NAME - also failing on master
- [POSSIBLE REGRESSION]: TASK_NAME - failure in file modified by previous auto-fix, needs human review

### Review Comments:
- [ADDRESSED]: @reviewer on File.java:45 - description
- [REPLIED]: @reviewer on Bar.java:100 - explained why not changed
- [FLAGGED]: @reviewer on Foo.java:45 [id:XXXXXXX] - needs your decision (reason)

### Actions Taken:
- Fixed N build failures, addressed N review comments
- Committed and pushed: "commit message"

### Next:
- Build is still running, will check again on next /loop invocation
```

### If no changes were made:

Still check the termination conditions below — if `BUILD_SUCCESS` and all reviews are addressed, delete the state file and report termination even though no code changes were made this iteration.

```
## PR Review Loop — Iteration (TIME)

**PR**: #PR_NUMBER (PR_URL) | **Branch**: BRANCH | **Mode**: interactive/autonomous
**Patch**: PATCH_ID | **Build Status**: STATUS | **Complete**: yes/no

No new failures or review comments to address.
Waiting for build to complete / new reviews to arrive.
```

---

## Termination Conditions

The workflow reports success and suggests stopping `/loop` when:

- **Green build**: Evergreen patch `status` is `"success"` (not just complete with no *new* failures)
- **All clear**: Build is green (`status == "success"`) and all review comments are addressed

On termination, delete the state file (`tmp/pr-review-loop-state.json`) entirely.

The workflow continues looping (via `/loop`) when:

- Build is still running
- Waiting for AI review bot results
- Waiting for human review feedback
- Just pushed fixes and waiting for new build

The workflow reports a problem and suggests user intervention when:

- A failure is too complex to fix automatically
- A review comment requires architectural decisions
- `git push` fails (e.g., remote has new changes)
- Evergreen API auth failure
- **10 fix iterations reached without a green build** — stop, report all remaining failures and what was attempted, and ask the user to take over

---

## GitHub API Reference

Use `gh api` for all GitHub operations. Replace `{REPO}` with your repo's owner/repo from `references/`:

| Operation | Endpoint |
|---|---|
| PR details | `GET /repos/{REPO}/pulls/{PR_NUMBER}` |
| Inline review comments | `GET /repos/{REPO}/pulls/{PR_NUMBER}/comments` (`per_page=100`) |
| General comments | `GET /repos/{REPO}/issues/{PR_NUMBER}/comments` (`per_page=100`) |
| Review summaries | `GET /repos/{REPO}/pulls/{PR_NUMBER}/reviews` |
| Reply to inline comment | `POST /repos/{REPO}/pulls/{PR_NUMBER}/comments/{comment_id}/replies` |
| Post general comment | `POST /repos/{REPO}/issues/{PR_NUMBER}/comments` |

---

## Error Handling

- **GitHub API**: Rate limiting (wait/retry), permission denied (verify auth), PR not found (check number)
- **Evergreen**: Auth failure (check token), task not found (may have been restarted)
- **Implementation**: Test failures (analyze/fix/rerun), build errors (fix), merge conflicts (notify user)
- **Interactive mode**: If user denies an action, skip that item and continue to the next phase

---

## Key Principles

1. **Interactive by default** — pause before pushing code or commenting on PRs unless autonomous mode is explicitly requested
2. **One push per invocation** — batch all fixes (build + review) into a single commit to minimize Evergreen build resets
3. **Build failures take priority** — fix CI first, address reviews during idle time
4. **State-based de-duplication** — use the state file to track handled failures and comments across invocations, and to detect regressions from previous auto-fixes
5. **Master comparison** — skip failures that also fail on master (unrelated to this branch)
6. **Fire-and-forget bot trigger** — trigger the review and move on; check results next invocation
7. **Leverage existing skills** — read and follow `evergreen-cicd` skill and its references for CI diagnosis guidance
8. **Report at every iteration** — always print status before exiting, even if no action was taken
9. **AI-assisted disclosure** — prefix all PR comments with `*[AI-assisted response]*`
10. **Do not resolve conversations** — that is the comment author's responsibility
11. **Clean up the lock you own** — if this invocation acquired `running_since`, clear it on every exit path. On termination, delete the entire state file. Never clear a lock owned by another invocation.

---

## Example Workflow (Interactive Mode)

> **Note**: This example uses MMS conventions (CLOUDP ticket prefix). Adjust based on your repo's `references/` file.

```
User: "Address the review comments on my PR"

Agent: [Fetches PR #12345, checks Evergreen, fetches review comments]

## PR Review Loop — Iteration (14:30 UTC)

**PR**: #12345 | **Branch**: CLOUDP-123456 | **Mode**: interactive
**Build Status**: success | **Complete**: yes

### Review Comments (2 new):

### Comment 1: @reviewer1 on UserService.java:45
**Feedback**: "Add null check before calling processUser()"
**Category**: Fix
**Proposed action**: Add Optional.ofNullable guard
**Files affected**: UserService.java

### Comment 2: @reviewer2 on AuthService.java:100
**Feedback**: "Should we extract this validation into a shared utility?"
**Category**: Reply (design question)
**Proposed action**: Explain current approach, suggest follow-up ticket

**Approve these actions?**

User: "Yes"

Agent: [Implements null check, runs tests, prepares replies]

### Ready to push
**Files changed**: UserService.java, UserServiceTest.java
**Review comments addressed**: 2
**Commit message**: "CLOUDP-123456: Add null check and respond to review feedback"

**Proposed replies**:
1. @reviewer1: "*[AI-assisted response]* Addressed. Added Optional.ofNullable guard."
2. @reviewer2: "*[AI-assisted response]* Good suggestion. The current scope is
   narrow enough that a utility isn't warranted yet — created CLOUDP-123457
   to revisit if more callers emerge."

**Approve commit, push, and replies?**

User: "Yes"

Agent: [Commits, pushes, posts replies]
"Changes pushed. 2 review comments addressed. PR ready for re-review."
```