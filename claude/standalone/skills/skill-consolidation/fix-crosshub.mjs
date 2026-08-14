#!/usr/bin/env node
// Post-test remediation: (1) prepend a provenance banner to every reference file so stale
// "use the X skill" pointers inside the verbatim copies are neutralized; (2) append a
// cross-hub map to each hub SKILL.md so a topic owned by a sibling hub is discoverable.
import fs from 'node:fs';
import path from 'node:path';

const m = JSON.parse(fs.readFileSync('consolidation-manifest.json', 'utf8'));
const SK = process.env.HOME + '/.claude/skills';
const hubs = Object.keys(m.hubs);

const BANNER_MARK = '<!-- hub-reference-banner -->';
function banner(hub, spoke) {
  return `${BANNER_MARK}
> **Reference file — part of the \`${hub}\` hub.** Formerly the standalone \`${spoke}\` skill.
> Sibling MongoDB sub-topics are now reference files under the four hubs (\`mongodb-expert\`,
> \`mongodb-atlas-expert\`, \`atlas-diagnostics-expert\`, \`mongodb-operations-expert\`) — **not**
> standalone skills. Ignore any "use the X skill" / \`related_skills\` / SKIP pointers below that
> name a bare \`mongodb-*\`/\`atlas-*\` skill; instead load that topic's \`references/<name>.md\`
> from the owning hub (see the hub's "Cross-hub map").

---
`;
}

let banners = 0;
for (const [hub, def] of Object.entries(m.hubs)) {
  for (const s of def.spokes) {
    const f = path.join(SK, hub, s.referenceFile);
    if (!fs.existsSync(f)) continue;
    let md = fs.readFileSync(f, 'utf8');
    if (md.includes(BANNER_MARK)) continue; // idempotent
    fs.writeFileSync(f, banner(hub, s.spoke) + '\n' + md);
    banners++;
  }
}

// Cross-hub map appended to each hub SKILL.md (idempotent)
const MAP_MARK = '<!-- cross-hub-map -->';
const crossMap = `
${MAP_MARK}
## Cross-hub map — where every MongoDB topic lives

All MongoDB knowledge is split across **four hubs** (plus \`mongodb-kb\` for KB-article lookups and
\`10gen\` for repo install/run). If a task's deep material is **not** in this hub's Sub-skill routing
table, it is a reference file under a sibling hub — **activate that hub or Read its \`references/<name>.md\` directly**.

| Hub | Owns | Example reference files |
| --- | --- | --- |
| \`mongodb-expert\` | Core data plane + **engine internals**: CRUD/MQL, aggregation, indexes, query performance, schema design, transactions, change streams, time-series, geospatial, views, BSON, error codes, connection strings, driver internals, **WiredTiger cache/eviction/checkpoint internals**, mongosh, database tools, multi-tenancy, sharding, replication, Compass | \`references/mongodb-wiredtiger-internals.md\`, \`mongodb-indexes-deep.md\`, \`mongodb-sharding.md\`, \`mongodb-replication.md\` |
| \`mongodb-atlas-expert\` | Atlas **cloud platform**: control plane, Atlas Search, Vector Search, Stream Processing, Charts, Data Federation, App Services, Triggers, Online Archive, Flex, networking, IAM/RBAC, Terraform, AKO | \`references/mongodb-atlas-search.md\`, \`mongodb-atlas-vector-search.md\` |
| \`atlas-diagnostics-expert\` | Live **diagnostics & performance**: ts-diag, FTDC, performance-troubleshooting symptom triage, benchmarking, monitoring/observability, capacity planning | \`references/mongodb-performance-troubleshooting.md\` |
| \`mongodb-operations-expert\` | **Ops & data movement**: backup/restore, DR, Ops Manager, upgrades, migration, mongosync, relational migrator, CDC, data lifecycle, security architecture, encryption, compliance, cost, Kafka/Spark connectors | \`references/mongosync.md\`, \`mongodb-backup-restore.md\` |

**High-overlap routing notes:**
- Performance **symptom triage** (high CPU, cache pressure, slow queries, latency spikes) starts at \`atlas-diagnostics-expert\`, but **storage-engine root-cause internals** (WiredTiger cache fill / dirty trigger / eviction threads / reconciliation / checkpoints) are owned by \`mongodb-expert\` — cross-load \`mongodb-expert/references/mongodb-wiredtiger-internals.md\` (and \`mongodb-wiredtiger.md\`) for depth.
- Migration symptoms vs migration **execution**: live-cluster diagnosis → \`atlas-diagnostics-expert\`; the migration/mongosync runbook → \`mongodb-operations-expert\`.
- Atlas Search/Vector **query syntax & index design** → \`mongodb-atlas-expert\`; the slowness *triage* of a running search → \`atlas-diagnostics-expert\`.
`;

let maps = 0;
for (const hub of hubs) {
  const f = path.join(SK, hub, 'SKILL.md');
  let md = fs.readFileSync(f, 'utf8');
  if (md.includes(MAP_MARK)) continue;
  fs.writeFileSync(f, md.replace(/\s*$/, '') + '\n' + crossMap);
  maps++;
}

console.log(`Banners prepended to reference files: ${banners}`);
console.log(`Cross-hub maps appended to hub SKILL.md: ${maps}`);
