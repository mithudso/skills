---
name: enrollment-eligibility-magi-vs-non-magi-medicaid
description: MAGI (Modified Adjusted Gross Income) vs non-MAGI income methodology for Medicaid — who uses which, how income is counted differently under each, the 5% FPL income disregard, deductions allowed/disallowed, asset tests, and why the distinction matters for eligibility. TRIGGER: MAGI Medicaid; non-MAGI Medicaid; MAGI vs non-MAGI; Medicaid income methodology; MAGI income counting; Medicaid 138% FPL; Modified Adjusted Gross Income Medicaid; Medicaid income disregard; Medicaid asset test; SSI-related Medicaid. SKIP: full Medicaid eligibility determination process → enrollment-eligibility-medicaid-eligibility-determination; NC Medicaid application → enrollment-eligibility-nc-medicaid-application.
version: "1.0.0"
updated: "2026-06-22"
category: health
tags: [medicaid, MAGI, non-MAGI, income-methodology, FPL, asset-test, ACA, CMS]
related_skills:
  - enrollment-eligibility-medicaid-eligibility-determination
  - enrollment-eligibility-nc-medicaid-application
  - enrollment-eligibility-aca-marketplace-enrollment
  - consumer-finance
---
# MAGI vs Non-MAGI Medicaid Income Methodology

## Why Two Systems?

The ACA standardized Medicaid income counting using tax-based rules (MAGI) for most groups. However, elderly and disabled individuals — whose Medicaid often covers long-term services, durable medical equipment, and complex benefits — retained pre-ACA income and asset methodologies (non-MAGI / SSI-related rules).

## MAGI Groups (Use MAGI Methodology)

- Children up to 19
- Pregnant women
- Parents and caretaker relatives
- Adults 19–64 (expansion group in expansion states)
- Medically needy (some states apply MAGI)

## Non-MAGI Groups (Use SSI-Related or State Plan Methodology)

- Elderly (65+) on Medicaid for long-term care / nursing home
- People with disabilities (unless in MAGI group based on income)
- Blind individuals
- Medically needy spend-down (varies by state)
- Individuals who receive SSI (often automatically eligible for Medicaid — 1634 states)

## MAGI Income Calculation

**Based on:** IRS Modified Adjusted Gross Income concept, applied household by household.

**What counts as income (MAGI):**
- Wages, salaries, tips
- Self-employment net income
- Unemployment compensation
- Social Security benefits (must include — even though they're excluded from taxable income for many people, they count for MAGI Medicaid)
- Retirement distributions, pensions
- Rental income (net)
- Capital gains
- Alimony received (for agreements pre-2019)
- Interest, dividends

**What is EXCLUDED from MAGI:**
- Child support received
- Supplemental Security Income (SSI) payments
- Gifts and inheritances
- Income of tax filer who is not in the household group
- Veterans' benefits (non-service-connected)
- Workers' compensation

**Household composition rules:**
- Based on IRS tax filing unit rules, not physical household
- Filer + dependents claimed on tax return = household
- Specific rules for non-filers, dependents claimed by someone outside household, pregnant women (counted as 2)

**5% FPL Income Disregard:**
- Every MAGI Medicaid eligibility standard has a 5% FPL disregard applied before comparing to limit
- Effectively means: stated income limit of 133% FPL + 5% disregard = 138% FPL effective threshold
- The 5% disregard replaced individualized income disregards from pre-ACA Medicaid

**No asset test for MAGI groups** — post-ACA, MAGI Medicaid groups have no resource/asset limit.

## Non-MAGI Income Calculation

**Based on:** Pre-ACA SSI methodology, which allows many more deductions.

**Income types counted:**
- Earned income (wages) with deductions
- Unearned income (Social Security, pensions, dividends)

**Deductions and disregards allowed (examples):**
- First $20/month of most income (general income disregard)
- First $65/month of earned income + half of remainder (earned income disregard)
- Work expenses related to disability
- Plan to Achieve Self-Support (PASS) deductions
- Impairment-related work expenses

**Result:** Non-MAGI gross income threshold may be lower on paper, but after disregards, net countable income is often much lower, expanding eligibility for disabled/elderly.

**Asset/Resource Test for Non-MAGI:**
- Resource limit: ~$2,000 individual / $3,000 couple (varies by state and program)
- Excluded assets: primary home (if residing there or intend to return), one vehicle (often), burial funds up to limit, certain life insurance
- Countable assets: savings accounts, investment accounts, second vehicles, cash value of most life insurance
- Some states have eliminated or raised resource limits through 1902(r)(2) waivers

## Practical Illustration

Person: Age 45, receives $800/month SSDI + earns $400/month part-time

**MAGI calculation:**
- Total income: $800 + $400 = $1,200/month
- Annual: $14,400
- 2025 FPL (single): $15,060/year → 138% = $20,783
- $14,400 < $20,783 → likely eligible under MAGI expansion

**Non-MAGI (if they qualified for that group):**
- SSDI $800: subtract $20 general disregard = $780 countable unearned
- Earned $400: subtract $65 disregard = $335; subtract half = $167.50 countable earned
- Total countable income: $780 + $167.50 = $947.50/month
- Compare against SSI-related Medicaid income standard (state-specific)

## State Variation

States may:
- Use more generous income methodologies via 1902(r)(2) waivers (non-MAGI)
- Expand MAGI income limits beyond federal minimums
- Apply medically needy / spend-down for non-MAGI groups when income exceeds standard

## Sources

- 42 CFR §435.603 — MAGI-based income methodology
- 26 U.S.C. §36B(d)(2) — MAGI definition (incorporated by reference)
- CMS MAGI Conversion: CMS guidance on "Converting Non-MAGI Rules to MAGI Equivalent" (2012)
- SSI income rules: 20 CFR Part 416, Subpart K
- CMS Medicaid eligibility: medicaid.gov/medicaid/eligibility/index.html
- Kaiser Family Foundation: "Understanding the ACA's Medicaid Expansion" (updated 2024)
