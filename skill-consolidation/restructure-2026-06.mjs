#!/usr/bin/env node
// One-off restructure (2026-06): split 3 mega-hubs into sub-hubs (parents kept as thin routers),
// create aws-cloud + security-review hubs, and fold 13 homeless skills into hubs/sub-hubs.
// Reuses the tiering lib so every moved reference stays engine-compatible (banner round-trips).
// Default = DRY-RUN. Pass --apply to write.
import fs from 'node:fs';
import path from 'node:path';
import { execFileSync } from 'node:child_process';
import { stripBanner, addBanner } from './tiering/lib.mjs';

const APPLY = process.argv.includes('--apply');
const SK = process.env.HOME + '/.claude/skills';
const CR = process.env.HOME + '/.claude/skill-consolidation';
const log = (...a) => console.log(...a);
const actions = [];

function origFromRef(p) { return stripBanner(fs.readFileSync(p, 'utf8')).replace(/^\s+/, ''); }
function origFromStandalone(p) { return fs.readFileSync(p, 'utf8'); } // no banner
function routingLine(orig) {
  // first non-empty content line of the description block
  const m = orig.match(/description:\s*>?\s*\n?\s*(.*)/);
  let line = (m ? m[1] : '').trim();
  if (!line && m) line = '';
  return (line || orig.split('\n').find(l => l.trim()) || '').slice(0, 95);
}
function writeFile(p, body) { actions.push(['write', p]); if (APPLY) { fs.mkdirSync(path.dirname(p), { recursive: true }); fs.writeFileSync(p, body); } }
function rmDir(p) { actions.push(['rmdir', p]); if (APPLY) fs.rmSync(p, { recursive: true, force: true }); }
function rmIfEmpty(p) { if (APPLY && fs.existsSync(p) && fs.readdirSync(p).length === 0) { fs.rmdirSync(p); actions.push(['rmdir-empty', p]); } }

function skillMd(name, description) {
  return `---\nname: ${name}\ndescription: >-\n${description.split('\n').map(l => '  ' + l).join('\n')}\norigin: local\n---\n\n# ${name}\n\n${description.split('TRIGGER')[0].trim()}\n\nThis hub routes to on-demand reference files under \`references/\`. See each spoke for depth.\n`;
}

// Produce one reference file for a spoke, return its manifest entry. src: {type:'ref',path}|{type:'standalone',dir}
function placeSpoke(spoke, destHubDir, destHub, hubListStr, src) {
  const orig = src.type === 'ref' ? origFromRef(src.path) : origFromStandalone(path.join(src.dir, 'SKILL.md'));
  const refRel = `references/${spoke}.md`;
  const refAbs = path.join(destHubDir, refRel);
  writeFile(refAbs, addBanner(orig, destHub, spoke, hubListStr));
  if (src.type === 'ref') { actions.push(['rm-oldref', src.path]); if (APPLY) fs.rmSync(src.path, { force: true }); }
  if (src.type === 'standalone') rmDir(src.dir);
  return { spoke, referenceFile: refRel, routingLine: routingLine(orig), srcBytes: Buffer.byteLength(orig) };
}

// ---- SPLIT + NEW-HUB + FOLD SPECS ------------------------------------------------
const D = (s) => s.replace(/\n\s+/g, ' ').trim();

