# trading-and-investing

**Category:** Security, Auth & Diagnostics
**Platform:** Claude
**Original Path:** claude/standalone/trading-and-investing

## Description
FOUNDATION + HUB for active trading & how US-retail financial markets work; routes to 16 spokes. Educational ONLY, NOT financial/investment/tax advice; trading can lose money (2026). TRIGGER: how markets work; asset classes (stocks, bonds, forex, commodities, crypto, ETFs, derivatives); participants (market makers, broker-dealers, exchanges, clearinghouses); order routing (PFOF, NBBO, T+1) & order types; sessions & circuit breakers; investing-vs-trading; trading styles; margin/short-selling/PDT; checking a trading broker; DEX/AMM/DeFi/on-chain trading; impermanent loss; MEV/sandwich attacks; ML/AI for trading; feature engineering; alt-data signals; RL for execution. SKIP: long-term passive/retirement investing (401k/IRA/Roth, index buy-and-hold, Social Security, fiduciary) → investing-and-retirement; budgeting → budgeting-and-saving; bank/CD/FDIC → personal-banking; income-tax → personal-income-taxes; investment-fraud recovery → identity-theft-and-credit-fraud.

---

# Trading & Investing (foundation + hub)

The front door for **active trading and how financial markets actually work** for a US retail participant. This skill does two jobs:

1. **Foundation** — it carries the shared overview every trading sub-topic builds on: the asset classes, the market participants, primary vs secondary markets, how an order travels from your broker to the market, market sessions and hours, the core investing-vs-trading distinction, and the risk/regulatory backbone. Deep treatments live in `references/` (loaded on demand).
2. **Hub / router** — it routes specific questions down to the family's 16 spokes (see the routing table). The spokes are the *intended* family; as each is built it owns its depth and this hub just points to it.

> **Educational information only — NOT financial, investment, tax, or legal advice, and NOT a recommendation to buy, sell, or hold anything.** All securities and trading carry **risk of loss**; you can lose money, and with leverage you can lose **more than you put in**. Past performance does not guarantee future results. **The large majority of active retail/day traders underperform a simple index or lose money** (see `references/investing-vs-trading.md`). Rules, products, tax treatment, and the vendor landscape change — every volatile claim here is stamped **as of 2026**; verify current facts with the primary regulator (SEC/Investor.gov, FINRA, CFTC, SIPC) before acting. For a personal decision, consult a licensed professional.

## Sibling skill — the passive/retirement seam (read this first if unsure)

This family is the **active** side: markets, instruments, execution, trading styles, and the risks of trading. Its sibling **`investing-and-retirement`** (in the consumer-finance family) owns the **passive, long-term, retirement** side — low-cost index buy-and-hold, 401(k)/403(b)/IRA/Roth mechanics and contribution limits, employer match and vesting, dollar-cost averaging for retirement, Social Security claiming, target-date defaults, and choosing a robo-advisor or fee-only fiduciary.

- "How do I start investing for retirement / which account / is this fund's expense ratio good / should I do a backdoor Roth / when do I claim Social Security" → **`investing-and-retirement`**.
- "How do markets work / how do I trade X / what's a limit order / is day trading worth it / how does margin work / what is options-/futures-/forex-/crypto-trading" → **this family**.

The two overlap on shared vocabulary (diversification, asset allocation, ETFs vs mutual funds, fees). When the *intent* is long-term passive wealth-building, defer to `investing-and-retirement`; when the intent is understanding markets or actively trading, stay here. `portfolio-theory-and-asset-allocation` (a spoke here) covers the *theory* (MPT, efficient frontier, factor models, rebalancing math); `investing-and-retirement` covers the *consumer how-to*.

## Routing table — the 16 spokes

> Spokes are the **intended family**. **Built** = reference file exists in `references/` and is listed in the Foundation references table below. **Unbuilt** = hub answers from this foundation until the spoke is built; check the available-skills list for a standalone installed skill of that slug first.

