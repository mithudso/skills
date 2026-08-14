<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `error-message-craft` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: error-message-craft
description: Voice and clarity for error messages in user-facing software — the "what / why / what to do next" triple, no-blame language, error-code naming conventions, status-code-vs-error-code terminology, i18n and localization considerations, conservative punctuation, and Nielsen Norman Group hostile-pattern avoidance. Grounded in Microsoft Writing Style Guide, Google Developer Documentation Style Guide, and NN/g error-message research. TRIGGER: "write an error message", "this error message is unclear", "improve this 500/4xx/5xx response", "how should we word this validation error", "make this error message kinder", "i18n our errors", "error code convention", "what does this CLI error say", "rewrite this stack trace banner", "design a form-validation error". SKIP: exception-handling code paths or retry logic (those are coding-patterns or backend-patterns concerns); status-page customer-incident messages (use incident-comms); long-form documentation tone (use technical-writing-craft); plain-language translation of complex content (use plain-language); microcopy unrelated to errors (use frontend-design).
---

# Error Message Craft

## Overview

Error messages are the most-read text most software ever produces. They are also the worst-edited. A login form may have its happy-path microcopy reviewed by three people; the "Invalid credentials" message ships with whatever the engineer typed at 02:00. The result is a class of writing — error text — that consistently fails users at the moment they most need to be helped, and consistently leaks frustration and blame into the product voice.

This skill governs the writing of error text specifically: HTTP error bodies, CLI error output, form-validation errors, in-product error banners, system-status notices, and the prose inside structured error objects (e.g., `{ "error": { "message": "..." } }`). The defining constraints come from three reference style guides that converge on the same triple: tell the user **what** went wrong, **why** it went wrong (when it helps them act), and **what to do next**. The Microsoft Writing Style Guide error chapter, the Google Developer Documentation Style Guide, and Nielsen Norman Group's error-message guidelines all reach the same architecture from different starting points.

Use this skill when:

- Writing or reviewing a new error message for a CLI, API response, form, or in-product banner.
- Standardizing error wording across a product surface that currently has 12 different phrasings of the same problem.
- Designing error-code conventions (numeric ranges, prefixes, status-code vs error-code terminology).
- Localizing error text for non-English audiences (i18n implications on tone, exclamation marks, pluralization).
- Auditing existing error text for blame language, exposed internals, or marketing voice.
- Writing the error-message section of a developer-facing API reference.

Skip this skill when:

- You are writing the exception-handling code itself (retry logic, error classification, log lines for operators). Use `coding-patterns`, `backend-patterns`, or `nodejs-observability`.
- You are writing customer-facing incident messages on a status page. Use `incident-comms`.
- You are writing general microcopy unrelated to errors (button labels, empty states, success toasts). Use `frontend-design`.
- You need plain-language rewriting for legal/accessibility compliance of long-form content. Use `plain-language`.

## Core Concepts

### 1. The "what / why / what to do next" triple

Every effective error message answers three questions in the user's order of need:

1. **What went wrong** — the observable failure, stated in user terms, not internal terms.
2. **Why it went wrong** — only when the cause helps the user decide what to do, and only when you can state it without speculation.
3. **What to do next** — a concrete action, ideally one click or one command.

This triple is the consensus across Microsoft's Error Message Guidelines, NN/g's error-message research, and the Google Developer Documentation Style Guide's form-validation guidance. Microsoft puts it as: state the problem clearly, explain what caused it (when helpful), and provide a solution. NN/g puts it as: be explicit, human-readable, precise, and give constructive advice.

The triple does not have to be three separate sentences. For a short field validation, it can collapse:

> "Email must include an @ symbol."

Here *what* (must include an @) implies *why* (current input lacks it) and *what next* (add an @). For an API 5xx, it often must be three clauses:

> "We couldn't save your changes (database unreachable). Wait 30 seconds and retry; if this persists, contact support with request ID `req_abc123`."

The pattern is: omit a clause only when the user can infer it from the other two. Never omit "what to do next" except when the user genuinely has no agency (in which case, say so explicitly — see Concept 6).

### 2. No-blame language — the system owns the failure

The most pervasive defect in real error text is blame language: "You entered an invalid email" instead of "Email must include an @ symbol." Blame language reads as the system pointing a finger at the user, which is exactly the opposite of what an error message should accomplish.

Three substitutions do most of the no-blame work:

- **"You [did wrong thing]" → "[Field] requires [thing]"** — passive on the user, active on the requirement.
- **"You failed to..." → "[Action] needs..."** — describes what is needed, not what was missed.
- **"Invalid input" → "[Field name] must be [format]"** — specific about the actual constraint, not a generic judgment.

NN/g's "Hostile Patterns in Error Messages" research is explicit on this: error messages should not blame the user, mock them with snark, expose technical details that imply user fault, or use ALL CAPS (read as shouting). The Microsoft style guide's tone guidance — "warm, relaxed, conversational" — exists specifically to crowd out the blame voice that ships by default.

A useful linguistic test: replace "you" with "the system" in the draft. If the sentence now reads as the system admitting a defect ("The system did not accept the email because it lacked an @"), the original was probably user-friendly. If it now reads as nonsense ("The system entered an invalid email"), the original was blame language and needs rewriting.

### 3. State the problem in user terms, not internal terms

The second pervasive defect is exposing implementation internals: stack-trace fragments, internal class names, database column names, third-party library messages. "Error: ECONNREFUSED 127.0.0.1:5432" tells the user nothing about what they did or what they should do.

The rewrite rule: every user-facing error message has a *user-vocabulary* form and a *log-vocabulary* form, and they are not the same string. The user sees:

> "Couldn't connect to the database. Try again in 30 seconds. If this persists, contact support and include the time of this error."

The log sees:

> `ERROR svc=auth method=login userId=12345 cause=ECONNREFUSED host=db-primary:5432 trace=...`

The Microsoft style guide is direct: "Do not use technical jargon. Use terminology that your audience understands." The Google Developer Documentation Style Guide reinforces this for developer-facing docs by giving precise naming for the rare technical terms that are appropriate ("status code" rather than "response code" or "error code"; numbers and names in code font).

For developer-audience error text (CLI tools, SDK errors), internal vocabulary is appropriate when the user *is* a developer and the term is documented. But "internal" still means "internal to the product or domain," not "internal to the codebase." `INVALID_AUTH_TOKEN` is appropriate; `ERR_AuthMiddlewareException_DTO_FAIL` is not.

### 4. Error codes — naming conventions and the status-code-vs-error-code distinction

Error codes are the most-litigated terminology choice in API design. The Google Developer Documentation Style Guide is explicit: "Call it a status code instead of a response code or error code, and put the number and the name in code font." This convention follows HTTP RFC usage and IETF documents that use "status code" as the official term for the 3-digit number.

Within that, structured error codes (the application-level identifier that travels alongside the HTTP status) follow a different naming convention. Real conventions across major APIs (Stripe, AWS, Google Cloud, GitHub):

- **UPPER_SNAKE_CASE** for the identifier (`INVALID_API_KEY`, not `invalid-api-key` or `invalidApiKey`).
- **Stable across versions** — once published, never rename. Add new codes; never repurpose old ones.
- **Categorized by prefix** — `AUTH_*` for authentication, `RATE_*` for rate-limiting, `VALIDATION_*` for input. A code is greppable in customer support tickets.
- **Documented** — every code has a docs page explaining what triggers it and what to do.

The error-code-in-message pattern is:

> "Invalid API key (code: `INVALID_API_KEY`). Check that your key is correctly set in the `Authorization` header. See [docs/errors/INVALID_API_KEY] for details."

The code is the bridge between the human message and the documentation. Never publish an error code without a documentation page.

### 5. Provide a concrete next action, or admit you cannot

The "what to do next" clause is where error messages most often go vague. "Please try again later" is the canonical failure: it gives no time hint, no escalation path, no diagnostic information.

A complete next-action clause names one of:

- **A retry hint with a time bound**: "Try again in 30 seconds." Not "try again later."
- **A specific input fix**: "Add an @ symbol to your email." Not "fix the email."
- **An escalation path with the diagnostic info**: "Contact support and include request ID `req_abc123`." Not "contact support."
- **A documentation link** for codes the user can act on: "See [docs/errors/RATE_LIMIT] for backoff guidance."

If genuinely no user action exists ("The database is being restored from backup"), say so explicitly and provide a status-page link:

> "Service is temporarily unavailable while we recover from an outage. No action is needed on your side. See [status.example.com] for updates."

This honesty pre-empts the inbound support ticket asking "what should I do?"

### 6. Conservative punctuation and tone — especially exclamation marks

The Google Developer Documentation Style Guide is specific: "Use exclamation points when an exclamation point is part of a specific error code or log message that must be matched exactly. However, be conservative with exclamation points in global documentation, as in some languages such as Japanese or Korean, exclamation marks can come across as overly emphatic or even shouting."

