<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `headline-craft` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: headline-craft
description: Craft, critique, and rewrite headlines for web articles, blog posts, marketing pages, news stories, op-eds, X/Twitter posts, and email subject lines. Covers the AIDA-headline relationship, the "you-frame" rule, headline-number conventions (odd numbers, BuzzSumo's top performers), curiosity-gap vs declarative headlines, news vs feature vs op-ed vs SEO-headline distinctions, the Upworthy/BuzzFeed lineage and where clickbait went wrong, A/B-test conventions, headline-length sweet spots by channel (web ~55–65 chars, X ~70, email subject ~30–50), and the subhead/deck pairing. TRIGGER when the user wants to write or rewrite a headline, title, subject line, page title, or article hed; when they ask "is this headline good," "make this more clickable," "what's a better title," or want A/B variants; when they're writing for a specific channel and need to know the character budget; when they want to evaluate clickbait risk. SKIP if the user only needs the one-sentence "headline test" rule (use sales-and-marketing-copy); if they want a full opinion piece (use op-ed-writing); if they want a launch announcement (use release-blog-and-launch-narrative); if they want a profile piece (use profile-writing); if it's a general writing critique with no headline focus (use writing-expert).
category: custom
tags: [writing, journalism, headlines, copywriting, seo, marketing]
---

# Headline craft

## Overview

Five times as many people read the headline as read the body copy — Ogilvy's most-quoted statistic, and the reason a headline is not a label that names the article but a **promise** that earns the click. A headline does three jobs at once: it tells the reader what they're about to get, gives them a reason to care, and signals what kind of piece it is (news, feature, opinion, marketing). When any of the three jobs fails, the click is either lost or — worse — bought with a promise the body can't keep.

This skill is for **writing, rewriting, and critiquing individual headlines** across the channels operators actually publish to: web articles, blog posts, marketing landing pages, news stories, op-eds, X/Twitter posts, LinkedIn posts, and email subject lines. The single-sentence "headline test" rule (would I click this if a stranger sent it to me?) is owned by `sales-and-marketing-copy`. This skill picks up where that test stops: when the answer is "no," what do you change?

## Core concepts

### 1. The five-times rule and the 80-cent budget

Ogilvy's empirical claim from print advertising — that five times more people read the headline than the body — has held up across digital channels with caveats. On the web, the multiple is closer to 8x because of social-feed scrolling and headline-only consumption. The practical implication is unchanged: **the headline is where 80% of the work happens**. A writer who finishes the body and dashes off a headline in 30 seconds has spent their effort backwards. Joanna Wiebe's discipline at Copyhackers: never write fewer than 10 headline variants before picking one. Ogilvy demanded 16. For high-stakes pieces (launch posts, op-eds, paid ads), 16 is still right.

### 2. AIDA mapped onto a headline

AIDA — Attention, Interest, Desire, Action — is usually framed as a full-page model, but it collapses cleanly onto the headline-and-deck pairing:

| AIDA stage | Where it lives | What it does |
|---|---|---|
| **Attention** | First 3 words of the headline | Stops the scroll. Concrete noun, surprising number, named person, or strong verb. |
| **Interest** | Rest of the headline | Promises a benefit, names a tension, or opens a curiosity gap. |
| **Desire** | Subhead / deck | Adds the second beat — specificity, proof, or who-it's-for. |
| **Action** | First line of body / CTA | Closes the loop the headline opened. |

If a headline can carry Attention and Interest alone, you don't need a deck. Most can't.

### 3. The "you-frame" rule

A headline pointed at the reader outperforms a headline pointed at the writer, the product, or an abstract topic. Compare:

- **Topic frame (weak):** "The future of distributed databases"
- **Product frame (weak):** "MongoDB announces new replication features"
- **You-frame (strong):** "What your replication setup is costing you (and how to fix it)"

The "you-frame" is not literally the word "you" — it is **the implied second-person subject**. If the reader cannot answer "what's in it for me?" from the headline alone, it has failed the test. Ogilvy: every headline must appeal to the reader's self-interest.

### 4. Numbers in headlines: odd, specific, and not too round

BuzzSumo's analysis of 100M headlines: numerals (digits, not spelled-out numbers) increase CTR by ~36%. Within numbered headlines:

- **B2C top performers:** 10, 5, 15, 7, 20, 6
- **B2B top performers:** 5, 10, 3, 7, 4, 6
- **Odd numbers** outperform even numbers for curiosity-driven content because they read as "researched" rather than "rounded." "7 lessons" feels found; "10 lessons" feels manufactured.
- **Specific numbers** outperform round numbers in evidence-driven contexts. "$1,247 saved" beats "$1,200 saved." "23% faster" beats "25% faster."
- **Avoid double numbers** in one headline. "5 ways to save 30% on..." dilutes both.

### 5. Curiosity-gap vs declarative headlines

Two dominant headline strategies, each with a failure mode:

**Declarative headlines** state the conclusion in the headline itself. "MongoDB beats Postgres on this workload by 4x." The reader gets the news; the click is for the proof. Good for hard news, technical posts, and anything where the body is evidence rather than story.

**Curiosity-gap headlines** withhold the conclusion to drive the click. "This database benchmark surprised everyone in the room." The reader clicks because the gap has been opened. Good for narrative pieces, profiles, and human-interest stories.

The failure mode of curiosity-gap is **clickbait**: a gap so wide the body can't close it. The 2014-era Upworthy/BuzzFeed playbook ("You won't believe what happened next") trained readers to distrust the form. By 2018 CTRs on those headlines had collapsed and major outlets had banned them. A modern curiosity-gap headline must be a **truthful gap** — the body must actually deliver the surprise the headline implies. Concrete-but-curious beats vague-curious every time. A 2024 study (PMC) confirmed that headlines conveying *just the right amount* of information maximize CTR; under-specified curiosity-gaps actually depress engagement.

### 6. Channel-specific length budgets

Headlines are not channel-portable. The same article needs different headlines for the homepage, social, search, and email.

| Channel | Sweet spot | Hard limit | Notes |
|---|---|---|---|
| **Web article H1** | 8–12 words | ~70 chars | Reads on-page; can breathe more. |
| **SEO title tag** | 50–60 chars | 60 chars (580–600 px) | Google truncates beyond pixel width; mobile shows up to ~70 chars. |
| **Meta description** | 150–160 chars desktop, 110–120 mobile | 160 chars | Not a headline but pairs with one. |
| **X / Twitter** | 60–70 chars | 280 chars | Under 70 chars lets retweeters add commentary. |
| **LinkedIn post** | 8–14 words | 150 chars before "see more" | Short headline + line break + body. |
| **Email subject line** | 30–50 chars | 60 chars before truncation | iPhone portrait shows ~35–40 chars. Front-load the hook. |
| **Push notification** | 25–40 chars | 50 chars iOS, 65 Android | Brutal — verb + benefit only. |

The pixel-width shift (Google moved off character counts in 2023) means modern SEO tooling measures pixels, not characters. Capital W is twice as wide as lowercase i. Treat character counts as estimates; verify with a pixel preview for anything search-critical.

### 7. News vs feature vs op-ed vs SEO headline registers

The same story takes four different headlines depending on the slot:

- **News headline:** Active verb, present tense, subject-verb-object. "Atlas adds vector search for Voyage embeddings." Designed to communicate the fact even if the reader doesn't click. AP style: no period, conservative with adjectives.
- **Feature headline:** Label or question form OK. "The quiet revolution in how databases handle vectors." Aims for tone and curiosity, not just facts.
- **Op-ed headline:** Argument-forward, often first-person-implied. "We're thinking about vector search wrong." Signals opinion through register; many publications italicize op-ed heds to differentiate.
- **SEO headline:** Keyword-front-loaded, declarative, designed for search-result skimming. "Vector search benchmarks: 12 databases compared (2026)." Often diverges sharply from the on-page H1.

A CMS that lets you set page title separately from H1 separately from social share title is exploiting this register difference. Most modern CMSes do. Use it.

### 8. The deck / subhead pairing

A subhead (also called a deck, dek, or sub-headline) is the second beat. It does one of three jobs:

1. **Add specificity** the headline omitted. Headline: "The case for synchronous replication." Deck: "Why eventual consistency cost us $400K last quarter."
2. **Identify the audience.** Headline: "A faster way to ship migrations." Deck: "For teams running >100 microservices on Kubernetes."
3. **Add evidence** that earns the headline's promise. Headline: "You're using indexes wrong." Deck: "A 6-month audit of 200 production clusters shows the same three mistakes."

Decks are usually 1.5–2× the length of the headline. They are not a place to restate the headline in different words — that wastes the second beat. If you can't write a deck that adds something, drop the deck.

### 9. A/B testing conventions

If you have the traffic for it (typically >5K impressions per variant per week), A/B test headlines. Conventions:

- **Test 2–3 variants**, not 10. Statistical power collapses with too many arms.
- **Vary one dimension at a time.** Number vs no-number. Declarative vs curiosity-gap. Long vs short. If you change everything, you learn nothing.
- **Hold the body constant.** A headline test with a different body is not a headline test.
- **Pick a metric and stick to it.** CTR is the obvious choice but pairs poorly with bounce rate. "Engaged sessions" (CTR × time-on-page) is a better single metric for editorial headlines.
- **Beware the optimization treadmill.** Headlines that win A/B tests trend toward higher curiosity-gap and lower truthfulness over many generations. Bake in a quality floor.

### 10. The Upworthy/BuzzFeed lesson

The 2012–2015 "curiosity-gap industrial complex" — Upworthy, BuzzFeed, ViralNova — proved you could 10x click-through rates with formulaic curiosity-gap headlines ("She walked into the room and what happened next will restore your faith in humanity"). It also proved the strategy was non-renewable. Reader trust eroded; Facebook's algorithm penalized the form by 2016; Upworthy's traffic collapsed 50%+ from peak. The lesson is not "don't use curiosity gaps" — it's that **headline strategies that exploit readers depreciate fast, and the depreciation is permanent at the publisher level**. A publication that burns trust to win one quarter loses it for years.

## Templates and examples

### Headline rewrite skeleton

When given a weak headline, run this checklist:

1. **Who is the implied subject?** If it's the writer or product, swap to the reader. (You-frame test.)
2. **What's the specific benefit or tension?** If it's abstract, replace with a concrete noun, number, or named thing.
3. **Is there a number?** If the piece has counted items, lead with the count. Odd, specific, BuzzSumo-top numbers.
4. **What channel?** Trim to the channel's character budget.
5. **News / feature / op-ed register?** Match the verb tense and tone.
6. **Curiosity-gap or declarative?** Pick deliberately. If curiosity, can the body deliver the promise?

### Example rewrites

| Weak | Stronger | Why |
|---|---|---|
| "The future of databases" | "What your database will look like in 2027 (and why it's not what you think)" | You-frame + curiosity gap + dated specificity. |
| "MongoDB launches Search Nodes" | "Search Nodes cut Atlas query latency 4x in production" | Verb-led, evidence-led, benefit-front. |
| "10 things to know about indexes" | "9 indexing mistakes I found in 200 production clusters" | Odd number, specific evidence, implied authority. |
| "An interview with our CEO" | "Why we killed our biggest product line" | Topic frame to tension frame. |
| "Improving developer experience" | "How we cut onboarding from 3 days to 90 minutes" | Concrete metric replaces abstract goal. |

### Op-ed headline patterns

- "Why [conventional wisdom] is wrong." ("Why 'serverless is cheaper' is a myth.")
- "The case for [unpopular position]." ("The case for synchronous replication.")
- "[We / Our industry] need[s] to stop [common practice]." ("We need to stop benchmarking on YCSB.")
- "[Named thing] is broken. Here's how to fix it." ("Vector benchmarks are broken. Here's how to fix them.")

### News headline patterns

- "[Subject] [active verb] [object]." ("Atlas adds vector indexes.")
- "[Subject] [verb], [secondary clause]." ("MongoDB acquires Voyage AI, expanding embedding portfolio.")
- "[Number] [things] [verb]." ("Three databases led 2026 benchmarks.")

### SEO headline patterns

- "[Keyword]: [benefit/comparison] ([year])" — "Vector databases: 12 options compared (2026)"
- "How to [task] in [tool]" — "How to enable Search Nodes in MongoDB Atlas"
- "[Number] [keyword] [qualifier]" — "8 indexing strategies for high-write workloads"

## Anti-patterns

- **The "we" headline.** "We released a new feature today." The reader doesn't care what *you* did; they care what they get.
- **The label headline.** "Distributed Systems Notes" or "Q3 Update." Names a folder, doesn't earn a click.
- **The double curiosity gap.** "The surprising thing we learned about a problem you didn't know you had." Compounding vagueness collapses the click rate.
- **The pun headline.** Ogilvy banned them in 1963 and the data still says he's right. Outside of headline-comedy publications, puns kill CTR.
- **Question headlines without payoff.** "Should you use a graph database?" If the answer is in the article, lead with the answer.
- **The bait-and-switch.** Headline promises X, body delivers Y. Worst-case for trust; eventually punished by every algorithm and every reader.
- **All caps or excessive punctuation.** ("BREAKING: This Will Change EVERYTHING!!!") Reads as spam; email providers downrank.
- **The keyword stuff.** "MongoDB Atlas Vector Search Database for AI RAG Embeddings in 2026" — written for crawlers, not humans, and modern search engines penalize it.
- **The under-specified curiosity gap.** "What we learned" without a hint of *what* you learned. Concrete-but-curious wins.
- **Decks that paraphrase the headline.** Wasted real estate.

## Decision heuristics

- **Is this for search?** Front-load the keyword, target 55–60 chars, use a year if the content is timely.
- **Is this for social?** Optimize for "would I retweet/repost this?" Headlines that are also takes get shared; headlines that are labels don't.
- **Is this for email?** Front-load the hook in the first 30 chars. Mobile preview is brutally short.
- **Is this a launch?** News register. Verb-led. Outcome in the headline, not the deck.
- **Is this an opinion piece?** Op-ed register. Argument in the headline, evidence in the body.
- **Is this clickbait-adjacent?** Apply the bait test: can the body deliver the headline's promise? If no, kill the headline, not the body.
- **Should this have a deck?** Yes if the headline alone leaves the reader uncertain who it's for or what they'll learn. No if it would just rephrase the headline.
- **How many variants to write?** 10 for routine. 16 for high-stakes. 3 for A/B test. 1 only if you've been writing headlines for 20 years.

## References

1. **Ogilvy, David. _Ogilvy on Advertising_ (1983)** — the source for the five-times rule, the 80-cent budget, the no-puns rule, and the appeal-to-self-interest principle. The headline chapter remains the single best 30 pages on the craft. [Notes summary](https://jimbouman.com/david-ogilvy-lessons-how-to-write/) | [Headline tips compilation](https://devedge-internet-marketing.com/2012/03/02/80-cents-of-dollar-spent-writing-headlines/)
2. **BuzzSumo headline studies** — the 100M-headline analysis behind the odd-numbers finding, the B2C/B2B number-performance rankings, and the curiosity-gap CTR data. [BuzzSumo resources](https://buzzsumo.com/resources/)
3. **Wiebe, Joanna. Copyhackers headline formulas** — the conversion-copywriting school's headline-formula doctrine, the "10 variants minimum" rule, and the voice-of-customer approach. [Copyhackers on headline formulas](https://copyhackers.com/2012/09/headline-formulas-and-the-science-of-high-converting-copywriting/)
4. **Zyppy / Search Engine Land studies on title tag pixel widths** — the modern SEO replacement for character-count rules. [Title tag length data](https://zyppy.com/title-tags/meta-title-tag-length/)
5. **PMC, "When curiosity gaps backfire" (2024)** — empirical study of 8,977 headline experiments on headline concreteness and click-through. [Study](https://www.ncbi.nlm.nih.gov/pmc/articles/PMC11704130/)
6. **Journalism University / WSJ headline-style differentiation** — the news-vs-feature-vs-opinion register distinctions and the italic-headline convention. [Headline formats](https://journalism.university/writing-and-editing-for-print-media/diverse-headline-formats-dynamic-storytelling/)
