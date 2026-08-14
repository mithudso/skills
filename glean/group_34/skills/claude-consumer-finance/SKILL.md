---
name: consumer-finance
version: "1.1.0"
updated: "2026-06-16"
category: custom
metadata:
  changelog:
    - "2026-06-16 sko v1.0.0->v1.1.0: Pass H 10/10 pos, 1/10 neg (predicted); fixed 5 Medium (M: desc trimmed 1064->945 chars; N: SKIP IDs canonical + venture-* expanded; B: insurance/mortgage seam note; F: Not-this-hub body redirect; K: em-dash density); 0 High"
description: "Hub for US personal/consumer finance — money management, coverage, and life planning; routes to 9 spokes. Educational only, NOT advice (as of 2026). TRIGGER: student loans (IDR/SAVE/PSLF, FAFSA, default); personal banking (FDIC/NCUA, HYSA, neobanks); individual income taxes (brackets, deductions, W-4, self-employment, NC); personal insurance (auto, home/renters, life, disability); health insurance (ACA, HSA, Medicare, Medicaid); medical bills & billing disputes; budgeting & saving; investing & retirement (401k, IRA, Roth, Social Security); estate planning & wills (POA, NC probate). SKIP: credit reports/scores, building credit, charge-offs, debt collectors, mortgages, auto loans, predatory lending, identity theft, bankruptcy, or US/NC credit-debt LAW → consumer-credit-and-debt; NC business/nonprofit/payroll → venture-nc-business-formation-tax, venture-nc-employer-payroll, venture-nc-nonprofit-formation; MongoDB/TAM → their own hubs."
tags:
  - consumer-finance
  - personal-finance
  - student-loans
  - banking
  - taxes
  - insurance
  - medical-debt
  - retirement
  - estate-planning
  - north-carolina
whenToUse:
  - "how do student loan repayment plans / forgiveness work"
  - "is my money safe at this bank / neobank (FDIC vs NCUA)"
  - "what tax bracket am I in / can I take this deduction or credit"
  - "how much car / home / life insurance do I need"
  - "how does an ACA marketplace plan or HSA work"
  - "I got a huge medical bill — what do I do"
  - "help me make a budget / build an emergency fund"
  - "how do I start investing / which retirement account"
  - "do I need a will / what happens in NC if I die without one"
related_skills:
  - consumer-credit-and-debt
---
# Consumer Finance (hub)

Front door for **US personal/household finance** — money management, coverage, life planning. Routes to 9 spokes. **Sibling hub `consumer-credit-and-debt`** owns credit, debt, collections, consumer-protection-law side; route there for credit scores, borrowing problems, collectors, governing statutes.

**Not this hub:** credit reports/scores, building credit, charge-offs, debt collectors (FDCPA), mortgages, auto loans, predatory lending, identity theft, bankruptcy, US/NC credit-&-debt law → `consumer-credit-and-debt`. NC business/employer/nonprofit → `venture-*` skills.

> **Educational only — not financial, tax, legal, or insurance advice.** Rules, limits, brackets, program availability change and are fact-specific. Every spoke date-stamped *as of 2026* with "verify current" pointers. For specific decisions, consult licensed professional (CPA/tax pro, fee-only fiduciary, licensed agent, or attorney).

## Routing table

Each spoke = on-demand reference under `references/`. When row matches, **Read the listed `references/` file** before answering at depth — table is router only.

| Spoke | Use when | Reference file |
| --- | --- | --- |
| `student-loans` | Federal vs private; IDR/SAVE/IBR/PAYE & new RAP; PSLF/forgiveness; deferment; default, rehab, consolidation; FAFSA; refinancing tradeoffs. | `references/student-loans.md` |
| `personal-banking` | Checking/savings/HYSA/CD; FDIC vs NCUA insurance; overdraft & Reg E; bank vs credit union; neobank/fintech safety; ChexSystems; Zelle/wire scams. | `references/personal-banking.md` |
| `personal-income-taxes` | Individual income tax: brackets, standard vs itemized, credits (EITC/CTC/education), W-4, self-employment/estimated, capital-gains basics, IRS process, NC flat tax. | `references/personal-income-taxes.md` |
| `personal-insurance` | Property/casualty + life + disability: auto (NC minimums), home/renters, replacement cost vs ACV, flood/NFIP, term vs whole life, umbrella. | `references/personal-insurance.md` |
| `health-insurance-and-coverage` | Plan mechanics (deductible/copay/OOP max/network); ACA marketplace & subsidies; HDHP+HSA; Medicare A–D & Medigap; Medicaid (NC expansion); COBRA; appeals. | `references/health-insurance-and-coverage.md` |
| `medical-debt-and-billing` | Understanding/disputing medical bill; No Surprises Act; hospital charity care (501(r)); negotiating; medical debt & credit reports; CareCredit traps. | `references/medical-debt-and-billing.md` |
| `budgeting-and-saving` | Build/fix budget (50/30/20, zero-based, envelope); emergency fund; automate savings; snowball vs avalanche payoff; free nonprofit counseling. | `references/budgeting-and-saving.md` |
| `investing-and-retirement` | Investing basics (diversification, index funds, fees, DCA); 401(k)/IRA/Roth/backdoor; Social Security claiming; spotting investment fraud; vetting advisor. | `references/investing-and-retirement.md` |
| `estate-planning-and-wills` | Will, financial & health-care POA, living will; NC intestacy; revocable trust vs will; beneficiary designations; NC probate / small-estate affidavit. | `references/estate-planning-and-wills.md` |

## Cross-cutting notes

- **Sibling hub — `consumer-credit-and-debt`.** Credit reports/scores, building credit, charge-offs/collections, debt collectors (FDCPA), mortgages & auto loans, predatory lending, identity theft, bankruptcy, US/NC credit-&-debt **law** live there. Hubs interlock: defaulted student loan or unpaid medical bill becomes credit/collections matter there; bankruptcy can discharge some debts. This hub covers **insurance** on auto/home (coverage); sibling covers **auto loans and mortgages** (financing/debt). Questions spanning both: start here for coverage, route there for loan.
- **North Carolina overlays** recur: NC flat income tax (`personal-income-taxes`), NC auto minimums & coastal wind/Beach Plan (`personal-insurance`), NC probate & intestacy (`estate-planning-and-wills`), NC wage-garnishment ban (sibling hub's NC-law spoke).
- **Volatile areas (re-verify before relying):** student-loan forgiveness/SAVE & new RAP plan (2026 flux); ACA enhanced subsidies (status uncertain as of 2026); medical-debt credit reporting (2024 CFPB rule vacated 2025, verify current status).
- **Boundary with `venture-*`:** those skills own **business/employer/nonprofit** formation, tax, payroll; this hub strictly **personal/household** finance.

## References

Each spoke maintains own "References / verify current" section citing primary authorities (IRS, NCDOR, studentaid.gov, FDIC/NCUA, CFPB, FTC, healthcare.gov, CMS/Medicare/Medicaid, SEC/FINRA/DOL/SSA, NC General Statutes at ncleg.gov, NC Judicial Branch). Start from relevant spoke.

<!-- cross-hub-map -->
## Cross-hub map — where every consumer-finance topic lives

Family split across these hubs. If deep material **not** in this hub's routing table, it's a reference under sibling hub below — **activate that hub or `Read` its `references/<name>.md` directly**. Every former standalone skill now a reference under one of these hubs (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `consumer-finance` | US Consumer / Personal Finance (hub) | `references/student-loans.md`, `references/personal-banking.md`, `references/personal-income-taxes.md`, `references/personal-insurance.md`, … |