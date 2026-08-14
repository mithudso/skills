<!-- Provenance: reference under the `big-bank-IT` skill. Created 2026-06-18 via /dr deep research (regulator primary + analyst + vendor sources). Educational, vendor-neutral technical orientation — NOT legal/compliance advice. Disconfirming evidence on multi-cloud, concentration risk, sovereign cloud, and repatriation is preserved. -->

# The Data-Center → Cloud Journey at Large Banks, and the Cloud-Architecture Constraints

`verified-as-of: 2026-06-18` (cloud-adoption figures, sovereign-cloud offerings, and the regulatory landscape are fast-moving — re-verify dated specifics before quoting to a customer).

> **Educational orientation, NOT legal/compliance advice.** This is the IT/infrastructure realization of supervisory expectations; the **regulatory map** (Basel/Dodd-Frank/FFIEC/GLBA, why banks buy) → `fsi-banking-regulatory-context`; the **DORA/CTP/TPRM *mechanics*** (the regimes as governance processes) → `enterprise-vendor-management-and-tprm`; **MongoDB networking/multicloud/compliance config** → `mongodb-aws-networking` / `mongodb-atlas-multicloud` / `mongodb-compliance`. Hyperscaler/vendor playbooks here are vendor-tier (factual but interested).

## Contents

