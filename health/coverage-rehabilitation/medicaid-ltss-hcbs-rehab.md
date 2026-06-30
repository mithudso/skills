---
name: medicaid-ltss-hcbs-rehab
description: >-
  Medicaid Long-Term Services and Supports (LTSS) and Home and Community-Based
  Services (HCBS) for rehabilitation and behavioral health: 1915(c) HCBS
  waivers, 1915(i) state plan HCBS, 1915(k) Community First Choice (CFC) option,
  Section 1115 HCBS waivers, HCBS settings rule (home-like, non-institutional),
  HCBS rehab services (supported employment, skills training, community living
  supports), Medicaid LTSS coverage of behavioral health in nursing facilities,
  PACE program rehab components, Money Follows the Person (MFP) transitions, and
  Medicaid managed care LTSS (MLTSS) models. Educational only, NOT
  medical/insurance advice; verified-as-of 2026-06-22.
  TRIGGER: Medicaid LTSS, Medicaid home and community-based services, Medicaid
  HCBS waiver, Medicaid 1915(c) waiver, 1915(i) Medicaid state plan HCBS,
  Community First Choice Medicaid, HCBS settings rule, Medicaid supported
  employment IPS, Medicaid skills training community living, Medicaid nursing
  facility behavioral health, PACE rehabilitation Medicaid, Money Follows the
  Person, Medicaid managed LTSS, MLTSS, Medicaid institutional vs community care,
  Medicaid rebalancing, Medicaid functional assessment LTSS.
  SKIP: Medicaid rehabilitative services state plan option → medicaid-rehabilitative-services;
  IMD exclusion → medicaid-imd-ccbhc; Medicaid enrollments/eligibility →
  enrollment-eligibility-*; NC specific HCBS waivers → nc-health-services-medicaid-waivers;
  Medicare home health → medicare-home-health-benefit.
version: "1.1.0"
updated: "2026-06-22"
category: health
tags: [medicaid, ltss, hcbs, rehabilitation, behavioral-health, community-based]
hub: health
---

# Medicaid LTSS and HCBS Rehabilitation Services

## Overview

Medicaid Long-Term Services and Supports (LTSS) comprises institutional care (nursing facilities, ICF/IID) and Home and Community-Based Services (HCBS). HCBS allows states to provide rehabilitation and behavioral health services to people who need LTSS-level support while living in the community — a policy priority called "rebalancing" away from institutional care.

## HCBS Authority Landscape

