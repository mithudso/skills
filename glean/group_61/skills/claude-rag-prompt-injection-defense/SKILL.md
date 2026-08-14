---
name: rag-prompt-injection-defense
description: >-
  Defending RAG and tool-using LLM apps against prompt injection, for builders
  and TAMs advising regulated (FSI) customers. TRIGGER: direct & indirect /
  cross-domain prompt injection via retrieved docs; the "lethal trifecta"
  (private data + untrusted content + exfiltration); RAG guardrails; defensive
  architectures (Dual-LLM, CaMeL "defense by design", spotlighting, tool-call
  allow-listing, human-in-the-loop for high-agency actions); OWASP LLM01
  injection, LLM02 info disclosure, LLM05 output handling, LLM06 excessive
  agency; why input filtering alone is insufficient. SKIP: embedding inversion /
  vector-store leakage -> embedding-inversion-threat-model; Atlas Vector Search
  tenant/PII isolation -> atlas-vector-search-pii-isolation; general RAG
  architecture -> ai-rag-retrieval; context-engineering hardening ->
  ai-mcp-sdk-prompting; app/web security -> security-review; bank AI governance ->
  bank-genai-model-risk-governance.
version: 1.1.0
updated: 2026-06-29
category: custom
whenToUse:
  - How do I defend a RAG or agentic LLM app against prompt injection?
  - What is the "lethal trifecta" and how do I break it?
  - How does indirect / cross-domain prompt injection work through retrieved documents?
  - What are the real defensive architectures (Dual-LLM, CaMeL, spotlighting, allow-listing)?
  - Which OWASP LLM Top 10 items apply to my RAG pipeline and how do I mitigate them?
  - Why isn't an input-filter / "injection classifier" enough on its own?
  - A bank customer asks how to make their Atlas Vector Search RAG assistant safe against injection.
keywords:
  - prompt injection defense
  - lethal trifecta Simon Willison
  - indirect prompt injection retrieved documents RAG
  - cross-domain prompt injection
  - dual LLM pattern
  - CaMeL defeating prompt injections by design
  - spotlighting delimiting datamarking
  - OWASP LLM Top 10 2025 LLM01 LLM02 LLM05 LLM06
  - excessive agency tool call allowlist
  - sensitive information disclosure exfiltration
  - human in the loop high agency actions
  - RAG guardrails input output guardrails
tags:
  - ai-security
  - prompt-injection
  - rag
  - llm-security
  - owasp-llm
  - fsi
  - security
  - tam
---

# RAG & Agent Prompt-Injection Defense

How to defend retrieval-augmented and tool-using LLM systems against prompt
injection. Written for builders and for TAMs advising regulated customers who
are deploying Atlas Vector Search + LLM assistants. **Core thesis: prompt
injection is not fully solvable by input filtering — you defend by
*architecture* that removes one leg of the lethal trifecta or constrains agency.**

## 1. The threat: injection vs jailbreak

- **Prompt injection** = untrusted input overrides the developer's instructions.
  Distinct from a **jailbreak** (coaxing the model past its safety training).
- **Direct injection**: the user types the attack.
- **Indirect / cross-domain injection**: the attack is **embedded in content the
  model retrieves** — a document chunk, web page, email, PDF, or a record returned
  by `$vectorSearch`. This is the RAG-specific danger: the *retrieved* text is
  treated as instructions. Authoritative framing: OWASP **LLM01:2025 Prompt
  Injection** (genai.owasp.org).

## 2. The lethal trifecta (Simon Willison)

Catastrophic data theft becomes possible when **all three** are present at once:

1. **Access to private data** (the RAG corpus, customer PII, tools that read secrets).
2. **Exposure to untrusted content** (retrieved docs, user input, third-party text).
3. **Ability to exfiltrate** (outbound HTTP, tool calls, rendering attacker-controlled
   links/images, sending email).

> Source: Simon Willison, "The lethal trifecta for AI agents" (simonwillison.net,
> 2025-06). His key claim: **remove any one leg and the data-theft class of attack
> collapses.** A RAG bank assistant that reads customer data *and* ingests untrusted
> documents *and* can call outbound tools has all three.

