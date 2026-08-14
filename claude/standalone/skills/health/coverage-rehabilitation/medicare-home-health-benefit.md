---
name: medicare-home-health-benefit
description: >-
  Medicare home health benefit — Part A and Part B coverage for skilled nursing,
  physical therapy, occupational therapy, speech-language pathology, home health
  aide, and medical social services delivered in the home; homebound status
  definition and documentation; Home Health Prospective Payment System (HHPPS)
  and PDGM (Patient-Driven Groupings Model); certification/recertification by
  physician; Review Choice Demonstration pre-claim review in select states; home
  health compare and star ratings; coverage of mental health/behavioral services
  in home health including psychiatric homebound criteria; no cost-sharing for
  home health; Medicare Advantage differences; expanded HHVBP model. Educational
  only, NOT medical/insurance advice; verified-as-of 2026-06-22.
  TRIGGER: Medicare home health, Medicare home health benefit, homebound status
  Medicare, Medicare home health aide, Medicare skilled nursing home, Medicare
  physical therapy at home, Medicare occupational therapy home, Medicare speech
  therapy home, PDGM home health, Medicare home health prospective payment, home
  health certification Medicare, home health recertification, Medicare home health
  prior authorization, Medicare home health cost sharing, Medicare Home Health
  Compare, Medicare home health aide hours, Medicare home health mental health,
  Medicare social worker home visit, psychiatric home health Medicare, Medicare
  homebound psychiatric, home health value-based purchasing.
  SKIP: Medicare SNF inpatient → medicare-snf-irf-coverage; Medicare outpatient
  PT/OT/speech clinic → medicare-outpatient-rehabilitation; Medicare hospice →
  see hospice; Medicaid home and community-based services → medicaid-ltss-hcbs-rehab.
version: "1.1.0"
updated: "2026-06-22"
category: health
tags: [medicare, home-health, rehabilitation, behavioral-health, coverage]
hub: health
---

# Medicare Home Health Benefit

## Overview

Medicare covers home health services for beneficiaries who are **homebound** and require **skilled care**. The benefit spans Parts A and B (with no defined episode limit), covers no deductible or coinsurance (for Original Medicare), and is administered through Certified Home Health Agencies (CHHAs).

## Homebound Status

**Regulatory basis:** 42 CFR §409.42(a); CMS Medicare Benefit Policy Manual (Pub. 100-02) Chapter 7, §30.1.

A beneficiary is "homebound" when:
- Leaving home requires **considerable and taxing effort** due to illness, injury, or functional limitation (needing assistance of another person or a supportive device); **OR** leaving home is **medically contraindicated**
- Absences from home are infrequent, brief, or for medical treatment / adult day programs / religious services / occasional outings
- Permanent homebound status is NOT required; condition may be temporary or expected to improve

**Psychiatric homebound criteria (CMS Benefit Policy Manual Ch. 7, §30.1):** A patient with a psychiatric illness is considered homebound if the illness is manifested **in part by a refusal to leave home**, OR if the condition is of such nature that **it would not be considered safe for the patient to leave home unattended**, even without physical limitations. This is a distinct basis for homebound status that does not require inability to ambulate. Example: a beneficiary recently hospitalized for major depression who refuses to leave home and has suicidal ideation may qualify as homebound on psychiatric grounds alone.

**Documentation burden:** The certifying physician must document the clinical basis for homebound status in the certification. Audits focus on this element; inadequate documentation is the most common cause of Medicare home health overpayments (OIG). [verified-as-of 2026-06-22]

## Covered Services

| Service | Requirement |
|---|---|
| Skilled nursing care | Intermittent or part-time; skilled nurse (RN or LPN) providing assessment, teaching, or complex wound care |
| Physical therapy | Must be medically necessary; can establish eligibility for other services |
| Occupational therapy | Maintain eligibility once established, even after SN/PT end |
| Speech-language pathology | Swallowing, communication, cognitive deficits |
| Home health aide | Only if skilled service also needed; personal care (bathing, dressing) — NOT custodial standalone |
| Medical social services | Counseling, resource coordination; social worker visits authorized under skilled plan |

**Mental health services in home health:** Home health aides cannot provide stand-alone psychiatric/mental health services. Psychiatric nursing can be covered as a skilled nursing service when provided by a **psychiatrically trained RN** — defined as an RN with special training and/or experience beyond the standard nursing curriculum; an RN with a Master's degree specializing in psychiatric or mental health nursing qualifies (LCD L34561). Covered psychiatric skilled nursing services include: observation/assessment, teaching/training, management and evaluation of a patient care plan, medication management, safety monitoring, and direct patient care for a diagnosed psychiatric condition. The principal diagnosis must map to an ICD-10 code in the **Behavioral Health** clinical grouping under PDGM.

## Not Covered Under Home Health

