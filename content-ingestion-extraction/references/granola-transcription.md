<!-- hub-reference-banner -->
> **Reference file — part of the `content-ingestion-extraction` hub.** Formerly the standalone `granola-transcription` skill.
> Sibling topics in this family are now reference files under the hubs (`content-ingestion-extraction`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: granola-transcription
version: 1.1.1
updated: 2026-05-31
description: >
  Granola AI meeting notes — bot-free audio capture, ProseMirror note model,
  AI summarization, REST API (get-documents, get-document-transcript,
  get-document-panels), WorkOS/API-key auth, MCP integration, native host bridge
  for Chrome extensions, Slack/Notion/HubSpot/Zapier integrations, corpus pipeline
  wiring, and polling sync patterns.
  TRIGGER: fetching Granola transcripts or documents via API or native host;
  building polling sync loops against api.granola.ai; parsing ProseMirror content;
  handling Granola 401/403/429 errors; wiring Granola data into a Chrome extension
  corpus pipeline; questions about Granola product features, pricing, or integrations.
  SKIP: real-time meeting capture during a live call (handled by the Granola desktop
  app itself); pre-recorded audio file transcription (Granola does not accept file
  uploads); general meeting transcription not involving Granola.
category: developer
tags: [granola, meeting-notes, transcription, api, mcp, chrome-extension, prosemirror, corpus]
related_skills: [plaud-integration, using-plaud-mcp, chrome-extension-expert, integration-clients]
whenToUse:
  - "How do I fetch Granola transcripts via the API?"
  - "Granola native host token path"
  - "Parse ProseMirror notes from Granola"
  - "Granola WorkOS token refresh"
  - "Build incremental sync loop for Granola documents"
  - "Granola 401 or 403 error"
  - "Set up Granola MCP in Claude Code"
  - "Granola pricing plans comparison"
  - "Wiring Granola into a Chrome extension corpus pipeline"
  - "Granola anti-patterns and common mistakes"
whenNotToUse:
  - "Real-time meeting capture: handled by the Granola desktop app, not the API"
  - "Pre-recorded audio transcription: Granola only supports live meeting capture"
  - "Interactive AI assistant context: use the official MCP integration instead of REST API"
---

# Granola: Comprehensive Integration & Product Reference

## When to Use This Skill

- Fetching transcripts, documents, or summaries via native host, manual token, MCP bridge, or API key
- Building polling sync loops against `api.granola.ai`
- Parsing ProseMirror note content or transcript utterances from Granola
- Handling 401/403/429 errors or transport failures against Granola endpoints
- Wiring Granola meeting data into a Chrome extension corpus pipeline
- Understanding Granola product features, pricing, or competitive positioning

## When NOT to Use This Skill

- **Free plan API access:** Granola's REST API requires a Business or Enterprise plan. Free plan requests return `403 Forbidden` — surface this to the user, do not retry.
- **Interactive AI assistant context:** Use the official MCP integration (`https://mcp.granola.ai/mcp`) instead of the REST API for Claude/ChatGPT/Cursor sessions.
- **Live meeting capture:** This skill covers post-meeting data retrieval and product knowledge. Real-time audio capture is handled by the Granola desktop app itself.
- **Pre-recorded audio processing:** Granola only supports live meeting transcription; it does not accept uploaded audio files.

---

## 1. Product Overview

### What Granola Is

Granola is an AI-powered meeting notepad that captures audio locally on your device (no bot joins the call), transcribes in real time, and produces structured AI-enhanced notes after the meeting ends. Founded in March 2023 by Chris Pedregal and Sam Stephenson (London, UK), the company reached unicorn status in March 2026 with a $125M Series C at a $1.5B valuation (led by Index Ventures and Kleiner Perkins). Previous rounds: $43M Series B (2025, led by Nat Friedman and Daniel Gross), bringing total funding to $192M.

**Key differentiator:** Granola uses a hybrid human+AI approach. During a meeting, you take notes in Granola's notepad. When the call ends, Granola combines your notes with the full transcript and uses LLMs to produce structured, detailed meeting notes. The AI prioritizes topics you flagged as important in your own notes, then fills in supporting details from the transcript. No other major competitor does this hybrid enhancement.

**Bot-free recording:** Granola captures system audio and microphone audio directly on your computer. No meeting bot appears in the participant list, no recording announcement interrupts the conversation. Participants do not see an attendee. Audio is streamed to the transcription provider in real time and is never stored — only the text transcript and your notes are saved.

### Supported Platforms

| Platform | Status | Notes |
|----------|--------|-------|
| macOS | Full support (primary) | macOS 13+, ideally 14.2+; native Electron app |
| Windows | Full support (June 2025) | Feature-equivalent to macOS |
| iOS | Supported (April 2025) | In-person meeting transcription via iPhone; phone call capture |
| Android | Not supported | No timeline announced |
| Web | View/edit only | Cannot capture or transcribe — viewing and editing existing notes only |

### Pricing (as of 2026)

| Plan | Price | Key Features |
|------|-------|-------------|
| **Basic** (Free) | $0 | AI meeting notes, 25-note history cap (total, not per month), AI chat, custom templates, multi-language, model training opt-out |
| **Business** | $14/user/month | Unlimited history, advanced AI models, integrations (Slack, Notion, HubSpot, Attio, Affinity, Zapier), MCP (full history), API access, team folders, centralized billing |
| **Enterprise** | $35/user/month | SSO (50+ users), enterprise API, org-wide auto-deletion periods, admin controls for sharing/API, team-wide model training opt-out, priority support, usage analytics |

**Plan change history:** The old "Pro" ($18/mo) and "Individual" plans were replaced in 2026 by the current Basic/Business/Enterprise structure.

### Competitive Positioning

| Feature | Granola | Otter | Fireflies | Fathom |
|---------|---------|-------|-----------|--------|
| Bot-free capture | Yes (core differentiator) | No (bot joins) | No (bot joins) | No (bot joins) |
| Hybrid human+AI notes | Yes (unique) | No | No | No |
| Free tier | 25 notes total | 300 min/month | 800 min total | Unlimited recordings, 5 AI summaries/month |
| Live collaborative transcription | No | Yes (unique) | No | No |
| Deep CRM automation | Via Zapier | No | Yes (field-level) | Limited |
| Conversation intelligence | Recipes/Chat | Basic | Advanced (sales) | Basic |
| Video/audio recording | No (text only) | Yes | Yes | Yes |

**Granola does NOT record or store audio or video.** It produces text-only output (transcripts + AI notes). This is a deliberate design choice for privacy but means there is no playback capability.

### Notable Customers

Vanta, Gusto, Thumbtack, Asana, Cursor, Lovable, Decagon, Mistral AI.

---

## 2. AI Engine

### Transcription Pipeline

Granola captures system audio + microphone audio on the device and streams it to cloud transcription providers. The app uses best-in-class providers — AssemblyAI for iOS beta builds and other providers (Deepgram reported by multiple sources) for production macOS/Windows. The specific provider may change; the app abstracts this.

**Audio flow:** Device microphone + system audio -> streamed to transcription provider -> real-time transcript returned -> displayed in-app. Audio is never stored on disk or in the cloud. Only the text transcript persists.

**Speaker diarization:**
- Desktop (macOS/Windows): Limited to "Me" (microphone) and "Them" (system audio). Live diarization of individual remote speakers is not yet supported.
- iOS: Supports face-to-face diarization but selects a single language per meeting.
- The `audio_source_type` field in transcripts reliably separates local mic (`"microphone"`) from remote participants (`"speaker"` / system audio), but named speakers are unreliable with 3+ participants.

**Language support (10+ languages):**
Desktop multi-language: English, French, German, Spanish, Italian, Portuguese, Dutch, Japanese, Russian, Hindi. Can switch languages mid-call with multi-language mode enabled.
iOS additional: Mandarin Chinese, Finnish, Korean, Polish, Turkish, Ukrainian, Vietnamese. iOS picks a single language per meeting.

**Internal jargon:** Users can add company-specific words to "Internal Jargon" in settings to improve recognition accuracy. Not available with multi-language mode.

### Summarization Engine

After a meeting ends, Granola combines the user's typed notes (used as priority anchors) with the full transcript and runs them through LLMs. The AI engine uses:
- **OpenAI** models (GPT-4 class) for summarization
- **Anthropic** models for chat and alternative processing
- **Google** models as additional options
- Users on Business/Enterprise can select their preferred model or use "Auto" selection

The output follows the user's chosen template structure, extracting action items with owners/deadlines, decisions, key discussion points, and any custom sections defined in the template.

**Multi-model support (Granola 2.0+):** Users can switch between reasoning models from OpenAI, Anthropic, and Google via a dropdown, or use Auto selection for the system's best pick.

### AI Chat (Granola Chat)

Granola Chat allows natural-language queries across meetings:
- **Single meeting:** Ask questions about one specific meeting
- **Folder/Space:** Query across all meetings in a shared folder or People/Companies record
- **Global:** Search your entire meeting history from the home screen
- Inline citations with jump-to-source links back to the specific transcript moment
- Can search notes, connect ideas across meetings, research online, and help find answers

### Recipes (Saved Prompts)

Recipes are saved chat prompts that process meeting notes through a specific lens:
- Access via `/` in Granola Chat
- 29+ built-in recipes created by productivity experts (Lenny Rachitsky, Matt Mochary, Ridd)
- Examples: "Coach me", "Prep me", "Write a brief", "Improve our sales process", "Extract Feature Requests", "Generate Follow-up Email", "Identify Recurring Pain Points", "Summarize Deal Risks", "Create Executive Summary"
- Users can create custom Recipes, share with colleagues, or keep private
- Recipes draw from your meeting transcript and notes for contextual output

---

## 3. Note Structure & Document Model

### ProseMirror Format

Granola stores notes in ProseMirror format — the same rich-text framework used by Notion, The New York Times, and Atlassian. Content is a JSON tree of typed nodes with optional marks (formatting).

```js
// Top-level document structure
{
  type: "doc",
  content: [
    {
      type: "heading",
      attrs: { level: 2, id: "abc123" },
      content: [{ type: "text", text: "Action Items" }]
    },
    {
      type: "bulletList",
      content: [
        {
          type: "listItem",
          content: [
            {
              type: "paragraph",
              content: [{ type: "text", text: "Send proposal by Friday" }]
            }
          ]
        }
      ]
    }
  ]
}
```

**Supported node types:** `doc`, `heading` (with `level` and `id` attrs), `paragraph`, `bulletList`, `listItem`, `codeBlock`, `hardBreak`, `text`.

### ProseMirror to Markdown Conversion

```js
function prosemirrorToMarkdown(node, depth = 0) {
  if (!node) return '';
  if (node.type === 'text') return node.text || '';
  const children = (node.content || []).map(c => prosemirrorToMarkdown(c, depth)).join('');
  switch (node.type) {
    case 'doc':        return children;
    case 'heading':    return '#'.repeat(node.attrs?.level || 1) + ' ' + children.trim() + '\n\n';
    case 'paragraph':  return children.trim() + '\n\n';
    case 'bulletList': return (node.content || []).map(i => prosemirrorToMarkdown(i, depth)).join('');
    case 'listItem':   return '  '.repeat(depth) + '- ' + (node.content || []).map(c => prosemirrorToMarkdown(c, depth + 1)).join('').trim() + '\n';
    case 'codeBlock':  return '```\n' + children + '\n```\n\n';
    case 'hardBreak':  return '\n';
    default:           return children;
  }
}
```

### Document Data Model

```js
{
  id: "uuid",
  title: "Weekly Team Sync",
  created_at: "2026-01-15T10:00:00.000Z",
  updated_at: "2026-01-15T11:30:00.000Z",
  last_edited_at: "...",              // fallback timestamp field
  workspace_id: "uuid",
  workspace_name: "...",

  // Notes content — check in priority order:
  notes_markdown: "## Action items\n- ...",  // preferred (plain markdown)
  notes: { type: "doc", content: [...] },    // ProseMirror JSON
  notes_plain: "plain text fallback",

  // AI-generated content
  summary: "...",           // flat text or nested object — flatten defensively
  overview: "...",
  last_viewed_panel: {      // most recent AI-generated panel
    content: { type: "doc", content: [...] },  // ProseMirror JSON
    title: "Meeting Summary"
  },

  // People
  people: [
    { name: "Alice", email: "alice@example.com", display_name: "Alice K." }
  ],

  // Calendar metadata
  google_calendar_event: { summary: "...", attendees: [...] },
  meeting_date: "...",
  sources: ["microphone", "system"],  // audio source types used

  // Organization
  folder: "...", folders: [...], folder_name: "...", breadcrumbs: [...]
}
```

### Transcript Data Model

```js
// Segment/utterance shape (field names vary across API versions)
{
  speaker_name: "Alice",              // or "speaker" — unreliable with 3+ people
  audio_source_type: "microphone",    // "microphone" = local user; "speaker" = system audio (remote)
  transcribed_text: "Let's start.",   // or "text"
  source: "microphone",              // alternate field name
  start: "2026-01-15T10:00:05.000Z", // ISO 8601 or ms offset; also "start_timestamp"
  end:   "2026-01-15T10:00:08.000Z", // also "end_timestamp"
  confidence: 0.97
}
```

**Normalize transcript response defensively:**
```js
function extractSegments(data) {
  if (Array.isArray(data)) return data;
  return data?.transcript ?? data?.segments ?? data?.utterances ?? [];
}
```

### Summary Panels

Panels are structured AI-generated note sections. Each panel:
```js
{
  title: "Action Items",           // or "Meeting Summary", "Decisions", etc.
  content_markdown: "- ...",       // preferred — plain markdown
  notes_markdown: "...",           // fallback
  content: { /* ProseMirror JSON */ }  // fallback
}
```
Extract: `panelsData?.panels ?? panelsData?.document_panels ?? []`

### Hybrid Note Display

Granola's UI distinguishes user input from AI enhancement:
- **Black text:** User's own typed notes
- **Gray text:** AI-generated enhancement from transcript
- Any user edit turns text black to indicate human input

---

## 4. Templates

### Structure

Templates define the AI output structure applied after meetings end. Three foundational elements:

1. **Purpose & Context:** Meeting objectives and background ("This is a customer discovery call...")
2. **Length & Style:** Detail levels, quote inclusion, emphasis ("Focus more on what they said, not what I said")
3. **Structure & Sections:** Expected groupings with prompts ("Key discussion points", "Decisions made", "Action items with owners and deadlines")

### Built-in Templates (29+)

Organized by meeting type: Project Kick-Offs, Pipeline Reviews, 1:1s, Customer Interviews, User Research, Sales Calls, IC Memos, PRD Meetings, Standups, Retros, and more. Template display is filtered by user type selected during onboarding.

### Custom Templates

- Create via "All Templates..." menu or Settings > "Manage templates"
- Templates are private by default; shareable with colleagues who have matching company email domains
- Free email providers (Gmail, Outlook.com) cannot use team sharing
- Applying a template regenerates notes using that structure
- Available on desktop only (macOS/Windows); iOS notes can receive templates after syncing to desktop
- Test iteratively against actual meetings to refine

---

## 5. Team Features & Collaboration

### Workspaces

Granola organizes work into workspaces — containers for notes, Spaces, folders, and meeting insights.

**Personal Workspace:** Private notes, individual organization, full admin control, Granola Chat search.

**Team Workspace:** Collaborative environment, shared notes/folders, collective templates and recipes, shared knowledge base. Notes remain private unless explicitly added to a shared folder or shared directly.

**Member Roles:**
- **Admins:** Manage workspace settings, billing, team structure, data export permissions, discoverability
- **Members:** Use workspace and view billing; cannot modify subscriptions or workspace settings

**Billing:** Per-workspace (not company-wide). Free and paid subscriptions can coexist across workspaces within one account. Per-seat billing from join date.

### Spaces (Team Meeting Intelligence)

Spaces transform Granola into a collaborative knowledge base for teams:
- **Cross-conversation search:** Ask natural-language questions, get answers synthesized from all meetings in a Space
- **Jump-to-source citations:** AI answers include links to the specific transcript moment
- **Permission controls:** Manage access to sensitive conversations
- **Folder organization:** Sales Calls, Customer Feedback, Hiring Loops, Weekly Syncs, etc.

**Use cases:** "Why are we losing this deal?" across all related sales calls. "What are users consistently asking for?" across hundreds of research interviews.

### Shared Notes & Folders

- Notes are never shared by default
- Share by: adding to a folder in team space, or sharing directly with someone
- Team folders enable: Slack auto-posting, Notion sync, Zapier triggers, cross-note AI queries
- **Browse View** (Business/Enterprise): Access public folders across your email domain
- Notes can be moved between workspaces via three-dot menu or Settings

### Briefs (Pre-Meeting Summaries)

Auto-generated pre-meeting context showing attendees, previous discussion topics, and current priorities. Available on Mac and Windows.

---

## 6. Calendar & Recording

### Calendar Integration

- **Google Calendar:** Native integration; primary calendar sync
- **Microsoft Outlook:** Native sign-in support added January 2026 (previously required syncing Outlook to Google Calendar first)
- Auto-detects upcoming meetings from calendar events containing video conferencing URLs
- Notifications for upcoming meetings; auto-opens Granola for scheduled calls
- Toggle secondary/shared calendars in "Coming Up" settings
- Force refresh calendar sync via sign-out/sign-in

### Audio Capture

- **System audio:** Captures all audio output from your computer (remote participants in Zoom, Meet, Teams, WebEx, Slack Huddles, etc.)
- **Microphone:** Captures your voice
- **No bot:** No participant joins the call; no recording announcement
- **No recording stored:** Audio streams to transcription provider and is discarded; only text transcript persists
- **Supported meeting platforms:** Zoom, Google Meet, Microsoft Teams, Cisco WebEx, Slack Huddles, FaceTime, WhatsApp, and any browser-based call
- **Manual mode:** Start recording manually for ad-hoc conversations without calendar events
- **iOS:** Place iPhone on conference table for in-person meeting transcription; phone call capture supported
- **Live-only:** Does not support importing or uploading pre-recorded audio files (MP3, WAV, etc.)

### Audio Permissions

**macOS:** System Settings > Privacy & Security > Microphone (toggle on for Granola) + Screen & System Audio Recording (enable for Granola). macOS 13+ required, 14.2+ recommended.

**Windows:** Settings > Privacy > Microphone > "Let desktop apps access your microphone" must be on. Disable audio enhancements; consider enabling mono audio mode for devices that conflict with stereo processing.

---

## 7. API & Developer Patterns

### Authentication Methods

Granola supports three authentication approaches:

**1. API Keys (Business/Enterprise plans)**
```
Authorization: Bearer grn_YOUR_API_KEY
```
- Format: `grn_` prefix followed by key string
- Personal API key: access your own notes + notes shared with you (Business+)
- Enterprise API key: workspace-level access to all Team Space notes (Enterprise only)
- Enterprise admin must enable API access: Settings > Workspace > General > API access for members
- Use for: server-side polling, Chrome extension native host bridges, automation scripts

**2. WorkOS Token Exchange (Desktop app internal auth)**
```
POST https://api.workos.com/user_management/authenticate
Body: { "client_id": "...", "grant_type": "refresh_token", "refresh_token": "..." }
Response: { "access_token": "jwt...", "refresh_token": "new_rotated_token", "expires_in": 3600 }
```
- Refresh tokens are **single-use** — each exchange invalidates the old token and issues a new one
- Access token TTL: ~1 hour
- The Granola desktop app handles token refresh automatically and writes to `supabase.json`
- **Never call WorkOS directly** from integrations — re-read `supabase.json` on 401

**3. MCP OAuth (Browser-based)**
- OAuth 2.0 with Dynamic Client Registration at `https://mcp.granola.ai/mcp`
- No API keys needed; credentials handled automatically via browser sign-in flow
- Use for: interactive AI assistant contexts (Claude, ChatGPT, Cursor)

### Token File Location

**macOS:** `~/Library/Application Support/Granola/supabase.json`
**Windows:** `%APPDATA%\Granola\supabase.json`

```json
{
  "workos_tokens": {
    "access_token": "eyJ...",
    "refresh_token": "wos_rt_..."
  }
}
```
**Caveat:** `workos_tokens` may be a nested JSON string — parse with `JSON.parse()` if `typeof tokens === 'string'`.

### REST Endpoints

All endpoints are `POST` with `Authorization: Bearer <token>` and `Content-Type: application/json`. Responses may be gzip-encoded — detect and decompress:

```js
async function fetchGranolaAPI(endpoint, body, token) {
  const res = await fetch(`https://api.granola.ai${endpoint}`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
      'User-Agent': 'Granola/5.354.0',
    },
    body: JSON.stringify(body),
  });
  if (!res.ok) return { error: res.status, body: null };
  const buf = await res.arrayBuffer();
  const bytes = new Uint8Array(buf);
  let text;
  // Detect gzip magic bytes 0x1f 0x8b
  if (bytes[0] === 0x1f && bytes[1] === 0x8b) {
    // Browser: use DecompressionStream API
    const ds = new DecompressionStream('gzip');
    const writer = ds.writable.getWriter();
    writer.write(bytes); writer.close();
    text = await new Response(ds.readable).text();
    // Node.js alternative: const zlib = require('zlib');
    // text = zlib.gunzipSync(Buffer.from(buf)).toString('utf-8');
  } else {
    text = new TextDecoder().decode(bytes);
  }
  return { error: null, body: JSON.parse(text) };
}
```

#### List Documents
```
POST https://api.granola.ai/v2/get-documents
Body: { "limit": 100, "offset": 0, "include_last_viewed_panel": true }
Headers: Authorization, Content-Type, User-Agent: Granola/5.354.0, X-Client-Version: 5.354.0
```
**Response shape** (normalize defensively — field names vary):
```json
{
  "docs": [...],          // or "documents", "items", "data.docs"
  "total_count": 342,     // or "total", "count", "pagination.total"
  "limit": 100,
  "offset": 0
}
```
**Important:** Returns only documents owned by the user. For shared documents and folder contents, use `get-documents-batch`.

#### Get Documents Batch
```
POST https://api.granola.ai/v1/get-documents-batch
Body: { "document_ids": ["uuid1", "uuid2"], "include_last_viewed_panel": true }
```
Recommended for fetching shared documents and folder contents by ID.

#### Get Transcript
```
POST https://api.granola.ai/v1/get-document-transcript
Body: { "document_id": "<uuid>" }
```
Returns array of utterance objects with `source`, `text`, `start_timestamp`, `end_timestamp`, `confidence`. Returns `null` / 404 when no transcript is available yet — treat as soft failure.

#### Get Panels (Structured Summary Sections)
```
POST https://api.granola.ai/v1/get-document-panels
Body: { "document_id": "<uuid>" }
```
Returns structured AI-generated note sections (action items, decisions, key moments).

#### Get Workspaces
```
POST https://api.granola.ai/v1/get-workspaces
Body: {}
```
Returns array with `id`, `name`, `created_at`, `owner_id`.

#### Get Document Lists (Folders)
```
POST https://api.granola.ai/v2/get-document-lists
Body: {}
```
Returns array with `id`, `title`/`name`, `created_at`, `workspace_id`, `documents`/`document_ids`, `is_favourite`.

### Rate Limits

| Surface | Limit |
|---------|-------|
| REST API | 25 req burst / 5 req/s sustained |
| MCP endpoint | ~100 requests/minute |
| 429 response | Exponential backoff: `Math.min(1000 * 2^attempt, 30000)` ms, max 5 retries |

### Pagination Pattern

```js
function extractDocuments(page) {
  if (Array.isArray(page)) return page;
  return page?.docs ?? page?.documents ?? page?.items
    ?? page?.data?.docs ?? page?.data?.documents ?? [];
}

