<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `voice-realtime-agents` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: voice-realtime-agents
title: Voice & Real-Time Agent Design
description: >
  The application / design layer of voice agents — how to build and operate a
  real-time conversational voice product, distinct from the model-internal
  speech architecture (Whisper encoders, Thinker-Talker, audio tokenization)
  owned by multimodal-llm-architecture. Covers the cascaded (STT→LLM→TTS) vs
  speech-native / speech-to-speech (S2S) architecture choice; real-time voice
  APIs (OpenAI Realtime gpt-realtime — WebRTC vs WebSocket, sideband server-side
  tools; Google Gemini Live native audio); turn detection & barge-in (server VAD,
  endpointing, semantic/model-integrated turn detection like Deepgram Flux,
  backchannel vs interruption, push-to-talk); the sub-800ms voice-to-voice
  latency budget (LLM TTFT ~70% of it, streaming everything, colocation); the
  streaming STT/TTS component landscape (Deepgram, AssemblyAI, Whisper;
  ElevenLabs, Cartesia); orchestration frameworks & telephony (LiveKit Agents,
  Pipecat, TEN; SIP/Twilio/PSTN); managed platforms (Vapi, Retell, Bland); voice
  UX (confidence-tiered confirmation, error recovery, consent); and voice-agent
  evaluation (4-layer framework, WER limits, interruption accuracy, MOS, TSR/FCR,
  simulated callers). TRIGGER: building a voice agent / voice bot / phone agent;
  "cascaded vs speech-to-speech", OpenAI Realtime API, Gemini Live, VAD / turn
  detection / endpointing / barge-in, voice latency budget, LiveKit / Pipecat /
  Vapi / Retell / Bland, streaming STT/TTS choice, telephony / SIP for agents,
  voice-agent evaluation. SKIP: model-internal speech architecture — audio
  tokenization, Whisper encoder, Thinker-Talker (use multimodal-llm-architecture);
  text chat/agent UI (use frontend-ui ▸ ai-native-ux-generative-ui); generic
  agent orchestration with no real-time/voice constraint (use agent-ecosystem).
origin: local
category: developer
version: "1.0"
updated: "2026-05-31"
tags:
  - voice-agents
  - realtime
  - speech-to-speech
  - stt
  - tts
  - latency
  - llm
  - agent
whenToUse:
  - "building or operating a real-time voice agent / phone bot"
  - "choosing cascaded (STT→LLM→TTS) vs speech-native / speech-to-speech architecture"
  - "using OpenAI Realtime API or Gemini Live API (transport, sessions, tools)"
  - "designing turn detection / endpointing / barge-in (VAD, semantic, Flux, backchannel)"
  - "engineering the sub-800ms voice-to-voice latency budget"
  - "selecting streaming STT/TTS and orchestration/telephony (LiveKit, Pipecat, Vapi, Retell, Bland, SIP)"
  - "designing voice UX (confirmation, error recovery, consent) and evaluating voice agents"
whenNotToUse:
  - "model-internal speech architecture (audio tokens, Whisper encoder, Thinker-Talker) — use multimodal-llm-architecture"
  - "text chat / generative UI design — use frontend-ui ▸ ai-native-ux-generative-ui"
  - "generic agent orchestration without a real-time/voice constraint — use agent-ecosystem"
related_skills:
  - multimodal-llm-architecture
  - ai-native-ux-generative-ui
  - agent-ecosystem
  - llm-inference-serving
---

# Voice & Real-Time Agent Design

A voice agent listens, reasons, and speaks inside a turn-taking loop tight enough
to feel human. This reference is the *application/operation* layer — architecture,
real-time APIs, turn-taking, latency, components, UX, eval. For *how the model
produces audio* (audio tokenization, Whisper encoders, Thinker-Talker), see
`multimodal-llm-architecture`.

The engineering north star is a **voice-to-voice latency budget under 800 ms**,
and the field has converged on **streaming everything, semantic turn detection,
and first-class barge-in** as table stakes.

## 1. Architecture: cascaded vs speech-native (S2S)

- **Cascaded** = an orchestration layer over swappable services: VAD → streaming
  STT → LLM → streaming TTS. Wins on **control** (a text intermediary for
  filters, audit, fallbacks), **debuggability** (isolate STT vs LLM vs TTS
  failures), **provider choice**, and **tool/reasoning** strength.
