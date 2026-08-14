<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-credit-and-debt` hub.** Formerly the standalone `predatory-lending-and-high-cost-credit` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-credit-and-debt`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: predatory-lending-and-high-cost-credit
description: >-
  Recognizing predatory, high-cost consumer credit and safe alternatives. Spoke
  of consumer-credit-and-debt; educational, NOT financial/legal advice; rules
  vary by state, change (as of 2026).
  TRIGGER: is a payday/auto-title/rent-to-own/pawn/"no credit check"/high-cost
  installment loan a bad idea; ~400% APR / rollover / debt-trap cycle;
  cash-advance or earned-wage-access (EWA) app; overdraft as credit; predatory
  signs (equity stripping, flipping, packing, balloon, prepay penalty); state
  36% rate-cap movement; Military Lending Act 36% MAPR; safer alternatives
  (PAL, CDFI, employer, counseling).
  SKIP: NC payday/usury -> north-carolina-credit-and-debt-law; TILA/HOEPA/MLA
  text -> us-consumer-credit-and-debt-law; mortgage detail ->
  home-mortgage-lending; PAL fees/terms/enrollment & building credit ->
  improving-and-rebuilding-credit; BNPL scoring -> credit-reports-and-scores;
  deposit-account & NSF fee mechanics -> personal-banking;
  defaulted-loan collector -> debt-collectors-and-fdcpa-rights.
metadata:
  version: 1.1.0
  updated: 2026-06-16
  changelog:
    - "2026-06-16 sko v1.0.0->v1.1.0 (--meta --no-sync) — 2 Medium structural: Pass I overdraft-seam asymmetry vs newly-expanded consumer-finance sibling set (+SKIP edge deposit-account & NSF fee mechanics -> personal-banking in desc + body; this skill owns overdraft-as-credit/APR framing, personal-banking owns the account mechanics); resolved the resulting mutual-SKIP [overdraft] loop into a one-way specificity gradient by dropping the shared token from the edge. desc 997 (<=1000). No content passes."
---

# Predatory Lending & High-Cost Credit

General educational information — **not financial or legal advice.** Rules, rate
caps, and product availability vary by state and change frequently. Figures are
**as of 2026**; verify current law and rates against the sources at the end and
your own state regulator before acting. This is a spoke of the
**consumer-credit-and-debt** hub.

The common thread across every product below: the price is structured to be
opaque (fees, not a stated APR), repayment is engineered so you re-borrow, and
the lender is protected against your default (your paycheck, your car, your
goods, or fine print). The defense is the same each time — **convert the cost to
an annual percentage rate (APR), ask what happens if you can't pay, and compare
against a safe alternative before signing.**

---

## How to spot a predatory or high-cost loan — checklist

Treat **two or more** of these as a red flag; several together is a debt trap:

- **Triple-digit APR**, or a price quoted only as a flat fee ("$15 per $100")
  that hides the APR. A $15-per-$100 two-week fee ≈ **~400% APR** (CFPB).
- **"No credit check," "guaranteed approval," "bad credit OK."** Legitimate
  lenders price risk; these phrases signal the lender is relying on your
  paycheck or collateral, not your ability to repay.
- **Repayment due in a single balloon** on your next payday (or one lump sum),
  rather than affordable installments — designed to force a rollover.
- **The lender takes a security interest in something you can't lose**: a
  post-dated check or ACH access to your bank account, your **car title**, your
  household goods, or your future wages.
- **Underwriting ignores ability to repay**: they don't check whether the
  payment fits your budget, only that they can collect.
- **Loan flipping / repeated refinancing**, rollover, or "renewal" offers that
  reset fees instead of reducing principal.
- **Packing**: single-premium credit insurance, "auto club," roadside, or
  GAP-type **add-ons** bundled into the loan and financed at the loan's rate.
