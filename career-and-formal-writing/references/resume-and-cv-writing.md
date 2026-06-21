<!-- hub-reference-banner -->
> **Reference file — part of the `career-and-formal-writing` hub.** Formerly the standalone `resume-and-cv-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: resume-and-cv-writing
description: Resume and CV writing — X-Y-Z achievement-bullet formula ("Accomplished X as measured by Y by doing Z" — Laszlo Bock / Google), quantification discipline, ATS parse-friendly formatting (single-column, standard headings, no tables/text-boxes/graphics), role tailoring, cover-letter complement, the skills-section debate, LinkedIn vs resume tone, academic CV vs industry resume distinction, employment-gap explanation, senior-IC vs management-track positioning, and the portfolio/resume convention split. TRIGGER when the user is writing or revising a resume, CV, cover letter, LinkedIn summary, achievement bullet, or asks about ATS, applicant tracking systems, X-Y-Z bullets, employment gaps, IC-vs-manager track resumes, academic-to-industry transitions, or "is my resume ATS-friendly". SKIP for grant or proposal writing (use proposal-and-grant-writing), generic technical/business prose (use writing-expert), sales or marketing conversion copy (use sales-and-marketing-copy — resume positioning is not conversion copy), structural rhetoric like BLUF or Minto pyramid (use rhetorical-frameworks-deep), executive-level comms unrelated to job search (use executive-comms).
---

# Resume and CV Writing

## Overview

A resume is read in 10-15 seconds on the first pass. It is parsed by an ATS (applicant tracking system) before any human sees it. It must survive the parse, then earn the human read, then earn the interview. Three different audiences (the parser, the recruiter skimming, the hiring manager evaluating) impose three different constraints on the same one or two pages.

This skill covers the genre-specific moves that recur across industry resumes, academic CVs, and cover letters. It is not a general writing reference — for prose craft, anti-AI-isms, BLUF, or audience calibration in business documents, see `writing-expert`. Use this skill for the specific decisions a job-seeking writer must make: how to write a bullet, how to format for an ATS, how to position for an IC vs management role, how to explain a gap, and when a CV is the right document.

When to invoke: the operator is writing or editing a resume, CV, cover letter, or LinkedIn profile section, or evaluating whether an existing document will pass an ATS or land an interview.

## When not to use

- The user is writing a grant proposal, NIH Specific Aims, logframe, or budget narrative → use `proposal-and-grant-writing`.
- The user wants generic prose craft, anti-AI-ism enforcement, BLUF/Minto/SCQA structure, or audience tone calibration for non-job-search writing → use `writing-expert`.
- The user wants sales pitches, landing-page copy, or conversion marketing — resume positioning is not conversion copy → use `sales-and-marketing-copy`.
- The user wants executive-level memo or board-deck phrasing unrelated to job applications → use `executive-comms`.
- The user wants structural rhetorical frameworks divorced from the resume genre → use `rhetorical-frameworks-deep`.

## Core Concepts

### 1. The X-Y-Z achievement formula

Popularized by Laszlo Bock, Google's former SVP of People Operations, in *Work Rules!*. The formula:

> **Accomplished [X] as measured by [Y] by doing [Z].**

- **X** — the achievement or outcome (not the task).
- **Y** — the measurement: percentage, dollar amount, users, time saved, headcount, latency.
- **Z** — the approach, method, or tools that produced it.

Weak: "Responsible for sales." Better: "Increased sales 25%." X-Y-Z: "Increased Q1 regional sales 25% ($1.2M ARR) by launching a partner-channel program across three Midwest accounts." The Y is non-negotiable — Google recruiters specifically flag the missing number.

### 2. Quantification discipline

Every bullet should contain a number unless physically impossible. Acceptable units: %, $, headcount, users, time saved, cycle-time reduction, NPS / CSAT delta, p95 latency, error-rate reduction, document count, training-hour count. When a metric is confidential, use a relative range ("reduced support ticket volume ~20%") or an order-of-magnitude proxy ("served 10,000+ users"). Never invent. The bullets you cannot quantify are the bullets you drop.

### 3. ATS parse-friendly formatting

The applicant tracking system parses your file before a human reads it. Failure modes (Jobscan, 2026 testing):

- **Two-column layouts and tables** scramble field order or drop content. Stick to single-column.
- **Custom or decorative headings** ("My Journey," "What I Bring") confuse the parser — it looks for "Experience," "Education," "Skills."
- **Text boxes, graphics, embedded images** are often dropped entirely.
- **Headers and footers** sometimes lose contact info; put name, phone, email in the body.
- **Non-standard fonts** can render as glyph junk. Use Arial, Calibri, Helvetica, or Georgia at 10-12pt.
- **PDF vs DOCX** — both work if text-selectable. Image-only PDFs (scanned, exported from design tools) fail the parse. Test with copy-paste.

### 4. Role tailoring

A single master resume cannot win across every application. For each target role: (a) reorder bullets so the top 3 in each job align with the target JD's stated priorities, (b) swap in keywords from the JD into the skills section, (c) rewrite the summary in the language the JD uses. Tailoring is the single highest-impact move; recruiter studies put callback-rate lift at 30-60% for tailored vs generic.

### 5. Cover letter as resume complement

The resume answers **what and where**. The cover letter answers **why and how** — why this company, why now, how your trajectory aligns. A cover letter that restates the resume is wasted. A cover letter that explains a non-obvious fit, a deliberate pivot, or a connection to the company's mission converts. Keep to ~3-4 paragraphs, one page, never two.

### 6. The skills section debate

Two schools of thought:

- **Pro skills section:** ATS keyword surface area, easy scan for recruiters, useful when career path is non-linear.
- **Con skills section:** "Skills" listed without context are unverifiable; better to demonstrate the skill in a bullet's Z component.

Synthesis: keep a skills section but limit it to (a) hard tools and technologies (Python, SQL, Figma, Kubernetes), (b) certifications, (c) languages. Do not list soft skills ("communication," "leadership") — they are not parseable signals and recruiters discount them. Move qualitative claims into achievement bullets where the Y proves the claim.

### 7. LinkedIn vs resume tone

The resume is third-person implicit and clipped ("Led a 5-engineer team that shipped..."). LinkedIn is first-person and conversational ("I lead a team that builds..."). LinkedIn is the long-form, searchable, networked surface — include the projects, the volunteer work, the writing, the recommendations. The resume is the focused, role-specific cut. Do not paste resume bullets verbatim into LinkedIn; rewrite for the medium.

### 8. Academic CV vs industry resume

A CV (curriculum vitae) is the academic standard. Several pages. Comprehensive: publications, presentations, grants, teaching, service, references. Used for postdocs, faculty, some government-lab and research-institute roles.

A resume is the industry standard. One page (early career) to two pages (senior, ≥10 years). Brief, achievement-focused, no publications list unless directly relevant, no references ("Available upon request" is dated — omit it).

Going from academia to industry: do not submit your CV. Translate. The 30-page CV becomes a 2-page resume that emphasizes transferable skills (project management of a multi-year study, grant writing as proposal writing, teaching as training and stakeholder communication), quantifies impact, and downplays publications unless they map to a research role.

### 9. Employment gap explanation

Unexplained gaps trigger reviewer suspicion. Explained gaps don't. Options:

- **Caregiving / family leave** — name it in one line on the resume, expand in the cover letter if relevant. LinkedIn's "Career Break" feature accepts this category officially.
- **Education / certification** — list as its own entry: "Independent study in [topic], earned [credential]."
- **Active job search + skill building** — frame as intentional: "Completed [course/credential], built [project], contributed to [open source]."
- **Health, personal, undisclosed** — a single line ("Career break, 2024-2025") is acceptable. Do not over-explain. Treat it as resolved, not as a problem.

Tone cue: how you frame the gap signals how the reader receives it. Apologetic framing reads as a red flag. Matter-of-fact framing reads as normal.

### 10. Senior-IC vs management-track resumes

A senior IC's resume emphasizes technical depth, architectural decisions, cross-team influence without authority, and impact through systems built. A manager's resume emphasizes team outcomes, headcount grown, attrition controlled, hiring loops run, business alignment, and impact through people developed.

Same person applying to both tracks needs **two resumes**. The IC version foregrounds artifacts and technical metrics (p95 latency, system uptime, lines of code shipped is not a metric — features shipped is). The manager version foregrounds team metrics (team size, eng-satisfaction score, promo rate, retention). Mixing the two signals confusion about which track the candidate wants.

### 11. Portfolio vs resume convention

Designers, front-end engineers, writers, researchers, and creatives present a portfolio alongside (not instead of) a resume. The resume is the structured summary. The portfolio is the evidence. The resume should link to the portfolio in the header. Do not embed portfolio images in the resume PDF — they break ATS and bloat file size. Engineers can substitute a GitHub or technical-blog link.

## Templates and Examples

### X-Y-Z bullet template

> [Action verb] [X — outcome] by [Z — approach/tool/method], [Y — measurement with unit].

### X-Y-Z bullet examples

- **Engineering:** Reduced p95 API latency from 850ms to 180ms by introducing a Redis read-through cache and refactoring the N+1 query pattern across three service boundaries.
- **Sales:** Grew enterprise pipeline 3x ($4.2M to $12.6M) in 12 months by repositioning the discovery motion around three industry verticals and rebuilding the qualification scorecard.
- **PM:** Shipped a self-serve onboarding flow that lifted activation from 32% to 51% (n=18,000 monthly signups) by sequencing three A/B-tested copy and UI changes over Q2-Q3.
- **Academic-to-industry translation:** Led a 4-year, $1.2M NIH-funded research program coordinating 6 institutions; managed cross-site IRB, data-sharing agreements, and quarterly milestone reviews — translatable as program management.

### Resume header template

```
[Full Name]
[City, State] | [phone] | [email] | [linkedin URL] | [portfolio/GitHub URL]
```

No street address. No photo (US convention; some European countries differ). No headshot. No "Objective" section — replaced by a 2-3 line professional summary.

### Section order (industry resume)

1. Header
2. Professional summary (2-3 lines, role-targeted)
3. Experience (reverse chronological)
4. Education
5. Skills (hard skills, tools, certifications, languages)
6. Selected projects / publications (optional, only if directly relevant)

### Academic CV section order

1. Header
2. Education
3. Academic appointments
4. Publications (peer-reviewed, then book chapters, then proceedings, then preprints)
5. Grants and funding
6. Invited talks and conference presentations
7. Teaching experience
8. Service (committees, peer review, editorial)
9. Awards and honors
10. References (named, with affiliations and contact)

### Cover letter skeleton

```
[Date]
[Hiring manager name, title, company]

