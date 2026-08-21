---
description: >-
  Confidence, self & identity psychology (psychology spoke) — replication-honest science of
  self-efficacy, self-esteem, impostor phenomenon, and identity; robust mechanism vs overstated
  claim. TRIGGER: Bandura self-efficacy (4 sources, mastery experience,
  collective/teacher efficacy); self-esteem (Baumeister 2003 vs Orth-Robins/Krueger-Vohs 2022
  debate, sociometer/Leary, contingencies of self-worth/Crocker, threatened egotism);
  impostor phenomenon (Clance & Imes, Tulshyan-Burey reframe); identity formation (Erikson,
  Marcia's 4 statuses, narrative identity/McAdams); self-concept & self-discrepancy (James,
  Higgins actual/ideal/ought); ego threat & self-affirmation (Steele); psychological safety
  (Edmondson). SKIP: growth mindset/grit/flow & "USE confidence" → applied-psychology;
  happiness/PERMA/SDT-wellbeing → psychology-positive; burnout/chronic-stress →
  psychology-stress-trauma; narcissism-as-disorder → psychology-clinical-personality;
  status/dominance/comparison → psychology-social.
name: psychology-confidence-identity
version: "1.0.1"
updated: "2026-06-23"
category: spoke
model: claude-opus-4-8
effort: high
hub: psychology
tags:
  - self-efficacy
  - bandura
  - self-esteem
  - impostor-syndrome
  - identity-formation
  - erikson
  - marcia-identity-statuses
  - self-concept
  - self-discrepancy
  - ego-threat
  - self-affirmation
  - psychological-safety
  - confidence
related_skills:
  - psychology
  - applied-psychology
  - psychology-positive
  - psychology-clinical-personality
  - psychology-social
  - psychology-stress-trauma
whenToUse:
  - "What does research actually say about self-efficacy / Bandura's 4 sources?"
  - "self-efficacy vs self-esteem vs self-confidence — what's the difference, mechanistically?"
  - "does high self-esteem cause success? (Baumeister 2003 vs Orth & Robins 2022)"
  - "sociometer theory, contingencies of self-worth, true vs contingent self-esteem"
  - "impostor syndrome — is it real, how prevalent, and is the construct sound?"
  - "the 'stop calling it impostor syndrome' / systemic-reframe debate"
  - "Erikson's psychosocial stages / identity vs role confusion"
  - "Marcia's identity statuses (diffusion/foreclosure/moratorium/achievement)"
  - "narrative identity (McAdams) vs status model; multicultural critiques of identity theory"
  - "self-concept, William James, Higgins self-discrepancy (actual/ideal/ought)"
  - "ego threat, self-affirmation theory (Steele), defensiveness mechanisms"
  - "psychological safety (Edmondson) — theory, Project Aristotle, and its critiques"
---

# Confidence, Self & Identity Psychology

A spoke of the **psychology** hub. Evidence-grounded, **replication-honest** reference for the
psychology of the self: how confidence is built (self-efficacy), how people value themselves
(self-esteem), how they doubt their own competence (impostor phenomenon), how identity forms
(Erikson/Marcia), how the self is represented (self-concept, self-discrepancy), what happens
when the self is threatened (ego threat, self-affirmation), and the team-level analog of
self-safety (psychological safety).

This skill owns the **theory and mechanism** of the self. "How do I *use* self-efficacy, growth
mindset, or psychological safety to drive motivation, adoption, learning, or trust?" →
**applied-psychology** (the hub's governing rule: *"how does it work / what does research say?"*
→ psychology; *"how do I use it?"* → applied-psychology).

This is one of psychology's **most marketing-distorted** corners. The self-esteem movement,
"impostor syndrome" as a pop diagnosis, and "psychological safety" as a management buzzword all
ran far ahead of the evidence. A genuine expert reference researches the **takedowns and
corrections** as hard as the theories — and that is what the §Replication sections below do.

---

## Quick Reference — Decision Table

