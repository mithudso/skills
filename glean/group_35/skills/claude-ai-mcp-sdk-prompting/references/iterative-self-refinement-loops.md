<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.**
> Sibling topics in this family are reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: iterative-self-refinement-loops
title: "Iterative Self-Refinement Loops"
version: "1.1.0"
updated: "2026-06-15"
description: >
  Methodology for designing iterative refine -> evaluate -> refine loops for LLM artifacts
  with principled stop conditions — the loop engineering UNDER concrete optimizers
  (ddo/pdo/sko/cdo, document-critique). Covers Self-Refine (Madaan 2023), Reflexion (verbal
  RL / self-reflection memory across attempts), CRITIC (tool-verified self-correction),
  self-consistency as a parallel loop, the Pan 2023 self-correction taxonomy, the critical
  honesty point that intrinsic self-correction of reasoning FAILS without external feedback
  (Huang 2023 oracle-feedback caveat), stop-condition design (convergence/diminishing-returns
  detection, oscillation & thrash detection, budget caps, severity-gated exit), external
  feedback that makes loops converge (tests, verifiers, judges, linters), and iterative
  refinement vs test-time compute (best-of-N, verifier-guided search).
  TRIGGER: designing or debugging an iterative refine/critique/revise loop for an LLM
  artifact; choosing a stop condition; "how many iterations/passes"; loop oscillates / output
  keeps growing / polishes errors; whether self-critique helps or you need an external
  verifier; sequential refinement vs best-of-N; building the convergence engine shared by
  ddo/pdo/sko/cdo/document-critique.
  SKIP: the concrete optimizers that USE this loop (pdo, sko, ddo/document-critique, cdo);
  test-time compute as a decoding/serving implementation — kernels, sampling, PRM training
  (reasoning-models, ai-llm-model-layer), though the refine-vs-best-of-N DESIGN choice stays
  here (Section 6); weight-updating/gradient RL incl. episodic-reward fine-tuning (agentic-rl),
  though Reflexion's in-context "verbal RL" (no weight update) is in scope; prompt wording
  technique (-> prompt-engineering).
whenToUse:
  - "Designing an iterative refine -> evaluate -> refine loop for an LLM artifact"
  - "Choosing a stop condition: convergence detection, oscillation/thrash detection, budget cap, severity gate"
  - "Deciding how many refinement iterations to run (k=1, k=2, cap at 3?)"
  - "Diagnosing a loop that oscillates, expands, or polishes errors into confident-sounding errors"
  - "Deciding whether self-critique actually helps or whether you need an external verifier"
  - "Wiring external-feedback signals (tests, linters, judges, tool verifiers) into a loop"
  - "Choosing sequential iterative refinement vs best-of-N / verifier-guided sampling"
  - "Building or reviewing the convergence engine shared by ddo/pdo/sko/cdo/document-critique"
metadata:
  changelog:
    - "2026-06-15 sko v1.0.0->v1.1.0 — Pass H N/A (reference file, hub owns routing); 1 High + 7 Medium fixed: added low_quality_plateau to the §5.5 stop predicate + reconciled the 6->7 count and 2+4+1 partition, operationalized the plateau detection rule, history-drop→oscillation mechanism + §5.2 cross-link, front-loaded the Self-Refine all-self-checkable caveat, SKIP rewrite (test-time-compute design-choice carve-out + agentic-rl/Reflexion verbal-RL seam), em-dash density 1.86->0.78 per 100 words, +RefineBench arXiv:2511.22173. All 10 sources verified real; GainNet exclusion confirmed clean"
---

# Iterative Self-Refinement Loops

**Iterative self-refinement** is the discipline of designing `generate -> evaluate -> refine` loops that improve an LLM artifact (prose, code, a prompt, a plan, an answer) across multiple passes, governed by a **principled stop condition** rather than a fixed guess. It is the loop-engineering layer that sits *under* every concrete optimizer in this hub's ecosystem — `prompt-deep-optimizer` (pdo), `skill-optimizer` (sko), `document-critique` / `ddo`, and `code-deep-optimizer` (cdo) all instantiate the same pattern. This reference is the shared methodology; those four are the productized applications.

