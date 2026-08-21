# JPMC — Cloud, Data-Platform & Modernization Strategy

> Provenance: reference under the `jpmorgan-chase-account-intelligence` skill. Built 2026-06-18 via `/dr` deep research (exa + WebSearch/WebFetch; primary sources). **verified-as-of: 2026-06-18.** Migration percentages, codenames, and data-center counts are **highly volatile** — re-verify before customer use.

## Contents
- Public-cloud posture (multi-cloud, AWS anchor)
- Internal / private cloud
- Workload-placement decision order
- Data platform: Data Products + Data Mesh
- Modernization & mainframe-exit
- Cloud-migration metrics (year-stamped)
- Data-architecture direction (streaming, lakehouse)
- TAM implications
- Sources

## Public-cloud posture — deliberate multi-cloud

JPMC explicitly describes a **hybrid, multi-cloud** strategy — sometimes framed publicly by its engineers as **"three-plus-one"** (three public hyperscalers + its own private cloud).[^1] The three public clouds are **AWS, Microsoft Azure, and Google Cloud**; Dimon has publicly confirmed JPMC uses all three.[^2][^3]

- **Strategic intent = avoid lock-in / concentration risk + resiliency.** JPMC's public-cloud leaders state it plainly: *"We want to avoid any lock-in or concentration risk."*[^4] Its infrastructure CIO: *"we are very concerned about concentration risk. If a cloud provider had an outage, the financial services industry cannot be impacted by that — it is critical infrastructure, and we don't want to be beholden or at risk with any single vendor."*[^5] **Confidence: fact** (consistent 2019→2026).
- **AWS is the most publicly documented (anchor) relationship.** The Global CIO at AWS re:Invent (2024) described going from ~100 apps on AWS (2020) to nearly 1,000 by 2023, including consumer Chase.com/mobile running **active-active across multiple AWS regions** ("two of three regions can fail without customer impact"), and **Chase UK built "from the ground up on AWS."** JPMC uses AWS Graviton, EC2, SageMaker, and Bedrock.[^6] **Confidence: fact.**
- **Forward signal (QCon London, March 2026):** JPMC engineers argued enterprises should "treat multi-cloud as a product, not a project," organizing around cross-cloud capabilities — driven partly by AI teams wanting Vertex AI + Bedrock + Azure OpenAI simultaneously. **Qualified** (engineer talk; reflects current thinking).[^7]

## Internal / private cloud

JPMC runs a substantial **private cloud**, historically code-named **"Gaia."** Origin: American Banker reported the "Gaia" internal-cloud codename in 2016;[^8] JPMC engineers presented "Project Gaia" (built on Cloud Foundry) at Cloud Foundry Summit 2017;[^9] a 2022 VMware Tanzu talk referenced the **"Gaia Kubernetes platform."**[^1] **Confidence on the codename: fact** (multiple independent sources 2016–2022). **Caveat:** no 2025–2026 primary source re-confirms the exact current internal branding — the platform is unambiguously ongoing, but treat the *name* "Gaia" as possibly evolved (flag as such to the customer).

- JPMC invested **~$2B to build new US private-cloud data centers** (announced ~2020/2021) and retains a large on-prem/private footprint deliberately for **regulatory/data-sovereignty** reasons and for workloads it judges it can secure better itself.[^5][^10]
- Data-center footprint: ~**32 data centers** consolidating toward **~17–20** "highly automated" facilities (the infra CIO calls 17 a "soft goal," no fixed deadline); new modules are liquid-cooling-ready for GPUs.[^5]
- Dimon's public framing (2024 letter): public cloud can provide *"30% additional efficiency if done correctly,"* while private cloud trades external innovation for control.[^11]

## Workload-placement decision order (stated publicly)

JPMC's infra leadership states an explicit ordering for where a workload runs:[^4][^5]

**legality/compliance → security/data-control → customer availability/resiliency → architecture → (last) cost optimization.**

Customer-facing / fast-time-to-market / lower-data-sensitivity apps tend to public cloud; ultra-low-latency (some trading) and most-sensitive data stay close/private. **For a TAM: cost is optimized LAST — lead with compliance, security, and resiliency, never TCO.**

## Data platform — Data Products + Data Mesh

JPMC's official technology blog documents a firmwide **Data Products + Data Mesh** architecture:[^12]
- Each data product gets its **own product-specific data lake** with **physical separation**; cloud-based storage + cataloging (the blog names AWS S3, Glue, Lake Formation, Athena).
- A cloud-based **Mesh Catalog** tracks cross-enterprise data flows.
- **"In-place consumption"** — share, don't copy — with granular column/record/value-level access control.

Corroborated by an AWS Big Data Blog co-authored with JPMC and by theCUBE coverage of JPMC's data-mesh talks.[^13] A separate **client-facing** product, "Fusion by J.P. Morgan" (Securities Services), exposes a data mesh to institutional investors via API/Snowflake/Databricks — **do not conflate** the internal firmwide mesh with the client product.[^14]

