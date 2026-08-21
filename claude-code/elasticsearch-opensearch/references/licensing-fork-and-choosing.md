# Licensing, the 2021 Fork, and Choosing Between Elasticsearch and OpenSearch

> Provenance: reference under the `elasticsearch-opensearch` skill. Researched 2026-06-18 via `/dr` deep research (multi-source: Elastic blog, OpenSearch.org, AWS Open Source blog, OSI, Linux Foundation, plus independent secondary). **verified-as-of 2026-06-18.** This is the most volatile part of the skill — re-verify versions and license terms before quoting.

## Contents

- The prehistory (2018–2019) — why the conflict started
- The 2021 license change (Apache-2.0 → SSPL + Elastic License v2)
- What SSPL is and why it is not OSI-approved
- The OpenSearch fork (2021) and the 2024 Linux Foundation move
- Elastic re-adding AGPLv3 (2024) — "Open Source. Again!"
- How to choose today
- Version landmarks

## The prehistory (2018–2019)

The fork did not come from nowhere. Starting June 2018, Elastic began intermingling proprietary "X-Pack" features into the formerly pure-Apache-2.0 codebase. AWS publicly objected, and in March 2019 (with Expedia and Netflix) launched **Open Distro for Elasticsearch**, a 100% Apache-2.0 distribution adding security, alerting, SQL, etc. — explicitly framed at the time as *"not a fork."*[^7][^8] AWS's Adrian Cockcroft framed it as "Keeping Open Source Open," arguing that mixing proprietary and open code meant "the majority of new Elasticsearch users are now, in fact, running proprietary software."[^7][^8] This is the direct precursor to both the license change and the eventual fork.

## The 2021 license change

On **January 14, 2021**, Elastic (Shay Banon) published "Doubling down on open, Part II," announcing that **starting with the 7.11 release**, the Apache-2.0-licensed source of **Elasticsearch and Kibana** would move to a **dual license: SSPL (Server Side Public License) OR the Elastic License**, at the user's choice.[^1] The 7.11 GA followed in early February 2021 (the license-header commit landed February 3, 2021).[^1][^5]

- **Date nuance:** Elastic's blog is dated **Jan 14, 2021**, but AWS's own materials repeatedly cite **"January 21, 2021."** Both refer to the same event/window; the Jan 14 blog is the primary announcement.[^11] *(Report both; do not "correct" one to the other.)*
- **Stated motive — AWS:** Elastic's reason was to protect its investment "by restricting cloud service providers from offering Elasticsearch and Kibana as a service without contributing back."[^1] A follow-up clarification joked about abuse "from you know who :)" — i.e., Amazon.[^2] AWS's own framing confirms the conflict: Elastic objected that AWS "was profiting from hosting open source Elasticsearch but not collaborating," and to AWS's use of the "Elasticsearch" trademark.[^21]

**Elastic License v2 (ELv2)**, introduced **February 2, 2021**, is a simplified "fair-code" license: free use/redistribution/modification with **three limitations** — (1) no providing it as a managed service, (2) no circumventing the license key for paid features, (3) no removing/obscuring licensing functionality.[^3] Because >90% of Elastic's downloads were already under the Elastic License, Elastic called the change "a non-event" for most users; the **default distribution is ELv2**, and Elastic stopped producing an Apache-2.0 ("OSS") distribution.[^3][^4][^6] In June 2021, ELv2 was rolled out to Logstash, Beats, Elastic Agent, APM-Server, and ECK as well.[^6]

## What SSPL is and why it is not OSI-approved

