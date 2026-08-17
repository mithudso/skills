<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `email-craft` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: email-craft
description: Craft and revise general-purpose professional emails — subject-line discipline, BLUF openings, single-ask-per-email, reply-all etiquette, intro requests, follow-up cadence, breakup emails, cold outreach (non-sales), thank-you notes, send-time decisions, phone-vs-email and Slack-vs-email judgment, signature hygiene, attachments vs links. TRIGGER when the user asks to "write an email", "draft a reply", "follow up by email", "craft a thank-you email", "intro request email", "breakup email", "cold email to a non-prospect", "should I email or call/Slack this", "fix this email", "rewrite this email", "what's a good subject line", "reply-all or just reply", "when should I send this email", "is my signature too long". SKIP for: sales-and-marketing-copy (sales/marketing campaigns, prospect nurturing, conversion-driven outreach); support-ticket-writing (customer support replies, ticket comments, internal-support routing); executive-comms (CEO/SVP-level announcements, board updates, all-hands prose); slack-messaging (any Slack channel/DM/thread message); writing-expert (long-form documents, blog posts, articles, prose editing without an email frame).
---

# email-craft

## Overview

General-purpose professional email writing. Optimized for the 90% of work email that isn't sales, isn't customer support, and isn't an exec announcement — internal notes, cross-team requests, intro asks, status updates, thank-yous, follow-ups, breakups, and judgment calls about whether email is even the right channel.

The skill's job is to make every email shorter, clearer, single-purpose, and easy to act on. It enforces five non-negotiables: a precise subject line, the BLUF rule, one ask per email, a deliberate To/Cc/Bcc decision, and a "should this even be an email" check before send.

## Core Concepts

### 1. The subject-line contract

The subject line is a promise. Open it like a label on a folder: it should tell the reader (a) what the email is about, (b) what role it plays (FYI, decision, action, question), and (c) whether it's time-sensitive — in roughly that order, in fewer than 50 characters where possible.

Two-part structure: `[Topic]: [Purpose]`. The colon does the work.

- "Q3 forecast: need your number by Thu noon"
- "Onboarding doc: review request (10 min)"
- "Reply-all on the deploy thread: please don't"
- "Intro: Alice <> Bob (cybersec)"

Prefixes that earn their keep: `Action:`, `Decision needed:`, `FYI:`, `Question:`, `Heads up:`, `Reminder:`, `Intro:`. Avoid `Following up`, `Checking in`, `Quick question`, `Hi`, `Update`. They tell the reader nothing and feel passive-aggressive when stale.

Length: 2-7 words wins. Open rates collapse past 7. Mobile clipping starts around 33 characters.

### 2. BLUF: Bottom Line Up Front

The first sentence states the conclusion, the ask, or the decision needed. Context follows. Reasoning follows context. Background last (or in a thread, or in a link).

```
Bad:  "Hey, hope you're well. I've been thinking about the Q3 forecast,
       and after reviewing the numbers from last quarter alongside the
       new pipeline data, I wanted to check in on..."

Good: "I need your Q3 forecast number by Thursday noon. Context below."
```

If the reader stops reading after sentence one, did they get the most important thing? If no, rewrite sentence one.

### 3. One ask per email

If the email has two asks, it has zero asks. Readers triage on the first one and forget the second. Split into two emails, each with its own subject line. The exception is a status update with no ask at all — that's allowed, but mark it `FYI:` so nobody hunts for a hidden question.

Test: can you summarize the ask in one sentence starting with a verb? "Please review X by Friday." "Can you forward this to Y?" "Please confirm you can attend Tuesday's review." If you can't, the ask isn't clear yet.

### 4. To / Cc / Bcc: a deliberate decision

- **To**: people who must act or reply. If they don't owe you a response, they don't belong on the To line.
- **Cc**: people who need to be informed but don't owe a reply. Use sparingly. Every Cc multiplies the email's social weight.
- **Bcc**: only two legitimate uses — (1) moving someone off a thread after introductions are done ("moving Alice to Bcc"), (2) hiding a large external distribution list. Bcc'ing your boss to surveil a colleague is corrosive and almost always gets discovered.
- **Reply-all**: default to No. Reply-all is correct when the answer materially changes what the other recipients should do. "Thanks!", "Got it!", "Sounds good!" are never reply-all answers.

