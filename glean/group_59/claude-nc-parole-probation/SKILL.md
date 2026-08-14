---
name: nc-parole-probation
description: >-
  NC parole, probation & community supervision hub — PRS, conditions, violations, electronic monitoring, CRV sanctions, DAC, reentry. Educational only, NOT legal advice. TRIGGER: NC probation/parole/post-release supervision PRSPC, NC CRV confinement-in-response-to-violation, NC absconder warrants ACE SOIU, NC ICAOS interstate compact, NC JMARC drug court, NC Justice Reinvestment Act JRA 2011, NC MAPP parole, NC medical geriatric parole GS 15A-1369, NC sex offender SORNA SBM, NC electronic monitoring GPS, NC day reporting DRC, NC-PASE risk assessment, NC evidence-based practices recidivism, NC Reentry 2030 SAVAN victim services, NC split sentence GS 15A-1351, NC probation length GS 15A-1342. SKIP: NC criminal trial/defense/pre-conviction → nc-criminal-defense (sentencing grids refs/nc-structured-sentencing.md; expungement refs/nc-expungement-record-sealing.md); SUD treatment/MAT/OTP → nc-drug-rehabilitation-treatment; specific case advice → consult a licensed NC attorney.
version: "1.2.2"
updated: "2026-06-22"
category: legal
model: claude-opus-4-8
effort: high
tags: [nc-parole-probation, legal, north-carolina, community-supervision, parole, probation, post-release-supervision, educational]
keywords:
  - NC probation
  - NC parole
  - NC post-release supervision
  - NC community supervision
  - NC probation violation
  - NC revocation hearing
  - NC Department of Adult Correction
  - NC supervised release
  - NC electronic monitoring GPS
  - NC ankle bracelet probation
  - NC house arrest
  - NC split sentence
  - NC special probation GS 15A-1351
  - NC swift certain proportionate sanctions
  - NC CRV confinement response violation
  - NC OPUS offender tracking
  - NC probation officer
  - NC Post-Release Supervision and Parole Commission PRSPC
  - NC drug court supervision JMARC
  - NC sex offender supervision SORNA
  - NC satellite based monitoring SBM
  - NC absconder warrants
  - NC ACE absconder capture enforcement
  - NC SOIU special operations intelligence
  - NC ICAOS interstate compact supervision
  - NC interstate compact probation transfer
  - NC Justice Reinvestment Act JRA 2011
  - NC MAPP mutual agreement parole program
  - NC medical geriatric parole GS 15A-1369
  - NC compassionate release prison
  - NC parole commission hearing process
  - NC probation supervision fee $40
  - NC probation financial conditions fee waiver
  - NC day reporting centers DRC ADRC
  - NC risk assessment tools NC-PASE
  - NC Static-99 sex offender risk assessment
  - NC evidence-based practices recidivism RNR
  - NC structured sentencing probation length
  - NC felony probation maximum 5 years
  - NC victim services SAVAN notification
  - NC Reentry 2030 strategic plan
  - NC local reentry councils LRC
  - NC intensive probation supervision IPS
  - NC reentry transitional housing
  - NC day reporting center substance abuse
  - NC probation financial conditions court costs
whenToUse:
  - "Questions about NC probation, parole, or post-release supervision"
  - "NC probation conditions (regular, intensive, supervised, special)"
  - "NC probation violation hearings, responses, and revocation"
  - "NC electronic monitoring, house arrest, or GPS ankle bracelet supervision"
  - "NC swift-certain-proportionate sanctions framework"
  - "NC Department of Adult Correction structure and community supervision"
  - "NC drug court or specialty court supervision (JMARC phases, conditions)"
  - "NC reentry after incarceration — programs, housing, services"
  - "NC sex offender supervision, SORNA tiers, SBM/satellite monitoring, registration"
  - "NC split sentence, suspended sentence, or active vs supervised release (GS 15A-1351)"
  - "NC CRV (Confinement in Response to Violation) — 90-day short confinement, two-CRV limit"
  - "NC absconder warrants — absconding definition, ACE teams, SOIU, consequences"
  - "NC ICAOS interstate compact — transferring probation/parole supervision to another state"
  - "NC Justice Reinvestment Act (JRA) 2011 — history, CRV creation, revocation limits, outcomes"
  - "NC MAPP mutual agreement parole program — structured milestones, pre-1994 eligibility"
  - "NC medical or geriatric parole / compassionate release (GS 15A-1369)"
  - "NC Parole Commission (PRSPC) — discretionary parole process, hearings, conditions"
  - "NC probation supervision fees ($40/month), fee waivers, EM device fees, court costs"
  - "NC day reporting centers (DRC/ADRC) — intermediate sanction, CBI, substance abuse services"
  - "NC risk assessment tools — NC-PASE predictive analytics, RNR model, risk classification tiers"
  - "NC evidence-based practices and recidivism reduction (Motivational Interviewing, Carey Guides, CBI)"
  - "NC probation/supervision length by offense class — GS 15A-1342, GS 15A-1343.2 default ranges"
  - "NC victim services in community supervision — SAVAN notification, victim input at parole hearings"
  - "NC Reentry 2030 strategic plan — EO 303, local reentry councils, 4 goals, housing targets"
  - "NC intensive probation supervision (IPS) — legacy IPS conditions, current risk-based tiers"
  - "NC reentry and transitional housing after incarceration"