| Construct | Core claim | Status / effect | Primary caveat |
|---|---|---|---|
| **Self-efficacy** (Bandura 1977/1997) | Domain-specific belief "I can execute *this task*" predicts effort, persistence, performance | **Robust.** Stajkovic & Luthans 1998: r≈.38 with work performance; Multon et al. 1991: ~14% of academic-achievement variance | Specificity matters — global "confidence" measures predict far worse; partly reverse-causal (past success builds efficacy) |
| **4 sources of efficacy** | Mastery > vicarious > verbal persuasion > physiological/affective | Ranking **replicates** repeatedly; mastery experience dominant | Source weights vary by domain; not strictly ordinal everywhere |
| **Collective / teacher efficacy** | Group-level efficacy predicts group performance | Supported, esp. teacher efficacy & education | Aggregation/measurement debates; can be inflated by shared method variance |
| **Self-esteem** (global self-worth) | Feeling good about oneself causes success, health, happiness | **Largely overstated.** Baumeister et al. 2003: benefits modest & mostly correlational; 2022 Orth-Robins/Krueger debate unsettled | Boosting self-esteem directly does **not** reliably cause outcomes; effects domain-dependent |
| **Sociometer theory** (Leary) | Self-esteem is a *gauge* of social inclusion, not a cause of outcomes | Well-regarded reframe; explains why approval moves everyone | A monitor, not a lever — "raising" the gauge needn't change the underlying belonging |
| **Contingencies of self-worth** (Crocker) | *What* your worth rests on matters more than its level; pursuit of self-esteem is costly | Influential; pursuit-of-esteem costs replicated in places | Mostly self-report/correlational; level still matters somewhat |
| **Impostor phenomenon** (Clance & Imes 1978) | Persistent self-doubt + fear of being exposed as a fraud despite success | **Real experience, shaky construct.** "70% of people" stat is folklore; measurement unstable | Not a DSM disorder; "syndrome" framing & individual-pathology lens contested (Tulshyan & Burey 2021) |
| **Erikson psychosocial stages** | 8 lifespan crises; adolescence = identity vs role confusion | Foundational, heuristically rich | **Not falsifiable as a stage theory**; based on mid-century Western males; little hard quantitative support |
| **Marcia identity statuses** | Crossing exploration × commitment → diffusion / foreclosure / moratorium / achievement | The **empirically tractable** operationalization; statuses & developmental drift replicate | Critics: not a faithful operationalization of Erikson; misses temporal continuity & narrative |
| **Narrative identity** (McAdams) | Identity = the internalized, evolving life *story* | Growing, complementary tradition; links to wellbeing (redemption sequences) | Mostly cross-sectional; causal direction underdetermined |
| **Self-discrepancy** (Higgins 1987) | Actual-vs-ideal gap → dejection; actual-vs-ought gap → agitation | Influential mapping of self-gaps to distinct emotions | Specific-emotion predictions **replicate inconsistently**; effects often small |
| **Self-affirmation** (Steele 1988) | Affirming an unrelated valued self-domain reduces defensiveness to threat | Many lab effects; some field successes (education) | **Heterogeneous & moderator-laden**; null/failed replications exist; not a panacea |
| **Psychological safety** (Edmondson 1999) | Shared belief the team is safe for interpersonal risk → learning, error reporting | Robust team-learning link; Project Aristotle popularized | **Construct creep**; cultural baselines differ; misused to excuse low accountability |

---

## 1. Self-Efficacy — the most robust construct here

**Albert Bandura** (1977, *"Self-efficacy: Toward a unifying theory of behavioral change"*;
1997, *Self-Efficacy: The Exercise of Control*) defined **perceived self-efficacy** as belief in
one's capability to **organize and execute the courses of action required to manage prospective
situations** — *"can I do this specific task?"* It grew out of his social-cognitive theory and
the principle of **reciprocal determinism** (person ↔ behavior ↔ environment interact).

**The load-bearing feature is *specificity*.** Self-efficacy is **task- and domain-specific**, not
a global trait. A person can have high math self-efficacy and low public-speaking self-efficacy.
This is the single most common error in casual use: treating self-efficacy as a synonym for
generalized "confidence." Measured at the right level of specificity, it predicts well; measured
globally, it predicts poorly.

### The four sources (the ranking that keeps replicating)

In **descending power**:

1. **Mastery experiences (enactive attainment)** — *direct, successful performance of the target
   behavior.* By far the strongest source; nothing builds the belief "I can do this" like having
   done it. (This is the mechanistic core of why graded, succeed-able challenges build confidence.)
