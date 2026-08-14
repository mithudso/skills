<!-- Provenance: reference under the `options-trading-and-strategies` spoke (trading-and-investing hub family). Created 2026-06-16 via /dr deep-research. Educational only — NOT financial/investment/tax advice. -->

# Volatility & Conceptual Options Pricing

US listed equity/index options, educational, **as of 2026** — **not advice**. Black-Scholes-Merton and put-call parity are presented as **conceptual scaffolding / intuition**, NOT formulas to derive or memorize. The *conceptual* relationships are stable; any live VIX level, skew steepness, or volatility-risk-premium magnitude is market-state-dependent and dated *as of 2026*.

## Contents

- [Implied volatility (IV)](#implied-volatility-iv)
- [Historical / realized volatility and the IV–HV gap](#historical--realized-volatility-and-the-ivhv-gap)
- [IV rank vs IV percentile](#iv-rank-vs-iv-percentile)
- [The volatility smile, skew, term structure & surface](#the-volatility-smile-skew-term-structure--surface)
- [Vega risk: the IV crush](#vega-risk-the-iv-crush)
- [The VIX](#the-vix)
- [Black-Scholes-Merton — intuition only](#black-scholes-merton--intuition-only)
- [Put-call parity — conceptual](#put-call-parity--conceptual)
- [References](#references)

## Implied volatility (IV)

**Implied volatility** is the volatility figure *backed out* of an option's current market price using a pricing model: you take the observed price and solve the model in reverse for the volatility input that reproduces it. IV therefore **cannot be directly observed** — it is inferred.[^schwab-iv][^wiki-iv][^optionsplaybook-vol] It is **forward-looking** (the market's collective estimate of *future* volatility over the option's life), expressed as an **annualized percentage**, and it moves continuously with supply and demand.[^schwab-iv][^wiki-iv]

- **How IV moves premium ("buying/selling vol"):** higher IV → richer (more expensive) premiums; lower IV → cheaper premiums, holding strike, spot, and time constant. Because IV is the discretionary input traders most disagree on, options are often discussed as direct exposure to volatility — buying premium is "buying vol," writing premium is "selling vol." (This is standard retail vernacular; **vega** is the precise sensitivity — see `the-greeks.md`.)[^schwab-iv]
- **Disconfirming check — IV is NON-directional.** A persistent retail misconception is that IV predicts *which way* the stock will move. IV estimates the expected **magnitude** (standard deviation) of movement, not direction.[^ebc-iv][^strike-iv] The practical consequence is the IV crush below: "right on direction, still lost."

## Historical / realized volatility and the IV–HV gap

**Realized (a.k.a. historical) volatility (RV/HV)** is **backward-looking**: the actual standard deviation of the underlying's past price moves over a chosen window (e.g., 20/30 days). It is computed directly from price history, needs no model, and is a measured fact rather than an inference.[^luxalgo-rviv][^macroption-vol] IV (forward-looking, model-derived) and RV (backward-looking, measured) answer **different questions**.

- **IV > RV** (the usual state): options are "expensive" relative to how much the stock has actually been moving — the market prices in more future movement than recently occurred. Premium *sellers* favored.[^luxalgo-rviv][^barchart-ivrange]
- **IV < RV**: options look "cheap" relative to realized movement. Premium *buyers* favored.[^luxalgo-rviv]
- **Volatility risk premium (VRP), qualify / as-of-2026:** IV tends to *exceed* subsequent realized volatility on average, because option buyers pay up for protection and sellers demand compensation for tail risk. One source quantifies this as IV exceeding RV by ~3–5 percentage points for large caps on average[^luxalgo-rviv] — treat the **specific figure as illustrative and regime-dependent**, not a constant.

## IV rank vs IV percentile

Both contextualize *today's* IV against that **same underlying's** trailing ~52 weeks (≈252 trading days) of IV — **never across different stocks** (45% IV is "high" for a utility, "low" for a biotech).[^tasty-ivr][^barchart-ivrank]

- **IV Rank (IVR)** = where current IV sits within its 52-week **high–low range**: `(Current IV − 52wk Low) / (52wk High − 52wk Low) × 100`. If XYZ's IV ranged 30–60 over the year and is now 45, IVR = 50%.[^tasty-ivr][^barchart-ivrank]
- **IV Percentile (IVP)** = the **percentage of days** in the trailing year on which IV closed *below* today's level: `(# days IV < Current IV) / 252 × 100`. IVP 80% means IV was lower on 80% of the past year's days.[^tasty-ivr][^menthorq-iv]
- **Disconfirming check — they are NOT the same and diverge.** IVR uses only the two *extreme* values (so one outlier spike can distort it); IVP uses the full distribution. Textbook divergence: a new IV of 15 sitting exactly between a 10 low and 20 high gives **IVR = 50%**, but if 15 exceeds three of four past observations, **IVP = 75%**.[^barchart-ivrank][^menthorq-iv] IVP is generally considered more robust (not hostage to outliers); IVR is simpler and more common on platforms.
- **Use:** high IVR/IVP → options "expensive" → traders lean toward premium-*selling* (credit) strategies; low → "cheap" → premium-*buying* (debit). Specific thresholds (e.g., "sell when IVR > 50") are **house conventions, not rules** (tentative).[^tasty-ivr]

## The volatility smile, skew, term structure & surface

Black-Scholes assumes one constant volatility per underlying, but in reality **IV varies by strike**. Plotting IV against strike (same expiry) is rarely flat:[^oic-skew][^wiki-smile]

- **Volatility smile** — a U-shape where *both* OTM puts and OTM calls carry higher IV than ATM; common in **FX** options, where both tails are hedged.[^oic-skew][^wiki-smile]
- **Volatility skew / smirk** — an asymmetric curve. For **equity indices**, IV is *highest at low strikes (OTM puts)* and declines as strike rises ("reverse skew" / downside smirk).[^oic-skew][^luxalgo-rviv]

**Why equity OTM puts trade at higher IV than OTM calls (crash/put-demand skew):** (1) persistent **portfolio-insurance demand** — institutions buy OTM puts as hedges, bidding up put IV; (2) **crash/jump risk** — equities gap down faster than they rise, so puts price discontinuous downside. The OIC: a smirk "can be the result of investors protecting their portfolios against a declining market, which traditionally exhibits larger and more severe moves to the downside," and skew "is a product of the forces of supply and demand."[^oic-skew][^luxalgo-rviv]

> **Disconfirming / historical finding — the skew is NOT a law of nature; it appeared after 1987.** Before the Oct 19, 1987 crash (Dow −22% in one session), index-option smiles were "relatively flat"; the persistent downside skew emerged *afterward* as the market permanently repriced fat-tail/crash risk. This is the cleanest evidence against BSM's constant-vol assumption.[^wiki-smile][^oic-skew] Equity skew is also markedly steeper than single-stock or FX skews.

**Term structure & the surface:** IV also varies by *time to expiry* — plotting ATM IV across expirations is the **term structure of volatility**, normally upward-sloping ("contango") in calm markets and inverting to "backwardation" (near-term IV spikes above longer-dated) during stress. Combining *both* dimensions — IV as a function of **strike AND expiry** — gives the **volatility surface**, the 3-D map desks use to quote a consistent, arbitrage-free set of prices.[^macroption-vix][^simtrade-curves] (Conceptual level only; surface construction/SVI is out of scope.)

## Vega risk: the IV crush

Long-premium positions (bought calls/puts) are **long vega** — they *gain* when IV rises and *lose* when IV falls, **independent of direction**. The classic illustration is the **post-earnings IV crush**: before a scheduled event, uncertainty inflates IV; once the news is out, that uncertainty premium **collapses almost instantly**, deflating option prices.[^spotgamma-crush][^sofi-crush]

> **Worked example (earnings IV crush).** A trader buys an ATM call for **$12** with IV at **60%** ahead of earnings. The company beats; the stock gaps **+4% to $208** — the trader was *right on direction*. But IV collapses **60% → 30%**, and the call is now worth only **~$9**. Net: **−$300 per contract despite a correct directional call**, because vega losses (a 30-point IV drop) overwhelmed the delta gain. *Teaching point:* IV is non-directional, and *when* you buy vol matters as much as which way you think the stock goes.[^spotgamma-crush]

(Mechanics of vega itself — units, sign, decay toward expiry — are owned by `the-greeks.md`.)

## The VIX

The **Cboe Volatility Index (VIX)** is a real-time index of the market's expected **30-day** volatility of the **S&P 500 (SPX)**. In Cboe's own words, it provides "a reading of constant, 30-day expected volatility of the S&P 500 Index," computed from "real-time, mid-quote prices of a strip of S&P 500 options (SPX)," and its output is "a non-directional, annualized expectation for the standard deviation of the S&P 500."[^cboe-vix] Plain-language: a VIX of 16 implies the market expects roughly a ±16% annualized move in the S&P 500.[^wiki-vix] It aggregates the weighted prices of near-term and next-term SPX puts and calls plus the risk-free rate; popularly the **"fear gauge,"** since it spikes when demand for downside protection surges.[^wiki-vix]

- **Model-free (qualify):** the VIX is described as a **model-free** implied-volatility measure — it does not invert Black-Scholes per option but is computed from a variance-swap-style weighting of the whole option strip. Cboe's plain-language page confirms the 30-day/non-directional/annualized properties; the "model-free" characterization rests on Cboe's methodology PDF and Wikipedia. Ground it in the methodology PDF if the point is load-bearing.[^cboe-vix-method][^wiki-vix]

## Black-Scholes-Merton — intuition only

Present BSM as **conceptual scaffolding** — the foundational mental model for *what drives an option's value*, not a formula retail traders compute by hand (platforms do it).

- **The five/six inputs:** (1) underlying price S; (2) strike K; (3) time to expiry T; (4) **volatility σ** — the one input not directly observable, which is where IV enters; (5) risk-free rate r; (6) dividends q (the Merton extension; vanilla 1973 Black-Scholes assumed none).[^simtrade-bsm][^brilliant-bsm]
- **Core assumptions:** underlying follows geometric Brownian motion → **lognormal prices**; **constant volatility** (and constant r) over the life; **no arbitrage**; **European exercise**; **frictionless markets** (no costs/taxes, continuous trading, borrow/lend at r).[^macroption-bsm][^simtrade-bsm][^essex-matei]
- **Well-known real-world failures (disconfirming):**
  - **Constant-vol is contradicted by the smile/skew** — if BSM held, IV would be identical across all strikes; the observed skew is direct empirical refutation, and it sharpened after 1987.[^essex-matei][^wiki-smile]
  - **Fat tails** — real equity returns have **fatter tails than the lognormal** predicts; extreme moves happen more often than the model assumes, which is exactly why deep-OTM puts are bid up.[^essex-matei]
  - **Early exercise / American options** — vanilla BSM is built for **European** exercise and does **not** handle the early-exercise optionality of American options. Merton (1973) showed BSM prices an American *call* on a **non-dividend** stock (early exercise is never optimal there), but the general American case — and American puts — fall outside it. Most listed US equity options are American, making this a practical, not academic, limitation.[^simtrade-bsm][^essex-matei]
- **Binomial / lattice models for early exercise:** the **binomial (Cox-Ross-Rubinstein) lattice** handles American exercise naturally — at every node you compare "hold" vs "exercise now" and take the maximum, walking backward through the tree. As steps → ∞, the binomial price **converges to Black-Scholes**, so the lattice is both the intuition-builder and the practical workhorse for options BSM can't close-form.[^optiontradingtips-binomial][^crr-binomial] *Takeaway:* BSM gives the conceptual skeleton; lattice/numerical methods extend it to the American options retail traders actually hold.

## Put-call parity — conceptual

For options on the same underlying at the **same strike (K) and same expiry**, the prices of a call, a put, the stock, and a risk-free bond are locked together by no-arbitrage:[^analystprep-pcp][^macroption-pcp]

> **C + PV(K) = P + S** — (long call + cash equal to the discounted strike) = (long put + the stock).

What it implies (state it, don't prove it):

- Two portfolios with **identical expiration payoffs must cost the same today** — otherwise a riskless arbitrage exists (buy the cheap side, sell the rich side).[^analystprep-pcp]
- It lets you build **synthetic positions** — e.g., synthetic long stock = long call + short put at the same strike; synthetic call = stock + put. Market makers price one option and *derive* the other so their two-sided quotes never create arbitrage.[^oic-pcp]
- It ties option prices to the **risk-free rate and dividends** (carry), not just volatility.
- **Disconfirming check — exact ONLY for European options.** American early exercise relaxes the clean equality into an **inequality / bounds**; the call–put gap sits *within* bounds rather than at a single value. Since most listed US equity options are American, parity holds *approximately* and is the conceptual backbone, not a hard equality.[^millersville-american]

## References

Authoritative (Cboe, OIC) and academic (university lecture notes) sources anchor the load-bearing claims; rigorous secondary references (Macroption, AnalystPrep) and brokerage/education sources corroborate. Market-state magnitudes are dated *as of 2026*.

[^schwab-iv]: Schwab, "Aligning Your Options Strategies With Implied Volatility," https://www.schwab.com/learn/story/aligning-your-options-with-implied-volatility — brokerage education. verified-as-of: 2026-06-16
[^wiki-iv]: Wikipedia, "Implied volatility," https://en.wikipedia.org/wiki/Implied_volatility — encyclopedic (well-sourced). verified-as-of: 2026-06-16
[^optionsplaybook-vol]: The Options Playbook, "What is Volatility," https://www.optionsplaybook.com/options-introduction/what-is-volatility — options education. verified-as-of: 2026-06-16
[^ebc-iv]: EBC Financial, "Implied Volatility" (non-directional misconception), https://www.ebc.com/forex/implied-volatility — broker explainer. verified-as-of: 2026-06-16
[^strike-iv]: strike.money, "Implied Volatility (IV)," https://www.strike.money/options/implied-volatility-iv — education (corroborating non-directionality). verified-as-of: 2026-06-16
[^luxalgo-rviv]: LuxAlgo, "Realized Volatility vs Implied Volatility" and "Volatility Smile vs Skew," https://www.luxalgo.com/blog/realized-volatility-vs-implied-volatility-key-differences/ — trading education. verified-as-of: 2026-06-16
[^macroption-vol]: Macroption, "Implied vs Realized vs Historical Volatility," https://www.macroption.com/implied-vs-realized-vs-historical-volatility/ — options reference. verified-as-of: 2026-06-16
[^barchart-ivrange]: Barchart, "Implied Volatility Range," https://www.barchart.com/options/implied-volatility-range — data provider. verified-as-of: 2026-06-16
[^tasty-ivr]: tastytrade/tastylive, "IV Rank vs IV Percentile" support article, https://support.tastytrade.com/support/solutions/articles/43000539059 — popularized IVR/IVP. verified-as-of: 2026-06-16
[^barchart-ivrank]: Barchart Education, "IV Rank vs IV Percentile," https://www.barchart.com/education/iv_rank_vs_iv_percentile — data-provider education (divergence example). verified-as-of: 2026-06-16
[^menthorq-iv]: MenthorQ, "IV Rank vs Percentile," https://menthorq.com/guide/iv-rank-vs-percentile/ — trading education. verified-as-of: 2026-06-16
[^oic-skew]: OIC, "Volatility Skew and Options — An Overview," https://www.optionseducation.org/news/volatility-skew-and-options-an-overview-1 — authoritative options education. verified-as-of: 2026-06-16
[^wiki-smile]: Wikipedia, "Volatility smile" (post-1987 emergence; FX smile vs equity skew), https://en.wikipedia.org/wiki/Volatility_smile — encyclopedic. verified-as-of: 2026-06-16
[^macroption-vix]: Macroption, "VIX Futures Curve" (term structure/contango-backwardation), https://www.macroption.com/vix-futures-curve/ — options reference. verified-as-of: 2026-06-16
[^simtrade-curves]: SimTrade, "Volatility curves — smiles & smirks" (surface), https://www.simtrade.fr/blog_simtrade/volatility-curves-smiles-smirks/ — finance education. verified-as-of: 2026-06-16
[^spotgamma-crush]: SpotGamma, "IV Crush Explained," https://support.spotgamma.com/hc/en-us/articles/15249330755859-IV-Crush-Explained — options-analytics provider (worked example). verified-as-of: 2026-06-16
[^sofi-crush]: SoFi, "Implied Volatility Crush," https://www.sofi.com/learn/content/implied-volatility-crush/ — brokerage education. verified-as-of: 2026-06-16
[^cboe-vix]: Cboe Insights, "What the VIX and VIX1D Indices Attempt to Measure," https://www.cboe.com/insights/posts/what-the-vix-and-vix-1-d-indices-attempt-to-measure-and-how-they-differ — exchange/authoritative. verified-as-of: 2026-06-16
[^cboe-vix-method]: Cboe, VIX Methodology white paper (model-free / variance-swap weighting), https://cdn.cboe.com/resources/vix/VIX_Methodology.pdf — exchange/primary. verified-as-of: 2026-06-16
[^wiki-vix]: Wikipedia, "VIX" (model-free; "fear gauge"; ±16% reading), https://en.wikipedia.org/wiki/VIX — encyclopedic. verified-as-of: 2026-06-16
[^simtrade-bsm]: SimTrade, "Black-Scholes-Merton option pricing model" (inputs, assumptions, American-call exception), https://www.simtrade.fr/blog_simtrade/black-scholes-merton-option-pricing-model/ — finance education. verified-as-of: 2026-06-16
[^brilliant-bsm]: Brilliant, "Black-Scholes-Merton," https://brilliant.org/wiki/black-scholes-merton/ — math/science education. verified-as-of: 2026-06-16
[^macroption-bsm]: Macroption, "Black-Scholes Assumptions," https://www.macroption.com/black-scholes-assumptions/ — options reference. verified-as-of: 2026-06-16
[^essex-matei]: M. Matei (University of Essex), "Black-Scholes: merits and shortcomings" (constant-vol failure, fat tails, European-only), https://www1.essex.ac.uk/economics/documents/eesj/matei.pdf — academic. verified-as-of: 2026-06-16
[^optiontradingtips-binomial]: OptionTradingTips, "Binomial Model," https://www.optiontradingtips.com/pricing/binomial-model.html — options education. verified-as-of: 2026-06-16
[^crr-binomial]: M. Brenndoerfer, "Binomial Tree Option Pricing (Cox-Ross-Rubinstein)" (convergence to BSM), https://mbrenndoerfer.com/writing/binomial-tree-option-pricing-cox-ross-rubinstein — explainer. verified-as-of: 2026-06-16
[^analystprep-pcp]: AnalystPrep (CFA), "Put-Call Parity," https://analystprep.com/cfa-level-1-exam/derivatives/put-call-parity/ — CFA prep. verified-as-of: 2026-06-16
[^macroption-pcp]: Macroption, "Put-Call Parity Formula," https://www.macroption.com/put-call-parity-formula/ — options reference. verified-as-of: 2026-06-16
[^oic-pcp]: OIC, "Put-Call Parity," https://www.optionseducation.org/advancedconcepts/put-call-parity — authoritative options education (synthetics). verified-as-of: 2026-06-16
[^millersville-american]: R. Buchanan (Millersville University), American options & put-call parity bounds, https://sites.millersville.edu/rbuchanan/book/American.pdf — academic (parity is an inequality for American options). verified-as-of: 2026-06-16
