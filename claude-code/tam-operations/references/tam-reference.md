<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `tam-reference` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: tam-reference
description: >
  MongoDB Premium Services TAM Operating Reference — role taxonomy (TAM/NTSE/IR/DCE), incident
  management (S1–S4 SLAs, JIMP, Break-Glass, PIR/RCA), account lifecycle (onboarding phases,
  offboarding), strategic initiatives (Straight-to-8 upgrade, multi-region DR, sharding), risk
  management (Risk Register, GWP), templates (Engagement Overview, TSP, Support Plan), tooling
  (Monday.com, FTDC, JIRA HELP/PROACTIVE), KPIs, and GTM/PLG model.
  TRIGGER: questions about TAM/NTSE roles or responsibilities, incident severity SLAs, JIMP
  structure, Break-Glass protocols, account review cadence (EBR/QBR/MOR), Straight-to-8 upgrade
  procedure, onboarding/offboarding steps, Monday.com initiative tracking, RCA structure,
  engagement maturity model, Premium Services operating procedures, or any MongoDB TAM workflow.
  SKIP: generating TAM documents for a named account → tam-operations (references/tam-account-reports.md),
  active case management and TS Tools API → tam-operations (references/case-tracker.md), Atlas cluster
  diagnostics (use `atlas-diagnostics-expert`), code review or implementation tasks (use domain-specific
  skills).
version: 1.2.0
updated: "2026-06-01"
origin: local
category: mongodb
tags:
  - tam
  - ntse
  - premium-services
  - incident-management
  - jimp
  - account-lifecycle
  - rca
  - monday-com
  - kpis
  - ebr-qbr
triggers:
  - TAM role
  - NTSE role
  - incident severity SLA
  - JIMP
  - Break-Glass
  - account onboarding
  - account offboarding
  - Straight-to-8
  - RCA structure
  - EBR QBR MOR
  - engagement overview template
  - support plan template
  - Monday.com initiative
  - risk register
  - Premium Services
related_skills:
  - tam-account-reports
  - case-tracker
  - atlas-diagnostics-expert
  - operator-report-generator
---

# MongoDB Premium Services TAM Operating Reference

Generated from `docs/tam-context.md` in `10gen/mdb-tam`. Start from the bundled context below; defer to that source file for exact terminology, workflows, and operating constraints.

## When to use this skill

Use when the request needs help with MongoDB TAM/NTSE/Premium Services operating guidance: roles, incident processes, account lifecycle, strategic initiatives, risk management, templates, tooling, or KPIs.

**Do not use** for:
- Generating TAM documents for a named account → tam-operations (references/tam-account-reports.md)
- Active case management and TS Tools API → tam-operations (references/case-tracker.md)
- Atlas cluster diagnostics → use `atlas-diagnostics-expert`
- Code review or implementation → use domain-specific skills

## Common TAM Mistakes

- Mixing Diataxis types: a runbook (how-to) that drifts into explanation, or a reference doc that tries to teach via tutorial.
- Using RYG health scores without operationalizing them — every score state needs a mandatory action.
- Presenting benchmarks without prescriptive recommendations.
- Score inflation: defaulting accounts to Green without data-backed justification.
- Writing for one audience when the deliverable serves multiple (engineers + executives). Use layered structure: exec summary → findings → technical appendix.

## Cross-Skill Routing

| Need | Route |
|---|---|
| Active case management and TS Tools API | tam-operations (references/case-tracker.md) |
| Atlas cluster diagnostics and troubleshooting | `atlas-diagnostics-expert` |
| Generate account review/support plan/JIMP | tam-operations (references/tam-account-reports.md) |
| Code review, frontend design, implementation | Domain-specific skill |

---

## Bundled Reference

Source: `docs/tam-context.md`

---
title: MongoDB Premium Services TAM Operating Reference
source_documents:
  - TAM-glean_context.md
  - TAM-gemini-context.md
synthesis_version: 1
synthesis_date: 2026-05-12
dedupe_policy: prefer_more_detail_and_data
audience: LLM retrieval and TAM/NTSE working reference
---

## 0. Scope

