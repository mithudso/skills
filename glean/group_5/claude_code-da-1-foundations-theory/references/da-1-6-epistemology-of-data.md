<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-6-epistemology-of-data` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-6-epistemology-of-data
description: >-
  The epistemology of data for data analysis — how data come to count as
  evidence and how warranted claims (knowledge) are produced from them. Covers
  the relational definition of data, the "raw data" myth and theory-ladenness,
  the data→information→knowledge chain and the DIKW critique, the "end of
  theory" debate, justification and validation of data-derived claims, and the
  fitness-for-purpose view of data quality.
  TRIGGER: questions about what data "really are," whether data are objective or
  raw, whether more data removes the need for theory/models, how a dataset earns
  the status of evidence, what it means to "know" something from data, why two
  analysts read the same dataset differently, the DIKW pyramid, data provenance
  as a basis for trust, or framing the philosophical assumptions behind an
  analysis.
  SKIP: hands-on data-cleaning/wrangling mechanics (no skill claim); concrete
  reproducibility/replicability practice and the replication crisis
  (da-1-6-3-reproducibility-replicability); formal logic of inductive vs.
  deductive inference (da-1-6-2-inductive-vs-deductive-reasoning); whether a
  measured association is causal (da-1-6-1-correlation-vs-causation); reliability
  & validity as measurement properties (da-1-2-3-reliability-validity);
  operationalization of constructs (da-1-2-4-operationalization-constructs);
  frequentist vs. Bayesian inference machinery
  (da-1-4-3-frequentist-vs-bayesian-paradigms).
---

# Epistemology of Data

Epistemology is the theory of knowledge: what knowledge is, how it is justified,
and where it comes from. The **epistemology of data** asks the same questions
about data specifically — when does a recorded observation count as *evidence*,
and how does an analysis turn that evidence into a *warranted claim* a person can
act on? This skill is scoped to data analysis as a discipline: it supplies the
conceptual frame an analyst stands on before any technique is applied. It is the
"why should anyone believe this?" layer underneath the lifecycle, the
statistics, and the tooling.

Treat this as the philosophical posture you adopt toward a dataset, not a
procedure you run.

## Why an analyst should care

Every analysis silently commits to a theory of knowledge. Pretending data are
neutral, objective, and self-interpreting leads to over-confident conclusions
and to the belief that scale alone settles questions. Naming the epistemic
assumptions makes an analysis auditable: you can state what the data are evidence
*of*, what would change the conclusion, and where the warrant is thin. The
recurring practical failures — drawing causal stories from convenience samples,
trusting a metric because it is precise, assuming a bigger dataset is a truer
one — are epistemology errors before they are statistical ones.

## Core positions

### 1. Data are relational, not given

The word "data" descends from Latin *datum*, "that which is given," which smuggles
in the idea that data sit in the world waiting to be collected. The dominant
contemporary view rejects this. Sabina Leonelli argues for a **relational**
account: an object counts as data only by virtue of its *function within a
process of inquiry* — something is treated as data when it is used as evidence
for a claim. The same artifact (a photograph, a log line, a survey response) can
be data in one investigation and not in another. There is no intrinsic,
context-free property that makes a thing "data." ([Leonelli, *What Counts as
Scientific Data? A Relational Framework*](https://philpapers.org/rec/LEOWCA);
[LSE Impact Blog interview](https://blogs.lse.ac.uk/impactofsocialsciences/2015/01/19/philosophy-of-data-science-series-sabina-leonelli/))

Consequence for analysis: "is this good data?" is not answerable in the
abstract. It is always "good data *for what purpose, in what inquiry*."

### 2. There is no "raw" data — data are theory-laden

A second pillar follows from the first. Data are never interpretation-free.
Before a single value is analyzed it has already been **selected** (these
fields, not those), **framed** by background assumptions, hypotheses, and the
instrument or schema that recorded it, and **shaped** by the categories the
collectors used. "Raw data" is therefore widely treated as an oxymoron — data
are loosely interpreted at least because they were chosen over other data and
sit inside a paradigm. ([*The epistemological foundations of data science: a
critical analysis*](https://philarchive.org/archive/DESTEF); [Synthese version,
Desai et al.](https://link.springer.com/article/10.1007/s11229-022-03933-2))

This is the data-analysis form of the long-standing **theory-ladenness of
observation** thesis: what you can record depends on the theory you bring.
Wolfgang Pietsch's work on data-intensive science shows the thesis persists even
when models are minimal — choices of features, encodings, and inclusion criteria
are themselves theoretical commitments. ([Pietsch, *Aspects of
theory-ladenness in data-intensive science*](https://core.ac.uk/download/pdf/33752483.pdf))

Practical reading: capturing provenance (who collected it, why, with what
instrument, under what definitions) is not metadata hygiene — it is what lets a
reader reconstruct the interpretation already baked into the numbers.

### 3. The data → information → knowledge chain, and the DIKW critique

A familiar mental model arranges **Data, Information, Knowledge, Wisdom** as a
pyramid (DIKW): data are bare symbols, information is data given context,
knowledge is justified/actionable information, wisdom is judgment about its use.
It is a useful first vocabulary for distinguishing "we recorded a number" from
"we know something." ([DIKW pyramid — Wikipedia](https://en.wikipedia.org/wiki/DIKW_pyramid))

But the model is **philosophically contested**, and an analyst should hold it
loosely. Critics (notably Frické) argue the hierarchy rests on dated
**operationalism and inductivism**, commits a logical error in treating
knowledge as merely accumulated information, and presupposes the very thing it
claims to build — to record data into a structured database you already need the
knowledge that fixes the categories, so "data" is not a theory-neutral
foundation. The clean upward escalator (more data → information → knowledge) does
not hold: whether interrogated data yields useful information depends entirely on
*how* it is interrogated and on the accuracy of the underlying data; information
may or may not become knowledge. ([Frické, *The knowledge pyramid: a critique of
the DIKW hierarchy*](https://experts.arizona.edu/en/publications/the-knowledge-pyramid-a-critique-of-the-dikw-hierarchy/);
[*Revising the DIKW Pyramid*, Law, Technology and Humans](https://lthj.qut.edu.au/article/view/1470))

Use DIKW to *name* the rungs (data vs. information vs. knowledge) but not to
assert that climbing happens automatically or that data is the bedrock.

### 4. The "end of theory" debate

Chris Anderson's 2008 "The End of Theory" claimed that at petabyte scale
correlation supersedes causation and models become obsolete — let the data
speak. This is the strongest empiricist position in data epistemology and it is
the one the literature most consistently rejects. The critical review of data
science argues that **theory cannot be removed even in a data-driven paradigm**:
data-driven results become credible only when *entrenched* in well-curated data,
established theory, and empirical validation; pragmatic predictive success alone
does not epistemically justify a claim (Symons & Alvarado). ([*The
epistemological foundations of data science: a critical
review*](https://link.springer.com/article/10.1007/s11229-022-03933-2);
[philarchive copy](https://philarchive.org/archive/DESTEF))

Stance to carry into an analysis: more data does not retire the need to say what
the data are evidence *for* and why the inference holds. Scale changes power, not
the logic of justification.

### 5. Justification and validation — how a claim earns belief

Epistemology's classic concern is *justified true belief*. For data, the
justification question becomes: what makes a data-derived claim believable? The
literature points to a layered answer rather than a single test —
- **Provenance and curation**: the claim is only as trustworthy as the data
  journey behind it (collection, transformation, integration), because each
  processing step changes what is effectively being treated as data. ([Leonelli,
  *Data-Centric Biology*](https://press.uchicago.edu/ucp/books/book/chicago/D/bo24957334.html))
- **Empirical validation and entrenchment** in existing theory, not predictive
  accuracy in isolation. ([Synthese review](https://link.springer.com/article/10.1007/s11229-022-03933-2))
- **Fitness for the intended use** (next section).

Note the boundary: *how* validation is then operationalized — replication,
hold-out tests, the mechanics of inductive vs. deductive inference, the
frequentist/Bayesian calculus — belongs to the sibling skills listed in SKIP.
This skill establishes *that* justification is layered and contextual and *why*.

### 6. Data quality is fitness-for-purpose, not an intrinsic property

Because data are relational, quality is judged against intended use, not in the
abstract. Data are "high quality" when **fit for their intended uses in
operations, decision making, and planning**, and specifically when they
correctly represent the **real-world construct** they refer to. Assessing this
means interrogating whether a source can support a *particular* analytic
requirement, given its provenance — which also bounds the generalizability of any
result. ([Data quality — Wikipedia](https://en.wikipedia.org/wiki/Data_quality);
[Razzaghi et al., fitness-based data-quality assessment, *Learning Health
Systems*](https://pmc.ncbi.nlm.nih.gov/articles/PMC8753309/))

The diversity of sources and the lack of unified standards make quality one of
the central epistemic challenges of data-intensive work. A precise, clean,
well-formatted dataset can still be epistemically worthless for a question it was
never fit to answer.

## Pitfalls (each is an epistemology error)

- **Treating data as raw/objective.** Forgetting the selection and framing
  already inside a dataset, then reporting findings as if the data "just showed"
  them.
- **Mistaking volume for truth.** Assuming a bigger dataset is automatically a
  more faithful one — the "end of theory" trap.
- **Mistaking precision for validity.** A number with many decimals or a tight
  confidence interval can be measuring the wrong construct.
- **Climbing DIKW on autopilot.** Assuming context-plus-data yields knowledge
  without asking how the data were interrogated.
- **Decontextualized quality claims.** Calling data "good" without naming the
  purpose; quality is relational.
- **Provenance amnesia.** Reusing a dataset for a new question without checking
  whether the conditions that produced it still license the new inference.

## A short worked illustration

A team has a clickstream dataset and is asked "do users prefer the new layout?"

- *Relational lens*: these logs are data *for engagement questions*; whether they
  are data for "preference" depends on whether clicks evidence preference at all.
- *Theory-ladenness*: a "click" was defined by the instrumentation team; bot
  filtering, session windows, and which events were even logged are theoretical
  choices that frame the answer before analysis starts.
- *DIKW caution*: aggregating clicks into "engagement rate" is the data→information
  step; calling a higher rate "preference" is a knowledge claim that needs an
  argued bridge, not an automatic one.
- *End-of-theory check*: a billion more clicks sharpen the estimate but do not
  establish that clicks mean preference.
- *Justification*: the claim earns belief through documented provenance, a
  validation step (the mechanics of which live in the sibling skills), and a
  statement of what would falsify it.
- *Fitness for purpose*: if the logs were collected to debug latency, they may be
  unfit to evidence preference no matter how clean they are.

The deliverable is not just a number but a stated chain from observation to
warranted claim, with its weak links named.

## Sources

1. Sabina Leonelli, *What Counts as Scientific Data? A Relational Framework* —
   https://philpapers.org/rec/LEOWCA
2. Sabina Leonelli, *Data-Centric Biology: A Philosophical Study* (Univ. of
   Chicago Press) —
   https://press.uchicago.edu/ucp/books/book/chicago/D/bo24957334.html ; LSE
   Impact Blog interview —
   https://blogs.lse.ac.uk/impactofsocialsciences/2015/01/19/philosophy-of-data-science-series-sabina-leonelli/
3. *The epistemological foundations of data science: a critical review/analysis*
   (Synthese) —
   https://link.springer.com/article/10.1007/s11229-022-03933-2 ; open copy —
   https://philarchive.org/archive/DESTEF
4. Wolfgang Pietsch, *Aspects of theory-ladenness in data-intensive science* —
   https://core.ac.uk/download/pdf/33752483.pdf
5. Martin Frické, *The knowledge pyramid: a critique of the DIKW hierarchy* —
   https://experts.arizona.edu/en/publications/the-knowledge-pyramid-a-critique-of-the-dikw-hierarchy/ ;
   *Revising the DIKW Pyramid*, Law, Technology and Humans —
   https://lthj.qut.edu.au/article/view/1470 ; DIKW pyramid (Wikipedia) —
   https://en.wikipedia.org/wiki/DIKW_pyramid
6. Data quality / fitness-for-purpose: Data quality (Wikipedia) —
   https://en.wikipedia.org/wiki/Data_quality ; Razzaghi et al., fitness-based
   data-quality assessment, *Learning Health Systems* —
   https://pmc.ncbi.nlm.nih.gov/articles/PMC8753309/
