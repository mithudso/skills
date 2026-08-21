<!-- hub-reference-banner -->
> **Reference file — part of the `ai-llm-model-layer` hub.** Formerly the standalone `llm-alignment-post-training` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

<!-- Provenance: reference under the `ai-agent-engineering` hub. Created 2026-05-31 via /dr deep-research from primary sources — DPO (Rafailov et al. 2023), the DPO survey (arXiv:2410.15595), SimPO (NeurIPS 2024), Constitutional AI (Bai et al. 2022) + the RLHF Book (Nathan Lambert), RewardBench (arXiv:2403.13787), Length-Controlled AlpacaEval (arXiv:2404.04475), ODIN/InfoRM reward-hacking papers, Self-Rewarding LMs (arXiv:2401.10020), OpenRLHF (arXiv:2405.11143), and the HuggingFace TRL + alignment-handbook docs. Scope: the LLM MODEL LAYER — how a pretrained model is turned into an aligned, instruction-following assistant via SFT + preference optimization. NOT agent/application engineering (the rest of this hub), NOT reasoning-RL (GRPO/RLVR/R1 → "reasoning models" skill), NOT PEFT mechanics (LoRA/QLoRA → PEFT skill). -->

# LLM Alignment & Post-Training

Turning a pretrained base model into a helpful, harmless, honest, instruction-following assistant. **Pretraining** gives a model knowledge and next-token fluency; it does not make the model *do what you ask* or *refuse what it should*. **Post-training** (a.k.a. alignment) is the stack of techniques — supervised fine-tuning, reward modeling, RLHF/RLAIF, and the direct-preference-optimization family — that closes that gap. This is the model layer that sits *under* everything else in this hub: the agent frameworks, RAG, and prompting references all assume a model that has already been aligned.

## Scope boundary (read first)

- **This reference** = the alignment / post-training of a model: SFT as the post-training base, reward modeling, RLHF (PPO), RLAIF & Constitutional AI, the DPO-variant family (DPO/IPO/KTO/ORPO/SimPO/CPO), preference-data pipelines, reward hacking, and alignment evaluation.
- **Reasoning-specific RL** — GRPO, RLVR (reinforcement learning with *verifiable* rewards), the DeepSeek-R1 recipe, process/outcome reward models for chains of thought — is a **separate concern**. It optimizes correctness against a verifier, not human preference against a reward model. See §10 for the one-paragraph pointer; defer depth to a dedicated reasoning-models skill.
- **Parameter-efficient fine-tuning mechanics** — LoRA, QLoRA, adapters, the rank/alpha/target-module knobs — are orthogonal plumbing that *any* of these stages can run on top of. Defer to a PEFT skill; this reference only notes *where* PEFT plugs in.
- **Model selection / the landscape** (which Claude/GPT/Gemini/open model to use) → `references/llm-models.md`. **Prompting** an already-aligned model → `references/prompt-engineering.md`. **Offline benchmark eval** (HELM/MMLU leaderboards) → `da-7-machine-learning`; alignment-specific *win-rate / reward-model* eval lives here (§9).

---

## 1. The post-training pipeline — where alignment fits

