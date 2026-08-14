<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `code-reviewer` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: code-reviewer
description: >
  Practical code-review reference covering reviewer mindset, workflow, comment quality, GitHub PR mechanics, and OWASP security checklists. Use when reviewing pull requests, auditing a codebase, deciding whether to approve/request-changes/comment, or running a security sweep.
  TRIGGER: "review this PR", "code review", "should I approve", "review this diff", "security review", "OWASP checklist", "comment on this change", "is this ready to merge".
  SKIP: user needs to write new code (not review existing); document critique target is prose not code (use document-critique).
origin: local
version: 1.1.0
updated: "2026-05-29"
category: developer
tags: [code-review, pull-request, github, owasp, security, reviewer-mindset]
related_skills:
  - coding-standards
  - coding-patterns
  - security-reviewer
  - document-critique
whenToUse:
  - reviewing a pull request or diff
  - deciding whether to approve, comment, or request changes
  - running a security sweep on code changes
  - writing effective review comments
  - handling reviewer pushback or disagreement
  - setting up CODEOWNERS or branch protection
whenNotToUse:
  - writing new code (not reviewing existing)
  - critiquing a prose document (use document-critique)
  - language/framework-specific correctness checks without a PR context
---

# Code Reviewer

Practical reference for reviewing pull requests: reviewer mindset, workflow, comment quality, GitHub mechanics, and OWASP security coverage.

**Sources:** Google Engineering Practices (reviewer standard, comments, speed, pushback), GitHub Docs (PR review states, CODEOWNERS), OWASP Code Review Guide v2.0 + Secure Coding Practices v2.0.1.

## Quick review rules

