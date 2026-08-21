---
name: medicare-advantage-behavioral-health
description: >-
  Medicare Advantage (Part C) behavioral health coverage — how MA plans cover
  mental health and SUD services differently from Original Medicare: MA BH
  supplemental benefits (in-home mental health support, caregiver support),
  MA prior authorization for BH services (2026 compliance timelines, CMS
  oversight), MA network adequacy for psychiatrists and BH providers, D-SNP
  behavioral health coordination requirements, MA SUD benefit coverage vs. OTP,
  MA inpatient psychiatric coverage and the 190-day limit applicability, MA
  telehealth mental health expansion, BH carve-in vs. carve-out within MA, MA
  Special Needs Plans for chronic conditions (C-SNPs for mental illness), and
  how MA BH coverage compares to Original Medicare + Medigap. Educational
  only, NOT medical/insurance advice; verified-as-of 2026-06-22.
  TRIGGER: Medicare Advantage mental health coverage, Medicare Advantage
  psychiatry, Medicare Advantage SUD treatment, Medicare Advantage behavioral
  health prior authorization, Medicare Advantage BH network, Medicare Advantage
  supplemental mental health benefits, Medicare Advantage telehealth mental
  health, D-SNP behavioral health, C-SNP mental illness Medicare Advantage,
  Medicare Advantage 190-day psychiatric limit, Medicare Advantage inpatient
  psych, Medicare Advantage OTP opioid treatment, Medicare Advantage vs original
  Medicare mental health, Medicare Advantage BH carve-out.
  SKIP: Medicare Advantage general structure → medicare-advantage-part-c; Original
  Medicare mental health/SUD → medicare-mental-health-sud-coverage; MA network
  adequacy standards → network-adequacy-behavioral-health; prior authorization
  reform general → mhpaea-prior-auth-discharge; D-SNP dual eligibility enrollment
  → enrollment-eligibility-dual-eligibility-medicare-medicaid.
version: "1.1.0"
updated: "2026-06-22"
category: health
tags: [medicare-advantage, behavioral-health, mental-health, SUD, coverage]
hub: health
---

# Medicare Advantage Behavioral Health Coverage

## Overview

Medicare Advantage (Part C) plans must cover all Original Medicare benefits, including behavioral health services, but they administer coverage differently — with their own networks, prior authorization requirements, and supplemental benefits. In practice, MA BH coverage varies dramatically by plan, creating both opportunities (supplemental benefits, integrated care) and risks (network restrictions, PA barriers, carve-out fragmentation).

## MA Must Cover: Original Medicare BH Benefits

MA plans must cover everything covered under Original Medicare Parts A and B:
- Inpatient psychiatric (Part A) — the 190-day lifetime limit applies to freestanding psychiatric hospitals (does NOT apply to distinct-part psychiatric units inside acute care hospitals)
- Outpatient psychotherapy and medication management (Part B)
- Partial hospitalization programs (PHP)
- Intensive outpatient programs (IOP) — added as a covered Original Medicare service effective January 1, 2024, via the CY2024 OPPS/ASC final rule (CMS-1786-FC); available in hospital outpatient departments, CMHCs, FQHCs, and RHCs [SOURCE: CMS CY2024 OPPS fact sheet]
- Opioid treatment programs (OTP) under the SUPPORT Act
- Psychiatric emergency care
- Mental health nursing home services (skilled nursing)

**Key MA difference**: MA can impose cost-sharing, prior authorization, and network restrictions on all of these. Original Medicare cannot impose PA on most outpatient BH services.

**As of 2022**: 98% of MA enrollees were in plans requiring prior authorization for some mental health/SUD services; 93% required PA for psychiatric hospital admissions; 91% for partial hospitalization. [SOURCE: KFF MA Mental Health Coverage analysis]

## MA Supplemental Behavioral Health Benefits

Supplemental-benefit authority for MA derives from the Balanced Budget Act of 1997 (Pub. L. 105-33), §§ 1851–1859 of the Social Security Act, codified at 42 CFR § 422.102. The BBA of 2018 (Pub. L. 115-123) further expanded this authority via SSA § 1852(a)(3) to allow non-primarily-health-related benefits for chronically ill enrollees. [SOURCE: Congress.gov BBA 1997; CMS SSBCI guidance April 2019]

