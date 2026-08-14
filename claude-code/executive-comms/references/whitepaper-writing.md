<!-- hub-reference-banner -->
> **Reference file — part of the `executive-comms` hub.** Formerly the standalone `whitepaper-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: whitepaper-writing
description: Long-form B2B thought-leadership and lead-generation whitepaper craft — the problem → solution → proof → CTA arc, executive-summary discipline, Gordon Graham's three whitepaper types (backgrounder, numbered-list, problem/solution), sponsor-vs-vendor-neutral tone calibration, research-citation hygiene (lighter than academic, heavier than blog), 8–20 page conventions, design-driven layouts, gated-vs-ungated decisions, the "case study within a whitepaper" pattern, and references-appendix structure. TRIGGER when user asks to write/outline/draft/critique a whitepaper, B2B research report, technical brief, position paper, lead-gen long-form asset, or vendor-research piece; questions about executive-summary structure, gated content strategy, whitepaper length budget, citation density, or the problem/solution whitepaper arc. SKIP for owned-blog launch posts (use release-blog-and-launch-narrative); short-form marketing copy/landing pages (use sales-and-marketing-copy); academic-level peer-reviewed papers (use academic-and-citation-writing); architecture decision docs (use rfc-and-design-docs); customer success stories (use case-study-writing).
category: custom
tags: [writing, external-comms, whitepaper, thought-leadership, lead-generation, long-form]
---

# Whitepaper Writing

## Overview

A whitepaper is a long-form, authoritative B2B document that educates a business reader on a problem and positions a solution. It is neither a brochure nor an academic paper. The form is calibrated:

- **Lighter than academic**: persuasion is allowed; citation density is moderate; the author has a point of view.
- **Heavier than a blog post**: data and methodology are visible; claims are sourced; the document is designed to be downloaded, archived, forwarded, and cited.
- **Length**: 8–20 pages typical; 8–12 is the sweet spot for engagement; 16–20 is appropriate for technical or research papers.
- **Voice**: third-person company-as-narrator. First person ("we") only inside a clearly marked sponsor "About" block or in customer-quote pull-outs.

Whitepapers are the most expensive and the most resilient asset in a B2B content portfolio. A good one runs for 18–36 months as a lead-gen workhorse. A bad one is read by no one and lives in a CMS folder.

Gordon Graham ("That White Paper Guy") identifies three structural types:

1. **Backgrounder**: explains a technology, methodology, or product in depth. Vendor-neutral education. Used to establish authority.
2. **Numbered-list**: "5 Reasons to Adopt X," "10 Pitfalls of Y." Easier to write; lower authority ceiling; good for top-of-funnel.
3. **Problem/solution**: describes an industry-wide problem, evaluates current solutions, presents a new approach. Most complex; highest credibility ceiling; the dominant form for B2B lead gen.

This skill focuses primarily on the **problem/solution** form, with notes for the other two.

## Core concepts

### 1. The problem → solution → proof → CTA arc

The dominant whitepaper structure. Every section has a job; every section is sized.

| Section | Job | Pages |
|---|---|---|
| Title page + abstract/exec summary | Sell the download, then prove it was worth it in 60 seconds | 1–2 |
| Industry context / problem definition | Establish the stakes and the size of the problem | 2–3 |
| Why current solutions fall short | Honest review of competing approaches; not a hit job | 2–3 |
| The proposed solution / new approach | Vendor-neutral description of the method | 2–4 |
| Proof: case study, data, methodology | Evidence layer; this is where the paper earns its credibility | 2–4 |
| Implementation considerations | Real-world adoption guidance; this is where readers convert | 1–2 |
| Call to action + about + references | Next step, sponsor block, citations | 1–2 |

The arc is **inverted from a sales pitch**. A pitch starts with the product; a whitepaper ends with the product (or a CTA toward it).

### 2. Executive-summary discipline

The executive summary is **not** the introduction. It is the entire paper, compressed. Decision-makers read only it. Done well, it triggers either the download or the next-step action; done badly, it sinks the paper regardless of content quality.

**Rules**:

- 1–2 pages, no more
- Bullet-pointable; written as if a busy reader will skim only the bolded leads
- States the problem, the proposed solution, the most compelling proof point, and the call to action — in that order
- Stands alone; do not require reading the body to understand it
- Written **last**, after the full paper is drafted, even though it appears first

**Skeleton**:

```text
Executive Summary