- 24-hour/continuous care
- Custodial care alone (ADL help, housekeeping)
- Personal care that is not part of a skilled plan
- Meals delivered to home
- Homemaker services

## Certification Requirements

**Regulatory basis:** 42 CFR §424.22.

- Physician (MD/DO) or allowed practitioner (NP/PA/CNS/certified nurse-midwife) must **certify** homebound status and need for skilled care
- Face-to-face encounter required **within 90 days before or 30 days after** home health start of care; may be conducted via telehealth
- Certification valid for up to **60-day certification periods**; however, payment occurs in **30-day billing periods** under PDGM
- Recertification required if skilled need continues beyond the certification period
- Home Health Agency must be **Medicare-certified**

## Payment: PDGM (Patient-Driven Groupings Model)

**Source:** CMS PDGM overview; 88 Fed. Reg. 77637 (Nov. 13, 2023) (CY 2024 final rule). Effective January 1, 2020.

Replaced the prior Therapy Threshold model:
- Episodes divided into **30-day billing periods** (payment unit); certifications remain 60-day periods
- Each 30-day period grouped across five dimensions → **432 case-mix groups**:
  1. Admission source (community vs. post-acute, 2 subgroups)
  2. Timing (early vs. late within certification, 2 subgroups)
  3. **Clinical grouping** (12 subgroups): Musculoskeletal Rehabilitation; Neuro/stroke Rehabilitation; Wounds; MMTA-Surgical Aftercare; MMTA-Cardiac and Circulatory; MMTA-Endocrine; MMTA-GI Tract and GU System; MMTA-Infectious Disease/Neoplasms/Blood-forming; MMTA-Respiratory; MMTA-Other; **Behavioral Health**; Complex Nursing Interventions
  4. Functional impairment level (3 subgroups)
  5. Comorbidity adjustment (3 subgroups)
- Removes the therapy visit threshold (old system incentivized therapy volume)
- Psychiatric nursing and social work visits for mental health diagnoses are captured in the **Behavioral Health** clinical grouping

**Note:** The acronym "BHHH" (Behavioral Health Home Health) sometimes appears in industry materials but is not a CMS-designated program name. The PDGM grouping is simply "Behavioral Health."

## Review Choice Demonstration (Pre-Claim Review)

**Source:** CMS Review Choice Demonstration page (cms.gov); extended June 1, 2024 for 5 additional years.

CMS implemented the **Review Choice Demonstration (RCD) for Home Health Services** — this is a pre-claim review / post-payment review program, **not** traditional prior authorization. Home health agencies in demonstration states choose from review options:
- **Choice 1:** Pre-claim review (submit before billing; provisional affirmation before claim submission)
- **Choice 2:** Post-payment review
- **Choice 3 (removed June 2024):** Minimal review with 25% payment reduction — eliminated from initial choice selections

**Demonstration states (as of 2026):** Illinois, Ohio, Texas, North Carolina, Florida, Oklahoma (Palmetto GBA MAC jurisdiction). Michigan and Massachusetts are **not** current demonstration states.

Agencies demonstrating 90%+ provisional affirmation rate (minimum 10 requests) may graduate to reduced review burden. The program does not alter the Medicare home health benefit or delay care.

## Value-Based Purchasing: Expanded HHVBP Model

**Source:** CMS Expanded HHVBP Model page; CY 2026 HH PPS Final Rule (90 Fed. Reg., Dec. 2, 2025).

The Expanded Home Health Value-Based Purchasing (HHVBP) Model launched January 1, 2022 and covers **all 50 states, the District of Columbia, and U.S. territories** — every Medicare-certified HHA participates. Payment adjustments of up to **±5%** based on performance (CY 2024 was performance year 2; CY 2026 is payment year 2). CMS compares HHAs within smaller- and larger-volume national cohorts.

## OASIS Assessment: Behavioral Health Items

**Source:** CMS OASIS-E Manual (updated January 1, 2024); OASIS-E1 (effective January 1, 2025); OASIS-E2 (effective April 1, 2026).

OASIS collects standardized behavioral health data elements at start of care, resumption of care, follow-up, and discharge:
- **PHQ-2 / PHQ-9** (Patient Health Questionnaire): depression screening and severity
- **Social isolation** items
- **Cognitive, behavioral, and psychiatric symptoms** items
- Functional status (ADL/IADL) items linked to PDGM functional impairment grouping

OASIS-E1 (2025) and OASIS-E2 (2026) added and refined BH-adjacent items per CMS quality measure updates.

## Cost Sharing

- **No deductible** and **no coinsurance** for Medicare-approved home health services under Original Medicare
- Home health aide visits: no cost sharing
- **Medicare Advantage**: may require copays for home health; check plan-specific Evidence of Coverage

## Medicare Advantage Differences