**Design move:** for high-value workloads, deliberately **break a leg** — e.g. no
outbound network/exfiltration channel from the data-bearing context; or isolate
untrusted content from the privileged tool-calling context.

## 3. Why input filtering is not enough

Injection-detection classifiers and "ignore previous instructions" regexes are
**probabilistic and bypassable** (encoding, obfuscation, multilingual, payload
splitting). Treat them as **defense-in-depth, not the control**. The durable
defenses are architectural (next section). Source: Willison; OWASP LLM01.

## 4. Defensive architectures (the durable controls)

- **CaMeL — "Defeating Prompt Injections by Design"** (Google DeepMind, arXiv
  2503.18813, 2025). A **privileged LLM** plans and emits a constrained program;
  a **quarantined LLM** handles untrusted data but **cannot trigger actions**; a
  capability/policy layer governs every tool call against data provenance. Untrusted
  text can never become a control-flow decision. Inspired by classic
  software-security (control/data-flow separation).
- **Dual-LLM pattern** (Willison, 2023): a privileged orchestrator never sees raw
  untrusted text; a quarantined LLM processes it and returns only structured,
  non-executable results. CaMeL is the hardened evolution of this.
- **Spotlighting / delimiting / datamarking** (Microsoft, 2024): mark untrusted
  content (delimiters, encoding, per-token markers) so the model can distinguish
  data from instructions. Mitigation, not a guarantee.
- **Tool-call allow-listing + least privilege**: enumerate permitted tools/params;
  deny by default; scope credentials per task. Directly attacks **LLM06 Excessive
  Agency**.
- **Human-in-the-loop for high-agency actions**: require explicit confirmation for
  state-changing or money-moving operations (wire, delete, send, external call).
- **Output handling**: never auto-execute, auto-render, or auto-follow model output.
  Encode/sanitize before it hits a browser, shell, SQL, or another system —
  **LLM05 Improper Output Handling**.

## 5. OWASP LLM Top 10 (2025) — the RAG-relevant subset

| ID | Risk | Primary mitigation |
|---|---|---|
| **LLM01** | Prompt Injection (direct + indirect) | Architectural isolation (CaMeL/Dual-LLM), spotlighting, least privilege |
| **LLM02** | Sensitive Information Disclosure | Minimize data in context; entitlement-filter retrieval; output redaction |
| **LLM05** | Improper Output Handling | Treat output as untrusted; encode/sanitize before downstream use |
| **LLM06** | Excessive Agency | Allow-list tools, least-privilege creds, human approval for high-impact actions |
| **LLM08** | Vector & Embedding Weaknesses | See `embedding-inversion-threat-model` |

Source: OWASP GenAI Security Project, "OWASP Top 10 for LLM Applications 2025"
(genai.owasp.org). ⚠️ The list is versioned annually — cite the year and re-check
the current edition.

## 6. What a MongoDB TAM should recommend to a bank

1. **Entitlement-filter retrieval** so the RAG corpus only returns documents the
   *requesting user* may see — push the access boundary into the `$vectorSearch`
   pre-filter (see `atlas-vector-search-pii-isolation`). This shrinks LLM02 blast radius.
2. **Break the trifecta**: no outbound exfiltration path from the PII-bearing context;
   isolate untrusted-document handling from privileged tool calls.
3. **Least-privilege tools + human approval** for any money-moving/state-changing action.
4. **Treat model output as untrusted** at every downstream boundary.
5. **Defense-in-depth** input/output guardrails on top — never as the only line.
6. Map the design to **OWASP LLM Top 10** and to the bank's model-governance program
   (`bank-genai-model-risk-governance`) so it survives validation/audit.

## Sources

- Simon Willison, "The lethal trifecta for AI agents," simonwillison.net (2025-06).
- Debenedetti et al., "Defeating Prompt Injections by Design" (CaMeL), arXiv:2503.18813 (2025).
- Microsoft, "Spotlighting" defense against indirect injection (arXiv:2403.14720, 2024).
- OWASP GenAI Security Project, "OWASP Top 10 for LLM Applications 2025," genai.owasp.org.

⚠️ AI-security guidance moves fast; CaMeL and the OWASP list are the most likely to
update. Re-verify the current OWASP edition and any newer "by-design" defenses before
relying on specifics in a customer deliverable.