The problem
[Industry] faces a [specific, sized problem] that [measurable consequence]. According to [source], [statistic that proves the problem is real].

Why current approaches fall short
[Approach A] addresses [partial scope] but [limitation]. [Approach B] [different limitation]. The gap is [specific unmet need].

A new approach
This paper proposes [solution category, not product name], built on [principle]. The approach is differentiated by [1–2 specific traits].

Evidence
A pilot deployment at [organization, possibly anonymized] reduced [metric] by [%] over [period]. Full methodology in Section 5.

What to do next
[Specific reader action]: download the implementation checklist, contact [sponsor], request a briefing, or read the full paper for the technical details.
```

### 3. Vendor-neutral vs sponsor-disclosed tone

Whitepapers are nearly always sponsored. Vendor neutrality does not mean pretending you have no point of view; it means the analysis would hold up if a competitor's analyst reviewed it.

**Vendor-neutral rules**:

- Name the problem before naming the product (or never name the product at all; let the CTA do that work)
- Discuss competing approaches honestly; their drawbacks should be real, not strawmen
- Cite third-party sources for problem-scoping statistics (analyst firms, regulators, peer-reviewed studies)
- Use the sponsor "About" block as the only place where company puffery is acceptable
- If the paper recommends a specific product, mark it clearly as the sponsored solution

**Sponsor-disclosure marking**:

- "This paper is sponsored by Acme Corp. The analysis was conducted by [author/firm]; the recommendations reflect Acme's perspective on the problem space."
- Place near the title page or in the executive summary footer.

A whitepaper that reads as a thinly disguised sales pitch is rejected by 60% of B2B decision-makers (per AdAge research cited by content marketing trade press). The neutrality discipline is not virtue signaling; it is a conversion ROI strategy.

### 4. Citation density and source hygiene

The right citation density is **2–6 cited claims per page** on average. Less and the paper reads as opinion; more and it reads as a literature review.

**Source hierarchy** (descending credibility):

1. Peer-reviewed academic research and government data (BLS, Census, FDA, etc.)
2. Major analyst firms (Gartner, Forrester, IDC, McKinsey) — but with the caveat that paywalled reports are often un-verifiable for the reader
3. Standards bodies (NIST, ISO, IEEE, W3C)
4. Industry surveys with disclosed methodology and sample size
5. The sponsor's own anonymized data (acceptable for proof sections; not for problem-scoping statistics)
6. Press articles, trade publications (acceptable for industry context; not for technical claims)
7. Vendor competitor whitepapers (use sparingly; cite to acknowledge prior art)

**Citation format**: footnotes or numbered-reference style; an appendix of references at the end. Inline parenthetical (Author, Year) is acceptable but visually heavy for designed PDFs.

**The "verify the verifiable" rule**: every claim of the form "X% of organizations report Y" must cite a source the reader can access. The single fastest way to disqualify a whitepaper is an uncited statistic.

### 5. Three whitepaper types — when to use each

| Type | Best for | Authority ceiling | Effort | Length |
|---|---|---|---|---|
| Backgrounder | Explaining a new category, methodology, or technology where the reader needs education before purchase | Medium-high | Medium | 6–12 pages |
| Numbered-list ("5 Pitfalls of...") | Top-of-funnel, easy-to-skim, viral-on-LinkedIn assets | Low-medium | Low | 4–8 pages |
| Problem/solution | Mid-to-bottom of funnel, where a buyer is evaluating approaches and needs to see proof | High | High | 10–20 pages |

The problem/solution form has the highest conversion-to-lead rate but also the highest writing cost — 8 to 12 weeks of effort typically, including research, drafting, design, and review cycles.

### 6. The "case study within a whitepaper" pattern

A whitepaper without a customer case study is a position paper; one with a case study is a sales asset. Embed a real (or anonymized) customer story in the proof section.

**Pattern**:

- Sidebar or dedicated section, 1–2 pages
- Customer name disclosed if approved; otherwise anonymized as "a Fortune 500 financial-services firm"
- Challenge → Approach → Result structure (mini case study)
- One direct customer quote, attributed
- Two-to-three quantified outcomes
- Linked to the full case study as a separate downloadable asset if available

The case-study sidebar should not duplicate the standalone case study; it should adapt and abbreviate it.

### 7. Length and design budget

Page-count is a tool, not a virtue. The right length is "as long as the argument requires, and not one page more."

**Practical budget**:

- 8–12 pages: standard problem/solution whitepaper for an established category
- 12–16 pages: technical whitepaper with extensive methodology or first-of-kind research
- 16–20 pages: research-grade paper with original primary data
- 4–8 pages: backgrounder or numbered-list
- 20+ pages: probably should be a multi-part series or an ebook; reading-completion rates collapse past 20

**Design budget**:

- Branded title page
- Table of contents (clickable in the PDF)
- 2–4 charts/infographics minimum
- Pull quotes every 2–3 pages to break density
- Sidebar boxes for case studies, definitions, technical asides
- Accessible PDF: alt text on images, semantic headings, readable on mobile

Pure-text whitepapers signal low investment. A flat, undesigned PDF reduces perceived authority before the first paragraph is read.

### 8. Gated vs ungated

The gating decision is a marketing-funnel decision, not a writing decision — but the writing should anticipate it.

**Gate when**:

- The paper contains original research with high differentiation value
- The funnel goal is lead capture; sales follow-up is staffed
- The audience is mid-or-bottom-funnel and willing to give an email for high value

**Ungate when**:

- The paper is foundational education that builds top-of-funnel authority
- SEO value is high; backlinks and search rank matter
- The audience is technical and form-allergic (e.g., developer-tools market)

**Hybrid pattern (modern best practice)**:

- First 2–3 pages (executive summary + first section) ungated, viewable on a page
- Full PDF gated behind a form
- Selected charts/quotes shareable on social without gate

Gated-form completion rates of 3–5x ungated engagement are typical for technical whitepapers — but only when the gated content is genuinely differentiated. Gating a generic whitepaper reduces total reach without producing usable leads.

### 9. Title and subtitle craft

A whitepaper title does two jobs: (1) earn the download click, (2) describe the paper accurately enough that it lives in search and citations for years.

**Title patterns that work**:

- "The State of [Industry] in 2026" — survey-based research paper
- "How [Audience] [Achieves Outcome]: A [Approach] Framework" — problem/solution
- "[N] Pitfalls of [Common Approach]" — numbered-list
- "[Technology] Explained: A Practical Guide for [Role]" — backgrounder

**Subtitle**: one-line argument or thesis. The subtitle is where the point of view lives.

> Title: **The State of Real-Time Data in Financial Services, 2026**
> Subtitle: *Why 73% of risk teams still operate on stale data — and the architectural shift required to fix it*

### 10. References, appendix, and methodology

A whitepaper without a references appendix reads as opinion. A whitepaper without a methodology appendix (for research papers) reads as marketing.

**References appendix**:

- Numbered or alphabetized
- Full citation: author, title, publisher, year, URL or DOI
- Listed at the end, on its own page(s)
- ~5–25 references typical; more than 40 starts to look academic

**Methodology appendix (for research-based papers)**:

- Sample size, geography, role/seniority breakdown
- Survey instrument or interview-protocol summary
- Time frame of data collection
- Statistical methods if applicable
- Disclosure of who funded the research and who conducted it

**Author bio**:

- Short paragraph: name, role, qualifications, prior work
- 1–2 sentences each for additional contributors
- LinkedIn or contact email for follow-up

## Templates and examples

### Full whitepaper outline (problem/solution, 10–12 pages)

```text
Title page (1 page)
  - Title + subtitle
  - Author name(s) + affiliation
  - Sponsor disclosure: "Sponsored by Acme Corp"
  - Publication date and version

