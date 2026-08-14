<!-- Provenance: reference under the `enterprise-vendor-management-and-tprm` skill. Created 2026-06-18 via /dr deep-research. Buyer-side governance lens for a B2B/SaaS seller (MongoDB TAM) selling into banks. Educational context, not legal/compliance advice. FAST-MOVING DOMAIN — re-verify all dates/designation statuses before quoting. -->

# Operational Resilience & Systemic Third-Party Regulation (Frontier)

The **regulatory pressure** that forces a bank to govern its technology/cloud vendors a certain way — DORA,
US OCC/FFIEC, the UK PRA/FCA Critical Third Parties regime, Basel — plus the cross-cutting **concentration /
fourth-party / exit-plan** frontier. This is *why* a bank's audit-rights, sub-processor-disclosure, and
exit-plan demands exist; a seller who knows the source argues less and prepares more.

> **FAST-MOVING — verify before quoting.** This is the most date-sensitive material in the skill. Every
> effective date, designation, and rule status is stamped **as of 2026** and several are explicitly flagged as
> unconfirmable. **Re-verify against the primary regulator source before telling a customer anything dated
> here.** Educational context, NOT legal/compliance advice. Cited inline as `[^n]`.

## Contents
- EU DORA (Digital Operational Resilience Act)
- US — OCC/FFIEC / interagency operational-resilience & concentration
- UK — PRA/FCA operational resilience + Critical Third Parties (CTP)
- Basel — Principles for Operational Resilience
- Concentration, fourth-party, exit plans & substitutability
- Disconfirming nuance
- Dates/statuses NOT independently confirmed
- References

---

## EU DORA (Digital Operational Resilience Act)

- **DORA became applicable on 17 January 2025** (Regulation (EU) 2022/2554). It applies to financial entities
  *and*, via an oversight framework, to ICT third-party service providers.[^dora1]
- **Five pillars:** (1) ICT risk management; (2) ICT-related incident management, classification & reporting;
  (3) digital operational resilience testing (incl. **TLPT** / threat-led penetration testing); (4) ICT
  third-party risk management; (5) information/intelligence sharing.[^dora1]
- **CTPP Oversight Framework & Lead Overseer.** The three ESAs (EBA, EIOPA, ESMA) designate ICT providers as
  **"critical" (CTPPs)**; one ESA acts as **Lead Overseer** (tied to the financial entities holding the
  largest asset share using that provider), with power to request information, investigate, and conduct on-site
  inspections (Art. 35) and issue recommendations.[^dora2]
- **First CTPP designations published 18 November 2025 — 19 providers, including AWS, Google Cloud, and
  Microsoft (Azure)** plus data-centre operators, telecoms, and fintech specialists; the list updates
  annually.[^dora3] **This is the single most important recent development for a cloud-vendor TAM** — the
  hyperscalers MongoDB Atlas runs on are now directly overseen.
- **Register of Information.** Each financial entity must maintain a register of **all** ICT third-party
  contractual arrangements (Art. 28(3)); the ESAs used these registers to drive the criticality assessment.
  This is the mechanism that forces a bank to enumerate every sub-processor — directly behind sub-processor-
  disclosure demands on MongoDB.[^dora4]
- **Contractual requirements (Art. 30).** Mandatory clauses in ICT contracts, stricter for contracts
  supporting "critical or important functions": full **audit / access / inspection rights** for the entity,
  its appointees, *and competent authorities* (incl. on-site); **exit strategies** with a mandated adequate
  transition period; **subcontracting** conditions and ICT-supply-chain monitoring; data protection/security,
  SLAs, incident-reporting/assistance obligations, termination rights. **This is the contractual engine behind
  audit-right, exit-plan, and sub-processor demands on a cloud vendor.**[^dora5]
- **TLPT (Arts. 26–27).** Threat-Led Penetration Testing required **at least every 3 years** for financial
  entities identified as significant by their competent authority (not all entities). The Eurosystem updated
  **TIBER-EU** to align with DORA's TLPT RTS (ECB, Feb 2025). TLPT is intelligence-led red-teaming, distinct
  from a standard pen test.[^dora6]

