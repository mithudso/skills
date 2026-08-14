# solve-case

**Category:** Science, Biology & Medicine
**Platform:** Claude
**Original Path:** claude/standalone/solve-case

## Description
End-to-end MongoDB/Atlas support-case solver. Orchestrates the mongodb-* and 10gen skills, the case and account MCPs, and the diagnostic and psychology agents into one workflow: identify the customer, run full troubleshooting and diagnosis (deep-diving any unknown), give a fact-based cited analysis, then craft a customer-psychology-informed reply bundled with blockers, tools used, who to talk to, an escalation call, and insights. TRIGGER: "solve this case", "solve case <id|url>", "work this case end to end", "diagnose and respond to this MongoDB/Atlas case", or taking a support case from intake to a drafted customer reply. SKIP: one narrow MongoDB symptom with no customer/response/escalation wrap-up (use atlas-diagnostics-expert or mongodb-kb); drafting a reply with no diagnosis (use content-and-marketing-writing); backtesting diagnosis methodologies (use diagnosis-methodology-backtest); building or auditing another skill (use skill-creator or skill-optimizer).

---

# Solve Case (MongoDB/Atlas, intake to drafted response)

This is an **orchestration skill**. It does not replace the specialist `mongodb-*` / `10gen`
skills; it sequences them, the case/account MCPs, and the diagnostic and psychology agents into one
repeatable workflow that takes a support case from "look at it" to a drafted, fact-checked,
psychology-aware customer reply plus an operations wrap-up.

## When to use
Trigger on **"solve this case"** (or `solve case <id|url>`, "work this case end to end"). Use when
the goal is the *whole* arc (customer, diagnosis, analysis, response, escalation), not a single
slice. For one narrow symptom, go straight to `atlas-diagnostics-expert` / `mongodb-kb` instead.

## Operating guardrails (read before acting)
- **Read-only diagnosis.** Inspect, query, and read; never run destructive operations (drops,
  deletes, writes, production changes, irreversible API calls). If a fix would require one, describe
  it; do not execute it.
- **Public-only when sharing externally.** Only `support.mongodb.com` **Public** KB URLs may appear
  in customer-facing text. Never paste Internal KB content to the customer.
- **The human owns the send.** This skill drafts; it does not send. End with a draft the TAM reviews.
- **Cite every factual claim** (Phase 5). Separate established fact from hypothesis.
- **PII/secrets.** Advise redaction of credentials, connection strings, and personal data before any
  external use.

## Inputs
- A **case id or URL**, an **account/customer name**, or "the currently tracked case."
- If none is given, resolve the active case via `mdb_case_assistant` (`mdb_case_get_tracking_state`
  or `mdb_case_list_account_cases`) and confirm which case in one line before proceeding.

## Workflow

Run the phases in order, carrying state forward. Each phase names the skills/MCPs/agents to use.
Depth lives in `references/`; load it when a phase needs it.

### Phase 0 — Intake & framing
Resolve the case identifier and restate the ask in one line (case, customer, presenting symptom,
severity if known). Set the guardrails above. Decide solo-vs-dispatch: for a broad or
multi-subdomain case, plan to dispatch the `uber-mongodb-diagnostician` agent in Phase 3; for a
focused one, drive the specialist skills directly.

### Phase 1 — Assemble the right expertise
Select only the skills that fit the case's symptoms; **do not load all 62 mongodb-* skills.**
Route from symptom to skill using `references/diagnostic-playbook.md`. Hubs to draw from:
`mongodb-expert` (data-plane/engine: CRUD, aggregation, indexes, schema, WiredTiger, sharding,
replication), `mongodb-atlas-expert` (control plane/platform), `atlas-diagnostics-expert` (live
perf, FTDC, ts-diag, Performance Advisor), `mongodb-operations-expert` (backup/DR, migration,
security, encryption, cost), `mongodb-kb` (symptom-to-article, error codes), `mongodb-docset-lookup`
(Manual text), and `10gen` (map a symptom to a 10gen diagnostic repo/tool plus install/run guidance).

### Phase 2 — Load the case & identify the customer
Via the **`mdb_case_assistant`** MCP: `mdb_case_get_case`, `mdb_case_get_case_comments`,
`mdb_case_get_account`, `mdb_case_get_help_ticket`, `mdb_case_get_case_stage`. Via
**`mdb_tam_account_context`**: `mdb_tam_corpus_search` and `mdb_tam_snapshot_latest` for account
history. Produce a short **customer profile**: who they are, Atlas tier / topology / server version,
business context, support tier/SLA, and relevant prior cases. For a deep account pull, dispatch the
**`account-data-collector`** agent.

### Phase 3 — End-to-end troubleshooting & diagnosis
Reconstruct the timeline (what changed, when, error codes, what the customer already tried). Form a
**rank-ordered list of root-cause hypotheses** with a confidence each. Gather diagnostic evidence:
- live read-only data via the **`mongodb`** MCP (`explain`, `collection-indexes`,
  `collection-schema`, `db-stats`, `mongodb-logs`, `aggregate`, `find`);
