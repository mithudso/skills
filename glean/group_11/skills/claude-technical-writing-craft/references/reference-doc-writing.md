<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `reference-doc-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: reference-doc-writing
description: Information-oriented documentation craft — the Diátaxis reference quadrant. Write technical descriptions of an API, configuration, schema, CLI, or protocol where the reader's job is to look something up, not learn or follow. Covers function signatures, parameter tables, return values, exhaustive coverage over narrative, strict consistency, alphabetical/logical ordering, search-discoverability, structure-mirrors-the-thing-described, the no-surprises rule, and the principle that examples illustrate the API rather than teach it. TRIGGER when the user asks to write reference docs, an API reference page, parameter/option table, schema documentation, CLI reference, function/method/class reference, configuration reference, or any doc whose primary purpose is rapid factual lookup by a working user. Also TRIGGER for OpenAPI/JSON-Schema descriptions, error code tables, and CLI flag enumeration. SKIP: hand-held lessons for newcomers (use tutorial-writing); goal-directed recipes (use howto-writing); concept/why/background (use explanation-doc-writing); REST/SDK endpoint authoring with auth+versioning context (use api-docs-craft — which references this skill for the reference page itself); sentence-level prose review (use technical-writing-craft).
---

# Reference Doc Writing — Diátaxis Reference Quadrant

## Overview

Reference documentation is **a description**. It tells the reader what something **is**, what its parts **are**, and what each part **does**. It does not teach, it does not advocate, and it does not narrate. The reader is at work, holds a specific lookup question (*"what does this flag default to?"*, *"what fields does this response return?"*, *"what errors can this throw?"*), and needs the answer in seconds.

Reference docs are information-oriented and sit in the theoretical + working quadrant of Diátaxis. The voice is neutral, the structure is regular, and the coverage is exhaustive. If a tutorial is a lesson and a how-to is a recipe, reference is **a map** — flat, accurate, complete.

The single most-violated rule: **reference does not interpret**. It states facts. Discussion belongs in explanation docs; storytelling belongs in tutorials; opinions belong nowhere in reference.

## Core Concepts

### 1. Architecture mirrors the thing described

If the product has modules, the reference has modules. If a class has methods, the doc has a section per method. If the CLI has subcommands, the doc tree has a page per subcommand. The doc topology is **isomorphic** to the API topology. Readers navigate by analogy to the code they already know.

### 2. Exhaustive coverage beats narrative

Every parameter is listed. Every return value is documented. Every error code is enumerated. Every configuration key appears. There is no "for brevity we omit…" — reference is the *only* place where the reader can be sure they have not missed an option. Omissions are bugs.

### 3. Strict consistency

Every function reference has the same sections in the same order. Every parameter table has the same columns. Every error code is formatted the same way. Consistency is not aesthetic — it lets a reader who has read one page **skim** the next page at 10× speed because the layout is predictable.

### 4. Neutrality

Reference does not say *"you'll usually want…"*, *"the recommended approach is…"*, or *"avoid this except in edge cases"*. Those judgments belong in how-tos and explanation docs. Reference states: *"`timeout`: integer, milliseconds, default `30000`, minimum `0`, maximum `600000`"*. The reader decides what to do with it.

### 5. Examples illustrate, never teach

Reference examples are **specimens**, not lessons. A canonical-call example next to a function reference shows the shape — argument positions, return shape, a representative success and a representative failure. It does **not** walk the reader through what the function is for or why they would use it (that's tutorial/explanation territory).

### 6. The no-surprises rule

Anything that could surprise a working reader must be stated explicitly. Default values. Units (ms vs seconds, bytes vs MiB). Whether a field is nullable. Whether an array can be empty. Whether the method mutates input. Whether order matters. Surprises are the failure mode of reference docs — every silent assumption is a future bug report.

### 7. Search-discoverability

Reference is read by **search**, not by table-of-contents traversal. The page title and the first sentence must contain the term the reader will type. Headings must be the names of the things they describe (`POST /v1/users`, `createSession(opts)`, `--max-retries`). Avoid clever titles. The page for `setTimeout` is named `setTimeout`.

### 8. Stable ordering

Pick an ordering rule per axis and apply it everywhere:

- **Alphabetical** for catalogs (CLI flags, config keys, error codes, environment variables).
- **Logical / call-order** for function references (constructor, then lifecycle methods, then utility methods).
- **Signature-order** for parameter tables (positional first in their declared order; keyword/optional after).

Whatever the rule, document it once and never violate it. Inconsistent ordering kills skim-speed.

### 9. Versioning is part of the contract

Each documented entity should say which version introduced it, which version deprecated it, and which version removed it. This is not optional in mature systems — readers debug against old SDKs constantly and need to know *"is this method in 4.2?"*.

### 10. The page is the unit, not the section

A reference page is a single navigable URL per entity (per endpoint, per function, per command). Cramming the entire surface area onto one mega-page defeats search and inflates load time. The exception is small flat catalogs (a 20-row error-code table) that benefit from being on one page for grep-ability.

## Templates

### Template — function/method reference

```markdown
# `<functionName>(<params>) → <return>`

<One-sentence factual statement of what the function does.>

**Since:** v1.4
**Stability:** stable
**Module:** `package.module`

## Parameters

| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| `opts.timeout` | `number` (ms) | no | `30000` | Maximum time before the call rejects with `TimeoutError`. Must be ≥ 0. |
| `opts.signal` | `AbortSignal` | no | — | If aborted, the call rejects with `AbortError`. |

## Returns

`Promise<Result>` — resolves with `Result` on success. Rejects with one of the errors below.

