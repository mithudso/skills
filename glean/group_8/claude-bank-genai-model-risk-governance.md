# bank-genai-model-risk-governance

**Category:** Science, Biology & Medicine
**Platform:** Claude
**Original Path:** claude/standalone/bank-genai-model-risk-governance

## Description
GenAI/LLM governance and model-risk framing for BANKS (FSI), for a MongoDB TAM advising banks on AI/LLM/vector-search. Educational, NOT legal/compliance advice. TRIGGER: model risk management for AI/LLM (SR 11-7 and its 2026 successor SR 26-2; effective challenge; model inventory); how MRM applies to non-deterministic GenAI & third-party/vendor models; NIST AI RMF 1.0 (GOVERN/MAP/MEASURE/MANAGE) and the Generative AI Profile (NIST-AI-600-1); EU AI Act (Reg 2024/1689) risk tiers, GPAI obligations, applicability timeline, credit-scoring as high-risk; validating a RAG/vector-search system. SKIP: company-specific account intel -> goldman-sachs-account-intelligence / jpmorgan-chase-account-intelligence; general FSI regulation -> fsi-banking-regulatory-context; vendor/third-party-risk -> enterprise-vendor-management-and-tprm; prompt-injection -> rag-prompt-injection-defense; security controls -> atlas-vector-search-pii-isolation / mongodb-operations-expert.

---

# Bank GenAI Model-Risk Governance

What a MongoDB TAM should understand about **AI/LLM governance and model risk in
banks**, to speak credibly when a bank customer wants to put an Atlas Vector
Search + LLM system into production. **Educational only — not legal or
compliance advice; the bank's model-risk and compliance functions own the
determinations.** **Core thesis: a bank cannot ship a GenAI/RAG system without a
governance wrapper (inventory, validation, monitoring, documentation); the TAM's
job is to make the MongoDB design *validatable and auditable*, not to interpret
the law.**

## 1. US model risk: SR 11-7 and its 2026 successor ⚠️

- **SR 11-7** (Fed/OCC *Supervisory Guidance on Model Risk Management*, 2011) is the
  long-standing US baseline. Core ideas every TAM should recognize:
  - **Model risk** = risk of adverse outcomes from model errors or misuse.
  - **Effective challenge** = critical, independent review by parties with competence,
    influence, and incentive.
  - **Three pillars**: (1) development/implementation/use, (2) **validation**, (3)
    **governance/policies/controls**.
  - **Model inventory** (every model registered) and **ongoing monitoring**.
