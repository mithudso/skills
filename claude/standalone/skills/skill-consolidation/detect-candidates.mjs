#!/usr/bin/env node
// detect-candidates.mjs — READ-ONLY detector for new hub-consolidation candidates.
//
// Lists top-level ~/.claude/skills/*/ dirs that are NOT already:
//   - a registered hub (owns a family in some *-manifest.json),
//   - a meta/infra tool that must never be hubbed (see EXCLUDE_LIST), or
//   - already a spoke folded under any *-manifest.json.
// Clusters the remainder into theme buckets by keyword/prefix heuristics and reports any
// bucket with >= READY_THRESHOLD members that is NOT already a registered family/hub
// as a "ready candidate".
//
// This script NEVER mutates skills. Detection only. Mutation is done by the scheduled
// auto-hub agent (see auto-hub-agent-prompt.md), and only when a real ready candidate exists.
//
// Usage:
//   node detect-candidates.mjs           # human-readable summary (default)
//   node detect-candidates.mjs --json    # machine-readable JSON

import fs from 'node:fs';
import path from 'node:path';
import os from 'node:os';
import { fileURLToPath } from 'node:url';
import { loadHubRegistry } from './hub-registry.mjs';
import { EXCLUDE_LIST } from './exclude-list.mjs';

const SKILLS_ROOT = path.join(os.homedir(), '.claude', 'skills');
const CONS_ROOT = path.dirname(fileURLToPath(import.meta.url));
const READY_THRESHOLD = 8; // a bucket with >= this many unhubbed members is "ready"

const args = new Set(process.argv.slice(2));
const asJson = args.has('--json');

// --- Meta / infra tools that must NEVER be hubbed ---------------------------------
// Hub skills themselves are detected dynamically from the manifests. The cross-cutting
// meta/infra exclude set is single-sourced in ./exclude-list.mjs (imported above),
// shared with audit-placement.mjs.

// --- Theme buckets (easy to extend: add { name, test } rows) ----------------------
// Each category is a regex (or predicate) tested against the skill dir name. First match
// wins, in declaration order, so put more-specific buckets before the catch-all.
const CATEGORIES = [
  {
    name: 'security',
    test: (n) => /(^|-)(security|secur|auth|vault|crypto|oauth|okta|compliance|http-security)(-|$)/.test(n),
  },
  {
    name: 'document-formats',
    test: (n) => /^(pdf|docx|pptx|xlsx|csv|json-advanced)$/.test(n)
      || /(document-format|drawio|spreadsheet|slide|markdown-doc)/.test(n),
  },
  {
    name: 'tam-ops',
    test: (n) => /(^|-)(tam|account|case|customer|operator|firedrill|autoremediation|incident|monday-board|slack-subscription)(-|$)/.test(n)
      || /(report-generator|file-consolidator|timeline-visualization)/.test(n),
  },
  {
    name: 'testing',
    test: (n) => /(^|-)(test|testing|vitest|e2e|webapp-test|backtest)(-|$)/.test(n)
      || /(integration-tester|webapp-testing)/.test(n),
  },
  {
    name: 'data-integration-extraction',
    test: (n) => /(scraping|dom-scraping|granola|plaud|iterative-retrieval|doc-store|doc-archaeology|database-migrations|databases-)/.test(n),
  },
  {
    name: 'catch-all',
    test: () => true,
  },
];

function categorize(name) {
  for (const c of CATEGORIES) if (c.test(name)) return c.name;
  return 'catch-all';
}

function listSkillDirs() {
  if (!fs.existsSync(SKILLS_ROOT)) return [];
  return fs
    .readdirSync(SKILLS_ROOT, { withFileTypes: true })
    .filter((d) => d.isDirectory() && !d.name.startsWith('.'))
    .map((d) => d.name)
    .sort();
}

function loadManifests() {
  // Delegates to the shared hub-registry loader (audit rec #3); consolidation-manifest.json
  // is excluded there exactly as it was here. Return shape preserved for downstream callers.
  const { hubs, spokes, families, manifestFiles } = loadHubRegistry({ consRoot: CONS_ROOT });
  return { hubs, spokes, families, manifestFiles };
}

