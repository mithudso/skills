<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-credit-and-debt` hub.** Formerly the standalone `bankruptcy-ch7-ch13` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-credit-and-debt`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: bankruptcy-ch7-ch13
description: >-
  US consumer bankruptcy: Chapter 7 (liquidation) vs Chapter 13
  (reorganization). Spoke of consumer-credit-and-debt; educational info, NOT
  legal advice (as of 2026). TRIGGER: should I file bankruptcy; Ch.7 vs Ch.13 /
  which chapter; means test / qualify for Ch.7; automatic stay stops
  garnishment/foreclosure/repossession/calls; what's discharged vs
  non-dischargeable (taxes, student loans, support, fraud); 341 meeting;
  trustee; exemptions; credit-counseling/debtor-ed; alternatives; student-loan
  discharge / undue hardship. SKIP: NC exemption/garnishment ->
  north-carolina-credit-and-debt-law; rebuilding a score after discharge ->
  improving-and-rebuilding-credit; how a bankruptcy ages/scores on a report ->
  credit-reports-and-scores; settling a debt WITHOUT filing (1099-C) ->
  charge-offs-collections-and-debt-resolution; stopping ONE collector /
  validation -> debt-collectors-and-fdcpa-rights; estate/probate & wills ->
  estate-planning-and-wills; statute text -> us-consumer-credit-and-debt-law.
metadata:
  version: 1.1.0
  updated: 2026-06-16
  changelog:
    - "2026-06-16 sko v1.0.0->v1.1.0 (--meta --no-sync) — 1 Medium structural: Pass I collision vs newly-expanded consumer-finance sibling set; +SKIP edge estate/probate & wills -> estate-planning-and-wills (reciprocal; that skill owns estate/probate, this skill owns the bankruptcy filing) in desc + body; desc trimmed to 997 (<=1000). No content passes."
---

# US Consumer Bankruptcy: Chapter 7 vs Chapter 13

> **This is general educational information, NOT legal, financial, or tax advice.** Bankruptcy is a federal court process that is **complex and intensely fact-specific**: outcomes turn on your income, your state's exemption law, the nature of each debt, and recent filings. **Means-test dollar figures and many rules change** (the income tables update twice a year; expense and debt-limit figures adjust periodically). Everything below is **as of 2026** and may be out of date. **Strongly consider consulting a licensed bankruptcy attorney** before filing or deciding not to; most offer a free initial consultation. Verify any current figure against the primary sources in **References** at the end.

This skill is a spoke of the **consumer-credit-and-debt** hub.

**How to answer with this skill (output guidance).** Stay **general and educational**. Do **not** render a verdict on whether a specific person qualifies for Chapter 7 or tell them to file (or not file) a particular chapter; instead, explain how the rules work and what factors matter. Lead with the not-legal-advice framing, **surface the licensed-attorney referral**, flag that means-test figures and debt limits **change** (point to References to verify current numbers), and route adjacent questions per the SKIP edges.

---

## What bankruptcy actually does

Bankruptcy is a legal process under federal law (Title 11, the **Bankruptcy Code**) that gives an honest but overwhelmed debtor a "fresh start." Two mechanisms do the work:

- **The automatic stay (11 U.S.C. § 362).** The *instant* you file the petition, an automatic stay takes effect, with no court order or hearing needed. It **immediately stops** most collection activity: lawsuits, **wage garnishment**, bank levies, **foreclosure** sales, **vehicle repossession**, utility shut-offs, and creditor phone calls/letters. It is one of the most powerful and immediate reasons people file. Limits: it does not stop most **criminal** proceedings or most **domestic-support (child/spousal support)** collection, a creditor can ask the court to "lift" the stay (e.g., to proceed with foreclosure on a home you are surrendering), and **repeat filers** may get a shortened stay or none.
- **The discharge.** A court order that **permanently wipes out personal liability** for the debts that qualify; creditors can never again try to collect a discharged debt. The discharge is the end goal. **Not every debt is dischargeable** (see the dischargeability list), and a **valid lien survives** discharge on property you keep: discharge erases the *personal obligation to pay*, not necessarily a creditor's right to the *collateral* (e.g., you can't keep a financed car and stop paying for it).

---

## Chapter 7 vs Chapter 13 at a glance

