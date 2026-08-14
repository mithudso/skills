<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-credit-and-debt` hub.** Formerly the standalone `home-mortgage-lending` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-credit-and-debt`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: home-mortgage-lending
description: >-
  US home-mortgage lending fundamentals (credit/underwriting-centered), NC notes:
  mortgage types (conventional, FHA, VA, USDA, jumbo) and fit; how score & DTI
  drive approval/pricing; PMI vs MIP; the TRID buying process; NCHFA assistance;
  refinancing. Spoke of consumer-credit-and-debt. Educational, not advice.
  TRIGGER: which mortgage type suits me; minimum credit score for an FHA / VA /
  conventional mortgage; DTI to qualify; PMI vs MIP / when PMI drops off;
  pre-qual vs pre-approval; underwriting / the 3 C's; Loan Estimate or Closing
  Disclosure / TRID; rate lock; NCHFA / down-payment assistance; rate-term vs
  cash-out refinance; NC attorney closing.
  SKIP: how scores are built / FICO factors -> credit-reports-and-scores;
  raising a score pre-application -> improving-and-rebuilding-credit; NC
  foreclosure / garnishment -> north-carolina-credit-and-debt-law;
  TILA/RESPA/ECOA statute text -> us-consumer-credit-and-debt-law;
  homeowners/hazard/flood insurance at closing -> personal-insurance.
metadata:
  changelog:
    - "2026-06-16 sko (--meta --no-sync) — 1 Medium structural: added reciprocal SKIP -> personal-insurance (homeowners/hazard/flood coverage required at closing, coverage side); desc trimmed to hold <=1000 (993). consumer-credit-and-debt anchor unchanged. Linter 0/0/0. No content passes; --no-sync (registry left stale)."
---

# Home-Mortgage Lending (US) — buyer & underwriting reference

**General educational information only. This is NOT financial, lending, tax, or
legal advice.** Mortgage rates, program rules, fees, and credit-score minimums
**change frequently and vary by lender**; every number below is **as of 2026**
and approximate. Lenders add their own stricter rules ("**overlays**") on top of
program minimums, so a program floor is rarely what you will actually be offered.
**Confirm current details with a licensed lender, the loan program, or a
HUD-approved housing counselor before relying on anything here.**

This is a **spoke of the consumer-credit-and-debt family**. It owns the *mortgage*
slice: loan types, mortgage underwriting, and the buying process. See the SKIP
list in the description for sibling skills that own how scores are built, raising
a score, NC foreclosure law, and federal statute text.

> Find a HUD-approved housing counselor (free/low-cost): the CFPB "Find a
> Counselor" tool, or HUD at 1-800-569-4287.

---

## 1. The 3 C's — how a mortgage is actually decided

Underwriting (manual or automated) evaluates the **3 C's**:

- **Credit:** your willingness/history of repaying: credit score(s), payment
  history, derogatories. *How scores are built lives in
  `credit-reports-and-scores`; raising one before you apply lives in
  `improving-and-rebuilding-credit`.*
- **Capacity:** your ability to repay: income, employment stability, and
  **debt-to-income (DTI)** ratio. This is the heart of §2 below.
- **Collateral:** the property: its value (via **appraisal**) and the
  **loan-to-value (LTV)** ratio. The home secures the loan.

A weakness in one C can sometimes be offset by strength in another (these are
"compensating factors," e.g., large reserves offsetting a higher DTI). Most
loans run through an **automated underwriting system (AUS)**: Fannie Mae's
**Desktop Underwriter (DU)** or Freddie Mac's **Loan Product Advisor (LPA)**,
which returns an Approve/Eligible (or Refer) recommendation that a human
underwriter then verifies against the documents.

---

## 2. DTI & underwriting (the capacity core)

**Debt-to-income (DTI)** is the single most decisive number after credit. Two
ratios:

- **Front-end (housing) ratio** = proposed total monthly housing payment
  ÷ gross monthly income. The housing payment is **PITI**: **P**rincipal,
  **I**nterest, property **T**axes, homeowners **I**nsurance, plus HOA dues and
  any mortgage insurance (PMI/MIP).
