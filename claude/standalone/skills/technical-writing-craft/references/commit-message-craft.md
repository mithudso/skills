<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `commit-message-craft` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: commit-message-craft
description: Author git commit messages that survive `git blame` years from now — Conventional Commits beyond the basics, Tim Pope's 50/72 rule, "why not what" body discipline, imperative-mood subjects, breaking-change footers (`BREAKING CHANGE:` and `!`), multi-commit storytelling for PR reviewers, fixup/squash/autosquash workflow hygiene, Signed-off-by/DCO trailers, and co-author/reviewed-by trailer conventions. TRIGGER when user asks "write a commit message", "review my commit", "format this commit", "break this into commits", "should this be a breaking change", "what type prefix", "feat vs fix vs chore", "how do I sign off", "squash these commits", "fixup", "autosquash", "conventional commits", "BREAKING CHANGE footer", "split this commit", "rewrite commit history", or pastes a `git log` for critique. SKIP when the task is prose writing not commit messages (use writing-expert); changelog or release-note generation rolling up many commits (use changelog-and-release-notes or changelogs-for-humans); release-blog or launch-narrative authoring (use release-blog-and-launch-narrative); PR description authoring (use pr-description-craft); RFC or design-doc drafting (use rfc-and-design-docs); general technical writing style (use technical-writing-craft); error-message UX in code (use error-message-craft).
---

# Commit Message Craft

## Overview

**What this skill covers.** Author commit messages that survive the test of `git log --oneline` and `git blame` years from now. Goes beyond "write good commits" platitudes into the mechanical rules (50/72, imperative mood, body wrap), the structural rules (Conventional Commits type/scope/subject/body/footer), the discipline rules (one commit = one logical change, fixup before push, squash before merge), and the cross-team rules (BREAKING CHANGE signaling, DCO sign-off, Co-authored-by trailers).

**When to use.**
- The user asks you to write, review, fix, or split a commit message.
- The user pastes a `git log`, `git show`, or `git diff` and asks for commit-message feedback.
- The user is about to run `git commit` and wants the message drafted from staged changes.
- The user is preparing a feature branch for merge and needs to clean up history (`git rebase -i`, `--autosquash`, squash-merge).
- The user asks how to mark a breaking change, add `Signed-off-by`, or attach `Co-authored-by`.

**When to skip.**
- Prose writing (emails, docs, reports) — use `writing-expert`.
- Rolling up many commits into user-facing release notes — use `changelog-and-release-notes`.
- Authoring the PR body itself — use `pr-description-craft`.
- Pre-implementation design docs — use `rfc-and-design-docs`.
- Sentence-level grammar/style review of prose — use `technical-writing-craft`.

## Core Concepts

### 1. The 50/72 rule (Tim Pope, 2008)

Subject line: hard limit 50 characters. Body: wrap at 72.

Why these numbers? Git's default tooling — `git log`, `git shortlog`, `git format-patch`, mailing-list bridges, GitHub's PR list — all assume the subject fits on one terminal line and the body has room for indentation. `git log --oneline` truncates anything past ~50. A 72-char body wrap leaves room for the 4-space indent `git log` applies, keeping everything under 80 columns.

**Mechanical rules:**
- Subject line ≤ 50 chars. Hard ceiling is 72, but cross 50 only when truly unavoidable.
- Blank line between subject and body — required. Many tools (`git log --pretty=%s%n%b`) rely on this separator.
- Body wrap at 72. Git wraps automatically only inside `git commit -m`'s string handling sometimes; use an editor that hard-wraps or wrap manually.
- No trailing period on the subject. The subject is a headline, not a sentence.

### 2. Imperative mood for the subject

Write the subject as a command the commit gives the codebase, not as a past-tense report of what you did.

Test it: prepend "If applied, this commit will ___". The result must be a grammatical English sentence.

