---
name: prompt-cost-estimator
description: Estimate how much a prompt will cost across Claude models (Fable 5, Opus 4.8/4.7/4.6, Sonnet 5/4.6, Haiku 4.5) and effort levels (low/medium/high/xhigh/max) before running it, as a model x effort cost matrix from exact token counts (offline heuristic fallback). Also recommends the optimal model/effort via a free no-network heuristic, and, only when asked, executes the prompt for real (y/N confirmation) against that pick or the cheapest option. TRIGGER whenever the user asks "how much will this prompt cost", "cost on Opus vs Sonnet", "estimate token cost", "compare model pricing", "cost per 1000 calls", "what model/effort should I use for this", "run this and tell me what it cost", or wants any price/spend/budget estimate for a Claude call, batch job, or agent loop, even without the word "cost". SKIP for billing/usage history (Anthropic Console), AWS infra costs (aws-startup-advisor), Claude API coding/migration (claude-api), or improving prompt content (prompt-deep-optimizer).
version: 1.3.0
updated: 2026-07-05
model: claude-sonnet-4-6
effort: medium
---

# Prompt Cost Estimator

Estimate the dollar cost of running a given prompt across every current Claude model and effort level, as a matrix — the default mode never spends anything. It can also recommend which model/effort actually fits the task, and, only when asked, run the prompt for real.

## How to run

The bundled script does everything. Script paths below are relative to this skill's base directory (given when the skill loads) — prefix them with it. Point the script at the prompt:

```bash
python3 <skill-dir>/scripts/estimate_cost.py --file prompt.txt
python3 <skill-dir>/scripts/estimate_cost.py --text "Summarize this repo's architecture"
python3 <skill-dir>/scripts/estimate_cost.py --tokens 12000   # already know the input count
python3 <skill-dir>/scripts/estimate_cost.py --file prompt.txt --quick   # fast, free, one-line ballpark
python3 <skill-dir>/scripts/estimate_cost.py --file prompt.txt --confirm-exec   # estimate, then ask y/N to actually run it
python3 <skill-dir>/scripts/estimate_cost.py --file prompt.txt --exec          # estimate, then run it immediately
python3 <skill-dir>/scripts/estimate_cost.py --file prompt.txt --recommend     # cheap heuristic: what model/effort fits this task
python3 <skill-dir>/scripts/estimate_cost.py --file prompt.txt --price-only    # bare dollar figure, nothing else — for scripts/hooks
```

### Quick, free estimate

When the user just wants a fast ballpark ("roughly how much", "is this cheap") and does not need the full matrix, use `--quick`. It never calls the network (free — no API key, no latency), picks the cheapest model at its lowest effort, and prints a single line:

```
~$0.01  (Haiku 4.5)
  14 in + ~2,000 out tok  [heuristic ...; rough — use full mode for a real quote]
```

`--quick` implies `--offline` and combines with `--calls`, `--models` (override the cheapest-model default), `--output-tokens`, and `--json`. Drop it and run the full table when the user is comparing models/efforts or needs an accurate quote.

### Recommending the optimal model/effort

