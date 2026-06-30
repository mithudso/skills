---
name: medicare-outpatient-rehabilitation
description: >-
  Medicare Part B coverage for outpatient rehabilitation services — physical therapy (PT), occupational therapy (OT), and speech-language pathology (SLP): eligibility, coverage settings, 80/20 cost-sharing, annual therapy thresholds and the KX modifier system, targeted medical review, documentation standards, Advanced Beneficiary Notice (ABN) obligations, and how Medicare Advantage differs. TRIGGER: Medicare outpatient physical therapy coverage, Medicare OT coverage, Medicare speech therapy, Part B therapy cap, KX modifier, therapy threshold 2025 2026, targeted medical review therapy, outpatient rehab Medicare eligibility, ABN therapy, GP GO GN modifiers. SKIP: SNF or IRF inpatient rehab → medicare-snf-irf-coverage; home health benefit → medicare-home-health; Medicare mental health outpatient → medicare-mental-health-sud-coverage; Medicaid outpatient therapy → medicaid-rehabilitative-services.
version: "1.1.0"
updated: "2026-06-22"
category: health/coverage-rehabilitation
tags: [medicare, outpatient-rehab, physical-therapy, occupational-therapy, speech-therapy, part-b, KX-modifier]
keywords:
  - Medicare physical therapy
  - Medicare occupational therapy
  - Medicare speech-language pathology
  - therapy threshold
  - KX modifier
  - targeted medical review
  - outpatient rehabilitation Part B
  - 20% coinsurance therapy
  - Advanced Beneficiary Notice
  - GP GO GN modifiers
  - Medicare Advantage therapy authorization
whenToUse:
  - "Does Medicare cover physical therapy?"
  - "What is the Medicare therapy cap or threshold?"
  - "KX modifier when and how to use"
  - "Medicare OT or speech-language pathology coverage"
  - "Targeted medical review for therapy services"
  - "Patient responsibility for outpatient therapy beyond threshold"
  - "GP GO GN therapy billing modifiers"
verified-as-of: "2026-06-22"
---

# Medicare Outpatient Rehabilitation Coverage (Part B)

Medicare Part B covers medically necessary outpatient physical therapy (PT), occupational therapy (OT), and speech-language pathology (SLP) services. Coverage is structured around an annual threshold system with documentation-based exceptions, not a hard cap that terminates coverage.

---

## Eligibility and Coverage Conditions

**Basic requirements** (all must be met): [^1][^2]

1. **Medicare Part B enrollment** — outpatient therapy is a Part B benefit.
2. **Physician or qualified practitioner order** — a physician, nurse practitioner, clinical nurse specialist, or physician assistant must certify medical necessity; for dates of service on or after January 1, 2025, a signed and dated referral/order from a qualifying practitioner can satisfy the initial certification requirement. (42 CFR §§ 410.60, 424.24)
3. **Skilled care standard** — the services must require the skills of a licensed therapist; custodial care alone does not qualify under traditional Medicare.
4. **Medically necessary** — services must be reasonable and necessary for diagnosis or treatment of a condition.

**Post-Jimmo Maintenance Standard**: Following the *Jimmo v. Sebelius* settlement (January 2013), Medicare covers skilled maintenance therapy in outpatient settings when the skills of a qualified therapist are needed to perform or supervise services safely and effectively — the maintenance standard does **not** require that functional improvement is expected. Maintenance of function or prevention/slowing of deterioration is sufficient, provided all other coverage criteria are met. This standard applies in SNF, home health, and outpatient therapy settings.[^3]

---

## Coverage Settings

Medicare Part B outpatient therapy is covered in: [^1][^2]

- Hospital outpatient departments
- Comprehensive Outpatient Rehabilitation Facilities (CORFs) — regulated under 42 CFR Part 485 Subpart B, must provide physician services, PT, and social/psychological services at a single fixed location
- Rehabilitation agencies (distinct from CORFs; also Medicare-certified outpatient therapy providers)
- Private therapy practices
- Skilled nursing facilities (for patients NOT in a Part A covered stay)
- Home settings (when the patient is not under a Part A home health episode)

---

## Cost Sharing

| Component | Amount (2026) |
|-----------|---------------|
| Part B deductible (annual) | $283 |
| Coinsurance | 20% of Medicare-approved amount |
| Medicare pays | 80% of approved amount |

The 2026 annual Part B deductible increased from $257 (2025) to $283, as published in the CMS November 2025 fact sheet and in the Federal Register (Nov. 19, 2025).[^6]

The 20% coinsurance applies per service visit after the annual deductible is met. There is no per-visit copay cap under Original Medicare.

---

## Annual Therapy Thresholds and the KX Modifier System

