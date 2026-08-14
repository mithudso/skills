<!-- hub-reference-banner -->
> **Reference file — part of the `da-data-engineering-platform` hub.** Formerly the standalone `da-30-data-governance-catalogs` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-30-data-governance-catalogs
description: >-
  Data governance, catalogs, and discovery as a discipline — DAMA-DMBOK & DCAM
  frameworks, data catalogs & metadata management, active metadata, data
  discovery/search, table- and column-level lineage, business glossaries, data
  stewardship & ownership roles, data classification & tagging, access
  governance/policy enforcement, data products & data mesh federated
  computational governance, data contracts (ODCS/ODPS), and the tooling
  landscape (DataHub, OpenMetadata, Amundsen, Atlan, Collibra, Alation, Unity
  Catalog, Microsoft Purview).
  TRIGGER: choosing or designing a data governance framework or operating model;
  standing up or evaluating a data catalog / metadata platform; defining
  steward/owner/custodian roles or a business glossary; implementing lineage,
  data classification, tagging, or access policies; designing data products /
  data mesh governance or data contracts; comparing catalog tools.
  SKIP: privacy law, consent, regulatory compliance specifics (use
  da-11-ethics-and-privacy); data reliability/SLA monitoring & freshness alerts
  as an operational practice (use da-19-data-observability); building ETL/ELT
  pipelines themselves (use da-13-data-engineering-and-pipelines).
---

# Data Governance, Catalogs & Discovery

## Overview

Data governance is the discipline of exercising authority, control, and shared
decision-making over the management of data assets — who can take what action,
on which data, under what circumstances, using what methods. Catalogs and
discovery are the operational layer that makes governance executable: a metadata
platform that inventories assets, attaches meaning (glossaries, tags,
classifications), traces movement (lineage), assigns accountability
(stewardship), and exposes it all through search so people can *find and trust*
data.

This skill covers governance **as a practice and an architecture**, not the
adjacent disciplines. Privacy law and ethics live in `da-11`; operational
reliability/freshness monitoring lives in `da-19`; pipeline construction lives
in `da-13`. Here the focus is: frameworks → metadata → discovery → lineage →
meaning → roles → classification → access → products/mesh → tooling.

Two macro-shifts define the 2024–2026 landscape:
1. **Passive → active metadata.** Catalogs stop being static inventories and
   become bidirectional orchestration layers that push metadata back into the
   stack to drive automation (Gartner; Atlan).
2. **Centralized → federated governance.** Data mesh reframes governance as
   *federated computational governance* — global rules enforced
   computationally, local ownership by domain teams (Dehghani; Fowler).

## Core Concepts

### 1. Governance frameworks: DAMA-DMBOK and DCAM
- **DAMA-DMBOK** (Data Management Body of Knowledge, DAMA International) organizes
  data management into **11 knowledge areas** rendered as the "DAMA wheel" with
  **Data Governance at the hub**: governance, architecture, modeling, storage &
  operations, security, integration & interoperability, document & content
  management, reference & master data, data warehousing & BI, metadata, and data
  quality. Governance is the connective function underpinning the other ten.
- **DCAM** (Data Management Capability Assessment Model, EDM Council) is a
  maturity-assessment standard organized into **eight core components** (data
  strategy, business case & funding, governance, architecture, technology,
  quality, operations, control environment). **DCAM v3** (released 2024) is the
  current standard. DCAM measures *how mature* your program is; DMBOK describes
  *what the disciplines are*. They are complementary — DMBOK for vocabulary and
  scope, DCAM for assessment and roadmap.
- Sibling: **CDMC** (Cloud Data Management Capabilities, EDM Council) extends
  DCAM-style assessment to cloud + sensitive-data controls.

### 2. Metadata management & active metadata
- **Metadata** = data about data, conventionally split into **technical**
  (schemas, types, partitions), **business** (definitions, glossary terms,
  ownership), and **operational** (run logs, freshness, query frequency).
