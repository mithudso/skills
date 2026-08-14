# skills

Clone-and-go Claude skills, commands, and a local semantic skill index.

## Quickstart

```bash
git clone https://github.com/10gen/tam-skills.git ~/.claude/skills
cd ~/.claude/skills && ./setup.sh
```

`setup.sh` is idempotent and safe to re-run. It installs dependencies
(Homebrew/Git/Node), prompts for API keys (written to `~/.claude/.env`, never
committed), lays down `~/.claude` config non-destructively, starts Ollama +
the embedding model, builds and embeds the skill index, bridges the legacy
`~/.claude/skill-consolidation` path to the vendored copy, merges MCP servers
(keyless servers stay disabled), and installs a nightly `git pull` + re-embed
cron job. Flags: `--skip-ollama --skip-service --skip-mcp --skip-keys
--skip-config --skip-cron`.

## Useful commands

- `/cfe <subject>` — Concept Family Explorer (map gaps, loop `/dr`, saturate)
- `/dr <topic>` — deep research that ends in an installed skill
- `/cdo`, `/ddo`, `/pdo`, `/sko`, `/dqo`, `/deso` — the deep-optimizer family
- `node skill-consolidation/gen-skills-index.mjs --search "..."` — semantic search

## Docs

Deeper technical docs live in [`skill-consolidation/`](skill-consolidation/):
`README.md` (index engine) and `INDEX-ARCHITECTURE.md` (full spec).

## Further reading

