<!-- hub-reference-banner -->
> **Reference file — part of the `content-and-marketing-writing` hub.** Formerly the standalone `audio-script-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: audio-script-writing
version: "1.0.0"
updated: "2026-05-29"
category: custom
tags:
  - writing
  - audio-writing
  - podcast
  - voice-ui
  - ivr
  - script
description: >
  Writing for the ear, not the eye — podcasts, voice-overs, voice-UI prompts (Alexa,
  Google Assistant, Siri shortcuts), IVR scripts, audiobook narration, and broadcast
  copy. Covers the cold-open pattern, ad-break placement, narration cadence, signposting,
  the "one idea per sentence" rule, parenthetical-elimination, pacing markers (commas,
  ellipses, line breaks as breath beats), the sample-dialog deliverable for voice
  conversation design, IVR pause budgets (100/250/500 ms), transcription-friendly
  writing, and podcast show-notes craft with chapter timestamps. References NPR script
  guidance, Google Assistant Conversation Design, Amazon Alexa Design Guide, the Mozilla
  i18n notes on voice strings, and IVR pacing standards.
  TRIGGER: "podcast script", "voice-over script", "VO script", "voice prompt",
  "Alexa skill", "Google Action", "Siri shortcut", "IVR script", "auto-attendant
  script", "writing for the ear", "narration", "audiobook", "cold open", "ad break",
  "show notes", "podcast chapters", "voice-UI", "conversation design", "audio script",
  "spoken-word script", "broadcast copy".
  SKIP: presentation delivery, vocal warmups, or speaking technique (use
  public-speaking-and-presentations); interview question design or moderating live
  conversation (use interview-and-conversational); plain-English customer prose with
  no audio target (use plain-language); accessibility captions or alt text for static
  images (use visual-writing).
triggers:
  - "podcast script"
  - "voice-over script"
  - "VO script"
  - "voice prompt"
  - "Alexa skill"
  - "Google Action"
  - "Siri shortcut"
  - "IVR script"
  - "auto-attendant"
  - "writing for the ear"
  - "narration"
  - "audiobook"
  - "cold open"
  - "ad break"
  - "show notes"
  - "podcast chapters"
  - "voice-UI"
  - "conversation design"
  - "audio script"
  - "broadcast copy"
skip:
  - presentation delivery or speaking technique → use public-speaking-and-presentations
  - interview moderation or live dialogue → use interview-and-conversational
  - plain-English customer prose without an audio target → use plain-language
  - accessibility captions or chart alt text → use visual-writing
related:
  - public-speaking-and-presentations
  - interview-and-conversational
  - plain-language
  - writing-expert
  - storytelling-and-narrative
---

# Audio Script Writing

Reference for writing that will be heard, not read. Covers podcasts, voice-overs,
voice-UI prompts, IVR scripts, audiobook narration, and the broadcast copy idiom.
The eye can re-scan a sentence; the ear gets one pass. Every rule in this skill
serves the listener's single-pass comprehension.

Deliver all responses in a direct, plain register. Avoid hedging and meta-commentary.

---

## When to use this skill

Activate when the user:

- Has prose that will be spoken aloud (podcast, voice-over, narration, ad read)
- Is designing a voice-UI prompt for Alexa, Google Assistant, or Siri Shortcuts
- Is writing or auditing an IVR script or auto-attendant flow
- Asks how to write a cold open, structure ad breaks, or pace a narration
- Needs podcast show notes with chapters and timestamps
- Is converting a written article into a spoken-word piece
- Wants the "writing for the ear" treatment applied to existing text
- Asks for a sample-dialog deliverable for a conversation-design project

Skip when:

- The task is **delivery** — pacing your own speech, breath control, gesture (use
  `public-speaking-and-presentations`)
- The task is **live conversation** — interview questions, moderation, dialogue
  design (use `interview-and-conversational`)
- The task is plain customer prose for the **eye**, not the ear (use `plain-language`)
- The task is image alt text or chart captions (use `visual-writing`)

---

## The one rule: writing for the ear is single-pass

A reader can re-scan; a listener cannot. Every other rule in this skill follows
from that constraint.

Listeners process roughly **150–160 words per minute** for narration, **130–145
wpm** for voice-UI prompts, and **110–125 wpm** for IVR. They lose comprehension
on:

1. Subordinate clauses stacked more than one deep
2. Parentheticals (the ear has no parens)
3. Anaphora across more than ~15 words ("it", "this", "that" referring back too far)
4. Numbers and proper nouns delivered without a beat of silence after
5. Homophones in ambiguous contexts (their/there, principal/principle)
6. Lists longer than three items without explicit numbering

---

## Core concept 1 — The one-idea-per-sentence rule

Eye prose can carry two or three coordinated ideas in a sentence. Ear prose
cannot.

**Eye version (fine on the page):**

> The migration, which we'd been planning since February, finally launched on
> Tuesday after a final round of testing that revealed two minor bugs we patched
> overnight.

**Ear version:**

> We'd been planning the migration since February. It finally launched on
> Tuesday. The last round of testing turned up two small bugs. We patched them
> overnight.

The eye version has 31 words and three subordinate clauses. The ear version
has four sentences, each carrying one fact. The total word count went up by
two words. The comprehension went up by an order of magnitude.

---

## Core concept 2 — Signposting

The ear cannot scroll back. The writer must signal structure out loud.

**Signpost vocabulary:**

- **First / second / third / finally** — explicit numbering
- **Here's the thing** — flag a key insight
- **Two reasons** — set up a numbered list
- **Quick recap** — a deliberate echo of earlier material
- **Stay with me** — flag complexity ahead
- **Coming up** — preview before an ad break
- **Back to the story** — return marker after a digression

A 20-minute podcast with no signposts forces the listener to track structure
themselves. A 20-minute podcast with signposts every 90–120 seconds carries
itself.

---

## Core concept 3 — No parentheticals, no em-dashed asides

The voice cannot draw parens. Two patterns to refactor:

**Refactor parentheticals into separate sentences:**

- Before: "The team (which had only formed in January) shipped on time."
- After: "The team shipped on time. They'd only formed in January."

**Refactor em-dash asides into a hard sentence break or remove:**

- Before: "Our latency — which had been climbing for three weeks — finally stabilized."
- After: "Our latency had been climbing for three weeks. It finally stabilized."

Em-dashes survive in print because the eye treats them as visual punctuation. In
audio, they become awkward pauses with no semantic payoff.

---

## Core concept 4 — Pacing markers as breath beats

Punctuation in a script is not grammar; it is breath direction.

| Mark | Pause length | Function |
|------|--------------|----------|
| Comma | ~100 ms | mid-sentence beat |
| Period | ~250 ms | sentence boundary |
| Em-dash or ellipsis | ~350–500 ms | dramatic pause |
| Paragraph break | ~750 ms–1.5 s | scene change |
| `[pause]` stage direction | author-specified | deliberate silence |
| `[beat]` | ~1 s | reset between ideas |

IVR engineers use the same budget formally: **100 ms after commas, 250 ms
after sentence-internal phrases, 500 ms at sentence end**, sometimes inserted as
silent WAV padding. Treat your script the same way — the punctuation is the
direction.

---

## Core concept 5 — The cold open

The first 30 seconds of a podcast are the cold open. The listener has not
committed yet. Three patterns that work:

**Pattern A — Drop into a scene:**

> It's 2 a.m. The paging system goes off. Marcus rolls over, reads the alert,
> and the alert is wrong.

**Pattern B — A question with stakes:**

> What would you do if your entire backup tier disappeared in the middle of a
> restore? That's what we're talking about today.

**Pattern C — A single startling fact:**

> Last year, 47% of incident retrospectives never produced a single action item.
> Today: why, and what to do about it.

**Anti-pattern — The throat-clear:**

> Hi everyone, welcome to the podcast. Today we're going to be talking about,
> uh, incident response.

The throat-clear loses 20–30% of listeners in the first 60 seconds. Cold-open
first, intro second.

---

## Core concept 6 — Ad break placement

Two principles:

1. **Place breaks at narrative seams.** Between segments, not mid-thought. A
   listener tolerates "we'll be right back" at a scene boundary; they resent it
   mid-sentence.
2. **Promise specificity across the break.** "Coming up: the one Postgres
   setting that broke production at 3 a.m." Vague promises ("more after this")
   permit a tab switch from which the listener does not return.

**Standard placements for a 30-minute episode:**

- One pre-roll (before cold open or just after, 0:00–1:00)
- One mid-roll at the act break (~12:00–15:00)
- One end-roll (~28:00) — lowest revenue tier; many shows skip

---

## Core concept 7 — Voice-UI prompts (Alexa, Google Assistant, Siri)

Conversation design adds three constraints on top of "writing for the ear":

1. **The prompt must end with an explicit prompt-for-input.** "Which would you
   like — small, medium, or large?" Not: "Let me know."
2. **Confirm without echoing.** "Adding milk to your list" is good; "I heard
   you say milk; I will now add milk to your list" is voice-UI throat-clearing.
3. **Three-option ceiling.** Lists of more than three options in a single
   prompt exceed working memory. Chunk into sub-menus or use category-first
   ("Food, household, or other?").

**Google Conversation Design** requires two artifacts as deliverables:

- **Sample dialogs** — scripts showing the happy path and the top 3–5 error
  paths
- **Flow diagram** — branches, including no-input, no-match, and help intents

**Alexa situational design** principles (from the Alexa Design Guide):

- Adapt to the speaker, the context, and the device
- Build incrementally; ship the optimal happy-path script first
- Test with real users; pseudo-fluency in prototypes hides real failures

---

## Core concept 8 — IVR scripts

IVR is voice-UI with a phone-grade audio channel and a stressed caller. Extra
rules:

1. **Most-likely path first.** "Press 1 for billing, 2 for support, 3 for
   everything else." Order by call-volume share, not alphabetical.
2. **Always offer a human escape.** "Press 0 at any time to reach an agent."
   Hiding this is the single highest source of caller anger.
3. **Confirm critical inputs.** Account numbers, phone numbers — read back at
   the cadence a person can write down. Group digits in 3-3-4 or 4-4 chunks.
4. **Never use jargon.** "Press 1 for IVR self-service" is broken on its face.
5. **Re-prompt twice, then escalate.** Caller silence twice in a row = route
   to agent. Do not loop endlessly.

---

## Core concept 9 — Transcription-friendly writing

Podcasts and broadcasts often get auto-transcribed for show notes, captions,
and SEO. Three patterns that fail transcription:

- **Numeric figures spoken as digits without spelling-out** ("twenty-twenty-six"
  vs "two-zero-two-six") — Whisper and similar engines mis-segment without context.
- **Specialist proper nouns delivered fast** — "WiredTiger", "ICU
  MessageFormat" — slow down, optionally spell once.
- **Crosstalk and parallel laughter** — transcription engines drop these entirely.
  Scripted shows can mark `[laughter]` or `[crosstalk]` for the editor.

Best practice: produce the script and the show notes from the same source so
the show notes pre-encode the right spelling and proper nouns.

---

## Core concept 10 — Podcast show notes craft

Show notes serve three masters: the listener, the search engine, and the
sponsor.

**Required sections:**

1. **Episode title** — keep under 60 characters; concrete, specific, no clickbait
2. **One-paragraph summary** — 50–80 words; uses the episode keyword once
3. **Chapter timestamps** — 4–8 chapters; `[12:15] How to find your first sponsor`
4. **Guest bio** — 1 paragraph, links to guest sites
5. **Mentioned resources** — every link spoken aloud, plus internal cross-links
6. **Transcript** — full or summary; search engines index this heavily
7. **Call to action** — newsletter, follow, review

**Anti-patterns:**

- Generic chapter labels ("Introduction", "Main discussion", "Conclusion")
- Single 600-word blob with no headers
- Missing CTA at the end

---

## Templates

### Podcast cold-open skeleton

```text
[COLD OPEN — 20-40 seconds]
[Drop into scene OR question with stakes OR startling fact]
[Single beat of silence]

