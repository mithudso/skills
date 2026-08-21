---
name: psychology-stress-trauma
description: >-
  Stress & trauma neuroscience (psychology spoke): the replication-honest science of how stress
  and trauma act on body and brain, separating robust mechanism from fragile claim.
  TRIGGER: how chronic stress/cortisol/the HPA axis harms the body; allostatic load (McEwen);
  Sapolsky's glucocorticoid cascade; PTSD neurobiology / why cortisol is sometimes LOW (Yehuda);
  hippocampal damage vs pre-existing vulnerability (Gilbertson twins); ACEs (Felitti/Anda); burnout
  vs depression (Maslach, ICD-11); resilience as the modal response (Bonanno); polyvagal theory /
  vagal tone (Porges vs Grossman & Taylor); fragile/harmful claims (CISD debriefing, trauma
  epigenetics, "Body Keeps the Score") vs PE/CPT/EMDR. SKIP: coping with or preventing burnout,
  recovery, emotion regulation -> applied-psychology; ADHD/autism ->
  psychology-neurodevelopmental; PERMA/post-traumatic growth -> psychology-positive; personality
  disorders/attachment -> psychology-clinical-personality; institutional betrayal / DARVO / survivor disbelief -> psychology-institutional-betrayal; diagnosis/medical advice -> a clinician.
version: "1.0.1"
updated: "2026-06-23"
category: spoke
model: claude-opus-4-8
effort: high
tags:
  - psychology
  - stress
  - trauma
  - PTSD
  - allostatic-load
  - HPA-axis
  - cortisol
  - neuroendocrinology
  - polyvagal-theory
  - vagal-tone
  - burnout
  - resilience
  - ACEs
  - replication
related_skills:
  - psychology
  - applied-psychology
  - psychology-clinical-personality
  - psychology-neurodevelopmental
  - psychology-positive
whenToUse:
  - "how chronic stress / cortisol / the HPA axis acts on body and brain (mechanism, not coping)"
  - "what is allostasis vs homeostasis, allostatic load, allostatic overload, McEwen"
  - "the acute stress response — SAM/catecholamine arm + HPA arm, CRH→ACTH→cortisol, MR/GR feedback, diurnal rhythm / cortisol awakening response"
  - "PTSD neurobiology — why cortisol is sometimes LOW (Yehuda), amygdala/hippocampus/mPFC fear circuit, fear conditioning & extinction"
  - "the hippocampal-volume cause-vs-pre-existing-vulnerability debate (Gilbertson monozygotic twins)"
  - "ACEs / adverse childhood experiences dose-response and its limits (Felitti & Anda)"
  - "burnout science — Maslach three dimensions, ICD-11 occupational phenomenon, burnout-vs-depression boundary"
  - "trauma resilience — Bonanno's trajectory model, resilience as the modal (not rare) response"
  - "polyvagal theory, vagal tone, RSA/HRV as a stress biomarker — and whether the physiology holds up (Grossman & Taylor)"
  - "evidence-based PTSD treatment (PE, CPT, EMDR) and which popular trauma claims are fragile or harmful"
  - "Sapolsky, glucocorticoid-cascade hypothesis, 'Why Zebras Don't Get Ulcers'"
  - "epigenetics of trauma / intergenerational transmission — and why it's contested"
---

# Stress & Trauma Neuroscience

Mechanistic, replication-honest reference for the biology and psychology of stress and trauma. This is the **"how does it work / what does the evidence actually show"** spoke of the `psychology` hub. It deliberately separates *robust mechanism* from *popular-but-fragile claim* — the stress/trauma literature is unusually prone to clinically beloved theories whose physiology is weak (polyvagal), to intuitions that the data invert (cortisol in PTSD), and to interventions that feel helpful but aren't (debriefing).

