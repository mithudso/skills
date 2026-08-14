<!-- hub-reference-banner -->
> **Reference file — part of the `content-ingestion-extraction` hub.** Formerly the standalone `document-deconstructor` skill.
> Sibling topics in this family are now reference files under the hubs (`content-ingestion-extraction`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: document-deconstructor
description: "Reverse-engineer any structured document into a reconstruction prompt with {{placeholder}} variables and a data-manifest.json mapping each placeholder to its source, extraction method, and last-known value. TRIGGER: 'deconstruct this document', 'templatize this report', 'make this document reproducible', 'make it auto-updatable', 'reverse-engineer this document', 'create a template from this document', 'extract data points from', 'document to prompt', 'make this repeatable', turn a one-off account review or QBR into an automated generation pipeline. SKIP: critique or review a document without templatizing it (use document-critique); generate a report from scratch (use operator-report-generator); edit a .docx file directly (use docx skill)."
version: "1.2.1"
updated: "2026-05-31"
origin: local
category: developer
tags: [document-analysis, template-generation, reconstruction-prompt, data-manifest, reverse-engineering, templatize, reproducible-documents, account-review, weekly-update, QBR]
whenToUse:
  - "deconstruct this document into a template"
  - "templatize this account review or QBR"
  - "make this document reproducible with fresh data"
  - "create a reconstruction prompt for this report"
  - "extract all data points and trace them to sources"
  - "turn this one-off document into an automated pipeline"
whenNotToUse:
  - "critique or review a document — use document-critique"
  - "generate a report from scratch — use operator-report-generator or writing-expert"
  - "edit a .docx file directly — use the docx skill"
  - "create a new template without an existing document — use writing-expert + tam-expertise"
related_skills: [document-critique, operator-report-generator, writing-expert, tam-expertise, case-mcp-server-guide]
---

# Document Deconstructor

Reverse-engineers any structured document into two artifacts that make the document reproducible with fresh data:

1. **`<name>.reconstruction-prompt.md`** — a complete prompt template with `{{placeholder}}` variables for every data point, ready to feed to an LLM with updated values.
2. **`<name>.data-manifest.json`** — maps each placeholder to its data source, extraction method, field path, and last-known value.

Together these artifacts turn a one-off document into a repeatable generation pipeline.

## When to use this skill

Activate when the user:

- Asks to "deconstruct," "templatize," or "reverse-engineer" a document
- Asks to "create a reconstruction prompt for" a document
- Says "make this document reproducible" or "make it auto-updatable"
- Provides a document path and asks to "make this repeatable"
- Wants to turn a manually written report into an automated generation pipeline
- Needs to extract all data points from a document and trace them to sources
- Wants to understand what data feeds into an existing deliverable

## When NOT to use

- The user wants to **critique or review** a document (use `document-critique`)
- The user wants to **generate a report from scratch** (use `operator-report-generator` or `writing-expert`)
- The user wants to **edit a .docx file** (use `docx`)
- The user wants to **create a new template without an existing document** (use `writing-expert` + `tam-expertise`)

## Supported document types

This skill handles any structured text document. Common types with domain-specific awareness:

| Document type | Typical sections | Primary data sources |
|---|---|---|
| **Account Review** | Executive summary, case portfolio, SLA metrics, architecture notes, engagement history, risk assessment, renewal timeline | Case MCP, Salesforce/CRM, Atlas metrics, Tableau exports, manual TAM notes |
| **TAM Engagement Overview** | Account profile, support plan, escalation contacts, architecture summary, open initiatives, meeting cadence | CRM, case history, architecture docs, calendar, org charts |
| **Joint Incident Management Plan (JIMP)** | Severity definitions, escalation matrix, communication channels, RCA process, SLA targets | Internal policy docs, customer agreements, contact directories |
| **TAM Support Plan** | Scope, deliverables, meeting schedule, escalation paths, SLA targets, success criteria | Contract docs, CRM, internal planning tools |
| **Weekly Update** | BLUF, case status, action items, blockers, upcoming milestones, metrics snapshot | Case MCP, action-item tracker, calendar, monitoring dashboards |
| **Case Analysis Report** | Case timeline, root cause, impact assessment, resolution steps, prevention recommendations | Case MCP, case comments, logs, KB articles, engineering threads |
| **QBR Deck** | Executive summary, SLA performance, case trends, architecture roadmap, risk register, recommendations | Tableau/BI exports, case history, SLA dashboards, roadmap docs |
| **Incident Post-Mortem** | Timeline, root cause, impact, detection, resolution, action items, lessons learned | Incident logs, monitoring, case threads, engineering analysis |

## Hard constraints

- **Two output files, always.** Every invocation produces both the reconstruction prompt and the data manifest. Never produce one without the other.
- **Self-contained prompt.** The reconstruction prompt must be executable by an LLM that has never seen the original document. No references to "the original" or "see above."
- **No data fabrication.** If a data point's source cannot be inferred, mark it `unknown` and flag it. Never guess a source and present it as `high` confidence.
- **Preserve original voice.** The reconstruction prompt must describe the original document's tone, not impose a new one.
- **Cross-validation is mandatory.** Step 7 validation must pass before delivering. If it fails, fix the mismatch before outputting.

## Input handling

The skill accepts documents in any of these forms:

| Input form | How to handle |
|---|---|
| **File path** (`.md`, `.txt`, `.json`, `.yaml`) | Read directly with the `Read` tool |
| **File path** (`.docx`) | Use the `docx` skill's reading workflow (pandoc or unpack) |
| **File path** (`.pdf`) | Use the `pdf` skill's reading workflow |
| **Raw text pasted in chat** | Treat the message content as the document; use "untitled-document" as the base name; ask the user for a preferred output directory |
| **Multiple files** | Deconstruct each file independently, then produce a combined manifest with a shared `source_summary` if the user asks for a unified view |
| **URL** | Fetch the content first (using `WebFetch` if available), then proceed as with raw text |

**Error handling:**

- If the file cannot be read (permission denied, binary, corrupt), report the error and stop. Do not guess at content.
- If the document is too short to meaningfully deconstruct (fewer than 3 data points), warn the user and ask whether to proceed.
- If the document format is unsupported, suggest converting it first and name the appropriate tool.

## Execution steps

Run these steps in order. Each step builds on the previous.

### Step 1: Read and inventory the document

Read the document using the `Read` tool (for text files) or the appropriate reader (for .docx, .pdf).

Produce an **inventory** with:

- **File path and format** (markdown, docx, PDF, text)
- **Document type** (match against the table above, or classify as "custom" with a description)
- **Total section count**
- **Word count estimate**
- **Date context** — any AS-OF date, reporting period, or temporal anchors in the document
- **Account/customer context** — if the document is about a specific account, note the account name, org ID, or other identifier

### Step 2: Analyze structure

Map the document's skeleton. For each section and subsection, record:

```
Section: <heading text>
  Level: <H1/H2/H3/...>
  Content type: <narrative | data-table | bullet-list | metric-block | timeline | checklist | mixed>
  Approximate word count: <number>
  Contains data points: <yes/no>
  Formatting conventions: <bold labels, inline code, ISO dates, severity tags, etc.>
```

Also identify document-level patterns:

- **Narrative patterns** — BLUF, chronological, severity-ordered, categorical
- **Citation style** — inline links, footnote references, case ID references, none
- **Table structure** — column headers, alignment, units, conditional formatting cues
- **Naming conventions** — how entities are referenced (full name vs. abbreviation, case ID format, cluster name format)
- **Recurring structural motifs** — e.g., every case has the same sub-bullet pattern (severity, status, owner, last update)

### Step 3: Extract data points

Walk every section and extract every discrete data point. A data point is any fact, number, date, name, identifier, URL, status, or measurement that would change if the document were regenerated with fresh data.

For each data point, record:

| Field | Description |
|---|---|
| `placeholder` | A descriptive `{{snake_case_name}}` variable name |
| `section` | Which section the data point appears in |
| `value` | The exact value as it appears in the document |
| `data_type` | `string`, `number`, `date`, `boolean`, `enum`, `url`, `identifier`, `list` |
| `context` | The surrounding sentence or label that gives the data point meaning |
| `is_derived` | Whether this value is computed from other data points (e.g., a count, a percentage, a trend) |
| `derivation` | If derived, the formula or logic (e.g., `count(open_cases where severity == 'S1')`) |

**Extraction rules:**

- **Dates** — extract both the value and the format used (ISO 8601, "May 26, 2026", "last Tuesday"). The placeholder should preserve the format in its name: `{{renewal_date_iso}}` vs. `{{renewal_date_display}}`.
- **Counts and aggregates** — always mark as `is_derived: true` and specify the derivation. Don't just extract "4 open cases"; extract the count and note it derives from filtering case data.
- **Enums and statuses** — note the full set of possible values if inferrable (e.g., severity: S1/S2/S3/S4).
- **Nested data** — for repeated structures (e.g., a table of cases), extract the template for one row and note the iteration pattern. Use array notation: `{{cases[].case_id}}`, `{{cases[].severity}}`.
- **Static text vs. data** — some text is structural (headings, boilerplate, transition sentences) and should stay in the prompt template verbatim. Only extract things that would change with fresh data.
- **Prose that embeds data** — sentences like "Goldman Sachs has 4 open S1 cases" contain both static framing and data. Extract the data points and leave the framing in the template: "{{account_name}} has {{s1_case_count}} open S1 cases."

### Step 4: Identify data sources

For each data point (or group of related data points), infer where the data came from. Use these source categories:

| Source category | Description | Typical access method |
|---|---|---|
| `case_mcp` | MongoDB case management system | `mdb_case_assistant` MCP tools: `get_case`, `list_account_cases`, `search`, `get_case_comments` |
| `salesforce` | Salesforce CRM (account, opportunity, contract) | Salesforce API, scraping, or export |
| `atlas_metrics` | MongoDB Atlas cluster metrics, billing, usage | Atlas Admin API, Atlas UI export |
| `tableau_export` | Tableau dashboard exports (CSV, screenshot) | Tableau Server API, manual download |
| `glean_search` | Enterprise search results | Glean MCP tools: `search`, `chat`, `read_document` |
| `corpus` | Local corpus store (IndexedDB + server dual-write) | Dashboard corpus APIs, `corpus-store/` modules |
| `monitoring` | Monitoring/alerting systems (Datadog, PagerDuty, Atlas alerts) | API or dashboard export |
| `calendar` | Calendar events, meeting schedules | Google Calendar API, manual |
| `manual_entry` | Manually written by the document author | Human input (cannot be automated without prompting) |
| `derived` | Computed from other data points | Calculation, aggregation, or inference |
| `contract` | Contract/agreement documents | Document store, legal system |
| `kb_article` | Knowledge base articles | KB search, Glean, internal wiki |
| `engineering` | Engineering tickets, PRs, design docs | Jira, GitHub, internal tools |
| `slack` | Slack messages, threads | Slack API, Slack export |
| `unknown` | Cannot determine the source | Flag for human review |

For each data point, record the source inference:

```json
{
  "placeholder": "open_case_count",
  "source": {
    "category": "case_mcp",
    "tool": "mdb_case_assistant.list_account_cases",
    "params": { "account_name": "{{account_name}}" },
    "field": "cases.filter(c => c.status !== 'Resolved').length",
    "confidence": "high",
    "reasoning": "Case counts in account reviews always come from the active case portfolio"
  }
}
```

**Confidence levels:**

- **high** — the source is unambiguous (case IDs come from the case system, Atlas cluster names come from Atlas)
- **medium** — the source is likely but could come from multiple places (a date could be from a contract or from CRM)
- **low** — educated guess; flag for human confirmation
- **unknown** — no reasonable inference possible

**When multiple sources are plausible**, list them in preference order and note the ambiguity.

### Step 5: Generate the reconstruction prompt

Build the `<name>.reconstruction-prompt.md` file. This is a complete, self-contained prompt that an LLM can execute to produce an equivalent document given fresh data values.

**Prompt structure:**

```markdown
# Reconstruction Prompt: <Document Title>

## Metadata
- **Original document:** <file path>
- **Document type:** <type>
- **Generated:** <timestamp>
- **Reporting period:** <if applicable>
- **Account:** <if applicable>

## Instructions

You are generating a <document type> for <context>. Follow the structure,
tone, and formatting conventions described below exactly. Replace every
`{{placeholder}}` with the corresponding value from the data manifest or
the fresh data provided.

### Tone and voice
<Describe the voice: formal/conversational, first-person/third-person,
technical level, audience>

### Formatting conventions
<Describe: heading levels, date formats, how case IDs are referenced,
table styles, bullet patterns, bold/italic usage>

### Section ordering
<Numbered list of sections in the exact order they appear>

## Document template

<The full document text with all data points replaced by {{placeholders}}.
Static text, transitions, and boilerplate remain verbatim.
Repeated structures use {{#each}} or equivalent iteration markers.>

## Data dependencies

| Placeholder | Source | Required | Notes |
|---|---|---|---|
| {{placeholder_name}} | <source description> | yes/no | <any caveats> |

## Generation checklist

- [ ] All {{placeholders}} resolved — no raw template markers in output
- [ ] Derived values recomputed from fresh data, not carried from the manifest
- [ ] Date formats match the conventions described above
- [ ] Tone matches: <voice description>
- [ ] Section count matches: <N> sections
- [ ] Tables populated with current data, not stale values
- [ ] No AI-isms (delve, leverage, robust, paradigm, seamless)
```

**Prompt writing rules:**

- The prompt must be **self-contained**. A reader with no knowledge of the original document should be able to follow it.
- Use `{{placeholder}}` syntax consistently. For arrays, use `{{#each cases}}...{{/each}}` or equivalent.
- Preserve the original document's voice and tone. Describe it explicitly in the Instructions section.
- Include conditional sections where the original document has them (e.g., "Include the Escalation section only if there are open escalations").
- The generation checklist acts as a self-validation contract.

### Step 6: Generate the data manifest

Build the `<name>.data-manifest.json` file. This is the machine-readable companion to the reconstruction prompt.

**Schema:**

```json
{
  "$schema": "document-deconstructor/data-manifest/v1",
  "document": {
    "title": "<document title>",
    "type": "<document type>",
    "original_path": "<file path>",
    "deconstructed_at": "<ISO timestamp>",
    "account": "<account name or null>",
    "reporting_period": "<period or null>"
  },
  "placeholders": {
    "placeholder_name": {
      "data_type": "string|number|date|boolean|enum|url|identifier|list",
      "section": "<section heading where it appears>",
      "context": "<the sentence or label surrounding the value>",
      "last_value": "<the value from the original document>",
      "is_derived": false,
      "derivation": null,
      "source": {
        "category": "<source category from the table>",
        "tool": "<MCP tool name, API endpoint, or 'manual'>",
        "params": {},
        "field": "<field path or extraction expression>",
        "confidence": "high|medium|low|unknown",
        "reasoning": "<why this source was inferred>"
      },
      "alternatives": [],
      "format": "<display format if relevant, e.g., 'ISO 8601', 'comma-separated', '$X.XM'>",
      "enum_values": null,
      "required": true,
      "notes": null
    }
  },
  "arrays": {
    "array_name": {
      "description": "<what this array represents>",
      "section": "<section heading>",
      "source": {
        "category": "<source category>",
        "tool": "<tool name>",
        "params": {},
        "field": "<field path to the array>"
      },
      "item_schema": {
        "field_name": {
          "data_type": "string",
          "field": "<path within each item>",
          "format": null
        }
      },
      "last_count": 0,
      "last_items": []
    }
  },
  "source_summary": {
    "source_category": {
      "placeholder_count": 0,
      "tools_used": [],
      "automation_feasibility": "full|partial|manual",
      "notes": null
    }
  },
  "automation_score": {
    "total_placeholders": 0,
    "automatable": 0,
    "semi_automatable": 0,
    "manual_only": 0,
    "score_percent": 0,
    "notes": "<interpretation of the score>"
  }
}
```

**Manifest rules:**

- Every `{{placeholder}}` in the reconstruction prompt must have a corresponding entry in the manifest.
- Every array iteration block must have a corresponding entry in `arrays`.
- The `source_summary` aggregates source usage so the user can see at a glance which data sources are needed.
- The `automation_score` quantifies how much of the document can be regenerated without human input.

### Step 7: Output and validate

1. Write both files to the same directory as the original document (or to a user-specified location).
2. Cross-validate (must pass before delivering):
   - Every `{{placeholder}}` in the prompt has a manifest entry
   - Every manifest entry has a corresponding `{{placeholder}}` in the prompt
   - No orphaned placeholders in either direction
   - `automation_score.total_placeholders` matches the actual count
   - All `{{#each}}` blocks have corresponding `arrays` entries in the manifest
   - **If validation fails:** fix the mismatch (add missing entries, remove orphans, correct counts), then re-validate. Do not deliver artifacts that fail cross-validation.
3. Report a summary to the user:
   - Total placeholders extracted
   - Automation score (what percentage can be auto-populated)
   - Source distribution (how many data points from each source)
   - Manual-entry items that need human input on each regeneration
   - Any `unknown` or `low` confidence source inferences that need review

## Worked example: Account Review snippet

**Original text:**
> Goldman Sachs currently has 4 open support cases, including 1 S1 (Case 01567493 - Rewards PROD Failover) and 2 S2 cases. Average MTTR over the past 90 days is 4.2 hours. The Premium Support contract renews on 2026-09-30 with an ARR of $2.1M.

**Extracted placeholders:**

| Placeholder | Value | Type | Source |
|---|---|---|---|
| `{{account_name}}` | Goldman Sachs | string | manual_entry / CRM |
| `{{open_case_count}}` | 4 | number (derived) | case_mcp: list_account_cases |
| `{{s1_case_count}}` | 1 | number (derived) | case_mcp: list_account_cases filtered by severity |
| `{{s1_case_id}}` | 01567493 | identifier | case_mcp: list_account_cases |
| `{{s1_case_title}}` | Rewards PROD Failover | string | case_mcp: get_case |
| `{{s2_case_count}}` | 2 | number (derived) | case_mcp: list_account_cases filtered by severity |
| `{{mttr_90d_hours}}` | 4.2 | number | derived: avg(resolution_time) over 90d window |
| `{{renewal_date}}` | 2026-09-30 | date | contract / salesforce |
| `{{arr_display}}` | $2.1M | string | salesforce / contract |

**Reconstruction prompt fragment:**
```
{{account_name}} currently has {{open_case_count}} open support cases,
including {{s1_case_count}} S1 (Case {{s1_case_id}} - {{s1_case_title}})
and {{s2_case_count}} S2 cases. Average MTTR over the past 90 days is
{{mttr_90d_hours}} hours. The Premium Support contract renews on
{{renewal_date}} with an ARR of {{arr_display}}.
```

**Data manifest fragment:**
```json
{
  "open_case_count": {
    "data_type": "number",
    "section": "Case Portfolio",
    "context": "currently has {{open_case_count}} open support cases",
    "last_value": 4,
    "is_derived": true,
    "derivation": "count(cases.filter(c => c.status !== 'Resolved'))",
    "source": {
      "category": "case_mcp",
      "tool": "mdb_case_assistant.list_account_cases",
      "params": { "account_name": "Goldman Sachs" },
      "field": "cases.filter(c => c.status !== 'Resolved').length",
      "confidence": "high",
      "reasoning": "Open case count is always derived from the active case portfolio"
    },
    "required": true
  }
}
```

## Source inference heuristics

Use these heuristics to infer data sources when the document does not explicitly cite them:

| Pattern in document | Likely source |
|---|---|
| Case IDs (8-digit numbers, e.g., 01567493) | `case_mcp` — `get_case` or `list_account_cases` |
| Case severity (S1/S2/S3/S4) and status | `case_mcp` — case record fields |
| Case comments, timeline entries | `case_mcp` — `get_case_comments` |
| Atlas cluster names, project names | `atlas_metrics` — Atlas Admin API |
| ARR, contract value, renewal dates | `salesforce` — opportunity/contract records |
| Account name, org ID, TAM assignment | `salesforce` or `corpus` — account records |
| SLA metrics (MTTR, MTTF, response time) | `derived` from case data, or `tableau_export` |
| Architecture diagrams, topology descriptions | `manual_entry` or `glean_search` |
| Meeting dates, cadence, attendees | `calendar` |
| Action items from prior meetings | `manual_entry` or `corpus` — meeting notes |
| KB article references | `kb_article` — `glean_search` |
| Jira ticket IDs (HELP-*, SERVER-*, etc.) | `engineering` — Jira API |
| Slack thread references | `slack` — Slack API |
| "As of <date>" temporal anchors | `derived` — reporting period, not a separate data point |
| Executive names, titles, org structure | `salesforce` or `manual_entry` |
| Version numbers (server, driver, Atlas) | `atlas_metrics` or `case_mcp` (case context) |
| Percentages, trends, deltas | `derived` — always specify the derivation formula |
| Tableau dashboard screenshots or chart data | `tableau_export` |

## Handling ambiguity

When the source of a data point is unclear:

1. **List all plausible sources** in the `alternatives` array of the manifest entry.
2. **Set confidence to `low` or `medium`** and explain the reasoning.
3. **Flag it in the output summary** so the user can confirm.
4. **Prefer the more automatable source** when two sources are equally likely. If the data could come from a manual note or from the case MCP, prefer the case MCP because it can be automated.

## Cross-skill integration

This skill works best when combined with:

- **`tam-expertise`** — for domain knowledge about account review structures, QBR formats, and TAM deliverables
- **`tam-reference`** — for MongoDB Premium Services operating procedures and terminology
- **`case-tracker`** — for understanding case data schemas and TS Tools API patterns
- **`operator-report-generator`** — for report template patterns, SBAR/BLUF conventions, and quality validation
- **`writing-expert`** — for voice/tone analysis and anti-AI-ism detection during prompt writing
- **`document-critique`** — to validate the reconstruction prompt itself before use
- **`case-mcp-server-guide`** — for correct MCP tool names and parameters when mapping case data sources
- **`tstools-reference`** — for TS Tools Support API endpoint details

## Anti-patterns to avoid

| Anti-pattern | Why it fails | Remedy |
|---|---|---|
| **Over-extracting static text** | Template becomes unreadable; every word is a placeholder | Only extract values that change between regenerations |
| **Under-extracting embedded data** | Prose sentences contain data points that get baked in as static text | Parse every sentence for facts that could change |
| **Vague source attribution** | `source: "database"` tells the user nothing actionable | Name the specific tool, endpoint, field path |
| **Missing derivation formulas** | Derived values without formulas can't be recomputed | Always specify the formula: `count(cases where severity == 'S1')` |
| **Ignoring conditional sections** | Template includes sections that should only appear when data warrants | Use conditional markers: `{{#if escalations.length > 0}}` |
| **Placeholder name collisions** | Two different data points with the same placeholder name | Namespace by section: `{{exec_summary_case_count}}` vs. `{{detail_case_count}}` |
| **Forgetting array structures** | A table of cases gets one placeholder instead of an iterable | Use `{{#each cases}}` with per-item field placeholders |
| **Prompt that requires the original** | The reconstruction prompt says "follow the original document's style" without describing it | The prompt must be fully self-contained |
| **Stale manifest values treated as current** | User regenerates using `last_value` instead of fetching fresh data | The `last_value` field is for reference only; the prompt must instruct fresh fetching |
| **No validation step** | Orphaned placeholders in prompt or manifest go undetected | Always cross-validate both artifacts before delivering |

## Common mistakes

- **Treating `last_value` as live data.** The manifest stores the value from the original document for reference. The reconstruction prompt must instruct fresh fetching, not reuse of stale values.
- **Extracting boilerplate as placeholders.** Section headings, transition phrases, and standard disclaimer text should stay verbatim in the template. Only values that change between regenerations become placeholders.
- **Flat placeholder lists for nested data.** A table with 5 case rows needs `{{#each cases}}` with per-row fields, not 5 x N individual placeholders like `{{case_1_id}}`, `{{case_2_id}}`.
- **Omitting the tone/voice section.** Without explicit voice guidance, the LLM regenerating the document will use its default voice, which rarely matches the original author's style.
- **Skipping the cross-validation step.** Orphaned placeholders (present in one artifact but not the other) are the single most common defect. Always run Step 7 validation.
- **Source confidence inflation.** Marking a source as `high` confidence when it is actually an educated guess. When in doubt, use `medium` and flag for review.

## Output file naming

Given an input document at `path/to/document-name.md`:

- Reconstruction prompt: `path/to/document-name.reconstruction-prompt.md`
- Data manifest: `path/to/document-name.data-manifest.json`

For documents without a file extension or with non-standard names, sanitize the name to `kebab-case` and use the same directory.

If the user specifies an output directory, use that instead.