[THEME / STING — 5-10 seconds]

[HOST INTRO — 30-45 seconds]
You're listening to [SHOW NAME]. I'm [HOST]. Today: [one-sentence promise of episode].
[Optional: name-of-guest credit]

[ACT 1 — main content]
...

[BREAK 1 — at act seam, around 12-15 minutes]
"Coming up: [specific tease]. Back in 30."
[Ad slot — 60-90 seconds]

[ACT 2 — main content continues]
...

[OUTRO — 30-60 seconds]
"That's [SHOW NAME] for this week. If you got something from this episode,
[concrete CTA]. Next time: [tease]. I'm [HOST]. Thanks for listening."
```

### Voice-over (commercial) script — 30 seconds

```text
[VO, warm]: It's six a.m. The coffee's already on. (beat)

[VO, conversational]: That's because Maker Coffee learns your routine. One
button, every morning, exactly how you like it. (beat)

[VO, slight smile]: Make tomorrow easier. (pause) Maker Coffee — at makercoffee.com.
```

Notes for the VO talent in brackets. Pauses marked. Total word count: ~40
words; reads in ~25 seconds with the marked pauses.

### Voice-UI prompt skeleton (Alexa / Google)

```text
INTENT: ADD_TO_LIST

[Opening prompt]
"What would you like to add?"

