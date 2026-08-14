# Writing Expert — Advanced Craft Reference

Deep reference for prose-craft concepts beyond the core SKILL.md rules: journalism structures (inverted pyramid, hourglass, bury-the-lede), achievement frameworks (STAR/SOAR/PAR/CAR), sentence-level devices (fragments, tricolon/isocolon, title vs sentence case, Oxford comma), cross-document consistency (style sheets and term banks), the show-don't-tell evidence rule, the Curse of Expertise, information scent, emphasis discipline (bold/italic/underline), and front-matter craft (TL;DR, kicker, nutgraf, deck, tabular-vs-prose). Loaded on demand from writing-expert when a task needs depth beyond the core skill. Worked examples and source citations are included per concept.

---

### Inverted pyramid (journalism)

The inverted pyramid orders information by decreasing importance: the most
load-bearing facts go in the first paragraph (the **lede**), supporting facts
follow, and background details land at the bottom. The structure originates
from newswire telegraphy — editors cut from the bottom when copy was too long
for the column, so the bottom had to be the least essential content.

**Rule:** the first paragraph must answer the reader's top question. If the
story were truncated after one paragraph, the reader should still get the
essential point. Each subsequent paragraph should add the next-most-important
detail.

**Worked example.** A 4-paragraph status update on a security incident:

1. (Lede) "Our SSO provider had a 47-minute outage today; all customer-facing
   logins failed between 14:03 and 14:50 UTC. No data was exposed."
2. (Supporting) Root cause: certificate rotation misconfiguration in our IdP.
3. (Context) Detection happened at 14:08 via synthetic monitoring; manual
   failover initiated at 14:32.
4. (Background) The IdP has been the auth gateway since 2024; this is the
   first outage of this duration.

A reader who stops after paragraph 1 still knows the impact and the safety
posture. A reader who reads to paragraph 4 gets full context.

**When to break the rule.** Inverted pyramid is wrong for:
- Persuasive prose (use Pyramid Principle or SCQA instead).
- Narrative writing where surprise is the point.
- Technical reference docs (use Diátaxis reference mode).

**References.** Associated Press *Stylebook*; Poynter Institute writing
fundamentals; Jakob Nielsen, "Inverted Pyramids in Cyberspace" (1996).

---

### STAR vs SOAR vs PAR vs CAR

Four closely related frameworks for narrating a specific past achievement or
incident. They share the same job (turn a vague claim into a concrete story)
but differ in emphasis.

| Framework | Expansion | Emphasis |
|-----------|-----------|----------|
| **STAR** | Situation, Task, Action, Result | Behavioral-interview standard; balanced |
| **SOAR** | Situation, Obstacle, Action, Result | Foregrounds the obstacle/conflict |
| **PAR** | Problem, Action, Result | Compressed; resume bullet form |
| **CAR** | Challenge, Action, Result | Same as PAR; "Challenge" framing |

**Rule of thumb.**
- Use **STAR** for behavioral interviews and full case studies — the Task
  step disambiguates "what was your role on the team."
- Use **SOAR** when the obstacle is the most interesting part of the story
  (post-mortems, leadership case studies).
- Use **PAR/CAR** for resume bullets and one-line achievement statements where
  brevity matters more than role clarity.

**Worked example.** Same achievement, four frames:

- STAR: "Our P99 read latency had drifted to 340ms (S). I was asked to bring
  it under 100ms before the Q3 customer review (T). I profiled the query
  path, added two compound indexes, and switched to read-from-secondary (A).
  P99 dropped to 62ms; the customer review went ahead on schedule (R)."
- SOAR: "Read latency was 340ms; the customer was threatening to escalate (S).
  The team had already tried index tuning twice without success (O). I
  re-profiled with the customer's actual query shape and found two missing
  compound indexes (A). P99 fell to 62ms (R)."
- PAR: "Reduced P99 read latency from 340ms to 62ms by re-indexing on
  customer query shape; saved an enterprise renewal."
- CAR: identical structure to PAR.

**When to break the rule.** STAR/SOAR/PAR/CAR don't fit narrative pieces
that depend on surprise, ambiguity, or unresolved tension. Use them when the
reader needs to extract evidence of competence from your story.

**References.** Bock, Laszlo. *Work Rules!* (Google hiring methodology, 2015);
US Department of Labor career-services PAR guidance.

---

### Sentence fragments (deliberate use)

A sentence fragment is a punctuated word group that lacks an independent
clause (no subject, no verb, or both). In default prose, fragments are
errors. In deliberate use, they create emphasis, rhythm, and voice that full
sentences cannot.

