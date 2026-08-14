<!-- hub-reference-banner -->
> **Reference file — part of the `career-and-formal-writing` hub.**
> Parent skill: `resume-and-cv-writing.md` §11 (Portfolio vs resume convention), which states the
> resume-level rule ("link to a GitHub or technical-blog URL in the header; never embed portfolio
> images in the PDF") but not how to curate the proof-of-work itself. This file covers that gap.

---

name: technical-portfolio-and-github-curation
description: >-
  Curating technical proof-of-work for engineers — pinned-repo strategy on a GitHub profile,
  README-as-sales-document conventions, the contribution-graph ("green squares") signal and why it is
  both overrated and partially real, live-demo hosting vs. code-only repos, choosing 3-5 projects to
  feature, when a dedicated portfolio website beats GitHub/LinkedIn, technical blog/talks/RFCs as
  signal, and open-source contribution vs. solo projects as a hiring signal. TRIGGER: GitHub profile
  for job hunting; pinned repositories; README for portfolio; contribution graph hiring; green squares
  vanity metric; live demo vs code repo; which projects to put on a portfolio; personal portfolio
  website vs GitHub; developer portfolio site sections; technical blog for career; conference talk
  portfolio signal; open source contributions vs side projects resume.
verified-as-of: 2026-07-21
---

# Technical Portfolio and GitHub Curation

## Overview

For an engineer, a GitHub profile (or portfolio site) is the writing sample a resume can only
summarize. The foundational source here is Marlow & Dabbish's 2013 CSCW study, which found hiring-side
interviewees treated GitHub activity traces as **more reliable than resume self-report**, because
they're third-party verifiable — not because they're exhaustive.[^1] Everything below follows from
that: curate for verifiable, low-effort-to-evaluate signal, not raw volume.

**Contents:** [1. Pinned-Repo Strategy](#1) · [2. README as Sales Document](#2) ·
[3. Contribution Graph](#3) · [4. Live-Demo Hosting](#4) · [5. Choosing 3-5 Projects](#5) ·
[6. Personal Portfolio Sites](#6) · [7. Blogs, Talks, RFCs](#7) · [8. OSS vs. Solo Projects](#8) ·
[9. Anti-Patterns](#9) · [10. Quick Reference](#10) · [References](#references)

---

## 1. Pinned-Repo Strategy {#1}

**The window is short.** Recruiter-facing sources converge on roughly 30 seconds for an initial GitHub
scan.[^2] Whatever isn't visible in that window is effectively invisible.

**Pinning fixes the ordering problem.** GitHub lets a user pin up to six repos above the default
reverse-chronological feed, which otherwise surfaces whatever was pushed most recently — often a
throwaway, not the strongest work.[^2]

**The six slots should tell a coherent story, not maximize count.** Six pins spanning a full-stack app,
a CLI tool, a data pipeline, and an API say "ships real projects across a range." Six near-identical
to-do-list clones say the opposite.[^2] Relevance to the target role beats stack purity.

**Quality over quantity applies profile-wide.** 10-15 well-documented repos beat a hundred empty forks;
archive or delete dead experiments rather than let them sit below the pins.[^2]

**The profile README is a distinct, high-value slot.** A public repo named exactly `username/username`
renders its README on the profile page, above the pinned grid — the single highest-visibility spot on
the whole profile.[^3] Use it for a short bio, current focus, and links to the 2-3 strongest projects —
not auto-generated stat widgets, which read as decoration substituted for substance.

---

## 2. The README as Sales Document {#2}

A project README in a portfolio context is a landing page: its job is converting a skimming reviewer
into a click on the live demo or the code.

**The 5/30/90 model:** five seconds to decide whether to keep reading, thirty to decide whether the
project is worth exploring, ninety (or a few minutes) to decide whether to dig into code.[^4]

- **0-5s:** title, one-line tagline, screenshot/GIF above the fold — often the only content seen at all.
- **5-30s:** the problem solved and why it matters, in plain language. "Task manager with real-time
  sync and offline support" beats "Todo app built with React," which tells the reviewer nothing they
  haven't seen fifty times.[^4]
- **30-90s:** live demo link (§4), feature list, enough install/usage detail to prove it runs.

**Recurring structure:** title + tagline, demo (link or GIF), features, install, usage, tech-stack
badges, contributing, license.[^4][^5] GitHub's own repository-best-practices docs frame README/license/
contributing/code-of-conduct as the documents that set project expectations — for a portfolio repo, the
README carries nearly all of that weight alone.[^5]

**Differentiate with substance, not a tech list.** What problem did this solve? What was the hardest
trade-off? What would change on a second pass? That last question signals the self-critique senior
engineers look for.[^4]

**A broken demo is worse than no demo** — it burns the trust the README just built. Test install/run
steps from a clean checkout, not from memory.[^4] Keep the README itself short — move deep architecture
notes to separate docs and link out.[^6]

**Baseline for context:** across public GitHub, only ~63% of repos have any README, 5.5% a contributor
guide, 2% a code of conduct.[^7] Clearing this bar decisively is a cheap way to stand out.

---

## 3. The Contribution Graph: Signal vs. Vanity Metric {#3}

**The overrated case (the stronger one):**
- Measures activity, not competence — a green graph of trivial whitespace/typo commits shows
  visibility, not skill.[^8]
- Textbook Goodhart's Law failure: a whole toolchain (streak scripts, backdated/empty commits) exists
  purely to keep the graph green, because it became a perceived screening filter.[^8]
- Only captures public activity — contract work, closed-source jobs, and private repos produce zero
  squares regardless of real output, undercounting anyone whose work isn't public.[^8]
- Streak-chasing has a documented burnout cost and rewards trivial commits over meaningful ones.[^8]

**The not-purely-noise case:** practitioners who source candidates directly report reading *pattern*,
not raw count — a scattered graph reads differently from a steady weekly cadence. Four or five squares
a week for six months signals durable habit; a 30-day pre-interview sprint reads as interview
performance, not practice.[^8]

**Guidance:** don't game it. Pattern over density is the only defensible read, and even that is
secondary to the pinned repos and their READMEs (§1-2). One critic's blunt summary: a recruiter who
stops at the graph has "completely missed the point" of what GitHub shows.[^8]

---

## 4. Live-Demo Hosting Conventions {#4}

**A code-only repo requires reconstruction work most reviewers skip.** A repo link shows a file tree
and README; the reviewer has to mentally simulate the running app or clone and build it. Under time
pressure that step gets skipped far more than completed — first-pass screening is a scan for "does this
demonstrably work," not a code review.[^9]

**Deployed demos remove that step.** A widely repeated (secondary-source, not a named study — treat as
directional) figure puts the share of hiring managers wanting a working deployed app, not just a repo,
around 84%.[^9][^10] The underlying mechanism isn't in dispute: a deployed project proves the candidate
carried work through build-ship-operate, not just "write code."[^9]

**Screening windows reinforce this** — estimated average time on a portfolio site during initial
screening is around 15 seconds, which a live demo fits and a clone-and-build workflow does not.[^9]

**Hosting has near-zero cost today:**

| Platform | Best for | Notes |
|---|---|---|
| GitHub Pages | Static HTML/CSS/JS, Jekyll | Free, no backend, simplest static setup |
| Netlify | Jamstack, forms, Git-based CI | Strong free tier, drag-and-drop or Git-connected |
| Vercel | React/Next.js, edge delivery | Free Hobby tier usually enough for a portfolio project |
| Render | Full-stack apps needing a real backend (Node/Python/Docker) | Free tier has cold starts; the practical Heroku alternative |

Common pattern: static site on Vercel/Netlify, individual backend-requiring demos on Render, all
cross-linked from the pinned README.[^11]

**Pair the demo with the repo — don't choose one.** Live link for the fast scan, clean repo for the
deeper technical follow-up.[^9][^12]

---

## 5. Choosing Which 3-5 Projects to Feature {#5}

**3-5 is a genuine cross-source consensus,** framed as depth over breadth: enough to show architecture,
trade-offs, and impact without wading through an uneven list; a longer shallow list reads as
unfocused.[^12] One source frames it directly: five unfinished projects are less convincing than one
finished end-to-end.[^12]

**Selection criteria, in priority order:**
1. **Relevance to the target role** — a project in-domain beats a flashier off-domain one, since it
   maps directly onto the job requirements.[^12]
2. **Range across the set** — frontend, backend, data, infra, a team project — not five variations on
   one skill.
3. **Unique value per project** — no slot should duplicate what another already proves.
4. **End-to-end completeness** — full-stack depth, a defensible data model, clean deployment, legible
   Git history for team projects.[^12]
5. **Live demo + a clear problem-solution-impact narrative** for each; a tech-stack list with no
   narrative undersells even a strong project.[^12]

**Prune actively.** Unpin/archive anything that no longer reflects current skill level. When a featured
project gets a real update, refresh its README and pin — a portfolio reads as a snapshot of current
ability, not a career museum.[^12]

---

## 6. Personal Portfolio Websites {#6}

**Core distinction:** a personal site answers "who is this," a portfolio answers "what can they do."
The strongest setup combines both rather than substituting one for the other.[^13]

**Skip a dedicated site when:** early-career with GitHub alone able to carry the weight; work has no
natural "gallery" to build a site around; or there's no real intent to maintain it — a stale site
listing dead projects hurts more than having none.[^13]

**Build one when:**
- **A single hub is needed** — one link (resume, LinkedIn, email signature) fanning out to GitHub, a
  blog, and a resume PDF, instead of forcing a reviewer to hunt across three profiles.[^13]
- **Deeper storytelling is wanted** than GitHub/LinkedIn conventions support — short case studies
  (problem, process, trade-offs, outcome) rather than a terse README or post.[^13]
- **Discoverability matters** — a personal site is indexable content the candidate controls, unlike a
  bare LinkedIn profile.[^13]

**Recommended sections, in order:** homepage (who/what/why in ~10 seconds, clear value prop) → about/bio
(fuller career narrative than a resume allows) → projects (same 3-5 discipline as §5, each a short case
study — "a portfolio without context is just a gallery"[^13]) → writing/talks (§7, if any) → contact
(direct email link, not just a form).

**Cross-link, don't silo:** reference the site from LinkedIn's Featured section, reference LinkedIn/
GitHub from the site's header, republish long-form writing to LinkedIn with a canonical link back — this
compounds discoverability instead of fragmenting presence.[^13]

---

## 7. Blog Posts, Talks, and RFCs as Portfolio Signal {#7}

**Technical writing is a durable but non-mandatory signal.** Multiple sources list blog posts,
articles, talks, and hackathon wins as portfolio strengtheners — published writing demonstrates domain
authority and the ability to explain, not just execute.[^14] A senior-level portfolio is sometimes
advised to mirror a conference-case-study shape: a few high-impact projects narrated with real
depth.[^14] Conference talks add a communication-skill signal a solo repo can't show at all.[^14]

**The credible counterpoint:** Julia Evans (jvns.ca), whose own blog is frequently cited as a genre
model, has argued against "everyone should blog" advice directly — most excellent developers she knows
don't blog, by a wide margin; of the developers she knows who do blog, most are excellent, which is
correlation, not requirement. Her stated motivation for writing is thinking, not career strategy.[^15]
**Takeaway:** a blog is a strong signal only if there's something real to say and the sustain-it
commitment is there; reviewers can tell writing done to think through a problem from writing done to
check a box.

**RFCs/design docs are the thinnest-evidenced category** — by definition usually not public. Practical
translation: convert real design-doc/RFC experience into a public artifact (a blog post walking through
a technical decision and its trade-offs) without disclosing anything proprietary; it does the same
signaling work.

**Placement:** a "Writing/Talks" section on a personal site (§6), or a pinned link if there's no site.
Don't force it — a dormant, once-a-year blog signals worse than no blog, the same way an abandoned side
project does.

---

## 8. Open-Source Contribution vs. Solo Projects {#8}

The single highest-leverage curation decision here, with the strongest evidence base of any section.

**The case for OSS contribution weighing more:**
- **Verifiability** — a merged PR is checkable directly (the diff, review comments, CI history), which
  is the mechanism behind the Marlow & Dabbish finding that GitHub activity traces are treated as more
  reliable than resume self-report.[^1]
- **Skills a solo project can't exercise** — reading unfamiliar code, following someone else's
  conventions, surviving external review across multiple rounds.[^16]
- **Proves the code survives a peer review that isn't the candidate's own** — a direct answer a
  self-reviewed solo repo can't give.[^16]
- **Differentiates from a crowded field** of near-identical clone projects (to-dos, weather wrappers,
  e-commerce clones) — a real merged external contribution reveals setup, collaboration, and
  communication all at once.[^16]
- **The strongest lever for career changers**, substituting for prior professional code they don't
  have.[^16]

**The genuine counterpoints:**
- **Context-dependent.** Large enterprises hiring for internal tooling often don't screen on public OSS
  history at all; it carries real weight at startups and dev-tools companies.[^16] Curate toward the
  employer, not a universal rule.
- **Some practitioners rank full ownership higher** — one informal but widely cited ranking places "own
  project as primary contributor, with real outside contributions" above "significant contributions to
  someone else's nontrivial project," above "solo project with only minor outside help," on the logic
  that architectural ownership is itself a signal OSS contribution doesn't fully replace. Others weight
  the two equally as long as the code is good.[^16]
- **Solo projects offer design freedom OSS contribution doesn't** — full choice of what to build,
  versus working inside someone else's architecture with little early foundational input.[^16]

**Synthesis:** treat the two as complementary within the same 3-5-project set (§5), not either/or.
Document either kind with quantified impact — "Implemented Redis caching (merged, PR #456), cutting
average API response time from 480ms to 220ms" is evidence; "Contributed to open source" is a label with
none.[^16] Execution quality inside the format matters more than which format was chosen.

**Caution on stars as a proxy:** don't let star count stand in for quality. An NDSS Symposium study
(~926,000 data points, PHP/Ruby/JavaScript) found only weak star-to-download correlation (0.14-0.47
depending on language) — stars measure visibility, not adoption or quality — and large-scale fake-star
manipulation, including malware-linked repos, has been documented at scale.[^17] A high star count is
worth mentioning; it isn't load-bearing evidence.

---

## 9. Anti-Patterns {#9}

| Anti-pattern | Why it fails |
|---|---|
| Leaving pins on GitHub's chronological default | Buries strongest work under whatever was pushed last |
| Six pinned repos that are near-identical clones | Wastes five of six slots; shows no range |
| README with only a title and stack badges | Fails the 5-second test; says nothing a hundred other repos don't |
| Broken or "coming soon" demo link | Worse than no demo — burns the trust the README built |
| Chasing streaks with backdated/trivial commits | Goodhart's Law failure; gamed graphs are increasingly recognized as gamed |
| Featuring 10+ projects at uneven polish | Dilutes attention from the 3-5 that would actually convert |
| Portfolio site abandoned mid-build or stale for years | A stale site with dead projects reads worse than no site |
| "Contributed to open source" with no PR link or metric | Unverifiable label; loses the exact reliability advantage OSS contribution offers |
| Citing star count as a quality stand-in | Weak correlation with real adoption; subject to fake-star manipulation |
| Blogging/speaking purely to pad a resume, then going dark | Reviewers spot box-checking; a dormant blog undercuts more than it helps |
| Embedding portfolio screenshots in the resume PDF | Breaks ATS parsing and bloats file size — link out (see `resume-and-cv-writing.md` §11) |
| Profile README filled with auto-generated stat widgets, no bio | Decoration substituted for substance in the highest-visibility spot on the page |

---

## 10. Quick Reference {#10}

1. Pin 6 repos showing range, not variations on one skill; archive dead experiments.
2. Give the `username/username` profile README a real bio and 2-3 project links, not stat badges.
3. Write each README to the 5/30/90 model: tagline+screenshot first, problem/why second, demo+features
   third.
4. Deploy a live demo for anything runnable (GitHub Pages / Netlify / Vercel / Render); pair it with the
   repo, don't replace it.
5. Narrow to 3-5 featured projects: relevance, range, and completeness over count.
6. Build a dedicated site only if it'll be maintained — otherwise GitHub + LinkedIn, done well, suffices.
7. Convert real design docs/RFCs/talks into a public post if there's something genuine to say; skip
   otherwise.
8. Weight OSS contributions and solo projects as complementary, both documented with quantified impact.
9. Ignore streak-gaming; let steady cadence speak for itself.
10. Never cite a star count as a quality claim on its own.

---

## References

[^1]: [Activity Traces and Signals in Software Developer Recruitment and Hiring — Marlow & Dabbish, CSCW 2013 (ACM DOI 10.1145/2441776.2441794)](https://www.researchgate.net/publication/262248832_Activity_traces_and_signals_in_software_developer_recruitment_and_hiring) — GitHub activity traces treated as more reliable hiring signals than resume self-report — verified-as-of: 2026-07-21
[^2]: [How Recruiters Actually Evaluate Your GitHub Profile — GitShare](https://gitshare.me/blog/how-recruiters-actually-evaluate-your-github-profile); [GitHub Recruiting: How to Find and Reach Engineers — Pin](https://www.pin.com/blog/github-recruiting/) — pinned-repo mechanics, 30-second scan estimate; practitioner consensus, not a controlled study
[^3]: [How To Create A GitHub Profile README — Monica Powell](https://aboutmonica.com/blog/how-to-create-a-github-profile-readme/) — mechanics of the `username/username` special repository
[^4]: [Writing a Great README — dev.to (binford2k)](https://dev.to/binford2k/writing-a-great-readme-5fn3); [Developer Portfolio Guide 2026 — Hakia](https://hakia.com/skills/building-portfolio/) — 5/30/90 attention model, reflective-question framing
[^5]: [Best practices for repositories — GitHub Docs](https://docs.github.com/en/repositories/creating-and-managing-repositories/best-practices-for-repositories) — official guidance on README/license/contributing roles
[^6]: [readme-checklist — ddbeck (GitHub)](https://github.com/ddbeck/readme-checklist/blob/main/checklist.md) — conciseness and "why not what" guidance
[^7]: [GitHub Octoverse — Grokipedia synthesis](https://grokipedia.com/page/GitHub_Octoverse) — 63% of public repos have a README, 5.5% a contributor guide, 2% a code of conduct; secondary synthesis of GitHub's Octoverse report — verified-as-of: 2026-07-21
[^8]: [Your GitHub Contribution Graph Means Absolutely Nothing — dev.to](https://dev.to/sylwia-lask/your-github-contribution-graph-means-absolutely-nothing-and-heres-why-2kjc); [How Developers Fake GitHub Contributions — Medium](https://medium.com/data-science-collective/developers-are-gaming-their-github-profiles-3f58f1f00c2a) — Goodhart's Law framing, streak-gaming mechanics, "shape over density" counter-read
[^9]: [Proof of Work for Developers — Fueler](https://fueler.io/blog/proof-of-work-for-developers-show-code-not-just-skills); [AI Project Demo vs GitHub: Fresher Interviews 2026 — FacePrep](https://faceprep.in/article/deploy-ai-project-streamlit-vs-github-for-fresher-interview-2026/) — 84% figure and 15-second window are secondary-blog sourced, no disclosed methodology; treat as directional
[^10]: [Developer Portfolio Guide 2026 — Hakia](https://hakia.com/skills/building-portfolio/) — corroborates live-demo preference
[^11]: [Best Free Hosting Platforms in 2026 — Appwrite](https://appwrite.io/blog/post/free-hosting-platform) — hosting comparison, free-tier specifics
[^12]: [How Many Projects Should I Have in My Portfolio? — Design Gurus](https://www.designgurus.io/answers/detail/how-many-projects-should-i-have-in-my-portfolio-as-a-software-engineer); [15 Software Engineer Portfolio Examples That Got Hired — ByAgentAI](https://byagentai.com/blog/software-engineer-portfolio-examples) — 3-5-project consensus and selection criteria
[^13]: [Personal Website vs. LinkedIn — CareerBldr](https://careerbldr.com/blog/personal-website-vs-linkedin/); [Why Your GitHub Profile Isn't a Portfolio — VibeCoders](https://www.vibecoders.bio/blog/github-profile-isnt-portfolio) — hub/discoverability rationale, recommended sections, abandonment-risk caution
[^14]: [How to Build a Tech Portfolio That Impresses Employers in 2026 — Tech Times](https://www.techtimes.com/articles/314992/20260311/how-build-tech-portfolio-that-impresses-employers-lands-you-job-2026.htm) — blog/talks as signal, conference-case-study framing
[^15]: [Julia Evans on Blogging — Chris Coyier](https://chriscoyier.net/2023/09/06/julia-evans-on-blogging/); [jvns.ca — Julia Evans](https://jvns.ca/) — counterpoint that blogging correlates with, but isn't required for, engineering excellence
[^16]: [Is Your Open Source Work Resume-Worthy? — Enhancv](https://enhancv.com/blog/open-source-on-resume/); [Why Open Source Contributions Matter In Hiring — daily.dev](https://recruiter.daily.dev/resources/open-source-contributions-matter-hiring/); [oss contributions vs side project — TeamBlind](https://www.teamblind.com/post/oss-contributions-vs-side-project-oufrxh7b) — verifiability/collaboration case, context-dependence, alternate ownership-first ranking, quantified-PR documentation example
[^17]: [The Fault in Our Stars: An Analysis of GitHub Stars as an Importance Metric for Web Source Code — NDSS Symposium](https://www.ndss-symposium.org/ndss-paper/auto-draft-490/) — ~926,000 data points, star-to-download correlation 0.14-0.47 by language; [How to Spot Fake GitHub Stars — dev.to](https://dev.to/alanwest/how-to-spot-fake-github-stars-before-they-burn-you-28op) — fake-star manipulation research