- Atlas diagnostic surfaces via `mdb_tam_diagnostics_build_cluster_url`,
  `mdb_tam_diagnostics_build_snapshot_url`, and `mdb_tam_diagnostics_status_summary` (FTDC, metrics);
- KB matches via `mongodb-kb` (and `mcp__plugin_mongodb_mongodb__search-knowledge`).

For breadth or cross-subdomain cases, **dispatch the `uber-mongodb-diagnostician` agent** (read-only;
returns rank-ordered hypotheses, evidence-to-collect, remediation, confidence, and citations) and
fold its output in. See `references/diagnostic-playbook.md` for the evidence checklist.

### Phase 4 — Deep-dive any unknown
For every component, error, feature, default, or version behavior you are not certain about, **stop
and resolve it** before relying on it: `mongodb-docset-lookup` (Manual), the matching `mongodb-*`
spoke, `mongodb-kb`, `context7` for driver/library docs, or web research as a last resort. Do not
hand-wave an unknown into the analysis.

### Phase 5 — Fact-based analysis
Write findings as **fact, with evidence and citation**: a KB Public URL, a Manual page, or the
specific metric / log line / explain stat that supports each claim. Tag confidence (high/medium/low)
and keep hypotheses explicitly labeled. Optionally verify with the **`mongodb-claim-validator`**
agent (technical claims) and/or **`tam-doc-validator`** (against corpus plus live case state).

### Phase 6 — Customer psychology + craft the response
Apply **`applied-psychology`** to shape the reply: trust repair (Mayer ABI, competence vs integrity),
avoid psychological reactance, match the stakeholder's stage, calibrate confidence honestly. Draft in
the TAM voice with **`content-and-marketing-writing`** (support-reply and escalation-handoff craft),
then de-robot with **`kill-the-AI-ism`** and audit with **`document-critique`**. For the highest-stakes
or trust-damaged cases, **dispatch the `customer-comms-psychologist` agent** (drafts AND pressure-tests
against behavioral science).

### Phase 7 — Wrap-up bundle
Emit the full deliverable using the template in `references/output-contract.md`:
1. **Fact-based analysis** (cited, confidence-tagged)
2. **Drafted customer response** (ready for TAM review; human owns send)
3. **Blockers** (what is missing or waiting: data, access, customer input, internal answer)
4. **Tools used** (skills, MCP tools, and agents actually invoked)
5. **Who to talk to** (internal roles/people: NTSE / IR / DCE / SME / Support escalation)
6. **Escalate?** yes/no plus the severity call (S1–S4 / SEV) and whether to file HELP / JIMP /
   Break-Glass / PIR-RCA, per the matrix in `references/output-contract.md`
7. **Insights** (patterns, proactive follow-ups, account-health signals)

## Skill / MCP / agent map

| Need | Use |
|---|---|
| Author/quality-gate this workflow | `claude-code-skills`, `skill-creator`, `skill-optimizer` |
| Core MongoDB engine/data-plane | `mongodb-expert` |
| Atlas platform/control-plane | `mongodb-atlas-expert` |
| Live perf/FTDC/diagnostic surface | `atlas-diagnostics-expert` |
| Backup/DR, migration, security, cost | `mongodb-operations-expert` |
| KB articles, error codes, Public URLs | `mongodb-kb` |
| Authoritative Manual text | `mongodb-docset-lookup` |
| 10gen diagnostic repos/tools | `10gen` |
| Case ops, escalation, SEV, who-to-talk-to | `case-mcp-server-guide`, `tam-operations` |
| Customer psychology | `applied-psychology` |
| Reply + handoff drafting / human voice | `content-and-marketing-writing`, `kill-the-AI-ism`, `document-critique` |
| **MCP:** case load/analyze/account/HELP | `mdb_case_assistant` |
| **MCP:** account context, corpus, Atlas diag URLs | `mdb_tam_account_context` |
| **MCP:** read-only data-plane diagnostics | `mongodb` (plugin) |
| **Agent:** deep root-cause diagnosis | `uber-mongodb-diagnostician` |
| **Agent:** gather customer/account data | `account-data-collector` |
| **Agent:** fact-check claims | `mongodb-claim-validator`, `tam-doc-validator` |
| **Agent:** psychology-tested reply | `customer-comms-psychologist` |
| **Agent:** post-incident RCA/postmortem | `incident-postmortem-drafter` |

## Degradation
If a skill, MCP, or agent is unavailable (case MCP offline, no Atlas access, and so on), **state what
is missing**, proceed with what you have, and record the gap as a **blocker** in Phase 7. Never invent
case data or evidence to fill a gap. If the case cannot be loaded at all, stop and report the blocker
rather than guessing.

## References
- `references/diagnostic-playbook.md`: symptom-to-skill/tool routing and the evidence checklist.
- `references/output-contract.md`: the Phase-7 deliverable template plus the escalation matrix
  (S1–S4 / SEV criteria, HELP/JIMP/Break-Glass/PIR triggers, who-to-talk-to roles).