- **SSPL** is a **source-available copyleft license created by MongoDB**, extending the GPL/AGPL framework. It permits free use, modification, and redistribution, **but** if you offer the software *as a service*, you must publicly release the source of your *entire management/service layer* under SSPL.[^1][^2] (fact — 3 sources)
- It is **not OSI-approved and not considered open source.** Elastic itself stated this directly: "it is not an OSI approved license and is not considered open source."[^2]
- **OSI's position** (published Jan 19, 2021): the SSPL "is Not an Open Source License"; OSI calls relicensing-then-still-calling-it-open "fauxpen source." Critically, **SSPL was submitted to OSI by MongoDB and then withdrawn (March 2019)** by license steward Eliot Horowitz "when it became clear that the license would not be approved," citing lack of community consensus around the copyleft provision. OSI's process page adds: "given the scope of copyleft in the [SSPL], it is not a license that anyone currently would be able to comply with."[^9][^10] (fact — multiple OSI primary sources)

## The OpenSearch fork (2021) and the 2024 Linux Foundation move

**Fork timeline:**
- **April 12, 2021** — AWS announced **OpenSearch**, a community-driven Apache-2.0 fork: **OpenSearch derived from Elasticsearch 7.10.2** and **OpenSearch Dashboards from Kibana 7.10.2** (the last ALv2 versions). Open Distro was folded into it. AWS warned the initial code was "alpha…not suitable for production." To preserve attribution the team took the full git history of the 7.10 branch, then stripped all X-Pack/Elastic-licensed code, branding, and telemetry.[^12][^17]
- **July 12, 2021** — **OpenSearch 1.0 GA** (first production-ready release).[^13][^14]
- **September 8, 2021** — **Amazon Elasticsearch Service renamed to Amazon OpenSearch Service** (still offering 19 ALv2 Elasticsearch versions, 7.10 and earlier).[^11]

**Linux Foundation move (2024):**
- **September 16, 2024** — AWS transferred OpenSearch to the newly formed **OpenSearch Software Foundation under the Linux Foundation**, announced at Open Source Summit Europe. Governance = a **governing board** (foundation members) plus an independent **Technical Steering Committee (TSC)** for technical decisions.[^18][^19][^20]
- **Founding members:** Premier = **AWS, SAP, Uber**; general members included Aiven, Aryn, Atlassian, Canonical, DigitalOcean, Eliatra, Graylog, Instaclustr by NetApp, Portal26.[^19][^20]
- **Scale at transfer:** >700 million downloads, thousands of contributors, 200+ maintainers.[^18][^20]
- **Stated motive:** shed the "AWS-controlled project" perception and remove a "headwind" for enterprises "reticent to invest heavily in vendor-controlled open-source projects."[^20][^23]

## Elastic re-adding AGPLv3 (2024) — "Open Source. Again!"

