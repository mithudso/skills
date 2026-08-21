# JPMC — Technology Scale (Spend, Headcount, Engineering Org)

> Provenance: reference under the `jpmorgan-chase-account-intelligence` skill. Built 2026-06-18 via `/dr` deep research (exa + WebSearch/WebFetch; primary sources). **verified-as-of: 2026-06-18.** Every dollar figure and headcount here is **point-in-time and re-stated annually** — re-verify against the latest company-update/Investor-Day before customer use.

## Contents
- Annual technology spend (dated trajectory)
- Run-the-bank vs investment split
- Technologist / engineer headcount (dated)
- Engineering-org footprint
- TAM implications
- Sources

## Annual technology spend (dated trajectory)

JPMC is the **largest corporate technology spender in financial services**, and Dimon frames technology as **"~10% of revenues … table stakes … for the rest of eternity."**[^1] The budget is re-stated each year (typically at the January company update and May Investor Day), so always cite the year:

| Fiscal year | Technology spend | Source tier |
|---|---|---|
| 2022 | ~$14B | press / Investor Day[^2] |
| 2023 | ~$15B (one ops view cited ~$15.3B) | press[^3] |
| 2024 | **~$17B** (up ~$1.5B YoY) | Investor Day transcript[^4] |
| 2025 | **~$18B** (up ~$1B / ~6% YoY) | Investor Day[^5] |
| **2026** | **~$19.8B** (up ~$2B / ~10% YoY; "nearly $20B") | 2026 company update / press[^6] |

The **~$19.8B figure for 2026 is the latest disclosed** and matches the Global CIO's current official bio.[^7] For scale comparison, public reporting puts Bank of America's tech budget around ~$12–14B.[^6] **Confidence: fact** (multi-source, primary + press); **volatility: high** (re-stated annually).

A negation check confirmed the spend has **risen every year** — there is no public evidence of a tech-budget cut. What *is* moderating is **modernization** spend and **headcount** growth, not total spend (see below and `jpmc-cloud-data-modernization.md`).[^5]

## Run-the-bank vs investment split

- **FY2024:** roughly **half "run the bank"** (infrastructure, software licenses, application/production support) and **half "investment,"** with the investment portion split into **~$4.5B products/platforms/UX** and **~$3.0B modernization** (cloud, dev efficiency, cyber/resiliency).[^8]
- **FY2025:** investment portion **~$8.0B total** (products/platforms/features ~$5.3B + other ~$2.7B), broken out by LOB (CCB $3.2B / CIB $3.7B / AWM $1.0B). JPMC stated it is **"past the point of peak modernization spend."**[^5][^9]

## Technologist / engineer headcount (dated)

- **Current (2026) official figure: ~65,000 technologists**, per the Global CIO's bio (the org she runs against the $19.8B budget).[^7]
- **~44,000 software engineers** specifically (Global CIO at AWS re:Invent, 2024-12).[^10]
- Trajectory and corroboration: ">60,000 technologists" in 2024 LOB CEO shareholder letters;[^11] a "63,000-person team" (2024 press / case study). The older "~50,000+ technologists" figure (2020-era) is correct as a historical anchor but is **superseded** — do not present it as current.[^7][^10]
- **40,000+ engineers** now use AI coding assistants (2025 Investor Day).[^9]
- A non-primary outlier (CIO.inc) cited "60,000 developers + 80,000 operations/call-center" attributed to President Daniel Pinto — broader than the engineer count; **treat as tentative** and prefer the bio's ~65,000 technologists / ~44,000 software engineers framing.

## Engineering-org footprint

Public primary statements (2024–2025) describe the estate as:[^10][^11]
- **>6,000 applications**
- **~1 exabyte of data** ("nearly an exabyte," Global CIO, AWS re:Invent 2024 — single high-profile quote, widely repeated; treat as qualified)
- **Multi-cloud** (AWS / Azure / GCP) **plus a substantial private cloud** (see `jpmc-cloud-data-modernization.md`)
- Moves **~$10 trillion/day**
- Provides technology serving **>90% of the Fortune 500** and handling a large share of US e-commerce transactions
- AI footprint: **~100 GenAI solutions in production**; LLM Suite used across the firm (see `jpmc-ai-strategy.md`)