[Capture happy path]
"Adding {item} to your shopping list. Anything else?"

[No-input retry 1]
"What should I add?"

[No-input retry 2]
"I didn't catch that. You can say 'add milk' or 'never mind'."

[No-match]
"Sorry, I missed that. What do you want to add?"

[Help]
"You can add items to your list one at a time. Just say the item — like
'add bananas' — and I'll add it. Say 'done' when you're finished."

[Exit]
"Okay, all set."
```

### IVR auto-attendant — 6-option menu

```text
[Greeting — 8 seconds]
"Thanks for calling Acme Support. This call may be recorded for quality."
[500 ms pause]

[Main menu — keep under 25 seconds total]
"To make this fast, press the number for your need."
[250 ms]
"For billing or your account — press 1."
[100 ms]
"For technical support — press 2."
[100 ms]
"To check on an existing ticket — press 3."
[100 ms]
"For all other questions — press 4."
[100 ms]
"To reach an agent at any time, press zero."
[500 ms]

[On no input, repeat once. On second no-input, route to agent queue.]
```

### Podcast show notes — full template

```markdown
# Episode 47: Why incident retros stall (and how to fix it)

**Listen:** [Apple] | [Spotify] | [Overcast] | [RSS]
**Length:** 38 minutes
**Guest:** Marcus Chen, SRE Lead at Northwind