const splits = [
  {
    parent: 'devops-infra', family: 'devops', manifest: 'devops-manifest.json',
    router: D(`DevOps / infrastructure / observability family ROUTER. Split into focused sub-hubs — route to:
      devops-linux-internals (kernel, boot, memory/NUMA, storage/filesystems, virtualization, io_uring, cgroups/namespaces, sandboxing, immutable Linux, privilege);
      devops-linux-admin (sysadmin, systemd, package management, shell scripting, host networking diagnostics);
      devops-containers-cicd (Docker, Kubernetes, CI/CD pipelines, Terraform/Kafka IaC, git workflows, library packaging);
      devops-observability (Node/OTel observability, Pino logging, Sentry, eBPF, Linux perf tracing). Pick the sub-hub matching the task.`),
    subhubs: [
      { id: 'devops-linux-internals', desc: D(`Linux kernel & OS-internals sub-hub (devops family).
        TRIGGER: Linux kernel architecture, boot/init, memory & NUMA, storage & filesystems, virtualization/KVM, io_uring async I/O,
        cgroups v2 & namespaces (PID/mount/net/user/cgroup), sandboxing & confinement, immutable/atomic Linux, Linux/macOS privilege model.
        SKIP: sysadmin/systemd/packaging/shell/host-networking → devops-linux-admin; containers/k8s/CI-CD/IaC → devops-containers-cicd;
        logging/tracing/metrics/eBPF/perf → devops-observability.`),
        spokes: [['linux-kernel-architecture','ref'],['linux-boot-init','ref'],['linux-memory-numa','ref'],['linux-storage-filesystems','ref'],['linux-virtualization','ref'],['io-uring-async-io','ref'],['linux-cgroups-namespaces','ref'],['linux-sandboxing-confinement','ref'],['immutable-atomic-linux','ref'],['linux-mac-privilege','ref']] },
      { id: 'devops-linux-admin', desc: D(`Linux/macOS sysadmin & host-operations sub-hub (devops family).
        TRIGGER: sysadmin troubleshooting (high CPU, OOM, disk full, DNS, port conflicts), systemd units/journald, package management,
        Bash/Zsh production shell scripting, host network diagnostics (ip, ss, dig, tcpdump, nmap, iptables/nftables).
        SKIP: kernel/memory/storage/namespaces internals → devops-linux-internals; Docker/Kubernetes/CI-CD/IaC → devops-containers-cicd;
        logging/tracing/metrics/eBPF/perf → devops-observability.`),
        spokes: [['linux-sysadmin','ref'],['systemd','ref'],['linux-package-management','ref'],['shell-scripting','ref'],['linux-networking-stack','ref']] },
      { id: 'devops-containers-cicd', desc: D(`Containers, orchestration, CI/CD, IaC & delivery sub-hub (devops family).
        TRIGGER: Docker/Dockerfile, image optimization, container runtime; Kubernetes workloads/networking/ingress; CI/CD pipeline design
        (GitHub Actions, reusable workflows, caching, release automation); Terraform/OpenTofu & Apache Kafka infra; git branching/merge/rebase,
        conventional commits, branch protection, monorepo CI, semantic-release; library packaging & distribution (npm/PyPI/crates, ESM/CJS, monorepo, semver, provenance).
        SKIP: Linux kernel/admin internals → devops-linux-internals / devops-linux-admin; logging/tracing/metrics → devops-observability.`),
        spokes: [['docker-containers','ref'],['kubernetes-networking','ref'],['cicd-pipelines','ref'],['terraform-kafka-infra','ref'],['git-workflows','standalone'],['code-packaging','standalone']] },
      { id: 'devops-observability', desc: D(`Observability, logging, tracing & performance sub-hub (devops family).
        TRIGGER: Node.js observability & OpenTelemetry instrumentation; Pino v9+ structured logging (child loggers, serializers, redaction);
        Sentry error monitoring & performance tracing; eBPF observability/networking/security (bpftrace/bcc, CO-RE/BTF/libbpf); Linux perf tracing & profiling.
        SKIP: container/CI-CD/IaC → devops-containers-cicd; Linux kernel/admin → devops-linux-internals / devops-linux-admin.`),
        spokes: [['nodejs-observability','ref'],['pino-structured-logging','ref'],['sentry-monitoring','ref'],['ebpf-observability','ref'],['linux-perf-tracing','ref']] },
    ],
  },
  {
    parent: 'programming-languages', family: 'lang', manifest: 'lang-manifest.json',
    router: D(`Programming-languages family ROUTER. Split into: lang-python (Python idioms, testing, typing, uv toolchain, packaging, CPython internals, pydantic);
      lang-js-ts (JavaScript/Node, TypeScript, Deno/Bun/edge runtimes, V8 internals, zod, JS debugging); lang-go-and-mobile (Go patterns, Kotlin/Compose Multiplatform).
      Route to the sub-hub for the language in question.`),
    subhubs: [
      { id: 'lang-python', desc: D(`Python language sub-hub (lang family).
        TRIGGER: Python idioms & best practices (3.12-3.14, async/asyncio, pattern matching, error handling); pytest/Hypothesis testing;
        static type checking (mypy, pyright, ty, pyrefly); the uv toolchain & packaging (pyproject, hatch); pydantic v2; CPython runtime internals & performance profiling;
        Python-in-browser/WASM. SKIP: JS/TS/Node/Deno/Bun → lang-js-ts; Go/Kotlin → lang-go-and-mobile; AI/ML frameworks → ai-* hubs.`),
        spokes: [['python-patterns','ref'],['python-testing','ref'],['python-static-type-checking','ref'],['python-supply-chain-security','ref'],['uv-python-toolchain','ref'],['pydantic-v2','ref'],['cpython-runtime-internals','ref'],['cpython-performance-profiling','ref'],['python-in-browser-wasm','ref']] },
      { id: 'lang-js-ts', desc: D(`JavaScript / TypeScript language sub-hub (lang family).
        TRIGGER: JavaScript & Node.js idioms/APIs; TypeScript (expert + advanced types); non-Node runtimes (Deno 2, Bun 1, edge/WinterTC);
        Node concurrency internals, native addons/N-API, TS+runtime features; V8 engine internals; zod schema validation; JS/Node/HTML/CSS debugging.
        SKIP: Python → lang-python; Go/Kotlin → lang-go-and-mobile; backend framework patterns (Express/Fastify/Nest/Hono) → software-engineering-patterns.`),
        spokes: [['javascript-nodejs','ref'],['javascript-runtimes-deno-bun-edge','ref'],['javascript-node-html-css-debugging-expert','ref'],['nodejs-concurrency-internals','ref'],['nodejs-native-addons-napi','ref'],['nodejs-typescript-and-runtime-features','ref'],['typescript-expert','ref'],['typescript-advanced-types','ref'],['zod-schema-validation','ref'],['v8-engine-internals','ref']] },
      { id: 'lang-go-and-mobile', desc: D(`Go & mobile/JVM language sub-hub (lang family).
        TRIGGER: Go idioms & patterns (goroutines, channels, error handling, Cobra CLIs, OTEL); Kotlin / Compose Multiplatform patterns.
        SKIP: Python → lang-python; JavaScript/TypeScript → lang-js-ts.`),
        spokes: [['go-patterns','ref'],['compose-multiplatform-patterns','ref']] },
    ],
  },
  {
    parent: 'ai-agent-engineering', family: 'ai-agent', manifest: 'ai-agent-manifest.json',
    router: D(`AI & agent-engineering family ROUTER. Split into: ai-agents-orchestration (agent frameworks, multi-agent, memory, planning, guardrails, coding/GUI agents, autonomous loops, eval);
      ai-rag-retrieval (RAG, iterative retrieval, vector/graph datastores); ai-llm-model-layer (training, fine-tuning, alignment/RLHF, compression, inference serving, transformer/multimodal architecture, model selection, observability);
      ai-mcp-sdk-prompting (MCP servers/builder, Anthropic SDK, prompt engineering, context engineering, LLM frameworks, tool-search, prompt lookup). Route to the matching sub-hub.`),
    subhubs: [
      { id: 'ai-agents-orchestration', desc: D(`AI agent building & orchestration sub-hub (ai-agent family).
        TRIGGER: agent ecosystems/frameworks (Claude SDK, LangGraph, CrewAI, OpenAI Agents); multi-agent orchestration & councils; A2A interop / Agent Cards;
        agent harness construction (action spaces, tool defs, observation format); agent memory architecture; planning patterns; reliability & guardrails;
        durable execution/state; agent workflow builders; autonomous loops; coding agents; computer-use/GUI agents; voice/realtime agents; eval-driven development.
        SKIP: RAG/retrieval → ai-rag-retrieval; model training/serving → ai-llm-model-layer; MCP/SDK/prompting → ai-mcp-sdk-prompting.`),
        spokes: [['a2a-interop','ref'],['agent-council','ref'],['agent-ecosystem','ref'],['agent-harness-construction','ref'],['agent-memory-architecture','ref'],['agent-planning-patterns','ref'],['agent-reliability-and-guardrails','ref'],['agent-state-and-durable-execution','ref'],['agent-workflow-builder_ai_toolkit','ref'],['autonomous-loops','ref'],['coding-agents','ref'],['computer-use-and-gui-agents','ref'],['multi-agent-orchestration','ref'],['voice-realtime-agents','ref'],['eval-driven-development','ref']] },
      { id: 'ai-rag-retrieval', desc: D(`RAG & retrieval sub-hub (ai-agent family).
        TRIGGER: Retrieval-Augmented Generation architecture (chunking, embeddings, hybrid search, reranking, query transformation, evaluation);
        advanced RAG patterns (self-RAG, corrective, GraphRAG, agentic, adaptive); iterative/multi-step retrieval; AI datastores (vector & graph databases).
        SKIP: agent orchestration → ai-agents-orchestration; model training/serving → ai-llm-model-layer; MCP/SDK/prompting → ai-mcp-sdk-prompting.`),
        spokes: [['advanced-rag-patterns','ref'],['rag-architecture','ref'],['iterative-retrieval','ref'],['ai-datastores','ref']] },
      { id: 'ai-llm-model-layer', desc: D(`LLM model-layer sub-hub (ai-agent family) — training, alignment, serving, architecture.
        TRIGGER: pretraining & scaling laws; distributed training; fine-tuning/PEFT; alignment & post-training (SFT, RLHF/PPO, DPO family); RLHF infrastructure; agentic RL;
        LLM compression (quantization/distillation/pruning); GPU kernels; inference serving/runtime (PagedAttention, batching, speculative decoding, TTFT/TPOT);
        transformer & multimodal architecture; reasoning models; model landscape/selection; LLM observability; continuous learning.
        SKIP: agent orchestration → ai-agents-orchestration; RAG → ai-rag-retrieval; MCP/SDK/prompting → ai-mcp-sdk-prompting.`),
        spokes: [['agentic-rl','ref'],['distributed-training','ref'],['llm-alignment-post-training','ref'],['llm-compression','ref'],['llm-fine-tuning-peft','ref'],['llm-gpu-kernels','ref'],['llm-inference-serving','ref'],['llm-pretraining-scaling-laws','ref'],['llm-routing-cascades','ref'],['rlhf-infrastructure','ref'],['transformer-architecture','ref'],['multimodal-llm-architecture','ref'],['reasoning-models','ref'],['llm-models','ref'],['llm-observability','ref'],['continuous-learning-v2','ref']] },
      { id: 'ai-mcp-sdk-prompting', desc: D(`MCP, SDK, prompting & context-engineering sub-hub (ai-agent family).
        TRIGGER: MCP servers & the MCP builder; Anthropic SDK; prompt engineering; LLM context engineering (CLAUDE.md, caching, agent memory, prompt-injection hardening);
        declarative LLM frameworks; multi-provider LLM integration review; AI programming languages; AI red-teaming tooling; MCP tool-search/discovery optimization; prompt discovery/lookup & registries.
        SKIP: agent orchestration → ai-agents-orchestration; RAG → ai-rag-retrieval; model training/serving → ai-llm-model-layer.`),
        spokes: [['mcp-builder','ref'],['mcp-servers','ref'],['anthropic-sdk','ref'],['prompt-engineering','ref'],['ai-languages','ref'],['declarative-llm-frameworks','ref'],['llm-integration-reviewer','ref'],['ai-redteaming-tooling','ref'],['llm-context-engineering','ref'],['mcp-tool-search-optimizer','standalone'],['prompt-lookup','standalone']] },
    ],
  },
];

