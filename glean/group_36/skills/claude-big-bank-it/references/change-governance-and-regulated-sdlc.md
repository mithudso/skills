<!-- Provenance: reference under the `big-bank-IT` skill. Created 2026-06-18 via /dr deep research (ITIL/Axelos + ISACA/COBIT + FFIEC + FCA + DORA/Accelerate + vendor-neutral engineering sources). Educational, vendor-neutral technical orientation — NOT advice. The two "DORA"s are disambiguated throughout. Disconfirming evidence preserved. -->

# Change Governance, Change-Freeze Windows, Segregation of Duties & the Regulated SDLC

`verified-as-of: 2026-06-18` (compliant-pipeline tooling and the DevOps-metrics research move; re-verify tooling claims).

> **Educational, vendor-neutral orientation, NOT advice.** SOX/regulatory specifics are pointer-level here (the control objective and its IT realization); the regulatory *map* → `fsi-banking-regulatory-context`; DR/failover mechanics → `resiliency-dr-and-sre.md`.
>
> **TERMINOLOGY GUARDRAIL — the two "DORA"s.** Throughout, **"DORA" without qualification = DevOps Research & Assessment** (the Google/Forsgren research behind the *Accelerate* book and the four key metrics). The **EU Digital Operational Resilience Act (also "DORA")** is a *regulation*, referenced only where named as a compliance regime — and disambiguated inline each time.

## Contents

