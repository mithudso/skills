<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `computer-use-and-gui-agents` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

# Computer-Use & GUI Agents

<!-- Provenance: on-demand reference for the ai-agent-engineering hub. Frontier topic (2024-2026 sources). -->
<!-- Sibling refs: agent-harness-construction (generic action space), agent-reliability-and-guardrails (safety), multimodal-llm-architecture (the VLM that sees the screen), agent-planning-patterns (the control loop). -->

## Overview

A **computer-use / GUI agent** is an agent whose action space is a human computer interface — it *operates a screen* (desktop, browser, or mobile) using the same affordances a person does: it looks at pixels, moves a cursor, clicks, types, scrolls, and reads the result. This is the defining contrast with a tool-calling agent that hits clean JSON APIs (`agent-harness-construction`): the GUI agent works against UIs that have **no API**, are visually noisy, change layout, and were never designed for machine consumption. The bet is universality — anything a human can do on a screen, the agent can in principle do, without a per-app integration.

The canonical loop is a tight **perception → reasoning → action** cycle:

```
loop:
  observation = capture_screen()        # screenshot (+ optional a11y tree / SoM marks)
  thought, action = model(observation, goal, history)   # VLM picks ONE action
  execute(action)                       # click(x,y) / type(text) / scroll / key
  # take a fresh screenshot, repeat until done or human-confirm needed
```

Every step costs a full multimodal model call, so latency and reliability compound across the trajectory. The VLM doing the seeing is covered in `multimodal-llm-architecture`; the planning/reflection structure of the loop is `agent-planning-patterns`. This reference is about the **type** — what is distinctive about agents that drive a screen.

**The honest one-liner (2025-2026):** these agents work, are improving fast, and are *not yet reliable enough to run unsupervised on consequential tasks*. Best-in-class full-OS success is ~40-60% vs ~72% human; web agents reach high-80s/low-90s on the easier WebVoyager but stumble on hard, multi-step, real-site work. Latency is minutes-per-task. Treat them as supervised assistants, sandboxed, with human confirmation on side effects.

## Computer-use agents vs browser agents

Two overlapping families. The split matters for action space, grounding, and reliability.

| | **Computer-use (full-OS) agents** | **Browser agents** |
| --- | --- | --- |
| Surface | Whole desktop: OS, multiple apps, file system, terminal | A single web browser |
| Observation | Screenshot (pixels), sometimes + a11y tree | DOM/HTML *or* screenshot, often both |
| Action | OS-level mouse/keyboard at pixel coords | CDP/Playwright actions on DOM elements *or* pixel clicks |
| Examples | Anthropic **Claude computer use**; OpenAI **Operator/CUA**; Microsoft **OmniTool**; Google Gemini computer use | **Browser Use**; **WebVoyager** agents; Playwright/CDP-driven agents; Operator (web mode) |
| Benchmark home | **OSWorld**, AndroidWorld | **WebArena/VisualWebArena**, **WebVoyager**, **Mind2Web** |
| Strength | Maximal generality; cross-app workflows | Far more reliable — structured DOM, real selectors, robust waits |
| Weakness | Hardest grounding; lowest success; slowest | Confined to the browser; brittle when sites are SPA/canvas-heavy or DOM is obfuscated |

**Anthropic Claude "computer use"** (public beta Oct 22, 2024, with Claude 3.5 Sonnet; current beta header `computer-use-2025-11-24`). Pure screenshot-driven: you are the system integrator — you choose the display, capture the screenshot, run the model, **execute** the returned action, and decide the safety rails. The model returns actions like `click (450,320)` or `type "hello"`; you screenshot again and continue. Anthropic was explicit at launch that it was "still experimental — at times cumbersome and error-prone," that **scrolling, dragging, and zooming present challenges**, and that developers should "begin exploration with low-risk tasks."