async function fetchAllDocuments(limit = 250) {
  const pageSize = Math.min(limit, 100);
  const pages = [];
  const seenIds = new Set();
  let offset = 0, totalHint = Infinity;

  while (offset < limit && offset < totalHint) {
    const page = await fetchDocuments(Math.min(pageSize, limit - offset), offset);
    const docs = extractDocuments(page);
    if (!docs.length) break;

    for (const doc of docs) {
      const key = doc.id || `${doc.title}|${doc.updated_at}`;
      if (seenIds.has(key)) continue;
      seenIds.add(key); pages.push(doc);
    }

    const total = page?.total_count || page?.total || page?.pagination?.total || 0;
    if (total > 0) totalHint = total;
    if (docs.length < pageSize) break;
    offset += pageSize; // advance by raw page size, not deduped count
  }
  return pages;
}
```

### Error Handling

| Code | Meaning | Action |
|------|---------|--------|
| 429 | Rate limit hit | Exponential backoff: `Math.min(1000 * 2^attempt, 30000)` ms, max 5 retries |
| 401 | Expired token | Re-read `supabase.json` once (Granola app auto-refreshes); retry once. For API keys, key may be revoked. |
| 403 | Plan not eligible | Surface to user: "Granola API requires Business or Enterprise plan." Do not retry. |
| 404 | No transcript/panels yet | Return `null`, continue — not all meetings have transcripts |
| 5xx | Server error | Retry up to 3x with 2s fixed delay; then surface error |

---

## 8. MCP Integration

### Official Granola MCP Server

**Endpoint:** `https://mcp.granola.ai/mcp`
**Transport:** Streamable HTTP
**Auth:** OAuth 2.0 with Dynamic Client Registration (browser-based sign-in)