| Spoke | Build state | Route here when the question is about… |
| --- | --- | --- |
| **stock-and-equity-trading** | Built | Trading individual stocks/equities and ETFs as instruments: share mechanics, order placement in practice, equity-specific tactics, IPOs/secondary offerings from a trader's seat, dividends/splits, equity screening. |
| **options-trading-and-strategies** | Built (also standalone skill) | Listed options: calls/puts, strike/expiration, the Greeks, single-leg and multi-leg strategies (spreads, straddles, covered calls, iron condors), assignment, options pricing. |
| **derivatives-futures-and-swaps** | Unbuilt | Futures, forwards, and swaps (non-options derivatives): contract specs, margin/mark-to-market, hedging vs speculation, interest-rate/FX/commodity/credit swaps, the broader derivatives landscape. |
| **technical-analysis** | Built | Reading price/volume: chart patterns, candlesticks, trend, support/resistance, indicators (RSI, MACD, moving averages, Bollinger), Dow theory, Elliott wave, market breadth. |
| **trading-strategies-and-styles** | Unbuilt | The trader's playbook by style: position/swing/day-trading/scalping mechanics, momentum/mean-reversion/breakout/trend-following systems, entries/exits, backtesting a discretionary strategy. |
| **algorithmic-and-quant-trading** | Unbuilt | Systematic/automated trading: strategy coding, backtesting frameworks, execution algos, market-making/stat-arb, signals, latency, the quant research workflow. |
| **crypto-and-digital-asset-trading** | Unbuilt | Trading crypto on centralized venues: spot vs perps/futures, exchanges and custody, stablecoins, spot crypto ETFs, order books, on-exchange mechanics (CEX). |
| **forex-and-currency-trading** | Unbuilt | The FX market: currency pairs, pips/lots, leverage, carry, sessions, majors/minors/exotics, retail spot-FX mechanics and the OTC structure. |
| **trading-risk-management** | Unbuilt | Sizing and protecting trading capital: position sizing, stop placement, risk-per-trade, R-multiples, drawdown control, Kelly, portfolio heat, margin/leverage management as a discipline. |
| **portfolio-theory-and-asset-allocation** | Unbuilt | The theory: Modern Portfolio Theory, the efficient frontier, CAPM/factor models, diversification math, correlation, rebalancing, strategic vs tactical allocation. |
| **trading-psychology-and-behavioral-finance** | Unbuilt | The mind: discipline, fear/greed, loss aversion, the disposition effect, overconfidence, prospect theory, cognitive biases, trading-journal and routine practice. |
| **fixed-income-and-bond-markets** | Unbuilt | Bonds in depth: the yield curve, duration/convexity, credit spreads, Treasuries/corporates/munis, bond pricing, rate risk, the repo/money markets. |
| **market-microstructure-and-execution** | Unbuilt | How trades really execute: order book dynamics, maker-taker, Reg NMS internals, dark pools/ATSs, latency, slippage, TWAP/VWAP, sub-penny price improvement, SIP vs direct feeds. |
| **defi-and-onchain-trading** | Built | On-chain/decentralized trading: DEXs, AMMs, liquidity pools, swaps, yield/staking, wallets, MEV, bridges — the on-chain counterpart to centralized crypto. |
| **ai-and-ml-for-trading** | Built | Machine learning applied to markets: predictive models, feature engineering on market data, sentiment/alt-data, reinforcement learning for execution, ML backtesting pitfalls (overfitting, look-ahead). |
| **trading-regulation-compliance-and-taxes** | Unbuilt | The rulebook & the IRS: SEC/FINRA/CFTC regimes, PDT/margin rules, wash sales, trader-tax-status, capital-gains mechanics for active traders, reporting, recordkeeping. |

## Foundation references (load on demand)

The shared foundation is split into focused references so a spoke can cross-reference one without pulling the whole hub:

