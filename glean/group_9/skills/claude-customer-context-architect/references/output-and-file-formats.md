## Output

- **Files emitted per run:**
  - Context file: `<CUSTOMER_SLUG>-context-<RUN_DATE>.md` at `<Drive>/<CUSTOMER_FOLDER>/contexts/`.
  - Per-case analysis files: `case-<CASE_NUMBER>-analysis-<RUN_DATE>.md` at `<Drive>/<CUSTOMER_FOLDER>/contexts/case-analyses/`.
  - Per-section reference files: `<section>-<RUN_DATE>.md` at `<Drive>/<CUSTOMER_FOLDER>/contexts/references/`. One file per customer-specific section (cases, aha, monday, monday-sync, entity-map, slack, slack-intervention, sentiment, chronology, tableau, todo). Each holds the raw notes and analysis that back the summarized section in the main context file.
  - Definitions sidecar: `<Drive>/<CUSTOMER_FOLDER>/definitions.md`. Overwrite **only if** the new content is a strict superset of the prior file; otherwise emit `definitions-<RUN_DATE>.md` alongside and leave the prior file untouched.
  - Accessibility sidecar: `context-v3-accessible.md` at `<Drive>/<CUSTOMER_FOLDER>/contexts/`. Screen-reader/audio-optimized rendering of the context file per `## Markdown structure: accessibility sidecar`. Emit fresh each run alongside the dated context file; overwrite the prior accessible sidecar (it is not dated). Does not count against the context-file length ceiling.
  - Create all directories if they do not exist. Prior files remain untouched.
- **Link rewriting:** case numbers → `https://hub.corp.mongodb.com/case/<CASE_NUMBER>` (verify against skill registry; fallback to plain number + warning); document titles → clickable links; `Analysis file` cell → relative link to newest per-case file or `—`.
- **Anchors:** GFM slugs (`## Foo bar` → `#foo-bar`).
- **Empty cells:** `—`.
- **Length ceilings:** context file ≤ 12,000 words (raised from 6,000 to accommodate the strategic/technical/relationship/knowledge intelligence bands and end-matter sections added in v3.0.0); appendix ≤ 1,500 words; corpus manifest ≤ 200 rows; item index ≤ 500 rows; per-case file ≤ 3,000 words; corpus summary ≤ 400 words; Monday sync manifest ≤ 100 rows; Slack intervention opportunities ≤ 10; TODO ≤ 200 rows. Per-section row caps are given inline in `## Markdown structure: context file`. The accessibility sidecar is a separate file and does not count against the context-file ceiling. When exceeded, drop lowest-confidence first (or lowest-tier for TODO) with `…N omitted`.

## Markdown structure: context file

Emit in this order. All anchor references use GFM slugs. Failed sections still emit their heading with `[Section Failed — <reason>]` as the body.

1. `## Table of contents`: GFM anchor links to every subsequent section; skip if the document is under 500 words.
2. `## Executive summary`: about the **account** (business + technical state, top risks). Prepend a `30-second` variant (50–75 words: status + single top risk/opportunity) and a `2-minute` variant (200–300 words: status, sentiment trend, top 3 items); the existing ≤ 250-word summary is the `5-minute` version. Include a renewal-risk line: `Renewal risk: <score>/100 (<band>) — Drivers: <top 2–3 factors>. Mitigations: <1–2 imperative actions>.` Append `(partial — missing: <inputs>)` when computed from incomplete inputs.

_Strategic & business intelligence band:_
3. `## Competitive intelligence`: competitor mentions from meetings, Slack, cases, SFDC opportunity competitor fields. Table: Competitor, Last Mentioned, Context, Sentiment, Source; plus a `Key themes` subsection (won-on / at-risk). Cap 30 rows.
4. `## Executive stakeholders`: VP+/C-suite contacts and engagement state. Table: Name, Title, Last Touchpoint, Topics, Sentiment, Engagement Trend, Notes; plus a Risks/Strengths callout list. Cap 20 rows.
5. `## Internal champions & detractors`: per-person relationship classification. Three sub-tables — Champions (Name, Role, Evidence, Engagement Strategy) · Neutrals (Name, Role, Notes) · Skeptics (Name, Role, Evidence, Mitigation Strategy); plus a relationship-strategy line. Cap 15 rows per sub-table. Change column.
6. `## Customer maturity assessment`: five-dimension scoring. Table: Dimension, Score (1–5 or —), Evidence, Trend vs Prior Quarter; plus `Overall maturity` and `Recommended focus` lines. Fixed 5 rows.
7. `## Expansion opportunities`: growth/upsell signals. Table: Opportunity, Estimated Value, Timeline, Status, Blocker. Cap 25 rows.

