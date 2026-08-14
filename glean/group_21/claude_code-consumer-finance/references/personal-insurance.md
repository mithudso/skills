<!-- hub-reference-banner -->
> **Reference file — part of the `consumer-finance` hub.** Formerly the standalone `personal-insurance` skill.
> Sibling topics in this family are now reference files under the hubs (`consumer-finance`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: personal-insurance
description: >-
  US consumer property/casualty + life + disability insurance, with North Carolina notes. consumer-finance spoke. Educational, not advice; "as of 2026," verify w/ NC DOI.
  TRIGGER: how insurance works (premium, deductible, coverage limit, exclusions, claims, underwriting); auto (BI/PD liability, collision, comprehensive, UM/UIM, med-pay, NC minimum limits); homeowners vs renters, HO-3, replacement cost vs ACV, flood is separate (NFIP), coastal NC wind/hail + the NC Beach Plan / NCIUA; life term vs whole/universal; disability own-occ vs any-occ; umbrella/excess liability; how much to buy; how to shop; common gaps and mistakes.
  SKIP: health / ACA / Medicare / Medicaid / HSA → health-insurance-and-coverage; auto loans/financing/APR, GAP insurance & auto F&I add-ons (warranties, service contracts) → auto-lending-and-financing; the medical BILLS / surprise billing / medical collections → medical-debt-and-billing; credit-based insurance SCORE mechanics → credit-reports-and-scores.
version: 1.1.0
updated: 2026-06-16
metadata:
  changelog:
    - "2026-06-16 sko v0->v1.0.0: Pass H 10/10 pos, 10/10 neg (predicted); fixed 5 Medium (M/G desc 1351->986 chars under soft cap; K em-dash density 2.47->0.63 per 100w; A1 NC UIM-at-floor phrasing; A2 liability notation 2-num->3-num BI/BI/PD; F1 when-not-to-file-a-claim) + Lows (A4 UM hit-and-run, B8 auto-loan SKIP hardening); 0 banned terms"
    - "2026-06-16 sko v1.0.0->v1.1.0 (--meta --no-sync): structural reconciliation — Pass A parent-hub anchor updated consumer-credit-and-debt (stale 'hub to be created') -> consumer-finance now installed, sibling-hub cross-link retained; A'/N resolvability closed (all SKIP/cross-ref ids resolve); desc unchanged at 986 chars (<=1000 cap); Pass I reciprocal peer edges into out-of-set siblings REPORTED not written"
---

# Personal Insurance

Practical, coverage-by-coverage reference for US consumer **property/casualty + life + disability** insurance, with **North Carolina** notes. Spoke of the **`consumer-finance`** hub. Its sibling hub, **`consumer-credit-and-debt`**, owns the credit/debt/collections side (auto loans, mortgages, credit scores, collectors); route there for the financing or debt angle on a covered loss.

> **FRAMING (read first).** This is **general educational information, NOT insurance, legal, tax, or financial advice.** Coverage forms, limits, exclusions, prices, and state minimums **change** and vary by insurer and policy. Everything here is **as of 2026**. Your own policy's declarations page and forms control; read them. Verify anything important with a **licensed agent** or the **NC Department of Insurance (ncdoi.gov)** before you act.

> **Health insurance is out of scope.** Medical/health coverage (ACA, Medicare, employer health, HSA) is a **separate sibling skill**; see SKIP routing in the frontmatter.

---

## 1. How insurance works (the shared vocabulary)

Every P&C, life, and disability policy runs on the same machinery:

- **Premium:** what you pay (monthly/semi-annual/annual) to keep the policy in force. Set by **underwriting/rating** (below).
- **Deductible:** what *you* pay out of pocket on a covered loss before the insurer pays. Can be a **flat dollar amount** (e.g., $1,000) or a **percentage** of a coverage limit (common for wind/hail and named-storm losses; see §3). **Higher deductible, lower premium**: it's how you and the insurer share risk.
- **Coverage limit:** the **maximum** the insurer will pay for a covered loss. Liability limits are usually written as three numbers, **BI per-person / BI per-accident / PD per-accident** (e.g., 50/100/50). In homeowners, the limits for other structures, personal property, and loss of use are usually set as **percentages of the dwelling (Coverage A) limit**.
- **Exclusions:** what the policy **does not** cover. Reading these matters more than reading the covered perils. Classic homeowners exclusions: **flood, earthquake, normal wear/maintenance, mold (often capped), intentional acts**.
- **Endorsement / rider:** an add-on that **modifies** the base policy (adds coverage, raises a sublimit, e.g., scheduling jewelry).
- **Named-peril vs open-peril (all-risk):** a **named-peril** form covers only perils it lists; an **open-peril** form covers everything **except** what it excludes. (HO-3 is open-peril on the structure, named-peril on contents; see §3.)

**The claims process (typical):** (1) ensure safety, mitigate further damage, and document with photos; (2) **report the claim** to the insurer promptly; (3) an **adjuster** investigates and estimates the loss; (4) you provide a **proof of loss / inventory**; (5) insurer pays **limit minus deductible**, subject to ACV vs replacement-cost rules (§3). For auto, **at-fault** facts and police reports drive who pays.

**Think before filing a small claim.** A claim near or below your deductible can *cost* you: it pays little or nothing, yet it can raise your premium and sits on your **C.L.U.E.** loss-history report for years (insurers price off claim frequency). For small losses, paying out of pocket is often cheaper than filing.

**How insurers underwrite & price.** Underwriting means assessing your risk, grouping you with similar risks, and deciding whether to insure you and at what price. Common inputs: **claims history** (insurers pull a **C.L.U.E.** loss-history report), property/vehicle characteristics, location, coverage/deductible choices, and, in most states, a **credit-based insurance score**.

> **Credit-based insurance score (cross-ref).** Many insurers use a credit-based insurance score in pricing where state law allows. The *mechanics of how that score is built* are out of scope here; see **credit-reports-and-scores**. (NC restricts but does not fully ban its use; specifics change, so verify with **ncdoi.gov**.)

---

## 2. Auto insurance

A personal auto policy bundles several distinct coverages:

| Coverage | What it pays for | Pays whom |
|---|---|---|
| **Bodily Injury (BI) liability** | Injuries you cause to *other* people | Others |
| **Property Damage (PD) liability** | Damage you cause to *others'* property | Others |
| **Collision** | Damage to *your* car from a crash/rollover, regardless of fault | You |
| **Comprehensive ("other than collision")** | *Your* car: theft, fire, vandalism, flood, hail, falling objects, animal strikes | You |
| **Uninsured Motorist (UM)** | You're hit by an at-fault driver with **no** insurance (or a hit-and-run) | You |
| **Underinsured Motorist (UIM)** | At-fault driver's limits are **too low** to cover your damages | You |
| **Medical Payments (MedPay)** | Medical bills for you/your passengers, regardless of fault | You |

- **Liability is the core.** Most states (and lenders, via collision/comprehensive) require it. **Collision and comprehensive are usually optional** but are typically **required by an auto lender/lessor** while there's a loan. (Lender add-ons like **GAP** belong to **auto-lending-and-financing**, not here.)
- **Limits matter more than the state minimum.** State minimums are floors, not recommendations. A serious injury claim can blow past them and leave **you** personally liable for the rest, so many advisors suggest carrying meaningfully higher BI limits (and UM/UIM to match).

### North Carolina specifics (as of 2026; verify at ncdoi.gov)

- **Required in NC:** **liability** *and* **uninsured motorist (UM)** coverage. Collision, comprehensive, UIM (above the floor), and MedPay are optional.
- **NC minimum liability limits, increased effective July 1, 2025** (for policies new/renewed on or after that date):
  - **Bodily injury: $50,000 per person / $100,000 per accident**
  - **Property damage: $50,000 per accident**
  - (The **prior** minimums, before July 1, 2025, were **30/60/25**. If you see an old policy or older guidance citing 30/60/25, it's outdated.)
- **UM/UIM:** UM is mandatory. Under the old **30/60/25** floor, a minimum-limits policy carried **no UIM** because, by statute, UIM only attaches when your BI limits *exceed* the state minimum (at the floor there is nothing to be "under"-insured against). The 2025 changes also revised how a vehicle is judged "underinsured," now based on the **total damages** the claimant sustained rather than purely comparing policy limits. Confirm your specific UIM trigger with your agent.
- A **credit-based insurance score** can factor into NC auto pricing where allowed; see **credit-reports-and-scores** for the score mechanics.

---

## 3. Homeowners & renters

### The HO-3 homeowners policy (the most common form)

HO-3 is **open-peril on the structure** (covers all causes of loss except those excluded) and **named-peril on personal property**. Its standard coverage parts:

- **Coverage A — Dwelling:** the house structure itself. Set this to the **cost to rebuild** (replacement cost), *not* market value or the mortgage balance.
- **Coverage B — Other Structures:** detached garage, fence, shed (typically ~10% of Coverage A).
- **Coverage C — Personal Property (contents):** your belongings (often ~50–70% of Coverage A). High-value items (jewelry, firearms, collectibles) have **sublimits**, so schedule them with an endorsement to fully cover.
- **Coverage D — Loss of Use (Additional Living Expense):** hotel, temporary rent, extra meals while your home is uninhabitable after a covered loss.
- **Liability (Coverage E)** and **Medical Payments to others (Coverage F):** if someone is injured on your property or you're sued for bodily injury/property damage you cause (including pets).

**Replacement Cost vs. Actual Cash Value (ACV), the single most important homeowners choice:**
- **Replacement cost (RCV):** pays to repair/replace with new, **no deduction for depreciation** (you typically get the depreciated amount first, then the rest after you actually repair/replace and submit receipts).
- **Actual cash value (ACV):** pays **replacement cost minus depreciation**, so you eat the depreciation. Cheaper premium, much smaller checks.
- Watch for **peril-specific ACV**: a policy can be RCV overall but pay **ACV on roof wind/hail losses**, so read the declarations.

### Renters insurance (HO-4)

Renters insurance covers your **personal property**, **personal liability**, and **loss of use**; it does **not** cover the building (that's the landlord's policy). It's inexpensive and badly under-bought. Same **ACV vs replacement-cost** choice applies, and **replacement-cost contents** is usually worth the small extra premium. A landlord's policy does **nothing** for your stuff or your liability.

### Why flood is separate (NFIP)

**Standard homeowners and renters policies EXCLUDE flood.** Flood damage is covered only by a **separate flood policy**, usually through the federal **National Flood Insurance Program (NFIP)** (or a private flood insurer). Key points:

- NFIP is a **single-peril** policy: it pays only for **direct physical loss from flood**.
- NFIP **building** coverage is typically paid on a **replacement-cost** basis for a primary residence meeting conditions; **contents** are generally **ACV**.
- There's usually a **30-day waiting period** before a new NFIP policy takes effect, so you can't buy it as a storm approaches.
- You can be flooded **outside** a high-risk zone; "not in a flood zone" ≠ "no flood risk." Find your risk and policies at **floodsmart.gov**.

### Coastal NC wind/hail & the NC Beach Plan (NCIUA)

In coastal North Carolina, **wind and hail** drive the market:

- **Windstorm/hail can be carved out or carry a separate, percentage deductible.** A **named-storm (hurricane) deductible** is commonly a **percentage of Coverage A** (e.g., 2–5% of dwelling value), not a flat dollar amount, and triggers when the National Weather Service issues a watch/warning for a named storm. On a $400k home a 5% deductible = **$20,000** out of pocket, so know yours before a storm.
- **The NC "Beach Plan" = the Coastal Property Insurance Pool, run by the NC Insurance Underwriting Association (NCIUA).** It's the **insurer of last resort** for wind/hail (and limited fire) coverage in the **18 eligible coastal counties** (and the narrower "beach area") when you can't get coverage in the standard market. It exists so coastal owners can get *some* coverage, not cheap coverage. Info: NCIUA (ncjua-nciua.org) or **ncdoi.gov**.
- **Mitigation can lower cost:** NC offers **FORTIFIED home** construction credits/discounts; ask your agent and see **ncdoi.gov**.

---

## 4. Life insurance

Two big families:

- **Term life:** pure death benefit for a set period (10/20/30 years). **Cheap**, no cash value, expires at term end. Covers a defined need (income replacement while kids are young, while a mortgage is outstanding). Example order of magnitude: a healthy 30-year-old might pay roughly **$160/yr for a 20-year, $250k term policy** (illustrative; quotes vary by health/age/insurer).
- **Permanent life (whole, universal, variable):** lasts your whole life **and** builds **tax-deferred cash value**, but premiums are **much higher** for the same death benefit. Universal life adds flexible premiums/death benefit tied to a cash account; variable ties cash value to investments.

**Who needs it / how much.** If **someone depends on your income**, you likely need coverage to replace it; it also covers **final expenses, debts, and estate costs**. Sizing rules of thumb (e.g., 10–12× income, or DIME: **D**ebt + **I**ncome replacement + **M**ortgage + **E**ducation) are starting points; match coverage to the actual dollars your dependents would need and the years they'd need them.

> **Term-vs-whole caveat.** "Buy term and invest the difference" is the common consumer-advocate stance: for most families with a **temporary need**, term covers far more for far less, and you invest the premium savings elsewhere. **Permanent** life is mainly for **lifelong** needs (e.g., estate liquidity, a special-needs dependent, certain business/tax situations), not as a primary "investment." Cash-value/whole policies are complex and high-commission; be skeptical of a sales pitch that frames whole life as an investment account. Verify suitability independently.

---

## 5. Disability insurance

Disability income insurance replaces **part of your paycheck** if illness/injury stops you from working; your earning ability is usually your biggest asset.

- **Short-term disability (STD):** short **elimination (waiting) period** (0–14 days), benefits for ~**3–6 months** (sometimes up to ~2 years). Often employer-provided.
- **Long-term disability (LTD):** longer elimination period (weeks to months; a **90-day** wait is common, and a **longer wait means a lower premium**), benefits for **years up to retirement age** (or life).
- **Benefit amount:** typically **~60% of pre-disability income** (private policies usually replace **50–70%**, because insurers leave a gap so you have incentive to return to work). Benefits from a policy **you** paid for with after-tax dollars are generally **tax-free**.

**Own-occupation vs. any-occupation, the key definition:**
- **Own-occ:** pays if you can't perform **your own occupation** (more protective, more expensive). A surgeon who can't operate but could teach still collects.
- **Any-occ:** pays only if you can't do **any gainful work you're reasonably suited for** (cheaper, stingier). Many group LTD plans start own-occ for ~24 months, then switch to any-occ.

---

## 6. Umbrella (personal excess liability)

A **personal umbrella policy** sits **on top of** the liability limits in your auto, homeowners, and renters policies and pays once those underlying limits are exhausted. It also covers some claims the base policies don't, such as **personal-injury offenses (libel, slander, false arrest)**.

- Sold in **$1M increments**; relatively cheap because it rarely pays.
- Insurers usually require you to **already carry high underlying limits** first (commonly ~**$250k auto BI** and ~**$300k homeowners liability**) before they'll sell an umbrella.
- **Who should consider it:** anyone with **assets/income or future earnings worth protecting** beyond their base liability limits (homeowners, landlords, people with pools/dogs/teen drivers, higher net worth).

---

## 7. How to shop & compare; common gaps and mistakes

**Shopping:**
- Get **3+ quotes** for the **same coverages, limits, and deductibles**; only apples-to-apples comparisons mean anything. Price isn't the only factor; check the insurer's **financial strength** and **claims/complaint record** (NC publishes complaint info via **ncdoi.gov**).
- **Bundle** auto + home/renters for multi-policy discounts; ask about deductible, mitigation (FORTIFIED, alarms), and good-driver credits.
- Re-shop and re-evaluate after **life events** (marriage, kid, home purchase, new car, big asset growth).

**Common gaps & mistakes:**
- **Underinsuring the dwelling:** insuring to market value or loan balance instead of **rebuild (replacement) cost**; leads to under-payment (and coinsurance penalties) at claim time.
- **Carrying only state-minimum auto liability:** a bad injury claim can exceed it and expose your personal assets.
- **Picking ACV to save a few dollars:** then getting depreciated payouts on a total loss. **Skipping renters insurance** entirely.
- **Assuming flood is covered:** it isn't; and ignoring the **30-day NFIP wait**.
- **Not knowing your wind/named-storm deductible** is a **percentage** (huge out-of-pocket on the coast).
- **Letting coverage lapse:** even a short gap can raise rates, void coverage when you have a loss, and (for auto) create legal/registration problems.
- **Treating whole life as an investment** without independent advice (§4).

---

## References / verify current (as of 2026)

Primary, authoritative sources — re-check these for current figures, especially **NC auto minimums** and **coastal wind** rules, which change:

- **NC Department of Insurance (NCDOI)** — auto & homeowners consumer pages, NC minimums, windstorm/hail, FORTIFIED credits, complaint data:
  - https://www.ncdoi.gov/consumers/auto-and-vehicle-insurance/basic-and-miscellaneous-auto-coverages
  - https://www.ncdoi.gov/changes-rating-automobile-insurance-policies-effective-july-1-2025
  - https://www.ncdoi.gov/consumers/homeowners-insurance/windstorm-and-hail
  - https://www.ncdoi.gov/consumers/homeowners-insurance
- **NC Insurance Underwriting Association (NCIUA / "Beach Plan", Coastal Property Insurance Pool)** — https://www.ncjua-nciua.org
- **NAIC (National Association of Insurance Commissioners)** — consumer insights & glossary (how insurance works, home/renters, life, disability, umbrella):
  - https://content.naic.org/consumer/auto-insurance.htm
  - https://content.naic.org/glossary-insurance-terms
  - https://content.naic.org/consumer/life-insurance.htm
  - https://content.naic.org/article/whats-umbrella-policy
- **Insurance Information Institute (III / Triple-I)** — Insurance Handbook, HO-3 sample form, coverage basics:
  - https://www.iii.org/sites/default/files/docs/pdf/Insurance_Handbook_20103.pdf
  - https://www.iii.org/sites/default/files/docs/pdf/HO3_sample.pdf
  - https://www.iii.org/article/what-are-principal-types-life-insurance
  - https://www.iii.org/article/what-are-types-disability-insurance
  - https://www.iii.org/article/what-umbrella-liability
- **FEMA / NFIP / FloodSmart** — flood is separate; RCV vs ACV; coverage summary:
  - https://www.floodsmart.gov
  - https://agents.floodsmart.gov/articles/actual-cash-value-replacement-cost-value-and-what-flood-insurance-covers
- **CFPB (Consumer Financial Protection Bureau)** — consumer-protection angle on insurance products tied to loans (e.g., force-placed/lender-placed insurance) — https://www.consumerfinance.gov

---

## Cross-references (related skills)

- **health-insurance-and-coverage** — anything health/medical: ACA, Medicare, Medicaid, employer health plans, HSAs. (Out of scope here.)
- **auto-lending-and-financing** — **GAP insurance** and auto **F&I add-ons** (extended warranties, service contracts) sold at the dealership.
- **medical-debt-and-billing** — the medical **bills** themselves, surprise/balance billing, and medical collections (vs. the insurance that pays them).
- **credit-reports-and-scores** — the **mechanics of credit-based insurance scores** (how the score is built and what moves it).
- **consumer-finance** — the parent **hub** for this spoke (the personal-finance router: banking, taxes, health coverage, budgeting, investing, estate planning).
- **consumer-credit-and-debt** — the **sibling hub** for the credit/debt/collections side (credit scores, charge-offs, collectors, mortgages, auto loans, bankruptcy, and US/NC credit-debt law).
