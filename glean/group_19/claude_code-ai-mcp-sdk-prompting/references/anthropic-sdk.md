<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.** Formerly the standalone `anthropic-sdk` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: anthropic-sdk
version: "1.1.0"
updated: "2026-05-29"
description: >
  Official TypeScript/Node.js client for the Anthropic Messages API. Covers messages.create,
  streaming (MessageStream + raw SSE), tool use / function calling, prompt caching
  (cache_control, ephemeral TTL, cache hit metrics), extended/adaptive thinking
  (effort parameter), vision (URL + base64), Batch API (client.beta.messages.batches),
  citations, and error handling (RateLimitError, APIStatusError 529, APIConnectionError).
  TRIGGER: code imports @anthropic-ai/sdk or anthropic; user asks about Claude API
  features (streaming, tool use, caching, thinking, batch, vision, citations); user needs
  to handle Anthropic-specific error types or tune model costs; user migrates between
  Claude model versions.
  SKIP: code imports openai / @azure/openai / @google-cloud/vertexai (use that provider's
  skill); provider-neutral or model-agnostic code with no Anthropic-specific calls;
  Bedrock/Vertex AI transport wrappers (same API shape, different auth — check those docs).
tags: [anthropic, claude, sdk, ai, llm, nodejs, typescript]
keywords: [anthropic, claude, "@anthropic-ai/sdk", messages.create, messages.stream, tool_use, cache_control, batch, message_batches, extended thinking, vision, citations, rate_limit_error, overloaded_error, APIError, budget_tokens, prompt caching]
related_skills:
  - llm-models
  - prompt-engineering
  - rag-architecture
---

# Anthropic Claude SDK (Node.js / TypeScript)

## Overview

`@anthropic-ai/sdk` is the official TypeScript/Node.js client for the Anthropic Messages API (Claude models). Node 18+ required.

## When not to use

- Code imports `openai`, `@azure/openai`, `@google-cloud/vertexai`, or another provider's SDK — use that provider's skill instead.
- Bedrock or Vertex AI wrappers (`@anthropic-ai/bedrock-sdk`) — same API shape but different auth/endpoint; check the Bedrock/Vertex docs for transport differences.
- Provider-neutral or model-agnostic code with no Anthropic-specific calls.

## Quick Setup

```bash
npm install @anthropic-ai/sdk
export ANTHROPIC_API_KEY=sk-ant-...
```

```typescript
import Anthropic from '@anthropic-ai/sdk';
const client = new Anthropic();           // picks up ANTHROPIC_API_KEY automatically
// const client = new Anthropic({ apiKey: process.env.ANTHROPIC_API_KEY });
```

## Current Model IDs

| Model | ID | Best for |
|---|---|---|
| Opus 4.7 | `claude-opus-4-7` | Complex reasoning, long-horizon agentic work |
| Sonnet 4.6 | `claude-sonnet-4-6` | Coding, agents, enterprise (frontier at scale) |
| Haiku 4.5 | `claude-haiku-4-5` | Speed, near-frontier, low cost |

Always pin the full ID string (e.g., `claude-sonnet-4-6`) — never use `latest`.

## Messages API

```typescript
const msg = await client.messages.create({
  model: 'claude-sonnet-4-6',
  max_tokens: 1024,
  system: 'You are a helpful assistant.',
  messages: [
    { role: 'user', content: 'Explain prompt caching in one paragraph.' }
  ],
});
// msg.content[0].type === 'text'
// msg.content[0].text
// msg.stop_reason: 'end_turn' | 'max_tokens' | 'tool_use' | 'stop_sequence'
```

Multi-turn: append each assistant response to `messages` before the next user turn.

## Streaming

```typescript
// Option A: MessageStream (helper methods, builds final message)
const stream = await client.messages.stream({
  model: 'claude-sonnet-4-6',
  max_tokens: 1024,
  messages: [{ role: 'user', content: 'Tell me a story.' }],
});
stream.on('text', (text) => process.stdout.write(text));
const final = await stream.finalMessage();

// Option B: raw async iterable (lower memory, no accumulation)
const raw = await client.messages.create({
  model: 'claude-sonnet-4-6',
  max_tokens: 1024,
  stream: true,
  messages: [{ role: 'user', content: 'Tell me a story.' }],
});
for await (const event of raw) {
  if (event.type === 'content_block_delta' && event.delta.type === 'text_delta') {
    process.stdout.write(event.delta.text);
  }
}
```

