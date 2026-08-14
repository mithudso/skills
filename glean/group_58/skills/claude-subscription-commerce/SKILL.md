---
name: subscription-commerce
description: "Subscription and recurring-revenue business models — recurring billing mechanics, MRR/ARR calculation, churn types (voluntary/involuntary/logo/revenue) and retention math, dunning management for failed payments, subscription-box logistics (curation, cadence, box-cost economics), and billing-platform choice (Recharge, Bold, Stripe Billing, Chargebee). TRIGGER: calculating or modeling MRR/ARR/churn/net-revenue-retention; designing a dunning/failed-payment recovery flow; choosing a subscription billing platform; subscription-box cadence/curation/logistics design; involuntary-churn (card-decline) reduction; subscription pricing/plan design (tiers, trials, pause/skip options). SKIP: general ecommerce unit economics (CAC/AOV) → ecommerce-fundamentals; dropship supplier/fulfillment mechanics → dropshipping-and-fulfillment; the specific dropship+subscription HYBRID operational model (inventory-timing conflicts) → dropshipping-and-fulfillment (references/dropship-subscription-hybrid.md); automated email/SMS retention flows → ecommerce-automation-and-tools; SaaS-specific product-led-growth motion → software-engineering-patterns."
version: "1.0.0"
updated: "2026-07-03"
category: business
whenToUse: Use when designing, modeling, or troubleshooting a recurring-billing/subscription-commerce business — churn, dunning, MRR, billing platforms, or subscription-box logistics.
keywords: [subscription, recurring billing, MRR, ARR, churn, involuntary churn, dunning, net revenue retention, subscription box, Recharge, Chargebee, Bold Subscriptions, Stripe Billing, cohort retention]
tags: [ecommerce, subscription, business]
related_skills: [ecommerce-fundamentals, dropshipping-and-fulfillment, ecommerce-automation-and-tools]
---

# Subscription Commerce

Sibling of `ecommerce-fundamentals` (unit economics), `dropshipping-and-fulfillment` (the dropship+subscription hybrid model lives there), and `ecommerce-automation-and-tools` (dunning/retention automation tooling).

## Quick reference — churn & revenue formulas

| Metric | Formula | Note |
|---|---|---|
| MRR | Sum of normalized monthly recurring revenue across active subscribers | Exclude one-time charges |
| ARR | MRR × 12 | Standard annualized view |
| Customer/logo churn | Customers lost in period / customers at period start | Watch cohort-level, not blended |
| Revenue churn | MRR lost in period / MRR at period start | Can exceed logo churn with plan downgrades |
| Net Revenue Retention (NRR) | (Starting MRR + expansion − contraction − churn) / Starting MRR | >100% means expansion outpaces churn |
| Involuntary churn share | Failed-payment cancellations / total cancellations | Often 20-40% of all subscription churn |

## Quick reference — dunning cadence pattern

| Attempt | Timing | Channel |
|---|---|---|
| 1st retry | 1-3 days after decline | Card network smart-retry / auto-retry |
| Customer notice | Immediately on decline | Email + SMS with update-payment link |
| 2nd-3rd retry | Days 5-7, then 10-14 | Retry + reminder |
| Grace period lapse | ~2-4 weeks (merchant-set) | Cancel or downgrade + win-back offer |

See `references/subscription-commerce-context.md` for full depth on churn diagnosis, billing-platform comparison, and subscription-box logistics.
