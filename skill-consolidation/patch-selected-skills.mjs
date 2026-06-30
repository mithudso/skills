import fs from 'node:fs';
const p = process.env.HOME + '/Documents/GitHub/mdb-context-hub/scripts/skill-pack.config.mjs';
let lines = fs.readFileSync(p, 'utf8').split('\n');
const before = lines.length;
// idempotent: strip any stray hub lines (repairs a prior mis-insert) AND remove the 5 deletions
const hubIds = ['document-formats', 'content-ingestion-extraction', 'tam-operations'];
const del = ['ebay-listing', 'gog', 'ordercli', 'dmux-workflows', 'etetoolkit'];
lines = lines.filter(l => ![...hubIds, ...del].some(id => l.includes(`{ id: '${id}',`)));
const hubs = [
  `  { id: 'document-formats', category: 'developer', priorityBucket: 'developer', tags: ['installed', 'developer', 'hub'], localContextPath: 'local-sources/document-formats/context.md', localManifestPath: 'local-sources/document-formats/manifest.yaml' },`,
  `  { id: 'content-ingestion-extraction', category: 'developer', priorityBucket: 'developer', tags: ['installed', 'developer', 'hub'], localContextPath: 'local-sources/content-ingestion-extraction/context.md', localManifestPath: 'local-sources/content-ingestion-extraction/manifest.yaml' },`,
  `  { id: 'tam-operations', category: 'custom', priorityBucket: 'custom', tags: ['installed', 'tam', 'hub'], localContextPath: 'local-sources/tam-operations/context.md', localManifestPath: 'local-sources/tam-operations/manifest.yaml' },`,
];
// anchor on the unique landmark, then find the '];' that closes SELECTED_SKILLS just above it
const metaIdx = lines.findIndex(l => l.includes('export const MANUAL_SKILL_METADATA'));
if (metaIdx < 0) { console.error('landmark not found'); process.exit(1); }
let closeIdx = -1;
for (let i = metaIdx; i >= 0; i--) { if (lines[i].trim() === '];') { closeIdx = i; break; } }
if (closeIdx < 0) { console.error('SELECTED_SKILLS close not found'); process.exit(1); }
lines.splice(closeIdx, 0, ...hubs);
fs.writeFileSync(p, lines.join('\n'));
console.log(`SELECTED_SKILLS patched: ${before} -> ${lines.length} lines; hubs inserted before line ${closeIdx + 1}`);
