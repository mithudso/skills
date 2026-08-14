<!-- hub-reference-banner -->
> **Reference file — part of the `da-2-data-analysis-lifecycle` hub.** Formerly the standalone `da-2-1-4-osemn` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-2-1-4-osemn
description: |
  Expert knowledge of the OSEMN process framework for data science projects.
  OSEMN (Obtain, Scrub, Explore, Model, iNterpret) is a practitioner-oriented
  five-stage lifecycle introduced by Hilary Mason and Chris Wiggins in their 2010
  "A Taxonomy of Data Science" post.

  TRIGGER: Use when a user asks about OSEMN, the "awesome" data science framework,
  the five-step data science process, or when structuring a data science project
  workflow using the Obtain/Scrub/Explore/Model/Interpret stages. Also trigger
  when contrasting OSEMN with CRISP-DM, KDD, SEMMA, or TDSP in the context of
  choosing a process framework.

  SKIP: Do not use for CRISP-DM (da-2-1-1-crisp-dm), KDD (da-2-1-2-kdd),
  SEMMA (da-2-1-3-semma), or TDSP (da-2-1-5-tdsp) — those frameworks have their
  own dedicated skills. Skip for general EDA technique questions that are not
  about the OSEMN workflow structure. Skip for deployment pipeline, MLOps, or
  production model serving questions; OSEMN does not cover those concerns.
---

# OSEMN — The Data Science Process Framework

## What is OSEMN?

OSEMN is a five-stage process framework that maps the typical workflow of a data
science project. It was introduced in 2010 by data scientists Hilary Mason
(then at bit.ly) and Chris Wiggins (Columbia University) in a blog post titled
"A Taxonomy of Data Science" on the now-archived dataists blog. The acronym is
pronounced "awesome" (or "possum") and stands for:

1. **O** — Obtain
2. **S** — Scrub
3. **E** — Explore
4. **M** — Model
5. **N** — iNterpret

The framework intentionally describes what data scientists *do* rather than how
a business *manages* a project. It is task-centric and practitioner-facing,
making it a natural companion to command-line and scripting workflows.
Jeroen Janssens structured the popular O'Reilly book *Data Science at the
Command Line* around OSEMN's five stages precisely because each stage maps
cleanly to a set of concrete technical tasks. [1][2]

---

## The Five Stages in Detail

### 1. Obtain

**Goal:** Acquire the raw data needed to answer the question.

**Typical activities:**
- Querying relational or NoSQL databases
- Calling REST APIs (e.g., Twitter API, Google Analytics)
- Web scraping (BeautifulSoup, Scrapy)
- Downloading public datasets (Kaggle, government open-data portals)
- Receiving file exports (CSV, Excel, JSON, Parquet)
- Generating synthetic data when real data is unavailable or sensitive

**Key concern:** data provenance — record *where* each dataset came from, the
timestamp of retrieval, and any access conditions. This is essential for
reproducibility and audits. The classic warning applies: "garbage in, garbage
out" (GIGO). A project can fail at this stage if source selection is wrong or
if the data does not actually represent the population of interest. [3][4]

**Automation note:** Mason and Wiggins emphasised automating acquisition through
scripts rather than manual downloads. This matters at scale and when the pipeline
must refresh on a schedule. [2]

---

### 2. Scrub

**Goal:** Transform raw, messy data into an analysis-ready form.

**Typical activities:**
- Identifying and handling missing values (imputation, removal, flagging)
- Removing or resolving duplicates
- Standardising formats (dates, currencies, units, categorical labels)
- Correcting data entry errors and outliers
- Joining or merging data from multiple sources
- Type casting (strings to numerics, etc.)
- Documenting every cleaning decision for reproducibility

**Time investment:** Scrubbing routinely consumes 40–60 % of total project time.
Mason and Wiggins themselves called it "the least sexy part of the analysis
process, but often one that yields the greatest benefits." [3][5]

