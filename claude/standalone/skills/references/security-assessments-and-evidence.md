<!-- Provenance: reference under the `enterprise-vendor-management-and-tprm` skill. Created 2026-06-18 via /dr deep-research. Buyer-side governance lens for a B2B/SaaS seller (MongoDB TAM) selling into banks. Educational context, not legal/compliance advice. -->

# Standardized Security Assessments & Evidence (the "Questionnaire Gauntlet")

The standardized questionnaires, certifications, and attestations a bank's TPRM team puts a vendor through,
and what each is actually worth. This is the part of buyer-side governance a seller's security/SA team meets
most often. For *MongoDB's own* certification scope, see `mongodb-compliance` (under `mongodb-operations-expert`);
this reference covers what a **buyer** asks for and does with the evidence. Claims dated **as of 2026** where
volatile; cited inline as `[^n]`.

## Contents
- Shared Assessments: SIG, SIG Lite, SCA
- Cloud Security Alliance: CAIQ, CCM, STAR
- SOC reports: SOC 1 vs 2 vs 3, Type I vs II, bridge letters, CUECs
- ISO/IEC 27001:2022 supplier controls and ISO/IEC 27036
- Penetration-test attestation letters
- Trust centers & questionnaire automation
- Disconfirming nuance (the "rubber stamp" problem)
- References

---

## Shared Assessments: SIG, SIG Lite, SCA

The **Shared Assessments Program**'s **SIG (Standardized Information Gathering)** questionnaire is the dominant
standardized third-party questionnaire, **refreshed annually** (the first Thursday in November).[^sig1]

- **Question counts (2025 release):** **SIG Lite ≈ 128**, **SIG Core ≈ 627**, **SIG Detail (full) ≈ 1,936**,
  spanning **21 risk domains across 4 control areas.** The 2025 update folded in mappings to DORA, NIS2, and
  NIST CSF 2.0.[^sig1] *(Counts revise annually — re-verify.)*
- **Core vs Lite = depth tiering.** SIG Lite is the shorter set for lower-risk vendors; SIG Core is the deeper
  questionnaire for higher-risk vendors / those with sensitive-data access. **Which one a vendor receives maps
  directly to the buyer's tiering decision** (see `references/tprm-lifecycle-and-risk-tiering.md`).[^sig2]
- **The SCA (Standardized Control Assessment)** is the verification/attestation layer *above* the SIG — an
  on-site or virtual procedure that *tests* controls rather than accepting self-reported answers; the
  "trust but verify" component.[^sig3]

## Cloud Security Alliance: CAIQ, CCM, STAR

The CSA's artifacts are the **cloud-specific** standardized assessment, highly relevant to a cloud service
like MongoDB Atlas:[^caiq1]

- **Cloud Controls Matrix (CCM) v4** — **197 control objectives across 17 domains** (the control framework).
- **Consensus Assessments Initiative Questionnaire (CAIQ) v4** — **261 questions** mapped to those 17 domains
  (the questionnaire form of the CCM), historically yes/no answers documenting IaaS/PaaS/SaaS controls.
- **STAR registry** has **two assurance levels**: **Level 1 = self-assessment** (vendor submits a completed
  CAIQ — free, publicly listed); **Level 2 = third-party audit / attestation** (independent certification,
  e.g., STAR Certification built on ISO 27001, or STAR Attestation built on SOC 2).[^caiq1] **Seller
  implication:** a Level 2 STAR attestation is far stronger evidence than a Level 1 self-assessment.

## SOC reports: SOC 1 vs 2 vs 3, Type I vs II, bridge letters, CUECs

The AICPA **SOC** family is the most-requested attestation set — and the distinctions matter:[^soc1]

