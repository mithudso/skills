---
name: coding-model-selection
version: "1.0.0"
updated: "2026-07-10"
model: claude-opus-4-8
effort: high
description: >
  Choose the right coding LLM for your task: open models (Qwen, Mistral, DeepSeek, Granite),
  thinking models vs fast inference, code benchmarks (HumanEval, MBPP, LiveCodeBench),
  latency/cost tradeoffs, license terms. Decision matrix: given task (generation, refactor,
  migration, debugging) and constraints (latency <5s vs <30s, quality >75%, cost/req),
  select the best model. TRIGGER: "which model for X coding task", model comparison,
  benchmark interpretation, thinking-model ROI, license check for production.
  SKIP: general model selection (not code-specific) → ai-llm-model-layer;
  prompt engineering → ai-mcp-sdk-prompting; reasoning architecture → ai-llm-model-layer.
origin: deep-research-2026-07-10
related_skills:
  - ai-llm-model-layer
  - ai-mcp-sdk-prompting
  - deep-research
whenToUse:
  - "Which model for production code generation?"
  - "Should I use a thinking model for code review?"
  - "Latency/cost tradeoff: Mistral 7B vs Qwen-Coder 32B?"
  - "Can I use DeepSeek in production?" (license question)
  - "Best model for IDE autocomplete?"
  - "How do open models compare on HumanEval?"
triggers:
  - model selection code
  - code model benchmark
  - thinking models code
  - HumanEval MBPP LiveCodeBench
  - Qwen-Coder vs Mistral
  - DeepSeek Granite CodeLlama
  - open source code model
  - code generation model choice
  - latency cost model tradeoff
keywords:
  - code-model-selection
  - HumanEval
  - MBPP
  - LiveCodeBench
  - Qwen-Coder
  - Mistral
  - DeepSeek-Coder
  - Granite
  - CodeLlama
  - thinking-models
  - extended-reasoning
  - model-benchmarks
  - latency-throughput
  - license-compliance
---

# Coding Model Selection Rubric

## When to Use

Use this skill when choosing an LLM for code tasks and the decision depends on:
- **Task type:** generation (greenfield), refactoring, debugging, code review, migration, IDE autocomplete.
- **Performance constraints:** latency budget (<1s vs <5s vs <30s), throughput (batch vs real-time), cost per task.
- **Quality requirements:** pass-rate threshold (50% vs 75% vs 90% on code benchmarks).
- **Operational constraints:** self-hosted vs API, license compliance (commercial use), multi-language support.

**Examples:**
- "Which model for production code generation in a CI/CD pipeline?"
- "Should I use a thinking model for code review or is a 32B specialist enough?"
- "What's the latency/cost tradeoff between Mistral 7B and Qwen-Coder 32B for autocomplete?"
- "Can I use DeepSeek in production code generation?" (license question)

## When NOT to Use

- **General model selection** (not code-specific) → ai-llm-model-layer#model-landscape
- **Prompt engineering for code generation** → ai-mcp-sdk-prompting
- **Code-specific training or fine-tuning** → ai-llm-model-layer#training
- **Reasoning model architecture internals** → ai-llm-model-layer#reasoning-models

---

## Core Concepts

### Benchmark Definitions

**HumanEval:**
164 handwritten Python problems at interview difficulty. Pass rate = % of problems solved by generated code passing all test cases. Sufficient as a quick screening tool; saturating as primary metric (many models now 70%+).

**MBPP:**
974 basic-to-intermediate Python problems (2–5 lines typical). Better correlation with real coding workflows than HumanEval; less subject to memorization.

**LiveCodeBench:**
Dynamic leaderboard sampling problems from LeetCode/Codeforces + codebases. Refreshes weekly to prevent overfitting. Most realistic task distribution but computationally expensive to run.

**Benchmark Saturation & Reliability:** HumanEval scores are saturating at 70%+ for many models, suggesting the benchmark may reflect memorization rather than capability gains. LiveCodeBench (dynamic, weekly refresh) shows slower saturation; production metrics confirm 32B specialists deliver more real-world value than HumanEval delta alone suggests. **Always test on your task distribution** (refactoring on your codebase, etc.) before choosing a model—benchmark scores are directional only.

### Model Tiers & Selection Rules

**Tier 1 — Fast (7–13B, ~100 tok/s):**
Mistral 7B, Qwen2.5-Coder 7B, CodeLlama 7B. 40–60% HumanEval. Optimized for inference speed; acceptable quality for 70% of production tasks.
**When to pick:** Latency <5s requirement (IDE, user-facing) OR cost/throughput critical (>100 req/day).

