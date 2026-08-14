---
name: kill-the-AI-ism
version: 1.2.0
updated: 2026-06-11
description: >
  Diagnostic skill for detecting and replacing generator artifacts ("AI-isms") in prose.
  Runs a structured heuristic sweep across four pattern categories — sentence-initial tells,
  structural tells, vocabulary tells, and formatting tells — and produces a findings report
  with severity flags and human-voice replacements for every hit. Distinct from
  writing-expert's banned-term enforcement: this skill fixes artifacts in an existing draft;
  writing-expert enforces rules during drafting. Output: annotated findings list + heuristic
  score summary (H1-H7). TRIGGER: "remove AI-isms", "make this sound human", "ChatGPT voice",
  "sounds like a bot", "de-robot this", "clean up generator artifacts" (voice/tone scope
  only; for broader structural/factual critique run document-critique after). SKIP: drafting
  from scratch or plain-language simplification with no AI-voice concern (use writing-expert);
  full document critique beyond voice/tone (use document-critique).
origin: local
category: custom
tags:
  - anti-ai-ism
  - voice
  - editing
  - prose-quality
  - human-writing
  - detection
triggers:
  - remove AI-isms
  - make this sound human
  - ChatGPT voice
  - AI-generated
  - this sounds like a chatbot
  - anti-AI-ism check
  - kill the AI tells
  - sounds like a bot
  - LLM artifacts
  - de-robot this
  - clean up generator artifacts
  - this sounds like it was written by AI
  - robot voice
skip:
  - drafting from scratch → use writing-expert
  - readability or plain-language simplification without AI-voice concern → use writing-expert
  - broad document critique beyond voice → use document-critique
whenToUse:
  - user says a draft "sounds like ChatGPT" or "sounds like a bot"
  - user wants to remove AI-isms, LLM artifacts, or generator tells from existing prose
  - user asks to make a document "sound more human"
  - user wants an anti-AI-ism audit or check before publishing
  - user says "de-robot this" or "kill the AI voice"
  - user has an AI-generated draft that needs voice normalization before use
  - user asks to clean up generator artifacts (voice scope — not full structural critique)
related_skills:
  - document-critique
---

# Kill the AI-ism

Diagnostic and replacement guide for generator artifacts in prose. A generator artifact is any
pattern a language model produces at high frequency that a human writer would use rarely or
never in that context. The goal is not to hide that AI was involved; it is to produce prose that
respects the reader's time and does not signal "I did not think hard about this."

Sources: Steven Pinker, *The Sense of Style* (2014); Joseph Williams, *Style: Lessons in Clarity
and Grace* (10th ed.); Ethan Mollick, "One Useful Thing" newsletter analyses of GPT/Claude
output patterns (2023-2025); the `writing-expert` skill (sibling — shares the Tier 1 ban list).

---

## Pattern catalog

### 1. Sentence-initial tells

These phrases appear at the start of a sentence or response. They cost words and communicate
servility, not thought.

| Pattern | Why it signals a generator | Human alternative |
|---|---|---|
| "I'll go ahead and" | Narrates the act of helping instead of helping | Delete. Start the sentence with the action. |
| "Let me" | Same as above | Delete, or rewrite as imperative/declarative |
| "I'd be happy to" | Emotional pre-validation | Delete. Just do it. |
| "Certainly!" | Filler affirmation with exclamation | Delete |
| "Great question!" | Sycophantic opener | Delete |
| "Of course" | Implies the request was obvious — often condescending | Delete |
| "Absolutely!" | Variant of Certainly! | Delete |
| "Sure!" | Same family | Delete |

**Rule:** if the first word of a response is an affirmation that adds no information, delete it
and start with the first substantive word.

### 2. Structural tells

These are layout and architecture patterns, not word-choice problems.

