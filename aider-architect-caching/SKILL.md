---
name: aider-architect-caching
description: >-
  Deep research on Aider's architect mode and token caching optimizations. Explains context windows, prompt caching, and the dual-model planning vs code execution workflow.
  TRIGGER: "explain aider architect mode", "aider token caching", "aider context window".
  SKIP: General Aider usage → aider-expert; executing shell commands → aider-run-execution.
whenToUse:
  - "explain Aider's architect mode"
  - "how does Aider caching work"
  - "aider context window management"
  - "optimize Aider token usage"
triggers:
  - aider architect
  - aider caching
  - aider context
version: 1.0.0
category: documentation
updated: 2026-08-18
model: claude-haiku-4-5
effort: low
keywords:
  - aider
  - architect mode
  - prompt caching
  - context window
  - dual-model
related_skills:
  - aider-expert
  - aider-run-execution
---

# Aider Architect Mode and Caching Optimizations

Aider's Architect mode introduces a two-step dual-model collaboration to improve the quality of complex code changes by separating reasoning from execution.

## When not to use

- General Aider usage questions → `aider-expert`
- Running shell commands via Aider → `aider-run-execution`

## Architect Mode: Planning vs Code Execution

### The Two Roles

1. **The Architect (Reasoning):** Designs the solution, focusing on logic and architecture without strict formatting constraints.
2. **The Editor (Execution):** Translates the Architect's plan into precise, valid code edits.

### Code Mode vs. Architect Mode

- **Code Mode (Default):** The model performs both reasoning and editing in a single turn. Faster and better for simple tasks.
- **Architect Mode:** Separates tasks to reduce cognitive load. A powerful reasoning model (e.g., OpenAI o1) acts as the Architect, while a highly precise model acts as the Editor.

### Practical Usage

- Trigger via `/architect` in the terminal or launch with `--architect`.
- Configure different models for each role to optimize for performance and cost.

## Token Caching & Context Window Management

Aider uses prompt caching to optimize context window usage, reduce costs, and improve response times by avoiding redundant token processing.

### What Aider Caches

Aider organizes chat history to cache static data:

1. **System Prompts:** Foundational instructions.
2. **Read-Only Files:** Static files added with `--read` or `/read-only`.
3. **Repository Map:** A concise representation of the project structure.
4. **Editable Files:** Files currently being modified.

### Enabling and Managing Caching

- **Enable:** Use `--cache-prompts` or set `cache-prompts: true` in the config (Anthropic and DeepSeek models).
- **Cache Keepalive:** Provider APIs often have a short TTL (e.g., 5 mins for Anthropic). Use `--cache-keepalive-pings N` to send periodic pings to keep the cache warm.

### Context Window Best Practices

1. **Be Selective:** Add only necessary files. Use `/drop` to remove irrelevant files.
2. **Monitor Usage:** Track usage and identify hogs using `/tokens`.
3. **Clear History:** Use `/clear` to start fresh if a conversation becomes sluggish, retaining essential context.
4. **Mental Model:** Treat the context window as a fixed-size cache. Keeping core context concise ensures the model doesn't evict crucial information.