- [Hub-and-spoke skill strategy](https://docs.google.com/document/d/1z3xVqyaRKNBYBmbm70o8ZAJFyD0fHYhQ-3o7zmGe79E/edit?tab=t.0#heading=h.dau39xa4rjjw) — why skills are organized as hubs with folded spokes instead of one flat list, and when to split/merge
- [How to use: concept-family-explorer, deep-research, skill-tree-architect, and the deep-optimizer family](https://docs.google.com/document/d/1J1n4xZFbynAq8MgF7dAq4vc6_nIG9DRnn9Xz60Q7fLo/edit?tab=t.0#heading=h.382lyy7dd8vs) — invocation, before/afters, and a routing cheatsheet for the meta-tooling
- [Skill-Library Indexing System — Technical Architecture](https://docs.google.com/document/d/10ZTBXtnU-eQ-pGQysQQ3AG7Hfn55HJLhcJrY3Uqe1_s/edit?tab=t.0#heading=h.qiwzxf642ge) — how `SKILLS-INDEX.{json,md}` and semantic search are built and kept fresh
- [Skills Ecosystem — Showcase & Field Guide](https://docs.google.com/document/d/1PedNawIj06dgVUfa5Kt96KZQXWM7kB9ZuZvm1_9BQn8/edit?usp=sharing) — a tour of what's in here and how it's meant to be used day to day
- [Skills Ecosystem — Showcase & Field Guide (additional section)](https://docs.google.com/document/d/1PedNawIj06dgVUfa5Kt96KZQXWM7kB9ZuZvm1_9BQn8/edit?tab=t.5v4qadbd4ux3) — same doc, different tab

## Skill families

125 top-level skills across 30 families (plus 455 topics folded into hub `references/`
files — the same content, just filed under the hub that owns it instead of getting its own
top-level directory). This section covers what each family/hub is for; **Standalone
skills** below covers everything that isn't part of a consolidated family.

Machine-readable version of this same structure: `skill-consolidation/SKILLS-INDEX.md`
(auto-generated, includes trigger phrases and cross-references — regenerate it with `node
skill-consolidation/gen-skills-index.mjs` after adding or editing skills).

### ai-agent
Agent frameworks, orchestration, and the full LLM stack beneath them. Router:
`ai-agent-engineering` (standalone, below). Hubs:
- `ai-agents-orchestration` — agent frameworks (Claude SDK, LangGraph, CrewAI), multi-agent orchestration, agent memory/planning/reliability, coding & GUI agents, autonomous loops, eval-driven development.
- `ai-llm-model-layer` — pretraining/scaling laws, fine-tuning/PEFT, RLHF/alignment, compression, inference serving, transformer/multimodal architecture, model selection, observability.
- `ai-mcp-sdk-prompting` — MCP servers & the MCP builder, the Anthropic SDK, prompt engineering, context engineering, tool-search optimization, iterative self-refinement loops.
- `ai-rag-retrieval` — RAG architecture (chunking, embeddings, hybrid search, reranking), advanced RAG patterns, vector/graph datastores.

### aws-cloud
- `aws-cloud` — AWS as a whole: IAM/EC2/S3/VPC (core), Lambda/Step Functions/CDK (serverless), Bedrock/SageMaker (AI/ML), PrivateLink/VPC endpoints, and AWS databases (RDS/Aurora/DynamoDB).

### blockchain
Chain internals, economics, and contract security. Hub:
- `blockchain` — how chains work beneath the app layer; routes to the spokes below.
- `bitcoin-protocol-expert` — UTXO model, Script/addresses, PoW, Lightning, ordinals/Runes.
- `blockchain-economics` — tokenomics, cryptoeconomics, fee markets, MEV.
- `ethereum-protocol-expert` — account model, EVM, EIP-1559, Proof of Stake/Gasper, rollups.
- `smart-contract-security` — reentrancy, access control, integer over/underflow, auditing EVM and Solana contracts.

### chrome-extension
- `chrome-extension-expert` — MV3 build/package/secure/test: service-worker lifecycle, manifest/permissions, cross-context messaging, native messaging, OAuth, storage, offscreen documents (17 references folded in).

### claude-code
- `claude-code-skills` — everything about Claude Code skills themselves: SKILL.md anatomy, authoring, discovery, distribution (plugins/marketplaces), management, and Claude Code workflow patterns (plan mode, worktrees, headless/CI).

### consumer-credit-and-debt
- `consumer-credit-and-debt` — US consumer credit/debt hub (credit reports, FDCPA, bankruptcy, mortgages, predatory lending, identity theft, NC-specific credit law; 11 folded spokes). Educational, not legal advice.

### consumer-finance
- `consumer-finance` — sibling hub for US personal finance (student loans, banking, income taxes, insurance, medical bills, budgeting, investing/retirement, estate planning; 9 folded spokes). Educational only.

### content-ingestion-extraction
- `content-ingestion-extraction` — acquiring and restructuring content from messy sources: doc archaeology, doc-store bootstrapping, Chrome DOM extraction, Granola/Plaud transcript ingestion.

### data-analytics
The full data-analysis stack, foundations through platform. Router: `data-analytics`
(standalone, below). Hubs:
- `da-1-foundations-theory` — measurement, probability, statistical inference, information theory, epistemology.
- `da-2-data-analysis-lifecycle` — process frameworks (CRISP-DM, KDD), problem framing, iteration, documentation.
- `da-3-data-acquisition-sampling` — data-source taxonomy, web scraping/APIs, survey/sampling design.
- `da-analytical-methods` — regression/ML, A/B testing and causal inference, forecasting, anomaly detection, and 12 more sub-techniques.
- `da-applied-and-communication` — visualization, dashboards, data ethics, product analytics, NLP, pricing analytics.
- `da-data-engineering-platform` — dbt, Spark, Airflow, DuckDB, lakehouse architecture, data governance.

### deep-research
- `deep-research` — multi-source cited research reports via firecrawl and exa (fallback: WebSearch/WebFetch); the `/dr` command.

### devops
The DevOps/infrastructure stack, kernel through CI. Router: `devops-infra` (standalone,
below). Hubs:
- `devops-containers-cicd` — Docker, Kubernetes, CI/CD pipeline design, Terraform/Kafka IaC, git workflows, library packaging.
- `devops-linux-admin` — sysadmin troubleshooting, systemd, package management, shell scripting, host/macOS network diagnostics.
- `devops-linux-internals` — kernel architecture, boot/init, memory/NUMA, storage, virtualization, cgroups/namespaces, sandboxing.
- `devops-observability` — Node/OTel instrumentation, Pino structured logging, Sentry, eBPF, Linux perf tracing.

### document-formats
- `document-formats` — generate/parse/convert PDF, DOCX, XLSX, PPTX, CSV, advanced JSON, draw.io, and Markdown in Python and Node.

### frontend
- `frontend-ui` — frontend/UI-UX design and build: design systems, HTML/CSS, accessibility (WCAG/ARIA), mobile/SwiftUI, usability heuristics, visual-design critique.

### integration
- `integration-clients` — third-party SaaS API clients: Jira, Monday.com, Slack, Salesforce, Glean, Aha! (REST/GraphQL, OAuth, webhooks, SDKs/CLIs).

### lang
Per-ecosystem language idioms and tooling. `lang-rust` is closely related but currently
files as a standalone skill (see below) rather than under this family. Hubs:
- `lang-go-and-mobile` — Go idioms (goroutines, channels, Cobra CLIs) and Kotlin/Compose Multiplatform.
- `lang-js-ts` — JavaScript/Node and TypeScript, non-Node runtimes (Deno/Bun/edge), Node concurrency internals, V8 internals.
- `lang-python` — Python 3.12+ idioms, async/structured concurrency, CLI/TUI apps, metaprogramming, pydantic v2, packaging.

### misc-catch-all
- `misc-catch-all` — grab-bag hub for topics too narrow to earn their own hub: 10gen, mongodb-kb, webcrypto-vault-reviewer, diagrams-as-code, skill-lookup, wordcloud generation, and more (17 folded spokes).

### mongodb
- `mongodb-expert` — data-plane and engine: CRUD/MQL, aggregation pipelines, index design, query optimization, WiredTiger internals, driver internals.
- `mongodb-atlas-expert` — Atlas platform: control plane, Admin API v2, Atlas CLI/Terraform/Kubernetes Operator, tiers, security, networking (27 references).
- `mongodb-operations-expert` — reliability and data movement: backup/DR, upgrades, migration patterns, encryption, compliance, cost optimization (18 references).
- `atlas-diagnostics-expert` — live diagnostics, performance, monitoring, and capacity planning (ts-diag, FTDC, benchmarking).

### networking
- `networking` — protocol-level DNS (resolution path, record types, DNSSEC, encrypted DNS) and network diagnostics; cloud-connectivity spokes are queued for a future pass.

### optimizers
The convergence-loop optimizer family: each hub runs a multi-pass audit on one artifact
type, applies every Medium+ fix, and verifies empirically before looping to convergence.
Full usage guide is in the ["How to use" doc](#further-reading) above. Router: `deep-optimizer`
(standalone, below); prose sibling `document-critique`/`ddo` files under the `writing`
family. Hubs:
- `code-deep-optimizer` (`/cdo`-style invocation) — source files and repos: an 18-pass audit (correctness, security, performance, tests) with a build/lint/test verify gate.
- `deep-query-optimizer` (`/dqo`) — SQL queries: sargability, index design, joins, pagination; verifies via `EXPLAIN`/`EXPLAIN ANALYZE` when a connection exists.
- `design-deep-optimizer` (`/deso`-style invocation) — UI/UX screens and mockups: an 11-pass critique (hierarchy, contrast, WCAG, brand-parity), with a re-render verify gate for code-backed designs.
- `prompt-deep-optimizer` (`/pdo`) — production prompts shipped in code: a 16-pass audit with an intent-preservation check and algorithm recommendation.
- `skill-optimizer` (`/sko`) — `SKILL.md` files: trigger-accuracy eval, cross-skill collision check, length budget, hub sync.

### physical-access-control
- `physical-access-control` — badges, readers, controllers, and the credential/RF stack beneath them: RFID/NFC, contactless smart cards, Wiegand/OSDP, mobile credentials.

### psychology
- `applied-psychology` — the operator/TAM-facing psychology hub: decision-making, trust, persuasion, learning, human-AI interaction, performance, personality. (The deeper clinical/social/developmental spokes file as standalone skills below rather than folding into this hub.)

### research-methodology
- `research-methodology` — academic/scientific research method: systematic reviews (PRISMA/GRADE), data annotation, FAIR data management, personal knowledge management, AI-assisted literature-review tools.

### security-review
- `security-review` — application and infrastructure security review: OWASP/ASVS web-app review, compliance auditing (secrets/supply-chain/PII), security headers, and the Okta platform security spokes (identity governance, ITDR, phishing-resistant auth, zero-trust device).

### software
- `software-engineering-patterns` — language-agnostic design/architecture/practice: API design, backend/microservices patterns, coding standards, code review, debugging, performance profiling, testing.

### tam-operations
- `tam-operations` — MongoDB TAM account operations: EBR/QBR deliverables, account health scoring, case/incident tracking, reporting automation, Monday.com/Slack ops.

### technical-instruction
- `technical-instruction` — teaching engineers: course/curriculum design, certification design, hands-on labs, teaching troubleshooting and diagnostic reasoning (10 folded spokes).

### trading-and-investing
- `trading-and-investing` — how US-retail markets work (asset classes, order routing, market participants); routes to 16 spokes. Educational only, not financial advice.
- `options-trading-and-strategies` — listed equity/ETF/index options for a US retail participant: the Greeks, strategies, IV. Educational only.

### venture
- `venture-business` — NC for-profit solo/lean founder: entity formation, business tax, payroll, real estate, marketing, funding.
- `venture-nonprofit-cause` — NC nonprofit/cause founder: 501(c)(3) formation, fundraising ops, cause marketing, donor engagement.

### writing
The largest family: every writing genre, plus the tools that critique and clean up prose.
- `writing-expert` — general prose craft, voice, and editing; frameworks (BLUF, Minto Pyramid, SCQA); anti-AI-ism; 18 sub-skills on demand.
- `career-and-formal-writing` — resumes/CVs, cover letters, academic prose, legal-adjacent drafting, survey design.
- `content-and-marketing-writing` — landing pages, press releases, newsletters, chatbot writing, support-ticket replies.
- `executive-comms` — board memos, exec summaries, pitch decks, negotiation prep, public speaking.
- `technical-writing-craft` — API docs, READMEs, runbooks, RFCs, PRDs, commit/PR messages, postmortems.
- `document-critique` (`/ddo`) — multipass document review with a convergence loop; see the optimizers section above.
- `ddo` — the `/ddo` command wrapper for `document-critique`.
- `kill-the-AI-ism` — detects and replaces generator artifacts ("AI-isms") in prose.
- `resume-and-cv-writing` — the X-Y-Z achievement-bullet formula, ATS-safe formatting, role tailoring.

## Standalone skills

Skills that either stand entirely on their own or are the top-level router for a family
listed above (routers are noted). Alphabetical; one line each.

| Skill | What it's for |
|---|---|
| `ai-agent-engineering` | Router → `ai-agent` family (agent orchestration, RAG, model layer, MCP/SDK/prompting) |
| `ai-assisted-copywriting-workflow` | Using LLMs for on-brand marketing copy: brief→draft→QA loop, variation generation, FTC compliance |
| `artifact-to-procedure` | Distills a completed artifact or solved problem into a reusable skill or playbook |
| `atlas-vector-search-pii-isolation` | PII isolation and entitlement boundaries in MongoDB Atlas Vector Search/RAG deployments |
| `automated-job-applications` | Automated/AI-assisted job-application tools and strategy (auto-apply bots, ATS optimization) |
| `bank-genai-model-risk-governance` | GenAI/LLM model-risk governance for banks (SR 11-7, NIST AI RMF, EU AI Act) |
| `bartending-career-track` | Bartending career progression, certifications, competitions, portfolio building |
| `big-bank-IT` | A large bank's IT estate: core banking, mainframe modernization, cloud posture, resiliency |
| `caveman` | Ultra-compressed communication mode (~75% token reduction, technical accuracy preserved) |
| `cloudflare-platform` | Cloudflare from the infrastructure angle: reverse-proxy DNS/TLS effects, Tunnel, Access, caching |
| `concept-family-explorer` | Gap-discovery layer above `/dr`: maps a subject's conceptual family and loops `/dr` over the gaps until saturated |
| `conversion-copywriting-and-voice-of-customer` | Mines customer language (reviews, tickets, transcripts) for evidence-based copy strategy |
| `data-analytics` | Router → `data-analytics` family (the six `da-*` hubs) |
| `database-proxies-query-middleware` | Database proxies and query-optimization middleware: pooling, routing, caching, firewalling |
| `deep-optimizer` | Router → `optimizers` family (code/query/design/prompt/skill deep-optimizers) |
| `devops-infra` | Router → `devops` family (linux internals/admin, containers/CI-CD, observability) |
| `diagnosis-methodology-backtest` | Blind, parallel, multi-agent backtest of competing diagnosis methodologies against ground truth |
| `direct-response-and-sales-letter-copywriting` | Long-form direct-response canon: market sophistication, sales-letter anatomy, VSL scripting |
| `distributed-systems-consensus` | Distributed-systems and consensus theory (CAP/FLP, Paxos/Raft, BFT) plus blockchain consensus |
| `effective-altruism-and-philanthropic-decision` | Evidence-based decisions about where to give and how to maximize philanthropic impact |
| `elasticsearch-opensearch` | Elasticsearch/OpenSearch and the ELK stack: engine internals, Query DSL, ILM, the 2021 fork |
| `embedding-inversion-threat-model` | Threat model for data leakage from vector embeddings; why embeddings aren't safe anonymization |
| `enterprise-vendor-management-and-tprm` | Buyer-side vendor governance and third-party risk, for a seller selling into large enterprises |
| `firecrawl-onboarding` | Onboarding/install router for Firecrawl (API keys, CLI install, usage-path routing) |
| `fsi-banking-regulatory-context` | FSI regulatory fluency for a B2B tech seller: how banking segments buy and why |
| `fundraising-and-donor-psychology` | Fundraising practice psychology: donor retention, the ask, recurring giving, stewardship |
| `generative-engine-optimization` | Structuring content to get cited by AI answer engines (AI Overviews, ChatGPT search, Perplexity) |
| `goldman-sachs-account-intelligence` | Public-source account intelligence on Goldman Sachs as a technology buyer |
| `harness-streamliner` | Audits/streamlines the Claude Code harness config to cut startup and token overhead |
| `health-behavior-change-and-donor-registration` | Health-behavior models (HBM, TPB, EPPM, COM-B) applied to organ-donor registration campaigns |
| `jpmorgan-chase-account-intelligence` | Public-source account intelligence on JPMorgan Chase as a technology buyer |
| `kql-kusto-query-language` | Kusto Query Language across Azure Data Explorer, Log Analytics, Sentinel, Defender XDR, Fabric |
| `lang-rust` | Rust language sub-hub: ownership/borrowing, traits/generics, async/tokio, Rust-for-blockchain |
| `legal` | Hub for NC criminal law, drug law, parole/probation, bail/bond, recording laws; routes to spokes |
| `logql-grafana-loki` | LogQL, the query language for Grafana Loki (label-indexed log/metric queries) |
| `nc-criminal-defense` | NC criminal law: bail/bond, drug law, DWI, juvenile justice, sentencing, expungement |
| `nc-drug-rehabilitation-treatment` | NC substance-use-disorder treatment law: harm reduction, MAT, drug courts, confidentiality |
| `nc-parole-probation` | NC parole, probation, and community-supervision law |
| `nc-recording-laws` | NC recording law: one-party consent, recording police, workplace/public-meeting recording |
| `okta-incident-response` | Guides MongoDB responders through Okta Tier-0 incident response and fire drills |
| `postgresql-expert` | PostgreSQL engine internals: MVCC, query planner, index design, WAL/replication, partitioning |
| `programming-languages` | Router → `lang` family (Python, JS/TS, Go/mobile) |
| `prompt-helper-optimizer` | One-off/exploratory prompt improvement (the `/ph` and `/phe` commands) |
| `psychology` | Psychology domain hub: clinical, developmental, social, and positive psychology |
| `psychology-clinical-personality` | Clinical/diagnostic depth on Cluster B personality pathology |
| `psychology-confidence-identity` | Confidence, self, and identity psychology — replication-honest science |
| `psychology-influence-depth` | Replication-honest science of influence, compliance, and resistance |
| `psychology-neurodevelopmental` | Deep neurodevelopmental psychology: ADHD and autism as research domains |
| `psychology-of-charitable-giving` | The empirical science of why people give (warm-glow, identifiable-victim effect, compassion fade) |
| `psychology-positive` | Positive psychology and wellbeing science |
| `psychology-social` | Social psychology: how groups, norms, and identity shape behavior |
| `psychology-stress-trauma` | Stress and trauma neuroscience |
| `rag-prompt-injection-defense` | Defending RAG/tool-using LLM apps against prompt injection (the "lethal trifecta") |
| `repo-bootstrapper` | Creates/upgrades a repo to the full mdb-tam standard (CLAUDE.md, docs suite, CI, ops infra) |
| `secrets-and-key-management` | The generate→store→transfer→use→rotate→destroy lifecycle for secrets; envelope encryption |
| `service-industry-job-hunting` | Job hunting for service-industry workers (bartenders, servers, hotel/event staff) |
| `service-industry-resume-and-interview` | Resume writing, formatting, and interview prep for bar/front-of-house roles |
| `skill-tree-architect` | Whole-tree architect for this repo's hub-and-spoke taxonomy; audits placement and description caps |
| `solve-case` | End-to-end MongoDB/Atlas support-case solver orchestrating the mongodb-* and 10gen skills |
| `splunk-platform-spl` | Splunk platform and SPL: architecture, search performance, knowledge objects |
| `telemetry-pipeline` | Vendor-neutral telemetry pipelines (Cribl, OTel Collector, Vector, Fluent Bit) between sources and destinations |
| `volunteer-and-prosocial-motivation` | The motivational science of volunteering and ambassadorship |

## Housekeeping note

A handful of top-level entries in this directory aren't skills and are excluded from the
counts above: `claude-config`, `health`, `loop-factory`, and `skill-consolidation` are plain
directories with no `SKILL.md`; `README.md`, `setup.sh`, `health.md`, and
`okta-incident-response.zip` are loose files. `backprop`, `build`, `check`, `context-canary`,
`deepen`, `grill`, `interface-kit`, `junior-to-senior`, `research`, `review`, `slop-detector`,
and `spec` are symlinks to `../../.agents/skills/<name>` that don't resolve on this machine
(the target directory isn't present) — worth a cleanup pass if those are meant to be live.
