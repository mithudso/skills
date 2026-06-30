#!/usr/bin/env node
// Referent-integrity updater for the skill hub/tiering system.
//
// The skill-optimizer seeds two kinds of cross-skill routing referents:
//   1. Frontmatter `related_skills:` — a YAML list of skill ids (block or inline `[a, b]` form).
//   2. Inline deferral tokens in description/body of the form `→ <skill-id>` or `→ `<skill-id>``
//      (often preceded by `SKIP:`).
// When a spoke is folded into a hub (becomes a cold references/<spoke>.md, top-level dir removed),
// or promoted/demoted by the tiering engine, these referents may point at a skill id that is no
// longer a top-level (indexed) skill -> routing breaks. This tool keeps referents resolvable.
//
// Model:
//   - spoke -> owningHub map is built from every *-manifest.json in skill-consolidation.
//   - A spoke is COLD if ~/.claude/skills/<spoke>/SKILL.md does NOT exist (folded), HOT if it does.
//   - Hubs and non-tiered skills are always resolvable (never rewritten as a target).
//
// Repairs (only on top-level ~/.claude/skills/*/SKILL.md; never reference files under references/):
//   - COLD spoke referent -> rewrite to owning hub:
//       related_skills entry  : <spoke>  ->  <hub>   (dedup; drop self-refs; drop exact dups)
//       inline arrow          : → <spoke>  ->  → <hub> (references/<spoke>.md)
//   - HOT spoke referent that currently points at the hub (reverse direction, idempotent):
//       restore the bare <spoke> id using the ledger so repeated repairs are stable.
//   - Every rewrite recorded in referents-ledger.json: {file, surface, from, to, spoke, hub}.
//
// Flags: --repair (do the scan; dry-run unless --apply) | --apply (write) | --quiet
//        --status (count dangling referents without writing)
// Idempotent and safe: only referent tokens / list entries are touched; never reference files,
// never hub routing tables, never non-referent content.

import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { loadSpokeMap as loadSpokeMapShared } from './hub-registry.mjs';

const HERE = path.dirname(fileURLToPath(import.meta.url));
const CONFIG_PATH = path.join(HERE, 'tiering', 'tier-config.json');
const LEDGER_PATH = path.join(HERE, 'referents-ledger.json');

const APPLY = process.argv.includes('--apply');
const QUIET = process.argv.includes('--quiet');
const STATUS = process.argv.includes('--status');
const REPAIR = process.argv.includes('--repair');
const log = (...a) => { if (!QUIET) console.log(...a); };

// ---- config + spoke map ----------------------------------------------------
const cfg = JSON.parse(fs.readFileSync(CONFIG_PATH, 'utf8'));
// Derive roots from this file's location (portable, survives `git pull`); the tier-config
// values are advisory only. Layout: <skillsRoot>/skill-consolidation/referents.mjs
const CONS_ROOT = HERE;
const SKILLS_ROOT = path.dirname(HERE);

// spoke -> hub across every manifest in the consolidation dir (a "tiered spoke" is any spoke
// in any manifest). We scan ALL *-manifest.json, not just tier-config's subset, so a referent
// to any folded spoke is repairable.
function loadSpokeMap() {
  // Delegates to the shared hub-registry loader (audit rec #3) so spoke->hub resolution
  // lives in one place. consolidation-manifest.json is excluded there as a known duplicate.
  return loadSpokeMapShared({ consRoot: CONS_ROOT });
}
const spokeHub = loadSpokeMap();
const isTieredSpoke = (id) => spokeHub.has(id);
const isCold = (spoke) => !fs.existsSync(path.join(SKILLS_ROOT, spoke, 'SKILL.md'));

// ---- ledger ----------------------------------------------------------------
function loadLedger() {
  try { return JSON.parse(fs.readFileSync(LEDGER_PATH, 'utf8')); } catch { return { rewrites: [] }; }
}
function saveLedger(ledger) {
  fs.writeFileSync(LEDGER_PATH, JSON.stringify(ledger, null, 2) + '\n');
}

// ---- top-level skill discovery (hubs + standalones; NOT reference files) ---
function topLevelSkillMds() {
  const out = [];
  for (const d of fs.readdirSync(SKILLS_ROOT)) {
    const p = path.join(SKILLS_ROOT, d, 'SKILL.md');
    if (fs.existsSync(p)) out.push({ dir: d, path: p });
  }
  return out;
}

// ---- frontmatter split -----------------------------------------------------
// Returns { fm, body, hasFm } where fm is the YAML text (without --- fences).
function splitFrontmatter(txt) {
  const m = txt.match(/^(---\r?\n)([\s\S]*?\r?\n)(---\r?\n)/);
  if (!m) return { pre: '', fm: '', post: txt, hasFm: false };
  const pre = m[1];
  const fm = m[2];
  const fence = m[3];
  const post = txt.slice(m[0].length);
  return { pre, fm, fence, post, hasFm: true };
}

