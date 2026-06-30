// hub-registry.mjs — single source of truth for reading the hub-and-spoke manifests.
//
// Recommendation #3 of the tooling-integration audit: /dr (Phase 0/2), concept-family-explorer,
// skill-optimizer (Pass O), referents.mjs, and detect-candidates.mjs each independently parsed
// `~/.claude/skill-consolidation/*-manifest.json` to resolve spoke -> hub. When the manifest shape
// changed, every call site drifted. This module centralizes that parse so there is one place to fix.
//
// A manifest has shape: { family, builtAt, hubs: { <hubId>: { title, keepExisting, spokes:
// [ { spoke, referenceFile, ... } | "<spoke-id>" ] } } }. `consolidation-manifest.json` is always a
// byte-for-byte duplicate of the most-recently-built family manifest (build.mjs writes both), so it
// is excluded by default — including it only re-sets identical spoke->hub entries.

import fs from 'node:fs';
import path from 'node:path';
import os from 'node:os';

// Derive from this module's own location (…/skill-consolidation/hub-registry.mjs) so the
// system is portable and survives `git pull`; fall back to the legacy path only if that fails.
export const DEFAULT_CONS_ROOT = (() => {
  try { return path.dirname(new URL(import.meta.url).pathname); }
  catch { return path.join(os.homedir(), '.claude', 'skills', 'skill-consolidation'); }
})();

/**
 * Load and parse every family manifest under consRoot.
 * @param {object} [opts]
 * @param {string} [opts.consRoot] - directory holding the *-manifest.json files.
 * @param {boolean} [opts.includeConsolidation] - include consolidation-manifest.json (default false; it is a duplicate).
 * @returns {{ spokeHub: Map<string,string>, hubs: Set<string>, spokes: Set<string>, families: Set<string>, manifestFiles: string[] }}
 */
export function loadHubRegistry(opts = {}) {
  const consRoot = opts.consRoot || DEFAULT_CONS_ROOT;
  const includeConsolidation = opts.includeConsolidation === true;
  const spokeHub = new Map();
  const hubs = new Set();
  const spokes = new Set();
  const families = new Set();
  const manifestFiles = [];

  let entries;
  try {
    entries = fs.readdirSync(consRoot);
  } catch {
    return { spokeHub, hubs, spokes, families, manifestFiles };
  }

  const files = entries.filter(
    (f) => f.endsWith('-manifest.json') && (includeConsolidation || f !== 'consolidation-manifest.json'),
  );

  for (const f of files) {
    let m;
    try {
      m = JSON.parse(fs.readFileSync(path.join(consRoot, f), 'utf8'));
    } catch {
      continue; // skip unparseable manifest rather than crash the caller
    }
    if (!m || !m.hubs) continue;
    manifestFiles.push(f);
    if (m.family) families.add(m.family);
    for (const [hub, def] of Object.entries(m.hubs)) {
      hubs.add(hub);
      for (const s of def.spokes || []) {
        const spoke = typeof s === 'string' ? s : s && s.spoke;
        if (spoke) {
          spokes.add(spoke);
          spokeHub.set(spoke, hub);
        }
      }
    }
  }
  return { spokeHub, hubs, spokes, families, manifestFiles };
}

/**
 * Resolve a skill id to its owning hub, or null if it is not a folded spoke.
 * A null result means "either a top-level skill or unknown" — the caller decides which.
 */
export function resolveSpokeHub(id, reg) {
  const r = reg || loadHubRegistry();
  return r.spokeHub.get(id) || null;
}

/** Convenience: the spoke -> hub Map only (matches the old loadSpokeMap() return shape). */
export function loadSpokeMap(opts = {}) {
  return loadHubRegistry(opts).spokeHub;
}
