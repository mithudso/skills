<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `software-architect` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: software-architect
description: >
  Software architecture methodology and documentation expert — ISO/IEC/IEEE 42010 concepts,
  arc42 template structure, C4 model diagrams, SEI ATAM evaluation, ADR writing, quality
  attributes and tradeoff analysis, and AWS/Azure/GCP well-architected reviews.
  TRIGGER: designing system architecture, writing ADRs, choosing C4 diagram levels,
  running an ATAM or well-architected review, documenting architecture with arc42, defining
  quality scenarios, decomposing system boundaries, modernizing or evolving an existing architecture.
  SKIP: AWS/GCP/Azure service configuration or deployment tasks (use aws-core or cloud-specific skills);
  MongoDB schema design (use mongodb-schema-design); API design (use api-design-patterns);
  CI/CD pipelines (use cicd-pipelines).
version: 1.2.0
last_updated: 2026-05-29
category: developer
tags:
  - software-architecture
  - arc42
  - c4-model
  - ATAM
  - ADR
  - quality-attributes
  - well-architected
  - system-design
  - architecture-documentation
triggers:
  - design system architecture
  - architecture documentation
  - ADR architecture decision record
  - C4 diagram system context container component
  - arc42 template sections
  - ATAM evaluation
  - well-architected review
  - quality attributes tradeoffs
  - architecture boundaries interfaces
  - decompose a system
  - modernize architecture
  - architecture review
  - technical debt architectural
  - stakeholder architecture concerns
related_skills:
  - aws-core
  - api-design-patterns
  - mongodb-schema-design
  - microservices-patterns
---

# Software Architect

Expert reference for software architecture methodology: system design, documentation (arc42, C4), evaluation (SEI ATAM), decision recording (ADRs), quality-attribute analysis, and well-architected cloud reviews.

## When to use this skill

- Designing or reviewing system architecture at any level (context, container, component)
- Writing or reviewing Architecture Decision Records (ADRs)
- Choosing the right C4 diagram level for the audience
- Running or preparing for a well-architected review (AWS, Azure, or GCP)
- Structuring an architecture document with arc42 sections
- Defining quality scenarios and analyzing tradeoffs with SEI ATAM
- Decomposing system boundaries and external interfaces
- Planning architecture modernization or managing technical debt

## When NOT to use this skill

- AWS/GCP/Azure service configuration, IAM, or deployment tasks → use `aws-core` or cloud-specific skills
- MongoDB schema or index design → use `mongodb-schema-design`
- REST/gRPC API design → use `api-design-patterns`
- CI/CD pipeline design → use `cicd-pipelines`
- Code-level design patterns (DI, CQRS, event sourcing at implementation level) → use `microservices-patterns` or `backend-patterns`

---

## How to apply this skill

Match the response format to the question type:

| Question type | Response shape | Success signal |
|--------------|----------------|----------------|
| "Help me design X" | Walk the 8-step workflow; produce arc42 sections or C4 diagrams as appropriate for the audience | User has a concrete next deliverable (a diagram, a section, an ADR) |
| "Review this architecture" | Apply ATAM steps 4–8 (identify approaches, utility tree, scenario analysis); produce a risk/tradeoff table | Risks, non-risks, sensitivity points, and tradeoff points are named |
| "Write an ADR for decision X" | Produce a Nygard-style ADR: Title, Status, Context, Decision, Consequences (positive + negative + neutral) | ADR is specific enough that a future developer understands why the alternative was rejected |
| "Which framework should I use?" | Identify the question as documentation, evaluation, or cloud-review; recommend the right tool from the methods table; explain caveat | Recommendation is grounded in the cited frameworks, not improvised |
| "Explain [concept]" | Explain using the cited source; note if the concept has scope limits | Explanation cites the source and names any known ambiguity |

**Source discipline:** all guidance must be grounded in the cited frameworks (arc42, C4, ISO 42010, SEI ATAM, AWS/Azure/GCP well-architected). If a question falls outside what the cited sources cover, say so explicitly rather than extrapolating. Do not invent framework guidance.

**Self-check before responding:** re-read the response and confirm it addresses the specific question asked — not the general topic area.

**Clarifying question policy:** if the question is ambiguous about which workflow phase or framework applies, ask exactly one targeted question before proceeding (e.g., "Are you looking to document the current architecture or evaluate it against quality goals?").

---

## Quick architecture rules

1. Start from **stakeholders, concerns, goals, and constraints** — not from technology. Both ISO 42010 and arc42 frame architecture description around stakeholders first.
2. Document with **multiple views** — no single diagram serves every audience. Use views that answer specific stakeholder questions.
3. Make **boundaries and interfaces explicit early**. arc42 treats context/scope as critical; C4 starts with system context before zooming inward.
4. Prefer documentation that is **clear, useful, and maintained** over exhaustive-but-stale. Documentation volume is not a quality signal.
5. Record **architecturally significant decisions** with rationale and consequences so future teams understand why the system looks the way it does.
6. Make quality attributes **concrete with scenarios** and analyze tradeoffs explicitly. Vague quality goals ("the system must be fast") cannot be evaluated.
7. Use the **simplest architecture that meets current needs**. Resist over-engineering; design for change means easy to modify, not pre-built for every scenario.
8. Treat architecture review as an **ongoing constructive conversation**, not a one-time audit.

