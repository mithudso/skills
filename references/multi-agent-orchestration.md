<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `multi-agent-orchestration` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agent-engineering` hub.** Sibling refs in this family: `agent-council` (debate/voting/consensus), `agent-ecosystem` (framework choice), `agent-planning-patterns` (decomposition/planning), `a2a-interop` (cross-agent protocol). Load those from `references/<name>.md`; ignore any "use the X skill" pointers that name a bare sibling.

# Multi-Agent Orchestration Topologies

How to **structure and coordinate multiple LLM agents** — the coordination-topology taxonomy, when multi-agent helps vs hurts, and the framework mapping. This is the structure lane; the *decision-making* lane (debate, voting, consensus, LLM-as-judge) lives in `agent-council`, deep framework selection in `agent-ecosystem`, the wire protocol for cross-vendor agents in `a2a-interop`, and tmux-based parallel agent harnesses in `dmux-workflows`.

## Overview

A "multi-agent system" (MAS) replaces one LLM loop with several LLM loops that have distinct prompts, tools, and (usually) separate context windows, plus a coordination scheme that routes work and merges results. The design space is captured by four dimensions ([Multi-Agent Collaboration Mechanisms survey, arXiv 2501.06322](https://arxiv.org/abs/2501.06322)): **actors** (which agents), **interaction type** (cooperation / competition / coopetition), **structure** (centralized / decentralized / distributed), and **coordination protocol** (how messages and control flow). Topology choice is the highest-leverage decision — a 2026 financial-MAS evaluation found coordination structure dominates outcomes ("coordination primacy") more than the underlying model ([arXiv 2603.27539](https://arxiv.org/pdf/2603.27539)).

Critical framing: **most "agentic" systems should be a single agent with good tools, not a MAS.** Anthropic's own agent guidance tells builders to "find the simplest solution possible, and only increase complexity when needed" — workflows and single agents first, multi-agent only when the task demands it ([Anthropic, Building Effective Agents](https://www.anthropic.com/engineering/building-effective-agents)). Read the [helps-vs-hurts decision](#the-helps-vs-hurts-decision) before reaching for any topology below.

## Workflow vs agent vs multi-agent (the altitude ladder)

Climb only as high as the task forces you. Lower rungs are cheaper, more reliable, and easier to debug.

1. **Single LLM call** — one prompt, no tools/loop.
2. **Workflow** — LLMs orchestrated through *predefined code paths*. Anthropic's five workflow patterns ([Building Effective Agents](https://www.anthropic.com/engineering/building-effective-agents)): **prompt chaining** (sequential steps, each output feeds the next), **routing** (classify input → specialized handler), **parallelization** (sectioning into independent subtasks, or voting = same task N times — voting/aggregation detail lives in `agent-council`), **orchestrator-workers** (a central LLM decomposes and delegates dynamically), **evaluator-optimizer** (generator + critic loop). Workflows are deterministic scaffolds; the control flow is fixed in code.
3. **Single autonomous agent** — one LLM dynamically directs its own tools/loop until done (see `autonomous-loops`).
4. **Multi-agent system** — multiple agents with separate contexts coordinate. Only here do the topologies below apply.

The orchestrator-worker *workflow* (fixed code orchestration) and the orchestrator-worker *multi-agent topology* (LLM orchestrator spawning LLM subagents) are the same shape at different altitudes — the distinction is whether the orchestrator is code or a model.

## Topology catalog

| Topology | Control | Context model | Best for | Main tradeoff |
| --- | --- | --- | --- | --- |
| **Single-threaded agent** | One loop | Continuous, unified | Coupled tasks, coding, anything needing consistent decisions | No parallelism; long tasks blow the context window |
| **Sequential pipeline** | Fixed order | Hand-off down the line | Staged transforms (extract→transform→summarize) | Error compounding; latency = sum of stages |
| **Parallel / sectioning** | Fan-out, then merge | Isolated per worker | Independent subtasks, breadth-first search | Merge conflicts; results may disagree |
| **Orchestrator-worker (supervisor)** | Central LLM manager | Orchestrator holds plan; workers isolated | Dynamic decomposition where subtask count/shape is unknown upfront | Orchestrator is bottleneck + single failure point; token multiplier |
| **Hierarchical (manager→sub-managers→workers)** | Multi-level supervisors | Tree; context narrows down levels | Large org-like task trees too big for one supervisor | Deep delegation chains; latency and error accumulate per level |
| **Network / peer-to-peer** | Any agent → any agent | Decentralized | Open-ended problems with no clear hierarchy | Hardest to reason about; can loop/thrash; non-determinism |
| **Handoff / swarm** | Decentralized, baton-passing | Full conversation travels with handoff | Routing among specialists (triage→billing→refunds) | No global planner; relies on each agent knowing when to hand off |
| **Group chat** | Shared conversation + speaker-selector | Shared transcript | Brainstorm/collaboration, emergent task division | Token blowup; can stall or go in circles without a strong manager |
| **Blackboard** | Shared store + control shell | Shared memory substrate, agents read/write | Decoupled specialists ("knowledge sources"), easy add/remove | Coordination logic moves into the control shell; shared-state contention |

### Supervisor / orchestrator-worker

A central agent receives the request, decides which worker to call, reads the result, then either calls the next worker or returns the answer ([LangGraph Multi-Agent Supervisor](https://reference.langchain.com/python/langgraph-supervisor)). The supervisor owns task delegation and result synthesis; workers stay isolated and specialized. This is the **default MAS topology** because it concentrates planning in one place (easier to reason about than peer-to-peer) while still parallelizing the leaf work. Use when subtask count or shape is **not known upfront** and must be decided dynamically — that dynamic decomposition is what separates the orchestrator-worker *topology* from a fixed parallelization *workflow* ([Building Effective Agents](https://www.anthropic.com/engineering/building-effective-agents)).

### Hierarchical (manager → workers)

Supervisors under supervisors. A top supervisor delegates to mid-level supervisors, each managing its own worker pool ([LangGraph multi-agent overview](https://machinelearningplus.com/gen-ai/langgraph-multi-agent-systems-supervisor-swarm-network/)). Reach for this only when a single supervisor's context or tool-routing surface gets too large; each extra level adds latency and another place for error to compound.

### Network / peer-to-peer

No manager — every agent can call every other agent ([LangGraph](https://towardsai.net/p/l/a-complete-guide-to-multi-agent-systems-in-langgraph-network-to-supervisor-and-hierarchical-models)). Maximum flexibility, minimum predictability. Justified only for genuinely open-ended problems with no natural hierarchy; expect harder debugging and a real risk of agents looping on each other.

### Handoff / swarm

Decentralized control by baton-passing. Each agent holds explicit **handoff** tools to transfer the active conversation to a peer; the system tracks the last-active agent so the next turn continues there ([LangGraph swarm](https://medium.com/@sameernasirshaikh/langgraph-swarm-vs-langgraph-supervisor-ce8194837d0a)). OpenAI's **Swarm** (experimental, 2024) reduced this to two primitives — a **routine** ("a list of instructions in natural language … along with the tools necessary to complete them") and a **handoff** ("an agent … handing off an active conversation to another agent"), implemented as a tool that *returns the next agent* while the full prior conversation persists ([OpenAI Cookbook: Orchestrating Agents](https://developers.openai.com/cookbook/examples/orchestrating_agents)). Swarm was succeeded in 2025 by the production **OpenAI Agents SDK**, which keeps the routine/handoff primitives and adds guardrails, tracing, and sessions ([openai/swarm](https://github.com/openai/swarm)). Best for customer-service-style routing among specialists; weak when a task needs a global plan no single specialist holds.

### Group chat (AutoGen / AG2) and Magentic-One

Agents post to a **shared conversation**; a speaker-selection policy (round-robin, model-chosen, or a manager) picks who talks next. AutoGen v0.4 reimplemented this on an asynchronous event-driven runtime ([microsoft/autogen](https://github.com/microsoft/autogen)). **Magentic-One** is a generalist orchestrator-led team built on AutoGen: an **Orchestrator** maintains a **Task Ledger** (verified facts, facts to look up, derived facts, educated guesses) and a **Progress Ledger** (self-reflection on progress + completion check), creates a plan, delegates to specialist agents (WebSurfer, FileSurfer, Coder, ComputerTerminal), and **dynamically revises the plan** when progress stalls ([Magentic-One, Microsoft Research](https://www.microsoft.com/en-us/research/articles/magentic-one-a-generalist-multi-agent-system-for-solving-complex-tasks/); [arXiv 2411.04468](https://arxiv.org/pdf/2411.04468)). The dual-ledger pattern is the notable transferable idea: separate *what we know* from *how the plan is progressing*.

### Blackboard

A 1970s pattern (Hearsay-II speech recognition) revived for LLMs: specialist agents ("knowledge sources") read from and write to a **shared blackboard**; a control shell picks who acts next based on board state ([CallSphere](https://callsphere.ai/blog/blackboard-architecture-multi-agent-systems-shared-knowledge-spaces)). Agents communicate *indirectly* through shared memory rather than direct messages, which decouples them — you add/remove specialists without touching the rest. Recent LLM work (bMAS) makes the blackboard the *exclusive* communication and memory substrate, broadcasting requests so each agent autonomously decides whether to contribute ([LLM-based Multi-Agent Blackboard System, OpenReview](https://openreview.net/pdf?id=egTQgf89Lm)). Trade: coordination intelligence concentrates in the control shell, and shared state becomes a contention point.

## Decomposition, delegation, and context model

- **Role specialization & decomposition.** Split by *capability* (researcher, coder, critic) or by *subproblem* (one agent per data source). The orchestrator must hand each subagent an explicit objective, output format, tool/source boundaries, and clear task limits — vague subagent prompts are a top failure cause ([Anthropic multi-agent research system](https://www.anthropic.com/engineering/multi-agent-research-system)).
- **Shared vs isolated context.** The central tension. *Isolated* contexts (orchestrator-worker, parallel) let each agent compress its own slice and exceed a single window in aggregate — but agents can't see each other's implicit decisions. *Shared* context (group chat, blackboard, single-thread) keeps everyone aligned but burns tokens and eventually overflows. See the [decision](#the-helps-vs-hurts-decision); context strategy detail is in `llm-context-engineering`.
- **Message passing.** Direct (handoff/swarm — conversation travels with the baton), broadcast (group chat — all see all), or indirect via shared store (blackboard). Cross-vendor/process message formats (Agent2Agent, ACP) are covered in `a2a-interop`.

## Worked case: Anthropic's multi-agent research system

The canonical orchestrator-worker-for-research case study ([How we built our multi-agent research system](https://www.anthropic.com/engineering/multi-agent-research-system); annotated by [Simon Willison](https://simonwillison.net/2025/Jun/14/multi-agent-research-system/)):

- **Architecture.** A **LeadResearcher** (Claude Opus 4) decomposes the query and spawns **parallel subagents** (Claude Sonnet 4) that each search independently with interleaved thinking, then return condensed findings. The lead synthesizes and decides whether to spawn more subagents or finish.
- **Why it works here.** Research is **breadth-first** — independent directions explored in parallel, each subagent compressing many tokens of search results into its own context window, so the system collectively exceeds a single window. Result: the multi-agent system **outperformed single-agent Opus 4 by 90.2%** on their internal research eval.
- **Token economics.** "Token usage by itself explains 80% of the performance variance" on BrowseComp; agents use **~4× the tokens** of chat, and **multi-agent systems use ~15× the tokens** of chat. So MAS only pays off on **high-value tasks where the quality gain justifies the cost** ([Anthropic](https://www.anthropic.com/engineering/multi-agent-research-system)).
- **Where they say it does NOT fit (in Anthropic's own words).** "Domains that require all agents to share the same context or involve many dependencies between agents are not a good fit for multi-agent systems today." Most coding "involves fewer truly parallelizable [sub]tasks than research," and "LLM agents are not yet great at coordinating and delegating to other agents in real time" ([via Slashdot summary](https://developers.slashdot.org/story/25/06/21/0442227/anthropic-deploys-multiple-claude-agents-for-research-tool---says-coding-is-less-parallelizable)).

## The helps-vs-hurts decision

This is a live, honest debate between two credible labs that shipped on opposite sides — present both.

**The skeptic's case — Cognition AI, "Don't Build Multi-Agents" (Walden Yan).** Two principles of *Context Engineering*: **(1) "Share context, and share full agent traces, not just individual messages."** **(2) "Actions carry implicit decisions, and conflicting decisions carry bad results."** The Flappy Bird example: ask two subagents to build a clone in parallel — one renders a Super-Mario-style background, the other an incompatible bird, because neither sees the other's implicit choices; the final agent inherits two miscommunications. Conclusion: prefer a **single-threaded linear agent** with continuous context; for long tasks, add a **compression model** that distills history into key decisions/events rather than forking into parallel agents. Stated exception: subagents are acceptable for **well-defined, read-only questions** (how Claude Code uses them) — never parallel *writes* that must stay mutually consistent ([Cognition, Don't Build Multi-Agents](https://cognition.ai/blog/dont-build-multi-agents)).

**The proponent's case — Anthropic (above) and LangChain.** For **breadth-first, parallelizable, read-heavy** work, isolated subagents are not just viable but win big (the 90.2% result). LangChain frames it as a context-management tradeoff: multi-agent splits context across agents to dodge single-window limits and reduce cross-talk, but pays in token cost and coordination overhead — build a MAS when the task is decomposable and "read-heavy," stay single-agent when subtasks are tightly coupled or write-heavy ([LangChain: How and when to build multi-agent systems](https://blog.langchain.com/how-and-when-to-build-multi-agent-systems/)).

**Reconciliation (the practitioner synthesis).** The two camps barely disagree once you split by **read vs write** and **coupling**:

```
Is the task decomposable into INDEPENDENT subtasks?
├─ No  → single-threaded agent (Cognition). Coupled work needs unified context.
└─ Yes → Are subtasks mostly READ/gather, or WRITE/produce-shared-artifact?
         ├─ Read/gather (research, multi-source lookup) → orchestrator-worker MAS (Anthropic) ✔
         └─ Write/shared artifact (most coding) → single agent, OR MAS only with
            strict interface contracts so writes can't conflict.
