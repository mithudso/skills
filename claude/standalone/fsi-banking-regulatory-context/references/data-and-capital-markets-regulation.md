<!-- Provenance: reference under the `fsi-banking-regulatory-context` skill. Created 2026-06-18 via /dr deep research (gov/standards primary sources + Big-4/law-firm explainers). Educational, pointer-level seller orientation — NOT legal or compliance advice. -->

# Data-Centric & Capital-Markets Regulation (the data-platform angle)

`verified-as-of: 2026-06-18` (model-risk guidance and Reg SCI are fast-moving — re-verify any dated claim).

> **Educational orientation, not legal/compliance advice.** Each item: what it is, who enforces it, and the "why it matters to a DATA-PLATFORM buyer/seller" angle. Pointer-level only. **Do not duplicate MongoDB-specific PCI/encryption config** — that is `mongodb-compliance`.

## Contents

- [BCBS 239 — risk-data aggregation (the most database-relevant rule)](#1-bcbs-239--risk-data-aggregation--reporting)
- [PCI-DSS — card-data security](#2-pci-dss--payment-card-industry-data-security-standard)
- [SR 11-7 → SR 26-2 — model risk (pointer only)](#3-sr-11-7--sr-26-2--model-risk-management-pointer-only)
- [GDPR & CCPA/CPRA in an FSI context](#4-gdpr--ccpacpra-in-an-fsi-context)
- [SEC & FINRA for capital markets (the data/tech angle)](#5-sec--finra-for-capital-markets--the-datatech-angle)
- [Sources](#sources)

## 1. BCBS 239 — risk-data aggregation & reporting

**What it is.** A set of **14 principles** (11 for banks, 3 for supervisors) issued by the **Basel Committee (BCBS / BIS) in January 2013**, a direct response to the crisis finding that many banks, including G-SIBs, **could not aggregate their risk exposures quickly or accurately.**[^1][^2] BCBS sets the principles but has **no direct enforcing power**; enforcement runs through national/regional supervisors (most visibly the **ECB** in the euro area, plus OCC/Fed/FDIC in the US).[^3][^9]

**The 14 principles (summary):**[^1][^7]
- **Governance & infrastructure:** (1) **Governance** (board oversight, documentation, validation); (2) **Data architecture & IT infrastructure**: design/build/maintain infrastructure that supports aggregation/reporting *in normal times and during stress*.
- **Aggregation capabilities:** (3) **Accuracy & integrity** (golden sources, data dictionary, controlled/automated aggregation); (4) **Completeness** (all material risk data across the group); (5) **Timeliness** (esp. in stress); (6) **Adaptability** (on-demand/ad-hoc requests).
- **Reporting practices:** (7) accuracy; (8) comprehensiveness; (9) clarity/usefulness; (10) frequency; (11) distribution.
- **Supervisory:** (12) review; (13) remedial actions; (14) home/host cooperation.

**Why it's the most database-relevant banking standard.** Principle 2 makes **data architecture and IT infrastructure an explicit precondition** for everything else.[^1] The principles operationalize, as supervisory expectations, the exact concepts a data platform sells against: integrated **data taxonomies/metadata**; **single identifiers** for legal entities/counterparties/customers/accounts; end-to-end **data lineage** (origin → capture → transform → final use, at attribute level); **data-quality controls across the lifecycle**; **automation over manual workarounds**; and accurate/complete/timely reporting **on demand and under stress**. The ECB's RDARR Guide explicitly calls for an integrated, documented group-level data architecture, data dictionaries/metadata repositories, and "complete and up-to-date data lineages on data-attribute level."[^3]

**As of 2026 — banks still struggle (high confidence).**
- A **January 2026 BIS/BCBS newsletter** states "meeting the intended outcomes of the Principles remains a continuous effort," and that **data lineage "remains a challenging component."**[^14]
- Early/later BIS progress reports found architecture (P2), adaptability (P6), and accuracy (P3) the lowest-compliance principles, with one finding **only 2 of 31 assessed banks fully compliant.**[^10][^21]
- The **ECB's 2024 Annual Report** notes "long-standing deficiencies" in RDARR capabilities for over a decade; ECB officials say adequate capabilities "are still the exception" and have threatened enforcement and **capital add-ons.**[^9][^13]

> **Why it matters to a data-platform seller.** BCBS 239 is the regulatory tailwind for *every* data-governance, lineage, metadata, single-source-of-truth, data-quality, and stress-resilient-reporting conversation in a bank. A persistent, decade-long, supervisor-escalated compliance gap means budget keeps flowing here, and the buyer's pain (legacy silos, fragmented estates, manual reconciliation, un-traceable lineage) maps directly to a modern governed data architecture. *(Confidence: fact, 3+ primary sources, 2026-dated.)*

## 2. PCI-DSS — Payment Card Industry Data Security Standard

**What it is.** A **contractual security standard** (*not* a government law) maintained by the **PCI Security Standards Council** (founded by the card brands). Compliance is enforced through **acquiring-bank and payment-network contracts.**[^17][^24] It governs **cardholder data (CHD)** and **sensitive authentication data (SAD)**.

**Version status as of 2026.** **v4.0** published March 2022; **v4.0.1** published June 2024 (a limited revision: clarifications, **no new/deleted requirements**); the future-dated requirements that were best-practice-only through **March 31, 2025 are now mandatory**; **v4.0.1 is the active version** and v3.2.1 is retired.[^24][^25][^26]

**The 12 requirements (pointer-level), in 6 control groups:**[^24][^30] (1) network security controls/segmentation; (2) secure configurations (no vendor defaults); (3) protect stored account data — **render PAN unreadable (encryption/tokenization), mask on display, never store SAD post-auth**; (4) strong cryptography for CHD in transit; (5) anti-malware; (6) secure systems/software; (7) restrict access by **need-to-know**; (8) identify & authenticate (v4 expanded **MFA to all access into the CDE**); (9) restrict physical access; (10) **log & monitor all access (audit trails)**; (11) test security regularly (incl. segmentation pen-tests); (12) policies/programs (incl. annual documented scope reviews + third-party governance).[^27][^29]

**Scope (CDE) and tokenization: the cost lever.** The **Cardholder Data Environment** is *every system that stores/processes/transmits CHD/SAD plus any connected system that could impact it.*[^24][^29][^30] **Scope reduction is the single most cost-effective control:** network/identity **segmentation** keeps non-CDE systems out of scope; **tokenization / P2PE** replace CHD with tokens to shrink the footprint. Note the token **vault itself becomes a high-value, in-scope CDE asset** requiring isolation, key management, and detokenization controls.[^28][^31] A common costly misconception: "we use a processor, so we're out of scope." A hosted/redirect model still requires the relevant SAQ and won't survive a QSA review without segmentation/P2PE evidence.[^24][^30]

> **Negation / criticism (qualified).** A recurring industry critique (anchored in **Verizon Payment Security Reports**) is that PCI compliance is a **point-in-time checkbox, not security**: organizations pass interim assessments but fall out of compliance between annual audits, and Verizon reported that none of the breached entities it investigated was fully PCI-compliant at the time of breach.[^18][^32][^34] Treat as a well-known critique resting largely on one investigative dataset, not an independent multi-source fact.

> **Why it matters to a data-platform seller.** PCI-DSS dictates *how and where card data is stored and segmented*: the reason buyers tokenize, isolate a CDE, encrypt at rest/in transit, mask PAN, control access by need-to-know, and keep audit trails. **Descoping (storing fewer card bytes) is the dominant cost driver.** Frame value in those terms, and keep it industry-level; vendor-specific PCI config is `mongodb-compliance`. *(Confidence: fact on version/scope; qualified on the "compliant ≠ secure" critique.)*

## 3. SR 11-7 → SR 26-2 — model risk management (pointer only)

**What model risk management (MRM) is.** Supervisory **guidance** (not an enforceable rule) on managing **"model risk"**: adverse consequences when models (quantitative methods turning inputs into estimates) are wrong or misused. Two core concepts: **effective challenge** (critical analysis by objective, competent, independent parties with standing to force change) and **model validation** (conceptual soundness; ongoing monitoring; outcomes analysis/back-testing).[^35][^36]

**Major 2026 update, fast-moving (independently re-verified during authoring).** On **April 17, 2026** the Fed, OCC, and FDIC jointly issued **revised guidance, SR 26-2 / OCC Bulletin 2026-13**, that **supersedes and replaces SR 11-7** (and the 2021 BSA/AML model-risk statement, SR 21-8).[^38][^39][^40] What a seller should know:
- A **risk-based, tailored** approach most relevant to banking organizations with **over $30 billion in total assets** (may also apply to smaller orgs with significant model risk).[^39][^40]
- A **narrower "model" definition** that explicitly excludes spreadsheet arithmetic and deterministic rule-based processes.[^39][^40]
- **Generative AI and agentic AI are explicitly out of scope** (Footnote 3) as "novel and rapidly evolving"; the guidance **still applies** to traditional statistical/quantitative models and **non-generative, non-agentic ML models.** The agencies plan a forthcoming **Request for Information** on AI/GenAI model risk.[^38][^39][^40]
- It is **guidance, not enforceable standards**; the OCC states non-compliance "will not result in supervisory criticism."[^39]
- Analyst-flagged consequence: a **governance gap**, where a GenAI/agentic layer sits *outside* MRM while the traditional components it orchestrates sit *inside*, so banks must govern GenAI via *broader* risk-management/third-party/operational frameworks. "Not a free pass": it still must be governed somewhere.[^40]

> **Why it matters (keep brief).** MRM governs *models, not databases*. The relevant hooks: models depend on **high-quality, well-governed input data** (BCBS 239 territory), a **model inventory**, **documentation/lineage**, and **vendor/third-party model oversight**. The 2026 AI/ML model-governance area is **fast-moving and currently in a regulatory gap**, useful context when a data platform underpins AI/ML workloads but not a primary data-platform compliance driver. *(Confidence: fact on the April 2026 replacement + GenAI carve-out; the governance-gap reading is analyst consensus — evolving.)*

## 4. GDPR & CCPA/CPRA in an FSI context

**GDPR (EU/UK).** Governs processing of EU residents' personal data, enforced by national DPAs coordinated via the EDPB. FSI-specific tensions:
- **Right to erasure (Art. 17) vs. financial recordkeeping (high confidence).** A deletion request appears to collide with multi-year retention duties (AML/KYC, MiFID II, etc.). The resolution is **Article 17(3)**: erasure **does not apply where processing is necessary to comply with a legal obligation** under EU/member-state law, or to establish/defend legal claims, so a firm can lawfully refuse erasure for data it is legally required to retain (e.g., AML's typical 5-years-from-relationship-end).[^45][^46][^50] Nuance: the legal-obligation basis is **limited to EU/member-state law**, so a US obligation (e.g., 17a-4) does not automatically qualify, a genuine cross-border conflict for global firms.[^48]
- **Architectural pattern this drives:** **partition/segregate data by retention basis** (AML-retained vs. operational vs. preference), automated deletion of non-retained data at account closure, segregated archives for legally retained data, and **anonymization** (which GDPR permits to keep indefinitely) as an alternative to deletion.[^49][^50] Directly a data-platform requirement: classify by retention basis, enforce differential lifecycle/TTL, prove proportionality.
- **Cross-border transfer & residency/localization (high confidence).** Beyond GDPR's own mechanisms (adequacy, SCCs), **sector-specific localization** forces FSI architecture: **China** (CSL/DSL/PIPL + PBOC, where core systems and personal financial info are generally stored/operated in China and cross-border transfer needs CAC assessment);[^51][^52] **India** (RBI 2018 *Storage of Payment System Data*, where payment data is stored only in India, cross-border-transaction copies repatriated/deleted within 24h, and the DPDP Act 2023 adds outbound oversight);[^53][^54] **Russia, Vietnam, Indonesia** and others. The US default favors cross-border flows with contractual safeguards.[^51][^54]

**CCPA / CPRA (California): the FSI nuance that surprises people (high confidence).** The CCPA (as amended by CPRA, effective Jan 1, 2023; enforced by the **California Privacy Protection Agency** + AG) does **NOT give financial institutions a blanket entity-level exemption**, only a **data-level exemption**: information collected/processed/sold/disclosed **pursuant to GLBA** is carved out, **but the institution itself is not.**[^56][^59] This contrasts with most other state privacy laws (Virginia/Colorado/Utah/Connecticut), which give a broader *entity-level* GLBA exemption (though some are narrowing).[^58][^60] Practical effect: a financial institution must determine, **data-element-by-element**, whether each piece arose from a financial product/service (GLBA-exempt) or from other sources (marketing/website/employee/B2B data, under CCPA). CCPA's "personal information" is far broader than GLBA's "NPI," so the same element can be exempt in one context and in-scope in another. The carve-out does **not** cover CCPA's data-breach private right of action.[^56][^58][^61]

> **Why it matters to a data-platform seller.** Privacy regimes turn into concrete requirements: **data classification/tagging by source and legal basis, differential retention/TTL and automated deletion, segregated/anonymized archives, region-pinned storage and residency, and dataset-level cataloging** of exempt vs. in-scope. The erasure-vs-retention tension and localization mandates are *architecture decisions*. *(Confidence: fact, multiple primary + law-firm sources.)*

## 5. SEC & FINRA for capital markets — the data/tech angle

*Scope note: only the data/technology obligations, not trading strategy or investor-facing securities law (→ `trading-and-investing`).*

- **Books-and-records / WORM (high confidence).** **SEC Rule 17a-4** requires broker-dealers to **preserve records for required periods in a tamper-resistant electronic form**, historically **WORM (Write Once, Read Many)**. **2022-23 amendments** (compliance May 3, 2023) added an **"audit-trail alternative"**: instead of WORM, a system may keep a **complete, time-stamped audit trail permitting recreation of an original record if modified/deleted** (capturing all modifications/deletions, timestamps, actor identity). WORM remains valid; systems must support prompt download in human-readable + usable electronic formats plus a redundant backup.[^63][^64][^66] **FINRA Rule 4511** requires members to make/preserve books and records **in a format/media that complies with SEA 17a-4**, with a **default 6-year retention.**[^69][^70]
- **Off-channel-communications enforcement wave (2021-2025), note the magnitude (high confidence).** A sprawling SEC (+ parallel CFTC) sweep penalized firms for failing to preserve business communications over **personal devices and unapproved apps (WhatsApp, iMessage, Signal, personal text/email)**, violating 17a-4 / Advisers Act 204-2. It began with **JPMorgan's $200M settlement (late 2021)** and ran through multiple waves (e.g., $1.1B+ across 16 firms in 2022; $392.75M across 26 firms in Aug 2024; $63M across 12 firms in Jan 2025), cumulatively **roughly $2-3 billion across ~60+ firms.**[^71][^72][^73][^74][^76]
- **Reg SCI — Systems Compliance and Integrity (medium-high confidence).** Adopted **2014**; governs the **technology resilience of critical US market infrastructure** ("SCI entities": exchanges, clearing agencies, large ATSs, plan processors) and their "SCI systems." Mandates capacity, integrity, resiliency, availability, security; written system inventories/lifecycle management; **BC/DR testing; annual SCI reviews; immediate SCI-event notification to the SEC.** A **2023 proposal would expand** the "SCI entity" definition (large broker-dealers, SBSDRs, certain exempt clearing agencies) and add cloud/third-party, cybersecurity, and annual-pen-test provisions.[^78][^79][^80] *As of 2026, treat the expansion as proposed/pending; verify final status.*
- **Consolidated Audit Trail (CAT) (high confidence).** Created under **SEC Rule 613 (Reg NMS)**: a **single central repository** capturing **the complete lifecycle of every order, quote, cancellation, modification, and execution in NMS securities and listed options across all US markets.** SROs and their member broker-dealers must record/report reportable events (next-trading-day), with millisecond+ clock sync and a **"daisy-chain" linkage** so an order's full lifecycle can be reconstructed. **No exemptions by firm size or trading type.**[^83][^84][^86][^87]

> **Why these drive data tech.** This cluster is the engine behind **immutable / append-only / WORM-or-audit-trail storage, long-horizon retention (6+ years), full audit trails, communications surveillance/capture, time-synchronized event logging, lineage/reconstruction of order lifecycles, and high-availability/DR resilience.** 17a-4 + FINRA 4511 = immutable, retrievable, retention-governed storage; the off-channel wave = a multi-billion-dollar cost of getting data retention wrong; Reg SCI = resilience/capacity/BC-DR/security for market-critical systems; CAT = massive, time-sequenced, linkable, queryable event data at market scale. *(Confidence: fact on 17a-4/4511/CAT/enforcement; Reg SCI expansion is proposed, qualify.)*

## Sources

[^1]: https://www.bis.org/publ/bcbs239.pdf — BCBS 239 full text (gov/standards; primary).
[^2]: https://www.bis.org/publ/bcbs239.htm — BCBS 239 landing page (gov/standards).
[^3]: https://www.bankingsupervision.europa.eu/ecb/pub/pdf/ssm.supervisory_guides240503_riskreporting.en.pdf — ECB Guide on effective RDARR (gov; lineage/architecture detail).
[^7]: https://www.pwc.com.au/financial-services/pdf/principles-for-effective-risk-data-aggregation-sept19.pdf — PwC: BCBS 239 essentials (analyst; 14-principle table).
[^9]: https://www.bankingsupervision.europa.eu/press/publications/annual-report/html/index.en.html — ECB Annual Report on supervisory activities 2024 (gov; "long-standing deficiencies").
[^10]: https://www.bis.org/bcbs/publ/d308.pdf — BIS: Progress in adopting the Principles, 2015 (gov/standards; lowest-compliance principles).
[^13]: https://kpmg.com/xx/en/our-insights/ecb-office.html — KPMG ECB Office: BCBS 239 / SREP commentary (analyst; ECB enforcement posture).
[^14]: https://www.bis.org/publ/bcbs_nl36.htm — BIS/BCBS Newsletter on BCBS 239 implementation, Jan 2026 (gov/standards; lineage still hard).
[^21]: https://www.bis.org/bcbs/publ/d559.pdf — BIS: later progress report (gov/standards; "only 2 of 31 fully compliant").
[^17]: https://www.pcisecuritystandards.org/document_library/ — PCI SSC Document Library (standards; primary).
[^18]: https://www.verizon.com/business/resources/reports/payment-security-report/ — Verizon Payment Security Report (analyst/trade; "checkbox" critique).
[^24]: https://www.pcisecuritystandards.org/ — PCI SSC + v4.0.1 practitioner guides (standards/trade; version/scope/tokenization).
[^25]: https://blog.pcisecuritystandards.org/just-published-pci-dss-v4-0-1 — PCI SSC: PCI DSS v4.0.1 published (standards; "no new/deleted requirements").
[^26]: https://www.securitymetrics.com/learn — SecurityMetrics: what changed with v4.0.1 (trade-press; future-dated reqs now live).
[^27]: https://www.securitymetrics.com/learn — v4.0.1 new requirements (MFA, scope reviews) (trade-press).
[^28]: https://datastealth.io/ — PCI DSS v4.0 tokenization requirements (trade-press; token-vault-as-CDE-asset).
[^29]: https://www.pcisecuritystandards.org/ — PCI DSS v4.0.1 compliance checklist (standards/trade; scope/segmentation/CDE).
[^30]: https://www.pcisecuritystandards.org/ — PCI DSS overview (standards; 12 requirements, CDE, scope minimization).
[^31]: https://datastealth.io/ — tokenization vault security controls (trade-press).
[^32]: https://csiac.org/articles/compliant-but-not-secure-why-pci-certified-companies-are-being-breached/ — Compliant but not Secure (analyst; Target/Heartland).
[^34]: https://www.securityweek.com/pci-dss-compliance-between-audits-declining-verizon/ — PCI compliance between audits declining (trade-press; Verizon sustainability data).
[^35]: https://www.federalreserve.gov/supervisionreg/srletters/SR2602.htm — Fed SR 26-2: Revised Guidance on Model Risk Management, Apr 17 2026 (gov primary; supersedes SR 11-7). **Independently re-verified during authoring.**
[^36]: https://occ.gov/news-issuances/bulletins/2011/bulletin-2011-12a.pdf — OCC Bulletin 2011-12a (original SR 11-7 text) (gov; effective challenge/validation).
[^38]: https://www.federalreserve.gov/supervisionreg/srletters/SR2602a1.pdf — SR 26-2 attachment: Supervisory Guidance on MRM, Apr 17 2026 (gov primary; GenAI carve-out Footnote 3). **Independently re-verified.**
[^39]: https://www.occ.gov/news-issuances/bulletins/2026/bulletin-2026-13.html + https://www.occ.gov/news-issuances/news-releases/2026/nr-occ-2026-29.html — OCC Bulletin 2026-13 + news release (gov primary; $30B relevance, "not enforceable," RFI forthcoming). **Independently re-verified.**
[^40]: https://www.sullcrom.com/insights/memo/2026/April/OCC-Fed-FDIC-Issue-Revised-Guidance-Model-Risk-Management — Sullivan & Cromwell memo (law-firm; change summary + governance-gap analysis).
[^45]: https://www.pwc.ch/en/insights/fs/mifidII-and-gdpr.html — PwC: MiFID II and GDPR (analyst; erasure vs retention).
[^46]: https://www.steel-eye.com/news/gdpr-vs-mifid-ii-record-keeping-vs-personal-data-protection — SteelEye: GDPR vs MiFID II (trade-press; Art. 17(3)).
[^48]: https://ir.lawnet.fordham.edu/jcfl/ — Fordham: reconciling US data preservation with GDPR erasure (academic; US-obligation ≠ GDPR basis).
[^49]: https://www.lexia.it/en/ — LEXIA: document retention and GDPR, 2026 (law-firm; segregation/anonymization).
[^50]: https://gdpr-info.eu/art-17-gdpr/ — GDPR Article 17 text (primary; erasure + 17(3) exemptions).
[^51]: https://assets.publishing.service.gov.uk/ — Frontier Economics: extent & impact of data localisation (paper/gov; cross-jurisdiction map).
[^52]: https://www.cliffordchance.com/ — Clifford Chance: PBOC data & cyber measures (law-firm; China FSI residency).
[^53]: https://www.rbi.org.in/Commonperson/english/scripts/FAQs.aspx?Id=2995 — RBI FAQ: Storage of Payment System Data (gov primary; India 24h repatriation).
[^54]: http://www.nishithdesai.com/ — Nishith Desai: cross-border compliance in fintech (law-firm; India vs US contrast).
[^56]: https://oag.ca.gov/privacy/ccpa — California AG: CCPA overview (gov primary).
[^58]: https://www.dwt.com/blogs/ — Davis Wright Tremaine: CCPA exempt-or-not (law-firm; data-element analysis).
[^59]: https://california-ccpa.org/section-1798-145-exemptions/ — CCPA §1798.145 exemptions (primary statute; GLBA data-level carve-out).
[^60]: https://www.nortonrosefulbright.com/ — Norton Rose Fulbright: GLBA vs state privacy laws (law-firm; entity vs data exemption).
[^61]: https://www.debevoise.com/ — Debevoise: CCPA compliance for FIs (law-firm; exemption-scope nuance).
[^63]: https://www.law.cornell.edu/cfr/text/17/240.17a-4 — 17 CFR §240.17a-4 (primary regulation text).
[^64]: https://www.sec.gov/investment/amendments-electronic-recordkeeping-requirements-broker-dealers — SEC: 17a-4 electronic-recordkeeping amendments (gov; WORM + audit-trail alternative).
[^66]: https://www.sec.gov/files/rules/final/2022/34-96034.pdf — SEC Final Rule 34-96034, 2022 (gov primary; amendments).
[^69]: https://www.finra.org/rules-guidance/rulebooks/finra-rules/4511 — FINRA Rule 4511 text (gov primary; 6-year default, media compliance).
[^70]: https://www.finra.org/rules-guidance/key-topics/books-records — FINRA: Books and Records (gov; retention summary).
[^71]: https://www.sec.gov/newsroom/press-releases/2024-98 — SEC: 26 firms pay $392.75M (gov primary).
[^72]: https://www.sec.gov/newsroom/press-releases/2024-18 — SEC: 16 firms pay $81M (gov primary).
[^73]: https://www.sec.gov/ — SEC off-channel enforcement releases (gov; cumulative ~$2B+). Corroborated by FT coverage.
[^74]: https://www.sec.gov/newsroom/press-releases/2025-6 — SEC: 12 firms pay $63.1M, Jan 2025 (gov primary).
[^76]: https://www.legaldive.com/ — Legal Dive: SEC off-channel cumulative tally (~$3B / ~60 firms) (trade-press).
[^78]: https://www.sec.gov/files/34-97143-fact-sheet.pdf — SEC: Reg SCI proposed expansion fact sheet (gov; six functions, cloud/3rd-party).
[^79]: https://www.sec.gov/files/rules/proposed/2023/34-97143.pdf — SEC: Reg SCI proposed rule, 2023 (gov primary).
[^80]: https://www.sec.gov/newsroom/press-releases/2023-53 — SEC: proposes to expand Reg SCI (gov; system inventory/lifecycle).
[^83]: https://catnmsplan.com/ — CAT NMS Plan official site (gov/SRO; live tech specs).
[^84]: https://www.sec.gov/rules-regulations/rule-613-consolidated-audit-trail — SEC: Rule 613 / CAT (gov primary; data scope, next-day reporting).
[^86]: https://www.sec.gov/files/rules/sro/finra/2017/34-79961.pdf — SEC: FINRA Rule 6800 Series filing (gov; daisy-chain linkage).
[^87]: https://www.finra.org/rules-guidance/key-topics/consolidated-audit-trail-cat — FINRA: CAT key topic (gov; "no exemptions by firm size").
