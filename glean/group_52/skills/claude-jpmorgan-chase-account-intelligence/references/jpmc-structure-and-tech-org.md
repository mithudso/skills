# JPMC — Corporate Structure, Lines of Business & How Technology Maps

> Provenance: reference under the `jpmorgan-chase-account-intelligence` skill. Built 2026-06-18 via `/dr` deep research (exa + WebSearch/WebFetch; primary sources). **verified-as-of: 2026-06-18.** PUBLIC-SOURCE intelligence; re-verify volatile facts (segment names, headcount, per-LOB budgets) before customer use.

## Contents
- The firm at a glance
- The reportable segments (CCB, CIB, AWM) + Corporate
- How technology maps to the organization (the matrixed model)
- TAM implications
- Sources

## The firm at a glance

JPMorgan Chase & Co. is the largest US bank by assets — **~$4.0 trillion in assets** and **~$345B stockholders' equity** at 2024 year-end[^1] — and reported **FY2025 revenue of ~$180.6B and net income of ~$57.0B**.[^2] Total workforce is **~318,000** (317,233 at FY2024-end; ~318,512 at FY2025-end; a 2026 interview referenced "319,000").[^3] These figures are point-in-time; treat the workforce number as "~318,000 as of 2026" and re-verify.

## The reportable segments (three) + Corporate

**Important correction to a common misconception:** JPMC reports **three** business segments, **not four**. "Corporate" is the **residual/non-segment bucket**, not a reportable LOB.[^4] Effective **Q2 2024**, JPMC combined the former **Corporate & Investment Bank** and **Commercial Banking** into a single segment, the **Commercial & Investment Bank (CIB)** (announced 2024-01-25).[^5]

| Segment | What it does | Rough scale (FY2024 unless noted) | Segment CEO (as of 2026 — verify) |
|---|---|---|---|
| **Consumer & Community Banking (CCB)** | US consumer & small-business banking, cards, home/auto lending, wealth for everyday customers; the "Chase" franchise | ~84M US consumers, ~7M small businesses; #1 US retail deposit share (~11%); ~32% ROE[^6] | Marianne Lake[^6] |
| **Commercial & Investment Bank (CIB)** | The merged wholesale powerhouse: global investment banking, markets/trading, securities services, global payments, commercial/corporate banking | ~$70B revenue, ~$25B net income, ~18% ROE; #1 in global IB fees; record Payments revenue ~$18B[^6] | Co-CEOs Troy Rohrbaugh & Doug Petno[^5][^6] |
| **Asset & Wealth Management (AWM)** | Asset management + private banking / wealth management | AUM ~$4.0 trillion (client assets higher); ~34% ROE[^6] | Mary Callahan Erdoes[^6] |
| **Corporate** (residual) | Treasury/CIO, firmwide functions, and reconciling items — not a reportable LOB | — | — |

**Volatility flag:** segment leadership shifts. At the January 2024 CIB-formation announcement, Jennifer Piepszak was initially named a CIB co-CEO but subsequently moved to a firmwide **COO** role; the current CIB co-CEOs are Rohrbaugh and Petno.[^5] Re-verify any segment-CEO name before using it.

## How technology maps to the organization — the matrixed model

This is the load-bearing answer for a TAM, and it is **"both centralized AND embedded,"** not one or the other:

**Centralized — a firmwide Global Technology organization.** Technology is run as a worldwide function under the **Global Chief Information Officer (Lori Beer as of 2026)**, who sits on the **~12-member Operating Committee** and is, per her official bio, "responsible for the firm's technology systems and infrastructure worldwide … supporting JPMorganChase's retail, wholesale and asset and wealth management businesses."[^7] Shared platforms, core infrastructure, cybersecurity, the cloud program, the database platform org, and engineering standards live here. A firmwide **Chief Data & Analytics Office** (Teresa Heitsenrether) owns data + AI across the firm (see `jpmc-ai-strategy.md` and `jpmc-leadership.md`).

