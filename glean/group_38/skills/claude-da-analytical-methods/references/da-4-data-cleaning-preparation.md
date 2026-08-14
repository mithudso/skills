<!-- hub-reference-banner -->
> **Reference file — part of the `da-analytical-methods` hub.** Formerly the standalone `da-4-data-cleaning-preparation` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-4-data-cleaning-preparation
description: >-
  Data cleaning and preparation as the disciplined transformation of raw,
  messy data into a model-ready table. Covers missing data theory (MCAR/MAR/
  MNAR per Rubin) and the imputation ladder (listwise, mean/median, KNN,
  IterativeImputer/MICE, miceforest, missingness indicators); outlier
  detection (z-score, IQR / Tukey fences, Isolation Forest, DBSCAN, LOF) and
  the remove/winsorize/transform decision; type coercion (pandas to_numeric,
  to_datetime, convert_dtypes, encoding detection via charset-normalizer,
  timezone normalization with DST-ambiguity handling); deduplication (exact,
  fuzzy via Levenshtein/Jaro-Winkler with RapidFuzz, probabilistic record
  linkage via Splink/Fellegi-Sunter, blocking keys); scaling and
  transformation (StandardScaler, MinMaxScaler, RobustScaler,
  PowerTransformer with Box-Cox vs Yeo-Johnson, log transforms); categorical
  encoding (one-hot, label/ordinal, target with Bayesian smoothing and CV to
  prevent leakage, frequency, binary); schema validation (Pandera, Great
  Expectations, Pydantic, JSON Schema); feature engineering basics
  (interaction terms, polynomial features, KBinsDiscretizer, date
  decomposition); text cleaning (Unicode NFC/NFKC normalization, case
  folding, tokenization, stopwords, lemmatization vs stemming); imbalanced
  datasets (SMOTE, ADASYN, Borderline-SMOTE, SMOTE-Tomek, undersampling,
  class_weight, imbalanced-learn); and the canonical anti-patterns
  (preprocessing-before-split leakage, silent dropna, over-cleaning,
  ignoring missingness mechanism, target encoding without CV).
  TRIGGER: someone asks "how do I handle missing data", "MCAR vs MAR vs
  MNAR", "should I impute or drop", "KNN vs MICE imputation", "how do I
  detect outliers", "Tukey fences vs z-score vs Isolation Forest", "should I
  remove outliers or winsorize", "how to parse messy dates / mixed
  encodings", "how to dedupe a customer table", "fuzzy matching library",
  "Splink record linkage", "StandardScaler vs RobustScaler", "Box-Cox vs
  Yeo-Johnson", "one-hot vs target encoding", "target encoding leakage",
  "high cardinality categorical", "Pandera vs Great Expectations", "schema
  validation in a pipeline", "feature engineering binning interaction
  terms", "Unicode NFC NFD normalization", "SMOTE ADASYN class imbalance",
  "class_weight balanced", "data leakage during preprocessing", "fit on
  train transform on test", "should I drop rows with any NaN", or asks for
  a cleaning playbook before EDA or modeling.
  SKIP: data acquisition / sampling design (use da-3-data-acquisition-
  sampling — that is where you decide what to collect and how to draw the
  sample; this skill assumes raw data is in hand). SKIP exploratory data
  analysis, profiling, and visualization (use da-5-exploratory-data-
  analysis — that comes after cleaning). SKIP statistical-inference theory
  (use da-1-4-statistical-inference-foundations). SKIP modeling itself —
  regression mechanics, ML model building, hyperparameter tuning belong to
  their own technique skills. SKIP MongoDB-specific schema validation (use
  mongodb-schema-design). SKIP SQL/aggregation cleaning that is really
  pipeline ETL (use mongodb-aggregation-pipeline or the relevant ETL
  skill). SKIP Excel-only cleaning (use xlsx). SKIP MongoDB type coercion
  inside drivers (use mongodb-bson-types). SKIP NLP modeling work beyond
  preprocessing (use a downstream NLP modeling skill). This skill is the
  da-curriculum section 4: raw → analysis-ready table.
---

# Data Cleaning and Preparation

## Overview

Data cleaning and preparation is the disciplined transformation of raw,
messy input data into a model-ready, analysis-ready table. It is the work
that sits between **da-3 (data acquisition and sampling)** and **da-5
(exploratory data analysis)**. The literature is consistent that this stage
absorbs 60-80 percent of a data analyst's time, and that errors made here
silently compound through every downstream step — missing-data mechanism
mis-specification biases every estimator; preprocessing-before-split leaks
the test distribution into the training pipeline; an overly aggressive
`dropna()` deletes the signal you came to find.

The discipline has four cores:

1. **Diagnose the mess.** Profile missingness, outliers, types, encodings,
   duplicates, schema drift. Decide whether a problem is a systematic
   signal or random noise *before* touching the data.
2. **Make principled fixes.** Each correction (imputation, scaling,
   encoding, deduplication) carries assumptions. Track them.
3. **Prevent leakage.** All transformations that "learn" from data
   (imputers, scalers, encoders, oversamplers) must be fitted on training
   data only, then applied to validation/test. Pipelines exist to make
   this mechanical.
4. **Validate.** Treat the cleaned dataframe as a contract. Schema
   validation (Pandera, Great Expectations) catches the next batch's
   regression before it kills the model.

This skill is curriculum-level. It does not replace technique-specific
skills (e.g. NLP modeling, MongoDB schema design); it gives the canonical
moves, the canonical failure modes, and the bibliography.

## Core Concepts

### 1. Missing Data: Mechanism Matters More Than Method

Donald Rubin's 1976 classification (Rubin, *Biometrika*) is still the
language of missing-data work. The three mechanisms determine which
imputation strategies are *valid* — not just convenient.

| Mechanism | Definition | Practical signal | Valid strategies |
|---|---|---|---|
| **MCAR** — Missing Completely At Random | P(missing) is independent of observed and missing values. Pure random. | Sensor dropout, random server crash. | Listwise deletion is unbiased (but inefficient). Mean / median imputation is unbiased for the mean but understates variance. |
| **MAR** — Missing At Random | P(missing) depends only on **observed** variables. | Survey respondents over 65 skip a salary question; we know age. | Conditional imputation: KNN, MICE, IterativeImputer. **Multiple imputation is the standard.** |
| **MNAR** — Missing Not At Random | P(missing) depends on the missing value itself. | High earners refuse to disclose salary. | No fully model-free fix. Requires explicit missingness model, sensitivity analysis, or domain-specific imputation. |

