---
name: code-benchmark-interpretation
version: "1.0.0"
updated: "2026-07-10"
description: >
  Interpret code-LLM benchmark scores and predict real-world performance.
  Covers HumanEval (164 problems, 70%+ saturation), MBPP (974 problems, real-workflow correlation),
  LiveCodeBench (dynamic weekly LeetCode/Codeforces sampling, overfitting detection).
  Explains benchmark bias, memorization risk, saturation plateaus, leaderboard methodology,
  confidence intervals, and task-distribution matching.
  
  Operator can read a model's scores on different benchmarks and infer:
  - Will this model work for MY code task?
  - Does this model overfit HumanEval?
  - When is LiveCodeBench more trustworthy than static benchmarks?
  - What confidence should I have in published leaderboard results?
  
  TRIGGER: "How do I read code benchmark scores?", "Is HumanEval saturated?",
  "Compare MBPP vs HumanEval", "Will this model work for my task?",
  "Benchmark memorization", "Leaderboard bias", "Should I use LiveCodeBench?",
  "Confidence interval on code model scores", "Task distribution mismatch",
  "Why do models plateau on HumanEval?", "OpenCompass vs CodeArena methodology".
  
  SKIP: general model evaluation (ai-llm-model-layer); choosing which model for production
  use without benchmark context (coding-model-selection); writing benchmark papers;
  building a code-execution test harness.

origin: ECC
related_skills:
  - ai-llm-model-layer
  - coding-model-selection
  - da-analytical-methods

whenToUse:
  - "Interpret benchmark scores for a code model"
  - "Should I trust this leaderboard result?"
  - "Will this model work for my coding task?"
  - "Compare models across different benchmarks"
  - "Understand benchmark saturation and overfitting"
  - "Evaluate real-world performance from scores"

triggers:
  - "code benchmark"
  - "HumanEval"
  - "MBPP"
  - "LiveCodeBench"
  - "benchmark saturation"
  - "benchmark memorization"
  - "leaderboard"
  - "OpenCompass"
  - "CodeArena"
  - "benchmark overfitting"
  - "benchmark bias"
  - "real-world performance"
  - "task distribution"

keywords:
  - code-benchmark
  - HumanEval
  - MBPP
  - LiveCodeBench
  - saturation
  - memorization
  - overfitting
  - leaderboard
  - benchmark-bias
  - confidence-interval
  - real-world-performance
  - task-distribution
  - OpenCompass
  - CodeArena
  - benchmark-methodology

references:
  - "HumanEval: Evaluating Large Language Models Trained on Code (OpenAI 2021)"
  - "MBPP: Problem Statement (Google ML-enhanced Search 2021)"
  - "LiveCodeBench: Holistic and Contamination-Free Evaluation of Large Language Models on Code (Li et al. 2024)"
  - "Saturation Analysis of Code Benchmarks"
  - "Memorization and Contamination in LLMs"
  - "Leaderboard Methodology Comparison (OpenCompass, CodeArena, HF Leaderboards)"
  - "Statistical Validity and Confidence Intervals in Benchmark Reporting"

---

# Code Benchmark Interpretation & Leaderboard Analysis

## Executive Summary

Code-LLM benchmarks (HumanEval, MBPP, LiveCodeBench) measure model capability on narrowly scoped programming tasks. Each has different saturation curves, memorization risk, task distributions, and real-world predictive value. A model's benchmark score alone **does not predict performance on your specific code task**. This skill teaches how to:

1. **Read benchmark scores with appropriate skepticism** — understand saturation, memorization risk, and confidence bounds
2. **Select the right benchmark for your evaluation goal** — HumanEval for broad capability, MBPP for production-like diversity, LiveCodeBench for overfitting detection
3. **Infer real-world performance** — match task distribution, account for domain shift, adjust for benchmark-specific bias
4. **Compare leaderboards critically** — understand methodology differences, reproducibility issues, statistical rigor