Premium Services is the layer of MongoDB support that sits on top of standard support. It combines advanced reactive support with proactive risk management, account-level technical leadership, and structured initiative execution. This document is the canonical operating reference for TAM and NTSE work on enterprise accounts.

## 1. Role Taxonomy

### 1.1 TAM (Technical Account Manager)
- Account-level technical DRI (Directly Responsible Individual). Acts as the "quarterback" coordinating all internal teams.
- Operates at Tier 2 (Strategic Implementation) and Tier 3 (Architecture & Governance).
- Owns: account strategy, consumption-inhibitor identification, complex initiative coordination, Risk Register, architectural recommendations ahead of scaling events, alignment of technical milestones with business outcomes, JIMP governance.

### 1.2 NTSE (Named Technical Services Engineer)
- Note: The acronym is sometimes expanded as "Named Technical Support Engineer." MongoDB's formal expansion is "Named Technical Services Engineer."
- Tactical subject-matter expert with deep customer-specific support ownership. Operates at Tier 1 (Core Maintenance) and Tier 2.
- Owns: high-context case management, proactive assessments, baseline reviews, weekly case trend presentations, RCA investigations, support case backlog, hands-on diagnostic troubleshooting, pre-launch architectural context acquisition.

### 1.3 IR (Incident Responder)
- First-line S1 technical response for Enhanced Support accounts. Leads the war room for high-severity events.
- Must complete Tier-0 IR certification before active duty.

### 1.4 DCE / RCE (Dedicated / Regional Consulting Engineer)
- High-touch engineering resources deployed to execute the TAM's strategic roadmap.
- Owns: hands-on design, systems integration, production implementation, validation against best practices, deep-technical remediation that requires direct infrastructure modification.

## 2. Service Taxonomy

Premium Services classifies work into four categories crossed with three tiers.

Categories: Software Development & Engineering; Infrastructure & Cloud Operations; Data Management & Analytics; Cybersecurity & Compliance.

Tiers:
- Tier 1 – Core Maintenance. Routine backups, hardware replacements. NTSE operates here.
- Tier 2 – Strategic Implementation. Cloud migrations, database optimizations. TAM and NTSE both operate here.
- Tier 3 – Architecture & Governance. Long-term roadmapping, regulatory policy creation. TAM operates here.

## 3. Incident Management

### 3.1 Severity Tiers, SLAs, and Cadence

| Tier | Definition | Response SLA | Update Cadence | Operational Protocol |
| --- | --- | --- | --- | --- |
| S1 / P0 (Critical) | Complete outage or critical production impact. No workaround. Large user base or critical financial functions affected. | < 15 min, 24x7 | 30–60 min | Break-Glass authorized. Mandatory 15-min bridge mobilization. |
| S2 / P1 (High) | Significant business impact. Service operational but severely degraded. Workaround may exist but is unsustainable. | < 30–60 min, 24x7 | 2 hr or at major milestones | Intensive diagnostic gathering: FTDC, log retrieval prior to remediation. |
| S3 / P2 (Normal) | Non-critical software bugs, partial loss of non-essential services, staging/dev environment issues. | Standard business hours / 24 hr | Dynamic | Architectural guidance, performance tuning, routine case tracking. |
| S4 / P3 (Low) | Minor operational bugs, cosmetic UI, pre-scheduled proactive events (e.g., holiday scaling). | Standard business hours | Dynamic | Logged for future sprints or queued for advisory review. |

### 3.2 Error Budgets and Error Envelopes

- 99.99% availability SLA = 52.56 min unmitigated downtime per calendar year.
- MTTR (Mean Time to Resolution) target: 18 min per incident.
- Mathematical ceiling: ~2.92 SLA-breaching incidents per year before contract penalties apply.
- Error Envelope: explicit pass/fail metric defined for planned high-risk operations (e.g., Goldman Sachs Rewards/Cookie regional failovers East-to-West). Defines max cumulative read/write unavailability, NotPrimary / PrimarySteppedDown thresholds, and P95/P99 latency degradation boundaries vs. baseline. Exceeding the envelope triggers mandatory rollback.

### 3.3 Break-Glass Protocols

Used during S1 for Tier-0 accounts (e.g., Okta, Goldman Sachs). RTO takes priority over diagnostic capture. The JIMP pre-authorizes these.