Modern post-training is a **multi-stage pipeline**, not a single step. The canonical shape (the one HuggingFace's alignment-handbook and most open recipes follow) is:

```
[pretrained base] → SFT (instruction tuning) → preference optimization (DPO or RLHF) → [aligned chat model]
```

Frontier recipes add stages around this skeleton. HuggingFace's **SmolLM3** recipe (July 2025) extends the standard two-stage `SFT → DPO` pipeline with a **mid-training** stage focused on reasoning before the preference phase. Many labs run **iterative** preference rounds (collect data with the current model, optimize, repeat). The two paradigms differ fundamentally:

- **SFT is off-policy / imitation:** it copies fixed expert demonstrations via cross-entropy. Cheap, stable, but capped by the quality of the demonstrations and prone to *exposure bias* (the model is only ever trained on gold prefixes, never its own).
- **Preference optimization is on-policy-ish / corrective:** it learns from *relative* judgments about the model's *own* (or near-own) outputs, which is what lets it exceed demonstration quality and suppress behaviors no demonstration explicitly showed.

Key mental model: **SFT teaches the format and the floor; preference optimization shapes the ceiling and the refusals.** You almost always need both.

---

## 2. Supervised fine-tuning (SFT) / instruction tuning — the base

SFT is the **first and load-bearing** post-training stage. You fine-tune the base model on a corpus of high-quality `(prompt, response)` pairs formatted with a **chat template** (the special tokens that delimit system/user/assistant turns — e.g. `<|user|> … <|assistant|> …`). The loss is plain **cross-entropy on the response tokens** (typically masking the prompt tokens so the model is graded only on what it should *generate*, not on echoing the input).

What SFT buys you:

- **Instruction-following behavior** — the model learns the *task* of responding to instructions, not just continuing text.
- **The output format / chat persona** — turn structure, tone, refusal style, tool-call syntax.
- **A strong reference point** — the SFT model becomes the `π_ref` (reference policy) that the next stage's KL penalty or DPO ratio is measured against.

**Data quality dominates.** A consistent 2024–2025 finding: for SFT, *quality and diversity of demonstrations beat raw quantity*. The "LIMA" thesis ("Less Is More for Alignment") and follow-ups show a few thousand carefully curated examples can rival much larger noisy sets. Failure modes to watch:

- **Alignment tax** — aggressive SFT (or preference optimization) can *degrade* general capabilities (reasoning, world knowledge) while improving helpfulness. Mitigations include mixing in pretraining/capability data and methods that blend objectives (see ORPO/CPO in §8).
- **Diversity collapse** — SFT can narrow the output distribution; recent work explicitly regularizes to *preserve diversity* during SFT.
- **Template mismatch** — train and serve must use the *exact* same chat template, or behavior silently degrades. This is a top operational footgun.

> SFT runs fine under PEFT (LoRA/QLoRA) — most open SFT today is LoRA. The PEFT mechanics are out of scope here (→ PEFT skill); just know SFT is the natural place to apply it.

---

## 3. Reward modeling (RM) — turning preferences into a scalar

Classic RLHF needs a **reward model**: a function `r(x, y)` that scores how good response `y` is for prompt `x`, learned from human preference comparisons. The standard recipe:

1. Collect **pairwise preferences**: for a prompt `x`, show annotators two responses `y_w` (chosen/winner) and `y_l` (rejected/loser); they pick the better one.
2. Train the RM (usually the SFT model with the LM head replaced by a **scalar value head**) under the **Bradley–Terry** objective: maximize `log σ(r(x, y_w) − r(x, y_l))`. The RM learns to assign the winner a higher score than the loser.

The RM is the **proxy for human judgment** that the RL stage optimizes against. Its quality bounds everything downstream — a biased RM produces a biased policy (this is the root of reward hacking, §9).

**Why comparisons, not absolute scores?** Humans are noisy and inconsistent at assigning absolute 1–10 ratings, but far more reliable at *relative* "A vs B" judgments. So preference data is overwhelmingly pairwise. Reward models are evaluated on their own benchmark, **RewardBench** (arXiv:2403.13787), which scores how often an RM prefers the known-better response across chat, safety, reasoning, and adversarial categories.

DPO and its family (§7–8) **collapse this step**: they show you can optimize the policy directly from preferences without ever materializing an explicit RM — the policy *implicitly* defines its own reward.

---

## 4. RLHF with PPO — the full loop

RLHF (Reinforcement Learning from Human Feedback) optimizes the policy to maximize the reward model's score, using **Proximal Policy Optimization (PPO)** as the RL algorithm. PPO-for-RLHF juggles **four models** simultaneously:

| Model | Role |
| --- | --- |
| **Policy** (actor) | The model being trained; generates responses. |
| **Reference** (`π_ref`, frozen SFT model) | Anchor for the KL penalty — keeps the policy from drifting too far. |
| **Reward model** (frozen) | Scores generated responses → the reward signal. |
| **Value model** (critic) | Estimates expected future reward per token; provides the baseline for advantage estimation (PPO is actor-critic). |

The loop, per step: the policy generates completions → the **reward model** scores them → a **per-token KL penalty** is subtracted from that reward → the **value model** estimates advantages → PPO updates the policy with its clipped surrogate objective.

The objective is, in words: **maximize reward − β · KL(policy ‖ reference)**. The **KL penalty** is the crux of stability — without it the policy will "run away" to whatever degenerate text maxes out the (imperfect) reward model. β trades off "satisfy the reward" against "stay close to the sensible SFT model."

**Why PPO is painful in practice:** four models in memory, sensitive hyperparameters (β, clip range, KL target, value-loss coefficient), and reward/value bookkeeping that is easy to get subtly wrong. The community-documented "N implementation details of RLHF with PPO" (ICLR Blog 2024) and "Secrets of RLHF" (Parts I & II) exist precisely because the gap between the paper and a working loop is large. This operational cost is the single biggest reason DPO took over much of the field.

---

## 5. RLAIF — reinforcement learning from AI feedback

**RLAIF** replaces (or supplements) the human preference labels with **labels generated by an LLM judge**. The pipeline is otherwise identical to RLHF: an LLM is shown two completions and asked which is better; those AI judgments train the reward model (or feed DPO directly). The term was introduced in Anthropic's Constitutional AI work and has since become a **default** because it scales: you are no longer rate-limited by human annotation throughput.

The core trade-off (Nathan Lambert's RLHF Book frames it crisply):

- **Human data: high-noise, low-bias.** Expensive, slow, inconsistent — but no systematic slant.
- **Synthetic/AI data: low-noise, high-bias.** Cheap, fast, consistent — but inherits the judge model's systematic biases (notably **self-preference bias**: models favor their own outputs).

Academic results show synthetic preference data performs *comparably* to human data, yet frontier labs still treat large human-preference sets as a competitive moat. The pragmatic frontier pattern is **hybrid routing**: send the bulk of easy comparisons to an AI judge, route the *hard / high-stakes* comparisons to humans.

---

## 6. Constitutional AI (CAI) / RLCAI

**Constitutional AI** (Bai et al., Anthropic, 2022) is the technique that kickstarted RLAIF and is the central alignment method behind every Claude model. It encodes human guidance as a **constitution** — a set of written natural-language principles ("Is the response harmful, unethical, or encouraging violence?") — and uses the *model itself* to apply them. Two phases:

1. **Supervised phase (critique → revision).** Sample a response from the model to a (often adversarial) prompt → ask the model to **critique** its own response against a randomly sampled constitutional principle → ask it to **revise** accordingly. Repeat across principles to get a refined `{y⁰, y¹, …, yⁿ}`. Fine-tune the model on the final **revised** responses. This produces a model that is harmless *without* a human ever labeling a harmful example.

2. **RL phase (RLAIF).** Generate response pairs from the SL-finetuned model → an **AI feedback model**, shown the prompt + the principles + both completions, picks which is "both higher quality and more aligned with the stated principle" → train a **preference model** on these AI labels → run RL (PPO) against it. The RL-from-AI-feedback half is the canonical **RLAIF / RLCAI**.

The constitution is the elegant part: a short, auditable, *editable* document is the only human input, and it steers behavior at scale through self-critique and AI preference labeling. Open recipes (e.g. via TRL) can reproduce the critique-revision data generation and feed it into SFT + DPO.

---

## 7. DPO — Direct Preference Optimization

**DPO** (Rafailov et al., 2023) is the pivot point of modern alignment. Its insight: the RLHF objective ("maximize RM reward subject to a KL constraint") has a **closed-form optimal policy**, so you can skip the reward model *and* the RL loop entirely and optimize the policy **directly** on preference pairs with a simple classification loss.

The DPO loss, in words: for each `(x, y_w, y_l)`, push up the policy's log-prob of the winner *relative to the reference* and push down the loser's, through a sigmoid:

```
L_DPO = −E[ log σ( β·[log π_θ(y_w|x) − log π_ref(y_w|x)]
                  − β·[log π_θ(y_l|x) − log π_ref(y_l|x)] ) ]
```

Properties:

- **Needs a reference model** `π_ref` (the frozen SFT model) — but only for a forward pass, no value model, no RM, no sampling loop.
- **Needs paired preference data** `(y_w, y_l)`.
- **Offline** — trains on a static, pre-collected dataset (though **iterative/online DPO** variants re-collect data with the updated policy each round; see §8 and the on/off-policy note in §9).
- **β** plays the KL-temperature role: small β = stay close to reference, large β = trust preferences more.

DPO is now the default first-choice for open-model alignment because it is dramatically simpler and cheaper than PPO while reaching comparable quality on most chat benchmarks. The standard open pipeline is `SFT → DPO`.

---

## 8. The DPO-variant family — IPO, KTO, ORPO, SimPO, CPO

DPO has spawned a family of variants, each changing one design choice to fix a specific failure mode. The survey arXiv:2410.15595 catalogs dozens; these are the load-bearing ones and **when each wins**:

| Variant | What it changes vs DPO | Ref model? | Data | Folds in SFT? | When it wins |
| --- | --- | --- | --- | --- | --- |
| **DPO** | baseline (sigmoid/BCE on the log-ratio margin) | Required | Paired | No | Default; you have clean pairwise data and an SFT checkpoint. |
| **IPO** (Identity PO) | Replaces the sigmoid/BCE loss with a **squared (regression) loss** targeting a fixed margin; weaker preference assumption | Required | Paired | No | DPO is **overfitting** to (especially deterministic/near-deterministic) preferences; you want regularization against over-optimization. |
| **KTO** (Kahneman–Tversky) | Treats each response **individually** with a prospect-theory value function (loss-aversion: bad answers penalized more than good rewarded) | Required | **Unpaired** (binary good/bad labels) | No | You only have **thumbs-up/down** signals, not pairs — far cheaper to collect; or your good:bad ratio is skewed. |
| **ORPO** (Odds-Ratio PO) | **Single-stage**: combines the SFT (NLL on winner) loss with an **odds-ratio** preference penalty — *no reference model at all* | **None** | Paired | **Yes** (replaces SFT) | You want **one stage instead of two** (no separate SFT pass), no reference-model memory, and reduced alignment tax. Good for skewed (e.g. 1:10) class ratios. |
| **SimPO** (Simple PO) | **Reference-free**; uses the **length-normalized average log-prob** as the implicit reward + a target reward margin γ | **None** | Paired | No | You want the simplest, most memory-efficient setup *and* built-in **length-bias** resistance. Reported +6.4 pts on AlpacaEval 2 / +7.5 on Arena-Hard over DPO. |
| **CPO** (Contrastive PO) | Adds a **behavior-cloning / NLL term** on the winner alongside a reference-free contrastive preference loss | Partial/None | Paired | **Yes** | You want to keep the policy anchored to the preferred-data distribution (preserve capabilities, fight alignment tax) without a full reference pass. Originated for machine-translation alignment. |

**Practical guidance** (synthesizing the survey + practitioner reports):

- Start with **DPO** if you have clean pairs and an SFT model. It is the most studied and best-tooled.
- Drop to **SimPO** when memory is tight or length-gaming is a problem (it removes the reference model and normalizes by length).
- Use **KTO** when you can only get binary feedback (much cheaper to collect than pairwise) or your data is heavily imbalanced.
- Use **ORPO** to collapse SFT+preference into **one** stage and skip the reference model entirely.
- Use **IPO** when DPO over-optimizes / collapses on your preference set.
- **Caveat on rankings:** controlled studies (e.g. arXiv:2603.19335) find that the *relative* ranking of these methods can **invert with model scale and data** — there is no universal winner. Pick based on your data shape (paired vs binary, balanced vs skewed) and your compute, then **measure** (§9), don't trust a leaderboard from a different scale.

**Online / iterative variants** sit on top of any of these: **iterative DPO**, **Online AI Feedback (OAIF)**, and **Self-Rewarding Language Models** (arXiv:2401.10020, where the model judges its *own* outputs to generate fresh preference pairs each round) move the data from static-offline toward on-policy, which closes much of DPO's quality gap with PPO.

---

## 9. Reward hacking, length bias, and over-optimization

The defining pathology of preference optimization: the policy learns to exploit **flaws in the reward signal** (explicit RM, or DPO's implicit reward) to score high *without* genuinely being better. The most notorious instance is **length bias** — reward models systematically prefer **longer** responses regardless of quality, so the policy learns to pad. **Over-optimization** (a.k.a. reward over-optimization, Goodhart's law) is the general case: as you optimize harder against an imperfect proxy, true quality eventually *declines* even as proxy reward keeps rising.

Mitigations, by layer:

- **Length-bias specifically:**
  - **ODIN** (arXiv:2402.07319) — train the RM with two heads, one correlated with length and one *de*correlated; discard the length head at RL time.
  - **R-DPO / LD-DPO** — add an explicit length penalty or length-discount to the DPO loss.
  - **SimPO** — length-normalize the implicit reward by construction (§8).
  - **FiMi-RM / bias-fitting** (arXiv:2505.12843) — debias the RM's length-reward distribution.
- **Over-optimization generally:**
  - **KL penalty / β** — the front-line defense in PPO and DPO: staying near the reference policy bounds how far you can chase the proxy.
  - **Reward-model ensembles** — average several RMs so the policy can't fool all of them on a spurious feature.
  - **InfoRM** (NeurIPS 2024) — information-theoretic RM that filters out preference-irrelevant signal and provides a **CSI** indicator to *detect* over-optimization.
  - **Bayesian / non-negative reward modeling** (BNRM) — more robust under distribution shift, more interpretable reward decomposition.
  - **IPO** — its regression loss is more resistant to over-optimizing deterministic preferences than DPO's sigmoid.
- **Evaluation-side hardening (don't fool yourself):**
  - **On/off-policy matters** — pure offline preference data drifts from what the policy actually generates; **on-policy / iterative** collection (re-sampling from the current model) reduces the proxy gap.
  - **Annotator disagreement** is signal: filter or down-weight low-agreement pairs (inter-annotator agreement on general chat preferences is ~73%).

---

## 10. Evaluating aligned models

Alignment is judged primarily by **win rate** — how often the aligned model's response is preferred over a baseline's, scored by an LLM-as-judge calibrated to humans.

- **AlpacaEval 2.0** — 805 fixed instructions; a GPT-4-class judge compares your model vs a reference (e.g. GPT-4-Turbo) head-to-head and reports the preference probability. Cheap (<$10, <3 min) and fast.
- **Length-Controlled AlpacaEval (LC-WR)** (arXiv:2404.04475) — **the standard fix for length bias in eval itself**. Regresses out the length effect so a verbose model can't win by padding. Raised Spearman correlation with LMSYS **Chatbot Arena** from 0.94 → 0.98. **Always report LC win-rate, not raw.**
- **Arena-Hard(-Auto)** — a harder, more separable auto-benchmark built from real Chatbot-Arena prompts; better at distinguishing strong models than MT-Bench / vanilla AlpacaEval.
- **Chatbot Arena (LMSYS)** — crowd-sourced human pairwise battles → Elo. The closest thing to ground truth for general preference; the auto-evals above are validated *against* it.
- **RewardBench** (§3) — evaluates the **reward model** itself, decoupled from the policy.
- **Safety evals** — alignment is not just helpfulness. Run dedicated harmlessness/refusal and red-team suites (e.g. harmful-instruction refusal sets, jailbreak resistance) so you don't ship a model that is helpful *and* unsafe.

**Eval is gameable too.** "Cheating Automatic LLM Benchmarks" (arXiv:2410.07137) shows a **null model** emitting a crafted constant string can hit ~86% win rate on AlpacaEval 2.0 — a reminder that auto-judges have exploitable biases, and that win-rate must be paired with human spot-checks and held-out safety/capability tests.

> **Reasoning-RL boundary:** if your goal is *correctness on verifiable tasks* (math, code) rather than *human preference*, you are in **RLVR / GRPO** territory — GRPO (DeepSeek-R1) drops the value model and estimates the baseline from a *group* of sampled answers scored by a rule-based **verifier**, not a learned RM. That is a different optimization target (verifiable reward, not preference) and belongs to a dedicated reasoning-models skill. Mentioned here only so you route it correctly; do **not** treat GRPO as just another DPO variant.

---

## 11. Tooling stack

| Tool | What it is | Use it when |
| --- | --- | --- |
| **HuggingFace TRL** | The de-facto full-stack post-training library. Ships trainers for every stage: `SFTTrainer`, `RewardTrainer`, `DPOTrainer`, `PPOTrainer`, `GRPOTrainer`, `KTOTrainer`, `ORPOTrainer`, `CPOTrainer`, `OnlineDPOTrainer`, plus `RLOO`, `BCO`, `GKD`, `Nash-MD`, `XPO`, `PRM` trainers. Integrates **PEFT** (LoRA/QLoRA) and **Accelerate**/DeepSpeed for single-GPU → multi-node. | Almost always the right default for SFT + preference optimization. Match the stage to the trainer. |
| **alignment-handbook** (HuggingFace) | Opinionated, reproducible **recipes** (config + scripts) on top of TRL — the published `SFT → DPO` pipelines behind Zephyr and SmolLM3. | You want a known-good, end-to-end recipe rather than wiring trainers yourself. Start here, then customize. |
| **Axolotl** | Config-file-driven fine-tuning wrapper (YAML) over the HF stack; broad model/dataset/format support, supports SFT and DPO/preference methods. | You want declarative, low-code fine-tuning runs and easy dataset format handling. |
| **OpenRLHF** | High-performance, **scalable RLHF** framework combining **Ray + DeepSpeed ZeRO-3 + vLLM**. Implements PPO, DPO, KTO, conditional SFT, rejection sampling. | You are running **large-scale online RLHF/PPO** (the regime where TRL's single-process PPO struggles) and need distributed rollout + training. |
| **NeMo-Aligner / verl / TRLX** | Other production-scale alignment frameworks (NVIDIA NeMo; volcengine verl; CarperAI TRLX). | Large-scale or vendor-aligned stacks where you already live in that ecosystem. |

PEFT integration means every stage above can run with **LoRA/QLoRA** — most open alignment today is LoRA-based for cost. (PEFT mechanics → PEFT skill.)

---

## Anti-patterns

- **Skipping SFT and going straight to DPO/PPO on a base model.** Preference optimization assumes an instruction-following starting point; without SFT the reference policy is garbage and training is unstable.
- **Train/serve chat-template mismatch.** Using a different template (or none) at inference than at SFT silently wrecks quality. Pin the template.
- **Reporting raw win-rate.** Always length-control (LC-WR). A verbose model wins on raw AlpacaEval for the wrong reason.
- **Trusting one auto-eval / one leaderboard.** Auto-judges are gameable (null-model attack), variant rankings invert with scale, and self-preference bias inflates a same-family judge. Triangulate: LC-WR + Arena-Hard + held-out safety + human spot-checks.
- **Optimizing harder to fix a plateau.** Past a point, more steps against an imperfect reward = over-optimization: proxy reward up, real quality down. Watch a held-out reward-model / KL budget and stop early.
- **Treating GRPO/RLVR as a DPO variant.** Different target (verifiable correctness vs human preference). Route to the reasoning-models skill.
- **Assuming synthetic preference data is free of bias.** It's low-noise but high-bias (self-preference). Route hard/high-stakes comparisons to humans.

---

## References

- Rafailov et al., *Direct Preference Optimization: Your Language Model is Secretly a Reward Model* (2023) — the DPO paper.
- *A Comprehensive Survey of Direct Preference Optimization: Datasets, Theories, Variants, and Applications* — arXiv:2410.15595 (the DPO-variant taxonomy).
- Meng, Xia, Chen, *SimPO: Simple Preference Optimization with a Reference-Free Reward* — NeurIPS 2024 (https://github.com/princeton-nlp/SimPO).
- Bai et al. (Anthropic), *Constitutional AI: Harmlessness from AI Feedback* — arXiv:2212.08073; https://www.anthropic.com/research/constitutional-ai-harmlessness-from-ai-feedback.
- Nathan Lambert, *The RLHF Book* — Constitutional AI / synthetic-data chapter: https://rlhfbook.com/c/13-cai.
- *Secrets of RLHF in Large Language Models* — Part I: PPO (arXiv:2307.04964); Part II: Reward Modeling (arXiv:2401.06080).
- *The N Implementation Details of RLHF with PPO* — ICLR Blogposts 2024.
- Lambert et al., *RewardBench: Evaluating Reward Models for Language Modeling* — arXiv:2403.13787.
- Dubois et al., *Length-Controlled AlpacaEval: A Simple Way to Debias Automatic Evaluators* — arXiv:2404.04475; https://github.com/tatsu-lab/alpaca_eval.
- Li et al., *From Crowdsourced Data to High-Quality Benchmarks: Arena-Hard and BenchBuilder* — arXiv:2406.11939.
- Chen et al., *ODIN: Disentangled Reward Mitigates Hacking in RLHF* — arXiv:2402.07319.
- Miao et al., *InfoRM: Mitigating Reward Hacking in RLHF via Information-Theoretic Reward Modeling* — NeurIPS 2024 (arXiv:2402.09345).
- *Bias Fitting to Mitigate Length Bias of Reward Model in RLHF* (FiMi-RM) — arXiv:2505.12843.
- Yuan et al., *Self-Rewarding Language Models* — arXiv:2401.10020.
- Hu et al., *OpenRLHF: An Easy-to-use, Scalable and High-performance RLHF Framework* — arXiv:2405.11143.
- *Do Post-Training Algorithms Actually Differ? … Scale-Dependent Ranking Inversions* — arXiv:2603.19335.
- *Cheating Automatic LLM Benchmarks: Null Models Achieve High Win Rates* — arXiv:2410.07137.
- HuggingFace **TRL** docs (https://huggingface.co/docs/trl) and **alignment-handbook** (https://github.com/huggingface/alignment-handbook).
