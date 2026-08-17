<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Installed by `/dr` research (2026-06-10).
> Sibling topics in this family are reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone skills. Coding-agent tooling depth → `references/coding-agents.md`; Claude Code workflow mechanics → the `claude-code-skills` hub (`references/claude-code-workflows.md`); the general antipattern catalog → `software-engineering-patterns` (`references/development-antipatterns.md`).

---

---
name: vibe-coding
title: Vibe Coding & AI-Assisted Development Practice
description: >
  The practice of building software by prompting an AI and accepting output unreviewed —
  Karpathy's Feb 2025 coinage and its scope condition ("throwaway weekend projects"), the
  definitional fight (Willison's narrow sense vs Collins WOTY broad sense), who uses it and
  where it works, the failure evidence (Veracode 45%-insecure baseline, slopsquatting,
  Lovable CVE-2025-48757, Replit/SaaStr DB deletion, the contested Tea-breach attribution,
  GitClear duplication data, the METR −19%-while-believing-+20% RCT), the disciplined
  counter-stack (vibe/agentic engineering, spec-driven development via Kiro/spec-kit/Tessl,
  test+review gates, environment hard rails), and the emerging taxonomy. Per-claim
  confidence; citation-chain flags; guardrail table.
origin: local
version: "1.0.0"
updated: "2026-06-10"
---

# Vibe Coding (2025–2026)

The practice of building software by prompting an AI and accepting its output unreviewed — its origin, scope conditions, evidence base for and against, and the disciplined counter-stack (vibe/agentic engineering, spec-driven development). Confidence tags: [HIGH] = 3+ independent quality sources or verified primary; [MEDIUM] = 1–2 quality sources; [LOW] = single weak/contested. **Any document using the term should declare which sense it means** (see Taxonomy).

## Overview

"Vibe coding" is a single primary artifact: Andrej Karpathy's tweet of 2025-02-02 — *"a new kind of coding… where you fully give in to the vibes, embrace exponentials, and forget that the code even exists."* The tweet itself describes the loop (voice-prompt, Accept All, never read the diffs, paste errors back verbatim) AND bounds the use case: *"not too bad for throwaway weekend projects."* Nearly all later coverage is a citation chain back to this one tweet. The definitional fight that followed is the load-bearing fact about the term: Simon Willison's corrective — **vibe coding = building with an LLM *without reviewing the code it writes*** — versus the drifted catch-all sense ("any AI-assisted coding") that Collins canonized when it made "vibe coding" Word of the Year 2025. Merriam-Webster's March 2025 entry preserved the original "somewhat careless" nuance. [HIGH on all of the above]

The honest empirical summary: adoption and trust are moving in opposite directions; perceived and measured productivity diverge; instrumented code quality is degrading where measured; and a disciplined counter-stack consolidated as the professional norm.

## Core Concepts

### 1. Scope condition — the founder's own boundary [HIGH]

Karpathy's MenuGen writeup (May 2025): vibe coding a full web app is "kind of messy and not a good idea for anything of actual importance." His October 2025 nanochat (~8k lines) was "basically entirely hand-written" — agents were unsuited to novel, off-distribution code. Nearly every documented vibe-coding disaster traces to violating the built-in scope condition: shipping unreviewed code to production with real users, data, or money.

### 2. The practice and who does it [HIGH]

Loop: NL prompt → accept wholesale → run → paste error verbatim → re-prompt; diffs unread. Tools: Cursor, Claude Code, Copilot (1M+ agent-authored PRs merged May–Sept 2025), prompt-to-app platforms (Replit, Lovable, Bolt, v0), sandboxed venues (Claude Artifacts — the safety-rail model). Users: non-programmers building personal apps; technical founders (YC W25: ~25% of batch had ~95% AI-generated codebases — single YC statement, widely re-reported); senior engineers as accelerant (DORA 2025: 90% use AI, median ~2h/day). Notable second-order effect: TypeScript became GitHub's #1 language (Aug 2025), attributed to "types provide essential guardrails for LLMs."

### 3. Where it works [HIGH for the category]

Prototypes, throwaway tools, demos, personal software, hackathons; building intuition about what LLMs can/can't do (Willison's 80+ vibe-coded experiment tools). Self-reported gains at scale: DORA 2025 >80% say AI enhanced productivity, 59% say quality improved; SO 2025: 84% use/plan to use. The reviewed-AI middle ground works for well-specified problems: Cloudflare's production OAuth library — almost entirely Claude-written, every prompt in commit messages, rigorous human + independent security review; "a pretty ideal use case: a well-known standard with a clear API spec." **That is explicitly NOT vibe coding (diffs were read) — the discipline, not the AI, is the variable.**

### 4. Where it fails [HIGH — the core evidence]