- [The estate today: hybrid is the steady state](#the-estate-today-hybrid-is-the-steady-state)
- [The supervisory guidance a bank's architects must satisfy](#the-supervisory-guidance-a-banks-architects-must-satisfy)
- [The four constraints that shape cloud architecture](#the-four-constraints-that-shape-cloud-architecture)
- [Sovereign cloud as the adoption consequence](#sovereign-cloud-as-the-adoption-consequence)
- [Repatriation: a selective counter-trend](#repatriation-a-selective-counter-trend)
- [Disconfirming evidence](#disconfirming-evidence)
- [Seller takeaways](#seller-takeaways)
- [Sources](#sources)

## The estate today: hybrid is the steady state

**Breadth is near-universal; depth is shallow — this is the single most important framing.** ~90%+ of FS firms "have some cloud" (Capgemini: 91% in 2023, up from 37% in 2020), but depth is far lower: **Accenture's 2025 index of 78 of the world's largest banks found only ~10% of core banking workloads have moved.**[^1][^7] A bank can truthfully say "we're a cloud company" while its systems of record still run on dedicated infrastructure. *(Confidence: fact.)*

**Hybrid is explicitly the end-state, not a way-station.** A 2026 analysis states the hybrid pattern (cloud for new/customer-facing/analytics/dev-test; cores on dedicated/private cloud) "is not a transitional state. It is the [steady] state for most large financial institutions"; Accenture calls hybrid "the preferred long-term model"; the **US Treasury's 2023 cloud report** found many larger institutions "plan to adopt a hybrid model."[^6][^7][^10] Gartner projects ~90% of organizations on hybrid/multi-cloud by 2027.[^12] *(Confidence: fact.)*

**The private-cloud layer is real and currently in upheaval.** Banks run substantial private cloud on **VMware** and **Red Hat OpenShift**. Broadcom's late-2023 VMware acquisition triggered a licensing "price shock" (perpetual licenses scrapped, mandatory subscriptions, per-core pricing; documented increases from ~2× to 8–15×, occasionally far higher), pushing banks toward OpenShift Virtualization (KubeVirt) — e.g. PNC, Emirates NBD, Morgan Stanley migrations.[^26][^27] *(Confidence: fact; specific bank examples are single-source each — qualified.)*

## The supervisory guidance a bank's architects must satisfy

A unifying principle runs through all blocs: **outsourcing to a cloud provider never transfers the bank's accountability** — the bank remains responsible "as if the services were performed by the bank itself."[^17][^30] *(Confidence: fact.)*

**US guidance lineage (note the rescissions — cite the current baseline):**
- **OCC Bulletin 2013-29** (third-party risk, 2013) and its cloud-specific FAQ supplement **OCC Bulletin 2020-10** (2020) were the foundational doctrine — **both were rescinded and replaced** on **June 6 2023** by the **Interagency Guidance on Third-Party Relationships: Risk Management** (OCC Bulletin 2023-17), issued jointly by the **OCC, Federal Reserve, and FDIC**. Cite the **2023 interagency guidance** as the current US baseline.[^3][^4][^30] *(Confidence: fact.)*
- **FFIEC "Security in a Cloud Computing Environment"** — a joint statement (Apr 30 2020); core message: cloud is a **shared-responsibility** model and "management should not assume that effective security and resilience controls exist simply because [systems] are operating in a cloud computing environment"; it "does not contain new regulatory expectations." The **FFIEC IT Examination Handbook** (esp. the Outsourcing booklet) remains the operational reference.[^2][^9] *(Confidence: fact.)*
- **FedRAMP relevance is indirect for commercial banks** — it governs cloud services handling *federal* data for *federal agencies* and "does not regulate private companies." Its relevance is as a *de facto security baseline* banks borrow from, and binding only where a bank touches government programs. Baker McKenzie confirms there are "no laws or regulations that expressly … govern banks' use of cloud services" in the US — banks operate under general safety-and-soundness + FFIEC/interagency guidance.[^24][^30] *(Confidence: fact.)*

**EU/UK guidance (architectural demands, not deep mechanics — `enterprise-vendor-management-and-tprm` owns the regulatory mechanics):**
- **EBA Guidelines on outsourcing (EBA/GL/2019/02)**, applicable from Sep 30 2019 — the European baseline; mandate documented **exit strategies**, require monitoring **concentration risk** (explicitly flagging cloud outsourcing as "dominated by a small number of … providers" that can become "a single point of failure"), and forbid a firm becoming an "empty shell."[^4][^5]
- **DORA (Regulation (EU) 2022/2554)**, applicable **Jan 17 2025** — demands a board-owned ICT third-party-risk strategy, registers of information, and an EU-level **oversight framework for Critical ICT Third-Party Providers (CTPPs)** with a formal ICT concentration-risk definition.[^13][^18][^33]
- **UK PRA/FCA: SS2/21** ("Outsourcing and third party risk management," PRA; effective Mar 31 2022, updated Dec 31 2024) implements the EBA GL and expands on exit plans — and separately the statutory **Critical Third Parties (CTP) regime** (final rules **PS16/24** + **SS6/24**, both Nov 12 2024) took effect **Jan 1 2025** under FSMA 2023, giving the Bank of England/PRA/FCA direct oversight of designated systemic providers. Both stress that designation **does not reduce the firm's accountability.**[^23][^25] *(Confidence: fact.)*

## The four constraints that shape cloud architecture

These are the questions a bank's architects will actually raise — not "is it secure," but:

**(1) Concentration risk — "hidden" and increasingly geopolitical.** Three hyperscalers (AWS, Azure, Google) hold ~two-thirds of global cloud infrastructure.[^13] The systemic concern is documented across regulators: the US Treasury 2023 report named market concentration a top-five challenge ("if an incident occurs at one CSP, it could affect many financial-sector clients concurrently"); a Bank of England 2020 survey found "most banks and insurers rely on just **two providers**."[^21][^22][^31][^19] The FSB highlights **"fourth-party" risk** (opacity from providers-of-providers). *(Confidence: fact.)*

**(2) Exit / stressed-exit strategy — documented AND tested.** This is the load-bearing demand. EBA GL para 107 requires a "sufficiently tested exit plan" for critical functions, including unplanned/stressed exit; DORA reinforces continuity through severe disruption of a critical provider.[^5][^33] The architect's task is to "prove reversibility through annual drills and documented extraction paths." *(Confidence: fact.)*

**(3) Data portability / interoperability — expected, but in tension with managed-service depth.** EBA GL requires arrangements that "facilitate the transfer of the outsourced function to another service provider"; the EU **Data Act** (egress-fee restrictions from Jan 12 2027) further lowers switching friction.[^4][^12] The engineering tension: true portability requires avoiding cloud-unique managed services, forfeiting much of cloud's value — leading some banks to *document* portability they don't exercise ("portability theatre").[^10] *(Confidence: fact, with the tension qualified.)*

**(4) Data residency / sovereignty — the fastest-moving constraint.** Beyond physical location, the binding concern is *jurisdiction*: the **US CLOUD Act (2018)** compels US-headquartered providers to produce data in their "possession, custody, or control" regardless of where it sits. This is why **encryption-key custody** matters architecturally:[^35][^41]
  - **BYOK** (bring your own key) — bank generates the key but uploads it to the provider's KMS; the provider retains technical access and **can be compelled to decrypt**.[^41][^43]
  - **HYOK / external key store** (hold your own key) — keys never leave the bank's HSM; the provider processes only ciphertext and "cannot produce readable data regardless of legal demands." AWS markets External Key Stores explicitly as HYOK for regulated workloads; it carries higher operational burden and limits usable services.[^42][^41] *(Confidence: fact.)*

## Sovereign cloud as the adoption consequence

The hyperscalers responded to EU sovereignty pressure with sovereign offerings (typically a **smaller service catalog at a ~15% premium**):[^38][^40]
- **AWS European Sovereign Cloud** — GA **Jan 15 2026**, first region in Brandenburg, Germany; physically/logically separate EU-resident operation, independent IAM/billing, ~90 services at launch; >€7.8B investment.[^38][^39]
- **Microsoft Sovereign Cloud** — Sovereign Public Cloud announced Jun 16 2025 across existing EU regions, plus partnership models like **Bleu** (Orange/Capgemini).[^40]
- **Oracle EU Sovereign Cloud** (2023, earliest); **Google** S3NS (Thales partnership).[^40] *(Confidence: fact on existence/structure; vendor-tier.)*

## Repatriation: a selective counter-trend

Headlines tout "86% plan repatriation," but the honest reading: these measure *any* workload movement. **IDC data shows only ~8–9% plan *full* repatriation**, while public-cloud spend still grew ~21.5% in 2025 (Gartner) — mathematically incompatible with a wholesale exit.[^12][^20] The accurate framing is **"right workload, right place"** / selective placement, driven by AI/GPU economics, data-gravity/egress costs, the VMware shock, and (for banks) DORA's demand for a credible on-prem fallback. *(Confidence: fact — the ~8–9% full-repatriation figure recurs across IDC-citing sources.)*

## Disconfirming evidence

1. **Multi-cloud is more aspiration than reality, and may be counterproductive.** Only ~31% of banks actually run multi-cloud (Accenture).[^7] The US Treasury report records the *sector's own* feedback that multi-cloud was "too technically complex and [the] operational risk too high"; leading analysts "recommend against such an approach for increasing operational resiliency."[^19][^6] Splitting compute and data across clouds can multiply data-layer cost 30×+. **Capital One went all-in on a single cloud (AWS)**, exiting its data centers — a regulated bank choosing deep integration over multi-cloud. So "go deep on one, with a tested exit" is a defensible competing posture; don't assume multi-cloud is universal best practice.
2. **Concentration risk is contested, not settled.** The FSB itself notes concentration "does not automatically pose systemic risks … concentration can reflect the quality, including the resilience, of a third party's services."[^19] The pro-hyperscaler argument (the big three are *more* reliable than most banks' own DCs) competes with the systemic-tail-risk argument (rarer but catastrophic correlated failures; hollowing out banks' in-house expertise). Honest synthesis: concentration **lowers idiosyncratic risk while raising systemic tail risk.**
3. **Sovereign cloud may be "sovereignty washing."** US-hyperscaler sovereign clouds achieve *operational* sovereignty (EU staff/infra/ops) but **cannot deliver *jurisdictional* sovereignty** while the parent is a US company subject to the CLOUD Act. The smoking gun: in **June 2025 Microsoft France's Chief Legal Officer testified under oath to the French Senate that he could *not* guarantee French data would never be passed to US authorities.**[^35][^36] For a bank architect, the practical takeaway: **only HYOK/external-key-store (and ultimately a non-US-jurisdiction provider) closes the gap** — BYOK and "data stays in the EU" do not.
4. **"Cloud-first" adoption stats are inflated by definitional looseness and lift-and-shift.** Capgemini cautions that high adoption "does not imply … full-scale or even effective migration" — most firms do lift-and-shift, forfeiting cloud's benefits. The breadth number is real; the depth it implies is not.[^1]

## Seller takeaways

- **Pitch a workload's placement, not a cloud-transformation thesis** — that thesis is empirically weak, and most workloads are still on-prem.
- **A run-anywhere product with a documentable exit/egress path is a compliance asset** — it directly de-risks the buyer's concentration/exit/DORA/CTP exposure. "Runs the same on-prem, in the bank's own cloud account, and in a sovereign region" is a first-class FSI buying criterion.
- **Know the key-custody story cold.** For sovereignty-sensitive deals, BYOK is not enough to answer the CLOUD Act question; be ready to discuss HYOK/external key stores and what they cost operationally.

## Sources

[^1]: Capgemini, World Cloud Report – Financial Services — https://www.capgemini.com/wp-content/uploads/2025/11/WCR_2026_Final-2MB-version.pdf — analyst — breadth 91%; lift-and-shift caveat; hybrid steady state.
[^2]: FFIEC Joint Statement "Security in a Cloud Computing Environment" (FDIC FIL-52-2020) — https://www.fdic.gov/news/financial-institution-letters/2020/fil20052a.pdf — primary regulator — shared responsibility; "no new regulatory expectations."
[^3]: OCC Bulletin 2020-10 (FAQs supplementing 2013-29, now rescinded) — https://www.occ.gov/static/rescinded-bulletins/bulletin-2020-10.pdf — primary regulator — cloud = third-party relationship.
[^4]: OCC Bulletin 2023-17, Interagency Guidance on Third-Party Relationships — https://www.occ.gov/news-issuances/bulletins/2023/bulletin-2023-17.html — primary regulator — rescinds 2013-29 & 2020-10; current US baseline.
[^5]: EBA Guidelines on outsourcing (EBA/GL/2019/02) — https://www.bde.es/f/webbde/INF/MenuHorizontal/Normativa/guias/EBA-GL-2019_02_EN.pdf — primary regulator — exit strategies, concentration risk, single-point-of-failure language; "sufficiently tested exit plan" (para 107).
[^6]: US Treasury, The Financial Services Sector's Adoption of Cloud Services — https://home.treasury.gov/system/files/136/Treasury-Cloud-Report.pdf — primary (govt) — hybrid plans; concentration top-5 challenge; multi-cloud "too complex."
[^7]: Accenture Banking Blog (2025 Rotation Index, 78 banks) — https://bankingblog.accenture.com/cloud-ai-banking-growth — analyst — only ~10% of core workloads in cloud; only 31% multi-cloud; hybrid preferred.
[^9]: FFIEC press release + OCC Bulletin 2020-46 (cloud joint statement) — http://www.ffiec.gov/news/press-releases/2020/pr-04-30 — primary regulator — confirms joint-statement scope/date.
[^10]: Cloud Computing in U.S. Finance has become operational (TechBullion, 2026) — https://techbullion.com/cloud-computing-in-u-s-finance-has-stopped-being-strategic-and-become-operational/ — trade — hybrid not transitional; "portability theatre."
[^12]: Cloud Repatriation 2026 Is a Statistical Illusion (Digital-Chiefs) — https://www.digital-chiefs.de/en/cloud-repatriation-2026-statistical-illusion/ — trade — 86% (any) vs 8–9% (full); 21.5% cloud growth; EU Data Act egress 2027.
[^13]: BaFin, When concentrations become a risk (DORA) — https://www.bafin.de/SharedDocs/Veroeffentlichungen/EN/Fachartikel/2025/fa_250107_DORA_Auslagerungen_und_Konzentrationsrisiken_en.html — primary regulator — DORA applicable 17 Jan 2025; CTPP oversight; ~two-thirds market share of 3 hyperscalers.
[^17]: Baker McKenzie, US Cloud Compliance Center — https://resourcehub.bakermckenzie.com/en/resources/cloud-compliance-center/na/united-states — legal — no US law expressly governs bank cloud use; accountability not transferred.
[^18]: EBA Guide on DORA oversight activities (JC 2025 29) — https://www.eba.europa.eu/sites/default/files/2025-07/0eb4ef25-1bce-4f17-b3d1-c46d59f59428/JC%202025%2029%20-%20DORA%20Oversight%20Guide.pdf — primary regulator — CTPP designation; ICT concentration-risk definition.
[^19]: FSB, Third-party dependencies in cloud services — https://www.fsb.org/uploads/P091219-2.pdf — primary (intl) — "fourth parties"; BoE "two providers"; multi-cloud judged too complex; concentration not automatically systemic.
[^20]: AlixPartners / IDC-citing repatriation analyses — https://www.alixpartners.com/media/gjql1l4g/alixpartners-cloud-repatriation-2026.pdf — analyst/trade — ~8–9% full repatriation; reversibility drills; DORA → on-prem fallback.
[^21]: Finextra / American Banker on Treasury report — https://www.americanbanker.com/news/many-risks-lurk-for-banks-in-the-cloud-treasury-report — trade — AWS/Google/Azure named; outsized bargaining power.
[^22]: PIFS/Harvard, Cloud Adoption and Concentration Risk — https://www.pifsinternational.org/wp-content/uploads/2023/04/PIFS-Cloud-Adoption-and-Concentration-Risk-in-the-Financial-Sector.pdf — analyst — concentration-risk framing.
[^23]: Bank of England PS16/24 (CTP final rules) — https://www.bankofengland.co.uk/prudential-regulation/publication/2024/november/operational-resilience-critical-third-parties-uk-financial-sector-policy-statement — primary regulator — CTP rules effective 1 Jan 2025; firm accountability preserved.
[^24]: FedRAMP scope clarification (overview) — https://www.fedramp.gov/program-basics/ — gov — FedRAMP governs federal-agency cloud use, not private companies.
[^25]: Bank of England SS2/21 — https://www.bankofengland.co.uk/prudential-regulation/publication/2021/march/outsourcing-and-third-party-risk-management-ss — primary regulator — SS2/21 effective Mar 31 2022 (updated Dec 31 2024); exit plans (Ch. 10).
[^26]: PNC's modernization with OpenShift Virtualization (Red Hat) — https://www.redhat.com/en/blog/pncs-infrastructure-modernization-journey-red-hat-openshift-virtualization — vendor — PNC/Emirates NBD/Morgan Stanley VMware-exit examples.
[^27]: VMware Exit Playbook (EU Cloud Patterns) / TechTarget — https://www.eucloudpatterns.eu/posts/vmware-exit-playbook/ — analyst/trade — Broadcom price increases 2×–15×; per-core pricing.
[^30]: OCC Bulletin 2023-17 detail / Bank Service Company Act logic — https://www.occ.gov/news-issuances/bulletins/2023/bulletin-2023-17.html — primary regulator — accountability remains with the bank.
[^31]: Bank of England, How reliant are banks and insurers on cloud outsourcing? (2020) — https://www.bankofengland.co.uk/bank-overground/2020/how-reliant-are-banks-and-insurers-on-cloud-outsourcing — primary regulator — IaaS "already highly concentrated."
[^33]: DORA RTS + EIOPA DORA page — https://www.eiopa.europa.eu/digital-operational-resilience-act-dora_en — primary/legal — DORA in force 17 Jan 2025; ICT third-party strategy.
[^35]: The sovereign cloud that isn't / "Trust me bro" (Heise / Computer Weekly) — https://www.computerweekly.com/feature/Is-cloud-data-sovereignty-all-just-a-case-of-Trust-me-bro — trade — CLOUD Act extraterritoriality; "sovereignty washing"; Microsoft French Senate admission.
[^36]: Sovereignty-washing analyses (CSO Online / SecureCloud) — https://www.csoonline.com/article/4184634/ — trade — Microsoft France CLO sworn testimony (June 2025); DORA-driven sovereign migration.
[^38]: Opening the AWS European Sovereign Cloud (AWS) — https://aws.amazon.com/blogs/aws/opening-the-aws-european-sovereign-cloud/ — vendor (primary) — GA 15 Jan 2026; Brandenburg; data/metadata residency; >€7.8B.
[^39]: AWS ESC coverage (The Register) — https://www.theregister.com/2026/01/15/aws_european_sovereign_cloud/ — trade — ~90 services at launch; EU-resident operation.
[^40]: Microsoft Sovereign Cloud / Oracle / Google sovereign offerings — https://blogs.microsoft.com/blog/2025/06/16/announcing-comprehensive-sovereign-solutions-to-help-european-organizations/ — vendor/trade — ~15% premiums; Microsoft Sovereign Public Cloud (Jun 16 2025); Bleu/S3NS.
[^41]: Differences between BYOK and HYOK (Fortanix) — https://www.fortanix.com/blog/differences-between-byok-and-hyok — vendor — BYOK (key in provider KMS) vs HYOK (key never leaves customer); HYOK for finance/gov.
[^42]: AWS KMS External Key Stores docs — https://docs.aws.amazon.com/kms/latest/developerguide/keystore-external.html — vendor (primary) — External Key Store = HYOK; "designed for regulated workloads."
[^43]: EBA/CLOUD-Act key control (Kiteworks) — https://www.kiteworks.com/regulatory-compliance/european-banks-eba-outsourcing-encryption/ — vendor — BYOK can be compelled under CLOUD Act; HYOK = provider holds only ciphertext.
