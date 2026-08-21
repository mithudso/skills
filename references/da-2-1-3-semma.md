<!-- hub-reference-banner -->
> **Reference file — part of the `da-2-data-analysis-lifecycle` hub.** Formerly the standalone `da-2-1-3-semma` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-2-1-3-semma
description: |
  Expert knowledge of the SEMMA process framework for data mining, as positioned within the
  Data Analysis Lifecycle taxonomy under Process Frameworks (alongside CRISP-DM, KDD, OSEMN, TDSP).

  TRIGGER: Use when the user asks about SEMMA specifically — what it stands for, how its five phases
  (Sample, Explore, Modify, Model, Assess) work, when to apply SEMMA vs other frameworks, how SEMMA
  maps to SAS Enterprise Miner nodes, or the practical trade-offs of using SEMMA on a data mining
  project. Also trigger when comparing SEMMA to CRISP-DM or KDD as process frameworks.

  SKIP: Questions about CRISP-DM, KDD, OSEMN, or TDSP as primary subjects — those have their own
  skills (da-2-1-1-crisp-dm, da-2-1-2-kdd, da-2-1-4-osemn, da-2-1-5-tdsp). Skip general data
  analysis lifecycle questions not tied to the SEMMA framework specifically. Skip questions about
  sampling as a statistical technique (covered under da-1-4-1-population-vs-sample).
---

# SEMMA: A Process Framework for Data Mining

## Position in the Taxonomy

```
Data Analysis
  └── Data Analysis Lifecycle (Process)
        └── Process frameworks
              ├── CRISP-DM  (da-2-1-1)
              ├── KDD       (da-2-1-2)
              ├── SEMMA     (da-2-1-3)  ← this skill
              ├── OSEMN     (da-2-1-4)
              └── TDSP      (da-2-1-5)
```

SEMMA is a five-phase process framework developed by SAS Institute for organizing data mining work. It is the most tightly coupled to a vendor tool (SAS Enterprise Miner) of the major process frameworks, and SAS itself describes it as "a logical organization of the functional tool set" of Enterprise Miner rather than a universal methodology. [1][2]

---

## What SEMMA Stands For

| Phase | One-line description |
|-------|----------------------|
| **S**ample | Extract a representative subset from the full dataset |
| **E**xplore | Understand data through statistics and visualization |
| **M**odify | Clean, transform, and engineer features for modeling |
| **M**odel | Apply modeling algorithms to the prepared data |
| **A**ssess | Evaluate models and compare candidates |

---

## Phase-by-Phase Detail

### 1. Sample

**What it covers.**
The analyst extracts a working dataset from the full data source. The sample must be large enough to hold sufficient signal for modeling but small enough for efficient processing — the classic size/speed trade-off for large operational databases.

Key activities:
- Identify source tables or files and merge as needed.
- Decide on a sampling strategy (random, stratified, cluster) to preserve the target distribution.
- Partition the dataset into training, validation, and test sets. Enterprise Miner's Data Partition node enforces this split explicitly. [3]
- Filter extreme outliers at the data-ingestion stage rather than later, reducing noise before expensive exploration.

**Pitfall.** Biased sampling is the most damaging error in the whole pipeline. A training set that over-represents a class or time window produces models that appear strong in Assess but fail in production. Always verify that the sample distribution matches the population on the key target variable before proceeding. [4]

---

### 2. Explore

**What it covers.**
The analyst forms a mental model of the data through statistical summaries and visualizations, looking for patterns, anomalies, and inter-variable relationships.

Key activities:
- Univariate analysis: distributions, min/max, quartiles, and frequency counts for each variable.
- Bivariate and multivariate analysis: correlations, scatter plots, and cross-tabulations.
- Variable importance scoring to prioritize candidates for the Model phase.
- Preliminary cluster analysis and association analysis can surface unsuspected groupings. [3]

Enterprise Miner's Explore tab includes nodes such as StatExplore, MultiPlot, and SomAnalysis for these tasks.

**Pitfall.** Confusing exploration with confirmation. Exploration is hypothesis-generating, not hypothesis-testing. Statistical patterns found during Explore are candidates for the model, not conclusions. Treating exploratory correlations as established facts leads to overfitted models. [4]

---