2. **Vicarious experience (modeling)** — seeing similar others succeed ("if they can, so can I").
   Strength rises with model–observer similarity.
3. **Verbal/social persuasion** — credible encouragement; weaker and fragile if unsupported by
   real performance.
4. **Physiological & affective states** — reading one's arousal/anxiety/fatigue as a signal of
   (in)capacity; reappraising arousal as readiness rather than fear matters.

Every subsequent meta-analytic test has reproduced Bandura's **ordinal claim that mastery
experience dominates**, though exact source weights shift by domain.

### Key distinctions to keep straight

- **Efficacy expectancy vs outcome expectancy** — "Can I perform the behavior?" vs "Will the
  behavior produce the outcome?" Both matter; Bandura argued efficacy is usually the binding
  constraint on initiating and persisting.
- **Self-efficacy vs self-esteem** — *capability belief about a task* vs *global self-worth.* They
  are routinely conflated and are **dissociable**: you can have high self-worth and low efficacy
  for a given task, or vice versa.
- **Self-efficacy vs self-confidence** — "self-confidence" is the vague lay term; self-efficacy is
  its operationalized, domain-anchored cousin.
- **Collective efficacy** — a group's shared belief in its conjoint capability; **teacher collective
  efficacy** is among the better-supported group-level versions and a strong correlate of school
  achievement.

### Replication / status — the strong end of the spectrum

- **Stajkovic & Luthans (1998)**, *Psychological Bulletin* — 114 studies (k=157, N≈21,616):
  weighted mean **r ≈ .38** between self-efficacy and work performance, **moderated by task
  complexity** (effect shrinks as tasks get more complex) and setting (lab > field).
- **Multon, Brown & Lent (1991)** — self-efficacy accounted for **~14% of variance** in academic
  achievement and ~12% in persistence: large for a single psychological predictor.
- **The honest caveat (Vancouver et al., 2001–2002).** At the *within-person* level the relationship
  can **reverse**: holding ability constant, *higher* self-efficacy sometimes predicts *worse*
  subsequent performance (over-confidence → reduced effort/complacency). The between-person positive
  correlation is real and robust; the naïve causal story ("just pump up efficacy") is not. Past
  performance is a major common cause — much of the efficacy→performance link is **performance →
  efficacy** running backward. Treat self-efficacy as a strong, specificity-dependent, partly
  reciprocal predictor, **not** a free lever.

---

## 2. Self-Esteem — the most overstated construct here

**Self-esteem** = a person's **global evaluation of their own worth** (Rosenberg's 1965 scale is
the field standard). The popular "self-esteem movement" (1980s–90s, esp. California's self-esteem
task force) assumed raising self-esteem would raise achievement and reduce social ills. The
science did not cooperate.

### The takedown — Baumeister, Campbell, Krueger & Vohs (2003)

This commissioned *Psychological Science in the Public Interest* review is the canonical citation:

- Self-esteem **correlates** with good outcomes but the correlations are **modest** and largely
  **non-causal** — high achievers develop high self-esteem more than high self-esteem produces
  achievement (**reverse/third-variable causation**).
- **No good evidence** that programs boosting self-esteem cause better grades, job performance, or
  health.
- High self-esteem **feels good** and modestly aids initiative and persistence, but it is **not** a
  major cause of life success.
- **Violence/aggression**: contrary to the "low self-esteem causes violence" myth, aggression is
  better predicted by **threatened egotism** — *high but unstable/narcissistic* self-views meeting
  challenge (Baumeister, Bushman & Boden 1996). Distinguish secure high self-esteem from **fragile
  high self-esteem / narcissism** (the latter → defensiveness and aggression; deep clinical
  narcissism → **psychology-clinical-personality**).

### The 2022 re-litigation (flag this — it is unsettled, not settled)

A *Psychological Bulletin* exchange reopened the question:

- **Orth & Robins (2022)** reviewed longitudinal/meta-analytic evidence and concluded self-esteem
  brings **modest but reliable benefits** across relationships, school, work, mental/physical health,
  and (lower) antisocial behavior — i.e., a partial pushback on the 2003 nihilism.
