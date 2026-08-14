# Ecommerce Automation and Tools — Context

This reference covers the tooling layer that automates ecommerce operations: routing
orders to the right fulfillment source, keeping stock in sync across channels,
repricing against marketplace competitors, running email/SMS lifecycle flows, gluing
apps together with no-code automation, and the Shopify app ecosystem plus headless
architecture that many of these tools plug into. It is a decision aid for tool
*selection and sequencing*, not a setup manual for any single product.

All named-vendor mechanics, pricing models, and benchmark figures below are grounded in
sources fetched or searched in **July 2026** (see `## References`). Vendor pricing and
marketplace statistics decay fast — treat every specific number as "as of the cited
source's date" and re-verify before quoting it to a customer.

## Contents
- Order routing
- Inventory sync
- Automated repricing
- Email/SMS lifecycle automation
- General workflow automation (Zapier/Make)
- Shopify app ecosystem and headless commerce
- Anti-patterns
- References

## Order routing

Order-routing tools decide, for each incoming order, which fulfillment source should
ship it — a specific warehouse, a 3PL, a dropship supplier, or Amazon FBA/Multi-Channel Fulfillment (MCF) — based on
rules like which locations stock the item, proximity to the recipient, shipping cost, or
delivery-speed promise. The category spans from lightweight label-and-rules tools up to
full distributed order-management systems (OMS/DOM).

**ShipStation** sits at the lighter end: it is primarily multi-carrier label automation
plus a rules engine. Automation Rules apply on import (when an order enters *Awaiting
Shipment*) and can auto-select carrier/service, assign package types, apply tags, and
split multi-item orders across shipments.[^4] Its **Auto-Routing** feature assigns an
order to a warehouse based on which locations stock the items and which is closest to the
recipient (documented as US/Canada-only), and it supports mixed fulfillment methods —
self-fulfill, FBA, 3PL/dropship — through the same rules layer.[^4] ShipStation is a good
fit when the core need is *shipping* automation with modest routing logic, not deep
inventory or purchasing control.

**Cin7** and **Extensiv Order Management** target higher-volume, multi-location,
multi-channel operations where routing is one function of a broader inventory/OMS
platform. Extensiv Order Management explicitly advertises rules-engine "algorithmic
fulfillment" with dropship/3PL routing, cross-docking, and Amazon multi-channel
fulfillment, plus auto-generated purchase orders driven by sales velocity, lead time, and
seasonality.[^3] Cin7 splits into two products: **Cin7 Core** (a ready-to-go inventory/OMS
for brands that don't need heavy customization) and **Cin7 Omni** (a highly configurable
omnichannel platform with native EDI and 3PL integrations and multi-entity support for
more complex operations).[^5] Both connect to hundreds of sales and accounting platforms
(Shopify, Amazon, WooCommerce, BigCommerce, QuickBooks, Xero).[^5]

> **Naming correction (important):** Earlier drafts referred to "Extensiv/Skubana" as if
> Skubana were still an independent product. It is not. Skubana was acquired by 3PL Central
> in 2021; 3PL Central rebranded to **Extensiv** in 2022, and the former Skubana product is
> now sold as **Extensiv Order Management** (also seen as "Extensiv Order Manager").[^1][^2]
> Skubana combined with 3PL Central, Scout, and CartRover under the Extensiv brand.[^1] Use
> "Extensiv Order Management (formerly Skubana)" on first mention and "Extensiv" thereafter;
> "Skubana" alone is a legacy name and will read as out-of-date to any current operator.

**Dropship-specific connectors** route an order directly to a named supplier per SKU
rather than to an internal warehouse. **DSers** is the canonical example on Shopify: it is
the official AliExpress dropshipping solution and an Alibaba Group partner (it succeeded
Oberlo after Oberlo shut down in 2022), and it automates the browse→import→auto-order flow
from AliExpress/1688/Alibaba suppliers to a Shopify (or WooCommerce/Wix) store.[^7] For
supplier sourcing, print-on-demand, and dropship margin mechanics, see
`dropshipping-and-fulfillment`; this skill covers only the routing/automation surface.

