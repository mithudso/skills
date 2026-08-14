# ai-assisted-copywriting-workflow

**Category:** Security, Auth & Diagnostics
**Platform:** Claude
**Original Path:** claude/standalone/ai-assisted-copywriting-workflow

## Description
Guide the end-to-end workflow for using LLMs to produce on-brand marketing copy: brand-voice prompting, brief→draft→human-QA loop, variation generation at scale, failure-mode guardrails, and FTC compliance. TRIGGER: use AI/LLMs for marketing copy while keeping brand voice; structure a copy brief for AI; QA or review AI-generated copy; human-in-the-loop copy workflow; generate headline/subject line variants for A/B testing; brand-voice drift / hallucinated stats / AI sameness; FTC rules for AI-generated ads (practical rules); publishing checklist for AI copy. SKIP: authoring a brand voice guide (no dedicated skill — use writing-expert); AI-persona/chatbot writing → content-and-marketing-writing; LLM prompt-engineering → ai-mcp-sdk-prompting; removing AI-voice tells → kill-the-ai-ism; prose craft frameworks (AIDA, PAS) → direct-response-and-sales-letter-copywriting; getting AI content cited by LLMs → generative-engine-optimization; FTC statute/fake-review rule detail → venture-marketing-strategy-local-seo.

---

# AI-Assisted Copywriting Workflow

> Practical. LLMs speed first drafts + variation; don't replace brand strategy, fact-check, or human judgment. Workflow keeps AI in lane.

## Why this matters (and what goes wrong)

AI copy fails predictably. 2026 audit: 15–25% numeric claims in AI content fabricated — fake reports (Digital Applied, 2026). Skip QA → hallucinated stats, stale product names (e.g., "Google AdWords" instead of "Google Ads"), slow voice homogenization only visible when catalog large enough to compare.

Fix = system, not better prompt. This skill cover that system.

---

## Phase 1 — Brand-voice prompting

### What to put in the prompt

Brand-voice prompt need five elements for reliable on-brand copy (Atom Writer, 2026):

| Element | What to specify |
|---|---|
| Persona | Who AI writes as ("You are a senior copywriter for [Brand], writing for [audience]") |
| Tone dimensions | Rate attributes: formal↔casual, dry↔warm, direct↔consultative |
| Vocabulary rules | Words to use; banned phrases; competitor terms to avoid |
| Structural rules | Sentence length range, paragraph format, heading style |
| Examples (few-shot) | 1–2 on-brand outputs; optionally 1 explicit off-brand negative |

3–5 before/after examples → 25–35% better style adherence vs rules-only (Atom Writer, 2026). Before: generic / after: on-brand.

### Awareness-stage-aware prompts

Match instructions to buyer stage:

- **Top of funnel (unaware/problem-aware):** Curiosity, relevance, category framing. No feature lists.
- **Middle of funnel (solution-aware):** Differentiation, proof points, comparison framing.
- **Bottom of funnel (product-aware):** Specificity, urgency, risk reduction, social proof.

Include awareness stage as explicit param in every copy prompt.

### Instruction drift

Output grows → model de-weights voice instructions (attention dilution). Counter:
1. Put critical rules at **start and end** of prompt (primacy + recency).
2. Break long-form into sections, not whole piece at once.
3. Add self-check: "Before finalizing, review whether this draft violates any of the banned phrases or tone rules above."

### System-level vs. user-level prompts

- **System prompt / custom GPT / Claude Project:** Up to ~2,000 words. Embed full brand-voice spec as reusable context.
- **Per-task user prompt:** Under ~500 words. Task-specific params (audience, channel, goal, format).

Keep brand spec separate from task brief — prevents spec crowding out task context.

---

## Phase 2 — The copy-brief → draft → human-QA loop

### The brief (human-authored, not AI-authored)

Copy brief for AI prompt include:

```
Product / offer: [what is being sold or promoted]
Target audience: [specific segment, job title, or persona]
Awareness stage: [unaware / problem-aware / solution-aware / product-aware]
Goal / desired action: [click, sign up, reply, purchase]
Key proof point or claim: [one verified, sourceable fact or differentiator]
Tone notes for this piece: [any channel or campaign-specific overrides]
Format: [length, structure, CTAs, character limits if relevant]
Forbidden territory: [competitors to avoid naming, claims we can't make, banned phrases]
```

Human write brief. AI generate draft. Division preserve strategic intent.

### Draft generation

Long-form: generate shorter chunks. Intro first, review, then section by section. Limits context loss, keeps voice consistent (Single Grain, 2025).

Short-form (subject lines, ads, social): request 5–10 variants per call. Structurally distinct — not paraphrases of same angle.

### Human-QA: the mandatory review passes

**Pass 1 — Factual verification (highest severity)**

Every stat, study, quote, spec, date, person ref must trace to primary source. No URL/doc → remove or reword. No shortcut. AI presents fabricated data with same confidence as verified (Prose Media, 2026; Agency Pro, 2026).

*Flag:* "X% say Y" with no named survey; stats from "a recent study"; stale data framed as current; non-existent integrations or features.

**Pass 2 — Brand-voice fidelity**

Read for voice, not just tone. Voice = what brand emphasizes, avoids, how it frames proof, what makes sales team cringe (Prose Media, 2026).

*Vocal-drift test:* Would stated author write this? No → rewrite, don't just soften.

*Structural tells:* Uniform 3–5 sentence paragraphs (AI default); 15–25 word sentences clustering; section intros summarizing what's coming.

**Pass 3 — AI-ism removal**

Scan for: *delve, leverage, tapestry, multifaceted, harness, streamline, furthermore, moreover, "it is important to note," "in today's landscape,"* em-dash overuse, over-capitalized abstracts. (See: kill-the-ai-ism for full sweep workflow.)

**Pass 4 — Strategic alignment**

Draft match brief's audience, awareness stage, desired action? Can be factually clean + on-voice but still wrong message for buyer.

**De-genericizing: human's highest-value contribution**

AI can't add what it doesn't have. Human job = injection, not correction:
- Personal anecdotes or customer stories
- Specificity about product not in public training data
- Emotional resonance calibrated to actual customer language (cross-ref: conversion-copywriting-and-voice-of-customer for VOC mining)
- POV distinguishing brand from every other vendor on topic

### QA checklist (printable)

```
FACTUAL
[ ] Every stat/number has a named, linkable source
[ ] Every person/company/study named actually exists
[ ] No "recent" data older than 3 years
[ ] Product features/specs match current product reality

VOICE
[ ] Reads like the brand, not like a template
[ ] No banned phrases
[ ] Tone matches awareness stage
[ ] No competitor brand terminology

AI ARTIFACTS
[ ] No AI-vocabulary tells (delve, leverage, etc.)
[ ] Sentence/paragraph length varies (not AI-uniform)
[ ] No summary-of-what's-about-to-be-said intros

STRATEGIC
[ ] Matches brief: audience, stage, goal
[ ] CTA fits the page intent
[ ] No claims broader than what we can support
[ ] Human-added specificity: story, proof, POV
```

---

## Phase 3 — Generating and testing variations at scale

### When variation generation makes sense

AI clearest ROI in copy = **volume for testing**: structurally distinct headline/subject line/CTA variants that'd take human hours (Pareto Performance, 2025; Single Grain, 2025).

Only valuable with traffic + measurement infrastructure. Oracle Marketing Cloud (2025): min ~1M active subscribers before ML-driven subject line optimization — below that, tests under-power.

### Structuring the variant generation prompt

```
Generate [N] structurally distinct [subject lines / headlines / CTAs] for:
- Product: [what]
- Audience: [who]
- Goal: [desired action]
- Dimension to vary: [e.g., benefit-led vs urgency-led vs social-proof-led]
- Character limit: [N]
- Brand voice constraints: [key rules]
- Forbidden: [competitor names, banned phrases, overclaims]

Label each variant with its structural angle (e.g., "value-first", "fear-of-missing-out", "specificity").
```

