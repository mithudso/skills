<!-- hub-reference-banner -->
> **Reference file — part of the `technical-instruction` hub.** Formerly the standalone `genai-education-instructional-design` skill.
> Sibling topics in this family are now reference files under the hubs (`technical-instruction`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: genai-education-instructional-design
description: >-
  Apply and evaluate GenAI for education (education layer, not LLM engineering). AI-assisted ID (ADDIE/SAM+AI, HITL; Articulate/Coursebox/Synthesia); LLM tutors/ITS (Khanmigo, LearnLM, Socratic/guardrailed, Bloom 2-sigma); AIG quality/automation bias; GenAI in CS (amplifier-vs-equalizer, cognitive offloading); hallucination, FERPA/COPPA, equity, AI-resistant assessments.
  TRIGGER: AI-assisted course design (AI required); AI tutors/ITS; Bloom 2-sigma; Socratic tutoring; ADDIE/SAM+AI; AIG quality; AI-resistant assessment (learning); GenAI in CS; cognitive offloading; FERPA/COPPA; AI equity.
  SKIP: LLM/agent engineering or RAG/prompt layer->ai-agents-orchestration,ai-llm-model-layer,ai-mcp-sdk-prompting; human-AI trust (no learning angle)->applied-psychology; cert psychometrics/AI-resistant cert->assessment-certification-design; ID without AI->instructional-design-course-architecture; delivery->technical-training-delivery; troubleshooting->teaching-troubleshooting-diagnostic-reasoning; family hub->technical-instruction.
version: "1.2.0"
updated: "2026-06-16"
category: custom
whenToUse:
  - "evaluating or choosing AI tools for course authoring (Articulate AI, Coursebox, Synthesia, Mindsmith)"
  - "designing a human-in-the-loop instructional design workflow with GenAI"
  - "assessing whether AI-generated learning objectives or quiz items meet quality standards"
  - "building or evaluating an AI tutor or ITS (Khanmigo, LearnLM, custom LLM tutor)"
  - "understanding the Bloom 2-sigma problem and what AI tutoring can realistically achieve"
  - "implementing Socratic or guardrailed LLM tutoring patterns"
  - "finding RCT or meta-analysis evidence on how well AI tutors improve learning outcomes"
  - "assessing GenAI impact on novice vs. expert learners in CS/programming courses"
  - "evaluating cognitive offloading risks from AI coding assistants"
  - "designing AI-resistant assessments or responding to AI-enabled academic integrity concerns"
  - "auditing data privacy compliance (FERPA, COPPA) for educational AI tool deployments"
  - "assessing equity implications of AI-powered learning tools"
keywords:
  - AI tutors
  - intelligent tutoring systems
  - ITS
  - Khanmigo
  - LearnLM
  - Bloom 2-sigma
  - generative AI instructional design
  - AI course authoring
  - automatic item generation
  - AIG
  - adaptive learning
  - AI personalized learning
  - AI assessment integrity
  - AI-resistant assessment
  - GenAI CS education
  - cognitive offloading
  - FERPA AI tools
  - pedagogical agents
  - Socratic tutoring LLM
  - ADDIE GenAI
  - educational AI hallucination
  - equity digital divide AI education
tags:
  - education
  - instructional-design
  - AI-tutoring
  - assessment
  - genai
  - cs-education
  - privacy
related_skills:
  - instructional-design-course-architecture
  - assessment-certification-design
  - technical-training-delivery
  - teaching-troubleshooting-diagnostic-reasoning
  - applied-psychology
  - ai-agents-orchestration
metadata:
  changelog:
    - "2026-06-16 sko v1.1.0->v1.2.0: Iter 2 — 3 High, 12 Medium fixed; description rewritten to 971 chars; source confidence labels corrected ([^6],[^13],[^15]); ARCHED tentative label added; VanLehn decision-table attribution clarified; meta-analysis qualifier actionable; explainability defined; Bastani dedup; corporate L&D expanded; em-dash density reduced; AI-resistant collision SKIP tightened; whenToUse abstract entry rewritten"
    - "2026-06-16 sko v1.0.0->v1.1.0: Iter 1 — 8 High, 19 Medium fixed; references extracted to references/genai-education-bibliography.md; description truncated to <1000 chars; category fixed; related_skills corrected"
---

# GenAI for Instructional Design & AI Tutors

**Domain:** Educational applications of generative AI (2024-2026) — lens is EDUCATION, not LLM engineering.
**Verified-as-of:** 2026-06-16

> **Scope limits.** This skill covers what AI does *for learners and designers*, not how AI works internally. For LLM/agent architecture, RAG, prompting technique, or model training, see the `ai-*` skill family. For ADDIE/SAM/course design with no AI component, see `instructional-design-course-architecture`. For psychometric mechanics (IRT, Angoff, DIF), see `assessment-certification-design`. For human trust calibration and cognitive bias in AI adoption (without a learning-outcome angle), see `applied-psychology`.

**Evidence confidence key** used throughout this skill:
- **Fact** — 3+ independent sources agree; treat as established finding.
- **Qualified** — 2 sources, or 1 strong RCT with known limits; use with stated caveats.
- **Tentative** — single study or preprint; directional signal only.

**Citations** are in `references/genai-education-bibliography.md`; inline `[^N]` markers map to that file.

---

## Contents

1. [AI-Assisted Instructional Design](#1-ai-assisted-instructional-design)
2. [AI Tutors & Intelligent Tutoring Systems](#2-ai-tutors--intelligent-tutoring-systems)
3. [AI for Assessment, Item Generation & Integrity](#3-ai-for-assessment-item-generation--integrity)
4. [Adaptive & Personalized Learning](#4-adaptive--personalized-learning)
5. [GenAI in Developer/CS Education](#5-genai-in-developercs-education)
6. [Risks, Guardrails & Governance](#6-risks-guardrails--governance)
7. [Decision Tables](#7-decision-tables)
8. [Anti-Patterns](#8-anti-patterns)

---

## 1. AI-Assisted Instructional Design

### Adoption Overview (verified-as-of: 2026-06-16)

By late 2024, 84% of instructional designers reported using ChatGPT in their work, up from ~5% a year earlier (practitioner survey, ~500 IDs).[^1] ChatGPT appeared in 26 of 35 peer-reviewed studies in a 2025 CITE Journal systematic review.[^2] Named commercial tools with AI course-authoring features include:

- **Articulate AI Assist** — generates outlines, lesson drafts, quiz questions, images, audio, and branching simulations within Storyline/Rise
- **Coursebox** — lesson, quiz, and video generation from source documents
- **Synthesia** — AI-avatar training videos from prompts or URLs, with quiz generation
- **Mindsmith** — AI-first storyboard-to-interactive eLearning with adaptive assessments and gamified scenarios
- **iSpring AI** — quiz generation, translation, and image integration
- **ThingLink Scenario Builder** (Fall 2024) — immersive virtual tours and training scenarios[^3][^4][^5]

Practitioner comparisons suggest GPT-4o is preferred for structured curriculum design and prerequisite mapping, while Claude Sonnet is preferred for longer-form content; these reports are from practitioner blogs, not controlled studies, so treat as anecdotal guidance only.[^6]

### Human-in-the-Loop (HITL) Workflows

HITL is the dominant recommended model across multiple sources.[^2][^7][^8] AI acts as a drafting co-pilot and independent reviewer, not a replacement for instructional judgment. Five-stage lifecycle:

1. **Strategy & Analysis** — human defines learner needs, context, and constraints; AI synthesizes survey data and skills-gap assessments
2. **AI-Assisted Drafting** — LLM generates objectives, outlines, scenario branches, quiz items
3. **SME/ID Expert Refinement** — human edits for accuracy, depth, organizational context, and learner appropriateness
4. **Governance Review** — copyright, privacy, accessibility, and equity review (see Section 6 for the full checklist)
5. **Continuous Feedback Loop** — learner performance data feeds revision cycles[^7]

A key bottleneck: without advanced prompting skill, AI outputs are shallow and generic. Workers skilled in prompt engineering complete ID projects up to 40% faster at higher quality; the skill gap means benefits accrue unevenly.[^1]

### ADDIE/SAM Augmentation (AI angle required)

| Phase | AI Contribution |
|---|---|
| Analysis | Survey summarization, skills-gap identification, needs assessment synthesis |
| Design | Learning-objective generation, prerequisite mapping, strategy selection |
| Development | Script generation, multimedia production, rapid prototyping, quiz generation |
| Implementation | Onboarding materials, communications drafting |
| Evaluation | Survey design, learner performance pattern analysis, dropout risk identification |

SAM's iterative design-develop-evaluate cycles map especially naturally to AI-assisted rapid prototyping; tools like Mindsmith operationalize this with fast storyboard iteration.[^9]

**ARCHED Framework** (AAAI 2025 workshop preprint): specialized AI agents in sequence — one generates pedagogically diverse options, another evaluates Bloom's-taxonomy alignment — with humans retaining final authority. Blind evaluation: ARCHED outputs achieved mean expert rating of 4.43/5; F1-score 91.8% for Bloom's level classification (Tentative — single preprint).[^10]

### Quality Concerns

- **Generic and shallow outputs** — the most frequently reported quality defect across multiple peer-reviewed studies; AI without learner context defaults to content-first rather than learner-centered design.[^2][^11]
- **Pedagogical misalignment** — most AI errors in teaching plans are instructional (wrong complexity level, mismatched detail), not factual.[^2]
- **Hallucination rate** — an estimated 53% of AI output errors in one educational case study were hallucinations.[^12]
- **Copyright and IP risks** — AI-generated content can inadvertently replicate copyrighted materials; governance frameworks such as 1EdTech AI-Generated Content Best Practices v1.0 recommend SME validation for all publisher-produced AI content.[^11]
- For anti-patterns on treating AI output as final, see Section 8.

---

## 2. AI Tutors & Intelligent Tutoring Systems

### The Bloom 2-Sigma Problem

Bloom (1984) found one-on-one expert tutoring combined with mastery learning raised student performance by ~2 standard deviations vs. conventional instruction — the average tutored student outperformed 98% of traditionally instructed peers.[^13] This remained economically unscalable for 40 years; AI is now the first serious candidate to address it at scale.[^14]

Note on [^13]: the bibliography entry cites a secondary blog source for Bloom 1984. The original peer-reviewed source is: Bloom, B. S. (1984). The 2 sigma problem: The search for methods of group instruction as effective as one-to-one tutoring. *Educational Researcher*, 13(6), 4–16. The 2-sigma finding is well-established across educational research; the direction of this claim is Fact-level, but readers should verify against the primary article.

VanLehn's 2011 meta-analysis (54 comparisons, 28 studies) showed pre-LLM step-based ITS achieved d=0.76 vs. no tutoring — nearly matching human tutoring (d=0.79), suggesting the practical human-tutor advantage is closer to 0.79 sigma, not 2.0.[^15] The d=0.76 figure is from the original peer-reviewed 2011 *Educational Psychologist* article; [^15] in the bibliography notes it was accessed via a 2015 secondary review. A 2025 meta-analysis of ITS across U.S. K-12 contexts found a significant but modest effect size of g=0.271.[^16]

### Major LLM-Based ITS (verified-as-of: 2026-06-16)

- **Khanmigo** (Khan Academy): deployed in 40+ districts with 28,000+ students. Multiple 2024 studies found no statistically significant learning outcome difference vs. traditional methods in some contexts; a qualitative evaluation found it did not reliably tailor tasks to individual ability or support metacognitive development.[^17][^18][^19]
- **LearnLM** (Google DeepMind): Gemini-family model fine-tuned for pedagogical quality. A 2025 UK RCT (N=165, 5 secondary schools) found students guided by LearnLM were 5.5 percentage points more likely to solve novel problems than students with human tutors alone; human tutors approved 76.4% of LearnLM drafts with zero or minimal edits.[^20][^21]
- **MATHia / Carnegie Learning**: extensively studied cognitive-tutor/LLM hybrid, deployed in thousands of K-12 schools.[^22]
- **GPT-4-based ITS feedback**: ~80% accuracy diagnosing student errors for College Algebra; degrades when responses contain multiple simultaneous errors.[^23]

### Evidence Summary

Confidence labels use the key defined at the top of this skill.

| Study | Design | Result | Confidence |
|---|---|---|---|
| Kestin et al. 2025, *Scientific Reports*, Harvard physics, N=194 | RCT | ~2x learning gains vs. active learning | Qualified — single elite university, 2-week duration, lower-order skills only [^24][^25] |
| LearnLM UK RCT (arXiv:2512.23633), N=165 | RCT | +5.5pp on novel problems vs. human tutoring alone | Qualified — preprint, Google-authored [^20][^21] |
| VanLehn 2011 (54 comparisons) | Meta-analysis | d=0.76 for step-based ITS vs. no tutoring | Fact — peer-reviewed [^15] |
| K-12 ITS meta-analysis arXiv:2511.04997 | Meta-analysis | g=0.271 positive effect | Qualified — preprint [^16] |
| Bastani et al. 2024, SSRN/PNAS-submitted | RCT | Unrestricted AI: +48% practice, -17% exam; guardrailed GPT tutor: on par with or above control | Fact — consistent with other guardrailed-vs-unguarded comparisons [^26] |
| PLOS ONE 2024 ChatGPT-as-tutor, RCT | RCT | Learning gains equivalent to human tutor-authored help | Fact — peer-reviewed [^27] |
| Brookings 2024 synthesis | Review | Cautious positive; effect sizes 0.2-0.6 in stronger studies | Qualified [^28] |

**Key nuance:** The Bastani et al. pattern is the most operationally important. Unrestricted AI access harms outcomes; guardrailed, Socratic AI tutoring delivers gains consistent with the ITS meta-analytic literature. The tool is not the problem; unconstrained use is. See Section 5 (Disconfirming Evidence) for the same finding in CS/developer education contexts.

**Interpreting mixed evidence:** When studies conflict (e.g., Khanmigo showing no gains vs. LearnLM showing gains), check: (1) study population and duration — short RCTs at elite institutions do not generalize; (2) whether guardrails were present; (3) what outcome was measured — procedural performance vs. conceptual understanding diverge reliably.

### Socratic & Guardrailed Tutoring Patterns

LLMs default to answer-giving without architectural enforcement.[^29] Key patterns:

- **Separate system prompt personas**: "hint giver," "question poser," "answer validator" — each with distinct instructions
- **Finite-state slot structure** (MWPTutor, 2024): pump moves → Socratic hints → increasingly directive prompts → answer as last resort[^29]
- **RAG-based course grounding**: constrain tutor responses to verified course material to reduce hallucination; RAG-constrained tutoring improved outcomes vs. unconstrained LLM in a 2025 study[^30]
- **Daily usage caps and metacognitive reflection prompts**: structured integration (compare → reflect → revisit) significantly reduces cognitive offloading vs. unguided access[^31]
- **Explicit fallibility disclosure**: warning students about AI fallibility measurably increases help-seeking from human instructors[^32]

### Pedagogical Agents: ITS Lineage

| Era | Key Development |
|---|---|
| Pre-LLM (1960s-2010s) | Computer-assisted instruction → hand-crafted cognitive models → ML+NLP integration (ALEKS, Carnegie Cognitive Tutor); AutoTutor's 17-year natural-language tutoring advantage documented[^22] |
| LLM era (2020s-) | Open-ended tutoring without domain-specific engineering; challenge reverses — must now *constrain* LLM flexibility into pedagogically sound behavior, not expand limited systems[^29] |

---

## 3. AI for Assessment, Item Generation & Integrity

### Automatic Item Generation (AIG) Quality

A systematic review of 60 papers using LLMs for AIG found **psychometric evaluations of generated items are absent in most papers** — quality claims rest on surface-level human ratings, not psychometric evidence.[^33][^34] This is the field's largest credibility gap.

Key documented issues:

- **Automation bias degrades item quality**: teachers reviewing AI-generated MCQs accepted flawed items at higher rates than self-authored items; a 2026 Frontiers peer-reviewed study found AI-assisted MCQ creation increased item-writing flaws across multiple workflow configurations.[^35]
- **Lexical overlap cueing bias**: stem wording telegraphs the correct answer, undermining discriminative validity.[^35]
- **Weaker quality for complex traits**: LLM-generated SJTs showed satisfactory reliability for behaviorally concrete traits but weaker psychometric quality for Honesty-Humility, Agreeableness, and Openness.[^33]
- **Hallucination in assessment content**: AI tools produce plausible but factually incorrect items with confident presentation; students often cannot distinguish hallucinated facts from correct ones.[^36]

For psychometric evaluation mechanics (IRT, Angoff, KR-20, DIF, validity studies), defer to `assessment-certification-design`.

### Assessment Integrity: AI-Enabled Cheating

- 88% of students (HEPI survey, 2025) reported using generative AI tools for assessments.[^37]
- AI detection tools improved from ~40% effectiveness (2020) to ~70% (2024) — still below the threshold for reliable enforcement.[^37]
- No single countermeasure is sufficient; layer redesign + process documentation + oral defense.[^38]

### AI-Resistant Assessment Formats

| Format | Evidence Level |
|---|---|
| Oral exams / live follow-up after written submission | Strong — multiple practitioner and research sources [^38][^39] |
| Audience-tailored assessments (specific local context) | Strong [^38] |
| Observational assessments of live work practice | Strong [^38] |
| Reflection on live events (role-plays, simulations) | Strong [^38] |
| Debate/panel participation | Moderate [^38] |
| Portfolio with process documentation | Moderate [^39] |
| Timed in-person assessments | Moderate — accessibility concerns [^39] |

---

## 4. Adaptive & Personalized Learning

### Platform Overview (verified-as-of: 2026-06-16)

AI-driven adaptive platforms use ML to model learner state, adjust content difficulty, sequence retrieval practice, and optimize review timing:

- **Duolingo Max**: ML-driven difficulty adjustment + LLM conversation ("Lilly" AI). Peer-reviewed studies report improvements across skills after 3 months; Duolingo's own research team authors many efficacy studies — independent replication is limited.[^40]
- **Century Tech**: adopted in 55+ countries, 10,000+ students; rigorous large-scale peer-reviewed efficacy studies are limited outside vendor-reported data.[^41]
- **ALEKS**: knowledge-space-theory adaptive assessment and learning; most extensively studied adaptive platform in higher education math.

A meta-analysis of 25 studies found adaptive learning technologies consistently enhance performance, with learner performance increasing in 59% of studies (Qualified — meta-analysis spans heterogeneous platforms and outcome measures; apply this finding as directional support, not a precise effect estimate; weight results from studies using validated outcome measures more heavily).[^42]

### Spaced Repetition & Mastery Progression

AI-enhanced spaced repetition predicts individual memory decay curves to personalize review timing. The SSP-MMC algorithm (trained on 220M memory behavior records) can reduce unnecessary reviews by 15-20% and improve long-term retention by ~10-15%.[^43]

**LECTOR** (2025 arXiv): LLM-enhanced concept-based adaptive spaced learning integrating concept mapping with LLM-generated practice items (tentative — single paper).[^44]

**Key design gap**: adaptive systems generate rich learner data but fail to translate diagnostic evidence into appropriate instructional actions reliably — the "diagnostic-pedagogical loop" problem.[^45]

**Corporate L&D context**: Corporate training differs from K-12/higher-ed in several ways that affect how AI tools apply. FERPA does not apply; instead, employee data is governed by employment contracts, state privacy law (CCPA, VCDPA, etc.), and vendor data-processing agreements. IP ownership of AI-generated training content typically vests in the employer under work-for-hire doctrine, but verify contract terms with each AI vendor. LRS/xAPI integration is standard in enterprise L&D stacks; AI-driven adaptive platforms (e.g., Docebo, Cornerstone, SAP SuccessFactors LMS with AI add-ons) ingest xAPI statements to personalize learning paths. Regulatory compliance training (OSHA, HIPAA, financial services) has documentation requirements that AI-generated content must satisfy — always validate AI-produced compliance content against the current regulatory text. Vendor efficacy claims for enterprise AI learning tools should be treated as Qualified at best; independent peer-reviewed studies are sparse as of 2026. For xAPI/LRS architecture and learning measurement in corporate contexts, see `learning-measurement-evaluation`.

---

## 5. GenAI in Developer/CS Education

### The Amplifier-vs-Equalizer Debate

The dominant research finding: GenAI in CS education acts as an **amplifier of existing advantage**, not an equalizer.[^46][^47][^48]

**The Widening Gap** (Lau et al., 2024, ACM ICER, 21 lab sessions with novice programmers): strong novice programmers used GenAI to produce intended code and filter bad suggestions; weak programmers experienced compounded metacognitive failures and developed false confidence. All 20 who attempted with AI completed the task, but completion masked the quality divide.[^46]

A 2024 systematic review (58 studies, 2022-2025): 94.83% of studies reported enhanced programming support from AI tools, but 65.52% also identified overreliance leading to superficial learning.[^47]

**Learner context note**: Evidence covers both CS students (Lau et al., Bastani et al.) and professional software developers (Anthropic RCT). Effects are directionally consistent but magnitude varies by prior expertise and task type.

### Cognitive Offloading Evidence

**Anthropic RCT (2026, n=52 software developers)**: AI users averaged 50% on comprehension tests vs. 67% for manual coders — approximately a two-letter-grade gap — with the deficit most pronounced on debugging tasks. High performers who used AI strategically (asking for explanations, using AI for conceptual questions only) were not penalized.[^49]

The mechanism: LLMs guide students through problems step-by-step so "each step was able to be rationalized as understanding," creating an **illusion of competence** that makes it harder to detect one's own knowledge gaps.[^46]

**Mitigation**: structured GenAI integration with metacognitive scaffolding (compare → reflect → revisit) significantly reduces cognitive offloading vs. unguided access (single 2026 Frontiers study).[^31]

### Disconfirming Evidence

- Codex/Copilot-assisted learners in a 2023 study did NOT show retention loss; learners with higher prior scores retained significantly more. Harm is not universal and appears tool- and task-specific.[^50]
- Bastani et al. "GPT Tutor" condition: structured, Socratic guardrails reversed the cognitive-offloading harm — students performed on par with or above the control group. (Full Bastani evidence in Section 2 Evidence Summary table.)[^26]

---

## 6. Risks, Guardrails & Governance

### Hallucination in Tutoring Contexts

LLMs hallucinate with confident, authoritative language that is difficult for students to detect without domain expertise. Hallucinated tutoring explanations can become embedded misconceptions.[^51]

A 2026 arXiv thematic analysis found >50% of student hallucination detection attempts relied on intuition, because AI content carries an "illusion of truthfulness."[^36]

- **Technical mitigation**: RAG-based architecture grounding responses in verified course material.[^30]
- **Behavioral mitigation**: explicitly warn students about AI fallibility — measurably increases help-seeking from human instructors.[^32]

### FERPA & Data Privacy (verified-as-of: 2026-06-16)

- 42% of US school districts using AI tools lack a Data Processing Agreement (DPA) with the vendor — a de facto FERPA compliance gap regardless of vendor security practices (CDT 2024 data).[^52]
- Faculty submitting student work to AI tools for grading may violate FERPA because the AI vendor becomes an unauthorized handler of education records.[^53]
- FTC finalized COPPA amendments in January 2025 shifting to opt-in consent for data collection from under-13s — directly affects K-12 AI tool deployments.[^54]
- FERPA's text predates cloud storage and has no cybersecurity requirements, creating an inherent compliance gap.[^55]

**Governance checklist** (must-have first):

1. **DPA required** — vendor Data Processing Agreement reviewed and signed before any student data flows to the tool
2. **Explicit student consent** — documented before using education records with AI
3. **Vendor data-use prohibition** — contract language prohibiting vendor from using student data to train commercial models
4. **Data minimization** — process only what is pedagogically necessary
5. **Periodic audits** — of vendor data practices and access logs, at least annually
6. **AI explainability for grading** — vendor must provide per-student explanation for any AI-generated grade or feedback: minimally, a rubric-aligned rationale (e.g., "score of 3/5 because criterion X was partially met"); a black-box score with no rationale does not satisfy this requirement

### Equity & Digital Divide

GenAI amplifies existing advantages:[^48][^56][^57]

- Well-resourced institutions benefit most — enterprise AI access, trained instructors, governance staff
- Community colleges (disproportionately serving underserved populations) typically cannot afford enterprise contracts
- AI systems trained on biased corpora underrepresent certain demographics; algorithmic bias documented in educational feedback tools
- Device and connectivity gaps remain the primary bottleneck for underserved students even when free-tier AI is available
- 50% of students report feeling less connected to teachers when AI is used in instruction[^56]

The "AI as equalizer" narrative has surface plausibility (free tiers, 24/7 availability) but breaks down in practice due to infrastructure gaps, language barriers, and algorithmic bias.[^48]

---

## 7. Decision Tables

### When to Use AI in Instructional Design

| Task | Use AI? | Caveat |
|---|---|---|
| First draft of learning objectives | Yes | Requires ID review against Bloom's levels and learner context |
| Final learning objectives | No — human review required | AI defaults to generic; organizational context missing |
| Module outline structure | Yes | Strong starting point; SAM iteration recommended |
| Branching scenario scripts | Yes | Effective with Articulate/Mindsmith; ID validates |
| Quiz item generation | Yes with caution | Automation bias risk; psychometric review required |
| Final psychometric validation | No | Requires human item analysis — defer to assessment-certification-design |
| Learner needs analysis | Assist only | AI can synthesize survey data; ID interprets |
| Copyright-sensitive content | No | Legal review required for AI-generated content |

### Choosing an AI Tutor Approach

| Scenario | Recommended Approach |
|---|---|
| Math/STEM procedural skills | Step-based ITS with Socratic hints (evidence: pre-LLM ITS meta-analysis d=0.76 for step-based systems vs. no tutoring[^15]; Socratic-hints pattern is the modern LLM implementation of that approach) |
| Open-ended conceptual learning | Guardrailed LLM tutor with RAG grounding and usage caps |
| Novice programmers | Structured scaffolding + metacognitive reflection; unguided access amplifies disadvantage |
| Language learning | LLM conversation agents (Duolingo Max model); spaced repetition layer |
| High-stakes assessment prep | AI formative feedback + human instructor validation at key checkpoints |
| Resource-constrained context | Verify device/connectivity access before deploying AI tutors |

### AI-Resistant Assessment Design

| Constraint | Best Countermeasure |
|---|---|
| Remote, unproctored exam | Audience-tailored + process documentation requirement |
| Written assignment | Add live oral follow-up; require reflection on specific classroom event |
| Programming assignment | Oral code walk-through; explain design decisions live |
| All exam types | Layer at least 2 of: redesign + process documentation + oral defense |

---

## 8. Anti-Patterns

- **Deploying an LLM tutor without guardrails**: unrestricted AI access consistently harms novice learners — false confidence, illusion of competence, exam score decline. Always implement Socratic constraints and usage caps.[^26][^46]
- **Treating AI-generated MCQs as ready-to-use**: automation bias causes teachers to over-accept flawed AI items. Every AI-generated item needs psychometric review.[^35]
- **Assuming AI equalizes access**: the amplifier-not-equalizer finding is consistent across CS education literature. Equity-focused deployment requires explicit infrastructure and training investment.[^46][^48]
- **FERPA compliance assumed from vendor claims**: 42% of districts lack DPAs. Verify independently; faculty tool use creates FERPA obligations regardless of institutional contracts.[^52][^53]
- **Using AI formative feedback for lower-ability learners without instructor review**: AI feedback is vague and impersonal for the students who most need precise guidance.[^58]
- **Treating AI-generated course content as final**: 1EdTech recommends SME validation for all publisher-produced AI content; human review is not optional.[^11]
- **Extrapolating from a single strong RCT**: the Harvard 2025 result (2x gains) has at least 6 identified methodological limits — 2 weeks, 1 course, 1 elite institution, lower-order skills, instructor-designed, small N. Always cite effect sizes alongside their constraints.[^24][^25]
- **Bypassing human review to speed delivery**: a 2024 practitioner survey (n~500) found IDs still handled only 2-4 concurrent projects post-AI adoption — AI had not materially expanded throughput, suggesting efficiency gains are offset by necessary review overhead.[^1]

---

*Full bibliography: [references/genai-education-bibliography.md](references/genai-education-bibliography.md)*
