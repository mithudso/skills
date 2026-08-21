# GS — Cloud Posture & Data Infrastructure

Public-source. ✅ fact · 🔶 inference · ⚠️ verify first.

## Polycloud strategy ✅

- Intentional **"polycloud"**: AWS, Google Cloud, Azure, and Oracle Cloud (OCI) under a single **orchestration layer**; avoids "lift and shift."
  - Source: Business Insider (2022-06).
- **Still heavily private-cloud-based "at balance."** CDO Neema Raphael: "we are at balance still heavily private cloud-based," with public-cloud-native pockets like Transaction Banking.
  - Source: computerweekly.com (CDO interview).

## AWS = deepest partner ✅

- **GS Financial Cloud for Data** (re:Invent, 2021-11) — built on AWS; sells GS data/analytics tooling to buy-side clients. Source: goldmansachs.com (2021-11-30).
- **Transaction Banking (TxB)** — built entirely on AWS, cloud-native, API-first. Source: americanbanker.com.
- **AWM Atlas Platform** runs on AWS or OCI; job posts cite multi-cloud security across AWS/GCP/Azure. (GS's internal "Atlas Platform" is a wealth-management platform name — NOT MongoDB Atlas; do not confuse the two.)

🔶 **MongoDB Atlas multi-cloud + cross-cloud clusters map directly to the polycloud thesis.** "Still heavily private cloud" → hybrid/self-managed MongoDB and on-prem→Atlas migration stories are both live. Atlas on AWS aligns with the deepest cloud relationship.

## Proprietary / internal data stack (know it; don't target it) ✅

- **SecDB** (Securities DataBase, 1993) + **Slang** (Securities Language, ~Python-like) — proprietary OO risk/pricing engine; **15–40M lines of Slang** (2017); still the primary risk system. Source: goldmansachs.com history; Wikipedia (Dubno); InfoQ.
- **Legend** — flagship data-modeling/governance platform, **open-sourced via FINOS (2020-10)**, uses the PURE language; Google Cloud integrations. Source: finos.org/legend.
- **GS Quant** — Python quant toolkit atop the risk-transfer platform.

## Data leadership / strategy

- ⚠️ **Neema Raphael — Chief Data Officer & Head of Data Engineering** (Partner). Pillars: platform engineering, content curation, governance. Mantra: "no AI strategy without a data strategy." Source: Business Insider (2024-03).
- Goal: a **multi-public-cloud "single data infrastructure backbone for analytics."** Source: computerweekly.com.

🔶 **Don't displace SecDB/Slang/Legend** — entrenched institutional core. Position MongoDB around **new cloud-native AWM/microservice workloads** and **AI/RAG data foundations**, not risk/pricing.
