<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-credit-and-debt` hub.** Formerly the standalone `charge-offs-collections-and-debt-resolution` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-credit-and-debt`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: charge-offs-collections-and-debt-resolution
description: 'Resolving defaulted US consumer debt: what a charge-off means, how debt moves creditor -> agency -> debt buyer, how to settle, "pay for delete," the 1099-C tax hit, and re-aging traps that restart the clock. Spoke of consumer-credit-and-debt; educational, NOT financial/tax/legal advice. TRIGGER: what a charge-off means / do I still owe it; settle a debt (lump-sum vs plan); what a debt buyer is; "pay for delete"; paid vs settled vs paid-in-full; 1099-C / canceled-debt / insolvency / Form 982; will paying restart the clock. SKIP: collector conduct, validation, being sued -> debt-collectors-and-fdcpa-rights; NC SOL, garnishment -> north-carolina-credit-and-debt-law; rebuilding a score -> improving-and-rebuilding-credit; statute text -> us-consumer-credit-and-debt-law; how a notation ages/moves a score -> credit-reports-and-scores; medical-bill charity care / prompt-pay discount / IRS 501(r) -> medical-debt-and-billing; affordable payoff planning -> budgeting-and-saving.'
metadata:
  changelog:
  - '2026-06-16 sko --meta --no-sync — structural pass vs new consumer-finance family; 0 High 0 Med (resolved 1 circular-SKIP Medium: medical "settlement" loop -> one-way specificity gradient). Added 2 SKIP edges (-> medical-debt-and-billing, -> budgeting-and-saving); trimmed description to 981 chars within 1000 cap.'
---

# Charge-offs, collections, and debt resolution

The lifecycle of a defaulted US consumer debt, from the moment a creditor writes
it off through how you actually resolve it. This is a **spoke of the
`consumer-credit-and-debt` hub** (do not edit the hub here).

> **Framing — read first.** This is **general educational information, not
> financial, tax, or legal advice**, and is current **as of 2026**. Laws,
> creditor policies, and IRS rules change. For any 1099-C / cancellation-of-debt
> or insolvency question consult a **tax professional or CPA**; for disputes,
> lawsuits, or statute-of-limitations questions consult a **consumer-law
> attorney**. Verify every figure against the primary sources at the end before
> relying on it.

---

## 1. What a charge-off actually is (and is NOT)

A **charge-off** is an *accounting* action by the original creditor, not a
forgiveness of what you owe.

- After a revolving account (e.g., a credit card) is roughly **180 days past
  due** — installment loans often sooner — the creditor moves the balance from
  an asset to a loss on its own books for regulatory/tax reasons. This is the
  "charge-off."
- **You still legally owe the full balance.** Charging off is the creditor's
  bookkeeping; it does not erase the debt, stop interest/fees where permitted,
  or end collection. The CFPB and FTC are explicit that paying or settling later
  is still expected.
- A charged-off account is typically reported to the credit bureaus with a
  **"charge-off"** status, one of the most damaging negative items on a report.

**On the credit report / how long.** A charge-off (and the related delinquency
or collection) can generally be reported for **up to seven years**. Under the
FCRA the seven-year clock runs from the **date of first delinquency** — the
month the account first went late and never recovered before charge-off — **not**
from the charge-off date and **not** reset by later collection activity. (How an
item is displayed/aged and how it moves a score belongs to
`credit-reports-and-scores`; the statute itself belongs to
`us-consumer-credit-and-debt-law`.)

> Charge-off (a credit-reporting / accounting event) is a **different clock**
> from the **statute of limitations** (how long you can be *sued*). They are
> unrelated, run from possibly different dates, and expire independently. Do not
> conflate them. (SOL detail + NC specifics -> the cross-referenced law skills.)

---

## 2. The collections lifecycle — know who you're dealing with

A defaulted debt usually travels down this chain. Where it sits changes who has
authority to settle and what they paid for it.

