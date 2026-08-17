<!-- hub-reference-banner -->
> **Reference file — part of the `executive-comms` hub.** Formerly the standalone `case-study-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: case-study-writing
description: Customer success story craft — the Challenge / Solution / Result template (BLUF for customer stories), customer-quote integration, anonymization decisions, "X% improvement" claim hygiene, before/after comparisons, multi-stakeholder voice (champion vs CFO vs end-user), the long-form / mid-form / one-pager derivative pipeline, written approval workflow, and the case-study-vs-testimonial-vs-reference-architecture distinction. TRIGGER when user asks to write/draft/review/critique a customer case study, customer success story, customer reference, before/after story, customer outcome write-up, or wants to evaluate a quote/metric/approval workflow for a case study; questions about whether something is a case study, a testimonial, or a reference architecture; questions about anonymizing a customer or hygiene of percentage claims. SKIP for product launch announcements (use release-blog-and-launch-narrative); short marketing/conversion copy (use sales-and-marketing-copy); long-form thought-leadership whitepapers (use whitepaper-writing); technical reference architectures for engineers (use rfc-and-design-docs); internal QBRs or account reviews (use executive-comms or tam-account-reports).
category: custom
tags: [writing, external-comms, case-study, customer-success, b2b-marketing]
---

# Case Study Writing

## Overview

A case study is a published narrative of one customer's outcome with a vendor's product or approach. It is the highest-conversion content asset in B2B marketing because it does the one thing prospects most need: prove that a buyer like them succeeded with this solution.

**The defining traits of a case study**:

- Specific, named (or carefully anonymized) customer
- Concrete before/after with quantified outcomes
- At least one direct quote from a named customer stakeholder
- A clear narrative arc: Challenge → Solution → Result
- Customer approval in writing before publication
- A non-marketing tone — the customer is the protagonist, not the vendor

A case study is **not**:

- A **testimonial** (a single quote, no narrative, no metrics — useful but lower-trust)
- A **reference architecture** (technical-design document explaining how to build, not a customer story)
- A **logo wall** (the customer permits its logo to appear; no narrative, no quote)
- A **press release** (the news of the customer signing — not the story of the outcome)
- A **product page** (vendor-voice description of capability)

The genre is old (consulting firms invented it in the 1960s) and remarkably stable. The 2026 evolution is the "long-form, mid-form, one-pager" derivative pipeline: one customer story produces three or four sizes of asset, each tuned to a different stage of the buyer funnel.

## Core concepts

### 1. Challenge / Solution / Result — the BLUF for customer stories

The dominant case study structure. Every section has a job.

| Section | Job | Length |
|---|---|---|
| Title + hero metric | Sell the click in five words; promise the outcome | 1 line + 1 stat |
| Challenge | What was broken before, why it mattered, what was at stake | 25% of body |
| Solution | What the customer adopted, how the implementation went | 35% of body |
| Result | Quantified outcomes; customer quote; what changed | 35% of body |
| About the customer + What's next | Context block + forward-looking statement | 5% of body |

**Front-load the result**. Place the hero metric in the title or subtitle, and again in the first paragraph. A reader who reads only the headline should know the outcome.

```text
TITLE:   How [Customer] Cut [Cost/Time/Risk Metric] by [X%] in [Time Period]
SUBTITLE: After replacing [legacy approach] with [new approach], [Customer] [achieved specific outcome]
```

### 2. Anonymization decisions

Not every customer can be named. The anonymization decision is binary and load-bearing.

**Name the customer when**:

- Written marketing approval is signed (PR + Legal + executive)
- The customer is referenceable (i.e., happy to be contacted by prospects)
- The story does not expose competitive intelligence, regulatory posture, or security architecture in ways the customer would later regret
- The customer logo will help, not hurt, your other prospects

**Anonymize when**:

- The customer is in a regulated industry where naming raises legal or compliance risk (financial services, healthcare, government)
- The story includes a sensitive admission ("we had a breach," "we lost a major deal")
- The customer is willing to share their story but not their name
- The customer is mid-acquisition or in a quiet period

**Anonymization patterns** (in descending specificity):

| Pattern | Example | Trust signal |
|---|---|---|
| Industry + size + region | "A Fortune 500 financial-services firm headquartered in the U.S. Northeast" | High |
| Industry + size | "A global retail chain with 1,200 stores" | Medium |
| Industry only | "A large healthcare provider" | Low |
| Generic ("a customer") | "One of our customers" | Very low; close to useless |

**A blind-but-verified anonymized case study** — where the publication is anonymous but a sales rep can name the customer in a 1:1 conversation under NDA — preserves both the customer's confidentiality and the prospect's trust.

### 3. The "X% improvement" claim hygiene

Quantified outcomes are the entire reason case studies convert. Every metric must be defensible.

**Rules**:

- **Specific baseline**: "down from 4.2 hours to 12 minutes" beats "much faster"
- **Specific time window**: "in the first 90 days" beats "over time"
- **Specific scope**: "across 1,200 trades per day" beats "in our trading operation"
- **Verifiable source**: the customer signed off on this number, in writing, with a name attached
- **No misleading aggregation**: "10x improvement" that is actually 10x on one metric and 0.8x on another is dishonest

**Defensible metric pattern**:

> "Before: customer support tickets averaged 4.2 hours to resolution.
> After: 12 minutes. Measured across 1,200 weekly tickets in Q3 2026.
> Source: Customer's internal support analytics, validated by [their team lead, name + title]."

**Anti-pattern**: vague hero metrics like "transformed our operations" or "10x productivity." Both are unverifiable and read as marketing.

**Special case — percentage gaming**:

| Honest | Misleading |
|---|---|
| "Reduced incident response time by 73%, from 45 minutes to 12 minutes" | "73% faster response" (no baseline) |
| "Saved $2.4M annually across three product lines" | "Saved millions" (no figure) |
| "Increased uptime from 99.5% to 99.95%" | "10x more reliable" (technically true; misleading) |

### 4. Quote integration

Every case study needs at least one direct customer quote. The quote is where the customer's voice carries the credibility the vendor's voice cannot.

**Quote rules**:

- Attributed to a named individual with a title (e.g., "Sarah Chen, Director of Engineering")
- Sounds like a human said it (not vendor-speak)
- Specific, not generic ("we cut our deploy time in half" beats "Acme has been a great partner")
- One feature-praise quote is acceptable; the strongest quotes describe outcomes
- 2–4 quotes per case study; one per stakeholder

**Quote-vetting test**: read it aloud. If it sounds like the vendor's PR team wrote it, the customer will eventually disclaim it.

**Example — bad quote**:

> "Acme is a true partner who has helped us leverage our data to drive transformational outcomes." — Director of IT

**Example — good quote**:

> "We used to spend two full days every month reconciling reports. Now it's a 20-minute job for one person. That gave my team a week back each month to work on actual product." — Sarah Chen, Director of Engineering

### 5. Multi-stakeholder voice

B2B decisions are multi-stakeholder. A case study that quotes only one champion misses the buying-committee dynamic.

**Voice palette**:

| Stakeholder | What they care about | Quote topic |
|---|---|---|
| Champion (the person who advocated for adoption) | Vindication; outcomes; the story | "Here's why we chose this and what changed" |
| End-user (the person who uses the product daily) | Daily experience; speed; pain reduction | "Here's what my work feels like now" |
| Finance (CFO or finance partner) | Cost, ROI, payback period | "Here's the financial picture" |
| Engineering/IT lead | Architecture, integration, reliability | "Here's how it fits with our other systems" |
| Executive sponsor (VP or C-level) | Strategic outcome, alignment with priorities | "Here's why this mattered to the company" |

A long-form case study (1,500+ words) should quote at least three of these. A one-pager may have only the champion. A mid-form will have two.

### 6. Long-form / mid-form / one-pager derivative pipeline

One customer story should yield multiple assets, each tuned to a different funnel stage.

| Format | Length | Funnel stage | Pages/words |
|---|---|---|---|
| Long-form case study | 1,200–2,500 words | Mid-funnel (active evaluation) | 3–6 pages |
| Mid-form case study | 600–1,000 words | Mid-to-top funnel (consideration) | 1–2 pages |
| One-pager / case study summary | 200–400 words | Sales enablement / outbound enclosure | 1 page |
| Customer video (2–3 min) | Scripted from long-form | Top-funnel awareness | N/A |
| Quote card (single quote + metric) | 50 words | Social media; sales decks | 1 graphic |
| Case study within a whitepaper | 1 page sidebar | Cross-asset use | embedded |

**Derivative discipline**:

- Approve **all derivative formats** in the original customer approval workflow. Customers hate finding out their story is on a billboard.
- Maintain a single source of truth (typically the long-form) and derive the others.
- Update derivatives when the long-form is revised.

### 7. The written approval workflow

Publishing a case study without written customer approval is malpractice. The workflow has named gates.

**Standard workflow**:

1. **Customer success identifies a candidate**: the customer is happy, has a quantifiable outcome, and is open to a story.
2. **Verbal pre-approval from a customer executive**: a Director-or-above says "yes, in principle."
3. **Discovery interview**: 30–60 minutes with the champion (and ideally the end-user or finance partner separately).
4. **Draft 1 to the customer**: shared with the champion for accuracy review.
5. **Customer Legal and PR review**: this is the long pole; budget 2–4 weeks. If the customer is regulated (finance, healthcare, government), budget 6–12 weeks.
6. **Final draft + approval signature**: written confirmation of the published text, the approved derivative formats (web, PDF, video, ads, conferences, etc.), and the approved channels.
7. **Publication**.
8. **Renewal**: re-confirm approval annually if the case study remains live; remove or anonymize if the customer relationship ends.

**Approval scope** — capture explicitly in the signed approval:

- Approved customer name and logo use
- Approved quoted individuals and their titles
- Approved metrics (each one)
- Approved channels (web, PDF download, paid ads, sales decks, conference talks)
- Approved derivatives (one-pager, video, social posts)
- Duration of approval (perpetual, or N-year)
- Renewal expectation

### 8. Case study vs testimonial vs reference architecture

These three asset types are routinely confused. The distinction is structural and matters for production effort.

| Aspect | Case study | Testimonial | Reference architecture |
|---|---|---|---|
| Length | 200–2,500 words | 1–3 sentences | 5–30 pages |
| Narrative arc | Challenge → Solution → Result | None | Technical decision rationale |
| Metrics | Required | Optional | Not the point |
| Customer quote | Multiple, attributed | The asset itself | Often anonymized |
| Production effort | 4–12 weeks | Hours | 4–8 weeks |
| Use | Sales/marketing primary asset | Sales decks, landing pages, social | Technical-buyer evidence |
| Approval workflow | Heavy (Legal, PR, exec) | Light (single signoff) | Medium (technical accuracy review) |

**Rule of thumb**: if it has a story arc and a number, it's a case study. If it's a quote, it's a testimonial. If it's a diagram and a deployment-pattern explanation, it's a reference architecture.

### 9. The before/after frame

The before/after frame is the single most reusable case study structure. It works because human readers map their own situation onto the "before" picture and project themselves into the "after" picture.

**Pattern**:

```text
Before:
  - Specific painful situation
  - Quantified cost (time, money, risk)
  - Named consequence (a missed quarter, a security incident, customer churn)