The single most important result in this literature is also the least comfortable one: **a model refining its own reasoning with no external feedback usually does not improve, and frequently gets worse** (Huang et al., 2023). Loop engineering is therefore mostly about (a) sourcing a *real* feedback signal and (b) knowing when to stop. The rest is plumbing.

---

## 1. The canonical loop

```
y0  = generate(x)                      # initial draft
for t in 0, 1, 2, ...:
    fb_t = evaluate(x, y_t)            # feedback / critique / verifier signal
    if stop(fb_t, t): break            # principled stop condition (Section 5)
    y_{t+1} = refine(x, history, fb_t) # revise, conditioned on prior outputs + feedback
return y_t
```

This is **Self-Refine's Algorithm 1** (Madaan et al., 2023) almost verbatim. Three knobs decide whether the loop works:

1. **Where `evaluate` gets its signal:** the model itself (intrinsic), the model plus a tool/verifier (tool-verified), or a separate critic/judge/human (external). This choice dominates everything else (Sections 3–4).
2. **What `refine` sees:** Self-Refine retains the *history* of prior outputs and feedback in the prompt so the model "learns from past mistakes and avoids repeating them." Dropping history invites oscillation: with no record of prior drafts, the refiner can re-introduce an error it already removed or revert the critic's last fix each pass (see 5.2).
3. **When `stop` fires:** the difference between a loop that converges and one that burns budget polishing or thrashing (Section 5).

---

## 2. The four core loop archetypes

| Archetype | Feedback source | Loop shape | Canonical paper | Converges reliably? |
|---|---|---|---|---|
| **Self-Refine** | Same model, self-feedback | Sequential: feedback -> refine, history-carried | Madaan et al. 2023 | Only on tasks with self-checkable structure |
| **Reflexion** | Task signal -> verbal self-reflection stored in memory | Sequential across *attempts* (episodes), memory-carried | Shinn et al. 2023 | Yes *when* a task/environment signal exists |
| **CRITIC** | External **tools** (interpreter, search, calculator) | Sequential: verify-with-tool -> critique -> revise | Gou et al. 2023/24 | Yes — external grounding is the point |
| **Self-Consistency** | None per-sample; aggregate by vote | **Parallel** sample N, majority/weighted vote | Wang et al. 2022/23 | Yes for tasks with a discrete answer |

### 2.1 Self-Refine (Madaan et al., 2023 — `arXiv:2303.17651`, NeurIPS 2023)
One model plays all three roles (generator, feedback provider, refiner) with three prompts and **no training**. The key design requirements the paper stresses: feedback must be **actionable and specific** ("the sentiment is neutral due to phrases like 'good'", not "make it better"), and the refine step must see the **full history**. Reported ~20% absolute average improvement across 7 tasks on GPT-3.5/GPT-4 — but every one of those 7 tasks (dialogue, code optimization, constrained generation) was deliberately self-checkable, with outputs the model *can* verify. The headline number is the optimistic end of the literature and does not generalize to free reasoning (see Section 3).

### 2.2 Reflexion (Shinn et al., 2023 — `arXiv:2303.11366`, NeurIPS 2023)
"Verbal reinforcement learning": instead of updating weights, the agent **verbally reflects** on a task feedback signal and writes the reflection into an **episodic memory buffer** that conditions the *next attempt*. The loop is across attempts/episodes, not within a single draft. Crucially, Reflexion is "flexible enough to incorporate various sources (external or internally simulated) of feedback" — its headline 91% pass@1 on HumanEval comes precisely because coding gives a **real external signal** (unit tests / interpreter). Reflexion is the bridge between pure self-critique and tool-verified correction: it is a memory architecture wrapped around whatever feedback you can get.