## Errors

| Code | Class | Thrown when |
|---|---|---|
| `TIMEOUT` | `TimeoutError` | `opts.timeout` elapsed before the call resolved. |
| `ABORTED` | `AbortError` | `opts.signal` was aborted. |
| `INVALID_ARG` | `TypeError` | A parameter failed type validation. |

## Example

```js
const r = await client.fetchThing(id, { timeout: 5000 });
```

## See also

- `<relatedFunction>` — for the streaming variant.
- [How to handle timeouts](…) — recipe-style guide.
```

### Template — CLI flag reference

```markdown
# `mytool deploy`

Deploy the current project to the configured target.

## Synopsis

```
mytool deploy [--target <name>] [--dry-run] [--force]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--target <name>` | string | `production` | Named target from `config.yaml`. |
| `--dry-run` | bool | `false` | Print planned actions; make no changes. |
| `--force` | bool | `false` | Skip the pre-deploy lock check. |

## Exit codes

| Code | Meaning |
|---|---|
| 0 | Success. |
| 2 | Configuration error. |
| 3 | Target unreachable. |

## See also

- `mytool rollback`
```

### Template — config key reference (alphabetical)

```markdown
# Configuration reference

All keys are read from `config.yaml`. Keys are case-sensitive.

## `database.host`

- **Type:** string
- **Default:** `localhost`
- **Since:** v1.0

Hostname of the primary database server.

## `database.port`

- **Type:** integer
- **Default:** `5432`
- **Range:** 1–65535

TCP port of the primary database server.

## `database.ssl.mode`

- **Type:** enum: `disable` | `require` | `verify-ca` | `verify-full`
- **Default:** `require`
- **Since:** v1.2

…
```

## Anti-Patterns

### AP-1 — Examples that teach instead of describe

> "Imagine you have a user signup form. You'd call `createUser` like this…"

That's a tutorial example. Reference examples show **the shape of the call**, not the use case. Strip the narrative.

### AP-2 — "Recommended" / "preferred" / "you should"

These are not reference statements. They are advice. Move them to a how-to ("How to choose between A and B") or an explanation ("Why we recommend A").

### AP-3 — Missing defaults, units, or nullability

Every parameter without a documented default, unit, or nullability is a footgun. The reader has to guess or read source. Reference exists precisely to eliminate that read.

### AP-4 — Incomplete error tables

> "Throws on failure."

What failures? Under what conditions? With what code? An unenumerated error list is a reference that has decided not to be reference.

### AP-5 — Mixed ordering

Parameters listed alphabetically on one page, by signature on another, by importance on a third. The reader's skim speed collapses. Pick one rule, apply it everywhere.

### AP-6 — Hidden behavior in prose

Critical facts buried in flowing prose ("Note that under conditions where the buffer is full, callers should be aware that…") instead of in the parameter table. Move every fact into a table cell or a bullet. Prose is for explanation docs.

### AP-7 — No since/deprecation metadata

Mature APIs accumulate features over versions. A reference without `Since: vX` forces readers to bisect changelogs. Add the version on every entity from day one.

## Decision Heuristics — is this actually reference?

1. **Will the reader arrive via search, looking for a specific name?** If yes → reference. If they'll browse a TOC sequentially → probably tutorial or explanation.
2. **Is the content exhaustive coverage of a surface area?** If you can plausibly omit items "for brevity", you are not writing reference.
3. **Is the voice neutral and factual?** If you find *"you'll want…"* or *"we recommend…"* — that's a how-to or explanation hiding inside.
4. **Does the structure mirror the code/system structure?** If you've reorganized things by "ease of learning" or "common workflow", you are blending in tutorial/how-to logic.
5. **Could two reference pages be written by two different authors and look identical in structure?** If not, the consistency contract is broken.

When the doc is not reference, switch quadrants:

- Newcomer onboarding → `tutorial-writing`
- Specific goal recipe → `howto-writing`
- Background / why / discussion → `explanation-doc-writing`
- REST/SDK endpoint authoring with auth, versioning, pagination, OpenAPI → `api-docs-craft` (uses this skill for the reference-page form itself)
- Sentence-level prose review → `technical-writing-craft`

## Cross-pollination notes

- **api-docs-craft** is the specialist for REST/GraphQL/SDK reference and links to this skill for the reference-page form.
- **mongodb-error-codes**, **case-mcp-server-guide**, and similar reference catalogs in this hub follow this skill's structural conventions implicitly — consider linking explicitly.
- **technical-writing-craft** holds prose review and Diátaxis theory.
- **doc-archaeology** is the right partner when reconstructing reference docs for an undocumented or legacy system.

## References

1. [Reference — Diátaxis](https://diataxis.fr/reference/) — Procida's canonical specification: reference is neutral, exhaustive, and its architecture mirrors the architecture of the thing described.
2. [Diátaxis — Start here](https://diataxis.fr/start-here/) — the theoretical+working quadrant placement and the working-user model.
3. [Reference guides — Divio Documentation](https://docs.divio.com/documentation-system/reference/) — Divio's "describe, and only describe" formulation, plus the structure-mirrors-product principle.
4. [Documentation Quadrants — The Grand Unified Theory of Documentation (Dunn)](https://dunnhq.com/posts/2023/documentation-quadrants/) — practitioner notes on reference exhaustiveness vs narrative.
5. [Diátaxis — A systematic approach to technical documentation authoring (BSSw.io)](https://bssw.io/items/diataxis-a-systematic-approach-to-technical-documentation-authoring) — research-software perspective on reference completeness and ordering.
