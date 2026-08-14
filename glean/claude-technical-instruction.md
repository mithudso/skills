# technical-instruction

**Category:** Science, Biology & Medicine
**Platform:** Claude
**Original Path:** claude/standalone/technical-instruction

## Description
Hub for teaching engineers — design, deliver, assess & measure technical training and engineering education; routes to 10 spokes. TRIGGER: course/curriculum/learning-path design (ADDIE/SAM/Dick&Carey/Gagné/UbD, learning objectives, scope-and-sequence); certification & credential design (job-task analysis, item writing, CTT/IRT psychometrics, cut scores, ISO 17024/NCCA, Open Badges); hands-on labs, participatory live coding, DevRel/customer academy; teaching troubleshooting & diagnostic reasoning (cognitive apprenticeship, productive failure, illness scripts); MongoDB University & certification program; training measurement & ROI (Kirkpatrick, Phillips ROI, xAPI/SCORM/cmi5, transfer); GenAI in education (AI-assisted ID, LLM tutors/ITS, AI-resistant assessment); expert-knowledge elicitation (CDM/ACTA/PARI/GDTA); skills taxonomies & competency frameworks (O*NET/ESCO/SFIA, CBE/CBT); is-training-the-right-fix / HPT (Gilbert BEM, EPSS, job aids). SKIP: learning SCIENCE (cognitive load, retrieval, spacing) → applied-psychology; single tutorial/Diátaxis → technical-writing-craft; how to debug → software-engineering-patterns; live MongoDB/Atlas diagnosis → atlas-diagnostics-expert; survey/Likert & measurement theory → career-and-formal-writing/da-1-foundations-theory.

---

# Technical Instruction & Engineering Education (hub)

Routing hub for the discipline of **teaching engineers**: how to design, deliver, assess, measure, and
AI-augment technical training and engineering education. This hub **routes** — it holds no depth of its
own; each spoke below is a full top-level skill you invoke directly. Pick the spoke that owns the task;
when a question spans two, start with the more specific one and let it defer back here.

## Routing table (the spokes)

| If the task is… | Go to |
|---|---|
| **Is training even the right fix?** Diagnose a performance gap and pick the intervention BEFORE designing training: ISPI HPT model, Gilbert's BEM (environment-first), Mager-Pipe flowchart, Rummler-Brache 3-levels/9-variables, the intervention taxonomy, **performance support / EPSS / DAPs**, job-aid design, the Five Moments of Need, learning in the flow of work, the performance-consultant role | **human-performance-technology** |
| Design a multi-lesson **course, curriculum, or learning path**; choose a process model (ADDIE, SAM, Dick & Carey); backward design / UbD; scope-and-sequence & prerequisite DAGs; write measurable **learning objectives**; constructive alignment; Gagné's Nine Events; Merrill's First Principles; format (self-paced/cohort/blended) | **instructional-design-course-architecture** |
| Design or accredit a **certification exam / credential**: job task analysis, test blueprints, MCQ/performance item writing, CTT/IRT psychometrics, **cut scores** (Angoff/Bookmark), ISO 17024 / NCCA, DIF, item-bank governance, Open Badges 3.0 / micro-credentials | **assessment-certification-design** |
| **Deliver** technical training: hands-on labs, participatory live coding, formative diagnostics, faded examples/Parsons problems, teaching SQL/NoSQL & programming misconceptions, mixed skill levels, DevRel/customer academy, partner/SI enablement, modality trade-offs, screencasts | **technical-training-delivery** |
| Teach the **pedagogy of troubleshooting & diagnostic reasoning**: cognitive apprenticeship, mental-model instruction, productive failure (Kapur), novice→expert trajectory, illness scripts, game-day/fire-drill as a *learning* event, diagnostic-competency assessment | **teaching-troubleshooting-diagnostic-reasoning** |
| **MongoDB University** platform + the **MongoDB certification** program: which cert/learning path, course catalog & labs, exam format/cost/proctoring/retakes, Credly badges, student/educator/partner discounts, positioning University to a customer team | **mongodb-university-certification** |
| **Measure** training/enablement effectiveness & ROI: Kirkpatrick Four Levels & New World model, Phillips ROI (Level 5), xAPI/SCORM/cmi5 & LRS, training transfer (Baldwin-Ford, LTSI), leading/lagging KPIs, L&D dashboards | **learning-measurement-evaluation** |
| Use **GenAI in education**: AI-assisted instructional design (HITL), LLM tutors/ITS (Khanmigo, LearnLM, Bloom 2-sigma), automatic item generation quality, GenAI in CS education, AI-resistant assessment, FERPA/COPPA & equity | **genai-education-instructional-design** |
| **Elicit expert cognitive knowledge** before course design: CDM / ACTA / PARI / GDTA interviews, expert blind-spot problem, Cognitive Demands Tables, Decision Requirements Tables, HTA vs CTA | **cognitive-task-analysis** |
| Build a **skills taxonomy / competency framework** (O*NET, ESCO, SFIA 9, Lightcast), KSAOs & job-task analysis, DACUM, **competency-based education/training (CBE/CBT)**, mastery-based progression, skills-gap analysis, the skills-framework-to-credential link | **skills-taxonomies-cbe** |

