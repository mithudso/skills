# legal

**Category:** Cloud, DevOps & Infrastructure
**Platform:** Claude Code
**Original Path:** claude-code/legal

## Description
Hub for NC criminal law, drug law, criminal defense, rehabilitation, parole/probation, bail/bond, and recording laws; routes to spokes. Educational, NOT legal or attorney advice (as of 2026). TRIGGER: NC criminal charges or arrest; drug schedules (GS Ch.90); bail/bond/pretrial release (GS 15A-531-547); NC criminal defense (indigent defense, open-file discovery, plea); structured sentencing; NC parole/probation/post-release supervision; NC drug rehabilitation, MAT, harm reduction, naloxone, syringe exchange; 42 CFR Part 2; Good Samaritan overdose immunity; bond forfeiture, bail bondsman licensure; NC recording laws, one-party consent, GS 15A-287, recording police NC, wiretapping NC. SKIP: family/estate/civil/tax/immigration/contract law → consult a licensed NC attorney; consumer debt → consumer-credit-and-debt; NC business/nonprofit → venture-business / venture-nonprofit-cause; NC real estate → venture-business (references/venture-nc-real-estate-law.md).

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