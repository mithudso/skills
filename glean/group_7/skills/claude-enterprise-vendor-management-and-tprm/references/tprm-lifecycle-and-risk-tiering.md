<!-- Provenance: reference under the `enterprise-vendor-management-and-tprm` skill. Created 2026-06-18 via /dr deep-research. Buyer-side governance lens for a B2B/SaaS seller (MongoDB TAM) selling into banks. Educational context, not legal/compliance advice. -->

# TPRM Lifecycle & Vendor Risk Tiering

How a large enterprise — especially a bank — runs **third-party risk management (TPRM)** over a vendor like
MongoDB: the lifecycle stages, how the vendor gets classified by criticality, what risk domains are assessed,
who owns the program, and how the vendor is monitored between assessments. This is the *buyer's* process; the
seller's job is to anticipate it. Claims are dated **as of 2026** where volatile; cited inline as `[^n]`.

## Contents
- The TPRM lifecycle (closed loop)
- Vendor risk tiering / criticality classification
- Inherent vs residual risk
- Risk domains assessed
- Fourth-party & concentration risk
- The TPRM operating model (three lines of defense)
- GRC / VRM tooling landscape
- Continuous monitoring & security-ratings services
- Disconfirming nuance
- References

---

## The TPRM lifecycle (closed loop)

TPRM is a **closed loop**, typically decomposed into **five to six named stages**; vendors collapse or expand
the count but the substance is consistent:[^life1]

1. **Planning / scoping** — identify the need, the data/access involved, and the inherent risk before
   selecting anyone.
2. **Due diligence & vendor selection** — assess candidate vendors (security, financial, resilience,
   compliance) proportionate to the inherent risk.
3. **Contracting** — negotiate risk-allocating terms (audit rights, security, sub-processors, exit). *Clause
   drafting is out of scope here → `legal-adjacent-writing`; the contract-instrument map is in
   `references/contract-instruments-map.md`.*
4. **Onboarding** — provision access, register the vendor, baseline controls.
5. **Ongoing monitoring** — reassess on a tier-driven cadence; watch continuous signals; manage findings.
6. **Offboarding / termination / exit** — wind down access, retrieve/destroy data, execute the exit plan.

For **US banks**, the authoritative spine for this lifecycle is the **June 2023 Interagency Guidance on
Third-Party Relationships: Risk Management** (Federal Reserve + FDIC + OCC, finalized June 6 2023; Federal
Register June 9 2023), which frames sound risk management "for all stages in the life cycle of third-party
relationships" and **replaced** the agencies' prior siloed guidance (OCC 2013-29, Fed SR 13-19, FDIC
2008).[^life2] It is **principles-based and risk-tiered, not prescriptive** — depth of diligence scales with
"the level of risk, complexity, and size" of the relationship — so two banks can run materially different
diligence on the same vendor.[^life3] (The *regulatory* framing of this and its operational-resilience
cousins is detailed in `references/operational-resilience-regulation.md`.)

## Vendor risk tiering / criticality classification

Tiering is the master control: **the tier dictates assessment depth, contractual requirements, and monitoring
cadence.** A common four-tier model:[^tier1]

| Tier | Typical trigger | Assessment depth | Monitoring cadence |
|---|---|---|---|
| **Tier 1 / Critical** | Access to highly sensitive data or business-critical systems | Comprehensive assessment + on-site/virtual audit + pen-test evidence + continuous monitoring | Continuous + at least annual deep review |
| **Tier 2 / High** | Material data/system access | Detailed questionnaire (e.g., SIG Core) + security ratings | ~Quarterly |
| **Tier 3 / Medium** | Limited access | Lightweight questionnaire (SIG Lite) + certificate review | Semi-annual / annual |
| **Tier 4 / Low** | Minimal access / commodity | Basic due diligence | Annual or event-triggered |

Tier is driven by **data sensitivity, system/data access, business criticality, and substitutability**.[^tier2]
A crucial nuance the best programs respect: **"criticality" ≠ "risk."** A vendor can be operationally critical
yet low-risk (well-controlled), or high-risk yet easily replaceable. Criticality drives *scope/depth,
monitoring frequency, contractual requirements, and escalation*; conflating it with risk into one number
distorts the program.[^tier2] **Seller implication:** which tier a bank assigns MongoDB determines whether it
gets a 128-question SIG Lite or a 600+-question SIG Core plus on-site audit (see
`references/security-assessments-and-evidence.md`).

## Inherent vs residual risk

The core scoring mechanic. **Inherent risk** = the baseline risk from the vendor's role and data access
*before* any controls are credited. **Residual risk** = what remains *after* the vendor's controls are
credited (conceptually, residual ≈ inherent × control effectiveness).[^resid1] Tiering usually keys off
**inherent** risk (you can score it before you've seen the vendor's controls); the assessment then evaluates
controls to derive **residual** risk, which drives the accept/mitigate/reject decision.[^resid1]