- **Speech-to-speech (S2S)** = one model ingests and emits audio in a single
  latent space (OpenAI `gpt-realtime`, Gemini Live). Wins on **latency**
  (200–300 ms) and **emotional prosody**, but is **opaque** and hard to gate.
- **As of mid-2026, most production agents are still cascaded** (S2S < ~15% H1
  2026, projected ~25–30% H2 2026 — directional forecast). Use **S2S for
  empathy-first** experiences; **cascaded for regulated / high-volume /
  tool-heavy**; or **hybrid** (S2S conversation, cascade for compliance-gated
  branches). A streaming cascade (Deepgram + a fast LLM + Cartesia) already hits
  sub-1s, narrowing S2S's main edge.

## 2. Real-time voice APIs

- **OpenAI Realtime API** (GA 2025, `gpt-realtime`): **WebRTC** for client media
  (browser/mobile), **WebSocket** for server/telephony media. Ephemeral client
  secrets keep permanent keys off the client. Sessions configured via
  `session.update`. **Keep tools/business logic server-side via a sideband
  control channel** (two connections to one session). GA adds **async function
  calling** with placeholder responses ("still waiting…") to suppress
  hallucination during pending calls.
- **Gemini Live API** (GA on Vertex): Gemini 2.5 Flash native audio, 30 voices /
  24 languages, emotion-aware, live S2S translation, multimodal (converse about a
  live visual stream), tool use + Google Search.

## 3. Turn detection, endpointing & barge-in (the hardest UX problem)

Three layers of increasing sophistication:

- **Server VAD** (e.g., Silero) — classifies each frame as speech/non-speech.
  Cheap but sees only energy, not meaning.
- **Endpointing** — silence-timeout on top of VAD. The latency trap: an 800 ms
  silence timeout adds ~1 s to *every* turn.
- **Semantic / model-based turn detection** — reads the partial transcript and
  infers completeness; can trigger *before* trailing silence, letting thresholds
  drop to **200–300 ms** without false interruptions. OpenAI exposes Server VAD +
  Semantic VAD modes.
- **Model-integrated (2026 frontier):** **Deepgram Flux** folds end-of-turn
  detection into the ASR model (~260 ms median end-of-turn) — collapsing
  ASR+VAD+endpointing into one.

**Barge-in:** keep turn detection live *during* playback; on user speech, **cancel
the TTS stream and hand control back to STT within ~200 ms**; echo cancellation
client-side; push-to-talk as fallback. **Backchannels** ("mm-hm") are the subtle
case — pure VAD misreads them; 2026 stacks use a turn-taking model to classify
backchannel vs barge-in.

## 4. Latency engineering (target < 800 ms)

Component latencies are **cumulative and sequential**; a naïve "each stage waits"
pipeline hits 2–3 s. Representative optimized budget: turn detection ~50 ms · STT
~150 ms · **LLM TTFT ~400 ms** · TTS first chunk ~150 ms · network ~50 ms ≈
**~800 ms**.

- **LLM inference is ~70% of the budget** — model selection is the biggest lever
  (fast tier: GPT-4o-mini ~400 ms TTFT, Claude Haiku ~360 ms).
- **Optimize TTFT, not total tokens** — TTS begins on the first complete sentence.
- **Mitigations:** streaming STT (−100–200 ms), streaming TTS (−200–400 ms),
  semantic endpointing (−200–400 ms), geographic colocation (cross-continent adds
  80–250 ms). (See `llm-inference-serving` for the model-serving side of TTFT.)

## 5. Components: streaming STT & TTS

- **STT (2026):** Deepgram leads voice-agent latency + end-of-speech (Nova-3 +
  Flux as the default); ElevenLabs Scribe v2 Realtime leads multilingual;
  AssemblyAI leads transcript intelligence; Whisper covers 57+ languages.
- **TTS:** positioning consensus — **ElevenLabs = realism/cloning**, **Cartesia =
  streaming latency**; Deepgram Aura and others compete on time-to-first-audio.
  (Specific TTFA benchmark numbers are single-vendor-published — directional.)

## 6. Orchestration frameworks & telephony

