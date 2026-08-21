<!-- hub-reference-banner -->
> **Reference file — part of the `technical-instruction` hub.** Formerly the standalone `cognitive-task-analysis` skill.
> Sibling topics in this family are now reference files under the hubs (`technical-instruction`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: cognitive-task-analysis
description: >-
  Elicit how experts actually think and decide using Cognitive Task Analysis (CTA) — the
  front-end knowledge-elicitation methodology that feeds instructional design. Covers CDM,
  ACTA (Militello & Hutton), PARI, GDTA (Endsley), think-aloud protocol, concept mapping.
  Addresses expert blind-spot problem (experts omit ~70% of decision steps unprompted), HTA
  vs CTA selection, Cognitive Demands Tables, Decision Requirements Tables.
  TRIGGER: cognitive task analysis; CDM or ACTA interview; capture tacit knowledge from SME;
  expert blind spot; HTA vs CTA; cognitive demands table; decision requirements table;
  front-end analysis for training; turn SME interviews into training; GDTA.
  SKIP: course/curriculum design -> instructional-design-course-architecture; teaching
  troubleshooting pedagogy -> teaching-troubleshooting-diagnostic-reasoning; software
  requirements -> software-engineering-patterns; psychometric exam ->
  assessment-certification-design. Family hub -> technical-instruction.

version: "1.1.0"
updated: "2026-06-16"
category: custom
tags:
  - cognitive-task-analysis
  - knowledge-elicitation
  - CTA
  - CDM
  - ACTA
  - tacit-knowledge
  - expert-performance
  - training-design
  - naturalistic-decision-making
  - instructional-design
related_skills:
  - instructional-design-course-architecture
  - teaching-troubleshooting-diagnostic-reasoning
  - assessment-certification-design
  - applied-psychology
  - technical-instruction
whenToUse:
  - "run a CDM or ACTA interview with a subject-matter expert"
  - "elicit tacit knowledge from an expert for training design"
  - "build a cognitive demands table or decision requirements table"
  - "decide whether HTA or CTA is the right front-end analysis method"
  - "address the expert blind-spot problem in knowledge capture"
  - "capture how experts notice cues and make decisions under uncertainty"
  - "turn SME interview transcripts into structured training inputs"
  - "conduct a knowledge audit using ACTA probe categories"
  - "apply GDTA to elicit situation-awareness requirements"
---

# Cognitive Task Analysis & Knowledge Elicitation

> Full context: `references/cognitive-task-analysis.md`

## What this skill covers

Cognitive Task Analysis (CTA) is the systematic methodology for capturing the unobservable cognitive processes — perceptual cues, mental models, decision heuristics, situation-assessment patterns — that expert performers use but cannot fully articulate unprompted. It is the analytic front end that sits upstream of course design (instructional-design-course-architecture) and of diagnostic-training pedagogy (teaching-troubleshooting-diagnostic-reasoning).

## Method selection decision table

| Scenario | Best method | Time cost |
|---|---|---|
| Retrospective incident, rich narrative data wanted | CDM (Critical Decision Method) | Medium–High |
| Practitioner is a non-researcher; needs fast results | ACTA (4-stage toolkit) | Medium |
| Procedural / troubleshooting task, paired experts | PARI | Medium |
| Situation-awareness and information requirements | GDTA (Endsley) | Medium |
| Real-time task capture, no recall distortion | Think-aloud (concurrent protocol) | Low–Medium |
| High-level mapping of expert knowledge structure | Concept mapping / knowledge audit | Low |
| Observable behavior is sufficient, cognition secondary | HTA (Hierarchical Task Analysis) | Low |

## Core principles

1. **The 70% problem.** In unaided verbal accounts, experts omit up to 70% of decision steps. Structured probing techniques exist precisely to surface what automaticity hides.
2. **Incidents over abstractions.** All major CTA methods anchor on specific events (CDM) or plausible scenarios (ACTA simulation interview) to bypass the expert's abstract, idealized self-description.
3. **Multiple passes, multiple experts.** Single-pass, single-expert interviews consistently underrepresent the knowledge space. Plan for 3-5 experts and at least 2 interview passes per method.
4. **Outputs must be actionable.** Raw transcripts are not CTA. The deliverable is a structured knowledge product (Cognitive Demands Table, Decision Requirements Table, or scenario set) that an instructional designer can directly consume.
5. **HTA is not sufficient for cognitive work.** HTA captures observable sequences; CTA captures the judgments, perceptual cues, and mental simulations that make those sequences succeed or fail under uncertainty.

## Quick guides

### CDM in 4 passes
1. **Incident selection** — identify a challenging, non-routine event the expert handled
2. **Timeline** — reconstruct the event chronologically without probing
3. **Decision-point probe sweep** — for each decision: What cues did you notice? What did you expect to happen? What alternative did you consider? What would a novice have missed?
4. **"What if" deepening** — introduce hypothetical variations to reveal dependencies and exception knowledge

### ACTA in 4 stages
1. **Task diagram** — decompose work into 3-6 major phases; flag cognitively demanding ones
2. **Knowledge audit** — probe each flagged phase using expertise characteristics: Big Picture, Pattern Recognition, Past & Future, Job Smarts, Improvisation, Self-Monitoring
3. **Simulation interview** — present a challenging scenario; elicit situational assessment, cue interpretation, and projected novice errors at each decision point
4. **Cognitive Demands Table** — synthesize all findings into the structured output

### Cognitive Demands Table columns
| Task step | Difficult cues | Why difficult | Errors/misinterpretations | Expert strategy | Training implication |
|---|---|---|---|---|---|

### Decision Requirements Table columns
| Decision point | Key cues (what triggers action) | Relevant knowledge/rules | Typical novice error | Critical uncertainty |
|---|---|---|---|---|

## Boundary: CTA vs adjacent methods

- **CTA elicits** the expertise. `teaching-troubleshooting-diagnostic-reasoning` **teaches** it. These two skills form a sequential pair; CTA comes first.
- **Instructional-design-course-architecture** consumes CTA outputs as front-end analysis. Route there once the Cognitive Demands Table is complete.
- General qualitative interviewing (semi-structured interviews, focus groups, ethnography) is broader and task-neutral. CTA is specifically task- and decision-anchored.

## Effect-size evidence

Meta-analysis of 12 surgical training studies (Stefanidis et al.): CTA-based training vs conventional — procedural knowledge SMD 1.36 (large), technical performance SMD 2.06 (large). 92.3% of studies showed CTA-based training improved surgical outcomes.