Grouping by structural angle → more genuinely distinct experiments than "10 different versions."

### Pre-flight screening before any variant goes live

Run brand-safety + compliance checks before deploy. Regulated industries: add legal/compliance gate. Never deploy unreviewed variants — volume makes off-brand/overclaiming easy to miss.

### Feeding results back

Test results only valuable if captured in voice spec. Angle consistently outperforms → document why, update few-shot examples in brand-voice prompt, note what to stop generating.

---

## Failure modes and guardrails

| Failure mode | Signal | Guardrail |
|---|---|---|
| Hallucinated statistics | Claims with no traceable source | Mandatory source-check on every numeric claim before publish |
| Brand-voice drift | Drafts sound like template, not brand | Per-client voice file (250–500 words of do/don't); QA against it, not general taste |
| AI sameness / homogenization | Consecutive posts read identically; catalog loses identity | Rotate structural angles; editor reads 3 random pieces from different dates to spot voice collapse |
| Over-claiming | Claims broader than product can support | Brief must specify what claims NOT permitted; pass 4 QA item |
| AI-vocabulary tells | "Delve into," em-dash abuse, uniform paragraph structure | Pass 3 artifact sweep; cross-ref kill-the-ai-ism |
| Stale knowledge | "Recent" data from before model cutoff | Flag any "recent" modifier; require publication date + link on all cited data |
| Source misattribution | Real stat credited to wrong publication or year | Verify independently even when source named — models confabulate publication details |

---

## FTC compliance and disclosure (2024–2026 — flag: recency)

**Status June 2026:** Actively evolving. Verify current FTC guidance before compliance decisions.

### What's settled

- **Fake reviews illegal.** FTC Final Rule on Consumer Reviews and Testimonials (effective Oct 2024) explicitly bans AI-generated fake reviews + testimonials — content for people who didn't use product (FTC, 2024). Operation AI Comply (Sep 2024) enforcement actions against AI fake-review generators.
- **Deceptive AI claims covered by existing law.** FTC Act unfair/deceptive prohibition applies regardless of AI use. No exemption (FTC, 2024).
- **Endorsement Guides apply to AI-assisted content.** Revised Guides (2023) require material connection disclosure. AI producing independent-consumer-voice look → disclosure risk.

### Practical rules

1. **Never fabricate testimonials.** Don't use AI to generate first-person reviews/quotes/endorsements for people who didn't provide them. Clearest legal line.
2. **Don't claim AI capabilities product lacks.** Overclaiming AI features → FTC exposure independent of copy quality.
3. **Disclosure norms for AI copy unsettled.** No FTC rule requiring disclosure that copy was AI-drafted (as of mid-2026). May change. Track FTC, consult counsel for regulated industries.
4. **Validate every factual claim independently.** False AI claims → deceptive advertising exposure under FTC Act. AI origin = no defense.

---

## Quick-reference decision table

| Task | Recommended workflow |
|---|---|
| One-off short-form draft (email, ad, social) | Brief → AI draft → Pass 1 + 2 + 3 → Human de-genericize → Publish |
| Long-form content (blog, whitepaper) | Brief → Section-by-section AI drafts → Pass 1 + 2 + 3 + 4 per section → Human inject specificity → Full-doc voice read → Publish |
| Headline / subject line variant generation | Brief → AI batch (5–10 structurally labeled variants) → Brand-safety screen → Pre-flight review → A/B test → Feed learnings back to voice spec |
| Brand-voice prompting setup | Build 5-element prompt (persona, tone dimensions, vocab rules, structure rules, 3–5 few-shot examples) → Test 10 outputs → Calibrate until ~70% need minimal editing |
| Reviewing AI copy you didn't brief | Run full QA checklist (4 passes) treating content as suspect until verified |

---

## Sources

1. **Atom Writer** — "Creating a Brand Voice Prompt for AI: Complete Template" (Feb 2026) — https://atomwriter.com/blog/brand-voice-ai-prompt-template/ — Five-element prompt structure; primacy/recency effects; few-shot adherence data.

2. **Atom Writer** — "The AI Writing Checklist: Quality Control Workflow" (May 2026) — https://www.atomwriter.com/blog/ai-writing-checklist-quality-control/ — AI vocabulary scan list; voice failure modes; hallucination as highest-severity risk.

3. **Digital Applied** — "AI Content Pipeline Anti-Patterns: Quality Failure Modes" (May 2026) — https://www.digitalapplied.com/blog/ai-content-pipeline-anti-patterns-quality-failure-modes-2026 — 15–25% hallucinated stat finding; voice collapse; stat fabrication as highest-severity anti-pattern.

4. **Prose Media** — "AI content QA checklist: how to catch errors before they embarrass your brand" (Mar 2026) — https://www.prosemedia.com/blog/ai-content-qa-checklist-before-publish — Brand voice vs. copyediting distinction; strategic drift.

5. **Agency Pro** — "AI Deliverables Quality Control: How Agencies Maintain Standards" (Apr 2026) — https://agencypro.app/blog/ai-deliverables-quality-control — Failure mode table; per-client voice file; source verification checklist.

6. **Single Grain / Eric Siu** — "LLM Brand Voice: How LLMs Interpret Tone in Content" (Dec 2025) — https://www.singlegrain.com/branding-2/how-llms-interpret-brand-tone-and-voice/ — RAG-augmented voice; approval tier design; section-by-section generation.

7. **Single Grain / Eric Siu** — "AI Ad Copy Testing at Scale Without Violating Brand Voice" (Dec 2025) — https://www.singlegrain.com/artificial-intelligence/ai-powered-ad-copy-testing-at-scale-without-violating-brand-voice/ — Hypothesis-to-test framework; brand-constraint translation; structured variant generation.

8. **The Brand Algorithm** — "How to Train AI to Write in Your Brand Voice (2026)" (May 2026) — https://www.the-brand-algorithm.com/guide-to-training-ai-for-brand-voice/ — 5-step human+AI workflow; short-form vs. long-form corpus sizing; 70% target for minimal-edit first drafts.

9. **Yotpo** — "AI Audit Checklist: Validate Content & Brand Voice" (Apr 2026) — https://www.yotpo.com/blog/ai-audit-checklist/ — Model homogenization / model collapse; Vectara Hallucination Leaderboard; semantic drift.

10. **Pareto Performance** — "Master AI for Paid Ads: A/B Testing at Massive Scale" (Sep 2025) — https://pareto-performance.com/mastering-ai-for-paid-ad-copy-a-b-testing-at-an-unprecedented-scale/ — Test structure types; metadata-tagged variant generation; traffic allocation.

11. **Oracle Marketing Cloud / Chad White** — "Using AI Subject Line & Copywriting Tools Successfully" (Jul 2025) — https://blogs.oracle.com/marketingcloud/using-ai-subject-line-and-copywriting-tools-successfully — Minimum audience size for ML testing; human control baselines still outperforming.

12. **FTC** — "Operation AI Comply" press release (Sep 2024) — https://www.ftc.gov/news-events/news/press-releases/2024/09/ftc-announces-crackdown-deceptive-ai-claims-schemes — Five enforcement actions; fake review tools; "no AI exemption from the laws on the books."

13. **FTC** — Final Rule on Consumer Reviews and Testimonials (Aug 2024) — https://www.govinfo.gov/content/pkg/FR-2024-08-22/pdf/2024-18519.pdf — Explicit prohibition on AI-generated fake reviews; prohibition on review hijacking.

14. **FTC** — Revised Endorsement Guides (2023) — https://www.ftc.gov/system/files/ftc_gov/pdf/p204500_endorsement_guides_in_2023.pdf — Updated guides covering social media and influencer endorsements; material connection disclosure.

15. **Search Engine Land / Andrew Holland** — "How to QA AI-generated content + free QA workflow checklist" (Mar 2026) — https://searchengineland.com/guide/qa-workflow-for-ai-generate-content — Automation vs. human review boundary; hallucination examples (product names, stat dates).