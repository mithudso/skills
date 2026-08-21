<!-- hub-reference-banner -->
> **Reference file — part of the `venture-business` hub.** Formerly the standalone `venture-small-business-planning` skill.
> Sibling topics in this family are now reference files under the hubs (`venture-business`, `venture-nonprofit-cause`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: venture-small-business-planning
description: >-
  Practical business-planning and financial-modeling guide for a solo/lean founder — business model
  design (Business Model Canvas vs Lean Canvas), idea/market validation (customer discovery, problem-solution
  and product-market fit, TAM/SAM/SOM, competitor analysis), the business plan document (lean one-pager vs
  traditional plan, and when a full plan is actually required), startup financial modeling (one-time costs,
  P&L/income statement, cash-flow projection, break-even, unit economics — CAC/LTV/contribution margin,
  pricing), funding options (bootstrapping, revenue, SBA loans/microloans, grants, friends & family,
  angel/VC) with NC small-business resources (SBA NC district, NC SBTDC, SCORE, NC IDEA grants), plus KPIs,
  milestones, a simple operating plan, and risk/assumptions discipline. General to any small business or
  startup, with notes on where a nonprofit differs (mission vs profit; grants/donations vs revenue; fiscal
  sponsorship). TRIGGER: questions about a business model, business plan, financial projections/model, unit
  economics, break-even, pricing, market sizing, validating an idea, or how/where to get funding for a small
  business or startup. SKIP: NC entity choice and tax filing (LLC/C-corp/S-corp, state tax, DBA) →
  venture-nc-business-formation-tax; nonprofit 501(c)(3)/charitable-solicitation/Form 990 legal and tax →
  venture-nc-nonprofit-formation; marketing strategy, local SEO, ad copy → venture-marketing-strategy-local-seo
  / venture-cause-nonprofit-marketing; spreadsheet/Excel mechanics and formulas → the pro library's xlsx /
  document-formats skill.
category: personal-venture
tags: [venture, business-planning, financial-modeling, startup, funding]
whenToUse:
  - Choosing or filling out a Business Model Canvas or Lean Canvas
  - Validating a business idea through customer discovery and market sizing
  - Deciding whether you need a lean one-pager or a full traditional business plan
  - Building a startup financial model (costs, P&L, cash-flow, break-even)
  - Working out unit economics (CAC, LTV, contribution margin) or setting prices
  - Choosing a funding path (bootstrap, loan, grant, friends & family, angel/VC)
  - Finding North Carolina small-business funding and free advising resources
  - Setting KPIs, milestones, and a simple operating plan for a lean venture
triggers:
  - business model canvas / lean canvas
  - business plan
  - financial projections / financial model
  - unit economics / CAC / LTV / break-even
  - market sizing / TAM SAM SOM / validate idea
  - startup funding / SBA loan / grant / NC IDEA
---

# Small-Business Planning & Financial Modeling

A founder-grade reference for turning an idea into a tested business model, a credible financial model,
and a funding plan you can execute. Written for a solo/lean founder, general to any small business or
startup. **Nonprofit notes** are flagged inline because the venture in focus is a nonprofit
(organ-donation awareness) and the logic differs in a few places — mission-vs-profit, funding via
grants/donations rather than revenue, and the fiscal-sponsorship option.

Boundaries: this is the *planning and money-modeling* layer. **Legal entity + tax** (LLC vs C-corp vs
S-corp, NC state tax, DBA) → `venture-nc-business-formation-tax`. **501(c)(3)/charitable-solicitation/
Form 990** mechanics → `venture-nc-nonprofit-formation`. **Marketing/SEO/ad copy** →
`venture-marketing-strategy-local-seo` (for-profit) or `venture-cause-nonprofit-marketing` (cause).
**Spreadsheet mechanics** (the actual .xlsx, formulas) → the xlsx / `document-formats` skill — this skill
says *what* to model, not which cell to type it in. **Deep fundraising artifacts** — the investor pitch
deck, an integrated 3-statement model, the cap table & dilution math, the SAFE/note/priced-round
decision, and the grant lifecycle — step up to `venture-startup-fundraising-deck-model`.

Work in roughly this order, but loop freely: **model → validate → plan → finance → operate → revisit
risks.** Early on, spend more time talking to customers than formatting documents.

---

## 1. Business model design: Business Model Canvas vs Lean Canvas

A *business model* is simply how your venture creates, delivers, and captures value. Sketch it on one
page before you write any plan. Two one-page tools dominate; pick by stage.

### Business Model Canvas (BMC) — the 9 blocks

Created by Alexander Osterwalder (Strategyzer). Best for an **existing business, a known model, or
corporate/structured planning** where partners and operations already matter. The nine blocks
(strategyzer.com, as of 2026-06):

1. **Customer Segments** — who you serve (mass market, niche, segmented, multi-sided).
2. **Value Propositions** — the bundle of products/services that solve a customer problem.
3. **Channels** — how you reach and deliver to customers (awareness → purchase → after-sales).
4. **Customer Relationships** — how you get, keep, and grow customers (self-serve, personal, community).
5. **Revenue Streams** — how you capture value (sales, subscription, usage, licensing, ads).
6. **Key Resources** — the critical assets (people, IP, capital, physical).
7. **Key Activities** — the most important things you must do well (production, problem-solving, platform).
8. **Key Partnerships** — suppliers and allies who supply resources or activities.
9. **Cost Structure** — what it costs to operate (cost-driven vs value-driven; fixed vs variable).

Right column (segments, value prop, channels, relationships, revenue) = the *customer/desirability*
side; left column (resources, activities, partners, costs) = the *infrastructure/viability* side.

### Lean Canvas — the startup variant

Adapted by Ash Maurya (LeanStack) **for early-stage startups searching for a model under high
uncertainty.** It keeps the one-page format but swaps four BMC blocks for ones that surface *risk*:
replaces **Key Partners + Key Activities → Problem + Solution**, and **Key Resources → Unfair
Advantage**, and **Customer Relationships → Key Metrics** (leanstack.com / Maurya, as of 2026-06). The
nine blocks: **Problem, Customer Segments, Unique Value Proposition, Solution, Channels, Revenue
Streams, Cost Structure, Key Metrics, Unfair Advantage.**

### Which to use

- **Brand-new idea, no proven model, testing whether a problem is real → Lean Canvas.** It forces you
  to name the problem and the riskiest assumptions first. This is almost always where a solo founder
  starts.
- **Established operation, or you need to map partners/operations/infrastructure → BMC.**
- They are not mutually exclusive: start on a Lean Canvas, graduate to a BMC once the model is proven.
- **Nonprofit note:** neither tool fits cleanly because "Revenue Streams" assumes paying customers.
  Adapt: treat the **beneficiary** (who you help) and the **funder/donor** as two distinct "segments" —
  usually different people — and rename Revenue Streams to **Funding Streams** (grants, donations,
  earned/program revenue). A *Mission/Impact* statement stands in for the profit motive (the "Mission
  Model Canvas" variant exists for this).

---

## 2. Idea & market validation

Validation is the discipline of replacing assumptions with evidence *before* you spend money. Steve
Blank's core point: a startup is not a small version of a big company — a big company *executes* a
known model; a startup *searches* for one (Blank, *The Four Steps to the Epiphany*; steveblank.com, as
of 2026-06).

### Customer discovery & interviews

- **Goal:** confirm a real, painful problem exists before you build. Test one hypothesis per round:
  *does this problem exist, and is it painful enough that people already try to solve it?*
- **The Mom Test (Rob Fitzpatrick):** ask about **specific past behavior, not hypothetical future
  behavior**. "Walk me through the last time you dealt with X" and "What do you use today?" reveal truth;
  "Would you buy this?" invites polite lies. Don't pitch; listen. Aim for ~15-30 conversations per
  segment early on, and look for patterns rather than validation of your idea.

### Problem-solution fit → product-market fit

Two milestones, in order:

1. **Problem-solution fit** — you have evidence the problem is real *and* your proposed solution
   plausibly solves it (often before you've built much). This is the end of customer discovery.
2. **Product-market fit (PMF)** — you've built something a defined market actually wants and uses, with
   evidence in retention and demand (people keep coming back; turning them off would upset them). PMF
   is the real goal; everything before it is search.

### Market sizing: TAM / SAM / SOM

A nested estimate of opportunity (forumvc.com, qubit.capital, as of 2026-06):

- **TAM (Total Addressable Market)** — total revenue if you captured 100% of demand. The vision.
- **SAM (Serviceable Available Market)** — the slice you can actually serve given geography, segment,
  channel. A subset of TAM.
- **SOM (Serviceable Obtainable Market)** — what you can realistically win in the near term (years
  1-3). A subset of SAM.

Two methods — **do both and cross-check:**

- **Top-down:** start from a published industry figure and filter down (e.g., $70B CRM market × 28%
  small-business share ≈ $19.5B TAM). Fast, but easy to inflate.
- **Bottom-up:** `Market size = (annual revenue per customer) × (number of potential customers)`. More
  credible; build it from real unit prices and a countable customer base.

**Discipline:** for an early startup, claiming more than ~1-5% of SAM in the first few years reads as
naive. Use top-down for vision, bottom-up for the credible path.

### Competitor analysis

List direct competitors (same solution), indirect (different solution, same problem), and the status
quo / "do nothing" — often your biggest competitor. For each: offering, price, strengths, weaknesses,
and the gap you exploit. Your **unfair advantage** (Lean Canvas) is what competitors can't easily copy —
proprietary data, network, brand, exclusive relationships, hard-won expertise.

**Nonprofit note:** market sizing becomes *need sizing* — how many people have the problem you address
(e.g., people on the transplant waitlist; unregistered eligible donors) — and *funding sizing* — the
pool of grants/donors realistically reachable. "Competitors" become other orgs serving the same cause;
collaboration is often smarter than competition. See `venture-organ-donation-system` for the donation
landscape and `venture-cause-nonprofit-marketing` for donor-journey design.

---

## 3. The business plan document

A business plan is the written articulation of the validated model plus how you'll execute and what
you need. The SBA recognizes two formats (sba.gov "Write your business plan", as of 2026-06).

### Lean (one-page) plan

A short, high-level summary — the SBA's lean format mirrors the BMC's nine blocks (key partnerships,
key activities, key resources, value proposition, customer relationships, customer segments, channels,
cost structure, revenue streams). Takes ~an hour, revised often. **Use it** as your internal roadmap,
for a simple model, or when you'll iterate frequently. For most solo founders this is the right
starting document.

### Traditional plan

A detailed 15-25 page document. **Most lenders and investors expect this format.** SBA's nine sections
and what each is *for*:

1. **Executive Summary** — what the company is and why it will succeed (write last; it's the hook).
2. **Company Description** — the problem you solve, who you serve, your edge.
3. **Market Analysis** — industry outlook, target market, competitive landscape (your §2 work).
4. **Organization & Management** — legal structure and who runs it (org chart). *Entity choice itself →
   `venture-nc-business-formation-tax` / `venture-nc-nonprofit-formation`.*
5. **Service or Product Line** — what you sell, the customer benefit, IP plans.
6. **Marketing & Sales** — how you'll attract and keep customers (*strategy depth →
   `venture-marketing-strategy-local-seo`*).
7. **Funding Request** — how much you need over ~5 years and exactly what for.
8. **Financial Projections** — historical (if any) + ~5-year forecast: P&L, cash flow, balance sheet.
9. **Appendix** — resumes, licenses, permits, letters of intent, detail tables.

### When you actually need a full plan

- **Applying for a loan or grant** — banks, SBA lenders, and most grantmakers require a traditional plan
  with detailed financials. The single most common forcing function.
- **Raising from investors** — they expect a plan and/or a pitch deck plus a model.
- **The business is complex, capital-intensive, or has co-founders** who need alignment.
- Otherwise a lean one-pager plus a solid financial model is enough to operate. Don't let a 50-page
  document substitute for talking to customers.

**Nonprofit note:** for grant applications, the "business plan" is often reframed as a **program/impact
plan** plus a **budget** and sometimes a logic model (inputs → activities → outputs → outcomes). Funders
care about mission, beneficiaries reached, and outcomes — not profit. SCORE and the SBTDC will still
help nonprofits with planning.

---

## 4. Financial modeling for a startup

A financial model is a forward projection of money in and out. For a solo founder it can live in one
spreadsheet with a few linked tabs. Build it from assumptions you can defend, and label every
assumption. (Spreadsheet mechanics → document-formats (references/xlsx.md) / `document-formats`; this section is *what* to model.)

### a) Startup / one-time costs

Separate **one-time** costs from **recurring (monthly)** costs (sba.gov "Calculate your startup costs",
as of 2026-06):

- **One-time:** equipment, deposits, initial inventory, incorporation/legal/licenses, logo/branding,
  website build, initial buildout.
- **Recurring (monthly):** rent, utilities, salaries/contractors, software subscriptions, ongoing
  marketing, insurance, loan payments, COGS.

A practical sizing rule:
`Total startup capital ≈ one-time costs + (monthly recurring × months to break-even) + a working-capital
reserve.` Plan **at least 12 months** of runway; lenders/investors compare projected costs to projected
revenue to judge viability, so present this cleanly.

### b) Income statement (P&L)

Shows **profitability over a period** on an accrual basis (revenue when *earned*, expenses when
*incurred*):

```
Revenue
− COGS (direct costs of delivering the product/service)
= Gross profit            (Gross margin % = gross profit / revenue)
− Operating expenses      (rent, salaries, marketing, software, G&A)
= Operating profit (EBIT)
− Interest & taxes
= Net profit (the "bottom line")
```

### c) Cash-flow projection — and why cash ≠ profit

The cash-flow statement shows **actual money moving in and out** — your liquidity and survival. **A
profitable business can still run out of cash** (HBS Online; ramp.com, as of 2026-06) because:

- **Accrual timing:** you book revenue when earned, but cash arrives later (accounts receivable / net-30
  invoices).
- **Balance-sheet movements never hit the P&L:** inventory purchases, equipment (capex), loan principal,
  and sales tax all consume cash without reducing reported profit.
- **Growth burns cash:** every new customer often costs money to acquire *before* they pay.

So model cash month-by-month, track the **closing cash balance**, and never let it go negative. Cash is
more important than profit in the short run — running out of cash is the most common way small
businesses die.

### d) Break-even analysis

The point where total revenue covers total costs (profit = 0).

`Break-even units = Fixed costs ÷ (Price per unit − Variable cost per unit)`

The denominator is the **per-unit contribution margin**. Break-even in dollars =
`Fixed costs ÷ contribution-margin %`. Knowing your break-even tells you the minimum volume to survive
and how price changes move that target.

### e) Unit economics — does one sale make money?

Strip the business down to a single customer/unit:

- **Contribution margin** = revenue per unit − variable costs per unit (the dollars left to cover fixed
  costs and profit).
- **CAC (Customer Acquisition Cost)** = total sales & marketing spend ÷ new customers acquired.
- **LTV (Lifetime Value)** = the **gross profit** (not revenue) a customer generates over their
  lifetime. A common form: `LTV = (avg revenue per customer × gross margin %) ÷ churn rate`.
- **LTV:CAC ratio** — health benchmark **~3:1 or higher** (you earn $3 of lifetime value per $1 spent
  acquiring). 1:1-2:1 means thin or no margin (ltvcacbook.com / Kruze, as of 2026-06).
- **CAC payback period** = `CAC ÷ (monthly revenue per customer × gross margin %)` — months to recover
  acquisition cost. Shorter is better; under ~12 months is comfortable for many small businesses.

If unit economics are underwater, growth makes losses *bigger*, not smaller. Fix the unit before you
scale spend.

### f) Pricing approaches

Four common approaches (quickbooks.intuit.com; Salesforce; bdc.ca, as of 2026-06):

- **Cost-plus:** total unit cost + a markup (cost $50 + 20% = $60). Simple; ignores customer value and
  competitors. A floor, not a strategy.
- **Value-based:** price to the value the customer perceives. Usually the most profitable; needs deep
  customer understanding and real differentiation.
- **Competitive:** price relative to rivals. Common where buyers compare on price; requires ongoing
  monitoring.
- **Dynamic:** adjust by demand/time. More relevant once you have volume and data.

Make sure price minus variable cost yields a contribution margin that lets you hit break-even at a
plausible volume. Most founders **underprice** — anchor on value, not just cost.

**Nonprofit note:** much "revenue" is donations/grants, which aren't priced like a product. But any
**earned-revenue** activity (event tickets, merchandise, fee-for-service programs) still needs unit
economics and pricing so it doesn't lose money. Model program *cost per beneficiary served* as your
unit metric, and donor acquisition cost as a CAC analogue.

---

### Worked mini financial-projection example

A solo founder launches a small subscription service. Illustrative numbers (round figures for clarity):

**Assumptions:** price $30/customer/month; variable cost to serve $6/customer/month → contribution
margin $24/customer (80% gross margin). Fixed monthly costs (rent, software, the founder's modest draw):
$3,000. One-time startup costs (laptop, website, incorporation, branding): $5,000. CAC $60/customer.

- **Break-even volume:** $3,000 ÷ $24 = **125 active customers** to cover monthly fixed costs.
- **Unit economics:** if average customer stays 20 months, LTV = $24 × 20 = $480. LTV:CAC =
  $480 ÷ $60 = **8:1** (healthy). CAC payback = $60 ÷ ($30 × 80%) = $60 ÷ $24 = **2.5 months**.
- **Cash vs profit, Month 1:** sign 40 customers. Revenue $1,200; variable $240; fixed $3,000;
  acquisition 40 × $60 = $2,400. **P&L:** 1,200 − 240 − 3,000 − 2,400 = **−$3,840** (a normal early
  loss). **Cash:** start with $20,000 raised, spend $5,000 one-time + $3,840 operating ≈ closing cash
  **$11,160**. At this burn the runway is the constraint, not profitability — watch the balance monthly.
- **Path:** at ~30 net new customers/month you cross 125 (break-even) around month 5-6 *if* churn stays
  low. The model's job is to show whether you reach break-even before cash runs out — and to flag the
  assumptions (churn, CAC, net adds) that decide the answer.

A model is not a prediction; its value is the *conversation about assumptions*. [Figures illustrative,
not benchmarks.]

---

## 5. Funding options + North Carolina resources

Match the funding source to your stage, your risk, and whether you can/should give up ownership.

### The options (and when each fits)

- **Bootstrapping / self-funding** — personal savings, reinvested revenue. Keeps 100% ownership and
  control; constrains speed. The default for most lean founders; appropriate whenever the business can
  grow on its own cash.
- **Revenue (customer-funded)** — let paying customers fund growth. The healthiest "funding" there is;
  pursue early via pre-sales, deposits, or a paid pilot.
- **Friends & family** — early, small, relationship-based. Cheap and fast but risks personal
  relationships; document terms in writing (loan vs equity) and be honest about the risk of total loss.
- **SBA loans & microloans** — government-*backed* loans issued by lenders/intermediaries. Good for
  founders who want capital **without giving up equity** and can service debt. Specifics below.
- **Grants** — non-dilutive (no equity, no repayment) but competitive and often restricted to a purpose.
  Great when you qualify; never count on them as your only plan. NC programs below.
- **Angel investors** — wealthy individuals investing early for equity. Fit when you have a high-growth,
  scalable model and need capital + mentorship beyond what debt provides; you give up ownership and gain
  a stakeholder.
- **Venture capital (VC)** — institutional equity for **high-growth, large-market** companies that can
  return many times the investment. Fits a small minority of startups; most small businesses are not VC
  candidates and shouldn't optimize for it. Expect dilution, a board, and pressure to scale fast.

A rough progression for a high-growth startup: bootstrap/F&F → grants/microloan → angel → VC. A
lifestyle or local small business may stop at bootstrap + a bank/SBA loan and never raise equity — which
is completely fine.

### SBA loan specifics (verify current terms with SBA/lender)

All figures from sba.gov, as of 2026-06 — **confirm before relying on them; terms change:**

- **SBA Microloan** — up to **$50,000** (avg ~**$13,000**), max term **7 years**, interest generally
  **8-13%**, via nonprofit community-based **intermediary lenders**. For working capital, inventory,
  supplies, furniture/fixtures, machinery, equipment; **cannot** pay existing debt or buy real estate.
  Best first stop for a small/early venture.
- **SBA 7(a)** — flagship general-purpose loan, up to **$5 million**; terms up to ~10 yrs
  (working capital/equipment) or 25 yrs (real estate).
- **SBA 504** — long-term, fixed-rate, for major fixed assets (real estate, heavy equipment), SBA-backed
  portion up to **$5 million** (some sources cite $5.5M for the debenture).
- **Recent change [verify]:** effective **July 4, 2026**, SBA doubled the cumulative cap so eligible
  borrowers can combine 7(a) + 504 for up to **$10 million** total (was $5M combined); small manufacturers
  get added flexibility (sba.gov, 2026-05-18). Confirm the rule is in force and applies to you.

### North Carolina small-business resources

- **SBA North Carolina District Office** — serves the entire state; offices in **Charlotte and
  Wilmington**. Helps with SBA funding programs, counseling, federal contracting certifications, and
  connects you to lenders and partners (sba.gov/district/north-carolina, as of 2026-06).
- **NC SBTDC (Small Business and Technology Development Center)** — the business/technology extension
  service of the UNC System, administered by NC State, in partnership with the SBA. **Free, confidential**
  management counseling; ~10 regional centers / 16 offices statewide plus specialty programs
  (sbtdc.org, as of 2026-06). Strong first call for one-on-one help.
- **SCORE** — national network of 13,000+ volunteer mentors offering **free, confidential** business
  mentoring; an SBA resource partner with NC chapters. Good for an experienced sounding board.
- **NC IDEA grants** — non-dilutive grants for NC startups (ncidea.org, as of 2026-06):
  - **NC IDEA MICRO** — **$10,000** non-dilutive, for very early companies validating assumptions /
    customer discovery; includes an ~8-week customer-discovery & product-launch program.
  - **NC IDEA SEED** — **$50,000** non-dilutive, released over 6-12 months in milestone-based tranches;
    for companies with proof of concept and an MVP. Eligibility: HQ'd in NC with majority of operations
    in-state, ≥1 full-time founder living in NC, addressing a large/growing market; *less likely* if
    revenue exceeds $250K (or $500K food/beverage) in the trailing 12 months or you've raised >$250K
    equity / >$1M grants. **You may apply to MICRO or SEED, not both, per cycle.** Cycles run twice a
    year (Spring opens ~late Jan; Fall ~late July).

### Nonprofit funding differs

A nonprofit's "funding stack" is **grants + donations + earned/program revenue**, not equity. Key points:

- **Grants & donations** are the core; build a diversified mix so no single funder can sink you. Donor
  acquisition and retention behave like CAC/LTV — model them.
- **Fiscal sponsorship** — a faster, cheaper alternative to standing up your own 501(c)(3) at the start:
  an existing 501(c)(3) (the sponsor) extends its charitable status to your mission-aligned project, so
  you can receive **tax-deductible donations and grants in weeks** while they handle fiduciary oversight.
  Typical sponsor fee **~5-15% of funds raised**; forming your own 501(c)(3) often runs **~$2,000-$20,000**
  and ongoing compliance cost (nptechforgood.com; gvng.org; missionedge.org, as of 2026-06). Common
  models: **Model A** (project is part of the sponsor) vs **Model C** (grantor-grantee). A pragmatic path
  is *fiscal sponsorship now → your own 501(c)(3) once funding justifies it.*
- For the legal/tax mechanics of 501(c)(3), charitable-solicitation registration, and Form 990, see
  `venture-nc-nonprofit-formation`. For donor-journey, campaigns, and Google Ad Grants, see
  `venture-cause-nonprofit-marketing`.

---

## 6. KPIs, milestones & a simple operating plan

### KPIs

KPIs track progress toward your objectives (your OKRs). **Pick the few that reflect your current
bottleneck**, not a dashboard of 20. Define each consistently, review on a fixed cadence, and assign each
to a decision owner (geckoboard.com; mercury.com, as of 2026-06). Common early-stage KPIs:

- **Cash runway** = current cash ÷ net monthly burn — *how many months until you're out of money.* The
  master survival metric. Aim for **12-18 months**.
- **Burn rate** — *gross burn* = total monthly cash out; *net burn* = cash out minus revenue in.
- **Revenue / MRR** and growth rate.
- **CAC, LTV, LTV:CAC, CAC payback** — the unit-economics set from §4.
- **Retention / churn** — usually a better PMF signal than raw acquisition.
- **Gross margin** and **contribution margin**.

### Milestones

Tie spending to **measurable milestones**, not the calendar. Investors and lenders want to see that each
dollar buys progress (problem-solution fit → first paying customers → break-even → PMF → growth).
Sequence: validate problem → MVP → first revenue → repeatable sales → break-even → scale.

### A simple operating plan

For a solo founder, one page is enough: the next **2-3 quarters'** objectives, the milestone for each,
the metric that proves it, the owner (you, until you hire), and the cash that funds it. Connect it to the
financial model so plan and money tell the same story. Review monthly; adjust as evidence arrives.

**Nonprofit note:** swap profit KPIs for **mission/impact KPIs** (e.g., donors registered, people
reached, outcomes delivered) *plus* the financial-health basics (runway, funding diversity, cost per
beneficiary). Beware **vanity metrics** (impressions, page views) that don't map to mission outcomes —
see `venture-cause-nonprofit-marketing`.

---

## 7. Risk & assumptions discipline

Every model and plan rests on assumptions; the failures come from the ones you never wrote down.

- **List your assumptions explicitly** — in the Lean Canvas, the financial model, and the plan. Mark each
  as *tested* or *untested*, and rank by how much damage a wrong guess does. Attack the riskiest untested
  assumption first (that's the point of customer discovery).
- **The riskiest assumptions are usually:** "people have this problem and will pay/act," your CAC, your
  churn/retention, and your pricing. These four decide whether the model is real.
- **Stress-test the model:** run a pessimistic case (higher CAC, higher churn, slower sales, longer
  collections). If the pessimistic case runs out of cash, you need more runway, lower burn, or faster
  validation *before* you commit.
- **Watch the cash, always.** Most small businesses die from running out of cash, not from a bad idea.
  Keep a reserve; know your runway every month.
- **Separate facts from forecasts** in any lender/funder document — credibility comes from showing your
  reasoning, not from big numbers.
- **Revisit on a cadence.** A plan is a living hypothesis; update it as evidence arrives. The plan that
  never changes is the plan nobody is using.

---

## Sources

Authoritative + practitioner sources consulted (as of 2026-06; verify program/loan/grant specifics
against the official source — figures change):

- SBA — *Write your business plan* (traditional vs lean): sba.gov/business-guide/plan-your-business/write-your-business-plan
- SBA — *Calculate your startup costs*: sba.gov/business-guide/plan-your-business/calculate-your-startup-costs
- SBA — *Microloans*: sba.gov/funding-programs/loans/microloans · *504 loans*: sba.gov/funding-programs/loans/504-loans
- SBA — *Doubles Cumulative 7(a)/504 Limit to $10M* (2026-05-18): sba.gov/article/2026/05/18/sba-doubles-cumulative-7a-504-loan-limit-10-million
- SBA — *North Carolina District Office*: sba.gov/district/north-carolina
- NC SBTDC: sbtdc.org · SCORE: score.org
- NC IDEA — *SEED*: ncidea.org/nc-idea-seed/ · grant-eligibility comparison: ncidea.org/comparison-of-grant-eligibility/
- Strategyzer — Business Model Canvas (Osterwalder): strategyzer.com
- LeanStack / Ash Maurya — Lean Canvas: leanstack.com/lean-canvas
- Steve Blank — Customer Discovery / *Four Steps to the Epiphany*: steveblank.com; Rob Fitzpatrick — *The Mom Test*
- Harvard Business School Online — *Cash Flow vs. Profit*: online.hbs.edu/blog/post/cash-flow-vs-profit
- TAM/SAM/SOM — Forum VC: forumvc.com/thought-pieces/understand-and-define-your-market-size; Qubit Capital
- Unit economics (LTV/CAC, contribution margin) — Kruze Consulting: kruzeconsulting.com/blog/unit-economics/; ltvcacbook.com
- Pricing strategies — QuickBooks, Salesforce, BDC: quickbooks.intuit.com/r/pricing-strategy/pricing-strategies/
- Startup KPIs / burn & runway — Geckoboard, Mercury
- Fiscal sponsorship vs 501(c)(3) — Nonprofit Tech for Good, GVNG, Mission Edge: nptechforgood.com

## Disclaimer

This skill is general educational guidance for planning and financial modeling — **not** financial,
investment, accounting, legal, or tax advice. Frameworks (BMC, Lean Canvas, unit economics) are tools,
not guarantees; every model depends on assumptions you must test. **All program, loan, grant, and
eligibility figures change** — verify current terms directly with the SBA, your lender or SBA-approved
intermediary, NC SBTDC, SCORE, and NC IDEA (and read the official program pages) before making decisions
or applying. Items marked **[verify]** were accurate to the cited source as of 2026-06 but should be
re-checked. For entity choice and tax, consult `venture-nc-business-formation-tax` or
`venture-nc-nonprofit-formation` and a qualified attorney/CPA. Talk to a professional before signing
loan documents, taking investment, or committing significant capital.
