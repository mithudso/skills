<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `interview-and-conversational` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: interview-and-conversational
version: "1.2.0"
updated: "2026-05-29"
category: custom
description: >
  Interview and conversational craft — question design, podcast/interview prep,
  journalism's 5W1H with deepening rules, motivational-interviewing patterns,
  behavioral interview (STAR), customer-discovery (Mom Test), and transcript cleanup.

  TRIGGER: "prep for an interview", "customer discovery call", "discovery questions",
  "5 whys", "motivational interviewing", "OARS framework", "interview script",
  "podcast prep", "behavioral interview", "STAR answer", "STAR method",
  "Mom Test", "how do I ask better questions", "transcript cleanup",
  "silence as a tool", "active listening markers", "funnel interview structure".

  SKIP: writing an essay or memo (use writing-expert); negotiation tactics or
  salary asks (use negotiation-and-persuasion); argument structure or debate prep
  (use rhetorical-frameworks-deep); conference keynote or public-speaking prep
  (use public-speaking-and-presentations); analyzing qualitative data after
  collection (use deep-research or deep-research-methods).
triggers:
  - prep for an interview
  - customer discovery
  - 5 whys
  - motivational interviewing
  - interview script
  - podcast prep
  - behavioral interview
  - STAR answer
  - STAR method
  - Mom Test
  - discovery questions
  - how do I ask better questions
  - transcript cleanup
  - OARS framework
  - active listening
skip_when:
  - writing an essay or memo (use writing-expert)
  - negotiation prep or salary asks (use negotiation-and-persuasion)
  - argument structure or debate prep (use rhetorical-frameworks-deep)
  - conference keynote or public-speaking prep (use public-speaking-and-presentations)
  - qualitative data analysis after collection (use deep-research or deep-research-methods)
related:
  - writing-expert
  - negotiation-and-persuasion
  - rhetorical-frameworks-deep
  - deep-research
  - deep-research-methods
  - public-speaking-and-presentations
---

# Interview and Conversational Craft

## When to Use This Skill

Use this skill when you need to: design interview or discovery questions, prep a podcast or show, structure a customer-discovery call, apply behavioral (STAR) patterns for hiring or job-seeking, apply the OARS motivational-interviewing framework, clean up or plan transcripts, or build a conversation prep checklist.

**If the user's request is ambiguous, ask exactly one targeted question before proceeding.** Typical clarifying questions:
- "What kind of output do you need — a question list, interview script, checklist, critique of existing questions, or coaching notes?"
- "Who is being interviewed and for what purpose (e.g., customer discovery, podcast, hiring, coaching)?"

**Output guidance by type:**
- **Question list:** 6–12 questions, roughly ordered open-to-narrow; each labeled with its purpose (discovery / clarifying / behavioral).
- **Interview script:** opening → 4–8 core questions → probe bank (3–5 follow-ups) → close; include stage directions (e.g., "pause four seconds before follow-up").
- **Checklist:** numbered items matching the Conversation Prep Checklist template below; expand each item with user-specific detail.
- **Critique of existing questions:** for each question, state the anti-pattern triggered (use Anti-Pattern Summary table), then a rewritten version.
- **Coaching notes:** paragraph form; cite the relevant framework (OARS, STAR, Mom Test, 5W1H) and explain *why* the recommended change works.

**Success criteria for Claude's output:** a good output names the framework or principle it draws on, uses specific past-tense behavioral prompts rather than hypotheticals, and can be executed by the user in a live conversation without improvisation. If the output contains a hypothetical question ("Would you...?") without explicit justification, it fails the quality bar.

**After producing output, re-read it against the goal stated in step 1 of the Conversation Prep Checklist. If the output does not address that goal, revise before sending.**

**Skip this skill** for: negotiation emails or salary asks (use `negotiation-and-persuasion`), argument or debate structure (use `rhetorical-frameworks-deep`), essay or memo writing (use `writing-expert`), or analyzing already-collected qualitative data (use `deep-research`).

---

## The Foundational Principle

Every interview is a directed conversation with a gap between what the interviewer knows and what they need to know. The job of question design is to close that gap without contaminating the data. You contaminate data by suggesting answers, signaling approval, or accepting abstract claims as if they were evidence.

---

## Question Design: Open vs. Closed

Open questions invite narration. Closed questions invite confirmation.

**Open:** "Walk me through what happened after the alert fired." Gets a story.
**Closed:** "Did you escalate the alert?" Gets yes or no.

The default should be open. Use closed questions only to confirm facts you believe you already understand: "So by 'slow' you mean response times above two seconds — is that right?"

Two-part questions ("What did you do, and why?") are a consistent anti-pattern. The interviewee answers one half, usually the easier one, and the other half disappears. Ask one question at a time.

