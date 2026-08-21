<!-- hub-reference-banner -->
> **Reference file — part of the `content-and-marketing-writing` hub.** Formerly the standalone `sales-and-marketing-copy` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: sales-and-marketing-copy
description: "Sales and marketing copywriting craft — landing pages, sales emails, ads, product positioning, conversion-focused prose. Covers AIDA, PAS, Schwartz's 5 awareness stages, above-the-fold hierarchy, email subject-line patterns, CTA microcopy, cold/nurture/breakup email skeletons, headline formulas, social proof typography, and common anti-patterns. TRIGGER: 'write a landing page', 'sales email', 'subject line', 'CTA copy', 'product positioning', 'value prop', 'above-the-fold', 'PAS framework', 'AIDA', 'conversion copy', 'ad copy'. SKIP: technical docs (use technical-writing-craft); executive memos (use executive-comms); doc review (use document-critique); negotiation (use negotiation-and-persuasion)."
version: "1.1.0"
updated: "2026-05-29"
category: custom
whenToUse:
  - "Write a landing page for my product"
  - "Help me write a cold sales email"
  - "What's a good subject line for this email"
  - "Write CTA copy for my signup button"
  - "How do I position my product against competitors"
  - "Rewrite my value prop to be more compelling"
  - "Apply the PAS framework to this pitch"
  - "My above-the-fold copy isn't converting"
  - "Write an AIDA sequence for this product"
  - "What headline formula should I use for this ad"
related_skills:
  - writing-expert
  - storytelling-and-narrative
  - rhetorical-frameworks-deep
  - negotiation-and-persuasion
---

# Sales and Marketing Copy

Reference for conversion-focused writing: landing pages, emails, ads, CTAs, and positioning copy. Sources: Schwartz *Breakthrough Advertising* (1966), Ogilvy *Confessions of an Advertising Man*, Sugarman *The Adweek Copywriting Handbook*, Cialdini *Influence*, Halbert/Carlton copy classics.

## How to use this skill

When invoked, follow this sequence before producing any copy:

1. **Identify the awareness stage** (Section 2) of the target audience. If unclear, ask: "Who is the primary reader — do they already know they have this problem, or are they discovering it for the first time?"
2. **Identify the conversion goal** — one goal per piece (landing page demo request, email reply, ad click). If multiple goals are stated, flag the conflict and recommend choosing one.
3. **Choose the framework** — AIDA for structured persuasion flows; PAS for pain-first framing; headline formulas for ad/email entry points.
4. **Produce copy or a skeleton** — for pages, produce a full section-by-section skeleton with placeholder copy and notes on what each section must accomplish. For emails, produce a complete draft. For subject lines or CTAs, produce 3–5 options with rationale.

When the request is ambiguous, ask one targeted question before producing copy. Do not ask more than one question per turn.

---

## 1. Classical Sequences

### AIDA — Attention / Interest / Desire / Action

The foundational direct-response sequence. Each stage earns permission to proceed to the next.

| Stage | Job | Failure mode |
|---|---|---|
| **Attention** | Stop the scroll or open the envelope | Headline too clever, benefit buried |
| **Interest** | Hold attention with relevance to the reader's world | Generic problem statement, no specificity |
| **Desire** | Build want — features converted into felt outcomes | Feature list instead of transformation |
| **Action** | Make the next step obvious and low-friction | Vague CTA ("Learn more"), multiple competing asks |

**Practical note:** AIDA is a *reader state* sequence, not a section-by-section template. A long-form sales page may cycle through the sequence multiple times before the final CTA.

#### AIDA — extended craft

**Origin and durability.** AIDA was articulated by E. St. Elmo Lewis in 1898 and codified across the early 20th century. Its survival across 125+ years of advertising — from print to direct mail to digital — is because it maps to a reader's decision pipeline, not to any medium. If you change channels but keep humans on the other end, AIDA still applies.