- **Passive metadata** sits in a static catalog and is read by humans. **Active
  metadata** is continuously analyzed, curated, and *pushed back* into tools to
  drive automation — e.g., auto-propagating a PII tag from a source column to
  every downstream table, or surfacing the most-queried tables first. Metadata
  moves **both directions**. Gartner projected ~30% of orgs adopting active
  metadata by 2026 with up to 70% faster time-to-delivery of data assets, and
  reframed metadata management as foundational to **AI readiness** in its 2025
  Magic Quadrant (first refresh in five years, Nov 2025).

### 3. Data discovery & search
- Discovery is the consumer-facing entry point: search across all assets, ranked
  by relevance, enriched with ownership, quality, popularity, and lineage so a
  user can judge **trustworthiness**, not just existence. Modern catalogs add
  natural-language/conversational search (Atlan), and **query-log ingestion**
  (Alation) that mines actual execution patterns to rank and recommend assets.

### 4. Data lineage — table-level and column-level
- **Table-level lineage** answers "which datasets feed which." **Column-level
  lineage** maps dependencies field-by-field, enabling precise impact analysis
  ("if I drop this column, what breaks?") and root-cause analysis.
- Lineage is typically derived by **SQL parsing**. Parser choice matters:
  DataHub uses **SQLGlot** (schema-aware, highest correct-lineage rate across
  dialects); OpenMetadata uses **sqllineage**; OpenLineage/Marquez uses
  **openlineage-sql**. Native column-level support exists for Snowflake,
  BigQuery, Databricks, and BI tools (Looker, Power BI, Tableau).

### 5. Business glossaries
- A **business glossary** is the controlled vocabulary of agreed business terms
  (e.g., "active customer") with definitions, owners, and relationships, linked
  to physical assets so technical columns inherit business meaning. Distinct
  from a **data dictionary** (technical, schema-level) and a **taxonomy/ontology**
  (hierarchical/semantic relationships). Stewards own glossary curation.

### 6. Data stewardship & ownership roles
- **Data Owner** — accountable (usually senior business role) for classification,
  protection, use, and quality of a data domain; results-focused; signs off on
  the glossary and access policy.
- **Data Steward** — responsible for quality, definitions, documentation,
  glossary, and lineage; task-focused; the day-to-day governance operator.
- **Data Custodian** — IT role; implements and maintains the storage/security
  controls the Owner specifies; handles access provisioning, incident review,
  platform monitoring.
- RACI mnemonic: Owner = Accountable, Steward = Responsible (for meaning/quality),
  Custodian = Responsible (for the technical controls).

### 7. Data classification & tagging
- **Classification** assigns sensitivity levels (public / internal / confidential
  / restricted) and compliance categories (PII, GDPR, HIPAA, etc.). Modern
  platforms auto-discover and tag sensitive data — Unity Catalog uses an
  **agentic/AI classifier** for continuous PII discovery. **Governed tags**
  enforce a controlled, consistent tag vocabulary (vs. free-form tags) so
  policies can key off them reliably.

### 8. Access governance & policy enforcement
- Move from per-object grants to **policy-as-data**. **ABAC**
  (attribute-based access control) evaluates tag-based conditions and applies
  **row filters** (which rows you see) and **column masks** (what values you see)
  automatically across catalogs/schemas — e.g., mask any column tagged `PII`
  unless the user is in `pii-readers`. Unity Catalog made ABAC row filters,
  column masks, governed tags, and data classification **GA in 2025**. Pair
  classification (find sensitive data) + tags (label it) + ABAC (enforce) for
  scalable, declarative governance.

### 9. Data products, data mesh & federated computational governance
- **Data mesh** (Zhamak Dehghani, ThoughtWorks 2019) is a sociotechnical
  approach to analytical data at scale, built on **four principles**:
  domain-oriented ownership, **data as a product**, self-serve data platform,
  and **federated computational governance**.
- **Federated computational governance** = a decision model led by a federation
  of domain + platform product owners with local autonomy, adhering to **global
  rules that are enforced computationally** (encoded into the platform, not
  enforced by committee). Goal: a mesh of independent data products that is
  interoperable, secure, and trustworthy.
- A **data product** is the smallest architectural unit encapsulating everything
  needed to share data (data + metadata + code + access + SLOs), owned by the
  domain team.