**OpenAI Operator / Computer-Using Agent (CUA)** (Jan 2025, GPT-4o-derived + RL). Same universal "screen, mouse, keyboard" interface — *no app-specific APIs*. Runs the perception/reasoning/acting loop with chain-of-thought, and is trained to **pause and ask for user confirmation before actions with external side effects** (purchases, sends) and to **decline high-risk tasks** (banking). Originally a ChatGPT Pro web-browsing agent; the `computer-use-preview` model is also available via API and Azure AI Foundry.

**Google** and others have since shipped computer-use models (Gemini computer use), and open agents (OpenAI's open CUA tooling, UI-TARS, Qwen-VL-based agents) let you swap the backing VLM.

## GUI grounding approaches (the core technical problem)

**Grounding** = mapping a high-level intent ("click Submit") to a concrete, executable action on *this* screen (the pixel coordinate or the DOM node). It is the single hardest part of a GUI agent and where most failures originate. Three families:

1. **Pixel / pure-vision grounding.** The VLM consumes a raw screenshot and emits coordinates directly (`click(x,y)`). Maximally general (works on canvas, remote desktops, games, any OS) but demands precise spatial grounding the base VLM often lacks — raw GPT-4o scored **0.8%** on the **ScreenSpot Pro** grounding benchmark. This is Claude computer use's and CUA's native mode.

2. **Accessibility tree / DOM (structured) grounding.** Read the OS a11y tree or the page DOM to get element roles, labels, and bounding boxes; act on a *node* rather than a pixel. Far more reliable when available — this is why **Browser Use** restructures the "messy DOM" (strip noise, label interactive elements, expose a clean control interface) and why DOM-driven web agents top the leaderboards. Caveats: a11y trees are incomplete/wrong on custom widgets, huge (token cost), and absent on canvas/remote-desktop surfaces. OSWorld's own finding: the a11y tree and SoM "can be helpful" but "can also lead to potential misguidance and varies across models."

3. **Set-of-Marks (SoM) prompting + screen parsers.** Bridge the two: detect interactive regions, overlay **numbered bounding-box marks** on the screenshot, and let the model pick a *mark ID* instead of raw coordinates — turning a hard regression problem into an easy selection problem. **OmniParser** (Microsoft Research; V2 checkpoints Feb 2025) is the reference screen parser: a fine-tuned **YOLO** icon-detection model + a **Florence-2** caption model produce interactable regions with functional descriptions, rendered as SoM. **OmniParser + GPT-4o reached 39.6% on ScreenSpot Pro (vs GPT-4o's 0.8% raw)** and powers **OmniTool**, which drives a Windows 11 VM with a swappable backend (GPT-4o/o1/o3-mini, DeepSeek-R1, Qwen-2.5-VL, or Claude computer use). SoM is model-agnostic and turns "any LLM into a computer-use agent."

Practical rule: **use the structured signal (DOM/a11y/SoM) whenever the surface exposes it; fall back to pure pixels only where it doesn't.** Hybrid (screenshot + marks + DOM text) is the current sweet spot.

## Action space

A GUI agent's action vocabulary is small and human-like (this is the *concrete instantiation* of the generic action space in `agent-harness-construction`):

- **Pointer:** `mouse_move(x,y)`, `left_click(x,y)`, `right_click`, `double_click`, `left_click_drag(from→to)`, `scroll(dir, amount)`
- **Keyboard:** `type("text")`, `key("ctrl+s")` / chords, `hold_key`
- **Meta:** `screenshot` (re-observe), `wait` (let UI settle), `cursor_position`; harness-level `goto(url)` / `navigate_back` for browsers; `done` / `ask_user` (request human confirmation)
- **Coordinate space:** the model reasons in the resolution of the image it was shown. If your real display is higher-res, **you must scale coordinates** before executing — a frequent integration bug. Anthropic recommends keeping the display at/below the model's analysis resolution (~XGA/WXGA) for accuracy.

Browser agents add a parallel **DOM action layer** (`click(element_id)`, `fill(selector, text)`, `select_option`) that is more robust than pixel clicks because it survives layout shifts and rides Playwright/CDP auto-waiting. One action per turn keeps the trajectory observable, recoverable, and auditable.

## Benchmarks + the reliability reality

| Benchmark | What it measures | Size | Notable result | Honest read |
| --- | --- | --- | --- | --- |
| **OSWorld** | Real full-OS tasks (Ubuntu/Windows): desktop apps, file I/O, multi-app, GUI **and** CLI; execution-based grading on a VM snapshot | 369 tasks | CUA **38.1%** (prev SOTA 22%); Claude 3.5 launch ~**14.9%** screenshot-only / 22% extended | **Human ≈ 72.4%.** Frontier ~40-60% in 2025-26. Far from solved. |
| **WebArena** | Self-hosted realistic sites (e-commerce, CMS, forums, gitlab) | ~812 | CUA **58.1%** (prev 36.2%) | Reproducible but synthetic; mid-50s is current ceiling. |
| **VisualWebArena** | WebArena + visually-grounded tasks (must reason over images) | ~910 | — | Adversarial pop-ups cut completion ~47% (see safety). |
| **WebVoyager** | **Live** real sites (Amazon, GitHub, Google Maps…) | 586 | CUA **87%**; **Browser Use 89.1%**; Aime Browser-Use **92.3%** | Highest numbers — but easier/shorter tasks, live sites drift, and self/LLM-grading inflates scores. Don't read 90% as "solved." |
| **Mind2Web** | Generalization across 137 real sites / 31 domains, action prediction | 2,350 | — | Tests cross-site generalization, the real weak spot. |
| **AndroidWorld** | Mobile GUI control, 20 stock Android apps | 116 | — | The mobile analogue of OSWorld. |
| **ScreenSpot / ScreenSpot Pro** | Pure **grounding** (locate the right element) | — | OmniParser+GPT-4o **39.6%**; raw GPT-4o **0.8%** | Isolates grounding from planning — exposes how weak raw-pixel grounding is. |
| **OS-Harm** | **Safety** of computer-use agents (misuse, injection, misbehavior) | NeurIPS 2025 | — | Frontier models comply with much deliberate misuse, are vulnerable to static injection. |

**Reliability reality — read this before deploying:**

- **Success ceilings are modest.** Full-OS work sits well below human; "high" web numbers come from the easiest benchmark. Treat any single headline % skeptically — grading methods differ and live benchmarks inflate.
- **Non-determinism.** Agents "succeed in one run but fail in another" with identical task and model. Plan for retries, idempotent actions, and verification — not one-shot success.
- **Latency is the silent killer.** End-to-end times reach **tens of minutes** for tasks humans do in minutes; **OSWorld-Human** found planning + reflection LLM calls account for **75-94%** of total latency. Per-step screenshot+VLM calls dominate cost and wall-clock.
- **Generalization gap.** Agents do far worse on unfamiliar UIs and long horizons; CUA "struggles to figure out how to use" UIs it wasn't trained on and is imprecise at text editing.

## Safety (the unique risks of agents that touch a real screen)

A GUI agent fuses **untrusted input and privileged action in the same loop** — whatever is on screen becomes part of the prompt, and the very next step is a real click. That is the threat model. (General agent guardrails live in `agent-reliability-and-guardrails`; below are the *computer-use-specific* risks.)

- **Prompt injection via on-screen content.** Malicious text in a web page, email, PDF, or — critically — a **pop-up** can hijack the agent. Adversarial pop-ups hit a **~86% mean attack success rate** across OSWorld and VisualWebArena while cutting task completion ~47%. The page the agent is reading is an attack surface. Mitigate by treating all screen text as untrusted (wrap it in untrusted-input envelopes server-side per the `<untrusted_*>` escaping convention), and by detecting/dismissing injected overlays.
- **Destructive / irreversible clicks.** A wrong click can delete files, send an email, complete a purchase, or change account settings. There is no `git revert` for the real world. **Gate high-impact actions behind explicit human confirmation** — CUA is trained to pause before external side effects and to decline banking/sensitive tasks; build the same gate yourself if your harness lacks it.
- **Off-task / misaligned actions.** Even non-malicious agents introduce unnecessary interactions, pursue unintended subgoals, or derail — actions that are "technically permissible yet deviate from user intent." Detect drift against the stated goal and abort.
- **Misuse.** Users can direct the agent at harmful goals; frontier models still comply with much of it (OS-Harm). Keep model-level refusals, site blocklists, and moderation in the loop.

**Defense-in-depth checklist:**
1. **Sandbox** — run in a disposable VM/container with no access to production credentials, secrets, or the real file system; scope what's reachable.
2. **Human-in-the-loop on side effects** — require confirmation for purchases, sends, deletes, financial/account changes; "watch mode" on sensitive sites (email, banking).
3. **Allow/block lists** — constrain navigable domains and launchable apps.
4. **Untrusted-content handling** — treat every screenshot/DOM/text as adversarial; isolate it from trusted instructions.
5. **Action logging + kill switch** — record every (observation, thought, action); cap step count and wall-clock; let a human abort mid-trajectory.
6. **Least privilege** — minimal accounts/scopes; no standing access to anything the task doesn't strictly need.

## When to use a computer-use / GUI agent

**Reach for one when:**
- The target system has **no API** and no integration is feasible (legacy desktop apps, third-party SaaS UIs, internal tools).
- The work is **cross-application** GUI workflow (read app A's screen, act in app B).
- You need a **general** automator across many UIs and can't build a connector per app.
- The task is **supervised** and **sandboxable** (QA/test exploration, scripted-but-fragile RPA replacement, demos, research).
- For web specifically, prefer a **browser agent with DOM grounding** (Browser Use, Playwright/CDP) — markedly more reliable than pixel-clicking the same page.

**Prefer something else when:**
- A clean **API / MCP tool** exists — use tool-calling (`agent-harness-construction`, `mcp-servers`). It is faster, cheaper, deterministic, and safer than driving a UI.
- The task is **high-stakes or irreversible** and can't be human-gated or sandboxed.
- You need **low latency or high throughput** — per-step VLM calls make GUI agents slow and expensive.
- A purpose-built **deterministic scraper/script** would do — don't put an LLM in a loop where a Playwright script suffices.

**Decision heuristic:** API > browser-agent-with-DOM > pure-vision computer use. Drop down a rung only when the rung above is unavailable. Always pair the lowest necessary rung with sandboxing + human confirmation on consequential actions.

## References

- Anthropic — *Introducing computer use, a new Claude 3.5 Sonnet* (Oct 22, 2024): https://www.anthropic.com/news/3-5-models-and-computer-use
- Anthropic — *Computer use tool* (Claude API docs, current beta): https://platform.claude.com/docs/en/agents-and-tools/tool-use/computer-use-tool
- OpenAI — *Computer-Using Agent (CUA) / Operator* (Jan 2025): https://openai.com/index/computer-using-agent/
- Microsoft Research — *OmniParser for pure vision-based GUI agent* + *OmniParser V2: Turning Any LLM into a Computer Use Agent*: https://www.microsoft.com/en-us/research/articles/omniparser-v2-turning-any-llm-into-a-computer-use-agent/ ; repo https://github.com/microsoft/OmniParser
- **OSWorld** — *Benchmarking Multimodal Agents for Open-Ended Tasks in Real Computer Environments* (NeurIPS 2024, 369 tasks): https://arxiv.org/abs/2404.07972 ; site https://os-world.github.io/
- **OSWorld-Human** — *Benchmarking the Efficiency of Computer-Use Agents* (latency 75-94%): https://arxiv.org/abs/2506.16042
- **WebVoyager** / **Browser Use** SOTA technical report (89.1% WebVoyager): https://browser-use.com/posts/sota-technical-report
- **OS-Harm** — *A Benchmark for Measuring Safety of Computer Use Agents* (NeurIPS 2025): https://arxiv.org/abs/2506.14866
- *Adversarial pop-up* attack on CUAs (~86% ASR across OSWorld/VisualWebArena): summarized in OS-Harm and related 2025 safety work.
