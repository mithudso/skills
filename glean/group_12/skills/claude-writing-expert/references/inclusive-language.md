<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `inclusive-language` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: inclusive-language
description: "Bias-free, inclusive language for technical writing, customer communications, internal docs, and code. Covers APA 7th edition bias-free guidelines (race/ethnicity, disability, age, gender, sexual orientation, religion), person-first vs identity-first language, gendered terms, race/ethnicity conventions, disability language, mental health language, age language, LGBTQ+ terminology, technical terminology shifts (blocklist/allowlist, primary/secondary, main), inclusive examples in technical writing, US-centric assumption checking, the lazy-generic test, and code/API naming audits. TRIGGER: user asks 'is this inclusive', needs bias-free language review, person-first vs identity-first guidance, replacing gendered terms, blocklist/allowlist vs blacklist/whitelist, inclusive examples in docs, LGBTQ+ terminology guidance, or a diversity language check on any document or codebase. SKIP: pure style or grammar (use writing-expert); plain language for reading-level reduction; legal or HR policy authoring (consult legal)."
origin: local
version: "1.3.0"
updated: "2026-05-29"
related_skills:
  - writing-expert
  - technical-writing-craft
  - document-critique
---

# Inclusive Language

Reference for bias-free, inclusive language across technical writing, customer communications, internal docs, and code. Draws on APA 7th ed. bias-free language guidelines, AP Stylebook race/ethnicity guidance, SPJ Diversity Toolbox, Google and Microsoft inclusive language style guides, and the NCDJ Disability Language Style Guide.

## When to use this skill

Activate when the user:

- asks whether text, docs, or code comments are inclusive or bias-free
- needs person-first vs identity-first language guidance
- wants to replace gendered defaults (chairman, guys, he/she)
- is reviewing technical terminology (blacklist/whitelist, master/slave)
- needs diverse, non-US-centric examples in technical writing or tutorials
- is localizing content and checking for cultural or regional assumptions
- asks about specific community conventions (e.g., Latine vs Latinx vs Latino, Deaf vs deaf)
- is auditing a codebase, API, or config for non-inclusive naming
- needs LGBTQ+ terminology guidance or is writing about sexual orientation or gender identity
- wants to check religious/spiritual neutrality in content

Skip when:
- the ask is purely style, grammar, or tone — use **writing-expert**
- the ask is reading-level simplification — use **plain-language**
- the content is a formal HR or legal policy — consult legal counsel

---

## Output format

When reviewing a document or codebase for inclusive language:
1. Produce a **find/replace table** — two columns: "Found" and "Suggested replacement", one row per issue.
2. Add a one-sentence rationale for each row that isn't self-evident.
3. If no issues are found, say so explicitly rather than producing an empty table.
4. For ambiguous terms (e.g., "tribe" used literally), flag as informational rather than a required change.

---

## 1. APA 7th Edition Bias-Free Language Principles

Source: *Publication Manual of the American Psychological Association*, 7th ed., §5.1–5.11.

Core principle: describe people at the appropriate level of specificity; acknowledge participation; choose language the people themselves use; avoid stigmatizing constructions.

**Four operating rules:**
1. **Specificity over generality.** "adults aged 65–75" beats "older adults"; "Somali refugees" beats "immigrants".
2. **Participation.** Frame people as active agents: "students who use wheelchairs" not "wheelchair-confined students".
3. **Community preference.** When a community is split, explain both terms and defer to the person's own choice.
4. **Avoid double standards.** If you wouldn't apply a label to one group, do not apply the parallel term to another.

---

## 2. Gender

### Replace gender-default-male terms

| Avoid | Use instead |
|---|---|
| chairman | chair, chairperson |
| manpower | workforce, staffing, personnel |
| mankind | humanity, people, humankind |
| man-hours | person-hours, labor hours |
| spokesman | spokesperson |
| stewardess | flight attendant |
| fireman | firefighter |
| policeman | police officer |
| freshman | first-year student |
| guys (mixed group) | everyone, folks, team, y'all |

### Singular they

Use **singular "they"** for:
- a person whose gender is unknown or irrelevant: "The user submits their request."
- a person who uses they/them pronouns.
- generic references in examples: "Ask your manager — they can approve it."

APA 7th ed. endorses singular they. Chicago 17th ed. also accepts it. Do not use "he/she", "(s)he", or alternating "he" and "she" chapters.

