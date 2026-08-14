<!-- hub-reference-banner -->
> **Reference file — part of the `venture-business` hub.** Formerly the standalone `venture-startup-fundraising-deck-model` skill.
> Sibling topics in this family are now reference files under the hubs (`venture-business`, `venture-nonprofit-cause`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: venture-startup-fundraising-deck-model
description: >-
  The raise-money / model-deeply layer for a lean founder, above the foundational
  business-model + unit-economics planning. Serves BOTH a for-profit founder (investor pitch
  deck, integrated 3-statement financial model, cap table & dilution, SAFE vs convertible note
  vs priced round) AND a nonprofit (the grant-writing lifecycle — prospect research, letter of
  inquiry, full proposal, logic model, project budget, reporting). Covers the standard 10-12 slide
  deck sequence, narrative arc, the deck-vs-memo question, and common deck mistakes; how the P&L,
  cash-flow statement, and balance sheet link in a driver-based model with scenario/sensitivity
  analysis; pre/post-money valuation, the option pool and the option-pool shuffle, post-money SAFE
  math, and worked seed-round dilution. TRIGGER: building or
  reviewing an investor pitch deck or investment memo; a 3-statement / integrated financial model;
  a cap table, dilution, pre/post-money valuation, option pool, SAFE, convertible note, or priced
  round; or the grant-writing lifecycle (LOI, full proposal, logic model, grant budget, reporting)
  for a nonprofit. SKIP: foundational business model / break-even / basic unit economics & lean
  projections → venture-small-business-planning; grant PROSE & persuasion writing, exec-summary
  craft → executive-comms; 501(c)(3) eligibility, charitable-solicitation registration & Form 990
  filing → venture-nc-nonprofit-formation.
category: personal-venture
tags: [venture, fundraising, pitch-deck, financial-model, cap-table, grants]
whenToUse:
  - Building or reviewing an investor pitch deck (the 10-12 slide sequence and narrative arc)
  - Deciding between a pitch deck and an investment memo, or building both
  - Constructing a 3-statement integrated financial model (P&L, cash flow, balance sheet linked)
  - Moving from a lean one-page projection to a real driver-based, scenario-tested model
  - Reading or building a cap table and computing dilution through a seed round
  - Choosing between a post-money SAFE, a convertible note, and a priced round
  - Sizing pre/post-money valuation and the option pool (and spotting the option-pool shuffle)
  - Running the grant-writing lifecycle for a nonprofit (LOI, proposal, logic model, budget, reporting)
triggers:
  - pitch deck / investor deck / investment memo
  - 3-statement model / integrated financial model / link P&L cash flow balance sheet
  - cap table / dilution / pre-money / post-money valuation
  - SAFE / convertible note / priced round / valuation cap / option pool
  - the ask / how much to raise / use of funds
  - grant proposal / letter of inquiry / LOI / logic model / grant budget
---

# Startup Fundraising: Deck, Model, Cap Table & Grants

The **raise-money / model-deeply** layer for a lean North Carolina founder. This sits one level
above `venture-small-business-planning` (the Business Model Canvas, basic unit economics, and the
lean one-page projection). When the founder is past "is this a business?" and into "how do I get
someone to fund it and how do I model what happens when they do?" — that's here.

It serves **two ventures at once**:

- **The for-profit** → investor pitch deck, a real integrated financial model, a cap table, and the
  SAFE / convertible-note / priced-round decision.
- **The nonprofit** (organ-donation awareness) → the grant-writing lifecycle: prospect research,
  letter of inquiry, full proposal, logic model, project budget, and reporting.

**Boundaries (cross-reference, do not duplicate):**

- Foundational business model, break-even, basic unit economics (CAC/LTV/contribution margin), and
  the *lean* one-page projection → **`venture-small-business-planning`**. This skill is the
  *advanced* layer; it assumes those exist and builds on them.
- **Prose craft** — the actual persuasive sentences of a grant narrative, the exec summary, the
  founder/shareholder letter, slide copy as writing → **`executive-comms`** (and `writing-expert`
  for general voice). This skill says *what goes in each section and why*; `executive-comms` makes
  the words land.
- **501(c)(3) recognition, NC charitable-solicitation registration, and Form 990 filing** (which
  together gate *grant eligibility*) → **`venture-nc-nonprofit-formation`**. Most foundation and
  government grants require 501(c)(3) status, so confirm that is in hand before chasing grants.
- Spreadsheet mechanics (the actual `.xlsx`, which cell, which formula) → the pro library's
  `document-formats` / xlsx skill. This skill says *what to model*, not where to type it.

Rough order of work, but loop freely: **deck ↔ model** (each disciplines the other) → **cap table**
(once you have a number to raise) → **grants** run on their own parallel track for the nonprofit.

> Not financial, legal, tax, or investment advice. SAFE terms, cap-table mechanics, and tax
> treatment must be confirmed with a startup attorney and a CPA; standards evolve. See the
> Disclaimer.

---

## 1. The investor pitch deck

A pitch deck's *only* job is to earn the next meeting. It is not a business plan, not the contract,
and not where you win the deal — it is the door. Investors skim it in **2-3 minutes**, so clarity
beats completeness. Y Combinator names three design principles: **legibility, simplicity, and
obviousness** — investors want clarity, traction, and a believable model, not design tricks or
filler [YC; Storydoc].

> **Boundary with `executive-comms`:** this section owns the *fundraising structure* — which slides,
> in what order, with what content and what numbers, and the common mistakes that get a deck passed
> on. The *narrative-arc writing craft* (the actual slide copy, the 10/20/30 rule, prose persuasion)
> lives in `executive-comms` (`references/pitch-deck-writing.md`). Build the spine here; write the
> words there.

### 1.1 The standard 10-12 slide sequence

Both the **Sequoia** template (10 slides + appendix) and the **YC** guidance converge on the same
spine. The canonical order [Sequoia; YC; Storydoc]:

| # | Slide | One job | Watch-outs |
|---|-------|---------|-----------|
| 1 | **Title / "what you do"** | Company name + one declarative sentence a non-expert understands ("Sequoia: company purpose") | Vague taglines; jargon |
| 2 | **Problem** | The pain, in one sentence investors instantly get | Inventing a problem; "vitamin not painkiller" |
| 3 | **Solution** | What you built, in plain language | Tech deep-dives; feature lists |
| 4 | **Why now** | The recent shift (tech, regulatory, behavioral) that makes this possible *today* | "No why-now" = why didn't this exist already? |
| 5 | **Market / TAM-SAM-SOM** | Size of the prize, **bottom-up** | The 1%-of-a-huge-market trap (see §1.3) |
| 6 | **Product** | How it works / demo / screenshots | Telling not showing |
| 7 | **Traction** | Evidence customers want it — revenue, users, growth, LOIs, waitlist | Vanity metrics; no trend line |
| 8 | **Business model** | How you make money (or, for impact, the funding model) | Hand-waving on pricing |
| 9 | **Go-to-market** | How you acquire customers and at what cost | "We'll do marketing" |
| 10 | **Competition** | Who else solves this + your wedge (a 2x2 or feature matrix) | "We have no competition" (red flag) |
| 11 | **Team** | Why *you* win — founder-market fit | Generic résumés; no relevant edge |
| 12 | **The ask** | Amount, use of funds, and the milestones it buys | Vague "raising a round"; no use-of-funds |

Put **deep financials, detailed metrics methodology, and extra team/tech in an appendix**, not the
main flow [Sequoia uses an appendix].

### 1.2 The narrative arc (order matters as much as content)

YC's structural rule: lead with **"what you do,"** end with **"the ask,"** and order the middle
**most-impressive-to-least.** If you have strong traction, move it right after the "what you do"
slide rather than burying it at slide 7 [YC; Storydoc]. The deck is a story — problem creates
tension, solution + why-now resolve it, traction proves it's real, the ask is the call to action.
The classic spine is **Problem → Solution → Why Now → Proof (traction/market) → The Plan
(model/GTM/team) → The Ask.**

### 1.3 Common mistakes

- **Top-down market sizing / the 1% fallacy.** "We only need 1% of a $100B market" is the single
  most credibility-destroying slide. Investors want **bottom-up** TAM/SAM/SOM built from real
  customer counts × price × conversion. Don't mix top-down and bottom-up on one slide [TechCrunch;
  Qubit; spectup].
