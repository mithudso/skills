# Output contract + escalation matrix

Loaded on demand by `solve-case` Phase 7. This is the shape of the final deliverable and the
rules for the escalation call.

## Phase-7 deliverable template

```markdown
## Case <id> — <customer> — <one-line presenting problem>
Severity: <S1–S4 or "TBD">   |   Status: <stage>   |   As of: <UTC timestamp>

### 1. Fact-based analysis
- <finding> — evidence: <KB Public URL | Manual page | metric/log/explain stat> — confidence: <high|med|low>
- <hypothesis, explicitly labeled> — what would confirm/refute it: <test>
- Root-cause assessment: <most-likely cause + why the evidence points there>

### 2. Drafted customer response  (human reviews before sending)
<reply in the TAM voice — acknowledges, states findings honestly, gives next steps,
sets expectations; Public KB links only; no Internal content>

### 3. Blockers
- <missing data / access / customer input / pending internal answer> — needed from: <who>

### 4. Tools used
- Skills: <...>   |   MCP tools: <...>   |   Agents: <...>

### 5. Who to talk to
- <role/person> — for <reason>   (see roles below)

### 6. Escalate?
- Decision: <yes/no>   |   Severity: <S1–S4 / SEV-n>
- Action: <file HELP ticket | open JIMP | request Break-Glass | start PIR-RCA | none>
- Why: <criterion met from the matrix below>

### 7. Insights
- <patterns, proactive follow-ups, account-health signals, prevention>
```

## Escalation matrix

Severity is the customer-impact call; the **action** is what you open. Confirm exact current
thresholds against `tam-operations`; these are the standing defaults.

| Severity | Signal | Typical action |
|---|---|---|
| **S1 / SEV-1** | Production down, data loss/at risk, security breach, no workaround | Immediate escalation; page IR; consider Break-Glass; customer comms cadence |
| **S2 / SEV-2** | Severe degradation, major feature unusable, tight workaround | Escalate to NTSE/IR; HELP ticket; frequent updates |
| **S3** | Moderate impact, workaround exists | Normal support flow; HELP ticket if it needs Engineering |
| **S4** | Question, minor issue, guidance | Handle directly; document |

**When to open which:**
- **HELP ticket** — the case needs MongoDB Engineering/Support input beyond TAM scope, or a
  bug/known-issue confirmation. (`mdb_case_get_help_ticket` to check existing.)
- **JIMP** (Joint Incident Management Plan) — coordinated multi-party incident on a Premium account.
- **Break-Glass** — emergency elevated access needed to mitigate an active S1.
- **PIR-RCA** — after a significant incident, to produce the post-incident review / root-cause analysis
  (consider the `incident-postmortem-drafter` agent).

## Who-to-talk-to roles (MongoDB Premium Services)

| Role | Talk to them for |
|---|---|
| **TAM** | Account ownership, relationship, coordination (you, usually) |
| **NTSE** (Named Technical Support Engineer) | Deep technical support, reproduction, Engineering bridge |
| **IR** (Incident Response) | Active SEV-1/2 coordination, paging, comms |
| **DCE** (Dedicated Consulting Engineer) | Architecture/design guidance, proactive engagements |
| **SME** | Subdomain depth (Search, sharding, drivers, security, etc.) |
| **Support escalation manager** | When SLA risk or customer escalation needs management attention |

Name the specific person from the account's `mdb_tam_account_context` snapshot when known; otherwise
name the role and note "identify the current <role> for this account."