| Pattern | Threshold for concern | Fix |
|---|---|---|
| Introductory throat-clearing | Any paragraph whose sole function is to restate what the user just asked | Delete the paragraph |
| "Here's a summary:" headers | Used before a passage the reader can already see is a summary | Delete the header |
| Excessive bullet nesting | More than 2 levels of indentation; bullets where prose would read faster | Flatten to prose |
| Uniform paragraph length | All paragraphs within 20% of mean length (see heuristic below) | Vary intentionally: one punchy, one developed |
| "In conclusion" | Announces an ending the reader can already see | Delete; let the last sentence land on its own |
| Numbered lists for non-sequential content | Numbers imply order or rank; bullets or prose serve better | Change to prose or bullets |
| Mirror-structure responses | Response section headers copy the user's question verbatim | Rewrite headers as claims, not echoes |

### 3. Vocabulary tells

#### Tier 1 — ban entirely (any occurrence is a finding)

delve, leverage, robust, paradigm, seamless, utilize, commence, facilitate, furthermore,
navigate, landscape, cutting-edge, holistic, plethora, "it's worth noting", "ultimately"
(as a summary throat-clearer), "in conclusion", "it's important to note",
"in today's rapidly evolving", "game-changer"

These twenty terms appear at 3-8x higher frequency in LLM output than in comparable
human-written technical or business prose (per 2023-2025 corpus analyses of GPT-4 and
Claude 2/3 outputs). A human writer reaching for "delve" or "leverage" in a technical
context should hear an alarm.

**Replacement table — Tier 1:**

| AI-ism | Human alternatives |
|---|---|
| delve | examine, look at, study, explore (or just start doing it) |
| leverage | use, apply, draw on |
| robust | strong, reliable, proven, durable (or be specific: "handles 10k req/s") |
| paradigm | model, approach, method, shift (be specific) |
| seamless | smooth, uninterrupted (or describe what specifically is smooth) |
| utilize | use |
| commence | start, begin |
| facilitate | help, enable, allow, run |
| furthermore | also, and, besides (or restructure to avoid the connector) |
| navigate | handle, manage, work through (or be literal if navigation is meant) |
| landscape | field, market, environment (or name what you mean) |
| cutting-edge | current, recent, new, leading (cite evidence instead) |
| holistic | whole, complete, end-to-end (or say what parts are included) |
| plethora | many, dozens of, a long list of (or give the count) |
| it's worth noting | [delete + state the note directly] |
| ultimately | [delete or replace with the actual conclusion: "so", "therefore", "as a result"] |
| in conclusion | [delete + let the final sentence stand alone] |
| it's important to note | [delete + just state the point] |
| in today's rapidly evolving | [delete] |
| game-changer | [be specific about what changed] |

#### Tier 2 — watch, replace if not precise

transform, innovative, dynamic, significant, comprehensive, multifaceted, nuanced, intricate

These words are not banned, but generators overuse them as intensifiers when the underlying
claim is weak. Test: can you replace the word with a specific fact? If yes, use the fact.
"This produces a significant improvement" → "This cuts latency by 40%."

#### Hedging adverbs

potentially, possibly, perhaps, arguably, somewhat, relatively

One hedge per 500 words is fine — uncertainty is real. More than three per 500 words
suggests the generator hedged to avoid being wrong rather than because the content is
genuinely uncertain. Collapse hedges into a single explicit uncertainty statement
("We don't know X yet") and remove the rest.

#### Empty quantifiers

"various", "numerous", "several", "many" — without a count or examples.

A human writer who knows the subject gives a number or an example. Replace with the count
("four approaches"), a representative example ("tools like X, Y, and Z"), or delete the
quantifier and name the thing directly.

### 4. Formatting tells

| Tell | Threshold | Fix |
|---|---|---|
| Every paragraph as a bullet list | >40% of lines are bullets in a prose document | Convert to flowing paragraphs; reserve bullets for genuinely enumerable items |
| Em-dash glut | >2 em-dashes per 100 words | Cut to 1/100 or fewer; use commas, parentheses, or recast the sentence |
| Bold key terms in every paragraph | Bold appears in every paragraph (not just in reference docs) | Reserve bold for critical warnings or lookup anchors only |
| "Note:" / "Tip:" admonitions in flowing prose | Any | Integrate the content into the sentence; admonition boxes belong in docs, not memos |
| Closing pleasantries | "I hope this helps!", "Let me know if you have any questions", "Feel free to ask" | Delete entirely |

