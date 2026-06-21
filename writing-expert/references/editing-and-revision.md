<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `editing-and-revision` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: editing-and-revision
version: "1.2.0"
updated: "2026-05-29"
description: >
  Multi-pass editing and revision craft for tightening prose after drafting is complete.
  TRIGGER: "edit this", "revise this", "tighten this", "cut by 30%", "make this shorter",
  "kill my darlings", "verb-first pass", "active voice scan", "trim the fat", "sharpen this".
  SKIP: drafting from scratch (writing-expert); structured critique with findings
  (document-critique); rhetorical structure (rhetorical-frameworks-deep);
  simplifying for a non-expert audience (plain-language).
triggers:
  - "edit this"
  - "revise this"
  - "tighten this"
  - "cut by 30%"
  - "make this shorter"
  - "kill my darlings"
  - "verb-first pass"
  - "active voice scan"
  - "cut the fat"
  - "trim this"
  - "sharpen this"
skip:
  - drafting from scratch → use writing-expert
  - structured review with findings → use document-critique
  - rhetorical architecture → use rhetorical-frameworks-deep
  - simplifying for a general/non-expert audience → use plain-language
related:
  - writing-expert
  - document-critique
  - technical-writing-craft
  - plain-language
  - rhetorical-frameworks-deep
sources:
  - Joseph M. Williams, "Style: Lessons in Clarity and Grace" (12th ed.)
  - William Strunk Jr. & E.B. White, "The Elements of Style"
  - Stephen King, "On Writing: A Memoir of the Craft"
  - Steven Pinker, "The Sense of Style: The Thinking Person's Guide to Writing in the 21st Century"
---

# Editing and Revision

**When invoked:** the user submits existing prose (inline or as a file path). If no text is provided, ask once: "Paste the text to edit, or give me a file path." For text over ~2,000 words, ask for a file path rather than inline paste.

**What to return:** the revised text, complete and drop-in ready, followed by a brief pass log noting which passes were applied and what was changed. Do not return only commentary — always return the edited text.

**Success criteria:** the revised text is shorter or clearer than the original, no meaning has been altered, and every applied pass is documented so the author can review changes.

**Known failure mode:** applying passes mechanically to deliberate stylistic choices — intentional fragments, rhetorical repetition, controlled sentence fragments for effect. Before cutting, check whether a pattern appears once (stylistic choice) or throughout (habit to fix).

---

Editing is not drafting. It is not reviewing. Drafting produces raw material; reviewing renders judgment. Editing transforms existing prose into what it should have been — by restructuring, cutting, and resurfacing what's already there.

The anti-pattern that kills both activities: editing while you draft. You get neither good flow nor good compression. Finish the draft first. Then edit cold.

---

## The Multi-Pass Methodology

Work top-down. Fix the largest problems before polishing sentences that may disappear anyway.

**Pass 1 — Structural.** Does the argument hold? Is the order defensible? Does each section earn its place? Cut entire sections before touching sentences.

**Pass 2 — Paragraph.** Does each paragraph have one job? Does the first sentence of each paragraph, read in sequence, produce a coherent argument? Reorder or merge before editing individual sentences.

**Pass 3 — Sentence.** Apply the verb-first pass, active-voice scan, and cohesion check at this level. Fix rhythm and length here.

**Pass 4 — Word.** Nominalizations, hedges, redundant pairs, filler phrases. This is the last pass, not the first.

The 24-hour rule: if the stakes are high, let the draft cool overnight before Pass 3 onward. Distance reveals what familiarity hides.

The five named passes in the "Targeted Editing Passes" section map to this hierarchy: Verb-First and Active-Voice belong to Pass 3 (sentence-level); Cut-30%, Paragraph-Coherence, and Topic-Sentence-First are Pass 2 work (paragraph-level). Run them in that order, not in isolation.

---

## Williams' Ten Lessons (Applied)

### 1. Actions in Verbs, Not Nouns — Kill Nominalizations

Nominalizations bury your action inside a noun formed from a verb: *conduct an investigation* instead of *investigate*, *make a recommendation* instead of *recommend*, *give consideration to* instead of *consider*.

Spotting them: look for Latinate nouns ending in *-tion*, *-ment*, *-ance*, *-ence*, *-ity*, *-ness* that have living verbs inside them. Pull the verb out and rebuild the sentence around it.

Before: *The committee engaged in a discussion of the proposal.*
After: *The committee discussed the proposal.*

Test every sentence: what is actually happening here? Name that action as a verb.

### 2. Characters in Subjects

