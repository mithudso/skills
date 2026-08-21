<!-- hub-reference-banner -->
> **Reference file — part of the `content-and-marketing-writing` hub.** Formerly the standalone `release-blog-and-launch-narrative` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: release-blog-and-launch-narrative
description: "Launch-blog and product-announcement narrative craft — the problem-to-solution-to-demo-to-CTA arc, the 'why now' paragraph, the 'what's not in v1' honest-disclosure section, embargo and sequencing, the cross-channel launch ladder (blog → tweet → email → community → analyst), launch-day comms tree, pre-launch teasers, the 'release notes accompany the announcement' convention, success criteria, and the relaunch problem. Distinct from conversion-focused sales copy: launches have a narrative arc and disclosure ethics. TRIGGER: 'write a launch blog post', 'product announcement', 'launch narrative', 'launch sequencing', 'embargo press release', 'launch day plan', 'why now paragraph', 'pre-launch teaser', 'announcement post', 'launch comms tree', 'GA announcement', 'beta to GA narrative', 'how do we announce this'. SKIP: a landing page or sales email aimed at converting an unaware audience (use sales-and-marketing-copy); an internal exec memo about a launch (use executive-comms); the engineering changelog or semver release notes (use changelog-and-release-notes); a user-facing changelog entry with screenshots (use changelogs-for-humans); the prose-craft pass on the draft (use writing-expert)."
version: "1.0.0"
updated: "2026-05-29"
category: custom
tags:
  - writing
  - launch-comms
  - narrative
  - product-marketing
parent_concept: "Writing and Documentation"
whenToUse:
  - "Write a launch blog post for a new product or feature"
  - "Structure an announcement post for GA"
  - "Help me draft the 'why now' paragraph"
  - "What goes in the 'what's not in v1' section"
  - "Plan the launch sequencing — blog, email, tweet, analyst"
  - "Draft an embargo for the press list"
  - "Build a launch-day comms tree"
  - "Write a pre-launch teaser sequence"
  - "How do we announce something we already shipped quietly six months ago"
  - "Coordinate release notes with the announcement"
related_skills:
  - sales-and-marketing-copy
  - executive-comms
  - changelog-and-release-notes
  - changelogs-for-humans
  - writing-expert
  - storytelling-and-narrative
  - rhetorical-frameworks-deep
triggers:
  - write a launch blog post
  - product launch announcement
  - launch narrative
  - launch sequencing
  - embargo press release
  - launch day plan
  - why now paragraph
  - pre-launch teaser
  - announcement post
  - launch comms tree
  - GA announcement
  - beta to GA narrative
  - how do we announce this
  - relaunch the same thing
  - launch ladder
---

# Release Blog and Launch Narrative

Reference for **launch-mode writing**: the blog post, press release, sequencing plan, and disclosure ethics that surround a public product announcement. A launch is a one-shot narrative event — distinct from a landing page (steady-state conversion), an exec memo (internal alignment), or a changelog (continuous release log).

Sources: Apple Newsroom press releases (apple.com/newsroom), Stripe Sessions and Stripe Blog product launches, Linear's blog launch posts, Vercel/Next.js Conf launch posts, Andy Raskin's "The Greatest Sales Deck I've Ever Seen" (Medium, Mission.org, 2016) and his strategic narrative framework, Geoffrey Moore *Crossing the Chasm* (3rd ed., HarperCollins, 2014), Lenny Rachitsky's launch writing on lennysnewsletter.com, Aakash Gupta product-launch playbooks (aakashg.com), Atlassian and Asana launch templates, First Round Review "37 minutes in an Amazon war room" (review.firstround.com), Apple/PR.co "Masterclass in Owned Media" analysis.

SKIP: a conversion-focused landing page (use `sales-and-marketing-copy`); an executive memo (use `executive-comms`); the engineering changelog or semver release notes (use `changelog-and-release-notes`); the user-facing changelog with screenshots (use `changelogs-for-humans`, the sibling skill); the prose-craft polish pass on a draft (use `writing-expert`).

---

## How to use this skill

When invoked, follow this sequence before writing:

1. **Classify the launch.** Is this a category-defining launch (new product, new category), a tier launch (GA after beta, "now generally available"), a feature drop inside an existing product, or a relaunch (re-announcing something already shipped)? Each one needs a different arc — see Section 1.
2. **Identify the "why now."** A launch post is not just an announcement; it is an argument for why this thing exists *and why it is being announced today*. If you cannot articulate "why now" in one sentence, do not write the post yet.
3. **Pick the audience tier.** Developer launch (technical depth, code samples), end-user launch (benefits, screenshots, video), executive/analyst launch (market positioning, scale, customer logos), or hybrid (most are hybrid, but one audience must be primary).
4. **Plan the sequence before drafting the post.** The blog post is one node in a launch ladder (Section 4). Drafting it in isolation produces a post that contradicts the tweet, the email, and the analyst briefing. Sketch the ladder first.
5. **Draft the post against the arc** (Section 2), then layer in the "why now" paragraph (Section 3), the honesty section (Section 5), and the call to action.
6. **Verify against the anti-patterns** (Section 8) before shipping.

If the caller has not stated the launch type, the audience, or the "why now," ask exactly one targeted question. Do not invent a "why now" — that is the load-bearing claim of the entire post.

---

## 1. Four launch archetypes

| Archetype | Example | Narrative arc | Risk |
|---|---|---|---|
| **Category-defining** | iPhone 2007, Stripe in 2011, Linear in 2019 | "The old game is over. A new game is starting. Here is the new game and our piece of it." | Overclaim; chasm failure; the new game isn't real |
| **GA / tier launch** | "Generally available," "Out of beta," "v1.0" | "We've been quietly proving this works. Now you can rely on it. Here's the proof." | Audience asks "wasn't this already out?"; weak proof |
| **Feature drop** | A new module, integration, model | "You already use [product]. Here is the thing you've been asking for." | Lost in the noise; not announcement-worthy |
| **Relaunch** | Rebrand, repositioning, re-emphasizing a quiet ship | "We shipped this. You may have missed it. Here is why it matters now." | Embarrassment; "we already did this" press cycle |

The archetype dictates everything downstream — arc, sequencing, embargo, what goes in v1's honesty section, and whether release notes accompany.

**Heuristic.** If you can't name the archetype in 10 seconds, the launch isn't ready. Stop and figure out which one this is before drafting.

---

## 2. The launch-post arc

The canonical structure for a launch blog post is a five-beat arc. Andy Raskin's strategic-narrative work, Apple Newsroom releases, and Stripe's launch posts all use variations of this; the labels differ but the shape is constant.

### Beat 1 — Name the shift (the "why now")

Open with a one-paragraph claim about a shift in the world. **Never open with the product.** Raskin's whole thesis: "the old game is over; a new game is starting." This is the paragraph that earns the reader's attention before the product is introduced.

| Strong | Weak |
|---|---|
| "Two years ago, paying online meant integrating with a bank. Today, it means writing seven lines of code." | "We're excited to announce a new feature." |
| "Customer support used to be a ticket queue. Now it's a real-time conversation, and ticket queues are slowing teams down." | "Today we're launching v2.0 of our help desk." |

If the "why now" is a recycled industry trend with no claim attached — "AI is everywhere" — it is not a "why now." It is a backdrop. Cut it.

### Beat 2 — Stakes (the cost of the old game)

What does the reader lose by staying with the old way? Quantify if possible. Apple's iPhone launch named the cost as "three separate devices in your pocket." Stripe's payments posts name the cost as "weeks of engineering time before you accept a dollar."

This is not a sales pitch about pain. It is the cost the reader is *already paying* and may not have named.

### Beat 3 — The reveal (here is what we built)

Now — and only now — name the product and the one-line description. Apple's pattern: "Today, Apple introduced [X], a [category] that [primary benefit]." This sentence should be the most-quoted sentence of the post.

Follow with three to five capability paragraphs, each starting with the user outcome, not the underlying mechanism. Code samples, screenshots, video go here, not earlier.

### Beat 4 — Proof

Customer logos, quotes, measurable outcomes, demo link, screenshots of the product in use. One outside voice is worth ten internal claims. Quantified claims ("3x faster," "deploys in 90 seconds") beat unquantified claims ("blazing fast") every time. If you cannot quantify, name a customer who can.