| **Dimension** | **Chapter 7: "liquidation" / straight bankruptcy** | **Chapter 13: "reorganization" / wage-earner's plan** |
|---|---|---|
| **Core idea** | Eligible debts wiped out quickly; non-exempt property can be sold to pay creditors. | A **3–5 year court-supervised repayment plan**; keep your property, pay back some/all debt from future income, discharge the rest at the end. |
| **Eligibility gate** | Must pass the **means test** if income is above the state median. | Must have **regular income** and debts **under the statutory debt limits** (set by statute and adjusted periodically — verify the current figures); cannot have had a recent disqualifying discharge. |
| **Typical timeline** | **~3–6 months** from filing to discharge. | Discharge comes only after **3–5 years** of completed plan payments. |
| **What happens to property** | A trustee may **liquidate (sell) non-exempt** property; **exempt** property is protected and most filers keep everything ("no-asset" cases are common). | You **keep your property**; instead you commit disposable income to the plan. |
| **Best when** | Income is low/moderate, mostly **unsecured** debt (cards, medical), little non-exempt property, no need to cure secured arrears. | You're **behind on a mortgage or car** and want to keep it (cure arrears over time); you're **over the means test**; you have non-exempt assets you want to protect; or you need lien stripping / cram-down. |
| **Discharge timing** | At case close (~60–90 days after the 341 meeting). | At **completion of all plan payments**. |
| **Stays on credit report** | **~10 years** from the filing date (FCRA cap). | Commonly **~7 years** (a credit-bureau practice for completed plans, not an FCRA rule — the 10-year cap technically covers all bankruptcies). |
| **Both require** | Pre-filing **credit counseling** + pre-discharge **debtor-education** course; a **petition + schedules**; a **341 meeting of creditors**; a **trustee**. | Same. |

---

## Chapter 7 (liquidation) in depth

- **The means test (11 U.S.C. § 707(b)).** Added by the 2005 **BAPCPA** reform to keep higher-income filers out of Chapter 7. Step 1: compare your **current monthly income** (your average gross monthly income over the **6 full calendar months before filing**, annualized) to the **median family income for your state and household size**. If you're **at or below** the median, you pass the presumption calculation (§ 707(b)(2)); a case can still rarely be challenged for bad faith under the "totality of circumstances" (§ 707(b)(3)). If you're **above**, Step 2 subtracts allowed living expenses (using **IRS National & Local Standards**) to compute disposable income; if what's left over five years clears a statutory dollar threshold, a presumption of "abuse" arises and you may be steered to Chapter 13 (or dismissed). *The abuse-presumption dollar thresholds adjust on a 3-year cycle, while the underlying median-income and IRS-standard inputs update more often (roughly twice a year) — so look up all current figures with the US Trustee Program rather than relying on any number quoted here.* Special circumstances can rebut the presumption.
  - **Married filers:** spouses may file **jointly or individually**; **household size and combined income** drive the median comparison, and a non-filing spouse's income can still count toward it (the "marital adjustment").