Three field rules:

- **You cannot test MAR vs MNAR from observed data alone.** Little's MCAR
  test can sometimes reject MCAR; nothing distinguishes MAR from MNAR
  without external information.
- **Include "many predictors of missingness" in the imputation model.**
  This makes the MAR assumption more plausible and reduces sensitivity to
  violations (van Buuren, *Flexible Imputation of Missing Data*).
- **Always add a missingness indicator** alongside imputation when the
  *fact* of missingness is potentially informative (e.g. critical-care
  data where "lab not ordered" is itself a signal).

#### The imputation ladder

1. **Listwise deletion** — drop any row with a missing cell. Defensible
   under MCAR with abundant data; throws away signal otherwise. The
   pandas `dropna()` default (`how='any'`, no `subset`) is the most
   common silent-data-loss anti-pattern.
2. **Constant / mean / median / mode** — `SimpleImputer`. Median is more
   robust to skew than mean. **Understates variance**; do not use for
   inference without combining with a missingness indicator.
3. **KNN imputation** — `KNNImputer`. Imputes each missing feature from
   the `n_neighbors` nearest rows that have the feature. Distance
   metrics: nan-euclidean (default). Great for small numerical datasets;
   computationally expensive for large ones.
4. **Iterative / chained-equations (MICE)** — `IterativeImputer`. Models
   each feature with missing values as a function of the others, cycles
   through features until convergence. Default estimator: BayesianRidge.
   Can swap in random forest (→ **miceforest** library). Scales
   `O(k·n·p³·min(n,p))` in features p; use `n_nearest_features` to tame.
5. **Multiple imputation (true MI)** — Repeat the imputation `m` times
   (typically 5-20), fit the downstream model on each completed dataset,
   pool estimates via Rubin's rules. The only method that *correctly*
   propagates imputation uncertainty into inference.

Benchmarks: in immunization-records demographics, MICE and miceforest
best preserved demographic proportions vs IterativeImputer; MICE took
~14 hours, IterativeImputer 2 minutes, miceforest 10 minutes for 15
imputations. Choose by dataset size and inference-vs-prediction
requirement.

### 2. Outlier Detection and Treatment

Outliers are observations distant from the bulk of the distribution.
The first question is never "how do I remove them" — it is **"why do
they exist."** A negative age is an error; a $50M transaction in a fraud
detection table is the entire reason for the project.

#### Detection methods

| Method | Mechanism | Best for | Caveats |
|---|---|---|---|
| **z-score** | $\|x - \mu\| / \sigma > 3$ | Univariate, approximately normal data | Mean/std are themselves outlier-sensitive; thresholds break under skew |
| **IQR / Tukey fences** | $x < Q_1 - 1.5 \cdot IQR$ or $x > Q_3 + 1.5 \cdot IQR$; extreme fences at 3·IQR | Univariate, robust to non-normal | Quartile-based; minimum sample ~ 30 |
| **Modified z-score (MAD)** | Uses median and median-absolute-deviation; threshold ~3.5 | Robust univariate alternative to z | Less standard; reach for it under heavy skew |
| **Isolation Forest** | Random partitioning; outliers isolated in few splits | Multivariate, large datasets, high dimensions | Hyperparameters (`contamination`) matter; not deterministic |
| **DBSCAN** | Density clustering; "noise" points = outliers | Multivariate clusters of varying shape | `eps` and `min_samples` sensitive; struggles in high dimensions |
| **Local Outlier Factor (LOF)** | Local density deviation | Multivariate, local-density problems | Computationally heavier than IF |

Practical heuristic: **z-score and IQR for univariate diagnostics;
Isolation Forest for multivariate.** Run several methods; treat
agreement as a signal.

#### Treatment: remove, winsorize, or transform

There is no universal best. The three choices:

- **Remove** — only when the value is *demonstrably* an error
  (negative age, future timestamp, impossible category). Otherwise this
  is a data-deletion anti-pattern.
- **Winsorize / cap** — pin values to a percentile (e.g. 1st / 99th).
  Preserves row count, retains rank order, attenuates leverage. Useful
  for linear models that are sensitive to extreme leverage. Less
  aggressive than removal but still distorts the tails.
- **Transform** — log, Box-Cox, Yeo-Johnson. Often the right move first
  for variables that grow exponentially (revenue, latency, file size,
  income). Tree-based models (RF, XGBoost) tolerate raw outliers
  better and rarely need scaling or transformation — choose treatment
  based on the downstream model family.

Decision tree:

```
Is the outlier an error?            → remove
Linear / distance-based model
  and outliers have high leverage?  → winsorize or log/Box-Cox
Heavy-tailed positive variable?     → log or Box-Cox (positive only) /
                                      Yeo-Johnson (allows zero/negative)
Tree-based downstream model?        → usually leave alone
```

### 3. Type Coercion, Date Parsing, Encoding, Timezones

pandas exposes four type-coercion entry points: `pd.to_numeric()`,
`pd.to_datetime()`, `.astype()`, and the modern `.convert_dtypes()`
(returns nullable Int64 / string dtype / pd.NA).

#### String → numeric
```python
pd.to_numeric(series, errors='coerce')  # invalid → NaN
pd.to_numeric(series, errors='raise')   # default, raises
pd.to_numeric(series, errors='ignore')  # returns input on failure
```
Default `errors='raise'` is correct for production; `coerce` is correct
for exploratory work where you want to count failures.

#### Date parsing
- Always pass `format=` when the format is known. The auto-inferred
  format **changes per row** ("dynamic format") which silently misparses
  mixed-style strings.
- `errors='coerce'` forces unparseable / out-of-bounds dates to `NaT`.
- For multi-format columns, use `format='mixed'` (pandas ≥ 2.0) or
  parse explicitly per pattern.

#### Encoding detection
For unknown encodings, prefer **`charset-normalizer`** (the modern
default; backs `requests`) or **chardet 7** (rewritten, 47× faster, 99.3%
accuracy on 2,517-file benchmark). chardet 5/6 → 7 is drop-in.
Multi-encoding files (concatenated logs, sloppy migrations) defeat every
detector — split by section first.

#### Timezone normalization
The two-step pattern is:
```python
s.dt.tz_localize('US/Pacific', ambiguous='infer',
                  nonexistent='shift_forward')
s.dt.tz_convert('UTC')
```
Pitfalls:

- **Ambiguous fall-back times** (02:30 in CET happens twice on DST end).
  Choices: `'infer'` (sequential), `bool-ndarray`, `'NaT'`, `'raise'`.
- **Nonexistent spring-forward times** (02:30 in CET doesn't exist on
  DST start). Choices: `'shift_forward'`, `'shift_backward'`, `'NaT'`,
  `'raise'`.
- **pytz vs dateutil vs zoneinfo** definitions diverge on edge zones;
  Python 3.9+ `zoneinfo` is now the reference. Persist timestamps as
  UTC; convert only for display.

### 4. Deduplication

Three regimes:

1. **Exact duplicates** — `df.drop_duplicates(subset=[...])`. Trivial
   but underused; always pass a `subset` for production tables.
2. **Fuzzy / approximate duplicates** — string-similarity metrics:
   - **Levenshtein** — edit distance (insertions / deletions /
     substitutions). General purpose.
   - **Jaro-Winkler** — favors common prefixes; classic for short
     person/company names with typos and transpositions.
   - **Damerau-Levenshtein** — adds transpositions to Levenshtein.
   - **Cosine / Jaccard over n-grams** — token-set similarity.
   - **Phonetic** (Soundex, Metaphone, Double Metaphone) — sound-alikes
     ("Smith" / "Smyth").
   - Libraries: **RapidFuzz** (C++-backed, faster than the legacy
     fuzzywuzzy by orders of magnitude); **jellyfish**; **textdistance**.
3. **Probabilistic record linkage** — for cross-source entity
   resolution.
   - **Fellegi-Sunter** (1969) is the canonical model: agreement on
     each comparison field gives a likelihood ratio (m-probability /
     u-probability); per-field log-likelihoods sum to a match score.
   - **Splink** (MoJ open source) is the production-grade
     implementation: DuckDB / Spark / Athena backends, EM algorithm
     for parameter estimation, scales to 100M+ records (1M records in
     ~2 minutes on a laptop with DuckDB).
   - **dedupe.io** is the alternative Python library; active-learning
     UI for labeling pairs.

#### Blocking is mandatory at scale

Comparing every record to every other record is O(n²). **Blocking
keys** (postal code, soundex of last name, first 3 chars of email)
restrict candidate pairs. Trade off:

- Tight blocks → fast, but miss matches that disagree on the block key.
- Loose blocks → high recall, but cost explodes.

Multiple block keys (union of candidate sets) is the standard cure.
Splink's "blocking rules" let you express this declaratively.

### 5. Scaling, Standardization, Transformation

The choice of scaler is dictated by:

1. **Downstream model's sensitivity** — distance-based (k-NN, SVM,
   k-means, PCA), gradient-based with L1/L2 regularization, and neural
   nets *require* scaling. Tree-based models do not.
2. **Outlier presence** — outliers wreck mean/std-based scalers.
3. **Shape of the distribution** — power transforms attack skew and
   heavy tails directly.

| Scaler / Transform | Formula | Use when | Avoid when |
|---|---|---|---|
| **StandardScaler** | $z = (x - \mu) / \sigma$ | Distance / gradient models, roughly normal features | Heavy outliers, skewed |
| **MinMaxScaler** | $(x - x_{min}) / (x_{max} - x_{min})$, default [0,1] | Bounded outputs needed (images, sigmoid activations) | New data can fall outside fit range; outliers compress majority |
| **MaxAbsScaler** | $x / |x|_{max}$ | Sparse data (preserves zeros) | Heavy outliers |
| **RobustScaler** | $(x - \text{median}) / IQR$ | Outliers present and not removable | Approximately normal data (use StandardScaler) |
| **Normalizer** | Scales each *row* to unit norm | Cosine-similarity downstream | Most tabular tasks |
| **PowerTransformer (Box-Cox)** | Estimates $\lambda$ by MLE | Strictly positive, skewed → Gaussian | Zero / negative values |
| **PowerTransformer (Yeo-Johnson)** | Generalization of Box-Cox | Any sign, skewed → Gaussian | Already symmetric |
| **QuantileTransformer** | Maps to uniform or normal via empirical CDF | Heavy-tailed, robust to outliers | Tiny samples (CDF unstable) |
| **log / log1p** | $\log(x)$ / $\log(1+x)$ | Right-skewed, positive (latency, revenue, counts) | Negatives (use log1p only for non-negatives) |

Always **fit on train, transform on test.** Pipelines enforce this.

### 6. Encoding Categorical Variables

| Method | Output cardinality | Order assumed | Leakage risk | Good for |
|---|---|---|---|---|
| **One-Hot** | k columns | No | None | Nominal, low cardinality (k < ~15), linear / NN models |
| **Ordinal / Label** | 1 integer column | Yes (assumes ordered) | None | True ordinal (size, education); tree models can tolerate fake order |
| **Target / Mean** | 1 numeric column | No | **High** — must use CV folds | High-cardinality nominal (zip, brand) |
| **Frequency / Count** | 1 numeric column | No | None | High-cardinality where frequency itself signals |
| **Binary** | log₂(k) columns | No | None | High-cardinality, want compactness over interpretability |
| **Hashing** | fixed h columns | No | None | Streaming / huge cardinality; tolerates collisions |
| **Embedding** | dense d-vector | No | Trained, not raw | Deep models with sufficient data |

**Target encoding leakage is the single most common cleaning mistake.**
The naive recipe — `df['cat_enc'] = df.groupby('cat')['y'].transform('mean')`
— exposes each row to its own target. Fix:

- **K-fold target encoding** — compute per-category mean using only
  *other* folds; encode the held-out fold from that.
- **Bayesian smoothing** — blend per-category mean with global mean,
  weighted by category count. Reduces overfit on rare categories.
- scikit-learn 1.3+ ships `TargetEncoder` with built-in CV and smoothing;
  prefer it over hand-rolled implementations.

### 7. Schema Validation

A cleaned dataframe is a contract. Schema validation lets you
*declare* the contract and enforce it on every batch.

| Tool | Scope | Engine support | Where it fits |
|---|---|---|---|
| **Pydantic v2** | Row/object level, type-driven | Python objects (not DataFrames natively) | API boundaries; pair with FastAPI |
| **Pandera** (0.29, Jan 2026) | DataFrame schemas, statistical hypotheses | pandas, Polars, Dask, Modin, PySpark, Ibis | In-process ML pipelines, unit tests |
| **Great Expectations** | DataFrame + multi-engine | pandas, Spark, SQL backends | Pipeline checkpoints, Data Docs, governance |
| **JSON Schema** | Record-level, language-agnostic | Any JSON-producing service | Cross-system contracts |
| **Soda Core / SodaCL** | SQL-first declarative checks | Warehouses (Snowflake, BigQuery, Redshift) | dbt / orchestration integration |