`--recommend` is a second, independent heuristic from the cost matrix: cheap (pure keyword + length matching on the prompt text, zero network calls) and fast, it answers "what should this actually run on" rather than "what's cheapest". It classifies the prompt into one of five tiers (mechanical/routine/analytical/long_horizon/frontier — same vocabulary as skill-optimizer's own Step 4.6 model-recommendation table) and prints the pick with a one-line rationale:

```
Recommended: Opus 4.8 @ high  [tier: analytical]
  why: review/design/diagnosis-style task needs real judgment
  to switch this session: /model claude-opus-4-8  /effort high
```

This is a heuristic, not a benchmark — say so if the user pushes on precision. It does not apply under `--quick` (that mode is deliberately "cheapest, no thinking") or with `--tokens` (needs the actual prompt text). If `--models` excludes the tier's default model, it substitutes the cheapest model still in scope and says so in the rationale.

**On "automatically switching":** the script cannot reach into a live Claude Code session and change its model — no tool does that. What it prints is the literal `/model`/`/effort` command pair; run those yourself (or ask the assistant to) if you want the *current* session to switch. Where "automatic" is real: `--exec`/`--confirm-exec` target the recommendation by default (see below), so the actual API call this script makes uses the optimal pick without you specifying it.

### Executing the prompt after the estimate

Two flags turn the estimate into an action:

- `--confirm-exec` — after printing the estimate, ask `Execute now with <model> @ <effort> ...? [y/N]:` and only run the prompt for real if the user answers yes. If stdin isn't a TTY (piped/non-interactive), it prints the question and skips rather than blocking.
- `--exec` — skip the question and run the prompt immediately.

Both target the `--recommend` pick by default (`--target recommended`); pass `--target cheapest` to restore pure cost-minimization instead. Both require actual prompt text (`--file`/`--text`, not `--tokens`) and real API credentials — this is a genuine, billed `messages.create` call, not a simulation. **Treat `--exec` as a spend-incurring action**: don't pass it on the user's behalf unless they asked to actually run the prompt, not just price it. On success the script prints the response plus actual token usage compared against the estimate; on failure (no credentials, bad model ID, an `anthropic` package too old to accept the `effort` kwarg, network error) it prints the error and leaves the estimate output untouched.

### Bare-number output for scripting

`--price-only` prints nothing but the dollar figure for `--target` (e.g. `$0.0031`) — no table, no recommendation, no exec, implies `--offline`. It exists for callers that need a single number with no parsing: hooks, status lines, other scripts. Use it whenever the caller is code, not a person reading the response.

### Automatic per-prompt recommendation (hook + status line)

`scripts/prompt_cost_hook.py` wraps the same free `--recommend` heuristic as a Claude Code `UserPromptSubmit` hook: on every prompt it computes the recommended model/effort and a price estimate (no network calls, milliseconds), writes `~/.claude/cost-estimator/<session_id>.json`, and emits the recommendation as `additionalContext` so it's visible at the start of every turn. A status line can read that state file to show the estimate below the prompt. This is the practical form of "auto-switch": no hook or tool can force-change a live session's model, so the hook surfaces the recommendation for the assistant or user to act on via `/model`/`/effort` — it does not flip the switch itself. Wiring this into `~/.claude/settings.json` and a status line script is a one-time setup step outside this skill's own files; read `scripts/prompt_cost_hook.py`'s docstring before wiring it, since it writes to a path outside the skill directory.

Useful knobs:

| Flag | Purpose |
| --- | --- |
| `--output-tokens N` | Expected base output at low effort (default 2000). Ask the user for a realistic value when the task shape is known — short classification (~100), prose/report (~2000), code generation (~5000+). |
| `--models a,b` | Restrict to specific model IDs |
| `--calls N` | Scale to N calls (batch jobs, agent loops) |
| `--cached` | Price input at cache-read rate (~0.1x) — use for repeated calls with a stable prefix |
| `--multipliers '{"high": 3.0}'` | Override effort multipliers with measured data |
| `--json` | Machine-readable output |
| `--offline` | Skip the count_tokens API, use the character heuristic |
| `--quick` | Fast, free one-line ballpark: implies `--offline`, cheapest model at its lowest effort (see below) |
| `--recommend` | Cheap, free, no-network heuristic: pick the optimal model/effort for this prompt and print the `/model`/`/effort` commands to switch to it |
| `--target {recommended,cheapest}` | Which row `--exec`/`--confirm-exec` acts on (default `recommended`) |
| `--confirm-exec` | After the estimate, ask y/N before actually running the prompt against `--target` |
| `--exec` | Skip the question, run the prompt immediately against `--target` — spend-incurring, see below |
| `--price-only` | Print only the dollar figure for `--target`, nothing else; implies `--offline` — for scripts/hooks |

## Token counting

The script tries the Anthropic `count_tokens` API first (exact, free, needs credentials — `ANTHROPIC_API_KEY` or an `ant auth login` profile). If unavailable it falls back to ~3.5 chars/token and labels the result as heuristic. Always tell the user which source was used; the heuristic can be off by ±20%, more on code or non-English text.

## Interpreting the matrix

- **Not every model supports every effort level.** The script already omits unsupported cells (Haiku 4.5 takes no effort parameter — shown as "-"; Opus 4.6 and Sonnet 4.6 lack `xhigh`; see the Effort levels column in `references/pricing.md`). Never invent a price for a missing cell.
- **Effort multiplies output spend, not input.** Effort levels (low → max) increase thinking + response tokens. The multipliers (1x/1.5x/2.5x/4x/6x) are heuristics — real spend varies by task. Say so when presenting results.
- **Input cost is usually the small term** for short prompts; for large-context prompts (100K+ tokens) input dominates and model choice matters more than effort.
- **Cache changes everything for repeated calls.** With a stable prefix, input re-reads cost ~0.1x. Show both `--cached` and uncached numbers when the user is planning a loop or batch.
- **Batch API is 50% off everything** — mention it when the workload is not latency-sensitive.
- **Sonnet 5 intro pricing** ($2/$10 per MTok) applies through 2026-08-31; the script handles this automatically and marks the row "(intro)".

## Presenting results

Show the full table, then lead with a one-line answer to the actual question: cheapest viable option, and cost at the model/effort the user is most likely to use. If the user gave no output estimate, state the assumption (2000 base output tokens) and offer to re-run with theirs. In `--quick` mode there is no table — relay the single-line ballpark and note it is rough, offering a full run for a precise quote. When the user asks which model/effort they should actually use (not just the cost), run `--recommend` and relay the tier + rationale, not just the model name — the "why" is what makes it actionable. Only pass `--exec` when the user explicitly asked to run the prompt, not merely to price it; when unsure, use `--confirm-exec` so the user gets a final yes/no before anything is billed.

Pricing is a cached snapshot (see `references/pricing.md` for the table and date). If precision matters or the user asks for "current" pricing, verify against https://platform.claude.com/docs/en/pricing.md before quoting numbers.
