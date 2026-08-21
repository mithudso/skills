<!-- hub-reference-banner -->
> **Reference file — part of the `da-3-data-acquisition-sampling` hub.** Formerly the standalone `da-3-data-collection-acquisition` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-3-data-collection-acquisition
description: |
  Overview and orientation skill for Data Collection & Acquisition — the third stage of the Data
  Analysis lifecycle (Data Analysis > Data Collection & Acquisition). Covers what data collection is,
  why it sits between problem framing and cleaning, the full landscape of collection approaches
  (primary vs secondary, qualitative vs quantitative, automated vs manual), the planning process,
  data quality dimensions, provenance documentation at collection time, and the decision logic for
  choosing a collection strategy. This is a parent/overview node: it maps the territory and directs
  you to specialist child nodes.

  TRIGGER: Use when the user asks "how do I collect data for my analysis", "what data collection
  method should I use", "what is primary vs secondary data", "what is the difference between
  surveys, experiments, APIs, and scraping at a high level (overview)", "how do I plan a data
  collection strategy", "what should I document when collecting data", "what are data quality
  dimensions at collection time", "how do I document data provenance", or any
  introductory/orientation question about acquiring data for analysis. Also trigger at the start of
  CRISP-DM "Data Understanding", OSEMN "Obtain", KDD "Selection", or TDSP "Data Acquisition" stages
  when the user needs to choose an approach before diving into specifics.

  SKIP: Once the user has chosen a specific collection mechanism and needs implementation depth,
  defer to specialist skills. API ingestion, web scraping implementation, CDC, streaming, ETL/ELT,
  survey design, sampling frame construction, sample size calculation, data contracts, file format
  selection, and validation tooling → da-3-data-acquisition-sampling. Data cleaning, missing values,
  deduplication, normalization → da-4-data-cleaning-preparation. Privacy, consent, ethics,
  GDPR/HIPAA compliance in data collection → da-11-ethics-and-privacy. Pure sampling theory and
  distributions → da-1-4-2-sampling-distributions-standard-error or da-1-3-probability-theory.
  Process framework selection (CRISP-DM vs OSEMN vs KDD) → da-2-1-process-frameworks.
---

# Data Collection & Acquisition: Overview