**Embedded — budget and technologists aligned into the LOBs.** A large share of technology *investment* spend and many technologists sit inside the businesses. In 2022 the Global CIO described the model explicitly: firmwide platforms/infrastructure are centralized, while **~$4 billion was "aligned to the firm's four major lines of business."**[^8] **Note the era:** that "four LOBs" quote predates the **Q2-2024 CIB merger** — at the time the four were CCB, the *former* Corporate & Investment Bank, Commercial Banking, and AWM. After the merger the firm reports **three** segments (CCB, CIB, AWM); the 2025 by-LOB technology-investment breakout below is on that three-segment basis. So "four" (2022) and "three" (2024+) are the same businesses re-grouped, not a discrepancy. The 2025 Investor Day broke technology **investment** out **by LOB: CCB ~$3.2B, CIB ~$3.7B, AWM ~$1.0B**.[^9] CCB alone is described as having its own CIO (Gill Haus) with a multi-billion-dollar budget and ~12,000+ technologists.[^10]

**Net for the TAM:** there is one global platform/standards authority *and* well-funded, semi-autonomous LOB technology orgs. A deal can be sponsored inside a LOB (where the budget and the business problem live) but will be governed by central data/infrastructure/security standards and, for databases, very likely a central database platform team. Know which side of the matrix your champion sits on.

## TAM implications

- **Two buying centers, one rulebook.** Sell the business value into the LOB (CCB/CIB/AWM); satisfy the central platform, data, security, and resiliency standards regardless of which LOB sponsors.
- **The database platform is likely central.** Public job postings describe a central NoSQL/MongoDB-Atlas platform org (see `jpmc-vendor-procurement-and-mongodb-signals.md`) — meaning a database decision may route through a shared-platform team, not only the LOB.
- **Scale is the context for everything.** A "small" JPMC workload is large by most standards; per-LOB tech investment alone runs in the billions.

## Sources

1. JPMorgan Chase & Co. Form 10-K, FY2024 (filed 2025-02), sec.gov — ~$4.0T assets; ~$345B equity; reportable-segment structure. (10-K, primary) — https://www.sec.gov/cgi-bin/browse-edgar?action=getcompany&CIK=0000019617&type=10-K
2. JPMorgan Chase 2025 Annual Report (released April 2026) — FY2025 ~$180.6B revenue / ~$57.0B net income. (Annual report, primary) — https://www.jpmorganchase.com/ir/annual-report
3. Banking Dive, "JPMorgan investor day: CEO Dimon, headcount, AI" (2025-05-19) — ~317,000 employees; CFO "resist headcount growth." (Press) — https://www.bankingdive.com/news/jpmorgan-investor-day-ceo-dimon-headcount-ai-banking/748479/
4. JPMorgan Chase 10-K FY2024 — three reportable segments (CCB, CIB, AWM); Corporate as residual. (10-K, primary) — https://www.sec.gov/cgi-bin/browse-edgar?action=getcompany&CIK=0000019617&type=10-K
5. PYMNTS, "JPMorgan forms expanded Commercial & Investment Bank" (2024-01-25) — Q2-2024 CIB merger; co-CEO leadership; Piepszak→COO. (Press) — https://www.pymnts.com/personnel/2024/jpmorgan-forms-expanded-commercial-investment-bank-makes-management-changes/
6. JPMorgan Chase 2024 Annual Report (PDF) — per-segment scale (CCB ~84M consumers; CIB ~$70B rev; AWM ~$4T AUM) and segment CEOs. (Annual report, primary) — https://www.jpmorganchase.com/ir/annual-report
7. Lori Beer official JPMorganChase bio — Global CIO; Operating Committee; "worldwide … retail, wholesale and asset and wealth management." (Leadership page, primary) — https://www.jpmorganchase.com/about/leadership/lori-beer
8. Banking Dive, "JPMorgan Chase IT: Lori Beer, $14 billion budget, 6,000 apps" (2022) — centralized-plus-LOB model; "~$4B aligned to the four LOBs." (Press) — https://www.bankingdive.com/news/jpmorgan-chase-IT-lori-beer-14-billion-budget-6000-apps/633135/
9. JPMorgan Chase 2025 Investor Day, full presentation (PDF) — technology investment by LOB: CCB $3.2B / CIB $3.7B / AWM $1.0B. (Investor Day deck, primary) — https://www.jpmorganchase.com/ir/investor-day
10. InformationWeek, "AI is driving a return to tech fundamentals, says Chase CIO" (2025-06-18) — CCB/Chase CIO Gill Haus; CCB technology org scale. (Press) — https://www.informationweek.com/it-leadership/ai-is-driving-a-return-to-tech-fundamentals-says-chase-cio