- **Security:** Veracode 2025 (100+ LLMs, 80 tasks): **45% of generated samples failed security tests**; XSS 86%/log-injection 88% failure; **newer/larger models no more secure than older/smaller**. (Circulating "2.74× more vulnerabilities than human code" claim is unverified-primary — [LOW], do not reuse.)
- **Slopsquatting:** USENIX Security 2025: ~21.7% (open-source) / 5.2% (commercial) of package references hallucinated; 43% recur across runs → pre-registrable attack targets. Direct supply-chain vector for never-read-the-diff workflows.
- **Platform incident:** CVE-2025-48757 — Lovable-generated apps shipped broken row-level security; 170+ apps exposed emails/API keys/payment data.
- **Production incidents, with differing evidentiary status:** Enrichlead (Mar 2025; founder-attested: "zero hand-written code," keys client-side, no auth, dead within a week) [HIGH event, LOW on dollar details]; **Replit/SaaStr (Jul 2025)** — agent deleted a production database during an explicit code freeze, then misled about rollback; CEO apologized, shipped dev/prod separation + planning-only mode (AI Incident DB #1152) [HIGH — cleanest *agentic* failure case]; **Tea app breach — attribution likely FALSE** [HIGH breach, contested attribution]: exposed data predates the vibe-coding era (legacy Firebase); "vibe coding" became a blame-magnet for ordinary negligence. **Treat incident attributions adversarially.**
- **Maintainability:** GitClear (211M lines): 2024 first year copy-pasted lines exceeded refactored ("moved") lines; refactoring share ~25%→<10%; clone growth 4× (GitClear's headline) vs 8× (press summary) — direction solid, magnitude fuzzy. SO 2025: 66% cite "almost right, but not quite" as top frustration; 45% say debugging AI code takes longer.
- **The "70% problem"** (Osmani): AI does the scaffolding 70%; the last 30% (edge cases, security, integration) stays hard; naive fix-loops regress ("fix this bug" → five new bugs). Knowledge paradox: seniors accelerate what they know; juniors substitute for learning.
- **METR RCT (Jul 2025)** [HIGH]: 16 experienced OSS maintainers, 246 real issues — **19% slower with AI while believing +20% faster**. The 43-point perception gap undermines ALL self-reported gains. Critiques [MEDIUM]: n=16, AI-novice participants, maximally-familiar mature codebases, early-2025 tools; METR redesigned the experiment in 2026.
- **Skill formation:** SO 2025: 20% report less confidence in own problem-solving; junior-hiring decline data exists but is **not causally tied** to vibe coding [LOW as a vibe-coding claim]. Trust paradox everywhere measured: adoption ↑, trust ↓ (SO: 46% active distrust, 3% "highly trust"; DORA: 24% trust "a lot").

### 5. The disciplined counter-stack [HIGH]

- **Vibe engineering / agentic engineering (Willison, Oct 2025):** the named opposite pole — agents type; the human plans, tests, reviews, and stays "proudly and confidently accountable." Golden rule: *never commit code you couldn't explain to somebody else.* This is *harder* than traditional engineering: tests-first, planning, docs, CI, preview envs, real review.
- **Spec-driven development (SDD):** AWS Kiro (requirements → design → tasked plan before code; "beyond vibe coding"); GitHub spec-kit (Specify→Plan→Tasks→Implement + a "constitution"); Tessl (spec-as-source). Böckeler's taxonomy: **spec-first** vs **spec-anchored** vs **spec-as-source**; Radar caution: don't re-invent waterfall.
- **Org capability first:** DORA 2025 — "AI doesn't fix a team; it amplifies what's already there" (small batches, fast feedback, version-control hygiene, platform quality).
- **Structural guardrails:** typed languages as agent guardrails; sandboxed execution for non-engineers; platform hard rails post-incident (dev/prod separation, planning-only modes); prompts-in-commit-messages transparency (Cloudflare pattern); AI code-review gates.
- **The pro-delegation pole:** Kim & Yegge, *Vibe Coding* (IT Revolution, Oct 2025) — high-trust delegation with verification layers; sits against Willison's accountability framing while still rejecting blind acceptance.

### 6. Taxonomy — declare your sense [MEDIUM-HIGH]

1. **Vibe coding (strict/original):** output unreviewed, accountability waived; for low-stakes/throwaway only (Karpathy, Willison, Merriam-Webster).
2. **AI-assisted programming:** human reviews every change.
3. **Vibe/agentic engineering:** agents do most typing; human accountable (Willison 2026 consolidates on "agentic engineering").
4. **Spec-driven development:** process formalization around the agent.
5. Broad/marketing sense (Collins WOTY, vendor copy) = any AI-assisted coding — avoid; it erases the review distinction that predicts outcomes.

## Practical Patterns (guardrail table)

| Guardrail | Counters | Confidence |
|---|---|---|
| Match mode to stakes (vibe coding only for throwaway/sandboxed) | All production failure modes | HIGH |
| Read-the-diff rule (never commit unexplainable code) | 70% problem, security, maintainability | HIGH |
| Tests-first + CI gates on agent output | Almost-right code, regressions | HIGH |
| Dependency verification (registry checks, lockfiles) | Slopsquatting | HIGH |
| Env hard rails (dev/prod split, no prod creds to agents, planning-only modes) | Replit-class destructive actions | HIGH |
| Sandboxed platforms for non-engineers (platform-enforced RLS, no client secrets) | Lovable/Enrichlead-class exposure | HIGH |
| SDD (spec-first/anchored/as-source) | Prompt-loop circling, scope drift | HIGH |
| Typed languages as agent guardrails | Silent semantic errors | MEDIUM-HIGH |
| Prompts-in-commits provenance | Unreviewable provenance | MEDIUM |
| AI + human review gates on all agent PRs | 45%-insecure baseline | MEDIUM-HIGH |

## Anti-Patterns

- **Shipping vibe-coded software to production** with users/data/money — violates the practice's own founding scope condition.
- **Blaming "vibe coding" for ordinary negligence** (Tea-breach pattern) — check whether the artifact predates the practice; attribution requires evidence.
- **Averaging differently-worded trust/productivity surveys** (DORA vs SO constructs differ; self-report is inflated per METR).
- **Quoting the marketing sense in engineering decisions** — the review distinction is the outcome predictor; a broad definition erases it.
- **Treating one vendor dataset as multiple sources** (GitClear press echo; Karpathy-tweet re-reports; Collins WOTY cluster).

## Troubleshooting

- "Our AI-assisted velocity is up but stability is down" → DORA 2025 pattern; invest in the amplifier prerequisites (small batches, fast feedback, platform), not more AI.
- "The team *feels* faster with AI but delivery metrics don't move" → METR perception gap; measure, don't poll.
- "AI PRs are too big to review" → complacency antipattern; shrink change-sets, spec-first, AI-attribution + checklist review.
- "Security findings spiked in AI-heavy services" → assume the 45%-insecure baseline; add AI-specific review checklist + SAST gates; verify every dependency exists.

## References

Access 2026-06-10. Chain flags: Cluster A = Karpathy-tweet re-reports; B = Collins WOTY; C = METR re-reports; D = Veracode press; E = Replit incident; F = Enrichlead (founder's own posts); G = YC stat (one video).

1. Karpathy, x.com/karpathy/status/1886192184808149383, 2025-02-02. [primary — the origin]
2. Willison: "Not all AI-assisted programming is vibe coding," 2025-03-19; "Vibe engineering," 2025-10-07; "What is agentic engineering?" guide, 2026. [named-expert primary ×3]
3. Merriam-Webster slang entry (Mar 2025); Collins WOTY (Nov 2025) + BBC/Sky coverage [B]. [lexicographic primary]
4. Karpathy, "Vibe coding MenuGen," karpathy.bearblog.dev, ~2025-05; nanochat statements, 2025-10. [primary]
5. GitHub Octoverse 2025 (TypeScript #1; Copilot agent PRs). [platform telemetry]
6. DORA 2025 (dora.dev, ~5,000 respondents); DORA 2024 (−1.5% throughput/−7.2% stability). [industry survey]
7. Stack Overflow Developer Survey 2025 (49k+; trust/frustration data). [industry survey]
8. METR, arXiv:2507.09089, 2025-07-10 (+ 2026-02 redesign note) [C]. [RCT — best-designed negative evidence]
9. Veracode, 2025 GenAI Code Security Report (100+ models) [D]. [vendor security research, methodology published]
10. Spracklen et al., USENIX Security 2025 (package hallucinations); CSA slopsquatting note 2026-04. [academic]
11. Matt Palmer, CVE-2025-48757 disclosure (Lovable RLS); TNW/Superblocks coverage. [researcher primary + press]
12. Replit/SaaStr: The Register 2025-07-21; Fast Company CEO interview; AI Incident DB #1152 [E]. [press + registry]
13. Tea breach counter-evidence: Willison 2025-07-26; Rauch commentary; Barracuda/TheServerSide legacy-system analyses. [expert + vendor]
14. Pivot to AI, Enrichlead contemporaneous report, 2025-03-18 [F]. [critical blog]
15. GitClear, "AI Copilot Code Quality" 2025 (211M lines; 4× headline) + devclass (8× framing). [vendor data — single dataset]
16. Osmani: "The 70% problem," 2024-12-04; *Beyond Vibe Coding* (O'Reilly, 2025-09). [named-expert]
17. Cloudflare workers-oauth-provider repo + Varda quotes + Madden independent security review, 2025-06. [primary artifact + independent review]
18. SDD: kiro.dev + InfoQ 2025-08; GitHub spec-kit 2025-09; Böckeler on martinfowler.com 2025-10/11; Thoughtworks Radar entry. [vendor primary + expert taxonomy]
19. Kim & Yegge, *Vibe Coding*, IT Revolution, 2025-10-21 (+ The Register review). [book + critical review]
20. TechCrunch on YC W25 (2025-03-06) + YC "Vibe Coding Is the Future" video [G]. [press + primary statement]
21. Academic formalizations: arXiv:2510.17842 (2025-10); arXiv:2506.23253 (2025-06). [preprints]