> **Boundary.** This spoke owns the *science*. The operator/coping angle (how to prevent a team from burning out, recover from stress, build resilience habits) lives in `applied-psychology` → `references/performance-and-resilience-psychology`. Emotion-regulation *technique* → `applied-psychology` → `references/emotion-and-affect-psychology`. Rule of thumb: **"how does it work / what's the evidence?" → here. "what do I do about it?" → applied-psychology.**

> **Not medical advice.** Educational synthesis of the research literature. Diagnosis and treatment require a licensed clinician.

---

## 1. The stress response has two arms (don't omit the fast one)

A real or *interpreted* threat to homeostasis triggers two coordinated outputs. Cortisol gets all the attention, but the fast arm fires first.

1. **SAM axis (sympathetic-adrenomedullary)**: milliseconds-to-seconds. The sympathetic nervous system and adrenal medulla release **catecholamines** (epinephrine, norepinephrine): tachycardia, raised blood pressure, glucose mobilization, the "fight-or-flight" surge. This is the arm popular accounts forget.
2. **HPA axis (hypothalamic-pituitary-adrenal)**: minutes-to-hours. Slower, hormonal, glucocorticoid-mediated, and self-limiting via negative feedback.

Allostatic-load indices integrate **both** arms (catecholamine excretion *and* cortisol), which is why a cortisol-only picture of "stress" is incomplete.

### The HPA cascade in detail
- **CRH** (corticotropin-releasing hormone) from the hypothalamic **paraventricular nucleus (PVN)** → portal vessels → anterior pituitary.
- CRH stimulates **ACTH** (adrenocorticotropic hormone) release from pituitary corticotrophs → systemic circulation → adrenal **cortex**.
- ACTH drives synthesis/secretion of **cortisol** (the primary human glucocorticoid; corticosterone in rodents).
- **Negative feedback**: cortisol inhibits CRH and ACTH at the hypothalamus and pituitary (and is modulated by the hippocampus and PFC), closing the loop. Dysregulated feedback, too much *or* too little, is the through-line of stress pathology.

### Two receptors, not one (and they matter)
Glucocorticoids act through two intracellular receptors with different affinities:
- **MR (mineralocorticoid / type-1 receptor)**: **high** affinity. Largely occupied at **basal/trough** cortisol; sets tonic HPA tone and baseline regulation.
- **GR (glucocorticoid / type-2 receptor)**: **low** affinity. Recruited as cortisol **rises** (stress, circadian peak); drives the additional inhibitory feedback and most stress-response gene effects.

The MR:GR *balance* is a better organizing concept than "cortisol level." This is also why PTSD's signature is best described in **receptor-sensitivity** terms (below), not in raw cortisol.

### Rhythm, not a set-point
Cortisol is secreted in **circadian** (24-h) and **ultradian** (pulsatile, ~hourly) rhythms. The **cortisol awakening response (CAR)**, a ~50%+ rise within ~30–40 min of waking, is a distinct, circadian-driven phenomenon thought to prime the day. Flattened diurnal slope and blunted CAR are markers of dysregulation. Practical upshot: a single spot cortisol is nearly uninterpretable without time-of-day context.

---

## 2. Allostasis & allostatic load (the organizing framework)

- **Homeostasis**: keeping a few life-critical variables (pH, O₂, core temp) in narrow ranges.
- **Allostasis** (Sterling & Eyer, 1988): **"stability through change"**: actively adjusting set-points *in anticipation of* demand, via the ANS, HPA, immune, and metabolic systems. Adaptation, not deviation.
- **Allostatic load** (McEwen & Stellar, 1993): the cumulative **physiological cost** of repeated allostatic activation: the wear-and-tear from chronic up-/down-regulation of stress mediators.
- **Allostatic overload**: the tipping point where load exceeds coping capacity and drives disease.

McEwen's four routes to load: (1) **frequent** stress hits; (2) **failure to habituate** to repeated stressors; (3) **failure to shut off** the response after the stressor ends; (4) **inadequate response** in one system causing compensatory **over**-activity in another (e.g., low cortisol → unrestrained inflammatory cytokines). Note the symmetry: **too little is as costly as too much.** The mediators that protect acutely (cortisol, catecholamines, cytokines) are the same ones that damage chronically — a dose/duration distinction, not a "good vs bad hormone" one.