#### Available Tools

| Tool | Purpose | Plan Restriction |
|------|---------|-----------------|
| `query_granola_meetings` | Chat with meeting notes, get insights | All plans |
| `list_meeting_folders` | Browse accessible folders | Business/Enterprise |
| `list_meetings` | Scan meeting history with metadata | All plans (Free: 30 days only) |
| `get_meetings` | Search meeting content and notes | All plans |
| `get_meeting_transcript` | Access raw transcripts | Business/Enterprise |
| `get_account_info` | Verify connected account details | All plans |

#### Setup by Client

**Claude Desktop (Paid Plans):**
Settings > Connectors > Search "Granola" > Connect and authenticate

**ChatGPT (Plus/Pro/Business/Enterprise):**
Settings > Search "Granola" > Connect, then use `@Granola` in chat

**Claude Code (Terminal):**
```bash
claude mcp add granola --transport http https://mcp.granola.ai/mcp
```
Then run `/mcp`, select granola, authenticate in browser.

**Any MCP Client:**
Use `https://mcp.granola.ai/mcp` with any client supporting Streamable HTTP transport and OAuth 2.0.

#### Limitations

- Accesses "My notes" space only; team folders not yet supported via MCP
- Free tier: last 30 days only, no transcript access
- Rate limit: ~100 requests/minute

