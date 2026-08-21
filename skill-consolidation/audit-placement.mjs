#!/usr/bin/env node
// audit-placement.mjs — READ-ONLY whole-tree placement & balance auditor.
//
// The gap this fills: detect-candidates.mjs finds NEW homeless families (by dir-name regex,
// over skills that are NOT already spokes). meta-validate.mjs checks ONE skill's structure.
// Neither re-examines an ALREADY-FOLDED spoke to ask "is it under the RIGHT hub?", nor scores
// the tree's balance against the 1536-char description cap. This script does exactly that —
// the tree-level analytic that `skill-tree-architect` orchestrates. It NEVER mutates anything.
//
// What it reports:
//   1. Hub balance      — per-hub description length vs the 1536 cap (% headroom), spoke count.
//                         Over-cap hubs MUST split or compress (keywords preserved). Very high
//                         spoke counts are split candidates; very low are merge/standalone hints.
//   2. Misplaced spokes — for every folded spoke, score its trigger vocabulary against EVERY hub.
//                         If a SIBLING hub fits better than its current owner by a margin, flag it
//                         as a re-file candidate (executable via `crossroute.mjs`). Heuristic —
//                         human judgment gates the move, exactly like detect-candidates.
//   3. Homeless best-fit — for each unhubbed, non-meta top-level skill, suggest its best-fit hub
//                         (or "(no strong hub)" → a new-family signal; cross-check detect-candidates).
//
// Usage:
//   node audit-placement.mjs                 # human-readable summary (default)
//   node audit-placement.mjs --json          # machine-readable JSON
//   node audit-placement.mjs --margin 0.12   # misplacement score margin (default 0.15)
//   node audit-placement.mjs --desc-cap 1536 # override the per-entry description cap
//
// Exit code is always 0 (advisory tool). Pair with meta-validate.mjs (hard gate) + detect-candidates.mjs.

import fs from 'node:fs';
import path from 'node:path';
import os from 'node:os';
import { fileURLToPath } from 'node:url';
import { loadHubRegistry } from './hub-registry.mjs';
import { EXCLUDE_LIST } from './exclude-list.mjs';

const SKILLS_ROOT = path.join(os.homedir(), '.claude', 'skills');
const CONS_ROOT = path.dirname(fileURLToPath(import.meta.url));

const argv = process.argv.slice(2);
const flags = new Set(argv.filter((a) => a.startsWith('--')));
const asJson = flags.has('--json');
function flagVal(name, dflt) {
  const i = argv.indexOf(name);
  return i >= 0 && argv[i + 1] && !argv[i + 1].startsWith('--') ? argv[i + 1] : dflt;
}
const MARGIN = parseFloat(flagVal('--margin', '0.15')) || 0.15;
const MIN_SCORE = 0.18; // a suggested hub must clear this absolute fit before it is worth surfacing
const DESC_CAP = (() => {
  const v = flagVal('--desc-cap', null);
  if (v) return parseInt(v, 10) || 1536;
  try {
    const s = JSON.parse(fs.readFileSync(path.join(os.homedir(), '.claude', 'settings.json'), 'utf8'));
    return s.skillListingMaxDescChars || s.maxSkillDescriptionChars || 1536;
  } catch { return 1536; }
})();
const HIGH_SPOKES = 14; // >= this many spokes → consider a sub-domain split (cf. HUB-STRATEGY ≥8 to form)

// Meta/infra tools that are intentionally standalone — never flagged as homeless-needing-a-hub.
// Single-sourced in ./exclude-list.mjs (imported above), shared with detect-candidates.mjs.

// A split-parent router (e.g. programming-languages → lang-*) is a thin top-level routing node whose
// children are HUBS, not spokes — so it is not in any manifest, yet it is legitimately top-level and
// must not be flagged as homeless. Identify by its self-describing router language.
function isRouter(desc) {
  return /\bROUTER\b/i.test(desc) || /\bSplit into\b/i.test(desc) || /Route to (?:the )?(?:matching |sub-?hub)/i.test(desc);
}

// Generic words carry no routing signal — drop them so overlap reflects DOMAIN vocabulary, not boilerplate.
const STOP = new Set([
  'the', 'and', 'for', 'with', 'that', 'this', 'your', 'you', 'are', 'not', 'use', 'used', 'using',
  'via', 'when', 'from', 'into', 'per', 'any', 'all', 'its', 'has', 'how', 'why', 'what', 'who',
  'skill', 'skills', 'hub', 'hubs', 'spoke', 'spokes', 'trigger', 'triggers', 'skip', 'reference',
  'references', 'expert', 'guide', 'task', 'tasks', 'user', 'users', 'claude', 'code', 'review',
  'reviewer', 'design', 'build', 'building', 'patterns', 'pattern', 'family', 'sub', 'domain',
  'covers', 'covering', 'including', 'incl', 'plus', 'over', 'across', 'their', 'them', 'one',
  'new', 'first', 'choose', 'choosing', 'matching', 'each', 'data', 'general',
]);

