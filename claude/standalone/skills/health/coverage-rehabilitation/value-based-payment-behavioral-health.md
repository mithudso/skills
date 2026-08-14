---
name: value-based-payment-behavioral-health
description: >-
  Value-based payment (VBP) models for behavioral health services in Medicaid
  and Medicare: pay-for-performance, shared savings, bundled payments, and
  global capitation applied to mental health and SUD treatment; HEDIS BH
  quality measures; BH integration in primary care VBP (collaborative care
  model billing); Certified Community Behavioral Health Clinic (CCBHC) PPS as
  VBP; Section 1915(a) managed care quality incentives; CMS Innovation Center
  BH models (CMMI); SUD episode-based payment; measurement challenges in BH VBP;
  Medicare ACO behavioral health integration; and BH VBP in commercial markets.
  Educational only, NOT medical/insurance advice; verified-as-of 2026-06-22.
  TRIGGER: value-based payment behavioral health, pay for performance mental
  health, Medicaid VBP mental health, behavioral health bundled payment,
  CCBHC prospective payment, collaborative care model billing G codes, BHI
  billing Medicare, Medicare collaborative care model, HEDIS behavioral health
  measures value based, CMS innovation center behavioral health, CMMI behavioral
  health model, ACO behavioral health, Medicaid quality incentive behavioral
  health, SUD episode payment, value-based contracting mental health, global
  capitation behavioral health.
  SKIP: CCBHC overall model → medicaid-imd-ccbhc; HEDIS measures as network
  adequacy tool → network-adequacy-behavioral-health; Medicaid managed care
  capitation structure → medicaid-managed-care-behavioral-health; Medicare
  outpatient BH billing codes → medicare-mental-health-sud-coverage.
version: "1.1.0"
updated: "2026-06-22"
category: health
tags: [value-based-payment, behavioral-health, medicaid, medicare, quality, VBP]
hub: health
---

# Value-Based Payment for Behavioral Health

## Overview

Value-based payment (VBP) ties reimbursement to quality and outcomes rather than volume. Applying VBP to behavioral health poses distinct challenges: outcomes are harder to measure, attribution is complex, and episode boundaries are unclear (unlike surgery). Despite these challenges, payers have moved aggressively toward BH VBP in both Medicare and Medicaid.

## BH VBP in Medicaid

### Quality Incentive Programs (QIPs) / Pay-for-Performance
- Medicaid MCO contracts commonly include quality withhold pools tied to HEDIS BH measures; specific withhold percentages vary by state contract
- Key HEDIS BH measures used in VBP (NCQA HEDIS MY2025) [NCQA HEDIS Measures; Medicaid 2025 BH Core Set]:
  - **FUH** (Follow-Up After Hospitalization for Mental Illness): 7-day and 30-day follow-up rates
  - **IET** (Initiation and Engagement of Substance Use Disorder Treatment): new SUD diagnosis → treatment within 14 days, engagement within 34 days
  - **AMM** (Antidepressant Medication Management): 12-week and 6-month medication continuation
  - **FUI** (Follow-Up After Emergency Department Visit for Substance Use): 7-day and 30-day follow-up
- MCOs may pass incentives down to BH providers through value-based contracts

### Episode-Based Payment (EBP) for SUD
- Some states have piloted SUD episode bundles: a defined "episode" of treatment with bundled payment covering all services in a performance window [Tennessee Health Care Innovation Initiative: tn.gov/tenncare/health-care-innovation.html]
- Tennessee's Health Care Innovation Initiative explicitly includes behavioral health episodes as part of retrospective episode-of-care payments; the model rewards providers for high-quality, efficient acute and behavioral treatments
- Episode design challenges: SUD is chronic, relapsing; episode boundary definition determines incentives and may inadvertently penalize providers for patient relapse
- **OUD episode payment**: some states bundling MAT initiation + sustained engagement into performance period