## Modernization & mainframe-exit

JPMC is **actively retiring legacy apps and migrating off legacy data centers, but deliberately keeps the mainframe for the most critical systems.** The infra CIO (2026): *"a lot of the systems that were using the mainframe have been modernized … but the most critical systems are still using it."* JPMC still hires COBOL programmers.[^5][^15] **For a TAM: "move off mainframe" is real but partial and selective — not a wholesale rip-and-replace.**

Legacy-app decommissioning is a long-running program: **>2,500 legacy applications decommissioned since 2017** (2023 Investor Day baseline).[^3] A 2026 press item cited a "500th legacy application decommissioned in Q1 2026" milestone — **treat as tentative** pending a primary JPMC confirmation.

## Cloud-migration metrics (year-stamped — HIGHLY VOLATILE)

JPMC re-reports these each Investor Day. Note JPMC uses **two different denominators** ("apps running large workloads in cloud" vs "apps on modern infrastructure incl. virtual servers"):[^3][^11][^5]

| Metric | 2023 | 2024 | 2025 (reported May 2025) |
|---|---|---|---|
| Apps running large workloads in public/private cloud | ~38% (infra in cloud) | ~50% | **~65%** |
| Apps on "modern infrastructure" (incl. virtual servers) | ~56% modern | ~80% (target) | **~80%** |
| Data in public/private cloud | — | ~70%→75% target | **~75%** ("landed") |