- **Restart API.** Restricted private "Kill Switch" endpoint. Use cases: unresponsive gray failures, WiredTiger cache eviction deadlocks, thread starvation. Authorized customer leads accept calculated in-flight data loss to force-restart mongod or mongos and immediately restore availability. Bypasses standard Atlas UI guardrails.
- **VM Replacement Protocol.** Bypasses cloud provider guardrails such as the AWS EBS 6-hour modification cooldown. Triggered when a cluster has correlated high IOWAIT and severe degradation due to an EBS volume locked in "optimizing" after a disk modification. TAM raises a P1 HELP ticket to Cloud Operations Engineering (COE). COE force-terminates the underlying EC2 VM and provisions a new host from a fresh snapshot. Restores service in 15–20 min.
- **Tier-0 Identity Isolation Protocol.** Triggered when identity infrastructure is compromised (Okta Super Admins, IdPs, Root Keys). Actions: immediate freeze of all admin changes in the cell, strict dual-custodian verification for any remedial action, automated revocation of OIDC tokens, programmatic pivot of Okta Sign-on Policies to require Phishing-Resistant MFA (WebAuthn / FIDO2) exclusively, immediate rotation of Okta Customer Support (OCS) tokens, move all coordination to pre-verified out-of-band channels.

### 3.4 JIMP (Joint Incident Management Plan)

Authoritative runbook used during live high-severity events to align MongoDB and customer response. Predefined Severity Mapping, Escalation Matrices, and Rules of Engagement.

JIMP roles:
- **Incident Commander (IC)** — overall strategy and delegation.
- **Operations Lead** — directs deep-technical diagnostic and mitigation execution. Pre-authorized to invoke Break-Glass.
- **Communications Lead (CL)** — manages internal and external stakeholder updates, insulating the technical team.
- **Scribe** — real-time chronological documentation of events.

### 3.5 Shift Handoff — "The Bridge"

Live synchronous meeting or dedicated tactical chat thread used during prolonged P0/P1 incidents that span multiple regions. The outgoing IC delivers a structured "State of Play": current business impact, technical actions executed this shift, residual risks, pending tasks. The incoming IC must explicitly accept the baton before the outgoing responder is relieved.

### 3.6 PIR / RCA (Post-Incident Review / Root Cause Analysis)

Initiated by the TAM after every S1 or S2 resolution.

- Preliminary draft within 24 hr of resolution.
- Formal retrospective meeting within 3–5 business days.
- Customer-facing incident report due within 5–7 days of resolution.

RCA structure:
1. Executive Summary.
2. Environment Details — cluster version, topology, scale.
3. Chronological Timeline in UTC, verified.
4. "5 Whys" Root Cause Analysis.
5. Corrective and Preventive Action Items with named owners and strict due dates.

## 4. Communication Protocols

- **Cases First, Slack Second.** The Support portal (Customer Hub / Support Hub) is the immutable system of record for substantive technical directives. Slack is for high-velocity coordination.
- **15-Minute Duplication Rule.** If a Slack inquiry requires more than 15 minutes of investigation or represents a recurring architectural query, the TAM must duplicate the directive into an official support case.
- OOO and offboarding require Monday.com updates, case handovers, calendar blocks, and Slack notices.

## 5. Account Lifecycle

### 5.1 Onboarding — Phase 1: Knowledge Transfer (Weeks 1–2)
- Audit account history.
- Review current cluster architecture.
- Examine active support case backlog.
- Evaluate Historical Issue Review docs for recurring debt patterns over the prior 4–6 months.
- Close with internal alignment sessions with AE and SA to harmonize technical strategy with commercial objectives.

### 5.2 Onboarding — Phase 2: Customer Introduction (Weeks 3–4)
- Formal introduction to customer executive sponsors and lead engineers.
- Establish synchronous cadences (weekly or bi-weekly syncs).
- Define boundaries and proactive expectations of the Premium Services partnership per the TAM Service Description.

### 5.3 Onboarding — Phase 3: Operational Alignment (Month 2+)
- TAM assumes full operational command.
- Prepare and deliver the first EBR.
- Take full ownership of Monday.com Initiative Trackers.
- Actively manage the JIMP.
- Continuous proactive risk identification (EOL agent deprecations, missing MFA, etc.).