### Medicaid Health Homes (Section 2703 / SSA §1945)
- Health homes receive PMPM (per-member per-month) care coordination payments for individuals with SMI or chronic conditions including BH [CMS Health Homes: medicaid.gov/medicaid/long-term-services-supports/health-homes]
- Enhanced FMAP: 90% federal match for the specific health home services in §1945 (comprehensive care management, care coordination, health promotion, transitional care, patient/family support, community referrals) — limited to first 8 quarters per enrollee; does NOT apply to underlying Medicaid services
- Quality measures required; performance can affect PMPM rates
- **BH integration**: health home providers must integrate and coordinate primary, acute, behavioral health, and long-term services — making BH-physical integration a core program requirement, not merely a metric

### CCBHC Prospective Payment System
- Authorized by Section 223 of the Protecting Access to Medicare Act (PAMA) of 2014 (P.L. 113-93); demonstration expanded by the Bipartisan Safer Communities Act of 2022 (P.L. 117-159) to allow up to 10 new states every two years beginning July 1, 2024 [Feldesman: feldesman.com/bipartisan-safer-communities-act-expands-medicaid-ccbhc-program/]
- Four PPS rate methodologies; QBPs **optional** under daily PPS-1 and PPS-3; **required** under monthly PPS-2 and PPS-4 [CMS CCBHC PPS/QBP guidance: medicaid.gov/medicaid/financial-management/certified-community-behavioral-health-clinic-ccbhc-demonstration/prospective-payment-system-pps-quality-bonus-payments-qbps; updated guidance Feb 15, 2024]
- QBP measures: CMS requires states using QBPs to include the first seven "comparative" quality measures for cross-state evaluation; states may add additional measures
- Quality measures commonly include: screening for mental health and SUD, FUH follow-up, depression remission, physical health screenings for individuals with SMI
- As of May 2026, HHS expanded the demonstration to 10 new states: AK, CO, HI, LA, MD, MS, MT, ND, WA, WV [HHS press release: hhs.gov/press-room/hhs-welcomes-10-new-states-into-ccbhc-medicaid-demonstration-program.html]

## BH VBP in Medicare

### Collaborative Care Model (CoCM) Billing
- CMS created billing codes for the Collaborative Care Model — evidence-based integration of BH into primary care [CMS MLN909432 BHI Services, Jan 2026: cms.gov/files/document/mln909432-behavioral-health-integration-services.pdf]
- **CoCM codes** (billed by treating PCP/NP/PA under their NPI):
  - **99492**: Initial month — ≥70 min of behavioral health care manager (BHCM) + clinical staff time; billed once per enrollment (or re-enrollment after 6+ month gap)
  - **99493**: Subsequent months — ≥60 min of BHCM + clinical staff time; primary ongoing code
  - **99494**: Add-on — each additional 30 min in same month; always accompanies 99492 or 99493
  - **G2214**: First 30 min of BHCM activity in a month; alternative base code for stable/brief months; mutually exclusive with 99492/99493 in same month
- **CoCM components**: treating provider (billing entity), behavioral health care manager (registry tracking, follow-up), psychiatric consultant (caseload review — no direct patient contact required), treat-to-target approach
- **Scope**: depression, anxiety, PTSD, SUD in primary care settings

### Behavioral Health Integration (BHI) in FFS Medicare
- **General BHI (99484 / G0323)**: ≥20 min/month of care management for BH conditions (including SUD) by clinical staff under physician direction; uses non-CoCM models; requires prior initiating visit (E/M, AWV, or IPPE within prior year) [CMS MLN909432 Jan 2026]
- **Psychiatric CoCM codes (99492–99494 + G2214)**: monthly time-based codes for the structured CoCM model with BHCM + psychiatric consultant
- These codes incentivize care integration; time-documentation requirements embed ongoing quality tracking in billing

