<!-- hub-reference-banner -->
> **Reference file — part of the `da-analytical-methods` hub.** Formerly the standalone `da-5-exploratory-data-analysis` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-5-exploratory-data-analysis
title: Exploratory Data Analysis (EDA)
version: "1.0.0"
updated: "2026-05-30"
category: data-analysis
origin: local
description: >
  Exploratory Data Analysis — the Tukey-rooted discipline of looking at data
  before modeling it. Covers univariate / bivariate / multivariate analysis,
  summary statistics, distribution checking, visual EDA, automated EDA tools
  (ydata-profiling, Sweetviz, Autoviz, D-Tale), time-series EDA, and the
  honest-EDA discipline (pre-registration, garden-of-forking-paths, HARKing).
  TRIGGER: user starts a new dataset analysis, asks "what does this data look like",
  needs distribution / outlier / correlation patterns surfaced, wants a profiling
  report, is choosing variables for modeling, debugging a model that surprises
  them in production, or asks about Tukey / EDA / data profiling.
  SKIP: explicit confirmatory hypothesis testing with a pre-registered protocol
  (use da-1-4-statistical-inference-foundations); fitting a specific model
  (da-6 / da-7); cleaning the data (da-4); building production dashboards
  (da-8-data-visualization); pure data cleaning (da-4-data-cleaning-preparation).
triggers:
  - exploratory data analysis
  - EDA
  - data profiling
  - data exploration
  - univariate analysis
  - bivariate analysis
  - multivariate analysis
  - summary statistics
  - distribution check
  - look at the data
  - ydata-profiling
  - sweetviz
  - autoviz
  - Tukey EDA
keywords:
  - EDA
  - data-profiling
  - univariate
  - bivariate
  - multivariate
  - QQ-plot
  - histogram
  - boxplot
  - scatter-plot
  - correlation
  - Pearson
  - Spearman
  - PCA
  - t-SNE
  - UMAP
  - ydata-profiling
  - autoviz
  - sweetviz
  - HARKing
  - garden-of-forking-paths
when_to_use:
  - You just received a new dataset and need to understand its shape
  - You want to surface outliers, missing patterns, or surprising distributions
  - You're choosing features for a model and need correlation / collinearity insight
  - A model behaves unexpectedly and you need to inspect input data
  - You want a one-shot profiling report (ydata-profiling, Sweetviz, Autoviz)
  - You're doing time-series EDA (STL decomposition, ACF/PACF, change points)
when_not_to_use:
  - You're cleaning data — use da-4-data-cleaning-preparation
  - You're running a confirmatory hypothesis test — use da-1-4-statistical-inference-foundations
  - You're fitting a regression / classification / clustering model — use da-6 / da-7
  - You're producing a publication chart — use da-8-data-visualization
related_skills:
  - da-1-3-probability-theory
  - da-1-4-statistical-inference-foundations
  - da-4-data-cleaning-preparation
  - da-6-statistical-modeling
  - da-8-data-visualization
---

# Exploratory Data Analysis

EDA is the discipline of looking at data *before* you model it. The phrase is John Tukey's (1977, *Exploratory Data Analysis*), and his framing still holds: EDA is **hypothesis-generating**, not **hypothesis-confirming**. The output of EDA is questions, not p-values.

## When to use this skill

Activate when the user:

- starts a new dataset and asks "what does this data look like"
- needs distribution / outlier / correlation patterns surfaced
- wants a profiling report (ydata-profiling, Sweetviz, Autoviz, D-Tale)
- is choosing variables for downstream modeling
- is debugging a model that surprises them in production
- asks about Tukey, EDA discipline, or honest-EDA practice

## When NOT to use this skill

- Confirmatory hypothesis testing with a pre-registered protocol → `da-1-4-statistical-inference-foundations`
- Fitting a specific model → `da-6-statistical-modeling` or `da-7-machine-learning`
- Cleaning the data → `da-4-data-cleaning-preparation`
- Building production dashboards or publication charts → `da-8-data-visualization`

---

## Tukey's EDA philosophy (1977)

Tukey separated **confirmatory** data analysis (CDA) from **exploratory** data analysis (EDA). The distinction matters because the two have opposite epistemics:

| Aspect | EDA | CDA |
|---|---|---|
| Goal | Find patterns | Test patterns |
| Hypothesis | Generated from the data | Pre-specified |
| p-values | Misleading (multiple looks) | Valid |
| Errors | Type-1 risk if treated as CDA | The point |
| Reporting | "Suggests further investigation" | "Confirms / rejects H0" |

The two failure modes when EDA pretends to be CDA: **HARKing** (Hypothesizing After Results Known) and the **garden-of-forking-paths**. If you tested every variable until something hit p < 0.05, that p-value is not a p-value — it's the cost of doing many silent tests.

The safe pattern: **separate the EDA dataset from the test dataset** (or pre-register the hypothesis before looking).

---

## Univariate analysis

Look at each variable alone. The goal: distribution shape, central tendency, spread, anomalies.

### Continuous variables