**Operationalization caveat (replication-relevant).** The *framework* is robust and widely supported, but the **allostatic load index** (typically a composite of primary mediators like cortisol, DHEA-S, epinephrine, norepinephrine plus secondary outcomes like BP, waist-hip, HbA1c, lipids, CRP) is operationalized **heterogeneously** across studies. Different biomarker panels and cut-points limit cross-study comparability. Treat AL as a strong organizing model whose specific composite scores are study-dependent.

---

## 3. Glucocorticoid-cascade hypothesis (Sapolsky)

Sapolsky's hypothesis: chronic glucocorticoid exposure produces a **feed-forward cascade**: sustained cortisol down-regulates hippocampal GRs → the hippocampus (a brake on the HPA axis) loses regulatory capacity → less feedback inhibition → still more cortisol → hippocampal neuronal damage/atrophy → worse regulation. A vicious cycle linking chronic stress to hippocampal aging and neuropsychiatric disease. Popularized in *Why Zebras Don't Get Ulcers*: the core insight is that humans uniquely activate a system built for acute physical emergencies in response to **psychological and chronic** stressors, so the same machinery that saves a zebra from a lion corrodes a human over decades.

**Replication status — partly revised.** The rodent evidence for glucocorticoid-induced hippocampal remodeling is strong. The strong *human* "neurotoxic cascade" claim has been **re-examined and softened**: hippocampal changes are real and correlate with cortisol and disorder, but causation/directionality in humans is contested (see §5). Hold the cascade as an influential, partly-supported model, not settled human fact.

---

## 4. PTSD neurobiology — and the cortisol surprise

PTSD is the canonical trauma disorder, but its biology **inverts the lay "stress = high cortisol" story.**

### The Yehuda finding: PTSD is often LOW / paradoxical cortisol
Rachel Yehuda's work shows PTSD is characterized not by elevated cortisol but by:
- **Enhanced glucocorticoid-receptor sensitivity** → **exaggerated negative feedback** (super-suppression on the low-dose dexamethasone suppression test — the *opposite* of melancholic depression, which shows feedback resistance/non-suppression).
- **Lower or normal** basal 24-h cortisol and flattened diurnal rhythm (less consistent than the GR-sensitivity finding).
- Increased lymphocyte GR number; greater GR responsiveness.

Mechanistically: central HPA *drive* (CRH) can be elevated while peripheral cortisol stays low because the feedback brake is over-tight. Low cortisol + enhanced GR sensitivity measured *in the acute aftermath of trauma* have been reported as **predictors** of later PTSD, suggesting a pre-trauma diathesis, not just a scar. The practical teaching point: **do not equate trauma or PTSD with high cortisol**: for PTSD the dysregulation more often runs the other way.

### The fear circuit
PTSD maps onto a well-replicated circuit imbalance:
- **Amygdala**: hyper-responsive (threat detection, fear acquisition).
- **Medial prefrontal cortex (esp. vmPFC) / anterior cingulate**: hypo-active (impaired fear *extinction* and top-down regulation — fear is learned fine but not unlearned).
- **Hippocampus**: implicated in context coding; deficits → failure to contextualize threat ("the danger is *here/now*" generalization), plus volume findings (§5).

Trauma memory is well-described by **Pavlovian fear conditioning**: a conditioned stimulus paired with the traumatic unconditioned stimulus elicits conditioned fear; **extinction** is new inhibitory learning, not erasure — which is exactly why exposure-based therapy (§8) works and why reminders can reinstate fear.

---

## 5. Hippocampal volume: damage, or pre-existing vulnerability? (an honest unsettled debate)

Smaller hippocampal volume is one of the most replicated neuroimaging correlates of PTSD. The *interpretation* is genuinely contested — flag this rather than asserting either side.

