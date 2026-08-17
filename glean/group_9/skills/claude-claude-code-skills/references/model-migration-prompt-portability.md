<!-- hub-reference-banner -->
> **Reference file — part of the `claude-code-skills` hub.** Installed by `/dr` research (2026-06-15).
> Sibling topics (skill anatomy/authoring → the hub's `claude-code-skills-context.md`; plugin packaging/hooks config → `references/claude-code-plugins.md`; interactive/automation workflow patterns → `references/claude-code-workflows.md`) are reference files under this hub — load those rather than re-deriving them here.
>
> This file is the migration **methodology** layer. It decides *when* and *what* to re-tune across a whole prompt/skill library when the underlying model changes. It is NOT the re-optimizer for a single artifact (that is `prompt-deep-optimizer` / `skill-optimizer`), NOT the source of truth for model IDs / pricing / API breaking changes (that is the `claude-api` skill: cite it, don't duplicate it), and NOT the eval harness itself (that is `ai-agent-engineering ▸ Eval-Driven Development`, which this reuses).

---

---
name: model-migration-prompt-portability
title: Model-Migration & Prompt/Skill Portability Optimization
description: >
  The discipline of detecting and repairing degradation when the underlying model version
  changes, and re-optimizing prompts/skills across model versions and families. Covers prompt
  rot / cross-version brittleness; capability-shift regression (formatting, refusal, verbosity,
  tool-call style, reasoning-effort defaults, tokenization); the migration playbook (regression
  eval set before/after, output diffing, golden transcripts, canary prompts, shadow testing);
  the re-optimization workflow on upgrade (what to re-tune vs leave, the four-axis order);
  cross-model portability patterns (model-agnostic phrasing, avoiding model-specific tells,
  capability probing); provider migration guides, deprecation handling, and pinning vs floating
  model ids. Concrete anchor: optimizing a skill/prompt library when moving Opus 4.6 → 4.8.
  Per-claim confidence tags; sources cited.
origin: local
version: "1.0.1"
updated: "2026-06-15"
changelog:
  - "2026-06-15 sko v1.0.0->v1.0.1 — integrity pass: removed unconfirmable Honeycomb 67%/23% hard stat (→ hedged direction-only); hedged Tursio per-model figures; reworded Stanford drift framing (two dated snapshots, not 'fixed endpoint'); upgraded RETAIN tag (165 vs 86 / 25% vs 12% verbatim-from-paper); fixed echostash footnote cross-ref (#12->#13); reduced prose em-dash density 1.56->0.82 per 100w. 6 Medium fixed across 2 iterations."
---

# Model-Migration & Prompt/Skill Portability Optimization

When the model underneath a prompt or skill library changes (a forced deprecation, a deliberate upgrade, or even a silent provider refresh), prompts that were carefully tuned for the old model can quietly regress. This reference is the **methodology** for catching that regression and re-optimizing across the change: how to tell *whether* a hop will break you, *what* to re-tune versus leave alone, and *how* to gate the migration so you find the breakage in a test harness instead of in production.

Confidence tags: **[HIGH]** = 3+ independent sources agree (incl. ≥1 provider primary source); **[MEDIUM]** = 2 sources or one strong source with caveats; **[LOW]** = single-source/contested. Sources and access dates in the References section. This is a fast-moving area, so provider specifics (model IDs, exact 400-error parameters, retirement dates) live in the `claude-api` skill and the official migration guides, which are always more current than this file.

## TRIGGER / SKIP

**TRIGGER** — load this reference when the task is about *the migration itself*, not one artifact:
- "we're moving from Opus 4.6 to 4.8 / from GPT-5.2 to GPT-5.5 — what breaks in our prompts/skills, and how do we re-optimize the whole library?"
- "a prompt that scored 94% last week scores 87% today and we didn't change it" (silent model refresh / prompt drift)
- "our JSON parser broke after the model upgrade"; "the model got more verbose / more literal / refuses more / stopped calling our tool after we upgraded"
- "set up a regression eval / golden transcripts / canary prompts before we flip the model"; "should we pin or float the model id?"; "deprecation deadline is coming — what's the migration playbook?"
- "make our prompt library model-agnostic / portable across Claude and GPT"

**SKIP** — route elsewhere:
- Re-optimize **one** prompt or skill in isolation → `prompt-deep-optimizer` / `skill-optimizer` (this skill decides *which* artifacts in a library to send there and *why*).
- Exact model IDs, pricing, which parameters now return `400`, retirement dates, the `/claude-api migrate` mechanical swap → the **`claude-api`** skill and the official Migration Guide. Cite them; never reproduce the table from memory.
- Build the eval *harness* / scorers / dataset infrastructure → `ai-agent-engineering ▸ Eval-Driven Development`. This skill consumes that harness for the before/after regression run.
- General prompt-engineering craft for a new model → the prompting hub. This file is only about the *delta* across a model change.

## The two failure modes (name them before you fix them) [HIGH]

Model-related prompt failures come in two distinct shapes, and the distinction drives the response (echostash; arXiv 2311.11123; models.news):

1. **Prompt regression (the loud one).** Explicit breakage when you migrate to a new model version. Parsers break, formats shift, behavior changes on the inputs your prompts were specifically tuned for. You know roughly when it happened (the flip), so you can gate it. The arXiv regression study found **58.8% of prompt+model combinations dropped accuracy across an API update, and 70.2% of those dropped by more than 5%** [HIGH]: regression is the norm, not the exception, on a non-trivial hop.
2. **Prompt drift (the quiet one).** Gradual, often unannounced behavioral change *within the same model id*, from quantization, serving-stack, or safety-tuning changes underneath you. The Stanford/UC-Berkeley "How Is ChatGPT's Behavior Changing over Time?" work documented large behavioral shifts between two dated snapshots served under the same GPT-4 model name (March vs June 2023) [HIGH] — the model changed underneath a stable id, which is exactly the drift case. You cannot point to a flip, so the only defense is **scheduled** regression runs on a cron, not change-triggered ones.

Survey context [LOW, unverified secondhand — do not cite as fact]: practitioner write-ups describe model-related production incidents as common while automated regression testing remains rare, and argue the low test coverage is part of why the felt failure rate is high. A specific "67% had an incident / 23% test" figure circulates in secondhand blog summaries but traces to no locatable primary survey (no Honeycomb publication carrying those numbers was found), so treat the *direction* (incidents common, coverage thin) as the only load-bearing claim here, not the percentages.

> **"Prompt rot."** The umbrella term practitioners use for both modes: a prompt that was good slowly or suddenly becomes worse without you editing it, because its quality was coupled to a model's behavior that has since changed.

## Why a prompt tuned for model A regresses on model B [HIGH]

The root cause is **coupling**: a tuned prompt encodes compensations for the *old* model's defaults, and those compensations become noise (or actively misfire) on the new model.

- **The dominant cross-generation shift is "more literal instruction-following."** This is corroborated independently across both major providers, which is unusual and makes it the single most reliable thing to expect on an upgrade:
  - OpenAI (GPT-4.1 Prompting Guide): *"GPT-4.1 is trained to follow instructions more closely and more literally than its predecessors… existing prompts optimized for other models may not immediately work with this model, because existing instructions are followed more closely and implicit rules are no longer being as strongly inferred."* [HIGH]
  - Anthropic (Migration Guide, Opus 4.6 → 4.7): *"Claude Opus 4.7 interprets prompts more literally and explicitly… it will not silently generalize an instruction from one item to another, and it will not infer requests you didn't make."* [HIGH]
  - Consequence: workarounds you wrote to *force* a behavior the old model under-did (extra examples to force generalization, redundant emphasis, all-caps) can cause the new model to over-attend to them. OpenAI: existing all-caps/"tips" tricks "could cause GPT-4.1 to pay attention to it too strictly." [MEDIUM]
- **The paradox of "better" models breaking prompts** [HIGH]. A model can improve on public benchmarks while regressing on *your* exact prompt patterns. echostash on the Feb-2026 GPT-4o→GPT-5 forced migration: *"GPT-5 was built to need less prompt engineering… In practice, prompts optimized for GPT-4o's quirks sometimes performed worse on GPT-5, not better."* The Tursio case (arXiv 2507.05573) is the canonical worked example: a regression suite passing fully on GPT-4-32k dropped a few points on GPT-4.1 and GPT-4.5-preview with the prompts unchanged, and a structured migration recovered it [MEDIUM — paper and direction confirmed; the exact per-model percentages were not re-verifiable at the source, so read them as small single-digit drops, not precise figures].

### The capability-shift regression checklist (what actually changes on upgrade) [HIGH]

These are the behavioral axes that move across a generation hop. Run your representative inputs against each (synthesized from Anthropic Migration Guide, the hidekazu worksheet, zro2.one, echostash):

| Axis | What shifts | Typical production failure |
| --- | --- | --- |
| **Output format** | Subtly different JSON structure / field ordering / whitespace; strict formatting instructions quietly deprioritized in favor of response quality | JSON / regex parsers break on day one. **Parse with a real parser, never substring-match.** |
| **Verbosity / response length** | Higher default reasoning effort tends to lengthen output; length-control prompts written for the old default over- or under-correct | Truncation against an unchanged `max_tokens`; cost/latency creep |
| **Tokenization** | A tokenizer change makes the *same text* count differently | `max_tokens` headroom wrong; budgets and truncation points off. Re-baseline with a token-count call. |
| **Refusal / safety posture** | Refusal rate and stop-reason distribution change; new stop reasons appear (`refusal`, `model_context_window_exceeded`) | Code that reads `content[0]` unconditionally crashes on a refusal |
| **Reasoning effort / "thinking" defaults** | Default thinking on/off and visibility change between generations; reasoning traces get shorter or are omitted by default | Chain-of-thought prompts produce shorter traces; rendered reasoning text empty unless you opt into a summary |
| **Tool-call style & frequency** | Tool-use rate shifts; newer models may reach for tools / sub-agents / memory *less* by default, preferring to reason; JSON escaping in tool args can differ | A tool you depend on stops being called; multi-agent flow spawns fewer sub-agents than the harness assumes |
| **Tone / persona** | Default register and "engaging-ness" shift | Brand-voice drift; over- or under-formality |

> For the *exact* parameter-level breaking changes (which params now return `400`, prefill removal, adaptive vs extended thinking, deprecation/retirement dates), go to the `claude-api` skill and the official Migration Guide — those are version-specific and change underneath this file.

## The migration playbook (gate the change in a harness, not in prod) [HIGH]

The sequence that vendor docs and practitioner postmortems converge on (Anthropic Migration Guide; zro2.one; echostash; EmberLM; models.news; hidekazu):

1. **Pin the current model to an exact snapshot** before you start, so you have a known-good baseline and an instant rollback. (See pinning section below.)
2. **Build / refresh a representative regression eval set first** — before touching any prompt. It *must predate the rewrite*: if you rewrite first, you can't tell whether drift came from the new model, the rewrite, or their interaction (AgentPatterns). Coverage = happy-path + edge cases + **the specific issues you previously fixed with prompt engineering** (turn every past incident into a permanent regression case) + output-format/parse compliance. Sizing: practitioners repeatedly report **50–100 carefully chosen cases catch more real regressions than thousands of synthetic ones** [HIGH]; the hidekazu worksheet uses ~40 as a concrete starting point.
3. **Score against the contract, not token-for-token.** LLMs are probabilistic; the question is "did the output still satisfy the contract?" not "did the bytes match?" Use property/assertion checks (schema validity, required fields, tool-call present, refusal-or-not) and LLM-as-judge for fuzzy quality — *not* exact-match. This is the `ai-agent-engineering ▸ Eval-Driven Development` machinery; reuse it.
4. **Run the full eval against the candidate model with the *existing* prompts** — change ONE thing at a time. Diff each output against the captured property, not against the old output token-for-token. This is the **golden-transcript** step: a stored set of (input → known-good behavior) pairs is your golden set; the diff localizes exactly which behaviors moved.
5. **Decide rewrite vs patch-forward per artifact** (next section) from what the diff actually shows — don't rewrite on faith.
6. **Re-optimize only the artifacts that regressed** (hand each to `prompt-deep-optimizer` / `skill-optimizer`), re-run the eval to convergence.
7. **Canary / gradual rollout.** Route a small slice of production traffic to the candidate behind a flag (e.g. 1% → 10% → 50% → 100%), watching the *same* metrics you measured offline (tool-call rate, refusal rate, token usage, latency) for a soak window (48h+ per stage is a commonly cited figure) at each stage. **Canary prompts** = a small set of high-signal probes you run continuously against prod to catch drift between full migrations.
8. **Keep the previous snapshot config for instant rollback**, update the prompt/skill registry + docs with the new pinned id, and set a deprecation alert for the version you just left.

### Shadow vs canary [MEDIUM]
- **Shadow testing** runs production traffic through *both* models simultaneously and diffs. It is the most comprehensive (catches regressions against real traffic, not just synthetic cases) but **doubles inference spend** during the window and needs dual-stream routing/compare infra.
- **Canary** sends only a small slice to the candidate. Cheaper, but only sees the traffic you route. Pick shadow when correctness is high-stakes and budget allows; canary otherwise.

## Re-optimization on upgrade: what to re-tune vs leave [HIGH]

**Classify the hop before deciding to rewrite anything.** Both providers now ship this as an explicit, *bounded* recommendation — the rewrite is conditional on observed drift, not assumed (AgentPatterns synthesizing OpenAI "Using GPT-5.5" + Anthropic Migration Guide). OpenAI's own upgrade tooling buckets upgrades three ways:

| Class | Choose when | Action |
| --- | --- | --- |
| **model-string-only** | Minor-version successor; prompts already short and task-bounded; no strict output-shape dependency. Anthropic: *"strong out-of-the-box performance on existing … prompts and evals."* | Swap the model id, keep prompts, **run the regression eval** and stop. Rewriting here discards a working prompt-eval coupling for no measured gain. |
| **model-string + light prompt rewrite** | Strict output shape; heavily scaffolded prompts; observed verbosity/literalism change; tool-heavy or multi-agent flow. Cross-generation hops default here. | Swap the string; rewrite *only* the prompts tied to the workflow risk; leave the rest. |
| **blocked** | Upgrade needs API-surface changes, parameter rewrites, or tool-handler rewiring (the `claude-api`-level stuff). | Report the blocker, do the API work first; do not improvise the prompt rewrite over a broken harness. |

### The four-axis re-tuning order (when you do rewrite) [HIGH]
When an artifact needs a rewrite, **start from the smallest prompt that preserves the product contract** (the externally observable behavior you owe the user: input domain, output shape, tool-call discipline, refusal posture, latency envelope): strip the inherited compensation layers *first*, re-tune second. Anthropic's checklist says exactly this: *"Re-baseline response length with existing length-control prompts removed, then tune explicitly"* and *"If you've added scaffolding to force interim status messages… try removing it."*

Then re-tune in this order (OpenAI "Using GPT-5.5"; the order is not arbitrary, because each stage destabilizes the ones below it):
1. **Reasoning effort.** Sets the depth of computation before anything downstream is observable. Tune it first or every later signal comes from the wrong substrate.
2. **Verbosity.** Higher effort lengthens output, so verbosity prompts written before effort is fixed over-correct. Anthropic's worked lever: *"to decrease verbosity, add: 'Provide concise, focused responses. Skip non-essential context, and keep examples minimal.'"*
3. **Tool descriptions.** Recalibrate against the new literal-interpretation profile and changed tool/sub-agent frequency. If the new model under-calls a tool, make the trigger condition explicit *in the system prompt and in the tool's own `description`* ("call this when the user asks about current prices or recent events") rather than relying on the old model's eagerness.
4. **Output format.** Fixed last, because format constraints surface most clearly once reasoning, length, and tool behavior are stable; re-tuning format against an unstable lower stage locks in transient artifacts.

> **LLM-as-critic re-write loop** (OpenAI cookbook Prompt Migration Guide) [MEDIUM]: a fast, repeatable way to produce the rewrite. Have the *new* model critique the old prompt for **ambiguity, undefined terms, conflicting instructions, and missing context/assumptions**, apply fixes, then re-run the eval. The tool RETAIN (EMNLP'24) automates the discovery half: it surfaces error "slices," and its refined prompts beat manual refinement (≈twice the errors found, 165 vs 86, and a 25% vs 12% quality gain in a fixed window) [MEDIUM: figures verbatim from the paper, but a single study]. This is the *generation* step; this skill owns the *gating* around it.

### When NOT to re-optimize [HIGH]
- **Minor-version successor with no eval drift** → model-string-only; a fresh rewrite is unmeasured churn.
- **Provider-managed harness** → for managed-agent products, "no changes beyond updating model name are required."
- **Audited / change-controlled prompts** → a prompt pinned by a regulatory regime, security review, or external sign-off can't be rewritten on every hop without re-running the audit; the certification cost dominates the tuning gain. Pin and defer.

## Cross-model portability patterns [MEDIUM]

To minimize re-work on every hop, reduce coupling up front:
- **Write the smallest contract-preserving prompt** and avoid model-specific *tells* — don't bake in the current model's quirks (its preferred JSON whitespace, its eagerness to call tools, all-caps "shouting" it happens to respond to). Those are the first things that break.
- **State implicit rules explicitly.** The cross-generation trend is toward literalism; a prompt that relies on a model inferring intent is the most fragile kind. Explicit beats implicit on every current model, so explicit phrasing is the portable default.
- **Decouple parsing from formatting.** Parse tool/JSON output with a real parser and validate against a schema; never substring-match. This survives format drift within *and* across models.
- **Capability probing before commit.** Don't assume a feature is present — probe it. Run a tiny suite that checks whether the candidate supports/uses the capabilities your prompt assumes (thinking on/off, tool-call willingness, structured-output mode, refusal posture) before you wire the full library to it.
- **Multi-model routing changes the calculus** [MEDIUM]: routing different tasks to different models (e.g. one model for analysis, another for code, a cheap one for high-volume) trades one big migration for more frequent, more *contained* ones — and amplifies whatever testing discipline you already have, for better or worse.

## Pinning vs floating model ids [HIGH]

- **Pin to an exact snapshot you have validated, and change it deliberately through the playbook**; never depend on a model string silently changing behavior underneath you. Note that some current Claude generations use a *dateless* id that is itself a pinned snapshot, not an evergreen auto-upgrading pointer, so verify in the `claude-api` skill which form you're using.
- **Floating / unversioned alias** = zero migration overhead and automatic improvements, but it's the **highest-risk** strategy for any system with parsers or structured-output needs: outputs can change with no action on your side, making fast regression detection nearly impossible. Reserve for low-stakes, non-parsed paths.
- **The pinning trap** [HIGH]: pinning without *ever testing forward* converts gradual drift into a cliff. Teams that pinned to GPT-4o in 2024 and never re-validated hit the Feb-2026 forced retirement far harder than teams that had been adapting incrementally. **Pin, plus scheduled forward-testing, plus a migration cadence** is the resilient combination; pinning alone only *defers* the risk.
- **Deprecation handling**: once a model is deprecated, requests past the retirement date fail outright, so migrate before the date. Watch the official deprecation table (in `claude-api` / the provider docs), set alerts the moment you adopt a snapshot, and treat the retirement date as a hard deadline that triggers the playbook above. Provider migration tooling (`/claude-api migrate`, OpenAI's `openai-docs` upgrade skill) automates the *mechanical* model-string swap and parameter fixes and then hands the *prompt review* to a human pass — the tooling does not do the re-optimization for you.

## Worked anchor: a prompt/skill library on Opus 4.6 → 4.8 [HIGH for sequence, defer specifics to `claude-api`]

1. **Read the migration guide layered, in order.** Breaking changes stack: 4.6 → 4.8 means applying the **4.7 changes and then the 4.8 changes** — *"read the official Migration Guide section for your exact target model, because breaking changes are layered."* (hidekazu). Get the exact 4.7 *and* 4.8 deltas from the `claude-api` skill.
2. **Pin 4.6** as the baseline and snapshot the registry; assemble/refresh the 50–100-case golden set across the library's representative skills, each scored by contract (schema valid, tool called, refusal-or-not, length bound), not exact match.
3. **Run the golden set on 4.8 with the existing (4.6-tuned) prompts unchanged.** Diff per behavior. Expect the literalism shift, possible tokenization/`max_tokens` headroom changes, changed thinking-default/visibility, and changed tool/sub-agent frequency to be the axes that move.
4. **Triage the diff into the three buckets.** Most short, task-bounded skills → model-string-only (just re-pin to 4.8 and keep the green eval). The few skills with strict output shapes, heavy length-control scaffolding, or tool/sub-agent assumptions → light rewrite. Anything needing parameter/harness changes (a now-`400` parameter, prefill, thinking config) → blocked: do that API work first via `claude-api`.
5. **For each rewrite candidate:** strip the 4.6-era compensation scaffolding, re-tune effort → verbosity → tool descriptions → output format against the golden cases, hand the artifact to `skill-optimizer` / `prompt-deep-optimizer`, re-run to convergence.
6. **Canary the upgraded library** (small traffic slice behind a flag), soak watching tool-call/refusal/token/latency, ramp, then re-pin the registry to the 4.8 snapshot, keep 4.6 config for rollback, and set a 4.6 deprecation alert.

## Key takeaways

- Name the failure mode first: **regression** (loud, on a flip; gate it) vs **drift** (quiet, within a pinned id; catch it with *scheduled* evals). Both are "prompt rot."
- A "better" model can regress *your* prompts; benchmark wins don't transfer to your contract. The reliable cross-generation shift is **more literal instruction-following**, which turns old compensation tricks into liabilities.
- **The eval set must predate the rewrite, change one thing at a time, and score the contract not the bytes.** 50–100 well-chosen cases beat thousands of synthetic ones.
- **Classify before rewriting:** model-string-only / light-rewrite / blocked. Re-optimize only what the diff shows regressed; for a rewrite, strip to the smallest contract-preserving prompt then re-tune **effort → verbosity → tool descriptions → output format**.
- **Pin, forward-test, set a cadence.** Pinning alone defers risk into a cliff; floating aliases are highest-risk for parsed outputs. Provider tooling does the mechanical swap; the prompt review is a human pass.
- This skill is the *gating methodology around* `prompt-deep-optimizer` / `skill-optimizer` (the re-writers) and `ai-agent-engineering ▸ Eval-Driven Development` (the harness); version-specific facts live in `claude-api`.

## References

Confidence: **[HIGH]** claims rest on ≥3 of the below incl. a provider primary source; **[MEDIUM]** on 2; **[LOW]** on a single study. All accessed **2026-06-15**.

**Provider primary sources (most current — always re-check; defer version specifics here):**
1. Anthropic — *Migration guide* (Claude API docs, platform.claude.com/docs/en/about-claude/models/migration-guide). Layered per-target-model breaking changes; the Opus 4.6→4.7 literalism/verbosity/subagent/tokenization checklist quoted throughout.
2. Anthropic — *Prompting best practices* and *Model deprecations* (platform.claude.com). Adaptive vs extended thinking; deprecated `temperature`/`top_p`/`top_k`→`400`; retirement table.
3. OpenAI — *GPT-4.1 Prompting Guide* (developers.openai.com/cookbook/examples/gpt4-1_prompting_guide). Source of the "more literal instruction-following … existing prompts may not immediately work" finding.
4. OpenAI — *Prompt Migration Guide* notebook (github.com/openai/openai-cookbook, examples/Prompt_migration_guide.ipynb). The LLM-as-critic rewrite loop (ambiguity / undefined terms / conflicts / missing context).

**Studies & papers:**
5. *Prompt Migration: Stabilizing GenAI Applications with Evolving Large Language Models*, arXiv 2507.05573 (Tursio case; GPT-4-32k→4.1/4.5 regression, migration testbed, three complexity tiers).
6. *(Why) Is My Prompt Getting Worse? Rethinking Regression Testing for Evolving LLM APIs*, arXiv 2311.11123 (58.8% of prompt+model combos regress; best-performing prompt changes per version; track both prompt and model versions).
7. *RETAIN: Interactive Tool for Regression Testing Guided LLM Migration*, EMNLP 2024 demo (aclanthology.org/2024.emnlp-demo.31). Error-slice discovery; defines "LLM migration."
8. Stanford / UC-Berkeley — *How Is ChatGPT's Behavior Changing over Time?* (arXiv 2307.09009, 2023). Documented drift between two dated snapshots served under the same model name (also surfaced via echostash, #13).

**Practitioner guides & incident data:**
9. AgentPatterns.ai — *Prompt-Rewrite Discipline on Cross-Generation Model Migration*. The three-bucket classification and the four-axis tuning order, synthesizing OpenAI "Using GPT-5.5" + Anthropic.
10. hidekazu-konishi.com — *Anthropic Claude Model Migration Guide* worksheet (2026-06-14). Per-dimension verification table; layered reading order; pinning guidance; ~40-case eval + canary steps.
11. zro2.one — *LLM API Versioning and Migration* guide (2026-03-20). Pin / eval / one-thing-at-a-time / gradual-rollout checklist; "scored the same on benchmarks" anti-pattern.
12. EmberLM — *Prompt Regression Testing: The Complete Guide for 2026*. Golden dataset, cron-scheduled regression for silent updates, CI gating on pass-rate threshold.
13. echostash — *Model updates and prompt stability: what the data shows* (2026-03-06). Regression-vs-drift distinction; four-strategy taxonomy (pin / float / shadow / gradual); 50–100-case finding. (Carries the secondhand "67%/23%" incident/coverage figures the body now flags as untraceable to a primary survey — vendor blog, treat as direction-only.)
14. models.news — *Prompt Versioning and Regression Testing Guide* (2026-06-08). "Did the output satisfy the contract?" framing; triggers for re-running the suite (model update, settings, new tool, policy change).