**Rule.** A fragment is acceptable when:
1. It punches a single idea with maximum compression. ("Not a chance.")
2. It echoes a preceding complete sentence and lands the close. ("We can ship
   this Friday. With three engineers. By dinner.")
3. It marks a deliberate pause or beat in voice-forward prose (sales copy,
   speeches, op-eds).

A fragment is wrong when:
1. The reader cannot resolve the missing subject/verb from context.
2. It appears so often that it stops being emphasis and becomes habit.
3. It substitutes for an argument the writer didn't make.

**Worked example.**

Weak: "We don't do six-month contracts. We don't do them. Not for any
customer."
Strong: "We don't do six-month contracts. Ever."

Weak (fragment as filler): "Important. We need to ship. Friday. With tests.
And docs."
Strong: "We ship Friday — with tests and docs. Non-negotiable."

**When to break the rule.** Almost never in formal regulatory writing, legal
prose, or API reference docs. Fine in marketing copy, op-eds, executive
narrative, and case studies where voice and emphasis matter.

**References.** Strunk & White §16 (fragment as deliberate device);
Stephen King, *On Writing* (fragments and rhythm).

---

### Title case vs sentence case

Two competing capitalization conventions for headlines, section titles, and
UI strings. Both are legitimate; mixing them within one document is the
error.

**Title case** capitalizes the first letter of every "major" word: nouns,
pronouns, verbs, adjectives, adverbs, and subordinating conjunctions; it
lowercases articles (a, an, the), short prepositions (of, in, on, for), and
coordinating conjunctions (and, but, or, nor). Example: "How to Write a
Status Report for a Skeptical Audience."

**Sentence case** capitalizes only the first word and any proper nouns.
Example: "How to write a status report for a skeptical audience."

**Rule of thumb.**

| Style guide / context | Preferred |
|-----------------------|-----------|
| Chicago Manual of Style (book titles, formal) | Title case |
| AP Stylebook (news headlines) | Title case (per AP's rules) |
| Microsoft Writing Style Guide | Sentence case (UI strings, headings) |
| Google Developer Documentation | Sentence case (headings, page titles) |
| Apple Human Interface Guidelines | Title case (some controls); sentence case elsewhere |
| Modern web / SaaS UI | Sentence case |

**Discipline:** pick one per document, declare it in your style sheet (see
"Style sheets and term banks" below), and hold it. The dominant trend in
2020s technical writing and SaaS UI is sentence case — it reads less
formal, scales better across translations, and avoids the "which words count
as major" debate.

**When to break the rule.** Proper nouns inside headings always retain their
own capitalization regardless of the surrounding style. Section numbers,
status indicators, and product names follow their canonical capitalization.

**References.** Chicago Manual of Style §8.159; AP Stylebook capitalization
chapter; Microsoft Writing Style Guide — Capitalization.

---

### Oxford comma policy

The **Oxford comma** (also serial comma) is the comma placed before the final
conjunction in a list of three or more items: "red, white, and blue" versus
"red, white and blue."

**The rule.** Pick a policy and hold it. The Oxford comma is required by the
Chicago Manual of Style, the MLA Handbook, the APA Publication Manual, the
US Government Printing Office, and most book publishers. It is omitted by
default in AP style (used in journalism) unless the sentence would be
ambiguous without it.

**Why the Oxford comma matters: the Maine dairy drivers case.**

In *O'Connor v. Oakhurst Dairy*, 851 F.3d 69 (1st Cir. 2017), Maine drivers
won an estimated $5 million in overtime back pay because of a missing Oxford
comma in the state's overtime law. The statute exempted from overtime:

> "The canning, processing, preserving, freezing, drying, marketing, storing,
> packing for shipment or distribution of [agricultural products]."

The drivers argued — and the First Circuit agreed — that "packing for
shipment or distribution" was a single activity (packing-for-shipment or
packing-for-distribution), not two separate activities (packing for shipment;
or distribution). Without an Oxford comma between "shipment" and "or
distribution," the law was ambiguous, and ambiguity in remedial legislation
runs in favor of the worker. The drivers, who distributed dairy but did not
pack, were therefore not covered by the exemption.

This is the canonical real-world case for the Oxford comma. Cite it when a
reviewer argues the serial comma is "just style."

**Worked example.**

Ambiguous: "I want to thank my parents, Ayn Rand and God."
Unambiguous: "I want to thank my parents, Ayn Rand, and God."

**When to break the rule.** AP-style newsrooms drop the Oxford comma by
default and add it only to resolve ambiguity. Outside news, default to
including it.

**References.**
- *O'Connor v. Oakhurst Dairy*, 851 F.3d 69 (1st Cir. 2017).
- Chicago Manual of Style §6.19.
- AP Stylebook — comma usage.

---

### Tricolon and isocolon (sentence rhythm)

Two rhetorical devices that exploit reader expectation about three-part
structures. Used deliberately, they make a sentence land harder than its
content alone would justify.

**Tricolon** is a series of three parallel clauses, phrases, or words.
Famous examples: "Veni, vidi, vici" (Caesar); "Government of the people, by
the people, for the people" (Lincoln); "Friends, Romans, countrymen" (Shakespeare).

**Isocolon** is a tricolon (or other multi-part series) in which each member
is the same syntactic length and structure. "I came, I saw, I conquered" is
both tricolon and isocolon — three independent clauses, each two words.

**Why threes work.** Cognitive research (Miller, "Magic Number Seven,"
1956; later refined by Cowan 2001) shows that three items feel complete,
two feel incomplete, four feel like a dump. Marketers exploit this with
"good, better, best" pricing tiers; speechwriters exploit it with three-beat
applause lines.

**The rule.** Use tricolon to land a closing claim, a tagline, or an emphatic
list. Use isocolon when the three parts are genuinely parallel in
significance — if one is more important, breaking parallelism signals that
imbalance to the reader.

**Worked example.**

Flat: "The migration will be faster, cheaper, and safer than the previous
approach, and it also has fewer dependencies."
Tricolon: "The migration is faster, cheaper, and safer." (drop the fourth
beat; let the three land)
Isocolon: "Faster to ship. Cheaper to run. Safer to operate." (three
parallel two-word phrases; maximum punch)

**When to break the rule.** When the natural list is two or four items,
do not force it into three. A four-item list ("Specificity, brevity, clarity,
honesty") is fine if the four items genuinely matter equally. Padding three
into four (or trimming four into three) to chase rhythm produces hollow
prose.

**References.** Lanham, Richard. *Analyzing Prose* (rhetorical figures);
Miller, George. "The Magical Number Seven, Plus or Minus Two." *Psychological
Review* 63(2), 1956; Cowan, Nelson. "The Magical Number 4." *Behavioral and
Brain Sciences* 24(1), 2001.

---

### Style sheets and term banks (cross-document consistency)

A **style sheet** is a project-specific reference that records every
capitalization, spelling, hyphenation, and terminology decision the team has
made — so the next writer doesn't relitigate them. A **term bank** is the
subset of the style sheet that names approved terms (and forbidden ones) for
a specific domain, product, or audience.

**Why this matters.** Without a style sheet, the same product gets written
as "log-in" on the marketing page, "login" in the docs, and "log in" in the
UI strings — all in the same week. Readers notice the drift, even when they
can't name it; their trust in the doc set erodes.

**Minimum viable style sheet** (one page, kept in the repo or a wiki):

```
Product name:        MongoDB Atlas  (not "Atlas", "mongo atlas", "MongoDB cloud")
Feature names:       Search Nodes (cap S, cap N)
                     Online Archive (cap O, cap A)
Action verbs:        "set up" (verb), "setup" (noun)
                     "log in" (verb), "login" (noun)
Numbers:             spell out one-nine; numerals for 10+; numerals always for units
Headings:            sentence case (per MSWG)
Oxford comma:        yes
Quotation marks:     curly in prose, straight in code
Date format:         YYYY-MM-DD (ISO 8601)
Banned terms:        delve, leverage, robust, paradigm  (see writing-expert Tier 1)
```

**Worked example.** A docs team that adopted a one-page style sheet found
that PR review time on documentation dropped 30% — reviewers stopped
arguing about "is it log-in or login" and shipped faster (Microsoft
content-engineering retrospective, 2023).

**When to break the rule.** Stylistic exceptions are fine when they're
deliberate and noted. A blog post in a more conversational voice may
override the style sheet for register; an internal Slack message doesn't
need to honor it at all. The rule is about external, durable artifacts.

**Term bank discipline.** When the style sheet is being used for accessibility,
inclusivity, or legal compliance, the inclusive-language and plain-language
skills also apply — coordinate term-bank decisions across all three.

**References.**
- Microsoft Writing Style Guide — Word lists; Style sheets.
- Chicago Manual of Style §2.55 (style sheets in editorial practice).
- Editorial Freelancers Association: style sheet templates.

---

## Output shape

| Request type | Return |
|---|---|
| Draft (new document) | The document only — no commentary unless the user asked for it |
| Edit (existing document) | The full edited document, then a one-line change summary (what changed and why) |
| Review / critique | Bulleted findings by section; no rewrite unless asked |
| Framework question | Concise answer + one concrete example; pointer to `references/writing-expert-context.md` for depth |
| AI-ism removal only | The cleaned text only; flag count in parentheses at the end |

## Common failure modes

| Symptom | Cause | Fix |
|---|---|---|
| Output sounds over-formal | Register misread — defaulted to C-suite | Re-check audience; lower to engineer/customer register |
| AI-isms reappear in rewrite | Self-check not run | Run Principle #8 scan on your own output before delivering |
| Answer drifts from user's goal | Scope creep during editing | Re-read user's stated goal; trim anything not requested |
| Paragraph structure breaks | Topic-sentence rule not applied | Apply `## Sentence and paragraph craft → Topic-sentence-first` |
| Too long | No length signal obtained | Invoke Principle #7 confirmation before writing |

## Review Handoff

After producing any document over 200 words, append: "Run `/document-critique` for argument strength, evidence gaps, tone consistency, and AI-ism detection."

## Handling Sensitive Content

Do not quote back, store, or echo customer names, revenue figures, PII, or confidential contract data beyond what is needed to complete the edit. If the input contains sensitive data, work with it in place and do not reproduce it in examples or explanations.

---

## Show, don't tell — the evidence rule for prose

**Rule.** "Don't tell me the moon is shining; show me the glint of light on broken glass." — Anton Chekhov, letter to his brother Aleksandr (May 10, 1886). Stephen King, in *On Writing* (2000), restated it as a working rule: descriptive *assertions* about a quality (busy, tense, impressive, complex) are weaker than *evidence* of that quality.

In technical and business prose the rule reads: replace adjectival claims with measurements, observations, or examples that *force* the reader to draw the conclusion themselves. "The migration was complex" is a *tell*. "The migration touched 17 services, blocked 2 release trains for 6 weeks, and required a 4-region cutover sequence" is a *show* — and the reader concludes "complex" without you saying it.

**Worked example — status report.**

- Tell: "Onboarding has been challenging."
- Show: "Three of the last five customers escalated within their first 60 days; the median time-to-first-value slipped from 14 days in Q3 to 31 days in Q4."

The *show* form gives the reader the evidence and the *tell* (challenging) becomes self-evident.

**Worked example — executive memo.**

- Tell: "The legacy system is fragile."
- Show: "The legacy system has paged on-call 11 times this quarter; the last three pages required cold restarts that broke active sessions."

**When to break it.**

- *Headlines, taglines, and BLUF lines*: a single declarative claim *is* the point. "Migration is risky" works in a one-line summary that the body then proves.
- *Definitions and glossary entries*: tell directly, do not dramatize.
- *Repeat assertions* in a long document: after evidence has been shown once, you may *tell* in later references for compactness.

**Diagnostic.** Highlight every adjective in the draft. For each, ask: "what observation, number, or example would make the reader say this word themselves?" If you cannot answer, the adjective is doing work the evidence should be doing. If you can, replace the adjective with that evidence.

**Anti-pattern.** Stacking adjectives ("robust, scalable, seamless platform") is the *opposite* of show-don't-tell — three adjectival tells with no evidence. See the **Anti-AI-ism Rules** section in `writing-expert` SKILL.md for why these specific words are usually filler.

**References.**

- Chekhov, A. Letter to Aleksandr Chekhov, May 10, 1886. Collected in *Letters of Anton Chekhov to His Family and Friends*. Translated by Constance Garnett, 1920.
- King, S. *On Writing: A Memoir of the Craft*. Scribner, 2000, pp. 173–177.

---

## The Curse of Expertise (distinct from the Curse of Knowledge)

**Rule.** Two related but distinct cognitive failures bite technical writers; treating them as one obscures the fix.

- *Curse of Knowledge* (Pinker, Heath brothers): once you know something, you cannot remember not knowing it. Symptom: assuming the reader already shares your context. This skill's `## Sentence and paragraph craft → Pinker Curse of Knowledge` section covers it.
- *Curse of Expertise* (Willingham; sometimes called the *expert blind spot*): with deep mastery, your reasoning becomes *automatic and proceduralized*. You skip steps not because you forgot the reader, but because *you no longer perform those steps consciously yourself*. Symptom: tutorials and explanations that have intermediate steps missing not from forgetfulness but because the writer's brain has chunked them away.

The fix is different. The Curse of Knowledge fix is *audience modeling*: imagine what the reader knows. The Curse of Expertise fix is *step decomposition*: write the procedure for a novice, then run it as a novice would, *recording every micro-decision*. The gaps you find are the chunks your expertise hid from you.

**Worked example — a deployment runbook.**

A Curse-of-Knowledge failure: "Connect to the cluster" (you forgot the reader does not know which cluster).
A Curse-of-Expertise failure: "Run the deploy script" (you remembered to name the script, but you skipped that an expert silently confirms the branch, the environment flag, the kubectl context, and the on-call coverage *before* running it — because those checks are now invisible to you).

**Worked example — code review.**

Curse of Knowledge: "use a coroutine" (reader doesn't know what a coroutine is).
Curse of Expertise: "use a coroutine" (you know that doing so requires picking a dispatcher, scoping the cancellation, and choosing a supervisor strategy — but you did not say so because you do those automatically).

**When to break it.** Reference documentation aimed at peer experts may legitimately compress to the chunked form — peers share the same automatisms. The rule applies to tutorials, onboarding docs, runbooks for rotations, and any document where the *reader's expertise level is below the writer's*.

**Diagnostic.** Read your procedure aloud while watching a colleague at the same level as the target reader attempt to follow it. Every time they hesitate, you have hit an expertise-chunked step. Names for these are "implicit prerequisites," "missing scaffolding," or "expert's blind alley."

**References.**

- Willingham, D. T. *Why Don't Students Like School? A Cognitive Scientist Answers Questions About How the Mind Works and What It Means for the Classroom*. Jossey-Bass, 2009 (chapter on the difference between novice and expert thinking).
- Heath, C., & Heath, D. *Made to Stick: Why Some Ideas Survive and Others Die*. Random House, 2007 (Curse of Knowledge framing).
- Hinds, P. J. "The curse of expertise: The effects of expertise and debiasing methods on prediction of novice performance." *Journal of Experimental Psychology: Applied* 5(2), 1999.

---

## Information scent — Pirolli & Card (information foraging)

**Rule.** Readers do not read; they *forage*. At every link, heading, sidebar callout, or paragraph break they evaluate whether the next chunk smells like the answer they are hunting. Strong *information scent* is text that lets a reader predict, accurately, what they will find by continuing. Weak scent is generic, abstract, or witty-but-empty.

Pirolli & Card (PARC, 1995, 1999) modeled web reading as analogous to predator foraging: readers spend cognitive effort proportional to expected information return. A heading like "Considerations" or "Notes" has near-zero scent — the reader cannot predict the payoff. A heading like "Choosing between batch and streaming ingest" or "What changed in 7.2 from 7.0" has strong scent.

**Worked example — table of contents in an architecture doc.**

- Low scent: "1. Introduction. 2. Background. 3. Approach. 4. Discussion. 5. Considerations. 6. Conclusion."
- High scent: "1. Why the current ingest path won't scale past 50k/s. 2. Two options: lock-free queue vs. partitioned consumer. 3. Why we picked the partitioned consumer. 4. Cutover plan and rollback triggers. 5. What we explicitly did not do (and why)."

The high-scent version lets a reader skim the TOC and jump precisely to the section that matches their question.

**Worked example — link text.**

- Low scent: "[Click here](...)", "[Read more](...)", "[Documentation](...)".
- High scent: "[How retry-after-backoff is configured per route](...)", "[Why we dropped support for HTTP/1.0 in 2.4](...)".

**When to break it.** Marketing and brand voice sometimes deliberately use *low-scent, high-curiosity* headings ("The mistake every CTO makes") to drive click-through. That is a different optimization — engagement, not foraging — and belongs in `sales-and-marketing-copy`, not in technical or report writing. Inside technical docs, low-scent headings are almost always a bug.

**Diagnostic.** Read only the headings, link text, and first sentences of each paragraph in the document. Can a target reader navigate to the answer to their most common question using only those signals? If not, your scent is too weak. See also `microcopy-and-ui-writing` for the UI variant of the same principle.

**References.**

- Pirolli, P., & Card, S. "Information foraging in information access environments." *Proceedings of CHI '95*, 1995. https://dl.acm.org/doi/10.1145/223904.223911
- Pirolli, P., & Card, S. "Information foraging." *Psychological Review* 106(4), 1999, pp. 643–675.
- Nielsen, J. "Information Scent: How Users Decide Where to Go Next." NN/g, 2003. https://www.nngroup.com/articles/information-scent/

---

## Bold, italic, underline — emphasis discipline

**Rule.** Emphasis devalues with use. Chicago Manual of Style §§7.49–7.55 treats italic as the workhorse for *load-bearing emphasis* (a word the reader must not skip; a technical term on first appearance; a title), and bold as a *structural* device (run-in heads, key terms inside a definition, scannable labels in lists). Underline in modern prose is reserved for hyperlinks; using it for emphasis fights the link affordance and is now considered a typographical error in screen text.

Microsoft's Writing Style Guide is more permissive of bold in UI and technical documentation (bold for UI labels, key concepts, the action a user must take). But the underlying principle is the same: *every emphasis convention you adopt must be used consistently, and used sparingly enough that emphasis still means something*.

**Operational rules.**

| Mark | Use for | Do not use for |
|---|---|---|
| *Italic* | technical terms on first appearance; titles of works; load-bearing emphasis; foreign phrases on first use | run-in heads in lists; long passages; substitute for bold in scannable docs |
| **Bold** | UI labels (Settings > Privacy); key terms in a definition; first-line scannable labels in bullet lists | load-bearing emphasis in flowing prose (use italic); decoration; second occurrences of the same term |
| Underline | hyperlinks; legal-style document titles in some house styles | emphasis in body prose; UI labels; key terms |
| `Code font` | identifiers, commands, filenames, exact strings the reader must type or recognize | proper nouns of products; generic technical concepts |

**Frequency budget.** As a heuristic: if more than ~5% of body words in a paragraph carry emphasis, no word is emphasized anymore. The reader's eye flattens the cues. Cut the weaker ones.

**Worked example — a definition list.**

- Bad: "**The retry policy** governs **how many times** the client will **retry** a **failed request** before **giving up**." (Five bold spans; nothing reads as emphasized.)
- Good: "**Retry policy** — the number of times the client retries a failed request before surfacing the error. Default is *three attempts with exponential backoff*; configurable via `client.retry.maxAttempts`."

**When to break it.** Marketing copy and certain UX patterns (banner alerts, success states) legitimately use bold for tonal weight rather than structural marking. Long-form fiction and personal essays may use italic for inner thought, foreign language, or stylistic voice — registers outside this skill's scope.

**Diagnostic.** Read the document. Each emphasis mark should answer: *what would the reader miss if I removed this italic / bold?* If the answer is "nothing meaningful," remove it.

**References.**

- *The Chicago Manual of Style*, 17th ed., §§7.49–7.55 (emphasis, italic, bold).
- Microsoft Writing Style Guide — "Text formatting." https://learn.microsoft.com/en-us/style-guide/text-formatting/formatting-text-in-instructions
- Butterick, M. *Practical Typography* — "Bold or italic." https://practicaltypography.com/bold-or-italic.html

---

## TL;DR conventions — when, where, what

**Rule.** TL;DR ("too long; didn't read") originated as Usenet/Reddit-era reader slang circa 2003 and migrated into Tim Ferriss-style blog conventions, then into corporate writing. It is now a load-bearing front-matter device, *not* a casual sign-off. Use it deliberately.

This skill's `## Front-matter conventions` table already specifies that TL;DR is *~1 paragraph at the top of long documents, stating the conclusion*. This subsection extends that with format, placement, and selection rules.

**When to add a TL;DR.**

- Documents over ~600 words where the reader's first question is "what is the answer?"
- Documents addressed to mixed-audience or skim-first readers (executives, on-call engineers, customers).
- Internal memos that propose a decision (the TL;DR is the proposed decision).
- Status reports and post-mortems (the TL;DR is the verdict).

**When to skip.**

- Documents under ~300 words — TL;DR is longer overhead than the document.
- Reference documentation organized for lookup, not reading — there is no "summary" of an API reference.
- Tutorials and step-by-step guides — the summary *is* the step list; a TL;DR would duplicate it.
- Narrative pieces where the payoff is intentional (release announcements with reveal arcs).

**Format.**

- **Label.** "TL;DR" remains the canonical form in technical/Internet writing. "Summary" or "Bottom line" works in more formal registers. "BLUF" (Bottom Line Up Front; military origin) is interchangeable in defense, government, and many corporate environments — see `executive-comms` for BLUF rules.
- **Length.** 1–3 sentences. If your TL;DR runs to a paragraph it is becoming an abstract; relabel and shorten.
- **Placement.** Immediately after the title and before any preamble. Never below the first heading.
- **Content.** State the *conclusion or recommendation*, not the topic. "TL;DR: we should defer the upgrade to Q3 because of the unresolved index bug" is correct. "TL;DR: this document discusses the upgrade timing" is wrong — it summarizes the *document*, not the *answer*.

**Worked example.**

- Wrong (topic summary): "TL;DR: This memo covers our analysis of the Q2 incident, root cause, and follow-up actions."
- Right (conclusion): "TL;DR: The Q2 incident was caused by an unbounded retry loop in the order service. The fix is shipped in 4.7.2. We are adding a circuit breaker in 4.8 to prevent recurrence."

**When to break it.** In legal-adjacent writing, regulated communications, and academic papers, a TL;DR is inappropriate — those registers use abstracts, executive summaries, or no summary at all. See `academic-and-citation-writing` and `legal-adjacent-writing`.

**References.**

- Wikipedia, "TL;DR." https://en.wikipedia.org/wiki/TL;DR
- Ferriss, T. *The 4-Hour Workweek* blog corpus — popularized TL;DR convention in business writing.
- AP Stylebook entry on "internet jargon" and Reddit's TL;DR convention archive on reddit.com/wiki.

---

## Kicker — the closing sentence

A *kicker* is the final sentence (or short paragraph) of a piece. In journalism craft it's named the kicker because while the lede drives the reader into the story, the kicker drives the story into the reader. Most business prose under-invests here; the last sentence is the one the reader carries away.

**The rule.** A kicker has three jobs, in order:

1. *Signal the end.* The reader should feel the door close — no ambiguity about whether more is coming.
2. *Nail the central point.* The reader should be able to recite the piece's takeaway after reading only the kicker.
3. *Resonate.* The line should be memorable — through surprise, a turn back to the lede, a small inversion, or an image. Not by being clever.

**Why it matters in business writing.** Reports, memos, post-mortems, and QBRs almost always die in the last paragraph — usually with a vague "We will continue to monitor and iterate." That's not a kicker; it's a sigh. A document whose last sentence is forgettable is a document whose conclusion is forgettable.

**Worked example — internal memo on a deprecation decision.**

- Weak (default ending): "We will continue to evaluate the situation and provide further updates as needed."
- Strong (kicker that closes the loop to the lede): "We started this conversation by asking whether to keep the legacy gateway alive for one more quarter. We're not."
- Strong (kicker by image): "The legacy gateway has been running on borrowed time for two years. We're returning it."

**Three reliable kicker patterns.**

1. **Full circle.** Return to a phrase, image, or claim from the lede. Resonance comes from the echo.
2. **The flat declarative.** A short, blunt sentence after a complex argument. The simplicity itself is the kick — Strunk & White's "omit needless words" applied to the last line.
3. **The named consequence.** Spell out the concrete outcome readers will feel — not "this will improve things" but "the support queue will be empty by 5pm Thursday."

**When to break it.** Reference documentation, API docs, and runbooks don't need kickers — the document ends when the steps end, and grafting a literary close onto a how-to is performative. Save kicker craft for prose where a reader follows an argument or narrative end-to-end (memos, reports, post-mortems, QBR commentary, blog posts, release notes).

**Anti-pattern.** Ending with a quote you couldn't write better yourself. Roy Peter Clark calls the trailing pull-quote ending "the default ending in American journalism" and argues it dilutes more than it lands. If a quote is doing the work the kicker should do, you haven't found the kicker yet — you've outsourced it.

**Diagnostic.** Read only the final sentence aloud, without the document in front of you. Can you tell what the document was about? Does the sentence feel like an ending or like the next paragraph fell off? If neither test passes, the kicker is missing.

**References.**

- Roy Peter Clark and the Poynter Institute — [Putting Endings First](https://www.poynter.org/reporting-editing/2004/putting-endings-first/), the canonical treatment of kicker craft.
- The Open Notebook — [Good Endings: How to Write a Kicker Your Editor — and Your Readers — Will Love](https://www.theopennotebook.com/2015/11/24/good-endings-how-to-write-a-kicker-your-editor-and-your-readers-will-love/).
- Longreads — [Sticking the Landing: On Kickers](https://longreads.com/2024/11/26/kickers-journalism-writing-craft/), modern examples of long-form kicker craft.

---

## Nutgraf — the "what this story is about" paragraph

The *nutgraf* (also nut graph, nut graf, or nutgraph — short for "nutshell paragraph") is the paragraph that tells the reader what the piece is about and why it matters. In journalism it sits between the lede and the body, typically the second to fifth paragraph. It exists because not every reader makes it past the opening hook — the nutgraf is the contract: *here is the story; here is why it's worth your time.*

**The rule.** Every piece longer than ~500 words needs a nutgraf. It must:

1. State the *what* in one or two sentences (the central claim, finding, or story).
2. State the *why now* — what makes this timely or relevant to *this reader, this week*.
3. Sit no later than the fifth paragraph. Past that point you've lost the readers who needed it.

**Why business writers should care.** The nutgraf is the answer to "so what?" — and "so what?" is the single most common silent objection from skim-readers. A memo, a status report, or a QBR section without a nutgraf forces readers to deduce relevance from context, and most of them won't bother.

**Worked example — internal status update.**

- Without a nutgraf (reader has to guess why this matters): "The latency regression in service X was traced to a misconfigured retry policy. Engineering has rolled out a fix. Postmortem to follow."
- With a nutgraf: "The latency regression in service X was traced to a misconfigured retry policy. **This is the third retry-related incident this quarter, and the pattern is now visible enough that we're making retry hygiene a Q3 reliability goal.** Engineering has rolled out a fix. Postmortem to follow."

The bolded sentence is the nutgraf — it tells the reader what the story is *about* (a pattern, not a one-off) and *why now* (a Q3 goal is being declared).

**Five-sentence nutgraf scaffold** (adapted from Poynter and Nieman Storyboard guidance):

1. The news or finding.
2. The context — what came before that makes this matter.
3. The stake — who is affected and how.
4. The "so what?" — the angle the piece will pursue.
5. The road map — a one-sentence preview of the rest of the piece (optional but useful in long memos).

**When to break it.** Pure reference documentation, API docs, runbooks, and tutorials don't need a nutgraf — their relevance is presumed by the reader having opened them. Skip the nutgraf in any piece a reader arrived at via active search, not via skimming.

**Diagnostic.** Delete the nutgraf candidate paragraph from your draft. Does the reader still know *why* they're reading this piece? If yes, the paragraph was filler. If no — you just found your nutgraf, and you should restore it and make sure it lands by paragraph five.

**References.**

- [Nut graph — Wikipedia](https://en.wikipedia.org/wiki/Nut_graph) — origin and canonical definition.
- Roy Peter Clark, Poynter Institute — [The nut graf tells the reader what the writer is up to](https://www.poynter.org/archive/2003/the-nut-graf-part-i/).
- Nieman Storyboard — [Nut grafs: Seven steps to score a winning story structure](https://niemanstoryboard.org/2021/10/26/nut-grafs-seven-steps-to-score-a-winning-story-structure/).

---

## Subhead / deck — the line under the headline

The *deck* (also dek, subhead, subheadline, or standfirst in British usage) is the secondary line of text that sits between the headline and the body. The headline (hed) gets the reader to stop. The deck gets them to read. Together they convince a skim-reader to invest 5–30 minutes in the piece.

**The rule.** A deck does three things:

1. Expands what the headline could not fit.
2. Answers the question the headline raised.
3. Earns the scroll — by promising something specific the reader doesn't already know.

**Sizing.** Heds are 6–10 words. Decks are 10–25 words. Anything longer is no longer a deck; it's the opening paragraph wearing the deck's hat.

**Worked example — an internal post-mortem.**

- Hed only: "Q2 Cache Eviction Incident"
- Hed + deck: "Q2 Cache Eviction Incident — **A 4-line config change took down checkout for 22 minutes. Here's how it happened and the three guardrails we're adding.**"

The deck does the work: it names the surprise (4 lines), the cost (22 minutes), the surface (checkout), and the promise (three guardrails). A reader who reads only the hed + deck already understands the piece — and is now more likely to read the rest.

**Deck-specific anti-patterns.**

- *Restating the hed.* "Q2 Cache Incident — A look at what happened in the Q2 cache incident." Adds zero information.
- *Generic context.* "Q2 Cache Incident — Sometimes infrastructure can be unpredictable." Adds atmospheric noise, no specifics.
- *Burying the promise.* "Q2 Cache Incident — On May 8 at 14:02 UTC, cache nodes began rejecting writes…" That's the lede, not the deck. Decks summarize; ledes narrate.

**When to break it.** Very short documents (under ~300 words) and reference pages (API docs, runbook entries) don't need decks — the reader is already committed by the time they arrive. Decks are for any prose long enough that a skim-reader needs convincing to scroll: blog posts, release announcements, post-mortems, memos that exceed two pages, white papers.

**Diagnostic.** Cover the body of the piece with your hand. Read only the hed and deck. Can you state the piece's central claim and what makes it worth reading? If not, the deck has work to do.

**References.**

- MasterClass — [How to Write a Subheading: 4 Tips for Writing a Dek](https://www.masterclass.com/articles/how-to-write-a-dek).
- River — [Hed and Dek in Journalism: Write Headlines That Get Clicks](https://rivereditor.com/blogs/generate-hed-dek-headline-subhead), modern usage with sizing guidance.

---

## Hourglass structure — the third way between inverted pyramid and narrative

The *hourglass* is a hybrid document structure named by Roy Peter Clark (Poynter, 1983) that combines the inverted pyramid (most important facts first) with a chronological narrative. It has three parts: **the top** (4–6 paragraphs of inverted pyramid summary), **the turn** (a transitional line that signals the story now starts at the beginning), and **the bottom** (a chronological narrative of how events unfolded).

It is the right structure when readers need both *the verdict fast* and *the story in order* — which is most post-mortems, most incident reports, most case studies, and many memos that explain a decision.

**The three structures compared.**

| Structure | Top | Middle | Bottom | Best for |
|---|---|---|---|---|
| **Inverted pyramid** | Conclusion / lead | Supporting facts (descending importance) | Background, context | Hard news, executive memos with a single decision |
| **Hourglass** | Summary lead (4–6 paragraphs) | The turn (transitional line) | Chronological narrative | Post-mortems, incident reports, case studies, investigative pieces |
| **Narrative ("zigzag" or chronological)** | Scene-setting hook | Events in chronological order, alternating perspective or action | Climax / resolution | Feature stories, customer narratives, founder origin pieces |

**Why the hourglass beats a pure inverted pyramid for post-mortems.** Inverted pyramid forces every fact into "is this more or less important than the previous fact?" That works for news, but post-mortems need to teach the *sequence* — what was true at 14:02, what changed at 14:07, what the on-call saw at 14:11. The hourglass gives readers the verdict fast (the top) *and* lets them learn from the timeline (the bottom).

**Worked example — incident report.**

- *The top (inverted pyramid):* What broke, when, who was affected, severity, current status, owner, time to detection, time to resolution, root cause class.
- *The turn:* "Here is how the incident unfolded, in order."
- *The bottom (chronology):* Timeline starting before the trigger event, walking through detection, escalation, mitigation, full resolution, and follow-up.

**Anti-pattern — burying the timeline at the top.** Putting the chronological narrative in the first paragraph is the most common hourglass failure. Readers who only have 90 seconds need the verdict, not the story. Lead with what happened, end with how.

**When to break it.** Use a pure inverted pyramid when there is no useful chronology (a decision memo doesn't need a timeline). Use a pure narrative when the chronology *is* the point and there is no verdict to deliver up front (a customer success story, a founder's-origin piece). Use the hourglass only when both halves carry weight.

**Diagnostic.** After drafting, ask: "Does a reader who reads only the first 4 paragraphs know everything they need to know to make a decision?" If yes, the top is doing its job. "Does a reader who needs to learn the timeline get it without re-reading the summary?" If yes, the bottom is doing its job.

**References.**

- Roy Peter Clark, Poynter Institute — [The hourglass: serving the news, serving the reader](https://www.poynter.org/reporting-editing/2003/the-hourglass-serving-the-news-serving-the-reader/), the original 1983 framework and its modern application.
- Pressbooks / CWI — [Newswriting Structures: The Inverted Pyramid and Beyond](https://cwi.pressbooks.pub/introductiontojournalismandnewswriting/chapter/chapter-5-newswriting-structures-the-inverted-pyramid-and-beyond/), the structure compared to the inverted pyramid and chronological narrative.

---

## Bury the lede — the anti-pattern and the deliberate delayed-lede technique

"Burying the lede" means opening a piece with secondary information, forcing the reader to dig (or skim, or quit) before reaching the central point. In business writing it is the single most common structural error — and it is also, very occasionally, a deliberate craft choice. Treat the two cases separately.

**The accidental bury (almost always wrong).** Symptoms: the first paragraph is throat-clearing, scene-setting, context, or housekeeping. The actual news, decision, or recommendation is in paragraph three or later. The reader thinks they know what the piece is about based on paragraph one, decides it's not for them, and leaves before discovering it actually was.

**The deliberate delayed lede (occasionally right).** A *feature lede* or *anecdotal lede* opens with a scene, a person, or a vignette and arrives at the central point a few paragraphs later. It is a narrative tool — used to build context, generate empathy, or surprise the reader with a turn. It works only when the delay is *earned*: the scene must be specific, vivid, and clearly heading somewhere.

**When the deliberate delayed lede is right.**

- Feature stories where the human angle is the point (founder profiles, customer narratives, longform incident reflection pieces).
- Pieces where the reader's surprise at the central point is part of the payoff.
- Audiences who arrived for the writing, not the news — readers of a long-form newsletter or magazine, not skimmers of a memo.

**When it's wrong (most business writing).**

- Memos with a recommendation. The recommendation goes in sentence one.
- Status reports. The status goes in sentence one.
- Post-mortems. The verdict goes in the first paragraph; the narrative goes after the turn (see hourglass above).
- Any document where the primary reader is an executive or on-call engineer who will stop reading after 30 seconds.

**Worked example — deliberate delayed lede done well.**

> *"It was 2:47 a.m. on a Tuesday when the third pager went off. By then the on-call had already restarted the cluster twice and was running out of ideas. He didn't know yet that the bug had been in production for eleven months. Here's how we missed it — and what changed after we found it."*

The lede is delayed by three sentences, but each sentence pulls the reader forward. The fourth sentence (the actual nutgraf) is the payoff.

**Worked example — accidental bury (rewrite).**

- Buried: "Q3 has been a busy quarter for the platform team, with several initiatives moving in parallel. Among the items the team has been tracking, the cache subsystem has received particular attention. After extensive analysis, we have concluded that the legacy cache should be deprecated."
- Unburied: "We're deprecating the legacy cache. Cutover by Q4. Migration guide below."

**Diagnostic.** Read paragraph one out loud. Can you state the central point of the piece from paragraph one alone? If no, the lede is buried — unless you're writing a feature, in which case ask: does paragraph one make me want paragraph two? If both fail, restructure.

**References.**

- MasterClass — [Bury the Lede: How to Avoid Burying the Lede in Your Writing](https://www.masterclass.com/articles/bury-the-lede-explained).
- Roy Peter Clark, Poynter Institute — [Lead vs. lede](https://www.poynter.org/reporting-editing/2019/lead-vs-lede-roy-peter-clark-has-the-definitive-answer-at-last/), on terminology and craft history.
- [Lead paragraph — Wikipedia](https://en.wikipedia.org/wiki/Lead_paragraph), straight vs delayed ledes.

---

## Tabular vs prose — when a table beats a paragraph

The decision to render information as a table or as prose is structural, not stylistic. The wrong choice forces the reader to do work the writer should have done.

**Use a table when:**

- You are comparing 3+ items across 2+ attributes. (A 2-item, 2-attribute comparison fits in one sentence: "X is faster; Y is cheaper.")
- The reader's primary action is lookup ("what's the rate limit for the free tier?"), not reading top-to-bottom.
- The items are parallel — they share a structure, and the parallelism itself carries information.
- Exact values matter (numeric thresholds, configuration options, version compatibility, pricing tiers).
- The same shape will be referenced repeatedly (a feature matrix, a permission grid, a tier-by-tier breakdown).

**Use prose when:**

- The information has a single point that doesn't need parallel structure to land.
- The relationships between items are causal or sequential, not parallel.
- The reader is expected to read once, top to bottom (a narrative, an argument, a recommendation).
- The data is fewer than ~4 facts — a sentence is shorter and easier than a 2×2 table.
- The "compare" intent is rhetorical, not analytical ("X is faster, but only if you ignore tail latency").

**Worked example — the same content, two renderings.**

*Prose (when there are only 2 options and 1 dimension):*

> "We considered both the managed service and the self-hosted approach. The managed service is faster to ship but locks us into the vendor's pricing; the self-hosted approach is the inverse."

*Table (when there are 3 options and 4 dimensions):*

| Option | Time to ship | Ongoing cost | Vendor lock-in | Operational burden |
|---|---|---|---|---|
| Managed (vendor A) | 2 weeks | $$$$ | High | None |
| Managed (vendor B) | 3 weeks | $$$ | Medium | Low |
| Self-hosted | 8 weeks | $$ | None | High |

Force-converting between the two breaks both. A four-cell table for two-option/one-dimension content is performative scaffolding. A 12-sentence paragraph for three-option/four-dimension content is unreadable.

**Anti-patterns.**

- *The decorative table.* A table with one column or one row is a list — render it as a list. A table with no parallel structure across rows is prose pretending to be a table.
- *The buried table.* If a comparison matters to the argument, the table goes where the comparison happens, not in an appendix. Tables in appendices are read by nobody.
- *The over-stuffed table.* Tables with more than ~7 columns become illegible; split into two tables, or move detail to footnotes/links.

**When to break it.** Tables in screen-reader-first contexts need extra care; a table that scans visually may be unreadable as a linear stream. For accessibility-heavy outputs, lean toward structured lists with consistent leading text — they degrade better. See `accessibility-writing` for screen-reader-friendly equivalents.

**Diagnostic.** If you can describe the comparison in one sentence without losing precision, use prose. If you find yourself writing "X is A, Y is B, Z is C, and meanwhile…" with parallel verbs across items, you've built a table in prose — promote it.

**References.**

- UC Irvine Writing Center — [The Dos and Don'ts of Using Tables and Figures in Your Writing](https://writingcenter.uci.edu/2024/04/29/the-dos-and-donts-of-using-tables-and-figures-in-your-writing/).
- BCcampus Pressbooks — [Technical Writing Essentials: Figures and Tables](https://pressbooks.bccampus.ca/technicalwriting/chapter/figurestables/).
