<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.** Formerly the standalone `prompt-lookup` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: prompt-lookup
version: "1.2.0"
updated: "2026-05-31"
description: >
  Prompt discovery, retrieval, management, versioning, and optimization. Uses the
  prompts.chat MCP server for search/get/improve. Covers prompt marketplaces
  (PromptBase, FlowGPT, prompts.chat, Anthropic Prompt Library), template management,
  versioning strategies, A/B testing, evaluation frameworks, and registry integration
  (LangChain Hub, PromptLayer, Langfuse, Agenta).
  TRIGGER: user asks for a prompt template, wants to search for prompts, mentions
  prompts.chat or PromptBase or FlowGPT, needs help with prompt versioning or A/B
  testing, wants to manage prompt templates across a team, or needs to integrate
  prompts into CI/CD or LLMOps pipelines.
  SKIP: writing or refining a single prompt from scratch
  (→ ai-agent-engineering (references/prompt-engineering.md));
  optimizing an existing prompt for production (→ prompt-helper-optimizer or
  prompt-deep-optimizer); building or reviewing a Claude Code skill (→ claude-code-skills).
whenToUse:
  - "find me a prompt template for X"
  - "search for prompts on prompts.chat"
  - "I need a prompt for SQL / code review / summarization"
  - "how do I version my team's prompts?"
  - "set up A/B testing for our prompts"
  - "manage prompts across my team"
  - "what prompts are on PromptBase or FlowGPT?"
  - "integrate prompts into our CI/CD pipeline"
  - "build a team prompt library"
  - "what prompt template schema should I use?"
related_skills:
  - ai-agent-engineering
  - prompt-helper-optimizer
  - prompt-deep-optimizer
  - claude-code-skills
triggers:
  - prompt template
  - prompt library
  - prompt marketplace
  - prompt registry
  - prompt versioning
  - prompt A/B testing
  - prompt management
  - prompts.chat
  - PromptBase
  - FlowGPT
  - prompt discovery
  - prompt optimization tool
  - prompt metadata
---

# Prompt Lookup

When the user needs AI prompts, prompt templates, or wants to improve existing prompts, use the prompts.chat MCP server. This skill covers prompt discovery, management, versioning, and optimization ecosystem.

## When to Use

- Searching for prompt templates ("Find me a code review prompt")
- Retrieving a specific prompt ("Get prompt XYZ")
- Improving an existing prompt ("Make this prompt better")
- Researching prompt marketplaces or registries
- Setting up prompt versioning or A/B testing
- Building a team prompt library
- Integrating prompts into CI/CD pipelines
- Understanding prompt metadata schemas

## When NOT to Use

- **Writing a single prompt from scratch** → `ai-agent-engineering` (references/prompt-engineering.md)
- **Optimizing an existing prompt for production** → `prompt-helper-optimizer` or `prompt-deep-optimizer`
- **Building a Claude Code skill** → `claude-code-skills`
- **General LLM API questions** → `ai-agent-engineering` (references/llm-models.md)

---

## MCP Tools (prompts.chat)

| Tool | Parameters | Use when |
|------|-----------|----------|
| `search_prompts` | `query`, `limit` (default 10, max 50), `type` (TEXT/STRUCTURED/IMAGE/VIDEO/AUDIO), `category`, `tag` | Finding prompts by keyword |
| `get_prompt` | `id` | Retrieving a specific prompt; fills `${variable}` slots |
| `improve_prompt` | `prompt`, `outputType` (text/image/video/sound), `outputFormat` (text/structured_json/structured_yaml) | Enhancing a user's existing prompt |

### Search → Retrieve workflow

```
1. Call search_prompts(query: "<user's need>", limit: 10)
2. Present results: title, description, author, category, tags, link
3. User selects → call get_prompt(id: "<id>")
4. If prompt has ${variable} slots, prompt user to fill required ones
5. If user wants improvement → call improve_prompt(prompt: "<text>")
```

### Error handling

| Error | Fix |
|-------|-----|
| Zero search results | Broaden query (remove filters, try synonyms); try Anthropic Prompt Library or PromptBase |
| `get_prompt` returns empty | Confirm ID; prompt may be removed. Suggest searching for alternatives |
| `improve_prompt` timeout | Return original; suggest manual refinement using `prompt-engineering` skill |
| Rate limiting | Wait briefly, retry once; if persistent, use manual recommendations from marketplace table |

---

## Prompt Marketplaces