- **Exemptions — what you keep.** You "exempt" property up to legal limits and the trustee can only liquidate **non-exempt** assets; whether you use the federal or your state's exemptions depends on your state (see the dedicated *Exemptions* section below — North Carolina, for one, requires the **state** set under **G.S. § 1C-1601**).
- **The trustee.** An impartial **case trustee** is appointed to review your petition, run the 341 meeting, and **sell any non-exempt property**, distributing proceeds to creditors by statutory priority. In most consumer cases there's nothing non-exempt to sell ("no-asset" cases).
- **Timeline.** Roughly **3–6 months**: file → automatic stay → 341 meeting (~21–40 days after filing) → discharge (~60–90 days after the 341 meeting if no objections).
- **Secured debt choices.** For a financed house or car you can generally **reaffirm** (keep paying), **redeem** (pay the item's current value in a lump sum), or **surrender** it. You can't keep collateral while discharging the debt on it. **Caution:** reaffirming makes you **personally liable again**, so a later default can leave you owing a deficiency; some debtors instead keep paying without reaffirming (an informal "ride-through"), but whether that survives is **jurisdiction- and lender-dependent** and post-BAPCPA the stay can end as to the collateral if you don't perform your stated intention — get advice before relying on it.

---

## Chapter 13 (reorganization) in depth

- **What it is.** A repayment plan for an **individual with regular income** (11 U.S.C. ch. 13). You propose a plan to pay creditors out of future income over **3 to 5 years**; the trustee collects your single monthly payment and disburses it.
- **Debt limits.** Eligibility caps your **secured and unsecured debt** (historically **two separate limits**, not one combined figure; a temporary single combined limit enacted in June 2022 by the Bankruptcy Threshold Adjustment and Technical Corrections Act **sunset in June 2024**, reverting to the two-tier structure). The figures are statutory and adjust periodically, so **verify the current limits** — see References.
- **Plan length (§ 1322(d)).** If your income is **below** the state median, the plan is normally **3 years**; if **above** the median, it must run **5 years**. **No plan may exceed 5 years.**
- **Who it suits / why choose it over Chapter 7:**
  - **Cure mortgage arrears to keep a home.** Chapter 13 **stops foreclosure** and lets you pay off the past-due balance over the plan while keeping current; Chapter 7 can't do this.
  - **You're over the means test.** Chapter 13 is the path when income is too high for Chapter 7.
  - **Protect non-exempt assets** you'd lose in a Chapter 7 liquidation.
  - **Lien stripping / cram-down.** You may "strip" a wholly unsecured junior mortgage (a second/HELOC with no equity behind it) and **cram down** certain other secured debts to the collateral's value (e.g., an older car loan), but **not** a mortgage on your principal residence, and car cram-down requires the loan to be old enough (a "910-day" rule for purchase-money car loans).
- **Co-debtor stay (§ 1301).** Unique to Chapter 13: creditors generally **cannot pursue a co-signer** on a consumer debt while the plan is active, a protection Chapter 7 does not give cosigners.
- **Discharge (§ 1328).** Granted **only after you complete all plan payments** (and certify support obligations are current and finish the debtor-education course). The Chapter 13 discharge is **somewhat broader** than Chapter 7's: it can reach a few debts that Chapter 7 can't (e.g., some willful-property-damage debts, debts incurred to pay non-dischargeable taxes, certain divorce **property-settlement** obligations). If you can't finish the plan, a **hardship discharge** is sometimes available.

---

## What bankruptcy will NOT erase (non-dischargeable debts)

Some debts survive a discharge in **both** chapters (Chapter 13 reaches a few more than Chapter 7). Commonly **non-dischargeable**:

- **Most student loans:** discharged only on a showing of **"undue hardship"** in a separate adversary proceeding (most courts apply the **Brunner** test). Historically very hard; **see the 2022 guidance note below.**
- **Recent income taxes:** income tax is dischargeable only if it's **old enough** (generally the return was **due more than ~3 years** ago), the **return was actually filed** (and not late beyond limits or fraudulent), and the tax was **assessed** in time. **Payroll/trust-fund taxes and most penalties are never dischargeable**, and there is **no discharge of post-petition tax**. You must **keep filing returns and paying current taxes** during the case. *(IRS guidance.)* → For tax mechanics generally, confirm with the IRS source in References.
- **Domestic support obligations:** **child support and alimony/spousal support** are non-dischargeable.
- **Debts from fraud, false financial statements, embezzlement, or larceny**, and recent luxury purchases / cash advances run up just before filing.
- **Debts for willful and malicious injury** to a person or property.
- **Death or personal injury from drunk/impaired driving (DUI/DWI).**
- **Most criminal fines, penalties, and restitution.**
- Certain other items: **post-petition** condo/HOA dues (pre-petition dues are dischargeable), debts left off the schedules, and debts from a prior bankruptcy where discharge was denied.

> Recent change, **student loans (date-stamped, November 2022):** the **DOJ** (Civil Division / US Trustee Program), in coordination with the **Department of Education**, issued **guidance (Nov 17, 2022)** giving government attorneys a clearer, fairer process for federal **Title IV** student-loan undue-hardship discharge. Borrowers complete a standardized **attestation form**; the government applies consistent factors and will **recommend (or not oppose) discharge** when the facts support it. DOJ/ED have reported that a large majority of borrowers using the process obtained **full or partial discharge**. This makes student-loan discharge **more attainable than the old reputation suggests**, but it is still case-by-case and **as of 2026** may have evolved; confirm current policy and consult an attorney.

---

## The process (both chapters)

1. **Pre-filing credit counseling (required).** Complete an **approved credit-counseling course within the 180 days before filing** (the eligibility rule is § 109(h); § 111 governs provider approval) from an agency on the **US Trustee Program's approved list**. (Narrow exceptions/waivers exist.)
2. **Petition & schedules.** File the bankruptcy petition plus detailed **schedules**: all assets, debts, income, expenses, contracts, and recent financial transactions, under penalty of perjury. **Accuracy matters**; errors or omissions can cost the discharge. Filing triggers the **automatic stay**. (There is a **court filing fee**; a Chapter 7 filer who can't afford it may apply for a **fee waiver** or installments under 28 U.S.C. § 1930(f).)
3. **Trustee assignment.** A trustee is appointed (Chapter 7: may liquidate non-exempt assets; Chapter 13: administers the plan and distributes payments).
4. **341 meeting of creditors (11 U.S.C. § 343).** A mandatory meeting (typically **21–40 days after filing**, often by phone/video) where the **trustee** (not a judge; judges are barred from attending) puts you under oath and asks about your debts, assets, and paperwork. Creditors may attend but usually don't. Usually brief.
5. **Pre-discharge debtor education (required).** Complete an **approved personal-financial-management ("debtor education") course** *after* filing and *before* discharge (a different course from the pre-filing counseling).
6. **Chapter 13 only:** the court holds a **plan-confirmation** hearing (§ 1325); then you make plan payments for 3–5 years.
7. **Discharge.** Chapter 7: ~60–90 days after the 341 meeting. Chapter 13: after the final plan payment.

**Refiling time-bars (the "recent disqualifying discharge").** You can't get a *second* discharge too soon after a prior one: roughly **8 years** between **Chapter 7 → Chapter 7** discharges, and shorter cross-chapter bars (commonly **4 years** Ch.7 → Ch.13, and **2 years** Ch.13 → Ch.13), measured filing-date to filing-date (§§ 727(a)(8), 1328(f)). You can usually still *file* — you just may not receive a discharge.

**Attorney vs. pro se.** You *can* file **"pro se"** (without a lawyer), and it's most feasible in a simple no-asset Chapter 7. But bankruptcy is technical: exemption planning, the means test, secured-debt strategy, and Chapter 13 plans are error-prone, and **mistakes can mean losing property or the discharge.** Chapter 13 in particular is very difficult to do well pro se. **A licensed bankruptcy attorney is strongly recommended.** ("Bankruptcy petition preparers" can type forms but cannot give legal advice.)

---

## Exemptions — federal vs state (quick orientation)

- The Bankruptcy Code provides a **federal exemption set** (§ 522(d)), but **each state decides** whether its residents use the federal set or are **limited to the state's own exemptions**.
- Roughly **half the states opt out** and **require state exemptions**; the rest let you choose federal or state. Domicile/residency timing rules determine which state's exemptions apply.
- **North Carolina requires the state exemptions** — NC filers cannot use the federal set; NC's are in **G.S. § 1C-1601** (homestead, motor vehicle, "wildcard," etc.).
- → For NC exemption **amounts and homestead/garnishment specifics**, use **north-carolina-credit-and-debt-law**. This skill stays at the federal/conceptual level.

---

## Credit impact and rebuilding

- A **Chapter 7** filing generally remains on a credit report for about **10 years** from the filing date — the FCRA's outer limit for any bankruptcy. A completed **Chapter 13** is commonly removed at about **7 years**, but that is a **credit-bureau practice, not an FCRA entitlement** (the 10-year cap technically applies to all bankruptcies). Individual discharged accounts also report as "included in bankruptcy."
- The score impact is large at first but **fades over time**, and many filers (freed of unmanageable debt) start rebuilding immediately. Post-bankruptcy credit (secured cards, credit-builder loans, on-time payments) can move a score up well before the public record ages off.
- → For the **rebuilding playbook after discharge**, use **improving-and-rebuilding-credit**. For **how the bankruptcy notation ages and affects a score**, use **credit-reports-and-scores**.

---

## Alternatives to bankruptcy (consider first / alongside)

Bankruptcy is powerful but serious; weigh these:

- **Nonprofit debt-management plan (DMP).** Through a reputable **nonprofit credit-counseling agency** (e.g., NFCC members): one consolidated monthly payment, often reduced interest, paid over ~3–5 years. No new loan, no bankruptcy record. → see **improving-and-rebuilding-credit** / the hub.
- **Debt settlement.** Negotiating to pay less than the full balance (often via lump sum). Can hurt credit, may trigger a **1099-C** (canceled debt can be taxable income), and for-profit settlement firms carry real risks. → see **charge-offs-collections-and-debt-resolution**.
- **Doing nothing if "judgment-proof."** If your only income/assets are protected (e.g., Social Security, and **wages that can't be garnished** in your state) and you have nothing non-exempt, creditors may have **no practical way to collect** even with a judgment. This is **not** debt forgiveness (the debt and collection attempts continue, and circumstances can change), but for some, **not filing** is rational. → state-specific garnishment/exemption analysis: **north-carolina-credit-and-debt-law**; stopping a specific collector: **debt-collectors-and-fdcpa-rights**.
- **Negotiating directly** with creditors (hardship programs, forbearance), or **workout** of a single secured debt (loan modification on a mortgage).

When unmanageable secured arrears, lawsuits/garnishment, or non-dischargeable-but-payable debt dominate, bankruptcy may still be the best tool, and **an attorney can compare these against Chapter 7 and 13 for your facts.**

---

## Cross-references (consumer-credit-and-debt hub)

- **north-carolina-credit-and-debt-law** — NC exemptions (G.S. § 1C-1601), who can be garnished in NC, NC foreclosure, NC statute of limitations.
- **improving-and-rebuilding-credit** — rebuilding a score after discharge; secured cards, credit-builder loans.
- **credit-reports-and-scores** — how a bankruptcy public record is reported, ages, and affects FICO/VantageScore.
- **charge-offs-collections-and-debt-resolution** — settling debts *without* filing; 1099-C canceled-debt tax; pay-for-delete.
- **debt-collectors-and-fdcpa-rights** — stopping/handling a specific collector short of bankruptcy; debt validation; being sued.
- **estate-planning-and-wills** (`consumer-finance` sibling family) — wills, powers of attorney, beneficiary designations, and **NC probate/estate administration**; this skill owns the *bankruptcy filing* itself, that skill owns the *estate/probate* process (a decedent's debts in NC → **north-carolina-credit-and-debt-law**).
- **us-consumer-credit-and-debt-law** — literal federal statute text (FCRA/FDCPA/ECOA) and enforcement.

---

## References / verify current (as of 2026 — figures and rules change)

Primary, authoritative sources — re-check these for current dollar amounts, debt limits, and policy:

- **US Courts — Bankruptcy Basics (the authoritative consumer overview):**
  - Chapter 7: https://www.uscourts.gov/court-programs/bankruptcy/bankruptcy-basics/chapter-7-bankruptcy-basics
  - Chapter 13: https://www.uscourts.gov/court-programs/bankruptcy/bankruptcy-basics/chapter-13-bankruptcy-basics
  - Discharge in Bankruptcy: https://www.uscourts.gov/court-programs/bankruptcy/bankruptcy-basics/discharge-bankruptcy-bankruptcy-basics
  - Process: https://www.uscourts.gov/court-programs/bankruptcy/bankruptcy-basics/process-bankruptcy-basics
- **DOJ — US Trustee Program (means test, approved courses, student-loan guidance):**
  - Means Testing (Census median-income tables + IRS National/Local Standards; income data effective **2026-04-01**, expense standards effective **2025-05-15**): https://www.justice.gov/ust/means-testing
  - Credit Counseling & Debtor Education approved providers: https://www.justice.gov/ust/credit-counseling-debtor-education-information
  - Student Loan Guidance (Nov 17, 2022 process + attestation form): https://www.justice.gov/ust/student-loan-guidance
- **IRS — bankruptcy & federal tax debt** (dischargeability conditions, keep filing/paying): https://www.irs.gov/businesses/small-businesses-self-employed/declaring-bankruptcy
- **CFPB — consumer bankruptcy basics & "Ask CFPB":** https://www.consumerfinance.gov/ (search "bankruptcy")
- **The statute itself — Bankruptcy Code, Title 11 U.S.C.** (§ 362 automatic stay; § 522 exemptions; § 523 nondischargeable debts; § 707(b) means test; §§ 1301/1322/1325/1328 Chapter 13): https://uscode.house.gov/
- **North Carolina exemptions** — N.C. Gen. Stat. **§ 1C-1601** (see north-carolina-credit-and-debt-law).

*This document is informational only and is not a substitute for advice from a licensed attorney about your specific situation.*
