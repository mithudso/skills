---
description: >-
  Executive & business communication hub — leadership, persuasion, decision-driving artifacts for C-suite, board, VP, investor, funder audiences. Covers BLUF, board memos, ICP briefs, status reports, decision-readiness. TRIGGER: board memo; exec summary; VP status update; make-this-exec-ready; 30-second pitch; ICP brief; decision memo; one-pager; OKRs; proposals & grants; pitch/investor decks; negotiation & persuasion prep; public speaking/presentations; keynote/conference talk; all-hands/speech writing; founder & shareholder letters; whitepapers; customer case studies. SKIP: general prose/voice/editing → writing-expert; software/product/eng docs → technical-writing-craft; marketing/PR/launch/newsletters → content-and-marketing-writing; career/academic/legal/policy/survey → career-and-formal-writing; analysis to build first → da-applied-and-communication.
name: executive-comms
version: 1.4.0
updated: "2026-06-01"
category: custom
related_skills:
  - writing-expert
  - technical-writing-craft
  - content-and-marketing-writing
  - career-and-formal-writing
whenToUse:
  - Draft a board memo or executive summary for C-suite or VP sign-off
  - Write a decision memo where the exec must approve, reject, or fund
  - Format a status update with red/yellow/green ratings for leadership
  - Create an ICP brief, one-pager, or headline deck for a VP or board audience
  - Compress a long document to a 90-second exec brief
  - Write OKRs, a grant or funding proposal, a pitch/investor deck, a founder letter, a whitepaper, or a customer case study
  - Prepare negotiation, persuasion, public-speaking, presentation, or speech material
whenNotToUse:
  - General prose, voice, tone, or editing without an audience change (use writing-expert)
  - Software, product, or engineering docs — specs, PRDs, RFCs, runbooks, API docs (use technical-writing-craft)
  - Marketing, PR, newsletters, launch or announcement copy (use content-and-marketing-writing)
  - Career, academic, legal, policy, or survey writing (use career-and-formal-writing)
---

# Executive Communications

Executive audiences differ from general readers in three ways: their attention budget is measured in seconds, not minutes; every document must carry a decision or an ask; and every material claim requires a number. A memo that would satisfy a manager will frustrate a VP. This skill covers the structures, frameworks, and anti-patterns that separate exec-ready writing from everything else.

See frontmatter `whenToUse` / `whenNotToUse` for routing rules.

---

## Sub-skill routing table

This hub consolidates 10 executive-comms sub-skills as on-demand reference files. When a task matches a row, **Read the listed `references/` file** before answering — do not rely on this table alone for depth.

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| One-pager writing | Craft single-page summaries — exec briefings, partner briefs, sales-enablement one-pagers, account briefs, project pitches, product one-pagers | `references/one-pager-writing.md` |
| OKR writing | Objectives and Key Results craft — the ambitious-but-not-impossible rule, measurable-KR test (numerator + denominator + deadline), OKR vs KPI distinction | `references/okr-writing.md` |
| Proposal & grant writing | Funder-aligned grant and proposal writing — problem statement, demonstrated need, theory of change, logical frameworks (logframes), SMART objectives | `references/proposal-and-grant-writing.md` |
| Pitch deck writing | Slide-by-slide narrative-arc craft for fundraising, sales, and partner pitch decks — the Guy Kawasaki 10/20/30 rule and standard slide order | `references/pitch-deck-writing.md` |
| Negotiation & persuasion | Tactical layer for persuasion-shaped writing and verbal-negotiation prep — Cialdini's 6 principles applied to written requests, Chris Voss tactics | `references/negotiation-and-persuasion.md` |
| Public speaking & presentations | Presentation and public-speaking craft — slide-doc vs PowerPoint structure, delivery mechanics, audience management | `references/public-speaking-and-presentations.md` |
| Speech writing | Write speeches meant to be spoken aloud — keynotes, all-hands addresses, conference talks, panel intros, toasts, departure and commencement speeches | `references/speech-writing.md` |
| Founder letter writing | Annual shareholder letters, founder updates, and "state of the company" letters in the Buffett / Bezos / Stripe lineage | `references/founder-letter-writing.md` |
| Whitepaper writing | Long-form B2B thought-leadership and lead-generation whitepaper craft — the problem → solution → proof → CTA arc, executive-summary discipline | `references/whitepaper-writing.md` |
| Case study writing | Customer success story craft — the Challenge / Solution / Result template, customer-quote integration, anonymization decisions | `references/case-study-writing.md` |