- **Open source:** **Pipecat** (Python pipeline/frame model; elegant for 1:1
  assistants; you assemble production deploy), **LiveKit Agents** (WebRTC-first;
  ships production transport — rooms, SFU, recording, SIP; best for
  multi-participant / latency-sensitive / telephony), **TEN** (third contender).
- **Managed (telephony-included):** **Vapi** (BYO-everything middleware; weaker
  native telephony), **Retell** (voice-quality + sub-500 ms leader; warm
  transfer, native SIP, HIPAA/SOC2/GDPR), **Bland** (API-first, high-volume
  outbound).
- **Telephony:** PSTN → WebRTC room via **SIP trunking** (buy a number on
  Twilio/Vonage/Telnyx, point its Voice URL at the framework's SIP URI); LiveKit
  supports DTMF, transfers, region pinning, noise cancellation.

## 7. Voice UX

Voice UI's constraint: **nothing is visible and the spoken sentence is gone the
moment it lands** — so error recovery is a *primary* discipline.

- **Confidence-tiered confirmation:** high → act + implicit confirm ("I've sent
  the invoice"); medium → clarify ("3 contacts named John — which?"); low →
  graceful fallback. Implicit confirmation beats explicit yes/no in enterprise.
- **Error recovery as trust:** users forgive the first error, doubt the second,
  give up by the third — state confusion + offer concrete next steps.
- **Consent before consequential action.** (The text/visual equivalents of these
  patterns are in `frontend-ui ▸ ai-native-ux-generative-ui`.)

## 8. Evaluation

A **4-layer framework**: (1) Infrastructure (audio quality, latency, ASR/TTS
perf), (2) Execution (intent, response accuracy, tool-calling), (3) User-behavior
(interruption handling, flow, sentiment), (4) Business outcome (containment,
FCR, escalation). Headline metrics: **WER** target <5% but it "ignores
interaction dynamics" (barge-in, endpointing) — don't evaluate on WER alone;
**interruption handling** (stop within ~200 ms, address >90%); **MOS** ~4.5/5 for
near-human; **Task Success / FCR** ~85%+. Gate releases with **simulated-caller
regression tests** and instrument **per component** (`call_id`/`turn_id`),
watching p95 tail latency.

## Practical patterns

1. Stream every stage; never block (start TTS on first sentence, LLM on first
   finalized STT segment).
2. Optimize LLM TTFT (≈70% of budget); pick a fast tier model, keep the prompt
   small.
3. Use semantic / model-integrated turn detection to push silence to 200–300 ms.
4. Make barge-in first-class; classify backchannel vs interruption.
5. Keep tools/business logic server-side (sideband channel).
6. Match transport to context (WebRTC client, WebSocket server/telephony).
7. Tier confirmations by confidence; treat error recovery as a primary surface.
8. Instrument per component; gate releases with simulated-caller regression.
9. Default cascaded for regulated/tool-heavy; reserve S2S for empathy-first.

## Anti-patterns

- Blocking (non-streaming) pipeline; over-long silence timeout; endpointing that
  holds past the fallback (tail-latency jump); no/poor barge-in or treating
  backchannels as barge-ins; retry cascades without hard timeouts; weak multi-turn
  context ("cancel that one" three turns deep); deploying S2S into regulated flows
  (no text checkpoint); WER-only evaluation.

## Cross-references

- **Model-internal speech** (audio tokenization, Whisper, Thinker-Talker) →
  `multimodal-llm-architecture`.
- **Text/visual AI UX equivalents** (streaming, confirmation, error UX) →
  `frontend-ui ▸ ai-native-ux-generative-ui`.
- **LLM-serving TTFT / fast inference** → `llm-inference-serving`.
- **General agent orchestration / tool design** → `agent-ecosystem`,
  `agent-harness-construction`.

## References

OpenAI Realtime API docs/blog; Google Gemini Live API docs. LiveKit (turn
detection, telephony) + Pipecat; Deepgram Flux. Latency: Hamming, Smallest.ai,
Twilio, Cresta. Components: futureagi, Softcery, AssemblyAI; Gradium TTS
benchmark. Architecture: Coval, Speko. UX: InfoWorld, Smashing, fuselabcreative.
Eval: Hamming, Cekura, Braintrust. Platforms: Retell/Vapi/Bland comparisons.
*(50 sources, 2024–2026; full URLs in the source research report.)*