### Avoid assumptions
- Do not assume a doctor, engineer, or CEO is male.
- Do not assume a nurse, teacher, or assistant is female.
- In examples and scenarios, vary pronouns or default to they.

---

## 3. Race and Ethnicity

Source: AP Stylebook race and ethnicity chapter; APA 7th ed. §5.7; SPJ Diversity Toolbox.

**Capitalization:**
- **Black** — capitalize when referring to people of African descent. AP and APA both endorse this (2020). Do not lowercase.
- **White** — AP capitalizes (2020); APA does the same for consistency. Follow the house style your org has chosen; if none exists, capitalize for parallelism with Black.
- **Indigenous, Native American, Alaska Native, Native Hawaiian** — capitalize.
- **Asian, Latino, Hispanic, Arab** — capitalize as ethnic/cultural adjectives.

**Specificity over pan-ethnic labels:**
- Prefer "Nigerian", "Kenyan" over "African" when nationality is known.
- Prefer "Mexican American", "Puerto Rican" over "Hispanic" when ethnicity is known.
- Prefer "Vietnamese American", "Korean" over "Asian" when specific.

**Latino/Latina/Latine/Latinx:**
- *Latino/Latina* — gendered Spanish terms; widely used and understood.
- *Latinx* — gender-neutral; used in some academic and activist contexts; a 2021 Pew survey found only 3% of U.S. Latinos use it personally; generates pushback in some communities.
- *Latine* — gender-neutral alternative growing in use among Spanish speakers in Latin America and Spain.
- **Guidance:** use the term the person or organization uses for themselves. In formal external writing, "Latino" (or "Latina/Latino") is the broadest safe default unless you know the audience prefers another term.

**Avoid "minority" as a catch-all:** it centers a white/majority reference point and aggregates very different groups. Use "people of color", "underrepresented groups", or the specific communities you mean.

**Avoid "non-white":** it defines groups only by what they are not.

---

## 4. Disability Language

Source: NCDJ Disability Language Style Guide; APA 7th ed. §5.11.

**Person-first vs identity-first:**

*Person-first* emphasizes the person before the condition: "person with autism", "person who is blind".
*Identity-first* centers the disability as part of identity: "autistic person", "Deaf person".

There is no universal rule. Communities differ:
- The **Deaf community** strongly prefers identity-first: "Deaf person" (capital D = cultural identity).
- Many in the **autistic community** prefer identity-first: "autistic person".
- Many in the **Down syndrome community** and cancer patient communities prefer person-first.
- When in doubt, follow the individual's preference. When writing about a community, follow that community's dominant convention or acknowledge the split.

**Specific term guidance:**

| Avoid | Use instead |
|---|---|
| suffers from, afflicted with | has, lives with, is diagnosed with |
| wheelchair-bound, confined to a wheelchair | uses a wheelchair, wheelchair user |
| the disabled, the blind, the deaf | people with disabilities, blind people / people who are blind |
| special needs | disability (specific) |
| mentally retarded | intellectual disability |
| crippled (pejorative) | person with a physical disability |
| hearing-impaired | hard of hearing (preferred by many), deaf |
| high-functioning / low-functioning autism | describe specific traits or support needs instead |

---

## 5. Mental Health Language

Source: APA 7th ed. §5.11; Mental Health America style guidance.

**Terms to avoid in all non-clinical contexts:**

| Avoid | Why / Alternative |
|---|---|
| crazy, insane, nuts | stigmatizing; use "unexpected", "difficult", "complex" |
| OCD about X | trivializes the condition; say "meticulous", "detail-oriented" |
| bipolar (meaning inconsistent) | trivializes; say "inconsistent", "unpredictable" |
| psychotic (meaning irrational) | use "irrational", "erratic" |
| schizophrenic (meaning contradictory) | use "contradictory", "inconsistent" |
| committed suicide | died by suicide (removes agency framing; AP Stylebook guidance) |
| successful suicide | died by suicide |

In clinical writing, use APA diagnostic terms accurately and only when clinically appropriate.

---

## 6. Age Language

Source: APA 7th ed. §5.10.

- Avoid "elderly" as a noun or adjective: use "older adults", "people aged 65 and older", or specific age ranges.
- Avoid "seniors" in formal writing (informal in some contexts, fine in others).
- Do not infantilize older adults: avoid "spry", "feisty", "little old lady".
- Do not infantilize young adults: "college-aged adults" not "kids".
- When age is relevant, be specific: "adults aged 18–30" is more useful than "young people".