**Pitfalls:**
- Undocumented cleaning steps make results impossible to reproduce.
- Aggressive outlier removal can bias models.
- Imputing values without understanding the *mechanism* of missingness
  (MCAR / MAR / MNAR) can introduce systematic error.
- Cleaning decisions made silently during EDA (stage 3) are often forgotten
  if not written down; treat scrubbing as a first-class logged step.

**Tooling (Python ecosystem):** pandas, pyjanitor, Great Expectations (validation),
dbt (SQL-based transformations). [4]

---

### 3. Explore

**Goal:** Understand the data's structure, distributions, and relationships before
committing to any model.

This stage corresponds to Exploratory Data Analysis (EDA). There is no hypothesis
being tested and no prediction being evaluated yet; the goal is to build intuition
about the data and surface facts that will inform modelling choices.

**Typical activities:**
- Summary statistics (mean, median, variance, percentiles)
- Distribution plots (histograms, KDE, Q-Q plots)
- Correlation matrices and scatter plots
- Time-series plots, seasonality checks
- Dimensionality reduction for high-dimensional data (PCA, t-SNE)
- Clustering to discover natural groupings
- Outlier and anomaly detection
- Checking statistical assumptions (normality, homoscedasticity)

**Key question to keep asking:** "What patterns violate my assumptions?"
Findings here directly shape the next stage: which features to engineer, which
model families are plausible, and which data quality issues remain. [3][5]

**Pitfall:** Confusing EDA with confirmatory analysis. EDA is generative; results
from EDA cannot be used as independent confirmation of a hypothesis tested on
the same data (this is "data snooping" or "p-hacking").

---

### 4. Model

**Goal:** Apply statistical or machine learning techniques to extract a formal
answer or prediction.

**Typical activities:**
- Feature engineering and selection
- Choosing a model family appropriate to the problem type:
  - Classification → logistic regression, random forest, gradient boosting
  - Regression → linear regression, XGBoost, neural networks
  - Clustering → k-means, DBSCAN, hierarchical
  - Dimensionality reduction → PCA, autoencoders
- Splitting data into train / validation / test sets
- Hyperparameter tuning (grid search, random search, Bayesian optimisation)
- Model evaluation (accuracy, AUC-ROC, RMSE, F1, etc.)
- Error analysis — examining where the model fails

**Key insight from Mason and Wiggins:** "Most of the impact will come from great
features, not great machine learning algorithms." Investing in the Scrub and
Explore stages pays off here. [3]

**Pitfall:** Overfitting through repeated evaluation on the same test set
(the "leaky test set" problem). The OSEMN framework does not prescribe a specific
evaluation methodology; practitioners must impose cross-validation discipline
themselves.

---

### 5. iNterpret

**Goal:** Turn model outputs back into answers to the original question and
communicate them to the people who need to act on them.