**Critical caveat (JPMC's own words):** *"70% of our data has been landed in the cloud, but … landing is not the same as having it be modernized and usable for AI and ML, so there's work to do there still."*[^16] Dimon (2025): *"Data is the hardest part … getting the data in the form that's usable is the hard part."*[^17] **The headline "75% of data in cloud" is "landed," not "AI-ready."** This is the live, unfinished business problem.

## Data-architecture direction (streaming, lakehouse)

- **Real-time / streaming is an explicit pillar.** 2025 Investor Day: a subset of data "will need to be streamed real-time," with progress on servicing/personalization; a senior data exec: *"we are focused on ensuring our data is available in milliseconds."*[^16][^18]
- **Apache Kafka / Confluent** is described as JPMC's "digital nervous system," used as the data-movement layer of the data mesh (managed Kafka Connect, multi-DC replication). Internal codenames in some syntheses are **tentative** — don't assert them.[^19]
- **Lakehouse / Databricks + interoperability.** JPMC is described as a "long-time Databricks customer"; a senior data exec stresses *"not all of our data will wind up in Databricks. Interoperability is very important."* Snowflake also appears as a supported channel.[^17][^20] **Qualified** — the stated direction is "central but federated," multi-platform interoperable.
- **Databases (public signal):** AWS re:Invent sessions show JPMC internal platforms on AWS using EKS, Lambda, RDS, Kinesis (streaming) and **"enterprise NoSQL"** on EC2; job posts repeatedly list NoSQL, Spark/Flink streaming, MSK/S3/EMR. **The architecture is cloud-native, microservices, event-driven, NoSQL-inclusive.** For the MongoDB-specific public signal, see `jpmc-vendor-procurement-and-mongodb-signals.md`.[^21]

## TAM implications

- **The data-platform conversation is the current one.** "Past peak modernization" + "data landed ≠ AI-ready" + a real-time/interoperability push = an opening for a modern operational-data story attached to AI-readiness.
- **Respect the concentration-risk doctrine.** JPMC values portability and multi-cloud optionality; a data layer that runs the same across AWS/Azure/GCP/private fits the doctrine. (Route MongoDB deployment-model specifics to `mongodb-atlas-expert`/`mongodb-operations-expert`.)
- **Resiliency is non-negotiable.** Active-active multi-region, DR/RTO/RPO evidence, and "the FSI can't be taken down by one cloud outage" are explicit values — bring that story.

## Sources

1. VMware Tanzu, "Improving JPMorgan Chase's developer experience on the cloud" (YouTube, 2022) — "three-plus-one" strategy; "Gaia Kubernetes platform." (Vendor talk, primary-adjacent) — https://www.youtube.com/watch?v=QOvBWlf7Cgg
2. DataCenterDynamics, "The IT wingspan of JPMorgan Chase & Co." (2025/2026) — Dimon confirms AWS/Azure/GCP; multi-cloud posture. (Press) — https://www.datacenterdynamics.com/en/analysis/the-it-wingspan-of-jpmorgan-chase-co/
3. Constellation Research, "JPMorgan Chase digital transformation, AI and data strategy" — migration metrics; ">2,500 legacy apps since 2017"; data mesh. (Press/analyst) — https://www.constellationr.com/insights/news/jpmorgan-chase-digital-transformation-ai-and-data-strategy-sets-generative-ai
4. SiliconANGLE, "A close look at JPMorgan's aggressive cloud migration" (2024-07-27) — "avoid any lock-in or concentration risk"; ~6,000 apps; 32→17 DCs. (Press) — https://siliconangle.com/2024/07/27/close-look-jpmorgans-aggressive-cloud-migration/
5. DataCenterDynamics (2025/2026), infra CIO interview — concentration risk; decision order; mainframe; 32→~17 DCs. (Press) — https://www.datacenterdynamics.com/en/analysis/the-it-wingspan-of-jpmorgan-chase-co/
6. CIO, "JPMorgan Chase builds ambitious AI foundation on AWS" (2024-12) — AWS app counts; Chase.com multi-region; Chase UK on AWS; Graviton/SageMaker/Bedrock. (Press) — https://www.cio.com/article/3616622/jpmorgan-chase-builds-ambitious-ai-foundation-on-aws.html
7. InfoQ, "Your multi-cloud strategy is a product problem" (2026-03) — QCon London JPMC talk; AI teams driving multi-cloud. (Press/conference) — https://www.infoq.com/news/2026/03/jpmc-multicloud-product-strategy/
8. American Banker, "Unexpected champion of public clouds: JPMorgan CIO Dana Deasy" (2016-09-22) — original "Gaia" internal-cloud codename. (Press) — https://www.americanbanker.com/news/unexpected-champion-of-public-clouds-jpmorgan-cio-dana-deasy
9. Cloud Foundry Summit 2017, "Project Gaia" session (Liste & Worgan) — Gaia private cloud on Cloud Foundry. (Conference, primary-adjacent) — https://cfsummit2017.sched.com/event/AJm9/
10. TechTarget, "Innovation with cost control" (2024-08-26) — private-cloud investment; Thought Machine Vault core; 32→~20 DCs; FinOps. (Press) — https://www.techtarget.com/searchcio/feature/JPMorgan-Chase-technology-goal-Innovation-with-cost-control
11. Jamie Dimon 2024 Letter to Shareholders (PDF) — "30% additional efficiency"; "complete the migration of our analytics data estate to the public cloud." (Annual report, primary) — https://www.jpmorganchase.com/ir/annual-report
12. JPMorganChase Technology blog, "Evolution of data mesh architecture" (2023-06-19) — data products, per-product data lake, in-place consumption, Mesh Catalog. (Blog, primary) — https://www.jpmorganchase.com/about/technology/blog/evolution-of-data-mesh-architecture
13. AWS Big Data Blog, "How JPMorgan Chase built a data mesh architecture" (2021) — data-mesh on AWS (co-authored w/ JPMC). (Vendor blog, primary-adjacent) — https://aws.amazon.com/blogs/big-data/
14. J.P. Morgan, "Fusion / Data Mesh" client product page — client-facing data mesh via API/Snowflake/Databricks. (Product, primary) — https://www.jpmorgan.com/insights/securities-services/data-solutions/data-access-catalog-mesh
15. CIO Dive, "How JPMorgan's infrastructure chief keeps the AI engine humming" (2025-01-30) — mainframe-to-quantum; hybrid AI compute. (Press) — https://www.ciodive.com/news/jpmorgan-chase-infrastructure-cio-ai-compute-strategy/738662/
16. JPMorgan Chase 2024 Investor Day, full transcript (PDF) — "landing is not the same as … modernized and usable for AI." (Investor Day transcript, primary) — https://www.jpmorganchase.com/ir/investor-day
17. Constellation Research, "Dimon on AI, data, cybersecurity" — "data is the hardest part"; "long-time Databricks customer." (Press) — https://www.constellationr.com/insights/news/jpmorgan-chases-it-ai-bets-where-returns-are
18. McKinsey, "Data in the age of AI: a conversation with Mark Birkhead of JPMorgan Chase" (2025-06) — "available in milliseconds"; AI-ready data. (Interview, primary exec quotes) — https://www.mckinsey.com/capabilities/tech-and-ai/our-insights/data-in-the-age-of-ai-a-conversation-with-mark-birkhead-of-jpmorgan-chase
19. Factor House, "How JPMorgan uses Apache Kafka in production" — Kafka "digital nervous system"; Confluent; managed Connect. (Press/vendor synthesis; codenames tentative) — https://factorhouse.io/articles/jpmorgan-kafka-architecture
20. JPMorgan Chase 2025 Investor Day, full transcript (PDF) — target-platform consolidation "about halfway"; real-time streaming progress. (Investor Day transcript, primary) — https://www.jpmorganchase.com/ir/investor-day
21. AWS re:Invent 2021, "How JPMC modernized its core risk-management platform" (FSI202, PDF) — internal platform on EKS/Lambda/RDS/Kinesis + "enterprise NoSQL" on EC2. (Conference deck, primary-adjacent) — https://d1.awsstatic.com/events/reinvent/2021/How_JPMC_modernized_its_core_risk_management_platform_FSI202.pdf