### ACOs and BH Integration
- Medicare Shared Savings Program (MSSP) ACOs have adopted BH quality metrics; high-need BH patients are often highest ACO cost, creating strong financial incentive for BH investment
- **ACO REACH** (replaced Global and Professional Direct Contracting in 2023; runs PY 2023–2026): global risk option permits Total Care Capitation — a risk-adjusted monthly payment for all covered services including BH; ACO REACH ends December 2026 [CMS ACO REACH fact sheet: cms.gov/newsroom/fact-sheets/accountable-care-organization-aco-realizing-equity-access-and-community-health-reach-model]
- **LEAD Model** (Long-term Enhanced ACO Design): replaces ACO REACH January 1, 2027; 10-year model; two-sided risk (global: 100% savings/losses; professional: 50%); does NOT use global capitation — risk-sharing operates against benchmarks, not capitated payments [CMS LEAD: cms.gov/priorities/innovation/innovation-models/lead]
- **Accountable Health Communities (AHC) model (CMMI, 2017–2023)**: tested community navigator screening/referral for HRSNs; found individuals with BH conditions less likely to resolve social needs; AHC quality measures are being considered for MSSP APM Performance Pathway [CMS AHC: cms.gov/priorities/innovation/innovation-models/ahcm]

### CMMI BH Innovation Models
- **Innovation in Behavioral Health (IBH) Model**: announced Jan 18, 2024; 8-year model covering adult Medicaid and Medicare enrollees with moderate-to-severe mental health conditions and/or SUD [CMS IBH: cms.gov/priorities/innovation/innovation-models/ibh]
  - Structure: 3-year Pre-Implementation Period (2025–2027) + 5-year Implementation Period (2028–2032); practices may receive care delivery funding from 2028
  - Payment: Integration Support Payment (ISP) expected ~$200–$220 PMPM (Medicare, risk-adjusted, prospective quarterly); Medicaid component designed by participating states with flexibility (e.g., CCBHC, Health Home frameworks)
  - Performance-based payment (PBP): pay-for-reporting in model years 4–5, shifting to pay-for-performance (up to 5% of ISP) from year 6 onward; upside-only (no clawback)
  - Current Cohort I states: Michigan, New York, South Carolina; up to 8 states total eligible
  - "No wrong door" approach: patient access to physical, BH, and HRSN supports regardless of care-entry point

## Measurement Challenges in BH VBP

1. **Outcome attribution**: BH outcomes (stable housing, employment, abstinence) depend on factors outside provider control
2. **Episode definition**: unlike surgical episodes, BH episodes lack natural boundaries
3. **Risk adjustment**: high-complexity BH populations need adequate risk adjustment or providers avoid them
4. **Data fragmentation**: BH encounter data often separate from medical claims; attribution requires linked data
5. **Long-time horizons**: meaningful BH outcomes (sustained recovery, relapse prevention) require multi-year follow-up vs. typical 1-year payment cycles
6. **Equity**: VBP models risk penalizing safety-net BH providers serving disadvantaged populations if quality measures don't risk-adjust adequately

## Commercial Market BH VBP

Commercial payers (Aetna, UHC, BCBS) developing BH VBP:
- BCBS Massachusetts Quality Tiering: BH providers tiered by quality/cost outcomes; higher tier = higher reimbursement
- Aetna BH partnership programs: shared savings with BH provider organizations meeting engagement and outcomes benchmarks
- Challenge: commercial BH VBP is voluntary; resistant provider organizations can opt out

## References

1. **CMS IBH Model** (Innovation in Behavioral Health) — model overview, timeline, payment structure, participants:
   https://www.cms.gov/priorities/innovation/innovation-models/ibh
   *Verified-as-of: 2026-06-22*

2. **CMS IBH Model Announcement** — Jan 18, 2024 press release:
   https://www.cms.gov/newsroom/press-releases/cms-announces-new-model-advance-integration-behavioral-health
   *Verified-as-of: 2026-06-22*