// ---- related_skills rewriting ----------------------------------------------
// Handles both block form (`related_skills:` then `  - id` lines) and inline array
// (`related_skills: [a, b, c]`). Returns { fm, changes:[{from,to,spoke,hub,reverse}] }.
function rewriteRelatedSkills(fm, selfId, ledger) {
  const changes = [];

  // Build reverse-lookup from ledger for this file's related_skills surface:
  // hub -> set of spokes that were rewritten to that hub (so we can restore HOT spokes).
  const ledgerForFile = (surface) => ledger.rewrites.filter(r => r.skillId === selfId && r.surface === surface);

  // --- inline array form ---
  // Match only the bracketed list on its own line; do NOT consume the trailing newline
  // (a greedy \s*$ in multiline mode would eat the \n separating this line from the fence).
  const inlineRe = /^(related_skills:[ \t]*)\[([^\]]*)\][ \t]*(?=\r?\n|$)/m;
  const inlineMatch = fm.match(inlineRe);
  if (inlineMatch) {
    const items = inlineMatch[2].split(',').map(s => s.trim()).filter(Boolean);
    const { next, ch } = transformList(items, selfId, ledgerForFile('related_skills'));
    for (const c of ch) changes.push(c);
    if (ch.length) {
      const rebuilt = `${inlineMatch[1]}[${next.join(', ')}]`;
      fm = fm.replace(inlineRe, rebuilt);
    }
    return { fm, changes };
  }

  // --- block form ---
  const lines = fm.split('\n');
  const startIdx = lines.findIndex(l => /^related_skills:\s*$/.test(l));
  if (startIdx === -1) return { fm, changes };
  // collect contiguous `  - id` entries
  let i = startIdx + 1;
  const entryIdx = [];
  const items = [];
  while (i < lines.length && /^\s*-\s+/.test(lines[i])) {
    const id = lines[i].replace(/^\s*-\s+/, '').trim();
    entryIdx.push(i); items.push(id); i++;
  }
  const indent = entryIdx.length ? (lines[entryIdx[0]].match(/^(\s*)-/)[1]) : '  ';
  const { next, ch } = transformList(items, selfId, ledgerForFile('related_skills'));
  for (const c of ch) changes.push(c);
  if (ch.length) {
    const rebuilt = next.map(id => `${indent}- ${id}`);
    lines.splice(startIdx + 1, entryIdx.length, ...rebuilt);
    fm = lines.join('\n');
  }
  return { fm, changes };
}

// Core list transform shared by inline + block forms.
// Forward: cold spoke -> hub (dedup, drop self, drop exact dup). Reverse: hub -> hot spoke (from ledger).
function transformList(items, selfId, ledgerEntries) {
  const changes = [];
  const out = [];
  const seen = new Set();
  // reverse restore: if a hub entry has a ledgered spoke that is now HOT, restore that spoke id.
  const hubToHotSpokes = new Map();
  for (const e of ledgerEntries) {
    if (e.reverse) continue;
    if (isTieredSpoke(e.spoke) && !isCold(e.spoke)) {
      if (!hubToHotSpokes.has(e.to)) hubToHotSpokes.set(e.to, []);
      hubToHotSpokes.get(e.to).push(e.spoke);
    }
  }
  for (let id of items) {
    let from = id, spoke = null, hub = null, reverse = false;
    if (isTieredSpoke(id) && isCold(id)) {
      // forward: cold spoke -> hub
      spoke = id; hub = spokeHub.get(id); id = hub;
      changes.push({ surface: 'related_skills', from, to: id, spoke, hub, reverse: false });
    } else if (hubToHotSpokes.has(id) && hubToHotSpokes.get(id).length) {
      // reverse: hub entry -> restore a now-HOT spoke (one per ledger entry, FIFO)
      spoke = hubToHotSpokes.get(id).shift();
      hub = id; from = id; id = spoke; reverse = true;
      changes.push({ surface: 'related_skills', from, to: id, spoke, hub, reverse: true });
    }
    // dedup + drop self-reference
    if (id === selfId) continue;
    if (seen.has(id)) continue;
    seen.add(id);
    out.push(id);
  }
  return { next: out, ch: changes };
}

