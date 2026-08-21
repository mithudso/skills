#!/usr/bin/env node
// Hub-and-spoke consolidation builder.
// For each hub: copy each spoke's SKILL.md verbatim into <hub>/references/<spoke>.md,
// extract a one-line routing description, emit a routing-table fragment, and record a manifest.
// Does NOT remove spoke dirs and does NOT rewrite hub frontmatter — those are deliberate later steps.
import fs from 'node:fs';
import path from 'node:path';

const mapPath = process.argv[2] || 'mongodb-mapping.json';
const m = JSON.parse(fs.readFileSync(mapPath, 'utf8'));
// skillsRoot is the parent of this script's dir (<skillsRoot>/skill-consolidation/build.mjs),
// derived for portability; the mapping file's skillsRoot field is advisory only.
const SK = path.dirname(path.dirname(new URL(import.meta.url).pathname));
const outDir = path.dirname(path.resolve(mapPath));

function frontmatter(md) {
  const mm = md.match(/^---\n([\s\S]*?)\n---/);
  return mm ? mm[1] : '';
}
function field(fm, name) {
  // handles `name: value`, `name: >- \n  folded`, `name: |\n  block`
  const re = new RegExp(`^${name}:\\s*(?:[|>][-+]?)?\\s*(.*)$`, 'm');
  const mm = fm.match(re);
  if (!mm) return '';
  let val = mm[1].trim();
  if (val) return val.replace(/^["']|["']$/g, '');
  // folded/block scalar: gather subsequent indented lines
  const lines = fm.split('\n');
  const idx = lines.findIndex(l => new RegExp(`^${name}:`).test(l));
  const buf = [];
  for (let i = idx + 1; i < lines.length; i++) {
    if (/^\s+\S/.test(lines[i])) buf.push(lines[i].trim());
    else break;
  }
  return buf.join(' ').replace(/^["']|["']$/g, '');
}
function oneLine(desc) {
  if (!desc) return '(no description)';
  // take text before TRIGGER/SKIP, first sentence-ish, cap length
  let s = desc.split(/\bTRIGGER\b|\bSKIP\b/)[0].trim();
  s = s.replace(/\s+/g, ' ');
  if (s.length > 160) s = s.slice(0, 157).replace(/[,;:\s]+\S*$/, '') + '…';
  return s;
}

const manifest = { family: m.family, builtAt: new Date().toISOString(), hubs: {} };
let copied = 0, missing = [];

for (const [hub, def] of Object.entries(m.hubs)) {
  const hubDir = path.join(SK, hub);
  const refDir = path.join(hubDir, 'references');
  fs.mkdirSync(refDir, { recursive: true });
  const rows = [];
  manifest.hubs[hub] = { keepExisting: def.keepExisting, title: def.title, spokes: [] };
  for (const spoke of def.spokes) {
    const src = path.join(SK, spoke, 'SKILL.md');
    if (!fs.existsSync(src)) { missing.push(spoke); continue; }
    const md = fs.readFileSync(src, 'utf8');
    const fm = frontmatter(md);
    const desc = field(fm, 'description');
    const refName = `${spoke}.md`;
    fs.writeFileSync(path.join(refDir, refName), md);
    copied++;
    const line = oneLine(desc);
    rows.push({ spoke, refName, line });
    manifest.hubs[hub].spokes.push({ spoke, referenceFile: `references/${refName}`, routingLine: line, srcBytes: Buffer.byteLength(md) });
  }
  // routing-table fragment for this hub
  const frag = [
    `<!-- ROUTING TABLE: ${hub} — auto-generated, edit descriptions as needed -->`,
    `## Sub-skill routing table`,
    ``,
    `This hub absorbs ${rows.length} former standalone skills as on-demand reference files. When a task matches a row, **Read the listed \`references/\` file** before answering — do not rely on this table alone for depth.`,
    ``,
    `| Sub-topic | When to load | Reference file |`,
    `| --- | --- | --- |`,
    ...rows.map(r => `| \`${r.spoke}\` | ${r.line.replace(/\|/g, '\\|')} | \`references/${r.refName}\` |`),
    ``
  ].join('\n');
  fs.writeFileSync(path.join(outDir, `${hub}.routing.md`), frag);
}

fs.writeFileSync(path.join(outDir, `${m.family}-manifest.json`), JSON.stringify(manifest, null, 2));
fs.writeFileSync(path.join(outDir, `consolidation-manifest.json`), JSON.stringify(manifest, null, 2));
console.log(`Reference files copied: ${copied}`);
console.log(`Missing spokes (skipped): ${missing.length ? missing.join(', ') : 'none'}`);
for (const [hub, h] of Object.entries(manifest.hubs)) console.log(`  ${hub}: ${h.spokes.length} references`);
console.log(`Routing fragments + consolidation-manifest.json written to ${outDir}`);