Readers understand sentences fastest when the grammatical subject is the person or thing doing the action. Passive constructions and nominalizations push characters off the stage.

Before: *It was determined by the review board that the application was deficient.*
After: *The review board found the application deficient.*

Ask: who is doing what? Put the doer in the subject slot.

### 3. Cohesion — Old Information Before New

Readers expect each sentence to open on ground already established and to close on something new. Violating this creates disorientation even when every sentence is technically correct.

Principle (Pinker calls it the "Given-New contract"): the beginning of a sentence anchors the reader to what they know; the end delivers new payload.

Check: can a reader finish your opening clause from the previous sentence? If they cannot, you may have scrambled the information order.

### 4. Emphasis — Put the Most Important Content at the End

The end of a sentence is its stress position. Readers unconsciously weight what lands there. Burying the key term in the middle of a long sentence wastes the position.

Before: *The discovery of the error, which had been overlooked for six months by the auditing team, was unexpected.*
After: *The auditing team had overlooked the error for six months — and the discovery shook us.*

If your sentence ends on a qualifier, a parenthetical, or filler, you have surrendered emphasis. Cut the tail or restructure.

### 5. Concision — Cut Metadiscourse, Hedges, Redundant Pairs

Three categories to eliminate on sight:

**Metadiscourse** — sentences about the text rather than about the subject: *It is important to note that*, *As I mentioned above*, *The purpose of this section is to*.

**Hedges** — epistemic padding that rarely reflects genuine uncertainty: *it could be argued that*, *it seems as if*, *in some sense*, *to a certain extent*.

**Redundant pairs** — *each and every*, *null and void*, *true and accurate*, *various and sundry*. One word does the job; the second adds nothing.

Strunk's rule applies: omit needless words. The test is not whether a word is defensible but whether the sentence is weaker without it.

### 6. Shape — Long Sentences Need Architecture

Long sentences are not inherently bad. Unmanaged long sentences are. A sentence that runs past 35 words must give readers a structural signal: a clear subject-verb core, a list with parallel grammar, or a subordination pattern with an obvious spine.

Diagnosis: find every sentence over two lines. Read it aloud. If you lose the thread before the period, break it or add a visible frame (colon, semicolon, numbered list).

### 7. Elegance — Rhythm, Balance, Parallelism

Parallelism is not a style preference; it is a comprehension tool. Readers expect grammatically matched items to be semantically matched. When you violate the expectation, readers waste effort figuring out the mismatch.

Before: *We need to improve morale, reduce attrition, and the onboarding process should be shortened.*
After: *We need to improve morale, reduce attrition, and shorten onboarding.*

Rhythm follows from cutting weak words and letting stressed syllables land where emphasis belongs. Reading aloud exposes clumsy cadence faster than any other technique.

### 8. Coherence — Global Argument Flow

Local sentence polish cannot rescue an incoherent argument. Before line-editing, map the skeleton: thesis, evidence, counterargument, resolution. If the skeleton is wrong, edit from there.

The paragraph-first-sentences test: extract only the opening sentence of each paragraph and read them in order. They should form a compressed but coherent version of your argument. If they do not — if a paragraph's first sentence is an example rather than a claim, if two paragraphs cover the same ground, if the sequence feels arbitrary — fix the structure before fixing the sentences.

### 9. Subordination — Correct Relative Clauses

Subordinate clauses focus and qualify claims. Misused, they obscure the main action. The test: is the content in the subordinate clause actually subordinate, or is it the main point dressed up as a modifier?

Before: *The proposal, which would restructure the department and realign three reporting lines, was tabled.*

If the restructuring is the story: *The proposal would restructure the department and realign three reporting lines — the board tabled it.*
If the tabling is the story: *The restructuring proposal was tabled.* (compress the relative clause to a noun phrase)

### 10. Ethics of Writing — Clarity as Honesty

Williams closes with something most style guides avoid: obscurity is not neutral. Bureaucratic prose, passive constructions that hide agency, nominalizations that diffuse accountability — these are choices, not accidents. *Mistakes were made* hides who made them. *The policy requires* hides who wrote the policy.

Clarity is a form of honesty. When you edit, ask: am I hiding something? Am I making this harder to understand than it needs to be? If the answer is yes, the fix is not a stylistic choice but an ethical one.

---

## Targeted Editing Passes

### The Verb-First Pass

Scan every sentence. Locate the main verb. If the verb is a form of *be* plus a nominalization (*is a reflection of*, *was a consideration in*), extract the hidden verb and rebuild. If the verb arrives after line two, restructure.

