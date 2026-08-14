# Dropshipping and Fulfillment — Context

## Contents
- Supplier sourcing and vetting
- Print-on-demand
- Fulfillment automation tooling
- Margin structure and pricing
- Anti-patterns
- References

## Supplier sourcing and vetting

Dropshipping removes inventory-holding risk by having a third-party supplier ship directly to the end customer, but this shifts risk into supplier reliability, quality control, and shipping time.

- **AliExpress (via DSers or similar connectors):** the historical default entry point — huge catalog, low minimums, but typically 1-4+ week shipping to Western markets unless the supplier uses a regional warehouse.
- **CJ Dropshipping:** offers private-label/branding options and China-based warehousing with often faster processing than raw AliExpress sourcing, plus product-sourcing/quality-check services.
- **Spocket / AutoDS:** curate US/EU-based suppliers for shorter transit times at a materially higher per-unit cost — the standard tradeoff between speed and margin.
- **Vetting checklist:** sample-order the product yourself before listing it; check supplier processing-time consistency (not just advertised shipping time); confirm return/defect policy; verify the supplier can sustain volume before scaling ad spend into a SKU.

## Print-on-demand

Print-on-demand (POD) is a dropshipping variant where the product (apparel, mugs, books, wall art) is manufactured only after the sale, using a customer's custom design.

- **Printful and Printify** are the dominant Shopify/Etsy-integrated POD providers; Printify operates a marketplace of print partners so a seller can choose the balance of cost vs. speed vs. print quality per product, while Printful is more vertically integrated with its own facilities.
- POD carries **zero inventory risk** (nothing is produced until ordered) but the highest per-unit cost of the fulfillment models compared here, and adds production time (typically 2-7 business days) on top of shipping time.
- Because there's no minimum order quantity, POD is well suited to testing many designs/niches cheaply before committing capital elsewhere.

## Fulfillment automation tooling

Order-routing apps (DSers for AliExpress/CJ, the native Printful/Printify Shopify apps, AutoDS) automatically push each storefront sale to the correct supplier as a purchase order, sync tracking numbers back to the storefront, and in some cases auto-adjust storefront pricing when a supplier's cost changes. This automation layer is what makes dropshipping operationally viable at volume — manually placing each supplier order does not scale past a handful of orders per day.

## Margin structure and pricing

Dropship margins are structurally thinner than owned-inventory retail because:
1. The supplier's per-unit price already includes their margin, leaving less room for the seller's markup before becoming uncompetitive.
2. Multiple sellers frequently list the identical AliExpress/CJ product, creating price competition on a commodity item.
3. Paid-ad CAC must be covered entirely by the remaining margin, which is why successful dropship operators often emphasize **niche differentiation, bundling, or branding** (private label via CJ, or POD custom designs) rather than competing purely on price against sellers of an identical generic product.

Typical dropship gross margins run 15-30%; POD often runs slightly higher (20-40%) because the design/brand adds differentiation the buyer can't price-shop identically elsewhere, but the base per-unit cost is also higher.

## Anti-patterns

- Listing a product with no sample-order quality check — the fastest way to accumulate refunds/chargebacks and marketplace/platform account risk.
- Advertising aggressively into a SKU before confirming supplier processing-time consistency at higher volume.
- Ignoring shipping-time expectations in ad creative/product copy — long-transit dropship items generate disproportionate support tickets and negative reviews if the customer wasn't told upfront.
- Treating margin as static — supplier prices on AliExpress/CJ-style marketplaces can change without notice; pricing automation (see `ecommerce-automation-and-tools`) is needed to avoid selling at a loss after a supplier cost increase.

## References
[^1]: CJ Dropshipping and DSers platform documentation on supplier sourcing and order automation (accessed 2026).
[^2]: Printful and Printify product/vendor documentation on print-on-demand fulfillment models (accessed 2026).
[^3]: Spocket and AutoDS platform documentation on US/EU supplier sourcing tradeoffs (accessed 2026).