---

## 7. Technical Terminology

Source: Google Developer Documentation Style Guide; Microsoft Writing Style Guide; IETF RFC 8890; Linux kernel mailing list terminology changes (2020).

### Replace deprecated technical terms

| Avoid | Use instead | Context |
|---|---|---|
| blacklist | blocklist, denylist | security, firewalls, spam filters |
| whitelist | allowlist, safelist, permitlist | security, access control |
| master / slave | primary / replica, leader / follower, primary / secondary, source / sink | databases, Git, I2C, LDAP |
| master branch | main, trunk, default | Git, version control |
| sanity check | confidence check, quick check, smoke test, coherence check | testing, review |
| dummy value | placeholder value, example value | code examples |
| native (meaning simple/fast) | bare-metal, built-in, platform-native | OS/platform docs |

**Notes:**
- "parent/child" remains widely accepted in tree structures; no community-wide push to replace it.
- "abort" vs "cancel": "abort" is still common in POSIX and system programming contexts; use "cancel" in end-user UI copy.

---

## 8. Sexual Orientation, Gender Identity, and Religion

Source: APA 7th ed. §5.8–§5.9; GLAAD Media Reference Guide.

### Sexual orientation and gender identity (SOGI)

- Use **sexual orientation** (not "sexual preference" — "preference" implies choice).
- Use **gender identity** (not "gender preference").
- **Gay, lesbian, bisexual, transgender, queer** — acceptable as adjectives: "gay man", "lesbian couple", "transgender person". Avoid as nouns: "a gay" or "a transgender".
- **Trans** — acceptable shorthand for transgender. Do not say "transgendered" (not a verb form).
- **Nonbinary** — one word; refers to gender identity outside the male/female binary. Use they/them unless the person specifies otherwise.
- **Intersex** — a biological variation; not a gender identity. Do not conflate with transgender.
- Avoid "born a man/woman" or "biologically male/female" in non-medical contexts — use "assigned male/female at birth" (AMAB/AFAB) when the distinction is clinically relevant.
- Do not use "lifestyle" to describe sexual orientation or gender identity — it implies choice.

| Avoid | Use instead |
|---|---|
| sexual preference | sexual orientation |
| transgendered | transgender |
| a transgender | a transgender person |
| homosexual (clinical, pejorative in many contexts) | gay, lesbian, or the person's own term |
| born a woman | assigned female at birth (AFAB), if medically relevant |

### Religion and belief

- Do not assume a default religion in examples, greetings, or holiday references.
- Capitalize names of religions, denominations, and holy texts: Muslim, Christian, Judaism, the Quran, the Bible, Dharma.
- Avoid using religious terms as generic metaphors: "that's my Bible for X" — use "that's my go-to reference for X".
- For holidays: use the specific holiday name when relevant; use "the holiday" or omit when the specific religion is not the point (see §13 for "Christmas Day").
- When listing holidays or time-off examples, use a non-US-centric set: include Eid al-Fitr, Diwali, Lunar New Year alongside Christmas and Hanukkah.

---

## 9. Code and API Naming Audits

Source: Google Developer Documentation Style Guide; Microsoft Writing Style Guide; Linux kernel mailing list (2020).

When auditing a codebase, API, config file, or schema for inclusive naming:

**What to flag:**
- Variable, function, or parameter names using `blacklist`/`whitelist` → rename to `blocklist`/`allowlist` or `denylist`/`safelist`
- Database field names, replica set configs, or replication terms using `master`/`slave` → `primary`/`secondary`, `leader`/`follower`
- Git branch names: `master` → `main`
- Config keys named `dummy_*` → `placeholder_*` or `example_*`
- Comment strings containing any term from §6 (Mental Health Language) or §13 (Common False Friends)

**What not to flag:**
- Third-party library names you do not control (note them as informational)
- Quoted strings that are user-supplied input (flag only if they appear in your own docs/examples)
- Historical commit messages (flag only if actively displayed in UI or docs)

**Suggested workflow:**
```
grep -rn "blacklist\|whitelist\|master\|slave\|dummy\|sanity" \
  --include="*.js" --include="*.ts" --include="*.py" \
  --include="*.yaml" --include="*.json" .
```
Review each match against the tables in §7 before flagging — context determines whether a match is a true positive.

---