### Community MCP Servers

**GranolaMCP (pedramamini):** Python-based, reads from local cache file (`cache-v3.json`) — no API calls, no network dependency, no auth required. Fast offline access. GitHub: `pedramamini/GranolaMCP`.

**granola-claude-mcp (cobblehillmachine):** Custom MCP connecting Granola data to Claude Desktop for AI-powered meeting intelligence.

**API vs MCP decision:** Use API keys (`grn_*`) for server-side polling, Chrome extension native host bridges, and automation scripts. Use MCP for interactive AI assistant contexts where OAuth browser flow is available.

---

## 9. Integrations

### Native Integrations (Business+ Plans)

**Slack:**
- One-way push from Granola to Slack channels
- Manual per-note sharing or automated folder rules (auto-post when note lands in folder)
- Setup: Settings > Integrations > Connect via OAuth
- Requires Google Workspace or Microsoft 365 accounts (not personal Gmail/Outlook.com)

**Notion:**
- One-way push; creates a Notion database on first connection
- Notes become independent database entries (not standalone pages) — enables filtering, properties, relations
- Manual sharing only (no auto-push); edits in Granola don't update Notion copies
- Setup: Settings > Integrations > Connect via OAuth

**HubSpot:**
- Attach meeting notes to HubSpot contact records
- Auto-suggests relevant contacts during sharing
- One-way push from share menu
- Business and Enterprise plans only

