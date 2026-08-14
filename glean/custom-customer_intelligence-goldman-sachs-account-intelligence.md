# goldman-sachs-account-intelligence

**Category:** Science, Biology & Medicine
**Platform:** Custom / Customer Intelligence
**Original Path:** custom/customer-intelligence/goldman-sachs-account-intelligence

## Description
Company-specific account intelligence on Goldman Sachs (GS) as a TECHNOLOGY BUYER, for a MongoDB TAM/seller. PUBLIC-SOURCE intel; verify before customer use (leadership/budgets/org-names decay fast). TRIGGER: prepping for a GS touchpoint, EBR, or account plan; GS segments (Global Banking & Markets, Asset & Wealth Management, Platform Solutions), tech spend/engineering scale, polycloud posture, AWM modernization, AI strategy (GS AI Assistant, Devin, OneGS), proprietary stack (SecDB/Slang, Legend/FINOS), tech leadership, consumer exit, MongoDB/Atlas-vs-DocumentDB signals, earnings signals. SKIP: FSI regulation (Basel, BCBS 239) -> fsi-banking-regulatory-context; bank GenAI model-risk governance -> bank-genai-model-risk-governance; big-bank IT archetype -> big-bank-IT; TPRM/procurement -> enterprise-vendor-management-and-tprm; MongoDB/Atlas config -> mongodb-* hubs; EBR/QBR/POV mechanics -> tam-operations; JPMorgan/Chase -> jpmorgan-chase-account-intelligence.

---

# Goldman Sachs — Account Intelligence (MongoDB TAM lens)

Company-specific profile of **Goldman Sachs as a technology buyer**, for a MongoDB
TAM/seller. **Public-source only** (earnings, engineering blogs, press, job
postings, conference talks). Leadership titles, headcounts, org structure, and
the Atlas-vs-DocumentDB footprint move fast — **re-verify every volatile item
with the account team before any customer-facing use.**

Confidence markers used below: ✅ well-sourced · 🔶 inference · ⚠️ verify first.

## Bottom line for the TAM (read this first)

1. **Primary play: AWM modernization** — Custody, Asset Transfers, Portfolio
   Management on Java/Spring Boot microservices already using MongoDB. 🔶
2. **Competitive reality: Atlas vs AWS DocumentDB is being decided internally.**
   GS job posts say "MongoDB clusters (Atlas or DocumentDB)" — differentiate on
   multi-cloud, sharding/indexing depth, Vector Search, and ops tooling. ⚠️
3. **Polycloud + a "single data backbone"** strategy maps directly to Atlas
   multi-cloud / cross-cloud clusters. ✅
4. **Do NOT target SecDB/Slang/Legend** — entrenched institutional core. Target
   **new cloud-native AWM/microservice workloads** and **AI/RAG data foundations**. 🔶
5. **Procurement is rigorous** — bring the full security/compliance + Terraform +
   private-networking story up front. ✅

## How GS makes technology decisions

- A large central **Core Engineering** platform group (3,000+ ⚠️) builds firmwide
  common platforms, so GS frequently **builds vs buys** and standardizes through
  platform teams. Land-and-expand means navigating those platform standards. 🔶
- **Security-as-code / GitOps** culture (policy-as-code, mTLS + IP allowlisting on
  Transaction Banking APIs) favors **Terraform-provisioned Atlas** — already
  referenced in their own job postings. ✅
- "**No AI strategy without a data strategy**" (CDO Neema Raphael ⚠️) is a
  ready-made entry narrative for data-foundation conversations.

## References

Detailed, sourced sections live in `references/`:

- `gs-structure-and-tech-org.md` — segments, consumer exit, engineering scale/culture.
- `gs-cloud-and-data.md` — polycloud, GS Financial Cloud, proprietary stack (SecDB/Slang, Legend).
- `gs-ai-strategy.md` — GS AI Assistant, Devin pilot, OneGS, agentic/RAG appetite.
- `gs-mongodb-signals.md` — confirmed MongoDB job-post signals, Atlas-vs-DocumentDB posture, the MongoDB sales map.
- `gs-leadership-and-procurement.md` — leadership map (verify titles), third-party-risk posture, earnings/tech-spend signals.

All claims are public-source as of the dates cited in the references; nothing
here reflects confidential Goldman Sachs information.