<!-- hub-reference-banner -->
> **Reference file — part of the `da-3-data-acquisition-sampling` hub.** Formerly the standalone `da-3-1-data-sources` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-3-1-data-sources
description: >
  Expert knowledge on data sources as a category within data collection &
  acquisition. Covers the full taxonomy: what a data source is, how sources are
  typed and inventoried across the primary/secondary/tertiary,
  internal/external, and structured/semi-structured/unstructured dimensions,
  how to build a source inventory, and how to assess a source's fitness for a
  given analytical purpose (relevance, completeness, freshness, accuracy, bias).

  TRIGGER: user asks about types of data sources, how to categorize or
  inventory sources, what counts as a data source, how to evaluate or choose
  data sources for a project, how to build a data source inventory, or how to
  assess source fitness. Also trigger when the question spans multiple
  classification dimensions at once (e.g., "what kinds of internal structured
  data do I have?") or when discussing the front-end of any data collection or
  acquisition plan where source selection is the central question.

  SKIP: when the question is *only* about the primary-vs-secondary distinction
  (→ da-3-1-1-primary-vs-secondary), *only* about internal-vs-external
  (→ da-3-1-2-internal-vs-external if installed), or *only* about
  structured/semi-structured/unstructured format
  (→ da-3-1-3-structured-semi-structured-unstructured if installed). Skip for
  downstream concerns: sampling design (→ da-3-data-acquisition-sampling), data
  cleaning (→ da-4-data-cleaning-preparation), or storage-layer/tooling
  questions.

version: "1.1.0"
updated: "2026-05-30"
category: "data-analysis"
tags:
  - data-sources
  - data-collection
  - taxonomy
  - primary-secondary
  - internal-external
  - structured-unstructured
  - data-inventory
  - fitness-for-purpose
related_skills:
  - da-3-data-collection-acquisition
  - da-3-data-acquisition-sampling
  - da-3-1-1-primary-vs-secondary
  - da-4-data-cleaning-preparation
---

# Data Sources: Taxonomy, Inventory, and Assessment

## How to Use This Skill

When triggered, determine which task the user needs:

- **Classify a source** → apply the three-dimension framework in the next section; state each applicable dimension explicitly. Example output: *"This CRM export is a primary source (you generated it), internal, and structured."*
- **Build or audit an inventory** → use the inventory table template in "Inventorying Data Sources"; fill in what is known, mark unknown fields as `?`, and flag gaps as action items.
- **Evaluate a source for a specific question** → work through the five fitness-for-purpose criteria in order; produce a go/no-go verdict in the form: *"[Source] is [fit / unfit / conditionally fit] for [question] because [criterion that passes or fails]. Action required: [none / document risk / find alternative]."* After drafting the assessment, re-read it to confirm each criterion is explicitly addressed and the verdict matches the criteria scores.
- **Spanning question (e.g., "what data do I have and which should I use?")** → classify each candidate source first, then run the fitness assessment on the most relevant candidates.

If the question targets only one dimension (originality, provenance, or format), defer to the appropriate sub-skill (see "Relationship to Sub-Skills").

## What Is a Data Source?

A data source is any origin point from which data can be collected for analysis.
The term encompasses both the *medium* (a database, a survey instrument, a sensor
feed, a published report) and the *process* that generated the data (direct
measurement, administrative record-keeping, secondary compilation). Because the
same dataset can serve as a primary source in one study and a secondary source in
another, source classification is always relative to the research question at hand
([University of Minnesota Crookston Library](https://crk.umn.edu/library/primary-secondary-and-tertiary-sources)).

## The Three Classification Dimensions

Data sources are typically described along three orthogonal dimensions. Each has
its own dedicated sub-skill; this skill treats them as a unified catalog view.

### 1. Originality: Primary / Secondary / Tertiary

**Primary sources** are first-hand accounts: raw survey responses, experimental
measurements, transaction logs, interviews, sensor readings, or any dataset the
analyst collected themselves. They carry maximum informational richness but also
maximum collection cost and data-quality responsibility.

**Secondary sources** analyze, synthesize, or repackage primary data. Published
research articles, government statistical summaries, market research reports, and
pre-compiled datasets all fall here. Secondary data eliminates the burden of
original collection and is often used to address research questions without the
time and cost investment of gathering original data
([UMN Crookston Library](https://crk.umn.edu/library/primary-secondary-and-tertiary-sources)).

**Tertiary sources** compile or index secondary sources. Encyclopedias, almanacs,
bibliographies, handbooks, and meta-aggregator databases belong here. Use them
for orientation and scoping (background on a topic, finding relevant primary or
secondary sources); do not cite them as direct evidence in quantitative analysis
or as the sole basis for a factual claim.

*The classification is context-dependent*: a government census is a primary source
to a demographer studying its methodology, but a secondary source to an economist
using its aggregate employment figures.

> **Single-dimension questions:** If the user asks *only* about the
> primary/secondary/tertiary distinction, defer to **da-3-1-1-primary-vs-secondary**
> rather than answering here.

### 2. Provenance: Internal / External

**Internal sources** are generated within the organization conducting the analysis.
Six canonical categories exist: operational data (transactions, inventory, orders),
human-resources data (payroll, performance records), infrastructure data (asset
registers, IT logs), communication data (emails, meeting notes), R&D data (project
outputs, test results), and customer data (CRM records, support tickets)
([Moor Insights & Strategy](https://moorinsightsstrategy.com/research-notes/enterprise-data-technology-part-1-internal-and-external-data-sources/)).

**External sources** originate outside the organization. Eight common categories:
market analysis and competitor intelligence, customer-generated content (reviews,
social media), economic indicators, regulatory and compliance documentation,
environmental and sensor data, third-party analyst reports, news feeds, and
demographic or geographic reference data
([Moor Insights & Strategy](https://moorinsightsstrategy.com/research-notes/enterprise-data-technology-part-1-internal-and-external-data-sources/)).

Unless an integration strategy is in place, internal and external data tend to live
in separate systems with different schemas, access controls, and update cadences.
Combining them requires explicit harmonization.

### 3. Format: Structured / Semi-Structured / Unstructured

**Structured data** follows a predefined schema with fixed fields, rows, and columns.
It is stored in relational databases and is directly queryable with SQL-family tools.
Examples: order tables, time-series metrics, financial ledgers
([Alation](https://www.alation.com/blog/structured-unstructured-semi-structured-data/)).

**Semi-structured data** uses organizational markers (tags, key-value pairs,
hierarchical nesting) without enforcing a rigid schema. JSON files, XML documents,
email headers, and NoSQL collections fall here. The schema can evolve without a
formal migration, but queries require parsing or path-based access patterns
([Alation](https://www.alation.com/blog/structured-unstructured-semi-structured-data/)).

**Unstructured data** carries no predefined format. Emails, documents, images,
audio, video, and free-text fields are the dominant types. Roughly 90% of
enterprise-generated data is unstructured
([Alation](https://www.alation.com/blog/structured-unstructured-semi-structured-data/)).
Effective use requires NLP, computer vision, or other content-extraction layers
before the data can be quantitatively analyzed.

## Inventorying Data Sources

An inventory is the formal catalog of all sources relevant (or potentially
relevant) to an analytical question. A minimal inventory entry records:

| Field | Description |
|---|---|
| Source name / identifier | Unique label for the dataset or feed |
| Source type | Primary / secondary / tertiary; internal / external; structured / semi / unstructured |
| Owner or custodian | Who controls access and updates |
| Coverage | Time range, geography, population, or subject domain covered |
| Update cadence | How often the data refreshes (real-time, daily, annual, one-time) |
| Access mechanism | Direct query, API, file transfer, licensed download, scrape |
| Known limitations | Gaps, exclusions, recording artifacts, or known biases |
| Provenance notes | How the data was originally generated |

A thoughtful inventory prevents the most common failure mode: being flooded with
unorganized data that is hard to collate, hard to analyze, and hard to keep secure
([Moor Insights & Strategy](https://moorinsightsstrategy.com/research-notes/enterprise-data-technology-part-1-internal-and-external-data-sources/)).

## Assessing Source Fitness for Purpose

Having a source in the inventory does not mean it is appropriate for a given
question. Apply these five criteria; a "no" on any one is a disqualifier or a risk
that must be documented before proceeding.

| Criterion | Question to ask | Example failure |
|---|---|---|
| Relevance | Does the source measure what the research question requires, or only a proxy? | Using website visit counts to infer purchase intent ([NICE](https://www.nice.org.uk/corporate/ecd9/chapter/assessing-data-suitability)) |
| Completeness | Are there systematic gaps (time windows, sub-groups, geographies) that could bias conclusions? | A customer dataset that excludes churned users, skewing retention estimates ([PMC](https://pmc.ncbi.nlm.nih.gov/articles/PMC7687257/)) |
| Freshness | Is the data current enough for the decision context? | Using a government survey with a two-year lag to inform a real-time pricing model ([NICE](https://www.nice.org.uk/corporate/ecd9/chapter/assessing-data-suitability)) |
| Accuracy and consistency | Are values within plausible ranges and definitions applied consistently across time or sites? | A CRM field where "close date" means contract sign in one region and invoice date in another ([PMC](https://pmc.ncbi.nlm.nih.gov/articles/PMC7687257/)) |
| Bias potential | Does the collection process systematically over- or under-represent certain groups or outcomes? | An administrative health database limited to insured patients that cannot generalize to the uninsured population ([PMC](https://pmc.ncbi.nlm.nih.gov/articles/PMC7687257/)) |

The SPIFD (Structured Process to Identify Fit-for-Purpose Data) framework provides
a step-by-step checklist for documenting feasibility assessments so the fitness
rationale is auditable and reproducible
([Marksman Healthcare](https://marksmanhealthcare.com/2025/01/24/enhancing-credibility-in-real-world-evidence-generation-through-space-and-spifd-frameworks/)).

## Common Pitfalls

**Treating type as fixed.** The same file can be a primary source (the raw export)
or secondary (the analyst's summary of it). Clarify which role the source plays in
the current analysis before labeling it.

**Conflating format with quality.** Structured data is not inherently more accurate
than unstructured; it is only more immediately queryable. An SQL table can contain
systematically wrong values; a set of interview transcripts can contain precise,
rich observations.

**Ignoring provenance for external sources.** External data often carries hidden
assumptions about how it was collected. Demographic benchmarks, for example, may
reflect a vendor's panel, not a probability sample of the population.

**Underweighting update cadence.** A source that was fit for purpose when the
project began may become stale during a multi-month analysis. Build refresh-date
checks into the inventory.

**Source proliferation without governance.** "Unless you take a thoughtful
approach, you could be flooded with unorganized data that is hard to collate, hard
to analyze and hard to keep secure"
([Moor Insights & Strategy](https://moorinsightsstrategy.com/research-notes/enterprise-data-technology-part-1-internal-and-external-data-sources/)).
Limit the active source set to what the research question requires.

## Relationship to Sub-Skills

This skill provides the catalog/overview layer. For depth on any single dimension:

- **da-3-1-1-primary-vs-secondary** — detailed treatment of the originality axis,
  including how to operationalize the distinction across disciplines.
- **da-3-1-2-internal-vs-external** *(if installed)* — governance, integration
  patterns, and cost/risk trade-offs for org-sourced vs. third-party data.
- **da-3-1-3-structured-semi-structured-unstructured** *(if installed)* — storage
  formats, tooling, and processing pipelines by structure type.

## Sources

1. [Primary, Secondary, and Tertiary Sources — University of Minnesota Crookston](https://crk.umn.edu/library/primary-secondary-and-tertiary-sources) — Library guide defining the three originality tiers with examples.
2. [Enterprise Data Technology Part 1: Internal and External Data Sources — Moor Insights & Strategy](https://moorinsightsstrategy.com/research-notes/enterprise-data-technology-part-1-internal-and-external-data-sources/) — Analyst research note cataloging six internal and eight external data categories.
3. [Structured, Unstructured, and Semi-Structured Data — Alation](https://www.alation.com/blog/structured-unstructured-semi-structured-data/) — Practical definitions and organizational context for all three format types.
4. [Assessing Data Suitability — NICE Real-World Evidence Framework](https://www.nice.org.uk/corporate/ecd9/chapter/assessing-data-suitability) — Regulatory-grade checklist for evaluating source fitness including relevance and freshness.
5. [Fitness for Purpose in Real-World Data Quality — PMC](https://pmc.ncbi.nlm.nih.gov/articles/PMC7687257/) — Peer-reviewed treatment of completeness, accuracy, bias, and missing-data assessment criteria.
6. [SPIFD Framework for Fit-for-Purpose Data — Marksman Healthcare](https://marksmanhealthcare.com/2025/01/24/enhancing-credibility-in-real-world-evidence-generation-through-space-and-spifd-frameworks/) — Step-by-step structured process for documenting data feasibility assessments.