- **Damage view** (consistent with Sapolsky): trauma/glucocorticoids shrink the hippocampus.
- **Vulnerability view** (**Gilbertson et al., 2002**, *Nature Neuroscience*): in monozygotic twins **discordant for combat exposure**, the *unexposed, never-traumatized* co-twins of combat-PTSD veterans **also had smaller hippocampi**, scaling with the affected twin's PTSD severity. Reading: small hippocampus is largely a **familial/pre-existing risk factor**, not purely a neurotoxic scar.
- **But not settled the other way either:** later twin/longitudinal work (e.g., environment-vs-genetics decompositions) has argued the **environment contributes more than genetics** to the smaller hippocampal volume in PTSD, reopening a role for trauma-related effects. Likely reality: **both** a heritable diathesis *and* trauma-related change, in proportions still being worked out.

Correct expert stance: "smaller hippocampus is a robust *correlate*; whether it is cause, consequence, or both is unresolved, with strong twin evidence for a pre-existing component."

---

## 6. ACEs — robust epidemiology, fragile at the individual level

The **Adverse Childhood Experiences (ACE) Study** (Felitti, Anda et al., CDC–Kaiser, 1998; ~17,000 adults) established a **graded dose-response**: more ACE categories → higher adult risk of heart disease, diabetes, COPD, depression, substance use, suicide attempts, and early mortality. One of public health's most influential and replicated population findings.

**Where to be careful (the "Critical Assessment at 20 Years," *Am J Prev Med* 2019):**
- The original sample was **unrepresentative** (largely white, insured, middle-class, college-educated, older HMO members) — narrowing generalizability.
- The ACE **checklist is narrow** (10 categories) — it under-counts adversity (omits poverty, community violence, racism, bullying) and **mis-represents the social distribution** of harm.
- Most important: the score is an **epidemiologic risk gradient, not an individual diagnostic/predictive instrument.** Using an ACE score to *screen and predict* an individual's outcomes is criticized — most high-ACE individuals do **not** develop the modeled outcomes (regression to population base rates), so individual screening risks both false alarms and determinism. Population truth ≠ individual prophecy.

---

## 7. Burnout — real, but not a tidy disease