- **Krueger, Vohs & Baumeister (2022)** replied that for **objective behavioral** outcomes the
  effects remain **variable and domain-dependent** — "feeling good without doing good"; the allure
  is "a matter of mind and memory, not behavior."
- **Net for an expert reference:** the *strong* causal claim of the self-esteem movement is dead.
  Whether self-esteem yields **modest** downstream causal benefits is **actively contested as of the
  mid-2020s**. Present both sides; do not assert either pole as settled.

### Three reframes that explain *why* the lever fails

- **Sociometer theory (Leary & Baumeister 2000).** Self-esteem is an internal **gauge of relational
  value / social inclusion**, evolved to track belonging. Artificially "raising the gauge" needn't
  change the underlying social reality, which is *why* self-esteem interventions disappoint. Leary's
  data suggest disapproval moves everyone's self-esteem — even people who **report** being immune.
- **Contingencies of self-worth (Crocker & colleagues).** What matters is **what your worth is staked
  on** (appearance, others' approval, virtue, academic competence, God's love…). Basing worth on
  *external/validation* contingencies is a chronic vulnerability; the **pursuit of self-esteem** has
  costs to learning, autonomy, relationships, and self-regulation. *Level* of self-esteem matters
  less than its *basis*.
- **True vs contingent self-esteem (Deci & Ryan, SDT).** "True" self-esteem grows organically from
  **autonomous, competent action in authentic relationships** (the SDT needs — see
  **psychology-positive** for SDT-as-wellbeing); **contingent** self-esteem rises and falls with
  ego-involving outcomes and is associated with fragility and ill-being.

---

## 3. The Impostor Phenomenon — real experience, contested construct

