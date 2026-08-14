# consumer-credit-and-debt

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude
**Original Path:** claude/standalone/consumer-credit-and-debt

## Description
Hub for US consumer credit & debt — credit health, lending, collections, and governing law (federal + NC); routes to 11 spokes. Educational, NOT financial or legal advice (as of 2026). TRIGGER: credit reports or scores (FICO/VantageScore, bureaus); building or rebuilding credit; charge-offs, collections, settling debt, 1099-C; a collector contacting/suing me (FDCPA); home-mortgage or auto lending; predatory/high-cost loans (payday/title/rent-to-own); identity theft & credit fraud; consumer bankruptcy (Ch.7/13); US credit statutes (FCRA/FDCPA/ECOA/TILA/FCBA/CROA/EFTA) or NC debt law. SKIP: personal/household finance — student loans, banking, income taxes, insurance (personal & health), medical bills, budgeting, investing/retirement, estate planning & wills → consumer-finance (sibling hub); NC business/nonprofit/real-estate → venture-nc-business-formation-tax / venture-nc-nonprofit-formation / venture-nc-real-estate-law; business credit/D&B; MongoDB / TAM / data-analytics → own hubs.

---

# Consumer Credit & Debt (hub)

Front door for **US personal consumer credit and debt**, federal + **North Carolina**. Routes to 11 spokes.

**When activated:** present routing table, identify best spoke, activate directly. Question names sub-topic (e.g., "FDCPA rights," "Chapter 7," "credit score") → activate matching spoke, skip hub summary. Broad or multi-spoke question → show routing table, ask one clarifying question.

> **Educational only — not financial, tax, or legal advice.** Statutes, rules, scoring models, dollar figures change and are fact-specific. Every spoke stamped *as of 2026* with "verify current law / consult licensed professional" pointer. For specific dispute, lawsuit, tax question, or filing: consult licensed attorney, HUD-approved counselor, or tax professional.

## How the family fits together (the cycle)

1. **Record:** credit report (bureaus collect)
2. **Rating:** credit score (derived from record)
3. **Damage events:** late payments, charge-offs, collections hit record
4. **Repair:** dispute errors, rebuild score, resolve debts
5. **Lenders:** price home and auto loans off rating
6. **Enforcement:** collectors and sometimes bankruptcy when things break down
7. **Law:** federal statutes + NC overlays govern every step above

Two law spokes (`us-consumer-credit-and-debt-law`, `north-carolina-credit-and-debt-law`) are backbone the topical spokes reference.

## Routing table

Each spoke = on-demand reference file under `references/`. Row matches → **Read the listed `references/` file** before answering at depth — table is router only.

| Spoke | Use when | Reference file |
| --- | --- | --- |
| `credit-reports-and-scores` | How report/score built; bureaus; FICO vs VantageScore; inquiries; freeze vs lock; item longevity. | `references/credit-reports-and-scores.md` |
| `improving-and-rebuilding-credit` | Raise/build/rebuild score; utilization; dispute errors; secured cards; rebuild after damage. | `references/improving-and-rebuilding-credit.md` |
| `charge-offs-collections-and-debt-resolution` | Charge-off meaning; settling; pay-for-delete; debt buyers; 1099-C canceled-debt tax; re-aging traps. | `references/charge-offs-collections-and-debt-resolution.md` |
| `debt-collectors-and-fdcpa-rights` | Collector calling/suing; debt validation; cease-contact; time-barred "zombie" debt; answering summons. | `references/debt-collectors-and-fdcpa-rights.md` |
| `home-mortgage-lending` | Mortgage types; score/DTI to qualify; PMI vs MIP; TRID process; NCHFA assistance; refinancing. | `references/home-mortgage-lending.md` |
| `auto-lending-and-financing` | Car loans; APR by credit tier; dealer vs credit-union financing; negative equity; GAP; repossession. | `references/auto-lending-and-financing.md` |
| `predatory-lending-and-high-cost-credit` | Payday/title/rent-to-own/pawn/"no-credit-check"; spotting predatory terms; rate caps; safe alternatives. | `references/predatory-lending-and-high-cost-credit.md` |
| `identity-theft-and-credit-fraud` | Someone used your identity; freezes/fraud alerts; IdentityTheft.gov recovery; FCRA 605B; tax ID theft. | `references/identity-theft-and-credit-fraud.md` |
| `bankruptcy-ch7-ch13` | Should you file; Ch.7 vs Ch.13; means test; automatic stay; dischargeable debts; alternatives. | `references/bankruptcy-ch7-ch13.md` |
| `us-consumer-credit-and-debt-law` | Federal statute text & rights: FCRA, FDCPA + Reg F, ECOA, TILA, FCBA, CROA, EFTA; CFPB/FTC enforcement. | `references/us-consumer-credit-and-debt-law.md` |
| `north-carolina-credit-and-debt-law` | NC overlays: NC Debt Collection Act (reaches original creditors), 3-yr SOL, wage-garnishment ban, exemptions, power-of-sale foreclosure, repossession, payday ban. | `references/north-carolina-credit-and-debt-law.md` |

## Cross-cutting notes

- **Sibling hub — `consumer-finance`.** Personal/household-finance side — student loans, personal banking, individual income taxes, personal & health insurance, medical bills/billing, budgeting & saving, investing & retirement, estate planning & wills — lives in **`consumer-finance`** hub; route there. Two hubs interlock: defaulted student loan or unpaid medical bill becomes credit/collections matter *here*; bankruptcy here can discharge debts from there. This hub owns everything credit-, debt-, collections-, lending-, and consumer-credit-law-related.
- **North Carolina is material overlay, not footnote.** NC bars wage garnishment for ordinary consumer debt; Debt Collection Act reaches *original* creditors (broader than federal FDCPA); SOL on most debt ~3 years (G.S. 1-52); effectively bans payday lending. NC-specific question → route to **north-carolina-credit-and-debt-law** alongside topical spoke.
- **Federal vs state:** **us-consumer-credit-and-debt-law** for statute text and federal rights; NC spoke for what changes in North Carolina.
- **Volatile area flagged:** 2024 CFPB medical-debt-off-credit-reports rule **vacated 2025** — US-law spoke carries current status; re-verify before relying.

## References

Each spoke maintains own "References / verify current law" section citing primary authorities (CFPB, FTC, US Code / eCFR, IRS, HUD/FHFA/VA, US Courts, NC General Statutes at ncleg.gov, NC DOJ). Start from relevant spoke.

<!-- cross-hub-map -->
## Cross-hub map — where every consumer-credit-and-debt topic lives

Family split across these hubs. Deep material **not** in this hub's routing table → reference file under sibling hub below — **activate that hub or `Read` its `references/<name>.md` directly**. Every former standalone skill now lives as reference under one of these hubs (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `consumer-credit-and-debt` | US Consumer Credit & Debt (hub) | `references/credit-reports-and-scores.md`, `references/improving-and-rebuilding-credit.md`, `references/charge-offs-collections-and-debt-resolution.md`, `references/debt-collectors-and-fdcpa-rights.md`, … |