BH-relevant supplemental benefits:

| Benefit | Description |
|---|---|
| In-home mental health support | Check-ins, safety monitoring, medication adherence support at home |
| Caregiver support | Respite services, caregiver training (especially for dementia caregivers) |
| Social isolation / loneliness programs | Structured social activities, companion programs |
| Telehealth mental health visits | Beyond what Original Medicare covers |
| Peer support services | Some MA plans offer peer support as supplemental benefit |
| Transportation to mental health appointments | Non-emergency medical transportation |
| Vision, hearing, dental (indirect BH) | Sensory loss linked to depression; dental pain affects mental health |
| Meals post-hospital | Addresses nutrition for BH recovery |

**Primarily Health-Related (PHR) Supplemental Benefits**: under 2019 CMS guidance (2019 MA Call Letter), MA plans expanded the PHR definition to include items that diagnose, prevent, or treat illness, compensate for physical impairments, or ameliorate functional/psychological impact — enabling housing navigation, food assistance, and transportation for enrollees with BH conditions.

**Special Supplemental Benefits for Chronically Ill (SSBCI)**: since 2020 (contract year 2020), MA can offer non-primarily health-related benefits for enrollees with specified chronic conditions, authorized under BBA 2018 amendment to SSA § 1852(a)(3), codified at 42 CFR § 422.102(f). SMI and SUD qualify. Eligibility requires an objective determination (health risk assessment or claims review — not self-attestation) that the enrollee is chronically ill (life-threatening or significantly limiting comorbid condition + high hospitalization risk + requires intensive care coordination). Qualifying benefits include:
- Food and produce
- General supports for living (housing, utilities assistance)
- Pest control
- Transportation for non-medical needs
[SOURCE: CMS SSBCI HPMS memo April 2019; Federal Register 2020-02085]

## Prior Authorization for BH in Medicare Advantage

Two major rules govern MA PA requirements for BH:

### CMS-4201-F (CY2024 MA Final Rule — effective plan year 2024)
The primary rule with substantive MA BH PA requirements. [SOURCE: Federal Register 2023-07115]
- **Emergency BH services may not be subject to prior authorization** — codified at 42 CFR 422.112
- MA organizations must establish a **utilization management (UM) committee** that annually reviews and approves all utilization management policies, including BH PA criteria
- Coverage criteria must align with national coverage determinations (NCDs), local coverage determinations (LCDs), and general Traditional Medicare coverage conditions — plans may not apply more restrictive criteria
- MA organizations must notify enrollees when a BH provider is dropped mid-year from the network
- CMS began adding BH facility types (outpatient behavioral health) as network adequacy specialty types requiring plan compliance

### CMS-0057-F (Interoperability and Prior Authorization Final Rule — published Feb 8, 2024)
Primarily an interoperability/FHIR API rule; key PA provisions: [SOURCE: Federal Register 2024-00895]
- Impacted payers must send **expedited (urgent) PA decisions within 72 hours**; **standard (non-urgent) PA decisions within 7 calendar days** — applies to all services including BH (excluding drugs)
- Requires payers to implement FHIR-based Prior Authorization API to share PA data with providers and patients
- Plans must track and report PA approval and denial rates
- **Compliance timeline**: operational provisions January 1, 2026; API requirements January 1, 2027

### Gold-Carding
Gold-carding (exempting high-approval-rate providers from PA) is a **voluntary insurer practice**, not mandated by federal rule. UnitedHealth Group launched a national gold card program in 2023–2024 exempting certain providers. Some plans have applied this to psychiatrists. [SOURCE: CMS CY2025 MA Final Rule fact sheet]

**BH PA problem in MA**: OIG's April 2022 report (OEI-09-18-00260) found that MA organizations denied medically necessary care at concerning rates, with some denials meeting Medicare coverage rules but failing to meet plan-added criteria. While not a BH-specific audit, BH services are among the most PA-intensive in MA plans.

## Network Adequacy for BH in MA

Governed by **42 CFR 422.116** (time/distance standards; four BH specialty types) and **42 CFR 422.112** (appointment wait-time standards). [SOURCE: eCFR 422.116; Federal Register 2023-07115]