- MA plans must cover home health and must comply with national coverage determinations (NCDs), local coverage determinations (LCDs), and Traditional Medicare coverage criteria
- MA plans may apply copays, require in-network CHHA, or impose prior authorization; the 2024 Interoperability and Prior Authorization Rule (CMS-0057-F, effective 2026) requires MA plans to respond to standard prior authorization requests within **7 calendar days** and expedited requests within **72 hours** [verified-as-of 2026-06-22]
- Plans may operationally differ in application of homebound criteria; CMS has required plans to meet traditional Medicare NCD/LCD standards

## Home Health Compare / Care Compare

CMS publishes Home Health star ratings on the Care Compare platform (CareCompare.cms.gov):
- Quality of patient care star ratings
- Patient survey (HHCAHPS) star ratings
- Claims-based outcomes (hospitalization, ED use)

The CY 2026 HHVBP model removed three HHCAHPS survey-based measures and added four new measures (three OASIS-based measures for bathing/dressing; one claims-based Medicare Spending per Beneficiary for Post-Acute Care measure).

## Behavioral Health Applications

- **Psychiatric home health nursing**: Covered skilled service when performed by a psychiatrically trained RN. Eligible conditions include schizophrenia, bipolar disorder, MDD, PTSD, dementia with behavioral symptoms, and other DSM-5-diagnosed psychiatric disorders (per LCD L34561). Services include medication reconciliation, psychoeducation, cognitive assessment, and safety planning.
- **Medical social work**: Connects patients with BH community resources, assists with advance directives; authorized under a skilled plan of care.
- **Telehealth in home health (Original Medicare):** Telehealth visits do **not** count as home health visits and are not billable as home health services. They may be billed separately under Part B when provided by a physician/practitioner. Pandemic-era telehealth flexibilities allowing any site (including the home) as originating site have been extended through **December 31, 2027** by subsequent appropriations acts. The face-to-face certification encounter may be conducted via telehealth (42 CFR §424.22). [verified-as-of 2026-06-22]

## References

1. **42 CFR §409.42** — Beneficiary qualifications (homebound confinement criterion): https://www.law.cornell.edu/cfr/text/42/409.42
2. **42 CFR §424.22** — Certification and face-to-face encounter requirements: https://www.law.cornell.edu/cfr/text/42/424.22
3. **CMS Medicare Benefit Policy Manual, Pub. 100-02, Chapter 7** — Home Health Services (homebound §30.1; psychiatric homebound §30.1.1; psychiatric nursing §40.2): https://www.cms.gov/Regulations-and-Guidance/Guidance/Manuals/Downloads/bp102c07.pdf
4. **LCD L34561 — Home Health Psychiatric Care** (psychiatric nurse qualifications, DSM-5 diagnosis requirement, psychiatric homebound criteria): https://www.cms.gov/medicare-coverage-database/view/lcd.aspx?lcdId=34561&ver=69
5. **CMS PDGM Overview** — 12 clinical groupings, 432 case-mix groups, behavioral health grouping: https://www.cms.gov/medicare/payment/prospective-payment-systems/home-health/home-health-patient-driven-groupings-model
6. **CY 2024 HH PPS Final Rule (CMS-1780-F)** — PDGM recalibration, HHVBP requirements: https://www.cms.gov/newsroom/fact-sheets/calendar-year-cy-2024-home-health-prospective-payment-system-final-rule-cms-1780-f
7. **CY 2025 HH PPS Final Rule (CMS-1803-F)**: https://www.cms.gov/newsroom/fact-sheets/calendar-year-cy-2025-home-health-prospective-payment-system-final-rule-fact-sheet-cms-1803-f
8. **CY 2026 HH PPS Final Rule (CMS-1828-F)** — HHVBP measure changes, OASIS-E2: https://www.cms.gov/newsroom/fact-sheets/calendar-year-cy-2026-home-health-prospective-payment-system-final-rule-cms-1828-f
9. **CMS Review Choice Demonstration for Home Health Services** — states (IL, OH, TX, NC, FL, OK), June 2024 extension, Choice 3 removal: https://www.cms.gov/data-research/monitoring-programs/medicare-fee-service-compliance-programs/prior-authorization-and-pre-claim-review-initiatives/review-choice-demonstration-home-health-services
10. **Expanded HHVBP Model** — all 50 states + DC + territories, ±5% payment adjustment: https://www.cms.gov/priorities/innovation/innovation-models/expanded-home-health-value-based-purchasing-model
11. **OASIS-E1 Manual** (effective January 1, 2025) — PHQ-2/9, behavioral health items: https://www.cms.gov/files/document/draft-oasis-e1-manual-04-28-2024.pdf
12. **CMS Interoperability and Prior Authorization Final Rule (CMS-0057-F)** — MA prior authorization timeframes (7-day standard / 72-hour expedited): https://www.cms.gov/newsroom/fact-sheets/cms-interoperability-prior-authorization-final-rule-cms-0057-f
