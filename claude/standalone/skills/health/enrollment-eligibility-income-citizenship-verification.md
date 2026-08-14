---
name: enrollment-eligibility-income-citizenship-verification
description: Income and citizenship/immigration verification for Medicare, Medicaid, and ACA marketplace enrollment — Federal Data Hub, DHS SAVE system, IRS data matching, self-attestation rules, data matching issues (DMIs), MAGI income documentation, asset verification for non-MAGI groups, and immigration status categories that confer eligibility. TRIGGER: Medicaid income verification; ACA income verification; citizenship verification health insurance; DHS SAVE Medicaid; Federal Data Hub; immigration status Medicaid; income documentation Medicaid; Medicaid data matching issue; ACA data matching issue; immigration health insurance eligibility; qualified immigrant Medicaid; undocumented Medicaid. SKIP: MAGI income methodology → enrollment-eligibility-magi-vs-non-magi-medicaid; ACA marketplace enrollment → enrollment-eligibility-aca-marketplace-enrollment; Medicaid determination process → enrollment-eligibility-medicaid-eligibility-determination.
version: "1.0.0"
updated: "2026-06-22"
category: health
tags: [income-verification, citizenship, immigration, DHS-SAVE, Federal-Data-Hub, MAGI, data-matching-issue, Medicaid, ACA, PRWORA]
related_skills:
  - enrollment-eligibility-medicaid-eligibility-determination
  - enrollment-eligibility-magi-vs-non-magi-medicaid
  - enrollment-eligibility-aca-marketplace-enrollment
  - enrollment-eligibility-nc-medicaid-application
---
# Income and Citizenship/Immigration Verification

## Federal Data Hub

Both the ACA marketplace and state Medicaid agencies use the **Federal Data Hub** — a CMS-operated data matching system connecting to:

- **IRS:** Prior-year tax return data (MAGI income, household composition)
- **SSA:** Social Security number verification, benefit amounts, citizenship confirmation
- **DHS SAVE:** Immigration status verification (see below)
- **VA:** VA benefit verification (for Medicaid)

**How it works:**
1. Applicant submits application with self-attested information
2. System queries Federal Data Hub in real time
3. If data matches: eligibility determination proceeds
4. If inconsistency: Data Matching Issue (DMI) triggered

## Income Verification

### Self-Attestation First

Federal policy: **accept self-attestation initially**. States and the marketplace must:
- Accept self-attested income as reasonable and compatible
- Not require documentation upfront if self-attestation is consistent with electronic data
- Give 90-day "inconsistency period" to resolve DMIs before terminating coverage

### MAGI Income Documentation (when required)

Primary income sources requiring documentation when queried:

| Income Type | Acceptable Documentation |
|---|---|
| Wages/salary | Recent pay stubs (last 30 days), employer letter |
| Self-employment | Most recent tax return, profit/loss statement, business bank records |
| Social Security | SSA award letter or statement |
| Unemployment | Benefit determination letter |
| Rental income | Tax return Schedule E, lease agreements |
| Retirement/pension | Award letter, 1099-R |
| Alimony | Divorce decree, bank statements |
| No income | Self-attestation; may require written statement |

### What MAGI Income Excludes

Critical for applicants: these do NOT count for MAGI:
- Child support received
- SSI payments
- Workers' compensation
- Gifts / inheritances
- Most veterans' benefits (non-service-connected)

### Non-MAGI Asset Documentation

For elderly/disabled (non-MAGI) Medicaid applicants:
- Bank statements (checking, savings) — typically last 3 months
- Investment account statements
- Certificate of deposit documentation
- Life insurance policy (cash value)
- Property records (if contesting that property is primary home)
- Transfer of assets records (for long-term care Medicaid — look-back period)

**Long-term care Medicaid 5-year look-back:** Applicants for nursing home Medicaid must document asset transfers in prior 60 months; inappropriate transfers may result in disqualification period.

## Citizenship Verification

### US Citizens

