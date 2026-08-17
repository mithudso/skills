<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `ai-collaboration-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: ai-collaboration-writing
version: "1.0.0"
updated: "2026-05-29"
category: custom
tags:
  - writing
  - ai-writing
  - human-ai-collaboration
  - prompt-craft
  - disclosure
description: >
  How writers work with LLMs as collaborators — not replacements — and ship
  prose that still sounds like themselves. Covers the four workflow modes
  ("I draft, AI edits", "AI drafts, I edit", "I prompt with a style sample",
  "AI redlines, I decide"), the voice-transfer problem (LLMs flatten
  distinctive voice), few-shot style-sample prompting, the redline-not-rewrite
  pattern, anti-AI-ism enforcement during AI use, the 2026 disclosure ethics
  consensus ("AI-assisted" vs "AI-written" vs invisible-use grammar/spell
  checking), the FTC and New York 2026 disclosure rules, journal and academic
  AI-declaration norms, and prompt-style-guide artifacts. References
  Anthropic's prompt engineering guide, OpenAI cookbook patterns, and Ethan
  Mollick's *Co-Intelligence*.
  TRIGGER: "use AI to edit", "AI drafted", "prompt for my style", "voice
  transfer", "AI flattens voice", "redline this", "AI disclosure",
  "AI-assisted vs AI-written", "ghostwriting with AI", "style sample prompt",
  "prompt-style-guide", "AI writing ethics", "FTC AI disclosure", "AI
  declaration", "how to write with Claude", "how to write with ChatGPT",
  "few-shot style example", "anti-AI-ism while using AI", "human-AI
  collaboration writing", "writing-with-AI workflow".
  SKIP: optimizing a production system prompt (use prompt-deep-optimizer);
  general prompt engineering reference (use prompt-engineering); writing
  prose without an AI collaborator (use writing-expert); plain-language for
  customers (use plain-language); building an AI app via the SDK (use
  claude-api).
triggers:
  - "use AI to edit"
  - "AI drafted"
  - "prompt for my style"
  - "voice transfer"
  - "AI flattens voice"
  - "redline this"
  - "AI disclosure"
  - "AI-assisted vs AI-written"
  - "ghostwriting with AI"
  - "style sample prompt"
  - "prompt-style-guide"
  - "AI writing ethics"
  - "FTC AI disclosure"
  - "AI declaration"
  - "how to write with Claude"
  - "how to write with ChatGPT"
  - "few-shot style example"
  - "anti-AI-ism while using AI"
  - "human-AI collaboration writing"
  - "writing-with-AI workflow"
skip:
  - production system-prompt optimization → use prompt-deep-optimizer
  - general prompt engineering reference → use prompt-engineering
  - prose without an AI collaborator → use writing-expert
  - plain-language simplification → use plain-language
  - building AI apps via SDK → use claude-api
related:
  - writing-expert
  - prompt-engineering
  - prompt-deep-optimizer
  - kill-the-AI-ism
  - editing-and-revision
  - plain-language
---

# AI-Collaboration Writing

Reference for writers working with LLMs as a collaborator — picking the right
workflow mode, preserving distinctive voice, applying anti-AI-ism discipline
mid-collaboration, and meeting the 2026 disclosure standards. The goal is
prose that still sounds like the human author. AI is the second writer in the
room, not the first.

Deliver all responses in a direct, plain register. Avoid hedging and meta-commentary.

---

## When to use this skill

Activate when the user:

- Wants AI to edit their draft without sanding off their voice
- Drafted with AI and now needs to humanize it
- Asks how to prompt with a style sample
- Reports that AI output sounds "generic" or "flattened"
- Asks whether to label output as "AI-assisted" or "AI-written"
- Needs to comply with the FTC AI disclosure rule, a journal AI-declaration
  policy, or a corporate AI policy
- Wants a prompt-style-guide artifact for their team
- Is choosing between AI as drafter vs AI as editor for a piece
- Wants the redline-not-rewrite pattern applied
- Is enforcing anti-AI-ism rules in their AI-collaborated output

Skip when:

- The task is optimizing a production system prompt (use
  `prompt-deep-optimizer`)
- The task is prompt-engineering reference, not writing collaboration (use
  `prompt-engineering`)
- The user is writing alone with no AI in the loop (use `writing-expert`)
- The task is simplification for non-technical audiences (use
  `plain-language`)
- The user is building an SDK-based application (use `claude-api`)

---

## The one rule: voice is human-load-bearing; AI is everything else

A writer's distinctive voice is the residue of their judgment. An LLM,
trained on the global mean of internet prose, regresses toward that mean.
Left to itself, AI flattens voice. The collaboration mode this skill
defends keeps the human responsible for voice and uses AI for tasks where
the mean is fine (mechanics, structure, length compression, format
conversion).

Voice is human. Everything that can be measured against a mean is fair game
for AI.

---

## Core concept 1 — The four collaboration modes

Pick the mode before opening the chat window.

| Mode | Who drafts | Who edits | Who decides | Best for |
|------|------------|-----------|-------------|----------|
| **I draft, AI edits** | Human | AI | Human | High-voice prose: essays, opinion, memos with author identity |
| **AI drafts, I edit** | AI | Human | Human | Boilerplate: status updates, release notes, structured docs |
| **I prompt with a style sample** | AI (constrained) | Either | Human | New work in the author's existing voice (newsletters, recurring docs) |
| **AI redlines, I decide** | Human | AI suggests changes | Human reviews suggestions | Manuscript-level critique without overwriting the source |

The wrong choice ruins the result. "AI drafts, I edit" on an op-ed produces
generic copy and a tired editor; "I draft, AI edits" on a release-note
matrix wastes hours.

---

## Core concept 2 — The "I drafted, now you edit" workflow

For voice-heavy prose. The human owns the draft; the AI does mechanical and
structural editing.

**Prompt skeleton:**

```text
You are an editor working on my draft. Do not rewrite it. Make four passes:

1. Mechanics — typos, grammar, punctuation, agreement errors.
2. Cohesion — flag sentences that don't connect to their predecessor.
3. Cuts — flag redundancy, hedging, and filler. Don't delete; suggest.
4. Anti-AI-isms — flag any sentence that sounds machine-written.

Return a numbered list of suggestions, not a rewritten draft. For each
suggestion, quote the source sentence and propose the edit. I will accept
or reject each one.

[paste the draft]
```

Two reasons this works:

1. The AI never has authority over the prose. The human accepts or rejects
   every edit.
2. The output is a list of diffs, not a replacement. The author sees what
   changed and why.

---

## Core concept 3 — The "AI drafted, now I edit" workflow

For low-voice prose where the goal is correctness and completeness, not
identity. Status updates, release notes, runbook stubs, internal wikis.

**Workflow:**

1. Prompt the AI to draft, with explicit structural constraints.
2. Receive the draft.
3. **Read every sentence aloud.** If a sentence has rhythm the human would
   not use, rewrite it.
4. Strip AI-isms (see core concept 8).
5. Add the one or two specific facts only the human knows (which broke,
   why, who fixed it).
6. Ship.

The risk: skipping step 3. AI draft + cursory skim + ship = the prose
everyone now recognizes as machine-written. The author hates it later, even
if no reader notices today.

---

## Core concept 4 — The voice-transfer problem

LLMs are trained on the mean of the internet. Average prose pulls them
toward five default tendencies:

1. **Hedging.** "It's worth noting that…", "It seems that…", "Arguably…"
2. **Triadic rhythm.** "Clear, concise, and effective." Three-item lists
   replace one strong word.
3. **Topic-sentence tax.** Every paragraph opens with a meta-sentence
   ("In this section, we will explore…").
4. **Latinate verbs.** "Utilize," "facilitate," "leverage," "demonstrate,"
   "exhibit."
5. **Bridge phrases.** "Furthermore," "moreover," "in addition,"
   "additionally."

A distinctive human voice violates several of these. AI editing or AI
drafting will regress them toward the mean unless the prompt actively
defends against it.

---

## Core concept 5 — The style-sample prompt

The most reliable way to keep an AI's output close to a specific voice is
**few-shot prompting with a style sample**. Show the AI 2–4 paragraphs of
the author's existing work and instruct it to match.

**Prompt skeleton:**

```text
Here are three paragraphs of my recent writing. Note the rhythm, sentence
length, vocabulary, and stance. Match this voice exactly.

SAMPLE 1:
[paste 1-2 paragraphs]

SAMPLE 2:
[paste 1-2 paragraphs]

SAMPLE 3:
[paste 1-2 paragraphs]

Now write a [memo / blog post / newsletter] on [topic] in this voice.