**Affinity:**
- Push to relationship records (Person, Company, Deal)
- Auto-suggests relevant entities
- API key-based authentication
- Popular with VC and sales teams

**Attio:**
- Push to Person, Company, or Deal records
- Workspace-wide setup: one connection benefits entire team
- OAuth-based authentication

### Zapier (Business+ Plans)

Connects to 8,000+ apps including Linear, Asana, Salesforce, Jira, Monday.com.

**Trigger types:**
1. **Note Added to Granola Folder:** Fires automatically when a note lands in a specific folder
2. **Note Shared to Zapier:** Fires when you manually push a note from the note sidebar

**Salesforce workaround:** No native Salesforce integration. Create a Zap with Granola trigger and Salesforce action to sync notes to Account, Opportunity, or Contact records.

### Glean Connector

Granola integrates with Glean enterprise search. The workspace Enterprise API key indexes all Team Space notes; per-user personal API keys index personal notes. Both become searchable via Glean.

---

## 10. Chrome Extension & Native Host Bridge

### Transport Options (Chrome Extension Context)

Three transports, tried in priority order:

| Transport | How | When to use |
|-----------|-----|-------------|
| Native host | `chrome.runtime.sendNativeMessage` to Python bridge | Best: reads `supabase.json` for token automatically |
| Manual token | Direct `fetch` to `api.granola.ai` with user-provided Bearer token | Fallback when native host unavailable |
| MCP bridge | POST to local/remote MCP endpoint | Alternative for MCP-capable environments |

