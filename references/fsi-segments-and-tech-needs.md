<!-- Provenance: reference under the `fsi-banking-regulatory-context` skill. Created 2026-06-18 via /dr deep research (firecrawl/exa + analyst/vendor/regulator sources). Educational seller orientation, NOT advice. -->

# FSI Segments & Differentiated Technology Needs

`verified-as-of: 2026-06-18` (vendor landscape and adoption figures are fast-moving — re-verify before quoting to a customer).

## Contents

- [Orientation: FSI is not one buyer](#orientation-fsi-is-not-one-buyer)
- [Core banking systems — the central fault line](#core-banking-systems--the-central-fault-line)
- [1. Retail / consumer banking](#1-retail--consumer-banking)
- [2. Commercial / corporate banking](#2-commercial--corporate-banking)
- [3. Investment banking & capital markets](#3-investment-banking--capital-markets)
- [4. Payments](#4-payments)
- [5. Wealth & asset management](#5-wealth--asset-management)
- [6. Insurance (adjacent segment)](#6-insurance-adjacent-segment)
- [Cross-segment patterns (seller takeaways)](#cross-segment-patterns-seller-takeaways)
- [Sources](#sources)

## Orientation: FSI is not one buyer

Financial services is a set of sub-verticals with materially different workloads, data shapes, latency tolerances, and buying cycles. Analysts consistently slice it into **banking** (retail/consumer + commercial/corporate), **capital markets / investment banking**, **payments**, and **wealth & asset management**, with **insurance** treated as an adjacent vertical.[^1][^2] As of 2026, payments is the largest fintech-investment vertical while capital-markets-tech and insurtech grow off a smaller base.[^2] The single most important technical fault line cutting across banking is the **core banking system**: whether it is a legacy mainframe core or a modern cloud-native one. *(Confidence: high — multi-analyst convergence on the taxonomy.)*

**Seller framing.** Don't sell to "a bank" generically. Identify the *segment* and the *workload*: a microsecond-latency trading store, a strong-consistency ledger, a high-throughput payment hub, and a wealth-data aggregation layer are genuinely different data problems that call for different answers.

## Core banking systems — the central fault line

**Concept.** The **core banking system (CBS)** is the system of record/ledger that processes deposits, loans, and payments. As of 2026 the market splits into two camps *(fast-moving vendor landscape — verify any vendor claim before quoting)*:[^7][^45]

- **Legacy / incumbent (modernized) cores:** Fiserv, FIS, Jack Henry, Temenos, Oracle FLEXCUBE, Infosys Finacle, TCS BaNCS, SAP Fioneer. Heavy compliance install base; the tier-one replacement cycle.[^3][^7][^45]
- **Cloud-native / next-gen cores:** Thought Machine (Vault), Mambu, Finxact, Tuum, 10x Banking, Pismo, nCino, Vodeno. Winning greenfield and digital-bank launches.[^6][^7][^45]

**The hard, well-evidenced truth: core modernization frequently disappoints.** IBM's 2025 study found **fewer than half of banking CIOs reported meaningful gains** from core modernization, and a majority struggled with update costs, risk management, and resilience *after* migration (e.g., ~73% found system-update costs harder to manage post-migration, citing vendor lock-in).[^46] McKinsey reports **only ~30% of core transformations** successfully migrated ledgers/products over the past decade.[^9] EY and others now push **progressive / parallel-core / strangler-fig modernization over "rip and replace."**[^8][^9][^47] *(Confidence: high — IBM, McKinsey, EY independently converge.)*

> **Seller implication.** Your opportunity is usually the **satellite/digital/data layer around the core**, not the ledger itself. Expect a risk-aware, scarred buyer who scrutinizes migration risk, lock-in, and resilience, not just features. The "composable core" pattern (keep/buy a commodity ledger, build differentiation around it) is the prevailing direction.

## 1. Retail / consumer banking

*(Business/technology side — not consumer personal-finance advice, which lives in `consumer-finance` / `personal-banking`.)*

- **Dominant workloads.** The CBS is the gravitational center; around it sit digital/mobile channels, card management, payments, onboarding/KYC, and fraud/risk.[^3][^4][^5]
- **Data shapes.** High-volume, **strong-consistency transactional (OLTP)** ledger workloads (ACID balances/entries), increasingly paired with **real-time event streams** (event-driven architectures) and a 360-customer-view analytical layer for personalization.[^4][^6]
- **Modernization pressure.** Moving off 1970s-80s monolithic/COBOL cores toward real-time, API-first, modular/cloud-native architectures, driven by faster product launch, personalization, embedded finance / Banking-as-a-Service, and open-banking APIs.[^3][^4][^6][^7]
- **Buying considerations.** Replacement is rare and high-stakes (18-36 months minimum for a modern core). Buyers weigh scalability, back-office integration, security, regulatory configurability, vendor track record, and migration risk above raw feature count.[^3][^8][^9]
- **Neobank caveat (well-evidenced).** Digital-only retail banks face structurally harder economics: the ECB found digital banks remain *less* profitable than traditional banks (higher deposit costs, high fixed IT/marketing spend), and a large majority of neobanks globally remain unprofitable as of early 2026; profitability correlates with holding a banking charter + lending products, not interchange alone.[^10][^11][^12] *(Confidence: high — regulator + multiple analyses converge.)*

## 2. Commercial / corporate banking

- **Dominant workloads.** Treasury/cash management, commercial lending, trade finance, and transaction banking, increasingly delivered as **embedded, real-time services inside the corporate client's own workflows** (liquidity analytics, automated treasury).[^13][^14]
- **Data shapes.** Lower transaction *count* than retail but higher *value* and far higher *complexity per transaction*: multi-entity hierarchies, multi-currency, entitlements, rich remittance/reference data (a key driver of ISO 20022 adoption). Corporate clients demand real-time visibility into liquidity/risk across fragmented systems.[^13]
- **Modernization pressure.** Legacy cores and fragmented data architectures limit scalability and speed; banks re-platform onboarding/servicing and add low-code configuration for phased modernization.[^13][^14] *(Confidence: medium-high — thinner independent coverage than retail/capital-markets.)*
- **Buying considerations.** Real-time liquidity/data integration, API connectivity to client ERPs/TMSs, and incremental modernization without disrupting high-value corporate relationships.[^13][^14]

## 3. Investment banking & capital markets

This is the segment whose technology needs diverge **most** from the rest of FSI.

- **Dominant workloads.** Order/execution management (OMS/EMS), market-data ingestion/distribution, **ultra-low-latency** trading and direct market access (DMA), pre-/post-trade risk, and post-trade clearing/settlement.[^15][^16][^17][^18] Latency requirements are extreme: co-located DMA gateways use **FPGA hardware acceleration**, and vendors advertise sub-200-microsecond round-trips and single-digit-microsecond time-series ingestion.[^17][^19][^20] *(These specific latency figures are single-vendor illustrative claims; treat as directional.)*
- **Data shapes.** High-frequency **time-series** data (ticks, quotes, order-book snapshots, trades) at **petabyte scale**, with a characteristic hot/warm/cold tiering (in-memory real-time → recent intraday on disk → historical archive for backtesting/compliance, the classic kdb+ RDB/IDB/HDB model). Specialized columnar/vector time-series engines dominate (kdb+/KX, ClickHouse, QuestDB, Deltix).[^19][^21][^22][^23] *(Confidence: high — multiple independent vendor + customer sources on the time-series/petabyte profile.)*
- **Modernization pressures.**
  - **T+1 settlement** (US/Canada/Mexico, May 2024) compressed post-trade timelines, forcing straight-through processing (STP) over manual workflows; DTCC reported same-day affirmation rising from ~73% to ~95% post-T+1.[^24][^25]
  - **Cloud adoption is real but uneven and contested.** The latency-sensitive front office (matching, co-located trading) largely stays on-prem/co-lo; research, backtesting, risk, and historical analytics move to cloud. Barriers (latency determinism, data sovereignty, migration complexity, cost) are genuine: "cloud networks were not originally designed for capital markets."[^21][^26] *(Confidence: high — barriers corroborated by analyst + specialist sources.)*
- **Buying considerations.** Latency determinism, co-location, FIX/native connectivity, time-series performance, and regulatory-deadline reliability outrank cost. This segment buys *specialized* infrastructure, not general-purpose platforms.

## 4. Payments

- **Dominant workloads.** Authorization/clearing/settlement across an expanding set of **rails** (card networks (Visa/Mastercard), ACH, wires (Fedwire/CHIPS), and **real-time rails (FedNow, The Clearing House RTP)**), plus issuing/acquiring processing, fraud prevention, and payment orchestration/hubs. The ecosystem splits into roles: cardholder → merchant → gateway → processor → acquirer → card network → issuer.[^27][^28][^30]
- **Data shapes.** Very high throughput with hard availability/latency demands, increasingly **bifurcated into batch (ACH) vs. always-on real-time (instant payments)**, often architected as separate execution engines behind a common payment hub.[^31] Scale (2025): ACH handled tens of billions of payments / tens of trillions of dollars; instant-payment rails are growing fast but off a smaller base.[^27]
- **Modernization pressures.**
  - **ISO 20022 migration** is the dominant theme. Fedwire cut over July 2025; FedNow and RTP launched on ISO 20022; SWIFT's cross-border MT→ISO coexistence window closed **November 2025**; downstream bank work continues into 2026-27.[^29][^32] A recognized risk: banks "check the compliance box" with translation layers without true modernization, adding fragile systems.[^32]
  - **Multi-rail complexity:** many instant-payment-enabled banks run *both* FedNow and RTP, duplicating fraud/compliance/exception stacks, driving demand for unified platforms.[^27]
- **Disconfirming note.** Despite the rails existing, US instant-payment adoption is *slower than headlines suggest*: a large share of institutions had not implemented sending (many were "receive-only"), held back by ROI at low volume, fraud risk (irreversible settlement), and the cost of *sending* over legacy cores.[^35] *(Confidence: high — multiple trade-press converge.)*
- **Buying considerations.** Throughput + 24/7/365 availability for instant rails, ISO 20022 data handling, multi-rail orchestration, real-time fraud, and decoupling channels from execution engines.[^27][^31]

## 5. Wealth & asset management

- **Dominant workloads.** Portfolio management/accounting, trading/rebalancing, custody integration, planning, CRM, billing, performance reporting, compliance, consolidating into **integrated advisor platforms** with a unifying data layer; robo-advisors automate portfolio management/rebalancing.[^36][^37][^38]
- **Data shapes.** The defining problem is **data aggregation/consolidation**: pulling custodial data, planning inputs, and third-party feeds into a "single source of truth." Volumes are modest vs. capital markets; the challenge is *integration and reconciliation across fragmented systems*. AI-assisted advisor copilots are an emerging 2026 workload.[^36][^37][^39] *(Confidence: medium-high — heavily vendor-sourced; capability claims are factual but not independently benchmarked.)*
- **Modernization pressure.** Replacing stitched-together point tools with consolidated platforms; data unification "without rebuilding the tech stack"; embedding AI into advisor workflows.[^36][^41]
- **Buying considerations.** Breadth of integrations (custodians, CRM, planning), data-layer quality, advisor UX, and adding AI without re-platforming.[^36][^37]

## 6. Insurance (adjacent segment)

Insurance is **adjacent**, not core banking, but technologically rhymes with the CBS modernization story.

- **Dominant workloads.** Policy administration systems (PAS), claims, underwriting, billing: the insurance "core." Guidewire and Duck Creek are the recognized P&C platform names; cloud-native PAS is the modernization target.[^42][^43][^44]
- **Data shapes & modernization.** Like banking, insurance runs on aging cores (trade sources estimate a large majority of mid-tier P&C carriers operate platforms built 15-25 years ago — *directional, not audited*), now being replaced by cloud-native PAS, with AI reshaping underwriting/claims/fraud.[^43]
- **Buying considerations.** Mirror banking: scalability, cloud-native vs. legacy, modular policy/billing/claims, AI-driven underwriting/claims automation.[^42][^43]

## Cross-segment patterns (seller takeaways)

1. **Legacy-core gravity is the universal modernization story.** Retail, commercial, and insurance all orbit an aging system-of-record and are all shifting from rip-and-replace toward incremental/parallel modernization because full replacement fails most of the time.[^9][^46][^47][^43]
2. **The real-time vs. batch split is now architectural.** Every segment is bifurcating into always-on real-time engines alongside legacy batch, often behind orchestration/hub layers.[^31][^27][^24]
3. **API-first, modular, event-driven, cloud-native** is the converged target-architecture vocabulary across banking, payments, and insurance.[^3][^4][^6][^31]
4. **Data shape dictates the datastore.** Strong-consistency OLTP ledgers (banking), petabyte time-series (capital markets), high-throughput multi-rail transactional (payments), aggregation/single-source-of-truth (wealth) are genuinely different problems.[^4][^22][^27][^36]
5. **Latency tolerance is the great differentiator:** capital-markets front-office (microseconds, FPGA, co-lo, on-prem) vs. wealth aggregation (batch-tolerant). This explains why capital markets resists cloud where others embrace it.[^26][^21][^19]
6. **Standardization waves force synchronized spend:** ISO 20022 (payments) and T+1 (capital markets) are regulator/standards-driven deadlines that compelled industry-wide simultaneous technology work.[^29][^24]
7. **"Modernized" ≠ "better outcomes":** buyers are increasingly skeptical and risk-aware after a decade of disappointing core programs; expect scrutiny on migration risk, lock-in, and resilience.[^46][^9]

## Sources

[^1]: https://www.mckinsey.com/industries/financial-services/our-insights — McKinsey Financial Services Insights (analyst). Segment taxonomy authority.
[^2]: https://www.deloitte.com/us/en/insights/industry/financial-services/financial-services-industry-outlooks/banking-industry-outlook.html — Deloitte 2026 Banking & Capital Markets Outlook (analyst).
[^3]: https://www.fisglobal.com/products/fis-modern-banking-platform — FIS Modern Banking Platform (vendor; factual product scope).
[^4]: https://www.intellectdesign.com/solutions/core-banking-solutions/consumer-core-banking/ — Intellect Consumer Core Banking (vendor; retail core workloads).
[^5]: https://www.santander.com/en/press-room/press-releases/2025/06 — Santander Gravity core migration (bank press release; self-reported figures, illustrative).
[^6]: https://mambu.com/en/insights — Mambu insights (vendor; cloud-native composable core).
[^7]: https://www.everestgrp.com/blog/ — Everest Group core-banking landscape 2026 (analyst/trade; vendor taxonomy).
[^8]: https://www.ey.com/en_nz/insights/financial-services/emeia/how-to-avoid-common-pitfalls-when-modernizing-banking-systems — EY: pitfalls modernizing core banking (analyst).
[^9]: https://www.mckinsey.com/capabilities/tech-and-ai/our-insights/tech-forward/how-to-get-a-core-banking-transformation-right-eight-mistakes-to-avoid — McKinsey: 8 core-transformation mistakes (analyst; ~30% success rate).
[^10]: https://www.ecb.europa.eu/press/financial-stability-publications/fsr/focus/2025/html/ecb.fsrbox202505_04~17b39a3c1a.en.html — ECB: digital-bank vs traditional profitability (regulator).
[^11]: https://www.fintechessential.com/industry/reports/neobank-profitability-problem — The Neobank Profitability Problem (trade-press).
[^12]: https://www.spark.money/research/neobank-business-model-analysis — Neobank Business Models (trade-press).
[^13]: https://bankingblog.accenture.com/commercial-banking-top-trends-2025 — Accenture commercial-banking trends (analyst).
[^14]: https://www.bottomline.com/resources — Bottomline: commercial-banking innovation (vendor).
[^15]: https://flextrade.com — FlexTrade FlexONE OEMS (vendor; OMS/EMS architecture).
[^16]: https://www.tcs.com/what-we-do/products-platforms/tcs-bancs — TCS BaNCS Capital Markets (vendor; post-trade/custody).
[^17]: https://www.raptorfintech.com — Raptor DMA / FPGA gateway (vendor; ultra-low-latency).
[^18]: https://iongroup.com/products/markets/xtp/ — ION XTP post-trade (vendor; clearing/settlement/STP).
[^19]: https://questdb.com/capital-markets/ — QuestDB for Capital Markets (vendor; sub-200µs latency claim).
[^20]: https://deltixlab.com — Deltix TimeBase (vendor; microsecond time-series claims).
[^21]: https://clickhouse.com/blog/qrt — ClickHouse / QRT petabyte platform (vendor + named customer).
[^22]: https://kx.com/products/kdb/ — kdb+ time-series DB (vendor; RDB/IDB/HDB tiering model).
[^23]: https://kx.com/products/kdb-x/ — KDB-X (vendor; tick-data workloads).
[^24]: https://posttrade360.com — Real cost of T+1 / DTCC affirmation stats (trade-press; DTCC 73%→95%).
[^25]: https://www.swift.com/securities/preparing-t1-settlement — SWIFT: preparing for T+1 (standards body).
[^26]: https://www.greenwich.com/market-structure-technology — Coalition Greenwich: capital-markets cloud adoption (analyst; cloud-limit barriers).
[^27]: https://investor.aciworldwide.com/news-releases — ACI multi-rail platform (vendor; cites network-operator volumes).
[^28]: https://stripe.com/us/resources/more/the-payment-industry-ecosystem-explained — Stripe: payment ecosystem (vendor explainer; roles).
[^29]: https://www.americanbanker.com/payments/news/fedwire-migrates-to-the-iso-20022-messaging-standard — American Banker: Fedwire ISO 20022 (trade-press).
[^30]: https://clearlypayments.com/blog/payment-systems-architecture-2026/ — Payment Systems Architecture 2026 (trade/vendor; stack layering).
[^31]: https://www.fisglobal.com — FIS Open Payment Framework (vendor; batch-vs-realtime engine separation).
[^32]: https://fintekcafe.com/iso-20022-migration-bank-tech-stack-impact/ — ISO 20022 reshaping bank tech stacks (trade-press; coexistence-window close).
[^35]: https://www.americanbanker.com/opinion/fraud-remains-a-major-concern-for-potential-fednow-participants — FedNow/RTP adoption barriers (trade-press; receive-only/fraud).
[^36]: https://orion.com — Orion / Denali Data Layer (vendor; advisor platform + data aggregation).
[^37]: https://www.fisglobal.com/products/fis-unity-wealth-platform — FIS Unity Wealth (vendor; 360° household view).
[^38]: https://betterment.com/advisors — Betterment Advisor Solutions (vendor; robo-advisor).
[^39]: https://assetvantage.com — Asset Vantage family-office software (vendor; single-GL multi-entity).
[^41]: https://qdeck.ai — Qdeck.ai (vendor; AI advisor copilot, emerging 2026 workload).
[^42]: https://insurtechdigital.com/top10/top-10-policy-administration-systems — Top 10 PAS (trade-press; PAS landscape).
[^43]: https://www.decerto.com/us/post/insurance-software-development-complete-2026-guide-for-insurers — Insurance Software 2026 Guide (trade/vendor; legacy-core estimate is directional).
[^44]: https://insurtechdigital.com — Guidewire vs Duck Creek comparison (trade-press; cloud-native PAS).
[^45]: https://dashdevs.com/blog/top-10-core-banking-solutions/ — Top Core Banking Solutions 2026 (trade-press; legacy-vs-cloud-native taxonomy).
[^46]: https://www.ibm.com/thought-leadership/institute-business-value/en-us/report/core-banking-modernization — IBM IBV: core-banking modernization study (analyst/vendor-research; <50% CIO gains).
[^47]: https://www.prosightfa.org/insights — Progressive Modernization (trade-press; strangler-fig/parallel-core).
