<!-- hub-reference-banner -->
> **Reference file — part of the `content-ingestion-extraction` hub.** Formerly the standalone `using-plaud-mcp` skill.
> Sibling topics in this family are now reference files under the hubs (`content-ingestion-extraction`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: using-plaud-mcp
title: Plaud MCP Tools — Quick Start
description: |
  Quick-start reference for the Plaud MCP server: 11 tools, 4 prompts, and 12 skill resources for managing Plaud AI recordings, transcripts, and memory.
  TRIGGER: using plaud MCP tools; find_recordings or get_recording calls; transcribing a Plaud recording via MCP; memory_search or memory_ingest; checking processing status; getting a presigned audio URL; MCP elicitation for cost-gated Plaud operations; resolving a file_ref by row number or short prefix.
  SKIP: hardware setup or device pairing (use plaud-integration skill); CLI usage of the plaud command (use using-plaud-cli skill); bulk export or advanced ingestion patterns beyond basic memory_ingest (use bulk-operations skill); deep transcript analysis workflows (use transcription skill).
category: custom
version: "1.1.0"
updated: "2026-05-29"
keywords:
  - Plaud
  - Plaud MCP
  - recordings
  - transcripts
  - memory search
  - find_recordings
  - transcribe
  - memory_ingest
  - get_content
  - ElevenLabs
  - MCP elicitation
  - file_ref
  - presigned URL
when_to_use:
  - "use Plaud MCP tools"
  - "find a Plaud recording"
  - "transcribe a recording via MCP"
  - "search Plaud memory"
  - "ingest recordings into memory"
  - "get Plaud processing status"
  - "resolve a file_ref by row number"
  - "daily briefing from Plaud"
  - "analyze a Plaud transcript"
  - "bulk summarize recent recordings"
related_skills:
  - content-ingestion-extraction
  - using-plaud-cli
---

# Plaud MCP Tools — Quick Start

Quick-start reference for the Plaud MCP server: 11 tools, 4 prompts, and 12 skill resources for managing Plaud AI recordings, transcripts, and memory. Each tool returns structured JSON with `next_steps` guidance.

## Prerequisites

- MCP server configured in `.mcp.json` — authenticated via GitHub OAuth (automatic on first connect)
- `ELEVENLABS_API_KEY` env var — required only for ElevenLabs transcription

## Tool quick reference

| Tool | Purpose | Domain skill |
|---|---|---|
| `find_recordings` | Browse and filter recordings | recording-search |
| `get_recording` | Full recording metadata | recording-search |
| `get_audio_url` | Presigned download URL (5 min expiry) | recording-search |
| `transcribe` | ElevenLabs transcript (cached or live) | transcription |
| `get_content` | AI summary or meeting notes | transcription |
| `trigger_processing` | Plaud server-side AI processing | transcription |
| `get_account_info` | Account profile and membership | — |
| `get_processing_status` | AI processing queue status | — |
| `list_languages` | Supported language codes | transcription |
| `memory_search` | Semantic search across transcripts | memory-search |
| `memory_ingest` | Bulk index recordings into memory | bulk-operations |

## Prompts

MCP clients can discover and invoke these workflow templates:

| Prompt | Purpose | Parameters |
|---|---|---|
| `daily_briefing` | Status overview — recent recordings, queue, memory | None |
| `analyze_transcript` | Deep analysis of a single recording | `file_ref` |
| `search_memory` | Guided dual-path search | `query`, `time_range` (default `"7d"`) |
| `bulk_summary` | Summarize multiple recordings | `count` (default `5`) |

## File ID resolution

All tools accepting `file_ref` support three formats:

| Format | Example | Notes |
|---|---|---|
| Full ID | `4f757af256ecba4fab502739c122dc78` | Always unambiguous |
| Short prefix | `4f75` | Must be unique among recent results |
| Row number | `3` | From last `find_recordings` call; persists within a session |

Call `find_recordings` once, then use row numbers in subsequent `transcribe`, `get_content`, and similar calls.

## Cost safety

Two tools prompt for confirmation via MCP elicitation before consuming external credits:

- `transcribe` — asks before calling ElevenLabs API (cached transcripts are free)
- `memory_ingest(source="all")` — asks before bulk operations that may trigger transcriptions

Older MCP clients that do not support elicitation skip the prompt and proceed directly.

## Skills as MCP resources

The server exposes plugin skills as MCP resources — any client can discover and read them without installing the plugin:

```
skill://recording-search/SKILL.md
skill://memory-search/SKILL.md
skill://transcription/SKILL.md
skill://bulk-operations/SKILL.md
skill://using-plaud-mcp/SKILL.md
skill://using-plaud-cli/SKILL.md
```

## Setup

The MCP server is configured in the plugin's `.mcp.json`. Authentication uses GitHub OAuth — the MCP client handles the flow automatically on first connect. See the plugin README for details.

## See also

| Skill | When to use it |
|---|---|
| **plaud-integration** | Hardware, AI engine, Developer Platform API, cloud/export, privacy — the deep product reference |
| **recording-search** | Browsing, filtering, and metadata retrieval patterns |
| **memory-search** | Semantic search across indexed transcripts |
| **transcription** | ElevenLabs transcription with cost awareness |
| **bulk-operations** | Batch ingestion, auto-polling, and export |
| **using-plaud-cli** | CLI usage guide (`plaud` command) |
