<!-- hub-reference-banner -->
> **Reference file — part of the `executive-comms` hub.** Formerly the standalone `negotiation-and-persuasion` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: negotiation-and-persuasion
description: Tactical layer for persuasion-shaped writing and verbal-negotiation prep. Covers Cialdini's 6 principles applied to written requests, Chris Voss's "Never Split the Difference" tactics (labeling, mirroring, calibrated questions), BATNA/ZOPA construction, Diamond's small-step concessions, soft-no language, anchoring, loss-vs-gain framing, and reciprocal concessions. TRIGGER: "draft a negotiation email", "how do I ask for X", "convince customer to Y", "frame this proposal", "Voss tactics", "BATNA", "soft no", "anchor the ask", "ask for a raise", "salary negotiation", "counter-offer email", "reciprocal concession". SKIP: argument structure (use rhetorical-frameworks-deep); exec memos (use executive-comms); marketing copy.
version: 1.1.0
updated: "2026-05-29"
related_skills:
  - rhetorical-frameworks-deep
  - writing-expert
  - executive-comms
sources:
  - "Cialdini, Robert. *Influence: The Psychology of Persuasion* (1984)"
  - "Voss, Chris. *Never Split the Difference* (2016)"
  - "Fisher, Roger & Ury, William. *Getting to Yes* (1981)"
  - "Diamond, Stuart. *Getting More* (2010)"
  - "Kahneman, Daniel. *Thinking, Fast and Slow* (2011)"
whenToUse:
  - User is drafting a request, proposal, or counter-offer email
  - User asks how to phrase a difficult ask — including salary, raise, or comp requests
  - User is preparing for a negotiation conversation
  - User wants to frame a proposal to move a decision
  - User is preparing for a salary or compensation negotiation
  - User asks about BATNA, ZOPA, anchoring, or Voss/Cialdini tactics while preparing for a negotiation
whenNotToUse:
  - Structuring an argument (use rhetorical-frameworks-deep)
  - Executive memo or briefing doc (use executive-comms)
  - Marketing or advertising copy
  - Definitional questions with no active negotiation or writing task
---

# Negotiation and Persuasion (Tactical Reference)

This skill operates at the **tactical layer** — the specific words, sequencing, and framing moves that get someone to say yes. It complements rhetorical-frameworks-deep (which handles argument structure) and executive-comms (which handles memo format). Use both when needed.

## How to use this skill

**If the user has not provided a draft or described their situation:** ask one targeted question before proceeding — e.g., "What are you trying to get, who is the other party, and what's their most likely objection?" Do not offer generic frameworks speculatively.

When invoked with a situation or draft, follow this sequence:

1. Silently work through the **pre-write checklist** (section 1) — if any answer is missing, surface it before drafting.
2. Identify which of the **Cialdini principles** are available to you.
3. Apply the relevant **Voss tactics** to the specific language.
4. Check the draft against the **anti-patterns** list.
5. Return a revised draft or a marked-up set of notes.

---

## 1. Pre-Write Checklist

Answer these before drafting. If you cannot answer them, your proposal is not ready.

- **My BATNA:** What will I do if this negotiation fails? Write it down. The clearer your BATNA, the less you need this deal.
- **Their BATNA:** What will they do if this fails? If their alternative is better than yours, you need to either improve your offer or improve their perception of yours.
- **ZOPA (Zone of Possible Agreement):** The range between the worst deal you'd accept and the worst deal they'd accept. If there is no overlap, no deal is possible — stop and reframe the scope.
- **Their real interests:** Not their stated position. Ask: what does a yes here solve for them? What does a no protect?
- **Legitimacy anchors:** What standards, precedents, benchmarks, or policies can you cite so your proposal feels objective rather than arbitrary? (Fisher/Ury: always negotiate on principle, not position.)

---

## 2. Cialdini's 6 Principles — Tactical Application

These are the psychological levers available in any written or verbal persuasion context. Each has a brief worked example in a written request.

### Reciprocity
Give something first — a concession, useful information, effort — before making your ask. People feel obligated to return favors.

> "I pulled together the usage data you mentioned needing for the QBR. Given you now have that context, I'd like to revisit the SLA timeline we discussed..."

Tactic: lead with an unrequested gift or piece of value. The ask comes after, not before.

### Commitment and Consistency
People honor prior commitments. Get a small yes first; larger yeses follow more easily because people want to be consistent with their past behavior.

> "Since we agreed last quarter that improving deployment speed was the top priority, this request is a natural next step toward that goal..."

Tactic: reference a prior agreement, stated goal, or earlier position the other party holds.

### Social Proof
People look to others' behavior under uncertainty. In a written request, name comparable peers, customers, or teams who have already done what you're asking.

> "Three other enterprise customers in your segment have moved to the dedicated cluster tier this quarter — here's what their scale looked like before the switch..."

Tactic: use specifics (number, names if shareable, segment) — vague social proof ("many customers") is weak.

### Authority
Expertise and credentials reduce friction. Cite data, cite your role, cite the source, or name the expert who informed your recommendation.

