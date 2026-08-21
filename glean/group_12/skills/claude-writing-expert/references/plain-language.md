<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `plain-language` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: plain-language
description: "Plain language craft for customer-facing prose, legal-adjacent docs, regulated writing, and accessibility-driven simplification. Covers the US Federal Plain Writing Act (2010), PlainLanguage.gov / PLAIN guidelines, reading-grade targeting (Flesch-Kincaid, Gunning Fog), common-word substitution, sentence-length discipline, and jargon translation. TRIGGER: 'make this plain English', 'simplify for customers', 'reading-grade level', 'Flesch-Kincaid score', 'plain language', 'remove jargon for non-technical audience', 'rewrite at grade 8', 'plain writing act compliance'. SKIP: technical writing for engineers with no accessibility/regulatory driver (use technical-writing-craft); general business writing for internal or executive audiences without a plain-language compliance requirement (use writing-expert); contracts or legal instruments where precise legal vocabulary is mandatory."
origin: local
version: "1.2.0"
updated: "2026-05-29"
keywords:
  - plain language
  - plain English
  - readability
  - Flesch-Kincaid
  - Gunning Fog
  - reading grade level
  - jargon removal
  - Plain Writing Act
  - consumer communications
  - regulated disclosures
  - customer-facing prose
tags:
  - writing
  - plain-language
  - readability
  - compliance
  - accessibility
whenToUse:
  - User asks to make text "plain English" or simplify it for customers or a general audience
  - User needs a reading-grade-level score or check (Flesch-Kincaid, Gunning Fog)
  - User wants to remove jargon for a non-technical or non-specialist audience
  - Document must comply with the Plain Writing Act or a plain-English regulatory requirement
  - User is drafting consent forms, terms of service summaries, patient-facing health docs, or regulated disclosures
  - User's draft is flagged as "too complex" or "too legalistic" by a stakeholder or reviewer
  - User asks about substituting complex words with simpler ones
  - User wants to target a specific reading grade level (e.g., "write at 8th grade")
whenNotToUse:
  - Writing for a peer technical audience without an accessibility or regulatory requirement — use technical-writing-craft
  - General business, executive, or internal prose without a plain-language compliance driver — use writing-expert
  - Legal instruments or contracts where precise terms of art are legally required
  - Academic writing where discipline-specific norms govern style
related_skills:
  - writing-expert
  - technical-writing-craft
  - document-critique
  - rhetorical-frameworks-deep
---

# Plain Language

Reference for customer-facing prose, regulated writing, legal-adjacent documents, and multi-audience communications where clarity is the primary goal.

## When to use this skill

- User asks to "make this plain English" or "simplify this for customers"
- User needs a reading-grade-level check or score (Flesch-Kincaid, Gunning Fog)
- User wants to remove jargon for a non-technical audience
- Document must comply with the Plain Writing Act (federal agency comms, government-facing materials)
- User is drafting consent forms, terms of service summaries, patient-facing health docs, public-sector communications, or regulated disclosures
- User's draft is flagged as "too complex" or "too legalistic" by a stakeholder

## When NOT to use this skill

- Writing for a peer technical audience where precision vocabulary is expected (use `technical-writing-craft`)
- Pure business prose with no regulatory or accessibility driver (use `writing-expert`)
- Legal or medical contexts where intentional precision requires technical vocabulary — plain language rules may underspecify the meaning
- Academic writing, where discipline-specific norms supersede readability targets

---

## Quick-start workflow

**If the user has not provided text to edit:** ask once — "Please share the text you'd like me to simplify." Do not proceed until text is provided.

**Before editing anything:** check §9 (When plain language is the wrong tool). If the document contains precision-critical clauses (legal terms of art, clinical terminology, regulatory definitions), flag those passages first and explain why they should be glossed rather than replaced.

**Then follow §12** (Workflow checklist) in order. Deliver output per §11 (Output format).

---

## 1. Legal and regulatory foundation

**Plain Writing Act of 2010 (Pub. L. 111-274; 5 U.S.C. § 301 note)**