Leading questions ("That must have been frustrating, right?") implant the emotional vocabulary you expect rather than the one they have. If they feel frustrated, they'll say so.

---

## The Five Whys

Sakichi Toyoda formalized the five-whys technique at Toyota in the early twentieth century as a root-cause discipline for manufacturing failures. The rule: when you receive an answer, ask why that answer is true. Repeat until you reach a cause you can act on or until the chain breaks because the answerer doesn't know.

In practice, the depth needed varies: sometimes two iterations reach the root; sometimes seven are needed. "Five" is a heuristic for "more than once." The technique fails if you accept the first causal explanation offered without testing it. The failure mode is stopping at a symptom and treating it as a cause.

Applied to interviews: "We had a lot of churn last quarter." Why? "Customers said onboarding was hard." Why was onboarding hard? And so on until you are talking about something specific enough to change.

---

## Journalism's 5W1H

Who / What / When / Where / Why / How

Each question type has a deepening move:

| Question | Surfaces | Deepening move |
|---|---|---|
| **Who** | Actors, stakeholders | Who else was affected? Who made the decision? Who noticed first? |
| **What** | Events | What specifically happened? What was different from normal? What changed just before? |
| **When** | Timing | What had just happened before that moment? Was the timing typical or anomalous? |
| **Where** | Context, location | Does the issue appear in other environments, regions, or teams? |
| **Why** | Intent, cause | Apply five-whys iteration; probe for real motivation behind the socially acceptable first answer. |
| **How** | Process | Ask them to show or demo rather than tell; ask what the next step would be. |

In journalism, the lede must answer all six. In discovery interviews, use 5W1H as a completeness check before closing: have you heard Who was involved, What happened, When and Where it happened, Why, and How it unfolded?

---

## Funnel Structure

Broad open → Narrow specifics → Confirm understanding.

Start wide to get the full landscape. The interviewee sets the frame; you have not yet introduced your assumptions. Then narrow with follow-up probes on the specific threads that matter most. Close the loop by summarizing and confirming: "So if I've understood this correctly..."

Skipping the wide opening and jumping to narrow specifics is a common error. You anchor on your hypothesis before the interviewee has told you what the actual shape of their experience is.

---

## Customer-Discovery Rules (Blank, Fitzpatrick)

Steve Blank's customer-development methodology ("The Four Steps to the Epiphany," 2005) holds that founders talk to customers to falsify hypotheses, not to validate them. The problem interviews come before the solution interviews.

Rob Fitzpatrick's "The Mom Test" (2013) names the core failure mode: asking questions your mother would answer encouragingly to spare your feelings. Questions like "Would you use this?" or "Do you think this is a good idea?" are Mom Test failures. They invite future opinions and social approval rather than past behavior and evidence.

The Mom Test rules:
1. Ask about their life, not your idea.
2. Ask about specifics in the past, not generalities or hypotheticals.
3. Talk less. Listen more.

**Past behavior beats future opinion every time.** "Have you ever paid someone to solve this?" beats "Would you pay for a solution?" The first is evidence. The second is a wish.

**The customer is the expert on their problem. You are the expert on the solution.** In discovery, your job is to understand their problem in their terms. Resist the impulse to explain your solution before you have fully understood the problem.

Red flag phrases that signal you are no longer in discovery mode:
- "What if we..." (you are pitching)
- "Have you considered..." (you are advising)
- "That's interesting, but..." (you are redirecting)

---

## Motivational Interviewing — OARS

Miller and Rollnick developed Motivational Interviewing (MI) originally for addiction counseling ("Motivational Interviewing: Helping People Change," 3rd ed., 2013). The OARS framework is now used across coaching, healthcare, and consulting.

**O — Open questions.** Invite elaboration. "What brings you here today?" not "Are you having trouble with X?"

**A — Affirmations.** Recognize genuine strengths or efforts. Not flattery — specific acknowledgment. "It sounds like you've been carrying this for a long time and still found a way to keep the team functional." Affirmations build safety, which is prerequisite to honest disclosure.

**R — Reflective listening.** Paraphrase what you heard. Two levels:
- Simple reflection: repeat back the content. "So the deployment failed at the last step."
- Complex reflection: name the emotion or meaning beneath the words. "It sounds like you felt responsible even though the decision wasn't yours."

Complex reflection is the harder skill. It requires you to hear what the person means, not just what they said.

**S — Summarizing.** Collect and synthesize what has been said before transitioning. "Let me make sure I have this: you saw the error Friday afternoon, escalated Saturday, and by Monday the team had already tried three different patches. Is that the sequence?"

Summaries do two things: confirm accuracy, and signal to the interviewee that you were actually listening.

---

## Active-Listening Verbal Markers

These markers invite expansion and signal engagement without implanting content:

- "Tell me more about that."
- "Say more about what you mean by [their word]."
- "So what I hear is... is that right?"
- "What happened next?"
- "How did that land for you?"
- "What made that particular moment stick in your memory?"

Avoid: "I know exactly what you mean" (you probably don't), "That's great" (evaluative), "Interesting" repeated as filler (signals you've stopped tracking).

---

## Silence as a Tool

Silence is the most underused instrument in interviewing. After asking a question, most people break silence within two seconds. The interviewee reads this as "my first answer was enough."

The four-second pause: after the interviewee finishes, wait four seconds before your next question. Uncomfortable for both parties. Reliably produces the second answer, which is almost always more honest and more detailed than the first.

Silence also works after pushback or emotional disclosure. Saying "take your time" and then waiting is not passive. It is a clear signal that depth is welcome.

Leaning into discomfort: when a question produces a deflection, the deflection itself is data. You can name it: "It seems like that's a hard one to answer." Do not fill the silence with an easier substitute question.

---

## Behavioral Interview Patterns — STAR

The STAR framework structures behavioral interview responses:
- **S — Situation:** the context and stakes
- **T — Task:** your specific responsibility
- **A — Action:** what you did, specifically (not "we")
- **R — Result:** the measurable or observable outcome

When asking behavioral questions, the prompt is: "Tell me about a time when you [competency]." This forces a specific past incident rather than a general claim about how the person operates.

Laszlo Bock ("Work Rules!", 2015, describing Google's hiring methodology) argues that behavioral questions with structured scoring rubrics are the highest-signal interview format available at scale. The anti-pattern is the hypothetical: "What would you do if...?" Hypothetical answers reveal how people believe they would behave, not how they do behave.

When answering a STAR question, the most common failure is spending too much time on S and T and too little on A and R. The interviewer is evaluating your actions and the outcome — keep the situation brief.

---

## Podcast and Show Prep

Pre-interview research: read or listen to three recent appearances by the guest. Note the three stories they tell often (they will tell them again). Prepare to go deeper or sideways off those stories rather than re-eliciting the same answers they've given ten times. The guest's best material is what hasn't been asked before.

Agree on: topic scope, rough duration, anything off-limits. Do not agree on specific questions. If the guest knows the questions in advance, you get rehearsed answers, not conversation.

The pre-interview call (15-20 minutes before recording) is where the actual rapport gets built. It also surfaces the story thread the guest most wants to tell, which you should let them tell in the recorded conversation rather than in the pre-call.

Script vs. verbatim transcript:
- Script: key questions, topic transitions, opening and close written out. You deviate freely mid-conversation.
- Verbatim transcript planning: write the full ideal flow, then treat it as dead on first contact with the guest. The value is in having thought through every transition, not in following it.

Landmines: know which questions the guest finds annoying or reductive before you ask them. You may ask them anyway, but knowing the reaction in advance is data.

---

## Interview Transcript Cleanup

Preserve disfluencies ("um," "like," hesitations) when:
- The disfluency itself signals uncertainty, searching, or emotional weight.
- The transcript is being used as raw research data.
- Authenticity of voice is part of the document's purpose (oral history, legal record).

Clean disfluencies when:
- The transcript is for publication, executive summary, or PR use.
- The subject is clearly fluent and the disfluencies are artifacts of spoken-word rhythm.
- The document is intended to represent the subject's ideas, not their verbal tics.

Never clean: factual content, word choices that carry precise meaning, hedging language the speaker intended as hedging. Changing "I think we might consider" to "We will" is not cleanup. It is fabrication.

---

## Conversation Prep Checklist

Before any substantive interview or discovery call:

1. **What is your goal?** One sentence. What must you know by the end that you don't know now?
2. **What is their goal?** Why did they agree to this? What do they want to get out of it?
3. **What is the one thing you must learn?** If the conversation gets cut to ten minutes, which question cannot be skipped?
4. **What assumptions are you testing?** List your top three. If the conversation confirms all three without friction, you were probably leading.
5. **What would falsify your current hypothesis?** If you cannot name this, you are not in discovery mode.
6. **Do you have your opening question ready?** It should be open, low-stakes, and get them talking about their own experience within thirty seconds.

---

## Anti-Pattern Summary

| Anti-pattern | Why it fails | Fix |
|---|---|---|
| Two-part question | Second part gets dropped | One question per turn |
| Leading question | Plants the answer | Remove your opinion from the framing |
| Future hypothetical | Reveals wishful thinking, not behavior | Ask about specific past incidents |
| Filling silence too fast | Cuts off the deeper answer | Four-second pause |
| Accepting abstractions | "We have a process for that" is not data | Ask for a specific example |
| Jumping to solution | Discovery ends prematurely | Keep asking "what else?" until nothing more comes |
| Yes/no when you wanted story | Closes the door | Rephrase as "walk me through" |

---

## "Yes, and" vs "Yes, but" — improv frames for interviews

Improv comedy's "yes, and" rule (Tina Fey, *Bossypants*, 2011; Sweet, *Something
Wonderful Right Away*, 1978; Second City foundational curriculum) is the
single most useful conversational discipline imported into interview craft.
The rule is simple: when the other person offers a thread, accept the
thread and add to it. Reject the thread and the conversation collapses.

### The two frames

| Frame | What it does | When to use it in an interview |
|-------|-------------|-------------------------------|
| **"Yes, and..."** | Accepts the interviewee's framing and extends it. Builds rapport and depth. | Discovery interviews, podcast conversations, OARS reflective work, customer research. The default. |
| **"Yes, but..."** | Acknowledges but pivots away. The "but" signals the previous answer was inadequate. | Rare — reserve for moments where the interviewee's framing is factually wrong and proceeding from it would derail the interview. |

### Why "yes, and" matters more in interviews than in negotiation

In negotiation, "yes, and" maintains relationship. In interviews, it
unlocks the *next layer of the answer*. The interviewee gives a surface
answer first; the second-level answer is the one the interviewer needs.
"Yes, and" gives the interviewee permission to keep going.

The interviewer's hidden enemy is the implicit "but" — the head shake, the
note-taking that signals "I'm waiting for a different answer," the follow-up
question that ignores what was just said. All of these end the thread.

### Worked example (customer-discovery interview)

Interviewee: "We tried adding more agents but the queue still grew."

"Yes, but" follow-up (low):
> "Did you also try changing the SLA targets?"

The interviewer has signaled that the agent answer was insufficient. The
interviewee will defend or retreat.

"Yes, and" follow-up (high):
> "So the agent capacity wasn't the bottleneck. What did the queue actually
> look like — was it growing at the same rate, faster, slower than before?"

The interviewer accepts the agent finding and extends into the queue
dynamics. The interviewee now describes the shape of the problem in more
detail — and may surface a constraint the agent count couldn't have fixed.

### Pairing with OARS

"Yes, and" pairs naturally with the **reflective listening** half of OARS
(see §"Motivational Interviewing — OARS"):

- Simple reflection + "and": "So the team felt the deadline was unfair —
  and how did that show up in the standup?"
- Complex reflection + "and": "It sounds like you felt responsible even
  though the decision wasn't yours — and that might be why you stayed late
  for two weeks?"

In both, the "and" invites the interviewee to confirm or correct, then
extend.

### "Yes, but" as a controlled tool

Sometimes the interviewee offers a frame you must reject:

- They have given a hypothetical, and your skill demands a specific past
  example: "Yes, and that's the general view — but I'd love a specific
  time when this happened, even if it wasn't dramatic."
- They have answered a different question: "Yes, that's useful context —
  but the thing I'm most curious about is X. Can we come back to that?"

In both, the "but" is doing real work — and it's followed by a forward
move, not a wall.

### Anti-pattern: stacked "yes, ands"

Three "yes, and" turns in a row, with no probe, signals that the
interviewer is afraid to push. The interviewee gets affirmed but never
challenged. The five-whys discipline and the four-second silence are the
counterweight — accept the thread, then dig.

### References

- Fey, T. *Bossypants* (2011) — chapter on Second City "yes, and" rules.
- Sweet, J. *Something Wonderful Right Away* (1978) — original Second City
  documentation.
- Miller, W. R. & Rollnick, S. *Motivational Interviewing* 3rd ed. (2013)
  — reflective listening compatible with "yes, and" framing.

---

## Routing Reminder

This skill covers question design, discovery interviews, STAR behavioral patterns, OARS motivational interviewing, podcast prep, and transcript cleanup. **Do not use for:** negotiation or salary asks (`negotiation-and-persuasion`), argument structure (`rhetorical-frameworks-deep`), essay writing (`writing-expert`), public speaking (`public-speaking-and-presentations`), or post-collection qualitative analysis (`deep-research`).

---

## Sources

- Fitzpatrick, Rob. *The Mom Test: How to Talk to Customers and Learn If Your Business Is a Good Idea When Everyone Is Lying to You.* 2013.
- Miller, William R. and Rollnick, Stephen. *Motivational Interviewing: Helping People Change.* 3rd ed. Guilford Press, 2013.
- Blank, Steve. *The Four Steps to the Epiphany: Successful Strategies for Products That Win.* 2005.
- Bock, Laszlo. *Work Rules!: Insights from Inside Google That Will Transform How You Live and Lead.* Twelve, 2015.
- Toyoda, Sakichi. Five-whys methodology, Toyota Production System, early twentieth century; formalized in Ohno, Taiichi. *Toyota Production System: Beyond Large-Scale Production.* Productivity Press, 1988.
