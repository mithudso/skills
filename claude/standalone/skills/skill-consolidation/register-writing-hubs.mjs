import fs from 'node:fs';
const SK = process.env.HOME + '/.claude/skills';
const HUBROOT = process.env.HOME + '/Documents/GitHub/mdb-context-hub';
const cfgPath = `${HUBROOT}/scripts/skill-pack.config.mjs`;
const hubs = ['writing-expert', 'technical-writing-craft', 'content-and-marketing-writing', 'career-and-formal-writing'];

// (a) refresh each hub's local-sources/context.md from its current SKILL.md body
for (const h of hubs) {
  const src = `${SK}/${h}/SKILL.md`;
  const lsCtx = `${HUBROOT}/local-sources/${h}/context.md`;
  if (!fs.existsSync(src)) { console.log(`SKILL.md MISSING for ${h}`); continue; }
  if (!fs.existsSync(`${HUBROOT}/local-sources/${h}`)) { console.log(`local-sources/${h} MISSING for ${h}`); continue; }
  const body = fs.readFileSync(src, 'utf8').replace(/^---\n[\s\S]*?\n---\n/, '').replace(/^\s+/, '');
  fs.writeFileSync(lsCtx, body);
  console.log(`refreshed local-sources/${h}/context.md (${body.length}B)`);
}

// (b) pin missing hubs into SELECTED_SKILLS (idempotent), before the MANUAL_SKILL_METADATA landmark
let cfg = fs.readFileSync(cfgPath, 'utf8');
let lines = cfg.split('\n');
const metaIdx = lines.findIndex(l => l.includes('export const MANUAL_SKILL_METADATA'));
let closeIdx = -1;
for (let i = metaIdx; i >= 0; i--) { if (lines[i].trim() === '];') { closeIdx = i; break; } }
if (closeIdx < 0) { console.error('SELECTED_SKILLS close not found'); process.exit(1); }
const toAdd = [];
for (const h of hubs) {
  if (cfg.includes(`id: '${h}'`)) { console.log(`${h} already pinned in SELECTED_SKILLS`); continue; }
  toAdd.push(`  { id: '${h}', category: 'custom', priorityBucket: 'custom', tags: ['installed', 'writing', 'hub'], localContextPath: 'local-sources/${h}/context.md', localManifestPath: 'local-sources/${h}/manifest.yaml' },`);
}
if (toAdd.length) { lines.splice(closeIdx, 0, ...toAdd); fs.writeFileSync(cfgPath, lines.join('\n')); }
console.log(`pinned ${toAdd.length} writing hubs into SELECTED_SKILLS`);