---

## Architecture workflow

Apply steps selectively based on the question — not all 8 steps are needed for every task.

1. **Capture goals, stakeholders, and constraints.** Requirements, quality goals, stakeholder expectations, and constraints before discussing structure. ([arc42](https://arc42.org/overview))
2. **Delimit the system and its context.** Business and technical context, communication partners, external interfaces. ([arc42 section 3](https://docs.arc42.org/section-3/))
3. **Choose a view set that fits the audience.** Match C4 diagram levels to the question; apply arc42 sections where the depth is justified. ([C4 diagrams](https://c4model.com/diagrams), [arc42 docs](https://docs.arc42.org/home/))
4. **Define decomposition and interactions.** Static structure, runtime behavior, deployment, and cross-cutting concepts. ([arc42](https://arc42.org/overview), [C4 abstractions](https://c4model.com/abstractions))
5. **Turn quality goals into measurable scenarios.** Use arc42 quality scenarios and SEI utility trees to make tradeoffs testable. ([arc42 section 10](https://docs.arc42.org/section-10/), [SEI ATAM](https://www.sei.cmu.edu/library/architecture-tradeoff-analysis-method-collection/))
6. **Evaluate the architecture.** ATAM-style scenario analysis: identify risks, sensitivity points, tradeoff points before implementation locks decisions in.
7. **Run a cloud-operational review when relevant.** Check AWS/Azure/GCP well-architected pillars for security, reliability, performance, operations, and cost.
8. **Record decisions, risks, and change history.** ADRs, technical debt, and documentation updates as the system evolves.

---

## Architecture description approaches

### ISO/IEC/IEEE 42010 concepts

ISO 42010 introduces a **conceptual model for architecture description**: architecture descriptions must be organized around viewpoints and explicitly state which conventions and frameworks they use. The practical takeaway: documentation should be *intentional about its viewpoints*, not an accidental pile of diagrams.

Note: ISO 42010:2011 is withdrawn; use it for conceptual grounding, not as a normative standard. ([overview](https://www.iso.org/standard/50508.html))

### arc42 documentation structure

arc42 is a pragmatic, open-source template for architecture communication organized into 12 sections: goals, constraints, context, solution strategy, building blocks, runtime, deployment, cross-cutting concepts, decisions, quality requirements, risks/technical debt, and glossary.

Key properties:
- **Systematic but flexible** — guidance is tagged as lean, thorough, or essential; use only what serves the audience.
- **Separates business context from technical context** — business context shows domain inputs/outputs; technical context shows channels, protocols, and hardware.
- **Section 9 (ADRs)** — recommends Nygard-style fields: Title, Context, Decision, Status, Consequences (positive, negative, and neutral — not just benefits).
- **Section 10 (quality requirements)** — distinguishes usage scenarios (runtime) from change scenarios (evolvability). ([arc42 overview](https://arc42.org/overview))

### C4 model

C4 is an **abstraction-first** approach using four abstraction levels: software systems, containers, components, and code. Most teams get sufficient value from system context + container diagrams alone — you do not need all four levels.

| Diagram | Shows | When to use |
|---------|-------|-------------|
| System context | System in its environment; users and neighboring systems | Always — first diagram to draw |
| Container | Applications, services, data stores inside the system | Architectural overview for developers |
| Component | Major components inside one container | Larger or complex containers only |
| Deployment | Infrastructure nodes and deployed elements | Operations and deployment review |

C4 is **notation independent** — diagrams must stand alone with titles, a legend, explicit element types, descriptions, and labeled relationships with direction and technology/protocol. ([C4 model](https://c4model.com/))

---

## SEI ATAM — Architecture evaluation

ATAM evaluates architectures **relative to quality attribute goals** and reveals how those goals interact and trade off. Its outputs are: risks, non-risks, sensitivity points, tradeoff points, and risk themes.

**Nine ATAM steps:**
1. Present the method
2. Present business drivers
3. Present the architecture
4. Identify architectural approaches
5. Generate a utility tree (quality factors → scenarios → priorities)
6. Analyze architectural approaches
7. Brainstorm and prioritize scenarios
8. Re-analyze against top scenarios
9. Present results

ATAM typically runs over three to four days with architects, a trained evaluation team, and stakeholder representatives. SEI notes that ATAM techniques can also be integrated **continuously into the design process** rather than saved for a late-stage review. ([SEI ATAM collection](https://www.sei.cmu.edu/library/architecture-tradeoff-analysis-method-collection/))

---

## Quality attributes and tradeoff analysis

- Quality requirements should be known in a **specific and measurable** way — "the system must handle 1,000 concurrent users with p99 latency < 200ms" not "the system must be fast."
- arc42 distinguishes **usage scenarios** (runtime behavior) from **change scenarios** (evolvability/change cost). Both must be covered.
- **Utility trees** structure quality-attribute work: decompose each quality goal into testable scenarios and rank by business priority.
- Tradeoffs are normal: Azure explicitly links every well-architected pillar to tradeoffs; ATAM is built to reveal how quality goals conflict rather than pretending they can all be optimized independently.

---

## ADR writing guide

An ADR (Architecture Decision Record) documents a single architecturally significant decision.

**When to write an ADR:**
- The decision is expensive to reverse
- Multiple reasonable alternatives existed
- The reasoning will not be obvious to future team members
- The decision affects other components or teams

**Nygard-style ADR fields:**

| Field | Purpose |
|-------|---------|
| Title | Short noun-phrase: "Use event sourcing for order history" |
| Status | Proposed / Accepted / Deprecated / Superseded |
| Context | Forces and constraints that led to this decision |
| Decision | The choice made, stated clearly |
| Consequences | Positive, negative, and neutral effects — not just benefits |

Store ADRs in version control alongside the code they govern. ([arc42 section 9](https://docs.arc42.org/section-9/))

**Example ADR (abbreviated):**

```
Title: Use PostgreSQL instead of DynamoDB for order history
Status: Accepted
Context: Order queries require multi-field filtering (customer, date range, status).
         DynamoDB's single-key model would require multiple GSIs and client-side joins.
Decision: Use PostgreSQL with a composite index on (customer_id, created_at).
Consequences:
  + Flexible queries without GSI proliferation
  + Schema migrations needed for structural changes (negative)
  + Team already has PostgreSQL expertise (neutral — no new training cost)
```

---

## Cloud well-architected reviews

All three major cloud providers frame well-architected reviews as **constructive improvement conversations**, not audits.

| Provider | Pillars | Review tool |
|----------|---------|-------------|
| **AWS** | Operational excellence, Security, Reliability, Performance efficiency, Cost optimization, Sustainability | AWS Well-Architected Tool |
| **Azure** | Reliability, Security, Cost optimization, Operational excellence, Performance efficiency | Azure Well-Architected Review + Azure Advisor |
| **Google Cloud** | Security, Reliability, Performance, Cost, Operations, Sustainability | Architecture Framework guidance |

**How to run a review:**
1. Select the provider framework matching the deployment target.
2. Work through pillar-by-pillar foundational questions.
3. Record high-risk findings with owners and remediation plans.
4. Treat results as a backlog, not a compliance checklist.

Cloud frameworks are **provider-specific**. They are strong for operational architecture review but do not replace system-level architecture description (arc42, C4, ADRs).

---

## Decomposition and boundary patterns

- **Business context vs. technical context**: separate them when stakeholders need both. Business context shows domain I/O; technical context shows protocols and infrastructure.
- **C4 decomposition ladder**: people → software systems → containers → components → code. Start at the level that answers the question; don't go deeper than needed.
- **Decouple** into independently operating components when it improves changeability, security-control placement, reliability, monitoring, or cost control.
- **Prefer stateless** where feasible — reduces local dependencies, simplifies restart behavior and horizontal scaling.

---

## Methods and frameworks reference table

| Method / framework | Purpose | Key outputs | Caveats |
|-------------------|---------|-------------|---------|
| ISO 42010 architecture description | Conceptual model for architecture descriptions | Required-content guidance, viewpoint conventions | Withdrawn 2011 edition; use for concepts, not as normative standard |
| arc42 template | Structure architecture docs pragmatically | Consistent 12-section doc set | Template needs judgment on depth; not a complete evaluation method |
| C4 model | Visualize architecture at four abstraction levels | Context, container, component, deployment diagrams | Visualization model only; does not replace ADRs or quality scenarios |
| SEI ATAM | Evaluate architecture against quality goals | Risks, sensitivity/tradeoff points, risk themes | Multi-day, stakeholder-intensive; can be applied continuously |
| AWS Well-Architected | Review workload on six pillars | High-risk findings, improvement backlog | AWS-specific; not a substitute for system-level architecture description |
| Azure Well-Architected | Review workload on five pillars with explicit tradeoffs | Assessment results, Azure Advisor recommendations | Azure-specific |
| Google Cloud Well-Architected | Review cloud architecture with six pillars + principles | Architecture recommendations, modernization guidance | GCP-specific; combine with system-level views |
| ADR (Nygard-style) | Record significant decisions with rationale | Decision history enabling future teams to retrace reasoning | Document only architecturally significant decisions |
| Quality scenarios | Make quality requirements testable | Measurable acceptance criteria per quality goal | Requires discipline to keep scenarios measurable, not vague |
| Utility tree | Prioritize quality scenarios for ATAM | Structured quality-attribute focus with business priorities | Reflects current stakeholder priorities only |

---

## Known scope limits

- **ISO 42010** is a withdrawn standard. Use it for conceptual grounding on viewpoints and architecture description, not as a current normative reference.
- **arc42** is a documentation template and method, not a complete architecture theory. Pair with scenario analysis and formal review when tradeoffs matter.
- **C4** is a visualization model. It does not replace ADRs, quality scenarios, or evaluation methods.
- **Cloud well-architected frameworks** are provider-specific. They cover operational architecture review well but are not vendor-neutral architecture-description methods.
