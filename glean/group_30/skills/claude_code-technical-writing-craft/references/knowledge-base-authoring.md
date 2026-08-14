<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `knowledge-base-authoring` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: knowledge-base-authoring
description: "Authoring craft for knowledge-base articles: writing for search discoverability (answer-the-question-in-the-title rule), one-question-per-article discipline, KB vs runbook vs FAQ distinctions, content lifecycle (draft → published → reviewed → archived), KCS (Knowledge-Centered Service) reuse-is-review patterns, deflection-metric-aware writing, escalation triggers ('if this didn't help…'), evergreen vs time-bound content, single-source-of-truth problems, KB taxonomy and faceted-search design, cross-linking patterns, MadCap Flare single-sourcing, and the Zendesk/Confluence/Notion/Salesforce-KB platform best practices. TRIGGER: user asks to write, structure, audit, or improve a knowledge-base article; user mentions 'KB article', 'help article', 'self-service content', 'KB taxonomy', 'KCS', 'deflection rate', 'support article', 'Zendesk article', 'Confluence KB', 'help center'; user wants to decide between KB, runbook, and FAQ; user wants a content lifecycle, review cadence, or archival policy; user wants to write a title that ranks in search. SKIP: maintaining and auditing an existing KB at scale for outdated content (use doc-archaeology); writing an internal-only runbook for engineers (use runbook-craft); writing user-facing UI microcopy (use microcopy-and-ui-writing); writing post-mortems or RCA (use postmortem-writing); writing executive or report prose (use writing-expert)."
origin: local
version: "1.0.0"
updated: "2026-05-29"
keywords:
  - knowledge base
  - KB article
  - help center
  - self-service content
  - KCS
  - knowledge-centered service
  - deflection rate
  - search discoverability
  - one question per article
  - KB taxonomy
  - faceted search
  - Zendesk
  - Confluence
  - Notion
  - Salesforce Knowledge
  - MadCap Flare
  - single-sourcing
  - evergreen content
  - article lifecycle
  - escalation trigger
tags:
  - writing
  - knowledge-base
  - support-content
  - self-service
  - documentation
  - taxonomy
whenToUse:
  - User asks to write a KB article, help-center article, or self-service article
  - User asks to audit or improve an existing KB article for search discoverability or clarity
  - User asks to choose between writing a KB article, a runbook, or a FAQ entry
  - User asks about KB taxonomy, faceted search, or article categorization
  - User asks about article lifecycle, review cadence, deprecation, or archival
  - User asks about KCS (Knowledge-Centered Service) practices: capture, structure, reuse, improve
  - User asks how to write a title that gets found in search
  - User asks about deflection metrics, escalation triggers, or "if this didn't help…" patterns
  - User is migrating content into Zendesk, Confluence Service Desk, Notion, Salesforce Knowledge, or MadCap Flare
  - User asks about evergreen vs time-bound content or "single source of truth" design
whenNotToUse:
  - User wants to audit an existing KB at scale for outdated content (use doc-archaeology)
  - User wants to write a runbook for engineers operating a system (use runbook-craft)
  - User wants to write product UI microcopy (use microcopy-and-ui-writing)
  - User wants post-mortem or RCA writing (use postmortem-writing)
  - User wants executive prose or business reports (use writing-expert)
  - User wants legal-adjacent contract clauses or disclosures (use legal-adjacent-writing)
related_skills:
  - doc-archaeology
  - runbook-craft
  - microcopy-and-ui-writing
  - writing-expert
  - support-ticket-writing
  - technical-writing-craft
  - plain-language
---

# Knowledge Base Authoring

Reference for writing knowledge-base articles, help-center content, and self-service support documentation. Covers the article-level craft (title, structure, voice, escalation) and the system-level design (taxonomy, lifecycle, single-source-of-truth, deflection metrics). Partner to `doc-archaeology` (which audits existing KBs for staleness) and `runbook-craft` (which writes step-by-step operational procedures for engineers).

