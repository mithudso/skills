---
name: legal
version: "1.3.0"
updated: "2026-06-22"
category: legal
model: claude-sonnet-4-6
effort: medium
description: "Hub for NC criminal law, drug law, criminal defense, rehabilitation, parole/probation, bail/bond, and recording laws; routes to spokes. Educational, NOT legal or attorney advice (as of 2026). TRIGGER: NC criminal charges or arrest; drug schedules (GS Ch.90); bail/bond/pretrial release (GS 15A-531-547); NC criminal defense (indigent defense, open-file discovery, plea); structured sentencing; NC parole/probation/post-release supervision; NC drug rehabilitation, MAT, harm reduction, naloxone, syringe exchange; 42 CFR Part 2; Good Samaritan overdose immunity; bond forfeiture, bail bondsman licensure; NC recording laws, one-party consent, GS 15A-287, recording police NC, wiretapping NC. SKIP: family/estate/civil/tax/immigration/contract law → consult a licensed NC attorney; consumer debt → consumer-credit-and-debt; NC business/nonprofit → venture-business / venture-nonprofit-cause; NC real estate → venture-business (references/venture-nc-real-estate-law.md)."
tags:
  - north-carolina
  - criminal-law
  - drug-law
  - bail-bond
  - parole
  - probation
  - criminal-defense
  - drug-rehabilitation
  - legal
whenToUse:
  - "I was arrested in NC — what are my rights"
  - "NC drug charge schedules or trafficking thresholds"
  - "bail, bond, or pretrial release conditions in NC"
  - "NC structured sentencing grid, felony class, or prior record level"
  - "NC public defender or indigent defense eligibility"
  - "open-file discovery in NC criminal cases"
  - "NC probation conditions, violation, or revocation"
  - "NC parole or post-release supervision"
  - "NC drug court, MAT, buprenorphine, or methadone clinic"
  - "NC naloxone, Good Samaritan law, or syringe exchange"
  - "NC bail bondsman license or bond forfeiture"
  - "NC pretrial services or cash-bail reform"
triggers:
  - "NC criminal"
  - "North Carolina arrest"
  - "drug schedule NC"
  - "bail bond NC"
  - "NC probation"
  - "NC parole"
  - "NC public defender"
  - "indigent defense NC"
  - "NC drug court"
  - "NC naloxone"
  - "Good Samaritan NC"
  - "syringe exchange NC"
  - "bond forfeiture NC"
  - "bail bondsman NC"
  - "pretrial release NC"
  - "structured sentencing NC"
  - "GS 15A"
  - "GS 90 drug"
  - "42 CFR Part 2"
whenNotToUse: "Family law, estate planning & wills, civil litigation, tax law, immigration, contract law — no NC spokes exist; consult a licensed NC attorney. Consumer debt/credit law → consumer-credit-and-debt. NC business/nonprofit → venture-business (references/venture-nc-business-formation-tax.md) / venture-nonprofit-cause (references/venture-nc-nonprofit-formation.md). NC real estate law → venture-business (references/venture-nc-real-estate-law.md)."
related_skills:
  - nc-criminal-defense
  - nc-drug-rehabilitation-treatment
  - nc-parole-probation
  - nc-recording-laws
  - consumer-credit-and-debt
metadata:
  changelog:
    - "2026-06-22 sko v1.2.0->v1.3.0 — 2 Medium: aligned drug-treatment-court (GS 7A-790) & involuntary-SUD-commitment (GS 122C-281) cites to spoke; model/effort opus/high->sonnet/medium (routing hub, light judgment); refs/->references/. Pass H 10/10 pos predicted."
    - "2026-06-22 sko v1.0.0->v1.1.0 — 4 Medium fixes: added model/effort (Step 4.6), corrected venture-nc-* SKIP targets to hub-aware form, removed duplicate ASCII domain-map tree, fixed category custom->legal. Pass H 10/10 pos, 0/10 neg (predicted)."
