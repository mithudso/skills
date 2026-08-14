#!/usr/bin/env node
// Family-agnostic post-build remediation: prepend a provenance banner to every reference file,
// and append a cross-hub map (generated from the manifest) to every hub SKILL.md.
// Usage: node fix-crosshub-generic.mjs <family-manifest.json>
import fs from 'node:fs';
import path from 'node:path';

const manifestPath = process.argv[2];
if (!manifestPath) { console.error('usage: fix-crosshub-generic.mjs <manifest.json>'); process.exit(1); }
const m = JSON.parse(fs.readFileSync(manifestPath, 'utf8'));
const SK = process.env.HOME + '/.claude/skills';
const hubs = Object.keys(m.hubs);

const BANNER_MARK = '<!-- hub-reference-banner -->';
const hubList = hubs.map(h => `\`${h}\``).join(', ');
function banner(hub, spoke) {
  return `${BANNER_MARK}
> **Reference file — part of the \`${hub}\` hub.** Formerly the standalone \`${spoke}\` skill.
> Sibling topics in this family are now reference files under the hubs (${hubList}) — **not** standalone
> skills. Ignore any "use the X skill" / \`related_skills\` / SKIP pointers below that name a bare sibling
> skill; load that topic's \`references/<name>.md\` from the owning hub (see the hub's "Cross-hub map").

---
`;
}

let banners = 0;
for (const [hub, def] of Object.entries(m.hubs)) {
  for (const s of def.spokes) {
    const f = path.join(SK, hub, s.referenceFile);
    if (!fs.existsSync(f)) continue;
    let md = fs.readFileSync(f, 'utf8');
    if (md.includes(BANNER_MARK)) continue;
    fs.writeFileSync(f, banner(hub, s.spoke) + '\n' + md);
    banners++;
  }
}

const MAP_MARK = '<!-- cross-hub-map -->';
const rows = Object.entries(m.hubs).map(([hub, h]) => {
  const examples = h.spokes.slice(0, 4).map(s => `\`references/${path.basename(s.referenceFile)}\``).join(', ');
  return `| \`${hub}\` | ${h.title || hub} | ${examples}${h.spokes.length > 4 ? ', …' : ''} |`;
});
const crossMap = `
${MAP_MARK}
## Cross-hub map — where every ${m.family} topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or \`Read\` its
\`references/<name>.md\` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
${rows.join('\n')}
`;

let maps = 0;
for (const hub of hubs) {
  const f = path.join(SK, hub, 'SKILL.md');
  if (!fs.existsSync(f)) continue;
  let md = fs.readFileSync(f, 'utf8');
  if (md.includes(MAP_MARK)) continue;
  fs.writeFileSync(f, md.replace(/\s*$/, '') + '\n' + crossMap);
  maps++;
}

console.log(`[${m.family}] banners: ${banners}, cross-hub maps: ${maps}`);
