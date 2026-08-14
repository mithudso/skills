# nc-drug-rehabilitation-treatment

**Category:** Psychology & Behavioral Science
**Platform:** Claude Code
**Original Path:** claude-code/nc-drug-rehabilitation-treatment

## Description
NC substance use disorder (SUD) treatment law hub — harm reduction, MAT, drug courts, confidentiality, involuntary commitment, Medicaid, facility licensing, collateral consequences, opioid settlement, recovery housing. NOT legal advice. TRIGGER: NC Good Samaritan GS 90-96.2; naloxone standing order GS 90-12.7; syringe exchange GS 90-113.27; fentanyl test strips NC; 42 CFR Part 2 SUD records; NC drug treatment courts GS 7A-790; involuntary SUD commitment GS 122C-281; MAT/OTP 42 CFR Part 8; methadone clinic NC; opioid treatment program NC; SAMHSA accreditation NC; DHHS SUD facility licensing; NC LME-MCO Medicaid behavioral health; MHPAEA parity SUD; buprenorphine X-waiver repeal; NC peer support specialist CPSS; SUD collateral consequences NC; SUD custody/CPS NC; NC opioid settlement funds MOA; NC recovery housing NCARR. SKIP: drug charge trial defense → nc-criminal-defense; probation/parole → nc-parole-probation; private health insurance → consumer-finance.

---

# NC Drug Rehabilitation & Treatment Law (Hub)

North Carolina law governing substance use disorder treatment, harm reduction, and recovery. **Educational only — NOT legal advice.** Rules change; verify with NC DHHS, NC courts, or a licensed NC attorney before acting.

## Routing Table