3. **CMS MLN909432 — Behavioral Health Integration Services** (Jan 2026) — 99484, G0323, 99492–99494, G2214 code requirements:
   https://www.cms.gov/files/document/mln909432-behavioral-health-integration-services.pdf
   *Verified-as-of: 2026-06-22*

4. **CMS CCBHC Demonstration — PPS & Quality Bonus Payments** — four PPS methodologies, QBP required vs. optional, quality measures:
   https://www.medicaid.gov/medicaid/financial-management/certified-community-behavioral-health-clinic-ccbhc-demonstration/prospective-payment-system-pps-quality-bonus-payments-qbps
   *Verified-as-of: 2026-06-22*

5. **CMS CCBHC PPS Updated Guidance** (Feb 15, 2024) — PPS-3/PPS-4 new rate options, special crisis services rates, QBP flexibilities:
   https://www.medicaid.gov/medicaid/financial-management/downloads/updated-ccbh-pps-guidance-02152024.pdf
   *Verified-as-of: 2026-06-22*

6. **HHS CCBHC 10-State Expansion** (May 2026) — AK, CO, HI, LA, MD, MS, MT, ND, WA, WV added:
   https://www.hhs.gov/press-room/hhs-welcomes-10-new-states-into-ccbhc-medicaid-demonstration-program.html
   *Verified-as-of: 2026-06-22*

7. **Bipartisan Safer Communities Act (P.L. 117-159) — CCBHC expansion** — Section 11001 amends PAMA §223, extends demonstration, allows 10 new states every 2 years beginning July 1, 2024:
   https://www.feldesman.com/bipartisan-safer-communities-act-expands-medicaid-ccbhc-program/
   *Verified-as-of: 2026-06-22*

8. **CMS Health Homes (Section 2703 / SSA §1945)** — PMPM structure, 90% enhanced FMAP limited to 8 quarters per enrollee, eligible services:
   https://www.medicaid.gov/medicaid/long-term-services-supports/health-homes
   *Verified-as-of: 2026-06-22*

9. **CMS ACO REACH Model** — replaces Global/Professional Direct Contracting (2023–2026); global option with Total Care Capitation; BH integration:
   https://www.cms.gov/newsroom/fact-sheets/accountable-care-organization-aco-realizing-equity-access-and-community-health-reach-model
   *Verified-as-of: 2026-06-22*

10. **CMS LEAD Model** — replaces ACO REACH Jan 1, 2027; 10-year two-sided risk model; risk-sharing against benchmarks (not capitation):
    https://www.cms.gov/priorities/innovation/innovation-models/lead
    *Verified-as-of: 2026-06-22*

11. **CMS Accountable Health Communities (AHC) Model** — 2017–2023 HRSN navigation; BH findings; AHC measures considered for MSSP APM pathway:
    https://www.cms.gov/priorities/innovation/innovation-models/ahcm
    *Verified-as-of: 2026-06-22*

12. **NCQA HEDIS Measures** — FUH, IET (Substance Use Disorder Treatment, updated name), AMM, FUI definitions and MY2025 specifications:
    https://www.ncqa.org/hedis/measures/
    *Verified-as-of: 2026-06-22*

13. **Medicaid 2025 BH Core Set** — CMS-mandated BH quality measures for Medicaid/CHIP managed care, including HEDIS BH measures:
    https://www.medicaid.gov/medicaid/quality-of-care/downloads/2025-bh-core-set.pdf
    *Verified-as-of: 2026-06-22*

14. **Tennessee Health Care Innovation Initiative** — episode-of-care payments including behavioral health; retrospective bundled payment for acute and BH conditions:
    https://www.tn.gov/tenncare/health-care-innovation.html
    *Verified-as-of: 2026-06-22*

15. **Manatt Health — IBH Model Overview** (Feb 2024) — ISP PMPM $200–$220, performance-based payment structure, upside-only:
    https://www.manatt.com/insights/newsletters/health-highlights/cmmis-innovations-promoting-physical-and-behavior
    *Verified-as-of: 2026-06-22*
