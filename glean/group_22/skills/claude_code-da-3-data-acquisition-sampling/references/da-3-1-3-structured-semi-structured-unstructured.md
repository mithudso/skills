<!-- hub-reference-banner -->
> **Reference file — part of the `da-3-data-acquisition-sampling` hub.** Formerly the standalone `da-3-1-3-structured-semi-structured-unstructured` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-3-1-3-structured-semi-structured-unstructured
description: |
  Expert knowledge on structured, semi-structured, and unstructured data as data source types in the context of data collection and acquisition. Covers definitions, the structuredness spectrum, source examples, collection approaches, storage fit, and analyst decision guidance including hybrid/ambiguous sources.

  TRIGGER:
  - User asks what type of data a given source is (structured, semi-structured, or unstructured)
  - User is evaluating or selecting data sources for a data collection or acquisition step
  - User asks how to acquire, ingest, or preprocess data of each type
  - User asks about the differences between structured, semi-structured, and unstructured data
  - User asks which storage system (data warehouse vs. data lake) fits their data sources
  - User asks about schema-on-write vs. schema-on-read in the context of source data types
  - User asks about the proportion or volume of each type in enterprise data
  - Questions about data taxonomy at the source level in a data analysis pipeline
  - User asks how to classify a data source that seems to span multiple categories

  SKIP:
  - Questions focused on database engines, RDBMS internals, or NoSQL systems as technology choices — defer to database-type skills (ai-datastores, mongodb-schema-design, etc.)
  - Questions about specific file formats or serialization protocols in depth (e.g., Parquet vs. ORC, Avro, Protocol Buffers) — defer to format-specific skills
  - Questions about ETL/ELT pipelines, data integration tools, or orchestration — defer to data engineering/pipeline skills (da-13-data-engineering-and-pipelines)
  - Questions about statistical analysis methods or modeling — defer to da-6-statistical-modeling or related skills
  - Questions about data quality rules, validation pipelines, or data contracts — defer to da-4-data-cleaning-preparation
  - Questions about web scraping mechanics, crawling tools, or HTML extraction techniques — defer to da-3-2-5-web-scraping-crawling (HTML as a data format is noted here for classification purposes only)

related_skills:
  - da-3-1-data-sources
  - da-3-data-collection-acquisition
  - da-3-2-collection-methods
  - da-3-1-1-primary-vs-secondary
  - da-3-1-2-internal-vs-external
  - da-13-data-engineering-and-pipelines
  - da-4-data-cleaning-preparation
  - da-3-2-5-web-scraping-crawling
---

# Structured, Semi-Structured, and Unstructured Data (Data Sources)

## Output Guidance

When this skill is active, structure responses as follows:
- **Classification question** ("Is X structured?"): give a one-sentence verdict, the key reason, and the recommended storage target.
- **Comparison question** ("What's the difference between..."): use the §8 Summary Comparison table as the spine, then add the 1–2 most relevant rows for the user's context.
- **Acquisition/ingestion question** ("How do I collect X?"): lead with the data type classification, then give the collection method and storage fit from the relevant section.
- **Hybrid/ambiguous source**: apply the §5 decision rule (dominant processing burden), state the classification, and note which part of the source drives the cost.

## Context in the Data Analysis Taxonomy

This concept sits at **Data Analysis > Data Collection & Acquisition > Data sources**. The classification answers a foundational question every analyst must settle before designing an acquisition pipeline: *what form does the data arrive in, and what does that imply for collection, storage, and preparation?*

---

## 1. The Structuredness Spectrum