---

## Output modes

When this skill is invoked, clarify the situation with one targeted question if ambiguous, then apply the matching mode:

| User provides | Default output mode |
|---|---|
| A topic or ask with no draft | Draft the full document in the appropriate format |
| An existing draft | Annotate inline using `[EXEC-NOTE: <finding>]` markers keyed to the decision-readiness checklist, then produce a clean rewritten version below the annotated draft |
| A format name only ("ICP brief", "board memo") | Produce a filled template with placeholder copy showing correct structure |
| "Critique this" or "make this exec-ready" | Run the decision-readiness checklist, return numbered findings, then produce a revised draft |

**Clarifying-question policy:** if the user's input is ambiguous (unclear audience, unclear decision required, unclear format), ask exactly one targeted question before drafting. Do not ask multiple questions at once.

**Self-check before delivering:** re-read the output and confirm it passes the decision-readiness checklist. If any item is unchecked, fix it before responding.

---

## Drafting sequence (from-scratch drafts)

Apply this sequence in order:

1. **Select the format** using the Format selection guide below.
2. **Write the BLUF** — one sentence, decision or ask, with a number.
3. **Draft the body** following the chosen format's structure (ICP, board memo, etc.).
4. **Run the decision-readiness checklist** before outputting. Fix any unchecked items.

---

## Format selection guide

Choose the format before drafting. The decision depends on read time available and action required.

| Situation | Best format | Why |
|---|---|---|
| Exec must approve budget or headcount | ICP brief (1 page) | Fastest path to yes/no; forces issue-context-proposal discipline |
| Board needs quarterly or strategic update | Board memo (3-10 pages + appendix) | Full record; supports fiduciary review |
| Meeting needs visual support | Headline deck (action-titled slides) | Scannable without narration |
| Async status update to VP | What/So What/Now What email (3 paragraphs) | No meeting required; decision at the top |
| Exec asks "where do we stand" mid-meeting | 30-second compression | Verbal; 4-step formula |

**Email as an exec format:** VP-level email follows the same BLUF rule as a memo. Subject line = the decision or ask. First sentence = the BLUF. Body = context + supporting data. Signature = next step with owner and date. Never bury the ask in paragraph 3.

---

## Audience model

### Attention budgets

| Format | Read time | What execs actually read |
|---|---|---|
| 90-sec brief | 90 s | Headline + first sentence of each section |
| 5-min readout | 5 m | Headline + section heads + bullets |
| 30-min deep dive | 30 m | Full doc including appendix |

Design for the shortest format the doc may encounter. A board packet will be read in 90 seconds by half the room regardless of length.

### What execs prioritize (in order)

1. Decision required: yes/no, approve/reject, fund/defer
2. Ask: what do you need from me, by when
3. Risk + mitigation: what can go wrong, what stops it
4. Quantified impact: revenue, headcount, time, cost
5. Context: only after the above four are satisfied

If items 1-4 are not in the first two sentences, the exec reads past them looking for something else or stops reading.

---

## Frameworks

### BLUF (Bottom Line Up Front)

The first sentence of every exec document is the decision, status, or ask, stated directly. Everything that follows supports it.

**Rule:** One sentence. Subject + verb + number or outcome. No qualifications in the first sentence.

