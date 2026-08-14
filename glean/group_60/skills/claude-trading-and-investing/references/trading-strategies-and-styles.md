<!-- Provenance: reference under the `trading-and-investing` hub. Mirrored from ~/.claude/skills/trading-and-investing/references/trading-strategies-and-styles.md by scripts/persist-spoke.mjs. -->

# Trading Strategies and Styles

> **Educational information only — NOT financial, investment, or trading advice.** Trading carries risk of loss; the large majority of active retail traders lose money. Past strategy performance does not guarantee future results. All volatile claims (statistics, regulatory rules, capital requirements) are stamped **as of 2026**; verify current details before acting.

## Contents

1. [Overview](#overview)
2. [Trading Styles by Time Horizon](#1-trading-styles-by-time-horizon)
3. [Momentum Strategies](#2-momentum-strategies)
4. [Mean-Reversion Strategies](#3-mean-reversion-strategies)
5. [Trend Following and Breakout Trading](#4-trend-following-and-breakout-trading)
6. [Backtesting a Discretionary Strategy](#5-backtesting-a-discretionary-strategy)
7. [Strategy Development Workflow](#6-strategy-development-workflow)
8. [Anti-Patterns Consolidated](#anti-patterns-consolidated)
9. [Cross-References](#cross-references)
10. [References](#references)

---

## Overview

Trading strategy is the ruleset that defines when to enter a position, when to exit, and how large to be. This reference covers the *discretionary and semi-systematic* side of that ruleset — the trader's playbook organized by time horizon and by the two dominant market dynamics any strategy must handle:

- **Trending markets** → momentum and trend-following strategies
- **Range-bound markets** → mean-reversion and pairs-based strategies

For the fully systematic/algorithmic side (execution algorithms, cointegration statistics, market making, factor models), see `references/algorithmic-and-quant-trading.md`.

**Failure base rate (verified-as-of: 2026-06-21)**: A 14-year study of 360,000+ individual traders found only 1% consistently beat the market, only 13% were still trading after 3 years, and fewer than 7% remained active by year 5 [^4][^5]. FINRA data shows 72% of day traders end the year with net losses [^4]. These figures set the context: an edge that survives statistical validation is genuinely rare.

---

## 1. Trading Styles by Time Horizon

### Scalping (seconds to minutes)

Scalpers execute 50–200+ trades per day, targeting gains of 0.05–0.2% per trade. The math is unforgiving: transaction costs must be minimized because the profit margin per trade is tiny. Scalpers typically require **80%+ win rates** to remain profitable given unfavorable risk-reward ratios (often 1:0.5 — risking $1 to make $0.50). The high win rate compensates for a small reward-to-risk ratio [^1][^2].

Key requirements: instruments with narrow bid-ask spreads (tight forex majors, liquid futures, zero-commission liquid equities with actual execution quality), direct access execution (Level 2 quotes, direct routing), and very fast execution infrastructure.

**Common failure mode**: Beginners fail scalping "not because of strategy, but because they ignore transaction costs" — a 0.1% spread in both directions on a 0.15% target profit leaves 0.05% before commissions and slippage [^1].

### Day Trading (intraday — all positions closed before market close)

Day trading avoids overnight gap risk but requires the most active attention and operationally the highest regulatory compliance burden [^3][^6].

**PDT rule (US, as of 2026)**: FINRA's Pattern Day Trader rule was revised. The original fixed-count threshold (4+ day trades within 5 business days required $25,000 minimum account balance in a margin account) was replaced by a risk-based intraday margin framework focused on real exposure rather than trade count [^6]. Many practitioners still recommend $50,000–$100,000 to trade actively while managing risk.

**Day trading failure statistics (verified-as-of: 2026-06-21)** [^4][^5]:
- Only 1% of traders consistently profitable across a 14-year study of 360,000+ traders
- ~13% still trading after 3 years; ~7% after 5 years
- 72% end the year with net losses per FINRA data
- Profitable day traders represent approximately 1.6% of all participants in an average year

### Swing Trading (days to weeks)

Swing traders capture trend "swings" — entering after a pullback in a trending market or at range extremes. The time horizon allows for meaningful risk-reward setups [^7][^8]:
- Standard target: 1:2 minimum risk-reward; strong setups targeting 1:3
- Common win rates: 45–65% depending on setup type
- Stop-losses placed below recent support; profit targets at next resistance or measured-move projections
- Practical size rule: risk 1–2% of account per trade; size the position based on the stop distance

### Position Trading (weeks to months)

Position traders hold through the full lifespan of a macro trend. Unlike buy-and-hold, position traders define a specific catalyst or technical trigger and manage risk with defined stops [^9][^10]:
- Time horizon: 1 month to 2 years
- Require capital sufficient to withstand larger pullbacks than shorter-term traders
- Best suited for major sector rotations, macro trend shifts, or post-breakout thesis trades

### Style Comparison

| Style | Trade Duration | Practical Capital Floor (US) | Win Rate Range | Target Risk-Reward |
|-------|---------------|------------------------------|----------------|-------------------|
| Scalping | Seconds–minutes | Low per trade; execution speed matters most | 80%+ required | 1:0.5 typical |
| Day Trading | Minutes–hours | $25K–$100K (margin account) | 40–55% | 1:1.5–1:2 |
| Swing Trading | Days–weeks | $5K–$25K | 45–65% | 1:2–1:3 |
| Position Trading | Weeks–months | $10K+ | 35–55% | 1:3–1:5+ |

Sources: [^11][^12][^13]

---

## 2. Momentum Strategies

### The Academic Foundation: Jegadeesh-Titman (1993)

Jegadeesh and Titman's landmark paper *"Returns to Buying Winners and Selling Losers"* established cross-sectional momentum as a systematic anomaly [^14][^15][^16]. Strategy construction:
- **Formation period**: prior 3–12 months of returns (12-month is most common)
- **Skip**: exclude the most recent month (the 12-1 strategy) to avoid the 1-month reversal from bid-ask bounce
- **Universe**: rank all stocks; buy top decile (past winners), short bottom decile (past losers)
- **Holding period**: 3–12 months

Results across all 16 tested combinations of formation/holding periods: approximately **1% per month** in excess return from the Winner Minus Loser (WML) portfolio.

A 2022 review covering 30 years of post-publication evidence confirmed: momentum remains "perhaps the most pervasive contradiction of the efficient market hypothesis," robust across asset classes, geographies, and time periods [^15].

### Cross-Sectional vs. Time-Series Momentum

**Cross-sectional (relative) momentum**: Ranks securities against each other. Long the top cohort, short the bottom cohort. Does not require any particular absolute market direction — a stock is "momentum" if it outperformed its peers, even in a down market.

**Time-series / absolute momentum** (Moskowitz, Ooi & Pedersen, 2012): Compares each security's return to its own past. If a security's 12-month return is positive, go long; if negative, go to cash or short. Tested across 58 futures markets with documented Sharpe ratios in excess of 1.0. Unlike cross-sectional momentum, it exits the asset class entirely during bad periods, providing crash protection [^17][^18].

**Dual Momentum** (Antonacci): Combines both — first use absolute momentum to decide whether to be in risk assets at all (if equities have negative 12-month absolute momentum, hold T-bills instead), then use relative momentum to pick between risk assets. Backtested to 1974, showing improved Sharpe and dramatically lower drawdowns versus buy-and-hold [^17][^18].

Evidence for momentum is exceptional: documented back to 1801 in equities (Geczy & Samonov 2015) and across commodities since 1223 (Greyserman & Kaminski 2014).

### Momentum Crashes: The Critical Risk

Daniel and Moskowitz (2016) documented that momentum strategies experience infrequent but severe **crash events** — periods of persistent large negative returns [^19][^20].

**Mechanism**: Momentum strategies implicitly short out-of-the-money options on market recovery. After a severe market decline, past losers (the short leg) are like deeply distressed options. When markets rebound sharply, these losers can surge 20–50%+ while past winners stagnate, generating massive losses on the short side.

**Crash conditions**:
- Following large market declines (panic states)
- When market volatility is elevated
- During sharp, V-shaped market reversals
- After bear markets end

**Historical examples**: 2001–2002 recovery; Spring–Summer 2009 (momentum strategies lost 60–80% in some implementations).

**Mitigation**: A dynamic strategy that forecasts momentum's conditional mean and variance can reduce exposure before predicted crashes — this approach nearly doubles the Sharpe ratio relative to static momentum [^19].

### Anti-Patterns in Momentum Trading

- Running momentum in isolation without crash hedging or volatility scaling
- Using less than 6-month lookback (too noisy) or more than 12-month (too slow to signal)
- Ignoring the 1-month reversal effect — always skip the most recent month in the formation period
- Trading momentum in illiquid stocks where execution costs eliminate the premium
- Crowding: momentum strategies use similar methodology; over-owned names amplify crashes when funds unwind simultaneously [^20][^21]

---

## 3. Mean-Reversion Strategies

### Core Concept

Mean reversion assumes that prices exhibiting extreme deviations from their historical average will return to that average. The strategy family includes range trading, Bollinger Band strategies, oversold/overbought systems, and spread/pairs trading [^22][^23].

**The regime dependency problem**: Mean reversion works only in range-bound markets. In a sustained trend, "oversold" can become more oversold for weeks. Misidentifying a trending market as a ranging one is the primary failure mode. ADX > 25 indicates a trending market (where mean-reversion logic fails); ADX < 20 suggests a ranging market (where it thrives) [^26].

### Bollinger Band Mean Reversion

Bollinger Bands (20-period moving average ± 2 standard deviations) create a dynamic envelope [^23]:
- Price touching the **lower band**: signals oversold conditions; potential long entry
- Price touching the **upper band**: signals overbought conditions; potential short entry
- **Band squeeze** (bands narrowing): indicates low volatility; often precedes a breakout rather than continued mean reversion

**Signal quality rule**: Mean-reversion signals on Bollinger Band touches are far less reliable when price is in a sustained trend. The touch is a trend-continuation signal, not a reversal signal.

### Pairs Mean-Reversion and Cointegration

For pairs/spread trading (the quantitative version), see `references/algorithmic-and-quant-trading.md` §3. The core distinction:
- **Correlation** ≠ **Cointegration**: Two securities can move together (correlation) without maintaining a stable long-run relationship. Pairs trading requires cointegration — a stationary linear combination of price series.
- **Z-score thresholds**: Enter at ±2.0 standard deviations; exit at 0 or ±0.5; stop-loss at ±3.0 (divergence beyond 3σ is more likely regime change than opportunity) [^24][^25]

### Half-Life of Mean Reversion

The Ornstein-Uhlenbeck process models mean-reverting series and yields a **half-life** parameter: the expected time for the spread to revert halfway to its mean [^24][^25]:
- Set maximum hold period at 2–3× the half-life
- Example: half-life of 5 days → exit position at 10–15 days maximum regardless of status
- Half-life over 30 days → the series may not be mean-reverting in any tradeable sense

### When Mean Reversion Fails

1. **Trending markets**: RSI can remain overbought for weeks in a strong uptrend. The "overbought" reading is a trend-continuation indicator, not a reversal signal.

2. **Regime change**: When a market transitions from range-bound to trending, the mean itself shifts. Shorting "overbought" in a secular bull market produces compounding losses.

3. **Asymmetric P&L profile**: Mean-reversion strategies often show 60–80% win rates — attractive-looking — but hide rare catastrophic losing trades when the mean shifts permanently. This profile is characterized as "picking up nickels in front of a steamroller" [^26].

4. **Cointegration breakdown**: Historically correlated pairs can decouple permanently due to mergers, regulatory changes, or business-model divergence. The 2008 financial crisis broke many historically stable pairs [^26][^27].

---

## 4. Trend Following and Breakout Trading

### The Evidence Base: A Century of Trend-Following Returns

AQR's landmark paper *"A Century of Evidence on Trend-Following Investing"* (Hurst, Ooi & Pedersen, 2017) studied 67 markets across 4 asset classes from January 1880 to December 2016 [^28][^29]:
- **Every decade since 1880** produced positive average returns from time-series momentum
- Performed well in **8 of 10 largest 60/40 portfolio drawdowns**
- Works across recessions/booms, war/peacetime, high/low inflation, high/low interest rates
- Low correlation to traditional asset classes — a genuine portfolio diversifier

### The Turtle Trading Rules (1983–1988)

Richard Dennis's experiment proved systematic trend-following could be taught from first principles. The rules (now public) [^30][^31][^32]:

**System 1** (shorter-term): Enter on a **20-day Donchian channel breakout** (new 20-day high or low). Skip if the previous System 1 signal in the same direction was a winner.

**System 2** (longer-term): Enter on a **55-day Donchian channel breakout**.

**Position sizing**: Size in **ATR units** — each unit risks approximately 1% of account per ATR of price movement. Add up to 4 units, adding at ½ ATR intervals as the trade moves favorably.

**Stops**: **2 ATR from entry** for all positions.

**Exits**: System 1 exits on a 10-day *opposite* breakout; System 2 exits on a 20-day opposite breakout.

Performance characteristics: **30–35% win rates** — most trades are small losses from false breakouts. Large winners cover all losses. The original Turtles reportedly returned $175M over 5 years at 80%+ annual compounding. Practically: missing even a handful of big trends per year can ruin annual performance. This is the "all or nothing" nature of trend-following.

### ATR-Based Entry and Breakout Confirmation

ATR (Average True Range) provides a volatility-normalized unit for execution [^33][^34]:
- **Stop placement**: 1.5–2 ATR below entry (long) / above entry (short)
- **Breakout confirmation**: require the breakout candle body to exceed 70% of its range (filters wicks-only breakouts)
- **Position sizing**: target dollar risk = 1 ATR × shares = fixed percentage of account (typically 1%)

**Reducing false breakouts**:
- The **two-candle rule**: require two consecutive closes beyond the breakout level — eliminates most single-candle fakeouts
- Higher timeframes (4H, daily) produce more reliable breakouts with lower false-signal rates than intraday
- Volume confirmation: a genuine breakout typically shows above-average volume at the breakout candle

### Trend Following in Futures vs. Equities

Original Turtle rules targeted diversified futures: currencies, US Treasuries, metals, crude oil, grains — 20–24 liquid markets. Futures advantages for trend-following:
- Inherent leverage (capital efficiency)
- No uptick rule for shorts
- No borrowing costs
- Lower transaction costs at institutional scale
- Macro trend exposure (rates, currencies, commodities) uncorrelated to equities

**In equities**: Pure trend-following on individual stocks is harder due to survivorship bias in backtests and the equity drift bias (markets trend upward overall, so short signals are systematically disadvantaged). **Cross-sectional momentum** (buying top decile by trailing return, shorting bottom decile) is the more robust equity adaptation.

### When Trend Following Fails

- **Choppy, range-bound markets**: Trend systems get whipsawed — entering on breakouts that immediately reverse. Multiple whipsaws in succession create extended drawdowns.
- **V-shaped recoveries**: The system misses the initial rally leg by design (needs the trend to establish), then enters late and gets stopped out on consolidation.
- **Post-GFC period (2009–2013)**: The SG Trend Index returned −1.8% over 5 years as central bank rate suppression removed macro trends in rates and currencies.
- **"Year of the whipsaw" (2023)**: Multiple asset classes produced repeated false breakouts; trend systems underperformed after three strong years (2020–2022) [^35][^36].
- **Structural reality**: Trend-following underperformed the S&P 500 in roughly half of all years historically. Its value is crash protection and low correlation, not consistent equity-beating returns.

---

## 5. Backtesting a Discretionary Strategy

*This section shares foundations with `references/algorithmic-and-quant-trading.md` §1 and §5. Cross-reference there for framework selection and the quant research pipeline. This section focuses on the discretionary → systematic translation problem.*

### The Translation Problem

The first and hardest step: converting discretionary judgment into explicit, replicable rules. Every vague decision must become unambiguous [^37][^38]:

- "It looks strong" → "Price is above the 50-day EMA AND 20-day ATR is above its 50-day average AND yesterday's range was in the top 25% of the last 20 days"
- "I'll exit if it feels wrong" → "Exit if price closes below the 8-day EMA for two consecutive days OR hits a hard stop at 1.5 ATR from entry"

**Rule completeness test**: If another person cannot replicate every trade identically from the written rules alone, the rules are not complete — you have a feeling, not a rule.

### Data and Sample Size Requirements

- **Minimum data**: 2 years; ideally 5–10 years spanning multiple market regimes (bull + bear + range-bound)
- **Minimum trade instances**: 100 trade occurrences (50 is the statistical floor; 200 is more reliable)
- **Regime coverage**: Include trending AND range-bound AND high-volatility periods in the backtest window
- **Survivorship bias**: For equity backtests, the dataset must include delisted stocks; testing only on companies that survived to the present inflates returns because bankruptcies and delistings (often the worst performers) are excluded [^39]

### In-Sample / Out-of-Sample Split

Reserve **30% of data** strictly for out-of-sample testing — do not look at it during development, do not modify rules based on its results [^39][^40]:

1. 70% in-sample: develop rules, optimize parameters
2. 30% out-of-sample: single validation test, never re-optimized against

Significant performance degradation in out-of-sample data (>30% worse) is a strong indicator of overfitting to the in-sample period.

### Walk-Forward Analysis (WFA)

Gold standard beyond a simple OOS split [^40][^41]:
1. Optimize on initial in-sample window → test on subsequent out-of-sample window
2. Slide both windows forward in time
3. Repeat across the full data history

This simulates how the strategy would have performed if developed and re-calibrated in real time. It prevents data snooping while testing robustness across multiple market regimes.

### Monte Carlo Simulation of Trade Sequences

Randomizes the order of individual trades across 1,000+ iterations to determine [^40]:
- Whether the observed equity curve depends on lucky trade sequencing
- The distribution of possible drawdowns and final returns
- The probability that observed performance resulted from skill vs. luck

A strategy that remains profitable across 1,000 randomized trade sequences has genuine edge independent of sequencing luck.

### Realistic Cost Modeling

Common backtesting failures [^39]:
- Ignoring bid-ask spread (particularly harmful for scalpers and small-caps)
- Using "fill at close" for end-of-day signals when actual execution is next open
- Ignoring market impact on positions >0.5% of average daily volume
- Treating zero-commission as zero-cost (PFOF affects execution quality even at $0 commission)

**Practical approach**: add 0.05–0.15% per trade to account for slippage; stress-test at 2× and 3× assumed slippage to measure fragility. Live results typically degrade 10–20% from backtested results due to execution differences.

### Key Backtesting Metrics

| Metric | Good Range | Warning Sign |
|--------|-----------|--------------|
| Profit Factor (gross profit / gross loss) | >1.5 | <1.0 = losing strategy |
| Maximum Drawdown | <15% | >25% = recovery difficulty |
| Sharpe Ratio | >1.0 | <0.5 = poor risk-adjusted |
| Sample Size | 200+ trades | <100 = statistically unreliable |
| Win Rate | Context-dependent | Evaluate with R-ratio, not in isolation |

### Expectancy Formula

The single most useful single-number metric for evaluating a strategy:

```
Expectancy = (Win Rate × Average Win) − (Loss Rate × Average Loss)
```

Expressed in R-multiples (where 1R = amount risked per trade):
- Above 0.2R = decent
- Above 0.5R = good
- Above 1.0R = excellent

**Why this matters**: A 40% win rate with 3:1 reward-risk has expectancy = 0.40×3 − 0.60×1 = **0.6R** (good). A 60% win rate with 1:1 reward-risk has expectancy = 0.60×1 − 0.40×1 = **0.2R** (decent). Win rate in isolation is meaningless — it must be paired with risk-reward [^42][^43].

---

## 6. Strategy Development Workflow

### Hypothesis-Driven Development

The antidote to data-mining and overfitting is starting with a *structural reason why an edge should exist* before looking at data [^44][^45]:

- **Risk premium hypothesis**: The strategy compensates for bearing a risk that other participants want to shed (e.g., trend following during crises compensates for drawdown risk that most investors can't stomach)
- **Behavioral mispricing hypothesis**: The strategy exploits systematic investor errors (momentum exploits underreaction; mean-reversion exploits overreaction)
- **Structural flow hypothesis**: The strategy takes the other side of systematic flows from non-profit-seeking participants (hedgers, index rebalancers, end-of-month window dressers)

Without a structural rationale, a backtest showing an edge is likely a noise pattern that will fail in live trading.

### Signal Quality Testing Before Full Backtest

Before building a full strategy, test the signal's information content in isolation [^44]:
- **Information Coefficient (IC)**: Correlation between the signal's prediction and subsequent returns. An IC of 0.05 (5%) is meaningful if consistent. (See `references/algorithmic-and-quant-trading.md` §5.2 for IC mechanics.)
- **IC over time**: Check for decay — signals with strong 1990s IC that fade post-2010 may be arbitraged away
- **IC in different regimes**: Does the signal work in both trending and range-bound environments, or only one?

### Alpha Decay

A fundamental challenge in strategy development: profitable signals attract capital, which arbitrages away the edge [^46][^47]:
- Signals are now "arbitraged away in hours or minutes" in highly competitive spaces
- Warning signs: backtest performs well in early periods but degrades sharply in recent years; live trading underperforms backtest by >30%
- Mitigation: faster execution infrastructure, proprietary data, or combining multiple signals with genuine independence

For empirical alpha decay rates (US: 5.6%/year; EU: 9.9%/year), see `references/algorithmic-and-quant-trading.md` §5.5.

### From Discretionary to Systematic: A Staged Approach

The shift is "less about changing strategies and more about changing how decisions are made" [^37][^38][^48]:

**Stage 1 — Rule documentation**: Write every entry/exit/filter/sizing rule in explicit if-then logic. If you cannot write it, you have a feeling, not a rule.

**Stage 2 — Manual backtest**: Replay charts bar-by-bar; record every trade that would have been triggered. Log: setup type, entry, stop, target, exit, P&L, R-multiple, market context.

**Stage 3 — Statistical validation**: After 50+ manual trades, calculate expectancy. After 100+ trades, assess statistical significance.

**Stage 4 — Code translation**: Convert explicit rules to code for automated backtesting on larger data samples. Stress-test parameters: if replacing a 20-day EMA with 15-day or 25-day collapses results, the edge was parameter-specific noise.

**Stage 5 — Paper trading → live progression**: 2–4 weeks paper trading, then live at 25–50% of intended size, scaling to full size only after live results match backtest metrics within 10–15%.

### Portfolio-Level Strategy Correlation

A portfolio of strategies should be evaluated on inter-strategy correlation, not individual Sharpe ratios [^49][^50]:
- Median correlation among strategies within the same quant firm: approximately 0.11
- Quant strategies vs. fundamental managers: correlation as low as 0.02
- A portfolio of trend-following + mean-reversion + momentum provides genuine diversification because they outperform in different regimes

**Regime dependency pattern**:
- Trend-following: best during macro moves and high-volatility trending markets
- Mean-reversion: best during range-bound, low-volatility conditions
- Cross-sectional momentum: best during trending risk-on environments

**Warning**: The assumption that uncorrelated strategies remain uncorrelated breaks under severe market stress (March 2020, August 2024 flash events) — strategies correlated toward losses temporarily when institutional deleveraging dominated [^49][^50].

---

## Anti-Patterns Consolidated

| Anti-Pattern | Domain | Consequence |
|---|---|---|
| Trading mean-reversion signals in a trending market | Mean reversion | Systematic losses; "overbought" is a trend-continuation signal in a trend |
| "Feels strong" rules not translated to explicit if-then logic | All discretionary | Unreplicable backtest; no measurable edge |
| Testing on survivorship-biased universe (current members only) | Backtesting | Inflated backtest returns; live underperformance |
| Look-ahead bias in signal generation | Backtesting | Real-world performance fails to replicate |
| Evaluating win rate without risk-reward ratio | All styles | High win rate can be a losing strategy |
| Running momentum without crash protection | Momentum | 50–80% drawdowns in V-recovery environments |
| Ignoring transaction costs at the strategy's actual trade frequency | All styles | Edge destroyed by slippage + commissions |
| Applying trend-following to choppy range-bound markets | Trend following | Whipsaws destroy capital in repeated small losses |
| Confusing correlation with cointegration in pairs trading | Pairs/mean reversion | Spread diverges indefinitely; large losses on "convergence" entries |
| Modifying rules after observing OOS results | Backtesting | OOS test contaminated; data snooping re-enters |
| Portfolio-building without measuring inter-strategy correlation | Strategy portfolio | "Diversified" strategies correlated at worst moments |
| Backtesting on <100 trades | All styles | Results statistically unreliable; any outcome is possible by chance |
| Assuming alpha is permanent once found | Strategy lifecycle | Edge decays as competitors discover and arbitrage the same pattern |

---

## Cross-References

- **Algorithmic and quant trading** (execution algorithms TWAP/VWAP/IS/POV, cointegration math, market making, quant research pipeline, factor models) → `references/algorithmic-and-quant-trading.md` (sister spoke — shares backtesting methodology foundation)
- **Technical analysis** (indicator mechanics — RSI, MACD, moving averages, chart patterns at depth) → `references/technical-analysis.md` (planned spoke)
- **Order lifecycle and execution** (how a retail order reaches the market; PFOF, slippage, execution quality) → `references/order-lifecycle-and-execution.md`
- **Investing vs. trading** (evidence on active trader underperformance, time horizon and goal distinctions) → `references/investing-vs-trading.md`
- **Statistical methodology** (IC, walk-forward analysis, Deflated Sharpe, FDR correction mechanics) → `da-analytical-methods` hub

---

## References

[^1]: [Pepperstone: Scalping Trading Guide](https://pepperstone.com/en/learn-to-trade/trading-guides/scalping-trading/) — scalping mechanics, win rate requirements | educational | verified: 2026-06-21

[^2]: [TMGM: Scalping Trading Strategy](https://www.tmgm.com/en/academy/trading-academy/scalping-trading-strategy) — scalping setup and cost considerations | practitioner blog | verified: 2026-06-21

[^3]: [Admiral Markets: Scalping vs Day Trading vs Swing Trading](https://admiralmarkets.com/education/articles/forex-strategy/scalping-vs-day-trading-vs-swing-trading) — style comparison | educational | verified: 2026-06-21

[^4]: [Vetted Prop Firms: What Percentage of Day Traders Lose Money](https://vettedpropfirms.com/what-percentage-of-day-traders-lose-money/) — 14-year study, 360K+ traders | practitioner research | verified: 2026-06-21

[^5]: [Bookmap: Why Most Day Traders Lose Money and What Winners Do Differently](https://bookmap.com/blog/why-most-day-traders-lose-money-and-what-the-winners-do-differently) — failure statistics, patterns | practitioner blog | verified: 2026-06-21

[^6]: [NerdWallet: Pattern Day Trader Rule Change 2026](https://www.nerdwallet.com/investing/news/pattern-day-trading-rule-change) | [Britannica Money: Pattern Day Trader Rule](https://www.britannica.com/money/pattern-day-trader-rule) — PDT rule mechanics | financial education | verified: 2026-06-21

[^7]: [Elearnmarkets: Essential Elements of Swing Trading](https://www.elearnmarkets.com/school/units/swing-trading/essential-elements-of-swing-trading) — entries, exits, risk-reward | educational | verified: 2026-06-21

[^8]: [Trade That Swing: The 1% Risk Rule](https://tradethatswing.com/the-1-risk-rule-for-day-trading-and-swing-trading/) — position sizing for swing traders | practitioner blog | verified: 2026-06-21

[^9]: [HeyGoTrade: What Is Position Trading Strategy](https://www.heygotrade.com/en/blog/what-is-position-trading-strategy/) — position trading overview | educational | verified: 2026-06-21

[^10]: [MarketClutch: Buy and Hold vs Position Trading](https://marketclutch.com/buy-and-hold-vs-position-trading/) — distinction between investing and position trading | educational | verified: 2026-06-21

[^11]: [Equiti: Compare Trading Styles — Day, Swing, and Position](https://www.equiti.com/sc-en/education/trading-strategies/compare-trading-styles-day-swing-and-position-trading/) — style comparison table | educational | verified: 2026-06-21

[^12]: [EdgeWonk: Choosing the Right Trading Style](https://edgewonk.com/blog/choosing-the-right-trading-style-day-trading-swing-trading-scalping-or-investing) — style selection framework | educational | verified: 2026-06-21

[^13]: [Admiral Markets: Style Comparison](https://admiralmarkets.com/education/articles/forex-strategy/scalping-vs-day-trading-vs-swing-trading) — capital and win rate survey | educational | verified: 2026-06-21

[^14]: [Jegadeesh-Titman 1993 / Wiley](https://onlinelibrary.wiley.com/doi/abs/10.1111/j.1540-6261.1993.tb04702.x) — original momentum paper | academic paper (seminal) | verified: 2026-06-21

[^15]: [Springer Nature: Returns to Buying Winners — 30 Years Later](https://link.springer.com/article/10.1007/s11408-022-00417-8) — 30-year post-publication review | academic paper | verified: 2026-06-21

[^16]: [SSRN: Jegadeesh-Titman Momentum](https://papers.ssrn.com/sol3/papers.cfm?abstract_id=227349) — working paper version | academic paper | verified: 2026-06-21

[^17]: [Robot Wealth: Dual Momentum Review](https://robotwealth.com/dual-momentum-review/) — Antonacci dual momentum, absolute + relative | practitioner blog | verified: 2026-06-21

[^18]: [Optimal Momentum (Antonacci)](https://www.optimalmomentum.com/) — dual momentum methodology | practitioner reference | verified: 2026-06-21

[^19]: [NBER w20439: Daniel-Moskowitz Momentum Crashes](https://www.nber.org/papers/w20439) | [ScienceDirect](https://www.sciencedirect.com/science/article/pii/S0304405X16301490) — momentum crash anatomy, dynamic hedge | academic paper | verified: 2026-06-21

[^20]: [AQR: Fact, Fiction, and Momentum Investing](https://www.aqr.com/-/media/AQR/Documents/Journal-Articles/JPM-Fact-Fiction-and-Momentum-Investing.pdf) — practical realities, crowding | practitioner paper | verified: 2026-06-21

[^21]: [Tandfonline: Crowding in Alternative Risk Premia](https://www.tandfonline.com/doi/full/10.1080/0015198X.2019.1600955) — crowding empirics and divergence premia | academic paper | verified: 2026-06-21

[^22]: [QuantInsti: Mean Reversion Strategies — Introduction and Building Blocks](https://blog.quantinsti.com/mean-reversion-strategies-introduction-building-blocks/) — mean reversion overview | educational | verified: 2026-06-21

[^23]: [CrossTrade: Bollinger Band Mean Reversion](https://crosstrade.io/learn/trading-strategies/bollinger-mean-reversion) — Bollinger Band mechanics | educational | verified: 2026-06-21

[^24]: [QuantVero: Mean Reversion Strategy Build and Test](https://www.quantvero.com/algo-trading/mean-reversion-trading-strategy/) — half-life, OU process | practitioner blog | verified: 2026-06-21

[^25]: [IBKR Quant News: Mean Reversion Strategies — Part I](https://www.interactivebrokers.com/campus/ibkr-quant-news/mean-reversion-strategies-introduction-trading-strategies-and-more-part-i/) — pairs mechanics, half-life | practitioner/educational | verified: 2026-06-21

[^26]: [HeyGoTrade: Mean Reversion in Trading](https://www.heygotrade.com/en/blog/mean-reversion-in-trading/) | [Quantum Algo: Mean Reversion Complete Guide](https://www.quantum-algo.com/blog/guides/mean-reversion-trading-complete-guide/) — regime failure modes | educational | verified: 2026-06-21

[^27]: [ForexTester: Mean Reversion Trading](https://forextester.com/blog/mean-reversion-trading/) — cointegration breakdown examples | practitioner blog | verified: 2026-06-21

[^28]: [AQR: A Century of Evidence on Trend-Following Investing](https://www.aqr.com/Insights/Research/Journal-Article/A-Century-of-Evidence-on-Trend-Following-Investing) | [SSRN](https://papers.ssrn.com/sol3/papers.cfm?abstract_id=2993026) — Hurst/Ooi/Pedersen 2017; 136 years of data | academic paper | verified: 2026-06-21

[^29]: [Alpha Architect: Time-Series Momentum — Historical Evidence](https://alphaarchitect.com/time-series-momentum-aka-trend-following-the-historical-evidence/) — AQR findings summary | practitioner/research | verified: 2026-06-21

[^30]: [Zaye Capital Markets: Turtle Trading Strategy](https://zayecapitalmarkets.com/turtle-trading-strategy/) — original Turtle rules walkthrough | educational | verified: 2026-06-21

[^31]: [Alchemy Markets: Turtle Trading Guide](https://alchemymarkets.com/education/strategies/turtle-trading-guide/) — Turtle system 1 and 2 mechanics | educational | verified: 2026-06-21

[^32]: [ResearchGate: Research on Quantitative Trading Strategies Based on the Turtle Trading Rule](https://www.researchgate.net/publication/370698720_Research_on_Quantitative_Trading_Strategies_Based_on_the_Turtle_Trading_Rule) — academic analysis of Turtle rules | academic paper | verified: 2026-06-21

[^33]: [LuxAlgo: False Breakout Strategies for Traders](https://www.luxalgo.com/blog/5-false-breakout-strategies-for-traders/) — breakout confirmation, ATR | educational | verified: 2026-06-21

[^34]: [ForexTester: Breakout Trading Strategy](https://forextester.com/blog/breakout-trading-strategy/) — breakout entry mechanics | educational | verified: 2026-06-21

[^35]: [Man Group: Is This Time Different — Trend Following](https://www.man.com/insights/is-this-time-different) — trend following drawdowns and recovery | asset manager research | verified: 2026-06-21

[^36]: [TiltFolio: When Trend Following Underperforms](https://www.tiltfolio.com/insights/when-trend-following-underperforms/) | [Top Traders Unplugged: December 2023 Performance Report](https://www.toptradersunplugged.com/trend-following-performance-report-december-2023/) — whipsaws 2023, post-GFC | practitioner research | verified: 2026-06-21

[^37]: [Monster Trading Systems: Systematic Trading Guide](https://www.monstertradingsystems.com/systematic-trading-guide/) — discretionary-to-systematic translation | practitioner blog | verified: 2026-06-21

[^38]: [Medium — Ayush Dhingra: How to Move from Discretionary to Systematic Trading](https://medium.com/@ayush.dh/how-to-move-from-discretionary-trading-to-systematic-quantitative-trading-64d3a189f6ea) — staged transition approach | practitioner blog | verified: 2026-06-21

[^39]: [TradeZella: Backtesting Trading Strategies — Complete Guide](https://www.tradezella.com/blog/backtesting-trading-strategies) | [GoatFundedTrader: 9 Steps to Backtest](https://www.goatfundedtrader.com/blog/backtesting-trading-strategies) — data requirements, OOS split, slippage | practitioner blog | verified: 2026-06-21

[^40]: [GoatFundedTrader: Best Practices for Backtesting](https://www.goatfundedtrader.com/blog/best-practices-for-backtesting-trading-strategies) — Monte Carlo, walk-forward | practitioner blog | verified: 2026-06-21

[^41]: [Surmount: Walk-Forward Analysis vs Backtesting](https://surmount.ai/blogs/walk-forward-analysis-vs-backtesting-pros-cons-best-practices) | [ArXiv: Walk-Forward Validation Framework](https://arxiv.org/html/2512.12924v1) — WFA methodology | practitioner + academic | verified: 2026-06-21

[^42]: [JournalPlus: Win Rate vs Risk-Reward](https://journalplus.co/learn/guides/win-rate-vs-risk-reward/) — expectancy mechanics | educational | verified: 2026-06-21

[^43]: [Pro Trader Dashboard: Expectancy Formula](https://protraderdashboard.com/blog/expectancy-formula-trading/) — expectancy calculation and R-multiples | educational | verified: 2026-06-21

[^44]: [SetupAlpha / Medium: Scientific Workflow for Generating Alpha in Quant Trading in 2026](https://medium.com/@setupalpha.capital/my-scientific-workflow-for-generating-alpha-in-quantitative-trading-in-2026-5238b26d4d95) — hypothesis-driven development, IC | practitioner blog | verified: 2026-06-21

[^45]: [Stanley Razor Blog: Developing and Backtesting Quantitative Strategies](https://stanleyrazor.github.io/portfolio-blog/posts/developing-and-backtesting-quantitative-strategies/) — hypothesis-first approach | practitioner blog | verified: 2026-06-21

[^46]: [Exegy: How to Stop Alpha Decay](https://www.exegy.com/alpha-decay/) — decay mechanisms, mitigation | practitioner research | verified: 2026-06-21

[^47]: [Medium — BT: The Truth About Automated Forex Strategies — Alpha Decay](https://turmanaube.medium.com/the-truth-about-automated-forex-strategies-alpha-decay-and-the-hunt-for-the-edge-d7782f4a3acd) — decay timeline, warning signs | practitioner blog | verified: 2026-06-21

[^48]: [QuantifiedStrategies: Discretionary vs Systematic Trading](https://www.quantifiedstrategies.com/discretionary-vs-systematic-trading/) | [QuantPedia: Combining Discretionary and Algorithmic Trading](https://quantpedia.com/combining-discretionary-and-algorithmic-trading/) — hybrid approaches | educational | verified: 2026-06-21

[^49]: [UBP: Systematic Multi-Strategy as a Portfolio Diversifier](https://www.ubp.com/en/news-insights/newsroom/systematic-multi-strategy-as-a-portfolio-diversifier/) — inter-strategy correlation data | asset manager research | verified: 2026-06-21

[^50]: [Goldman Sachs AM: How Quant Strategies Drive the Alpha-Enhanced Approach](https://am.gs.com/en-fi/advisors/insights/article/2025/how-quant-strategies-drive-the-alpha-enhanced-approach-to-equity-investing) — multi-strategy portfolio construction | asset manager research | verified: 2026-06-21