## When to use this skill

- Writing a new KB / help-center / self-service article
- Auditing or rewriting an existing article for search discoverability or clarity
- Choosing between KB article, runbook, FAQ, and inline-help microcopy
- Designing a KB taxonomy, faceted-search structure, or categorization scheme
- Setting up an article lifecycle (draft → review → publish → review → archive)
- Implementing or running KCS (Knowledge-Centered Service)
- Tuning titles for search and deflection rate
- Designing escalation triggers (the "if this didn't help…" pattern)
- Migrating content between Zendesk, Confluence Service Desk, Notion, Salesforce Knowledge, MadCap Flare, or similar

## When NOT to use this skill

- Auditing an existing KB at scale for outdated content → `doc-archaeology`
- Writing a runbook (engineer-facing, action-focused, technical) → `runbook-craft`
- Writing UI microcopy (button labels, tooltips, form copy) → `microcopy-and-ui-writing`
- Writing post-mortems → `postmortem-writing`
- Writing executive prose → `writing-expert`

## The 9-point KB article test

Every KB article should pass these tests before publishing:

1. **The title is the question the user typed.** Not "Password Reset" — "How do I reset my password?" or "Reset your password." Search words on the page; titles that match queries rank higher.
2. **The article answers exactly one question.** If the article covers two questions, split it.
3. **The first paragraph is the answer.** A user who reads only the first paragraph should leave with the answer. Details follow.
4. **The reader can self-identify in the first 50 words.** "This article is for [Atlas Free Tier users / Enterprise customers / administrators]." Disqualify out-of-scope readers fast.
5. **Steps are numbered, prerequisites are explicit.** No buried "first you must…" in the middle of the steps.
6. **There is an escalation path.** "If this didn't solve the problem, [contact support / try X / see Y]." Never end a KB article in a dead end.
7. **The article is dated and versioned.** "Last updated: 2026-05-29." Stale-looking content erodes trust even when correct.
8. **There is at least one cross-link.** Either to a related article, a runbook, or the next logical step. Articles in isolation feel orphaned.
9. **The article has at least one search-friendly keyword cluster the title doesn't cover.** "Forgot password" "lost password" "can't sign in" are different queries that lead to the same article — cover them in the body or via labels.

## Core concepts

### 1. KB vs runbook vs FAQ vs inline help

The four self-service content types serve different audiences and answer different questions. Choose deliberately:

| Type | Audience | Answers | Length | Example |
|---|---|---|---|---|
| **KB article** | End user, customer, new employee | "How do I do X?" "What does Y mean?" "Why is Z happening?" | 300–1,500 words | "How do I reset my password?" |
| **Runbook** | Engineer, SRE, support agent under pressure | "What are the exact steps to remediate incident X?" | 200–2,000 words; step-focused | "Restart the cluster after OOM" |
| **FAQ** | Anyone with a quick question | "Short answer to common question." | 1–3 sentences per Q | "Do you support SSO? Yes — see…" |
| **Inline help / tooltip** | User in the product right now | "What does this field mean?" | 1–15 words | "Email address used for account recovery." |

The most common mistake: writing a runbook and publishing it in the customer KB. Customers don't care about your internal kubectl steps. Mirror: writing a customer KB article and stuffing it into engineering's runbook repo. Engineers want commands, not introductions.

Pick the format based on who is reading and what they're doing when they read it.

### 2. The "answer the question in the title" rule

Knowledge-base articles compete in two search engines: the platform's internal search (Zendesk, Confluence, Salesforce, Notion) and external search (Google). Both rank titles heavily.

Title patterns that work:

- **The question form.** "How do I export my data?" "Why am I getting a 429 error?" "What is a replica set?"
- **The action form.** "Export your data." "Resolve a 429 error." "Set up replica sets."
- **The error-or-symptom form.** "'Connection refused' when connecting to MongoDB Atlas." "Backup failing with E11000 error."

Title patterns that fail:

- **Topic-only.** "Data Export." (Not specific enough; doesn't match a user's mental model.)
- **Product-only.** "Atlas Search." (Discovery yes; article-level title no.)
- **Marketing voice.** "Effortlessly Manage Your Data Anywhere!" (Doesn't match what people search.)
- **Internal jargon.** "FTDC Configuration Issues" (when the user doesn't know what FTDC is).

The "shadow the search query" test: open your support ticketing system, find the 20 most common ticket titles for the topic, and align your KB article titles to those queries.

### 3. The "one article, one question" rule

A KB article should answer exactly one question. The reasons:

- **Search ranking.** An article that answers two questions ranks weaker for each than two focused articles would.
- **Reuse in KCS.** Support agents linking the article in a ticket want a single-purpose article. A multi-purpose article forces the agent to caveat: "see the second section."
- **Maintenance.** Updates to question A may not affect question B. Splitting reduces the blast radius of edits.
- **Deflection measurement.** You can't measure deflection per question if two questions live in one article.

The exception: composite articles that explicitly chain ("Setting up X involves three steps: 1, 2, 3"). The composite article is the entry point; each step is a separate article.

### 4. KCS (Knowledge-Centered Service) — the core idea

KCS is a methodology from the Consortium for Service Innovation. Core principles relevant to authoring:

- **Capture in the moment.** Create or improve a KB article during the support interaction, not as a separate task afterward.
- **Structure for reuse.** Every article should be findable, scannable, and link-able. A support agent should be able to paste a link into a ticket and trust the article will resolve.
- **Reuse is review.** Every time an article is reused, it implicitly proves usefulness. Demand-driven improvement: high-reuse articles get refined; never-reused articles get archived.
- **Improve in the flow of work.** When an agent finds an article incomplete or wrong, they fix it then and there — not as a separate ticket.
- **Demand-driven KB.** Don't write articles for problems no one has. Write what tickets ask for. Track reuse, not page count.

The KCS lifecycle of an article:

1. **Draft / Working.** Visible to authors only.
2. **Validated for internal use.** Reviewed by a peer or KCS coach; visible to support agents.
3. **Validated for external use.** Reviewed for customer-facing clarity, brand voice, no internal jargon; visible to customers.
4. **Mature.** High reuse, low edit rate.
5. **Archived.** No longer relevant or replaced by a better article.

### 5. Deflection rate and the metrics that follow

**Deflection rate** is the percentage of support inquiries that get resolved by self-service content without contacting a human. It is the single most-watched KB metric.

The deflection-aware writing moves:

- **Make the answer findable in 30 seconds.** Users who can't find an answer give up and contact support — counting as a deflection failure.
- **Make the answer skim-able.** Bullet points, bold key terms, numbered steps. Walls of prose drive users to the contact button.
- **Make the next-step explicit.** "If you completed these steps and the error persists, [contact support with the error code from step 4]." This is paradoxical but real: a clear escalation path increases deflection by giving users confidence to try the steps fully.
- **Front-load disqualification.** "This article is for Atlas Dedicated clusters. For Free Tier, see [link]." The wrong-audience user bounces fast and finds the right article.

Secondary metrics that pair with deflection:

- **Article usefulness rating.** Thumbs up/down at the article footer.
- **Searches with no clicks.** A query that shows results but nothing gets clicked indicates poor title/snippet match.
- **Searches with no results.** A content gap.
- **Article-to-ticket ratio.** For a given topic, how many tickets cite the article vs. how many ignore it.

### 6. Article lifecycle and review cadence

A KB article has a lifecycle. Without explicit lifecycle management, KBs decay: articles describe the 2020 UI, link to dead pages, or reference deprecated features.

A workable lifecycle:

| Stage | Trigger | Action |
|---|---|---|
| Draft | Author creates | Peer review; assign reviewer |
| Published — fresh | First publish | Indexed; available; reviewed at 90 days |
| Published — mature | Reviewed; high reuse | Reviewed annually |
| Published — at risk | Low usefulness rating, low reuse, contradicted by ticket data | Triaged: rewrite, merge, or archive |
| Deprecated | Feature removed or replaced | Banner: "This article describes a deprecated feature. See [new article]." |
| Archived | No longer accessible | Not in search; preserved as historical reference only |

**Review cadence by content type:**

| Content type | Cadence | Reason |
|---|---|---|
| UI walkthroughs with screenshots | 90 days | UIs change |
| API references | On every API change | Mismatch = bug report |
| Pricing/plans | On every pricing change | Wrong info costs revenue |
| Conceptual ("What is X?") | Annual | Stable but should refresh |
| Troubleshooting | Quarterly | New failure modes emerge |
| Compliance / regulatory | On every reg change | Stale = risk |

### 7. Single source of truth (SSoT)

The classic KB problem: the same information lives in three places — the customer help center, the internal wiki, and a sales-enablement deck — and they drift apart. Three approaches:

1. **MadCap Flare or similar single-sourcing.** Author content once in a topic-based system; publish to multiple outputs (PDF, online help, internal KB, customer KB). Variables and conditional tags differentiate audience versions of the same source.
2. **Linked references.** One article is canonical; others link to it with no duplication. Brittle: link rot, fragmented user experience.
3. **One canonical home, syndicated reads.** Publish in one place; embed via includes or API into others. Used by Notion (database views) and Confluence (excerpt-include).

The hard part is not the tooling — it is the discipline. Without a published rule ("API behavior lives in /docs/api; if you find it elsewhere, link, do not copy"), drift returns.

### 8. KB taxonomy and faceted search

A good KB has two organizational axes:

- **Hierarchical category.** Top-level concept buckets (Getting Started, Account, Billing, API, Troubleshooting, Security).
- **Facets / tags / labels.** Cross-cutting attributes (product, plan tier, OS, user role, region, severity).

Faceted search lets users filter independently of category — "show me Account articles that apply to Enterprise plan on Linux." This is far more powerful than hierarchy alone but requires consistent tagging.

Taxonomy design principles:

- **No more than 7 top-level categories.** Beyond that, users skip the category list and search.
- **Categories are mutually exclusive.** "Billing" or "Account" — pick one for an article about a billing-related account question; tag it with both labels but file it under one.
- **Tags are non-exclusive.** A "How to upgrade your plan" article can be tagged "billing," "account," "self-service," "all-tiers."
- **No "Miscellaneous" or "Other."** If you need it, your taxonomy is wrong.
- **Test by user task.** Take five top user tasks; the path to each from the home page should be obvious within two clicks.

### 9. Cross-linking patterns

Articles in isolation feel orphaned and rank weaker. Productive cross-link patterns:

- **Related articles** at the footer: 3–5 hand-picked or algorithmically suggested.
- **Prerequisite link** at the top: "Before starting, complete [Set up your account]."
- **Next-step link** at the bottom: "After this, you may want to [Connect from your app]."
- **Inline cross-references** mid-article: "See [Connection string format] for the full grammar."
- **Escalation cross-reference** at the bottom: "If this didn't help, [contact support] or [open a ticket]."

The discipline: every article must link to at least one other article. New orphan articles are flagged.

### 10. Evergreen vs time-bound content

- **Evergreen** content stays accurate over time: conceptual articles ("What is sharding?"), reference articles ("MongoDB connection-string format"), policy articles ("Our SLA").
- **Time-bound** content has an expiration: feature announcements ("New in 2026 Q1"), pricing pages, beta program articles, compliance attestations.

Time-bound articles need:

- A visible "applies as of YYYY-MM-DD" line.
- An automatic expiration that flips status to "needs review" after N months.
- An archival rule.

Mixing evergreen and time-bound in the same article causes drift: the evergreen part survives, the time-bound part rots, and the article reads as half-right.

## Templates

### Template 1: KB article skeleton

```markdown
# [Action or question form title]

**Last updated:** YYYY-MM-DD
**Applies to:** [audience scope, e.g., Atlas Dedicated, Enterprise plan, all OS]

## Summary

[1–3 sentences. The reader who reads only this should leave with the answer.]

## Before you begin

- [Prerequisite 1]
- [Prerequisite 2]
- [Required permissions]

## Steps

1. [Step 1 — verb-first imperative]
2. [Step 2]
3. [Step 3]

## Verify

[How to confirm the outcome. A command, a UI state, a result.]

## Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| [Symptom 1] | [Cause] | [Fix or link to dedicated article] |
| [Symptom 2] | [Cause] | [Fix or link] |

## If this didn't help

- Try [link to closely related article].
- See [link to deeper-dive article].
- [Open a support ticket] and include the error code from step N.

## Related articles

- [Related 1]
- [Related 2]
- [Related 3]
```

### Template 2: FAQ entry

```markdown
### Does [product] support [feature]?

[1–3 sentence answer.] For details, see [link to deeper article].
```

### Template 3: KCS in-flight capture (when a ticket reveals a content gap)

```markdown
# [Question as the customer phrased it, lightly cleaned up]

**Status:** Working (internal only — promote to validated after peer review)
**Captured from:** Ticket #12345
**Author:** [agent]
**Captured at:** YYYY-MM-DD

## Quick answer

[The agent's resolution, in 1–3 sentences.]

## Context (when this applies)

[The conditions under which this question shows up.]

## Resolution steps

1. [Step]
2. [Step]

## Notes for the next reviewer

- [Anything the author wasn't sure about; flagged for KCS coach]
- [Anything that might generalize beyond this one ticket]
```

### Template 4: Escalation block at the article foot

```markdown
## If this didn't solve the problem

Try these in order:

1. **Confirm prerequisites.** Re-check the "Before you begin" section above. The most common reason these steps fail is a missing prerequisite.
2. **Search for the specific error message.** If you saw an error code, search the KB for the exact code.
3. **Try the related troubleshooting article:** [link]
4. **Contact support.** When you contact support, include:
   - The error code or message text
   - The step number where the failure occurred
   - Your account ID and cluster name
   - The timestamp (UTC) of the failure
```

### Template 5: Deprecated-article banner

```markdown
> **Deprecated.** This article describes [feature] which was deprecated on
> YYYY-MM-DD and removed on YYYY-MM-DD. For current guidance, see
> [the replacement article].
>
> This article remains available for users on [legacy version / older
> tier], who should still be able to follow these steps until [end date].
```

### Template 6: Taxonomy starter for a B2B SaaS KB

```
Top-level categories (max 7):
  - Getting started
  - Account & billing
  - Product features
  - Integrations
  - Troubleshooting
  - Security & compliance
  - API & developer

Facet dimensions:
  - product line (if multi-product)
  - plan tier (Free / Pro / Enterprise)
  - audience (admin / end-user / developer)
  - OS / platform (Linux / macOS / Windows / iOS / Android / web)
  - severity (informational / impact / outage)
  - lifecycle (current / deprecated / beta)
```

### Template 7: KB article QA checklist

```
[ ] Title is in question-form or action-form.
[ ] Title matches the search query a user would type.
[ ] First paragraph answers the question.
[ ] Article scope ("Applies to:") is explicit in the first 50 words.
[ ] One article, one question — if two, split.
[ ] Steps are numbered, prerequisites are at the top.
[ ] At least one screenshot if the article describes UI flow.
[ ] At least one verification step ("how to confirm it worked").
[ ] Escalation block at the bottom ("if this didn't help…").
[ ] At least one related-article cross-link.
[ ] Last-updated date and reviewer assigned.
[ ] Tags or labels applied per the taxonomy.
[ ] No internal jargon (FTDC, BUWAF, internal acronyms) without expansion.
[ ] No dead links.
[ ] Reading-level appropriate for audience (use the plain-language skill if needed).
```

### Template 8: Lifecycle-state transition guide

```
Author creates article
  └── status: DRAFT
       ↓ (peer review passes)
  └── status: VALIDATED_INTERNAL
       ↓ (KCS coach or content reviewer passes)
  └── status: PUBLISHED (visible to customers)
       ↓ (every 90 days for UI articles, annual for stable)
  └── REVIEW_DUE — triggered by date
       ↓ (still relevant, still accurate)
  └── status: PUBLISHED (refreshed date stamp)

  OR:

  └── status: AT_RISK (low usefulness, ticket data contradicts)
       ↓ (decision: rewrite | merge | archive)

  OR:

  └── status: DEPRECATED (feature removed)
       ↓ (banner added, replacement linked)
       ↓ (after replacement is stable, archive)
  └── status: ARCHIVED (removed from search; historical only)
```

## Anti-patterns

1. **Title that is just the feature name.** "Data Export." Doesn't match user search behavior.
2. **Two questions in one article.** "How to set up SSO and configure user provisioning." Split.
3. **No "Applies to" scope.** Article describes Free Tier behavior but is also linked from Enterprise tickets. Both audiences are frustrated.
4. **Prerequisites in step 3.** "First, install the CLI. (Note: also requires admin access — see step 5 if you haven't enabled that yet.)" Reorder.
5. **Dead-end article.** The article describes the happy path but offers no escalation when the user's situation diverges.
6. **Last-updated date in the year 2022.** Stale-looking content kills trust even when correct. If still accurate, refresh the date.
7. **Article with no cross-links.** Orphan content ranks poorly and feels lonely.
8. **Marketing voice in a KB article.** "Effortlessly unlock the power of…" Users searching for help don't want marketing copy.
9. **Wall-of-text answers without structure.** Bullet points, numbered steps, and bold key terms speed scanning. Walls of prose drive users to "contact support."
10. **"Coming soon" content.** A KB article describing a feature not yet released. If the feature is coming next week, fine. If it's a year out, this is marketing, not KB.
11. **Article with screenshots that don't match the current UI.** Either update screenshots on UI release, or don't include them.
12. **Conflating runbooks with KB articles.** A customer-facing KB article should never include "ssh to the bastion and tail the log." That belongs in a runbook.
13. **Using "click here" for the cross-link.** Use descriptive link text. (See `accessibility-writing` for full discussion.)
14. **No version differentiation.** Article applies to "MongoDB" without specifying 6.0, 7.0, 8.0. Users on the wrong version get wrong instructions.
15. **Tagging everything with everything.** Over-tagging defeats faceted search; every facet matches everything and filters nothing.

## Decision heuristics

| Situation | Choice |
|---|---|
| Question gets asked 5+ times a week in tickets | Write a KB article. Capture in KCS flow next time it comes in. |
| Question gets asked once a year | Don't write a dedicated article. Reference in a broader article if relevant. |
| Article is hit a lot but rated "not helpful" | Rewrite, don't add to it. The structure is wrong. |
| Article is rated "helpful" but rarely surfaces in search | Title and tagging problem. Fix the title; review labels. |
| Question requires a 2,000-word answer | You probably have multiple questions. Split. |
| Question and answer fit in 3 sentences | FAQ entry, not full KB article. |
| Question is "what command do I run for incident X?" | Runbook, not KB article. |
| Question is "what does this checkbox mean?" | Inline help / tooltip, not KB article. |
| Feature is in beta | Either don't publish, or publish with a "Beta" banner and a "subject to change" disclaimer. |
| Feature is deprecated | Replace with a banner pointing to the new path; keep article available until usage drops. |
| Article hasn't been reviewed in 18+ months | Triage: rewrite, merge, or archive. |
| Two articles cover overlapping topics | Merge if they answer the same question. Cross-link if they answer adjacent questions. |
| Information lives in three places | Make one canonical; the others link, do not copy. |

## Platform notes

### Zendesk Guide

- Article titles are heavily weighted in Zendesk search.
- Labels (article-level) drive faceted search and filter ranking; use 5–10 consistent labels per article.
- Use community-content separation: articles for canonical answers, community posts for variability.

### Confluence Service Desk / Confluence

- KB spaces should have a flat or shallow page structure; deep nesting hides articles from search.
- Page labels are the primary faceted dimension; add labels generously but consistently.
- Confluence's search ranks page labels and titles strongly; use page properties for structured metadata.

### Notion

- Database-as-KB pattern: each article is a database entry with status, owner, last-updated, and tag properties.
- Use views (filtered by status, audience, product) to expose different lenses on the same database.
- Notion's search is weaker than Zendesk/Confluence; invest in titles and consistent property values.

### Salesforce Knowledge

- Data categories are the hierarchy; channels (Internal/Customer/Partner/Public) control visibility.
- Article types differentiate templates (FAQ vs How-To vs Troubleshooting).
- Validation rules can enforce required fields (audience, scope, prerequisites).

### MadCap Flare

- Topic-based authoring: one topic per question.
- Variables for product names, version numbers, customer names — change once, propagate everywhere.
- Conditional tags differentiate output: internal-only sentences vs customer-facing.
- Snippets for shared sentences across topics (e.g., a common preamble).

## Cross-skill notes

- **Use `doc-archaeology` for KB-wide audits.** That skill scans an existing KB for stale, redundant, contradicting, or orphan content. This skill writes new articles and improves individual ones.
- **Use `runbook-craft` for engineer-facing runbooks.** A runbook is for the on-call SRE under pressure. A KB article is for the end user.
- **Use `plain-language` for reading-level targeting.** Customer-facing KB articles typically aim for Flesch-Kincaid grade 8.
- **Use `microcopy-and-ui-writing` for inline help and tooltips.** When the answer fits in a tooltip, write a tooltip, not a KB article.
- **Use `support-ticket-writing` for the ticket-side of the same conversation.** A KB article and a ticket macro often answer the same question; both should match.
- **Use `accessibility-writing`** for headings, link text, alt text, and color-independence inside KB articles.

## References

1. Consortium for Service Innovation, *KCS v6 Practices Guide — Section 2 The KCS Practices*: https://library.serviceinnovation.org/KCS/KCS_v6/KCS_v6_Practices_Guide/030
2. Atlassian, *What is knowledge-centered service (KCS)?*: https://www.atlassian.com/itsm/knowledge-management/kcs
3. Zendesk, *Best practices: Developing content for your knowledge base*: https://support.zendesk.com/hc/en-us/articles/4408831743258-Best-practices-Developing-content-for-your-knowledge-base
4. Zendesk, *Best practices: Four steps to a streamlined knowledge base*: https://support.zendesk.com/hc/en-us/articles/4408883031706
5. Zendesk, *8 knowledge base article templates that work*: https://www.zendesk.com/blog/knowledge-base-article-template/
6. Salesforce, *What is Knowledge-Centered Service (KCS)?*: https://www.salesforce.com/service/knowledge-centered-service/
7. Atlassian Confluence, *Write, search, and share knowledge base articles*: https://confluence.atlassian.com/spaces/SERVICEDESKSERVER/pages/956713329/Write+search+and+share+knowledge+base+articles
8. SearchUnify, *Top KCS Metrics for Customer Support Teams*: https://www.searchunify.com/resource-center/short-articles/top-kcs-metrics-to-measure-knowledge-adoption-in-customer-support-teams/
9. MadCap Software, *Single Sourcing Explained: The Power of MadCap Flare*: https://www.madcapsoftware.com/blog/single-sourcing-explained-the-power-of-madcap-flare/
10. Supportbench, *Runbooks vs KB Articles: Use Cases*: https://www.supportbench.com/runbooks-vs-kb-articles-when-each-is-the-right-tool/
