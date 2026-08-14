# Vendored Superpowers Skills — Inventory

Upstream: https://github.com/obra/superpowers (MIT). Sync only against release tags.

**Last synced tag: v6.1.1** (2026-07-13)

## Vendored skills

| Skill | Repo category | Upstream path |
|---|---|---|
| brainstorming | before-you-code | skills/brainstorming |
| writing-plans | before-you-code | skills/writing-plans |
| executing-plans | before-you-code | skills/executing-plans |
| dispatching-parallel-agents | before-you-code | skills/dispatching-parallel-agents |
| systematic-debugging | while-you-code | skills/systematic-debugging |
| subagent-driven-development | while-you-code | skills/subagent-driven-development |
| verification-before-completion | before-you-merge | skills/verification-before-completion |
| receiving-code-review | before-you-merge | skills/receiving-code-review |
| finishing-a-development-branch | before-you-merge | skills/finishing-a-development-branch |

Repo layout per skill: `.agents/skills/<category>/<name>/skills/<name>/SKILL.md`, supporting `.md` files under `references/`, executable helpers under `scripts/`. Upstream keeps supporting files flat next to SKILL.md.

## Pruned skills — do NOT re-add; redirect references

Removed in commit `3fb5ab2` ("Prune environment-locked obra skills"); rationale in that commit message.

| Upstream skill | Why pruned | Redirect references to |
|---|---|---|
| test-driven-development | RED/GREEN loop hardcoded Jest/npm | Inline TDD wording ("write a failing test first"), no skill link |
| using-git-worktrees | Duplicated tooling/git-workflows | `git-workflows` |
| requesting-code-review | MMS has a richer equivalent | `subagent-driven-development`'s `references/code-reviewer.md` (we keep a vendored copy of upstream's code-reviewer.md there) |
| using-superpowers | Plugin-runtime bootstrap, meaningless standalone | Drop the reference (including `../using-superpowers/references/` platform notes) |
| writing-skills | Plugin-specific authoring flow | Drop the reference |

## Standing customizations (apply on every sync)

Global, all skills:

1. **Frontmatter** — keep `source`, `license`, `mongodb:` (team/owner/internal) fields; never add `version` or `when_to_use`. Adopt upstream `description` changes when they occur.
2. **De-namespacing** — `superpowers:<name>` → `<name>` everywhere.
3. **References layout** — supporting `.md` files live in `references/`, links rewritten accordingly.
4. **No harness-specific tool names** where upstream drops them (e.g. upstream itself moved from "TodoWrite" to "create todos" — prefer the neutral wording).

Per skill:

| Skill | Customization to preserve |
|---|---|
| brainstorming | Visual Companion removed entirely (checklist item, section, `visual-companion.md`, `scripts/`) — the companion server does not ship in this repo (commit `f0d9393`). SKILL.md explicitly links `references/spec-document-reviewer-prompt.md` (upstream leaves it unreferenced; our validator flags orphans) |
| writing-plans | Task Structure template is language-neutral: placeholder paths, prose "paste concrete test code in the target language" instead of Python-only examples, multi-runner command examples (`go test`, `pytest`, `npm test`, `resmoke`) (commit `3fb5ab2`). SKILL.md explicitly links `references/plan-document-reviewer-prompt.md` (same orphan reasoning as brainstorming) |
| dispatching-parallel-agents | Example scenarios are a mix of build-fix / log-investigation / script-authoring instead of upstream's `.test.ts`-only examples (commit `3fb5ab2`) |
| executing-plans | Integration section points to `git-workflows` for worktree setup |
| systematic-debugging | MongoDB reference files `references/server.md` + `references/cloud.md` and their index entries in SKILL.md; TDD-skill pointer replaced with inline "write a failing test/repro first" wording; `references/condition-based-waiting.md` drops upstream's pointer to the excluded `condition-based-waiting-example.ts` |
| subagent-driven-development | `references/code-reviewer.md` (vendored from upstream requesting-code-review) is the final whole-branch review template; links point there instead of `../requesting-code-review/`. Upstream's `scripts/` (review-package, task-brief, sdd-workspace) are vendored as-is — re-add the executable bit when committing (`git add --chmod=+x`) |
| verification-before-completion | "Fern-Specific Verification" section (finalize-branch pointer) |
| finishing-a-development-branch | `bazel test` included in the test-runner examples list |
| receiving-code-review | none beyond globals |

## Excluded upstream files (do not vendor)

- `CREATION-LOG.md`, `test-academic.md`, `test-pressure-*.md` (upstream's skill-authoring test fixtures)
- brainstorming: `visual-companion.md`, `scripts/` (companion server not shipped)
- requesting-code-review: everything except `code-reviewer.md` (vendored into subagent-driven-development)
- systematic-debugging: `condition-based-waiting-example.ts`, `find-polluter.sh` (kept out to date; revisit if SKILL.md starts linking them)

## Origin history

- Imported ~upstream v5.0.7: commits `4d0dc49` (before-you-code), `0e3d954` (systematic-debugging + server/cloud refs), `0594ce8`, `fd5b540` (before-you-merge); PR #6 added the execution-chain skills, PR #13 the query-integration symlinks, PR #28 the `skills/<name>/SKILL.md` layout refactor.
- Customization commits: `3fb5ab2` (prune + stack-neutral scrub), `f0d9393` (remove visual companion), `60a4ecd` (frontmatter fixes, references/ move), `b342d3e` (drop when_to_use), `c4d29be` (internal flag).
- v5.0.7 → v6.1.1 sync: 2026-07-13 (this skill's first use).
