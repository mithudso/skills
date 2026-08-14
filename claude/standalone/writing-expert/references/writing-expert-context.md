# Writing Expert — Reference Context

Comprehensive reference for technical, business, and report writing. Covers document types, frameworks, formatting, tone, anti-AI-ism enforcement, and data storytelling.

---

## 1. Technical Writing

### Structure
Every technical document answers: What is this? Who is it for? What can I do with it? How do I do it?

**Universal template:** Purpose (one sentence) → Prerequisites → Core content → Verification → Troubleshooting.

### Document types

**Architecture docs:** Start with a context diagram (boxes and arrows showing major components). Follow with: component inventory, data flow, deployment topology, security boundaries, and decision log. Keep decisions and their rationale — they're the most valuable part.

**API documentation:** For each endpoint: HTTP method + path, description (one sentence), request parameters (table with name/type/required/description), request body (JSON example), response (JSON example with status codes), error codes, and a curl example. Document edge cases, not just the happy path.

**Runbooks:** Imperative mood, one action per step, numbered. Every step has a verification ("You should see..."). Include rollback instructions. Assume the reader is stressed and sleep-deprived — optimize for scannability.

**README:** Title → one-paragraph description → quick start (fewest possible steps to working) → full installation → usage examples → configuration → contributing → license. The README is marketing — if the quick start has more than 5 steps, simplify.

