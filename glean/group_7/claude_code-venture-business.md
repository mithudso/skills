# venture-business

**Category:** AI, Agents & Prompt Engineering
**Platform:** Claude Code
**Original Path:** claude-code/venture-business

## Description
Hub for a North Carolina FOR-PROFIT solo/lean founder doing it themselves — entity, money, real estate, marketing; routes to 10 spokes. Educational, NOT legal/tax/financial advice (as of 2026). TRIGGER: NC entity formation & choice (LLC/PLLC/S-corp/C-corp, sole prop), EIN, DBA, registered agent, NC business tax (sales/use, withholding, franchise, corporate income), CTA/FinCEN BOI, annual report; becoming an employer & running NC payroll (SUTA/DES, workers comp, 941/940, I-9, W-2-vs-1099 misclassification); entity lifecycle (dissolution, conversion, foreign qualification, reinstatement, PTET/SALT-cap); business planning, financial modeling, unit economics, break-even, pricing, funding (SBA/grants/NC IDEA); pitch deck, 3-statement model, cap table, SAFE/priced round; NC real-estate licensing & brokerage law (NCREC, GS 93A, BIC, trust accounts, agency, fair housing); real-estate marketing & investing (cap rate, BRRRR, flip, DSCR, 1031, NAR settlement); property-management law, short-term/vacation rentals, seller financing; founder marketing strategy & local SEO (Google Business Profile, CAN-SPAM/TCPA/FTC); Canva brand & print collateral. SKIP: nonprofit/cause/501(c)(3)/fundraising/organ-donation → venture-nonprofit-cause (sibling hub); personal finance (individual income tax, banking, insurance, investing, estate) → consumer-finance; consumer credit/debt → consumer-credit-and-debt.

---

# NC For-Profit Venture (hub)

The front door for a **North Carolina solo/lean for-profit founder** handling their own formation, money, real estate, and marketing. This hub routes to 10 reference spokes. Its **sibling hub, `venture-nonprofit-cause`,** owns the 501(c)(3)/charitable side; route there for nonprofit formation, fundraising, cause marketing, and the organ-donation cause work.

**When activated:** match the question to a row, then **Read the listed `references/` file** before answering at depth — the table is only a router. Everything here is **educational information, not legal, tax, or financial advice**; rules and figures change and are date-stamped *as of 2026* — verify with the NC Secretary of State, NCDOR, IRS, NCREC, or a licensed professional before acting.

## Routing table

| Spoke | Use when | Reference file |
| --- | --- | --- |
| `venture-nc-business-formation-tax` | Entity choice (LLC/PLLC/S-corp/C-corp/sole prop), EIN, registered agent, DBA, NC business tax, CTA/BOI, annual report. | `references/venture-nc-business-formation-tax.md` |
| `venture-nc-employer-payroll` | First hire & NC payroll: NC-BR withholding, SUTA/DES, new-hire reporting, workers' comp, FICA/FUTA, 941/940, I-9, W-2 vs 1099. | `references/venture-nc-employer-payroll.md` |
| `venture-nc-entity-lifecycle` | Dissolution, conversion (sole-prop→LLC, LLC→corp), foreign qualification, reinstatement, the PTET/SALT-cap election. | `references/venture-nc-entity-lifecycle.md` |
| `venture-small-business-planning` | Business model (Canvas/Lean), idea/market validation, TAM/SAM/SOM, unit economics, break-even, pricing, funding options. | `references/venture-small-business-planning.md` |
| `venture-startup-fundraising-deck-model` | Investor pitch deck, integrated 3-statement model, cap table & dilution, SAFE/convertible/priced round; nonprofit grant lifecycle. | `references/venture-startup-fundraising-deck-model.md` |
| `venture-nc-real-estate-law` | NC real-estate licensing & brokerage law: NCREC, GS 93A, provisional broker, BIC, trust accounts, agency, RPOADS, fair housing. | `references/venture-nc-real-estate-law.md` |
| `venture-real-estate-marketing-investing` | RE marketing (listings, farming, IDX/SEO, NAR settlement) & investing (cap rate, cash-on-cash, BRRRR, flip, DSCR, 1031, REITs). | `references/venture-real-estate-marketing-investing.md` |
| `venture-real-estate-advanced` | NC property-management law, short-term/vacation rentals (NC Vacation Rental Act, ADR/RevPAR), creative/seller financing. | `references/venture-real-estate-advanced.md` |
| `venture-marketing-strategy-local-seo` | Founder marketing strategy, STP/positioning, funnel/AARRR, local SEO & Google Business Profile, marketing compliance (CAN-SPAM/TCPA/FTC). | `references/venture-marketing-strategy-local-seo.md` |
| `venture-canva-founder-brand-stack` | DIY brand/design studio in Canva: Brand Kit/Hub, Magic Studio, marketing & print collateral, print specs (CMYK/bleed/300 DPI). | `references/venture-canva-founder-brand-stack.md` |

## Cross-cutting notes

- **Sibling hub — `venture-nonprofit-cause`.** Anything 501(c)(3): nonprofit formation, charitable-solicitation registration, Form 990, fundraising operations (Ad Grants, CRM, email), cause/donor marketing, and the organ-donation cause venture. The two hubs share one founder's toolkit; start here for the commercial side. (Note: `venture-startup-fundraising-deck-model` also carries the nonprofit grant-writing lifecycle, kept with the for-profit fundraising material.)
- **Personal vs business boundary.** Individual income tax, personal banking/insurance/investing, and estate planning are **personal** finance → `consumer-finance`. This hub is strictly the **business/employer** side. Consumer credit, debt, and collections law → `consumer-credit-and-debt`.
- **North Carolina-specific** throughout: NC SoS filings, NCDOR tax, NCREC licensing, NC Vacation Rental Act, NC IDEA / SBTDC / SBA-NC resources. Verify current rules at the primary source before relying.

<!-- cross-hub-map -->

## Cross-hub map — where every venture topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `venture-business` | NC For-Profit Venture: Formation, Funding, Real Estate & Marketing (hub) | `references/venture-nc-business-formation-tax.md`, `references/venture-nc-employer-payroll.md`, `references/venture-nc-entity-lifecycle.md`, `references/venture-small-business-planning.md`, … |
| `venture-nonprofit-cause` | NC Nonprofit & Cause Venture: Formation, Fundraising, Cause Marketing & Organ Donation (hub) | `references/venture-nc-nonprofit-formation.md`, `references/venture-nonprofit-fundraising-ops.md`, `references/venture-cause-nonprofit-marketing.md`, `references/venture-donor-community-engagement.md`, … |