### Beat 5 — Call to action and honesty

Two things, in this order:

1. **Clear CTA** — sign up, talk to sales, read the docs, try the demo. One primary CTA, one secondary at most. If you list four CTAs the reader picks zero.
2. **What's not in v1** (Section 5). The honest-disclosure section. Skipping this is the most common launch mistake and the one that destroys credibility fastest.

---

## 3. The "why now" paragraph

The "why now" is a load-bearing claim that does three things in one paragraph:

1. **Names a shift** that has already occurred or is occurring (not a prediction).
2. **Establishes the cost** of operating as if the shift hadn't happened.
3. **Creates urgency** that is *not* "limited time offer" urgency. It is "you're already behind if you ignore this" urgency.

**Template.**

> Until recently, [old default] was the only way to [task]. But [shift A] and [shift B] have changed what's possible. Teams that adapt are [outcome]. Teams that don't are [cost]. Today we're announcing [product] to help you [verb] in this new world.

**Worked example (developer-tool launch).**

> Until recently, deploying a web app meant owning servers, configuring autoscalers, and writing infrastructure code that had nothing to do with your product. But the rise of edge networks and the maturation of serverless primitives have made server management an implementation detail. Teams that lean into this ship five times more often. Teams that don't are still in standup meetings about Kubernetes upgrades. Today we're launching [Product] to give your team the deploy experience the new world has been waiting for.

**Anti-pattern: the "AI is everywhere" why-now.** "AI is transforming everything" is not a why-now — it has no specific shift, no cost, no urgency the reader can act on. If your why-now would work verbatim for ten unrelated products, rewrite it.

**Heuristic.** A working "why now" paragraph should fail the substitution test: try replacing your product name with a competitor's. If the paragraph still works, your why-now is generic. Add specifics until the swap breaks the paragraph.

---

## 4. The cross-channel launch ladder

A launch is not a blog post. A launch is a sequenced cascade where each channel quotes, links to, or amplifies the previous one. The canonical ladder, from highest fidelity to lowest:

```text
            ┌──────────────────────────┐
T-2 weeks   │ Analyst & press briefing │  (embargoed)
            │ Customer reference calls │
            └────────────┬─────────────┘
                         │
T-3 days    ┌────────────▼─────────────┐
            │ Pre-launch teaser        │  (cryptic)
            │ Waitlist priority email  │
            └────────────┬─────────────┘
                         │
T=0 hour    ┌────────────▼─────────────┐
            │ Blog post (canonical)    │  ← single source of truth
            │ Press release (newsroom) │
            │ Release notes published  │
            └────────────┬─────────────┘
                         │
T+0:05      ┌────────────▼─────────────┐
            │ Founder/CEO tweet thread │  (links to post)
            │ Company X/LinkedIn post  │
            └────────────┬─────────────┘
                         │
T+1 hour    ┌────────────▼─────────────┐
            │ Customer email blast     │  (benefit-led)
            │ In-product announcement  │
            └────────────┬─────────────┘
                         │
T+1 day     ┌────────────▼─────────────┐
            │ Community post (Hacker   │
            │  News, Reddit, Discord)  │  ← founder shows up
            │ Customer-facing webinar  │
            └────────────┬─────────────┘
                         │
T+2 to +7d  ┌────────────▼─────────────┐
            │ Partner co-marketing     │
            │ Long-form deep-dive blog │
            │ Podcast/conference talk  │
            └──────────────────────────┘
```

**Sequencing rules.**

1. **The blog post is the canonical source.** Every other piece links to it. The tweet quotes a line from it. The email summarizes it. The press release is a structured version of it. If any two channels contradict each other, fix the blog post first.
2. **Analyst and press briefings are under embargo** — they get the post early in exchange for a specific publish time. Embargoes only work if every recipient explicitly agrees in writing; an unsolicited "EMBARGOED UNTIL" header is not binding.
3. **Customer references go before press, not after.** The press briefing should already include "and here are three customers who'll talk to you" — chasing references after the announcement is amateur hour.
4. **Hacker News / community posts come from a real human**, ideally the founder or PM, posted from a personal account. A corporate account submitting its own launch reads as astroturf.
5. **Release notes accompany the announcement.** When the blog post says "today we're launching X," the docs and the changelog must already reflect X. A reader who follows the "read the docs" CTA to a 404 will not return.