Data collection is "the systematic process of gathering, measuring and analyzing information from
various sources" to support decision-making, hypothesis testing, and model building
([GeeksforGeeks — Methods of Data Collection](https://www.geeksforgeeks.org/data-analysis/methods-of-data-collection/)).
It sits between problem framing (defining what question you are answering) and data cleaning
(making the collected data usable), and the choices made here determine what analytical questions
remain answerable downstream.

**Output from this skill:** a recommended collection approach with rationale, and a pointer to the
relevant child skill for implementation depth.

---

## 1. Why Collection Strategy Comes First

The collection phase answers: *What data exists or can be gathered, from whom or what, using which
mechanism, and under what constraints?*

Poor collection decisions propagate forward. Biased sampling produces biased estimates even with
perfect cleaning. An API that returns only a 30-day window cannot support year-over-year trend
analysis. A survey that omits a demographic group cannot generalize to that group. Fixing these
problems after collection is often impossible
([IBM — Data Quality Issues and Challenges](https://www.ibm.com/think/insights/data-quality-issues)).

Before selecting a method, establish three things:

1. **Research objective** — What specific question must the data answer? Objectives drive variable
   selection and unit of analysis ([Sapien — How to Create a Data Collection Plan](https://www.sapien.io/blog/how-to-create-a-data-collection-plan)).
2. **Target population** — Who or what are observations drawn from? The collection method must
   reach that population; if it cannot, you have coverage bias. Example: a web survey reaches
   internet users, not the general public — this gap matters if your question is about all adults.
3. **Constraints** — Time, budget, legal/ethical boundaries, data volume, and required freshness
   all constrain which methods are viable
   ([SurveyCTO — Create a Data Collection Plan in 7 Steps](https://www.surveycto.com/data-collection-quality/data-collection-plan-seven-steps/)).

---

## 2. Primary vs Secondary Data

The most fundamental split in data collection:

| Dimension | Primary Data | Secondary Data |
|---|---|---|
| Origin | Collected by you, for this study | Pre-existing; collected by someone else |
| Fit | Tailored to your exact question | May not match your construct or population |
| Cost/time | High | Low |
| Quality control | You own it | You inherit their decisions |
| Typical examples | Surveys, experiments, interviews, direct sensor capture | Published reports, public datasets, existing databases, scraped web data |

Primary data is preferred when the question is novel, when no existing dataset covers the
right population or time period, or when you need to control measurement instruments. Secondary
data is preferred when cost or time matters and an adequate existing source exists
([QuestionPro — Data Collection Methods](https://www.questionpro.com/blog/data-collection-methods/),
[GeeksforGeeks](https://www.geeksforgeeks.org/data-analysis/methods-of-data-collection/)).

**Combining both:** many analyses use primary data for the core construct and secondary data for
covariates or benchmarks. When combining, reconcile any differences in population coverage,
collection period, and variable definitions before merging — undocumented differences between
primary and secondary sources are a frequent source of analysis errors.

---

## 3. The Collection Method Landscape

### 3.1 Human-Generated Primary Data

| Method | Best For | Key Limitation |
|---|---|---|
| **Surveys / questionnaires** | Attitudes, self-reported behaviors, demographics at scale | Response bias, recall error, low response rates |
| **Interviews** | Motivations, mental models, complex processes | Does not scale; observer bias risk |
| **Focus groups** | Product feedback, concept testing, group reactions | Results not generalizable to broader populations |
| **Observation** | Behavioral data without self-report; field research | Labor-intensive; subjects may alter behavior when observed |
| **Experiments** | Establishing causal relationships | Expensive; artificial conditions; ethical constraints on manipulation |

([QuestionPro](https://www.questionpro.com/blog/data-collection-methods/),
[GeeksforGeeks](https://www.geeksforgeeks.org/data-analysis/methods-of-data-collection/),
[Airbyte](https://airbyte.com/data-engineering-resources/data-collection-techniques))

### 3.2 System- and Machine-Generated Data

| Method | Best For | Key Limitation |
|---|---|---|
| **API integration** | Real-time feeds, social media, financial data, any service with a defined contract | Rate limits, auth management, provider availability |
| **Web scraping** | Publicly available data when no API exists; price monitoring, content aggregation | Terms of service / legal risk; brittle to layout changes |
| **Database extraction / CDC** | Operational data from relational or document databases; Change Data Capture (CDC) for row-level change events | Requires schema knowledge and appropriate permissions |
| **IoT / sensor data** | Predictive maintenance, environmental monitoring, health tracking; continuous high-frequency streams | Requires edge or stream processing infrastructure (MQTT, Kafka, Flink) |
| **Log and event data** | User behavior analysis from application logs and event streams | Requires parsing and normalization before analysis |
| **Transaction data** | Purchase behavior, financial records — structured and high-fidelity | Captures completed transactions only; not intent or context |

([Pingax — Data Collection Methods: APIs, Web Scraping, and More](https://pingax.com/data-analytics/stage-2-data-collection/data-sources-and-techniques/data-collection-methods-apis-web-scraping-and-more/),
[Airbyte](https://airbyte.com/data-engineering-resources/data-collection-techniques))

### 3.3 Secondary / Existing Data Sources

| Source Type | Examples | Watch For |
|---|---|---|
| Published datasets | Government census, IMF/World Bank, UCI ML Repository, Harvard Dataverse | Recency lag; granularity gaps |
| Internal records | CRM, support tickets, financial ledgers, HR data | Access controls; privacy constraints |
| Research literature data | Data from published studies (for meta-analysis) | Validity depends on original study quality |

([GeeksforGeeks](https://www.geeksforgeeks.org/data-analysis/methods-of-data-collection/))

---

## 4. Qualitative vs Quantitative Collection

The qualitative/quantitative distinction cuts across the methods above and affects downstream
analysis choices:

| Dimension | Quantitative | Qualitative |
|---|---|---|
| Goal | Measure, count, generalize | Understand, interpret, explore |
| Sample strategy | Probability sampling for representativeness | Purposive/theoretical sampling for information richness |
| Sample size | Large (power analysis driven) | Small (data saturation driven) |
| Data form | Numbers, ratings, counts | Text, narratives, themes |
| Common methods | Structured surveys, experiments, automated logs | Interviews, focus groups, observation, document analysis |

Mixed-methods designs use both — primary quantitative collection for breadth, qualitative for
depth — and are common in applied research
([PMC — Sampling in Qualitative Research](https://www.ncbi.nlm.nih.gov/pmc/articles/PMC5774281/)).

---

## 5. Choosing a Collection Method: Decision Logic

Work through these questions in order:

1. **Does the required data already exist?** If a public or internal dataset covers your
   population, time range, and variables adequately, use secondary data first. Validate its
   quality before investing in primary collection.

2. **Does the question require causal inference?** If yes and a randomized controlled trial (RCT)
   or controlled experiment is feasible and ethical, prioritize experimental design. Otherwise,
   use observational data with appropriate causal methods (see `da-12-ab-testing-causal-inference`).

3. **Is the target phenomenon attitudinal, behavioral, or system-generated?**
   - Attitudinal or opinion (what people think/believe) → survey or interview
   - Behavioral at a human level (what people do, captured by a person or observer) → observation, interview
   - Behavioral at scale or system-generated (clickstreams, logs, sensor telemetry) → automated collection: APIs, logs, sensors, CDC
   - Both attitudinal and behavioral → mixed methods

4. **What are the freshness and volume requirements?** Real-time or near-real-time needs push
   toward streaming APIs, CDC, or IoT pipelines. Historical batch analysis can use bulk exports.

5. **What are the legal and ethical constraints?** Data involving personal information requires
   informed consent, data minimization, and compliance with applicable law. Web scraping requires
   review of terms of service. See `da-11-ethics-and-privacy` for depth.

6. **What is the budget?** Surveys and experiments carry direct cost. API-based collection may
   have per-call pricing. Primary collection always costs more than reusing secondary data.

**Worked example:** A researcher wants to understand why customers churn. The question is attitudinal
("why do they leave?") and behavioral ("when do they leave?"). No existing internal dataset captures
churn reasons (only churn events). Decision: primary qualitative collection (exit interviews) for
the "why", combined with secondary internal data (CRM records) for the "when and who". Next step:
`da-3-data-acquisition-sampling` for survey design and CDC setup.

---

## 6. Planning a Data Collection Effort

A data collection plan should document
([SurveyCTO](https://www.surveycto.com/data-collection-quality/data-collection-plan-seven-steps/),
[Sapien](https://www.sapien.io/blog/how-to-create-a-data-collection-plan)):

- **Objective and research questions** — What decision or hypothesis does this data support?
- **Variables and operationalization** — For each construct, how will it be measured? (See
  `da-1-2-4-operationalization-constructs`.)
- **Target population and sampling frame** — Who/what are observations drawn from, and how will
  the sample be selected?
- **Collection instrument or mechanism** — Survey instrument, API endpoint, scraper specification,
  sensor configuration.
- **Collection timeline and responsible parties**
- **Data storage and access controls**
- **Metadata and provenance documentation** (see Section 7 below)

---

## 7. Data Quality and Provenance at Collection Time

Errors introduced during collection are harder to correct than errors introduced later, because
the ground truth is gone. Key quality dimensions to control at collection:

| Dimension | What to check | Risk if neglected |
|---|---|---|
| **Accuracy** | Does the instrument measure what it claims? (Survey wording, sensor calibration, schema mapping) | Systematic bias that cleaning cannot remove |
| **Completeness** | Are all required records and fields captured? | Missing data that can only be imputed, not recovered |
| **Consistency** | Same definitions, units, and formats across sources and time periods? | Integration errors when combining sources |
| **Timeliness** | Data captured at the right time relative to the phenomenon? (Recall bias in surveys; stale API snapshots) | Measurement error tied to temporal mismatch |
| **Representativeness** | Does the sample reflect the target population? (Sampling bias, coverage bias, response bias) | Estimates that do not generalize |

([IBM — Data Quality Issues and Challenges](https://www.ibm.com/think/insights/data-quality-issues),
[IEEE DataPort — How to Ensure Data Quality in Research Datasets](https://ieee-dataport.org/news/how-ensure-data-quality-research-datasets))

### Provenance documentation

Data provenance records the origin and history of data: who collected it, when, from what source,
using what instrument or protocol, and what transformations have been applied. Provenance
documentation enables auditability, reproduction of results, and assessment of fitness for a
different use case
([Secoda — Best Practices for Documenting Data Provenance](https://www.secoda.co/blog/best-practices-for-documenting-data-provenance)).

Minimum provenance fields at collection time:
- Source system or instrument name and version
- Collection date and time (with timezone)
- Collection method (API endpoint, survey ID, sensor ID)
- Population and sampling approach
- Known limitations or caveats

Treat metadata as infrastructure, not documentation: capture it automatically where possible and
store it alongside the data from the start
([Secoda](https://www.secoda.co/blog/best-practices-for-documenting-data-provenance)).

---

## 8. Common Pitfalls

| Pitfall | Consequence | Mitigation |
|---|---|---|
| Collecting data before defining the question | Data does not answer the actual research need | Define question and variables first |
| Convenience sampling instead of probability sampling | Unrepresentative results; cannot generalize | Define sampling frame; use random selection |
| No pilot or pre-test of the instrument | Systematic measurement errors go undetected until it is too late | Pilot with a small subset before full collection |
| Missing provenance documentation | Cannot reproduce or audit results | Document source, method, date at collection time |
| Ignoring response/non-response bias | Systematic skew in estimates | Track non-response; compare respondents to population profile |
| Collecting more variables than needed | Privacy risk; processing overhead; collection burden on subjects | Collect only what the question requires (data minimization) |
| No validation at ingest | Bad records enter the pipeline silently | Apply schema checks and range assertions at the point of entry |

---

## 9. Where This Sits in the Data Analysis Lifecycle

```
Problem framing          → da-2-2-problem-framing
  ↓
Data Collection &        ← YOU ARE HERE (da-3-data-collection-acquisition)
Acquisition (overview)
  ↓
  ├─ API / scraping / CDC / streaming / ETL / ELT / surveys / sampling design
  │   → da-3-data-acquisition-sampling (implementation depth)
  │
  └─ Ethics, consent, GDPR, privacy → da-11-ethics-and-privacy
  ↓
Data Cleaning &          → da-4-data-cleaning-preparation
Preparation
  ↓
Exploratory Analysis     → da-5-exploratory-data-analysis
```

---

## Sources

1. [GeeksforGeeks — Methods of Data Collection](https://www.geeksforgeeks.org/data-analysis/methods-of-data-collection/) — Primary vs secondary; interview, questionnaire, observation, experiment, focus group methods
2. [QuestionPro — Data Collection Methods](https://www.questionpro.com/blog/data-collection-methods/) — Overview of primary/secondary split and quantitative/qualitative types
3. [Airbyte — Data Collection Techniques](https://airbyte.com/data-engineering-resources/data-collection-techniques) — 10 modern techniques including IoT, API, CDC, AI text analysis, transaction tracking
4. [Pingax — Data Collection Methods: APIs, Web Scraping, and More](https://pingax.com/data-analytics/stage-2-data-collection/data-sources-and-techniques/data-collection-methods-apis-web-scraping-and-more/) — API, scraping, and automated method details
5. [IBM — Data Quality Issues and Challenges](https://www.ibm.com/think/insights/data-quality-issues) — Quality dimensions, bias types, collection-time error patterns
6. [IEEE DataPort — How to Ensure Data Quality in Research Datasets](https://ieee-dataport.org/news/how-ensure-data-quality-research-datasets) — Standardized protocols, validation, metadata, security
7. [Secoda — Best Practices for Documenting Data Provenance](https://www.secoda.co/blog/best-practices-for-documenting-data-provenance) — Provenance documentation, lineage, metadata as infrastructure
8. [SurveyCTO — Create a Data Collection Plan in 7 Steps](https://www.surveycto.com/data-collection-quality/data-collection-plan-seven-steps/) — Planning process, alignment with objectives, field protocols
9. [Sapien — How to Create a Data Collection Plan](https://www.sapien.io/blog/how-to-create-a-data-collection-plan) — Objective alignment, population definition, method selection
10. [PMC — Series: Practical guidance to qualitative research, Part 3](https://www.ncbi.nlm.nih.gov/pmc/articles/PMC5774281/) — Qualitative sampling and data collection; saturation; iterative process
