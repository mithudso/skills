<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `alarm-scheduler-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: alarm-scheduler-patterns
description: >
  Chrome MV3 extension alarm scheduling expert: tiered severity-based polling, cooldowns,
  debouncing, dynamic backoff (exponential + jitter), tab deduplication, persistence across
  service-worker restarts, and missed-alarm burst suppression. Covers the metronome pattern
  (single alarm + per-resource interval checks) and all chrome.alarms API constraints.
  TRIGGER: user mentions chrome.alarms, alarm scheduling, or polling in Chrome extensions;
  designing periodic background tasks in MV3; tiered refresh or severity-based polling;
  debugging alarms not firing or alarm persistence issues; exponential backoff in extensions;
  setTimeout/setInterval failing after service worker restart.
  SKIP: Node.js cron jobs or server-side scheduling (use job-scheduling-patterns); non-Chrome
  browser extension alarms; general service worker questions not involving chrome.alarms.
version: 1.2.0
updated: "2026-05-29"
category: developer
tags:
  - chrome-extension
  - alarms
  - scheduling
  - polling
  - mv3
  - cooldown
  - service-worker
  - tiered-refresh
keywords:
  - chrome.alarms
  - alarm scheduling
  - chrome extension
  - MV3
  - service worker
  - tiered refresh
  - severity-based polling
  - exponential backoff
  - jitter
  - cooldown
  - debounce
  - tab deduplication
  - metronome pattern
  - onInstalled
  - onStartup
  - alarm persistence
  - burst suppression
whenToUse:
  - User mentions chrome.alarms, alarm scheduling, or polling in a Chrome extension
  - Designing periodic background tasks in a Manifest V3 extension
  - Implementing tiered refresh or severity-based polling intervals
  - Debugging alarms not firing, alarm persistence issues, or SW restart problems
  - Implementing exponential backoff or dynamic interval scheduling in an extension
  - Preventing duplicate alarm-triggered work across multiple tabs
whenNotToUse:
  - Node.js cron jobs or server-side scheduled tasks — use job-scheduling-patterns
  - General MV3 service worker questions without a scheduling/alarm angle — use mv3-service-worker-expert
  - Non-Chrome browser extension APIs
related_skills:
  - mv3-service-worker-expert
  - chrome-storage-patterns
  - chrome-dev
  - ops-registry-patterns
  - job-scheduling-patterns
globs:
  - "**/service-worker*.js"
  - "**/background*.js"
  - "**/alarm*.js"
  - "**/scheduler*.js"
  - "**/tiered-refresh*.js"
---

# Alarm Scheduler Patterns

## Core Constraints (Chrome 120+)

- Minimum interval: **30 seconds** (`periodInMinutes: 0.5`)
- Maximum active alarms per extension: **500**
- Alarms may be cleared on browser restart — always re-register in `onInstalled` AND `onStartup`
- Creating an alarm with an existing name replaces it (built-in debounce)
- Missed alarms during device sleep fire at most once on wake

**Required manifest permission:**
```json
{ "permissions": ["alarms", "storage"] }
```

## Decision Tree

| Problem | Solution |
|---|---|
| Need periodic work that survives SW restart | `chrome.alarms` — never `setInterval` |
| Many resources with different urgency | Metronome + per-resource interval check |
| Too many resources for per-resource alarms (> ~100) | Metronome (1 alarm for all) |
| Alarm fires but work duplicates | Cooldown gate with `chrome.storage.session` |
| Need sub-30s polling temporarily | `setTimeout` inside alarm handler + storage keep-alive |
| Multiple tabs watching same resource | Centralized registry + single fetch + broadcast |
| Alarm disappears after browser restart | Dual registration: `onInstalled` + `onStartup` |
| API is down — avoid hammering | Exponential backoff with jitter |

## Key Patterns (Summaries)

### Metronome vs Per-Resource Alarms

Use a **single metronome alarm** (e.g., every 1 minute) that checks per-resource intervals on each tick, rather than one alarm per resource. Reasons: 500-alarm limit, batching efficiency, severity-based priority sorting.

```javascript
// One alarm drives everything
await chrome.alarms.create('mca-tiered-refresh', { periodInMinutes: 1 });
// On each tick, check which resources are due (see context file for full implementation)
```

### Naming Conventions

| Pattern | Example | Use Case |
|---|---|---|
| `{ext}-{purpose}` | `mca-prune-session` | Static singleton alarm |
| `{feature}:{action}:{id}` | `firedrill:drip:scenario-3` | Dynamic per-resource alarm |
| `{ext}-{tier}-refresh` | `mca-tiered-refresh` | Tiered polling dispatch |

### Top-Level Listener (Critical Rule)

Listeners **must** be registered at the top level of the service worker, never inside async functions or conditionals. Chrome uses this registration to decide whether to wake the worker.

```javascript
// CORRECT
chrome.alarms.onAlarm.addListener(async (alarm) => { ... });

// WRONG — may miss events after SW restart
async function setup() { chrome.alarms.onAlarm.addListener(handler); }
setup();
```

### Persistence Pattern

```javascript
async function ensureAlarms() {
  await chrome.alarms.create('mca-tiered-refresh', { periodInMinutes: 1 });
  // add other alarms here
}

chrome.runtime.onInstalled.addListener(async () => { await ensureAlarms(); });
chrome.runtime.onStartup.addListener(async () => { await ensureAlarms(); });
```

### Storage Choice for Scheduling State

| Store | Survives SW Death | Survives Restart | Use For |
|---|---|---|---|
| `chrome.storage.session` | Yes | No | Cooldown timestamps, locks, transient tab maps |
| `chrome.storage.local` | Yes | Yes | Alarm configs, retry counts, scheduling metadata |
| In-memory variable | No | No | Per-execution concurrent-run guards only |

## Anti-Patterns

1. **`setInterval`/`setTimeout` for scheduled work** — cancelled on SW termination; use `chrome.alarms`.
2. **Per-resource alarms at scale** — hits 500-alarm limit; use the metronome pattern.
3. **Relying on `onInstalled` alone** — alarms cleared on restart; always add `onStartup` registration.
4. **Conditional listener registration** — Chrome won't wake the SW for listeners it never saw registered; always register unconditionally, execute conditionally.
5. **No jitter in multi-instance setups** — all instances fire simultaneously; add randomized `delayInMinutes`.
6. **Ignoring wake-burst after sleep** — check a `lastRun` timestamp before processing to suppress duplicates.

## Troubleshooting

| Symptom | Check |
|---|---|
| Alarm not firing | `"alarms"` in `manifest.json`? Interval ≥ 0.5 min? Listener at top level? Alarm count < 500? |
| Handler runs but does nothing | Log `alarm.name` — name mismatch? Cooldown active? Lock held? Error swallowed? |
| Alarms disappear after update | Extension updates clear all alarms — handle `onInstalled` with `details.reason === 'update'` |
| Need faster than 30s | Use `setTimeout` inside the alarm handler while SW is awake; keep alive with `storage.session.set` |

## Full Reference

For all code implementations — tiered interval calculation, cooldown/lock patterns, tab deduplication, exponential backoff, burst mode, persistence/rehydration — read `references/alarm-scheduler-context.md`.
