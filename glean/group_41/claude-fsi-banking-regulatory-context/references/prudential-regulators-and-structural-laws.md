<!-- Provenance: reference under the `fsi-banking-regulatory-context` skill. Created 2026-06-18 via /dr deep research (gov primary sources + Big-4/law-firm explainers). Educational, pointer-level seller orientation — NOT legal or compliance advice; gives no compliance instructions. -->

# Prudential Regulators & Structural Banking Laws (pointer-level)

`verified-as-of: 2026-06-18` (the 2025-2026 US capital + supervision landscape is actively changing — re-verify any dated claim before relying on it).

> **Educational orientation, not legal/compliance advice.** Each item below is "what it is, who enforces it, the one-line reason it shapes technology." It is *not* a compliance checklist and makes no representation about any institution's obligations. Route specifics to the customer's own legal/risk function.

## Contents

- [The regulators and who supervises whom](#the-regulators-and-who-supervises-whom)
- [The FFIEC IT Examination Handbook](#the-ffiec-it-examination-handbook--the-examiners-tech-rulebook)
- [The examination culture: MRAs / MRIAs](#the-examination-culture-mras--mrias)
- [Basel III / Basel IV (capital & liquidity)](#1-basel-iii--basel-iv-capital--liquidity-adequacy)
- [Dodd-Frank Act](#2-dodd-frank-act-2010)
- [Sarbanes-Oxley (SOX)](#3-sarbanes-oxley-act-sox-2002)
- [Gramm-Leach-Bliley Act (GLBA)](#4-gramm-leach-bliley-act-glba-1999)
- [The seller's mental model](#how-the-pieces-fit-the-sellers-mental-model)
- [Sources](#sources)

## The regulators and who supervises whom

US banking supervision is **fragmented across multiple federal agencies plus state regulators**; a bank's primary federal regulator depends on its charter type and Fed membership. *The same software at two banks can face different supervisory expectations depending on the charter.*[^1][^2][^13]

| Regulator | What it is / remit | Why it matters to a tech buyer/seller |
|---|---|---|
| **OCC** (Office of the Comptroller of the Currency) | Independent bureau of the **U.S. Treasury**; charters and supervises **national banks, federal savings associations, and federal branches of foreign banks**. Funded by assessments on banks.[^1][^2] | Has a dedicated IT/security examination function; its examiners scrutinize bank technology and third-party arrangements. Selling to a national bank = selling into an OCC-supervised control environment.[^1][^2][^11] |
| **Federal Reserve** (Board + 12 Reserve Banks) | The central bank. The Board is responsible for supervision; the 12 Reserve Banks examine. Supervises **bank/financial holding companies and state-chartered member banks**, with a systemic lens.[^3] | Supervises the *holding-company* level (where enterprise IT/data/risk decisions often sit) and runs stress tests — both major drivers of risk-data and modeling technology.[^3][^8] |
| **FDIC** (Federal Deposit Insurance Corporation) | Provides **deposit insurance** (at least **$250,000** per depositor, per ownership category, per insured bank) and is the primary federal supervisor of **state-chartered non-member banks**; also the failed-bank resolution authority.[^4][^13] | As insurer/resolver, cares intensely about operational resilience, recoverability, and data integrity — pulling through to backup/DR, records, and continuity technology.[^4] |
| **NCUA** (National Credit Union Administration) | Charters/supervises **federal credit unions** and runs the **NCUSIF** share-insurance fund — the credit-union analog of the FDIC.[^14] | Credit unions are a distinct, NCUA-supervised buyer segment with their own supervisory priorities; same safety/soundness + insurance logic.[^14] |
| **CFPB** (Consumer Financial Protection Bureau) | Created by Dodd-Frank for **consumer financial protection** (vs. prudential safety-and-soundness).[^8][^14] | *As of 2026 its posture is sharply contested* — substantially curtailed in early 2025 with litigation following, and signaling a narrowed supervision cycle. **Tentative / volatile — verify before relying.**[^14] |
| **FFIEC** (Federal Financial Institutions Examination Council) | Not a supervisor — an **interagency coordinating body** promoting uniform examination. Members: FRB, CFPB, FDIC, NCUA, OCC, and the State Liaison Committee.[^11][^13] | Publishes the **FFIEC IT Examination Handbook**, the de-facto reference examiners use to evaluate bank technology (below).[^11][^12] |

## The FFIEC IT Examination Handbook — the examiner's tech rulebook

The **IT Examination Handbook ("InfoBase")** is a series of booklets prepared *for examiners* to assess an institution's IT risk management. Booklets state they "do not impose new requirements"; they describe principles/practices examiners review, and each agency applies them under its own authority.[^11][^12] Two recent booklets show the direct technology hook:

- **Architecture, Infrastructure, and Operations (AIO)**, issued **June 30, 2021**: explicitly addresses **cloud computing, microservices, AI/ML, zero-trust architecture, and IoT**, and its scope covers depository institutions, nonbanks, holding companies, **and third-party service providers.**[^11][^12] *A software/cloud vendor's product is in-scope of how examiners assess the bank's architecture, and the vendor can itself be examined as a third-party service provider.*
- **Development, Acquisition, and Maintenance**, issued **September 2024**: covers SDLC, IT project management, and **supply-chain risk for systems/components.**[^11][^12]

> **Why it matters to tech.** The IT Examination Handbook is the single most useful public document for a seller to understand *what examiners expect* of bank technology. *(Confidence: high — FFIEC/OCC/FDIC primary sources.)*

## The examination culture: MRAs / MRIAs

US banks live under **continuous supervisory examination** by their primary regulator. When examiners find a deficiency they escalate via **supervisory findings**:[^9][^10]

- **MRA ("Matter Requiring Attention"):** the standard vehicle communicating criticisms/required corrective actions to a bank's management and board.[^9][^10]
- **MRIA ("Matter Requiring Immediate Attention"):** a **more urgent** variant (notably at the Fed) for significant safety-and-soundness risk, significant noncompliance, or **repeat criticisms** that escalated due to inaction. Some agencies also use **MRBA** (Matters Requiring Board Attention).[^9][^10]

Nuances for the *culture*, not the mechanics: MRAs/MRIAs are largely an **informal supervisory convention** (not statutory) that grew up in the examination process, and they are **confidential supervisory information** banks generally cannot disclose publicly.[^9][^10]

> **Why it matters to tech.** This examination-and-findings culture is *the* reason banks demand auditability, evidence, and demonstrable controls from technology and vendors; an open MRA on a system or a vendor-management gap creates real urgency and budget. **Note (late 2025-2026):** the Fed and other agencies have been overhauling supervision and issuing new "supervisory operating principles," with industry critique of inconsistent MRA practices. The framework holds; specific practices are in flux. *(Qualify / verify.)*[^9]

## 1. Basel III / Basel IV (capital & liquidity adequacy)

**What it is.** The **Basel framework** is a set of *international* banking standards set by the **Basel Committee on Banking Supervision (BCBS)**, hosted at the **Bank for International Settlements (BIS)** and endorsed by the G20. The standards are **not law in themselves**; each jurisdiction implements them via its own rules. **Basel III** (post-2008) strengthened both capital and liquidity: more/higher-quality capital against risk-weighted assets; the **Liquidity Coverage Ratio (LCR)** (survive a 30-day stress) and **Net Stable Funding Ratio (NSFR)** (stable funding over one year).[^5][^7] The **"Basel III endgame"** (loosely "Basel IV") is the **final batch of the 2017 Basel III reforms**, overhauling how credit, operational, and market risk are measured.[^5][^7]

**Why it matters to tech.** Computing capital/liquidity ratios, risk-weighting assets, running stress scenarios, and demonstrating the numbers to examiners is a **massive data-aggregation and computation problem**, a primary driver of demand for risk-data platforms, lakehouses, calculation engines, and lineage tooling. *(The data-quality standard for this, BCBS 239, is in `data-and-capital-markets-regulation.md`.)*

> **As of 2026, US implementation is in flux (verify any number).** Independently confirmed against Federal Reserve / OCC primary sources during authoring: on **March 19, 2026** the Fed, OCC, and FDIC jointly issued **three proposals** to modernize the capital framework (the Basel III endgame re-proposal plus G-SIB-surcharge and broader-capital proposals). The endgame proposal would replace the current **dual-calculation** requirement for the largest (Category I/II) banks with a **single "expanded risk-based approach,"** eliminating the advanced (internal-models) approach for credit risk; the agencies state aggregate capital would **"modestly decrease,"** while remaining well above pre-crisis levels. Comments were due **June 18, 2026.**[^15][^16][^17] This is a substantial softening from the **2023** proposal (which would have raised aggregate capital for the largest banks materially and failed to reach consensus).[^5] **Direction is fact; precise final calibration and finalization status are tentative; re-verify.**

## 2. Dodd-Frank Act (2010)

**What it did (pointer-level).** The post-2008 systemic-risk reform. Key structural pieces:[^8]
- **Stress testing / CCAR:** supervisory stress tests (DFAST) + the Fed's **Comprehensive Capital Analysis and Review** of large banks' capital adequacy/planning.
- **Volcker Rule:** restricts proprietary trading (partially rolled back in 2020).
- **Created the CFPB and the FSOC** (Financial Stability Oversight Council, chaired by the Treasury Secretary).

**Why it matters to tech.** Stress testing/CCAR is, like Basel, a **data-and-modeling engine**, driving scenario analytics, model-risk infrastructure, and on-demand enterprise risk-data assembly. *(Model-governance specifics → `data-and-capital-markets-regulation.md`.)*

> **Disconfirming note (qualify).** Dodd-Frank was **not repealed or "gutted."** The most material change was the **2018 EGRRCPA (S.2155)**, which raised the enhanced-prudential-standards ("SIFI") threshold from **$50B to $250B** in assets (with Fed discretion for $100B-$250B banks). Further deregulation efforts continue into 2025-2026, but broad statutory repeal faces political constraints. **Verify current status.**[^8]

## 3. Sarbanes-Oxley Act (SOX, 2002)

**What it requires.** Enacted after Enron/WorldCom; governs **public-company financial-reporting integrity.** The load-bearing provision for technology is **Section 404**, requiring management (and, for larger "accelerated filers," an external auditor under 404(b)) to establish, maintain, assess, and attest to the effectiveness of **Internal Controls Over Financial Reporting (ICFR).** Section 302 assigns executive responsibility.[^15s]

**Why it matters to tech (the big one for IT governance).** Because financial reports run on IT systems, SOX 404 pulls **IT General Controls (ITGCs)** into scope. It is the prime structural reason banks (and all public companies) demand from systems and vendors:[^15s]
- **Auditability / evidence:** controls an independent auditor can test and document.
- **Change management:** controlled, logged, approved changes to systems touching financial data.
- **Segregation of duties (SoD):** no single person controls multiple critical financial tasks; access controls that enforce it.
- Access management and DR/continuity for financial systems.

> **Disconfirming note (qualify).** SOX 404, especially 404(b) external audit, is persistently criticized as costly, *disproportionately* burdening smaller companies (a 2025 GAO report reconfirms the asymmetry). The requirements remain in force.[^15s]

## 4. Gramm-Leach-Bliley Act (GLBA, 1999)

**What it requires.** GLBA (which also repealed Glass-Steagall's separation of banking/securities/insurance) protects consumers' **nonpublic personal information (NPI)** via two rules:[^16g]
- **Safeguards Rule:** maintain an information-security program with administrative, technical, and physical safeguards, **and ensure affiliates and service providers protect customer information** in their care (FTC's updated Safeguards Rule, effective 2023, adds specific elements).
- **Privacy Rule:** disclose information-sharing practices and offer an **opt-out** of certain sharing.

Enforcement splits: the **FTC** for many non-bank financial institutions; the **federal banking agencies** for the institutions they supervise.[^16g]

**Why it matters to tech.** GLBA is a foundational driver of **data-protection, encryption, access-control, and vendor-security** requirements, and its explicit **service-provider** obligation is exactly why bank customers push security/data-handling requirements down onto their technology vendors. *(GDPR/PCI-DSS → `data-and-capital-markets-regulation.md`; this is the US prudential-adjacent data-protection baseline.)*

## How the pieces fit (the seller's mental model)

1. **Charter determines the supervisor** (OCC / Fed / FDIC / NCUA / state) → who examines using **FFIEC** common standards + the **IT Examination Handbook.**
2. **Examination culture** (MRAs/MRIAs) makes **auditability, evidence, and demonstrable controls** non-negotiable.
3. **The structural laws create the demand:** Basel III + Dodd-Frank stress testing → **risk-data & computation** platforms; SOX 404 → **auditability, change-management, segregation-of-duties**; GLBA → **data protection & vendor security.**
4. **Vendors are in-scope:** FFIEC AIO and GLBA both reach **third-party service providers**, so your product is, in effect, part of the bank's examined control environment.

## Sources

[^1]: https://www.occ.gov/about/who-we-are/index-who-we-are.html — Who We Are, OCC (gov; remit).
[^2]: https://www.occ.gov/about/what-we-do/index-what-we-do.html — What We Do, OCC (gov; supervisory mission).
[^3]: https://www.federalreserve.gov/aboutthefed/fedexplained/supervision-regulation.htm — The Fed Explained: Supervision & Regulation (gov).
[^4]: https://www.fdic.gov/resources/deposit-insurance/understanding-deposit-insurance — Understanding Deposit Insurance, FDIC (gov; $250k limit, FDIC role).
[^5]: https://www.federalreserve.gov/newsevents/speech/bowman20260312a.htm — Bowman speech on Basel III / bank capital, Mar 12 2026 (gov primary; capital-framework direction). Corroborated by analyst pieces (Arnold & Porter advisory, Mar 2026).
[^7]: https://www.bis.org/bcbs/ — BIS / Basel Committee (gov/standards; LCR & NSFR definitions).
[^8]: https://www.federalreserve.gov/publications/comprehensive-capital-analysis-and-review-questions-and-anwers.htm + https://www.brookings.edu/articles/no-dodd-frank-was-neither-repealed-nor-gutted-heres-what-really-happened/ + https://www.congress.gov/bill/115th-congress/senate-bill/2155 — Fed CCAR Q&A (gov); Brookings on Dodd-Frank (analyst); Congress.gov S.2155 (gov; $50B→$250B threshold).
[^9]: https://www.federalreserve.gov/frrs/guidance/supervisory-considerations-for-the-communication-of-supervisory-findings.htm + https://bpi.com/the-3-letters-at-the-heart-of-bank-supervision-dysfunction/ — Fed guidance on communicating supervisory findings (gov); Bank Policy Institute on MRA/MRIA (industry; mechanics + critique).
[^10]: https://www.gao.gov/assets/gao-19-352.pdf — GAO: Bank Supervision (gov watchdog; corroborates MRA/MRIA definitions).
[^11]: https://ithandbook.ffiec.gov/ — FFIEC IT Examination Handbook InfoBase (gov; member-agency list, Handbook purpose, AIO scope).
[^12]: https://www.occ.gov/news-issuances/bulletins/2024/bulletin-2024-26.html + https://www.occ.gov/news-issuances/bulletins/2021/bulletin-2021-30.html — OCC bulletins on the FFIEC IT Handbook booklets (gov; 2024 Dev/Acq/Maint, 2021 AIO).
[^13]: https://cable.tech/resources/a-guide-to-understanding-u-s-banking-regulators — Guide to US Banking Regulators (trade-press; corroborates the charter/regulator split — primary sources lead).
[^14]: https://ncua.gov/ + https://www.consumerfinance.gov/about-us/the-bureau/ — NCUA (gov); CFPB "The Bureau" (gov). The contested 2025-2026 CFPB posture is flagged tentative (corroborated by trade-press, e.g. PYMNTS).
[^15]: https://www.federalreserve.gov/newsevents/pressreleases/bcreg20260319a.htm — Federal Reserve: agencies request comment on three capital proposals, Mar 19 2026 (gov primary; single-calculation expanded risk-based approach, "modestly decrease," comments due Jun 18 2026). **Independently re-verified during authoring.**
[^16]: https://www.occ.gov/news-issuances/bulletins/2026/bulletin-2026-9.html — OCC Bulletin 2026-9: Regulatory Capital re-proposal (gov primary; Category I/II single framework). **Independently re-verified during authoring.**
[^17]: https://www.federalreserve.gov/newsevents/pressreleases/bowman-statement-20260319.htm — Bowman statement on the three capital proposals, Mar 19 2026 (gov primary).
[^15s]: https://files.gao.gov/reports/GAO-25-107500/index.html + https://corpgov.law.harvard.edu/2022/09/22/sarbanes-oxley-%C2%A7-404-at-twenty/ — GAO-25-107500 on SOX compliance costs (gov); Harvard corpgov "SOX §404 at Twenty" (academic). SOX 404 ICFR/ITGC/SoD linkage corroborated across multiple Big-4 explainers.
[^16g]: https://www.ftc.gov/business-guidance/privacy-security/gramm-leach-bliley-act + https://www.ftc.gov/legal-library/browse/rules/safeguards-rule — FTC: GLBA; FTC: Safeguards Rule (gov; Safeguards + Privacy Rules, service-provider obligation).
