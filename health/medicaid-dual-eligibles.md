---
name: medicaid-dual-eligibles
description: >-
  Dual eligibles — individuals enrolled in both Medicare and Medicaid — eligibility categories (full duals, partial duals, QMB/SLMB/QI), benefit coordination (Medicare primary, Medicaid secondary), Dual Eligible Special Needs Plans (D-SNPs), Integrated Care (PACE, financial alignment demos), cost-sharing protections, NC dual enrollment. TRIGGER: dual eligible Medicare Medicaid, dual enrolled Medicare Medicaid, QMB Medicare Savings Program, SLMB Medicare Savings Program, D-SNP dual special needs plan, PACE program, Medicaid pays Medicare premiums, full dual partial dual, Medicaid wrap-around benefits, Medicare low income subsidy LIS, integrated care duals, NC dual eligible. SKIP: Medicare parts/eligibility standalone → enrollment-eligibility-medicare-eligibility-criteria; Medicaid eligibility standalone → enrollment-eligibility-medicaid-eligibility-determination; NC managed care → nc-health-services-medicaid-managed-care.
version: "1.0.0"
updated: "2026-06-22"
category: health
tags: [medicaid, Medicare, dual-eligible, QMB, D-SNP, PACE, integrated-care, cost-sharing]
related_skills:
  - enrollment-eligibility-medicaid-eligibility-determination
  - enrollment-eligibility-medicare-eligibility-criteria
  - nc-health-services-medicaid-managed-care
  - medicaid-fmap-federal-financing
  - nc-health-services-medicaid-waivers
---

# Dual Eligibles: Medicare-Medicaid Coordination

> **Educational only — NOT health or legal advice.** Benefit rules, cost-sharing protections, and enrollment procedures change frequently; verify with your State Health Insurance Assistance Program (SHIP) or state Medicaid agency.

**verified-as-of: 2026-06-22**

---

## 1. Who Are Dual Eligibles?

**Dual eligibles** are individuals enrolled in BOTH Medicare AND Medicaid simultaneously — approximately **12–13 million people** nationally (about 20% of Medicaid enrollees, 18% of Medicare enrollees). They tend to be older, sicker, and lower-income than typical Medicare or Medicaid enrollees alone.

### Why They Have Both
- Medicare covers elderly (65+) and people with disabilities (after 24-month SSDI waiting period)
- Medicaid covers low-income individuals; many elderly/disabled people have incomes and assets below Medicaid thresholds
- These populations overlap significantly

---

## 2. Categories of Dual Eligibility

### Full Duals
Enrolled in both Medicare and full Medicaid. Medicaid covers **cost-sharing** (premiums, deductibles, copays), **wrap-around benefits** not covered by Medicare (dental, vision, hearing, long-term care), and services until Medicare exhausts coverage.

### Partial Duals — Medicare Savings Programs (MSPs)
Medicaid pays some but not all Medicare costs. Four federal MSP categories:

| MSP | Income Limit | What Medicaid Pays |
|---|---|---|
| **QMB** (Qualified Medicare Beneficiary) | ≤100% FPL | Part A + Part B premiums, deductibles, and copays |
| **SLMB** (Specified Low-Income Medicare Beneficiary) | 100–120% FPL | Part B premium only |
| **QI** (Qualifying Individual) | 120–135% FPL | Part B premium only (limited slots, first-come) |
| **QDWI** (Qualified Disabled and Working Individual) | ≤200% FPL, working | Part A premium only |

> **QMB Key Rule**: QMB enrollees **cannot be billed** for Medicare cost-sharing (deductibles, copays) — providers who accept Medicare must accept Medicaid as payment in full. This is federal law (42 U.S.C. § 1396a(n)); QMB billing protections are frequently violated in practice.

---

## 3. Benefit Coordination

For full duals, coordination follows a strict order:

1. **Medicare pays first** (primary payer) for services Medicare covers
2. **Medicaid pays second** (secondary payer) for:
   - Medicare cost-sharing
   - Services Medicare doesn't cover (LTSS, dental, vision, certain drugs)
3. Member pays nothing or minimal amount

### Why Coordination Is Complex
- Different coverage rules, different claims systems
- Part D (drug coverage) interacts with Medicaid's "wrap" for drug costs
- Low-Income Subsidy (LIS/"Extra Help") auto-qualifies full duals for Part D cost-sharing assistance
- Long-term care: Medicaid is primary payer for nursing home costs once Medicare skilled nursing coverage (SNF) exhausts

---

## 4. Dual Eligible Special Needs Plans (D-SNPs)

