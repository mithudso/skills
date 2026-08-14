<!-- hub-reference-banner -->
> **Reference file — part of the `devops-containers-cicd` hub.** Formerly the standalone `git-workflows` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: git-workflows
description: >
  Git branching strategies, commit conventions, merge strategies, CI integration,
  monorepo patterns, cherry-picking, branch protection, and release automation.
  Covers trunk-based development, GitFlow, GitHub Flow, GitLab Flow, conventional
  commits, semantic versioning, merge vs rebase vs squash, monorepo tooling (Nx,
  Turborepo), git hooks (Husky, Lefthook), and branch protection rulesets.
  TRIGGER: choosing a branching model; setting up conventional commits or commitlint;
  designing merge/rebase/squash policy; configuring branch protection or rulesets;
  wiring semantic-release / release-please / Changesets; setting up git hooks
  (Husky, Lefthook, pre-commit); affected-based monorepo CI (Nx, Turborepo);
  cherry-picking and backport strategy.
  SKIP: CI/CD pipeline and deployment design beyond git triggers (use devops-infra);
  package publishing and distribution mechanics (use code-packaging); language-agnostic
  software design patterns (use software-engineering-patterns).
triggers:
  - git workflow
  - branching strategy
  - trunk-based development
  - gitflow
  - github flow
  - gitlab flow
  - conventional commits
  - commitlint
  - semantic versioning
  - semver
  - release-please
  - semantic-release
  - merge vs rebase
  - squash merge
  - monorepo
  - Nx affected
  - Turborepo
  - cherry-pick
  - git hooks
  - husky
  - lefthook
  - pre-commit hook
  - branch protection
  - rulesets
  - feature flags git
version: 1.0.1
updated: 2026-05-31
category: developer
whenToUse:
  - Choosing a branching model (trunk-based, GitFlow, GitHub Flow, GitLab Flow) for a team
  - Setting up conventional commits and commitlint enforcement
  - Designing a merge vs rebase vs squash policy for pull requests
  - Configuring GitHub branch protection rules or rulesets on a protected branch
  - Wiring release automation with semantic-release, release-please, or Changesets
  - Installing and configuring git hooks via Husky, Lefthook, or pre-commit
  - Setting up affected-based monorepo CI with Nx or Turborepo
  - Planning a cherry-pick or backport strategy across release branches
related_skills:
  - devops-infra
  - software-engineering-patterns
  - code-packaging
---

# Git Workflows -- Expert Reference

## 1. Overview

This skill covers the Git workflow patterns used by modern
development teams. It addresses branching strategies, commit conventions, merge
strategies, CI/CD integration, monorepo tooling, cherry-picking, git hooks, branch
protection, and release automation. Use this skill when choosing a branching model,
configuring commit enforcement, designing merge policies, or setting up automated
releases.

**Decision framework -- choosing a workflow:**

| Factor | Trunk-Based | GitFlow | GitHub Flow | GitLab Flow |
|---|---|---|---|---|
| Release cadence | Continuous | Scheduled | Continuous | Staged |
| Team size | Any (best < 50) | Medium-large | Small-medium | Medium-large |
| Environment count | 1 (prod) | 2+ | 1 | 2+ (staging, prod) |
| Versioned releases | Via tags | Via release branches | Via tags | Via env branches |
| Complexity | Low | High | Low | Medium |

---

## 2. Branching Strategies

### 2.1 Trunk-Based Development (TBD)

Trunk-based development is a branching strategy where all developers commit small,
incremental changes to a single shared branch (main/trunk) at least once a day. It
is the workflow most aligned with continuous integration and continuous delivery.

**Core rules:**

- Every developer integrates into trunk at least daily.
- Feature branches, when used, are short-lived (hours, never more than 1-2 days).
- Code that is not ready for users is hidden behind feature flags, not on long-lived branches.
- The trunk is always in a releasable state; automated tests gate every merge.
- Releases are cut directly from trunk (via tags or automated pipelines).

**Feature flags in TBD:**

Feature flags decouple deployment from release. Incomplete features are merged to
main within a flag path and activated later. Flag types include:

| Flag type | Lifecycle | Example |
|---|---|---|
| Release toggle | Short (days-weeks) | `--withOneClickPurchase` |
| Experiment / A-B | Medium (weeks) | `showNewCheckoutFlow` |
| Ops toggle | Long (permanent) | `enablePartnerSync` |
| Permission toggle | Long (permanent) | `allowBetaUsers` |

Flag best practices:

- Keep flags short-lived; automate cleanup to remove unused flags.
- Document each flag's purpose, owner, and planned removal date.
- Use dependency injection or abstraction layers rather than scattered `if` checks.
- Fan out CI after unit tests to cover meaningful flag permutations.
- Store runtime flag state in durable distributed stores (Consul, Etcd, LaunchDarkly).

