# FSI / Big-Bank Buying Behavior — deep dive

> Reference of the `enterprise-b2b-buyer-psychology` spoke (applied-psychology hub). Load for an FSI / large-bank deal. Companion to the main `../SKILL.md` (read Section 6 there first).
>
> `verified-as-of: 2026-06-18` for all regulatory and statistical claims. **The regulatory landscape is fast-moving — verify each item against the primary regulator before relying on it.** Sourced primarily from regulator text (federalreserve.gov, occ.gov/occ.treas.gov, fdic.gov, ffiec.gov, eur-lex.europa.eu, bankingsupervision.europa.eu) plus practitioner accounts of bank sales cycles.

## Contents
- Why FSI is the hardest enterprise buyer
- Procurement & vendor-risk as gatekeepers
- US third-party / vendor-risk regulation
- Model risk (SR 11-7 and the 2026 revision)
- EU DORA and the exit-plan/concentration mandate
- Security, data residency, audit, concentration & exit
- Conservatism, herd behavior, "proof not promise"
- Disconfirming evidence (are banks faster now? is lock-in regulation overstated?)
- Vendor implications (selling a database into a flagship bank)
- Sources

## Why FSI is the hardest enterprise buyer

Every dynamic in the main spoke intensifies in a large regulated financial institution: a **larger buying center with more independent vetoes**, **longer cycles**, and, the distinctive feature, **risk and lock-in are regulatory obligations, not merely psychology.** The result is that **bank buying is assurance-led, not demo-led: satisfy risk before expecting a commercial decision.**[^fsisalesmotion]

## Procurement & vendor-risk as gatekeepers

Large banks run the most formalized procurement in the economy: structured RFPs with scoring rubrics, reference checks, and **5–8 stakeholders across business, technology, risk, compliance, and procurement**, each often holding **independent veto power.**[^fsisalesmotion][^fsivarga][^a16z] Cycle facts (ranges vary by institution; treat as QUALIFIED):

- Enterprise bank cycles commonly run **9–18 months** (some sources 6–18).[^fsisalesmotion][^fsispecs]
- The **vendor risk assessment alone takes 3–6 months**; CISO/InfoSec assessment **6–8 weeks even with full cooperation**, then legal/compliance with no fixed timeline.[^fsispecs]
- **"The sponsor is your champion; procurement is your buyer"** — the business sponsor cannot approve alone; they must clear risk/InfoSec/compliance/legal gates first.[^fsivarga]

## US third-party / vendor-risk regulation (as of 2026)

**Interagency Guidance on Third-Party Relationships: Risk Management (June 2023).** Issued jointly by **OCC, Federal Reserve, FDIC** — FRB **SR 23-4**, **OCC Bulletin 2023-17**, **FDIC FIL-29-2023**.[^ig-frb][^ig-occ][^ig-fdic] It **replaced** the agencies' prior separate guidance (OCC 2013-29, FRB SR 13-19, FDIC FIL-44-2008) and sets out a **risk-management life cycle**: planning → due diligence / third-party selection → contract negotiation → ongoing monitoring → termination.[^ig-frb] **Legal nuance:** it is **supervisory guidance that "does not have the force and effect of law and does not impose new requirements"** — yet examiners review third-party risk management in standard supervision, so it shapes bank behavior strongly.[^ig-frb]

**OCC Heightened Standards (12 CFR Part 30, Appendix D).** Applies to **"covered banks", generally ≥ $50 billion in average total consolidated assets.**[^hs-cfr][^hs-occ] Requires a **formal written risk-governance framework** designed by independent risk management and approved by the board, with a **three-lines-of-defense** model. Strategic-risk factors examiners weigh explicitly include **"risks associated with new activities or technologies, particularly those that are innovative or unproven"** — i.e., adopting a novel vendor is itself a supervised risk.[^hs-large]

**FFIEC outsourcing guidance.** The **FFIEC IT Examination Handbook** ("Outsourcing Technology Services" booklet; Management booklet §III.C.8) requires a comprehensive outsourcing risk-management process and states the controlling principle that outsourced relationships **"should be subject to the same risk management, security, privacy, and other policies… as if conducted in-house,"** with regulators retaining **statutory authority to examine a third party's activities.** Third-party management is framed as **an enterprise-wide governance issue, not a technology issue.**[^ffiec1][^ffiec2][^ffiec3]

## Model risk (SR 11-7 and the 2026 revision)

