---
name: medicare-mental-health-sud-coverage
description: >-
  Medicare coverage for mental health and substance use disorder (SUD) treatment: Part A inpatient psychiatric (190-day lifetime limit), Part B outpatient psychotherapy cost-sharing, partial hospitalization programs (PHP), intensive outpatient programs (IOP), opioid treatment programs (OTP) under the SUPPORT Act, telehealth mental health rules, newly covered provider types (LMHCs, MFTs), parity gap (Medicare exempt from MHPAEA), Medicare Advantage prior authorization. TRIGGER: Medicare mental health coverage, Medicare psychiatry, Medicare psychotherapy cost sharing, Medicare inpatient psychiatric 190-day limit, Medicare partial hospitalization program PHP, Medicare intensive outpatient IOP, Medicare OTP opioid treatment program, SUPPORT Act Medicare, methadone Medicare Part B, buprenorphine Medicare OTP, Medicare SUD treatment, Medicare parity MHPAEA, Medicare addiction counselor coverage, Medicare telehealth mental health. SKIP: Medicaid mental health → medicaid-rehabilitative-services; MHPAEA enforcement private insurance → mhpaea-enforcement; CCBHC model → medicaid-imd-ccbhc; NC SUD law → nc-drug-rehabilitation-treatment; Medicare SNF/IRF → medicare-snf-irf-coverage.
version: "1.1.0"
updated: "2026-06-22"
category: health/coverage-rehabilitation
tags: [medicare, mental-health, SUD, OTP, PHP, IOP, SUPPORT-act, parity, part-b]
keywords:
  - Medicare mental health
  - Medicare psychotherapy
  - inpatient psychiatric 190-day limit
  - partial hospitalization program
  - intensive outpatient program
  - opioid treatment program Medicare
  - SUPPORT Act Medicare OTP
  - Medicare SUD coverage
  - Medicare parity gap
  - MHPAEA Medicare exempt
  - Medicare addiction counselor
  - Medicare telehealth mental health
  - methadone Medicare
  - buprenorphine OTP Medicare
whenToUse:
  - "Does Medicare cover mental health therapy or psychiatry?"
  - "Medicare inpatient psychiatric hospitalization limits"
  - "Medicare partial hospitalization program coverage"
  - "Does Medicare cover intensive outpatient for mental health?"
  - "Medicare coverage for opioid treatment programs / methadone clinic"
  - "SUPPORT Act and Medicare OTP benefit"
  - "Is Medicare subject to MHPAEA parity?"
  - "Medicare telehealth mental health requirements"
  - "Licensed mental health counselors billing Medicare"
verified-as-of: "2026-06-22"
---

# Medicare Mental Health and SUD Coverage

Medicare covers mental health and substance use disorder (SUD) treatment across multiple benefit components, but with significant restrictions not found in its medical/surgical coverage — most notably the 190-day inpatient psychiatric lifetime limit and exemption from federal parity law.

---

## Medicare's Parity Gap

**Medicare is NOT subject to MHPAEA** (Mental Health Parity and Addiction Equity Act of 2008). [^1]

The law explicitly excludes Medicare, meaning:
- Medicare's utilization management, coverage exclusions, and reimbursement rates for mental health and SUD can be more restrictive than for medical/surgical treatment.
- The 50% outpatient mental health coinsurance (before 2010) was a direct expression of this gap; it was partially corrected by MIPPA 2008.
- Advocates continue to push for legislation extending MHPAEA to Medicare.
- **Medicare Advantage plans** are also not subject to MHPAEA, though they cannot be less generous than Original Medicare.

---

## Part A: Inpatient Psychiatric Coverage

### General Hospitals vs. Freestanding Psychiatric Hospitals [^1][^2]

| Setting | Benefit | Lifetime Limit |
|---------|---------|----------------|
| General acute hospital (psychiatric unit) | Standard Part A benefit period | No lifetime limit |
| Freestanding psychiatric hospital | Part A benefit period rules | **190-day lifetime limit** |

**The 190-day limit**: Medicare Part A covers up to 190 days of inpatient psychiatric care in a freestanding psychiatric hospital during a beneficiary's *lifetime*. This is Medicare's only lifetime inpatient limit — all other inpatient benefits reset by benefit period. It does not apply to psychiatric care in general hospitals.

**Cost sharing (2026)**:
- Days 1–60: $1,736 deductible per benefit period
- Days 61–90: $434/day coinsurance
- Days 91–150 (lifetime reserve): $868/day
- Beyond: 100% patient responsibility

Services covered during inpatient psychiatric stay: psychiatric evaluation, medication management, individual therapy, group therapy, milieu therapy, social work.

---

## Part B: Outpatient Mental Health Coverage

### Cost Sharing [^1][^2]