- [Change management & release governance (ITIL, CAB, RFCs)](#change-management--release-governance-itil-cab-rfcs)
- [Change-freeze / blackout windows](#change-freeze--blackout-windows)
- [Segregation of duties (SoD), maker-checker, PAM, break-glass](#segregation-of-duties-sod-maker-checker-pam-break-glass)
- [The regulated SDLC — change controls as SOX ITGCs](#the-regulated-sdlc--change-controls-as-sox-itgcs)
- [Audit trail / evidence for examiners](#audit-trail--evidence-for-examiners)
- [The core tension: deployment frequency vs bank change control](#the-core-tension-deployment-frequency-vs-bank-change-control)
- [Disconfirming evidence](#disconfirming-evidence)
- [Seller takeaways](#seller-takeaways)
- [Sources](#sources)

## Change management & release governance (ITIL, CAB, RFCs)

**ITIL classifies every change into three types, and the type determines the approval path.** A **standard change** is "low-risk, pre-authorized … well understood and fully documented … implemented without needing additional authorization"; a **normal change** is assessed/approved/scheduled by risk; an **emergency change** "must be introduced as soon as possible" (e.g. a security patch or major-incident fix).[^3][^4] ITIL 4 renamed the practice "change enablement" to signal it should help changes move forward, not just gate them.[^4] *(Confidence: fact.)*

**The Change Advisory Board (CAB)** assesses and authorizes normal changes (an Emergency CAB / ECAB handles urgent ones); a **Request for Change (RFC)** is the ticket of record. Notably, ITIL 4 itself now frames the CAB as **advisory, not a blocker** — it "focuses on risk, conflicts, and readiness."[^1][^4] *(Confidence: fact.)*

**Why a bank's deployment cadence is slow and gated:** each production change must be ticketed, risk-assessed, approved by someone *other than the requester*, scheduled into an approved window, tested in non-production, and evidenced — these steps are mandated as IT general controls, not merely best practice (below). Production access is restricted: developers typically cannot push to production directly.[^6][^17] *(Confidence: fact.)*

## Change-freeze / blackout windows

**A change freeze (blackout window, code freeze, moratorium) is a scheduled period in which production deployments are blocked, imposed at times of elevated business risk.** Tooling enforces it: ServiceNow models them as "Blackout" schedule exceptions; GitLab exposes a `$CI_DEPLOY_FREEZE` variable so pipelines auto-pause; freeze policies are commonly defined as code with cron windows + explicit break-glass approver lists.[^11][^12][^13] *(Confidence: fact.)*

**Banks/financial firms freeze around quarter/year-end financial close, peak retail (Black Friday/holiday), tax/regulatory cutoffs, and market-sensitive periods.** "A banking system might enforce a complete lockdown during end-of-year financial processing"; Mastercard's payments arm runs a holiday code freeze to "ensure payment gateway stability during peak shopping days." Durations of ~72–96 hours are common; emergency changes are explicitly carved out so a freeze never blocks an outage fix.[^11][^12][^33][^34] **Why:** maximize stability when the cost of an incident is highest by removing the largest source of incidents — change itself. *(Confidence: fact; "~72–96h norm" is practitioner consensus, not a standard — qualified.)*

## Segregation of duties (SoD), maker-checker, PAM, break-glass

**SoD = "the principle that no user should be given enough privileges to misuse the system on their own,"** enforced statically (conflicting roles) or dynamically (at access time). The classic model separates **authorization, custody, recordkeeping, and reconciliation/verification.**[^16][^17] *(Confidence: fact.)*

**The IT realization is specific and auditor-checked:** "developers cannot migrate code to production, DBAs cannot access business applications … access administrators cannot approve their own access requests." The canonical control: **"no single employee can unilaterally deploy a code change into production"**; "the person who writes the SQL cannot be the same person who approves it."[^17][^18] *(Confidence: fact.)*

**Maker-checker / four-eyes is the banking-specific expression — but the terms aren't synonyms** (terminology genuinely varies): **maker-checker / four-eyes** is an *atomic* control (maker initiates, checker approves before execution); **SoD** is *broader/structural* (a continuous process split into phases across the lifecycle); **dual control** = two people act *simultaneously* on one highly sensitive task. "Role separation [SoD] is the design foundation; two-person review and maker-checker enforce it at a specific moment."[^19][^20][^21] *(Confidence: qualified.)*

**Separation of dev/test/prod and least privilege are the enabling controls** ("developers have full access to dev/test, but no direct access to production").[^6] **PAM** (Privileged Access Management) and **break-glass** govern residual privileged access — the modern pattern is **just-in-time (JIT)** access (time-limited, task-specific, revoked after use → "zero standing privilege"), and **break-glass** emergency access "must be heavily logged, time-limited, and reviewed after use" — "a break-glass pathway with automatic logging and follow-up review," not "a back door around process."[^25][^36][^37][^38] *(Confidence: fact.)*

**Tie to SOX 404 ITGC:** SoD is one of the core IT general control (ITGC) families (logical access, change management, SoD, IT operations) that "establish whether the systems supporting financial reporting can be trusted." Weak SoD causes auditors to "reduce reliance on system-dependent controls [and] expand substantive testing."[^28] (SOX statute specifics → `fsi-banking-regulatory-context`.) *(Confidence: fact.)*

## The regulated SDLC — change controls as SOX ITGCs

**Change & release controls map directly onto SOX ITGC domains; auditors trace any production change end-to-end.** The four domains an auditor checks for in-scope (financial-reporting) systems: **change management** (every change authorized before execution, linked to a ticket; writer ≠ approver), **access controls** (no shared admin; RBAC tied to an IdP; periodic access reviews), **audit trail** (full history of what ran, who requested/approved, justification, timestamp), and **program development** (changes flow dev→staging→prod with testing evidence; no direct-to-prod). COBIT anchors the same controls (change management = BAI06; SoD = DSS06.03). FFIEC's IT handbook spells out the bank-examination version.[^18][^24][^29][^31] *(Confidence: fact.)*

**The central adaptation for regulated DevOps: enforce separation by the pipeline, not by separate teams.** "While one person may not manually push a button to deploy their own code, the 'person' doing the deployment can be the automated pipeline itself" — shifting the control "from a manual, human-based separation to a systematic, automated one that is far more reliable and auditable."[^40] Mechanically: identity-verified commit → mandatory peer review (approver ≠ author) → automated build into an artifact repo no human can manually upload to → automated tests + security scans (SAST/SCA/SBOM) → change-approval gate → deployment by a peer-reviewed script, never by hand → complete commit-to-prod audit trail.[^16][^17][^40] **Maturity matters** — auditors test *enforcement*, not existence: "They do not ask whether pipelines exist. They ask whether pipelines enforce control."[^9] **Policy-as-code** (OPA/Rego, Kyverno) is the dominant evidence/enforcement mechanism; real banks evidence this (DBS embeds policy-as-code "guardrails"; Capital One translated SoD into "immutable stage gates" that auto-collect compliance evidence).[^42][^43][^44] *(Confidence: fact.)*

## Audit trail / evidence for examiners

**Every production change must be evidenced because examiners reconstruct it after the fact.** The audit-trail ITGC requires "the actual SQL statement, who requested it, who approved it, the business justification, and the timestamp." The auditor's test is sample-and-trace: "an auditor should be able to pick any change in production and easily trace [it back to] the approval ticket … in one click, not an archaeology project." SOC 2 frames the same under CC8.1, prizing "system-generated evidence over self-reported."[^17][^18][^10] **The "evidence-by-design" principle:** compliance artifacts (approval logs, scan results, deployment history, commit→build→artifact→prod traceability) should fall out of normal delivery as a byproduct.[^9][^10] **Immutability/WORM** is a hard requirement in parts of FS: SEC Rule 17a-4 requires WORM (or, since 2022, an "audit-trail alternative" maintaining "a complete time-stamped audit trail of all modifications").[^23] *(Confidence: fact.)*

## The core tension: deployment frequency vs bank change control

**DORA's (DevOps Research & Assessment) four key metrics** define delivery performance: deployment frequency, lead time for changes, change failure rate, time to restore service. DORA's repeated finding: **speed and stability are not a trade-off** — top performers excel at all four.[^1][^2][^7] *(Confidence: fact.)*

**DORA's most directly relevant — and most disconfirming — finding for banks:** external change approval (CAB / senior-manager sign-off) is **negatively correlated with delivery performance and shows *no* correlation with change failure rate.** From *Accelerate*: "approval by an external body (such as a manager or CAB) simply doesn't work to increase the stability of production systems … It is, in fact, worse than having no change approval process at all"; the 2019 report found such organizations **2.6× more likely to be low performers.** DORA's recommended substitute is **lightweight peer review at code check-in + a deployment pipeline** — explicitly framed as a *better* way to meet the SoD control objective.[^2][^7][^8] *(Confidence: fact — primary source + corroboration.)*

**This collides with the bank's reality**, where SOX ITGC, FFIEC examination, and (in the EU) the **Digital Operational Resilience Act** *require* demonstrable independent change authorization and SoD.[^18][^24][^9] The convergent resolution is **not to abandon the controls but to relocate them**: satisfy independent-approval and SoD *intent* via enforced peer review + automated gates in the pipeline (stronger, system-generated evidence than a CAB meeting), and reserve the human CAB for a strategic role (notification, coordination, trade-off decisions). DBS and Capital One are the canonical proof points that a regulated bank *can* reach high deployment frequency while remaining auditable.[^8][^43][^44] *(Confidence: fact / qualified.)*

## Disconfirming evidence

1. **CABs/external approval don't improve stability and slow delivery** — DORA's data shows external approval is *worse than no approval process* for stability and makes orgs 2.6× more likely to be low performers.[^2][^7][^8] (Practitioner summary: "Change Advisory Boards Don't Work.")[^22]
2. **Change freezes/moratoriums frequently *increase* risk by deferring it into a big-bang batch** — freezes mean "changes pile up … massive, unwieldy releases," a "large-batch anti-pattern for stability," and block security patches/bug fixes.[^47][^48][^49] *(Caveat: the critique targets routine/seasonal freezes; even critics concede event-scoped freezes can be justified when revenue/regulatory exposure genuinely peaks, and banks face hard regulatory cutoffs the pure-CD argument doesn't dissolve.)*
3. **SoD-as-separate-teams provides *weaker* assurance than proponents claim** — deployers who didn't write the code "cannot provide meaningful technical review … A developer intent on introducing a subtle bug … can satisfy all process controls while still achieving their goal"; SoD also "cannot address … collusion."[^16] The deeper point: the goal is that *no one person can unilaterally and unverifiably change production* — which automated pipelines satisfy without team silos.
4. **The friction is often a *design* failure, not inherent** — "security controls don't slow teams down by default; poorly designed security controls do." Applying the same heavyweight gate to low- and high-risk changes alike pushes engineers to batch changes (raising risk) or work around controls. ISACA's constructive version: where perfect SoD is infeasible, **compensating controls** (independent authorization/verification) are widely accepted — SoD is a risk-management dial, not an absolute.[^32]
5. **"You build it, you run it" (YBIYRI) vs SoD is genuinely contested** — proponents argue YBIYRI is compatible with SoD/PCI-DSS/SOX (it's "a myth that governance standards require [team-based] SoD"); the friction is real only where automation is immature (an emergency hotfix that should take 10 minutes takes 120 because the deploy button is behind a separate team).[^6] *(Confidence: qualified.)*

## Seller takeaways

- **Respect the control plane — don't pitch "go faster."** The credible message is "our deployment model *fits* your regulated SDLC and produces audit evidence" (system-generated approval logs, commit→prod traceability, policy-as-code-friendly).
- **Speak SoD and audit-trail fluently.** Know maker-checker, PAM/JIT, break-glass, and ITGC change-management — and that the modern bank enforces separation *in the pipeline*. A product whose change/access model maps to those controls clears a real procurement gate.
- **Anticipate change-freeze windows** in any go-live or migration plan — quarter/year-end and peak-period freezes are non-negotiable, and an emergency-change carve-out exists for genuine outages.

## Sources

[^1]: Change Enablement in One Page (ITIL 4) — https://itiligence.co.uk/change-enablement-in-one-page-itil-r-4/ — practitioner — ITIL change types; CAB-as-advisory.
[^2]: 2019 Accelerate State of DevOps Report — https://dora.dev/research/2019/dora-report/2019-dora-accelerate-state-of-devops-report.pdf — primary research — four key metrics; "2.6× more likely low performers" with external approval.
[^3]: ITIL 4 Foundation Glossary (PeopleCert/Axelos) — https://www.itconcepts.ch/wp-content/uploads/itil4-foundation-glossary-january-2019.pdf — standard body — definitions of change, standard/emergency change, RFC.
[^4]: IT Change Management: ITIL Framework (Atlassian) — https://www.atlassian.com/itsm/change-management — vendor (neutral content) — standard/normal/emergency change; CAB definition.
[^6]: Access controls & separation of duties best practices (Liquibase) — https://docs.liquibase.com/solution-guides/access-controls-separation-of-duties/access-controls-separation-of-duties-best-practices.html — vendor — dev/test/prod separation; writer ≠ approver ≠ deployer; YBIYRI friction.
[^7]: DORA | Streamlining change approval — https://dora.dev/capabilities/streamlining-change-approval/ — primary research — CAB negative impact; peer-review substitute; strategic CAB role.
[^8]: Accelerate (Forsgren/Humble/Kim) — https://itrevolution.com/product/accelerate/ — primary book — "worse than no change approval process"; peer review to meet SoD goal.
[^9]: Auditor's Guide to CI/CD Security (regulated-devsecops.com) — https://regulated-devsecops.com/start-here/ — specialist — "pipelines enforce control" auditor stance; EU-DORA/NIS2 mapping; maturity model.
[^10]: SOC 2 Trust Service Criteria Mapped to Pipeline Controls — https://regulated-devsecops.com/ci-cd-governance/soc-2-trust-service-criteria-mapped-to-pipeline-controls/ — specialist — CC8.1 change mgmt; system-generated evidence preference.
[^11]: Harness Deployment Freezes / Blackout Windows — https://www.harness.io/blog/harness-deployment-freezes — vendor — blackout-window definition; finance/holiday freeze examples.
[^12]: Release Freeze Strategy & Enforcement — https://beefed.ai/en/release-freeze-strategy-enforcement — practitioner — freeze triggers (regulatory cutoffs, financial close); GitLab `$CI_DEPLOY_FREEZE`.
[^13]: Create a Change Blackout Window (ServiceNow / The Snowball) — https://thesnowball.co/how-to/create-change-blackout-window — practitioner — ITSM blackout config; recurring quarter-end freezes; emergency carve-out.
[^16]: Separation of duties as separate teams (anti-pattern) — https://migration.minimumcd.org/docs/anti-patterns/separation-of-duties-antipattern/ — engineering (disconfirming) — separate-team SoD = weaker assurance + bottleneck; "a script deploys, reviewed independently."
[^17]: SOX IT Controls: Section 404 (UINAT) — https://uinat.com/compliance/sox-it-controls-guide/ — practitioner — IT-specific SoD list; "no single employee can unilaterally deploy"; trace-to-ticket.
[^18]: Database Change Management for SOX Compliance (Bytebase) — https://www.bytebase.com/blog/database-change-management-sox-compliance/ — vendor — four ITGC domains; auditor questions; trace-to-ticket.
[^19]: Maker-checker (Wikipedia) — https://en.wikipedia.org/wiki/Maker-checker — tertiary — maker-checker/4-eyes definition.
[^20]: Four-Eyes vs Maker-Checker vs SoD (Latch) — https://latchworkflow.com/blog/four-eyes-principle-maker-checker-segregation-of-duties/ — trade — distinguishes the three; "SoD is the design foundation."
[^21]: Four-Eyes Principle (Chequedb) — https://chequedb.com/resources/blog/four-eyes-principle-foundations — trade — SoD vs four-eyes vs dual control; audit-evidence point.
[^22]: Change Advisory Boards Don't Work (Octopus) — https://octopus.com/blog/change-advisory-boards-dont-work — vendor — corroborates Accelerate's anti-CAB finding.
[^23]: SEC 17a-4 / financial audit trails (Adaptive Query) — https://www.adaptivequerysystems.com/cryptographic-governance/financial-audit-trails — specialist — SEC 17a-4 WORM + 2022 audit-trail alternative.
[^24]: FFIEC / regulator system-development controls — https://ithandbook.ffiec.gov/it-booklets/development-and-acquisition/ — regulator — documented request→test→prod approvals "to ensure segregation of duties"; emergency-change controls.
[^25]: IaC Patterns for Regulated Trading Systems (behind.cloud) — https://behind.cloud/infrastructure-as-code-patterns-for-regulated-trading-system — specialist — author/reviewer/approver/deployer/auditor roles; break-glass logged+time-limited+reviewed.
[^28]: SOX ITGCs / system-dependent controls (Schneider Downs) — https://schneiderdowns.com/our-thoughts-on/sox-it-general-controls-system-dependent-controls/ — CPA — ITGC families; SoD/change/access = highest-risk; auditor reduces reliance if weak.
[^29]: Auditing IT Risk: Change Mgmt & AppDev (ISACA Journal) — https://www.isaca.org/resources/isaca-journal/past-issues/2011/auditing-it-risk-associated-with-change-management-and-application-development — ISACA — COBIT change-control objectives, RACI.
[^31]: Implementing SoD (ISACA Journal) — https://www.isaca.org/resources/isaca-journal/issues/2016/volume-3/implementing-segregation-of-duties-a-practical-experience-based-on-best-practices — ISACA — COBIT DSS06.03/EDM04.02; compensating controls when SoD infeasible.
[^32]: Compensating controls / SoD as a dial (ISACA) — https://www.isaca.org/resources/isaca-journal/issues/2016/volume-3/implementing-segregation-of-duties-a-practical-experience-based-on-best-practices — ISACA — compensating controls widely accepted where SoD infeasible.
[^33]: What Is Code Freeze? (Qodex) — https://qodex.ai/blog/what-is-code-freeze-in-software-development — trade — banking "complete lockdown during end-of-year financial processing."
[^34]: Mastercard Holiday Code Freeze Best Practices — https://www.mastercard.com/global/en/business/payments/holiday-code-freeze.html — payments network (primary) — payments-firm rationale for holiday freeze.
[^36]: PAM Strategy for Banking (Hammer IT) — https://hammeritconsulting.com/privileged-access-management-banking/ — trade — FFIEC/GLBA/PCI mandates; JIT; eliminate standing access.
[^37]: Privileged Account Management for Financial Services (NIST SP 1800-18) — https://www.nccoe.nist.gov/financial-services/privileged-account-management — NIST/NCCoE — PAM definition; privileged/emergency/service accounts.
[^38]: What Is Just-In-Time Access (Palo Alto Networks) — https://www.paloaltonetworks.com/cyberpedia/what-is-just-in-time-access-jit — vendor — JIT; break-glass production access; zero standing privilege.
[^40]: SOX Compliance for Software Delivery (Harness) — https://www.harness.io/harness-devops-academy/sox-compliance-for-software-delivery-explained — vendor — "the deployer can be the automated pipeline"; manual→systematic; policy-as-code.
[^42]: Policy as Code DevSecOps for Regulated Teams — https://makitsol.com/policy-as-code-devsecops-for-regulated-teams/ — trade — OPA/Rego + Kyverno; rules out of PDFs into PRs/CI.
[^43]: Powering Capital One's Microservices with CI/CD (Capital One Tech) — https://www.capitalone.com/tech/software-engineering/realigning-devops-practices-to-support-microservices/ — bank primary — SoD → immutable stage gates; auto evidence collection.
[^44]: How DBS is driving IT automation (Computer Weekly) — https://www.computerweekly.com/news/365535432/CW-Innovation-Awards-How-Singapores-DBS-Bank-is-driving-IT-automation — press — policy-as-code RBAC guardrails for regulatory compliance.
[^47]: Why Change Moratoriums Don't Work (Cloud Artisan) — https://cloudartisan.com/posts/2025-08-04-why-change-moratoriums-dont-work/ — practitioner (disconfirming) — freeze → pile-up → post-freeze disasters; blocks security patches.
[^48]: The Problem with Maintenance Windows and Change Freezes (Mangoteque) — https://blog.mangoteque.com/2024/02/28/change-freezes-and-maintenance-windows/ — practitioner (disconfirming) — freezes = large-batch anti-pattern.
[^49]: The Real Cost of Holiday Feature Freeze (CroCoder) — https://www.crocoder.dev/blog/real-cost-of-holiday-feature-freeze — practitioner (disconfirming) — measured Big-Bang effect; bug spikes.
