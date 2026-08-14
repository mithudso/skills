---
description: >-
  Technical & product writing hub — software docs, specs, engineering comms. Sentence mechanics & structure (Given/New, active voice, imperative mood, nominalization/zombie nouns, heading discipline, cohesion/coherence, RFC 2119, Diátaxis, show-first); style-guide grounded (Google, Microsoft WSG, Strunk & White, Williams). TRIGGER: technical writing, API docs, README, readability, passive voice, heading structure; how-to & tutorials; KB articles; reference & explanation docs; runbooks/on-call; specs; PRDs; RFCs & design docs; user stories & acceptance criteria; meeting minutes & ADRs; changelogs & release notes; error messages; microcopy/UI writing; commit & PR messages; code/agent plans; postmortems; incident/status-page comms. SKIP: general prose/voice/editing → writing-expert; exec/business/persuasion → executive-comms; marketing/PR/newsletters → content-and-marketing-writing; career/academic/legal/policy/survey → career-and-formal-writing.
name: technical-writing-craft
origin: local
version: "1.4.0"
updated: "2026-05-30"
keywords:
  - technical-writing-craft
  - Given/New principle
  - nominalization
  - zombie nouns
  - active voice
  - imperative mood
  - sentence rhythm
  - heading discipline
  - cross-reference
  - forward link
  - cohesion
  - coherence
  - RFC 2119
  - Diátaxis
  - show-first
  - verb tense
  - Williams Style
  - Strunk White
  - Google Developer Docs
  - Microsoft Writing Style Guide
  - runbook prose
  - API docs prose
  - how-to writing
  - tutorial writing
  - knowledge-base authoring
  - reference docs
  - explanation docs
  - spec writing
  - PRD
  - RFC
  - design docs
  - user story
  - acceptance criteria
  - meeting minutes
  - decision log
  - changelog
  - release notes
  - error message craft
  - microcopy
  - UI writing
  - commit message
  - PR description
  - code plan
  - agent plan
  - postmortem
  - incident comms
related_skills:
  - writing-expert
  - executive-comms
  - content-and-marketing-writing
  - career-and-formal-writing
whenToUse:
  - "improving prose quality in an API reference, SDK guide, README, runbook, or RFC"
  - "applying active voice, imperative mood, or second-person instructions to technical docs"
  - "fixing zombie nouns / nominalizations in documentation"
  - "applying the Given/New principle to improve sentence flow in a doc"
  - "structuring headings — parallel structure, no orphan H3s, action vs noun headings"
  - "choosing between cross-reference (dependency) and forward link (navigational) in a doc"
  - "applying RFC 2119 normative keywords to a spec or API contract"
  - "auditing or improving prose discipline within a Diátaxis section — tutorial voice, how-to imperative, reference dryness, explanation connectives (to choose which mode to assign to a new page, load references/api-docs-craft.md)"
  - "improving the show:tell ratio in code documentation"
  - "make this doc more readable, clean up my docs, this section is unclear"
---

# Technical Writing Craft

Reference for sentence-level and document-level craft in API docs, runbooks,
RFCs, and README files. Style-guide grounded; distinct from business writing.

Primary sources cited throughout:
- **Williams** — Joseph Williams, *Style: Lessons in Clarity and Grace*, 12th ed.
- **S&W** — Strunk & White, *The Elements of Style*, 4th ed.
- **Google** — Google Developer Documentation Style Guide (developers.google.com/style)
- **MSWG** — Microsoft Writing Style Guide (learn.microsoft.com/style-guide)

Extended reference topics (verb tense, numbers conventions, "show your work",
and four editing types) live in `references/extended-topics.md`.

---

## Sub-skill routing table

