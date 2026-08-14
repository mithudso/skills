---
name: jpmorgan-chase-account-intelligence
description: >-
  Company-specific account intelligence on JPMorgan Chase (JPMC) as a
  TECHNOLOGY BUYER, for a MongoDB TAM/seller working the account. PUBLIC-
  SOURCE intel; verify before customer use (leadership/budgets/org-names decay
  fast). TRIGGER: prepping for a JPMC/Chase/J.P. Morgan touchpoint, EBR, or
  account plan; JPMC LOB structure (CCB/CIB/AWM), tech
  spend/headcount/engineering scale, cloud posture (multi-cloud + private),
  data-mesh, modernization, AI strategy (LLM Suite, CDAO), tech leadership,
  outages/consent orders, procurement/third-party-risk (SMCR, BYOC), public
  MongoDB/Atlas signals, and earnings/Investor-Day tech signals. SKIP: FSI
  regulation generalities (Basel, SR 11-7, BCBS 239) -> fsi-banking-
  regulatory-context; generic big-bank IT archetype -> big-bank-IT;
  TPRM/procurement as a discipline -> enterprise-vendor-management-and-tprm;
  Goldman Sachs account -> goldman-sachs-account-intelligence;
  MongoDB/Atlas config -> mongodb-* hubs; EBR/QBR/POV mechanics -> tam-
  operations; a Chase retail/depositor/FDIC question -> personal-banking.
version: 1.1.0
updated: 2026-06-18
category: custom
whenToUse:
  - I'm preparing for a JPMorgan Chase / Chase / J.P. Morgan customer touchpoint, account plan, or EBR and need a company-specific tech-buyer profile.
  - What are JPMC's lines of business (CCB, CIB, AWM, Corporate) and how does its technology organization map to them?
  - How much does JPMorgan Chase spend on technology, how many technologists does it have, and how big is its engineering org?
  - What is JPMC's cloud posture (AWS/Azure/GCP multi-cloud vs private cloud), data-mesh/data-platform strategy, and mainframe/legacy-exit status?
  - What is JPMorgan Chase's AI strategy (LLM Suite, Chief Data & Analytics Office, AI value claims, coding assistants)?
  - Who are JPMC's key technology leaders (Global CIO, Chief Data & Analytics Officer, infrastructure CIO, AI research head) right now?
  - What outages or regulatory consent orders shape JPMC's technology and data spend?
  - How does JPMorgan Chase engage vendors (Supplier Minimum Control Requirements, BYOC/self-hosting, concentration risk) and how hard is it to sell in?
  - Is there any public signal of a MongoDB or Atlas relationship at JPMorgan Chase?
  - What did recent JPMC earnings / Investor Day say about technology and AI investment?
keywords:
  - JPMorgan Chase account intelligence
  - JPMC technology buyer profile
  - Chase J.P. Morgan technology strategy
  - JPMorgan lines of business CCB CIB AWM Corporate
  - JPMorgan technology spend budget annual
  - JPMorgan technologists headcount engineering org
  - JPMorgan multi-cloud AWS Azure GCP private cloud Gaia
  - JPMorgan data mesh data products modernization mainframe
  - JPMorgan LLM Suite generative AI CDAO
  - JPMorgan Lori Beer Global CIO Teresa Heitsenrether
  - JPMorgan consent order WhatsApp off-channel trade surveillance
  - JPMorgan Supplier Minimum Control Requirements BYOC
  - JPMorgan MongoDB Atlas signals Feinsmith
  - JPMorgan Investor Day earnings technology AI investment
  - account-based selling financial services TAM
tags:
  - account-intelligence
  - jpmorgan-chase
  - fsi
  - financial-services
  - banking
  - sales-enablement
  - enterprise-gtm
  - tam
  - account-planning
  - mongodb-tam
---

# JPMorgan Chase — Account Intelligence for a MongoDB TAM