_Technical & operational intelligence band:_
8. `## Technology stack`: categorized non-MongoDB technology surface. Table: Category, Technologies, Notes; plus an `Integration opportunities` list. Cap 50 rows.
9. `## Performance & scale trends`: quantitative trends. Table: Metric, Current, Prior Quarter, Change, YoY Growth, Anomalies; plus `Trends` and `Forecasts` lines. Cap 30 rows.
10. `## MongoDB feature adoption`: Atlas feature availability vs actual use. Table: Feature, Available, In Use, Adoption Rate, Opportunity; plus an `Underutilized features` list. Cap 30 rows. Change column.
11. `## Migration timeline` (conditional — omit heading entirely when no migration detected): header line `Type · Target go-live · Rollback plan`; stage table: Stage, Planned Date, Actual Date, Status, Blockers; plus a `Δ timeline` note. Cap 30 rows.
12. `## Anti-pattern recurrence analysis`: table: Anti-Pattern, Occurrences (180d), Cases, Recommendation; plus a high-priority action callout. Cap 30 rows.
13. `## SLA/SLO commitments`: table: Commitment, Source, Status, Performance, Notes; plus an `Upcoming deadlines` list. Cap 30 rows.

_Relationship & communication intelligence band:_
14. `## Communication velocity`: MongoDB↔customer responsiveness. Table: Direction, Avg Response Time, Trend, Issues; plus going-dark + asymmetric-urgency notes + assessment. Cap 20 rows.
15. `## Commitment tracker`: promises extracted from meetings/cases. Table: Commitment, Owner, Due Date, Status, Days Overdue, Impact; plus promise-keeping rate + action-required line. Cap 50 rows. Change column.
16. `## Meeting participant trends`: per-attendee attendance over the last 6 customer-facing meetings. Table: Attendee, Role, Meetings Attended (Last 6), Trend, Notes; plus pattern-insight callouts. Cap 30 rows. Change column.

_Knowledge & learning band:_
17. `## Reusable solutions library`: resolved-case fixes for reuse; vector-similarity auto-suggestion is enabled from these rows, not run at generation time. Table: Problem Type, Case, Solution Summary, Reusability, Artifacts. Cap 100 rows.
18. `## Customer-specific gotchas`: unusual constraints; flag `>2x` rows `[Repeatedly forgotten]`. Table: Constraint, Impact, Source, First Mentioned, Last Mentioned. Cap 50 rows. Change column.
19. `## Training needs assessment`: gaps mapped to MongoDB training offerings. Table: Gap Area, Evidence, Recommended Training, Est. Impact, Delivered, Outcome; plus a `Priority` line. Cap 20 rows. Change column.
20. `## Value delivered (for QBR narrative)`: MongoDB wins for QBR use. Table: Achievement, Date, Impact, Customer Acknowledgment, Suggested QBR Story; plus a one-line QBR value narrative. Cap 25 rows. Change column.

