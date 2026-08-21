<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `eval-driven-development` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: eval-driven-development
title: Eval-Driven Development for LLM Apps
description: >
  The BUILD-TIME discipline of evaluating LLM/agent applications — the
  analyze→measure→improve loop with **error analysis as the engine**, distinct
  from runtime production monitoring (llm-observability) and academic benchmark
  leaderboards (da-7). Covers eval-driven development as a practice (Hamel Husain's
  "your AI product needs evals", 60-80% of effort on error analysis; the analyze→
  measure→improve loop; "evaluators for errors you discover, not imagine"); the
  Three Gulfs (Specification / Generalization / Comprehension); error analysis &
  qualitative coding (open/axial coding of traces, the SME "benevolent dictator");
  criteria drift & "Who Validates the Validators?" (Shankar UIST 2024, EvalGen /
  SPADE assertion synthesis); eval levels (code/assertion vs LLM-as-judge vs human;
  offline dev-time vs online feedback flywheel); LLM-as-judge depth (position /
  verbosity / self-preference bias, calibration vs human labels via Cohen's kappa /
  Krippendorff's alpha / TPR-TNR, binary-beats-Likert, pairwise vs pointwise);
  golden datasets & synthetic eval-data (silver→gold, 50-100 curated inputs);
  metric design (assertions, rubric scoring, pairwise/Elo, pass@k vs pass^k,
  faithfulness/groundedness); CI/CD eval gating & regression suites; agent/
  trajectory evaluation (tool-call correctness, earliest-critical-decision); and
  the tooling landscape (promptfoo, DeepEval, Ragas, LangSmith, OpenAI Evals,
  Inspect/AISI, Arize Phoenix, Langfuse, Braintrust). TRIGGER: "evals for my LLM
  app", eval-driven development, error analysis, LLM-as-judge design/calibration/
  bias, golden dataset / synthetic eval data, rubric vs binary scoring, pass@k vs
  pass^k, eval-gated CI / regression suite, agent/trajectory eval, choosing an
  eval tool (promptfoo/DeepEval/Ragas/Braintrust/LangSmith). SKIP: production
  runtime tracing / live monitoring / OTel dashboards (use llm-observability — EDD
  only owns the feedback flywheel into the dataset); academic leaderboard harnesses
  HELM/MMLU/LLM-as-judge-benchmarking (use da-analytical-methods ▸
  da-7-machine-learning); prompt-optimization algorithms that consume eval signals
  (use prompt-deep-optimizer / prompt-helper-optimizer); RAG-specific retrieval
  eval depth (use advanced-rag-patterns).
origin: local
category: developer
version: "1.0"
updated: "2026-05-31"
tags:
  - evals
  - eval-driven-development
  - llm-as-judge
  - error-analysis
  - testing
  - llm
  - agent
whenToUse:
  - "building an eval system / eval-driven development for an LLM or agent app"
  - "error analysis — turning real failures into evaluators (open/axial coding)"
  - "designing or calibrating an LLM-as-judge (bias, kappa, binary-vs-Likert)"
  - "curating a golden dataset or generating synthetic eval data (silver→gold)"
  - "metric design (assertions, rubric, pairwise, pass@k vs pass^k)"
  - "eval-gated CI/CD and regression suites for prompts/agents"
  - "agent/trajectory evaluation (tool-call correctness, multi-step success)"
  - "choosing an eval tool (promptfoo, DeepEval, Ragas, LangSmith, Braintrust)"
whenNotToUse:
  - "production runtime tracing / live monitoring / OTel dashboards — use llm-observability"
  - "academic benchmark harnesses (HELM/MMLU) — use da-analytical-methods ▸ da-7-machine-learning"
  - "prompt-optimization algorithms that consume eval signals — use prompt-deep-optimizer"
  - "RAG-specific retrieval-quality eval depth — use advanced-rag-patterns"
related_skills:
  - ai-llm-model-layer
  - ai-rag-retrieval
  - prompt-deep-optimizer
  - ai-agents-orchestration
---

# Eval-Driven Development for LLM Apps

The build-time evaluation discipline. The founding thesis (Hamel Husain): the
biggest divider between reliable AI systems and "YOLO" development is a robust
eval system — in practice **60-80% of build effort goes to error analysis and
evaluation**, not prompt-writing. Evals are part of the inner loop the way
debugging is part of software development. The loop is **Analyze (read traces /
error analysis) → Measure (write evaluators for *real* failures) → Improve (fix
prompt/retrieval/architecture) → repeat**, and the rule that anchors it: **write
evaluators for errors you *discover*, not errors you *imagine*.**

> **Boundary:** runtime tracing / live production monitoring is `llm-observability`
> (EDD's only online interest is the *feedback flywheel* that turns sampled prod
> failures into new offline test cases). Generic leaderboard harnesses (HELM/MMLU)
> are `da-7-machine-learning` — the "generic off-the-shelf metrics" anti-pattern is
> precisely the boundary marker.

## Error analysis — the engine

Borrowed from grounded-theory qualitative research:
- **Open coding** — read 30-50 traces, write freeform notes on each failure (no
  fixed taxonomy yet). Do it *yourself* first (the domain-expert "benevolent
  dictator") — don't outsource the first pass.
- **Axial coding** — cluster the notes into failure-mode categories (an LLM can
  help cluster; a human refines).
- Write evaluators only for failure modes that (a) recur and (b) aren't fixable by
  a trivial prompt change. Error analysis "sorts failures into the right bucket so
  you don't build evaluators for problems a prompt change would have solved."
- **Rule of thumb:** ~30 min reading 20-50 outputs after any significant change.

**The Three Gulfs (Shankar/Husain):** *Specification* (gap between what you want
and what you told the model), *Generalization* (works on seen examples vs the long
tail), *Comprehension* (what the system does vs what you understand it to do).
Error analysis primarily attacks Comprehension; better specs attack Specification;
robust data + retrieval attack Generalization.

**Criteria drift** ("Who Validates the Validators?", Shankar UIST 2024): you need
criteria to grade outputs, but grading outputs is what helps you define the
criteria — so evaluation criteria are often **not definable a priori**; they
emerge from reading real outputs. This is the theoretical basis for
error-analysis-first. Companion systems: **EvalGen** (mixed-initiative — generate
candidate assertions + judge prompts, have the human grade a subset, select the
implementations that best align) and **SPADE** (synthesize data-quality assertions
from prompt-version history).

## Eval levels & the offline/online split

| Level | What | When |
|---|---|---|
| **Code/assertion** | Deterministic: exact-match, regex, JSON-schema, contains, latency, cost | First line; CI gate |
| **LLM-as-judge** | A model scores against a rubric or picks a pairwise winner | Subjective/open-ended quality — *must be validated* |
| **Human** | Domain-expert labels | The source of truth judges are calibrated against |

**Offline (dev-time)** runs against a golden set in CI / the inner loop;
**online** samples live traffic and feeds failures back into the dataset. The two
form a flywheel — online failures become offline test cases (the EDDOps closed
loop).

## LLM-as-judge depth

**Biases:** *position* (favors a slot, often first — can shift >10% on code tasks;
mitigate by swapping positions and averaging), *verbosity/length* (longer scores
higher), *self-preference* (a model over-rewards its own output — use a *different*
model family as judge).

**Calibration / validation against humans:** build a human-labeled gold set;
measure judge-vs-human agreement; iterate the judge prompt. Use **Cohen's kappa**
(two raters) / **Krippendorff's alpha** (many) for *agreement beyond chance* — a
judge can correlate perfectly yet be systematically too harsh/lenient, which
correlation misses and kappa catches. Report **TPR/TNR** per class. For ordinal
(Likert) scores use Kendall's tau / Spearman. **Binary beats Likert** as the
default — easier to calibrate, easier inter-annotator agreement, standard
classification metrics apply (well-validated rubrics still have a place for
nuanced domain quality — "always binary" is an over-simplification). **Pairwise**
(A-vs-B) is often more reliable than **pointwise** for subjective quality, at
O(n²) cost + its own position bias.

## Metric design

Assertions/exact-match (structured output) · rubric/criteria scoring · pairwise
preference (Bradley-Terry/Elo) · **pass@k** (≥1 of k passes — capability/coverage)
vs **pass^k** (*all* k pass — reliability/consistency, the one that matters for
agents that must not fail intermittently) · RAG metrics (faithfulness/groundedness,
answer relevance, context precision/recall — depth in `advanced-rag-patterns`).

## Golden datasets & synthetic eval data

10-20 examples is enough to *track* iterative improvement; a proper golden set is
~50-100 curated, expert-labeled inputs (scale with complexity/risk). Balance
typical + edge cases; include tricky regressions. **Synthetic — silver→gold:**
generate "silver" synthetic inputs, promote to "gold" via SME review +
evaluator-agreement + bias audit; DeepEval Synthesizer / Phoenix *evolve* inputs to
widen coverage. Always validate synthetic data against human judgment.

## CI/CD & agent evaluation

- **Eval-gated deploys:** a prompt/model/code change must not regress the golden
  set beyond a threshold before merge; every fixed bug becomes a permanent
  regression test. The **dev inner loop**: edit → offline eval → read failures →
  repeat, fast and local, before any CI gate.
- **Agent/trajectory eval:** beyond final-answer correctness — tool-call
  correctness (right tool, right args), step ordering, plan quality, multi-step
  success; find the *earliest critical decision* that triggered a failure cascade.
  pass^k is especially relevant. (Coding-agent eval depth → coding-agents; agent
  harness → `ai-agents-orchestration` (references/agent-harness-construction.md).)

## Tooling landscape

Teams converge on **two tools**: a CI-gating framework + an annotation/regression
platform.
- **promptfoo** — YAML-driven eval + red-team (note: acquired by OpenAI 2026-03);
  prompt/model matrix + CI. **DeepEval** — "pytest for LLMs", 50+ metrics, agent
  eval. **Ragas** — RAG-specific metrics. **OpenAI Evals** — open registry.
  **Inspect** (UK AISI) — safety/capability evals.
- Platforms: **Braintrust**, **LangSmith** (LangChain-coupled), **Arize Phoenix**
  (open-core, eval + monitoring), **Langfuse** (OSS observability + eval — straddles
  `llm-observability`). *Treat any single-vendor "alternatives" ranking as
  directional, not neutral.*

## Anti-patterns

- **Vibes-based eval** — eyeballing with no dataset/metric.
- **No error analysis** — evaluators for imagined, not observed, failures.
- **Overfitting the eval set** — train/test discipline applies.
- **Unvalidated judge** — trusting an LLM judge never calibrated against humans.
- **Generic off-the-shelf metrics** — vendor defaults instead of metrics derived
  from *your* error analysis (this is the boundary vs academic benchmarks).
- **EDD taken too literally** (contested) — strict "write all evals before any
  prompt" fights criteria drift; the reconciling synthesis: stand up the *harness*
  early, derive *specific evaluators* from error analysis on real outputs.

## Cross-references

- **Runtime tracing / live production monitoring / online eval dashboards** →
  `ai-llm-model-layer` (references/llm-observability.md) (EDD owns only the failure→dataset flywheel).
- **Academic leaderboard harnesses (HELM/MMLU)** → `da-analytical-methods ▸
  da-7-machine-learning`.
- **Prompt-optimization algorithms that consume eval signals** →
  `prompt-deep-optimizer`, `prompt-helper-optimizer`.
- **RAG retrieval-quality eval (Ragas depth)** → `ai-rag-retrieval` (references/advanced-rag-patterns.md),
  `rag-architecture`.
- **Agent harness / coding-agent eval** → `ai-agents-orchestration` (references/agent-harness-construction.md),
  `coding-agents`.

## References

Hamel Husain "Your AI Product Needs Evals" + the evals-FAQ (Husain/Shankar);
Shankar et al. "Who Validates the Validators?" (UIST 2024, arXiv:2404.12272) +
SPADE (2401.03038); Pragmatic Engineer "pragmatic guide to LLM evals"; the EDDOps
process-model paper (2411.13768); Eugene Yan "Evaluating LLM-Evaluators"; Evidently
LLM-as-judge guide; OpenAI eval-driven-system-design cookbook + Evals; Langfuse
error-analysis posts; DeepEval / Ragas / Arize / Braintrust docs. *(18 sources,
2024-2026; full URLs in the source research report.)*
