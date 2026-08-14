<!-- hub-reference-banner -->
> **Reference file — part of the `technical-instruction` hub.** Formerly the standalone `skills-taxonomies-cbe` skill.
> Sibling topics in this family are now reference files under the hubs (`technical-instruction`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: skills-taxonomies-cbe
title: "Skills Taxonomies & Competency-Based Education"
version: "1.0.0"
updated: "2026-06-16"
category: technical-instruction
description: "Design, navigate, and apply skills taxonomies and competency frameworks (O*NET, ESCO, SFIA 9, Lightcast), build competency models using KSAOs and job task analysis, and implement Competency-Based Education (CBE) and Competency-Based Training (CBT) programs where mastery replaces seat-time. Covers the skills-framework-to-credential link — these are the competency substrates that micro-credentials and certifications certify against. TRIGGER: O*NET or ESCO skills taxonomy; SFIA 9 proficiency levels; building a competency framework from KSAOs; competency modeling or job task analysis; DACUM method; competency-based education or CBE direct assessment; mastery-based progression; skills gap analysis; skills-based hiring or skills-based organization; skills ontology or skills graph; micro-credential competency alignment; Open Badges or 1EdTech CLR substrate. SKIP: psychometrics, item writing, cut scores, or accreditation of the EXAM that certifies a competency → assessment-certification-design (that skill owns the examination instrument; this skill owns the competency/skills framework it certifies against); course architecture, curriculum sequencing, or learning-objective taxonomy → instructional-design-course-architecture; MongoDB University learning paths → mongodb-university-certification; Kirkpatrick evaluation or training ROI measurement → learning-measurement-evaluation; broad HR performance management without a learning design angle."
whenToUse: "Use when working with skills taxonomies (O*NET, ESCO, SFIA, Lightcast), designing or evaluating competency frameworks, building CBE/CBT programs, analyzing skills gaps, supporting skills-based hiring initiatives, or connecting competency evidence to credential issuance."
keywords:
  - O*NET
  - ESCO
  - SFIA
  - competency framework
  - competency modeling
  - KSAOs
  - CBE
  - competency-based education
  - CBT
  - mastery-based progression
  - direct assessment
  - skills taxonomy
  - skills gap analysis
  - skills-based hiring
  - skills-based organization
  - skills ontology
  - DACUM
  - job task analysis
  - micro-credential substrate
  - Open Badges
  - 1EdTech
  - Lightcast
  - skills graph
  - proficiency levels
tags:
  - technical-instruction
  - competency
  - skills
  - education
  - workforce-development
  - credentials
---

# Skills Taxonomies & Competency-Based Education

<!-- Family hub → technical-instruction -->

This skill covers the **competency and skills framework substrate** — the classification systems, modeling methods, and learning approaches that define what skills are, how they are organized, and how demonstrated competency replaces seat-time as the unit of educational progress. Certifications and micro-credentials certify *against* these frameworks; assessment design is a separate concern (see SKIP section).

**Verified-as-of: 2026-06-16**

## Contents

- [Quick Reference — Framework Decision Table](#quick-reference)
- [Core Concepts](#core-concepts)
  - [O*NET: US Occupational Taxonomy](#onet)
  - [ESCO: EU Skills Ontology](#esco)
  - [SFIA 9: IT/Digital Proficiency Framework](#sfia-9)
  - [Skills Ontologies and Skills Graphs](#skills-ontologies)
  - [KSAOs: The Competency Modeling Unit](#ksaos)
  - [Competency Framework Design (JTA, DACUM, Proficiency Scales)](#framework-design)
  - [Competency-Based Education (CBE) and CBT](#cbe)
  - [Skills-Based Hiring and Skills-Based Organizations](#sbo)
  - [Micro-credential Substrate and Credential Infrastructure](#micro-credential-substrate)
- [Practical Patterns](#practical-patterns)
- [Anti-Patterns and Failure Modes](#anti-patterns)
- [Boundary: This Skill vs. Assessment-Certification-Design](#boundary)
- [References](#references)

---

## Quick Reference — Framework Decision Table {#quick-reference}

| Use Case | Best Framework | Notes |
|---|---|---|
| Classify US jobs, map occupations to skills | O*NET (onetcenter.org) | 1,016 occupational titles; authoritative DOL source |
| EU labor market matching, multilingual | ESCO v1.2.1 | 3,039 occupations; 13,939 skills; 28 languages; free API |
| IT/digital career ladders with proficiency levels | SFIA 9 | 7 levels (Follow → Set strategy); 120+ skills; Oct 2024 |
| Enterprise skills intelligence / commercial use | Lightcast Open Skills | 33,000+ skills; adjacency data; publicly available |
| Define human attributes for a job/role | KSAOs (K, S, A, Other) | Foundation of I/O psychology job analysis |
| Build a competency framework from scratch | DACUM + proficiency anchors | 8–12 expert-worker panel; outputs duties + tasks map |
| Replace seat-time with demonstrated mastery | CBE / direct assessment | Mastery = fixed; time = variable |
| Close an organizational skills gap | Skills gap analysis | Current-state inventory → future-state needs heat map |
| Connect a badge/cert to a competency | 1EdTech CLR + Open Badges 3.0 | Technical substrate for competency-aligned credentials |

---

## Core Concepts

### O*NET: US Occupational Taxonomy {#onet}

**O*NET (Occupational Information Network)** is maintained by the US Department of Labor and is the primary source of occupational information in the United States.[^1] It covers 1,016 occupational titles representing 55,000+ jobs, using a six-digit SOC code plus two additional digits for granularity. The current database is O*NET 30.3 (released May 2026).

**Content Model — Three Dimensions:**

*Worker Domain*
- **Abilities (52 variables):** Organized across four domains — cognitive, psychomotor, physical, and sensory — at three abstraction tiers: 52 specific → 15 intermediate → 4 general[^2]
- **Skills (35 skills):** Grouped into 7 categories: content, process, social, complex problem solving, technical, systems, and resource management
- **Knowledge (33 areas):** Domain content knowledge (e.g., mathematics, customer service, medicine)
- **Work Styles (21 variables):** Personality tendencies affecting performance

*Job Domain*
- **Work Activities:** 19,000+ task statements → 2,000+ Detailed Work Activities → 325 intermediate → 41 generalized work activities
- Work context and environment variables

*Market Domain*
- Labor market projections, occupational outlook

**O*NET 30.3 Database Counts (verified May 2026):**[^1]
- Software Skills: 31,821 rows (includes 176 Hot Technologies)
- Transferable Skills: 44,700 rows (new in 30.3)
- Essential Skills: 17,880 rows (new in 30.3)
- Knowledge: 59,004 rows
- Abilities: 92,976 rows

O*NET publishes machine-readable competency frameworks in CTDL JSON-LD and ASN schema formats, enabling integration with digital credential platforms.

**Limitations (disconfirming evidence):** A National Academies review found occupations are updated infrequently — an average of only 71 observations per occupation from a single time point; data for some occupations was last collected in 2006.[^4] The 1,016 categories are considered too broad for precision ML applications. Critics also document vague item wording and limited predictive validity for some use cases.[^5]

---

### ESCO: EU Skills Ontology {#esco}

**ESCO (European Skills, Competences, Qualifications and Occupations)** is governed by the European Commission (DG EMPL) and functions as "a dictionary, describing, identifying and classifying professional occupations and skills relevant for the EU labour market."[^6] Current version: ESCO v1.2.1 (December 2025).

**Three-Pillar Architecture:**[^7]
1. **Occupations pillar** — 3,039 professional roles, hierarchically organized, mapped to ISCO-08 (International Standard Classification of Occupations)
2. **Skills/Competences pillar** — 13,939 skills linked to occupations via essential/optional relationships; both hierarchical and associative inter-skill relationships
3. **Qualifications pillar** — formal credentials mapped to both occupations and skills, enabling cross-border equivalency

**Key Technical Properties:**
- Available as Linked Open Data (RDF/SKOS format) and via REST API (free)
- Translated into 28 languages (all EU official languages + Icelandic, Norwegian, Ukrainian, Arabic)
- v1.2 (May 2024) used both human expertise and AI techniques to add 35 new occupations, 42 new skills, and 196 new knowledge concepts[^8]

**Relationship to National Frameworks:** ESCO maps to the European Qualifications Framework (EQF), enabling cross-border credential recognition. National employment agencies (Germany's BA, France's ROME) publish ESCO alignment mappings.

**Limitation:** Despite ESCO's open API, API connections between ESCO and commercial HR platforms are not yet standardized — a skills defined in ESCO does not map semantically cleanly to the same skill in Workday or LinkedIn.[^17]

---

### SFIA 9: IT/Digital Proficiency Framework {#sfia-9}

**SFIA (Skills Framework for the Information Age)** is maintained by the SFIA Foundation (global non-profit) and provides a two-axis competency matrix for digital and IT professionals.[^9] SFIA 9 was published October 2024; SFIA 10 is in active consultation.

**The 7 Levels of Responsibility (verbatim):**[^9]

| Level | Name | Core Character |
|---|---|---|
| 1 | **Follow** | Routine tasks under close supervision; requires guidance |
| 2 | **Assist** | Routine supervision; uses discretion for routine problems |
| 3 | **Apply** | Varied tasks; standard methods; manages own work |
| 4 | **Enable** | Diverse complex activities; delegates; works autonomously |
| 5 | **Ensure, advise** | Authoritative guidance; accountable for significant outcomes |
| 6 | **Initiate, influence** | Significant organisational influence; high-level decisions |
| 7 | **Set strategy, inspire, mobilise** | Highest level; overall vision and strategy |

**Five Generic Attributes (apply at every level):**
- Autonomy — independence and accountability for results
- Influence — reach and impact of decisions
- Complexity — range and intricacy of responsibilities
- Business Skills / Behavioural Factors — communication, planning, problem-solving
- Knowledge — depth and breadth required

**Skill Categories (6):** Strategy & architecture; Change & transformation; Development & implementation; Delivery & operation; People & skills; Relationships & engagement.

SFIA 9 covers 120+ professional skills; not every skill exists at all 7 levels — each skill specifies the levels at which it is applicable. SFIA 9 added 26 new skills and 60 new level-skill combinations over SFIA 8; it also introduced AI skills.

**Organizational Use:** Organizations build role profiles specifying 8–12 must-have and nice-to-have skills at defined levels. Role profiles enable structured career pathways, gap analysis, and L&D targeting. BCS (British Computer Society) distributes SFIAplus, which adds behavioral competencies atop SFIA's technical skills.[^10]

**Limitations:** Peer-reviewed examination of SFIA v7 found: absence of universal certification criteria for objective skill assessment, complexity in skill/proficiency mapping, limited soft-skills representation, and limited scope for automating skill management.[^11] Self-rating without calibration is the common practice and creates inflation risk. Version incompatibility between major releases (v7 not backward compatible with v6) creates interoperability problems across organizations.[^9]

---

### Skills Ontologies and Skills Graphs {#skills-ontologies}

A **skills ontology** is a formal representation of skills and their relationships: broader/narrower (hierarchical), related (associative), and equivalent (synonymy). Implemented as a graph database, nodes represent skills and edges represent relationships.[^12]

**Skill Adjacency:** Skills that frequently co-occur in job postings or tend to be learned together have high adjacency scores. Adjacency enables inference of likely-but-unverified skills and "adjacent role" recommendations — a person with skill A is likely to have skill B if A and B are highly adjacent.[^13]

**Major Platforms (verified 2026):**

*Lightcast (formerly EMSI + Burning Glass Technologies)*
- Maintains the largest integrated occupational data asset; 33,000+ Open Skills in a 3-tier hierarchy: 31 top-level Categories → Subcategories → individual skills[^13]
- Methodology: hybrid NLP/ML over hundreds of millions of job postings + human expert review
- Open Skills publicly available; API exposes adjacency matrices
- Josh Bersin describes Lightcast as "the Switzerland of the SkillsTech market" — neutral data provider used by competing HCM platforms[^15]

*Workday Skills Cloud*
- Universal skills ontology built as a knowledge graph within Workday HCM
- 2024: added LLM-based skill inference — infers skills from brief job-title inputs, post-processes against the ontology, ranks by relevance[^14]

*Eightfold AI*
- AI-native talent intelligence; infers skills from career histories; distinctive capability is adjacent skills inference for role-shift recommendations

**Fragmentation persists:** Despite open standards (ESCO, O*NET, WEF Global Skills Taxonomy), commercial platforms maintain proprietary ontologies with inconsistent semantics. HR Open Standards Consortium (2024) found API connections not yet standardized — a "data analysis" skill on LinkedIn is not semantically identical to one in Workday.[^17]

---

### KSAOs: The Competency Modeling Unit {#ksaos}

**KSAOs (Knowledge, Skills, Abilities, and Other Characteristics)** is the canonical framework in industrial-organizational (I/O) psychology for specifying the human attributes required for effective job performance.[^18]

**Four Components:**
- **Knowledge:** Understanding of concepts, facts, and domain content (formal and experiential)
- **Skills:** Learned/acquired technical or manual proficiencies — demonstrated behaviors
- **Abilities:** Relatively stable, underlying cognitive or physical capacities (what cognitive tests measure)
- **Other Characteristics:** Personality traits, attitudes, values, certifications, work styles

KSAOs are elicited through job analysis — task inventories, SME panels, structured interviews — and validated through empirical study. When validated, KSAO models ground evidence-based talent decisions across the full employee lifecycle: job descriptions, selection assessments, training needs analysis, performance standards, and certification requirements.

**Competency dictionaries:** Large organizations and government agencies maintain standardized KSAO libraries from which job-specific profiles are assembled.

**Critical limitation:** I/O psychology research explicitly states competency modeling "is considerably less rigorous and often fails to achieve the levels of reliability and validity attained by job analytic methods."[^21] Competencies are "heterogeneous constructs" that bundle multiple KSAOs, making valid assessment design difficult. Competency frameworks with more than 10–12 competencies become unwieldy; common practice produces 15–20+ competency frameworks that practitioners cannot reliably use.[^21]

---

### Competency Framework Design (JTA, DACUM, Proficiency Scales) {#framework-design}

**Job Task Analysis (JTA)** is the systematic process of identifying all tasks performed in a role, the frequency and criticality of each, and the KSAOs required for each task. JTA outputs ground curriculum design, certification blueprints, and hiring criteria.

**DACUM (Developing A Curriculum)** is the most widely validated JTA facilitation method.[^22][^23] A panel of 8–12 expert workers guided by a skilled facilitator works through seven steps:
1. Orient panel members to the method
2. Review and confirm the scope of the occupation
3. Identify major duties (typically 8–12)
4. Identify tasks under each duty (typically 6–20 per duty)
5. Review and refine all duty/task statements
6. Sequence statements in performance order
7. Identify which tasks are entry-level requirements

The DACUM chart maps all duties and tasks for a role, with each task statement following the "action verb + object + condition" format. Originated at Ohio State University (DACUM International Training Center); widely used by community colleges, professional associations, government agencies, and corporations.[^23][^24]

**Proficiency scales:** Competency frameworks define 3–5 proficiency levels per competency, typically with behavioral anchors (Behaviorally Anchored Rating Scales, BARS) that describe observable behavior at each level. This creates the bridge from abstract competency definitions to assessable performance evidence.

**Building a competency framework — general sequence:**
1. Conduct job task analysis (DACUM or equivalent)
2. Group tasks into competency clusters
3. Define KSAOs required per cluster
4. Write proficiency level anchors (BARS)
5. Validate with SMEs and incumbents
6. Pilot assessment against standards
7. Map to external frameworks (O*NET, SFIA, ESCO) for portability

---

### Competency-Based Education (CBE) and CBT {#cbe}

**Competency-Based Education (CBE)** is "a teaching and learning approach in which curricula are designed around competencies (outcomes), and learner progression is based solely on the performance-based demonstration of those competencies, not the amount of time spent on task."[^26] Competency-Based Training (CBT) is the same model applied in workplace/vocational contexts.

**The core inversion:** Traditional education fixes time (a semester) and varies what students learn. CBE inverts this — **time is variable** and **mastery is the fixed target**.[^27]

**Two US implementation tracks:**[^28][^29]

*Credit-based CBE:* CBE within the existing credit-hour system. No special federal approval required. The large majority of CBE programs operate this way.

*Direct assessment CBE:* Measures progress entirely by demonstrated competency with no credit-hour denominator. Requires explicit Department of Education authorization under 34 CFR 668.10. To disburse Title IV federal student aid, an institution must: (1) receive accreditor approval including credit-hour equivalency methodology; (2) ensure "regular and substantive interaction" between students and instructors; (3) maintain programs of at least 10 weeks; (4) obtain final ED authorization. Only a handful of US institutions have direct assessment authorization — it is rare in practice despite interest.[^29]

**C-BEN (Competency-Based Education Network):** The primary US practitioner community for CBE. Provides consulting, technical assistance, research, policy advocacy, and convening (CBExchange conference). Developed the first Quality Principles & Standards for CBE programs (2017) through a collaborative process involving 100+ institutional stakeholders.[^26]

**CBE Mastery Progression:**
- Learner is assessed against a defined competency standard
- If mastery is not demonstrated, learner receives additional instruction and is re-assessed
- Progression continues iteratively until the standard is met
- Learner advances when ready, not when the semester ends

**Critical limitations and disconfirming evidence:**
- Federal audits found accreditor reviews of direct assessment programs were inconsistent — auditors found HLC failed to require sufficient information or adequately assess credit-hour equivalency claims[^29]
- C-BEN advises institutions to add "at least one year" to timelines just for direct assessment approval — a structural adoption barrier
- Critics argue CBE's outcomes focus may eliminate "much of what makes a college education valuable, including intensive interaction with Ph.D.s and peers"
- Equity stratification risk: CBE programs may funnel lower-income students into narrower vocational tracks
- Large class sizes prevent the individualized feedback loops mastery-based progression requires[^30]
- Implementation quality problems are documented as more prevalent than definitional flaws (the medical CBE literature distinguishes these)[^31]

---

### Skills-Based Hiring and Skills-Based Organizations {#sbo}

**Skills-based hiring** removes degree requirements and replaces them with skill-based criteria — work sample tests, structured skills assessments, demonstrated competency evidence. Major adopters include IBM, Google, Apple, and the US federal government.[^32]

**Adoption statistics (2025, qualified — single-source ecosystem risk):**[^32]
- 85% of companies claim to use skills-based hiring (up from 81% in 2024)
- 76% of employers say it produces better outcomes than degree-based hiring
- When degree requirements removed: reported 19× larger qualified candidate pool; 50% reduction in time-to-hire

**Critical disconfirming evidence:** Harvard/Burning Glass research found that while 85% of companies claim skills-based hiring, at large firms fewer than 1 in 700 hires (0.14%) are actually non-degree graduates hired through skills-based pathways. The claim and the practice are substantially disconnected.[^33]

**The skills-based organization (SBO):** Deloitte defines an SBO as one that "places skills and human capabilities at the heart of talent strategies, creating a new operating model for work and the workforce."[^34] The SBO concept deconstructs jobs into tasks and reconstructs work based on skills, aiming for greater organizational agility. In a 2024 Workday survey, 55% of 2,300 business leaders reported beginning the transition.

**Skills gap analysis methodology:**
1. Build a skills inventory (current-state: what skills employees have, at what proficiency level)
2. Define future-state requirements (roles, projects, strategic direction)
3. Map delta: gap heat maps by role × skill × proficiency
4. Prioritize by role criticality and gap severity
5. Build L&D, hiring, or reskilling plans to close priority gaps

**Disconfirming evidence on SBO:** Nearly 70% of companies state skills-based hiring is "too difficult to implement right now," with the barrier being not technology but lack of clarity on how to build a skills program.[^33] Deloitte's own follow-up research shows many organizations are struggling to realize tangible value from skills-based approaches despite significant investment.[^35] Josh Bersin's analysis characterizes implementation as delivering consistent outcomes only in early-career hiring.[^36]

---

### Micro-credential Substrate and Credential Infrastructure {#micro-credential-substrate}

Skills frameworks are the **competency substrate** that micro-credentials and certifications certify against. A certification asserts that a holder has demonstrated mastery of a defined competency — the certification exam (see `assessment-certification-design`) is the instrument, but the competency framework is what the instrument certifies against.

**1EdTech Consortium (formerly IMS Global) — Technical Stack:**[^37][^38]
- **Comprehensive Learner Record (CLR Standard):** Aggregates and portably represents learning achievements across institutions and providers
- **Open Badges v3.0:** Granular digital credential issuance aligned to specific competencies; integrates with skills frameworks so badges reference specific competency standards
- **CASE (Competencies and Academic Standards Exchange):** Machine-readable competency framework interchange format enabling systems to ingest, map, and align to any competency framework

**The "certifying against a competency" relationship:**
- A competency framework (O*NET, SFIA, custom) defines: what the competency is, what proficiency looks like at each level, how evidence is produced
- A credential program uses this to: write assessment criteria, design evidence requirements, issue badges aligned to specific framework competencies
- An employer/hiring system uses this to: validate that a credential-holder actually demonstrates required skills

**Stackable credentials:** Learners accumulate micro-credentials (each certifying mastery of a bounded competency) that aggregate toward larger qualifications. The CLR Standard makes this stacking technically portable across institutions.

**1EdTech employer research (qualified — advocacy source):**[^38] 87% of HR leaders have heard of digital badges; 48% have encountered badges in candidate screening; 34% use competency-over-degree hiring strategies.

---

## Practical Patterns

**Pattern 1 — Framework selection by context:**
- US workforce development or hiring → O*NET as the base taxonomy
- EU/cross-border labor market → ESCO (with API)
- IT/digital career ladders → SFIA 9 (most granular proficiency model for tech)
- Enterprise talent intelligence → Lightcast (commercial-grade, adjacency data)
- Learning design / certification → custom KSAO framework grounded in JTA/DACUM

**Pattern 2 — Building a competency framework step-by-step:**
Run DACUM analysis (8–12 expert workers, facilitated, ~2 days) → identify 8–12 major duty clusters → write KSAO specifications per cluster → define 3–5 proficiency level anchors (BARS) → validate with larger SME group → map to O*NET or SFIA for external portability → pass to assessment designers (see `assessment-certification-design`) for exam/badge development.

**Pattern 3 — Connecting CBE to credentialing:**
Define competencies via JTA → write direct assessment criteria → align criteria to 1EdTech CASE framework → issue Open Badges v3.0 with competency alignment → aggregate to CLR for learner portability.

**Pattern 4 — Skills gap analysis:**
Build role-skill matrix (job architecture) → collect employee skills self-assessment + manager calibration → compare to target proficiency by role → generate gap heat map → prioritize by business criticality → trigger L&D, reskilling, or hiring interventions.

**Pattern 5 — SFIA 9 role profiling:**
Select 8–12 skills from SFIA skill catalogue → for each skill, specify required level (1–7) at "must have" or "nice to have" → document role profile → use for job description, career progression, gap analysis, and L&D targeting.

---

## Anti-Patterns and Failure Modes

**Framework overload:** Using O*NET, ESCO, SFIA, and a custom framework simultaneously without a mapping layer. Choose one primary framework; map others as needed. Lightcast functions as a neutral bridge.

**Taxonomy size explosion:** Ontologies growing beyond 10,000–50,000 skills create signal-to-noise problems. Overlap, redundancy, and over-granularity reduce practical utility. Prioritize manageable, well-curated frameworks over exhaustive ones.

**Competency framework bloat:** Frameworks with more than 10–12 competencies become unwieldy. Common failure is producing 15–20+ competency frameworks that practitioners cannot reliably assess, leading to compliance-only use and eventual abandonment.

**Self-rating inflation in SFIA:** Without calibration mechanisms (panel review, 360 evidence, structured behavioral interview), SFIA self-ratings drift upward over time, reducing signal quality.

**CBE seat-time substitution without mastery infrastructure:** Organizations label programs "competency-based" but retain fixed time blocks and summative-only assessment without iterative mastery loops — producing CBE in name only.

**Skills-based hiring claim vs. practice gap:** The 85%/0.14% gap (Harvard/Burning Glass) is the canonical failure mode. Announcing skills-based hiring without changing hiring workflows, assessment tools, or ATS configuration produces no change in outcomes.

**DACUM without follow-through:** Conducting DACUM analysis but not translating the chart into training designs, competency frameworks, or assessment criteria — the analysis becomes a bookshelf artifact.

**Version incompatibility in SFIA:** Organizations on SFIA 7 and SFIA 9 cannot directly compare role profiles without a mapping exercise. Plan migrations deliberately.

**Interoperability assumption:** Assuming that because ESCO, O*NET, and a commercial platform all have "data analysis" skills they mean the same thing. Semantic alignment requires explicit mapping; it does not happen automatically.

---

## Boundary: This Skill vs. Assessment-Certification-Design {#boundary}

This skill covers the **competency/skills substrate**: what skills/competencies are, how they are classified (O*NET, ESCO, SFIA), how they are modeled (KSAOs, JTA, DACUM), how CBE uses them as the basis for mastery-based progression, and how they serve as the framework micro-credentials certify against.

The `assessment-certification-design` skill covers the **examination instrument**: item writing (MCQ and performance-based), psychometrics (IRT, KR-20, reliability, validity), cut scores (Angoff, Bookmark), test blueprints, ANSI/ISO 17024 and NCCA accreditation, and micro-credential mechanics (Open Badges, stackable credentials). That skill presupposes a defined competency framework; this skill builds the framework it tests against.

**The handoff point:** A competency/skills framework is defined here → the assessment instrument certifying mastery of that framework is designed in `assessment-certification-design`.

---

## References

[^1]: O*NET 30.3 Database. US Department of Labor/ETA. onetcenter.org. Released May 2026. Authoritative database row counts and content model structure. https://www.onetcenter.org/database.html

[^2]: O*NET Content Model. onetcenter.org. Detailed breakdown of abilities (52), skills (35), knowledge (33), work activities hierarchy. https://www.onetcenter.org/content.html

[^3]: O*NET OnLine Help: Overview. onetonline.org. Overview of occupational categories and SOC alignment. https://www.onetonline.org/help/onet/

[^4]: National Academies of Sciences. *A Database for a Changing Economy: Review of the Occupational Information Network (O*NET)*. Chapter 6: Update cycle and sample size concerns. Average 71 observations per occupation; some data from 2006. https://www.nationalacademies.org/read/12814/chapter/6

[^5]: ResearchGate. "The O*NET content model: strengths and limitations." Peer-reviewed critique documenting vague item wording, content duplication, limited predictive validity. https://www.researchgate.net/publication/295834355_The_ONET_content_model_strengths_and_limitations

[^6]: ESCO. "What is ESCO?" European Commission DG EMPL. ESCO v1.2.1 (December 2025). Three-pillar architecture; 3,039 occupations; 13,939 skills; 28 languages. https://esco.ec.europa.eu/en/about-esco/what-esco

[^7]: ESCO. "Occupations pillar." ESCO hierarchy and ISCO-08 mapping. https://esco.ec.europa.eu/en/about-esco/escopedia/escopedia/occupations-pillar

[^8]: ESCO. "ESCO v1.2 is live!" May 2024 release announcement. 35 new occupations, 42 new skills, AI-assisted methodology. https://esco.ec.europa.eu/en/news/esco-v12-live

[^9]: SFIA Foundation. SFIA 9 Levels of Responsibility. sfia-online.org. Published October 2024. Verbatim level names and descriptions; 7-level model; 5 generic attributes. https://sfia-online.org/en/sfia-9/responsibilities

[^10]: BCS. "SFIAplus IT Skills Framework." BCS extension adding behavioral competencies to SFIA. https://www.bcs.org/it-careers/sfiaplus-it-skills-framework/

[^11]: ScienceDirect. "An examination of SFIA version 7." Peer-reviewed limitations: no universal certification criteria, complexity, limited soft-skills, automation barriers. https://www.sciencedirect.com/science/article/abs/pii/S0268401219305237

[^12]: Lightcast Open Skills Taxonomy. lightcast.io. 33,000+ skills; 3-tier hierarchy; adjacency data; publicly available. https://lightcast.io/open-skills

[^13]: Lightcast. "The Lightcast Open Skills Taxonomy" (white paper, August 2023). Hybrid NLP + human review methodology; adjacency matrix structure. https://4906807.fs1.hubspotusercontent-na1.net/hubfs/4906807/The%20Lightcast%20Open%20Skills%20Taxonomy%20Aug%202023.pdf

[^14]: Workday Engineering. "Skill Inference: Building an LLM-based Service in the Workday Skills Cloud." Medium/Workday, 2024. LLM-based skill inference pipeline. https://medium.com/workday-engineering/skill-inference-building-an-llm-based-service-in-the-workday-skills-cloud-47c9cce9f7bd

[^15]: Josh Bersin. "Lightcast: The Switzerland of the Global SkillsTech Market." joshbersin.com, May 2023. Market positioning; Lightcast as neutral provider. https://joshbersin.com/2023/05/lightcast-the-switzerland-of-the-global-skillstech-market/

[^16]: Workday. "The Foundation of the Workday Skills Cloud." Workday Blog. Knowledge graph architecture and ML inference overview. https://blog.workday.com/en-us/foundation-workday-skills-cloud.html

[^17]: HR Open Standards Consortium. "HR in 2024: From XML to APIs — Standardizing Skills and Competency Data." 2024. API interoperability gap; semantic inconsistency across platforms. https://www.hropenstandards.org/news/hr-in-2024-from-xml-to-apis---standardizing-skills-and-competency-data-for-the-future

[^18]: APA Dictionary of Psychology. "Knowledge, skills, abilities, and other characteristics (KSAOs)." American Psychological Association. Canonical definition. https://dictionary.apa.org/knowledge-skills-abilities-and-other-characteristics

[^19]: cogn-iq.org. "Competency Frameworks: From Job Analysis to Validated Assessment." Methodology: JTA → KSAO → proficiency anchors → validation. https://www.cogn-iq.org/blog/competency-frameworks/

[^20]: SIOP/Purdue. "Best Practices in Competency Modeling." Best practices document; competency dictionary; behavioral-level definitions. https://apps.it.purdue.edu/sites/Home/DirectoryApi/Files/0e3546ac-6f92-4520-bc41-d38da187f4e0/Download

[^21]: assess.com. "What is a KSAO?" I/O psychology critique: competency modeling less rigorous than formal job analysis; heterogeneous construct problem; >12 competency failure mode. https://assess.com/ksao-and-assessment/

[^22]: ERIC. "DACUM: A Proven and Powerful Approach to Occupational Analysis." ERIC ED346248. Origin at Ohio State University; 7-step facilitation method. https://eric.ed.gov/?id=ED346248

[^23]: DACUM International / OSU. "Applications of DACUM Analysis." Duty/task chart methodology and use cases. https://dacum.osu.edu/wp-content/uploads/2025/08/Applications-of-DACUM-Analysis.pdf

[^24]: PMC/NIH. "Job Analysis and Curriculum Design Using DACUM." PMC11240485. Peer-reviewed validation of DACUM in health professions. https://www.ncbi.nlm.nih.gov/pmc/articles/PMC11240485/

[^25]: PubMed. "DACUM: a versatile competency-based framework for staff development." PMID 11840017. Cross-sector application. https://pubmed.ncbi.nlm.nih.gov/11840017/

[^26]: C-BEN. "What is Competency-Based Education?" c-ben.org. Verbatim definition; Quality Principles & Standards (May 2017). https://www.c-ben.org/competency-based-education/

[^27]: D2L. "The Complete Guide to Competency-Based Education." Time-as-variable vs. seat-time explanation. https://www.d2l.com/blog/the-complete-guide-to-competency-based-education/

[^28]: California Community Colleges Chancellor's Office. "Direct Assessment CBE." Credit-based vs. direct assessment distinction; Title IV implications. https://www.cccco.edu/About-Us/Chancellors-Office/Divisions/Educational-Services-and-Support/da-cbe

[^29]: US Department of Education. "Direct Assessment Competency-Based Programs." 34 CFR 668.10; dual approval requirements (accreditor + ED). https://www.ed.gov/laws-and-policy/higher-education-laws-and-policy/higher-education-policy/direct-assessment-competency-based-programs

[^30]: cloudassess.com. "What is Competency-Based Education? Principles & Components." CBE equity and implementation challenges. https://cloudassess.com/blog/what-is-competency-based-education/

[^31]: PMC/NIH. "Outcomes are what matter: competency-based medical education." 2023. Implementation vs. definitional critique distinction in CBME. https://www.ncbi.nlm.nih.gov/pmc/articles/PMC11240485/

[^32]: TestGorilla. "The State of Skills-Based Hiring 2025 Report." Adoption statistics; outcomes data. https://www.testgorilla.com/skills-based-hiring/state-of-skills-based-hiring-2025/

[^33]: The Interview Guys. "The State of Skills-Based Hiring in 2025: 85% claim it, only 1 in 700 are affected." Harvard/Burning Glass 0.14% finding; 70% cite implementation difficulty. https://blog.theinterviewguys.com/the-state-of-skills-based-hiring/

[^34]: Deloitte. "The skills-based organization: A new operating model for work and the workforce." 2024. SBO definition; 55% of leaders beginning transition. https://www.deloitte.com/content/dam/insights/articles/2024/us175310_consulting-the-skills-based-org-report/di-the-skills-based-organization-report.pdf

[^35]: Deloitte. "Rethinking skills-based talent models: 4 paths to business value." 2024. Value realization struggles; implementation challenges. https://www.deloitte.com/us/en/insights/topics/talent/creating-value-with-skills.html

[^36]: Josh Bersin Company. "The Journey to the Skills-Based Organization: What Works." 2024. Sober reality; consistent outcomes only in early-career. https://26502993.fs1.hubspotusercontent-eu1.net/hubfs/26502993/Guide%20pdfs/JBC_The%20Journey%20to%20the%20Skills-Based%20Organization.pdf

[^37]: 1EdTech. "New Open Badges 3.0 Standard provides enhanced security and mobility." CLR + Open Badges + CASE as the micro-credential technical stack. https://www.1edtech.org/1edtech-article/new-open-badges-30-standard-provides-enhanced-security-and-mobility/411060

[^38]: 1EdTech. "Research Shows Growing Employer Support for Competency Frameworks, Talent Analytics and Skills." 2024. Employer adoption statistics (advocacy source — note bias). https://www.1edtech.org/article/research-shows-growing-employer-support-competency-frameworks-talent-analytics-and-skills

[^39]: ERIC/EJ1326834. "Competency-Based Education: Theory and Practice." Historical and theoretical foundations. https://files.eric.ed.gov/fulltext/EJ1326834.pdf

[^40]: ILO/skillsforemployment. "European Skills/Competences, Qualifications and Occupations." Policy context; EQF mapping; national framework relationships. https://www.skillsforemployment.org/sites/default/files/2024-01/edmsp1_212824.pdf