### 5. The "would this fit in Slack" test

Before sending, ask: could this be a Slack DM or channel post? If yes, and the reader is a Slack-active colleague, send Slack. Email is for:

- Anything you want a searchable, durable record of (commitments, decisions, agreements).
- External recipients who don't share your Slack workspace.
- Long-form content that needs structure (>3 paragraphs or a list of asks).
- Anyone whose Slack hours don't overlap yours.
- Anything that warrants a CC chain across orgs.

Slack is for: quick questions with same-day decay, thread-shaped collaboration, casual heads-ups, anything you'd hate to see archived for years.

### 6. The "phone vs email" judgment

Pick up the phone (or schedule a 10-minute call) when:

- The topic is emotional, sensitive, or carries firing/conflict/HR risk.
- You're on email round 3+ on the same thread without resolution.
- You need a real-time decision and the recipient is fast on phone but slow on email.
- The misread risk is high (tone-sensitive feedback, layoffs, performance, breakups).
- Cultural or trust capital is at stake — apology, condolence, a hard ask.

Email is fine for: structured asks, anything that needs a paper trail, anything async-friendly, anything with a deadline more than a day out.

A useful heuristic: if you've drafted the email three times and it still feels wrong, the answer is a call.

### 7. Follow-up cadence

Default cadence for an unanswered ask:

- Day 0: original email
- Day 3-4 (business days): polite nudge, reply on the same thread, no new subject line, BLUF the ask again in one sentence
- Day 7-10: third nudge, change channel (Slack DM, calendar invite, walk by the desk)
- Beyond: stop. The non-answer is the answer. Decide whether to escalate, route around, or drop the ask.

Never send a fourth email-only nudge — it reads as desperate, and the issue is no longer the email; it's the relationship or priority.

### 8. The breakup email

Used to gracefully close a thread that has gone nowhere — vendor calls you don't want, conversations that lost relevance, asks the other party has clearly chosen not to engage with. Pattern:

1. Acknowledge no response.
2. Make it the recipient's call to revive (low pressure).
3. State your default action ("I'll assume X and move on").
4. Leave the door open without begging.

```
Subject: Closing the loop on [topic]

Hi [name] — I haven't heard back on this, so I'll assume it's not a
priority right now and stop following up. If that changes, drop me a
line and I'll pick it back up. Thanks for considering it.
```

### 9. Cold outreach (general-purpose, non-sales)

You want a 10-minute call with a stranger — for a career question, a research interview, a partnership exploration, an intro to someone in their network. Rules:

- Subject line names the connection or the specific reason. ("From the Rust Belt podcast — quick research question")
- First sentence states who you are and how they're connected to your ask (one sentence).
- Second sentence is the ask itself, with a quantified time cost. ("Could I get 15 minutes on your calendar in the next two weeks?")
- One paragraph of context max.
- Make it easy to say no. ("Totally fine if not — no need to reply.")
- No attachments on a cold email. Link to a resume / project / GitHub / portfolio if relevant.

### 10. Intro request emails

Two flavors, very different mechanics.

**Asking for an intro** (you want a connector to introduce you to a third party):

- Write a "forwardable email" the connector can paste with one click.
- Make the ask in the third party's language, not yours.
- Give the connector an explicit out: "Totally fine if this isn't a fit — feel free to pass."

**Receiving / making an intro** (the double opt-in):

- Always double-opt-in before connecting two strangers. Ask both sides if they're open to it before sending the intro mail.
- The actual intro mail: two sentences on each person, why you think they'd benefit, then "I'll let you two take it from here."
- Both parties reply-all, then move the connector to Bcc to free them from scheduling churn.

### 11. Thank-you notes

Send within 24-48 hours of the favor. Specific beats generic. Name the thing, name the impact, name what's next.

```
Subject: Thank you — the X review made the difference

The notes you sent on the API design saved us a full redesign cycle.
We shipped the v2 spec yesterday with your concurrency concerns baked
in. I owe you a coffee — and a draft of the v3 RFC when it lands.
```

Avoid: "Just wanted to say thanks!" (lazy), "Couldn't have done it without you!" (vacuous), thank-yous that smell like setups for the next ask.