### 3. Modify

**What it covers.**
Raw data rarely enters a model directly. The Modify phase transforms the Explore output into a clean, well-structured analytical dataset (sometimes called an analytical base table, or ABT).

Key activities:
- **Missing value handling:** mean/median/mode imputation, indicator variables for missingness, or deletion, depending on the mechanism (MCAR/MAR/MNAR).
- **Outlier treatment:** Winsorizing, capping, or explicit outlier flags.
- **Variable creation:** derived features via business logic (e.g., ratio of claim amount to policy limit), interaction terms, or polynomial features.
- **Variable transformation:** log-transforms for skewed distributions, binning continuous variables for decision trees, label encoding or one-hot encoding for categorical variables.
- **Dimensionality reduction:** PCA, factor analysis, or self-organizing maps to reduce collinear inputs. Enterprise Miner's SomAnalysis node lives here. [3]

**Pitfall.** Leaking information from the test partition into the training set during this phase — for example, computing mean imputation on the full dataset before splitting. All transformations that depend on the data distribution (means, quantiles, encodings) must be fit on the training partition only, then applied to validation/test. [4][5]

---

### 4. Model

**What it covers.**
The analyst fits one or more predictive or descriptive models to the prepared data. SEMMA is notably technique-agnostic at this phase — the framework does not prescribe which algorithm to use.

Common algorithms applied here:
- Supervised: logistic regression, decision trees (C5.0, CART), neural networks, gradient boosting, support vector machines, random forests.
- Unsupervised: k-means clustering, hierarchical clustering, association rules.

Enterprise Miner organizes this as the Model tab, with nodes for each algorithm family. Multiple models are typically trained in parallel so that the Assess phase can compare them head-to-head. [3]

**Pitfall.** Training only one model type and treating it as the answer. The Assess phase exists precisely to compare candidates; building a single model short-circuits this. A second common error is tuning hyperparameters on the test set rather than the validation set, which yields optimistically biased accuracy estimates. [4]

---

### 5. Assess

**What it covers.**
Models trained on the training partition are evaluated on the held-out validation (and ultimately test) partition, and results are weighed against business requirements.

Assessment outputs:
- **Classification:** ROC/AUC, Kolmogorov–Smirnov statistic, lift charts, profit/cost matrices, misclassification rates.
- **Regression:** RMSE, MAE, R², residual plots.
- **Clustering:** silhouette score, Dunn index, within-cluster sum of squares.
- **Comparison:** Enterprise Miner's Model Comparison node renders side-by-side lift charts and profit analyses for all candidate models. [3]

The Assess step also asks whether the best-performing model meets the original project goal. If not, the loop returns to an earlier phase — additional feature engineering in Modify, a different algorithm in Model, or broader sampling in Sample.

**Pitfall.** Selecting a model solely on statistical metrics while ignoring operational constraints (inference latency, interpretability requirements, fairness across subgroups). A model that scores best on AUC but takes 500ms per prediction may be unusable in a real-time use case. [4][5]

---

## Iterative vs. Sequential Execution

Although presented as a linear acronym, SEMMA is intended to be iterative. [1][4] Practitioners regularly cycle:

- Assess → Modify: poor model performance triggers additional feature engineering.
- Assess → Sample: class imbalance discovered during assessment requires resampling.
- Explore → Sample: anomalies found during exploration reveal data quality issues at the source.

The framework does not specify a formal feedback-loop notation the way CRISP-DM's circular diagram does, but the iterative expectation is present in SAS's own documentation.

---

## Positioning Against Other Process Frameworks

| Dimension | SEMMA | CRISP-DM | KDD |
|-----------|-------|----------|-----|
| Phase count | 5 | 6 | 9 |
| Business understanding phase | No | Yes | Partial (domain knowledge) |
| Deployment guidance | No | Yes | No |
| Tool coupling | SAS Enterprise Miner | Vendor-neutral | Vendor-neutral |
| Adoption (2020 survey) | ~1% | ~49% | ~11% |
| Primary strength | Streamlined model-build cycle | Full project lifecycle | Research-grade rigor |

[1][4][5]