whenNotToUse:
  - "Seeking actual legal advice for a specific case — consult a licensed NC attorney"
  - "NC criminal trial, defense, or pre-conviction procedure → nc-criminal-defense"
  - "NC felony/misdemeanor sentencing grids, prior record level → nc-criminal-defense (references/nc-structured-sentencing.md)"
  - "NC record expungement eligibility and process → nc-criminal-defense (references/nc-expungement-record-sealing.md)"
  - "SUD treatment law, MAT/OTP regulation, harm reduction, 42 CFR Part 2, involuntary commitment → nc-drug-rehabilitation-treatment"
  - "Federal probation or parole (non-NC) — no skill available; consult federal resources"
metadata:
  disclaimer: "Educational only, NOT legal advice. For any specific legal matter, consult a licensed North Carolina attorney."
  hub_family: legal
  parent_hub: legal
  spokes:
    - nc-absconder-warrants
    - nc-crv-confinement-response-violation
    - nc-day-reporting-centers
    - nc-drug-court-supervision
    - nc-electronic-monitoring-gps
    - nc-evidence-based-practices-recidivism
    - nc-icaos-interstate-compact
    - nc-intensive-probation-supervision
    - nc-justice-reinvestment-act-jra
    - nc-mapp-mutual-agreement-parole
    - nc-medical-geriatric-parole
    - nc-parole-commission-prspc
    - nc-probation-financial-conditions
    - nc-reentry-2030-strategic-plan
    - nc-reentry-transitional-housing
    - nc-risk-assessment-tools
    - nc-sex-offender-supervision
    - nc-split-sentence-special-probation
    - nc-structured-sentencing-supervision-length
    - nc-victim-services-community-supervision
  changelog:
    - "2026-06-22 sko+pdo v1.2.0->v1.2.2 — 3-iter convergence loop: iter1 G1-High desc 1788->989 chars; iter2 K-Medium JRA acronym-only->full 'Justice Reinvestment Act JRA 2011', drop probation-fee term to stay <=1000; iter3 verified 0 Medium+ findings, all 20 refs linked, em-dash 0/100w, category legal, hub_family legal"
    - "2026-06-22 sko v1.1.1->v1.2.0 — High: rebuilt Hub Contents routing table for all 20 actual reference files (old 11-spoke table referenced non-existent files); updated metadata.spokes to match 20 on-disk references; Medium: category custom->legal; M2: description TRIGGER expanded with 17 new topic terms; M3: whenToUse expanded from 10 to 26 entries covering all new references; M4: keywords expanded from 15 to 45; M6: version bumped; Low: parent_hub nc-criminal-defense->legal; N2: nc-drug-rehabilitation-treatment added to related_skills; Cross-hub SUD treatment redirect added; N3: SKIP updated with SUD treatment path"
    - "2026-06-22 sko v1.0.0->v1.1.0 — Pass H 10/10 pos, 0/10 neg (predicted); 5 Medium fixed (M: desc 1377->979 chars; G: category legal->custom, model/effort added; K: em-dash density 2.15->0/100; A: SKIP targets hub-aware form; N: whenNotToUse federal redirect added); Pass O seeded nc-criminal-defense with nc-parole-probation deferral; hub sync complete (312 skills)"
related_skills:
  - nc-criminal-defense
  - nc-drug-rehabilitation-treatment
  - consumer-credit-and-debt
  - consumer-finance
  - legal