### 12. Send-time decisions

Inbox attention follows a daily curve. Tuesday-Thursday, 10am-2pm recipient-local time, is the highest-attention window in most studies. Monday morning is buried under weekend backlog. Friday afternoon is read-and-archived without action. Weekends are read-and-deferred-to-Monday, by which point your email is buried.

Decision rule:

- Internal recipient, ongoing relationship → send when drafted. Async wins beat clever scheduling.
- External recipient, senior leader, or first contact → schedule for Tue-Thu 10am recipient-local.
- Drafted after 9pm local → always schedule for 9am next workday. Never send a time-stamp that signals "I'm working at midnight and so should you."

### 13. Signature hygiene

A good signature is 4 lines or fewer:

```
First Last
Title, Team
Company  |  Phone (if relevant)
```

What to cut: motivational quotes, every social handle, pronoun + title + team + team motto + 6 logos + a banner ad for the company offsite. Pronouns are fine on one line if your culture uses them. Mobile signatures should say "Sent from mobile — short replies are short" or nothing — never "Please excuse typos sent from my Galaxy Ultra Quantum."

Internal email after the first reply: drop the signature entirely. Repeating it on every reply is noise.

### 14. Attachment vs link

Default to a link. Attachments are for:

- External recipients on org boundaries where shared drives don't reach.
- Final signed PDFs and legal artifacts you need to preserve in the recipient's inbox.
- Files under 1 MB that the recipient will want offline.

Everything else: shared-drive link with view or comment permissions verified before send. Three lifetime checks before clicking send on any attachment: (1) is the permission right? (2) does the filename include the version or date? (3) would I be horrified if it landed in the wrong inbox?

### 15. Tone, register, and the read-aloud test

Read every email aloud before sending. If a sentence sounds passive-aggressive, it is. If it sounds curt, soften it with one warming word. If it sounds like a corporate press release, you're hiding behind formality.

Calibration sliders:

- "Just" — usually deletable. ("I just wanted to ask..." → "Could you...")
- "Sorry to bother you" — deletable. Substitute respect for the reader's time with brevity.
- "Per my last email" — passive-aggressive, regardless of intent. Quote the relevant line directly instead.
- "Hope this finds you well" — fine once per recipient per quarter; tired after that.
- "Thanks in advance" — fine when sincere, presumptuous when used to pre-extract a yes.

Match the register the recipient set, half a notch warmer. Mirror is the rule.

## Templates and Examples

### Template: action-required ask

```
Subject: [Topic]: [verb] needed by [date]

Hi [name],

I need [the specific thing] by [date/time]. Here's the context: [1-3
sentences max]. The link / file: [link]. Reply with [the shape of the
answer you need] — happy to jump on a 10-minute call if that's faster.

Thanks,
[name]
```

### Template: FYI status update (no ask)

```
Subject: FYI: [Topic] status as of [date]

Quick status update, no action needed:

- [What just happened]
- [What's next]
- [Anything blocking, with owner]

Full details: [link to doc or thread].

— [name]
```

### Template: forwardable intro request

```
Subject: Quick favor — intro to [person]?

Hi [connector],

If you think it's a fit, would you be open to introducing me to
[person] at [company]? I'm [one-line context on why I'd be useful to
them and them to me].

If easier, here's a forwardable below — feel free to paste, edit, or
ignore. Zero pressure.

—

[forwardable starts here]
Hi [person] — I'm working on [thing] and would love 15 minutes of
your time on [specific question]. [Connector] suggested you might
have a useful angle. Open to a quick call in the next two weeks?
```

### Template: polite follow-up nudge

```
Subject: (reply on existing thread, do not change subject)

Hi [name] — bumping this in case it slipped. The ask is still: [one
sentence]. No rush if it's not a priority — let me know if you'd
prefer I route this elsewhere.
```

### Template: breakup / close-the-loop

```
Subject: Closing the loop on [topic]

Hi [name] — I haven't heard back, so I'll assume [topic] isn't a fit
right now and stop following up. If that changes, just drop me a
line. Appreciate you considering it.
```

### Template: thank-you with named impact

```
Subject: Thank you — [the specific thing] made the difference

[Specific named impact, 1-2 sentences]. I'm grateful — and [the next
thing I'll do that compounds the favor, if any]. Owe you a [coffee /
drink / favor in kind].
```