> "Our database performance team analyzed your ops logs and identified two specific indexes driving the slow queries. Their recommendation is..."

Tactic: distinguish between title-authority (your position) and expertise-authority (your evidence). Expertise-authority is stronger and more durable.

### Liking
People say yes to people they like or feel connected to. In writing, this means: personalize the opening, use the reader's name, mirror their language, and show you listened to a prior conversation.

> "After our call last Tuesday, I went back and re-read the issues you flagged around data locality — you were right that we underweighted that constraint..."

Tactic: acknowledge something the other person said or did before making your ask.

### Scarcity
Availability limits create urgency. Use honestly — invented scarcity backfires badly when discovered.

> "We have budget allocated for this engagement through end of Q2; after that it rolls back into the central pool and would require a re-approval cycle..."

Tactic: real constraints (deadlines, budget cycles, headcount slots) are more convincing than artificial ones. Name the mechanism, not just the deadline.

---

## 3. Chris Voss Tactics in Writing

Voss's framework comes from FBI hostage negotiation. In business writing, these tactics work because they lower the other party's defensiveness and keep them talking.

### Tactical Empathy
Demonstrate that you understand their position before making your ask. This is not agreement — it is acknowledgment.

> "I understand this ask is landing at a difficult time in your planning cycle, and that resourcing decisions like this require sign-off you may not have yet..."

Rule: state their difficulty accurately. Vague empathy ("I know this is hard") lands hollow. Specific empathy ("I know this requires CFO sign-off and your CFO review is in three weeks") lands.

### Mirroring
Repeat the last 2–3 words of their last statement as a question. In writing, this means restating their concern back to them, slightly widened, to invite elaboration.

> "You mentioned the timeline was the main constraint. ...the main constraint?"

> "You mentioned the main concern is timeline. I want to make sure I understand what 'timeline' means here — is it the delivery date, the internal review cycle, or both?"

### Labeling
Name the emotion or concern you observe: "It seems like...", "It sounds like...", "It looks like..."

> "It seems like the concern isn't the cost itself, but whether the ROI case is strong enough to justify the internal conversation..."

Rule: use "It seems like" or "It sounds like," not "I think you feel." The "It" construction feels observational, not presumptuous. After labeling, go silent (in writing: end the paragraph there).

### The Calibrated Question
Open-ended "how" or "what" questions that put the problem in the other party's hands without triggering defensiveness. Avoid "why" — it implies accusation.

> "How am I supposed to make this work with the current scope?"

> "What would need to be true for this to move forward?"

> "How does this fit with what your team is trying to accomplish?"

These questions invite the other party to solve the problem with you rather than defend against you.

---

## 4. Anchoring in Proposals

**Lead with the bigger ask, with justification, then fall back to your real position.** This is anchoring: the first number or scope shapes the entire negotiation.

- Ask for more than you expect to get.
- The anchor must include a rationale — unanchored high numbers come across as uninformed.
- Your "fallback" position should be your actual goal.

> "Ideally we'd move to dedicated clusters across all three environments to eliminate contention risk. If budget is a constraint, starting with production and preproduction would capture 80% of the benefit at roughly half the cost."

Anti-pattern: presenting only your actual position first, then negotiating down. You have already lost anchoring advantage.

---

## 5. Loss vs. Gain Framing

Kahneman and Tversky (Prospect Theory): losses feel approximately 2x more powerful than equivalent gains. Frame proposals in terms of what the other party stands to lose by not acting, not only what they gain by acting.

**Gain frame:** "This migration would reduce your p99 query latency by 40%."

**Loss frame:** "Without this migration, you're leaving 40% query latency improvement on the table every day — and the current architecture will become the bottleneck before your next growth cycle."

Rule: use both in the same proposal — gain frame for vision, loss frame for urgency. Loss frame alone sounds alarmist; gain frame alone undersells urgency.

---

## 6. Soft-No Language and the "No, AND" Construction

Rejections are unavoidable. How you phrase them determines whether the conversation continues.

**Soft no:** "I'd love to do this, but our current capacity is committed through Q3. The earliest we could scope this is Q4, when we have a team rotation coming in."

**"No, AND":** Redirect rather than just refuse. "We can't accelerate the original scope, and what we can do is deliver the highest-priority module first so you have something in production before the deadline."

The "no, AND" construction keeps the momentum. It says: I'm still working toward your goal, within the constraint.

---

## 7. Stuart Diamond — Concessions and Standards

Diamond's "Getting More" core argument: most negotiation fails because parties try to win rather than learn. Three tactics from this framework:

**Reciprocal concessions:** Concessions should always be conditional and explicitly traded. Structure: "If you can do X [their concession], we can do Y [your concession]."

> "If you can extend the implementation window by two weeks, we can include the additional index tuning pass at no extra charge."

Anti-pattern: giving a concession without receiving one. Unilateral concessions signal that you had room all along, which invites further demands.

**Small-step concessions:** Never move from position A to position C in one jump. Move from A to B, confirm their acknowledgment, then to C. Each small move feels earned and builds momentum.

