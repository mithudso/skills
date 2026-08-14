<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `api-docs-craft` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: api-docs-craft
description: Author REST/HTTP API documentation that developers can actually integrate against — Diátaxis framework (tutorial / how-to / reference / explanation by Daniele Procida), endpoint reference patterns (Stripe and Twilio gold-standard), error-code documentation with troubleshooting, copy-pasteable code samples in multiple languages, the "deprecation banner" pattern (RFC 8594 `Sunset`/`Deprecation` headers), OpenAPI/Swagger UX (three-pane layout, "Try It" widgets, language switchers), API versioning strategies (URI vs header vs date-pinning), and the reference-vs-tutorial separation. TRIGGER when user asks "write API docs", "document this endpoint", "OpenAPI spec for this", "Swagger UI structure", "Diátaxis", "tutorial vs how-to", "what should my error docs include", "deprecation banner", "Sunset header", "code samples in docs", "versioning strategy for docs", "Stripe-style docs", or pastes an endpoint definition asking for the docs page. SKIP when designing the API itself (use api-design-patterns or backend-patterns); writing prose at the sentence level (use technical-writing-craft or writing-expert); authoring developer-marketing or landing-page copy (use writing-expert with sales-and-marketing-copy); generating release notes that aggregate API changes (use changelog-and-release-notes or changelogs-for-humans); authoring an architectural RFC for an API (use rfc-and-design-docs); crafting in-product error-message strings (use error-message-craft); authoring runbooks for the team operating the API (use runbook-craft).
---

# API Documentation Craft

## Overview

**What this skill covers.** Produce REST/HTTP API documentation that meets the gold-standard set by Stripe and Twilio: a Diátaxis-organized structure (tutorials for learners, how-to guides for task-doers, reference for lookups, explanation for the curious), endpoint reference pages that include request/response examples in multiple languages, an error-code catalog with troubleshooting steps, deprecation signaling via RFC-8594 headers and visible banners, and an OpenAPI/Swagger interactive layer that lets developers try the API without leaving the docs.

**When to use.**
- The user is writing docs for a new endpoint, error code, SDK, or webhook.
- The user pastes an OpenAPI/Swagger spec and asks for the doc-page treatment.
- The user is restructuring an existing docs site (typical trigger: "our docs are a mess").
- The user is announcing a deprecation or sunset for an endpoint or version.
- The user is choosing a versioning strategy (URI path vs header vs date-pinned).
- The user is designing the layout/IA of a developer portal.

**When to skip.**
- Designing the API surface itself — use `api-design-patterns` or `backend-patterns`.
- Sentence-level editing of prose in the docs — use `technical-writing-craft` or `writing-expert`.
- Marketing copy for the developer portal homepage — use `writing-expert` with `sales-and-marketing-copy`.
- Generating customer-facing release notes that aggregate many API changes — use `changelog-and-release-notes`.
- Architectural RFCs for an API that doesn't yet exist — use `rfc-and-design-docs`.
- Sentence-level technical writing not specific to APIs — use `technical-writing-craft`.

## Core Concepts

### 1. Diátaxis — the four-quadrant framework

Daniele Procida's Diátaxis framework (diataxis.fr) organizes documentation by *what the reader is trying to do*, not by what the author wants to say. Four quadrants on two axes:

|                       | **Practical steps** (action) | **Theoretical knowledge** (cognition) |
| --------------------- | ----------------------------- | -------------------------------------- |
| **When learning**     | **Tutorials**                 | **Explanation**                        |
| **When working**      | **How-to guides**             | **Reference**                          |

**Tutorial (learning-oriented):**
- For a beginner who has never used the system.
- Guarantees a successful outcome by holding the user's hand the whole way.
- Concrete, specific, opinionated. "Build your first webhook handler in 10 minutes."
- Promises: "If you follow these steps, you will end up with a working X."

**How-to guide (task-oriented):**
- For a competent user who already knows the system but has a specific task.
- Goal-oriented, often a sequence of decisions ("if your auth uses OAuth, do A; if API key, do B").
- "How to rotate an API key without downtime."