Para 1 (the why): one sentence on what role you're applying for, one or two on why this company / why now. Cite something specific to the company.

Para 2 (the evidence): one or two of your strongest accomplishments that map directly to the JD's top requirements. Quantify.

Para 3 (the fit): connect your trajectory to where the role / team is going. Show that you've thought about the company's next 12 months, not just its job posting.

Closing: clear ask for an interview, contact info, thanks.
```

## Anti-Patterns

1. **Responsibility bullets instead of achievement bullets.** "Responsible for managing the customer-success team" describes the job description, not the candidate. Use X-Y-Z.
2. **Unquantified bullets across an entire role.** If you cannot put a number on any of three bullets under a job, you have not thought hard enough — or the role belongs lower on the resume.
3. **Two-column "designer-y" templates.** Beautiful in Figma, mangled by the ATS. Single-column or no callback.
4. **Soft skills as standalone skills entries.** "Strong communicator, team player, detail-oriented." Recruiters discount these. Demonstrate in a bullet's Z.
5. **The 30-page CV submitted for an industry role.** A signal that the candidate has not translated. Convert to a 2-page resume.
6. **Identical resume for every application.** Untailored resumes lose to tailored ones with weaker raw credentials.
7. **Apologetic gap framing.** "I was unfortunately out of work due to..." The tone primes the reader to see the gap as a problem.
8. **Mixing IC and manager track signals.** "Led architecture decisions and grew the team from 4 to 9." Pick one resume per track.
9. **Photos, color blocks, infographics on a US resume.** Drops ATS parse rate, introduces unconscious-bias risk on the human side.
10. **Cover letter that restates the resume.** Wasted page. The cover letter must add something the resume cannot.
11. **References on the resume.** Outdated. Omit. Provide on request.
12. **Embedded portfolio images.** Link to the portfolio; do not embed.

## Decision Heuristics

- **CV or resume?** Academia, postdocs, research-institute roles, some government labs → CV. Everything else → resume. Government labs vary; check the posting.
- **One page or two?** Less than 10 years' experience → one page. 10+ years → two pages, never more. Senior executives sometimes go three; risky.
- **PDF or DOCX?** PDF is safer (formatting preserved) if you generated it from a text-based source. DOCX if the application portal explicitly asks for it. Never submit an image-only PDF.
- **How tailored?** At minimum: rewrite the summary, reorder the top 3 bullets per role, swap 3-5 keywords in the skills section. 15-30 minutes per application is the floor.
- **Skills section length?** 8-15 items max. More than that signals you don't know what's relevant.
- **Bullets per role?** 3-5 for recent / relevant roles, 1-2 for older / less relevant. Never zero — every listed role needs at least one bullet of evidence.
- **Summary or no summary?** Yes if you are pivoting (industry, function, or seniority) or if the resume needs to argue a non-obvious fit. No if your title and most-recent role already make the case.
- **Should the cover letter template match the resume template visually?** Yes — same font, same header, same color (if any). Signals attention to detail.

## References

- [The XYZ Method Resume — Teal](https://www.tealhq.com/post/xyz-resume) — practical breakdown of the Bock formula with examples.
- [Google Recruiters Say Using the X-Y-Z Formula Will Improve Your Odds — Inc.](https://www.inc.com/bill-murphy-jr/google-recruiters-say-these-5-resume-tips-including-x-y-z-formula-will-improve-your-odds-of-getting-hired-at-google.html) — origin of the formula in Laszlo Bock's *Work Rules!*.
- [Anatomy of an ATS-Friendly Resume Format (2026 Checklist) — Jobscan](https://www.jobscan.co/blog/20-ats-friendly-resume-templates/) — current ATS-pass formatting rules and template testing.
- [5 Critical ATS Resume Formatting Mistakes to Avoid in 2026 — Jobscan](https://www.jobscan.co/blog/ats-formatting-mistakes/) — the specific layout failures (tables, two-column, custom headings) that drop parse rate.
- [The Resume vs. Curriculum Vitae (CV) — Harvard FAS Mignone Center for Career Success](https://careerservices.fas.harvard.edu/blog/2023/08/28/the-resume-vs-curriculum-vitae-cv/) — Harvard's canonical academic-vs-industry distinction.
- [Academic CVs to Industry Resume — University of Colorado Boulder Career Services](https://www.colorado.edu/career/graduate-students/careers-industry/resumes/academic-cvs-industry-resume) — translation guide for the academia-to-industry transition.
- [Individual Contributor vs Management Track: Career Path Guide 2026 — Hakia](https://hakia.com/careers/ic-vs-management/) — current IC-vs-manager track distinctions and resume implications.
- [How to Explain Employment Gaps in Resumes — AIApply](https://aiapply.co/blog/explaining-employment-gaps-in-resumes-and-interviews) — gap-framing options across resume and LinkedIn.
