<!-- hub-reference-banner -->
> **Reference file — part of the `executive-comms` hub.** Formerly the standalone `public-speaking-and-presentations` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: public-speaking-and-presentations
version: "1.2.0"
updated: "2026-05-29"
description: >
  Presentation and public-speaking craft — slide-doc vs PowerPoint structure,
  Kawasaki 10/20/30 rule, Duarte Resonate framework, talk arc design, Q&A
  handling, rehearsal discipline, demo prep, and virtual vs in-person calibration.
  TRIGGER: "prep a talk", "slide deck", "presentation structure", "10/20/30 rule",
  "Duarte Resonate", "headline deck", "Q&A handling", "rehearsal plan", "demo prep",
  "how to open/close a talk".
  SKIP: writing a written document (use writing-expert); executive communication
  strategy — what to say to whom (use executive-comms); creating slide *content*
  about a specific technical topic (use domain skill); interview prep
  (use interview-and-conversational); negotiation (use negotiation-and-persuasion).
triggers:
  - prep a talk
  - deck for
  - slide deck
  - presentation
  - 10/20/30 rule
  - Duarte
  - headline deck
  - talk structure
  - Q&A handling
  - rehearsal plan
  - demo prep
  - how to open a talk
  - how to close a talk
  - presenting to
skip:
  - writing a written document only (use writing-expert + executive-comms)
  - creating slide content about a specific topic — what to put on the slide, not how to present (use domain skill)
  - executive communication strategy — what to say to whom, not how to present (use executive-comms)
  - interview prep (use interview-and-conversational)
  - negotiation prep (use negotiation-and-persuasion)
related:
  - writing-expert
  - executive-comms
  - storytelling-and-narrative
  - interview-and-conversational
---

# Public Speaking and Presentations

## How to Use This Skill

When a user asks for presentation or public-speaking help, apply the guidance below. Default response shape:

- **For deck-building requests** ("help me build a deck," "outline my talk"): produce a structured plan — audience analysis answers, one-sentence core message, talk arc with time allocations, headline-deck titles for each major section, and a recommended opening pattern.
- **For critique requests** ("review my slides," "is my structure right"): produce a gap analysis against the frameworks in this skill, ordered by impact, with specific fixes.
- **For coaching requests** ("how do I handle Q&A," "what's the Duarte framework"): produce direct, specific guidance from the relevant section; do not pad with unasked-for coverage of other sections.

**Clarifying questions:** if the user's request is ambiguous on slot length, audience, or purpose, ask exactly one targeted question before proceeding. Examples of useful clarifying questions: "How long is your slot?" / "Who is the primary audience — technical practitioners or executives?" / "Is this a pitch, a keynote, or an internal update?" Do not ask multiple questions at once.

**Self-check before responding:** confirm your output covers the user's actual question, names the specific principle or framework you are drawing on, and gives at least one concrete action (not just a description of a concept).

## Audience Analysis Before You Open PowerPoint

Every structural and design decision depends on your answers to four questions. Answer them before touching a slide tool.

1. **Who is in the room?** Job function, seniority, and domain expertise determine vocabulary, assumed context, and decision-making authority. A room of engineers and a room of VPs require different entry points for the same topic.
2. **What do they already believe?** Start from their current position, not yours. If they are skeptical, open with evidence. If they are already bought-in, skip the justification and go straight to the "how."
3. **What one action do you want them to take?** A talk with no specific ask produces no specific result. Name the action: approve budget, change a process, try a tool, share a finding with their team.
4. **What is the one thing they must remember?** If they forget everything else in 48 hours, what sentence should survive? Write that sentence before slide 1. Every slide either supports it or gets cut.

If you cannot answer all four, stop and gather the information. Building slides before answering them wastes everyone's time.

## Slide-Doc vs Presentation Slides

These are two completely different artifacts with different jobs. Confusing them is the single most common presentation mistake.

**Presentation slides** (Reynolds / "Presentation Zen" model) exist to support a live speaker. The speaker carries the argument; the slide carries visual reinforcement. One idea per slide. Minimal text. Large images or diagrams. Without the speaker the deck is nearly meaningless — that is a feature, not a bug.

**Slide-docs** (Tufte-influenced, common in consulting) are standalone documents formatted as slide-shaped pages. They contain complete sentences, dense supporting data, footnotes, and annotations. They can be emailed, read without a presenter, and used as a leave-behind. McKinsey and BCG reports are slide-docs.