1. **Original creditor, in-house collections.** The lender's own recovery team,
   pre- or post-charge-off. They own the debt and have the most flexibility.
2. **Assigned / contingency agency.** A third-party collector working the debt
   *on behalf of* the creditor for a commission. The creditor still **owns** it;
   the agency is paid a percentage of what it recovers.
3. **Debt buyer (sold debt).** The creditor **sells** the charged-off account —
   often in bulk portfolios for **pennies on the dollar** — to a debt-buying
   company that now *owns* it and collects for its own account. Because their
   cost basis is tiny, debt buyers often have wide room to settle, but the data
   they buy can be thin or stale, so **verify the debt** before engaging.

**Practical point:** A sold debt may appear on your report **twice**: the
original creditor's charge-off (often showing a $0 balance once sold) and the
debt buyer's collection entry. That is generally permissible if the original
shows it was transferred/sold and zeroed; duplicate *active* balances are not.

> Stopping the calls, demanding **debt validation**, collector misconduct, and
> being sued are **collector-rights** topics -> `debt-collectors-and-fdcpa-rights`.

---

## 3. Decision flow — validate, then negotiate in writing, then handle tax

```
        ┌─────────────────────────────────────────────────────────────┐
        │  A debt is charged off / in collections / a debt buyer calls  │
        └─────────────────────────────────────────────────────────────┘
                                   │
                                   ▼
   STEP 1 — VALIDATE FIRST (before you say or pay anything)
   • Is this debt really mine, the right amount, not already paid?
   • Is it within the statute of limitations (could I be sued)?  Is it past it
     (time-barred)?  ── SOL math -> NC / federal law skills; validation steps
     and the validation notice -> debt-collectors-and-fdcpa-rights.
   • CAUTION: do NOT make a payment OR admit the debt is yours until you know
     the SOL status; it can RESTART the clock (see §8).
                                   │
                                   ▼
   STEP 2 — NEGOTIATE, AND GET IT IN WRITING *BEFORE* PAYING
   • Decide lump-sum vs payment plan (see §4).
   • Get the agreement (amount, that it settles the ENTIRE debt, and how it
     will be reported) in a SIGNED letter from the current owner BEFORE you
     send a dollar.  Keep the letter and proof of every payment forever.
                                   │
                                   ▼
   STEP 3 — UNDERSTAND THE TAX CONSEQUENCE
   • Forgiven/settled balance of $600+ may generate a 1099-C and be TAXABLE
     income.  Check whether the insolvency or bankruptcy exclusion applies
     (Form 982).  Talk to a tax pro.  (see §7)
                                   │
                                   ▼
                          Resolve → then rebuild
                  (rebuilding the score -> improving-and-rebuilding-credit)
```

---

## 4. Settling the debt — lump-sum vs payment plan

Some creditors and collectors will accept **less than the full balance** to
close the account. The CFPB notes you can pay in full or **propose a repayment
plan**, and the FTC notes some collectors will **accept less than what you owe**.

**Lump-sum settlement.** A one-time payment for a fraction of the balance, in
exchange for the debt being marked resolved. Often the cheapest total outcome and
the cleanest to document. Settlement amounts vary widely by age, owner, and your
leverage; **debt buyers** (tiny cost basis) frequently settle for a meaningful
discount, but there is **no guaranteed percentage** — anyone promising a fixed
"X cents on the dollar" outcome is overselling. Treat any percentage you see
online as illustrative only.

**Payment plan.** Structured installments. Useful if you can't raise a lump sum,
but watch for (a) re-aging risk if the account was near or past the SOL (§8),
(b) interest/fees continuing to accrue, and (c) a missed payment voiding the
deal.

**Non-negotiables for either path:**

