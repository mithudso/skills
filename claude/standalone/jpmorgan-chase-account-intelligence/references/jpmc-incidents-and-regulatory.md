# JPMC — Incidents & Regulatory Orders That Shape Technology Spend

> Provenance: reference under the `jpmorgan-chase-account-intelligence` skill. Built 2026-06-18 via `/dr` deep research (exa + WebSearch/WebFetch; SEC/CFTC/OCC/Fed primary text). **verified-as-of: 2026-06-18.** Factual, public record only — not legal advice. The two highest-stakes orders were re-verified against primary regulator text during authoring.

## Contents
- Why this matters to a TAM
- Public tech incidents
- Regulatory orders touching data/controls/records
- The through-line: what each event funds
- Cross-reference note (regulation generalities)
- Sources

## Why this matters to a TAM

JPMC's technology investment in **resiliency, surveillance/archiving, and data governance/lineage** is not abstract — it is partly *court-ordered*. Knowing the specific public events lets you connect a data-platform conversation (auditability, immutability, DR, lineage, access control) to JPMC's own documented obligations and risk posture. Keep it factual; do not editorialize about the bank's conduct.

## Public tech incidents

- **Chase mobile/online-banking outage — April 24, 2024.** Chase.com and the Chase Mobile app went down for several hours; Downdetector peaked at ~7,000 reports (majority mobile-login). Chase attributed it to **"an internal issue"** and restored service the same day. **Fact.**[^1] *TAM relevance: availability/resiliency is the dominant operational theme; JPMC's framing ("internal issue") aligns with its public emphasis on multi-region active-active and DR.*
- **"Infinite money glitch" / check-fraud episode — late August 2024.** A viral TikTok trend exploited the deposit-to-clearing gap at Chase ATMs (classic check-kiting). Chase closed the loophole within days; on **2024-10-28 JPMorgan began filing civil suits** against high-dollar offenders (e.g., a Houston defendant owing ~$290,939). **Fact.**[^2] *TAM relevance: reinforces investment in real-time fraud controls — a long-running ML/AI win JPMC cites (see `jpmc-ai-strategy.md`).*

## Regulatory orders touching data / controls / records

- **SEC $125M + CFTC $75M = $200M — recordkeeping / "off-channel" (WhatsApp) communications — 2021-12-17.** JPMorgan Securities **admitted** firm-wide failures to preserve business communications sent via personal devices/text/WhatsApp/personal email (including by senior supervisors), violating Exchange Act §17(a) and Rule 17a-4. **Critically for tech spend:** both orders' undertakings explicitly require assessing the **technological solutions** for record retention, **electronic-communications surveillance**, and **archiving** of approved-channel messages. JPMC's own 8-K confirmed the resolution. **Fact** (re-verified against sec.gov and cftc.gov primary text).[^3][^4] *This was the first resolution in the industry-wide Wall Street off-channel sweep; it directly drives surveillance/archiving/records-retention technology investment.*
- **OCC $250M + Federal Reserve ~$98.2M = ~$348.2M — trade-surveillance / risk-data — 2024-03-14.** The OCC assessed $250M **plus a cease-and-desist order**; the Fed added ~$98.2M. Findings: gaps in trading-venue coverage and surveillance **"without adequate data controls,"** failing to surveil **billions of instances of trading activity across at least 30 global trading venues** (conduct 2014–2023). The OCC C&D ties remediation to the bank's **Data Risk Management Policy**, automated venue/trade-data reconciliation, and trading-venue/data governance; JPMC must obtain OCC non-objection before onboarding new venues. JPMC said it self-identified the issue and found no client/market harm. **Fact** (re-verified against occ.gov and federalreserve.gov primary text).[^5][^6] *This is fundamentally a **data-controls / risk-data-aggregation** story — it directly motivates data-governance, lineage, and reconciliation tooling.*
- **OCC $250M — fiduciary internal controls — 2020-11-24.** A **separate, earlier** matter: deficient internal controls over the fiduciary business. **Qualified.**[^7] **⚠️ Do not conflate the two $250M OCC penalties** — 2020 (fiduciary controls) vs 2024 (trade-surveillance data controls). They are different matters.

## The through-line: what each event funds

| Event | What it pressures JPMC to invest in |
|---|---|
| 2021 off-channel order | Communications **surveillance + archiving + records retention** (immutable/WORM-style, long retention, monitoring) |
| 2024 trade-surveillance order | **Data controls, lineage, reconciliation, data governance** across trading venues |
| 2024 outage | **Availability / resiliency** (multi-region active-active, DR/RTO/RPO) |
| 2024 fraud episode | **Real-time fraud detection / controls** (ML/AI) |

For a data-platform seller, the recurring requirements that fall out of these — **auditability, immutability/retention, encryption, lineage, granular access control, DR evidence** — are exactly what JPMC's Supplier Minimum Control Requirements demand (see `jpmc-vendor-procurement-and-mongodb-signals.md`).

## Cross-reference note (regulation generalities)

The *regulations* behind these orders — Exchange Act §17(a)/Rule 17a-4, FINRA recordkeeping, BCBS 239 risk-data-aggregation principles, SR 11-7 model risk, FFIEC examination culture — are explained generically in **`fsi-banking-regulatory-context`**. This reference covers **JPMC's specific orders and what they fund**, not the regulatory background. When you need "what is BCBS 239 / 17a-4 / the FFIEC handbook," route there.

## Sources

1. PIX11 / AZCentral / Mirror, Chase outage coverage (2024-04-24/25) — outage scope; Downdetector ~7,000; "internal issue." (Press) — https://pix11.com/news/chase-bank-app-website-outage-affects-thousands/
2. CNBC, "JPMorgan suing customers over 'infinite money glitch'" (2024-10-28) — episode + lawsuits. (Press) — https://www.cnbc.com/2024/10/28/jpmorgan-suing-customers-over-infinite-money-glitch.html
3. SEC press release 2021-262 (2021-12-17) — $125M; JPMS admitted; surveillance/archiving/tech-remediation undertakings. (SEC, primary) — https://www.sec.gov/newsroom/press-releases/2021-262
4. CFTC press release 8470-21 (2021-12-17) — $75M; archiving/monitoring undertakings. (CFTC, primary) — https://www.cftc.gov/PressRoom/PressReleases/8470-21
5. OCC news release nr-occ-2024-25 (2024-03-14) — $250M + cease-and-desist; surveillance "without adequate data controls"; ≥30 venues; Data Risk Management Policy remediation. (OCC, primary) — https://www.occ.gov/news-issuances/news-releases/2024/nr-occ-2024-25.html
6. Federal Reserve enforcement20240314a (2024-03-14) — ~$98.2M; ~$348.2M combined. (Federal Reserve, primary) — https://www.federalreserve.gov/newsevents/pressreleases/enforcement20240314a.htm
7. OCC news release nr-occ-2020-159 (2020-11-24) — earlier, distinct $250M fiduciary-controls penalty. (OCC, primary) — https://www.occ.gov/news-issuances/news-releases/2020/nr-occ-2020-159.html