function frontmatterBlock(md) {
  const m = md.match(/^---\n([\s\S]*?)\n---/);
  return m ? m[1] : null;
}
// Extract the full description value (block scalars + folded multi-line), mirroring meta-validate.mjs.
function descriptionText(md) {
  // Prefer the frontmatter block; for banner-led reference copies fall back to the whole file.
  const fm = frontmatterBlock(md);
  const src = fm || md;
  const lines = src.split('\n');
  const i = lines.findIndex((l) => /^description:/.test(l));
  if (i < 0) {
    // No frontmatter description (e.g. a banner-only reference copy) — use the first prose paragraph.
    const prose = md.replace(/^>.*$/gm, '').replace(/^#.*$/gm, '').replace(/^---[\s\S]*?---/, '').trim();
    return prose.slice(0, 600);
  }
  const first = lines[i].replace(/^description:\s*/, '');
  const parts = [];
  if (!/^[|>]/.test(first.trim()) && first.trim()) parts.push(first.trim().replace(/^["']|["']$/g, ''));
  for (let j = i + 1; j < lines.length; j++) {
    if (/^[A-Za-z_][\w-]*:(\s|$)/.test(lines[j])) break; // next top-level key (inline OR bare block-parent like whenToUse:) ends the value
    parts.push(lines[j].trim());
  }
  return parts.join(' ').trim();
}
function tokenize(text) {
  const out = new Set();
  for (const t of (text || '').toLowerCase().match(/[a-z0-9][a-z0-9+.-]{2,}/g) || []) {
    const w = t.replace(/^[-.]+|[-.]+$/g, '');
    if (w.length >= 3 && !STOP.has(w)) out.add(w);
  }
  return out;
}
function readIf(p) { try { return fs.readFileSync(p, 'utf8'); } catch { return null; } }
function pct(n, d) { return d ? Math.round((n / d) * 100) : 0; }

// Resolve a folded spoke's source text: HOT (top-level dir) first, else its reference copy
// (flat references/<spoke>.md, or nested references/<spoke>/SKILL.md).
function spokeText(spoke, hub) {
  const hot = readIf(path.join(SKILLS_ROOT, spoke, 'SKILL.md'));
  if (hot) return { text: hot, where: 'hot' };
  const flat = readIf(path.join(SKILLS_ROOT, hub, 'references', `${spoke}.md`));
  if (flat) return { text: flat, where: 'reference' };
  const nested = readIf(path.join(SKILLS_ROOT, hub, 'references', spoke, 'SKILL.md'));
  if (nested) return { text: nested, where: 'reference-nested' };
  return { text: null, where: 'missing' };
}

function scoreAgainstHubs(tokens, hubTokens) {
  // recall of the candidate's domain tokens within each hub's vocabulary, normalized by candidate size.
  const scores = [];
  for (const [hub, toks] of hubTokens) {
    if (!tokens.size) { scores.push([hub, 0]); continue; }
    let hit = 0;
    for (const w of tokens) if (toks.has(w)) hit++;
    scores.push([hub, hit / tokens.size]);
  }
  scores.sort((a, b) => b[1] - a[1]);
  return scores;
}

function main() {
  const reg = loadHubRegistry({ consRoot: CONS_ROOT });
  const { spokeHub, hubs, spokes } = reg;

  // --- hub vocabulary + cap balance ---
  const hubTokens = new Map();
  const hubBalance = [];
  for (const hub of [...hubs].sort()) {
    const md = readIf(path.join(SKILLS_ROOT, hub, 'SKILL.md'));
    const desc = md ? descriptionText(md) : '';
    hubTokens.set(hub, tokenize(desc));
    const spokeCount = [...spokeHub.values()].filter((h) => h === hub).length;
    hubBalance.push({
      hub,
      present: !!md,
      descLen: desc.length,
      capPct: pct(desc.length, DESC_CAP),
      overCap: desc.length > DESC_CAP,
      spokeCount,
      splitCandidate: spokeCount >= HIGH_SPOKES,
    });
  }
  hubBalance.sort((a, b) => b.descLen - a.descLen);

  // --- misplaced folded spokes ---
  const misplaced = [];
  const spokeIssues = [];
  for (const [spoke, owner] of [...spokeHub.entries()].sort()) {
    const { text, where } = spokeText(spoke, owner);
    if (!text) { spokeIssues.push({ spoke, owner, issue: 'reference-copy-missing' }); continue; }
    const toks = tokenize(descriptionText(text));
    if (toks.size < 4) continue; // too little signal to judge placement
    const scores = scoreAgainstHubs(toks, hubTokens);
    const ownerScore = (scores.find(([h]) => h === owner) || [owner, 0])[1];
    const [bestHub, bestScore] = scores[0];
    if (bestHub !== owner && bestScore >= MIN_SCORE && bestScore - ownerScore >= MARGIN) {
      misplaced.push({
        spoke, owner, where,
        ownerScore: +ownerScore.toFixed(3),
        suggested: bestHub,
        suggestedScore: +bestScore.toFixed(3),
        runnerUp: scores[1] ? { hub: scores[1][0], score: +scores[1][1].toFixed(3) } : null,
      });
    }
  }
  misplaced.sort((a, b) => (b.suggestedScore - b.ownerScore) - (a.suggestedScore - a.ownerScore));

  // --- homeless / standalone best-fit ---
  let live = [];
  try {
    live = fs.readdirSync(SKILLS_ROOT, { withFileTypes: true })
      .filter((d) => d.isDirectory() && !d.name.startsWith('.')).map((d) => d.name).sort();
  } catch { /* skills root unreadable */ }
  const homeless = [];
  const routers = [];
  for (const name of live) {
    if (hubs.has(name) || spokes.has(name) || EXCLUDE_LIST.has(name)) continue;
    const md = readIf(path.join(SKILLS_ROOT, name, 'SKILL.md'));
    if (!md) continue;
    const desc = descriptionText(md);
    if (isRouter(desc)) { routers.push(name); continue; } // legitimately top-level routing node
    const toks = tokenize(desc);
    const scores = scoreAgainstHubs(toks, hubTokens);
    const [bestHub, bestScore] = scores[0] || [null, 0];
    homeless.push({
      skill: name,
      bestHub: bestScore >= MIN_SCORE ? bestHub : null,
      bestScore: +(bestScore || 0).toFixed(3),
      note: bestScore >= MIN_SCORE ? 'fits an existing hub' : '(no strong hub) — new-family signal',
    });
  }
  homeless.sort((a, b) => b.bestScore - a.bestScore);

  const result = {
    generatedAt: new Date().toISOString(),
    skillsRoot: SKILLS_ROOT,
    descCap: DESC_CAP,
    margin: MARGIN,
    counts: {
      hubs: hubs.size,
      foldedSpokes: spokes.size,
      overCapHubs: hubBalance.filter((h) => h.overCap).length,
      misplacedSpokes: misplaced.length,
      spokeIssues: spokeIssues.length,
      homeless: homeless.length,
      routers: routers.length,
    },
    hubBalance, misplaced, spokeIssues, homeless, routers,
  };

  if (asJson) { process.stdout.write(JSON.stringify(result, null, 2) + '\n'); return; }

  const L = [];
  L.push('Skill-tree placement & balance audit (read-only, heuristic)');
  L.push('='.repeat(58));
  L.push(`Generated: ${result.generatedAt}`);
  L.push(`Desc cap: ${DESC_CAP} chars   Misplacement margin: ${MARGIN}`);
  L.push(`Hubs: ${hubs.size}   Folded spokes: ${spokes.size}`);
  L.push('');
  L.push('HUB BALANCE (description length vs cap, spoke count):');
  for (const h of hubBalance) {
    const flags = [];
    if (h.overCap) flags.push('OVER-CAP — split/compress');
    else if (h.capPct >= 90) flags.push('near cap');
    if (h.splitCandidate) flags.push(`${h.spokeCount} spokes — split candidate`);
    if (!h.present) flags.push('SKILL.md missing');
    L.push(`  ${String(h.capPct).padStart(3)}%  ${String(h.descLen).padStart(4)}c  ${String(h.spokeCount).padStart(2)}sp  ${h.hub}${flags.length ? '   <= ' + flags.join('; ') : ''}`);
  }
  L.push('');
  if (spokeIssues.length) {
    L.push(`SPOKE INTEGRITY ISSUES (${spokeIssues.length}) — also caught by meta-validate:`);
    for (const s of spokeIssues) L.push(`  - ${s.spoke} (owner ${s.owner}): ${s.issue}`);
    L.push('');
  }
  L.push(`MISPLACED SPOKE CANDIDATES (${misplaced.length}) — re-file via crossroute.mjs after human review:`);
  if (!misplaced.length) L.push('  (none cleared the margin — placement looks coherent)');
  for (const m of misplaced) {
    L.push(`  - ${m.spoke}: now under ${m.owner} (${m.ownerScore}) → fits ${m.suggested} (${m.suggestedScore})` +
      (m.runnerUp ? `  [runner-up ${m.runnerUp.hub} ${m.runnerUp.score}]` : ''));
  }
  L.push('');
  L.push(`HOMELESS / STANDALONE BEST-FIT (${homeless.length}) — cross-check detect-candidates.mjs for ≥8 families:`);
  for (const h of homeless) {
    L.push(`  - ${h.skill}: ${h.bestHub ? `${h.bestHub} (${h.bestScore})` : h.note}`);
  }
  if (routers.length) L.push(`Routers (top-level routing nodes, OK — not homeless): ${routers.join(', ')}`);
  L.push('');
  L.push('Heuristic — token-overlap scoring is a SIGNAL, not a verdict. Human judgment gates every move.');
  L.push('Next: meta-validate.mjs <hub> (hard gate) · detect-candidates.mjs (new families) · crossroute.mjs (execute a re-file).');
  process.stdout.write(L.join('\n') + '\n');
}

main();