_Operational spine:_
21. `## Corpus summary`: ≤ 400 words about the **document set** (types, themes, time coverage, gaps).
22. `## Entity map`: table: Kind, Value, Notes, Change. Kinds: `org-id`, `project-id`, `cluster`, `alias`, `owner-tam`, `owner-csm`, `owner-ntse`, `owner-ir`.
23. `## Definitions`: table: Term, Definition, Source, Change. Also emitted to `definitions.md` per Output.
24. `## High-confidence facts`: `[High] <fact> — <link>`.
25. `## Open questions`: materially important only.
26. `## Major incidents and themes`: chronological.
27. `## Cases`: table: Case #, Subject, Owner, Sev, Created, Last update, Status, Summary (≤ 40 words), Root cause, Resolution, Next steps, Analysis file, Change.
28. `## AHA feature requests`: table: ID, Title, Status, Value, Customer requirements, Source, Change.
29. `## Monday.com`: table: Item ID, Group, Owner, Status, Timeline, Notes, Staleness, Source, Change.
30. `## Monday sync manifest`: table: Item ID, Field, Old value, New value, Direction (`pull`|`push`|`propose`), Outcome (`written`|`failed; <reason>`|`—`). Empty body when the sync unit didn't run.
31. `## Slack and document context`: grouped by topic.
32. `## Slack intervention opportunities`: prompt rows per `## Slack intervention audit`; cap 10.
33. `## TODO`: prioritized aggregated TODOs per `## Unified TODO compilation`; add an `Auto-Nag` column (`—` | `7d sent <DATE>` | `14d sent <DATE>` | `21d escalated <DATE>`) populated by `## Workflow automation` (5.3). Cap 200 rows.
34. `## Sentiment and contacts`: trend + contact list with last-interaction date.
35. `## Chronology`: timeline of milestones.
36. `## Visual timeline`: single Mermaid ```mermaid gantt``` block rendering `## Chronology` events (cases, milestones, sentiment shifts, team changes, incidents), color-coded by type. `[Section Failed]` body if Chronology failed.
37. `## Knowledge health assessment`: corpus staleness + tribal-knowledge gaps (relates to `## Corpus manifest`). Two sub-tables — `### Stale documents` (Document, Last Accessed, Last Modified, Recommendation) · `### Tribal knowledge` (Topic, Evidence, Recommendation); plus an `### Actions` list. Cap 50 rows per sub-table.
38. `## Low-confidence data`: `[Low] <signal> — <link>`.
39. `## Data conflicts`.
40. `## Historical`: dropped-but-preserved entries with last-seen date and Change status.
41. `## Missing sources`.
42. `## Ambiguities blocking completion`: only if unresolved after Clarification protocol.
43. `## Corpus manifest`: table: Document, Type, Date, Source path, 1-line summary. Cap 200.
44. `## Item index`: table: Item, Type, First seen, Last seen, Sources, Change. Timestamps anchor to source-document dates so re-runs are stable. Cap 500.
45. `## Appendix`: deduped raw bullets, ≤ 1,500 words. Distinct from item index (appendix = raw statements; index = provenance).

_End matter (internal, not customer-facing):_
46. `## Context drift analysis`: monthly delta vs 30/60/90-day-prior context files. Table: Dimension, Prior (30d ago), Current, Change, Significance; plus key-drift callouts. `[Section Failed — insufficient context-file history]` when no 30-day-prior file exists. Cap 15 rows.
47. `## Workflow automations`: internal status block — per-automation Status (`active` | `disabled — <reason>`) + Trigger for case pre-analysis (5.1), followed by a table: Case, Date, Suggestion, Confidence, Outcome (recent pre-analyses; cap 10). 5.2 (meeting prep) and 5.5 (case routing) emit nothing here.
48. `## Context Q&A capability`: internal status block — Status, Index name (`ctx-<CUSTOMER_SLUG>`), Query interface (`query_customer_context(question, account_id)`), 3 sample queries.
49. `## Confidence legend`: inline tags — High · Medium · Low · Needs Verification — with the source/recency criteria from `## Definitions`.

## Markdown structure: per-case analysis file

1. `## Case summary`: 3–5 sentences.
2. `## Facts established`: bullets with links.
3. `## Technical analysis`: expectation-vs-reality gaps explicit.
4. `## Unknowns and open questions`.
5. `## Blockers`.
6. `## Recommended next steps`: numbered, owner named.
7. `## Escalation guidance`: evidence-based.
8. `## Customer response draft` (optional; skip if Closed).
9. `## Additional insights`.
10. `## Sources`.
11. `## Practical troubleshooting plan`.
12. `## Customer validation commands` (optional).
13. `## Meeting summary` (optional).
14. `## Meeting flow` (optional).
15. `## Decisions made` (optional).
16. `## Action items` (optional).
17. `## Meeting transcript` (optional).
18. `## Missing sources`.

## Markdown structure: definitions.md sidecar

Standalone file at `<Drive>/<CUSTOMER_FOLDER>/definitions.md`. Mirrors the context file's `## Definitions` table with a header:

```
# Customer Definitions — <CUSTOMER_NAME>

Last updated: <RUN_DATE>

| Term | Definition | Source | Change |
```

Overwrite policy per Output block.

## Markdown structure: accessibility sidecar

