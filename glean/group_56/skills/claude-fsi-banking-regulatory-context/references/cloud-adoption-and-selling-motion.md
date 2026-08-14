<!-- Provenance: reference under the `fsi-banking-regulatory-context` skill. Created 2026-06-18 via /dr deep research (gov + analyst + vendor + practitioner sources). Educational seller orientation — NOT legal/compliance advice. Disconfirming evidence on cloud depth/repatriation is preserved. -->

# FSI Cloud-Adoption Posture & Selling Into a Regulated Industry

`verified-as-of: 2026-06-18` (cloud-adoption figures, sovereign-cloud, and the regulatory landscape are fast-moving — re-verify dated specifics).

> **Educational orientation, not legal/compliance advice.** Part A = where FSI cloud adoption stands; Part B = the seller's reality (proof burden, references, security expectations, build-vs-buy). Treat hyperscaler/vendor playbooks as vendor-tier (factual but interested). TPRM/DORA *mechanics* → `enterprise-vendor-management-and-tprm`; buyer-psychology mechanics → `applied-psychology`.

## Contents

- [PART A — FSI cloud-adoption posture](#part-a--fsi-cloud-adoption-posture-as-of-2026)
- [PART B — Selling into a regulated industry](#part-b--selling-into-a-regulated-industry)
- [Sources](#sources)

## PART A — FSI cloud-adoption posture (as of 2026)

### A1. Near-universal *presence*, shallow *depth*

The headline "90%+ of banks use cloud" is true but misleading: there are two different numbers, and conflating them is the most common seller error.
- **Breadth (has *some* cloud): ~90%+.** A US Treasury report (citing a 2021 ABA survey) found >90% of surveyed banks keep at least some data/apps/operations in cloud, but >80% described themselves as in "adoption" or "early adoption," and **only ~5% called their cloud use "mature."**[^1]
- **Depth (share of workloads actually *in* cloud): low.** The same Treasury-cited survey estimated **~8% of all banking workloads** were cloud-based; Accenture's 2025 index of 78 large banks found **only ~10% of *core* workloads** had moved.[^1][^4] A vendor-cited "59% of workloads" figure is an outlier — treat as tentative/optimistic.[^7]

The ambition-reality gap is real: fewer than ~40% of FS executives are highly satisfied with cloud outcomes, with **lift-and-shift** (porting the monolith without re-architecting) a common culprit.[^2][^3][^8] *(Confidence: fact — multi-source, gov + analyst.)*

> **Seller implication.** Don't pitch to "the bank's cloud maturity"; pitch to a *specific workload's* placement decision. Most workloads are still on-prem and many migrated ones underdeliver. Win on *this* workload, not a cloud-transformation thesis.

### A2. Hybrid + multi-cloud is the durable norm

Hybrid/multi-cloud dominates as strategy: surveys put **~82% of FS firms on multi-cloud or hybrid** (one split it ~41% multi-cloud / ~41% hybrid).[^9][^10] But **true** multi-cloud (workloads genuinely spread across >1 hyperscaler) is rarer: Accenture found only ~31% of banks have adopted a multi-cloud strategy.[^4] Many "hybrid" estates are really *containers on VMs spanning on-prem + one public cloud + edge*: permanent coexistence, not a transition state.[^12] *(Confidence: fact that hybrid is the steady state; qualify the multi-cloud share — Accenture 31% vs LSEG 41% reflect aspiration vs operational reality.)*

> **Seller implication.** "Runs the same everywhere" (portability across on-prem + multiple clouds) is a first-class buying criterion in FSI, reinforced by regulator exit-strategy demands (A3).

### A3. Concentration risk: regulators are actively driving portability / exit-strategy demand

The fastest-moving area and a genuine tailwind for portability-positioned vendors. Regulators worry that FS dependence on a handful of hyperscalers creates **single points of failure with no substitutability**: the UK stood up a statutory **Critical Third Parties (CTP) regime** (effective 1 Jan 2025) letting authorities designate and directly oversee systemically critical providers, and under **DORA** (in force 17 Jan 2025) EU regulators designate critical ICT third parties.[^13][^14][^18] ~72% of cloud services in Europe are provided by US firms.[^17] The regulatory output a seller will encounter: demands for **documented, *tested* exit strategies**, **concentration-risk analysis** (incl. "hidden concentration" where several SaaS vendors all run on the same underlying cloud), and explicit **portability/interoperability** expectations.[^19][^21][^22] Honest nuance: industry bodies note *mandating* portability is technically fraught and that portability is "typically not fast enough to meet recovery targets," an exit tool rather than a failover tool.[^20] *(Confidence: fact, gov primary sources, that regulators drive exit/portability/multi-cloud demand; fast-moving.)*

> **Seller implication.** A regulator-driven "can you leave?" question now sits inside technical evaluations. Standard data formats, containerized/abstracted deployment, and run-anywhere portability *de-risk the buyer's DORA/CTP exposure*, so position a credible, documentable exit/egress path as a compliance asset. (Detailed TPRM/DORA mechanics → `enterprise-vendor-management-and-tprm`.)

### A4. Sovereign cloud (the adoption consequence)

Residency *drivers* are in `how-regulation-drives-buying-behavior.md`; the adoption consequence is what matters here. All three hyperscalers launched **sovereign offerings in 2025-2026** (e.g., the AWS European Sovereign Cloud reached GA in early 2026 as a separate EU legal entity with dedicated IAM/billing; Microsoft and Oracle fielded sovereign clouds + external key management).[^23][^24][^25] A key buyer-sophistication point: critics note vendors don't always explain *how* an EU-board/customer-managed-key arrangement reduces the risk of foreign-government access under the CLOUD Act.[^17] Sovereign regions are typically a **smaller service catalog at a premium.** *(Confidence: fact on the offerings' existence, vendor-tier; re-verify specifics.)*

> **Seller implication.** A portable product that *runs inside* a sovereign region or the customer's own DC (rather than depending on a hyperscaler-proprietary managed service that may not exist there) is differentiated. Simplifying the residency story accelerates a deal stuck in "fear of messing up" paralysis.

### A5. Repatriation / cost-control counter-trend (disconfirming evidence)

There is a genuine counter-current to "all-in on cloud," but it's narrow and easy to overstate in both directions. It's **real but selective**: surveys suggest ~80% of organizations expect *some* repatriation of *some* workloads within a year, but **only ~8-9% intend *full* workload repatriation**. Most move *select* elements (production data, backups, specific compute), not wholesale exits.[^27][^28][^29] Drivers: unpredictable steady-state cost at scale, data security/sovereignty, and **AI-infrastructure economics**; FS and healthcare lead repatriation due to scale + regulation.[^31][^32][^33] A strong skeptic counter-voice argues much "repatriation" is *optimization theater*, and that the real cost is now running two environments (a hybrid complexity tax).[^34][^35] *(Confidence: fact that repatriation is selective/hybrid, not a cloud reversal; magnitude is contested — vendor surveys skew high.)*

> **Seller implication.** The honest framing is **"right workload, right place,"** not cloud-vs-on-prem ideology. A product that runs *anywhere* lets the customer make placement decisions without re-platforming: the hedge both the repatriation crowd and the cloud-first crowd want. Don't argue the bank should be all-in on cloud; that thesis is empirically weak.

### A6. Mainframe / core-banking persistence

The core didn't move and largely still hasn't: analyst/vendor sources indicate **>40% of banks still run COBOL-based cores and 45 of the top 50 banks still rely on mainframe** for mission-critical operations; back-office/core (core banking, risk & compliance) is the *slowest* to migrate.[^36][^37][^38] Why: core replacement runs **$100-500M over 3-7 years with substantial failure risk**, large banks spend 70-85% of IT budgets just maintaining legacy, COBOL talent is retiring, and these systems "cannot fail."[^36][^39] *(Confidence: fact — analyst + vendor convergent; figures well-triangulated.)*

> **Seller implication.** The migratable opportunity is the *satellite/digital/data layer around the core*, not the ledger. Target net-new digital workloads, real-time data, and the modernization "wrap": the "composable core" pattern (B4).

## PART B — Selling into a regulated industry

### B1. The proof burden: evidence, not assertions

In FSI the security/compliance review is often the *real* sales gate, run by a third-party-risk team of trained auditors, not the line-of-business buyer.[^41] The standard proof artifacts a TAM/pre-sales will be asked for:
- **Security questionnaires:** **SIG** (Shared Assessments — SIG Lite ~150 Qs / SIG Full 1,000+ Qs; *favored by financial services*) and **CAIQ** (Cloud Security Alliance, ~280 Qs, maps to the Cloud Controls Matrix); they largely ask the same underlying things in different orders.[^42][^43]
- **Certifications/attestations:** **SOC 2 Type II** (the baseline, which must cover the in-scope product over ≥6-12 months, with a **bridge letter** to the present), **ISO 27001**, **PCI-DSS v4.0.1** (contractual, non-negotiable if card data is processed), **FedRAMP** for US public sector (independent 3PAO assessment + annual pen test).[^44][^45][^46][^47][^48]
- **Penetration-test evidence:** banks expect summaries/results; contracts increasingly grant **customer-initiated pen-test rights.**[^49]
- **Right-to-audit clauses & regulatory-grade SLAs:** audit & access rights for the firm *and* its regulators, defined RTO/RPO, incident-notification windows (often 2-24 hours), and on-site assessment rights, explicitly "not a questionnaire."[^49][^50]

> **"Is SOC 2 enough?": No.** SOC 2 is the *floor*, not the finish line. Banks' TPRM teams require supplementary documentation (pen-test results, tested BC/DR and IR plans, insurance certificates, their proprietary questionnaire), broader Trust Service Criteria (Availability + Processing Integrity effectively mandatory), and SOC 2 says nothing about **model risk** (SR 11-7 / the 2026 SR 26-2 successor), contract terms, or whether the bank should buy this category at all. OCC guidance is explicit that banks must judge sufficiency themselves and scale diligence to risk.[^44][^45][^52] Never answer a model-risk, residency, or resilience question by "pointing at the SOC 2."[^52] *(The full security-questionnaire + attestation gauntlet, from the buyer-governance side, is `enterprise-vendor-management-and-tprm`.)*

> **Seller implication.** Build the **evidence bank *before* the questionnaire arrives**: a trust center with SOC 2 Type II (current bridge letter), ISO 27001, PCI scope, pen-test summary, BC/DR test evidence, IR plan, sub-processor list, and a pre-filled SIG/CAIQ. Proactively offering these compresses the risk-committee cycle.

### B2. References & risk-transfer: "who else like me uses this"

In risk-averse FSI, **trust is the primary currency** and a *named, in-industry* reference outweighs marketing.[^54][^55] In-industry *specificity* matters more than logo count: "used by $2B community banks in the Southeast" beats "trusted by 50+ banks"; peer validation must match asset size and geography.[^56] References function as **risk-transfer**, since a bank wants proof a *peer regulated firm* already absorbed that risk successfully; a single tier-one reference "unlocks dozens of conversations."[^55][^57] **Borrowed credibility** works when direct references aren't available: partnering with core providers (FIS, Fiserv, Jack Henry) who refer you, or citing analyst validation (Gartner et al.).[^58] *(Confidence: fact — convergent across GTM sources; a16z "B2FI" is canonical.)*

> **Seller implication.** Land *one* lighthouse FSI reference early and invest heavily in its referenceability; it is the single highest-impact asset in the motion. Match the reference to the prospect's segment, not just the industry. As a TAM, your existing customers' referenceability is a strategic account asset, so cultivate it deliberately.

### B3. Regulatory-grade security expectations (buyer-expectation level)

What a bank's security team will demand of *any* vendor handling regulated data (expectation level; vendor-specific config → `mongodb-compliance`):
- **Encryption at rest, in transit, and increasingly *in use*.** Baseline AES-256 at rest, TLS 1.3 in transit; backups encrypted with segregated keys. DORA explicitly mandates encryption at rest, in transit, and where necessary in use.[^59][^60][^61]
- **Key management / BYOK, and the BYOK-vs-customer-controlled distinction (sophisticated buyers know it).** Banks want **customer-managed keys in FIPS 140-2/3 HSMs** with independent rotate/revoke/audit. Regulators (EBA/BaFin) and security teams distinguish **BYOK** (provider can still technically decrypt → compellable under foreign disclosure law) from **customer-*controlled*/HYOK keys** (provider only ever sees ciphertext).[^59][^62][^63][^64]
- **Isolation / tenancy:** clear tenant-separation architecture, private-connectivity options, multi-tenant risk controls.[^59][^60]
- **Resilience / DR:** documented and *tested* RTO/RPO, annual live failover tests, BC/DR evidence; vendor resistance to providing this is a red flag.[^49][^50]
- **Data-deletion guarantees:** secure deletion / **crypto-shredding**; return *or* deletion of data on termination in documented formats; key destruction as the mechanism to render data inaccessible across backups/replicas (supports right-to-erasure).[^60][^63][^64]

> **Seller implication.** Be fluent in the **BYOK vs. customer-controlled-key** distinction. Getting it wrong in front of a bank security team is disqualifying; getting it right signals you understand *their* regulatory exposure. Bring a tenancy diagram, encryption spec (algorithms + TLS version), and DR test results to technical reviews.

### B4. Build-vs-buy at banks: the "is this differentiating?" lens

Large banks historically **built** for maximum control/customization (a custom core can exceed $50M and take years), viable only with substantial resources.[^66][^67] The decisive 2025-2026 lens is now **"is this core/differentiating?"**, and the dominant pattern is the **"Composable Core": buy commodity, build differentiation.** COTS for non-differentiating back-office (general ledger, KYC, regulatory reporting); custom microservices for customer-facing experience, AI, and proprietary products.[^68][^69][^70] What banks are *not* building: core systems from scratch and foundational AI models, because "the real strategic questions sit higher in the stack."[^71] AI is making the decision *more* complex (lowering build cost while raising data's strategic value), and practical constraints (fragmented data, talent gaps) further limit internal build.[^71] *(Confidence: fact — convergent analyst + vendor; "composable/hybrid" corroborated by neutral sources.)*

> **Seller implication.** Win by positioning your product as the **commodity/foundational layer the bank should *not* waste scarce engineers building**, while explicitly *enabling* the differentiation they build on top (APIs, data platform, flexibility). Frame against the *internal build alternative*, not just competitors. A platform that is *both* a reliable bought commodity *and* a build-on-able foundation hits the composable sweet spot.

### B5. The multi-stakeholder, long-cycle sale, and where a TAM/pre-sales fits

*(Behavioral drivers — risk-aversion, auditability, SoD, long cycles — are in `how-regulation-drives-buying-behavior.md`; buyer-psychology mechanics → `applied-psychology`. Here, just the GTM shape and the role.)* Structured RFPs with scoring rubrics; **5-10 stakeholders** across business, technology, risk, compliance, procurement; **9-18 month** cycles; budgets swayed by exam findings. Upside: once in, banks renew for years.[^54][^55] The technical/pre-sales role (TAM, Solutions Architect) is a **client-facing technical partner who translates the product into the bank's "decision system,"** drives adoption/value realization, works *with the bank's security team to match its standards*, and acts as trusted advisor rather than vendor.[^72][^73][^75] "No one owns the problem internally" kills more FSI deals than product gaps.[^75]

> **Seller implication.** A TAM's three highest-impact contributions in FSI are (1) being the proof-and-trust orchestrator (B1/B3), (2) owning referenceability (B2), and (3) translating product value into language risk/compliance/procurement can defend internally. For the POC/POV *execution mechanics*, see `tam-operations` (enterprise POC/POV motion).

## Sources

[^1]: https://home.treasury.gov/system/files/136/Treasury-Cloud-Report.pdf — US Treasury: FS sector adoption of cloud services (gov; ~8% workloads, 5% mature).
[^2]: https://www.capgemini.com/insights/research-library/world-cloud-report-financial-services/ — Capgemini World Cloud Report FS 2025 (analyst; ambition-reality gap, 12% innovators).
[^3]: https://www.sogeti.us/research-and-insight/world-cloud-report-financial-services-2025/ — Sogeti/Capgemini World Cloud Report FS 2025 (analyst).
[^4]: https://bankingblog.accenture.com/cloud-ai-banking-growth — Accenture 2025 Rotation Index (analyst; 10% core, 31% multi-cloud).
[^7]: https://www.technavio.com/report/private-and-public-cloud-market-in-the-financial-services-industry-analysis — Technavio FSI cloud market (analyst; "59%" outlier — tentative).
[^8]: https://www.ccgcatalyst.com/thought-leadership/commentary/the-great-rebuild-core-modernization-data-infrastructure-and-cloud/ — CCG Catalyst: The Great Rebuild (trade-press; "changed the road not the vehicle").
[^9]: https://www.lseg.com/en/media-centre/press-releases/2025 — LSEG Global Cloud Survey (analyst/vendor; 82% hybrid/multi-cloud).
[^10]: https://www.lseg.com/en/insights/data-analytics — LSEG: cloud strategies in FS (analyst/vendor; 41/41 multi/hybrid split).
[^12]: https://www.nutanix.com/blog — Nutanix Enterprise Cloud Index: FS (vendor; "hybrid reality… permanently").
[^13]: https://www.bankofengland.co.uk/prudential-regulation/publication/2024/november/operational-resilience-critical-third-parties-to-the-uk-financial-sector-policy-statement — BoE PS16/24 Critical Third Parties (gov; statutory CTP regime, effective 1 Jan 2025).
[^14]: https://www.bankofengland.co.uk/prudential-regulation/publication/2024/november/operational-resilience-critical-third-parties-to-the-uk-financial-sector-supervisory-statement — BoE/PRA/FCA SS6/24 (gov).
[^17]: https://www.thebanker.com/ — "Cloud conundrum leaves banks with FOMU" (trade-press; 72% EU cloud = US firms; sovereignty paralysis).
[^18]: https://ecipe.org/publications/cloud-resilience-and-security/ — ECIPE: Cloud Resilience and Security (think-tank; designated critical ICT providers; lock-in).
[^19]: https://www.dlapiper.com/en/insights/publications/2024/06/new-ecb-guidelines-on-outsourcing-cloud-services — DLA Piper on ECB cloud-outsourcing guidance (law-firm; stressed-exit/insolvency BCP testing).
[^20]: https://www.afme.eu/publications/position-papers — GFMA/AFME public-cloud portability (industry body; portability ≠ failover — disconfirming nuance).
[^21]: https://regtechanalyst.com/ — How DORA redefines ICT exit planning (trade-press; weak exit = resilience problem).
[^22]: https://www.regulation-dora.eu/blog — Cloud exit strategies & concentration risk under DORA (trade-press; hidden/geographic concentration).
[^23]: https://aws.amazon.com/blogs/aws/opening-the-aws-european-sovereign-cloud/ — AWS European Sovereign Cloud (vendor; dedicated EU IAM/billing). GA in early 2026 corroborated by trade press.
[^24]: https://pretius.com/blog/oracle-sovereign-cloud — Oracle Sovereign Cloud overview (trade-press/vendor; AWS ESC GA, premium, service-catalog size).
[^25]: https://blogs.microsoft.com/blog/2025/06/16/announcing-comprehensive-sovereign-solutions-empowering-european-organizations/ — Microsoft sovereign solutions (vendor; EU-controlled access, external key management).
[^27]: https://www.crn.com/news/networking/2025/f5-ceo-ai-is-sparking-cloud-repatriation-movement — F5 CEO on repatriation (trade-press; IDC ~80% expect some repatriation; AI driver).
[^28]: https://biztechmagazine.com/article/2025/08/why-some-workloads-are-coming-home-case-cloud-repatriation — Why some workloads come home (trade-press/IDC; only 8-9% full repatriation).
[^29]: https://www.puppet.com/blog/cloud-repatriation — Cloud repatriation 2025 (vendor/IDC; selective not wholesale).
[^31]: https://www.crn.com/news/networking/2025/f5-ceo-ai-is-sparking-cloud-repatriation-movement — (as [^27]; AI-infrastructure driver).
[^32]: https://www.opentext.com/ — The cloud-repatriation shift (vendor; security 51% / cost 39% drivers).
[^33]: https://digitaldigest.com/cloud-repatriation-strategy-2025/ — Why companies bring data home (trade-press; FS + healthcare lead).
[^34]: https://www.lastweekinaws.com/blog/cloud-repatriation-is-getting-complicated/ — Corey Quinn (practitioner skeptic; "optimization theater," hybrid complexity tax).
[^35]: https://www.cloudmagazin.com/en/ — Cloud repatriation: when it makes sense (trade-press; egress/skills/capex hidden costs).
[^36]: https://aws.amazon.com/blogs/industries/modernizing-core-banking-systems-a-strategic-guide-for-financial-leaders/ — AWS: modernizing core banking (vendor; >40% COBOL cores, 45/50 mainframe, 70-85% IT budget on legacy).
[^37]: https://www.rocketsoftware.com/ — IDC/Rocket: mainframe in FS (analyst, vendor-sponsored; 20% banks ≥50% workloads on mainframe through 2025).
[^38]: https://www.fstech.co.uk/ — Celent: hybrid multi-cloud world (analyst, vendor-sponsored; back-office/core slowest).
[^39]: https://www.softwareseni.com/ — Database dynasties & language longevity (trade-press; core replace $100-500M/3-7yr, 0.01% error catastrophic).
[^41]: https://blog.getagency.com/articles/what-enterprise-banks-expect-fintech-soc-2 — What banks expect from your SOC 2 (trade-press/vendor; TPRM = trained auditors; SOC 2 insufficient).
[^42]: https://safe.security/resources/blog/ — Vendor security questionnaire best practices (vendor; SIG vs CAIQ, question counts).
[^43]: https://www.sharedassessments.org/ — Shared Assessments SIG (standards/industry; SIG Lite vs Full).
[^44]: https://blog.getagency.com/articles/what-enterprise-banks-expect-fintech-soc-2 — (as [^41]; supplementary docs beyond SOC 2, Availability + PI mandatory).
[^45]: https://www.myabt.com/blog/ — Tech due diligence on fintech vendors (trade-press; SOC 2 Type II ≥6mo + bridge letter; OCC lifecycle).
[^46]: https://blog.getagency.com/articles/soc-2-for-payment-processing-companies — SOC 2 for payment processing (trade-press/vendor; PCI ≠ SOC 2 substitute).
[^47]: https://fedramp.gov/ — FedRAMP (gov; Rev5 agency path, 3PAO).
[^48]: https://www.fedramp.gov/resources/documents/ — FedRAMP penetration-test guidance (gov; annual 3PAO pen test).
[^49]: https://algoy.com/ — What US/UK/EU regulators examine (trade-press; on-site not questionnaire, RTO/RPO, 2-24hr notification, audit rights).
[^50]: https://www.fdic.gov/ — Interagency TPRM / Shared Assessments comment (gov; SOC/SIG in RFP, on-site for critical, continuous monitoring).
[^52]: https://aloan.ai/guides/ — SOC 2 Type II for commercial lending AI (trade-press/vendor; SOC 2 silent on model risk, contract terms, buy-decision).
[^54]: https://salesmotion.io/blog/how-to-sell-to-financial-services — How to sell to financial services (trade-press; exam cycles, 5-8 stakeholders, reference other banks).
[^55]: https://abmatic.ai/blog/ — ABM strategy for UK FS 2026 (trade-press/vendor; trust primary currency; one reference unlocks dozens).
[^56]: https://danishleadco.io/blog/proof-led-outreach-for-fintech-startups-selling-to-banks — Proof-led outreach selling to banks (trade-press/vendor; match asset-size/geography; partner-borrowed credibility).
[^57]: https://www.fintechtris.com/blog/ — Mastering B2B fintech enterprise sales (trade-press; "consultants not sellers"; shared risk).
[^58]: https://a16z.com/b2fi-demystifying-software-sales-into-financial-institutions/ — a16z "B2FI" (VC/analyst; core-provider channel FIS/Fiserv/Jack Henry; analyst credibility).
[^59]: https://defensive.cloud/ — RFP vendor-security checklist 2026 (trade-press; AES-256/TLS 1.3, BYOK/HSM FIPS 140-2/3, KMS logs).
[^60]: https://learn.daydream.ai/templates/ — Vendor encryption standards template (vendor; CMEK/BYOK, tenant isolation, crypto-shredding).
[^61]: https://www.kiteworks.com/regulatory-compliance/ — European banks: meet EBA guidelines (vendor; DORA encryption at rest/transit/in-use).
[^62]: https://www.kiteworks.com/regulatory-compliance/eba-encryption-key-control-guidelines/ — EBA encryption key control (vendor; key control without provider assistance).
[^63]: https://www.kiteworks.com/secure-file-sharing/customer-controlled-encryption-keys-banking/ — Customer-controlled keys in banking (vendor; key destruction = deletion mechanism).
[^64]: https://www.kiteworks.com/third-party-risk/ — Data sovereignty third-party sharing (vendor; BYOK ≠ customer-owned keys, CLOUD Act).
[^66]: https://www.cedaribsifintechlab.com/ — Modern core banking: build/buy/partner (trade-press; custom core >$50M).
[^67]: https://plumery.com/ — Buy vs build vs buy-and-build (vendor; hybrid "buy and build").
[^68]: https://vlinkinfo.com/blog/custom-vs-cots-banking-modernization — Custom vs COTS banking modernization 2026 (trade-press/vendor; composable core).
[^69]: https://www.ncino.com/ — The build vs buy dilemma (vendor, interested; hybrid; focus on differentiators).
[^70]: https://www.zartis.com/ — How digital banks scale: build vs buy (trade-press/vendor; build/buy matrix by component).
[^71]: https://team8.vc/report/build-vs-buy-for-banks-in-the-age-of-ai/ — Build vs buy for banks in the age of AI (VC/analyst; not building cores or foundation models; data/talent constraints).
[^72]: https://www.tealhq.com/ — JPMorgan Client Solutions Architect (primary JD; pre-sales technical adoption, client↔internal liaison).
[^73]: https://careers.cognizant.com/ — Cognizant Pre-Sales Solutions Architect BFS (primary JD; trusted-advisor, sales↔technical liaison).
[^75]: https://www.stacybishop.com/articles/how-to-sell-fintech-to-banks — How fintech founders sell to banks (practitioner; "no one owns the problem internally" kills deals).
