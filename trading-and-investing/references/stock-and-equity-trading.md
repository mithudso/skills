<!-- Provenance: spoke reference under the `trading-and-investing` hub. Created 2026-06-16 via /dr deep-research (6 concepts; ~40 independent sources, regulator/exchange/index-provider/peer-reviewed grade). Educational only — NOT financial/investment/tax advice. Volatile claims dated "as of 2026". -->

# Stock & Equity Trading (spoke)

The **equity** spoke of the `trading-and-investing` family: trading individual stocks (and equity ETFs as instruments) from a US retail participant's seat — what a share *is*, how equities trade and are priced, how to read a company's fundamentals, the corporate actions that move a stock, the primary market that creates new shares, the long/short mechanics specific to equities, and the central question of **active stock-picking vs passive index investing**.

> **Educational information only — NOT financial, investment, tax, or legal advice, and NOT a recommendation to buy, sell, or hold any security.** Equities carry **risk of loss**; a stock can fall to **$0** (you lose 100%), and **short selling can lose more than you invest** (losses are theoretically unlimited). The evidence below is blunt that **active stock-picking underperforms a low-cost index for the large majority of participants**. Rules, tax treatment, and the vendor/market landscape change — every volatile claim here is stamped **as of 2026**; verify current facts with the primary regulator (SEC/Investor.gov, FINRA) before acting. For a personal decision, consult a licensed professional.

**This spoke builds on the hub foundation — read those first, not re-derived here:** asset-class definition of equities → `references/asset-classes-and-instruments.md`; exchanges, clearing/settlement (T+1), bid-ask spread, primary vs secondary market → `references/market-participants-and-structure.md`; order types and routing (PFOF, NBBO) → `references/order-lifecycle-and-execution.md`; sessions/circuit breakers → `references/market-sessions-and-venues.md`; the investing-vs-trading distinction and the underperformance evidence overview → `references/investing-vs-trading.md`; margin, the Pattern-Day-Trader rule and its 2026 replacement, short-selling overview, SIPC → `references/trading-risks-and-protections.md`.

**Peer-spoke handoffs (SKIP edges — go there for depth):**
- Listed **options** on equities (calls/puts, the Greeks, spreads, the gamma squeeze mechanics) → `options-trading-and-strategies`.
- **Futures, forwards, swaps** and commodity-futures depth → `derivatives-futures-and-swaps`.
- **Chart-based** entry/exit timing, candlesticks, RSI/MACD, support/resistance → `technical-analysis`.
- **Order-routing / microstructure** depth (dark pools, Reg NMS internals, maker-taker, slippage, TWAP/VWAP) → `market-microstructure-and-execution`.
- **Long-term passive / retirement** investing — 401(k)/IRA/Roth, target-date funds, dollar-cost averaging for retirement, choosing a robo or fee-only fiduciary, picking a fund by expense ratio for a retirement portfolio → the peer skill `investing-and-retirement`. **This spoke owns the active side and the active-vs-passive evidence; it hands off the passive playbook.**
- Position sizing / stops / Kelly → `trading-risk-management`; the trader's tax detail (wash sales, trader-tax-status, capital-gains mechanics) → `trading-regulation-compliance-and-taxes`.

## Contents