**When to use TBD:** Teams practicing CI/CD, deploying frequently, with strong
automated test coverage. Google, Meta, and most high-performing engineering orgs
use trunk-based development.

> Sources: [Atlassian -- Trunk-Based Development](https://www.atlassian.com/continuous-delivery/continuous-integration/trunk-based-development), [trunkbaseddevelopment.com -- Feature Flags](https://trunkbaseddevelopment.com/feature-flags/), [FeatBit -- TBD Feature Flags 2025](https://www.featbit.co/articles2025/trunk-based-development-feature-flags-2025), [Flagsmith -- TBD Guide](https://www.flagsmith.com/blog/trunk-based-development)

### 2.2 GitFlow

GitFlow, introduced by Vincent Driessen in 2010, is a branching model designed for
projects with scheduled releases and multiple versions in production.

**Branch types:**

| Branch | Lifetime | Branches from | Merges into |
|---|---|---|---|
| `main` | Permanent | -- | -- |
| `develop` | Permanent | `main` | -- |
| `feature/*` | Temporary | `develop` | `develop` |
| `release/*` | Temporary | `develop` | `main` + `develop` |
| `hotfix/*` | Temporary | `main` | `main` + `develop` |

**Workflow:**

1. New features branch from `develop` into `feature/<name>`.
2. Completed features merge back into `develop` via PR.
3. When ready for release, a `release/<version>` branch is cut from `develop`.
4. Final testing, bug fixes, and version bumps happen on the release branch.
5. The release branch merges into both `main` (tagged) and `develop`.
6. Hotfixes branch from `main`, are fixed, then merge into both `main` and `develop`.

**When to use GitFlow:** Projects with explicit version numbers (mobile apps,
SDKs, libraries), regulated environments requiring release documentation, or teams
managing multiple supported versions simultaneously.

**When NOT to use GitFlow:** Continuously deployed web apps and services. Vincent
Driessen himself added a 2020 reflection noting that GitFlow was designed before
GitHub-style CI/CD existed, and simpler models are preferable for web applications
with continuous delivery.

> Sources: [nvie.com -- A Successful Git Branching Model](https://nvie.com/posts/a-successful-git-branching-model/), [Atlassian -- GitFlow Workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow), [git-flow cheatsheet](https://danielkummer.github.io/git-flow-cheatsheet/), [AWS -- Branches in GitFlow](https://docs.aws.amazon.com/prescriptive-guidance/latest/choosing-git-branch-approach/branches-in-a-gitflow-strategy.html)

### 2.3 GitHub Flow

GitHub Flow is a lightweight branching model built for teams that deploy frequently.
The main branch is always production-ready.

**Workflow:**

1. Create a branch from `main` for every feature or bugfix.
2. Commit changes to the branch with descriptive messages.
3. Open a Pull Request for code review and discussion.
4. Automated CI runs tests on the PR.
5. After approval and passing checks, merge into `main`.
6. Deploy immediately from `main`.

**Key properties:**
- Only one long-lived branch (`main`).
- No release branches -- releases are tags or automated from main.
- PRs are the unit of change, review, and discussion.
- Ideal for continuous deployment with no staging gates.

**When to use:** Small-to-medium teams, web applications, SaaS products, any project
with a single production version and continuous deployment.

> Sources: [Ei Square -- Gitflow vs GitHub Flow vs GitLab Flow](https://www.eisquare.co.uk/blogs/how-to-choose-your-branching-strategy), [FUEiNT -- GitHub Flow vs GitLab Flow](https://fueint.com/blog/github-gitlab), [GitKraken -- Git Branch Strategy](https://www.gitkraken.com/learn/git/best-practices/git-branch-strategy)

### 2.4 GitLab Flow

GitLab Flow combines GitHub Flow's simplicity with environment-based promotion
branches for teams that need controlled deployments.

**Environment branch model:**

```
feature/* --> main --> staging --> production
```

- Features merge into `main` (integration).
- `main` merges into `staging` for pre-production validation.
- `staging` merges into `production` for live deployment.
- Deployment promotion is an explicit merge, not an automatic pipeline trigger.

**Release branch model (alternative):**

For projects shipping versioned releases (mobile apps, on-premise software):

```
feature/* --> main --> release/1.0 --> release/1.1
```

Bug fixes are cherry-picked upstream from release branches back to main.

**When to use:** Teams needing explicit environment promotion, regulated industries,
projects with staging/production separation, or projects with both continuous
integration and scheduled releases.

> Sources: [FUEiNT -- GitHub Flow vs GitLab Flow](https://fueint.com/blog/github-gitlab), [Medium -- Comparing Git Branching Strategies](https://medium.com/thedevproject/comparing-git-branching-strategies-git-flow-vs-github-flow-vs-gitlab-flow-2e1dd28be103), [DevOps Knowledge Hub -- Branching Strategy Comparison](https://devops.aibit.im/article/git-branching-strategy-comparison)

---

## 3. Commit Conventions

### 3.1 Conventional Commits Specification (v1.0.0)

The Conventional Commits specification provides a lightweight convention for
structured commit messages that integrate with semantic versioning.

**Format:**

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Required elements:**
- `type`: A noun describing the category (feat, fix, etc.)
- `description`: A short summary after the colon and space

**Optional elements:**
- `scope`: Contextual section in parentheses, e.g., `feat(parser):`
- `body`: Longer explanation, separated by a blank line
- `footer(s)`: Metadata like `BREAKING CHANGE:`, `Reviewed-by:`, `Refs:`

**Standard types:**

| Type | Purpose | SemVer impact |
|---|---|---|
| `feat` | New feature | MINOR bump |
| `fix` | Bug fix | PATCH bump |
| `docs` | Documentation only | None |
| `style` | Formatting, whitespace | None |
| `refactor` | Code restructuring, no behavior change | None |
| `perf` | Performance improvement | None |
| `test` | Adding or fixing tests | None |
| `build` | Build system or external dependencies | None |
| `ci` | CI configuration and scripts | None |
| `chore` | Maintenance tasks | None |

**Breaking changes:** Indicated by `!` after type/scope or by a `BREAKING CHANGE:`
footer. Both correlate with a MAJOR version bump. BREAKING CHANGE must be uppercase.

**Examples:**

```
feat(auth): add OAuth 2.0 login flow

Implements the full OAuth 2.0 authorization code flow with PKCE
for the authentication module.

Refs: #1234
```

```
fix!: correct rate-limiter token bucket overflow

BREAKING CHANGE: rate limiter now rejects requests exceeding
the configured burst size instead of queuing them.
```

### 3.2 Commit Message Enforcement with commitlint

commitlint checks commit messages against a configured rule set. The standard
config `@commitlint/config-conventional` enforces the Angular convention.

**Installation:**

```bash
npm install --save-dev @commitlint/{cli,config-conventional}
echo "export default { extends: ['@commitlint/config-conventional'] };" > commitlint.config.mjs
```

**Integration with Husky:**

```bash
npx husky init
echo 'npx --no -- commitlint --edit "$1"' > .husky/commit-msg
```

**Custom rules (commitlint.config.mjs):**

```js
export default {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [2, 'always', [
      'feat', 'fix', 'docs', 'style', 'refactor',
      'perf', 'test', 'build', 'ci', 'chore', 'revert'
    ]],
    'scope-case': [2, 'always', 'kebab-case'],
    'subject-max-length': [2, 'always', 72],
    'body-max-line-length': [2, 'always', 100],
  },
};
```

> Sources: [Conventional Commits v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/), [AverageDevs -- Conventional Commits in Git](https://www.averagedevs.com/blog/conventional-commits-git), [ShakilTech -- Conventional Commits Complete Guide](https://blog.shakiltech.com/conventional-commits/), [Jeff Bailey -- What Are Conventional Commits](https://jeffbailey.us/blog/2025/09/28/what-are-conventional-commits/)

---

## 4. Semantic Versioning and Release Automation

### 4.1 Semantic Versioning (SemVer)

SemVer uses the format `MAJOR.MINOR.PATCH`:

- **MAJOR**: Incompatible API changes (breaking changes).
- **MINOR**: Backward-compatible new functionality.
- **PATCH**: Backward-compatible bug fixes.

Pre-release versions use suffixes: `1.0.0-alpha.1`, `1.0.0-beta.2`, `1.0.0-rc.1`.

### 4.2 Automated Release Tools

**semantic-release:**
Fully automated: analyzes commit messages (Conventional Commits), determines the
next version, generates changelogs, creates git tags, and publishes to npm/GitHub.

```yaml
# .github/workflows/release.yml
name: Release
on:
  push:
    branches: [main]
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-node@v4
        with:
          node-version: 22
      - run: npm ci
      - run: npx semantic-release
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          NPM_TOKEN: ${{ secrets.NPM_TOKEN }}
```

**release-please (Google):**
Creates/updates a Release PR with changelog preview. Merging the PR triggers
publishing. Better for teams wanting human review of release notes.

```yaml
# .github/workflows/release-please.yml
on:
  push:
    branches: [main]
jobs:
  release-please:
    runs-on: ubuntu-latest
    steps:
      - uses: googleapis/release-please-action@v4
        with:
          release-type: node
```

**Changesets:**
Developer-driven: requires explicit changeset files describing changes. Best for
monorepos and teams wanting collaborative changelog creation.

```bash
npx changeset        # developer creates changeset
npx changeset version # bumps versions + updates changelogs
npx changeset publish # publishes to npm
```

**Comparison:**

| Tool | Automation level | Human review | Monorepo support | Best for |
|---|---|---|---|---|
| semantic-release | Full | None | Plugin-based | Single-package, CI/CD |
| release-please | Semi | PR review | Built-in | Teams wanting review |
| Changesets | Manual | Explicit | Native | Monorepos, fine control |

> Sources: [semantic-release GitHub](https://github.com/semantic-release/semantic-release), [MerginIT -- semantic-release Guide 2025](https://merginit.com/blog/29062025-automated-multi-platform-releases), [Oleksii Popov -- NPM Release Automation](https://oleksiipopov.com/blog/npm-release-automation/), [AWS -- Semantic Versioning](https://aws.amazon.com/blogs/devops/using-semantic-versioning-to-simplify-release-management/)

---

## 5. Merge Strategies

### 5.1 Merge Commits (`git merge`)

Creates a new commit with two parents, preserving the full branch history.

```
main:    A---B-------M
              \     /
feature:       C---D
```

**Pros:** Full history preserved, easy to revert entire features, shows exactly
when integration happened.

**Cons:** Noisy history with many merge commits, harder to bisect.

**Use when:** Merging shared branches, team-level integration, preserving context
of collaboration.

### 5.2 Squash Merge (`git merge --squash`)

Combines all feature branch commits into a single commit on the target branch.

```
main:    A---B---S (squashed C+D)
              \
feature:       C---D
```

**Pros:** Clean linear history on main, each PR is one atomic commit, easy to
revert a whole feature.

**Cons:** Loses granular commit history, feature branch becomes orphaned (no
merge ancestor).

**Use when:** PRs with many WIP commits, enforcing one-commit-per-feature policy.

### 5.3 Rebase (`git rebase`)

Replays feature branch commits on top of the target branch tip, creating new
commit hashes.

```
Before:  A---B (main)
              \
               C---D (feature)

After:   A---B---C'---D' (feature, rebased)
```

**Pros:** Clean, linear history without merge commits, easier to bisect, each
commit is individually meaningful.

**Cons:** Rewrites history (new hashes), dangerous on shared branches, can
cause conflicts commit-by-commit.

**The Golden Rule of Rebasing:** Never rebase commits that have been pushed to
a shared/public branch. Only rebase your own local, unpushed work.

### 5.4 Interactive Rebase (`git rebase -i`)

Allows reordering, squashing, editing, and dropping commits before pushing.

```bash
git rebase -i HEAD~5
# pick, squash, fixup, reword, edit, drop
```

Common workflow: clean up a feature branch before opening a PR.

### 5.5 Phased Integration Model (Recommended)

Modern best practice combines strategies:

1. **During development:** Use rebase to keep feature branch current with main.
2. **Before PR:** Interactive rebase to clean up WIP commits.
3. **At merge time:** Use squash merge or merge commit depending on team policy.

**Decision matrix:**

| Scenario | Strategy |
|---|---|
| Local feature branch, updating from main | `git rebase main` |
| Cleaning up commits before PR | `git rebase -i` |
| Merging PR (simple feature) | Squash merge |
| Merging PR (complex, multi-commit feature) | Merge commit |
| Shared branch integration | Merge commit (never rebase) |
| Hotfix to main | Merge commit |

> Sources: [Mitchell Hashimoto -- Merge vs Rebase vs Squash](https://gist.github.com/mitchellh/319019b1b8aac9110fcfb1862e0c97fb), [Atlassian -- Merging vs Rebasing](https://www.atlassian.com/git/tutorials/merging-vs-rebasing), [DataCamp -- Git Merge vs Git Rebase](https://www.datacamp.com/blog/git-merge-vs-git-rebase), [DEV Community -- Git Workflow Best Practices 2025](https://dev.to/_d7eb1c1703182e3ce1782/git-workflow-best-practices-2025-team-proven-strategies-1eg6)

---

## 6. Cherry-Picking Strategies

### 6.1 When to Cherry-Pick

- Backporting bug fixes from main to release branches.
- Applying hotfixes to production while development continues.
- Extracting specific changes from a feature branch for early release.
- Moving commits between divergent branches that cannot be fully merged.

### 6.2 Techniques

**Single commit:**
```bash
git cherry-pick abc1234
```

**Multiple specific commits:**
```bash
git cherry-pick abc1234 def5678 ghi7890
```

**Commit range (exclusive start):**
```bash
git cherry-pick abc1234..ghi7890   # excludes abc1234
git cherry-pick abc1234^..ghi7890  # includes abc1234
```

**No-commit mode (consolidate multiple picks):**
```bash
git cherry-pick --no-commit abc1234 def5678
# resolve conflicts, then:
git commit -m "backport: apply fixes X and Y from main"
```

### 6.3 Best Practices

- **Use `-x` flag:** Appends "(cherry picked from commit ...)" to the message.
  Essential for audit trails when backporting to release branches.
- **Keep working tree clean** before cherry-picking to avoid partial states.
- **Enable `rerere`** (reuse recorded resolution) for frequently repeated picks:
  `git config --global rerere.enabled true`.
- **For noisy branches:** Create a clean backport branch, rewrite to essentials,
  then range-pick. Slower upfront but produces fewer conflicts.
- **Prefer merge over cherry-pick** when the entire branch content is needed.
  Cherry-pick is for surgical extraction only.

### 6.4 Anti-Patterns

- Cherry-picking the same logical change to many branches without tracking --
  use a backport label and automation instead.
- Cherry-picking merge commits without understanding the `-m` parent flag.
- Using cherry-pick as a substitute for proper branch integration.

> Sources: [Git SCM -- git-cherry-pick](https://git-scm.com/docs/git-cherry-pick), [TheLinuxCode -- Cherry-Pick Multiple Commits](https://thelinuxcode.com/how-i-cherry-pick-multiple-commits-in-git-lists-ranges-and-real-world-backports/), [DataCamp -- Git Cherry-Pick](https://www.datacamp.com/tutorial/git-cherry-pick), [aCompiler -- Mastering Git Cherry Pick 2025](https://acompiler.com/git-cherry-pick/)

---

## 7. Git Hooks and CI Enforcement

### 7.1 Hook Types

| Hook | Trigger | Common use |
|---|---|---|
| `pre-commit` | Before commit is created | Lint, format, secret scan |
| `commit-msg` | After message is written | Validate conventional commits |
| `pre-push` | Before push to remote | Type-check, test, dead code scan |
| `post-merge` | After merge completes | Install deps, rebuild |
| `pre-rebase` | Before rebase starts | Warn about shared branches |
| `post-checkout` | After checkout/switch | Rebuild if deps changed |

If a hook exits with non-zero status, Git cancels the operation.

### 7.2 Hook Managers

**Husky (v9+):**

```bash
npm install --save-dev husky
npx husky init
# Creates .husky/ directory with a sample pre-commit hook
echo 'npx lint-staged' > .husky/pre-commit
echo 'npx --no -- commitlint --edit "$1"' > .husky/commit-msg
```

Husky uses the `prepare` npm script to install hooks on `npm install`. Note: some
CI environments skip the `prepare` script; ensure hooks are not required in CI
(CI should run linters/tests directly).

**Lefthook:**

```yaml
# lefthook.yml
pre-commit:
  parallel: true
  commands:
    lint:
      glob: "*.{js,ts,tsx}"
      run: npx eslint {staged_files}
    format:
      glob: "*.{js,ts,tsx,json,md}"
      run: npx prettier --check {staged_files}

commit-msg:
  commands:
    commitlint:
      run: npx commitlint --edit {1}

pre-push:
  commands:
    typecheck:
      run: npx tsc --noEmit
    test:
      run: npx vitest run --changed
```

Lefthook advantages: single Go binary (no Node dependency), parallel hook execution,
built-in staged file filtering, one YAML config replaces Husky + lint-staged.

**pre-commit (Python):**

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v5.0.0
    hooks:
      - id: trailing-whitespace
      - id: check-yaml
      - id: detect-private-key
  - repo: https://github.com/commitizen-tools/commitizen
    rev: v4.1.0
    hooks:
      - id: commitizen
        stages: [commit-msg]
```

### 7.3 lint-staged Integration

Run linters only on staged files for performance:

```json
{
  "lint-staged": {
    "*.{js,ts,tsx}": ["eslint --fix", "prettier --write"],
    "*.{json,md,yml}": ["prettier --write"]
  }
}
```

> Sources: [DEV Community -- Git Hooks Guide 2025](https://dev.to/arasosman/git-hooks-for-automated-code-quality-checks-guide-2025-372f), [0xDC.me -- pre-commit and lefthook](https://0xdc.me/blog/git-hooks-management-with-pre-commit-and-lefthook/), [Steve Kinney -- Git Hooks with Lefthook](https://stevekinney.com/courses/self-testing-ai-agents/git-hooks-with-lefthook), [jsdev.space -- Git Hooks with Husky](https://jsdev.space/howto/git-hooks-husky/)

---

## 8. Branch Protection and Rulesets

### 8.1 GitHub Branch Protection Rules

Protect critical branches (main, release/*) from unreviewed or untested changes.

**Key settings:**

| Setting | Effect |
|---|---|
| Require pull request before merging | No direct pushes to protected branch |
| Required approving reviews | N approvals before merge (configurable) |
| Dismiss stale reviews | Re-review needed after new commits |
| Require status checks to pass | CI must succeed before merge |
| Require branches to be up to date | Branch must include latest main |
| Require signed commits | Only verified GPG/SSH signatures |
| Require linear history | Only squash or rebase merges |
| Block force pushes | Prevent history rewriting |
| Restrict who can push | Limit direct push access |

### 8.2 GitHub Rulesets (Modern Alternative)

Rulesets replace older branch protection with a stackable, composable model.
Multiple rulesets can target the same branches with additive rules.

**Additional ruleset-only capabilities:**

- **Restrict file paths**: Block commits modifying specific directories.
- **Restrict file extensions**: Block commits with certain file types.
- **Restrict file size**: Prevent large file commits.
- **Require code scanning results**: Block merges with security alerts above threshold.
- **Require deployments to succeed**: Changes must deploy to staging before merge.
- **Required reviewers by file pattern**: Route reviews to specific teams.

### 8.3 Recommended Configuration (Production Branch)

```
Required:
  - Pull request before merging
  - 1+ approving reviews (2 for main in large teams)
  - Dismiss stale reviews on new commits
  - Require status checks: [lint, test, build, security-scan]
  - Require branches to be up to date
  - Block force pushes
  - Require linear history (squash or rebase only)
  - Require code scanning (critical/high severity)

Optional (high-security):
  - Require signed commits
  - Require deployments to succeed
  - CODEOWNERS file with required reviewers
```

> Sources: [GitHub Docs -- Managing Branch Protection](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/managing-a-branch-protection-rule), [GitHub Docs -- Available Rules for Rulesets](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets), [Arnica -- GitHub Branch Protection Guide](https://www.arnica.io/blog/what-every-developer-needs-to-know-about-github-branch-protection), [Hadosec -- 10 Rules of GitHub Branch Protection](https://www.hadosec.com/blog/github-branch-protection/)

---

## 9. Monorepo Workflows

### 9.1 Monorepo Fundamentals

A monorepo houses multiple projects/packages in a single repository. Git workflow
implications include larger diffs, cross-package dependencies, and the need for
affected-based CI to avoid running all tests on every change.

### 9.2 Tooling Comparison

| Tool | Language | Key feature | Best for |
|---|---|---|---|
| Nx | TypeScript | Affected graph, generators | Large enterprise, polyglot |
| Turborepo | TypeScript | Remote caching, simple config | JS/TS projects |
| pnpm workspaces | N/A | Built-in, no extra tooling | Lightweight setups |
| Lerna | TypeScript | Publishing, changelog | Legacy monorepos |
| Bazel | Multi | Hermetic builds, massive scale | Google-scale projects |

### 9.3 Affected-Based CI

Only run tasks for packages changed since the base branch.

**Turborepo:**
```bash
npx turbo run test --filter='...[origin/main]'
```

**Nx:**
```bash
npx nx affected --target=test --base=origin/main --head=HEAD
```

**GitHub Actions with Nx:**
```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for affected detection
      - uses: nrwl/nx-set-shas@v4
      - run: npx nx affected --target=test
      - run: npx nx affected --target=build
      - run: npx nx affected --target=lint
```

**Important:** Both tools need git history to compare changes. The default
`actions/checkout` does a shallow clone (`fetch-depth: 1`). Set `fetch-depth: 0`
for affected detection.

### 9.4 Monorepo CI Patterns

- **Path-based filtering**: Use GitHub Actions `paths` filter to skip workflows
  for unchanged packages.
- **Dynamic matrix**: Generate a matrix of affected packages, run each in parallel.
- **Remote caching**: Turborepo and Nx both support remote caches (Vercel, Nx Cloud)
  to share build artifacts across CI runs and developer machines.
- **Task dependencies**: `turbo.json` and `nx.json` define task pipelines so
  `build` runs before `test` automatically.

### 9.5 Monorepo Branching Considerations

- Use path-scoped CODEOWNERS for per-package review requirements.
- Consider package-scoped conventional commit scopes: `feat(api): ...`, `fix(ui): ...`.
- Release automation tools like Changesets handle per-package versioning natively.
- Avoid GitFlow in monorepos -- trunk-based or GitHub Flow works better.

> Sources: [WarpBuild -- GitHub Actions Monorepo Guide](https://www.warpbuild.com/blog/github-actions-monorepo-guide), [DevToolBox -- Monorepo Tools 2026](https://viadreams.cc/en/blog/monorepo-tools-2026/), [Feature-Sliced Design -- Monorepo Architecture 2025](https://feature-sliced.design/blog/frontend-monorepo-explained), [Nx Blog -- PNPM Workspaces with Nx](https://nx.dev/blog/setup-a-monorepo-with-pnpm-workspaces-and-speed-it-up-with-nx)

---

## 10. Anti-Patterns and Common Mistakes

### 10.1 Branching Anti-Patterns

- **Long-lived feature branches**: Branches older than a week accumulate merge
  conflicts and integration risk. Prefer short-lived branches with feature flags.
- **Rebasing shared branches**: Rewriting history on branches others are using
  forces everyone to re-sync and loses their local commits.
- **No branch protection on main**: Allows accidental force pushes, untested
  code, and unreviewed changes to reach production.
- **Using GitFlow for web apps**: Adds unnecessary complexity for continuously
  deployed services.

### 10.2 Commit Anti-Patterns

- **Giant commits**: Single commits touching dozens of files across multiple
  concerns. Break into logical atomic commits.
- **Meaningless messages**: "fix", "wip", "asdf". Use conventional commits.
- **Mixing formatting with logic changes**: Makes review harder and bisect useless.
  Separate formatting commits from behavioral changes.
- **Committing secrets**: Use pre-commit hooks with tools like `detect-secrets`
  or `gitleaks` to prevent credential leaks.

### 10.3 Merge Anti-Patterns

- **Always squash**: Loses valuable context for complex multi-step changes.
  Use squash for simple PRs, merge commits for complex ones.
- **Never rebasing**: Leads to a tangled history full of merge commits.
  Rebase local branches before opening PRs.
- **Cherry-picking as standard integration**: Cherry-pick is surgical; use
  merge/rebase for standard branch integration.

### 10.4 CI Anti-Patterns

- **Running all tests on every change** in a monorepo. Use affected-based execution.
- **Skipping hooks in CI**: CI should run the same checks as local hooks, but
  directly (not via Husky's prepare script).
- **No status checks on protected branches**: Allows broken code to merge.

---

## 11. Quick Reference Cheat Sheet

### Branching Strategy Selection

```
Q: How often do you deploy?
  Continuously --> Q: Multiple environments?
    No  --> GitHub Flow
    Yes --> GitLab Flow
  Scheduled releases --> Q: Multiple supported versions?
    No  --> GitHub Flow + tags
    Yes --> GitFlow
  Several times daily --> Trunk-Based Development
```

### Common Git Commands by Workflow

```bash
# Start a feature (any workflow)
git checkout -b feat/my-feature main

# Update feature branch (rebase for clean history)
git fetch origin && git rebase origin/main

# Interactive cleanup before PR
git rebase -i HEAD~5

# Squash merge a PR (GitHub CLI)
gh pr merge --squash

# Cherry-pick with provenance
git cherry-pick -x abc1234

# Tag a release
git tag -a v1.2.0 -m "Release 1.2.0"
git push origin v1.2.0

# View branch protection (GitHub CLI)
gh api repos/{owner}/{repo}/branches/main/protection
```

### Conventional Commit Quick Reference

```
feat(scope): add new capability        --> MINOR bump
fix(scope): correct broken behavior    --> PATCH bump
feat!: redesign authentication API     --> MAJOR bump
docs: update installation guide        --> no bump
chore: update dependencies             --> no bump
```

---

## 12. Cross-References

- **devops-infra**: CI/CD pipeline design, GitHub Actions, deployment strategies
- **software-engineering-patterns**: code quality conventions, naming, immutability, design patterns
- **code-packaging**: package versioning, publishing, distribution
- **security-reviewer**: code review and PR review checklists with a security focus

---

## 13. References

### Branching Strategies
1. [Atlassian -- Trunk-Based Development](https://www.atlassian.com/continuous-delivery/continuous-integration/trunk-based-development)
2. [trunkbaseddevelopment.com -- Feature Flags](https://trunkbaseddevelopment.com/feature-flags/)
3. [FeatBit -- TBD Feature Flags 2025](https://www.featbit.co/articles2025/trunk-based-development-feature-flags-2025)
4. [Flagsmith -- TBD Guide](https://www.flagsmith.com/blog/trunk-based-development)
5. [nvie.com -- A Successful Git Branching Model](https://nvie.com/posts/a-successful-git-branching-model/)
6. [Atlassian -- GitFlow Workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow)
7. [git-flow cheatsheet](https://danielkummer.github.io/git-flow-cheatsheet/)
8. [AWS -- Branches in GitFlow](https://docs.aws.amazon.com/prescriptive-guidance/latest/choosing-git-branch-approach/branches-in-a-gitflow-strategy.html)
9. [Ei Square -- Gitflow vs GitHub Flow vs GitLab Flow](https://www.eisquare.co.uk/blogs/how-to-choose-your-branching-strategy)
10. [FUEiNT -- GitHub Flow vs GitLab Flow](https://fueint.com/blog/github-gitlab)
11. [GitKraken -- Git Branch Strategy](https://www.gitkraken.com/learn/git/best-practices/git-branch-strategy)

### Commit Conventions
12. [Conventional Commits v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/)
13. [AverageDevs -- Conventional Commits in Git](https://www.averagedevs.com/blog/conventional-commits-git)
14. [ShakilTech -- Conventional Commits Complete Guide](https://blog.shakiltech.com/conventional-commits/)
15. [Jeff Bailey -- What Are Conventional Commits](https://jeffbailey.us/blog/2025/09/28/what-are-conventional-commits/)

### Semantic Versioning and Release Automation
16. [semantic-release GitHub](https://github.com/semantic-release/semantic-release)
17. [MerginIT -- semantic-release Guide 2025](https://merginit.com/blog/29062025-automated-multi-platform-releases)
18. [Oleksii Popov -- NPM Release Automation](https://oleksiipopov.com/blog/npm-release-automation/)
19. [AWS -- Semantic Versioning](https://aws.amazon.com/blogs/devops/using-semantic-versioning-to-simplify-release-management/)

### Merge Strategies
20. [Mitchell Hashimoto -- Merge vs Rebase vs Squash](https://gist.github.com/mitchellh/319019b1b8aac9110fcfb1862e0c97fb)
21. [Atlassian -- Merging vs Rebasing](https://www.atlassian.com/git/tutorials/merging-vs-rebasing)
22. [DataCamp -- Git Merge vs Git Rebase](https://www.datacamp.com/blog/git-merge-vs-git-rebase)
23. [DEV Community -- Git Workflow Best Practices 2025](https://dev.to/_d7eb1c1703182e3ce1782/git-workflow-best-practices-2025-team-proven-strategies-1eg6)

### Cherry-Picking
24. [Git SCM -- git-cherry-pick](https://git-scm.com/docs/git-cherry-pick)
25. [TheLinuxCode -- Cherry-Pick Multiple Commits](https://thelinuxcode.com/how-i-cherry-pick-multiple-commits-in-git-lists-ranges-and-real-world-backports/)
26. [DataCamp -- Git Cherry-Pick](https://www.datacamp.com/tutorial/git-cherry-pick)
27. [aCompiler -- Mastering Git Cherry Pick 2025](https://acompiler.com/git-cherry-pick/)

### Git Hooks
28. [DEV Community -- Git Hooks Guide 2025](https://dev.to/arasosman/git-hooks-for-automated-code-quality-checks-guide-2025-372f)
29. [0xDC.me -- pre-commit and lefthook](https://0xdc.me/blog/git-hooks-management-with-pre-commit-and-lefthook/)
30. [Steve Kinney -- Git Hooks with Lefthook](https://stevekinney.com/courses/self-testing-ai-agents/git-hooks-with-lefthook)
31. [jsdev.space -- Git Hooks with Husky](https://jsdev.space/howto/git-hooks-husky/)

### Branch Protection
32. [GitHub Docs -- Branch Protection Rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/managing-a-branch-protection-rule)
33. [GitHub Docs -- Available Rules for Rulesets](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets)
34. [Arnica -- GitHub Branch Protection Guide](https://www.arnica.io/blog/what-every-developer-needs-to-know-about-github-branch-protection)
35. [Hadosec -- 10 Rules of GitHub Branch Protection](https://www.hadosec.com/blog/github-branch-protection/)

### Monorepo Workflows
36. [WarpBuild -- GitHub Actions Monorepo Guide](https://www.warpbuild.com/blog/github-actions-monorepo-guide)
37. [DevToolBox -- Monorepo Tools 2026](https://viadreams.cc/en/blog/monorepo-tools-2026/)
38. [Feature-Sliced Design -- Monorepo Architecture 2025](https://feature-sliced.design/blog/frontend-monorepo-explained)
39. [Nx Blog -- PNPM Workspaces with Nx](https://nx.dev/blog/setup-a-monorepo-with-pnpm-workspaces-and-speed-it-up-with-nx)
