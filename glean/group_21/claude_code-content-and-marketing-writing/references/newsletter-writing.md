<!-- hub-reference-banner -->
> **Reference file — part of the `content-and-marketing-writing` hub.** Formerly the standalone `newsletter-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: newsletter-writing
description: Craft a recurring email newsletter — cadence rhythm, "from your editor" voice, intro paragraph that sets the week's lens, what-I'm-reading sections, PS footers, subject-line variation, sponsorship integration, segment vs broadcast tradeoffs, growth vs retention metrics, and Substack/Beehiiv/ConvertKit platform conventions. TRIGGER: "write my newsletter issue", "newsletter intro", "what should this week's lens be", "set up a weekly newsletter", "newsletter subject line variation", "newsletter cadence weekly vs monthly", "what I'm reading section", "PS footer for newsletter", "newsletter retention vs growth", "Substack/Beehiiv/ConvertKit voice", "Stratechery/Lenny/Money Stuff style", recurring email-publication craft. SKIP: one-off transactional or personal email (use email-craft); product changelog or release notes (use changelogs-for-humans); internal exec memo or all-hands (use executive-comms); conversion-focused sales/marketing email or landing copy (use sales-and-marketing-copy); customer case study (use case-study-writing); founder letter to investors or annual shareholder letter (use founder-letter-writing); general prose editing with no recurring-publication framing (use writing-expert).
---

# Newsletter Writing

## Overview

A newsletter is not "an email." It is a **recurring, named publication** with a voice, a cadence, and a contract with the reader. The contract: *I will show up at the same time, with the same lens, and you will get something worth your inbox.* Break the cadence or break the lens and the contract dies. Readers don't unsubscribe, they just stop opening.

This skill covers the craft of *recurring* email publications in the Substack / Beehiiv / ConvertKit lineage: the format Stratechery (Ben Thompson) productized, Lenny Rachitsky scaled to 1M+ subscribers, Matt Levine (Money Stuff) made daily, and Ryan Broderick (Garbage Day) made weekly. The unit of work is **the issue**, but the unit of strategy is **the year of issues**.

Use this when the user is writing or designing a publication that ships on a schedule, not a one-off email. For single transactional sends, use `email-craft`. For product update notes (even if they go out by email), use `changelogs-for-humans`. For internal exec memos, use `executive-comms`. For conversion-focused funnel emails, use `sales-and-marketing-copy`.

## Core Concepts

### 1. Cadence is a contract, not a preference

Pick a cadence you can hold **for years**, not months. Weekly Tuesday 9am is the dominant pattern (Lenny, Stratechery weekly free, Garbage Day) because (a) Tuesday/Wednesday land in the inbox slack between Monday catchup and Thursday wind-down, and (b) 9am local lands after the morning standup wave. Monthly first-of-month works for deeper, slower lenses (year-in-review formats, longform essays). Daily (Money Stuff) only works if the topic *churns daily* and you have a full-time commitment.

Heuristic: **the cadence sets the depth ceiling.** A weekly demands one strong idea. A monthly demands a thesis. A daily demands a beat. Don't promise weekly and deliver monthly content; readers feel the dilution.

### 2. The "from your editor" voice

The newsletter voice is **first-person, conversational, opinionated, and named.** It is not a brand voice; it is a person's voice. Ben Thompson writes as Ben. Lenny writes as Lenny. Matt Levine writes as Matt. Even at company newsletters, the best ones name an editor ("From Lenny's desk," "Notes from the editor"). This is the single biggest voice lever; abandon the corporate "we" and write as **I**.

Concrete markers of the voice: contractions, asides in parentheses, occasional self-deprecation, naming your own uncertainty ("I'm not sure this is right, but..."), and the occasional sentence fragment. The voice should sound like the writer is emailing one friend, not broadcasting.

### 3. The intro paragraph sets the week's lens

Every issue opens with **a lens**: a 2–4 sentence intro that tells the reader *how to read this issue*. It is not a summary. It is a framing. "This week I want to talk about why X is happening, because last week's Y made it impossible to ignore." The lens does three jobs: (1) signals you have a point of view, (2) gives the issue a center of gravity, (3) lets a skim-reader decide to invest 4 minutes.

Stratechery does this in the opening paragraph of every essay. Lenny does it in the deck-style summary at the top. Money Stuff does it with a one-liner section header that ironically frames the day's market drama.

### 4. Subject-line variation across a year

You will write ~52 subject lines a year on a weekly. They cannot all use the same formula or open rates collapse (the brain learns to skip the pattern). Rotate across at least four archetypes:

- **Curiosity gap**: "The pricing mistake every PM makes"
- **Number/list**: "7 onboarding flows I stole this week"
- **Named entity**: "What Stripe's new pricing tells us"
- **First-person declarative**: "I changed my mind about A/B testing"

Avoid clickbait. It spikes opens but kills retention, and **retention is the metric that matters** (see #9). Short subject lines (under ~50 chars) consistently outperform long ones in mobile preview panes.

### 5. The "what I'm reading" / curation section convention

A persistent secondary section ("What I'm reading this week," "Things I clicked," "Links I liked" in the Tim Ferriss / Lenny / Patrick Collison style) does two jobs: (a) gives readers value when the main essay isn't for them, (b) positions you as a curator-with-taste, which is a different and complementary authority to writer-with-thesis. 3–7 items is the sweet spot. Each item gets one line of *your* commentary, not a summary of the linked piece. Without commentary, it reads like a feed; with commentary, it reads like taste.

### 6. The PS footer convention

The PS at the end of a newsletter is a power slot. It is read at a disproportionately high rate because readers who reach the end are your most engaged segment. Use it for: a soft CTA ("hit reply and tell me what you'd add"), a low-pressure ask (job board, referral program), a personal note, or a teaser for next week. Do **not** stack three PSes. One PS, one job. The single-job PS is also where ConvertKit-style creators put the "reply to this email" CTA that drives the reply-rate metric (which mailbox providers use as a positive engagement signal).

### 7. Segment vs broadcast tradeoffs

Broadcasting (everyone gets the same issue) keeps voice and cadence consistent and is what readers actually signed up for. Segmenting (paid vs free, by topic interest, by tenure) sounds clever but fragments the editorial product and multiplies the writing load. Default to broadcast. Reserve segmentation for: (a) free vs paid tiers (Stratechery, Lenny premium), (b) onboarding sequence for new subscribers (welcome flow + greatest hits), (c) re-engagement campaigns for cold subscribers (90+ days no open). Avoid topic-interest segmentation unless your list is >50k and your topics genuinely don't overlap.

### 8. Sponsorship integration ethics

If you take sponsors, the ethical baseline is: (a) clearly labeled as sponsored ("Together with [sponsor]" or "A word from this week's sponsor"), (b) visually distinct from editorial, (c) never wrapped inside the editorial argument, (d) you would use or recommend the product yourself, (e) the sponsor does not get editorial review of the main piece. The Lenny / Morning Brew / Stratechery model puts the sponsor in a fixed slot (top or upper-third) so the reader's pattern-match is "this is the ad, the editorial follows." Hiding sponsorship inside editorial (the "native advertorial" pattern) is the fastest way to burn reader trust.

### 9. Growth vs retention metrics

**Retention is the real metric.** Beehiiv's 2026 state-of-newsletters argument: a list of 3,000 deeply engaged readers beats 30,000 ghosts every time. Track:

- **Open rate**: directional only; mailbox providers in 2025–2026 have degraded its signal value
- **Click rate**: better signal than opens
- **Reply rate**: strongest engagement signal mailbox providers use
- **Retention curve**: % of subscribers still opening at week 4, 12, 26, 52
- **Unsubscribe rate per issue**: a single spike frequently diagnoses one bad issue
- **Forward / share rate**: proxy for "would I send this to a friend"

Growth without retention is vanity. Hit "publish" on an issue you would forward yourself. That's the bar.

### 10. Platform lineage: Substack vs Beehiiv vs ConvertKit

- **Substack**: writer-first, social discovery via Notes and the recommendation network (Lenny credits this as one of his biggest growth drivers). Best for: solo writers building audience, paid subscriptions, embedded community.
- **Beehiiv**: operator-first, built by Morning Brew alumni. Strong on growth tooling (referral program, ad network, segmentation), SEO web archive. Best for: newsletters monetized via sponsorship or operated as media businesses.
- **ConvertKit (Kit)**: creator/marketer-first, strong automation and tag-based segmentation, weaker discovery. Best for: creators with a product (course, book, software) where the newsletter is a funnel, not the product.

Pick based on **what the newsletter is for**: audience → Substack; media business → Beehiiv; product funnel → Kit.

## Templates and Examples

### Newsletter issue skeleton (weekly)

```
Subject: [varies across 4 archetypes; see Core Concept #4]
Preview text: [40–90 chars that complement, don't repeat, the subject]

[Sponsor block: fixed slot, clearly labeled, ~50 words, optional]

Hey everyone,

[Intro / lens paragraph: 2–4 sentences. What is this week about?
 Why now? What's the one thing the reader should hold in their head
 as they read?]

## [Main section heading: declarative, not cute]

[Main body: one strong idea, 400–1,200 words depending on cadence
 depth. Subheads every 200–300 words. First-person voice. Examples.
 At least one number, name, or specific. End with the "so what."]

## What I'm reading this week

- [Link]. [One sentence of YOUR commentary, not a summary.]
- [Link]. [Your take.]
- [Link]. [Your take.]

## [Optional secondary section: tool of the week, reader question,
   chart of the week, etc. Consistent slot, not every issue needs it]

— [First name]

PS: [One thing. A soft CTA, a personal note, or a teaser. Not three.]
```

### Sample intro paragraphs (the "lens")

**Curiosity-led:**
> "Three different founders asked me the same question this week, which usually means I should write about it. The question: when do you stop being the salesperson? Here's where I landed."

**News-pegged:**
> "Stripe shipped a pricing change on Tuesday that almost nobody covered, and I think it's the most interesting thing that happened in payments this year. Let me explain why."

**Counter-thesis:**
> "Everyone is writing about why agents are eating SaaS. I want to write the opposite case, not because I think it's right, but because I haven't seen it argued, and arguments you can't find are usually the ones worth making."

### Sample PS footers

- "PS: hit reply and tell me what I got wrong. I read every reply."
- "PS: we're hiring two PMs. Forward this to someone good."
- "PS: next week, the onboarding teardown I've been promising for a month. Finally."

## Anti-Patterns

- **The everything-bagel issue.** Five disconnected sections, no lens, no thesis. Cut to one.
- **The corporate "we."** Kills the voice. Use I.
- **The cadence drift.** "Weekly" becomes biweekly becomes monthly becomes guilt-stricken silence. Either re-commit to weekly or publicly downshift to monthly and reset expectations.
- **The summary-only "what I'm reading."** Links without your take are a feed reader, not a newsletter.
- **The three-PS pileup.** Each PS competes with the others. Pick one.
- **Clickbait subjects.** Spike opens, kill retention, kill long-term inbox placement.
- **Hidden sponsorship.** Sponsor wrapped inside editorial. Trust-burning.
- **Over-segmentation.** Splitting your list into eight interest cohorts when you have 4,000 subscribers. You're shipping eight half-newsletters.
- **The intro that summarizes the issue.** The intro should *frame*, not preview. Save the body for the body.
- **Tracking opens as the north-star metric.** Use retention curve and reply rate. Opens are degraded signal in 2026.

## Decision Heuristics

- **Weekly vs monthly?** If the topic moves weekly and you can sustain it for 2+ years, weekly. Otherwise monthly. Don't try biweekly; readers can't internalize a biweekly rhythm.
- **Subject line working?** Read it aloud. If it sounds like a marketer wrote it, rewrite. If it sounds like you texted a friend the headline, ship.
- **Cut or keep this section?** If you removed it for three issues, would readers email you asking where it went? If no, cut.
- **Sponsor fit?** Would you use this product? Would you mention it to a friend without being paid? If no on either, decline.
- **Free vs paid tier?** Free tier should be self-sufficiently great. Paid tier is for people who already love the free tier and want more depth or community, not for "the real content."
- **What goes in the PS?** The one thing the reader should do, remember, or click after they close the email. If you can't name it in five words, the PS isn't ready.
- **Open rate dropped 8% this week, panic?** No. Check reply rate and retention. If those are stable, it's mailbox-provider noise. If reply rate also dropped, look at the subject line and the lens paragraph.

## References

- [Lenny's Newsletter — Lenny Rachitsky on Substack](https://www.lennysnewsletter.com/): canonical playbook for the weekly operator newsletter; growth via Substack's recommendation network; the "what I'm reading" section convention.
- [Stratechery by Ben Thompson](https://stratechery.com/about/): the subscription-newsletter business model Substack's founders cite as their inspiration; the strategic-framework lens; the long-form essay format.
- [The State of Newsletters 2026 — Beehiiv Blog](https://www.beehiiv.com/blog/the-state-of-newsletters-2026): retention vs growth metrics, the degradation of open rate as signal, the case for treating newsletters as media companies.
- [How To Measure Real Email Engagement in 2026 — Beehiiv Blog](https://www.beehiiv.com/blog/email-engagement-metrics): the engagement-rate stack (click + reply + forward + retention) over opens-only.
- [Lenny Rachitsky on building a consistent writing habit — Substack On Substack](https://on.substack.com/p/how-to-create-consistent-writing-habit-lenny): cadence as a years-long commitment, quality × consistency.
- Matt Levine, *Money Stuff* (Bloomberg): the daily-beat newsletter as a model; section-header voice; the "one ironic frame per topic" pattern.
- Ryan Broderick, *Garbage Day*: the weekly-culture newsletter; first-person voice; lens-driven intros.