- **SOC 1** = controls relevant to a customer's **financial reporting** (ICFR).
- **SOC 2** = **security / trust** controls, based on the **Trust Services Criteria** (Security/"common
  criteria" is mandatory; Availability, Processing Integrity, Confidentiality, Privacy are optional add-ons).
  Private, NDA-gated, contains detailed test results — **this is what a bank's TPRM team actually wants.**
- **SOC 3** = a **public**, general-use summary of a SOC 2 (always Type II; opinion + assertion + system
  description, but **no detailed test results**) — the shareable teaser.

The numbers are **different reports, not a maturity ladder.**[^soc1]

- **Type I vs Type II:** **Type I** = control *design* at a single point in time; **Type II** = design *and
  operating effectiveness* over a period (typically 6–12 months). **Type II is the de-facto bank standard;
  many buyers now reject Type I.**[^soc1]
- **Bridge / gap letter:** covers the period between a SOC 2's report end-date and the buyer's review date. It
  is a **management assertion (NOT an independent audit)**, typically covering up to ~3 months, summarizing
  whether anything material changed. Buyers and auditors know it carries **less weight** than the audited
  report.[^soc2]
- **What a TPRM team does with a SOC 2:** reads the **auditor opinion type** (unqualified vs qualified),
  checks the **scope / Trust Services Criteria covered**, scrutinizes the **exceptions/deviations** noted,
  confirms the **report period** (and whether a bridge letter fills the gap to today), and reviews the
  **Complementary User Entity Controls (CUECs)** — the controls the *bank itself* must operate for the
  vendor's controls to be effective.[^soc3] *(CUEC review is standard practice; treat as `QUALIFIED`.)*

## ISO/IEC 27001:2022 supplier controls and ISO/IEC 27036

- **ISO/IEC 27001 (the ISMS certification)** is the global counterpart to SOC 2. **As of 2026 the current
  edition is ISO/IEC 27001:2022**, and the transition deadline from the 2013 version was **October 2025** — so
  a vendor presenting a 2013-edition certificate is stale; flag it. Supplier-relationship controls now live in
  Annex A organizational controls **A.5.19–A.5.23** (supplier relationships; security in supplier agreements;
  ICT supply-chain security; monitoring/review of supplier services; cloud-service security).[^iso1]
- **ISO/IEC 27036 ("Cybersecurity — Supplier relationships")** is the supplier-relationship-*specific* standard
  a procurement/risk team may cite. It is a four-part standard: **Part 1** overview & concepts (27036-1:2021),
  **Part 2** requirements (27036-2:2022), **Part 3** ICT supply-chain guidance, **Part 4** cloud-service
  security. It defines requirements for *both* acquirer and supplier to define, operate, and monitor the
  relationship — i.e., it formalizes "how a buyer should govern a vendor like MongoDB."[^iso2]

## Penetration-test attestation letters

Buyers usually accept a **pen-test attestation / summary letter, not the full report.** The letter confirms
testing occurred, gives high-level scope and outcome, and states whether findings were remediated — *without*
exposing exploit paths, IPs, payloads, or screenshots (which would create risk if circulated). The maxim:
**"the letter is for showing; the full report is for fixing."** The same pattern applies to vulnerability-scan
evidence — summary over raw output.[^pentest1] **Seller implication:** a vendor should have a current
attestation letter ready; offering the full report is rarely necessary and usually unwise.

## Trust centers & questionnaire automation

The vendor-side answer to the gauntlet, and a recognized **sales-cycle accelerator**:[^trust1]

- A **trust center** (SafeBase, Whistic Trust Center Exchange, Conveyor, Vanta) is a self-service portal where
  prospects pull SOC 2 / ISO 27001 / pen-test letters / CAIQ on demand — reducing inbound questionnaires
  (SafeBase claims 74%+ reduction).
- **Questionnaire-response automation** (AI-assisted, drawing from the trust center / knowledge base) handles
  the formal questionnaires that still arrive, cutting turnaround "from days to minutes."
- They are **complementary, not substitutes**: trust centers reduce *volume*; automation accelerates the
  *remainder*. A complete, current, audit-grade evidence package is both a sales accelerator and the antidote
  to the "rubber stamp" critique below.[^trust1]

## Disconfirming nuance (the "rubber stamp" problem)

This is the single most important strategic insight in this reference, and it cuts both ways for a seller:

- **SOC 2 ≠ "secure."** Industry sources are explicit: SOC audits "are not a silver bullet and do not
  guarantee a breach won't happen"; evidence reflects a *point in time* and configurations drift afterward.
  Do not let a customer (or your own org) treat a SOC 2 as proof of security — it is scoped assurance of
  *controls during a window*.[^nuance-soc2]
- **Security questionnaires are widely seen as a low-trust rubber stamp.** Hard numbers: only **~34% of TPRM
  professionals believe questionnaire responses**, yet **~81% of organizations report that the vast majority
  of returned questionnaires claim *perfect* compliance.** Self-attestations are "rarely conducted by an
  outside auditor." Questionnaire fatigue drives copy-paste answers; reviewers lack capacity to scrutinize
  hundreds of responses. **CSA itself published a 2025 piece calling questionnaires "a familiar but ineffective
  norm."**[^nuance-quest]
- **The gauntlet is real and fragmented.** TPRM teams juggle *multiple parallel* standards (SIG, CAIQ, plus
  bespoke customer-proprietary forms); no single standard has won, so **a bank may still send MongoDB its own
  custom questionnaire regardless of available SOC 2 / ISO / CAIQ evidence.**[^nuance-frag]

**Net strategic takeaway for selling into banks:** the buyer-side trend (as of 2026) is *skepticism of
self-attested questionnaires* and a *shift toward verified evidence* (SOC 2 Type II, ISO 27001:2022, STAR
Level 2, pen-test attestation letters) plus *continuous outside-in monitoring* (see
`references/tprm-lifecycle-and-risk-tiering.md`). A maintained, audit-grade evidence package in a trust center
is the best lever a seller has.

---

## References

[^sig1]: SIG 2025 version, counts (Lite 128 / Core 627 / Detail 1,936), annual cadence, 21 domains, DORA/NIS2/CSF mappings. Shared Assessments — standards-body. https://sharedassessments.org/blog/2025-sig/ ; https://sharedassessments.org/about-sig/ — verified-as-of 2026-06-18
[^sig2]: SIG Core vs Lite depth tiering. UpGuard; WorkStreet — vendor-doc. https://www.upguard.com/blog/sig-questionnaire ; https://www.workstreet.com/blog/sig-lite
[^sig3]: SCA (Standardized Control Assessment) as verification layer above SIG. Shared Assessments; CentralEyes — standards-body/vendor-doc. https://sharedassessments.org/about-sig/ ; https://www.centraleyes.com/sig-security-questionnaire/
[^caiq1]: CCM v4 (197 objectives/17 domains), CAIQ v4 (261 questions), STAR Level 1 self-assessment vs Level 2 attestation. Cloud Security Alliance — standards-body. https://cloudsecurityalliance.org/artifacts/cloud-controls-matrix-v4 ; https://cloudsecurityalliance.org/research/topics/caiq ; https://cloudsecurityalliance.org/blog/2024/02/17/the-csa-cloud-controls-matrix-and-consensus-assessment-initiative-questionnaire-faqs — verified-as-of 2026-06-18
[^soc1]: SOC 1/2/3 distinctions; Type I vs II; Trust Services Criteria (Security mandatory); SOC 3 public/no test results. Secureframe; LinfordCo (CPA) — vendor-doc/practitioner. https://secureframe.com/hub/soc-2/soc-1-vs-soc-2-vs-soc-3 ; https://linfordco.com/blog/trust-services-critieria-principles-soc-2/ ; https://linfordco.com/blog/soc-2-vs-soc-3/
[^soc2]: Bridge/gap letter = management assertion (not an audit), ~3 months, less weight. StrikeGraph; Iris — vendor-doc. https://www.strikegraph.com/blog/what-is-a-bridge-letter-in-a-soc-2-report ; https://heyiris.ai/blog/what-is-a-soc-2-bridge-letter
[^soc3]: What a TPRM team does with a SOC 2 (opinion, scope, exceptions, period, CUECs). LinfordCo; Secure Controls Framework — practitioner (CUEC review QUALIFIED). https://linfordco.com/blog/busting-soc-2-myths/ ; https://securecontrolsframework.com/grc-fundamentals/common-cybersecurity-frameworks/trust-services-criteria-soc-2
[^iso1]: ISO/IEC 27001:2022 current; Oct-2025 transition deadline; Annex A 5.19–5.23 supplier controls. ISMS.online; URM; HighTable — vendor-doc/practitioner. https://www.isms.online/iso-27001/annex-a-2022/5-19-information-security-supplier-relationships-2022/ ; https://www.urmconsulting.com/blog/iso-27001-2022-a-5-organisational-controls-supplier-management ; https://hightable.io/iso-27001-annex-a-controls-reference-guide/ — verified-as-of 2026-06-18
[^iso2]: ISO/IEC 27036 four-part supplier-relationship standard (acquirer + supplier requirements). ISO; IEC; HCLTech — standards-body/practitioner. https://www.iso.org/standard/82060.html ; https://www.iec.ch/blog/managing-supplier-relationships-cyber-security ; https://www.hcltech.com/blogs/iso-27036-provides-standardized-framework-acquirer-and-supplier-relationships
[^pentest1]: Pen-test attestation letter vs full report ("letter for showing, report for fixing"). Packetlabs; Cybri; Adversis — vendor-doc. https://www.packetlabs.net/posts/letters-of-attestation ; https://cybri.com/blog/a-guide-to-penetration-test-letters-of-attestation/ ; https://www.adversis.io/blogs/attestation-letter-assessment-summary-guidance
[^trust1]: Trust centers + questionnaire automation as complementary sales accelerators; 74%+ inbound reduction claim. SafeBase; The Hacker News; Tribble — vendor-doc/practitioner. https://safebase.io/resources/security-questionnaires ; https://thehackernews.com/2024/07/how-trust-center-solves-your-security.html ; https://tribble.ai/blog/best-ai-trust-center-security-portal-platforms-2026/ — verified-as-of 2026-06-18
[^nuance-soc2]: SOC 2 ≠ guarantee of security; point-in-time, config drift. LinfordCo (CPA); Secureframe — practitioner (disconfirming). https://linfordco.com/blog/busting-soc-2-myths/ ; https://secureframe.com/hub/soc-2/what-is-soc-2
[^nuance-quest]: ~34% believe questionnaire responses / ~81% claim perfect compliance; CSA "ineffective norm." Vanta; CSA; RiskRecon — vendor-doc/standards-body (disconfirming). https://www.vanta.com/resources/security-questionnaires-are-ineffective ; https://cloudsecurityalliance.org/blog/2025/04/02/why-security-questionnaires-are-a-familiar-but-ineffective-norm-for-assessing-risk ; https://blog.riskrecon.com/the-problem-with-security-questionnaires — verified-as-of 2026-06-18
[^nuance-frag]: Fragmented multi-questionnaire fatigue; bespoke forms persist despite available standards. SecurityScorecard; CSA — vendor-doc/standards-body (disconfirming). https://securityscorecard.com/blog/the-four-questionnaires-your-tprm-team-is-managing-and-struggling-to-keep-up-with/ ; https://cloudsecurityalliance.org/blog/2025/04/02/why-security-questionnaires-are-a-familiar-but-ineffective-norm-for-assessing-risk
