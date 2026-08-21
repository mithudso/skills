<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `rhetorical-frameworks-deep` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: rhetorical-frameworks-deep
description: Deep reference for rhetorical and argument frameworks in technical/business writing — Aristotle (ethos/pathos/logos), Toulmin (claim/data/warrant), Pyramid Principle, MECE, Cialdini's 6 principles, narrative transportation, stasis theory. TRIGGER: "structure this argument", "make this more persuasive", "what framework should I use", "frame this for an exec", "argue for/against X in writing", "organize my argument", "build the case for", "how do I convince". SKIP: prose-only style edits (use writing-expert); doc review/critique (use document-critique); verbal negotiation tactics; "what is X" definitional questions with no writing task attached.
version: 1.2.0
updated: "2026-05-29"
related_skills:
  - writing-expert
  - document-critique
  - executive-comms
whenToUse:
  - User asks to structure an argument or proposal
  - User asks which framework to use for a writing task
  - User is preparing a persuasive document (proposal, EBR, escalation, design doc)
  - User wants to audit their draft for logical flaws before sending
whenNotToUse:
  - Pure stylistic editing with no argument-structure question
  - Doc review / fact-check (use document-critique)
  - Verbal negotiation tactics (no registered skill — advise user)
---

# Rhetorical Frameworks (Deep Reference)

Good technical and business writing persuades through structure, not just accuracy. This skill provides working-reference definitions and concrete examples for eight foundational frameworks. Each framework solves a different problem: some diagnose where a disagreement lives, some organize the argument's skeleton, some explain why readers accept or resist claims.

## How to use this skill

**If the user has not provided a document or draft:** ask one targeted question before proceeding — e.g., "What document or argument are you working on, and who is your audience?" Do not explain frameworks speculatively until you have the user's actual content or goal.

When invoked with content, follow this sequence:

1. **Diagnose the situation first.** Use stasis theory (§7) to identify which question is actually in dispute. If the audience hasn't accepted the facts or quality of the problem, framework selection for the solution argument is premature.
2. **Select one primary framework** that fits the document type and audience. A one-page exec note calls for Pyramid Principle; a detailed proposal benefits from Toulmin; a categorization exercise needs MECE.
3. **Apply the framework to the user's actual content** — produce a concrete restructured outline, rewritten opening, or annotated draft section, not just an explanation of the framework. Work only with content the user has provided; do not invent data, claims, or examples they did not supply.
4. **Check for distortions** (§8) in the user's existing draft before or after restructuring.
5. **Self-check before delivering:** confirm (a) the framework is named, (b) the rationale is 2-4 sentences, and (c) every restructured claim traces back to something the user provided.
6. **Hand off** to `writing-expert` for prose drafting or `document-critique` for a post-draft review pass.

