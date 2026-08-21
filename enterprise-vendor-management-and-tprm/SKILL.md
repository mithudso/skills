---
name: enterprise-vendor-management-and-tprm
description: >-
  Buyer-side enterprise vendor governance, procurement & third-party/operational-resilience risk — how large enterprises (esp. BANKS) govern, assess & manage tech vendors, for a B2B/SaaS seller (e.g. a MongoDB TAM) selling INTO them. Educational, NOT legal advice. TRIGGER: TPRM lifecycle & vendor risk tiering, fourth-party/concentration risk; the security-questionnaire gauntlet (SIG, CAIQ/STAR, SOC 2 Type II, ISO 27001/27036, trust centers); procurement & sourcing (Kraljic, RFI/RFP/RFQ, Ariba/Coupa, MFN); SLAs/service credits, vendor scorecards, buyer-run supplier QBRs, vendor consolidation; contract stack (MSA/SOW/DPA, precedence); op-resilience regs (DORA, OCC/FFIEC, UK PRA/FCA CTP, Basel, exit plans). SKIP: clause DRAFTING → career-and-formal-writing; MongoDB compliance config → mongodb-operations-expert; security-control review → security-review; why-banks-buy/FSI → fsi-banking-regulatory-context; buying-committee psychology → applied-psychology; seller-side QBR/POV → tam-operations.
version: 1.1.0
updated: 2026-06-18
category: business-and-gtm
whenToUse: >-
  Use when you need the BUYER-side view of how a large enterprise — especially a bank or other regulated
  financial institution — governs, assesses, procures, and manages its technology vendors, and what that
  means for a vendor (e.g. a MongoDB TAM) selling into them. Covers the TPRM lifecycle and vendor risk
  tiering, the standardized security-assessment / questionnaire gauntlet (SIG, CAIQ, SOC 2, ISO 27001/27036,
  pen-test attestations, trust centers), enterprise procurement & strategic sourcing (category management,
  Kraljic, RFx, e-procurement suites, negotiation levers), SLAs/scorecards/buyer-run QBRs/vendor
  rationalization, the high-level commercial-contract-instrument stack (MSA/SOW/order form/DPA), and the
  operational-resilience & concentration-risk regulatory frontier (DORA, OCC/FFIEC, UK PRA/FCA CTP, Basel).
  NOT for drafting contract clauses, MongoDB product compliance config, or generic security-control review.
keywords:
  - third-party risk management
  - TPRM
  - vendor risk management
  - vendor due diligence
  - vendor risk tiering
  - security questionnaire
  - Shared Assessments SIG
  - CAIQ
  - SOC 2 Type II
  - ISO 27036
  - enterprise procurement
  - strategic sourcing
  - category management
  - Kraljic matrix
  - RFP RFI RFQ
  - source-to-pay
  - Ariba Coupa
  - vendor scorecard
  - SLA service credits
  - vendor consolidation
  - MSA SOW order form DPA
  - operational resilience
  - DORA
  - critical third party
  - concentration risk
  - fourth-party risk
  - exit plan substitutability
  - selling into banks
  - trust center
tags:
  - vendor-management
  - tprm
  - procurement
  - operational-resilience
  - banking
  - gtm
  - third-party-risk
  - enterprise-sales
---

# Enterprise Vendor Management & Third-Party Risk (Buyer-Side, for a Seller)

The **buyer-side governance and regulatory lens** on how large enterprises — especially **banks and other
regulated financial institutions** — govern, onboard, assess, procure, and manage their technology vendors.
The point of view is that of a **B2B/SaaS seller** (concretely, a MongoDB Technical Account Manager selling
into a bank) who needs to understand the machinery on the *other* side of the table: the third-party-risk
team, the procurement function, the legal/contracting desk, and the regulators standing behind all three.

> **What this skill is NOT.** This is governance + regulatory *literacy*, not legal advice and not a how-to
> for drafting contracts or configuring MongoDB. See **Scope & cross-references** below — several adjacent
> skills own the pieces this one deliberately defers.

> **Freshness.** The regulatory and tooling claims here are dated **as of 2026** and move fast (DORA
> oversight, the UK CTP regime, GRC vendor M&A). Treat anything stamped `verified-as-of` as needing
> re-confirmation before you quote it to a customer. Sources are cited inline as `[^n]` → `## References`.

---

## The one-paragraph orientation

