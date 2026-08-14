#!/usr/bin/env node
// Cross-route a single spoke INTO an existing hub without disturbing the hub's existing spokes.
// Copies the spoke reference (flat, or nested with --nested for multi-file spokes), appends the
// spoke to the existing <family>-manifest.json hub, and inserts one routing-table row into the hub
// SKILL.md (after the last routing row, before the cross-hub map). Does NOT delete the spoke dir
// and does NOT trigger the Write-tool optimizer hook (pure node fs, like fix-crosshub).
// Usage: node crossroute.mjs <family-manifest.json> <hubName> <spoke> [--nested]
import fs from 'node:fs';
import path from 'node:path';

const [, , manifestPath, hub, spoke] = process.argv;
const nested = process.argv.includes('--nested');
const SK = process.env.HOME + '/.claude/skills';
if (!manifestPath || !hub || !spoke) { console.error('usage: crossroute.mjs <manifest> <hub> <spoke> [--nested]'); process.exit(1); }

function frontmatter(md) { const m = md.match(/^---\n([\s\S]*?)\n---/); return m ? m[1] : ''; }
function field(fm, name) {
  const re = new RegExp(`^${name}:\\s*(?:[|>][-+]?)?\\s*(.*)$`, 'm');
  const mm = fm.match(re); if (!mm) return '';
  let val = mm[1].trim(); if (val) return val.replace(/^["']|["']$/g, '');
  const lines = fm.split('\n'); const idx = lines.findIndex(l => new RegExp(`^${name}:`).test(l));
  const buf = []; for (let i = idx + 1; i < lines.length; i++) { if (/^\s+\S/.test(lines[i])) buf.push(lines[i].trim()); else break; }
  return buf.join(' ').replace(/^["']|["']$/g, '');
}
function oneLine(desc) {
  if (!desc) return '(no description)';
  let s = desc.split(/\bTRIGGER\b|\bSKIP\b/)[0].trim().replace(/\s+/g, ' ');
  if (s.length > 150) s = s.slice(0, 147).replace(/[,;:\s]+\S*$/, '') + '…';
  return s;
}

const srcDir = path.join(SK, spoke);
const srcSkill = path.join(srcDir, 'SKILL.md');
if (!fs.existsSync(srcSkill)) { console.error(`spoke ${spoke} has no SKILL.md`); process.exit(1); }
const md = fs.readFileSync(srcSkill, 'utf8');
const routingLine = oneLine(field(frontmatter(md), 'description'));

// 1. copy reference
let referenceFile;
const refDir = path.join(SK, hub, 'references');
fs.mkdirSync(refDir, { recursive: true });
if (nested) { fs.cpSync(srcDir, path.join(refDir, spoke), { recursive: true }); referenceFile = `references/${spoke}/SKILL.md`; }
else { fs.writeFileSync(path.join(refDir, `${spoke}.md`), md); referenceFile = `references/${spoke}.md`; }

// 2. append to manifest (dedupe)
const m = JSON.parse(fs.readFileSync(manifestPath, 'utf8'));
if (!m.hubs[hub]) { console.error(`hub ${hub} not in ${manifestPath}`); process.exit(1); }
m.hubs[hub].spokes = m.hubs[hub].spokes.filter(s => s.spoke !== spoke);
m.hubs[hub].spokes.push({ spoke, referenceFile, routingLine, srcBytes: Buffer.byteLength(md) });
fs.writeFileSync(manifestPath, JSON.stringify(m, null, 2));

// 3. insert routing-table row into hub SKILL.md (after last routing row, before cross-hub map)
const skp = path.join(SK, hub, 'SKILL.md');
let sk = fs.readFileSync(skp, 'utf8');
if (sk.includes(`\`${spoke}\``) && sk.includes(referenceFile)) { console.log(`${spoke} → ${hub}: row already present`); }
else {
  const lines = sk.split('\n');
  let mapIdx = lines.findIndex(l => l.includes('<!-- cross-hub-map -->'));
  if (mapIdx < 0) mapIdx = lines.length;
  let lastRow = -1;
  for (let i = 0; i < mapIdx; i++) if (/^\|\s*`[^`]+`\s*\|.*\|\s*$/.test(lines[i])) lastRow = i;
  const row = `| \`${spoke}\` | ${routingLine.replace(/\|/g, '\\|')} | \`${referenceFile}\` |`;
  if (lastRow >= 0) lines.splice(lastRow + 1, 0, row);
  else { console.error(`no routing table found in ${hub}`); process.exit(1); }
  fs.writeFileSync(skp, lines.join('\n'));
}
console.log(`crossrouted ${spoke} → ${hub} (${referenceFile})`);
