<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** New `/dr`-researched reference (concept-family-explorer run, 2026-06-15).
> Sibling topics in this family are reference files under the hubs — **not** standalone skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling skill; load that topic's `references/<name>.md` from the owning hub.

---

---
name: automated-program-repair
title: Automated Program Repair & Code Auto-Remediation
description: >
  Code-level automated program repair (APR) and code auto-remediation as a field: the four-family taxonomy (search/heuristic GenProg, constraint/semantics SemFix-Angelix-Nopol, template/pattern TBar, learning- and LLM-based patch generation), LLM/agentic repair paradigms (fine-tune/prompt/procedural/agentic, RepairAgent), the Agentless localize->repair->validate pipeline as a mechanism reference, patch overfitting and the test-suite-as-weak-oracle problem, spectrum-based fault localization (Tarantula/Ochiai/DStar/Barinel), and industrial self-healing CI (Getafix, SapFix/Sapienz/Infer, Dependabot/Renovate auto-remediation). TRIGGER: "automated program repair", "generate-and-validate / GenProg", "semantics vs template vs learning-based repair", "how does RepairAgent / Agentless work", "my auto-generated patch passes tests but is wrong / patch overfitting", "plausible vs correct patch", "spectrum-based fault localization / Ochiai", "SWE-bench repair evaluation hygiene", "how Getafix/SapFix fix bugs", "Dependabot/Renovate auto-fix PRs". SKIP: designing the agent's plan->edit->test->repair self-repair LOOP, Agentless-as-an-agent-architecture, or SBFL inside an agent -> coding-agents; TRAINING a repair agent (SWE-RL, RLVR reward) -> agentic-rl; debugging/root-cause method, code review, diagnostic registry -> software-engineering-patterns; runtime/infrastructure self-healing (MAPE-K, Kubernetes, AIOps) -> self-healing-systems-autonomic-computing; skill authoring/quality -> claude-code-skills / skill-optimizer.
category: developer
---

# Automated Program Repair & Code Auto-Remediation

## Overview

Automated Program Repair (APR) is the field that takes a buggy program plus a *correctness oracle* (most often a test suite) and produces a source-code patch that makes the program pass the oracle. This skill is the **conceptual + mechanism reference** for code-level repair: how patches are generated, how faults are localized, why "passing tests" is not the same as "correct", and how the field industrialized into self-healing CI and dependency-update bots. It is **not** an agent-loop design skill — when the question is how to build the plan->edit->test->repair loop of a coding agent, route to `coding-agents`; when it is how to *train* a repair model, route to `agentic-rl`.

The authoritative spine is the ACM Computing Surveys 2024 APR taxonomy (Huang et al., *Evolving Paradigms in Automated Program Repair*, dl.acm.org/doi/10.1145/3696450) and the 2025 LLM-APR survey (arxiv.org/html/2506.23749v1). The field has pivoted hard toward LLM-based repair; the central conceptual axis, as the 2024 survey characterizes it, is *optimizing template selection (probabilistic models)* vs *learning high-level repair knowledge (neural/LLM)*.

## When to Use This Skill

| Task | Examples |
|---|---|
| Place a repair technique in the taxonomy | search vs constraint vs template vs learning-based |
| Reason about LLM/agentic repair | fine-tune vs prompt vs procedural pipeline vs agentic (RepairAgent) |
| Explain the localize->repair->validate pipeline | Agentless three-phase mechanism, reranking, reproduction tests |
| Diagnose a bad auto-patch | patch overfitting, plausible-vs-correct, weak-oracle |
| Pick/explain a fault-localization metric | Tarantula/Ochiai/Jaccard/DStar/Barinel suspiciousness |
| Map industrial repair systems | Getafix, SapFix/Sapienz/Infer, Dependabot, Renovate |
| Set up repair evaluation hygiene | dataset leakage, undertested functions, SWE-bench caveats |

## Core Concepts

