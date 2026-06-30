---
name: fsi-banking-regulatory-context
description: >-
  Financial-services-industry (FSI) fluency for a B2B technology SELLER (e.g. a
  MongoDB TAM): how banking segments differ as tech buyers, the regulation
  shaping their buying, and WHY banks buy that way. Educational, NOT legal
  advice. TRIGGER: selling/positioning tech into banks, capital markets,
  payments, wealth/asset mgmt, or insurance; a segment's tech needs (core
  banking, low-latency trading, ISO 20022); pointer-level "what is Basel III,
  Dodd-Frank, SOX, GLBA, FFIEC, BCBS 239, PCI-DSS, SR 11-7, SEC 17a-4, FINRA, CAT
  for a bank"; why FSI buyers are risk-averse and want auditability, data
  residency, segregation-of-duties; FSI cloud adoption, proof burden (SOC
  2/SIG/CAIQ), build-vs-buy. SKIP: CONSUMER finance/banking -> consumer-finance,
  personal-banking; MongoDB compliance config -> mongodb-compliance;
  third-party-risk/DORA -> enterprise-vendor-management-and-tprm; GenAI/AI
  model-risk governance (SR 11-7/SR 26-2, NIST AI RMF, EU AI Act for an LLM
  deployment) -> bank-genai-model-risk-governance; buying-committee
  psychology -> applied-psychology; trading/markets -> trading-and-investing.
version: 1.1.0
updated: 2026-06-18
category: custom
whenToUse:
  - How do I sell or position technology into a bank, and why are banks such conservative/slow buyers?
  - What regulations shape how a capital-markets firm, payments processor, or bank buys technology?
  - What is Basel III / Dodd-Frank / SOX / GLBA / FFIEC / BCBS 239 / PCI-DSS / SEC 17a-4 at a pointer level for a banking tech seller?
  - Why do FSI buyers demand auditability, data residency, and segregation-of-duties from vendors?
  - How does a retail bank's tech need differ from a payments or capital-markets firm's?
  - What's the proof burden (SOC 2, SIG/CAIQ, references) and cloud-adoption posture when selling into a regulated financial institution?
  - When do banks build vs buy, and how does FFIEC examination culture shape vendor expectations?
keywords:
  - financial services industry
  - FSI vertical
  - banking technology buyer
  - core banking systems
  - capital markets technology
  - payments ISO 20022
  - prudential regulators OCC Federal Reserve FDIC
  - Basel III endgame
  - Dodd-Frank SOX GLBA FFIEC
  - BCBS 239 risk data aggregation
  - PCI-DSS
  - SEC 17a-4 FINRA CAT WORM
  - selling into regulated industry
  - bank cloud adoption
  - build vs buy banks
  - data residency sovereignty
  - SOC 2 SIG CAIQ proof burden
tags:
  - fsi
  - financial-services
  - banking
  - regulatory
  - sales-enablement
  - enterprise-gtm
  - vertical
  - capital-markets
  - payments
  - compliance-context
---

# FSI Banking Industry & Regulatory Context for Technology Sellers

> **What this is.** Industry + regulatory **fluency for a B2B technology seller** (a MongoDB Technical Account Manager, solutions architect, or pre-sales engineer) approaching the **financial-services industry (FSI)** as a technology *buyer*. It answers three questions: **how the segments differ** as tech buyers, **the regulatory map** that shapes their buying, and **why banks buy technology the way they do** (conservatively, with heavy proof and long cycles).
>
> **What this is NOT.** This is **educational orientation, not legal or compliance advice.** It gives no compliance instructions and makes no representation about any specific institution's obligations. Regulatory facts are pointer-level ("what it is, who enforces it, why it matters to technology"). **Fast-moving claims are dated "as of 2026" and must be re-verified** before you rely on them with a customer — the 2025-2026 US regulatory landscape is unusually fluid (Basel re-proposal, model-risk guidance rewrite, sovereign-cloud, deregulation). When in doubt, route the customer to their own legal/compliance/risk function.

## How to use this skill (routing)

This is a hub. Load the reference that matches your need:

| If you need… | Read |
|---|---|
| The **segments** of FSI and how each one's technology/data needs differ (retail, commercial, capital markets, payments, wealth/asset mgmt, insurance) + the core-banking vendor landscape | `references/fsi-segments-and-tech-needs.md` |
| The **prudential regulators** (OCC, Fed, FDIC, NCUA, CFPB, FFIEC) and the **structural laws** (Basel III/IV, Dodd-Frank, SOX, GLBA) + the examination/MRA culture | `references/prudential-regulators-and-structural-laws.md` |
| The **data-centric & capital-markets regulations** that shape data architecture (BCBS 239, PCI-DSS, SR 11-7→SR 26-2 pointer, GDPR/CCPA-in-FSI, SEC 17a-4 / FINRA 4511 / Reg SCI / CAT) | `references/data-and-capital-markets-regulation.md` |
| **Why** regulation makes banks buy the way they do (risk-aversion, auditability, data residency, segregation-of-duties, long cycles, change/resilience) | `references/how-regulation-drives-buying-behavior.md` |
| The **FSI cloud-adoption posture** and the **"selling into a regulated industry"** motion (proof burden, references, security expectations, build-vs-buy) | `references/cloud-adoption-and-selling-motion.md` |