**Stage transitions are the failure points.** Most weak copy doesn't fail *inside* a stage — it fails at a *transition*. Attention earned by a hook collapses when Interest opens with a feature list instead of relevance. Desire earned by transformation language collapses when Action is a vague "Learn more" instead of a concrete next step. Audit transitions specifically:

- **A → I transition.** The line immediately after the headline must reference the reader's world, not the product. If your subhead names the product before the problem, you have lost the transition.
- **I → D transition.** The pivot from "this problem is real" to "and here's the better future" must happen on the strongest proof point you have — a number, a name, a before/after.
- **D → A transition.** The CTA must be the *easiest* next step, not the *biggest* one. "Book a 15-minute call" beats "Schedule a demo" beats "Contact sales."

**Worked example — landing page hero.**

- *Attention:* "Your build queue is 20 minutes long. Ours is 90 seconds." (Specific, comparative, names the pain.)
- *Interest:* "Every push triggers a fresh container in under a second. No warm-up, no cold start, no shared runners." (Relevance — addresses the named pain with mechanism.)
- *Desire:* "Teams shipping on Vercel moved their CI from 18 min to 90s without changing a line of YAML." (Transformation, specific, social proof embedded.)
- *Action:* "Connect a repo — first 100 builds free." (Concrete next step, low friction, named outcome.)

**When to break it.** AIDA breaks down at two extremes:

- *Stage 5 (Most-Aware) audiences* don't need A/I/D — they need only Action. A landing page for an audience who already knows your product and price should skip to the offer. Inserting Interest and Desire stages to a ready buyer feels like manipulation.
- *Brand-only campaigns* (no conversion goal) don't follow AIDA — they aim at recall, not action. Coca-Cola spots are not AIDA. If you have no conversion goal, do not force AIDA scaffolding.

**Diagnostic — the four-question audit.**

1. Open the piece. Read only the opening 2 seconds' worth. Did it earn the next 2 seconds? (Attention)
2. Read the next 10 seconds. Did it name the reader's world specifically enough that they think "this is for me"? (Interest)
3. Read the body. Is there a concrete transformation — number, before/after, named customer — not just a feature list? (Desire)
4. Look at the CTA. Is the next step concrete, low-friction, and obviously the *one* thing to do? (Action)

A piece that scores 4/4 ships. 3/4 needs a transition rewrite. 2/4 or below needs a structural rewrite.

**References (AIDA-specific).**