Selection rule of thumb:

- Pure Python / DataFrames / type hints / "validations as unit tests"
  → **Pandera**.
- Multi-engine pipelines (Spark + SQL + pandas), shared expectation
  suites, human-readable docs, governance → **Great Expectations**.
- API boundary → **Pydantic**.
- Cross-org JSON contracts → **JSON Schema**.
- Many teams stack all three: Pydantic at the API edge, Pandera in
  ML / notebook code, GE in production warehouses.

### 8. Feature Engineering Basics

Out of scope is the entire ML-feature-store ecosystem; in scope are
the cleaning-adjacent transformations:

- **Binning / discretization** — `KBinsDiscretizer` (uniform, quantile,
  kmeans strategies); capture non-linearity for linear models, reduce
  outlier leverage. Risk: throws away granularity; over-binning
  destroys signal.
- **Interaction terms** — multiplicative or boolean conjunctions of
  features (e.g. `num_bedrooms * sqft`). Driven by hypothesis, not
  exhaustive search.
- **Polynomial features** — `PolynomialFeatures(degree=2)` produces
  squares, cubes, and cross-terms. Explodes feature space; regularize.
- **Date decomposition** — from a timestamp extract year, month,
  day-of-week, hour, is_weekend, days_since_epoch, business-day
  count, holidays. "House age" (current_year − year_built) is a more
  informative encoding than the raw year.
- **Cyclical encoding** — for hour-of-day, day-of-year, etc., use
  `sin(2π·t/T)` and `cos(2π·t/T)` so distance respects circularity
  (23:00 is close to 01:00).

### 9. Text Cleaning

For modeling, "text cleaning" is a contracting concept — modern
transformer pipelines tolerate raw text far better than classical NLP
pipelines did. Even so, the canonical text-cleaning pipeline is:

1. **Unicode normalization** — choose **NFC** for "composed" form (most
   storage / display) or **NFKC** for "compatibility composed"
   (collapses width variants, ligatures, etc., for search). NFD/NFKD
   are decomposed forms used by some downstream tools. The trap:
   "é" can be one codepoint (U+00E9) or two (U+0065 + U+0301);
   Python `str` equality treats them as different. Always normalize
   before comparing.
2. **Case folding** — `str.casefold()` is stronger than `str.lower()`
   for non-ASCII (e.g. "ß" → "ss").
3. **Whitespace, control characters, zero-width chars** — strip
   U+200B, U+FEFF, control codes.
4. **Tokenization** — word-level (NLTK, spaCy), subword (BPE,
   WordPiece for BERT-family), or character. Modern LLM pipelines use
   subword tokenizers from the model card; do not hand-roll.
5. **Stopword removal** — task-dependent. Helps bag-of-words /
   topic modeling; harms transformer fine-tuning and sentiment /
   sarcasm tasks where stopwords carry meaning.
6. **Stemming vs lemmatization** — stemming (Porter, Snowball) crudely
   chops suffixes; lemmatization (spaCy, NLTK WordNet) uses vocabulary
   and POS for proper base form. Lemmatization is slower, more
   accurate; modern pipelines often skip both in favor of subword
   tokenization that handles morphology implicitly.

For full NLP modeling work, defer to a downstream NLP modeling skill —
this section is preprocessing only.

### 10. Imbalanced Datasets

Imbalance is a model-side problem more often than a data-side
problem; the cleaning surface is the *resampling* choice.

| Technique | Family | What it does | Caveats |
|---|---|---|---|
| **Random oversampling** | Oversample | Duplicate minority rows | Pure overfit risk on minority |
| **SMOTE** | Oversample (synthetic) | Linear interpolation between minority points and their k-NNs | No noise handling; can synthesize through majority regions |
| **Borderline-SMOTE** | Oversample (synthetic) | Only synthesize near the decision boundary | More targeted than vanilla SMOTE |
| **ADASYN** | Oversample (synthetic) | Density-weighted: synthesize more for hard-to-learn minority points | More aggressive than SMOTE; can overfit to noise |
| **SMOTE-Tomek** | Hybrid | SMOTE oversample, then Tomek-link removal of boundary pairs | Cleans the synthetic noise; common production default |
| **SMOTE-ENN** | Hybrid | SMOTE then Edited-Nearest-Neighbors cleaning | Aggressive cleaner |
| **Random undersampling** | Undersample | Drop majority rows | Information loss; cheap; works when majority is huge |
| **Tomek links** | Undersample | Remove pairs on the boundary | Cleans, doesn't balance fully |
| **NearMiss** | Undersample | Selects majority points near minority | Picks informative majority rows |
| **class_weight='balanced'** | No resampling | Loss weight $\propto 1/n_{class}$ in training | Often *the* right answer; no preprocessing change needed |

Two iron rules:

- **Resample only the training fold, never the test fold.** All
  imbalanced-learn samplers expose `fit_resample(X_train, y_train)`;
  test set must keep the true imbalance to give an honest score.
- **Try `class_weight='balanced'` before SMOTE.** It is leak-proof,
  has no hyperparameters, and is competitive with sophisticated
  oversamplers in many tabular benchmarks. Recent comparative work
  (XGBoost on financial distress) found SMOTE-Tomek and
  Bagging-SMOTE marginally best; class_weight remained competitive.

Library: **`imbalanced-learn`** (pip `imbalanced-learn`, conda-forge),
scikit-learn-compatible.

### 11. Anti-Patterns

The catalogue of cleaning sins, ranked by frequency:

1. **Preprocessing before train/test split.**
   - Symptom: scaler / imputer / encoder fit on the full dataset, then
     split. The test set has leaked into the transformation.
   - Fix: split first; wrap all data-learning transformations in a
     `Pipeline`; use `cross_validate` / `GridSearchCV` on the pipeline.
2. **Target encoding without cross-validation.**
   - Symptom: per-category mean computed on all rows; each row sees
     its own y. Train score soars, holdout collapses.
   - Fix: K-fold target encoding or scikit-learn 1.3+ `TargetEncoder`.
3. **Silent `dropna()`.**
   - Symptom: `df.dropna()` deletes 80% of rows because one
     non-critical column is sparse.
   - Fix: always pass `subset=[...critical columns]`; profile
     missingness before dropping.