Output format: deliver a restructured outline or rewritten opening (scoped to one section or the document's top-level structure — not a full rewrite unless the user requests it), followed by a brief rationale (2-4 sentences) naming which framework was applied and why.

## When to use this skill

- You have a position to defend and need to choose an argument architecture before writing.
- You are asked to "make this more persuasive" or "reframe this for leadership."
- You are producing a proposal, escalation, EBR, design doc, or change-request that must move a skeptical audience.
- You want to audit your own draft for logical distortions before sending.

## When NOT to use this skill

- You only need style or grammar help — use `writing-expert`.
- You have a finished draft and want a structured critique — use `document-critique`.
- The persuasion is verbal: a meeting, negotiation, or conversation (no writing-specific skill registered; advise directly).
- The goal is story structure for a narrative piece (no writing-specific skill registered; use `writing-expert` for narrative drafting).

---

## Framework reference

### 1. Aristotle's three appeals

*Rhetoric* (c. 350 BCE) identifies three channels through which speakers gain assent.

| Appeal | Core question | How it works in technical writing |
|--------|---------------|-----------------------------------|
| **Ethos** (credibility) | "Why should I trust you?" | Cite your track record, acknowledge the scope of your data, disclose limitations honestly. A TAM who leads with three months of account telemetry before recommending a change signals credibility before the ask. |
| **Pathos** (emotion / values) | "Why should I care?" | Connect the recommendation to what the reader already values: uptime, revenue, team safety, career outcome. Technical writers suppress pathos instinctively; that is often a mistake with executive audiences. |
| **Logos** (logic) | "Does this follow?" | Data, causal chains, analogies, and formal argument structure. The trap: most technical writers over-rely on logos and present data without first establishing ethos or a clear stake. |

**Worked example.** A TAM drafting a "migration urgency" note might open with:
- *Ethos*: "I have reviewed your last 90 days of cluster metrics and the Atlas release notes for the version you are running."
- *Pathos*: "A version-level CVE that matches your topology was disclosed last week; I want to make sure your team hears about it before it surfaces as an incident."
- *Logos*: "The EOL date is 2026-07-01. Patch testing in staging typically takes 3-4 weeks. That leaves one sprint of margin."

All three in the first three sentences — before a single recommendation.

---

### 2. Toulmin model

Toulmin, *The Uses of Argument* (1958). Six components; the first three are load-bearing in almost every argument.

| Component | Definition |
|-----------|-----------|
| **Claim** | The conclusion you want the reader to accept. |
| **Data** | The evidence you are standing on. |
| **Warrant** | The principle that connects data to claim. Often implicit; making it explicit is the single most clarifying move. |
| **Backing** | Evidence that the warrant itself is reliable. |
| **Qualifier** | The conditions under which the claim holds ("usually," "in cases where"). |
| **Rebuttal** | Conditions under which the claim would not hold. Acknowledging them preempts reader objections. |

**Worked example (TAM escalation).**
- *Claim*: This account needs a dedicated support path for the next 30 days.
- *Data*: Four P1 tickets in six weeks; two required Director-level escalation; average time-to-resolution was 3x the SLA.
- *Warrant*: When resolution patterns deviate this far from SLA norms, standard routing is insufficient and a dedicated path reduces risk of churn.
- *Backing*: Internal study shows accounts with persistent SLA breaches churn at 2.4x the baseline rate within 6 months (cite internal data).
- *Qualifier*: This applies to accounts in enterprise tiers, not developer-tier relationships.
- *Rebuttal*: If the ticket pattern is attributable to a single one-time migration event, standard routing may resume after that event closes.

The rebuttal is not weakness; it is the evidence that you have thought clearly.

---

### 3. Pyramid Principle (Minto)

Minto, *The Pyramid Principle* (1987). Argument travels top-down: answer first, support follows.

**SCQA scaffold:**

1. **Situation** — shared, uncontested context. What everyone already knows.
2. **Complication** — the tension or change that makes the situation unstable.
3. **Question** — the question the complication forces the reader to ask.
4. **Answer** — your recommendation, stated immediately.

Everything below the answer is evidence arranged in a pyramid: main supporting points, then sub-points beneath each.

**When to use it.** Any document where the audience has limited time and high trust in your judgment: executive summaries, status updates, escalation memos, design proposals. Do not bury your recommendation in paragraph six.

**TAM example sketch.**
- *Situation*: Atlas M30 cluster, 80% CPU headroom consumed for the past 14 days.
- *Complication*: Month-end batch jobs run next Friday; past patterns suggest a 2x spike.
- *Question*: How should we respond before the batch window opens?
- *Answer*: Recommend a one-tier vertical scale-up this week, reverting after batch completes.

---

### 4. MECE

*Mutually Exclusive, Collectively Exhaustive* — a structuring principle from consulting practice (McKinsey, c. 1970s; popularized by Minto and Rasiel).

- **Mutually exclusive**: categories do not overlap. If two boxes on your slide could both contain the same item, the structure is ambiguous.
- **Collectively exhaustive**: together, the categories account for all relevant cases. If an important case has no home, there is a gap.

**Common violations in technical writing:**

| Violation | Example |
|-----------|---------|
| Overlap | "Performance issues" and "slow queries" as separate root-cause categories — slow queries are a subset of performance issues. |
| Gap | A runbook lists three failure modes but omits network partitions in a distributed system. |
| False balance | Three tiny sub-categories of one problem paired against one large category of another, implying equal weight. |

**When MECE is most important.** Root-cause taxonomies, options analyses, decision trees, and any slide where the structure is meant to be exhaustive. MECE is not required for narrative text — forcing it there produces bullet-point prose.

---

### 5. Cialdini's six principles

Cialdini, *Influence: The Psychology of Persuasion* (1984). Originally studied in sales; each principle maps cleanly to internal technical writing.

| Principle | Core mechanism | Internal writing application |
|-----------|----------------|------------------------------|
| **Reciprocity** | People return favors. | Share genuinely useful data, analysis, or a heads-up before making your ask. A TAM who sends a monthly cluster-health summary builds reciprocity credit. |
| **Commitment / consistency** | People act in line with prior commitments. | Cite the reader's own stated goals. "In Q3 you prioritized 99.99% uptime — this proposal supports that commitment." |
| **Social proof** | People follow what similar others do. | "Four other accounts in your industry made this change in the past six months." Customer references and peer benchmarks operate here. |
| **Authority** | People defer to credible experts. | Credentials, tenure, institutional backing. Citing Atlas documentation, a MongoDB engineer, or a known peer in the reader's domain. |
| **Liking** | People agree with people they like. | Genuine shared context, acknowledgment of the reader's constraints, appropriate warmth in framing. Not flattery — alignment. |
| **Scarcity** | People value what is rare or time-limited. | "EOL date is July 1. After that, CVE patches are not backported." Scarcity should be real, not manufactured — readers notice the difference. |

---

### 6. Narrative transportation

Green and Brock, "The Role of Transportation in the Persuasiveness of Public Narratives," *Journal of Personality and Social Psychology* (2000). Transportation theory holds that when readers become absorbed in a story, counter-arguing drops off. The narrative world temporarily displaces the reader's critical frame.

**Why this matters for technical writing.** Data can be disputed; a story about a specific, named incident is harder to argue with because the reader processes it as witness testimony rather than abstract claim. A TAM writing about downtime risk who opens with "Last quarter, a similar cluster at a comparable account experienced a 4-hour outage during month-end close — here is what the ticket trail looked like" is using transportation, not just statistics.

**When stories work better than data.** Use narrative transportation when the audience is skeptical of aggregate statistics, when the risk is hard to visualize from numbers alone, or when you need to create emotional weight behind a technically correct but dry recommendation. The story does not replace the data; it opens the reader to receiving the data. Pair them: story first, statistics second.

---

### 7. Stasis theory

Classical origin (Hermagoras, c. 150 BCE); developed in *Rhetorica ad Herennium* and Cicero's *De Inventione*. Stasis theory diagnoses where a disagreement actually lives. Before you can persuade, you must know which question is actually in dispute.

| Stasis | The question at issue | Typical signal |
|--------|-----------------------|----------------|
| **Fact (conjectural)** | Did this happen / is this true? | "I don't think we actually have a performance problem." |
| **Definition** | What is this? What category does it belong to? | "This is a configuration issue, not a bug." |
| **Quality** | How serious is it? Was it right or wrong? | "Yes, it is slow — but it is within acceptable range." |
| **Policy** | What should we do about it? | "Granted it is serious — but should we scale up or re-index?" |

**How to use it.** Read your reader's objections and identify which stasis they occupy. If they dispute *fact*, no amount of policy argument will land — you must resolve the fact question first. If they dispute *quality*, cite impact data and benchmarks before proposing action. Arguing at the wrong stasis is the most common reason a technically correct proposal fails to persuade.

**TAM example.** A customer who says "I don't see why this is urgent" is at the *quality* stasis, not the *policy* stasis. Presenting migration steps (policy) will not move them. Present blast-radius data and SLA implications (quality) first.

---

### 8. Common distortions to detect in your own writing

These are logos-emphasis errors — the argument looks rational but contains a structural flaw.

- **Appeal to authority without evidence**: "MongoDB recommends X" — does the official documentation actually say that, or are you generalizing?
- **Sunk-cost framing**: "We have already invested six months in this approach" — past investment is not a reason to continue; expected future value is.
- **False dichotomy**: "We either scale up or we accept downtime" — are those genuinely the only two options?
- **Availability bias in risk framing**: Citing the most recent incident as representative of base-rate risk when it may be an outlier.
- **Overprecision**: Stating projections to two decimal places when the underlying data has high variance — false precision erodes trust when readers notice it.
- **Hedging that undermines the claim**: "This might possibly be worth considering" — qualifiers are legitimate (Toulmin), but stacking them signals you do not believe your own argument.

Before sending a high-stakes document, scan it explicitly against this list.

---

### 9. Anaphora and epistrophe (deliberate repetition)

Two of the oldest figures of speech in the classical rhetorical canon. Both
trade on the same insight: when an audience hears the same word or phrase
land in the same position across consecutive clauses, the repetition signals
that the speaker chose every word — and the repeated word becomes the
emotional spine of the passage.

| Figure | Position of the repetition | Example |
|--------|---------------------------|---------|
| **Anaphora** | Same word/phrase at the **start** of consecutive clauses | "We shall fight on the beaches, we shall fight on the landing grounds, we shall fight in the fields and in the streets..." (Churchill, 1940) |
| **Epistrophe** | Same word/phrase at the **end** of consecutive clauses | "...of the people, by the people, for the people." (Lincoln, 1863) |

### Why they work

The repetition does three things at once:
1. **Cognitive marking** — the repeated phrase becomes a memory anchor;
   listeners and readers recall the repeated phrase long after they forget
   the surrounding prose.
2. **Implied parallelism** — repeating the lead-in (or tail) signals that
   the items in the series are coordinate. The audience does not have to
   work out the structure; it is given.
3. **Build of intensity** — three or more repetitions accumulate weight.
   The third instance lands harder than the first because the pattern is now
   visible and the audience anticipates it.

### Worked example (TAM escalation memo)

Flat:
> "The customer's renewal is at risk. Their service tickets are aging.
> Their executive sponsor stopped responding to our outreach last week,
> and their usage is down 30%."

Anaphora ("We watched..."):
> "We watched their service tickets age past SLA. We watched their executive
> sponsor stop responding to outreach. We watched their usage fall 30% in
> sixty days. The renewal is at risk because we watched — and did not act."

The anaphora compresses four distinct facts into one rhythmic build, and the
final sentence weaponizes the figure to indict the team's response. The same
four facts, listed without anaphora, read as a status update; with anaphora,
they read as an argument for urgent action.

### Rules for use in business and technical writing

- **Three is the minimum.** Two repetitions are pattern-noise; three creates
  the figure. Four is the standard ceiling — five starts to sound like
  speechwriting (which may or may not be what you want).
- **Use sparingly.** One anaphora or epistrophe per document is plenty. More
  than one and the writing tips from "rhetorically deliberate" into
  "performative."
- **Reserve for the stakes you're actually claiming.** Anaphora in a status
  update about cluster CPU utilization is theater. Anaphora in a memo
  arguing that an account is about to churn is appropriate emphasis.
- **Avoid in apology and bad-news letters.** Repetition reads as rhetorical
  flourish — wrong register when the audience expects directness.

### When to break the rule

- Marketing copy, opinion pieces, and keynote speeches use anaphora and
  epistrophe far more freely; the audience expects voice-forward prose.
- Legal writing has its own repetition rules (defined-term consistency); do
  not import anaphora into a contract.

### References

- Aristotle. *Rhetoric*, Book III on figures (c. 350 BCE).
- Lanham, R. *Analyzing Prose* (2003), ch. 7 — figures of repetition.
- Churchill, W. "We Shall Fight on the Beaches." House of Commons, 4 June 1940.
- Lincoln, A. Gettysburg Address, 19 November 1863.

---

### 10. Tricolon and isocolon (sentence rhythm)

Two related figures from the same classical inventory. Both exploit the
three-part rhythm that human cognition treats as complete.

**Tricolon** — a series of three parallel words, phrases, or clauses.
"Veni, vidi, vici." "Of the people, by the people, for the people."
"Friends, Romans, countrymen."

**Isocolon** — a tricolon (or other multi-part series) in which the parts
match in length and syntactic structure. "I came, I saw, I conquered" is both
tricolon (three parts) and isocolon (each part: pronoun + past-tense verb,
exactly two words).

### Why three works

Cognitive research (Miller, "The Magical Number Seven, Plus or Minus Two,"
1956; Cowan 2001 update to four) shows that small sets of three items feel
complete in a way that two-item lists feel incomplete and four-item lists
feel like dumps. Speechwriters, ad copywriters, and political messaging
teams have exploited this for two thousand years.

### Worked example

Flat (four-beat, uneven):
> "The migration will be faster, cheaper, safer, and it also has fewer
> dependencies and a simpler rollback story."

Tricolon (three beats, clean):
> "The migration is faster, cheaper, and safer."

Isocolon (three parallel two-word phrases):
> "Faster to ship. Cheaper to run. Safer to operate."

The isocolon version costs four extra words but lands harder than either
alternative because each beat is exactly the same length — the parallelism
becomes audible.

### Use in technical and business writing

- **Recommendations and benefit lists** — when you have three benefits,
  resist the urge to pad to four. When you have four, ask which two combine.
- **Closing sentences** — a tricolon at the end of an argument anchors the
  reader's memory of the argument. Use it sparingly; once per document.
- **Status report key findings** — three findings feel actionable; five feel
  like a dump.

### When to break the rule

- The natural list is genuinely two items or genuinely four. Padding three
  into four or trimming four into three to chase rhythm produces hollow
  prose.
- Reference and catalog content (API param lists, error tables) should be
  exhaustive, not rhythmic — the audience is scanning, not absorbing.

### Anaphora + tricolon combined

When anaphora wraps a tricolon, the result is the highest-density rhetorical
move available in English. Churchill's "We shall fight on the beaches, we
shall fight on the landing grounds, we shall fight in the fields..." is
anaphora over a tricolon. Reserve this combination for moments where the
stakes are genuinely high.

### References

- Aristotle. *Rhetoric*, Book III.
- Lanham, R. *Analyzing Prose* (2003), figures of repetition.
- Miller, G. "The Magical Number Seven, Plus or Minus Two." *Psychological
  Review* 63(2), 1956.
- Cowan, N. "The Magical Number 4 in Short-Term Memory." *Behavioral and
  Brain Sciences* 24(1), 2001.

---

## Composition with sibling skills

- After choosing a framework here, hand off to `writing-expert` for drafting.
- After drafting, hand off to `document-critique` for the review pass.
- For face-to-face persuasion (not writing), no dedicated skill is registered; advise the user directly on verbal tactics.
- For story-driven narrative framing, use `writing-expert` and ask it to apply narrative transportation principles.

---

## References

- Aristotle. *Rhetoric* (c. 350 BCE). Trans. W. Rhys Roberts. Penguin Classics.
- Toulmin, S. E. *The Uses of Argument*. Cambridge University Press, 1958.
- Minto, B. *The Pyramid Principle: Logic in Writing and Thinking*. Minto International, 1987.
- Cialdini, R. B. *Influence: The Psychology of Persuasion*. Harper Business, 1984 (rev. 2007).
- Green, M. C., & Brock, T. C. "The role of transportation in the persuasiveness of public narratives." *Journal of Personality and Social Psychology* 79(5), 701-721, 2000.

---

## ABT (And, But, Therefore) — the minimum-viable argument spine

**Rule.** Randy Olson's ABT compresses any persuasive structure into three connectives: *and* (shared premises the audience already accepts), *but* (the contradiction, friction, or gap the argument exists to resolve), *therefore* (the conclusion or call to action). ABT is rhetorically adjacent to Toulmin (data → warrant → claim) and Minto (situation → complication → resolution), but it sits at a lower level of granularity: it is the *one-sentence argument* the longer structure expands.

**Mapping to other frameworks already in this skill.**

| Layer | ABT slot | Minto | Toulmin | SCQA |
|---|---|---|---|---|
| Setup | And | Situation | Data | Situation |
| Tension | But | Complication | Rebuttal/Qualifier | Complication, Question |
| Resolution | Therefore | Resolution / Key line | Claim | Answer |

ABT is the *spine*; Minto and Toulmin are the *skeleton built around it*. If your one-sentence ABT does not hold up, the long-form argument cannot either.

**Worked example — RFC for a sharded migration.**

- And: "We are at 85% capacity on the single replica set and read latency is within SLA today."
- But: "Projected growth at current ingestion rate breaches the disk ceiling in 14 months and the read-after-write QPS curve crosses the documented single-RS ceiling at month 11."
- Therefore: "We must commit to a sharded topology in this quarter to leave a safety margin for index rebuilds and chunk migrations before either ceiling hits."

A reader who agrees with the ABT will read the full RFC for *how*, not *whether*. A reader who disagrees will point at the *but* — which is exactly where you want disagreement to land.

**When to break it.**

- *Pure exposition* (reference documentation, glossaries): no argument, no ABT.
- *Multi-claim documents* (long technical proposals): use a top-level ABT plus subordinate ABTs per section. Trying to compress a 30-page proposal into one ABT will force false precision.
- *Inductive vs deductive register*: ABT is naturally deductive (conclusion follows logically). For inductive register (gather evidence, then conclude), invert to "We see X, and we see Y, *but* taken together they imply Z, *therefore*..." — the *but* still has to do real contradictory work, otherwise it is filler.

**Diagnostic.** Strip the document to its ABT. Read the ABT aloud to a colleague who has not seen the document. Ask: "what would change your mind?" If they cannot name a counter-claim that would falsify the *therefore*, the argument is not yet load-bearing — it is an assertion dressed as a deduction.

**References.**

- Olson, R. *Houston, We Have a Narrative: Why Science Needs Story*. University of Chicago Press, 2015.
- Olson, R., Barton, D., & Palermo, B. *Connection: Hollywood Storytelling Meets Critical Thinking*. Prairie Starfish, 2013 (introduces ABT in workshop form).
- Story Circles Narrative Training materials, science-needs-story.com.