> ## ⚠️ Verify-as-of / refresh-before-use (READ FIRST)
> This is **public-source account intelligence**, current **as of 2026-06-18**. It is **not** confidential, not insider information, and not legal/investment advice — every fact here comes from public filings, official JPMC pages, Investor-Day materials, transcripts, or reputable press, and each load-bearing claim is footnoted.
>
> **This content decays fast.** Before you use any of it with the customer, re-verify the volatile classes below against a current primary source:
> - **Leadership names & titles** — JPMC's AI-research/AI-strategy seats turned over in early-to-mid 2026 (see `references/jpmc-leadership.md`). Treat every name as a point-in-time snapshot.
> - **Technology spend & headcount** — the tech budget is re-stated annually (Jan company-update / May Investor Day). The figure used here (~$19.8B for 2026) will be superseded.
> - **Cloud-migration percentages** — JPMC re-reports app/data/infra-cloud percentages each Investor Day; they move every year.
> - **AI metrics** — LLM Suite user counts and AI-value claims change quarterly.
> - **Org names** — segment names changed in 2024 (CIB merger); internal platform codenames (e.g. "Gaia") may have evolved.
>
> When a fact is load-bearing for a customer conversation, click through to the cited source and confirm it is still current. Where this skill says a fact is **unconfirmed** or **single-source**, do not assert it to the customer.

## What this is (and is NOT)

**This is** the **company-specific** profile of JPMorgan Chase as a *technology buyer*, written for a MongoDB Technical Account Manager (or AE/SA) who works the account: its structure, the scale and shape of its technology organization, its cloud/data/AI strategy, who runs technology, what regulatory and operational events shape its spend, how it buys from vendors, and what is publicly known about MongoDB at JPMC.

**This is NOT** a primer on banking-industry regulation or the generic big-bank IT archetype. Those are owned by sibling skills — this skill *cross-references* them and stays JPMC-specific (see **Cross-references** below).

## How to use this skill (routing)

This is a hub. Load the reference that matches your need. **Every reference carries its own numbered `## Sources` list and `verified-as-of: 2026-06-18` stamp.**

| If you need… | Read |
|---|---|
| JPMC's **corporate structure & lines of business** (CCB, CIB, AWM, Corporate), firm scale, and **how the technology org maps** to the LOBs (the matrixed Global Technology + LOB-aligned model) | `references/jpmc-structure-and-tech-org.md` |
| JPMC's **technology scale** — annual tech spend (dated), technologist/engineer headcount, the engineering-org footprint, and the run-vs-invest split | `references/jpmc-tech-scale.md` |
| JPMC's **cloud, data-platform & modernization strategy** — multi-cloud (AWS/Azure/GCP) + private cloud, data mesh / data products, mainframe-exit, migration metrics, streaming/lakehouse direction | `references/jpmc-cloud-data-modernization.md` |
| JPMC's **AI strategy** — LLM Suite, the Chief Data & Analytics Office, AI value/ROI claims, key use cases | `references/jpmc-ai-strategy.md` |
| JPMC's **technology leadership** — current named leaders (Global CIO, CDAO, infra CIO, LOB CIOs, AI research) each dated, plus volatility flags | `references/jpmc-leadership.md` |
| JPMC's **incidents & regulatory orders** that shape tech/data spend (outages, the off-channel/WhatsApp order, the 2024 trade-surveillance data-controls order) | `references/jpmc-incidents-and-regulatory.md` |
| JPMC's **vendor / procurement / third-party-risk posture** (Supplier Minimum Control Requirements, the 2025 supplier open letter, BYOC/self-hosting appetite, named partnerships) **and the public MongoDB/Atlas signals** + recent earnings/Investor-Day tech signals | `references/jpmc-vendor-procurement-and-mongodb-signals.md` |

## The 60-second orientation (BLUF)

1. **JPMC is the largest US bank and the largest corporate technology spender in finance.** Its 2026 technology budget is **~$19.8 billion** (up ~10% YoY), run by a Global CIO who sits on the 12-person Operating Committee; the firm employs **~65,000 technologists** within a **~318,000**-person workforce.[^spend][^beerbio] **(Both figures volatile — re-verify; see `jpmc-tech-scale.md`.)**

