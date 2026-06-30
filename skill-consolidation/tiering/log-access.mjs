#!/usr/bin/env node
// PostToolUse hook target. Reads the hook JSON payload on stdin and logs a skill "access"
// when a hub reference file is Read (cold access) or a tiered spoke is invoked via Skill (hot access).
// MUST be fast and MUST NEVER fail the session — all errors are swallowed and it always exits 0.
import { loadConfig, loadSpokeMap, appendLog, loadState, saveState } from './lib.mjs';

function readStdin() {
  return new Promise(res => {
    let d = ''; let done = false;
    const t = setTimeout(() => { if (!done) { done = true; res(d); } }, 800);
    process.stdin.on('data', c => (d += c));
    process.stdin.on('end', () => { if (!done) { done = true; clearTimeout(t); res(d); } });
    process.stdin.on('error', () => { if (!done) { done = true; clearTimeout(t); res(d); } });
  });
}

try {
  const raw = await readStdin();
  const payload = JSON.parse(raw || '{}');
  const tool = payload.tool_name || payload.toolName;
  const input = payload.tool_input || payload.toolInput || {};
  const cfg = loadConfig();
  const spokeMap = loadSpokeMap(cfg);

  let hit = null; // { spoke, hub, via }
  if (tool === 'Read' && (input.file_path || input.path)) {
    const fp = input.file_path || input.path;
    const m = fp.match(/skills\/([^/]+)\/references\/(.+)\.md$/);
    if (m) {
      const spoke = m[2];
      if (spokeMap.has(spoke)) hit = { spoke, hub: spokeMap.get(spoke).hub, via: 'reference-read' };
    }
  } else if (tool === 'Skill') {
    const name = input.skill || input.command || input.name;
    if (name && spokeMap.has(name)) hit = { spoke: name, hub: spokeMap.get(name).hub, via: 'skill-invoke' };
  }

  if (hit) {
    const ts = new Date().toISOString();
    appendLog({ ts, ...hit });
    const state = loadState();
    const s = state[hit.spoke] || { tier: 'cold', accessCount: 0 };
    s.lastAccess = ts;
    s.accessCount = (s.accessCount || 0) + 1;
    state[hit.spoke] = s;
    saveState(state);
  }
} catch { /* never break the session */ }
process.exit(0);
