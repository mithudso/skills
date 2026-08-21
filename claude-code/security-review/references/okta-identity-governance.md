<!-- hub-reference-banner -->
> **Reference file — part of the `security-review` hub.** Spoke created by /dr deep research
> (cluster C5, 2026-06-11). Sibling topics in this family are reference files under the hubs
> (`security-review`), **not** standalone skills. Load sibling topics from the owning hub's
> `references/<name>.md` (see the hub's routing table).

---

---
name: okta-identity-governance
version: 1.0.0
updated: 2026-06-11
description: >
  Okta Identity Governance (OIG): access certification campaign design
  (resource/user/entitlement campaigns, security access reviews, preconfigured
  campaigns), Access Requests (conditions vs request types, approval
  sequences, Slack/Teams), Entitlement Management (bundles, policies,
  connectors, disconnected apps), separation-of-duties rules and their
  evaluation caveats, governance reporting and the auditor reporting package
  for SOX/SOC 2/ISO evidence, packaging/licensing reality, lifecycle/Workflows
  integration, certification anti-patterns (rubber-stamping, scope explosion).
  TRIGGER: designing or auditing OIG certification campaigns; access request
  approval workflows in Okta; OIG SOD rules; access-review audit evidence
  from Okta; OIG vs SailPoint/Saviynt-class IGA; OIG licensing.
  SKIP: SCIM provisioning (okta-scim-lifecycle-security); tenant hardening
  (okta-admin-hardening); app integration and management-API dev
  (okta-expert); PAM depth; non-Okta IGA platform implementation.
category: security
tags: [okta, oig, identity-governance, iga, access-certification, access-requests, entitlement-management, separation-of-duties, sox, soc2, audit-evidence, rubber-stamping]
related_skills: [software-engineering-patterns, tam-operations]
whenToUse:
  - "Design or review an OIG access certification campaign"
  - "Build access request approval workflows (conditions, approval sequences)"
  - "Configure separation-of-duties rules for entitlements in Okta"
  - "Set up Entitlement Management (bundles, policies, connectors)"
  - "Produce SOX/SOC 2/ISO access-review evidence from OIG"
  - "Decide whether OIG suffices or a dedicated IGA platform is needed"
  - "Understand OIG licensing and packaging (what the add-on includes)"
  - "Diagnose rubber-stamping or scope explosion in certification campaigns"
whenNotToUse:
  - "SCIM provisioning mechanics and lifecycle security: use okta-scim-lifecycle-security"
  - "Admin console/tenant hardening, super-admin reduction: use okta-admin-hardening"
  - "Okta app integration, auth servers, hooks, API development: use okta-expert"
  - "Privileged access management depth (Okta Privileged Access): separate cluster"
origin: local
---

# Okta Identity Governance (OIG)

Expert reference for Okta's IGA layer: certification campaigns, access
requests, entitlement management, separation of duties, and the reporting
that turns those processes into audit evidence, plus an honest account of
where OIG ends and dedicated IGA platforms begin. Synthesized from Okta
first-party docs, Okta-practitioner writing, and independent and competitor
analyses. verified-as-of: 2026-06-11.

## When to use this reference

- Designing, reviewing, or troubleshooting OIG access certification campaigns
- Building self-service access requests with approval routing in Okta
- Configuring entitlements, bundles, resource collections, or SOD rules
- Assembling SOX / SOC 2 / ISO 27001 access-review evidence from OIG
- Evaluating OIG against SailPoint/Saviynt-class IGA; explaining OIG licensing

## When NOT to use this reference

- SCIM provisioning mechanics and deprovisioning security: `okta-scim-lifecycle-security` (this hub)
- Tenant/admin-console hardening, Govern Okta Admin Roles internals: `okta-admin-hardening` (this hub)
- Okta platform development (APIs, auth servers, Terraform): `okta-expert` (this hub)

---

## Packaging and licensing reality

OIG is a **subscription add-on to an existing Okta Workforce Identity
tenant**, not a standalone product. The add-on bundles three offerings:
**Lifecycle Management (LCM), Okta Workflows, and Access Governance**. Access
Governance (certifications, requests, entitlements, SOD) is a feature of OIG
and is *not* sold on its own.[^9][^1][^15] Activating it adds governance
surfaces to the existing Admin Console without disturbing SSO/MFA config.[^15]

Pricing facts (volatile; verified-as-of: 2026-06-11):[^9][^22]

| Item | Status |
| --- | --- |
| Workforce suites | Starter $6/user/mo; Core Essentials $14; Essentials $17; Professional and Enterprise quote-only. Billed annually, $1,500 annual contract minimum |
| Identity Governance in suites | Included from the Essentials suite up (alongside LCM and Privileged Access); an add-on elsewhere; absent from Starter/Core Essentials |
| Unit model | Per-user/per-month, suite-unified license counts (exceptions: Workflows flow counts, Okta Privileged Access resource units) |
| Not in the sticker price | Implementation services, premium support, advanced analytics, custom connectors[^22] |

Procurement implications: OIG presupposes the Okta platform ("an add-on to
the required Workforce Identity suite"[^17]), so its cost case is strongest
when Okta is already the IdP; and since the add-on includes LCM and
Workflows, teams paying for those separately should re-price as a bundle.
Independent reviewers rate OIG "expensive for what it delivers at the
mid-market tier"[^16] — negotiate against the suite, not the add-on list.

OIG became generally available in 2022; capability has accreted since (user
campaigns 2023-06, Entitlement Management GA 2024-01-22, preconfigured
campaigns release 2025.01.0).[^15][^12][^11][^13] Access Requests descends
from Okta's atSpoke acquisition rather than organic development, a lineage
critics cite when explaining seams between the Access Requests web app and
the rest of the Admin Console.[^16]

## Architecture: how the pieces fit

The three offerings divide the work: **LCM** handles joiner/mover/leaver
fulfillment (directory integration, profile mapping, OIN/SCIM provisioning,
On-Premises Provisioning), **Workflows** automates custom actions (including
ones triggered by OIG-specific System Log events), and **Access Governance**
supplies the governance surfaces proper:[^1]

- **Access Certifications**: recurring/one-off campaigns plus event-driven security access reviews, with automatic remediation
- **Access Requests**: self-service requests routed through approval sequences or request types
- **Entitlement Management**: app entitlements as first-class Okta objects, assignable by policy or individually
- **Separation of duties**: allow/block rules over entitlement combinations
- **Resource collections, owners, labels**: cross-app role-like groupings, ownership for approval/review routing, and labels (default `Crown Jewel`, `Privileged`) for campaign scoping[^1][^3]
- **Reports, APIs, System Log events** for evidence and extensibility[^1]

Okta positions the stack as least-privilege enforcement that doubles as audit
machinery: campaigns "help you meet your audit and compliance requirements or
professional standards like SOC2 and SOX."[^1] License reclamation from
certifying inactive users is a budget side benefit (vendor case-study
figures are directional only).[^10]

## Access certification campaigns: design

### Campaign types

| Type | Scope question | Use it for |
| --- | --- | --- |
| Resource campaign | Who has access to these resources? | App/group/entitlement owners; compliance evidence per system[^12] |
| User campaign (added 2023-06) | What access does this user have? | Manager-led User Access Reviews across a person's grants[^12] |
| Security access review | Is this user's access appropriate, now? | Event-driven review after an incident or anomaly; launch manually or via API/trigger[^2] |
| Preconfigured campaigns (release 2025.01.0) | Templates | "Discover inactive users" (top 5 apps by inactive users, 90-day scope) and "Okta administrator review" (admin-role entitlements on the Admin Console app)[^13] |

Resource campaigns certify groups, apps, entitlements, entitlement bundles,
collections, and (with Okta Privileged Access) service accounts; user
campaigns can include or exclude individually-assigned vs group/rule-derived
assignments and exclude named apps/groups.[^2][^12]

### Reviewer design (where campaigns succeed or fail)

- Reviewer types: manager, custom attribute, group owner, resource owner, or named users, with a **fallback reviewer** when resolution fails. The manager type requires `managerId` (or the referenced attribute) to carry the manager's Okta username/ID; unresolved managers silently route to the fallback, so validate manager-attribute hygiene before launch.[^3]
- Resource-owner routing falls back in chain: entitlement/bundle owner → app owner → fallback; collection owner → fallback.[^3]
- **Disable self-review** so users never certify their own access.[^3]
- **Multilevel review**: up to two levels per campaign; level 2 sees level 1's decision and justification, items reach level 2 only after level 1 acts, and remediation applies to the *final* reviewer's decision. Slow level-1 reviewers stall the campaign, so set notification cadence accordingly.[^3]
- **Delegates** keep campaigns moving when reviewers are out; the same mechanism covers access-request approvals.[^3][^7]
- Slack reviewer notifications require the Slack workspace integration.[^3]

### Decision support and remediation

Reviewers decide Approve / Revoke / Reassign per item.[^12] Okta remediates
automatically per campaign settings (remove the grant on revoke, or take no
action and require manual remediation; preconfigured campaigns default to
*no* remediation action, which you should usually change).[^2][^13] Reviewer
context panels show assignment date, last access, and usage;[^13]
**Governance Analyzer** (Early Access, verified-as-of: 2026-06-11) adds
insights (SOD conflicts, usage history, past decisions, profile changes) plus
approve/revoke recommendations, and **Smart Review** groups items into steps
by common attribute to cut repetitive decisions.[^2] These target
rubber-stamping directly; see Anti-patterns below.

### Campaign hygiene (Okta's own best practices)

Self-explanatory names and purpose-bearing descriptions (reviewers see both);
resources verified active; recurring-campaign considerations reviewed; the
**Create auditor reporting package** checkbox enabled on audit-driven
resource campaigns; entitlement reviews require the app's Governance Engine
enabled and entitlements created first.[^3]

## Access Requests and approval workflows

Two request models coexist; choosing wrong creates rework:[^6]

| Aspect | Conditions (newer; Okta says start here) | Request types (legacy atSpoke model) |
| --- | --- | --- |
| Definition lives | On the app's profile page in the Admin Console | In the Access Requests web app, owned by Access Requests teams |
| Requester UX | App tile in the End-User Dashboard catalog | Access Requests web app, Slack, or Teams (plus dashboard tiles with Unified requester experience) |
| Approval logic | Reusable **approval sequences**: questions, approval tasks, custom tasks, Workflows delegated-flow triggers | Per-request-type step flows: questions, tasks, timers |
| Scope control | Who can request, level, **duration (auto-revoke on expiry)**; reuses Okta/AD groups and entitlement bundles | Configuration lists binding request types to resources |
| External systems | Jira/ServiceNow automation **not available**[^7] | Jira and ServiceNow actions supported[^6] |
| Privacy | Requests forced private (admins + involved parties only)[^7] | Team-visible |
| Admin roles | Supports admin role bundles with time-bound grants (Govern Okta Admin Roles) | Cannot manage admin-role bundles or admin-granting groups[^6] |

Operational gotchas worth designing around: group profile changes used in
conditions take effect only after disabling/re-enabling the condition
(dashboard tiles otherwise refresh twice daily); approval steps assigned to a
group snapshot membership at request time, so approvers added later cannot
act; `managerId` must be populated for manager-scoped request-on-behalf
rules; requesters cannot reopen condition-managed requests. Escalation (to
the task assignee's manager) and delegates exist for stalled tasks.[^7]

Slack and Microsoft Teams integration lets users submit and approve from
chat,[^6] and entitlement bundles surface with end-user-friendly descriptions
so people request the level they need rather than guessing at group
names.[^11] SOD interacts with requests at submission time; see the race
condition below.

## Entitlement Management

Entitlement Management promotes app entitlements (permissions, roles,
licenses inside downstream apps) to **first-class Okta objects** alongside
users, groups, and apps: visible on the user, requestable through Access
Requests, certifiable through campaigns.[^14][^4] Capabilities:

- **Assignment by policy or individually** from the Admin Console, replacing group-per-permission sprawl in Universal Directory.[^4]
- **Bundles** group entitlements into requestable, role-like units; policy-driven assignment keys off profile attributes and updates with profile changes.[^11]
- **Connectors**: GA shipped five enhanced out-of-the-box governance connectors (Salesforce, Google Workspace, Box, NetSuite, Office 365), with more since; SCIM integrations can carry entitlements bidirectionally so an assignment in Okta materializes in the app.[^11][^1][^15]
- **Disconnected apps**: import/update entitlements via CSV, Okta Workflows, or custom connectors when no API exists.[^11][^14]
- **Resource collections** define cross-app sets of apps+entitlements assigned, requested, and certified as one unit — Okta's nearest construct to a business role.[^1][^14]
- Reports: User Entitlements, Past/Active Campaign Details and Summary, Past Access Requests.[^4]

Scope honestly: only a fraction of Okta's 600+ integrations are
governance-deep (entitlement-aware) connectors, a structural point
competitors press;[^17][^10] everything else needs SCIM-with-entitlements,
Workflows, or CSV. Fine-grained certification inside complex apps (e.g.,
Salesforce profile/permission-set nuance) can exceed what the stock connector
represents, pushing teams to custom entitlement modeling.[^14]

## Separation of duties (SOD)

SOD rules **allow, allow-with-oversight, or block specific entitlement
combinations** for apps with the Governance Engine enabled — the canonical
example being create-invoice + approve-payment.[^5] Enforcement is
two-pronged:[^5]

- **Preventative** (Access Requests): block, or allow with custom settings, requests whose grant would create a rule conflict; the Past Access Requests (Conditions) report exposes conflicts via a Conflict name column.
- **Remediative** (Access Certifications): campaigns surface existing conflicts to reviewers (configurable contextual info; also a Governance Analyzer insight) for revocation.[^5][^2]

**Known race condition (design around it):** request-time SOD evaluation
checks only the requester's *existing* entitlement assignments at submission.
Concurrent or rapid-fire open requests for mutually conflicting entitlements
raise no warning and are not blocked. Okta's documented mitigations: check a
requester's other open requests before granting, and run recurring
certification campaigns that review SOD conflicts.[^5] SOD depends on
Entitlement Management (rules are entitlement-combination rules), so apps
without governance-enabled entitlements are outside its reach.[^5][^4]
Cross-application SOD (conflicts spanning two apps' entitlements) is the
classic dedicated-IGA differentiator competitors claim against OIG;[^18]
verify current OIG behavior against your actual conflict matrix before
committing audit language to it.

## Governance reporting: SOX / SOC 2 / ISO evidence

What auditors generally require from an access-review program: evidence that
reviews occurred on schedule, who reviewed, what was decided, that
revocations happened, SoD enforcement, and complete audit trails — with
quarterly cadence typically expected for systems material to financial
reporting (SOX §404(a) management assessment, §404(b) external
attestation).[^21][^20] OIG's evidence surfaces:

1. **Auditor reporting package** (per resource campaign, opt-in checkbox at creation): five reports — Campaign scope; Resource access at launch; Resource access at completion; Campaign actions (decisions); Resource access changes launch→complete. Constraints: the campaign must run >5 hours for the scope report; the access-changes report covers at most the final 90 days of a longer campaign; campaigns including service accounts or resource collections drop the three resource-access reports.[^8][^3]
2. **Standing reports**: Past Campaign Details/Summary, Active Campaign Details/Summary, Past Access Requests, User Entitlements.[^4][^15]
3. **System Log events** for every review decision and request action: streamable to SIEM, usable as Workflows triggers (e.g., ticket every revoke).[^12][^1]
4. **APIs** for campaigns, entitlements, requests, and reports when evidence must land in a GRC system.[^1]

Gaps to plan for: practitioner reviews complain that campaign-result export
to CSV/Excel is limited and request history is hard to mine without external
reporting; prototype the export path early (single practitioner source;
verify in your tenant).[^15] Independent comparisons add that non-editable
timestamped audit artifacts and full-stack campaign coverage (beyond
IdP-connected apps) are where auditors at larger regulated shops push past
native IdP governance.[^16] The control is only as good as its remediation: a
review that flags access but never removes it fails
operational-effectiveness testing.[^20][^21]

## OIG vs dedicated IGA (SailPoint/Saviynt class)

**When OIG suffices** (practitioner consensus, including sources with no Okta
stake): the application estate is predominantly behind Okta SSO/provisioning,
governance needs are standard (periodic certifications, self-service
requests, basic SOD, audit reporting), and speed-to-value matters — OIG
deploys in days-to-weeks vs the multi-month implementations typical of
dedicated platforms.[^16][^10]

**Where the boundary bites** (independent analysis):[^16]

- **Ecosystem boundary**: OIG governs what Okta can see. Apps outside the Okta catalog, on-prem systems without connectors, AWS/GCP IAM, and identities that never touch Okta are out of scope without extra tooling.
- **Depth**: less granular visibility into permissions *inside* downstream apps and in-app activity than dedicated platforms; review delegation/exception granularity and audit-report depth lag at the regulated high end.
- **Mid-market price-to-depth** complaints recur.

**Competitor-published claims** (SailPoint: 250+ bidirectional governance
connectors, scale stats, "only a tiny fraction" of Okta connectors manage
entitlements; Saviynt: converged platform, stronger SoD, role lifecycle, and
a claim that OIG "lacks" SoD that is demonstrably stale given OIG's shipped
SOD feature) are a checklist of what to verify, not findings; both pages
exist to sell against Okta.[^17][^18][^5]

Decision heuristic: bucket every system that must be governed for compliance
by integration path (Okta governance connector / SCIM with entitlements /
Workflows-CSV / impossible). If the "impossible" bucket holds material
systems (mainframe, deep ERP SoD, multi-cloud IAM), OIG alone will not carry
the audit; pair it with a dedicated platform (SailPoint has historically run
as the deep-governance layer beside Okta-as-IdP[^17]). Re-evaluate
periodically: OIG's feature velocity since 2022 has retired several earlier
objections (SOD, entitlement campaigns, auditor packages).[^11][^5][^8]

## Lifecycle and Workflows integration

Governance and lifecycle are designed as one loop: LCM grants birthright
access on joiner events; Access Requests handles ad-hoc grants between
lifecycle events (time-boxed, auto-expiring via conditions); certifications
and SOD prune what accumulated; Workflows stitches custom fulfillment for
whatever has no connector.[^1][^6] Concretely:

- OIG-specific System Log events (request approved, campaign completed, entitlement revoked) can trigger Workflows for ticketing, notifications, or custom deprovisioning.[^1][^15]
- Approval sequences can invoke Workflows **delegated flows** mid-request.[^6]
- Entitlement changes propagate through SCIM-with-entitlements or Workflows; mover events re-evaluate policy-assigned entitlements via profile attributes.[^1][^11]
- For provisioning/deprovisioning mechanics (SCIM, OPP, orphaned access) see `okta-scim-lifecycle-security`; for governing Okta *admin* roles see `okta-admin-hardening` (both in this hub).

## Anti-patterns in certification programs

### Rubber-stamping

The dominant failure mode: reviewers bulk-approve without scrutiny, which
"completely negates any access controls currently in place"; auditors
increasingly treat no-justification bulk approval as a finding.[^19] Root
causes converge across independent sources: reviewer time burden, lack of
access-domain understanding, review fatigue at scale (10 reports × ~200 apps
= thousands of items),[^19] and poor data quality (ambiguous group names,
missing ownership, jargon entitlements) that makes blind approval the
rational move.[^20]

OIG-native mitigations, mapped to cause:

| Cause | Mitigation in OIG |
| --- | --- |
| No context | Reviewer context panels (assignment date, last access, usage); Governance Analyzer insights/recommendations (EA); customizable reviewer context[^13][^2] |
| Wrong reviewer | Resource-owner/group-owner reviewer types; fallback hygiene; disable self-review[^3] |
| Fatigue | Smart Review grouping; narrow campaign scopes; preconfigured inactive-user campaigns to clear deadwood first[^2][^13] |
| Opaque names | Friendly display names/descriptions for groups and entitlement bundles[^7][^11] |
| No accountability | Multilevel review with visible level-1 justification; decision + justification in System Log and reports[^3][^12] |

### Scope explosion

Campaigns that certify everything-everywhere-annually produce fatigue and
theater. Counter-design: risk-tier the estate (the `Crown Jewel` /
`Privileged` labels exist for this[^3]); review privileged and
material-to-financials access quarterly and the long tail less often;[^20][^21]
prefer several narrow campaigns (one app family, one admin surface) over one
giant one; and use event-driven security access reviews for
incident-adjacent checks instead of widening periodic campaigns.[^2] The
preconfigured templates encode this: top-5-inactive-apps, not all apps.[^13]

### Other recurring failures

- **Weak remediation**: decisions without enforced revocation. Wire campaign remediation to auto-remove (and ticket manual cases via Workflows); mark manual remediation complete so the auditor package's completion report reflects reality.[^8][^20]
- **Spreadsheet evidence**: manual UARs in spreadsheets fail timestamp/traceability expectations and are consistently flagged in SOX assessments.[^21]
- **Unreviewed non-human identities**: service accounts carry privilege but get skipped; OIG certifies service accounts only when managed in Okta Privileged Access — otherwise they need a separate control.[^2][^20]
- **Annual-only cadence**: event-driven reviews (role change, transfer, termination, incident) catch what anniversary campaigns miss.[^20]

## Review checklist (condensed)

1. Licensing: confirm what the org owns (OIG add-on vs Essentials+ suite; Workflows flow count).[^9]
2. Campaign inventory: types, cadence, risk-tiered scope; auditor package enabled on audit campaigns.[^3][^8]
3. Reviewer hygiene: manager attributes populated, fallback reviewers, self-review disabled, delegates configured.[^3]
4. Requests: conditions preferred for new builds; durations set so access auto-expires; group-snapshot and tile-refresh gotchas handled.[^6][^7]
5. Entitlements: Governance Engine on material apps; bundles with human-readable descriptions; disconnected-app import path chosen.[^4][^11]
6. SOD: rules cover the org's actual toxic combinations; concurrent-request race mitigated by recurring SOD-aware campaigns.[^5]
7. Evidence: export path to GRC/auditors prototyped; remediation completion tracked; System Log streamed.[^8][^15]
8. Boundary: list the material systems OIG cannot govern, with a named compensating control for each.[^16]

## Known ambiguities and guardrails

- Feature availability moves monthly and varies by edition/EA-GA status
  (Governance Analyzer was Early Access at verification time); confirm in
  *your* org before writing policy or audit language.
  verified-as-of: 2026-06-11.[^2]
- Pricing/packaging figures are list-price snapshots; enterprise quotes
  diverge. verified-as-of: 2026-06-11.[^9][^22]
- Several capability criticisms (export limits, bulk-assignment instability,
  ITSM gaps) are single-practitioner reports from a vendor with a competing
  product; treat as hypotheses to test in a trial tenant.[^15]
- Competitor comparison pages contain at least one verifiably stale claim
  about OIG; never cite them as evidence of OIG's current capabilities
  without checking the docs.[^17][^18][^5]
- Public Sector Service (FedRAMP-class) OIG has documented functional
  limitations relative to commercial; check the limitations page for gov
  tenants.[^1]
- Sibling spokes named here (`okta-scim-lifecycle-security`) belong to the
  same /dr family pack and may land in this hub's `references/` slightly
  after this file; until then, treat those pointers as forward references.

## References

All sources fetched and verified 2026-06-11.

1. [^1] Okta Documentation, "Identity Governance" (OIG overview: LCM + Workflows + Access Governance, components, SOC2/SOX positioning), vendor product docs. <https://help.okta.com/oie/en-us/content/topics/identity-governance/iga.htm>
2. [^2] Okta Documentation, "Access Certifications" (campaigns vs security access reviews, personas, Governance Analyzer EA, Smart Review), vendor product docs. <https://help.okta.com/oie/en-us/content/topics/identity-governance/access-certification/iga-access-cert.htm>
3. [^3] Okta Documentation, "Best practices for creating campaigns" (reviewer types, fallback chains, multilevel review, labels, auditor package checkbox), vendor product docs. <https://help.okta.com/oie/en-us/content/topics/identity-governance/access-certification/best-practices-create-campaign.htm>
4. [^4] Okta Documentation, "Entitlement Management" (policy/individual assignment, bundles, report list), vendor product docs. <https://help.okta.com/oie/en-us/content/topics/identity-governance/em/entitlement-mgt.htm>
5. [^5] Okta Documentation, "Separation of duties" (rule semantics, preventative/remediative enforcement, request-time evaluation race caveat), vendor product docs. <https://help.okta.com/oie/en-us/content/topics/identity-governance/sd/separation-of-duties.htm>
6. [^6] Okta Documentation, "Access Requests" (conditions vs request types, approval sequences, omnichannel, admin-role limits), vendor product docs. <https://help.okta.com/en-us/content/topics/identity-governance/access-requests/ar-overview.htm>
7. [^7] Okta Documentation, "Access request conditions" (setup, duration auto-revoke, privacy, group-snapshot and tile-refresh gotchas, escalation, delegates), vendor product docs. <https://help.okta.com/en-us/content/topics/identity-governance/access-requests/rcar-conditions.htm>
8. [^8] Okta Documentation, "Auditor reporting package" (five reports, >5h constraint, 90-day window, service-account/collection exclusions), vendor product docs. <https://help.okta.com/oie/en-us/content/topics/identity-governance/auditor-reporting/auditor-report-pkg.htm>
9. [^9] Okta, "Plans and Pricing" (suite prices, OIG add-on composition footnote, $1,500 minimum, unified license model), vendor pricing page. <https://www.okta.com/pricing/>
10. [^10] Okta, "Identity Governance" product page (600+ integrations claim, customer ROI figures), vendor marketing. <https://www.okta.com/products/identity-governance/>
11. [^11] Okta Blog, "Announcing Entitlement Management for Okta Identity Governance," 2024-01-22 (GA date, five launch connectors, bundles, disconnected resources), vendor product blog. <https://www.okta.com/blog/product-innovation/announcing-entitlement-management-for-okta-identity-governance/>
12. [^12] Okta Support KB, David Edwards, "User Access Reviews in Okta Identity Governance," 2023-07-25 (user campaigns added 2023-06, resource-vs-user framing, System Log/Workflows hooks), vendor KB, practitioner-authored. <https://support.okta.com/help/s/article/user-access-reviews-in-okta-identity-governance?language=en_US>
13. [^13] IAMSE.blog, David Edwards, "Preconfigured Access Certification Campaigns in Okta Identity Governance," 2025-02-11 (release 2025.01.0, two templates, default no-remediation, context panels), Okta-employee practitioner blog. <https://iamse.blog/2025/02/11/preconfigured-access-certification-campaigns-in-okta-identity-governance/>
14. [^14] IAMSE.blog, "OIG Entitlement Management" topic hub (entitlements as first-class objects, SoD intro, disconnected-app imports, collections, Salesforce fine-grained limits), Okta-employee practitioner blog. <https://iamse.blog/identity-governance-and-administration-iga/oig-entitlement-management/>
15. [^15] Multiplier, "Okta Identity Governance: The Good, The Bad, and The Catch" (2022 launch, add-on model, report types, export/ITSM criticisms; competing-vendor analysis), independent vendor analysis. <https://multiplierhq.com/blog/okta-identity-governance-the-good-the-bad-and-the-catch>
16. [^16] Zluri, "Entra ID Governance vs Okta Identity Governance," 2026-05 (acquisition origin, ecosystem-boundary framing, when native IdP governance suffices), independent vendor analysis. <https://www.zluri.com/eye-on-identity/entra-id-governance-vs-okta-identity-governance-iga-comparison>
17. [^17] SailPoint, "SailPoint vs Okta" (add-on positioning, 250+ connector claim, entitlement-depth critique), competitor comparison; bias noted. <https://www.sailpoint.com/compare/sailpoint-vs-okta>
18. [^18] Saviynt, "Saviynt vs Okta" (converged-platform and SoD claims, including a stale "lacks SoD" assertion), competitor comparison; bias noted. <https://saviynt.com/solution-comparison/saviynt-vs-okta>
19. [^19] Clarity Security, "Why Managers Rubberstamp User Access Reviews," 2023-07-20 (definition, three root causes, risk framing), independent IGA vendor analysis. <https://claritysecurity.com/clarity-blog/why-managers-rubberstamp-uars/>
20. [^20] Palo Alto Networks Cyberpedia, "What Is Access Certification?" (process steps, certification types, challenges, review-frequency norms), independent vendor education. <https://www.paloaltonetworks.com/cyberpedia/what-is-access-certification>
21. [^21] BalkanID, "User Access Reviews for SOX Compliance," 2025-09-30 (SOX §404(a)/(b), auditor expectations, quarterly cadence, SoD toxic combinations, spreadsheet trap), independent IGA vendor guide. <https://www.balkan.id/post/how-user-access-reviews-help-organizations-achieve-sox-compliance>
22. [^22] UnderDefense, "Okta Pricing 2026: Ultimate Guide" (independent corroboration of suite pricing, OIG bundle footnote, excluded-cost list), independent pricing guide. <https://underdefense.com/industry-pricings/okta-pricing-ultimate-guide-for-security-products/>