- Lewis, E. St. Elmo. Original AIDA formulation, 1898. See Wikipedia, "[AIDA (marketing)](https://en.wikipedia.org/wiki/AIDA_(marketing))" for primary-source lineage.
- [The AIDA Model — Smart Insights](https://www.smartinsights.com/traffic-building-strategy/offer-and-message-development/aida-model/) — modern application across digital channels.
- [AIDA Copywriting Framework — SaaS Funnel Lab](https://www.saasfunnellab.com/essay/aida-copywriting-framework/) — 2025 SaaS-specific examples and stage-transition diagnostics.

---

### PAS — Problem / Agitate / Solve

Pain-first frame. Effective when the audience already feels the problem but has not connected it to your solution.

1. **Problem** — Name the problem precisely. Specificity signals that you understand the reader's world. ("Your email open rates dropped 40% after iOS 15.")
2. **Agitate** — Deepen the pain. Enumerate the downstream costs: lost revenue, wasted time, damaged relationships, compounding risk. Do not manufacture fear; surface consequences the reader already suspects but has not named.
3. **Solve** — Present your product/service as the clear, logical resolution. The transition from agitation to solution should feel like relief.

**Schwartz warning:** Agitation without a credible solve creates distrust, not conversion. Match the intensity of the pain to the scale of the promised solution.

---

## 2. Schwartz's 5 Stages of Awareness

Eugene Schwartz (*Breakthrough Advertising*, 1966) argues that the right hook strategy is determined by the reader's current awareness stage, not by your product's features.

| Stage | Reader state | Lead with | Avoid |
|---|---|---|---|
| **1 — Unaware** | Does not know they have a problem | Story, provocative idea, striking contrast | Product name, category jargon |
| **2 — Problem-Aware** | Feels the pain but doesn't know solutions exist | Precise pain description (PAS) | Solution before agitation |
| **3 — Solution-Aware** | Knows solutions exist but not yours | Mechanism differentiation | Generic category claim |
| **4 — Product-Aware** | Knows your product but hasn't committed | Proof, guarantee, objection removal | Repeating features they already know |
| **5 — Most-Aware** | Ready; just needs a compelling offer | The offer itself | Lengthy explanation |

**Application rule:** Before writing the first word, identify the awareness stage of the majority of your traffic. Mismatching stage is the most common structural error — writing a Stage 5 headline to a Stage 1 audience destroys conversion.

**Hook examples by stage:**
- Stage 1: *"Most founders never find out why their second round fell apart."*
- Stage 2: Name the problem so precisely that the reader thinks you're describing them personally. This is where PAS is most powerful.
- Stage 3: Focus on mechanism — *how* your approach solves the problem differently.
- Stage 4: Lead with the strongest proof element: case study, guarantee, or removed objection.
- Stage 5: The headline can be as direct as *"$97. Full access. Cancel anytime."*

---

## 3. Landing Pages

### One Promise Per Page

Every landing page should have exactly one conversion goal. A page that asks for a demo, a newsletter signup, and a download forces readers to make three decisions. Each additional ask reduces the probability of any one of them happening. (Cialdini: reducing choice reduces decision fatigue and increases compliance.)

### Above-the-Fold Hierarchy

Everything visible before the scroll must resolve the reader's first question: *"Is this for me, and is it worth my time?"*

```
Headline          — The single biggest benefit or the sharpest problem
Subheadline       — Expand the headline; add specificity or credibility
Value proposition — 2–4 bullets: outcomes, not features
CTA (primary)     — One button, imperative verb + outcome
Social proof      — Logo row or a single killer testimonial with attribution
```

**Order is load-bearing.** A credible headline buys attention for the subhead. The subhead earns the scroll. Moving social proof above the value prop can work for high-awareness audiences (Stage 4–5) but confuses Stage 2–3 readers who don't yet trust the frame.

### Below-the-Fold Structure

For long-form pages (SaaS, high-ticket services):

1. Problem section (PAS agitation)
2. Mechanism — how it works, at a level that builds confidence without overwhelming
3. Proof — case studies with concrete numbers, not adjectives
4. Objection handling — address the two or three real objections; ignoring them does not make them disappear
5. Offer summary — restate what they get, at what price, by when
6. Risk reversal — guarantee, trial period, cancellation policy
7. Final CTA — repeat the primary action; do not introduce a new ask

---

## 4. Email Subject Lines

### Constraints

- Hard limit: 50 characters (preview pane + mobile). Aim for 35–45.
- The subject line's only job is to earn the open. It does not summarize the email.
- Preview text is the second subject line. Write it intentionally; do not let it default to "View in browser."

### Patterns That Drive Opens

| Pattern | Example | Works because |
|---|---|---|
| **Curiosity gap** | "The pricing mistake we almost shipped" | Creates an open loop the reader must close |
| **Specificity** | "3 changes that lifted our conversion 22%" | Numbers signal credible, scannable content |
| **FOMO / deadline** | "Closes Friday — last cohort seats" | Scarcity and urgency (Cialdini: scarcity principle) |
| **Social proof** | "How Figma cut their onboarding time in half" | Authority + proof before the click |
| **Direct question** | "Is your retention rate fixable?" | Forces a yes/no answer that self-selects engaged readers |
| **Contrarian** | "Why most A/B tests are a waste of time" | Interrupts the expected narrative; triggers curiosity |

**Anti-pattern:** Question marks and exclamation points in the same subject line signal desperation. Choose one device per subject line.

---

## 5. Email Skeletons

### Cold Email (3-sentence structure)

Effective cold email is not about your product. It is about the reader's problem.

```
Sentence 1 (Problem + context): One observation about a specific problem
                                 the prospect has, grounded in their world.
Sentence 2 (Proof + relevance): One concrete result you produced for
                                 a comparable company or situation.
Sentence 3 (Ask):               A single, low-friction next step —
                                 never a demo request in the first touch.
```

Example:
> Noticed your support team is using Zendesk for internal escalations — most teams outgrow that setup around 50 agents and end up with duplicated tickets. We helped [Competitor] cut internal escalation time by 40% in six weeks. Would a 10-minute breakdown of how they did it be useful?

### Nurture Email

Goal: advance the reader through the awareness funnel over multiple touches.

- Lead with a piece of value (insight, framework, short case study) before the soft ask
- Each email has one topic and one CTA
- The CTA should be a micro-commitment that naturally precedes the macro-commitment (read this → watch this → book this)

### Breakup Email

Sent after 3–4 unresponded touches. The breakup email often has the highest reply rate of the sequence.

```
Subject: Should I close your file?

Acknowledge the silence without blame. Offer an easy "no."
Optionally provide one last piece of value or alternative resource.
Make the opt-out frictionless.
```

Example:
> Haven't heard back, so I'll assume the timing isn't right — no problem at all. If that changes, you know where to find us. In the meantime, [resource link] might be useful.

---

## 6. Headline Formulas

| Formula | Template | Example |
|---|---|---|
| How-to | "How to [achieve outcome] [without objection]" | "How to close enterprise deals without a procurement battle" |
| Number list | "[N] [things] that [outcome]" | "7 SaaS onboarding flows that doubled activation" |
| Question | "[Question that names reader's pain]?" | "Why are 60% of your activated users gone in 30 days?" |
| Contrarian | "Why [common belief] is [wrong/costing you X]" | "Why longer free trials are killing your conversion" |
| Why X is... | "Why [product/approach] is [surprising claim]" | "Why boring copywriting outperforms clever in B2B" |
| Benefit-first | "[Big outcome] for [specific audience]" | "Enterprise-grade security for teams under 10" |

**Ogilvy rule:** "On the average, five times as many people read the headline as read the body copy." Test headlines first. The body copy is irrelevant if the headline loses the reader.

---

## 7. CTA Microcopy

### Imperative verb + outcome, not a noun

| Weak | Strong |
|---|---|
| Submit | Get my free report |
| Download | Download the 2026 Benchmark |
| Register | Save my seat |
| Learn more | See how it works |
| Subscribe | Send me the weekly digest |

**Rules:**
- First-person phrasing ("Get my..." vs. "Get your...") consistently outperforms second-person in A/B tests — Joanna Wiebe (Copyhackers) research finding.
- Every CTA should resolve the reader's micro-anxiety: *What happens when I click this?* The button label answers that question.
- One primary CTA per section. A secondary CTA (e.g., "or schedule a demo") can exist but must be visually subordinate.

---

## 8. Social Proof Typography

### The formula: Who + What they do + What they got (concrete number)

Bad:
> "Great product, highly recommend." — Sarah K.

Good:
> "We reduced customer churn by 18% in the first quarter after switching. The reporting alone saved my team four hours a week." — Sarah Kim, Head of Customer Success, Lattice (500-person SaaS)

**Elements:**
- Full name (not initials — initials signal invented testimonials)
- Title and company (establishes that this person is comparable to your target reader)
- Concrete result (number or time frame — adjectives are not proof)
- Photo where possible (Cialdini: liking and social proof compound when the face is visible)

**Placement:** Put the most powerful testimonial directly below the CTA on the above-the-fold section, or immediately after the price on a sales page (addresses the highest-anxiety moment).

---

## 9. Sales Psychology

### Anchoring on price

Present the highest tier or full price before the discounted or monthly price. The first number sets the psychological anchor. Reversing the order (low-to-high) makes the high price feel punitive.

### Urgency without manipulation

Legitimate urgency: a real deadline (cohort close date, founding pricing end, capacity limit). Manufactured urgency ("Offer expires in 24 hours!" on a page that resets every 24 hours) erodes trust on subsequent visits and signals that no urgency actually exists. Ogilvy: "The consumer is not a moron; she is your wife."

### The Fear / Hope / Dream Triangle (Schwartz)

Every product sells by resolving a tension between what readers fear, what they hope for, and what they dream about.

| Appeal | Copy frame | Best awareness stage |
|---|---|---|
| **Fear** | Loss aversion — "You're leaving $X on the table" | Stage 2 — problem must be felt before fear copy lands |
| **Hope** | Achievable near-term outcome — "Be ready in 30 days" | Stage 3–4 — reader is evaluating options |
| **Dream** | Identity transformation — "Become the kind of founder who..." | Stage 5 — reader is already sold; reinforce buyer identity |

---

## 10. Voice and Brand Alignment

### When to break grammar rules

Informal contractions, sentence fragments, and em-dashes are grammatically incorrect and often right for sales copy:

- "We don't do six-month contracts. Ever." (fragment — emphasis)
- "You're not behind. You're just early." (contraction — warmth)
- "The result: a 3x lift in qualified pipeline, in week two." (colon + comma instead of double em-dash — same pace, lower punctuation density)

Rules for breaking rules:
1. The break must serve the voice, not mask weak thinking.
2. Every fragment must be intentional, not an error of omission.
3. Fragments are stronger when followed by a complete sentence that pays them off.

### Jargon as an awareness signal

Domain jargon is appropriate for Stage 4–5 copy (product-aware audiences) and creates in-group recognition. It is alienating in Stage 1–2 copy where the goal is connection with a reader who may not yet know the category.

---

## 11. Anti-Patterns

| Anti-pattern | Why it fails | Fix |
|---|---|---|
| Corporate jargon ("synergy", "leverage", "ecosystem") | Signals generic positioning; no reader feels seen | Replace with the specific outcome in plain language |
| Undifferentiated value props ("world-class", "innovative", "best-in-class") | Every competitor says the same thing; no decision signal | Name the mechanism that delivers the result |
| Hedging language ("may help", "can potentially") | Reduces perceived confidence; signals legal timidity | Commit to the specific claim with proof behind it |
| Multiple CTAs per page section | Forces a decision about which decision to make | One CTA per section; visual hierarchy enforces priority |
| Features without outcomes | Readers do not buy features; they buy the life the feature enables | "256-bit encryption" → "Your data stays yours, always" |
| Passive voice in headlines | Reduces energy and ownership | Active construction; name who does what to whom |
| Long-form without pattern interrupts | Reader fatigue; attention drops after 300 words of prose | Break with subheads, bullets, pull quotes, and white space every 150–200 words |

---

## 12. Sentence fragments (deliberate use)

In sales copy, sentence fragments are not errors. They are rhythm tools. A
fragment slows the eye, punches a beat, and forces the next sentence to land
harder than a syntactically complete paragraph would.

### The rule

A fragment earns its place in sales copy when it:

1. Compresses an idea to its essential punch ("Not gonna happen.")
2. Echoes and lands the close of a longer preceding sentence ("We rebuilt the
   pipeline. From scratch.")
3. Marks a beat the eye can rest on between long persuasion paragraphs.

A fragment is wrong when:
1. It substitutes for a claim the writer didn't make ("Important. Big. Now.")
2. It clusters too tightly — three fragments in a row become noise.
3. It hides ambiguity. If the reader can't reconstruct the subject/verb from
   context, the fragment fails.

### Worked example

Flat (all complete sentences):
> "Our pricing is straightforward. There are no hidden fees, no surprise
> charges, and no escalators. You pay one price for as long as you stay."

With fragments (rhythm + emphasis):
> "Our pricing is straightforward. No hidden fees. No surprise charges. No
> escalators. One price — for as long as you stay."

The fragment cluster lands harder because each "No" forces a beat. Reading
aloud, you hear the three-beat punch before the resolution sentence.

### Headline fragments

Headlines are almost always fragments by structure ("3 changes that lifted
conversion 22%"). That is correct. Headlines compete for attention against
other headlines, not against grammar textbooks.

### Subject-line fragments

Email subject lines work the same way — most high-performing subject lines
are fragments. "The pricing mistake we almost shipped" (no main verb) outperforms
the syntactically complete "We almost shipped a pricing mistake." The
fragment is more compressed and reads less corporate.

### When to break the rule

- Regulated copy (financial disclosures, healthcare claims) often requires
  complete sentences for compliance review. Fragments may trigger flags.
- B2B enterprise pitches to procurement audiences benefit from formal,
  complete-sentence prose — fragments can read as too marketing-flavored.
- Long-form sales pages can use fragments freely but should keep at least
  60–70% complete sentences for narrative momentum.

### References

- Strunk & White, *Elements of Style*, §16.
- Sugarman, J. *The Adweek Copywriting Handbook* (2007), ch. 4 (slippery slide
  technique — short sentences and fragments pull the reader forward).

---

## 13. The 60-second elevator pitch

The elevator pitch is the shortest unit of sales narrative: 60 seconds, ~150
spoken words, designed to deliver a complete value claim in the time it
takes an elevator to reach the executive floor. Originally an in-person
sales artifact, it now appears in pitch decks (slide 1), cold-email opens,
website hero sections, and founder introductions.

### The 60-second structure

| Beat | Time | What to deliver |
|------|------|-----------------|
| Hook | 5–10 s | A specific, named problem the listener already feels |
| Audience | 5 s | Who this is for (your ICP, in one phrase) |
| Solution | 15 s | What you do, in one sentence — mechanism-level, not "platform" |
| Differentiation | 10 s | Why your approach beats the obvious alternative |
| Proof | 10 s | One number, one customer, one outcome |
| Ask | 5–10 s | The next step you want — never "let's talk soon" |

Total: ~60 seconds at conversational pace (~150 words).

### Worked example (B2B SaaS founder)

> "Enterprise sales teams lose 30% of qualified deals because nobody at the
> deal post-mortem can reconstruct what the prospect actually said three
> months ago. (Hook) We work with $50M–$500M ARR sales orgs (Audience) and
> build a deal-memory layer that turns every call, email, and ticket into
> a searchable narrative tied to the deal record. (Solution) Unlike CRM
> notes, which depend on reps writing them, our system captures the source
> material automatically — so the memory exists whether the rep updates
> the CRM or not. (Differentiation) Lattice ran us against their last 12
> months of closed-lost deals; we surfaced a recurring procurement
> objection their team had never aggregated, and it changed how they
> qualify deals. (Proof) Could I send you a 90-second video that shows
> exactly what we'd surface for one of your deals? (Ask)"

That is 150 words and lands every beat.

### When to break the rule

- The pitch is for an executive who has already named the problem; skip the
  hook and start at Audience.
- The pitch is internal (project pitch to your boss); proof and ask should
  swap weight — your boss cares more about the ask (resource, timeline) than
  the proof.
- High-context audiences (investors familiar with your space) can absorb a
  shorter pitch — 30 seconds, no hook, just solution + traction + ask.

### Common failure modes

- **No specific problem.** Hook becomes "Companies struggle with X" — generic
  and skippable.
- **Buzzword solution.** "AI-powered platform" tells the listener nothing.
  Name the mechanism.
- **Soft ask.** "Would love to grab coffee sometime" is not an ask. "Could
  I send you a 90-second video that shows X for your specific situation?" is.

### References

- Bartlett, M. *The Diary of a CEO* — founder pitch interviews (~150-word
  rule across multiple guests).
- Pink, D. *To Sell Is Human* (2012) — the modern elevator pitch as
  "Twitter pitch" and "one-question pitch" reformulations.

---

*Sources: Schwartz, Eugene. Breakthrough Advertising (1966). Ogilvy, David. Confessions of an Advertising Man (1963). Sugarman, Joseph. The Adweek Copywriting Handbook (2007). Cialdini, Robert. Influence: The Psychology of Persuasion (1984). Halbert, Gary. The Boron Letters (2013). Carlton, John. Kick-Ass Copywriting Secrets of a Marketing Rebel. Pink, Daniel. To Sell Is Human (2012).*

---

## The headline test — "would you click your own headline?"

**Rule.** David Ogilvy: "On the average, five times as many people read the headline as read the body copy. When you have written your headline, you have spent eighty cents out of your dollar" (*Ogilvy on Advertising*, 1983). Headlines do most of the work; the test for whether yours pays off is brutally simple: *if you saw this headline in someone else's feed, would you click it?*

Joanna Wiebe (Copyhackers) operationalized it: write the headline last, after the rest of the copy is drafted. Then run the click test against three rivals: your previous headline, your competitor's headline, and a deliberately weak placeholder. If yours does not beat all three on a 5-second click-decision, rewrite.

**The four-question click test.**

1. *Specificity.* Does the headline name a specific outcome, audience, number, or mechanism, or is it abstract?
2. *Curiosity gap.* Does the headline open a loop the reader needs the body copy to close?
3. *Self-relevance.* Can the target reader say "that is me" or "that is my problem" within 2 seconds?
4. *Promise.* Is there a concrete payoff for clicking, or is the headline a topic label?

A headline that scores 4/4 is publishable. A 3/4 is workshoppable. A 2/4 or below should be rewritten.

**Worked example — B2B SaaS landing page.**

- Weak: "Modern data infrastructure for the cloud era." (Score: 1/4 — vague, no curiosity, no specific audience, no payoff.)
- Strong: "Migrate from Postgres to MongoDB in 48 hours without dropping a write." (Score: 4/4 — specific source/target, time bound, names the fear, implicit mechanism.)

**Worked example — blog post.**

- Weak: "Our thoughts on the future of databases."
- Strong: "Why we replaced our analytics warehouse after 18 months — and what we'd do differently."

**Diagnostic.** Open the page in an incognito window. Look at only the headline for 2 seconds. Ask yourself out loud: "What would I get if I clicked this?" If your answer is vague, hedged, or "I'm not sure" — the headline is failing the test.

**When to break it.** Pure brand campaigns (Coca-Cola, Apple's "Think Different") legitimately use abstract headlines because the brand already carries the meaning. You almost certainly are not Coca-Cola or Apple. For direct-response B2B copy, specificity beats art every time.

**References.**

- Ogilvy, D. *Ogilvy on Advertising*. Crown, 1983, chapter on headlines.
- Wiebe, J. *Copyhackers* — "How to write headlines that get clicks." https://copyhackers.com/2016/05/headline-formulas/
- Sullivan, L. *Hey Whipple, Squeeze This: A Guide to Creating Great Ads*. Wiley, current edition.

---

## Subhead craft — earning the scroll after the headline

The *subhead* (also dek, subheadline, supporting headline) is the line directly under the headline. It does the work the headline couldn't fit. In conversion copy, the subhead converts a stop into a scroll — and the difference between a 0.2% and a 2% conversion rate often lives in that one line.

### The subhead's three jobs in conversion copy

1. **Expand the promise.** The headline made a claim. The subhead specifies it. ("Faster CI" → "Builds in 90 seconds, not 18 minutes.")
2. **Add the proof axis.** The strongest subheads carry a number, a name, or a comparison the headline couldn't fit. ("From 18 min to 90s — without changing a line of YAML.")
3. **Pre-answer the first objection.** The most common reaction to a strong headline is "yeah, but…" The subhead earns the scroll by anticipating and neutralizing it. ("No new tooling, no team retraining, no migration script.")

### Sizing for conversion copy

- **Headline:** 5–10 words. Above-the-fold heroes lean shorter.
- **Subhead:** 10–25 words. One line on desktop; up to two on mobile.
- **The combined hed+subhead read time:** under 4 seconds. If a visitor can't extract the proposition in 4 seconds, the hero is failing.

### The headline-subhead handoff (the failure point)

Most weak conversion copy doesn't fail at the headline or the subhead — it fails at the *handoff*. Two reliable patterns:

1. **Claim → Mechanism.** The headline makes the claim; the subhead names the mechanism. ("Ship database changes without dropping a write. Online schema migrations run as background ops, with zero application code changes.")
2. **Problem → Outcome.** The headline names the problem; the subhead names the better future. ("Your CI is the slowest part of your stack. Cut build time 12x with a one-line config change.")

### Worked example — B2B SaaS landing page hero

**Weak (headline + subhead repeat each other):**

- Hed: "Modern data infrastructure for the cloud era."
- Sub: "Build, scale, and operate your data platform on a modern foundation."

Both lines say nothing twice. Score: 0/3 jobs.

**Strong (handoff: claim → mechanism, with proof embedded):**

- Hed: "Migrate from Postgres to MongoDB in 48 hours."
- Sub: "Without dropping a write — and without retraining your engineers on a new query language."

The headline makes the claim. The subhead specifies (no dropped writes, no retraining), pre-answers two objections (data loss, learning curve), and earns the scroll. Score: 3/3.

### Worked example — email subject line + preheader

The subject line and preheader are a hed+dek pair for email. Same rules apply:

- Subject: "Your CI is the slowest part of your stack."
- Preheader: "Three customers cut build time from 18 min to 90s last quarter. Here's the config change."

The subject names the problem; the preheader adds the proof axis (three customers, specific numbers) and previews the payoff (the config change).

### Anti-patterns specific to conversion subheads

1. **Restating the headline.** "Faster builds — Build faster with our platform." Adds nothing. Cut.
2. **Listing features.** "Faster builds — Includes parallel runners, persistent caches, and auto-scaling." The reader hasn't agreed to care about features yet. Lead with outcomes.
3. **Hedging.** "Faster builds — Most teams see significant improvements." Vague qualifiers ("most," "significant") collapse credibility. Use specific numbers or none.
4. **The brand-voice subhead with no information.** "Faster builds — Welcome to the future of CI." Brand-only campaigns can do this. Direct-response copy cannot.
5. **The CTA-disguised subhead.** "Faster builds — Click here to learn more." The subhead is not the CTA. The CTA is the CTA.

### When to break it

Pure brand campaigns (no conversion goal) can use a single line with no subhead — Apple's "Think Different" had no dek, and didn't need one. You almost certainly need one. If you have a conversion goal, the subhead is doing work; don't skip it.

### Diagnostic — the four-second test

Cover everything below the hero. Read the headline and subhead in four seconds. Three questions:

1. Do I know what is being offered?
2. Do I know who it's for?
3. Do I know what would change for me if I scrolled?

3/3 ships. 2/3 needs a subhead rewrite. 0–1 needs both lines rewritten.

### References

- MasterClass — [How to Write a Subheading: 4 Tips for Writing a Dek](https://www.masterclass.com/articles/how-to-write-a-dek).
- River — [Hed and Dek in Journalism: Write Headlines That Get Clicks](https://rivereditor.com/blogs/generate-hed-dek-headline-subhead).
- See also: `writing-expert` skill — *Subhead / deck — the line under the headline* for the general-purpose treatment with editorial (non-conversion) examples.
