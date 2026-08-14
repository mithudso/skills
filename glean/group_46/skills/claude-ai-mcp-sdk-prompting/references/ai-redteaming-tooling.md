<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.** Formerly the standalone `ai-redteaming-tooling` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: ai-redteaming-tooling
title: AI Red-Teaming & Security-Testing Tooling
description: >
  The OFFENSIVE security-testing discipline for LLM/agent apps — systematically
  generating adversarial inputs to find jailbreaks, prompt injection,
  data/system-prompt exfiltration, and unsafe tool use, then feeding findings
  into fixes and regression tests. Distinct from runtime guardrails/defense
  (input/output filters, dual-LLM/CaMeL, the lethal trifecta) which live in
  agent-reliability-and-guardrails. Covers the tooling landscape (Garak scanner,
  Microsoft PyRIT framework, promptfoo red-team, Giskard, Meta Purple Llama /
  CyberSecEval); the attack taxonomy (direct vs indirect prompt injection,
  jailbreak families — DAN/roleplay, payload splitting, encoding, many-shot,
  Crescendo; exfiltration); automated attack-generation algorithms (GCG, PAIR,
  TAP tree-of-attacks, red-teamer LLMs); benchmarks & datasets (AdvBench,
  HarmBench, JailbreakBench, AgentDojo for tool-use injection, OWASP LLM Top 10
  as a test checklist); commercial continuous red-teaming platforms
  (Lakera/Gandalf→Cisco, Robust Intelligence→Cisco AI Defense, Mindgard,
  HiddenLayer); and process/governance (OWASP GenAI Red Teaming Guide 4 phases,
  MITRE ATLAS, NIST AI RMF, red-team-in-CI). TRIGGER: red-teaming or
  security-testing an LLM/agent app; "test for jailbreaks / prompt injection";
  Garak / PyRIT / promptfoo red-team / Giskard / CyberSecEval; GCG / PAIR / TAP /
  Crescendo / many-shot; AgentDojo / JailbreakBench / HarmBench; OWASP LLM Top 10
  testing; continuous red-teaming / red-team-in-CI. SKIP: DEFENSIVE guardrails —
  input/output filters, dual-LLM/CaMeL, the lethal trifecta mitigation (use
  agent-reliability-and-guardrails); source-code security review / SAST (use
  security-reviewer); compliance policy audit (use security-compliance-auditor).
origin: local
category: developer
version: "1.0"
updated: "2026-05-31"
tags:
  - ai-security
  - red-teaming
  - prompt-injection
  - jailbreak
  - llm
  - agent
  - testing
  - owasp-llm
whenToUse:
  - "red-teaming / adversarially testing an LLM or agent application"
  - "choosing a red-team tool (Garak scanner vs PyRIT framework vs promptfoo vs Giskard)"
  - "understanding jailbreak families and direct vs indirect prompt injection (offensive)"
  - "using automated attack generation (GCG, PAIR, TAP, red-teamer LLMs)"
  - "running red-team benchmarks (AdvBench, HarmBench, JailbreakBench, AgentDojo)"
  - "standing up continuous red-teaming / red-team-in-CI"
  - "mapping testing to OWASP LLM Top 10 / MITRE ATLAS / NIST AI RMF"
whenNotToUse:
  - "defensive guardrails (filters, dual-LLM/CaMeL, lethal-trifecta mitigation) — use agent-reliability-and-guardrails"
  - "source-code security review / SAST — use security-reviewer"
  - "compliance/policy audit — use security-compliance-auditor"
related_skills:
  - agent-reliability-and-guardrails
  - security-reviewer
  - security-compliance-auditor
---

# AI Red-Teaming & Security-Testing Tooling

The **offensive** half of LLM/agent application security: generate adversarial
inputs, find where the system fails in ways you don't want, and convert each
finding into a regression test. The defensive counterpart — input/output
guardrails, dual-LLM/CaMeL, the lethal-trifecta *mitigation* — lives in
`agent-reliability-and-guardrails`; this reference treats those as *testing
targets* and points there for fixes.

