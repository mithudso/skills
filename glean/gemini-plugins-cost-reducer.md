# cost-reducer

**Category:** AI, Agents & Prompt Engineering
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/cost-reducer/skills/cost-reducer

## Description
Use when auditing Claude Code or agent setup for cost reduction opportunities. Analyzes SKILL.md and CLAUDE.md files for bloat, identifies hook automation opportunities, flags expensive model choices, and surfaces parallelization gaps.

---

# Agent Cost Reducer

Audit a Claude Code or agent-skills setup and produce a prioritized list of cost-reduction recommendations across four dimensions: instruction bloat, missing hook automation, model right-sizing, and parallelization gaps.

## Required Inputs

You can run this against:
- **A specific repo** — analyze its `.claude/` directory and any `SKILL.md` / `CLAUDE.md` files
- **This marketplace repo** — analyze `.agents/skills/` for instruction bloat and model choices
- **A user's global config** — analyze `~/.claude/` for hook opportunities

Ask the user which scope to use if not specified.

---

## Step 1: Inventory

Collect the files you'll analyze:

```bash
# CLAUDE.md files in repo
find . -name "CLAUDE.md" -not -path "*/node_modules/*"

# Agent instruction files
find . -name "SKILL.md" -not -path "*/node_modules/*"

# Installed agents
ls ~/.claude/agents/ 2>/dev/null

# Existing hooks
cat ~/.claude/settings.json 2>/dev/null | python3 -m json.tool
cat .claude/settings.json 2>/dev/null | python3 -m json.tool

# Recent conversation transcripts (hook opportunities)
ls ~/.claude/projects/ 2>/dev/null | head -5

# Check if caveman is installed at the project level
claude plugin list --scope project 2>/dev/null | grep -q caveman && echo "caveman: project-installed" || echo "caveman: not project-installed"
```

Read each file and build a mental inventory before proceeding to any analysis step.

---

## Step 2: Instruction Bloat (SKILL.md / CLAUDE.md)

For each instruction file, flag:

### 2a. Inline content that belongs in `references/`

Content that belongs in `references/` instead of inline:
- Tables of known values (distro names, status codes, enum lists)
- Example outputs longer than ~10 lines
- Repo-specific paths or commands that only apply to one project

Mark each as: **"Move to references/ — saves ~N tokens"**

### 2b. Redundant or obvious instructions

Flag instructions that restate default model behavior or add no constraint:
- "Read the file before editing" (Claude Code already does this)
- "Make sure the code compiles" (no mechanism to enforce it differently)
- "Be thorough" / "Think carefully" (filler that costs tokens without changing behavior)

Mark each as: **"Remove — no behavioral delta"**

### 2c. Structural inefficiency

- `description` fields longer than 2 sentences — truncate to the trigger condition + one-line outcome
- Duplicate guidance that appears in multiple files — consolidate into one place
- `references/` files that are loaded eagerly when they could be loaded on demand

---

## Step 3: Hook Automation Opportunities

Scan for patterns worth automating as hooks instead of re-prompting each session.

### Common patterns to look for

| Pattern | Hook type | Example |
|---------|-----------|---------|
| "After every file edit, run the linter" | `PostToolUse` on `Edit`/`Write` | `eslint --fix $FILE` |
| "Before committing, run tests" | `PreToolUse` on `Bash` matching `git commit` | `npm test` |
| "After tool use, log to Honeycomb" | `PostToolUse` on `mcp__.*__skill.*` | telemetry command |
| "When Claude stops, show summary" | `Stop` hook | notification script |
| "Validate SKILL.md on every save" | `PostToolUse` on `Write` matching `SKILL.md` | `skill-validator check` |

### How to identify them

Look for:
1. Instructions in CLAUDE.md that say "always X before Y" or "every time you Z, also W"
2. Steps that appear verbatim in multiple skill workflows
3. Manual setup steps the user runs at the start of each session (from recent transcripts if available)

For each opportunity, output:

```
Hook opportunity: <description>
Type: PreToolUse | PostToolUse | Stop
Matcher: <tool name or regex>
Command: <shell command>
Estimated sessions saved: high | medium | low
```

---

## Step 4: Model Right-Sizing

For each agent (in `~/.claude/agents/` or `.agents/skills/*/agents/`), check the `model:` field.

### Downgrade candidates

Flag agents using `opus` or `inherit` where `sonnet` or `haiku` would suffice:

| Agent characteristics → recommended model |
|---|
| Pure file reading + summarization → `haiku` |
| Search, grep, or glob only → `haiku` |
| Simple single-step transformations → `haiku` |
| Multi-step reasoning, code generation, debugging → `sonnet` |
| Deep architectural analysis, complex trade-offs → `opus` |

Flag `inherit` unconditionally — it ties cost to the parent session model, which is non-deterministic.

Output for each candidate:

```
Agent: <name>
Current model: opus | inherit
Recommended: sonnet | haiku
Reason: <one line>
Estimated saving: ~Nx per invocation (opus→sonnet ≈ 5x, sonnet→haiku ≈ 5x)
```

---

## Step 5: Parallelization Gaps

Look for sequential agent calls or tool chains that have no dependency between steps.

### Where to look

1. Skill workflows that call multiple subagents in sequence (e.g., "first run agent A, then run agent B with the result of A") — if A's output isn't actually needed by B, they can run in parallel
2. CLAUDE.md instructions that say "first do X, then Y" where Y doesn't depend on X's output
3. Skills that read multiple independent files one at a time instead of in a single parallel batch

### Output format

```
Parallelization gap: <location>
Current: sequential calls to [A, B, C]
Can parallelize: [A, B] (C still needs A's output)
Estimated latency saving: ~Nx wall time
```

---

## Step 6: Report

Produce a prioritized report sorted by estimated impact (tokens saved × frequency):

```markdown
## Agent Cost Reduction Report

### Quick Wins (low effort, immediate saving)
1. **[Hook]** Automate lint-on-save — saves re-prompting ~N times/day
2. **[Model]** Downgrade `<agent>` from opus to haiku — 25x cheaper per call
3. **[Bloat]** Remove 3 filler instructions from CLAUDE.md — saves ~150 tokens/session

### Medium Effort
4. **[Bloat]** Move distro table from `variant-right-sizer` inline to `references/` — saves ~400 tokens when skill isn't active
5. **[Parallel]** Fan out [A, B] in `pr-review-loop` — saves ~40s wall time per run

### Longer Term
6. **[Bloat]** Consolidate duplicate "submit patch" guidance across 3 Evergreen skills into shared `references/evergreen-common.md`

### Complementary Tools
- **[`/caveman`](https://www.skills.sh/juliusbrussee/caveman/caveman)** — compresses Claude's response verbosity by ~75% for the session. Complements structural fixes above with session-level savings. Only include this if caveman is not project-installed (i.e. `caveman: not project-installed` from the Step 1 inventory check). Include it even if the user has caveman globally — the recommendation is to add it to the repo's project setup.

### Summary
| Dimension | Issues found | Est. token saving/session |
|-----------|-------------|--------------------------|
| Instruction bloat | N | ~X tokens |
| Hook automation | N | ~Y re-prompts |
| Model right-sizing | N agents | ~Zx cost multiplier |
| Parallelization | N gaps | ~Ws wall time |
```

---

## Tips

- Focus on the highest-token skills first — run `wc -w` across all SKILL.md files to rank them
- The combined SKILL.md token budget across this marketplace is 10,000 tokens; CI warns when exceeded
- `haiku` is the right default for any agent that only reads and summarizes — it's 25x cheaper than `opus` and the quality difference is negligible for retrieval tasks
- A well-targeted `PostToolUse` hook beats a CLAUDE.md instruction every time — hooks run unconditionally; instructions can be missed
- For session-level token savings, recommend the [`caveman` skill](https://www.skills.sh/juliusbrussee/caveman/caveman) — it compresses Claude's responses by ~75% by stripping filler language. Complementary to this skill's structural recommendations.