1. **Generate-and-validate (G&V) repair & GenProg.** The original paradigm: mutate/recombine candidate edits (genetic programming), keep candidates that pass the test suite. The test suite *is* the oracle; this is the root of every downstream overfitting problem (ACM CSUR 2024 taxonomy).

2. **The four-family APR taxonomy.** (a) **Search/heuristic-based** — GenProg, TrpAutoRepair search an edit space; (b) **Constraint/semantics-based** — SemFix, Angelix, Nopol use symbolic execution + component-based program synthesis to *synthesize* a repair satisfying a constraint (angelic forests); (c) **Template/pattern-based** — TBar applies curated fix templates ranked by a probabilistic model; (d) **Learning-based** — neural/LLM models that learn high-level repair knowledge. The 2024 survey frames "optimizing template selection" vs "learning repair knowledge" as the core axis.

3. **LLM-APR design paradigms (2022-2025).** The 2025 survey of 63 systems splits LLM repair into four paradigms: **fine-tuning**, **prompting**, **procedural pipelines** (fixed localize->repair->validate stages), and **fully agentic frameworks** (the LLM freely interleaves tools). Named systems: SWE-agent, AutoCodeRover, RepairAgent, OpenHands.

4. **RepairAgent — the agentic-repair baseline (ICSE 2025).** First autonomous LLM-agent APR system: an LLM + a set of repair tools (read-file, apply-patch, run-tests) and — per the paper's architecture — a middleware orchestrator and a finite-state machine guiding tool invocation, with a *dynamically updated prompt*. It freely interleaves information-gathering, repair-ingredient collection, and validation. On Defects4J it autonomously repaired 164 bugs (39 not fixed by prior work) at ~270k tokens/bug (~14 cents on GPT-3.5). *NOTE: the agent-loop mechanics belong in `coding-agents`; here it is a reference point for cost/effectiveness.*

5. **Agentless localize->repair->validate (FSE 2025).** A deliberately *non-agentic* three-phase pipeline that often beats open-ended agents on cost/accuracy: (a) **hierarchical localization** (file -> class/function -> edit location, LLM + embedding IR), (b) **repair** (generate multiple diff candidates), (c) **validation** using existing regression tests + LLM-generated reproduction tests, then majority-voting/reranking to select the submitted patch. 32% on SWE-bench Lite at ~$0.70/issue; >50% on SWE-bench Verified with Claude 3.5 Sonnet. This is the canonical contrast to agents; *Agentless-as-an-agent-architecture routes to `coding-agents`.*

6. **Patch overfitting & the weak-oracle problem.** The foundational FSE 2015 result (Smith et al., *Is the Cure Worse Than the Disease?*): patches that pass the repair test suite are, on held-out tests, **as likely to break as to fix** undertested functionality. Patch quality is proportional to repair-suite coverage; low-coverage suites (the common case) produce more overfitting patches, and patch minimization does **not** fix it. This is the *plausible vs correct* distinction and the reason "passes tests" is not "is fixed".

7. **Overfitting detection.** A 2024+ subfield: dynamic approaches (Opad fuzzing + the O-measure), static/learning-based classifiers, and oracle-enhancement (auto-generating stronger tests). Frontier: **LLM-as-judge patch correctness** (LLM4PatchCorrect, execution-free critics, PatchDiff) to filter plausible-but-wrong patches without running tests.

8. **Spectrum-based fault localization (SBFL) — the front half.** Rank program entities by *suspiciousness* computed from pass/fail test spectra over the tuple (cef, cnf, cep, cnp) using similarity coefficients: **Tarantula, Ochiai, Jaccard, DStar, Barinel, Op2, Kulczynski**. Empirically **Ochiai/DStar > Tarantula** (Abreu et al. on Ochiai/Jaccard; Wong et al. on DStar). SBFL feeds patch generation; modern LLM hierarchical localization (Agentless) is its learning-based successor. *SBFL-inside-an-agent routes to `coding-agents`.*

