<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** A new reference added for the AI-native UX / generative-UI design layer.
> Sibling topics in this family are now reference files under the hubs (`frontend-ui`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: ai-native-ux-generative-ui
title: AI-Native UX & Generative UI Design
description: >
  Design the interface and interaction layer of LLM/AI applications — the
  "design" half of building with LLMs, distinct from the model/agent
  engineering layer. Covers streaming UX (TTFT as the key metric, markdown
  buffering, smooth streaming, stop/regenerate, streaming a11y), generative UI
  (tool-call → typed message parts → components; Vercel AI SDK useChat /
  streamObject / useObject; the paused RSC streamUI path; Thesys C1 / OpenUI
  schema-driven runtimes), latency masking & perceived performance (skeletons,
  optimistic UI, thinking indicators), trust & calibration UI (claim-adjacent
  citations, confidence cascades vs numeric scores, the reasoning-trace-display
  caveat, anthropomorphization & disclaimers), refusal / error / fallback UX
  (the Google PAIR error taxonomy, "be human not machine", graceful
  degradation), human-AI steering (prompt controls, suggestions, feedback
  affordances, mixed-initiative), and the four governing design frameworks
  (Microsoft 18 Guidelines for Human-AI Interaction, Google PAIR People+AI
  Guidebook, Apple Generative AI HIG, the Shape of AI pattern library).
  TRIGGER: designing a chat / assistant / agent UI; "streaming UX", "time to
  first token UX", "generative UI", "Vercel AI SDK useChat / streamObject",
  rendering tool calls as components, AI loading / skeleton / typing states,
  showing citations or confidence or reasoning in an AI UI, AI error / refusal /
  "I don't know" UX, prompt-input and feedback (thumbs / regenerate / edit)
  design, applying Microsoft HAX / Google PAIR / Apple AI HIG / Shape of AI
  patterns. SKIP: the backend/model/agent engineering layer — RAG, agent
  orchestration, observability, guardrails (use the ai-agent-engineering hub);
  generic visual design, color, typography, layout with no AI-specific
  interaction (use web-design / ui-ux-pro-max); WCAG/ARIA review of a non-AI UI
  (use accessibility-ux-reviewer); model-internal speech/vision architecture
  (use ai-agent-engineering ▸ multimodal-llm-architecture).
origin: local
category: developer
version: "1.0"
updated: "2026-05-31"
tags:
  - frontend
  - ux
  - ui
  - ai
  - llm
  - generative-ui
  - streaming
  - design
keywords:
  - AI-native UX
  - generative UI
  - streaming UX
  - time to first token
  - Vercel AI SDK
  - useChat
  - streamObject
  - trust calibration
  - citations UX
  - AI error UX
  - Microsoft HAX guidelines
  - Google PAIR
  - Apple Generative AI HIG
  - Shape of AI
whenToUse:
  - "designing a chat, assistant, copilot, or agent interface"
  - "implementing streaming responses and the UX around them (TTFT, stop, regenerate, smooth streaming, a11y)"
  - "building generative UI — rendering model tool-calls or structured output as live components"
  - "deciding how to show citations, confidence, or reasoning traces without inflating false trust"
  - "designing AI error, refusal, and graceful-degradation states"
  - "designing prompt input, suggestions, and feedback affordances for steerability"
  - "grounding an AI product against Microsoft HAX / Google PAIR / Apple HIG / Shape of AI"
whenNotToUse:
  - "backend/model/agent engineering (RAG, orchestration, observability, guardrails) — use ai-agent-engineering"
  - "generic visual design with no AI interaction — use web-design / ui-ux-pro-max"
  - "WCAG/ARIA review of a non-AI UI — use accessibility-ux-reviewer"
related_skills:
  - frontend-design
  - accessibility-ux-reviewer
  - ai-agent-engineering
---

# AI-Native UX & Generative UI Design

The interface and interaction layer for LLM/AI products — the **design** half of
building with LLMs. The `ai-agent-engineering` hub owns the model, retrieval,
orchestration, and reliability layers; this reference owns what the user
actually sees and touches. AI-native UX has consolidated into a recognizable
pattern language; three findings dominate.

1. **Streaming is the foundational UX primitive.** Users perceive streaming
   responses as ~40% faster than buffered ones at *identical* total latency, so
   **time-to-first-token (TTFT)** — not completion time — is the metric to
   design around.
2. **Generative UI has standardized on tool-call → typed `parts` → components.**
   The model emits a tool call; the result streams back as a typed message part
   that maps to a rendered component. The Vercel AI SDK is the reference
   implementation.
3. **The trust/calibration layer is genuinely contested.** Citations and
   "show-your-work" reasoning traces can *inflate unjustified* trust; design
   them skeptically.

Four authoritative frameworks govern the space (§7): Microsoft's 18 Guidelines
for Human-AI Interaction, Google's People + AI (PAIR) Guidebook, Apple's
Generative AI HIG, and the community Shape of AI pattern library.

## 1. Streaming UX

**TTFT is the design metric.** Perceptual thresholds: <0.5s feels instant,
0.5–1s responsive, >1.5–2s sluggish unless the UI explains the wait. Streaming
TTFT is typically 200–500ms vs 5–30s for a buffered full response — and the
perceived-speed gain (~40%) holds even when total time is unchanged.

**The hard part is rendering, not receiving.** Transport is usually
Server-Sent Events (SSE) over a persistent connection (WebSockets for
bidirectional/voice). Named rendering techniques:

- **Markdown buffering ("render only complete structures")** — defer rendering
  of unclosed bold/italic markers, code fences, partial table rows, and
  incomplete list items so formatting doesn't flicker as ambiguous characters
  arrive.
- **Code-highlight strategy** — defer syntax highlighting until block
  completion, or progressive highlight with a 200–300ms debounce; detect
  language from the opening fence.
- **Layout-thrash prevention** — min-heights on response containers; grow
  vertically without sibling reflow; scroll the response area independently.
- **Smooth streaming** — batch/animate token reveal rather than rendering each
  raw chunk. This is the *perceived*-TTFT lever product teams control
  independently of infra TTFT.

**Controls.** A **Stop** button must be visible during streaming and must
preserve partial output (it also saves token cost); **Regenerate** and **Copy**
appear post-completion; offer a **rendered/raw toggle**.

**Streaming a11y.** `aria-live="polite"` + `aria-atomic="false"` so screen
readers announce new content without re-reading; `aria-busy="true"` during the
stream; batch announcements every 2–3s, not per token; never steal input focus
mid-stream.

## 2. Generative UI

**Definition:** letting the model go beyond text to "generate UI." The current
recommended pattern is **tool-call-driven rendering**:

- **Tool definition** — a `description`, a Zod **`inputSchema`**, and an
  `execute` returning structured data shaped like component props.
- **API route** — `streamText()` with the tools, returning
  `toUIMessageStreamResponse()`.
- **Typed message `parts`** — each message carries a `parts` array. Text parts
  `{ type:'text', text }`; tool parts `{ type:'tool-{name}', state, input/output }`.
  Tool-part **states** are the canonical render switch: `input-available`
  (loading) → `output-available` (data) / `output-error`.
- **Client render** — switch on `part.type` + `part.state`; spread
  `part.output` into the component.

**Use `useChat` (client-side), not the paused RSC path.** AI SDK RSC
(`streamUI` / `createStreamableUI`) streams React Server Components from the
server, but **its development is officially paused** — Vercel recommends AI SDK
UI for production. Decision rule: RSC/Server Actions → `streamUI`; client hooks
→ `streamText` + `useChat`, or `streamObject` + `useObject`.

**Structured-output-to-component.** `streamObject` (core) constrains output to a
Zod/JSON schema and streams **partial objects** as fragments arrive; the client
`useObject` (still `experimental_`) progressively builds the object so the UI
renders incrementally — best for large response structures.

**`useChat` state machine.** `status` ∈ `submitted | streaming | ready | error`
is the canonical source for skeleton / typing / stop affordances; the hook also
exposes `sendMessage`, `regenerate`, `stop`, `error`, `clearError`,
`addToolOutput`, `setMessages`.

**Schema-driven runtimes beyond Vercel.** **Thesys C1** is an OpenAI-compatible
middleware returning a structured UI spec (forms/tables/charts/layouts) rendered
by its React SDK; it also publishes **OpenUI** as an open generative-UI
standard. (Vendor productivity claims — "80% less frontend code" — are marketing,
unverified; the architecture is sound.)

## 3. Latency masking & perceived performance

- **Skeletons / shimmer / placeholders** during the pre-first-token gap — "if
  nothing visible happens, users think submit failed."
- **Thinking / typing indicators** and honest queue messaging ("high demand,
  your request is queued") — *transparency beats silence*.
- **Optimistic UI** — render the user's turn immediately plus an assistant
  placeholder before the stream starts.
- **Progressive disclosure** — users read early tokens while later ones arrive;
  streaming *is* a latency mask.

## 4. Trust & calibration UI (the contested layer)

Goal: **calibrated trust** — neither over- nor under-reliance. (For the
psychology underneath this — automation bias, algorithm aversion/appreciation,
the trust-calibration curve — see `applied-psychology ▸ human-ai-interaction-psychology`.)

- **Citations / source attribution.** Make citations prominent and *separate*
  from the response; place them **adjacent to the specific claim**; use
  meaningful labels (article/publication titles); link to the relevant section;
  set the expectation upfront that sources may be wrong. Perplexity is the
  exemplar (inline numbered, retrieve-then-synthesize). **Contested:** users
  rarely click citations yet still gain *false* confidence from their presence.
- **"Show your work" / reasoning traces — contested.** Shape of AI's "Stream of
  Thought" and shipping reasoning models display step-by-step traces; **NN/g
  advises against presenting them as explanation** because traces are often
  post-hoc rationalizations, not faithful accounts. Resolution: collapsible
  reasoning aids steerability/debugging but should **not** be framed as
  ground-truth explanation.
- **Confidence display.** Prefer **behavioral confidence cascades** (high → act
  + confirm; medium → ask to clarify; low → state uncertainty, narrow scope)
  over raw numeric percentages, which are debated.
- **Anthropomorphization & disclaimers.** Avoid first-person phrasing implying
  human thought and avoid personas/backstories that inflate perceived
  capability; put **specific, actionable** disclaimers ("double-check AI
  outputs") *near the input*, not buried in a footer.

## 5. Refusal / error / fallback UX

**Error taxonomy (Google PAIR).** *System limitations* (can't give a right/any
answer) · *context errors* (working as intended but broke the user's mental
model) · *user-perceived vs invisible* failstates. Sources: prediction/training
errors, input errors, relevance errors (low confidence / bad timing),
system-hierarchy conflicts.

**Principle: "be human, not machine."** Address mistakes with humility; make
failure *safe, boring, and a natural part of the product.*

**Graceful degradation** answers three questions — *what still works? what
doesn't? what next?* Offer 2–3 recovery options (retry / wait / reduced mode),
let users **revert** an AI decision, and route to a **human or manual path**
when the AI can't solve it. "I don't know" should suggest a next step (rephrase,
narrow scope); guardrail-blocks should acknowledge + explain + offer an
alternative, never a cryptic refusal. Maps to HAX **G9/G10** (efficient
correction; scope when in doubt).

## 6. Human-AI steering & feedback

**Prompt-input design — the 4 uses of prompt controls (NN/g):** (1) increase
**discoverability** (show upload/feature affordances), (2) **educate/inspire**
(conversation starters, prompt libraries — ~19–30% of users don't know GenAI's
capabilities), (3) **set constraints/scope** up front (source pickers, output
modes), (4) **facilitate follow-ups** (77% of conversations have >1 exchange —
offer edit / regenerate / quick modifiers like "shorter"). Use standard icons
**with labels**; group controls by purpose; don't override standard gestures.

**Suggestions & mixed-initiative.** System-generated prompt suggestions and
below-answer follow-ups reduce effort and drive continuation; proactively ask a
framing question instead of launching a generic response.

**Feedback affordances.** Standard icons (favorite, thumbs-up/down, copy, share)
at the bottom of each reply; **edit-in-place** of generated artifacts;
**Regenerate** and **Variations** for side-by-side comparison. Google PAIR's
three rules: align feedback with model improvement, communicate its value +
time-to-impact, and keep easy opt-out.

## 7. Governing frameworks

- **Microsoft — 18 Guidelines for Human-AI Interaction** (the canonical rubric,
  4 phases): *Initially* G1 make clear what it can do, G2 how well. *During* G3
  time on context, G4 show relevant info, G5 social norms, G6 mitigate bias.
  *When wrong* G7 efficient invocation, G8 efficient dismissal, G9 efficient
  correction, G10 scope when in doubt, G11 explain why. *Over time* G12 remember,
  G13 learn, G14 update cautiously, G15 granular feedback, G16 convey
  consequences, G17 global controls, G18 notify of changes. (HAX Toolkit also
  ships Design Patterns, a Design Library, and a prioritization Workbook.)
- **Google — People + AI (PAIR) Guidebook.** Six chapters: User Needs + Defining
  Success · Data Collection + Evaluation · Mental Models · Explainability + Trust
  · Feedback + Control · Errors + Graceful Failure, each with patterns *and*
  anti-patterns.
- **Apple — Generative AI HIG.** Transparency-first: communicate where the app
  uses AI, manage expectations, maintain user control, build in safety
  guardrails on input/output.
- **Shape of AI** — community pattern library, ~60 named patterns in 6
  categories: **Wayfinders** (start the first prompt), **Inputs/Prompt Actions**,
  **Tuners** (refine via context/params), **Governors** (human-in-the-loop:
  Action Plan, Citations, Stream of Thought, Verification), **Trust Builders**
  (Disclosure, Footprints, Watermark, Consent), **Identifiers** (Avatar, Name,
  Personality).

**Agent UX specifically:** surface the agent's intended plan before execution,
reveal tool calls (Shape of AI Action Plan / Footprints / Stream of Thought),
and gate consequential actions behind confirmation (Verification; HAX G11).

## Practical patterns (decision rules)

- Stream by default; design for TTFT; mask the pre-first-token gap with a
  skeleton; animate token reveal.
- Generative UI = tool-calls → typed `parts` → components via `useChat` +
  `streamText`; `streamObject`/`useObject` + Zod for structured-output-to-UI;
  avoid the paused RSC `streamUI`.
- Citations inline and claim-adjacent + an honest disclaimer near the input;
  treat reasoning traces as steering aids, not explanations; prefer confidence
  cascades over numbers.
- Classify failures (system / context / invisible); give 2–3 recovery paths +
  revert + human handoff; never a cryptic refusal.
- Make steering cheap: discoverable prompt controls, suggestions, standard
  feedback icons, mixed-initiative framing.

## Anti-patterns

- Optimizing total latency instead of TTFT; buffering the whole response.
- Rendering raw token chunks (flicker) or unbuffered markdown (formatting jank);
  layout thrash from unbounded containers.
- Using the paused RSC `streamUI` path for new production work.
- Citations/reasoning shown as proof of correctness (false-confidence inflation).
- Numeric confidence scores presented as precise truth.
- Cryptic refusals / "I don't know" with no next step; burying disclaimers in a
  footer; anthropomorphic personas that oversell capability.
- Per-token screen-reader announcements; stealing focus during a stream.

## Cross-references

- **Backend/model/agent layers** (RAG, orchestration, observability, guardrails,
  multimodal/voice internals) → `ai-agent-engineering` hub references.
- **Generic UI craft** (color, type, layout, design systems) → `web-design`,
  `ui-ux-pro-max`, `frontend-design` (this hub).
- **Accessibility audit** of the resulting UI → `accessibility-ux-reviewer`.
- **The psychology of trust/reliance** behind §4 →
  `applied-psychology ▸ human-ai-interaction-psychology`.

## References

Vercel AI SDK — Generative UI, `useChat`, `streamUI` (paused), structured data;
GitHub discussion #2162 (streamUI vs streamObject decision). Redis — streaming
LLM responses & TTFT. TheFrontKit — streaming UI techniques (markdown buffering,
layout thrash, a11y). NN/g — Explainable AI in chat, Prompt Controls, Prompt
Suggestions, GenAI UX research agenda. Microsoft Research / HAX Toolkit — 18
Guidelines for Human-AI Interaction. Google PAIR — People+AI Guidebook (Feedback
+ Control, Errors + Graceful Failure). Apple — Generative AI HIG. Shape of AI —
pattern library. Thesys C1 / OpenUI — schema-driven generative-UI runtime.
OrangeLoops / ZipTie — trust, confidence cascades, Perplexity citation patterns.
*(Full URLs in the source research report; 24 sources, 2024–2026.)*
