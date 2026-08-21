<!-- hub-reference-banner -->
> **Reference file — part of the `da-2-data-analysis-lifecycle` hub.** Formerly the standalone `da-2-4-documentation-reproducibility` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-2-4-documentation-reproducibility
version: "1.1.0"
updated: "2026-05-30"
description: |
  Documentation and reproducibility practices as they apply to the data analysis lifecycle: computational notebooks (Jupyter, R Markdown, Quarto), data and code versioning, data lineage and provenance tracking, environment management, and workflow documentation.

  TRIGGER: Use when the conversation involves making a data analysis reproducible; choosing or structuring computational notebooks; documenting analytical decisions, transformations, or data flow; version-controlling datasets, code, or environments; tracking data lineage or provenance through a pipeline; setting up DVC, MLflow, or similar tools for analysis reproducibility; reviewing a notebook for hidden-state or cell-ordering issues; building a README or workflow document for an analysis project; parameterizing notebooks for repeatable runs; testing notebooks for output consistency.

  SKIP: Defer to `da-1-6-3-reproducibility-replicability` for philosophical distinctions between computational reproducibility vs. scientific replicability. Defer to general software documentation skills for API docs, changelog writing, or project documentation not tied to data analysis workflows. Defer to `da-13-data-engineering-and-pipelines` for pipeline orchestration and engineering patterns beyond documentation concerns.

related_skills:
  - da-1-6-3-reproducibility-replicability
  - da-13-data-engineering-and-pipelines
  - da-2-1-process-frameworks
  - da-4-data-cleaning-preparation
---

# Documentation & Reproducibility in the Data Analysis Lifecycle

## What This Concept Covers

Within the data analysis lifecycle, documentation and reproducibility refers to the practices that let an analyst—or someone else—re-run an analysis on the same data and get identical outputs, and re-read the work and understand every decision that was made. The concept has three interlocking parts: **computational notebooks** (how you write and share the analysis), **versioning** (tracking changes to code, data, and environments over time), and **lineage/provenance** (recording where data came from and how it was transformed).

