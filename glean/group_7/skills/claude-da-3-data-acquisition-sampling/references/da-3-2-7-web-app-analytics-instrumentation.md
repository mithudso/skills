<!-- hub-reference-banner -->
> **Reference file — part of the `da-3-data-acquisition-sampling` hub.** Formerly the standalone `da-3-2-7-web-app-analytics-instrumentation` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-3-2-7-web-app-analytics-instrumentation
version: "1.1.0"
updated: "2026-05-30"
description: |
  Event instrumentation for web and mobile apps. Covers the full lifecycle from
  taxonomy design through SDK implementation: track plan authorship (Iteratively /
  Avo / RudderStack), event naming conventions, property schemas, identity
  management (anonymous IDs, user IDs, alias/merge), session boundaries, and
  consent/privacy wiring for GA4, Snowplow, PostHog, Segment, Amplitude,
  Mixpanel, Heap, and Adobe Analytics.

  TRIGGER: user is designing or auditing an analytics tracking plan; naming
  events or properties; implementing analytics.identify / analytics.track SDK
  calls; setting up anonymous-to-known identity resolution; configuring session
  start/end logic; wiring consent mode or cookie banners to an analytics SDK;
  choosing between Segment, Amplitude, Mixpanel, PostHog, Snowplow, Heap, GA4,
  or Adobe Analytics for instrumentation; asking about event governance, data
  quality validation, or Avo/RudderStack tracking-plan tooling; instrumenting
  iOS/Android mobile apps with analytics SDKs.

  SKIP: raw API ingestion without user-facing events → use da-3-2-6-apis-data-feeds.
  Privacy law compliance and GDPR legal obligations → use da-11-ethics-and-privacy.
  Statistical analysis of collected data → use da-5-exploratory-data-analysis or
  da-6-statistical-modeling.
related_skills:
  - da-3-2-6-apis-data-feeds
  - da-11-ethics-and-privacy
  - da-3-data-collection-acquisition
  - da-3-2-collection-methods
---

# 3.2.7 Web/App Analytics Instrumentation

## What this is

Analytics instrumentation is the practice of deliberately recording user
interactions and system signals from web pages and mobile apps so that product,
growth, and data teams can answer behavioral questions. The discipline spans four
concerns: (1) deciding *what* to record (taxonomy), (2) naming and shaping those
records consistently (conventions and schemas), (3) tying records to individuals
across sessions and devices (identity), and (4) honoring consent before firing
any of it (privacy wiring).

---

## 1. The Track Plan — Single Source of Truth

