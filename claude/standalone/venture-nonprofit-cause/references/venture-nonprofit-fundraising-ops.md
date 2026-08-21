<!-- hub-reference-banner -->
> **Reference file — part of the `venture-nonprofit-cause` hub.** Formerly the standalone `venture-nonprofit-fundraising-ops` skill.
> Sibling topics in this family are now reference files under the hubs (`venture-business`, `venture-nonprofit-cause`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: venture-nonprofit-fundraising-ops
description: >-
  Hands-on fundraising & marketing-OPERATIONS layer for a lean cause/nonprofit (built for an NC
  organ-donation org driving donor-registry sign-ups, general to any small charity) — the build-and-run
  layer beneath the strategy sibling. TRIGGER: building, structuring, or fixing a Google Ad Grants
  account — ad-group structure, keyword themes, negatives, the 5% CTR rule, the Maximize-Conversions
  $2-CPC-cap exception, conversion tracking, Smart Bidding, or a suspension; email lifecycle automation
  — welcome series, donor/registrant journeys, segmentation, cadence, send-time, SPF/DKIM/DMARC
  deliverability, or choosing a lean ESP; donor CRM selection (Zeffy, Little Green Light, Bloomerang,
  Neon, DonorPerfect); measurement — GA4 key events, UTM conventions, registry-signup conversion
  events, offline/"dark"-conversion attribution limits, or a simple holdout/lift test. SKIP: cause
  STRATEGY, donor-journey DESIGN, storytelling ethics, Ad Grants ELIGIBILITY →
  venture-cause-nonprofit-marketing; for-profit marketing + local SEO + CAN-SPAM/TCPA →
  venture-marketing-strategy-local-seo; design/Canva/print → venture-canva-founder-brand-stack;
  501(c)(3)/charitable-solicitation → venture-nc-nonprofit-formation; donation-system/registry mechanics
  → venture-organ-donation-system; marketing-mix-modeling → da-applied-and-communication.
category: personal-venture
tags: [venture, nonprofit, fundraising-ops, google-ad-grants, crm, email-marketing]
whenToUse:
  - Structuring or repairing a Google Ad Grants account — campaigns, ad groups, ads, keywords, negatives, sitelinks
  - Hitting or recovering the 5% CTR rule, switching to Maximize Conversions to unlock bids above the $2 CPC cap, or fixing a suspended grant
  - Setting up Google Ads / GA4 conversion tracking so a registry sign-up or donation counts as a meaningful conversion
  - Building email lifecycle automations — welcome series, new-donor/registrant journeys, lapsed re-engagement — with segmentation and cadence
  - Fixing email deliverability (SPF/DKIM/DMARC, spam-complaint rate) or choosing a lean nonprofit ESP
  - Selecting a donor/constituent CRM for a 1–3 person team (Zeffy, Little Green Light, Bloomerang, DonorPerfect, Neon)
  - Instrumenting measurement — GA4 key events, a UTM naming convention, conversion events for registry sign-ups
  - Measuring real impact under attribution limits — offline/"dark" conversions, a simple holdout or geo-lift test
triggers:
  - Google Ad Grants account structure, ad groups, keywords, negative keywords, CTR, suspension
  - Maximize Conversions, Smart Bidding, $2 CPC cap, conversion tracking setup for the grant
  - email welcome series, donor journey, automation, segmentation, cadence, send time
  - email deliverability, SPF, DKIM, DMARC, spam complaint rate, nonprofit ESP choice
  - donor CRM selection, Bloomerang, DonorPerfect, Neon CRM, Little Green Light, Zeffy
  - GA4 key events, conversion events, UTM parameters, registry sign-up tracking
  - offline / dark conversion attribution, holdout test, geo lift, incrementality
  - how do I set up, structure, or fix Ad Grants, email automation, CRM, or analytics
---

# Nonprofit Fundraising & Marketing Operations

The **operational execution layer** for a lean cause nonprofit's fundraising/marketing — the buttons,
structures, and settings beneath the strategy. **Strategy lives in the siblings; this skill builds
the machine.** Use it when the question is *how do I set this up, structure it, or fix it*, not *what
should our campaign say* (→ `venture-cause-nonprofit-marketing`) or *how does local SEO work*
(→ `venture-marketing-strategy-local-seo`).

Worked example threaded throughout: **"Register Carolina" — a North Carolina organ-donor-registry
sign-up drive run by a 2-person nonprofit**, primary conversion = a completed donor-registration
referral (offsite at the state/Donate Life registry), secondary = newsletter signup.

> **Rules change. Verify every Ad Grants / Gmail / GA4 / ESP rule against the official current docs
> before you act on it.** Program rules below carry as-of dates; treat them as a starting checklist,
> not gospel. See **Disclaimer**.

---

## 1. Google Ad Grants — Operational Playbook

The Ad Grant gives an approved 501(c)(3) up to **$10,000/month** in Google Search ads (≈ **$329/day**)
[Whole Whale; Getting Attention, as of 2026-02]. It is *constrained* free search inventory: text ads
only, Search Network only, and a thicket of compliance rules that **auto-deactivate** non-compliant
accounts. Eligibility and whether to use the grant at all are strategy (→ sibling); below is how to
build and keep the account.

### 1.1 Account & campaign architecture

Mirror your site's intent themes. A clean starting structure for Register Carolina:

```
Account
├── Campaign: "Become a Donor" (registration intent)
│   ├── Ad Group: register organ donor      → /register
│   ├── Ad Group: how to sign up donor       → /register
│   └── Ad Group: donate life NC             → /register
├── Campaign: "Donation Myths / Learn"  (education intent)
│   ├── Ad Group: organ donation myths       → /myths
│   └── Ad Group: who can donate              → /eligibility
└── Campaign: "Ways to Give"  (support intent)
    ├── Ad Group: donate to organ nonprofit  → /donate
    └── Ad Group: volunteer ambassador        → /volunteer
```

**Hard structural minimums** (suspension triggers if missing) [Getting Attention 2026-02; Big Sea;
official Ad Grants policy `support.google.com/nonprofits/answer/9314402`]:

- **≥ 2 ad groups per campaign** (no flat 1-ad-group campaigns).
- **≥ 2 active text ads per ad group** (3 responsive search ads is better — more headline rotation).
- **≥ 2 sitelink assets** at account or campaign level (e.g., *Register* · *Myths* · *Volunteer* · *Donate*).
- Keep each ad group **tightly themed** (3–15 closely related keywords) so ads stay relevant → higher
  Quality Score → higher CTR.

### 1.2 Keyword themes + the negative-keyword discipline

**Keyword rules** [Getting Attention 2026-02; official policy]:
- **No single-word keywords** — except your **brand name**, **medical/condition terms**, and recognized
  charitable terms like *donate*.
- **No "overly generic" keywords** (e.g., `free videos`, `e-books`, bare `students`).
- **No keyword with Quality Score 1 or 2** — pause/remove; keep QS **≥ 3**. Audit weekly.

**Keyword themes for the drive:** `register as organ donor`, `how to become organ donor NC`,
`sign up donate life`, `organ donor registry north carolina`, `is organ donation safe`,
`organ donation myths`. Use mostly **phrase** and **exact** match; broad match leaks irrelevant
traffic that **tanks CTR** (the #1 deactivation cause).

**Negative keywords are CTR insurance — build the list before launch and prune monthly:**
- **Job/career seekers:** `jobs`, `salary`, `careers`, `volunteer hours` (if you don't want them).
- **Wrong donation type:** `blood`, `plasma`, `marrow`, `sperm`, `egg`, `hair`, `car`, `clothes`,
  `food bank` (keep your drive about organ/eye/tissue).
- **Researchers/students writing papers:** `essay`, `statistics pdf`, `quizlet`, `definition`.
- **Commercial/medical-service intent:** `transplant cost`, `surgeon`, `near me clinic`.
- Run the **Search Terms report weekly**; every irrelevant query that got a click is a future negative.

### 1.3 The 5% CTR survival rule — and the Maximize-Conversions escape hatch

**The rule:** maintain **account-wide CTR ≥ 5% every month**. Miss it **two consecutive months →
temporary deactivation** [official policy 9314402; Getting Attention 2026-02]. CTR = clicks ÷
impressions across the whole account for the month. This is the single most common reason small grant
accounts die.

**How to stay above 5%:** tight ad-group themes; phrase/exact match; aggressive negatives; pause
low-CTR keywords; write specific, intent-matched RSAs; **pause underperformers rather than let them
drag the account average**.

**The bidding mechanic that matters most** — the **$2.00 max CPC cap and its exception:**
- Standard/manual bid strategies (Manual CPC, Enhanced CPC, Maximize Clicks, Target CPA) are **capped
  at a $2.00 max CPC** — often too low to win competitive auctions, starving the account of clicks.
- **Switching to `Maximize Conversions` (or `Maximize Conversion Value` / Target ROAS) REMOVES the $2
  cap**, letting Google bid **$4–$12+** per click and compete normally [Media Cause; AboveX Digital;
  Whole Whale; Digital Tabby — multi-source, as of 2026]. This is *the* lever that makes an Ad Grant
  actually spendable.
- **Prerequisite:** Maximize Conversions requires **valid conversion tracking** (see 1.4). Accounts
  created after ~April 2019 are effectively expected to run conversion-based **Smart Bidding** anyway
  [Whole Whale]. So: **set up conversion tracking first, then switch to Maximize Conversions.**

> **Operator takeaway:** Tight structure + negatives keep you **above 5% CTR**; Maximize Conversions
> + working conversion tracking let you **escape the $2 cap and actually spend the grant**. You need both.

### 1.4 Conversion tracking setup

Policy requires **≥ 1 meaningful conversion per month**, and the conversion must be a **genuine action**
(donation, registry sign-up, volunteer form, newsletter signup) — *not* every pageview or click
[official policy 9314402; Getting Attention 2026-02]. Two clean paths:

1. **GA4 key event → import to Google Ads** (recommended for a lean org; one analytics source of
   truth — see §4). Mark the registry-referral / thank-you event as a **key event** in GA4, then
   import it as a conversion into the linked Google Ads account.
2. **Google Ads tag directly** via Google Tag / GTM firing on a thank-you page or button click.

For Register Carolina, the registry completion happens **offsite** (state/Donate Life registry you
don't control). Workarounds, best→worst: (a) a **confirmation/thank-you page back on your own domain**
after the referral; (b) **outbound-click tracking** on the "Go to the registry" button (a *proxy*
conversion — counts the intent, not the completion); (c) newsletter signup as a tracked secondary
conversion so the account always has ≥ 1 conversion. Document which one you're counting so you don't
fool yourself.

### 1.5 Common suspension / deactivation causes (the deactivation checklist)

[official policy 9314402; Nonprofit Megaphone; Big Sea; AboveX — as of 2026]
- **CTR < 5% for two consecutive months** (most common).
- **Missing structure:** < 2 ad groups/campaign, < 2 ads/ad group, < 2 sitelinks.
- **Banned keywords:** single-word, overly generic, or **QS 1–2** left active.
- **No conversion tracking** / no meaningful conversion reported in a month.
- **Mission/relevance violations:** ads or landing pages off-mission, low-quality site,
  commercial/deceptive content, broken or slow landing pages.
- **Geo-targeting** not set to where you actually serve.
- **Not completing the annual program survey** (sent to the account login email).
- Google sends a **warning ≥ 7 days before** a suspension action for repeat violations — watch that inbox.

**Reactivation:** fix the violation, then request reactivation through the Ad Grants support form /
the in-account flow. Don't just re-enable; remediate first or it recurs.

---

## 2. Email Lifecycle Automation

Email is the highest-ROI owned channel for nonprofits: nonprofits raised **~$2.40 per subscriber** in
email revenue and **~$54 per 1,000 fundraising messages** in 2025 [M+R Benchmarks 2025]. The strategy
skill owns *donor-journey design and message ethics*; this section owns the **automation plumbing,
segmentation mechanics, cadence, and deliverability**.

### 2.1 The core automated journeys (build these first)

| Journey | Trigger | Sequence (lean version) |
|---|---|---|
| **Welcome series** | new newsletter/registrant subscribe | 1) instant welcome + what to expect, 2) the mission/impact story (day 2–3), 3) one concrete action — register / share / give (day 5–7) |
| **New-donor thank-you** | first donation | 1) instant tax-receipt + heartfelt thanks, 2) "here's your impact" (day 3–5), 3) soft second-touch / steward (day 14+) |
| **Registrant-confirmation** | completed/clicked registry referral | 1) congrats + "you're a donor" identity reinforcement, 2) "tell your family" (most registries need next-of-kin awareness), 3) become an ambassador / share |
| **Lapsed re-engagement** | no open/click in 6–12 mo | 1) "we miss you" + impact, 2) preference/frequency choice, 3) last-touch before sunsetting |

