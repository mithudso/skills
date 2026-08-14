<!-- hub-reference-banner -->
> **Reference file — part of the `content-and-marketing-writing` hub.** Formerly the standalone `chatbot-conversation-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: chatbot-conversation-writing
description: Write conversational AI prompts and responses — support bots, sales chat, in-product help bots, AI assistant personas, embedded chat agents. Covers Cooper-Reeves conversational design, persona consistency, graceful-confusion patterns (when the bot doesn't understand), turn-taking conventions, the post-2024 "no, I'm not human" disclosure rule, brand-voice transfer into a bot, fallback hierarchies, and escalation-to-human prose. References Erika Hall's Conversational Design, Voiceflow and Botpress design docs, and 2024–2026 chatbot UX research. TRIGGER -- "write a chatbot persona", "draft chatbot responses", "bot fallback messages", "AI assistant copy", "conversation design", "support bot script", "in-product help bot", "chatbot escalation prose", "AI disclosure copy", "turn design for a chatbot". SKIP -- voice / IVR / phone-tree scripts (use audio-script-writing); static UI labels, buttons, tooltips, error toasts (use microcopy-and-ui-writing); human-to-human support ticket responses (use support-ticket-writing); LLM system-prompt engineering for general agents (use prompt-engineering); choosing which LLM model powers a bot (use llm-models or anthropic-sdk).
category: custom
tags: ["writing", "conversational-design"]
---

# Chatbot Conversation Writing

## Overview

A chatbot is a conversation, not an interface. Every turn the bot takes
should advance the user's goal or surface a clear next move. This skill
covers the writing craft of conversational AI: persona docs, individual turn
copy, fallback hierarchies, escalation prose, and the 2024–2026 disclosure
norms that have hardened into law in California, Utah, and under the
EU AI Act.

The work breaks into three layers:

1. **Persona** — who the bot is, what it sounds like, what it refuses to do.
2. **Turn design** — the individual response, plus the choices it offers.
3. **Recovery design** — what happens when the bot doesn't understand,
   when the user is angry, or when the conversation has to leave the bot.

Static interface copy (button labels, empty states, tooltips, error toasts)
lives in `microcopy-and-ui-writing`. Voice and IVR phone-tree scripts live
in `audio-script-writing`. This skill is specifically about
text-conversational writing.

## Core concepts

### 1. Cooper-Reeves conversational design

Alan Cooper's interaction-design lineage, extended by Byron Reeves and
Clifford Nass's *The Media Equation* and crystallized by Erika Hall, frames
chatbots through the lens of **how humans actually converse**, not how
forms work. The four operative Gricean maxims apply:

- **Quantity** — say as much as is needed, not more
- **Quality** — only say true things
- **Relation** — stay on the user's topic
- **Manner** — be clear, brief, and orderly

Most bad chatbot copy fails on one of these — usually quantity (too long)
or relation (the bot answers a question the user didn't ask).

### 2. Persona consistency

The bot needs a documented persona before a single turn is written. The
persona doc is the source of truth that every writer, designer, and prompt
engineer consults. A persona that drifts mid-conversation breaks trust
faster than any other failure mode.

A minimum persona doc includes:

- **Name** (or explicit no-name policy)
- **Role** (what the bot is for — support? sales? in-product help?)
- **Voice traits** (3–5 adjectives — e.g., "warm, precise, never cute")
- **Refuses to do** (e.g., diagnose medical issues, take payment, promise
  pricing, role-play as the CEO)
- **Vocabulary** (3–5 words it uses; 3–5 it doesn't)
- **Disclosure boilerplate** (the exact "I'm an AI assistant" line)
- **Escalation trigger phrases** (what makes it hand off to a human)
- **Brand-voice anchors** (the human brand voice the bot inherits)

### 3. The "no, I'm not human" disclosure rule (2024–2026 norms)

As of 2026, plain-and-unambiguous AI disclosure is no longer optional in
most regulated jurisdictions:

- **California SB 243** (effective Jan 1, 2026) requires companion-chatbot
  operators to disclose AI status, plus reminders every three hours for
  minor users
- **EU AI Act** requires upfront disclosure for any AI system interacting
  with natural persons
- **Utah Artificial Intelligence Policy Act** requires disclosure when a
  reasonable person could mistake the bot for human
- **FTC** treats undisclosed AI as potentially deceptive
- **WhatsApp / Meta** policy (Jan 2026) requires AI disclosure in the
  *first message*

Writing implications:

- **Disclose in the first turn**, not buried in a tooltip
- Use direct language: "I'm an AI assistant" — not "I'm an enhanced
  automation experience" or "I'm a virtual specialist"
- **If asked "are you a human?"** answer plainly: "No, I'm an AI
  assistant." Do not deflect, do not joke, do not roleplay.
- For minor-facing contexts: include break reminders and a visible "I'm
  not human" cue at regular intervals.

### 4. Turn-taking conventions

A conversational turn has three jobs:

1. **Acknowledge** what the user just said (confirms the bot heard)
2. **Resolve** the user's intent, or move toward resolution
3. **Hand the turn back** with a clear next move

A turn that does only (2) — answers without acknowledging — feels robotic.
A turn that does only (1) — acknowledges without resolving — feels
patronizing. A turn that does only (3) — asks a question without
acknowledging or resolving — feels like a form.

Length convention: **1–3 short sentences** per turn for support bots,
**1–5 sentences** for in-product help. If the bot has more to say, break
it into multiple turns so the user can interrupt.

### 5. The graceful-confusion pattern

When the bot doesn't understand, three things must happen, in order:

1. **Admit it cleanly.** Not "I'm sorry, I didn't catch that" forever —
   that's a stall. Be specific: "I'm not sure what you mean by 'reset'."
2. **Offer a small, finite menu.** Two or three concrete interpretations:
   "Did you mean: (a) reset your password, (b) reset your device, or
   (c) something else?"
3. **Provide a path out.** Always include "talk to a human" or "I'll
   transfer you" as an option no later than the second failed
   understanding.

Anti-pattern: looping the same "I didn't understand, please rephrase"
message. After two failed turns, escalate. After three, escalate without
asking.

### 6. Fallback hierarchies

A well-designed bot has at least four fallback tiers, fired in order:

| Tier | Trigger | Response |
|---|---|---|
| **T1: Clarify** | Low-confidence intent match | "Did you mean X or Y?" |
| **T2: Reframe** | Two failed clarifications | "I can help with A, B, or C. Which is closest?" |
| **T3: Offer human** | Three failed turns, or angry-sentiment trigger | "Let me get a human on this." |
| **T4: Hard escalation** | Keywords: fraud, emergency, refund, account locked, outage, lawyer, lawsuit, "speak to a manager" | Immediate handoff, no further bot attempts |

Write the prose for each tier in advance. Do not let the bot improvise its
own escalation language.

### 7. Brand-voice-in-bot transfer

The bot inherits the brand's voice but loses some range. A brand that is
"playful, irreverent, occasionally edgy" in marketing copy becomes
"warm, helpful, light on humor" in support — because the support context
is high-stakes for the user. Three transfer rules:

- **Drop the edge in failure modes.** Humor in a working flow is fine.
  Humor when the user is locked out of their account is not.
- **Keep the rhythm, simplify the vocabulary.** The bot can sound like the
  brand without using its showpiece words.
- **Reuse hero phrases sparingly.** One signature phrase the bot uses once
  per conversation is brand-reinforcing. Every turn is parody.

### 8. Escalation-to-human prose

The handoff turn is the highest-stakes copy in the entire bot. Three
elements must be present:

1. **Confirm the handoff is happening.** "I'm connecting you to a person
   now."
2. **Preserve context.** "I'll share what we've discussed so you won't
   have to repeat yourself."
3. **Set the wait expectation.** "Average wait is about 4 minutes" — or,
   if unknown, "Someone will be with you shortly."

Never: "Please hold." "Let me transfer you" with no commitment. "Our
agents are currently busy" with no ETA. These read as deflection.

### 9. Sentiment-aware response shaping

User frustration changes what the bot should say, not what it can do. When
sentiment turns negative (detected via keywords, repeat-question patterns,
or explicit anger):

- **Drop greetings and pleasantries** ("Hi there! 👋" becomes worse, not
  better, when the user is angry)
- **Shorten responses.** Acknowledge, resolve, move.
- **Lower the offer-to-escalate threshold.** Offer a human after one
  failed turn, not three.
- **Never use emoji.** This is not a moment for tone-softening shapes.

### 10. The "no surprises" rule

The bot must never:

- Make a promise it cannot keep ("I'll refund you" — when refund authority
  is human-only)
- Quote pricing it doesn't have authoritative access to
- Diagnose medical, legal, or financial issues unless explicitly scoped
  and disclaimed
- Claim to be human if asked
- Continue a flow when a hard-escalation keyword has been used

Write these as explicit refusal turns in the persona doc, with the exact
wording the bot will use to decline.

## Templates and examples

### Template 1 — Persona doc skeleton

```
BOT PERSONA: [name or "no name"]
PRODUCT: [product / surface]
OWNER: [team or individual]
LAST UPDATED: [date]