Go sentence by sentence. This pass is not skimmable.

### The Cut-30% Pass

Set a target: 30% of words on the cutting room floor. Not a suggestion — a constraint. The constraint forces prioritization. You cannot preserve everything, so you are forced to identify what actually matters.

Tactics: cut every sentence that restates a point already made. Cut every transition that explains what you are about to say rather than saying it. Cut the last sentence of introductions and the first sentence of conclusions — both are usually restatements.

King: "Kill your darlings, kill your darlings, even when it breaks your egocentric little scribbler's heart." The sentences you love most are often the ones that exist to display your cleverness rather than to serve the reader.

The Faulkner formulation: *In writing, you must kill all your darlings.* A darling is any sentence you would fight to keep. Ask why you would fight. If the answer is "because I wrote it well" rather than "because it is load-bearing," cut it.

### The Active-Voice Pass

Search the draft for every *was* and *were* followed by a past participle (*was written*, *were submitted*, *was determined*). For each instance, decide:

1. Is the actor unknown or genuinely irrelevant? Keep passive.
2. Is the actor known and relevant? Rewrite in active voice, naming the actor.
3. Is the passive hiding accountability? This is the ethics test from Lesson 10. Rewrite.

Not all passive is wrong. Passive is correct when the patient matters more than the agent: *The bridge was designed in 1887* is fine if the year matters and the designer does not.

### The Paragraph-Coherence Pass

Print or copy only the first sentence of each paragraph. Read those sentences as a standalone document. Ask:

- Does the sequence form an argument?
- Are any two sentences redundant?
- Does the sequence skip a necessary step?
- Does any sentence function as an example rather than a claim?

Restructure at the paragraph level before you do any sentence editing.

### The Topic-Sentence-First Pass

For every paragraph, ask: is the point in the first sentence, or is it buried in sentence three after setup and context? If buried, move it to the front. Setup belongs after the claim, not before it — unless you are writing a mystery.

### Final Self-Check

After applying all passes, re-read the revised text against the original. Confirm: no meaning was altered, no subject lost its referent, no argument step was deleted. If the revision changed something it should not have, restore the original passage.

---

## Editing AI-Generated Text

AI output arrives fluent but structurally flat — it tends toward balanced hedges, neutral verbs, and paragraphs that restate rather than advance. Apply the same passes, but expect more of these specific patterns:

- **Nominalization density is high.** LLMs favor *provide an explanation of* over *explain*, *make a determination* over *decide*. The verb-first pass is the highest-leverage move on AI text.
- **Emphasis is weak.** Sentences end on qualifiers (*in many cases*, *depending on the context*). Move the load-bearing term to the end.
- **Topic sentences are buried or absent.** Paragraphs open with context or transition rather than claim. The topic-sentence-first pass is essential.
- **The cut-30% target is easily met.** AI text typically contains 20–40% restatement. Apply it without guilt.

AI text rarely needs structural editing — the argument order is usually serviceable. Sentence and word passes are where the work is.

---

## The Floor Rule — When to Stop Cutting

Cutting has a floor. A sentence compressed past its meaning is not tight; it is broken. Stop when:

- A reader would need to re-read to recover the missing context.
- The sentence no longer names its subject (pronouns without clear antecedents).
- Parallel structure has been broken to save words.
- The tone has shifted from direct to terse in a way that reads as dismissive.

The test: read the cut sentence cold, 24 hours later. If you have to reconstruct what it meant, you cut too far. Restore the minimum needed for clarity.

---

## Reading Aloud as an Editing Tool

Every sentence that sounds wrong when spoken is wrong. Reading aloud catches:

- Run-on sentences (you run out of breath)
- Rhythm clashes (stressed syllables landing on weak words)
- Missing transitions (you stumble at a join)
- Repetition (you hear the same word twice in three lines)
- Empty throat-clearing (you feel the urge to skip it)

Read slowly. Read what is on the page, not what you intended to write.

---

## Additional Editing Passes

### The Verb-Tense Pass

In technical and business prose, three tenses each have a default job. Mixing
them without reason confuses readers about whether they are being told
history, behavior, or plan.

**The rule:**

| Tense | Default use |
|-------|-------------|
| **Past** | Specific historical events, completed incidents, post-mortem narrative, changelog entries. "We migrated the cluster on 2025-03-14." |
| **Present** | Current system behavior, API semantics, reference facts, anything still true. "The function returns a cursor." |
| **Future** | Plans, scheduled deprecations, commitments with a date. "The `/v1` endpoints will be removed on 2026-12-31." Reserve "will" for genuine commitments, not predictions. |