**Typical activities:**
- Linking quantitative results back to the business or research question
- Quantifying uncertainty and communicating confidence intervals
- Translating accuracy metrics into business impact (e.g., "92 % accuracy
  reduces customer churn cost by $1.2 M / year" rather than just the number)
- Building visualisations and narratives for non-technical stakeholders
- Deciding whether findings are actionable or whether more data is needed
- Documenting limitations and conditions under which conclusions hold

**Three questions the iNterpret stage must answer (from original framework):**
1. Is the data representative? (selection bias check)
2. Are the right features included? (feature relevance)
3. What is the hypothesis space? (scope of conclusions)

**Pitfall:** Reporting raw model performance without business context. Stakeholders
need to know what the result *means for them*, not just the metric value. [4][5]

---

## How OSEMN Fits in the Taxonomy

OSEMN sits alongside CRISP-DM, KDD, SEMMA, and TDSP as a **process framework**
within the Data Analysis Lifecycle. Each framework covers roughly the same
territory; their differences are in emphasis and origin:

| Framework | Origin | Emphasis | Stages |
|---|---|---|---|
| OSEMN | Mason & Wiggins, 2010 | Practitioner tasks, scripting/automation | 5 |
| CRISP-DM | Consortium (SPSS, NCR, Daimler), 1996 | Business context, iterative evaluation | 6 |
| KDD | Fayyad et al., 1996 | Research / knowledge discovery | 9 |
| SEMMA | SAS Institute | Tool-guided workflow within SAS Enterprise Miner | 5 |
| TDSP | Microsoft, 2016 | Team collaboration, DevOps integration | 5 |

OSEMN is the most lightweight and technology-agnostic. Unlike CRISP-DM, it has
no explicit "Business Understanding" phase — a common criticism (see Limitations
below). Unlike KDD, it avoids academic formalism. Unlike SEMMA, it is not tied
to any vendor toolchain. [2][3]

---

## When to Use OSEMN

**Good fit when:**
- The team is small (one to a few data scientists) working on a contained problem
- The work is primarily analytical / scripting-oriented rather than engineering
- A shared vocabulary for project stages is needed without heavyweight process overhead
- Teaching or onboarding practitioners to data science workflow concepts
- The output is a one-time analysis or report rather than a continuously deployed model

**Less suitable when:**
- The project requires explicit business stakeholder alignment before data work begins
  (CRISP-DM's Business Understanding phase handles this better)
- Model deployment and ongoing production monitoring are central concerns
  (OSEMN implicitly assumes a one-time deliverable; it does not address MLOps)
- Large cross-functional teams with defined roles need coordination structures
  (TDSP provides project templates and role definitions; OSEMN does not)
- The workflow is highly iterative, requiring frequent jumps between stages
  (OSEMN reads as linear/waterfall; practitioners must impose their own iteration)

---

## Known Limitations

1. **No problem-definition phase.** OSEMN begins with Obtain, but a prior step —
   "Should this project happen at all, and what outcome am I driving?" — is absent.
   The AOSEMN extension proposed by some practitioners adds an "Ask" stage before
   Obtain to address this gap. [5]

2. **No deployment phase.** After iNterpret, the model is presumably done; there
   is no guidance for putting it into production, monitoring drift, or retraining.

3. **Linear / waterfall framing.** Real projects cycle back from Explore to Scrub,
   or from Model back to Explore. The acronym implies a fixed order that rarely
   holds in practice.

4. **No team structure.** OSEMN says nothing about roles, handoffs, or collaboration,
   which matters as data science has become more of a team discipline.

---

## Worked Example: Retail Foot-Traffic Analysis

| Stage | What happens |
|---|---|
| Obtain | Download 6 months of IoT sensor foot-traffic counts from store API; export transaction records from the data warehouse |
| Scrub | Fix 12 % missing sensor readings (forward-fill within business hours); deduplicate transaction rows; align timestamps to a common timezone |
| Explore | Plot hourly traffic distributions; find a 3 pm congestion spike on weekdays; correlate traffic with weather data |
| Model | Fit a time-series model to forecast traffic; optimise shelf-placement using the predicted congestion pattern |
| iNterpret | Recommend adjusted staffing schedule; project 15 % reduction in checkout queue time; present to operations manager with before/after scenario charts |

---

## Sources

[1] Janssens, J. (2021). *Data Science at the Command Line, 2nd Edition*, Chapter 11.
O'Reilly Media. https://www.oreilly.com/library/view/data-science-at/9781492087908/ch11.html

[2] Data Science PM. "OSEMN Data Science Life Cycle."
https://www.datascience-pm.com/osemn/

[3] DataRundown. "Data Science Life Cycle: CRISP-DM and OSEMN Frameworks."
https://datarundown.com/data-science-life-cycle/

[4] Faisal, S. "OSEMN Framework: The 5-Step Data Science Pipeline Explained."
https://faisalsir.com/blog/data-analytics/osemn-framework-the-data-scientists-obsessive-process/

[5] Harshvardhan. "Data Science Process Frameworks."
https://harshvardhan.blog/data-science-process-frameworks