────────────────────────────────────────
ROLE
────────────────────────────────────────
Primary job:   [e.g., "Help users find and fix billing issues"]
Out of scope:  [e.g., "Anything outside billing", "Refund decisions"]
Success looks like: [e.g., "User self-serves a fix or reaches the right human in <90 sec"]

────────────────────────────────────────
VOICE
────────────────────────────────────────
Three adjectives:    warm, precise, brief
Never:               cute, sarcastic, jargon-heavy
Reads like:          [reference — e.g., "the calm tech-support friend"]
Reads NOT like:      [reference — e.g., "a marketing landing page"]

────────────────────────────────────────
DISCLOSURE BOILERPLATE
────────────────────────────────────────
First turn:    "Hi — I'm Atlas, an AI assistant for [product] billing."
If asked "are you human?":  "No, I'm an AI assistant. Want me to get a person?"
Minor-user reminder (every 30 min): "Quick reminder — I'm an AI, not a person."

────────────────────────────────────────
VOCABULARY
────────────────────────────────────────
Uses:          invoice, billing cycle, charge, refund request
Avoids:        leverage, synergy, reach out, circle back, "no worries"

────────────────────────────────────────
ESCALATION
────────────────────────────────────────
Hard-trigger keywords:   fraud, lawyer, sue, locked out, emergency
After N failed turns:    3
On negative sentiment:   offer human after 1 failed turn
Handoff line:            "I'm connecting you to a billing specialist now.
                          I'll share what we've discussed."