- Requires US federal executive agencies to write new or substantially revised covered documents in plain language.
- "Covered documents" = documents necessary for obtaining federal benefits or services, filing taxes, complying with federal requirements, or understanding federal policy (and published after July 1, 2011).
- Does not apply directly to private companies, but is the benchmark standard courts, regulators, and consumer-protection agencies reference when evaluating disclosure clarity.
- SEC "Plain English" rule (Rule 421(d), 17 C.F.R. § 230.421) independently mandates plain English in prospectus cover pages — active voice, short sentences, concrete terms, no legal/business jargon.
- CFPB, FDA patient labeling, and CMS Medicare communications have adopted equivalent plain-language requirements.

**PLAIN — Plain Language Action and Information Network**

- Federal interagency group that publishes `plainlanguage.gov`; the authoritative US government guidance.
- Core checklist: use "you" and "we"; active voice; short sentences; common words; lists for enumerable items; helpful headings; no unnecessary jargon.

---

## 2. UK plain English movement

Sir Ernest Gowers wrote *Plain Words* (1948, updated as *The Complete Plain Words*, 1954) for UK civil servants, arguing that bureaucratic prose obscures meaning and wastes reader time. Key Gowers principles:

- Prefer the familiar word over the unfamiliar.
- Prefer the concrete word over the abstract.
- Prefer the short word over the long.
- Prefer the Saxon word over the Latinate where both are natural.
- Never use a longer construction when a shorter one will do.

Plain English Campaign (UK, founded 1979) continues this tradition, awarding "Crystal Mark" accreditation to clear public documents. Its anti-patterns catalog overlaps heavily with PlainLanguage.gov.

---

## 3. Readability formulas

### Flesch Reading Ease (Flesch, 1948)

```
Score = 206.835 − (1.015 × ASL) − (84.6 × ASW)
```

Where:
- **ASL** = average sentence length (words per sentence)
- **ASW** = average syllables per word

Score interpretation:
| Score | Reading ease | Grade level | Typical audience |
|-------|--------------|-------------|------------------|
| 90–100 | Very easy | 5th grade | Basic consumer |
| 70–80 | Easy | 6th grade | Marketing copy, plain consumer |
| 60–70 | Standard | 7th–8th grade | General public target |
| 50–60 | Fairly difficult | 10th–12th grade | Technical docs acceptable |
| 30–50 | Difficult | College | Specialized professional |
| 0–30 | Very confusing | College grad | Expert-only |

**Targets:** general public → 60–70 (7th–8th grade). Marketing/consumer → 70+ (6th grade or easier). Technical documentation for professionals → 50–60 acceptable.

### Flesch-Kincaid Grade Level (Kincaid et al., 1975)

```
FKGL = (0.39 × ASL) + (11.8 × ASW) − 15.59
```

Returns a US school grade level (e.g., 8.2 = 8th grade, 2nd month). Microsoft Word computes this natively (Review → Document Stats). Google Docs requires a third-party add-on such as Readable or Grammarly.

**Targets:**
- General public communications: **Grade 8 or below**
- Consumer marketing: **Grade 6 or below**
- Technical docs for professional audience: **Grade 10–12 acceptable**
- Healthcare patient materials (per health literacy research): **Grade 6 or below**

### Gunning Fog Index (Gunning, 1952)

```
Fog = 0.4 × (ASL + (100 × complex_words / total_words))
```

Where **complex words** = words with 3+ syllables, excluding proper nouns, compound words, and common -ed/-es suffixes.

Fog score equals approximate years of formal education required. Target: **Fog ≤ 12** for general audience; **Fog ≤ 8** for consumer-facing docs.

### Practical calculation

You do not need to hand-compute these. Use:
- Microsoft Word: Review → Spelling & Grammar → Document Stats (shows FKGL + Flesch)
- Hemingway Editor (`hemingwayapp.com`): visual grade-level highlighting
- `readable.com` or `readability-score.com`: multi-formula batch scoring
- Python `textstat` library: `textstat.flesch_kincaid_grade(text)`

---

## 4. Core plain-language rules (PlainLanguage.gov)

### 4.1 Address the reader directly
Use "you" for the reader, "we" for your organization. Avoid impersonal constructions ("the applicant must" → "you must"; "it is required that" → "you must").

### 4.2 Active voice
Place the actor before the action. "The committee approved the request" not "The request was approved by the committee." Passive is acceptable when the actor is unknown or irrelevant — but keep it rare.