## US — OCC/FFIEC / interagency operational-resilience & concentration

- **June 2023 Interagency Guidance on Third-Party Relationships: Risk Management** (Fed + FDIC + OCC; Federal
  Register 9 June 2023; OCC Bulletin 2023-17; FDIC FIL-29-2023; Fed SR 23-4) **replaced** each agency's prior
  guidance (OCC 2013-29 + 2020-10 FAQs, Fed SR 13-19, FDIC 2008), creating one consistent supervisory approach.
  *(TPRM-lifecycle mechanics → `references/tprm-lifecycle-and-risk-tiering.md`; here it's the regulatory
  framing.)*[^us1]
- **2020 "Sound Practices to Strengthen Operational Resilience"** — interagency paper (OCC Bulletin 2020-94;
  Fed SR 20-24; FDIC). Targets large banks (≥$250B total assets, or ≥$100B with other risk characteristics);
  draws on operational-risk, BCM, third-party-risk, cyber, and recovery/resolution standards. The current US
  "operational resilience" baseline — **guidance, not a binding rule.**[^us2]
- **FFIEC "Architecture, Infrastructure, and Operations" (AIO) booklet** issued **30 June 2021** (OCC Bulletin
  2021-30; Fed SR 21-11), rescinding the 2004 "Operations" booklet and extending risk-management standards
  across outsourced activities. The FFIEC IT Handbook also retains dedicated **Outsourcing Technology
  Services** and **Business Continuity Management** booklets.[^us3]
- **A dedicated US bank "operational resilience" RULE remains exploratory / NOT proposed (as of mid-2026).**
  Then–Acting Comptroller Michael Hsu (speech 12 March 2024) said the OCC was "exploring baseline operational
  resilience requirements for large banks with critical operations, including third-party service providers" —
  but no ANPR/NPR has been confirmed issued, and 2025–26 prudential messaging tilted toward "rightsizing
  regulation." **Treat any specific US operational-resilience rule as not-yet-existent.**[^us4]
- **Cloud concentration as a US systemic concern.** The US Treasury's Feb 2023 report on cloud adoption in
  financial services flagged concentration, limited CSP bargaining leverage for smaller banks, and
  incident-response coordination gaps; Treasury + FBIIC/FSSCC set up workstreams (e.g., a Cloud Executive
  Steering Group) rather than hard rules. FSOC annual reports and OFR have repeatedly named third-party/cloud
  concentration as a vulnerability.[^us5]

## UK — PRA/FCA operational resilience + Critical Third Parties (CTP)

- **Operational-resilience policy.** FCA **PS21/3** "Building operational resilience" + PRA **SS1/21**. Firms
  must identify **Important Business Services (IBS)**, set **impact tolerances**, and **map and test** their
  ability to stay within tolerance.[^uk1]
- **Hard deadline: 31 March 2025** — by which in-scope firms had to have completed mapping/testing and made
  investments so they can **remain within impact tolerances** for each IBS (after a transitional period that
  began 31 March 2022). The regime is now in its ongoing/business-as-usual supervisory phase.[^uk2]
- **Critical Third Parties (CTP) regime.** Statutory powers came via **FSMA 2023** (new s.312L FSMA 2000)
  letting **HM Treasury designate** a provider as a CTP if disruption "could threaten the stability of, or
  confidence in, the UK financial system." Final rules published **12 November 2024** (PRA **PS16/24** / FCA
  **PS24/16**; PRA **SS6/24**), **in force 1 January 2025**. Regulators (PRA/FCA/BoE) gain oversight powers
  over designated CTPs, including the ability to prohibit a firm's use of a deficient CTP. CTP powers extend
  *only* to services provided to the financial sector, not the provider's entire business.[^uk3]
- **No UK CTPs designated yet (status caveat).** As of the latest confirmable sources, **HM Treasury had not
  yet designated any CTPs** — the rules apply to a provider only once its HMT designation order takes effect
  (with its own transitional periods). **Re-verify current designation status before citing.**[^uk4]
