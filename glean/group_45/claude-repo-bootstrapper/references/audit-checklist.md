# Repo Bootstrapper — Audit Checklist

**Before anything else: read the audit ledger.** If `docs/repo-bootstrap-audit-*.md` files exist, read the most recent one first. Findings marked RETRACTED there are off-limits without new evidence. Findings marked Deferred are your starting backlog. When parallel subagents report gaps, verify each against the artifact (grep/read it) before remediating — agents systematically over-report (counting one block of a chunked index as the total; reading "coverage gap marks" as "undocumented call").

When auditing an existing repo, check each file. Severity rubric:
- **blocking** = file is referenced elsewhere but missing
- **major** = file present but missing a required section
- **medium** = section present but stale relative to the codebase
- **minor** = formatting / naming drift

## Workflow files

**`.github/copilot-instructions.md`**
- [ ] `## Default Execution Strategy` block is first and contains these 4 rules: (1) read relevant files before editing, (2) run tests after every change, (3) never invent commands or paths, (4) follow the workflow log rule (prompts.md + memory.md + version bump)
- [ ] `## Build, Test, and Validation Commands` section exists with runnable commands
- [ ] `## High-level Architecture` section exists
- [ ] `## Key Conventions` section exists with ≥6 rules

**`CLAUDE.md`**
- [ ] `## Repository shape` section exists
- [ ] `## Commands` section matches copilot-instructions commands
- [ ] Workflow log rule present (prompts.md + memory.md + version bump)

**`AGENTS.md`**
- [ ] Table of every agent under `.claude/agents/` or `.github/agents/` with name, scope (account-vs-global), model, when-to-use
- [ ] Cross-link to the agent files
- [ ] Notes when an agent depends on env vars or external auth

**`GEMINI.md`**
- [ ] Either contains Gemini-specific rules OR points readers at `CLAUDE.md` with a note that conventions are shared
- [ ] Lists any Gemini-only tool mappings (e.g., `activate_skill`)

**`memory.md`**
- [ ] Versioned format (`## vN - date`)
- [ ] Current version matches the repo's canonical version file
- [ ] Active task / completed / next steps all present in latest entry
- [ ] File under ~200 KB; if larger, **medium** — rotate older sections to `docs/archive/` via `scripts/rotate-workflow-logs.mjs` (never rotate while an editor swap file like `.memory.md.swp` is live; the tool refuses)
- [ ] Same size check applies to `prompts.md`

## Dotfile meta

**`.editorconfig`**
- [ ] `root = true` at top
- [ ] At minimum: `indent_style`, `indent_size`, `end_of_line`, `charset`, `trim_trailing_whitespace`, `insert_final_newline`

**`.gitignore`**
- [ ] Excludes `node_modules/`, `.env`, `.env.local`, `.DS_Store`, build outputs, IDE caches
- [ ] Does NOT exclude any of the manifest files
- [ ] Repo-specific paths covered (per language / framework)

**`.gitattributes`**
- [ ] `* text=auto eol=lf` for cross-platform consistency
- [ ] Binary markers on non-text assets

**`.nvmrc` / `.node-version` / `.tool-versions`** *(Node repos only)*
- [ ] Present when the repo runs Node (or any pinned toolchain)
- [ ] Matches `package.json` engines field

**`.env.example`**
- [ ] Lists every env var read by the codebase (cross-check via `grep -rE "process.env.[A-Z_]+" src/`)
- [ ] No real secrets — placeholder values only
- [ ] Each var has a one-line comment describing its purpose

**`.vscode/`**
- [ ] `settings.json` configures the project's formatter + linter
- [ ] `extensions.json` recommends the formatter, language servers, MCP integration
- [ ] `launch.json` defines at least one debug config per runtime
- [ ] `mcp.json` mirrors `.mcp.json` (shape may vary by editor version)

**`.mcp.json`**
- [ ] Lists every MCP server the repo's agents depend on
- [ ] Each entry has `command` or `httpUrl` plus, where applicable, `env`

**`.github/workflows/`**
- [ ] At least one workflow runs tests on PR
- [ ] At least one workflow runs lint/type-check on PR
- [ ] Workflow files use pinned action versions (sha-pinned in security-sensitive repos)

**`.github/dependabot.yml`**
- [ ] Configured for the package ecosystems the repo uses
- [ ] Sets a schedule

**`.github/CODEOWNERS`**
- [ ] Path-based assignments cover at least the top-level directories
- [ ] No stale handles

