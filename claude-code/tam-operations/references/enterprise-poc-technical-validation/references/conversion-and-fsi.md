# POC conversion & FSI/regulated constraints (deep reference)

Sub-reference of **`enterprise-poc-technical-validation`** (a spoke under the `tam-operations` hub).
Load this file for §8 (POC-to-production conversion + the conversion-statistics provenance table) and
§9 (running POCs under FSI/regulated constraints). The primary `SKILL.md` holds the playbook,
taxonomy, criteria, MAP, roles, scoping, failure-mode catalog, and checklists. The full citation list
lives at the bottom of this file.

> Evidence labels match the primary: **[consensus]** = 3+ independent practitioner sources;
> **[common]** = 2 sources / one strong authority; **[contested]** = sources disagree (contradiction
> preserved); **[stat: verify]** = a specific number with weak or vendor-sourced provenance,
> directional only. Volatile claims dated **as of 2026**; regulatory specifics **verify current**.

## Contents

1. [POC-to-production conversion & the technical-win → commercial-win gap](#8-poc-to-production-conversion--the-technical-win--commercial-win-gap)
2. [Conversion statistics: provenance & skepticism](#conversion-statistics--provenance--skepticism-read-before-quoting)
3. [Running POCs under FSI / regulated constraints](#9-running-pocs-under-fsi--regulated-constraints)
4. [FSI POC do / don't](#fsi-poc-do--dont)
5. [References (full citation list)](#references-full-citation-list)

---

## 8. POC-to-production conversion & the technical-win → commercial-win gap

**The frame** [consensus][^care][^momenta][^backdrop]: "the pilot isn't the finish line; it's the
moment the real buying decision begins." John Care's strict definition of a **Technical Win**: "when
you are informed, *in writing*, that your solution has been accepted and judged superior to the
competition." Most orgs use a far looser definition (merely completing a POC), often with "no proof."
Critically, **the Technical Win is not the deal**: Care calls a loose TW "the equivalent of *not
losing*, which is very different from winning"; "no-one gets paid for a TW; the full business win
involves the transfer of money."

**Why a technical success specifically fails to convert** [consensus]:

1. It resolves the *technical* evaluator's risk but not **finance's commercial risk, legal's
   contractual risk, or IT leadership's governance risk.**
2. The people who ran the eval **cannot approve the purchase** ("dies in committee").
3. The finding was **never translated into the dollar the executive is measured on** (e.g., a POC that
   *found the root cause* but never connected it to the revenue problem the exec owned).
4. **Production cost was never scoped.** The technology is the cheap part; connectors, security review,
   workflow redesign, and training are not. [stat: verify] Production is often quoted at ~3–5× the
   prototype (a repeated vendor estimate).

**What bridges the gap** [consensus]:

- **Economic buyer in the room, especially at the readout.** "The readout meeting is the gate that
  converts a POC into a business decision"; book it *at kickoff* so the decision happens when momentum
  is highest.
- **Value tied to dollars in the buyer's language,** in a finance-auditable model the buyer owns ("when
  Finance can reproduce the number without you, approvals accelerate").
- **Pre-agreed next step / conditional close:** "If we hit these KPIs by [date], we move to contract
  review," secured *before* kickoff. This prevents the "great POV, but we need to think about it" stall.
- **Parallel tracks** (technical + business case advance together) and a fast transition to MAP/contract
  while post-success momentum is high.

**The "innovation budget / tourist money" problem** [consensus][^idc]: many pilots are "born at the
board level" and "highly underfunded or not funded at all; trickle-down economics," so the POC happens
"not because of a strong business case." The pilot becomes "an institutional hedge; it signals action
without committing to the cost or accountability of production. No one explicitly kills it. It just
never moves forward." **Detect early:** a buyer hesitant to define success metrics or name the economic
buyer is likely spending tourist money. On the *vendor* side, counting paid pilots as ARR inflates true
ARR when conversion is <50%.

### Conversion statistics: provenance & skepticism (READ BEFORE QUOTING)

There is **no single trustworthy "POC-to-production conversion rate."** Numbers come from different
populations (DoD AI vs enterprise GenAI vs B2B SaaS), measure different things (project failure vs
value-delivery vs production deployment), and the most viral one is contested. [contested]

| Statistic | Source & tier | What it actually measures | Verdict |
|---|---|---|---|
| "80% of AI projects fail" | **RAND, Aug 2024**: primary research, 65 practitioner interviews, DoD-oriented | Failing to deliver *intended value* (84% cited leadership causes) | Credible primary source, but a *value-failure* rate, not a conversion rate |
| "88% never reach production (4 of 33 POCs)" | **IDC / Lenovo**: vendor-partnered research | AI POCs reaching wide-scale deployment | Widely repeated; vendor-partnered → **qualify** |
| "95% of (gen)AI pilots fail" | **MIT NANDA, July 2025**: self-described *preliminary* working paper, not peer-reviewed | On its own funnel ≈ **75% of pilots** vs a harsh 6-month-P&L bar; the *same report* shows ~83% pilot-to-implementation for general tools, ~67% with vendors | **Contested / over-read**; do not cite as "19 of 20" |
| "30% of GenAI projects abandoned after POC by end-2025" (later ≥50%) | **Gartner, July 2024**: analyst *forecast* | A prediction of post-POC abandonment | Primary Gartner forecast; qualify as a forecast |
| "78% of Global 2000 IT execs: <half of POCs reach production" | **Sapphire Ventures, 2020**: VC-run survey | POC→production, enterprise IT | Genuine primary survey predating the AI hype; good corroboration that **<50% conversion is the historical norm** |
| "85% of AI projects fail (poor data quality)" | Attributed to "Gartner 2025" across vendor blogs | — | **Provenance weak / likely misattributed; do not cite as fact** |
| B2B SaaS POC conversion "50–90%"; "60% with structure vs 30–40% without" | Vendor blogs | POC/pilot → contract | Vendor-sourced, uncorroborated → **tentative** |

**Synthesis:** the *direction* is well-supported across independent populations: **most enterprise
pilots/POCs do not convert to production**, and "<50% convert" (Sapphire 2020) predates the AI-specific
80–95% claims. The precise rate is population-dependent and the most-quoted figure (95%) is the least
sound. The optimistic counter-read: the same MIT data shows vendor-partnered pilots ~67% and general
tools ~83%; failures concentrate in *internal builds* and *flashy use cases*, not the technology.
**Quote RAND / IDC / Sapphire with caveats; treat 95% as contested.**

---

## 9. Running POCs under FSI / regulated constraints

What is **different and harder** at a regulated buyer (bank, insurer, financial services). All
regulatory specifics are **as of 2026, verify current**; do not assert SR-letter numbers or rule
clauses without checking.

> Cross-reference: for FSI buyer fluency (how bank segments buy, why they are risk-averse, the
> regulatory map at pointer level), see the `fsi-banking-regulatory-context` skill. This section is the
> POC-execution overlay: how those constraints reshape the technical-validation motion.

**The structural fact** [consensus][^canapi][^msftfsi]: at a regulated buyer, **the POC is not the
gate; the third-party risk management (TPRM) / vendor security review is**, and it runs on a clock
measured in **months**. A POC "does not remove any required vendor-management step; it only creates a
streamlined version of them up front." So the POC and the security review **coexist**; the POC does not
waive the review. The winning play is to run **assurance, legal, and technical validation as parallel
workstreams from day one.**

**1. Security-review sequencing.** The bank's TPRM lifecycle (per the **2023 Interagency Guidance on
Third-Party Relationships**, OCC/FRB/FDIC joint, *verify current*): **Planning → Due Diligence &
Selection → Contract Negotiation → Ongoing Monitoring → Termination.** Information-security due
diligence (your infosec program, SDLC, **penetration-test results**) sits *before* the bank commits.
The questionnaire stack: **SIG** (SIG-Lite ~150–200; SIG Full **800+**), **CAIQ** (~261; submit once to
the **CSA STAR Registry** and reuse), VSAQ, the bank's custom DDQ, with **SOC 2 Type II** (Type I
"rarely accepted"), ISO 27001, pen-test summaries, and a DPA demanded alongside. [consensus] [stat:
verify] Vanta-sourced figures ("~78% of companies report security reviews delayed a deal," "~67% of
B2B deals require a questionnaire before close") are widely re-reported from one survey; directional.
The classic failure is scheduling the InfoSec review *after* contract execution.

**2. Data handling: you will not get production data.** [consensus, tier-1 regulator + corroboration]

- **GLBA / NPI:** banks must protect nonpublic personal information and **oversee service providers by
  contract** (Safeguards Rule).
- **PCI DSS v4.0:** production account data, including PANs, **must not be used in test/development
  environments** (practitioner sources cite Req 6.5.4 / 3.x; *verify the exact clause against the PCI
  SSC standard*).
- **GDPR / residency:** sharing even *masked* production data can still trigger GDPR (the masking
  process touches personal data, and residual re-identification risk remains); only **born-synthetic**
  data (no real data at any stage) clearly escapes the obligation.
- **Working distinction:** **default to synthetic** for the POC/sandbox layer; reserve **masked
  production** for later stages where exact fidelity matters (rare events, edge cases). Where data may
  live: on-prem, the **customer's own VPC/tenant**, or a vendor environment *only* if contractually
  covered and de-scoped of real data.
- **Contrarian / critical caveat (important)** [tier-1, FCA]: **synthetic data is NOT a free pass.** The
  FCA is explicit that synthetic data **complements, not replaces, live operational data**, and the
  privacy–utility tradeoff means more-realistic data carries higher re-identification risk. Practitioner
  skeptics note **no major regulator has endorsed synthetic data for AML model validation,** so
  compliance leaders still rely on masked/anonymized production data for true model tuning. **TAM
  implication:** synthetic is the right *POC-stage* answer; do **not** promise it satisfies *production
  model validation.*

**3. Environment provisioning.** [consensus][^msftfsi] A bank sandbox is **"not a lightly governed dev
subscription"**; it is a **fully isolated tenant** with **no connectivity to corporate
networks/identity/production** and **no production or customer data; only synthetic or public
datasets**, run with *lighter controls*. "Isolation, not internal governance alignment, is the critical
prerequisite for safe experimentation." Production-like (enterprise-ready) environments add
non-negotiable baseline controls (private endpoints only, customer-managed keys, managed identities/no
static secrets, multi-region resiliency), which is exactly why a sandbox result is not a production
guarantee. Provisioning friction: identity/privileged-access management is a **prerequisite**, not a
follow-up; landing-zone/account provisioning can take **months** at a large bank; change-advisory-board
approval adds weeks. **Accelerator:** buying via a **cloud marketplace** (AWS/Azure/GCP) private offer
can ride the bank's existing vetting + committed spend and collapse onboarding "from months to days."

**4. Timeline realities & the regulatory bodies.** [Medium: practitioner estimates / tier-1 for rule
existence] Total bank sales cycle **6–18 months**, run as parallel tracks. New-cloud-service enablement
inside a large bank "commonly requires multiple months even when operating efficiently." SOC 2 Type II
itself takes 6–12 months to *obtain*, a prerequisite you must already hold. The instruments that shape a
regulated POC (all **verify current**):

| Instrument | Body | What it forces |
|---|---|---|
| **2023 Interagency Guidance on Third-Party Relationships** (OCC Bulletin 2023-17) | OCC + FRB + FDIC | 5-stage TPRM lifecycle; infosec + pen-test due diligence before engagement; heightened oversight for "critical activities" |
| **SR 11-7** Model Risk Management | FRB + OCC | Any model used for decisions must be validated, governed, *understood*; push vendors for conceptual summaries + outcomes testing; the focal point for AI/ML POCs |
| **FFIEC** guidance | FFIEC | IT / third-party-risk examination expectations |
| **NYDFS 23 NYCRR §500.11** | NY DFS | A written third-party security policy; vendors must meet minimum cybersecurity practices |
| **DORA** | EU / ESAs | In force 17 Jan 2025; Register of Information on ICT third-party contracts; Critical Third-Party Provider oversight; TLPT |
| **SS2/21** Outsourcing & **SS1/21** Operational Resilience | UK PRA / BoE (+ FCA) | Materiality assessment, due diligence, BCP/exit for material outsourcing (2026 op-resilience policy statements issued; verify) |
| **EU AI Act** | EU | Classifies e.g. credit scoring as high-risk AI → sharper explainability constraints |

**The quantified anti-pattern** [tier-1-grade][^msftfsi]: "**Four sequential reviews of two weeks each
produce two months of elapsed time with days of actual work.**" The fix is to **parallelize**:
architecture design alongside questionnaire completion, IaC alongside control validation, threat
modeling alongside sandbox experimentation.

**Contrarian view on the sandbox** [consensus]: a strong skeptic camp argues bank pilots systematically
mislead because they "succeed by avoiding reality": clean curated extracts, no integration, no
compliance, friendly users; "if the pilot and production use different data, the pilot result is not a
reliable predictor." **TAM implication:** design the regulated POC to *test the production-hard things*
(legacy-core integration path, governance accountability, explainability, messy real-pipeline data) and
make those first-class exit criteria, not a clean demo.

### FSI POC do / don't

| Do | Don't |
|---|---|
| Kick off the security review **in parallel on day 1** | Treat security as a post-pilot formality (the #1 failure mode) |
| Hold **SOC 2 Type II + pen-test + ISO 27001** before outreach | Start your SOC 2 at deal time (it takes 6–12 months) |
| Keep a master **SIG/CAIQ answer library**; STAR-register CAIQ | Re-author every questionnaire from memory (inconsistency = red flags) |
| Scope the POC to an **isolated sandbox + synthetic/masked data** | Ask for or accept **production NPI / PCI / PII** in a POC |
| Position synthetic data as the **POC-stage** answer | Promise synthetic data will satisfy **production model validation** |
| Run the POC in the **customer's VPC/tenant** where possible | Demo only in your pre-tuned vendor sandbox and call it production-ready |
| Pre-package your **4th-party / subprocessor** docs | Surprise the diligence team with an undisclosed cloud-subprocessor chain |
| Build **explainability + governance** into AI POCs from the outset | Optimize for demo accuracy and defer the regulator/explainability question |
| Use **cloud-marketplace private offers** to shorten procurement | Default to a months-long net-new onboarding path when a shortcut exists |
| Expect **6–18 months**; parallelize workstreams | Promise a leadership close date that assumes serial reviews |

---

## References (full citation list)

Evidence base is overwhelmingly practitioner lore (SE books, presales communities, vendor blogs); the
closest-to-empirical sources are flagged. Volatile claims dated **as of 2026**; regulatory specifics
**verify current**. These footnotes back claims in BOTH the primary `SKILL.md` and this sub-reference.

[^care]: John Care, *Mastering Technical Sales: The Sales Engineer's Handbook* (4th ed., Artech House) and "The Technical Win." Canonical SE authority: POC 7-phase structure, "technical win" vs business win, evaluation as the parent term. https://masteringtechnicalsales.com/
[^cohan]: Peter Cohan, *Great Demo!*: Vision Generation vs Technical Proof, "Other Forms of Proof," footprint↔proof-effort, believable environments. https://greatdemo.com/ (methodology described as validated across thousands of demos; closest to empirical here).
[^homerun]: Homerun Presales, "Evaluation Plans: The Key to Great POCs and POVs" / "POCs and POVs That Actually Win Deals." Presales-tooling authority. https://www.homerunpresales.com/blog/poc-best-practices-evaluation-plans-the-key-to-great-pocs-and-povs (2026).
[^guideflow]: Guideflow, "POV vs POC: guide for SaaS teams" and "Sales POC that closes." Vendor practitioner; clean POC/POV comparison + "where the doubt sits" test. https://www.guideflow.com/blog/pov-vs-poc (2026).
[^steerlab]: Steerlab, "What Is a Proof of Value (POV)" / "What Is a Proof of Concept (POC)." Practitioner; sharp technical-vs-business framing. https://www.steerlab.ai/blog/what-is-proof-of-value-pov (2026).
[^vivun]: Vivun, "POC vs. POV: A Seller's Guide." Vendor; key disconfirming source ("terms often used interchangeably… same core purpose"). https://www.vivun.com/blog/poc-vs-pov-a-sellers-guide (2025).
[^rework]: rework.com, "PoCs That Predict Success and PoCs That Waste Everyone's Time." Buyer/RevOps perspective (independent counterweight); written-criteria, vendor-vs-buyer definition conflict, MAP "when they help and hurt." https://resources.rework.com/insights/saas-buying/pocs-that-predict-success (2026).
[^30mpc]: Armand Farrokh, 30 Minutes to President's Club, "My 4-Step POC Framework." Leading sales newsletter; controllable exit criteria, "avoid the POC," conditional-close question. https://www.30mpc.com/newsletter/my-4-step-poc-framework
[^nedl]: Ashish Jaiman, nēdl Labs, "Escaping Pilot Purgatory: Why Technical Success is a Vanity Metric." Practitioner (ex-Microsoft); "tourist money," conditional close, good/bad criterion examples. https://nedllabs.com/blog/escaping-pilot-purgatory (2025).
[^pointer]: Pointer Strategy, POC scope & success-criteria enablement certification standard (explicit failure definition). Practitioner.
[^presalespulse]: PreSales Pulse, "POC Management Playbook for SEs." Presales practitioner; "When to Walk Away" kill criteria, 2–6 wk time-box, provisioning lead time. https://presalespulse.com/careers/poc-management-playbook/ (2026).
[^meddicc]: MEDDICC LTD / Andy Whyte: Champion, Economic Buyer, Decision Criteria pages. Dominant enterprise qualification framework (used here for role vocabulary). https://meddicc.com/what-is-meddpicc/ (cross-ref `tam-commercial-metrics` for scoring).
[^forcemgmt]: Force Management, "Coaches vs. Champions," "Leverage Your Champion," Command of the Message / "How MEDDICC Helps Win with Decision-Makers." Value-selling authority; 3-part champion test, above/below-the-line metrics, cost of not funding. https://www.forcemanagement.com/
[^spotlight]: Spotlight.ai, value-engineering / Business Value Assessment series. Vendor (value-consulting) but the most detailed on BVA mechanics: cost-of-inaction, hard/soft/revenue split, value-hypothesis sequencing. (2026).
[^growthhub]: The Growth Hub, "Sales Autopsy: The 'Flawless' Pilot that Stayed a Pilot." Practitioner case; "won on accounting," Quantified Risk-of-Inaction. https://thegrowthhub.me/sales-autopsy-technical-vs-financial-gap/ (2026).
[^prolifiq]: Prolifiq, Mutual Action Plan guidance. Practitioner; buyer-ownership as the strongest deal-health signal. https://www.prolifiq.com/post/mutual-action-plan-template
[^rafiki]: Rafiki AI, "Mutual Action Plan: Co-build the Path to Close." Practitioner; reverse-timeline-from-go-live, the champion-CFO litmus test. https://getrafiki.ai/sales-strategy/mutual-action-plan-template-co-build-path-to-close/
[^itsjustrevenue]: It's Just Revenue, Mutual Action Plan + "Pilot-to-Production Conversion." Practitioner (independent-leaning); "the plan is the qualification"; readout discipline. https://www.itsjustrevenue.com/insights/mutual-action-plan
[^prospeo]: Prospeo, SE role guide, "MAP with teeth," technical-vs-economic buyer. Practitioner. https://prospeo.io/
[^nobel]: Nobel Recruitment / Modern Presales, SE-owns-technical / AE-owns-commercial split; SE-vs-TAM lifecycle. Practitioner. https://www.modernpresales.com/blog/proof-of-value-vs-proof-of-concept (2026).
[^challenger]: CEB/Gartner, *The Challenger Customer*: Mobilizer/Talker/Blocker stakeholder profiles (n≈1,460 / ~5,000 stakeholders). Empirically grounded. https://a.sfdcstatic.com/content/dam/www/ocms/assets/pdf/misc/The-Challenger-Customer.pdf (+ Gartner B2B Buying Journey: 6–11 stakeholders, six non-linear buying jobs).
[^pulse]: Pulse RevOps, "The Competitive Knockout Session" / "Incumbent Displacement Map." Sales-training practitioner; "change what gets evaluated," incumbent wins the spec-sheet bake-off. https://pulserevops.com/ (+ Stealery, The Revenue Architect, SCOUT; 6+ convergent sources).
[^infoweek]: InformationWeek, "How to Make Vendor Technology Bakeoffs Work." Industry press; bake-off = simulated-production, buyer's own processes, ~2 finalists, post-RFP. https://www.informationweek.com/software-services/how-to-make-vendor-technology-bakeoffs-work
[^itbrief]: Rajesh Agnihotri, IT Brief, "Stop confusing demos with POCs." Practitioner op-ed; "a great demo earns the right to a POC; a great POC earns the right to a PO," Definition of Done, POC purgatory. https://itbrief.asia/story/stop-confusing-demos-with-pocs-your-pipeline-depends-on-it (2026).
[^momenta]: Momenta Ventures / Product Market Pro, "Your Champion Can't Buy Your Product" / "Why the Deal Stalls After the Pilot Succeeds." VC/operator; single-threading, "pilot tests whether the product works, not whether the org can buy." https://momenta.vc/insights/your-champion-cant-buy-your-product-and-thats-why-growth-stops (2026).
[^backdrop]: Backdrop (getbackdrop.ai), "Proof of Concept: How to Structure One That Closes" / "Technical Win." Vendor; readout-as-gate, EB at readout, agree-in-writing. https://www.getbackdrop.ai/learn/proof-of-concept (2026).
[^idc]: Ashish Nadkarni (IDC), via CIO.com, "88% of AI pilots fail to reach production." Trade journalism + analyst; the "4-of-33" stat + "trickle-down / underfunded board-level" tourist-money quote. https://www.cio.com/article/3850763/ (2025). RAND, "The Root Causes of Failure for AI Projects" (RRA2680-1, 2024): primary research, "80% fail." https://www.rand.org/pubs/research_reports/RRA2680-1.html . MIT NANDA, "The GenAI Divide" (2025): contested 95% claim. Gartner press release (Jul 2024): "30% abandoned after POC." Sapphire Ventures (2020): "78% say <half of POCs reach production." https://sapphireventures.com/blog/over-50-of-proof-of-concepts-fail-heres-how-to-fix-yours/
[^canapi]: Canapi, "Demystifying the Bank Vendor Management Process." Fintech-VC practitioner; "a POC does not remove… vendor management steps." https://www.canapi.com/insight/demystifying-the-bank-vendor-management-process-part-i
[^msftfsi]: Microsoft, "Cloud Service Enablement Framework for Regulated Financial Services." GSIB-grounded standards-grade doc; isolated-tenant sandbox, no prod data, sequential-vs-parallel ("four 2-week reviews = two months"), "multiple months to enable services." https://learn.microsoft.com/en-us/industry/financial-services/service-enablement-framework (updated 2026). Regulator/standards anchors (verify current): OCC Bulletin 2023-17 / 2023 Interagency Guidance on Third-Party Relationships (https://www.occ.gov/news-issuances/bulletins/2023/bulletin-2023-17.html); FRB/OCC SR 11-7 Model Risk Management; NYDFS 23 NYCRR 500; EU DORA (in force 17 Jan 2025, ESMA); UK PRA/BoE SS2/21; FTC GLBA Safeguards Rule; PCI DSS v4.0 (verify clause). FCA, "Synthetic Data to Support Financial Services Innovation" (synthetic complements, not replaces, live data). https://www.fca.org.uk/publication/call-for-input/synthetic-data-to-support-financial-services-innovation.pdf . Finextra, "The Pilot Graveyard" (explainability/legacy-integration walls). https://www.finextra.com/blogposting/31141/
[^forrester]: Forrester (Lisa Singer), "Proof Is the Product: How Trials and POCs Have Become a Real Go-To-Market Motion." Analyst; the 2026 shift toward outcome-driven proof. https://www.forrester.com/blogs/proof-is-the-product-how-trials-and-pocs-have-become-a-real-go-to-market-motion/ (2026).
