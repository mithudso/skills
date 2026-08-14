<!-- hub-reference-banner -->
> **Reference file — part of the `executive-comms` hub.** Formerly the standalone `pitch-deck-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: pitch-deck-writing
description: Slide-by-slide narrative-arc craft for fundraising, sales, and partner pitch decks. Covers the Guy Kawasaki 10/20/30 rule (10 slides, 20 minutes, 30-point font), the Sequoia Capital / Reid Hoffman pitch deck structure (the canonical Series A template), the Y Combinator seed deck guide, the problem-solution-traction-team-ask arc, the Amazon six-page narrative memo as a deck alternative, "deck vs narrative deck" tradeoff, appendix slide convention, "what I most want you to remember" closer, the no-cliffhanger rule, and the "designed by a designer, written by a writer" collaboration. TRIGGER when user says "write a pitch deck", "investor deck", "sales deck", "Series A deck", "seed deck", "Sequoia template", "Guy Kawasaki", "10/20/30 rule", "fundraising slides", "demo day pitch", "deck narrative", "founder pitch", "partner pitch deck", or asks for slide-by-slide structure of a persuasive presentation. SKIP for one-page executive memos (use executive-comms), long-form whitepapers and analyst reports (use whitepaper-writing), sales/marketing copy that isn't a deck (use sales-and-marketing-copy), PRDs and product specs (use prd-writing or spec-writing), executive comms that aren't pitches (use executive-comms), conference talks and keynotes (use public-speaking-and-presentations), and post-mortems or status reports (use postmortem-writing or writing-expert).
category: custom
tags: [writing, business-writing, pitch-deck, fundraising, presentations]
---

# Pitch Deck Writing

## Overview

A pitch deck is a slide-based artifact that compresses a fundraising, sales, or partnership case into a tight narrative arc — usually 10 to 14 slides — designed to be skimmed in minutes and remembered in seconds. The genre has matured around a small number of canonical structures: Guy Kawasaki's 10/20/30 rule, Sequoia Capital's template (developed by Reid Hoffman and the Sequoia partners), and Y Combinator's seed-deck guidance.

A pitch deck is **not** a slideshow of everything you know. It is a sequenced argument: each slide raises the question the next slide answers. Investors read in scan-mode — average viewing time per slide is reportedly under 10 seconds in the first pass. If a slide's point doesn't land in that window, it doesn't land.

The Amazon counter-tradition: Bezos banned PowerPoint internally in 2004, replacing decks with six-page narrative memos. The narrative form forces clearer thinking ("there is no way to write a six-page narratively structured memo and not have clear thinking"), but the bullet-deck format remains dominant outside Amazon. The "narrative deck" hybrid — dense, full-sentence slides meant to be read silently before discussion — splits the difference.

This skill covers the slide-deck form. For prose memos, fall back to `writing-expert` and `executive-comms`.

## Core Concepts

### 1. The 10/20/30 rule (Kawasaki)

Guy Kawasaki's rule for fundraising decks:
- **10 slides** — the optimal number; humans can't absorb more than 10 concepts in a meeting
- **20 minutes** — the maximum runtime, allowing 40 minutes of a 60-minute slot for discussion
- **30-point font** — minimum; forces brevity and prevents the "read off the slide" anti-pattern

The 30-point rule is the most violated. If you can't say it in 30-pt type, you have too many words on the slide. The audience reads the slide instead of listening to you.

### 2. The canonical Sequoia / Reid Hoffman structure

Sequoia Capital's pitch deck template is the global standard for Series A. It consists of 10–12 slides in a specific order:

1. **Title / Company purpose** — one-line description of what the company does
2. **Problem** — the customer pain, with evidence
3. **Solution** — your product as the answer to that pain
4. **Why now** — timing argument; why this couldn't have been built 5 years ago and can't wait 5 more
5. **Market size** — TAM/SAM/SOM, bottom-up preferred
6. **Competition** — who else is in the space, how you differ
7. **Product** — screenshots, demo, key features
8. **Business model** — how you make money
9. **Team** — why this team can win
10. **Traction** — what you've already proven (metrics, customers, growth)
11. **Financials** — current state and projections
12. **Ask** — how much you're raising, on what terms, what you'll do with it

Slide 4 ("Why now") is the slide most founders skip and the slide most VCs care about. It is where timing risk is addressed.

### 3. The Y Combinator seed deck

Y Combinator's seed-stage variant is shorter (10–12 slides) and tighter than Sequoia's Series A template. The order shifts to put traction earlier when you have it, and the deck is optimized for clarity-over-completeness:

1. **Title** — name, one-line description, contact
2. **Problem**
3. **Solution**
4. **Why now**
5. **Market size**
6. **Product**
7. **Traction** (or "Why we will win" if pre-traction)
8. **Business model**
9. **Competition**
10. **Team**
11. **Financials / Ask**
12. **Vision** (optional closer)

YC's guidance: visually boring is fine. Large text, one idea per slide, no animations. The deck should be readable in 3 minutes.

### 4. The narrative arc — problem → solution → traction → team → ask

Strip away the template specifics and every effective pitch deck follows the same five-beat arc:

1. **Problem** — establish the pain. The reader nods.
2. **Solution** — show your answer. The reader believes the answer is plausible.
3. **Traction** — prove it's working. The reader believes the answer is real.
4. **Team** — prove you can execute. The reader believes you specifically can win.
5. **Ask** — close. The reader knows what you want and what they get.

Everything else (market size, competition, product detail, business model, financials) is supporting evidence for these five beats. When you trim a deck, ask: which beat does this slide serve? Slides that don't map to a beat should be cut or moved to the appendix.

### 5. The appendix convention

Strong decks have an appendix — 5 to 20 extra slides beyond the main 10–12 — that are NOT shown in the live pitch but are sent with the deck for VCs to dig into. The appendix is where you put:

- Detailed financial model
- Customer logos and case studies
- Technology architecture
- Hiring plan
- Detailed competitive analysis
- Cohort retention data
- Press / awards

Convention: appendix slides are clearly labeled (e.g., "Appendix: Unit Economics") so the reader knows they're in supplementary territory. The main deck is for the meeting; the appendix is for follow-up.

### 6. The "no cliffhanger" rule

Pitch decks are NOT mystery novels. Do not save the reveal. Slide 1 should already answer "what does this company do" at headline level. Slide 2 sharpens the problem; slide 3 makes the solution concrete. By slide 3, the reader knows your full thesis.

Anti-pattern: opening with anecdotes, market trends, or origin stories that delay the company description until slide 5. VCs lose patience by slide 3 if they don't know what you do.

### 7. "What I most want you to remember" — the closer

Reid Hoffman's pitch advice: end the deck with one slide that names the single thing you want the audience to walk away remembering. Not a recap of all 10 slides — one sentence, one number, one image. Examples:

- "We are growing 30% month-over-month with $0 paid acquisition."
- "Every Fortune 100 bank is on our waitlist."
- "This is the team that built and sold X to Y."

The closer slide is what survives the next 24 hours in the investor's mental working set. Without it, the deck dissolves into vague impressions.

### 8. Narrative deck vs traditional deck (Amazon style)

Two formats sit at opposite ends of the deck genre:

**Traditional deck (Kawasaki / Sequoia / YC):**
- Big headlines, sparse text
- 30-pt font, ≤10 slides
- Presenter narrates; slides support
- Read in minutes

**Narrative deck (sometimes called "Amazon-style" though Amazon banned decks entirely):**
- Full sentences, dense paragraphs on slides
- Meant to be read silently in the room before discussion
- 6–10 slides but each is essentially a page of prose
- Common in operating reviews, partner meetings, board decks

The Amazon six-pager itself is **not** a deck — it's a memo. But the "narrative deck" hybrid (dense slides read in silence) borrows the same discipline of forcing the author to write complete arguments rather than bullet fragments.

Use the traditional deck for fundraising and sales pitches with a presenter. Use the narrative deck (or a memo) for board reviews, strategy reviews, and internal alignment meetings.

### 9. "Designed by a designer, written by a writer"

The collaboration norm in serious pitch decks: a writer drafts the narrative and the words; a designer renders the slides. Founders often do both and the deck shows it — either the design is amateur and the message is buried, or the design is slick and the message is empty.

The minimum-viable separation:
- **Writer's draft** — slide titles + body copy in a doc (not slides yet). Iterate the narrative.
- **Designer's pass** — render approved copy into slides with consistent visual hierarchy, type, color, and imagery.
- **Founder's read-through** — read every slide title in sequence. The titles alone should tell the story.

If the slide titles don't form a coherent narrative, the deck has a structural problem no design can fix.

### 10. The slide-title-as-thesis rule

Every slide has a title. The title is a thesis, not a category. Compare:

- **Category title (weak):** "Market Size"
- **Thesis title (strong):** "$47B TAM growing 22% annually; we serve the underserved bottom 60%"

The reader who skims only the titles should still get the argument. This forces you to know what each slide is actually saying.

## Templates and Examples

### Sequoia / Series A deck — slide-by-slide

| # | Slide | Thesis-title example | Body content |
|---|---|---|---|
| 1 | Title / Purpose | "Acme — the operating system for indie retailers" | Company name, one-line description, presenter, date |
| 2 | Problem | "Indie retailers lose 8% of revenue to inventory errors" | Customer pain with evidence — quote, stat, support volume |
| 3 | Solution | "One scan; the cloud reconciles the rest" | Product as the answer to slide 2 |
| 4 | Why now | "Cheap POS hardware + computer vision finally make this affordable" | Timing / enabling-tech argument |
| 5 | Market size | "$47B TAM, 320K stores in the US alone" | TAM/SAM/SOM, bottom-up |
| 6 | Competition | "Big POS players ignore < 5-store chains; we own this segment" | 2x2 or comparison chart |
| 7 | Product | "Setup in 12 minutes; first reconciliation in 1 hour" | Screenshots, demo flow, top 3 features |
| 8 | Business model | "$99/store/month + 0.3% transaction fee" | Pricing, ACV, expansion mechanics |
| 9 | Team | "Founders shipped this exact problem at $LASTCO" | Photos + 1-line credentials each |
| 10 | Traction | "412 stores, $2.1M ARR, 14% MoM growth" | Logos, metrics chart, customer quotes |
| 11 | Financials | "Path to $20M ARR by end of next year" | Top-line projections + key drivers |
| 12 | Ask | "$8M Series A to scale GTM in North America" | Amount, use of funds in 3–4 buckets |
| (closer) | What to remember | "412 stores, 14% MoM, zero churn." | One number / one line |
| (appendix) | Various | — | Unit economics, cohort retention, hiring plan, etc. |

### Y Combinator seed deck — slide-by-slide

| # | Slide | Body content |
|---|---|---|
| 1 | Title | Company name + one-line + contact |
| 2 | Problem | Customer pain with one piece of evidence |
| 3 | Solution | Your product, plainly stated |
| 4 | Why now | Timing / enabling shift |
| 5 | Market size | TAM, one number, source |
| 6 | Product | Hero screenshot + 3 bullet features |
| 7 | Traction | Best metric you have, even if small |
| 8 | Business model | One-line of how money flows |
| 9 | Competition | 2x2 or table — your differentiation |
| 10 | Team | Founder photos + line of credentials |
| 11 | Ask | Round size, use of funds |
| 12 | Closer | The one thing |

### Sales deck — slide-by-slide

Adjust the arc for selling a product to a buyer (not selling equity to an investor):

1. **Title** — your company + buyer's name
2. **Their world today** — the buyer's status quo (acknowledge it; this is the problem slide reframed)
3. **The shift** — what's changing in their market / regulations / customer expectations
4. **What winners are doing** — how leaders in their space are responding
5. **Why most fail** — the trap of incremental approaches
6. **Our solution** — your product as the way to win
7. **How it works** — 3-step concept, not feature list
8. **Proof** — customer logos, case study, quantified results
9. **Investment** — pricing, packaging
10. **Next step** — pilot, POC, signature

This "Challenger" sale structure (popularized by Andy Raskin and others) borrows the same five-beat arc but reframes Slide 2 around the buyer's reality rather than the seller's problem statement.

### Pitch deck draft as prose (writer's pass)

Before designing slides, write the deck as a doc. Each slide gets a heading and 2–4 sentences:

```markdown
## Slide 1 — Title
Acme — the operating system for indie retailers. Mitchell Hudson, May 2026.

