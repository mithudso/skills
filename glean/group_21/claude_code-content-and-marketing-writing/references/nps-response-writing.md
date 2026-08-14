<!-- hub-reference-banner -->
> **Reference file — part of the `content-and-marketing-writing` hub.** Formerly the standalone `nps-response-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: nps-response-writing
description: Writing responses to customer feedback at scale — NPS verbatims, CSAT comments, support survey replies, and app-store reviews. Covers the thank/acknowledge/route/commit template, detractor vs passive vs promoter calibration, closed-loop response cadence, the "I'll route this to engineering" anti-pattern, escalation-trigger language for save calls, confidentiality vs name-attribution choices, and platform-specific app-store reply patterns. TRIGGER: "respond to NPS", "write back to a detractor", "reply to this app-store review", "NPS verbatim response", "CSAT comment reply", "closed-loop email", "save-call language", "what do I say to a 4-star reviewer", "promoter thank-you that doesn't sound canned", "passive nudge", "customer survey reply template". SKIP: reactive support ticket replies for an active issue (use support-ticket-writing); marketing copy or campaign emails (use sales-and-marketing-copy); executive communications and board-level customer updates (use executive-comms); designing survey questions themselves rather than replying to them (use survey-builder or survey-question-writing); generic prose editing (use writing-expert).
---

# NPS Response Writing

## Overview

This skill covers how to write replies to customer feedback collected at scale — NPS verbatim comments, CSAT free-text, post-support survey replies, and public app-store reviews. The job is not to "respond to a ticket." The job is to close the loop on a structured feedback signal: acknowledge the specific point the customer made, commit to a specific next step (or explicitly decline to), and route the signal where it belongs internally without burying the customer in process.

Three things make these responses different from support replies:

1. **The customer did not ask for help.** They volunteered feedback. A reply that treats the verbatim as a support ticket misreads the contract.
2. **The score is the spine.** Detractor (0–6), passive (7–8), and promoter (9–10) responses calibrate differently. Same template across all three reads as form-mail.
3. **Cadence is short.** Under 48 hours for detractors, under 72 hours for app-store replies. After 72 hours, recovery rates fall off sharply and replies start to feel like form letters arriving from the void.

Use this skill when drafting individual replies, when writing reply templates a team will reuse, or when calibrating a closed-loop program's response style. Skip it for active support tickets, marketing campaigns, or designing the survey instrument itself.

## Core Concepts

### 1. The thank/acknowledge/route/commit template

Every NPS-style reply has four moves in this order:

1. **Thank** — one line, specific to the act of giving feedback, not the score. ("Thanks for taking the time to write that out" beats "Thanks for being a customer.")
2. **Acknowledge** — quote or paraphrase the specific point. Proves you read it. One sentence.
3. **Route** — name what happens next internally. Be specific: "I'm sending this to our billing team's weekly review on Thursday" beats "I'll route this to the right team."
4. **Commit** — one concrete next step you (the responder) own, with a date or trigger. If you cannot commit anything, say that explicitly and explain why; do not fake a commitment.

Skipping any of the four reads as evasive. The "Acknowledge" step is the one most often skipped, and the most damaging to skip.

### 2. Detractor vs passive vs promoter calibration

| Tier | Score | Customer state | Response goal | Tone | Length |
|------|-------|----------------|---------------|------|--------|
| Detractor | 0–6 | Frustrated, considering churn or already churning | Understand the specific failure; offer a save conversation when any save-call trigger is present (see Core Concept 5) | Plain, accountable, no marketing voice | 4–8 sentences |
| Passive | 7–8 | Satisfied but unenthusiastic, easily poached | Convert: ask one targeted question about what would have made it a 9 | Curious, low-friction, no apology | 3–5 sentences |
| Promoter | 9–10 | Enthusiastic, possible reference/case study | Deepen the relationship; harvest specifics for marketing/CS | Warm, specific, no upsell | 2–4 sentences |

The dominant calibration error is treating passives like detractors (over-apologizing and triggering reflective frustration) or treating detractors like passives (under-acknowledging and confirming the score).

### 3. The "I'll route this to engineering" anti-pattern

"I'll pass this along to the team" is the most common closed-loop failure mode. It signals that the responder has no agency and that the feedback is going into a black box. Replace it with one of:

- A named destination: "I'll add this to our March platform-stability review."
- A named owner: "Our billing PM, Priya, will see this on Monday."
- A named mechanism: "I'm filing this as part of our quarterly detractor synthesis that the product team reads end-to-end."
- An explicit decline: "I can't promise this gets prioritized — we're locked on Q2 commitments — but I've logged it under [theme] so it surfaces in our next planning cycle."

Specific routing builds trust even when the answer is "no." Vague routing destroys trust even when the answer is "yes."

### 4. Closed-loop response cadence

- **Detractor first response:** within 24–48 hours. Past 72 hours, the customer has already told three colleagues and the recovery window narrows sharply.
- **Detractor save call (if warranted):** within 5 business days of first response.
- **Passive reply:** within 1 week.
- **Promoter reply:** within 2 weeks; a delayed promoter reply is rarely fatal.
- **App-store reply:** within 72 hours; the top quartile of publishers reply within 30 minutes.
- **Loop close-out** (confirming resolution back to the customer): within 30 days of original feedback. This is the step most programs skip; closing the loop here is what produces score lift on re-survey.

Measure **loop closure rate**, not response rate. Sending a reply is not closing the loop.

### 5. When an NPS verbatim becomes a save call

Escalate from a written reply to a live conversation (phone, video, in-person) when any of these appear in the verbatim:

- Explicit churn language: "evaluating alternatives," "renewal is up," "we're moving to [competitor]."
- A specific named individual who is angry (often a champion or budget-holder).
- A claim of broken commitment ("you promised us X").
- Compounding issues ("this is the third time…").
- Revenue-at-risk signal in B2B: an enterprise account, a multi-product customer, or a strategic logo.
- Legal/compliance language: "data breach," "GDPR," "contract violation," "lawsuit."

The escalation language in the written reply itself should be plain: "This sounds like more than I can resolve in email — can I get 15 minutes with you this week?" Avoid the words "jump on a call," which most enterprise buyers find irritating.

### 6. Confidentiality vs name attribution

Default to **confidential treatment** of the verbatim unless the customer has signaled otherwise. Specifically:

- Do not quote a customer's verbatim back to them in marketing copy without explicit written permission.
- Do not attach a customer's name to a feature request that gets shared in internal forums where sales reps or partners may see it.
- Do attach the name internally when routing to the account owner or CSM, so they can contextualize.
- For app-store replies, never reference internal account data ("I see you're on the Pro plan") — the reply is public.
- If you want to use the verbatim as a quote or case study, ask in the reply: "Would you be open to us sharing a version of this with our product team / on our site? Happy to edit it down with you first."

### 7. App-store review response patterns

Public replies differ from private NPS replies in three ways:

1. **Audience is multiple.** Future browsers read the reply too. Write for them, not only for the reviewer.
2. **No account context.** You cannot assume you know who they are or what they paid.
3. **Algorithmic weight.** Both Apple and Google Play favor responsive publishers; iOS 18.4+ AI-generated review summaries pull from review text, so individual replies indirectly shape the public summary.

For 1- and 2-star reviews: acknowledge specifically, name the fix or workaround if there is one, give a non-public channel for follow-up (support email, in-app help), and avoid blanket apologies. For 4- and 5-star reviews: short thank-you, name the specific feature they praised if they named one, do not ask for an upgrade or a referral in the reply.

### 8. The honest-decline reply

Sometimes the right answer is "we are not going to fix this, and here is why." This is harder to write than a deflection, but it converts detractors at a higher rate than vague promises do. Structure:

1. Acknowledge the specific point.
2. State the trade-off plainly: "We deliberately chose X over Y because [reason that respects the customer's intelligence]."
3. Offer the closest workaround you can.
4. Leave the door open: "If this stays a blocker for you, I want to know — that's the signal that moves it up the list."

Honest declines build more trust than performative empathy.

### 9. Before-sending self-check (and the ambiguity rule)

Before sending any drafted reply, re-read it once against this five-point check. If any answer is "no," revise before sending.

1. Does the first non-greeting line acknowledge a specific thing the customer said (not just thanks for feedback)?
2. Is the routing destination named (person, meeting, mechanism), not "the team"?
3. Is there exactly one concrete next step that the responder owns?
4. Does the tone match the tier (no apology to a passive, no upsell to a promoter, no marketing voice to a detractor)?
5. Could a colleague reading the reply tell which tier the customer is in without seeing the score?

**Ambiguity rule:** if the verbatim is unclear about what the customer is upset about, satisfied with, or asking for, do not guess. Draft the reply with one targeted clarifying question ("When you said 'the dashboard feels heavy,' do you mean slow to load, hard to navigate, or visually overwhelming?") and skip the route/commit moves until the customer replies. Guessing wrong on a detractor verbatim is the single fastest way to confirm the score.

## Templates and Examples

### Detractor (NPS 3, B2B SaaS verbatim: "Support takes forever and your billing portal is broken")

> Hi Sarah,
>
> Thanks for being direct in the survey — that's more useful to us than a polite 7.
>
> Two specific things you flagged: support response time, and the billing portal. On support: I pulled your last four tickets and the median first-response was 19 hours, which is outside our SLA for your plan. That's on us. On billing: the portal had a known issue with invoice downloads through last Tuesday; it's fixed now, but you would have hit it.
>
> I'm doing two things. First, I've flagged your account to our support lead, Marco, who will personally route your next three tickets and confirm response time. Second, I'd like 20 minutes with you this week or next — not a sales conversation, just to hear what "fast enough" looks like for your team, so I can take that back to our support planning.
>
> Tuesday or Thursday afternoon work?
>
> — Mitch

### Passive (NPS 8, consumer app verbatim: "It's good, just nothing exciting")

> Hi Alex,
>
> Thanks for the note — "good but not exciting" is genuinely useful for us, more than a 10 with no comment would be.
>
> If you don't mind a quick question: what would have made it a 9 or 10 for you? One thing, off the top of your head, is enough.
>
> No pressure to reply — but if you do, it goes to the product team's monthly read.
>
> — Mitch

### Promoter (NPS 10, B2B verbatim: "Your onboarding team saved our launch")

> Hi Priya,
>
> Thank you — I forwarded your note to Janet's onboarding team this morning and they were genuinely chuffed. The "saved our launch" line is the one they're going to remember.
>
> If you're ever open to a 20-minute conversation with our product team about what worked, I'd love to set it up — no marketing ask attached. Either way, glad the launch landed.
>
> — Mitch

### App-store 1-star reply (consumer iOS)

> Sorry the import flow blocked you on a 200MB file — that's a real bug, not the experience we want. We shipped a fix in v4.7.2 last week that raises the limit to 2GB and adds a progress indicator. If it's still failing for you, email support@[app].com and reference review code [if platform supports] — we'll look at it directly.

### App-store 5-star reply (consumer iOS)

> Thank you — and specifically thanks for calling out the dark-mode redesign, that was a six-week project and the team will love hearing this landed.

### Honest decline (Detractor verbatim: "Why don't you support [feature X]?")

> Hi Jamie,
>
> Fair question, and the short answer is we made a deliberate trade-off. We chose to invest in [Y] this year instead of [X] because our usage data showed [Y] blocking 4x as many customer workflows. That's a real cost to you, and I'm not going to pretend otherwise.
>
> The closest workaround today is [workaround], which gets you most of the way there but not all of it. If that's not enough — and the gap is costing you real time — please tell me. That's the signal that moves [X] up the priority list, and I want it on record.
>
> — Mitch

## Anti-Patterns

- **The form-letter thank-you.** "Thank you for your valuable feedback. We are always striving to improve." This reads as bot output and confirms the detractor's view.
- **The unprompted apology.** Saying "I'm sorry" to a passive (NPS 8) signals you read them as a detractor, which is mildly insulting.
- **The route-to-the-void.** "I've passed this along to the team" with no named destination, owner, or mechanism.
- **The re-survey ask before resolution.** Asking a detractor to re-rate before you have fixed anything. This reads as score-gaming.
- **The marketing voice.** Brand-voice phrases ("we're committed to delivering delight") in a reply to a frustrated customer. Drop the voice; use plain English.
- **The upsell smuggled into a promoter reply.** Promoters are not pre-qualified leads. A reply that pivots to an upgrade ask burns the goodwill.
- **The public airing of private context.** Quoting account data, plan tier, or internal ticket IDs in a public app-store reply.
- **The over-promise.** Committing to a fix you cannot actually ship. Better to honest-decline than to over-promise and miss.
- **The "jump on a call" reflex.** Some customers want a written response; do not escalate to a call as a default, only when the verbatim signals it.
- **The 14-day-late reply.** A reply that arrives outside the recovery window often reactivates the frustration rather than closing the loop.

## Decision Heuristics

- **Score < 7 + churn language → save call, not email.** Skip the written reply and go straight to a 15-minute live conversation request.
- **Score < 7 + no churn language → written reply within 48h + named-owner routing.**
- **Score 7–8 → one targeted question, no apology.** Ask what would have made it a 9. Stop there.
- **Score 9–10 + specific praise → name the team or person they praised in the reply.** It is the highest-leverage warmth move available.
- **Score 9–10 + generic praise → short thank-you, optional reference-program ask.** Do not over-fish.
- **Public review (app store, Trustpilot, G2) → write for the next reader, not only the reviewer.** Assume the reply is read by 100x more people than the reviewer.
- **Verbatim mentions a named individual on your team → tell that person before you reply.** Both for accuracy and for morale (positive or negative).
- **Verbatim mentions a competitor by name → do not name the competitor back.** Acknowledge the alternative exists, address the gap, move on.
- **You cannot commit a fix → use the honest-decline template, not vague-routing.**
- **You are not the right responder (wrong account team, wrong region) → hand off explicitly with the customer cc'd**, do not silent-forward. "Looping in Janet who owns your account — Janet, the context is in the survey link."

## References

- [Introducing the Net Promoter System | Bain & Company](https://www.bain.com/insights/introducing-the-net-promoter-system-loyalty-insights/) — Reichheld and Bain on the inner-loop closed-loop feedback model.
- [The Net Promoter System's "Inner Loop": Individual Learning and Connections with Customers | Bain & Company](https://www.bain.com/insights/the-net-promoter-systems-inner-loop/) — How frontline employees should close the loop one customer at a time.
- [Net Promoter 3.0 | Bain & Company](https://www.bain.com/insights/net-promoter-3-0/) — Reichheld's updated framing on earned growth and loop closure as the operational test.
- [Net Promoter Score (NPS): The Ultimate Guide | Qualtrics](https://www.qualtrics.com/articles/customer-experience/net-promoter-score/) — Promoter/passive/detractor segment definitions and verbatim-analysis patterns.
- [What is a Detractor and How Do You Turn Them into Promoters? | Qualtrics](https://www.qualtrics.com/experience-management/customer/detractors/) — Detractor recovery mechanics.
- [Net Promoter Score Guide | Medallia](https://www.medallia.com/net-promoter-score/) — Medallia's CX-program framing on verbatim categorization and loop closure.
- [Closing the Customer Feedback Loop | Bain & Company](https://www.bain.com/insights/closing-the-customer-feedback-loop-newsletter) — Bain on cadence, ownership, and the difference between response rate and loop closure rate.
- [Ratings, Reviews, and Responses | Apple Developer](https://developer.apple.com/app-store/ratings-and-reviews/) — Apple's guidance on responding to App Store reviews; tone and platform constraints.
- [How to respond to app store reviews the right way in 2026 | MobileAction](https://www.mobileaction.co/guide/how-to-respond-to-app-store-reviews/) — Current best-practice timing and content patterns for public app-store replies.
- [Closed Loop Feedback (CX) Best Practices + Examples | CustomerGauge](https://customergauge.com/blog/close-the-loop) — B2B-oriented closed-loop cadence, save-call escalation, and recovery-window data.