## The 60-second orientation (BLUF)

1. **FSI is not one buyer.** Retail banking, commercial banking, capital markets, payments, and wealth/asset management have *materially different* data shapes, latency tolerances, and buying logic. The biggest technical fault line is the **core banking system** (legacy mainframe vs. cloud-native) and how much else orbits it. Capital markets (microsecond latency, petabyte time-series) is the outlier that diverges most from the rest. [seg]

2. **The regulator follows the charter.** A bank's primary federal supervisor (OCC / Federal Reserve / FDIC / NCUA) depends on its charter; they examine using common **FFIEC** standards and the **FFIEC IT Examination Handbook**. The same product can face different supervisory expectations at two different banks. [reg]

3. **Examination culture is the engine.** Banks live under continuous supervisory exams; findings escalate through **MRAs/MRIAs** (Matters Requiring Attention / Immediate Attention). This is *the* reason auditability, evidence, and demonstrable controls are non-negotiable in technology decisions. [reg]

4. **A handful of regulations are squarely about data infrastructure.** **BCBS 239** (risk-data aggregation, lineage, single-source-of-truth) is the most database-relevant banking standard — and banks still struggle with it after a decade. **PCI-DSS** governs how card data is stored and segmented (tokenize to descope). **SOX 404** drives auditability, change-management, and segregation-of-duties in IT. **SEC 17a-4 / FINRA 4511 / CAT** drive immutable/WORM-or-audit-trail storage, long retention, and surveillance. [data]

5. **Banks buy conservatively because the payoff is asymmetric.** A botched change → a regulatory finding, consent order, or outage (large, visible, career-ending); moving slowly → diffuse, deferred cost. The honest 2026 model is **"banks move fast where they can *evidence* control, slowly where they can't"** — not "banks are dinosaurs." [buy]

6. **Selling in is a proof-and-trust motion, not a feature demo.** Expect 9-18 month cycles, 5-10 stakeholders with veto power, security questionnaires (SIG/CAIQ), certifications (SOC 2 Type II, ISO 27001, PCI), pen-test evidence, right-to-audit and exit clauses, and a disproportionate weight on **in-industry references**. SOC 2 is the floor, not the finish line. [sell]

## Cross-references (do not duplicate these here)

- **MongoDB-specific compliance configuration** (Atlas certifications, FedRAMP/AtlasGov, HIPAA BAA, PCI scoping in Atlas, BYOK, Queryable Encryption, audit logging config) → **`mongodb-compliance`** (also reachable via the `mongodb-operations-expert` hub, which owns encryption/CSFLE/Queryable Encryption). This skill stays vendor-neutral and industry-level.
- **Third-party risk management (TPRM) and DORA *mechanics*** (vendor-governance lifecycle, vendor risk tiering, the security-questionnaire gauntlet from the buyer side, DORA RTS / critical-third-party regimes as a process, procurement/contract-instrument stack) → **`enterprise-vendor-management-and-tprm`** (the buyer-side sibling). This skill only points to "operational resilience / third-party risk" as the *adjacent regime that explains why* assurance is heavy.
- **Enterprise POC/POV technical-validation *mechanics*** (POC vs POV vs pilot vs bake-off, success/exit criteria, mutual action plans, technical-win→commercial-win, regulated-POC specifics) → **`tam-operations`** (enterprise POC/POV motion; the dedicated spoke is `enterprise-poc-technical-validation`). This skill covers *why* the regulated buyer demands a heavy pilot, not how to run it.
- **Buyer psychology of buying committees** (DMU/buying-center, Gartner "no decision", champion/mobilizer, switching-cost fear) → **`applied-psychology`** (enterprise B2B buyer psychology reference). This skill covers the *institutional/regulatory* constraints, not the cognitive mechanics.
- **CONSUMER finance / personal banking / personal credit** (deposit accounts, FDIC insurance for a depositor, budgeting, credit scores) → **`consumer-finance`**, **`personal-banking`**, **`consumer-credit-and-debt`**. Those serve a *consumer*; this skill serves an *enterprise seller*.
- **Trading strategy, market mechanics, and securities regulation from a trader/investor lens** → **`trading-and-investing`** family. This skill covers SEC/FINRA only as they shape a *capital-markets firm's technology buying*.
- **General enterprise architecture / cloud well-architected practice** → **`software-engineering-patterns`**, **`aws-cloud`**.

## Provenance & verification

- Built **2026-06-18** via `/dr` deep research (firecrawl/exa + government primary sources). Each reference carries its own numbered `## Sources` list and a `verified-as-of` stamp on volatile claims.
- The two highest-risk, most time-sensitive claims were independently re-verified against primary sources during authoring: the **March 19, 2026 Basel III re-proposal** (federalreserve.gov press release `bcreg20260319a`; OCC Bulletin 2026-9; comments due June 18, 2026) and **SR 26-2 / OCC Bulletin 2026-13** replacing SR 11-7 on **April 17, 2026** (federalreserve.gov `SR2602`; OCC `nr-occ-2026-29`).
- **Re-verify any 2026-dated claim before customer use.** This domain moves quarterly.
