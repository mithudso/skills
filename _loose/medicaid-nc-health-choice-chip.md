---
name: medicaid-nc-health-choice-chip
description: >-
  NC Health Choice — North Carolina's CHIP program — eligibility (ages 6–18, 133–210% FPL), benefits, cost-sharing, enrollment, relationship to Medicaid, and federal CHIP financing (enhanced FMAP). TRIGGER: NC Health Choice, CHIP North Carolina, Children's Health Insurance Program NC, NC Health Choice eligibility, NC Health Choice premium, NC Health Choice income limit, NC Health Choice 2025, CHIP enhanced match NC, NC child uninsured, NC Medicaid vs Health Choice. SKIP: full Medicaid eligibility for children → enrollment-eligibility-medicaid-eligibility-determination; MAGI methodology → enrollment-eligibility-magi-vs-non-magi-medicaid; NC Medicaid Managed Care → nc-health-services-medicaid-managed-care.
version: "1.0.0"
updated: "2026-06-22"
category: health
tags: [medicaid, CHIP, NC-Health-Choice, children, eligibility, North-Carolina]
related_skills:
  - enrollment-eligibility-medicaid-eligibility-determination
  - nc-health-services-medicaid-managed-care
  - medicaid-aca-expansion
  - medicaid-fmap-federal-financing
---

# NC Health Choice: North Carolina's CHIP Program

> **Educational only — NOT health or legal advice.** Income thresholds and benefits are subject to change; verify current eligibility and enrollment with NCDHHS or HealthCareConnect.nc.gov.

**verified-as-of: 2026-06-22**

---

## 1. What Is CHIP / NC Health Choice?

**CHIP (Children's Health Insurance Program)** is a federal-state program (Title XXI of the Social Security Act, enacted 1997) that covers uninsured children in families with incomes above Medicaid limits but too low to afford private insurance.

**NC Health Choice** is North Carolina's CHIP program, administered by NCDHHS Division of Health Benefits (DHB). NC operates CHIP as a **Medicaid expansion CHIP** — benefits are provided through the same Medicaid delivery system (Standard Plans), not a separate program.

---

## 2. Eligibility

### Age
- Children ages **6 through 18**
- Children under 6 are covered by Medicaid up to 133% FPL (not CHIP)

### Income (2026)
| FPL Range | Coverage |
|---|---|
| 0–133% FPL (children under 6) | Medicaid |
| 0–133% FPL (children 6–18) | Medicaid |
| 133%–210% FPL (children 6–18) | **NC Health Choice (CHIP)** |
| >210% FPL | Potentially marketplace; no public coverage |

> **Note**: NC's upper income limit has been 210% FPL for most categories. Federal law allows states to set CHIP upper limits up to 300% FPL.

### Citizenship/Residency
- Must be a US citizen or qualified immigrant
- NC resident
- Uninsured (access to affordable employer-sponsored insurance disqualifies, unless premium assistance pathway)

### Other Rules
- Cannot be covered by Medicaid (Medicaid covers up to 133% FPL)
- Parents/legal guardians apply on behalf of the child

---

## 3. Benefits

NC Health Choice provides **comprehensive child health coverage** similar to Medicaid:

- Well-child visits (EPSDT — Early and Periodic Screening, Diagnostic, and Treatment)
- Immunizations
- Primary and specialty care
- Emergency care
- Mental health and substance use disorder services
- Dental care
- Vision care
- Prescription drugs
- Hospital (inpatient and outpatient)
- Rehabilitation services

### Key Distinction from Private Insurance
NC Health Choice covers the full EPSDT benefit, which is broader than most commercial plans — it covers any medically necessary service for children identified in a screening.

---

## 4. Cost-Sharing

Unlike Medicaid (where cost-sharing is minimal), NC Health Choice has modest cost-sharing:

| Service | Cost-Sharing |
|---|---|
| Office visits | $5 copay |
| Specialist visits | $10 copay |
| Emergency room | $25 copay (waived if admitted) |
| Prescription drugs | $1–$3 copay |
| Annual premium | None (NC does not charge premiums for NC Health Choice as of 2026) |

> Federal law caps total annual cost-sharing at 5% of family income for CHIP families.

---

## 5. Enrollment Process

1. Apply online at **HealthCareConnect.nc.gov** or paper application through DSS county office
2. Medicaid/NC Health Choice uses a unified application — the system determines eligibility for Medicaid first, then NC Health Choice if Medicaid criteria not met
3. No waiting period (federal law prohibits CHIP waiting periods for children)
4. 12-month continuous eligibility — once enrolled, child stays enrolled for 12 months regardless of income changes (prevents coverage gaps)

---

## 6. Federal Financing (Enhanced FMAP)

CHIP receives a higher federal match than regular Medicaid:

| Program | Federal Match |
|---|---|
| Regular Medicaid (NC, FY2026) | ~68% |
| CHIP enhanced match (NC) | ~84% |

The CHIP enhanced FMAP is calculated by increasing the regular FMAP by 15 percentage points, but not exceeding 85%.

CHIP funding is capped at annual allotments (not open-ended like Medicaid), requiring Congressional reauthorization. CHIP has been reauthorized multiple times; the most recent long-term reauthorization was in 2018 (Pub. L. 115-120).

---

## 7. Relationship to Medicaid

NC operates **CHIP as Medicaid expansion** (as opposed to a separate CHIP program). This means:
- NC Health Choice enrollees are enrolled in Standard Plan managed care (same MCOs as Medicaid)
- Benefits package mirrors the Medicaid benefit
- Single application and eligibility determination process

Some states operate separate CHIP programs with different benefit packages; NC's integration simplifies the eligibility-to-enrollment pipeline.

---

## 8. Express Lane Eligibility

NC uses **Express Lane Eligibility (ELE)** — when a family's income is already verified through another public benefit program (SNAP, WIC, NSLP free/reduced lunch), the data can be used to fast-track CHIP/Medicaid enrollment without duplicative verification. This reduces administrative burden and increases enrollment among eligible children.

---

## 9. Anti-Patterns

**"NC Health Choice is only for families that earn too much for Medicaid"** — Correct in substance: it covers the 133–210% FPL band for children 6–18. Families near the lower bound may qualify for Medicaid instead; the unified application determines which program applies.

**"CHIP and NC Health Choice are different programs"** — NC Health Choice IS NC's CHIP program (Medicaid expansion CHIP type).

**"CHIP requires a waiting period"** — Federal law prohibits waiting periods for CHIP children. NC does not impose one.

---

## 10. References

- NCDHHS, NC Health Choice overview: [medicaid.ncdhhs.gov/nc-health-choice](https://medicaid.ncdhhs.gov/beneficiaries/nc-health-choice)
- HealthCareConnect.nc.gov (enrollment portal): [healthcareconnect.nc.gov](https://www.ncdhhs.gov/divisions/dhb/nc-health-choice)
- KFF, CHIP Program Overview: [kff.org/medicaid/fact-sheet/childrens-health-insurance-program-chip](https://www.kff.org/medicaid/fact-sheet/childrens-health-insurance-program-chip/)
- Medicaid.gov, CHIP overview: [medicaid.gov/chip](https://www.medicaid.gov/chip/index.html)
- Social Security Act § 2101 et seq. (Title XXI, CHIP); 42 U.S.C. § 1397aa et seq.
- CHIP Reauthorization Act (CHIPRA) 2009; Pub. L. 115-120 (2018 reauthorization)
