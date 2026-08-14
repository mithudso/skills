#!/usr/bin/env node
// Skill auto-tiering engine (S3 Intelligent-Tiering analogue).
//   PROMOTE cold -> hot : materialize ~/.claude/skills/<spoke>/SKILL.md from the hub reference
//                         (banner stripped) so it re-enters the always-on skill index.
//   DEMOTE  hot  -> cold: refresh the reference from the (possibly edited) standalone, then remove
//                         the standalone dir so it leaves the index. Content always survives as the reference.
// Flags: --apply (perform changes; default is dry-run) | --status (print the tier table) | --quiet
//        --demote <spoke> (targeted on-demand demote of one HOT spoke — e.g. a skill-tree-architect
//        fold — bypasses the idle/LRU triggers but reuses the same drift-preserving demote path;
//        dry-run unless --apply; exits non-zero if the spoke is unknown, not hot, or blocked)
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { execFileSync } from 'node:child_process';
import { loadConfig, loadSpokeMap, loadState, saveState, readLog, stripBanner, addBanner } from './lib.mjs';

const APPLY = process.argv.includes('--apply');
const QUIET = process.argv.includes('--quiet');
const STATUS = process.argv.includes('--status');
const DEMOTE_IDX = process.argv.indexOf('--demote');
const DEMOTE_TARGET = DEMOTE_IDX !== -1 ? process.argv[DEMOTE_IDX + 1] : null;
if (DEMOTE_IDX !== -1 && (!DEMOTE_TARGET || DEMOTE_TARGET.startsWith('--'))) {
  console.error('[skill-tier] --demote requires a spoke id'); process.exit(1);
}
const log = (...a) => { if (!QUIET) console.log(...a); };

const cfg = loadConfig();
const spokeMap = loadSpokeMap(cfg);
const state = loadState();
const now = Date.now();
const events = readLog();

// recent access count per spoke within the promotion window
const winMs = cfg.windowDays * 86400000;
const recent = {};
const lastAccess = {};
for (const e of events) {
  const t = Date.parse(e.ts);
  if (now - t <= winMs) recent[e.spoke] = (recent[e.spoke] || 0) + 1;
  if (!lastAccess[e.spoke] || t > lastAccess[e.spoke]) lastAccess[e.spoke] = t;
}

function isHot(spoke) { return fs.existsSync(spokeMap.get(spoke).skillMdAbsPath); }
function hotSpokes() { return [...spokeMap.keys()].filter(isHot); }

function hubList(hub) {
  // sibling hubs in the same family (for banner refresh)
  const fam = spokeMap.get([...spokeMap.keys()].find(s => spokeMap.get(s).hub === hub)).family;
  const hubs = new Set([...spokeMap.values()].filter(v => v.family === fam).map(v => v.hub));
  return [...hubs].map(h => `\`${h}\``).join(', ');
}

function promote(spoke) {
  const info = spokeMap.get(spoke);
  if (!fs.existsSync(info.refAbsPath)) return `skip ${spoke} (no reference)`;
  if (isHot(spoke)) return `already-hot ${spoke}`;
  const original = stripBanner(fs.readFileSync(info.refAbsPath, 'utf8'));
  if (APPLY) {
    fs.mkdirSync(info.skillDirAbsPath, { recursive: true });
    fs.writeFileSync(info.skillMdAbsPath, original);
    state[spoke] = { ...(state[spoke] || {}), tier: 'hot', promotedAt: new Date().toISOString(), lastAccess: state[spoke]?.lastAccess || new Date().toISOString() };
  }
  return `PROMOTE ${spoke} -> hot (index)`;
}

function demote(spoke, reason) {
  const info = spokeMap.get(spoke);
  if (!isHot(spoke)) return `not-hot ${spoke}`;
  if (!fs.existsSync(info.refAbsPath)) return `BLOCK demote ${spoke} (reference missing — refusing to remove)`;
  if (APPLY) {
    // preserve any drift: refresh the reference from the current standalone, then remove the standalone
    const cur = fs.readFileSync(info.skillMdAbsPath, 'utf8');
    const refreshed = addBanner(stripBanner(cur), info.hub, spoke, hubList(info.hub));
    fs.writeFileSync(info.refAbsPath, refreshed);
    fs.rmSync(info.skillDirAbsPath, { recursive: true, force: true });
    state[spoke] = { ...(state[spoke] || {}), tier: 'cold', demotedAt: new Date().toISOString() };
  }
  return `DEMOTE ${spoke} -> cold (${reason})`;
}

