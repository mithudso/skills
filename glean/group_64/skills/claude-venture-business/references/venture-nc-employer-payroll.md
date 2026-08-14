<!-- hub-reference-banner -->
> **Reference file — part of the `venture-business` hub.** Formerly the standalone `venture-nc-employer-payroll` skill.
> Sibling topics in this family are now reference files under the hubs (`venture-business`, `venture-nonprofit-cause`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: venture-nc-employer-payroll
description: >-
  Becoming an employer and running payroll in North Carolina for a lean for-profit or
  nonprofit hiring its first W-2 employee(s): the federal + NC setup, taxes, filings,
  and labor-law duties, with a first-hire checklist. TRIGGER: hiring a first NC
  employee; setting up payroll; registering for NC withholding (NC-BR) or unemployment
  tax (DES/SUTA); new-hire reporting; needing workers' comp (3+ employees); FICA, FUTA,
  Form 941/940, W-4, or I-9; NC minimum wage, overtime, or final paycheck; 1099-vs-W-2
  worker misclassification; choosing a payroll provider (Gusto/QuickBooks/ADP);
  nonprofit payroll. SKIP: entity formation, the entity's own EIN, or business
  income/sales/franchise tax → venture-nc-business-formation-tax; 501(c)(3)
  recognition, charitable solicitation, Form 990 → venture-nc-nonprofit-formation;
  business planning or funding → venture-small-business-planning; consumer
  credit/debt/collections law → the consumer-credit-and-debt family.
category: personal-venture
tags: [venture, north-carolina, employer, payroll, hr-compliance]
whenToUse:
  - Hiring your first W-2 employee in North Carolina and need the full federal + state setup sequence
  - Registering for NC income-tax withholding (Form NC-BR) or NC unemployment insurance (DES / SUTA)
  - Figuring out new-hire reporting, workers' compensation, or which payroll forms you must file
  - Computing or understanding federal payroll taxes (FICA, FUTA, Form 941, Form 940) and onboarding paperwork (W-4, I-9)
  - Applying NC wage-and-hour rules (minimum wage, overtime, final paycheck, pay frequency, wage-change notice)
  - Deciding whether a worker is a W-2 employee or a 1099 contractor and avoiding misclassification penalties
  - Choosing a payroll provider (Gusto, QuickBooks Payroll, ADP) for a lean for-profit or nonprofit
  - Handling nonprofit-specific payroll questions (still owes payroll tax; 501(c)(3) UI reimbursement election)
triggers:
  - "I'm hiring my first employee in NC, what do I have to do?"
  - "How do I register for NC withholding / set up payroll taxes?"
  - "Do I need workers' comp in North Carolina?"
  - "What's the difference between a 1099 contractor and a W-2 employee?"
  - "How much is FUTA / FICA / NC unemployment tax?"
  - "When is my final paycheck due to a terminated employee in NC?"
  - "Which payroll provider should I use: Gusto, QuickBooks, or ADP?"
  - "Does my nonprofit have to pay payroll taxes?"
version: 1.0.1
updated: 2026-06-16
metadata:
  changelog:
    - "2026-06-16 sko v1.0.0->v1.0.1: fixed 1 High (Pass M: description 1961->951 chars, under 1000 Glean cap) + 1 Medium (Pass K: em-dash density 2.30->0.92/100 words, 'robust' filler removed); Pass H predicted 10/10 pos, 0/10 neg-fp; 0 banned terms; body ~5.7k tok"
---

# NC Employer Obligations & Payroll

The "first employee" layer for a North Carolina founder. Once you decide to put someone
on payroll as a **W-2 employee** (not a 1099 contractor; see classification below), you
trigger a stack of **federal + NC state** registrations, taxes, filings, and labor-law
duties. This skill is the operational checklist and reference for that stack.

**Scope boundaries (cross-reference, don't duplicate):**
- Forming the entity, getting the EIN for the *entity itself*, and business income/sales/franchise
  tax → **venture-nc-business-formation-tax**.
- 501(c)(3) recognition, charitable-solicitation license, Form 990 → **venture-nc-nonprofit-formation**.
- Business plan, financial model, funding → **venture-small-business-planning**.
- Personal credit/debt/collections law → the **consumer-credit-and-debt** family.

> **Not legal, tax, or accounting advice.** Rates, wage bases, thresholds, and deadlines
> change (often annually). Verify every figure against the primary source cited, and consult
> a CPA or employment attorney before acting. Figures below are stated with as-of dates.

---

## 1. Becoming an Employer: the registration stack

Before the first paycheck, you generally need **four** things in place:

| # | What | Who | How |
|---|------|-----|-----|
| 1 | **Federal EIN** | IRS | Free, online, instant via the `irs.gov` EIN Assistant. (Most entities already have one from formation; see venture-nc-business-formation-tax. Required to deposit payroll taxes and file 941/940/W-2.) |
| 2 | **NC withholding account** | NCDOR | **Form NC-BR** or the NCDOR online business registration → get an NC Withholding Tax ID. |
| 3 | **NC unemployment (UI/SUTA) account** | NC DES | **Form NCUI 604** (Employer Status Report) / register in **NCSUITS** online. |
| 4 | **Workers' comp coverage** | private insurer / NCIC | Required at **3+ employees** (see §4). |

Plus per-hire paperwork (**W-4, I-9**, state new-hire report; see §3, §5) and a **payroll system**
to actually compute, pay, and file (see §8).

> **Nonprofit note:** A 501(c)(3) is **still an employer** for payroll-tax purposes. It
> withholds and remits income tax and FICA exactly like a for-profit, and pays FUTA only if not
> otherwise exempt (501(c)(3)s are **exempt from FUTA**; see §6). The 501(c)(3) twist is on the
> *state UI* side: it can elect **reimbursement** instead of contributions (see §2).

**Sources:** NCDOR Business Registration, `ncdor.gov/registration`; NC-BR form page,
`ncdor.gov/taxes-forms/sales-and-use-tax/.../nc-br-business-registration-application...`;
NC DES Employers, `des.nc.gov/employers` (all as of 2026-06).

---

## 2. NC Unemployment Insurance (UI / SUTA)

State unemployment tax, administered by the **NC Division of Employment Security (DES)**.
This is an **employer-only** tax, never withheld from employee wages.

**Are you a liable employer?** (meet *any* of these; DES, as of 2026-06)
- **General business:** paid **$1,500+** in wages in any calendar **quarter**, **OR** employed
  **at least one worker in 20 different weeks** in a calendar year.
- **501(c)(3) nonprofit:** **4 or more** workers in the U.S. (≥1 in NC) in **20 different weeks**
  in this or last year.
- **Agricultural:** 10+ workers on a day in 20 different weeks, **OR** $20,000+ gross wages in a quarter.
- **Domestic/household:** $1,000+ in wages in a calendar quarter.
- Or **voluntarily elect** liability.

**Register:** Form **NCUI 604** (Employer Status Report) or via **NCSUITS** (the NC State
Unemployment Insurance Tax System, `des.nc.gov`). DES mails a PIN to manage the account online.

**Rates & wage base** (DES Tax Rate Information, as of 2026-06):
- **Standard beginning (new-employer) rate: 1.000%.** New employers pay this until they have
  enough history (generally after ~2 years) to receive an experience-based rate.
- **Taxable wage base** (per employee per year, resets annually): **$34,200 for 2026** (was
  **$32,600 for 2025**). You pay SUTA only on each employee's wages up to this cap.
- **Experienced-employer rate range:** ~**0.06% to 5.76%** (varies with the employer's reserve
  ratio and UI Trust Fund condition).

**Quarterly filing, Form NCUI 101** (Quarterly Tax and Wage Report):
- File **every quarter even if no tax is due / no employees that quarter** ("zero" report).
- Due the **last day of the month after quarter-end:** Q1 → **Apr 30**, Q2 → **Jul 31**,
  Q3 → **Oct 31**, Q4 → **Jan 31**.
- Employers with **10+ employees must file electronically**; under 10 may file on paper.

**Nonprofit (501(c)(3)) reimbursement election:** A liable 501(c)(3) may elect to be a
**reimbursing employer**. Instead of paying quarterly UI *contributions*, it reimburses DES
**dollar-for-dollar** for any unemployment benefits actually paid to its former employees. This
can save money for orgs with stable, low-turnover staffing but carries the risk of a large bill
after a layoff. [UNVERIFIED — confirm the current election mechanics, deadline, and any bonding
requirement with DES before electing.]

**Sources:** DES Tax Rate Information, `des.nc.gov/employers/tax-rate-information`; DES "Are
You Required to Pay Taxes?", `des.nc.gov/employers/are-you-required-pay-taxes`; DES Quarterly
Tax & Wage Report, `des.nc.gov/employers/file-adjust-or-review-quarterly-tax-wage-report`
(all as of 2026-06).

---

## 3. Federal New-Hire Reporting (to the NC New Hire Directory)

Federal law (PRWORA 1996) and **N.C.G.S. 110-129.2** require **every** employer (public,
private, **nonprofit**, and government) to report each **newly hired or re-hired** employee to
the **NC New Hire Directory** (run by NC DHHS Child Support Services).

- **Deadline: within 20 calendar days** of the employee's hire/start date (not business days).
- **Report:** employee name, address, SSN, hire date + employer name, address, FEIN.
- **How:** online, electronic file, mail, or fax via `ncnewhires.ncdhhs.gov`.
- **Penalty:** states may impose up to **$25 per unreported new hire**, rising to **$500** if there's
  an employer-employee conspiracy not to report.

> This is separate from, and *in addition to*, UI registration and tax filings. It's how the
> state matches employees against child-support and benefit-fraud databases.

**Sources:** NC New Hire Directory — `ncnewhires.ncdhhs.gov/law` and `/reporting_fundamentals`
(as of 2026-06).

---

## 4. Workers' Compensation (NC Industrial Commission, Ch. 97)

The **NC Workers' Compensation Act (G.S. Chapter 97)**, enforced by the **NC Industrial
Commission (NCIC)**, requires coverage once you hit the threshold.

**Threshold: 3 or more employees.** Businesses (corporation, LLC, partnership, sole prop)
**regularly employing 3+ employees** must carry workers' comp insurance **or** be approved as
self-insured.

**Counting rules** (NCIC, as of 2026-06):
- **Corporate officers** count toward the 3 even if they elect to exclude *themselves* from coverage.
- Sole proprietors, partners, and LLC members are **not automatically** counted as employees.
- A worker the business calls an "independent contractor" can be **reclassified as an employee**
  by the NCIC based on the degree of control (see §7; classification matters here too).
- An employer can be liable for an **uninsured subcontractor's** injured workers.

**Key exemptions / exceptions:** certain railroad workers; **casual** workers outside the trade or
business; **domestic servants** in a private home; **farm laborers** where fewer than 10 full-time
non-seasonal workers; federal employees in NC; certain commissioned agricultural sellers. **Any
employer using radiation** must carry coverage **regardless of headcount.**

**Penalties for going without required coverage** (as of 2026-06):
- **Civil:** the NCIC assesses **$1 per employee per day**, with a **$50/day minimum** and a
  **$100/day maximum**, for each day uninsured. [UNVERIFIED — sources gave the daily formula but
  conflicting min figures ($20 vs $50); confirm the current statutory min/max with the NCIC.]
  The penalty cannot reach back more than **3 years** before the NCIC first assesses it.
- **Criminal:** willful failure = **Class H felony**; negligent failure = **Class 1 misdemeanor**;
  responsible corporate officers can be personally liable.

**Get coverage:** through a commercial insurer / agent, the assigned-risk market (NCRB), or an
approved self-insurance program. Most lean employers buy a policy from a carrier or via their
payroll/PEO provider.

**Sources:** NCIC Employer Information — `ic.nc.gov/wcinsrqmt.html`; G.S. Chapter 97 —
`ncleg.gov/.../Chapter_97.pdf` (as of 2026-06).

---

## 5. Per-Hire Onboarding Paperwork (federal)

Collect these **before / on the first day** for every W-2 employee:

- **Form W-4 (federal):** the employee's withholding certificate that drives federal income-tax
  withholding. NC has its own equivalent, **Form NC-4**, for state withholding. Keep on file
  (don't send to IRS/NCDOR unless requested).
- **Form I-9, Employment Eligibility Verification (USCIS):** verify identity + work authorization.
  - Employee completes **Section 1 by the first day** of work.
  - Employer reviews acceptable documents and completes the employer section **within 3 business
    days** of the start date.
  - **Current form:** the **08/01/23 edition**; editions valid through **05/31/2027** (electronic
    I-9 systems must be on the 05/31/2027-expiration version by **07/31/2026**).
  - **E-Verify** (electronic confirmation against SSA/DHS) is **voluntary** federally and for most
    private NC employers, but some state contractors must use it. [UNVERIFIED — confirm any current
    NC E-Verify mandate by employer size for your situation.]
- Keep I-9s **separate** from the personnel file; retain per the USCIS retention rule (later of
  3 years after hire or 1 year after termination).

**Sources:** USCIS I-9 — `uscis.gov/i-9`; E-Verify — `e-verify.gov`; IRS Form W-4 — `irs.gov`
(as of 2026-06).

---

## 6. Federal Payroll Taxes (FICA, FUTA, withholding)

Three federal pieces ride on every paycheck. Your payroll system handles the math and deposits;
you (or it) file the returns.

**A. FICA (Social Security + Medicare)** (shared employer/employee, 2026 figures):
- **Total 7.65% each side** = **6.2% Social Security** + **1.45% Medicare**.
- **Social Security wage base 2026: $184,500** (up from $176,100 in 2025); SS tax stops above this.
  **Medicare has no wage cap.**
- **Additional Medicare Tax 0.9%** withheld from employee wages over **$200,000/yr** (employee-only;
  employer does **not** match it).

**B. Federal income-tax withholding:** based on the employee's W-4 and IRS withholding tables
(IRS **Pub 15 / 15-T**).

**C. FUTA (Federal Unemployment Tax)** (employer-only):
- **6.0% on the first $7,000** of each employee's wages, but employers in a state in good standing
  get a **5.4% credit**, making the **effective rate 0.6%** (= **$42/employee/yr** max). NC is
  normally a full-credit state.
- **501(c)(3) nonprofits are exempt from FUTA.**

**Federal returns:**
- **Form 941** (quarterly): reports withheld income tax + both halves of FICA. Due the last day of
  the month after each quarter (Apr 30 / Jul 31 / Oct 31 / Jan 31). (Very small employers may be
  assigned **annual Form 944** instead.)
- **Form 940**: annual FUTA return, due **Jan 31**.
- **Deposits** (the actual money) are made on a **monthly or semiweekly** schedule via **EFTPS**,
  determined by your lookback-period liability, *separate from* the return due dates.
- **W-2** to each employee and **W-3** transmittal to SSA by **Jan 31**.

**Sources:** IRS Pub 15 (Circular E) and Pub 15-T — `irs.gov`; 2026 Social Security wage base
(SSA announcement, as of 2026-06). FUTA rate per IRS Form 940 instructions.

---

## 7. NCDOR State Income-Tax Withholding

Withhold **NC income tax** from employee wages and remit to NCDOR.

- **Register:** Form **NC-BR** / NCDOR online business registration → NC Withholding Tax ID. NCDOR
  assigns your **filing frequency** in the confirmation letter.
- **Rate:** NC has a **flat individual income-tax rate of 3.99% for 2026** (down from 4.25% in 2025,
  per S.L. 2023-134, effective for wages paid on/after 2026-01-01). Withholding tables (**Form NC-30**)
  build in the standard deduction, so the **percentage-method withholding rate is ~4.09% for 2026**.
- **Filing frequency & forms** (NCDOR, as of 2026-06):
  - Withhold **< $250/month average** → **quarterly**, **Form NC-5**.
  - Withhold **$250–$1,999/month average** → **monthly**, **Form NC-5**.
  - Withhold **$2,000+/month average** → **semiweekly**, deposit with **Form NC-5P**; reconcile
    quarterly on **Form NC-5Q**.
- **Annual reconciliation: Form NC-3**, filed **with W-2s and 1099s by Jan 31**, electronically.

**Sources:** NCDOR Withholding Tax — `ncdor.gov/taxes-forms/withholding-tax` and Withholding Tax
FAQ; 2026 Income Tax Withholding Tables and Instructions for Employers (NC-30); NCDOR Individual
Income Tax Rate Schedules (all as of 2026-06).

---

## 8. Wage & Hour: NC Wage and Hour Act + FLSA

Two layers apply: the **federal Fair Labor Standards Act (FLSA)** and the **NC Wage and Hour Act
(NCWHA)**, enforced by the **NC Department of Labor (NCDOL)**. When they differ, the rule **more
protective of the employee** generally governs.

- **Minimum wage:** NC's minimum wage **tracks the federal $7.25/hour** (as of 2026-06). NC has no
  higher state minimum and **no local minimum-wage** preemption carve-out: $7.25 statewide.
- **Overtime:** **1.5× the regular rate** for hours over **40 in a workweek** for **non-exempt**
  employees (FLSA + NCWHA). NC has **no daily overtime** rule.
- **Exempt vs non-exempt:** "exempt" (no overtime) requires meeting an FLSA **duties test** (executive,
  administrative, professional, outside sales, certain computer roles) **and** a **salary threshold**.
  Job title alone never makes someone exempt. [UNVERIFIED — the federal exempt salary threshold has
  been in flux through 2024–2026 litigation; confirm the current DOL threshold before classifying
  anyone as exempt.]
- **Pay frequency:** employers must set **regular paydays** and pay **at least monthly** (wages may be
  paid daily/weekly/bi-weekly/semi-monthly/monthly). Salaried employees may be paid monthly.
- **Final paycheck:** a terminated or departing employee must be paid all wages due **on or before the
  next regular payday** (by trackable mail if the employee requests). NC does **not** require immediate
  same-day final pay.
- **Wage-change / pay notice:** notify employees **in writing (or by posted notice) at least 24 hours**
  before **reducing** wages or changing wage benefits (a wage *increase* needs no advance notice).
  Provide written notice of pay rate, payday, and benefits at hire.
- **Required posters:** display the NCDOL Wage and Hour Act notice and federal labor-law posters.

**Sources:** NCDOL Wages & Promised Wages — `labor.nc.gov/workplace-rights/...`; NCDOL Wage and Hour
Act notice; FLSA via US DOL — `dol.gov` (all as of 2026-06).

---

## 9. W-2 Employee vs 1099 Contractor (classification)

The single highest-risk payroll decision. **You don't get to choose by preference or by contract;**
the working relationship dictates the status. Misclassification is a **top IRS *and* DOL enforcement
priority** and exposes you to **back taxes, penalties, back overtime/minimum-wage, and the NCIC's
power to reclassify** for workers'-comp purposes.

**IRS common-law test ("right to control"),** weighed across **3 categories** (no single factor decides):
1. **Behavioral control:** do you direct *how, when, where* the work is done? Training? Detailed
   instructions? (→ employee)
2. **Financial control:** who provides tools/equipment, who bears profit/loss risk, is the worker
   free to seek other clients, how is pay structured (salary/hourly vs by-the-job)? (→ contractor if
   the worker bears business risk)
3. **Relationship:** written contract, benefits, permanency, is the work a **core part of your
   business**? (→ employee if integral and open-ended)

**Strong "this is really an employee" red flags:** you set their hours, require on-site work, supply
their equipment, control task methods, bar them from other clients, and there's no defined end date.

**If genuinely unsure:** either party can file **IRS Form SS-8** for an official determination (can
take 6+ months). NC agencies (DES for UI, NCIC for comp, NCDOL for wages) apply their **own
control-based tests** and can each reach a different result, so a worker can be a contractor for one
program and an employee for another.

**Practical rule for a lean founder:** if the person does ongoing, directed, core work under your
control, treat them as a **W-2 employee**. The cost of getting it wrong dwarfs the convenience of a 1099.
**1099-NEC** is for genuinely independent vendors (their own business, own tools, own risk, multiple clients).

**Sources:** IRS Independent Contractor (Self-Employed) or Employee + Form SS-8 — `irs.gov`; US DOL
Misclassification — `dol.gov` (as of 2026-06).

---

## 10. Payroll Providers for a Lean Org

You *can* run payroll by hand (compute withholding, deposit via EFTPS, file 941/940/NC-5/NCUI 101,
issue W-2s), but for a first hire a provider that **files and deposits taxes for you** is almost
always worth it and reduces misclassification/late-filing risk.

| Provider | Best for | Notes (as of 2026-06) |
|----------|----------|------------------------|
| **Gusto** | First-timers wanting payroll **+ HR/benefits/onboarding** in one clean UI | ~**$49/mo base + $6/employee/mo** (Simple plan; base rose from $40 in Mar 2026). Handles federal + NC filings, new-hire reporting, W-2/1099. Popular with very small orgs and nonprofits. |
| **QuickBooks Payroll** | Orgs **already on QuickBooks Online** | Core ~**$50/mo + $6/employee**; Premium/Elite tiers higher. Native QBO accounting integration is the main draw. |
| **ADP (RUN)** | Multi-state, faster-growing, or HR-heavy orgs | Deep compliance + benefits marketplace; **opaque pricing**, often a **setup fee** and contract. More than most first-hire orgs need. |

**What any provider should do for you:** withhold and **deposit** federal + NC taxes; file **941/940,
NC-5/NC-5Q, NCUI 101**; issue **W-2/W-3 and 1099s**; submit **new-hire reports**; and often broker
**workers'-comp** (pay-as-you-go) and benefits.

> **Nonprofit note:** providers handle nonprofit payroll fine, but **tell them you're a 501(c)(3)**
> so they suppress **FUTA** (you're exempt) and correctly handle a **state UI reimbursement election**
> if you make one. These are common setup mistakes.

**Sources:** vendor comparison pages and pricing (Gusto, Intuit QuickBooks, ADP), as of 2026-06.
Confirm current pricing/features directly with each vendor.

---

## First-Hire Checklist (NC, for-profit or nonprofit)

**Before posting the role**
- [ ] Decide **W-2 vs 1099** honestly using the IRS control test (§9). If it's ongoing directed core
      work → W-2.
- [ ] Confirm you'll be at **3+ employees** → line up **workers' comp** (§4).

**Set up accounts (one-time)**
- [ ] **EIN** in place (IRS); usually already done at formation.
- [ ] **NCDOR withholding** account via Form **NC-BR** / online (§7).
- [ ] **NC DES UI** account via Form **NCUI 604** / NCSUITS (§2).
- [ ] **Workers' comp** policy bound if 3+ employees (§4).
- [ ] Choose a **payroll provider** (§10).
- [ ] *Nonprofit:* tell payroll provider you're a **501(c)(3)** (FUTA-exempt; consider UI
      **reimbursement** election).

**Each new hire**
- [ ] Collect **Form I-9** (Section 1 day one; employer part within 3 business days) (§5).
- [ ] Collect **Form W-4** (federal) and **Form NC-4** (state) (§5, §7).
- [ ] **Report to the NC New Hire Directory within 20 days** (§3).
- [ ] Provide **written pay/payday/benefits notice**; post NCDOL + federal posters (§8).

**Ongoing**
- [ ] Run payroll on a **regular payday**, withhold correctly, **deposit** taxes on schedule (§6, §7).
- [ ] File **Form 941** (quarterly), **NC-5/NC-5Q** (per frequency), **NCUI 101** (quarterly, even if zero).
- [ ] File **Form 940** (FUTA, annual); *skip if 501(c)(3)*.
- [ ] Year-end: **W-2/W-3** + **NC-3** (with W-2s/1099s) by **Jan 31**.
- [ ] Pay **final paychecks** by the next regular payday; give **24-hr written notice** before any wage cut.

---

## Sources

All accessed 2026-06. Verify figures against the primary source before relying on them.

- **NCDOR** — Business Registration `ncdor.gov/registration`; Form NC-BR page; Withholding Tax
  `ncdor.gov/taxes-forms/withholding-tax` + Withholding Tax FAQ; 2026 Income Tax Withholding Tables
  & Instructions (Form NC-30); Individual Income Tax Rate Schedules.
- **NC DES (Division of Employment Security)** — Tax Rate Information
  `des.nc.gov/employers/tax-rate-information`; "Are You Required to Pay Taxes?"
  `des.nc.gov/employers/are-you-required-pay-taxes`; Quarterly Tax & Wage Report (NCUI 101)
  `des.nc.gov/employers/file-adjust-or-review-quarterly-tax-wage-report`; NCSUITS.
- **NC New Hire Directory (NC DHHS)** — `ncnewhires.ncdhhs.gov/law`, `/reporting_fundamentals`;
  N.C.G.S. 110-129.2.
- **NC Industrial Commission** — Employer Information `ic.nc.gov/wcinsrqmt.html`; G.S. Chapter 97
  `ncleg.gov`.
- **NCDOL** — Wage and Hour Act / Promised Wages `labor.nc.gov`; NCDOL Wage and Hour notice.
- **IRS** — Pub 15 (Circular E) & 15-T; Forms 941, 940, 944, W-2/W-3, W-4; Independent Contractor
  vs Employee + Form SS-8; EIN Assistant — all `irs.gov`.
- **US DOL** — FLSA and worker-misclassification guidance — `dol.gov`.
- **USCIS / E-Verify** — Form I-9 `uscis.gov/i-9`; `e-verify.gov`.
- **SSA** — 2026 Social Security wage base announcement.
- **Payroll providers** — Gusto, Intuit QuickBooks Payroll, ADP comparison/pricing pages.

## Disclaimer

This is general educational information, **not** legal, tax, accounting, or HR advice, and does not
create any professional relationship. Employment and tax rules (rates, wage bases, thresholds,
filing frequencies, exempt-salary levels, and deadlines) change frequently and depend on your
specific facts. Before hiring or running payroll, verify current requirements directly with **NC DES,
NCDOR, the NC Industrial Commission, NCDOL, the IRS, USCIS/E-Verify, and the US DOL**, and consult a
**CPA and/or a North Carolina employment attorney**. Items marked **[UNVERIFIED]** were not confirmed
against a primary source during research and must be checked before use.