Standalone file `context-v3-accessible.md` at `<Drive>/<CUSTOMER_FOLDER>/contexts/`. Screen-reader-first, audio-narration-ready. Emit fresh every run alongside the dated context file; never overwrite the dated context file. `## Security constraint` still applies.

1. `# Accessible Customer Context — <CUSTOMER_NAME>` + `Last updated: <RUN_DATE>`.
2. `## Table of contents`: GFM-slug anchor links to every section.
3. Strict single-increment heading hierarchy (no skipped levels).
4. Per section: a plain-language summary in short sentences; expand acronyms on first use with a pronunciation guide (e.g., "SLA (ess-ell-ay)").
5. Every table preceded by a one-line alt-text description of what it conveys.
6. High-contrast status markers as text labels, not color alone: `[BLOCKER]`, `[HIGH RISK]`, `[OK]`.
7. `## Pronunciation guide`: acronym → spoken form list.

## Version history

- **v3.2.1 (2026-07-03):** Corrected the output directory name from `context-file/` to `contexts/` (the real TAM shared-drive convention) and introduced `CUSTOMER_FOLDER` as a distinct derived variable — the resolved Drive engagement folder name (often a short internal code, e.g. `GS` for Goldman Sachs, N.A.), never the SFDC canonical `CUSTOMER_NAME` or `CUSTOMER_SLUG`. Applied throughout all path templates (context file, per-case files, references, definitions.md, accessibility sidecar).
- **v3.2.0 (2026-07-03):** Replaced the AHA feature-request source with the real retrieval procedure — a corp-SSO-gated Aha! shared report split alphabetically by account name, with the known G–Z report URL registered and the A–F range flagged for one-time `/dr` discovery. Documented the SSO-redirect access constraint (credential-less fetch fails) and the not-pre-filtered-by-customer full-text search step; kept the prior Drive-search behavior as an explicit fallback rather than replacing it outright.
- **v3.1.1 (2026-07-03):** Corrected the Glean connector reference — `glean_default` was a broken/disabled server; replaced throughout with the working `claude.ai Glean (via MCP)` connector and its real tools (`search`, `read_document`, `chat`), with `owner:"<name>"`/`from:"<name>"`/date filters cited for the document-discovery and QBR-owner-search use cases. Dropped the invented in-skill `authenticate` step (this connector has none) in favor of a connected/not-connected precondition, and added a permission-aware-results rule (an empty Glean result is not a failure).
- **v3.1.0 (2026-07-03):** Wired Glean explicitly into `## Sources and fallback policy` (item 4) as the fallback query layer for account context no native source can retrieve — calendar entries (item 9, now Glean-backed), document discovery (`## MongoDB team discovery` failure handling), and meeting-attendee history for Relationship intelligence 3.4. Added a fallback-provenance rule (`Source: Glean (fallback)`, capped at Medium confidence, never overrides a canonical source).
- **v3.0.1 (2026-07-03):** skill-optimizer pass — extracted the six detailed v3.0.0 execution units (strategic/technical/relationship/knowledge intelligence, workflow automation, presentation generation) to `references/intelligence-and-automation-units.md` for length budget, leaving at-a-glance summaries + pointers in the execution model; body reduced ~16k→~12.6k tokens. No behavior change.
- **v3.0.0 (2026-07-02):** Added 30 intelligence & automation enhancements across 6 categories, each integrated into both the execution model and the context-file markdown structure; raised the context-file length ceiling from 6,000 to 12,000 words; added inline confidence tagging (High/Medium/Low/Needs Verification), a Confidence legend, and an accessibility sidecar (`context-v3-accessible.md`).
  - **Strategic:** competitive intelligence, renewal-risk scoring, expansion opportunities, executive stakeholder map, customer maturity model.
  - **Technical:** technology stack, migration timeline (conditional), performance & scale trends, anti-pattern recurrence, SLA/SLO registry.
  - **Relationship:** communication velocity, unacknowledged value delivered, commitment tracker, meeting participant patterns, champions/detractors.
  - **Knowledge:** reusable solutions library, customer-specific gotchas, knowledge decay, training needs, feature adoption.
  - **Workflow:** case pre-analysis, meeting prep, stale-TODO auto-nag, context-drift detection, intelligent case routing.
  - **Presentation:** executive-summary variants, visual timeline, context Q&A, confidence scoring, accessibility mode.