9. **Industrial self-healing CI.** **Getafix** (Meta) learns recurring fix patterns from past human commits and produces human-like fixes in seconds for Infer static-analysis findings (~42% developer acceptance). **SapFix/Sapienz/Infer** is the end-to-end industrial loop: test+static-analysis -> localize -> template/learned-pattern fix -> validate (developer + Sapienz tests) -> propose to a human reviewer (human-in-the-loop, *not* auto-deploy). The pattern: detect -> localize -> patch -> validate -> propose.

10. **Autofix / dependency-bot auto-remediation.** **Dependabot security updates** raise a PR upgrading a vulnerable dependency to the minimum patched version without breaking the dependency graph, linked to the alert. **Renovate** (Mend) opens update PRs across ecosystems and surfaces merge-confidence signals (age, adoption, pass rates). These are *oracle-driven by version metadata*, not tests — the supply-chain auto-remediation cousin of APR.

11. **Repository-scale defects & evaluation hygiene.** Persistent open problems from both surveys: verifying *semantic* correctness beyond test suites, repository-scale (not single-method) defects, weak oracles / data leakage in SWE-bench-style evaluation (ground-truth leakage, misleading issue descriptions), and LLM cost/latency for agentic repair. Watch for unfair experiment comparisons and dataset overlap.

## References Outline

A full skill would carry these sub-files under `references/`:

- `apr-taxonomy.md` — the four families with canonical systems (GenProg, SemFix/Angelix/Nopol, TBar, learning-based), the optimize-templates vs learn-knowledge axis.
- `llm-and-agentic-repair.md` — the four LLM-APR paradigms; RepairAgent FSM+tools+middleware; SWE-agent/AutoCodeRover/OpenHands as references; cost/token baselines; explicit hand-off note to `coding-agents` for loop design.
- `agentless-pipeline.md` — hierarchical localization, multi-candidate diff repair, reproduction-test + regression validation, reranking/majority-voting; the agent-vs-pipeline cost/accuracy debate.
- `fault-localization-sbfl.md` — the (cef,cnf,cep,cnp) model, the coefficient catalog and a comparison table, Ochiai/DStar superiority, LLM localization as successor.
- `patch-overfitting.md` — plausible-vs-correct, coverage-correlated quality, Opad/O-measure, LLM-as-judge correctness, evaluation hygiene/leakage.
- `industrial-and-autofix.md` — Getafix learned patterns, SapFix/Sapienz/Infer loop, Dependabot security/version PRs, Renovate merge-confidence gating.

## Key Landmarks & Citations

- Huang et al., *Evolving Paradigms in Automated Program Repair: Taxonomy, Challenges, and Opportunities*, ACM Computing Surveys 2024 — dl.acm.org/doi/10.1145/3696450 (authoritative taxonomy).
- *A Survey of LLM-based Automated Program Repair*, arXiv 2025 — arxiv.org/html/2506.23749v1 (63 systems, four paradigms).
- Bouzenia et al., *RepairAgent*, ICSE 2025 — first autonomous LLM-agent APR; Defects4J 164 bugs.
- Xia & Zhang, *Agentless*, FSE 2025 — arxiv.org/pdf/2407.01489 (localize->repair->validate, SWE-bench numbers).
- Smith et al., *Is the Cure Worse Than the Disease? Overfitting in Automated Program Repair*, FSE 2015 — the weak-oracle/overfitting foundation.
- Abreu, Zoeteweij & van Gemund, *On the Accuracy of Spectrum-based Fault Localization* (Ochiai/Jaccard); Wong et al., *The DStar Method for Effective Software Fault Localization*, IEEE Transactions on Reliability 2014 — the SBFL coefficient catalog and Ochiai/DStar effectiveness.
- Meta Engineering, *Getafix* (2018) and SapFix/Sapienz — industrial self-healing CI.
- GitHub Docs, *Dependabot security updates*; Renovate (Mend) — supply-chain auto-remediation.