**Troubleshooting guides:** Symptom-first organization (what the user sees, not what's broken internally). Decision tree format: "If you see X, try Y. If that doesn't work, check Z." Include the actual error messages users will copy-paste into search.

### Principles
- **Precision over brevity** — "The function returns `null` when the input array is empty" beats "Handles edge cases."
- **Examples are mandatory** — Every concept needs at least one concrete example.
- **Progressive disclosure** — Quick start first, deep details later. Don't front-load complexity.
- **Audience adaptation** — A developer tutorial reads differently than an operations runbook. State the audience explicitly.

---

## 2. Business Writing

### BLUF (Bottom Line Up Front)
Military-origin principle: put the conclusion in the first sentence. The reader should understand the message if they only read the first paragraph.

**Structure:** Bottom line → context → supporting details → action required → deadline.

**Example:**
> We need to increase the support team by 2 FTEs by Q3 to maintain SLA compliance. Case volume grew 34% QoQ while staffing remained flat. Without additional capacity, we project SLA breach rates will exceed 15% by August.

### Minto Pyramid Principle
Answer first, then group supporting arguments into 3-5 mutually exclusive, collectively exhaustive (MECE) categories, each backed by data.

```
Answer: We should migrate to Atlas.
├── Cost: 23% reduction in total infrastructure spend
├── Reliability: 99.995% uptime SLA vs current 99.9%
└── Velocity: 40% faster provisioning for new environments
```

### SCQA (Situation-Complication-Question-Answer)
For problem-analysis documents:
- **Situation:** What's the current state? (undisputed facts)
- **Complication:** What changed or went wrong?
- **Question:** What do we do about it? (the reader's implicit question)
- **Answer:** Our recommendation with evidence.

### Document types

**Executive summaries:** Max one page. 3-5 key metrics, delta-focused (what changed, not what is). Action items with owners and dates. No background the reader already knows.

**Status reports:** TL;DR → completed → in progress → blocked → next period → metrics. Use the same structure every time — consistency lets readers scan for changes.

**Proposals / business cases:** Problem statement → proposed solution → alternatives considered → cost analysis → timeline → risks → recommendation. Quantify everything possible.

**Meeting minutes:** Date, attendees, decisions made (numbered), action items (owner + date), open questions. Skip the discussion recap — decisions and actions are what matter.

---

## 3. Report Writing

### Account reviews / QBRs
**Structure:** Executive summary → account health scorecard → key metrics with trends → support case analysis → project status → risks and opportunities → action plan → appendix.

**Data storytelling pattern:** Setup (what we expected) → Conflict (what actually happened) → Resolution (what we did or recommend). Every metric needs context: "Latency is 230ms" is data. "Latency increased from 150ms to 230ms after the shard migration, exceeding the 200ms SLA target" is a story.

### Incident post-mortems
**Blameless by design.** Focus on systems, not people. Use "the monitoring system failed to alert" not "the on-call engineer missed the alert."

**Required sections:** Impact summary (duration, affected users, data loss), timeline (UTC timestamps), root cause (specific, not vague), contributing factors (what made the impact worse), remediation (what we fixed, what we're still fixing with owners/dates), lessons learned (what we'll do differently).

### Data-driven narratives
- Lead with the insight, not the data: "Customer churn doubled in Q2" not "Here is a chart of Q2 churn rates."
- Comparisons create meaning: absolute numbers alone are meaningless without benchmarks, trends, or targets.
- Use "compared to" language: "Response time improved 40% compared to Q1" or "3x faster than the industry median."
- Round numbers for readability: "$2.3M" not "$2,287,419.38" in prose. Exact figures go in tables.

---

## 4. Document Formatting

### Markdown best practices
- **Heading hierarchy:** H1 = document title (only one). H2 = major sections. H3 = subsections. Never skip levels (no H1 → H3).
- **Bullets vs numbered lists:** Bullets for parallel items with no order. Numbers for sequential steps or ranked items.
- **Tables vs prose:** Use tables for structured comparisons (features, metrics, options). Use prose for causal reasoning, narratives, and explanations.
- **Code blocks:** Always specify the language for syntax highlighting. Use inline code for identifiers, file paths, and CLI commands within prose.
- **Bold for emphasis:** Sparingly — max 2-3 bolded phrases per section. Bold the key takeaway, not every term.
- **Links:** Descriptive text, never "click here." Use relative links for internal references.

### Document length guidelines
| Document type | Target length |
|--------------|--------------|
| Slack message | 1-3 sentences |
| Email | 5-10 sentences |
| Status report | 1 page |
| Executive summary | 1 page |
| Post-mortem | 2-3 pages |
| Architecture doc | 3-10 pages |
| Runbook | 1-5 pages (per procedure) |
| QBR / account review | 5-15 pages |

### When to use which format
- **Bullet points:** Parallel items, checklists, requirements lists, action items
- **Numbered steps:** Procedures, prioritized lists, timelines
- **Tables:** Feature comparisons, metrics, option analysis, reference data
- **Prose paragraphs:** Explanations, reasoning, narratives, context-setting
- **Diagrams:** Architecture, data flow, decision trees, timelines

---

## 5. Tone and Voice

### The register spectrum
**Formal:** Executive communications, legal/compliance docs, customer escalations. Full sentences, no contractions, no slang. "We recommend proceeding with Option B."

**Semi-formal:** Internal reports, cross-team emails, documentation. Contractions OK, clear and direct. "We'd recommend Option B based on the cost analysis."

**Professional-conversational:** Slack, internal wiki, team updates. Short sentences, direct, personality OK. "Option B wins on cost. Here's the breakdown."

**Conversational:** Chat, casual internal comms. Fragments OK, colloquial language OK. "Go with B. Way cheaper."

### Empathy in customer-facing writing
- Acknowledge the impact before diving into technical details: "I understand this is affecting your production workload."
- Use "we" language for shared ownership: "Let's work through this together" not "You need to fix your configuration."
- State timeline expectations: "I'll have an update by 3pm UTC" not "I'll look into it."
- Close with a concrete next step, not platitudes.

### Active vs passive voice
**Active (default):** Subject performs the action. "The migration script updates the schema." Clearer, shorter, more engaging.

**Passive (when appropriate):** When the actor is unknown ("The error was introduced in v2.3"), when the focus is on the object ("All customer data was backed up before the migration"), or in formal/legal contexts.

---

## 6. Writing Frameworks

### STAR (Situation, Task, Action, Result)
For case studies, achievement narratives, and interview-style writeups. Especially effective for demonstrating impact.

### SCQA (Situation, Complication, Question, Answer)
For problem analysis, proposals, and consulting-style memos. Creates narrative tension that holds attention.

### Minto Pyramid
For executive communications and decision documents. Answer first, supporting evidence arranged in logical groups.

### 5W1H (Who, What, When, Where, Why, How)
For investigative writing, incident reports, and requirements gathering. Ensures completeness.

### Inverted Pyramid
For announcements, news-style updates, and time-pressed readers. Most important information first, background last.

---

## 7. Anti-AI Writing Guide

### Why this matters
LLM-generated text has statistically detectable patterns that reduce credibility with experienced readers. The goal is not deception — it's producing prose that respects the reader's time and doesn't trigger "this is AI slop" reactions.

### Vocabulary rules

**Tier 1 — Always replace (5-20x more frequent in AI text):**

| Never use | Use instead |
|-----------|------------|
| delve | examine, explore, look at |
| leverage | use |
| robust | strong, reliable, durable |
| paradigm | model, approach, method |
| seamless | smooth, easy, integrated |
| utilize | use |
| facilitate | help, enable, support |
| commence | start, begin |
| furthermore | also, and |
| it's important to note | (delete — just state the fact) |
| in today's rapidly evolving | (delete — add nothing) |
| game-changer | (describe what specifically changed) |
| navigate (metaphorical) | handle, manage, deal with |
| landscape | field, space, area |
| cutting-edge | new, latest, advanced |
| holistic | complete, comprehensive, full |

**Tier 2 — Flag when clustered (2+ in one section):**
harness, foster, resonate, empower, unlock, drive (metaphorical), journey, ecosystem, transform, innovative, dynamic, significant

**Tier 3 — Structural patterns:**
- Excessive em dashes (target: 0 per 1,000 words in professional writing)
- Uniform paragraph lengths (vary between 1-5 sentences)
- Formulaic openings ("In the world of...", "When it comes to...", "In an era of...")
- Hedge-stacking ("could potentially", "may eventually", "might arguably")
- Generic conclusions ("In conclusion, X remains crucial for Y")
- Chatbot artifacts ("I hope this helps!", "Let me know if you have questions!", "Let's dive in!")
- Bold overuse (more than 3 per section)
- Significance inflation on routine events

### Five principles for human-sounding prose
1. **Vary sentence length.** Mix short fragments with longer constructions. Three words. Then a longer sentence that develops the thought with specific detail and a concrete example drawn from the actual situation.
2. **Be concrete.** Replace vague claims with numbers, names, dates, examples. "Performance degraded" → "P99 latency jumped from 50ms to 340ms after the 14:32 deploy."
3. **Have a voice.** State preferences and reactions where appropriate. Start sentences with "But" and "And." Use incomplete sentences for emphasis.
4. **Cut the neutrality.** If the piece takes a position, take it. Don't hedge with "could potentially" — say "will" or "won't."
5. **Earn emphasis.** Don't tell the reader something is interesting. Make it interesting by presenting the substance.

### Synonym cycling (don't do it)
Humans repeat the clearest word. If "developers" is the right term, use it three times rather than cycling through "developers... engineers... practitioners... builders." Forced variation reads as thesaurus abuse and is a strong AI tell.

### Context tolerance matrix
| Content type | Tier 1 | Tier 2 | Tier 3 | Em dashes |
|-------------|--------|--------|--------|-----------|
| Investor/exec email | Zero tolerance | Zero tolerance | Strict | 0 |
| Customer summary | Zero tolerance | Flag clusters | Moderate | 0-1 |
| Technical blog | Replace obvious | Relaxed | Moderate | 0-2 |
| Internal docs | Replace obvious | Relaxed | Relaxed | OK |
| Slack/chat | Relaxed | Relaxed | N/A | OK |

---

## 8. Data Storytelling

### The narrative arc for data
1. **Setup:** What was the expectation or baseline? ("We targeted 99.9% uptime for Q2")
2. **Conflict:** What happened that deviated? ("Three unplanned outages brought us to 99.7%")
3. **Resolution:** What did we do or recommend? ("We implemented automated failover, which has prevented 4 potential outages since")

### Principles
- **Lead with the insight, not the chart.** "Customer churn doubled" is the headline. The chart is supporting evidence.
- **Comparisons create meaning.** "$2.3M" means nothing alone. "$2.3M, a 34% increase over Q1 and 2x our forecast" tells a story.
- **Round for prose, precise for tables.** Write "$2.3M" in paragraphs, put "$2,287,419" in the data table.
- **Every metric needs a "so what."** Don't just report numbers — explain why they matter and what action they imply.
- **Trend > snapshot.** Show direction of change, not just current state. "Improving from 85% to 94% over 6 months" beats "Currently at 94%."

---

---

## 9. QBR Deep Dive

Design for a 50-60 minute meeting: 30 minutes status, 30 minutes decisions. Distribute a 5-8 page pre-read 48 hours ahead. Select 5-7 KPIs (more creates noise). For every miss, explain the "why." Use trend arrows and color-coding for instant visual context. Restate client goals using their exact words.

**Structure:** Client goals restated → progress against goals with metrics → key wins (specific outcomes, not vague claims) → risks with proposed mitigations → recommended next steps → appendix with detailed data.

---

## 10. Post-Mortem Deep Dive

Focus on learning, not filing. Use blameless language: ask "How did it make sense for someone to do what they did?" not "Why did they do that?" Get granular: don't say "traffic spike," say "traffic increased 340% in 90 seconds."

**Timeline format:** What seemed unusual to participants at each moment, not just what happened. Include what went well, not just what failed. Address systemic conditions in action items, not individual behavior.

---

## Sources

- [Claude Skills for Knowledge Extraction & Report Writing - Medium](https://medium.com/ai-simplified-in-plain-english/claude-skills-for-knowledge-extraction-report-writing-the-2026-enterprise-playbook-b50ebcd2f71d)
- [How To Write Without Sounding Like AI - George Kao](https://georgekao.medium.com/how-to-write-without-sounding-like-ai-e2e0d5930adb)
- [The "Anti-AI" Writing Cheat Sheet - Tom Orbach](https://www.marketingideas.com/p/the-anti-ai-writing-cheat-sheet)
- [avoid-ai-writing skill - GitHub](https://github.com/conorbronsdon/avoid-ai-writing)
- [Report Writing Skill - MCP Market](https://mcpmarket.com/tools/skills/report-writing-1)
- [Don't Write Like AI - Blake Stockton](https://www.blakestockton.com/dont-write-like-ai-1-101-negation/)
- [Best Claude Code Skills 2026 - Firecrawl](https://www.firecrawl.dev/blog/best-claude-code-skills)
- [Claude AI for Business & Startups Guide - Badal Khatri](https://www.badalkhatri.com/blog/claude-ai-skillset-for-business-and-startups-the-complete-2026-guide)
- [Google Technical Writing One](https://developers.google.com/tech-writing/one)
- [Google Developer Documentation Style Guide](https://developers.google.com/style)
- [API Documentation Best Practices - Postman](https://www.postman.com/api-platform/api-documentation/)
- [Pyramid Principle - Slideworks](https://slideworks.io/resources/the-pyramid-principle-mckinsey-toolbox-with-examples)
- [SCQA Framework - Management Consulted](https://managementconsulted.com/scqa-framework/)
- [Incident Postmortem Best Practices - Pragmatic Engineer](https://blog.pragmaticengineer.com/postmortem-best-practices/)
- [Data Storytelling Frameworks - Beautiful.ai](https://www.beautiful.ai/blog/data-storytelling-that-works-5-proof-backed-frameworks-for-communicating-insights-clearly)
- [QBR Templates - ClearPoint Strategy](https://www.clearpointstrategy.com/blog/quarterly-business-review-templates)
- [AI-Generated Prose vs Human Writing - Reuters/Oxford](https://reutersinstitute.politics.ox.ac.uk/news/how-ai-generated-prose-diverges-human-writing-and-why-it-matters)
