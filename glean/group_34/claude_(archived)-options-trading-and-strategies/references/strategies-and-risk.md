<!-- Provenance: reference under the `options-trading-and-strategies` spoke (trading-and-investing hub family). Created 2026-06-16 via /dr deep-research. Educational only — NOT financial/investment/tax advice. -->

# Core Options Strategies & Their Payoff / Risk Profiles

US listed equity options, educational, **as of 2026** — **not advice**. Standard contract = 100 shares. "Breakeven" is **at expiration**. Strategy mechanics are **stable**; the figures below are structural, not market-dependent. Broker-specific items (approval tiers, margin formulas, commissions) and the PDT-rule change are **volatile** and flagged.

> **Risk lead.** A **long** option's worst case is the premium paid (you cannot lose more). A **naked short call**'s loss is **theoretically unlimited**; a **naked short put**'s loss is large but bounded (down to a zero stock price). "Income" strategies (covered calls, short puts) are routinely mis-sold as "safe" — they are not (see [Disconfirming findings](#disconfirming-findings--common-myths)).

## Contents

- [Organizing principle: defined- vs undefined-risk](#organizing-principle-defined--vs-undefined-risk)
- [Single-leg directional (long call / long put)](#single-leg-directional-long-call--long-put)
- [Stock + option combinations](#stock--option-combinations)
- [Vertical spreads](#vertical-spreads)
- [Volatility strategies (straddle / strangle)](#volatility-strategies-straddle--strangle)
- [Neutral defined-risk premium collection](#neutral-defined-risk-premium-collection)
- [Time-based spreads (calendar / diagonal)](#time-based-spreads-calendar--diagonal)
- [Assignment, pin risk & early assignment around ex-dividend](#assignment-pin-risk--early-assignment-around-ex-dividend)
- [Options approval levels & margin](#options-approval-levels--margin)
- [Bid/ask spreads, liquidity & multi-leg orders](#bidask-spreads-liquidity--multi-leg-orders)
- [Tax pointer](#tax-pointer)
- [Disconfirming findings — common myths](#disconfirming-findings--common-myths)
- [References](#references)

## Organizing principle: defined- vs undefined-risk

The single most decision-relevant lens. **Defined-risk** = the worst-case loss is known and capped at trade entry. **Undefined-risk** = the loss is open-ended or "large down to zero," not bounded by a long option or by cash/stock. The dividing line is almost always whether a short option is *covered* (by a long option — a spread — or by stock/cash) versus left *naked*.[^oic-longcall][^schwab-vertical]

- **Long options (buyers)** are always defined-risk: max loss = premium paid. You **cannot** lose more than the premium.[^oic-longcall][^oic-longput]
- **Naked short call = theoretically UNLIMITED loss** (the stock can rise without bound; you must buy at any price to deliver at the strike).[^oic-nakedcall][^wiki-nakedcall]
- **Naked short put = large but bounded loss** = (strike − premium) × 100, realized if the stock goes to zero. Often *called* "undefined" because it is uncapped relative to the premium and can be ruinous, but it is technically floored at zero.[^oic-csp]

## Single-leg directional (long call / long put)

**Long call** — the core bullish bet.
- **Construction:** buy 1 call. **Outlook:** bullish.
- **Max profit:** theoretically unlimited (rises with the stock above the strike). **Max loss:** premium paid (100% of it if it expires OTM). **Breakeven:** strike + premium. **Risk: defined.**[^oic-longcall]

**Long put** — the core bearish / hedge bet.
- **Construction:** buy 1 put. **Outlook:** bearish.
- **Max profit:** strike − premium (maximized if the stock → 0). **Max loss:** premium paid. **Breakeven:** strike − premium. **Risk: defined.**[^oic-longput]

## Stock + option combinations

**Covered call / buy-write** (long 100 shares + short 1 call).
- **Outlook:** neutral to moderately bullish. **Max profit:** (call strike − stock cost) + premium — **upside is capped** at the strike. **Max loss:** substantial — the stock can fall toward zero; loss ≈ stock cost − premium. **Breakeven:** stock cost − premium. **Risk:** floored by zero (not naked), but carries **full downside stock risk** plus assignment risk (shares called away). ⚠️ **Not** a "safe income" or downside-protection strategy — the premium is a thin cushion, not a hedge.[^oic-coveredcall][^proshares-cc]

**Cash-secured put** (short 1 put + cash set aside to buy 100 shares) — "get paid to wait" / acquisition.
- **Outlook:** neutral to moderately bullish; willing to buy the stock at the strike. **Max profit:** premium received. **Max loss:** strike − premium (large; realized if the stock craters — you are obligated to buy at the strike no matter how low). **Breakeven:** strike − premium. **Risk:** profile is **almost identical to a covered call** — same large downside, capped gain. The cash makes it *secured*, not low-risk.[^oic-csp]

**Protective put** (long shares + long put) — portfolio insurance / "married put."
- **Outlook:** bullish but wants downside insurance. **Max profit:** unlimited (retains upside, lagging an unhedged holder by the premium). **Max loss:** limited to (stock cost − put strike) + premium — the put floors the downside. **Breakeven:** stock cost + premium. **Risk: defined downside.** Cost = premium drag.[^oic-protput]

**Collar** (long shares + long put + short call) — cheap/zero-cost insurance, both ends capped.
- **Construction:** hold stock, buy a protective put (floor), sell a call (ceiling) to fund the put. **Max profit:** capped at the call strike. **Max loss:** capped at the put strike. **Breakeven:** stock cost ± net debit/credit. **Risk: defined on both sides.**[^oic-collar]

## Vertical spreads

Two same-type, same-expiry options at different strikes. **Debit spreads** (you pay; bullish call / bearish put) — time decay works against you. **Credit spreads** (you collect; bullish put / bearish call) — time decay works for you. Max loss **and** max profit are both capped and known at entry; the long leg caps the loss a naked short would otherwise leave open.[^oic-bullcall][^oic-bullput][^schwab-vertical]

| Spread | Direction | Debit/Credit | Construction | Max profit | Max loss | Breakeven |
| --- | --- | --- | --- | --- | --- | --- |
| **Bull call** | Bullish | Debit | Buy lower-strike call, sell higher-strike call | (strike width) − net debit | net debit paid | long-call strike + net debit |
| **Bear put** | Bearish | Debit | Buy higher-strike put, sell lower-strike put | (strike width) − net debit | net debit paid | higher strike − net debit |
| **Bull put** | Bullish | Credit | Sell higher-strike put, buy lower-strike put | net credit received | (strike width) − net credit | short-put strike − net credit |
| **Bear call** | Bearish | Credit | Sell lower-strike call, buy higher-strike call | net credit received | (strike width) − net credit | short-call strike + net credit |

All four: **defined risk, defined reward.**[^oic-bullcall][^oic-bullput]

## Volatility strategies (straddle / strangle)

**Long straddle** (buy call + buy put, **same** strike/expiry) — big-move bet, direction-agnostic.
- **Outlook:** expects a large move either way (e.g., earnings); long volatility. **Max profit:** unlimited on the upside, very large on the downside. **Max loss:** both premiums paid (**defined**). **Breakevens:** strike + total premium (up); strike − total premium (down).[^oic-straddle][^fidelity-straddle]

**Long strangle** (buy OTM call + buy OTM put, different strikes) — cheaper than a straddle, needs a bigger move.
- **Max profit:** unlimited up / large down. **Max loss:** both premiums (**defined**). **Breakevens:** call strike + total premium; put strike − total premium.[^fidelity-strangle]

**Short straddle / short strangle** (sell the same legs) — range-bound premium collection.
- **Max profit:** premiums received (capped). **Max loss: UNLIMITED** (the short naked call leg → infinite upside risk; large downside risk on the put leg). **Risk: undefined.** The dangerous mirror image of the long versions.[^macroption-shortstraddle]

## Neutral defined-risk premium collection

**Iron condor** (short OTM put spread + short OTM call spread) — profit if the stock stays in a range.
- **Construction:** a short strangle with protective long wings (short call + higher long call; short put + lower long put). **Outlook:** range-bound / low volatility. **Max profit:** net credit (both spreads expire worthless). **Max loss:** (wider wing width) − net credit (**defined** — the long wings cap each side). **Breakevens:** short-put strike − credit (down); short-call strike + credit (up). Wider profit zone but lower max profit than the iron butterfly.[^tradestation-condor][^oa-ironfly]

**Iron butterfly** (short ATM put + short ATM call at the same strike, with long OTM wings) — pin a price, defined risk.
- **Max profit:** net credit; achieved only if the stock pins the short strike at expiry. **Max loss:** wing width − net credit (**defined**). **Breakevens:** short strike ± credit. Higher max profit but a narrower, harder-to-hit profit zone than the condor.[^tradestation-condor][^oa-ironfly]

**Long butterfly spread** (buy 1 lower strike, sell 2 middle strike, buy 1 higher strike — all calls OR all puts, same expiry) — bet the stock pins the middle strike.
- **Outlook:** neutral; expects the price near the middle strike at expiry. **Max profit:** wing width − net debit (at the middle strike). **Max loss:** net debit paid. **Breakevens:** lower strike + debit; upper strike − debit. **Risk: defined.**[^fidelity-butterfly][^britannica-butterfly]

## Time-based spreads (calendar / diagonal)

**Calendar (horizontal) spread** — exploit time decay / term structure.
- **Construction:** **same strike, different expirations** — typically sell the near-term, buy the longer-term ("long calendar"). "Horizontal" = same row (strike), different columns (expiry). **Outlook:** neutral near-term; wants the underlying near the strike so the short near-dated leg decays faster than the long-dated leg. **Mechanics:** long calendar = **net positive theta** and long vega; trades the horizontal volatility skew. **Risk:** loss limited to the net debit — a debit, **defined-risk** position.[^fidelity-calendar][^oic-putcalendar]

**Diagonal spread** — calendar + vertical hybrid.
- **Construction:** **different strikes AND different expirations** ("diagonal" = different row and column). Combines a directional lean (the strike difference) with the term-structure/theta play. **Outlook:** mildly directional while harvesting time decay; e.g., a long diagonal call is a capital-efficient stand-in for a covered call ("poor man's covered call"). **Risk:** long diagonal is a **defined-risk** debit position.[^fidelity-calendar-vp]

## Assignment, pin risk & early assignment around ex-dividend

- **Assignment risk (general):** any short option can be assigned, forcing you into (covered call) or out of cash into (cash-secured/naked put) a stock position.[^oic-exercise]
- **Early assignment risk:** American-style equity options can be exercised *any time* before expiry — most likely on a short call just before an **ex-dividend** date, or on deep-ITM shorts.[^ibkr-pin][^schwab-divrisk]
- **PIN RISK:** when a short option sits *right at the strike* at expiration, you cannot know whether it will be assigned, so you cannot hedge precisely. Because equity options trade through Friday but the assignment notice can arrive over the weekend, you may discover Monday that you were assigned — leaving an **unexpected, unhedged stock position over the weekend** exposed to Monday's open. Standard mitigation: **close or roll** any short option pinned near the strike before expiration rather than letting it ride.[^wiki-pin][^ibkr-pin]

> **Ex-dividend early-assignment mechanics (load-bearing).** The long-call holder exercises early to **capture the dividend** when doing so beats holding the option — specifically when the call's **remaining extrinsic (time) value is less than the dividend**. The short call holder is then assigned and must deliver shares and pay the dividend. Risk **peaks in the session immediately before the ex-dividend date**.[^oa-divassign][^schwab-divrisk]
>
> **Worked example.** Stock at $110; short **105 call** worth $5.50 = $5.00 intrinsic + **$0.50 extrinsic**; upcoming **dividend $0.75**. Since the dividend ($0.75) **>** the extrinsic value ($0.50), early exercise/assignment the day before ex-div is **likely** — exercising captures $0.75 while forfeiting only $0.50 of time value.[^oa-divassign] **Quick check:** compare the **same-strike put**'s price to the dividend — if the put (all extrinsic, since it's OTM) is worth *less* than the dividend, the call will likely be exercised (put-call-parity logic). **Avoid** by closing the short call before ex-div, or rolling up/out (at added cost).[^oa-divassign]
>
> **Puts** have their own early-exercise drivers: short puts face early exercise mainly when **deep ITM with little/no extrinsic value**, and **higher interest rates** raise the propensity (exercising a deep-ITM American put frees cash that earns interest to expiration).[^schwab-rates]

## Options approval levels & margin

**The concept of tiered approval is universal; the names and numbers VARY by broker (volatile, as of 2026).** Every broker tiers options permissions, each unlocking progressively riskier strategies — but labels and counts differ materially (e.g., Robinhood ~2 tiers and does not allow naked selling/undefined-risk spreads; Fidelity has 5 levels; tastytrade has 3 word-named tiers; Schwab/E*TRADE/IBKR each label differently). Treat any specific number/name as **broker-and-date specific**.[^scotia-levels][^oa-levels][^stockbrokers-levels]

A **typical ladder** (names differ by broker):

| Level | Typically unlocks |
| --- | --- |
| **1** | Covered calls and cash-secured puts (fully covered by stock/cash) |
| **2** | Long calls and long puts (often long straddles; some brokers add debit spreads) |
| **3** | Spreads / defined-risk multi-leg (debit & credit spreads, iron condors) — "everything except naked" |
| **4** | Naked / uncovered call and put selling — highest risk, highest approval, largest capital/margin |

The application assesses **suitability** — trading **experience**, **annual income**, **net worth / liquid net worth**, and **investment objectives** (FINRA Rule 2360). Brokers commonly attach account-size minimums to higher tiers; those minimums are broker-set, not regulatory.[^finra2360][^oa-levels]

**Margin & options (as of 2026):**

- In a **cash account** you can do fully-covered strategies: long single-leg calls/puts, **covered calls** (≥100 shares per call), and **cash-secured puts** (full strike × 100 set aside). The only short-option strategies generally permitted in a cash account are the covered call and the cash-secured put.[^tasty-margin][^fidelity-margin]
- **Disconfirming check — "do all options strategies need a margin account?"** No. Multi-leg **spreads** and any uncovered short generally require a **margin account**, but covered calls and cash-secured puts do not.[^tasty-margin]
- A defined-risk spread's collateral is capped at the spread width minus credit received (small, known). Selling **naked** options requires *substantially* larger, formula-based margin (a percentage of the underlying notional, often many multiples of a spread's requirement), sits at the **highest approval tier**, and is exposed to margin calls.[^tasty-margin][^fidelity-margin]
- **PDT / day-trading pointer (volatile):** frequent intraday options trading historically triggered the **Pattern Day Trader** rule ($25k minimum in a margin account). FINRA **replaced** PDT with a risk-based **intraday-margin** framework **effective June 4, 2026** (a $2,000 floor remains for *leveraged* trading; transition through Oct 20, 2027). Pre-2026 "$25k to day-trade" guidance is outdated — verify your broker's current policy. Detail → `trading-risk-management` and the hub's `trading-risks-and-protections` reference.[^finra-pdt][^etrade-pdt]

## Bid/ask spreads, liquidity & multi-leg orders

- **Why options spreads are wider than stock spreads:** a stock trades as one instrument, but its options fragment across dozens-to-hundreds of contracts (every strike × expiry × call/put), splintering liquidity; options are also harder to price/hedge, so market makers demand wider spreads. Liquid index/mega-cap options can quote a few cents; thin names can show spreads of 50–100%+ of the option's price.[^schwab-spreads][^nasdaq-liquidity]
- **OI and volume as liquidity signals:** high volume = active two-sided participation today; high open interest = an established, deep secondary market. Both tend to tighten spreads; a common retail rule of thumb is to prefer contracts with volume in the hundreds-to-thousands and open interest in the thousands.[^oa-oivol][^tradingblock-liq]
- **The spread is a fee** paid on **both** entry and exit, so the round-trip cost is roughly double the half-spread; on illiquid contracts a market order can fill 50–100% worse than mid. Because spreads are a *percentage* of small premiums, this friction often dwarfs commissions — prioritize **liquid** options over chasing low commissions.[^oa-slippage]
- **Order placement:** price off the **mid** and "walk" the limit toward the natural price; a **marketable limit** (set at/just through the opposite quote) seeks an immediate fill while capping the worst price. **Market orders on thin options are discouraged.**[^public-orders]
- **Multi-leg orders** are submitted as **one combined order priced at a net debit or net credit**; all legs fill together at your limit (or not at all), so you are not "legged" into a partial position.[^investingdaily-multileg]

## Tax pointer

(POINTER only — full detail → `trading-regulation-compliance-and-taxes`.) Options carry specific US tax treatment: gains are **short-term vs long-term** by holding period; the **wash-sale rule can apply**; **Section 1256 60/40** treatment (60% long-term / 40% short-term, marked-to-market) applies to broad-based **INDEX** options (SPX, NDX, RUT, VIX) — **NOT** equity/ETF options (AAPL, SPY); and **exercise/assignment adjusts the cost basis** of the resulting stock (premium adds to/subtracts from basis). This is a pointer; detail belongs to the sibling tax skill.[^cboe-tax][^greentrader-1256]

## Disconfirming findings — common myths

1. **"Covered calls are safe income / give downside protection" — FALSE / heavily qualified.** Covered calls provide *little* downside protection; if the stock drops steeply, losses can far exceed the premium. The premium is a thin cushion, not a hedge; the real costs are full downside stock exposure plus **opportunity cost** (capped upside in a rally).[^proshares-cc][^chase-cc]
2. **"Selling puts is reckless while covered calls are safe" — FALSE premise.** A cash-secured/short put and a covered call have essentially the **same** risk/reward profile (OIC: the covered-call profile is "almost identical to a true short put"). Put-selling only *appears* riskier when traders use leverage and over-size — a position-sizing problem (→ `trading-risk-management`), not an inherent property.[^oic-csp][^optionstrategist-cc]
3. **"You can lose more than the premium on a long option" — FALSE.** A long call or put buyer's maximum loss is strictly the premium paid, no matter how far the underlying moves.[^oic-longcall][^oic-longput]
4. **"Naked calls and naked puts are equally risky" — FALSE.** A naked call = theoretically *unlimited* loss; a naked put = *bounded* loss of (strike − premium) × 100 (the stock can't go below zero). Both are dangerous, but only the naked call is truly unlimited.[^oic-nakedcall][^wiki-nakedcall]

## References

Anchored on the OIC's authoritative strategy library, corroborated by ≥2 brokers/educators (Schwab, Fidelity, TradeStation, Option Alpha) per strategy; FINRA grounds the approval/margin facts. Strategy mechanics are stable; broker-tier names, margin formulas, and the PDT change are dated *as of 2026*.

[^oic-longcall]: OIC, "Long Call," https://www.optionseducation.org/strategies/all-strategies/long-call — authoritative options education. verified-as-of: 2026-06-16
[^oic-longput]: OIC, "Long Put," https://www.optionseducation.org/strategies/all-strategies/long-put — options education. verified-as-of: 2026-06-16
[^oic-coveredcall]: OIC, "Covered Call (Buy/Write)," https://www.optionseducation.org/strategies/all-strategies/covered-call-buy-write — options education. verified-as-of: 2026-06-16
[^oic-csp]: OIC, "Cash-Secured Put," https://www.optionseducation.org/strategies/all-strategies/cash-secured-put — options education (profile ≈ covered call; bounded downside). verified-as-of: 2026-06-16
[^oic-protput]: OIC, "Protective Put (Married Put)," https://www.optionseducation.org/strategies/all-strategies/protective-put-married-put — options education. verified-as-of: 2026-06-16
[^oic-collar]: OIC, "Collar (Protective Collar)," https://www.optionseducation.org/strategies/all-strategies/collar-protective-collar — options education. verified-as-of: 2026-06-16
[^oic-bullcall]: OIC, "Bull Call Spread (Debit Call Spread)," https://www.optionseducation.org/strategies/all-strategies/bull-call-spread-debit-call-spread — options education. verified-as-of: 2026-06-16
[^oic-bullput]: OIC, "Bull Put Spread (Credit Put Spread)," https://www.optionseducation.org/strategies/all-strategies/bull-put-spread-credit-put-spread — options education. verified-as-of: 2026-06-16
[^oic-nakedcall]: OIC, "Naked Call (Uncovered Call / Short Call)," https://www.optionseducation.org/strategies/all-strategies/naked-call-uncovered-call-short-call — options education (unlimited risk). verified-as-of: 2026-06-16
[^oic-straddle]: OIC, "Long Straddle," https://www.optionseducation.org/strategies/all-strategies/long-straddle — options education. verified-as-of: 2026-06-16
[^oic-putcalendar]: OIC, "Long Put Calendar Spread (Put Horizontal)," https://www.optionseducation.org/strategies/all-strategies/long-put-calendar-spread-put-horizontal — options education. verified-as-of: 2026-06-16
[^oic-exercise]: OIC, Options Exercise & Assignment FAQ, https://www.optionseducation.org/referencelibrary/faq/options-exercise — options education. verified-as-of: 2026-06-16
[^schwab-vertical]: Schwab, "Bullish & Bearish Vertical Options Spreads," https://www.schwab.com/learn/story/bullish-bearish-vertical-options-spreads — brokerage education. verified-as-of: 2026-06-16
[^schwab-divrisk]: Schwab, "Ex-Dividend Dates: Understanding Dividend Risk," https://www.schwab.com/learn/story/ex-dividend-dates-understanding-dividend-risk — brokerage education. verified-as-of: 2026-06-16
[^schwab-rates]: Schwab, "How Interest Rate Movements Affect Options Prices" (put early-exercise drivers), https://www.schwab.com/learn/story/how-interest-rate-movements-affect-options-prices — brokerage education. verified-as-of: 2026-06-16
[^fidelity-straddle]: Fidelity, "Long Straddle" strategy guide, https://www.fidelity.com/learning-center/investment-products/options/options-strategy-guide/long-straddle — brokerage education. verified-as-of: 2026-06-16
[^fidelity-strangle]: Fidelity, "Long Strangle" strategy guide, https://www.fidelity.com/learning-center/investment-products/options/options-strategy-guide/long-strangle — brokerage education. verified-as-of: 2026-06-16
[^fidelity-butterfly]: Fidelity, "Long Butterfly Spread with Calls," https://www.fidelity.com/learning-center/investment-products/options/options-strategy-guide/long-butterfly-spread-calls — brokerage education. verified-as-of: 2026-06-16
[^fidelity-calendar]: Fidelity, "Long Calendar Spread with Calls," https://www.fidelity.com/learning-center/investment-products/options/options-strategy-guide/long-calendar-spread-calls — brokerage education. verified-as-of: 2026-06-16
[^fidelity-calendar-vp]: Fidelity Viewpoints, "Calendar spreads" (calendar & diagonal terminology), https://www.fidelity.com/viewpoints/active-investor/calendar-spreads — brokerage education. verified-as-of: 2026-06-16
[^fidelity-margin]: Fidelity, "Margin requirements" FAQ (options strategies), https://www.fidelity.com/trading/faqs-margin — brokerage. verified-as-of: 2026-06-16
[^macroption-shortstraddle]: Macroption, "Short Straddle Payoff" (unlimited risk), https://www.macroption.com/short-straddle-payoff/ — options reference. verified-as-of: 2026-06-16
[^tradestation-condor]: TradeStation, "Iron Condors & Butterflies," https://www.tradestation.com/learn/options-education-center/iron-condors-butterflies/ — brokerage education. verified-as-of: 2026-06-16
[^oa-ironfly]: Option Alpha, "Iron Butterfly," https://optionalpha.com/strategies/iron-butterfly — options education. verified-as-of: 2026-06-16
[^britannica-butterfly]: Britannica Money, "Butterfly option spread," https://www.britannica.com/money/butterfly-option-spread — education (corroborating). verified-as-of: 2026-06-16
[^wiki-nakedcall]: Wikipedia, "Naked call" (unlimited loss), https://en.wikipedia.org/wiki/Naked_call — encyclopedic (corroborating). verified-as-of: 2026-06-16
[^wiki-pin]: Wikipedia, "Pin risk," https://en.wikipedia.org/wiki/Pin_risk — encyclopedic (pin risk, weekend assignment). verified-as-of: 2026-06-16
[^ibkr-pin]: Interactive Brokers, "Understanding Special Exercises and Pin Risk," https://www.interactivebrokers.com/campus/traders-insight/securities/options/understanding-special-exercises-and-pin-risk/ — brokerage education. verified-as-of: 2026-06-16
[^oa-divassign]: Option Alpha, "Dividend Assignment Risk" (extrinsic-vs-dividend rule + worked example), https://optionalpha.com/learn/dividend-assignment-risk — options education. verified-as-of: 2026-06-16
[^scotia-levels]: Scotiabank/Scotia iTRADE, "What is Level 1/2/3/4 options trading," https://help.scotiabank.com/article/what-is-level-1-2-3-4-options-trading — broker (tier ladder). verified-as-of: 2026-06-16
[^oa-levels]: Option Alpha, "Options Approval Levels," https://optionalpha.com/learn/options-approval-levels — options education (ladder concept; names vary; application factors). verified-as-of: 2026-06-16
[^stockbrokers-levels]: StockBrokers.com, Fidelity vs Robinhood comparison (per-broker tier counts), https://www.stockbrokers.com/compare/fidelityinvestments-vs-robinhood — broker comparison. verified-as-of: 2026-06-16
[^finra2360]: FINRA Rule 2360 (options account approval & supervision; suitability info collected), https://www.finra.org/rules-guidance/rulebooks/finra-rules/2360 — regulator. verified-as-of: 2026-06-16
[^tasty-margin]: tastytrade, "Margin vs Cash Accounts" (cash-secured put / covered call are the only short strategies in a cash account; spreads need margin), https://tastytrade.com/learn/accounts/account-resources/margin-vs-cash-accounts/ — brokerage education. verified-as-of: 2026-06-16
[^finra-pdt]: FINRA, "Intraday Margin Requirements" Investor Insights (PDT replaced by intraday margin, effective 2026-06-04), https://www.finra.org/investors/insights/intraday-margin-requirements — regulator. verified-as-of: 2026-06-16
[^etrade-pdt]: E*TRADE, "Pattern Day Trading rule change," https://us.etrade.com/knowledge/library/margin/pattern-day-trading-rule-change — brokerage (corroborating the 2026 change). verified-as-of: 2026-06-16
[^schwab-spreads]: Schwab, "Large Bid/Ask Options Spreads in Volatile Markets," https://www.schwab.com/learn/story/large-bidask-options-spreads-volatile-markets — brokerage education. verified-as-of: 2026-06-16
[^nasdaq-liquidity]: Nasdaq, "How Options Liquidity Impacts Trading," https://www.nasdaq.com/articles/how-options-liquidity-impacts-trading-2017-10-30 — exchange/education (liquidity fragmentation). verified-as-of: 2026-06-16
[^oa-oivol]: Option Alpha, "Options Volume vs Open Interest," https://optionalpha.com/learn/options-volume-vs-open-interest — options education. verified-as-of: 2026-06-16
[^tradingblock-liq]: TradingBlock, "Options Liquidity," https://www.tradingblock.com/blog/options-liquidity — brokerage education. verified-as-of: 2026-06-16
[^oa-slippage]: Option Alpha, "Wide Bid/Ask Spreads & Slippage Are Costing You" (spread as a hidden round-trip fee), https://optionalpha.com/podcast/wide-bidask-spreads-slippage-are-costing-you-3840-each-year — options education. verified-as-of: 2026-06-16
[^public-orders]: Public.com, "Market vs Limit Orders for Options" (marketable limit; avoid market orders on thin options), https://help.public.com/en/articles/8852337-market-vs-limit-orders-for-options — brokerage education. verified-as-of: 2026-06-16
[^investingdaily-multileg]: Investing Daily, "Multi-Leg Options Orders" (single net debit/credit fill), https://www.investingdaily.com/44469/multi-leg-options-orders/ — education. verified-as-of: 2026-06-16
[^cboe-tax]: Cboe, "Index Options Benefits & Tax Treatment" (Section 1256 60/40 for index options, not equity), https://www.cboe.com/tradable_products/index-options-benefits-tax-treatment/ — exchange. verified-as-of: 2026-06-16
[^greentrader-1256]: Green Trader Tax, "Section 1256 Contracts" (scope; equity options excluded; wash-sale interaction), https://greentradertax.com/trader-tax-center/tax-treatment/section-1256-contracts/ — tax specialist (corroborating; detail deferred). verified-as-of: 2026-06-16
[^proshares-cc]: ProShares, "Covered Call ETFs: The Myth of Downside Protection," https://www.proshares.com/browse-all-insights/insights/covered-call-etfs-the-myth-of-downside-protection — issuer research (covered-call myth). verified-as-of: 2026-06-16
[^chase-cc]: Chase, "What are covered calls" (opportunity cost / capped upside), https://www.chase.com/personal/investments/learning-and-insights/article/what-are-covered-calls — brokerage education. verified-as-of: 2026-06-16
[^optionstrategist-cc]: The Option Strategist, "Covered Calls vs Selling Puts — Same Strategy, Different Name," https://optionstrategist.com/blog/2025/04/covered-calls-vs-selling-puts-same-strategy-different-name — options education (corroborating the equivalence). verified-as-of: 2026-06-16