When to use which:
- You are presenting live to an audience who will watch you → presentation slides
- You are sending a document for async reading or executive review → slide-doc
- You must do both → build the slide-doc first, then strip it into a presentation deck; never conflate the two into a "slidument" (Reynolds' term for the worst of both)

The core failure mode is the slidument: bullet-heavy slides that are too dense to look good on screen and too sparse to be useful to a reader. Avoid it by deciding the audience's mode of consumption before opening PowerPoint.

## Kawasaki 10/20/30 Rule

Guy Kawasaki's rule for investor/pitch contexts (blog post 2005; "The Art of the Start 2.0", 2015):

- **10 slides** — the human mind cannot absorb more than ten concepts in a single sitting in a pitch context
- **20 minutes** — even if you have a full hour, present in 20 minutes; leave the rest for discussion
- **30-point minimum font** — if you cannot fit your message in 30pt type, you have too many words

The 10 canonical slides for a pitch: problem, solution, business model, underlying magic/technology, go-to-market, competitive analysis, team, projections/milestones, status/timeline, call to action.

This rule is pitch-specific. For technical talks and keynotes, adapt the spirit: more slides are fine as long as each slide holds one idea, total runtime fits your slot with 10% margin, and no slide requires the audience to squint.

**Slides-per-minute rule of thumb for non-pitch contexts:** plan 1–2 minutes per slide. A 30-minute conference talk supports 15–30 slides; a 60-minute workshop supports 30–45. If you are running faster than 1 minute per slide consistently, your slides are too granular. If you are spending more than 3 minutes on a single slide, split the idea or move supporting detail to a backup slide.

## Nancy Duarte's Resonate Framework

From "Resonate" (Wiley, 2010). The central insight: every great speech moves back and forth between "what is" and "what could be."

**The contrast pattern:**
1. Establish the current reality (what is) — ground the audience in shared facts
2. Describe the future possibility (what could be) — create tension and desire
3. Repeat this oscillation; each cycle raises the stakes
4. End with the "new bliss" — the world after the audience acts

The pattern mirrors Joseph Campbell's hero's journey. The speaker is not the hero; the audience is the hero. The speaker is the mentor (Yoda, not Luke). Reframing your role this way changes how you write every slide and every sentence.

**The STAR moment** (Something They'll Always Remember): one designed moment of high contrast, surprise, or visceral demonstration. Steve Jobs pulling the MacBook Air from an envelope. Design this moment deliberately; it is not accidental.

**The sparkline** (Duarte's visualization): a graph that plots "what is" (low) vs "what could be" (high) across the arc of a talk. Sketch your own sparkline before writing slides — it forces you to plan the emotional pacing, not just the content sequence.

## Three-Act Talk Structure

Adapted from dramatic structure (Aristotle, Syd Field) for technical and business talks:

**Act 1 — Setup (10–15% of total time):**
- Establish stakes: why should this audience care, right now?
- Define the problem or question the talk will answer
- Signal what the audience will be able to do differently by the end

**Act 2 — Exploration (70–75% of total time):**
- Develop the argument, evidence, or narrative in chunks of roughly equal weight
- Each chunk ends with a mini-resolution before the next conflict or question opens
- Vary the energy within Act 2: primacy-recency research shows audiences are most alert at the start and end of a talk, with attention dipping in the middle — so front-load your most complex material into the early part of Act 2, and place a re-energizing moment (story, demo, striking visual) two-thirds of the way through to lift attention before the close

**Act 3 — Resolution (10–15% of total time):**
- Do not introduce new material here
- Synthesize, do not summarize (synthesis is "therefore"; summary is "and then")
- End on a single, memorable sentence or image the audience can carry out of the room

A common failure: Act 2 runs so long that Act 3 is cut. Budget Act 3 explicitly and protect it.

## The Headline-Deck Approach (McKinsey / BCG)

Every slide title is a complete, standalone sentence stating the takeaway, not a topic label.

Bad title: "Market Analysis"
Good title: "The North American market is growing 12% YoY and outpacing EMEA in every segment."

Rules:
- Subject + verb + object (or subject + verb + modifier)
- Active voice whenever possible
- The title should be true even if the slide body is removed
- A reader who reads only slide titles should follow the full argument
- Body content (charts, tables, bullets) is evidence that supports the title sentence

This approach forces clarity before design. If you cannot write the title sentence, you do not yet know what the slide is arguing. Fix the thinking, then build the slide.

## Visual Hierarchy on Slides

From Reynolds' "Presentation Zen" and Duarte's "slide:ology":

- **One idea per slide.** If you feel the urge to use "and" in a slide title, split the slide.
- **Picture or diagram primary, words secondary.** Images process 60,000x faster than text (commonly cited; cognitive load research). Use full-bleed photography or a single large diagram when the concept benefits from it.
- **White space is not wasted space.** Clutter signals unclear thinking. Removing elements almost always improves a slide.
- **Contrast drives attention.** Size, color, weight, and position all create hierarchy. Make sure your most important element is visually dominant.
- **Consistent template.** Font pairing (one serif + one sans, or two weights of one sans), two or three accent colors, and a grid. Deviate only to signal something exceptional.
- Avoid clip art, default SmartArt shapes, and decorative drop shadows — they consume visual attention without adding information.

## Speaker Notes Discipline

Speaker notes are a rehearsal tool, not a teleprompter script.

Write down:
- The transition sentence that moves from the previous slide to this one
- Any statistic or proper noun you might blank on under pressure
- The one-sentence takeaway for this slide (matches the title)
- A reminder of any physical action: "pause here," "click to next build," "hold for reaction"

Do not write down:
- Full paragraphs of prose you intend to read aloud
- Everything on the slide restated in sentence form
- Filler you would never say in a real conversation

If you need the notes to get through a talk, you have not rehearsed enough.

## The 6x6 Rule — When It Helps and When It Misleads

The 6x6 rule (no more than 6 bullet points per slide, no more than 6 words per bullet) is a useful forcing function for people who default to walls of text. It prevents the worst sliduments.

It is misguided when applied universally:
- A single high-quality sentence is better than six 6-word fragments that lose meaning through compression
- Technical talks often need one precise statement, not six stubs
- Duarte and Reynolds both prefer eliminating bullets entirely in favor of images + brief labels

Use 6x6 as a ceiling-check, not a target. The real target is one idea per slide, presented at whatever word count makes that idea clear.

## Talk Openings — Five Patterns

Anderson ("TED Talks", 2016) and Duarte ("Resonate", 2010) both emphasize the first 60 seconds as disproportionately important — attention is highest and trust is being built.

1. **Startling fact or statistic:** drop a number the audience does not expect; pause to let it land; then explain it
2. **Story:** open with a specific scene — one person, one moment, sensory detail; no "let me tell you a story about..."
3. **Contrarian claim:** state something the audience believes that you are about to challenge; creates immediate tension
4. **Direct question:** ask a question the audience cannot immediately answer; creates an information gap they want filled
5. **Scene-setting:** describe the world as it will look after your idea is adopted; then pivot back to how we get there

Do not open with: your name and title (the host introduced you), an apology, "great to be here," agenda slides, or a joke that is not genuinely funny.

## Talk Closings — Five Patterns

The last 90 seconds are the second-highest-attention moment. End deliberately.

1. **Callback to opening:** return to the story, image, or question you opened with; now it has new meaning
2. **Call to action:** one specific, concrete thing the audience can do in the next 48 hours
3. **Single image or quote:** full-bleed slide with no text except one sentence or one attributed quote; let it sit
4. **Summary triple:** three short phrases that encapsulate the talk's three main points; the rule of three aids memory
5. **The future vision:** describe the world three years from now if the audience acts; end on the "what could be" peak

Do not close with "thank you" as your last slide. Do not close by reading a summary of everything you said. Do not trail off into "so... yeah... that's it."

## Q&A Handling

**The bridging technique** (political communications, widely taught): when a hostile or off-topic question arrives, acknowledge it, then bridge to your main message. "That's a fair question about X. What I can tell you is that the core issue is Y, and here's why..."

**The parking lot:** for questions that are important but outside scope, say explicitly: "That deserves a real answer and I don't want to shortchange it. Let me take that offline." Write it down visibly. Follow up.

**"I don't know" as a complete answer:** say it cleanly, without hedging or improvising a guess. Add "I'll find out and get back to you by [date]." This builds more credibility than a wrong answer delivered confidently.

**Repeat or rephrase before answering:** buys you two seconds to think and ensures the full audience heard the question.

**Controlling the clock:** designate a time-check person or set a visible timer. "We have time for two more questions." Then keep that commitment.

**Hostile questions:** do not match the energy. Lower your voice slightly. Answer the question that is beneath the question. Never argue.

## Rehearsal Discipline

Minimum effective rehearsal for a 30-minute talk:
1. **Three times out loud, alone** — not in your head, not mumbling; full volume, full sentences, full transitions
2. **Once to a real human** — a colleague, partner, or anyone willing to watch; their confusion is signal
3. **Video review of yourself** — watch on mute first (check body language and gesture), then with sound (check pace and filler words)

Common filler words to eliminate: "um," "uh," "you know," "like," "so basically," "right?" Replace with silence. Silence is not awkward; it is authoritative.

Pacing target: 125–150 words per minute for comprehension. Technical content: 110–125 wpm. Record and measure if unsure.

Time your talk at 90% of the slot. If the slot is 30 minutes, rehearse to 27. You will run long under pressure.

## Demo Discipline

**Pre-record a fallback.** For any live demo, have a screen recording ready that you can switch to if connectivity, credentials, or the demo gods fail. Announce it matter-of-factly: "I have a recording as a fallback — let's try live first."

**Narrate while clicking.** Never go silent during a demo. Explain what you are about to do, what you are doing, and what the audience should be watching for. Silence during a demo feels like failure even when nothing has gone wrong.

**Demo-gods mitigation:**
- Test the exact network and machine you will present on, the day before
- Disable notifications (macOS: Do Not Disturb; Windows: Focus Assist)
- Close all non-demo apps and browser tabs
- Use a dedicated browser profile with only the required extensions
- Have a backup laptop or hotspot

**Keep demos short.** A demo that illustrates one idea in 90 seconds is more effective than a walkthrough that illustrates ten in six minutes.

## Virtual vs In-Person Calibration

**Gaze targets:** in-person, scan the room in thirds (left, center, right), making 3–5 second eye contact with individuals. Virtual: look at the camera lens, not at faces on screen. Place the camera at eye level; looking up into a laptop camera is unflattering and signals submission.

**Gesture amplitude:** in-person, gestures can be wide and expressive. Virtual, keep gestures in the frame — roughly shoulder-width. Gestures that go off-camera are invisible and create unease.

**Audio quality:** virtual, audio quality matters more than video quality. A USB condenser microphone, a quiet room, and a soft surface behind you (bookcase, curtain) are worth more than a 4K webcam in an echo chamber. Test your audio before every important virtual session.

**Energy calibration:** virtual audiences are more fatigued and distracted. Increase vocal variety (pace, volume, pitch) noticeably above what feels natural in a room — the camera and compression flatten affect, so what reads as energetic in person reads as flat on screen. Shorten slide dwell times, tighten segments to under 8 minutes, and add interactive prompts (polls, chat questions, direct name-calling) every 5–7 minutes.

**Slides on virtual:** share a window, not your full desktop. Use Presenter View so you can see notes while the audience sees slides. Pause screen-sharing during genuine discussion — shared slides during Q&A reduce eye contact and create a passive mode.

## Anti-Patterns to Eliminate

- **Reading slides verbatim.** If the slide says it and you say it, one of you is redundant. The audience reads faster than you speak; they will finish before you and disengage.
- **Walls of bullets.** Each bullet resets the audience's attention to near-zero. Prefer one strong sentence over seven weak bullets.
- **"As you can see here..."** If it is visible, do not narrate its existence. Tell them what it means.
- **Apologizing for slides.** "Sorry this is busy" tells the audience the slide is bad and you know it. Fix the slide or cut it.
- **Running long.** Going over time disrespects every person in the room and every speaker after you. It also signals poor preparation. End early or on time, always.
- **The vague closing.** "So, in conclusion, we looked at X, Y, and Z..." is not a close. It is a summary trailing into silence. Design a real ending.
- **Animations for their own sake.** Fly-in, zoom, and spin effects draw attention to themselves, not to your content. Use transitions only to signal structure (new section) or to reveal a build deliberately.

## Sources

- Reynolds, Garr. "Presentation Zen" (2nd ed., New Riders, 2011)
- Duarte, Nancy. "Resonate" (Wiley, 2010)
- Duarte, Nancy. "slide:ology" (O'Reilly, 2008)
- Kawasaki, Guy. "The Art of the Start 2.0" (Portfolio/Penguin, 2015); "The 10/20/30 Rule of PowerPoint" (blog post, 2005)
- Anderson, Chris. "TED Talks: The Official TED Guide to Public Speaking" (Houghton Mifflin Harcourt, 2016)
- McKinsey & Company. Slide writing conventions (internal style guide, widely cited)
- BCG. Pyramid Principle-adjacent headline-deck standards (widely cited)