Rules:
- No hedging phrases ("it's worth noting", "arguably")
- No triadic lists where one strong word works
- No topic sentences that announce the section
- No latinate verbs (utilize, facilitate, demonstrate)
- No bridge phrases (furthermore, moreover)
- Sentences average 14-18 words, ranging 5-30
- Use semicolons sparingly; prefer hard sentence breaks
```

This usually gets within 80% of the author's voice. The remaining 20%
needs human editing.

---

## Core concept 6 — The redline-not-rewrite pattern

For manuscript-level critique without losing the original prose.

**Prompt skeleton:**

```text
Read the following draft. Don't rewrite it. Produce a redline review:

1. List the 5 strongest sentences (quote them, one line each).
2. List the 5 weakest sentences (quote them, one line each).
3. Identify the 3 places the argument actually moves.
4. Identify the 3 places the argument stalls.
5. Suggest one structural change at the paragraph level, if any.
6. Suggest one sentence-level change for each of the 5 weakest sentences.

Do not propose a new draft.

[paste draft]
```

The output is a critique, not a replacement. The author keeps authorship.

---

## Core concept 7 — The prompt-style-guide artifact

For teams using AI to draft repetitive prose (status updates, release
notes, customer emails), a **prompt-style-guide** is a reusable artifact
that captures the team's voice rules and feeds them to every prompt.

**Structure:**

```markdown
# Acme Engineering Writing — Prompt Style Guide

## Voice
- Direct, concrete, no marketing register.
- Sentences average 14-18 words.
- Active voice unless the agent is unknown.
- No hedging. State, don't suggest.

## Banned words (Tier 1 — never use)
- delve, leverage, robust, paradigm, seamless, utilize, commence,
  facilitate, furthermore, navigate, landscape, cutting-edge, holistic,
  ecosystem, harness (verb), foster (verb)

## Banned patterns
- Triadic lists used as filler ("clear, concise, and effective")
- Topic sentences that announce the section
- Bridge phrases ("furthermore", "moreover", "in addition")
- "It's worth noting that..." and synonyms

## Structural rules
- Open with the load-bearing fact (BLUF).
- One idea per paragraph.
- Code, command, or number in the first 100 words when relevant.
- No emoji unless explicitly requested.

## Tone calibration
- Engineer audience: assume domain knowledge.
- Executive audience: lead with the business impact.
- Customer audience: explain mechanics in plain language.

## Output format
- Markdown headings, bullet lists, code fences as appropriate.
- No frontmatter.
- No author-disclaimer / no AI-disclaimer in the output text itself.
```

Paste this guide into every AI session that produces team prose. It is the
single highest-leverage intervention.

---

## Core concept 8 — Anti-AI-ism enforcement during AI use

When using AI as a drafter, expect the model to insert AI-isms; expect to
strip them. The Tier 1 ban list from the `writing-expert` and
`kill-the-AI-ism` skills:

**Banned (Tier 1):** delve, leverage, robust, paradigm, seamless, utilize,
commence, facilitate, furthermore, navigate, landscape, cutting-edge,
holistic, ecosystem, harness (verb), foster (verb)

**Banned phrases:**

- "It's worth noting that…"
- "In today's fast-paced world…"
- "Navigating the landscape of…"
- "It's important to remember that…"
- "Furthermore, it should be noted…"
- "In this article, we will explore…"
- "Let's dive into…"

**Banned structural patterns:**

- Three-bullet triadic conclusions ("In summary, X is clear, concise, and
  effective")
- Topic sentences that announce ("In this section…")
- Closing paragraphs that recap the entire piece

Two enforcement approaches:

1. **Prompt-time:** put the ban list in the prompt-style-guide.
2. **Edit-time:** ctrl-F the banned words in the AI output; rewrite or
   delete the sentence.

---

## Core concept 9 — Disclosure ethics (2026 consensus)

A practical hierarchy:

| Use level | Disclosure required? | Label |
|-----------|----------------------|-------|
| Grammar / spell check | No | None |
| Reference / citation formatting | No | None |
| Format conversion (Markdown → HTML) | No | None |
| AI-suggested edits, accepted/rejected by human | Usually no | "Edited with AI assistance" if disclosure norm in venue |
| AI drafted, heavily rewritten by human | Yes | "AI-assisted" |
| AI drafted, lightly edited by human | Yes | "AI-written, human-reviewed" |
| AI generated, published as-is | Yes (often required) | "AI-generated" |

**Regulatory landscape, 2026:**

- **FTC (USA):** disclosure required for AI-generated endorsements,
  reviews, and material claims in advertising.
- **New York State (June 2026):** explicit AI-disclosure requirements for
  influencer and brand AI use.
- **Academic journals:** most require an "AI-use declaration" naming the
  model and the role (drafting, editing, analysis). AI cannot be listed
  as an author.
- **News organizations:** majority require disclosure when AI generated a
  meaningful portion of a piece; "edited with AI" is generally not
  required.

**Two-line disclosure template (general purpose):**

```text
This piece was [drafted / edited / fact-checked] with assistance from
[model name]. The author wrote the [draft / final version] and is
responsible for all content, including any errors.
```

---

## Core concept 10 — Choosing the right model and feature for the task

Not all AI work needs the most powerful model. Match capability to task:

| Task | Model tier | Feature |
|------|-----------|---------|
| Spelling, grammar, mechanics | Smallest | None |
| Status update / release note draft | Mid | None |
| Style-sample matching | Mid–top | Few-shot prompting |
| Manuscript critique | Top | Extended thinking if available |
| Voice-preserving edit pass | Top | Style-sample + ban-list prompt |
| Multi-hour collaboration on a long piece | Top | Memory / project context if available |

This is a writing skill, not an AI app-building skill — but the writer
benefits from knowing that smaller models suffice for mechanical work,
freeing budget for the voice-critical passes.

---

## Templates

### Universal "AI edits my draft" prompt

```text
You are an editor. Do not rewrite my draft. Suggest changes as a numbered
diff list. For each suggestion: quote the source sentence verbatim, then
show the proposed change, then give a one-line reason.

