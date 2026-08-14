---
name: medicaid-managed-care-behavioral-health
description: >-
  Medicaid managed care for behavioral health services: behavioral health
  carve-outs vs integrated managed care carve-ins, Medicaid managed care
  organization (MCO) requirements for BH, mental health and SUD service
  availability, parity compliance in Medicaid managed care, state plan vs
  waiver authority for managed care, network adequacy for BH in MCOs, care
  coordination and integration requirements, encounter data and quality
  measurement (HEDIS BH measures), capitation rate-setting for BH, and
  behavioral health carve-out managed care organizations (BHOs). Educational
  only, NOT medical/insurance advice; verified-as-of 2026-06-22.
  TRIGGER: Medicaid managed care behavioral health, Medicaid BH carve-out,
  Medicaid behavioral health carve-in, Medicaid MCO mental health coverage,
  Medicaid SUD treatment managed care, Medicaid behavioral health organization
  BHO, parity Medicaid managed care, Medicaid MCO network adequacy behavioral
  health, Medicaid behavioral health capitation, Medicaid MCO care coordination
  BH, HEDIS behavioral health Medicaid, Medicaid managed behavioral health
  integration, Medicaid MCO BH service requirements.
  SKIP: Medicaid FFS behavioral health state plan → medicaid-rehabilitative-services;
  MHPAEA parity details → mhpaea-prior-auth-discharge; Medicaid LTSS managed
  care → medicaid-ltss-hcbs-rehab; NC specific managed care → nc-health-services-medicaid-managed-care;
  CCBHC model → medicaid-imd-ccbhc.
version: "1.1.0"
updated: "2026-06-22"
category: health
tags: [medicaid, managed-care, behavioral-health, carve-out, coverage]
hub: health
---

# Medicaid Managed Care and Behavioral Health

## Overview

As of July 1, 2024, 78% of Medicaid beneficiaries (over 66 million enrollees) receive care through risk-based MCOs [KFF Medicaid Managed Care Tracker, 2024; MACPAC MACStats Feb 2026]. How states structure BH services within managed care — carve-in vs. carve-out, integrated vs. specialty — profoundly shapes access, accountability, and outcomes. Behavioral health services (MH and SUD) are among the most commonly carved-out benefit categories, alongside dental and NEMT [KFF 10 Things to Know About Medicaid Managed Care].

## Carve-Out vs. Carve-In Models

### Behavioral Health Carve-Out
- BH services (mental health and/or SUD) managed by a **separate behavioral health organization (BHO)** from the primary MCO
- Rationale: BH needs specialized expertise, catastrophic cost risk management, dedicated network
- Types:
  - **Partial carve-out**: specialty BH (inpatient psychiatric, SUD residential) to BHO; routine BH in MCO
  - **Full carve-out**: all BH services to BHO; primary care MCO handles physical health only
- **Accountability gap**: care coordination between physical health MCO and BHO historically weak; hospital readmissions fall through gaps

### Behavioral Health Carve-In (Integrated Managed Care)
- BH services included within the same MCO capitation as physical health
- Growing trend since ACA and CMS emphasis on integration
- MCO must demonstrate BH network and expertise; often subcontracts with BH specialty network
- Better positioned for care integration, comorbidity management, and whole-person care
- **NC model**: NC Medicaid transitioned to integrated managed care (Standard Plans) in 2021–2023, with BH I/DD Tailored Plans for high-need populations (launched July 2024)

## Federal Regulatory Framework