- **EU–UK cooperation.** The ESAs and UK regulators (BoE/PRA/FCA) signed an **MoU on cross-border oversight of
  critical ICT providers on 14 January 2026**, signaling coordinated oversight of the same global
  hyperscalers.[^uk5]

## Basel — Principles for Operational Resilience

- **BCBS "Principles for Operational Resilience" published 31 March 2021** (BIS document d516), promoting a
  principles-based approach so banks can withstand, adapt to, and recover from severe disruptions.[^basel1]
- **Seven principles:** (1) governance; (2) operational risk management; (3) business continuity planning &
  testing; (4) mapping of interconnections & interdependencies of critical operations; (5) **third-party
  dependency management**; (6) incident management; (7) resilient ICT incl. cybersecurity.[^basel1]
- **Relationship to PSMOR.** Published concurrently with a revised **"Principles for the Sound Management of
  Operational Risk" (PSMOR)**; the Committee frames operational resilience as the *outcome* of effective
  operational-risk management — the two are complementary.[^basel2]
- **Significance.** Basel is the soft-law backbone; DORA, the UK regime, and US sound practices are regional
  implementations of the same third-party-dependency / ICT-resilience ideas — useful for a TAM to frame "this
  is global, not just EU."[^basel2]

## Concentration, fourth-party, exit plans & substitutability

- **Cloud concentration is an explicit systemic concern.** AWS + Azure + Google Cloud hold ~63% of the global
  cloud market; one analysis cited >65% of EU financial entities using at least two of the three for critical
  functions. The 18 Nov 2025 DORA CTPP list is the first concrete regulatory act on this.[^conc1]
- **CrowdStrike outage (19 July 2024).** A faulty update bricked ~8.5M Windows machines (~70% of Fortune 500
  affected). Widely cited as the canonical real-world single-point-of-failure event — and notably *not* a
  cloud-hyperscaler failure but a single-vendor software-update failure, which **broadened the concentration
  conversation beyond IaaS** to any widely-deployed software/agent. The FCA published observations/lessons on
  31 October 2024, tying it to operational-resilience and CTP priorities.[^conc2]
- **Exit strategies / stressed exit / substitutability** (the contractual demand chain). The **EBA Guidelines
  on Outsourcing (Feb 2019)** already require a documented, *tested* exit strategy for critical/important
  functions covering **stressed exit** (provider failure/degradation), an explicit **substitutability
  assessment** (easy / difficult / impossible), reintegration feasibility, and consideration of **step-in
  rights**; DORA Art. 30 hardens the contractual transition-period requirement.[^conc3] **This is exactly what
  becomes a contractual demand on a vendor: transition assistance, data portability/export, and a workable
  exit.** (How this lands as a deal friction — incl. source-code/SaaS escrow — is in the SKILL.md
  Selling-in playbook.)
- **Fourth-party / sub-contracting chain (the *regulatory* hook).** Fourth-party and concentration risk are
  defined once in `references/tprm-lifecycle-and-risk-tiering.md` (their definitional home) — this bullet adds
  only what the *regulation* does with them: DORA's Register of Information and Art. 30 sub-contracting
  conditions force banks to look *past* their direct vendor into the full ICT supply chain, so a vendor is
  asked to disclose *its* critical sub-processors as a contractual obligation.[^conc4]

## Disconfirming nuance

- **Regulators do NOT automatically treat concentration as systemic.** The UK PRA/FCA/BoE (aligned with the
  FSB) stated their CP "does not consider that concentration in the provision of third-party services…
  automatically pose[s] systemic risks." So the demand is *resilience/oversight* of the provider, not a hard
  cap on using a popular vendor.[^nuance-conc]
- **Regulators explicitly do NOT expect "perfect" exit plans.** Given hybrid/multi-cloud complexity,
  regulators increasingly acknowledge perfect exits aren't always technically/economically feasible — a useful
  counter to a customer who frames an exit-plan ask as "we must be able to fully leave you tomorrow."
  *(Single strong source — `TENTATIVE`.)*[^nuance-exit]