- **Back-end (total) ratio** = (PITI + **all** other monthly debt payments:
  car loans, student loans, minimum credit-card payments, child support, etc.)
  ÷ gross monthly income. Back-end is the one lenders weigh most.

**The ~43% guideline.** 43% back-end DTI is the most-cited underwriting
threshold and a useful planning anchor. Its origin: the CFPB's original
**General Qualified Mortgage (QM)** definition capped DTI at 43%. **Important
nuance (as of 2026):** a 2021 CFPB amendment **replaced** the General QM's hard
43% DTI cap with a **price-based (APR-vs-APOR) threshold**, so 43% is no longer a
bright-line federal rule for General QM. In practice many lenders and AUS still
treat the low-to-mid 40s as a comfort zone, and stronger files routinely go
higher: **conventional AUS approvals can reach ~50% back-end** with
compensating factors, and **FHA can go higher still** with strong credit/
reserves. Treat 43% as a planning rule of thumb, not a wall.

**Income & asset documentation (typical "full doc" file):**

- **Income:** 2 years W-2s, recent pay stubs (30 days), and (for
  self-employed/commission/variable income) **2 years personal (and often
  business) tax returns**, plus year-to-date P&L. Lenders use a 2-year *average*
  for variable income.
- **Assets:** 2 months of bank/brokerage statements to source the **down
  payment + closing costs** and verify **reserves**. Large or irregular deposits
  must be sourced/explained; gift funds need a **gift letter**.
- **Reserves:** months of PITI you could cover after closing from remaining
  assets. Often 0–2 months for a primary residence, more for investment
  properties, multi-unit, jumbo, or as a compensating factor.

**Stability matters as much as amount:** underwriters want a documented,
likely-to-continue income stream (job gaps, recent career changes, and unverified
income all draw scrutiny).

---

## 3. Mortgage types — and who each suits

| Type | Backed by | Typical min down | Mortgage insurance | Best for |
|---|---|---|---|---|
| **Conventional / conforming** | Fannie Mae / Freddie Mac (GSE rules) | 3% (first-time/low-income programs) – 5%+ | **PMI** if <20% down; **cancellable** | Solid credit, want to drop MI later |
| **FHA** | HUD / FHA-insured | 3.5% (score ≥580) | **MIP** (upfront + annual); often **life-of-loan** | Lower scores / thinner credit / smaller down |
| **VA** | Dept. of Veterans Affairs-guaranteed | **$0** | **No** MI (one-time **funding fee**) | Eligible veterans, service members, surviving spouses |
| **USDA (Section 502 Guaranteed)** | USDA Rural Development | **$0** | **Guarantee fee** (upfront + annual) | Low/moderate income in eligible **rural** areas |
| **Jumbo (non-conforming)** | Private (exceeds GSE limit) | 10–20%+ | Varies by lender | High-cost areas / loans above the conforming limit |

**Conforming vs jumbo line.** A **conforming** loan is at or below the FHFA
**conforming loan limit**: **$832,750** for a one-unit home in 2026 (baseline),
up to **$1,249,125** in designated high-cost areas (higher still in AK, HI, Guam,
and the USVI). Above that it's a **jumbo** loan, which is non-conforming, held in
private portfolios, and underwritten more strictly (higher scores, bigger down
payments, more reserves). *Re-verify the current year's limit at fhfa.gov; it
resets annually.*

**Fixed vs adjustable rate.** **Fixed-rate** (rate locked for the full term,
commonly 30- or 15-year) is what most borrowers choose for payment certainty.
**ARMs** (e.g., 7/6) start lower then adjust to an index + margin; they suit
shorter expected tenures and carry rate-reset risk.

---

## 4. Credit score & DTI drive approval AND price

Two borrowers can both be approved yet get very different rates. Pricing is
**risk-based**, driven mainly by **credit score × LTV**.