### 4.3 Short sentences
Target average sentence length of 15–20 words. No single sentence should exceed 40 words. The 95th-percentile sentence in a document should stay under 25 words. Break compound sentences at conjunctions; break serial-clause sentences into lists.

### 4.4 Common words
Prefer the everyday word. See §5 for the substitution table. When a technical term is unavoidable, define it on first use in plain parenthetical: "latency (the delay before data transfers begin)".

### 4.5 Lists for lists
When you have three or more parallel items, use a bulleted or numbered list rather than a run-on sentence. Use numbers when sequence matters; bullets otherwise. Keep list items grammatically parallel.

### 4.6 Informative headings
Use question-form or descriptive headings ("What documents do you need?" rather than "Documentation"). Headers should allow a scanner to find their section without reading surrounding text.

### 4.7 One idea per paragraph
Open each paragraph with the main point. Do not bury the conclusion. Paragraph length: 3–5 sentences for general audiences; rarely exceed 7.

---

## 5. Common-word substitution table

| Avoid | Use instead |
|-------|-------------|
| utilize | use |
| in the event that | if |
| prior to | before |
| subsequent to | after |
| terminate | end / stop |
| commence | start / begin |
| facilitate | help / support |
| endeavor | try |
| obtain | get |
| approximately | about |
| at this point in time | now |
| due to the fact that | because |
| in order to | to |
| with regard to | about / on |
| in accordance with | under / per |
| pursuant to | under / following |
| demonstrate | show |
| sufficient | enough |
| additional | more / extra |
| component | part |
| methodology | method |
| impacted | affected |
| leverage (non-financial) | use |
| implement | put in place / carry out |
| initiate | start |
| finalize | finish |
| ascertain | find out |
| be cognizant of | know |
| provide assistance to | help |
| on a regular basis | regularly |
| in close proximity to | near |

---

## 6. Concrete vs abstract nouns

Abstract nouns (nominalization) are verbs or adjectives turned into nouns. They add syllables, distance the actor, and inflate grade level.

| Abstract (nominalization) | Concrete (verb or adjective form) |
|---------------------------|-----------------------------------|
| provide an explanation for | explain |
| conduct an investigation of | investigate |
| make a determination | determine / decide |
| give consideration to | consider |
| have an impact on | affect / change |
| perform an analysis of | analyze |
| reach a conclusion | conclude |
| be in violation of | violate |

Rule: when a noun ends in -tion, -ment, -ance, -ence, -ity, -ism, or -al, check whether the underlying verb or adjective is more direct.

---

## 7. Jargon translation rules

1. **Never assume vocabulary.** If a term appears in a dictionary but not in everyday conversation, treat it as jargon for a general audience.
2. **Define on first use.** Parenthetical plain-English gloss after first mention: "indemnification (protection against legal claims)".
3. **Consider a footer glossary** for documents with 5+ technical terms that cannot be avoided.
4. **Acronyms:** spell out every acronym on first use, even familiar ones. Different audiences know different acronyms.
5. **Industry-specific terms:** if the term has a plain synonym, use the synonym. If not (e.g., "HIPAA"), spell it out and describe what it governs.
6. **Test jargon against the "stranger on the street" standard:** would someone outside your industry understand this word without prompting?

---

## 8. Anti-patterns

| Anti-pattern | Description | Fix |
|--------------|-------------|-----|
| Throat-clearing | Opening sentences that delay the main point ("It is important to note that...") | Delete or move main point to sentence 1 |
| Doubled words | Legal doublets that add nothing ("aid and assist", "null and void", "terms and conditions") | Keep one word |
| Latinate preference | Choosing long Latinate words over short Anglo-Saxon ones ("utilize" vs "use") | See §5 |
| Zombie nouns | Nominalized verbs that hide the actor (see §6) | Restore the verb |
| Passive stacking | Multiple passive clauses in sequence, obscuring all actors | Rewrite with active subject |
| Hedging forests | Excessive qualifiers ("it may be possible that", "in some cases") that erode confidence | Remove or rephrase as specific condition |
| Wall-of-text paragraphs | >7-sentence paragraphs that make the reader work to find the point | Break at logical sub-topics; use lists |
| Circular definitions | "Termination means the act of terminating" | Give a concrete plain-English meaning |