### Statutory Basis [^4][^5]

The current threshold-plus-KX-modifier system was created by **Section 50202 of the Bipartisan Budget Act of 2018 (BBA 2018)**, codified as Section 1833(g)(7) of the Social Security Act. It replaced the former hard therapy caps. The targeted medical review process was established under the Medicare Access and CHIP Reauthorization Act of 2015 (MACRA).

### Threshold Amounts [^4][^5]

| Service Group | 2025 Threshold | 2026 Threshold |
|---------------|---------------|---------------|
| PT + SLP combined | $2,410 | $2,480 |
| OT (separate) | $2,410 | $2,480 |

Thresholds are indexed annually by the Medicare Economic Index (MEI) and reset January 1 each calendar year.

**These are NOT hard caps** — Medicare will continue to pay above the threshold if the KX modifier is appended.

### The KX Modifier [^4][^5]

When therapy charges exceed the annual threshold, providers must append the **KX modifier** to claims. It serves as the provider's attestation that:
- Services are medically necessary
- They require therapist-level skills
- Documentation in the medical record justifies continued treatment

> **Important**: The KX modifier does **not** require attestation that the patient will achieve functional improvement. Attestation is that services are medically necessary and require skilled care — consistent with the *Jimmo* maintenance standard. (CMS: "confirmation that services are medically necessary as justified by appropriate documentation in the medical record.")

**Requirements**:
- No special documentation must be submitted with the claim itself.
- The provider is responsible for ensuring adequate support exists in the patient's medical record.
- Documentation must include: baseline functional measurements, measurable goals, regular progress notes, and objective outcome measures.

**Without the KX modifier**, claims above the threshold are denied automatically.

### Targeted Medical Review Threshold [^5]

| Level | Amount | Notes |
|-------|--------|-------|
| KX modifier threshold | $2,480 (PT/SLP), $2,480 (OT) (2026) | Above = KX required |
| Targeted Medical Review (TMR) threshold | $3,000 | May trigger manual audit |

- When expenses exceed $3,000, claims may be selected for targeted medical review (Targeted Probe & Educate / TPE).
- Not all claims above $3,000 are reviewed — selection based on statistical deviation, provider history, geographic variation, and random selection.
- The $3,000 TMR threshold applies through 2028, after which it will be indexed annually to MEI. (BBA 2018 §50202)
- Providers facing review must submit: clinical rationale, objective outcome data, treatment modifications, and documentation that alternatives were considered.

---

## Billing Modifiers

These therapy-type modifiers remain required on outpatient therapy claims: [^1]

| Modifier | Service |
|----------|---------|
| GP | Physical therapy services |
| GO | Occupational therapy services |
| GN | Speech-language pathology services |

**G-code functional reporting (FLR) eliminated**: Effective January 1, 2019, CMS eliminated the functional limitation reporting requirement (nonpayable G-codes and severity modifiers), citing undue administrative burden. Functional outcome reporting is no longer required for Medicare outpatient rehabilitation billing.[^7]

---

## Advanced Beneficiary Notice (ABN)

When a provider anticipates Medicare may not cover a service (e.g., patient approaching threshold, questionable medical necessity), an **ABN must be issued before services are rendered**. [^1]

The ABN:
- Explains what service may not be covered and why
- Provides cost estimates
- Gives the patient the choice to receive the service and accept financial responsibility, or decline

Failure to issue an ABN before delivering a non-covered service removes the provider's ability to bill the patient for that service.

---

## Medicare Advantage (Part C) Differences [^1][^5]

Medicare Advantage (MA) plans:
- Must cover all services covered under Original Medicare (including outpatient therapy).
- Must apply the same therapy threshold amounts as Original Medicare ($2,480 for 2026) — plans cannot set lower thresholds.
- **May require prior authorization** for outpatient therapy services — CMS has issued rules shortening MA prior authorization response timeframes (CMS-4205-F / CMS-0057-F); verify current plan-specific timelines with the insurer.
- May use step therapy requirements (e.g., requiring lower-cost interventions first).
- Some plans have different network requirements for therapy providers.

**Action item**: Always verify MA plan-specific requirements; they vary significantly by plan and year.

---

## What Medicare Does NOT Cover (Outpatient Rehab Context)

- Custodial or maintenance care that does not require therapist skills (unless the *Jimmo* maintenance standard applies — skilled maintenance is covered).
- Therapy where no skilled need exists and no medically necessary goals remain — note: functional plateau alone does not disqualify coverage if skilled maintenance of the patient's condition requires therapist skills.
- Services that are not medically necessary per clinical guidelines.
- Personal comfort items or exercise equipment.
- Most gym memberships (though some MA plans include SilverSneakers as a supplemental benefit).