**Four BH specialty types subject to network adequacy evaluation** (effective contract year 2024):
1. Psychiatry
2. Clinical psychology
3. Licensed clinical social work
4. Inpatient psychiatric facility services

**Appointment wait-time standards** (42 CFR 422.112; effective contract year 2024):
- Emergency/urgently needed services: **immediately**
- Non-emergency services requiring medical attention: **within 7 business days**
- Routine and preventive care: **within 30 business days**

**Telehealth credit**: MA plans receive a **10 percentage-point credit** toward time/distance standards for clinical psychology and clinical social work when they provide additional telehealth benefits from those specialty types. For contract year 2025+, CMS added "Outpatient Behavioral Health" as an additional facility specialty eligible for this 10-point credit. [SOURCE: CMS CY2024 MA Final Rule; CMS CY2025 MA Final Rule]

**In-network vs. out-of-network**: MA HMOs have no OON coverage for non-emergency BH; MA PPOs have OON at higher cost-sharing.

## D-SNP Behavioral Health Coordination

Dual Eligible Special Needs Plans (D-SNPs) serve Medicare-Medicaid dual eligibles, who have extremely high rates of SMI and SUD:
- D-SNPs required to have care coordination structures addressing both Medicare and Medicaid BH
- Highly Integrated D-SNPs (HIDE D-SNPs) and Fully Integrated D-SNPs (FIDE D-SNPs) have deeper integration
- **FIDE D-SNPs**: receive capitation from both Medicare and Medicaid; cover LTSS, BH, and acute services in one plan — best positioned for whole-person BH care
- FIDE D-SNPs required to have interdisciplinary care team including BH professional

## C-SNPs for Chronic Mental Illness

Chronic Condition Special Needs Plans (C-SNPs) can specialize in serving beneficiaries with:
- Chronic mental illness (specifically schizophrenia, bipolar disorder, MDD)
- SUD (less common)
C-SNPs must have: specialized BH expertise, targeted BH care management, network with adequate BH providers
Less common than D-SNPs; some MA markets have SMI-specific C-SNPs.

## MA vs. Original Medicare + Medigap for BH

| Factor | Original Medicare + Medigap | Medicare Advantage |
|---|---|---|
| Network | Any Medicare provider; no network restrictions | In-network required (HMO) or higher OON cost (PPO) |
| PA requirements | Minimal PA on outpatient BH | PA for inpatient psych, PHP, IOP commonly required (98% of enrollees face some BH PA) |
| Cost predictability | Medigap covers most cost-sharing; predictable | Lower premium but PA denials, OOP risk |
| Supplemental BH benefits | None | May include mental health support programs |
| Telehealth | CY2024+ expansions | More telehealth flexibility; 98% of MA enrollees had access (2022) |
| OTP/SUD | OTP covered under SUPPORT Act | OTP covered but PA may apply (85% of MA enrollees require PA for OTP) |
| Best for | Complex BH needs requiring specialist access; frequent hospitalizations | Relatively stable BH needs; enrolled in integrated D-SNP |

## Mental Health Parity in Medicare Advantage

**MHPAEA does not apply to Medicare or Medicare Advantage.** The Mental Health Parity and Addiction Equity Act of 2008 (MHPAEA) governs commercial insurance markets (individual, small/large group), employer-sponsored plans, and Medicaid managed care — but not Medicare programs. [SOURCE: CMS MHPAEA fact sheet]

MA beneficiaries have a more limited cost-sharing parity protection: the Medicare Improvements for Patients and Providers Act of 2008 (MIPPA) equalized cost-sharing for outpatient mental health services with outpatient medical/surgical services under Part B, ending the historical 50% cost-sharing for MH. This applies through MA.

**Parity-adjacent protections under CMS rules** (not MHPAEA):
- **CMS-4201-F (CY2024)**: MA organizations must establish care coordination programs including BH services, explicitly framed as advancing "parity between behavioral health and physical health services and whole-person care." Coverage criteria must align with Traditional Medicare NCDs/LCDs — plans cannot apply more restrictive criteria to BH than to medical/surgical services.
- **PA non-discrimination**: CMS has used its MA oversight authority to ensure BH services are not systematically disadvantaged relative to other services, including the emergency BH PA prohibition.