When a bank evaluates a vendor like MongoDB, **three machines run in parallel** and a fourth stands behind
them. (1) **Third-party risk management (TPRM)** classifies the vendor by criticality and subjects it to a
risk assessment proportionate to that tier — security, financial, resilience, concentration, compliance. (2)
**Procurement** runs a sourcing process, files the vendor on a portfolio matrix, and applies negotiation
leverage to drive price and terms. (3) **Legal/contracting** imposes the bank's own paper — an MSA stack with
audit rights, exit plans, and liability allocation. (4) Behind all three sit **financial regulators** (EU
DORA, US OCC/FFIEC, UK PRA/FCA, Basel) whose operational-resilience and third-party rules the bank must
*flow down* to the vendor as contractual demands. A seller who understands all four moves faster and gets
surprised less.

---

## Routing table — load the reference for the task

| If the question is about… | Read |
|---|---|
| TPRM lifecycle stages; vendor risk **tiering/criticality**; inherent vs residual risk; risk domains; fourth-party & concentration risk; the TPRM operating model & three-lines-of-defense; GRC/VRM tooling; security-ratings continuous monitoring | `references/tprm-lifecycle-and-risk-tiering.md` |
| The **security-questionnaire gauntlet**: Shared Assessments SIG/SIG Lite/SCA; CSA CAIQ/CCM/STAR; SOC 1/2/3, Type I vs II, bridge/gap letters, CUECs; ISO/IEC 27001:2022 supplier controls & ISO/IEC 27036; pen-test attestation letters; trust centers & questionnaire automation | `references/security-assessments-and-evidence.md` |
| Enterprise **procurement & strategic sourcing**: purchasing vs procurement vs sourcing vs category management; Kraljic matrix; RFI/RFP/RFQ; reverse auctions; S2P/P2P suites & punchout/cXML; procurement-vs-business-buyer dynamic & incentives; negotiation levers; SLAs/OLAs & service credits; vendor scorecards; buyer-run supplier QBRs; vendor consolidation/rationalization | `references/procurement-sourcing-and-negotiation.md` |
| The high-level **contract-instrument map**: MSA, SOW, order form, DPA, SCCs, NDA, BAA; the document stack & order of precedence; vendor-paper vs customer-paper. (MAPPING only — clause drafting is out of scope) | `references/contract-instruments-map.md` |
| **Operational-resilience & systemic third-party regulation**: EU DORA (+ CTPP oversight, Register of Information, TLPT); US OCC/FFIEC/interagency-2023 + Sound Practices; UK PRA/FCA operational resilience + Critical Third Parties (CTP) regime; Basel Principles for Operational Resilience; post-CrowdStrike concentration scrutiny; exit plans / stressed exit / substitutability | `references/operational-resilience-regulation.md` |
| **"So what for a vendor selling in?"** — the parallel-gates playbook, regulatory flow-down, deal-acceleration via evidence package, and the post-sale TAM angle | this file, **Selling-in playbook** below |

---

## Scope & cross-references (what this skill defers, and to whom)

This skill is the **buyer-side governance + regulatory** lens. It deliberately stops at several boundaries and
hands off to a sibling skill. Cross-reference these rather than duplicating them:

- **Contract-clause DRAFTING / wording** — indemnity language, limitation-of-liability construction,
  precedence-clause text, force-majeure wording → **`legal-adjacent-writing`** (a reference under the
  `career-and-formal-writing` hub). This skill *maps* the instruments (what an MSA/DPA/SOW is and where it
  sits); it never tells you how to write a clause.
- **MongoDB-specific compliance** — Atlas SOC 2 / ISO / PCI / FedRAMP scope, CSFLE/Queryable Encryption,
  shared-responsibility config, audit-ready Atlas architecture → **`mongodb-compliance`** (a reference under
  `mongodb-operations-expert`). This skill covers what a *buyer* does with a SOC 2; that one covers what
  MongoDB's own certifications *are*.
- **Generic security-control review / auditing** — OWASP, control testing, secrets/supply-chain/PII review →
  **`security-reviewer`** / **`security-compliance-auditor`** (references under `security-review`).
