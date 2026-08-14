<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Installed by `/dr` research (2026-06-10).
> Sibling topics (code smells in review context → `references/code-reviewer.md`; naming/quality conventions → `references/coding-standards.md`; design patterns → `references/coding-patterns.md`; debugging → `references/debugging-strategies.md`) are reference files under this hub.

---

---
name: development-antipatterns
title: Development Antipatterns
description: >
  Catalog of software development anti-patterns with detection signs and remediations —
  classics (God object/Blob, spaghetti, big ball of mud, golden hammer, lava flow,
  copy-paste, magic numbers, premature optimization with the full Knuth quote),
  architecture (distributed monolith, microservice envy, vendor lock-in and its inverted
  twin, accidental complexity), process (cargo cult, analysis paralysis, death march,
  bikeshedding, design by committee), and AI-era additions 2024–2026 (AI-amplified
  duplication/DRY erosion, complacency with AI-generated code, slopsquatting, prompt-and-pray,
  contradictory delivery-metric evidence). Includes the smells→refactorings mapping, detector
  disagreement data, the contested empirical record (Sjøberg/Hall vs Khomh/Palomba), when an
  antipattern is the rational choice, and criticism of antipattern catalogs themselves.
origin: local
version: "1.0.0"
updated: "2026-06-10"
---

# Development Antipatterns

Catalog of software development anti-patterns — design, architecture, process, and the 2024–2026 AI-era additions — with detection signs, remediations, the empirical evidence for and against the catalog itself, and when an "antipattern" is actually the rational choice. Confidence tags: [HIGH] = 3+ independent sources; [MEDIUM] = 2 or qualified; [LOW] = single-source/contested.

## Overview

"Antipattern" has a stricter definition than "bad thing": Andrew Koenig (JOOP, 1995) — *"just like a pattern, except that instead of a solution it gives something that looks superficially like a solution, but isn't one."* Brown, Malveau, McCormick & Mowbray's *AntiPatterns* (1998) made it a literary form: symptoms, root causes, and a **refactored solution** (40 antipatterns across developmental/architectural/managerial categories). Fowler's key qualifier is load-bearing throughout: **the same solution can be a pattern in one context and an antipattern in another.**

⚠️ **Citation-chain warning:** a large fraction of all antipattern web content (SourceMaking, c2 mirrors, listicles) traces to the single 1998 book; apparent multi-source agreement on definitions is often one source repeated. The genuinely independent legs: Riel 1996 (god class), Foote & Yoder 1997 (mud), Knuth 1974, Feynman 1974, Parkinson 1957/Kamp 1999, Conway 1968, Yourdon 1997, the empirical corpus (Sjøberg/Hall/Khomh/Palomba/Yamashita), Kapser & Godfrey, DORA 2024/25, METR 2025.

## Core Concepts

### Classic design/code antipatterns [HIGH]

- **God object / Blob** — independent pre-1998 root: Riel's heuristic 3.2 ("be suspicious of an abstraction named Driver/Manager/System/Subsystem"). One class hoards control/data. Remediation: Extract Class, redistribute responsibilities.
- **Spaghetti code** — convoluted control flow; rooted in the structured-programming debates (Dijkstra 1968). Remediation: incremental structuring under tests; rewrite only in extremis.
- **Big Ball of Mud** — Foote & Yoder 1997, with a preserved nuance the folklore dropped: the authors say their patterns "**are not anti-patterns** in the customary sense" — mud is *frequently the rational outcome* of throwaway code, piecemeal growth, and business pressure. Responses: consolidation phases, RECONSTRUCTION, Fowler's sacrificial architecture.
- **Golden hammer** — Kaplan's 1964 law of the instrument; entered IT as an antipattern via Brown et al. Detection: one technology in every proposal.
- **Lava flow / dead code** [MEDIUM — essentially single-origin, Brown et al.] — hardened dead-but-feared code from R&D-to-production; remediation: architecture + config management discipline, then delete (VCS is the net).
- **Copy-paste programming** [HIGH, contested] — Fowler/Beck call duplication the #1 smell, BUT Kapser & Godfrey ("'Cloning Considered Harmful' Considered Harmful", WCRE 2006/EMSE 2008) found **up to 71% of clones had a positive maintainability impact** (forking-as-sandbox, templating, hardware variants). Clone-elimination is default-good, not always-good — evaluate intent before refactoring.
- **Magic numbers** [MEDIUM] — unexplained literals; Replace Magic Number with Symbolic Constant; trivially lintable.
- **Premature optimization** [HIGH, widely misquoted] — Knuth 1974 in full: "We should forget about small efficiencies, say about 97% of the time… **Yet we should not pass up our opportunities in that critical 3%**" (and: an easily-obtained 12% improvement "is never considered marginal"). The antipattern is *unmeasured* optimization; weaponizing the quote to justify never profiling is the mirror-image antipattern.