- **Prepayment penalties** (you're punished for paying early/fast).
- **Mandatory arbitration / class-action waiver** buried in the contract.
- **Equity stripping** (mortgage context): a loan sized to your home equity, not
  your income — the lender profits from foreclosure, not repayment.
- **Pressure, rushed signing, blank fields, or "sign now or lose the rate."**

> Quick conversions: monthly fee ÷ amount × 12 ≈ APR. A "$30 per $100 per two
> weeks" fee ≈ ~780% APR. Anything you must roll over to afford is a trap.

---

## Per-product summaries

### Payday loans (the canonical debt trap)
- **Mechanics:** small, short-term, very high-cost loan (typically **$500 or
  less**), due in full on your next payday (usually ~2 weeks). You write a
  post-dated check or authorize an ACH debit for principal + fee.
- **Cost:** fees commonly **$10–$30 per $100** borrowed; a $15/$100 two-week loan
  ≈ **almost 400% APR** (CFPB).
- **The trap:** the loan is due as a lump sum you usually can't spare, so you
  "roll over"/renew — paying only the fee to extend. CFPB: **4 of 5 payday loans
  are rolled over or renewed within 14 days**, and most are made to borrowers who
  renew so often they pay **more in fees than they originally borrowed**. Example
  (CFPB): a $300 loan at $45/2 weeks costs $360 in fees over four months and you
  **still owe the original $300**.
- **State patchwork:** legality and caps vary widely. **North Carolina
  effectively bans payday lending** (its rate caps make the product unlawful;
  in-state APRs of 390–780% prompted the ban) — for NC specifics see
  **north-carolina-credit-and-debt-law**.

### Auto title loans
- **Mechanics:** short-term, high-cost loan secured by your **vehicle's title**
  (often lent at ~25–50% of the car's value), typically due as a single payment
  in ~30 days. A common **25%-per-month** finance charge works out to
  **~300% APR**.
- **The risk that makes them worse than payday:** if you can't pay, the lender
  can **repossess the car even if you've made partial payments**, and you may
  owe past-due payments, the full remaining balance, plus repossession/storage/
  sale/attorney costs to get it back. Some lenders install **GPS trackers and
  starter-interrupt devices** to disable and locate the car (FTC).

### Rent-to-own (RTO / lease-to-own)
- **Mechanics:** you rent furniture, appliances, or electronics with an option to
  own after a series of payments. Not always disclosed as financing, so the
  effective cost is easy to miss.
- **Watch for:** **total of payments far exceeding the cash price**; what happens
  if you pay late or miss — you can **lose the item with no equity**, and it may
  be repossessed; whether it helps or hurts your credit (FTC). Compare the
  all-in total to buying outright, layaway, or a secured-card purchase.

### Pawn loans
- **Mechanics:** a loan secured by a physical item (jewelry, electronics, tools)
  you hand over; redeem by repaying principal + fees by the deadline, or forfeit
  the item. **No personal liability beyond the item** and no credit-report
  furnishing. But monthly finance charges (varying by state) commonly run
  **~10–25% per month, which is ~120–300% APR**, and you risk losing property
  worth more than the loan. Convert the monthly fee (fee ÷ amount × 12) before
  borrowing.

### High-cost installment / "no credit check" loans
- **Mechanics:** longer-term, larger-balance loans (often $1,000–$10,000+) at
  triple-digit APRs marketed as a "safer" alternative to payday — but the long
  term means you can pay **far more in interest than principal**. Frequently
  **packed** with single-premium credit insurance and add-ons, and structured to
  encourage **flipping** (refinancing that resets fees).

### Overdraft as a form of credit
- When a bank covers a transaction that overdraws your account and charges a flat
  fee, that fee is effectively the price of a very short-term loan, so a small
  overdraft repaid in days carries an **enormous effective APR**. Example: a $35
  fee to cover a $20 overdraft repaid in 5 days works out to **over 12,000% APR**.
  EWA and cash-advance repayments that ACH-debit your
  account can **trigger overdraft/NSF fees** if funds are short (CFPB), stacking
  cost on cost.

### Earned-wage access (EWA) / cash-advance apps — *frame as of 2026*
- **What they are:** apps (e.g., Earnin, Dave, Brigit, MoneyLion) that advance a
  portion of wages you've "already earned" before payday, repaid by ACH on
  payday. Marketed as fee-free or "tip"-based.
- **Real cost (CFPB Data Spotlight):** users of employer-sponsored products take
  out an **average of ~27 advances/year**, and the typical employer-sponsored
  advance carries an **APR over 100%** once **expedite fees** ($1–$5.99, avg
  ~$3.18), **subscription fees** (up to $14.99/mo), and **"tips"** are counted.
- **Regulatory status is unsettled, so date-stamp it:** in **July 2024** the CFPB
  proposed an interpretive rule treating many paycheck-advance products as TILA
  credit (requiring cost disclosure). That direction then **reversed**: the CFPB
  **rescinded the 2020 EWA advisory opinion in January 2025**, and a **December
  2025 advisory opinion** took the position that genuine EWA is **not**
  Regulation Z credit (that same December 2025 action also **withdrew the July
  2024 proposed interpretive rule**). The federal treatment flipped twice in
  about a year; **as of 2026 it remains in flux and state EWA laws vary**, so
  re-verify before relying on any disclosure regime.

### Buy now, pay later (BNPL) — brief
- Short-term "pay-in-4" point-of-sale installment plans. Main risks: **easy
  over-extension and stacking** multiple plans, late fees, and disputes when a
  return collides with a payment schedule. **For how BNPL is reported or scored,
  see credit-reports-and-scores.**

---

## Hallmarks of predatory lending (terms to know)

- **Equity stripping**: sizing a (usually mortgage) loan to your collateral
  equity rather than ability to repay; the lender wins on foreclosure.
- **Loan flipping**: repeated refinancing that generates new fees/points without
  net benefit to the borrower.
- **Packing**: bundling overpriced add-ons (single-premium credit insurance,
  auto club, GAP) into the financed amount.
- **Balloon payments**: small payments then a large lump sum that forces a
  refinance or default.
- **Prepayment penalties**: fees for paying off early, locking in interest.
- **Mandatory arbitration / class-action waivers**: stripping access to courts.
- **Abusive add-ons & junk fees**: application, "processing," or "membership"
  fees that inflate the real cost.

### HOEPA high-cost mortgages (brief)
The Home Ownership and Equity Protection Act (an amendment to TILA, Regulation Z)
designates certain **high-cost mortgages** by APR, points-and-fees, or
prepayment-penalty thresholds, and then **bans or restricts** the worst
features — balloon payments, prepayment penalties, and loan flipping — and
**requires pre-loan homeownership counseling**. (Conceptual summary; for the
statutory thresholds and rule text see **us-consumer-credit-and-debt-law**.)

---

## Military Lending Act — 36% MAPR for servicemembers

The **Military Lending Act (MLA)** protects **active-duty servicemembers and
their dependents** ("covered borrowers") on most consumer credit:

- **36% MAPR cap**: the **Military Annual Percentage Rate** caps the all-in cost
  at 36%, and unlike a normal APR it **folds in most fees and certain add-on
  charges** (e.g., application/participation fees and credit-insurance premiums),
  so a "low rate, high fees" structure can't evade it.
- **No mandatory arbitration**, no class-action waiver of SCRA/other rights.
- **No prepayment penalties**; **no mandatory military allotment** as a repayment
  condition.
- **Required military-specific disclosures.**

The 36% MAPR is also the benchmark of the broader **state rate-cap movement**
(below). For the underlying statute/DoD rule text, see
**us-consumer-credit-and-debt-law**.

---

## Safe alternatives (lead with these)

- **Credit-union Payday Alternative Loans (PALs):** small loans from federal
  credit unions, expressly designed as a low-cost payday substitute — **no
  balloon payment** and strict limits on loan count/term; far cheaper than payday
  or title loans (NCUA). *(For PAL eligibility and how to qualify/build credit
  with one, see **improving-and-rebuilding-credit**.)*
- **Ask your creditor for more time / a payment plan**: many will grant an
  extension or revised schedule before you turn to high-cost credit (FTC).
- **Nonprofit lenders & CDFIs** (Community Development Financial Institutions) —
  mission-driven small-dollar and emergency loans at fair rates.
- **Employer hardship / paycheck or assistance programs**: often a cheaper way to
  bridge a shortfall than a third-party EWA app.
- **Nonprofit credit counseling**: a reputable NFCC/accredited agency can set up
  a budget or debt-management plan. *(See **improving-and-rebuilding-credit**;
  beware advance-fee "credit repair" scams.)*
- **Local assistance**: 211/community aid, utility/medical hardship programs,
  faith- and community-based emergency funds.

---

## Frontier / things in motion (date-stamped, verify)

- **State 36% rate-cap movement (as of 2026):** a growing number of states cap
  small-loan APRs at ~36% all-in (the MLA benchmark), and some have done so by
  **ballot measure**. The map keeps changing — check your state regulator.
- **CFPB Payday Rule history:**
  - **2017**: CFPB finalized the *Payday, Vehicle Title, and Certain High-Cost
    Installment Loans* rule with two parts: **(a) mandatory-underwriting /
    ability-to-repay** provisions for short-term and balloon loans, and **(b)
    payment provisions** (lender can't keep re-attempting ACH withdrawals after
    **two consecutive failed attempts** without new authorization).
  - **July 7, 2020**: CFPB **revoked the ability-to-repay / mandatory-
    underwriting provisions**; the **payment provisions were retained** (and
    separately ratified).
  - **Litigation**: industry challenges (incl. *CFSA v. CFPB*) delayed the
    payment provisions' compliance date for years. Confirm the **current
    compliance/enforcement status** before relying on it.
- **EWA federal status**: see the EWA section above; **in flux as of 2026.**

---

## References / verify current

Primary sources (U.S. federal regulators and NC DOJ). Re-verify — rules and rates
change.

- CFPB: *What is a payday loan?*:
  https://www.consumerfinance.gov/ask-cfpb/what-is-a-payday-loan-en-1567/
- CFPB: *Why is the APR higher than the interest rate on my payday loan?*:
  https://www.consumerfinance.gov/ask-cfpb/what-is-an-annual-percentage-rate-apr-and-why-is-it-higher-than-the-interest-rate-for-my-payday-loan-en-1625/
- CFPB: *4 of 5 payday loans are rolled over or renewed*:
  https://www.consumerfinance.gov/about-us/newsroom/cfpb-finds-four-out-of-five-payday-loans-are-rolled-over-or-renewed/
- CFPB: *Payday Lending Rule (12 CFR Part 1041)* hub:
  https://www.consumerfinance.gov/rules-policy/regulations/1041/
- CFPB: *Payday rule — 2020 Revocation Rule* (ability-to-repay revoked; payment
  provisions retained, July 7, 2020):
  https://www.consumerfinance.gov/rules-policy/final-rules/payday-vehicle-title-and-certain-high-cost-installment-loans-revocation-rule/
- CFPB: *Consumer use of payday, auto title, and pawn loans* (Making Ends Meet):
  https://www.consumerfinance.gov/data-research/research-reports/consumer-use-of-payday-auto-title-and-pawn-loans-insights-making-ends-meet-survey/
- CFPB: *Data Spotlight: Developments in the Paycheck Advance Market* (EWA ~100%+
  APR, ~27/yr, fees/tips):
  https://www.consumerfinance.gov/data-research/research-reports/data-spotlight-developments-in-the-paycheck-advance-market/
- CFPB: *Paycheck advance / EWA rulemaking* (2024 proposal; status in flux):
  https://www.consumerfinance.gov/rules-policy/rules-under-development/consumer-credit-offered-to-borrowers-in-advance-of-expected-receipt-of-compensation-for-work/
- CFPB: *Military Lending Act (MLA)* overview & rights (36% MAPR):
  https://www.consumerfinance.gov/consumer-tools/military-financial-lifecycle/military-lending-act-mla/
  and https://www.consumerfinance.gov/ask-cfpb/what-are-my-rights-under-the-military-lending-act-en-1783/
- FTC: *What To Know About Payday and Car Title Loans* (title repossession,
  GPS/starter-interrupt, alternatives):
  https://consumer.ftc.gov/articles/what-know-about-payday-and-car-title-loans
- FTC: *Renting to Own, Lease-to-Own, Layaway, BNPL* (RTO total-cost risks):
  https://consumer.ftc.gov/articles/rent-own-lease-own-layaway-buying-over-time
- NCUA: *Payday Alternative Loans (PALs)* / payday rule guidance:
  https://ncua.gov/regulation-supervision/letters-credit-unions-other-guidance/cfpb-issues-amendments-payday-vehicle-title-and-certain-high-cost-installment-loans-rule
- NC DOJ: *Payday Loans* (NC ban; 390–780% APR warning; complaints
  1-877-5-NO-SCAM):
  https://ncdoj.gov/protecting-consumers/credit-and-debt/payday-loans/
- NC DOJ: *Predatory Loans* (NC rate caps, predatory-lending consumer guidance):
  https://ncdoj.gov/protecting-consumers/mortgages-home-loans/predatory-loans/

*Also useful (verify): National Consumer Law Center (nclc.org) reports on high-
cost lending and the state 36% rate-cap movement.*

---

### Cross-references
- **north-carolina-credit-and-debt-law**: NC payday ban, NC usury / small-loan
  caps (Ch. 53), and NC-specific repossession/deficiency detail.
- **us-consumer-credit-and-debt-law**: literal TILA/HOEPA/Reg Z and MLA statute &
  rule text and thresholds.
- **home-mortgage-lending**: mortgage types, underwriting, and the home-buying
  process (this skill keeps only the predatory-mortgage angle: equity stripping,
  the HOEPA high-cost summary).
- **improving-and-rebuilding-credit**: credit-union PAL details, nonprofit credit
  counseling, and building credit instead of borrowing high-cost.
- **credit-reports-and-scores**: how BNPL and these products report and score.
- **personal-banking** (`consumer-finance` sibling family): the deposit-account
  mechanics of overdraft and NSF — fee schedules, the Reg E overdraft opt-in,
  how the account works, and choosing a bank/credit union. This skill keeps only
  the *overdraft-as-credit* angle (the buried APR when a flat fee prices a tiny,
  days-long advance); the account side is personal-banking's.
- **debt-collectors-and-fdcpa-rights**: when a high-cost loan defaults and goes to
  collections.