## 1. Benchmark Anatomy & Characteristics

### HumanEval (OpenAI 2021)

- **Size:** 164 problems (Python)
- **Task types:** Function implementation from docstring + test cases (~300 tokens avg)
- **Complexity range:** Small functions (50–1000 LOC) / competitive programming level
- **Saturation point:** ~70–74% (models cluster tightly above this); improvements become harder
- **Memorization risk:** HIGH — dataset was published in 2021, large models trained post-2022 likely memorized it
- **Real-world alignment:** Poor — too focused on short isolated problems; lacks refactoring, debugging, long-context editing, maintenance
- **Leaderboard bias:** Rewards short-solution elegance; penalizes verbose-but-correct code
- **Use this for:** Quick broad capability check; historical comparison; detecting catastrophic failures

### MBPP (Google, published 2021, extended 2023)

- **Size:** 974 problems (core: 500 + extended: 474)
- **Task types:** Diverse: string manipulation, algorithms, data structures, math, utilities
- **Complexity range:** 10–50 LOC typical; less competitive-programming focused than HumanEval
- **Saturation point:** ~70–75% (similar to HumanEval but slower approach due to task variety)
- **Memorization risk:** MEDIUM-HIGH — published dataset, but 5x larger makes memorization harder per problem
- **Real-world alignment:** Better — includes practical utility functions, library-like coding patterns
- **Leaderboard bias:** Slightly favors straightforward implementations; less emphasis on golf/elegance
- **Use this for:** Production readiness screening; more robust generalization signal than HumanEval alone

### LiveCodeBench (Li et al. 2024)

- **Size:** ~500 problems, weekly additions from LeetCode/Codeforces (growing)
- **Task types:** Real competitive programming + practical coding tasks
- **Complexity range:** Medium (LeetCode Easy–Medium level)
- **Saturation point:** Intentionally designed to avoid saturation; moves weekly
- **Memorization risk:** VERY LOW — dynamic dataset means models trained before week N cannot have memorized week N problems
- **Real-world alignment:** Medium — better than HumanEval, worse than production codebases; competitive-programming biased
- **Leaderboard bias:** Requires efficient solutions; harder for interpretability-first models
- **Use this for:** Overfitting detection; recency-valid model ranking; detecting train-test contamination

## 2. Saturation & Performance Plateaus

### Why Models Plateau on Static Benchmarks

- **Task distribution ceiling:** Once models understand the benchmark's distribution (e.g., "most solutions are ≤30 lines"), further improvements require architectural changes, not scale
- **Loss of signal:** At 70%+ accuracy, remaining failures are often edge cases or adversarial examples, not systematic capability gaps
- **Diminishing returns from scale:** Adding parameters or compute yields <1% accuracy gain per 10x scale at saturation
- **Benchmark contamination:** Post-2022 models likely saw HumanEval in training data → reported scores reflect partial memorization, not pure capability

### Saturation Across Benchmarks

| Benchmark | Saturation Onset | Typical "Top" Performance | Rate of Plateau |
|-----------|------------------|--------------------------|-----------------|
| HumanEval | ~65% | 72–75% | Steep (1 model per 1% improvement) |
| MBPP | ~68% | 75–78% | Moderate (harder due to diversity) |
| LiveCodeBench | Moving target | N/A (no plateau intended) | By design prevents saturation |

**Key insight:** If two models score 72% and 74% on HumanEval, the 2% gap is **not** reliable signal—might be memorization variance or test-case luck. Use MBPP or LiveCodeBench to disambiguate.

## 3. Memorization & Data Contamination Risk

### High-Risk Scenarios

1. **Model trained after benchmark publication** (HumanEval: 2021, MBPP: 2021, LiveCodeBench: weekly)
   - Check model training data cutoff date
   - If training data includes CommonCrawl or GitHub (post-benchmark-publication), memorization is plausible

