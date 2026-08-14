# Ecommerce Fundamentals — Context

> Grounding note: every volatile figure below (market shares, fee ranges,
> benchmarks) carries an inline `[^n]` footnote to a real source fetched or
> surfaced via web search in July 2026. Market-share and fee numbers decay —
> re-verify against the cited primary source before quoting to a customer.

## Contents
- Key benchmarks at a glance
- Platform models
- Unit economics deep dive
- Fulfillment models
- Payment processing
- Cart & checkout optimization
- Marketplace vs DTC
- Anti-patterns
- References

## Key benchmarks at a glance

Headline figures used across the sections below, with their footnotes. All are
volatile — treat as of 2025–2026 and re-verify against the cited source.

| Metric | Figure (2025–2026) | Source |
|---|---|---|
| Amazon share of US *ecommerce* sales | ~40% (Statista: 40.5%) | [^5][^6] |
| Amazon share of *total* US retail | ~7% | [^6] |
| Amazon share of US *marketplace* GMV | ~64% (down from ~71% peak, 2022) | [^7] |
| US retail ecommerce, 2025 | ~$1.23T (+5.4% YoY) | [^8] |
| WooCommerce share of ecommerce systems (W3Techs) | ~48.6% (~8.2% of all sites) | [^4] |
| Shopify App Store apps | 16,000+ | [^1] |
| Avg documented cart-abandonment rate (Baymard) | 70.22% | [^24] |
| Top *actionable* abandonment cause | extra costs too high — 39% | [^24] |
| Typical consumer credit-card interchange | ~1.5–3% (CNP > card-present) | [^13][^15] |
| Amazon referral fee, most categories | 8–15% (full range ~6–45%) | [^28][^29] |
| LTV:CAC guardrail (borrowed from SaaS) | ≥3:1 | [^9] |
| Legacy Visa chargeback-monitoring threshold | 0.9% dispute ratio + 100 disputes/mo | [^16] |

## Platform models

The first architectural decision is *who owns the infrastructure and the PCI
burden*. That choice cascades into cost structure, customization ceiling, and
how much engineering the merchant must staff.

**Hosted SaaS platforms** (Shopify, BigCommerce): the vendor manages
infrastructure, uptime, and PCI compliance; the merchant customizes via
themes and apps. Shopify dominates SMB-to-mid-market DTC on the strength of
its app ecosystem — the Shopify App Store now lists **over 16,000 apps**[^1]
(up from a few thousand in 2021) — and its accelerated checkout, Shop Pay.
Shopify's own commissioned study (run by a "big-three" management consultancy)
reports Shop Pay converts **up to 50% better than guest checkout** and
outpaces other accelerated checkouts by at least 10%, with its mere presence
on a storefront driving a ~5% conversion lift.[^2][^3] Treat the 50% figure as
vendor-sponsored and directional, not an independent benchmark.

**Self-hosted / open-source** (WooCommerce on WordPress, Magento Open Source):
full code control, no platform transaction fee, but the merchant owns hosting,
security patching, and the full PCI scope. WooCommerce is the single most
widely deployed ecommerce technology by site count: W3Techs reports it runs on
**~8.2% of all websites** and holds **~48.6% of all ecommerce systems** it
surveys (as of July 2026).[^4] BuiltWith-derived trackers put its share of the
top-1M sites lower (single digits), so "most-installed ecommerce software
globally" is defensible on total live-site count but is *not* the same as
"most-used among high-traffic sites."[^4]

**Enterprise / composable commerce** (Adobe Commerce/Magento Enterprise,
BigCommerce Enterprise, commercetools, Salesforce Commerce Cloud): built for
multi-brand/multi-region catalogs, B2B pricing tiers, and headless/composable
architecture (MACH — Microservices, API-first, Cloud-native, Headless). The
tradeoff is implementation cost and integration engineering.

**Marketplaces** (Amazon, Etsy, Walmart Marketplace, eBay): the seller lists
into an existing demand pool instead of building one. Keep three denominators
distinct so you do not cross-wire the numbers:
- Amazon is estimated at **just over 40% of US *ecommerce* sales** (Statista put
  it at 40.5% in 2025; EMARKETER projected it to surpass 40%).[^5][^6]
- That is only **~7% of *total US retail*** (online + offline).[^6]
- Among US *marketplace* GMV (gross merchandise value) specifically, Amazon's
  share is higher still — EMARKETER data put it around 64% in 2025, down from a
  ~71% peak in 2022.[^7]