| Tool | What it shows | When to use |
|---|---|---|
| **Histogram** | Distribution shape, modes | First pass, always |
| **KDE** (kernel density estimate) | Smoothed distribution | When bin choice matters |
| **Boxplot** | Median, IQR, outliers | Quick outlier check |
| **Violin plot** | Distribution + boxplot in one | Comparing groups |
| **QQ plot** | Departure from a reference distribution | Check normality assumption |

### Summary statistics (default panel)

- **Center**: mean, median, mode (and *why mean ≠ median* when they diverge — skew)
- **Spread**: standard deviation, IQR, range, MAD (median absolute deviation, robust)
- **Shape**: skewness, kurtosis (be careful — kurtosis of normal is 3 in some conventions, 0 in others; "excess kurtosis" subtracts 3)
- **Tails**: percentiles (p1, p5, p95, p99) — much more useful than min/max for tail behavior
- **Missingness**: count and % of NA/null/empty (links to `da-4-data-cleaning-preparation`)

### Categorical variables

- **Frequency table** with counts + proportions
- **Bar chart** (not pie, especially not for > 3 categories — see `da-8-data-visualization`)
- **High-cardinality flag**: if a categorical has > 50 unique values, treat as either an identifier (don't model directly) or apply target encoding / frequency encoding before modeling

---

## Bivariate analysis

Two variables at a time. Goal: relationships, correlations, conditional patterns.

### Continuous × Continuous

- **Scatter plot** — the workhorse. Add alpha or use hexbin for > 5k points.
- **Correlation coefficients**:
  - **Pearson** — linear association; assumes both variables are roughly normal; sensitive to outliers
  - **Spearman** — rank-based; monotonic association; robust to outliers
  - **Kendall's τ** — also rank-based; smaller for the same data than Spearman; preferred in small samples or with many ties
- **Anscombe's quartet** is the canonical reminder: four datasets, all with identical mean, variance, and Pearson r — but radically different shapes. **Plot the data.**

### Categorical × Continuous

- **Boxplot grouped by category** (or violin)
- **Strip plot / swarm plot** for small samples
- **Group means + confidence intervals**

### Categorical × Categorical

- **Cross-tab** with counts and row/column percentages
- **Mosaic plot** — area encodes joint frequency
- **Chi-squared test of independence** (technically CDA, but acceptable as an EDA pointer)

---

## Multivariate analysis

Three or more variables at once.

- **Pair plot** (seaborn `pairplot`, ggplot2 `GGally::ggpairs`) — all pairwise scatters + univariate diagonals in a grid. Good up to ~6 variables; beyond that becomes unreadable.
- **Correlation heatmap** — `df.corr()` rendered as a heatmap. Use a diverging palette centered at zero. Cluster the rows/columns for blockwise structure.
- **Parallel coordinates plot** — one polyline per row across normalized axes; reveals multivariate clusters.
- **PCA for EDA** — first 2-3 principal components plotted as scatter; color by class label or cluster. Different goal than PCA-for-modeling: here you want to see structure, not reduce dimensions.
- **t-SNE / UMAP** — nonlinear dimensionality reduction. UMAP is faster, preserves global structure better, and is the modern default for visual cluster inspection. Caution: t-SNE / UMAP distances are **not** meaningful; only neighborhood structure is.
- **Conditional plots** (faceted small multiples) — split scatters by a third categorical variable.

---

## Distribution checking

Before modeling, check whether your variables match the distributional assumptions of the chosen model.

- **QQ plot** — points on the line → matches reference distribution. Heavy tails curve away at the ends.
- **Shapiro-Wilk** — formal normality test for n < 5000. Very sensitive — will reject "approximately normal" at large n. Use as a pointer, not a gate.
- **Anderson-Darling** — more sensitive to tails than Shapiro-Wilk.
- **Kolmogorov-Smirnov** — general distribution-match test (any reference distribution).

Rule of thumb: don't use a formal normality test to decide whether to use a parametric method. **Plot the QQ plot.** The eye is better than the test, and the central limit theorem usually saves you anyway above n=30 unless distributions are pathological.

---

## Time-series EDA

Time-series data has structure regular EDA misses.

- **Line plot** with time on x-axis (always start here)
- **STL decomposition** (Seasonal-Trend decomposition using Loess) — separates the series into trend + seasonality + residual
- **ACF / PACF** plots — autocorrelation and partial autocorrelation. Tells you what AR/MA order to try if modeling with ARIMA
- **Change-point detection** — `ruptures` (Python), `bcp` (R), or visual inspection
- **Stationarity tests** — ADF (Augmented Dickey-Fuller), KPSS. Same caveat as Shapiro-Wilk: pointer, not gate.
- **Rolling statistics** — rolling mean / rolling std with a windowed view. If both wander, the series is non-stationary.
- **Seasonal subseries plot** — one line per season showing values across years
- **Calendar heatmap** — date × hour (or day) heatmap for daily/hourly patterns

---

## Automated EDA tools

These produce a profiling report in one line of code. Use them for the first-pass overview; do NOT use them as a substitute for thinking.

| Tool | Strength | Watch out for |
|---|---|---|
| **ydata-profiling** (formerly pandas-profiling) | Most thorough; correlations, warnings, missing-value patterns | Slow on > 100k rows; report can be 5+ MB HTML |
| **Sweetviz** | Beautiful HTML; train-vs-test comparison built-in | Less comprehensive than ydata |
| **Autoviz** | Auto-selects chart types per variable | Sometimes picks wrong chart type |
| **D-Tale** | Interactive in-browser pandas dataframe explorer | Live exploration, not for handoff reports |
| **`df.describe()`** (pandas / Polars) | Built-in summary stats | First thing you should always do |
| **`skim` (R) or `skimr`** | One-line richer summary than base R | R-only |

A reasonable default workflow: `df.head()` → `df.info()` → `df.describe()` → ydata-profiling on a sample if the dataset is large.

---

## Practical EDA recipes

### Recipe: new tabular dataset, 5-minute first pass

```python
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt

df = pd.read_csv("data.csv")

# 1. Shape and types
print(df.shape, df.dtypes.value_counts())

# 2. Missingness pattern
print(df.isna().mean().sort_values(ascending=False).head(10))

# 3. Summary
print(df.describe(include='all').T)

# 4. Univariate distributions for numerics
df.select_dtypes("number").hist(bins=30, figsize=(15, 10))
plt.tight_layout(); plt.show()

# 5. Correlation heatmap
sns.heatmap(df.corr(numeric_only=True), cmap='RdBu_r', center=0, annot=True, fmt='.2f')
plt.show()

# 6. Pairplot if ≤ 6 numeric columns
if df.select_dtypes("number").shape[1] <= 6:
    sns.pairplot(df.select_dtypes("number"))
    plt.show()
```

### Recipe: time-series 5-minute first pass

```python
import statsmodels.tsa.seasonal as smt

# 1. Plot the series
df['value'].plot(figsize=(15, 4))

# 2. STL decomposition
stl = smt.STL(df['value'], period=24).fit()
stl.plot()

# 3. ACF/PACF
from statsmodels.graphics.tsaplots import plot_acf, plot_pacf
plot_acf(df['value'], lags=48); plot_pacf(df['value'], lags=48)

# 4. Stationarity check
from statsmodels.tsa.stattools import adfuller
print(adfuller(df['value'].dropna()))
```

---

## Anti-patterns

1. **Reporting EDA findings as confirmed effects.** "We found a strong correlation between X and Y" without test/holdout. The p-value is meaningless after EDA.
2. **HARKing** — Hypothesizing After Results Known. You looked, you found, you wrote the paper as though you'd predicted it. This is fraud in the formal sense.
3. **Garden of forking paths** — every analytic decision (which outliers to drop, which transform to apply, which model family) is a branching point. With enough branches, something will look "significant." Pre-register the analysis path.
4. **Treating EDA visualizations as final deliverables.** EDA charts are for *you*. Production charts need design (use `da-8-data-visualization`).
5. **Ignoring missingness patterns.** A 90%-missing column is information itself. Drop it, but understand why.
6. **Correlation matrices on mixed-scale data.** Pearson on a binary variable is technically defined but rarely what you want; use point-biserial or just stratify.
7. **Trusting automated profiling without verification.** Tools like ydata-profiling will silently truncate to a sample without telling you in some configurations.
8. **Time-series EDA with shuffled rows.** Always sort by time first.

---

## EDA discipline (honest practice)

Two practices separate honest EDA from p-hacking:

1. **Pre-registration**. Before looking, write down the planned analysis. Deviations get reported as exploratory. The OSF (Open Science Framework) is one common registry. See also `da-1-6-3-reproducibility-replicability`.
2. **Train/test split before EDA**. Do EDA on the training set; reserve the test set untouched. Pretend the test set does not exist until you have a final model and a final report.

When neither of those is realistic, at least **flag** in the writeup that the analysis is exploratory. The right phrase is "this analysis suggests further investigation"; the wrong phrase is "this analysis confirms".

---

## References

1. Tukey, J. W. (1977). *Exploratory Data Analysis*. Addison-Wesley. The original.
2. Wickham, H., & Grolemund, G. (2017). *R for Data Science*, Chapter 7 — EDA.
3. McKinney, W. (2022). *Python for Data Analysis* (3rd ed.), Chapter 10.
4. Tufte, E. R. (2001). *The Visual Display of Quantitative Information*. Graphics Press. Background on visual EDA.
5. Gelman, A., & Loken, E. (2014). "The garden of forking paths." *American Statistician*. The forking-paths paper.
6. Kerr, N. L. (1998). "HARKing: Hypothesizing After the Results are Known." *Personality and Social Psychology Review*.
7. ydata-profiling: [https://github.com/ydataai/ydata-profiling](https://github.com/ydataai/ydata-profiling)
8. Sweetviz: [https://github.com/fbdesignpro/sweetviz](https://github.com/fbdesignpro/sweetviz)
9. UMAP: McInnes, L. et al. (2018). *UMAP: Uniform Manifold Approximation and Projection*. arXiv:1802.03426.
10. Anscombe, F. J. (1973). "Graphs in statistical analysis." *American Statistician* — the four datasets.