## TAM implications

- **Budget is not the constraint; proof and standards are.** With a ~$19.8B budget, JPMC can fund what it values — the gate is the Supplier Minimum Control Requirements and architectural fit, not price (see `jpmc-vendor-procurement-and-mongodb-signals.md`). Cost is optimized last.
- **"Past peak modernization" reframes the pitch.** The infrastructure-migration wave is maturing; the live spend is moving to **application code + data modernization for AI-readiness** — the data-platform conversation (see `jpmc-cloud-data-modernization.md`).
- **Always date your numbers.** If you quote a tech-spend or headcount figure to the customer, state the year and re-verify — these are the fastest-decaying facts in this profile.

## Sources

1. JPMorgan Chase 2025 Investor Day, closing Q&A (PDF) — Dimon: technology is "table stakes … ~10% of revenues." (Investor Day transcript, primary) — https://www.jpmorganchase.com/ir/investor-day
2. Banking Dive, "JPMorgan Chase IT: Lori Beer, $14 billion budget" (2022) — ~$14B 2022 budget. (Press) — https://www.bankingdive.com/news/jpmorgan-chase-IT-lori-beer-14-billion-budget-6000-apps/633135/
3. Constellation Research, "JPMorgan Chase digital transformation, AI and data strategy" — ~$15.3B 2023 tech budget; vendor-rationalization lever. (Press/analyst) — https://www.constellationr.com/insights/news/jpmorgan-chase-digital-transformation-ai-and-data-strategy-sets-generative-ai
4. JPMorgan Chase 2024 Investor Day, firm-overview transcript (PDF) — FY2024 ~$17B tech spend; $4.5B products/platforms + $3B modernization. (Investor Day transcript, primary) — https://www.jpmorganchase.com/ir/investor-day
5. Banking Dive, "JPMorgan investor day" (2025-05-19) — FY2025 ~$18B (~6% increase); "past peak modernization spend." (Press) — https://www.bankingdive.com/news/jpmorgan-investor-day-ceo-dimon-headcount-ai-banking/748479/
6. Business Insider / AOL, "JPMorgan to spend almost $20 billion" (2026-02-24) — FY2026 ~$19.8B (~10% increase); BofA ~$14B comparison. (Press) — https://www.aol.com/news/jpmorgan-spend-almost-20-billion-000403027.html
7. Lori Beer official JPMorganChase bio — ~$19.8B budget; ~65,000 technologists. (Leadership page, primary) — https://www.jpmorganchase.com/about/leadership/lori-beer
8. TechTarget, "JPMorgan Chase technology goal: innovation with cost control" (2024-08-26) — FY2024 run-vs-invest ~50/50; $4.5B products; $3B modernization. (Press) — https://www.techtarget.com/searchcio/feature/JPMorgan-Chase-technology-goal-Innovation-with-cost-control
9. JPMorgan Chase 2025 Investor Day, full presentation (PDF) — investment ~$8.0B; by-LOB split; 40K+ engineers using AI assistants. (Investor Day deck, primary) — https://www.jpmorganchase.com/ir/investor-day
10. JPMorganChase Technology blog, "Global CIO showcases advancements at AWS re:Invent" (2024-12) — ~44,000 software engineers; ~exabyte of data; >90% of Fortune 500. (Blog, primary) — https://www.jpmorganchase.com/about/technology
11. JPMorgan Chase 2024 LOB CEO letters to shareholders (PDF) — ">60,000 technologists"; 6,000 apps; ~exabyte data. (Shareholder letter, primary) — https://www.jpmorganchase.com/ir/annual-report