---

Most incident retrospectives never produce a single action item. In this
episode, Marcus Chen walks us through the three failure modes he's seen across
60+ post-incident reviews — and the one structural change that fixed all three.

---

## Chapters

- [00:00] Cold open: the 3 a.m. page that wasn't
- [01:45] Why retros stall: the three failure modes
- [09:20] The "action-item owner" rule
- [14:00] Mid-roll
- [15:30] Designing a retro template that survives Q4
- [25:00] The "30-day check" — does anything actually change?
- [33:00] Listener question: when to skip the retro entirely
- [36:30] What Marcus is reading

## Resources mentioned

- Marcus's retro template (GitHub)
- Etsy "Blameless PostMortems" essay
- Google SRE workbook chapter 15
- Northwind incident-response playbook

## About Marcus

Marcus Chen leads the SRE team at Northwind. Before that, he ran on-call
operations at two infrastructure startups. Find him at marcuschen.dev.

## Follow / subscribe

- New episodes every Tuesday
- Newsletter at example.com/newsletter
- Leave a review on Apple Podcasts — it helps more than you think
```

---

## Anti-patterns

| Anti-pattern | Why it fails | Fix |
|--------------|--------------|-----|
| Long sentences with nested clauses | Ear loses the through-line | One idea per sentence |
| Em-dashes and parens in spoken prose | No audio equivalent | Hard sentence break |
| Reading numbers fast | Listener can't write them down | Chunk + slow + read-back |
| Throat-clear intros ("Hi everyone, welcome…") | Loses 20–30% in 60s | Cold open first |
| Ad breaks mid-thought | Listeners resent it | Place at act seams |
| Voice-UI prompt without input request | Caller doesn't know it's their turn | End with explicit prompt |
| IVR with no zero-out | Caller rage | Always offer human escape |
| Generic show-note chapter labels | No SEO, no scannability | Concrete labels with keywords |
| Reading the URL out loud verbatim | "h-t-t-p-colon-slash-slash" is torture | Say the brand; put URL in notes |
| Crosstalk in scripted shows | Transcription breaks | Stage-direct or re-record |
| Words that look spelled but sound different | "principal" vs "principle" | Pre-read aloud; swap homophones |
| Lists of 4+ options in a voice prompt | Exceeds working memory | Chunk to 3 or use sub-menus |

---

## Decision heuristics

**Audio-target detection — when to apply this skill:**

- Will a human voice (or TTS) speak the final output? → Yes, this skill
- Is it captions for someone else's audio? → No, use `visual-writing`
- Is it interview prep where you're moderating live? → No, use `interview-and-conversational`
- Is it customer-facing prose that happens to be readable aloud? → Probably `plain-language`

**Cold open vs. throat-clear:**

- Will the first sentence make a listener stop scrolling? → Cold open
- Does the first sentence introduce you or the show? → Throat-clear; rewrite

**Punctuation pass:**

- Read every sentence aloud. If you ran out of breath, break the sentence.
- Mark every parenthetical for rewrite.
- Mark every em-dash; keep only those that earn the dramatic pause.

**Voice-UI vs IVR vs narration:**

- Constrained turn-taking, single-task → voice-UI
- Phone, distressed user, menus → IVR
- Continuous speech, narrative arc → podcast / VO / audiobook

**Ad placement:**

- Is the listener mid-thought? → Move the break
- Is the next segment teased specifically? → Keep
- Is the tease vague ("more after this")? → Rewrite the tease

---

## References

- W3C / IAB show-note conventions; NPR Training: ["How to write a mean
  script"](https://www.npr.org/2023/02/17/1156597697/how-to-write-script-speech-podcast-voice)
- Adonde Media: ["The Podcast Script: Writing for the
  Ear"](https://medium.com/@AdondeMedia/the-podcast-script-writing-for-the-ear-95fab0d9c5be)
- Google: [Conversation Design](https://developers.google.com/assistant/conversation-design/welcome)
- Amazon: Alexa Design Guide (situational design principles)
- Call Center Studio: [Best Practices for IVR
  Scripts](https://callcenterstudio.com/blog/best-practices-for-writing-an-effective-and-professional-ivr-script/)
- Bobby Owsinski: *The Recording Engineer's Handbook* (5th ed.) — production
  notes on narration capture, mic technique that informs script length
- Buzzsprout: ["How to Write Podcast Show
  Notes"](https://www.buzzsprout.com/blog/podcast-show-notes)
- Exemplary AI: ["Podcast Show Notes 101: SEO
  guide"](https://exemplary.ai/blog/podcast-show-notes-boost-seo)

---

## Related skills

- `public-speaking-and-presentations` — for delivery technique, not writing
- `interview-and-conversational` — for live-dialogue design
- `plain-language` — for customer-facing prose
- `writing-expert` — for general prose craft
- `storytelling-and-narrative` — for narrative arc inside an audio piece