## 10. Inclusive Examples in Technical Writing

Source: Google Developer Documentation Style Guide; Microsoft Writing Style Guide.

### Names in examples
- Use a diverse range of names rather than defaulting to Anglo-American male names.
- Avoid: John, Jane, Bob, Alice (fine occasionally, but over-used).
- Recommended rotation: Alex, Sam, Casey, Jordan, Priya, Wei, Amara, Mateo, Ingrid, Daisuke.
- Vary names across examples in the same document; do not assign tech roles only to Western names.

### Personas and roles
- Assign technical roles (developer, sysadmin, architect) to characters of any gender/name.
- Assign non-technical roles (customer, end-user, manager) to characters of any gender/name.
- Do not assign all authoritative roles (admin, architect) to male-coded names.

### Pronouns in examples
- Use **they/them** for any generic user unless a named persona has an established pronoun.
- "The user submits their credentials" — not "his credentials".

---

## 11. US-Centric Assumptions

Common assumptions to remove from global or externally-facing docs:

| Assumption | Fix |
|---|---|
| MM/DD/YYYY date format | Use ISO 8601 (YYYY-MM-DD) or spell out the month |
| US phone format (555) 867-5309 | Use E.164 international format: +15558675309 |
| ZIP code | postal code |
| Social Security Number | government-issued ID number |
| 9-to-5 business hours | specify timezone; say "business hours in your region" |
| Dollar sign $ | specify currency: USD, EUR, etc. |
| English as universal default | note language/locale assumptions explicitly |
| "domestic" meaning US | "in [country]" or "local" |

Localization test: read each example and ask "does this work for a reader in Germany, Japan, and Nigeria?" If not, make it explicit or more general.

---

## 12. The Lazy-Generic Test

Before finalizing any example, scenario, or persona, run this check:

1. **Gender:** Did you default to male pronouns or male-coded names?
2. **Race/ethnicity:** Did you default to Anglo-American names? Are all technical experts in the example white?
3. **Location:** Is this example US-specific when it doesn't need to be?
4. **Ability:** Does the scenario assume sighted, hearing, or able-bodied users without acknowledging accessibility?
5. **Age:** Did you assume a young adult without reason?

If you answer yes to any of these, revise one or more examples.

---

## 13. Common False Friends

| Term | Guidance |
|---|---|
| "Christmas Day" | In external global content, use "December 25" or "the winter holiday" unless writing specifically about Christmas |
| "sanity check" | Debate exists; "confidence check" or "smoke test" is always safe |
| "grandfathered" | Some orgs flag as racially coded; alternatives: "legacy-approved", "pre-existing exception", "retained from prior policy" |
| "cakewalk" | Historical association with slavery; use "easy task", "straightforward" |
| "peanut gallery" | Historically tied to segregation-era seating; use "critics", "the audience", "skeptics" |
| "tribe" | Acceptable in its literal anthropological sense; avoid as a metaphor for work teams |
| "guru", "ninja", "rockstar" | Avoid as job titles or competency descriptors; use specific role names |

---

## 14. Style sheets and inclusive-language term banks

A **style sheet** is a project-specific reference recording every
capitalization, spelling, hyphenation, and terminology decision a team has
made. A **term bank** is the subset of that style sheet that explicitly
governs which terms are approved, which are banned, and which require
context-sensitive replacement. For inclusive language, the term bank is the
load-bearing artifact — without it, the same writer fixes "master/slave" in
one PR and ships it again in the next.

### Why a single document is not enough

A team can adopt every rule in §§1–13 and still ship non-inclusive copy if
the rules live only in memory. The first reviewer enforces them; the
second reviewer (new hire, contractor, summer intern) does not — because
they were never told. A style sheet checked into the repo, the wiki, or the
shared drive makes the inclusive-language rules durable across personnel
turnover.

### Minimum viable inclusive-language term bank

One page, kept in the repo at `/docs/style-sheet.md` or the team wiki:

```
# Inclusive language — approved/banned terms

## Banned (replace on sight)
blacklist          → blocklist, denylist
whitelist          → allowlist, safelist
master (replication, branch) → primary, leader, main
slave              → secondary, follower, replica
dummy value        → placeholder, example
sanity check       → confidence check, smoke test
suffers from       → has, lives with
wheelchair-bound   → uses a wheelchair
the disabled       → people with disabilities
guys (mixed group) → everyone, folks, team, y'all
manpower           → workforce, staffing
chairman           → chair, chairperson
crazy / insane     → unexpected, difficult, complex
committed suicide  → died by suicide
sexual preference  → sexual orientation
transgendered      → transgender

## Context-dependent (review before merging)
abort              → "cancel" in end-user UI; "abort" acceptable in POSIX/system contexts
parent/child       → acceptable in tree/process structures
tribe              → acceptable in literal anthropological sense; avoid as metaphor for teams
grandfathered      → some orgs ban; alternatives: legacy-approved, pre-existing exception
sanity check       → acceptable in some legacy contexts; default to "smoke test"

## Approved capitalizations
Black              (referring to people of African descent) — cap
Indigenous         — cap
Deaf (cultural)    — cap; "deaf" lowercase for the medical condition
Latino, Latina, Latine — cap
LGBTQ+, BIPOC      — all caps
```

### Discipline

- Commit the term bank to the repo. PR reviewers cite specific banned terms
  by reference: "blocked term per style sheet §banned: 'blacklist'."
- Hook the term bank into CI: a grep against the banned column on every PR
  is a cheap, durable enforcement. The §9 audit workflow becomes
  automated.
- Update the term bank as community guidance evolves. Set a quarterly
  review on the calendar to check for changed APA/AP/NCDJ guidance.
- Cross-link with related style sheets. The plain-language and
  writing-expert skills also benefit from a shared style sheet — combine
  inclusive-language entries with capitalization, hyphenation, and tone
  rules so the team has one document, not three.

### Worked example: cross-doc drift the term bank prevents

Without a term bank, a team can ship:
- "blocklist" in the new API docs (because the doc author read this skill).
- "whitelist" in the legacy admin guide (because nobody told the contractor).
- "master branch" in the onboarding doc (because the writer learned Git in 2018).

All three are produced by competent people. The drift comes from the
absence of a single canonical reference. The term bank closes it.

### When to break the rule

- **Quoted user input or external documents** retain their original
  terminology — the term bank governs your prose, not what you quote.
- **Third-party API names you do not control** stay as the vendor named them
  (MySQL replicas are "replicas" in MySQL docs but the underlying API still
  exposes "SHOW SLAVE STATUS" until the vendor renames it). Flag as
  informational rather than a required change.
- **Historical commits and git history** remain as written; rewriting
  history to update terminology causes its own problems. Going forward,
  hold the line.

### Composition with other skills

| Skill | What it adds to the style sheet |
|-------|--------------------------------|
| **writing-expert** | Tone, capitalization, Oxford comma policy, anti-AI-ism Tier 1 list |
| **plain-language** | Reading-grade target, jargon substitutions |
| **technical-writing-craft** | Heading conventions, RFC 2119 keyword usage, verb tense rules |
| **inclusive-language** (this skill) | The banned/approved term bank for bias-free language |

A team's full style sheet is the union of all four. Maintaining them as one
document — not four — reduces friction.

### References

- Microsoft Writing Style Guide — Word lists; A-Z word list.
  `https://learn.microsoft.com/en-us/style-guide/a-z-word-list-term-collections/`
- Google Developer Documentation Style Guide — Word list.
  `https://developers.google.com/style/word-list`
- Chicago Manual of Style §2.55 — style sheets in editorial practice.
- Editorial Freelancers Association — style sheet templates.
- APA Style — Word choice and reducing bias. `https://apastyle.apa.org`

---

## Sources

- APA (2020). *Publication Manual of the American Psychological Association*, 7th ed., §5.1–5.11. https://apastyle.apa.org/style-grammar-guidelines/bias-free-language
- AP Stylebook (2020). Race-related coverage guidance. https://www.apstylebook.com
- SPJ Diversity Toolbox. https://www.spj.org/dtb.asp
- Google Developer Documentation Style Guide — Inclusive documentation. https://developers.google.com/style/inclusive-documentation
- Microsoft Writing Style Guide — Bias-free communication. https://learn.microsoft.com/en-us/style-guide/bias-free-communication
- NCDJ Disability Language Style Guide (2021). https://ncdj.org/style-guide/
- Mental Health America. Language and Mental Health. https://mhanational.org
- GLAAD Media Reference Guide (current edition). https://glaad.org/reference
- IETF RFC 8890 — The Internet is for End Users (terminology note).
- Linux kernel mailing list (2020). Inclusive terminology commit — https://git.kernel.org/