| Spoke | Use when | Reference file |
|---|---|---|
| `nc-harm-reduction-laws` | Good Samaritan overdose immunity (GS 90-96.2), naloxone standing orders (GS 90-12.7), syringe exchange programs (GS 90-113.27), fentanyl test strip status | [references/nc-harm-reduction-laws.md](references/nc-harm-reduction-laws.md) |
| `nc-mat-otp-law` | Medication-assisted treatment law — 42 CFR Part 8, OTP certification, methadone clinics, X-waiver elimination, buprenorphine prescribing rules, DEA telehealth (2025 final rule), naltrexone/Vivitrol | [references/nc-mat-otp-law.md](references/nc-mat-otp-law.md) |
| `nc-drug-treatment-courts` | NC Recovery Courts (Drug Treatment Courts) — GS 7A-790 et seq. (Article 62), therapeutic jurisprudence, adult/DWI/veterans/mental health/family/juvenile court types, eligibility, program phases, deferred prosecution (GS 15A-1341(a2)), discharge/dismissal (GS 15A-1341(a5)), structured sentencing interaction, LME-MCO funding | [references/nc-drug-treatment-courts.md](references/nc-drug-treatment-courts.md) |
| `nc-sud-confidentiality` | 42 CFR Part 2 SUD record confidentiality vs HIPAA, single TPO consent (2024 rule), SUD counseling notes separate consent, court order requirements for law enforcement disclosure, NC GS 122C-52/53 (S.L. 2023-95), compliance deadline February 2026 | [references/nc-sud-confidentiality.md](references/nc-sud-confidentiality.md) |
| `nc-involuntary-sud-commitment` | Involuntary SUD commitment in NC — GS 122C-281–288 (Part 8; §122C-280 is RESERVED), petition by anyone, two-part criteria (substance abuser + dangerous), 24hr examination, 10-day district court hearing, 90-day inpatient / 180-day outpatient commitment, patient rights, appeal to Court of Appeals | [references/nc-involuntary-sud-commitment.md](references/nc-involuntary-sud-commitment.md) |
| `nc-behavioral-health-medicaid` | NC LME-MCO system (Alliance, Partners, Trillium, Vaya), Tailored Plans (July 1 2024) for severe SUD/MH/IDD/TBI, ASAM levels of care (SAIOP 2.1, SACOT 2.5, residential 3.1–3.7), 1115 SUD Demonstration Waiver (CCPs eff. Jan 1 2026), prior authorization, MHPAEA parity compliance, Children and Families Specialty Plan, NCCARE360 referral platform | [references/nc-behavioral-health-medicaid.md](references/nc-behavioral-health-medicaid.md) |
| `nc-sud-facility-licensing` | NC SUD treatment facility licensing — GS 122C, DHSR MHLCS, 10A NCAC 27G service categories (SAIOP .4400, SACOT .4500, residential .3400, detox .3200, OTP .3600, medically monitored .3700), Certificate of Need, OTP dual-track approval (SOTA → DHSR → DEA → SAMHSA → accreditation), deemed status GS 122C-81, Medicaid enrollment separate | [references/nc-sud-facility-licensing.md](references/nc-sud-facility-licensing.md) |
| `nc-peer-support-specialist` | NC Certified Peer Support Specialist (CPSS) — NCCPSS Program at UNC-Chapel Hill SSW, 50-hour training requirement, 3-hour ethics training (eff. July 1 2024), 2-year recertification, scope of practice, Medicaid billing via LME-MCO/licensed org only, SUD-specific roles (MAT, harm reduction, recovery courts), recovery coach vs CPSS distinction | [references/nc-peer-support-specialist.md](references/nc-peer-support-specialist.md) |
| `nc-drug-collateral-consequences` | Collateral consequences of NC drug convictions — SNAP tiered ban (6-mo Class H/I; lifetime Class G+), Work First/TANF bar, HB 564 reform, federal public housing (42 USC 13661), occupational licensing (GS 93B-8.1 — nexus test, no automatic denial, 30-day response right), Certificate of Relief (GS 15A-173.2), FAFSA Simplification Act (aid bars eliminated July 2023), voting rights restoration on supervision completion | [references/nc-drug-collateral-consequences.md](references/nc-drug-collateral-consequences.md) |
| `nc-sud-family-law` | SUD in NC family law — GS Chapter 50 custody (best interests), GS Chapter 7B AND proceedings, substance-exposed newborns + Plan of Safe Care (CARA Act/CAPTA, GS 7B-302), TPR grounds (GS 7B-1111(a)(1) neglect, (a)(6) SUD incapability, (a)(2) foster care 12 months), Family Treatment Court interface, 42 CFR Part 2 in AND proceedings | [references/nc-sud-family-law.md](references/nc-sud-family-law.md) |
| `nc-opioid-settlement-funds` | NC opioid settlement funds — NC MOA (85% to counties/municipalities, 15% state), Exhibit A high-impact uses (naloxone, MAT, recovery housing, peer support, syringe services, post-overdose response, criminal justice diversion), Exhibit B broader uses (requires collaborative planning), CORE-NC transparency platform, state (NCDHHS/UNC) 15% spending | [references/nc-opioid-settlement-funds.md](references/nc-opioid-settlement-funds.md) |
| `nc-recovery-housing-law` | NC recovery housing — no state licensure requirement for basic sober living, NCARR certification (NARR state affiliate; 4 levels: peer-run → supervised), Oxford House NC (~320 houses), FHA/ADA disability protections (limit local zoning restrictions), when recovery housing crosses into DHSR licensure threshold (GS 122C-3(14)(b)), Medicaid/opioid settlement funding pathways | [references/nc-recovery-housing-law.md](references/nc-recovery-housing-law.md) |

## Cross-Hub Map

- **Criminal defense procedure** (NC drug charges, controlled substances act, trafficking, Good Samaritan as a defense at trial) → `nc-criminal-defense`
- **Supervised release violations** (probation revocation for drug use, post-release supervision, drug court supervision) → `nc-parole-probation`
- **Personal finance / health insurance** (private health plan coverage, deductibles, HSA) → `consumer-finance`
- **Health behavior change** (behavioral models for addiction treatment engagement) → `health-behavior-change-and-donor-registration`

## Key NC Statutes

| Statute | Topic |
|---|---|
| GS 90-96.2 | Good Samaritan / overdose immunity |
| GS 90-12.7 | Naloxone / opioid antagonist authority |
| GS 90-113.27 | Syringe exchange programs |
| GS 122C-281 – 122C-287 | Involuntary SUD commitment (Part 8, Chapter 122C; §122C-280 is reserved) |
| GS 7A-790 et seq. (Article 62) | Drug Treatment Courts / NC Recovery Courts |
| GS 90-113.22 / GS 90-113.23 | Drug paraphernalia definition and offenses |

## Key Federal Law

| Regulation/Law | Topic |
|---|---|
| 42 CFR Part 2 | SUD record confidentiality (stricter than HIPAA) |
| 42 CFR Part 8 | OTP certification and methadone regulation |
| MHPAEA (2008, amended 2024) | Mental health/SUD parity in insurance |
| ADA / Section 504 | Disability protections for persons in recovery |
| SUPPORT Act (2018, reauthorized) | Federal SUD funding authorization |