- **"Regulation alone is unlikely to materially reduce concentration."** Skeptics argue cloud economics keep
  driving banks to the largest platforms, that split agency mandates make CSPs hard to regulate, and that
  hyperscalers may be *more* resilient than in-house alternatives — so concentration may be a rational
  trade-off.[^nuance-skeptic]
- **DORA is criticized as prescriptive and burdensome** (per-entity/per-country incident reporting; overlap
  with NIS2 / PCI DSS 4.0; fines up to 2% of total annual worldwide turnover for entities, and CTPPs facing
  periodic penalty payments up to 1% of average daily worldwide turnover). Relevant because banks pass this
  burden onto vendors as documentation/evidence demands.[^nuance-burden]

## Dates/statuses NOT independently confirmed (handle with care)

1. **A US "operational-resilience frameworks in place by July 1 2025" deadline** — a single weak vendor source
   claimed this; it conflicts with the "no binding US rule exists" picture. **Do not cite a US op-res
   deadline.**[^us4]
2. **Whether HM Treasury has designated ANY UK CTPs yet** — best evidence says *none* as of finalized-rules
   reporting, but no positive mid-2026 confirmation of zero designations was found. **Re-verify.**[^uk4]
3. **The full 19-provider DORA CTPP roster** beyond AWS / Google Cloud / Microsoft — the "19 total,
   hyperscalers included" count is well-corroborated; the complete named list sits in the ESMA designated-CTPP
   PDF and is not individually verified here.[^dora3]

*(Multiply-confirmed and treated as FACT: DORA application 17 Jan 2025; UK op-res deadline 31 Mar 2025; UK CTP
rules in force 1 Jan 2025; Basel principles 31 Mar 2021; June 2023 US interagency guidance; 18 Nov 2025 CTPP
designation; 14 Jan 2026 EU–UK MoU.)*

---

## References