The same applies to all punctuation that carries tonal weight in some locales:

- **Exclamation marks**: avoid in error text. The error is already attention-getting.
- **ALL CAPS for emphasis**: never. Reads as shouting in all locales. Use bold or italics in rich-text contexts.
- **Ellipses (…) for trailing-off**: avoid. Reads as passive-aggressive uncertainty.
- **Multiple question marks**: never. "??? Couldn't find it" reads as mockery.
- **Smiley/sad emoji**: never in error text. Reads as performative in failure contexts.

The default tone is calm and direct. The Microsoft style guide phrases this as "warm and relaxed by using a conversational tone." NN/g's video on efficient error messages reinforces: the message should be polite, precise, and give constructive advice — none of which require exclamation points.

### 7. Internationalization — pluralization, formal/informal register, length budget

Error messages get translated. The translation will be longer than the English source (Germanic languages run 30-50% longer; some Asian-language translations run shorter but with very different tonal connotations). The Locize and Lingoport i18n guides emphasize: error messages that work in English can become unreadable, untranslated, or culturally tone-deaf in other locales.

The writing implications:

- **Avoid contractions** where translation may be awkward ("can't" → "cannot" is safer for translators).
- **Avoid idioms** ("hit a snag", "ran into a wall"). They do not translate.
- **Avoid concatenation in code** ("Error: " + fieldName + " is invalid"). Concatenation breaks grammatical agreement in gendered or inflected languages. Use full-sentence templates with placeholders the translator can rearrange.
- **Plan for pluralization complexity**. English has two plural forms (one/many); Arabic has six; Russian has three. A naive `1 file` / `n files` pattern breaks in target locales. Use ICU MessageFormat or equivalent.
- **Avoid culturally-bound politeness markers**. "Please" is a default English politeness marker that can read as obsequious in translation; "we're sorry to bother you, but..." reads as performative in many locales. Plain direct phrasing translates more cleanly.
- **Budget for length**. A 40-character English error may need to fit in a 60-character German translation in the same UI slot. Either design the slot to expand or write shorter English.

The Microsoft style guide's tone guidance and the i18n literature converge: plain, direct, jargon-free English translates well. Clever or idiomatic English does not.

### 8. NN/g hostile patterns — what never to ship

NN/g's "Hostile Patterns in Error Messages" article is a catalog of failure modes that consistently damage user trust. The patterns to systematically avoid:

- **Mockery**: "Oops! Something went terribly wrong! 😱" — reads as performative panic.
- **Self-deprecation**: "Our bad! We messed up." — undermines confidence; the user wants competence, not a confession.
- **Blame deflection**: "Your browser is too old." — even if true, frame it as a system requirement, not a user judgment.
- **Vague hedging**: "An unexpected error occurred." — uninformative.
- **Exposed internals**: stack traces, internal class names, raw DB errors.
- **Marketing voice in failure**: "Thanks for your patience as we work to deliver an amazing experience!" — cardinal sin.
- **Forced positivity**: "Great news! There was an error." — surreal and untrusted.

The NN/g error-message scoring rubric grades real error text against these patterns. The single most effective edit pass on existing error text is: read every message aloud, ask "would I say this to a customer in person?", and rewrite anything that fails the test.

### 9. Form-validation errors — inline, contextual, specific

Form-validation errors have their own conventions because they appear next to the field that failed:

- **Inline, near the field**, not at the top of the form. NN/g: "show an actionable error message below or next to the problem field."
- **Specific to the field's actual requirement**, not generic. "Password must be at least 12 characters" not "Invalid password."
- **Suggest the fix when possible**. NN/g's "City and ZIP code don't match" example: offer a button for the city that matches the ZIP the user entered.
- **Don't validate prematurely**. Showing "email is invalid" while the user is still typing the @ symbol is hostile. Validate on blur or submit.
- **Color and icon are not enough**. The text must carry the meaning; color is a redundancy, not a substitute (accessibility requirement).
- **Reading level**. NN/g: write error messages at a 7-8th grade reading level (Flesch-Kincaid) or lower.

For CLI tools, the equivalent of inline form-validation is the error pointing to the exact argument that failed:

```
Error: --port requires a number between 1 and 65535. Got: "abc"
```

The error names the flag, names the constraint, names what was received.

### 10. Logging vs user-facing — separate strings, separate audiences

This is the meta-rule that subsumes several others: the string a user sees is not the string the log captures. They serve different audiences with different needs.