**Why welcome first:** welcome emails average **~80% open rates** and welcome-series emails pulled a
**1.6% CTR vs 0.59%** for standard fundraising emails in 2025 — roughly **3×** [CauseVox; M+R 2025].
New subscribers are at peak interest; capture it automatically.

### 2.2 Segmentation (the minimum useful cuts)

Don't over-engineer. For a lean org, segment by:
- **Role:** donor vs registrant vs volunteer vs general subscriber (different asks).
- **Donor value/recency:** first-time vs repeat vs lapsed; major vs small (drives ask size and
  stewardship depth).
- **Engagement:** opened/clicked recently vs dormant (protects deliverability — see 2.4).
- **Behavior:** event attendees, advocacy-action takers, registry-referral clickers.

Build segments as **dynamic/auto-updating** lists on a tag/custom-field, not hand-curated lists.

### 2.3 Cadence & send-time

- **Frequency:** nonprofits sent **~50 emails/subscriber/year** on average in 2025 (~31 of them
  fundraising) [M+R 2025] — i.e., roughly **weekly**. A lean org can do well at **2–4 sends/month**
  (1 newsletter/impact + 1–2 asks) plus the automated journeys. **Let subscribers self-select
  frequency** via a preference center.
- **Send-time:** there is **no universal best time**; benchmark "best days" are weak signals. **Use
  your ESP's send-time optimization** and **A/B test** your own list. Start mid-morning on a weekday,
  then test. [UNVERIFIED — specific "best day/time" claims vary by source and audience; treat as a
  hypothesis to test, not a rule.]