Executive summary (1–2 pages)
  - Problem (1 paragraph)
  - Why current approaches fall short (1 paragraph)
  - New approach (1 paragraph)
  - Evidence preview (1 paragraph)
  - Call to action (1 sentence)

Section 1: The industry problem (2–3 pages)
  - Opening with a vivid scene or statistic
  - Definition: what exactly is the problem
  - Scope: how widespread, how costly, citing 3–5 external sources
  - Stakes: what changes if it is not solved

Section 2: Why current approaches fall short (2–3 pages)
  - Approach A: description, where it works, where it does not
  - Approach B: description, where it works, where it does not
  - Approach C: description, where it works, where it does not
  - The unmet need that motivates a new approach

Section 3: A new approach (2–4 pages)
  - The principle or architecture being proposed (vendor-neutral framing)
  - How it works: a diagram or two; a short technical explanation
  - How it differs from prior approaches in specific terms

Section 4: Evidence (2–4 pages)
  - Case study sidebar: real customer with quantified outcome
  - Pilot data or third-party validation
  - Methodology pointer to appendix

Section 5: Implementation considerations (1–2 pages)
  - Adoption prerequisites
  - Common pitfalls and mitigations
  - Suggested rollout sequence

Section 6: Call to action (1 page)
  - Specific next step (download a checklist, request a briefing, read a related paper)
  - About the sponsor (60–120 words, third person)
  - About the authors