## Reference files (load on demand)

Each spoke is now an on-demand reference file under `references/`. When the routing table sends you to a spoke, **Read its `references/` file** before answering at depth.

| Spoke | Reference file |
| --- | --- |
| `human-performance-technology` | `references/human-performance-technology.md` |
| `instructional-design-course-architecture` | `references/instructional-design-course-architecture.md` |
| `assessment-certification-design` | `references/assessment-certification-design.md` |
| `technical-training-delivery` | `references/technical-training-delivery.md` |
| `teaching-troubleshooting-diagnostic-reasoning` | `references/teaching-troubleshooting-diagnostic-reasoning.md` |
| `mongodb-university-certification` | `references/mongodb-university-certification.md` |
| `learning-measurement-evaluation` | `references/learning-measurement-evaluation.md` |
| `genai-education-instructional-design` | `references/genai-education-instructional-design.md` |
| `cognitive-task-analysis` | `references/cognitive-task-analysis.md` |
| `skills-taxonomies-cbe` | `references/skills-taxonomies-cbe.md` |

## Cross-hub map (cede to these neighbors — this hub does NOT own them)

| Topic | Owner | Why it is theirs, not this hub's |
|---|---|---|
| The **science** of how people learn — cognitive load, retrieval practice, spacing, deliberate practice, expertise | **applied-psychology** (`references/learning-and-expertise-psychology.md`) | This hub *applies* learning science; it does not own the mechanisms. |
| Writing one **tutorial, how-to, reference, or explanation** doc (Diátaxis CRAFT) | **technical-writing-craft** (`references/tutorial-writing.md`) | Single-doc craft is below course-architecture altitude. |
| **How to actually debug / run a root-cause analysis** on code or systems | **software-engineering-patterns** | This hub teaches people *to* debug; it does not do the debugging. |
| Diagnosing a **live MongoDB/Atlas** problem; MongoDB technical content (MQL, indexes, schema, perf) | **atlas-diagnostics-expert**, **mongodb-expert**, **mongodb-atlas-expert** | Owns the technical diagnosis & content the training teaches. |
| Enablement as an **account deliverable** (EBR/QBR, support plan, customer health) | **tam-operations** | Account-level packaging of enablement, not the instruction itself. |
| **Survey / Likert / NPS** item wording | **career-and-formal-writing** | Survey items are a formal-writing genre, not certification psychometrics. |
| **Research measurement theory** (reliability, validity, levels of measurement) | **da-1-foundations-theory** | Foundational measurement theory beneath any technique. |

## Notes

- `mongodb-university-certification` is **cross-cutting**: it lives in this family (it is a learning/cert
  surface) and is also reachable from the MongoDB hubs (`mongodb-expert`, `atlas-diagnostics-expert`,
  `mongodb-atlas-expert`) for the MongoDB-customer enablement angle. Its primary home is this hub; the
  MongoDB hubs carry a cross-edge to it.
- Each spoke remains a standalone top-level skill (this hub routes, it does not fold). Invoke a spoke
  directly when you already know the sub-domain; use this hub when you do not.

<!-- cross-hub-map -->
## Cross-hub map — where every technical-instruction topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `technical-instruction` | Technical Instruction & Engineering Education (hub) | `references/human-performance-technology.md`, `references/instructional-design-course-architecture.md`, `references/assessment-certification-design.md`, `references/technical-training-delivery.md`, … |