- ⚠️ **MAJOR, RECENT, VERIFY-FIRST CLAIM:** research indicates SR 11-7 was **superseded
  by SR 26-2 on 2026-04-17**, and that the successor guidance treats **GenAI/agentic AI
  as out of its direct scope** (handled under other guidance), among other updates.
  *This is a high-impact, very recent change sourced to federalreserve.gov (SR2602) and
  secondary coverage (ABA Banking Journal, Moody's, Baker Tilly).* **Do not state this to
  a customer without pulling the primary Fed/OCC letter and confirming current scope** —
  treat SR 11-7's principles as the durable baseline and SR 26-2 as the item to verify.

**The hard part for GenAI either way:** LLMs are **non-deterministic**, often
**third-party/vendor** models (the bank didn't train them), and behavior shifts with
prompts/data. Traditional statistical validation (backtesting a fixed function) doesn't
map cleanly — banks substitute eval suites, red-teaming, guardrail testing, human review,
and tight monitoring. A RAG system's "model" includes the retriever, the embedding model,
and the prompt — all in scope.

## 2. NIST AI Risk Management Framework (voluntary, widely adopted) ✅

- **AI RMF 1.0** (NIST, Jan 2023), four functions:
  - **GOVERN** — culture, policies, accountability (cross-cutting).
  - **MAP** — context, intended use, and risks.
  - **MEASURE** — analyze/benchmark/track risks (evals, metrics, red-teaming).
  - **MANAGE** — prioritize, treat, and monitor risks.
- **Generative AI Profile — NIST-AI-600-1** (July 2024): a companion enumerating
  GenAI-specific risks (e.g. **confabulation/hallucination**, data privacy, information
  security, harmful bias, dangerous content, IP) with suggested actions mapped to the
  four functions. This is the most practical checklist to map a RAG design against.
- ✅ Voluntary but the de-facto common language; US banks frequently align their AI
  governance to it even where not strictly required.

## 3. EU AI Act — Regulation (EU) 2024/1689 ⚠️ (timeline-sensitive)

- **Risk tiers**: *unacceptable* (banned), *high-risk* (strict obligations), *limited*
  (transparency), *minimal*. For banks, **credit scoring / creditworthiness assessment is
  listed high-risk (Annex III)** — triggering risk management, data governance, logging,
  human oversight, accuracy/robustness, and conformity obligations.
- **GPAI (general-purpose AI) models**: providers face transparency/documentation duties;
  **systemic-risk GPAI** (very capable models) face added obligations. A bank *using* a
  third-party LLM is typically a **deployer**, with its own (lighter) duties — but
  high-risk use cases pull in the full obligation set.
- **Timeline** ⚠️: phased applicability **2025–2027** (prohibitions and AI-literacy first;
  GPAI obligations and high-risk requirements later). **Confirm the exact date each
  obligation bites against EUR-Lex / the EU AI Office before advising** — dates have
  already shifted in public discussion.
- **Extraterritorial + heavy penalties** (up to the higher of a large fixed sum or a % of
  global turnover). A global bank with EU customers is in scope.

## 4. How they intersect for a global bank 🔶

A US/EU/global bank must satisfy **all** applicable regimes at once: SR 11-7/SR 26-2-style
MRM (US prudential), NIST AI RMF (common control language), and the EU AI Act (binding for
EU use). Common denominators a TAM can lean on:

- **Model/AI inventory** — the system must be registered and documented.
- **Documentation & lineage** — what model, what data, what retrieval, what prompt, what
  evals; reproducible.
- **Validation/eval + ongoing monitoring** — including drift, guardrail efficacy, and the
  RAG retrieval quality.
- **Human oversight** for consequential decisions.

## 5. What a MongoDB TAM should actually do

- **Make the design validatable & auditable**: stable index/version metadata, audit
  logging (who retrieved what), reproducible retrieval config, and clear data lineage from
  source → embedding → index → context.
- **Map the MongoDB controls to the frameworks**: encryption/CMEK, RBAC, private
  networking, entitlement-filtered retrieval (`atlas-vector-search-pii-isolation`), and the
  injection/leakage threat models (`rag-prompt-injection-defense`,
  `embedding-inversion-threat-model`) — these are concrete evidence for MEASURE/MANAGE and
  for high-risk obligations.
- **Stay in your lane**: surface the governance hooks and supply evidence; the bank's MRM,
  legal, and compliance teams make the determinations. Route firm regulatory questions to
  `fsi-banking-regulatory-context` and procurement/TPRM to
  `enterprise-vendor-management-and-tprm`.

## Sources

- Federal Reserve / OCC, **SR 11-7** "Guidance on Model Risk Management" (2011),
  federalreserve.gov; **SR 26-2** (2026-04-17) ⚠️ verify primary letter.
- NIST, **AI RMF 1.0** (2023) and **Generative AI Profile, NIST-AI-600-1** (2024-07),
  nist.gov.
- **EU AI Act**, Regulation (EU) 2024/1689, EUR-Lex; EU AI Office guidance.

⚠️ This domain is fast-moving and jurisdiction-specific. Every dated/scoped claim here —
especially SR 26-2 and EU AI Act applicability dates — must be re-verified against the
primary source before any customer-facing use, and nothing here is legal advice.