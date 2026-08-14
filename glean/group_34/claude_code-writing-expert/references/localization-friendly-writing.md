<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `localization-friendly-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: localization-friendly-writing
version: "1.0.0"
updated: "2026-05-29"
category: custom
tags:
  - writing
  - localization
  - i18n
  - l10n
  - translation
  - icu-messageformat
description: >
  Writing source-language text (typically English) that will translate cleanly
  into 30+ locales. Covers the "translation-friendly English" rules (short
  sentences, no idioms, no metaphors, no cultural references, predictable
  subject-verb-object structure), ICU MessageFormat syntax (plural, select,
  selectordinal, date/number formatters), CLDR plural categories
  (zero/one/two/few/many/other), placeholders that survive translation (named,
  not positional; never concatenated), translator-comment discipline,
  pseudo-localization testing (30–40% expansion, accented Latin, mirrored RTL),
  avoiding text-in-images, RTL-friendly layouts, and key-naming conventions for
  i18n string files. References Unicode CLDR, ICU User Guide, Mozilla l10n
  developer guidelines, W3C i18n WG, and i18next best practices.
  TRIGGER: "translation friendly", "localizable string", "i18n", "l10n",
  "ICU MessageFormat", "plural rules", "CLDR", "pseudo-localization",
  "pseudolocale", "RTL", "right-to-left", "string concatenation", "placeholder
  in translation", "string key naming", "translator comment", "Fluent",
  "gettext", "po file", "MessageFormat 2", "selectordinal", "string expansion",
  "untranslatable", "won't translate", "translating into German Japanese
  Arabic", "internationalization".
  SKIP: plain-language simplification for English-only audiences (use
  plain-language); accessibility text alternatives (use visual-writing or
  accessibility-ux-reviewer); writing the underlying prose for tone (use
  writing-expert); time-zone / currency / number formatting at the engineering
  layer with no source-text concern (use api-design-patterns).
triggers:
  - "translation friendly"
  - "localizable string"
  - "i18n"
  - "l10n"
  - "ICU MessageFormat"
  - "plural rules"
  - "CLDR"
  - "pseudo-localization"
  - "pseudolocale"
  - "RTL"
  - "right-to-left"
  - "string concatenation"
  - "placeholder"
  - "string key naming"
  - "translator comment"
  - "Fluent"
  - "gettext"
  - "MessageFormat 2"
  - "selectordinal"
  - "string expansion"
  - "internationalization"
skip:
  - plain-language for English-only audiences → use plain-language
  - alt text or accessible chart captions → use visual-writing
  - general prose tone or voice → use writing-expert
  - engineering-side timezone / currency / number formatting → use api-design-patterns
related:
  - plain-language
  - writing-expert
  - visual-writing
  - inclusive-language
  - frontend-design
---

# Localization-Friendly Writing

Reference for writing source-language strings that translate cleanly into 30+
locales. The work happens before the translator sees the file: a string that
respects ICU MessageFormat, CLDR plural rules, placeholder discipline, and
the "no idioms" rule costs the translator minutes; a string that violates
any of those costs hours and breaks at runtime.

Deliver all responses in a direct, plain register. Avoid hedging and meta-commentary.

---

## When to use this skill

Activate when the user:

- Writes UI strings, copy, or notifications that will be translated
- Asks about ICU MessageFormat, plurals, gender selection, or ordinals
- Needs CLDR plural-category coverage (zero/one/two/few/many/other)
- Has a string that "doesn't translate" and wants it diagnosed
- Wants pseudo-localization applied as a smoke test
- Needs translator comments or context fields
- Is naming keys for an i18n string file (`messages.en.json`,
  `strings.po`, `.ftl`, `.resw`)
- Is auditing strings before sending to a translation vendor
- Asks about RTL-friendly layout choices in copy (date formats, mirrored
  punctuation, neutral phrasing)
- Wants the "translation-friendly English" rewrite of an existing block

Skip when:

- The audience is monolingual English; the goal is reading-level (use
  `plain-language`)
- The task is alt text or chart captions (use `visual-writing`)
- The task is general prose tone (use `writing-expert`)
- The task is engineering-side date/timezone/currency formatting with no
  source-text bearing (use `api-design-patterns`)

---

## The one rule: write so a translator can reorder, expand, and replace

Every translation operation needs three freedoms:

1. **Reorder** — subject-verb-object in English is not subject-verb-object in
   Japanese or German. The translator must move words.
2. **Expand** — German, Russian, Finnish run 30–40% longer than English.
   French and Spanish run 15–25% longer. The UI must accommodate.
3. **Replace** — Plural forms, gendered forms, formal/informal address — the
   translator must select among multiple forms based on runtime values.

Every rule in this skill exists to preserve those three freedoms.

---

## Core concept 1 — The translation-friendly English rules

10 rules. Each is a hard rule; violating one will break at least one locale.

1. **One sentence, one idea.** Compound sentences with subordinate clauses
   become unparseable in OV languages.
2. **Subject-verb-object, in that order.** Avoid sentence-initial adverbial
   phrases and dangling participles.
3. **No idioms.** "Hit the ground running" has no German equivalent; the
   translator will invent a literal version that sounds nonsensical.
4. **No metaphors.** "Move the needle" requires a needle, which requires a
   gauge, which requires the metaphor to land. It doesn't, in 27 of 30
   locales.
5. **No cultural references.** No baseball, no Thanksgiving, no "as American
   as apple pie," no Marvel cinematic universe.
6. **Avoid humour and wordplay.** Puns are untranslatable by definition.
7. **No abbreviations the reader must decode.** "Q1," "EOY," "ASAP" — spell
   them out at first use.
8. **No phrasal verbs where a single verb works.** "Set up the account"
   becomes "create the account."
9. **No latinate jargon.** "Utilize" → "use." "Initiate" → "start."
10. **Active voice as the default.** "The system saves the file" not "The
    file is saved by the system." Some languages mark passives in a way
    that breaks UI placement.

---

## Core concept 2 — ICU MessageFormat: plural

ICU MessageFormat is the industry standard for runtime string selection
based on numeric value. The plural argument selects on CLDR plural category.

```icu
{count, plural,
  =0 {No items}
  one {# item}
  other {# items}
}
```

Reading this:

- `count` — the variable
- `plural` — the selector type
- `=0`, `=1` — explicit literal matches (used for "no items" or "one item"
  special cases)
- `one`, `other` — CLDR plural categories
- `#` — the value of `count`, formatted with the locale's number format

**Common mistake — only English plural categories:**

```icu
{count, plural,
  one {# item}
  other {# items}
}
```

This works for English. It silently breaks Russian, Polish, Arabic. Russian
needs `one`, `few`, `many`, `other`. Arabic needs all six.

---

## Core concept 3 — CLDR plural categories

The full set of CLDR plural categories, per Unicode CLDR:

| Category | Used by example languages | Example values |
|----------|----------------------------|----------------|
| `zero` | Arabic, Welsh, Latvian | 0 (in some languages) |
| `one` | English, Spanish, French | 1 |
| `two` | Arabic, Welsh | 2 |
| `few` | Russian, Polish, Czech | 2–4 (Russian: 2, 3, 4, 22, 23, 24…) |
| `many` | Russian, Polish, Welsh | 5+ in some languages, plus large numbers in French |
| `other` | every language; required fallback | Everything else |

Critical facts:

- **`other` is required.** Every plural block must include `other`. It is the
  CLDR-mandated fallback.
- **English uses only `one` and `other`.** Don't assume that means others do.
- **Japanese, Korean, Chinese, Thai use only `other`.** No plural marking.
- **Arabic uses all six** — `zero`, `one`, `two`, `few`, `many`, `other`.
- **The `=N` literal match overrides categories.** Use it for special-case
  copy ("No messages" rather than "0 messages") without breaking pluralization
  elsewhere.

The source-language author writes only the English categories. CLDR-aware
translation tooling expands the form to whatever the target locale requires.

---

## Core concept 4 — ICU select and selectordinal

**`select`** is a switch over a string variable, typically used for gender:

```icu
{gender, select,
  female {She updated her profile.}
  male {He updated his profile.}
  other {They updated their profile.}
}
```

`other` is required, even in `select`. Treat unknown / non-binary / unset
gender by routing to `other`, not to one of the binary branches.

**`selectordinal`** is for ordinal forms (1st, 2nd, 3rd):

```icu
{place, selectordinal,
  one {#st place}
  two {#nd place}
  few {#rd place}
  other {#th place}
}
```

English ordinals use `one`/`two`/`few`/`other`. Russian ordinals use only
`other`. Welsh uses six categories.

---

## Core concept 5 — Placeholders that survive translation

Three rules:

1. **Use named placeholders, not positional.** `{username}` survives word
   reorder. `%s %s` does not — the translator cannot tell which is which.
2. **Never concatenate.** `"Hello, " + username + "!"` forces English word
   order on every locale. Instead: `"Hello, {username}!"`.
3. **Always provide a comment describing the placeholder.** A translator who
   sees `{count}` does not know if it is items, dollars, percent, or seconds.

**Concatenation anti-pattern:**

```javascript
// BAD
const msg = t('error.prefix') + ' ' + filename + ' ' + t('error.suffix');

// GOOD
const msg = t('error.full', { filename });
// strings.en.json: { "error.full": "Could not save file {filename}." }
```

**Rich-text interpolation (React, Vue, etc.):** use the framework's
placeholder-as-component pattern, not string concatenation of JSX:

```jsx
// BAD
<>Click <a href="…">here</a> to continue</>

// GOOD (i18next Trans component)
<Trans i18nKey="continue.prompt">
  Click <a href={url}>here</a> to continue.
</Trans>
// strings.en.json: { "continue.prompt": "Click <1>here</1> to continue." }
```

The translator gets one whole sentence and decides where the link goes.

---

## Core concept 6 — Translator comments

Every non-trivial string gets a translator comment. The comment answers
three questions:

1. **What is this?** UI element type (button, error, tooltip, heading).
2. **What does the placeholder mean?** `{count}` = unread messages, integer
   ≥ 0.
3. **Where does it appear?** Modal title? Inline help? Toast?

**ICU MessageFormat with translator comment (i18next syntax):**

```json
{
  "inbox.unread": {
    "message": "{count, plural, =0 {No unread messages} one {# unread message} other {# unread messages}}",
    "description": "Inbox header label, top of the messages page. {count} is the integer number of unread messages, 0 or greater."
  }
}
```

**Fluent (`.ftl`) with comment:**

```fluent
# Inbox header label. Shown at the top of the messages page.
# $count (Number): The integer number of unread messages.
inbox-unread = {$count ->
    [0] No unread messages
   *[other] {$count} unread {$count ->
       [one] message
      *[other] messages
   }
}
```

**gettext (`.po`) with comment:**

```po
#. Inbox header label at the top of the messages page.
#. count is the integer number of unread messages, 0 or greater.
msgid "{count} unread messages"
msgid_plural "{count} unread messages"
msgstr[0] ""
msgstr[1] ""
```

---

## Core concept 7 — Pseudo-localization

Pseudo-localization is a smoke test that runs before any human translator
sees the strings. It transforms the source English into a pseudolocale that:

- **Expands every string 30–40%** to surface truncation bugs (German /
  Russian / Finnish proxy)
- **Replaces ASCII characters with accented Latin equivalents** to surface
  hardcoded ASCII assumptions
- **Wraps every string with sentinels** like `[!! … !!]` to surface
  un-extracted hardcoded strings
- Optionally **mirrors the string** to surface RTL bugs

Example transformations:

| Original | Expanded accented | Sentinel-wrapped | Mirrored |
|----------|-------------------|------------------|----------|
| `Save` | `[!! Šåvé !!]` | `[!! Save !!]` | `evaS` |
| `Settings` | `[!! Šéttîñgš (extended) !!]` | `[!! Settings !!]` | `sgnitteS` |
| `Welcome, {name}!` | `[!! Wélçömé, {name}! (with extra text to pad length) !!]` | `[!! Welcome, {name}! !!]` | `!{name} ,emocleW` |

Mozilla Fluent has built-in pseudolocale support; React, Vue, and Android
have third-party libraries. **Run pseudo-localization in CI** before every
release to catch hardcoded strings and truncated UI.

---

## Core concept 8 — RTL-friendly writing

Right-to-left languages (Arabic, Hebrew, Persian, Urdu) force three writing
considerations:

1. **Avoid baked-in directional assumptions.** "Click the arrow on the right"
   becomes wrong in Arabic. Prefer "Click the arrow next to the search box."
2. **Use logical CSS properties at the engineering layer.** That is a CSS
   concern, not a writing concern, but the writer must not write copy that
   contradicts the layout once mirrored.
3. **Numbers stay LTR inside RTL text.** This is automatic in Unicode
   bidi, but copy that says "version 2.3.1" should not rely on directional
   formatting outside the bidi algorithm.

---

## Core concept 9 — Avoiding text in images

Text in images cannot be translated without re-rendering the image. Three
rules:

1. **No baked text in screenshots, hero images, marketing banners.** Use
   HTML overlays.
2. **Iconography, not labels-in-pixels.** A magnifying-glass icon is
   universal. The word "Search" inside the icon image is not.
3. **Diagrams: separate the diagram from the labels.** Provide labels as
   text overlays or as `<text>` elements in SVG so the translation pipeline
   can substitute them.

When text in an image is unavoidable (corporate logo, brand wordmark): treat
it as untranslated. Document the exception.

---

## Core concept 10 — Key naming conventions

A good string key is stable, semantic, and namespaced. Patterns:

```text
GOOD (semantic, namespaced):
  inbox.unread.label
  settings.security.two_factor.toggle
  errors.network.timeout.body
  modal.delete_account.title

BAD (content-derived, fragile):
  "Save changes"           # key changes every time copy changes
  "msg1", "label2"          # opaque
  "areYouSureYouWantTo"     # camelCase + content-derived

ALSO GOOD (Fluent flat namespace):
  inbox-unread-label
  settings-security-2fa-toggle
```

Three rules:

1. **The key describes the role, not the content.** `confirm.save` not
   `save_changes`.
2. **Namespace by feature, then by sub-feature.** `inbox.compose.title`.
3. **Don't bury locale in the key.** The key is locale-invariant; the file
   is locale-specific (`messages.en.json`, `messages.de.json`).

---

## Templates

### Translation-friendly rewrite (before / after)

| Source (untranslatable) | Translation-friendly |
|-------------------------|----------------------|
| "Hit the ground running with our new dashboard" | "Get started fast with the new dashboard" |
| "We've got your back" | "We support you" |
| "Move the needle on engagement" | "Increase engagement" |
| "Q4 hit out of the park" | "Q4 exceeded the target" |
| "Welcome aboard!" | "Welcome to {productName}" |
| "Setting things up..." | "Setting up your account..." |
| "We're cooking up something special" | "We are preparing a new feature" |
| "Click here to learn more" | "Open the {pageName} page" |

### ICU plural — full template

```icu
{messageCount, plural,
  =0 {You have no new messages.}
  one {You have one new message.}
  other {You have # new messages.}
}
```

### ICU select — gender-aware

```icu
{authorGender, select,
  female {{authorName} updated her profile.}
  male {{authorName} updated his profile.}
  other {{authorName} updated their profile.}
}
```

### ICU nested plural + select — donation acknowledgment

```icu
{donorCount, plural,
  =0 {No donations yet.}
  one {{donorGender, select,
    female {She donated #.}
    male {He donated #.}
    other {They donated #.}
  }}
  other {{donorCount, number, currency} from # donors.}
}
```

### Translator comment block (universal pattern)

```text
KEY: inbox.unread.header

CONTEXT:
  Type:        UI string — page header
  Surface:     Top of the messages page (between top nav and message list)
  Trigger:     Always visible when inbox is open
  Visibility:  Logged-in users

PLACEHOLDERS:
  {count}     Integer, 0 or greater. The number of unread messages
              for the current user.

NOTES FOR TRANSLATOR:
  - "Inbox" here refers to the in-app messaging inbox, not email.
  - Some locales may prefer "messages" over "items"; use whichever
    is natural.
  - This string can wrap to two lines on narrow viewports.
```

### Pseudo-localization CI check (shape)

```bash
# Pseudo-localize all en-US strings into the pseudolocale "qps-ploc"
npm run i18n:pseudo-localize -- --source=en --target=qps-ploc

# Run UI tests against the pseudolocale to catch:
#   - truncation (text overflow)
#   - hardcoded ASCII (non-extracted strings)
#   - non-extracted strings (no sentinel wrap)
#   - RTL bugs (when paired with --mirror)
npm run test:e2e -- --locale=qps-ploc
```

### Key-naming pattern — i18next + Fluent + gettext

```text
i18next:
  "inbox.unread.header": "{count, plural, ...}"

Fluent:
  inbox-unread-header = {$count ->
      [0] No unread messages
     *[other] {$count} unread messages
  }

gettext:
  msgctxt "inbox.unread.header"
  msgid "{count} unread messages"
  msgstr ""
```

---

## Anti-patterns

| Anti-pattern | Why it fails | Fix |
|--------------|--------------|-----|
| Concatenating UI fragments | Word order is fixed by English | Single string with named placeholders |
| Positional `%s %s` | Translator can't reorder | Named `{username}`, `{date}` |
| Only `one` and `other` plural categories | Russian, Arabic, Polish break | Author all CLDR forms via tooling |
| No translator comment | Translator guesses; gets it wrong | Comment every non-trivial string |
| Idioms in source | No literal translation | Rewrite to the underlying meaning |
| Hard-coded ASCII assumption | Accented chars in pseudolocale break layout | Pseudo-localize in CI |
| Text baked into images | Untranslatable without re-rendering | HTML overlay or SVG `<text>` |
| Content-derived keys (`"Save changes"`) | Key changes every copy edit | Semantic, namespaced keys |
| Sentence built from two strings | Translator can't see the whole | Single message; placeholders inside |
| Selecting on gender without `other` | Non-binary / unset breaks at runtime | Always include `other` branch |
| Abbreviations only ("Q1", "EOY") | Unrecognized in 27 of 30 locales | Spell out at first use |
| "Click here" links | Forces English word order | Rich-text interpolation around the link |
| Right-side / left-side language | Wrong after RTL mirror | Reference by neighbour element |

---

## Decision heuristics

**Does this string need to be localized?**

- Is it user-visible? → Yes, localize
- Is it a developer-facing log? → No, leave in English
- Is it a brand wordmark? → No (untranslated by policy)
- Is it a tooltip? → Yes
- Is it a server-side error string that reaches a user? → Yes

**Plural or select?**

- Variable is numeric and grammar depends on it → `plural` (or
  `selectordinal` for 1st/2nd/3rd)
- Variable is a string with a small fixed set of values → `select`
- Variable affects ordinal grammar → `selectordinal`

**Placeholder discipline:**

- One placeholder per concept → name it
- Position-only placeholders (`%1$s`, `%2$d`) → migrate to named
- Concatenation in source code → rewrite as a full message

**Pseudo-localization smoke test ordering:**

1. Expand to surface truncation
2. Sentinel-wrap to surface un-extracted strings
3. Mirror to surface RTL bugs
4. Run UI tests against the pseudolocale before any human translation

**Key naming:**

- Stable across copy edits? → Semantic, role-based
- Reflects the surface? → Yes (namespace by feature)
- Encodes the locale? → No (locale lives in the file, not the key)

---

## References

- Unicode CLDR: [Plural
  Rules](https://cldr.unicode.org/index/cldr-spec/plural-rules) and the
  full CLDR repository
- ICU: [Formatting
  Messages](https://unicode-org.github.io/icu/userguide/format_parse/messages/)
  — canonical MessageFormat syntax
- Mozilla L10n: [Best practices for
  developers](https://mozilla-l10n.github.io/documentation/localization/dev_best_practices.html)
  and the Mozilla General Localization Style Guide
- W3C i18n Activity: `w3.org/International/` — articles on bidi, language
  declaration, and source-text best practices
- i18next: [Best
  Practices](https://www.i18next.com/principles/best-practices) and
  [Interpolation](https://www.i18next.com/translation-function/interpolation)
- Project Fluent: `projectfluent.org` — Mozilla's modern MessageFormat
  alternative
- MessageFormat 2 (ECMA-402 working draft) — the successor to ICU
  MessageFormat, in late-stage draft as of 2026
- Google Material Design: i18n guidance (RTL mirroring, number formatting)
- Phrase: ["10 Common Software Localization
  Mistakes"](https://phrase.com/blog/posts/10-common-mistakes-in-software-localization/)

---

## Related skills

- `plain-language` — when the same source must also hit a reading-grade
  target
- `writing-expert` — for tone and voice in the source
- `visual-writing` — when localizable text appears in captions or alt text
- `inclusive-language` — when gendered or culturally biased terms appear
- `frontend-design` — for the layout side of RTL and string expansion