**Worked example.**

Mixed (confusing):
> The service caches tokens for 5 minutes. It used the original TTL of 1 hour
> until we will reduce it after seeing leakage incidents.

Repaired (each tense earns its job):
> The service caches tokens for 5 minutes. The original TTL was 1 hour; we
> reduced it in Q3 2025 after a leakage incident. The TTL will become
> configurable in v3.

**Pass tactic.** Highlight every verb. For each, ask: is this a past event, a
present truth, or a future plan? If three sentences in a row drift between
tenses without a reason, the prose lost its temporal frame. Restore it.

**When to break the rule.** Narrative sections (post-mortems, retrospectives)
sometimes use past tense throughout for cohesion, even when describing
currently-true facts. Tutorials may use present-tense imperative ("You run
the command. The system prints the result.") as a stylistic choice.

**Reference.** Google Developer Documentation Style Guide — Verb tense.
`https://developers.google.com/style/tense`

---

### The Pronoun-Antecedent Hygiene Pass

A pronoun without a clear antecedent forces the reader to backtrack or guess.
This is one of the highest-frequency hidden errors in technical prose — the
writer knows what "it" refers to, the reader does not.

**The rule.** Every "it," "this," "that," "they," and "those" must have an
unambiguous antecedent within the immediately preceding sentence (or, in
exceptional cases, within the preceding two sentences if there are no
competing candidates).

**The danger pronouns:**
- **It** — most commonly ambiguous when the previous sentence contains two
  possible referents.
- **This** / **That** — referring to "the previous statement" is ambiguous
  when the previous statement had multiple clauses.
- **They** — fine for groups; ambiguous when the previous sentence named
  more than one group.

**Worked examples.**

Ambiguous:
> "The migration introduced a new index, which slowed write throughput. It
> needs to be addressed."

What is "it"? The migration? The new index? The slowdown? The reader cannot
tell from this passage.

Repaired:
> "The migration introduced a new index, which slowed write throughput. The
> slowdown needs to be addressed."

Ambiguous:
> "We rolled back the deploy and reverted the schema migration. This
> resolved the issue."

Did the rollback resolve it, the schema revert, or both together?

Repaired:
> "We rolled back the deploy and reverted the schema migration. The combined
> action resolved the issue." — or — "The schema revert resolved the issue."

**Pass tactic.** Search the draft for every "this," "that," "it," "they" at
the start of a sentence. For each, ask: would a reader who joined the
document at this paragraph know what this pronoun refers to? If not, replace
with a specific noun.

**When to break the rule.** Conversational prose, marketing copy, and short
internal Slack messages tolerate looser antecedents because the context is
visible to the reader. Formal docs, RFCs, API references, and legal-adjacent
prose require strict antecedent hygiene.

**Reference.** Williams, *Style*, ch. 5 (Cohesion); Strunk & White §11.

---

### The Noun-Stack Pass

A "noun stack" is a chain of three or more nouns piled in front of a final
noun, with no prepositions to clarify the relationships. Common in technical
prose because nouns feel formal; the cost is that the reader must guess which
noun modifies which.

**The rule.** Limit noun stacks to two consecutive nouns. When you find
three or more, break the stack by inserting prepositions (of, for, in) or by
extracting the verbs hiding inside the nouns.

**Worked examples.**

Stacked (hard to parse):
> "Customer satisfaction score improvement initiative results review meeting."

Six nouns; the reader must guess which modifies which. Is it a meeting that
reviews results? A review of meeting results? A satisfaction initiative for
score improvement?

Unstacked:
> "Meeting to review the results of the initiative for improving customer
> satisfaction scores."

Longer, but each relationship is explicit.

Stacked:
> "Cluster replication lag monitoring alert threshold configuration."

Unstacked:
> "Configuration of the alert threshold that monitors replication lag in the
> cluster."

Or, better, restructure to verbs: "How to configure the alert that fires
when replication lag exceeds the threshold."

**Pass tactic.** Scan headings and section titles especially — noun stacks
breed in titles. Count consecutive nouns. If three or more, rewrite.

**When to break the rule.** Some technical terms are stable compound nouns
("connection pool exhaustion," "two-factor authentication setup") and
unstacking them produces stilted prose. The rule applies to ad-hoc stacks,
not to canonical terms-of-art.

**Reference.** Pinker, S. *The Sense of Style* (2014), ch. 4 — noun pile-ups;
PlainLanguage.gov — write strings of nouns into prose.

---

### The Read-Aloud Test (formalized)

The earlier "Reading Aloud as an Editing Tool" section is a habit; this is
the test you can run on every paragraph. Pinker calls it the "ear test."
Mid-twentieth-century stylists (Zinsser, *On Writing Well*; King, *On
Writing*) all reach for the same instrument: when prose fails the ear, it
will fail the eye.

**The rule.** Before declaring any document done, read every sentence aloud
at conversational pace. The ear catches what the eye glosses.

**What the read-aloud test catches:**

| Signal | What it reveals |
|--------|-----------------|
| You run out of breath before the period | Sentence too long; split it. |
| You stumble on a transition | Missing connective tissue between clauses. |
| Two consecutive sentences sound identical in rhythm | Lack of variation; reader will skim. |
| You hear the same word twice in three lines | Unintentional echo; replace one. |
| You feel the urge to skip a phrase | Throat-clearing or filler. Cut it. |
| A pronoun makes you re-parse the sentence | Antecedent unclear. Replace with a noun. |
| The sentence sounds like a flowchart | Nominalizations and noun stacks. |

**Pass tactic.** Read the entire document aloud once, slowly, then again at
normal pace. Mark anywhere the read fails — do not stop to fix. Then go
back and fix the marks together. Trying to read-aloud-then-fix one sentence
at a time loses the cumulative ear for the document's rhythm.

**Quiet alternative.** When you cannot read aloud (open office, library,
public transit), subvocalize: move your lips as if speaking but make no
sound. The ear test still works; you just look slightly mad while doing it.

**Audio alternative.** Use a screen reader, text-to-speech (macOS `say`,
NaturalReader, ElevenLabs read-back). Synthetic voices catch different
problems than your own voice — the rhythm flattens, exposing structural
weaknesses your reading might rescue.

**When to break the rule.** Highly formulaic docs (API reference, error code
catalogs) do not benefit from the ear test — they exist to be scanned, not
read.

**References.**
- Pinker, S. *The Sense of Style* (2014), ch. 6 (the ear test).
- Zinsser, W. *On Writing Well*, 30th anniversary ed. (1976/2006).
- King, S. *On Writing* (2000), "On the writer's toolbox" — reading aloud as
  the final pass.

---

## Quick Diagnostic Checklist

Before declaring a revision done, check:

- [ ] Structural pass complete — argument skeleton valid
- [ ] Every paragraph has a first sentence that states its point
- [ ] First sentences in sequence form a coherent argument
- [ ] No nominalizations where a verb is available
- [ ] Every sentence's subject names the actor
- [ ] Key content lands at sentence end, not middle
- [ ] Every *was/were + past participle* audited
- [ ] Metadiscourse, hedges, and redundant pairs cut
- [ ] Long sentences (>35 words) have visible architecture
- [ ] Parallel items are grammatically parallel
- [ ] 30% cut target met, or conscious decision to stop short
- [ ] Read aloud complete — no stumbles remain

---

## "Kill your darlings" — the discipline (Quiller-Couch attribution)

This skill already references the King formulation. The full provenance and a fuller working procedure follow.

**Rule.** Sir Arthur Quiller-Couch coined the line in *On the Art of Writing* (1916), lecture 12: "If you here require a practical rule of me, I will present you with this: 'Whenever you feel an impulse to perpetrate a piece of exceptionally fine writing, obey it — whole-heartedly — and delete it before sending your manuscript to press. Murder your darlings.'" Faulkner is often credited with "kill your darlings"; the closer attribution is Quiller-Couch, with King's *On Writing* (2000) restoring the phrase to wide circulation.

The rule is *revision discipline*, not anti-craft. The instruction is: identify the passages you would fight to keep, *examine why*, and cut anything you would fight to keep for reasons of ego rather than service to the reader.

**Operational procedure during revision.**

1. **Tag your darlings.** On the second editing pass, mark every passage where your gut says "this is the best writing in the piece." Tag them: in margin notes, in a colored highlight, in a separate file.
2. **Apply the load-bearing test.** For each tagged passage, ask three questions:
   a. *Function.* What argumentative, narrative, or informational work does this passage do that no other passage does?
   b. *Substitution.* Could a plainer 1–2-sentence version do the same work?
   c. *Ego.* Am I keeping this because it serves the reader, or because I am proud of it?
3. **Cut on a 2-of-3.** If two of three answers point to "non-load-bearing," cut it. Move the cut text to a scraps file — most darlings find a home in a later piece where they *are* load-bearing.

**Worked example — an architecture doc.**

Tagged darling: "The system, like a Swiss watch, ticks with the unhurried precision of long engineering."

- Function: tonal flourish; no argumentative or informational work.
- Substitution: "The system is reliable" (uglier; less satisfying to write; but does the same job).
- Ego: yes, you like the cadence.

Verdict: cut. The reader of an architecture doc does not want Swiss watches; they want to know whether the system is reliable under their workload.

**Counter-example — when the darling stays.**

Tagged darling: "We could not reproduce the outage in staging because staging does not see the customer's traffic shape."

- Function: load-bearing. It explains why the post-mortem cannot offer a clean reproducer, which is the question the reader will ask next.
- Substitution: any plainer phrasing loses the diagnostic content.
- Ego: not really; the sentence is workmanlike.

Verdict: keep.

**When to break it.** The rule applies to prose with a job to do (technical writing, business writing, journalism, expository nonfiction). In personal essays, memoir, fiction, and brand-voice marketing copy, voice and music may be the load. The rule is *test before cutting*, not *cut every flourish*.

**Diagnostic.** After cutting, re-read the surrounding paragraph. If it now reads as flat or under-served, the darling may have been doing more work than your three-question test caught. Restore it — but with conscious justification this time.

**References.**

- Quiller-Couch, A. *On the Art of Writing*. Cambridge University Press, 1916. Public domain — Project Gutenberg: https://www.gutenberg.org/ebooks/13892
- King, S. *On Writing: A Memoir of the Craft*. Scribner, 2000.

---

## Numbers conventions — when to spell out, when to use numerals

**Rule.** Numbers conventions diverge across house styles. The two dominant US guides differ on the most common case.

| Range | AP Stylebook | Chicago Manual (§9) |
|---|---|---|
| Zero through nine | Spell out (one, two, ... nine) | Spell out (zero, one, ... one hundred and beyond) |
| 10 and above | Numerals (10, 11, 100, 1,000) | Numerals **for technical/scientific work**; spell out **for general prose** through one hundred and round hundreds/thousands |
| Start of sentence | Spell out always | Spell out always (or recast) |
| Ages, percentages, dimensions, money, exact times | Numerals always | Numerals always |
| Adjacent numbers ("five 10-year-olds") | Use a numeral and a spelled form to avoid "5 10-year-olds" | Same |

**Pick one style and stick to it.** The most common failure is unconscious mixing — "five servers" in one paragraph and "5 servers" in the next.

**Operational rules across styles.**

- Always use numerals for: percentages (12%), money ($45), measurements (3 ms, 8 GB), exact times (3:15 p.m.), version numbers (v2.4), citations and page references (p. 27), and any number paired with a unit.
- Always spell out at the start of a sentence — or recast: "Twelve customers responded" or "We received responses from 12 customers."
- For very large numbers, mix: "1.2 million," "$3.4 billion." Spelling "one billion two hundred million" is almost always wrong.
- Hyphenate compound numbers under one hundred when spelled out: *twenty-one*, *ninety-nine* (Chicago §9.13).
- Decades: AP "the 1990s" or "the '90s"; Chicago "the 1990s" (no apostrophe). Never "1990's" (possessive form).

**Worked example.**

Inconsistent: "The team migrated 12 services in five weeks, with 1,200 ms median latency and only 3 incidents reported."

Chicago-consistent (technical): "The team migrated 12 services in 5 weeks, with 1,200 ms median latency and 3 incidents reported."

AP-consistent (general): "The team migrated 12 services in five weeks, with 1,200 ms median latency and three incidents reported."

**When to break it.** Tables, charts, and code samples always use numerals — no spell-out, ever. Direct quotations preserve the original. Legal and financial documents have their own conventions; defer to the house style or jurisdiction.

**References.**

- *Associated Press Stylebook*, current edition, "numerals" entry.
- *The Chicago Manual of Style*, 17th ed., Chapter 9 ("Numbers"), §§9.1–9.69.
- Microsoft Writing Style Guide — "Numbers." https://learn.microsoft.com/en-us/style-guide/numbers

---

## Em dash, en dash, hyphen — the three marks

**Rule.** Three horizontal marks with distinct jobs. Treat them as different punctuation, not as length variants of each other.

| Mark | Width | Use |
|---|---|---|
| Hyphen `-` | shortest | compound modifiers (high-frequency, well-known), word breaks at line ends, prefixes when needed (re-cover vs recover), phone numbers, URLs |
| En dash `–` | medium (en width) | ranges (pp. 12–34, 2001–2010, Mon–Fri), score lines (3–2), open compounds where one element is two words (New York–London flight) |
| Em dash `—` | longest (em width) | parenthetical asides, abrupt break in thought, dialogue interruption |

Chicago §§6.78–6.91 is the canonical reference. AP differs slightly: AP uses spaces around the em dash ("the team — three engineers and a manager — shipped on time"); Chicago does not ("the team—three engineers and a manager—shipped on time"). Pick one and stay consistent across the document.

**The em-dash discipline trap.** The em dash is ergonomically tempting because it accepts almost any grammatical relationship — comma, colon, parenthesis, semicolon — without the writer having to choose. This is why em-dash overuse is the most common signature of *both* AI-generated prose and human first drafts. The craft rule, per `writing-expert`'s anti-AI-isms section: ≤1 em dash per 100 words in polished human prose; ≤1 per 500 words for AI-ism detection thresholds. Most uses can be replaced by a comma, colon, period, or set of parentheses.

**Worked example.**

Overused: "The migration—planned for Q2—was delayed by an index bug—not the network issues we expected—and shipped in Q3."

Disciplined: "The migration, planned for Q2, was delayed by an index bug (not the network issues we expected) and shipped in Q3."

Em dash earning its place: "We rolled back—immediately, with the customer still on the phone."

**Range vs hyphen confusion.** "Pages 12-34" with a hyphen is wrong; "pages 12–34" with an en dash is correct. Most modern editors auto-convert; do not undo the conversion.

**Worked example — compound modifiers.**

Wrong: "post-mortem-style analysis," "post mortem style analysis."
Right: "post-mortem-style analysis" (Chicago) — or recast to avoid the stack: "an analysis in post-mortem style."

**When to break it.** Email, chat, code comments, and informal writing tolerate hyphens where en dashes would be technically correct. The discipline applies to publication-grade prose.

**Diagnostic.** Search for ` - ` (space-hyphen-space) and ` -- ` in the draft. Almost every match is a place an em dash, en dash, or restructured sentence belongs instead.

**References.**

- *The Chicago Manual of Style*, 17th ed., §§6.78–6.91 (dashes).
- AP Stylebook — "punctuation: dashes" entry.
- Butterick, M. *Practical Typography* — "Hyphens and dashes." https://practicaltypography.com/hyphens-and-dashes.html

---

## Quotation conventions — single, double, curly, straight

**Rule.** Quotation marks have two axes of variation: *single vs double* (which mark for the outer layer) and *straight vs curly* (the glyph form). House style picks one combination; the writer's job is consistency.

| Convention | Outer | Inner | Region |
|---|---|---|---|
| US (Chicago, AP, Microsoft) | Double `"` | Single `'` | "She said, 'Run it now,' and left." |
| UK (Oxford, Guardian) | Single `'` (often) | Double `"` (often) | 'She said, "Run it now," and left.' |
| Technical/code contexts | Straight `"` `'` | n/a | Code, file paths, identifiers |
| Body prose | Curly `"` `'` `'` `'` | n/a | Modern published prose |

**Punctuation inside vs outside the closing quote — the US/UK split.**

- US convention: periods and commas always go *inside* the closing quote, regardless of logic. "The flag is `--retry`," not "The flag is `--retry`,". Colons, semicolons, dashes, question marks, and exclamation points follow logical placement (inside if part of the quote, outside if not).
- UK and "logical" convention: punctuation goes inside only if it belongs to the quoted material. This is also the standard in computer code and command-line documentation, where putting a comma *inside* the quote would change the meaning of the quoted string.

The technical-writing exception is important: when quoting code, identifiers, commands, or strings the reader will type or paste, *always use logical (UK) placement* even in US house style. Otherwise you ship documentation that tells the reader to type `--retry,` instead of `--retry`.

**Straight vs curly.**

- Body prose, published output: curly quotes (`"` `"` `'` `'`). Most editors auto-convert.
- Code blocks, inline code, exact strings: straight quotes (`"` `'`). The reader must be able to copy and paste without breakage.
- Markdown source code: many writers leave the source straight and let the renderer curl on output; this is fine if rendering is consistent.

**Worked example — mixing the rules.**

Wrong: "Run `kubectl get pods,`" he said. (Comma inside changes the command.)
Right: "Run `kubectl get pods`," he said. (US convention overridden for the code; comma outside the backtick-quoted command.)

**Scare quotes and irony quotes.** Use sparingly. "We have 'tested' the failover plan" reads as sarcasm. If you mean *tested*, write *tested*; if you mean *not really tested*, write *we have not yet tested*.

**When to break it.** Direct quotation of source material preserves the original's quotation style. Legal and academic citation systems have specific conventions; follow them. Code and CLI documentation always uses straight quotes inside the code itself.

**References.**

- *The Chicago Manual of Style*, 17th ed., §§6.9–6.18 and §13 (quotations).
- AP Stylebook — "punctuation: quotation marks."
- Microsoft Writing Style Guide — "Quotation marks." https://learn.microsoft.com/en-us/style-guide/punctuation/quotation-marks

---

## Track changes / suggesting mode — etiquette

**Rule.** When editing someone else's document in a collaborative tool (Word Track Changes, Google Docs Suggesting Mode, Microsoft 365 Co-authoring, Notion comments, GitHub PR review), you are working in another person's draft. The etiquette is asymmetric — the editor adapts to the author's stage, not the other way around.

**Tool conventions.**

- **Microsoft Word Track Changes.** Each edit is attributed by user; accepted/rejected one at a time or in bulk; comments are anchored to a selection. Show All Markup is the editor's default; Simple Markup is the author's preferred reading view.
- **Google Docs Suggesting Mode.** Edits appear as colored insertions/strikethroughs; the author accepts or rejects via the side icon. Same anchored-comment semantics as Word.
- **GitHub PR review.** Line-level "Suggestion" blocks let the reviewer propose exact replacement text the author can commit with one click. Use suggestions for *small, mechanical* changes (typo, formatting, naming); use review comments for *larger structural concerns*.
- **Notion / Confluence comments.** Both support anchored comments. Notion's "edit history" is per-block; Confluence has page-level versioning.

**Etiquette rules across tools.**

1. **Match the author's stage.** A first-draft author needs *high-level direction* (structure, missing arguments, audience mismatch); a near-final author needs *line edits and polish*. Sending 80 comma-level edits to a first-draft author is a category error.
2. **Distinguish "must" from "consider."** Tag each suggestion. STC's technical-editing handbooks recommend explicit labels: *fix* (factual/grammatical error), *suggest* (improvement), *consider* (style preference). Without labels, the author cannot triage.
3. **Don't rewrite voice.** Suggested edits should preserve the author's voice unless voice itself is the problem (and then say so out loud, in a comment, not by silently rewriting). Imposing your voice through track changes is the most common reason authors reject editorial passes.
4. **Group related changes.** A single anchored comment that says "this paragraph would land harder with the conclusion in the first sentence" is more useful than five separate strikethroughs.
5. **Explain non-obvious cuts.** A deletion of three sentences without a comment looks like a power move. Add a one-line rationale.
6. **Resolve, don't hide.** Once a comment is addressed, mark it resolved. Don't delete it — the author may want the history. Microsoft and Google both retain resolved comments by default; preserve that.
7. **Accept your own typos.** Reviewers introduce typos. Always re-read your own track-changes additions before publishing.
8. **Sign off.** When you are done with a review pass, leave a top-level comment ("done — 12 suggestions, 4 fixes, 3 considers") so the author knows you are out of the document.

**Worked example — a bad review.**

A reviewer takes a 3,000-word draft architecture doc and:

- silently rewrites the introduction to match their own voice
- changes "use" to "utilize" 14 times (regression toward AI-ism)
- adds 47 comma changes with no labels
- leaves no top-level summary

The author opens the doc, sees a wall of red, and cannot tell what is load-bearing. The review is worse than no review.

**Worked example — a good review.**

Same draft. The reviewer:

- adds 3 anchored comments tagged FIX (factual errors)
- adds 8 anchored comments tagged SUGGEST with one-line rationales
- proposes 2 GitHub-style suggestions for clean replacements of awkward sentences
- leaves a top-level comment: "Structure is solid. Three fixes are blocking; eight suggestions are taste. Approve after fixes."

**When to break it.** A document in true crisis (legal exposure, customer-facing error, executive submission in an hour) may justify direct rewrite without ceremony — but make the rewrite visible (rename the file with a v-bump, or send a separate "I had to redline this hard, here is what changed" note). Silent overwrites destroy trust.

**References.**

- Microsoft 365 — "Track changes in Word." https://support.microsoft.com/en-us/office/track-changes-in-word-197ba630-0f5f-4a8e-9a77-3712475e806a
- Google Docs Help — "Suggest changes." https://support.google.com/docs/answer/6033474
- STC (Society for Technical Communication) — *Technical Editing in the 21st Century*, 2017 (review-etiquette chapters).
- Carolyn Rude & Angela Eaton, *Technical Editing*, 5th ed., Pearson, 2010 — Chapter on collaborative editing.