- **Maslach's three dimensions** (the dominant model, operationalized by the MBI): (1) **emotional exhaustion**; (2) **depersonalization / cynicism**; (3) **reduced personal accomplishment / efficacy**.
- **ICD-11** lists **burn-out** as an **"occupational phenomenon," explicitly NOT a medical condition**: filed under "factors influencing health status," and **scoped to the occupational context only** (don't extend it to general life fatigue).
- **Burnout-vs-depression boundary debate (unresolved):** exhaustion-dominant burnout overlaps heavily with depression on symptoms and biomarkers; some researchers argue burnout is not cleanly separable from depression/anxiety, while others defend a distinct, work-contextualized construct. Treat the boundary as **genuinely contested.**
- Burnout's physiology is **not** a clean "high-cortisol" story; findings are mixed (some hypo-cortisol, flattened CAR), consistent with the §2 point that chronic dysregulation runs in either direction. A useful (replication-honest) lens is **Job Demands–Resources (JD-R)**: strain rises when demands chronically outstrip resources/recovery.

---

## 8. Resilience is the modal response — and what actually treats trauma

### Bonanno: resilience is common, not rare
**George Bonanno's** prospective trajectory work reframes the field. Across many potentially-traumatic-event studies, four trajectories recur, and **resilience is the *modal* (most common) outcome**:
- **Resilience**: stable healthy functioning (~**two-thirds**; ~65% pooled across studies).
- **Recovery**: initial elevation then return to baseline (~20%).
- **Chronic**: persistent dysfunction/PTSD (~10%).
- **Delayed**: symptoms emerging later (~9%).

Key methodological point: these proportions only emerge from **prospective** designs (assessing people *before*, or population-representatively, not only treatment-seekers). Convenience/clinical samples massively over-estimate chronic dysfunction. Implication: trauma exposure ≠ trauma disorder; most people are resilient, and pathologizing the majority (or assuming everyone needs intervention) is itself a documented error (see CISD below).

### Evidence-based PTSD treatment
Strongest evidence (APA and VA/DoD clinical practice guidelines) is for **trauma-focused psychotherapies**:
- **Prolonged Exposure (PE)**: first-line; repeated imaginal/in-vivo exposure drives extinction learning (§4).
- **Cognitive Processing Therapy (CPT)**: first-line; restructures trauma-related cognitions.
- **Trauma-focused CBT**: strong support.
- **EMDR**: effective and recommended, but typically graded **somewhat lower / second-tier** in the most recent APA guideline (more outcome variability; the eye-movement component's specific contribution is debated — many argue it works *via* exposure).

---

## 9. Polyvagal theory & vagal tone — clinically beloved, physiologically contested

This is the **single most important replication call in this domain.** Polyvagal theory (PVT) is enormously popular in trauma therapy and somatic circles; its core physiology is **largely rejected by comparative physiologists.** Represent both.

### What Porges' theory claims
- A **three-circuit hierarchy** of autonomic response, evolutionarily layered: (1) ancient **dorsal vagal** (unmyelinated) → immobilization/shutdown/"freeze-collapse"; (2) **sympathetic** → mobilization/fight-flight; (3) newest **ventral vagal** (myelinated) → the **social engagement system** (safety, connection, calm).
- The **social engagement system**: an evolved linkage between the myelinated vagus (heart) and the cranial nerves controlling face/voice/head muscles — so safety, facial expression, and prosody co-regulate.
- **Neuroception**: subconscious, continuous appraisal of cues as safe / dangerous / life-threatening, gating which autonomic circuit is active.
- **Respiratory sinus arrhythmia (RSA)** is offered as a readout of **"vagal tone."**

### Why physiologists reject the core premises (Grossman & Taylor; Grossman 2023)
- **"The myelinated cardiac vagus is a uniquely *mammalian* evolutionary innovation"** — **false.** Myelinated cardiac vagal fibers and a dorsal/ventral motor distinction are documented in **non-mammalian** vertebrates (including some teleost/air-breathing fishes in Taylor's comparative work). The evolutionary/comparative-anatomy keystone of PVT is contradicted by the comparative literature.
- **The clean "dorsal = shutdown / ventral = social-safety" functional split** is an oversimplification not supported by vagal neuroanatomy/physiology.
- **RSA as a clean index of "vagal tone"** is disputed: Porges' camp insists RSA must *not* be corrected for respiration; Grossman & Taylor show RSA is strongly **co-modulated by respiratory parameters** (rate, depth), so uncorrected RSA conflates breathing with cardiac vagal control. ("Vagal tone" via RSA is not the unambiguous quantity the theory needs.)
- Net: a 2023 critique frames the five basic premises as facing **"fundamental challenges and likely refutations."** Defenders reply that PVT is a generative clinical *heuristic* even if the physiology is wrong — but a heuristic's clinical usefulness is **not** evidence its mechanistic claims are true.

### The honest position on vagal tone / HRV
Strip away PVT and a **real, modest** signal remains: **HRV** (esp. vagally-mediated metrics like RMSSD, HF-HRV) is **associated** with stress, depression, and anxiety — reduced HRV correlates with these states and with cardiovascular risk. But (umbrella reviews): the association is **non-specific** (HRV shifts with many conditions, meds, fitness, age, breathing), modest in effect, and **not a validated diagnostic biomarker** for any psychiatric disorder. HRV **biofeedback** shows small-to-moderate benefit for stress/anxiety in meta-analyses. So: vagal/HRV effects are **real but oversold**, and they do **not** validate polyvagal theory's anatomy or its three-state map.

**How to talk about it:** "Polyvagal theory is a popular clinical framework whose specific evolutionary and physiological claims are contested or refuted by comparative physiologists (Grossman & Taylor). The underlying observation that vagal/parasympathetic activity (indexed imperfectly by HRV/RSA) tracks stress and emotion is real but non-specific. Use 'autonomic/vagal regulation' language; avoid stating PVT's mechanisms as established fact."

---

## 10. Replication scorecard (cite this when challenged)

| Claim / construct | Status | What's solid vs fragile |
|---|---|---|
| Two-arm acute stress response (SAM + HPA); CRH→ACTH→cortisol; MR/GR feedback; diurnal/CAR rhythm | **Robust** | Core endocrinology, well established |
| Allostasis / allostatic load **framework** (McEwen; Sterling & Eyer) | **Robust** | Framework strong; the AL **index** is operationalized heterogeneously — composite scores are study-dependent |
| Glucocorticoid-cascade hypothesis (Sapolsky) | **Partly revised** | Rodent hippocampal remodeling solid; strong human "neurotoxic cascade" softened/contested |
| PTSD = **enhanced GR sensitivity / exaggerated feedback**, often LOW cortisol (Yehuda) | **Robust (and counter-intuitive)** | GR-sensitivity/super-suppression replicates; hypocortisolemia less consistent. Inverts the lay "high cortisol" story |
| PTSD fear circuit (amygdala↑, vmPFC↓, hippocampus) | **Robust** | Well-replicated; extinction-learning model underlies PE |
| Smaller hippocampus in PTSD: cause vs **pre-existing vulnerability** | **Unsettled** | Correlation robust (Gilbertson twins → strong vulnerability evidence); later work argues environment contributes too. Don't assert either alone |
| ACEs dose-response (Felitti/Anda) | **Robust at population level / fragile at individual** | Graded epidemiologic risk replicates; original sample unrepresentative, checklist narrow, **individual screening/prediction criticized** |
| Burnout (Maslach 3-D; ICD-11 occupational phenomenon) | **Construct accepted; boundaries contested** | ICD-11 = phenomenon, **not** a disease; burnout-vs-depression separability genuinely debated |
| Resilience as **modal** trauma response (Bonanno) | **Robust (with prospective designs)** | Replicates in prospective samples; clinical/convenience samples over-state chronicity |
| PE / CPT / trauma-focused CBT for PTSD | **Robust** | Strongest guideline-recommended evidence |
| EMDR for PTSD | **Effective; lower-graded** | Works, but more variable; eye-movement-specific mechanism debated (likely exposure) |
| **Polyvagal theory** core premises (Porges) | **Contested / largely refuted in physiology** | Myelinated-vagus-uniquely-mammalian claim is **false**; dorsal/ventral functional split + RSA-as-vagal-tone disputed (Grossman & Taylor 2007; Grossman 2023). Popular clinically ≠ mechanistically valid |
| HRV / "vagal tone" as a **diagnostic** stress/mental-health biomarker | **Real but non-specific & oversold** | Associations replicate but are non-specific; not a validated diagnostic; biofeedback small-to-moderate |
| **CISD / psychological debriefing** (one-session) | **Ineffective, possibly HARMFUL** | Cochrane (2002) + meta-analyses: no benefit; some studies show **worse** PTSD vs control. Don't recommend single-session debriefing |
| **Transgenerational epigenetic** trauma transmission (Yehuda FKBP5) | **Highly contested** | Small samples (n≈22 offspring, n≈9 controls), confounds (in-utero, cultural/genetic), **awaits independent confirmation**; transgenerational inheritance in mammals broadly disputed. Flag hard |
| **Triune-brain** ("reptilian/mammalian/human") framing | **Obsolete** | Discredited neuroanatomy; underpins much pop-trauma writing — don't use it as mechanism |
| *"The Body Keeps the Score"* (van der Kolk) specific claims | **Popular, partly overstated** | Peer-reviewed critique (BJPsych Bulletin): some neurobio/treatment claims outrun the evidence, rely on cross-sectional data and outdated models (triune brain). Useful as cultural touchstone, **not** a citation of record |

---

## 11. Practical applications (mechanism → use)

- **Auditing trauma/wellness content for over-claiming.** Flag: triune-brain language, "trauma is stored in the body" stated as settled mechanism, polyvagal "ventral/dorsal state" maps presented as established neuroscience, "stress = high cortisol" (often inverted in PTSD), and single-session debriefing recommendations.
- **Interpreting cortisol data.** Insist on **time-of-day** and **rhythm** (diurnal slope, CAR), not spot values; remember MR/GR balance and that **low** cortisol can signal dysregulation (e.g., PTSD-type enhanced feedback) — not "low stress."
- **Reading HRV / "vagal tone" claims.** Treat as a **non-specific** autonomic correlate, not a diagnosis or proof of polyvagal theory. Useful for trend-tracking and biofeedback; not a biomarker of any specific disorder.
- **Framing trauma outcomes honestly.** Expect **resilience as the default** (~⅔). Avoid pathologizing the majority; reserve "PTSD" for the chronic-trajectory minority; route distress to **PE/CPT/trauma-focused CBT**, the guideline-strongest treatments.
- **Using ACEs responsibly.** Powerful for **population** risk and advocacy; **not** an individual crystal ball — high ACE ≠ destiny.
- **Burnout conversations.** Use Maslach's three dimensions and the JD-R demands>resources framing; note ICD-11 treats it as occupational, and the depression boundary is blurry — don't over-medicalize or over-dismiss.
- **Distinguishing mechanism from coping.** When the user wants *interventions* (recovery, prevention, regulation skills), hand off to `applied-psychology` — keep this spoke to "what's happening and how strong is the evidence."

---

## SKIP routing (sibling spokes & hubs)

| If the user actually wants… | Route to |
|---|---|
| How to **prevent burnout / recover / cope / build resilience** (operator tactics) | `applied-psychology` → `references/performance-and-resilience-psychology` |
| **Emotion-regulation technique**, reappraisal, de-escalation, affect labeling | `applied-psychology` → `references/emotion-and-affect-psychology` |
| **ADHD / autism** neurodevelopmental stress, masking, sensory load | `psychology-neurodevelopmental` |
| **Wellbeing / flourishing / PERMA / post-traumatic growth** as positive-psych | `psychology-positive` |
| **Personality disorders / attachment / dissociation** as clinical-personality depth | `psychology-clinical-personality` |
| **Behavior-change / habit / nudge** to act on stress findings | `applied-psychology` → `references/behavior-change-psychology` |
| **Decision-making under stress**, bias, dual-process | `applied-psychology` → `references/behavioral-decision-making` |
| Clinical **diagnosis or personal medical advice** | a licensed clinician (out of scope here) |

---

## Key theorists / sources (anchor names)
- **Bruce McEwen**: allostatic load; stress mediators; brain remodeling (McEwen & Stellar 1993; *Annual Review* / PNAS / NEJM syntheses).
- **Peter Sterling & Joseph Eyer**: coined **allostasis** (1988).
- **Robert Sapolsky**: glucocorticoid-cascade hypothesis; *Why Zebras Don't Get Ulcers*.
- **Rachel Yehuda**: low-cortisol / enhanced-GR-sensitivity model of PTSD; (contested) intergenerational FKBP5 epigenetics.
- **Mark Gilbertson** (with Pitman) — monozygotic-twin hippocampal-volume vulnerability study (2002).
- **Vincent Felitti & Robert Anda**: ACE Study (1998).
- **Christina Maslach**: burnout three-dimensional model / MBI.
- **George Bonanno**: resilience-as-modal-response trajectory model.
- **Stephen Porges**: polyvagal theory (theory's author).
- **Paul Grossman & Edwin W. Taylor**: the principal physiological **refutation** of polyvagal theory (2007; Grossman 2023).
