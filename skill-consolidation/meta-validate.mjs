#!/usr/bin/env node
// meta-validate.mjs — deterministic structural linter for the skill ecosystem.
//
// This is the "fill the gaps" half of skill-optimizer's `--meta` (structural-only) mode:
// it covers the meta-work that NO existing script owns. It does NOT re-derive spoke->hub
// state (it imports hub-registry.mjs) and it does NOT touch content — prose quality, trigger
// rewriting, and peer seeding stay with skill-optimizer; deep frontmatter audit stays with
// Pass G; tool-search discoverability stays with mcp-tool-search-optimizer.
//
// Checks (all deterministic):
//   1. file/folder + naming    — SKILL.md present, frontmatter `name` == dir basename, kebab-case,
//                                 hub has a references/ dir.
//   2. manifest schema         — owning *-manifest.json parses and matches the expected shape.
//   3. spoke-copy-exists       — every spoke in the manifest has its declared referenceFile copy —
//                                 flat references/<spoke>.md OR nested references/<spoke>/SKILL.md
//                                 (a missing copy is HIGH: deleting that spoke dir would lose it).
//   4. dangling routing rows   — every routing-table row points to an existing reference file
//                                 (local rows resolve under this hub; cross-hub-map rows resolve
//                                 under the sibling hub named in column 1); every local reference
//                                 appears in the routing table (orphan check).
//   5. circular SKIP           — heuristic mutual-hard-SKIP detection (A -> B and B -> A).
//   6. tier-config presence    — a hub's family manifest is registered in tiering/tier-config.json
//                                 (`--register-tier --apply` inserts it idempotently).
//   7. description cap          — frontmatter `description` length <= per-entry cap (default 1536,
//                                 or settings.json skillListingMaxDescChars / --desc-cap). Over-cap
//                                 descriptions truncate in the index; fix by splitting a hub or
//                                 compressing prose WITHOUT dropping spoke keywords.
//
// Usage:
//   node meta-validate.mjs <skill-id-or-path> [--json] [--register-tier] [--apply]
//        [--cons-root <dir>] [--skills-root <dir>]
//
// Exit code: 1 if any High finding (CI/skill gate), else 0. `--json` prints a machine-readable
// report on stdout regardless. Only `--register-tier --apply` writes anything (tier-config.json).

import fs from 'node:fs';
import path from 'node:path';
import os from 'node:os';
import { fileURLToPath } from 'node:url';
import { loadHubRegistry } from './hub-registry.mjs';

// Optional real YAML parser: used for a full frontmatter-validity check when resolvable, but never
// required (this dir has no node_modules). Deep frontmatter audit is skill-optimizer Pass G's job.
let yaml = null;
try { yaml = (await import('js-yaml')).default; } catch { /* not installed here — light checks only */ }

// ---- args -----------------------------------------------------------------
const argv = process.argv.slice(2);
const flags = new Set(argv.filter((a) => a.startsWith('--')));
const positional = argv.filter((a) => !a.startsWith('--'));
function flagVal(name) {
  const i = argv.indexOf(name);
  return i >= 0 && argv[i + 1] && !argv[i + 1].startsWith('--') ? argv[i + 1] : null;
}
const asJson = flags.has('--json');
const doRegisterTier = flags.has('--register-tier');
const doApply = flags.has('--apply');

// Default to this script's own dir (portable, survives `git pull`); --cons-root still overrides.
const CONS_ROOT = flagVal('--cons-root') || path.dirname(new URL(import.meta.url).pathname);

// skillsRoot is the parent of the skill-consolidation dir under the fixed repo layout
// (<skillsRoot>/skill-consolidation); --skills-root still overrides for non-standard checkouts.
const SKILLS_ROOT = flagVal('--skills-root') || path.dirname(CONS_ROOT);

const target = positional[0];
if (!target) {
  console.error('usage: node meta-validate.mjs <skill-id-or-path> [--json] [--register-tier] [--apply]');
  process.exit(2);
}
// One target per invocation. A space-joined list (e.g. a mangled shell loop) used to fall through
// as a single bogus skill id and surface as an opaque error — fail loudly instead.
if (positional.length > 1) {
  console.error(`meta-validate: expected exactly one target, got ${positional.length}: "${positional.join(' ')}". Run once per skill (loop in the shell).`);
  process.exit(2);
}