- **Projections that contradict the rest of the deck.** If your model shows 1,000 customers by
  year 2 but your market slide implies a SAM of 500, you've lost the room. The deck, the model, and
  the unit economics must tell **one consistent story** [pitchdeckguide].
- **Vanity metrics over momentum.** Investors want believable numbers that align with your product
  and GTM (MRR, CAC, LTV, burn) and a **trend**, not a single big number [spectup; goingvc].
- **No why-now, no team edge, no use-of-funds.** Each is a standard reason to pass.

### 1.4 The deck-vs-memo question

A live debate as of 2025-26. A **pitch deck** is visual and brief — built to *capture attention and
get the meeting*. An **investment memo** is narrative and rigorous — it "does the work when you're
not in the room," letting investors evaluate asynchronously and *close* the deal [Visible.vc;
Qubit]. They are increasingly used **together**: Rippling famously raised a $45M Series A with **no
formal pitch deck**, leading with an investor memo — but *accompanied by 46 slides* of metrics and
methodology footnotes [Rippling]. Guidance: build the **deck to open**, and a **memo (or a tight
data room) to close**; for a seed raise, a strong deck + a one-to-two-page memo is plenty. Memo
prose craft → `executive-comms`.

---

## 2. The 3-statement integrated financial model

This is the **handoff from the lean projection** (the sibling skill's single-tab P&L + cash-flow
sketch) to a *real* model where the three core statements are **linked** and move together. A
3-statement model combines the **income statement (P&L)**, the **cash-flow statement**, and the
**balance sheet** into one dynamic forecast driven by assumptions [CFI; Vena].

