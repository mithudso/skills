<!-- hub-reference-banner -->
> **Reference file — part of the `research-methodology` hub.** Formerly the standalone `ai-assisted-literature-review` skill.
> Sibling topics in this family are now reference files under the hubs (`research-methodology`, `deep-research`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: ai-assisted-literature-review
title: "AI-Assisted Literature Review & Research-Tool Landscape"
description: >-
  AI research-tool landscape and the provenance/trust layer above hands-on deep research. TRIGGER: choosing which AI research tool for a job (Elicit, Consensus, scite, SciSpace, Undermind, Semantic Scholar/S2ORC, ResearchRabbit/Connected Papers, PaperQA2, Ai2 Asta, OpenAI/Gemini/Perplexity DR); building evidence/extraction tables; AI-assisted review methodology and its limits (hallucinated citations, shallow recall, recency gaps, screening tail-risk); citation grounding/faithfulness in RAG (ALCE, attribution, NLI entailment); claim-to-source traceability, content credentials/C2PA, auditable/reproducible AI synthesis. SKIP: executing web research w/ firecrawl+exa -> deep-research; the HOW of research thinking -> deep-research-methods; the Firecrawl paper-finding workflow -> firecrawl-research-papers; eval-driven LLM-app dev -> eval-driven-development; formal evidence-synthesis methodology + research integrity (PRISMA/GRADE) -> systematic-review-research-integrity.
category: custom
version: 2
updated: "2026-06-15"
related_skills:
  - deep-research-methods
  - systematic-review-research-integrity
  - deep-research
  - firecrawl-research-papers
  - ai-agents-orchestration
whenToUse:
  - "choosing which AI research tool for a job (Elicit, Consensus, scite, Undermind, SciSpace, PaperQA2)"
  - "building evidence/extraction tables across many papers"
  - "AI-assisted review methodology and its limits (hallucinated citations, recall, recency)"
  - "citation grounding/faithfulness in RAG (ALCE, NLI entailment, correctness vs faithfulness)"
  - "claim-to-source traceability, content credentials/C2PA, auditable/reproducible AI synthesis"
keywords:
  - AI literature review
  - Elicit
  - Consensus
  - scite
  - Undermind
  - Semantic Scholar
  - deep research agent
  - citation faithfulness
  - ALCE
  - C2PA
  - research provenance
tags:
  - research
  - ai-tools
  - provenance
  - deep-research
---

<!-- Provenance: spoke under deep-research-methods. Created via /dr deep-research (cluster D). -->

# AI-Assisted Literature Review & Research-Tool Landscape

This skill is the **tooling-and-trust frontier** that sits ABOVE the hub's execution skills.
It does NOT execute research (→ `deep-research`), does NOT teach research thinking
(→ `deep-research-methods`), and is NOT the Firecrawl paper-finding workflow
(→ `firecrawl-research-papers`). It answers two questions: **which AI research tool for which
job**, and **how do you trust and reproduce an AI-assisted synthesis**.

## When this skill applies

Use it when you are:
- Picking a tool from the AI-research stack for a specific job (discover vs synthesize vs verify vs extract).
- Building an evidence or extraction table across many papers.
- Designing an AI-assisted review and reasoning about where AI helps vs silently fails.
- Judging whether an AI-synthesized claim is actually grounded in its cited source.
- Standing up an auditable, reproducible research trail so a reviewer can trust the output.

Route away when: you just need to *run* web research now (`deep-research`); you want the
cognitive method of decomposition/synthesis (`deep-research-methods`); you want Firecrawl's
managed paper-finding flow (`firecrawl-research-papers`); you are build-time evaluating an
LLM app (`eval-driven-development`); or you need formal evidence-synthesis methodology and
research-integrity rigor — PRISMA/PROSPERO/GRADE/RoB, meta-analysis, QRPs
(`systematic-review-research-integrity`). This skill is the AI-research-*tool* landscape and
the provenance/trust of AI synthesis; `systematic-review-research-integrity` is the
formal-methodology sibling for trusting a *body of evidence*.

## 1. Core Concepts

**The stack has three jobs, not one.** Treat "AI research tools" as three distinct layers;
most failures come from using a tool for the wrong layer.

1. **Discover** — find the right papers (citation-graph and semantic-search tools).
2. **Synthesize** — read across many papers and answer a question with citations
   (deep-research agents and academic QA systems).
3. **Verify** — check that a specific claim is true, supported, and not refuted
   (smart-citation and faithfulness tools).

**The data backbone is shared.** Most academic tools sit on the **Semantic Scholar Academic
Graph (S2AG)** (~214M papers, ~2.49B citations) and its full-text corpus **S2ORC**, or on
**OpenAlex**. Per-tool counts differ (Elicit reports ~125–138M, SciSpace ~280M) because each
filters or augments the backbone differently — the index size is not the blind spot, the
*corpus coverage* is. Tools that share a backbone share its gaps: STEM/biomed/CS are well
covered; humanities, law, non-English, and very recent (<6-month, low-citation) work are
under-represented. "Different tool, same gap."

**Grounding ≠ truth, and correctness ≠ faithfulness.**
- *Grounded* = the answer is built only from attached sources, with in-text citations, and the
  system abstains when sources are insufficient.
- *Citation correctness* = the cited source entails the statement (an NLI check).
- *Citation faithfulness* = the model actually **relied** on that source to produce the claim.
  A citation can be correct yet unfaithful (decorative post-hoc citation). Auditing only
  correctness over-credits systems (Wallat et al., 2025).

**Provenance is a spectrum.** From weakest to strongest: inline links → claim→passage
entailment → semantic provenance graph (claim↔evidence with relation type) → cryptographically
signed content credentials (C2PA). Match the rung to the stakes: internal/exploratory →
inline links; a shareable synthesis → claim→passage entailment (§4); audit-grade or published
output → semantic graph (§7) plus C2PA (§6).

## 2. The Tool-by-Tool Landscape

### Discovery & citation-graph tools
- **Semantic Scholar** — the discovery substrate. AI-relevance ranking + **TLDR** one-line
  abstract summaries; free API and bulk datasets. Best for fast scanning, and as the corpus
  under everything else.
- **Connected Papers** — one seed paper → one similarity graph (co-citation + bibliographic
  coupling, *not* a strict citation map). Best for rapidly seeing the shape/neighborhood of a
  field from a known paper. Session-based; free tier ~5 graphs/month.
- **ResearchRabbit** — collection-based ("Spotify for papers"); builds evolving citation
  networks around a growing library, syncs with Zotero, emails on new citing work. Free. Best
  for ongoing/longitudinal lit reviews and author tracking.
- *Connected Papers vs ResearchRabbit*: session vs collection. Many researchers run Scholar
  (search) + Connected Papers (map) + ResearchRabbit (monitor) together.

### Synthesis & academic-QA tools
- **Elicit** — the structured-extraction workhorse. Searches ~125–138M papers and populates
  **extraction tables**: you define columns (sample size, population, intervention, outcome,
  effect size) and it extracts across dozens of papers. Has a **Systematic Review flow** with
  override of inclusion criteria / screening / extracted fields. Best for evidence tables,
  scoping reviews, rapid evidence synthesis. Weakness: re-wording the question changes the
  papers it cites and its conclusions.
- **Consensus** — built only on peer-reviewed literature. Best for **yes/no/claim** questions:
  the **Consensus Meter** aggregates whether evidence says Yes / No / Possibly, plus a
  Claims-and-Evidence table. Ask a validatable question, not a topic.
- **Undermind** — autonomous multi-agent deep search with citation-graph reasoning; iteratively
  refines its search. Documented to beat one-shot embedding search on recall *and* precision
  (Tay, 2025). Best for finding seed papers / exploratory deep search. Not for formal SR.
- **SciSpace** — large corpus (~280M); "Deep Review" scores ~85%+ on ScholarQA-CS2. Strong at
  single-paper Q&A plus multi-paper synthesis.
- **Ai2 Asta / Scholar QA** — open-source, attribution-first RAG; generates organized reports
  with per-claim evidence and literature-review tables; publicly reports which papers it cites
  most. Best for transparent, reproducible synthesis.
- **FutureHouse PaperQA2** — accuracy-over-everything agentic RAG (open source); matched/beat
  PhD experts on LitQA2 and contradiction detection. Key design move: **re-ranking + contextual
  summarization (RCS)** before generation, plus metadata-aware retrieval and an abstain option.
  Best for high-stakes, full-text scientific QA where you can tolerate cost/latency.

### Verification tools
- **scite (Smart Citations)** — classifies each citation as **Supporting**, **Contrasting**, or
  **Mentioning**, with the citing sentence. Best for fact-checking your bibliography — has a key
  study been refuted or retracted before you cite it? Raw citation counts are meaningless;
  citation *context* is the signal.

### General deep-research agents (web-grounded)
- **OpenAI Deep Research** — strongest reasoning depth / structured synthesis; can fabricate
  sources buried in long reports; premium-priced.
- **Google Gemini Deep Research** — broad web breadth, tight Workspace/Docs export; prone to
  SEO-bias and occasionally cites sources that don't directly support the claim.
- **Perplexity (Pro) Deep Research** — fastest (~minutes), free tier, strong inline citation
  transparency; tuned toward authoritative sources. Less depth.

> **Cross-cutting limit:** independent benchmarks (DeepResearch Bench) find even strong agents
> are "fundamentally limited by poor citation grounding," weak at visualization and at very
> specialized technical sub-domains, and degrade on rapidly-changing/recent topics.

### Selection matrix

| Job to be done | Reach for | Avoid using for |
|---|---|---|
| Map the shape of a field from one paper | Connected Papers | Answering a question |
| Ongoing lit-review you'll revisit | ResearchRabbit (+ Zotero) | One-off deep answers |
| Build an evidence/extraction table | **Elicit** | Yes/no claim questions |
| "Does X cause/help Y?" claim question | **Consensus** | Open-ended exploration |
| Best-recall seed-paper / exploratory deep search | **Undermind** | Formal SR with full audit trail |
| Highest-accuracy full-text scientific QA | **PaperQA2**, Asta Scholar QA | Speed/cost-sensitive scanning |
| Check if a paper has been refuted | **scite** | Discovery; connections only |
| Broad multi-source web synthesis | OpenAI / Gemini / Perplexity DR | Anything needing guaranteed sourcing |
| Programmatic corpus access at scale | Semantic Scholar API / S2ORC / OpenAlex | Hand exploration |

Rule of thumb: **discover with graph tools, synthesize with an attribution-first agent, verify
with scite + entailment, and never trust one tool for all three.**

## 3. The AI-Assisted Review Workflow (and where it breaks)

A classic review is **search → screen → extract → synthesize**. AI changes each stage but moves
the human's job rather than removing it.

| Stage | What AI does well | The failure mode to guard |
|---|---|---|
| **Search** | Semantic/agentic search finds papers keyword search misses | Recall is *not* complete; recency gaps; corpus blind spots |
| **Screen** | Ranks likely-included studies to the top | **Tail risk**: an eligible study can rank deep in the ordering (one report found a relevant abstract only ~96% of the way down the ranked list, Ref 10) → can't stop early if you need *all* eligible studies |
| **Extract** | Auto-fills evidence tables (PICO, effect sizes) | Hallucinated cells: confidently invents randomization status, outcomes not in the text |
| **Synthesize** | Drafts cited, organized summaries; detects contradictions | Unsupported statements, fabricated citations, one-sidedness, prompt-sensitivity of conclusions |

**Operating rules (from the methods literature):**
- **Use AI at specific stages, not to automate the whole review.** For a true systematic review,
  humans still screen every abstract for full recall; AI ranking reduces *effort ordering*, not
  *required coverage*.
- **Force verbatim grounding.** In extraction/screening prompts, require the model to quote the
  exact sentence supporting each decision — turns a hallucination check into a glance.
- **Treat conclusions as prompt-sensitive.** Re-run the question 2–3 different wordings; if cited
  papers or conclusions move, the synthesis is unstable.
- **Disclose AI use.** Record tool identity, model version, prompts, stage of use, and
  human-verification steps — the emerging **PRISMA-trAIce** checklist / **AITDI** index.
- **AI tools lack validated reliability.** None of the consumer tools have established
  validity/accuracy evidence for SR-grade work — say so in limitations.

## 4. Citation Grounding & Faithfulness

How to *measure* whether an AI synthesis is actually supported by its sources.

**ALCE (Gao et al., EMNLP 2023)** is the reference frame. For each statement:
- **Citation recall** = is the statement *entailed* by the concatenation of its cited passages?
  (via an NLI model; rooted in the **AIS** "Attributable to Identified Sources" framework). High
  recall → well-supported.
- **Citation precision** = are there *irrelevant* citations attached? High precision → less
  human-review burden.
- Report the **F1** of the two.

**Stricter variants:** passage-grounded correctness (score only statements derivable from
retrieved passages, ignoring parametric memory); attribution groundedness (recall/precision on
the *citations*); and **correctness ≠ faithfulness** (Wallat et al., 2025) — also require the
model to *causally rely* on the cited document (alter the cited doc → the claim should change;
add irrelevant docs → attribution shouldn't change). Decorative citations pass correctness but
fail faithfulness.

**Practical grounding checklist for any AI synthesis you receive:**
1. Does each load-bearing claim carry an inline citation?
2. Open 3–5 citations at random — does the source actually state the claim (entailment), or
   merely touch the topic?
3. Are there claims with *no* citation that should have one?
4. Does the system abstain ("Insufficient Information") when it should, or always answer?
5. Are sources real and resolvable (no fabricated DOIs / dead links)?

## 5. Auditing Deep-Research Agents (failure modes, quantified)

**DeepTRACE** (Narayanan Venkit et al., 2025) turns community-identified failures into **eight
measurable dimensions** by decomposing the answer into statements and building a citation matrix
(statement × source cited) and a factual-support matrix (statement × source supports):

| # | Dimension | What it measures | Acceptable / Problematic |
|---|---|---|---|
| 1 | **One-Sidedness** | Skew toward one stance on debate queries | balanced ↔ one-sided |
| 2 | **Overconfidence** | Fraction of high-confidence, unhedged statements | low is good |
| 3 | **Relevant Statements** | % of statements addressing the query | ≥90 good / <70 bad |
| 4 | **Uncited Sources** | % of listed sources never cited | <5 good / ≥10 bad |
| 5 | **Unsupported Statements** | % of relevant statements no listed source supports | <10 good / ≥25 bad |
| 6 | **Source Necessity** | % of sources uniquely supporting a claim | 80–100 good |
| 7 | **Citation Accuracy** | % of statement→source citations that truly support | ≥90 good / <50 bad |
| 8 | **Citation Thoroughness** | % of *possible* accurate citations actually included | >50 good / <20 bad |

**Headline findings (expectations, not guarantees):** generative search + deep-research agents
are frequently **one-sided and overconfident** on debate queries; include large fractions of
statements **unsupported by their own listed sources**; deep-research configs reduce overconfidence
and can hit high thoroughness, but **citation accuracy ranges only 40–80%**. Two named failure
modes: **statement hallucination** (text deviates from the cited source) and **citation
hallucination** (the reference itself is fabricated). Takeaway: more sources ≠ stronger grounding;
always spot-check.

## 6. Content Provenance & Credentials (C2PA)

When the *artifact* (report, figure, image) must carry verifiable origin, use **C2PA Content
Credentials** — "a nutrition label for digital content."

- **Assertions** — statements about the asset: origin, edit actions, use of AI.
- **Claim** — assertions + content **bindings** bundled together.
- **Claim signature** — the claim is cryptographically signed by the **claim generator**,
  producing tamper-evidence.
- **Manifest / Content Credential** — the signed unit; the **Manifest Store** is the asset's
  full provenance and can reference other manifests (provenance chains).
- **Bindings**: *hard* (cryptographic hash → "this exact asset"); *soft* (invisible watermark) so
  credentials survive metadata stripping.

**Trust model:** trust is anchored in the **identity of the signer** of the claim — not in the
content itself. For AI-assisted research, C2PA is the standards-grade answer to "prove how this
was made" — use it for published figures/reports needing a cryptographic origin trail; use the
lighter mechanisms in §7 for the day-to-day claim→source trail.

## 7. Building an Auditable Research Trail

The goal (claim-level auditability literature): a synthesis is **auditable** if an independent
reviewer can verify each claim with **E_verify ≪ E_generate**, using only (1) the output, (2) a
structured provenance graph, and (3) the cited sources — *without* re-running the research.

**Action-trace logs are not enough.** W3C **PROV** (entities, activities, agents) records *what
the agent did*, but not *which source substantiates each claim*. Aim for **semantic provenance**:
a claim↔evidence graph that preserves relation type and surfaces contradictions.

**A practical, tool-agnostic trail:**
1. **Lock intent up front** — question, inclusion/exclusion criteria, admissible sources, success
   metric (a brief/preregistration). Log deviations.
2. **Bind every key claim to evidence** — claim → source ID → quoted passage → relation
   (supports / contradicts / context). One row per claim.
3. **Record the tool chain** — tool identity, model version, prompt, retrieval params, date
   (corpora drift). This is the PRISMA-trAIce/AITDI disclosure, operationalized.
4. **Verify before you trust** — run ALCE-style entailment spot-checks (§4) and/or a backward
   claim-tracer like **VeriTrail** (extract claims, verify each backward through the generation
   DAG, return Fully/Not-Fully-Supported/Inconclusive + an evidence trail).
5. **Make it replayable** — append-only ledger (JSONL with content hashes) so the synthesis can
   be re-derived and diffed.
6. **Sign the artifact if stakes warrant** — attach C2PA Content Credentials (§6).

**Minimal evidence-table schema:**

| claim_id | claim text | source (DOI/URL) | quoted passage | relation | tool+model | date | verified? |
|---|---|---|---|---|---|---|---|

## 8. Pitfalls

- **Wrong-layer tool use.** Asking Connected Papers to answer a question, or trusting a
  deep-research agent to *discover* exhaustively. Match tool to job (§2 matrix).
- **Trusting fluency as evidence.** A polished report is not a grounded one; citation accuracy can
  be 40–80% even when the prose is clean.
- **Counting citations instead of reading them.** A 100-citation paper may be cited because it was
  *wrong*. Use scite's supporting/contrasting context.
- **Decorative citations.** Correct-but-unfaithful attributions pass an NLI check yet the model
  never relied on the source. Demand verbatim grounding.
- **Fabricated sources buried in long reports.** Always resolve a sample of DOIs/links.
- **Action logs masquerading as provenance.** "Here's what the agent did" ≠ "here's which source
  supports this claim." Keep a semantic claim→evidence trail.

(The §3 operating rules — stop-screening tail risk, prompt-sensitivity, AI-use disclosure — and
the §1 corpus blind spots are the other standing failure modes; they are stated once where they
apply rather than repeated here.)

## References

1. Ai2 — AstaBench benchmarking blog: https://allenai.org/blog/astabench
2. Singh et al. — Ai2 Scholar QA (ACL 2025 Demo): https://aclanthology.org/2025.acl-demo.49.pdf
3. Skarlinski et al. — PaperQA2 / superhuman synthesis (2024): https://arxiv.org/pdf/2409.13740
4. Narayanan Venkit et al. — DeepTRACE (arXiv 2509.04499, 2025): https://arxiv.org/abs/2509.04499
5. Gao et al. — ALCE, citation generation/eval (EMNLP 2023): https://aclanthology.org/2023.emnlp-main.398.pdf · repo: https://github.com/princeton-nlp/alce
6. Wallat et al. — Correctness is not Faithfulness in RAG (2025): https://staff.fnwi.uva.nl/m.derijke/wp-content/papercite-data/pdf/wallat-2025-correctness.pdf
7. C2PA Content Credentials Explainer v2.4: https://spec.c2pa.org/specifications/specifications/2.4/explainer/Explainer.html
8. Semantic Scholar API: https://www.semanticscholar.org/product/api · S2ORC: https://github.com/allenai/s2orc
9. UBC Library — AI search tools comparison 2025: https://wiki.ubc.ca/AI-powered_search_tools_at-a-glance_comparison,_2025
10. JMIR (2026) — Humans still need to review all abstracts: https://formative.jmir.org/2026/1/e82896
11. BMC Med Res Methodol (2025) — Elicit for systematic review: https://link.springer.com/article/10.1186/s12874-025-02528-y
12. JMIR AI (2025) — PRISMA-trAIce checklist: https://ai.jmir.org/2025/1/e80247
13. Microsoft Research — VeriTrail (2025): https://www.microsoft.com/en-us/research/blog/veritrail-detecting-hallucination-and-tracing-provenance-in-multi-step-ai-workflows/
14. Aaron Tay — Academic Deep Research (2025): https://aarontay.substack.com/p/why-i-think-academic-deep-research