---

## 5. The "what's not in v1" honesty section

Every launch has scope cuts. Audiences know this. The post that pretends otherwise loses credibility on first use; the post that names the cuts gains credibility before first use.

**Where it goes.** Near the end of the post, after the CTA, under a header like "What's not in v1," "What's next," or "Known limitations." Not buried. Not in a separate FAQ. In the post.

**Template.**

> v1 deliberately ships without [X], [Y], and [Z]. We chose this scope to get [primary benefit] in your hands sooner. [X] is on the roadmap for [quarter]. If [X] or [Y] is a hard requirement for you today, [alternative path — wait, contact us, use this workaround].

**What this section does for the reader.**

- Lets engineers in evaluation mode answer "can I use this for my case yet?" without booking a call.
- Removes "they overclaimed" as a future complaint.
- Pre-empts the inevitable Hacker News commenter who finds the missing feature.

**What this section does for the writer.**

- Forces a real scope conversation with the product team before the post ships, surfacing claims that are not yet true.
- Distinguishes the post from competitor announcements that pretend their v1 does everything.

**Anti-pattern.** A "known limitations" section that lists only minor things ("doesn't support Bulgarian localization yet") while omitting the big missing feature everyone will ask about. Readers notice. Name the elephant.

---

## 6. Embargo and sequencing mechanics

An **embargo** is an agreement with press or analysts: you give them the news early, they hold publication until your stated time. Embargoes are professional trust, not legal contracts. Rules that actually work:

1. **Get explicit written consent before sharing.** "Are you willing to take this under embargo until [date/time TZ]?" Wait for "yes." Only then send the briefing materials.
2. **State the embargo at the top of every doc.** Bold, red, top of page: `EMBARGOED UNTIL [DATE] [TIME] [TIMEZONE]`. Repeat in the email subject and body.
3. **Give a useful lead time.** 48–72 hours is standard for press, 1–2 weeks for analyst firms (Gartner, Forrester) who need to write longer pieces.
4. **Provide a press kit.** Logo, product screenshots, executive headshots, a one-page fact sheet, a quote from the CEO, customer quote(s) on the record. Not optional.
5. **Plan for embargo breaks.** Someone will publish early. Have a Plan B: pre-stage the blog post in draft, be ready to publish 2–6 hours early if a break is widely reported.
6. **Stripe Sessions / Vercel Ship / Apple keynote pattern.** When a launch event is the embargo lift, all materials drop at the moment the keynote announces them, including blog posts, docs updates, pricing pages, and changelog entries.

---

## 7. Launch-day comms tree

A launch-day comms tree is the human-routing diagram for the four hours around T=0. It answers: who watches what, who decides what, who responds where.

```text
                        ┌──────────────────────┐
                        │   Launch Commander   │  (one person, full authority)
                        │   - Go/no-go         │
                        │   - Embargo break    │
                        └──────────┬───────────┘
                                   │
        ┌──────────┬───────────┬───┴───────┬────────────┬──────────────┐
        ▼          ▼           ▼           ▼            ▼              ▼
   ┌────────┐ ┌────────┐ ┌──────────┐ ┌─────────┐ ┌──────────┐ ┌────────────┐
   │ Eng    │ │ Press  │ │ Social   │ │ Support │ │ Customer │ │ Hacker     │
   │ on-call│ │ relay  │ │ ops      │ │ triage  │ │ success  │ │ News watch │
   │        │ │        │ │          │ │         │ │          │ │            │
   │ - Site │ │ - Mon. │ │ - Sched. │ │ - New   │ │ - Top    │ │ - Founder  │
   │   up   │ │   pubs │ │   posts  │ │   tickt │ │   accts  │ │   replies  │
   │ - 404s │ │ - Quote│ │ - DMs    │ │ - Bug   │ │   FYI    │ │ - Mod      │
   │ - 5xx  │ │   sched│ │ - Reply  │ │   surge │ │ - QBR    │ │   issues   │
   │        │ │        │ │   queue  │ │   spike │ │   slot   │ │            │
   └────────┘ └────────┘ └──────────┘ └─────────┘ └──────────┘ └────────────┘
```