const newHubs = [
  { hub: 'aws-cloud', family: 'aws-cloud', manifest: 'aws-cloud-manifest.json',
    desc: D(`Amazon Web Services hub — core infra, serverless/event-driven, AI/ML services, and the AWS database ecosystem.
      TRIGGER: IAM (policies, roles, SCPs, federation), AWS CLI/SSO, boto3/aws-sdk, EC2/EBS/S3, VPC design, Well-Architected, CloudWatch/CloudTrail (aws-core);
      Lambda, EventBridge, Step Functions, API Gateway, ECS/EKS Fargate, CDK/SAM/CloudFormation, serverless patterns (aws-serverless);
      Bedrock & SageMaker — foundation models, agents, RAG, guardrails, training/endpoints (aws-ai-ml); AWS databases — RDS/Aurora/DynamoDB/DocumentDB/Neptune, DocumentDB-vs-Atlas, CockroachDB, IndexedDB (databases-aws-cockroach-indexeddb).
      SKIP: non-AWS cloud; MongoDB Atlas platform → mongodb-atlas-expert.`),
    spokes: ['aws-core', 'aws-serverless', 'aws-ai-ml', 'databases-aws-cockroach-indexeddb'] },
  { hub: 'security-review', family: 'security-review', manifest: 'security-review-manifest.json',
    desc: D(`Application & infrastructure security-review hub.
      TRIGGER: security review of web apps / Chrome extensions / backend services (OWASP Top 10, ASVS, CSP/CORS, secrets, session/auth) (security-reviewer);
      codebase security & compliance audit — secrets/credentials, supply-chain, PII/logging, AI-policy, repo governance/SSDLC gates (security-compliance-auditor);
      HTTP security headers — Helmet.js, Content-Security-Policy, HSTS, COOP/COEP, X-Frame-Options (http-security-headers);
      Okta identity platform — Identity Engine, OAuth/OIDC, auth servers, management APIs, security posture (okta-expert).
      SKIP: Web Crypto/vault code review → webcrypto-vault-reviewer; designing new auth flows → web-auth-patterns.`),
    spokes: ['security-reviewer', 'security-compliance-auditor', 'http-security-headers', 'okta-expert'] },
];