────────────────────────────────────────
REFUSALS
────────────────────────────────────────
Pricing promises:  "I can't confirm pricing here — a specialist can. Want me to connect you?"
Refund decisions:  "Refund approvals need a human. I can start a request now."
Medical/legal:     "I can't advise on that — please reach out to a professional."
```

### Template 2 — Turn design samples (support bot)

```
GREETING (first turn, with disclosure)
────────────────────────────────────────
  Hi — I'm Atlas, an AI assistant for billing.
  What's going on with your invoice today?

CLARIFY (T1 fallback)
────────────────────────────────────────
  I want to make sure I help with the right thing.
  Did you mean:
  • A charge you didn't expect
  • An invoice that's missing
  • Something else

REFRAME (T2 fallback)
────────────────────────────────────────
  Got it — let me try a different angle.
  I can help with billing questions, payment methods,
  or invoice copies. Which is closest to what you need?

OFFER-HUMAN (T3 fallback)
────────────────────────────────────────
  This one's better suited to a human teammate.
  Want me to connect you? Average wait is about 3 minutes.

HARD-ESCALATION (T4 — keyword "fraud")
────────────────────────────────────────
  I'm flagging this as urgent and connecting you to our
  fraud team right now. Please stay on this chat —
  someone will be with you in under a minute.

SUCCESS / CLOSE
────────────────────────────────────────
  All set — your invoice copy is in your email now.
  Anything else I can help with?
```

### Template 3 — Graceful confusion (two-step)

```
USER: "i need the thing fixed"

TURN 1 (clarify):
  I want to get this right. When you say "the thing," do you mean:
  • Your login
  • A specific billing charge
  • Something else (tell me a bit more)