This hub absorbs 22 former standalone skills as on-demand reference files. When a task matches a row, **Read the listed `references/` file** before answering — do not rely on this table alone for depth.

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `api-docs-craft` | Author REST/HTTP API documentation that developers can actually integrate against — Diátaxis framework (tutorial / how-to / reference / explanation by… | `references/api-docs-craft.md` |
| `howto-writing` | Problem-oriented step-by-step documentation craft — the Diátaxis how-to quadrant. Write goal-directed recipes for competent users who already know the… | `references/howto-writing.md` |
| `tutorial-writing` | Learning-oriented documentation craft — the Diátaxis tutorial quadrant. Write scaffolded "follow me" lessons that guarantee success at every step, build a… | `references/tutorial-writing.md` |
| `knowledge-base-authoring` | Authoring craft for knowledge-base articles: writing for search discoverability (answer-the-question-in-the-title rule), one-question-per-article… | `references/knowledge-base-authoring.md` |
| `reference-doc-writing` | Information-oriented documentation craft — the Diátaxis reference quadrant. Write technical descriptions of an API, configuration, schema, CLI, or protocol… | `references/reference-doc-writing.md` |
| `explanation-doc-writing` | Understanding-oriented documentation craft — the Diátaxis explanation quadrant. Write the "discussion" docs that build mental models, give background… | `references/explanation-doc-writing.md` |
| `runbook-craft` | Runbook-specific writing craft — execution-safety constraints for documents that must work under pressure on a sleep-deprived on-call engineer at 03:00.… | `references/runbook-craft.md` |
| `spec-writing` | Engineering implementation spec craft — the contract-shaped genre that says WHAT a system must do, distinct from PRDs (which describe what to build at the… | `references/spec-writing.md` |
| `prd-writing` | Product Requirements Document (PRD) craft — the PM-owned genre that defines what to build and why, before engineering proposes how. Covers the "problem /… | `references/prd-writing.md` |
| `rfc-and-design-docs` | RFC and design document craft — Stripe RFCs, Google design docs, Squarespace RFCs, | `references/rfc-and-design-docs.md` |
| `user-story-and-acceptance-criteria` | Author user stories and acceptance criteria that survive grooming and ship as working software — Mike Cohn's "As a / I want / So that" template (with the… | `references/user-story-and-acceptance-criteria.md` |
| `meeting-minutes-and-decision-log` | Capture craft for meeting notes, action items, decision logs, attendee tracking, agenda-vs-minutes distinction, ADR (Architecture Decision Record) form for… | `references/meeting-minutes-and-decision-log.md` |
| `changelog-and-release-notes` | Changelog and release-notes craft — Keep a Changelog spec, semver communication obligations, | `references/changelog-and-release-notes.md` |
| `changelogs-for-humans` | User-facing changelog craft — the end-user counterpart to the developer changelog. Benefit-led language (not feature-list language), screenshots and GIFs… | `references/changelogs-for-humans.md` |
| `error-message-craft` | Voice and clarity for error messages in user-facing software — the "what / why / what to do next" triple, no-blame language, error-code naming conventions… | `references/error-message-craft.md` |
| `microcopy-and-ui-writing` | In-product UX writing for buttons, empty states, tooltips, form-field labels, validation messages, system messages, password-reset flows, onboarding tours… | `references/microcopy-and-ui-writing.md` |
| `commit-message-craft` | Author git commit messages that survive `git blame` years from now — Conventional Commits beyond the basics, Tim Pope's 50/72 rule, "why not what" body… | `references/commit-message-craft.md` |
| `pr-description-craft` | Author pull-request descriptions that get merged fast — the "what / why / how / test" template, GitHub/GitLab/Bitbucket PR templates, the "TL;DR for… | `references/pr-description-craft.md` |
| `code-plan-writing` | Translates specifications, feature requests, and refactors into structured implementation plans for human developers and AI coding agents. Produces… | `references/code-plan-writing.md` |
| `agent-plan-writing` | Write execution plans for AI agent workflows — multi-agent orchestration, | `references/agent-plan-writing.md` |
| `postmortem-writing` | The writing-specific craft for postmortems — disciplined Five Whys in prose, blameless framing language, timeline reconstruction in UTC, contributing… | `references/postmortem-writing.md` |
| `incident-comms` | The writing surface of incident communication — status-page entries, customer-facing incident updates, internal SEV channels, and the cadence/voice shifts… | `references/incident-comms.md` |
| `extended-topics` | Extended depth — verb tense, "show your work" reasoning transparency, and code-doc show-first rules beyond the main file | `references/extended-topics.md` |

---

## 1. Style guides — scope and emphasis

| Guide | Best for | Key emphasis |
|-------|----------|--------------|
| **Google Dev Docs** | API reference, SDK guides, CLI docs | Second-person imperative, present tense, scannable structure |
| **Microsoft WSG** | Product docs, UI strings, procedures | Warm-but-professional tone, accessibility, global readability |
| **Strunk & White** | Any prose; universal sentence discipline | Omit needless words; prefer the active voice |
| **Williams** | Clarity at sentence and paragraph level | Given/New, cohesion, nominalization, stress position |

Do not treat any single guide as sovereign. Use Google and Microsoft for
structural and style decisions; use Williams and S&W for sentence mechanics.

---

## 2. Williams' Given/New principle

*Source: Williams, ch. 4 "Emphasis"*

Every sentence has two jobs:

- **Subject (topic position)** — carries what the reader already knows or has
  just read. This is *given* information.
- **Predicate (stress position)** — carries what is new, surprising, or the
  point. This is *new* information.

Violating Given/New forces readers to reparse. They start a sentence expecting
to orient on something familiar, but the subject introduces a concept they
haven't seen yet.

**Broken Given/New:**
> The retry logic is governed by a configurable backoff policy. Three fields
> control the policy: `initialDelay`, `maxDelay`, and `multiplier`.

The second sentence begins with "Three fields" — new information in topic
position — before the reader has the policy structure in mind.

**Repaired:**
> The retry logic is governed by a configurable backoff policy. The policy
> exposes three fields: `initialDelay`, `maxDelay`, and `multiplier`.

Now "The policy" links back to the previous sentence's stress position.

**Coherence test (Williams):** read only the opening words of each sentence in
a paragraph. They should form a chain — each one echoes or extends the last.
If the chain breaks, the paragraph loses coherence regardless of how correct
each sentence is in isolation.

---

## 3. Sentence rhythm

Good technical prose is not uniformly short. It alternates.

**Long–short–long pattern:**
A longer sentence with full context, then a short punchy anchor, then a longer
elaboration. The short sentence provides cognitive relief and marks the landing
point. *S&W §12*: "Choose a suitable design and hold to it."

**The comma test:**
Read the sentence aloud. Every pause where you want to breathe is a candidate
comma. If the pauses outnumber three, split the sentence. A sentence requiring
four tracked clauses costs the reader more than two sentences cost.

**Parsing-load reduction (Williams, ch. 4):**
Readers parse faster when subject and verb are close together. Distance between
subject and verb is the single biggest source of hard-to-read technical prose.

- Broken: "The request, after authentication via the token endpoint and after
  the rate-limit bucket has been checked, is forwarded to the backend."
- Repaired: "After authentication and rate-limit checks, the request is
  forwarded to the backend."

Move long modifiers before or after the main clause, not between subject and verb.

---

## 4. Sentence diagramming for clarity edits

When a sentence resists quick repair, parse it into its grammatical skeleton to
find the buried verb. The five-slot diagram:

```
[Agent] + [Verb] + [Object] + [Modifier] + [Modifier]
```

Steps:
1. Find the main verb. Strip `to-be` forms (`is`, `was`, `are`) and
   nominalized verbs first — the real action is usually hiding inside a noun.
2. Find the grammatical subject. Ask: is this the *agent* — the entity
   performing the action?
3. If the subject is not the agent, locate the buried agent (often inside a
   prepositional phrase or a nominalization) and promote it to subject.
4. Rebuild the sentence: `[Agent] [Verb] [Object]`, then reattach modifiers.

**Example — applying the five-slot diagram:**

Original: "The implementation of rate limiting provides protection of the API
against abuse."

Slot analysis:
- Grammatical subject: "The implementation" — not the agent; it's a nominalized verb
- Main verb: "provides" — weak; the real action is buried in "protection"
- Buried agent: "rate limiting" (inside the subject nominalization)
- Buried verb 1: "implement" → nominalized as "implementation"
- Buried verb 2: "protect" → nominalized as "protection"

Rebuilt: `[Rate limiting] [protects] [the API] [from abuse]`

Result: "Rate limiting protects the API from abuse."

Two nominalizations eliminated, sentence length halved. *S&W §14*: "Use the
active voice. Do not use the passive where you can use the active."

---

## 5. Active voice — the rule and the three exceptions

**Default rule:** active voice (*S&W §14*, *Google: Active voice*).
Subject performs the action. The agent is visible and accountable.

> **Active:** The SDK throws `AuthError` when the token expires.
> **Passive:** `AuthError` is thrown when the token expires.

The active version names what throws — the SDK. That matters to the reader.

**The three cases where passive is correct:**

1. **No actor exists.** The process completes without a human or system
   agent. "The index is built during cluster startup." Nothing performs the
   build; it is a state transition.

2. **Actor is not relevant.** The identity of the actor adds no value and
   distracts from the point. In a runbook step: "The payload is validated
   before the request is accepted" — who validates it is the server internals;
   the operator does not act on that information.

3. **Focus should be on the receiver.** When the receiver of the action is
   the topic, passive keeps it in subject position and honors Given/New.
   "The config file is read on startup" — if the paragraph is about the config
   file, passive keeps it correctly in topic position.

Outside these three cases, default to active. *Google Style Guide*: "Use
active voice for most sentences." *MSWG*: "Passive voice makes it harder to
identify who or what performs an action."

---

## 6. Imperative mood and second person

*Source: Google Style Guide — "Write in second person"; MSWG — "Address the reader" and "You"*

Instructions must use the imperative. No softening, no optionals, no
description-as-instruction. Both Google and Microsoft recommend second-person
("you") over third-person ("the user", "the developer") — second person is
direct and reduces the gap between instruction and actor. However, "you" +
indicative mood softens instructions into suggestions. Use "you" + imperative
or the bare imperative for required steps; reserve "you can" for genuinely
optional actions.

| Pattern | Appropriate use |
|---|---|
| Imperative ("Run X.") | Default for all required procedural steps |
| "You can run X." | Genuinely optional; the reader may skip this |
| "You should run X." | Avoid — "should" implies optional when "must" is meant |
| "It is possible to configure retries." | Avoid — use "Configure retries by setting…" |
| "Users may wish to enable debug mode." | Avoid — use "Enable debug mode by…" |
| "The user runs X." | Avoid — third person distances the instruction |

The soft variants ("you can", "you may", "it is possible to") imply the action
is optional. Use them only when it genuinely is.

**Present tense for state descriptions:** describe what the system does, not
what it will do. "The function returns a cursor" not "The function will return
a cursor." *Google Style Guide*: "Use present tense."

---

## 7. Show, don't tell — example before explanation

Code docs suffer from the opposite of narrative fiction's problem: they explain
and then, maybe, show. Reverse the order.

**Tell-then-show (weaker):**
> The `retry` option controls how many times the client retries a failed
> request before returning an error. It accepts an integer.
> ```js
> client = new Client({ retry: 3 });
> ```

**Show-then-tell (stronger):**
> ```js
> client = new Client({ retry: 3 });
> ```
> `retry` sets the maximum number of retries before the client returns an error.
> Defaults to `0` (no retries).

The example anchors the reader's mental model before the explanation fills in
the edges. This mirrors how developers actually read docs: they scan to an
example first, then read the surrounding prose only if the example doesn't
answer the question. *Google Style Guide*: "When documenting a function, show
an example first."

**Rule A — minimum viable example first.** Every concept's documentation
should open with the *shortest runnable* example that uses the concept.

**Rule B — show outputs, not just code.** If the example produces a result,
show the result inline:

```js
const users = await client.users.list({ active: true });
// → [{ id: "u_1", email: "ada@example.com", active: true }, ...]
```

**Rule C — show errors with their resolution.** Pair a malformed request with
the error response and the corrected request. The reader sees both the failure
and the fix.

**Diagnostic.** Open the doc and read only the code blocks. If you can
complete the task from the code alone, the show:tell ratio is right. If the
code is incomprehensible without prose, you have under-invested in the show side.

**References.**
- Google Developer Documentation Style Guide — "Code samples." https://developers.google.com/style/code-samples
- Daniele Procida, "What nobody tells you about documentation" — show-first applies most strongly in tutorials and how-tos.

---

## 8. Nominalization — zombie nouns

*Source: Williams, ch. 3 "Actions"; S&W §13*

Nominalizations are verbs or adjectives converted to nouns. Williams calls them
nominalizations; the colloquial term "zombie nouns" (Helen Sword, *Stylish Academic
Writing*, 2012) captures the same idea — nouns that drain the life out of prose
by burying the action.

| Zombie noun | The verb hiding inside |
|---|---|
| "make a calculation" | calculate |
| "perform an analysis" | analyze |
| "provide an indication" | indicate |
| "give consideration to" | consider |
| "conduct an investigation" | investigate |
| "offer a description of" | describe |

**Rule:** when you write "make/give/provide/perform/conduct/offer + noun",
replace the whole phrase with the verb.

- "The logger performs a flush of the buffer on shutdown." → "The logger flushes
  the buffer on shutdown."
- "This section provides a description of the authentication flow." → "This
  section describes the authentication flow." (Better: cut the meta-sentence
  and just describe it.)

The "provides a description of" pattern is especially common in tech docs because
writers mistake it for formal register. It is not formal; it is wordless.

---

## 9. Throat-clearing intros — cut them

*Source: Williams, ch. 5 "Cohesion and Coherence"; S&W §13 "Omit needless words"*

Throat-clearing is a sentence (or paragraph) that announces what is about to
be said instead of saying it.

**Throat-clearing:**
> "This document covers the architecture of the authentication system. It will
> explain the components, their relationships, and how they interact."

**What to do instead:** start with the first real sentence of content.
> "The authentication system has three components: the token endpoint, the
> session cache, and the refresh daemon."

Opening meta-announcements delay the point, inflate the word count, and train
readers to skip the first paragraph. S&W §13: "Omit needless words."

---

## 10. Heading discipline

*Source: Google Style Guide — "Headings"; MSWG — "Headings and titles"*

**Every heading earns its place.** A heading that could be deleted without
affecting navigation or comprehension should be deleted.

**No orphan H3s.** An H3 must not appear without a preceding H2 in the same
section. An H3 alone under an H2 with no sibling H3s is a sign the H3 is
either the H2's content (promote it) or the H2 is unnecessary (flatten it).

**Headings are not sentences.** Avoid terminal punctuation except question marks
on FAQ entries. "Authentication Flow" not "Authentication Flow." — drop the period.

**Parallel structure across siblings.** Sibling headings at the same level
should follow the same grammatical pattern:
- Good: "Configure TLS", "Enable Logging", "Set Retry Policy" (all imperative)
- Bad: "Configuring TLS", "Enable Logging", "How to Set Retry Policy"

**Action headings for procedures, noun headings for reference.** Procedural
sections use verb phrases ("Install the SDK", "Run the migration"). Reference
sections use noun phrases ("Configuration Options", "Error Codes").

**The heading test:** read just the headings as a list. They should outline
the document's argument. If the list is incoherent or redundant, the document's
structure is broken regardless of how good the body prose is.

---

## 11. Cross-references vs forward links

These are not interchangeable. Choose deliberately.

| Type | Definition | Use when |
|---|---|---|
| **Cross-reference** | Points to content elsewhere that the reader needs *now* to complete the current task | The current section depends on that content; omitting the link leaves a gap |
| **Forward link** | Points to content the reader will need *later*, after completing the current task | The current section is complete on its own; the link is navigational aid |

**Cross-reference** example (required, not optional):
> "Before running the migration, configure your connection string
> (see [Connection Strings](../connection-strings.md))."

**Forward link** example (navigational only):
> "For advanced retry configuration, see [Retry Policy](../retry-policy.md)."

**Anti-pattern:** using "see also" as a dump for every tangentially related
page. Audit every cross-reference: is it a dependency or a suggestion?

---

## 12. Cohesion vs coherence

*Source: Williams, ch. 5 "Cohesion and Coherence"*

**Cohesion** is sentence-to-sentence flow within a paragraph. Each sentence
connects back to the previous one via shared vocabulary, pronoun reference, or
topic continuity.

**Coherence** is the whole-document argument — the logical progression from
premise to conclusion, or from problem to solution.

You need both. A document can be:
- **Cohesive but incoherent:** smooth prose leading nowhere. Common in
  tutorial introductions that flow nicely but have no argument.
- **Coherent but not cohesive:** a correct structure with jarring transitions.
  Reads like a bullet list turned into paragraphs.

**Cohesion check:** in each paragraph, does every sentence link back to the
previous one? (Given/New principle, §2 above.)

**Coherence check:** read only the first sentence of each section. Do they
sequence into a complete argument?

---

## 13. RFC 2119 normative keywords in technical docs

When an API reference, SDK guide, or protocol spec is going to be consumed by
implementers who must conform to it, lift the normative-keyword convention from
IETF RFC 2119 (1997) as clarified by RFC 8174 (2017).

### The rule

Only the **uppercase** forms (MUST, MUST NOT, SHOULD, SHOULD NOT, MAY) carry
conformance meaning. Lowercase versions in the same document are ordinary prose.
Declare the convention in a Conventions/Notation section:

> The key words MUST, MUST NOT, SHOULD, SHOULD NOT, and MAY in this document
> are to be interpreted as described in RFC 2119 and RFC 8174 when, and only
> when, they appear in all capitals.

### Worked example (API reference)

Ambiguous:
> The client should send the `Authorization` header on every request.

Conformance-testable:
> Clients **MUST** include an `Authorization: Bearer <token>` header on every
> request to `/v2/*` endpoints. Clients **SHOULD** retry on `503` responses
> using exponential backoff. Clients **MAY** cache the OAuth token until five
> minutes before its `exp` claim.

A QA engineer can write a test against each MUST/SHOULD/MAY line.

### When to break the rule

- README files and tutorials: uppercase keywords read as shouting in conversational prose.
- Internal team runbooks: team review usually resolves ambiguity without RFC 2119.
- Use RFC 2119 when the audience is external implementers or conformance testing depends on the wording.

### References

- `https://www.rfc-editor.org/rfc/rfc2119` — RFC 2119, Bradner 1997.
- `https://www.rfc-editor.org/rfc/rfc8174` — RFC 8174, Leiba 2017 (uppercase clarification).

---

## 14. Diátaxis: tutorial vs how-to vs reference vs explanation

The Diátaxis framework (Daniele Procida) categorizes technical documentation
into four modes, each with a different reader posture and a different writing
discipline. Mixing modes within a single document is the most common cause of
docs that are simultaneously too long and unhelpful.

For questions about *which* Diátaxis mode to apply when structuring a documentation site, use `api-docs-craft`. This section covers writing discipline *within* each mode.

### The four modes

| Mode | Reader posture | Author posture | What it answers |
|------|---------------|----------------|-----------------|
| **Tutorial** | Learning by doing, hand-held | Trusted teacher | "Take me through my first success" |
| **How-to guide** | Working on a specific task, knows the basics | Guide who shows the path | "How do I accomplish X?" |
| **Reference** | Looking up a specific fact | Authoritative catalog | "What is the exact signature/spec of X?" |
| **Explanation** | Trying to understand, not doing | Tutor, contextualizer | "Why does it work this way?" |

### The discipline

- **Tutorials** are end-to-end, ordered, and guaranteed to work as written. They do not branch. They use the imperative. They teach by producing a result.
- **How-to guides** assume the reader has goals and basic competence. They are goal-shaped ("How to enable TLS"). They can branch on conditions. They omit what is irrelevant to the named task.
- **Reference** docs are accurate, complete, and dry. They follow a consistent template (parameters, return value, errors, example). They do not teach; they catalog.
- **Explanation** docs discuss design, history, alternatives, and trade-offs. They use connectives like "because," "therefore," and "in contrast." They do not include step-by-step procedures.

### Worked example

A single section titled "Configuring authentication" that mixes all four modes overloads the reader.

Diátaxis split:
- **Tutorial:** "Your first authenticated request" (15 minutes, guaranteed result).
- **How-to:** "How to rotate an OAuth client secret."
- **Reference:** "Authentication configuration options" (table of every parameter).
- **Explanation:** "Why we chose OAuth 2.0 over API keys."

### When to break the rule

- Very small projects can keep one README that mixes modes — the overhead of four files outweighs the structural benefit.
- An onboarding doc may legitimately interleave tutorial and explanation in a controlled way.

### References

- `https://diataxis.fr` — Procida, D. *The Documentation System*.

---

## 15. Common anti-patterns quick-reference

| Anti-pattern | Rule violated | Repair |
|---|---|---|
| Nominalization ("provide an indication") | Williams ch. 3 | Replace with the verb: "indicate" |
| Throat-clearing intro | S&W §13, Williams ch. 5 | Delete; start with the first content sentence |
| Passive hiding actor | S&W §14 | Promote the agent to subject |
| "You can run X" for a required step | Google Style, MSWG | Imperative: "Run X." |
| Explain-then-example order in code docs | Google Style | Example first, then explanation |
| Orphan H3 | Google Style | Flatten or add sibling H3s |
| Subject–verb separation by long modifier | Williams ch. 4 | Move modifier to front or end |
| Uniform sentence length | Williams ch. 4 | Vary: long–short–long |
| "See also" link dump | Google Style | Audit: cross-reference (dependency) or forward link (navigational)? |
| "This document will explain…" meta-announcement | S&W §13 | Delete; write the thing |
| Broken Given/New (new info in topic position) | Williams ch. 4 | Restructure so subjects carry prior context |
| Mixed verb tenses without purpose | Google Style | Present for behavior, past for history, future for commitments |

---

*Extended reference: verb tense discipline, numbers conventions, "show your work" transparency, and four editing type definitions → `references/extended-topics.md`*

---

## Appendix: The nutgraf in technical writing

The *nutgraf* (or nut graph) is a journalism convention — the paragraph after the opening that tells the reader what the piece is about and why it matters. Technical writers should adopt it, with one adaptation: in API docs, runbooks, and RFCs the nutgraf isn't about *timeliness*, it's about *scope and prerequisites*. The reader needs to know what they'll learn, what they need to know first, and whether this is the right page for their question.

**The technical-writing nutgraf scaffold (3–4 sentences, no more):**

1. **What this page covers.** One sentence stating the scope.
2. **Who it's for.** The audience filter — language, role, expertise level, prerequisites.
3. **What success looks like.** What the reader will be able to do after reading it (tutorials/how-tos) or what they will understand (explanation/reference).
4. **What this page is *not*.** A single sentence with a forward-link to the adjacent page. Optional but high-leverage — saves readers from skimming a page that doesn't match their intent.

**Worked example — README opening:**

- Without a nutgraf: "Welcome to project X! Project X is a tool that helps you manage configuration. Read on for installation instructions."
- With a nutgraf: "**This README is the quickstart for project X.** It assumes you have Node 20+ installed and a working npm. By the end of the next 5 minutes you will have a running local server. For production deployment, see DEPLOY.md. For architecture, see ARCHITECTURE.md."

**Worked example — runbook opening:**

- Without a nutgraf: "When the deploy fails, follow these steps."
- With a nutgraf: "**This runbook covers recovery from a failed production deploy.** It is for on-call engineers with deploy permissions. By the end you will have either rolled forward to a green deploy or rolled back to the previous version. For partial deploy degradation (some pods green, some red), see runbooks/partial-deploy.md instead."

**Why "what this is not" earns its line.** The single biggest source of wasted reader time in technical docs is reading the wrong page. A one-sentence forward-link in the nutgraf saves readers two minutes of skimming and a frustrated trip back to the index.

**When to skip.** One-line reference entries (a single endpoint description, a single env-var entry) don't need a nutgraf — the heading does the work. Anything longer than ~150 words needs one.

**Diagnostic.** Read your opening. Can the reader answer "is this the right page for me?" without scrolling? If no, the nutgraf is missing or doing the wrong job. See also `writing-expert` for the general-purpose nutgraf framework.

**References.**

- [Nut graph — Wikipedia](https://en.wikipedia.org/wiki/Nut_graph) — origin and canonical definition.
- Daniele Procida — *What nobody tells you about documentation* — the tutorial/how-to/explanation/reference quadrant; nutgraf intent shifts by quadrant.
- Google Developer Documentation Style Guide — opening paragraph guidance.

---

## Appendix: Tabular vs prose in technical writing

In technical docs the table-or-prose decision is sharper than in business writing because *lookup* is the dominant reader behavior — most API doc readers arrived via search and need to find one fact fast.

**Always use a table when:**

- Listing endpoints, methods, status codes, or error codes with parallel attributes.
- Documenting configuration options (key, type, default, description).
- Showing version compatibility, feature support across platforms, or tier-by-tier pricing.
- Documenting permission/role matrices.

**Always use prose when:**

- Explaining *why* a behavior is the way it is (architecture decisions, design rationale).
- Walking through a sequence of operations (tutorials, how-tos — use ordered lists, not tables).
- Stating a single fact ("This endpoint is idempotent.") — a one-row table is a bug.

**Three reliable table-shapes for technical docs.**

1. **The option matrix.**

   | Option | Type | Default | Description |
   |---|---|---|---|
   | `timeout_ms` | integer | `30000` | Request timeout in milliseconds. |
   | `retry_max` | integer | `3` | Maximum retry attempts on 5xx errors. |

2. **The error-code grid.**

   | Code | Meaning | Recommended action |
   |---|---|---|
   | `409 CONFLICT` | The resource was modified concurrently. | Re-fetch and retry with the new ETag. |

3. **The compatibility matrix.**

   | Driver version | Server 6.x | Server 7.x | Server 8.x |
   |---|---|---|---|
   | v4.x | yes | yes | no |
   | v5.x | no | yes | yes |

**Technical-writing-specific anti-patterns.**

- *Tables of code samples.* Code samples belong in fenced code blocks, not table cells — line-wrapping in tables breaks code readability.
- *Tables that hide important detail in footnotes.* If a footnote is load-bearing for a row, the row belongs in prose.
- *Sortable-looking tables in static docs.* If a table appears to invite sorting but is rendered as plain HTML/markdown, document the default sort or make it actually sortable.

**Diagnostic.** Run a screen reader against the table. If the reading order doesn't make sense as a linear stream of "row 1, column 1: X. Row 1, column 2: Y," restructure as a definition list or structured list. See `accessibility-writing` for screen-reader-friendly equivalents.

**References.**

- BCcampus Pressbooks — [Technical Writing Essentials: Figures and Tables](https://pressbooks.bccampus.ca/technicalwriting/chapter/figurestables/).
- Microsoft Writing Style Guide — Tables (accessibility and column-count guidance).
- Google Developer Documentation Style Guide — Tables.

<!-- cross-hub-map -->
## Cross-hub map — where every writing topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `writing-expert` | Writing Expert (prose craft, voice, style, editing) — hub | `references/editing-and-revision.md`, `references/rhetorical-frameworks-deep.md`, `references/storytelling-and-narrative.md`, `references/headline-craft.md`, … |
| `technical-writing-craft` | Technical & Product Writing (docs, specs, engineering comms) — hub | `references/api-docs-craft.md`, `references/howto-writing.md`, `references/tutorial-writing.md`, `references/knowledge-base-authoring.md`, … |
| `executive-comms` | Executive & Business Communication (leadership, persuasion, decks) — hub | `references/one-pager-writing.md`, `references/okr-writing.md`, `references/proposal-and-grant-writing.md`, `references/pitch-deck-writing.md`, … |
| `content-and-marketing-writing` | Content, Marketing & External Comms (PR, newsletters, launch, social) | `references/sales-and-marketing-copy.md`, `references/press-release-writing.md`, `references/crisis-pr-writing.md`, `references/newsletter-writing.md`, … |
| `career-and-formal-writing` | Career, Academic, Legal & Formal Writing | `references/resume-and-cv-writing.md`, `references/cover-letter-writing.md`, `references/job-description-writing.md`, `references/performance-review-writing.md`, … |