---
# NC Parole, Probation & Community Supervision Hub

North Carolina community supervision: parole, probation, post-release supervision, and the systems that govern them. **Educational only, NOT legal advice.** For any specific legal matter, consult a licensed North Carolina attorney.

## Context

Under NC's **Structured Sentencing Act** (effective Oct 1, 1994), traditional parole was largely abolished for offenses committed after that date. Most people completing active sentences are released under **Post-Release Supervision (PRS)**, a mandatory non-discretionary release period governed by the NC Post-Release Supervision and Parole Commission (PRSPC). Traditional parole remains for pre-1994 offenses and certain categories (DWI, parole-eligible misdemeanors). The **NC Department of Adult Correction (DAC)**, formerly DPS, administers field supervision through probation/parole officers. The 2011 **Justice Reinvestment Act (JRA)** fundamentally restructured probation revocation, creating the CRV system and reducing revocations by ~50%.

## Hub Contents

### Commission, Parole & Release

| Reference | Topic |
|---|---|
| [nc-parole-commission-prspc](references/nc-parole-commission-prspc.md) | PRSPC structure (GS 143B-1490), discretionary parole vs PRS, hearings, conditions, MAPP, medical release |
| [nc-mapp-mutual-agreement-parole](references/nc-mapp-mutual-agreement-parole.md) | MAPP three-party agreement (PRSPC + DAC + inmate), pre-Oct-1994 eligibility, milestone structure, 958 eligible offenders (2023) |
| [nc-medical-geriatric-parole](references/nc-medical-geriatric-parole.md) | GS 15A-1369 (Article 84B): terminally ill, geriatric (age 55+), permanently disabled; 15–45 day PRSPC review; 2024 statistics (4 of 32,000 released) |

### Probation Conditions & Types

| Reference | Topic |
|---|---|
| [nc-split-sentence-special-probation](references/nc-split-sentence-special-probation.md) | GS 15A-1351: active confinement capped at 1/4 max imposed sentence, 72-hour reporting rule, 5-year max probation period, DWI rules |
| [nc-intensive-probation-supervision](references/nc-intensive-probation-supervision.md) | IPS formally repealed by JRA for post-Dec-2011 sentences; 1,700+ legacy IPS offenders; current intensive supervision via risk tiers/NC-PASE; curfew delegation GS 15A-1343.2(e)(f) |
| [nc-structured-sentencing-supervision-length](references/nc-structured-sentencing-supervision-length.md) | GS 15A-1342 max 5-year felony probation; GS 15A-1343.2 default ranges by offense class/punishment type; extension requires judicial findings |
| [nc-probation-financial-conditions](references/nc-probation-financial-conditions.md) | Supervision fee $40/month (GS 15A-1343(c1)); EM device fee $90 one-time + daily cost; restitution; fee waiver good cause standard; 2025 court costs chart |

### Violations & Sanctions

| Reference | Topic |
|---|---|
| [nc-crv-confinement-response-violation](references/nc-crv-confinement-response-violation.md) | CRV: 90-day confinement for technical violations; two-CRV limit; JRA 2011 rules (GS 15A-1344); Burke/Robeson/North Piedmont CRV centers; felony vs misdemeanor CRV |
| [nc-absconder-warrants](references/nc-absconder-warrants.md) | Absconding defined; ACE absconder capture enforcement teams; SOIU Special Operations and Intelligence Unit; ~12,000 active warrants statewide; 1-888-646-0024 tip line |
| [nc-drug-court-supervision](references/nc-drug-court-supervision.md) | JMARC (Judicially Managed Accountability and Recovery Courts); GS 7A-272 jurisdiction; probation conditions GS 15A-1343(b1)(2b); 4-phase structure; graduation; deferred prosecution GS 15A-1341 |

### Electronic Monitoring & Technology

| Reference | Topic |
|---|---|
| [nc-electronic-monitoring-gps](references/nc-electronic-monitoring-gps.md) | House arrest with EM (GS 15A-1343); GPS probation/parole; SCRAM alcohol monitoring; vendor landscape; device fees; inclusion/exclusion zones; tamper detection; NC DOA contract 915C |

### Sex Offender Supervision

| Reference | Topic |
|---|---|
| [nc-sex-offender-supervision](references/nc-sex-offender-supervision.md) | SORNA tier system (GS 14-208 et seq.); sheriff registration; residency/employment restrictions; satellite-based monitoring (SBM, GS 14-208.40); Static-99 risk assessment; lifetime supervision |