- **August 29, 2024** — Shay Banon published "Elasticsearch Is Open Source. Again!" announcing that **AGPLv3** (an **OSI-approved** license) would be added as a **third option** alongside SSPL and ELv2 for the free portions of the Elasticsearch/Kibana source. Nothing was removed — it is purely additive; existing ELv2/SSPL users are unaffected; binary distributions and the Apache-2.0 client libraries are unchanged.[^24][^4]
- **Exact version/date:** the change was "expected to take place **before the 8.16 release** is generally available." The AGPLv3-support commit landed **September 13, 2024**; **Elasticsearch 8.16 GA was November 12, 2024**.[^4][^25][^26][^27]
- **Banon's reasoning:** "Open source is in my DNA." He argued the 2021 change had successfully mitigated the market-confusion risk over three years, so Elastic now "feel[s] comfortable adding AGPL," explicitly tying it to AI/RAG adoption and vector search. He pointedly noted MongoDB *used to be* AGPL and Grafana *moved to* AGPL, pre-empting "AGPL isn't real open source" objections.[^24][^25]
- **Scope caveat:** AGPL applies to a **subset** of the free source code — *not* the paid features, which stay ELv2; "releases will continue to be under the Elastic License."[^4][^25]
- **Causation flag (commentary, not Elastic's stated reason):** multiple outlets noted the AGPL move coincided closely (~Sept 16, 2024) with OpenSearch's Linux-Foundation transfer and speculated it was partly a competitive response to OpenSearch gaining a neutral governance home. Treat the **coincident timing** as fact and the **causal link** as speculation.[^23]

## How to choose today (verified-as-of 2026-06-18)

Below the 7.10 fork point the engines are essentially identical; above it they have meaningfully diverged.[^31] Note: several feature comparisons cite vendor-authored benchmarks — flagged below.

**Licensing implications (the cleanest differentiator):**
- **OpenSearch = Apache 2.0** throughout: permissive, OSI-approved, vendor-neutral procurement-friendly, embeddable, can be offered as a managed service.[^30][^31]
- **Elasticsearch = tri-license source (AGPLv3 / SSPL / ELv2)**; default binaries are **ELv2** — you **cannot** offer it as a managed service or circumvent paid-feature gating, and AGPLv3's network-copyleft is itself a constraint some orgs avoid.[^4][^31]

**Feature divergence (verified-as-of 2026-06):**
- **Query languages:** Elasticsearch has **ES|QL** (piped, deeply integrated into Kibana). OpenSearch has **PPL** (Piped Processing Language, Splunk-SPL-inspired) + a SQL plugin. Both support Query DSL. (ES|QL is Elastic-only — see `query-dsl-aggregations-relevance.md`.)
- **Semantic search / ML:** Elasticsearch is generally **ahead** for turnkey AI — **ELSER** (Elastic Learned Sparse EncodeR), native model inference, the Inference API, the retriever/RRF framework, `semantic_text`. OpenSearch uses **Neural Search + ML Commons** (bring-your-own / remote models) and has pushed hard on agentic capabilities (native MCP support in 3.x).[^28][^29][^30]
- **Vector / k-NN:** both do HNSW + hybrid (BM25+vector) + RRF. Elasticsearch: dense vectors to ~4,096 dims, BBQ (Better Binary Quantization). OpenSearch: FAISS + NMSLIB + Lucene engines, to ~16,000 dims, GPU acceleration in 3.x.[^28][^29][^30]
- **Security plugin (OpenSearch's biggest cost argument):** **OpenSearch ships RBAC, document/field-level security, audit logging, and anomaly detection free by default.** In Elasticsearch these (plus SSO, ML, Cross-Cluster Replication, searchable snapshots) sit behind **paid tiers**.[^28][^30]
- **Observability:** Elasticsearch pushes its own polished Elastic Agent + APM; OpenSearch is **OpenTelemetry-first** (Data Prepper + OTel collectors, no proprietary APM agents).[^28][^30]

**Performance — contested (flag bias):** Elastic's own benchmarks claim Elasticsearch is significantly faster (a 2023 post claimed 40%–140%; later posts claim large vector-search advantages). **These are vendor-authored.** An **independent Trail of Bits assessment (March 2025)** reached the opposite conclusion — OpenSearch 2.17.1 ~1.6x faster on the Big5 workload and ~11% faster on vector search than Elasticsearch 8.15.4 — and an independent methodology critique argued Elastic's 2023 results were dominated by network latency / cached aggregations. **Net: performance claims are benchmark-dependent and both vendors cherry-pick; treat as inconclusive and test on your own workload.**[^32][^33]

**Managed service & cost:**
- OpenSearch → **Amazon OpenSearch Service** (plus Aiven, Instaclustr, etc.); software cost **$0 at any scale** (infrastructure-only); generally cheaper on AWS.[^28][^30][^31]
- Elasticsearch → **Elastic Cloud** (AWS/Azure/GCP) + Elastic Cloud Serverless; freemium with a paid ceiling (Platinum/Enterprise for security, ML, CCR).[^28][^30]

**Decision heuristic (synthesized):**[^29][^30][^31]
- **Choose Elasticsearch** if you need the latest ML/vector/semantic features first (ELSER, ES|QL, Inference API for RAG); want Elastic Observability/Security as a polished suite; want Elastic Cloud SLA/support; and SSPL/ELv2/AGPL constraints don't affect your use (typical internal deployment).
- **Choose OpenSearch** if you need **Apache 2.0** (procurement, embedding, offering-as-a-service); want **free** built-in security/RBAC; are deep in AWS where OpenSearch Service pricing decides it; value **Linux Foundation vendor-neutral governance**; or are cost-sensitive on log/observability at scale.
- **Migration caveat:** index formats and ML/vector setups have diverged; cross-migration generally requires reindexing.[^31]

## Version landmarks (verified-as-of 2026-06-18)

**Elasticsearch:**
- **8.x:** GA Feb 2022; extended a final time to **8.18 / 8.19** (GA alongside 9.x) so laggards could get features.
- **9.0 GA — April 15, 2025** (built on **Lucene 10**), released simultaneously with 8.18.
- **Current line as of mid-2026: 9.x** (e.g., 9.4.x tagged by June 2026 per GitHub releases / endoflife.date).[^27]

**OpenSearch:**
- **1.0 GA — July 12, 2021** (forked from ES 7.10.2).[^13]
- **3.0 GA — May 6, 2025** (first major since 2022; built on **Lucene 10**; JDK 21).[^34]
- **Current line as of mid-2026: 3.x** (e.g., 3.7.0 ~June 2026 per OpenSearch docs version history); 2.x in maintenance.[^34]

**Cross-note:** Elasticsearch 9.0 and OpenSearch 3.0 both jumped to **Lucene 10** within ~3 weeks of each other (April 15 vs May 6, 2025) — a clean anchor for "they still share the Lucene foundation."

## References

1. Elastic — "Doubling down on open, Part II" (Jan 14 2021): https://www.elastic.co/blog/licensing-change — primary license-change announcement (7.11, SSPL+Elastic dual). blog
2. Elastic — license-change clarification (Jan 19 2021): https://www.elastic.co/blog/license-change-clarification — SSPL "not OSI approved…not considered open source," Amazon. blog
3. Elastic — Elastic License v2 (Feb 2 2021): https://www.elastic.co/blog/elastic-license-v2 — three limitations. blog
4. Elastic — licensing FAQ: https://www.elastic.co/pricing/faq/licensing — 7.11/2021 SSPL+ELv2; Sept 2024 AGPLv3 before 8.16; no Apache distro. docs
5. Elastic — license-header commit (Feb 3 2021): https://github.com/elastic/elasticsearch (commit a92a647…) — SSPL+Elastic-2.0 headers, OSS distro removed. docs/commit
6. Elastic — ELv2 update (Jun 3 2021): https://www.elastic.co/blog/elastic-license-update — ELv2 to Logstash/Beats/Agent/APM/ECK. blog
7. AWS — "Keeping Open Source Open: Open Distro" (Cockcroft, Mar 11 2019): https://aws.amazon.com/blogs/opensource/keeping-open-source-open-open-distro-for-elasticsearch/ — proprietary intermingling since June 2018. blog
8. AWS Insider — Open Distro launch + Cockcroft quotes: https://awsinsider.net/articles/2019/03/12/elasticsearch.aspx — secondary. news
9. OSI — "The SSPL is Not an Open Source License" (Jan 19 2021): https://opensource.org/blog/the-sspl-is-not-an-open-source-license — fauxpen, SSPL withdrawn before approval. blog (OSI)
10. OSI — license review process + MongoDB SSPL withdrawal (Horowitz, 2019): https://opensource.org/licenses/review-process — "not a license anyone currently could comply with." docs/forum (OSI)
11. AWS — "Announcing Amazon OpenSearch Service" (Sep 8 2021): https://aws.amazon.com/blogs/aws/announcing-amazon-opensearch-service-which-supports-opensearch-10/ — cites "Jan 21 2021"; service rename. blog (AWS)
12. AWS — "Introducing OpenSearch" (Apr 12 2021): https://aws.amazon.com/blogs/opensource/introducing-opensearch/ — fork from ES/Kibana 7.10.2, ALv2, "alpha." blog (AWS)
13. AWS — "OpenSearch 1.0 launches" (Jul 12 2021): https://aws.amazon.com/blogs/opensource/opensearch-1-0-launches/ — GA, production-ready. blog (AWS)
14. InfoQ — OpenSearch 1.0 GA: https://www.infoq.com/news/2021/08/opensearch-elastic-amazon/ — secondary confirmation. news
17. InfoQ — Amazon OpenSearch fork mechanics (Apr 2021): https://www.infoq.com/news/2021/04/amazon-opensearch/ — full 7.10 git history, X-Pack/telemetry stripped. news
18. Linux Foundation — OpenSearch Software Foundation press (Sep 16 2024): https://www.linuxfoundation.org/press — transfer, >700M downloads, 200+ maintainers. news (LF)
19. AWS + OpenSearch.org — LF transfer & governance: https://aws.amazon.com/blogs/opensource/aws-welcomes-the-opensearch-foundation/ ; https://opensearch.org/blog/building-the-future-of-opensearch-together/ — TSC + governing board, members. blog
20. TechCrunch — AWS brings OpenSearch under Linux Foundation (Sep 16 2024): https://techcrunch.com/2024/09/16/aws-brings-opensearch-under-the-linux-foundation-umbrella/ — SAP/Uber premier, "shed AWS-driven perception." news
21. DevClass — AWS hands OpenSearch to LF: https://www.devclass.com/ — Elastic–AWS trademark history. news
23. AWS Insider / SiliconAngle — "vendor-controlled," "reticent to invest" skepticism: https://awsinsider.net/articles/2024/09/17/opensearch-changes-hands.aspx — disconfirming/commentary. news
24. Elastic — "Elasticsearch Is Open Source. Again!" (Banon, Aug 29 2024): https://www.elastic.co/blog/elasticsearch-is-open-source-again — AGPL added, pre-empts critiques. blog
25. Business Wire + Elastic DevRel newsletter — AGPL "before 8.16 GA," vector-search quote: https://www.businesswire.com/news/home/20240829537786/en/ — press. news/blog
26. Elastic — AGPLv3-support commit (Sep 13 2024): https://github.com/elastic/elasticsearch (commit a59c182…). docs/commit
27. Elastic — 8.16 GA (Nov 12 2024) https://www.elastic.co/blog/whats-new-elastic-8-16-0 ; 9.0 GA (Apr 15 2025, Lucene 10); GitHub releases (9.4.x); https://endoflife.date/elasticsearch (as-of Jun 2026). blog/docs
28. BigDataBoutique — OpenSearch vs Elasticsearch compared (2026): https://bigdataboutique.com/blog/opensearch-vs-elasticsearch-compared — ES|QL vs PPL, ML Commons vs ELSER, vector engines. blog (consultancy)
29. tech-insider.org / knowi.com — feature/pricing tables (BBQ, ELSER, dims, free RBAC): https://tech-insider.org/elasticsearch-vs-opensearch-2026-2/ . blog
30. signoz.io / intellizu / oktopeak — "choose X if" tables, OTel-first vs Elastic Agent, free-vs-paid security: https://signoz.io/comparisons/elasticsearch-vs-opensearch/ . blog
31. jusdb.com — concise decision guide; "below 7.10 identical," reindex caveat: https://www.jusdb.com/compare/elasticsearch-vs-opensearch . blog
32. Elastic — performance-gap benchmark (2023, vendor, biased): https://www.elastic.co/blog/elasticsearch-opensearch-performance-gap — "40–140% faster." blog (vendor — flagged)
33. Trail of Bits (independent, disconfirming, Mar 6 2025): https://blog.trailofbits.com/2025/03/06/benchmarking-opensearch-and-elasticsearch/ — OpenSearch 2.17.1 ~1.6x faster Big5, ~11% faster vector vs ES 8.15.4; + blunders.io methodology critique. blog (independent)
34. OpenSearch — 3.0 GA (May 6 2025, Lucene 10): https://opensearch.org/blog/unveiling-opensearch-3-0/ ; version history https://docs.opensearch.org/latest/version-history/ (as-of Jun 2026; 3.7.0). docs/blog