## Anti-Patterns

- **The wall of text** — five paragraphs, no whitespace, no headings, no BLUF. Recipient bails at line 3.
- **The hidden ask** — buried in paragraph 4, after backstory. If you wanted it answered, you'd have led with it.
- **The two-ask email** — readers pick one, drop the other, both stall. Always split.
- **Reply-all "thanks!"** — multiplies inbox noise by N for zero value. Reply just the sender.
- **Subject-line drift** — "RE: RE: FW: RE: Q3" after the conversation has pivoted to Q4. Change the subject when the topic changes; start a new thread.
- **Per my last email** — quote the line instead, or pick up the phone.
- **The midnight send** — signals burnout culture. Schedule for next morning.
- **The motivational-quote signature** — five lines of motivational quote = five lines of "I take myself seriously and so should you."
- **The PDF as the first contact** — a heavy attachment on a cold email is an instant delete.
- **Bcc'ing the boss as surveillance** — corrosive, almost always discovered, ends careers.
- **The "circling back" loop** — third nudge with no new information. The non-answer is the answer; change channel or drop it.
- **The apology pre-amble** — "So sorry to bother you, I know you're super busy, I hate to ask but..." Two lines of nothing. Just ask.

## Decision Heuristics

**Should I send email or Slack?**
- Needs paper trail, external recipient, or > 3 paragraphs → email
- Same-day decay, internal, casual → Slack
- Senior leader who lives in email → email even if Slack would feel natural

**Should I send email or call?**
- Emotional / sensitive / HR-adjacent → call
- Round 3+ with no resolution → call
- Tone-misread risk high → call
- Structured ask with no emotional content → email

**Should I reply or reply-all?**
- Reply unless the response materially changes what others on the thread should do
- "Thanks!" / "Got it!" → never reply-all
- New decision that re-routes the project → reply-all

**Should I send now or schedule?**
- Internal, ongoing relationship → send now
- External, senior, or first contact → schedule Tue-Thu 10am recipient-local
- Drafted past 9pm → schedule for 9am next workday

**Should I attach or link?**
- Default link
- External org boundary + needs to live in their inbox → attach
- Signed legal / final PDF → attach
- Anything still editable → link

**Should I CC them?**
- Only if they're worse off not knowing
- Cc multiplies social weight; budget Cc lines as if each one cost a dollar

**Should I write the third follow-up?**
- No. Change channel, escalate, or close the loop. Three email-only nudges max.

## The Five Cs of email — the canonical quality check

The Five Cs (Clear, Concise, Complete, Correct, Courteous) is a long-standing business-communication rubric — taught by the British Council, professional-writing programs, and corporate communications curricula for decades. It is not a writing framework; it is a *pre-send checklist*. Run it on any email that matters (anything to a stakeholder, an executive, a customer, or anyone external).

### The five Cs, in send order

1. **Clear.** Can the reader understand the message on a single read? Vague antecedents ("this", "that", "it" without a referent), passive voice that hides the actor, multi-clause sentences that bury the verb — these are the clarity killers. The clarity test: would someone skimming on a phone get the meaning right?
2. **Concise.** Is every sentence earning its place? Adverbs ("really", "actually", "basically"), hedge phrases ("I just wanted to", "I think maybe"), and warm-up clauses ("Hope this finds you well, and I wanted to circle back on…") add length without adding signal. Concision is respect for the reader's time, not brevity for its own sake.
3. **Complete.** Does the email contain everything the reader needs to act — the ask, the deadline, the owner, the relevant context, the path forward if blocked? Incomplete emails generate reply threads. Reply threads generate confusion. A complete email pre-empts the three follow-up questions the reader would have asked.
4. **Correct.** Spelling, grammar, names, titles, numbers, dates, attachments. A misspelled recipient name in the salutation costs more credibility than a typo in the body. Triple-check names, numbers, and the attached file (the "I forgot the attachment" follow-up is the most common Five-Cs failure).
5. **Courteous.** Is the tone proportionate to the relationship and the stakes? Courteous is not the same as effusive. Courteous means respecting the reader's standing, time, and context. A blunt note can be courteous; a saccharine one rarely is.