Bans:
- Do not introduce any of these words: delve, leverage, robust, paradigm,
  seamless, utilize, commence, facilitate, furthermore, navigate,
  landscape, cutting-edge, holistic, ecosystem.
- Do not introduce hedging ("it's worth noting", "arguably").
- Do not introduce triadic filler ("clear, concise, and effective").
- Do not introduce bridge phrases ("furthermore", "moreover").

Output format:
1. [Quote] → [Proposed change] — [Reason]
2. [Quote] → [Proposed change] — [Reason]
...

DRAFT:
[paste]
```

### Universal "AI drafts in my voice" prompt

```text
Match the voice of the following sample exactly. Note the sentence rhythm,
length variation, vocabulary, and stance.

VOICE SAMPLE:
[paste 200-400 words of your own writing]

WRITE: a [length-target] [piece type] on [topic].

RULES:
- Sentence length: average 14-18 words; range 5-30; no run-ons.
- Vocabulary: Anglo-Saxon-rooted preferred; minimize latinate.
- Voice: active. Subject-verb-object.
- Banned: delve, leverage, robust, paradigm, seamless, utilize, commence,
  facilitate, furthermore, navigate, landscape, cutting-edge, holistic.
- No triadic filler. No topic-announcement sentences. No bridge phrases.
- One idea per paragraph. Open paragraphs with the load-bearing fact.
- End on a specific image or concrete detail, not a recap.

Return only the draft. No preamble. No closing remarks.
```

### Universal "AI redlines my manuscript" prompt

```text
You are a manuscript reviewer. Do not rewrite anything. Produce a redline
review with these six sections:

1. STRONGEST 5 SENTENCES — quote each, one line.
2. WEAKEST 5 SENTENCES — quote each, one line. Why each fails.
3. STRUCTURE — list the 3 places the argument moves; the 3 places it stalls.
4. VOICE — does this sound like one author, or like committee prose? Cite.
5. ANTI-AI-ISMS — list any sentence that reads as machine-generated. Cite.
6. ONE STRUCTURAL CHANGE — propose at most one paragraph-level reorg.

Do not produce a rewritten draft. I will edit the source.

MANUSCRIPT:
[paste]
```

### Two-line AI disclosure (general)

```text
This [piece / report / post] was drafted with assistance from [model name].
The author wrote the final version, made all editorial decisions, and is
responsible for all content.
```

### Prompt-style-guide skeleton (team artifact)

```markdown
# [Team] Prompt Style Guide v[N]

## Voice
- [3-5 bullets]

## Banned words (Tier 1)
- [list]

## Banned patterns
- [list]

## Structural rules
- [list]

## Tone by audience
- Engineer: [...]
- Executive: [...]
- Customer: [...]