### 5.4 Execute Phase (Months 2–6)
- Execute initiatives with Tier-0 service focus.
- Continuously update Support Plans, initiative docs.
- Feed case learnings into product requests.

### 5.5 Realign / Renew Phase (Months 6–12)
- Compare progress against goals.
- Reassess risk posture.
- Use QBRs to reset initiatives for the next cycle.

### 5.6 Offboarding

Required when a TAM rotates off an account:
- Schedule formal sync with the incoming TAM to transfer nuance not captured in documentation.
- Close all open workstreams in Monday.com.
- Update the Account Health Scorecard (Green / Yellow / Red).
- Conduct KT sessions on active projects and long-term architectural strategy.
- Send an Executive Summary Email to customer stakeholders summarizing specific technical achievements, cost optimizations, and systemic stabilizations delivered during the TAM's tenure.
- Author Success Story documentation (Section 9.4).

## 6. Strategic Initiatives

### 6.1 Definition
Long-term structured projects designed to solve underlying technical challenges and accelerate customer objectives. Tracked rigorously in Monday.com.

### 6.2 Straight-to-8 Upgrade (v4.4 → v8.0)
Non-standard major-version upgrade that bypasses 5.0, 6.0, 7.0. Used for large Tier-1 workloads (example: Goldman Sachs Marcus "Cookie" and SEM clusters).

Execution sequence:
1. Orderly shutdown of all application instances running S7 (legacy) binaries.
2. Stage new S8 binary package to `/opt/app/bin/s8_stage`.
3. SHA-256 checksum verifying integrity.
4. Update system symlinks to new S8 path.
5. Run `s8_init.sh` in "check-only" mode to validate environment.

FCV Pinning safety mechanism:
- Core binaries upgraded to 8.0 while Feature Compatibility Version stays pinned at 7.0.
- Allows rapid binary downgrade to 4.4 / 7.0 without a full logical restore.
- "Point of No Return" is the explicit command to unpin the FCV to 8.0.

Guardrails:
- Block DDL operations via automated scripts during the window.
- WiredTiger tickets temporarily hardcoded to 128 to suppress caching regressions.
- Driver alignment is mandatory: certify the Java 4.10 driver on JRE 11+. Expunge legacy 3.x drivers from the classpath to prevent NoSuchMethodError during binary init.
- Aggressively tune `maxConnectionsPerHost` to match the 4.10 driver's threading model.

Known issue (Morgan Stanley SEM2 deployment): new shard nodes triggered `ECONNREFUSED` on custom ports 8872 and 8874 when contacting the Config Server Replica Set (CSRS). Requires coordinated intervention from the customer's network engineering team before the deployment proceeds.

Common upgrade error codes: `ERR_DB_LOCK_CONTENTION_80`, `AUTH_TOKEN_PROPAGATION_DELAY`.

### 6.3 Multi-Region Topology Migration / Cloud Outage DR