## Risk domains assessed

A bank's assessment is **multi-dimensional, not just cybersecurity**.[^dom1] Typical domains:

- **Information / cyber security** (the most visible, but only one axis)
- **Financial viability** (will the vendor still exist in 3 years?)
- **Business continuity / operational resilience** (can it withstand and recover from disruption?)
- **Compliance / regulatory** (does it meet the obligations the bank must flow down?)
- **Concentration** (is the bank over-dependent on this vendor or its upstream?)
- **Reputational** (does association create reputational risk?)
- **Fourth-party / subcontractor** (what about the vendor's vendors?)

**Seller implication:** a strong security posture is necessary but not sufficient — a bank will also probe
MongoDB's financial health, BC/DR, sub-processors, and concentration profile.

## Fourth-party & concentration risk

These are explicit banking concerns, increasingly hard contractual demands:[^fourth1]

- **Fourth-party (Nth-party) risk** — your vendor's subcontractors. For MongoDB Atlas, the most salient
  fourth parties are the underlying hyperscalers (AWS / Azure / GCP). Banks increasingly require the primary
  vendor to **disclose material sub-processors** and accept **flow-down obligations** onto them.[^fourth1]
- **Concentration risk** — many institutions depending on the same upstream provider, creating *systemic*
  exposure. Vendor diversification does not help if the diversified vendors all sit on the same hyperscaler.
  This is the bridge to the operational-resilience regulatory frontier (DORA's CTPP oversight, etc. — see
  `references/operational-resilience-regulation.md`).[^fourth1]

## The TPRM operating model (three lines of defense)

Banks run TPRM through a **dedicated TPRM / vendor-risk office plus a business "relationship owner,"** mapped
onto the **three-lines-of-defense** model:[^model1]

- **First line** — the business / relationship owner who uses the vendor and owns the risk day-to-day.
- **Second line** — the TPRM / risk-management function that sets policy, runs assessments, and challenges.
- **Third line** — internal audit, providing independent assurance.

*(The three-lines mapping is widely-described standard practice; treat as `QUALIFIED` rather than single-source
fact.)*[^model1] **Seller implication:** the person sending MongoDB a questionnaire (second line) is usually
*not* the person who wants the product (first line). The first-line relationship owner is the seller's ally on
value and urgency; the second line is paid to find risk and slow things down.

## GRC / VRM tooling landscape

Banks run TPRM on **GRC / vendor-risk-management (VRM) platforms.** Named examples **as of 2026** (verify —
this market consolidates fast): **OneTrust, ServiceNow VRM, ProcessUnity, Archer, Venminder**, and others.
**Vendor-fact note:** **Prevalent was acquired by Mitratech in October 2024** and is now marketed under
Mitratech — so a "standalone Prevalent" reference is dated.[^tool1] Archer is the long-standing configurable
enterprise-GRC backbone; ProcessUnity is a dedicated VRM tool; OneTrust / ServiceNow are broader platforms.[^tool1]
Treat the precise roster as volatile and re-verify before quoting.

## Continuous monitoring & security-ratings services

**Continuous monitoring ≠ point-in-time assessment.** Between the annual questionnaire / SOC 2 snapshot,
banks subscribe to **security-ratings services** — **BitSight, SecurityScorecard, RiskRecon** — that generate
**"outside-in" scores from externally observable data, requiring no input from the rated company.** BitSight,
for example, cites 120+ data sources and ~25 risk vectors across four categories (compromised systems,
security diligence, user behavior, public breach disclosures); scores update automatically as issues are
remediated.[^mon1] This gives the bank a *between-assessment* signal a once-a-year questionnaire cannot.

**Seller implication:** a bank may ding MongoDB on an outside-in rating the seller did not even know was being
computed — and (see nuance) the ding may be a false positive or a shared-cloud (fourth-party) artifact.

## Disconfirming nuance

- **Security ratings have a credibility problem.** CISOs have historically questioned their value due to high
  **false-positive rates** and limited ROI; documented cases include scores depressed by false positives that
  need manual correction, and "thousands of false positives" on cloud auto-scaling / load-balancer setups.
  Breach-correlation validation is contested. So an outside-in score is a *signal*, not ground truth — and a
  seller can legitimately contest a rating tied to shared cloud infrastructure.[^nuance-ratings]
- **Criticality and inherent risk are routinely conflated**, distorting tiering if scored as one number.[^tier2]
- **The post-2023 US banking guidance is deliberately non-prescriptive** — there is no single mandated
  questionnaire or checklist, so diligence depth varies widely bank to bank.[^life3]
- **Self-attested questionnaires are low-trust** (only ~34% of TPRM pros believe them) — covered in depth in
  `references/security-assessments-and-evidence.md`.

---

## References

[^life1]: TPRM lifecycle as a 5–6 stage closed loop incl. offboarding. Aravo; Apexanalytix; Neotas — vendor-doc. https://aravo.com/blog/six-phases-of-the-tprm-lifecycle/ ; https://www.apexanalytix.com/resources/blog/third-party-risk-management-lifecycle/ ; https://www.neotas.com/third-party-risk-management-tprm-lifecycle/
[^life2]: June 2023 Interagency Guidance on Third-Party Relationships: Risk Management; lifecycle framing; replaced prior guidance. OCC Bulletin 2023-17; Federal Register; FDIC — regulator-primary. https://www.occ.gov/news-issuances/bulletins/2023/bulletin-2023-17.html ; https://www.federalregister.gov/documents/2023/06/09/2023-12340/interagency-guidance-on-third-party-relationships-risk-management ; https://www.fdic.gov/news/press-releases/2023/pr23047a.pdf — verified-as-of 2026-06-18
[^life3]: 2023 guidance is principles-based, risk-tiered, non-prescriptive; reaches subcontractors. InnReg; Aravo unpacking — practitioner/vendor-doc. https://www.innreg.com/blog/interagency-guidance-third-party-risk-management ; https://aravo.com/blog/unpacking-the-final-interagency-guidance-on-third-party-relationships-risk-management/
[^tier1]: Four-tier vendor model; assessment depth + monitoring cadence by tier. UpGuard; Mitratech; ProcessUnity — vendor-doc. https://www.upguard.com/blog/what-is-vendor-tiering ; https://mitratech.com/resource-hub/blog/third-party-risk-scoring-tiering/ ; https://www.processunity.com/third-party-risk-management/inherent-risk/
[^tier2]: Tier drivers (data sensitivity, access, criticality, substitutability); criticality ≠ risk. Atlas Systems; Mitratech — vendor-doc. https://www.atlassystems.com/blog/how-to-determine-vendor-criticality ; https://mitratech.com/resource-hub/blog/third-party-risk-scoring-tiering/
[^resid1]: Inherent vs residual risk; residual ≈ inherent × control effectiveness; tiering keys off inherent. ProcessUnity — vendor-doc. https://www.processunity.com/resources/blogs/vendor-inherent-risk-residual-risk/ ; https://www.processunity.com/third-party-risk-management/inherent-risk/
[^dom1]: Multi-dimensional risk domains (security, financial, BC, compliance, concentration, reputational, fourth-party). Onspring; Safe Security; Hellios — vendor-doc. https://onspring.com/resources/guide/guide-what-is-third-party-risk-management-tprm/ ; https://safe.security/resources/blog/2026-guide-to-third-party-risk-management-tprm/ ; https://hellios.com/third-party-risk-management-what-it-is-and-how-to-build-it
[^fourth1]: Fourth-party / Nth-party and concentration risk in banking; sub-processor disclosure + flow-down. Aravo; RiskLedger; Ncontracts — vendor-doc. https://aravo.com/blog/tprm-fourth-party-risk-management-and-concentration-risk-in-banking/ ; https://riskledger.com/resources/securing-your-suppliers-suppliers ; https://www.ncontracts.com/nsight-blog/managing-fourth-party-risk-what-you-need-to-know
[^model1]: TPRM operating model + three-lines-of-defense (QUALIFIED — widely described). Onspring; Hellios — vendor-doc. https://onspring.com/resources/guide/guide-what-is-third-party-risk-management-tprm/ ; https://hellios.com/third-party-risk-management-what-it-is-and-how-to-build-it
[^tool1]: GRC/VRM tooling landscape; Prevalent acquired by Mitratech Oct 2024. Mitratech/GlobeNewswire press; BitSight FI platform guide — vendor-doc/press. https://mitratech.com/resource-hub/pressreleases/acquisition-of-preparis-and-prevalent/ ; https://www.bitsight.com/guides/best-third-party-risk-management-platforms-financial-institutions — verified-as-of 2026-06-18
[^mon1]: Security-ratings services (BitSight/SecurityScorecard/RiskRecon) outside-in scoring; 120+ sources / ~25 vectors. BitSight; UpGuard comparison — vendor-doc. https://www.bitsight.com/security-ratings ; https://www.upguard.com/compare/bitsight-vs-securityscorecard — verified-as-of 2026-06-18
[^nuance-ratings]: Security-ratings false-positive / ROI / breach-correlation criticism. BankInfoSecurity; UpGuard — press/vendor-doc (disconfirming). https://www.bankinfosecurity.com/bitsight-securityscorecard-panorays-lead-risk-ratings-tech-a-25326 ; https://www.upguard.com/compare/bitsight-vs-securityscorecard