### 10. Data contracts & specifications
- A **data contract** is an enforceable agreement between producer and consumer
  covering schema, semantics, quality, and SLAs. **ODCS** (Open Data Contract
  Standard, v3.x, governed by **Bitol**, a Linux Foundation AI & Data project;
  originated at PayPal) defines schema-level executable contracts. **ODPS**
  (Open Data Product Specification) is broader — design, publish, discover,
  monetize, and govern data products as business value units, and can reference
  ODCS contracts inline or by URL. Use ODCS for the interface; ODPS for the
  product wrapper.

## Tools & Frameworks

| Tool | Type / License | Strengths | Notes |
| --- | --- | --- | --- |
| **DataHub** | OSS (Apache 2.0) | Event-driven, Kafka + Elasticsearch + graph (JanusGraph/Neo4j) + GMS; SQLGlot column-level lineage; federated metadata services | Engineering-heavy; "context platform" positioning |
| **OpenMetadata** | OSS (Apache 2.0) | Simple 4-component stack (MySQL/Postgres + Elasticsearch, no graph DB); 90+ connectors; 700+ JSON Schemas; discovery+governance+quality+lineage in one | Maintained by Collate; weekly releases |
| **Amundsen** | OSS (Lyft) | Minimalist, search-relevance-first discovery | Lighter governance; popular in eng orgs |
| **Apache Atlas** | OSS | Hadoop-ecosystem lineage & classification | Legacy/Hadoop-centric |
| **Atlan** | Commercial | Active metadata, NL/conversational search, automated column-level lineage; 4–6 wk setup | Gartner MQ + Forrester Wave Leader 2025 |
| **Collibra** | Commercial | Governance orchestration, formal stewardship workflows | Best for regulated enterprises; 3–9 mo rollout |
| **Alation** | Commercial | Query-log-driven active metadata, analytics-first | 6–12 wk; cross-system lineage gaps reported |
| **Unity Catalog** | Databricks (OSS core) | Native classification, governed tags, ABAC row/column policies (GA 2025) | Enforcement engine inside Databricks |
| **Microsoft Purview** | Azure | Cross-Azure technical metadata + discovery + classification | Pairs with Unity Catalog (Purview = discovery, UC = enforcement) |

Frameworks: **DAMA-DMBOK** (scope/vocabulary), **DCAM v3** + **CDMC** (maturity
assessment), **data mesh** (federated operating model), **ODCS/ODPS** (contracts
& product specs).

## Methodology — standing up governance + a catalog

1. **Frame the operating model.** Decide centralized vs. federated (mesh). Map
   domains. Assign Owner/Steward/Custodian per domain (RACI). Use DCAM to
   baseline current maturity and set a roadmap.
2. **Pick the scope that delivers value first.** Start with the highest-value or
   highest-risk domains, not a boil-the-ocean catalog of everything.
3. **Ingest technical metadata.** Connect sources (warehouses, lakes, BI, dbt)
   to the catalog. Auto-harvest schemas + lineage. Verify column-level lineage
   coverage for your key dialects.
4. **Layer meaning.** Build the business glossary; link terms to physical
   assets. Stewards curate.
5. **Classify & tag.** Run automated sensitive-data classification; apply
   *governed* (controlled-vocabulary) tags, not free-form.
6. **Enforce access declaratively.** Define ABAC policies keyed on tags (mask
   PII, row-filter by region). Test that policies propagate across schemas.
7. **Activate the metadata.** Wire automation: tag propagation along lineage,
   freshness/popularity signals into search ranking, push back to source tools.
8. **Operationalize.** Stewardship rituals, glossary review cadence, data
   contracts (ODCS) on critical interfaces, product specs (ODPS) for shared
   data products.
9. **Measure.** Coverage (% assets cataloged/owned/classified), adoption (search
   usage, time-to-find), trust (% certified assets), policy compliance.

## Practical Patterns

- **Certify, don't catalog everything.** A "verified/certified" badge on
  trusted assets beats 100% coverage of unmanaged junk. Discovery is about
  trust, not census.