### Architecture antipatterns [HIGH]

- **Distributed monolith** — services split on paper, coupled in deployment/data/sync chains. Detection: lockstep deploys; multiple services on the same DB tables; >2-hop synchronous chains; one library bump redeploys everything; one feature PR touches five repos. Consequence: "operational complexity of microservices with the coupling of a monolith — costs of both, benefits of neither." Remediation: DB-per-service/single-writer, events (outbox/pub-sub), contract versioning — or honestly re-consolidate.
- **Microservice envy** — ThoughtWorks Radar: **Hold** (since 2015); microservices trade development complexity for operational complexity and require CD/DevOps maturity; monolith-first.
- **Vendor lock-in — and its inverted twin** — Hohpe (martinfowler.com, 2019): lock-in is a 2×2 trade-off (switching cost × unique utility) in eight flavors; "you're also locked into what you built yourself." ThoughtWorks names the overcorrection "Generic cloud usage" (lowest-common-denominator cloud + home-grown abstraction layers) an antipattern too.
- **Accidental complexity** — Brooks 1986: accidental difficulties are self-inflicted and removable; essential complexity is the problem's own. Umbrella over pattern-overuse and framework accretion.
- **Abstraction inversion** [MEDIUM — folklore-grade sourcing] — re-implementing low-level features on top of high-level ones (semaphores on Ada rendezvous; fetch-whole-table-to-filter-in-app).

### Process / organizational antipatterns

- **Cargo cult programming/engineering** [HIGH] — Feynman 1974 ("the form is perfect… but no airplanes land") → Jargon File 1991 (ritual code without understanding) → McConnell 2000 (imitating the *forms* of successful organizations without the substance). Detection: elements no one can justify. Remediation: demand causal explanation; delete what can't be explained, tests as net.
- **Analysis paralysis** [MEDIUM-HIGH] — perfection-seeking analysis gridlock; remediation: iterative development, defer detail until knowledge exists.
- **Death march** [HIGH] — Yourdon 1997: parameters ≥50% beyond norms, "compensated" by forced overtime. Preserved criticism: DeMarco — death marches cluster on "characteristically insignificant" projects (important ones get resourced); Weir disputes the "norm" claim.
- **Bikeshedding** [HIGH] — Parkinson 1957 (triviality law) popularized for software by Kamp's 1999 FreeBSD email. Detection: debate volume ∝ 1/importance. Remediation: name it, timebox, delegate trivial decisions.
- **Design by committee** [MEDIUM] — strongest adjacent primary: Conway 1968 (designs copy communication structures; overpopulated design orgs produce poor designs). Remediation: small empowered design authority; inverse Conway maneuver.
- **Feature/scope creep, hero culture** [MEDIUM — practitioner consensus, weak primary literature] — uncontrolled growth; celebrated firefighting as broken-process signal (bus-factor risk, burnout).

### AI-era antipatterns (2024–2026)

