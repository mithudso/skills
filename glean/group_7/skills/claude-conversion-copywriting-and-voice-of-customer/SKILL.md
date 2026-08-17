---
name: conversion-copywriting-and-voice-of-customer
title: Conversion Copywriting and Voice of Customer Research
version: 1.2.0
updated: 2026-06-17
description: >-
  Mine customer language and build evidence-based conversion copy strategy:
  VOC research (reviews, support tickets, transcripts, Reddit, G2);
  Jobs-to-be-Done message extraction; message hierarchy; hypothesis-driven
  copy tests. TRIGGER: find customer language for copy; mine reviews/Reddit/
  G2/Trustpilot for VOC; they-say/we-say gap; messaging hierarchy; JTBD copy
  inputs; post-purchase survey for objections; "what almost stopped you"
  survey; what to A/B test (headline, value prop, CTA, lead); objection
  inventory; message-to-market match; stages of awareness; HiPPO copy.
  SKIP: prose frameworks, AIDA/PAS, CTA wording → content-and-marketing-writing;
  writing researched copy as prose → content-and-marketing-writing; long-form
  sales letters, VSL scripts → direct-response-and-sales-letter-copywriting;
  A/B statistics → da-analytical-methods; persuasion theory →
  applied-psychology; prose editing → writing-expert;
  offer/value-prop design → offer-design-and-value-proposition.
keywords:
  - voice of customer
  - VOC research
  - review mining
  - message mining
  - jobs to be done
  - JTBD copy
  - message hierarchy
  - conversion copywriting
  - objection inventory
  - post-purchase survey
  - what almost stopped you
  - copy testing
  - hypothesis-driven copy
  - stages of awareness
  - HiPPO copy
  - Joanna Wiebe
  - Copyhackers
  - MarketingExperiments
  - message-to-market match
  - five stages of awareness
  - Eugene Schwartz
whenToUse:
  - "How do I find out what language my customers actually use to describe their problem?"
  - "Mine Amazon, G2, Reddit, or Trustpilot reviews for copy language"
  - "Build a messaging hierarchy from customer research data"
  - "Design a post-purchase survey to surface objections"
  - "What should I A/B test in my copy first — headline, CTA, or value prop?"
  - "Run a they-say/we-say audit on my landing page copy"
  - "Extract copy angles from JTBD interviews or outcomes research"
  - "Identify and rank objections for a landing page"
  - "How do I use the five stages of awareness to prioritize copy angles?"
whenNotToUse:
  - Writing or rewriting prose structures (use content-and-marketing-writing)
  - Expressing researched copy as finished prose (use content-and-marketing-writing)
  - Writing long-form sales letters, VSL scripts, or direct-response pages (use direct-response-and-sales-letter-copywriting)
  - Running A/B test statistical analysis (use da-analytical-methods)
  - Selecting persuasion tactics like social proof or scarcity (use applied-psychology)
  - General prose or voice editing (use writing-expert)
  - Designing the offer, value proposition, or positioning strategy (use offer-design-and-value-proposition)
related_skills:
  - content-and-marketing-writing
  - direct-response-and-sales-letter-copywriting
  - misc-catch-all
  - applied-psychology
  - da-analytical-methods
  - writing-expert
metadata:
  changelog:
    - "2026-06-17 sko v1.1.0->v1.2.0: fixed 4 High (dead SKIP routing targets sales-and-marketing-copy/persuasion-and-influence-psychology → content-and-marketing-writing/applied-psychology in description/banner/whenNotToUse), 5 Medium (ICE multiply vs sum, sample size power assumption + two-tailed correction, rehab-center lift claim inline citation, missing SKIP/whenNotToUse edge to direct-response-and-sales-letter-copywriting, missing related_skills entry for direct-response-and-sales-letter-copywriting)"
    - "2026-06-17 sko v1.0.0->v1.1.0: fixed 2 High (banner dead skill ID, description >1000 chars), 10 Medium (wrong SKIP IDs, missing related_skills/whenNotToUse entries, overclaim, undefined terms, B2B/enterprise gap, saturation criterion, post-test decision guidance, sources block moved to references/)"
---
# Conversion Copywriting and Voice of Customer Research

> **Spoke of the Writing and Documentation family.** Covers research/optimization beneath copy craft. Prose frameworks, AIDA/PAS, landing page structure, CTA wording → **content-and-marketing-writing**. A/B stats, causal-lift modeling → **da-analytical-methods**. Persuasion mechanisms (ELM, social proof, reciprocity) → **applied-psychology**.

---

## Core Principle

Key rule: **prospect's words beat yours every time.** Joanna Wiebe (coined "conversion copywriter", Copyhackers) frames process as *listening before writing* — copy already exists in reviews, support tickets, call transcripts, forums. Surface it, organize into hierarchy, test systematically. Research-and-testing discipline, not writing discipline.