function main() {
  const live = listSkillDirs();
  const { hubs, spokes, families, manifestFiles } = loadManifests();

  const excludedReasons = {};
  const remaining = [];
  for (const name of live) {
    if (hubs.has(name)) {
      excludedReasons[name] = 'registered-hub';
      continue;
    }
    if (EXCLUDE_LIST.has(name)) {
      excludedReasons[name] = 'meta-infra';
      continue;
    }
    if (spokes.has(name)) {
      excludedReasons[name] = 'already-spoke';
      continue;
    }
    remaining.push(name);
  }

  // Cluster the remainder.
  const buckets = {};
  for (const name of remaining) {
    const cat = categorize(name);
    (buckets[cat] ||= []).push(name);
  }

  // Families the user has explicitly chosen NOT to hub (the scheduler must never auto-consolidate these).
  const HELD_FAMILIES = new Set(['tam-ops']);

  // A bucket is a "ready candidate" if it has >= threshold members AND its category name
  // is not already a registered family AND is not on the user's hold list.
  // (catch-all is never a candidate — it's noise.)
  const candidates = [];
  for (const [cat, members] of Object.entries(buckets)) {
    if (cat === 'catch-all') continue;
    if (HELD_FAMILIES.has(cat)) continue; // user-held: surfaced in report but never auto-hubbed
    const alreadyFamily = families.has(cat) || hubs.has(cat);
    if (members.length >= READY_THRESHOLD && !alreadyFamily) {
      candidates.push({ category: cat, count: members.length, members: members.sort() });
    }
  }
  candidates.sort((a, b) => b.count - a.count);

  const result = {
    generatedAt: new Date().toISOString(),
    skillsRoot: SKILLS_ROOT,
    readyThreshold: READY_THRESHOLD,
    counts: {
      liveSkillDirs: live.length,
      registeredHubs: hubs.size,
      foldedSpokes: spokes.size,
      excludedMetaInfra: Object.values(excludedReasons).filter((r) => r === 'meta-infra').length,
      remaining: remaining.length,
    },
    manifests: manifestFiles.sort(),
    registeredFamilies: [...families].sort(),
    buckets: Object.fromEntries(
      Object.entries(buckets).map(([k, v]) => [k, v.sort()])
    ),
    readyCandidates: candidates,
  };

  if (asJson) {
    process.stdout.write(JSON.stringify(result, null, 2) + '\n');
    return;
  }

  // Human summary
  const L = [];
  L.push('Skill hub-candidate detector (read-only)');
  L.push('='.repeat(42));
  L.push(`Generated: ${result.generatedAt}`);
  L.push(`Skills root: ${SKILLS_ROOT}`);
  L.push('');
  L.push(`Live skill dirs: ${result.counts.liveSkillDirs}`);
  L.push(`  - registered hubs:      ${result.counts.registeredHubs}`);
  L.push(`  - folded spokes:        ${result.counts.foldedSpokes}`);
  L.push(`  - excluded meta/infra:  ${result.counts.excludedMetaInfra}`);
  L.push(`  - remaining (unhubbed): ${result.counts.remaining}`);
  L.push('');
  L.push(`Registered families: ${result.registeredFamilies.join(', ') || '(none)'}`);
  L.push('');
  L.push(`Clusters of remaining unhubbed skills (ready threshold = ${READY_THRESHOLD}):`);
  for (const [cat, members] of Object.entries(result.buckets).sort((a, b) => b[1].length - a[1].length)) {
    const flag = result.readyCandidates.some((c) => c.category === cat) ? '  <== READY CANDIDATE' : '';
    L.push(`  [${members.length}] ${cat}${flag}`);
    L.push(`        ${members.join(', ')}`);
  }
  L.push('');
  if (result.readyCandidates.length === 0) {
    L.push('No ready candidates (no unhubbed cluster reached the threshold).');
  } else {
    L.push(`READY CANDIDATES (>= ${READY_THRESHOLD}, not yet a hub/family):`);
    for (const c of result.readyCandidates) {
      L.push(`  - ${c.category} (${c.count}): ${c.members.join(', ')}`);
    }
  }
  process.stdout.write(L.join('\n') + '\n');
}

main();