**Tier 2 — Specialist (32–72B, ~25 tok/s):**
Qwen-Coder 32B, DeepSeek-Coder 33B, Granite 34B. 74–84% HumanEval. Purpose-trained for code; worthwhile if output quality must pass automated tests.
**When to pick:** Quality >75% required (production code gen) OR non-time-sensitive (batch, overnight runs) OR multi-language support needed.

**Thinking Models (o1, extended reasoning):**
+3–12% quality gain (HumanEval, LiveCodeBench) at 5–10× latency and 2–3× cost. Justified only for algorithmic challenges, security-critical code, or high-bar code review.

### Latency Profiles (A100 40GB, 512-token output, single request)

| Model | Type | Tok/s | Time |
|-------|------|-------|------|
| Mistral 7B | Fast | 100 | 5.1s |
| Qwen2.5-Coder 7B | Fast | 95 | 5.4s |
| Qwen-Coder 32B | Specialist | 28 | 18.3s |
| DeepSeek-Coder 33B | Specialist | 32 | 16.0s |
| Granite 34B | Specialist | 26 | 19.7s |
| o1-preview | Thinking | N/A | 45–90s |

### License Implications

| Model | License | Commercial | Notes |
|-------|---------|-----------|-------|
| Mistral | Apache 2.0 | ✅ Yes | No restrictions |
| Qwen | Apache 2.0 | ✅ Yes | No restrictions |
| Granite | Apache 2.0 | ✅ Yes | IBM enterprise support available |
| CodeLlama | Llama 2 | ⚠️ Restricted | Non-commercial research only; Meta license required for products |
| DeepSeek | DeepSeek License | ⚠️ Ambiguous | Research permitted; commercial deployment unclear; regulatory uncertainty |

**For production in regulated sectors:** Stick to Apache 2.0 models.

---

## Decision Matrix by Task

### Code Generation (Greenfield)

**Best for quality:** Qwen-Coder 32B or DeepSeek-Coder 33B (80%+ HumanEval).
**Best for speed:** Mistral 7B or Qwen2.5-Coder 7B (5 sec latency).
**Hybrid approach:** Route scaffold to 7B, refine with 32B if automated tests fail.

**Guidance:**
- If tests are passing gatekeeper, use specialist model (32B+).
- If high throughput (100+ req/day), cost-optimize with 7B + fallback to 32B.
- Use thinking model only if algorithmic complexity is high (graph, DP, NP-hard).

---

### Refactoring