### Native Host Token Path

**macOS:** `~/Library/Application Support/Granola/supabase.json`
**Windows:** `%APPDATA%\Granola\supabase.json`

```json
{ "workos_tokens": { "access_token": "...", "refresh_token": "..." } }
```
The `workos_tokens` value may be a nested JSON string — parse with `JSON.parse()` if `typeof tokens === 'string'`.

**Token TTL:** WorkOS access tokens expire in ~1 hour. The Granola desktop app handles refresh automatically and writes a new token to `supabase.json`. The extension should never call WorkOS directly — just re-read `supabase.json` on each sync cycle or after a 401 response.

### Local Cache File

The Granola desktop app maintains a local cache at `cache-v3.json` in the app data directory. Community tools (GranolaMCP) read this cache directly for offline access. The local SQLite database is encrypted (SQLCipher) — there is no direct DB access.

### Polling / Sync Pattern

Use timestamp-gated incremental sync:

```js
const LAST_SYNC_KEY = 'granola_last_sync';

async function syncGranola({ force = false, limit = 250 } = {}) {
  const lastSync = force ? null : await getLastSyncTime();
  const syncStarted = new Date().toISOString();

  const documents = await fetchAllDocuments(limit);
  if (!documents.length) return { imported: 0, skipped: 0, errors: 0 };

  let imported = 0, skipped = 0, errors = 0;
  for (const doc of documents) {
    const updatedAt = doc.updated_at || doc.last_edited_at || doc.created_at || '';
    if (lastSync && updatedAt && Date.parse(updatedAt) < Date.parse(lastSync)) {
      skipped++; continue;
    }
    try {
      await importDoc(doc);
      imported++;
    } catch (e) { errors++; }
  }

  // Only advance watermark if documents were returned
  if (documents.length > 0) await setLastSyncTime(syncStarted);
  return { imported, skipped, errors };
}
```

### Stable Meeting Key

Use `granola:<doc.id>` as the hash input so updates overwrite rather than duplicate:
```js
const hashInput = docId ? `granola:${docId}` : `granola-fallback:${title}:${createdAt}`;
```

### Account Inference

Match meetings to accounts using:
1. Folder/path hints from `doc.folder`, `doc.breadcrumbs` — matched against per-account folder terms
2. Full-text search across title, summary, people, calendar event
3. Fallback: string matching on participant email domains

### Transport Unavailable Detection

```js
function isTransportUnavailable(errMsg) {
  return /native messaging host|manual granola token is not configured|mcp bridge endpoint is not configured/i
    .test(errMsg);
}
// Soft failure -> disabled: true, reason: <hint>. Hard failure -> throw.
```

---

## 11. Anti-Patterns & Common Mistakes