// ---- helpers (regex frontmatter reader, mirrors build.mjs conventions) -----
function frontmatterBlock(md) {
  const mm = md.match(/^---\n([\s\S]*?)\n---/);
  return mm ? mm[1] : null;
}
function field(fm, name) {
  const re = new RegExp(`^${name}:\\s*(?:[|>][-+]?)?\\s*(.*)$`, 'm');
  const mm = fm.match(re);
  if (!mm) return '';
  const val = mm[1].trim();
  return val.replace(/^["']|["']$/g, '');
}
const KEBAB = /^[a-z0-9]+(?:-[a-z0-9]+)*$/;

// Full description length incl. block scalars (`>`/`|`) and folded multi-line values — `field()`
// only grabs the first line, which under-counts. Used by the cap-compliance check.
function descriptionLen(fmText) {
  const lines = fmText.split('\n');
  const i = lines.findIndex((l) => /^description:/.test(l));
  if (i < 0) return 0;
  const first = lines[i].replace(/^description:\s*/, '');
  const parts = [];
  if (!/^[|>]/.test(first.trim()) && first.trim()) parts.push(first.trim().replace(/^["']|["']$/g, ''));
  for (let j = i + 1; j < lines.length; j++) {
    if (/^[A-Za-z_][\w-]*:(\s|$)/.test(lines[j])) break; // next top-level key (inline value OR bare block-parent like whenToUse:) ends the value
    parts.push(lines[j].trim());
  }
  return parts.join(' ').trim().length;
}
// Per-entry description cap. Honors settings.json overrides; defaults to the harness default 1536.
const DESC_CAP = (() => {
  if (flagVal('--desc-cap')) return parseInt(flagVal('--desc-cap'), 10) || 1536;
  try {
    const s = JSON.parse(fs.readFileSync(path.join(os.homedir(), '.claude', 'settings.json'), 'utf8'));
    return s.skillListingMaxDescChars || s.maxSkillDescriptionChars || 1536;
  } catch { return 1536; }
})();

const findings = [];
function add(check, level, message, fix) {
  findings.push({ check, level, message, ...(fix ? { fix } : {}) });
}

// ---- resolve the target to a SKILL.md path + identity ----------------------
const registry = loadHubRegistry({ consRoot: CONS_ROOT });

let skillPath = null;
let skillId = null;
let dirName = null; // basename of the skill directory (for top-level skills)
let isFoldedSpoke = false;

if (target.includes('/') || target.endsWith('.md')) {
  skillPath = path.resolve(target);
  dirName = path.basename(path.dirname(skillPath));
  skillId = dirName;
} else {
  skillId = target;
  const topLevel = path.join(SKILLS_ROOT, skillId, 'SKILL.md');
  if (fs.existsSync(topLevel)) {
    skillPath = topLevel;
    dirName = skillId;
  } else {
    // maybe a folded spoke: ~/.claude/skills/<hub>/references/<id>.md
    const hub = registry.spokeHub.get(skillId);
    if (hub) {
      const ref = path.join(SKILLS_ROOT, hub, 'references', `${skillId}.md`);
      if (fs.existsSync(ref)) { skillPath = ref; isFoldedSpoke = true; }
    }
  }
}

// isHub is referenced by report(); compute it before the not-found guard so the early-exit path doesn't hit a TDZ.
const isHub = registry.hubs.has(skillId);

if (!skillPath || !fs.existsSync(skillPath)) {
  add('file/folder', 'High', `Cannot locate SKILL.md for "${target}" (looked under ${SKILLS_ROOT}).`);
  report(); // report() calls process.exit(), so execution stops here
}

const md = fs.readFileSync(skillPath, 'utf8');
const fm = frontmatterBlock(md);

// ---- check 1: file/folder + naming ----------------------------------------
if (fm === null) {
  add('frontmatter', 'High', `${skillPath} has no parseable \`---\` frontmatter block.`);
} else {
  if (/\t/.test(fm)) add('frontmatter', 'High', 'Literal tab in YAML frontmatter (hard YAML error).');
  if (yaml) {
    try { yaml.load(fm); } catch (e) { add('frontmatter', 'High', `Frontmatter fails YAML parse: ${e.message}`); }
  }
  const nameField = field(fm, 'name');
  if (!nameField) {
    add('naming', 'Medium', 'Frontmatter is missing a `name:` field.');
  } else {
    if (!KEBAB.test(nameField)) add('naming', 'Medium', `\`name: ${nameField}\` is not kebab-case.`);
    if (!isFoldedSpoke && dirName && nameField !== dirName) {
      add('naming', 'High', `Frontmatter \`name: ${nameField}\` != directory "${dirName}" — the id won't resolve.`);
    }
  }
}

// hub must own a references/ dir
let refDir = null;
if (isHub && !isFoldedSpoke) {
  refDir = path.join(path.dirname(skillPath), 'references');
  if (!fs.existsSync(refDir) || !fs.statSync(refDir).isDirectory()) {
    add('file/folder', 'High', `Hub "${skillId}" has no references/ directory at ${refDir}.`);
    refDir = null;
  }
}

// ---- locate the owning family manifest (hub only) --------------------------
let ownerManifestFile = null;
let ownerFamily = null;
let ownerHubDef = null;
if (isHub) {
  for (const f of registry.manifestFiles) {
    try {
      const m = JSON.parse(fs.readFileSync(path.join(CONS_ROOT, f), 'utf8'));
      if (m.hubs && m.hubs[skillId]) { ownerManifestFile = f; ownerFamily = m.family; ownerHubDef = m.hubs[skillId]; break; }
    } catch { /* skipped by registry already */ }
  }
}

// ---- check 2: manifest schema (hub only) -----------------------------------
if (isHub && ownerHubDef) {
  if (!ownerFamily) add('manifest-schema', 'Medium', `${ownerManifestFile} has no top-level "family" string.`);
  if (!Array.isArray(ownerHubDef.spokes)) {
    add('manifest-schema', 'High', `${ownerManifestFile} hub "${skillId}" has no spokes[] array.`);
  } else {
    ownerHubDef.spokes.forEach((s, i) => {
      const spoke = typeof s === 'string' ? s : s && s.spoke;
      if (!spoke) { add('manifest-schema', 'Medium', `${ownerManifestFile} spoke[${i}] has no resolvable id.`); return; }
      // Two conventional forms are valid: flat `references/<spoke>.md` and nested
      // `references/<spoke>/SKILL.md` (a folded standalone skill kept as a dir). Only flag neither.
      if (typeof s === 'object' && s.referenceFile &&
          s.referenceFile !== `references/${spoke}.md` &&
          s.referenceFile !== `references/${spoke}/SKILL.md`) {
        add('manifest-schema', 'Low', `spoke "${spoke}" referenceFile "${s.referenceFile}" is neither references/${spoke}.md nor references/${spoke}/SKILL.md.`);
      }
    });
  }
}

// ---- check 3: spoke-copy-exists-before-delete ------------------------------
if (isHub && refDir && Array.isArray(ownerHubDef?.spokes)) {
  for (const s of ownerHubDef.spokes) {
    const spoke = typeof s === 'string' ? s : s && s.spoke;
    if (!spoke) continue;
    // Honor the manifest's declared referenceFile (flat or nested references/<spoke>/SKILL.md),
    // matching what the tiering engine (lib.mjs) resolves; default to the flat form.
    const refRel = (typeof s === 'object' && s.referenceFile) ? s.referenceFile : `references/${spoke}.md`;
    const copy = path.join(path.dirname(skillPath), refRel);
    if (!fs.existsSync(copy)) {
      add('spoke-copy', 'High',
        `Spoke "${spoke}" is in the manifest but ${refRel} is MISSING — do NOT delete its standalone dir; rebuild the copy first (build.mjs).`);
    }
    const standalone = path.join(SKILLS_ROOT, spoke, 'SKILL.md');
    if (fs.existsSync(standalone) && fs.existsSync(copy)) {
      add('spoke-copy', 'Low', `Spoke "${spoke}" still has a standalone dir AND a reference copy — safe to fold (delete the standalone dir).`);
    }
  }
}

// ---- check 4: dangling routing rows + orphan references --------------------
if (isHub && refDir) {
  // The regex matches a single backtick-wrapped `references/X.md` cell. It catches two 3-col tables
  // that share this shape, distinguished by whether col 1 is a known hub other than the target:
  //   - local Sub-skill routing table — col 1 is a SPOKE; col 3 is the reference under THIS hub
  //     (load-bearing — it names the file to Read), so it resolves against the target's own dir.
  //   - Cross-hub map (single-file rows) — col 1 is a SIBLING HUB; col 3 is an example reference
  //     that lives under THAT hub, so it resolves under SKILLS_ROOT/<col1>/, not locally.
  // The target's own cross-hub-map row (col1 === skillId) falls to the local branch (same dir).
  // NOTE: we deliberately do NOT validate MULTI-file cross-hub-map cells. That "example reference
  // files" column is free-form illustrative prose with inconsistent formats across hubs
  // (`references/X.md`, `<hub>/references/X.md`, `…/Y.md`, bare `Y.md`); treating it as a validatable
  // list produced false positives (path mangling on hub-prefixed/ellipsis forms) and spurious
  // orphans (a hub's own cross-hub self-row flips the orphan check on). It is an "e.g." column,
  // not a contract — left unvalidated by design.
  const rowRe = /^\|\s*`([^`]+)`\s*\|[^|]*\|\s*`(references\/[^`]+)`\s*\|/gm;
  const rowRefs = new Set();
  let mm;
  while ((mm = rowRe.exec(md)) !== null) {
    const col1 = mm[1];
    const refRel = mm[2];
    const isCrossHubRow = registry.hubs.has(col1) && col1 !== skillId;
    const baseDir = isCrossHubRow ? path.join(SKILLS_ROOT, col1) : path.dirname(skillPath);
    // Only local rows seed the orphan set — a sibling hub's example file is not a local reference
    // and must not mask a genuinely orphaned local reference that happens to share its basename.
    if (!isCrossHubRow) rowRefs.add(path.basename(refRel));
    if (!fs.existsSync(path.join(baseDir, refRel))) {
      const where = isCrossHubRow ? `under sibling hub "${col1}"` : 'locally';
      add('dangling-row', 'High', `Routing-table row for "${col1}" points to ${refRel} (${where}), which does not exist.`);
    }
  }
  // orphan references: a *.md in references/ not present in any routing row
  let refFiles = [];
  try { refFiles = fs.readdirSync(refDir).filter((f) => f.endsWith('.md')); } catch { /* refDir already flagged */ }
  for (const f of refFiles) {
    // The hub's own context file (`<hub>-context.md` / `*-context.md`) is provenance, not a spoke —
    // it never has a routing row by design, so it is not an orphan.
    if (/-context\.md$/.test(f) || f === `${skillId}-context.md`) continue;
    if (rowRefs.size > 0 && !rowRefs.has(f)) {
      add('dangling-row', 'Medium', `references/${f} exists but has no row in the hub routing table (orphan reference).`);
    }
  }
}

// ---- check 5: circular / mutual-hard-SKIP (same-topic only) ----------------
// The defect (per skill-optimizer Pass O) is A and B deferring to each other for the SAME topic —
// an unresolvable loop. Healthy sibling hubs defer to each other for DIFFERENT topics, which is the
// routing mesh working; flagging those would be noise. So extract each `-> <id>` edge WITH the
// topic phrase preceding the arrow, and flag a mutual edge only when the two directions share a
// topic token.
const STOP = new Set(['the', 'and', 'for', 'use', 'using', 'via', 'with', 'that', 'this', 'design', 'general', 'when', 'not', 'see', 'skip']);
function skipEdges(text) {
  const skipClause = text.split(/\bSKIP\b/i)[1] || '';
  const edges = [];
  for (const seg of skipClause.split(/[;\n•]|(?:^|\s)-\s/)) {
    const arrowRe = /(?:->|→)\s*`?([a-z0-9]+(?:-[a-z0-9]+)*)`?/g;
    let mm;
    while ((mm = arrowRe.exec(seg)) !== null) {
      const before = seg.slice(0, mm.index).toLowerCase();
      const toks = new Set(before.split(/[^a-z0-9]+/).filter((t) => t.length > 2 && !STOP.has(t)));
      edges.push({ id: mm[1], toks });
    }
  }
  return edges;
}
if (fm) {
  // Scan ONLY the frontmatter SKIP clause — the declared routing contract — not the whole body.
  // Body prose (e.g. "High-overlap routing notes" that legitimately names BOTH directions of a
  // broad<->specific boundary) is not a deferral edge; including it produced false "mutual SKIP"
  // loops on healthy bidirectional boundaries. The authoritative SKIP arrows live in the
  // frontmatter description, so a real loop declared there is still caught — this narrows scope,
  // it does not relax the match.
  const myEdges = skipEdges(fm);
  const seen = new Set();
  for (const e of myEdges) {
    if (e.id === skillId || seen.has(e.id)) continue;
    const peerPath = path.join(SKILLS_ROOT, e.id, 'SKILL.md');
    if (!fs.existsSync(peerPath)) continue; // only resolvable top-level peers
    const peerFm = frontmatterBlock(fs.readFileSync(peerPath, 'utf8'));
    const back = (peerFm ? skipEdges(peerFm) : []).filter((pe) => pe.id === skillId);
    const shared = back.flatMap((b) => [...b.toks].filter((t) => e.toks.has(t)));
    if (shared.length > 0) {
      seen.add(e.id);
      add('circular-skip', 'Medium',
        `Mutual SKIP on shared topic [${[...new Set(shared)].slice(0, 4).join(', ')}]: "${skillId}" <-> "${e.id}" defer to each other for the same topic — resolve into a one-way specificity gradient.`);
    }
  }
}

// ---- check 6: tier-config presence (hub only) ------------------------------
if (isHub && ownerManifestFile) {
  if (!tierConfig) {
    add('tier-config', 'Medium', `No tier-config.json at ${TIER_CONFIG_PATH}; cannot confirm "${ownerManifestFile}" is tiered.`);
  } else if (!Array.isArray(tierConfig.manifests) || !tierConfig.manifests.includes(ownerManifestFile)) {
    if (doRegisterTier && doApply) {
      tierConfig.manifests = Array.isArray(tierConfig.manifests) ? tierConfig.manifests : [];
      tierConfig.manifests.push(ownerManifestFile);
      fs.writeFileSync(TIER_CONFIG_PATH, JSON.stringify(tierConfig, null, 2) + '\n');
      add('tier-config', 'Low', `Registered "${ownerManifestFile}" in tier-config.json manifests[] (--register-tier applied).`);
    } else {
      add('tier-config', 'Medium',
        `"${ownerManifestFile}" is not in tier-config.json manifests[] — this family won't auto-tier. Fix: node meta-validate.mjs ${skillId} --register-tier --apply`);
    }
  }
}

// ---- check 7: description cap-compliance -----------------------------------
// The frontmatter `description` is the always-on triggering signal; the body loads only AFTER the
// hub is selected. So a description over the per-entry cap is truncated in the index (lost triggers)
// and inflates every-session token cost. Fix is NEVER to drop spoke keywords — split a hub into
// themed sub-hubs, or compress prose while preserving every enumerated keyword.
if (fm) {
  const dlen = descriptionLen(fm);
  if (dlen > DESC_CAP) {
    add('description-cap', 'Medium',
      isHub
        ? `description is ${dlen} chars > cap ${DESC_CAP} — split this hub into themed sub-hubs (distributing the spoke enumeration) or compress the prose while keeping every spoke keyword; do not drop enumeration.`
        : `description is ${dlen} chars > cap ${DESC_CAP} — compress under the cap without dropping trigger keywords.`);
  }
}

// ---- report ----------------------------------------------------------------
function report() {
  const summary = { High: 0, Medium: 0, Low: 0 };
  for (const f of findings) summary[f.level] = (summary[f.level] || 0) + 1;
  if (asJson) {
    process.stdout.write(JSON.stringify(
      { target, skillId, skillPath, isHub, isFoldedSpoke, ownerManifestFile, ownerFamily, summary, findings }, null, 2) + '\n');
  } else {
    console.log(`meta-validate: ${skillId}  (${isHub ? 'hub' : isFoldedSpoke ? 'folded spoke' : 'standalone'})`);
    console.log(`  ${summary.High} high · ${summary.Medium} medium · ${summary.Low} low`);
    for (const f of findings) console.log(`  [${f.level.toUpperCase().padEnd(6)}] ${f.check}: ${f.message}`);
    if (findings.length === 0) console.log('  clean — no structural findings.');
  }
  process.exit(summary.High > 0 ? 1 : 0);
}
report();