US retail ecommerce totaled roughly **$1.23 trillion in 2025** (+5.4% YoY),
per Digital Commerce 360.[^8] At that scale, a marketplace presence is close to
mandatory in many categories even when DTC is the primary channel.

## Unit economics deep dive

- **CAC (Customer Acquisition Cost):** total sales + marketing spend ÷ new
  customers in the period. *Blended CAC* (all channels, including organic) and
  *paid CAC* (paid channels only) tell very different stories — always specify
  which.
- **AOV (Average Order Value):** revenue ÷ orders. Levers: bundling,
  free-shipping thresholds set just above current AOV, post-purchase upsells,
  and volume discounts.
- **Contribution margin:** revenue − COGS − variable costs (payment
  processing, pick/pack/ship, returns, and ad spend attributable to the sale).
  This is the number that must fund CAC recovery. Gross margin alone overstates
  profitability because it ignores fulfillment and payment costs.
- **LTV (Lifetime Value):** cumulative contribution margin per customer over
  the relationship. For repeat-purchase brands, a common approximation is
  `LTV = AOV × purchase frequency × average lifespan × margin`.
- **LTV:CAC ratio:** the standard efficiency benchmark. The widely repeated
  **≥3:1** target originates in **SaaS**, not ecommerce — David Skok's
  "SaaS Metrics 2.0" argues the best SaaS businesses run higher than 3, "sometimes
  as high as 7 or 8," and SaaS Capital benchmark work puts the *healthy median*
  near 3:1.[^9][^10] DTC operators borrow the heuristic: <1:1 means the business
  loses money per customer acquired; a very high ratio (>5:1) can signal
  *under*-investment in growth. Because it is a borrowed rule of thumb, treat it
  as a directional guardrail, not a physical constant.
- **CAC payback period:** months of contribution margin needed to recover CAC.
  Shorter payback reduces the cash-flow risk of paid acquisition.
- **CAC creep:** as paid channels (Meta, Google Shopping) saturate, CAC tends to
  rise. Apple's App Tracking Transparency (ATT, iOS 14.5, 2021) degraded ad
  targeting and measurement; reported **iOS user-acquisition costs rose ~20–30%**
  afterward, and academic analysis (Aridor et al., 2024) documents measurable
  medium-term effects on advertiser performance.[^11][^12] Note this evidence is
  strongest for *app-install / mobile-measurement* cost — read across to
  ecommerce paid-social CAC as directional, not a clean measured figure. The
  strategic response is to shift weight toward owned channels (email/SMS,
  loyalty, organic/SEO) to protect LTV:CAC.

## Fulfillment models

| Model | Description | Best for | Tradeoff |
|---|---|---|---|
| In-house / DTC warehousing | Merchant owns pick-pack-ship | Full control, brand-quality unboxing | Capital-intensive; doesn't scale elastically |
| 3PL (third-party logistics) | Outsourced warehousing + shipping (ShipBob, ShipMonk, etc.) | Scaling brands avoiding fixed warehouse costs | Less control over pack experience; per-unit fees |
| FBA (Fulfillment by Amazon) | Amazon warehouses and ships for Prime-eligible listings | Amazon-marketplace sellers wanting the Prime badge | Amazon owns the customer relationship + packaging; layered fees |
| Dropshipping | Supplier ships directly to the buyer; merchant holds no inventory | Low-capital testing of SKUs/niches | Longer shipping, thinner margins, less QC (see `dropshipping-and-fulfillment`) |

Hybrid models are common: a brand might run 3PL for core SKUs and dropship the
long-tail/test SKUs.

**FBA fee mechanics.** FBA is not a flat percentage cut. Amazon charges a
*per-unit fulfillment fee* (roughly $3–$10+ for standard-size items, scaling with
size and weight) plus *monthly storage fees* — on top of the category referral
fee.[^29] For a low-priced SKU the fixed per-unit fee can dominate the margin, so
model FBA per-SKU rather than as a blanket rate. This is why the marketplace
"take" compounds (see Anti-patterns): referral + fulfillment + storage + optional
ad fees stack on the same sale.

## Payment processing