| Topic | Covers | Reference |
| --- | --- | --- |
| `asset-classes-and-instruments` | The major asset classes a retail participant can access — equities, fixed income, forex, commodities, crypto/digital assets, pooled vehicles (ETFs vs mutual vs index vs target-date), and derivatives at overview depth. What each *is* and how retail accesses it. | `references/asset-classes-and-instruments.md` |
| `market-participants-and-structure` | The cast (retail, institutions, market makers, broker-dealers, exchanges, clearinghouses/CCPs), what a clearinghouse does and why it cuts counterparty risk (DTCC/NSCC, OCC), settlement & the T+1 cycle, primary vs secondary markets, exchange vs OTC, and the bid-ask spread. | `references/market-participants-and-structure.md` |
| `order-lifecycle-and-execution` | How a retail order travels broker→venue→execution→clearing→settlement; order types (market/limit/stop/stop-limit) and each one's risk; payment for order flow (PFOF) and the best-execution duty; the NBBO; the $0-commission landscape. | `references/order-lifecycle-and-execution.md` |
| `market-sessions-and-venues` | Where things trade (NYSE/Nasdaq/CME/Cboe/OTC), US equity sessions & extended hours, the 2026 move toward overnight/24-hour equity trading, how forex/crypto/futures hours differ, market holidays, and circuit breakers (market-wide + LULD). | `references/market-sessions-and-venues.md` |
| `investing-vs-trading` | The central distinction — time horizon, goal, what drives decisions, the trading-styles spectrum, active vs passive, the high-level tax angle, and the **load-bearing evidence** that most active traders underperform (Barber & Odean; Taiwan & Brazil day-trader studies). | `references/investing-vs-trading.md` |
| `trading-risks-and-protections` | The risk/regulatory backbone — the Pattern Day Trader rule **and** its 2026 replacement, margin/leverage, short selling, leveraged forex/crypto/derivatives, structural & behavioral risks, and investor protection (SIPC, BrokerCheck, the SEC mandate, the standard disclaimer). | `references/trading-risks-and-protections.md` |
| `stock-and-equity-trading` | Stock & equity trading spoke — order handling, equity-specific mechanics, and how the foundation applies to single-name and ETF equity trades. | `references/stock-and-equity-trading.md` |
| `options-trading-and-strategies` | Listed equity/ETF/index options for a US retail participant — calls/puts, the Greeks, implied vol, and defined- vs undefined-risk strategies (spreads, straddles, condors). | `references/options-trading-and-strategies.md` |
| `defi-and-onchain-trading` | On-chain/decentralized trading: AMM x·y=k formula, concentrated liquidity (v3), impermanent loss math, price impact/slippage, MEV sandwich attacks and defense (private RPC, intent protocols), CoW Protocol/UniswapX/1inch Fusion, perpetual DEXs (GMX/dYdX/Hyperliquid), cross-chain bridges (lock-and-mint, burn-and-mint, bridge risks), self-custody wallets (MetaMask, Rabby, Ledger, Trezor, ERC-4337 AA). | `references/defi-and-onchain-trading.md` |
| `ai-and-ml-for-trading` | Machine learning applied to markets: feature engineering on OHLCV/microstructure data (info bars, fractional differentiation, PIT normalization), gradient boosting for cross-sectional equity signals (Gu/Kelly/Xiu, SHAP, purged CV), neural networks for financial time series (918-experiment architecture comparison), alternative data taxonomy (satellite, transaction, NLP), LLMs in finance (FinBERT, BloombergGPT, hallucination risk), RL for execution (Nevmyvaka, Almgren-Chriss baseline, RL failure modes), ML backtesting pitfalls (look-ahead bias taxonomy, HLZ multiple testing, McLean & Pontiff factor decay, PBO), production systems (IC monitoring, drift detection, layered risk controls). | `references/ai-and-ml-for-trading.md` |
| `technical-analysis` | Reading price and volume: chart types (candlestick/OHLC/line), candlestick patterns (doji/hammer/engulfing/morning star), chart patterns (head-and-shoulders, double top/bottom, triangles, flags, wedges), trend analysis and trendlines, support and resistance levels, key indicators (RSI, MACD, Bollinger Bands, SMA/EMA, ATR, stochastics), volume indicators (OBV, VWAP, Volume Profile), market breadth (A/D line, McClellan Oscillator, Zweig Breadth Thrust), Dow Theory (6 tenets, 3 trend types, 3 phases), Elliott Wave (5-wave impulse, A-B-C correction, Fibonacci guidelines), multi-timeframe analysis, empirical evidence and academic critique (Fama/EMH, BLL 1992, Sullivan et al. 1999, Park & Irwin 2007, Lo et al. 2000, Adaptive Market Hypothesis). | `references/technical-analysis.md` |

