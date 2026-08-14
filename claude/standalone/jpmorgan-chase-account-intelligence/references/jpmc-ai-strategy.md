# JPMC — AI Strategy (LLM Suite, CDAO, Value Claims, Use Cases)

> Provenance: reference under the `jpmorgan-chase-account-intelligence` skill. Built 2026-06-18 via `/dr` deep research (exa + WebSearch/WebFetch; primary sources). **verified-as-of: 2026-06-18.** AI user-counts and value claims change **quarterly** — re-verify before customer use.

## Contents
- LLM Suite (internal GenAI assistant)
- The AI / data-science organization (CDAO)
- AI value / ROI claims (dated)
- Key AI use cases
- TAM implications
- Sources

## LLM Suite — the internal generative-AI assistant

**Name and purpose confirmed.** "LLM Suite" is JPMC's proprietary internal GenAI platform, giving employees secure access to large language models for drafting, summarization, document analysis, presentation generation, earnings-transcript analysis, and data insights. Released **summer 2024** (a 2023 pilot is referenced).[^1][^2] It won American Banker's 2025 "Innovation of the Year" Grand Prize and was still actively cited by the Global CIO in mid-2026 — **it is real and active, not discontinued** (negation check passed).[^2][^3]

**User count — the single most volatile figure; always date it.** Trajectory: ~140,000 (mid-2024) → ~200,000 onboarded within ~8 months of the July-2024 debut → "230,000+ globally" (early 2026).[^1][^4] Best framing for the TAM as of late-2025/2026:
- **~250,000 employees have access** (the whole firm except branch/call-center staff).[^5]
- **~150,000 weekly active** (Dimon, Bloomberg TV, Oct 2025).[^6]
- About **half of those with access use it daily.**[^5]

**Model-agnostic / multi-provider — confirmed.** LLM Suite routes to multiple external foundation-model providers, explicitly **OpenAI and Anthropic** (named repeatedly).[^5][^7] **Caution:** other providers (e.g., Google/Gemini, Llama) were **not** named in the sources reviewed — do **not** assert them. The accurate claim is "model-agnostic, routing to multiple external providers including OpenAI and Anthropic." Eight major upgrades since launch; now positioned as an "AI hub" with custom assistants, document analysis, and data viz; ~3–6 hours/week saved per user (self-reported).[^4]

## The AI / data-science organization (CDAO)

- **Chief Data & Analytics Office (CDAO)** — created as a **firmwide org in 2024** to consolidate all data + AI under one umbrella, led by **Teresa Heitsenrether (Chief Data & Analytics Officer)**, who reports to Dimon and sits on the Operating Committee.[^8] (See `jpmc-leadership.md` for the full leadership table.)
- **Derek Waldron — Chief Analytics Officer**, oversees the overall AI program and owns LLM Suite (in role since 2023).[^9]
- **Machine Learning Center of Excellence** — still exists, **inside the CDAO** (confirmed via active JPMC job postings). **Qualified.**[^10]
- **JPMorganChase AI Research** — the named fundamental-research group (data mining, ML, cryptography, explainability, human-AI interaction; ~100+ researchers). **Leadership changed in early 2026 — see `jpmc-leadership.md`.**
- **AI talent scale:** Dimon cited **~2,000 people "working on AI"**; JPMC has ranked **#1 on the Evident AI Index** multiple years running.[^6][^8]

## AI value / ROI claims (dated trajectory)

| When | Claim | Source |
|---|---|---|
| 2023 | Target ~**$1.5B** AI/personalization business value by year-end (raised from $1B) | American Banker[^11] |
| 2024 Investor Day | AI use-case value ~**$1B–$1.5B**; President Pinto raised the bar to ~**$2B** expected | Evident[^12] |
| 2025 Investor Day | **+35% YoY** increase in value from AI/ML (CCB), productivity gains, flat fraud cost despite ~12% CAGR in attacks | Investor Day (primary)[^13] |
| **Oct 2025 (latest)** | Dimon: **~$2B/yr spend producing ~$2B/yr of benefit** — *"for $2B of expense, we have about $2B of benefit"* | Bloomberg TV / eMarketer[^6][^14] |

**Confidence: fact** (multi-source for the trajectory and the latest figure). **Volatility: high.** Use-case scale: ~100 GenAI solutions in production plus "several hundred" total AI use cases (2025), with a stated plan to reach **~1,000 use cases by 2026**.[^13][^15]

## Key AI use cases (for conversation mapping)

