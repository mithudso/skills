<!-- Provenance: reference under the `big-bank-IT` skill. Created 2026-06-18 via /dr deep research (analyst + vendor + gov + primary architecture sources). Educational, vendor-neutral technical orientation — NOT advice. Disconfirming evidence on core-replacement and next-gen-core adoption is preserved. -->

# Core Banking Systems & the System-of-Record vs System-of-Engagement Split

`verified-as-of: 2026-06-18` (vendor landscape, COBOL/mainframe statistics, and core-transformation figures are fast-moving and largely vendor-origin — re-verify and present as attributed ranges before quoting to a customer).

> **Educational, vendor-neutral orientation, NOT advice.** The headline mainframe/COBOL numbers are *directionally* solid but *numerically* soft and mostly IBM/Micro-Focus/Reuters-origin (they recirculate through trade press, so multiple "sources" often collapse to one). Present them as ranges, attributed and dated. WHY banks buy conservatively → `fsi-banking-regulatory-context`; modernization *patterns* (the Rs, strangler-fig) → `batch-and-mainframe-modernization.md`; MongoDB product fit → the `mongodb-*` hubs.

## Contents

- [What a core banking system actually is and does](#what-a-core-banking-system-actually-is-and-does)
- [The mainframe estate reality](#the-mainframe-estate-reality)
- [COBOL and the legacy-maintenance budget](#cobol-and-the-legacy-maintenance-budget)
- [The packaged core-banking vendor landscape](#the-packaged-core-banking-vendor-landscape)
- [System-of-record vs system-of-engagement (SoR/SoE)](#system-of-record-vs-system-of-engagement-sorsoe)
- [Why core replacement is rare and risky](#why-core-replacement-is-rare-and-risky)
- [Seller takeaways](#seller-takeaways)
- [Sources](#sources)

## What a core banking system actually is and does

A **core banking system (CBS)** is the back-end software holding a bank's **ledger of record** — the authoritative source of truth for every balance, account, and posted transaction. Across vendor-neutral definitions the consistent functional decomposition is: **account master/management** (open/maintain/close deposit and loan accounts, store customer/KYC data); **transaction processing/posting** (validate and post every debit/credit, keeping balances consistent across channels); **general ledger** (double-entry bookkeeping producing the authoritative financials); **deposits** (interest accrual on savings/CDs); **loans & credit** (origination, disbursement, interest, repayment, collections); and **compliance/reporting** (regulatory reports, audit trails, reconciliation, AML feeds).[^1][^2][^3] *(Confidence: fact — multiple independent sources converge on this module list.)*

A notable design trend is the **"lean core" / "thin core"** model: the core handles only the most fundamental functions (general ledger + transaction posting) and everything else (payments, KYC, product management, channels) is handled by specialist applications connected via APIs.[^1] This is the same impulse as the SoR/SoE split below. *(Confidence: qualified.)*

## The mainframe estate reality

Banks' cores overwhelmingly run on **IBM Z** mainframes (z/Architecture hardware running **z/OS**), with mission-critical business logic in **COBOL** (often alongside PL/I), transaction processing via **CICS**, data in **DB2**, and batch orchestration via **JCL**.[^4] *(Confidence: fact.)*

The widely-repeated **"~45 of the top 50 banks still run mainframe"** claim is real but the exact number drifts by source and date — treat it as a range. IBM's Institute for Business Value cites **43 of the top 50** (ranked by S&P, Nov 2022); a 2026 IBM article cites **44 of the top 50**; an integrator (DXC) cites **45 of the top 50**; older figures cite "92 of the top 100."[^4][^5][^6][^7][^8] Mainframes are credited with **~68–72% of the world's transactional workloads** and **~90% of card transactions**, with availability claims up to "eight-nines."[^4][^5][^6] *(Confidence: qualified — repeated across many sources but largely IBM-origin; treat the direction as high-confidence and the precise numbers as soft.)*

**Why banks stay:** transactional integrity, pervasive encryption, raw throughput, and decades of hardened business logic. The deterrent to moving is cost and risk, not technical impossibility.[^8][^9] *(Confidence: fact.)*

## COBOL and the legacy-maintenance budget

**COBOL lines-of-code estimates vary widely — cite a range, not a point.** The most-repeated figure is **~220 billion lines in production** (a Reuters/Gartner-lineage estimate); the Open Mainframe Project estimated ~250B and Micro Focus/OpenText ~800B in separate surveys. The credible band is **~220B–800B**, the spread driven by sampling method and whether distributed COBOL is counted.[^9][^10][^11][^12][^13] Pervasiveness stats (COBOL touches ~95% of ATM transactions, ~$3T daily commerce) recirculate from a few vendor/survey origins — directionally credible, numerically soft.[^9][^10][^11] *(Confidence: qualified.)*

**The "~70% of IT budget on legacy maintenance" claim is well-supported by multiple independent analysts.** Accenture's 2026 banking report states banks spend **~70%** of IT budget keeping old technology alive (and that tech costs grew 4× faster than revenue over 15 years); RS2 independently reports 70%; Backbase cites 70–80%.[^14][^15][^16][^17] This maps to the industry "**Run-the-Bank vs Change-the-Bank (RTB:CTB)**" framing (~70:30 average; top-quartile target ~50:50). McKinsey notes banking spends 6–12% of revenue on IT — more than any other major industry.[^18][^19] *(Confidence: fact — multiple independent analysts converge on ~70%.)*

**COBOL-talent retirement risk is real and consistently quantified:** the most-cited baseline is an average COBOL-developer age of **~55–58 with ~10% retiring annually**; a late-2025 survey found **71% of mainframe teams understaffed and 54% underfunded**.[^20][^21][^22] *(Confidence: fact.)*

## The packaged core-banking vendor landscape

The market is widely described as **bifurcated** into legacy/incumbent cores and cloud-native/next-gen cores.[^23][^24][^25][^26] *(Confidence: fact.)*

**Legacy / incumbent cores:**
- **US-dominant "Big Three":** **Fiserv** (DNA, Premier, Signature; dominant in US community banks/credit unions), **FIS** (Profile, IBS, **Systematics** — the platform many top-50 banks run on IBM Z), **Jack Henry** (SilverLake; strong in credit unions/regional banks). These three dominate US core processing.[^24][^27][^28] *(Confidence: fact.)*
- **Global Tier-1/2:** **Temenos** (Transact/T24; frequently cited global leader by customer count), **Oracle FLEXCUBE**, **Infosys Finacle**, **TCS BaNCS**, **Finastra**.[^23][^24][^26] *(Confidence: fact for the set; qualified for exact market-share percentages, which vary by analyst and segment.)*

**Cloud-native / next-gen cores** (built post-2010; microservices, API-first, real-time, products-as-configuration/code):
- **Thought Machine (Vault Core)** — products defined as Python **"smart contracts"**; real-time ledger; multi-cloud; customers include Lloyds, Standard Chartered, JPMorgan Chase.[^25][^26][^31]
- **Mambu** — SaaS-only composable core; strong in lending/savings (N26, ABN AMRO).[^23][^26]
- **Finxact** — real-time event-driven core, **acquired by Fiserv (2022)**.[^26][^27]
- **10x Banking (SuperCore)** — targets Tier-1; runs at Chase UK and Westpac.[^25][^26]
- **Pismo** — Latin-American cloud core, **acquired by Visa**.[^26] *(Confidence: fact for the vendor set and customers.)*

## System-of-record vs system-of-engagement (SoR/SoE)

The **SoR/SoE distinction originates with Geoffrey Moore** (of *Crossing the Chasm*) in the AIIM white paper *"Systems of Engagement and the Future of Enterprise IT"* (**published January 2011**). Moore framed **Systems of Record** as data-centric, operational, reliable, secure, tightly-coupled, stateful; **Systems of Engagement** as user-experience-centric, provisional, reversible, loosely-coupled, stateless.[^34][^35][^36] *(Confidence: fact — primary source.)*

**Applied to banking:** the **core banking system is the SoR** — the immutable, strong-consistency ledger of record — while the digital/channel/data layer (mobile and web apps, onboarding, servicing, orchestration) is the **SoE** built *around* it. IBM's own banking reference architecture explicitly names a "Core systems and systems of record" layer (general ledger, payment processing, customer master — "often run on mainframe") beneath the channels/engagement layers; McKinsey's stack similarly separates a core-systems/data layer, a decisioning layer, and an engagement layer.[^3][^37][^39] *(Confidence: fact — multiple independent sources incl. primary vendor architecture docs.)*

**The immutability / ledger-of-record property** is the core's defining characteristic: an append-only, procedurally/cryptographically verifiable, strong-consistency ledger (AWS's QLDB core-banking reference frames the ledger as "the immutable system-of-record," separating reads from writes).[^40] *(Confidence: fact.)*

**Why banks wrap rather than replace** — the **"composable core" / "hollow the core" / "safe wrapper"** patterns: keep the core as the stable ledger and build modern services around it. McKinsey's **"progressive modernization"** (retain the legacy platform but progressively minimize it, using a strangler pattern) is described as the strategy "most banks have pursued" as the lower-risk option.[^33][^41] The recurring framing: **"the core remains the system of record; the channel/digital layer becomes the system of execution."**[^38] *(Confidence: fact / qualified.)*

**The data-layer-around-the-core opportunity:** the highest-value, lowest-risk move is often building a digital/orchestration/data interface layer (golden-source data access, streaming/CDC, data mesh, API gateway) rather than ripping out the ledger. McKinsey reported one bank cut a data-modernization program's cost by two-thirds and timeline from 5→3 years by focusing on **access to golden-source data via interface layers** instead of centralizing all data.[^41] *(Confidence: fact.)* This is the conceptual seam where document/operational data stores (an Operational Data Layer in front of the core) sit — see `data-platform-and-integration.md`.

## Why core replacement is rare and risky

**Cost and timeline:** McKinsey's analysis of 50+ CBS transformations gives the canonical figures — integration/replacement for a mid-size bank can exceed **$50M**; for larger banks **$300–400M is "not unheard of"**; planning alone can take up to 3 years and implementation 4–10 years, with observed **~100% cost overruns and 50–100% timeline overruns**.[^41][^44] The "$100–500M, 3–7 year" framing is well-supported as a central estimate. *(Confidence: fact.)*

**The "~30% succeeded" claim is McKinsey's, and consistent across a decade of their research:** *"over the past decade, only about 30 percent of CBS transformations succeeded in carrying out a complete migration of ledgers and products to a new system."*[^44] *(Confidence: fact — but note it is one institutional source repeated, not three independent houses; weight accordingly.)*

**The IBM 2025 finding is confirmed.** IBM's IBV report **"The 94% core banking problem"** (Sept 2025, with BIAN; ~700 leaders at banks >$10B assets) found **94% of core-modernization projects exceed timelines** and **"less than half of banking CIOs reported meaningful gains"** (improvements reported by only 49% client experience, 46% data access, 37% resilience, 31% risk management, 27% update cost); in some cases operational capabilities *deteriorated*.[^46][^47][^49] *(Confidence: fact — IBM primary source.)*

### Disconfirming evidence (preserve these)

1. **COBOL/mainframe permanence is genuinely contested as of 2026.** On Feb 23–24 2026, **IBM shares fell ~13% (~$40B)** after Anthropic published Claude Code tooling positioned to translate/modernize COBOL — the market briefly priced the moat as evaporating.[^50][^51] *…But the more defensible position is the rebuttal:* Gartner's Matt Brasier — *"Modernizing COBOL has been a technically solved problem for a while. The real problem is that the costs of modernization are high and the ROI is low"*; **translation ≠ modernization** (the hard parts are data-architecture redesign, runtime replacement, transaction integrity, testing, regulatory sign-off).[^50][^51] Gartner predicts **75% of AI-powered mainframe-exit vendors will pivot or cease by 2030.**[^53]
2. **Next-gen cloud-native cores have limited Tier-1 *full-production* track records.** A 2026 buyer's guide warns "no Tier-1 bank in full production" on some next-gen cores and that "a pilot with 10,000 accounts is fundamentally different from 10 million."[^25] Counter-evidence that adoption *is* progressing: Lloyds/Standard Chartered/JPMorgan on Thought Machine, Chase UK/Westpac on 10x.[^25][^31][^33]
3. **The stat ecosystem is circular.** The "43/44/45 of top 50," "~70% of workloads," and "220B lines" figures trace back to a few IBM/Micro-Focus/Reuters origins. Treat the direction as high-confidence, precise numbers as qualified.

## Seller takeaways

- **Don't sell to "a bank" — sell to a layer.** SoR (the ledger): don't touch. SoE/data layer: the real opportunity. Frame around a *specific workload's* data problem.
- **The wrap-don't-replace logic is your strongest, most-defensible talking point.** The core works as a ledger; the risk/ROI math (McKinsey ~30%, $50–400M, 3–10yr; IBM 2025 <half see gains) is *why* banks build the engagement/data layer around it — which is precisely where new architecture and new datastores land.
- **Meet a scarred buyer.** Expect scrutiny of migration risk, vendor lock-in, and post-migration resilience over feature count.

## Sources

[^1]: What is Core Banking? — https://www.opendue.com/glossary/core-banking — trade — core functions; lean/thin core.
[^2]: Core Banking System glossary — https://www.spark.money/glossary/core-banking-system — trade — SoR/single-source-of-truth; six modules.
[^3]: Evolution of banking systems of engagement — https://plumery.com/the-evolution-of-banking-systems-of-engagement/ — trade/vendor — SoR vs SoE in banking; attributes Moore 2011.
[^4]: How mainframes are rewriting the AI playbook (IBM) — https://www.ibm.com/downloads/documents/us-en/1227c12d3a38b190 — vendor — 43 of top 50 banks; ~72% transactional workloads; 90% card transactions.
[^5]: Unlocking the full potential of mainframes (IBM IBV) — https://www.ibm.com/downloads/documents/us-en/10c31775c85402a2 — analyst/vendor — 43 of top 50; 70% transactional workloads; eight-nines availability.
[^6]: The world still runs on mainframes (IBM Think) — https://www.ibm.com/think/news/world-runs-on-mainframes — vendor — 44 of top 50 banks (2026 framing); IBM Z definition.
[^7]: Mainframe 2020 (MIT Tech Review, sponsored) — https://wp.technologyreview.com/wp-content/uploads/2020/09/Mainframe-2020-A-catalyst-for-transformation_0920.pdf — trade/sponsored — 44 of top 50, 92 of top 100 (older figures).
[^8]: Why banks still rely on COBOL-driven mainframes (DXC) — https://dxc.com/insights/knowledge-base/blogs/why-banks-still-rely-on-cobol-driven-mainframe-systems — vendor/integrator — "45 of top 50"; >40% of banks use COBOL as core tech.
[^9]: Why Mainframes Remain Vital in BFS (Global Banking & Finance) — https://www.globalbankingandfinance.com/the-outlook-for-bank-mainframes-in-2022-and-beyond/ — trade — >200B lines COBOL; transactional-integrity rationale.
[^10]: Legacy code modernization case study — https://case-studies.ai/use-cases/engineering-and-research/ER-002-legacy-code-modernization/ — trade — 220–344B COBOL lines; CBA AU$1B/5yr; 3–7yr timelines.
[^11]: Why Banks Still Use COBOL (DataField.Dev) — https://datafield.dev/blog/why-banks-still-use-cobol.html — trade — 220B lines; 95% ATM / $3T daily.
[^12]: COBOL the incessant number cruncher (Planet Mainframe) — https://planetmainframe.com/2024/03/cobol-the-incessant-number-cruncher/ — trade — 250B vs 800B vs 220B estimates.
[^13]: 800 billion lines of COBOL (The Stack) — https://www.thestack.technology/cobol-in-daily-use/ — trade — Micro Focus survey: 800B lines.
[^14]: Banks spending 70% of IT budgets on legacy tech (Accenture, via Reseller News) — https://www.reseller.co.nz/article/4142040/banks-spending-70-of-it-budgets-on-legacy-tech-accenture.html — analyst (Accenture 2026) — 70% on legacy; tech costs grew 4× faster than revenue.
[^15]: Banks waste 70% of IT budgets on legacy tech (bobsguide / RS2) — https://www.bobsguide.com/banks-waste-70-of-it-budgets-on-legacy-tech-is-ai-the-answer/ — trade (RS2) — corroborates 70%.
[^16]: 70% of bank IT budgets go to maintaining legacy tech (Digit/RS2) — https://www.digit.fyi/70-of-bank-it-budgets-go-to-maintaining-legacy-tech/ — trade (RS2) — corroborates 70%.
[^17]: Banking legacy technical debt (Backbase) — https://www.backbase.com/blog/modernizing-legacy-tech-building-platforms-that-grow-with-your-bank — vendor — 70–80% of IT budget on legacy.
[^18]: Change-the-Bank vs Run-the-Bank (KnowMBA) — https://knowmba.com/concepts/change-the-bank-vs-run-the-bank — trade — RTB:CTB ~70:30 vs ~50:50 top quartile.
[^19]: Managing bank IT spending (McKinsey) — https://www.mckinsey.com/capabilities/mckinsey-digital/our-insights/tech-forward/managing-bank-it-spending-five-questions-for-tech-leaders — analyst — banking IT = 6–12% of revenue (highest of any industry).
[^20]: COBOL Developer Retirement vs Training Pipeline (Hypercubic) — https://www.hypercubic.ai/insights/cobol-developer-retirement-rate-vs-training-pipeline — trade — avg age 58, 10%/yr; 71% understaffed / 54% underfunded.
[^21]: The $3 Trillion Code Nobody Knows How to Fix (Metaintro) — https://www.metaintro.com/blog/cobol-developer-shortage-legacy-systems-career-opportunity-2026 — trade — avg age 55; 85% universities dropped COBOL.
[^22]: How FS companies maintain mainframes as COBOL experts retire (BizTech) — https://biztechmagazine.com/article/2025/04/how-financial-services-companies-can-maintain-mainframes-cobol-experts-retire — trade — 71% understaffed / 54% underfunded.
[^23]: Core Banking Systems & Banktech Providers 2026 (Open Banking Tracker) — https://www.openbankingtracker.com/banktech-providers — trade — vendor directory; legacy vs cloud-native split.
[^24]: Top Core Banking Solutions market share (Verified Market Research) — https://www.verifiedmarketresearch.com/blog/best-retail-core-banking-systems/ — analyst — Temenos/Oracle/Finacle/FIS/Fiserv shares (retail segment).
[^25]: Buyer's Guide: Cloud-Based Core Banking Systems (Finantrix) — https://www.finantrix.com/buyer-guides/cloud-based-core-banking-systems — trade/analyst — born-in-cloud vs modernized-incumbent; "no Tier-1 in full production" caution.
[^26]: Core Banking Systems 2026 guide (Crassula) — https://crassula.io/guides/core-banking/ — trade — legacy-vs-cloud-native comparison; Mambu/Thought Machine/Finxact/10x/Pismo profiles.
[^27]: Market Structure of Core Banking Services Providers (Federal Reserve Bank of Kansas City) — https://www.kansascityfed.org/research/payments-system-research-briefings/market-structure-of-core-banking-services-providers/ — gov/central bank — "Big Three" dominance; Finxact acquired by Fiserv 2022.
[^28]: Retail Core Banking Systems: NA Community Bank Edition (FIS/Datos) — https://www.fisglobal.com/-/media/fisglobal/files/pdf/report/retail-banking-core-banking-systems-na-community-bank-edition.pdf — analyst/vendor — community-bank five-vendor market.
[^31]: Thought Machine Review 2026 (SectorPunk) — https://sectorpunk.com/en/reviews/thought-machine — trade — Vault smart-contract architecture; Tier-1 customers; talent scarcity.
[^33]: Banks' core technology conundrum (McKinsey) — https://www.mckinsey.com/industries/financial-services/our-insights/banks-core-technology-conundrum-reaches-an-inflection-point — analyst — progressive modernization most-pursued; adoption patterns.
[^34]: New Geoffrey Moore White Paper (AIIM blog, Jan 2011) — https://info.aiim.org/aiim-blog/newaiimo/2011/01/19/new-geoffrey-moore-white-paper-on-future-of-enterprise-it — primary — confirms Moore authored "Systems of Engagement and the Future of Enterprise IT," Jan 2011.
[^35]: Systems of Record and Systems of Engagement task force (AIIM, Oct 2010) — https://info.aiim.org/aiim-blog/newaiimo/2010/10/20/systems-of-record-and-systems-of-engagement — primary — origin of the SoR/SoE concept.
[^36]: The Future of Enterprise IT: SoR Meet SoE (Moore, LinkedIn) — https://www.linkedin.com/pulse/20121112224653-110300724-the-future-of-enterprise-it-systems-of-record-meet-systems-of-engagement — primary (author) — SoR vs SoE attributes.
[^37]: Beyond digital transformations / AI bank of the future (McKinsey) — https://www.mckinsey.com/industries/financial-services/our-insights/beyond-digital-transformations-modernizing-core-technology-for-the-ai-bank-of-the-future — analyst — engagement / decisioning / core-systems-&-data layers.
[^38]: Core-Dependent Architectures No Longer Enough (NLS Banking) — https://nlsbanking.com/blog/core-dependent-architectures-are-no-longer-enough-for-digital-banking/ — trade/vendor — "core = system of record; channel layer = system of execution."
[^39]: IBM Banking reference architecture (IBM Cloud Docs) — https://cloud.ibm.com/docs/industry-ref-arch?topic=industry-ref-arch-banking-app — vendor (primary architecture) — explicit "Core systems and systems of record" layer (often on mainframe).
[^40]: Building a core banking system with Amazon QLDB (AWS) — https://aws.amazon.com/blogs/industries/building-a-core-banking-system-with-amazon-quantum-ledger-database/ — vendor (primary) — ledger as "immutable system-of-record."
[^41]: How banks can achieve next-generation legacy modernization / Next-gen core platforms (McKinsey) — https://www.mckinsey.com/capabilities/mckinsey-digital/our-insights/tech-forward/how-banks-can-achieve-next-generation-legacy-modernization — analyst — progressive modernization + strangler; golden-source interface layer cut cost 2/3, timeline 5→3.
[^44]: How to get a core banking transformation right (McKinsey, Feb 2022) — https://www.mckinsey.com/capabilities/tech-and-ai/our-insights/tech-forward/how-to-get-a-core-banking-transformation-right-eight-mistakes-to-avoid — analyst (primary) — verbatim "only about 30 percent of CBS transformations succeeded"; 50+ analyzed; cost/timeline overruns.
[^46]: The 94% core banking problem (IBM IBV, Sept 2025) — https://www.ibm.com/thought-leadership/institute-business-value/en-us/report/core-banking-modernization — analyst (primary) — "less than half of banking CIOs reported meaningful gains."
[^47]: 2025 The voice of the makers (IBM IBV) — https://www.ibm.com/thought-leadership/institute-business-value/en-us/report/core-banking-modernization-makers — analyst — "fewer than half of CIOs reporting substantial improvements"; with BIAN.
[^49]: Barriers in IT resilience & core banking modernisation (Compliance Corylated) — https://www.compliancecorylated.com/news/barriers-remain-in-quest-for-it-resilience-and-core-banking-modernisation/ — trade — IBM sub-figures (49/46/37/31/27%).
[^50]: IBM's $40B wipeout: translating COBOL ≠ modernizing it (NovaLogiq) — https://novalogiq.com/2026/02/25/ibms-40b-stock-wipeout-is-built-on-a-misconception-translating-cobol-isnt-the-same-as-modernizing-it/ — trade — Feb 2026 ~$40B/13% drop; Brasier "ROI low"; translation ≠ modernization.
[^51]: 2 Reasons IBM's COBOL Fears Are Overblown (Motley Fool) — https://www.fool.com/investing/2026/03/26/2-reasons-ibms-cobol-fears-are-overblown/ — trade — testing/validation/regulatory sign-off is the real barrier.
[^53]: Most AI mainframe-migration vendors expected to fail by 2030 (Gartner, via Neowin) — https://www.neowin.net/news/most-ai-powered-mainframe-migration-vendors-expected-to-fail-by-2030-gartner-warns/ — analyst (Gartner) — 75% of AI mainframe-exit vendors to pivot/cease by 2030.
