<!-- hub-reference-banner -->
> **Reference file — part of the `research-methodology` hub.** Formerly the standalone `research-data-annotation-coding` skill.
> Sibling topics in this family are now reference files under the hubs (`research-methodology`, `deep-research`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: research-data-annotation-coding
title: "Data Annotation, Labeling & Qualitative Coding"
description: >-
  Turning raw data into reliable coded/labeled data — unifying ML/data annotation and qualitative research coding. TRIGGER: annotation guidelines or a label schema; choosing/interpreting an inter-annotator/inter-coder agreement metric (Cohen's/Fleiss' kappa, Krippendorff's alpha); active learning or weak supervision (Snorkel); tooling (Label Studio, Prodigy, Argilla); validating LLM-as-annotator labels vs human gold; thematic analysis (Braun & Clarke) or grounded theory (Charmaz); building a coding codebook; coding interviews/cases/notes; low agreement or kappa misuse. SKIP: stats on coded data -> da-1-foundations-theory; modeling/EDA/ML selection -> da-*; training/serving/picking a model -> ai-llm-model-layer; LLM-as-JUDGE calibration or coding agent traces for evals (not annotator) -> eval-driven-development; kappa for PRISMA screening -> systematic-review-research-integrity; dataset codebook/dictionary -> research-data-management-fair. Produces labels/codes, not analyses.
category: custom
version: 1.1.0
updated: 2026-06-15
whenToUse:
  - "writing annotation guidelines or a label schema (single vs multi-label, spans) for an ML/eval dataset"
  - "choosing or interpreting an agreement metric (Cohen's/Fleiss' kappa, Krippendorff's alpha)"
  - "active learning or weak supervision (Snorkel) to cut labeling cost; annotation tooling"
  - "validating LLM-as-annotator / LLM-assisted labels against human ground truth"
  - "qualitative coding: thematic analysis, grounded theory, codebooks; coding interviews/cases/notes"
  - "diagnosing low agreement, single-annotator gold, annotator quality/spam, or kappa misuse"
keywords:
  - data annotation
  - data labeling
  - inter-annotator agreement
  - Krippendorff alpha
  - Cohen kappa
  - active learning
  - weak supervision
  - LLM annotation
  - thematic analysis
  - grounded theory
  - qualitative coding
  - codebook
tags:
  - research
  - annotation
  - data
  - qualitative
---

# Data Annotation, Labeling & Qualitative Coding

Turning raw data into **reliable coded/labeled data** is one discipline, not two. ML
annotation and qualitative research coding solve the same problem (apply a stable
scheme to raw observations, consistently enough that someone else gets the same result),
and they share the same machinery: a scheme, a guideline/codebook, pilot iteration,
multiple coders, disagreement analysis, and reliability checks. This skill is the **craft
of producing trustworthy labels and codes**. It is *not* the statistical analysis you run
on coded data, the ML model you train on labels, or the eval you run on a model; those
route away (see SKIP in the description).

Use it whether you are labeling an ML/eval dataset (sentiment, NER — named-entity recognition, toxicity, preference
pairs) or coding qualitative data (interview transcripts, support-case threads, meeting
notes, open-ended survey responses).

## Core concepts: the two halves are one workflow

| Stage | ML / data annotation | Qualitative coding |
| --- | --- | --- |
| Scheme | Label schema (classes, attributes, relations) | Codebook / coding frame, OR emergent codes |
| Instructions | Annotation guidelines | Codebook definitions + memos |
| Orientation | Schema usually fixed up front (deductive) | Inductive (data-driven) vs deductive (theory-driven) |
| People | Annotators / crowd / SME | Coders / researchers |
| Reliability | Inter-annotator agreement (IAA) | Inter-coder agreement (ICA) / reliability — *or deliberately none* |
| Resolve | Adjudication → gold labels | Consensus, OR reflexive interpretation |
| Cost control | Active learning, weak supervision, LLM assist | Theoretical sampling, focused coding, CAQDAS, LLM assist |

The single most important framing decision is **deductive vs inductive**: do you bring a
predefined scheme to the data, or build the scheme from the data? Most ML annotation is
deductive (decide the classes, then label). Qualitative coding can be either: grounded
theory is radically inductive, reflexive TA is usually inductive, codebook/framework/template
analysis is closer to deductive. This decision governs everything downstream, **including
whether agreement metrics even apply** (see Qualitative coding, below).

## Annotation guidelines & label-schema design

Guidelines are the contract that makes independent annotators converge. They are a
*deliverable*, not a preamble, and the first version is a **hypothesis to be tested**.

A usable guideline contains, at minimum:
- **A label schema** — every class, what each means, and the relationships/attributes.
  State positive *and* negative definitions.
- **Edge-case rules** — what to do when an item does not fit cleanly.
- **10–20 worked examples and counterexamples** — annotators learn from examples far
  faster than from abstract prose. Examples are the spec.
- **The annotation unit** — token, span, sentence, document, turn, whole case.
- **The label structure** — one label per item or **multi-label**; flat or **hierarchical**;
  and for span tasks, whether spans may **overlap or nest**. This choice changes how you
  measure agreement: multi-label and span tasks need a unit-aware set-overlap or span-level
  α, not item-level κ (see below).
- **A conflict-resolution / "when unsure" procedure**.
- **A version number**, linked to the dataset version it produced.

**The non-negotiable practice is the pilot.** Before scaling, have 2–3 annotators
(ideally *not* the guideline authors) label ~50 items, compute agreement with the
task-appropriate chance-corrected coefficient (see the selection guide below), review every
disagreement together, revise the guideline, and repeat. The gate to scale is not a fixed
number but a trend: iterate until a fresh batch stops surfacing new guideline gaps and
agreement has stabilized (for subjective tasks that plateau may sit at moderate, not near
1.0; the bands below are soft). Then scale a pilot of 200–500
items to expose tooling bottlenecks, realistic throughput, and cost-per-label before
committing thousands. In production, keep **~10% overlap** between annotators for ongoing
IAA monitoring, and treat guidelines as a living, versioned document. A trick from Snorkel:
encode candidate rules as labeling functions and run them over hand-labeled data to surface
underspecified cases mechanically.

## Inter-annotator / inter-coder agreement: metrics & selection

Raw percent agreement is misleading because annotators agree by chance: with 90% of
items in one class and two labels, two random annotators "agree" ~82% of the time.
**Chance-corrected** coefficients fix this:

    coefficient = (A_o − A_e) / (1 − A_e)

where `A_o` is observed agreement and `A_e` is expected-by-chance agreement; the
coefficients differ in how they model chance.

**Selection guide:**

| Metric | Data type | Annotators | Notes / when to use |
| --- | --- | --- | --- |
| **Percent / raw agreement** | Any | Any | Report it, but never alone — no chance correction. |
| **Cohen's κ** | Nominal | Exactly 2 | Models per-annotator bias (separate marginals). Unstable when class frequencies or rater biases differ. |
| **Scott's π** | Nominal | 2 | Like Cohen's but assumes a single shared distribution (no per-coder bias). |
| **Fleiss' κ** | Nominal | ≥3 | Generalizes Scott's π via pairwise agreement; classic for crowdsourcing. Assumes equal #ratings per item. |
| **Krippendorff's α** | Nominal / ordinal / interval / ratio | ≥2 | Most general: handles missing data, varying #raters per item, and any measurement scale via a configurable distance metric. Preferred for ordinal/Likert, incomplete designs, and modern NLP. |

Key equivalences: for **nominal data with no missing values**, Krippendorff's α equals
Fleiss' κ; for the **2-annotator nominal complete case**, α equals **Scott's π** (α shares the
single-pooled-distribution chance model of π/Fleiss, not Cohen's separate-marginals model — so
it collapses to Cohen's κ only when the two coders' marginals coincide). Differences
appear precisely when you have missing data, ordinal/continuous scales, or unequal coverage.
For ordinal annotations (emotion, opinion, severity), use **weighted** coefficients (weighted
κ for 2 coders, or Krippendorff's α with an ordinal distance). An off-by-one disagreement
should cost less than an off-by-three.

**Interpretation bands** (rule of thumb, *not* law): <0.20 poor, 0.20–0.40 fair, 0.40–0.60
moderate, 0.60–0.80 good/substantial, 0.80–1.00 very good. Treat these as soft; the
"acceptable" threshold depends on task subjectivity and stakes.

**Report like a professional:** the metric name and scale, the confusion matrix, the number
of items/annotators, confidence intervals (bootstrap), and, most usefully, a
**disagreement analysis**. Low agreement is a *diagnosis*, not just a number: it means
insufficient training, an under-specified guideline, or a genuinely ill-posed task. Choose
the coefficient from the *task* (scale, #annotators, missingness), never by trying several
and reporting the highest.

**Adjudication** turns multiple annotations into a single gold label (majority vote, or an
expert resolving conflicts). Adjudication can cost as much as primary annotation, and both
annotators can be wrong by chance, so gold is "best available," not infallible truth.

**Annotator quality control** is distinct from agreement, and essential the moment you use a
crowd. Plain majority vote treats every annotator as equally reliable; instead, qualify
annotators with a screening task, seed **gold / honeypot check items** throughout the work to
catch spam and drift, and remove or down-weight low-accuracy annotators. When you have
repeated labels per item, model-based aggregation that estimates each annotator's reliability
(**Dawid–Skene**, **MACE**) recovers truer labels than majority vote, especially with noisy
or adversarial crowds.

## Active learning & weak supervision: cutting labeling cost

When unlabeled data is abundant but labeling is expensive, two complementary techniques
reduce human effort:

- **Active learning (AL)** — train on a small seed set, then let the model *choose* which
  unlabeled items a human should label next (the human is the "oracle"), retrain, repeat.
  Pool-based query strategies: **uncertainty sampling** (least-confident, smallest margin,
  or max entropy), **representativeness/diversity** (cover the input space), and **hybrid**.
  AL reaches a target accuracy with far fewer labels; it pairs naturally with model-in-the-loop
  tools (Prodigy is built around it).
- **Weak supervision / data programming (Snorkel)** — instead of hand-labeling, write
  **labeling functions (LFs)**: noisy, programmatic rules/heuristics/pattern/KB-lookups that
  vote or abstain on each item. Snorkel learns a **generative label model** of LF accuracies
  and correlations **with no ground truth** (from their agreements/disagreements), emits
  **probabilistic labels**, then trains a **discriminative model** on those labels that
  *generalizes beyond* the LFs' coverage. Ideal when domain knowledge is expressible as rules.

AL and weak supervision compose: weak supervision bootstraps coverage cheaply; active
learning spends scarce human labels where the model is most uncertain.

## LLM-assisted labeling & its validation

LLMs are now first-line *assistants* for both annotation and qualitative coding — fast and
cheap, and on simpler or more objective tasks with a clear codebook, LLM–human agreement can
be high; on subjective or minority-class tasks it degrades sharply, so measure it per task.
**They are not ground truth.** Treat an LLM annotator like a new, untrusted crowdworker who must earn
trust on *your* task:

- **Always validate against human-labeled gold, per task.** Without gold you cannot detect
  that the model is systematically wrong. LLM errors tend to be **correlated** across samples,
  producing biased labels that look internally consistent. Report LLM-vs-human agreement with
  the *same* chance-corrected metric (κ / α) and a confusion matrix; for binary verdicts,
  report κ or MCC (Matthews correlation coefficient) alongside accuracy (accuracy alone is
  unreliable on imbalanced criteria).
- **Do not train a small model purely on LLM labels and assume parity.** Classifiers trained
  on LLM-generated labels show performance drops, unstable predictions, and early plateaus —
  worse on complex tasks and minority classes. Filtering/ensembles only partly help. A model
  fine-tuned on real human labels still tends to beat one trained on synthetic labels.
- **"Human-in-the-loop" is not automatic safety.** When humans merely *review* LLM suggestions
  on subjective tasks, they often over-accept them (biasing the "human-reviewed" gold toward
  the model), and review may not even save time. If you use HITL (human-in-the-loop), blind the reviewer to the
  model's label on a held-out audit slice to measure over-reliance.
- **Use the LLM where it is genuinely strong:** flagging *likely label errors* in existing
  datasets, pre-annotation to speed humans, and as a *second coder/auditor*. **Self-verification**
  and **cross-model verification** reduce variance and surface systematic errors.
- **Watch for LLM-as-judge biases** (position, verbosity, self-preference) when an LLM scores
  outputs; calibrate against human labels first. (For eval-specific depth — golden datasets,
  judge calibration, criteria drift — defer to the eval-driven-development skill.)

## Qualitative coding methods

Qualitative coding assigns interpretive labels to segments of text. Pick the approach to
match your research values and question, and recognize that the schools differ on whether
reliability metrics apply at all.

**Inductive vs deductive; semantic vs latent.** Inductive = codes emerge bottom-up; deductive
= codes come top-down from theory/a prior frame. Semantic (manifest) coding stays at surface
meaning; latent coding interprets underlying ideas/assumptions.

**Thematic analysis (Braun & Clarke).** The most widely used qualitative method; "TA" is an
umbrella, not one procedure. The **six phases** (recursive, not linear):
1. Familiarisation with the data.
2. Generating initial codes (systematically code the whole dataset).
3. Generating/searching for themes (collate codes into candidate patterns).
4. Reviewing themes (check against coded data and the full dataset).
5. Refining, defining and naming themes.
6. Writing the report.

The **tripartite distinction is load-bearing** for how you handle reliability:
- **Coding-reliability TA** (Boyatzis; Guest et al.) — fixed codebook agreed up front, ≥2
  independent coders, **inter-coder agreement (statistically tested) is the quality signal**.
- **Codebook TA** (incl. framework, template analysis) — middle ground: structured codebook
  used to map/chart, but embracing subjectivity; agreement metrics usually *not* the point.
- **Reflexive TA** (Braun & Clarke's own) — coding is organic/evolving, themes are *patterns
  of shared meaning* developed *through* coding, researcher subjectivity is a **resource not a
  threat**, and a single coder is fine. **Computing κ/α here is a category error** — multiple
  coders act as "critical friends," not consensus-machines.

So: *which TA* determines whether inter-coder agreement is appropriate. Don't bolt a kappa
onto a reflexive analysis to look rigorous.

**Grounded theory (Charmaz, constructivist).** A *whole methodology* aimed at building theory
from data. Defining practices: simultaneous data collection and analysis; constructing
codes/categories from data not preconceived hypotheses; the **constant comparative method**;
**theoretical sampling** (collect more data to develop emerging categories, not for
representativeness); **memo-writing** throughout; literature review *after* an independent
analysis. Coding moves **initial** (line-by-line — Charmaz favors **gerunds** to capture
action) → **focused/selective** → **theoretical** coding. **Axial** coding (relate categories
to subcategories) is the cornerstone of the *Strauss & Corbin* tradition; constructivist GT
treats it as optional and Charmaz is wary of its rigidity, so do not present it as a required
Charmaz phase. Distinguish *coding for topics* (nouns) from *coding with gerunds*
(process, action).

## Codebooks

A codebook is the qualitative analogue of an annotation guideline: for each code, a name, a
definition, **inclusion and exclusion criteria**, example quotes, and a "do not confuse
with…" note. Codebooks suit deductive/codebook TA, framework/template analysis, and any
project needing inter-coder reliability or team coding. They are *less* central to reflexive
TA and emerge *late* in grounded theory. Version the codebook and record when codes are added,
merged, or split — the same traceability discipline as a label schema.

## Tooling

**ML/data annotation:**
- **Label Studio** — broadest open-source platform; text/image/audio/video, ML-backend
  pre-annotation, native **multi-annotator IAA tracking**, native RLHF pairwise template,
  Python SDK. Default for mixed-modality teams.
- **Prodigy** (Explosion/spaCy; paid, local) — **active-learning-first**, scriptable,
  keyboard-driven; best-in-class for NER/text-classification *with a starting model*.
- **Doccano** (open-source) — simple, fast text classification, sequence labeling (NER), and
  seq2seq; great for small teams; minimal IAA.
- **Argilla** (open-source; Hugging Face) — built for **LLM feedback / preference / RLHF** data
  and dataset curation; push/pull to HF Hub.
- Also: **brat** (entity/relation), **Potato** (YAML, no-code, agreement metrics, agent/LLM
  eval), **CVAT** (CV/video), **Snorkel Flow** (programmatic labeling), managed crowds.

**CAQDAS (Computer-Assisted Qualitative Data Analysis Software):**
- **NVivo** — institutional standard; large/complex and mixed-methods; computes **kappa for
  inter-rater reliability**; steeper learning curve.
- **MAXQDA** — most usable; strongest **mixed methods**; good audio/video.
- **ATLAS.ti** — built for **grounded theory / theory building**; strong **AI-assisted coding**;
  its ICA tool was built *with Krippendorff* and warns ICA "crosses the qual-quant divide" and
  demands stricter rules than ordinary coding.
- **Dedoose** — web-based, low-cost, mixed-methods.
- Prefer tools supporting the **REFI-QDA (Rotterdam Exchange Format Initiative) / codebook exchange** format so codebooks and
  projects move between packages.

For ICA in CAQDAS: run it *after* the code system is stable and all codes are defined, with
independent coders on a *subset* of the data — not as a one-click afterthought.

## Pitfalls (anti-patterns)

- **Vague guidelines / no pilot.** Shipping abstract rules with no worked examples and scaling
  before piloting. The pilot is the highest-ROI quality step.
- **Single-annotator "gold."** One annotator (or one model) is not ground truth — you have no
  way to estimate reliability. Use overlap + adjudication.
- **Kappa misuse.** Reporting raw % agreement as reliability; picking the coefficient that gives
  the highest number; ignoring the **kappa paradox** (high agreement but low κ under class
  imbalance); using nominal κ on ordinal data; treating an interpretation band as a hard
  pass/fail.
- **Treating LLM labels as ground truth.** Skipping human validation; training a downstream model
  on LLM labels and assuming parity; trusting "human-reviewed" LLM labels without measuring
  over-acceptance.
- **Coding without a codebook (where one is required).** Team or deductive coding with no shared,
  versioned codebook → uninterpretable, irreproducible results.
- **Category error: forcing IRR onto reflexive TA / GT.** Computing κ/α on a reflexive or
  grounded-theory analysis to "prove rigor." Those traditions draw rigor from reflexivity, memos,
  an audit trail, and a coherent account — not agreement coefficients.
- **Untracked guideline/codebook drift.** Changing the scheme mid-project without versioning →
  you can't trace a label back to the rules in force.
- **Wrong unit / unstable scheme.** Computing agreement before the code system is stable, or
  mixing annotation units, makes every reliability number meaningless.

## References

1. James, *Counting on Consensus: Selecting the Right Inter-annotator Agreement Metric for NLP* (2026). https://arxiv.org/html/2603.06865
2. Artstein & Poesio, *Inter-Coder Agreement for Computational Linguistics*, Computational Linguistics 34(4) (2008). https://aclanthology.org/J08-4004.pdf
3. Bamman, *Annotation* lecture, Berkeley INFO 256. https://people.ischool.berkeley.edu/~dbamman/info256_slides/7.annotation.pdf
4. Snorkel AI, *Data annotation guidelines and best practices* (2022). https://snorkel.ai/blog/data-annotation/
5. Ratner et al., *Snorkel: Rapid Training Data Creation with Weak Supervision*, VLDB (2018). https://www.vldb.org/pvldb/vol11/p269-ratner.pdf
6. Ratner, Varma, Hancock, Ré et al., *Weak Supervision: A New Programming Paradigm for ML*, Stanford SAIL (2019). https://ai.stanford.edu/blog/weak-supervision/
7. Tharwat & Schenck, *A Survey on Active Learning*, Mathematics 11(4):820 (2023). https://www.mdpi.com/2227-7390/11/4/820
8. Horych et al., *The Promises and Pitfalls of LLM Annotations in Dataset Labeling* (NAACL Findings 2025). https://aclanthology.org/2025.findings-naacl.75.pdf
9. Pangakis & Wolken, *grounding automated annotation in human judgment* (2024). https://arxiv.org/pdf/2409.09467
10. *Feeding LLM Annotations to BERT Classifiers at Your Own Risk* (2025). https://arxiv.org/html/2504.15432
11. *Just Put a Human in the Loop? Investigating LLM-Assisted Annotation for Subjective Tasks* (ACL Findings 2025). https://aclanthology.org/2025.findings-acl.1323.pdf
12. Braun & Clarke, *Doing reflexive TA* (six phases). https://www.thematicanalysis.net/doing-reflexive-ta/
13. Braun & Clarke, *Using thematic analysis in psychology*, QRP 3(2) (2006). https://educationaldevelopment.uams.edu/wp-content/uploads/sites/57/2025/01/9-Thematic_analysis.pdf
14. Braun & Clarke, *Thematic analysis: Choosing a suitable approach* (coding-reliability vs codebook vs reflexive TA), SRA. https://the-sra.org.uk/SRA/Blog/ThematicanalysisChoosingasuitableapproach.aspx
15. Charmaz, *Constructing Grounded Theory* (2006). https://www.sxf.uevora.pt/wp-content/uploads/2013/03/Charmaz_2006.pdf
16. ATLAS.ti, *Inter-Coder Agreement (ICA)* manual. https://manuals.atlasti.com/Mac/en/manual/ICA/InterCoderAgreemenIntroduction.html
17. Univ. of Michigan (Potato), *Open-Source Annotation Tools Compared* (2026). https://www.potatoannotator.com/docs/guides/annotation-tools-compared