- **Get every term in writing, signed, BEFORE you pay.** Per the CFPB: *"If you
  agree to a repayment or settlement plan, get the plan and the debt collector's
  promises in writing before you make a payment."* Per the FTC: *get a signed
  letter from the collector that says the amount you're paying settles the
  entire debt — and you no longer owe anything for that debt.*
- **Confirm who owns the debt** and that they have authority to settle it. Pay
  the **current owner**, not a prior holder.
- **Keep the letter and a record of every payment** indefinitely. Sold/resold
  debts resurface; your written settlement is your proof.
- **Pin down how it will be reported** (see §5) — "settled," "paid in full,"
  deletion, balance to $0 — in the same letter.
- Verify the SOL status first (§8) so a payment doesn't reopen lawsuit exposure.

> A **debt-settlement company** generally can't get better terms than you can
> get yourself, can't guarantee savings, and may charge fees and tell you to
> stop paying creditors (which deepens the damage). The CFPB distinguishes this
> from nonprofit credit counseling. DIY negotiation is usually the better first
> move.

---

## 5. Paid vs settled vs "paid in full" — what the report says

How a resolved account is *labeled* matters, and you should negotiate the label
in your written agreement:

- **Paid in full:** you paid the entire balance owed. The charge-off/collection
  remark generally **remains** on the report for its seven years, but now shows a
  zero balance and "paid." Best of the realistic outcomes.
- **Settled / settled for less than full balance / "paid, settled":** you paid
  less than owed and the creditor accepted it. The account shows **resolved but
  for less than the full amount**; the FTC notes *"if you settle the debt, some
  collectors will report that on your credit report to show you didn't pay the
  full amount."* The negative item still ages off at seven years.
- **Charged off, balance owing:** unresolved. Most damaging.

Paying or settling **does not delete** the history and does **not** reset the
seven-year clock to your benefit; the item still ages off based on the original
date of first delinquency. The score *impact* of each notation ->
`credit-reports-and-scores`.

---

## 6. "Pay for delete" — what it is and why it's shaky

**Pay for delete** = you offer to pay (often a settlement) **in exchange for the
collector deleting the negative tradeline** from your credit reports entirely.

Reality check:

- **Accurate negative information generally cannot be required to be removed.**
  The CFPB: *"You generally cannot have negative information removed from your
  credit report if it is accurate,"* and warns to *"beware of anyone who claims
  that they can remove information from your credit report that's current,
  accurate, and negative — it's probably a credit repair scam."*
- **Bureau and furnisher policies discourage it.** The major credit bureaus'
  furnisher agreements and reporting guidelines (e.g., the Metro 2 framework)
  call for accurate reporting and discourage deleting accurate accounts in
  exchange for payment. Many furnishers will refuse outright.
- **Efficacy is mixed and unenforceable.** Some smaller collectors/debt buyers
  *will* agree informally, but there's no obligation, and a verbal promise is
  worthless. **If you attempt it, get the deletion promise in the same signed
  letter, before paying** — though even then a furnisher may decline to action
  the bureaus.
- A debt buyer can only delete **its own** tradeline; it cannot remove the
  **original creditor's** charge-off entry.

Bottom line: useful to *ask for* in writing, never to *count on*. Disputing
genuine **errors** is a separate, free right (-> `us-consumer-credit-and-debt-law`
and `credit-reports-and-scores`); do not pay anyone to do that for you.

---

## 7. The 1099-C — canceled debt can be taxable income

This is the step people forget. **When a creditor forgives or cancels debt, the
IRS generally treats the forgiven amount as taxable income to you.**

- **The $600 trigger.** A creditor (applicable financial entity) generally must
  file **Form 1099-C, Cancellation of Debt,** when it cancels **$600 or more** of
  debt and an identifiable event occurs. You may get a copy in the mail.
- **It's generally taxable.** The IRS: *"In general, if your debt is canceled,
  forgiven, or discharged for less than the amount owed, the amount of the
  canceled debt is taxable"* and is reported as income for the year of
  cancellation. **This applies to the forgiven portion of a settlement** — settle
  a $10,000 debt for $4,000 and the ~$6,000 forgiven can be reportable income.