## Slide 2 — Problem
Indie retailers (320K stores in the US) lose 8% of revenue to inventory
errors. Existing POS systems cost $40K+ and require dedicated IT. Owners
told us in 32 of 32 interviews: "I know my stock is wrong; I just can't
afford to fix it."

## Slide 3 — Solution
Acme is a $99/month app that scans your shelf with any phone camera and
reconciles inventory against your POS in 60 seconds. First scan in 12 minutes
from install.

## Slide 4 — Why now
Phone cameras hit the resolution / low-light threshold in 2024. Compute
cost for real-time inventory CV dropped 80% since 2023. We can do for $99
what cost $40K three years ago.

[...]
```

Iterate this doc until the slide titles alone tell the story. Then design.

## Anti-patterns

1. **The mystery opening** — first 3 slides don't say what the company does. Investors give up by slide 3. *Fix:* one-line description in slide 1 title; slide 2 problem; slide 3 solution.

2. **Wall-of-text slides** — paragraphs of body copy presented live. Audience reads instead of listening; presenter is redundant. *Fix:* 30-pt font rule; if it doesn't fit, simplify or split.

3. **No "why now"** — deck skips the timing argument. VCs assume the idea was already tried and failed. *Fix:* slide 4 explicitly addresses the enabling shift.

4. **No traction (when you have some)** — burying real customer / revenue / growth data deep in the deck or appendix. *Fix:* if traction is your strongest signal, move it to slide 4 or 5 — even before market size.

5. **Generic team slide** — "20+ years experience" instead of specific credentials tied to this problem. *Fix:* name companies, name outcomes ("scaled X from $1M to $50M ARR at LastCo").

6. **Competition denial** — "we have no competitors." VC reads this as "you don't understand the market." *Fix:* every market has competitors including the status quo and "do nothing." Name them.

7. **TAM inflation** — "$2T global market" with no bottom-up math. Triggers immediate skepticism. *Fix:* bottom-up TAM = (# customers) × (ACV) with both numbers defended.

8. **No ask** — the deck ends with vague "let's talk." VC doesn't know what you want. *Fix:* explicit slide stating round size, use of funds in 3–4 buckets, and milestones the round will hit.

9. **Designed-but-not-written** — gorgeous slides with vague copy. Common when the founder is design-fluent. *Fix:* write the prose first; design last.

10. **Identical to template** — Sequoia template, default fonts, default colors. Forgettable. *Fix:* use the template as scaffolding; the copy and visual identity are yours.

11. **Cliffhanger close** — last slide is "Thank you" or "Questions?" Wasted prime real estate. *Fix:* "What I most want you to remember" slide with one number / one line.

12. **No appendix** — every supporting detail is in the main deck OR no supporting detail exists for follow-up. *Fix:* main deck = 10–12 slides for the meeting; appendix = the deep dive that VCs ask for.

## Decision Heuristics

| Situation | Use |
|---|---|
| Seed round to investors | YC seed structure, 10–12 slides |
| Series A to investors | Sequoia / Reid Hoffman structure, 12 + appendix |
| Series B+ to investors | Sequoia structure + traction-heavy, financials promoted |
| Enterprise sales to a buyer | "Challenger" sales arc (status quo → shift → solution) |
| Partner / BD meeting | Sales arc but lighter on price; heavier on integration value |
| Internal board review | Narrative deck (dense slides) or six-pager memo |
| Demo day (3-minute pitch) | 10 slides MAX, 30-pt font, problem→solution→traction→ask only |
| Conference keynote | Different genre — use `public-speaking-and-presentations` |

**Deck vs memo test:** is there a live presenter? Yes → deck. No → memo. Mixed (deck circulated for async review) → narrative deck or supplement deck with a 2-page summary.

**Slide-count test:** if you can't tell the story in 12 slides, the story isn't clear yet. Cut, don't add.

**Title-only test:** read only the slide titles in sequence. Do they tell the story? If no, the structure is broken.

**Closer test:** if your audience remembered only ONE slide tomorrow, which one would you want it to be? That slide goes last (or second-to-last before the ask).

## Cross-references

- **Adjacent skills:** `executive-comms` (single-stakeholder exec messaging), `sales-and-marketing-copy` (ad copy, landing pages), `public-speaking-and-presentations` (delivery and conference talks), `writing-expert` (prose craft for slide copy), `storytelling-and-narrative` (narrative-arc theory).
- **Upstream:** `audience-modeler` to define who reads the deck and what they care about.
- **Companion:** for the narrative-deck/memo variant, pair with `executive-comms` and `writing-expert`.

## References

1. Guy Kawasaki, "The 10/20/30 Rule of PowerPoint" — the canonical 10 slides / 20 minutes / 30-point-font rule. https://guykawasaki.com/the_102030_rule/
2. Sequoia Capital, "Writing a Business Plan" / pitch deck template — the canonical Series A structure (Title, Problem, Solution, Why Now, Market, Competition, Product, Business Model, Team, Financials). https://www.sequoiacap.com/article/writing-a-business-plan/
3. Y Combinator Startup Library, "How to build your seed round pitch deck" — Michael Seibel's seed-deck guidance. https://www.ycombinator.com/library/2u-how-to-build-your-seed-round-pitch-deck
4. Y Combinator Startup Library, "How to build a great Series A pitch and deck." https://www.ycombinator.com/library/8d-how-to-build-a-great-series-a-pitch-and-deck
5. CNBC, "Why Jeff Bezos makes Amazon execs read 6-page memos at the start of each meeting" — Amazon's narrative-memo counter-tradition to PowerPoint. https://www.cnbc.com/2018/04/23/what-jeff-bezos-learned-from-requiring-6-page-memos-at-amazon.html
6. Reid Hoffman / Greylock partners — public pitch-deck writeups including the LinkedIn Series B deck annotated by Hoffman. https://www.reidhoffman.org/linkedin-pitch-to-greylock/

## When to Use This Skill

Use when the user asks to: write or critique a pitch deck (investor, sales, partner), choose between deck structures (Sequoia vs YC vs Kawasaki), sequence slides into a narrative arc, write slide titles as theses, structure an appendix, write a "what to remember" closer, decide between traditional deck and narrative deck, or coordinate writer/designer collaboration on a deck.

Skip when the user asks for: one-page executive memos (use executive-comms), long-form whitepapers (use whitepaper-writing), product requirements (use prd-writing), engineering specs (use spec-writing), sales/marketing ad copy outside a deck (use sales-and-marketing-copy), conference keynote design (use public-speaking-and-presentations), or post-mortems and status reports (use postmortem-writing or writing-expert).