| Mistake | Fix |
|---------|-----|
| Calling WorkOS token refresh directly | Don't — the Granola app refreshes the token; just re-read `supabase.json` on 401 |
| Re-using a refresh token after exchange | Refresh tokens are single-use; always persist the new token after each exchange |
| Hard-failing on missing transcript | 404 = no transcript yet; return `null`, continue |
| Retrying on 403 | Plan not eligible — surface to user, stop retrying |
| Assuming `doc.docs` always exists | Normalize: check `docs`, `documents`, `items`, `data.docs` — use `extractDocuments()` |
| Treating all speakers as identified | Only `audio_source_type` reliably separates local mic from remote; named speakers unreliable with 3+ people |
| Not deduplicating paginated results | API may return overlapping items; track `seenIds` Set by doc ID |
| Using raw page count for offset | Advance offset by `pageSize` (raw), not deduped count — avoids premature loop exit |
| Advancing sync watermark on empty response | Only advance when `documents.length > 0` |
| Using local DB directly | Granola's SQLite DB is SQLCipher-encrypted — API only |
| Expecting video/audio playback | Granola does not store audio or video; text-only output |
| Assuming MCP accesses team folders | Official MCP accesses "My notes" only; team folders not yet supported |
| Using personal Gmail/Outlook.com for team features | Slack integration, template sharing require Google Workspace or Microsoft 365 accounts |
| Importing pre-recorded audio | Granola only supports live meeting transcription; no file upload/import |

---

## 12. Privacy, Security & Compliance

### Data Handling

- **No bot joins meetings** — audio captured locally; participants do not see an attendee
- **Audio never stored:** Streamed to transcription provider in real time, discarded after transcription. No playback capability.
- **Data at rest:** US-hosted AWS VPC, encrypted at rest and in transit, daily backups
- **Notes private by default:** Must explicitly share via folder or direct share

### Certifications

- **SOC 2 Type 2** certified (July 2025, independently audited)
- **GDPR** compliant with Data Processing Agreement available
- **Not HIPAA compliant** — verify before healthcare use

### Model Training

- Third-party AI providers (OpenAI, Anthropic) cannot use customer data for model training
- Granola trains on anonymized data only; opt-out available in Settings
- Enterprise: model training disabled by default; org-wide opt-out

### Enterprise Controls

- SSO enforcement (50+ users)
- Configurable transcript auto-deletion periods
- Admin controls for meeting link sharing
- Data export and note-transfer permissions
- Organization-wide notification of Granola usage
- API access controls per workspace

### Security Advisory History

**TRA-2025-07 (Tenable, March 2025):** Hard-coded AssemblyAI API key exposed in unauthenticated `get-feature-flags` endpoint. Scope: 333 beta testers of iOS TestFlight build only. macOS production unaffected (uses different transcription provider). Remediated same day: key revoked, endpoint secured. No unauthorized access confirmed beyond Tenable's controlled test.

### Consent

- Granola requires manual activation per meeting — does not auto-join or auto-record
- macOS: "Let others know you're using Granola" sends an auto-consent message to participants
- Users are responsible for notifying participants per local recording consent laws
- Wrap untrusted content in tagged envelopes before LLM processing:
  ```
  <untrusted_transcript>{{ transcript_text }}</untrusted_transcript>
  ```

---

## 13. Troubleshooting

### Transcription Issues

| Problem | Solution |
|---------|----------|
| No transcription at all | Restart Granola; verify Microphone + Screen & System Audio Recording permissions (macOS) or desktop app microphone access (Windows) |
| Transcription stops after 4-5 minutes | Check: Bluetooth/USB device stability, computer sleep settings; try built-in audio |
| Poor audio quality / garbled text | Verify default audio device matches across system and meeting app; disable USB amplifiers or audio mixing software |
| Multi-language transcription inaccurate | Internal jargon not supported with multi-language mode; disable if single-language meeting |
| VDI/Citrix: only captures microphone | Install Granola locally instead of in virtual environment |

### Authentication Errors

| Problem | Solution |
|---------|----------|
| 401 Unauthorized | Re-read `supabase.json` (Granola app auto-refreshes on token expiry); if using API key, verify key is not revoked |
| 403 Forbidden | Plan not eligible for API access; upgrade to Business ($14/user/month) or Enterprise |
| MCP auth fails | Ensure browser can open for OAuth flow; try disconnecting and reconnecting in AI tool settings |
| "Native messaging host not found" | Run `install.sh <extension-id>` from native-host directory; restart Chrome |

### Integration Issues

| Problem | Solution |
|---------|----------|
| Slack: "not available for personal accounts" | Use Google Workspace or Microsoft 365 email, not personal Gmail/Outlook.com |
| Notion: edits not syncing | One-way push only; edits in Granola don't update Notion copies after sharing |
| Zapier: missing workspace selector | Update Zap to latest version; select workspace, save and test |
| HubSpot/Attio: token expired | Disconnect and reconnect in Settings > Integrations |
| Affinity: auth failed | Regenerate API key in Affinity; paste new key in Settings > Integrations |
| Calendar not syncing | Sign out and sign in to re-authenticate calendar connection; check that events contain video conferencing URLs |
| MCP: "only returns last 30 days" | Free tier limitation; upgrade to Business for full history |

### Network Issues

- Firewall/VPN/proxy may block transcription — check `status.granola.ai` for service incidents
- Corporate networks: may need IT to whitelist Granola domains (TLS inspection can break connections)
- Try switching DNS to `8.8.8.8` as diagnostic step

### Support Contact

Email: hey@granola.so — include app version, OS details, integration name, failure time, and screenshots.

---

## 14. Product Milestones (Changelog)

