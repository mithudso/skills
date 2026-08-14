<!-- Provenance: reference under the `options-trading-and-strategies` spoke (trading-and-investing hub family). Created 2026-06-16 via /dr deep-research. Educational only — NOT financial/investment/tax advice. -->

# Options Fundamentals & Contract Mechanics

US listed equity/ETF and index options, for a retail participant, **as of 2026**. Educational only — **not advice**. This domain is largely **stable**; volatile items (0DTE growth, product launches) are flagged *as of 2026*.

## Contents

- [Calls vs puts; the buyer/writer split](#calls-vs-puts-the-buyerwriter-split)
- [Contract terms & the 100-share multiplier](#contract-terms--the-100-share-multiplier)
- [Intrinsic vs extrinsic (time) value](#intrinsic-vs-extrinsic-time-value)
- [Moneyness (ITM / ATM / OTM)](#moneyness-itm--atm--otm)
- [Exercise style: American vs European](#exercise-style-american-vs-european)
- [The OCC, exercise & assignment](#the-occ-exercise--assignment)
- [The options chain](#the-options-chain)
- [Open interest vs volume](#open-interest-vs-volume)
- [Contract tenors: LEAPS, weeklies, 0DTE](#contract-tenors-leaps-weeklies-0dte)
- [Settlement: shares vs cash](#settlement-shares-vs-cash)
- [References](#references)

## Calls vs puts; the buyer/writer split

A **call** gives its buyer the **right, not the obligation, to buy** 100 shares of the underlying at the strike price on or before expiration. A **put** gives the buyer the **right, not the obligation, to sell** 100 shares at the strike.[^sec-bulletin][^finra] There are four basic positions: long call, short call, long put, short put.

- The **buyer / holder (long side)** pays the **premium** up front and holds the right. **Maximum loss for a long option = the premium paid** (you cannot lose more, no matter how far the underlying moves against you).[^sec-bulletin]
- The **seller / writer (short side)** **receives the premium** and takes on the **obligation**: a short **call** writer is obligated to **sell/deliver** shares if assigned; a short **put** writer is obligated to **buy/take delivery of** shares if assigned. The premium received is the writer's maximum profit; the writer keeps it whether or not the option is exercised — but the *loss* can be far larger than that premium (see `strategies-and-risk.md` on naked options).[^sec-bulletin]

## Contract terms & the 100-share multiplier

- **Strike price** — the fixed price at which the holder may buy (call) or sell (put) the underlying.[^sec-bulletin]
- **Expiration date** — the last date the option can be exercised. The standard convention for monthly options is the **third Friday of the expiration month**, which is the last trading day.[^cboe-weeklys][^oic-exercise] *(Historical note: older SEC/OCC text phrased expiration as "the Saturday after the third Friday," reflecting a legacy OCC technical expiration date; in practice trading stops on the **third Friday** and that is the operative last-trading-day for retail.)*
- **Premium** — the price per share paid for the contract. Because one standard contract covers 100 shares, **total cost = quoted premium × 100**.[^sec-bulletin][^oic-exercise]
- **Contract multiplier** — the US-listed standard is **100 shares per equity/ETF option contract**, enforced by the OCC and the exchanges.[^occ-spec] **Caveat:** corporate actions (splits, mergers, stock dividends, spinoffs) can create **adjusted (non-standard) contracts** that deliver something other than 100 shares of a single security.[^oic-adjust]

> **Worked example.** A call quoted at **$2.50** costs **$2.50 × 100 = $250** for one contract (plus any fees). If you sell one covered call at **$2.50**, you collect **$250**.

## Intrinsic vs extrinsic (time) value

Premium decomposes as: **Premium = Intrinsic value + Extrinsic (time) value.**[^fidelity-pricing]

- **Intrinsic value** = how far in-the-money the option is. Call intrinsic = max(stock − strike, 0). Put intrinsic = max(strike − stock, 0). **Out-of-the-money options have zero intrinsic value** — they are *all* extrinsic.[^fidelity-pricing][^cme-moneyness]
- **Extrinsic (time) value** = Premium − Intrinsic value; it reflects the time remaining and the implied volatility priced in. (The *drivers* — time decay and volatility — are owned by `the-greeks.md` and `volatility-and-pricing.md`; this page just names the decomposition.)

> **Worked example.** A **call** struck at **$50** with the stock at **$55** has **$5 intrinsic value/share = $500 per contract**. If that call trades at **$7.50**, the extra **$2.50/share ($250)** is extrinsic (time) value. A **put** struck at **$50** with the stock at **$45** likewise has **$5 intrinsic ($500/contract)**.

## Moneyness (ITM / ATM / OTM)

Moneyness describes the strike relative to the underlying price. It is a property of the **contract**, not a statement about whether your *trade* is profitable or which side you are on.[^cme-moneyness][^macroption-money]

| | Call | Put |
| --- | --- | --- |
| **In-the-money (ITM)** | stock **above** strike | stock **below** strike |
| **At-the-money (ATM)** | stock ≈ strike | stock ≈ strike |
| **Out-of-the-money (OTM)** | stock **below** strike | stock **above** strike |

Only ITM options carry intrinsic value.

## Exercise style: American vs European

- **American-style** — exercisable **any** market day up to and including expiration.[^oic-exercise]
- **European-style** — exercisable **only at expiration** (the last trading day).[^oic-exercise]
- **Which is which (US):** essentially **all listed single-name equity and ETF options are American-style and physically (share) settled**. Most **cash-settled broad-based index options (e.g., SPX, and the mini XSP)** are **European-style**.[^cboe-eu][^oic-exercise] European/cash settlement on SPX **eliminates early-assignment risk**.[^cboe-eu]
- **Disconfirming check — "are all equity options American?"** Single-stock/ETF options: yes, American. But index options are a different product and are typically European, so the blanket claim "all options are American" is **false**. (Equity *weeklies* are still American and physically settled.)[^cboe-eu][^oic-exercise]

## The OCC, exercise & assignment

- **The OCC (Options Clearing Corporation)** is the central clearinghouse and guarantor for all US listed options; it processes exercises and assignments and stands between buyer and seller so neither relies on the other's creditworthiness.[^occ-spec][^oic-exercise]
- **Exercise (long):** the holder gives **irrevocable instructions** to exercise; an American holder can do this on any open market day. Exercising a call makes you a stockowner (you must have the capital or margin to buy the shares); exercising a put sells your shares at the strike.[^oic-exercise]
- **Assignment (short):** the OCC **randomly assigns** exercise notices to clearing-member firms; each firm then assigns one of its short customers (on a random or otherwise fair, documented basis). You cannot predict *which* short position gets assigned.[^oic-exercise][^thismatter]
- **Exercise-by-exception ("Ex-by-Ex") / auto-exercise:** at expiration the OCC **automatically exercises options that are ITM by the threshold amount — $0.01 per contract** — unless the clearing member submits contrary instructions.[^oic-exercise][^thismatter] **Nuance:** because members (and customers, via their broker) can override with do-not-exercise instructions, it is **not strictly "automatic"** — hence the name "exercise by exception," and an individual broker's threshold may differ from the OCC's $0.01.[^oic-exercise]

> **Practical cost.** If a long ITM option auto-exercises into a stock position you cannot fund — or a short is assigned and you lack the shares/cash — the broker typically **force-liquidates at the next open using market orders**, often at poor prices. Holding ITM options into expiration without the capital or instructions to manage them is a real, avoidable cost.[^ibkr-exercise]

## The options chain

An **options chain** is the grid of all listed contracts for an underlying, organized by **expiration date** and **strike price**, conventionally with **calls on one side and puts on the other**.[^sofi-oi] Core columns: **Bid** (highest price a buyer will pay / what you receive when selling), **Ask** (lowest price a seller will accept / what you pay when buying), **Last** (most recent trade), **Volume** (today's contracts traded), and **Open Interest** (total outstanding contracts).

## Open interest vs volume

| Metric | What it measures | Resets? |
| --- | --- | --- |
| **Volume** | Contracts traded **during the current session** | Yes — resets to zero each day |
| **Open interest (OI)** | Total **outstanding** contracts not yet closed, exercised, or expired | No — updates once after the close, moves slower |

OI rises by one when a buyer and seller **both open** a new contract; it falls when both **close**. **Volume can exceed OI** when the same contracts are traded back and forth intraday.[^sofi-oi][^oa-oivol] Both correlate with liquidity (see `strategies-and-risk.md`), but they are **distinct metrics** — conspicuously so for 0DTE, where huge same-day volume produces little or no *overnight* open interest.[^spotgamma-0dte]

## Contract tenors: LEAPS, weeklies, 0DTE

- **LEAPS (Long-Term Equity AnticiPation Securities):** options with expirations **more than one year out** (typically up to ~2–3 years). Mechanically identical to standard options (100-share multiplier, calls/puts); they carry **more extrinsic value** because of the long tenor, and **rho** matters more for them (see `the-greeks.md`).[^oic-leaps][^fool-leaps]
- **Weeklies (Cboe's "Weeklys"):** short-dated series typically **listed Thursday and expiring the following Friday**, filling the weeks around standard third-Friday expirations. Equity/ETF weeklies are **American-style and physically settled**.[^cboe-weeklys]
- **0DTE (zero-days-to-expiration):** an option **on its final trading day** (0 days left). For products with **daily** expirations — notably **SPX**, plus other major indices/ETFs — there is a 0DTE contract on nearly every trading day. **Disconfirming nuance:** large single stocks **without** daily expirations (e.g., NVDA, TSLA, AAPL) do **not** get a fresh daily 0DTE; for them "0DTE" just means the **nearest weekly Friday** on its expiration day. Daily 0DTE is an index/major-ETF phenomenon.[^spotgamma-0dte]

> **0DTE / SPX growth — as of 2026 (volatile, verify):** 0DTE options were roughly **a quarter of all US listed options volume in 2025** (up from ~21–22% in 2024), and for **SPX specifically** 0DTE has grown to a **majority of SPX volume** (reported in the ~47–59% range depending on source/period, averaging on the order of ~2M+ contracts/day).[^cboe-state][^vol-report][^oa-0dte] Treat the precise percentage as **directional and as-of-2026**, not a fixed constant. Cboe's own research argues 0DTE's rise has *not* materially increased realized S&P 500 volatility — note that is the exchange's framing of a publicly contested debate.[^cboe-impact]

## Settlement: shares vs cash

- **Equity & ETF options: physically (share) settled** — exercise/assignment delivers the actual 100 shares (a buyer of an exercised call receives shares and pays strike × 100; an assigned call writer delivers shares).[^occ-spec][^oic-exercise]
- **Most broad-based index options (SPX, XSP, etc.): cash-settled** — at expiration the ITM amount is paid/received in **cash** (there is no deliverable index "share").[^cboe-eu]
- **Disconfirming check — "do options always settle in shares?"** No — index options cash-settle; only equity/ETF options deliver shares.[^cboe-eu]

## References

Authoritative tiers first (SEC/Investor.gov, FINRA, OCC, OIC, Cboe); brokerage and education sources corroborate. Volatile facts dated *as of 2026*.

[^sec-bulletin]: SEC / Investor.gov, Investor Bulletin: An Introduction to Options, https://www.investor.gov/introduction-investing/general-resources/news-alerts/alerts-bulletins/investor-bulletins-63 — regulator. verified-as-of: 2026-06-16
[^finra]: FINRA, "Options" investor overview, https://www.finra.org/investors/investing/investment-products/options — regulator. verified-as-of: 2026-06-16
[^occ-spec]: OCC, Equity Options Product Specifications (100-share unit, deliverable, physical settlement), https://www.theocc.com/clearance-and-settlement/clearing/equity-options-product-specifications — clearing/docs. verified-as-of: 2026-06-16
[^oic-exercise]: OIC / Options Industry Council, Options Exercise & Assignment FAQ (American vs European, random assignment, exercise-by-exception / $0.01 ITM), https://www.optionseducation.org/referencelibrary/faq/options-exercise — authoritative options education. verified-as-of: 2026-06-16
[^oic-adjust]: OIC, "Splits, Mergers, Spinoffs & Bankruptcies" (adjusted contracts), https://www.optionseducation.org/referencelibrary/faq/splits-mergers-spinoffs-bankruptcies — options education. verified-as-of: 2026-06-16
[^fidelity-pricing]: Fidelity, "Understanding options pricing" (intrinsic vs extrinsic, worked example), https://www.fidelity.com/learning-center/trading-investing/understanding-options-pricing — brokerage education. verified-as-of: 2026-06-16
[^cme-moneyness]: CME Group, "Calculating Options Moneyness and Intrinsic Value," https://www.cmegroup.com/education/courses/introduction-to-options/calculating-options-moneyness-and-intrinsic-value — exchange education. verified-as-of: 2026-06-16
[^macroption-money]: Macroption, "Option Moneyness" (ITM/ATM/OTM for calls and puts), https://www.macroption.com/option-moneyness/ — education. verified-as-of: 2026-06-16
[^cboe-eu]: Cboe, Mini-SPX (XSP) / SPX European-style, cash-settled product pages and fact sheet, https://www.cboe.com/tradable_products/sp_500/mini_spx_options/european_style/ — exchange. verified-as-of: 2026-06-16
[^thismatter]: thismatter.com, "Options Exercise, Assignment, and Settlement" (OCC, random assignment, ex-by-ex), https://thismatter.com/money/options/trading-exercise-assignment.htm — education (corroborating). verified-as-of: 2026-06-16
[^ibkr-exercise]: Interactive Brokers, "Exercise and Assignment" trading lesson, https://www.interactivebrokers.com/campus/trading-lessons/exercise-and-assignment/ — brokerage education. verified-as-of: 2026-06-16
[^sofi-oi]: SoFi, "Open Interest in Options" (chain layout; OI vs volume), https://www.sofi.com/learn/content/open-interest-options/ — education. verified-as-of: 2026-06-16
[^oa-oivol]: Option Alpha, "Options Volume vs Open Interest," https://optionalpha.com/learn/options-volume-vs-open-interest — education. verified-as-of: 2026-06-16
[^spotgamma-0dte]: SpotGamma, 0DTE volume/open-interest relationship and the single-stock 0DTE nuance, https://support.spotgamma.com/hc/en-us/articles/15350782964755 — options-analytics education. verified-as-of: 2026-06-16
[^oic-leaps]: OIC, "LEAPS Overview," https://www.optionseducation.org/optionsoverview/leaps-overview — options education. verified-as-of: 2026-06-16
[^fool-leaps]: The Motley Fool, "LEAPS" definition, https://www.fool.com/terms/l/leaps/ — education (corroborating). verified-as-of: 2026-06-16
[^cboe-weeklys]: Cboe, "Available Weeklys" (listing/expiration rules; American, physically settled), https://www.cboe.com/available_weeklys — exchange. verified-as-of: 2026-06-16
[^cboe-state]: Cboe Insights, "The State of the Options Industry 2025" (0DTE share of volume), https://www.cboe.com/insights/posts/the-state-of-the-options-industry-2025 — exchange. verified-as-of: 2026-06-16
[^vol-report]: Traders Magazine, VOL Report citing Cboe 2025 figures, https://www.tradersmagazine.com/vol-report/vol-report-0dte-flex-options-are-2025-heroes/ — industry (corroborating Cboe). verified-as-of: 2026-06-16
[^oa-0dte]: Option Alpha, "The Rise of SPX 0DTE Trading — Analyzing Volume Trends," https://optionalpha.com/blog/the-rise-of-spx-0dte-trading-analyzing-volume-trends — education. verified-as-of: 2026-06-16
[^cboe-impact]: Cboe Insights, "Evaluating the Market Impact of SPX 0DTE Options" (exchange's view that realized volatility has not significantly risen), https://www.cboe.com/insights/posts/volatility-insights-evaluating-the-market-impact-of-spx-0-dte-options/ — exchange (contested public debate; flagged as exchange perspective). verified-as-of: 2026-06-16