// folds into EXISTING hubs (append to existing manifests)
const folds = [
  { hub: 'software-engineering-patterns', family: 'software', manifest: 'software-manifest.json', spokes: ['testing-and-vitest-expert', 'webapp-testing', 'repo-pattern-scanner'] },
  { hub: 'tam-operations', family: 'tam-operations', manifest: 'tam-operations-manifest.json', spokes: ['monday-board-audit', 'slack-subscription-auditor'] },
  { hub: 'content-ingestion-extraction', family: 'content-ingestion-extraction', manifest: 'content-ingestion-extraction-manifest.json', spokes: ['using-plaud-mcp'] },
  { hub: 'mongodb-operations-expert', family: 'mongodb', manifest: 'mongodb-manifest.json', spokes: ['database-migrations'] },
];

// new single-hub families (fold one spoke, create manifest)
const newFamilyFolds = [
  { hub: 'claude-code-skills', family: 'claude-code', manifest: 'claude-code-manifest.json', spokes: ['claude-code-plugins'] },
  { hub: 'deep-research', family: 'deep-research', manifest: 'deep-research-manifest.json', spokes: ['deep-research-methods'] },
];

// ---- EXECUTE ----------------------------------------------------------------------
const newManifestFiles = [];

// SPLITS
for (const sp of splits) {
  const siblingIds = sp.subhubs.map(s => '`' + s.id + '`').join(', ');
  const manifestHubs = {};
  for (const sub of sp.subhubs) {
    const destDir = path.join(SK, sub.id);
    const entries = sub.spokes.map(([spoke, type]) => {
      const src = type === 'ref'
        ? { type: 'ref', path: path.join(SK, sp.parent, 'references', `${spoke}.md`) }
        : { type: 'standalone', dir: path.join(SK, spoke) };
      return placeSpoke(spoke, destDir, sub.id, siblingIds, src);
    });
    writeFile(path.join(destDir, 'SKILL.md'), skillMd(sub.id, sub.desc));
    manifestHubs[sub.id] = { keepExisting: false, title: sub.id, spokes: entries };
  }
  // parent -> thin router; drop emptied references dir
  const parentMd = path.join(SK, sp.parent, 'SKILL.md');
  const cur = fs.readFileSync(parentMd, 'utf8');
  const nameM = cur.match(/^name:\s*(.*)$/m);
  writeFile(parentMd, skillMd(sp.parent, sp.router).replace(`name: ${sp.parent}`, `name: ${nameM ? nameM[1] : sp.parent}`));
  rmIfEmpty(path.join(SK, sp.parent, 'references'));
  // write family manifest
  const mf = { family: sp.family, builtAt: new Date().toISOString(), hubs: manifestHubs };
  writeFile(path.join(CR, sp.manifest), JSON.stringify(mf, null, 2) + '\n');
  newManifestFiles.push(sp.manifest);
}