4. **Mean / median imputation with no missingness indicator.**
   - Symptom: variance is understated; downstream model treats the
     imputed mean as if it were observed.
   - Fix: pair imputation with `MissingIndicator`; for inference, use
     multiple imputation.
5. **Removing outliers without diagnosis.**
   - Symptom: $|z| > 3$ stripped wholesale; key signal (fraud, churn,
     anomaly) deleted.
   - Fix: confirm errors vs legitimate extreme observations; prefer
     transform / winsorize unless the value is impossible.
6. **Mean-encoding before splitting.** Same as (2) for a different
   encoder.
7. **Resampling (SMOTE / undersampling) on the full dataset before
   splitting.**
   - Symptom: synthetic minority points in the test fold; test score
     is contaminated.
   - Fix: resample only training fold; `imblearn.Pipeline`.
8. **Imputing with the global mean, then computing correlations.**
   - Mean imputation collapses variance and biases correlations toward
     zero.
9. **Ignoring the missingness mechanism.**
   - Treating MNAR data with MAR-only methods produces biased
     estimates. Document the assumed mechanism.
10. **Type coercion with `errors='ignore'` silently keeping strings in
    a numeric column.**
    - Use `coerce` (produces NaN you can audit) or `raise` (fail fast).
11. **Auto-inferring date formats across mixed-format columns.**
    - pandas changes the inferred format per row. Pass an explicit
      `format=` or use `format='mixed'`.
12. **Dropping duplicates without `subset`.**
    - Drops only true row-clones; misses customer-table dupes that
      differ in trivial columns (timestamps, IDs).

## Tools and Frameworks

| Tool | Category | When to reach for it |
|---|---|---|
| **pandas** (`to_numeric`, `to_datetime`, `convert_dtypes`, `dropna`, `fillna`, `drop_duplicates`) | Core dataframe ops | Default for tabular Python |
| **NumPy** | Numerical core | Lower-level array operations |
| **scikit-learn `sklearn.impute`** (`SimpleImputer`, `KNNImputer`, `IterativeImputer`, `MissingIndicator`) | Imputation | Standard sklearn pipeline imputation |
| **`miceforest`** | Imputation | Fast random-forest MICE for large datasets |
| **`statsmodels.imputation.mice`** | Imputation | Proper multiple imputation with Rubin pooling for inference |
| **scikit-learn `sklearn.preprocessing`** (`StandardScaler`, `MinMaxScaler`, `MaxAbsScaler`, `RobustScaler`, `Normalizer`, `PowerTransformer`, `QuantileTransformer`, `KBinsDiscretizer`, `PolynomialFeatures`, `OneHotEncoder`, `OrdinalEncoder`, `TargetEncoder` ≥ 1.3) | Scaling, transform, encoding | Standard sklearn pipeline transforms |
| **`category_encoders`** | Categorical encoding | Wider menu than sklearn (binary, hashing, James-Stein, leave-one-out, M-estimator) |
| **`charset-normalizer` / `chardet 7`** | Encoding detection | Reading files of unknown encoding |
| **`python-dateutil` / `pendulum` / `zoneinfo`** | Date/timezone | Date parsing; modern Python prefers stdlib `zoneinfo` |
| **`RapidFuzz`** | Fuzzy string matching | Levenshtein / Jaro-Winkler at speed (C++ backend) |
| **`jellyfish`** | String similarity & phonetics | Levenshtein, Damerau, Jaro, Soundex, Metaphone |
| **`recordlinkage`** | Record linkage | Full pipeline: blocking + compare + classify |
| **`Splink`** (MoJ) | Probabilistic record linkage | Fellegi-Sunter at scale; DuckDB / Spark / Athena |
| **`dedupe`** (dedupe.io library) | Probabilistic dedupe | Active-learning UI, smaller datasets |
| **Sklearn `IsolationForest`, `LocalOutlierFactor`, `EllipticEnvelope`** | Outlier detection | Multivariate ML-based outlier scoring |
| **`PyOD`** | Outlier detection | 40+ outlier methods (classical + DL) in one API |
| **`imbalanced-learn`** | Imbalanced data | SMOTE, ADASYN, Borderline-SMOTE, SMOTE-Tomek/ENN, undersamplers, `imblearn.Pipeline` |
| **`pandera`** | DataFrame validation | Type-hinted pandas / Polars / PySpark schemas |
| **`Great Expectations`** | Data quality | Multi-engine pipelines, governance, Data Docs |
| **`pydantic` v2** | Object-level validation | API boundaries, FastAPI, config |
| **`Soda Core` / `SodaCL`** | SQL-first data quality | Warehouse checks in dbt / orchestration |
| **`Monte Carlo` / `Bigeye` / `Anomalo`** | Data observability | Production monitoring, lineage, anomaly alerting |
| **`OpenRefine`** | Interactive cleaning | Small-data exploratory cleaning; faceting, clustering |
| **`Trifacta` (now Alteryx Designer Cloud)** | Visual data prep | Business-analyst self-service |
| **`spaCy` / `NLTK` / `regex`** | Text cleaning | Tokenization, lemmatization, regex with Unicode support |
| **`ftfy`** | Text cleaning | Fixes mojibake (UTF-8-decoded-as-Latin-1, etc.) |
| **`unicodedata` (stdlib)** | Text normalization | NFC/NFD/NFKC/NFKD |

## Methodology: A Cleaning Workflow

The defensible workflow:

```
0. Snapshot raw data (immutable copy + hash).
1. Profile (shape, dtypes, % missing per column,
   distinct counts, distribution sketches).
2. Decide what each column is supposed to mean
   (data dictionary; this is the contract).
3. Split train / val / test FIRST.
4. Build a sklearn Pipeline (or pandera-validated
   function) containing:
     a. Type coercion (to_numeric, to_datetime, astype)
     b. Encoding detection if reading raw files
     c. Deduplication (or pre-pipeline if cross-source)
     d. Missingness analysis → imputer + indicator
     e. Outlier treatment (decide remove / winsorize /
        transform per column)
     f. Scaling / power transform (numeric)
     g. Categorical encoding (with CV for target enc)
     h. Feature engineering (binning, interactions, dates)
     i. Imbalanced resampling (training fold only;
        imblearn.Pipeline)
5. Fit pipeline on train; transform val and test.
6. Validate cleaned schema (pandera / Great Expectations).
7. Persist the fitted pipeline alongside the model
   (joblib / pickle / ONNX).
8. Document assumptions (missingness mechanism,
   outlier policy, encoding choices).
```