### 2.1 The three statements and how they link

Build order is conventionally **income statement → balance sheet → cash-flow statement**, then wire
the links [CFI; Wall Street Prep]:

1. **Income statement → cash-flow statement.** Net income (the P&L "bottom line") is the **top line
   of the cash-flow statement's** Cash-from-Operations section [CFI].
2. **Income statement → balance sheet.** Net income *not* paid out as dividends flows into
   **retained earnings** on the balance sheet (Retained Earnings_end = Retained Earnings_begin +
   Net Income − Dividends) [Wall Street Prep].
3. **Cash-flow statement → balance sheet.** Changes in **working-capital accounts** (A/R, inventory,
   A/P) drive operating cash flow; the **ending cash balance** from the bottom of the cash-flow
   statement becomes the **cash line** at the top of the balance sheet [CFI; FE Training].
4. **The test:** the balance sheet must **balance** (Assets = Liabilities + Equity) in every period.
   If it doesn't, a link is wrong. [UNVERIFIED — exact treatment varies by template] Interest
   expense computed on average debt creates **circularity** (interest → net income → cash → debt
   paydown → interest); handle with an iterative-calc toggle or a "circularity switch."

```
  ┌──────────────────┐   net income    ┌────────────────────┐
  │ Income Statement │ ───────────────▶ │ Cash-Flow Statement │
  └──────────────────┘                  └────────────────────┘
        │  retained earnings                     │ ending cash
        ▼                                         ▼
  ┌─────────────────────────────────────────────────────────┐
  │           Balance Sheet  (Assets = Liab + Equity)        │
  └─────────────────────────────────────────────────────────┘
```

### 2.2 Driver-based assumptions

A credible model is **driver-based**: outputs flow from a small set of explicit assumptions you can
defend, not hard-coded numbers. Typical drivers — new customers/month, churn, ACV/price, sales-rep
productivity, gross margin %, headcount plan, A/R and A/P days, capex. Keep all assumptions on **one
clearly labeled tab** so an investor (or you, in six months) can flex one number and watch the
statements respond [CFI; Vena]. This is also where the sibling skill's **unit economics** (CAC,
LTV, contribution margin) plug in as drivers — see `venture-small-business-planning`.

### 2.3 Scenario & sensitivity analysis

- **Scenarios:** model **Base / Upside / Downside** by swapping assumption *sets* (e.g., growth
  rate and churn together). Investors expect at least a base and a conservative case.
- **Sensitivity:** flex **one** driver across a range to see which inputs the outcome is most
  sensitive to (a one- or two-variable data table). Identifies what to de-risk first.
- **The number that matters most early:** **runway / cash-out date** (months of cash at current
  burn). The model's job at seed is less "predict year-5 revenue" (no one believes it) and more
  "show that the raise buys enough runway to hit the next milestone." Tie the **use-of-funds** on
  the deck's ask slide directly to this.

> Reality check: a seed-stage model is a **statement of assumptions and ambition**, not a forecast.
> Investors read it to judge whether you *think clearly about drivers*, not to trust the year-5
> number.

---

## 3. Cap table & dilution

The **cap table** records who owns what (founders, option pool, investors) in shares and
percentages. Fundraising changes it; understanding the math protects the founder's stake.

### 3.1 Pre-money, post-money, and price per share