**Pauline Clance & Suzanne Imes (1978)** coined the **impostor phenomenon** (note: *phenomenon*,
not "syndrome") to describe high-achieving women who, despite objective success, **dismiss it as
luck/charm/error** and live in fear of being **exposed as frauds**. Clance later built the **Clance
Impostor Phenomenon Scale (CIPS)**, still the dominant measure (alongside Harvey's IP scale and
Leary's). It is **not a clinical diagnosis** — it appears in no DSM/ICD.

### What is solid vs what is folklore

- **Solid:** the *experience* is genuine and common; it correlates reliably with **anxiety,
  depression, low self-esteem, perfectionism, neuroticism, and burnout**, especially in evaluative,
  high-stakes, minority-status contexts.
- **Folklore:** the ubiquitous **"70% of people experience impostor syndrome"** stat is **not a
  rigorous population estimate** — it traces to a loose secondary claim, not a representative study.
  Prevalence figures swing wildly (systematic reviews report anywhere from ~9% to ~82%) precisely
  **because measurement is inconsistent** — different scales, different cutoffs, mostly convenient
  samples (students, clinicians). Cite ranges, not the round number.

### The construct critique (flag this prominently)

- **Measurement instability** — no agreed cutoff, weak consensus on factor structure; many studies
  treat a *continuous, normal* experience as if it were a discrete condition. Recent work (e.g.
  Walker & Saklofske 2023) is still trying to build psychometrically sound instruments — itself a
  signal the construct was under-validated for decades.
- **The pathologizing critique — Tulshyan & Burey (2021, HBR, "Stop Telling Women They Have
  Impostor Syndrome").** The concept was developed without accounting for **systemic racism,
  classism, and gender bias**; it took a near-universal discomfort and **pathologized it as an
  individual defect**, disproportionately aimed at women and people of color. For many, feeling like
  an outsider **isn't a distortion — it's an accurate read of an exclusionary environment.** The
  reframe: fix the **environment** (which legitimizes diverse leadership styles), not the individual.
- **Treatment gap** — systematic reviews find **almost no rigorously evaluated treatments**
  specifically for impostor feelings; most "interventions" are educational/normalizing.

**Expert framing:** describe it as a **real, distressing, common experience** that is *partly an
internal cognitive pattern and partly a signal of genuinely non-inclusive contexts* — and resist
both the pop-diagnostic "syndrome" framing and the inflated prevalence claims.

---

## 4. Identity Formation — Erikson and Marcia

### Erikson's psychosocial theory

**Erik Erikson** (1950, *Childhood and Society*; 1968, *Identity: Youth and Crisis*) recast
development as **eight lifespan psychosocial crises**, each a tension to resolve:

1. Trust vs Mistrust (infancy) · 2. Autonomy vs Shame/Doubt · 3. Initiative vs Guilt ·
4. Industry vs Inferiority · 5. **Identity vs Role Confusion (adolescence)** · 6. Intimacy vs
Isolation · 7. Generativity vs Stagnation · 8. Ego Integrity vs Despair.

Stage 5 is the namesake: adolescents must integrate childhood identifications into a coherent,
continuous sense of **who they are** — Erikson's hallmark is the **sense of temporal-spatial
continuity** ("the same person across time and contexts"). He introduced **identity crisis** and the
**psychosocial moratorium** (a sanctioned period of role experimentation).

> **Status caveat.** Erikson's framework is **foundational and heuristically powerful but not a
> falsifiable, well-quantified stage theory.** The fixed sequence and discrete crises have **little
> hard empirical confirmation**; the model was built from clinical/biographical work on **mid-20th-
> century Western, largely male** lives. Treat it as a generative lens, not validated mechanism.

### Marcia's identity statuses — the operationalization

**James Marcia (1966)** made Erikson's identity testable by crossing **two dimensions**:

- **Exploration** — has the person actively considered alternatives? and
- **Commitment** — has the person settled on values/goals/roles?

Yielding **four identity statuses** (assessed via the **Identity Status Interview**):

| | Low commitment | High commitment |
|---|---|---|
| **Low exploration** | **Diffusion** (no exploration, no commitment — least mature, least stable) | **Foreclosure** (commitment *without* exploration — adopted from parents/authority) |
| **High exploration** | **Moratorium** (actively exploring, not yet committed — the "crisis" state) | **Achievement** (committed *after* exploration — most mature) |

### Replication / status — what holds, what doesn't

- **Holds:** the statuses are recoverable; longitudinal latent-class work confirms **status
  *trajectories*** and supports **Waterman's developmental hypothesis** — more achievers and fewer
  diffusions in middle-to-late vs early-to-middle adolescence. Refinements found **two moratoriums**
  (classic "searching" vs a more troubled/ruminative form), feeding modern **dual-cycle models**
  (Luyckx, Crocetti: *commitment-making, identification-with-commitment, exploration-in-breadth,
  -in-depth, ruminative exploration*).
- **Doesn't / contested:** critics argue the status model is **not a faithful operationalization of
  Erikson** — it captures exploration/commitment but **drops Erikson's core sense of temporal
  continuity** and treats identity as a *static category* rather than a *developing structure*.
- **Multicultural & gender critique:** the achievement→good, foreclosure→bad valuation reflects a
  **Western individualistic** assumption; in collectivist cultures or contexts with **structured
  rites of passage**, "foreclosure" (committing via family/community without solo exploration) may be
  adaptive and normative. Early gender work also misfit women's identity-via-relationships pathways.

### Narrative identity — the complementary tradition

**Dan McAdams** argues identity is the **internalized, evolving life story** a person authors to
give life unity and purpose — a *level of personality* above traits and characteristic adaptations.
Status (categorical) and narrative (story-structural) approaches are **complementary**, both
descended from Erikson. Narrative work links **redemptive sequences** (suffering → growth) and
**agency/communion** themes to wellbeing — but is largely **cross-sectional**, so causal direction
is underdetermined.

---

## 5. Self-Concept & Self-Discrepancy

### The self-concept lineage

**William James (1890)** split the self into the **"I"** (the self-as-knower, the stream of
experience) and the **"Me"** (the self-as-known: material, social, and spiritual selves) and gave the
field its first **self-esteem formula**: *self-esteem = successes / pretensions* — you raise it by
succeeding **or** by lowering your aspirations. **Self-concept** is the broader **cognitive
representation** of who one is (multiple, context-dependent self-schemas); **self-esteem** is the
**evaluative** charge on it; **identity** is its socially-situated, continuity-bearing organization.
Keep the three distinct: *representation* (self-concept), *evaluation* (self-esteem), *organization/
continuity* (identity).

### Higgins' self-discrepancy theory (1987)

**E. Tory Higgins** mapped **gaps between self-states to distinct emotions**:

- Three **self-domains**: **actual** (who I am), **ideal** (who I *want* to be — hopes/aspirations),
  **ought** (who I *should* be — duties/obligations).
- Two **standpoints**: own and significant-other's.
- **Predictions:** **actual–ideal** discrepancy → absence of positive outcomes → **dejection**
  emotions (disappointment, sadness, dejection); **actual–ought** discrepancy → presence of negative
  outcomes → **agitation** emotions (guilt, anxiety, self-contempt). Ties to his later
  **promotion vs prevention** regulatory-focus theory (ideals→promotion/gains, oughts→prevention/
  safety).

> **Replication caveat.** The *broad* idea (self-gaps cause distress) is well-supported, but the
> theory's **signature specific prediction** — that ideal-discrepancies selectively cause *sadness*
> while ought-discrepancies selectively cause *anxiety* — has **replicated inconsistently**; several
> studies find both discrepancies predict general distress without the clean emotional
> dissociation, and effects are often small. Cite the framework for its mapping logic, not as a
> precise emotional dial.

---

## 6. Ego Threat & Self-Affirmation

**Ego threat** = information or events implying personal **inadequacy** (failure feedback, social
rejection, stereotype-relevant challenge) that menace the global sense of self-worth/integrity and
trigger **defensiveness** — biased information processing, self-justification, derogating the source,
sometimes aggression (the *threatened-egotism* route to violence, §2).

**Self-affirmation theory — Claude Steele (1988).** People are motivated to maintain **global
self-integrity** (a sense of being *adequate, moral, and in control*) — note this is *global*, so a
threat in one domain can be neutralized by **affirming an unrelated valued domain** (values, roles,
relationships). Writing about a cherished value before a threat reliably **reduces defensiveness** and
opens people to threatening-but-true information, with documented downstream effects on health
messaging and education (the well-known "values-affirmation" classroom interventions).

> **Replication caveat (important).** Self-affirmation effects are **real but heterogeneous and
> heavily moderated** — they work best for people **higher in the threatened domain / under genuine
> threat**, can be **null or backfire** otherwise (e.g., affirming the *same* domain as the threat;
> affirming low-resource individuals can demotivate). Several lab and field **failures to replicate**
> exist, and early education results were sometimes larger than later, better-powered ones. Treat
> self-affirmation as a **context-dependent buffer**, not a robust universal intervention.

---

## 7. Psychological Safety — the team-level analog

**Amy Edmondson (1999)** defined **team psychological safety** as a **shared belief that the team is
safe for interpersonal risk-taking** — that speaking up, asking questions, admitting errors, or
proposing wild ideas won't be punished or humiliated. Her hospital study produced the field's iconic
counterintuitive finding: **better** teams reported **more** medication errors — not because they
erred more, but because they felt **safe enough to surface** them. Safety predicts **team learning
behavior**, which predicts **performance**.

**Project Aristotle (Google, 2012–2015)** studied 180+ teams and reported **psychological safety as
the single largest predictor** of team effectiveness, ahead of dependability, structure/clarity,
meaning, and impact — which catapulted the concept into mainstream management.

> **Replication / critique (flag this).**
> - **Construct creep** (Newman, Donohue & Eva 2017) — "psychological safety" has been stretched to
>   mean almost any nice climate, blurring it with trust, cohesion, and engagement; measurement and
>   discriminant validity suffer.
> - **It is not the same as comfort or low standards** — Edmondson is emphatic that safety must pair
>   with **high accountability**; the dangerous misuse is invoking "psych safety" to excuse
>   **low-performance** cultures. The 2×2 is safety × standards: high-high = *learning zone*.
> - **Project Aristotle is industry research, not a controlled study** — single-company, internally
>   analyzed, never peer-reviewed; its "single largest predictor" headline is **suggestive, not
>   established causal fact**, and sits inside social psychology's broader replication problems.
> - **Cultural baselines differ** — the scale's anchors and speaking-up norms vary across national
>   and organizational cultures, complicating cross-context comparison.
>
> Net: the **core link (safety → learning/error-reporting → performance) is well-supported**; the
> **buzzword inflation around it is not.**

---

## Boundary / SKIP routing

| If the question is really about… | Route to |
|---|---|
| **How to *use* confidence / growth mindset / grit / flow / psych safety to drive motivation, adoption, learning, or trust** | **applied-psychology** (operator lens — the hub's "how do I use it?" rule) |
| **Growth mindset (Dweck) & grit (Duckworth) replication depth** | **applied-psychology** (performance & resilience reference) — small/contested effects flagged there |
| **Happiness, PERMA, eudaimonia, meaning, or SDT *as wellbeing theory*** | **psychology-positive** |
| **Burnout, chronic stress, allostatic load, cortisol** | **psychology-stress-trauma** |
| **Narcissism / fragile high self-esteem *as a personality disorder*, attachment, object relations** | **psychology-clinical-personality** |
| **Status, dominance vs prestige, social comparison, in-group identity** | **psychology-social** |
| **Trust/rapport/psychological safety *applied to a customer or stakeholder*** | **applied-psychology** (trust & psychological-safety reference) |
| **Persuasion, compliance, influence theory** | **psychology-influence-depth** |

**The governing rule:** *"How does it work / what does research say?"* → **psychology** (this spoke).
*"How do I use it on a person/team/customer?"* → **applied-psychology**.

---

## Replication scorecard (cite-honest summary)

| Strong / robust | Mixed / contested | Fragile / overstated / corrected |
|---|---|---|
| Self-efficacy → performance (specificity-dependent; r≈.38 work) | Self-esteem causal benefits (2022 Orth-Robins vs Krueger/Vohs — open) | The self-esteem *movement*'s strong causal claim (Baumeister 2003) |
| 4-sources ranking (mastery dominant) | Self-affirmation (real but moderator-heavy; some null) | "70% have impostor syndrome" (folklore, not a population estimate) |
| Psych safety → team learning → performance | Project Aristotle's "biggest predictor" headline (industry, not peer-reviewed) | Higgins' *specific* ideal=sadness / ought=anxiety dissociation (inconsistent) |
| Marcia statuses recover + developmental drift | Narrative identity ↔ wellbeing (cross-sectional, causal direction open) | "Low self-esteem causes violence" (it's *threatened egotism*) |
| Sociometer & contingencies-of-self-worth logic | Within-person self-efficacy can *reverse* (Vancouver) | Erikson stages as a *falsifiable, validated* stage theory |

---

## Core sources

- **Bandura, A. (1977, 1997)** — self-efficacy theory; *Self-Efficacy: The Exercise of Control*.
- **Stajkovic, A. & Luthans, F. (1998)**, *Psych. Bulletin* — self-efficacy ↔ work performance meta (r≈.38).
- **Multon, Brown & Lent (1991)** — self-efficacy ↔ academic outcomes meta (~14% variance).
- **Vancouver et al. (2001, 2002)** — within-person reversal / reciprocal causation caveat.
- **Baumeister, Campbell, Krueger & Vohs (2003)**, *PSPI* — "Does High Self-Esteem Cause…?" takedown.
- **Orth & Robins (2022)** + **Krueger, Vohs & Baumeister (2022)**, *Psych. Bulletin* — the modern debate.
- **Leary & Baumeister (2000)** — sociometer theory.
- **Crocker & colleagues (2003, 2011)** — contingencies of self-worth; **Deci & Ryan** — true vs contingent.
- **Clance & Imes (1978)** — impostor phenomenon; **CIPS**; **Walker & Saklofske (2023)** — measurement.
- **Tulshyan & Burey (2021)**, *HBR* — "Stop Telling Women They Have Impostor Syndrome" (systemic reframe).
- **Erikson, E. (1950, 1968)** — psychosocial stages; identity vs role confusion.
- **Marcia, J. (1966)** — identity statuses; **Waterman**, **Luyckx**, **Crocetti** — dual-cycle refinements.
- **McAdams, D.** — narrative identity.
- **James, W. (1890)** — I/Me, self-esteem = success/pretensions.
- **Higgins, E.T. (1987)** — self-discrepancy theory; regulatory focus.
- **Steele, C. (1988)** — self-affirmation theory / self-integrity.
- **Edmondson, A. (1999)** — psychological safety & team learning; **Newman, Donohue & Eva (2017)** — construct critique; **Google Project Aristotle (2015)**.