Triggered by catastrophic regional collapses (example: AWS me-central-1 / me-south-1 AZ outages impacting Darwinbox, DTC Ehail Mobility, Moody's Analytics).

Pattern: reconfigure a stalled 3-node cluster in me-central-1 to a 2-1-2 distribution across me-central-1, eu-central-1 (Frankfurt), and ap-south-1 (Mumbai).

Operational notes:
- Cloud provider APIs frequently exhibit instability during these moves: snapshot throttling, persistent `AWSEnsureNetworkPermissionsAppliedMove` Atlas plan failures.
- To bypass AWS snapshot throttling, COE and TAMs may abandon automated restoration in favor of manual `rsync` or File Copy Based Initial Sync (FCBIS).
- Cross-region moves fracture PrivateLink and VPC peering. Pre-provision new private endpoints in the target region before node addition.
- Post-migration verification: run `dig +trace` from peered VPCs to confirm public DNS hostnames resolve to new private IPs, ensuring traffic is not routed over the public internet.

### 6.4 Sharding, Rebalancing, High-Volume Initial Sync

Used to reclaim large amounts of fragmented disk space or resolve orphan accumulation on donor shards.

Triage signal: balancer throughput collapsing from ~40 GB/day to ~3 GB/day. Check config server logs for "No available shards to take chunks for zone" errors.

Mitigation: temporarily disable balancing on non-essential collections to prioritize core asset migration.

Initial Sync sequencing:
- Logical Initial Sync (LIS) on the primary donor node.
- Optimized Initial Sync (OIS) on remaining secondaries.

OOM mitigation (example: Wise Plc mongosync migration): COE tunes memory parameters:
- `maxIndexBuildMemoryUsageMegabytes = 10GB`
- `minSnapshotHistoryWindowInSeconds = 1`
- Enable `tcmallocAggressiveMemoryDecommit` to aggressively flush fragmented memory pages back to the OS.

### 6.5 AI Agent Golden Path Program

Pre-validated architectural route for customers building AI agents. Leverages approved internal SDKs and platforms like Vertex AI. NTSEs maintain the infrastructure; standardized telemetry and logging allow rapid debugging without deciphering bespoke environments. Field failures feed back into the automation pipeline.

## 7. Risk Management

- **Risk Register.** Living document of technical and process risks. Risks sourced from health checks, incident retrospectives, and consulting reports, then mapped to specific mitigations and initiatives.
- **Risk-Driven Method.** Named methodology that places the Risk Register at the center of TAM operations.
- **Get-Well Plan (GWP).** Targeted stabilization plan used to transition at-risk accounts from low-trust post-incident states to stable well-run environments.
- Proactive risk monitoring covers EOL agent deprecations, missing MFA, scaling vulnerabilities.

## 8. Workload Pre-Onboarding Questionnaire

Administered before deployment of any net-new workload or major architectural migration.

Technical requirements:
- Deployment type: Atlas Managed, Self-Managed, or Hybrid.
- Cloud Provider: AWS, Azure, or GCP.
- Geographical data residency regions.
- Initial Cluster Tier sizing based on RAM/CPU (M10, M30, M80).
- Multi-region fault tolerance and replica-set configuration (minimum 3 nodes).

Data profiling:
- Workload type: Transactional OLTP vs. Analytical OLAP.
- Initial data volume in TB; 12-month projected growth curve.
- Peak throughput in Ops/sec with explicit read/write ratios.
- Average document size in KB.
- Primary access patterns: exact match, text search, range scans. Used to recommend indexing strategy.

Security:
- IP Access Lists.
- VPC Peering architecture.
- PrivateLink requirements.
- Regulated data classification: PII, PCI, HIPAA.
- Customer Key Management.

## 9. Templates and Documents

### 9.1 TAM Engagement Overview
Single source of truth for the entire account team. Required sections:
- TL;DR and Timeline. Executive summary of the customer's engineering orientation; continuously updated chronological timeline of critical events, contract modifications, major architectural shifts.
- Current Status Scorecard. Traffic-light indicators (Positive, Mostly Positive, Some Caution, Caution, Immediate Need) mapped to the Five Foci.
- Team Index. Customer Technical Leads mapped to MongoDB TAM, NTSE, DCE, SA, CSM.
- Architecture & Incident Routing. Links to Atlas production clusters, QBR schedule, S1/S2 War Room procedures vs. S3/S4 case routing.

Five Foci (note: gemini source uses slightly inconsistent labels across sections — "Initiatives, Issues, Product, GTM, Operations" in one place and "Issues, Initiatives, Product, Commercial, Operations" in another. Treat GTM and Commercial as the same focus).

### 9.2 Technical Success Plan (TSP) / TAM Engagement Roadmap
Dynamic visualization deployed during Monthly Technical Reviews. Deconstructs multi-year objectives into phases. Each priority lists:
- Scope.
- Hyperlinked Monday.com Initiative IDs.
- Customer Next Steps and MongoDB Next Steps (explicitly assigned).
- Timeline projections.
- Priority ranking: High / Medium / Low.

### 9.3 Support Plan
Master scaffold that maps customer business outcomes to specific Premium Services initiatives with measurable success criteria.

### 9.4 Success Story Documentation

Required at the close of a major initiative or before TAM rotation. Narrative structure:
1. Customer Overview and Business Goals.
2. Start-of-Engagement state — operational pain points, architectural blockers.
3. The Solution — what the TAM/NTSE executed.
4. Measurable Outcomes — explicit and quantifiable. Examples: 40% reduction in SEV 1/2 incidents, exact infrastructure cost savings from node rightsizing, latency reductions in ms, workload-onboarding expansion metrics.

## 10. Runbooks

### 10.1 T-7 to T-1 Pre-Flight Checks (before major maintenance)
1. Cluster and Replication Integrity. Atlas UI shows HEALTHY. Run `db.adminCommand({ replSetGetStatus: 1 }).initialSyncStatus`. Replication lag stable. Oplog window sized to survive the maintenance duration. No background FCBIS or schema migrations consuming disk IOPS.
2. Backup Verification. Recent successful snapshots present. If Continuous Cloud Backups are active, verify integrity of latest restore points. Manually trigger one final out-of-band safety snapshot before maintenance start.
3. Feature Compatibility Version Audit. Cluster FCV matches required baseline (example: FCV 7.0). Rollback safety net is viable.
4. Application Freeze Protocols. Written confirmation from the customer that no DDL operations, driver deployments, or heavy non-essential batch-analytics jobs will run concurrently with the upgrade.
5. Infrastructure Prerequisites (T-1). Hosts on supported OS and kernel (example: RHEL 8; avoid problematic Linux Kernel 6.19). Audit and security agents healthy. Disk and IOPS utilization with minimum 20% free headroom.

### 10.2 Change Stream / Oplog Diagnostic (errors: ChangeStreamHistoryLost, "Resume of change stream was not possible")
1. Extract diagnostic logs for the exact requested resume-token timestamp generated by the consumer application.
2. Compare timestamp against the shard's `earliestOptime` and `latestOptime` boundaries.
3. If requested resume timestamp is older than `earliestOptime`, oplog history is unrecoverable server-side; the global stream fails entirely.
4. Examine connection metrics. Confirm cross-region consumer watchers use `readPreference=secondaryPreferred`. Assess for Mutex contention during subsequent CPU peaks.
5. Instruct the customer to forcibly restart failing change-stream consumers with a new `startAtOperationTime` set to the earliest surviving oplog entry.

### 10.3 Ops Manager / Kubernetes Backup Failures (cascading 503 errors)
1. Structure investigation linearly: Global Load Balancer → Ops Manager services → Kubernetes Backup Agents.
2. Collect `mms0.log` and `backup-daemon.log` from all Ops Manager nodes behind the LB. Pull Automation and Backup Agent logs from the target replica set.
3. Interrogate K8s worker node logs for OOMKill or CPU throttling. Confirm JVM `Xmx` is accurately tuned to deployed pod memory specs.
4. Validate with the customer's network team whether the Global LB is modifying or rewriting HTTP responses between Backup Agents and the Ops Manager control plane.

## 11. Tooling

### 11.1 Monday.com (canonical source of truth)
- Parent Items = Initiatives. Sub-Items = Milestones.
- Required free-form-text update syntax:
  - Activity Date
  - Updates (synthesized high-density key points)
  - Next Steps (explicitly assigned owners)
  - Time Spent (optional hourly logging)
- Weekly: TAM finalizes status and publishes a sentiment indicator (Positive, Mostly Positive, Some Caution, Caution, Immediate Need) across the Five Foci.
- External (customer-facing) boards: standardized name prefix per template convention. Sensitive internal data segregated into a hidden "MDB Files" column restricted from external guests. Dedicated "Shared Files" column for bilateral document transfer. Mirror columns (Product Status, Product Manager, etc.) dynamically pull data from internal boards.

### 11.2 Diagnostic Telemetry
- **FTDC (Full Time Diagnostic Data Capture).** Primary artifact for hardware saturation and execution bottlenecks. NTSEs query FTDC for memory fragmentation thresholds, tcmalloc current allocated bytes during massive index builds, CPU steal percentages, disk IOWAIT saturation, getmore cursor operation efficiency.
- **Splunk Logs / Lumberjack.** Real-time infrastructure event tracking.
- **ts-diag, org-report, Atlas CLI.** Standard local diagnostic and reporting tools.

### 11.3 Jira Ecosystem
- **HELP tickets.** Trigger direct intervention from Cloud Operations Engineering (COE). Examples: bypass stalled autoscaling plans, adjust underlying EBS volumes, force VM replacement.
- **PROACTIVE tickets.** Track preemptive anomaly detection that flags clusters before total failure.

### 11.4 Atlas Control Plane
Used by IRs during Break-Glass. Access must be verified before assuming active duty.

## 12. KPIs and Success Metrics

### 12.1 Response and MTTR
- S1 response target: 15 min.
- S2 response target: 30–60 min.
- MTTR optimization floor: 18 min for critical service restorations.

### 12.2 Satisfaction and Execution
- CSAT / NPS for managed interactions: > 4.5 / 5.0.
- Customer-facing RCA delivery window: 5–7 days post-resolution.

### 12.3 Adoption and Velocity
- Feature Adoption Rate: 60% of licensed users engaging with advanced capabilities within 90 days.
- Time to Value (TTV) for the customer's first key action: < 48 hr.
- DAU / MAU ratio: > 40%.
- Expansion Velocity: Seed → Team tier in < 6 months.

### 12.4 Operational Optimization
- Documented dollar savings from targeted infrastructure rightsizing.
- Percentage improvement against CIS security benchmarks.

## 13. Engagement Maturity Model

| Level | Operational Characteristics | TAM Focus |
| --- | --- | --- |
| 1 – Reactive (Initial) | Incident-driven. High volume of low-complexity tickets. Continuous firefighting. | Basic stabilization. Aggressive case management. Establish initial communication cadences. |
| 2 – Managed (Standard) | Predictable cadences. Documented Runbooks and SOPs in use. | Identify recurring patterns. Trend analysis. Address low-hanging technical debt. |
| 3 – Proactive (Advanced) | Predictive service alerts. Focus shifted from recovery to capacity planning and upgrade readiness. | Pre-emptive Health Checks. Risk Assessments. Architecture reviews before scaling. |
| 4 – Strategic (Optimized) | TCO and ROI optimized. Customer participates in product feedback loops. | Trusted advisor embedded in the long-term roadmap. Drive continuous operational excellence. |

## 14. Account Reviews

Cadence:
- Weekly: case review sheets, NTSE status calls, proactive assessments.
- Monthly: Monthly Operational Review (MOR), Monthly Technical Review (uses TSP).
- Quarterly: Quarterly Business Review (QBR), value confirmation, initiative reviews.
- Per major engagement: Executive Business Review (EBR).

EBR / MOR deck flow (four phases, in order):
1. Executive Overview. Bottom Line Up Front. Overall account health and primary focus areas since last engagement.
2. Operational Review. Ticket volumes, priority distributions (P1 / P2 / P3 trends), MTTR metrics, SLA compliance statistics.
3. Strategic Alignment. Technical Roadmap, execution status of ongoing initiatives, technical-debt retirement, Account Health Scorecard mapping relationship health against business value realization.
4. Forward Planning. Next steps, upcoming product beta participation, ownership for upcoming milestones.

## 15. AI Operational Constraints (Tier-0 environments)

Allowed:
- Authorized AI usage such as the Okta Copilot Assistant for data synthesis and knowledge retrieval.
- Ingesting high-volume AWS CloudTrail or Okta System Logs to isolate "Impossible Travel" or "Anomalous Admin Geography" patterns.
- Drafting initial "T0 Alert" executive notifications that translate complex telemetry into business-risk summaries.
- Acting as a rapid query engine for navigating extensive SOP documentation.

Prohibited:
- AI systems with autonomous write-action capabilities. Any CLI command, API payload, or configuration script generated by an AI assistant must be manually validated and executed by a human Incident Commander.
- Uploading Customer Support Information (CSI) — log files, architectural diagrams, proprietary code — into unapproved external LLMs. Violation grounds for immediate PIP escalation.

## 16. Professional Standards and Quality Bar

- Maintain active visibility in Slack channels. Acknowledge requests promptly.
- Do not share technical documentation without articulating specific context and intended next steps.
- Translate broad technical directives into structured execution plans. Maintain accurate metadata, milestones, and status indicators in standardized Monday.com Initiative Trackers.
- Do not deploy unauthorized custom-built internal tools or "shadow IT" dashboards. AI wrappers and uploads of CSI to external platforms are severe policy violations.
- Establish clear ownership during weekly account syncs. Lead discussions with a predefined structured priority list, not reactive external prompting.
- TAM and CSM/AE co-own the account strategy and drive renewals. TAM collaborates with Professional Services (PS) to identify and scope consulting opportunities.

## 17. GTM and Product-Led Growth (PLG)

- Methodology: "Land and Expand." Stabilize initial low-friction adoption, prove architectural viability, then advocate for seat expansion and onboarding of mission-critical workloads.
- Adjacent-TAM penetration target: 15% market penetration in targeted verticals.
- Economic target: 3:1 LTV:CAC during account expansion phases.
- TAM intervention is explicitly credited with influencing revenue growth and retention across the organization's top-tier accounts. Quantified via Success Stories (Section 9.4).

## 18. IR Readiness Checklist (for engineers rotating into the Tier-0 IR pool)

- Proficiency with the Okta IR Operational Checklist: shift-start verifications, live execution protocols, shift-end handoff documentation.
- Verified access to PagerDuty, incident tracking systems, and the Atlas control plane before assuming active duty.
- Authorization to execute Break-Glass operations without administrative friction.

## 19. Named Customer References

| Customer | Context |
| --- | --- |
| Okta | Tier-0; Break-Glass authorizations; Identity Isolation Protocol; Okta Copilot Assistant; Norberto Leite is named Customer Lead. |
| Goldman Sachs (Marcus / Rewards / Cookie / SEM) | Tier-0; Error Envelope failovers; Straight-to-8 target workload. |
| Morgan Stanley (SEM2) | ECONNREFUSED on custom ports 8872 / 8874 during Straight-to-8 deployment. |
| Darwinbox, DTC Ehail Mobility, Moody's Analytics | AWS me-central-1 / me-south-1 outage impact; multi-region topology migrations. |
| Zomato (`chat.message_store`) | Sharding / rebalancing example. |
| Cimpress Schweiz GmbH | Sharding / rebalancing example. |
| Wise Plc | mongosync migration; OOM mitigation via memory parameter tuning. |
| Trend Micro (`mgcp-prod`) | ChangeStreamHistoryLost diagnostic runbook example. |

## 20. Glossary

| Term | Definition |
|---|---|
| AE | Account Executive |
| COE | Cloud Operations Engineering |
| CSAT | Customer Satisfaction |
| CSI | Customer Support Information |
| CSM | Customer Success Manager |
| CSRS | Config Server Replica Set |
| DCE | Dedicated Consulting Engineer |
| DRI | Directly Responsible Individual |
| EBR | Executive Business Review |
| EBS | Elastic Block Store (AWS) |
| EOL | End of Life |
| FCBIS | File Copy Based Initial Sync |
| FCV | Feature Compatibility Version |
| FTDC | Full Time Diagnostic Data Capture |
| GWP | Get-Well Plan |
| IC | Incident Commander |
| IR | Incident Responder |
| JIMP | Joint Incident Management Plan |
| KT | Knowledge Transfer |
| LIS | Logical Initial Sync |
| LTV:CAC | Lifetime Value to Customer Acquisition Cost |
| MOR | Monthly Operational Review |
| MTTR | Mean Time to Resolution / Recovery |
| NPS | Net Promoter Score |
| NTSE | Named Technical Services Engineer |
| OCS | Okta Customer Support |
| OIDC | OpenID Connect |
| OIS | Optimized Initial Sync |
| OLAP | Online Analytical Processing |
| OLTP | Online Transaction Processing |
| OOM | Out of Memory |
| PIP | Performance Improvement Plan |
| PIR | Post-Incident Review |
| PLG | Product-Led Growth |
| PS | Professional Services |
| QBR | Quarterly Business Review |
| RCA | Root Cause Analysis |
| RCE | Regional Consulting Engineer |
| RTO | Recovery Time Objective |
| SA | Solutions Architect |
| SLA | Service Level Agreement |
| SLO | Service Level Objective |
| SOP | Standard Operating Procedure |
| TAM | Technical Account Manager |
| TCO | Total Cost of Ownership |
| TSP | Technical Success Plan |
| TTV | Time to Value |
| VPC | Virtual Private Cloud |