- **Pre-money valuation** = the company's agreed value *before* new money goes in.
- **Post-money valuation** = **pre-money + new investment**.
- **Investor ownership %** = **investment ÷ post-money** [Carta; cake equity — two sources agree].
- **Price per share** = pre-money ÷ fully-diluted shares *before* the round.
- **New shares issued** = investment ÷ price per share.

A new round **dilutes** every existing holder's percentage (their share *count* is unchanged; the
denominator grows).

### 3.2 The option pool (and the "option-pool shuffle")

Investors usually require an **employee option pool** (seed pools commonly **~11-13%**, median
~11.8% [dowhatmatter]). The trap, the **"option-pool shuffle"**: investors demand the new/expanded
pool be carved out of the **pre-money** — so the shares are created *before* the round, which
**dilutes founders, not the incoming investor**, while the investor's target % stays fixed
[Venture Hacks; Kruze; Alder VC].

> **Worked shuffle example** [synthesized from Alder VC / dowhatmatter]: agree to a **20% pool on a
> $5M post-money** with $1M allocated to the pool. The investor fixes their stake at 20% ($1M). The
> $1M pool comes out of *your* side. Founders go from **80% → 60%** — you paid for the entire pool.
> Negotiating lever: size the pool to the *actual* next-18-months hiring plan, not a round number,
> and push for as much of it as possible to come post-money.

### 3.3 SAFE (post-money) vs convertible note vs priced round

As of 2025 the early-stage default is overwhelmingly the **post-money SAFE**: SAFEs were **~90% of
all pre-seed rounds** on Carta in Q1 2025 and **64% of seed rounds** (vs 27% priced, 10% notes);
**~87% of SAFEs** issued were **post-money** [Carta State of Pre-Seed Q3 2025].

| | **Post-money SAFE** | **Convertible note** | **Priced round** |
|---|---|---|---|
| Legal nature | Not debt; converts on a trigger | **Debt** — interest + maturity | Equity sold now at a set price |
| Interest | None | **~4-8%/yr** (median ~7%), compounds & converts | n/a |
| Maturity | None — outstanding until a trigger | **~18-36 months**; repayable if no conversion | n/a |
| Valuation | **Valuation cap** (± discount) | Cap and/or discount | Negotiated pre-money |
| Ownership clarity | **Fixed at signing**: investment ÷ post-money cap = % | Unknown until conversion | Known immediately |
| Cost / speed | Cheap, fast, few docs | Moderate | Expensive, slow (full term sheet, legal) |

*Sources: Carta; CRV; cake equity; Qubit — cross-checked across ≥2.*

**Post-money SAFE mechanics** [Carta; cake equity]:
- **Ownership = investment ÷ post-money valuation cap.** $500K on a $10M post-money cap = **5%**,
  fixed at signing.
- The **post-money cap includes** the SAFE itself, founder shares, and any pre-conversion option
  pool — but **excludes the future priced round** and (critically) **other SAFEs don't dilute each
  other's %**. Each *new* post-money SAFE dilutes founders and existing holders directly.
- **Discount** (e.g., 20%) and/or **MFN** clauses may also apply.
- **Founder cost of clarity:** post-money SAFEs are **more dilutive to founders in later rounds**
  than the old pre-money SAFEs, precisely because they fix the investor's % and push all subsequent
  dilution onto the founders [Carta]. Top founders in 2025 kept median seed dilution to ~19%, often
  structuring SAFEs to stay under an ~18% threshold [Carta].

### 3.4 Worked dilution example (priced seed round)

A clean priced-round waterfall (round numbers chosen so every figure reconciles; the *method*
follows cake equity / CRV):

| Holder | Shares | % (pre) | | % (post) |
|---|---:|---:|:-:|---:|
| Founder A | 3,500,000 | 43.8% | → | **35.0%** |
| Founder B | 3,500,000 | 43.8% | → | **35.0%** |
| Option pool | 1,000,000 | 12.5% | → | **10.0%** |
| **Pre-round fully diluted** | **8,000,000** | 100% | | |
| New seed investor | 2,000,000 | (new) | → | **20.0%** |
| **Post-round total** | **10,000,000** | | | **100%** |

**The math:** raise **$2M** at **$8M pre-money** → **$10M post-money**. Price/share = $8M ÷
8,000,000 = **$1.00**. New shares = $2M ÷ $1.00 = **2,000,000**. Post-round total = 8M + 2M =
**10,000,000**. Investor % reconciles **both ways**: 2M ÷ 10M shares = **20%** = $2M ÷ $10M
post-money. Each founder goes 43.8% → 35.0% — **8.8 points**, a **20% erosion** of their stake.

