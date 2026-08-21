---
name: medicaid-fmap-federal-financing
description: >-
  FMAP (Federal Medical Assistance Percentage) — the formula that determines the federal share of Medicaid spending; regular FMAP, enhanced FMAPs (expansion 90/10, CHIP, FMAP floors/ceilings), maintenance of effort, allotments vs open-ended funding, COVID-era enhanced FMAP, NC FMAP rate. TRIGGER: FMAP Medicaid, Federal Medical Assistance Percentage, federal match Medicaid, Medicaid federal share, state share Medicaid, enhanced FMAP, FMAP formula, Medicaid block grant vs FMAP, FMAP floor ceiling, NC FMAP rate, Medicaid financing federal state, Medicaid expenditure CMS. SKIP: ACA Medicaid expansion → medicaid-aca-expansion; CHIP financing → medicaid-nc-health-choice-chip; NC Medicaid Managed Care → nc-health-services-medicaid-managed-care.
version: "1.0.0"
updated: "2026-06-22"
category: health
tags: [medicaid, FMAP, federal-financing, federal-match, CMS, state-share, budget]
related_skills:
  - medicaid-aca-expansion
  - medicaid-nc-health-choice-chip
  - enrollment-eligibility-medicaid-eligibility-determination
  - nc-health-services-medicaid-managed-care
---

# FMAP: Federal Medical Assistance Percentage

> **Educational only — NOT financial or legal advice.** FMAP rates are published annually by CMS in the Federal Register; state budget and policy impacts require current official figures.

**verified-as-of: 2026-06-22**

---

## 1. What Is FMAP?

The **Federal Medical Assistance Percentage (FMAP)** is the share of each state's Medicaid expenditures that the federal government reimburses. It is calculated annually using a formula based on each state's **per capita income** relative to the national average.

**Formula** (Social Security Act § 1905(b)):
```
FMAP = 1 - (State per capita income² / National per capita income²) × 0.45
```

- **Floor**: No state gets less than 50% federal match
- **Ceiling**: No state gets more than 83% (per statutory rule)
- **Average**: ~57% nationwide (federal pays ~57 cents of every Medicaid dollar)

Wealthier states (higher per capita income) → lower FMAP (federal pays less)
Poorer states (lower per capita income) → higher FMAP (federal pays more)

This makes Medicaid **counter-cyclical and redistributive** — states with lower incomes receive more federal support.

---

## 2. NC FMAP Rate

NC's FMAP varies year to year based on income data (3-year average lag). As of FY2026:
- **Regular FMAP (NC)**: approximately **67–68%** (federal pays ~$0.67 per $1.00 of regular Medicaid spending)
- This means NC pays approximately **32–33%** of regular Medicaid costs from state funds

NC is a mid-range state — not the highest FMAP (which goes to states like Mississippi ~77%) or the lowest (California, New York, Massachusetts, Connecticut, New Hampshire, Delaware, Wyoming get the 50% floor).

---

## 3. Enhanced FMAPs

Several federal programs use higher match rates than the regular FMAP:

| Program | Federal Match | Statutory Authority |
|---|---|---|
| **ACA Medicaid expansion** (adult group 138% FPL) | 90% | ACA § 2001 |
| **CHIP** | FMAP + 15 pp (capped at 85%) | SSA § 2105 |
| **Community First Choice** (1915(k) HCBS) | FMAP + 6 pp | ACA § 2401 |
| **FMAP for certain home/community-based services** (ARP) | FMAP + 10 pp | ARPA 2021 (temporary) |
| **Qualifying individuals** (QI Medicare Savings) | 100% | OBRA 1990 |
| **Family planning services** | 90% | SSA § 1905(a) |

The 90/10 split for ACA expansion enrollees is the most significant: the federal government pays 90% of costs for adults added through ACA expansion, creating strong state incentive to expand.

---

## 4. Open-Ended Entitlement Structure

A critical feature of Medicaid financing: it is an **open-ended entitlement**. There is no federal cap on how much the federal government pays states. If states spend more (more enrollees, more expensive services), the federal government automatically pays its FMAP share.

