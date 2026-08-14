---
name: ecommerce-fundamentals
description: "Ecommerce business hub — platform selection (Shopify/BigCommerce/WooCommerce/Magento vs Amazon/Etsy marketplaces), unit economics (CAC, AOV, contribution margin, LTV:CAC), fulfillment models (3PL, FBA, dropship, in-house), payments (gateways, interchange, chargebacks, PCI DSS), checkout optimization, marketplace-vs-DTC tradeoffs. TRIGGER: choosing a platform; unit-economics (CAC/AOV/LTV/margin); fulfillment/3PL choice; payment gateway/interchange/chargeback; cart abandonment; marketplace vs DTC strategy. SKIP: recurring billing/churn/dunning → subscription-commerce; dropship sourcing/margin → dropshipping-and-fulfillment; order routing/inventory sync/email-SMS automation → ecommerce-automation-and-tools; checkout FRONTEND → frontend-ui; payment SECURITY review → security-review; marketing copy → content-and-marketing-writing."
version: "1.0.0"
updated: "2026-07-03"
category: business
whenToUse: Use when advising on ecommerce platform choice, unit economics, fulfillment strategy, payments, checkout conversion, or marketplace vs DTC decisions.
keywords: [ecommerce, shopify, bigcommerce, woocommerce, magento, marketplace, amazon, unit economics, CAC, AOV, LTV, contribution margin, fulfillment, 3PL, FBA, payment gateway, interchange, chargeback, PCI DSS, checkout abandonment, DTC]
tags: [ecommerce, business, retail]
related_skills: [subscription-commerce, dropshipping-and-fulfillment, ecommerce-automation-and-tools]
---

# Ecommerce Fundamentals

Hub for the ecommerce-business skill cluster. See `references/` for depth; siblings: `subscription-commerce`, `dropshipping-and-fulfillment`, `ecommerce-automation-and-tools`.

## Quick reference — platform decision

| Need | Platform | Why |
|---|---|---|
| Fast launch, low technical overhead | Shopify | Hosted, largest app ecosystem, best checkout conversion benchmarks |
| Content-heavy or WordPress-native site | WooCommerce | Full CMS control, self-hosted, plugin flexibility |
| Enterprise B2B or complex catalog | BigCommerce / Adobe Commerce (Magento) | Native B2B features, multi-storefront, open architecture |
| Reach existing buyer intent, no site needed | Amazon / Etsy / Walmart Marketplace | Built-in traffic, but margin and customer-data tradeoffs |
| Headless/composable at scale | Shopify Hydrogen, BigCommerce headless, commercetools | Custom frontend, API-first backend |

## Quick reference — core unit-economics formulas

| Metric | Formula | Healthy signal |
|---|---|---|
| CAC | Total acquisition spend / new customers acquired | Lower than 1/3 of LTV |
| AOV | Total revenue / number of orders | Rising over time via bundling/upsell |
| Contribution margin | (Revenue − COGS − variable costs) / Revenue | Must cover CAC within payback window |
| LTV:CAC | Customer lifetime value / CAC | ≥3:1 is the common benchmark |
| CAC payback period | CAC / (monthly contribution margin per customer) | <12 months for most DTC brands |

See `references/ecommerce-fundamentals-context.md` for fulfillment models, payment processing, checkout optimization, and marketplace-vs-DTC tradeoffs in full depth.
