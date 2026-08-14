<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `brand-voice-guide-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: brand-voice-guide-writing
description: >
  Authoring craft for the brand voice guide artifact — the document other writers will
  consult. Covers the four-section template (Voice, Tone, Vocabulary, Examples), the
  "we are X, we are not Y" voice-attributes framing, the brand-voice spectrum (3-axis
  chart: formal↔casual, serious↔playful, respectful↔irreverent), the tone-shifts
  matrix (different tones for different surfaces/situations), the do/don't pairs
  convention, voice-guide versioning and review discipline, the voice-review process
  for new content, and the Mailchimp/Slack/Atlassian/Twilio/MailerLite lineage.
  Distinct from style-sheet authoring inside writing-expert — this skill is about
  writing the GUIDE that codifies voice for an organization or product.
  TRIGGER: "write a brand voice guide", "voice and tone document", "voice attributes",
  "voice spectrum chart", "tone shifts matrix", "do/don't pairs", "voice guide template",
  "we are not", "Mailchimp-style voice guide", "voice chart", "brand voice playbook",
  "how should writers use our voice", "voice review process", "voice guide template".
  SKIP: writing the prose itself in a specific voice (use writing-expert, sales-and-marketing-copy,
  or microcopy-and-ui-writing); a style sheet of grammar/spelling decisions for a single
  document (writing-expert covers style-sheet conventions); brand strategy or naming
  (out of scope); accessibility writing rules (use accessibility-writing).
origin: local
category: custom
version: "1.0.0"
updated: "2026-05-29"
tags:
  - writing
  - brand-voice
  - style-guide
related_skills:
  - writing-expert
  - sales-and-marketing-copy
  - microcopy-and-ui-writing
  - technical-writing-craft
  - policy-and-governance-writing
---

# Brand Voice Guide Writing

Authoring craft for the brand voice guide *as an artifact*. The output is a document
other writers will open, reference, and apply — not a style sheet for one piece, not a
brand strategy deck. Source lineage: Mailchimp Content Style Guide, Slack Brand Voice,
Atlassian Design Voice and Tone, Twilio Paste voice docs, MailerLite voice doc, Margot
Bloomstein *Content Strategy at Work* (2012) and BrandSort method, Nicole Fenton &
Kate Kiefer Lee *Nicely Said* (2014).

## When to use this skill

- Your org or product needs a voice guide and one doesn't exist.
- A voice guide exists but writers ignore it (usually means it's abstract, dated, or unwieldy).
- You're refreshing a voice guide because the brand, audience, or product surface area shifted.
- You need to add a "voice review" gate to your content pipeline and need the rubric.

## Overview

A brand voice guide answers four questions for every writer who opens it:

1. **Who do we sound like?** (Voice — the constant.)
2. **How do we shift in different situations?** (Tone — the variable.)
3. **Which words do we use, and which do we never use?** (Vocabulary.)
4. **What does this look like in practice?** (Examples — usually do/don't pairs.)

If a guide doesn't make these four questions answerable in under 10 minutes of reading,
it will be ignored. The single biggest failure mode of voice guides is being too long,
too abstract, and undated — so writers don't trust the guidance is current.

## Core Concepts

### 1. The four-section skeleton

The Mailchimp/Atlassian/Twilio canon converges on four sections, in this order:

1. **Voice** — three to five attributes that describe the constant personality. State each
   as "We are X. We are not Y." (Bloomstein's BrandSort framing: include the *not*.)
2. **Tone** — a matrix showing how voice flexes by situation. Rows = surfaces or moments;
   columns = recommended tonal shifts.
3. **Vocabulary** — preferred words, banned words, and grammar/style decisions specific to
   the brand (capitalization of product names, preferred contractions, oxford comma, etc.).
4. **Examples** — at least 5 do/don't pairs spanning the highest-traffic surfaces (homepage
   hero, error messages, support replies, marketing email, social post). Real examples,
   not paraphrased.

Anything not falling into one of these four buckets probably belongs in a different
document (brand strategy, content strategy, accessibility guide, editorial calendar).

### 2. Voice attributes — the "we are X, we are not Y" pattern

Bloomstein's BrandSort method asks teams to sort attribute cards into three piles: who
we are, who we'd like to be, and who we're *not*. The "we're not" pile is load-bearing —
without it, attributes drift toward generic ("friendly, professional, trustworthy")
which says nothing.

The canonical format for the Voice section:

> **Confident, not arrogant.** We make declarative claims when we know we're right. We
> don't oversell, hedge with weasel words, or talk down to readers who know less than us.
>
> **Warm, not saccharine.** We sound like a competent colleague who likes their job. We
> don't use exclamation marks to perform enthusiasm or call readers "friend."

Three to five attributes is the sweet spot. Two is thin. Six or more is unmemorable.

### 3. The brand-voice spectrum chart (typically 3 axes)

The Nielsen Norman Group axes are: formal↔casual, serious↔funny, respectful↔irreverent,
enthusiastic↔matter-of-fact. Most voice guides pick 3 of these 4 and place a marker on
each axis. The marker is the "voice fingerprint."

```text
Formal        |---------------●---|        Casual
Serious       |---●---------------|        Playful
Matter-of-fact|---------●---------|        Enthusiastic
```

Why this beats a paragraph: a writer can hold three positions in working memory and
auto-correct a sentence against them in seconds. A paragraph of attributes can't be
recalled mid-draft.

**Anti-pattern.** Putting all markers in the middle. If your voice is "moderately
formal, moderately serious, moderately enthusiastic," your voice is "no voice."

### 4. The tone-shifts matrix

Voice is constant. Tone shifts. The matrix is a small grid showing how the voice flexes
for the situations writers actually face. Mailchimp's canonical example shifts tone for
"user just succeeded" vs "user just failed" vs "user is reading help docs."

| Situation | Tone shift | Example |
|---|---|---|
| User succeeded (deploy passed) | Brief, warm, low-key | "Deployed in 47s. You're live." |
| User failed (deploy errored) | Direct, calm, no blame | "Deploy failed at build step. Logs below." |
| Onboarding (first-run) | Patient, encouraging, no jargon | "First time? We'll walk you through it." |
| Marketing (new feature) | Confident, specific, no hype | "Indexes now rebuild without a write pause." |
| Legal / policy update | Plain, precise, no marketing | "We updated our privacy policy on May 1." |

Three to seven rows. More than that and writers won't read it.

### 5. Do/don't pairs (the "this but not that" convention)

Fenton & Kiefer Lee's *Nicely Said* popularized the "This But Not That" exercise. Each
pair shows the same intent rendered two ways. The "don't" is realistic — usually a
common drift, not a strawman.

> **Do:** "We paused your build. Resume it when ready."
> **Don't:** "Oops! Looks like your build got paused. No worries — just click Resume
> whenever you'd like to keep going! 🚀"

The don't side teaches more than the do side. Pick don'ts from real drafts your team
has actually shipped; never invent strawmen.

### 6. Vocabulary section: preferred, banned, and decisions

Three sub-lists, each short:

- **Preferred** — words that carry the voice. ("ship" not "release"; "you" not "the user".)
- **Banned** — words that break the voice. ("leverage", "synergy", "seamless", "robust",
  "delightful", "passionate".) Include the canonical AI-tell list if applicable.
- **Decisions** — product-name capitalization, oxford comma, contractions allowed, em dash
  vs en dash, "we" vs company name, headline case vs sentence case.

Keep each list to 10–20 entries. Long lists go unread.

### 7. Voice-guide lineage and what to borrow

| Source | What to borrow |
|---|---|
| Mailchimp Content Style Guide | The voice-vs-tone framing; the situational tone table; the "we are X, not Y" attribute format. |
| Slack Brand Voice | "Clear over clever" as a guiding principle; the bias toward shorter sentences. |
| Atlassian Design Voice and Tone | Principles ("clear, familiar, inclusive") tied directly to examples; product-surface tone shifts. |
| Twilio Paste (voice and tone) | Embedding voice rules directly inside the design system so writers and designers share one source. |
| MailerLite voice doc | Brevity — a voice guide can be 3 screens, not 30 pages. |
| Bloomstein *Content Strategy at Work* | The BrandSort method; making "we are not" part of the artifact. |
| Fenton & Kiefer Lee *Nicely Said* | The "This But Not That" exercise; do/don't pairs as the central teaching device. |

### 8. Dating and versioning discipline

Every voice guide should display, in the top metadata:

- **Version** (semver or date)
- **Last reviewed** date
- **Owner** (a person, not a team — a team owns nothing)
- **Next review** date (no more than 12 months out)

A voice guide without a date is dead on arrival. Writers assume undated docs are stale
and route around them. Schedule a review every 6–12 months even if "nothing changed" —
the review itself is the signal of trust.

### 9. The voice-review process

A voice guide that isn't applied is decoration. The lightest workable enforcement:

1. **Self-check** — writer runs the do/don't pairs against their draft before publishing.
2. **Peer review** — one other writer scans the draft for tone-matrix fit (5 min, not 30).
3. **Owner gate** — the voice owner reviews high-visibility surfaces (homepage hero,
   onboarding, error copy, anything customer-facing during incidents).

Document this process inside the guide. If it's not in the guide, it doesn't exist.

### 10. The "one voice, many writers" tension

Voice guides exist because more than one person writes for the brand. Tension shows up
when a senior writer's instinct disagrees with the guide. Two rules:

- **The guide wins on consistency.** Even if a senior writer's version is "better," the
  guide's version is the org's version. Inconsistency erodes brand more than mediocrity.
- **The guide changes when enough drafts argue against it.** If three writers in six
  months overrode the same rule, the rule is wrong. Update the guide; don't fight the
  drafts.

## Templates

### Template 1: The four-section skeleton

```markdown
# [Brand] Voice Guide