### 1915(c) HCBS Waivers (the workhorse)
- States request CMS approval to waive comparability, statewide, and freedom-of-choice requirements (SSA § 1915(c))
- Must serve people who would otherwise require institutionalization in a nursing facility, ICF/IID, or hospital — the cost-neutrality test caps average per-capita HCBS spend at 100% of institutional equivalent ([MACPAC, 1915(c) waivers](https://www.macpac.gov/subtopic/1915-c-waivers/); verified 2026-06-22)
- Services: personal care, respite, supported employment, adult day services, assistive technology, skills training, behavior supports, community transition services, crisis stabilization ([CMS national overview](https://www.cms.gov/outreach-and-education/american-indian-alaska-native/aian/ltss-ta-center/info/national-overview-1915-c-waivers))
- Initial waiver approval: 3 years; renewals up to 5 years; dually eligible waivers may receive 5-year initial approval
- Enrollment caps allowed (waiver = not entitlement); waiting lists common
- **Behavioral health 1915(c) waivers** cover: individual skills training, behavioral consultation, mental health peer support, crisis respite, supported housing navigation ([MACPAC, BH services under HCBS waivers](https://www.macpac.gov/subtopic/behavioral-health-services-covered-under-hcbs-waivers-and-spas/))
- **NC examples**: NC Innovations (I/DD — renewed through June 30, 2029), CAP-DA (older adults), CAP-C (children medically fragile) ([NC Medicaid Innovations Waiver](https://medicaid.ncdhhs.gov/beneficiaries/nc-innovations-waiver); verified 2026-06-22)

### 1915(i) State Plan HCBS (entitlement, no waiver required)
- Eligible individuals: functional need below nursing-facility level
- Allows states to define target population via functional criteria (42 CFR Part 441, Subpart M; [eCFR 441.710](https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-C/part-441/subpart-M))
- HCBS services added to state plan — no enrollment cap (entitlement class)
- Less utilized by states than 1915(c) due to entitlement risk; useful for BH populations

### 1915(k) Community First Choice (CFC) Option
- Established by ACA Section 2401; available to states since October 1, 2011 ([CMS CFC fact sheet](https://www.cms.gov/newsroom/fact-sheets/community-first-choice-option-section-1915k); [medicaid.gov CFC page](https://www.medicaid.gov/medicaid/home-community-based-services/home-community-based-services-authorities/community-first-choice-cfc-1915-k))
- Provides attendant care services and supports to individuals at institutional level of care living in the community
- Enhanced FMAP: additional **6 percentage points** above regular FMAP for CFC services (verified 2026-06-22)
- Covers: hands-on assistance, supervision, skills building for ADLs/IADLs
- Settings requirements for CFC governed by 42 CFR 441.530 ([eCFR](https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-C/part-441/subpart-K/section-441.530))

### Section 1115 HCBS Demonstrations
- States use 1115 waivers to test innovative HCBS models
- Recent examples: Coverage of community transition services, SDOH supports (housing navigation, medically tailored meals), crisis stabilization in community settings
- **ARPA Section 9817**: provided a temporary **10 percentage point** FMAP increase for HCBS (not an 1115 authority; applicable earn period April 1, 2021–March 31, 2022; spending window extended to March 31, 2025) ([CMS press release](https://www.cms.gov/newsroom/press-releases/cms-issues-guidance-american-rescue-plan-funding-medicaid-home-and-community-based-services); [CMS spending deadline extension](https://www.cms.gov/newsroom/press-releases/hhs-extends-american-rescue-plan-spending-deadline-states-expand-and-enhance-home-and-community); verified 2026-06-22)

## HCBS Settings Rule (Published 2014; Compliance Deadline March 2023)

- Final rule published January 16, 2014 ([Federal Register 2014-00487](https://www.federalregister.gov/documents/2014/01/16/2014-00487/medicaid-program-state-plan-home-and-community-based-services-5-year-period-for-waivers-provider-payment))
- Settings requirements codified across multiple CFR sections: **42 CFR 441.301(c)(4)** (1915(c) waiver settings characteristics; (c)(5) settings that do NOT qualify), 42 CFR 441.530 (CFC/1915(k) settings), 42 CFR 441.710 (1915(i) settings) ([LII/Cornell CFR 441.301](https://www.law.cornell.edu/cfr/text/42/441.301); verified 2026-06-22)
- Home-like characteristics required: integrated in broader community, access to community activities, choice of roommate, privacy/dignity
- Presumptively institutional settings require CMS "heightened scrutiny" approval: settings on grounds of or adjacent to institutions, or settings that isolate individuals from the broader community
- Compliance deadline: initially 2020, extended to 2022, then to **March 2023** (COVID extension); enforcement ongoing — as of late 2023, 43 states reported full implementation in at least one waiver; most others under corrective action plans ([KFF settings rule implementation](https://www.kff.org/medicaid/how-are-states-implementing-new-requirements-for-medicaid-home-and-community-based-services/); verified 2026-06-22)
- **Impact on BH**: adult foster homes, group homes, and day programs serving people with mental illness must meet settings rule to bill HCBS

## Behavioral Health HCBS Services

| Service | Description |
|---|---|
| Individual Skills Training | Teach independent living, ADLs, social skills, symptom management |
| Supported Employment (IPS model) | Individual Placement and Support; integrated competitive employment with coaching |
| Assertive Community Treatment (ACT) | Intensive team-based community support; included as HCBS in some state plans |
| Peer Support | Peer specialists delivering recovery support in community settings |
| Crisis Stabilization (community) | Non-hospital crisis residential; short-stay community crisis programs |
| Behavioral Supports / PBS | Positive Behavioral Supports; applied behavior analysis for I/DD-BH comorbidity |
| Housing Navigation | Supported housing search, application, transitions |
| Benefits Counseling | Employment and benefits navigation (WorkIncentive Counseling) |

## LTSS in Nursing Facilities — Behavioral Health

Medicaid covers BH services in Nursing Facilities (NFs) if the facility provides them, but:
- Preadmission Screening and Resident Review (PASRR) required for individuals with MI, I/DD, or related conditions — authorized under SSA § 1919(e)(7); governed by **42 CFR Part 483, Subpart C (§§ 483.100–483.138)** ([eCFR 483.100](https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-G/part-483/subpart-C/section-483.100); verified 2026-06-22)
- PASRR Level I screen for all new admissions to Medicaid-certified NFs (on or after January 1, 1989); Level II evaluation for those with identified MI or I/DD
- PASRR Level II identifies need for specialized services beyond NF care
- Specialized services (psychiatric, behavioral, psychosocial rehabilitation) must be arranged outside the per-diem NF rate
- **Mentally Ill / Intellectually Disabled (MI/ID) exemptions**: residents admitted before January 1, 1989 grandfathered; residents needing nursing care for acute conditions temporarily exempt (42 CFR 483.106)

## PACE (Program of All-inclusive Care for the Elderly)

- Serves dual-eligible individuals ≥55 needing nursing-facility level of care; governed by 42 CFR Part 460 ([eCFR Part 460](https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-G/part-460))
- Capitated model: PACE organization receives combined Medicare + Medicaid payment
- PACE must provide all covered Medicare and Medicaid services, including:
  - PT, OT, speech therapy (onsite or by contract)
  - Mental health services and psychiatric consultation
  - Social work counseling and care management
  - Adult day center with rehabilitative programming
- **BH integration in PACE**: interdisciplinary team includes mental health professional; medication management; behavioral interventions for dementia

## Money Follows the Person (MFP)

- Demonstration grant first authorized by Deficit Reduction Act of 2005 (§ 6071); extended multiple times
- Most recent extension: **Consolidated Appropriations Act, 2023 (P.L. 117-328)** extended MFP through **FY2027** — NOT ARPA ([Congress.gov, CAA 2023 Medicaid provisions](https://www.congress.gov/crs-product/R47821); verified 2026-06-22)
- Helps Medicaid beneficiaries transition from institutional settings (NF, ICF/IID, hospital) to HCBS after **60+ consecutive days** of institutional residence (minimum stay reduced from 90 to 60 days by Consolidated Appropriations Act, 2021, effective January 26, 2021; days for short-term rehabilitative services now count)
- Enhanced FMAP for first year of community services: formula yields rates in the range of 75%–90% (individual state rate = state FMAP + half the difference between 100% and state FMAP) ([CMS MFP](https://mdctmfp.cms.gov/))
- States demonstrate that BH supports can enable community living for long-stay NF or IMD residents

## Medicaid Managed LTSS (MLTSS)

States increasingly enroll LTSS-eligible individuals in managed care (governed by 42 CFR Part 438; network adequacy standards at 42 CFR 438.68):
- MCO receives capitation to provide both LTSS and behavioral health
- **Care coordination**: MCO must coordinate HCBS, behavioral health, and acute care
- CMS managed care regulations (42 CFR 438) require MLTSS contracts to include network adequacy, care coordination, and quality standards; 2016 and 2024 rules expanded access requirements ([Federal Register, Ensuring Access 2024-08363](https://www.federalregister.gov/documents/2024/05/10/2024-08363/medicaid-program-ensuring-access-to-medicaid-services))
- **Financial risk in MLTSS**: capitation must be actuarially sound; HCBS underpayment risk if MLTSS not structured correctly

## Key Policy Tensions

1. **Waiting lists vs entitlement**: 1915(c) allows caps; 1915(i) does not — states resist 1915(i) due to open-ended cost
2. **BH integration gap**: most LTSS systems were designed for I/DD or aging; adults with serious mental illness (SMI) often poorly served
3. **Provider workforce**: direct care worker shortage limits HCBS expansion despite policy intent
4. **Housing**: HCBS can fund housing navigation but not rent itself (HCBS as housing costs exception is narrow); 1115 proposals test bridge subsidies

## References

1. MACPAC. "1915(c) Waivers." <https://www.macpac.gov/subtopic/1915-c-waivers/> (verified 2026-06-22)
2. MACPAC. "Behavioral Health Services Covered Under HCBS Waivers and 1915(i) SPAs." <https://www.macpac.gov/subtopic/behavioral-health-services-covered-under-hcbs-waivers-and-spas/> (verified 2026-06-22)
3. CMS. "National Overview of 1915(c) HCBS Waivers." <https://www.cms.gov/outreach-and-education/american-indian-alaska-native/aian/ltss-ta-center/info/national-overview-1915-c-waivers> (verified 2026-06-22)
4. Federal Register. "Medicaid Program; State Plan Home and Community-Based Services … and Home and Community-Based Setting Requirements …" Jan. 16, 2014 (Document No. 2014-00487). <https://www.federalregister.gov/documents/2014/01/16/2014-00487/medicaid-program-state-plan-home-and-community-based-services-5-year-period-for-waivers-provider-payment> (verified 2026-06-22)
5. KFF. "How Are States Implementing New Requirements for Medicaid Home- and Community-Based Services?" <https://www.kff.org/medicaid/how-are-states-implementing-new-requirements-for-medicaid-home-and-community-based-services/> (verified 2026-06-22)
6. medicaid.gov. "Community First Choice (CFC) 1915(k)." <https://www.medicaid.gov/medicaid/home-community-based-services/home-community-based-services-authorities/community-first-choice-cfc-1915-k> (verified 2026-06-22)
7. CMS. "Community First Choice Option Section 1915(k) — Fact Sheet." <https://www.cms.gov/newsroom/fact-sheets/community-first-choice-option-section-1915k> (verified 2026-06-22)
8. CMS. "CMS Issues Guidance on American Rescue Plan Funding for Medicaid HCBS." <https://www.cms.gov/newsroom/press-releases/cms-issues-guidance-american-rescue-plan-funding-medicaid-home-and-community-based-services> (verified 2026-06-22)
9. CMS. "HHS Extends American Rescue Plan Spending Deadline for States to Expand and Enhance HCBS." <https://www.cms.gov/newsroom/press-releases/hhs-extends-american-rescue-plan-spending-deadline-states-expand-and-enhance-home-and-community> (verified 2026-06-22)
10. Congress.gov (CRS). "Consolidated Appropriations Act, 2023 (P.L. 117-328): Medicaid and CHIP Provisions." <https://www.congress.gov/crs-product/R47821> (MFP extended through FY2027; verified 2026-06-22)
11. CMS. "Money Follows the Person Data Collection Tool." <https://mdctmfp.cms.gov/> (verified 2026-06-22)
12. eCFR. "42 CFR Part 483, Subpart C — Preadmission Screening and Annual Review (PASRR)." <https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-G/part-483/subpart-C> (verified 2026-06-22)
13. eCFR. "42 CFR 441.530 — Home and Community-Based Setting (CFC)." <https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-C/part-441/subpart-K/section-441.530> (verified 2026-06-22)
14. eCFR. "42 CFR 441.710 — State Plan HCBS under 1915(i)." <https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-C/part-441/subpart-M/section-441.710> (verified 2026-06-22)
15. NC Medicaid. "NC Innovations Waiver." <https://medicaid.ncdhhs.gov/beneficiaries/nc-innovations-waiver> (renewed through June 30, 2029; verified 2026-06-22)
16. LII / Cornell Law School. "42 CFR § 441.301 — Contents of request for a waiver." <https://www.law.cornell.edu/cfr/text/42/441.301> (HCBS settings at (c)(4); verified 2026-06-22)
17. eCFR. "42 CFR Part 460 — Program of All-inclusive Care for the Elderly (PACE)." <https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-G/part-460> (verified 2026-06-22)
18. Federal Register. "Medicaid Program; Ensuring Access to Medicaid Services" (Final Rule, 2024). <https://www.federalregister.gov/documents/2024/05/10/2024-08363/medicaid-program-ensuring-access-to-medicaid-services> (MLTSS network adequacy; verified 2026-06-22)