## Practical Patterns

### Pattern 1: A leak-proof preprocessing pipeline

```python
from sklearn.compose import ColumnTransformer
from sklearn.pipeline import Pipeline
from sklearn.impute import SimpleImputer, KNNImputer
from sklearn.preprocessing import (
    StandardScaler, OneHotEncoder, TargetEncoder, RobustScaler
)
from sklearn.model_selection import train_test_split

X_train, X_test, y_train, y_test = train_test_split(
    X, y, test_size=0.2, stratify=y, random_state=42)

numeric = Pipeline([
    ('impute', KNNImputer(n_neighbors=5)),
    ('scale',  RobustScaler()),
])

nominal_low = Pipeline([
    ('impute', SimpleImputer(strategy='constant', fill_value='__MISSING__')),
    ('encode', OneHotEncoder(handle_unknown='ignore')),
])

nominal_high = Pipeline([
    ('impute', SimpleImputer(strategy='constant', fill_value='__MISSING__')),
    ('encode', TargetEncoder(smooth='auto', cv=5)),  # sklearn ≥ 1.3
])

preprocess = ColumnTransformer([
    ('num',  numeric,      numeric_cols),
    ('low',  nominal_low,  low_card_cols),
    ('high', nominal_high, high_card_cols),
])

pipe = Pipeline([
    ('prep',  preprocess),
    ('model', SomeEstimator()),
])
pipe.fit(X_train, y_train)
print(pipe.score(X_test, y_test))
```

All learning steps (imputers, scalers, encoders) are fitted **only on
training data** inside the cross-validation fold. No leakage.

### Pattern 2: Imbalanced + scaling, properly composed

```python
from imblearn.pipeline import Pipeline as ImbPipeline
from imblearn.over_sampling import SMOTE
from sklearn.preprocessing import StandardScaler

pipe = ImbPipeline([
    ('scale',  StandardScaler()),
    ('smote',  SMOTE(random_state=42)),  # train-fold only by construction
    ('model',  LogisticRegression(class_weight=None)),  # or 'balanced' w/o SMOTE
])
```

`imblearn.Pipeline` (not `sklearn.Pipeline`) is required so that
`SMOTE.fit_resample` is called only on training folds during CV.

### Pattern 3: Schema validation with Pandera

```python
import pandera as pa
from pandera.typing import Series, DataFrame

class CustomerSchema(pa.DataFrameModel):
    customer_id: Series[int]   = pa.Field(unique=True, gt=0)
    email:       Series[str]   = pa.Field(str_matches=r'^[^@]+@[^@]+$')
    age:         Series[int]   = pa.Field(ge=0, le=120, nullable=True)
    plan:        Series[str]   = pa.Field(isin=['free', 'pro', 'enterprise'])
    signup_at:   Series[pa.DateTime]

@pa.check_types
def transform_customers(df: DataFrame[CustomerSchema]) -> DataFrame[CustomerSchema]:
    ...  # transformations preserve the contract
```

### Pattern 4: Probabilistic dedupe with Splink

```python
import splink.duckdb.duckdb_comparison_library as cl
from splink.duckdb.linker import DuckDBLinker

settings = {
    "link_type": "dedupe_only",
    "blocking_rules_to_generate_predictions": [
        "l.postcode = r.postcode",
        "substr(l.surname,1,1) = substr(r.surname,1,1) and l.dob = r.dob",
    ],
    "comparisons": [
        cl.exact_match("email"),
        cl.jaro_winkler_at_thresholds("first_name", [0.9, 0.7]),
        cl.jaro_winkler_at_thresholds("surname", [0.9, 0.7]),
        cl.levenshtein_at_thresholds("postcode", [1, 2]),
    ],
}

linker = DuckDBLinker(df, settings)
linker.estimate_u_using_random_sampling(max_pairs=1e6)
linker.estimate_parameters_using_expectation_maximisation(
    "l.email = r.email")
preds = linker.predict()
clusters = linker.cluster_pairwise_predictions_at_threshold(preds, 0.95)
```

### Pattern 5: Datetime parsing for "mixed garbage"

```python
# Wrong: silent per-row format change
dates = pd.to_datetime(s)

# Right: explicit format or 'mixed'
dates = pd.to_datetime(s, format='mixed', errors='coerce')
n_bad = dates.isna().sum() - s.isna().sum()
if n_bad:
    log.warning("dropped %d unparseable dates", n_bad)
```

## Troubleshooting

| Symptom | Likely cause | Fix |
|---|---|---|
| Test score wildly worse than CV score | Preprocessing leakage | Wrap all data-learning steps in a `Pipeline`; refit on train fold only |
| Target encoder gives perfect train accuracy, terrible holdout | Encoded with full data including row's own y | Use sklearn ≥ 1.3 `TargetEncoder` (built-in CV) or K-fold encoding manually |
| Imputed dataset has near-zero correlations | Mean / median imputation collapsed variance | Switch to MICE / `IterativeImputer`; add `MissingIndicator` |
| Model bias against a demographic after imputation | MAR/MNAR assumption violation per subgroup | Stratify imputation; include subgroup in conditional imputer |
| `pd.to_datetime` returns object dtype not datetime64 | Mixed formats; pandas fell back to object | Pass explicit `format=`, use `format='mixed'` (pandas ≥ 2.0), or per-pattern parsing |
| `UnicodeDecodeError` on CSV read | Wrong encoding declared | Use `charset-normalizer` / `chardet 7` to detect; `pd.read_csv(..., encoding=detected)` |
| Two timestamps that should compare equal don't | One is naive, one tz-aware; or different tz | Localize naive → convert to UTC for storage / comparison |
| Strings `"café"` and `"café"` not equal | NFC vs NFD codepoint difference | `unicodedata.normalize('NFC', s)` before comparison |
| Fuzzy match misses obvious typo | Wrong metric (Levenshtein for transpositions; Jaro-Winkler for prefix typos) | Try `RapidFuzz`'s Jaro-Winkler; for phonetic matches add Metaphone |
| Splink memory blow-up | Blocking rule too loose → quadratic candidate pairs | Add more restrictive blocking rules; union multiple narrow rules |
| Holdout score on minority class drops after SMOTE | SMOTE applied to full data including test | Use `imblearn.Pipeline`; SMOTE only inside training fold |
| Class_weight='balanced' worse than SMOTE+model | Severe imbalance + complex boundary | Try SMOTE-Tomek / ADASYN; consider focal loss for NNs |
| Pandera schema fails on production batch but not on dev | Distribution drift; new category value | Add `Field(isin=...)` checks; treat schema violations as monitoring signal |
| StandardScaler-transformed feature still extreme | Heavy-tailed feature; outliers dominate mean/std | Switch to RobustScaler or QuantileTransformer; consider log1p first |
| Box-Cox fails with `ValueError: Data must be positive` | Box-Cox is positive-only | Use Yeo-Johnson (handles zero / negative) |
| Tree model performance unchanged by scaling | Expected; trees are scale-invariant | Skip scaling for tree models; save compute |