Appendix A: References (1 page)
Appendix B: Methodology (if research-based, 1 page)
```

### Numbered-list whitepaper outline (5–8 pages)

```text
Title page (1 page)
  "5 Costly Pitfalls of [Common Approach] — and How to Avoid Them"

Introduction (1 page)
  Frame the problem; preview the five pitfalls

Pitfall 1 — Pitfall 5 (3–5 pages, ~1 page each)
  - Pitfall name (bolded)
  - What it looks like in practice
  - Why it persists
  - How to avoid it
  - Mini example or quote

Conclusion + call to action (1 page)
```

### Backgrounder outline (6–10 pages)

```text
Title page
"[Technology/Methodology] Explained: A Practical Guide for [Role]"

Introduction (1 page)
  Why this technology now; who should care

Section 1: Foundations (2–3 pages)
  Definitions and core concepts; vendor-neutral

Section 2: How it works (2–3 pages)
  Architecture or methodology overview

Section 3: Use cases (1–2 pages)
  Real-world applications with brief examples

Section 4: Selection criteria (1 page)
  How a buyer should evaluate options

Call to action + about (1 page)
```

## Anti-patterns

| Anti-pattern | Why it fails | Fix |
|---|---|---|
| Sales-pitch tone throughout | Rejected by 60% of B2B decision-makers; destroys credibility | Vendor-neutral problem analysis; sponsor visible only in About block and CTA |
| Executive summary is just the intro | Decision-makers read only the exec summary; if it lacks the proof and CTA, the paper is wasted | Write the exec summary last; ensure it stands alone |
| Uncited statistics | The single fastest credibility killer | Every "X% of..." claim cites a verifiable source |
| Strawman competitor analysis | Reads as dishonest; readers detect it | Describe competing approaches with their actual strengths; argue against them on real grounds |
| First-person voice ("we believe...") | Mixes whitepaper voice with sponsor voice | Third-person company-as-narrator; first person only in About block |
| 24-page paper with no design | Reading-completion collapses past 20; no design reduces perceived authority | 8–12 pages is the engagement sweet spot; invest in design |
| No case study, no proof | The paper is a position paper, not a sales asset | Embed a case study sidebar; cite real customer outcomes |
| No references appendix | Reads as opinion | 5–25 numbered references at the end |
| Methodology hidden or absent (research paper) | Procurement and technical readers reject findings they cannot validate | Methodology appendix with sample size, geography, dates |
| Gating top-of-funnel education content | Reduces reach without producing leads | Gate only differentiated, mid-or-bottom-funnel content |
| Title too clever to surface in search | Whitepapers live for years; SEO matters | Descriptive title; clever subtitle |
| Author bio is a marketing paragraph | Reads as inflated | One paragraph, factual: name, role, qualifications, prior work |

## Decision heuristics

**Which whitepaper type to write**

- Buyer is unaware of the category → backgrounder
- Buyer is aware but does not understand the failure modes → numbered-list
- Buyer is comparing approaches and you have a differentiated point of view → problem/solution

**Length budget**

- 4–8 pages: numbered-list, backgrounder, or single-claim research note
- 8–12 pages: standard problem/solution; the default
- 12–16 pages: technical or research-heavy
- 16–20 pages: original primary research with full methodology
- 20+ pages: split into a series; consider an ebook format

**Gated or ungated**

- Original primary research, mid-or-bottom-funnel audience, lead-capture staffed → gated
- Foundational education, top-of-funnel, SEO-driven → ungated
- Default for sponsored B2B whitepapers in 2026 → hybrid (executive summary on-page, full PDF gated)

**Citation density target**

- 2 sources per page minimum
- 6 sources per page maximum
- Methodology appendix carries unlimited density for original research

**Customer case study inclusion**

- Always include when the proof section is the dominant claim
- Anonymize if the customer cannot give written approval in your publication window
- Match the case study's depth to the whitepaper's depth (sidebar for short paper; full embedded section for long paper)

**When to bring in an outside author or analyst**

- The paper relies heavily on industry data the sponsor does not own → commission an analyst (Forrester, IDC, or a respected boutique)
- The audience is highly technical and the sponsor's internal voice is too marketing-heavy → use a technical author or known practitioner
- The credibility lift from a third-party byline outweighs the cost → byline an outside expert

## Cross-skill collisions and routing

- **owned-blog launch post** (short, narrative, marketing voice) → `content-and-marketing-writing` (references/release-blog-and-launch-narrative.md)
- **short-form marketing/conversion copy or landing page** → `content-and-marketing-writing` (references/sales-and-marketing-copy.md)
- **academic-level rigor with peer-reviewed citation standards** → `career-and-formal-writing` (references/academic-and-citation-writing.md)
- **architecture decision document or RFC for internal engineering** → rfc-and-design-docs
- **customer success story as a standalone asset** → `case-study-writing`
- **a press release announcing a whitepaper's publication** → `content-and-marketing-writing` (references/press-release-writing.md)
- **executive summary that is a one-page summary of a longer internal report** → `executive-comms`

## References

- [How to plan a problem/solution white paper — That White Paper Guy (Gordon Graham)](https://thatwhitepaperguy.com/how-to-plan-a-problem-solution-white-paper/)
- [The 3 Types of White Papers and When to Use Each One — CopyEngineer](https://copyengineer.com/3-types-of-white-papers/)
- [How To Write an Exceptional White Paper for Your B2B Brand — Content Marketing Institute](https://contentmarketinginstitute.com/content-creation-distribution/how-to-write-an-exceptional-white-paper-for-your-b2b-brand-examples)
- [Whitepapers in B2B Content Marketing in 2026 — Digital Marketing Knight](https://www.digitalmarketingknight.com/whitepapers-in-b2b-content-marketing/)
- [Gated vs. Ungated White Papers: A Complete B2B Guide — Zaphyr](https://zaphyrpro.com/gated-vs-ungated-white-papers-guide)
- [The marketer's guide to writing a persuasive B2B white paper — Bolder](https://www.bolderagency.com/journal/the-marketers-guide-to-writing-a-persuasive-b2b-white-paper)
- [White Papers for Dummies (Gordon Graham) — book reference](https://www.amazon.com/White-Papers-Dummies-Gordon-Graham/dp/1118496922)