> Founders typically expect **~15-25% dilution per priced round** [cake equity / CRV]. Across
> seed + Series A + B, plan for founders to be well under 50% by Series B — model the *full stack*,
> not one round in isolation.

---

## 4. The grant-writing lifecycle (for the nonprofit)

For the **organ-donation nonprofit**, "fundraising" largely means **grants** plus donations.
**Prerequisite:** most foundation and government grants require **501(c)(3)** status — and NC
**charitable-solicitation registration** before you fundraise. Confirm both first →
`venture-nc-nonprofit-formation`. The **persuasive writing itself** (compelling need statement,
case for support, exec summary) → `executive-comms`. This section is the **lifecycle and the
artifacts**.

### 4.1 Prospect research (find the right funders)

Before writing a word, find funders whose **mission, geography, and grant size** fit yours. A
mismatched application is wasted effort. Tools: **Candid / Foundation Directory**, GrantStation,
Instrumentl, federal **Grants.gov**, and for an organ-donation cause, **Donate Life** affiliates,
health-system community-benefit funds, and state/HRSA programs. Score prospects on **alignment**
(do they fund this cause?), **fit** (right program area + dollar range?), and **accessibility** (do
they accept unsolicited requests, or invite-only?).

### 4.2 Letter of inquiry (LOI)

Many funders use a **two-stage process** starting with an LOI (a.k.a. letter of intent / interest /
concept paper) — a brief, formal pitch requesting **permission to submit a full proposal** [Candid;
GrantStation]. Keep it **≤ 3 pages**; a good LOI is often *harder* to write than the full proposal
because it must compress everything [Candid; Instrumentl]. Include: organization background +
mission + capacity; the **need/problem**; the **proposed solution**; goals, objectives, activities;
evaluation approach; and a budget summary. Candid recommends building a **logic model first** so the
LOI's elements line up cleanly [Candid].

### 4.3 Full proposal

If the LOI shows alignment, the funder invites a **complete proposal** [Candid]. The four common
grant-proposal documents are the **proposal narrative, the logic model, the project budget, and
supporting attachments** [Candid]. The narrative typically covers: statement of need (data-backed),
project description (goals → SMART objectives → activities), the logic model, evaluation plan,
organizational capacity, sustainability beyond the grant, and the budget + justification. Narrative
**writing craft** → `executive-comms`.

### 4.4 The logic model

The funder-standard planning artifact, formalized by the **W.K. Kellogg Foundation** and required
by most US federal, state, and private funders [Kellogg; Sopact]. A one-page, left-to-right chain:

```
  Inputs/Resources → Activities → Outputs → Outcomes (short → medium → long) → Impact
```

