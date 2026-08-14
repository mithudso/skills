# Technical Writing Craft — Extended Reference Topics

Load this file when the user needs depth on verb tense, "show your work" reasoning transparency, or the code-doc–specific show-first rules beyond what is in the main SKILL.md.

---

## Verb tense discipline

In technical writing, three tenses each have a default job. Mixing them without reason confuses readers about whether they are being told history, behavior, or plan.

### The rule

| Tense | Use for |
|-------|---------|
| **Present** | Current system behavior, API semantics, reference facts. "The function returns a cursor." "The cache evicts entries older than 5 minutes." |
| **Past** | Specific historical events, completed migrations, incident timelines, changelog entries. "We migrated to PostgreSQL in Q2 2024." "The outage lasted 47 minutes." |
| **Future** | Plans, scheduled changes, deprecations with a date. "The `/v1` endpoints will be removed on 2026-12-31." Reserve `will` for genuine commitments. |

### Worked example

Mixed tense (confusing):
> The service caches tokens for 5 minutes. It used the original TTL of 1 hour
> until we will reduce it after seeing token leakage incidents.

Repaired (each tense earns its job):
> The service **caches** tokens for 5 minutes. The original TTL **was** 1 hour;
> we **reduced** it in Q3 2025 after a leakage incident. The TTL **will
> become** configurable in v3.

Now: present describes current behavior, past describes a completed change, future describes a real commitment.

### When to break the rule

- Narrative sections (post-mortems, retrospectives) may use past tense throughout for cohesion, even when describing currently-true facts.
- Tutorials may use present-tense imperative ("You run the command. The system prints the result.") as a stylistic choice.

### References

- Google Developer Documentation Style Guide — Verb tense: present tense for state descriptions. `https://developers.google.com/style/tense`
- Microsoft Writing Style Guide — Verbs.

---

## "Show your work" — reasoning transparency

The math-classroom rule applies to docs that explain how a system arrives at a result: don't print only the answer. When a function returns a non-obvious value, a configuration produces a non-obvious effect, or a recommendation ranks options, show the reader the inputs, the steps, and the inference.

### The rule

For every claim in a doc that a careful reader might dispute or want to verify, expose one of:

1. The data the claim rests on (with a citation or example).
2. The computation or transformation that produced it.
3. The principle or constraint that ranks the options.

### Worked example

Bare conclusion (bad):
> Use `mongosh --quiet` in scripts.

Shown reasoning (good):
> Use `mongosh --quiet` in scripts. Without `--quiet`, `mongosh` prints a
> three-line banner to stdout on startup. Most scripts pipe `mongosh` output
> to `jq` or `grep`; the banner pollutes the pipeline and causes parse errors
> when downstream tools expect JSON.

The reader who knows why the banner is a problem now has the reasoning to apply the rule in new contexts (e.g., CI logs, CLI wrappers).

### When to break the rule

- API reference entries can be bare — the catalog format encodes that the reader is looking up a fact, not asking why.
- Quick-start tutorials suppress reasoning to keep the time-to-first-success short. Push reasoning into a linked explanation doc.

### References

- Polya, G. *How to Solve It*. Princeton University Press, 1945.
- Industrial Empathy. "Design Docs at Google." `https://www.industrialempathy.com/posts/design-docs-at-google/`

---

## Numbers conventions — technical-writing variant

Technical writing has stricter numbers conventions than general business prose. The rule is *every measured, indexed, counted, or referenceable number is a numeral, regardless of magnitude.*

**Always numerals in technical writing.**

- Measurements with units: 3 ms, 8 GB, 100 ms p99, 2 vCPU
- Versions: v1.2.3, Postgres 16, MongoDB 7.0
- Error codes, ports, identifiers: HTTP 500, port 27017, exit code 1
- API parameters and limits: max_retries=3, batch_size=1000, timeout=30s
- Counts of countable items in test or benchmark prose: "1 of 12 shards failed"
- Page numbers, table numbers, figure numbers, equation numbers
- Decimal numbers and percentages, always: 0.5, 99.9%, 1.2x
- Dates and times in ISO formats: 2026-05-29, 14:30 UTC

**Spell out only at the start of a sentence (or recast to avoid it).**

Wrong: "5 nodes failed before recovery." — never start with a numeral.
Right: "Five nodes failed before recovery." — or recast: "Recovery completed after 5 nodes failed."

Recast is almost always better — it puts the number in numeral form where the reader expects it for a measurement.

**Special case — counts of nontechnical entities.**

In running prose, counts of *teams, customers, people, regions, environments* can still spell out small numbers if the surrounding house style is Chicago-leaning ("three engineers reviewed the proposal"). But the moment the count is a *metric* ("3 engineers / 12-person team = 25% review coverage"), switch to numerals. Consistency within a paragraph beats absolute consistency across the document.

**Comma vs space vs underscore for large numbers.**

- US technical writing: comma every three digits ("1,200,000").
- SI / international scientific: thin space every three digits ("1 200 000"). Most authors substitute a regular space.
- Code and configuration: no separator ("1200000") or underscores in source languages that support them (`1_200_000` in Rust, Python, Ruby).
- Never the European decimal point as thousands separator ("1.200.000") in English-language documentation.

**Worked example.**

Poor: "The query took five hundred ms to complete on a 2 node cluster running version seven point zero."
Better: "The query took 500 ms to complete on a 2-node cluster running MongoDB 7.0."

**When to break it.** Marketing-facing technical content (release announcements, blog headlines) may relax the rule slightly to read more naturally; but body copy in the same blog post should hold the technical convention.

**References.**

- *The Chicago Manual of Style*, 17th ed., Chapter 9 (Numbers); §9.7 "Technical and scientific writing" gives the numerals-throughout rule.
- Microsoft Writing Style Guide — "Numbers." https://learn.microsoft.com/en-us/style-guide/numbers
- Google Developer Documentation Style Guide — "Numbers and measurements." https://developers.google.com/style/numbers
- IEEE Editorial Style Manual — section on numerals in technical prose.

---

## Four editing types — decision reference

When technical writing work is requested, first classify which editing type is needed before applying any pass. Running the wrong type wastes effort and can break prose that is already sound.

| Type | Focus | When to use | Do NOT use for |
|---|---|---|---|
| **Developmental / Structural** | Structure, arc, section ordering, logical flow | The document's argument or organization is wrong | Polishing prose in a well-structured doc |
| **Line editing** | Sentence-level word choice, syntax, voice, flow | Structure is sound; individual sentences are weak, unclear, or off-voice | Correcting grammar/spelling |
| **Copy editing** | Grammar, spelling, punctuation, style consistency | Structure and sentences are final; need mechanical cleanup | Reordering sections or rewriting sentences |
| **Proofreading** | Final error catch in formatted/laid-out output | Document is in its final format (PDF, HTML, published page) | Any content-level changes |

**Decision rule:** Run them top-down. Do not proofread a document that needs developmental editing. Never line-edit a section that will be cut.

*Sources: Jane Friedman, "The Differences Between Line Editing, Copy Editing, and Proofreading" — https://janefriedman.com/the-differences-between-line-editing-copy-editing-and-proofreading/; BlueLeaf Editing, "Editing 101: Understanding the 6 Different Types of Editing" — https://blueleafediting.com/understanding-the-different-types-of-editing/*
