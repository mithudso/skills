<!-- hub-reference-banner -->
> **Reference file — part of the `da-applied-and-communication` hub.** Formerly the standalone `da-35-synthetic-data-generation` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-35-synthetic-data-generation
description: >-
  Synthetic data generation as a discipline — generating artificial records that
  preserve the statistical structure of real data for privacy, augmentation,
  testing/dev, and class rebalancing. Covers tabular synthesis (Gaussian copula,
  CTGAN, TVAE, CopulaGAN, CART/sequential), the SDV + SDMetrics ecosystem, deep
  generative methods (GANs, VAEs, diffusion/TabDDPM), class-imbalance resampling
  (SMOTE family, ADASYN) and when generative beats resampling, differentially
  private synthesis (DP-GAN, PATE-GAN, marginal-based PrivBayes/MST, SmartNoise),
  the fidelity-vs-utility-vs-privacy trade-off and its evaluation (TSTR, distance
  to closest record, membership-inference risk), a text/image synthesis overview,
  and the regulatory picture (GDPR, ICO, NIST). TRIGGER: "synthetic data",
  "generate fake/artificial data", "CTGAN", "TVAE", "SDV", "Synthetic Data Vault",
  "SDMetrics", "Gaussian copula synthesis", "TabDDPM", "diffusion for tabular",
  "differentially private synthetic data", "DP-GAN", "PATE-GAN", "PrivBayes",
  "MST synthesizer", "SmartNoise synthesizers", "TSTR", "train on synthetic test
  on real", "membership inference on synthetic data", "is synthetic data GDPR
  compliant / anonymous", "synthetic data for testing/dev", "rebalance classes
  with generated data". SKIP: plain SMOTE/ADASYN resampling with no synthesis-as-
  discipline question (→ da-4-data-cleaning-preparation); differential privacy /
  k-anonymity theory in the abstract (→ da-11-ethics-and-privacy); feature
  engineering / encoding (→ da-17-feature-engineering-and-feature-stores);
  training a general GAN/VAE/diffusion model for images as a modeling task with no
  data-generation goal (→ da-7-machine-learning).
---

# Synthetic Data Generation

Synthetic data is artificial data produced by a model fit to real data, designed
to reproduce the real data's statistical properties (marginals, correlations,
joint structure) without being a copy of any real record. The goal is to make
synthetic records *useful* for a downstream task while breaking the one-to-one
link to real individuals.

This skill treats synthetic data generation as a discipline: when to use it, which
generator to reach for, how to make it private, and — most importantly — how to
prove it is good enough on the three axes that always compete: **fidelity, utility,
privacy**.

## When to use this skill

Use it when the task is to *create* data rather than analyze existing data:
share data you cannot share in raw form, augment a too-small or imbalanced
training set, populate a test/dev/CI environment without production PII, or
release a public benchmark from a sensitive source.

Do **not** use it for plain class-rebalancing with off-the-shelf SMOTE (that is a
data-prep step — see `da-4`), abstract differential-privacy/k-anonymity theory
(`da-11`), or feature encoding (`da-17`). Those skills are adjacent; this one owns
the *generation pipeline and its evaluation*.

---

## Core concepts

### 1. Why synthetic data (the four motivations)