**Rules for the tree.**

- **One commander, full authority.** Not a committee. The commander can pull the launch, push embargo, escalate. Everyone else routes through them.
- **Pre-approved message templates** for the top five failure modes: site outage, embargo break, major bug discovered post-launch, viral negative thread, press misquote. Drafted *before* T=0, signed off, ready to fire.
- **Hot list** of issues being actively watched, posted on a shared screen or pinned channel. New ones appended in real time.
- **End-of-day debrief** scheduled before launch day starts. What worked, what broke, what we'll change for the next launch.
- **PPR balance** (promotion, pressure, response) — don't spend all your bandwidth on proactive promotion if a reactive issue is brewing.

---

## 8. Anti-patterns

**The "we're excited to announce" opener.** Lazy, generic, replaceable in any post. Cut. Open with the shift.

**The "this changes everything" overclaim.** It doesn't. Even if it did, claiming so directly trains the reader to discount you. Show, don't tell — let the reader conclude this changes everything from the proof section.

**The "product-first" structure.** Specs, then features, then maybe a sentence about who cares. Inverts the arc. Move the "why now" and the stakes to the top.

**The hidden "what's not in v1."** Burying scope cuts in an FAQ or a docs page tagged "limitations" is the same as omitting them, and readers treat it as deception when they hit the wall.

**The "AI/cloud/edge is everywhere" why-now.** A "why now" that works for any product is a "why now" for no product. Specify.

**The post that doesn't match the docs.** Blog post says "supports X." Docs are silent on X, or docs say "X coming soon." Ship-blocker. Fix the docs or the post before launch.

**The relaunch problem.** Re-announcing the same thing six months later as if new. Readers and journalists remember. If you must re-announce (e.g., "now generally available," "now with a major new capability"), be explicit: lead with what's new *this time*, link to the original announcement, never imply this is a first-time launch.

**The unprovoked superlative.** "The world's best," "the most powerful," "industry-leading" without a citation. Cut or attribute (Gartner, customer quote, benchmark with methodology).

**The four-CTA ending.** Sign up, talk to sales, read the docs, join the webinar, follow us on X, watch the demo, download the whitepaper. The reader picks zero. One primary CTA, one secondary, stop.

**The blog post without an author byline.** Launches benefit from human attribution — CEO, PM, founding engineer. "By the Acme Team" is a signal you're hiding behind the brand.

**The launch without success criteria.** If you don't know what success looks like (signups in week 1, press hits in tier-1 outlets, customer references generated), you can't tell if the launch worked, which means the next launch repeats the same mistakes.

**The Hacker News no-show.** Posting the link and disappearing while comments pile up. Founders/PMs who answer comments in the first 4 hours get a measurably better thread.

---

## 9. Success criteria for a launch

Define before launch day. Track for 30 days. Sample tiers:

| Tier | Metric | Source |
|---|---|---|
| **Reach** | Unique pageviews on the blog post in week 1 | Web analytics |
| **Reach** | Tier-1 press hits (named outlets agreed in advance) | PR log |
| **Engagement** | Hacker News front page, time on front page, top comment sentiment | HN / manual review |
| **Conversion** | Signups attributable to launch (UTM, referrer) | Analytics + product DB |
| **Conversion** | Sales-qualified leads from launch sources | CRM |
| **Adoption** | New-product activation rate in week 1, week 4 | Product analytics |
| **Sentiment** | Net inbound mentions (social, community) — count and sentiment | Listening tool / manual |
| **Internal** | Did the comms tree hold? Embargo respected? Site uptime? | Postmortem |

The number of metrics is less important than the **pre-commit**: write the success criteria down before launch, and review them at the 30-day postmortem.

---

## 10. Decision heuristics