**Reference (information-oriented):**
- For a developer who already knows what they want and is looking up the exact shape.
- Comprehensive, accurate, terse. No narrative.
- "POST /v1/charges — request fields, response fields, error codes."

**Explanation (understanding-oriented):**
- For a curious developer who wants the *why* behind a design.
- Discursive, opinionated, free to digress. "Why our IDs are prefixed."
- This quadrant is the most-skipped. Without it, the developer can't predict behavior when the docs don't cover their exact case.

**The fatal mistake:** mixing types on one page. A reference page that drifts into tutorial steps wastes the lookup-er's time and confuses the learner.

### 2. Endpoint reference page — the canonical structure

A reference page for an endpoint should contain, in order:

1. **HTTP method + path** as the page title (e.g., `POST /v1/charges`).
2. **One-paragraph summary** of what the endpoint does.
3. **Request:** required vs optional parameters, each with type, constraints, and one-line description.
4. **Request body schema** (for `POST`/`PUT`/`PATCH`) — JSON shape with field-by-field annotations.
5. **Response schema** — success shape, including all fields. Use the same field-annotation pattern.
6. **Status codes:** every status this endpoint can return, with the trigger condition.
7. **Errors:** structured error codes, each linked to the error catalog (see Concept 4).
8. **Code samples:** the same call in 4-8 languages, language-switchable.
9. **"Try it" widget** (if the docs platform supports it).
10. **Idempotency, rate limits, scopes/permissions** — anything else the developer needs to *successfully* call this endpoint.
11. **Related endpoints** — links to the predecessor/successor calls in typical workflows.

### 3. Code samples — discipline, not decoration

Code samples are *the* thing developers copy. Make them paste-and-run.

**Rules:**
- **Multiple languages.** Stripe's docs ship samples in curl, Node, Python, Ruby, PHP, Go, Java, .NET. At minimum, ship curl + one strongly-typed language + one dynamic language.
- **Language switcher synced across the page.** When the reader clicks "Python" in one sample, all samples on the page switch to Python. Stripe and Twilio both do this.
- **Use environment variables, not literal secrets.** `$STRIPE_API_KEY`, not `sk_test_abc123`. Show the env-var name and link to the "Where do I get this key?" page.
- **Show the response.** The sample isn't complete without an example response inline.
- **Match the API version pinned in the docs.** If the reader is browsing v2024-10-01 docs, samples must be valid against that version.
- **Runnable in isolation.** A sample should not assume context from an earlier sample on the page. Each sample includes its imports/requires.

**Anti-sample:**
```js
// missing import; mystery variable; no response shown
const result = stripe.charges.create({ amount, currency });
```

**Good sample:**
```js
// Node.js
const Stripe = require('stripe');
const stripe = new Stripe(process.env.STRIPE_API_KEY);

const charge = await stripe.charges.create({
  amount: 2000,        // amount in cents
  currency: 'usd',
  source: 'tok_visa',  // test-mode token; see /testing
  description: 'Charge for jenny.rosen@example.com',
});

// charge.id => "ch_3O..."
// charge.status => "succeeded"
```

### 4. Error documentation — the error catalog

