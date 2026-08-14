#!/usr/bin/env node
// sanitize-settings.mjs — maintainer tool.
// Reads the real ~/.claude/settings.json and emits a shareable template with:
//   • every secret value replaced by a ${PLACEHOLDER} (by env-key name or context)
//   • every absolute home path replaced by the __HOME__ token (setup.sh renders it)
// It NEVER prints a secret value: the change report shows only key-paths + placeholders.
//
// Usage:  node sanitize-settings.mjs [srcSettings] [outTemplate]
// Default src:  ~/.claude/settings.json
// Default out:  ~/.claude/skills/claude-config/settings.json.template

import { readFileSync, writeFileSync, mkdirSync } from 'node:fs';
import { homedir } from 'node:os';
import { dirname, join } from 'node:path';

const HOME = homedir();
const SRC = process.argv[2] || join(HOME, '.claude', 'settings.json');
const OUT = process.argv[3] || join(HOME, '.claude', 'skills', 'claude-config', 'settings.json.template');

const KEY_SECRET = /(secret|token|password|api[_-]?key|client[_-]?secret|access[_-]?key)/i;
// Value patterns that indicate a real credential (used only as a backstop net).
const VAL_SECRET = /(^eyJ[A-Za-z0-9_-]{10,}$)|(^sk-[A-Za-z0-9]{10,})|(xox[abpr]-[A-Za-z0-9-]+)|(^ghp_[A-Za-z0-9]{20,})|(AKIA[0-9A-Z]{16})|(^[A-Fa-f0-9]{32,}$)|(^[A-Za-z0-9+/]{40,}={0,2}$)/;

const changes = [];
const homeRe = new RegExp(HOME.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'g');

function tokenizeHome(s) {
  // collapse this machine's home AND any /Users/<name> prefix to __HOME__
  return s.replace(homeRe, '__HOME__').replace(/\/Users\/[^/"'\s]+/g, '__HOME__');
}
function isPlaceholder(s) { return /^\$\{[^}]+\}$/.test(s); }

// Walk recursively. `inMcpEnv` carries the env-key name when we are inside mcpServers.*.env.
function walk(node, path) {
  if (Array.isArray(node)) {
    for (let i = 0; i < node.length; i++) {
      const v = node[i];
      // mcpServers.*.args: the value right after a "-t" flag is the Monday API token.
      if (typeof v === 'string' && i > 0 && node[i - 1] === '-t' && !isPlaceholder(v)) {
        node[i] = '${MONDAY_API_TOKEN}';
        changes.push({ path: `${path}.${i}`, placeholder: '${MONDAY_API_TOKEN}', reason: 'after -t flag' });
        continue;
      }
      node[i] = walk(v, `${path}.${i}`);
    }
    return node;
  }
  if (node && typeof node === 'object') {
    const isEnv = /\.env$/.test(path) || path.endsWith('.env');
    for (const k of Object.keys(node)) {
      const child = node[k];
      if (typeof child === 'string' && isEnv && !isPlaceholder(child) &&
          (KEY_SECRET.test(k) || VAL_SECRET.test(child))) {
        node[k] = '${' + k + '}';
        changes.push({ path: `${path}.${k}`, placeholder: node[k], reason: 'mcp env secret' });
        continue;
      }
      node[k] = walk(child, path ? `${path}.${k}` : k);
    }
    return node;
  }
  if (typeof node === 'string') {
    let s = tokenizeHome(node);
    // Backstop: a real-looking credential anywhere OUTSIDE permission rules.
    const isRule = /^(permissions\.)/.test(path);
    if (!isRule && !isPlaceholder(s) && VAL_SECRET.test(s)) {
      const name = 'SECRET_' + path.replace(/[^A-Za-z0-9]+/g, '_').toUpperCase();
      changes.push({ path, placeholder: '${' + name + '}', reason: 'backstop value-pattern' });
      s = '${' + name + '}';
    }
    return s;
  }
  return node;
}

const raw = readFileSync(SRC, 'utf8');
const obj = JSON.parse(raw);
const out = walk(obj, '');

mkdirSync(dirname(OUT), { recursive: true });
writeFileSync(OUT, JSON.stringify(out, null, 2) + '\n');

console.log(`wrote ${OUT.replace(HOME, '~')}`);
console.log(`secret leaves replaced: ${changes.length}`);
for (const c of changes) console.log(`  ${c.path}  ->  ${c.placeholder}   (${c.reason})`);

// Self-check: fail loudly if any credential or raw home path survived.
const result = readFileSync(OUT, 'utf8');
const residualSecret = /eyJ[A-Za-z0-9_-]{20,}|sk-[A-Za-z0-9]{16,}|xox[abpr]-|ghp_[A-Za-z0-9]{20,}|AKIA[0-9A-Z]{16}/.test(result);
const residualHome = result.includes('/Users/');
console.log(`\nself-check  residual-credential=${residualSecret}  residual-home-path=${residualHome}`);
if (residualSecret || residualHome) { console.error('!! SANITIZE FAILED — not safe to commit'); process.exit(1); }
console.log('OK — template is clean');