---

## 9. When plain language is the wrong tool

The "When NOT to use this skill" section above governs skill routing. This section covers in-document judgment: when you are already editing a document and encounter content that should stay complex.

Do not simplify these passages:

- **Precision-critical legal clauses** — terms of art with established legal meaning ("indemnification", "force majeure", "fiduciary duty") should be defined, not replaced. Substituting a plain synonym may change the legal meaning or create ambiguity. Action: add a parenthetical gloss, not a word swap.
- **Clinical and diagnostic terminology** — "myocardial infarction" has a clinical precision "heart attack" does not. In patient-facing materials, keep the clinical term and add a plain explanation: "myocardial infarction (heart attack)". Do not drop the technical term.
- **Peer-to-peer technical specifications** — when the document audience is confirmed as domain experts, precision vocabulary is correct. Oversimplification degrades accuracy.
- **Regulatory definitions** — defined terms in statutes and regulations must match the statutory definition exactly. Do not rephrase; add a plain-English explanation in a separate sentence or sidebar.

---

## 10. Audience-appropriate calibration

"Plain" is not a single register. Calibrate by audience:

| Audience | Target FKGL | Sentence target | Vocabulary guidance |
|----------|-------------|-----------------|---------------------|
| General public / consumers | ≤ 8 | ≤ 20 avg | Everyday words; no acronyms unreduced |
| Healthcare patients | ≤ 6 | ≤ 15 avg | Define every clinical term |
| Marketing / acquisition copy | ≤ 6 | ≤ 15 avg | Conversational; second-person throughout |
| Board members / executives | ≤ 10 | ≤ 22 avg | Business vocab OK; no tech jargon |
| Regulated disclosures (SEC, CFPB) | ≤ 8 | ≤ 20 avg | Active voice; plain-English definitions required |
| Technical docs for professionals | 10–12 | ≤ 25 avg | Domain vocabulary expected; define only on first use |

Plain-for-engineers uses short sentences and active voice but retains technical precision. Plain-for-board-members avoids engineering jargon but uses finance vocabulary. Do not apply a single universal grade target across all audiences.

---

## 11. Output format and clarifying questions

### When the audience is unknown

If the user has not identified an audience, ask one question before editing:

> "Who is the primary reader — general public, customers, patients, executives, or a regulated/compliance context?"

Do not ask about format, tone, or length until the audience is confirmed. Audience determines the grade-level target (§10), which drives all other decisions.

### Deliverable shape

When returning a plain-language edit, provide:

1. **Rewritten text** — the full edited passage as plain prose (not a code block unless the original was code). Ready to copy-paste directly.
2. **Before/after readability scores** — FKGL or Fog Index for the original and revised text. Calculate if a tool is available; estimate from sentence length and word complexity if not, and label the estimate clearly: "(estimated: ~FKGL 9 → ~FKGL 7)". If unable to estimate, state "score unavailable — verify with Hemingway Editor or Word's Document Stats."
3. **Change summary** — a brief list of what changed: sentences split, words swapped, passives converted, nominalizations restored to verbs. Cap at 5 bullets.

If the document is long (>500 words), offer to work section by section rather than rewriting the whole thing at once.

### Worked example

**Original (FKGL ~13):**
> "In the event that the aforementioned documentation requirements are not fulfilled prior to the commencement of the review period, the applicant's submission will be deemed insufficient for further consideration."

**Plain-language rewrite (FKGL ~7):**
> "If you do not submit the required documents before the review period starts, we will not consider your application."

**Change summary:**
- "In the event that" → "If"
- "aforementioned documentation requirements are not fulfilled" → "you do not submit the required documents"
- "prior to the commencement of" → "before"
- Passive ("will be deemed") → active ("we will not consider")
- 37-word sentence → 18 words

---

## 12. Workflow checklist

When asked to apply plain language to a document:

1. **Score first.** Run a readability score before editing so you have a baseline (FKGL and Fog Index).
2. **Identify the audience.** Confirm reading-grade target from §10 before editing.
3. **Attack sentence length.** Find sentences over 30 words. Break each one.
4. **Swap jargon.** Work through §5 and §7. Flag unavoidable terms for definition.
5. **Eliminate nominalizations.** Scan for -tion/-ment/-ance suffixes (§6).
6. **Convert passive voice.** Identify actor-action pairs; front-load the actor.
7. **Test lists.** Confirm all enumerable items (3+) are bulleted or numbered.
8. **Re-score.** Run readability again. Confirm grade target is met.
9. **Spot-check for precision loss.** Review defined terms and confirm no meaning changed.