1. **Approve once a change definitely improves overall code health** — not perfection, but net improvement ([Google standard](https://google.github.io/eng-practices/review/reviewer/standard.html)).
2. **Start with purpose and high-level design** before line comments; reject early if the design is wrong ([Google navigate](https://google.github.io/eng-practices/review/reviewer/navigate.html)).
3. **Review for design, functionality, complexity, tests, naming, comments, style, and documentation** — not just syntax ([Google what to look for](https://google.github.io/eng-practices/review/reviewer/looking-for.html)).
4. **Review every line in your scope**; pull in broader context when the diff alone is insufficient.
5. **Be kind, specific, and explicit about why** a change is needed; comment on the code, not the person ([Google comments](https://google.github.io/eng-practices/review/reviewer/comments.html)).
6. **Respond within one business day** — fast response beats perfect uninterrupted throughput ([Google speed](https://google.github.io/eng-practices/review/reviewer/speed.html)).
7. **Prefer small, reviewable PRs**; ask for a split or send high-level feedback first on large changes.
8. **Use GitHub review states deliberately**: Comment (non-blocking), Approve (merge-ready), Request changes (must fix before merge).
9. **Always run a security sweep** on data handling, trust boundaries, auth, session, access control, output encoding, logging, and configuration.

## Review workflow

| Step | Action | Why |
|------|---------|-----|
| 1. Understand intent | Read the change description; determine whether the change belongs at all | Reject bad designs early before wasting line-review time |
| 2. Check main design | Review the architectural center before peripheral details | Major design problems invalidate many smaller comments |
| 3. Review the rest | Tests-first can clarify intended behavior before reading implementation | Consistent, deliberate ordering |
| 4. Inspect behavior | Edge cases, concurrency, race/deadlock risk, over-engineering | Hard to catch by execution alone |
| 5. Verify tests + docs | Unit/integration tests, updated docs for changed flows | Production changes should ship their tests |
| 6. Security checklist | OWASP coverage across auth, sessions, encoding, secrets, config | Manual review catches what scanners miss |
| 7. Submit review state | Pending → Comment/Approve/Request changes | Platform workflow |
| 8. Resolve disagreement | Reconsider if author is right; explain code-health rationale; escalate if needed | Avoid "clean it up later" drift |

## Reviewer mindset

- **Core goal:** improve overall code health over time, not demand perfection in every patch.
- **Technical facts outrank opinion.** Style-guide requirements are binding; preferences outside the guide are not merge blockers.
- **Code review is a mentoring surface.** Educational comments are encouraged but must be marked non-blocking.
- **OWASP lens:** manual secure code review finds issues automated scanners miss; it belongs in the SDLC, not as a last-minute add-on.

## What to look for

### Non-security areas

| Area | Key questions |
|------|---------------|
| Overall design | Does this change belong here? Does it integrate well? Is now the right time? |
| Functionality | Does it work from both user and developer perspectives, including edge cases? |
| Concurrency | Race conditions, deadlocks, missing synchronization? |
| Complexity | Is it solving the actual current problem, not speculative future ones? |
| Tests | Do they fail for real breakage, avoid false positives, make useful assertions? |
| Names + comments | Descriptive names, comments that explain *why*, updated docs for changed flows? |

### Security areas (OWASP checklist coverage)

**Input validation** · **Output encoding** · **Authentication + password management** · **Session management** · **Access control** · **Cryptographic practices** · **Error handling + logging** · **Data protection** · **Communication security** · **System configuration** · **Database security** · **File management** · **Memory management** · **General coding practices**

Key OWASP principles: validate all untrusted input on a trusted system; use contextual output encoding; centralize auth/authz controls; use framework session primitives; avoid sensitive data in errors; keep systems patched.

## Writing effective review comments

- **Comment on the code, not the person.** "This function does X when Y should happen" not "You wrote this wrong."
- **Explain why** when the rationale isn't obvious — state the code-health, maintainability, or user-impact reason.
- **Balance pointing out the problem with prescribing the exact solution.** Concrete suggestions help; the author is still responsible for the fix.
- **Label comment strength:** `Nit:`, `Optional:` / `Consider:`, `FYI:` so authors know what must be fixed vs what is informational.
- **Call out what is done well** — reinforcement is part of mentoring culture.
- When code is hard to understand, ask for clarification or simplification rather than silently guessing intent.

## Balancing quality, speed, and scope

- Respond within **one business day** unless in focused work — wait for a natural break.
- **LGTM with comments** is valid when remaining items are minor or the author is trusted to handle them; make that intent explicit.
- **Ask for PR splits** on large changes, or at minimum provide high-level design comments first.
- Prefer **style-only changes in separate PRs** to keep reviews understandable and rollback-safe.

## GitHub PR review mechanics

| Feature | Use for |
|---------|---------|
| **Comment** | Non-blocking feedback, informational notes |
| **Approve** | Merge-ready; clearly improves code health |
| **Request changes** | Issues that must be fixed before merge |
| **Suggested changes** | Exact patch proposals the author can apply in one click |
| **Viewed markers** | Track file-by-file progress; file collapses until changed again |
| **CODEOWNERS** | Auto-request owners for files; gate merge with branch protection |

Notes:
- Pending comments are visible only to the reviewer until submitted.
- **Request changes** only blocks merge when protected branches or rulesets are configured.
- PR authors cannot approve their own PRs.
- CODEOWNERS files must be on the base branch and follow GitHub syntax constraints.

## Standards and best practices

### Acceptance threshold
Approve when the change **definitely improves code health**, even if not perfect. Do not approve changes that introduce known maintainability decline, unjustified complexity, or quality regression outside emergencies.

### Review scope control
- Keep major style-only changes separate from functional changes.
- If a PR is too large, split or deliver high-level comments first.

### Test and documentation expectations
- Production changes should bring tests in the same change (except emergencies).
- Doc updates are required when the change alters how users build, test, use, release, remove, or deprecate code.

### Security discipline
- Don't postpone security review because automated tooling exists.
- Use OWASP checklist areas as the minimum frame for every application code review.

### Maintainability
- Favor simplification over speculative flexibility.
- Insist on fixing newly introduced complexity now, not "later," unless a real emergency.
- For discovered-but-unblockable debt: file a bug and leave a TODO tied to it.

## Known scope limits

- **Google reviewer guide** is process- and philosophy-heavy, not language-specific. It doesn't replace language/framework best-practice docs.
- **GitHub Docs** describe workflow mechanics and repository controls, not the team's quality bar.
- **OWASP Code Review Guide v2.0** uses the 2013 Top 10 era taxonomy — use it for structured coverage, not as the only contemporary exploit reference.
- **OWASP Secure Coding Practices v2.0.1** project page is archived; the PDF remains valid as a checklist source.