## How to answer from this hub (output contract)

When answering directly from this foundation: (1) for any advice-seeking question, **lead with the educational-only / not-advice framing** and that trading carries risk of loss; (2) **cite the primary regulator** (SEC/Investor.gov, FINRA, CFTC, SIPC) for load-bearing facts, and stamp volatile facts *as of 2026* with a "verify current" pointer; (3) when the question belongs to a spoke whose **reference file exists** in the Foundation references table above, `Read` its `references/<slug>.md` and answer from it — if a standalone installed skill also exists for that spoke, the standalone skill takes precedence; if the Read fails or the reference is silent on the specific question, degrade to foundation depth; (4) when the question belongs to a spoke **not yet in the Foundation references table**, first check the available-skills list for a standalone installed skill of that spoke's slug — if found, invoke it; otherwise name the destination spoke and answer at foundation depth; (5) when a query spans two spokes, identify the **primary spoke** (the one that owns the action) and answer from it, citing the secondary spoke for supplementary detail. Keep US-centric, flagging where a fact (T+1, PDT, SIPC, Reg NMS, CFTC caps) is US-only.

## First, the account (the gate before anything else)

Everything downstream runs through a **brokerage account** opened with a broker-dealer (after identity verification, funded from a bank). The one account choice that gates the rest is **cash vs margin**: a **cash account** trades only settled funds (no borrowing); a **margin account** lets the broker **lend you money** against the account as collateral — which is what enables short selling and triggers the margin and Pattern-Day-Trader rules in `references/trading-risks-and-protections.md`. Account *type* (individual/joint/retirement-wrapper) sits on top of that; a retirement-wrapper account points to the passive sibling, `investing-and-retirement`.

## The one-paragraph foundation (if you only read the hub)

A **retail participant** buys and sells **instruments** (equities, bonds, FX, commodities, crypto, pooled funds like ETFs/mutual funds, and derivatives) through a **broker-dealer**, which routes the order to an **execution venue** (an exchange like NYSE/Nasdaq, a wholesaler/market maker, or an ATS); the trade then **clears** through a central counterparty (NSCC for US equities, OCC for listed options) and **settles** the next business day (**T+1**, as of 2026). The issuer raises money only once, in the **primary market** (an IPO or new issue); everything after is the **secondary market**, where investors trade among themselves and prices are discovered. **Investing** means holding for the long term to build wealth through compounding; **trading** means buying and selling frequently to profit from price moves — and the evidence is strong that **most active retail traders underperform or lose money**, so trading capital should be money you can afford to lose. Markets are regulated (SEC, FINRA, CFTC) and your *brokerage account* is protected against broker failure by **SIPC** — but **nothing protects you against market losses**.

## Cross-cutting notes

**`as of 2026` volatile claims to re-verify** (detail in the references) — one line each, highest-stakes first:

- **Pattern Day Trader rule is being replaced** (highest-stakes): FINRA amended Rule 4210 to drop the $25k/PDT designation for a new **intraday-margin** framework, **effective June 4, 2026** (transition through Oct 20, 2027). Most third-party sites still describe the legacy $25k rule — verify your broker's *current* policy.
- **Settlement is T+1** (since May 28, 2024); **SIPC** limits are **$500k / $250k cash** — both change, verify at the regulator.
- **Spot Bitcoin ETFs** (approved Jan 2024) and **spot Ether ETFs** (Jul 2024) exist and have grown large.
- **Nasdaq 23-hour trading** was SEC-approved (Apr 10, 2026), launch targeted H2 2026; **CME launched 24/7 crypto futures** (May 29, 2026).
- US retail **forex leverage** is capped (~50:1 majors / ~20:1 others) by CFTC/NFA; **$0 stock/ETF commissions** are standard; capital-gains breakpoints change yearly.

Other seams:

- **Commodities routing.** Commodities are an asset class in the foundation, but there is **no dedicated commodities spoke** — route commodity *futures* depth to `derivatives-futures-and-swaps` and commodity *ETP/ETF* depth to `stock-and-equity-trading`.
- **Tax seam.** This family's `trading-regulation-compliance-and-taxes` spoke owns the *trader's* tax detail (wash sales, trader-tax-status, active-trader capital-gains mechanics). The consumer-finance family's `personal-income-taxes` owns ordinary individual filing. Keep deep tax out of the foundation — it states only the short-vs-long-term-gains distinction at a high level.
- **Fraud seam.** Spotting an investment scam and the "is this a Ponzi / guaranteed-returns" pattern is shared with `investing-and-retirement`; recovery *after* a scam → `consumer-credit-and-debt` (references/identity-theft-and-credit-fraud.md).
- **Blockchain seam.** Trader-side DeFi (executing swaps on DEXs, managing liquidity positions, MEV defense, perp DEX trading, bridge mechanics for moving assets) is owned here in `defi-and-onchain-trading`. Protocol mechanics, cryptoeconomics, and the academic/research layer (AMM math derivation, MEV supply-chain/PBS theory, tokenomics, consensus) → `blockchain` / `blockchain-economics`.
- **Data/quant seam.** The math/ML *techniques* behind quant and ML-for-trading (regression, time-series, feature engineering, backtesting statistics) draw on the `da-*` data-analysis hubs; the trading spokes own the market application, the `da-*` hubs own the method.
- **International note.** This family is **US-centric**. Other major markets (LSE, Euronext, Tokyo, Hong Kong, Shanghai) run their own hours, regulators, settlement cycles, and tax regimes; the foundation flags where the US specifics (T+1, PDT, SIPC, Reg NMS, CFTC leverage caps) are US-only and would differ abroad.

## References

Each foundation reference carries its own cited `## References` / `## Sources` section anchored to primary authorities — **SEC / Investor.gov, FINRA, CFTC, SIPC, DTCC, OCC, the Federal Reserve / BIS, S&P Dow Jones Indices (SPIVA)**, and the peer-reviewed household-finance literature (**Barber & Odean**; the **Taiwan** and **Brazil** day-trader studies). Volatile facts are dated *as of 2026* with "verify current" pointers. Start from the relevant reference above; as the 16 spokes are built, each will carry its own deeper citations.

<!-- cross-hub-map -->
## Cross-hub map — where every trading-and-investing topic lives

All active-trading material lives under this single hub. Built reference files are listed below; for unbuilt spokes, check whether a standalone installed skill exists for that slug and invoke it, or answer at foundation depth.

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `trading-and-investing` | Trading & Investing — all active trading, markets, and risk topics | `references/asset-classes-and-instruments.md`, `references/market-participants-and-structure.md`, `references/order-lifecycle-and-execution.md`, `references/market-sessions-and-venues.md`, `references/investing-vs-trading.md`, `references/trading-risks-and-protections.md`, `references/stock-and-equity-trading.md`, `references/options-trading-and-strategies.md`, `references/defi-and-onchain-trading.md`, `references/ai-and-ml-for-trading.md`, `references/technical-analysis.md` |