// NEW HUBS
for (const nh of newHubs) {
  const destDir = path.join(SK, nh.hub);
  const entries = nh.spokes.map(spoke => placeSpoke(spoke, destDir, nh.hub, '`' + nh.hub + '`', { type: 'standalone', dir: path.join(SK, spoke) }));
  writeFile(path.join(destDir, 'SKILL.md'), skillMd(nh.hub, nh.desc));
  const mf = { family: nh.family, builtAt: new Date().toISOString(), hubs: { [nh.hub]: { keepExisting: false, title: nh.hub, spokes: entries } } };
  writeFile(path.join(CR, nh.manifest), JSON.stringify(mf, null, 2) + '\n');
  newManifestFiles.push(nh.manifest);
}

// FOLDS into existing manifests
function appendToManifest(manifestFile, hub, entries) {
  const p = path.join(CR, manifestFile);
  const m = JSON.parse(fs.readFileSync(p, 'utf8'));
  if (!m.hubs[hub]) m.hubs[hub] = { keepExisting: true, title: hub, spokes: [] };
  const have = new Set(m.hubs[hub].spokes.map(s => s.spoke));
  for (const e of entries) if (!have.has(e.spoke)) m.hubs[hub].spokes.push(e);
  actions.push(['manifest-append', `${manifestFile}:${hub} += ${entries.map(e => e.spoke).join(',')}`]);
  if (APPLY) fs.writeFileSync(p, JSON.stringify(m, null, 2) + '\n');
}
for (const f of folds) {
  const destDir = path.join(SK, f.hub);
  const entries = f.spokes.map(spoke => placeSpoke(spoke, destDir, f.hub, '`' + f.hub + '`', { type: 'standalone', dir: path.join(SK, spoke) }));
  appendToManifest(f.manifest, f.hub, entries);
}

