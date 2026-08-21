# JPMC — Vendor/Procurement Posture, MongoDB Signals & Earnings Tech Signals

> Provenance: reference under the `jpmorgan-chase-account-intelligence` skill. Built 2026-06-18 via `/dr` deep research (exa + WebSearch/WebFetch; primary sources incl. JPMC supplier-control docs). **verified-as-of: 2026-06-18.** The MongoDB section was deliberately stress-tested against over-claiming — read its characterization carefully.

## Contents
- Supplier base & procurement scale
- Third-party-risk posture (the proof burden)
- The April-2025 supplier open letter (BYOC/self-hosting)
- Named technology partnerships
- Procurement/sourcing behavior signals
- Public MongoDB / Atlas signals (carefully characterized)
- Recent earnings / Investor-Day technology signals
- TAM implications
- Sources

## Supplier base & procurement scale

JPMC runs a very large, formally governed supplier program: a **procurement budget >$23B/yr**, with **>$2B/yr spent with diverse suppliers** under a three-decade supplier-diversity program (Tier-1/Tier-2 structure).[^1] This implies rigorous, multi-stage vendor governance for any technology purchase.

## Third-party-risk posture (the proof burden)

JPMC publishes **Supplier Minimum Control Requirements (SMCR)** — a public, versioned (2020, 2025 updates) control catalog incorporated by reference into every supplier's Master Agreement.[^2] It is the concrete "proof burden" a data-platform deal faces. SMCR covers:
- **Data protection** — encrypt JPMC Confidential data **in transit and at rest** (including backups and when shared with sub-processors).
- **Technology/business resiliency** — approved **RTO/RPO and Maximum Tolerable Downtime**, DR recovery plans and **evidence of testing**.
- **Software-supply-chain security** — SBOM-style integrity/authenticity; open-source licensing/inventory.
- **Fourth-party / subcontractor** identification and monitoring.
- **Vulnerability management** and **incident/event management integrated into JPMC's own incident process.**

**For a TAM: this is the checklist your data-platform story must satisfy** — encryption, DR/RTO/RPO test evidence, fourth-party transparency, audit logging, vulnerability handling. (For how MongoDB/Atlas meets each, route to `mongodb-operations-expert` / `mongodb-atlas-expert`; this skill does not give product config.)

## The April-2025 supplier open letter (BYOC / self-hosting)

In **April 2025**, JPMC's CISO org published an unusually pointed **"open letter to our third-party suppliers."**[^3] Key, highly TAM-relevant points:
- JPMC's third-party providers suffered **"a number of incidents"** over the prior three years, forcing JPMC to **isolate compromised providers**.
- Suppliers must prioritize security **"equal to or above launching new products,"** and move beyond **"annual compliance checks"** to **continuous, demonstrable control evidence** and secure-by-default configurations.
- It explicitly names **confidential computing, customer self-hosting, and bring-your-own-cloud (BYOC)** as preferred controls for SaaS, and warns that AI / AI-agent services **"amplify"** supply-chain risk.

**Implication:** JPMC has a public appetite for **self-hosted / BYOC** deployment and skepticism of opaque multi-tenant SaaS — directly shaping how you'd position MongoDB Atlas (in JPMC's own cloud accounts, with private endpoints) versus self-managed MongoDB Enterprise Advanced.

## Named technology partnerships (public)