| Date | Milestone |
|------|-----------|
| March 2023 | Company founded (Chris Pedregal, Sam Stephenson) |
| July 2024 | Slack integration launched |
| February 2025 | Notion export; MCP server launched |
| April 2025 | iOS app launched; Granola Chat rebuilt (agentic, recipes, citations) |
| May 2025 | Granola 2.0 (team folders, Spaces, Briefs); Windows support |
| June 2025 | Windows public launch |
| July 2025 | SOC 2 Type 2 certified; team folders; Zapier integration |
| September 2025 | Recipes launched; "Shared with Me" view |
| December 2025 | @Mentions functionality |
| January 2026 | Microsoft sign-in (native Teams/Outlook) |
| March 2026 | Series C ($125M, $1.5B valuation); personal API + enterprise API launched |

---

## 15. Quick Reference

| Task | Endpoint / Method |
|------|-------------------|
| List recent meetings | `POST /v2/get-documents` |
| Get transcript | `POST /v1/get-document-transcript` |
| Get summary panels | `POST /v1/get-document-panels` |
| Get documents by ID | `POST /v1/get-documents-batch` |
| Get workspaces | `POST /v1/get-workspaces` |
| Get folders | `POST /v2/get-document-lists` |
| Auth token source (native) | macOS: `~/Library/Application Support/Granola/supabase.json`; Windows: `%APPDATA%\Granola\supabase.json` |
| API key format | `grn_YOUR_API_KEY` |
| MCP endpoint | `https://mcp.granola.ai/mcp` |
| Rate limit (API) | 25 req burst / 5 req/s sustained |
| Rate limit (MCP) | ~100 req/min |
| 401 response | Re-read `supabase.json` or verify API key |
| 403 response | Plan not eligible — surface to user, do not retry |
| 404 on transcript/panels | No transcript yet — return `null`, continue |
| Status page | `https://status.granola.ai` |
| Support email | hey@granola.so |

---

## 16. Sources

1. [Granola official website](https://www.granola.ai/)
2. [Granola pricing page](https://www.granola.ai/pricing)
3. [Granola security page](https://www.granola.ai/security)
4. [Granola API docs](https://docs.granola.ai/introduction)
5. [Granola MCP docs](https://docs.granola.ai/help-center/sharing/integrations/mcp)
6. [Granola MCP blog post](https://www.granola.ai/blog/granola-mcp)
7. [Granola templates docs](https://docs.granola.ai/help-center/taking-notes/customise-notes-with-templates)
8. [Granola transcription docs](https://docs.granola.ai/help-center/taking-notes/transcription)
9. [Granola workspaces docs](https://docs.granola.ai/help-center/workspaces)
10. [Granola multi-language support](https://help.granola.ai/article/multi-language)
11. [Granola integrations guide](https://www.granola.ai/blog/granola-integrations-complete-guide-connecting-meeting-tools)
12. [Granola integration troubleshooting](https://www.granola.ai/blog/granola-integration-troubleshooting-common-issues-solutions)
13. [Granola transcription troubleshooting](https://docs.granola.ai/help-center/troubleshooting/transcription-issues)
14. [Granola Recipes blog](https://www.granola.ai/blog/say-hello-to-recipes)
15. [Granola templates & recipes blog](https://www.granola.ai/blog/meeting-recipes-repeatable-formats)
16. [Granola 2.0 blog](https://www.granola.ai/blog/two-dot-zero)
17. [Granola Spaces launch](https://www.createwith.com/tool/granola/updates/granola-introduces-spaces-for-team-wide-meeting-intelligence)
18. [Granola Series C (TechCrunch)](https://techcrunch.com/2026/03/25/granola-raises-125m-hits-1-5b-valuation-as-it-expands-from-meeting-notetaker-to-enterprise-ai-app/)
19. [Granola Series C blog](https://www.granola.ai/blog/series-c)
20. [Granola product updates/changelog](https://www.granola.ai/updates)
21. [Reverse-engineered Granola API (getprobo)](https://github.com/getprobo/reverse-engineering-granola-api)
22. [Reverse engineering Granola for Obsidian (Joseph Thacker)](https://josephthacker.com/hacking/2025/05/08/reverse-engineering-granola-notes.html)
23. [Tenable TRA-2025-07 advisory](https://www.tenable.com/security/research/tra-2025-07)
24. [Granola post-mortem: AssemblyAI key exposure](https://docs.granola.ai/help-center/policies/security-reports/post-mortem-assembly-ai-api-key-exposure)
25. [GranolaMCP community server](https://github.com/pedramamini/GranolaMCP)
26. [Granola encrypted DB context](https://www.shadow.do/blog/granola-encrypted-its-local-database-heres-why-that-matters----and-what-to-use-instead)
27. [Granola Glean connector docs](https://docs.glean.com/connectors/native/granola/)
28. [Granola Zapier integrations](https://zapier.com/apps/granola/integrations)
29. [Granola vs Otter vs Fireflies vs Fathom (2026)](https://www.useluminix.com/reports/industry-analysis/ai-meeting-notes-comparison-granola-vs-otter-vs-fireflies-vs-fathom-2026)
30. [Zapier: What is Granola?](https://zapier.com/blog/granola-ai/)
31. [Granola pricing analysis (get-alfred)](https://get-alfred.ai/blog/granola-pricing)
32. [Granola pricing benchmarks (costbench)](https://costbench.com/software/ai-meeting-assistants/granola/)
33. Repo implementation: `native-host/granola_host.py`, `src/background/granola.js`

## Related Skills

- **plaud-integration** — Plaud AI recorder hardware and software ecosystem (complementary meeting capture tool with hardware devices)
- **using-plaud-mcp** — Quick-start for the Plaud MCP server
- **chrome-extension-expert** — Native host / chrome.runtime.sendNativeMessage protocol patterns used by Granola's native bridge
- **integration-clients** — Shared third-party API client patterns (auth, retry, pagination) for wiring Granola data downstream