---

## Part 1: Voice of Customer (VOC) Research

### What VOC Research Is

VOC = systematic collection of prospect's own language — words they use for their problem, desired solution, anxieties, outcomes. Goal: close **"they-say/we-say" gap** between how prospect talks and how marketing describes it.

**Sources ranked by richness:**

| Source | Signal type | Best for |
|--------|-------------|----------|
| Customer interviews | Unprompted emotional language | Core messaging, headline candidates |
| Post-purchase surveys | Friction buyers overcame | Objection copy, FAQ |
| Sales-call transcripts / Gong | Real objections in real words | Objection-handling copy |
| Support tickets / chat logs | Confusion, frustration language | Clarity fixes, FAQ, onboarding copy |
| Session replays (Hotjar, FullStory) | In-context hesitation/confusion signals | Friction points, unclear copy sections |
| Product reviews (Amazon, G2, Trustpilot, Yelp) | Scaled VOC; no interviewer effect | Review mining (see below) |
| Reddit / Quora / community forums | Unguarded problem language | Awareness-stage copy, blog hooks |
| Lost-deal CRM notes | Why-not-us language | Competitive differentiation copy |
| On-site exit surveys | In-the-moment friction | Page-level objection fixes |

**B2B/enterprise note:** Sparse public reviews (enterprise software, niche B2B) → skip to customer interviews regardless of customer count. Review-vs-interview fork applies to B2C/SMB markets with meaningful third-party review coverage.

### Review Mining at Scale

Review mining ("comment mining" / "message mining") = reading third-party reviews of your product, competitors, or adjacent solutions to find language patterns. Wiebe's foundational example: rehab center headline "If you think you need rehab, you do" came verbatim from Amazon book review on addiction (Wiebe, Copyhackers, Oct 2014; Havice, CXL, Aug 2015). Outperformed control ("Your addiction ends here") by large margin on CTA clicks.

**Review mining process:**

1. **Identify right review corpus.** Not necessarily reviews of your exact product — audience talking about the problem/solution space. Books, adjacent tools, services in same decision category — all valid.
2. **Collect with intent.** Look for: problem descriptions, desired outcomes, specific fears/anxieties, language about competitors/alternatives, "aha" moments after using solution, recurring phrases.
3. **Tag and theme.** Code each excerpt: Problem | Desired Outcome | Objection/Anxiety | Benefit/Transformation | Differentiator. Spreadsheet or affinity map.
4. **Build frequency ranking.** More often a theme appears in uncoached language → higher it belongs in message hierarchy.
5. **Swipe high-specificity phrases verbatim.** Worth swiping when specific enough only someone who lived the problem would write it — e.g. "I've tried six other tools and they all..." Vague phrases ("easy to use", "great product") not worth swiping — too generic to differentiate.

**Where to mine:**
- Amazon: search problem space or adjacent books/products
- G2, Capterra, Trustpilot, Yelp: direct competitors and adjacent categories
- Reddit: `site:reddit.com "[problem keyword]"` searches
- App Store / Google Play reviews
- Public community forums (Quora, niche Slack communities with public archives)

**Saturation signal:** stop when three consecutive review passes yield no new themes in any of the five tag categories. Practically: ≥50 tagged excerpts before declaring saturation on well-reviewed product; 20 minimum for sparse markets.

### The "They-Say/We-Say" Audit

Audit existing copy against VOC data before mining new sources:

| Column A: They say | Column B: We say | Gap? |
|--------------------|-----------------|------|
| "I'm drowning in spreadsheets" | "Streamlined workflow management" | Yes |
| "I needed it done by Friday and it actually worked" | "Fast implementation" | Yes |
| "It's the first tool my whole team actually uses" | "Collaborative platform" | Yes |

Every gap = rewrite candidate. Column A copy = test-ready headline.

---

## Part 2: Jobs-to-Be-Done Message Extraction

### The JTBD Framework for Copy Inputs

Jobs-to-Be-Done (Christensen/Ulwick) reframes copy question from "what features to describe?" to "what progress is customer trying to make?" Produces more specific, emotionally resonant copy.

**Three layers of the job:**

| Layer | Question | Copy use |
|-------|----------|----------|
| Functional job | What is customer trying to GET DONE? | Feature bullets, value prop body |
| Emotional job | How do they want to FEEL? | Headline, subheadline, lead |
| Social job | How do they want to be SEEN? | Transformation copy, testimonials |

Functional job gets you *considered*; emotional and social jobs get you *chosen*. Copy addressing only functional job (feature lists) loses to copy naming emotional job ("stop feeling anxious on Sunday night") in competitive markets.

### Wiebe's Messaging Hierarchy Process