2. **Benchmark problems appear in public GitHub repos** → crawled into training corpora
   - HumanEval: ~50% of problems appeared in GitHub repos by 2023 (Sobania et al.)
   - MBPP: slightly more dispersed, but still ~30–40% coverage

3. **Model trained on code-completion datasets curated from same domains** (LeetCode, Codeforces)
   - GPT-3.5+ and Claude-3+ models: trained on mixed public code corpora, likely overlap

### Mitigation Strategies

- **LiveCodeBench is your friend** — weekly updates mean training cutoff can be verified; memorization is nearly impossible
- **Negative control:** Compare model's accuracy on HumanEval vs. LiveCodeBench
  - If HumanEval >> LiveCodeBench (e.g., 75% vs. 50%), strong sign of memorization
  - If similar, suggests genuine capability

- **MBPP extended set:** Subset published later; better memorization resistance than core
- **Local verification:** Run model on unreleased problems from same distribution (e.g., LeetCode private problem set)

## 4. Leaderboard Methodology & Reproducibility

### OpenCompass (Open-source Chinese evaluation platform)

- **Benchmark coverage:** HumanEval, MBPP, LiveCodeBench, several others
- **Evaluation method:** Code execution + output matching (strict)
- **Confidence reporting:** Pass@1, Pass@10, Pass@100 (K-runs with beam search)
- **Reproducibility:** High (open-source, model weights disclosed, runs may be replicable)
- **Gotchas:** Different Python environment versions → subtle import/library version mismatches possible

### CodeArena (Meta/AI Alignment Research Center)

- **Benchmark coverage:** Custom curated; focus on code reasoning (not just execution)
- **Evaluation method:** Partial credit; reasoning clarity scored alongside correctness
- **Confidence reporting:** Detailed rubric-based scores; not just pass/fail
- **Reproducibility:** Medium (proprietary evaluation framework)
- **Gotchas:** Mixes correctness with interpretation → harder to isolate pure code capability

### Hugging Face Leaderboards

- **Benchmark coverage:** HumanEval, MBPP, sometimes custom datasets
- **Evaluation method:** Varies by leaderboard; often pass@1 (single generation, no beam)
- **Confidence reporting:** Usually pass@1 only; no confidence intervals published
- **Reproducibility:** Low (often unclear how many runs, seed handling, environment setup)
- **Gotchas:** Rankings can flip week-to-week due to statistical noise (small N runs)

### Statistical Rigor Comparison

| Leaderboard | Confidence Intervals | Multiple Runs | Reported Variance | Seed Control |
|-------------|----------------------|---------------|-------------------|--------------|
| OpenCompass | Yes (implicit in Pass@K) | Yes (via beam search) | Sometimes | Often |
| CodeArena | Rubric-based (no CI) | Single generation | No | Unclear |
| HF Leaderboards | No | Usually single | Not published | Often uncontrolled |

**Bottom line:** OpenCompass > CodeArena > HF Leaderboards for statistical trustworthiness.

## 5. Confidence Intervals & Statistical Validity

### Pass@K Metric (Standard for HumanEval/MBPP)