These are not academic niceties. Studies find that fewer than 26% of shared Jupyter notebooks can re-execute all cells successfully in shared environments, and fewer than 5% produce identical results to stored outputs when run from scratch ([Pimentel et al., PMC 2021](https://pmc.ncbi.nlm.nih.gov/articles/PMC8106381/)). The failures trace back to concrete, fixable practices.

---

## 1. Computational Notebooks: Literate Programming for Analysis

### Tools Overview

A computational notebook combines code, its outputs, and human-readable narrative in one document. The concept originates in Knuth's "literate programming" (1984) — the idea that a program should explain its logic to a human reader, not just instruct a machine. Applied to data analysis, the notebook is the primary artifact: it shows what was done, why, and what it found.

Three tools dominate:

- **Jupyter Notebooks** — language-agnostic via kernels (Python, R, Julia, many more); widely used for exploratory analysis, ML, and scientific computing ([PLOS Computational Biology, 2019](https://journals.plos.org/ploscompbiol/article?id=10.1371%2Fjournal.pcbi.1007007))
- **R Markdown** — R-centric, tight RStudio integration, mature ecosystem of output formats (HTML, PDF, Word, slides)
- **Quarto** — the successor to R Markdown; natively supports Python, R, Julia, and JavaScript; produces the same output formats; designed around reproducibility by default (each render starts from a clean slate) ([UBC MDS course materials](https://ubc-mds.github.io/DSCI_521_platforms-dsci_book/lectures/6-rmarkdown-quarto-slides-ghpages.html))

Quarto's clean-slate render model is the clearest structural advantage over Jupyter: because the document always re-executes from top to bottom in a fresh environment, the hidden-state problem described below cannot arise by construction.

### The Hidden-State Problem

The single most common reproducibility failure in Jupyter notebooks is **out-of-order or skipped cell execution**. A user who runs cells in an ad-hoc order and then saves the notebook leaves the kernel state inconsistent with the visual narrative. When someone else opens the file and runs all cells top-to-bottom, they get different results or outright errors.

Measured at scale:
- 76.88% of shared notebooks contain execution-counter skips (cells re-run or deleted mid-session without resetting) ([Pimentel et al., 2021](https://pmc.ncbi.nlm.nih.gov/articles/PMC8106381/))
- 36.45% show evidence of out-of-order execution
- NameError exceptions increase 1.27–1.39x in notebooks with skips, pointing to undefined variables left in kernel memory

The fix is procedural: **always use "Restart Kernel and Run All Cells" before sharing or committing a notebook**. This exposes errors that an incremental run would hide.

### Notebook Authoring Rules (from published guidelines)

The following rules are synthesized from two peer-reviewed treatments ([PLOS Computational Biology 2019](https://journals.plos.org/ploscompbiol/article?id=10.1371%2Fjournal.pcbi.1007007); [Briefings in Bioinformatics 2023](https://academic.oup.com/bib/article/24/6/bbad375/7326135)):

1. **Tell a story, not a log.** The notebook should have a beginning (question/context), middle (analysis steps with rationale), and end (conclusions). Code alone cannot convey why a particular transformation was chosen or what a result means.

2. **Document process, not just results.** Record explorations including dead ends. Include analysis dates and who ran it. Add explanatory text continuously, not after the fact.

3. **One meaningful step per cell.** Cells exceeding ~100 lines should be split. Use markdown headers to delineate sections. Each cell should perform one coherent operation.

4. **Modularize reused code.** Wrap repeated logic in functions; consider extracting frequently reused code to a separate module. Avoid copy-pasting cells.

5. **Record all dependencies.** Use `pip freeze > requirements.txt` or `conda env export > environment.yml`. Print package versions inside the notebook using tools like `watermark`. This is the most commonly neglected step: only 12.5% of repositories in one large study declared any dependencies ([Pimentel et al., 2021](https://pmc.ncbi.nlm.nih.gov/articles/PMC8106381/)).

6. **Use version control.** Commit notebooks to Git. Use `nbdime` (or Quarto/R Markdown's text-based formats) to get meaningful diffs — standard `git diff` on Jupyter `.ipynb` JSON is unreadable.

7. **Build toward a pipeline.** Design exploratory notebooks so they can become parameterized pipelines. Put variable declarations at the top. Before every commit, restart the kernel and run all cells top to bottom; if anything breaks, the notebook was not in a reproducible state.

8. **Share data alongside code.** For small datasets, commit them to the repo. For larger ones, use Zenodo, Figshare, or OSF with a DOI. Document upstream steps when only a processed version can be shared.

9. **Make notebooks runnable by others.** Use Binder or Docker for zero-install cloud environments. Store static HTML/PDF renderings alongside the live notebook so reviewers can read without running.

### Parameterized Notebooks

For analyses that must run repeatedly against different inputs (date ranges, cohorts, parameter sweeps), **Papermill** injects parameters into a Jupyter notebook at runtime and records them in output metadata, making each run independently traceable. Usage:

```bash
papermill analysis.ipynb output/run_2026-05.ipynb -p start_date 2026-05-01 -p cohort "enterprise"
```

Each output notebook is a self-contained record of one run. Quarto and R Markdown achieve the same effect through parameterized YAML headers (`params:` block), making them equally suitable for scheduled or scripted re-execution.

### Notebook Testing

A notebook that passes "Restart & Run All" manually is still not tested against future changes. Two tools address this:

- **nbval** — runs a notebook via pytest and asserts that each cell's output matches stored outputs; flags cells whose output changes between runs
- **pytest-notebook** — similar approach with richer assertion options and integration with CI pipelines

Add notebook tests to CI so that changes to upstream data, helper functions, or library versions surface as test failures rather than silent result drift.

---

## 2. Versioning: Code, Data, and Environments

### Why Three Separate Versioning Concerns

Code versioning with Git is well understood. Data analysis adds two additional axes that Git alone does not handle well:

- **Data versioning** — raw data files change (new extracts, corrections, schema changes). Git is not designed for large binaries or fast-changing tabular data.
- **Environment versioning** — an analysis that worked in `pandas==1.3` may silently change behavior in `pandas==2.0`. Capturing and pinning the compute environment is as important as capturing the code.

### Code Versioning

Standard Git practices apply: commit frequently with descriptive messages, use branches for exploratory work, write a README that explains the project structure and how to reproduce the analysis. Pair with Git tagging at publication or handoff points so that specific versions can be retrieved unambiguously ([KDnuggets best practices](https://www.kdnuggets.com/best-practices-for-version-control-in-data-science-projects)).

Notebook-specific note: Jupyter `.ipynb` files are JSON and produce large, noisy diffs when outputs are stored. Configure `nbstripout` as a pre-commit hook to strip outputs before committing, or use Quarto/R Markdown text-based sources which diff cleanly.

### Data Versioning

**DVC (Data Version Control)** is the standard open-source tool for versioning datasets and models alongside code. DVC stores lightweight pointer files (`.dvc`) in Git while actual data lives in configurable remote storage (S3, GCS, Azure Blob, SFTP, etc.). A `dvc repro` command re-runs the full pipeline with any stage's dependencies tracked.

Key DVC concepts:
- `dvc add <file>` — tracks a data file, creates a `.dvc` pointer committed to Git
- `dvc run` / `dvc.yaml` — defines pipeline stages with explicit inputs, outputs, and commands
- `dvc push/pull` — syncs data to/from remote storage
- Each Git commit + DVC state together define a fully reproducible snapshot

The integration between DVC and MLflow addresses the full lifecycle: DVC manages data and pipeline execution while MLflow logs experiment parameters, metrics, and artifacts ([LearnCodePro tutorial](https://www.learncodepro.com/tutorials/data-science-machine-learning-ai/model-deployment-mlops/data-versioning-dvc-experiment-tracking-mlflow)). Tying dataset versions to code commits is critical — without this link, it becomes impossible to know which data produced which result ([LabelYourData, 2026](https://labelyourdata.com/articles/machine-learning/data-versioning)).

**Alternative tools by use case:**

| Tool | Best for |
|---|---|
| **LakeFS** | Git-like branching and rollback for data lakes (S3, GCS, Azure) |
| **Delta Lake** | ACID transactions and time-travel queries on Spark tables |
| **DoltDB** | SQL-queryable, Git-committed relational data; diff and merge on table rows |

Use DVC when your data lives outside a data lake and your pipeline is script-based. Use LakeFS or Delta Lake when data lives in a cloud lake and you need branch-level isolation during development.

### Environment Versioning

Three primary options ([Briefings in Bioinformatics 2023](https://academic.oup.com/bib/article/24/6/bbad375/7326135)):

- **Conda environments** — `environment.yml` captures all package versions; broadly portable across Linux/macOS/Windows for data science workloads
- **Docker containers** — full OS-level capture; most portable for sharing with others; widely adopted due to community resources and pre-built images; minimal runtime overhead
- **Guix** — offers reproducible builds with cryptographic guarantees; less common outside bioinformatics

The minimum viable approach is a pinned `requirements.txt` or `environment.yml`. The gold standard for long-term reproducibility is a Docker image whose build instructions are committed alongside the analysis code.

---

## 3. Data Lineage and Provenance

### Definitions and Distinction

**Data provenance** focuses on origins: where did this data come from, and who collected or generated it?

**Data lineage** describes the full lifecycle: source → transformations → destination, including every intermediate step ([Atlan lineage guide](https://atlan.com/know/data-lineage-tracking/)).

The two terms are often conflated, but the distinction matters in practice. Provenance answers "what is the raw material and can I trust it?" while lineage answers "how did it become what I'm analyzing, and what else depends on it?"

### Why Lineage Matters in Analysis

When an upstream dataset is corrected or a transformation bug is found, lineage identifies every downstream analysis affected — often within minutes rather than days of manual auditing. Without lineage records, finding all affected outputs requires tracing every notebook and report by hand, and some will be missed.

Lineage also supports compliance: regulators under GDPR, HIPAA, and similar frameworks increasingly expect auditable proof of data provenance ([IBM data lineage overview](https://www.ibm.com/think/topics/data-lineage)).

### Lineage Capture Methods

1. **Manual documentation** — README files, data dictionaries, transformation logs. Fragile and quickly becomes stale.
2. **Pipeline-native lineage** — tools like **dbt** and **Apache Airflow** produce lineage graphs as a byproduct of defining the pipeline. This is the most reliable approach for structured ETL/ELT work.
3. **Query parsing** — automated tools analyze SQL to infer joins and transformations (column-level lineage).
4. **Log-based tracking** — extracts lineage from database write-ahead logs, CDC streams.
5. **API-driven capture** — pulls metadata from cloud platforms (Snowflake, Databricks, BigQuery).

### The OpenLineage Standard

**OpenLineage** (hosted by LF AI & Data Foundation) is the open standard for pipeline-level lineage metadata. It defines a JSON event format emitted by running jobs (Airflow tasks, Spark jobs, dbt models) that any compatible backend can consume. Tools like Apache Airflow, Spark, dbt, and Flink all have native OpenLineage integrations ([OpenLineage GitHub](https://github.com/OpenLineage/OpenLineage)). The resulting graph provides end-to-end provenance across heterogeneous platforms without requiring a single vendor's stack.

For more formal provenance modeling, the **W3C PROV** family of specifications provides a domain-agnostic data model built around three core concepts: entities (data artifacts), activities (processes that create or transform them), and agents (people or systems responsible). PROV underpins many academic and scientific provenance systems ([W3C PROV paper, ACM 2013](https://dl.acm.org/doi/10.1145/2452376.2452478)).

### Column-Level vs. Table-Level Lineage

Table-level lineage shows which datasets connect to which. Column-level lineage follows individual fields through derivations and transformations — essential for regulated data (to demonstrate that a PII field was correctly masked end-to-end) and for diagnosing data quality issues in specific calculated fields.

---

## 4. Workflow Documentation: What to Record Beyond the Notebook

A well-authored notebook captures the analytical narrative. The surrounding project infrastructure needs its own documentation:

| Document | Purpose | Minimum content |
|---|---|---|
| **README** | Entry point for anyone new to the project | Project question; how to reproduce from raw data to outputs; directory structure; links to data sources and dependencies |
| **Data dictionary / codebook** | Define every variable | Name, type, units, allowed values, and source field for each variable |
| **Decision log** | Record analytical choices and rejected alternatives | Why a threshold was chosen, why a model was selected, what alternatives were considered and discarded |
| **Change history** | Track methodology and data updates over time | When data was refreshed, when methodology changed, when outputs were last validated |

The decision log is the document most often omitted and most often regretted. It lives outside the notebook narrative because it captures reasoning that doesn't fit the linear story: "we considered approach X but chose Y because Z" belongs in the log, not in a code comment.

---

## 5. Common Pitfalls

| Pitfall | Consequence | Fix |
|---|---|---|
| Running notebook cells out of order without restarting | Results differ between author's machine and reviewer's | Always use "Restart & Run All" before committing |
| Missing or unpinned dependencies | Analysis fails or produces different results months later | Pin all dependencies; export environment on every significant change |
| Hardcoded absolute paths | `FileNotFoundError` on any other machine | Use relative paths anchored to project root; use an `.env` file for data root |
| No data version recorded | Can't identify which snapshot of the data produced which result | Tag data versions in DVC; store hash/checksum in notebook header |
| Transformation logic in notebook only, not documented | Next analyst can't tell why a filter or recoding was applied | Write narrative justification in markdown cells; maintain decision log |
| Treating lineage documentation as optional | Downstream breakage not discovered until stakeholders report wrong numbers | Adopt pipeline-native lineage from the start (dbt, Airflow); don't try to backfill |
| Storing secrets in notebooks | Credentials leak when notebooks are shared or published | Use environment variables or a secrets manager; never hardcode tokens or API keys in cells |

---

## 6. Quick Reference by Situation

**Starting a new analysis project:**
- Set up a Git repo with a clear directory structure (`data/raw/`, `data/processed/`, `notebooks/`, `src/`)
- Create `environment.yml` or `requirements.txt` on day one
- Write a README stub; fill it in as you work, not at the end
- If data is large or shared, initialize DVC

**Choosing a notebook format:**
- Quarto if you want guaranteed-clean rendering and multi-language support
- R Markdown if the analysis is R-only and the team is already RStudio-native
- Jupyter if you need an interactive, kernel-rich environment or the team uses Python exclusively

**Auditing an existing notebook for reproducibility:**
- Check execution counters — do they run 1, 2, 3, ... or are there gaps and jumps?
- Run "Restart & Run All" — does it complete without errors?
- Check for hardcoded paths and missing imports
- Verify `requirements.txt` or `environment.yml` exists and is current

**Documenting data transformations:**
- Use dbt models with `description:` fields for SQL-based transforms
- Use DVC `dvc.yaml` stage definitions to record pipeline steps and dependencies
- Log key transformation parameters as MLflow run metadata if the analysis involves model training

---

## Key Sources

1. [Ten Simple Rules for Writing and Sharing Computational Analyses in Jupyter Notebooks](https://journals.plos.org/ploscompbiol/article?id=10.1371%2Fjournal.pcbi.1007007) — Rule et al., PLOS Computational Biology, 2019. The canonical reference for notebook best practices.

2. [Understanding and Improving the Quality and Reproducibility of Jupyter Notebooks](https://pmc.ncbi.nlm.nih.gov/articles/PMC8106381/) — Pimentel et al., IEEE/ACM MSR, 2019 (PMC 2021). Empirical study of 1.4M notebooks; source of the <26% / <5% reproducibility statistics.

3. [Five Pillars of Computational Reproducibility](https://academic.oup.com/bib/article/24/6/bbad375/7326135) — Briefings in Bioinformatics, 2023. Framework covering literate programming, version control, environment management, data sharing, and documentation.

4. [Data Lineage Tracking: Complete Guide](https://atlan.com/know/data-lineage-tracking/) — Atlan, 2026. Practical guide to lineage types, capture methods, and implementation patterns.

5. [OpenLineage (GitHub)](https://github.com/OpenLineage/OpenLineage) — LF AI & Data Foundation. The open standard for pipeline-level lineage metadata across Airflow, Spark, dbt, and other tools.

6. [Reproducible and Trustworthy Workflows for Data Science](https://ubc-dsci.github.io/reproducible-and-trustworthy-workflows-for-data-science/) — UBC MDS course. Practical coverage of Quarto, Git, DVC, and environment management in teaching context.

7. [W3C PROV Family of Specifications](https://dl.acm.org/doi/10.1145/2452376.2452478) — Moreau et al., ACM EDBT, 2013. The foundational standard for provenance metadata modeling.