---

## 13. Hemingway-style short-sentence discipline

Ernest Hemingway's prose taught a generation of journalists, copywriters, and
plain-language editors a single rule: short sentences, simple words, active
voice, almost no adverbs. The Hemingway Editor app
(`https://hemingwayapp.com`) codifies the same rules as automated scoring:
flag sentences that are hard to read, sentences that are very hard to read,
adverbs, passive voice, and complex-word substitutions.

### The rule

| Target | Rule |
|--------|------|
| Average sentence length | ≤ 14 words (Hemingway app default), aligned with FKGL grade 6–8 |
| Longest sentence | ≤ 25 words; never exceed 35 |
| Adverbs | Limit to about 1 per 100 words (`-ly` ending and intensifiers like *very*, *really*, *quite*) |
| Passive voice | < 20% of sentences |
| Complex words | Substitute when a shorter word does the same work (see §5) |

### Worked example

Long-sentence draft (FKGL ~14):
> "The patient should be carefully monitored by the nursing staff over the
> course of the next twenty-four hours for any signs of complications that
> might potentially develop following the procedure, and any unusual symptoms
> should be promptly reported to the attending physician."

Hemingway pass (FKGL ~7):
> "Nurses will watch the patient for the next 24 hours. Watch for any signs
> of trouble after the procedure. Call the doctor right away if anything
> looks wrong."

Changes:
- 41-word sentence broken into three (12, 11, 12 words).
- Passive "should be monitored" → active "Nurses will watch."
- "Carefully," "potentially," "promptly" — three adverbs removed; the verbs
  carry the meaning.
- "Complications that might develop" → "signs of trouble."
- Latinate vocabulary ("attending physician") → plain ("doctor").

### When to break the rule

- Legal terms of art and clinical diagnostic terminology stay precise even
  when long (see §9). Hemingway pressure applied to a legal definition can
  break the meaning.
- A controlled long sentence carrying parallel structure (Williams ch. 6 —
  visible architecture) is fine even at 35 words — readers parse parallels
  faster than unstructured prose of the same length.
- Marketing copy can drop below 8 words to land a rhythm; do not feel
  obligated to hit 14 every time.

### References

- Hemingway, E. *A Farewell to Arms* (1929); *The Sun Also Rises* (1926).
- Hemingway Editor. `https://hemingwayapp.com` — automated rule application.
- Zinsser, W. *On Writing Well* (1976) ch. 2 — simplicity.

---

## 14. The noun-stack problem

A "noun stack" (also called a noun pile-up) is a chain of three or more nouns
in a row, modifying a final noun, with no prepositions to indicate the
relationships. Common in regulatory and technical prose because nouns feel
formal; the cost is forced re-parsing by the reader.

### The rule

Limit noun stacks to two consecutive nouns. When you find three or more,
break the stack by:

1. Inserting prepositions (of, for, in, to, by).
2. Pulling out a verb hiding inside one of the nouns.
3. Splitting the stack into two sentences.

### Worked examples

Stacked (regulatory prose):
> "Health insurance premium subsidy eligibility determination notice."

Five nouns in a row before "notice." A reader has to mentally parenthesize.

Unstacked:
> "Notice of the determination of your eligibility for a subsidy on your
> health insurance premium."

Or, plainer still:
> "We have decided whether you qualify for help paying your health insurance."

Stacked (technical prose):
> "Database replication lag monitoring alert threshold configuration."

Unstacked:
> "How to configure the alert that fires when database replication lag
> exceeds the threshold."

### Why it matters for plain language

Noun stacks add cognitive load disproportionate to their length: a 6-word
stack can be slower to read than a 12-word unstacked sentence. They also
inflate the Gunning Fog index because most components are 3+ syllables.

### When to break the rule

Some compound terms are canonical and unstacking produces stilted prose:
"two-factor authentication setup," "connection pool exhaustion." The rule
targets ad-hoc stacks, not standard terms of art.