**Version:** 1.0 — Last reviewed: YYYY-MM-DD — Owner: [Name] — Next review: YYYY-MM-DD

## 1. Voice

We sound like ___, not ___. Three to five attributes, each in "we are X, not Y" form.

- **[Attribute 1], not [opposite].** [One sentence of what this means in practice.]
- **[Attribute 2], not [opposite].** [One sentence.]
- **[Attribute 3], not [opposite].** [One sentence.]

### Voice fingerprint (spectrum)

```text
Formal       |--------●----------|  Casual
Serious      |-●-----------------|  Playful
Matter-of-fact|----------●-------|  Enthusiastic
```

## 2. Tone

Voice is constant. Tone shifts by situation. Use this matrix.

| Situation | Tone | Example opener |
|---|---|---|
| [Situation 1] | [Tone shift] | "[Example]" |
| [Situation 2] | [Tone shift] | "[Example]" |
| [Situation 3] | [Tone shift] | "[Example]" |

## 3. Vocabulary

**Preferred:** [10–20 words/phrases that carry the voice.]
**Banned:** [10–20 words/phrases that break the voice.]
**Decisions:** [Capitalization, punctuation, contractions, etc.]

## 4. Examples

### Do/Don't pairs

> **Do:** "[Example]"
> **Don't:** "[Example that drifts off-voice.]"