- **AI-amplified duplication / DRY erosion** [MEDIUM-HIGH — single primary dataset (GitClear 211M lines), direction independently corroborated]: 2024 was the first year copy/pasted lines exceeded moved (refactoring-signature) lines; refactoring-associated lines fell ~25%→<10% (2021→2024). ⚠️ All 2025 press on this is ONE dataset re-reported. Independent corroboration of mechanism: FSE 2025 found commercial AI generators emit Type-1/2 clones up to 7.5% with vulnerability-propagation risk.
- **Complacency with AI-generated code** [HIGH] — ThoughtWorks Radar technique-antipattern: vigilance decays after positive experiences; agents enlarge change-sets, worsening reviewability; failure modes: automation bias, sunk cost, anchoring, review fatigue. Surveys [MEDIUM]: ~59% use AI code they don't fully understand (Clutch 2025); DORA 2025: 90% use AI, ~30% report little/no trust. Remediations: AI-specific review checklists (does every API exist? server-side authz? real business rules?), AI-attribution in PRs, TDD + static analysis in-loop.
- **Hallucinated dependencies → slopsquatting** [HIGH] — USENIX Security 2025 (16 models): ~20% of samples recommended nonexistent packages (5.2% commercial, 21.7% open-source); 43–58% of hallucinated names recur, making them registrable attack targets. Remediations: registry existence/age checks in CI, lockfiles, policy gates on AI-added dependencies.
- **Prompt-and-pray** [MEDIUM-HIGH] — business logic living entirely in prompts/chained LLM calls (O'Reilly Radar 2025; independently named by Rasa): chaining five <90%-reliable steps fails >40% of the time; undebuggable; violates least-agency. Remediation: "structured automation" — LLM for understanding, deterministic versioned workflow for business logic.
- **System-level delivery effects** [HIGH, contradictory by design]: DORA 2024 — +25% AI adoption associated with −1.5% throughput, −7.2% stability; DORA 2025 — throughput turned positive, **stability still negative** ("AI is an amplifier"); METR RCT 2025 — experienced OSS maintainers **19% slower with AI while believing +20% faster**; vs GitHub/Peng +55% and Google RCT +21% on greenfield/boilerplate tasks. Reconciliation: task- and codebase-maturity-dependent; self-reported productivity is unreliable.

### Criticism of antipattern catalogs themselves [HIGH]

- **Seemann's falsifiability test:** a true antipattern must be improvable *without significant trade-offs*; otherwise it's a disliked trade-off, not an antipattern. "Considered harmful" rhetoric smuggles universality into context-dependent judgments.
- **Antipattern dogma as cargo cult** [MEDIUM-HIGH]: reflexive Singleton bans producing over-complicated DI "serving theoretical purity"; "pattern obsession is itself an anti-pattern" (and Gamma is on record regretting Singleton's inclusion in GoF).
- **The empirical record is thin:** smelly classes do change/break more (Khomh 2009/12; Palomba 2018, 17,350 instances), BUT controlled/size-adjusted studies find effects vanish or are <10% and inconsistent — Sjøberg 2013: none of 12 smells significant after adjusting for file size and churn ("file size and revisions explain almost all variance"); Hall 2014: "**arbitrary refactoring is unlikely to significantly reduce fault-proneness and may increase it**"; field developers neither mention nor act on smells (Yamashita & Moonen). Best-supported use of the catalog: a shared *communication vocabulary*, not a defect-prediction science. Prioritize by size, churn, and smell *combinations*, not any single smell's presence.

## Detection / Remediation quick table

| Antipattern | Detection signs | Remediation | Acceptable when |
|---|---|---|---|
| God class / Blob | *Manager/System/Driver* name; hoards data+control; PMD/DECOR/JDeodorant flag | Extract Class; Move Method/Field | Thin façade; very small apps |
| Duplicated code | Clone detectors; copy/paste > moved trend | Extract Method/Class, Pull Up | Sandbox forks, cross-platform variants (≤71% of clones can be benign) |
| Spaghetti code | Convoluted flow, goto-dense legacy | Tests, then incremental structuring | Tiny scripts; generated FSMs |
| Big ball of mud | No discernible architecture | Consolidation phases; sacrificial architecture | Prototypes; survival-mode economics |
| Premature optimization | Optimizing before profiling | Profile first; optimize the measured 3% | Measured hot paths; easy 12% gains |
| Golden hammer | Same tech in every proposal | Requirements-driven selection | Deliberate stack standardization |
| Distributed monolith | Lockstep deploys; shared DB; sync chains | DB-per-service, async events, or re-consolidate | Time-boxed mid-migration states |
| Microservice envy | Services for fashion; no CD maturity | Monolith-first; 1–2 services | — |
| Vendor lock-in | High switching cost, no unique utility | Hohpe 2×2 trade-off analysis | High unique utility (accepted lock-in) |
| Cargo cult | "We've always done it this way" | Demand causal explanation; delete unexplainable | — |
| Analysis paralysis | Analysis never converges | Iterative increments | Safety-critical shifts the threshold |
| Death march | ≥50% beyond norms + forced overtime | Triage, renegotiate, or walk away | Never as steady state |
| Bikeshedding | Debate ∝ 1/importance | Name it; timebox; delegate | — |
| Complacency w/ AI code | Large AI PRs merged fast; author can't explain diff | AI-review checklist; attribution; TDD gates | Labeled throwaway prototypes |
| Slopsquatting exposure | Unknown packages in AI diffs | Registry checks, lockfiles, CI policy gates | — |
| Prompt-and-pray | Business logic only in prompts | LLM understands, deterministic workflow executes | Low-stakes demos |

## Troubleshooting (applying the catalog)

- **"Is this an antipattern or a trade-off?"** Run Seemann's test: can it be improved without significant trade-offs? If no — it's a trade-off; document it (ADR) instead of "fixing" it.
- **"Which smells do we fix first?"** Not by smell name: by file size × churn × smell combinations (the only factors with robust empirical backing). Tools disagree wildly (inter-tool agreement 0–71%, kappa ≈ 0 in places) — treat detector output as candidates, not verdicts.
- **"The team labels everything an antipattern."** That's the meta-antipattern: catalog-as-dogma. Require the labeler to name the context condition under which the construct would be acceptable; if they can't, the label is rhetoric.

## References

Access 2026-06-10. Chain flags: [BMM98] = Brown et al. *AntiPatterns* 1998 descendants; [FOW] = Fowler/Beck *Refactoring*; [GC] = single GitClear dataset; [USX] = single USENIX hallucination study.

1. Koenig, "Patterns and Antipatterns," JOOP 8(1), 1995 (via c2 wiki history + Fowler bliki). [primary]
2. Brown, Malveau, McCormick, Mowbray, *AntiPatterns*, Wiley 1998 (+ antipatterns.com, hillside.net). [primary; root of chain BMM98]
3. Foote & Yoder, "Big Ball of Mud," PLoP 1997 — hillside.net/plop/plop97/Proceedings/foote.pdf. [primary, definitive]
4. Riel, *Object-Oriented Design Heuristics*, 1996 (heuristics 3.2/3.4). [primary, independent of BMM98]
5. Knuth, "Structured Programming with go to Statements," ACM Comp. Surveys 6(4), 1974. [primary]
6. Fowler & Beck, *Refactoring* ch.3 (1999; 2nd ed. 2018) + refactoring.com/catalog. [primary; root of FOW]
7. Kapser & Godfrey, "'Cloning Considered Harmful' Considered Harmful," WCRE 2006 / EMSE 2008; + 2021 retrospective. [academic, definitive counter-evidence]
8. Sjøberg et al., "Quantifying the Effect of Code Smells on Maintenance Effort," IEEE TSE 39(8), 2013. [academic, null-after-controls]
9. Hall et al., "Some Code Smells Have a Significant but Small Effect on Faults," ACM TOSEM 2014. [academic]
10. Khomh et al. WCRE 2009/2012; Palomba et al. EMSE 2018 (17,350 instances). [academic, correlational]
11. Yamashita & Moonen, EMSE/IST 2013 + developer survey (simula.no). [academic]
12. Detection tooling: Moha & Guéhéneuc, "DECOR," IEEE TSE 2010; Fontana et al., JOT 2012 (kappa≈0); UFMG tool-agreement study (0–71.43%); JDeodorant (Tsantalis). [academic]
13. ThoughtWorks Technology Radar: "Microservice envy" (Hold); "Generic cloud usage"; "Complacency with AI-generated code." [industry assessment]
14. Hohpe, "Don't get locked up into avoiding lock-in," martinfowler.com, 2019-09. [expert]
15. Feynman, "Cargo Cult Science," 1974; Jargon File v2.5.1, 1991; McConnell, IEEE Software 2000. [primary ×3, independent]
16. Yourdon, *Death March*, 1997/2003; Weir (InformIT 2002) + DeMarco (OOPSLA 2001) critiques. [primary + criticism]
17. Kamp, bikeshed email, 1999 (phk.freebsd.dk); Parkinson 1957; Conway, "How Do Committees Invent?", 1968. [primary]
18. GitClear, "AI Copilot Code Quality" 2025 (211M lines). [industry data; GC — single dataset, all press derivative]
19. Spracklen et al., "We Have a Package for You!," USENIX Security 2025. [academic; USX]
20. Bowne-Anderson & Nichol, "Beyond Prompt and Pray," O'Reilly Radar, 2025-01-21; Rasa "Process Calling," 2025-05. [expert, independent pair]
21. DORA 2024 + 2025 reports (dora.dev); METR, "Measuring the Impact of Early-2025 AI…", 2025-07 (arXiv:2507.09089). [survey + RCT]
22. Seemann, "Some thoughts on anti-patterns," blog.ploeh.dk, 2019-01-21; Dave Thomas, "Design Patterns Are Not Design," PragDave 2025; ACCU Overload Singleton debate (Radford; Bashir). [expert criticism]
23. FSE 2025, "An Empirical Study of Code Clones from Commercial AI Code Generators." [academic]