### References

- Pinker, S. *The Sense of Style* (2014), ch. 4 — "noun piles."
- PlainLanguage.gov — "Avoid noun strings."
  `https://www.plainlanguage.gov/guidelines/words/avoid-noun-strings/`

---

## Sources

- Plain Writing Act of 2010, Pub. L. 111-274 (5 U.S.C. § 301 note)
- PlainLanguage.gov — Federal Plain Language Guidelines, Revision 1 (May 2011)
- Flesch, R. (1948). "A new readability yardstick." *Journal of Applied Psychology*, 32(3), 221–233.
- Kincaid, J.P., Fishburne, R.P., Rogers, R.L., & Chissom, B.S. (1975). *Derivation of New Readability Formulas for Navy Enlisted Personnel* (Research Branch Report 8-75). Naval Technical Training Command.
- Gunning, R. (1952). *The Technique of Clear Writing*. McGraw-Hill.
- Gowers, E. (1954). *The Complete Plain Words*. HMSO (updated edition, Penguin, 2014).
- SEC Rule 421(d), 17 C.F.R. § 230.421 — Plain English Requirements for Prospectuses.

---

## SMOG and the FOG companion (extending the existing reading-level coverage)

This skill already covers Flesch Reading Ease, Flesch-Kincaid Grade Level, and Gunning Fog in detail. SMOG (Simple Measure of Gobbledygook) is the third widely cited reading-level formula and the one most often used in healthcare, legal, and government plain-language compliance. It is added here as a companion, not a replacement.

### SMOG (McLaughlin, 1969)

**Rule.** Harry McLaughlin's SMOG formula was designed specifically as a *quick, conservative* readability estimate aimed at *complete comprehension* — McLaughlin's premise was that the FOG and Flesch-Kincaid formulas predict the grade level where a reader *partially* understands the text, while SMOG predicts the grade level where a reader *fully* understands it. SMOG is therefore the formula of choice in healthcare patient education (used by CDC, NIH, NCI) and in many state-government and insurance plain-language compliance audits.

**Formula.**

```
SMOG grade = 1.0430 * sqrt(30 * polysyllables / total_sentences) + 3.1291
```

Where *polysyllables* are words of 3+ syllables, counted across exactly 30 sentences (10 at the start, 10 in the middle, 10 at the end). For texts under 30 sentences, McLaughlin published a conversion table; most modern tools handle the small-sample case automatically.

**Quick-estimate version** (the one in widespread clinical use):

```
SMOG ≈ sqrt(polysyllable count in 30 sentences) + 3
```

**Worked example.**

Sample: 30 sentences from a patient-discharge instruction sheet. Polysyllable count: 49.

- SMOG (exact): 1.0430 × sqrt(30 × 49 / 30) + 3.1291 = 1.0430 × sqrt(49) + 3.1291 = 1.0430 × 7 + 3.1291 ≈ 10.43
- SMOG (quick): sqrt(49) + 3 = 10

Verdict: requires roughly a 10th-grade reading level for *full* comprehension. CDC clear-communication guidance targets grade 6–8 for general public health materials; this document needs simplification.

**How SMOG compares to FOG and Flesch-Kincaid.**

| Formula | Counts | Predicts |
|---|---|---|
| Flesch-Kincaid Grade Level | sentence length + syllables/word | Approximate US grade for *partial* comprehension |
| Gunning FOG | sentence length + % "hard words" (3+ syllables, excluding common suffixes) | Approximate US grade for *partial* comprehension |
| SMOG | polysyllabic word density across 30 sentences | Approximate US grade for *complete* comprehension |

SMOG typically scores 1–2 grades *higher* than Flesch-Kincaid on the same passage — not because the text is harder, but because SMOG is asking a stricter comprehension question.

**When to use which.**

- *General business and tech writing*: Flesch-Kincaid is fine. Most editors (Word, Hemingway) report it natively.
- *Government plain-language compliance*: PlainLanguage.gov accepts Flesch-Kincaid or Gunning FOG; both are referenced in agency style guides.
- *Health, legal, insurance, K-12 educational materials*: use SMOG. CDC, NIH/NCI, and the American Medical Association recommend SMOG for patient-facing materials specifically. Target grade 6–8 for general public; grade 5 for low-literacy populations.
- *Regulated financial disclosures* (SEC plain-English rule): no single formula is mandated, but the rule's effective practice combines Flesch Reading Ease (target ≥60) with sentence-length and active-voice rules.