- **Around a registry drive:** sequence announce → story/impact → reminder → last-call; don't blast
  the whole list on the same day at the same hour.

### 2.4 Deliverability (the part that silently kills campaigns)

As of 2024–2025, **Gmail/Yahoo/Microsoft reject unauthenticated bulk mail outright** [PowerDMARC;
Red Sift; Chronos — as of 2026]. Non-negotiable setup:

- **SPF** — publish a DNS TXT record authorizing your sending IPs/ESP.
- **DKIM** — enable signing in your ESP and publish the DKIM key in DNS (cryptographic signature).
- **DMARC** — publish a DMARC record; **`p=none` is an acceptable start**, but progress toward
  `p=quarantine` → `p=reject` [PowerDMARC; Red Sift].
- **Spam-complaint rate:** keep **below 0.3%** (Gmail's enforcement ceiling); **aim for < 0.1%**
  [Chronos 2026]. One-click **list-unsubscribe** is required for bulk senders.
- **List hygiene:** authenticate the sending domain (not just the ESP's), warm new domains slowly,
  **suppress chronically dormant addresses** (they drag complaint/bounce rates and reputation).

> CAN-SPAM / consent law (who you may legally email, opt-out timing) is **compliance → live in
> `venture-marketing-strategy-local-seo`**, not here.

### 2.5 Lean nonprofit ESPs

Pick on: free/discount tier, automation/journey builder, segmentation, deliverability, and whether it
**writes back to your CRM**. [Zeffy; Constant Contact; softailed — as of 2026]

| ESP | Lean-org fit | Notes |
|---|---|---|
| **MailerLite** | Best budget Mailchimp-style pick | Clean automation builder, usable free tier, ~30% nonprofit discount |
| **Brevo (ex-Sendinblue)** | Volume-based pricing | Free tier by *sends* not contacts — good for big list/low frequency |
| **Mailchimp** | Familiar, feature-rich | Free tier exists but **per-contact pricing scales fast** |
| **Constant Contact** | Hand-holding / phone support | No free plan; 20–30% nonprofit prepay discount |
| **HubSpot (free)** | CRM + email together | Free email + CRM, ~2,000 sends/mo, unlimited contacts |
| **Zeffy** | If it's also your donation tool | $0 model; email + donation forms in one (see §3) |

**Decision rule:** if email and donations should share one contact record and you're truly tiny,
favor an **all-in-one (Zeffy)**; if you need richer automation/segmentation, favor **MailerLite/Brevo**
and sync to your CRM.

---

## 3. Donor / Constituent CRM Selection

A donor CRM is the **system of record** for people, gifts, and interactions — distinct from your ESP
(email) and your ad accounts. A lean org needs: contact + gift history, **donation forms that
auto-write to the record**, basic segmentation/tags, receipting, simple reports, and a path to grow.
Don't buy enterprise fundraising suites you can't staff.

| CRM | Pricing (as of 2026) | Best for | Watch-outs |
|---|---|---|---|
| **Zeffy** | **$0 at every tier** (optional donor-paid tips fund it) | Volunteer-run / 1–3 people; donation forms + basic CRM + email in one | Lighter CRM depth; you trade features for free [Zeffy; LiveImpact] |
| **Little Green Light (LGL)** | ~**$45/mo** (~$486/yr) + card fees | Simplicity + affordability; classic small-shop donor DB | Fewer built-in marketing automations [Zeffy compare; G2] |
| **Bloomerang** | from ~**$125/mo** (+~1% platform fee on its giving tools) | **Donor retention** focus; clear health/retention dashboards | Pricier; more than a micro-org needs [Bloomerang; LiveImpact] |
| **Neon CRM (Neon One)** | from ~**$99/mo**, tiered | Growing org wanting events/memberships/automation | Tiered features; cost climbs with needs [Neon One] |
| **DonorPerfect** | **per-user**, quote-based | Mature, scalable, lots of integrations | Per-seat pricing; heavier to learn [LiveImpact] |

**Lean-org decision path:**
1. **Truly tiny / no budget / want one tool for forms+CRM+email →** start with **Zeffy** (free) and
   keep it until you outgrow it.
2. **Want a real donor database, modest budget, simplicity over features →** **Little Green Light**.
3. **Retention is the priority and budget exists →** **Bloomerang**.
4. **Scaling fast — events, memberships, automation →** **Neon** or **DonorPerfect**.

**Integration reality:** your **CRM ↔ ESP ↔ donation forms ↔ GA4** should connect. The fewer
hand-exports, the fewer errors. All-in-ones (Zeffy, Bloomerang's suite) reduce glue; best-of-breed
stacks need native integrations or a sync (e.g., Zapier). Confirm the integration exists *before* you
commit — this is the #1 thing lean orgs get stuck on later.

---

## 4. Measurement & Instrumentation (GA4 + UTMs + lift)

You can't optimize the funnel you can't see. Build the **minimum reliable measurement** and respect
its limits — especially for **offline / "dark" conversions** like a registry sign-up you can't track
to completion. (Marketing-mix-modeling depth → `da-applied-and-communication`.)

### 4.1 GA4 conversion events (key events)

In GA4 (2026), the model is **Event → Key Event → Conversion** [official GA4 docs
`support.google.com/analytics/answer/13965727`]:
- **Any event can be marked a Key Event** in Admin → Key events.
- A **Conversion** is created from a key event and gives a consistent measure across **GA4 and Google
  Ads** (this is what feeds Smart Bidding in §1.4).
- **Marking is NOT retroactive** — it only counts from the toggle date. Set it up **before** a campaign,
  not after [Optimize Smart; Netpeak — as of 2026].

For Register Carolina, define key events: `registry_referral_click` (outbound to the registry),
`registration_confirmed` (if you can get a return/thank-you page), `newsletter_signup`,
`donation_complete`. Mark the meaningful ones as key events → import to Google Ads.

### 4.2 UTM naming convention (write it down, enforce it)

UTMs are how GA4 attributes traffic to a campaign. The five (+id) parameters:
`utm_source`, `utm_medium`, `utm_campaign`, `utm_content`, `utm_term` (+ `utm_id` for reconciliation)
[Blue Hills Digital; UTM.io — as of 2026]. **Consistency beats cleverness** — pick one casing and
stick to it:

- **Rules:** **all lowercase**, **no spaces** (use `_` or `-`), one documented vocabulary.
- **Convention:** `utm_campaign={year}_{initiative}` e.g. `2026_register_carolina`.
- **Examples:**
  - Email ask: `?utm_source=newsletter&utm_medium=email&utm_campaign=2026_register_carolina&utm_content=cta_button`
  - Facebook post: `?utm_source=facebook&utm_medium=social&utm_campaign=2026_register_carolina&utm_content=story_video`
  - Printed flyer QR: `?utm_source=flyer&utm_medium=qr&utm_campaign=2026_register_carolina`
- **Cross-domain gotcha:** if a donation/registry page is on a **different domain** (a processor or the
  state registry), configure **cross-domain measurement** or those conversions show up as **referrals
  from the processor** and your email/ads/UTM credit is lost [RallyUp; Kissmetrics — as of 2026].

### 4.3 Attribution limits & the "dark"/offline conversion problem

GA4 is good at **digital clicks** but **cannot stitch multi-touch, cross-device, and offline
journeys** — a social→email→offline path gets mis-credited; offline actions (a call, an in-person
sign-up, a registry completion you can't tag) aren't captured natively [Kissmetrics; Improvado — as
of 2026]. For an organ-donor drive, the *true* outcome (a completed state-registry record) is often
**unobservable** to your analytics. Mitigations:
- Track the **best available proxy** (outbound click, return thank-you page) and **label it honestly**
  as a proxy, not a confirmed registration.
- Add an **"How did you hear about us?"** field on forms to capture self-reported/offline attribution.
- Use **offline conversion import** (upload CRM/registry outcomes back to Google Ads) only if you can
  actually get the data — usually you can't for a third-party registry. [UNVERIFIED for the specific
  case of importing third-party registry completions — depends entirely on the registry's data-sharing,
  which most don't offer.]

### 4.4 A simple holdout / lift test (proving incrementality on a budget)

Attribution tells you *which click got credit*; a **holdout** tells you whether the campaign **caused
new outcomes** [Improvado; Measured; fusepoint — as of 2026]. The gold standard is a randomized
controlled comparison; lean practical versions:

- **Audience holdout (email/digital):** randomly **hold back ~10%** of an eligible segment from a
  campaign; compare conversion/registry-click rate **exposed vs held-out**. Cheap, runs inside most ESPs.
- **Geo lift (for offline/broadcast/flyers):** run the push in some **NC regions/counties** and not in
  matched others; compare registry-signup or web-conversion rates. Good when you can't randomize people.
- **Window:** measure over **~21–30 days** (and 4–12 weeks for slower/offline effects) — conversions
  lag exposure, so a too-short window understates lift [Improvado; Triple Whale].
- **Lean reality:** a tiny list won't reach statistical significance on small effects — use the
  holdout for **direction and gut-check**, not precision, and pool across campaigns over time.

---

## Quick operational checklist (Register Carolina)

- [ ] Ad Grants: ≥ 2 ad groups/campaign, ≥ 2 ads/group, ≥ 2 sitelinks, phrase/exact match, negatives loaded.
- [ ] Ad Grants: conversion tracking live → **switch to Maximize Conversions** to clear the $2 CPC cap.
- [ ] Ad Grants: weekly Search-Terms + Quality-Score audit; keep account CTR ≥ 5%; complete annual survey.
- [ ] Email: welcome + new-donor + registrant + lapsed journeys built; dynamic segments; preference center.
- [ ] Email: SPF + DKIM + DMARC published; complaint rate < 0.1%; one-click unsubscribe.
- [ ] CRM: chosen on team size/budget; donation forms auto-write to records; CRM↔ESP↔GA4 connected.
- [ ] GA4: key events set **before** launch; documented lowercase UTM convention; cross-domain configured.
- [ ] Measurement: one honest proxy for the offsite registry conversion + a 10% holdout for the next push.

---

## Sources

Verify time-sensitive program rules against the official current docs before acting.

- **Google Ad Grants — Policy Compliance Guide** (official): `support.google.com/nonprofits/answer/9314402` — 5% CTR, structure minimums, keyword rules, conversion tracking, annual survey. (fetched 2026-06)
- **Google for Nonprofits — Account management policy** (official): `support.google.com/nonprofits/answer/117827`.
- **Getting Attention — "Google Ad Grants Rules: 9 Compliance Policies"**: structure/keyword/CTR detail. (updated 2026-02-18)
- **Media Cause — "Should Nonprofits Use the Maximize Conversions Bid Strategy?"**: $2 CPC cap exemption. (as of 2026)
- **AboveX Digital — Ad Grants bid strategies / suspension**; **Whole Whale — Ad Grants policy updates**; **Digital Tabby — real max CPC**: corroborate Maximize Conversions removing the $2 cap & Smart Bidding expectation. (as of 2026)
- **Nonprofit Megaphone; Big Sea — Ad Grants suspension/requirements**: deactivation causes & reactivation. (as of 2026)
- **M+R Benchmarks 2025** (`mrbenchmarks.com`): email frequency (~50/yr), revenue/subscriber (~$2.40), welcome-series CTR.
- **CauseVox — nonprofit welcome series**: ~80% welcome open rate, journey structure. (as of 2026)
- **PowerDMARC; Red Sift; Chronos — Gmail/Yahoo 2026 sender requirements**: SPF/DKIM/DMARC, p=none→reject, complaint-rate thresholds.
- **Zeffy; Constant Contact; softailed — nonprofit ESP & CRM comparisons** (2026): ESP/CRM pricing & fit.
- **Bloomerang; Neon One; LiveImpact; G2 — donor CRM comparisons** (2026): CRM pricing/positioning.
- **GA4 — Conversions vs key events** (official): `support.google.com/analytics/answer/13965727`; **Optimize Smart; Netpeak** — key-event setup, non-retroactivity. (as of 2026)
- **Blue Hills Digital; UTM.io; RallyUp — UTM conventions & GA4 for nonprofits** (2026): UTM taxonomy, cross-domain referral pitfall.
- **Kissmetrics — GA4 missing conversions / attribution limits** (2026): offline & multi-touch blind spots.
- **Improvado; Measured; Triple Whale; fusepoint — holdout & incrementality testing** (2026): holdout/geo-lift method, measurement windows.

## Disclaimer

Educational operations guidance, not legal, tax, or financial advice, and not an official statement of
any platform's terms. **Google Ad Grants, Gmail/Yahoo sender, GA4, and ESP/CRM rules and prices change
frequently** — every program rule, threshold, and price here carries an as-of date and **must be
re-verified against the official current documentation** before you rely on it. Items marked
**[UNVERIFIED]** were not confirmed against a primary source. For 501(c)(3) status, charitable-
solicitation registration, and CAN-SPAM/TCPA consent compliance, see the formation and marketing-
compliance siblings; for cause strategy, donor-journey design, and storytelling ethics, see
`venture-cause-nonprofit-marketing`.