// NEW single-hub family folds (create manifest)
for (const f of newFamilyFolds) {
  const destDir = path.join(SK, f.hub);
  if (APPLY) fs.mkdirSync(path.join(destDir, 'references'), { recursive: true });
  const entries = f.spokes.map(spoke => placeSpoke(spoke, destDir, f.hub, '`' + f.hub + '`', { type: 'standalone', dir: path.join(SK, spoke) }));
  const mf = { family: f.family, builtAt: new Date().toISOString(), hubs: { [f.hub]: { keepExisting: true, title: f.hub, spokes: entries } } };
  writeFile(path.join(CR, f.manifest), JSON.stringify(mf, null, 2) + '\n');
  newManifestFiles.push(f.manifest);
}

// Remove programming-languages hub block from software-manifest.json (its spokes moved to lang-*)
{
  const p = path.join(CR, 'software-manifest.json');
  const m = JSON.parse(fs.readFileSync(p, 'utf8'));
  if (m.hubs['programming-languages']) { delete m.hubs['programming-languages']; actions.push(['manifest-edit', 'software-manifest.json: removed programming-languages hub block']); if (APPLY) fs.writeFileSync(p, JSON.stringify(m, null, 2) + '\n'); }
}

// tier-config: add new manifests
{
  const p = path.join(CR, 'tiering', 'tier-config.json');
  const cfg = JSON.parse(fs.readFileSync(p, 'utf8'));
  for (const mf of newManifestFiles) if (!cfg.manifests.includes(mf)) cfg.manifests.push(mf);
  actions.push(['tier-config', 'manifests += ' + newManifestFiles.join(', ')]);
  if (APPLY) fs.writeFileSync(p, JSON.stringify(cfg, null, 2) + '\n');
}

log(`\n[restructure]${APPLY ? ' APPLIED' : ' DRY-RUN'} — ${actions.length} actions`);
const counts = {};
for (const [k] of actions) counts[k] = (counts[k] || 0) + 1;
log('  by type:', JSON.stringify(counts));
log('  new manifests:', newManifestFiles.join(', '));
if (APPLY) {
  try { const out = execFileSync(process.execPath, [path.join(CR, 'referents.mjs'), '--repair', '--apply', '--quiet'], { encoding: 'utf8' }); log('  referent repair:', out.trim() || 'ok'); }
  catch (e) { console.error('  referent repair FAILED:', e.message); }
}
