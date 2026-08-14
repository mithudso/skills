<!-- hub-reference-banner -->
> **Reference file — part of the `technical-instruction` hub.** Formerly the standalone `mongodb-university-certification` skill.
> Sibling topics in this family are now reference files under the hubs (`technical-instruction`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: mongodb-university-certification
description: >-
  Covers MongoDB University platform, certifications, and technical enablement paths. TRIGGER: which MongoDB cert to pursue; learning path for DBA, developer, or Atlas admin; course catalog, in-browser labs, or free content; exam format, cost, proctoring, or retake policy; study strategy for Associate Developer / DBA / Atlas Administrator / Data Modeler; Credly digital badges; student, educator, or partner/SI discounts; positioning University to a customer team; building an onboarding plan via MongoDB University. SKIP: actual MQL/aggregation/indexing/schema content → mongodb-expert; live Atlas perf diagnostics → atlas-diagnostics-expert; certification psychometrics, item writing, cut-score methodology → assessment-certification-design; instructional design theory → technical-training-delivery or instructional-design-course-architecture; 10gen diagnostic tooling → 10gen; Atlas platform technical depth (tier architecture, networking, Admin API) → mongodb-atlas-expert. Education hub → technical-instruction.
version: "1.0.0"
updated: "2026-06-16"
category: mongodb
metadata:
  changelog:
    - "2026-06-16: Initial creation via gap-fill + sko convergence loop (3 iter, 34 Medium+ fixed, CLEAN exit, hub synced)"
triggers:
  - MongoDB certification
  - MongoDB University
  - Associate Developer exam
  - Associate DBA exam
  - Associate Atlas Administrator exam
  - Associate Data Modeler exam
  - MongoDB certification learning path
  - MongoDB skill badge
  - learn.mongodb.com
  - MongoDB Credly badge
  - MongoDB exam cost
  - MongoDB exam prep
  - MongoDB instructor-led training
  - MDB100
  - MDB200
  - MDB300
  - MongoDB student program
  - MongoDB educator program
  - MongoDB partner certification
  - certification discount
  - MongoDB exam proctoring
related_skills:
  - mongodb-expert
  - atlas-diagnostics-expert
  - mongodb-atlas-expert
  - assessment-certification-design
  - technical-training-delivery
  - instructional-design-course-architecture
  - 10gen
---

# MongoDB University & Certification Program

Expert reference for MongoDB's education platform (learn.mongodb.com), its four active Associate certifications, the free Skill Badges program, instructor-led training, and academic/partner enablement programs. Use this skill when advising customers, teams, or MongoDB TAMs on learning and credentialing paths.

All facts verified against official learn.mongodb.com pages as of 2026-06-16. Volatile details (exam question counts, upcoming beta exams) are stamped accordingly.

## Output Format

- **Single-cert recommendation:** Return the relevant Quick Reference table row plus a one-paragraph rationale covering role fit and prep path.
- **Certification comparison:** Return a side-by-side table covering role, format, time, cost, and status.
- **Enablement plan:** Return a numbered phased plan (Skill Badges → Learning Path → Exam → ILT if needed).
- **TBD or unverified data points:** Direct the user to https://learn.mongodb.com/pages/certification-program or certification@mongodb.com. Do not estimate or fabricate exam question counts.

## Quick Reference: Certification Decision Table

| Role | Certification | Questions | Time | Price | Status (2026-06-16) |
|------|--------------|-----------|------|-------|---------------------|
| Application Developer | **Associate Developer** | 53 MCQ | 75 min | $150 | Active, GA |
| Data Modeler / Architect | **Associate Data Modeler** | 70 MCQ (60 scored, 10 unscored) | 105 min | $150 | Active, GA |
| Atlas Platform Admin | **Associate Atlas Administrator** | TBD | 120 min | $150 | Active, GA |
| Self-Managed DBA | **Associate DBA** | TBD | TBD | $150 | **Revision in progress** — beta sign-up open, new version coming; old exam retired |

All exams: online proctored via **ProctorU** (switched from Examity April 2024), English, no prerequisites, pass/fail result with domain score report, Credly digital badge upon passing.

## Quick Reference: Who Gets a Discount?

| Path | Discount |
|------|----------|
| Complete any certification learning path on University | 50% off exam (first attempt only; retakes are full price) |
| GitHub Student Developer Pack (verified student) | 100% off (free) one exam per learning path completed |
| MongoDB for Educators Program (verified faculty) | 100% off after completing Developer or DBA path |
| Attend a MongoDB.local event + complete Developer or DBA path | 75% off exam |
| MongoDB Partner Associate Certification Workshops (channel / SI) | No-cost workshop + certification prep |

## Core Concepts

### 1. MongoDB University Platform

**URL:** https://learn.mongodb.com
**Cost:** All self-paced content is free. Exams are paid ($150 each).
**Platform rebuild:** November 2022 relaunch on a new LMS with modern UX; earned Brandon Hall Group Silver Excellence Award for Best Advance in Creating an Extended Enterprise Learning Program (2023).

**Content types:**
- **Units** — short instructional modules (15–60 min) with video, code recap, and knowledge-check quizzes
- **Labs** — hands-on exercises in an in-browser development environment (no local install needed); can be accessed standalone or embedded in courses
- **Courses** — structured collections of units around a theme
- **Learning Paths** — curated sequences of courses that build toward a certification exam; completing a path unlocks a 50% exam discount
- **Learning Bytes** — short explainer videos for quick topic lookup
- **Skill Badges** — 60–90 min focused credentials with a 10-question assessment (see Section 3)
- **Instructor-Led Training Events** — live virtual sessions (see Section 5)

**Language support:** Subtitle translations available in Chinese (Simplified and Traditional), Korean, Spanish, German, Japanese, Italian, French, Portuguese (as of the November 2022 platform relaunch).

**Scale:** Hundreds of courses, labs, and video modules spanning all MongoDB domains (as of 2026-06-16).

### 2. Certification Program

**Program overview page:** https://learn.mongodb.com/pages/certification-program
**Certification contact:** certification@mongodb.com

MongoDB certifications are industry-recognized credentials tied to specific Associate-level roles. They include a Credly digital badge, inclusion in the Credly Talent Directory (recruiter-visible), and a digital certificate.

#### Active Certifications (verified 2026-06-16)

**A. MongoDB Associate Developer**
- Page: https://learn.mongodb.com/pages/mongodb-associate-developer-exam
- Validates: Building beginner-to-intermediate applications using MongoDB as backing database
- Languages: C#, Java, Node.js, Python, PHP — shared core questions plus language-specific questions for each variant. Each language variant is a separate exam purchase ($150 each); a developer testing in two languages pays $150 × 2.
- Format: 53 MCQ, online proctored, 75 min, $150 per language variant
- Recommended prep: Complete the language-specific MongoDB Developer Learning Path (~20 hours each)
- Study materials: Exam Study Guide, Practice Questions, Certification Deep Dive video

**B. MongoDB Associate Data Modeler**
- Page: https://learn.mongodb.com/pages/mongodb-associate-data-modeler-exam
- Validates: Data modeling expertise for modern MongoDB applications — relational-to-document model migration, schema design patterns, anti-patterns, optimization
- Format: 70 MCQ (60 scored, 10 unscored), online proctored, 105 min, $150; available since March 18, 2024
- Recommended prep: MongoDB Data Modeling Path (~4.25 hours focused path)
- Study materials: Associate Data Modeler Exam Study Guide, Practice Questions

**C. MongoDB Associate Atlas Administrator**
- Page: https://learn.mongodb.com/pages/mongodb-associate-atlas-administrator-exam
- Validates: Designing, operating, and managing deployments with MongoDB Atlas — security, monitoring, performance optimization, version upgrades
- Format: Online proctored, 120 min, $150; launched as beta ~May 2024, now GA; updated learning path released May 29, 2026 (path v1 still accessible)
- Recommended prep: MongoDB Atlas Administrator Path (13 hours, current path). New learners should use the current path. Learners mid-way through v1 may complete it — both versions unlock the 50% exam discount.
- Note (2026-06-16): A free beta exam is available through July 14, 2026 (verify current pricing at learn.mongodb.com before advising customers — this window may have passed)

**D. MongoDB Associate DBA (Self-Managed)**
- Page: https://learn.mongodb.com/pages/mongodb-associate-database-administrator-exam
- Validates: Building, supporting, and securing MongoDB self-managed infrastructure
- **Status (2026-06-16): Revision in progress.** The old exam is retired; beta sign-up is open at the page above. Candidates can start the "MongoDB Database Admin Path (Self-Managed)" now for a 50% discount when the new exam launches.
- Recommended prep: MongoDB Database Admin Path (Self-Managed) — 13 hours; covers CRUD, aggregation, indexing, query optimization, sharding, monitoring, cluster reliability, data resilience, security (AuthN/AuthZ, network, encryption at rest)

#### Exam Administration Details

- **Proctoring vendor:** ProctorU (switched from Examity April 2024); 24/7 availability
- **Scheduling:** Done through the learner dashboard on learn.mongodb.com after purchase
- **Retake policy:** 15-day cooling-off period between retakes; retakes cost $150 (same as original). The 50% learning-path discount applies to the first attempt only; retakes are charged at the full $150 price. Confirm current retake pricing at certification@mongodb.com before advising a candidate who earned a discount voucher.
- **Rescheduling:** Free cancel/reschedule up to 24 hours before the scheduled time
- **Accommodations:** Extended time available for candidates with disabilities or non-English speakers; request at certification@mongodb.com (72-hour advance notice required)
- **Score reporting:** Pass/fail with domain-level breakdown. Pass threshold is not disclosed globally; each Exam Study Guide documents domain weightings but not the raw cut score.
- **Recertification / expiry:** No fixed expiration schedule is published as of 2026-06-16. Direct customers to https://learn.mongodb.com/courses/program-guide for current policy.
- **Credly:** Badges issued by Credly (support via https://support.credly.com); certified candidates appear in Credly Talent Directory

### 3. MongoDB Skill Badges (Free Micro-Credentials)

**Introduced:** February 2025
**Format:** 60–90 minutes of self-paced content + 10-question assessment; free; Credly badge upon passing
**Catalog includes (2026-06-16):**
- MongoDB Overview (beginner; for any role)
- MongoDB Basics for Students
- Data Modeling Skills for Developers (series of 4 badges: Relational to Document Model, Schema Design Patterns, Schema Anti-Patterns, Schema Design Optimization)
- Indexing Design Fundamentals
- Query Optimization
- Fundamentals of Data Transformation (aggregation pipeline)
- Search Fundamentals
- Vector Search Fundamentals
- Vector Search Performance
- Building GenAI Apps (Learning Badge Path, 7.5 hours)
- Deploying and Evaluating GenAI Apps (Learning Badge Path)

**Benefit for customers:** Skill Badges suit quick, focused upskilling without the commitment of a full certification. Each badge earns Credly Talent Directory inclusion, useful for teams demonstrating MongoDB fluency to management or for individuals updating LinkedIn profiles.

**TAM positioning note:** Skill Badges are a low-barrier first step. Recommend them to new customer teams before progressing toward full certifications.

### 4. Learning Paths by Role

| Role | Recommended Path | Hours | Discount Unlocked |
|------|-----------------|-------|-------------------|
| New developer (Python) | MongoDB Python Developer Path | ~20 hr | 50% Developer exam |
| New developer (Node.js) | MongoDB Node.js Developer Path | ~20 hr | 50% Developer exam |
| New developer (Java) | MongoDB Java Developer Path | ~20 hr | 50% Developer exam |
| New developer (C#) | MongoDB C# Developer Path | ~20 hr | 50% Developer exam |
| New developer (PHP) | MongoDB PHP Developer Path | ~20 hr | 50% Developer exam |
| Data architect / modeler | MongoDB Data Modeling Path | ~4.25 hr | 50% Data Modeler exam |
| Atlas admin | MongoDB Atlas Administrator Path | ~13 hr | 50% Atlas Admin exam |
| Self-managed DBA | MongoDB Database Admin Path (Self-Managed) | ~13 hr | 50% DBA exam (new version) |
| GenAI developer | MongoDB GenAI Developer Path | ~10 hr | Skill Badge credentials (no exam discount) |
| MongoDB.local attendee | {Learning} Path to .local | ~4 hr | 75% off exam (at event) |

All paths are free; assessments and labs are included.

### 5. Instructor-Led Training (ILT)

**URL:** https://learn.mongodb.com/pages/public-training-calendar
**Delivery:** Virtual (Zoom + Instruqt hands-on labs); private sessions available for teams (7+ learners)
**Subscription model:** Training Subscription for unlimited live + on-demand access; Precision Learning Programs (PLP) for 20+ learner cohorts with post-assessments, detailed reports, and certification tracking. TAMs can request a PLP proposal through the MongoDB training team; completion reports and certification tracking are delivered to the customer's designated admin on learn.mongodb.com.

**Core public curriculum (each day = 8 hours including labs):**

| Course | Content | Prerequisite |
|--------|---------|--------------|
| MDB100: MongoDB Database & Security | Atlas setup, document model, CRUD, security fundamentals | None |
| MDB200: MongoDB Optimization & Performance | Indexing theory, profiling, logs, metrics, Atlas Search/Vector Search intro | MDB100 |
| MDB300: MongoDB Production Readiness | Replication, sharding, backup options, cluster architecture | MDB100+MDB200 |
| DEV400: MongoDB Developer Extension | Advanced developer topics | MDB100-300 |

**Operations / Atlas curriculum:**
- OFA400: Atlas Admin (8 hours)
- OFA500: Atlas API (8 hours)
- OFS400/500: Ops Manager Admin / API (8 hours each)

**Recommended sequence for customers:** All engineers complete MDB100 → MDB200 → MDB300 before successive courses. This sequence aligns with the Developer and DBA certification learning paths.

**Multiple languages available** for private sessions (coordinate with MongoDB representative).

### 6. Academic & Partner Programs

#### MongoDB for Students
- Accessed via GitHub Student Developer Pack
- Benefits: $50 Atlas credit + one free certification exam per completed learning path
- Eligibility: degree/diploma-granting program enrollment with verifiable school email or documentation
- URL: https://www.mongodb.com/students

#### MongoDB for Educators (Academia Program)
- Benefits: $500 Atlas credit + free certifications after completing Developer or DBA learning path
- Eligibility: faculty at accredited institutions (high schools, colleges, universities); bootcamp instructors (eligibility case-by-case — direct to academia@mongodb.com)
- Contact: academia@mongodb.com
- Curriculum resources provided: MongoDB 101, Introduction to Modern Databases, Querying in Non-Relational Databases, MongoDB Aggregation Framework (all free, modular, with hands-on labs)
- 1,000+ universities in the MongoDB for Educators community

#### MongoDB Partner / SI Enablement
- **MongoDB Partner Associate Certification Workshops** (channel partners, SIs, government resellers): No-cost workshops; cover MongoDB core capabilities, Atlas for Government, vector search / gen AI use cases. Note: 4.8 NASBA CPE credits were offered through Carahsoft for a specific 2025 government-sector workshop series (Source 19); this is a Carahsoft-administered event feature, not a standing program benefit. Confirm CPE availability with the MongoDB training team before citing it to partners.
- Partners can access University content and certifications for their technical teams
- For enterprise team training with tracking and ROI reporting, recommend Precision Learning Programs (PLP)

### 7. TAM Positioning Playbook

**When a customer team is new to MongoDB:**
1. Start with free Skill Badges (MongoDB Overview + Data Modeling) — under 90 minutes, no exam required
2. Recommend the role-appropriate Learning Path (Developer or DBA/Atlas Admin)
3. After path completion, the 50% discount makes the $150 exam effectively $75
4. For large teams (20+), propose a Precision Learning Program with tracking and ROI reporting

**When a customer wants to validate a team's proficiency:**

| Team Type | Recommended Certification |
|-----------|--------------------------|
| Developer teams | Associate Developer (choose language variant) |
| Atlas platform/ops teams | Associate Atlas Administrator |
| Self-managed DBA teams | Start the DBA learning path now; wait for refreshed exam |
| Data architects | Associate Data Modeler |

**When a customer is evaluating MongoDB before committing:**
MDB100 (1-day ILT) requires the least time commitment for a technical evaluator; it covers core strengths, Atlas setup, CRUD, and security — no prior MongoDB knowledge required.

**When a partner or SI is involved:**
- Point them to Partner Associate Certification Workshops (no-cost; CPE credits are not a standing program benefit — see Section 6 for details before citing to partners)
- For teams building a MongoDB practice, recommend Training Subscription for unlimited access

**When a customer is building AI/GenAI applications:**
- Recommend the MongoDB GenAI Developer Path, Vector Search Fundamentals Skill Badge, and the Building GenAI Apps badge path (7.5 hours)

**When a customer fails an exam:**
Remind them of the 15-day wait; point them to the domain score breakdown in their result report to identify weak areas; recommend targeted Skill Badges or ILT sessions for those domains before retaking.

**For customers in non-English-speaking regions:**
All certification exams are in English only. Non-English speakers may request extended time at certification@mongodb.com (72 hours advance notice). University self-paced content is available with subtitles in Chinese (Simplified/Traditional), Korean, Spanish, German, Japanese, Italian, French, and Portuguese — recommend starting there before exam registration.

**Common gotchas:**
- The Associate DBA exam is **under revision** (2026-06-16). Do not register customers for it; recommend the learning path and beta sign-up instead.
- The Atlas Admin Path has **two versions** (v1 and current). Direct new learners to the current path; v1 remains valid for learners already enrolled.
- Retakes require a 15-day wait at the same $150 cost — budget-conscious customers should prepare thoroughly before attempting.
- When designing a MongoDB certification program from scratch (item writing, JTA, cut scores, program governance), route to `assessment-certification-design`, not this skill.
- For questions about *how to deliver* MongoDB ILT (facilitation, lab design, live coding logistics), route to `technical-training-delivery`.

## Sources

1. MongoDB Associate Developer Exam page — https://learn.mongodb.com/pages/mongodb-associate-developer-exam (verified 2026-06-16)
2. MongoDB Associate DBA Exam page — https://learn.mongodb.com/pages/mongodb-associate-database-administrator-exam (verified 2026-06-16)
3. MongoDB Associate Data Modeler Exam listing — https://learn.mongodb.com/courses/mongodb-associate-data-modeler-exam (verified 2026-06-16)
4. MongoDB University catalog — https://learn.mongodb.com/catalog (verified 2026-06-16)
5. MongoDB Atlas Administrator Path page — https://learn.mongodb.com/learning-paths/mongodb-atlas-administrator-path (verified 2026-06-16)
6. MongoDB Developer Learning Paths overview — https://learn.mongodb.com/pages/mongodb-developer-learning-paths (verified 2026-06-16)
7. MongoDB Database Admin Path (Self-Managed) — https://learn.mongodb.com/learning-paths/mongodb-dba-certification-learning-path (verified 2026-06-16)
8. New Atlas Administrator Learning Path and Certification blog — https://www.mongodb.com/company/blog/news/new-atlas-administrator-learning-path-and-certification (May 2024, updated Oct 2025)
9. A Year of Thrill: Celebrating the New MongoDB University — https://www.mongodb.com/company/blog/year-of-thrill-celebrating-new-mongodb-university (Nov 2023, updated Mar 2025)
10. Introducing Next Generation of MongoDB Education — https://www.mongodb.com/blog/post/introducing-next-generation-mongodb-education (Nov 2022)
11. MongoDB for Students page — https://www.mongodb.com/students (verified 2026-06-16)
12. MongoDB for Educators (Academia) page — https://www.mongodb.com/academia (verified 2026-06-16)
13. MongoDB Skill Badges introduction — https://www.mongodb.com/community/forums/t/introducing-mongodb-skill-badges (Feb 2025)
14. MongoDB GenAI Learning Badges blog — https://www.mongodb.com/company/blog/news/introducing-two-mongodb-generative-ai-learning-badges (Oct 2024)
15. MongoDB Retake Policy (10gen Education Zendesk) — https://10geneducation.zendesk.com/hc/en-us/articles/31457469358477-Retake-Policy (updated Dec 2025)
16. Overview of the Testing Experience (10gen Education Zendesk) — https://10geneducation.zendesk.com/hc/en-us/articles/31456057240973 (verified 2026-06-16)
17. MongoDB Instructor-Led Training page — https://www.mongodb.com/services/training (verified 2026-06-16)
18. MongoDB Public Training Calendar — https://learn.mongodb.com/pages/public-training-calendar (verified 2026-06-16)
19. MongoDB Partner Certification Workshop Series 2025 (Carahsoft) — https://carahevents.carahsoft.com/Event/Details/576135-Web
20. DataCamp MongoDB Certification Guide — https://www.datacamp.com/blog/mongodb-certification (Aug 2024) [secondary/third-party source; used only for corroboration of publicly-available exam facts, not as authoritative policy reference]
21. MongoDB Atlas Learning Hub — https://www.mongodb.com/resources/product/platform/atlas-learning-hub (verified 2026-06-16)
22. MongoDB Data Modeling Skill Badges path — https://learn.mongodb.com/learning-paths/mongodb-data-modeling-skills-for-developers (verified 2026-06-16)