After:
  - Specific changed situation
  - Quantified improvement
  - Named consequence (a hit quarter, no security incident, retained customers)
```

**Visual representation**: side-by-side table, before/after photo, or simple two-column comparison. The visual is almost always reusable as a social-media graphic.

### 10. Tone calibration: the customer is the hero

The single hardest tone discipline in case study writing: the vendor must not be the hero. The customer is.

**Vendor-as-hero (anti-pattern)**:

> "Acme's revolutionary platform transformed BigCo's operations, delivering 73% productivity gains through industry-leading capabilities."

**Customer-as-hero (correct)**:

> "BigCo's engineering team cut deploy time from 4 hours to 22 minutes by replacing their legacy CI system with [solution]. The team — led by Sarah Chen — chose the new platform after evaluating four alternatives in a six-week bake-off."

The vendor enters the story only when it's structurally required (which product was chosen, how implementation went). The customer drives the narrative; the vendor is the tool they used.

## Templates and examples

### Long-form case study skeleton (1,200–2,500 words)

```text
TITLE: How [Customer] [Hero Outcome] with [Solution Category]
SUBTITLE: One-line context: industry, size, and what changed

Hero metric box (sidebar, top of page):
  - [X%] reduction in [metric]
  - [$Y] saved annually
  - [Z weeks] payback period