**Bad:** "This memo summarizes the findings of our Q2 review of the enterprise renewal pipeline and explores several possible paths forward for the executive team's consideration."

**Good:** "We need $2.4M approved by June 15 to close 3 enterprise renewals at risk of churn."

The bad example tells the exec what the memo is about. The good example gives them the decision. All supporting detail (which accounts, why they are at risk, what the mitigation is) follows in the body.

**Bad:** "The platform migration has experienced some challenges related to timeline."

**Good:** "Platform migration is 6 weeks late; go-live moves from July 1 to August 12 unless we add 2 engineers this sprint."

---

### Headline-deck conventions

Every slide title is an action sentence, not a noun phrase. The title carries the conclusion; the body carries the evidence.

**Noun-phrase title (bad):** "Customer X Renewal Status"

**Action title (good):** "Renew Acme Corp now — $1.8M ARR at risk if we miss June 30"

**Noun-phrase title (bad):** "Engineering Capacity Overview"

**Action title (good):** "We are 4 engineers short of shipping v2 on schedule"

A deck of action titles can be read as a stand-alone executive summary. A deck of noun-phrase titles cannot.

---

### ICP brief (Issue / Context / Proposal)

Three paragraphs, each required, none exceeding 5 sentences.

**Paragraph 1 — Issue:** State the problem in one or two sentences. Quantify it. Identify who owns it.

**Paragraph 2 — Context:** Provide the 3-4 facts an exec needs to evaluate the proposal. Not background history; only facts that change which option is correct.

**Paragraph 3 — Proposal:** State the recommended action, the owner, the date, and the cost or resource ask.

**Example:**

> *Issue:* Our largest enterprise segment (38% of ARR) renews on a 12-month cycle; 14 accounts representing $6.2M expire in Q3 with no renewal motion started.
>
> *Context:* Average time-to-close for enterprise renewals is 67 days. Q3 close date is August 31. As of today, June 1, we have 91 days — enough to save all 14 if we start now. The renewal team is currently at 60% capacity; we have budget headroom of $140K.
>
> *Proposal:* Approve hiring 2 contract renewal managers at $70K each for 90 days, starting June 15. Owner: VP Sales. Expected outcome: protect $6.2M ARR at a cost of $140K (2.3% of at-risk revenue).

---

### Board memo structure

| Section | Length | Purpose |
|---|---|---|
| Cover page | 1 page | Topic, date, presenter, classification |
| Executive summary | 1 page | BLUF + 3-5 key findings + recommendation |
| Body | 3-8 pages | Problem / Options / Recommendation with data |
| Appendix | Unlimited | Supporting data, methodology, prior context |

The executive summary must stand alone. Board members who read only the first page must be able to vote or defer with full understanding. Move anything that does not help a decision into the appendix.

---

### "Tell me in 30 seconds" recursive compression

Use this four-step reduction when an exec asks for the short version mid-meeting or when testing whether a draft is decision-ready.

1. **Extract the decision.** What single yes/no is required?
2. **Name the ask.** Who needs to do what?
3. **Attach one number.** What is the cost, risk, or gain?
4. **Name the risk.** What happens if we do nothing or decide wrong?

If you cannot complete all four steps in 30 seconds, the document is not exec-ready. Rewrite the opening.

**Before:** A 400-word memo about migration delays, team capacity, vendor contract terms, and risk mitigation options.

**After (30-second version):** "Migration is 6 weeks late. I need 2 engineers reassigned from Project X by Friday or go-live slips to August 12, costing $180K in delayed revenue recognition."

---

### What / So What / Now What

Use this pattern for status updates, incident readouts, and recurring reporting.

| Layer | Question | Example |
|---|---|---|
| What | What happened? (fact, number, event) | "Support ticket volume rose 40% week-over-week." |
| So What | Why does it matter to this audience? | "At this rate, SLA breach risk appears in 8 days." |
| Now What | What action is required, by whom, by when? | "We are adding 4 agents on Monday; I need your sign-off by EOD today." |