This contrasts with:
- **Block grants**: Fixed federal payment regardless of state spending — states bear all cost risk above the block amount
- **Allotments** (CHIP): Annual caps on federal funding; states that exhaust their allotment must fund additional CHIP costs entirely with state dollars

### Why Open-Ended Matters
- Medicaid is automatically countercyclical: recessions increase enrollment (job loss → income drop → eligibility), and federal spending rises proportionally
- Block grant proposals periodically appear in Congress as a way to reduce federal Medicaid spending by capping federal liability

---

## 5. Maintenance of Effort (MOE)

As a condition of receiving enhanced FMAP (especially during COVID and ARP periods), states must meet **maintenance of effort (MOE)** requirements:
- Cannot reduce Medicaid eligibility below pre-COVID/pre-expansion levels during the MOE period
- Cannot disenroll eligible enrollees continuously enrolled at the time of a public health emergency

The **COVID-era continuous enrollment requirement** (Families First Coronavirus Response Act, March 2020) prohibited states from disenrolling Medicaid enrollees in exchange for an enhanced FMAP (+6.2 percentage points). This ended March 31, 2023.

### Medicaid Unwinding (2023–2024)
After the continuous enrollment ended, states conducted **Medicaid redeterminations** (renewals) for all enrollees. NC and other states found:
- Millions of enrollees were removed ("unwound") — some legitimately income-ineligible, some due to procedural/administrative errors
- CMS monitored state-level unwinding data and pressured states to improve accuracy
- NC completed its unwinding period by 2024

---

## 6. Federal Financing: Where the Money Goes

Federal Medicaid spending (~$900 billion annually nationally) flows to states as **reimbursement claims** (not upfront grants):
1. State pays providers (or managed care plans) for Medicaid services
2. State files quarterly claims with CMS
3. CMS reimburses the state's FMAP share
4. Administrative costs: separate FMAP (usually 50–75%) for state administration

This means states must have liquidity to pay claims before receiving federal reimbursement — a significant cash flow consideration for states.

---

## 7. Policy Debates

**Block Grant Proposals**: Periodically debated in Congress (most recently as AHCA 2017). Would replace open-ended FMAP with fixed per-capita or aggregate federal payment. Would shift financial risk to states, likely requiring cuts to eligibility, benefits, or provider rates during downturns.

**Enhanced FMAP Sunsets**: Temporary enhanced rates (COVID, ARPA) expire, requiring states to absorb higher costs. States sometimes expand programs while enhanced rates are in effect, then face fiscal stress when they expire.

**FMAP Reform**: Some argue the current formula doesn't adequately account for state fiscal capacity or health care cost differences.

---

## 8. Anti-Patterns

**"Medicaid is a block grant"** — False. Medicaid is an open-ended federal entitlement. Block grant is a different financing model; periodically proposed but never enacted.

**"The federal government pays 50% of Medicaid"** — 50% is the floor for the wealthiest states. Most states receive more. NC gets ~67-68% regular FMAP; expansion enrollees cost NC only ~10%.

**"States can cut Medicaid whenever they want"** — Federal maintenance-of-effort requirements (tied to enhanced FMAP) restrict states' ability to reduce eligibility during certain periods.

---

## 9. References

- Social Security Act § 1905(b) — FMAP formula
- CMS Federal Register, FMAP rates (published annually): [federalregister.gov/search FMAP](https://www.federalregister.gov/search?query=federal+medical+assistance+percentages)
- KFF, Medicaid Financing and Spending: [kff.org/medicaid/state-indicator/medicaid-spending](https://www.kff.org/medicaid/state-indicator/medicaid-spending/)
- CMS, FMAP information by state: [medicaid.gov/medicaid/finance/medicaid-federal-financial-participation](https://www.medicaid.gov/medicaid/finance/medicaid-federal-financial-participation/index.html)
- MACPAC (Medicaid and CHIP Payment and Access Commission): [macpac.gov](https://www.macpac.gov)
- Families First Coronavirus Response Act (FFCRA), Pub. L. 116-127 (2020)
- American Rescue Plan Act (ARPA), Pub. L. 117-2 (2021)