- Verified through SSA (match of SSN, name, date of birth)
- If SSA cannot verify: must provide documentary evidence
  - US birth certificate
  - US passport or passport card
  - Certificate of Naturalization (N-550/N-570)
  - Certificate of Citizenship (N-560/N-561)
  - US Consular Report of Birth Abroad

### Immigration Status — DHS SAVE System

**DHS SAVE (Systematic Alien Verification for Entitlements):**
- Web-based system used by federal/state agencies to verify immigration status
- Used by Medicaid agencies, state CHIP programs, and ACA marketplace
- Queries DHS databases in real time or near-real time

**Immigration categories that qualify for Medicaid/CHIP** (selected):

| Category | Medicaid/CHIP Eligibility |
|---|---|
| Lawful Permanent Resident (LPR) | After 5-year bar (most programs); some states waive with state funds |
| Refugee | No 5-year bar; immediate eligibility |
| Asylee | No 5-year bar; immediate eligibility |
| Cuban/Haitian Entrant | No 5-year bar |
| Amerasian | No 5-year bar |
| Trafficking victim (T visa) | No 5-year bar |
| Withholding of Deportation/Removal | After 5-year bar |
| Special Immigrant Juvenile | After 5-year bar (some states earlier with state funds) |
| VAWA self-petitioner | After 5-year bar |
| COFA (Compact of Free Association) | Varies by state |

**Not eligible for federal Medicaid:**
- Undocumented immigrants (Emergency Medicaid only)
- F, J, H, B visa holders (non-immigrant, temporary)
- DACA recipients (may be eligible for state-funded programs in some states)

**ACA Marketplace:** Lawfully present immigrants can enroll regardless of 5-year bar (different rule from Medicaid). DACA recipients newly eligible as of 2024.

## Data Matching Issues (DMIs)

### What triggers a DMI?

- Income self-attested differs significantly from IRS data (prior year vs. projected current year common mismatch)
- SSN not found in SSA records
- Citizenship/immigration status cannot be verified
- Household composition inconsistency

### Resolution Process

1. Applicant receives notice of DMI with explanation
2. 90-day inconsistency period to submit documentation
3. During 90 days: coverage/benefits continue for Medicaid; APTC continues for marketplace
4. Submit documentation via:
   - Marketplace: upload via HealthCare.gov or mail
   - Medicaid: submit to county DSS or state Medicaid agency
5. If resolved: enrollment confirmed with original effective date
6. If unresolved at 90 days: 
   - Marketplace: loss of APTC (must pay full premium)
   - Medicaid: loss of coverage

### Common DMI Pitfalls

- **Income change:** Self-employed with variable income; income this year differs greatly from prior tax return — provide current-year documentation
- **Name discrepancy:** Name on application doesn't match SSA record exactly (married name, hyphenation)
- **Newly arrived immigrants:** SAVE cannot yet find record — provide immigration documents directly
- **Recent citizenship:** Naturalized citizen not yet in SSA database — provide naturalization certificate

## Reasonable Compatibility Standard

Agencies must apply the "reasonable compatibility" standard: if self-attested information is reasonably compatible with electronic data, accept it without requiring documentation. "Reasonably compatible" generally means income within 10% of data source, or there is a plausible explanation for the difference.

States cannot require documentation as a routine step — documentation only required when self-attestation is not reasonably compatible with data hub results.

## Sources

- 42 CFR §435.952 — Medicaid income and resource verification
- 42 CFR §457.380 — CHIP eligibility verification
- 45 CFR §155.315 — Marketplace income/citizenship verification
- CMS: Federal Data Hub overview (cms.gov/cciio/resources/data-resources)
- DHS SAVE: uscis.gov/save
- PRWORA §401–435 — Immigration eligibility bars
- CMS SHO letter: guidance on income verification reasonable compatibility standard (2013)
- INA §212(a)(4) — Public charge consideration (separate from Medicaid eligibility but often confused with it)