- `Add retry logic to S3 uploader` — passes (If applied, this commit will add retry logic…).
- `Added retry logic to S3 uploader` — fails (If applied, this commit will added…).
- `Adds retry logic to S3 uploader` — fails the test too, even though it reads fine in `git log`. Stay imperative for consistency with Git's own auto-generated messages ("Merge branch 'foo'", "Revert 'bar'").
- `Fixed bug` — fails on every axis: past tense, vague, no scope.

### 3. Conventional Commits: structure beyond the basics

```
<type>[optional scope][!]: <description>

[optional body]

[optional footer(s)]
```

**Types (the canonical set):**
- `feat` — a new feature (correlates to MINOR in SemVer).
- `fix` — a bug fix (correlates to PATCH).
- `docs` — documentation only.
- `style` — formatting, whitespace, missing semis; no code-behavior change.
- `refactor` — code change that neither fixes a bug nor adds a feature.
- `perf` — performance improvement.
- `test` — adding or fixing tests.
- `build` — build system or external dependencies (npm, gradle, Docker).
- `ci` — CI configuration (GitHub Actions, CircleCI, Evergreen).
- `chore` — maintenance that doesn't fit elsewhere; release commits.
- `revert` — reverts a previous commit; body should reference the reverted SHA.

**Scope.** Optional parenthesized noun naming the affected area: `feat(auth):`, `fix(api):`, `refactor(db):`. Scope is a hint for the reader, not a taxonomy — keep it consistent within a repo but don't over-engineer.