- [1. Equity instruments & market structure](#instruments)
  - [Common vs preferred stock](#common-preferred)
  - [Share classes & dual-class structures](#share-classes)
  - [How shares trade: NYSE vs Nasdaq](#nyse-nasdaq)
  - [Major US indices & how they are weighted](#indices)
  - [Reading a stock quote](#quote)
  - [Market capitalization, tiers, float vs shares outstanding](#market-cap)
  - [GICS sectors & industries](#gics)
- [2. Fundamental equity analysis basics](#fundamentals)
  - [Valuation multiples (P/E, earnings yield, PEG)](#multiples)
  - [EPS — basic vs diluted](#eps)
  - [P/B and book value](#pb)
  - [Revenue vs earnings, margins](#revenue-earnings)
  - [Dividend yield, payout ratio, coverage](#yield-payout)
  - [Balance-sheet basics](#balance-sheet)
  - [Earnings season, guidance & analyst ratings](#earnings-season)
  - [Investing styles: growth, value, dividend](#styles)
- [3. Dividends, buybacks & stock splits (corporate actions)](#corporate-actions)
  - [Cash-dividend mechanics & the four dates](#dividend-dates)
  - [DRIPs, special & stock dividends](#drips)
  - [Share buybacks](#buybacks)
  - [Stock splits & reverse splits](#splits)
- [4. Corporate actions & the primary market](#primary-market)
  - [IPOs](#ipos)
  - [Direct listings](#direct-listings)
  - [SPACs](#spacs)
  - [Follow-on offerings & the "secondary" terminology trap](#followon)
  - [Lockup periods](#lockups)
- [5. Equity-specific trading mechanics: long vs short](#long-short)
  - [Going long](#long)
  - [Short selling & the locate/borrow requirement](#short)
  - [Short interest & the short squeeze](#squeeze)
  - [Holding horizon: long-term vs active](#horizon)
- [6. Active equity trading vs passive index investing](#active-passive)
- [References](#references)

---

## 1. Equity instruments & market structure
<a id="instruments"></a>

### Common vs preferred stock
<a id="common-preferred"></a>

A **stock** is a security giving fractional **ownership** in a company. Two main classes:

- **Common stock** entitles the owner to **vote** at shareholder meetings (typically one vote per share) and to receive **dividends** when declared. In a liquidation, common holders are paid **last** — after creditors, bondholders, and preferred holders.[^1]
- **Preferred stock** usually carries **no voting rights**, receives its (typically fixed) dividend **before** common holders, and ranks **ahead of** common in bankruptcy/liquidation.[^1] It is a **hybrid** security — senior to common but **subordinate to bonds** — and because the dividend is usually fixed and the share often perpetual, its price behaves like a long-duration bond (it moves inversely with interest rates).[^2]

Preferred sub-types (combinable): **cumulative** (skipped dividends accrue and must be paid before any common dividend) vs **non-cumulative** (a skipped dividend is simply lost); **callable/redeemable** (issuer can buy it back, usually at a premium); **convertible** (holder can convert to common at a preset ratio); **participating** (can receive extra beyond the fixed rate).[^2] Most retail "stock trading" means **common** stock; preferred is closer to fixed income and is covered at overview depth in the asset-classes foundation.

### Share classes & dual-class structures
<a id="share-classes"></a>

A company can issue **multiple classes of common stock** with unequal voting power — a **dual-class** (or multi-class) structure that concentrates control with founders/insiders via "super-voting" shares.[^3]

- **Alphabet:** `GOOGL` = Class A (1 vote); `GOOG` = Class C (no vote); Class B (10 votes/share) is held by insiders and does not trade publicly.[^3]
- **Berkshire Hathaway:** `BRK.A` (full voting, very high price) vs `BRK.B` (created 1996, a small economic fraction of an A share with proportionally far less voting power).[^3]

**The governance debate:** critics (notably the Council of Institutional Investors) argue dual-class violates "one share, one vote," entrenches managers, and weakens takeover/activist discipline; proponents say it protects founder vision and long-term strategy.[^3]

**Index eligibility — a load-bearing, time-stamped reversal** *(as of 2026)*: in **July 2017** S&P Dow Jones Indices barred companies with multiple share classes from **new** addition to the S&P Composite 1500 (existing members like Alphabet were grandfathered); FTSE Russell set a free-float voting threshold the same year. In **April 2023, S&P Dow Jones Indices reversed the ban** — multi-class companies are again eligible for the S&P 500/1500 if they meet all other criteria.[^4] Re-verify before relying, as index policy changes.

### How shares trade: NYSE vs Nasdaq
<a id="nyse-nasdaq"></a>

The two dominant US listing venues run **different market models** (the deeper microstructure of *how* an order executes lives in `market-microstructure-and-execution`):

- **NYSE = an auction market with Designated Market Makers (DMMs).** Each listed stock is assigned to **one** DMM who manages the order book and runs the opening and closing auctions, with a formal obligation to maintain a "fair and orderly market" and supply liquidity to offset imbalances. NYSE is a **hybrid** model and retains a physical trading floor. (The DMM role was called the "specialist" before 2008.)[^5]
- **Nasdaq = a dealer/electronic market with multiple competing market makers.** No single specialist; many dealers post competing two-sided quotes and orders route to the best price. Liquidity comes from dealer competition.[^5]

**Listing requirements** differ — NYSE generally sets a higher bar; Nasdaq is tiered (Global Select / Global / Capital Market) and typically less costly — and the specific dollar thresholds change (the SEC approved new NYSE/Nasdaq **minimum-price** rules in 2025–2026), so treat any specific figure as volatile and check the current listing standards.[^5]

### Major US indices & how they are weighted
<a id="indices"></a>

The three headline US indices use **three different weighting methods** — a common point of confusion:

- **S&P 500 — float-adjusted market-cap weighted, committee-selected.** Weight ∝ float-adjusted market value (only publicly tradable shares count). It is **~500 companies but not simply "the 500 largest" and not purely rules-based** — S&P's Index Committee exercises **full discretion**; quantitative criteria only earn a place on an eligibility list.[^6] *(As of 2026)* it actually carries **~503 stock listings** because a few companies (e.g., Alphabet) have two included share classes; the count drifts.[^6]
- **Nasdaq-100 — modified market-cap weighted, the 100 largest non-financial Nasdaq companies.** Each constituent's share count is **capped at 3× its free-float** to limit low-float dominance, with capping rules that scale down the largest weights; reconstituted annually in **December** with quarterly rebalances, plus an off-cycle **"special rebalance"** if concentration thresholds are breached.[^7]
- **Dow Jones Industrial Average — PRICE-weighted, 30 stocks.** Higher-**priced** stocks have more influence regardless of company size; DJIA = (sum of the 30 share prices) ÷ the **Dow Divisor** (adjusted for splits/component changes; now <1, so it acts as a multiplier).[^8] Price-weighting is **widely criticized** as conceptually inferior — a stock split arbitrarily cuts a component's weight, and a high-priced mid-size company can outweigh a far larger low-priced one.[^8]

Briefly: the **Russell 2000** is the small-cap benchmark (the Russell 1000 large-cap, Russell 3000 ≈ the whole investable market), all float-cap weighted; **total-market** indices (CRSP US Total Market, Wilshire 5000) aim to capture the entire US market.[^9] *(As of 2026, FTSE Russell moved Russell reconstitution from annual to semi-annual — verify the effective schedule.)*[^9]

### Reading a stock quote
<a id="quote"></a>

A retail quote shows, at minimum:[^10]

- **Bid** — the highest price a buyer is currently willing to pay; **ask (offer)** — the lowest a seller will accept; the **bid-ask spread** (ask − bid) is **narrow** when liquidity is high and **wide** when it is low (the spread is a real round-trip cost — see the order-lifecycle foundation).
- **Last** — the price of the most recent executed trade.
- **Volume** — shares traded (today or average daily); a liquidity/interest gauge.
- **Day range** and **52-week range** — the low–high over the session and the trailing year.
- **Market cap** — current price × shares outstanding.

### Market capitalization, tiers, float vs shares outstanding
<a id="market-cap"></a>

**Market cap = current share price × shares outstanding.**[^11] Size **tiers are conventions, not an official standard** — providers differ — so treat thresholds as approximate *(as of 2026)*: FINRA-cited bands run **mega-cap > $200B, large $10B–$200B, mid $2B–$10B, small ~$250M–$2B, micro ~$50M–$250M, nano < $50M**; larger caps tend to be less volatile.[^11] Morningstar instead uses a **relative** method (top 70% of total US market cap = large, next 20% = mid, etc.).[^11]

**Float vs shares outstanding:** **shares outstanding** = all issued shares (including restricted insider/locked blocks); **public float** = the freely tradable subset (= outstanding − restricted/closely-held). Float ≤ shares outstanding always, and index providers weight on **float**, not total shares.[^12]

### GICS sectors & industries
<a id="gics"></a>

The **Global Industry Classification Standard (GICS)** — maintained jointly by **MSCI and S&P Dow Jones Indices** — sorts every public company by its principal business into a four-tier hierarchy: **Sector → Industry Group → Industry → Sub-Industry**.[^13] *(As of 2026, last revised March 2023)* the structure is **11 sectors, 25 industry groups, 74 industries, 163 sub-industries**.[^13] The **11 sectors:** Energy, Materials, Industrials, Consumer Discretionary, Consumer Staples, Health Care, Financials, Information Technology, Communication Services, Utilities, Real Estate.[^13] (Real Estate was split from Financials in 2016; Communication Services was created in 2018.)[^13] GICS is the standard vocabulary for "what sector is this stock in" and for sector ETFs/screens.

---

## 2. Fundamental equity analysis basics
<a id="fundamentals"></a>

Fundamental analysis values a company from its business and financial statements — distinct from **technical analysis** (price/volume charts → `technical-analysis`). The basics below are analytical *metrics*; they are not a valuation method on their own and every multiple has limits.

### Valuation multiples (P/E, earnings yield, PEG)
<a id="multiples"></a>

- **P/E ratio** = share price ÷ EPS (equivalently market cap ÷ net earnings) — "how much you pay per $1 of earnings."[^14] **Trailing** P/E uses actual last-twelve-months EPS (hard data); **forward** P/E uses projected next-12-month EPS (depends on estimates that can be wrong).[^14] A **high** P/E signals high growth expectations (or overvaluation and more volatility); a **low** P/E signals a perceived value play (or trouble).[^14]
- **Earnings yield** = EPS ÷ price = 1 ÷ P/E, as a percentage — lets you compare a stock's earnings return against bond yields.[^15]
- **PEG** = P/E ÷ earnings growth rate; a rule of thumb treats **PEG < 1** as cheap relative to growth.[^15]
- **Limits (load-bearing):** P/E is **negative/meaningless for unprofitable companies** (analysts fall back on price-to-sales); it uses accounting earnings that one-offs can distort; it ignores balance-sheet leverage and growth; and a **low P/E can be a "value trap"** — cheap for a deteriorating reason.[^16]

### EPS — basic vs diluted
<a id="eps"></a>

**EPS** = (net income − preferred dividends) ÷ shares outstanding — earnings available to **common** holders.[^17] **Basic EPS** uses actual common shares; **diluted EPS** adds all dilutive potential shares (options, warrants, convertible debt/preferred) as if issued, so **diluted EPS ≤ basic EPS**.[^17] US GAAP (ASC 260) requires public companies to present both with equal prominence.[^17] Dilution matters because option/convertible-heavy companies report a meaningfully lower diluted number.

### P/B and book value
<a id="pb"></a>

**P/B** = price ÷ book value per share, where **book value = shareholders' equity (assets − liabilities)**.[^18] P/B is most meaningful for **asset-heavy** businesses (banks, insurers, real estate, manufacturers) and **understates** value for **intangible-heavy / asset-light** firms (tech, pharma, services) because R&D, brands, and software are mostly **expensed, not capitalized** — so book value is understated and P/B looks very high.[^18] Buybacks above book value, and even negative book equity (e.g., long-time repurchasers), further distort P/B.[^18]

### Revenue vs earnings, margins
<a id="revenue-earnings"></a>

**Top line** = total **revenue**; **bottom line** = **net income** (after all costs, interest, taxes).[^19] Strong revenue doesn't guarantee profit if costs rise in step. The **margin** ladder strips successive cost layers: **gross margin** (revenue − COGS), **operating margin** (after SG&A and R&D), **net margin** (after interest, taxes, one-offs).[^19] Both revenue growth (demand/scale) and earnings growth (cost discipline → shareholder returns) matter; profitless revenue and margin-only growth are each incomplete signals.[^19]

### Dividend yield, payout ratio, coverage
<a id="yield-payout"></a>

- **Dividend yield** = annual dividend per share ÷ price (the cash return on the price paid).[^20]
- **Payout ratio** = dividends ÷ earnings (the share of earnings paid out vs retained). Commonly cited as **sustainable ~30–60%** (conventions, not hard rules; stable sectors like utilities run higher), **caution 60–80%**, **danger > 80%**, and **> 100% = paying out more than it earns** (a red flag).[^20]
- **Dividend coverage ratio** = earnings ÷ dividends (the inverse); > 1 means the dividend is earned.[^20]
- **Disconfirming — the high-yield "dividend trap":** a yield far above peers is usually a warning that the **price has fallen** on deteriorating fundamentals and the market expects a **cut** (which typically drops the price further). A high yield is not free income.[^21]

### Balance-sheet basics
<a id="balance-sheet"></a>

The accounting identity is **Assets = Liabilities + Equity**; equity is the **residual** owners' claim after liabilities.[^22] Equity holders are paid **last**, after creditors — the conceptual reason equity is riskier than debt. Two everyday gauges: **debt-to-equity** (total liabilities ÷ equity; higher = more borrowed-capital reliance) and the **current ratio** (current assets ÷ current liabilities; > 1 covers near-term obligations).[^22] Higher leverage magnifies losses to equity in downturns because debt's fixed claims rank ahead of the residual equity claim.[^22]

### Earnings season, guidance & analyst ratings
<a id="earnings-season"></a>

Public companies file an **annual Form 10-K** (audited) and **quarterly Form 10-Q** (generally unaudited) with the SEC, with material events on **Form 8-K**; all are public on **EDGAR**.[^23] Around each quarter-end ("earnings season") a company releases results and hosts an **earnings call**. Analysts publish **consensus estimates**; results above are **beats**, below are **misses**.[^24]

**Load-bearing nuance:** because a stock's value is the present value of **future** cash flows, the market prices on **changes in expectations** — so **forward guidance often moves a stock more than the headline number**. A company can beat on earnings yet fall sharply if it cuts guidance, and high-multiple stocks swing violently when guidance resets the earnings their valuation was built on.[^24] **Analyst ratings** (Buy/Outperform, Hold/Neutral, Sell/Underperform) express **expectations, not guarantees**; a consensus rating averaging many analysts is more reliable than any single call (note many scales count **1 = Strong Buy** up to 5 = Strong Sell).[^24]

### Investing styles: growth, value, dividend
<a id="styles"></a>

Three analytical lenses — not mutually exclusive:[^25]

- **Growth:** high revenue/earnings growth, **high multiples** (P/E 30–50+; unprofitable names judged on price-to-sales), little/no dividend (profits reinvested), concentrated in tech/innovation, higher volatility.
- **Value (Graham/Buffett):** buy below estimated **intrinsic value** with a **margin of safety** (the cushion against estimation error); **low multiples**; Buffett evolved toward "high-quality businesses at fair prices" with durable moats.
- **Dividend / income:** companies that pay and ideally **grow** dividends; favors dividend **growth** over chasing the highest (possibly trap) yield; **Dividend Aristocrats** = S&P 500 members that have raised dividends for **≥25 consecutive years** *(as of 2026, a record ~69 companies — re-verify the count, which changes yearly)*.[^25]

**Cyclicality** *(as of 2026, volatile)*: **growth outperformed value for much of the 2010s–early 2020s** (low rates, tech leadership), but leadership rotates (value led in 2022; growth led 2023–24; a value rotation began in 2025). Value's worst stretches have historically preceded its best — frame style leadership as cyclical, not permanent.[^26] Buffett himself rejects a rigid growth-vs-value split, calling growth a component of value.[^25]

---

## 3. Dividends, buybacks & stock splits (corporate actions)
<a id="corporate-actions"></a>

How a company returns capital or restructures its shares — the *mechanics* (the yield/payout *metrics* are in §2).

### Cash-dividend mechanics & the four dates
<a id="dividend-dates"></a>

A cash dividend runs through four dates:[^27]

1. **Declaration date** — the board announces the dividend, amount, record date, and payment date.
2. **Ex-dividend date (ex-date)** — the day the stock starts trading **without** the right to the upcoming dividend. **Buy before the ex-date → you get the dividend; buy on/after → the seller keeps it.**[^27]
3. **Record date** — you must be on the company's books as a holder to receive the dividend.
4. **Payment date** — the dividend is actually paid.

**Load-bearing, T+1-dependent fact** *(as of 2026)*: since the US moved to **T+1 settlement (May 28, 2024)**, the **ex-dividend date and record date are generally the SAME business day** for regular-way dividends (under the old T+2 the ex-date was one business day before the record date). The clearing utility **DTCC** confirms this directly; the ex-date shifts one business day earlier only if the record date is not a business day.[^28] On the ex-date the share price typically opens **lower by approximately the dividend** (the cash will leave the company), though the real move is only approximate (taxes, ticks, and ordinary movement intervene) — so **"buy just before the ex-date for free money" is a myth**: the price drop offsets the dividend (often with a tax cost).[^28] In theory (Modigliani–Miller, frictionless markets) dividend policy is value-neutral; real-world frictions (taxes, signaling) are why dividends matter in practice.[^29] **Qualified** dividends are taxed at lower long-term-gains rates if an IRS holding-period test is met; **ordinary** dividends at ordinary rates — deep tax → `trading-regulation-compliance-and-taxes` / the consumer-finance `personal-income-taxes` skill.[^27]

### DRIPs, special & stock dividends
<a id="drips"></a>

- **DRIP (dividend reinvestment plan):** automatically reinvests cash dividends into more (often fractional) shares, usually commission-free, compounding the share count over time. Company-sponsored DRIPs may offer a small discount; brokerage "synthetic" DRIPs reinvest even without a formal company plan.[^30]
- **Special dividend:** a **one-time** payment (often large) after an unusually profitable period or asset sale — not part of the regular cadence.[^31]
- **Stock dividend:** paid in **additional shares** rather than cash; raises share count and is economically similar to a small split (generally not immediately taxable, with exceptions).[^31]

### Share buybacks
<a id="buybacks"></a>

A **buyback (repurchase)** returns capital by **reducing shares outstanding** — remaining holders own a larger slice, and **EPS rises mechanically** (same earnings ÷ fewer shares) even with no operational improvement.[^32] Three methods: **open-market repurchase** (most common, ~95%, under the SEC **Rule 10b-18** safe harbor that limits daily volume/price), **tender offer** (buy a set block at a stated price, usually at a premium), and **accelerated share repurchase (ASR)** (a bank delivers a large block upfront for immediate share-count reduction).[^32]

**Buybacks vs dividends:** buybacks are **flexible** (pause/resize without the penalty markets attach to a dividend cut); they create **no immediate shareholder tax** (deferred to sale) vs dividends taxed when paid; dividends **signal** stable recurring earnings while buybacks often signal management thinks shares are **undervalued**.[^32] *(As of 2026)* a **1% excise tax** on net stock buybacks by public companies (IRC §4501, Inflation Reduction Act) applies to repurchases after **Dec 31, 2022**, with the taxable amount **netted** against new stock issued the same year; Treasury/IRS **finalized the regulations in November 2025 at the 1% rate**.[^33] A **proposal to raise it to 4%** (the Stock Buyback Accountability Act of 2026) exists but is **not enacted** — treat as tentative and re-verify.[^33] **Criticism** is real and widely held (EPS/financial engineering, executive-comp incentives tied to EPS, opportunity cost vs R&D/wages, leveraged-buyback fragility), though defenders rebut it.[^34]

### Stock splits & reverse splits
<a id="splits"></a>

- **Forward split** (e.g., 2-for-1, 10-for-1): multiplies share count and divides price proportionally — **no change to market cap, ownership %, voting, or fundamental value.**[^35] Companies split to lower the **nominal** price for accessibility/liquidity, but **fractional shares** (now standard at major brokers) largely remove that rationale.[^35] Recent high-profile splits: **NVIDIA 10-for-1**, **Chipotle 50-for-1**, **Broadcom 10-for-1** (all 2024).[^35] **Disconfirming:** splits **do not create value**; there is a small positive **announcement effect** (~2–4% abnormal return) attributed to **signaling**, not value creation, with mixed long-term evidence.[^36]
- **Reverse split** (e.g., 1-for-10): consolidates shares into fewer at a proportionally higher price; **no value change.** The usual purpose is to lift the price back above the **$1.00 minimum bid** for continued NYSE/Nasdaq listing (Nasdaq triggers a deficiency below $1 for 30 consecutive business days; cure by trading ≥ $1 for 10 consecutive days).[^37] *(As of 2026)* the SEC approved **2024–2025 Nasdaq/NYSE rule changes restricting reverse-split abuse** — a company that recently reverse-split (within the prior year) is **not** eligible for a new bid-price cure period.[^37] Reverse splits carry a **distress stigma** and often see continued selling pressure.[^37]

---

## 4. Corporate actions & the primary market
<a id="primary-market"></a>

How shares are **first created and sold** (the **primary** market) — distinct from the **secondary** market where investors trade existing shares (see the market-participants foundation).

### IPOs
<a id="ipos"></a>

An **initial public offering** is a company's first sale of shares to the public.[^38] In a **firm-commitment** underwriting, investment-bank **underwriters** (a lead **bookrunner** plus a **syndicate**) buy the shares from the issuer and resell them, bearing inventory risk. The process: file a **Form S-1** registration statement with the SEC (which reviews **disclosure adequacy**, not merits), run a **roadshow**, **build a book** of indications of interest, set the **offer price**, and typically go effective the morning of pricing.[^38] A **greenshoe / over-allotment** option lets underwriters sell **up to 15% more shares** and buy them back at the offer price for **price stabilization**.[^39]

**Disconfirming, load-bearing for retail:** IPO shares **skew heavily institutional** — roughly **~90% to institutions, ~10% to individuals** — and brokerages that offer retail IPO access impose asset minimums (often **$100k–$500k**) with **no guarantee** of an allocation.[^40] FINRA **Rules 5130/5131** bar selling new issues to "restricted persons" (industry insiders) and bar "spinning" (rewarding executives with allocations).[^41] First-day **"pops"** are well documented (long-run US average first-day return **~19%**), but studies find IPOs tend to **underperform the broad market over ~3 years**, with the effect partly shrinking against size/style-matched benchmarks.[^42] A **quiet period** and **gun-jumping** rules (Securities Act §5) restrict promotional communication around the filing.[^38]

### Direct listings
<a id="direct-listings"></a>

In a (traditional) **direct listing**, **no new shares are created and no capital is raised** — **existing** shares are sold directly on the exchange with **no underwriter**; a financial advisor and the exchange facilitate **price discovery via an opening auction** (a published **reference price** is an indication, not an offer price).[^43] The company still files an S-1 and clears SEC review. **The SEC approved primary (capital-raising) direct listings** — NYSE in **December 2020**, Nasdaq in **May 2021** — letting a company issue new shares without a firm-commitment underwriting.[^44] **Pros:** avoids the ~6–7% underwriting spread and the IPO-pop "money left on the table," immediate holder liquidity, and **often no lockup**. **Cons:** no capital raised (classic version), more **price volatility** without underwriter stabilization, and no built-in institutional demand book — suited mainly to large, well-known brands (Spotify 2018, Slack 2019, Coinbase 2021).[^43] *(As of 2026)* despite the 2020–21 rule approvals, the **primary capital-raising direct listing has seen essentially no adoption** — verify against current data before stating definitively.[^44]

### SPACs
<a id="spacs"></a>

A **SPAC (special-purpose acquisition company / "blank-check" company)** is a shell that **IPOs to raise cash into a trust** (usually $10/unit), then has roughly **two years (up to three)** to find and merge with a private target (the **de-SPAC**); if no deal closes it liquidates and returns trust cash.[^45] Public holders get **redemption rights** (redeem for the pro-rata trust ≈ $10) at the vote — but the SEC cautions that if you bought above $10 you only get the **pro-rata trust amount, not your purchase price**.[^45] The sponsor typically takes a **~20% "promote"** (founder shares for nominal cost) plus warrants — a dilution and **incentive-misalignment** source (the promote can reward closing **any** deal).[^45]

*(As of 2026, historical)* SPACs boomed in **2020–2021** (2021 alone: ~613 SPAC IPOs, >$160B, over half of US listings) then **busted in 2022+** (a large share of 2021 SPACs liquidated).[^46] **Post-merger de-SPAC stocks have broadly underperformed**, with most falling below the $10 redemption value — directly answering "are SPACs a bad investment": on average, poorly.[^46] The SEC **adopted final SPAC rules in January 2024**: enhanced disclosure of sponsor compensation/conflicts/dilution, **de-SPAC↔IPO alignment** (the target becomes a co-registrant with **Section 11 liability**), required disclosure of **projections** and removal of the **PSLRA safe harbor** for them; notably, the SEC **did not adopt a specific underwriter-liability rule**, instead issuing **guidance** that it reads "underwriter" broadly in a de-SPAC.[^47]

### Follow-on offerings & the "secondary" terminology trap
<a id="followon"></a>

An **already-public** company can issue/sell more shares in a **follow-on offering**:[^48]

- **Dilutive (primary) follow-on:** the company creates and sells **new** shares; proceeds go **to the company**; existing holders are **diluted** (lower EPS).
- **Non-dilutive (secondary) follow-on:** **existing** shares (held by insiders/large holders) are sold; **no new shares**, proceeds go **to the sellers**.

Either way the price tends to fall, via **mechanical dilution** and **negative signaling** (issuing equity can imply management thinks shares are fully valued); seasoned-equity-offering announcement returns are typically **negative**.[^48] **Shelf registration (Form S-3)** lets a qualified company pre-register and sell shares "off the shelf" over a 3-year window; an **at-the-market (ATM)** offering dribbles shares into the market at prevailing prices.[^48]

**The terminology trap (clarify explicitly):** the **"secondary market"** is where investors trade **already-issued** shares among themselves (everyday NYSE/Nasdaq trading) — no money to the issuer. A **"secondary offering"** is a **primary-market event** — a registered sale of a large block of **existing** stock by big holders. Worse, people loosely call a **dilutive follow-on** a "secondary offering" even though a strict *secondary* offering sells **existing** (non-dilutive) shares. Prefer the precise phrasing: **"dilutive follow-on (new shares)"** vs **"non-dilutive secondary (insiders sell existing shares)"**, and reserve **"secondary market"** for investor-to-investor trading.[^48]

### Lockup periods
<a id="lockups"></a>

A **lockup** is a **contractual** restriction (imposed by underwriters, not an SEC mandate) barring insiders/early investors from selling for a set window after the IPO — **commonly 90–180 days (180 is standard)**.[^49] Its purpose is to prevent a post-IPO **supply shock** and to **signal** insider confidence. When the lockup **expires**, the tradable float can jump sharply and the stock often sees **price pressure** around the expiration. **Direct listings frequently have no traditional lockup** (no underwriter to impose one), letting insiders sell from day one.[^49]

---

## 5. Equity-specific trading mechanics: long vs short
<a id="long-short"></a>

Margin, the Pattern-Day-Trader rule, and the short-selling overview are in the `trading-risks-and-protections` foundation; this section is the **equity-specific** depth.

### Going long
<a id="long"></a>

**Going long** = buying shares to profit if the price rises (plus any dividends). **Maximum loss is limited to the amount invested** because a share price cannot fall below **$0** (a 100% loss at most).[^50] This is the default retail position and the symmetric counterpart to shorting.

### Short selling & the locate/borrow requirement
<a id="short"></a>

**Short selling** = selling shares you **borrow** (from your broker, sourced from its inventory, other clients' margin accounts, or another lender), hoping to **buy them back ("cover") lower** and return them — you profit from a **decline**.[^51] It **requires a margin account** (Regulation T sets a short's initial margin at **150%** of the position value — the proceeds plus 50% of the value from your own funds).[^51]

**The asymmetric risk is load-bearing:** a short's gains are **capped** (the stock can only fall to $0) but **losses are theoretically UNLIMITED** because the price can rise without bound — **you can lose more than the proceeds you received.**[^51] Short example: short 100 shares at $10 ($1,000 proceeds); if it rises to $100 you spend $10,000 to cover — a $9,000 loss, nine times the proceeds.[^51]

**The locate/borrow requirement (Regulation SHO):** **Rule 203(b)(1)** bars a broker from accepting a short-sale order unless it has borrowed the security or has **"reasonable grounds to believe the security can be borrowed"** and delivered by settlement — the **"locate,"** documented before the sale.[^52] Brokers keep an **"easy-to-borrow"** list (low/zero fee) vs **"hard-to-borrow" (HTB)** names that carry a **borrow fee / stock-loan rate** (dynamic; HTB names can run very high *(as of 2026, illustrative)* and reprice sharply).[^53] The lender can **recall** the shares anytime, forcing a **buy-in** at an unfavorable price.[^53] **Naked short selling** (selling without a locate) is **generally prohibited** and produces **failures to deliver (FTDs)**; Reg SHO's close-out rule (Rule 204) and the **threshold-securities** list address persistent fails — though not every transient FTD is itself a violation.[^52] Unique short costs also include **margin interest** and **paying the lender any dividends** ("in lieu of" payments) on the borrowed stock.[^51] **Rule 201 (the alternative uptick rule)** is a short-sale circuit breaker: if a stock falls **≥10% intraday from the prior close**, short sales for the rest of that day and all the next may only execute **above the national best bid**.[^52]

> **Disconfirming note:** despite its reputation, the academic consensus is that short selling **improves price discovery and market efficiency** (short-sale bans demonstrably hurt it); short sellers are often informed traders correcting overvaluation.[^54] Separately, short selling is **hard for retail** — you can be right on the thesis yet lose everything if you are early (borrow/margin costs accrue and a squeeze can force you out).[^51]

*(As of 2026)* the SEC's institutional short-position reporting regime (**Rule 13f-2 / Form SHO**, adopted 2023) has been **repeatedly delayed and is not yet in effect** (compliance extended to January 2028 after a 2025 court remand) — re-verify, as this is the most volatile fact in this section.[^55]

### Short interest & the short squeeze
<a id="squeeze"></a>

**Short interest** = shares currently sold short and not yet covered; often expressed as **% of float** (>20% is treated as very high) or as **days-to-cover / short-interest ratio** (short interest ÷ average daily volume).[^56] FINRA publishes aggregate short interest **twice monthly**.[^56] High short interest signals bearish sentiment **and** squeeze potential.

A **short squeeze** is a self-reinforcing loop: a heavily shorted stock rises, **forcing shorts to cover (buy back)** to limit losses or meet margin calls, which pushes the price **higher**, forcing more covering — intensified by rising borrow fees, recalls, and margin calls.[^57] The canonical retail example is **GameStop (GME), January 2021**: short interest reached **~140% of float** (possible because lent shares are re-lent in a chain), Reddit-driven retail buying drove GME from ~$17 to an intraday **~$483** on Jan 28, 2021, and the loop broke when brokers restricted buying; AMC ran in parallel.[^57] **Contested nuance:** the **SEC's 2021 staff report** concluded the GME rise was sustained mainly by **positive sentiment and direct buying — not primarily short-covering**, a finding some academics dispute — present it as a genuine debate.[^57] A **gamma squeeze** (options dealers delta-hedging short calls, amplifying the move) is **mechanically distinct** and co-occurred with GME — for the options mechanics see `options-trading-and-strategies`.[^57]

### Holding horizon: long-term vs active
<a id="horizon"></a>

Equity holding ranges from **buy-and-hold for years** (low turnover/cost; the default) to **frequent/active trading**.[^58] The **holding period sets the capital-gains line** at a high level: a gain on a position held **one year or less is short-term (ordinary income rates)**; held **more than one year it is long-term (preferential 0/15/20% rates)** — deep mechanics (wash sales, trader-tax-status) → `trading-regulation-compliance-and-taxes`.[^59] Frequent intraday trading triggers the **Pattern-Day-Trader rule** (and its 2026 replacement) — see `trading-risks-and-protections`, not re-derived here.[^58]

---

## 6. Active equity trading vs passive index investing
<a id="active-passive"></a>

This is the spoke's central decision, and the **seam** with the peer skill `investing-and-retirement`.

**The distinction:** **active** equity investing tries to **beat** the market through **security selection** (stock-picking) and **market timing**; **passive** investing buys a broad index fund/ETF and holds it to **match** the market at minimal cost.[^60] Per the SEC, passive "usually translates into less trading … more favorable income-tax consequences … and lower fees."[^60]

**The load-bearing evidence that active stock-picking underperforms** (educational, not advice — but central to this spoke):

- **SPIVA scorecards (S&P Dow Jones Indices)** — a large **majority of active US equity funds underperform their benchmark** over long horizons. *(As of 2026, SPIVA US Year-End 2025)* **~79% of active large-cap funds lagged the S&P 500 in 2025** (the 16th consecutive year the majority lagged), and over **15 years no US equity category had a majority of active managers beat** their benchmark; ~92% lagged over 20 years.[^61] These percentages update twice a year — re-verify against the current scorecard.
- **Morningstar Active/Passive Barometer** — an **independent** study with a different method reaches the same conclusion: only ~21% of active funds survived **and** beat comparable passive funds over the 10 years through 2025.[^62] (This is the strongest rebuttal to "SPIVA is rigged.")
- **Persistence** — top-performing funds rarely stay top: SPIVA's Persistence Scorecard finds top-quartile funds have little better than random odds of repeating, consistent with **luck over skill**.[^63]
- **Why (mechanism):** **Sharpe's "Arithmetic of Active Management"** — before costs the average active dollar earns the market return, so **after costs the average active dollar must underperform** the average passive dollar.[^64] **Bessembinder** — **most individual stocks underperform Treasury bills**, and a tiny fraction (~4% of companies) account for **all** net market wealth creation, so missing the few big winners (easy when stock-picking) badly lags the index that holds them all.[^65]
- **Retail day-trader evidence** — **Barber & Odean** ("Trading Is Hazardous to Your Wealth": the most active households earned ~11.4% vs ~17.9% for the market); **Taiwan** (>80% of day traders lose money, <1% consistently profitable) and **Brazil** (~97% of persistent day traders lost money) studies all converge.[^66]

**The cost angle:** active US-equity funds average ~**0.60%** expense ratio vs ~**0.11%** for passive (asset-weighted) — a structural headwind that **compounds** over decades — plus the **tax drag** of higher turnover (active mutual funds distribute far more capital gains than mostly-passive ETFs).[^67]

**The contrarian view (a genuine debate, not settled):** active can win in **less-efficient** segments — *(as of 2026)* in SPIVA YE-2025 a **majority of active small-cap and emerging-market-debt funds beat** their benchmarks — but the edge is **short-horizon and does not persist** (no category over 15 years).[^68] A separate, unresolved **"passive bubble" / price-discovery** debate (Michael Burry-style critiques vs rebuttals) concerns **second-order market-structure** effects of indexing — keep it separate from the individual-investor question, where the evidence favors low-cost passive.[^69] Rules-based **factor/smart-beta** sits between the two as a middle path.[^68]

**Where the boundary sits (the design rule):** the dividing line is **intent**, not asset class.
- **Stays in this spoke** (intent = actively trade / understand markets): stock selection and equity analysis; trading tactics, frequency, turnover; the **active-vs-passive debate and the underperformance evidence itself** (it explains why most active trading fails); the cost/fee/tax-drag argument as a reason active underperforms; the "passive bubble" market-structure debate.
- **Defers to `investing-and-retirement`** (intent = build long-term retirement wealth passively): **how** to pick/hold a broad index fund or ETF for the long run; **retirement accounts** (401(k)/403(b), Traditional/Roth IRA, contribution limits, vesting); **target-date funds**, dollar-cost averaging **for retirement**, glidepaths; choosing a **robo-advisor or fee-only fiduciary**; expense-ratio-driven **fund selection for a retirement portfolio**.

In one line: **active = trying to *beat* the market (this spoke explains the approach and the sobering evidence); passive buy-and-hold for retirement = *matching* the market cheaply (the peer skill owns the how-to).** This matches the hub's existing SKIP edge to `investing-and-retirement`.

---

## References

Sources are graded by tier: **regulator/SRO** (SEC/Investor.gov, FINRA, Congress.gov CRS, IRS, DTCC), **exchange/index-provider** (NYSE, Nasdaq, S&P Dow Jones Indices, MSCI, FTSE Russell), **peer-reviewed/academic**, and **reputable financial education** (corroborating). Volatile facts are dated *as of 2026* with a "verify current" pointer. The SEC's own pages (sec.gov) intermittently blocked automated fetch during research; in those cases the identical Investor.gov mirror and/or multiple tier-1 law-firm summaries were used and are cited.

[^1]: SEC / Investor.gov — Stocks (common vs preferred, voting, dividend priority, liquidation order). https://www.investor.gov/introduction-investing/investing-basics/investment-products/stocks (regulator)
[^2]: S&P Global, "Practice Essentials: US Preferreds" https://www.spglobal.com/spdji/en/documents/education/practice-essentials-us-preferreds.pdf ; Fidelity, "Convertibles and preferreds" https://www.fidelity.com/learning-center/investment-products/fixed-income-bonds/convertibles-and-preferreds (index-provider; financial education)
[^3]: Council of Institutional Investors, "Dual-Class Stock" https://www.cii.org/dualclass_stock ; Harvard Law School Forum on Corporate Governance, dual-class analysis https://corpgov.law.harvard.edu/2024/10/16/re-thinking-the-hostility-towards-dual-class-share-structures-when-dual-class-shares-work-better/ ; U.S. News, "GOOG vs GOOGL" https://money.usnews.com/investing/articles/goog-vs-googl-stock-difference (institutional body; academic; financial education)
[^4]: Davis Polk, "S&P Dow Jones reopens its indices to companies with multiple share classes" (April 2023 reversal) https://www.davispolk.com/insights/client-update/sp-dow-jones-reopens-its-indices-companies-multiple-share-classes ; Wilson Sonsini, 2017 multi-class index limitation + grandfathering https://www.wsgr.com/en/insights/major-stock-index-providers-to-limit-inclusion-of-multi-class-companies-what-it-means-and-why-it-matters.html (law firm — primary analysis of index-provider policy) [verify as of 2026]
[^5]: NYSE Market Model (DMM/auction) https://www.nyse.com/market-model ; NYSE, "Market Making and the NYSE DMM Difference" https://www.nyse.com/data-insights/market-making-and-the-nyse-dmm-difference ; Akerman, "SEC Approves Changes to NYSE and Nasdaq Minimum Price Rules" (2025–2026) https://www.akerman.com/en/perspectives/sec-approves-changes-to-nyse-and-nasdaq-minimum-price-rules-what-public-companies-need-to-know.html (exchange; law firm) [thresholds volatile as of 2026]
[^6]: S&P Dow Jones Indices, US Indices Methodology (float-cap weighting; committee discretion) https://www.spglobal.com/spdji/en/methodology/article/sp-us-indices-methodology/ ; Corporate Finance Institute, "S&P 500 Index" (committee discretion) https://corporatefinanceinstitute.com/resources/equities/sp-500-index/ (index-provider; financial education) [~503 listings, as of 2026]
[^7]: Nasdaq, official Nasdaq-100 Index Methodology (modified cap weighting, 3× free-float cap, December reconstitution, special rebalance) https://indexes.nasdaq.com/docs/Methodology_NDX.pdf (index-provider, primary methodology)
[^8]: S&P Dow Jones Indices, Dow Jones Averages Methodology (price weighting; the divisor) https://www.spglobal.com/spdji/en/methodology/article/dow-jones-averages-methodology/ ; Stanford SIEPR, "The Dow Jones Average: The Impact of Fixing Its Flaws" (price-weighting criticism) https://siepr.stanford.edu/publications/working-paper/dow-jones-average-impact-fixing-its-flaws (index-provider; academic)
[^9]: FTSE Russell (LSEG), "Four Decades of the Russell Reconstitution" https://www.lseg.com/content/dam/ftse-russell/en_us/documents/research/four-decades-russell-reconstitution.pdf (index-provider) [annual→semi-annual reconstitution change, verify as of 2026]
[^10]: Charles Schwab, "How to read a stock quote" https://www.schwab.com/learn/story/learn-how-to-read-stock-quotes ; WallStreetZen, "How to Read a Stock Quote" https://www.wallstreetzen.com/blog/how-to-read-a-stock-quote/ (financial education)
[^11]: Schwab, "How well do you know market cap?" (FINRA-cited tiers) https://www.schwab.com/learn/story/how-well-do-you-know-market-cap ; Wikipedia, "Market capitalization" / small-cap (no-consensus caveat; Morningstar relative method) https://en.wikipedia.org/wiki/Market_capitalization (financial education; tertiary) [tier bands non-standard, as of 2026]
[^12]: Cornell Law (Wex), "public float" (SEC definition) https://www.law.cornell.edu/wex/public_float ; RBC Direct Investing, "Shares outstanding vs float" https://www.rbcdirectinvesting.com/learn/en/di/hubs/investing-academy/article/shares-outstanding-vs-float/jk0dlma0 (legal reference; financial education)
[^13]: MSCI, GICS https://www.msci.com/indexes/index-resources/gics ; S&P Dow Jones Indices, GICS Methodology https://www.spglobal.com/spdji/en/documents/methodologies/methodology-gics.pdf ; Wikipedia, GICS (11/25/74/163 structure, March 2023 revision) https://en.wikipedia.org/wiki/Global_Industry_Classification_Standard (index-providers; tertiary corroboration) [structure as of 2026]
[^14]: Britannica Money, "P/E ratio" (trailing vs forward) https://www.britannica.com/money/price-to-earnings-ratio ; Corporate Finance Institute, "Price Earnings Ratio" https://corporatefinanceinstitute.com/resources/valuation/price-earnings-ratio/ (financial education)
[^15]: Nasdaq glossary, "earnings yield" https://www.nasdaq.com/glossary/e/earnings-yield ; Corporate Finance Institute, "Earnings Yield" / PEG https://corporatefinanceinstitute.com/resources/valuation/earnings-yield/ (exchange; financial education)
[^16]: Finance Strategists, "Limitations of the P/E Ratio" (value trap, leverage-blind, accounting distortions) https://www.financestrategists.com/wealth-management/fundamental-vs-technical-analysis/limitations-of-the-price-earnings-ratio/ ; New Constructs, "P/E Ratios Are Misleading" https://www.newconstructs.com/why-pe-ratios-are-misleading/ (financial education; research firm)
[^17]: Deloitte, "Roadmap: Earnings per Share" (basic vs diluted, ASC 260 equal prominence) https://dart.deloitte.com/USDART/home/publications/roadmap/earnings-per-share ; Wall Street Prep, "Diluted EPS" https://www.wallstreetprep.com/knowledge/diluted-eps/ (Big-Four accounting authority; financial education)
[^18]: Corporate Finance Institute, "Price to Book Ratio" https://corporatefinanceinstitute.com/resources/valuation/price-to-book-ratio/ ; O'Shaughnessy Asset Management, "Negative Equity, Veiled Value, and the Erosion of Price-to-Book" (intangibles/buyback distortion) https://www.osam.com/Commentary/negative-equity-veiled-value-and-the-erosion-of-price-to-book (financial education; quant research)
[^19]: Harvard Business School Online, "How to Read & Understand a Profit and Loss Statement" (margin ladder) https://online.hbs.edu/blog/post/income-statement-analysis ; Corporate Finance Institute, "Revenue vs Income" https://corporatefinanceinstitute.com/resources/accounting/revenue-vs-income/ (academic; financial education)
[^20]: SmartAsset, "What Is a Good Dividend Payout Ratio?" https://smartasset.com/financial-advisor/dividend-payout-ratio ; Investing.com Academy, "Dividend Coverage Ratio" https://www.investing.com/academy/stock-picking/dividend-coverage-ratio/ (financial education)
[^21]: Morningstar, "How to Avoid Dividend Traps" https://www.morningstar.com/stocks/how-avoid-dividend-traps ; Motley Fool, "Dividend Yield Trap" https://www.fool.com/investing/how-to-invest/dividend-stocks/dividend-yield-trap/ (financial education)
[^22]: Corporate Finance Institute, "Debt to Equity Ratio" + "Balance Sheet" (accounting identity, residual equity, current ratio) https://corporatefinanceinstitute.com/resources/accounting/debt-to-equity-ratio-formula/ ; SEC / Investor.gov, "Beginners' Guide to Financial Statements" https://www.sec.gov/reportspubs/investor-publications/investorpubsbegfinstmtguide (financial education; regulator)
[^23]: SEC / Investor.gov — Form 10-K / Form 10-Q (filing cadence, audited vs unaudited, EDGAR) https://www.investor.gov/introduction-investing/investing-basics/glossary/form-10-k (regulator)
[^24]: Charles Schwab, "Forward Guidance" / earnings https://www.schwab.com/learn/story/what-to-know-earnings-season ; Britannica Money, "forward guidance" https://www.britannica.com/money/forward-guidance ; TipRanks, "Analyst ratings" glossary https://www.tipranks.com/glossary/a/analyst-ratings (financial education) [guidance-move magnitudes illustrative]
[^25]: Ryan O'Connell (CFA), "Growth vs Value Investing" + "Value Investing" https://www.ryanoconnellfinance.com/ ; VanEck, "Dividend investing" https://www.vaneck.com/us/en/blogs/income-investing/ ; S&P Dow Jones Indices, Dividend Aristocrats methodology (25-year criterion) https://www.spglobal.com/spdji/en/indices/dividends-factors/sp-500-dividend-aristocrats/ (CFA practitioner; asset manager; index-provider) [Aristocrat count ~69, as of 2026]
[^26]: Vanguard, value-vs-growth research https://corporate.vanguard.com/content/corporatesite/us/en/corp/articles/value-vs-growth-stocks.html ; J.P. Morgan Asset Management, "Guide to the Markets" value/growth https://am.jpmorgan.com/us/en/asset-management/adv/insights/market-insights/guide-to-the-markets/ ; Morningstar, value/growth rotation https://www.morningstar.com/markets (asset-manager research; financial education) [year-by-year leadership volatile, as of 2026]
[^27]: SEC / Investor.gov — "Ex-Dividend Dates: When Are You Entitled to Stock and Cash Dividends" (the four dates; buy-before-ex rule) https://www.investor.gov/introduction-investing/investing-basics/glossary/ex-dividend-dates-when-are-you-entitled-stock-and ; IRS Topic No. 404, Dividends (qualified vs ordinary, holding-period test) https://www.irs.gov/taxtopics/tc404 (regulator; IRS)
[^28]: DTCC, "T+1 Dividend Processing FAQ" (ex-date = record date under T+1) https://www.dtcc.com/-/media/Files/PDFs/T2/T1-Dividend-Processing-FAQ.pdf ; Charles Schwab, "Ex-Dividend Dates and Dividend Risk" (approximate price drop) https://www.schwab.com/learn/story/ex-dividend-dates-understanding-dividend-risk (clearing utility — primary; financial education) [T+1 effective May 28, 2024; as of 2026]
[^29]: NYU Stern / Damodaran, "Dividend Irrelevance" (Modigliani–Miller) https://pages.stern.nyu.edu/~adamodar/New_Home_Page/invfables/dividirrelevance.htm (academic)
[^30]: Charles Schwab, "How a Dividend Reinvestment Plan (DRIP) Works" https://www.schwab.com/learn/story/how-dividend-reinvestment-plan-works ; Motley Fool, "Dividend Reinvestment" https://www.fool.com/investing/stock-market/types-of-stocks/dividend-stocks/dividend-reinvestment/ (financial education)
[^31]: Interactive Brokers, "Dividend Types" (special, stock vs cash) https://www.interactivebrokers.com/campus/trading-lessons/dividend-types/ ; SmartAsset, "Stock Dividend vs Cash Dividend" https://smartasset.com/financial-advisor/stock-dividend-vs-cash-dividend (broker; financial education)
[^32]: Skadden, "Revisiting Share Repurchases" (OMR/tender/ASR; Rule 10b-18) https://www.skadden.com/insights/publications/2022/03/revisiting-share-repurchases-in-volatile-times ; CFA Institute Research Foundation, "Stock Buybacks" literature review (buybacks vs dividends; EPS mechanics; signaling) https://rpc.cfainstitute.org/sites/default/files/-/media/documents/book/rf-lit-review/2022/rflr-stock-buybacks.pdf (law firm; academic/industry body)
[^33]: Congressional Research Service (Congress.gov), Report R47397, "The 1% Excise Tax on Stock Repurchases" (rate, netting, non-deductibility, Jan 1 2023 effective) https://www.congress.gov/crs-product/R47397 ; Sullivan & Cromwell, "Stock Buyback Tax" final regulations (finalized Nov 21, 2025 at 1%) https://www.sullcrom.com/insights/memo/2025/November/Stock-Buyback-Tax (CRS — primary/non-partisan; law firm) [1% confirmed; 4% proposal tentative; as of 2026]
[^34]: Harvard Law School Forum on Corporate Governance, "Share Buybacks and Executive Compensation: Assessing Key Criticisms" https://corpgov.law.harvard.edu/2023/07/16/share-buybacks-and-executive-compensation-assessing-key-criticisms/ (academic)
[^35]: SoFi, "What Is a Stock Split?" https://www.sofi.com/learn/content/what-is-a-stock-split/ ; Britannica Money, "Fractional share investing" (undercuts split rationale) https://www.britannica.com/money/fractional-share-investing ; Bloomberg, "What is a stock split" (NVDA/AVGO/CMG 2024) https://www.bloomberg.com/news/articles/2024-06-01/what-is-a-stock-split-what-to-know-about-nvidia-s-nvda-split (financial education; reputable media)
[^36]: ScienceDirect, "Do stock splits signal undervaluation?" (announcement effect = signaling, not value) https://www.sciencedirect.com/science/article/abs/pii/S2214635016000113 (peer-reviewed)
[^37]: Skadden, "SEC Approves Nasdaq Rule Change" (2024–2025 reverse-split restrictions; $1 min bid, 30-day trigger, 10-day cure) https://www.skadden.com/insights/publications/2024/11/sec-approves-nasdaq-rule-change ; K&L Gates, "The State of Play for Reverse Stock Splits by Nasdaq- and NYSE-Listed Issuers" https://www.klgates.com/The-State-of-Play-for-Reverse-Stock-Splits-by-Nasdaq-and-NYSE-Listed-Issuers-1-23-2025 (law firms) [2025 rule change; as of 2026]
[^38]: SEC / Investor.gov — IPO glossary; gun-jumping/quiet period (Securities Act §5) https://www.investor.gov/introduction-investing/investing-basics/glossary/initial-public-offering-ipo ; Cornell Law (Wex), "gun jumping" https://www.law.cornell.edu/wex/gun_jumping (regulator; legal reference)
[^39]: Corporate Finance Institute, "Over-Allotment (Greenshoe) Option" (up to 15%, stabilization) https://corporatefinanceinstitute.com/resources/equities/overallotment-greenshoe-option-ipo/ (financial education)
[^40]: Fidelity, "IPO share allocation process" (~90% institutional / ~10% retail; asset minimums; no guarantee) https://www.fidelity.com/learning-center/trading-investing/trading/ipo-share-allocation-process ; Charles Schwab, "How IPO Shares Are Priced and Allotted" https://www.schwab.com/learn/story/getting-slice-how-ipo-shares-are-priced-and-allotted (broker/financial education)
[^41]: FINRA Rule 5131 (new-issue allocations, anti-spinning) https://www.finra.org/rules-guidance/rulebooks/finra-rules/5131 ; InnReg, "FINRA Rule 5130 Explained" (restricted persons) https://www.innreg.com/resources/finra-rules/5130-restrictions-on-the-purchase-and-sale-of-initial-equity-public-offerings (SRO — primary; financial education) [Apr 2026 FINRA 5130/5131 amendments proposed — verify]
[^42]: Jay R. Ritter (University of Florida), IPO Data / "Initial Public Offerings: Underpricing" (long-run ~19% first-day average; 3-year underperformance) https://site.warrington.ufl.edu/ritter/ipo-data/ (academic — authoritative IPO dataset) [first-day-return figures volatile, as of 2026]
[^43]: Corporate Finance Institute, "Direct Listing — Overview, Pros/Cons" https://corporatefinanceinstitute.com/resources/equities/direct-listing/ ; Fenwick, "SEC Approves Nasdaq Rule Change Allowing Direct Listings with a Capital Raise" https://www.fenwick.com/insights/publications/sec-approves-nasdaq-rule-change-allowing-direct-listings-with-a-capital-raise (financial education; law firm)
[^44]: Covington, "SEC Approves NYSE Rule Allowing Primary Issuances In Direct Listings" (NYSE Dec 2020) https://www.cov.com/en/news-and-insights/insights/2021/01/sec-approves-nyse-rule-allowing-primary-issuances-in-direct-listings ; Paul Weiss, "SEC Approves Nasdaq Rule Change Permitting Primary Direct Floor Listings" (Nasdaq May 2021) https://www.paulweiss.com/practices/transactional/capital-markets/publications/sec-approves-nasdaq-rule-change-permitting-primary-direct-floor-listings?id=40140 (law firms) [primary-DL low adoption, verify as of 2026]
[^45]: SEC / Investor.gov, "What You Need to Know About SPACs — Updated Investor Bulletin" (lifecycle, trust, redemption, ~20% promote, warrants, sponsor-conflict caution) https://www.investor.gov/introduction-investing/general-resources/news-alerts/alerts-bulletins/investor-bulletins/what-you (regulator)
[^46]: Yale Journal on Regulation, "Was the SPAC Crash Predictable?" (boom/bust, post-merger underperformance, dilution) https://www.yalejreg.com/bulletin/was-the-spac-crash-predictable/ ; Morgan Lewis, "A Look at the SPAC Market in 2022" https://www.morganlewis.com/pubs/2022/04/a-look-at-the-spac-market-in-2022 (academic; law firm) [2020-22 figures historical, as of 2026]
[^47]: Davis Polk, "SEC adopts final SPAC rules" (Jan 2024: disclosure, co-registrant/§11 liability, projections, PSLRA safe-harbor removed, underwriter-liability guidance not a rule) https://www.davispolk.com/insights/client-update/sec-adopts-final-spac-rules ; Ropes & Gray, "SEC Adopts New Rules Regarding SPAC Transactions" https://www.ropesgray.com/en/insights/alerts/2024/02/sec-adopts-new-rules-regarding-spac-transactions (law firms) [SEC 2024 SPAC rules; verify effective dates as of 2026]
[^48]: Cornell Law (Wex), "follow-on offering," "secondary offering," "secondary market" (dilutive vs non-dilutive; terminology trap) https://www.law.cornell.edu/wex/follow-on_offering ; Perkins Coie, "Public Company Handbook Ch.12: Follow-On Offerings and Shelf Registrations" (S-3 shelf, ATM) https://perkinscoie.com/public-company-handbook-chapter-12-follow-offerings-and-shelf-registrations (legal reference; law firm)
[^49]: Wikipedia, "Lock-up period" (90–180 days, 180 standard, contractual not statutory) https://en.wikipedia.org/wiki/Lock-up_period ; SoFi, "Navigating the IPO Lock-Up Period" (expiration supply event) https://www.sofi.com/learn/content/navigating-ipo-lock-up-period/ (tertiary; financial education)
[^50]: SEC / Investor.gov, "Introduction to Short Sales" (long-position loss limited to amount invested, contrast with unlimited short loss) https://www.investor.gov/introduction-investing/general-resources/news-alerts/alerts-bulletins/investor-bulletins-51 (regulator)
[^51]: SEC / Investor.gov, "Introduction to Short Sales" (borrow/cover mechanics, margin account, 150% Reg T, unlimited loss, dividends owed to lender) https://www.investor.gov/introduction-investing/general-resources/news-alerts/alerts-bulletins/investor-bulletins-51 ; Bankrate, "Short Selling: How To Short Sell Stocks" (worked loss example) https://www.bankrate.com/investing/short-selling-how-to-short-a-stock/ (regulator; financial education)
[^52]: SEC, Regulation SHO / short-sale rules (Rule 203(b)(1) locate; Rule 204 close-out; threshold securities; Rule 201 alternative uptick) https://www.sec.gov/rules-regulations/staff-guidance/trading-markets-frequently-asked-questions-7 ; Congressional Research Service, "Short Selling: Background and Policy Issues" (IF12400) https://www.congress.gov/crs-product/IF12400 (regulator; CRS — primary)
[^53]: Interactive Brokers Campus, "The Risks of Shorting: Borrow Fees" (ETB/HTB, stock-loan rate, recall/buy-in) https://www.interactivebrokers.com/campus/traders-insight/securities/short-selling/the-risks-of-shorting-series-part-ii-borrow-fees/ (broker/financial education) [borrow-fee levels illustrative, as of 2026]
[^54]: Managed Funds Association, "An Introduction to Short Selling" (price-discovery benefit; ban harms) https://www.mfaalts.org/wp-content/uploads/2025/04/MFA-Updated-Intro-to-Short-Selling-Research-Paper.pdf ; NBER, "Short-Sale Constraints and Overpricing" https://www.nber.org/reporter/spring05/short-sale-constraints-and-overpricing (industry body; academic)
[^55]: Morgan Lewis, "Short-Sale Reporting on Form SHO: Compliance Date Further Extended to 2028" (Rule 13f-2 not yet in effect) https://www.morganlewis.com/pubs/2025/12/short-sale-reporting-on-form-sho-compliance-date-further-extended-to-2028 (law firm) [13f-2/Form SHO delayed to Jan 2028; most volatile fact, verify as of 2026]
[^56]: FINRA, Short Interest Reporting (bi-monthly cadence) https://www.finra.org/filing-reporting/regulatory-filing-systems/short-interest ; Charles Schwab, short-interest / days-to-cover https://www.schwab.com/learn/story/what-is-short-interest (SRO — primary; financial education)
[^57]: SEC, "Staff Report on Equity and Options Market Structure Conditions in Early 2021" (GME; concluded sentiment/buying, not primarily short-covering) https://www.sec.gov/files/staff-report-equity-options-market-struction-conditions-early-2021.pdf ; Columbia Law School Blue Sky Blog, "An Academic Critique of the SEC's GameStop Report" https://clsbluesky.law.columbia.edu/2022/02/22/an-academic-critique-of-the-secs-gamestop-report/ (regulator — primary; academic critique) [GME attribution contested]
[^58]: Charles Schwab, buy-and-hold vs active trading https://www.schwab.com/learn/story/stock-trading-vs-investing-whats-difference ; FINRA, margin/pattern-day-trader key topic https://www.finra.org/rules-guidance/key-topics/margin-accounts (financial education; SRO)
[^59]: IRS Topic No. 409, "Capital Gains and Losses" (short-term ≤1 yr ordinary; long-term >1 yr preferential 0/15/20%) https://www.irs.gov/taxtopics/tc409 (IRS — primary)
[^60]: SEC / Investor.gov, "Mutual Funds and ETFs" (active vs passive definitions; passive = lower trading/tax/fees) https://www.investor.gov/introduction-investing/investing-basics/investment-products/mutual-funds-and-exchange-traded-funds-etfs/mutual-funds (regulator)
[^61]: S&P Dow Jones Indices, SPIVA US Scorecard (Year-End 2025; majority of active large-cap funds lag the S&P 500; no category beats over 15 yr) https://www.spglobal.com/spdji/en/spiva/article/spiva-us/ (index-provider — primary) [percentages update twice yearly; as of 2026]
[^62]: Morningstar, "US Active/Passive Barometer" (independent; ~21% of active funds survived and beat passive over 10 yr through 2025) https://www.morningstar.com/business/insights/research/active-passive-barometer (independent benchmark study) [as of 2026]
[^63]: S&P Dow Jones Indices, "SPIVA US Persistence Scorecard" (top performers rarely persist; luck over skill) https://www.spglobal.com/spdji/en/spiva/article/us-persistence-scorecard/ (index-provider — primary)
[^64]: William F. Sharpe, "The Arithmetic of Active Management" (Financial Analysts Journal, 1991) https://web.stanford.edu/~wfsharpe/art/active/active.htm (peer-reviewed)
[^65]: Hendrik Bessembinder, "Do Stocks Outperform Treasury Bills?" (Journal of Financial Economics, 2018; most stocks lag T-bills; ~4% of firms create all net wealth) https://papers.ssrn.com/sol3/papers.cfm?abstract_id=2900447 (peer-reviewed)
[^66]: Barber & Odean, "Trading Is Hazardous to Your Wealth" (Journal of Finance, 2000) http://faculty.haas.berkeley.edu/odean/papers/returns/individual_investor_performance_final.pdf ; Chague, De-Losso & Giovannetti, "Day Trading for a Living?" (Brazil; ~97% of persistent day traders lose) https://papers.ssrn.com/sol3/papers.cfm?abstract_id=3423101 (peer-reviewed/academic)
[^67]: Morningstar, "US Fund Fee Study" (active ~0.59–0.60% vs passive ~0.11%, asset-weighted) https://www.morningstar.com/business/insights/blog/funds/us-fund-fee-study ; SEC, "How Fees and Expenses Affect Your Investment Portfolio" https://www.sec.gov/investor/alerts/ib_fees_expenses.pdf (research; regulator) [fee figures 2024 study; verify as of 2026]
[^68]: S&P Dow Jones Indices SPIVA US YE-2025 (active majority beat in small-cap & EM debt short-horizon; none over 15 yr) https://www.spglobal.com/spdji/en/spiva/article/spiva-us/ (index-provider — primary) [segment figures volatile, as of 2026]
[^69]: ETF Strategy, "Passive investing a 'bubble', says Michael Burry" https://www.etfstrategy.com/passive-investing-a-bubble-says-big-short-investor-michael-burry-10449/ ; Elm Wealth, "A Few Quick Responses to Michael Burry's Index-bubble Remarks" (rebuttal) https://elmwealth.com/michael-burry-index-bubble/ (reputable financial commentary — both sides of an unresolved debate)