- **Privacy / data sharing.** Replace a protected dataset with a synthetic
  surrogate so analysts, vendors, or the public can work without touching PII.
  Gartner projected synthetic data would let organizations avoid 70% of
  privacy-violation sanctions by 2025, and estimated ~60% of data used in AI/
  analytics projects in 2024 would be synthetically generated
  ([Gartner, 2024](https://www.gartner.com/en/newsroom/press-releases/2024-06-27-safeguarding-privacy-with-synthetic-data);
  [MIT Sloan, 2023](https://mitsloan.mit.edu/ideas-made-to-matter/what-synthetic-data-and-how-can-it-help-you-competitively)).
- **Augmentation.** Fill gaps, enlarge small datasets, balance rare classes, or
  simulate scenarios that are scarce in reality (fraud, equipment failure).
- **Testing / dev / CI.** Populate non-prod environments at scale with
  realistic-but-fake records; Gartner's 2024 Data Masking Market Guide tells
  teams to prefer test-data-management tools that emit synthetic records
  ([Gartner via K2view, 2024](https://www.k2view.com/gartner-report-synthetic-data-generation)).
- **Class rebalancing.** Generate minority-class examples so a classifier sees a
  less skewed distribution (overlaps with the SMOTE/generative discussion below).

> Inference, not fact: the "synthetic surpasses real by 2030" line is a vendor/
> analyst projection, not a measured outcome. Treat market-size and adoption
> numbers as directional.

### 2. Tabular synthesis methods

Tabular data is the hard case: mixed types, non-Gaussian/multimodal continuous
columns, and highly imbalanced categoricals. Main families:

- **Gaussian copula.** Transform each column to a standard-normal space (via its
  CDF), fit a multivariate Gaussian to capture the correlation structure, sample,
  then invert back to the original margins. Fast, transparent, stable on small
  data; weak on highly non-linear dependencies
  ([SDV docs](https://docs.sdv.dev/sdv/single-table-data/modeling/synthesizers/gaussiancopulasynthesizer)).
- **CTGAN (Conditional Tabular GAN).** Uses **mode-specific normalization** for
  multimodal continuous columns and a **conditional generator + training-by-
  sampling** to handle imbalanced categoricals. Outperformed Bayesian-network and
  prior GAN baselines on ≥87% of test datasets in the original study
  ([Xu et al., NeurIPS 2019](https://proceedings.neurips.cc/paper_files/paper/2019/file/254ed7d2de3b23ab10936522dd547b78-Paper.pdf);
  [CTGAN repo](https://github.com/sdv-dev/CTGAN)).
- **TVAE (Tabular VAE).** Variational autoencoder from the same NeurIPS 2019
  paper; often strong on smaller datasets where diversity matters, and trains
  more stably than a GAN.
- **CopulaGAN.** Hybrid that applies the copula transform before GAN modeling
  (experimental in SDV).
- **CART / sequential synthesis.** Non-deep, classical approach (popularized by
  the R `synthpop` package): synthesize one column at a time, each conditioned on
  the already-synthesized columns, using a CART regression/classification tree.
  Interpretable, fast, the default in much official-statistics work.

### 3. Deep generative methods

- **GANs.** Generator vs. discriminator adversarial game; high fidelity on
  complex distributions but training instability and mode collapse are real
  risks. CTGAN is the tabular workhorse.
- **VAEs.** Encoder/decoder with a latent prior; smoother, more stable training,
  good diversity, sometimes blurrier/less sharp than GANs. TVAE is the tabular
  instance.
- **Diffusion models.** Iterative denoising from a Gaussian prior. **TabDDPM**
  (ICML 2023) applies diffusion to any tabular dataset across numerical and
  categorical features and beat prior SOTA on several benchmarks; concurrent work
  includes STaSy and CoDi, with latent-space and mixed-type variants (TabSyn,
  TabDiff) following
  ([Kotelnikov et al., TabDDPM, ICML 2023](https://proceedings.mlr.press/v202/kotelnikov23a/kotelnikov23a.pdf);
  [survey repo, 2024-25](https://github.com/Diffusion-Model-Leiden/awesome-diffusion-models-for-tabular-data)).

### 4. Class imbalance: resampling vs. generative

- **SMOTE** interpolates new minority points along line segments between a
  minority sample and its k nearest minority neighbors — based on *local* info,
  not the global minority distribution
  ([imbalanced-learn / Chawla et al.](https://machinelearningmastery.com/smote-oversampling-for-imbalanced-classification/)).
- **ADASYN** adapts how many synthetic points to make per minority sample using a
  density criterion, generating more in *harder-to-learn* regions
  ([imbalanced-learn docs](https://machinelearningmastery.com/smote-oversampling-for-imbalanced-classification/)).
- **Generative (GAN/VAE/diffusion)** can capture complex, non-linear joint
  structure that SMOTE/ADASYN interpolation misses; hybrids like SMOTE→GAN refine
  unrealistic SMOTE points
  ([MDPI Mathematics, 2023](https://www.mdpi.com/2227-7390/11/16/3605);
  [SMOTified-GAN, 2021](https://arxiv.org/pdf/2108.03235)).
- **Rule of thumb.** Low-dimensional numeric data, small budget → SMOTE/ADASYN.
  High-cardinality categoricals, strong non-linear interactions, or a privacy
  requirement → generative. Beware: SMOTE on the full dataset before train/test
  split leaks information — resample inside the CV fold only.

### 5. Differentially private synthesis

Plain synthetic data is **not** automatically private; formal guarantees require
differential privacy (DP) built into training.

- **DP-GAN / DP-SGD GANs.** Add calibrated noise + gradient clipping to the
  discriminator's gradients; the post-processing theorem then makes the generator
  (and its output) DP.
- **PATE-GAN.** Replaces the discriminator with a PATE-style ensemble of teacher
  discriminators whose noisy aggregated votes train a student; this makes the
  generator DP by post-processing and often gives better utility than DP-GAN at
  the same epsilon
  ([Jordon, Yoon, van der Schaar, ICLR 2019](https://openreview.net/pdf?id=S1zk9iRqF7)).
- **Marginal-based (PrivBayes, MST).** Privately measure low-order marginals with
  the Laplace/Gaussian mechanism, then reconstruct a distribution (a Bayesian
  network for **PrivBayes**; a graphical model over a maximum-spanning-tree of
  marginals for **MST**). **MST won the 2018 NIST DP Synthetic Data Challenge**;
  both are noted as strong against membership-inference while keeping high utility
  ([Zhang et al., PrivBayes, 2017](https://dl.acm.org/doi/10.1145/3134428);
  [McKenna, Miklau, Sheldon, MST, 2021](https://arxiv.org/pdf/2301.08844);
  [NIST DP Synthetic Data Challenge](https://pages.nist.gov/privacy_collaborative_research_cycle/pages/techniques.html)).
- **SmartNoise (OpenDP).** Microsoft + Harvard OpenDP toolkit; the
  `smartnoise-synth` library exposes DP synthesizers (MST, PrivBayes/PATE-style,
  AIM, etc.) with a uniform `fit()`/`sample()` API
  ([Microsoft, 2021](https://opensource.microsoft.com/blog/2021/02/18/create-privacy-preserving-synthetic-data-for-machine-learning-with-smartnoise/);
  [SmartNoise synth docs](https://docs.smartnoise.org/synth/index.html)).

> Caveat: high epsilon (weak privacy budget) hollows out the guarantee. Recent
> work shows MST/PrivBayes at high epsilon still leak
> ([arXiv 2402.06699, 2024](https://arxiv.org/html/2402.06699v1)). DP is only as
> strong as the epsilon you actually set.

### 6. Tools & frameworks

| Tool | Role |
| --- | --- |
| **SDV (Synthetic Data Vault)** | Python one-stop library; single-table (GaussianCopula, CTGAN, TVAE, CopulaGAN), multi-table (HMA), and sequential (PAR) synthesizers; metadata-driven `fit()`/`sample()` ([SDV docs](https://docs.sdv.dev/sdv); [repo](https://github.com/sdv-dev/SDV)) |
| **SDMetrics** | SDV's evaluation library: a Quality Report (column shapes, column-pair trends), Diagnostic Report (validity/structure), and privacy metrics ([SDV docs](https://docs.sdv.dev/sdv)) |
| **CTGAN (sdv-dev/ctgan)** | Standalone CTGAN + TVAE package ([repo](https://github.com/sdv-dev/CTGAN)) |
| **synthcity** | van der Schaar lab library spanning tabular/time-series/survival + privacy benchmarks ([arXiv 2301.07573](https://arxiv.org/pdf/2301.07573)) |
| **SmartNoise-synth (OpenDP)** | DP synthesizers with formal guarantees ([docs](https://docs.smartnoise.org/synth/index.html)) |
| **imbalanced-learn** | SMOTE/ADASYN and variants for class rebalancing |
| **synthpop (R)** | CART/sequential synthesis, official-statistics standard |

### 7. Evaluation: fidelity vs. utility vs. privacy

The three axes trade off — pushing one usually costs another (especially privacy
vs. fidelity/utility under DP). Measure all three; never report only one.

- **Fidelity** — does synthetic *look like* real?
  - Per-column shape: Kolmogorov–Smirnov (continuous), total-variation / chi-sq
    (categorical) — SDMetrics "column shapes".
  - Pairwise structure: correlation-similarity, contingency similarity — SDMetrics
    "column-pair trends".
  - Validity/structure: ranges, types, key uniqueness (SDMetrics Diagnostic).
- **Utility** — does synthetic *work like* real for the task?
  - **TSTR (Train on Synthetic, Test on Real):** train the model only on
    synthetic, evaluate on a held-out *real* test set; compare to TRTR
    (train-real-test-real). TSTR ≈ TRTR means high utility
    ([AWS ML blog, 2023](https://aws.amazon.com/blogs/machine-learning/how-to-evaluate-the-quality-of-the-synthetic-data-measuring-from-the-perspective-of-fidelity-utility-and-privacy/);
    [TSTR overview](https://www.emergentmind.com/topics/train-on-synthetic-test-on-real-tstr)).
- **Privacy** — can an attacker recover real records?
  - **Distance to Closest Record (DCR):** flag synthetic rows that are near-copies
    of training rows.
  - **Membership inference attack (MIA):** can an adversary tell whether a given
    real record was in the training set? Low attacker advantage = better privacy
    ([health eval framework, Frontiers Digital Health, 2025](https://www.frontiersin.org/journals/digital-health/articles/10.3389/fdgth.2025.1576290/full)).

### 8. Text & image synthesis (overview)

- **Text.** Instruction-tuned LLMs generate synthetic text/labels from prompts
  (frameworks: DataGen, MagPie); useful for NLP augmentation but quality/diversity
  and label noise are real limits, and DP-via-API methods (Aug-PE/Private
  Evolution) exist for private text
  ([LLM text-classification synthesis, 2023](https://ar5iv.labs.arxiv.org/html/2310.07849);
  [DP synthetic text via APIs, 2024](https://arxiv.org/pdf/2403.01749)).
- **Images.** GANs gave limited augmentation benefit; the diffusion era
  (Stable Diffusion, DALL·E, Imagen) sharply improved fidelity and dropped
  per-image cost, making large-scale synthetic image augmentation viable
  ([Synthetic Data in 2024 review](https://www.timlrx.com/blog/synthetic-data-in-2024-progress-opportunities-challenges/);
  [VAE/GAN/diffusion image study, 2025](https://www.mdpi.com/2313-433X/11/8/252)).

### 9. Regulatory context

Synthetic data is **not automatically anonymous** and not automatically out of
GDPR scope. Fully synthetic data that meets the anonymisation bar escapes GDPR;
partially synthetic data usually remains personal data. The ICO's March 2025
anonymisation guidance stresses that effective anonymisation is a high, hard-to-
meet bar and requires a documented re-identification-risk assessment
([GDPR Local](https://gdprlocal.com/synthetic-data-under-gdpr/);
[ICO guidance summary, RPC 2025](https://www.rpclegal.com/snapshots/data-protection/summer-2025/ico-publishes-new-guidance-on-anonymisation-and-pseudonymisation/);
[NIST DP synthetic-data program](https://pages.nist.gov/privacy_collaborative_research_cycle/pages/techniques.html)).
Practical implication: treat a privacy claim as something you must *measure* (DCR +
MIA, ideally under a DP budget), not something the word "synthetic" grants.

---

## Methodology (end-to-end pipeline)

1. **Define the goal first.** Privacy share, augmentation, test data, or
   rebalancing? The goal sets which axis (privacy vs. utility vs. fidelity) you
   optimize and which evaluation gates you must pass.
2. **Profile and build metadata.** Column types, primary/foreign keys, datetime
   formats, constraints (ranges, monotonic, computed columns). SDV needs explicit
   metadata; bad metadata is the #1 cause of bad output.
3. **Split before you fit.** Hold out a real test set *before* training the
   synthesizer so TSTR and MIA are honest.
4. **Pick a generator by data + constraint.** Start with Gaussian copula (fast
   baseline); escalate to CTGAN/TVAE for non-linear structure; consider
   diffusion (TabDDPM) for highest fidelity; switch to a DP synthesizer
   (PATE-GAN, MST/PrivBayes via SmartNoise) the moment a formal privacy guarantee
   is required.
5. **Fit, then enforce constraints.** Apply business rules / valid ranges; reject
   or post-process invalid rows.
6. **Evaluate on all three axes.** Fidelity (KS / TV, correlation similarity),
   utility (TSTR vs. TRTR), privacy (DCR + membership inference). For DP runs,
   report epsilon alongside the metrics.
7. **Iterate against the binding constraint.** If privacy fails, lower epsilon /
   regularize; if utility fails, increase capacity/epochs or change family. Expect
   to trade.
8. **Document.** Generator, hyperparameters, epsilon, seed, metric values, and the
   re-identification-risk assessment — required for any GDPR/anonymisation claim.

## Practical patterns

- **Copula-first ladder.** Baseline with GaussianCopula; only move to CTGAN/TVAE/
  diffusion if fidelity/utility metrics demand it. Don't start with a GAN.
- **Resample inside the fold.** For class imbalance, run SMOTE/ADASYN (or a
  conditional generator) only on the training fold, never before the split.
- **DP by construction, not by hope.** If you need a privacy guarantee, choose a
  DP synthesizer and set epsilon explicitly; do not retrofit privacy onto a
  non-DP generator with post-hoc filtering alone.
- **Conditional sampling for rare classes.** Use CTGAN's conditional generator (or
  SDV conditional sampling) to oversample specific categories deliberately.
- **Report TSTR next to TRTR.** A TSTR score in isolation is meaningless; the gap
  to TRTR is the signal.
- **Keep the metadata in version control.** It is the contract between real and
  synthetic; regenerate from it, don't hand-edit output.

## Anti-patterns

- **Calling synthetic data "anonymous" with no test.** Always run DCR + MIA; for a
  regulatory claim, use DP and document the risk assessment.
- **Evaluating fidelity only.** A perfect KS score can still memorize records
  (great fidelity, terrible privacy) or be useless downstream (poor utility).
- **SMOTE before the split / on test data.** Classic leakage; inflates CV scores.
- **High epsilon = "private."** Large epsilon (e.g., ≫1) often provides little real
  protection; MST/PrivBayes have demonstrated leakage at high epsilon.
- **GAN by default on small/simple data.** Unstable training and mode collapse;
  copula or TVAE is usually a better, faster baseline.
- **Ignoring constraints/keys.** Synthesizers happily emit out-of-range values,
  broken foreign keys, and violated business rules unless constrained.
- **One synthetic draw treated as ground truth.** Sample multiple datasets;
  metrics vary run to run, especially for GANs.

## Troubleshooting

- **Mode collapse / low diversity (GAN).** Switch to TVAE or diffusion; add the
  conditional generator; tune batch size / PacGAN-style packing.
- **Categorical cardinality blows up training time.** High-cardinality columns are
  CTGAN's weak spot — group rare categories, or use copula/CART.
- **DP synthetic data is useless (utility floor).** Epsilon is too small or
  marginals too high-order; raise epsilon within policy, lower marginal order
  (MST), or switch PATE-GAN ↔ marginal-based.
- **TSTR ≫ worse than TRTR.** Fidelity gap in the columns the model relies on;
  inspect per-column shape + correlation reports to find which features drifted.
- **Synthetic rows duplicate real rows (DCR ~0).** Overfitting the synthesizer;
  reduce epochs/capacity, add DP, or increase training data.
- **Invalid rows (out of range, broken keys).** Add SDV constraints / metadata;
  post-process and re-validate.

## References

- Xu, Skoularidou, Cuesta-Infante, Veeramachaneni — *Modeling Tabular Data using
  Conditional GAN (CTGAN)*, NeurIPS 2019.
  https://proceedings.neurips.cc/paper_files/paper/2019/file/254ed7d2de3b23ab10936522dd547b78-Paper.pdf
- Jordon, Yoon, van der Schaar — *PATE-GAN: Generating Synthetic Data with
  Differential Privacy Guarantees*, ICLR 2019.
  https://openreview.net/pdf?id=S1zk9iRqF7
- Kotelnikov et al. — *TabDDPM: Modelling Tabular Data with Diffusion Models*,
  ICML 2023. https://proceedings.mlr.press/v202/kotelnikov23a/kotelnikov23a.pdf
- McKenna, Miklau, Sheldon — *Winning the NIST Contest: MST*, 2021.
  https://arxiv.org/pdf/2301.08844
- Zhang, Cormode, Procopiuc, Srivastava, Xiao — *PrivBayes: Private Data Release
  via Bayesian Networks*, ACM TODS 2017. https://dl.acm.org/doi/10.1145/3134428
- SDV — *Synthetic Data Vault docs & synthesizers* (2024-25).
  https://docs.sdv.dev/sdv | https://github.com/sdv-dev/SDV
- CTGAN library. https://github.com/sdv-dev/CTGAN
- SmartNoise / OpenDP synthesizers (Microsoft + Harvard, 2021-25).
  https://opensource.microsoft.com/blog/2021/02/18/create-privacy-preserving-synthetic-data-for-machine-learning-with-smartnoise/
  | https://docs.smartnoise.org/synth/index.html
- NIST Privacy Collaborative Research Cycle — DP synthetic-data techniques (2025).
  https://pages.nist.gov/privacy_collaborative_research_cycle/pages/techniques.html
- AWS ML Blog — *Evaluate synthetic data: fidelity, utility, privacy* (2023).
  https://aws.amazon.com/blogs/machine-learning/how-to-evaluate-the-quality-of-the-synthetic-data-measuring-from-the-perspective-of-fidelity-utility-and-privacy/
- Frontiers in Digital Health — *Comprehensive evaluation framework for synthetic
  tabular data in health* (2025).
  https://www.frontiersin.org/journals/digital-health/articles/10.3389/fdgth.2025.1576290/full
- *High Epsilon Synthetic Data Vulnerabilities in MST and PrivBayes* (2024).
  https://arxiv.org/html/2402.06699v1
- MDPI Mathematics — *GANs vs SMOTE for imbalanced data* (2023).
  https://www.mdpi.com/2227-7390/11/16/3605 | SMOTified-GAN (2021)
  https://arxiv.org/pdf/2108.03235
- imbalanced-learn / SMOTE & ADASYN guide.
  https://machinelearningmastery.com/smote-oversampling-for-imbalanced-classification/
- Gartner — *Safeguarding Privacy with Synthetic Data* (2024).
  https://www.gartner.com/en/newsroom/press-releases/2024-06-27-safeguarding-privacy-with-synthetic-data
- ICO anonymisation guidance summary, RPC (2025).
  https://www.rpclegal.com/snapshots/data-protection/summer-2025/ico-publishes-new-guidance-on-anonymisation-and-pseudonymisation/
- GDPR Local — *Synthetic Data Under GDPR* (2024-25).
  https://gdprlocal.com/synthetic-data-under-gdpr/
- *Synthetic Data in 2024* review. https://www.timlrx.com/blog/synthetic-data-in-2024-progress-opportunities-challenges/
- synthcity (van der Schaar lab). https://arxiv.org/pdf/2301.07573