- **Gateway vs processor vs acquirer:** the *gateway* (Stripe, Braintree,
  Authorize.net, Shopify Payments) captures and encrypts card data; the
  *processor* moves the transaction through the card networks; the *acquiring
  bank* settles funds to the merchant. Flat-rate providers (Stripe, Shopify
  Payments) bundle all three plus their markup into a single published rate.

- **Interchange fees:** the largest component of card-processing cost, set by
  the card networks (Visa/Mastercard) and paid to the *card-issuing* bank. US
  rates vary by card type and how the card is presented. Representative
  published Visa rates: CPS/Retail (card-present) ~**1.51% + $0.10**; rewards
  cards run higher, e.g. Signature Preferred ~**2.5% + $0.10**; regulated debit
  is capped low (about **$0.21 + 0.05%**) under the Durbin Amendment.[^13][^14][^15] The
  practical takeaway holds: interchange for typical consumer credit cards sits
  in the **~1.5–3%** band, with **card-not-present (ecommerce) transactions
  priced above card-present** for the same card because they carry more fraud
  risk.[^13][^15] Interchange is only one layer of what the merchant pays. The
  all-in **merchant discount rate** = interchange (to the issuing bank) + network
  assessments (to Visa/Mastercard) + processor markup. Interchange is the largest
  and least negotiable layer; only the processor markup is genuinely negotiable —
  which is why flat-rate providers can profitably bundle all three into one
  published rate for small merchants.

- **Chargebacks:** a cardholder-initiated reversal through their issuing bank
  rather than a merchant refund. CNP (card-not-present) ecommerce carries
  structurally higher chargeback exposure than in-person retail. Visa's legacy
  Dispute Monitoring Program (VDMP) flagged merchants at roughly a **0.9%
  dispute-to-transaction ratio *and* 100+ monthly disputes**; exceeding it meant
  fees and remediation timelines, and sustained excess risked processor
  termination.[^16] As of 2025 Visa is consolidating VDMP/VFMP into the **Visa
  Acquirer Monitoring Program (VAMP)**, with updated thresholds effective mid-2025
  — verify current numbers against Visa's VAMP materials, as this is actively
  changing.[^17]

- **PCI DSS:** the Payment Card Industry Data Security Standard governs how
  merchants handle cardholder data. Fully outsourcing payment capture to a
  compliant provider — hosted checkout, redirect, or iframe (Shopify Payments,
  Stripe Checkout) so **card data never touches the merchant's systems** —
  qualifies the merchant for **SAQ A**, "the smallest possible subset of PCI DSS
  requirements."[^18][^19] Important dated nuance: **PCI DSS v4.0.1 tightened
  SAQ A eligibility and added requirements** (e.g., protecting the payment page
  from script tampering / e-skimming), so SAQ A still sharply reduces scope but
  is no longer as light as it once was — the "TPSP-embedded payment page"
  criteria and new script/iframe controls now apply.[^20][^21] SAQ A remains a
  short questionnaire (roughly two dozen controls) versus the full
  self-hosted-card-field scope.[^21]

- **BNPL (Buy Now, Pay Later):** Klarna, Afterpay, and Affirm integrate as
  alternative payment methods. Providers *market* large AOV lifts — Klarna's
  merchant base has reported a **~45% AOV increase**, and the CFPB's 2022 market
  report catalogs vendor claims like "41% increase in average order value" and
  "30% increase in conversion."[^22][^23] Read these as **provider-reported
  marketing claims** (the CFPB frames them exactly that way), not independently
  audited results. BNPL shifts default/fraud risk to the provider in exchange
  for a merchant fee that is typically higher than card interchange.

## Cart & checkout optimization

**Baseline abandonment.** Baymard Institute's running meta-analysis of 50
studies puts the **average documented cart-abandonment rate at 70.22%**
(list last updated Sept 2025).[^24] A large share of that is unavoidable
"window-shopping" behavior, so the actionable question is *why* users who
intended to buy dropped out.

**Reasons for abandonment** — from Baymard's 2025 survey of 1,026 US online
shoppers (multiple selections allowed). Note the two framings:[^24]
- The single most-selected reason is **"just browsing / not ready to buy" (43%)**
  — Baymard flags this as largely *unavoidable*, so don't optimize against it.
