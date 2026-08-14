<!-- hub-reference-banner -->
> **Reference file — part of the `content-ingestion-extraction` hub.** Formerly the standalone `plaud-integration` skill.
> Sibling topics in this family are now reference files under the hubs (`content-ingestion-extraction`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: plaud-integration
version: 1.1.1
updated: 2026-05-31
description: >
  Plaud AI recorder ecosystem — hardware devices (Note, Note Pro, NotePin,
  NotePin S specs and selection guide), Plaud Intelligence AI engine, Developer
  Platform API (OAuth 2.0, API keys, regional endpoints), official and community
  MCP servers, CLI, AutoFlow/Zapier automation, privacy/compliance (HIPAA, SOC 2,
  GDPR), Chrome extension integration patterns, and anti-patterns.
  TRIGGER: working with Plaud recording devices; fetching Plaud transcripts or
  summaries via API or MCP; building Plaud-connected Chrome extension features;
  configuring the Plaud MCP server; comparing Plaud hardware devices; questions
  about Plaud subscription plans, privacy policy, or compliance certifications.
  SKIP: general meeting transcription not involving Plaud hardware or API (use
  granola-transcription for Granola); Plaud MCP tool-call usage details (use
  using-plaud-mcp); audio file transcription from non-Plaud sources.
category: developer
tags: [plaud, recording, transcription, mcp, api, chrome-extension, hardware, ai-notes]
related_skills: [granola-transcription, using-plaud-mcp, chrome-extension-expert, ai-agent-engineering]
whenToUse:
  - "Which Plaud device should I buy for phone call recording?"
  - "How do I authenticate with the Plaud Developer Platform API?"
  - "Set up Plaud MCP server in Claude"
  - "Plaud bearer token extraction for community MCP"
  - "Build Chrome extension that ingests Plaud recordings"
  - "Plaud AutoFlow trigger setup"
  - "Plaud API regional endpoints"
  - "Is Plaud HIPAA compliant?"
  - "Plaud NotePin vs NotePin S differences"
  - "Zapier integration with Plaud"
whenNotToUse:
  - "Granola meeting notes (use granola-transcription)"
  - "Plaud MCP tool usage details (use using-plaud-mcp)"
  - "Compliance posture audit of Plaud data flows (use security-compliance-auditor)"
---

# Plaud Integration

Plaud is the world's leading AI note-taking hardware brand, combining purpose-built recording devices with an AI-powered software ecosystem for transcription, summarization, and meeting intelligence. The product line spans card-form-factor recorders (Note series) and wearable capsules (NotePin series), all paired with the Plaud App (iOS/Android), Plaud Web, Plaud Desktop, a Developer Platform, MCP server, and CLI.

---

## 1. Hardware Devices

### 1.1 Plaud Note (Original)

The card-sized recorder designed for planned recording sessions -- meetings, calls, lectures.

| Spec | Value |
|------|-------|
| Form factor | Credit-card sized, ultra-slim |
| Thickness | 0.12 in (2.99 mm) |
| Weight | 1.06 oz (30 g) |
| Battery | 400 mAh; 30 hours continuous recording; 60-day standby; 2-hour charge |
| Storage | 64 GB local |
| Microphones | 2 MEMS + 1 VPU (Vibration Processing Unit) |
| Voice range | 9.84 ft (3 m) |
| Connectivity | BLE 5.2, Wi-Fi |
| Recording modes | Dual-mode -- in-person and phone call (manual switch) |
| Call recording | Attach to phone back; VPU captures speaker vibrations for both sides |
| Price | $159 USD |
| Accessories | Magnetic case, magnetic ring, charging cable |

### 1.2 Plaud Note Pro

The flagship card recorder with expanded mic array, AMOLED display, and Smart Dual Mode.

| Spec | Value |
|------|-------|
| Dimensions | 85.6 x 54.1 x 2.99 mm (3.37 x 2.13 x 0.12 in) |
| Weight | 1.06 oz (30 g) |
| Battery | 500 mAh; 30 h (Enhance mode) / 50 h (Endurance mode); 60-day standby |
| Storage | 64 GB local + unlimited cloud |
| Microphones | 4 MEMS + 1 VPU; AI beamforming |
| Voice range | 16.4 ft / 5 m (Enhance) ; 9.84 ft / 3 m (Endurance) |
| Connectivity | Bluetooth BLE 5.4, dual-band Wi-Fi (2.4/5 GHz) |
| Display | 0.95" AMOLED, up to 600 nits, Gorilla Glass |
| Recording modes | Smart Dual Mode -- auto-detects phone call vs in-person |
| Materials | Aluminum alloy + Corning Gorilla Glass |
| Price | $189 USD |
| Accessories | Magnetic case, magnetic ring, charging cable |

**Key differentiators from Note:** InstantView AMOLED display, 4-mic array with 5 m range, Smart Dual Mode (no manual switching), dual-band Wi-Fi, Endurance mode for 50-hour battery life.

### 1.3 Plaud NotePin

The wearable AI memory capsule for hands-free, all-day in-person capture.

| Spec | Value |
|------|-------|
| Dimensions | 51 x 21 x 11 mm (2 x 0.83 x 0.43 in) |
| Weight | 0.58 oz (16.6 g) without magnetic pin |
| Battery | 270 mAh; 20 hours continuous recording; 40-day standby |
| Storage | 64 GB local (~480 hours of audio) |
| Microphones | 2 dual-MEMS with beamforming |
| Connectivity | BLE 5.2, Wi-Fi |
| Recording modes | In-person only (no VPU for phone calls) |
| Colors | Cosmic Gray, Lunar Silver, Sunset Purple |
| Price | $169 USD |
| Accessories | Magnetic pin, clip, charging dock, USB-C cable |

**Dual-MEMS technology:** The two-microphone configuration achieves a 3 dB SNR improvement through coherent signal averaging -- audio signals combine with 6 dB gain while noise increases by only 3 dB. Beamforming and spatial audio capture reduce background noise. MEMS sensors are resistant to mechanical vibration and temperature shock.

### 1.4 Plaud NotePin S

Current-generation wearable with tactile record button and enhanced battery.

| Spec | Value |
|------|-------|
| Dimensions | 51 x 21 x 11 mm |
| Weight | 0.61 oz (17.4 g) without magnetic pin |
| Battery | 320 mAh; 20 hours continuous recording; 40-day standby |
| Storage | 64 GB local (~480 hours of audio) |
| Microphones | 2 dual-MEMS with AI speech enhancement |
| Connectivity | BLE, Wi-Fi |
| Recording modes | In-person only |
| Price | $179 USD |
| Accessories | Magnetic pin, clip, lanyard, wristband, charging dock, USB-C cable |

**Key differentiators from NotePin:** Tactile record button (long-press to start, short-press to highlight), larger 320 mAh battery, wristband and lanyard accessories included, Press to Highlight feature for real-time human-AI alignment.

### 1.5 Device Selection Guide

| Scenario | Recommended Device |
|----------|--------------------|
| Phone calls + in-person meetings | Note Pro (Smart Dual Mode) |
| Phone calls on a budget | Note (manual mode switch) |
| All-day wearable, hands-free capture | NotePin S (tactile button) |
| Wearable on a budget | NotePin |
| Maximum recording range (conference rooms) | Note Pro (5 m Enhance mode) |
| Maximum battery life | Note Pro (50 h Endurance mode) |
| Multimodal capture (audio + text + images) | Note Pro (InstantView display) |

---

## 2. Software Ecosystem

### 2.1 Plaud App (iOS / Android)

The primary companion app for device pairing, recording management, and AI features.

**Setup flow:**
1. Download from App Store (iOS) or Google Play (Android)
2. Create a Plaud account
3. Open the app and tap "Connect" in the top-left corner
4. Select your device model (Note, Note Pro, NotePin, NotePin S)
5. Press and hold the device record button until the white light flashes
6. Select the device from the discovered list to pair

**Recording management:**
- Auto-upload over Wi-Fi or cellular after recording sync
- Auto-generation of transcripts and summaries with AI-selected language, template, and LLM
- Manual recording: start/stop from the app or from the device button
- Import external audio files (MP3, WAV) for transcription

**AI features in-app:**
- Transcription with speaker labels and timestamps
- Summary generation (auto or custom template selection)
- Mind map generation from any recording
- Ask Plaud chat interface for Q&A against recordings
- Press to Highlight (NotePin S / Note Pro) marks priority moments for AI
- Text input and image attachment for multimodal context

### 2.2 Plaud Web (web.plaud.ai)

Browser-based access to the full Plaud ecosystem, no installation required.

- Upload audio files directly for transcription
- View, manage, and export recordings, transcripts, and summaries
- Sharing via public link (anyone with URL) or invite-only link (specific email addresses)
- Cross-platform sync: same account, subscription, and cloud files across App, Web, and Desktop
- MCP bearer token extraction point for community MCP servers

### 2.3 Plaud Desktop (macOS / Windows)

Desktop application for ambient meeting capture without joining as a bot.

- Captures system audio during Zoom, Teams, Google Meet, and other meeting tools
- No bot joins the call -- records locally from computer audio
- Auto-detect meeting start or manual recording trigger
- Press to Highlight during meetings to mark priority audio moments
- Inline text notes added as AI context for richer summaries
- Screenshot capture during meetings for visual context (slides, diagrams)
- Recordings sync to Plaud cloud for transcription and summarization

---

## 3. Plaud Intelligence (AI Engine)

Plaud Intelligence is the AI layer powering transcription, summarization, and conversational Q&A. Current version: Intelligence 3.0 (launched October 2025).

### 3.1 Transcription Engine (ASR)

| Feature | Detail |
|---------|--------|
| Languages | 112 languages supported |
| Speaker diarization | Auto speaker labeling with word-level timestamps |
| Custom vocabulary | 10+ built-in industry glossaries (medical, legal, finance, etc.) |
| Glossary matching | Prioritizes glossary-matched terms even with similar pronunciations |
| Auto-detection | Smart detection of language and number of speakers |
| Noise handling | Trained specifically for offline audio: handles background noise, overlapping speakers, distance from mics |
| Formatting | Automatic punctuation and paragraph formatting |

### 3.2 Summary Generation

| Feature | Detail |
|---------|--------|
| Templates | 10,000+ professional summary templates |
| 360-degree View | Multiple role-specific summaries from one conversation (action items for sales, insights for managers, strategic overviews for leadership) |
| Auto-generation | AI selects optimal template, language, and LLM model automatically |
| Custom generation | Manual selection of template, model, and language |
| Mind maps | Visual structured overview of conversation topics and relationships |
| Multimodal input | Audio, text notes, images, and highlights all feed into summary context |

### 3.3 AI Models

Plaud Intelligence runs on leading LLMs. The available models rotate; as of May 2026 the options include:
- GPT-5.2 (OpenAI)
- Claude Sonnet 4.5 (Anthropic)
- Gemini 3 Pro (Google)

Users can select the model during custom generation, or let auto-generation choose the optimal model. Check the Plaud App or Web for the current model list -- Plaud updates available models without notice.

### 3.4 Ask Plaud

Interactive Q&A against your recordings with reference-based answers.

- **Single-file mode:** Ask questions about a specific recording. Tap timestamps in answers to jump to that audio moment.
- **Cross-file mode (Ask Across All Files):** Search and get AI answers drawn from all transcribed files.
- **Reference-based answers:** Every response is grounded in and traceable back to original audio -- no hallucinated content.
- **Pre-built actions:** Summarize, Get Insights, Generate To-Dos, Write Email.
- **Smart suggestions:** Follow-up questions recommended by AI for deeper exploration.
- **Save as note:** Pin any AI answer as a persistent note attached to the recording.

### 3.5 Press to Highlight

Available on NotePin S and Note Pro. A short press on the device (or tap in the app) marks an audio moment as a priority cue. The AI weights highlighted moments higher when generating summaries and answering questions. Enables real-time human-AI alignment during live conversations.

---

## 4. Cloud, Export, and Sharing

### 4.1 Cloud Sync (Private Cloud Sync / PCS)

| Setting | Behavior |
|---------|----------|
| PCS **off** (default) | Recordings uploaded for transcription, then server data deleted immediately after processing is complete and results are transferred back to the phone |
| PCS **on** | Data encrypted and stored in cloud; persists until user explicitly deletes or disables PCS |

When PCS is enabled, data is accessible across Plaud App, Plaud Web, and Plaud Desktop. PCS is **required** for MCP server access.

### 4.2 Export Formats

**Audio:**
| Format | Availability |
|--------|-------------|
| MP3 | Plaud App + Plaud Web (compressed) |
| WAV | Plaud App only (lossless) |

**Transcripts and summaries:**
| Format | Content |
|--------|---------|
| TXT | Plain text, no timestamps |
| SRT | Subtitle format with timestamps (for video captioning) |
| DOCX | Word document with optional timestamps and speaker labels |
| PDF | Formatted document with optional timestamps and speaker labels |

Plaud supports 27+ export format combinations across audio, transcript, and summary content types.

### 4.3 Sharing

- **Public link:** Anyone with the URL can view
- **Invite-only link:** Restricted to specific Plaud account email addresses
- **Email delivery:** AutoFlow can email transcripts and summaries automatically

---

## 5. Subscription Plans

| Plan | Transcription | Price | Target User |
|------|--------------|-------|-------------|
| **Starter** (free) | 300 min/month | $0 | Occasional recording |
| **Pro** | 1,200 min/month (20 h) | $99.99/year or $17.99/month | Regular meeting volume |
| **Unlimited** | Up to 24 h/day | $239.99/year or $29.99/month | Doctors, lawyers, journalists, sales |
| **Team** | Up to 24 h/day per seat | Per-seat pricing | Organizations with admin management |

**Extra Minutes top-up:** 600, 3,000, or 6,000-minute packages available without changing plan.

Every activated Plaud device includes the free Starter plan (300 min/month) with full access to all Plaud Intelligence features.

---

## 6. Developer Platform

### 6.1 Status and Access

The Plaud Developer Platform launched October 2025 and is in **private beta**. Request access at dev.plaud.ai. Documentation lives at docs.plaud.ai.

### 6.2 Authentication

| Method | Use Case |
|--------|----------|
| OAuth 2.0 | User-facing SaaS integrations |
| API key (Client ID + Secret Key) | Server-to-server calls |
| Bearer token (from web.plaud.ai session) | MCP / CLI / community tools |
| Webhook signature validation | Verifying `Plaud-Signature` header on inbound webhooks |

All credentials are **region-scoped** -- every call must target the matching regional host.

### 6.3 Regional Endpoints

| Region | Host | Status |
|--------|------|--------|
| US | `platform-us.plaud.ai` | Available |
| JP | `platform-jp.plaud.ai` | Available |
| EU | `platform-eu.plaud.ai` (Frankfurt) | Available |
| SG | `platform-sg.plaud.ai` | Available |

### 6.4 API Endpoints

| Category | Purpose |
|----------|---------|
| Token | `POST /api/oauth/api-token` -- generate API token from client credentials |
| Devices | `GET /api/devices/` -- list paired devices |
| Files | List, get, rename, delete recordings |
| Transcripts | Retrieve transcription text with timestamps and speaker labels |
| Summaries | Retrieve AI-generated summaries |
| Metadata | Recording metadata (duration, date, device, language) |
| ASR | Server-side transcription submission (manual transcription path) |
| Webhooks | Receive transcript/summary completion events |

### 6.5 Embedded SDK (iOS / Android)

For developers building apps that pair directly with Plaud hardware.

**Integration flow:**
1. Create app in developer console, receive Client ID + Secret Key
2. Install the Plaud SDK (iOS 13.0+ / Android API 21+)
3. Pair device over BLE (3 SDK calls to connect)
4. Control recording sessions (start, stop, pause)
5. Transfer audio from device to app
6. Choose transcription path:
   - **Automatic:** Enable webhook -- Plaud processes audio and delivers transcript in payload
   - **Manual:** Disable auto-transcription, receive raw audio, submit to Plaud ASR endpoint or your own provider

SDK capabilities: firmware updates, device fleet management, usage monitoring. Starter template app available on GitHub (github.com/Plaud-AI/plaud-template-app).

### 6.6 Python API Client (plaud-ai)

Community-maintained Python wrapper for the Developer Platform API.

```bash
pip install plaud-ai
```

```python
from plaud_ai import PlaudAIAPIClient, PlaudAIAPIToken, PlaudAIDevicesAPI

# Generate token
api_token = PlaudAIAPIToken()
token = api_token.generate_token(client_id="...", secret_key="...")

# Initialize client
client = PlaudAIAPIClient(api_token=token.data)
devices = PlaudAIDevicesAPI(client)
response = devices.list()

# Webhook signature validation
from plaud_ai.webhooks import is_valid_signature
valid = is_valid_signature(body, signature_header, webhook_secret)
```

---

## 7. MCP Server Integration

### 7.1 Official Plaud MCP (Connector Directory)

The recommended path for Claude users. One-click OAuth setup, no token extraction needed.

**Remote MCP server URL:** `https://mcp.plaud.ai/mcp`

**Setup for Claude Web/Desktop:**
1. Open Claude settings and go to Connectors
2. Click "Add custom connector"
3. Fill in: Name = `Plaud`, Remote MCP server URL = `https://mcp.plaud.ai/mcp`
4. Click Add -- Claude opens the Plaud authorization page
5. Sign in with your Plaud account and click Authorize
6. Plaud tools are now available in all conversations

**Cross-surface:** Connecting once in Claude Web or Desktop automatically enables Plaud across Claude Web, Claude Desktop, and Claude Code. The connection is tied to your Claude account, not the device.

**Prerequisite:** Plaud Cloud Sync (PCS) must be enabled.

**Supported clients:** Claude (Web, Desktop, Code), ChatGPT Web, Cursor, Windsurf, Gemini CLI, Zed.

### 7.2 mcp-server-plaud (PyPI -- Community Python)

For users who want local control or need to configure environment variables.

```bash
pip install mcp-server-plaud
# or
uvx mcp-server-plaud
```

**Claude Desktop config** (`~/Library/Application Support/Claude/claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "plaud": {
      "command": "uvx",
      "args": ["mcp-server-plaud"],
      "env": {
        "PLAUD_TOKEN": "bearer eyJ...your_token_here"
      }
    }
  }
}
```

**Environment file method:**
```bash
mkdir -p ~/.config/plaud
echo 'PLAUD_TOKEN=bearer eyJ...your_token_here' > ~/.config/plaud/.env
```

**Important:** Include the `bearer ` prefix (with the space) when setting `PLAUD_TOKEN`.

**EU region:** Set `PLAUD_API_DOMAIN=https://web-eu.plaud.ai` in env.

**Token extraction from web.plaud.ai:**
1. Log in to web.plaud.ai
2. Open browser DevTools (Network tab)
3. Find any API request and copy the `Authorization` header value
4. The token is a long JWT string starting with `eyJ...`

**Available tools:**

| Tool | Purpose |
|------|---------|
| `plaud_verify_token` | Check if auth token is valid |
| `plaud_list_recordings` | List recordings with pagination |
| `plaud_get_transcript` | Full transcript + AI summary for a recording |
| `plaud_search_recordings` | Search recordings by title keyword |
| `plaud_get_recent` | Recent recordings with full transcripts |

### 7.3 plaud-toolkit (npm -- Community TypeScript)

Monorepo with three packages: `@plaud/core`, `@plaud/cli`, `@plaud/mcp`. Status: Alpha, MIT license.

**Authentication:** Email + password login; tokens obtained automatically and last ~300 days with silent refresh within 30 days of expiry.

**MCP tools:**

| Tool | Purpose |
|------|---------|
| `plaud_list_recordings` | List all recordings |
| `plaud_get_transcript` | Retrieve transcription for a recording |
| `plaud_get_recording_detail` | Full recording metadata |
| `plaud_user_info` | Account profile |
| `plaud_get_mp3_url` | Audio download URL |

Credentials stored in `~/.plaud/config.json` (mode 0600).

### 7.4 Plaud MCP Plugin (Extended Toolset)

When using the Plaud MCP plugin or Connector Directory, additional tools may be available:

| Tool | Purpose |
|------|---------|
| `find_recordings` | Browse/filter recordings |
| `get_recording` | Full recording metadata |
| `get_audio_url` | Presigned download URL (5-min or 24-hr expiry) |
| `transcribe` | ElevenLabs transcript (cached or live) |
| `get_content` | AI summary or meeting notes |
| `trigger_processing` | Plaud server-side AI processing |
| `get_account_info` | Account profile and membership |
| `get_processing_status` | AI processing queue status |
| `list_languages` | Supported language codes |
| `memory_search` | Semantic search across indexed transcripts |
| `memory_ingest` | Bulk index recordings into memory store |

**File ID resolution** (all tools accepting `file_ref`):
- Full ID: `4f757af256ecba4fab502739c122dc78`
- Short prefix: `4f75`
- Row number: `3` (from last `find_recordings` output; persists within session)

**Cost safety:** Tools that consume external credits (`transcribe`, `memory_ingest`) prompt for confirmation via MCP elicitation before proceeding.

For detailed tool usage, session state, and MCP prompts, see the **using-plaud-mcp** sibling skill.

### 7.5 CLI Tool

The `plaud` CLI (via `@plaud/cli` or community packages) provides terminal-based recording management.

| Command | Purpose |
|---------|---------|
| `plaud login` | Authenticate and store credentials |
| `plaud files` / `plaud list` | List recordings (ID, title, date, duration) |
| `plaud transcript <id>` | Retrieve transcription |
| `plaud audio <id>` | Get temporary download URL (expires 24h) |
| `plaud search <query>` | Search recordings by keyword |
| `plaud download <id> <path>` | Save audio file locally |
| `plaud sync <folder>` | Sync all recordings to a local folder |
| `plaud export` | Export transcripts to local files |

---

## 8. Workflow Automation

### 8.1 AutoFlow (Native)

Built-in automation that monitors incoming recordings and auto-processes them.

**How it works:**
1. Set trigger conditions (e.g., keyword spoken in first 60 seconds of recording)
2. Configure outputs: transcription, summary, or both
3. Set email recipient for automatic delivery
4. When a synced/uploaded recording matches the trigger, Plaud processes and delivers

**Use cases:** Client meeting transcripts to email, standup summaries to inbox, automatic processing of all recordings without manual trigger.

### 8.2 Zapier Integration

**Triggers:** Transcripts generated, summaries generated (fires on transcription, re-transcription, summarization, or re-summarization completion).

**Popular zaps:**
- Archive transcript + summary to Google Drive / OneDrive
- Create Notion database entry from summary
- Post recap to Slack / Teams channel
- Create monday.com item or Jira issue from key takeaways
- Surface action items in Todoist, Asana, or other task managers
- Update CRM records (Salesforce, HubSpot)
- AI step (Claude/ChatGPT) to extract action items, owners, and due dates from summary

Plaud integrates with 9,000+ apps on Zapier.

### 8.3 MCP-Powered Workflows

With the Plaud MCP connected to Claude:
- Summarize a specific meeting by date or title
- Draft follow-up emails from meeting content
- Extract action items across one or more recordings
- Compare discussions across meetings (planning vs retro)
- Generate daily task summaries with preliminary research
- Execute agentic workflows (extract action items and sync to Asana or HubSpot)

### 8.4 Third-Party CRM Integrations

**ScribbleIQ** extracts structured CRM data from Plaud transcripts:
- Contacts, companies, opportunities, action items
- Auto-create or update records in Salesforce, HubSpot, or Zoho

---

## 9. Privacy, Security, and Compliance

### 9.1 Compliance Certifications

| Standard | Status |
|----------|--------|
| GDPR | Compliant |
| ISO 27001 | Certified |
| ISO 27701 | Certified |
| SOC 2 Type II | Certified |
| HIPAA | Compliant |
| EN 18031 | Certified |

### 9.2 Encryption

| Layer | Standard |
|-------|----------|
| In transit | HTTPS / TLS 1.2 or above |
| At rest | AES-256 encryption |
| Application layer | Secondary encryption with unique per-user keys |

### 9.3 Data Storage Model

See Section 4.1 for the PCS on/off behavior table. Data is stored and processed in the nearest regional data center:
- US (AWS-West)
- Europe (Frankfurt)
- Singapore
- Japan

### 9.4 AI Training Policy

- Plaud does **not** train on user data by default
- Explicit opt-in required before any data is used for model improvement
- All LLM subprocessors are contractually bound to zero data retention via DPAs

### 9.5 Data Access

Plaud personnel do not routinely access user data. Limited access occurs only for: technical support, error diagnosis, security investigations, or legal compliance.

### 9.6 Recording Consent

- Users must obtain consent from all participants before recording where required by law
- Users warrant that uploaded audio does not infringe third-party rights
- Two-party consent jurisdictions (California, Illinois, many EU countries) require explicit agreement from all recorded parties

### 9.7 Token Security

- Bearer tokens for MCP/CLI expire periodically (community tokens from web.plaud.ai)
- Official MCP via Connector Directory uses OAuth -- no manual token management
- Store tokens in secure locations (`~/.plaud/config.json` with mode 0600, or `~/.config/plaud/.env`)
- Re-extract from web.plaud.ai if you receive 401 errors
- Never commit tokens to version control

---

## 10. Open-Source and Community Tools

### 10.1 OpenPlaud

Self-hosted AI transcription interface for Plaud devices. Independent open-source project (AGPL-3.0), not affiliated with Plaud Inc.

- Connects to your Plaud account and transcribes with your own AI keys
- Supports OpenAI, Anthropic, Groq, OpenRouter, Together, Fireworks, LM Studio, Ollama, vLLM
- Can run Whisper or local Llama models entirely on your own machine
- AES-256-GCM encryption for API keys, bearer tokens, and user content
- Self-hosted (free) or hosted option available
- Website: openplaud.com / GitHub: openplaud/openplaud

### 10.2 Applaud

Free, open-source audio transcription and summarization tool. Provides transcription, flashcards, questions, and summaries. Supports local LLM models via Ollama.

### 10.3 python-plaud-ai

Community Python API wrapper for the Plaud Developer Platform. MIT license.
- GitHub: DmytroLitvinov/python-plaud-ai
- PyPI: `pip install plaud-ai`
- Covers: token generation, device listing, webhook signature validation

---

## 11. Integration Patterns

### 11.1 Browser Extension Integration

For extensions that ingest Plaud data into a local corpus:

1. **MCP bridge:** Use a native-messaging host to call Plaud MCP tools from the extension service worker
2. **Transcript ingestion:** Call `find_recordings` + `get_content` to pull transcripts and summaries
3. **Local storage:** Store in IndexedDB or a local database; optionally mirror to a backend
4. **Incremental sync:** Track last-synced recording ID or timestamp to avoid re-fetching

### 11.2 Server-Side Pipeline

For backends that maintain a database mirror of Plaud recordings:

1. **Webhook receiver:** Register a Plaud webhook endpoint to receive transcript/summary completion events
2. **Polling fallback:** Periodically call `plaud_list_recordings` or `plaud_get_recent` to catch missed webhooks
3. **Normalization:** Map Plaud transcript JSON (speaker labels, timestamps) to your document schema
4. **Deduplication:** Use Plaud file ID as the dedup key; check before inserting

### 11.3 Plaud-to-LLM Pipeline

For feeding Plaud transcripts into custom LLM workflows:

1. Retrieve transcript via API or MCP (`get_content` or `plaud_get_transcript`)
2. Chunk by speaker turns or time windows for long recordings
3. Wrap untrusted transcript text in XML tags (e.g., `<untrusted_meeting_transcript>`) with HTML-entity escape to prevent prompt injection
4. Submit to LLM with task-specific system prompt (action items, email draft, CRM update)

### 11.4 Zapier / Webhook Pipeline

See Section 8.2 for Zapier trigger details and popular zaps. For custom webhook pipelines:

1. Register a webhook endpoint with Plaud (Developer Platform, private beta)
2. Validate inbound payloads using the `Plaud-Signature` header and your webhook secret
3. Parse transcript JSON (speaker labels, timestamps) and route to downstream systems
4. Optionally add an AI step to extract structured data before routing to CRM or task management

---

## 12. Anti-Patterns and Limitations

| Anti-Pattern | Why It Fails | Correct Approach |
|-------------|-------------|-----------------|
| Recording phone calls with NotePin | NotePin lacks VPU; only captures speaker audio, not both sides | Use Note or Note Pro (attach to phone back for VPU capture) |
| Skipping PCS and expecting MCP to work | MCP requires cloud-synced data | Enable Private Cloud Sync before using MCP |
| Hardcoding bearer tokens | Tokens expire periodically | Use OAuth via Connector Directory, or implement token refresh |
| Ignoring region scoping | API calls fail with wrong-region errors | Match API host to account region |
| Treating Ask Plaud answers as source-of-truth without checking references | AI answers are reference-based but still AI-generated | Always verify by clicking timestamp references back to original audio |
| Running AutoFlow without keyword triggers | Processes every recording, consuming transcription minutes | Set keyword triggers to filter relevant recordings |
| Expecting real-time transcription | Plaud processes after recording sync, not during | Transcripts available after upload + AI processing completes |
| Mixing cloud sync regions | Data may not appear across platforms | Keep all surfaces on the same account and region |

---

## 13. Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| MCP returns 401 | Bearer token expired | Re-extract token from web.plaud.ai, or re-authorize via Connector Directory |
| No recordings visible in MCP | PCS not enabled | Enable Plaud Cloud Sync in mobile app settings |
| Transcript empty | Processing still in progress | Check `get_processing_status` or wait for webhook callback |
| Wrong region errors | Credential/host mismatch | Ensure API calls target the region matching your credentials |
| Audio URL expired | Presigned URL past TTL (5 min or 24 h depending on method) | Re-request via `get_audio_url` or `plaud audio <id>` |
| CLI "command not found" | Package not installed | Install via `npm i -g @plaud/cli` or use `npx @plaud/mcp` |
| Device not pairing | Bluetooth not active or device not in pairing mode | Hold record button until white light flashes; ensure Bluetooth is on |
| NotePin call recording is muffled | NotePin captures only ambient audio, not VPU vibrations | Put call on speaker and place NotePin close, or use Note/Note Pro instead |
| AutoFlow not triggering | Keyword not spoken in first 60 seconds | Ensure the trigger keyword is said within the first minute of recording |
| Plaud Desktop not detecting meetings | Meeting app audio routing misconfigured | Check system audio settings; ensure Plaud Desktop has microphone/audio permissions |
| Transcription minutes exhausted | Hit plan limit | Purchase Extra Minutes top-up or upgrade plan |
| Webhook payloads not arriving | Webhook URL unreachable or signature validation failing | Verify endpoint is publicly accessible; validate `Plaud-Signature` header with webhook secret |

---

## 14. Cross-References

### Sibling Skills
- **using-plaud-mcp** -- Detailed MCP tool usage, session state, file-ID resolution, cost safety, prompts, and elicitation patterns
- **chrome-extension-expert** -- Chrome extension architecture and native-messaging bridge patterns for ingesting Plaud data
- **ai-agent-engineering** -- Building agentic workflows that consume Plaud transcripts
- **security-compliance-auditor** -- Auditing compliance posture (GDPR, HIPAA, SOC 2) for Plaud data flows
- **granola-transcription** -- Alternative meeting transcription patterns for comparison

### External References
1. Plaud Official Site: https://www.plaud.ai/
2. Plaud Device Comparison: https://www.plaud.ai/pages/plaud-device-comparison
3. Plaud Note Pro Product Page: https://www.plaud.ai/products/plaud-note-pro
4. Plaud NotePin S Product Page: https://www.plaud.ai/products/plaud-notepin-s
5. Plaud Intelligence (AI Engine): https://www.plaud.ai/pages/plaud-intelligence
6. Plaud Trust Center: https://www.plaud.ai/pages/trust
7. Plaud Developer Platform: https://dev.plaud.ai/
8. Plaud Developer Docs: https://docs.plaud.ai/
9. Plaud MCP Documentation: https://docs.plaud.ai/documentation/plaud_app/mcp
10. Plaud Support Center: https://support.plaud.ai/
11. Plaud MCP Support Article: https://support.plaud.ai/hc/en-us/articles/57751078986265-Plaud-MCP
12. Plaud Cloud Sync Support: https://support.plaud.ai/hc/en-us/articles/51820671018265-Private-Cloud-Sync
13. Plaud AI Plan Pricing: https://www.plaud.ai/pages/plaud-ai-plan-pricing
14. Plaud Intelligence 3.0 Launch: https://www.plaud.ai/blogs/news/plaud-intelligence-3-0-launch
15. Plaud MCP & CLI Announcement: https://www.plaud.ai/blogs/news/introducing-plaud-mcp-and-cli
16. Plaud Dual-MEMS Technology: https://www.plaud.ai/blogs/articles/plaud-notepin-dual-mems
17. Plaud Desktop Download: https://support.plaud.ai/hc/en-us/articles/53792807283225-Download-the-Plaud-Desktop
18. Plaud Zapier Integration: https://zapier.com/apps/plaud/integrations
19. Community MCP Server (PyPI): https://pypi.org/project/mcp-server-plaud/
20. Community Toolkit (TypeScript): https://github.com/sergivalverde/plaud-toolkit
21. Python API Client: https://github.com/DmytroLitvinov/python-plaud-ai
22. OpenPlaud (Self-Hosted): https://github.com/openplaud/openplaud
23. Plaud Template App (SDK): https://github.com/Plaud-AI/plaud-template-app
24. Plaud Data Handling FAQ: https://support.plaud.ai/hc/en-us/articles/49922008621337
25. Plaud AI Data Usage Policy: https://global.plaud.ai/pages/ai-data-usage-transparency-policy

## Related Skills

- **granola-transcription** — Granola AI meeting notes (complementary software-only meeting capture tool, bot-free)
- **using-plaud-mcp** — Quick-start for the Plaud MCP server (subset of this skill's MCP section)
- **chrome-extension-expert** — Chrome extension and native-messaging host patterns used by Plaud's native bridge
- **ai-agent-engineering** — Agentic pipelines for downstream meeting action-item extraction and routing