---

## Detection heuristics (run on a draft)

Run these mechanically before doing a word-level pass. Each one produces a flag or a pass.

**H1 — Em-dash density (per-context table — authoritative for every consumer of this skill)**
Count em-dashes (—) and divide by word count. Structural use in headings and table
separators is exempt in every context. Apply the row matching the document's context:

| Context | Pass | Flag |
|---|---|---|
| Formal external docs (customer-facing deliverables, RFCs, KB articles, published prose) | ≤1 per 300 words | >1 per 300 words |
| Polished prose (default) | ≤1 per 100 words | >2 per 100 words |
| Internal docs / Slack | Relaxed — no numeric threshold | Flag only when density hurts readability |

Strict AI-ism detection target: ≤1 per 500 words — apply this bar only when the explicit
goal is detecting generator authorship (the writing-expert Tier 3 detection case) rather
than editing prose.

**H2 — Bullet density**
Count lines that begin with a bullet or list marker. Divide by total lines.
- Pass: ≤40%
- Flag: >40% in a prose document (not a reference table or checklist)

**H3 — Sentence-initial affirmation**
Scan the first word or phrase of each paragraph or response block.
Flag any: "Certainly", "Absolutely", "Great", "Sure", "Of course", "I'd be happy",
"Let me", "I'll go ahead".
- Pass: 0 occurrences
- Flag: any

**H4 — Tier 1 banned-term density**
Count occurrences of the 20 Tier 1 terms.
- Pass: 0 per 1000 words
- Flag: any occurrence

**H5 — Hedge-adverb density**
Count: potentially, possibly, perhaps, arguably, somewhat, relatively.
- Pass: ≤3 per 500 words
- Flag: >3 per 500 words

**H6 — Paragraph-length variance**
Compute mean paragraph length in words. Flag if every paragraph is within 20% of the mean.
Suspicious uniformity — humans vary sentence rhythm instinctively; generators don't.
- Pass: at least one paragraph deviates >20% from mean
- Flag: all paragraphs within 20% of mean

**H7 — Closing pleasantry scan**
Check the final sentence of the document or response.
Flag: "hope this helps", "let me know", "feel free", "don't hesitate".
- Pass: 0 occurrences
- Flag: any

---

## The "would a human write this?" sniff test

When heuristics pass but something still feels off, ask:

1. Does this sentence say anything a reader could not already see from context?
   If no — delete it.
2. Would a subject-matter expert write this exact phrasing in an email to a peer?
   If no — rewrite it.
3. Is every paragraph the same shape (topic sentence, two elaborations, transition)?
   If yes — break the pattern in at least one paragraph.
4. Does the document anticipate every possible question before the reader asks it?
   Generators over-complete; humans leave some things to the reader.
   If yes — cut the anticipatory sections and let the prose breathe.

---

## When AI-isms are correct

Do not apply this skill when:

- The interface is explicitly an AI assistant (chatbot UI, support widget). Phrases like
  "I'd be happy to help" are expected and appropriate in that context.
- The document must disclose AI involvement (regulated content, AI-generated content labels).
  Authentic disclosure is not an artifact.
- The author's personal voice genuinely uses these phrases. Voice is an input, not an output.

---

## Compose with sibling skills

- **writing-expert** — enforces Tier 1 ban during drafting and handles multi-pass tightening,
  readability, and rhythm; use kill-the-AI-ism to audit after the fact or when the draft came
  from another source, then writing-expert for brevity and grade-level targeting
- **document-critique** — broader structural/factual review; kill-the-AI-ism handles voice
  only, document-critique handles everything else