**Current (2010–present)**: After the Part B deductible ($257 in 2025; $283 in 2026; indexed annually), Medicare pays **80%** and the beneficiary pays **20%** coinsurance — same as other Part B services. [verified-as-of: 2026-06-22, source: CMS 2026 Medicare Parts A & B Premiums and Deductibles]

**Historical note**: Before 2010, Medicare charged 50% coinsurance for outpatient mental health (vs. 20% for other Part B). MIPPA 2008 phased this in over 2010–2014 to achieve cost-sharing parity.

### Covered Services [^2][^3]

- Individual and group psychotherapy
- Psychiatric diagnostic evaluations
- Medication management (by psychiatrists and prescribing practitioners)
- Psychological testing
- Family counseling (when primary purpose is beneficiary's mental health treatment)
- Crisis intervention services

### Covered Provider Types [^1][^2]

As of **January 1, 2024** (expansion from Consolidated Appropriations Act, 2023 / CY 2024 Physician Fee Schedule Final Rule): [^2a]

| Provider | Direct Medicare billing | Rate |
|----------|------------------------|------|
| Psychiatrists (MD/DO) | Yes | 100% PFS |
| Psychologists (PhD/PsyD) | Yes | 100% PFS |
| Clinical social workers (LCSW) | Yes | 75% of psychologist rate |
| Nurse practitioners, clinical nurse specialists, PAs | Yes | 85% of PFS |
| **Licensed Mental Health Counselors (MHC)** | **Yes (new 2024)** | **75% of psychologist rate** |
| **Marriage and Family Therapists (MFT)** | **Yes (new 2024)** | **75% of psychologist rate** |

Newly added MHC and MFT providers must hold a master's or doctoral degree in their field, complete at least 2 years or 3,000 hours of post-master's supervised clinical experience, and be licensed/certified by the state in which they practice. [verified-as-of: 2026-06-22, source: CMS MFT/MHC FAQ]

**Note**: The skill previously listed "Addiction counselors (where state-licensed)" as a new 2024 provider type — this is not accurate as a distinct Medicare billing category. Addiction counselors providing services within a Medicare-enrolled OTP bill under the OTP bundled payment, not as independent Part B providers.

**Provider access gap**: Only ~60% of psychiatrists accept new Medicare patients (vs. 81% of general practitioners); ~7.5% have opted out of Medicare entirely — affecting access despite coverage expansion.

### Telehealth Mental Health [^1][^4]

**Current status (2026)**: Pandemic-era telehealth flexibilities were extended through **December 31, 2027** by the Consolidated Appropriations Act, 2026. The in-person visit requirement is **not yet in effect**. [verified-as-of: 2026-06-22, source: CMS Telehealth FAQ updated 2/26/2026]

**Upcoming in-person requirement (effective after December 31, 2027)**:
- **Initial in-person visit** required within 6 months before initiating telehealth mental health services for new patients.
- **Periodic in-person visits** required every 12 months thereafter to continue telehealth mental health coverage.
- Beneficiaries who began receiving mental health services on or before December 31, 2027 are considered "established" and subject only to the annual (not initial 6-month) in-person requirement after that date.

**Geographic and originating site restrictions**:
- Permanently removed for SUD telehealth services effective **July 1, 2019** (SUPPORT Act, Sec. 2001). [^4]
- Permanently removed for general mental health telehealth services by the **Consolidated Appropriations Act, 2021**.

Video + audio required; audio-only may qualify in limited circumstances.

---

## Partial Hospitalization Programs (PHP)

### Coverage and Eligibility [^3][^5]

**Coverage**: Medicare Part B covers PHPs when ordered by a psychiatrist or physician trained in psychiatric diagnosis and treatment.

**Level of care requirement**: minimum ~20 hours of therapeutic services per week (structured programming 4–8 hours/day).

**Two eligibility pathways**:
1. Discharged from inpatient psychiatric hospitalization AND PHP substitutes for continued inpatient care.
2. Patient would be at reasonable risk of requiring inpatient hospitalization without PHP.

**Cost sharing**: Part B rules — 20% coinsurance after deductible.

**Payment structure**: CMS pays two PHP APCs per provider type:
- APC for days with 3 services/day
- APC for days with ≥4 services/day

**Settings**: Hospital-based or Community Mental Health Center (CMHC)-based PHPs.

---

## Intensive Outpatient Programs (IOP)

### Coverage (New in 2024) [^1][^4a][^5]

Medicare added IOP as a distinct covered benefit beginning **January 1, 2024**.

**Level of care threshold**: minimum **9 hours/week** of intensive outpatient mental health services.

**Physician certification**: required every other month confirming continued IOP level of care is appropriate.

**Cost sharing**: Part B — 20% coinsurance after deductible.

IOP fills the critical gap between PHP (20+ hours/week) and standard outpatient psychotherapy (1 hour/week), enabling step-down care from PHP without return to standard office visits.

---

## Opioid Treatment Programs (OTP): Medicare Part B Coverage

### SUPPORT Act — Medicare OTP Benefit [^6][^7]

The **Substance Use Disorder Prevention that Promotes Opioid Recovery and Treatment for Patients and Communities Act (SUPPORT Act, 2018)** — Section 2005 — created the Medicare OTP benefit:

- **Effective**: January 1, 2020.
- Medicare Part B now covers OUD treatment services provided by Medicare-enrolled OTPs.
- **First time** methadone (oral) became covered under Medicare Part B (methadone is not covered under Part D for OUD).

### Covered OTP Services [^6][^7]

- MOUD medications: **methadone** (oral), **buprenorphine** (oral, injectable, implantable), **naltrexone** (injectable), naloxone, nalmefene
- SUD counseling (individual and group)
- Individual and group psychotherapy
- Toxicology testing
- Periodic assessments (in-person and telehealth)

### Bundled Payment Structure [^6][^7]

Medicare pays OTPs using a **bundled prospective payment** (per episode/week), not fee-for-service:

| Bundle component | Basis |
|-----------------|-------|
| Rate variation | Based on medication type |
| Included services | Medication administration, counseling, assessments, toxicology |
| Telehealth | Permitted for counseling and assessments |

**Bundle types by medication**:
- Methadone (oral) — highest bundle (requires daily dispensing)
- Buprenorphine (oral) — lower bundle (can be prescribed to take home)
- Buprenorphine (injectable/implantable) — specific bundle rates
- Naltrexone (injectable) — separate bundle rate

### OTP Cost Sharing [^6]

- **No copayment** for OTP services — the beneficiary copayment is waived (zero) for all HCPCS codes billed by OTPs. [verified-as-of: 2026-06-22, source: CMS OTP Billing & Payment]
- **Part B deductible does apply** — as mandated for all Part B services by section 1833(b) of the Act; it is not waived for OTP.
- For dually eligible beneficiaries (Medicare + Medicaid), state Medicaid agencies are liable for the Part B deductible, subject to applicable limits.
- The zero-copayment design was intended to remove financial barriers to MOUD treatment access.

### Office-Based OUD Treatment (Non-OTP) [^6]

Separate from the OTP bundle:
- Buprenorphine prescribed in office-based settings (not OTPs) is covered under Medicare Part D for the prescription, and the prescribing visit under Part B.
- Naloxone covered under Part D.
- Behavioral health components billed under standard Part B outpatient mental health.

---

## Medicare Advantage (Part C) Differences [^1]

MA plans impose significant additional utilization management on mental health and SUD:
- **Prior authorization**: ~94% of MA plans require it for inpatient psychiatric stays; ~85% for mental health specialty services.
- **Network adequacy**: Less than one-quarter of county psychiatrists typically included in MA networks.
- **Daily copayments**: Some plans charge copays starting on Day 1 of inpatient psychiatric stay.
- **IOP and PHP**: Coverage varies by plan; prior auth commonly required.
- OTP: MA plans must cover OTP services but may add administrative requirements.

---

## Part D: Mental Health and SUD Medications [^1][^2]

- Most antidepressants, antipsychotics, anxiolytics covered under Part D.
- **2025 change**: Part D out-of-pocket cap of **$2,000/year** — significant for high-cost psychiatric medications.
- Prior authorization and step therapy requirements common for new psychiatric drug initiations.
- Out-of-pocket before cap: <$10 for generics; $40–$100+ or 40–50% coinsurance for brand-name drugs.
- Methadone for OUD: NOT covered under Part D (covered only through OTP bundled payment under Part B).

---

## Coverage Gaps and Remaining Issues

| Gap | Current Status |
|-----|---------------|
| 190-day psychiatric hospital lifetime limit | Remains in effect; proposals to repeal have not passed |
| MHPAEA parity | Medicare explicitly excluded; legislative fix needed |
| Psychiatrist shortage / opt-outs | Coverage exists; access remains restricted |
| Residential SUD treatment | Generally not covered under Original Medicare |
| Peer support specialists | Not covered under Medicare FFS |
| Recovery support services | Not covered |

---

## Anti-Patterns / Common Mistakes

- **Assuming MHPAEA applies to Medicare**: it does not; coverage disparities between mental health and medical/surgical are legally permissible under Medicare.
- **190-day limit applies to all psych inpatient**: it only applies to freestanding psychiatric hospitals, not general hospital psychiatric units.
- **Methadone under Part D**: Part D does not cover methadone for OUD; only Part B OTP bundled payment covers it.
- **Assuming LMHCs/MFTs couldn't bill Medicare before 2024**: this changed January 1, 2024.
- **No IOP coverage before 2024**: IOP was added as a new benefit in 2024; prior gap is a common knowledge artifact.
- **Telehealth mental health without in-person visit**: the in-person requirement is **not yet in effect** as of 2026 — telehealth flexibilities were extended through December 31, 2027 by the Consolidated Appropriations Act, 2026. The requirement will take effect after December 31, 2027.

---

## References

[^1]: KFF. "FAQs on Mental Health and Substance Use Disorder Coverage in Medicare." https://www.kff.org/mental-health/issue-brief/faqs-on-mental-health-and-substance-use-disorder-coverage-in-medicare/ (verified-as-of: 2026-06-22) — authoritative overview of parity gap, 190-day limit, cost sharing, OTP, IOP, telehealth, provider expansion, MA prior auth statistics.

[^2]: CMS. "2026 Medicare Parts A & B Premiums and Deductibles." https://www.cms.gov/newsroom/fact-sheets/2026-medicare-parts-b-premiums-deductibles (verified-as-of: 2026-06-22) — Part A deductible $1,736; days 61–90 coinsurance $434; lifetime reserve $868/day; Part B deductible $283 in 2026.

[^2a]: CMS. "Marriage and Family Therapists (MFTs) and Mental Health Counselors (MHCs)." https://www.cms.gov/medicare/payment/fee-schedules/physician-fee-schedule/marriage-family-therapists-mental-health-counselors (verified-as-of: 2026-06-22) — credential requirements; billing at 75% of clinical psychologist rate; effective January 1, 2024.

[^3]: Medicare.gov. "Mental Health Care (Outpatient)." https://www.medicare.gov/coverage/mental-health-care-outpatient; "Inpatient Mental Health Care." https://www.medicare.gov/coverage/mental-health-care-inpatient; "Partial Hospitalization." https://www.medicare.gov/coverage/mental-health-care-partial-hospitalization (verified-as-of: 2026-06-22).

[^4]: CMS. "Medicare Telehealth FAQ (updated 2/26/2026)." https://www.cms.gov/files/document/telehealth-faq-updated-02-26-2026.pdf (verified-as-of: 2026-06-22) — in-person requirement delayed to after December 31, 2027 by Consolidated Appropriations Act, 2026; SUPPORT Act Sec. 2001 removed SUD geographic restrictions as of July 1, 2019; Consolidated Appropriations Act, 2021 permanently removed geographic restrictions for general mental health telehealth.

[^4a]: Medicare.gov. "Intensive Outpatient Program Services." https://www.medicare.gov/coverage/mental-health-care-intensive-outpatient-program-services (verified-as-of: 2026-06-22) — 9 hours/week minimum; physician confirmation every other month; effective January 1, 2024.

[^5]: CMS. "CY 2024 Medicare Hospital Outpatient Prospective Payment System and Ambulatory Surgical Center Payment System Final Rule (CMS 1786-FC)." https://www.cms.gov/newsroom/fact-sheets/cy-2024-medicare-hospital-outpatient-prospective-payment-system-and-ambulatory-surgical-center-0 (verified-as-of: 2026-06-22) — IOP established as new Medicare benefit effective January 1, 2024; settings include hospital outpatient departments, CMHCs, FQHCs, RHCs.

[^6]: CMS. "Opioid Treatment Programs (OTP)." https://www.cms.gov/medicare/payment/opioid-treatment-program (verified-as-of: 2026-06-22) — SUPPORT Act Sec. 2005 created OTP benefit effective January 1, 2020; zero copayment; Part B deductible applies per §1833(b); bundled payment by medication type.

[^7]: CMS. "Medicare Benefit Policy Manual Chapter 17 — Opioid Treatment Programs (OTPs)." https://www.cms.gov/files/document/chapter-17-opioid-treatment-programs-otps.pdf (verified-as-of: 2026-06-22) — coverage criteria, bundled payment structure, cost-sharing rules for OTP services.

[^8]: Congress.gov. "H.R.6 — SUPPORT for Patients and Communities Act, 115th Congress." https://www.congress.gov/bill/115th-congress/house-bill/6 (Pub. L. 115-271, signed October 24, 2018) — Sec. 2005 created Medicare OTP Part B benefit; Sec. 2001 removed SUD telehealth geographic restrictions.

[^9]: CMS. "Calendar Year (CY) 2024 Medicare Physician Fee Schedule Final Rule." https://www.cms.gov/newsroom/fact-sheets/calendar-year-cy-2024-medicare-physician-fee-schedule-final-rule (verified-as-of: 2026-06-22) — MFT/MHC Medicare enrollment and billing; IOP benefit implementation.

[^10]: KFF. "What to Know About Medicare Coverage of Telehealth." https://www.kff.org/medicare/what-to-know-about-medicare-coverage-of-telehealth/ (verified-as-of: 2026-06-22) — telehealth extension history, CAA 2026 extension through 2027, mental health in-person requirement delayed.