Never deliver a "What" without a "So What." Status updates that report facts without consequences force the exec to do analysis they expect you to have done.

---

### Ladder of Inference (Argyris)

The Ladder of Inference describes how reasoning moves from raw data to action. Executives often receive documents that start at the top (conclusions, proposed actions) without showing the rungs. This produces distrust or second-guessing.

| Rung | Label | What to show the reader |
|---|---|---|
| 1 | Observable data | What actually happened, unfiltered |
| 2 | Selected data | What you chose to focus on and why |
| 3 | Interpreted meaning | What you think it means |
| 4 | Assumptions | What you took for granted |
| 5 | Conclusions | What you decided is true |
| 6 | Beliefs | The generalized view you now hold |
| 7 | Actions | What you propose to do |

Ground every recommendation explicitly in layer 1 and name the selection criteria at layer 2. When you skip rungs, add a bridging sentence: "We focused on accounts over $500K because smaller accounts have a separate renewal motion."

**Example that skips rungs (bad):** "We recommend exiting the SMB segment."

**Example that shows rungs (good):** "SMB accounts (< $50K ARR) represent 42% of ticket volume but 6% of ARR. Support cost per dollar of SMB revenue is 4x enterprise. We recommend pausing new SMB acquisition while we evaluate a self-serve model."

---

## Status report patterns

### Red / Yellow / Green semantics

| Color | Means | Obligates the writer to |
|---|---|---|
| Red | Decision or escalation needed now | Name the ask, the date, and the consequence of inaction |
| Yellow | On track but risk present | Name the specific risk, the mitigation, and the owner |
| Green | On track, no action needed | One sentence of confirmation; no narrative needed |

**Common red/yellow mistakes:**

- Marking yellow when the project is red to avoid bad news. Yellow means "I have a plan for the risk." No plan means red.
- Writing a green update with a paragraph of caveats. Caveats belong in yellow. Green is a clean signal.
- Red status without an ask. Every red item must say: "I need [specific thing] from [specific person] by [specific date]."

---

## Quantification discipline

- Every material claim pairs with a number, or is explicitly marked qualitative.
- Use the actual number, not "significant," "substantial," or "meaningful."
- Prefer rates and absolutes together: "40% increase (from 50 to 70 tickets/day)" rather than "40% increase."
- Where a number is not yet available, write: [TBD by Owner on Date]. Do not omit the placeholder.
- One decimal place is usually enough. Precision beyond that signals false confidence.
- **Numbers must come from the user's source data.** Do not invent figures to satisfy the quantification rule. If a number is unavailable, use the [TBD] placeholder and flag it explicitly.

**Bad:** "We have seen strong adoption of the new feature."

**Good:** "Feature adoption reached 62% of active accounts in the first 30 days, up from a 45% target."

---

## Anti-patterns

| Anti-pattern | Problem | Fix |
|---|---|---|
| Bullets without verbs ("Customer satisfaction") | Noun bullets carry no claim | Add a verb and a number: "Increased CSAT 15 points to 87" |
| "As you may know" | Condescending filler | Delete the phrase; make the claim directly |
| "It is worth noting" | Passive hedging | Note it directly with a subject and verb |
| Multiple recommendations buried in body | Exec may read only the summary and miss them | Promote to the BLUF; one primary recommendation per document |
| Hedging adverbs (potentially, possibly, perhaps) | Signals you have not committed to a view | Quantify the uncertainty or cut the hedge |
| Jargon density above 5% | Forces exec to decode instead of decide | Define on first use; replace with plain language where possible |
| Passive voice in recommendation sentences | Hides the owner | "The team should consider..." becomes "VP Eng decides by Friday" |
| Abstract risk statements ("there may be challenges") | Gives exec nothing to act on | Name the specific risk, the probability, and the dollar or time impact |