Every error code the API can return needs:
- The HTTP status (`400`, `402`, `429`, `500`...).
- The machine-readable error code (`card_declined`, `rate_limited`, `invalid_request_error`).
- The human-readable message (and whether it's safe to surface to end users).
- The trigger condition — what causes this error.
- The remediation — what the developer should do.

Structure these as a separate **Errors** reference section (not inline in every endpoint). Each endpoint references the error codes it can return; the catalog is the source of truth.

**Error catalog entry example:**
```
## card_declined  (402 Payment Required)

The card was declined by the issuer. The response includes a
`decline_code` (e.g., `insufficient_funds`, `lost_card`) that
narrows the reason.

**Common decline_codes:**
- `insufficient_funds` — the card has insufficient available balance.
- `lost_card` — the card was reported lost. Do NOT surface the reason
  to the cardholder; show a generic "card declined" message.
- `do_not_honor` — generic issuer decline. Retry with a different card.

**Remediation:**
- For `insufficient_funds`, prompt the customer for a different card.
- For `lost_card` / `stolen_card`, show generic decline messaging only.
- For transient codes (`issuer_not_available`), exponential-backoff retry.

**Related:** [Decline codes reference](#) · [Testing declines](#)
```

### 5. Deprecation banner pattern — RFC 8594 `Sunset` and `Deprecation`

When an endpoint is deprecated, two layers signal it:

**Layer 1 — HTTP headers (RFC 8594 + the Deprecation header draft):**
```
HTTP/1.1 200 OK
Deprecation: Sun, 11 Nov 2026 23:59:59 GMT
Sunset: Sat, 11 May 2027 23:59:59 GMT
Link: <https://api.example.com/docs/migration/v1-to-v2>; rel="sunset"
```
- `Deprecation` — the date the endpoint *was* deprecated (or a boolean `true`).
- `Sunset` — the date the endpoint will stop responding.
- `Link rel="sunset"` — points to the migration guide.
- After Sunset passes, the server typically returns `410 Gone`.

**Layer 2 — visible banner in the docs:**
Every deprecated endpoint's docs page leads with a styled callout box:
```markdown
> ⚠️ **Deprecated 2026-11-11. Sunsets 2027-05-11.**
>
> This endpoint will return `410 Gone` after the sunset date.
> Migrate to [`POST /v2/charges`](../v2/charges) — see the
> [migration guide](../migration/v1-to-v2).
```
The banner sits *above* the endpoint title. It's the first thing a reader sees, regardless of how they got to the page (search, link, bookmark).

**Two stages of deprecation:**
1. **Deprecated but operational** — endpoint still works; `Deprecation` header set; docs banner shown. No `Sunset` header yet (or `Sunset` is far in the future).
2. **Sunset announced** — `Sunset` header set with a concrete date; docs banner updated with the date; migration guide must be complete.

### 6. API versioning strategies

Three dominant strategies; each has doc implications:

**URI path versioning (`/v1/`, `/v2/`):**
- Stripe's older API, most public REST APIs.
- Pros: visible in URL, easy to route, easy to document (one section per version).
- Cons: encourages clients to hard-code paths; whole-version migration is heavyweight.
- **Doc pattern:** separate doc trees per version; sticky version selector at the top.

**Header versioning (`Accept: application/vnd.example.v2+json` or `X-API-Version: 2`):**
- GitHub's older API.
- Pros: URLs stay clean; clients can pin to a header value.
- Cons: harder to test from a browser; less discoverable.
- **Doc pattern:** version selector picks the header value; samples show the header in every request.

**Date-pinned versioning (`Stripe-Version: 2024-10-01`):**
- Stripe's current API.
- Pros: granular (every breaking change gets a date); accounts pin to a date and stay there.
- Cons: harder to mentally model; changelog density is high.
- **Doc pattern:** every page shows "you are viewing the docs for `2024-10-01`"; an "API versions" changelog page lists every dated change.

**The breaking-change contract.** Whatever versioning strategy: pick a deprecation window (Stripe announces ~12 months ahead) and stick to it. Document it on a single "Versioning" page that every endpoint links to.

### 7. OpenAPI / Swagger UX — three-pane layout and "Try It"

Stripe pioneered the three-pane layout that the industry copied:

| Left pane          | Center pane             | Right pane            |
| ------------------ | ----------------------- | --------------------- |
| Navigation         | Prose + parameters     | Live code samples     |
| (table of contents) | (the reference content) | + interactive request |

**Why three panes:**
- The reader doesn't lose their place in the table of contents.
- The samples are always visible while the reader scans the prose.
- The "Try It" widget can stay pinned in the right pane while the reader scrolls the center.

**"Try It" widgets:**
- Pre-fill with the reader's test-mode API key (after a one-click "use my keys" flow).
- Send the request from the doc page; show the live response below.
- Surface errors clearly with links to the error-catalog entry.

**OpenAPI-driven docs:**
- The OpenAPI/Swagger YAML/JSON spec is the source of truth.
- Doc-rendering tools (Redoc, Mintlify, Bump.sh, Fern, Scalar) generate the reference layer from the spec.
- Tutorials, how-to guides, and explanations are still hand-written prose — OpenAPI generates the *reference* quadrant only.

**Common pitfalls in OpenAPI docs:**
- Missing `description` fields on parameters → renders as a wall of types with no explanation.
- No `example` values → "Try It" widgets are unusable until the reader fills every field by hand.
- Stale spec → samples don't match real responses. Spec must be generated from server code (or vice versa, contract-first) so it can't drift.

### 8. The "explanation" quadrant — where most docs starve

Most docs sites do reference and (sometimes) tutorials well. They skip explanation entirely, and developers pay for it.

Examples of explanation pages that earn their keep:
- "Why our IDs are prefixed" — explains `ch_`, `cus_`, `pi_` so debuggers don't have to guess.
- "How idempotency keys work" — explains the model so developers stop sending random keys per retry.
- "Date-based versioning explained" — explains the pinning model so developers understand why their old code keeps working.
- "How webhooks retry" — explains the back-off curve so developers don't reinvent it.

Explanation pages are essays. They can be long, opinionated, and discursive. They link out to the reference for lookup. They never become how-to guides.

### 9. The "Getting started" path — tutorial-first

Every API doc site needs a single, opinionated "Getting started" tutorial that takes a brand-new reader from `npm install` to a successful first API call in ≤ 10 minutes.

**The 10-minute tutorial:**
1. Install the SDK (one command).
2. Set an environment variable with the test API key (one command + a link to "where do I get this?").
3. Make a first call (one copy-paste).
4. See the response. Confirm success.

This is *the* highest-leverage doc page on the entire site. It's the first impression. Maintain it like a product feature: re-run it monthly to make sure every step still works.

### 10. Search and information architecture

The docs site's search is the entry point for the majority of returning developers. Treat it as a product.

- Index every endpoint, every error code, every parameter name.
- Boost recent / current-version pages over older versions.
- Surface the page *type* (tutorial / how-to / reference / explanation) in the search result. The same word ("webhook") appears in all four quadrants; readers want to know which kind of page they're picking.
- Track the no-result and zero-click queries — those are the gaps in the docs.

## Templates and Examples

### Template — endpoint reference page

```markdown
# POST /v1/charges

Create a charge against a card, customer, or source.

## Request

`POST https://api.example.com/v1/charges`

### Headers
| Header           | Required | Description                              |
| ---------------- | -------- | ---------------------------------------- |
| Authorization    | yes      | `Bearer $API_KEY`                        |
| Idempotency-Key  | no       | A unique key to dedupe retried requests. |
| Stripe-Version   | no       | Date-pinned version; defaults to account pin. |

### Body
| Field         | Type    | Required | Description                              |
| ------------- | ------- | -------- | ---------------------------------------- |
| `amount`      | integer | yes      | Positive amount in the smallest currency unit (e.g., cents). |
| `currency`    | string  | yes      | 3-letter ISO currency code, lowercase.   |
| `source`      | string  | yes*     | Card token, source ID, or customer ID. \*Required unless `customer` is given. |
| `customer`    | string  | yes*     | Existing customer ID.                    |
| `description` | string  | no       | Free-form text shown on receipts and the Dashboard. ≤ 1000 chars. |

## Response

**200 OK** — returns a `Charge` object.

```json
{
  "id": "ch_3O...",
  "object": "charge",
  "amount": 2000,
  "currency": "usd",
  "status": "succeeded",
  "created": 1715000000
}
```

## Errors

| Status | Code                | When it happens                          |
| ------ | ------------------- | ---------------------------------------- |
| 400    | invalid_request_error | Required field missing or malformed.   |
| 401    | authentication_error  | API key invalid or revoked.            |
| 402    | card_declined         | Card was declined. See `decline_code`. |
| 429    | rate_limited          | Account exceeded its rate limit.       |
| 500    | api_error             | Transient server error; safe to retry. |

→ Full error reference: [Errors](../errors)

## Code samples

```bash
curl https://api.example.com/v1/charges \
  -u "$API_KEY:" \
  -d amount=2000 \
  -d currency=usd \
  -d source=tok_visa
```

```python
import os, stripe
stripe.api_key = os.environ["API_KEY"]
charge = stripe.Charge.create(
    amount=2000, currency="usd", source="tok_visa"
)
```

## Idempotency
Pass `Idempotency-Key` to make this call safe to retry. Same key
within 24 hours returns the same response.

## Related
- [Refunds](../refunds) — refund a charge.
- [Disputes](../disputes) — handle a chargeback.
```

### Template — error catalog entry

```markdown
## rate_limited (429 Too Many Requests)

The account exceeded its rate limit. Response includes a
`Retry-After` header with the number of seconds to wait.

**Remediation:**
- Honor the `Retry-After` header; back off and retry.
- For sustained high traffic, request a rate-limit increase through
  your account dashboard.
- Implement exponential backoff with jitter for unattended retries.

**Common causes:**
- Bulk imports without throttling.
- Webhook retry storms after a long downtime.
- A misconfigured retry loop sending the same request without backoff.
```

### Template — Diátaxis-organized docs site structure

```
docs/
├── tutorials/                        # Learning-oriented
│   ├── getting-started.md           # The 10-minute first call
│   ├── building-a-webhook.md
│   └── first-payment-flow.md
├── guides/                          # Task-oriented (how-to)
│   ├── rotate-api-key.md
│   ├── handle-card-declines.md
│   ├── migrate-v1-to-v2.md
│   └── set-up-idempotency.md
├── reference/                       # Information-oriented
│   ├── api/
│   │   ├── charges.md
│   │   ├── customers.md
│   │   └── webhooks.md
│   ├── errors.md
│   ├── status-codes.md
│   └── api-versions.md
├── explanation/                     # Understanding-oriented
│   ├── why-prefixed-ids.md
│   ├── idempotency-model.md
│   ├── versioning-philosophy.md
│   └── webhook-retry-semantics.md
└── changelog.md                     # Versioning changelog
```

### Example — deprecation banner in a reference page

```markdown
> ⚠️ **Deprecated 2026-11-11. Sunsets 2027-05-11.**
>
> After **2027-05-11**, requests to this endpoint return `410 Gone`.
> Migrate to [`POST /v2/charges`](../v2/charges).
> See the [v1 → v2 migration guide](../guides/migrate-v1-to-v2).

# POST /v1/charges
```

### Example — OpenAPI snippet with rich descriptions

```yaml
paths:
  /v1/charges:
    post:
      summary: Create a charge
      description: |
        Charges a card, customer, or source for a fixed amount in
        the smallest currency unit (e.g., cents for USD).
      parameters:
        - in: header
          name: Idempotency-Key
          required: false
          description: |
            A unique key to make this request safe to retry. Same
            key within 24h returns the same response.
          schema:
            type: string
            example: "a1b2c3d4-uuid-..."
      requestBody:
        required: true
        content:
          application/x-www-form-urlencoded:
            schema:
              type: object
              required: [amount, currency]
              properties:
                amount:
                  type: integer
                  minimum: 1
                  description: Positive integer in the smallest currency unit.
                  example: 2000
                currency:
                  type: string
                  pattern: "^[a-z]{3}$"
                  description: 3-letter ISO currency code, lowercase.
                  example: "usd"
      responses:
        "200":
          description: Charge created.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Charge"
              example:
                id: "ch_3O..."
                object: "charge"
                amount: 2000
                currency: "usd"
                status: "succeeded"
        "402":
          description: Card was declined.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Error"
              example:
                error:
                  type: "card_error"
                  code: "card_declined"
                  decline_code: "insufficient_funds"
                  message: "Your card has insufficient funds."
```

## Anti-Patterns

- **Mixing Diátaxis quadrants on one page** — a reference page that drifts into tutorial-style narrative wastes the lookup-er's time and confuses the learner who needs the whole flow.
- **Code samples without runnable context** — missing imports, mystery variables, no example response. Always show the imports, the env-var sourcing, and the response shape.
- **Single-language code samples** — assumes the reader uses your favorite stack. Ship curl + at least one strongly-typed language + one dynamic language.
- **No error catalog** — error documentation scattered across endpoint pages with no central index. Errors are reusable concepts; document them once and link.
- **Deprecation messages only in the changelog** — the developer doesn't read the changelog when they hit the deprecated endpoint. The deprecation must be visible on the endpoint's reference page *and* in the response headers.
- **OpenAPI spec without `description` or `example` fields** — renders as a wall of types. Every parameter needs a one-line description; every schema needs an example.
- **Stale code samples** — samples that worked against an older version. Either auto-generate from a tested codebase or run the samples in CI.
- **No "Getting started" tutorial** — a 10-page reference index is not a substitute for a 10-minute "your first call" walkthrough.
- **Treating the docs as an afterthought** — Stripe famously made docs quality part of engineering performance review. Docs that ship after the API ship late, are inaccurate, and rot fast.
- **Versioning that isn't documented** — readers can't tell what version they're reading. Show the version selector on every page.

## Decision Heuristics

- **Which Diátaxis quadrant am I writing?** Ask: *Is the reader learning, working on a task, looking up a fact, or trying to understand?* If unsure, the page probably wants to be two pages.
- **URI vs header vs date versioning?** Public, broadly-consumed REST API → URI path is the most discoverable. Highly-evolving API with frequent micro-changes → date-pinned. Internal API with type-safe clients → header.
- **Inline error doc on the endpoint, or link to the catalog?** Catalog. The endpoint page lists *which* errors it returns; the catalog explains each one. Avoid duplication.
- **OpenAPI-generated or hand-written reference?** Generated, from a single source of truth (spec-first or code-first). Hand-written reference drifts within a sprint.
- **One tutorial or many?** One "Getting started" canonical tutorial that everyone walks through. Additional tutorials only for fundamentally different first-paths (e.g., "Getting started with webhooks" vs "Getting started with charges").
- **When to ship a deprecation?** Two-stage: announce deprecation 6-12 months before sunset. Ship `Deprecation` header + docs banner immediately. Ship `Sunset` header once the date is fixed. Don't sunset without both headers shipped for at least the announced window.
- **How many code-sample languages?** Minimum: curl + Node/Python + one strongly-typed language (Go, Java, .NET). Stripe ships 8; most teams can't maintain that many. Pick what you can keep tested.
- **Where does the "explanation" content go?** A separate top-level section, not woven into reference. Readers should be able to skip it; readers who need it should be able to find it.

## References

- [Diátaxis — official site by Daniele Procida (diataxis.fr)](https://diataxis.fr/) — canonical statement of the four-quadrant framework: tutorial / how-to / reference / explanation.
- [Stripe API Reference](https://docs.stripe.com/api) — the industry gold-standard three-pane layout, multi-language samples, and idempotency/versioning patterns.
- [Why Stripe's API Docs Are the Benchmark (apidog.com)](https://apidog.com/blog/stripe-docs/) — analysis of the Stripe doc patterns worth copying.
- [Twilio API Documentation](https://www.twilio.com/docs) — code-snippet repository and multi-language curl/Node/Python/Java patterns; companion gold-standard to Stripe.
- [RFC 8594 — The Sunset HTTP Header Field](https://datatracker.ietf.org/doc/html/rfc8594) — canonical spec for the `Sunset` header used in deprecation flows.
- [The Deprecation HTTP Header Field (IETF draft)](https://greenbytes.de/tech/webdav/draft-ietf-httpapi-deprecation-header-latest.html) — companion `Deprecation` header spec.
- [OpenAPI Specification 3.1 (Swagger)](https://swagger.io/specification/) — schema language that drives generated reference docs.
- [Fern — API documentation best practices guide](https://buildwithfern.com/post/api-documentation-best-practices-guide) — current (2026) best-practices roundup including OpenAPI workflows.
- [Zalando RESTful API Guidelines — Deprecation](https://github.com/zalando/restful-api-guidelines/blob/main/chapters/deprecation.adoc) — production reference for staged-deprecation policy and header use.
- [Simon Willison — The Diátaxis documentation framework](https://simonwillison.net/2021/Aug/21/diataxis/) — pragmatic walk-through with examples of mis-quadranted pages.