**Key gaps vs. CRISP-DM.**
SEMMA does not include a Business Understanding phase or a Deployment phase. This means the framework says nothing about:
- Defining project success criteria before touching data.
- Stakeholder sign-off on the analytical question.
- Operationalizing the model into a production system.
- Monitoring model drift over time.

For projects where these concerns are substantial, CRISP-DM or TDSP are more complete guides. SEMMA is best understood as covering the core technical loop of a data mining project rather than the full project lifecycle.

---

## When SEMMA Is an Appropriate Choice

SEMMA fits well when:

1. **The organization runs SAS Enterprise Miner.** The framework maps directly to the software's node tabs (Sample, Explore, Modify, Model, Assess). Practitioners who know the tool already understand the framework.
2. **The project scope is narrowly technical.** When business understanding and deployment are handled through separate project management structures, SEMMA's focused scope is an asset rather than a gap.
3. **Quick model prototyping within a known domain.** The five-step checklist is fast to communicate and easy to track without heavyweight documentation.
4. **Teaching introductory data mining.** The acronym is memorable and its phases map cleanly onto the hands-on work of building a first model.

SEMMA is a poor fit when:
- The project requires formal documentation of business goals and success criteria.
- The end goal includes deploying and monitoring a production model.
- The team is not using SAS software (the framework becomes ambiguous outside Enterprise Miner).
- Organizational stakeholders need to see a process map that covers the full project lifecycle.

---

## Worked Example: Customer Churn

To make the phases concrete, consider a telecom company predicting which customers will cancel service within 90 days.

| Phase | What the analyst actually does |
|-------|-------------------------------|
| Sample | Pull a random 20% stratified sample of the 50M-row customer table, maintaining the 4% churn rate. Partition 70/15/15 training/validation/test. |
| Explore | Plot churn rate by contract type, tenure decile, and product bundle. Identify that 38% of customers lack a recorded payment method — flag for Modify. |
| Modify | Impute missing payment-method indicator with a new category "Unknown." Create a feature: days since last service call. Log-transform monthly charges (right skew). One-hot encode contract type. |
| Model | Train logistic regression (baseline), gradient boosted tree, and neural network. Use 5-fold cross-validation on the training partition with validation partition held out. |
| Assess | Compare lift charts at top-10% target depth. GBT achieves 3.8× lift vs. 2.1× for logistic regression. Business team accepts GBT given acceptable latency. |

---

## Common Pitfalls Summary

1. **Biased sample** — training distribution does not match production distribution.
2. **Exploratory findings treated as conclusions** — correlations from Explore promoted to model features without validation.
3. **Target leakage in Modify** — transformations fit on the full dataset, not the training partition alone.
4. **Single-model trap in Model** — only one algorithm built, skipping the comparative value of Assess.
5. **Metrics-only Assess** — statistical performance accepted without checking operational and fairness constraints.
6. **Treating SEMMA as a full project lifecycle** — omitting business understanding or deployment because the framework does not mention them.

---

## Sources

[1] Wikipedia contributors, "SEMMA," *Wikipedia, The Free Encyclopedia*.
    https://en.wikipedia.org/wiki/SEMMA

[2] SAS Institute Inc., "Data Mining and SEMMA," *SAS Enterprise Miner Documentation*, SAS 9.4.
    https://support.sas.com/documentation/cdl/en/emcs/66392/HTML/default/n0pejm83csbja4n1xueveo2uoujy.htm

[3] SAS Institute Inc., "SAS Enterprise Miner: SEMMA Node Reference," ibid. (Enterprise Miner node organization by tab).

[4] Azevedo, A. and Santos, M.F. (2008). "KDD, SEMMA and CRISP-DM: A Parallel Overview."
    *Proceedings of the IADIS European Conference on Data Mining*, pp. 182–185.
    https://www.researchgate.net/publication/268770881_A_Comparative_Study_of_Data_Mining_Process_Models_KDD_CRISP-DM_and_SEMMA

[5] Starburst Data, "SEMMA vs CRISP-DM: The Differences," *Starburst Blog*.
    https://www.starburst.io/blog/semma-vs-crisp-dm/

[6] HandWiki contributors, "SEMMA," *HandWiki*.
    https://handwiki.org/wiki/SEMMA
    (cites Rohanizadeh & Moghadam, 2009, for industrial procedure applications)
