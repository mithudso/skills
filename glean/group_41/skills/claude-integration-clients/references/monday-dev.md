<!-- hub-reference-banner -->
> **Reference file — part of the `integration-clients` hub.** Formerly the standalone `monday-dev` skill.
> Sibling topics in this family are now reference files under the hubs (`integration-clients`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: monday-dev
version: 2.1.0
last_updated: 2026-05-29
description: >-
  Monday.com platform architecture, Apps Framework, CLI (mapps), Vibe Design
  System, column types, board views, automations, WorkDocs, dashboards, AI
  features, pricing, and marketplace publishing.
  TRIGGER: user builds a monday app; uses mapps CLI; chooses column types or
  board views; configures automations or workflow builder; publishes to the
  monday marketplace; uses Vibe components or @vibe/mcp; asks about platform
  hierarchy, pricing tiers, or monday dev/CRM/service products.
  SKIP: GraphQL queries, column value JSON formats, pagination, webhooks, rate
  limits, or MCP server config → use monday-api instead.
category: developer
tags: [monday, platform, apps-framework, cli, mapps, vibe, automations, ai, marketplace, developer]
whenToUse:
  - building or scaffolding a monday.com app with the Apps Framework
  - using the mapps CLI (app:create, code:push, app:promote, code:logs)
  - choosing between board views (Kanban, Gantt, Timeline, Chart, etc.)
  - selecting or comparing column types for a board design
  - configuring automations or the workflow builder
  - using Vibe Design System components (@vibe/core, @vibe/icons, @vibe/mcp)
  - publishing or updating an app in the monday marketplace
  - choosing between monday work management, CRM, dev, or service products
  - understanding monday AI credits, AI columns, or Agent Factory
  - deploying an app to monday Code hosting (Cloud Run)
  - setting up OAuth for a monday app
  - using the monday-sdk-js or @mondaycom/apps-sdk
whenNotToUse:
  - GraphQL queries or mutations against api.monday.com/v2 → monday-api
  - column value JSON formats, webhooks, or rate limits → monday-api
  - monday MCP server configuration → monday-api
  - Chrome extension integrating with monday → monday-api
related_skills:
  - monday-api
  - monday-board-audit
---

# Monday.com Platform, CLI, and Apps Framework

## When to use this skill

Use when working with the monday.com platform, building monday apps, using the monday CLI (`mapps`), choosing monday views or column types, configuring automations, or publishing to the monday marketplace. For GraphQL API queries/mutations, column value JSON formats, webhooks, rate limiting, and MCP server configuration, use the companion `monday-api` skill instead.

---

## 1. Platform Architecture

### 1.1 Structural hierarchy

```
Account
  └── Workspaces (logical containers)
        ├── Folders (optional grouping)
        │     └── Sub-folders
        ├── Boards (rich tables — the primary building block)
        │     ├── Groups (row sections within a board)
        │     │     └── Items (individual rows)
        │     │           └── Subitems (child rows)
        │     ├── Columns (typed data fields — 40+ types)
        │     └── Views (table, kanban, gantt, etc.)
        ├── Dashboards (multi-board aggregation with widgets)
        ├── WorkDocs (collaborative documents)
        └── WorkForms (form-to-board data collection)
```

- **Workspace** is the top-level container. Every board, dashboard, and WorkDoc lives inside a workspace.
- **Boards** come in three visibility types: Main (all workspace members), Private (invited members only), Shareable (external guests).
- **Subitems** add a second layer under items. Multi-level boards (2025-10+) support up to 5 layers.
- **Groups** are collapsible sections that organize items visually.

### 1.2 Board types

| Type | Description |
|---|---|
| Classic board | Standard item/group/column structure |
| Multi-level board | Up to 5 layers of subitems (2025-10+, initially Portfolio only) |
| High-level board | Overview board linking to low-level detail boards |
| Low-level board | Detailed task board linked from a high-level board |

### 1.3 Products

| Product | Purpose |
|---|---|
| **monday work management** | Project and task management (the core product) |
| **monday CRM** | Sales pipeline, deal tracking, customer management |
| **monday dev** | Software development sprints, bug tracking, releases |
| **monday service** | IT service management, ticketing |
| **monday campaigns** | AI-powered marketing campaign management (2025+) |

---

## 2. Column Types

### 2.1 Essential columns

| Column | API Type ID | Description |
|---|---|---|
| Status | `status` | Color-coded labels tracking item state |
| People | `people` | Assign team members or teams |
| Numbers | `numbers` | Numeric values with optional units |
| Text | `text` | Free-form text content |
| Date | `date` | Single date (optional time) |
| Timeline | `timeline` | Date range (from/to) for durations |
| Connect Boards | `board_relation` | Link items across boards |
| Mirror | `mirror` | Reflect data from connected board columns |
| Monday Doc | `doc` | Embedded document column |

### 2.2 Common columns

| Column | API Type ID | Description |
|---|---|---|
| Checkbox | `checkbox` | Boolean checked/unchecked |
| Dropdown | `dropdown` | Single or multi-select from predefined labels |
| Link | `link` | URL with display text |
| Email | `email` | Email address with display text |
| Phone | `phone` | Phone number with country code |
| Location | `location` | Lat/lng coordinates with address |
| Files | `file` | File attachments |
| Tags | `tags` | Shared labels across boards |
| Rating | `rating` | Star rating (1-5) |
| Vote | `vote` | Team voting on items |

### 2.3 Time and tracking columns

| Column | API Type ID | Description |
|---|---|---|
| Hour | `hour` | Time of day (HH:MM) |
| Week | `week` | Calendar week selection |
| Time Tracking | `time_tracking` | Start/stop timer with logged hours |
| Creation Log | `creation_log` | Auto-records who created the item and when |
| Last Updated | `last_updated` | Auto-records last modification |
| Auto-Number | `auto_number` | Sequential numbering |

### 2.4 Advanced columns

| Column | API Type ID | Description |
|---|---|---|
| Formula | `formula` | Calculated values from other columns |
| Dependencies | `dependency` | Task dependency relationships |
| Progress Tracking | `progress` | Visual completion percentage |
| Button | `button` | Triggers custom actions |
| World Clock | `timezone` | Time zone display |

### 2.5 AI-powered columns (2025+)

All AI columns consume AI credits from the account pool.

| Column | Purpose |
|---|---|
| Write with AI | Generate content from prompts |
| Summarize | Condense lengthy content |
| Summarize Updates | Condense item update threads |
| Translate | Translate text between languages |
| Assign Labels | Auto-categorize with AI |
| Sentiment Detection | Analyze emotional tone |
| Extract from File | Pull data from attached documents |
| Custom AI Prompt | Execute personalized AI requests |
| Prioritize with AI | Rank items by importance |

---

## 3. Board Views

| View | Description | Best for |
|---|---|---|
| Table | Default spreadsheet-style grid | General data management |
| Kanban | Cards grouped by status column | Workflow visualization |
| Timeline | Horizontal bar chart of date ranges | Project scheduling |
| Gantt | Waterfall-style timeline with dependencies | Project management |
| Calendar | Date-based event layout | Scheduling and deadlines |
| Chart | Bar, pie, and other graph types | Data visualization |
| Workload | Team capacity and task distribution | Resource planning |
| Map | Geographic pin placement | Location-based data |
| Form | Input collection tied to board items | Data collection |
| Files Gallery | Grid of file attachments | Media management |
| Cards | Visual card representation of items | Quick overviews |
| Pivot Boards | Cross-tabulation analysis | Data analysis |

Additional views available through the Apps Marketplace. Dashboards with AI capabilities require Pro or Enterprise plans.

---

## 4. Automations and Integrations

### 4.1 Automation anatomy

Every automation consists of:
1. **Trigger** — the event that starts it (column change, item creation, date arrival)
2. **Condition** — optional filter the trigger must meet
3. **Action** — what happens (move item, notify, update column, send email)

### 4.2 Automation builder (2025+)

- Dynamic data from triggers reusable across all pickers
- Both items and subitems in the same automation
- Cross-board automations linking boards together
- Custom automation templates saveable for account-wide reuse

### 4.3 Automation limits by plan

| Plan | Monthly actions |
|---|---|
| Basic | 0 (no automations) |
| Standard | 250 |
| Pro | 25,000 |
| Enterprise | 250,000 |

### 4.4 Integration types

Built-in integrations: Slack, Gmail, Google Calendar, Microsoft Teams, Outlook, Jira, GitHub, Salesforce, HubSpot, Zapier, and 200+ more via the marketplace.

---

## 5. WorkDocs, WorkForms, and Dashboards

### 5.1 WorkDocs

Collaborative real-time documents within the workspace:
- Real-time co-editing with comments and mentions
- Embedded boards, dashboards, videos, and rich media
- Text-to-item conversion: highlight text and click "+ item" to create board items
- Version history for tracking changes

### 5.2 WorkForms

Customizable forms that feed submissions into boards:
- Brand-customizable design (logo, colors, cover image)
- Shareable via link (no monday account needed)
- Conditional visibility, validation rules, and required fields
- Password protection available

### 5.3 Dashboards

Multi-board data aggregation:
- 15 built-in widgets: numbers, chart, timeline, battery, workload, and more
- Custom widgets via Dashboard Widget app features
- New widget types (2026-04+): APP_FEATURE (embed custom app functionality) and LISTVIEW (cross-board tabular list)

---

## 6. Monday AI Features

### 6.1 AI capabilities

| Feature | Description |
|---|---|
| **monday sidekick** | AI assistant with full business context |
| **AI Blocks** | Building blocks for columns, automations, and workflows |
| **Agent Factory** | Create personalized AI agents that work autonomously |
| **monday magic** | Transform natural language prompts into ready-to-use work solutions |
| **monday vibe** | AI app builder — create apps without writing code |
| **Formula Builder** | Generate formula column expressions from natural language |

### 6.2 AI credits model (2026+)

| Plan | Minimum monthly credits |
|---|---|
| Basic | 1,000 |
| Standard | 2,000 (up to 8,000) |
| Pro | 3,000 (up to 20,000) |
| Enterprise | Custom |

AI credits are mandatory for new customers as of May 2026. Enable AI: Admin > Permissions > AI Connectors.

---

## 7. Pricing Tiers (monday work management)

| Plan | Price/seat/mo (annual) | Key features |
|---|---|---|
| **Free** | $0 | Up to 2 seats, 3 boards, 200 items, no automations |
| **Basic** | ~$9 | Unlimited boards and items, no automations or integrations |
| **Standard** | ~$12 | Automations (250/mo), integrations, timeline/Gantt, guest access |
| **Pro** | ~$19 | Automations (25K/mo), time tracking, formula, dependencies |
| **Enterprise** | Custom | Automations (250K/mo), advanced permissions, SSO, HIPAA |

- Minimum 3 seats on paid plans (multiples of 5)
- Annual billing: 18% discount over monthly
- Each product (work management, CRM, dev, service) has its own pricing tiers

**Developer accounts:** Free, up to 10 seats, 25K monthly automation actions, Pro/Enterprise features for development.

---

## 8. Monday CLI (`mapps`) — Complete Reference

### 8.1 Installation and initialization

```bash
npm install -g @mondaycom/apps-cli
mapps init -t <API_TOKEN>     # creates .mappsrc config file
```

### 8.2 App management

```bash
mapps app:create -n "My App" [-d TARGET_DIR]
mapps app:list
mapps app:scaffold ./my-app <TEMPLATE_NAME>
mapps app:deploy [-a APP_ID] [-d DIRECTORY] [-v VERSION_ID] [-f] [-z REGION]
mapps app:promote -i <VERSION_ID> -a <APP_ID>
```

### 8.3 Monday Code hosting

```bash
mapps code:push -d <BUILD_DIR> -i <VERSION_ID> [-a APP_ID] [-f] [-c] [-s] [-z REGION]
mapps code:status -i <VERSION_ID> [-z REGION]
mapps code:logs -i <VERSION_ID> -t <TYPE> [-r REGEX] [-f START] [-e END] [-z REGION]
mapps code:env -m [list-keys|set|delete] -i APP_ID [-k KEY] [-v VALUE]
mapps code:secret -m [list-keys|set|delete] -i APP_ID [-k KEY] [-v VALUE]
```

### 8.4 Scheduler

```bash
mapps scheduler:create -a APP_ID -s "0 * * * *" -u "/endpoint" -n "job-name"
mapps scheduler:list -a APP_ID
mapps scheduler:run -a APP_ID -n "job-name"
mapps scheduler:delete -a APP_ID -n "job-name"
```

### 8.5 Features, versions, and storage

```bash
mapps app-features:create -a APP_ID -i VERSION_ID -n "feature" -t FEATURE_TYPE
mapps app-features:list [-a APP_ID] [-i VERSION_ID]
mapps app-version:list -i APP_ID
mapps storage:export [-a APP_ID] [-c CLIENT_ACCOUNT_ID] [-d FILE_PATH]
mapps storage:search [-a APP_ID] [-c CLIENT_ACCOUNT_ID] [-t SEARCH_TERM]
mapps manifest:export [-a APP_ID] [-i VERSION_ID]
mapps manifest:import -p MANIFEST_PATH [-a APP_ID]
mapps tunnel:create -a APP_ID [-p PORT]
```

### 8.6 Regional deployment

The `-z` flag supports regions: `us`, `au`, `eu`, `il`. IL region deployment is mandatory for Live promotion (2025+).

---

## 9. Apps Framework

### 9.1 App feature types

| Feature Type | Placement |
|---|---|
| **Board View** | Views Center tab on boards |
| **Item View** | Updates section of an item |
| **Board Menu** | Board context menu |
| **Dashboard Widget** | Dashboard widget picker |
| **Custom Object** | Left-pane workspace menu |
| **Account Settings View** | Platform-wide settings |
| **Doc Action** | WorkDoc toolbar |
| **AI Assistant** | Various AI surfaces |
| **Integration** | Automation center |
| **Workspace Template** | Template gallery (requires at least one other feature) |

### 9.2 Monday Code hosting infrastructure

| Feature | Description |
|---|---|
| Key-value Storage | Per-customer data persistence |
| Secure Storage | Encrypted data management |
| Document DB | Structured data (MongoDB-compatible) |
| Message Queue | Decouple services for scalability |
| Cron Scheduling | Scheduled job execution |
| Logging | Built-in Logger SDK, accessible via CLI |
| Security Scanning | CVE detection during deployment |
| Multi-region | US, AU, EU, IL |

**Supported languages:** Node.js and Python have official SDKs. Additional templates for C#, Django, Go, Java, PHP, and Ruby.

### 9.3 App lifecycle

```
Draft version  → develop & test → promote →  Live version
  ├── Create app (app:create)
  ├── Add features (app-features:create)
  ├── Deploy code (code:push)
  ├── Connect feature builds (app-features:build)
  ├── Set env vars + secrets (code:env, code:secret)
  ├── Configure manifest + OAuth scopes
  ├── Test with active version
  └── Promote when ready (app:promote)
```

**Key rules:**
- Always build and deploy on draft versions
- Never push directly to live
- IL region deployment is mandatory for Live promotion (2025+)

### 9.4 OAuth flow for apps

1. Configure scopes and redirect URLs in Developer Center
2. Redirect user to `https://auth.monday.com/oauth2/authorize?client_id=...&redirect_uri=...&state=...`
3. User approves; monday redirects with a temporary authorization code (valid 10 minutes)
4. Your backend exchanges the code at `https://auth.monday.com/oauth2/token`
5. Access token grants API access scoped to the user until they uninstall

**Friction-free install (2025+):** Add `force_install_if_needed=true` to the auth URL.

### 9.5 monday-sdk-js (client-side)

```bash
npm install monday-sdk-js
```

Current version: 0.5.8. Provides seamless authentication in embedded contexts (iframes).

```js
import mondaySdk from 'monday-sdk-js';
const monday = mondaySdk();
monday.setApiVersion('2025-07');

// Get board context
const { data } = await monday.get('context');
console.log(data.boardId, data.itemId);

// Authenticated API call (no manual token needed)
const res = await monday.api(`{ boards(ids: [${data.boardId}]) { name } }`);

// Session token for backend verification
const { data: token } = await monday.get('sessionToken');
```

**Note:** The server-side GraphQL client in `monday-sdk-js` is deprecated. Use `@mondaydotcomorg/api` for server-side GraphQL queries.

### 9.6 Required OAuth scopes

Common scopes: `boards:read`, `boards:write`, `items:read`, `items:write`, `users:read`, `webhooks:write`, `account:read`, `docs:read`, `docs:write`, `assets:read`.

Declare every scope the app uses in the manifest — missing scopes cause runtime permission errors.

### 9.7 Marketplace publishing

Review dimensions: documentation completeness, legal compliance, UI/UX design (use Vibe), product functionality, workspace template requirements (requires at least one additional feature type).

Monetization: subscription-based pricing, discount and trial extensions, analytics dashboard.

---

## 10. Vibe Design System

```bash
npm install @vibe/core @vibe/icons @vibe/style
```

```tsx
import '@vibe/core/tokens';
import { Button, TextField, Flex, Icon } from '@vibe/core';
import { Search } from '@vibe/icons';

<Button kind="primary">Primary</Button>
<Button kind="secondary">Secondary</Button>
<Icon icon={Search} size={24} label="Search" />
```

**MCP access:** `@vibe/mcp` (v4.0.0) exposes component metadata, examples, accessibility guidance, icons, tokens, and migration analyzers (v3-migration, v4-migration, dropdown-migration).

**Key rule:** Use Vibe components instead of custom UI primitives. Monday's app review process favors Vibe for UI consistency.

---

## 11. MCP Surfaces

| Surface | Package / Mode | Use for |
|---|---|---|
| **Platform MCP** | `@mondaydotcomorg/monday-api-mcp` (default mode) | Board/item/doc/form/dashboard CRUD, user/workspace operations |
| **Apps MCP** | `@mondaydotcomorg/monday-api-mcp --mode apps` | App lifecycle: create, version, feature, deploy, promote |
| **Vibe MCP** | `@vibe/mcp` | Component metadata, examples, accessibility, icons, tokens, migration |

MCP is preinstalled on all monday.com accounts. Admin verification: Admin > Permissions > AI Connectors.

---

## 12. Decision Guide

1. **Operating on monday data** (boards, items, docs, forms) → Platform MCP / GraphQL API (`monday-api` skill)
2. **Building or managing monday apps** → Apps MCP + CLI (`mapps`) + Developer Center
3. **Building UI in a monday app** → Vibe MCP + `@vibe/core` components
4. **Automating workflows** → Automations / Workflow Builder / Integration recipes

**Common mistakes:**
- Using Apps MCP for board/item CRUD (use Platform MCP instead)
- Using Platform MCP for app version/deploy work (use Apps MCP instead)
- Guessing CLI commands without checking `monday_apps_get_development_context`
- Building custom UI primitives instead of using Vibe components
- Pushing code directly to live versions (always draft → promote)
- Forgetting OAuth scope declarations (causes runtime permission errors)
- Missing IL region deployment before Live promotion

---

## 13. Sources

1. [Understanding monday.com's structural hierarchy](https://support.monday.com/hc/en-us/articles/7278527605906)
2. [Available column types](https://support.monday.com/hc/en-us/articles/115005310285)
3. [Board views](https://support.monday.com/hc/en-us/articles/360001267945)
4. [Get started with monday AI](https://support.monday.com/hc/en-us/articles/11512670770834)
5. [Plans and pricing](https://support.monday.com/hc/en-us/articles/4405633151634)
6. [CLI documentation](https://developer.monday.com/apps/docs/command-line-interface-cli)
7. [What is a monday app?](https://developer.monday.com/apps/docs/intro)
8. [monday SDK introduction](https://developer.monday.com/apps/docs/introduction-to-the-sdk)
9. [OAuth and Permissions](https://developer.monday.com/apps/docs/oauth)
10. [Get started with monday code hosting](https://developer.monday.com/apps/docs/get-started)
11. [Vibe design system docs](https://developer.monday.com/apps/docs/vibe-design-system)
12. [monday.com MCP GitHub](https://github.com/mondaycom/mcp)
13. [Automations overview](https://support.monday.com/hc/en-us/articles/360001222900)
14. [Get started with monday workdocs](https://support.monday.com/hc/en-us/articles/360021702939)
15. [Developer plans and pricing](https://developer.monday.com/apps/docs/plans-and-pricing)
