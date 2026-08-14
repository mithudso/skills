<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-credit-and-debt` hub.** Formerly the standalone `auto-lending-and-financing` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-credit-and-debt`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: auto-lending-and-financing
description: >-
  US auto loans / car financing + NC repo notes. Spoke of
  consumer-credit-and-debt. Educational, not advice.
  TRIGGER: auto loan / car financing; APR by credit tier; dealer vs CU financing;
  dealer markup / buy rate / reserve / F&I office; negotiate/refinance a car loan;
  pre-approval; negative equity / upside-down; 72/84-month loan; GAP insurance;
  extended warranty / VSC; buy-here-pay-here / no-credit-check; starter-interrupt
  / GPS kill-switch; lease vs buy; repossession / deficiency; redeem.
  SKIP: scores/FICO/inquiries -> credit-reports-and-scores; raising a score ->
  improving-and-rebuilding-credit; NC statute/garnishment/SOL ->
  north-carolina-credit-and-debt-law; home loans -> home-mortgage-lending;
  statute text -> us-consumer-credit-and-debt-law; charge-off/post-repo collector
  -> charge-offs-collections-and-debt-resolution /
  debt-collectors-and-fdcpa-rights; predatory credit ->
  predatory-lending-and-high-cost-credit; auto liability/collision/comprehensive
  coverage -> personal-insurance.
metadata:
  changelog:
    - "2026-06-16 sko (--meta --no-sync) — 1 Medium structural: added reciprocal SKIP -> personal-insurance (auto liability/collision/comp coverage; GAP & F&I add-ons stay owned here). Resolved circular-SKIP [gap,add,ons] flagged by meta-validate into a one-way gradient by dropping the GAP/F&I parenthetical from the edge. Desc trimmed to <=1000 (1000). Anchor unchanged. Linter 0/0/0. No content passes; --no-sync."
---

# Auto Lending & Financing

Practical buyer reference for US auto loans. Spoke of the **consumer-credit-and-debt**
hub. **General educational information, NOT financial, lending, or legal advice.**
Rates, terms, and tiers vary by lender, state, and time and change frequently —
all figures are **as of 2026** and must be re-verified (see References below).

---

## 1. How an auto loan is structured

- **Secured by the vehicle.** The car is collateral. The lender holds a **lien**
  (is the lienholder on the title) until the loan is paid off; default lets the
  lender **repossess** (Section 9). This is why auto-loan APRs are far below
  unsecured personal-loan/credit-card APRs.
- **Principal, term, payment.** You finance the price (less down payment and any
  trade-in equity), plus any rolled-in fees, tax, and add-ons. **Term** = number
  of months; common terms run 36–84 months. Longer term = lower monthly payment
  but **more total interest** and longer time underwater (Section 6).
- **APR vs interest rate.** The **interest rate** is the cost of borrowing the
  principal. The **APR (annual percentage rate)** includes the interest rate
  **plus certain lender fees**, so it is the better number for comparing offers.
  APR ≥ interest rate. TILA/Regulation Z requires the APR to be disclosed (for
  statute text -> `us-consumer-credit-and-debt-law`).
- **Simple-interest amortization.** Almost all US auto loans are **simple
  interest**: interest accrues daily on the **current outstanding balance**.
  Early in the loan most of each payment goes to interest; later, most goes to
  principal. Practical consequences:
  - Paying **early or extra** reduces principal and the interest you'll owe —
    there is usually **no benefit lost** (unlike precomputed-interest loans).
  - Paying **late** means more days of accrued interest, so more of your next
    payment goes to interest and less to principal.
  - Watch for a **prepayment penalty** (rare on mainstream auto loans, more
    common in subprime/BHPH) and for **precomputed interest** (Rule of 78s),
    which penalizes early payoff — both are red flags.

---

## 2. Credit-tier pricing (how APR scales with score)