---

## Documentation Best Practices

For services near or above threshold: [^5]

1. **Initial evaluation** — baseline functional measurements with standardized outcome tools.
2. **Plan of care** — specific, measurable functional goals with target timelines.
3. **Progress notes** — objective measurements demonstrating response to treatment per visit or weekly.
4. **Recertification** — periodic review of goals and plan alignment.
5. **Discharge summary** — outcomes achieved vs. goals, home program.

Inadequate documentation is the most common driver of therapy claim denials on targeted review.

---

## Anti-Patterns / Common Mistakes

- **Treating $2,480 as a hard cap**: it is a threshold — covered beyond with KX modifier if medically necessary.
- **Not issuing ABN pre-service**: removes right to bill patient if Medicare denies.
- **Forgetting GP/GO/GN modifiers**: claims rejected without them.
- **Missing KX modifier above threshold**: automatic denial, not an exception process.
- **Assuming MA plans work like Original Medicare**: prior auth and step therapy requirements vary.
- **Conflating KX with improvement expectation**: KX attests to skilled need and medical necessity — improvement is not required (post-*Jimmo* maintenance therapy is covered).
- **Denying coverage at plateau**: functional plateau alone does not end coverage if skilled maintenance care is needed to prevent deterioration.
- **Applying improvement standard when maintenance is needed**: *Jimmo* standard permits skilled maintenance therapy coverage in outpatient settings.

---

> **Scope note**: PDPM (Patient-Driven Payment Model) and its impact on SNF therapy billing are outside this skill's scope — see `medicare-snf-irf-coverage` for SNF Part A therapy. This skill covers Part B outpatient therapy only.

---

## References

[^1]: CMS. "Therapy Services." https://www.cms.gov/medicare/coding-billing/therapy-services (verified-as-of: 2026-06-22)

[^2]: eCFR. "42 CFR § 410.60 — Outpatient physical therapy services: Conditions." https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-B/part-410/subpart-B/section-410.60 | "42 CFR § 424.24 — Requirements for medical and other health services furnished by providers under Medicare Part B." https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-B/part-424/subpart-B/section-424.24 | "42 CFR Part 485 Subpart B — Conditions of Participation: CORFs." https://www.ecfr.gov/current/title-42/chapter-IV/subchapter-G/part-485/subpart-B (verified-as-of: 2026-06-22)

[^3]: CMS. "Jimmo Settlement." https://www.cms.gov/medicare/settlements/jimmo | "FAQs Regarding Jimmo Settlement Agreement." https://www.cms.gov/center/special-topic/jimmo-settlement/faqs (maintenance standard applies in SNF, home health, and outpatient therapy settings; improvement not required) (verified-as-of: 2026-06-22)

[^4]: APTA. "Medicare Payment Thresholds for Outpatient Therapy Services." https://www.apta.org/your-practice/payment/medicare-payment/coding-billing/therapy-cap (2026 threshold: $2,480 PT/SLP combined, $2,480 OT; BBA 2018 §50202, SSA §1833(g)(7)) (verified-as-of: 2026-06-22)

[^5]: Noridian Medicare JF Part B. "Per-Beneficiary KX Modifier Thresholds." https://med.noridianmedicare.com/web/jfb/specialties/outpatient-therapy/per-beneficiary_kx_modifier_thresholds | CMS RAC. "0228 — Therapy Claims Billed with KX Modifier." https://www.cms.gov/data-research/monitoring-programs/medicare-fee-service-compliance-programs/medicare-fee-service-recovery-audit-program/approved-rac-topics/0228-therapy-claims-billed-kx-modifier-medical-necessity-documentation-requirements (TMR threshold $3,000 through 2028; BBA 2018 §50202; MACRA 2015) (verified-as-of: 2026-06-22)

[^6]: CMS. "2026 Medicare Parts A & B Premiums and Deductibles." https://www.cms.gov/newsroom/fact-sheets/2026-medicare-parts-b-premiums-deductibles | Federal Register. "Medicare Part B Monthly Actuarial Rates, Premium Rates, and Annual Deductible Beginning January 1, 2026" (Nov. 19, 2025). https://www.federalregister.gov/documents/2025/11/19/2025-20251/medicare-program-medicare-part-b-monthly-actuarial-rates-premium-rates-and-annual-deductible (2026 Part B deductible: $283; 2025: $257) (verified-as-of: 2026-06-22)

[^7]: CMS. "Functional Reporting." https://www.cms.gov/medicare/coding-billing/therapy-services/functional-reporting (G-code FLR eliminated effective January 1, 2019; CY 2019 Physician Fee Schedule Final Rule) (verified-as-of: 2026-06-22)