- **Tag-driven policy.** Classify → governed tag → ABAC. One policy ("mask
  `PII`") covers thousands of objects and auto-applies to new ones.
- **Propagate along lineage.** Use column-level lineage to auto-inherit
  classifications/tags downstream so a single source label cascades.
- **Glossary terms as the bridge.** Bind business terms to physical columns so
  non-technical users search in business language.
- **Federated rules, central platform.** In mesh, encode global rules
  computationally in the self-serve platform; let domains own products within
  those rails.
- **Contracts on the boundaries.** Put ODCS data contracts on cross-domain /
  producer-consumer interfaces where breakage is expensive; don't contract
  everything.
- **Query logs for relevance.** Rank search and recommend assets by actual usage
  (Alation-style), not alphabetical or last-modified.

## Anti-Patterns

- **Catalog as a graveyard.** A one-time bulk ingest with no stewardship, no
  owners, stale within months. Governance is a continuous program, not a
  project.
- **Passive metadata only.** Treating the catalog as a read-only wiki — no
  automation, no push-back into tools. Metadata that doesn't *act* decays.
- **Free-form tag sprawl.** Uncontrolled tags (`pii`, `PII`, `personal`,
  `sensitive`) make policies unreliable. Use governed/controlled vocabularies.
- **Governance by committee bottleneck.** Central team must approve every
  change — kills velocity. Federate ownership; enforce computationally.
- **Owner/steward/custodian conflation.** One overloaded "data person" can't be
  accountable, responsible for meaning, *and* run the platform. Separate the
  roles or accountability evaporates.
- **Table-level lineage where column-level is needed.** Impact analysis on a
  schema change is guesswork without field-level lineage.
- **Tool-first, model-last.** Buying Collibra/Atlan before defining domains,
  roles, and policies yields shelfware. Operating model first.
- **Boil-the-ocean rollout.** Cataloging every asset before any are governed.
  Start narrow, prove value, expand.

## Troubleshooting

- **Lineage is incomplete / missing columns.** Check parser/dialect support
  (SQLGlot vs sqllineage), ensure the catalog has schema context, and confirm
  the connector ingests query history (not just DDL). Dynamic SQL and
  `SELECT *` degrade column-level resolution.
- **Sensitive data slipping through.** Automated classification missed it —
  re-run/expand classifiers, add custom patterns, and propagate along lineage so
  derived columns inherit the tag.
- **ABAC policy not applying.** Verify the object actually carries the governed
  tag the policy keys on, that the policy is at the right catalog/schema scope,
  and that classification ran before policy evaluation.
- **Low catalog adoption.** Usually a trust/relevance problem: no owners, no
  certification, poor search ranking. Add ownership, certify key assets, rank by
  query-log popularity, link glossary terms.
- **Purview ↔ Unity Catalog drift.** Schema/lineage/classification out of sync —
  confirm the connector/API sync cadence; remember Purview is
  discovery/technical-metadata and UC is the enforcement plane.
- **Glossary nobody uses.** Terms not linked to physical assets, or no steward
  cadence. Bind terms to columns and put glossary review in the stewardship
  ritual.
- **Mesh governance chaos.** Global rules defined but not *computational* —
  encode them into the self-serve platform; "agreement" docs don't enforce.

## References

**Frameworks**
- DAMA International — DAMA-DMBOK (Data Management Body of Knowledge): https://www.dama.org/cpages/body-of-knowledge ; https://www.damadmbok.org/ (2024)
- Snowflake — "DAMA-DMBOK Explained": https://www.snowflake.com/en/fundamentals/data-governance/framework/dama-dmbok/ (2024)
- Atlan — "DAMA DMBOK Framework Guide": https://atlan.com/dama-dmbok-framework/ (2025)
- EDM Council — "Announcing DCAM v3": https://edmcouncil.org/announcement/announcing-dcam-v3-meet-the-new-standard-for-your-data/ (2024)
- EDM Council — DCAM framework: https://edmcouncil.org/frameworks/dcam/ (2024)
- Snowflake — "DCAM Explained": https://www.snowflake.com/en/fundamentals/data-governance/framework/dcam/ (2024)

**Metadata & active metadata**
- Gartner — Magic Quadrant for Metadata Management Solutions (Nov 19, 2025): via Informatica/Atlan: https://www.informatica.com/metadata-management-magic-quadrant.html ; https://atlan.com/gartner-magic-quadrant-for-metadata-management/ (2025)
- Gartner — Market Guide for Active Metadata Management: https://www.gartner.com/en/documents/4004082 (2024)
- Atlan — "Gartner Active Metadata Management Guide": https://atlan.com/gartner-active-metadata-management/ (2026)
- OvalEdge — "Active Metadata Management": https://www.ovaledge.com/blog/active-metadata/ (2024)

**Lineage**
- DataHub — "SQL Lineage: How DataHub's Column-Level Parser Works": https://datahub.com/blog/extracting-column-level-lineage-from-sql/ (2024)
- DataHub Docs — Lineage feature guide: https://docs.datahub.com/docs/features/feature-guides/lineage (2025)
- OpenMetadata Docs — "How Column-Level Lineage Works": https://docs.open-metadata.org/latest/how-to-guides/data-lineage/column (2025)

**Discovery, glossary, stewardship**
- Atlan — "Alation vs Collibra vs OpenMetadata vs Atlan": https://atlan.com/alation-vs-collibra-vs-openmetadata-vs-atlan/ (2025)
- Atlan — "16 Best Data Catalog Tools": https://atlan.com/data-catalog-tools/ (2026)
- DQOps — "Data Owner vs Data Steward vs Data Custodian": https://dqops.com/data-owner-data-steward-data-custodian-roles/ (2024)
- EWSolutions — "Data Stewardship Roles": https://www.ewsolutions.com/data-stewardship-roles-a-complete-guide/ (2024)
- Atlan — "Data Stewardship 101": https://atlan.com/data-stewardship-101/ (2025)

**Classification & access governance**
- Databricks — "Find Sensitive Data at Scale with Data Classification in Unity Catalog": https://www.databricks.com/blog/find-sensitive-data-scale-data-classification-unity-catalog (2025)
- Databricks — "ABAC row filtering and column masking... GA in Unity Catalog": https://www.databricks.com/blog/abac-row-filtering-and-column-masking-policies-governed-tags-and-data-classification-are-now (2025)
- Microsoft Learn — Unity Catalog Data Classification: https://learn.microsoft.com/en-us/azure/databricks/data-governance/unity-catalog/data-classification (2025)

**Data mesh, products & contracts**
- Zhamak Dehghani — "Data Mesh Principles and Logical Architecture" (martinfowler.com): https://martinfowler.com/articles/data-mesh-principles.html (2020)
- Dehghani — *Data Mesh: Delivering Data-Driven Value at Scale*, O'Reilly (2022)
- Atlan — "Data Mesh Principles (Four Pillars)": https://atlan.com/data-mesh-principles/ (2025)
- Starburst — "Federated Computational Governance": https://www.starburst.io/blog/data-mesh-book-bulletin-principle-of-federated-computational-governance/ (2024)
- Bitol / Linux Foundation — Open Data Contract Standard (ODCS) v3.x: https://github.com/bitol-io/open-data-contract-standard ; https://bitol-io.github.io/open-data-contract-standard/ (2025)
- Open Data Product Specification (ODPS) — "ODPS vs ODCS": https://blog.opendataproducts.org/when-standards-collide-clarifying-odps-and-odcs-in-the-data-product-landscape-c2978f9c13d9 (2025)

**Tooling architecture**
- DataHub Docs — Architecture Overview: https://docs.datahub.com/docs/architecture/architecture (2025)
- Atlan — "OpenMetadata Explained": https://atlan.com/openmetadata-explained/ (2025)
- TheDataGuy — "Open-Source Data Governance Frameworks: OpenMetadata, DataHub, Atlas, Amundsen": https://thedataguy.pro/writing/2025/08/open-source-data-governance-frameworks/ (2025)
- OpenMetadata Standards: https://openmetadatastandards.org/ (2025)