### WCAG 3 draft thresholds (current as of 2026)

The W3C WCAG 3 working draft proposes explicit reading-level targets as a Level Bronze conformance item: *content should be readable at a secondary-education level (US grade 9 / lower secondary)* for general-audience pages, with Level Silver targeting *primary-education level* (US grade 6) for safety-critical content. These are draft thresholds and not yet normative, but they align with the existing CDC clear-communication targets and are already showing up in enterprise accessibility audits.

**References.**

- McLaughlin, G. H. "SMOG Grading: A New Readability Formula." *Journal of Reading* 12(8), 1969, pp. 639–646. https://www.semanticscholar.org/paper/SMOG-Grading-a-New-Readability-Formula-McLaughlin
- CDC — "Simply Put: A Guide for Creating Easy-to-Understand Materials." https://www.cdc.gov/healthliteracy/pdf/Simply_Put.pdf
- W3C — "WCAG 3 Working Draft." https://www.w3.org/TR/wcag-3.0/ (check the current draft for the active reading-level requirement)
- Python `textstat` library — implements `textstat.smog_index(text)`, `textstat.gunning_fog(text)`, `textstat.flesch_kincaid_grade(text)`.

---

## The Curse of Expertise — the plain-language angle

This skill already references the Curse of Knowledge through its jargon and noun-pile sections. The Curse of *Expertise* is the closely related second failure that plain-language writers hit specifically. The two are distinct.

- **Curse of Knowledge** (Pinker; Heath brothers): once you know a fact, you cannot remember not knowing it. The fix is *audience modeling* — imagine what the reader knows.
- **Curse of Expertise** (Willingham; Hinds 1999): with deep mastery, your reasoning becomes *automatic and chunked*. You skip steps not because you forgot the reader, but because *those steps are no longer conscious for you*. The fix is *step decomposition* — write the procedure for a novice and run it as a novice would, capturing every micro-decision the chunked version hid.

The plain-language failure mode is acute because the chunks are not jargon — they are *unstated prerequisites*. A patient instruction sheet that says "take one tablet daily with food" hides every micro-decision an expert makes silently: which meal counts as "with food," what to do about the missed dose, what counts as a side-effect worth reporting, what "daily" means when the dose is at 11 PM and the day rolls over. None of those gaps are jargon — they are *expert chunks*.

**Plain-language fix.**

1. *Write the procedure for the target reading grade* using the techniques already in this skill.
2. *Have a novice perform the task while you watch silently.* Every hesitation, every wrong turn, every clarifying question is a Curse-of-Expertise gap your draft did not address.
3. *Add a step or example for each gap.* Plain-language drafts almost always need *more* steps, not fewer. The instinct to compress for grade-level conflicts with the need to surface chunked steps; resolve the conflict in favor of clarity.

**Worked example — a discharge instruction.**

Expert-chunked: "Resume normal activity as tolerated."
Novice-tested rewrite: "After 48 hours, you can walk around the house and do light tasks like making meals. Wait 1 week before lifting anything over 10 pounds. Wait 2 weeks before driving. Call your doctor if any activity causes pain, dizziness, or shortness of breath."

The original is grade 8; the rewrite is *also* grade 8 by SMOG, but it has decomposed three expert chunks (timing, weight limits, warning signs).

**Diagnostic for plain-language drafts.** After scoring for reading grade, run the *novice task test*: hand the draft to someone at the target audience's level and ask them to perform or describe the procedure. Note every clarifying question. Each one is a chunk the expert author did not see.

**References.**

- Willingham, D. T. *Why Don't Students Like School?* Jossey-Bass, 2009.
- Hinds, P. J. "The curse of expertise: The effects of expertise and debiasing methods on prediction of novice performance." *Journal of Experimental Psychology: Applied* 5(2), 1999.
- Heath, C. & Heath, D. *Made to Stick*. Random House, 2007 (Curse of Knowledge framing).
- See also: `writing-expert` skill — *The Curse of Expertise (distinct from the Curse of Knowledge)* for general-purpose framing.