[Minimum 5 pairs spanning the highest-traffic surfaces.]

## How to use this guide

1. Self-check against the do/don't pairs.
2. Peer review for tone-matrix fit.
3. Owner review for [list of high-visibility surfaces].

## Changelog

- YYYY-MM-DD — v1.0 — Initial release.
```

### Template 2: The voice-spectrum chart (ASCII)

```text
Voice fingerprint

Formal        |------●------------|  Casual           (lean casual)
Serious       |----●--------------|  Playful          (mostly serious)
Respectful    |---------●---------|  Irreverent       (centered; flex by surface)
Matter-of-fact|--------------●----|  Enthusiastic     (lean enthusiastic)
```

Place a single marker on each axis. Annotate with a one-word descriptor if the position
is non-obvious. Four axes is the maximum a writer can recall mid-draft; three is better.

### Template 3: Sample do/don't pairs

> **Surface: Error message after a failed deploy**
> **Do:** "Deploy failed at the build step. Logs below."
> **Don't:** "Oops! Something went wrong! 😬 No worries — let's get you back on track."

> **Surface: Onboarding first-run greeting**
> **Do:** "First time here? Two-minute setup."
> **Don't:** "Welcome to the future of [category]! We're so excited to have you on board."

> **Surface: Homepage hero subhead**
> **Do:** "Ship database changes without dropping a write."
> **Don't:** "Leverage our robust, seamless platform to unlock your data potential."

> **Surface: Incident status post**
> **Do:** "We saw elevated 5xx errors from 14:02–14:11 UTC. Service is restored. Postmortem to follow."
> **Don't:** "We experienced a brief hiccup today. Thanks for your patience — you're amazing!"

> **Surface: Support reply to a frustrated customer**
> **Do:** "You're right — that should have worked. I reproduced it on my end. Fix is in progress; I'll update you within 4 hours."
> **Don't:** "I totally understand your frustration! Let me dig into this for you and circle back ASAP."

## Anti-Patterns

1. **No "we are not."** Lists three positive attributes ("friendly, professional,
   trustworthy") with no opposites. Reader can't tell what off-voice looks like.
2. **Spectrum markers all centered.** "Moderately formal, moderately playful" is no
   guidance. Force a position on each axis.
3. **Don't-side strawmen.** Don't pairs that are obviously bad teach nothing. Use real
   drafts your team has shipped.
4. **No date, no owner.** Writers won't trust an undated guide. Always show last-reviewed
   and next-review dates plus a named owner.
5. **30+ page voice guide.** Anything past 10 pages won't be read. Move depth into a
   linked appendix; keep the main guide skim-friendly.
6. **Mixing brand strategy with voice.** Mission, vision, positioning belong in a
   different document. Voice is about *language*, not strategy.
7. **No tone matrix.** "Voice without tone" is the most common gap — writers can match
   personality but not situation. Always include the matrix.
8. **Banning words without preferring others.** "Don't say leverage" without "say use"
   leaves writers stuck mid-sentence. Pair every ban with a substitute.
9. **No "how to use this guide" section.** The review process must be inside the artifact
   or it won't happen.
10. **Copying Mailchimp's voice verbatim.** Their voice fits Mailchimp. Borrow the
    *structure*, write your own *attributes*.

## Decision Heuristics

- **Do we need a voice guide at all?** Yes, if more than one person writes for the
  brand. No, if it's a solo founder writing everything (just be consistent with yourself).
- **3 axes or 4?** 3 unless a fourth axis (typically respectful↔irreverent) is
  load-bearing for the brand. More than 4 axes is unrecallable.
- **How many tone-matrix rows?** 3 minimum (success, failure, neutral). 7 maximum. Aim
  for the 5–6 surfaces that account for 80% of customer-facing copy.
- **How many do/don't pairs?** Minimum 5. Maximum ~15. Each pair must teach a distinct
  rule; redundant pairs dilute the guide.
- **Length cap?** A working voice guide fits on ~10 letter-sized pages or ~3 screens of
  web content. Anything longer is reference, not guidance — move it to an appendix.
- **Refresh cadence?** 6 months for an early-stage brand or a rebrand; 12 months for
  steady-state. Always schedule the next review when you ship the current version.
- **Versioning?** Semver if you treat the guide like code; date-based (YYYY.MM) is fine
  for marketing-owned guides. Either way, dated and changelogged.
- **Where does it live?** Inside the design system if one exists (Twilio Paste model);
  otherwise a single Notion/Confluence page with a stable URL. PDFs go stale and lose
  their URL — avoid.

## References

- [Mailchimp Content Style Guide — Voice and Tone](https://styleguide.mailchimp.com/voice-and-tone/)
- [Atlassian Design — Voice and Tone](https://atlassian.design/content/voice-and-tone/) (canonical examples of principle-tied-to-example)
- [Twilio Paste — Voice and Tone](https://paste.twilio.design/foundations/content/voice-and-tone) (voice rules embedded in a design system)
- [Slack Brand Voice — analysis and lineage](https://www.eliteasia.co/brand-tone-of-voice-examples/) (third-party summary; Slack's full guide is internal)
- Bloomstein, M. *Content Strategy at Work*. Morgan Kaufmann, 2012 — introduces BrandSort and the message-architecture method.
- [Appropriate Inc. — Margot Bloomstein consulting / BrandSort](https://appropriateinc.com/) — the canonical "we are not" framing.
- Fenton, N. and Kiefer Lee, K. *Nicely Said: Writing for the Web with Style and Purpose*. New Riders, 2014 — the "This But Not That" exercise; voice vs tone primer.
- [Nielsen Norman Group — The Four Dimensions of Tone of Voice](https://www.nngroup.com/articles/tone-of-voice-dimensions/) — the canonical four-axis spectrum.
- [Brand Voice Playbook — Brand Vision](https://www.brandvm.com/post/brand-voice-playbook) — recent (2026) practitioner overview of voice-guide structure.