- **Core banking — Thought Machine (Vault Core):** Sept 2021 deal to replace the US retail core across CCB; JPMC is **also an investor** (Thought Machine's $200M Series C, 2021); inducted into JPMC's 2022 "Hall of Innovation"; multi-year migration. **Fact.**[^4]
- **Cloud — AWS, Microsoft Azure, Google Cloud (all three):** multi-cloud; AWS is the anchor (see `jpmc-cloud-data-modernization.md`). **Fact.**
- **Data/AI vendors — Databricks, Snowflake, MongoDB:** JPMC publicly named as a customer of all three by its Head of Technology Strategy (see MongoDB section). Internal platforms referenced publicly include **JADE** (JPMorgan Chase Advanced Data Ecosystem) and **Infinite AI**; **Apache Kafka/Confluent** and **Apache Iceberg** appear in architecture/job specs. **Fact (named as customer); qualified on internal-platform detail.**[^5]
- **Payments — Renovite (acquired Sept 2022):** cloud-native payments tech; had been a JPMC vendor since 2021; folded into J.P. Morgan Payments to build next-gen merchant acquiring (explicitly to counter Stripe/Block). **Fact.**[^6]

## Procurement/sourcing behavior signals

JPMC publicly displays a **"vendor → POC/load-test → acquire-if-strategic"** pattern: it ran **Renovite** as a vendor before buying it outright ("too important not to own"), and **performance-tested Thought Machine's Vault under simulated peak load** before scaling.[^4][^6] **Qualified.** *TAM relevance: expect rigorous POCs / load-testing and a build-vs-buy-vs-acquire calculus before any strategic vendor is scaled.*

## Public MongoDB / Atlas signals (carefully characterized)

**Bottom line:** There is a **confirmed public signal that JPMorgan Chase USES MongoDB (including MongoDB Atlas at internal-platform scale)** — corroborated by a senior-executive statement, a production AI-agent talk, and numerous JPMC job postings. **There is NO public MongoDB customer case study, press release, or commercial-relationship/"customer win" naming JPMorgan Chase.** Characterize the relationship as **"publicly-evidenced usage,"** not a publicly confirmed commercial/reference relationship. **Do not assert deal size, commercial terms, satisfaction, or strategic status — none are public.**

**Strongest signals (JPMC uses MongoDB):**
- **Senior-executive attribution — strongest source.** Constellation Research quotes **Larry Feinsmith (MD & Head of Technology Strategy, Innovation & Partnerships, JPMorganChase)** stating JPMC is **"a customer of Databricks, Snowflake and MongoDB [and] has multiple platforms."** This is the single most authoritative public statement that JPMC is a MongoDB customer. **Qualified→fact for "uses."**[^5]
- **Internal MongoDB Atlas platform team.** A JPMorganChase **"Lead Software Engineer – Cloud Databases"** posting describes a **"strategic MongoDB Atlas team"** that **"architect[s] and engineer[s] the control plane for our globally used MongoDB Atlas SaaS offering — provisioning, automation, scaling, integrations & self-service."** **This is the highest-value lead for a TAM** — it implies JPMC runs **Atlas as an internal, self-service database-as-a-service at global scale, with its own platform/SRE org around it.** **Strong, but single primary posting — treat as a lead to validate, not a confirmed account fact.**[^7]
- **Production AI-agent talk.** At LangChain "Interrupt" (May 2025), JPMorgan Chase **Private Bank** engineers described a multi-agent investment-research system ("Ask David") whose doc-search agent **"gets the data from MongoDB."** Confirms MongoDB in at least one production GenAI workload. **Qualified.**[^8]
- **Multiple JPMC job postings** require MongoDB skills across LOBs (Corporate, Infrastructure Platforms, CIB/Prime Finance) and geographies (US/UK/India) — on-prem replica sets & sharded clusters **and** Atlas (autoscaling, multi-region, private endpoints, RBAC, encryption, PITR, DR). **Fact that MongoDB/Atlas skills are broadly required.**[^9]

**Disconfirming checks (run deliberately):**
- **No JPMorgan Chase case study on mongodb.com.** MongoDB's FSI case-study library surfaces **Wells Fargo, Nationwide, Lombard Odier, Macquarie**, and an unnamed "major Canadian bank" — but **no JPMorgan Chase** case study, press release, or named reference. **Confirmed absence.**[^10]
- **No MongoDB earnings-call or press mention** of JPMorgan Chase as a customer surfaced.
- **Third-party tech-footprint trackers** (technologychecker.io, landbase) list "JP Morgan Chase" under MongoDB/Atlas — **inference-based aggregators; weak; do not lean on them.**

**Precise framing to use internally:** *"JPMC publicly uses MongoDB and Atlas (exec-confirmed + an internal Atlas platform team); there is no public case study or commercial-relationship detail. The internal Atlas control-plane team is a strong lead to validate against our own account records."*

## Recent earnings / Investor-Day technology signals (2025–2026)

- **2026 technology spend ≈ $19.8B**, up ~10% YoY ("nearly $20B," largest in the industry); 2025 ~$18B, 2024 ~$17B.[^11] (See `jpmc-tech-scale.md`.)
- **Total 2026 expenses projected ~$105B** (up ~$9B / +10% from ~$96B in 2025); ~a quarter of the increase is **tech and "tech-adjacent."** Dimon framed it as the **"most competitive banking environment in ≥20 years,"** citing PayPal/Stripe/Block.[^11][^12]
- **"Past peak modernization" — focus shift to application code + data.** CFO Barnum: JPMC has *"shifted focus from infrastructure modernization to modernizing the underlying application code and data"* to be **AI-ready.**[^13] **This is the most direct TAM signal: the next modernization wave is the application/data layer.**
- **AI adoption:** LLM Suite ~150,000 weekly users; ~1,000 AI use cases (see `jpmc-ai-strategy.md`).[^12]
- **"Vendor rationalization"** is an explicit, recurring efficiency lever **within a rising total budget** (a buying lens even as spend grows).[^5]

## TAM implications

- **Win the SMCR before you win the deal.** Encryption, DR/RTO/RPO evidence, fourth-party transparency, audit logging, vulnerability handling — have the proof package ready.
- **Bring a self-hosted/BYOC/in-their-cloud-account answer.** It matches JPMC's stated supplier preference.
- **Validate the Atlas lead internally.** The "strategic MongoDB Atlas team" posting is the single most actionable signal — check it against the real account before assuming scope.
- **Expect a rigorous POC and a build-vs-buy lens.** JPMC load-tests strategic vendors and rationalizes its vendor base even while growing spend.
- **Sequence: compliance/security/resiliency first, cost last.** Mirror JPMC's own decision order.

## Sources

1. JPMorganChase newsroom, supplier-diversity stories — >$23B procurement budget; >$2B diverse spend. (Newsroom, primary) — https://www.jpmorganchase.com/about/suppliers/supplier-diversity
2. JPMorganChase Supplier Minimum Control Requirements (PDF; 2020 & 2025 updates) — encryption, RTO/RPO/DR evidence, software-supply-chain, fourth-party, incident integration. (Primary) — https://www.jpmorganchase.com/about/suppliers
3. JPMorganChase Technology blog, "An open letter to our third-party suppliers" (April 2025) — supplier incidents; continuous control evidence; confidential computing / self-hosting / BYOC. (Blog/newsroom, primary) — https://www.jpmorganchase.com/about/technology/blog/open-letter-to-our-suppliers
4. American Banker, "JPMorgan Chase moving retail bank's core system to cloud" (2021) + BusinessWire, Thought Machine 2022 Hall of Innovation — Vault Core deal; JPMC investor; load-tested. (Press) — https://www.americanbanker.com/news/jpmorgan-chase-moving-retail-banks-core-system-to-cloud
5. Constellation Research, "JPMorgan Chase digital transformation, AI and data strategy" — Feinsmith: "customer of Databricks, Snowflake and MongoDB"; JADE/Infinite AI; vendor rationalization; ~$15.3B (2023) tech budget. (Press/analyst) — https://www.constellationr.com/insights/news/jpmorgan-chase-digital-transformation-ai-and-data-strategy-sets-generative-ai
6. JPMorganChase newsroom, "JPM to acquire Renovite Technologies" (2022) + CNBC (2022-09-12) — Renovite "vendor→acquire" pattern. (Newsroom + press) — https://www.jpmorganchase.com/ir/news/2022/jpm-to-acquire-renovite-technologies-inc
7. JPMorganChase, "Lead Software Engineer – Cloud Databases" posting (Built In London) — "strategic MongoDB Atlas team"; control plane for "globally used MongoDB Atlas SaaS offering." (Job posting, primary — single source) — https://builtinlondon.uk/job/lead-software-engineer-cloud-databases/4557472
8. LangChain "Interrupt" (YouTube, 2025) — JPMC Private Bank "Ask David" agent "gets the data from MongoDB." (Conference talk) — https://www.youtube.com/watch?v=yMalr0jiOAc
9. JPMorganChase MongoDB job postings (Dice, Built In, foundit) — on-prem replica sets/sharded clusters + Atlas (PITR/DR/private endpoints/RBAC) across LOBs/geos. (Job postings, primary) — https://www.dice.com/
10. MongoDB customer case studies (mongodb.com) — Wells Fargo, Nationwide, Lombard Odier, Macquarie, unnamed Canadian bank — **negation: no JPMorgan Chase case study.** (MongoDB site) — https://www.mongodb.com/solutions/customer-case-studies/wells-fargo
11. American Banker, "JPMorgan invests for new age of competition amid AI fears" (2026) — ~$19.8B tech; ~$105B expenses; ~$9B increase. (Press) — https://www.americanbanker.com/
12. JPMorgan Chase 2026 Bernstein Strategic Decisions Conference transcript (PDF) + CNBC (2026-04-06) — 1,000 use cases; ~150K weekly LLM users; "most competitive environment." (Transcript + press) — https://www.jpmorganchase.com/ir
13. JPMorgan Chase 2026 Company Update, full transcript (PDF) — ~$19.8B tech spend; "past peak modernization"; application code + data focus. (Investor/transcript, primary) — https://www.jpmorganchase.com/ir