### Medicaid Managed Care Final Rule (2016, updated 2024)
- **42 CFR Part 438** governs Medicaid MCO contracts [eCFR, current text]
- Original 2016 final rule: 81 FR 27498 (May 6, 2016) [GovInfo FR-2016-05-06/2016-09581]; major 2024 update: CMS-2439-F, 89 FR [page unconfirmed] (May 10, 2024), effective July 9, 2024 [Federal Register document 2024-08085]
- Requires states to ensure MCOs provide all Medicaid-covered services, including BH
- Key 2024 provisions for BH:
  - **Network adequacy**: states must establish time/distance and appointment wait time standards for BH providers; routine outpatient MH/SUD appointments must not exceed **10 business days** from date of request (90% compliance rate required; requirement takes effect for contract rating periods beginning on or after **July 9, 2027** — 3 years from the rule's effective date) [Georgetown CCF Aug 2024; CMS-2439-F]
  - **LTSS coordination**: MCOs serving LTSS enrollees must coordinate with HCBS providers
  - **Quality measurement**: states must implement quality strategy including BH measures

### Parity in Medicaid Managed Care
- MHPAEA applied to Medicaid MCOs via **2016 final rule** (81 FR 18390–18445, March 30, 2016), codified at **42 CFR Part 438 Subpart K (§§ 438.900–438.920)** — not a 2014 rule [Federal Register 2016-06876; eCFR 42 CFR 438 Subpart K; confirmed 81 FR 18390]
- Compliance with Subpart K required by October 2, 2017 for MCOs, PIHPs, and PAHPs [42 CFR 438.920]
- Managed care contracts must comply with parity; states must include parity requirements in contracts
- MCO BH NQTLs (network adequacy standards, PA criteria, utilization management criteria) must be no more restrictive than for analogous physical health conditions [42 CFR 438.910]
- Separate September 9, 2024 MHPAEA final rule from HHS/DOL/Treasury added NQTL comparative analysis content requirements [CMS press release, Sept 2024]
- **Enforcement gap**: Medicaid managed care parity enforcement historically weaker than commercial insurance; 2024 rules strengthen state obligations

## Service Requirements

MCOs must cover all Medicaid BH services included in the state plan and any applicable waiver:
- Inpatient psychiatric (subject to IMD exclusion limitations)
- Outpatient mental health (individual/group therapy, medication management, assessment)
- SUD treatment continuum (prevention, outpatient, intensive outpatient, residential, medication-assisted treatment)
- Peer support services (if state has elected under state plan)
- Crisis services (mobile crisis, crisis stabilization): ARP Section 9813 (P.L. 117-2, enacted March 11, 2021) created an optional state plan benefit for community-based mobile crisis intervention services with **85% enhanced FMAP** for the first 12 fiscal quarters (April 1, 2022–March 31, 2027) [Medicaid.gov mobile crisis page; CMS SHO #21-008]
- CCBHC services: made a permanent optional Medicaid state plan benefit by **Consolidated Appropriations Act, 2024 (P.L. 118-42, Sec. 209)**, signed March 9, 2024 — any state may now add CCBHC services without a demonstration waiver [Medicaid.gov CCBHC Demonstration; Congress.gov R48075]

## Network Adequacy Requirements

Per **42 CFR 438.68** [eCFR current text], states must set time/distance and appointment wait time standards for:
- Mental health professionals (psychiatrists, psychologists, LCSWs, LPCs — adult and pediatric separately)
- SUD treatment facilities and providers
- States must consider geographic location, distance, travel time, and characteristics of specific Medicaid populations [42 CFR 438.68]
- **2024 rule** added mandatory **appointment wait time standards**: routine outpatient MH/SUD ≤ 10 business days from request (compliance = 90% of requests fulfilled within standard; takes effect for contract rating periods beginning on or after **July 9, 2027** — 3 years from effective date) [CMS-2439-F; Georgetown CCF Aug 2024]
- The general 2-year transition applies to certain other network adequacy provisions; the appointment wait-time standard specifically uses a 3-year transition [CMS-2439-F]
- **Telehealth**: MCOs may count telehealth providers toward network adequacy; 2024 rule acknowledges telehealth but requires in-person access for some services

## Encounter Data and Quality

- MCOs must submit encounter data for all services including BH [42 CFR Part 438 Subpart E]
- HEDIS BH measures in Medicaid managed care [NCQA HEDIS measures; CMS reporting requirements]:
  - **Antidepressant Medication Management (AMM)**: adherence to depression medication for 6 months
  - **Follow-Up After Hospitalization for Mental Illness (FUH)**: follow-up within 7 days and 30 days of discharge [NCQA FUH measure]
  - **Follow-Up After Emergency Department Visit for Mental Illness (FUM)**
  - **Initiation and Engagement of Alcohol and Other Drug Abuse or Dependence Treatment (IET)**
  - **Follow-Up After High-Intensity Care for SUD (FUI)**
- States may include BH quality in MCO quality incentive payments (QIPs)
- 2024 final rule strengthened quality measurement and external quality review (EQR) requirements [89 FR 34002]

## Care Coordination and Integration Requirements

2024 Medicaid managed care rule (89 FR 34002) strengthened requirements [Georgetown CCF Aug 2024]:
- MCOs must implement care management for high-need members (including BH/physical health comorbidities)
- "Comprehensive risk stratification" must identify BH needs
- **Transitions of care**: MCOs must ensure timely follow-up after BH hospitalization discharge; HEDIS FUH measure tracks 7-day and 30-day post-discharge follow-up rates [NCQA FUH; 89 FR 34002]
- **Health homes** (Section 2703 of ACA): state plan option integrating care management for people with chronic conditions including SMI; states receive **90% enhanced FMAP** for health home services for the **first 8 fiscal quarters** (10 quarters for SUD-focused plans approved after Oct 1, 2018) [Medicaid.gov Health Homes FAQ Dec 2017; Medicaid.gov Health Home Information Resource Center]

## Capitation Rate-Setting for BH

- BH capitation cells separated from physical health in many states for transparency
- Actuarially sound rates required; BH is high-volatility (catastrophic utilization spikes)
- **Risk corridors and reinsurance**: states often include stop-loss for high-cost BH episodes
- SUD residential (if state obtained IMD 15-day exception) increases actuarial complexity

## State Variations

| Approach | Examples |
|---|---|
| Full BH carve-out | Massachusetts (MassHealth/MBHP manages BH for PC ACOs and PCC Plan) [Mass.gov MBHP EQR 2024] |
| BH carve-in with specialty subcontract | North Carolina Standard Plans, Oregon, Arizona |
| BH/Tailored Plans for high-acuity | North Carolina Tailored Plans (launched July 1, 2024 after delays) [NCDHHS press release June 27, 2024], Oregon CCOs |
| County-based managed care | California Regional Centers for I/DD |

Note: NC Tailored Plans serve ~210,000 beneficiaries with SMI, SED, severe SUD, I/DD, or TBI; launch was delayed twice (from April 2023, then October 2023) before the July 1, 2024 launch [NCDHHS Feb 2023; July 2023; June 2024 press releases].

## References

Verified-as-of: 2026-06-22

1. **42 CFR Part 438 – Managed Care** (current text, incl. Subpart K parity and § 438.68 network adequacy)
   https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-C/part-438

2. **2024 Medicaid Managed Care Final Rule (CMS-2439-F)** — FR doc 2024-08085, May 10, 2024, eff. July 9, 2024 (FR starting page unconfirmed from primary source; document number verified)
   https://www.federalregister.gov/documents/2024/05/10/2024-08085/medicaid-program-medicaid-and-childrens-health-insurance-program-chip-managed-care-access-finance

2a. **2016 Medicaid Managed Care Final Rule** — 81 FR 27498, May 6, 2016
    https://www.govinfo.gov/app/details/FR-2016-05-06/2016-09581

3. **2016 MHPAEA Medicaid Managed Care Final Rule** — 81 FR 18390–18445, March 30, 2016 (codified at 42 CFR Part 438 Subpart K, §§ 438.900–438.920); verified via GovInfo search
   https://www.federalregister.gov/documents/2016/03/30/2016-06876/medicaid-and-childrens-health-insurance-programs-mental-health-parity-and-addiction-equity-act-of

4. **KFF – 10 Things to Know About Medicaid Managed Care** (enrollment data, carve-out prevalence; cites 78% figure as of July 1, 2024)
   https://www.kff.org/medicaid/10-things-to-know-about-medicaid-managed-care/

5. **MACPAC – Percentage of Medicaid Enrollees in Managed Care by State (FY 2022); MACStats Feb 2026**
   https://www.macpac.gov/publication/percentage-of-medicaid-enrollees-in-managed-care-by-state-and-eligibility-group/

6. **Georgetown CCF – A Closer Look at the Access Provisions in Final Medicaid Managed Care Rule** (May 15, 2024; covers 10-business-day BH wait time standard)
   https://ccf.georgetown.edu/2024/05/15/a-closer-look-at-the-access-provisions-in-final-medicaid-managed-care-rule/

7. **Georgetown CCF – Explanation of Final Medicaid Managed Care and Access Rules** (Aug 7, 2024)
   https://ccf.georgetown.edu/2024/08/07/an-explanation-of-final-medicaid-managed-care-and-access-rules/

8. **ARP Section 9813 – Mobile Crisis (P.L. 117-2)** — 85% FMAP, 12 quarters, April 2022–March 2027
   https://www.medicaid.gov/medicaid/benefits/behavioral-health-services/state-option-provide-qualifying-community-based-mobile-crisis-intervention-services

9. **Consolidated Appropriations Act 2024 (P.L. 118-42, Sec. 209)** — made CCBHC permanent Medicaid state plan option
   https://www.medicaid.gov/medicaid/financial-management/certified-community-behavioral-health-clinic-ccbhc-demonstration
   CRS summary: https://www.congress.gov/crs-product/R48075

10. **Medicaid Health Homes (Section 2703)** — 90% enhanced FMAP for 8 quarters (10 for SUD plans)
    https://medicaid.gov/medicaid/ltss/health-homes/index.html
    FAQ: https://www.medicaid.gov/state-resource-center/medicaid-state-technical-assistance/health-home-information-resource-center/downloads/health-homes-faq-12-18-17.pdf

11. **NCQA – Follow-Up After Hospitalization for Mental Illness (FUH)** (7-day and 30-day measures)
    https://www.ncqa.org/report-cards/health-plans/state-of-health-care-quality-report/follow-up-after-hospitalization-for-mental-illness-fuh/

12. **NCDHHS – NC Medicaid Launches Behavioral Health and I/DD Tailored Plans July 1** (June 27, 2024)
    https://www.ncdhhs.gov/news/press-releases/2024/06/27/nc-medicaid-launches-behavioral-health-and-intellectualdevelopmental-disabilities-tailored-plans

13. **Massachusetts MBHP External Quality Review Technical Report 2024**
    https://www.mass.gov/doc/massachusetts-behavioral-health-partnership-mbhp-eqr-technical-report-2024-0/download

14. **CMS – MHPAEA Fact Sheet and 2024 Final Rule (Sept 9, 2024)**
    https://www.cms.gov/marketplace/private-health-insurance/mental-health-parity-addiction-equity
