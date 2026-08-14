# tam-operations

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude Code
**Original Path:** claude-code/tam-operations

## Description
TAM operations hub — account deliverables, customer health, case & incident ops, reporting automation, Customer Dashboard support-data. TRIGGER: TAM deliverables — EBR/QBR, account review, support plan, engagement overview, JIMP, weekly update, case-analysis report; account health, churn risk, NRR, sentiment, health-scoring; MongoDB Premium Services operating model (TAM/NTSE/IR/DCE roles, S1–S4 SLAs, JIMP/Break-Glass/PIR-RCA, Straight-to-8, risk register, Monday tracking); account-artifact collection; shift-handoff (SBAR), meeting-prep, operator reports (freshness, BLUF); customer file consolidation; active case tracker (Support API, severity, diagnostic catalog, LLM summarization); case/event timeline viz; incident response (SEV, IC/comms, postmortems, SLO/error budgets, on-call, MTTD/MTTR); Monday.com board audit; Slack subscription audit. SKIP: vendor API client dev → integration-clients; MongoDB technical → mongodb-* hubs.

---

# TAM Operations

Hub for the **MongoDB Technical Account Manager's operational toolkit** — producing account
deliverables, scoring customer health, running case and incident operations, automating reports, and
integrating the Customer Dashboard's support data. Each former standalone skill is now an on-demand
reference: when a task matches a row below, **`Read` that reference file before answering**.

Boundary: this hub owns the **TAM operator's own work** — deliverables, health, case/incident process,
reporting, and the support-API integration the dashboard depends on. When the question is about
MongoDB/Atlas *technical* depth, the *prose craft* of a deliverable, or *MCP tooling* mechanics,
defer to the sibling hubs listed in `whenNotToUse`.

<!-- ROUTING TABLE: tam-operations -->
## Sub-skill routing table

This hub absorbs 14 flat reference files plus 2 nested sub-skills (16 spokes total). When a task
matches a row, **Read the listed `references/` file** before answering — do not rely on this table
alone for depth.

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `tam-expertise` | TAM deliverables, health/churn/NRR, success frameworks (BLUF/STAR/Pyramid/SCQA), 30-60-90 onboarding, escalation/stakeholder dynamics | `references/tam-expertise.md` |
| `tam-commercial-metrics` | Commercial-metric definitions & benchmarks — NRR/GRR formulas + worked examples + segment benchmarks, MEDDPICC framework & scoring (incl. expansion/renewal application), dated 2025-2026 SaaS retention/churn/expansion figures | `references/tam-commercial-metrics.md` |
| `tam-reference` | MongoDB Premium Services TAM operating reference — TAM/NTSE/IR/DCE roles, S1–S4 SLAs, JIMP/Break-Glass/PIR-RCA, lifecycle, Straight-to-8, risk register, Monday tracking | `references/tam-reference.md` |
| `tam-account-reports` | Generate a named-account document from live MCP data — account review, support plan, engagement overview, JIMP, weekly update, case-analysis report | `references/tam-account-reports.md` |
| `account-health-scorer` | Health-scoring algorithms — weighted composites, signal categories, grading heuristics, time-series trending, anomaly detection | `references/account-health-scorer.md` |
| `account-artifacts-collector` | Parallel data collection across MCP + local sources, persisting typed JSON artifacts for report generation / agent consumption | `references/account-artifacts-collector.md` |
| `operator-report-generator` | Operator report engines — shift handoff (SBAR), meeting prep, data-freshness scoring, BLUF, quality validation, template-driven rendering | `references/operator-report-generator.md` |
| `customer-file-consolidator` | Collect, dedupe, version-consolidate, and semantically analyze a customer's local files into a unified TAM briefing, then ingest to corpus | `references/customer-file-consolidator.md` |
| `case-tracker` | Build/extend the Customer Dashboard active case tracker — TS Tools API, case schema, severity model, subscription flows, LLM summarization, diagnostic catalog | `references/case-tracker.md` |
| `case-timeline-visualization` | Vanilla-JS temporal/event-sequence visualization for dashboards — DOM/SVG/Canvas, zoom/scroll, accessibility, performance | `references/case-timeline-visualization.md` |
| `incident-response` | Incident lifecycle — SEV classification, IC/comms/scribe roles, blameless postmortems, SLO/SLI/SLA + error budgets, on-call, MTTD/MTTR, chaos prep | `references/incident-response/SKILL.md` |
| `autoremediation` | Self-healing systems — retry/circuit-breaker, automated recovery, graceful degradation, canary rollback, AI-assisted repair loops | `references/autoremediation.md` |
| `firedrill-integration-tester` | Firedrill/game-day validation — scenario execution, safety/abort, scoring rubrics, multi-step agent orchestration via the mdb_case_assistant firedrill tools | `references/firedrill-integration-tester/SKILL.md` |
| `ts-tools-support-api` | TS Tools Support API implementation patterns — auth cascade, response normalization, getCaseBundle, retry/backoff, enrichment, JIRA extraction, Socket.IO | `references/ts-tools-support-api.md` |
| `tstools-reference` | TS Tools Support API endpoint reference — all 92 endpoints, request/response shapes, cookie-vs-Bearer auth, Socket.IO event model, bootstrap payload | `references/tstools-reference.md` |
| `monday-board-audit` | Autonomous Monday.com board audit — status/priority/link/title updates, 14-day freshness, gap-analysis item creation, markdown audit log | `references/monday-board-audit.md` |
| `slack-subscription-auditor` | Slack workspace/app audit — slash-command inventory, dead-command & orphaned-subscription detection, bot-scope least-privilege, rate limits | `references/slack-subscription-auditor.md` |

<!-- cross-hub-map -->
## Cross-hub map — where to route tasks outside this hub

All 16 TAM operations spokes live in this hub. Tasks that fall outside the Sub-skill routing table
above belong in a sibling hub. Use `whenNotToUse` as the routing guide.

| Sibling hub | Route tasks about… |
| --- | --- |
| `mongodb-expert` / `mongodb-atlas-expert` / `atlas-diagnostics-expert` | MongoDB/Atlas data-plane, engine, platform, or live diagnostics |
| `mongodb-kb` | MongoDB error-code or KB-article lookup |
| `case-mcp-server-guide` | Starting/using the local case MCP server; choosing `mdb_case_*` tools |
| `writing-expert` / `technical-writing-craft` / `executive-comms` | Prose craft, voice, or editing of any TAM deliverable |
| `da-*` hubs | Metric modeling, funnels, cohorts, statistical forecasting |
| `content-ingestion-extraction` | Deconstructing/templatizing an existing document into a reconstruction prompt |