- Excluding browsers, the top **actionable** cause is **extra costs too high
  (shipping/tax/fees) — 39%**, followed by:
  - Delivery too slow — 21%
  - Didn't trust the site with card info — 19%
  - Forced to create an account — 19%
  - Checkout too long / complicated — 18%
  - Return policy unsatisfactory — 15%
  - Site errors / crashes — 15%
  - Couldn't see total cost up front — 14%
  - Not enough payment methods — 10%
  - Card declined — 8%

**The checkout-length gap.** Baymard's usability testing finds an ideal
checkout can be as short as **12–14 form elements (7–8 fields)**, yet the
average US checkout shows **~23.5 form elements** by default — leaving a
documented **20–60% reduction** available on most flows. Baymard estimates a
well-designed checkout can lift conversion ~35% on large sites, and puts the
orders recoverable across the US + EU from checkout-usability issues alone at
roughly **$260 billion**.[^24]

**Levers:** guest checkout as the default; upfront/transparent shipping and tax
estimates; one-page or single-scroll checkout; express-wallet buttons (Shop Pay,
Apple Pay, Google Pay, PayPal); progress indicators; autofill/address
validation; and mobile-first form design. Mobile now drives **the majority of
ecommerce traffic (60%+ of website visits)**, yet **desktop typically converts
at a higher rate (~1.7×)** — so mobile UX is a traffic-capture problem, not a
reason to deprioritize desktop conversion.[^25][^26]

**Cart recovery:** abandoned-cart email/SMS sequences (see
`ecommerce-automation-and-tools`) recover a meaningful minority of carts when
timed within the first hours. Aggregator benchmarks cite abandoned-cart email
conversion in the **high-single-digit-to-low-teens percent** range (one
aggregator reports ~10.7%); treat these as vendor/aggregator estimates, not
audited figures, and measure your own baseline.[^27]

## Marketplace vs DTC

| Dimension | Marketplace (Amazon/Etsy) | DTC (owned site) |
|---|---|---|
| Customer-data ownership | Marketplace owns the relationship; seller gets limited buyer data | Merchant owns full first-party data (more valuable post-ATT) |
| Traffic | Built-in demand, but pay-to-play (PPC ads + referral fees) | Must build demand (paid + organic + owned) |
| Fees | Referral fee **8–15% for most categories** (6–45% across the full range; e.g., Amazon Device Accessories 45%), plus fulfillment/FBA and optional ad fees[^28][^29] | Payment processing (~1.5–3% interchange + markup) + full CAC |
| Margin | Layered fees compress margin | Higher gross-margin potential, but bears full CAC |
| Brand control | Limited — templated listings, no custom checkout | Full control over storefront, upsell, loyalty |
| Trust / conversion | Inherits marketplace trust; often converts browsers faster | Must earn trust independently (reviews, badges, guarantees) |

Amazon's referral fee is **8% or 15% per category for most goods**, but the
published schedule spans **6% to 45%** depending on category (jewelry, for
example, is tiered ~20% then 5%; Amazon Device Accessories is 45%).[^28][^29]
Most maturing DTC brands run **both** channels deliberately: marketplace for
reach and liquidity, owned site for margin, data, and brand equity — treating
the marketplace as a top-of-funnel/awareness channel that feeds retargeting to
the owned site.

## Anti-patterns

- **Optimizing AOV/CAC in isolation without contribution margin.** A
  "successful" campaign can still lose money per unit once fulfillment, returns,
  and payment costs are subtracted.
- **Treating marketplace sales as free distribution.** Referral fees
  (8–15% typical[^28]) plus per-unit FBA fulfillment fees, storage, and ad spend
  stack. *Illustrative estimate:* for many sellers the combined marketplace take
  lands in the ~25–40%-of-sale-price range — but this is a synthesis of the
  component fees above, **not** a single published figure; model it from your
  own category's actual referral + FBA + storage + ad rates.
- **Ignoring interchange/chargeback costs when modeling margin**, especially for
  high-ticket or high-fraud-risk categories where CNP interchange and dispute
  exposure are highest.[^13][^16]
- **Adding account-creation friction to checkout without measuring its
  abandonment cost** — forced account creation is a top-5 actionable
  abandonment cause (19% of shoppers).[^24]
- **Quoting the ≥3:1 LTV:CAC target as if it were ecommerce-native.** It is a
  borrowed SaaS heuristic;[^9] validate against your own payback and margin math.

## References