**`.github/PULL_REQUEST_TEMPLATE.md`**
- [ ] Sections for: What changed / Why / Tests / Risk / Rollback

**`.github/ISSUE_TEMPLATE/`**
- [ ] At minimum `bug_report.md` and `feature_request.md`

**`.github/SECURITY.md`**
- [ ] How to report vulnerabilities
- [ ] Points at `docs/SECURITY.md` for the threat model
- [ ] Response-time SLA

**`LICENSE`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`**
- [ ] Present
- [ ] `CONTRIBUTING.md` links to `docs/DEVELOPMENT.md`

## Documentation suite

**`docs/ARCHITECTURE.md`**
- [ ] Context section covers external actors and integrations
- [ ] Container/component diagram present (ASCII or linked image)
- [ ] At least one ADR documented

**`docs/TESTING.md`**
- [ ] All test suites listed with exact run commands
- [ ] Coverage targets stated
- [ ] CI gates described

**`docs/external-calls.md`** *(load-bearing for the auto-remediation contract)*
- [ ] One row per external call (HTTP, MCP, child_process, WebSocket, etc.)
- [ ] Each row has: file:line, call target, logger handle, test file, doc paragraph, retry policy
- [ ] A row missing two or more of (log, test, docs, retry) is at least **medium** in the audit report

**`docs/runbooks/`**
- [ ] One file per recurring operational procedure (incident response, rollback, on-call handoff)
- [ ] Each runbook follows: precondition → action steps → verification → escalation

## Operations registry

**`server/src/lib/operations-registry.js`** (canonical in `10gen/mdb-tam`; older repos may carry `mcp-server/src/operations-registry.ts`)
- [ ] Every external call (HTTP, MCP, child_process, file I/O, fs.watch) has an `OperationEntry` in `OPS_REGISTRY`
- [ ] Each entry has: `id`, `kind`, `transport`, `target`, `description`, `sourceFile`, `readOnly`, `verifyDatastore`, `retryPolicy`
- [ ] **Entry truthfulness:** `transport`/`target`/`description` match what the source file actually does — read the sourceFile; a Google-Docs feed registered as `transport: 'local' / target: 'local-store'` is a **major** finding
- [ ] **Code coverage:** every error code thrown by the sourceFile (`grep "code: '"` it) appears in the entry's error-code/remediation map; codes thrown but unmapped are **medium**
- [ ] Each entry has at least one `remediation` with `code`, `hint`, and `safety` (repo vocabulary may be `safe-to-retry` / `not-retriable` / `skip-project` — auth/config failures must NOT be retry-safe)
- [ ] Operations with `verifyDatastore: true` include a `DATASTORE_VERIFY_FAILED` error code — **`kind: sync` only; do not demand it for `kind: check`/`report` ops** (retracted finding M3, 2026-06-01)
- [ ] `listOperations`, `getOperation`, `runOperation`, `auditOperations`, `findRemediationByCode` are exported
- [ ] Registry is wired into the service layer (re-exports) and server (MCP tool registrations)
- [ ] `scripts/generate-ops-registry-doc.mjs` is the sole writer of `docs/operations-registry.json`; `ops:doc:check` runs in CI

**`docs/operations-registry.json`**
- [ ] Generated by `writeRegistryDoc()` or equivalent `tam_ops_write_doc` tool
- [ ] Contains `operations`, `audit`, and `status` sections
- [ ] Audit section shows per-operation pass/fail against the 5-standard contract

**`docs/tool-inventory.json`**
- [ ] Generated from runtime `instrumentedRegisterTool` metadata, not hand-maintained
- [ ] `totalTools` matches the actual count of registered tools
- [ ] Each entry has `name`, `title`, `description`, `annotations`, `kind`, `domain`
- [ ] CI check verifies the committed file matches the live server (count + name diff)

**CI drift check** (in `.github/workflows/ci.yml` or equivalent)
- [ ] After tests pass, a step boots the server, builds the live tool inventory, and compares against the committed `docs/tool-inventory.json`
- [ ] Fails with a clear message naming which tools are missing or extra
- [ ] Error message tells the developer how to regenerate

## Index hygiene

**`docs/high_signal_file_index.json` + `docs/llm-repo-index.json`**
- [ ] `node scripts/check-doc-indexes.mjs` (or equivalent) passes: zero dead paths in either index — dead entries make LLM retrieval worse than no index (observed: 90% dead, 4379/4857 entries, 2026-06-09)
- [ ] The checker is wired into a CI job so index rot fails the build
- [ ] `llm-repo-index.*` is regenerated (never hand-edited): `python3 scripts/generate_llm_repo_index.py`
- [ ] `high_signal_file_index.json` is curated: prune dead entries with `--prune`, then add entries for significant new files following the existing entry schema; note that the file is a LIST OF BLOCKS each holding `files[]` — count entries at the file level, not the block level
- [ ] New workspace components (e.g. a new `packages/*`) appear in BOTH indexes, `docs/codebase-overview.md`, `docs/COMPONENTS.md`, `CLAUDE.md` repository-shape, and `docs/ARCHITECTURE.md`

## Skill files

For each `docs/*-context.md`:
- [ ] Corresponding skill exists at `~/.claude/skills/<id>/SKILL.md`
- [ ] `.skillfish.json` (at `~/.claude/skills/<id>/.skillfish.json`) has `source: "manual"` and a `sha` field containing the SHA-256 of that skill's `SKILL.md` file at the time it was last synced

## Extended checks (run per-repo-type)

| Check | Run when |
|---|---|
| Agent inventory, coding patterns, concept tree, MCP health, dependency audit, skill coverage | All repos |
| Customer data in VCS | All repos (TAM repos especially) |
| Native host manifest validation | Repo contains `native-host/` directory |
| Offscreen document audit | Repo is a Chrome MV3 extension (`"manifest_version": 3` in `manifest.json`) |
| Auth surface inventory | Repo has OAuth, cookie-based auth, API tokens, or a vault |

### Customer data in VCS
- [ ] `git ls-files` scan for customer-named directories/files (account names, Slack/meeting exports, CRM dumps, `*_context_modules_*.json`, initiative-board XLSX)
- [ ] Tracked customer data is at least a **medium** finding: surface it with a `git rm --cached` + `.gitignore` recommendation and let the operator decide — never auto-delete or untrack customer files yourself
- [ ] Check generated indexes don't leak customer content (filenames in indexes are acceptable for private repos; flag anything beyond filenames)

### Agent inventory (`AGENTS.md`)
- [ ] Every `.md` file in `.claude/agents/` has a corresponding row in `AGENTS.md`
- [ ] Each row includes: name, description, model, when-to-invoke, and required env/auth
- [ ] AGENTS.md is auto-generated from agent frontmatter — not hand-maintained

### Coding patterns registry sync
- [ ] Extract patterns matching coding-patterns catalog categories (error-handling, retry, caching, auth, validation, async, data-access, performance)
- [ ] For each pattern, check if a matching entry exists in `coding-patterns/registry.json` in the hub
- [ ] Flag new patterns for `tam_coding_pattern_save` and stale patterns for update

### Concept tree sync
- [ ] Verify the repo's primary domains are represented in `concept-tree/tree.json`
- [ ] Flag any domain covered by the repo but missing from the concept tree
- [ ] Check concept `researchedAt` dates — flag concepts older than 90 days as stale

### MCP server health check
- [ ] If the repo ships an MCP server, run `--self-test` and verify it boots
- [ ] Enumerate registered tools and compare count against the last known value in `docs/MCP.md`
- [ ] Flag any tool mentioned in docs but not registered, or registered but not documented

### Dependency audit
- [ ] Run `npm audit --json` (or equivalent) and flag critical/high vulnerabilities
- [ ] Run `npm outdated --json` and flag major version bumps available
- [ ] Check that `package-lock.json` exists and is committed
- [ ] Verify `engines.node` matches `.nvmrc` / `.node-version`

### Native host manifest validation (Chrome extension repos)
- [ ] For each native host in `native-host/`, verify the JSON manifest is valid
- [ ] Check that `allowed_origins` contains the current unpacked extension ID
- [ ] Verify the `path` field points to an existing executable
- [ ] Check that `install.sh` exists and is executable

### Offscreen document audit (Chrome MV3 repos)
- [ ] `Reason` enum value is valid and matches actual use
- [ ] Lifetime management: created on demand, closed when no longer needed
- [ ] Communication uses `chrome.runtime.sendMessage`
- [ ] Flag offscreen documents that are created but never closed (leak risk)

### Skill coverage gap analysis
- [ ] List the repo's technology stack (package.json deps, imports, manifest permissions)
- [ ] Compare against installed skills at `~/.claude/skills/`
- [ ] Flag any technology used by the repo with no corresponding skill

### Auth surface inventory (Chrome extension repos)
- [ ] For each auth surface (OAuth, cookies, API tokens, vault): documented, has a health-check mechanism, failure mode documented