Two consensus lessons frame the field (Microsoft's "red-teaming 100 GenAI
products"): **automation augments but does not replace human red-teamers**, and
**red-teaming must be continuous, not one-off** — and notably, *safety
benchmarks cannot substitute for red-teaming* (benchmarks measure known
behaviors; red-teaming finds novel failures).

## Tooling landscape — scanners vs frameworks

| Tool | Type | What it does | Owner |
|---|---|---|---|
| **Garak** | Scanner ("Nessus for LLMs") | Curated probe battery → vuln report; 5-part plugin arch (Probes / Detectors / Generators / Harnesses / Evaluators), 20+ probe modules | NVIDIA, Apache-2.0 |
| **PyRIT** | Framework ("Metasploit for LLMs") | Composable, *adaptive* multi-turn attack orchestration: Targets / Converters / Scorers / Orchestrators; now in Azure AI Foundry Red Teaming Agent | Microsoft, MIT |
| **promptfoo** | Scanner + eval, CI-native | 50+ vuln types, OWASP/NIST/ATLAS presets, **fails the build** on regression; runs locally | promptfoo, OSS |
| **Giskard** | Scanner | Security (injection, leakage) + quality (hallucination, bias, sycophancy) scan for LLM/RAG | Giskard, Apache-2.0 |
| **Purple Llama** | Models + benchmark | Llama Guard / Prompt Guard (defensive) + **CyberSecEval** (offensive benchmark: insecure-code + cyberattack compliance; v3 adds visual injection) | Meta |

Scanner = run a fixed battery, get a report (Garak, Giskard, promptfoo). Framework
= script your own adaptive attacks (PyRIT). promptfoo straddles eval + red-team
and is the reference "red-team-in-CI" implementation.

## Attack taxonomy (offensive)

- **Direct vs indirect prompt injection.** *Direct* = user input overrides/reveals
  the system prompt (often conflated with jailbreaking). **Indirect** = malicious
  instructions smuggled via untrusted content the model ingests (RAG docs, web
  pages, files, tool outputs) — the **attacker is a third party**, not the user.
  Indirect injection (Greshake et al., arXiv:2302.12173) is the highest-impact
  *agentic* surface and the core of the lethal trifecta.
- **Jailbreak families.** DAN ("Do Anything Now") roleplay/dual-persona,
  persona/"admin" escalation, **payload splitting**, **encoding/obfuscation**
  (Base64, multilingual, hiding in code comments), multi-turn gradual escalation.
- **Many-shot jailbreaking** (Anthropic, 2024) — hundreds-to-thousands of
  in-context examples of bad behavior, then the harmful ask; effectiveness
  follows a **power law** with shot count, newly feasible thanks to long context.
- **Crescendo** (multi-turn) — benign-looking turns exploit the model's tendency
  to follow its own recent output; up to ~100% ASR on many tasks.
- **Exfiltration** — a successful injection in a tool-enabled app exfiltrates data
  via tools (writing to a public repo, DNS, markdown-image rendering). Maps to
  OWASP **LLM02** (Sensitive Info Disclosure) and **LLM07** (System Prompt
  Leakage).
- **The lethal trifecta** (Willison, 2025) — private-data access + untrusted
  content + external communication = structurally exploitable. As a *test target*
  it's indirect-injection + tool-poisoning cases; the *defense* is in
  `agent-reliability-and-guardrails`.

## Automated attack generation

- **GCG (Greedy Coordinate Gradient)** — token-level **white-box**; optimizes an
  adversarial suffix to force "Sure, here is…". High success, very query-heavy;
  suffixes can transfer across models.
- **PAIR (Prompt Automatic Iterative Refinement)** — an **attacker LLM** refines
  semantically coherent prompts; **black-box**, query-efficient.
- **TAP (Tree of Attacks with Pruning)** — black-box; attacker LLM grows and
  prunes a tree of candidate prompts; vs GPT-4o, jailbreaks 16% more prompts than
  PAIR with 60% fewer queries. Productized by Robust Intelligence/Cisco for
  "algorithmic red-teaming."
- **Red-teamer LLMs / fuzzing** — the modern pattern: an attacker LLM generates,
  mutates, and scores attacks in a loop (PyRIT orchestrators, TAP, PAIR,
  Crescendomation). Hybrid GCG+PAIR adds up to +33pp ASR. *Note (Microsoft): in
  real products, prompt engineering often beats gradient attacks — don't fixate
  on white-box.*

## Benchmarks & datasets

- **AdvBench** (2023) — original harmful-behaviors/strings set (shipped with GCG).
- **HarmBench** (CAIS, 2024) — standardized automated-red-team + robust-refusal
  framework, broad coverage incl. copyright/multimodal.
- **JailbreakBench** (NeurIPS 2024) — 100 harmful + 100 benign behaviors, a
  **validated classifier**, artifact repo, public leaderboard.
- **AgentDojo** (2024) — the key **tool-use / agent injection** benchmark: 97
  realistic tasks, 629 security cases (banking, Slack, travel, workspace);
  metrics = benign utility, utility-under-attack, ASR.
- **OWASP LLM Top 10 (2025) as a test plan** — LLM01 Prompt Injection · LLM02
  Sensitive Info Disclosure · LLM03 Supply Chain · LLM04 Data/Model Poisoning ·
  LLM05 Improper Output Handling · LLM06 Excessive Agency · LLM07 System Prompt
  Leakage · LLM08 Vector/Embedding Weaknesses · LLM09 Misinformation · LLM10
  Unbounded Consumption.

## Commercial continuous red-teaming

- **Lakera** — *Gandalf* (gamified injection challenge; 1M+ players, 80M+
  adversarial prompts feeding threat intel) + *Lakera Red* (pre-prod attack
  simulation). Acquired by Cisco (2025).
- **Robust Intelligence → Cisco AI Defense** — "first algorithmic red-teaming"
  using TAP; now Cisco "AI Validation."
- **Mindgard** — automated **DAST-AI** for LLMs/agents/multimodal; continuous.
- **HiddenLayer** — ML model protection + threat intel.

## Process & governance

- **OWASP GenAI Red Teaming Guide (2025) — four phases:** (1) model evaluation,
  (2) implementation testing (guardrails in place), (3) infrastructure/system
  assessment, (4) runtime behavior analysis.
- **Framework mapping (complementary):** OWASP LLM Top 10 = dev/secure-coding;
  **MITRE ATLAS** (v5.1.0: 16 tactics / 84 techniques / 32 mitigations) =
  ops/threat-modeling; **NIST AI RMF** (Govern/Map/Measure/Manage) = governance.
- **Red-team in CI + reporting:** run the suite in CI/CD and fail the build on
  regression (promptfoo); document each case with goal, prompt, model version,
  and observed behavior as a reproducible test.

## Anti-patterns

- **One-off red-teaming** — models/prompts/threats drift; testing must be
  continuous.
- **Testing only direct injection** — ignoring indirect injection / tool
  poisoning, the highest-impact agentic surface.
- **No regression after fixes** — not re-running the attack that found the bug.
- **Benchmarks as a substitute** for red-teaming (they measure known behaviors).
- **Over-relying on gradient attacks** when prompt engineering is cheaper and
  often more effective.
- **Automation-only** — scanners miss novel attacks human creativity finds.

## Cross-references

- **Defensive guardrails** (input/output filters, dual-LLM/CaMeL, lethal-trifecta
  mitigation, Llama Guard/Prompt Guard as deployed defenses) →
  `agent-reliability-and-guardrails`.
- **Source-code security review / SAST / insecure-code generation review** →
  `security-reviewer`.
- **Compliance & policy audit** → `security-compliance-auditor`.
- **Tool-use injection in context** → `agent-harness-construction`,
  `advanced-rag-patterns` (RAG/tool poisoning surface).

## References

NVIDIA Garak (GitHub + README); Microsoft PyRIT (Security Blog 2024; Azure AI
Foundry); promptfoo red-team docs + CI/CD; Giskard LLM scan; Meta Purple Llama /
CyberSecEval (arXiv:2312.04724). Attacks: Greshake indirect injection
(2302.12173); TAP (2312.02119); Crescendo (2404.01833); Anthropic Many-Shot.
Benchmarks: JailbreakBench (2404.01318), AgentDojo (2406.13352). Process: OWASP
LLM Top 10 2025 + GenAI Red Teaming Guide; MITRE ATLAS; NIST AI RMF; "Lessons
from Red Teaming 100 GenAI Products" (2501.07238). Platforms: Lakera/Gandalf,
Cisco AI Defense (Robust Intelligence), Mindgard. *(46 sources, 2024–2026; full
URLs in the source research report.)*