---

## Decision-readiness checklist

Before sending any exec communication, confirm yes to each item:

- [ ] Can the reader say yes or no in 90 seconds?
- [ ] Is the ask stated in the first two sentences?
- [ ] Is the decision-blocker (data gap, approver, budget authority) named explicitly?
- [ ] Is the next-step owner named with a date?
- [ ] Does every material claim have a number, or is it explicitly marked qualitative?
- [ ] Are all recommendations in the opening, not buried in the body?
- [ ] Is red/yellow/green status accompanied by its obligated disclosures?
- [ ] Have hedging adverbs and filler phrases been removed?

If any item is unchecked, revise before sending.

---

## Composition workflow

Use these related skills in sequence for highest-quality output:

1. **Draft:** use `writing-expert` with exec audience framing. Apply BLUF, the format selection guide, and the decision-readiness checklist from this skill as the drafting constraints.
2. **Self-check:** re-read the draft against the decision-readiness checklist above. Fix any unchecked item before proceeding.
3. **Review:** use `document-critique` with the decision-readiness checklist as the rubric.
4. **Persuasion framing:** if the memo needs to move a skeptical audience, load `references/negotiation-and-persuasion.md` for the tactical layer, or `writing-expert` for deep ethos/logos/pathos rhetorical sequencing.

---

## The 60-second elevator pitch (executive-audience variant)

The executive elevator pitch converts a 60-second informal moment into a calendar follow-up, a budget approval, or a decision. Structure (total ~150 words):

| Beat | Time | Content |
|------|------|---------|
| Decision frame | 5 s | The single yes/no you want from the listener |
| Status | 10 s | Where things stand today, with one number |
| Risk or stake | 10 s | What goes wrong if no decision is made |
| Path forward | 15 s | The recommendation, in active voice |
| Ask | 10 s | Approval, sign-off, budget, owner, date |
| Margin | 10 s | Silence for response |

**Worked example (TAM brief to VP of Customer Success):**

> "I need 10 minutes on your calendar this week to walk through the Acme renewal. The renewal is August 31, $2.4M ARR. Acme's CSAT dropped from 78 to 51 after three S1 incidents. Without Director-level outreach, churn risk is 60% — we lose the account and the renewal motion stalls into Q4. I want to propose a 90-day recovery plan: weekly exec sync, doubled TAM hours, and one engineering credit pool. I need your approval on the engineering credit and a 30-minute Acme exec call you'd attend in week two."

That is 110 words and lands every beat.

**Key rules:** Decision frame is first. End with a specific ask ("I need your approval on X"), never a soft close. One number per beat. If the moment shrinks, drop to the 30-second compression in the "Tell me in 30 seconds" section above.

---

## References

1. Minto, Barbara. *The Pyramid Principle: Logic in Writing and Thinking*. Pearson Education, 1987. — Origin of BLUF and top-down document structure for business writing.
2. Argyris, Chris. "Teaching Smart People How to Learn." *Harvard Business Review*, May-June 1991. — Source of the Ladder of Inference and grounding conclusions in observable data.
3. US Army. *Mission Command: Command and Control of Army Forces*, Field Manual FM 6-0. Headquarters, Department of the Army, 2003. — Institutional source of BLUF as a writing doctrine for high-stakes, time-constrained communication.

---

## TL;DR conventions in executive communications

The label signals register; the rule (lead with the conclusion) is universal.

| Label | Register | When to use |
|---|---|---|
| BLUF | Military, defense, government, enterprise corporate | High-stakes operational decisions |
| Bottom Line | US corporate, board-facing, finance | Neutral default for mixed boards |
| TL;DR | Tech-native, internal engineering memos | IC-mixed channels; signals "I respect your time" |
| Executive Summary | Formal documents (board packets, M&A decks, regulator submissions) | When the summary will circulate separately |