| Platform | Model | Strengths |
|----------|-------|-----------|
| **prompts.chat** | Free open-source | Curated "Act as..." collection; MCP-accessible |
| **Anthropic Prompt Library** | Free (60+ prompts) | Official Claude prompts; demonstrates XML-structured formatting |
| **PromptBase** | Paid marketplace (500K+ listings) | Largest general marketplace; coding, SQL, automation prompts |
| **FlowGPT** | Free / donation-based | Community ChatGPT workflow prompts; budget-friendly |
| **PromptHero** | Free with premium | Image generation prompts (Stable Diffusion, Midjourney, DALL-E) |
| **God of Prompt** | Paid bundles | Premium packs for business, marketing, coding |

**When to recommend:**
- Quick free lookup → prompts.chat (MCP-integrated), Anthropic Prompt Library
- Production-quality paid → PromptBase, God of Prompt
- Image generation → PromptHero
- Claude-specific patterns → Anthropic Prompt Library first

---

## Prompt Template Schema

A well-structured template contains:

```yaml
id: "code-review-v2"
title: "Code Review Assistant"
version: "2.1.0"
author: "team-platform"
description: "Reviews code for bugs, style, and security"
model_compatibility:
  - claude-sonnet-4-6-20250514
  - gpt-4.1
variables:
  - name: language
    type: string
    required: true
  - name: code
    type: string
    required: true
  - name: focus_areas
    type: array
    required: false
    default: ["bugs", "style", "security"]
tags: ["coding", "review", "quality"]
template: |
  <role>You are a senior {{language}} developer performing a code review.</role>
  <instructions>Review the following code for: {{focus_areas}}</instructions>
  <code>{{code}}</code>
```

### Variable Interpolation Patterns

| Pattern | Syntax | Used by |
|---------|--------|---------|
| Mustache/Handlebars | `{{variable}}` | LangChain, PromptHub |
| Dollar-brace | `${variable}` | prompts.chat, CLI tools |
| F-string | `{variable}` | Python LangChain, DSPy |
| Jinja2 | `{% if %}...{% endif %}` | Ansible, some LLMOps platforms |

---

## Anti-Patterns

### Discovery
| Anti-Pattern | Fix |
|-------------|-----|
| Copy-paste from internet without provenance | Use a registry; track source and version |
| Team members maintain private prompt collections | Centralize in shared registry with search |
| Assuming paid = quality without testing | Always eval against your use case before adopting |
| Dozens of similar prompts with minor variations | Consolidate; use variables and conditional logic |

### Engineering
| Anti-Pattern | Fix |
|-------------|-----|
| Old examples contradict current instructions | Audit example set whenever instructions change |
| 8–10+ distinct instructions in one prompt | Split into chained prompts or structured sections |
| Complex workflow crammed into one massive prompt | Decompose into prompt chains or agent steps |
| Not tracking which model version a prompt targets | Record `model_compatibility` in metadata; re-eval on model upgrades |

---

## Claude-Specific Patterns

Claude 4.x follows instructions literally — it does exactly what you ask, nothing more. Prompts designed for Claude 3.x may need updating.

```xml
<system>You are a {{role}} assistant.</system>
<instructions>{{task_instructions}}</instructions>
<context>{{relevant_context}}</context>
<examples>
  <example>
    <input>{{example_input}}</input>
    <output>{{example_output}}</output>
  </example>
</examples>
<input>{{user_input}}</input>
```

**Anthropic Prompt Library categories:** Coding, Writing, Analysis, Business, Education, Creative, Operations, Research.

---

## Discovery Workflow

```
1. User describes need
2. search_prompts via MCP
3. If insufficient results, check:
   a. Anthropic Prompt Library (Claude-specific)
   b. PromptBase (paid, high quality)
   c. FlowGPT (free, community)
   d. LangChain Hub (developer-focused)
4. Evaluate found prompt against user's model and use case
5. Improve via improve_prompt (prompts.chat) or tam_optimize_prompt (mdb_context_hub MCP)
6. Save to local registry for reuse via tam_save_prompt
```

## Building a Team Prompt Library

1. Audit: collect all prompts scattered across code, Slack, docs
2. Standardize: adopt a metadata schema with id, version, variables, model_compatibility
3. Centralize: choose Git-native, registry-first, or hybrid pattern
4. Add evaluation: require eval scores before promoting to production
5. Enable discovery: tag, categorize, make searchable
6. Govern updates: version control with review process
7. Schedule maintenance: quarterly review for staleness and model compatibility

For extended platform comparisons, versioning workflows, A/B testing procedures, and integration code patterns, see `references/prompt-management-platforms.md`.