2. **Tech is matrixed: one firmwide Global Technology org + budget/teams aligned into the lines of business.** JPMC reports **three** segments — **Consumer & Community Banking (CCB)**, **Commercial & Investment Bank (CIB**, the 2024 merger of the former CIB and Commercial Banking**)**, and **Asset & Wealth Management (AWM)** — plus **Corporate** as the residual. Shared platforms, infrastructure, cyber and engineering standards are centralized under the Global CIO; technologists and a large share of investment spend sit inside each LOB.[^seg][^idsplit] (See `jpmc-structure-and-tech-org.md`.)

3. **Posture is deliberate hybrid multi-cloud.** JPMC uses **all three hyperscalers (AWS — the anchor relationship — plus Azure and Google Cloud)** *and* a substantial **private cloud**, explicitly to avoid **concentration risk** in critical financial infrastructure. Its stated workload-placement decision order is **legality/compliance → security/data-control → availability/resiliency → architecture → cost (optimized LAST)**.[^multicloud][^dcd] (See `jpmc-cloud-data-modernization.md`.)

4. **The next modernization wave is the application/data layer — squarely the data-platform conversation.** JPMC says it is **"past peak modernization spend"** on infrastructure and is shifting to **modernizing application code and data to be AI-ready**. It runs a firmwide **Data Products + Data Mesh** model and is pushing into **real-time/streaming (Kafka/Confluent) and lakehouse (Databricks) with cross-platform interoperability** — while openly admitting **"data landed in cloud ≠ data that's modernized/usable for AI."**[^pastpeak][^datamesh] (See `jpmc-cloud-data-modernization.md`.)

5. **AI is a board-level, quantified program.** The internal **LLM Suite** GenAI assistant is available to **~250,000 employees** (late-2025/2026; ~150,000 weekly active per Dimon, Oct 2025) and is **model-agnostic**, routing to multiple external providers **including OpenAI and Anthropic**. All data + AI sits under a firmwide **Chief Data & Analytics Office (CDAO)**. Dimon's latest public AI-value framing: **~$2B/yr spend producing ~$2B/yr of benefit**.[^llmsuite][^aivalue] (See `jpmc-ai-strategy.md`.)

6. **Regulatory and operational events drive specific tech spend.** A **$200M SEC+CFTC** off-channel-communications order (2021) and a **~$348M OCC+Fed** trade-surveillance order (2024, explicitly a *data-controls* failure) both carry court-ordered **surveillance/archiving and data-governance** remediation — and public 2024 outages and a fraud episode reinforce **resiliency** investment.[^whatsapp][^surveil] (See `jpmc-incidents-and-regulatory.md`.)

7. **Selling in is a heavy proof-and-trust motion — and JPMC publicly prefers self-hosted/BYOC for SaaS.** Its **Supplier Minimum Control Requirements** mandate encryption, approved RTO/RPO and DR evidence, fourth-party transparency, and incident integration; an April-2025 **open letter to suppliers** demands continuous control evidence and explicitly favors **confidential computing, customer self-hosting, and bring-your-own-cloud (BYOC)**. **Cost is optimized last.**[^smcr][^openletter] (See `jpmc-vendor-procurement-and-mongodb-signals.md`.)

8. **MongoDB at JPMC is publicly *evidenced usage*, not a publicly confirmed commercial/reference relationship.** A senior JPMC executive (Head of Technology Strategy) publicly called JPMC **"a customer of Databricks, Snowflake and MongoDB,"** a JPMC job posting describes a **"strategic MongoDB Atlas team"** running an internal Atlas control plane, and a Private Bank GenAI talk shows a production agent reading from MongoDB. **There is NO MongoDB case study, press release, or earnings mention naming JPMC.** Treat internal Atlas adoption as a strong lead to validate, not a confirmed account fact.[^feinsmith][^atlasteam] (See `jpmc-vendor-procurement-and-mongodb-signals.md`.)

## TAM positioning takeaways (how to use the above)