**Format rules (3 sentences max):**
1. Recommendation — not the topic. "We should approve $4M in capex for the European DR build-out."
2. Cost or trade-off. "The migration costs $1.8M one-time but saves $4.2M/year."
3. What is required of the reader, by when. "We need a go/no-go from CFO and CTO by 30 June."

**Placement:** Top of the document, before any heading. Bold or visually set apart. One per document. 75-word cap — anything longer is an Executive Summary.

**Worked example:**

Weak: *TL;DR: This document presents an analysis of our infrastructure costs and proposes a path forward.*

Strong: *Bottom Line: We should consolidate from three cloud providers to one (AWS) over 9 months. The migration costs $1.8M one-time but saves $4.2M/year and unifies SRE on-call. We need a go/no-go from CFO and CTO by 30 June to start contracting.*

---

## Bury the lede — the cardinal sin of executive comms

An exec who cannot extract the central point from sentence one of paragraph one will not read paragraph two. Two cases:

1. **Accidental bury** — the writer opens with context, history, or housekeeping; the recommendation is in paragraph three. Always wrong.
2. **Deliberate delayed lede** — a narrative opening that earns its delay through tension. Legitimate only in crisis comms (empathy before facts) and leadership essays (founder letters, keynotes) where the narrative is the argument.

**Accidental bury disguises to cut on sight:**
- "As you know, the platform team has been working on…" — throat-clearing tell
- "Per our discussion in last week's offsite, and consistent with Q1 priorities…" — attribution before content
- "Beginning in Q3, the team began evaluating…" — chronology lede
- "Thank you for your support of the initiative…" — gratitude before the answer

**Worked examples:**

Accidental bury: *"Q3 has been a productive quarter for the platform team, with multiple initiatives running in parallel — including a significant evaluation of our legacy infrastructure footprint. As you know, the legacy cache has been a topic of recurring discussion across architecture reviews and on-call retrospectives. After consultation with the SRE and product teams, and consistent with our reliability priorities for the year, we have arrived at a decision."* — 80 words, still no decision.

BLUF rewrite: *"Recommendation: deprecate the legacy cache by end of Q4. Migration plan attached; sign-off requested by Friday."*

Deliberate delayed lede (layoff context): *"Last week we asked 8% of the company to leave. That decision was mine, and it was harder than any I've made in seven years of running this team. Here's how we got here, what it means for the work ahead, and what comes next."* — each sentence earns its place.

**3-test diagnostic** — read paragraph one. Can a reader state: (1) the recommendation, (2) what is being asked of them, (3) what changes if they approve or reject? Pass = 3/3. Restructure before sending if 0–1.

<!-- cross-hub-map -->
## Cross-hub map — where every writing topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `writing-expert` | Writing Expert (prose craft, voice, style, editing) — hub | `references/editing-and-revision.md`, `references/rhetorical-frameworks-deep.md`, `references/storytelling-and-narrative.md`, `references/headline-craft.md`, … |
| `technical-writing-craft` | Technical & Product Writing (docs, specs, engineering comms) — hub | `references/api-docs-craft.md`, `references/howto-writing.md`, `references/tutorial-writing.md`, `references/knowledge-base-authoring.md`, … |
| `executive-comms` | Executive & Business Communication (leadership, persuasion, decks) — hub | `references/one-pager-writing.md`, `references/okr-writing.md`, `references/proposal-and-grant-writing.md`, `references/pitch-deck-writing.md`, … |
| `content-and-marketing-writing` | Content, Marketing & External Comms (PR, newsletters, launch, social) | `references/sales-and-marketing-copy.md`, `references/press-release-writing.md`, `references/crisis-pr-writing.md`, `references/newsletter-writing.md`, … |
| `career-and-formal-writing` | Career, Academic, Legal & Formal Writing | `references/resume-and-cv-writing.md`, `references/cover-letter-writing.md`, `references/job-description-writing.md`, `references/performance-review-writing.md`, … |