The foundational model-risk guidance is **SR 11-7 / OCC Bulletin 2011-12** (FDIC adopted 2017 via FIL-22-2017): models must be managed with **conceptual-soundness validation, ongoing monitoring, outcomes/back-testing analysis, and governance.**[^mr-fdic][^mr-occ2011] **CRITICAL UPDATE (as of 2026):** the agencies issued **revised model-risk guidance in April 2026 — FRB SR 26-2 and OCC Bulletin 2026-13** — clarifying a risk-based, tailored approach (and, per related OCC guidance, that community banks aren't required to perform annual validation). The revised guidance still **explicitly covers vendor/third-party models, which must be validated like in-house models.**[^mr-occ2026][^mr-frb2026] **This is fresh and still settling, verify current text before relying.** Direct vendor impact: **if a bank uses your platform to build or run models, your component can be pulled into model validation.**

## EU DORA and the exit-plan / concentration mandate

Unlike US supervisory guidance, **DORA (Digital Operational Resilience Act, Regulation (EU) 2022/2554)** is **binding regulation**, applied from **17 January 2025.**[^dora] Key ICT third-party provisions:

- **Art. 28** — maintain a **strategy on ICT third-party risk**, a **Register of Information** on all ICT arrangements, **pre-contracting risk assessment**, and **documented exit strategies** for services supporting **critical or important functions.**[^dora-ch5][^dora-cyadviso]
- **Art. 29** — **mandatory ICT concentration-risk assessment** (single-provider and sector-level) **before** entering new arrangements.[^dora-ch5][^dora-cyadviso]
- **Art. 30** — **mandatory contractual clauses**; for critical/important functions, *additional* clauses on **audit rights, subcontracting controls, data return/access, and exit assistance with a mandated adequate transition period.**[^dora-art30][^dora-ch5]
- **Concentration risk is handled deliberately, not by hard caps:** DORA explicitly **declines to impose rigid caps** on ICT exposures, opting for a "flexible and gradual approach," with a **Lead Overseer** monitoring **critical ICT third-party providers (CTPPs)** at the systemic level.[^dora]
- **Level-2 RTS (Delegated Reg. (EU) 2024/1773 & 2024/1774)** make the **exit plan mandatory per critical contractual arrangement**, require it to be **realistic and tested against "unforeseen and persistent service interruptions,"** and require **effective audit/access rights** for the entity, its auditors, and competent authorities, including **pooled audits.**[^dora-rts][^dora-ivass]

**Corroborating EU supervisory weight:** the **EBA Guidelines on Outsourcing (EBA/GL/2019/02)** require **unrestricted audit/inspection rights cascaded to subcontractors** and a risk-based approach to **data location**; the **ECB Guide on outsourcing cloud services (2025)** treats **vendor lock-in, concentration (provider/geography/functionality), data residency, and exit strategy** as explicit supervisory review items.[^eba][^ecb]

## Security, data residency, audit, concentration & exit

- **Security/InfoSec review is the dominant time sink.** The fix: a **"vendor trust pack" (SOC 2, pen-test results, BCP/DR, DPA, pre-negotiated MSA) built *before* the first enterprise conversation.**[^fsispecs][^fsibishop]
- **Data residency / sovereignty:** banks must control **where data is stored and processed**; supervisors advise a **list of permitted countries** factoring legal/political risk (including third-country legislator reach). Cloud providers respond with **region selection, attestation reports (SOC/ISO), and commitments not to move data outside the chosen region.**[^ecb][^res-aws]
- **"Concentration risk" + "exit plan" make lock-in a REGULATORY concern, not psychological.** DORA Arts. 28/29/30 *legally require* concentration-risk assessment and **credible, tested exit plans** for critical functions; the ECB names vendor lock-in as a risk to scrutinize. **"An untested exit is not credible."** For a regulated FI, inability to leave a vendor is a **supervised deficiency** — fundamentally different from the psychological switching-cost framing in the main spoke's Section 4.[^dora-ch5][^ecb][^dora-cloudexit]

## Conservatism, herd behavior, "proof not promise"

FSI buyers increasingly demand **"proof, not promise", referenceable outcomes, case studies, and third-party validation are decisive**; "banks buy **defensible execution, not innovation hype.**"[^fsibishop][^ccgroup] a16z notes founders face **internal competition from a risk-averse culture and buyers who don't know how to evaluate a startup.**[^a16z] The conservative-committee dynamic prioritizes **risk avoidance over innovation speed**, and the "**reference-able peer bank**" instinct (a form of herd behavior, *the* "nobody got fired" heuristic applied to peer institutions) is a primary risk reliever.[^fsidanish]

## Disconfirming evidence

**Are banks actually faster buyers now?** Partially, but deliberately, not at SaaS speed. Bain reports FSIs embedding procurement earlier and shifting category strategy "toward innovation, speed, resilience"; KPMG/Accenture report many FIs self-identifying as "innovators" and planning **more open-source over proprietary vendors** to cut cost and reduce concentration; cloud adoption (partial/complete) jumped **37%→91% (2020→2023).**[^neg-bain][^neg-kpmg][^neg-cloud] **BUT** independent buyer research finds deals got **bigger and more deliberate, not faster** — more deals over £500k, cycles still 6–12 months, with compliance/diligence front-loaded.[^ccgroup] Cloud-marketplace standard contracts (e.g., AWS Marketplace) are a real accelerant by **standardizing legal terms** so not every tool triggers a fresh legal cycle.[^neg-endava] **Net: faster on-ramps and rising appetite, but governance gates remain the binding constraint.**

**Is vendor lock-in regulation overstated?** There is a genuine industry counter-argument. The **GFMA/SIFMA Cloud Portability** white paper argues mandated **multi-cloud/portability "could… likely double implementation costs"** and may be impractical (banks long accepted mainframe concentration); **PIFS/Harvard** cites Treasury/MAS feedback that multi-cloud is "too technically complex." ECIPE refines it: **market concentration itself isn't the core risk, *un-contestability* (structural barriers to switching) is.**[^neg-gfma][^neg-pifs][^neg-ecipe] **Counter-counterpoint:** American Banker and supervisors warn banks **overestimate their ability to migrate workloads rapidly** and that exit strategies must be **regularly tested.**[^neg-ambanker][^ecb] **Net: the *requirement* (assess concentration, hold a tested exit plan) is real and binding under DORA; the *prescription* of multi-cloud/portability is contested, present both sides.**

## Vendor implications (selling a database into a flagship bank)

Synthesis of the above (the MongoDB-specific application is reasoned synthesis from cited sources):

1. **Develop a Mobilizer-grade champion, then arm them for an *assurance* fight, not a feature fight.** The champion's hardest internal battles are with risk, InfoSec, compliance, and procurement. Pre-build a consensus toolkit *tilted to risk* — business case/ROI **plus** a vendor trust pack (SOC 2, pen-test summary, BCP/DR, DPA, data-residency/region story, pre-negotiated MSA).[^fsibishop][^fsispecs]
2. **De-risk = clear the regulatory gates proactively.** Map your offering to the bank's obligations: third-party risk life cycle (Interagency Guidance / FFIEC); concentration risk + a credible, *testable* exit plan (DORA Arts. 28/29/30); data residency / region pinning + audit rights cascaded to subcontractors (EBA/ECB); and, if the platform runs models, model-validation support (SR 11-7 / the 2026 SR 26-2 revision).
3. **For a database vendor specifically, exit/portability and data-residency are existential selling points.** Demonstrable **data egress, no proprietary one-way doors, multi-region/data-locality controls, and a rehearsed migration path** turn DORA's "credible exit plan" mandate from a blocker into a **differentiator.**[^dora-cloudexit][^neg-gfma]
4. **Engage risk/InfoSec/procurement early and in parallel,** not after the demo, the assurance phase dominates the 9–18-month cycle.[^fsisalesmotion][^fsispecs]
5. **Multi-thread by role** — banks have 5–8 vetoes; one champion who leaves or can't satisfy the CISO sinks the deal. Cover EB, champion, technical/security evaluator, compliance, procurement, end-user, *with* the champion (air cover), not around them.
6. **Lead with "defensible execution" and referenceable peer banks** matched to size/regulatory profile.[^fsibishop][^ccgroup]

## Sources

[^fsisalesmotion]: Salesmotion, "How to Sell to Financial Services Companies" (assurance-led; 5–8 stakeholders; vendor risk 3–6 months): https://salesmotion.io/blog/how-to-sell-to-financial-services . Tier: blog.
[^fsivarga]: James Varga wiki, "How to sell fintech products to banks" ("sponsor is your champion; procurement is your buyer"): https://wiki.jamesvarga.com/knowledge/fintech-gtm/selling-to-banks . Tier: blog.
[^fsispecs]: FintechSpecs, "Why Fintech SaaS Sales Cycles Are Slow" (9–18 months; CISO 6–8 weeks; vendor trust pack): https://fintechspecs.com/blog/why-fintech-saas-sales-cycles-are-slow/ . Tier: blog.
[^fsibishop]: Stacy Bishop, "How Fintech Founders Sell to Banks" (defensible execution; trust pack; proof over promise): https://www.stacybishop.com/articles/how-to-sell-fintech-to-banks . Tier: blog.
[^fsidanish]: DanishLeadCo, "How to Accelerate FinTech Sales with Buying Committees" (risk-avoidance over speed; conservative committee): https://danishleadco.io/blog/how-to-accelerate-fintech-sales-with-buying-committees . Tier: blog.
[^a16z]: a16z (Haber et al.), "B2FI: Demystifying Software Sales Into Financial Institutions" (2023) (risk-averse culture; buyers can't evaluate startups): https://a16z.com/b2fi-demystifying-software-sales-into-financial-institutions/ . Tier: blog (VC).
[^ig-frb]: Federal Reserve, "Interagency Guidance on Third-Party Relationships" (SR 23-4 landing; life cycle; "not law" but examined): https://www.federalreserve.gov/frrs/guidance/interagency-guidance-on-third-party-relationships.htm ; PDF: https://www.federalreserve.gov/supervisionreg/srletters/sr2304a1.pdf . Tier: regulator.
[^ig-occ]: OCC Bulletin 2023-17, "Third-Party Relationships: Interagency Guidance on Risk Management" (2023): https://www.occ.gov/news-issuances/bulletins/2023/bulletin-2023-17.html . Tier: regulator.
[^ig-fdic]: FDIC FIL-29-2023 (2023): https://www.fdic.gov/news/financial-institution-letters/2023/fil23029.html . Tier: regulator.
[^hs-cfr]: 12 CFR Part 30, Appendix D (OCC Heightened Standards), Cornell LII: https://www.law.cornell.edu/cfr/text/12/appendix-D_to_part_30 . Tier: regulator/legal.
[^hs-occ]: OCC, "Corporate and Risk Governance," Comptroller's Handbook (written risk-governance framework; three lines of defense): https://www.occ.treas.gov/publications-and-resources/publications/comptrollers-handbook/files/corporate-risk-governance/pub-ch-corporate-risk.pdf . Tier: regulator.
[^hs-large]: OCC, "Large Bank Supervision," Comptroller's Handbook (≥$50B covered banks; new/unproven tech as strategic risk): https://www.occ.treas.gov/publications-and-resources/publications/comptrollers-handbook/files/large-bank-supervision/pub-ch-large-bank-supervision.pdf . Tier: regulator.
[^ffiec1]: FFIEC IT Handbook, "Outsourcing Technology Services — Introduction": https://ithandbook.ffiec.gov/it-booklets/outsourcing-technology-services/introduction.aspx . Tier: regulator.
[^ffiec2]: FFIEC IT Handbook, "Outsourcing — Risk Assessment and Requirements" ("same… as if conducted in-house"): https://ithandbook.ffiec.gov/it-booklets/outsourcing-technology-services/risk-management/risk-assessment-and-requirements . Tier: regulator.
[^ffiec3]: FFIEC IT Handbook, "Management III.C.8 Third-Party Management" (enterprise-wide governance issue): https://ithandbook.ffiec.gov/it-booklets/management/iii-it-risk-management/iiic-risk-mitigation/iiic8-third-party-management . Tier: regulator.
[^mr-fdic]: FDIC FIL-22-2017 adopting SR 11-7 (validation, monitoring, governance): https://www.fdic.gov/news/financial-institution-letters/2017/fil17022.pdf . Tier: regulator.
[^mr-occ2011]: OCC Bulletin 2011-12 / SR 11-7 attachment: http://business.cch.com/bfld/OCC-2011-12-04042011.pdf . Tier: regulator.
[^mr-occ2026]: OCC Bulletin 2026-13, "Model Risk Management: Revised Guidance" (2026; risk-based; covers vendor/third-party models): https://www.occ.gov/news-issuances/bulletins/2026/bulletin-2026-13.html . Tier: regulator (2026 revision — verify).
[^mr-frb2026]: FRB SR 26-2 attachment, revised "Supervisory Guidance on Model Risk Management" (2026): https://www.federalreserve.gov/supervisionreg/srletters/SR2602a1.pdf . Tier: regulator (2026 revision — verify).
[^dora]: EUR-Lex, Regulation (EU) 2022/2554 (DORA) — binding; applied 17 Jan 2025; no rigid caps; Lead Overseer/CTPP oversight: https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX%3A32022R2554 . Tier: regulator.
[^dora-ch5]: DORA Chapter V (managing ICT third-party risk; Arts. 28–29 strategy, Register, pre-contract assessment, exit strategies, concentration assessment): https://digital-operational-resilience.org/regulations/digital-operational-resilience-act/chapter-v/ . Tier: docs (regulation text).
[^dora-art30]: DORA Art. 30 (key contractual provisions; additional clauses for critical/important functions; exit assistance + transition period): https://www.aiact-info.eu/regulation/dora/article/30/key-contractual-provisions . Tier: docs (regulation text mirror).
[^dora-cyadviso]: CyAdviso, "DORA Third-Party ICT Risk: Art. 28–30 Guide for 2026" (corroborating summary): https://www.cyadviso.com/dora-third-party-ict-risk . Tier: blog.
[^dora-rts]: EC, DORA RTS Delegated Reg. 2024/1773 explanatory + Art. 10 exit plan (mandatory per critical arrangement; tested): https://ec.europa.eu/finance/docs/level-2-measures/dora-regulation-rts--2024-1531_en.pdf . Tier: regulator.
[^dora-ivass]: IVASS mirror, Delegated Reg. (EU) 2024/1773 (audit/access rights; pooled audits): https://www.ivass.it/normativa/internazionale/internazionale-ue/regolamenti-europei/rd-2024-1773/Delegated_Regulation_2024_1773_of_13_march_2024.pdf . Tier: regulator.
[^eba]: EBA Guidelines on Outsourcing (EBA/GL/2019/02), Google Cloud mapping (audit rights cascaded to subcontractors; data location): https://services.google.com/fh/files/misc/gcp_eba_outsourcing_guidelines_mapping.pdf . Tier: regulator mapping.
[^ecb]: ECB, "Guide on outsourcing cloud services to cloud service providers" (2025) (lock-in, concentration, data residency, exit as supervisory items): https://www.bankingsupervision.europa.eu/ecb/pub/pdf/ssm.supervisory_guides202507.en.pdf . Tier: regulator.
[^dora-cloudexit]: Regulation-DORA, "Cloud Exit Strategies & Concentration Risk Under DORA (2026)" ("an untested exit is not credible"): https://www.regulation-dora.eu/blog/cloud-exit-strategy-concentration-risk-dora . Tier: blog.
[^res-aws]: AWS, "HKMA Cloud Adoption Compliance Guide" (data residency; region control; attestation): https://d1.awsstatic.com/onedam/marketing-channels/website/public/HongKong-FinServ-ComplianceGuide-CloudAdoption.pdf . Tier: vendor.
[^ccgroup]: CCGroup (D. Lowther), "From promise to proof: how fintech buying is changing" (buyer survey: bigger/more deliberate; proof decisive): https://ccgrouppr.com/blog/from-promise-to-proof-how-fintech-buying-is-changing/ . Tier: blog.
[^neg-bain]: Bain & Co., "Procurement in Financial Services: A Tech-Powered Value Creator" (procurement earlier; innovation/speed/resilience): https://www.bain.com/insights/procurement-in-financial-services-a-tech-powered-value-creator/ . Tier: consulting.
[^neg-kpmg]: KPMG Global Tech Report 2026: Financial Services (innovator self-identification; open-source shift): https://assets.kpmg.com/content/dam/kpmgsites/mx/pdf/2026/04/global-tech-report-fs-2026.pdf . Tier: consulting.
[^neg-cloud]: Capgemini World Cloud Report – Financial Services (37%→91% cloud adoption 2020→2023): https://paymentsindustryintelligence.com/wp-content/uploads/2023/11/World-Cloud-Report-Financial-Services.pdf . Tier: industry report.
[^neg-endava]: Endava, "Beyond the RFP: How FS Firms Are Accelerating AI Adoption" (marketplace standard contracts compress legal cycles): https://www.endava.com/insights/articles/beyond-the-rfp . Tier: consulting.
[^neg-gfma]: GFMA/SIFMA, "Public Cloud Portability" white paper (multi-cloud may double cost / impractical; mainframe precedent): https://www.sifma.org/wp-content/uploads/2025/04/GFMA-Cloud-Portability-Paper-FINAL.pdf . Tier: industry.
[^neg-pifs]: PIFS-Harvard (H. Nadler), "Cloud Adoption in the Financial Sector and Concentration Risk" (multi-cloud "too technically complex"): https://www.pifsinternational.org/wp-content/uploads/2023/04/PIFS-Cloud-Adoption.pdf . Tier: academic/industry.
[^neg-ecipe]: ECIPE, "Cloud Resilience and Security: Why Exit, Portability, and Lifecycle Design Matter" (un-contestability, not concentration, is the core risk): https://ecipe.org/publications/cloud-resilience-and-security/ . Tier: think tank.
[^neg-ambanker]: American Banker, "Banking on a single cloud platform? It's time to rethink the risk" (banks overestimate migration ability; test exit): https://www.americanbanker.com/opinion/banking-on-a-single-cloud-platform . Tier: trade press.
