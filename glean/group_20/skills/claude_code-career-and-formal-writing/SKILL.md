---
name: career-and-formal-writing
description: >-
  Career, academic, legal & formal writing craft — genres with rigid conventions or high-stakes
  review by recruiters, ATS parsers, counsel, or governance bodies. TRIGGER: resume/CV; ATS
  optimization (Workday/Greenhouse/Lever/Taleo, keyword/JD mirroring, ATS-safe formatting);
  cover letter; JD or pay-transparency posting; performance/360 review; academic prose, thesis,
  citations (APA/Chicago/MLA/IEEE/Vancouver); legal-adjacent drafting (NDAs, MSAs, ToS,
  disclaimers); policy & governance docs; survey question design (Likert, NPS, bias-free).
  SKIP: general prose/voice → writing-expert; developer docs → technical-writing-craft;
  executive/board/persuasion → executive-comms; marketing/PR/newsletter →
  content-and-marketing-writing; auto-apply tools → automated-job-applications;
  hospitality job search → service-industry-job-hunting; service-industry resume, trail shift,
  or bar/restaurant interview → service-industry-resume-and-interview.
category: custom
tags: [writing, career, academic, legal, formal]
whenToUse:
  - Write or critique a resume or CV (ATS-safe formatting included)
  - Optimize a resume for ATS parsing — Workday/Greenhouse/Lever/Taleo, keyword/JD mirroring, scoring tools
  - Draft or tailor a cover letter for a job application
  - Author a job description that attracts candidates and meets pay-transparency law
  - Write a performance review, self-review, or peer/360 feedback
  - Write academic prose or format citations (APA, Chicago, MLA, IEEE, Vancouver)
  - Draft legal-adjacent prose for counsel review (NDAs, ToS, privacy notices, disclaimers)
  - Author a corporate, infosec, or governance policy document
  - Design bias-free survey questions, Likert scales, or NPS wording
related_skills:
  - writing-expert
  - technical-writing-craft
  - executive-comms
  - content-and-marketing-writing
  - automated-job-applications
model: Codex-opus-4-8
effort: high
version: "1.1.0"
updated: "2026-06-23"
---
# Career, Academic, Legal & Formal Writing

Genres here share one trait: governed by **external conventions you don't invent**. Resume judged against recruiter/ATS expectations; citation right or wrong by published style manual; privacy notice read by counsel and regulator; policy audited; survey question measured for bias. Craft here is less about voice, more about meeting a fixed standard precisely.

This hub covers **eight domains** across **9 on-demand reference files**: resumes/CVs, cover letters, job descriptions, performance reviews, academic/citation writing, legal-adjacent drafting, policy/governance writing, and survey question design. The ninth file is an ATS deep-dive extending the resume domain.

## How this hub works

When a task matches a row in the routing table below, **Read the listed `references/<name>.md` before producing any draft, citation, clause, or rubric** — the reference carries frameworks, templates, formulas, and anti-patterns. Answer from the table alone only for a one-line "which approach applies" pointer.

Legal-adjacent and policy/governance rows carry a hard constraint: draft text for review by qualified counsel or governance owner — do not substitute for it. Flag that boundary in any output touching those domains.

## Sub-skill routing table

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `resume-and-cv-writing` | Resume/CV writing — X-Y-Z achievement bullets, quantification discipline, ATS-safe formatting overview, resume vs CV distinction, reverse-chronological vs functional layouts | `references/resume-and-cv-writing.md` |
| `ats-resume-optimization` | ATS deep dive — load *in addition to* `resume-and-cv-writing` when the task is specifically about machine parsing, keyword scoring, or per-vendor rules; covers Workday, Greenhouse, Lever, Taleo, scoring tools (Jobscan, Resume Worded, Teal), and service/hourly worker guidance | `references/ats-resume-optimization.md` |
| `cover-letter-writing` | Craft persuasive cover letters — three-part skeleton (why this role / why you / why now), tailoring to JD, avoiding resume-restatement | `references/cover-letter-writing.md` |
| `job-description-writing` | Recruiting-facing JD craft — must-have vs nice-to-have split, inclusive language, US pay-transparency salary-range disclosure | `references/job-description-writing.md` |
| `performance-review-writing` | Write performance reviews evidence-based, fair, actionable — SBI/STAR evidence pattern, bias-mitigation language, rating-to-narrative alignment | `references/performance-review-writing.md` |
| `academic-and-citation-writing` | Academic prose, scholarly argument, citation discipline — APA 7th, Chicago 17th, MLA 9th, IEEE, Vancouver; thesis structure, source integration | `references/academic-and-citation-writing.md` |
| `legal-adjacent-writing` | Draft legal-adjacent prose for counsel review — disclaimers, MSAs, NDAs, ToS, privacy notices, security-incident language ⚠ counsel review required | `references/legal-adjacent-writing.md` |
| `policy-and-governance-writing` | Author prescriptive policy docs — corporate, infosec, acceptable-use policies, governance frameworks, RACI ownership ⚠ governance-owner review required | `references/policy-and-governance-writing.md` |
| `survey-question-writing` | Bias-free survey question design — Likert scale construction, NPS wording (Reichheld 2003), single-select / multi-select / open response, avoiding double-barreled and leading items | `references/survey-question-writing.md` |

## Cross-hub note

One of five hubs in the `writing` family (alongside `writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`). Route out when the task leaves the formal/career/academic/legal domain:

- General prose, voice, tone, word-level editing → `writing-expert`
- Software, API, product, developer docs → `technical-writing-craft`
- Executive, board, business, or persuasion writing → `executive-comms`
- Marketing, PR, launch, newsletter, social copy → `content-and-marketing-writing`
- Auto-apply tools, LinkedIn Easy Apply automation, job-bot strategy, mass-application ops → `automated-job-applications` (not a writing-family hub; a job-search ops skill)

When a task spans two hubs (e.g., executive bio that is part profile, part resume), load the closest-fit reference here and pull the sibling hub for the other half.

<!-- cross-hub-map -->
## Cross-hub map — where every writing topic lives

Family split across these hubs. If the task's deep material is **not** in the routing table above, it is a reference file under a sibling hub — **activate that hub or `Read` its `references/<name>.md` directly**. Every former standalone skill is now a reference under one of these hubs (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `writing-expert` | Prose craft, voice, style, editing | `references/editing-and-revision.md`, `references/rhetorical-frameworks-deep.md`, `references/storytelling-and-narrative.md`, … |
| `technical-writing-craft` | Technical & product writing — docs, specs, engineering comms | `references/api-docs-craft.md`, `references/howto-writing.md`, `references/tutorial-writing.md`, … |
| `executive-comms` | Executive & business comms — leadership, persuasion, decks | `references/one-pager-writing.md`, `references/okr-writing.md`, `references/proposal-and-grant-writing.md`, … |
| `content-and-marketing-writing` | Content, marketing & external comms — PR, newsletters, launch, social | `references/sales-and-marketing-copy.md`, `references/press-release-writing.md`, `references/crisis-pr-writing.md`, … |
| `career-and-formal-writing` | Career, academic, legal & formal writing (this hub) | `references/resume-and-cv-writing.md`, `references/ats-resume-optimization.md`, `references/cover-letter-writing.md`, … |