A **track plan** (also called a tracking plan or measurement spec) is a living
document — typically a spreadsheet, YAML file, or dedicated tool entry — that
declares every event your product emits, its properties, their types, and which
screen or code path fires it. ([Amplitude Docs](https://amplitude.com/docs/data/data-planning-playbook))

Track plans serve three teams at once:

- **Product / analytics**: defines the questions the data must answer
- **Engineering**: gives implementation details down to parameter names and types
- **Data governance**: provides the schema against which incoming events are
  validated before they reach the warehouse

### 1.1 Track-plan tooling

| Tool | Model | How it enforces the plan |
|---|---|---|
| **Avo** | SaaS; collaborative web UI | Generates typed SDK wrappers from the plan; lint rules flag naming violations at authoring time; publishes plan to RudderStack, Segment, or Amplitude ([Avo Docs](https://www.avo.app/docs/publishing/publishing/rudderstack)) |
| **Iteratively** (now part of Amplitude Data) | SaaS | Generates typed wrappers; integrates CI validation; merged into Amplitude's Data product ([Amplitude](https://amplitude.com/blog/analytics-tracking-practices)) |
| **RudderStack Tracking Plans** | Self-hosted or cloud | Proactive monitoring: blocks or flags events that violate the declared schema at ingestion ([RudderStack Docs](https://www.rudderstack.com/docs/data-governance/tracking-plans/)) |
| Spreadsheet / YAML | DIY | No enforcement; requires manual code review and discipline |

**Avo + RudderStack integration**: Avo publishes its tracking plan directly into
RudderStack's Tracking Plan API on every branch merge, closing the loop between
spec and enforcement. ([Avo Blog](https://www.avo.app/blog/optimize-your-event-driven-infrastructure-with-rudderstack-and-avo))

### 1.2 What belongs in a track plan entry

Each event row should include:

- Event name (canonical, exact-case version)
- Human description (one sentence: what user action triggers this)
- Required vs optional property list, each with: name, type, allowed values/regex
- Platform(s) where it fires (web, iOS, Android, server)
- Owner (team or person responsible for keeping it current)
- Status (planned / implemented / deprecated)

Track only events that answer a defined business question. Over-instrumentation
obscures answers just as readily as under-instrumentation; every event adds
pipeline cost and maintenance debt.
([Amplitude Blog](https://amplitude.com/blog/analytics-tracking-practices))

---

## 2. Event Taxonomy Design

A taxonomy is the classification scheme that organizes events into a coherent,
navigable structure. Without one, event lists become unsearchable noise after a
few dozen entries.

### 2.1 Object–Action (noun–verb) model

The most widely adopted taxonomy pattern names events as `[Object] [Past-tense
Verb]`, e.g., `Product Viewed`, `Checkout Completed`, `Subscription Cancelled`.
([Amplitude Blog](https://amplitude.com/blog/event-taxonomy))

Alternatives:
- `[Verb]_[Object]` in snake_case: `checkout_completed`, `button_clicked`
- `[Verb] [Object]` in present tense: `view product`, `complete checkout`

Pick one and enforce it across the entire product. Mixing conventions across
teams is the single largest source of duplicate event debt.
([Wudpecker Blog](https://www.wudpecker.io/blog/simple-event-naming-conventions-for-product-analytics))

### 2.2 Levels of specificity

Two anti-patterns:

- **Too generic**: `Button Clicked` with a `button_name` property — loses
  funnel-analysis capability; every query requires filtering
- **Too specific**: `Red Buy Button on Homepage Clicked` — explodes the event
  namespace, makes cohort analysis impossible

A practical middle ground: `[Feature Area] [Object] [Action]`, e.g.,
`Checkout Payment Method Selected`. Reserve hyper-specific events for
instrumentation of tightly scoped A/B tests.
([WarpDriven Analytics](https://warpdriven.ai/en/blog/industry-1/ecommerce-analytics-event-taxonomy-best-practices-2025-51))

### 2.3 Casing rules (pick one, enforce everywhere)

| Convention | Example | Common in |
|---|---|---|
| Title Case | `Product Viewed` | Amplitude, Mixpanel display names |
| snake_case | `product_viewed` | Segment, RudderStack wire format |
| camelCase | `productViewed` | Mobile SDKs, some CDP schemas |

Avo's global namespace enforces chosen casing, catches near-duplicates, and
flags convention violations at authoring time — not at query time.
([Avo Docs](https://www.avo.app/docs/data-design/best-practices/naming-conventions))

---

## 3. Property Schemas

Properties are key–value pairs that travel with each event and describe the
specific instance. They are the primary axis for segmentation and funnel
analysis.

### 3.1 Event properties vs user properties

- **Event properties** describe *this occurrence*: `price`, `currency`,
  `payment_method`, `item_count`. They do not persist after the event.
- **User properties** describe *the actor*: `plan_tier`, `signup_cohort`,
  `country`. They persist until explicitly changed and apply to all future
  events from that user.
  ([Amplitude Docs](https://amplitude.com/docs/data/data-planning-playbook))

Never send the same field as both an event property and a user property — the
analytic semantics differ and the confusion compounds downstream.
([Avo Docs](https://www.avo.app/docs/data-design/best-practices/naming-conventions))

### 3.2 Type discipline

| Data | Correct type | Wrong type |
|---|---|---|
| Price | Number (float) | String `"9.99"` |
| Currency | String ISO-4217 (`"USD"`) | Number |
| Timestamp | ISO-8601 string or Unix ms integer | Free-form date string |
| Boolean flags | Boolean | String `"true"` |
| IDs | String (preserves leading zeros) | Number |

Amplitude recommends a hard cap of **20 properties per event** and 20 user
properties per user to keep cardinality manageable.
([Amplitude Blog](https://amplitude.com/blog/analytics-tracking-practices))

### 3.3 Super properties / global context

Mixpanel, PostHog, and most CDPs allow "super properties" (also called context
properties or global traits) set once and automatically appended to every
subsequent event: `app_version`, `platform`, `locale`, `environment`. Register
these at SDK init rather than repeating them per event.
([Mixpanel Docs](https://docs.mixpanel.com/docs/tracking-methods/sdks/javascript#setting-super-properties))

### 3.4 PII prohibition

Never send personally identifiable information (name, email address, phone,
government ID) as event properties or user traits directly to analytics
platforms. If an identifier is needed for joining, send a hashed or opaque
internal ID. ([Amplitude Blog](https://amplitude.com/blog/analytics-tracking-practices))

---

## 4. Identity Management

Analytics identity answers: "Is this the same person across sessions, devices,
and sign-in states?"

### 4.1 The three-state model

1. **Anonymous (pre-identify)**: SDK generates a random UUID stored in a
   first-party cookie or localStorage. Segment calls this `anonymousId`. Every
   event carries it from the first page load. It persists across sessions until
   cookies are cleared or `reset()` is called.
   ([Saber Glossary](https://www.saber.app/glossary/anonymous-id))

2. **Identified**: User authenticates. `analytics.identify(userId, traits)`
   is called. From this point, events carry both `anonymousId` and `userId`.
   The platform links all prior anonymous events to the user where supported.
   ([Statsig Blog](https://www.statsig.com/blog/demystifying-identity-resolution))

3. **Aliased / merged**: When a user crosses systems (e.g., created an account
   on mobile first, now on web), `analytics.alias(newId, previousId)` merges
   two identity records into one canonical profile. Segment and Mixpanel support
   this; Amplitude uses a merge API. The merge is permanent and irreversible in
   Segment and Mixpanel; Amplitude allows identity un-merge via its admin API.
   Call it exactly once.

### 4.2 Session boundaries

A **session** groups events into a single continuous visit. Platforms define
session end differently:

| Platform | Default session timeout |
|---|---|
| GA4 | 30 minutes of inactivity (configurable) |
| Amplitude | 30 minutes of inactivity |
| Mixpanel | Rolling 30-minute window |
| Snowplow | User-defined; event-driven session heartbeat available |
| PostHog | 30-minute inactivity; configurable in SDK |

Session IDs should be generated client-side as UUIDs and attached as an event
property (`session_id`) when platforms do not auto-manage sessions. This is
required for server-side or hybrid instrumentation where the SDK cannot
maintain a browser cookie.

`analytics.reset()` in Segment (and equivalent calls in other SDKs) clears both
`anonymousId` and `userId`, starting a fresh anonymous session. Call it on
explicit sign-out, never on page navigation.

### 4.3 Cross-device identity

No client-side mechanism reliably links the same user across browsers and
devices. Options in increasing reliability:

1. **Email-based matching** at the warehouse (requires consent and an email
   capture event) — lowest reliability; depends on users providing the same
   email in every context
2. **Probabilistic matching** offered by Segment's Identity Resolution, Adobe
   Analytics's CDA, or Mixpanel's Identity Merge — confidence-weighted, not
   exact; infers identity from shared signals (IP, user-agent, timing)
3. **Deterministic matching** via log-in events synced to a customer data
   platform (CDP) — highest reliability; links records only when the user
   authenticates with a known identifier

---

## 5. Consent and Privacy Wiring

Instrumentation must not fire until the user's consent state is known and, where
required by law, affirmatively granted. This section covers the integration
pattern; the legal obligations live in `da-11-ethics-and-privacy`.

### 5.1 GA4 Consent Mode v2

Google's Consent Mode v2 (mandatory for EU users as of March 2024) passes
consent signals — `analytics_storage`, `ad_storage`, `ad_user_data`,
`ad_personalization` — to the GA4 tag before any hit fires. When consent is
denied, GA4 fires cookieless "pings" for modeling purposes rather than dropping
data entirely. Server-side consent mode (via GTM server container) is more
reliable than client-side because browser extensions cannot suppress it.
([SecurePrivacy](https://secureprivacy.ai/blog/server-side-consent-mode-for-ga4-how-to-track-analytics-while-respecting-privacy))

### 5.2 Snowplow consent events

Snowplow ships first-class consent event schemas (part of its Enhanced Consent
plugin): `consent_granted`, `consent_withdrawn`, `consent_expired`,
`consent_selected`. These events are tracked alongside behavioral data in the
same pipeline, giving a complete audit trail.
([Snowplow Docs](https://docs.snowplow.io/docs/sources/web-trackers/tracking-events/consent-gdpr/))

### 5.3 PostHog

PostHog does not capture any events until the consent state is resolved. On
denial it records a privacy-preserving hash count (not linkable to a person)
rather than dropping the data point entirely — preserving aggregate funnel
metrics without violating user choice. PostHog Cloud EU hosts exclusively in
Frankfurt, keeping EU data in-region.
([PostHog Docs](https://posthog.com/docs/privacy/gdpr-compliance))

Cookieless mode (PostHog `persistence: 'memory'`) stores no cookies or
localStorage; anonymous IDs reset on page load. This sacrifices cross-session
attribution for maximum privacy compliance.

### 5.4 Segment and other CDPs

Segment does not implement consent gating natively; you must gate the
`analytics.load()` call or use a middleware to block destinations selectively
based on consent category. Typical pattern:

```js
if (consentGranted('analytics')) {
  analytics.load(WRITE_KEY);
}
// Or use Segment's consent-aware middleware to conditionally forward
// events to specific destinations based on consent categories.
```

### 5.5 General consent integration pattern

1. Load your consent management platform (CMP) before the analytics SDK.
2. On `consentGranted` callback: initialize SDK, set persistent `anonymousId`.
3. On `consentDenied`: skip SDK init or run in cookieless / sampling-only mode.
4. On `consentWithdrawn`: call `analytics.reset()`, flush local storage, signal
   the backend to delete the user's data record.
5. Store consent state in a first-party cookie with a short expiry (≤13 months
   per GDPR guidance); re-prompt on expiry.

---

## 5b. Mobile and Server-Side Instrumentation Notes

### Mobile (iOS / Android)

Mobile SDKs manage `deviceId` differently from web SDKs. On iOS, Amplitude and
Mixpanel default to IDFV (Identifier for Vendor) as the device ID; on Android
they use a generated UUID persisted to shared preferences. Key differences from
web:

- **App backgrounding resets session clocks.** Most mobile SDKs pause the
  session timer when the app enters background and resume on foreground; if the
  gap exceeds the timeout, a new session starts even with no user action.
- **Reinstall resets `deviceId`.** A fresh install generates a new device ID,
  breaking the historical continuity for that device. Use `userId` (from
  authenticated state) as the stable cross-install identifier.
- **Offline queuing.** Mobile SDKs buffer events locally and flush on next
  connectivity. Timestamps are set at event creation, not at flush — ensure your
  warehouse ingestion pipeline respects `time` (event creation) vs
  `server_received_time` (flush time).

### Server-side instrumentation

Server-side calls bypass the browser entirely, which eliminates ad-blocker
suppression but requires the caller to supply identity context explicitly:

- Pass `userId` (from your auth session) on every server-side `track` call.
- Generate a `session_id` server-side and attach it as an event property; the
  server has no browser cookie to rely on.
- Never send server-side events with only an `anonymousId` unless the user has
  not authenticated — the anonymous ID from the browser and the server will
  diverge unless explicitly threaded through.
- For Segment, use the HTTP API or a server-side source (Node, Python, Go SDKs);
  for Amplitude, use the HTTP API v2 directly.

---

## 6. Platform Quick Reference

| Platform | Identity model | Session handling | Consent support | Open-source? |
|---|---|---|---|---|
| **GA4** | ClientID (cookie) + UserID | Automatic, 30 min | Consent Mode v2 | No |
| **Segment** | anonymousId + userId + alias | No native session; attach `session_id` | Middleware / destinations | Partial (Analytics.js) |
| **Amplitude** | deviceId + userId | Automatic, 30 min | Opt-out flag | No |
| **Mixpanel** | distinct_id; alias/identify merge | Rolling 30 min | EU residency option | JS SDK OSS |
| **PostHog** | anonymousId + distinctId | 30 min, configurable | Built-in consent gating | Yes |
| **Snowplow** | domain_userid (cookie) + user_id | Configurable via heartbeat | Enhanced Consent plugin | Yes |
| **Heap** | Identity auto-captured (retroactive) | Automatic | GDPR opt-out API | No |
| **Adobe Analytics** | Experience Cloud ID (ECID) | Server-side session | CMP integration + CDA | No |

---

## 7. Implementation Checklist

Use this when instrumenting a new feature or auditing existing instrumentation.

**Before writing any code**
- [ ] Is the event in the track plan with description, properties, and owner?
- [ ] Does the event name follow the agreed naming convention?
- [ ] Are all property types declared (not inferred)?
- [ ] Does any property risk carrying PII?

**During implementation**
- [ ] Anonymous ID is generated at SDK init, before any events fire.
- [ ] `identify()` is called exactly once per sign-in (not on every page load).
- [ ] `alias()` / merge is called at most once per user lifetime.
- [ ] Session ID is attached where the platform does not manage it automatically.
- [ ] Super properties / global context registered at init, not per event.

**Consent**
- [ ] SDK init is gated on consent resolution — no events before consent is known.
- [ ] Consent withdrawal triggers `reset()` and backend deletion signal.
- [ ] Consent events are themselves tracked for audit purposes.

**Validation**
- [ ] Events fire in development and appear in the platform's debugger/live view.
- [ ] Required properties are present on every event instance.
- [ ] No events fire in test/staging environments against production write keys.

---

## 8. Common Pitfalls

**Tracking too much**: every event adds pipeline cost and query noise. Start with
the five to ten events that directly answer current product questions.
([Amplitude Blog](https://amplitude.com/blog/event-taxonomy))

**Inconsistent naming across platforms**: if web sends `product_viewed` and
mobile sends `Product Viewed`, cross-platform funnels break silently. Enforce a
single canonical name in the track plan and generate SDK wrappers (Avo) rather
than hand-coding calls.

**Calling `identify()` before consent is obtained**: passing a user ID to an
analytics platform before consent constitutes personal-data processing under
GDPR without a legal basis.

**Aliasing repeatedly**: calling `alias` more than once for a user merges
unrelated identity records in platforms that apply the call globally (Mixpanel,
older Segment behavior). Call it exactly once, immediately after the first
sign-up event.

**Missing `reset()` on sign-out**: without a reset, the next user on a shared
device inherits the previous user's `anonymousId`, corrupting funnel analysis.

---

## Sources

1. [Amplitude Data Planning Playbook](https://amplitude.com/docs/data/data-planning-playbook) — canonical track-plan structure and taxonomy guidance
2. [Amplitude Blog: Analytics Tracking Best Practices](https://amplitude.com/blog/analytics-tracking-practices) — property limits, PII rules, naming discipline
3. [Amplitude Blog: Event Taxonomy Foundation](https://amplitude.com/blog/event-taxonomy) — object–action model and taxonomy anti-patterns
4. [Avo Docs: Naming Conventions](https://www.avo.app/docs/data-design/best-practices/naming-conventions) — casing, format, tense, vocabulary standardization
5. [Avo Blog: RudderStack Integration](https://www.avo.app/blog/optimize-your-event-driven-infrastructure-with-rudderstack-and-avo) — track-plan publishing workflow
6. [RudderStack Docs: Tracking Plans](https://www.rudderstack.com/docs/data-governance/tracking-plans/) — ingestion-layer schema enforcement
7. [Statsig: Demystifying Identity Resolution](https://www.statsig.com/blog/demystifying-identity-resolution) — anonymous-to-known identity state machine
8. [Saber Glossary: Anonymous ID](https://www.saber.app/glossary/anonymous-id) — anonymousId mechanics, persistence, and privacy
9. [SecurePrivacy: Server-Side Consent Mode for GA4](https://secureprivacy.ai/blog/server-side-consent-mode-for-ga4-how-to-track-analytics-while-respecting-privacy) — GA4 Consent Mode v2 architecture
10. [Snowplow Docs: Consent & GDPR Tracking](https://docs.snowplow.io/docs/sources/web-trackers/tracking-events/consent-gdpr/) — Enhanced Consent plugin event schemas
11. [PostHog Docs: GDPR Compliance](https://posthog.com/docs/privacy/gdpr-compliance) — consent gating, cookieless mode, EU hosting
12. [Wudpecker: Event Naming Conventions](https://www.wudpecker.io/blog/simple-event-naming-conventions-for-product-analytics) — practical naming patterns comparison
13. [WarpDriven: Ecommerce Event Taxonomy 2025](https://warpdriven.ai/en/blog/industry-1/ecommerce-analytics-event-taxonomy-best-practices-2025-51) — specificity spectrum, funnel-verb patterns