## Output format
- [Markdown / plain text / structured fields]

## Length budget
- [target word count or section budget]

## Disclosure
- [in-document AI disclosure policy, if any]
```

---

## Anti-patterns

| Anti-pattern | Why it fails | Fix |
|--------------|--------------|-----|
| AI drafts a voice-heavy essay | Voice flattens to internet mean | "I draft, AI edits" mode |
| AI edits without a ban list | Inserts more AI-isms than it removes | Prompt-style-guide with banned words |
| Single prompt: "Make this better" | AI rewrites and erases voice | "Suggest a numbered diff list, don't rewrite" |
| No style sample on creative tasks | Output drifts to generic | Paste 200-400 words of prior work |
| Disclosure missing where required | FTC / journal / venue violation | Default to disclose for "AI drafted" cases |
| Disclosure on grammar-check only | Over-claims; misleading | Only disclose when AI generated content |
| "Polish this draft" without rules | AI rewrites voice into mean | List banned words and patterns in the prompt |
| Trusting the AI's self-edit | Model can't see its own AI-isms | Human ctrl-F the ban list afterward |
| Drafting in AI then claiming sole authorship | Ethics violation | Disclose appropriately |
| Same prompt for every team writer | Voice diverges or homogenizes wrong | Team prompt-style-guide as the shared artifact |
| Letting AI add closing recap | AI-ism signal | Strip closing recaps in edit pass |
| "Make it more engaging" | Triggers triadic filler & breathless tone | Specify the change concretely |

---

## Decision heuristics

**Pick the workflow mode:**

- Is the prose voice-load-bearing (essay, opinion, customer-facing letter)?
  → "I draft, AI edits"
- Is the prose mechanical (release note, status update, schema doc)?
  → "AI drafts, I edit"
- Is the user shipping new prose in a known recurring voice (newsletter,
  recurring memo)? → "Style-sample prompt"
- Is the user reviewing their own manuscript? → "AI redlines"

**Use AI as drafter only when:**

- The structure is standard and largely template-driven
- The vocabulary is technical, low-stakes, low-voice
- The human will read every sentence aloud after

**Voice-preservation check:**

- Read the output aloud. If you would not say it that way, rewrite it.
- Grep for Tier 1 banned words. Strip them.
- Count triadic lists. Cut at least half.

**Disclosure decision:**

- Did AI generate text that survives in the final output? → Disclose
- Did AI only check grammar / format? → No disclosure required
- Is the venue an academic journal, FTC-regulated advertising, NY
  state-regulated content, or a publication with an explicit AI policy?
  → Follow the venue's required format

**Prompt-style-guide rollout:**

- One team, one document
- Update quarterly as voice evolves
- Test the guide by giving an AI a fresh topic and the guide; review output

---

## References

- Anthropic: [Prompting best practices](https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/claude-prompting-best-practices),
  the prompt engineering overview, and Anthropic's writing-with-Claude docs
- OpenAI Cookbook: prompt patterns, including few-shot and persona-grounded
  prompting
- Ethan Mollick: *Co-Intelligence: Living and Working with AI* (2024) and
  his "One Useful Thing" Substack — practical writer-with-AI workflows
- FTC: [AI Disclosure Rules for Creators
  (2026)](https://blog.promise.legal/startup-central/ftc-ai-disclosure-rules-creators-2026/)
- New York State: [AI Disclosure Law (June
  2026)](https://humanadsai.com/blog/new-york-ai-disclosure-law-2026)
- ICMJE (medical) and journal-side AI declaration guidance: e.g., [Journal
  of Clinical Question generative-AI
  declaration](https://jclinque.com/news/generative-ai-declaration-for-manuscript-submissions)
- AZHIN: [Disclosure and Attribution of AI and Writing
  Tools](https://azhin.org/cummings/disclosure-attribution-ai)
- arXiv (2024–2026): empirical studies on reader perception of AI disclosure
  in writing (e.g., 2410.04545, 2604.27129)

---

## Related skills

- `writing-expert` — for the prose craft this skill defends
- `kill-the-AI-ism` — for the surgical removal of AI-isms in output
- `prompt-engineering` — for prompt-engineering reference
- `prompt-deep-optimizer` / `prompt-helper-optimizer` — for production
  prompt files, not collaboration prompts
- `editing-and-revision` — for the post-AI edit pass
- `plain-language` — when AI-assisted prose must hit a reading-grade target