Lenders price by **credit tier**. Industry tiers (Experian/VantageScore bands)
and **representative average APRs are as of Q4 2025** (Experian *State of the
Automotive Finance Market*); your offer depends on lender, term, down payment,
**loan-to-value (LTV: the amount financed vs. the car's value)**, income, and
the vehicle:

| Tier | Score band | New-car APR | Used-car APR |
|---|---|---|---|
| **Super-prime** | 781–850 | ~5% | ~7% |
| **Prime** | 661–780 | ~6.5% | ~9% |
| **Nonprime** | 601–660 | ~9.5% | ~14% |
| **Subprime** | 501–600 | ~13% | ~18–19% |
| **Deep subprime** | 300–500 | ~15–16% | ~21%+ |

(Overall Q4 2025 averages were roughly **6.4% new / 11.3% used** across all
tiers.) Takeaways:
- The spread from super-prime to subprime is enormous — often **8–15+
  percentage points**, which can mean **thousands to tens of thousands** of extra
  dollars over the loan.
- A **higher score, larger down payment, shorter term, and lower LTV** all push
  your rate down. Even a small score improvement that bumps you into the next
  tier can cut the rate sharply.
- **How scoring works, FICO vs VantageScore, what moves a score** is out of
  scope here -> `credit-reports-and-scores`; **raising your score before
  applying** -> `improving-and-rebuilding-credit`.

---

## 3. Direct vs dealer (indirect) financing

**Direct lending** — you get a loan **directly** from a bank, credit union, or
online lender, ideally **pre-approved before you shop**. You walk in knowing your
rate and a max amount; the car purchase and the financing are separate
negotiations. **Credit unions** are frequently the cheapest source.

**Dealer-arranged (indirect) financing** — the dealer's **F&I (Finance &
Insurance) office** takes your application and forwards it to one or more lenders
(often ~5: banks, credit unions, captive/nonbank finance companies). The dealer
is the middleman; the loan is then assigned to the lender.

### The dealer rate markup ("reserve")

This is the single most important fee to understand (CFPB):

- The lender quotes the dealer a **buy rate** — the rate the lender is actually
  willing to fund. *"A buy rate is the interest rate that a financial institution
  quotes to the dealer when you apply for dealer-arranged financing."* (CFPB)
- The dealer may then offer **you** a higher **contract rate**. The difference
  is the **dealer markup / dealer reserve / dealer participation**, and the
  lender **shares that extra interest revenue with the dealer**.
- A markup of a few percentage points over the life of the loan can cost
  hundreds to thousands of dollars — for the **same borrower and same lender**.

### What this means for you

- **The interest rate is negotiable.** *"You can negotiate the terms of your auto
  loan"* and dealers/lenders *"are not required to offer the best rates available"*
  (CFPB). Bring a pre-approval and make the dealer **beat** it.
- Negotiate the **car price, the financing, and any trade-in separately** so a
  "good" monthly payment can't hide a marked-up rate or a long term.
- Dealer financing **can** occasionally win (manufacturer **subvented/captive**
  promo rates like 0–2.9% APR on specific new models). Compare the **APR and
  total cost**, not the monthly payment, and read whether a low APR forces you to
  give up a cash rebate.

---

## 4. Shopping & pre-approval

1. **Check your credit first** and fix obvious errors before applying
   (`improving-and-rebuilding-credit`).
2. **Get pre-approved** at 2–3 places (your bank, a **credit union**, an online
   lender) so you know your real rate and budget before the lot.
3. **Rate-shop inside the de-duplication window.** Multiple **auto-loan**
   inquiries within a short window (commonly **14 days**, up to **45 days** on
   newer scoring models) are bundled and counted as **a single inquiry** for
   scoring, so shopping many lenders barely dents your score. Mechanics of
   inquiries and the shopping window -> `credit-reports-and-scores`.
4. At the dealer, let them try to beat your pre-approval, then take the best
   **APR/total-cost** offer. Decline add-ons you don't want (Section 5).

### Refinancing an existing auto loan