### 2.3 CRITIC (Gou et al., 2023/24 — `arXiv:2305.11738`, ICLR 2024)
The thesis is in the title and in the repo's own words: *"we find that LLMs' Self-Verification and Self-Correction are unreliable; and we propose CRITIC, which enables LLMs to validate and rectify themselves through interaction with external tools."* The loop verifies the current output with an **appropriate tool** (search engine for facts, code interpreter for code, calculator for math), produces a tool-grounded critique, then revises. CRITIC is the constructive counterpart to Huang 2023: self-correction works *when you give it teeth.*

### 2.4 Self-Consistency (Wang et al., 2022/23 — `arXiv:2203.11171`, ICLR 2023)
Not a refinement loop at all, but a **parallel** sampling strategy. Sample N diverse chain-of-thought paths, then **marginalize over reasoning paths by majority (or verifier-weighted) vote** over the final answers. It belongs here as the contrast case: it improves answers *without* any iterate-on-the-same-draft step, and it is the conceptual seed of best-of-N test-time compute (Section 6). Big gains on discrete-answer reasoning (GSM8K +17.9%, etc.); inapplicable where there is no countable answer to vote on.

### 2.5 The map: Pan et al. (2023/24) self-correction taxonomy
`arXiv:2308.03188` (TACL 2024) is the field's survey. Its organizing axis is **when** correction happens (**training-time**, **generation-time**, **post-hoc**) and **who** produces the feedback (the model itself vs an external system). Use it to place any new technique: most of what this reference covers is *generation-time/post-hoc* correction with *automated* feedback, the slice "most practical and deployable with minimal human feedback."

---

## 3. The critical honesty point: intrinsic self-correction of reasoning fails

**Huang et al., 2023 — "Large Language Models Cannot Self-Correct Reasoning Yet"** (`arXiv:2310.01798`, ICLR 2024, Google DeepMind). This is the load-bearing caveat for the whole discipline; do not design a loop without it.

- It defines **intrinsic self-correction**: the model fixes its initial answer "based solely on its inherent capabilities, **without the crutch of external feedback**."
- Finding: on reasoning, intrinsic self-correction **does not help and often *degrades* performance**; "the performance after self-correction even deteriorates."
- The demolition of prior optimism: the apparent gains in earlier self-correction work came from **oracle labels** (the loop was told *whether* the answer was wrong, or *when* to stop). *"The improvements in these studies result from using oracle labels to guide the self-correction process, and the improvements vanish when oracle labels are not available."*

**The oracle-feedback caveat, stated plainly:** if your loop's stop signal or its critique secretly depends on knowing the right answer, your benchmark is measuring the oracle, not the model. In production you usually do not have the oracle. So:

- A self-critique loop with **no external signal** is, for genuine reasoning, **hope, not signal**. It tends to **polish errors into confident-sounding errors**.
- Corroborating evidence: *RefineBench* (Lee et al., 2025, `arXiv:2511.22173`) found self-refinement without external feedback yielded **+1.8 pp or less over five iterations**, while **guided** refinement (real external checker) approached near-perfect; and models "routinely halt early due to overconfidence even when errors remain." The "Dark Side of Intrinsic Self-Correction" line of work (2024) shows models can spiral into excessive "think" loops and talk themselves out of correct answers (recency/refinement bias).

**Design rule:** the value of a refinement loop is roughly the value of its feedback signal. If you cannot point to where real signal enters, expect zero-to-negative lift and a tendency to over-edit.

---

## 4. External-feedback signals that make loops actually converge

The fix for Section 3 is to give `evaluate` a grounded signal. In rough order of reliability:

1. **Executable / verifiable:** unit tests, type checker, compiler, linter, schema validator, a calculator, a runnable example. Binary or near-binary, hard to game. This is why coding loops (Reflexion's HumanEval, CRITIC's program synthesis) converge.
2. **Tool-verified:** search engine for fact-checking, a knowledge base, a retrieval call, a code interpreter (CRITIC's whole premise). Grounds claims against an external world.
3. **Separate-model judge / verifier:** an LLM-as-judge or a trained verifier/reward model scoring the output. Weaker (can share the actor's blind spots) but applicable to prose/design where no test harness exists. Best practice: **separate the critic from the actor** (different prompt, different call; do not let the critic see the actor's chain of reasoning) to reduce shared bias.
4. **Self-critique with named criteria:** the weakest. Only trust it when the artifact is **self-checkable** (does this JSON match the schema? does this code run?) and when the critique prompt **names specific criteria and demands quoted evidence + an explicit verdict.** Never for free reasoning.

A practical hierarchy for stop signals: **prefer a real validator's PASS/FAIL over the model's "looks good to me."** The model's self-reported "no further improvements possible" is poorly calibrated (Kadavath et al., 2022, on calibration).

---

## 5. Stop-condition design (the core engineering)

A loop without a disciplined stop condition either ships too early or thrashes. Combine **multiple** of these; never rely on a single self-reported flag.

### 5.1 Convergence / diminishing-returns detection
Track an observable signal across passes and stop when it flattens:
- **Change velocity:** size of the edit/diff from `y_t` to `y_{t+1}`. Once edits go cosmetic, the easy errors are gone.
- **Content similarity:** high similarity between consecutive outputs => stabilized.
- **No-progress / no-improvement detector:** if `y_{t+1}` is byte-identical (or trivially different) to `y_t`, stop. In practice this fires earlier than expected — once the easy edits are exhausted, the model often returns a near-identical draft.
- **Findings velocity** (the optimizer pattern, as in document-critique/ddo): when consecutive critique passes surface near-identical findings, the output has stabilized.

> **Convergence != correctness.** These signals measure whether the output has *stopped changing*, not whether it is *right*. A loop can converge to a confidently wrong, self-consistent answer. For high-stakes outputs, pair convergence detection with a Section-4 external checker.

### 5.2 Oscillation & thrash detection
- **Oscillation / loop_detected:** output alternates between two near-identical versions (the actor reverts to a prior draft each time the critic pushes the other way). Without a loop detector you "burn the full budget alternating between two answers." Detect by hashing recent outputs and watching for repeats.
- **Expansion / scope drift:** output **grows every pass** instead of stabilizing. Symptom of a runaway "add more" critique. Cap output size or gate on net-new value, not just "any change."
- **Low-quality plateau** — the convergence signals agree (the diff has gone cosmetic) but quality is still below bar. *Detect by:* the diff/similarity check fires **while** the latest critique still reports open findings at or above the gate, or an absolute judge/validator score sits below a preset floor. Distinguish it from 5.3: a severity-gate exit means the findings cleared; a plateau means they persist while edits have stopped helping. The approach needs a **restart/redesign, not more passes** (different prompt, different feedback source, or escalate to a human).

### 5.3 Severity-gated exit (the optimizer convention)
The pattern used by document-critique/ddo/cdo/pdo/sko: **loop until no Medium-or-higher findings remain.** Each pass produces severity-ranked findings; the loop exits only when the highest open severity drops below the gate. This makes "done" a property of the artifact, not a fixed iteration count — but it **must** be paired with a hard cap (5.4) so a stubborn Medium can't loop forever.

### 5.4 Budget caps (the non-negotiable backstop)
Always set a **hard max-round limit** as a cost fallback, independent of every other signal. Practitioner consensus is strong and convergent:
- **Default to k=1**; add a second iteration only if your eval *shows* the second pass routinely helps.
- **Treat k=3 as a maximum** for unguided self-critique. "Past k=3, revisions usually drift without improving"; iteration three "usually does not" beat iteration two on eval.
- The shape isn't a model defect; it's **information theory**: critique without external grounding is bounded by what actor and critic already jointly know; once the easy errors are squeezed out in one or two passes, the next pass has nothing to find.
- When you *do* have a strong external verifier (tests), you can afford a higher cap because each pass has real signal to act on.

### 5.5 A composable stop predicate
```
stop(state) := any of:
    severity_gate_met(state)         # no Medium+ findings  (5.3) — primary "success" exit
 or validator_pass(state)           # tests/lint/schema PASS (4)  — strongest success signal
 or no_progress(y_t, y_{t-1})       # cosmetic/identical edit (5.1)
 or oscillation_detected(history)   # alternating outputs    (5.2)
 or expansion_detected(history)     # output growing         (5.2)
 or low_quality_plateau(state)      # converged but below quality bar (5.2) — escalate/redesign
 or t >= max_iters                  # hard budget cap        (5.4) — always present
```
The first two are *success* exits; the middle four are *give-up / redesign* exits (`no_progress` plus the three thrash signals); the last is the *cost* backstop. A production loop wires all seven.

---

## 6. Iterative refinement vs test-time compute (best-of-N, verifier-guided)

These are two different ways to "spend more inference compute for a better answer," and they trade off.

- **Sequential / iterative refinement** = revise the *same* draft over multiple passes (Self-Refine, Reflexion). Modifies the *proposal distribution* by asking the model to improve its own prior output.
- **Parallel test-time compute** = sample *many independent* candidates and **select** with a verifier or vote (self-consistency = vote; best-of-N = verifier/reward-model picks the top). Modifies *how the verifier is used*.

**Snell et al., 2024 — "Scaling LLM Test-Time Compute Optimally"** (`arXiv:2408.03314`) unifies these as two axes (*refine the proposal distribution* vs *search against a verifier*) and shows the **right choice depends on problem difficulty**:
- **Best-of-N / sampling** is more effective on **easier** problems and at **higher** compute budgets.
- **Sequential revision** and **beam/tree search against a process reward model (PRM)** are more effective on **harder** problems and at **lower** budgets.
- A **process reward model** (scores each intermediate *step*, not just the final answer) enables verifier-guided **tree search**, beating naive best-of-N by up to ~4x compute efficiency; on FLOPs-matched comparisons, optimal test-time compute let a small model beat a 14x larger one on problems where it had non-trivial success.

**Engineering takeaway:** they compose. A strong loop often = **best-of-N to widen the proposal** (parallel) + **a verifier to select** + **a short sequential refine pass on the selected candidate**. Pick the mix by difficulty and budget: easy/cheap => lean on sampling+vote; hard/expensive => lean on verifier-guided sequential refinement. (Parallel test-time compute *as a decoding/model-layer concern* belongs to `reasoning-models` / `ai-llm-model-layer`; here it's the design alternative to a refine loop.)

---

## 7. Decision guidance — should I even run a loop, and what kind?

```
Do you have a REAL external feedback signal (tests / verifier / tool / judge / human)?
  NO  -> Is the artifact self-checkable in structure (schema, runnable code, format)?
          YES -> Run a SHORT self-critique loop (k<=2), criteria-named, no-progress exit.
          NO  -> For genuine REASONING: do NOT loop on self-critique (Huang 2023).
                 Prefer self-consistency (vote over N samples) or a single best draft.
                 Get a signal before looping.
  YES -> Is the answer discrete / countable?
          YES & easy/cheap   -> Best-of-N + verifier-weighted vote (parallel).
          YES & hard         -> Verifier-guided (PRM) search + short sequential refine.
          Free-form artifact -> Tool/judge-in-the-loop refine (CRITIC-style),
                                severity-gated, separate actor/critic, hard cap.
Across attempts with a task signal & you want carryover learning -> add Reflexion memory.
Always: set a hard max-iteration cap; detect oscillation & expansion; convergence != correct.
```

**One-line heuristics**
- Default to **one** refinement pass; earn the second with eval evidence; cap at three for unguided loops.
- The loop is only as good as its feedback signal — if you can't name where signal enters, expect over-editing.
- **Separate the critic from the actor**; make critiques name criteria + quote evidence + emit an explicit verdict.
- **Convergence is not correctness** — a self-consistent loop can stabilize on a confident wrong answer.
- Prefer a **validator's PASS/FAIL** over the model's "looks good"; the model's self-stop is poorly calibrated.
- **Oscillation and expansion mean restart/redesign**, not more passes.

---

## 8. Relationship to the concrete optimizers (read this, don't reinvent them)

This reference is the **shared loop-engineering methodology**. The hub ecosystem already ships productized loops that implement it — use them rather than rebuilding:

- **`prompt-deep-optimizer` (pdo)** — multi-pass audit + convergence loop for production prompts.
- **`skill-optimizer` (sko)** — convergence-loop quality gate for skill files.
- **`document-critique` / `ddo`** — multipass prose critique that applies Medium+ fixes and loops to convergence (the severity-gated pattern in 5.3 is theirs).
- **`code-deep-optimizer` (cdo)** — review-and-fix loop verified by build/lint/tests (its tests-as-feedback is the Section-4 ideal).

When you need the *loop* (stop conditions, oscillation handling, feedback-signal choice), read this. When you need to *run* one against a real artifact, invoke the matching optimizer.

---

## Sources

Primary papers (preferred):
- **Madaan et al. (2023)** — *Self-Refine: Iterative Refinement with Self-Feedback.* `arXiv:2303.17651`, NeurIPS 2023. https://arxiv.org/abs/2303.17651 · project: https://selfrefine.info/
- **Shinn et al. (2023)** — *Reflexion: Language Agents with Verbal Reinforcement Learning.* `arXiv:2303.11366`, NeurIPS 2023. https://arxiv.org/abs/2303.11366
- **Gou et al. (2023/24)** — *CRITIC: Large Language Models Can Self-Correct with Tool-Interactive Critiquing.* `arXiv:2305.11738`, ICLR 2024. https://openreview.net/forum?id=Sx038qxjek
- **Huang et al. (2023)** — *Large Language Models Cannot Self-Correct Reasoning Yet.* `arXiv:2310.01798`, ICLR 2024 (Google DeepMind). https://arxiv.org/abs/2310.01798 — **the oracle-feedback caveat.**
- **Wang et al. (2022/23)** — *Self-Consistency Improves Chain of Thought Reasoning in Language Models.* `arXiv:2203.11171`, ICLR 2023. https://arxiv.org/abs/2203.11171
- **Snell et al. (2024)** — *Scaling LLM Test-Time Compute Optimally can be More Effective than Scaling Model Parameters.* `arXiv:2408.03314`. https://arxiv.org/abs/2408.03314
- **Pan et al. (2023/24)** — *Automatically Correcting Large Language Models: Surveying the Landscape of Diverse Automated Correction Strategies.* `arXiv:2308.03188`, TACL 2024. https://aclanthology.org/2024.tacl-1.27/ — the field survey/taxonomy.

Supporting / practitioner consensus:
- **Lee et al. (2025)** — *RefineBench: Evaluating Refinement Capability of Language Models via Checklists.* `arXiv:2511.22173` — guided >> unguided refinement; early over-confident halting.
- *Understanding the Dark Side of LLMs' Intrinsic Self-Correction* (2024), `arXiv:2412.14959` — over-thinking spirals, refinement/recency bias.
- **Kadavath et al. (2022)** — *Language Models (Mostly) Know What They Know* — calibration of self-reported confidence (why "looks good to me" is a weak stop signal).
- Practitioner field reports (agent self-critique loops): convergent guidance on k=1–2 default, cap at 3, `no_progress` + `loop_detected` halts, separating actor/critic.

> **Caveat on sourcing:** a 2026-dated "optimal stopping / GainNet" result surfaced during research on a non-standard preprint host and could not be corroborated against a recognized venue; it was **deliberately excluded**. The stop-condition guidance here rests on the peer-reviewed primaries above plus broadly convergent practitioner consensus, not on that item.