- **You must report it whether or not you receive the form**, and you're
  responsible for the correct amount even if the 1099-C is wrong.

**Exclusions that can wipe out the tax (talk to a tax pro):**

- **Insolvency exclusion.** You're **insolvent** when your **total liabilities
  exceed the fair market value of your total assets** immediately before the
  cancellation. You can exclude the canceled debt **up to the amount of your
  insolvency**. This is the most common lifeline for someone settling distressed
  debt.
- **Bankruptcy.** Debt discharged in a **Title 11 bankruptcy** is excluded.
- **Other exclusions/exceptions** include qualified principal-residence
  indebtedness, certain student-loan cancellations, and gifts.

**The mechanics:** to claim an exclusion you generally file **IRS Form 982,
Reduction of Tax Attributes Due to Discharge of Indebtedness** (e.g., check the
insolvency box and report the excluded amount), and **IRS Publication 4681**
walks through canceled debt, the insolvency worksheet, and Form 982 with
examples.

> 1099-C / insolvency math is genuinely tricky (asset/liability valuation,
> attribute reduction). **Get a tax professional** before filing. This skill
> flags the issue; it does not give tax advice.

---

## 8. Re-aging and restart traps — the most expensive mistakes

Two different clocks can be **restarted**, and people trip them by accident:

1. **Statute-of-limitations restart (lawsuit exposure).** The CFPB: *"Making a
   partial payment or acknowledging you owe an old debt, even after the statute
   of limitations expired, may restart the time period."* A single payment — or
   even a *promise to pay* or written admission — on a **time-barred** debt can
   **revive** the creditor's right to sue you. **This is the big one.** Before
   you pay or acknowledge an old debt, confirm the SOL status (NC specifics ->
   `north-carolina-credit-and-debt-law`; the doctrine ->
   `us-consumer-credit-and-debt-law`).

2. **Credit-report re-aging (illegal kind).** The seven-year reporting clock runs
   from the **original date of first delinquency** and is **not** supposed to be
   reset by new activity, by a sale to a debt buyer, or by a partial payment.
   When a collector improperly reports a **newer** delinquency date to keep an
   old item on your report longer, that's unlawful **re-aging** of the
   tradeline; dispute it (-> `credit-reports-and-scores` /
   `us-consumer-credit-and-debt-law`).

**Defensive rules:**

- **Validate and date the debt before paying or admitting anything.**
- Don't let a collector talk you into a "good-faith payment" on an old account
  without first checking whether it's time-barred.
- Get the reporting outcome and the date treatment in your written agreement.
- If you're near the seven-year report fall-off or the SOL line, get advice
  before acting; a payment can cost you more than it saves.

---

## 9. Quick playbook (TL;DR)

1. **Charge-off ≠ forgiven.** You still owe it; it reports ~7 years from first
   delinquency.
2. **Find out who owns it** — original creditor, agency, or debt buyer — and
   **validate** the debt.
3. **Check the SOL first.** Don't pay or admit an old debt until you know whether
   it's time-barred; a payment can **restart** the lawsuit clock.
4. **Negotiate** lump-sum or plan. **Get a signed agreement BEFORE paying**,
   stating it settles the entire debt and how it'll be reported.
5. **"Pay for delete" is a nice-to-ask, never-to-count-on.** Accurate items
   usually can't be forced off.
6. **Plan for a 1099-C.** Forgiven $600+ may be taxable; check the **insolvency**
   exclusion and **Form 982** with a **tax pro**.
7. **Keep every document forever.** Then rebuild (->
   `improving-and-rebuilding-credit`).

---

## Cross-references (sibling spokes of `consumer-credit-and-debt`)

- **`debt-collectors-and-fdcpa-rights`** — stopping collectors, the validation
  notice, debt-validation letters, collector misconduct, being sued.