## Anti-Patterns (Quick Reference)

1. Preprocessing before train/test split.
2. Target encoding without K-fold CV.
3. `df.dropna()` with no `subset`.
4. Mean imputation with no missingness indicator (then trusting CIs).
5. Removing outliers without diagnosis.
6. SMOTE / resampling on the full dataset.
7. Mixing pytz, dateutil, and zoneinfo timezones without conversion.
8. `pd.to_datetime` on mixed-format columns without `format=`.
9. `errors='ignore'` silently keeping non-numeric strings.
10. `drop_duplicates()` without `subset` on a customer table.
11. PowerTransformer Box-Cox on data with zeros / negatives.
12. Ignoring the cyclical nature of hour-of-day with linear encoding.
13. One-hot encoding a 10,000-category feature in a linear model.
14. Validating only on dev data, never on the next batch.

## References

### Missing data
- van Buuren, S. *Flexible Imputation of Missing Data* (2nd ed.) —
  [§1.2 MCAR/MAR/MNAR](https://stefvanbuuren.name/fimd/sec-MCAR.html)
- Rubin, D. B. (1976). "Inference and missing data." *Biometrika* 63.
- [Mechanisms of Missing Data (MCAR, MAR, MNAR) — APXML](https://apxml.com/courses/intro-feature-engineering/chapter-2-handling-missing-data/missing-data-mechanisms)
- [Missing data mechanisms: MCAR, MAR, MNAR (UCL Discovery)](https://discovery.ucl.ac.uk/id/eprint/10150989/2/Pham_Missing%20data,%20part%202.%20Missing%20data%20mechanisms_AAM.pdf)
- [scikit-learn — 7.4 Imputation of missing values](https://scikit-learn.org/stable/modules/impute.html)
- [Comparing Multiple Imputation Methods (MICE vs miceforest vs IterativeImputer) — PMC](https://www.ncbi.nlm.nih.gov/pmc/articles/PMC12380239/)
- [Iterative Imputation for Missing Values — MachineLearningMastery](https://machinelearningmastery.com/iterative-imputation-for-missing-values-in-machine-learning/)
- [Missing Indicators in Feature Engineering — APXML](https://apxml.com/courses/intro-feature-engineering/chapter-2-handling-missing-data/missing-value-indicators)
- [On Missingness Features in Machine Learning for Critical Care — PMC](https://pmc.ncbi.nlm.nih.gov/articles/PMC8701717/)
- [Assumptions and analysis planning in studies with missing data — PMC](https://pmc.ncbi.nlm.nih.gov/articles/PMC10396404/)

### Outliers
- [Outliers Detection Using IQR, Z-score, LOF and DBSCAN — Analytics Vidhya](https://www.analyticsvidhya.com/blog/2022/10/outliers-detection-using-iqr-z-score-lof-and-dbscan/)
- [Outlier Detection Methods — IQR, Z-Score & Statistical Tests](https://statsolvepro.com/outlier-detection-methods/)
- [A Practical Guide to Outlier Detection and Treatment](https://medium.com/@inkollusrivarsha0287/a-practical-guide-to-outlier-detection-and-treatment-in-data-science-43e853614fe8)
- [Don't Throw Away Your Outliers — Towards Data Science](https://towardsdatascience.com/dont-throw-away-your-outliers-c37e1ab0ce19/)
- [How to Make Your ML Models Robust to Outliers — KDnuggets](https://www.kdnuggets.com/2018/08/make-machine-learning-models-robust-outliers.html)

### Type coercion, dates, encodings, timezones
- [pandas.to_numeric documentation](https://pandas.pydata.org/docs/reference/api/pandas.to_numeric.html)
- [pandas.to_datetime documentation](https://pandas.pydata.org/docs/reference/api/pandas.to_datetime.html)
- [pandas.Series.dt.tz_localize](https://pandas.pydata.org/docs/reference/api/pandas.Series.dt.tz_localize.html)
- [pandas Time Zone Handling](https://tedboy.github.io/pandas/timeseries/timeseries14.html)
- [Avoid Inferring Dates In Pandas](https://medium.com/@dvdpedros/avoid-this-shortcut-when-parsing-dates-in-pandas-a175bd39ec73)
- [Mastering Pandas dtypes — codepointtech](https://codepointtech.com/mastering-pandas-dtypes-the-complete-guide/)
- [Charset Detection in Python: chardet, cchardet, and charset-normalizer](https://bytetunnels.com/posts/charset-detection-python-chardet-cchardet-charset-normalizer/)
- [charset_normalizer documentation](https://charset-normalizer.readthedocs.io/en/latest/)
- [chardet (Python character encoding detector)](https://github.com/chardet/chardet)

### Deduplication / record linkage
- [Fuzzy Matching 101: 2026 Guide — MatchDataPro](https://matchdatapro.com/fuzzy-matching-101-a-complete-guide-for-2025/)
- [Fuzzy Matching 101 — Data Ladder (2026)](https://dataladder.com/fuzzy-matching-101/)
- [SQL Fuzzy Match Guide — Interview Query (2025)](https://www.interviewquery.com/p/sql-fuzzy-matching)
- [Deduplicating 7M records in 2 minutes with Splink — Robin Linacre](https://medium.com/data-science-collective/deduplicating-7-million-records-in-two-minutes-with-splink-4b1a87035a85)
- [An Interactive Introduction to Record Linkage (Fellegi-Sunter) — Robin Linacre](https://www.robinlinacre.com/intro_to_probabilistic_linkage/)
- [The Fellegi-Sunter Model — Splink documentation](https://moj-analytical-services.github.io/splink/topic_guides/theory/fellegi_sunter.html)
- [Splink: MoJ's open source library — GOV.UK](https://www.gov.uk/government/publications/joined-up-data-in-government-the-future-of-data-linking-methods/splink-mojs-open-source-library-for-probabilistic-record-linkage-at-scale)
- [Splink on PyPI](https://pypi.org/project/splink/)

### Scaling and transformation
- [scikit-learn 7.3 Preprocessing data](https://scikit-learn.org/stable/modules/preprocessing.html)
- [Compare the effect of different scalers on data with outliers — scikit-learn](https://scikit-learn.org/stable/auto_examples/preprocessing/plot_all_scaling.html)
- [PowerTransformer — scikit-learn](https://scikit-learn.org/stable/modules/generated/sklearn.preprocessing.PowerTransformer.html)

### Categorical encoding
- [Encoding Categorical Variables: One-Hot vs Label & Beyond — Unidata](https://unidata.pro/blog/encoding-categorical-variables-one-hot-vs-label/)
- [Categorical Encoding Practical Guide — Let's Data Science](https://letsdatascience.com/blog/categorical-encoding-a-practical-guide-to-one-hot-label-and-target-methods)
- [Ordinal and One-Hot Encodings for Categorical Data — MachineLearningMastery](https://machinelearningmastery.com/one-hot-encoding-for-categorical-data/)
- [How to Do Target Encoding Without Data Leakage — Medium](https://medium.com/@prathik.codes/how-to-do-target-encoding-without-data-leakage-the-right-way-280bd24fbc81)
- [Target Encoder: A powerful categorical encoding method — Train in Data](https://www.blog.trainindata.com/target-encoder-a-powerful-categorical-encoding-method/)
- [Sampling Techniques in Bayesian Target Encoding — arXiv 2006.01317](https://arxiv.org/pdf/2006.01317)

### Schema validation
- [Pandera Python Data Validation Guide 2026](https://pythondatabench.com/article/data-validation-python-pandera-practical-guide)
- [Data validation in Python: Pandera and Great Expectations — endjin](https://endjin.com/blog/a-look-into-pandera-and-great-expectations-for-data-validation)
- [The data validation landscape in 2025 — aeturrell](https://aeturrell.com/blog/posts/the-data-validation-landscape-in-2025/)
- [5 Python Data Validation Libraries — KDnuggets](https://www.kdnuggets.com/5-python-data-validation-libraries-you-should-be-using)

### Feature engineering
- [Feature Engineering 2025: Transforming Raw Data — madrigan](https://blog.madrigan.com/en/blog/202511291136/)
- [10 Best Techniques for Feature Engineering — TheTechThinker](https://thetechthinker.com/feature-engineering-in-machine-learning/)
- [Interaction Terms and Polynomial Features — APXML](https://apxml.com/courses/applied-data-science/chapter-2-practical-feature-engineering/interaction-polynomial-features)
- [Working with numerical data — Google ML Crash Course](https://developers.google.com/machine-learning/crash-course/numerical-data)

### Text cleaning
- [Text Preprocessing: Complete Guide to Tokenization, Normalization & Cleaning](https://mbrenndoerfer.com/writing/text-preprocessing-nlp-tokenization-normalization)
- [Text Normalization: Unicode Forms, Case Folding & Whitespace](https://mbrenndoerfer.com/writing/text-normalization-unicode-nlp)
- [Text Cleaning, Normalization and Representation — NLP for Social Science](https://nlp4ss.jeju.ai/en/session02/lecture1.html)
- [Text Preprocessing in NLP — GeeksforGeeks](https://www.geeksforgeeks.org/nlp/text-preprocessing-for-nlp-tasks/)

### Imbalanced data
- [SMOTE for Imbalanced Classification — MachineLearningMastery](https://machinelearningmastery.com/smote-oversampling-for-imbalanced-classification/)
- [Handling Imbalanced Datasets with SMOTE, ADASYN & Class Weighing — Coderspacket](https://coderspacket.com/posts/handling-imbalanced-datasets-with-smote-adasyn-class-weighing/)
- [Imbalanced Learning: SMOTE, Class Weighting And Beyond — ML4Devs](https://www.ml4devs.com/what-is/imbalanced-data-model-training/)
- [Comparative Analysis of Resampling Techniques (XGBoost financial distress) — MDPI Mathematics 2025](https://www.mdpi.com/2227-7390/13/13/2186)
- [Optimal Data Augmentation Ratio for ADASYN — arXiv 2510.18252](https://arxiv.org/pdf/2510.18252)

### Anti-patterns and pipelines
- [scikit-learn — Common pitfalls and recommended practices](https://scikit-learn.org/stable/common_pitfalls.html)
- [How to Avoid Data Leakage When Performing Data Preparation — MachineLearningMastery](https://machinelearningmastery.com/data-preparation-without-data-leakage/)
- [Prevent Data Leakage in Cross-Validation — Medium](https://medium.com/@pacosun/split-smart-avoiding-data-leakage-in-cross-validation-894863a93553)
- [Mastering Pandas DataFrame.dropna() — Bomberbot](https://www.bomberbot.com/python/mastering-pandas-dataframe-dropna-a-comprehensive-guide-for-data-cleaning-and-preprocessing/)

### Modern data quality ecosystem
- [15 Best Data Quality Tools for 2026 — Mammoth](https://mammoth.io/blog/data-quality-tools/)
- [Best Data Quality Tools for 2026: Selection Guide — Atlan](https://atlan.com/know/data-quality/top-tools/)
- [Monte Carlo vs Bigeye vs Soda in 2026 — The SaaS Podium](https://thesaaspodium.com/monte-carlo-bigeye-soda-comparison/)
- [The 2026 Open-Source Data Quality and Observability Landscape — DataKitchen](https://datakitchen.io/blog/the-2026-open-source-data-quality-and-data-observability-landscape/)

## Related skills

- **da-3-data-acquisition-sampling** — what to collect and how to sample
  it; runs before this skill.
- **da-5-exploratory-data-analysis** — profiling, visualization,
  summary statistics; runs after this skill.
- **da-1-2-measurement-theory** — the levels-of-measurement question
  ("can I take a mean of this?") underlies every encoding and scaling
  choice here.
- **da-1-3-probability-theory** — the distributional assumptions behind
  z-score, Box-Cox, MICE.
- **da-2-data-analysis-lifecycle** — where this stage sits in the
  CRISP-DM / KDD / OSEMN / TDSP frameworks.
- **mongodb-schema-design** — when the upstream store is MongoDB.
- **mongodb-bson-types** — BSON-specific type coercion at the driver
  layer.
- **xlsx** — when the upstream is Excel and tabular cleaning is the goal.