**Selection heuristic:**
- Shipping-label automation with light routing, single or few warehouses → ShipStation.[^4]
- Multi-warehouse / multi-channel / higher volume with purchasing and inventory control →
  Cin7 (Core for simpler, Omni for complex) or Extensiv Order Management.[^3][^5]
- Pure dropship, supplier-per-SKU → DSers or another dropship connector.[^7]

## Inventory sync

Selling the same SKUs across Shopify, Amazon, eBay, Walmart, TikTok Shop, and other
channels simultaneously creates **overselling risk** unless stock levels sync across every
channel in near-real-time after each sale. The failure mode this category exists to
prevent is concrete: selling the last unit on two channels at once because neither channel
knew about the other's sale, forcing a cancellation, a refund, and a marketplace-metrics
hit (Amazon in particular penalizes cancellation/defect rates).

Dedicated inventory-management SaaS centralizes the authoritative stock count and pushes
updates outward to every connected channel:

- **Cin7** tracks inventory, orders, stock locations, and sales channels in real time
  across online, retail, wholesale, and POS, and connects to 700+ platforms; Omni adds
  native EDI/3PL and multi-entity support for brands managing stock across many
  channels.[^5]
- **Extensiv Order Management** provides multi-channel inventory visibility and control
  across all warehouses, 3PLs, dropshippers, and FBA distribution centers, treating them as
  one pooled network.[^3]
- **Linnworks** connects inventory directly to orders, warehouses, listings, and shipping
  so that a sale on any channel instantly updates stock everywhere and triggers the right
  fulfillment workflow; it syncs listings to actual stock levels and offers two-way
  connections to Amazon, eBay, Shopify, Walmart, TikTok Shop, and dozens more channels.[^6]
  Its rules-based automation also routes orders by channel, stock location, carrier, and
  packaging — so in practice Linnworks spans both this section and Order routing.[^6]

Lighter-weight Shopify-native apps handle simpler two-or-three-channel setups where a full
OMS would be overkill. The decision axis is channel count and volume: a couple of channels
and modest order volume can live on a Shopify-native sync app; five-plus channels,
multiple warehouses/3PLs, or marketplace-metric sensitivity push you toward a dedicated OMS
(Cin7 / Extensiv / Linnworks).[^3][^5][^6]

The non-negotiable design property is **a single source of truth** for the available-to-sell
count, with every channel treated as a subscriber to it rather than an independent ledger.
Two "syncing" systems that both think they own the count will drift and oversell.

## Automated repricing

Repricing splits into two distinct needs that are often conflated:

### 1. Marketplace competitive repricing (Amazon / Walmart)

Tools watch competitor pricing and Buy Box status and adjust the seller's price
algorithmically within a merchant-set **floor (min) and ceiling (max)** to win or retain
the Buy Box. Two philosophies exist, and the difference is decision-relevant:

- **Rule-based repricers** (e.g., **RepricerExpress**) react to the market by explicit
  rules. The seller sets a min (floor) and max (ceiling) per SKU; the tool never uploads a
  price outside that band, and rules dictate how to compete (e.g., match, beat by a set
  amount, or compete with the next-best seller).[^9] RepricerExpress lets you define
  behavior for edge cases — what to do when a competitor prices *below* your min (Do Not
  Reprice, Go to Min, Go to Max, or Compete with Next Best Seller) and what to do *while you
  already hold the Buy Box* (by default it will not change your price when you have the Buy
  Box, to avoid risking it; advanced settings allow "increase price" or "reprice as
  normal").[^9] The floor/ceiling band is the primary margin guardrail.

- **AI / game-theory repricers** (e.g., **Seller Snap**) model competitors as players who
  react to your moves and aim to win the Buy Box *at the highest sustainable price* rather
  than the lowest. Seller Snap markets a "cooperative" strategy — described as a "wait my
  turn" approach where, instead of undercutting the current Buy Box holder, the repricer
  holds back, lets that seller raise prices, and takes its turn at the higher price — to
  avoid triggering a race-to-the-bottom price war.[^8] The pitch is explicitly *against*
  naive rule-based undercutting.[^8]

The practical takeaway: rule-based tools are transparent and predictable but can get
dragged into price wars if rules are set aggressively; AI/game-theory tools aim to dampen
price wars but are more of a black box. **Either way, a floor price is mandatory** — see
Anti-patterns.

### 2. Supplier-cost-driven repricing (dropshipping-specific)

When a dropship supplier's cost changes, the storefront price must adjust to protect
margin. Some dropship-order-routing apps include cost-change repricing as a feature;
otherwise it requires custom rule automation — e.g., a Zapier/Make scenario watching a
supplier cost feed and updating Shopify prices, or a custom script. This is a different
problem from Buy Box repricing (there is no competitor to react to — you are protecting a
target margin against input-cost drift) and is covered further in
`dropshipping-and-fulfillment`.

## Email/SMS lifecycle automation

**Klaviyo** (email-first, with native and integrated SMS) and **Postscript** and
**Attentive** (SMS-first) are the dominant ecommerce lifecycle-marketing platforms.
Postscript is Shopify-focused SMS/MMS; Attentive is a broader conversational platform
offering two-way, real-time personalized SMS *and* email across retail/ecommerce.[^14]
Lifecycle "flows" are automated, trigger-based sequences (as opposed to one-off
"campaigns"). Klaviyo's own data shows flows punch far above their volume: they generate a
large share of email revenue from a small share of sends because they fire at
high-intent moments.[^12]

### Standard automated flows
- **Abandoned cart / abandoned checkout:** triggered when a cart or checkout is started but
  not completed within a set window; a short multi-touch sequence.
- **Welcome series:** triggered on first signup/purchase to introduce the brand and often
  deliver a first-purchase incentive.
- **Post-purchase:** order confirmation, shipping updates, then a delayed review request or
  replenishment reminder.
- **Win-back / lapsed:** triggered after a period of purchase inactivity, often paired with
  an incentive or new-product highlight.
- **Browse abandonment:** triggered on product-page views without an add-to-cart (lower
  intent than cart abandonment, so typically lighter touch).
- **Subscription-specific flows** (payment-failure/dunning notices) live in
  `subscription-commerce` for the churn/dunning *strategy*; this skill covers the general
  email/SMS *tooling* used to send them.

### Abandoned-cart timing (corrected)

Earlier drafts asserted a generic "1 hour, 24 hours, 3 days" email sequence. That is not
what the platforms' own guidance says, and email vs SMS timing genuinely differ — a
decision-useful nuance:

- **Email (Klaviyo):** Klaviyo's pre-built abandoned-cart flow uses a **4-hour** time delay
  before the first message. Its guidance recommends sending the first message roughly
  **2–4 hours** after checkout starts, a second message **20–48 hours** later, and finds
  that **2–3 messages total** yields optimal performance; longer delays suit
  high-consideration/high-ticket items.[^10]
- **SMS:** SMS best practice is a much tighter first touch than email — the first text
  within roughly the **first 15–30 minutes** while intent is hot, personalized with the
  shopper's name, product, and a direct checkout link, then one next-day follow-up
  (sometimes with a small incentive), then stop. (This is widely-cited SMS timing guidance;
  the Postscript figures below are the benchmark data, not this timing rule.)[^13]

So the corrected rule of thumb is: **email leads with hours-to-a-day cadence over 2–3
touches; SMS opens within the first half hour and stays short.** Do not copy an email
cadence onto an SMS flow.

### Why abandoned cart is the flagship flow — with real benchmarks

Abandoned-cart is consistently the highest-return automated flow because it catches
demonstrated purchase intent. Published, on-page benchmark figures (verify against the
live report before quoting — these decay):

- **Klaviyo (email):** across measured stores, the abandoned-cart flow shows an **average
  revenue per recipient (RPR) of ~$3.65** and an **average placed-order (conversion) rate
  of ~3.33%**, the highest of any flow type; top-10% performers reach **~$28.89 RPR**.[^11]
  Across all flow types the average RPR is much lower, which is why abandoned-cart and
  welcome are the priority flows to build first.[^11][^12]
- **Postscript (SMS):** across 17,000+ Shopify stores in its 2026 report (2025 data),
  abandoned-cart texts earn roughly **$3.52–$10.95 earnings per message (EPM)** at the
  25th–75th percentile, with a **conversion rate of ~3.97%–7.84%** and click-through of
  ~9.5%–17.3%; the overall median revenue-per-message across all message types is about
  **$0.98**, and median subscriber lifetime value about **$71.20**.[^13]

(Note: some third-party write-ups circulate a "34x SMS ROI" / "$8.11 EPM" figure attributed
to Postscript — those specific numbers are **not** on Postscript's own benchmarks page and
should not be cited as Postscript's. Use the on-page ranges above.[^13])

The design implication: because these flows are so high-return, the failure mode is almost
never "not worth doing" — it is over-messaging (see Anti-patterns).

## General workflow automation (Zapier/Make)

**Zapier** and **Make** (formerly Integromat) provide no-code "trigger → action"
automation that connects apps lacking a native integration — e.g., "new Shopify order →
post to Slack → append a row to Google Sheets → create a task in a project tool." Both are
excellent glue for lower-volume or non-time-critical logic. The two price on **different
units**, and the unit shapes when each stops being economical:

- **Zapier is task-metered.** Every action step that runs counts as a **task** (the
  trigger typically does not; multi-step Zaps consume one task per action step per run).
  The free tier allows on the order of **~100 tasks/month**, and paid plans scale up into
  the hundreds-of-thousands to millions of tasks/month.[^15][^21] When you exceed a plan's
  task limit, Zapier switches to **pay-per-task** billing up to a hard cap of **3× your
  plan's task limit** (base allocation + up to 2× more in overage) before workflows
  stop.[^15] Lower tiers also poll less frequently, adding latency to trigger detection.

- **Make is credit/operation-metered.** Each module action in a scenario consumes a
  **credit** (formerly "operation") — reading, searching, creating, updating, transforming,
  or iterating each cost one, so a 10-step scenario burns ~10 credits per run, and a complex
  run can consume "anywhere from two credits to thousands."[^16] The **free** plan allows
  **1,000 credits/month** with a 15-minute minimum interval between scheduled runs; the
  **Core** plan is **$12/month for 10,000 credits** with scheduling down to 1-minute
  intervals and unlimited active scenarios.[^16] (Router and error-handler modules consume
  zero credits.[^16]) Make's per-step granularity often makes complex multi-step scenarios
  cheaper than the equivalent multi-task Zap, which is part of why heavier automation users
  gravitate to it.

**The volume tradeoff (the load-bearing point):** because *every step* consumes an
allowance, a high-volume, latency-sensitive flow — real-time multi-channel inventory sync
being the canonical example — will burn through task/credit allowances quickly *and* hit
polling-interval latency, and typically outgrows no-code tooling. At scale, move that flow
to a purpose-built SaaS (an OMS/inventory platform) or custom middleware, and keep Zapier/
Make for the low-frequency, non-critical glue they excel at.[^15][^16] Both platforms have
also been layering AI/agentic features on top of the same metering model, but the
underlying per-step billing is what governs the economics.

## Shopify app ecosystem and headless commerce

### The app ecosystem

Nearly every function in this skill — routing, inventory sync, repricing, email/SMS,
reviews, upsells, dropship connectors — is available as an installable Shopify app rather
than custom code. The Shopify App Store is one of the largest app marketplaces of any
hosted commerce platform: third-party App Store trackers counted roughly **11,900 apps as
of late 2024** (the count fluctuates because Shopify periodically purges low-quality
listings under stricter review), with about **87% of merchants using apps** and the
average store installing around **six**.[^20] (These are third-party-tracker figures, not
a Shopify-official headcount, and the tracked total moves over time — cite with the as-of
date.) This app depth is why Shopify remains the default recommendation in
`ecommerce-fundamentals` for merchants who want fast time-to-market without heavy
engineering: you assemble capabilities from apps instead of building them.

### Headless / composable commerce

**Headless commerce** decouples the customer-facing frontend from the commerce backend via
APIs, so the storefront UI is built and deployed independently of the platform that owns
carts, catalog, checkout, and orders.

On Shopify, the first-party headless stack is:
- **Storefront API** (and the Customer Account API) — the GraphQL data layer a custom
  frontend queries for products, carts, and checkout.[^17]
- **Hydrogen** — Shopify's React-based headless framework, built on the open-source **React
  Router** framework, shipping components/hooks/utilities optimized for the Storefront
  API.[^17] (Hydrogen is a toolkit for building the storefront; it is not itself a hosted
  theme.)
- **Oxygen** — Shopify's global edge hosting for Hydrogen storefronts, a worker-based
  JavaScript runtime built on Cloudflare's open-source **workerd**, included on all paid
  Shopify plans at no extra cost.[^18] A common alternative is a custom Next.js/Vercel
  frontend against the Storefront API instead of Hydrogen/Oxygen.[^17][^18]

Beyond a single vendor's headless stack sits **MACH** — a broader composable-architecture
philosophy promoted by the **MACH Alliance**, standing for **M**icroservices,
**A**PI-first, **C**loud-native SaaS, and **H**eadless.[^19] The idea is to assemble
best-of-breed, independently deployable, API-connected SaaS components (a separate CMS,
search, cart, checkout, PIM, etc.) rather than buy one monolithic suite — **commercetools**
is a canonical MACH-aligned commerce backend. MACH maximizes
flexibility and swappability at the cost of integration and operational complexity: you
own the seams between components.[^19]

**When headless/MACH is justified vs. not:** it trades implementation complexity and
ongoing engineering cost for full control over frontend performance, bespoke UX, and
multi-brand/multi-region frontend flexibility that a themed hosted storefront can't match.
Choose it when there is an actual, articulated frontend-customization or performance
requirement — not by default. A themed Shopify storefront plus apps serves the large
majority of merchants at a fraction of the engineering overhead (see Anti-patterns).

## Anti-patterns

- **No-code tool in a high-volume, latency-sensitive path.** Wiring Zapier/Make into
  real-time multi-channel inventory sync (or any high-frequency critical flow) instead of a
  purpose-built sync tool leads to task/credit burn, polling-interval delay, task-queue
  backlog, and ultimately overselling. Reserve no-code for low-frequency, non-critical
  glue.[^15][^16]
- **Repricing with no floor price.** Automated marketplace repricing without a min/floor can
  race a price to an unprofitable level during a competitor price war — the exact scenario
  floor/ceiling bands and game-theory "cooperative" strategies exist to prevent.[^8][^9]
- **No frequency capping across lifecycle flows.** Running welcome + abandoned cart +
  browse-abandon + post-purchase + win-back concurrently without global frequency capping
  causes message fatigue that raises unsubscribe/opt-out rates and, for SMS, wastes
  per-message cost on the highest-return channel. Cap sends per subscriber across *all*
  flows, not per flow.[^11][^13]
- **Copying email cadence onto SMS.** SMS abandoned-cart best practice is a much faster
  first touch (commonly cited as within the first ~15–30 minutes) with one follow-up;
  reusing an email's multi-hour, multi-touch cadence[^10] on SMS misses the intent window
  and over-messages.
- **Two systems both owning the stock count.** Running two "sync" tools that each treat
  themselves as the source of truth guarantees drift and overselling. Designate one
  authoritative available-to-sell ledger.[^3][^5][^6]
- **Headless for its own sake.** Choosing headless/MACH without a concrete
  frontend-customization or performance requirement multiplies engineering overhead versus
  a themed hosted storefront plus apps, for benefits the business won't use.[^17][^19]
- **Treating "Skubana" as a current product name.** It was rebranded to Extensiv Order
  Management in 2022; using the legacy name signals stale knowledge to operators.[^1][^2]

## References

[^1]: Etail Solutions — "Extensiv Order Management (Formerly Skubana)" integration page. States Skubana was acquired by 3PL Central in 2021 and rebranded to Extensiv in 2022; Skubana combined with 3PL Central, Scout, and CartRover under the Extensiv brand. https://www.etailsolutions.com/integrations/extensiv-order-management-formerly-skubana (fetched July 2026)
[^2]: Supply & Demand Chain Executive — "3PL Central Acquires Skubana to Expand Order Management for E-Commerce" (2021 acquisition announcement). https://www.sdcexec.com/transportation/3pl-4pl/news/21366595/3pl-central-3pl-central-acquires-skubana-to-expand-order-management-for-ecommerce (searched July 2026)
[^3]: Extensiv — Order Management System product page and "Extensiv Order Management (formerly Skubana)" capability description: rules-engine "algorithmic fulfillment," dropship/3PL routing, cross-docking, Amazon MCF, multi-channel inventory visibility across warehouses/3PLs/dropshippers/FBA, auto-generated POs from sales velocity/lead time/seasonality. https://www.extensiv.com/order-management-system (searched/fetched July 2026)
[^4]: ShipStation Help Center — "Automation Basics," "Introduction to Automation Rules," and Auto-Routing documentation: rules apply on import to Awaiting Shipment; auto-select carrier/service/package; Auto-Routing by stock location and recipient proximity (US/Canada); Auto-Split; support for self-fulfill/FBA/3PL/dropship. https://help.shipstation.com/hc/en-us/articles/360025870192-Automation-Basics (searched July 2026)
[^5]: Cin7 — "Cin7 Core or Omni: Which One Is Right for Your Business?" and product/solutions pages: Core (ready-to-go inventory/OMS) vs Omni (customizable omnichannel, native EDI/3PL, multi-entity); real-time inventory across online/retail/wholesale/POS; 700+ integrations. https://www.cin7.com/blog/cin7-core-or-omni-which-one-is-right-for-your-business/ (searched July 2026)
[^6]: Linnworks — Order Management, Inventory Management, and Multichannel Listings feature pages: inventory linked to orders/warehouses/listings/shipping so a sale anywhere updates stock everywhere; rules-based routing by channel/stock location/carrier/packaging; two-way connections to Amazon, eBay, Shopify, Walmart, TikTok Shop, and more. https://www.linnworks.com/features/order-management/ (searched July 2026)
[^7]: DSers — Shopify App Store listing ("DSers: Dropship+AliExpress+AI"): official AliExpress dropshipping solution and Alibaba Group partner; succeeded Oberlo (which shut down in 2022); automates browse→import→auto-order fulfillment from AliExpress/1688/Alibaba suppliers to Shopify/WooCommerce/Wix. https://apps.shopify.com/dsers (searched July 2026)
[^8]: Seller Snap — "Amazon AI Algorithmic Repricer" product pages: game-theory-based repricing aiming to win the Buy Box at the highest sustainable price; "cooperative" / "wait my turn" strategy to avoid price wars; contrasted explicitly with naive rule-based undercutting. https://sellersnap.io/amazon-ai-algorithmic-repricer/ (searched July 2026)
[^9]: RepricerExpress Help Center — "Choosing a repricing rule and setting Min and Max prices" and related pricing-rule docs: per-SKU min (floor)/max (ceiling); never uploads prices outside the band; basic vs advanced rules; behavior options when competitor prices below min (Do Not Reprice / Go to Min / Go to Max / Compete with Next Best Seller) and when holding the Buy Box (default no change; advanced allows increase / reprice as normal). https://help.repricerexpress.com/hc/en-us/articles/360000044357-Choosing-a-repricing-rule-and-setting-Min-and-Max-prices (searched July 2026)
[^10]: Klaviyo Help Center — "How to create an abandoned cart flow": pre-built flow uses a 4-hour delay before the first message; recommends first message 2–4 hours after checkout starts, a second 20–48 hours later, and 2–3 messages total for optimal performance; longer delays for high-consideration items. https://help.klaviyo.com/hc/en-us/articles/115002779411 (fetched July 2026)
[^11]: Klaviyo — "Abandoned Cart Benchmark Report: Rates & Statistics": abandoned-cart flow average revenue per recipient ~$3.65 and average placed-order (conversion) rate ~3.33% (highest of any flow); top-10% RPR ~$28.89; industry variation. https://www.klaviyo.com/blog/abandoned-cart-benchmarks (fetched July 2026)
[^12]: Klaviyo Help Center — "Benchmarks for flow emails reference": flow vs campaign performance, RPR by flow type; flows generate a disproportionate share of email revenue from a small share of sends. https://help.klaviyo.com/hc/en-us/articles/360033669452 (searched July 2026)
[^13]: Postscript — "SMS Benchmarks" (2026 report, 2025 data, 17,000+ Shopify stores): abandoned-cart earnings per message ~$3.52–$10.95 (25th–75th percentile), conversion ~3.97%–7.84%, CTR ~9.5%–17.3%; overall median revenue per message ~$0.98; median subscriber LTV ~$71.20. (The benchmarks page reports these stats; it does not contain a "34x ROI"/"$8.11 EPM" figure. The 15–30-minute first-touch timing cited in the body is general SMS best practice, not a figure from this page.) https://postscript.io/sms-benchmarks (fetched July 2026)
[^14]: Attentive — product site: AI-powered conversational commerce platform offering two-way, real-time personalized SMS and email for retail/ecommerce. https://www.attentive.com/ (searched July 2026)
[^15]: Zapier Help Center — "How pay-per-task billing works in Zapier": on exceeding a plan's task limit, Zapier switches to pay-per-task for the rest of the period; maximum pay-per-task usage is capped at 3× the plan's task limit (base allocation + up to 2× overage); pay-per-task rates change effective July 15, 2026. https://help.zapier.com/hc/en-us/articles/15279018245901-How-pay-per-task-billing-works-in-Zapier (fetched July 2026)
[^16]: Make — "Pricing & Subscription Packages": credits (formerly operations) model; each module action consumes one credit; a run can use "anywhere from two credits to thousands"; Free = 1,000 credits/month with a 15-minute minimum run interval; Core = $12/month for 10,000 credits with 1-minute scheduling and unlimited active scenarios; Router/error-handler modules consume zero credits. https://www.make.com/en/pricing (fetched July 2026)
[^17]: Shopify — Hydrogen documentation and Storefront API: Hydrogen is Shopify's React-based headless framework built on the open-source React Router framework, with API clients for the Storefront and Customer Account APIs. https://shopify.dev/docs/api/hydrogen/latest (searched July 2026)
[^18]: Shopify — "Getting started with Hydrogen and Oxygen": Oxygen is Shopify's global edge deployment platform for Hydrogen storefronts, a worker-based JavaScript runtime built on Cloudflare's open-source workerd, available on all paid plans at no extra cost. https://shopify.dev/docs/storefronts/headless/hydrogen/getting-started (searched July 2026)
[^19]: MACH Alliance — "MACH at Different Levels" / insights hub: MACH = Microservices, API-first, Cloud-native SaaS, Headless; a non-profit promoting composable, best-of-breed, independently deployable, API-connected architecture. https://machalliance.org/insights-hub/mach-at-different-levels (fetched July 2026)
[^20]: Uptek — "Shopify App Store Statistics": ~11,905 apps as of Nov 14, 2024 (count fluctuates with Shopify's quality purges); ~87% of merchants use apps; average merchant installs ~6 apps. Third-party App Store tracker data, not a Shopify-official figure. https://uptek.com/shopify-statistics/app-store/ (fetched July 2026)
[^21]: Zapier — Plans & Pricing: task-metered plans; free tier on the order of ~100 tasks/month, scaling into hundreds-of-thousands-to-millions of tasks on paid tiers; lower tiers poll less frequently. https://zapier.com/pricing (searched July 2026)