**LLPAs (conventional).** Fannie/Freddie assess **Loan-Level Price Adjustments**:
upfront, risk-based price hits keyed to your **representative credit score and
LTV** (plus loan features like cash-out, investment property, condo). Lower
score and/or higher LTV = larger LLPA = higher rate or more cost. This is why a
760-score borrower is quoted meaningfully better than a 660-score borrower on the
same loan. (The LLPA matrix lives in the Fannie Mae Selling Guide; it is revised
periodically.)

**Typical minimum scores by program (as of 2026 — floors, not what you'll be
offered; lender overlays vary and are usually higher):**

- **FHA:** **580** to get the **3.5%-down** option; **500–579** requires **10%
  down**. (Many lenders overlay to 600–640 regardless of the FHA floor.)
- **Conventional:** historically a **~620** minimum representative score.
  **Update (as of 2026):** **Fannie Mae eliminated its hard 620 minimum
  representative-score requirement effective Nov 16, 2025** for DU-underwritten
  purchase/refinance decisions; DU now evaluates credit holistically rather than
  gating on a single floor. **Lender overlays still commonly enforce ~620–640**,
  and pricing (LLPAs) still penalizes lower scores heavily, so this is more a
  manual-floor change than a green light at any score.
- **VA:** **no statutory minimum credit score** — the VA does not set one.
  **Lender overlays fill the gap**, commonly **580–620+**.
- **USDA:** **no published program minimum**; the AUS (GUS) drives it. Lenders
  commonly want **~640** for streamlined approval.

**Takeaway:** raising a score before applying can both *qualify* you and *lower
your rate*. The how-to for that is `improving-and-rebuilding-credit`.

---

## 5. Down payment, PMI vs MIP

**PMI (Private Mortgage Insurance) — conventional loans, <20% down.** Protects
the *lender*, paid by *you*, and (crucially) **cancellable**:

- **Automatic termination:** the servicer must end PMI when the balance is first
  *scheduled* to reach **78% LTV of the original value** (and you're current).
- **Borrower-requested cancellation:** you may *request* cancellation at **80%
  LTV** (scheduled or via extra payments), subject to good payment history and
  servicer conditions (sometimes a new appraisal). This is your right under the
  federal **Homeowners Protection Act (HPA)**.
- A "**final-termination**" backstop requires PMI to end by the loan's
  amortization midpoint if it hasn't already.

**MIP (Mortgage Insurance Premium) — FHA loans.** Two parts: an **upfront MIP**
(financed into the loan, ~1.75% of the base amount) **plus annual MIP** paid
monthly. The catch is **duration**:

- Most FHA loans today with the minimum (3.5%) down carry MIP for the **life of
  the loan**; it does **not** auto-cancel at 78% the way PMI does.
- Only if you put **≥10% down** does annual MIP drop off after **11 years**.
- The common way out of life-of-loan FHA MIP is to **refinance into a
  conventional loan** once you have ~20% equity and qualifying credit.

**This MIP-duration difference is a core FHA-vs-conventional trade-off:** FHA is
easier to qualify for, but its mortgage insurance is often permanent; PMI on a
conventional loan goes away.

---

## 6. The buying process

1. **Pre-qualification vs pre-approval.**
   - *Pre-qualification*: a quick, informal estimate based on **stated**
     (unverified) info. Weak; not a commitment.
   - *Pre-approval*: the lender **pulls credit and reviews documents** and
     issues a conditional letter for a specific amount. Sellers take it
     seriously; get this before house-hunting.
2. **Application & the Loan Estimate (LE).** Once you give the **six items** that
   legally trigger an application (name, income, SSN, property address,
   estimated property value, loan amount), the lender must deliver the **Loan
   Estimate within 3 business days**. The LE is a standardized **3-page** form —
   use it to **compare offers across lenders** (rate, APR, monthly payment,
   closing costs, cash to close).
3. **Processing, appraisal & underwriting.** The lender orders an **appraisal**
   (independent opinion of value protecting the lender's collateral; a low
   appraisal can derail or reprice the deal) and underwrites the file against the
   3 C's, clearing "conditions."
4. **Rate lock.** Locking fixes your rate for a set window (e.g., 30/45/60 days)
   against market moves; **floating** leaves it exposed. Locks can be extended
   (usually for a fee).
5. **The Closing Disclosure (CD) & closing.** You must receive the **5-page
   Closing Disclosure at least 3 business days before closing**; compare it
   against your LE. **Certain changes restart the 3-day clock** (APR increase
   beyond tolerance, a prepayment penalty added, or a loan-product change). At
   closing you sign and (for purchases) take ownership.

> **TRID** = the **TILA-RESPA Integrated Disclosure** rule, which created the LE
> and CD (replacing the old GFE/TIL and HUD-1). The *statutory* TILA/RESPA text
> itself lives in `us-consumer-credit-and-debt-law`; this skill covers the
> practical forms and timing.

---

## 7. First-time buyer & down-payment assistance — North Carolina (NCHFA)

The **North Carolina Housing Finance Agency (NCHFA)** runs the state's main
buyer-assistance programs, offered **through participating lenders** (not by
NCHFA directly). As of 2026:

- **NC Home Advantage Mortgage™** — a stable fixed-rate mortgage (conventional,
  FHA, VA, or USDA) with **down-payment assistance up to 3% of the loan amount**,
  for **first-time *and* move-up** buyers. Common eligibility: credit score
  **≥640**, income within program limits (an overall cap around **$152,000** for
  this product), and DTI within limits.
- **NC 1st Home Advantage Down Payment™** — for **first-time buyers and military
  veterans**, a larger **down-payment assistance amount** (publicized at
  **$15,000**; some materials reference earlier $8,000 tiers — verify the current
  amount). County/household-size income limits apply.
- The DPA is structured as a **0%-interest deferred second mortgage** with this
  forgiveness schedule: **no forgiveness in years 1–10**, then **20%/year in
  years 11–15** (fully forgiven at the end of year 15). If you sell, refinance, or
  transfer the home before then, you repay only the portion **not yet forgiven**
  (the full balance in years 1–10; the declining remainder in years 11–14).
- **Mortgage Credit Certificate (MCC)** — historically offered by NCHFA to
  convert part of mortgage interest into a **federal tax credit** for first-time
  buyers; confirm current availability.

*Always confirm current amounts, score floors, and income caps at **nchfa.com** —
these change.* Federal first-time-buyer help also flows through **FHA, VA, and
USDA** loans (§3).

---

## 8. Refinancing basics

Replacing your existing mortgage with a new one:

- **Rate-and-term refinance** — change the **interest rate and/or term** (e.g.,
  30→15 year, or to drop the rate). Goal: lower payment or pay off faster. A
  conventional refi to **shed FHA life-of-loan MIP** once you have ~20% equity is
  a classic move (§5).
- **Cash-out refinance** — borrow **more than you owe** and take the difference
  in cash, tapping home equity. Higher rates/LLPAs than rate-term, lower max LTV,
  and the new loan re-underwrites your full 3 C's.
- **Streamline refinances** exist for some government loans (**FHA Streamline**,
  **VA IRRRL**) with reduced documentation when going FHA→FHA or VA→VA.
- Refinancing has **closing costs**; the usual test is the **break-even point**
  (months to recoup costs from monthly savings) versus how long you'll keep the
  loan. A primary-residence refinance also carries a **3-business-day right of
  rescission** (you can cancel within 3 days of signing).

---

## 9. North Carolina specifics (cross-references, not duplication)

- **NC is an attorney-closing (attorney-state) jurisdiction.** A **licensed NC
  attorney** must conduct the real-estate closing / title work; you cannot close
  a NC home purchase through a non-attorney escrow/title agent the way some
  Western states do. Budget for attorney fees in closing costs.
- **NC foreclosure is non-judicial "power-of-sale."** Most NC mortgages use a
  **deed of trust** with a power-of-sale clause, foreclosed through a **clerk of
  superior court hearing** rather than a full lawsuit. **The detailed NC
  foreclosure process, statute of limitations, garnishment, and homestead
  exemptions live in `north-carolina-credit-and-debt-law`; go there for that
  depth. This skill only flags that NC differs.**

---

## 10. Cross-references (consumer-credit-and-debt family)

- **How credit reports & scores are built; FICO factors; the FHFA mortgage-score
  transition →** `credit-reports-and-scores`
- **Raising / rebuilding a score *before* you apply (utilization, disputes,
  timelines) →** `improving-and-rebuilding-credit`
- **NC foreclosure-law detail, statute of limitations, garnishment, exemptions →**
  `north-carolina-credit-and-debt-law`
- **Federal statute text — TILA/Reg Z, RESPA, ECOA/Reg B, adverse-action,
  fair-lending →** `us-consumer-credit-and-debt-law`

*(Reciprocal pointers: those sibling skills should point mortgage-specific
questions — loan types, PMI/MIP, DTI/underwriting, LE/CD, NCHFA, refinance —
back here to `home-mortgage-lending`.)*

---

## References / verify current (primary sources, as of 2026)

Rates, fees, score floors, and program rules change — **re-verify before
relying.**

**CFPB — Owning a Home, TRID/LE & CD, PMI, QM**
- Loan options overview: https://www.consumerfinance.gov/owning-a-home/loan-options/
- Conventional loans: https://www.consumerfinance.gov/owning-a-home/conventional-loans/
- FHA loans: https://www.consumerfinance.gov/owning-a-home/fha-loans/
- Special loan programs (VA/USDA): https://www.consumerfinance.gov/owning-a-home/special-loan-programs/
- Loan Estimate & Closing Disclosure (Know Before You Owe): https://www.consumerfinance.gov/owning-a-home/closing-disclosure/ and https://www.consumerfinance.gov/know-before-you-owe/
- When do I get a Closing Disclosure (3-day rule): https://www.consumerfinance.gov/ask-cfpb/when-do-i-get-a-closing-disclosure-en-179/
- Removing PMI / HPA: https://www.consumerfinance.gov/ask-cfpb/when-can-i-remove-private-mortgage-insurance-pmi-from-my-loan-en-202/
- HPA examination procedures: https://www.consumerfinance.gov/compliance/supervision-examinations/homeowners-protection-act-hpa-or-pmi-cancellation-act-examination-procedures/
- Ability-to-Repay / Qualified Mortgage rule: https://www.consumerfinance.gov/rules-policy/final-rules/ability-to-pay-qualified-mortgage-rule/
- What is a Qualified Mortgage: https://www.consumerfinance.gov/ask-cfpb/what-is-a-qualified-mortgage-en-1789/
- Find a housing counselor: https://www.consumerfinance.gov/find-a-housing-counselor/

**HUD / FHA**
- FHA mortgage limits & info: https://www.hud.gov/program_offices/housing/sfh
- Single Family Housing Policy Handbook 4000.1: https://www.hud.gov/program_offices/administration/hudclips/handbooks/hsgh

**VA home loans**
- Eligibility: https://www.va.gov/housing-assistance/home-loans/eligibility/
- Funding fee & closing costs: https://www.va.gov/housing-assistance/home-loans/funding-fee-and-closing-costs/
- Loan limits & entitlement: https://www.va.gov/housing-assistance/home-loans/loan-limits/

**Fannie Mae / Freddie Mac / FHFA**
- Fannie Mae credit-score requirements (Selling Guide B3-5.1-01): https://selling-guide.fanniemae.com/sel/b3-5.1-01/general-requirements-credit-scores
- Fannie Mae eligibility & pricing (LLPA matrix): https://singlefamily.fanniemae.com/originating-underwriting/mortgage-products/eligibility-pricing
- FHFA conforming loan limit values: https://www.fhfa.gov/data/conforming-loan-limit

**USDA Rural Development**
- Single Family Housing Guaranteed Loan Program: https://www.rd.usda.gov/programs-services/single-family-housing-programs/single-family-housing-guaranteed-loan-program

**North Carolina Housing Finance Agency (NCHFA)**
- Home buyers landing: https://www.nchfa.com/home-buyers
- NC Home Advantage Mortgage: https://www.nchfa.com/home-buyers/buy-home/nc-home-advantage-mortgage
- NC 1st Home Advantage Down Payment: https://www.nchfa.com/home-buyers/home-buyer-mortgage-products/nc-1st-home-advantage-down-payment