if (STATUS) {
  const rows = [...spokeMap.keys()].map(s => ({
    spoke: s, hub: spokeMap.get(s).hub, tier: isHot(s) ? 'HOT' : 'cold',
    recent: recent[s] || 0, last: lastAccess[s] ? new Date(lastAccess[s]).toISOString().slice(0, 10) : '-'
  })).filter(r => r.tier === 'HOT' || r.recent > 0);
  console.log('Tiered spokes with activity (HOT or recently accessed):');
  if (!rows.length) console.log('  (none yet — all cold; access a reference to start accumulating signal)');
  for (const r of rows) console.log(`  [${r.tier}] ${r.spoke}  hub=${r.hub}  recent=${r.recent}  last=${r.last}`);
  console.log(`\nTotal tiered spokes: ${spokeMap.size} | currently HOT (in index): ${hotSpokes().length} | maxHot: ${cfg.maxHot}`);
  // PAIR-DRIFT detection: while a spoke is HOT, its standalone SKILL.md and the hub reference
  // (banner stripped) form a pair that only demote() re-syncs — surface any divergence here.
  for (const s of hotSpokes()) {
    const info = spokeMap.get(s);
    if (!fs.existsSync(info.refAbsPath)) continue;
    const ref = stripBanner(fs.readFileSync(info.refAbsPath, 'utf8')).split('\n');
    const cur = fs.readFileSync(info.skillMdAbsPath, 'utf8').split('\n');
    let n = Math.abs(ref.length - cur.length);
    for (let i = 0; i < Math.min(ref.length, cur.length); i++) if (ref[i] !== cur[i]) n++;
    if (n) console.log(`  PAIR-DRIFT ${s} (${n} differing lines)`);
  }
  process.exit(0);
}

const actions = [];
if (DEMOTE_TARGET) {
  // Targeted on-demand demote (--demote <spoke>): the deterministic path for folding a HOT spoke
  // (skill-tree-architect Phase 3 step 4). Bypasses the idle/LRU triggers but reuses the same
  // drift-preserving demote() (standalone -> reference sync, then removal).
  if (!spokeMap.has(DEMOTE_TARGET)) { console.error(`[skill-tier] --demote: unknown spoke '${DEMOTE_TARGET}'`); process.exit(1); }
  actions.push(demote(DEMOTE_TARGET, 'manual fold'));
} else {
  // 1) PROMOTE cold spokes over threshold
  for (const spoke of spokeMap.keys()) {
    if (!isHot(spoke) && (recent[spoke] || 0) >= cfg.promoteThreshold) actions.push(promote(spoke));
  }
  // 2) DEMOTE idle hot spokes
  for (const spoke of hotSpokes()) {
    const la = lastAccess[spoke] || Date.parse(state[spoke]?.promotedAt || 0);
    if (la && (now - la) > cfg.demoteAfterDays * 86400000) actions.push(demote(spoke, `idle >${cfg.demoteAfterDays}d`));
  }
  // 3) Enforce maxHot via LRU eviction
  const hot = hotSpokes();
  if (hot.length > cfg.maxHot) {
    const byLru = hot.sort((a, b) => (lastAccess[a] || 0) - (lastAccess[b] || 0));
    for (const spoke of byLru.slice(0, hot.length - cfg.maxHot)) actions.push(demote(spoke, 'maxHot LRU eviction'));
  }
}

if (APPLY) saveState(state);
log(`[skill-tier]${APPLY ? '' : ' DRY-RUN'} ${actions.length} action(s); HOT now: ${hotSpokes().length}/${cfg.maxHot}`);
for (const a of actions) log('  ' + a);
if (!actions.length) log('  (no tier changes)');

// After any applied promote/demote, keep cross-skill routing referents resolvable: a promoted spoke
// becomes HOT (restore its bare id) and a demoted spoke becomes COLD (repoint to its owning hub).
// referents.mjs is idempotent, so running it on every --apply (even with no tier changes) is safe.
if (APPLY) {
  const REFERENTS = path.join(path.dirname(fileURLToPath(import.meta.url)), '..', 'referents.mjs');
  try {
    const out = execFileSync(process.execPath, [REFERENTS, '--repair', '--apply', '--quiet'], { encoding: 'utf8' });
    if (out.trim()) log(out.trim());
    else log('[skill-tier] referent repair: ok');
  } catch (e) {
    console.error('[skill-tier] referent repair FAILED:', e.message);
  }
}

// Targeted mode is scriptable: exit non-zero when the requested demote did not happen
// (spoke not hot, or blocked because its hub reference is missing).
if (DEMOTE_TARGET && !actions[0]?.startsWith('DEMOTE')) process.exit(1);