**Finding their standards:** Ask what the other party's own policy or criteria say about this situation — then hold them to it. People find it hard to refuse a proposal that meets their own stated standards.

> "You mentioned your team's standard is a 6-month implementation timeline for migrations of this scale. Our proposal is well within that envelope — what would need to change for it to fit your intake process?"

---

## 8. Anti-Patterns

Avoid these in written or verbal negotiation:

| Anti-pattern | Why it fails |
|---|---|
| Ultimatums ("take it or leave it") | Eliminates the other party's ability to say yes without losing face. |
| Comparison shopping in your opening | "Competitor X offers this for less" — signals you're not committed and invites them to match the competitor and end the conversation. |
| "As we discussed" without documentation | Implies they agreed to something they may not recall agreeing to. Creates defensiveness. |
| Naming your floor first when uncertain | Destroys your anchoring position. Only name numbers when you hold the anchoring advantage. |
| Vague empathy | "I know this is difficult" without specifics reads as performative. |
| Conceding without a condition | Trains the other party that concessions are free. |

---

## 9. Quick-Reference Sequences

**For a written request (proposal, ask, counter-offer):**
1. Lead with a gift or acknowledgment (reciprocity / liking)
2. State their concern accurately (tactical empathy / labeling)
3. Anchor high with rationale
4. Offer the fallback position as a "win"
5. Name a calibrated question to invite them in
6. Close with a loss-frame urgency note if deadline is real

**For verbal negotiation prep:**
1. Write down your BATNA and their probable BATNA
2. Identify the ZOPA
3. Prepare three calibrated questions
4. Prepare two labels for their likely objections
5. Identify the legitimacy anchor you'll cite (policy, precedent, benchmark)
6. Decide your anchor number and your walkaway point

---

## 10. "Yes, and" vs "Yes, but" — improv frames for negotiation

Improv comedy taught business communication a rule that maps cleanly onto
negotiation: the response that accepts the previous turn and extends it is
the response that keeps the conversation alive. Reject the previous turn
and the scene collapses; reject the previous turn in a negotiation and the
deal collapses.

### The two frames

| Frame | What it does | When to use it |
|-------|-------------|----------------|
| **"Yes, and..."** | Accepts the other party's framing and adds to it. Builds momentum. Signals collaboration. | Discovery, brainstorming, joint problem-solving, expanding the pie. |
| **"Yes, but..."** | Acknowledges the previous turn but pivots away from it. Signals that the previous turn was inadequate. | Almost never — except deliberately, when you must reject the framing. |

### Why "yes, and" beats "yes, but"

Tina Fey (*Bossypants*, 2011) and the Second City improv canon teach that
the human ear hears "but" as an eraser: everything before "but" gets
discarded. "I hear you on the timeline, **but** we can't move that deadline"
is psychologically identical to "We can't move that deadline." The
acknowledgment was free; the rejection landed alone.

"Yes, and" keeps the acknowledgment alive: "I hear you on the timeline, **and**
the deadline is fixed for reasons I want to walk you through" treats both
clauses as load-bearing. The other party feels heard and the constraint
gets stated.

### Worked example (TAM, customer pushing for free engineering support)

"Yes, but" framing (low):
> Customer: "We need your engineering team to investigate this issue
> end-to-end, on our behalf."
> TAM: "Yes, but that's outside the scope of your current support
> contract."

The customer hears: rejection. The "yes" was a courtesy.

"Yes, and" framing (high):
> Customer: "We need your engineering team to investigate this issue
> end-to-end, on our behalf."
> TAM: "Yes, and the fastest path to that level of involvement is the
> dedicated investigation track in our Premium tier — let me walk you
> through what that engagement looks like and how it'd map to this issue."

The acknowledgment is preserved. The constraint (engineering effort lives
in a higher tier) is delivered as a path, not a wall.

### "No, AND" — Voss's variant

Chris Voss's "No, AND" construction (see §6 above) is a special case of
"Yes, and" in disguise: the refusal is paired with a forward motion. "We
can't accelerate the original scope, AND what we can do is deliver the
highest-priority module first." Same psychological effect — the
conversation stays alive.

### When "Yes, but" is the right move

Some moments require explicit redirection:
- The other party has offered a frame that is factually wrong and
  proceeding from it will lead the conversation off-track.
- The previous turn contained a hidden ask that you cannot accept and that
  a soft acknowledgment would be read as implicit agreement to.
- The other party is using "yes, and" repeatedly to grow scope past your
  walkaway point — they have weaponized the frame.

In those cases, a clean "I want to push back on the framing — here's why..."
is more honest than a coded "Yes, but."

### Pre-write checklist (adding to §1)

Before drafting any negotiation response, ask:
- Does my first word acknowledge their last move?
- If I'm about to write "but," can I write "and" instead and keep the
  substance?
- Is my constraint phrased as a path forward, or as a wall?

### References

- Fey, T. *Bossypants* (2011) — chapter on Second City "yes, and" rules.
- Sweet, J. *Something Wonderful Right Away* (1978) — origin of Second City
  improv principles.
- Voss, C. *Never Split the Difference* (2016) — the "No, AND" variant.
