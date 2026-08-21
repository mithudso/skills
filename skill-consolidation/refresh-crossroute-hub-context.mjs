import fs from 'node:fs';
const SK = process.env.HOME + '/.claude/skills';
const LS = process.env.HOME + '/Documents/GitHub/mdb-context-hub/local-sources';
const hubs = ['frontend-ui', 'chrome-extension-expert', 'ai-agent-engineering', 'programming-languages', 'writing-expert'];
for (const h of hubs) {
  const src = `${SK}/${h}/SKILL.md`;
  if (!fs.existsSync(src)) { console.log(`SKILL.md MISSING for ${h}`); continue; }
  const md = fs.readFileSync(src, 'utf8');
  const body = md.replace(/^---\n[\s\S]*?\n---\n/, '').replace(/^\s+/, '');
  const destDir = `${LS}/${h}`;
  if (!fs.existsSync(destDir)) { console.log(`local-sources/${h} MISSING (skipping ${h})`); continue; }
  const dest = `${destDir}/context.md`;
  const prev = fs.existsSync(dest) ? fs.readFileSync(dest, 'utf8') : '';
  fs.writeFileSync(dest, body);
  const newRows = (body.match(/^\|\s*`[^`]+`\s*\|/gm) || []).length;
  console.log(`${h}: context.md ${prev.length}B -> ${body.length}B (${newRows} routing rows)`);
}