**D-SNPs** are Medicare Advantage plans specifically designed for dual eligibles. They:
- Must coordinate with Medicaid
- May offer enhanced benefits (more dental, vision, transportation, meal delivery)
- Receive a higher capitation rate from CMS due to complexity of the population
- Are required to offer unified appeals processes (integrated with Medicaid)

**Integrated D-SNPs (HIDE-SNPs and FIDE-SNPs)** have deeper integration with Medicaid managed care — a single plan manages both Medicare and Medicaid benefits for the member.

In NC (as of 2026): Several D-SNP options are available, including through Standard Plan MCOs that also participate in Medicaid managed care.

---

## 5. PACE — Program of All-Inclusive Care for the Elderly

**PACE** (Program of All-Inclusive Care for the Elderly) is a fully integrated dual-eligible program:
- Eligible: age 55+, nursing home-eligible, able to live safely in community
- One organization provides ALL Medicare and Medicaid services
- Services provided through adult day health center + home care + specialists
- Capitated payment: PACE org receives both Medicare and Medicaid monthly payments, bears full risk
- Strong evidence base: reduces nursing home placement, hospitalizations

NC has PACE sites in several major metro areas including Raleigh-Durham, Charlotte, Greensboro.

---

## 6. Financial Alignment Demonstrations

CMS and states have piloted **financial alignment demonstrations** to integrate Medicare and Medicaid financing for duals. Two models:
1. **Capitated model**: Single managed care plan receives blended Medicare-Medicaid payment for duals
2. **Managed FFS model**: Fee-for-service with care coordination

NC has not operated a financial alignment demonstration; most NC duals are served through Standard Plans (Medicaid) + their choice of Medicare coverage (original Medicare or MA plan).

---

## 7. Long-Term Care and Medicaid for Duals

Medicaid is the dominant payer for **long-term services and supports (LTSS)** — nursing facility care, home and community-based services. For dual eligibles in nursing homes:
- Medicare covers first 20 days at 100% (SNF benefit)
- Days 21–100: Medicare copay ~$200/day; Medicaid typically covers this for full duals
- Day 101+: Medicare SNF benefit exhausts; Medicaid becomes primary payer for nursing facility costs

This creates the major Medicaid cost: ~30–35% of all Medicaid spending goes to LTSS, concentrated in the dual-eligible population (who represent ~40% of total Medicaid spending despite being 20% of enrollees).

---

## 8. NC Dual Enrollment (2026 Estimates)

- Approximately **230,000–250,000** NC Medicaid enrollees are dual eligibles
- Majority are elderly (65+); smaller proportion are younger adults with disabilities (SSI/SSDI)
- Full duals in NC access LTSS primarily through Medicaid Direct (fee-for-service) or Tailored Plans if they have serious MH/IDD/SUD/TBI needs
- Medicare Savings Programs (QMB/SLMB/QI) enrollment is administered by NCDHHS Division of Health Benefits

---

## 9. Anti-Patterns

**"Dual eligibles are fully covered with no out-of-pocket costs"** — Full duals have very low or zero cost-sharing, but partial duals (QMB, SLMB) have narrower protections. QMB enrollees can still be illegally billed despite protection — they should report billing to their State Medicaid office.

**"D-SNPs are always better than original Medicare + Medicaid"** — D-SNPs offer advantages but may have network restrictions; not universally superior. PACE is often a better option for nursing home-eligible individuals who want to stay in the community.

**"Medicare always pays first"** — True for services Medicare covers. For services only Medicaid covers (e.g., LTSS beyond SNF benefit), Medicaid is primary.

---

## 10. References

- Medicaid.gov, Dual Eligible Individuals overview: [medicaid.gov/medicaid/eligibility/dual-eligibles](https://www.medicaid.gov/medicaid/eligibility/dual-eligibles/index.html)
- CMS, Medicare-Medicaid Coordination Office (MMCO): [cms.gov/medicare-medicaid-coordination](https://www.cms.gov/Medicare-Medicaid-Coordination)
- KFF, Dual Eligible Beneficiaries (2025 chartbook): [kff.org/medicaid/report/medicare-medicaid-dual-eligibles](https://www.kff.org/medicaid/report/medicare-medicaid-dual-eligible-beneficiaries/)
- Medicare Savings Programs: [medicare.gov/basics/costs/help/medicare-savings-programs](https://www.medicare.gov/basics/costs/help/medicare-savings-programs)
- CMS PACE program overview: [cms.gov/Medicare/Health-Plans/pace](https://www.cms.gov/Medicare/Health-Plans/pace)
- 42 U.S.C. § 1396a(n) — QMB billing protection
- NCDHHS, Medicare Savings Programs: [medicaid.ncdhhs.gov/beneficiaries/medicare-savings-programs](https://medicaid.ncdhhs.gov/beneficiaries/medicare-savings-programs)