- **Inputs** — resources you commit (staff, funds, volunteers, partners).
- **Activities** — what the program *does* (e.g., DMV registry drives, faith-community workshops).
- **Outputs** — direct, *countable* products (e.g., # events held, # people reached).
- **Outcomes** — the *change* produced, sequenced short → medium → long (e.g., # new registered
  donors → shifted attitudes → higher community registration rate).
- **Impact** — the long-run mission goal.

Distinguish from a **theory of change**, which maps the *causal mechanisms and assumptions* behind
why activities produce outcomes; the logic model is the operational/accountability view, the ToC is
the evaluation-design view [Sopact; La Piana].

### 4.5 The project budget

Tie every dollar to the project; the **budget and narrative must tell the same story** or reviewers
flag it [fundsforNGOs; technicalwriterhq]. Structure:

- **Direct costs** — attributable to the project: project staff salaries, materials, travel, events.
- **Indirect costs (overhead)** — org-wide costs not tied to one project: admin salaries, rent,
  utilities. Funders cap these, commonly **10-20% of direct costs**; the federal **de minimis rate
  is 10% of Modified Total Direct Costs (MTDC)** as of Oct 1, 2023, and indirect is allowed only
  with a federally approved rate agreement [Uniform Guidance 2 CFR 200; National Council of
  Nonprofits].
- **Budget justification/narrative** — for each line: who, what it covers, why it's needed, and how
  the amount was calculated [fundsforNGOs; RSP Wisconsin].

### 4.6 Reporting & stewardship

The award is the *start* of the relationship, not the end. Funders expect **progress and final
reports** tying spending and results back to the **outcomes in your logic model**, organized to
their reporting format (federal grants follow **2 CFR 200** categorization, which many private
funders adopt) [fundsforNGOs; OJP]. Good stewardship — accurate reports, hitting stated outcomes,
proactive communication — is what renews and grows funding. Report *writing* craft →
`executive-comms`.

---

## 5. Cross-skill map

| Need | Skill |
|------|-------|
| Business Model Canvas, break-even, basic unit economics, **lean** projection | `venture-small-business-planning` |
| The **prose** of a grant narrative, exec summary, founder/shareholder letter, memo, deck copy | `executive-comms` |
| 501(c)(3) recognition, NC charitable-solicitation registration, Form 990 (grant **eligibility**) | `venture-nc-nonprofit-formation` |
| For-profit NC entity choice & tax (LLC/C-/S-corp) — the entity that *holds* the cap table | `venture-nc-business-formation-tax` |
| The actual spreadsheet (.xlsx, formulas, cells) | `document-formats` / xlsx |
| General prose voice & editing | `writing-expert` |

---

## Sources

- Y Combinator, pitch deck design (legibility/simplicity/obviousness; order; the ask): ycombinator.com/library (via Leland, Storydoc, SaaStr summaries)
- Sequoia Capital, pitch deck template (10 slides + appendix; company purpose; why now): slidebean.com/templates/sequoia-capital-pitch-deck; uvic.ca pitch-deck-template
- Storydoc, "How to Write a YC Pitch Deck": storydoc.com/blog/y-combinator-pitch-deck-examples
- TechCrunch, "5 critical pitch deck slides most founders get wrong" (top-down TAM trap)
- Qubit Capital; spectup; goingvc; pitchdeckguide, on TAM/SAM/SOM, market-sizing & traction mistakes
- Visible.vc, "Investment Memo: Template, Examples"; Rippling, "Series A Pitch Deck and Memo" (deck-vs-memo)
- Corporate Finance Institute, "What is a 3-Statement Model?": corporatefinanceinstitute.com/resources/financial-modeling/3-statement-model/
- Wall Street Prep, "How are the Three Financial Statements Linked?"; FE Training, linking statements; Vena Solutions, 3-statement model
- Carta, "Pre-money vs post-money SAFEs," "What is a SAFE?," "Priced Rounds," "Convertible Securities," "State of Pre-Seed Q3 2025": carta.com/learn & carta.com/data
- CRV, "SAFE vs Convertible Note"; cake equity, "Equity Dilution" & "SAFE vs Convertible Note" (worked example, formulas)
- Venture Hacks, "The Option Pool Shuffle"; Kruze Consulting; Alder VC, on option-pool dilution math
- dowhatmatter.com, seed-stage option-pool benchmarks (~11.8% median)
- Candid, "Four common grant proposal documents," LOI knowledge base: candid.org; learning.candid.org
- Instrumentl; GrantStation, letter of inquiry guidance
- W.K. Kellogg Foundation, Logic Model Development Guide; Sopact; La Piana, on logic model vs theory of change
- fundsforNGOs; Professional Grant Writers; technicalwriterhq; RSP Wisconsin, on grant budgets, direct/indirect costs, justification
- U.S. Uniform Guidance **2 CFR 200** (de minimis 10% MTDC indirect rate, eff. Oct 1 2023); National Council of Nonprofits; OJP Grants 101, on reporting & cost rules

## Disclaimer

This skill is **educational information, not financial, legal, tax, or investment advice**, and does
not create any professional relationship. Fundraising instruments and their consequences are
high-stakes and jurisdiction- and fact-specific:

- **SAFE, convertible-note, and priced-round terms, cap-table entries, and dilution outcomes** must
  be reviewed with a **startup/securities attorney** before you sign anything; the worked examples
  here are illustrative and simplified (they omit liquidation preferences, anti-dilution, SAFE
  stacking interactions, pro-rata rights, and conversion edge cases).
- **Financial-model outputs** are only as good as their assumptions; have a **CPA** review tax
  treatment, accounting, and any projection shared with investors.
- **Securities laws** govern who you may solicit and how (e.g., accredited-investor rules,
  Regulation D, state blue-sky law) — get counsel before raising from anyone.
- **Grant eligibility, indirect-cost rules, and tax-exempt compliance** turn on current IRS, NC,
  and funder requirements that **change**; verify against the funder's guidelines, IRS, and the NC
  Secretary of State, and see `venture-nc-nonprofit-formation`.

Market data (Carta SAFE/dilution percentages, option-pool benchmarks, indirect-cost rates) reflects
**2024-2026** snapshots and shifts over time — re-verify before relying on a specific figure. Items
marked **[UNVERIFIED]** were not confirmed against a primary source and should be checked
independently.