- **Lead with compliance, security, resiliency, and auditability — not TCO.** JPMC literally optimizes cost last. The data-platform story that wins is the one that survives the Supplier Minimum Control Requirements: encryption at rest/in transit, DR/RTO/RPO evidence, granular access control, audit logging, fourth-party transparency.[^smcr][^dcd]
- **Have a self-hosted / BYOC / in-their-cloud-account answer ready.** JPMC's public supplier guidance favors confidential computing, self-hosting, and BYOC for SaaS — relevant to how you'd position MongoDB Atlas (in JPMC's own cloud accounts, with private endpoints) vs. self-managed MongoDB Enterprise Advanced.[^openletter] *(For the MongoDB deployment-model and compliance specifics, route to the `mongodb-atlas-expert` / `mongodb-operations-expert` hubs — this skill does not cover product config.)*
- **Map your conversation to the LOB *and* the central platform org.** Budget and engineers live in CCB / CIB / AWM, but shared data/infra standards and the database platform org are central. Know which you're talking to (`jpmc-structure-and-tech-org.md`).
- **Anchor to the "AI-ready data" narrative.** JPMC's own executives say data quality / fit-for-purpose is the unfinished, multiyear bottleneck for AI. That is the live business problem a modern operational-data story can attach to (`jpmc-cloud-data-modernization.md`, `jpmc-ai-strategy.md`).
- **Re-verify leadership before naming anyone.** Beer / Heitsenrether / Waldron / Alves / Haus were stable and well-sourced as of mid-2026; the AI-research and AI-strategy seats were actively churning (`jpmc-leadership.md`).

## Do-not-say / disambiguation checklist (read before a touchpoint)

These are the easy ways to mis-speak in front of JPMC. Each is sourced in the references; this box just collects them so you don't trip on one live.

- **Don't assert a MongoDB commercial relationship.** Public evidence shows JPMC *uses* MongoDB/Atlas; there is **no** public case study, deal size, or commercial-terms detail. Say "publicly-evidenced usage," and validate scope against your own account records. (`jpmc-vendor-procurement-and-mongodb-signals.md`)
- **Don't name churning AI leaders without re-verifying.** Beer/Heitsenrether/Waldron/Alves/Haus were stable as of mid-2026; the AI-research and AI-strategy seats turned over in early-to-mid 2026 and are single-source. (`jpmc-leadership.md`)
- **Don't say JPMC has "four reportable segments."** It reports **three** (CCB, CIB, AWM); **Corporate** is the residual bucket, not a reportable LOB. (`jpmc-structure-and-tech-org.md`)
- **Don't conflate the two $250M OCC penalties.** 2020 (fiduciary controls) vs 2024 (trade-surveillance *data controls*) are different matters. (`jpmc-incidents-and-regulatory.md`)
- **Don't conflate the internal firmwide data mesh with "Fusion by J.P. Morgan"** (the client-facing Securities Services product). (`jpmc-cloud-data-modernization.md`)
- **Don't state "Gaia" as the current private-cloud name with certainty.** The codename is well-sourced through 2022 but may have evolved; the platform is ongoing. (`jpmc-cloud-data-modernization.md`)
- **Don't present "75% of data in cloud" as AI-ready.** JPMC's own executives call it "landed," not modernized/usable for AI. (`jpmc-cloud-data-modernization.md`)
- **Don't quote a tech-spend or headcount number without its year.** They are re-stated annually and decay fast. (`jpmc-tech-scale.md`)
- **Don't name only OpenAI for LLM Suite.** It is model-agnostic, routing to multiple providers **including OpenAI and Anthropic**; other providers are unconfirmed. (`jpmc-ai-strategy.md`)

## Cross-references (do NOT duplicate these here)

- **FSI industry segments, regulators, and the structural laws** (Basel III, Dodd-Frank, SOX, GLBA, FFIEC, **BCBS 239**, PCI-DSS, SR 11-7 / SR 26-2, SEC 17a-4 / FINRA CAT) and *why banks buy conservatively* → **`fsi-banking-regulatory-context`**. This skill names JPMC-specific orders and maps them to spend; it does **not** re-explain the regulations. (When this skill mentions "risk-data aggregation" or "records retention," the regulatory background lives there.)
- **The generic big-bank IT archetype** (the patterns common to *any* large global bank's technology org — mainframe-plus-cloud, build-vs-buy, the central-plus-federated model in the abstract) → **`big-bank-IT`** (sibling). This skill is the **named-company instance**; the archetype is the pattern.
- **Third-party-risk-management & procurement *as a discipline*** (TPRM lifecycle, vendor risk tiering, the security-questionnaire gauntlet from the buyer side, DORA/operational-resilience regimes, the MSA/SOW/DPA contract stack) → **`enterprise-vendor-management-and-tprm`**. This skill covers *JPMC's specific* supplier controls and posture, not the general methodology.
- **MongoDB product, Atlas configuration, and compliance config** (Atlas tiers/networking/private endpoints, BYOK, Queryable Encryption, audit logging, FedRAMP/compliance posture) → the **`mongodb-atlas-expert`** and **`mongodb-operations-expert`** hubs (the latter owns encryption/CSFLE/Queryable Encryption and compliance posture). This skill stays at the account-intelligence level and does not give MongoDB configuration advice.
- **TAM deliverable mechanics** (EBR/QBR structure, account-review templates, the enterprise POC/POV technical-validation motion, account health scoring) → **`tam-operations`**. This skill is the account *content*, not the deliverable *format*.
- **Enterprise B2B buyer psychology** (buying committees, champion/mobilizer, "no decision") → **`applied-psychology`**. This skill covers JPMC's institutional facts, not the cognitive mechanics.
- **CONSUMER banking / personal finance** (a Chase depositor's accounts, FDIC insurance, credit) → **`personal-banking`**, **`consumer-finance`**. Those serve a consumer; this serves an enterprise seller.

## Provenance & verification

- Built **2026-06-18** via `/dr` deep research using `exa` + `WebSearch`/`WebFetch` against **primary public sources** (JPMorganChase 10-K & annual reports, Jamie Dimon's shareholder letters, Investor-Day decks/transcripts, the JPMC technology blog and newsroom, SEC/CFTC/OCC/Federal Reserve releases and order text, JPMC supplier-control documents, and named-executive interviews/keynotes), corroborated by reputable press (Reuters, CNBC, American Banker, Banking Dive, Constellation Research, DCD, TechTarget, CIO, InfoQ, Fortune). **60+ independent sources** across seven concept areas.
- Two of the highest-stakes, most-volatile claims were independently re-verified during authoring against primary regulator text: the **2021 SEC ($125M) + CFTC ($75M)** off-channel order (sec.gov press release 2021-262; cftc.gov 8470-21) and the **2024 OCC ($250M) + Federal Reserve (~$98.2M)** trade-surveillance/data-controls order (occ.gov nr-occ-2024-25; federalreserve.gov enforcement20240314a).
- The MongoDB signal was deliberately stress-tested for over-claiming: the absence of any JPMC case study on mongodb.com was confirmed, and the relationship is characterized as **publicly-evidenced usage, not a confirmed commercial relationship.**
- **Re-verify any 2026-dated claim before customer use.** This account moves every quarter; leadership and spend figures are the fastest-decaying.

[^spend]: 2026 JPMC Company Update transcript (~$19.8B technology spend, up ~10% YoY), Jan 2026; corroborated by American Banker, "JPMorgan invests for new age of competition amid AI fears" (2026) and Business Insider/AOL "almost $20 billion" (2026-02-24). See `references/jpmc-tech-scale.md`.
[^beerbio]: Lori Beer official JPMorganChase bio (~$19.8B budget, ~65,000 technologists), jpmorganchase.com/about/leadership/lori-beer (accessed 2026); Fortune via Yahoo (2026-04-29). Workforce ~318,000: JPMC 10-K FY2024 / Banking Dive Investor Day coverage (2025-05-19).
[^seg]: JPMC 10-K FY2024 (filed 2025-02; three reportable segments — CCB, CIB, AWM — with Corporate as residual; CIB formed by combining the former Corporate & Investment Bank and Commercial Banking effective Q2 2024), sec.gov; PYMNTS (2024-01-25). See `references/jpmc-structure-and-tech-org.md`.
[^idsplit]: 2025 JPMC Investor Day presentation (technology investment by LOB: CCB $3.2B, CIB $3.7B, AWM $1.0B) and Banking Dive (2022) on the centralized-plus-LOB-aligned model (~$4B aligned to the LOBs). See `references/jpmc-structure-and-tech-org.md`.
[^multicloud]: JPMC infrastructure leaders on multi-cloud + concentration risk: DataCenterDynamics, "The IT wingspan of JPMorgan Chase & Co." (2025/2026); SiliconANGLE (2024-07-27); GeekWire (2019). See `references/jpmc-cloud-data-modernization.md`.
[^dcd]: DataCenterDynamics (2025/2026), JPMC infra CIO on the compliance→security→resiliency→architecture→cost decision order. See `references/jpmc-cloud-data-modernization.md`.
[^pastpeak]: "Past peak modernization spend" / shift to application code + data for AI-readiness: CFO Jeremy Barnum, 2025 Investor Day (Banking Dive 2025-05-19) and 2026 Company Update transcript. See `references/jpmc-cloud-data-modernization.md`.
[^datamesh]: JPMC Technology blog, "Evolution of data mesh architecture" (2023-06-19); 2024 Investor Day transcript ("landing is not the same as … modernized and usable for AI"). See `references/jpmc-cloud-data-modernization.md`.
[^llmsuite]: LLM Suite: JPMC Technology blog (2025-06-03); American Banker (2025-05-22, 2025-07-07); user-count trajectory and OpenAI+Anthropic routing per AI Magazine (2025-10-03), Dimon on Bloomberg TV (2025-10-27). See `references/jpmc-ai-strategy.md`.
[^aivalue]: Dimon, Bloomberg TV (2025-10-27): "~$2B of expense, ~$2B of benefit"; eMarketer (2025-10-09). Trajectory from $1.5B (2023) target. See `references/jpmc-ai-strategy.md`.
[^whatsapp]: SEC press release 2021-262 ($125M) and CFTC 8470-21 ($75M), 2021-12-17; both orders' undertakings reference surveillance/archiving technology remediation. See `references/jpmc-incidents-and-regulatory.md`.
[^surveil]: OCC nr-occ-2024-25 ($250M + cease-and-desist) and Federal Reserve enforcement20240314a (~$98.2M), 2024-03-14; findings cite trading-activity surveillance gaps "without adequate data controls." See `references/jpmc-incidents-and-regulatory.md`.
[^smcr]: JPMorganChase Supplier Minimum Control Requirements (public PDF; 2020 & 2025 updates). See `references/jpmc-vendor-procurement-and-mongodb-signals.md`.
[^openletter]: JPMC Technology blog, "An open letter to our third-party suppliers" (April 2025) — continuous control evidence; confidential computing / self-hosting / BYOC preference. See `references/jpmc-vendor-procurement-and-mongodb-signals.md`.
[^feinsmith]: Larry Feinsmith (MD, Head of Technology Strategy, Innovation & Partnerships, JPMorganChase), quoted by Constellation Research: JPMC "a customer of Databricks, Snowflake and MongoDB." See `references/jpmc-vendor-procurement-and-mongodb-signals.md`.
[^atlasteam]: JPMorganChase "Lead Software Engineer – Cloud Databases" posting describing a "strategic MongoDB Atlas team" / internal Atlas control plane (Built In London); corroborated by multiple JPMC MongoDB job postings and a LangChain "Interrupt" (2025) Private Bank agent talk. No JPMC case study exists on mongodb.com (confirmed absence). See `references/jpmc-vendor-procurement-and-mongodb-signals.md`.
