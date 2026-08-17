<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-6-3-reproducibility-replicability` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-6-3-reproducibility-replicability
description: >-
  Epistemology-of-data concept covering reproducibility and replicability as
  criteria for what counts as trustworthy knowledge from data analysis: the
  precise NASEM/ACM/Turing-Way definitions, the reproduce/replicate/robust/
  generalize distinction, the replication crisis and its causes, and the
  computational practices that make an analysis re-runnable.
  TRIGGER: questions about what "reproducible" vs "replicable" mean; whether a
  result/finding/analysis can be trusted or re-run; the replication crisis,
  p-hacking, researcher degrees of freedom, HARKing as epistemic problems;
  designing an analysis so others can verify it; reproducibility standards,
  badges, or reporting requirements; "same data same code" vs "new data".
  SKIP: hands-on environment pinning, Docker/conda/lockfile mechanics, or
  CI pipeline authoring as an engineering task (defer to coding/devops skills);
  experiment-design and sampling theory as their own statistical topics; data
  versioning as a tooling how-to; correlation-vs-causation (da-1-6-1) and
  inductive-vs-deductive reasoning (da-1-6-2) which are sibling epistemology
  nodes; statistical significance testing mechanics (defer to inference skills).
---

# Reproducibility & Replicability (Epistemology of Data)

This skill sits under **Data Analysis > Foundations & Theory > Epistemology of
data**. It treats reproducibility and replicability as *epistemic criteria* —
the standards by which a claim derived from data earns the status of reliable
knowledge — not as a DevOps how-to. Use it to reason about whether a result
should be believed, what would make it verifiable, and why a field's body of
findings might be unreliable.

## The core distinction (get the terms right)

The single most common error is using "reproduce" and "replicate"
interchangeably. The authoritative framing fixes two axes: **same vs. different
data**, and **same vs. different analysis**.

The U.S. National Academies (NASEM) 2019 consensus report set the reference
definitions:

- **Reproducibility** = "obtaining consistent results using the same input
  data; computational steps, methods, and code; and conditions of analysis."
  This is essentially *computational* reproducibility: re-run the original
  data through the original code and get the original numbers
  ([NASEM 2019, ch.3](https://www.nationalacademies.org/read/25303/chapter/3)).
- **Replicability** = "obtaining consistent results across studies aimed at
  answering the same scientific question, each of which has obtained its own
  data." Replication collects *new* data
  ([NASEM 2019, ch.3](https://www.nationalacademies.org/read/25303/chapter/3)).

The defining contrast: **reproducibility re-uses the original data and code;
replicability gathers new data**
([NASEM 2019, ch.3](https://www.nationalacademies.org/read/25303/chapter/3)).

The Turing Way extends this into a 2x2 matrix that adds two more terms, which
is the most useful mental model for data analysis
([Turing Way, Definitions](https://book.the-turing-way.org/reproducible-research/overview/overview-definitions/)):

| | **Same data** | **Different data** |
|---|---|---|
| **Same analysis** | **Reproducible** | **Replicable** |
| **Different analysis** | **Robust** | **Generalisable** |

- **Reproducible** — same analysis on the same data gives the same answer.
- **Replicable** — same analysis on *different* data gives a qualitatively
  similar answer.
- **Robust** — *different* analysis pipeline (e.g. R vs. Python) on the same
  data gives a similar answer.
- **Generalisable** — a finding that is both replicable and robust; not
  dependent on a particular dataset or a particular pipeline
  ([Turing Way, Definitions](https://book.the-turing-way.org/reproducible-research/overview/overview-definitions/)).

### A terminology landmine

The labels are *not* universal. The ACM historically defined the two terms the
opposite way from NASEM, then in 2020 swapped "Results Reproduced" and "Results
Replicated" to align with the ISO/NISO convention
([ACM, Artifact Review and Badging — Current](https://www.acm.org/publications/policies/artifact-review-and-badging-current);
[ACM, New Changes to Badging Terminology](https://www.acm.org/publications/badging-terms)).
Because the vocabulary is contested across communities, **state which axis you
mean** ("same-data-same-code" vs. "new-data") rather than relying on the word
alone. The confusion has its own literature (Plesser, "Reproducibility vs.
Replicability: A Brief History of a Confused Terminology", cited in the
[Turing Way](https://book.the-turing-way.org/reproducible-research/overview/overview-definitions/)).

ACM's operational definitions are worth knowing because they tie the concept to
*artifacts* — the digital objects (data, code, scripts) behind a result:

- **Repeatability** — the same team reliably repeats its own result.
- **Reproducibility** — an *independent* team gets the same result using *the
  authors' own* artifacts.
- **Replicability** — an independent team gets the result using artifacts they
  developed *completely independently*
  ([ACM, Artifact Review and Badging — Current](https://www.acm.org/publications/policies/artifact-review-and-badging-current)).

## Why this is an epistemology question

Reproducibility is the minimum bar: if you cannot even re-derive a number from
its own data and code, the claim is not yet knowledge — it is an
unverified assertion. Replicability is the higher, scientific bar: a finding
that survives new data is more likely to reflect a real effect rather than an
artifact of one dataset. The two are *complementary criteria of warrant*, which
is why they belong under epistemology of data rather than under tooling.

NASEM frames the sources of *non*-reproducibility along an epistemic line
([NASEM 2019, ch.3](https://www.nationalacademies.org/read/25303/chapter/3)):

- **Potentially helpful / unavoidable sources** — intrinsic variability in the
  system, measurement limits, the genuine difficulty of complex phenomena.
  These are part of normal discovery, not failures.
- **Unhelpful / avoidable sources** — inadequate design, sloppy conduct,
  incomplete reporting, bias, or misconduct. These are inefficiencies that
  impede knowledge.

The practical upshot: a non-reproducible result is not automatically *wrong*,
but it is *unwarranted* until the source of variation is identified.

## The replication crisis (the cautionary core)

A large fraction of published findings fail to replicate. In psychology,
roughly 40% of studies in top journals did not replicate
([An Introduction to Data Analysis, App. E.1](https://michael-franke.github.io/intro-data-analysis/app-94-replication-crisis.html)).
The crisis is the empirical evidence that the replicability criterion is
routinely *not met* across fields, including biomedicine and economics
([Replication crisis — Wikipedia](https://en.wikipedia.org/wiki/Replication_crisis)).

Mechanisms that inflate false-positive findings — the things a data analyst
must guard against:

- **Researcher degrees of freedom / p-hacking** — exploiting flexibility in
  data collection and analysis (dropping outliers, choosing subgroups, switching
  tests) until *p* crosses the threshold
  ([An Introduction to Data Analysis, App. E.1](https://michael-franke.github.io/intro-data-analysis/app-94-replication-crisis.html)).
- **HARKing** — Hypothesizing After Results are Known: presenting a post-hoc
  finding as if it were predicted in advance
  ([Replication crisis — Wikipedia](https://en.wikipedia.org/wiki/Replication_crisis)).
- **Publish-or-perish incentives** — pressure for significant, novel results
  drives questionable research practices and a literature skewed toward false
  positives
  ([An Introduction to Data Analysis, App. E.1](https://michael-franke.github.io/intro-data-analysis/app-94-replication-crisis.html)).

Important nuance, reported as Low-confidence because sources differ on
magnitude: some analyses argue that even large reductions in p-hacking would
yield only modest gains in replicability, implying that broader scientific
culture and incentive structures — not just individual analyst behavior — are
the binding constraint
([Questionable research practices may have little effect on replicability, PMC7561355](https://pmc.ncbi.nlm.nih.gov/articles/PMC7561355/)).

## How to make a data analysis reproducible

Reproducibility is the part most under an analyst's direct control. The
literature converges on **five pillars of computational reproducibility**
([The five pillars of computational reproducibility, PMC10591307](https://www.ncbi.nlm.nih.gov/pmc/articles/PMC10591307/)):

1. **Literate programming** — interleave code, results, and narrative (e.g.
   notebooks, R Markdown) so the analysis and its explanation stay in sync.
2. **Code version control and sharing** — track every change to the analysis
   code and publish it.
3. **Compute environment control** — pin language/library versions so the same
   code produces the same result on another machine.
4. **Persistent data sharing** — make the input data available and citable.
5. **Documentation** — describe inputs, steps, and decisions completely enough
   to re-run.

Two analysis-specific points that recur across sources:

- **Set the random seed.** Any step with algorithmic randomness (resampling,
  cross-validation splits, stochastic models) must initialize its pseudo-random
  generator with a fixed value, or the result is not bit-for-bit reproducible
  ([Harnessing the best software tools for reproducible analysis](https://hashadatascience.com/reproducible_data_science/)).
- **Capture the environment.** Containerization plus continuous integration
  (e.g. "continuous analysis") can automatically re-run the pipeline whenever
  code or data change, catching reproducibility breaks early
  ([Reproducibility of computational workflows is automated using continuous analysis, PMC6103790](https://pmc.ncbi.nlm.nih.gov/articles/PMC6103790/)).

These tooling notes are included only as *what reproducibility requires*. For
the engineering how-to (Docker, conda/Poetry lockfiles, DVC, CI authoring),
defer to coding/devops skills — that is out of scope here.

## Reporting standards that operationalize the criteria

The epistemic criteria become actionable through reporting and badging
requirements:

- **NASEM Recommendation 6-1** — provide a "clear, specific, and complete
  description" of the methodology: all procedures, measurements, variables,
  data-analysis decisions, statistical methods, and uncertainty characterization
  ([NASEM 2019, ch.3](https://www.nationalacademies.org/read/25303/chapter/3)).
- **NASEM Recommendation 4-1** — convey "clear, specific, and complete
  information about any computational methods and data products": input data,
  detailed method descriptions, and the computational environment
  ([NASEM 2019, ch.3](https://www.nationalacademies.org/read/25303/chapter/3)).
- **ACM artifact badging** — independent reviewers audit the artifacts behind a
  paper and award badges (Artifacts Available / Evaluated / Results Reproduced /
  Results Replicated), making verifiability a visible, rewarded property
  ([ACM, Artifact Review and Badging — Current](https://www.acm.org/publications/policies/artifact-review-and-badging-current)).

## Decision checklist for an analyst

1. **Which claim am I making?** Reproducible (same data/code), replicable (holds
   on new data), robust (holds under a different pipeline), or generalisable
   (both)? Name the axis explicitly.
2. **Is it even reproducible yet?** Can someone re-run my data + code and get my
   numbers? If not, fix that before claiming anything.
3. **Did I fix the seed and pin the environment?** Otherwise "reproducible" is
   not literally true.
4. **Did I pre-specify the analysis, or am I HARKing / p-hacking?** Distinguish
   confirmatory from exploratory results and say which this is.
5. **Did I report enough** (per NASEM 6-1 / 4-1) for an independent party to
   verify and to attempt replication?

## Common pitfalls

- Treating "reproducible" and "replicable" as synonyms — they answer different
  epistemic questions.
- Assuming non-reproducible = wrong; it may stem from unavoidable variability
  ([NASEM 2019, ch.3](https://www.nationalacademies.org/read/25303/chapter/3)).
- Claiming reproducibility without a fixed seed or pinned environment.
- Presenting an exploratory, p-hacked finding as a confirmatory one.
- Trusting a single non-replicated study as established knowledge.

## Sources

1. [NASEM, *Reproducibility and Replicability in Science* (2019), ch.3](https://www.nationalacademies.org/read/25303/chapter/3) — authoritative definitions, sources of non-reproducibility, Recommendations 4-1 and 6-1.
2. [The Turing Way — Definitions](https://book.the-turing-way.org/reproducible-research/overview/overview-definitions/) — reproducible/replicable/robust/generalisable 2x2 matrix.
3. [ACM — Artifact Review and Badging (Current)](https://www.acm.org/publications/policies/artifact-review-and-badging-current) and [Badging Terminology change](https://www.acm.org/publications/badging-terms) — artifact-based definitions and the terminology swap.
4. [The five pillars of computational reproducibility (PMC10591307)](https://www.ncbi.nlm.nih.gov/pmc/articles/PMC10591307/) — five practical pillars.
5. [An Introduction to Data Analysis, App. E.1 — Psychology's replication crisis](https://michael-franke.github.io/intro-data-analysis/app-94-replication-crisis.html) — replication failure rates, p-hacking, researcher degrees of freedom.
6. [Replication crisis — Wikipedia](https://en.wikipedia.org/wiki/Replication_crisis) — cross-field scope, HARKing.
7. [Continuous analysis (PMC6103790)](https://pmc.ncbi.nlm.nih.gov/articles/PMC6103790/) — automating reproducibility via containers + CI.
8. [Questionable research practices may have little effect on replicability (PMC7561355)](https://pmc.ncbi.nlm.nih.gov/articles/PMC7561355/) — contrarian evidence on p-hacking's impact (Low-confidence nuance).
9. [Harnessing the best software tools for reproducible analysis](https://hashadatascience.com/reproducible_data_science/) — random seed / determinism.