- **Can't name the archetype** (category, GA, feature, relaunch) → not ready to write.
- **Can't write the "why now" without it being generic** → not ready to write.
- **The docs don't yet match the post** → not ready to ship.
- **No release notes / changelog entry ready** → not ready to ship (and use `changelog-and-release-notes` for the developer changelog, `changelogs-for-humans` for the user-facing one).
- **No customer reference / quote / logo** → consider holding for a beta-to-GA arc when you have one.
- **The "what's not in v1" section is empty** → you have not had a real scope conversation yet; have it now.
- **Press list has not confirmed embargo in writing** → no embargo; treat as public.
- **No launch commander named** → no launch day; name one.
- **It's a relaunch** → say so explicitly in the post. Don't pretend.
- **The CTA list has more than two items** → cut to one primary, one secondary.
- **Author byline is "the Team"** → find a human to put their name on it.

---

## 11. Templates

### 11.1 Launch blog post skeleton

```markdown
# [Product name]: [one-line description that names the new game]

By [Human Author], [Title] · [DATE]

[BEAT 1 — Why now.] Until recently, [old default]. But [shift]. [Cost of old default.]

[BEAT 2 — Stakes.] Teams that operate as if [shift] hasn't happened are paying for it in [specific cost]. Over the past [period], we've watched this play out at [N] customers.

[BEAT 3 — Reveal.] Today, we're launching [Product], a [category] that [primary benefit] in [time/effort unit].

## [Capability 1, framed as outcome]
[Two to four sentences, ending with a screenshot, code sample, or short demo.]

## [Capability 2, framed as outcome]
[Same.]

## [Capability 3, framed as outcome]
[Same.]

## How customers are using it

> "[Specific, quantified quote.]"
> — [Name, Title, Company]

[Optional: two more named-customer quotes.]

## Get started

[Primary CTA, one verb: Start building / Try it / Read the docs / Talk to us.]
[Secondary CTA, optional: View the demo / Read the deep-dive.]

## What's not in v1

[Honest paragraph about scope cuts. Roadmap commitment for top one or two.]

---

Related: [release notes link] · [docs link] · [pricing link]
```

### 11.2 Embargoed press-release skeleton

```text
EMBARGOED UNTIL [DAY], [DATE] AT [TIME] [TIMEZONE]
DO NOT PUBLISH BEFORE THIS DATE/TIME.

FOR IMMEDIATE RELEASE [post-embargo]

[CITY] — [Company] today announced [Product], a [category] that
[one-sentence primary benefit].

[Why now: shift + stakes, condensed to one paragraph.]

[Reveal: product description with two to three capability paragraphs.]

"[CEO or product lead quote that says something specific, not 'we're excited']."
— [Name, Title, Company]

"[Customer quote with a measurable outcome.]"
— [Name, Title, Customer company]

[Availability: when, where, pricing, geographic limits.]

[Boilerplate "About [Company]" paragraph.]

Press contact: [Name, email, phone]
Press kit: [URL — logos, screenshots, headshots, fact sheet]
```

### 11.3 Launch-ladder timing template

```text
T-14 days   Analyst briefings scheduled (Gartner, Forrester, IDC)
T-10 days   Press list confirmed; embargo agreements in writing
T-7  days   Customer-reference calls; quotes finalized; press kit assembled
T-5  days   Pre-launch teaser content drafted; founder thread written
T-3  days   Press receives embargoed briefing + materials
T-2  days   Internal all-hands; comms-tree dry run; pre-approved templates signed off
T-1  day    Final post review; docs/changelog/pricing pages staged
T-12 hr     Site monitoring on; on-call engineer briefed
T-1  hr     Launch commander confirms go/no-go
T=0         Blog post publishes; press releases pubd; release notes live; tweet thread fires
T+5  min    CEO/founder X post + LinkedIn post
T+1  hr     Customer email blast; in-product announcement
T+4  hr     Founder appears on Hacker News thread, answers comments
T+1  day    Community AMA / Discord / Reddit
T+7  days   Long-form deep-dive follow-up post; postmortem of launch ops
T+30 days   Success-criteria review; what we'll change next launch
```

### 11.4 "What's not in v1" template