---
# NC Legal (hub)

Front door for **North Carolina criminal law, drug law, defense, rehabilitation, parole/probation, and bail/bond**. Routes to 3 top-level spokes and their reference sub-files.

> **Educational only — NOT legal or attorney advice.** NC statutes, court rules, and sentencing guidelines change. All spokes stamped *as of 2026*. For any specific case, charge, or proceeding: consult a **licensed NC attorney**.

**When activated:** identify the best-fit spoke from the routing table and activate it directly. Multi-domain question → show routing table and ask one clarifying question (e.g., "Is this about the criminal charge itself, supervision conditions, or treatment access?").

---

## Routing table

### Top-level spokes

| Spoke | Owns | Activate when |
|---|---|---|
| `nc-criminal-defense` | Arrest/charge process, rights, plea, trial, Alford plea, expungement, NC criminal procedure (GS 15A), 4th/5th/6th Amendment | User faces NC criminal charge or asks about defense process |
| `nc-drug-rehabilitation-treatment` | MAT/OTP regulation (42 CFR Part 8), buprenorphine, SAMHSA, NC DHHS SUD licensing, drug treatment court (GS 7A-790), involuntary SUD commitment (GS 122C-281), 42 CFR Part 2 confidentiality, Medicaid/LME-MCO, peer support | User asks about NC addiction treatment, MAT, or SUD system |
| `nc-parole-probation` | NC probation conditions/violation/revocation, parole, post-release supervision (PRS), DAC/OPUS, electronic monitoring, swift-certain sanctions, reentry | User on or facing NC supervision |
| `nc-recording-laws` | One-party consent (GS 15A-287), criminal/civil penalties, recording police (First Amendment/Sharpe), attorney-client recording, AI/Plaud devices, video/voyeurism (GS 14-202), body cameras, ECPA, courthouse rules | User asks about NC recording laws, wiretapping, surveillance, or audio/video consent |

### Reference sub-files

**`nc-criminal-defense/references/`:**

| Reference | Owns |
|---|---|
| `nc-bail-and-bond.md` | Bail types, pretrial release conditions, GS 15A-531–547, secured/unsecured bond, DV bond (GS 15A-534.1), preventive detention, Iryna's Law (SL 2025-93) |
| `nc-bail-bondsmen-licensure.md` | Bail bondsman NC DOI license, surety bail, surrender/off-bond, GS Ch. 58 Art. 71 |
| `nc-bond-forfeiture-remission.md` | Bond forfeiture process, remission/set-aside, FFA, estreatment |
| `nc-pretrial-services-reform.md` | Pretrial services agencies, risk assessment tools, cash-bail reform debate, cite-and-release |
| `nc-controlled-substance-schedules.md` | GS Ch.90 Schedules I–VI, possession vs PWIMSD, trafficking thresholds, mandatory minimums, drug paraphernalia |
| `nc-structured-sentencing.md` | Felony/misdemeanor grids, prior record level, dispositional ranges, split sentences, habitual felon |
| `nc-open-file-discovery.md` | NC open-file discovery model, Brady/Giglio obligations, GS 15A-903 |
| `nc-indigent-defense.md` | NC Indigent Defense Services (IDS), public defender eligibility, assigned counsel rates |

**`nc-drug-rehabilitation-treatment/references/`:**

| Reference | Owns |
|---|---|
| `nc-harm-reduction-laws.md` | NC Good Samaritan law (GS 90-96.2), naloxone standing order (GS 90-12.7), syringe exchange (GS 90-113.27), fentanyl test strips |

---

## Coverage gaps

The following NC legal domains have no spokes — direct to a licensed NC attorney:
- Family law (divorce, custody, child support)
- Estate planning, wills, probate
- Civil litigation & tort law
- NC tax law
- Immigration law
- Contract & business law
- NC real estate law → see `venture-business` (references/venture-nc-real-estate-law.md) for formation context