Replacing your current loan with a new one from another lender. It can help
when:
- **Your credit improved** since you bought (you've moved up a tier — Section 2),
- **Rates dropped**, or
- **You took a marked-up dealer rate** and want to escape it (often refinanceable
  within the first months through a bank or credit union).

Compare the **new APR and total remaining interest**, not just the payment.
Caution: extending the **term** to lower the payment can re-extend how long you
stay underwater (Section 6), and check for any prepayment penalty on the old
loan.

---

## 5. Add-ons (GAP, extended warranties / service contracts)

The F&I office sells **optional** products that **increase your loan amount** (so
you also pay interest on them). CFPB: *"You are not required to purchase add-on
products to secure financing through the dealer and you may have the right to
cancel them later."* *"Dealers are not the only sellers of most of these
products"* — shop them.

- **GAP (Guaranteed Asset Protection).** If the car is **totaled or stolen** and
  you owe **more than its insurance-payout value**, GAP can cover the
  **"gap."** Worth considering when you're **upside-down** — small down payment,
  long term, or a fast-depreciating vehicle. Often **much cheaper** from your
  **own insurer or credit union** than the dealer's price; can usually be
  **canceled for a prorated refund** if you pay down/refinance. Not useful once
  you have positive equity.
- **Extended warranty / vehicle service contract (VSC).** Pays for certain
  repairs after the factory warranty ends. Frequently **marked up heavily** and
  **negotiable**; many overlap the factory warranty, exclude common failures, or
  go unused. Compare third-party/manufacturer pricing and the actual coverage and
  deductible. **Credit insurance** (life/disability) is likewise optional and
  usually a poor value.
- Rule of thumb: **add-ons are negotiable and cancellable.** Decide on each one
  **on its own merits**, never as a condition of getting the loan, and never let
  them be slipped into the payment unexplained.

---

## 6. Negative equity / "upside-down" loans

You're **upside-down (underwater)** when you **owe more than the car is worth** —
common early in a long loan because cars depreciate faster than the balance falls.

- **Rolling negative equity into a new loan.** Dealers will offer to pay off your
  old underwater loan and **add the shortfall to the new loan**. This stacks old
  debt on a new (also-depreciating) car, **instantly** putting you deeper
  underwater and raising payments and total interest. Avoid; if unavoidable,
  minimize the amount rolled and keep the term short.
- **Long terms (72–84 months)** lower the monthly payment but **keep you
  underwater longer**, pile on interest, and raise the odds you'll trade in still
  upside-down — feeding the cycle. They also pair dangerously with low/no down
  payment.
- **Mitigation:** larger down payment, **shorter term**, buy less car, keep the
  car past payoff, and consider **GAP** (Section 5) while underwater.

---

## 7. Subprime, buy-here-pay-here (BHPH) & device traps

For thin-file or damaged credit, mainstream lenders may decline, pushing buyers
toward higher-cost channels — handle with care.

- **Buy-here-pay-here / "no credit check."** The **same dealer sells the car AND
  originates and services the loan**, typically targeting **subprime** buyers
  (CFPB). Expect **very high APRs**, weekly/biweekly in-person payments, older
  high-mileage cars at inflated prices, and **hidden charges** (CFPB cited a
  *required* repair warranty and a *required* GPS payment-reminder device baked
  into the deal). Subprime loans from finance companies/BHPH default far more
  often than bank loans.
- **Starter-interrupt / GPS "kill-switch" devices.** Many subprime/BHPH lenders
  install devices that let them **remotely disable the car** or track it. CFPB
  enforcement has found vehicles wrongly disabled. Know whether a device is
  installed and what triggers it.
- **Safer alternatives:** a **credit-union** auto loan (including credit-builder
  options), a **larger down payment**, buying a cheaper car outright, or
  **rebuilding credit first** (`improving-and-rebuilding-credit`) then financing.
- High-cost lending patterns generally (rate caps, debt traps, MLA) ->
  `predatory-lending-and-high-cost-credit`.

---

## 8. Leasing vs buying (brief)

- **Buying (financing):** you **own** the car; build equity; no mileage limits;
  cheaper long-term if you keep it past payoff; you bear depreciation and resale.
- **Leasing:** you pay for the **use** (roughly the depreciation over the term)
  plus rent/fees; **lower monthly payments**; you **don't own** it at lease end
  (unless you buy it out). Watch **mileage caps** (often **~15,000/yr**, with
  per-mile overage fees), charges for **excess wear and damage**, **early-
  termination** penalties, and required insurance (Federal Reserve, *Keys to
  Vehicle Leasing*). Best for drivers who want a new car every few years and stay
  within mileage; usually worse for high-mileage or long-keep buyers.

---

## 9. Repossession & deficiency (UCC Article 9 + NC)

When you default on a **secured** car loan, the lienholder can repossess. The
framework is **UCC Article 9** (adopted in every US state, with state-specific
overlays).

### Self-help repossession

- **Self-help repossession is allowed** without going to court **as long as the
  repossessor does not "breach the peace"** (no force, no breaking into a closed
  garage, no confrontation). **No advance notice** is generally required before
  the car is taken.
- Default can be triggered by missed payments **or other contract breaches**
  (e.g., **letting required insurance lapse**) (NC DOJ).

### After repossession — your rights

- **Right to redeem.** You can usually **redeem** (get the car back) by paying
  the **full balance plus repossession costs** before the lender sells it. Some
  contracts/states also allow **reinstatement** (catch up the missed payments +
  fees) — check your contract and state law.
- **Notice and sale.** The lender must sell the car in a **commercially
  reasonable** manner and must **notify you** of a **public sale** (date, time,
  place — you may attend and bring bidders) or of the **date after which** a
  **private sale** will occur (NC DOJ).
- **Surplus or deficiency.** Sale proceeds are applied to repossession costs and
  the loan. If the sale **exceeds** what you owe, the lender must **refund the
  surplus**. If it **falls short**, the remaining balance is the **deficiency
  balance**, which you still owe and can be **sued for** (NC DOJ).

### North Carolina specifics (summary — verify in the NC skill)

- **NC DOJ** confirms self-help repossession, no required advance notice, the
  redeem/notice/sale/deficiency framework above, and that a **surplus** must be
  refunded.
- A deficiency can become a **court judgment**. As high-level orientation only:
  NC is unusually debtor-protective on collection of such a judgment (wage
  garnishment is generally not available, though bank funds may be reachable) —
  the exact reach, exemptions, and procedure are NC-statute questions for
  `north-carolina-credit-and-debt-law`, not this skill.
- For NC statute citations (UCC Article 9 as enacted, the **NC Debt Collection
  Act**, the **3-year statute of limitations** (G.S. 1-52), NC garnishment and
  property/homestead exemptions, and the full deficiency procedure) ->
  **`north-carolina-credit-and-debt-law`**.

### After the deficiency

- A repossession and any deficiency **hurt your credit** and the deficiency can
  be **charged off and sold** to a debt buyer or sent to collections.
- **Settling the deficiency, pay-for-delete, 1099-C canceled-debt tax** ->
  `charge-offs-collections-and-debt-resolution`.
- **A collector contacting you, debt validation, being sued, FDCPA rights** ->
  `debt-collectors-and-fdcpa-rights`.

---

## Cross-references

- **How credit scores/reports are built, FICO vs VantageScore, hard/soft
  inquiries, the rate-shopping de-dup window** -> `credit-reports-and-scores`
- **Raising/rebuilding a score before applying, credit-builder tools** ->
  `improving-and-rebuilding-credit`
- **NC repossession/deficiency statute detail, NC garnishment, NC SOL,
  exemptions** -> `north-carolina-credit-and-debt-law`
- **Home loans / mortgages** -> `home-mortgage-lending`
- **FCRA / ECOA / TILA-Reg Z / FDCPA statute text and enforcement** ->
  `us-consumer-credit-and-debt-law`
- **Settling a deficiency/charge-off, 1099-C** ->
  `charge-offs-collections-and-debt-resolution`
- **Dealing with a collector after repossession** ->
  `debt-collectors-and-fdcpa-rights`
- **High-cost/predatory lending generally** ->
  `predatory-lending-and-high-cost-credit`

---

## References / verify current (as of 2026)

Rates, tiers, and rules change — re-check these **primary** sources:

- **CFPB:** Auto loans (consumer tools):
  https://www.consumerfinance.gov/consumer-tools/auto-loans/
- **CFPB:** Dealer-arranged vs bank/credit-union financing:
  https://www.consumerfinance.gov/ask-cfpb/what-is-the-difference-between-dealer-arranged-and-bank-financing-en-759/
- **CFPB:** What is a "buy rate"? (dealer markup):
  https://www.consumerfinance.gov/ask-cfpb/what-is-a-buy-rate-for-an-auto-loan-en-727/
- **CFPB:** What is a Finance and Insurance (F&I) department?:
  https://www.consumerfinance.gov/ask-cfpb/what-is-a-finance-and-insurance-fi-department-en-747/
- **CFPB:** "No credit check" / buy-here-pay-here auto loans:
  https://www.consumerfinance.gov/ask-cfpb/what-is-a-no-credit-check-or-buy-here-pay-here-auto-loan-en-887/
- **CFPB:** Data Point: Subprime Auto Loan Outcomes by Lender Type:
  https://files.consumerfinance.gov/f/documents/cfpb_subprime-auto_data-point_2021-09.pdf
- **FTC:** Financing or Leasing a Car:
  https://consumer.ftc.gov/articles/financing-or-leasing-car
- **Federal Reserve:** Keys to Vehicle Leasing (lease vs buy):
  https://www.federalreserve.gov/pubs/leasing/
- **Experian:** State of the Automotive Finance Market (tier APRs; quarterly):
  https://www.experian.com/blogs/ask-experian/average-car-loan-interest-rates-by-credit-score/
- **NC DOJ:** Car Repossession:
  https://ncdoj.gov/protecting-consumers/automobiles/car-repossession/