**Practical implication**: MA enrollees with BH needs cannot invoke MHPAEA rights. Their recourse is MA grievance/appeals, CMS oversight, and state insurance regulation where applicable — a weaker enforcement posture than commercial insurance parity.

## BH Carve-Out Within MA

Some MA organizations carve out BH to a behavioral health subcontractor:
- Administrative advantage: specialized BH utilization management
- Care integration risk: fragmented coordination between physical health MA and BH carve-out
- CMS contract requirements: MA organization remains accountable regardless of subcontracting; must ensure BH subcontractor meets all MA requirements

## References

1. CMS. **CY2024 MA and Part D Final Rule (CMS-4201-F) Fact Sheet** — BH network adequacy, UM committee, emergency BH PA prohibition.
   https://www.cms.gov/newsroom/fact-sheets/2024-medicare-advantage-and-part-d-final-rule-cms-4201-f

2. Federal Register. **CY2024 MA Final Rule (2023-07115)** — 42 CFR 422.112/422.116 wait-time and network adequacy standards for BH.
   https://www.federalregister.gov/documents/2023/04/12/2023-07115/medicare-program-contract-year-2024-policy-and-technical-changes-to-the-medicare-advantage-program

3. CMS. **CMS-0057-F Interoperability and Prior Authorization Final Rule** — 72-hr expedited / 7-day standard PA decision timelines; FHIR API requirements; compliance dates Jan 2026/2027.
   https://www.cms.gov/newsroom/fact-sheets/cms-interoperability-prior-authorization-final-rule-cms-0057-f

4. Federal Register. **CMS-0057-F (2024-00895)** — full rule text.
   https://www.federalregister.gov/documents/2024/02/08/2024-00895/medicare-and-medicaid-programs-patient-protection-and-affordable-care-act-advancing-interoperability

5. CMS. **CY2025 MA and Part D Final Rule (CMS-4205-F) Fact Sheet** — outpatient BH added as network adequacy specialty; gold-carding context.
   https://www.cms.gov/newsroom/fact-sheets/contract-year-2025-medicare-advantage-and-part-d-final-rule-cms-4205-f

6. CMS. **CY2024 OPPS/ASC Final Rule (CMS-1786-FC)** — IOP added as Original Medicare-covered service effective Jan 1, 2024.
   https://www.cms.gov/newsroom/fact-sheets/cy-2024-medicare-hospital-outpatient-prospective-payment-system-ambulatory-surgical-center-payment

7. CMS. **SSBCI HPMS Guidance Memo (April 2019)** — special supplemental benefits for chronically ill; SMI/SUD qualifying conditions; eligibility determination requirements.
   https://www.cms.gov/medicare/health-plans/healthplansgeninfo/downloads/supplemental_benefits_chronically_ill_hpms_042419.pdf

8. Federal Register. **CY2021 MA Final Rule (2020-02085)** — SSBCI implementation details, 42 CFR 422.102(f).
   https://www.federalregister.gov/documents/2020/02/18/2020-02085/medicare-and-medicaid-programs-contract-year-2021-and-2022-policy-and-technical-changes-to-the

9. KFF. **Mental Health and Substance Use Disorder Coverage in Medicare Advantage Plans** — PA rates (98%/93%/91%/85%), supplemental benefit uptake (6%), OON access stats.
   https://www.kff.org/mental-health/mental-health-and-substance-use-disorder-coverage-in-medicare-advantage-plans/

10. OIG HHS. **Some Medicare Advantage Organization Denials of Prior Authorization Requests Raise Concerns About Beneficiary Access to Medically Necessary Care (OEI-09-18-00260, April 2022)**
    https://oig.hhs.gov/reports/all/2022/some-medicare-advantage-organization-denials-of-prior-authorization-requests-raise-concerns-about-beneficiary-access-to-medically-necessary-care/

11. Medicare.gov. **Inpatient Mental Health Care Coverage** — 190-day lifetime limit applicability to freestanding psychiatric hospitals.
    https://www.medicare.gov/coverage/mental-health-care-inpatient

12. CMS. **Mental Health Parity and Addiction Equity Act (MHPAEA) Fact Sheet** — scope of MHPAEA (commercial/employer/Medicaid); Medicare/MA explicitly excluded.
    https://www.cms.gov/newsroom/fact-sheets/mental-health-parity-and-addiction-equity-act-2008-mhpaea