Joanna Wiebe's methodology (Copy School / Copyhackers, CXL 2016):

1. **Research:** Mine reviews, interviews, surveys. Collect raw language.
2. **Synthesis — not summary.** Find patterns. Problems mentioned most? Outcomes wanted most? Rank by frequency of unprompted mention.
3. **Build message hierarchy.** Most frequently mentioned desired outcome or problem → primary message (above fold). Secondary messages fill body. Tertiary → FAQs or objection-handling copy.
4. **Wire the hierarchy.** Produce page wireframe or outline using hierarchy before writing prose. Hierarchy = structural decision, not writing decision.
5. **Swipe, combine, refine.** Paste verbatim VOC into wire as placeholder copy. Apply formulas (PAS, AIDA) to shape prose — words stay close to research.

**Messaging hierarchy template:**

```
Primary message:  [Most frequent desired outcome — exact customer language]
Primary proof:    [Evidence that supports primary message]
Objection 1:      [Most common anxiety] → Counter-objection
Objection 2:      [Second anxiety] → Counter-objection
Secondary benefit: [Second most frequent outcome]
Differentiator:   [Why you vs. the alternative they were using]
CTA context:      [What stage of awareness is the visitor at? Match lead.]
```

### Objection Inventory

Before writing landing page, produce exhaustive objection inventory. Sources: post-purchase surveys, sales-call recordings, support tickets, G2 reviews (negative and mixed).

**Objection mapping table:**

| Objection (their words) | Frequency | Counter-objection placement | Format |
|------------------------|-----------|----------------------------|--------|
| "I don't know if it works for [use case]" | High | Near CTA | Testimonial from that use case |
| "Seems expensive for what it does" | High | Pricing section | Value anchor / ROI statement |
| "Takes too long to set up" | Medium | Above fold (if conversion killer) | "Live in X minutes" proof |
| "I've tried similar tools before" | Low | Mid-page | Differentiation copy / comparison |

Place counter-objections at friction point — not buried in general FAQ.

---

## Part 3: The Five Stages of Awareness (Research Application)

Eugene Schwartz's five stages of awareness (*Breakthrough Advertising*, 1966) = **copy structure** tool AND **research classification** tool. Tag review/survey responses by awareness stage to identify which copy problems to solve.

| Stage | Prospect knows... | Research signal | Copy implication |
|-------|-------------------|-----------------|-----------------|
| Unaware | Nothing about the problem | No industry vocabulary; describes symptoms | Lead with story or startling claim |
| Problem aware | Has problem; no solution exists | "I've been struggling with X for years" | Lead with problem validation |
| Solution aware | Solutions exist; not yours | Mentions categories, not brands | Lead with unique mechanism |
| Product aware | Knows your product; hasn't bought | Comparison language, hesitation language | Lead with differentiator + objection handling |
| Most aware | Ready to buy; needs offer details | Price/trial/guarantee questions | Lead with direct offer |

**Practical use:** Tag VOC excerpts by awareness stage. If 80% of customer review data is "solution aware" but landing page leads with "most aware" direct offer → message-to-market mismatch. **Fix:** rewrite lead to match dominant awareness stage — introduce unique mechanism before asking for sale.

---

## Part 4: Survey Instruments for VOC

### The "What Almost Stopped You" Post-Purchase Survey

Highest single survey question for conversion research. Fielded immediately after purchase (confirmation page or post-purchase email):

> **"What, if anything, almost stopped you from completing your purchase today?"**

Why works: buyers who completed purchase overcame same objections blocking non-buyers. They answer honestly because purchase tension is resolved. Non-buyers won't respond to survey — recent buyers are best proxy for their hesitations.

Variants:
- "What was your biggest concern before buying from us?"
- "What almost stopped you from going ahead?"
- "Is there anything you almost needed to know before you decided?"

**Actionable output:** Categorize answers into objection buckets → rank by frequency → map to page element where objection not currently addressed → brief copy test.

Other high-value post-purchase questions (Conversion Rate Experts, 2024):
1. "How would you describe this product to a friend?" → headline language
2. "What problem were you trying to solve?" → messaging hierarchy inputs
3. "What made you choose us over the alternative?" → differentiator copy

### On-Site VOC Survey Triggers

| Trigger | Question | Insight |
|---------|----------|---------|
| Exit intent (pricing page) | "What's stopping you from signing up today?" | Pricing/value objections |
| 30-second dwell (homepage) | "What are you hoping to find on this page?" | Relevance gap |
| Post-signup | "What led you to sign up today?" | Pull-toward language |
| Churn / cancellation | "What's the main reason you're leaving?" | Product-promise gap |

---

## Part 5: Copy Testing Discipline

### What to Test (and in What Order)

Testing priority ladder — broad-to-narrow:

| Priority | Element | Why first |
|----------|---------|-----------|
| 1 | Radical redesign / full page | Tests entire message strategy; biggest possible signal |
| 2 | Value proposition / main headline | Largest single-element lever (MarketingExperiments, 2011) |
| 3 | Lead / above-the-fold block | Message-to-market match for traffic temperature |
| 4 | CTA text + context | "Try for free" vs. "See [specific outcome]" |
| 5 | Objection-handling block | Reduces abandonment at friction points |
| 6 | Social proof placement | Trust at the right moment |
| 7 | Single variable (button text, color, word swap) | Only after higher-order elements validated |

**Never start with single-variable tests on unvalidated page strategy.** Headline test on page with fundamentally mismatched message-to-market → noise. Fix strategy first.

### Hypothesis Format

Every copy test needs falsifiable hypothesis (Craig Sullivan / CXL format):

```
We believe that [COPY CHANGE]
for people [AUDIENCE SEGMENT]
will make [METRIC] improve
because [REASON GROUNDED IN VOC RESEARCH].
We will know this when we see [DATA SIGNAL].
```

Example:
> We believe that replacing "Streamlined workflow management" with "Stop drowning in spreadsheets" (from post-purchase survey data, n=87) for small-team project managers will increase trial signup rate because it matches exact language buyers used to describe their problem. We will know this when we see a statistically significant lift in trial signups with ≥95% confidence.

### Avoiding HiPPO Copy

HiPPO = "Highest Paid Person's Opinion." Copy written by committee or opinion = dominant failure mode. Antidote:

1. Source every copy claim to a VOC datum (quote, review, survey response)
2. Document what research said before writing
3. Let test decide — not review round

**Test prioritization with ICE scoring:** Rank hypotheses by Impact × Confidence × Ease (each 1–10). Multiply; higher products get priority. *Impact*: how much could this move primary metric? *Confidence*: how strong is VOC evidence? *Ease*: how fast to implement and measure? High-confidence tests (multiple VOC sources, high frequency) beat creative hunches.

### Sample Size Guidance (Decision Table)

| Baseline conversion rate | Minimum detectable effect (10% lift) | Approx. visitors needed per variant* |
|--------------------------|--------------------------------------|--------------------------------------|
| 1% | 0.1pp | ~160,000 |
| 2% | 0.2pp | ~80,000 |
| 5% | 0.5pp | ~31,000 |
| 10% | 1.0pp | ~16,000 |

*Two-proportion z-test, 80% power, α=0.05 (two-tailed). Use calculator (e.g. Evan Miller's) for exact baseline and MDE.

**Rule of thumb:** Can't reach required sample in under 4 weeks at current traffic → run qualitative copy testing instead — user interviews on copy, moderated comprehension checks, or five-second test (show page 5 seconds; ask what they remember). Never run underpowered tests. Statistical mechanics/significance calculations → **da-analytical-methods**.

### Reading Test Results

| Result | Action |
|--------|--------|
| Challenger wins (statistically significant) | Roll out challenger; brief next test on next priority ladder item |
| No significant difference | Declare learning (change didn't matter); move on; don't re-run without new hypothesis |
| Challenger loses | Analyze *why* (did copy create new objection?); don't abandon research — data may point to better angle |
| Inconclusive (low traffic) | Switch to qualitative methods; don't extend test indefinitely hoping for significance |

---

## Quick Reference: VOC Research Decision Table

| Situation | Recommended method | Time to data |
|-----------|-------------------|-------------|
| New product, no customers; B2C/SMB with reviews | Mine competitor / adjacent reviews + Reddit | 4–8 hrs |
| New product, no customers; B2B / sparse reviews | JTBD interviews with 5–7 target-market prospects | 1–2 days |
| Existing product, <100 customers | Customer interviews (5–7) | 1–2 days |
| Existing product, 100–500 customers | Post-purchase survey + review mining | 2–4 hrs setup |
| Existing product, 500+ customers | All-source synthesis (surveys + reviews + support tickets + session replays) | Full research sprint |
| Landing page not converting; no clear why | "What almost stopped you" survey + exit-intent survey | 2–4 weeks to collect |
| Need to prioritize what to test first | Frequency-rank all VOC themes; run radical redesign test first | — |

---

## Sources and References

See `references/sources.md` for full annotated citation list (16 sources).

Key sources: Wiebe / Copyhackers (review mining, messaging hierarchy, VOC methodology); CXL / Laja / Havice (ResearchXL model, copy testing, hypothesis format); MarketingExperiments / MECLABS (headline testing, A/B test progression); Schwartz *Breakthrough Advertising* (five stages of awareness); Conversion Rate Experts (post-purchase survey questions); JTBD.one / Boysen (outcome-driven copy generation).