Key stream event types: `message_start`, `content_block_start`, `content_block_delta`, `content_block_stop`, `message_delta`, `message_stop`.

**Streaming with tool use:** when Claude calls a tool during a stream, tool input arrives as `input_json_delta` events (not `text_delta`). Check `event.delta.type` to distinguish (extend Option B's loop):

```typescript
// Extends Option B — raw is the async iterable from client.messages.create({ stream: true, ... })
let partialJson = '';
for await (const event of raw) {
  if (event.type === 'content_block_delta') {
    if (event.delta.type === 'text_delta') {
      process.stdout.write(event.delta.text);
    } else if (event.delta.type === 'input_json_delta') {
      partialJson += event.delta.partial_json;  // accumulate partial tool input JSON
    }
  }
}
// After stream ends, parse: JSON.parse(partialJson) gives the tool's input object
```

## Tool Use / Function Calling

```typescript
const tools: Anthropic.Tool[] = [{
  name: 'get_weather',
  description: 'Get current temperature for a city.',
  input_schema: {
    type: 'object' as const,
    properties: { city: { type: 'string', description: 'City name' } },
    required: ['city'],
  },
}];

const resp = await client.messages.create({
  model: 'claude-sonnet-4-6',
  max_tokens: 1024,
  tools,
  // tool_choice: { type: 'auto' }          // default
  // tool_choice: { type: 'any' }           // must use one of the tools
  // tool_choice: { type: 'tool', name: 'get_weather' }  // force specific tool
  messages: [{ role: 'user', content: 'What is the weather in Paris?' }],
});

// Handle tool call
if (resp.stop_reason === 'tool_use') {
  const toolUse = resp.content.find(b => b.type === 'tool_use') as Anthropic.ToolUseBlock;
  const result = await callMyFunction(toolUse.name, toolUse.input as { city: string });

  // Send result back — include the full assistant turn + tool result
  const finalResp = await client.messages.create({
    model: 'claude-sonnet-4-6',
    max_tokens: 1024,
    tools,
    messages: [
      { role: 'user', content: 'What is the weather in Paris?' },
      { role: 'assistant', content: resp.content },
      {
        role: 'user',
        content: [{
          type: 'tool_result',
          tool_use_id: toolUse.id,
          content: JSON.stringify(result),
        }],
      },
    ],
  });
  console.log(finalResp.content[0]);
}
```

## Prompt Caching

Caching cuts costs up to 90% and latency up to 80%. Cache reads cost 10% of standard input price.

```typescript
// system must be an array of content blocks (not a plain string) to use cache_control
const resp = await client.messages.create({
  model: 'claude-sonnet-4-6',
  max_tokens: 1024,
  system: [{
    type: 'text',
    // Place cache_control on the LAST block of the stable prefix you want cached.
    // Everything before this marker is included in the cache.
    text: 'You are an expert on a 50,000-word legal document...\n<document>...full text...</document>',
    cache_control: { type: 'ephemeral' },
  }],
  messages: [{ role: 'user', content: 'Summarize section 3.' }],
});
// resp.usage.cache_creation_input_tokens — tokens written to cache (1.25x cost)
// resp.usage.cache_read_input_tokens     — tokens read from cache (0.10x cost)
```

**Key rules:**
- Minimum cacheable prefix: **1,024 tokens** (Sonnet/Haiku); **4,096 tokens** (Opus).
- Default TTL: 5 minutes (refreshed on each read). Extended TTL: `{ type: 'ephemeral', ttl: '1h' }`.
- Cache is workspace-scoped (as of Feb 2026).
- Tools and system prompts are the best caching targets; cache tool definitions when using many tools.

## Extended Thinking

Newer models use `adaptive` thinking (replaces deprecated `budget_tokens`).

```typescript
// Adaptive thinking (Opus 4.6+, Sonnet 4.6+)
const resp = await client.messages.create({
  model: 'claude-sonnet-4-6',
  max_tokens: 16000,
  thinking: {
    type: 'adaptive',
    effort: 'high',       // 'low' | 'medium' | 'high'
    // display: 'omitted' // omit thinking text from stream while preserving signature for multi-turn
  },
  messages: [{ role: 'user', content: "Sketch a proof of Fermat's last theorem." }],
});

// Response includes thinking blocks + text blocks
for (const block of resp.content) {
  if (block.type === 'thinking') {
    console.log('[thinking]', block.thinking);
  } else if (block.type === 'text') {
    console.log(block.text);
  }
}
```

**Notes:**
- `effort: 'high'` elicits more thinking; complex queries also trigger more thinking automatically.
- Set `display: 'omitted'` to suppress thinking text in streaming responses while keeping the opaque signature needed for multi-turn continuity.
- Thinking blocks CAN be cached in prior assistant turns, counting as input tokens on cache read.
- For large thinking budgets (>32k tokens), use the Batch API to avoid network timeouts.

## Vision

```typescript
import Anthropic from '@anthropic-ai/sdk';
import fs from 'fs';

const client = new Anthropic();

// URL source (preferred — avoids re-sending bytes on every multi-turn)
const respUrl = await client.messages.create({
  model: 'claude-sonnet-4-6',
  max_tokens: 1024,
  messages: [{
    role: 'user',
    content: [
      { type: 'image', source: { type: 'url', url: 'https://example.com/chart.png' } },
      { type: 'text', text: 'Describe this chart.' },
    ],
  }],
});

// Base64 source (for private/local images not reachable by URL)
const imageData = fs.readFileSync('chart.png').toString('base64');
const respB64 = await client.messages.create({
  model: 'claude-sonnet-4-6',
  max_tokens: 1024,
  messages: [{
    role: 'user',
    content: [
      {
        type: 'image',
        source: { type: 'base64', media_type: 'image/png', data: imageData },
      },
      { type: 'text', text: 'What does this show?' },
    ],
  }],
});
```

Supported formats: `image/jpeg`, `image/png`, `image/gif`, `image/webp`. Base64 images are re-sent in full on every turn — prefer URL or the Files API for multi-turn conversations.

## Batch API

Process up to 10,000 requests asynchronously at ~50% cost reduction. Results available within 24h.
The Message Batches API is under the `beta` namespace: `client.beta.messages.batches`.

```typescript
// 1. Create batch
const batch = await client.beta.messages.batches.create({
  requests: [
    {
      custom_id: 'req-001',
      params: {
        model: 'claude-sonnet-4-6',
        max_tokens: 512,
        messages: [{ role: 'user', content: 'Summarize: ...' }],
      },
    },
    {
      custom_id: 'req-002',
      params: {
        model: 'claude-sonnet-4-6',
        max_tokens: 512,
        messages: [{ role: 'user', content: 'Translate: ...' }],
      },
    },
  ],
});
console.log('Batch ID:', batch.id);

// 2. Poll until done (results typically available in minutes to hours)
let status = batch;
while (status.processing_status !== 'ended') {
  await new Promise(r => setTimeout(r, 60_000));   // poll every 60s
  status = await client.beta.messages.batches.retrieve(batch.id);
}

// 3. Stream results
for await (const entry of await client.beta.messages.batches.results(batch.id)) {
  if (entry.result.type === 'succeeded') {
    console.log(entry.custom_id, entry.result.message.content[0]);
  } else {
    console.error(entry.custom_id, entry.result.type, entry.result.error);
  }
}
```

Combine with prompt caching for maximum cost savings.

## Citations

Cite sources from grounded documents. `cited_text` does NOT count toward output tokens.

```typescript
const resp = await client.messages.create({
  model: 'claude-sonnet-4-6',
  max_tokens: 1024,
  messages: [{
    role: 'user',
    content: [
      {
        type: 'document',
        source: {
          type: 'content',
          content: [
            { type: 'text', text: 'Anthropic was founded in 2021...' },
            { type: 'text', text: 'The company focuses on AI safety...' },
          ],
        },
        title: 'Anthropic Overview',
        citations: { enabled: true },
      },
      { type: 'text', text: 'When was Anthropic founded?' },
    ],
  }],
});

// Response contains text blocks with citation arrays
for (const block of resp.content) {
  if (block.type === 'text') {
    console.log(block.text);
    // block.citations[].document_title, .cited_text, .start_char_index
  }
}
```

Citations must be enabled on all documents or none within a single request.

## Error Handling

The SDK auto-retries 429 and 5xx errors up to 2 times by default. Increase with `new Anthropic({ maxRetries: 5 })`. For custom backoff control (e.g., jitter, per-error-type delays), disable SDK retries and implement your own:

```typescript
import Anthropic from '@anthropic-ai/sdk';

// Pass maxRetries: 0 to disable SDK auto-retry and control backoff yourself
const client = new Anthropic({ maxRetries: 0 });

async function callWithRetry(params: Anthropic.MessageCreateParamsNonStreaming) {
  const maxAttempts = 4;
  let lastErr: unknown;
  for (let attempt = 0; attempt < maxAttempts; attempt++) {
    try {
      return await client.messages.create(params);
    } catch (err) {
      lastErr = err;
      if (err instanceof Anthropic.RateLimitError) {
        // 429: exponential backoff
        await new Promise(r => setTimeout(r, Math.pow(2, attempt) * 1000));
      } else if (err instanceof Anthropic.APIStatusError && err.status === 529) {
        // 529 overloaded_error: exponential backoff
        await new Promise(r => setTimeout(r, Math.pow(2, attempt) * 2000));
      } else if (err instanceof Anthropic.APIConnectionError) {
        // Network error: short exponential backoff
        await new Promise(r => setTimeout(r, Math.pow(2, attempt) * 500));
      } else {
        // Non-retryable: auth errors, invalid requests, TypeErrors, etc.
        throw err;
      }
    }
  }
  throw lastErr; // rethrow after exhausting retries
}
```

| Error class | HTTP | When |
|---|---|---|
| `Anthropic.AuthenticationError` | 401 | Bad API key |
| `Anthropic.PermissionDeniedError` | 403 | Insufficient permissions |
| `Anthropic.NotFoundError` | 404 | Model/resource not found |
| `Anthropic.RateLimitError` | 429 | Tokens/requests per minute exceeded |
| `Anthropic.APIStatusError` (status=529) | 529 | API overloaded — transient |
| `Anthropic.APIConnectionError` | — | Network failure |
| `Anthropic.APIError` | any | Base class for all API errors |

Ramp up traffic gradually to avoid acceleration limits (sharp usage spikes trigger 429s).

## Cost Optimization Checklist

| Technique | Savings | When |
|---|---|---|
| Prompt caching (system / tools) | Up to 90% on cache hits | Repeated large contexts |
| Batch API | ~50% vs synchronous | Non-time-sensitive bulk work |
| Haiku instead of Sonnet/Opus | 5–20x cheaper | Simple tasks, high volume |
| Right-size `max_tokens` | Reduces OTPM bill | You know expected output length |
| Cache tool definitions | Cache read price | Many tools repeated across turns |
| Extended TTL (`ttl: '1h'`) | Fewer cache misses | Sessions >5 min |

## Common Mistakes

| Mistake | Fix |
|---|---|
| Using `client.messages.batches` | Batch API is under `client.beta.messages.batches` |
| Sending `base64` images every multi-turn | Use URL source or Files API |
| Ignoring `stop_reason === 'tool_use'` | Always check stop reason before parsing response text |
| Caching tokens below minimum | 1,024 (Sonnet/Haiku) or 4,096 (Opus) minimum; smaller prompts aren't cached |
| Using deprecated `budget_tokens` | Use `thinking: { type: 'adaptive', effort: '...' }` on Opus 4.6+, Sonnet 4.6+ |
| Polling batch results with no delay | Poll every 60s minimum |
| Not checking `cache_read_input_tokens` | Verify caching works by inspecting `resp.usage` |
| Plain string in `system` with `cache_control` | `system` must be an array of content blocks to attach `cache_control` |

## Sources

- [Anthropic SDK TypeScript README](https://github.com/anthropics/anthropic-sdk-typescript/blob/main/README.md)
- [Messages API docs](https://docs.anthropic.com/en/api/messages-examples)
- [Streaming docs](https://docs.anthropic.com/en/docs/build-with-claude/streaming)
- [Tool use docs](https://docs.anthropic.com/en/docs/agents-and-tools/tool-use/implement-tool-use)
- [Prompt caching docs](https://docs.anthropic.com/en/docs/build-with-claude/prompt-caching)
- [Extended thinking docs](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
- [Vision docs](https://docs.anthropic.com/en/docs/build-with-claude/vision)
- [Batch processing docs](https://docs.anthropic.com/en/docs/build-with-claude/batch-processing)
- [Citations docs](https://docs.anthropic.com/en/docs/build-with-claude/citations)
- [Errors docs](https://docs.anthropic.com/en/api/errors)
- [Models overview](https://docs.anthropic.com/en/docs/about-claude/models/overview)
