<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `agent-workflow-builder_ai_toolkit` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: agent-workflow-builder_ai_toolkit
version: 1.1.0
updated: 2026-05-29
description: >
  Generates, enhances, develops, and deploys AI agent applications and workflows
  using Microsoft Agent Framework. TRIGGER: user asks to create, scaffold, build,
  modify, fix, trace, monitor, debug, evaluate, measure, or deploy AI apps,
  agents, or workflows using the Microsoft Agent Framework or Azure AI Foundry.
  SKIP: agent frameworks other than Microsoft Agent Framework unless the user
  explicitly requests a migration or comparison.
origin: local
tags: [microsoft-agent-framework, azure-ai, foundry, agent-creation, python, dotnet, tracing, evaluation]
related_skills: [agent-ecosystem, agent-plan-writing, mcp-servers]
---

# Building AI Agent / Workflow (Microsoft Agent Framework)

## Critical instructions

- **Interpret intent**: Accurately capture what the user wants. Execute specific capabilities or multiple as needed. Ask if unclear.
- **SDK exclusivity**: Use **Microsoft Agent Framework** exclusively. Do not apply if the user explicitly asks for a different SDK or package.

## When NOT to use this skill

- User has a different agent framework (LangGraph, CrewAI, Claude SDK, etc.) — use `agent-ecosystem`
- User wants to plan agent orchestration architecture without building — use `agent-plan-writing`
- User is building an MCP server — use `mcp-servers`

## Core responsibilities

1. **Agent creation** — Generate AI agent code with best practices
2. **Existing agent enhancement** — Update, fix, add features, add debug support
3. **Model selection** — Recommend and compare AI models
4. **Tracing** — Integrate tracing for debugging and monitoring
5. **Evaluation** — Assess agent performance and quality
6. **Deployment** — Go to production via Foundry

## Core principle

**Language:** Python by default unless specified. **.NET** is also supported.

**Microsoft Agent Framework:** Unified Python/.NET SDK for enterprise AI agents and workflows with type safety, checkpointing, and multi-agent orchestration.

## Toolbelt

| Category | Tool | Description |
| --- | --- | --- |
| Code generation | `aitk-get_agent_model_code_sample` | Basic snippets (agent, workflow, chat) |
| Code generation | `githubRepo` | Search `microsoft/agent-framework` for patterns |
| Code generation | `aitk-agent_as_server` | Best practices for HTTP server wrapping |
| Debugging | `aitk-add_agent_debug` | Dev tools and VS Code configs for Agent Inspector |
| Python env | `getPythonEnvironmentInfo`, `configurePythonEnvironment`, `installPythonPackage`, `getPythonExecutableCommand` | Manage Python environment and dependencies |
| Models | `aitk-get_ai_model_guidance` | Expert advice on model selection |
| Models | `aitk-list_foundry_models` | List deployed Foundry models |
| Operations | `aitk-get_tracing_code_gen_best_practices` | Tracing setup guidance |
| Operations | `aitk-evaluation_planner` | Evaluation strategy planning |
| Operations | `aitk-evaluation_agent_runner_best_practices` | Evaluation runner best practices |
| Operations | `aitk-get_evaluation_code_gen_best_practices` | Evaluation code guidance |

---

## Agent creation

**When to use:** User asks to "create", "scaffold", "start", or "build" a new agent or workflow.

### SDK setup

**Python** — Requires Python 3.10+. Pin the version to avoid breaking changes:

```bash
pip install agent-framework-azure-ai==1.0.0b260107 agent-framework-core==1.0.0b260107
```

**.NET** — Use `--prerelease`. Remind user in generated documentation:

```bash
dotnet add package Microsoft.Agents.AI.AzureAI --prerelease
dotnet add package Microsoft.Agents.AI.OpenAI --prerelease
dotnet add package Microsoft.Agents.AI.Workflows --prerelease
```

### Options

- **More samples:** If the scenario is specific, call `githubRepo` to search for samples before generating.
- **Minimal / test only:** Skip Agent-as-Server, Debug, and verification steps.
- **Deferred config:** Skip configuration and remind the user to update `.env` later.

### Creation workflow

```
- [ ] Gather context (Samples, Model, Server, Debug)
- [ ] Create implementation plan
- [ ] Select model & configure environment
- [ ] Implement code (Agent-as-Server pattern)
- [ ] Install dependencies
- [ ] Verify startup (Run-Fix loop)
- [ ] Documentation & Handoff
```

**Step 1 — Gather context.** Call tools from the Toolbelt. For standard new agent requests:

