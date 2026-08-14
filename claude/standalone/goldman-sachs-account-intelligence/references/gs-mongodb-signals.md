# GS — MongoDB Signals & Sales Map

Public-source. ✅ fact · 🔶 inference · ⚠️ verify first.

## Confirmed public MongoDB usage signals ✅

From Goldman Sachs **engineering job postings** (verifiable public artifacts):

- **AWM Engineering** roles list **MongoDB** in the stack alongside Java/Spring Boot,
  Kafka, React, AWS/OCI. Workloads named: **Custody, Asset Transfers, Portfolio
  Management.**
- One posting specifies maintaining **"MongoDB clusters (Atlas or DocumentDB)"** —
  the single most important competitive signal in this profile.
  - Source: Goldman Sachs careers / aggregators (Built In, eFinancialCareers), 2025–2026.

⚠️ Job posts prove *a* MongoDB footprint and intent, not deployment size, spend, or
that Atlas (vs DocumentDB) won. Treat as directional. No public MongoDB customer case
study for GS was found.

## The competitive crux: Atlas vs AWS DocumentDB ⚠️

GS's deep AWS relationship makes **DocumentDB** the default-path competitor. When GS
says "Atlas or DocumentDB," differentiate Atlas on:

- **Multi-cloud / cross-cloud clusters** (DocumentDB is AWS-only) — fits polycloud.
- **Sharding, advanced indexing, aggregation depth, and version currency** (DocumentDB
  emulates an older MongoDB API surface).
- **Atlas Vector Search + native Voyage embedding/rerank** for the AI/RAG agenda.
- **Operational tooling** (Performance Advisor, online archive, Terraform/Atlas Operator).

## Sales map 🔶

| Target | Why | Skill to pair |
|---|---|---|
| **AWM modernization** (Custody, Asset Transfers, PM) | Confirmed Java/Spring + MongoDB | mongodb-atlas-expert |
| **AI/RAG data foundation** | GS AI Assistant, Devin, OneGS, ~1M prompts/mo | atlas-vector-search-pii-isolation, rag-prompt-injection-defense, bank-genai-model-risk-governance |
| **Polycloud data backbone** | "single data infrastructure backbone" across clouds | mongodb-operations-expert (multi-cloud) |

**Avoid:** SecDB/Slang/Legend/GS Quant (entrenched proprietary core).
