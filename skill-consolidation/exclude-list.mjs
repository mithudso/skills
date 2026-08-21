// exclude-list.mjs — single source of truth for meta/infra skills that must NEVER be hubbed
// (never clustered as a new-family candidate, never flagged homeless-needing-a-hub).
//
// Imported by detect-candidates.mjs and audit-placement.mjs — the same single-source pattern
// both scripts already use for hub-registry.mjs. Replaces two hand-synced copies that had
// drifted (detect-candidates.mjs was missing skill-creator, skill-tree-architect, and
// concept-family-explorer); this is the union set (extracted 2026-06-11).
//
// Note: skill-creator has no ~/.claude/skills/ dir today (it is a plugin skill) — its entry is
// defensive parity carried over from audit-placement.mjs.
export const EXCLUDE_LIST = new Set([
  // prompt-* / optimizer / skill meta tools
  'prompt-helper-optimizer', 'prompt-deep-optimizer', 'prompt-lookup',
  'skill-optimizer', 'skill-lookup', 'skill-creator', 'skill-tree-architect',
  'claude-code-skills', 'claude-code-plugins', 'concept-family-explorer',
  // research / repo meta tools
  'deep-research', 'deep-research-methods', 'dr',
  'repo-bootstrapper', 'repo-file-analyzer', 'repo-pattern-scanner',
  'git-workflows',
  // distinct-mechanism standalones called out in HUB-STRATEGY.md
  '10gen', 'mongodb-kb',
]);
