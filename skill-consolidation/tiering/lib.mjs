// Shared helpers for the skill auto-tiering system.
import fs from 'node:fs';
import path from 'node:path';

export const TIER_DIR = path.dirname(new URL(import.meta.url).pathname);
export const CONFIG_PATH = path.join(TIER_DIR, 'tier-config.json');
export const LOG_PATH = path.join(TIER_DIR, 'access-log.jsonl');
export const STATE_PATH = path.join(TIER_DIR, 'tier-state.json');
export const BANNER_MARK = '<!-- hub-reference-banner -->';

export function loadConfig() {
  const cfg = JSON.parse(fs.readFileSync(CONFIG_PATH, 'utf8'));
  // Derive the roots from this file's own location so the system is portable
  // across machines/users and survives `git pull` — the values stored in
  // tier-config.json are treated as advisory only. Layout is fixed by the repo:
  //   <skillsRoot>/skill-consolidation/tiering/lib.mjs
  const consRoot = path.dirname(TIER_DIR);            // .../skill-consolidation
  cfg.consolidationRoot = consRoot;
  cfg.skillsRoot = path.dirname(consRoot);            // .../skills
  return cfg;
}

// Build spoke -> {hub, referenceFile, refAbsPath, skillDirAbsPath} from the family manifests.
export function loadSpokeMap(cfg) {
  const map = new Map();
  for (const mf of cfg.manifests) {
    const manifest = JSON.parse(fs.readFileSync(path.join(cfg.consolidationRoot, mf), 'utf8'));
    for (const [hub, def] of Object.entries(manifest.hubs)) {
      for (const s of def.spokes) {
        map.set(s.spoke, {
          hub,
          family: manifest.family,
          referenceFile: s.referenceFile,
          refAbsPath: path.join(cfg.skillsRoot, hub, s.referenceFile),
          skillDirAbsPath: path.join(cfg.skillsRoot, s.spoke),
          skillMdAbsPath: path.join(cfg.skillsRoot, s.spoke, 'SKILL.md')
        });
      }
    }
  }
  return map;
}

export function loadState() {
  try { return JSON.parse(fs.readFileSync(STATE_PATH, 'utf8')); } catch { return {}; }
}
export function saveState(state) {
  fs.writeFileSync(STATE_PATH, JSON.stringify(state, null, 2));
}

// Reference file = provenance banner + original SKILL.md. Strip the banner to recover the original.
export function stripBanner(md) {
  if (!md.includes(BANNER_MARK)) return md;
  return md.replace(/^<!--\s*hub-reference-banner\s*-->[\s\S]*?\n---\n/, '').replace(/^\s+/, '');
}
// Re-add a minimal banner when refreshing a reference from a (possibly edited) hot standalone.
export function addBanner(originalMd, hub, spoke, hubList) {
  const banner = `${BANNER_MARK}
> **Reference file — part of the \`${hub}\` hub.** Formerly the standalone \`${spoke}\` skill.
> Sibling topics in this family are now reference files under the hubs (${hubList}) — **not** standalone
> skills. Ignore any "use the X skill" / \`related_skills\` / SKIP pointers below that name a bare sibling
> skill; load that topic's \`references/<name>.md\` from the owning hub (see the hub's "Cross-hub map").

---
`;
  return banner + '\n' + originalMd;
}

export function appendLog(entry) {
  fs.appendFileSync(LOG_PATH, JSON.stringify(entry) + '\n');
}
export function readLog() {
  try {
    return fs.readFileSync(LOG_PATH, 'utf8').trim().split('\n').filter(Boolean).map(l => {
      try { return JSON.parse(l); } catch { return null; }
    }).filter(Boolean);
  } catch { return []; }
}