- **`north-carolina-credit-and-debt-law`** — NC statute of limitations,
  garnishment rules, wage exemptions.
- **`improving-and-rebuilding-credit`** — rebuilding a score after the debt is
  resolved.
- **`us-consumer-credit-and-debt-law`** — the text/mechanics of FCRA, FDCPA, and
  other federal statutes; the re-aging and time-barred-debt doctrines.
- **`credit-reports-and-scores`** — how a charge-off/collection is displayed,
  aged, and how each notation moves a score.

---

## References / verify current (primary sources, as of 2026)

Confirm any figure against these before relying on it; URLs and rules change.

- **CFPB — How do I negotiate a settlement with a debt collector?**
  https://www.consumerfinance.gov/ask-cfpb/how-do-i-negotiate-a-settlement-with-a-debt-collector-en-1447/
  (get the plan/promises in writing before paying)
- **CFPB — Can debt collectors collect a debt that's several years old?**
  https://www.consumerfinance.gov/ask-cfpb/can-debt-collectors-collect-a-debt-thats-several-years-old-en-1423/
  (statute of limitations; partial payment/acknowledgment may **restart** it)
- **CFPB — How long does information stay on my credit report?**
  https://www.consumerfinance.gov/ask-cfpb/how-long-does-information-stay-on-my-credit-report-en-323/
  (most negative info ~7 years; bankruptcy up to 10)
- **CFPB — Is it possible to remove accurate but negative information from my
  credit report?**
  https://www.consumerfinance.gov/ask-cfpb/is-it-possible-to-remove-accurate-negative-information-from-my-credit-report-en-1249/
  (accurate negatives generally can't be removed; "pay for delete" / credit-repair
  scam warning)
- **CFPB — What is the difference between credit counseling and debt settlement,
  debt consolidation, or credit repair?**
  https://www.consumerfinance.gov/ask-cfpb/what-is-the-difference-between-credit-counseling-and-debt-settlement-debt-consolidation-or-credit-repair-en-1449/
- **FTC — Debt Collection FAQs** (consumer.ftc.gov/articles/debt-collection-faqs)
  https://consumer.ftc.gov/articles/debt-collection-faqs
  (settling for less; **get a signed settlement letter**; settling may be
  reported as not-paid-in-full; paying old debt may not erase report history)
- **FTC — How To Get Out of Debt**
  https://consumer.ftc.gov/articles/getting-out-debt
  (get any settle/manage agreement in writing; how it affects credit)
- **IRS — Topic No. 431, Canceled Debt – Is It Taxable or Not?**
  https://www.irs.gov/taxtopics/tc431
  (canceled debt generally taxable; bankruptcy & insolvency exclusions; Form 982)
- **IRS — Instructions for Forms 1099-A and 1099-C**
  https://www.irs.gov/instructions/i1099ac
  (creditor files Form 1099-C for canceled debt of **$600 or more** + identifiable
  event)
- **IRS — About Form 982 (Reduction of Tax Attributes Due to Discharge of
  Indebtedness)** and **Instructions for Form 982**
  https://www.irs.gov/forms-pubs/about-form-982  •
  https://www.irs.gov/instructions/i982
  (how to claim the insolvency/bankruptcy exclusion)
- **IRS — Publication 4681 (Canceled Debts, Foreclosures, Repossessions, and
  Abandonments)**
  https://www.irs.gov/publications/p4681
  (insolvency worksheet, exclusions, Form 982 examples)
- **IRS — Taxpayer Advocate Service: I Have a Cancellation of Debt or Form
  1099-C**
  https://www.taxpayeradvocate.irs.gov/get-help/general/cancellation-of-debt/

*This document is general educational information, not financial, tax, or legal
advice (as of 2026). Consult a tax professional for 1099-C/insolvency questions
and an attorney for disputes or statute-of-limitations questions.*