THEN gate on economics: is the task valuable enough to pay the ~15× token cost?
   └─ No → drop to a workflow or single agent regardless of decomposability.
```

Decide multi-agent when **all** hold: subtasks are independent, mostly read/exploratory, exceed one context window or one tool surface, and the task value clears the ~15× token premium. Otherwise prefer a single agent or a fixed workflow. Consensus/voting topologies (multiple agents answering the *same* question to cross-check) are a separate axis — see `agent-council`.

## Framework map (brief — defer deep choice to `agent-ecosystem`)

| Framework | Primary topologies | Coordination primitive |
| --- | --- | --- |
| **LangGraph** | Supervisor, swarm, hierarchical, network, custom | Stateful graph (nodes=agents, edges=control); `langgraph-supervisor` & `langgraph-swarm` prebuilts ([docs](https://reference.langchain.com/python/langgraph-supervisor)) |
| **OpenAI Agents SDK** (ex-Swarm) | Handoff/swarm | Handoff = tool returning next agent; full conversation persists ([cookbook](https://developers.openai.com/cookbook/examples/orchestrating_agents)) |
| **Microsoft AutoGen / AG2** | Group chat, orchestrator (Magentic-One) | Async event-driven messaging; speaker selection; Task/Progress ledgers ([repo](https://github.com/microsoft/autogen)) |
| **CrewAI** | Crews (autonomous) + Flows (deterministic) | **Crews** = role-based autonomous collaboration; **Flows** = event-driven step orchestration; nest either in the other ([CrewAI](https://crewai.com/crewai-flows)) |

CrewAI's split is the clean mental model for the whole space: **Flows = deterministic orchestration (the workflow rungs), Crews = autonomous reasoning (the MAS rungs)** — choose where you want autonomy and where you want determinism ([CrewAI docs](https://docs.crewai.com/en/introduction)). For tmux/terminal-multiplexed parallel agent runs (vs in-framework topologies), see `dmux-workflows`.

## Anti-patterns

- **MAS-by-default.** Reaching for multiple agents before trying a single agent or a fixed workflow. Start simple; add agents only when a measured limit forces it ([Building Effective Agents](https://www.anthropic.com/engineering/building-effective-agents)).
- **Parallel writers, no shared context.** Forking subagents to *produce* parts of one coherent artifact — the Flappy-Bird failure. Parallelize reads, serialize writes (or impose strict interface contracts).
- **Vague subagent briefs.** Spawning subagents without explicit objective, output format, tool boundaries, and limits → duplicated or conflicting work ([Anthropic](https://www.anthropic.com/engineering/multi-agent-research-system)).
- **Ignoring the token multiplier.** Deploying a ~15× MAS on low-value or latency-sensitive tasks. Gate on task value; a chat-grade query should never hit a research swarm.
- **Deep hierarchies / unbounded peer networks.** Each delegation level and each open peer edge adds latency and a compounding-error surface, and invites loops. Keep trees shallow; bound message hops.
- **Lossy hand-offs.** Passing only the last message (not the relevant trace) across a handoff/pipeline boundary, so the receiver re-derives or contradicts prior decisions ([Cognition](https://cognition.ai/blog/dont-build-multi-agents)).
- **No termination/merge story.** Group-chat and network topologies that lack a stop condition or a deterministic result-merger stall or thrash.

## References

- Anthropic — [How we built our multi-agent research system](https://www.anthropic.com/engineering/multi-agent-research-system) (orchestrator-worker, 90.2%, ~15× tokens, when-not-to-fit). *JS-heavy page; numbers cross-checked against the [Simon Willison annotation](https://simonwillison.net/2025/Jun/14/multi-agent-research-system/) and [Slashdot summary](https://developers.slashdot.org/story/25/06/21/0442227/anthropic-deploys-multiple-claude-agents-for-research-tool---says-coding-is-less-parallelizable).*
- Anthropic — [Building Effective Agents](https://www.anthropic.com/engineering/building-effective-agents) (5 workflow patterns; simplest-solution-first).
- Cognition AI (Walden Yan) — [Don't Build Multi-Agents](https://cognition.ai/blog/dont-build-multi-agents) (context engineering, single-threaded, Flappy Bird).
- LangChain — [How and when to build multi-agent systems](https://blog.langchain.com/how-and-when-to-build-multi-agent-systems/) (context-management tradeoff).
- LangGraph — [Multi-Agent Supervisor](https://reference.langchain.com/python/langgraph-supervisor) + [supervisor/swarm/network overview](https://machinelearningplus.com/gen-ai/langgraph-multi-agent-systems-supervisor-swarm-network/).
- OpenAI — [Orchestrating Agents: Routines and Handoffs (cookbook)](https://developers.openai.com/cookbook/examples/orchestrating_agents); [openai/swarm](https://github.com/openai/swarm) (succeeded by the Agents SDK).
- Microsoft — [Magentic-One](https://www.microsoft.com/en-us/research/articles/magentic-one-a-generalist-multi-agent-system-for-solving-complex-tasks/) ([arXiv 2411.04468](https://arxiv.org/pdf/2411.04468)); [microsoft/autogen](https://github.com/microsoft/autogen).
- CrewAI — [Flows](https://crewai.com/crewai-flows) and [Introduction / Crews](https://docs.crewai.com/en/introduction).
- Surveys — [Multi-Agent Collaboration Mechanisms (arXiv 2501.06322)](https://arxiv.org/abs/2501.06322); [coordination-primacy financial-MAS taxonomy (arXiv 2603.27539)](https://arxiv.org/pdf/2603.27539); blackboard: [bMAS (OpenReview)](https://openreview.net/pdf?id=egTQgf89Lm).

*Cross-references: `agent-council` (debate/voting/consensus — the same-question cross-check axis), `agent-ecosystem` (deep framework selection), `agent-planning-patterns` (task decomposition & planning), `a2a-interop` (cross-vendor agent wire protocols), `dmux-workflows` (tmux parallel agent harnesses), `autonomous-loops` (single-agent loop design), `llm-context-engineering` (context compression/sharing).*