- **WHY banks buy / FSI fluency** — how banking segments differ as tech buyers, the regulatory pointer-map
  (Basel III, Dodd-Frank, SOX, GLBA, BCBS 239, SR 11-7, 17a-4…), why FSI buyers are risk-averse and run long
  reviews → **`fsi-banking-regulatory-context`**. That skill is the *seller's FSI literacy*; this skill is the
  *mechanics* of the vendor-governance and resilience-regulation machinery the bank runs. They are designed as
  a pair — route segment/why-they-buy questions there, route TPRM/DORA/procurement *mechanics* here.
- **Buying-committee PSYCHOLOGY** — buying center / DMU, Gartner buying group & "no decision," champion /
  mobilizer, switching-cost & lock-in fear → **`applied-psychology`** (enterprise B2B buyer psychology). This
  skill covers the *roles and incentives* of procurement vs the business buyer; that one covers the
  *psychology* of how the committee decides.
- **Seller-side deliverables** — seller-run QBR/EBR, account reviews, support plans, the POC/POV technical-
  validation motion → **`tam-operations`**. Note the mirror-image trap: a **buyer-run supplier QBR**
  (procurement reviewing the vendor against a scorecard) is covered *here*; a **seller-run QBR** (the TAM
  presenting value to the customer) is `tam-operations`.

### When NOT to use this skill at all

Beyond routing adjacent *topics* to the siblings above, do not fire this skill when the *whole question* is
out of frame:

- **Consumer / retail banking or personal finance** (a person's bank account, credit, mortgage) → the
  `consumer-finance` / `consumer-credit-and-debt` families. This skill is about *enterprises buying technology
  vendors*, not individuals using a bank.
- **A non-regulated SMB buyer with no formal TPRM or procurement function** — most of this skill (regulatory
  flow-down, vendor tiering, the questionnaire gauntlet) assumes a mature enterprise/bank buyer; for a small
  informal buyer it is mostly inapplicable.
- **Seller-internal deal-desk / quota / pricing-approval mechanics** (the vendor's *own* internal process for
  approving a discount) → `tam-operations` or sales-ops; this skill is about the *buyer's* machinery.

---

## Selling-in playbook (cross-cutting synthesis — the "so what")

This is the load-bearing payoff for a TAM. The detailed mechanics live in the references; this is how they
combine into what actually happens in a deal and an account.

### The deal runs as parallel gates, and risk/security is usually the long pole
A large-enterprise deal commonly runs **~6 months to ~2 years** with **8–12 stakeholders** (of whom roughly
6–10 form the core buying committee — see `references/procurement-sourcing-and-negotiation.md`), and the
procurement + security-review phase is now the **longest single phase**, not a formality.[^selling1][^selling2]
Roughly concurrent gates: technical evaluation → business case → **procurement negotiation** → **security /
TPRM review** → **legal redlines** → **vendor-master / S2P onboarding**. Deals most often stall *between
technical evaluation and business case* and *during procurement and security review*.[^selling1] Practitioner
figures (treat as directional, not audited): legal reviews ~21–45 days; security assessments commonly
~60–90 days; each questionnaire stall costs ~2–8 weeks of pipeline.[^selling2][^selling3] **As of 2026,**
SOC 2 Type II and SSO support are repeatedly cited as the two most common *hard blockers* — the deal does
not move without them.[^selling3]

### Banks are slower because their diligence is a *regulatory mandate*, and they flow their obligations down to you
A non-regulated enterprise assesses a vendor because it wants to. A bank assesses because a regulator
*requires* it to, across a defined lifecycle (planning → due diligence → contracting → ongoing monitoring →
termination).[^selling4] The decisive consequence for a seller is **regulatory flow-down**: the bank passes
its own regulatory obligations to the vendor as **contract requirements**, not negotiable niceties. Under EU
DORA Art. 30 (and analogous US/UK guidance), critical-function vendor contracts must carry **audit / access /
inspection rights** (for the bank *and* its regulator), **sub-processor controls and disclosure**,
**incident-notification and assistance**, **business-continuity / resilience** obligations, and a documented
**exit strategy with a transition period**.[^selling4][^or-dora-contract] So the bank's TPRM questionnaire,
its audit-rights clause, its sub-processor-approval demand, and its exit-plan ask are all downstream of a
regulator — arguing them away on commercial grounds rarely works; the better move is to *satisfy them with
ready evidence*. (Detail: `references/operational-resilience-regulation.md`.)

### A maintained evidence package is the single biggest accelerant — and the antidote to the "rubber-stamp" critique
The buyer-side trend (as of 2026) is **skepticism of self-attested questionnaires** — only ~34% of TPRM
professionals say they believe questionnaire responses, yet ~81% of returned questionnaires claim *perfect*
compliance, and CSA itself calls questionnaires "a familiar but ineffective norm."[^sa-vanta][^sa-csa] The
counter-move is **verified, pre-assembled evidence**: SOC 2 Type II, ISO 27001:2022, a STAR Level 2
attestation, a pen-test *attestation letter*, and a completed CAIQ — published in a **trust center**
(SafeBase / Whistic / Conveyor / Vanta) so prospects self-serve. Trust centers reduce questionnaire
*volume* (one vendor, SafeBase, claims 74%+ fewer inbound — see `references/security-assessments-and-evidence.md`),
and questionnaire-response automation accelerates the *remainder*; they are complementary, not
substitutes.[^sa-trust][^selling5] **Caveat:** these acceleration
figures originate from the GRC vendors selling the tooling — read them as motivated estimates, and note a
*stale* trust center actively hurts credibility.[^selling5] (Detail:
`references/security-assessments-and-evidence.md`.)

**Evidence-package readiness checklist** (the recurring high-value output — what to keep current in the trust
center; map each artifact to what it proves and how often it goes stale):

| Artifact | What it proves to a bank's TPRM team | Freshness cadence |
|---|---|---|
| SOC 2 **Type II** report | Controls operated effectively over a period (not just designed) | Annual; bridge letter covers the gap to "today" |
| ISO/IEC **27001:2022** certificate | Certified ISMS (flag any stale 2013-edition cert) | 3-yr cycle + annual surveillance |
| CSA **STAR** (Level 2 attestation preferred) | Cloud-specific controls, independently attested | Aligned to the underlying ISO/SOC 2 |
| Completed **CAIQ** | Cloud control answers in a standard form | Refresh on material change |
| **Pen-test attestation letter** (not the full report) | Testing occurred + findings remediated, without exposing exploit detail | Per test cycle (commonly annual) |
| **Sub-processor list** | Fourth-party transparency (e.g., the hyperscaler under Atlas) | Keep current; notice on change (DPA obligation) |
| **BC/DR + exit/portability** summary | Resilience and a workable stressed exit | Reviewed at least annually |

### Where the vendor sits on the buyer's mental map predicts how it gets treated
Procurement files every purchase on a portfolio matrix (the **Kraljic** model: profit-impact × supply-risk).
A vendor classified as **"leverage"** (high spend, low switching risk) gets relentless price competition and
multi-bidding; a vendor classified as **"strategic"** (high impact, high switching risk) gets partnership and
relationship leverage.[^proc-kraljic] In parallel, banks run **vendor consolidation / rationalization** to
cut supplier count and concentrate spend on **preferred** vendors — which *rewards the incumbent that wins
preferred status* and *threatens the tail vendor whose spend gets migrated away*.[^proc-consol] The seller's
strategic goal: be the platform spend **consolidates onto**, not the one consolidated away. (Detail:
`references/procurement-sourcing-and-negotiation.md`.)

### Expect specific, recurring contract frictions — and recognize which are regulator-driven vs commercial
- **Liability cap & indemnity** (highest-stakes commercial fight): banks push for higher or *super-caps* (a
  raised cap, often a multiple of fees, for specified breach types) and *uncapped carve-outs* for data-breach,
  confidentiality, and IP infringement, because an indemnity subordinated to a low general cap is "effectively
  worthless."[^selling-liab] *Commercial* — negotiable.
- **Audit rights, sub-processor disclosure/approval, exit/portability** (and sometimes source-code/SaaS
  escrow): largely *regulator-driven* table-stakes for banks.[^selling4][^or-dora-contract][^selling-escrow]
  Note escrow is contested for pure SaaS (you need live data + environment, not just code), so a bank's
  escrow demand may be more checkbox than true resilience control.[^selling-escrow]
- **Data residency / localization**: standard global data-processing terms can violate local sovereignty
  rules; in-region hosting is often a paid tier.[^selling-residency]
- **MFN / price-protection** and **renewal caps**: commercial; narrow MFN to list pricing, prospective-only.
- **"We use our own paper"** + a **bespoke security questionnaire even when you have a trust center**: banks
  impose their own MSA template and often still send a custom questionnaire because their TPRM team has its
  own requirements.[^selling-paper]

### The post-sale TAM angle: the bank's obligations are *continuous*, so you become the standing channel
A bank's flow-down obligations — audit rights, ongoing monitoring, sub-processor-change notice, annual
re-assessment, resilience/exit testing — are *recurring*, not one-time.[^selling4] So after the sale, the TAM
becomes the standing channel for **recurring evidence/audit requests**, **annual re-assessments**,
**sub-processor-change notifications**, **buyer-run supplier QBRs against a scorecard**, and **exit /
continuity conversations**. Owning these proactively (a current evidence package, a known sub-processor list,
a credible continuity story) is renewal and expansion protection, not overhead.[^selling6]

---

## Anti-patterns (seller-side, when navigating buyer governance)

- **Treating the security review as a checkbox.** As of 2026 it is a *revenue gate*; an incomplete or slow
  response is a cited reason buyers drop vendors.[^selling3]
- **Arguing regulator-driven asks on commercial grounds.** Audit rights, sub-processor disclosure, and exit
  plans flow down from DORA/OCC/PRA — satisfy them with evidence rather than negotiating them away.[^selling4]
- **Letting a SOC 2 stand in for "secure."** A SOC 2 is scoped assurance of controls during a window; it is
  explicitly *not* a guarantee, and configurations drift after the report period.[^sa-soc2myth]
- **Over-relying on self-attested questionnaires (either direction).** Only ~34% of TPRM pros trust them;
  lead with verified evidence.[^sa-vanta]
- **Confusing a buyer-run supplier QBR with a seller-run QBR.** The former is procurement *grading you*
  against a scorecard with renewal/consolidation consequences; clarify whose review it is.[^proc-qbr]
- **Confusing MSA vs order form roles.** Commercial terms (price, quantity, term) live in the order form;
  durable legal terms live in the MSA — and precedence is *not* always "MSA wins."[^contract-precedence]
- **Assuming concentration = automatic regulatory veto.** Regulators (UK/FSB) state concentration does not
  *automatically* pose systemic risk; the demand is resilience/oversight, not a ban on popular
  vendors.[^or-conc-nuance]

---

## References

[^selling1]: Enterprise sales-cycle structure (6–24 mo, 8–12 stakeholders, stall points). Arcade / Vieu / Gain / Rework enterprise-sales analyses — practitioner/sales-authority. https://www.arcade.software/post/enterprise-sales-cycle ; https://try.vieu.com/blog/enterprise-sales-cycle
[^selling2]: Legal-review and security-assessment timings; questionnaire-stall cost. Safe Security 2026 TPRM guide; Security Boulevard — practitioner (figures directional). https://safe.security/resources/blog/2026-guide-to-third-party-risk-management-tprm/ ; https://securityboulevard.com/2026/04/9-critical-security-questionnaire-items-that-stall-enterprise-saas-deals/ — verified-as-of 2026-06-18
[^selling3]: Security review as a revenue gate; SOC 2 Type II + SSO as hard blockers; buyers drop slow vendors. Security Boulevard; Sprinto; SecurityPal — practitioner. https://securityboulevard.com/2026/04/9-critical-security-questionnaire-items-that-stall-enterprise-saas-deals/ ; https://sprinto.com/journey/sales-blockers/what-is-client-security-questionnaire/ — verified-as-of 2026-06-18
[^selling4]: Bank vendor diligence as a 5-phase regulatory mandate; regulatory flow-down into vendor contracts. US interagency third-party guidance explainers (Ncontracts, Elliott Davis); ABA fourth-party. https://www.ncontracts.com/nsight-blog/interagency-guidance-on-third-party-relationships/ ; https://www.elliottdavis.com/insights/u-s-regulators-prioritize-vendor-management-2025
[^selling5]: Trust-center acceleration claims (and the "vendor-motivated / stale center hurts" caveats). Secureframe; SafeBase; TrustCloud — vendor-doc (motivated). https://secureframe.com/blog/what-is-a-trust-center ; https://www.trustcloud.ai/grc/the-power-of-transparency-how-a-trust-center-can-accelerate-enterprise-sales-and-build-credibility/ — verified-as-of 2026-06-18
[^selling6]: Technical Account Manager / CSM post-sale scope (health checks, QBRs, renewals, escalations) and tie to recurring bank obligations. TAM role descriptions + interagency ongoing-monitoring phase. https://www.ncontracts.com/nsight-blog/interagency-guidance-on-third-party-relationships/
[^selling-liab]: SaaS liability cap defaults and the buyer push for uncapped data-breach/IP carve-outs (indemnity worthless if capped low). Borders Law Group; CloudNuro; TermsFeed — law-firm/practitioner. https://borderslawgroup.com/saas-indemnity-provisions-5-things-to-watch-for/ ; https://www.cloudnuro.ai/blog/saas-liability-cap
[^selling-escrow]: Software/SaaS escrow as a bank exit demand, and the contest over whether it delivers continuity for a live SaaS. Escode; Vorys (skeptic) — practitioner/law-firm. https://www.escode.com/resources/what-is-software-escrow/ ; https://www.vorys.com/publication-SaaS-Recovery-and-Escrow-Alternatives-Planning-for-Continuity-on-a-Rainy-Day-in-the-Cloud
[^selling-residency]: Data residency vs sovereignty vs localization; in-region hosting as a paid tier. Alation; WorkOS — practitioner. https://www.alation.com/blog/data-residency-vs-sovereignty-vs-localization-buyers-guide/
[^selling-paper]: Banks impose their own paper and still send bespoke questionnaires. Secureframe; Workstreet; Advantage-FI — practitioner. https://secureframe.com/blog/what-is-a-trust-center ; https://info.advantage-fi.com/en-us/contract-negotiations
[^sa-vanta]: ~34% of TPRM pros trust questionnaire responses; ~81% of responses claim perfect compliance. Vanta — vendor-doc (disconfirming). https://www.vanta.com/resources/security-questionnaires-are-ineffective — verified-as-of 2026-06-18
[^sa-csa]: CSA calls questionnaires "a familiar but ineffective norm." Cloud Security Alliance — standards-body (disconfirming). https://cloudsecurityalliance.org/blog/2025/04/02/why-security-questionnaires-are-a-familiar-but-ineffective-norm-for-assessing-risk
[^sa-trust]: Trust center reduces questionnaire volume; automation accelerates the remainder. SafeBase; The Hacker News writeup — vendor-doc/practitioner. https://safebase.io/resources/security-questionnaires ; https://thehackernews.com/2024/07/how-trust-center-solves-your-security.html
[^sa-soc2myth]: SOC 2 ≠ guarantee of security; point-in-time, config drifts. LinfordCo (CPA firm); Secureframe — practitioner (disconfirming). https://linfordco.com/blog/busting-soc-2-myths/ ; https://secureframe.com/hub/soc-2/what-is-soc-2
[^proc-kraljic]: Kraljic portfolio matrix (profit-impact × supply-risk) and quadrant treatment. CIPS; Art of Procurement; Mindtools — professional-body/practitioner. https://www.cips.org/intelligence-hub/supplier-relationship-management/kraljic-matrix
[^proc-consol]: Vendor consolidation/rationalization; preferred-vendor tiers; incumbent-vs-challenger implication. Zycus; Ramp; Amazon Business — vendor-doc/practitioner. https://www.zycus.com/glossary/what-is-vendor-consolidation
[^proc-qbr]: Buyer-run supplier QBR = customer evaluating the vendor against a scorecard (mirror of seller-run QBR). LearnHowToSource; Supply Chain Mechanic — practitioner. https://blog.learnhowtosource.com/quarterly-business-review/
[^contract-precedence]: Order-of-precedence defaults (often MSA-wins, sometimes order-form-wins for its scope). Law Insider; fynk; Sirion — practitioner/CLM. https://www.lawinsider.com/clause/order-of-precedent ; https://fynk.com/en/clauses/order-of-precedence/
[^or-dora-contract]: DORA Art. 30 mandatory contract clauses (audit, exit, sub-contracting) for critical/important functions. Bird & Bird; DORA Art. 30 portal — law-firm/official-text. https://www.twobirds.com/en/insights/2024/global/dora-some-insights-on-contractual-clauses-in-agreements-between-financial-entities ; https://www.digital-operational-resilience-act.com/Article_30.html — verified-as-of 2026-06-18
[^or-conc-nuance]: Regulators (UK PRA/FCA/BoE, aligned with FSB) state concentration does not *automatically* pose systemic risk. Lexology summary of UK CP; FSB. https://www.lexology.com/library/detail.aspx?g=99f48f00-94c2-45ba-b94a-2d7e0fca51cc — verified-as-of 2026-06-18