// ---- inline arrow rewriting ------------------------------------------------
// Forward:  → <cold-spoke>            -> → <hub> (references/<spoke>.md)
//           → `<cold-spoke>`          -> → `<hub>` (references/<spoke>.md)
// Reverse:  → <hub> (references/<spoke>.md)  -> → <spoke>   when spoke is now HOT.
function rewriteInlineArrows(body, selfId, ledger) {
  const changes = [];

  // Reverse first: collapse "<hub> (references/<spoke>.md)" back to bare spoke when HOT.
  // Match optional backticks around hub id.
  body = body.replace(/→(\s*)`?([a-z0-9][a-z0-9-]*)`?\s*\(references\/([a-z0-9][a-z0-9-]*)\.md\)/g,
    (full, sp1, hub, spoke) => {
      if (isTieredSpoke(spoke) && !isCold(spoke)) {
        changes.push({ surface: 'inline', from: full.trim(), to: `→${sp1}${spoke}`.trim(), spoke, hub, reverse: true });
        return `→${sp1}${spoke}`;
      }
      return full; // spoke still cold -> leave the hub pointer in place
    });

  // Forward: bare arrow target that is a cold tiered spoke.
  // Negative lookahead avoids re-matching the "(references/...)" form we may have just produced.
  body = body.replace(/→(\s*)(`?)([a-z0-9][a-z0-9-]*)(`?)(?!\s*\(references\/)/g,
    (full, sp1, bt1, id, bt2) => {
      if (isTieredSpoke(id) && isCold(id)) {
        const hub = spokeHub.get(id);
        const rep = `→${sp1}${bt1}${hub}${bt2} (references/${id}.md)`;
        changes.push({ surface: 'inline', from: full.trim(), to: rep.trim(), spoke: id, hub, reverse: false });
        return rep;
      }
      return full;
    });

  return { body, changes };
}

// ---- driver ----------------------------------------------------------------
function run() {
  if (!REPAIR && !STATUS) {
    console.error('usage: node referents.mjs --repair [--apply] [--quiet] | --status');
    process.exit(2);
  }
  const ledger = loadLedger();
  const files = topLevelSkillMds();
  let totalRewrites = 0, totalFiles = 0, dangling = 0;
  const newLedgerEntries = [];

  for (const f of files) {
    const selfId = f.dir;
    const orig = fs.readFileSync(f.path, 'utf8');
    const { pre, fm, fence, post, hasFm } = splitFrontmatter(orig);

    let fmOut = fm, bodyOut = hasFm ? post : orig;
    const fileChanges = [];

    if (STATUS) {
      // count dangling referents (cold-spoke targets in related_skills or inline arrows)
      const fmText = hasFm ? fm : '';
      const bodyText = hasFm ? post : orig;
      // related_skills (block + inline)
      const rsInline = fmText.match(/^related_skills:\s*\[([^\]]*)\]/m);
      const rsIds = [];
      if (rsInline) {
        rsInline[1].split(',').forEach(s => { const v = s.trim(); if (v) rsIds.push(v); });
      } else {
        const lines = fmText.split('\n');
        const start = lines.findIndex(l => /^related_skills:\s*$/.test(l));
        if (start !== -1) {
          for (let j = start + 1; j < lines.length && /^\s*-\s+/.test(lines[j]); j++) {
            rsIds.push(lines[j].replace(/^\s*-\s+/, '').trim());
          }
        }
      }
      for (const id of rsIds) if (isTieredSpoke(id) && isCold(id)) dangling++;
      for (const m of bodyText.matchAll(/→\s*`?([a-z0-9][a-z0-9-]*)`?(?!\s*\(references\/)/g)) {
        if (isTieredSpoke(m[1]) && isCold(m[1])) dangling++;
      }
      continue;
    }

    if (hasFm) {
      const rs = rewriteRelatedSkills(fmOut, selfId, ledger);
      fmOut = rs.fm;
      for (const c of rs.changes) fileChanges.push({ ...c, skillId: selfId });
    }
    const ar = rewriteInlineArrows(bodyOut, selfId, ledger);
    bodyOut = ar.body;
    for (const c of ar.changes) fileChanges.push({ ...c, skillId: selfId });

    if (!fileChanges.length) continue;

    const rebuilt = hasFm ? (pre + fmOut + fence + bodyOut) : bodyOut;
    totalFiles++;
    totalRewrites += fileChanges.length;
    for (const c of fileChanges) {
      newLedgerEntries.push({ file: f.path, skillId: selfId, surface: c.surface, from: c.from, to: c.to, spoke: c.spoke, hub: c.hub, reverse: !!c.reverse });
      log(`  [${c.reverse ? 'restore' : 'repair '}] ${selfId} (${c.surface}): ${c.from}  ->  ${c.to}`);
    }
    if (APPLY) fs.writeFileSync(f.path, rebuilt);
  }

  if (STATUS) {
    console.log(`[referents] dangling referents (cold spokes still named directly): ${dangling}`);
    process.exit(0);
  }

  if (APPLY && newLedgerEntries.length) {
    ledger.rewrites.push(...newLedgerEntries.map(e => ({ ...e, at: new Date().toISOString() })));
    saveLedger(ledger);
  }
  log(`[referents]${APPLY ? '' : ' DRY-RUN'} ${totalRewrites} referent rewrite(s) across ${totalFiles} file(s).`);
}

run();