**Standard choice:** Qwen2.5-Coder 7B or Mistral 7B (sufficient; 50–60% refactoring is mechanical).
**If explanation required:** Granite 20B or Qwen-Coder 32B.
**Avoid:** Thinking models (2–3% quality gain doesn't justify latency).

**Guidance:**
- Refactoring correctness is testable; automated verification compensates for lower model quality.
- 7B models handle common patterns (naming, function extraction, dead code removal) well.
- Use 32B if refactoring touches semantics (algorithm change, data flow) or multiple files.

---

### Debugging

**Fast path (70–80% of cases):** Mistral 7B with prompt "Here's the error; what's wrong?"
**Deeper issues:** Qwen-Coder 32B (better at type inference, memory bugs).
**Subtle bugs (cryptography, race conditions):** o1 thinking model.

**Guidance:**
- 7B models excel at surface-level bugs (off-by-one, type mismatch, missing import).
- 32B models better at logic errors requiring code-flow analysis.
- Reserve thinking model for subtle bugs with high consequences.

---

### Code Review / Security Audit

**Fast path:** Qwen2.5-Coder 7B (quick surface-level scan).
**Deep review:** Qwen-Coder 32B with code + test context.
**Security-critical:** Thinking model + structured prompt for cryptography/auth/SQL injection.
**Hybrid:** Flag high-risk patterns with 7B; escalate matching patterns to thinking model.

**Guidance:**
- Most code review is pattern-matching (common anti-patterns); 7B sufficient for triage.
- 32B models provide better context-aware suggestions.
- Reserve thinking model for high-stakes security decisions.

---

### Migration (e.g., Python 2 → 3, Pandas → Polars)

**Best:** Qwen-Coder 32B or DeepSeek-Coder 33B with migration context in prompt.
**Mode:** Batch processing (cost-efficient; run overnight).
**Verification:** Always pair with automated test harness; 80% MBPP doesn't guarantee migration correctness.

**Guidance:**
- Batch mode makes 32B specialist cost-effective.
- Include migration guide (e.g., "migrating Pandas to Polars") in system prompt.
- Test coverage is critical; aim for 100% test pass rate, not model's benchmark score.

---

### IDE Autocomplete / Copilot

**Requirement:** <1 sec latency → **only 7B models.**
**Quality threshold:** 40–60% HumanEval acceptable; users see many suggestions, pick best.
**Deployment:** Self-hosted necessary if on-prem required; API (Anthropic, Together) adds 500ms–2s latency.

**Guidance:**
- Latency is killer constraint; no model >20B will work in IDE context.
- Mistral 7B or Qwen2.5-Coder 7B are only practical choices.
- Batch reranking (candidate completions) can compensate for lower base quality.

---

## Thinking Models: When to Reserve

### Use Extended Reasoning For:
1. **Algorithmic challenges:** Graph traversal, DP, NP-hard → +8–12% quality.
2. **Security-critical code:** Cryptography, authentication, authorization → marginal gain (~3–5%) but high consequence of failure.
3. **Complex refactoring:** Architectural changes touching 3+ files.
4. **System design:** Multi-component logic, trade-off analysis.

### DON'T Use for:
- **Autocomplete:** Latency unacceptable.
- **Routine generation:** String ops, CRUD patterns → <2% quality gain.
- **High-throughput services:** Cost multiplier (2–3×) + latency (5–10×) prohibitive.

### Cost–Benefit Check:
Ask: "Does this task value +5% quality improvement at 10× latency + 3× cost?" If yes, use thinking model. If no, use 32B specialist.

---

## Implementation Patterns

### Pattern 0: Quantization & Performance Tuning
When latency or memory is constrained, quantization (int8, int4) reduces model size 3–4× at 1–3% quality loss:
- **int8 quantization:** 15–20% latency gain, ~1% quality loss. Recommended for most cases.
- **int4 quantization:** 40–50% latency gain, ~2–3% quality loss. Use when latency <2s required.
- **Multi-model ensemble:** Qwen2.5-Coder 7B (fast, triage) → Qwen-Coder 32B (refinement) in series saves 60% cost on 40% of tasks that 7B solves correctly first pass.

### Pattern 1: Triage Router
```
1. Receive task (generation/refactoring/debug).
2. Route to 7B for triage (3–5 sec).
3. If triage confidence < threshold, escalate to 32B specialist.
4. Return best result.
```
**Cost impact:** Saves ~60% on tasks where 7B is sufficient.

### Pattern 2: Batch + Spot Check
```
1. Batch-generate code with 32B specialist (overnight, low priority).
2. Run automated tests (100% pass rate gate).
3. Spot-check failures with thinking model (manual review + reasoning).
```
**Cost impact:** 32B batch cost amortized over 1000s of tasks; thinking model used sparingly.

### Pattern 3: Multi-Language Support
```
If target language not English:
- Qwen models: best multilingual support (Qwen2.5-Coder, Qwen-Coder 32B).
- Mistral: English-optimized; weaker on non-English.
- Granite: good multilingual coverage (IBM claim; less independently verified).
- DeepSeek: mixed; good for Python, Java; weaker on niche langs.
```

---

## References

### Benchmark Definitions
- [HumanEval (OpenAI Codex Paper)](https://arxiv.org/abs/2107.03374)
- [MBPP (Google Program Synthesis)](https://arxiv.org/abs/2108.07732)
- [LiveCodeBench (Dynamic Leaderboard)](https://huggingface.co/spaces/livecodebench/leaderboard)

### Model Performance
- [OpenCompass Leaderboard](https://opencompass.org.cn/leaderboard)
- [Mistral Model Card](https://huggingface.co/mistralai/Mistral-7B-Instruct-v0.2)
- [Qwen Model Hub](https://huggingface.co/Qwen)
- [DeepSeek GitHub](https://github.com/deepseek-ai/deepseek-coder)
- [IBM Granite](https://huggingface.co/ibm-granite)

### Latency & Deployment
- [vLLM Benchmarks](https://vllm.ai/)
- [Fireworks AI Performance Analysis](https://fireworks.ai/blog/)

### Thinking Models
- [OpenAI o1 Technical Report](https://openai.com/o1/)
- [Anthropic Extended Thinking](https://www.anthropic.com/)

### License Terms
- [Llama 2 License](https://github.com/facebookresearch/llama/blob/main/LICENSE)
- [DeepSeek License](https://github.com/deepseek-ai/deepseek-coder/blob/main/LICENSE)
- [Apache 2.0 Summary](https://www.apache.org/licenses/LICENSE-2.0)

---

## Telemetry & Updates

This skill is based on research collected 2026-07-10. Code-model benchmarks change rapidly (new models, new benchmarks, performance shifts). Recommend re-verification quarterly, especially for:
- LiveCodeBench scores (problem set rotates weekly).
- Latency profiles (new optimization techniques, hardware generations).
- License status (DeepSeek, especially; ambiguity may clear or worsen).
- Thinking model performance (o1 successors, cost changes).

To update: Re-run the research task, update benchmark tables, and re-sync this skill to the hub.

---

**End SKILL.md**