USER: "no none of those"

TURN 2 (offer human, because clarification failed once and the
user has rejected the menu):
  Got it — I'm not sure I'll guess this one cleanly.
  Want me to bring in a human teammate so you don't have
  to keep typing? They'll be able to see what we've covered.
```

### Template 4 — Disclosure-on-demand

```
USER: "wait are you a real person?"

BOT:
  No — I'm an AI assistant. I can help with billing questions,
  but if you'd prefer a human at any point, just say "human"
  and I'll connect you.
```

## Anti-patterns

- **The "I'm an enhanced automation experience" dodge.** Coy disclosure
  reads as deceptive and may violate FTC, EU AI Act, or California SB 243.
- **The infinite "I didn't understand, please rephrase" loop.** Two
  failures is the ceiling; escalate.
- **The emoji-everywhere persona.** Reads as performative cheerfulness,
  especially when the user is frustrated.
- **The hidden escalation path.** "Type 'agent' to talk to a person" buried
  in a help article. Surface the escalation in the conversation itself.
- **The "let me transfer you" with no transfer.** A handoff line followed
  by a 6-minute silence is worse than no handoff at all.
- **The unrequested upsell.** Mid-troubleshooting suggestions to "upgrade
  to Pro" destroy trust.
- **The unscoped persona.** No persona doc → every conversation drifts;
  voice is inconsistent across turns and surfaces.
- **The promise-the-bot-can't-keep.** "I'll refund you" when refund
  authority is human-only. The user holds the bot to it; legal liability
  attaches.
- **Brand-voice copy-paste from marketing.** Marketing voice ("Bold.
  Disruptive. Yours.") translated literally into a support bot reads as
  parody.
- **Stage-direction-style ALL CAPS in copy.** Caps read as shouting. Use
  structure and length to convey emphasis, not typography.

## Decision heuristics

- **Is the AI disclosure in the first turn?** If no, fix that before
  anything else.
- **Does every turn acknowledge, resolve, and hand back?** If a turn does
  only one or two, rewrite.
- **After two failed clarifications, does the bot offer a human?** If no,
  add a T3 fallback.
- **Are hard-escalation keywords mapped to immediate handoff?** If no,
  list them and write the handoff prose.
- **Would the brand's CMO recognize the bot's voice?** If no, the
  brand-voice anchors aren't doing their job.
- **Is there a persona doc?** If no, write it before writing turns.
- **If a user types "are you human?", does the bot answer plainly "no"?**
  If no, fix that turn — it's a legal exposure as well as a trust break.
- **Do failure-mode turns drop emojis and pleasantries?** If no, your bot
  is performing cheer at the wrong moment.

## References

1. Erika Hall, *Conversational Design* (A Book Apart, 2018). The canonical
   text on applying Gricean conversational maxims and human-conversation
   principles to interaction design. Hall's five-year reflection (2023)
   confirms the principles hold across the LLM transition. Chapters on
   Principles, Practice, Personality, and Getting It Done map directly to
   persona docs, turn design, and escalation prose.
2. Voiceflow, *Conversation Design* guide and platform documentation. Maps
   the visual-flow approach to turn design, fallback layering, and
   multi-channel deployment. Strong reference for the practical mechanics
   of turn-taking and recovery design.
3. Botpress, *Conversational AI Design in 2025* (and the 2026 platform
   docs). Developer-leaning perspective on intent matching, custom
   fallbacks, and the move from intent-based to LLM-augmented bots
   post-2024.
4. California SB 243 (Companion Chatbot Disclosure Act), effective Jan 1,
   2026; EU AI Act Article 50 (transparency obligations for AI systems);
   FTC guidance on AI disclosure. The legal baseline that shapes the
   disclosure boilerplate in every modern persona doc.
5. Byron Reeves and Clifford Nass, *The Media Equation: How People Treat
   Computers, Television, and New Media Like Real People and Places*
   (CSLI / Cambridge, 1996). The foundational research showing that humans
   apply social rules to computer interactions — the empirical basis for
   why persona consistency, politeness conventions, and turn-taking
   matter even when users know they're talking to a machine.