**Description.** Imperative mood, lowercase first letter (this differs from Tim Pope's "Capitalized" rule — Conventional Commits prefers lowercase after the colon; pick one convention per repo and enforce it). No period.

### 4. Breaking changes: `!` and `BREAKING CHANGE:` footer

Two ways to signal a major-version-bumping breaking change:

**Method 1 — exclamation mark in the header:**
```
feat(api)!: remove deprecated /v1/users endpoint
```
The `!` MUST appear before the colon. If `!` is used, `BREAKING CHANGE:` MAY be omitted from the footer.

**Method 2 — footer token:**
```
feat(api): drop XML response format

BREAKING CHANGE: Clients sending Accept: application/xml now receive
406 Not Acceptable. Migrate to application/json or set
X-Legacy-XML-Fallback: 1 (sunset 2026-12-01).
```

**Rules from the spec:**
- `BREAKING CHANGE` is case-sensitive — uppercase only. `BREAKING-CHANGE` (hyphen) is also accepted as a synonym.
- A breaking change in ANY type triggers a MAJOR version bump — even `chore!:` or `docs!:` if it removes documented behavior.
- Use BOTH `!` and `BREAKING CHANGE:` footer when the breakage needs explanation longer than the subject can hold. The `!` is the index; the footer is the migration guide.

### 5. The "why not what" rule for the body

The diff already shows *what* changed. The body must explain *why*.

A diff says: `+ const TIMEOUT_MS = 30_000;`. The body must say: "Glean's p99 search latency crossed 28s during the May 12 incident. 30s gives 2s of slack before the Cloudflare 524 cuts the connection." That sentence is the entire reason the commit exists, and it cannot be reconstructed from the diff.

**Body content checklist:**
- What problem did this commit solve, or what behavior did it enable?
- Why this approach and not an alternative? (One sentence on the rejected alternative is often enough.)
- Any non-obvious side effect or trade-off the reviewer should know about?
- Links to tickets, design docs, or incident reports — but link out, don't paraphrase.

**Body anti-content:**
- Don't restate the subject line.
- Don't narrate the diff line by line ("Changed line 47 to call `retry()` instead of `fetch()`"). Reviewers can read the diff.
- Don't include implementation walkthroughs that belong in code comments.

### 6. Multi-commit storytelling for PR reviewers

A PR's commit list is a table of contents. Reviewers who can't load the whole diff at once read by commit. Plan the sequence:

1. **Prep commits first** — refactors, renames, type-only changes that have no behavior change. These shrink the diff for the substantive commits that follow.
2. **The substantive commit** — the one that does the actual feature/fix. The body of *this* commit carries the heavy reasoning.
3. **Test commit(s)** — if tests are separated, name them `test: add coverage for ...`. Many teams require tests in the same commit; respect the local convention.
4. **Cleanup commits last** — docs, changelog, version bumps.

Each commit on its own should build and pass tests (the `git bisect`-friendly property). If you have to break this for a multi-step rename, document the bisect-skip in a body line.

### 7. Fixup / squash / autosquash discipline

`git commit --fixup=<sha>` writes a commit message starting with `fixup! <original-subject>`. Later, `git rebase -i --autosquash <base>` reorders these fixups under their targets and marks them `fixup` (drop the message) or `squash` (keep both messages).

**Workflow:**
1. Develop normally; when you find a flaw in an earlier commit, do `git commit --fixup=abc123` instead of writing a new free-form commit.
2. Enable `git config rebase.autosquash true` (or pass `--autosquash` explicitly).
3. Before pushing the PR update, run `git rebase -i origin/main` and let autosquash reorganize.
4. Force-push the cleaned branch (`git push --force-with-lease`).

**Fixup vs squash:**
- `--fixup` discards the fixup's commit message; only the original is kept. Use for typo fixes, missed `console.log` removal, lint repairs.
- `--squash` keeps both messages joined; the editor opens for you to merge them. Use when the "fix" adds meaningful nuance.

**Never push fixup commits to a shared branch as-is.** They are scaffolding for your local history, not part of the merge story. CI hooks can block pushes containing `fixup!` or `squash!` subjects.

### 8. Signed-off-by and the Developer Certificate of Origin

The `Signed-off-by:` trailer is a per-commit affirmation that you wrote (or have the right to contribute) this code under the project's license. Required by Linux kernel, Node.js (moving to required), Kubernetes, many CNCF projects, Pi-hole, cert-manager.

```
fix(net): handle ECONNRESET during initial TLS handshake

The current code path treats a reset before the ServerHello as a
generic IO error, masking the actual TLS issue.

Signed-off-by: Mitchell Hudson <mitch.hudson@mongodb.com>
```

**Rules:**
- Use `git commit -s` (or `--signoff`) to add automatically — Git uses your `user.name` and `user.email`.
- The DCO covers the *commit*, not the *PR*. Every commit on a DCO-protected branch must have its own sign-off.
- DCO replaces the heavier CLA in most modern OSS projects. The text you are signing is at https://developercertificate.org/.

### 9. Other trailers worth knowing

- `Co-authored-by: Name <email>` — GitHub credits both authors in the contribution graph and "Co-authored by" UI. Used for pair programming and Claude-assisted commits.
- `Reviewed-by: Name <email>` — Linux kernel convention; the reviewer adds it after review.
- `Reported-by: Name <email>` — credits the bug reporter.
- `Fixes: <sha> ("<subject>")` — Linux kernel pattern linking a fix to the commit that introduced the bug. GitHub does not parse this; for issue-closing use `Closes #123` or `Fixes #123`.
- `Closes #123`, `Fixes #123`, `Resolves #123` — GitHub auto-closes the issue on merge into the default branch.

**Trailer ordering.** Trailers go at the end, separated from the body by a blank line. Order is conventional but not enforced; common pattern: issue links first, then `Co-authored-by`, then `Signed-off-by` last.

### 10. Splitting commits — one logical change per commit

A commit should be the smallest atomic change that still builds and passes tests. Sign that you're commit-splitting correctly with these heuristics:

- The subject line uses "and" or a comma joining two ideas → split.
- The diff touches unrelated modules with no shared rationale → split.
- Reverting the commit would require also reverting unrelated improvements → split.
- The body needs more than one paragraph just to list distinct things → split.

Use `git add -p` (patch mode) or `git add -i` (interactive) to stage hunks selectively. `git reset -p` to unstage hunks. `git commit --patch` to launch staging from within the commit flow.

## Templates and Examples

### Template — minimal Conventional Commit

```
<type>(<scope>): <imperative description in 50 chars or less>
```

### Template — Conventional Commit with body and footers

```
<type>(<scope>): <imperative description>

<Why this change exists. What problem it solves or what behavior it
enables. Wrap at 72 columns. Multiple paragraphs are fine.>

<Any non-obvious side effects, trade-offs, or follow-ups.>

Closes #1234
Co-authored-by: Pat Reviewer <pat@example.com>
Signed-off-by: Mitchell Hudson <mitch.hudson@mongodb.com>
```

### Example — fix with explanation

```
fix(s3-uploader): retry on TLS handshake reset

The uploader treated any error before bytes were sent as fatal. AWS
NLB occasionally resets the TLS handshake under load (observed in
incident #4821, ~3% of uploads during the 14:00 UTC spike). Retrying
twice with 250ms / 750ms backoff eliminates the symptom without
masking genuine connection failures.

Tested against the recorded NLB-reset trace; all 200 captured
sessions now succeed within the retry budget.

Closes #4821
Signed-off-by: Mitchell Hudson <mitch.hudson@mongodb.com>
```

### Example — breaking change with `!` and footer

```
feat(api)!: drop /v1 endpoints

/v1 was sunset 2026-04-01 per the deprecation banner shipped in
v3.2.0 (six months ago). The /v2 contract is a strict superset; the
migration guide lives at docs/migration/v1-to-v2.md.

BREAKING CHANGE: All /v1/* routes now return 410 Gone. Clients still
sending /v1 traffic will see immediate failures — verify your
client's base URL is set to /v2 before deploying this version.

Closes #2901
```

### Example — refactor (no behavior change)

```
refactor(case-enricher): extract severity-rank pure helper

Moves the severity-string-to-numeric mapping out of the enricher
into its own module so the dashboard renderer can use the same
ordering without re-importing the enricher.

No behavior change. Same inputs produce same outputs.
```

### Example — revert

```
revert: feat(auth): enable WebAuthn registration flow

This reverts commit 8f3a2b1c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9.

WebAuthn registration triggered a 100% failure for users on iOS 17.2
with attached security keys (incident #5102). Reverting while we
investigate the platform interaction. Re-enable tracked in #5103.
```

### Example — fixup that will autosquash

```
fixup! feat(s3-uploader): retry on TLS handshake reset
```

(No body — the autosquash rebase will fold this into the original commit and discard this subject.)

### Example — commit produced by AI assistant with attribution

```
feat(dashboard): add ownerless-case acknowledgement popup

Surfaces cases without an assigned owner in a non-modal popup so
the on-call TAM can claim them in one click instead of digging
through the unfiltered queue.

Co-authored-by: Claude <noreply@anthropic.com>
Signed-off-by: Mitchell Hudson <mitch.hudson@mongodb.com>
```

### Example — splitting a "kitchen sink" commit

Before (one commit):
```
update auth, fix typos in readme, bump deps, refactor logger
```

After (four commits in dependency order):
```
chore(deps): bump @anthropic-ai/sdk from 0.39.0 to 0.41.0
refactor(logger): split level-filtering into pure helper
fix(auth): handle expired refresh-token race in token-rotator
docs(readme): correct three typos in install section
```

## Anti-Patterns

- **"Fixed bug" / "Updates" / "WIP"** — useless. The subject line is the most-read part of the commit; spend a sentence on it.
- **Past-tense subjects** — `Added`, `Fixed`, `Refactored`. Breaks Tim Pope's "If applied" test and clashes with Git's own auto-messages.
- **Restating the diff in the body** — "Changed `let` to `const` on line 47." Reviewers read diffs; bodies are for *why*.
- **Stacking unrelated changes in one commit** — defeats `git bisect`, makes `git revert` either useless or destructive.
- **Burying breaking changes in the body without `!` or `BREAKING CHANGE:`** — release-tooling (semantic-release, release-please) won't bump the major version, and downstream consumers won't see the warning.
- **`feat!:` for additive features without removals** — `!` means *breaking*, not *exciting*. Adding a new endpoint isn't breaking; removing a parameter is.
- **`Co-authored-by` without the trailing blank line** — GitHub silently ignores the trailer. Trailers must be in the trailer block.
- **Pushing `fixup!` / `squash!` commits to the PR branch** — many teams' CI blocks merging until autosquash collapses them. Run `git rebase -i --autosquash` before pushing.
- **Mixing `chore` and `fix` for the same change** — pick one. Bug fixes are `fix`. Dependency bumps are `chore(deps)`. Routine maintenance is `chore`.

## Decision Heuristics

- **`feat` or `fix`?** New capability the user didn't have before → `feat`. Behavior that was supposed to work but didn't → `fix`. A `feat` that fixes a long-standing usability gap might still be a `feat`; ship it as a feature so it shows up in MINOR-bump release notes.
- **`refactor` or `chore`?** If the change is in production source paths and could in theory affect behavior → `refactor`. If it's in tooling, config, or CI → `chore` or `ci` or `build`.
- **Body or no body?** If the diff is one line of obvious cleanup, the subject is the whole commit. If a reviewer might ask "why?" → write a body.
- **One commit or many?** Run the bisect test: if a future bisect lands on this commit and breaks, will the body tell the bisector enough to know whether the breakage is real? If no, split.
- **Squash-merge or merge-commit?** Squash-merge collapses the PR's commit storytelling into one message — only worth it for tiny PRs. For PRs with deliberate commit sequencing, prefer `git merge --no-ff` or rebase-merge so the structure survives.
- **Add `BREAKING CHANGE:` to the footer when…** the change removes, renames, or alters the semantics of any documented public interface — even if a deprecation banner shipped six months ago. Removal-after-deprecation is still a breaking change for users who didn't migrate.
- **Sign off with `-s` when…** contributing to any CNCF project, Linux kernel, or any repo whose CONTRIBUTING.md or DCO bot enforces it. Cheap to make a default with `git config commit.gpgsign` and an alias.

## References

- [Conventional Commits v1.0.0 specification](https://www.conventionalcommits.org/en/v1.0.0/) — canonical spec for `<type>(<scope>): <description>` plus `BREAKING CHANGE:` and `!`.
- [Tim Pope — A Note About Git Commit Messages (tbaggery.com, 2008)](https://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html) — original 50/72 rule, imperative mood, and body discipline.
- [Chris Beams — How to Write a Git Commit Message (cbea.ms)](https://cbea.ms/git-commit/) — the seven canonical rules of a great commit message, expanding on Pope.
- [Conventional Commits Cheatsheet (qoomon, GitHub)](https://gist.github.com/qoomon/5dfcdf8eec66a051ecd85625518cfd13) — practical reference of types, breaking-change syntax, and footer patterns.
- [Linux Foundation — Developer Certificate of Origin (DCO)](https://wiki.linuxfoundation.org/dco) — the text you sign with `Signed-off-by:`.
- [Developer Certificate of Origin text (developercertificate.org)](https://developercertificate.org/) — what each `Signed-off-by` line attests to.
- [thoughtbot — Auto-squashing Git Commits](https://thoughtbot.com/blog/autosquashing-git-commits) — `--fixup` and `--autosquash` workflow.
- [Atlassian Bitbucket — Git Best Practices for Commit Messages](https://community.atlassian.com/forums/Bitbucket-articles/Git-Best-Practices/ba-p/1628803) — imperative-mood guidance and structural conventions.
- [GitHub Docs — Creating a commit with multiple authors (Co-authored-by)](https://docs.github.com/en/pull-requests/committing-changes-to-your-project/creating-and-editing-commits/creating-a-commit-with-multiple-authors) — Co-authored-by trailer format.