```markdown
## What's not in v1

We deliberately scoped v1 to [primary use case] so we could ship it in [time]
instead of [longer time]. That means a few things are not in this release:

- **[Capability X]** is on the roadmap for [quarter]. If you need [X] today,
  [workaround / alternative product / contact us].
- **[Capability Y]** is supported only for [subset]; full support arrives in
  [quarter].
- **[Integration Z]** is not yet built. We'll prioritize it based on demand —
  [feedback link / request URL].

If [X], [Y], or [Z] is a blocker for your use case, we'd rather you know now
than discover it during evaluation. [Contact link.]
```

### 11.5 Founder tweet-thread skeleton

```text
1/ [The why-now in a single tweet. No emoji "🚀". The shift, named clearly.]

2/ [Stakes — the cost of the old default. One concrete example.]

3/ Today we're launching [Product]. [One-line description.]

4–7/ [One tweet per capability. Each ends with a screenshot, GIF, or short video.]

8/ [Customer quote with name and company.]

9/ [What's not in v1 — one sentence acknowledging scope.]

10/ [CTA + link to the canonical blog post.]
```

---

## 12. References

- Andy Raskin, "The Greatest Sales Deck I've Ever Seen" (Medium / Mission.org, 2016) — five-element strategic narrative framework.
- Andy Raskin, "The Making of a Great Strategic Narrative" (LinkedIn, andyraskin.com).
- Geoffrey A. Moore, *Crossing the Chasm* (3rd ed., HarperCollins, 2014) — positioning template ("For/that need/is a/that — Unlike/provides"), whole-product concept, target-segment dominance.
- Apple Newsroom (apple.com/newsroom) — press-release format, "Today, Apple introduces..." sentence pattern, owned-media discipline.
- Stripe Blog launches (stripe.com/blog, stripe.com/newsroom) — sequencing around Stripe Sessions, capability-as-outcome capability headers.
- Linear blog and changelog (linear.app) — launch-post arc with embedded changelog cadence.
- Vercel changelog and launch posts (vercel.com/changelog) — daily-cadence model and Next.js Conf launch sequencing.
- Lenny Rachitsky, lennysnewsletter.com — operator-grade launch playbooks, podcast interviews on launches.
- Aakash Gupta, "9 Product Launch Strategies" (aakashg.com).
- First Round Review, "My Launch Lessons from 37 Minutes in an Amazon War Room" (review.firstround.com) — war-room and comms-tree mechanics.
- PR.co, "A Masterclass in Owned Media: How Apple's Newsroom Dominates PR" (pr.co/blog).
- Asana, *Product Marketing Launch Template* (asana.com/templates/product-marketing-launch) — timeline mechanics.
- Atlassian, *Product Launch Timeline* (atlassian.com/agile/product-management/product-launch-timeline) — stage definitions.
- PRLab, "What Is an Embargoed Press Release?" (prlab.co) — embargo mechanics.
- Taggbox, "A Guide to Social Media War Room" (taggbox.com/blog) — launch-day monitoring.
- WHOOP Engineering, "From War Room to War Floor" (engineering.prod.whoop.com) — engineering-side launch coordination.

---

## 13. Cross-skill boundaries

- **`sales-and-marketing-copy`** — conversion copy for steady-state pages (landing pages, ads, nurture emails, CTA microcopy). The launch post is a one-shot narrative event with disclosure ethics; use this skill, not that one, for the announcement itself.
- **`executive-comms`** — internal exec memos, board updates, all-hands talking points around the launch. The blog post is external; the exec memo is internal. Different audience, different artifact.
- **`changelog-and-release-notes`** — the developer-facing changelog and semver release notes that *accompany* the launch. Always pair: launch post + changelog entry, published the same hour.
- **`changelogs-for-humans`** — the user-facing changelog with screenshots and benefit-led entries, the sibling to the developer changelog. The launch post calls back to this; the changelog calls forward to the launch post.
- **`writing-expert`** — the prose-craft pass on the draft (anti-AI-isms, sentence flow, nominalizations). Run after this skill produces the structural draft.
- **`storytelling-and-narrative`** and **`rhetorical-frameworks-deep`** — deeper theory of narrative arcs, persuasion frameworks. Use when the launch needs a deeper strategic-narrative rebuild.
