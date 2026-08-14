<!-- Provenance: reference under the `big-bank-IT` skill. Created 2026-06-18 via /dr deep research (Martin Fowler primary + AWS/Microsoft/Gartner + McKinsey + FCA/PRA primary + reputable press). Educational, vendor-neutral technical orientation — NOT advice. Disconfirming evidence on strangler-fig, automated translation, and transformation success rates preserved. -->

# Batch Windows & Mainframe/Legacy Modernization Patterns (and Legacy-Exit Risk)

`verified-as-of: 2026-06-18` (modernization tooling — automated COBOL→Java, the AI-migration vendor landscape — moves fast; re-verify tooling claims).

> **Educational, vendor-neutral orientation, NOT advice.** The core-banking *vendor landscape* ("what cores exist") is owned by `core-banking-and-sor-soe.md`; this reference covers batch + the modernization *strategy taxonomy* + strangler-fig/encapsulation + legacy-exit risk. The 24x7-vs-batch tension overlaps with `latency-payments-and-247.md` (touched there lightly; owned here).

## Contents

- [The batch window: what EOD batch is, and why it's shrinking](#the-batch-window-what-eod-batch-is-and-why-its-shrinking)
- [The modernization-strategy taxonomy — the "Rs"](#the-modernization-strategy-taxonomy--the-rs)
- [Incremental patterns: strangler fig & encapsulation](#incremental-patterns-strangler-fig--encapsulation)
- [Parallel run / progressive / coexistence modernization](#parallel-run--progressive--coexistence-modernization)
- [Legacy-exit RISK: why big-bang fails, and the coexistence tax](#legacy-exit-risk-why-big-bang-fails-and-the-coexistence-tax)
- [Disconfirming evidence](#disconfirming-evidence)
- [Seller takeaways](#seller-takeaways)
- [Sources](#sources)

## The batch window: what EOD batch is, and why it's shrinking

**What runs in a bank's overnight batch.** A bank's core traditionally orchestrates a tightly-sequenced "nightly batch cycle" — controlled on the mainframe by **JCL** job streams — that performs transaction posting (the *final post* to the ledger, vs. the intraday *memo post* that only adjusts available balance), overdraft processing, **interest accrual**, fee assessment, **statement generation**, regulatory reporting, GL updates, **reconciliation**, and production of **settlement files** (ACH/NACHA, Fedwire/SWIFT).[^1][^3] *(Confidence: fact — corroborated across an educational COBOL reference and vendor core docs.)*

**Why batch exists as a design choice (not just a flaw):** batch/settlement windows historically created *friction that doubled as a stabilizer* — netting, intraday float, and cut-offs gave banks time to source funding, catch errors, and reconcile before anything became irrevocable. "Payment files were validated in batch windows and settled the next business day. If something broke, you had hours to fix it before it really mattered."[^5][^6] *(Confidence: fact.)*

**Why the window is shrinking:** instant-payment schemes (SEPA Instant, FedNow, RTP, Pix, UPI) execute and settle in seconds, 24/7/365, with *no downstream window to catch a mistake*; always-on settlement removes the netting buffer. The **Bank of England's 2026 consultation** on extending RTGS/CHAPS toward near-24/7 notes participants' overnight **batch systems "would need to complete within two hours or otherwise run asynchronously"** — direct primary evidence that the batch window is now a regulated, actively-negotiated constraint.[^6][^8] *(Confidence: fact.)*

**Batch as a hard modernization constraint:** legacy cores are "optimized for batch, not event-driven or high-velocity streaming," and even where a payments engine runs continuously, *upstream/downstream systems reintroduce batch via posting and reconciliation cycles.* The pragmatic industry answer is usually **not** to rip out the core but to add a **real-time operational/gateway layer** in front of it (bidirectional ISO 20022 transformation, a "shielding layer" exposing account data via APIs) — explicitly an encapsulation/strangler move.[^9][^11] *(Confidence: fact for the constraint; qualified for the gateway prescription.)*

## The modernization-strategy taxonomy — the "Rs"

The "Rs" are a **portfolio decision model** — each application gets a path based on cost, risk, and value. Two lineages exist and are often conflated:[^13]
- **Gartner's modernization options (7):** *encapsulate, rehost, replatform, refactor, re-architect, rebuild, replace* (application modernization).[^15][^16]
- **The "6/7 Rs of cloud migration"** (AWS, echoed by Microsoft/IBM): *retain, retire, rehost, relocate, repurchase, replatform, refactor* (Microsoft uses *rebuild* for repurchase).[^17][^18] *(Confidence: fact.)*

The legacy/mainframe-relevant patterns:
- **Rehost ("lift and shift")** — redeploy "as is" without recompiling. For mainframes, concretely **hardware/middleware emulation** (emulate the legacy instruction set so legacy binaries run unchanged on Linux/x86 — IBM Z, Unisys ClearPath, IBM i; CICS/IMS/VSAM/DB2; COBOL/PL-I/RPG).[^20]
- **Replatform ("lift and reshape")** — new runtime, minimal code change. For mainframes, **recompiling COBOL/PL-I on x86** (AWS + Rocket/Micro Focus, TmaxSoft OpenFrame), keeping 3270 screens/JCL, often temporarily retaining DB2; data converted EBCDIC→ASCII.[^21][^22]
- **Refactor / re-architect** — change code to improve structure (refactor) or decompose the monolith into services (re-architect). For mainframes, **automated COBOL→Java translation** (AWS Transform via Blu Age; IBM watsonx Code Assistant for Z) "while preserving business logic," with automated equivalence testing.[^20][^24][^25]
- **Replace / repurchase ("drop and shop")** — eliminate the old app, buy a packaged/cloud-native product (the vendor landscape is in `core-banking-and-sor-soe.md`).[^14]
- **Retain, Retire, Encapsulate** — Retain = leave and revisit; Retire = decommission, preserve data; **Encapsulate** = "leave it running, wrap behind an API" (Gartner's first option).[^15][^16] *(Confidence: fact.)*

## Incremental patterns: strangler fig & encapsulation

**Strangler fig (Martin Fowler).** Named from rainforest figs that germinate in a host tree and gradually envelop it: *gradually build a new system around the edges of the old, moving behavior over until the legacy system is "strangled" and can be decommissioned.* Fowler updated the page (Aug 2024) and reframed it around four activities: understand desired outcomes; break the problem into parts; deliver the parts; change the organization to sustain it. The canonical mechanism is a **façade/proxy** that intercepts inbound requests and routes to legacy or new — initially mostly legacy, shifting incrementally. AWS frames the same as **transform → coexist → eliminate**, keeping the monolith for rollback.[^26][^27][^30] The headline rationale is **reduced risk** (value lands steadily; even a partial migration can yield strong ROI).[^27] *(Confidence: fact — Fowler primary + AWS/Microsoft independently.)*

**Event interception** is the underlying machinery — exploit existing integration points ("technical seams") to intercept state-changing events and route some to new components (even pre-commit DB triggers); **asset capture** migrates a subset of business assets, and Fowler stresses **bidirectional (reverse) migration** to de-risk.[^31][^32] *(Confidence: fact.)*

**Encapsulation (distinct from strangler).** Wrap the legacy system in a modern API layer (REST/GraphQL) *without changing the underlying code*; expose only what other systems need. It's the **lowest-risk, fastest-payoff** pattern and "the most underused" — but it **does not address technical debt**, carries the legacy maintenance burden indefinitely, and the API is only as scalable as what's behind it. The supporting machinery is the **façade** and **anti-corruption layer (ACL)** (translates between systems so the new design isn't contaminated by legacy semantics; Microsoft pairs strangler-fig with an ACL during coexistence).[^29][^35][^36] *(Confidence: fact.)*

## Parallel run / progressive / coexistence modernization

In core banking, the industry has converged: **phased and coexistence-led ("dual-core") approaches are now the dominant model**, especially for large/complex banks; big-bang replacement "is no longer the norm."[^38][^39] **Parallel run** keeps legacy and new running simultaneously with synchronized data until the new core is proven, providing a fallback (at the cost of duplicate systems); **phased migration** delivers segments (by product/customer/geography) sequentially; the common **progressive model** (open new accounts on the new core, migrate existing in waves) is explicitly the **strangler-fig pattern applied to core banking.**[^38][^40] A reported quality bar: at tier-1 banks, parallel-run for major product lines typically runs **60–90 days** and is treated as mandatory.[^41] *(Confidence: fact that coexistence dominates; the 60–90-day figure is single-source — qualified.)*

## Legacy-exit RISK: why big-bang fails, and the coexistence tax

**The "~30%" core-transformation success rate** (McKinsey): only ~30% of CBS transformations succeeded in fully migrating ledgers and products over the past decade — inside McKinsey's broader, repeatedly-replicated finding that **<30% of digital/organizational transformations succeed.**[^42][^43][^44] *(Confidence: fact — appears across multiple McKinsey publications; one institutional source repeated.)*

**TSB 2018 — the canonical cautionary tale.** In April 2018 TSB migrated ~5.2M customers to a new core (Proteo4UK) in a largely single-event ("big bang") weekend cutover. The data migrated, but the platform "immediately experienced technical failures"; internet/mobile banking were "almost unusable," all branches affected, fraud spiked, and BAU wasn't restored until December 2018.[^46][^47]
- **Slaughter & May's 262-page independent review (Nov 2019):** *"TSB did not give sufficient consideration to whether a largely single-event migration was the right choice"*; live proving "was not carried out at sufficient scale"; "the platform was not ready … and Sabis was not ready to operate the platform."[^47][^50]
- **Regulatory enforcement (Dec 2022):** FCA + PRA fined TSB a **total of £48.65m** (FCA £29.75m; PRA £18.9m) for operational-risk, governance, and **outsourcing** failures; TSB paid £32.7m in customer redress and ~£330m+ in total post-migration charges.[^52][^54] *(Note: TSB partly disputed the report, attributing much disruption to two data centres "configured inconsistently despite being specified to be identical" — a contradiction worth preserving.)*[^46] *(Confidence: fact — primary FCA/PRA notices + S&M report.)*

**Why banks now prefer progressive over rip-and-replace:** big-bang is "strongly discouraged for tier-1/2 institutions" and has "encountered regulatory pushback"; FCA/PRA operational-resilience expectations (PS21/3) effectively require the ability to **revert to the previous system state**, which means the legacy core must stay live throughout — hard to defend in a single-window cutover.[^39][^41] *(Confidence: fact / qualified.)*

**The coexistence / hybrid tax** — progressive modernization is not free: **data synchronization** during transition is "complex and error-prone"; a **people/skills split** between maintaining legacy and building new; **"adapter hell"** (each strangled service needs bespoke rerouting + its own rollback plan); **vendor lock-in** (mitigated by exit clauses, open APIs, multi-cloud mandates). McKinsey/industry name **data migration, regulatory compliance, and organizational change management** as the three most-underestimated risk areas.[^34][^40][^58] *(Confidence: fact for the categories.)*

## Disconfirming evidence

1. **Strangler-fig's signature failure mode — the "Permanent Hybrid" / "Pareto Stall."** Teams migrate easy peripheral modules first (illusion of progress), then stall as the core's dependency density rises and cost-to-migrate climbs while value-added falls — leaving a hybrid that combines "the rigidity of the legacy system with the distributed complexity of the new one." Mitigation: define decommission/"kill-switch" criteria *before* the first line of code, and budget for the disproportionate "final mile."[^59] Thoughtworks independently lists incomplete-migration→hybrid as a core risk.[^34]
2. **When strangler-fig does NOT apply** (Microsoft): when requests can't be intercepted; when you can't access/modify the legacy source; when the system is small enough that wholesale replacement is simpler; or when you must decommission quickly.[^29]
3. **Automated COBOL→Java refactoring reliability is genuinely contested.** A transpiler vendor's "testing paradox": comprehensive tests let you refactor with *any* tool (so you don't need the LLM), but most shops *lack* such tests — so they can't validate LLM output; and *"scale is the killer"* (at millions of LOC even a 1% error rate is tens of thousands of defects). IBM's own ICSE-2026 research reports its watsonx-Code-Assistant-for-Z pipeline at *median ~80% structural / >75% functional* quality — strong but explicitly **not 100%**, and weak on middleware like CICS. Net: automated refactoring is real and shipping, but "preserve business logic" is a vendor aspiration hinging on testing rigor; deterministic transpilation is a credible competing school.[^24][^25][^62][^63]
4. **Most successful real-time-payments programs *don't* replace the core at all** — "the teams making the most progress rarely start by replacing their core"; they add a real-time gateway/operational layer (encapsulation), implying many banks deliberately choose **not** to exit legacy.[^11]
5. **Transformation-failure base rates are partly an execution artifact** — McKinsey's ~30% triples (~79%) for organizations taking a rigorous, fully-implemented approach, cutting against fatalistic "core transformation is doomed" readings.[^43]

## Seller takeaways

- **Never propose a rip-and-replace.** The credible posture aligns with wrap/encapsulate/strangler — land value in the layer *around* the core, with a documented rollback story.
- **Speak to the batch constraint.** "We can sit in front of the core as a real-time layer without disturbing the EOD batch" is a resonant, low-fear message.
- **Know the TSB story.** It's the reference every bank board uses to justify caution; aligning with progressive/coexistence modernization signals you understand their risk posture.

## Sources

[^1]: Banking and Payment Systems in COBOL — https://datafield.dev/learning-cobol-programming/part-07/chapter-34/key-takeaways.html — educational — EOD batch lifecycle (memo vs final post; accrual, statements, reconciliation, GL) and JCL orchestration.
[^3]: Oracle Banking — Batch Processes (BOD/EOD) — https://docs.oracle.com/cd/E64763_01/html/CI/CI10_Batch.htm — vendor product doc — concrete EOD/BOD job list (accruals, statements, profit posting).
[^5]: Always-On Settlement in Banking — https://www.dunnixer.com/insights/information/banking/us/always-on-settlement-in-banking — analytical — batch/settlement windows as stabilizer; 24/7 removes maintenance windows & netting buffer.
[^6]: Real-Time Rails Changed the Physics of Payments (Mambu) — https://mambu.com/en/insights/articles/real-time-rails-changed-the-physics-of-payments — vendor — SCT Inst/FedNow/Pix/UPI 24/7; batch window as error-catch buffer; core needs downtime for EOD.
[^8]: Extending RTGS and CHAPS settlement hours (Bank of England CP, 2026) — https://www.bankofengland.co.uk/paper/2026/cp/extending-rtgs-and-chaps-settlement-hours-next-steps — central bank (primary) — overnight batch must complete in ~2 hrs or run async; settlement-hours vs change-windows trade-off.
[^9]: Real-Time Payments Adoption Challenges — https://www.dunnixer.com/insights/information/banking/us/real-time-payments-adoption-challenges-banks — analytical — RTP as operating-model shift; upstream/downstream reintroduce batch.
[^11]: Why Real-Time Payments Break Batch-Era Architectures (GridGain) — https://www.gridgain.com/resources/blog/real-time-payments-digital-wallets-architectures — vendor — real-time gateway/shielding layer in front of mainframe; "rarely start by replacing the core."
[^13]: Legacy Modernization Approaches: The 7 R's (Saigon Technology) — https://saigontechnology.com/blog/legacy-modernization-approaches-7-rs/ — trade — distinguishes Gartner app-modernization Rs from AWS/IBM cloud-migration Rs.
[^14]: Use the 7 Rs to develop an app modernization strategy (TechTarget) — https://www.techtarget.com/searchcloudcomputing/tip/Use-the-7-Rs-to-develop-an-app-modernization-strategy — trade — AWS-expanded Rs; repurchase economics & data import.
[^15]: Gartner 7-options (summary) — https://www.gartner.com/smarterwithgartner/7-options-to-modernize-legacy-systems — analyst (landing) — encapsulate/rehost/replatform/refactor/rebuild/replace; encapsulate = wrap behind API.
[^16]: The 6 Rs of application modernization (Microsoft Learn) — https://learn.microsoft.com/en-us/azure/app-modernization-guidance/plan/the-6-rs-of-application-modernization — vendor authoritative — rehost/replatform/refactor/rebuild/retire/retain; when to refactor/rebuild.
[^17]: AWS Mainframe Modernization — https://aws.amazon.com/mainframe-modernization/ — vendor — replatform (Micro Focus/Rocket) vs refactor (AWS Transform/agentic AI COBOL→Java).
[^18]: Microsoft 6 Rs (as above) — included in [^16].
[^20]: Demystifying Legacy Migration Options to AWS — https://aws.amazon.com/blogs/apn/demystifying-legacy-migration-options-to-the-aws-cloud/ — vendor — hardware/middleware emulation = rehost; automated-refactoring scope (IBM Z/Unisys/IBM i; CICS/IMS/VSAM/DB2).
[^21]: AWS Mainframe Modernization — Replatform (Prescriptive Guidance) — https://docs.aws.amazon.com/prescriptive-guidance/latest/replatform-mainframe-apps-shared-db2/mainframe-modernization.html — vendor — recompile COBOL/PL-I, temporarily retain DB2.
[^22]: TmaxSoft OpenFrame replatforming on AWS — https://aws.amazon.com/blogs/apn/how-to-succeed-at-large-scale-mainframe-replatforming-with-tmaxsoft-openframe-on-aws/ — vendor (partner) — recompile on x86, EBCDIC→ASCII.
[^24]: Enterprise-Scale COBOL-to-Java Translation (IBM Research / ICSE 2026) — https://research.ibm.com/publications/enterprise-scale-cobol-to-java-translation-llms-augmented-with-program-analysis — peer-reviewed/vendor research — WCA for Z pipeline; median >80% structural / >75% functional.
[^25]: Quality Evaluation of COBOL to Java Code Transformation (arXiv, IBM) — https://arxiv.org/html/2507.23356v1 — peer-reviewed — LLM-as-judge limits; weak on CICS; "not a silver bullet."
[^26]: Strangler Fig — Martin Fowler (2024 rewrite) — https://martinfowler.com/bliki/StranglerFigApplication.html — canonical primary — origin metaphor; gradual replacement; four activities.
[^27]: Original Strangler Fig Application — Martin Fowler — https://martinfowler.com/bliki/OriginalStranglerFigApplication.html — primary — reduced-risk rationale; Event Interception as fundamental strategy.
[^29]: Strangler Fig Pattern (Microsoft Azure Architecture Center) — https://learn.microsoft.com/en-us/azure/architecture/patterns/strangler-fig — vendor authoritative — façade/proxy; anti-corruption layer; explicit "when NOT to use."
[^30]: Strangler fig pattern (AWS Prescriptive Guidance) — https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-decomposing-monoliths/strangler-fig.html — vendor — transform/coexist/eliminate; keep monolith for rollback.
[^31]: Event Interception — Martin Fowler — https://martinfowler.com/articles/patterns-legacy-displacement/event-interception.html — primary — technical seams; route state-changes to new components.
[^32]: Asset Capture — Martin Fowler — https://martinfowler.com/bliki/AssetCapture.html — primary — migrate subset of assets; bidirectional/reverse migration to de-risk.
[^34]: Embracing the Strangler Fig pattern (Thoughtworks) — https://www.thoughtworks.com/en-us/insights/articles/embracing-strangler-fig-pattern-legacy-modernization-part-one — trade — challenge list: interdependencies, data sync, incomplete-migration hybrid.
[^35]: Legacy System Modernization Approaches (AltexSoft) — https://www.altexsoft.com/blog/legacy-system-modernization-approaches/ — trade — encapsulation defined; indefinite legacy maintenance burden.
[^36]: Application Modernization Strategies (ArielSoftwares) — https://www.arielsoftwares.com/application-modernization-strategies/ — trade — encapsulate as lowest-risk/fastest-payoff & "most underused."
[^38]: Core banking migration strategies (10x Banking) — https://www.10xbanking.com/insights/core-banking-migration-strategies-choosing-the-right-path-to-a-4th-generation-platform — vendor — phased/coexistence now dominant; big-bang no longer norm; parallel-run & phased defs.
[^39]: 10 Key Areas for Successful Core Banking Modernization (Oliver Wyman) — https://www.oliverwyman.com/our-expertise/insights/2025/may/next-gen-core-banking-modernization.html — analyst — "dual core" coexistence preferred; big-bang higher risk & regulatory pushback.
[^40]: Core Banking in the Cloud (Finantrix) — https://www.finantrix.com/articles/core-banking-in-the-cloud — trade — progressive migration = strangler fig; vendor-lock-in mitigations; 3 underestimated risks.
[^41]: Core Banking Modernization Roadmap (Scrums.com) — https://www.scrums.com/guides/complete-guide-to-legacy-software-modernization — trade — TSB cautionary tale; PS21/3 revert requirement; 60–90-day parallel-run at tier-1.
[^42]: How to get a core banking transformation right (McKinsey) — https://www.mckinsey.com/capabilities/tech-and-ai/our-insights/tech-forward/how-to-get-a-core-banking-transformation-right-eight-mistakes-to-avoid — analyst — "only ~30% of CBS transformations succeeded."
[^43]: Losing from day one: successful transformations (McKinsey) — https://www.mckinsey.com/capabilities/people-and-organizational-performance/our-insights/successful-transformations — analyst — 30% success "hasn't budged"; rigorous approach → 79%.
[^44]: Why most digital banking transformations fail (McKinsey) — https://www.mckinsey.com/capabilities/tech-and-ai/our-insights/tech-forward/why-most-digital-banking-transformations-fail-and-how-to-flip-the-odds — analyst — only 30% of digital-banking transformations succeed.
[^46]: TSB Board publishes independent review of 2018 IT Migration (TSB) — https://www.tsb.co.uk/news-releases/slaughter-and-may.html — primary (company) — staged-then-"main migration event"; TSB's partial disagreement (data-centre misconfiguration).
[^47]: TSB programme pulled apart in report on IT meltdown (Computer Weekly) — https://www.computerweekly.com/news/252474170/TSB-programme-pulled-apart-in-report-on-IT-meltdown — press — "platform not ready … Sabis not ready"; testing too small-scale.
[^50]: TSB neglected to assess capabilities of main IT provider (Computer Weekly) — https://www.computerweekly.com/news/252474751/TSB-neglected-to-assess-capabilities-of-main-IT-provider-in-its-failed-system-migration — press — 262-page report; "big bang" without understanding risks.
[^52]: TSB fined £48.65m for operational resilience failings (FCA) — https://www.fca.org.uk/news/press-releases/tsb-fined-48m-operational-resilience-failings — regulator (primary) — total £48.65m; £32.7m redress; BAU by Dec 2018.
[^54]: PRA Final Notice to TSB (Bank of England) — https://www.bankofengland.co.uk/-/media/boe/files/prudential-regulation/regulatory-action/final-notice-from-pra-to-tsb-bank.pdf — regulator (primary) — PRA £18.9m; breach of Fundamental Rules 2 & 6; undue reliance on intragroup provider.
[^58]: Pros and cons of the Strangler architecture pattern (Red Hat) — https://www.redhat.com/en/blog/pros-and-cons-strangler-architecture-pattern — vendor — "adapter hell"; ongoing routing/network management; per-instance rollback.
[^59]: The Permanent Hybrid Trap (Moonello) — https://www.moonello.com/insights/the-strangler-that-never-strangles-preventing-the-permanent-hybrid-trap — trade — Pareto Stall; define "kill switch" before coding.
[^62]: LLMs and COBOL Modernization: Where Transpilation Fits (Heirloom Computing) — https://heirloomcomputing.com/llms-and-cobol-modernization/ — vendor (opinionated) — the "testing paradox"; "scale is the killer"; deterministic transpilation scales linearly.
[^63]: Accelerate mainframe modernization with AWS Transform (AWS) — https://aws.amazon.com/blogs/migration-and-modernization/accelerate-mainframe-modernization-with-aws-transform-a-comprehensive-refactor-approach/ — vendor — automated COBOL→Java "functional equivalence"; iterative validation can extend timelines.