- **Fraud detection/prevention** — the oldest use (since ~2012); flat fraud cost despite rising attacks.[^13]
- **Customer personalization** — e.g., a personalized mobile landing page cited as +25% engagement.[^13]
- **CIB / markets** — pricing, hedging, trading analytics.[^13]
- **Credit decisioning** and **operations automation.**[^13]
- **Software engineering** — **40,000+ engineers using AI coding assistants**, with coding-productivity gains cited up to ~20%.[^13][^15]

## TAM implications

- **AI is funded and quantified at the top.** This is not exploratory; it's a board-level program with an ROI discipline. A data story that accelerates "AI-ready data" attaches to a live, funded priority.
- **The bottleneck JPMC names is data, not models.** Its own executives say data quality / fit-for-purpose is the multiyear blocker (see `jpmc-cloud-data-modernization.md`). That is the operational-data conversation.
- **Model-agnostic posture matters.** LLM Suite deliberately routes across providers; JPMC values optionality and interoperability — consistent with its multi-cloud doctrine.
- **Do not over-state which models or how many users.** Quote OpenAI+Anthropic only; date the user count; re-verify.

## Sources

1. American Banker, "How JPMorganChase democratized employee access to gen AI" (2025-05-22) — LLM Suite ~200K users; Waldron quotes. (Press) — https://www.americanbanker.com/news/how-jpmorganchase-democratized-employee-access-to-gen-ai
2. JPMorganChase Technology blog, LLM Suite award post (2025-06-03) — 0→200K in ~8 months; purpose. (Blog, primary) — https://www.jpmorganchase.com/about/technology
3. American Banker, "Chase's Lori Beer is #8 on the Most Innovative…" (2026-06-01) — Beer current CIO; LLM Suite ~250K. (Press) — https://www.americanbanker.com/
4. The Digital Banker, "JPMorgan Chase's LLM Suite drives AI transformation" (2026-03-06) — "230K+ globally"; 8 upgrades; 3–6 hrs/week. (Press) — https://www.thedigitalbanker.com/
5. AI Magazine, "JPMorgan, OpenAI, Anthropic evolving banking operations" (2025-10-03) — ~250K access; OpenAI+Anthropic; ~half daily. (Press) — https://aimagazine.com/
6. Yahoo Finance, "Jamie Dimon declares JPMorgan Chase's…" (2025-10-27) — Dimon "$2B expense / $2B benefit"; ~150K weekly users. (Press) — https://finance.yahoo.com/
7. European Financial Review, "JPMorgan expands AI push with LLM Suite" (2025-10-01) — OpenAI+Anthropic; agents next. (Press) — https://www.europeanfinancialreview.com/
8. MarketsMedia, "JPMorganChase makes data AI-ready" (2025-06-17) — CDAO formed 2024; Heitsenrether; firmwide CDO. (Press) — https://www.marketsmedia.com/
9. McKinsey, "JPMorgan Chase's Derek Waldron on building an AI-first bank culture" (2025-10-29) — Waldron Chief Analytics Officer; LLM Suite; ROI discipline. (Interview, primary-ish) — https://www.mckinsey.com/
10. Built In, "Applied AI ML Director — Machine Learning Center of Excellence" (JPMC job listing) — ML CoE sits in CDAO. (Job posting, primary) — https://builtin.com/
11. American Banker, "JPMorgan Chase aims to create $1.5 billion in value with AI by year-end" (2023-05-30) — $1.5B 2023 target. (Press) — https://www.americanbanker.com/news/jpmorgan-chase-aims-to-create-1-5-billion-in-value-with-ai-by-yearend
12. Evident Insights, "Find me the AI money" (2024-09-19) — $1B–$1.5B (2024 ID); Pinto raised to $2B; 140K LLM Suite. (Press) — https://evidentinsights.com/
13. JPMorgan Chase 2025 Investor Day, full presentation + transcript (PDF) — +35% AI/ML value; fraud; personalization +25%; 40K+ engineers w/ AI assistants. (Investor Day, primary) — https://www.jpmorganchase.com/ir/investor-day
14. eMarketer, "AI arms race: JPMorgan $2 billion" (2025-10-09) — $2B spend/$2B benefit; ~2,000 AI staff. (Press) — https://www.emarketer.com/
15. CNBC, "JPMorgan Chase fully AI-connected megabank" (2025-09-30) — use-case scale; ~1,000 by 2026 plan. (Press) — https://www.cnbc.com/2025/09/30/jpmorgan-chase-fully-ai-connected-megabank.html