- **User string**: short, action-oriented, no internals, includes correlation ID for support.
- **Log string**: full context, structured fields, internal IDs, stack traces, correlation ID matching the user string.

The correlation ID (request ID, trace ID, error ID) is the bridge. The user message includes it ("...request ID `req_abc123`"); the log includes it (`reqId=req_abc123`). Support can correlate the two.

The implementation implication: error objects in code should carry both strings (or a code that resolves to both). The Microsoft style guide is explicit: "Write a separate error message for each known cause of the error. Do not use a single, generic message to explain every possible reason for the error." That principle applies twice — once for users, once for logs.

## Templates and Patterns

### Central artifact: a complete error message

```
Couldn't save your changes.

The database is temporarily unavailable. Wait 30 seconds and retry.

If this persists, contact support and include the request ID below.
Request ID: req_abc123
Code: DB_UNAVAILABLE
```

Four parts: **what** (couldn't save), **why** (DB unavailable, stated in user terms), **what next** (wait + retry, with escalation path), **support handoff** (correlation ID + code).

### API error response body

```json
{
  "error": {
    "code": "INVALID_API_KEY",
    "message": "The API key you provided is not valid. Check that your key is set correctly in the Authorization header.",
    "status": 401,
    "request_id": "req_abc123",
    "docs": "https://docs.example.com/errors/INVALID_API_KEY"
  }
}
```

Structured fields: `code` (UPPER_SNAKE_CASE, stable), `message` (human-readable, no internals), `status` (HTTP status code, in the body for clients that lose the header), `request_id` (correlation), `docs` (link).

### CLI error output

```
Error: --port requires a number between 1 and 65535.

  Received: "abc"
  Expected: integer in range 1..65535

Try: --port 8080
See: my-cli --help port
```

Names the flag, the constraint, the actual input, the suggested fix, the help anchor.

### Form-validation error (inline)

```
Email
[ alice@                              ]
⚠ Email must include an @ symbol and a domain (for example, alice@example.com).
```

Inline below the field, names the constraint, gives an example.

### Rewriting a blame-language error

| Before (blamey, internal-leaking, vague) | After (user-voiced, specific, actionable) |
|---|---|
| "You entered an invalid email." | "Email must include an @ symbol (for example, alice@example.com)." |
| "Error: ECONNREFUSED" | "Couldn't connect to the database. Try again in 30 seconds. If this persists, contact support with request ID `req_abc123`." |
| "An unexpected error occurred." | "Something went wrong on our side (code: `INTERNAL_500`). The error has been logged with ID `req_abc123`. Please retry; if the issue persists, contact support." |
| "Invalid input." | "First name must be 1-50 characters and cannot contain digits." |
| "Oops! 😱 Our bad!" | "Couldn't save your changes. We've logged the issue (id: `req_abc123`). Try again, or contact support if this persists." |
| "Please try again later." | "Try again in 30 seconds. If this persists for more than 5 minutes, see [status.example.com]." |

### Error code documentation page (template)

```markdown
# INVALID_API_KEY

**Status code**: 401

**What it means.** The API key sent in the `Authorization` header is not
recognized by our authentication service. Either the key was mistyped,
revoked, or belongs to a different environment.

**Common causes.**
- The key was rotated and the new key has not been deployed.
- A test-environment key is being sent to the production endpoint.
- The `Authorization` header is missing the `Bearer` prefix.

**How to fix.**
1. Verify the key in your dashboard at [...]
2. Confirm the environment (test/live) matches the endpoint.
3. Ensure the header is `Authorization: Bearer <key>`.

**Related**: `EXPIRED_API_KEY`, `REVOKED_API_KEY`.
```

## Anti-Patterns

- **Marketing voice in failure**: "Thanks for your patience as we work to deliver an amazing experience!" — never.
- **Exposed internals**: stack traces, raw DB errors, internal class names. Bridge with a correlation ID; do not leak the trace.
- **Generic catch-all**: a single "Something went wrong" message for every cause. Write a separate message per known cause.
- **Premature validation**: showing "email is invalid" while the user is typing the @ symbol. Validate on blur or submit.
- **Punctuation theater**: exclamation marks, ALL CAPS, ellipses for trailing-off, multiple question marks, panicky emoji.
- **Vague time hints**: "try again later." Always give a number.
- **Untranslatable idioms**: "hit a snag", "ran into a wall." Plain phrasing translates.
- **Concatenated strings in code**: `"Error: " + field + " is invalid"`. Breaks i18n. Use full-template strings.
- **Color-only signaling**: red text without the word "Error" or an icon with a text label. Accessibility failure.
- **Documented-but-unfindable codes**: shipping `INVALID_API_KEY` with no docs page.

## Decision Heuristics

- **When to include the cause clause ("why")**: include when it changes what the user should do. "Database unreachable" → user should retry. "Subscription expired" → user should renew. If the cause doesn't change the action, drop it.
- **When to include a correlation ID**: always for 5xx errors; usually for unexpected 4xx errors; rarely for expected validation errors.
- **When to suggest a fix**: always, except when no user action is possible (then state that explicitly and link to status).
- **When to link to docs**: always when the error has a documented code; never as a substitute for an inline explanation of what to do right now.
- **When to call this skill from another**: any time `backend-patterns`, `api-design-patterns`, `frontend-design`, or `coding-patterns` produces user-facing error text, route the text through `error-message-craft` for a writing pass. Any time `writing-expert` is asked about a CLI/API/form error, route here.
- **When to escalate to `plain-language`**: when the audience includes non-developers, regulatory contexts (banking, healthcare), or accessibility-mandated reading levels below 7th grade.
- **When to escalate to `incident-comms`**: when the "error" is a multi-customer outage rather than a single-user failure. Status-page voice is a different genre.
- **When to break the rules**: developer-facing CLI tools where the user is a developer and the convention is dense, terse output. The triple (what/why/next) still applies; the wording can compress.

## References

- [Microsoft Writing Style Guide — Error Message Guidelines (Win32)](https://learn.microsoft.com/en-us/windows/win32/debug/error-message-guidelines) — the canonical Microsoft guidance: state the problem, explain the cause when helpful, provide a solution, avoid jargon.
- [Microsoft — Writing style (Windows apps)](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) — the broader writing-style guidance including the conversational-tone rule.
- [Microsoft — Formatting text in instructions](https://learn.microsoft.com/en-us/style-guide/procedures-instructions/formatting-text-in-instructions) — text formatting conventions that apply to error messages.
- [Microsoft — UX guidelines for errors (Business Central)](https://learn.microsoft.com/en-us/dynamics365/business-central/dev-itpro/developer/devenv-error-handling-guidelines) — developer-facing error-handling guidance.
- [Microsoft — Windows Admin Center UI text and design style guide](https://learn.microsoft.com/en-us/windows-server/manage/windows-admin-center/extend/guides/ui-text-style-guide) — concrete examples of error and notification copy.
- [Google Developer Documentation Style Guide](https://developers.google.com/style) — the home of the developer-doc style guide referenced for error code naming and punctuation.
- [Google Developer Style Guide — Notices](https://developers.google.com/style/notices) — formatting for cautions, warnings, errors.
- [Google Developer Style Guide — Periods and other end punctuation](https://developers.google.com/style/periods) — exclamation-mark guidance and i18n implications.
- [Google Developer Style Guide — Word list](https://developers.google.com/style/word-list) — "status code" vs "response code" vs "error code" terminology.
- [Google Developer Style Guide — Write accessible documentation](https://developers.google.com/style/accessibility) — accessibility implications for error text.
- [NN/g — Error-Message Guidelines](https://www.nngroup.com/articles/error-message-guidelines/) — explicit, human-readable, polite, precise, constructive advice.
- [NN/g — 10 Design Guidelines for Reporting Errors in Forms](https://www.nngroup.com/articles/errors-forms-design-guidelines/) — inline, contextual, specific form-validation errors.
- [NN/g — Hostile Patterns in Error Messages](https://www.nngroup.com/articles/hostile-error-messages/) — catalog of error-text failure modes to avoid.
- [NN/g — An Error Messages Scoring Rubric](https://www.nngroup.com/articles/error-messages-scoring-rubric/) — quantitative rubric for auditing existing error text.
- [NN/g — Usability Heuristic 9: Help Users Recognize, Diagnose and Recover from Errors](https://www.nngroup.com/videos/usability-heuristic-recognize-errors/) — Nielsen's heuristic that underpins the field.
- [Locize — What is i18n? (2026 edition)](https://www.locize.com/blog/what-is-i18n/) — i18n fundamentals including error-message translation.
- [Lingoport — Fix i18n Bugs: Best Practices for Software Localization](https://lingoport.com/blog/i18n-bugs/) — concrete i18n bugs to avoid in error text.
- [Localizely — Guidelines for error handling and localization of error messages](https://localizely.com/blog/error-messages/) — practical i18n guidance for error strings.