Section 1: At a glance (3–5 bullets, 80 words)
  - Customer: [name and one-line description]
  - Challenge: [one line]
  - Solution: [one line]
  - Result: [one line — the hero metric]

Section 2: The Challenge (300–500 words)
  - Opening scene: a specific moment of pain
  - The systemic problem behind it
  - What was at stake (revenue, risk, time)
  - Why prior approaches failed
  - Quote from the champion describing the pain

Section 3: The Solution (400–700 words)
  - The evaluation process (briefly)
  - What was chosen and why
  - Implementation: the timeline, the team, the integration
  - Quote from a technical or end-user stakeholder

Section 4: The Result (400–700 words)
  - The quantified outcomes (each with baseline, time window, and source)
  - The qualitative changes
  - Forward-looking statement: what the customer plans next
  - Quote from an executive sponsor or finance partner

Section 5: About [Customer] (50–80 words)
  - Standard "about" block
  - Logo
  - Optional: a "Why we chose [vendor]" sidebar

Footer:
  - Approved-by line
  - Date of publication
  - Optional: link to one-pager PDF
```

### Mid-form case study skeleton (600–1,000 words)

```text
TITLE: [Customer] [Hero Outcome] (one line)

Hero metric box
  - [Single most compelling metric]

Lead paragraph (80 words):
  Five Ws + the hero metric