- **Definition:** Model generates K independent samples; solution "passes" if ANY sample solves the problem
- **Formula:** Pass@k = 1 - (# unsolved / N) × (1 - k/N) [approximation]
- **Interpretation:**
  - Pass@1 = single-attempt pass rate (what matters for real use)
  - Pass@10, Pass@100 = theoretical upper bound (not realistic for deployment)
  - Gap between Pass@1 and Pass@K = how often model has idea but needs retries

### Confidence Interval Calculation

For a benchmark with N=164 problems (HumanEval), observed pass rate p̂ = 72%:

- **95% CI (Normal approximation):** 72% ± 1.96×√(0.72×0.28/164) ≈ 72% ± 7%
- **Practical meaning:** True capability is somewhere in [65%, 79%]
- **Implication:** A 2% score difference is **not statistically significant** at N=164; need 5%+ gap to be confident

For MBPP (N=974):
- **CI:** 72% ± 2.8% (tighter due to larger N)
- **Result:** 2% differences become marginally significant

**Key insight:** HumanEval's small N (164) means leaderboard rankings are noisy. MBPP and LiveCodeBench provide better signal.

### Reproducibility & Seed Effects

- Single-seed runs (common in HF leaderboards) have ±3–5% variance just from sampling randomness
- Reported scores without confidence intervals are **not trustworthy for model comparison**
- Always ask: "How many runs?" and "What's the reported std dev?"

## 6. Task Distribution Mismatch & Real-World Prediction

### Why Benchmark Scores Don't Predict Production Performance

1. **Code length:** HumanEval: ~300 tokens; production: 100–10,000+ tokens
   - Models degrade with context (lost-in-the-middle phenomenon)
   - Benchmark tests short problem-solving; production requires long refactoring

2. **Task type skew:** Benchmarks emphasize algorithm-implementation; production emphasizes:
   - Refactoring existing code
   - Debugging and test-driven fixes
   - Library API selection
   - Documentation & type annotations
   - Integration with unfamiliar codebases

3. **Editing vs. generation:** HumanEval = write from scratch; production = modify/complete
   - Models often better at editing than greenfield generation

### Predicting Real-World Performance

**If you need to pick a model for production coding:**

1. **Collect benchmark scores on HumanEval, MBPP, LiveCodeBench**
   - If all three within ±3% → model is consistent, good signal
   - If HumanEval >> others → likely memorization; expect worse real-world performance

2. **Match task distribution to your domain**
   - Web dev → MBPP string/API tasks more predictive than HumanEval algorithms
   - Systems programming → neither benchmark is ideal; test locally
   - Data wrangling → MBPP better than HumanEval

3. **Adjust for domain shift**
   - Benchmark performance to real-world ≈ benchmark_score × (0.7–0.85)
   - E.g., model with 75% on MBPP → expect 53–64% first-try success on your codebase
   - Gap is due to: unfamiliar APIs, integration complexity, style/idiom differences

4. **Use LiveCodeBench as tie-breaker**
   - Models ranked A > B on HumanEval might reverse on LiveCodeBench
   - LiveCodeBench rank is more predictive of multi-week performance

## 7. Benchmark Bias & Limitations

### Inherent Biases

- **Language bias:** Most benchmarks are Python-first (HumanEval, MBPP, LiveCodeBench); TypeScript/Rust/Go models tested on translations
- **Test-case bias:** Edge cases in benchmark tests ≠ edge cases in your domain
- **Solution-style bias:** Benchmarks reward compact solutions; your codebase may value readability
- **Recency bias:** Benchmarks frozen at publication; new language features, libraries not tested

### Task-Specific Failures

- **String/regex:** Benchmarks over-sample (common in MBPP); models may underperform on non-string code
- **I/O & state:** Benchmarks mostly pure functions; models struggle with stateful code
- **Debugging:** No benchmarks test "fix a failing test" (models often better than scores suggest)
- **Refactoring:** No benchmarks for "improve this code" (models significantly better than on HumanEval)

**Implication:** Don't generalize a benchmark score to all code tasks. Benchmark X predicts performance on tasks *like* benchmark X.

## 8. Benchmark Selection Guide

**Use HumanEval if:**
- Quick broad capability check needed
- Comparing models on well-known benchmark
- Communicating with others familiar with HumanEval
- ⚠️ Don't use for final production decision (too saturated, high memorization risk)

**Use MBPP if:**
- More robust generalization signal needed
- Task distribution closer to production (APIs, utilities, string/data manipulation)
- Want to avoid HumanEval's saturation ceiling
- Good balance of diversity and reproducibility

**Use LiveCodeBench if:**
- Overfitting detection required (compare HumanEval to LiveCodeBench)
- Want most recent model ranking (weekly updates prevent obsolescence)
- Need statistically fresh benchmark (uncontaminated training data)
- Production model selection (best predictor of real-world performance)

**Multi-benchmark strategy (recommended):**
```
Score all three benchmarks.
If HumanEval ≥ MBPP + 5% → suspect memorization
If LiveCodeBench ≥ MBPP → model generalizes well to new tasks
If MBPP ≥ LiveCodeBench by ≤2% → model is production-ready
Decision → LiveCodeBench rank > MBPP > HumanEval
```

## 9. Common Pitfalls & How to Avoid Them

| Pitfall | Mistake | Correction |
|---------|---------|-----------|
| Single-benchmark decision | "Model A scored 75% on HumanEval, so hire it" | Use MBPP + LiveCodeBench confirmation |
| Ignoring confidence intervals | Ranking models by <2% score gap | Only trust ≥3% gaps on HumanEval, ≥2% on MBPP |
| No memorization check | Assuming high HumanEval score = capability | Compare to LiveCodeBench; check training date |
| Task distribution mismatch | "Model scored 80% on HumanEval, so it'll do great refactoring" | Benchmark score ≠ production score; apply domain adjustment |
| Ignoring leaderboard methodology | Comparing HF leaderboard to OpenCompass directly | Use same leaderboard or adjust for methodology |
| Single seed runs | Trusting leaderboard score from 1 run | Demand multiple runs; distrust single-seed results |
| Extrapolating to new tasks | "HumanEval predicts performance on Advent of Code" | Only generalize within similar task distribution |

## 10. Advanced Topics

### Estimating Performance on Unreleased Tasks

If you have 2–3 released benchmark scores for a model, you can estimate real-world performance:

1. **Fit a domain-adjustment curve:**
   - E.g., HumanEval → MBPP → LiveCodeBench → local task
   - Typical drop: 3–5% per jump in realism

2. **Account for task similarity:**
   - If your tasks are 80% "algorithm" + 20% "string manipulation", weight accordingly
   - Benchmark scores apply *directly* only to near-identical tasks

### Designing Your Own Benchmark

If public benchmarks don't fit your domain:

- Start with HumanEval + MBPP as baseline for reproducibility
- Add 20–50 problems *from your actual codebase* (anonymized)
- Ensure ≥20 samples per task type (statistics require N ≥ 20)
- Run multiple seeds (≥3); report mean ± std dev
- Compare model scores on public + private benchmarks to infer domain transfer

## References

Key papers and resources cited:

- OpenAI. "Evaluating Large Language Models Trained on Code." arXiv preprint arXiv:2107.03374 (2021). — HumanEval definition
- Google. "MBPP: A Mostly Basic Programming Problems Dataset." — Dataset and initial benchmarks
- Li et al. "LiveCodeBench: Holistic and Contamination-Free Evaluation of Large Language Models on Code." (2024). — Recency and memorization analysis
- Sobania et al. "Memorization in Large Language Models for Code." — Training data contamination evidence
- OpenCompass Team. "OpenCompass: A Comprehensive Evaluation Platform for LLMs." — Leaderboard methodology
- Anthropic, OpenAI benchmark reports. — Confidence interval reporting

---

## Quick Lookup

**Q: My model scores 73% on HumanEval. Is this good?**
- A: It's average; likely in saturation region. Check MBPP + LiveCodeBench. If those are ≤68%, suspect memorization.

**Q: Should I trust HF leaderboard rankings?**
- A: No single-seed run. Use OpenCompass or run your own if possible.

**Q: How much worse will my model perform on real code vs. HumanEval?**
- A: Expect ~15–30% drop (benchmark → production). Adjust per domain and task distribution match.

**Q: What's the minimum benchmark score to ship a model?**
- A: MBPP ≥70% + LiveCodeBench ≥65%. HumanEval is secondary.

**Q: Is my model memorizing HumanEval?**
- A: If HumanEval − LiveCodeBench ≥ 5%, likely yes. Verify by checking training date vs. benchmark publication date.
