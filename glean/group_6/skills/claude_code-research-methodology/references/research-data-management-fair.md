<!-- hub-reference-banner -->
> **Reference file — part of the `research-methodology` hub.** Formerly the standalone `research-data-management-fair` skill.
> Sibling topics in this family are now reference files under the hubs (`research-methodology`, `deep-research`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: research-data-management-fair
title: "Research Data Management & FAIR Data Stewardship"
description: >-
  Practitioner guide to documenting, licensing, and depositing research datasets so they are Findable, Accessible, Interoperable, Reusable (FAIR). Covers FAIR principles, Data Management Plans (DMPs; NIH DMS/NSF), dataset metadata standards (DataCite, Dublin Core, schema.org/Dataset), repository selection (Zenodo, Dryad, Figshare, OSF, re3data), persistent identifiers (DOIs, ORCID), data licensing (CC0 vs CC-BY), versioning, provenance, the research-data lifecycle. TRIGGER: documenting/depositing a dataset; writing a DMP or data-sharing plan; making data FAIR/reproducible/citable; choosing a repository; minting a DOI; building a dataset codebook/data dictionary/README (variable definitions for reuse); picking a data license. SKIP: coding/annotation codebooks & label schemas (inter-annotator agreement, thematic coding) -> research-data-annotation-coding; engagement/Drive doc taxonomy -> content-ingestion-extraction; pipelines/warehousing -> da-*; operational DB backup -> mongodb-*.
category: custom
version: 1.2.0
updated: "2026-06-15"
whenToUse:
  - "documenting or depositing a research/analysis dataset"
  - "writing a Data Management Plan (DMP/DMSP) or data-sharing plan"
  - "making data FAIR, reproducible, or citable; choosing a repository; minting a data DOI"
  - "building a dataset codebook / data dictionary (variable & value definitions for reuse) or dataset README"
  - "picking a data license (CC0 vs CC-BY); metadata standards for a dataset"
keywords:
  - FAIR principles
  - research data management
  - data management plan
  - DataCite
  - data repository
  - persistent identifier
  - dataset codebook
  - data dictionary
  - data license
  - data provenance
tags:
  - research
  - data-management
  - fair
  - data-stewardship
metadata:
  changelog:
    - "2026-06-15 sko v1->v1.2.0 — Pass H 10/10->10/10 pos, 9/10->10/10 neg; 14 fixed (codebook collision SKIP edge, NIH 2026 DMS format, DataCite 4.7, file/folder gap, body scope-out, §3 table)"
---

# Research Data Management & FAIR Data Stewardship

Research Data Management (RDM) is the practice of organizing, documenting, storing, and
preserving the datasets a project produces or reuses, so that they remain understandable
and usable long after the analysis is done. FAIR is the dominant quality standard for that
work: data and metadata should be **F**indable, **A**ccessible, **I**nteroperable, and
**R**eusable — by machines as well as humans.

This skill is for the practitioner who has a *research or analysis dataset in hand* and must
make it reusable, citable, and reproducible: where to put it, how to describe it, how to
license it, and what "good" looks like.

**Out of scope — use instead:**
- *Coding/annotation codebooks and label schemas* (the rules an annotator applies to **produce** coded/labeled data, inter-annotator agreement, thematic coding) → `research-data-annotation-coding`. This skill's codebook is the **dataset-documentation** sense (§5): describing variables in data you already have.
- *Engagement / Drive document taxonomy* (numbered folder trees, `_meta` indexes) → `content-ingestion-extraction` (references/doc-store-bootstrapper.md).
- *Production data pipelines, warehousing, and platform* (ETL/ELT, dbt, lakehouse, governance catalogs) → the `da-*` data-engineering skills.
- *Operational database storage, backup, and DR* → the `mongodb-*` skills.

> FAIR is not the same as "open." FAIR data can be access-controlled (A1.2 allows
> authentication/authorization). The goal is *reusability under clearly stated conditions*,
> not unconditional disclosure. Sensitive data can be FAIR while staying restricted.

## 1. Core concept: the FAIR principles

FAIR was published by Wilkinson et al. (2016, *Scientific Data*) and is maintained as a
living document by GO FAIR. The principles deliberately emphasize **machine-actionability** —
software should be able to find, access, and interpret data with minimal human help. The
principles apply to three entity types: **data**, **metadata** (information *about* the data),
and **infrastructure** (the searchable resources that index them).

**Findable**
- **F1.** (Meta)data are assigned a globally unique and persistent identifier (PID). *The most
  important clause* — most other FAIR properties depend on having a PID (e.g. a DOI).
- **F2.** Data are described with rich metadata (see R1).
- **F3.** Metadata explicitly include the identifier of the data they describe.
- **F4.** (Meta)data are registered/indexed in a searchable resource (repository, catalog,
  Google Dataset Search).

**Accessible**
- **A1.** (Meta)data are retrievable by their identifier over a standardized protocol.
  - **A1.1** the protocol is open, free, universally implementable (HTTP(S), not a proprietary API).
  - **A1.2** the protocol supports authentication/authorization *where necessary* — this is how
    restricted data stays FAIR.
- **A2.** Metadata remain accessible **even when the data are gone** (a "tombstone" landing page
  for withdrawn data).

**Interoperable**
- **I1.** (Meta)data use a formal, accessible, shared, broadly applicable knowledge-representation
  language (e.g. RDF, JSON-LD; controlled vocabularies; not a bespoke encoding).
- **I2.** Vocabularies used in metadata are themselves FAIR.
- **I3.** (Meta)data include qualified references to other (meta)data (typed links: *isSupplementTo*,
  *isDerivedFrom*).

**Reusable**
- **R1.** (Meta)data are richly described with plural, accurate, relevant attributes.
  - **R1.1** released with a clear, accessible **data-usage license** (see §6).
  - **R1.2** associated with **detailed provenance** — where the data came from, who made it,
    how it was processed (see §7).
  - **R1.3** meet **domain-relevant community standards** (the metadata schema/format your field
    expects).

A common self-check is the **FAIR maturity** lens: don't ask "is it FAIR (yes/no)," ask "how FAIR,
on each axis." Tools like FAIRsharing and F-UJI score datasets against these clauses.

## 2. Data Management Plans (DMPs)

A **Data Management Plan** (NIH calls it a **Data Management & Sharing Plan / DMSP**) is a short
formal document describing what you will do with data *during and after* a project. Most US federal
grants and many foundations and journals now require one. Write it at the *start* — planning ahead
is where the payoff is. **DMPTool** (dmptool.org) ships funder-specific templates and prompts.

**NIH** — the 2023 **DMS Policy** stands (an approved plan is a *term and condition of award*), but
the *plan format changed* in 2026 (NOT-OD-26-046). For due dates **on or after May 25, 2026** use the
**2026 Pilot DMS Plan format** — a simpler structured form (mostly yes/no questions plus a
data-type → repository table), not the old narrative. Required elements:
1. A **table** listing the **key scientific data types** anticipated and, for each, the **repository
   (or an example)** where it will be managed and shared (≤100 words). NIH expects established repositories.
2. Yes/No: are the data covered by the **GDS (genomic data sharing) Policy**, and can you meet its
   accelerated-timeline and Institutional-Certification expectations? (If no → element 4.)
3. Yes/No: will the data be **shared** and **preserved** in that repository?
4. If you answered "no" to 1–3 or sharing is otherwise limited, **describe the limitation** and the
   ethical/legal/technical reason (≤300 words).

The narrative *six-element* format (data type; tools/software/code; standards; preservation/access/
timelines; access-distribution-reuse considerations; oversight) is the **pre-May-25-2026 (legacy)
format** — still accepted for earlier due dates, and its content areas still inform the new form.

**NSF** asks for analogous content under PAPPG: data types/products, metadata standards, access &
sharing policies (with privacy/IP protections), reuse/redistribution/derivatives policy, and
archiving/preservation plans. If a project produces no data, you submit a justification in place of
the plan.

> Practitioner move: budget for data management. NIH expects preservation costs to be paid *during*
> the performance period (e.g. long-term repository fees paid before the grant ends).

## 3. Metadata standards for datasets

Metadata is what makes F2, F3, I1, and R1 real. Three generic standards cover almost every case:

| Standard | What it is / where it is read | Use it for |
| --- | --- | --- |
| **DataCite** (v4.7, Mar 2026; check schema.datacite.org/versions for the latest) | The schema behind data **DOIs**; most repositories populate it when you mint a DOI | The DOI record for any deposited dataset |
| **Dublin Core** | 15-element lowest-common-denominator vocabulary; DataCite publishes an official **DataCite → Dublin Core** crosswalk, so a good DataCite record degrades gracefully into DC | Interop/exchange fallback |
| **schema.org/Dataset** (JSON-LD) | What **Google Dataset Search** and web crawlers read; reputable repositories emit it automatically | Open-web findability of a dataset on your own landing page (`name`, `description`, `creator`, `license`, `distribution`, `identifier`, `keywords`) |

DataCite properties have three obligation levels (Mandatory / Recommended / Optional):
- **Mandatory (must supply):** `Identifier` (the DOI), `Creator`, `Title`, `Publisher`,
  `PublicationYear`, `ResourceType` (e.g. `Dataset`, `Software`).
- **High-value Recommended:** `Subject` (keywords), `Contributor`, `Date`, `RelatedIdentifier`
  (typed links to papers/software/versions, which is I3/R1.2 in practice), `Description`
  (especially `descriptionType="Abstract"`), `Rights` (the license).
- If a mandatory value is genuinely unavailable, DataCite defines machine-recognizable codes
  (`:unkn`, `:none`, `:null`) rather than leaving it blank.

**Rule of thumb:** pick the **domain/community standard** first (R1.3) — e.g. ISA-Tab/MIAME (omics),
DDI (social-science surveys), EML (ecology), CF conventions (climate). DataCite/DC/schema.org sit
*underneath* as the generic discovery layer; the domain schema carries the science.

## 4. Repositories & persistent identifiers

### Choosing a repository — tiered decision

Follow the widely-used **domain → institutional → generalist** tier (ARDC, NIH/GREI, Turing Way):

1. **Domain / discipline-specific repository first.** Most findable and reusable to your peers; built
   for your data types and community metadata standards. Examples: GenBank, GEO (genomics), PDB
   (structures), PANGAEA (earth/environment), ICPSR (social science). Discover via **re3data.org**
   and **FAIRsharing.org**.
2. **Institutional repository** if no domain fit (and some institutions require at least a *metadata
   record* even when data live elsewhere).
3. **Generalist repository** as the fallback: **Zenodo** (CERN-backed, free, DOI, large files,
   GitHub release integration), **Dryad** (curated, publication fee, CC0-only), **Figshare** (free
   deposit, DOIs, many formats), **OSF** (project workspace + DOIs for public components),
   **Mendeley Data**, **Vivli** (clinical).

**Selection criteria** (NIH/GREI checklist — score candidates on these):
- Assigns a **persistent identifier** (DOI) and a stable landing page.
- **Long-term sustainability** / preservation commitment and documented **retention policy**.
- **Curation / quality-assurance** services.
- Free, easy **access**; supports **broad and measured reuse** with clear use guidance.
- Supports your **file formats** and **size**; offers needed **access controls** (embargo,
  restricted, application-mediated) for sensitive data.
- Records **data provenance / versioning**.
- Trustworthiness — prefer repositories certified under **CoreTrustSeal** or aligned to the **TRUST**
  principles (Transparency, Responsibility, User focus, Sustainability, Technology).

> Match the repository to funder/journal/community requirements — those often *mandate* a specific
> repository (e.g. genomic data to an NIH-designated archive).

### Persistent identifiers (PIDs)

- **DOI (Digital Object Identifier)** — the standard PID for datasets/software, resolvable and
  citable, minted via DataCite by the repository. Satisfies F1 and underpins citability.
- **Concept DOI vs version DOI** — Zenodo (and others) mint a *version DOI* for each deposit and one
  *concept DOI* that always resolves to the latest version. **Cite the version DOI** for reproducibility;
  advertise the concept DOI for "always-current."
- **ORCID** — a PID for *researchers*. Put ORCIDs on creators/contributors so credit is unambiguous and
  machine-linkable.
- **ROR** (Research Organization Registry, for organizations), **RAiD** (Research Activity Identifier,
  for projects), and **RRID** (Research Resource Identifier, for reagents/tools/software) round out the
  PID ecosystem; DataCite added `Project`/`Award` resource types in 4.6 and `RAiD` as a
  related-identifier type in 4.7.

A dataset DOI lets others cite data as a first-class research output. Provide a ready-to-paste
**data citation** (Creator (Year). *Title*. Publisher. Version. DOI) on the landing page.

## 5. Codebooks, data dictionaries & READMEs

Metadata describes the dataset as a whole; **codebooks/data dictionaries** describe it *variable by
variable*; a **README** orients a stranger to the files. All three are what make R1/reusability concrete
— "documentation that lets someone with no prior knowledge of the project understand exactly what the
data mean" (UK Data Service).

### Data dictionary / codebook

A **data dictionary** gives a meaningful description of each variable and value. A **codebook** is the
survey-oriented cousin (adds question text and response codes).

> This is the **dataset-documentation** codebook — it describes variables in data you *already have*, so
> a reuser understands them. Do not confuse it with a *coding/annotation codebook* (the label schema and
> rules an annotator applies to **produce** coded data, paired with inter-annotator agreement); that is
> `research-data-annotation-coding`.

For each variable record (ICPSR fields):
- **Variable name** — the column name/number (mnemonic `EMPLOY1`, pattern `VAR001`, or question-based `Q2b`).
- **Variable label** — short human description; for survey/instrument items copy the exact question or measure wording, for derived variables state the derivation.
- **Question text** — the exact survey wording for any item-derived variable.
- **Values** — the coded values present (`1,2,3,4,5`).
- **Value labels** — what each code means (`1=Excellent … 5=Poor`).
- **Units** and **data type** (integer, date, categorical, free text).
- **Summary statistics** — frequencies/% for categoricals; min/max/median for continuous (a quick
  integrity check for reusers).
- **Missing-data codes** — *every* missing convention, including system-missing and blank (e.g.
  `-1 = Refused`, `-9 = Not applicable`). Undocumented missingness silently biases reuse.
- **Skip / universe patterns** — which population a variable applies to.
- **Notes / provenance** — derivations, transformations, source citation for copyrighted instruments.

> The test: a data dictionary is good when every variable is **self-explanatory to someone outside your
> research group**.

### README (one per dataset)

The widely-adopted **Cornell Data Services** convention is **one plain-text README per dataset**, a
standalone file covering:
- **General information** — title; abstract/summary; creators + ORCIDs + contacts; date/location of
  collection; funding; the DOI/citation.
- **Methodological information** — how data were collected/generated/processed; instruments/software
  and versions; quality-assurance steps.
- **Data & file overview** — list every file, its format, relationships between files, naming
  conventions, and the version (with a brief changelog).
- **Data-specific information** — variable lists / pointer to the data dictionary; codes and
  abbreviations; units; missing-value codes.
- **Sharing & access** — license (R1.1), any access restrictions/embargo, links to related
  publications/software (RelatedIdentifier).

Keep README/codebook in **open, plain formats** (`.txt`, `.md`, `.csv`) so they outlive proprietary tools.

### File & folder organization

A reuser navigates the deposit before they read a byte of data, so structure and names are part of the
documentation:
- **File names**: no spaces or special characters (use `_` or `-`); date components in ISO-8601
  (`YYYYMMDD`) so they sort chronologically; semantic, consistent stems (`survey-2025_clean_v2.csv`,
  not `final_FINAL.csv`); keep names short but self-describing.
- **Folder layout**: separate `raw/` (never edited) from `processed/`/`derived/`; keep `docs/`
  (README, codebook, DMP), `scripts/`/`code/`, and `data/` distinct so provenance is legible.
- **One documented convention, stated in the README.** The naming/folder scheme itself goes in the
  README's "Data & file overview" so the pattern is explicit, not guessed.
- **Manifest for large deposits**: a file-list with a checksum (MD5/SHA-256) per file lets a reuser
  verify nothing was truncated or corrupted in transit (fixity).

## 6. Data licensing

A license (FAIR R1.1) tells reusers (and their lawyers and software) what they may do. Without one,
reuse is legally uncertain and the data effectively isn't reusable.

**The core rule (OpenAIRE / Creative Commons):**
- **Datasets / databases → CC0** (public-domain dedication). It waives copyright *and* sui-generis
  database rights, removing all uncertainty and preventing **attribution stacking** (the practical mess
  when many attribution-required datasets are combined). With CC0 you *request* citation as a scholarly
  norm rather than legally requiring it. Some generalist repos (Dryad) are CC0-only.
- **Creative/"work"-like outputs, or where attribution must be legally enforced → CC-BY 4.0.** CC is
  recognized as the interoperable "gold standard" — one condition (attribution), most easily combined.
- **Avoid NC (non-commercial) and ND (no-derivatives)** for open research data — they break Open-Access
  compliance and limit reuse/integration.
- **Database-specific alternatives:** **ODbL** and **ODC-BY** apply only to sui-generis database rights
  and the structure, *not* the contents — narrower than CC; use only with a clear reason.

The 2026 Creative Commons guidance: pick **CC0 or CC-BY**, and expose license + attribution as both
**human- and machine-readable** metadata so the license travels with the data.

> For **sensitive / personal data**, licensing is not enough: de-identify (pseudonymize/anonymize),
> place behind controlled access (A1.2), and document consent terms before deposit. License the
> *shareable* version.

## 7. Versioning & provenance

**Provenance (R1.2)** is the data's origin story: who generated/collected it, how it was processed, prior
publication, and which parts derive from someone else's (possibly transformed) data. Capture the
**workflow** that produced the dataset — ideally machine-readable (W3C PROV, or at minimum a methods
section in the README + scripts).

**Versioning** keeps reuse reproducible:
- Treat the deposited dataset as **immutable**; publish corrections as a **new version** with its own
  version DOI (cite versions, not "latest").
- Maintain a **changelog** and a clear version label in metadata (`Version` property) and README.
- Link versions and related outputs with typed `RelatedIdentifier`s (`IsNewVersionOf`,
  `IsSupplementTo`, `IsDerivedFrom`).
- Version **code and environment** alongside data (Git + a tagged release archived to Zenodo;
  `requirements`/`environment` capture) so the analysis is rerunnable.

## 8. The research-data lifecycle

The **DCC Curation Lifecycle Model** frames RDM as a cycle, not a final step, and its central lesson is
that *the metadata and documentation needed for reuse must be captured at creation/collection, not bolted
on at the end.* A pragmatic lifecycle:

1. **Plan** — write the DMP; choose standards and a target repository up front; plan consent and
   de-identification.
2. **Create / collect** — capture descriptive + structural + technical metadata *as you go*; start the
   README and data dictionary on day one.
3. **Process** — clean, validate, transform; record every transformation (provenance).
4. **Analyze** — keep analysis code under version control; preserve the environment.
5. **Appraise & select** — decide what to keep long-term vs dispose (per documented policy/legal needs).
6. **Preserve / ingest** — deposit the **final, cleaned, quality-assured, documented** version in the
   repository; convert to open/preservation formats; get the DOI.
7. **Access & share** — apply the license, set any embargo/access controls, announce the dataset
   (data citation, link from the paper, schema.org markup).
8. **Reuse** — others (and future-you) discover, cite, and build on it, closing the loop.

## 9. Practical FAIR checklist

Before you call a dataset done:

- [ ] **PID minted** — DOI via a repository; ORCIDs on all creators (F1).
- [ ] **Right repository** — domain-first, else institutional, else generalist; trustworthy (CoreTrustSeal/TRUST); supports your formats + access needs (F4, A1).
- [ ] **Rich metadata** — DataCite mandatory fields + Abstract + keywords + RelatedIdentifiers; domain standard applied (F2, I1, R1.3).
- [ ] **README** — one per dataset, plain text, covering general/methods/files/variables/sharing.
- [ ] **Data dictionary / codebook** — every variable: name, label, values, value labels, units, missing codes; self-explanatory to an outsider.
- [ ] **License attached** — CC0 for datasets (or CC-BY where attribution must be enforced); machine-readable (R1.1).
- [ ] **Provenance documented** — collection + processing workflow; derivations and sources (R1.2).
- [ ] **Open formats** — CSV/TSV/Parquet, plain-text/Markdown docs; not proprietary-only.
- [ ] **Versioned** — version DOI for the cited snapshot; changelog; code + environment archived.
- [ ] **Sensitive data handled** — de-identified; access controls + consent documented before deposit (A1.2).
- [ ] **Metadata persists** — landing page survives even if data are withdrawn (A2).

## 10. Anti-patterns

- **Data dumped without documentation.** A `.zip` of cryptically-named CSVs with no README/dictionary is *not* reusable, however "open." Documentation is the deliverable, not an afterthought.
- **No persistent identifier.** A Dropbox/Drive/personal-website link rots and isn't citable; mint a DOI in a real repository (F1/F4).
- **Proprietary-only formats.** `.sav`/`.dta`/`.xlsx`/vendor binaries exclude anyone without that tool and risk obsolescence; always ship an open companion (CSV/Parquet/plain text).
- **No license / ambiguous "contact me" terms.** Legal uncertainty kills reuse and breaks machine-actionability; state CC0 or CC-BY explicitly.
- **Undocumented missing-data codes.** A bare `-9` or `99` silently corrupts every downstream analysis.
- **"Latest version" citations.** A mutable concept DOI (or bare URL) breaks reproducibility; cite the specific version DOI.
- **FAIRwashing / "open = FAIR".** Posting a file publicly is not FAIR without PID, rich metadata, standards, and a license; conversely, restricted data *can* be FAIR.
- **Metadata only in the paper.** Methods/variable definitions buried in a PDF (not machine-readable, not attached to the data) fail I1 and R1.
- **Deferring documentation to "after the project."** Capture metadata at creation; reconstructing it later is lossy or impossible.
- **Ignoring the domain standard.** Generic Dublin Core where your field has ISA-Tab/DDI/EML/CF loses the discipline-specific richness reusers expect (R1.3).

## References

- Wilkinson, M.D. et al. (2016). *The FAIR Guiding Principles for scientific data management and
  stewardship.* Scientific Data 3:160018. https://www.nature.com/articles/sdata201618
- GO FAIR. *FAIR Principles* (living document + per-principle pages). https://www.go-fair.org/fair-principles/
- DataCite Metadata Working Group (2026). *DataCite Metadata Schema 4.7.*
  https://datacite-metadata-schema.readthedocs.io/en/latest/
- NIH. *Writing a Data Management & Sharing Plan* (2023 policy; 2026 Pilot plan format per NOT-OD-26-046).
  https://grants.nih.gov/policy-and-compliance/policy-topics/sharing-policies/dms/writing-dms-plan
- NSF. *Preparing Your Data Management and Sharing Plan* (PAPPG). https://www.nsf.gov/funding/data-management-plan
- DMPTool. *General Guidance.* https://dmptool.org/general_guidance
- ARDC. *Guide to Choosing a Data Repository.* https://ardc.edu.au/resource/guide-to-choosing-a-data-repository/
- NIH/GREI & FAIRsharing. *Generalist Repository Comparison.* https://fairsharing.org/GeneralRepositoryComparison
- The Turing Way. *Research Data Management.* https://book.the-turing-way.org/reproducible-research/rdm/rdm-repository/
- OpenAIRE. *Research data: how to license.* https://www.openaire.eu/research-data-how-to-license/
- Creative Commons (2026). *Licensing Best Practices for the Sharing of Scientific Data.*
  https://creativecommons.org/2026/04/20/licensing-best-practices-for-the-sharing-of-scientific-data/
- ICPSR / MCW RDM. *Codebooks and Data Dictionaries.* https://mcw.libguides.com/c.php?g=1288089&p=9469367
- Cornell Data Services / UVa Library. *Metadata & Documentation (README).* https://guides.lib.virginia.edu/RDM/metadata-and-documentation
- DCC. *Curation Lifecycle Model.* https://dcc.ac.uk/faq/dcc-curation-lifecycle-model
- UK Data Service. *Research data management / Prepare your data for deposit.* https://ukdataservice.ac.uk/learning-hub/research-data-management/