[^1]: Shopify App Store — homepage states "over 16,000 apps." https://apps.shopify.com/ (via web search, July 2026).
[^2]: Shopify, "Shop Pay: The fastest accelerated checkout on the internet" — "as much as 50% better conversion compared to guest checkout," "~5% lift" from presence; figures from a Shopify-commissioned study. https://www.shopify.com/shop-pay (via web search, July 2026).
[^3]: Shopify (Enterprise blog), "Shopify Checkout is the best-converting in the world. Here's why." https://www.shopify.com/enterprise/blog/shopify-checkout (via web search, July 2026).
[^4]: W3Techs, "Usage statistics of WooCommerce" — ~8.2% of all websites; ~48.6% of ecommerce systems surveyed (as of July 3, 2026). https://w3techs.com/technologies/details/cm-woocommerce (fetched July 2026). Cross-checked against BuiltWith-derived trackers (single-digit % of top-1M sites), which is why "most-installed by total site count" ≠ "most-used among high-traffic sites."
[^5]: Statista, "Market share of leading retail e-commerce companies in the United States 2025" — Amazon 40.5%, Walmart 9.2%. https://www.statista.com/statistics/274255/market-share-of-the-leading-retailers-in-us-e-commerce/ (via web search, July 2026).
[^6]: EMARKETER, "Amazon will surpass 40% of US ecommerce sales this year" — Amazon >40% of US ecommerce; its share of *total* US retail ~7.0% in 2025. https://www.emarketer.com/content/amazon-will-surpass-40-of-us-ecommerce-sales-this-year (via web search, July 2026).
[^7]: EMARKETER data reported via LinkedIn (Sky Canaves) — Amazon's share of US *marketplace* ecommerce ~64% in 2025, down from a ~71% peak in 2022. https://www.linkedin.com/posts/sky-canaves_marketplaces-will-account-for-42-of-online-activity-7377412297612328960-dSf5 (via web search, July 2026).
[^8]: Digital Commerce 360, "2025 U.S. ecommerce sales mark fourth straight year of single-digit growth" — ~$1.234T in 2025, +5.4% YoY. https://www.digitalcommerce360.com/article/us-ecommerce-sales/ (via web search, July 2026).
[^9]: David Skok, "SaaS Metrics 2.0 — A Guide to Measuring and Improving What Matters" (forEntrepreneurs) — best SaaS LTV:CAC >3, "sometimes as high as 7 or 8." https://www.forentrepreneurs.com/saas-metrics-2/ (via web search, July 2026).
[^10]: Monetizely, "The LTV/CAC Ratio" citing SaaS Capital benchmark — median healthy SaaS LTV/CAC ~3:1. https://www.getmonetizely.com/articles/the-ltvcac-ratio-your-north-star-metric-for-saas-success (via web search, July 2026).
[^11]: AdAction, "Mobile App User Acquisition Cost: Benchmarks, Formula & Strategy" — iOS acquisition costs up ~20–30% since ATT/iOS 14.5. https://www.adaction.com/blog/mobile-app-user-acquisition-cost (via web search, July 2026).
[^12]: Aridor, Che, Hollenbeck, Kaiser, McCarthy et al. (2024), "Evidence from Apple's App Tracking Transparency" — quantifies medium-term economic effects of ATT on advertiser performance (Toulouse School of Economics working paper). https://www.tse-fr.eu/sites/default/files/TSE/documents/sem2024/eco_platforms/aridor2024.pdf (via web search, July 2026).
[^13]: Visa, "Visa USA Interchange Reimbursement Fees" (published rate schedule PDF). https://usa.visa.com/content/dam/VCOM/download/merchants/visa-usa-interchange-reimbursement-fees.pdf (via web search, July 2026).
[^14]: Mastercard, "Merchant Interchange Rates" (US). https://www.mastercard.com/us/en/business/support/merchant-interchange-rates.html (via web search, July 2026).
[^15]: Clearly Payments, "Interchange Rates in USA for Credit Card Processing" — e.g., Visa CPS/Retail 1.51% + $0.10, rewards 1.65–2.5% + $0.10. https://www.clearlypayments.com/interchange-rates-in-usa/ (via web search, July 2026).
[^16]: Checkout.com, "What is the Visa Dispute Monitoring Program? (VDMP)" — early-warning threshold ~0.9% dispute ratio and 100 disputes/month. https://www.checkout.com/blog/what-is-the-visa-dispute-monitoring-program (via web search, July 2026).
[^17]: Visa, "Visa Acquirer Monitoring Program (VAMP) Fact Sheet 2025" — VAMP consolidates prior fraud/dispute monitoring programs; threshold updates effective June 2025. https://corporate.visa.com/content/dam/VCOM/corporate/visa-perspectives/security-and-trust/documents/visa-acquirer-monitoring-program-fact-sheet-2025.pdf (via web search, July 2026).
[^18]: PCI Security Standards Council, "Best Practices for Securing E-commerce" (Information Supplement, PDF) — SAQ A is "the smallest possible subset of PCI DSS requirements" for merchants who fully outsource to a PSP. https://www.pcisecuritystandards.org/pdfs/best_practices_securing_ecommerce.pdf (via web search, July 2026).
[^19]: Lexology / PCI guidance summary — SAQ A eligibility for merchants who fully outsource website payment operations (URL redirect or iframe). https://www.lexology.com/library/detail.aspx?g=54973ba8-a0f8-435c-b4a8-2e9c8e79a975 (via web search, July 2026).
[^20]: PCI Security Standards Council blog, "FAQ Clarifies New SAQ A Eligibility Criteria for E-Commerce Merchants" — v4.0.1 eligibility now keyed to a TPSP/processor-embedded payment page. https://blog.pcisecuritystandards.org/faq-clarifies-new-saq-a-eligibility-criteria-for-e-commerce-merchants (via web search, July 2026).
[^21]: TrustedSec, "The Hidden Trap in the PCI DSS SAQ A Changes" — SAQ A is ~21 controls (moving to ~19), but v4.0.1 adds payment-page script/anti-skimming requirements for SAQ A merchants. https://trustedsec.com/blog/the-hidden-trap-in-the-pci-dss-saq-a-changes (via web search, July 2026).
[^22]: CNBC, "Why retailers are embracing buy now, pay later financing services" — Klarna's merchant base reports a ~45% increase in AOV. https://www.cnbc.com/2021/09/25/why-retailers-are-embracing-buy-now-pay-later-financing-services.html (via web search, July 2026).
[^23]: US Consumer Financial Protection Bureau, "Buy Now, Pay Later: Market Trends and Consumer Impacts" (Sept 2022, PDF) — catalogs *provider-reported* claims of "41% increase in average order value," "30% increase in conversion," "85% higher merchant AOV." https://files.consumerfinance.gov/f/documents/cfpb_buy-now-pay-later-market-trends-consumer-impacts_report_2022-09.pdf (via web search, July 2026).
[^24]: Baymard Institute, "50 Cart Abandonment Rate Statistics" — 70.22% average documented abandonment rate (avg of 50 studies; list last updated Sept 22, 2025) and 2025 reasons-for-abandonment survey of 1,026 US online shoppers; also checkout-length benchmarks (ideal 12–14 form elements vs ~23.5 average) and ~35% conversion-lift potential. https://baymard.com/lists/cart-abandonment-rate (fetched July 2026).
[^25]: OuterBox, "Mobile eCommerce Statistics" — mobile now makes up 60%+ of website traffic. https://www.outerboxdesign.com/articles/digital-marketing/mobile-ecommerce-statistics/ (via web search, July 2026).
[^26]: Smart Insights, "E-commerce conversion rate benchmarks — 2025 update" — desktop conversion runs ~1.7× smartphone. https://www.smartinsights.com/ecommerce/ecommerce-analytics/ecommerce-conversion-rates/ (via web search, July 2026).
[^27]: Mailmend, "Cart Abandonment Recovery Statistics" (aggregator) — abandoned-cart emails ~44.76% open rate, ~10.7% conversion rate. https://mailmend.io/blogs/cart-abandonment-recovery-statistics (via web search, July 2026). Treat as aggregator estimate.
[^28]: Amazon Seller Central, "Selling on Amazon fee schedule" — referral fee rates by category (e.g., 8% at/under a price threshold, 15% above; category tiers vary). https://sellercentral.amazon.com/help/hub/reference/external/G200336920 (via web search, July 2026).
[^29]: Amazon, "Standard selling fees" (sell.amazon.com/pricing) — referral fees range across categories (most 8–15%; full range ~6–45%; e.g., Amazon Device Accessories 45%), plus FBA fulfillment and storage fees. https://sell.amazon.com/pricing (via web search, July 2026).