### Interstate & Policy

| Reference | Topic |
|---|---|
| [nc-icaos-interstate-compact](references/nc-icaos-interstate-compact.md) | ICAOS transfers across state lines; NC DAC Interstate Compact Office; mandatory vs discretionary transfers; sending/receiving state roles; travel permits; GS 148 Article 4B |
| [nc-justice-reinvestment-act-jra](references/nc-justice-reinvestment-act-jra.md) | JRA 2011 (SL 2011-192): CRV creation, revocation limits, $560M savings, 10 prisons closed, 175 PO positions funded, revocations down 50%; 2025 SPAC evaluation |

### Evidence-Based Practices & Assessment

| Reference | Topic |
|---|---|
| [nc-risk-assessment-tools](references/nc-risk-assessment-tools.md) | NC-PASE as primary risk classification tool; Static-99 for SBM; RNR model; criminogenic needs domains; risk levels (high/moderate/low); supervision frequency; COMPAS/LSI-R background |
| [nc-evidence-based-practices-recidivism](references/nc-evidence-based-practices-recidivism.md) | RNR model, NC-PASE predictive analytics, Carey Guides, Motivational Interviewing, Cognitive Behavioral Intervention (CBI), 8 EBP principles, NC DAC EBP program, SPAC JRA outcomes |

### Intermediate Sanctions & Programming

| Reference | Topic |
|---|---|
| [nc-day-reporting-centers](references/nc-day-reporting-centers.md) | DRC/ADRC: intensive community supervision alternative; CBI; substance abuse services; GED/employment readiness; urine drug screens; cost-effective jail alternative; county-operated programs |

### Victim Services

| Reference | Topic |
|---|---|
| [nc-victim-services-community-supervision](references/nc-victim-services-community-supervision.md) | NC SAVAN automated victim notification; DAC Victim Support Services; 12 mandatory notification events; GS 15A-837; victim input at parole hearings; restitution enforcement; ncsavan.org |

### Reentry

| Reference | Topic |
|---|---|
| [nc-reentry-2030-strategic-plan](references/nc-reentry-2030-strategic-plan.md) | Executive Order 303 (Jan 2024): 4 goals, 26 objectives, 133 strategies; LRCs (28 of 100 counties → all 100); 1,800 housing units/year target; DAC Rehabilitation and Reentry Division; 2025 progress report |
| [nc-reentry-transitional-housing](references/nc-reentry-transitional-housing.md) | DAC Division of Rehabilitation and Reentry; Reentry 2030 initiative; local reentry councils; recovery housing; barriers to housing; Second Chance Act; NC Statewide Reentry Council |

---

## Key Resources

- [NC Department of Adult Correction](https://www.dac.nc.gov/) (official DAC site)
- [NC Post-Release Supervision and Parole Commission](https://www.dac.nc.gov/divisions/post-release-supervision-and-parole-commission) (Commission rules and procedures)
- [NC General Statutes Chapter 15A](https://www.ncleg.gov/Laws/GeneralStatuteSections/Chapter15A) (Criminal Procedure Act)
- [NC General Statutes Chapter 15C / 15E](https://www.ncleg.gov/Laws/GeneralStatuteSections/Chapter15C) (monitoring/supervision law)
- [UNC School of Government NC Criminal Law Blog](https://nccriminallaw.sog.unc.edu/) (case law updates)
- [NC Sentencing Commission](https://www.nccourts.gov/commissions/sentencing-and-policy-advisory-commission) (data, reports, training)
- [ICAOS Rules & Bench Book](https://www.interstatecompact.org/) (interstate compact resources)

## Cross-Hub Notes

- **Pre-conviction supervision** (pretrial release, bail, bond conditions) → `nc-criminal-defense`
- **Sentencing grids and prior record level** → `nc-criminal-defense` (references/nc-structured-sentencing.md)
- **Record expungement/sealing after supervision ends** → `nc-criminal-defense` (references/nc-expungement-record-sealing.md)
- **Employment/housing barriers post-supervision** may also touch `consumer-finance` (budgeting, benefits) or `consumer-credit-and-debt` (debt from supervision fees)
- **SUD treatment law** (drug court program law, MAT/OTP regulation, harm reduction statutes GS 90-96.2/90-113.27, 42 CFR Part 2 confidentiality, involuntary commitment) → `nc-drug-rehabilitation-treatment`