Challenge (~200 words):
  Brief problem statement; one quote

Solution (~250 words):
  What was chosen, how implementation went; one quote

Result (~300 words):
  Quantified outcomes; one final quote

About + CTA (~50 words)
```

### One-pager case study skeleton (200–400 words)

```text
Header:
  [Customer Name + Logo]
  [Hero Outcome in one line]

Three boxes side-by-side:
  CHALLENGE: [50 words]
  SOLUTION: [50 words]
  RESULT: [50 words + the hero metric]

One pull quote:
  "[Best customer quote, 25–40 words]"
  — Name, Title

About [Customer]: [40 words]
CTA: [10 words + URL]
```

### Anonymized case study skeleton

Same structure as named, but:

- Title: "How a Fortune 500 Financial-Services Firm Cut..."
- Replace customer name with the agreed anonymization pattern throughout
- Quote attribution: "Director of Engineering at a global retail chain"
- About block: industry + size + region only
- Maintain a confidential customer log internally with the true identity for sales-rep reference under NDA

## Anti-patterns

| Anti-pattern | Why it fails | Fix |
|---|---|---|
| Vendor-as-hero voice | Buyers read it as ad copy and discount it | Customer is the protagonist; vendor is the tool |
| Vague hero metric ("transformed operations") | Unverifiable; reads as marketing fluff | Specific number with baseline, time window, and scope |
| Quotes that sound like the PR team wrote them | Customer will eventually disclaim; reads as fake | Vet quotes: would a human actually say this? |
| Only one stakeholder quoted | Misses the buying-committee dynamic | Quote 2–3 stakeholders: champion, end-user, finance/exec |
| No written approval | Legal/PR risk; customer relationship damage when discovered | Sign a written approval covering name, quotes, metrics, channels, derivatives |
| Publishing approved content in unapproved channels | Customer feels misled; relationship damage | Approval scope explicit in signed doc; track derivatives |
| Skipping the "Challenge" section | The case study becomes a feature list | Always lead with the painful before-state |
| Aggregating metrics in misleading ways ("10x") | Reads as honest; is actually misleading | One specific metric per claim; show the baseline |
| Using "leading," "transformational," "revolutionary" | AI-ism / marketing-ism that signals low rigor | Specific, factual language; let the numbers carry the weight |
| Anonymization too vague ("one of our customers") | Reads as a fabricated story | Specify industry + size + region at minimum |
| No "what's next" section | The story ends; the relationship reads as static | Include a forward-looking statement: customer's next phase |
| Renewing case studies indefinitely without re-approval | Customer may have changed, churned, or moved on | Annual approval renewal for live case studies |
| Mixing case study with reference architecture | Confuses buyer; neither asset does its job | Two separate assets, each tuned to its audience |

## Decision heuristics

**Named or anonymized**

- Customer happy + signed approval + referenceable → named
- Customer happy but in a regulated industry or competitively sensitive → anonymized with high specificity (industry + size + region)
- Customer story is great but no approval likely → write internally for sales-rep talk track; do not publish

**Which derivative format to lead with**

- If sales is requesting it for outbound → one-pager first
- If marketing is staffing a content campaign → long-form first, derive the others
- If the customer wants minimal exposure → mid-form only; no video, no one-pager

**How many metrics to include**

- 1 hero metric (in the title and hero box)
- 2–4 supporting metrics in the body
- More than 5 starts to feel like a dashboard rather than a story

**When to include a video**

- The customer is willing and available
- The customer's voice and presence will add credibility (e.g., an admired industry figure)
- The product/outcome benefits from visual demonstration

**When the case study is "done"**

- The Challenge / Solution / Result arc is complete
- Every metric has a baseline, time window, scope, and source
- All quotes are attributed to a named individual with a title
- The customer has signed approval covering text, quotes, metrics, channels, and derivatives
- The hero metric appears in the title and the first paragraph
- A final fact-check pass has been done with the champion

**Refresh cadence for live case studies**

- Annual: re-approve and re-verify metrics
- 18 months: consider a refresh interview with new outcomes
- 24 months: usually time to retire or substantially rewrite

## Cross-skill collisions and routing

- **Product launch announcement / news release** → `content-and-marketing-writing` (references/press-release-writing.md) (case study can be quoted from in the release)
- **Owned-blog launch narrative** → `content-and-marketing-writing` (references/release-blog-and-launch-narrative.md)
- **Short marketing/conversion copy or landing page** → `content-and-marketing-writing` (references/sales-and-marketing-copy.md)
- **Long-form thought-leadership whitepaper** → whitepaper-writing (case study often appears as a sidebar inside)
- **Reference architecture / technical-buyer evidence** → rfc-and-design-docs
- **Internal QBR or executive account review** → `executive-comms` or `tam-account-reports`
- **Customer testimonial (a single quote, no narrative)** → handled as a sub-asset of this skill, but does not need its own case study
- **Customer logo wall (logo permission only, no story)** → not in scope; handled via standard logo-use approval

## References

- [How to Turn a Case Study into a Customer Success Story — HubSpot Marketing Blog](https://blog.hubspot.com/marketing/customer-success-story)
- [B2B Case Study Template: Problem-Solution-Results Framework — Libril](https://libril.com/blog/b2b-case-study-template)
- [How to Write an Effective Case Study in 8 Steps — Uplift Content](https://www.upliftcontent.com/blog/how-to-write-a-case-study/)
- [Customer Case Study: A Proven Framework to Build Trust and Revenue — Testimonial.to](https://testimonial.to/resources/customer-case-study)
- [Case Studies Had Their Run. In 2026? It's Going To Be All About Customer Evidence. — UserEvidence](https://userevidence.com/blog/case-studies-had-their-run-in-2026-its-going-to-be-all-about-customer-evidence/)
- [Writing a B2B Case Study: A Template and Best Practices — Liger Marketing](https://ligermarketing.com/writing-a-b2b-case-study-a-template-and-best-practices/)
- [A Roundup of Case Study Examples Every Marketer Should See — HubSpot](https://blog.hubspot.com/marketing/case-study-examples)
