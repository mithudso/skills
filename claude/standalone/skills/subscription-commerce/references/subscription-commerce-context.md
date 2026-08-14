# Subscription Commerce — Context

## Contents
- Recurring billing mechanics
- MRR/ARR and cohort math
- Churn taxonomy and diagnosis
- Dunning management
- Subscription-box logistics
- Billing platform comparison
- Anti-patterns
- References

## Recurring billing mechanics

A subscription-billing engine must handle: tokenized card-on-file storage (via a PCI-scoped vault, not raw card storage), scheduled charge generation, proration on plan changes, tax recalculation per billing cycle, and automatic retry orchestration on failure. Platforms like Stripe Billing, Chargebee, Recurly, and — for Shopify-native merchants — Recharge and Bold Subscriptions abstract this so merchants don't build a billing engine from scratch.[^1]

## MRR/ARR and cohort math

MRR must be **normalized**: an annual plan is divided by 12, a quarterly plan by 3, and one-time fees (setup, shipping) are excluded so MRR reflects only recurring revenue. ARR is simply MRR × 12 and is mainly a communication convenience for investors, not a distinct calculation.

Cohort analysis (grouping subscribers by signup month and tracking retention over subsequent months) is the standard way to see whether churn is improving, because blended churn rates can mask a recent-cohort improvement hidden by legacy-cohort decay.

## Churn taxonomy and diagnosis

- **Voluntary churn:** the customer actively cancels (dissatisfaction, price, no longer needs product). Diagnosed via cancellation-flow surveys and exit interviews.
- **Involuntary churn:** the subscription lapses because a payment failed (expired card, insufficient funds, bank fraud-block) and was never actively cancelled. Industry analyses commonly attribute a large minority (often cited in the 20-40% range) of total subscription churn to this cause — meaning dunning quality is a direct retention lever, not just an operational nicety.[^2]
- **Logo churn vs revenue churn:** logo churn counts canceled accounts; revenue churn counts lost MRR. A low-value-plan-heavy cancellation wave can show alarming logo churn with modest revenue churn, or vice versa with a single large-account loss.
- **Net Revenue Retention (NRR):** the single most-watched metric by subscription-business investors because NRR >100% means the existing customer base grows revenue even with zero new signups (via upsell/expansion outweighing churn+contraction).

## Dunning management

Effective dunning combines:
1. **Smart/adaptive retry logic** — network-level tools (Stripe's Smart Retries, Visa/Mastercard account-updater services) time retries around likely-successful windows (e.g., after a payday) rather than fixed intervals.
2. **Card-updater services** — automatically refresh expired/reissued card numbers on file via the card networks, preventing a chunk of involuntary churn before it starts.
3. **Multi-channel customer notification** — email plus SMS with a direct "update payment method" link measurably outperforms email-only in recovery rate.
4. **Grace periods** — keeping access active for a defined window during retries reduces the "feels like I was punished" friction that drives voluntary churn triggered by a failed-payment experience.
5. **Win-back offers** — a discount or paused-not-cancelled option presented at the lapse point can recover a portion of subscribers who would otherwise churn.

## Subscription-box logistics

Subscription boxes (curated recurring physical-product shipments — e.g., beauty, snacks, hobby-kit boxes) layer physical fulfillment complexity onto billing complexity:

- **Curation cadence:** monthly is standard; the curation/sourcing lead time (often 60-90 days for private-label or licensed products) must be planned well ahead of the billing cycle it ships against.
- **Box-cost economics:** must account for product cost, box/packaging cost, pick-pack labor, and shipping — all before contribution margin is calculated; box companies typically target a specific perceived-retail-value multiple (e.g., "$50+ value for $25") to drive perceived deal quality and reduce churn.
- **Billing-to-fulfillment timing:** the billing date and the "cutoff for this cycle's box contents" date are usually different — customers who subscribe after the cutoff roll into the next cycle, which must be communicated clearly to avoid support tickets and early churn.
- **Skip/pause/swap features:** giving subscribers control (skip a month, swap an item) is consistently associated with lower voluntary churn than a rigid ship-every-month model, because it removes the "I have too much of this already" cancellation trigger.

## Billing platform comparison

| Platform | Best for | Note |
|---|---|---|
| Recharge | Shopify-native subscription/box merchants | Deep Shopify checkout integration, large app-store presence |
| Bold Subscriptions | Shopify merchants wanting simpler setup | Smaller feature set than Recharge, often cheaper |
| Chargebee | Mid-market to enterprise SaaS + ecommerce hybrid | Strong metered/usage billing, revenue-recognition tooling |
| Stripe Billing | Developer-first, custom storefronts, headless | Most flexible via API, most implementation effort |
| Recurly | Enterprise subscription businesses | Strong dunning/analytics focus, enterprise support |

## Anti-patterns

- Treating all churn as voluntary and only surveying cancellers — ignores the involuntary-churn lever, which is often the cheapest to fix (dunning improvements have high ROI).
- Reporting blended churn without cohort breakdown, masking whether retention is actually improving.
- No skip/pause option, forcing customers into a binary keep-or-cancel decision that inflates voluntary churn.
- Billing cutoff dates not communicated, generating support load and early cancellations from confused first-cycle subscribers.

## References
[^1]: Stripe Billing, Chargebee, and Recharge product documentation on recurring-billing architecture (accessed 2026).
[^2]: Involuntary-churn share estimates, Recurly "State of Subscriptions" research and Baremetrics/ProfitWell subscription-analytics benchmark reports (accessed 2026).