[^dora1]: DORA applicable 17 Jan 2025; scope; five pillars. EIOPA; ESMA; Faegre Drinker; InnReg — regulator-primary/law-firm. https://www.eiopa.europa.eu/digital-operational-resilience-act-dora_en ; https://www.esma.europa.eu/esmas-activities/digital-finance-and-innovation/digital-operational-resilience-act-dora — verified-as-of 2026-06-18
[^dora2]: CTPP oversight framework; ESAs designate; Lead Overseer powers (Art. 35). EIOPA; FMA Austria; Noerr — regulator-primary/law-firm. https://www.eiopa.europa.eu/digital-operational-resilience-act-dora/dora-oversight_en ; https://www.noerr.com/en/insights/an-overview-of-the-dora-supervisory-framework — verified-as-of 2026-06-18
[^dora3]: First CTPP designations 18 Nov 2025 — 19 providers incl. AWS/Google Cloud/Microsoft. EBA; EIOPA; Morgan Lewis; AWS self-confirmation — regulator-primary/law-firm/vendor. https://www.eba.europa.eu/publications-and-media/press-releases/european-supervisory-authorities-designate-critical-ict-third-party-providers-under-digital ; https://aws.amazon.com/blogs/security/aws-designated-as-a-critical-third-party-provider-under-eus-dora-regulation — verified-as-of 2026-06-18
[^dora4]: Register of Information (Art. 28(3)) drove criticality assessment. EBA — regulator-primary. https://www.eba.europa.eu/publications-and-media/press-releases/esas-publish-report-landscape-ict-third-party-providers-eu
[^dora5]: DORA Art. 30 mandatory contract clauses (audit/access/regulator rights, exit + transition period, subcontracting). Securiti; Bird & Bird; DORA Art. 30 portal — practitioner/law-firm/official-text. https://www.twobirds.com/en/insights/2024/global/dora-some-insights-on-contractual-clauses-in-agreements-between-financial-entities ; https://www.digital-operational-resilience-act.com/Article_30.html — verified-as-of 2026-06-18
[^dora6]: TLPT (Arts. 26–27) ≥ every 3 years for significant entities; TIBER-EU updated for DORA (ECB Feb 2025). ECB; FinancialRegulations.eu; Yogosha — regulator-primary/practitioner. https://www.ecb.europa.eu/press/intro/news/html/ecb.mipnews250211.en.html — verified-as-of 2026-06-18
[^us1]: June 2023 interagency third-party guidance; replaced OCC 2013-29 / Fed SR 13-19 / FDIC 2008. Federal Register; OCC; FDIC — regulator-primary. https://www.federalregister.gov/documents/2023/06/09/2023-12340/interagency-guidance-on-third-party-relationships-risk-management ; https://www.occ.gov/news-issuances/bulletins/2023/bulletin-2023-17.html
[^us2]: 2020 "Sound Practices to Strengthen Operational Resilience" (OCC 2020-94; Fed SR 20-24); large-bank scope; guidance not rule. OCC; Federal Reserve — regulator-primary. https://www.occ.gov/news-issuances/bulletins/2020/bulletin-2020-94.html
[^us3]: FFIEC AIO booklet 30 June 2021 (OCC 2021-30; Fed SR 21-11); Outsourcing + BCM booklets. OCC; FFIEC IT Handbook — regulator-primary. https://www.occ.gov/news-issuances/bulletins/2021/bulletin-2021-30.html ; https://ithandbook.ffiec.gov/it-booklets/architecture-infrastructure-and-operations
[^us4]: US op-res rule exploratory/not-proposed; Hsu 12 Mar 2024 speech; 2025-26 "rightsizing"; July-2025 deadline UNCONFIRMED. OCC; ABA Banking Journal; FDIC — regulator-primary/press (status QUALIFIED/uncertain). https://www.occ.gov/news-issuances/news-releases/2024/nr-occ-2024-23.html ; https://bankingjournal.aba.com/2024/03/occs-hsu-banking-agencies-considering-changes-to-operational-resilience-requirements/
[^us5]: US cloud-concentration systemic concern; Treasury Feb 2023 report; FSOC/OFR. FSOC 2024 Annual Report; OFR 2025 Annual Report; AWS summary — regulator-primary/vendor. https://home.treasury.gov/system/files/261/FSOC2024AnnualReport.pdf ; https://www.financialresearch.gov/annual-reports/files/OFR-AR-2025.pdf
[^uk1]: UK op-res — FCA PS21/3 + PRA SS1/21; IBS, impact tolerances, map & test. FCA; Bank of England/PRA — regulator-primary. https://www.fca.org.uk/publications/policy-statements/ps21-3-building-operational-resilience ; https://www.bankofengland.co.uk/-/media/boe/files/prudential-regulation/supervisory-statement/2021/ss121-march-22.pdf
[^uk2]: UK op-res hard deadline 31 Mar 2025 (remain within impact tolerances). Sidley; Complyport — law-firm/practitioner. https://www.sidley.com/en/insights/newsupdates/2025/01/uk-operational-resilience-rules-are-you-ready-for-31-march-2025
[^uk3]: UK CTP regime — FSMA 2023 / s.312L; HMT designation; final rules 12 Nov 2024 (PRA PS16/24 / FCA PS24/16 / SS6/24); in force 1 Jan 2025; FS-only scope. FCA; Bank of England; Morgan Lewis; Kemp IT Law — regulator-primary/law-firm. https://www.fca.org.uk/publications/policy-statements/ps24-16-operational-resilience-critical-third-parties-uk-financial-sector ; https://www.bankofengland.co.uk/prudential-regulation/publication/2024/november/operational-resilience-critical-third-parties-to-the-uk-financial-sector-policy-statement — verified-as-of 2026-06-18
[^uk4]: No UK CTPs designated yet (status caveat — re-verify). Lexology; Ashurst — law-firm (QUALIFIED). https://www.lexology.com/library/detail.aspx?g=99f48f00-94c2-45ba-b94a-2d7e0fca51cc ; https://www.ashurst.com/en/insights/uk-ctp-regime-final-rules-from-regulators/ — verified-as-of 2026-06-18
[^uk5]: EU–UK MoU on cross-border oversight of critical ICT providers, 14 Jan 2026. FCA; EIOPA; Covington — regulator-primary/law-firm. https://www.fca.org.uk/news/statements/uk-and-eu-regulators-sign-memorandum-understanding-strengthen-oversight-critical-third-parties — verified-as-of 2026-06-18
[^basel1]: BCBS Principles for Operational Resilience (d516, 31 Mar 2021); seven principles. BIS; RiskBusiness — regulator-primary/practitioner. https://www.bis.org/bcbs/publ/d516.htm
[^basel2]: Op-res published concurrently with revised PSMOR; resilience as outcome of op-risk mgmt. McCarthy Tétrault; Monocle — law-firm/consultancy. https://www.mccarthy.ca/en/insights/publications/basel-committee-issues-principles-operational-resilience-and-risk
[^conc1]: Cloud concentration ~63% AWS/Azure/GCP; >65% of EU FEs use ≥2; DORA list first concrete act. Bristows; Lexology; FStech — law-firm/press. https://inquisitiveminds.bristows.com/post/102hg8t/cloud-concentration-systemic-risk ; https://www.fstech.co.uk/fst/EU_Names_19_Critical_Tech_Providers_Under_DORA_Including_AWS_Google_And_Microsoft.php
[^conc2]: CrowdStrike 19 Jul 2024 (~8.5M machines, ~70% Fortune 500); single-vendor not hyperscaler; FCA lessons 31 Oct 2024. US GAO; CSA; Burges Salmon — regulator-primary/standards-body/law-firm. https://www.gao.gov/products/gao-24-107733 ; https://www.burges-salmon.com/our-thinking/crowdstrike-incident-regulatory-implications-and-lessons-learned/
[^conc3]: Exit strategy / stressed exit / substitutability / step-in (EBA Guidelines on Outsourcing Feb 2019; DORA Art. 30 transition). Pinsent Masons; DLA Piper; EBF — law-firm/industry-body. https://www.pinsentmasons.com/out-law/analysis/exit-planning-what-outsourcing-banks-need-to-know ; https://www.ebf.eu/wp-content/uploads/2020/09/Cloud-exit-strategy-Testing-of-exit-plans.pdf
[^conc4]: Fourth-party / sub-contracting chain — banks look past direct vendor into the ICT supply chain. GloCert; DORA Art. 30 — practitioner/official-text (QUALIFIED). https://www.glocertinternational.com/resources/guides/dora-ict-third-party-risk-contracting-and-exit/
[^nuance-conc]: Regulators (UK/FSB) say concentration not *automatically* systemic. Lexology (UK CP summary); FSB — law-firm/regulator (disconfirming). https://www.lexology.com/library/detail.aspx?g=99f48f00-94c2-45ba-b94a-2d7e0fca51cc ; https://www.fsb.org/uploads/PIFS.pdf
[^nuance-exit]: Regulators do not expect "perfect" exit plans (multi-cloud complexity). Microsoft Cloud Blog — vendor (TENTATIVE, single strong source). https://www.microsoft.com/en-us/microsoft-cloud/blog/banking/2026/02/02/managing-concentration-risk-and-exit-requirements-a-framework-for-financial-institutions/ — verified-as-of 2026-06-18
[^nuance-skeptic]: "Regulation alone unlikely to reduce concentration"; hyperscalers may be more resilient. Bristows; Carnegie Endowment Cloud Governance — law-firm/think-tank (disconfirming). https://inquisitiveminds.bristows.com/post/102hg8t/cloud-concentration-systemic-risk ; https://cloud.carnegieendowment.org/cloud-governance-issues/effects-of-cloud-market-concentration/
[^nuance-burden]: DORA prescriptive/burdensome; fines up to 2% turnover (entities) / 1% avg daily turnover (CTPPs). Infosecurity Magazine; Protiviti — press/consultancy (disconfirming). https://www.infosecurity-magazine.com/news/dora-financial-firms-compliance/ ; https://www.protiviti.com/us-en/whitepaper/dora-compliance-untangling-key-hurdles-implementation
