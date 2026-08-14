<!-- Provenance: reference under the `options-trading-and-strategies` spoke (trading-and-investing hub family). Created 2026-06-16 via /dr deep-research. Educational only — NOT financial/investment/tax advice. -->

# The Option Greeks (Risk Sensitivities)

US listed equity options, educational, **as of 2026** — **not advice**. The Greeks are the (theoretical) risk sensitivities of an option's price to its inputs. Conceptual finance, largely **stable**. All per-contract dollar figures assume the standard **100-share multiplier**.

## Contents

- [Sign-convention master key](#sign-convention-master-key)
- [Delta — directional sensitivity](#delta--directional-sensitivity)
- [Gamma — how fast delta changes](#gamma--how-fast-delta-changes)
- [Theta — time decay](#theta--time-decay)
- [Vega — implied-volatility sensitivity](#vega--implied-volatility-sensitivity)
- [Rho — interest-rate sensitivity](#rho--interest-rate-sensitivity)
- [Position / portfolio Greeks](#position--portfolio-greeks)
- [Second-order Greeks (naming only)](#second-order-greeks-naming-only)
- [The central tradeoff: long vs short premium](#the-central-tradeoff-long-vs-short-premium)
- [References](#references)

## Sign-convention master key

Load-bearing and confirmed across sources: for a **long** option (call or put), regardless of type: **+gamma, +vega, −theta**. A **short** option flips all three: **−gamma, −vega, +theta**. **Delta**'s sign depends on type and direction (below).[^oa-gamma-sign][^tradingblock-gamma]

## Delta — directional sensitivity

**What it measures:** the theoretical change in an option's premium per **$1 move in the underlying**. The OIC: delta is "a theoretical estimate of how much an option's premium may change given a $1 move in the underlying."[^oic-delta]

- **Ranges:** long **call** delta runs **0 → +1.00**; long **put** delta runs **0 → −1.00**. Long calls = +delta, short calls = −delta; long puts = −delta, short puts = +delta. ATM ≈ ±0.50; deep ITM → ±1.00 (the option moves nearly dollar-for-dollar with the stock); deep OTM → ~0.[^oic-delta][^cme-greeks]
- **Hedge ratio / share-equivalent:** delta is also the **hedge ratio** — the "share-equivalent" exposure. A +0.40 call behaves like owning ~40 shares; to delta-hedge one such long call you would short ~40 shares.[^oic-delta]
- **The "approximate probability of finishing ITM" — and the caveat:** traders commonly read delta as the **approximate** probability of expiring ITM (a 0.30-delta option ≈ ~30% chance).[^oic-delta][^schwab-delta] **The caveat is real:** delta is *not* the true probability. It corresponds to a risk-neutral construct and is not exactly N(d₂) (the actual risk-neutral exercise probability), and the *real-world* probability depends on the asset's true drift, not the risk-free rate. It is close in practice — especially for near-term options, which is why the shortcut is popular — but state it as an **approximation**, never a true probability.[^macroption-delta-prob]

> **Worked example.** A call with delta **0.40** gains ~**$0.40** per $1 up-move in the stock → **~$40 per contract** (0.40 × $1 × 100). A put with delta **−0.40** *loses* ~$0.40 (−$40/contract) if the stock rises $1.

## Gamma — how fast delta changes

**What it measures:** the **rate of change of delta per $1 move** in the underlying — the option's convexity, or "the delta of the delta."[^schwab-gamma]

- **Where it's highest:** **ATM and near expiration**. A small move can flip an ATM option's moneyness, so its delta reacts most; gamma → ~0 deep ITM/OTM. ATM gamma spikes dramatically into the final days (flattish ~30 DTE, explosive ~1 DTE).[^schwab-gamma][^macroption-2nd]
- **Long vs short:** **long options are long (positive) gamma** — deltas move in your favor (rising as the stock rises, falling as it falls). **Short options are short (negative) gamma** — deltas move against you.[^schwab-gamma][^tradingblock-gamma]
- **Short-gamma risk:** negative gamma near expiration is among the highest-risk option exposures. A stock oscillating around the short strike forces a delta-hedger to **buy high / sell low** while rebalancing, bleeding money — the core of the "steamroller" intuition below.[^schwab-gamma]
- **Nuance (qualify):** "gamma is highest ATM" holds in normal/low-vol conditions, but the **peak can shift** off-ATM: as IV rises, ATM gamma *falls* while OTM/ITM gamma rises, flattening the curve. The clean ATM peak is most pronounced at low vol and short DTE.[^macroption-2nd]

> **Worked example.** A call with delta 0.50 and gamma **0.05**: after a $1 up-move, new delta ≈ **0.55**; after a $1 down-move ≈ **0.45**. Delta itself moves 5 points per $1.

## Theta — time decay

**What it measures:** the dollar value an option loses **per calendar day**, all else equal. The OIC: theta "represents, in theory, how much an option's premium may decay per day with all other pricing factors remaining the same."[^oic-theta]

- **Decay is NON-LINEAR and accelerates.** The OIC is explicit: "Theta or time decay is not linear. The theoretical rate of decay will tend to increase as time to expiration decreases."[^oic-theta] ATM time value decays roughly with the **square root of remaining time** (a 2-month ATM ≈ √2 ≈ 1.41× a 1-month ATM), so decay steepens sharply into expiration; the final weeks carry the bulk of an ATM option's total time decay.[^oa-theta][^days-theta]
- **Long vs short:** **long options pay theta (negative)** — decay erodes the premium you hold. **Short options collect theta (positive)** — you keep premium as it decays. ATM options carry the most theta exposure (most extrinsic value to lose); deep ITM/OTM decay slowly.[^oic-theta]

> **Worked example.** The OIC's case: $50 stock, $50-strike call at $3 with theta **0.05** loses ~**$0.05/day** (~$5/contract), all else equal.[^oic-theta] Decay is not constant — an option might shed ~$3/day at 60 DTE but ~$12/day at 7 DTE.

## Vega — implied-volatility sensitivity

**What it measures:** premium change per **1 percentage-point (1 vol point / 100 bps) change in implied volatility**. The OIC: vega "measures the amount of increase or decrease in premium based on a 1% (100 basis points) change in the implied volatility assumption."[^oic-vega] (Vega is not a true Greek letter but is universally grouped with the Greeks.)

- **Largest for ATM and longer-dated options.** The OIC: "Longer-term options tend to have higher Vega than near-term options"; ATM options are most vega-sensitive.[^oic-vega]
- **Long vs short:** **long options are long vega** (gain when IV rises); **short options are short vega** (gain when IV falls). A correct directional move can still lose if IV collapses enough to offset it (the IV-crush lesson in `volatility-and-pricing.md`).[^tradingblock-gamma]

> **Worked example.** The OIC's case: a 12-month call, IV 30%, **vega 0.15**, value $4; IV rises 2 points to 32% → premium rises 0.15 × 2 = **$0.30** to ~$4.30 (≈ +$30/contract).[^oic-vega]

## Rho — interest-rate sensitivity

**What it measures:** premium change per **1 percentage-point change in interest rates**. The OIC: rho is "the measure of an option's sensitivity to interest rate changes."[^oic-rho] Calls have **positive** rho (higher rates → higher call premiums); puts have **negative** rho.[^oic-rho]

- **Usually the least significant for short-dated equity options; matters more for LEAPS.** Rate changes "impact longer-term options much more than near-term ones."[^oic-rho]
- **Don't overstate "irrelevant" (qualify).** Rho is **not** universally negligible. In the higher-rate regime of 2022–2024 its significance rose, and it matters more for **longer-dated and ITM** options. The accurate framing is "least significant *for short-dated equity* options," not "always negligible."[^schwab-rates]

> **Worked example.** The OIC's case: a $100 call with rho **+0.45**, rates 3% → 4% → premium rises ~**$0.45**.[^oic-rho]

## Position / portfolio Greeks

The Greeks are **additive**: a multi-leg position's net Greek = the **quantity-weighted, signed sum** of each leg's Greek. So a spread has a **net delta, net gamma, net theta, net vega**. **Position delta** expresses the whole position's directional exposure as a share-equivalent (e.g., net +250 delta ≈ long 250 shares). A bull call spread, for instance, has lower net delta than the naked long call because the short leg offsets part of the long leg's Greeks. Traders manage the portfolio's aggregate "Greek profile" — e.g., trimming net vega before an earnings event, or rebalancing net delta toward neutral.[^tradingblock-gamma]

## Second-order Greeks (naming only)

Sensitivities **of the Greeks themselves** (second-order). At awareness depth only — not derived here:

- **Vanna** — sensitivity of **delta to a change in IV** (equivalently, vega to a change in spot).[^macroption-2nd]
- **Charm** ("delta decay") — sensitivity of **delta to the passage of time**.[^macroption-2nd]
- **Vomma** (a.k.a. volga) — sensitivity of **vega to a change in IV**.[^macroption-2nd]

They matter most to advanced/professional desks and affect OTM and shorter-dated options most.

## The central tradeoff: long vs short premium

- **Long premium (net buyer):** long gamma + long vega, but **pays theta**. You profit from large/fast moves and from rising IV; you bleed time value in a quiet, range-bound tape.[^tradingblock-gamma]
- **Short premium (net seller):** **collects theta**, but is short gamma + short vega. Steady income in calm markets; exposed to accelerating losses on a sharp move (negative gamma) or a vol spike (negative vega).[^tradingblock-gamma]
- **"Picking up pennies in front of a steamroller":** the canonical description of short-premium / short-gamma selling — the small, steady theta "pennies" can be wiped out by a single large adverse move ("the steamroller"), because short gamma makes losses accelerate, worst of all near expiration in short-dated/weekly options.[^theta-steamroller]

## References

Authoritative options education (OIC) first; exchange (CME) and brokerage (Schwab) corroborate; rigorous secondary sources (Macroption) ground the probability and second-order nuances.

[^oic-delta]: OIC, "Delta," https://www.optionseducation.org/advancedconcepts/delta — authoritative options education. verified-as-of: 2026-06-16
[^oic-theta]: OIC, "Theta," https://www.optionseducation.org/advancedconcepts/theta — options education (non-linear decay). verified-as-of: 2026-06-16
[^oic-vega]: OIC, "Vega," https://www.optionseducation.org/advancedconcepts/vega — options education. verified-as-of: 2026-06-16
[^oic-rho]: OIC, "Rho," https://www.optionseducation.org/advancedconcepts/rho — options education. verified-as-of: 2026-06-16
[^cme-greeks]: CME Group, "Option Delta / The Greeks" education, https://www.cmegroup.com/education/courses/introduction-to-options.html — exchange education (delta ranges; corroborating). verified-as-of: 2026-06-16
[^schwab-delta]: Schwab, "Options Delta, Probability and Other Risk Analytics," https://www.schwab.com/learn/story/options-greeks-delta — brokerage education (probability-proxy framing). verified-as-of: 2026-06-16
[^macroption-delta-prob]: Macroption, "Delta of Calls/Puts and Probability of Expiring ITM" (delta ≠ true ITM probability), https://www.macroption.com/delta-calls-puts-probability-expiring-itm/ — rigorous secondary. verified-as-of: 2026-06-16
[^schwab-gamma]: Schwab, "Options Greeks: Gamma Explained," https://www.schwab.com/learn/story/gamma-options — brokerage education. verified-as-of: 2026-06-16
[^macroption-2nd]: Macroption, "Second Order Greeks" and gamma-peak nuance, https://www.macroption.com/second-order-greeks/ — secondary education. verified-as-of: 2026-06-16
[^tradingblock-gamma]: TradingBlock, "Long Gamma vs Short Gamma" and long/short sign conventions, https://www.tradingblock.com/blog/long-gamma-vs-short-gamma — brokerage education (sign conventions). verified-as-of: 2026-06-16
[^oa-gamma-sign]: Option Alpha, options Greeks sign conventions (long = +gamma/+vega/−theta), https://optionalpha.com/learn/option-greeks — education. verified-as-of: 2026-06-16
[^oa-theta]: Option Alpha, "0DTE Options Time Decay" (accelerating theta), https://optionalpha.com/blog/0dte-options-time-decay — education. verified-as-of: 2026-06-16
[^days-theta]: Days to Expiry, "Theta Decay Accelerates Near Expiration" (√time decay), https://daystoexpiry.com/ — education (corroborating). verified-as-of: 2026-06-16
[^schwab-rates]: Schwab, "How Interest Rate Movements Affect Options Prices" (rho relevance in higher-rate regimes / for LEAPS), https://www.schwab.com/learn/story/how-interest-rate-movements-affect-options-prices — brokerage education. verified-as-of: 2026-06-16
[^theta-steamroller]: ThetaTrend, "Picking Up Pennies in Front of the SPX Steamroller," https://thetatrend.com/ — practitioner education (concept corroboration). verified-as-of: 2026-06-16