Data source types are not three discrete buckets — they occupy a **continuum of structuredness**, measured by how readily a machine can parse and interpret the semantics without human intervention. ([ResearchGate, continuum diagram, 2014](https://www.researchgate.net/figure/The-continuum-of-structuredness-in-terms-of-human-and-machine-understandable-information_fig1_266269866))

```
Unstructured ←——————————————————→ Structured
raw video/audio   HTML/XML/JSON   relational tables
free text         log files       sensor CSV feeds
images            emails          OLAP cubes
```

The three named categories are practical shorthand for regions on this spectrum, not hard boundaries.

---

## 2. Structured Data

### Definition
Structured data is organized within a **predefined schema**: fixed columns with declared data types, rows representing records. It obeys the tabular model associated with relational databases. ([Databricks, 2024](https://www.databricks.com/blog/structured-vs-unstructured-data); [Splunk, 2024](https://www.splunk.com/en_us/blog/learn/data-structured-vs-unstructured-vs-semi-structured.html))

### Characteristics
- Fixed, schema-on-write format — structure is enforced at ingestion time
- Immediately queryable with standard tools (SQL, BI platforms) without reformatting
- Quantitative fields dominate: numbers, dates, booleans, short categorical strings
- High consistency across records within the same source
- Low flexibility: adding a new field requires a schema migration

### Canonical Source Examples
| Source | What makes it structured |
|---|---|
| Relational database tables (CRM, ERP, POS) | Fixed column definitions, enforced types |
| Spreadsheets with typed columns | User-defined but consistent schema |
| IoT sensor readings logged in tabular form | Numeric values at fixed timestamps |
| Web server access logs (structured format) | Predefined field order per line |
| Financial transaction records | Standard ledger fields: amount, date, account |

([AWS, 2024](https://aws.amazon.com/compare/the-difference-between-structured-data-and-unstructured-data/); [Alation, 2024](https://www.alation.com/blog/structured-unstructured-semi-structured-data/))

### Collection and Ingestion
Direct extraction via SQL queries, JDBC/ODBC connectors, or API endpoints that return typed payloads. Schema validation can be applied at ingestion without transformation. ETL pipelines typically load directly into a data warehouse.

### Storage Fit
Data warehouses (Snowflake, BigQuery, Redshift, traditional RDBMS). OLAP cubes for multidimensional aggregation workloads.

### Analysis Readiness
High. Analysts and business users can query and visualize structured data directly with minimal preprocessing. ([Alation, 2024](https://www.alation.com/blog/structured-unstructured-semi-structured-data/))

---

## 3. Semi-Structured Data

### Definition
Semi-structured data **does not conform to the rigid tabular schema of relational databases** but carries **tags, markers, or metadata that separate and describe its semantic elements**, giving it a self-describing quality. ([Wikipedia, Semi-structured data](https://en.wikipedia.org/wiki/Semi-structured_data); [Altexsoft, 2024](https://www.altexsoft.com/blog/semi-structured-data/))

The term originates in database research (the Object Exchange Model, or OEM, predated XML as an early formalization). With the Internet's expansion, formats like XML and later JSON became the dominant semi-structured exchange mechanisms because organizations needed flexibility beyond fixed-schema tables. ([Wikipedia](https://en.wikipedia.org/wiki/Semi-structured_data))

### Characteristics
- **Self-describing**: data carries its own field names (e.g., `"name": "John Doe"` in JSON)
- **Flexible schema**: records from the same source may have different fields present or absent
- **Hierarchical / nested**: supports tree or graph structures not expressible in flat tables
- **Partial consistency**: some structure is predictable; other parts vary by record
- Schema-on-read is the natural processing model — structure is applied at query time
- Eliminates the object-relational impedance mismatch that plagues mapping nested objects to tables

([Wikipedia](https://en.wikipedia.org/wiki/Semi-structured_data); [Splunk, 2024](https://www.splunk.com/en_us/blog/learn/data-structured-vs-unstructured-vs-semi-structured.html); [Altexsoft, 2024](https://www.altexsoft.com/blog/semi-structured-data/))

### Common Formats
| Format | Primary acquisition context |
|---|---|
| **JSON** | REST APIs, web browsers, IoT device payloads, mobile app events |
| **XML** | Enterprise system integrations, document exchange, SOAP services |
| **YAML** | Configuration files, DevOps pipelines (Kubernetes, Docker Compose) |
| **HTML** | Web scraping, content extraction from pages |
| **Email (with headers)** | Message metadata structured; body unstructured |
| **Log files with key=value pairs** | Application and system logs |
| **NoSQL document records** | Event streams, content management, social data |

([Altexsoft, 2024](https://www.altexsoft.com/blog/semi-structured-data/))

### Collection and Ingestion
APIs returning JSON or XML are the most common acquisition path. Web scraping produces HTML. IoT devices often emit JSON over MQTT or HTTP. Ingestion requires a parser for the specific format, plus handling of optional/missing fields and varying nesting depths.

### Storage Fit
Document stores (MongoDB, Couchbase), column stores with schema-on-read (Apache Hive, AWS Athena on S3), cloud data platforms supporting native JSON/XML columns (Snowflake VARIANT, BigQuery JSON). Also fits data lakes as intermediate staging before normalization.

### Analysis Readiness
Medium. The inherent structure (keys, tags) allows parsing and field extraction without full manual annotation, but query complexity is higher than SQL against flat tables. Tooling gap: no query language as universal as SQL for semi-structured data. ([Wikipedia](https://en.wikipedia.org/wiki/Semi-structured_data))

### Pipeline Role
Semi-structured data also functions as a **transition format** between raw collection and analysis-ready storage:

1. **Source to lake**: Raw API and IoT payloads arrive as JSON/XML and land in a data lake in semi-structured form.
2. **Lake to warehouse**: ETL/ELT flattens nested fields into structured tables for BI consumption, or retains the flexible form for schema-on-read analytics.
3. **Unstructured enrichment**: Adding metadata tags to images, audio, or documents produces a hybrid — see §5 for classification guidance.

---

## 4. Unstructured Data

### Definition
Unstructured data has **no predefined organizational format** — it arrives in its native form, and meaning must be extracted rather than read from a schema. ([AWS, 2024](https://aws.amazon.com/compare/the-difference-between-structured-data-and-unstructured-data/); [Splunk, 2024](https://www.splunk.com/en_us/blog/learn/data-structured-vs-unstructured-vs-semi-structured.html))

### Characteristics
- No schema enforced at creation time
- Qualitative, contextual, or sensory in nature
- High information density but low machine-readability without preprocessing
- Requires domain-specific extraction techniques per data type (NLP for text, CV for images, ASR for audio)
- Volume creates the largest acquisition burden: estimates place unstructured data at **80–90% of all enterprise data growth** ([Splunk, 2024](https://www.splunk.com/en_us/blog/learn/data-structured-vs-unstructured-vs-semi-structured.html); [Alation, 2024](https://www.alation.com/blog/structured-unstructured-semi-structured-data/))

### Canonical Source Examples
| Source | Primary extraction method |
|---|---|
| Free-text documents (PDFs, Word files) | NLP, named entity recognition, embedding |
| Emails | NLP, header parsing (hybrid: headers are semi-structured) |
| Social media posts and comments | NLP, sentiment analysis |
| Audio recordings (calls, meetings) | Speech-to-text, then NLP |
| Video files | Computer vision, object detection, transcription |
| Images (photos, medical scans, satellite) | Computer vision, OCR |
| Open-ended survey responses | NLP, topic modeling |
| System logs (free-form messages) | Regex extraction, log parsing |

([Databricks, 2024](https://www.databricks.com/blog/structured-vs-unstructured-data); [AWS, 2024](https://aws.amazon.com/compare/the-difference-between-structured-data-and-unstructured-data/))

### Collection and Ingestion
Data must be gathered from file systems, message queues, APIs (social platforms, communication tools), or streaming pipelines. Raw form is preserved in a data lake. Preprocessing — tokenization, vectorization, OCR, speech-to-text — is required before analysis.

### Storage Fit
Data lakes (S3, GCS, ADLS), content management systems, digital asset management systems, vector databases (after embedding). Lakehouse architectures unify storage with governed access. ([Databricks, 2024](https://www.databricks.com/blog/structured-vs-unstructured-data))

### Analysis Readiness
Low without preprocessing. High processing cost; requires data scientists or ML engineers rather than general analysts. Insights are richer and contextual — captures sentiment, narrative, and qualitative nuance that structured sources miss. ([Alation, 2024](https://www.alation.com/blog/structured-unstructured-semi-structured-data/))

---

## 5. Hybrid and Ambiguous Sources

Many real-world sources span categories. Classify by the **dominant processing burden**, not the container format.

| Source | Why it spans categories | Classify as |
|---|---|---|
| Email | Headers are structured (To, From, Date, Subject); body is unstructured free text | Unstructured (body drives the analysis cost) |
| Server log with key=value pairs | Format is semi-structured; individual message fields may be free text | Semi-structured (parseable without NLP) |
| PDF with embedded tables | Document container is unstructured; embedded tables are structured | Unstructured (requires extraction step first) |
| Video with JSON sidecar metadata | Media content is unstructured; metadata is semi-structured | Semi-structured if metadata is the analysis target; Unstructured if content is |
| CSV exported from a NoSQL store | File format is flat/structured; field values may contain nested JSON strings | Semi-structured (nested values need a second-pass parse) |
| Social media post with structured engagement counts | Text body is unstructured; like/share counts are structured | Mixed — collect both; store separately |

**Decision rule for hybrids**: identify which part of the source answers your analysis question. Apply the classification of that part to drive storage, tooling, and preprocessing decisions. ([AWS, 2024](https://aws.amazon.com/compare/the-difference-between-structured-data-and-unstructured-data/))

---

## 6. Decision Framework for Data Acquisition

When evaluating a data source, ask:

| Question | Classification / Action |
|---|---|
| Does the source have a fixed, known schema? | Structured — target a data warehouse directly |
| Does the source self-describe its fields with tags or keys? | Semi-structured — parse and store with schema-on-read |
| Does the source require human interpretation to extract meaning? | Unstructured — plan a preprocessing pipeline |
| Do you need fast SQL analytics immediately? | Prioritize structured or flatten semi-structured first |
| Do you need qualitative or contextual insight? | Unstructured is unavoidable; budget for ML/NLP |
| Are schema changes frequent or unknown at collection time? | Prefer a semi-structured or unstructured landing zone |

([Alation, 2024](https://www.alation.com/blog/structured-unstructured-semi-structured-data/); [Databricks, 2024](https://www.databricks.com/blog/structured-vs-unstructured-data))

---

## 7. Common Pitfalls

**Misclassifying at acquisition time**
Logs and emails feel "structured enough" but have highly variable formats. Treating them as structured without a robust parser causes silent data loss when the format changes.

**Underestimating semi-structured complexity**
JSON looks simple, but deeply nested arrays and optional fields multiply parsing complexity. Schema drift (field names or types changing between API versions) is a common failure mode.

**The "data swamp" problem**
Dumping all unstructured and semi-structured data into a data lake without metadata cataloguing, lineage, and governance turns the lake unusable. ([Databricks, 2024](https://www.databricks.com/blog/structured-vs-unstructured-data))

**Ignoring unstructured sources**
Analysts who only collect structured data miss 80–90% of available enterprise data. Social media sentiment, call recordings, and clinical notes often contain the most predictive signals for business questions. ([Splunk, 2024](https://www.splunk.com/en_us/blog/learn/data-structured-vs-unstructured-vs-semi-structured.html))

**Assuming structure implies quality**
Structured data can still contain nulls, outliers, encoding errors, and stale records. The tabular form makes *finding* quality issues easier, not *preventing* them.

**One-size-fits-all storage**
Schema-on-write (data warehouse) and schema-on-read (data lake) serve different needs. Forcing unstructured sources into a warehouse schema prematurely destroys information. ([Splunk, 2024](https://www.splunk.com/en_us/blog/learn/data-structured-vs-unstructured-vs-semi-structured.html))

---

## 8. Summary Comparison

| Dimension | Structured | Semi-Structured | Unstructured |
|---|---|---|---|
| Schema | Fixed, predefined | Self-described via tags/keys | None |
| Consistency | High across records | Partial; optional fields common | Low |
| Machine readability | High | Medium | Low (requires preprocessing) |
| Typical formats | SQL tables, CSV, spreadsheets | JSON, XML, YAML, HTML, logs | Text, PDF, audio, video, images |
| Primary storage | Data warehouse | Document store / data lake | Data lake / vector DB |
| Query approach | SQL | Schema-on-read, document query | ML/NLP/CV post-extraction |
| Analyst access | Direct | Moderate tooling required | Data scientist required |
| Volume share in enterprise | ~10–20% | Growing with APIs/IoT | ~80–90% |
| Insight type | Quantitative, operational | Transactional + contextual | Qualitative, contextual, narrative |

---

## Sources

1. [Databricks — Structured vs. Unstructured Data](https://www.databricks.com/blog/structured-vs-unstructured-data) — Definitions, critical differences, lakehouse architecture guidance, pitfalls
2. [Splunk — Structured, Semi-Structured & Unstructured Data](https://www.splunk.com/en_us/blog/learn/data-structured-vs-unstructured-vs-semi-structured.html) — Definitions, examples, analyst impact, strategic guidance
3. [Alation — Structured, Unstructured, Semi-Structured Data](https://www.alation.com/blog/structured-unstructured-semi-structured-data/) — Definitions, collection and processing methods, decision framework
4. [Altexsoft — Semi-Structured Data](https://www.altexsoft.com/blog/semi-structured-data/) — Semi-structured formats, flexible schema characteristics, challenges
5. [Wikipedia — Semi-Structured Data](https://en.wikipedia.org/wiki/Semi-structured_data) — Formal definition, historical context, OEM/XML/JSON technical background
6. [AWS — Structured vs. Unstructured Data](https://aws.amazon.com/compare/the-difference-between-structured-data-and-unstructured-data/) — Definitions, storage options, semi-structured hybrids
7. [ResearchGate — Continuum of Structuredness](https://www.researchgate.net/figure/The-continuum-of-structuredness-in-terms-of-human-and-machine-understandable-information_fig1_266269866) — Spectrum model from academic research