### How the Five Cs interact

The Cs are not independent. They are in tension, and choosing well between them is the craft:

- **Concise vs Complete.** Short emails risk omitting context; complete emails risk bloat. Resolve by ranking: complete trumps concise when an action is requested; concise trumps complete in FYI threads.
- **Courteous vs Concise.** A one-line "no" to a senior request reads as curt. A three-line "no" with one sentence of context lands as courteous and still concise.
- **Clear vs Courteous.** Hedge language ("I was thinking maybe we could possibly…") is often a courtesy reflex that destroys clarity. The fix: state the position clearly, then add one short sentence of context for warmth.

### Worked example — applying the Five Cs

**Draft (fails 4 of 5):**

> "Hi team, hope everyone's having a good week! I just wanted to reach out and circle back on the conversation from our last meeting regarding the Q3 forecast. I think it would be really helpful if we could maybe figure out a time to talk about the numbers and possibly identify some next steps. Let me know what you think when you get a chance!"

Diagnostic:
- *Clear?* No. What numbers? What conversation? What next steps?
- *Concise?* No. 60 words to say "schedule a meeting."
- *Complete?* No. No proposed time, no owner, no agenda.
- *Correct?* Presumably yes (no errors).
- *Courteous?* Performatively so — the warmth is filler, not respect.

**Rewrite (passes 5/5):**

> "Hi team — I need 30 minutes this week to lock the Q3 forecast. Thursday 2pm or Friday 10am work for me. Agenda: walk through the three line items still open from last week's review, decide owners, and confirm the number we report on Monday. Let me know which slot works, or propose another."

40 words. Clear (specific subject, specific time options). Concise (no filler). Complete (subject, ask, options, agenda, fallback). Correct (presumably). Courteous (acknowledges the reader's calendar with options).

### When to break the Cs

- **Disaster emails.** When something has gone wrong, courtesy and completeness override concision. Apologize properly; explain fully; offer the next step.
- **Outbound to people who don't know you.** Courteous goes up; concise stays the same; complete goes down (don't dump your full context on a cold reader).
- **Internal humor and team rapport.** A team that has worked together for years can drop the courtesy-padding without seeming curt — context is doing the work the words used to do. Don't apply the Five Cs as a corporate-stiffness mandate.

### Diagnostic — the pre-send Five-C scan (30 seconds)

Before clicking send, scan in this order:

1. Read the subject line. Does it carry the email's purpose?
2. Read sentence one. Is the ask, decision, or answer there (BLUF)?
3. Scan for adverbs and hedges. Cut what isn't earning.
4. Scan recipient names, attached files, dates, and numbers for correctness.
5. Read the closing. Is it proportionate — neither curt nor saccharine?

A passing scan ships. A failing scan goes back to draft.

### References

- [The 5 Cs of email writing — British Council](https://www.britishcouncil.in/blog/5-cs-email-writing) — canonical introduction; covers all five with examples.
- [Why the 5 C's are Still Relevant for Emails — Mamma](https://email.mamma.com/why-the-5cs-are-still-relevant-for-emails) — modern reframing for hybrid-work email volume.
- [5 C's for Effective Email Communication — Soft Skills Guide](https://softskillsguide.com/2024-soft-skills-guide/c2-business-communications/s1-foundational-communications-skills/l2-5cs-to-effective-email-communication/) — workplace-communication context and pre-send checklist framing.

## References

1. [The Essential Guide to Crafting a Work Email — Harvard Business Review](https://hbr.org/2015/07/the-essential-guide-to-crafting-a-work-email)
2. [5 Tips for Writing Professional Emails — Harvard Business Review](https://hbr.org/2022/08/5-tips-for-writing-professional-emails)
3. [Ask a Manager — Alison Green](https://www.askamanager.org/)
4. [Email Etiquette and the Perils of "Reply All" — Harvard Business Publishing](https://store.hbr.org/product/email-etiquette-and-the-perils-of-reply-all/H007EV)
5. [Email Subject Line Best Practices for 2026 — Truelist](https://truelist.io/blog/email-subject-line-best-practices)
6. [40 Email Best Practices Every Professional Needs to Know — EmailAnalytics](https://emailanalytics.com/email-best-practices/)