Required: `aitk-get_agent_model_code_sample`, `aitk-agent_as_server`, `aitk-add_agent_debug`, `aitk-get_ai_model_guidance`, `aitk-list_foundry_models`.

Recommended: `githubRepo` for advanced patterns (MCP, Multimodal, Assistants API, Responses API, Copilot Studio, Anthropic, Reflection, Switch-Case, Fan-out/Fan-in, Loop, Human-in-Loop).

**Step 2 — Create implementation plan.** Think through a detailed step-by-step plan. Output the high-level steps so the user knows what will happen.

**Step 3 — Select model & configure environment.** Decide on the model before coding — this determines the correct client/credential patterns.

- If the user hasn't specified a model, transition to the Model Selection capability.
- Create/update `.env` (do not overwrite existing variables):
  ```
  FOUNDRY_PROJECT_ENDPOINT=<project-endpoint>
  FOUNDRY_MODEL_DEPLOYMENT_NAME=<model-deployment-name>
  ```
- Always output what is configured, the file location, and how to change it later.

**Step 4 — Implement code.**

- **Server mode:** Implement the Agent-as-Server pattern (HTTP) unless "Minimal" was requested.
- **Debug:** Apply dev tools and add `.vscode/launch.json` and `.vscode/tasks.json` from `aitk-add_agent_debug`.
- **Patterns:** Use context from the gather step to structure the agent or workflow.

**Step 5 — Install dependencies.**

1. Generate/update `requirements.txt`.
2. Check/configure Python environment using `configurePythonEnvironment`.
3. Install packages using `installPythonPackage` or the terminal command for the correct executable.

**Step 6 — Verify startup (Run-Fix loop).**

Enter a run-fix loop: run → if unexpected error, fix → rerun → repeat until no startup/init error.

1. Run the main entrypoint (HTTP Server).
2. If startup fails: fix error → rerun.
3. If startup succeeds: stop server immediately.

Guardrails:
- DO perform a real run to catch startup errors early (static syntax check is not enough)
- DO clean up after verification — if you started the HTTP server, stop it
- DO ignore environment/auth/connection/timeout errors; focus only on startup/init errors
- Do NOT wait for user input
- Do NOT create separate test code or scripts
- Do NOT mock configuration

**Step 7 — Documentation & handoff.**

- Create/update `README.md`.
- Remind the user of next steps:
  - Debug / F5 to try and verify the app locally
  - Tracing setup to monitor and troubleshoot runtime issues

---

## Existing agent enhancement

**When to use:** User asks to "update", "fix", "add feature", "enhance", or "improve" an existing agent.

**Principles** (for Microsoft Agent Framework — do not change other SDKs unless explicitly asked):

1. **Context first:** Explore the codebase to understand existing architecture, patterns, and dependencies before making changes.
2. **Explore & gather:** Use the Toolbelt to gather context while respecting existing types.
3. **Respect tech stack:** Do not migrate or change other SDKs unless explicitly requested.
4. **Respect existing types:** Keep existing types like `*Client`, `*Credential`, etc. No migration unless requested.
5. **New features:** Follow the Creation Workflow (server mode, debug support).
6. **Verify:** Use Step 6 from the Creation Workflow to ensure no regressions.

---

## Model selection

**When to use:** User asks to "configure", "change", or "recommend" a model, or asks "which model" to use. Also triggered automatically during Agent Creation.

Use `aitk-get_ai_model_guidance` and `aitk-list_foundry_models`.

- For production-quality agents, recommend Foundry model(s).
- The user's existing model deployment is a quick start, but not necessarily the best choice — recommend based on user intent and model capabilities.
- Always explain your recommendation and show alternatives even if not yet deployed.
- If no Foundry project/model is available, recommend creating one via the Foundry extension.

---

## Tracing

**When to use:** User asks to "monitor", "trace", or improve "observability".

Use `aitk-get_tracing_code_gen_best_practices` to retrieve best practices, then apply them to instrument the code.

---

## Evaluation

**When to use:** User asks to "improve performance", "measure", or "evaluate" the agent.

1. **Planning first:** Use `aitk-evaluation_planner` to clarify metrics, test dataset, and runtime.
2. **Runner:** Use `aitk-evaluation_agent_runner_best_practices` for collecting responses from test datasets.
3. **Code:** Use `aitk-get_evaluation_code_gen_best_practices` for evaluation code generation.

---

## Deployment

**When to use:** User asks to "deploy", "publish", or "go production" with the agent.

1. Ensure the app is wrapped as an HTTP server (if not, use `aitk-agent_as_server` first).
2. Execute VS Code command: [Microsoft Foundry: Deploy Hosted Agent](azure-ai-foundry.commandPalette.deployWorkflow).
