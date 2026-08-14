#!/usr/bin/env node
// bench-search.mjs — A/B retrieval methods over a labeled query set: a KEYWORD baseline (the
// literal token overlap a regex TRIGGER/SKIP router relies on), SEMANTIC cosine over each
// embedding corpus (0.6b live + optional larger 4b), and HYBRID (alpha*cosine + (1-alpha)*keyword
// coverage) over each. Reports P@1, P@3, MRR per method + mean query-embedding latency. Read-only:
// loads the corpora + offline index, never writes. Queries are deliberately paraphrased to measure
// the recall the regex layer structurally cannot get.
//   node bench-search.mjs                 # table + per-query first-hit ranks
//   node bench-search.mjs --json
//   node bench-search.mjs --alpha=0.92    # hybrid blend weight (default 0.92, the production value)
// The primary corpus is the live production SKILLS-EMBEDDINGS.json (now 4b); an optional smaller-model
// comparison corpus (0.6b) is read from SKILLS_06_CORPUS — when absent, the 0.6b rows are skipped.
// Override paths via env: SKILLS_4B_CORPUS (primary), SKILLS_06_CORPUS (comparison).
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { build, embedText, cosine, tokens, STOPWORDS, kwCoverage, OLLAMA_HOST } from "./gen-skills-index.mjs";

const HERE = path.dirname(fileURLToPath(import.meta.url));
const CORPUS_4B = process.env.SKILLS_4B_CORPUS || path.join(HERE, "SKILLS-EMBEDDINGS.json"); // production (now 4b)
const CORPUS_06 = process.env.SKILLS_06_CORPUS || "/tmp/skills-0.6b/SKILLS-EMBEDDINGS.json"; // optional comparison

// query → acceptable expected skill ids (a hit if any expected id is in the top-k).
const LABELS = [
  { q: "why is my atlas cluster pegged at 100% cpu on a read-heavy workload", expect: ["atlas-diagnostics-expert", "mongodb-atlas-expert", "deep-mongodb-mql-query-optimizer"] },
  { q: "how do I make a react widget usable for screen reader users", expect: ["interface-kit", "accessibility-ux-reviewer", "react-nextjs"] },
  { q: "stream row changes out of a database into kafka", expect: ["terraform-kafka-infra", "telemetry-pipeline", "mongodb-operations-expert"] },
  { q: "help me put together a CV for a cocktail bar job", expect: ["service-industry-resume-and-interview"] },
  { q: "i got pulled over and charged with drunk driving in north carolina", expect: ["nc-criminal-defense"] },
  { q: "tighten up this system prompt so the agent behaves consistently", expect: ["prompt-deep-optimizer", "prompt-helper-optimizer"] },
  { q: "this select statement does a full table scan and is slow", expect: ["deep-query-optimizer"] },
  { q: "critique the visual hierarchy of my landing page mockup", expect: ["design-deep-optimizer"] },
  { q: "start a charity with tax exempt status in NC", expect: ["venture-nonprofit-cause", "venture-nc-nonprofit-formation"] },
  { q: "track down a memory leak in my node server", expect: ["javascript-node-html-css-debugging-expert", "javascript-nodejs"] },
  { q: "chunking and reranking documents for retrieval augmented generation", expect: ["ai-rag-retrieval"] },
  { q: "build a model context protocol server", expect: ["ai-mcp-sdk-prompting"] },
  { q: "github actions that builds a docker image and ships it to kubernetes", expect: ["devops-containers-cicd"] },
  { q: "my dns lookups intermittently come back servfail", expect: ["networking"] },
  { q: "visitors get a 520 error from my site behind the orange cloud", expect: ["cloudflare-platform"] },
  { q: "explain taproot and how lightning channels settle", expect: ["bitcoin-protocol-expert"] },
  { q: "the borrow checker keeps rejecting my code over lifetimes", expect: ["lang-rust"] },
  { q: "structured concurrency with task groups in python", expect: ["lang-python", "python-concurrency"] },
  { q: "our identity provider is down, run the sev1 playbook", expect: ["okta-incident-response"] },
  { q: "prepare a quarterly business review for a key account", expect: ["tam-operations"] },
  { q: "should I pick elasticsearch or opensearch after the license change", expect: ["elasticsearch-opensearch"] },
  { q: "my SPL search times out over a wide time range", expect: ["splunk-platform-spl"] },
  { q: "dispute a charge-off on my credit report", expect: ["consumer-credit-and-debt", "charge-offs-collections-and-debt-resolution"] },
  { q: "why do first year donors stop giving", expect: ["fundraising-and-donor-psychology"] },
  { q: "use self determination theory to drive product adoption", expect: ["applied-psychology"] },
  { q: "audit my web app for owasp vulnerabilities", expect: ["security-review"] },
  { q: "proofread and fact check this runbook before publishing", expect: ["document-critique"] },
  { q: "review and fix bugs across this whole repository", expect: ["code-deep-optimizer"] },
  { q: "call the jira rest api from my app", expect: ["integration-clients", "jira-extension-client"] },
  { q: "pick a regression model and evaluate it", expect: ["da-analytical-methods"] },
];

const rankOf = (ranked, expect) => { for (let i = 0; i < ranked.length; i++) if (expect.includes(ranked[i])) return i + 1; return Infinity; };
const argVal = (name, def) => {
  const i = process.argv.indexOf(name);
  if (i >= 0 && process.argv[i + 1]) return process.argv[i + 1];
  const p = process.argv.find((a) => a.startsWith(name + "="));
  return p ? p.slice(name.length + 1) : def;
};
const loadCorpus = (p) => {
  try { const c = JSON.parse(fs.readFileSync(p, "utf8")); return { model: c.model, dim: c.dim, vecs: Object.entries(c.vectors || {}).filter(([, r]) => Array.isArray(r.embedding)) }; }
  catch { return null; }
};

async function main() {
  const asJson = process.argv.includes("--json");
  const alpha = Number(argVal("--alpha", "0.92"));
  const idx = await build();
  const byId = new Map(idx.skills.map((s) => [s.id, s]));
  const bags = new Map(idx.skills.map((s) => [s.id, tokens([s.name, s.summary, (s.triggers || []).join(" ")].join(" "))]));

  const c4b = loadCorpus(CORPUS_4B); // primary (production)
  const c06 = loadCorpus(CORPUS_06); // optional comparison
  if (!c4b) { console.error(`missing primary corpus at ${CORPUS_4B} — run: node gen-skills-index.mjs --embed`); process.exit(2); }

  const unknown = [...new Set(LABELS.flatMap((l) => l.expect).filter((id) => !byId.has(id)))];
  if (unknown.length) console.error("WARN unknown expected ids:", unknown.join(", "));

  const tieRank = (id) => byId.get(id)?.rankAccess30d || 0;
  const sortIds = (arr) => arr.sort((a, b) => b.s - a.s || tieRank(b.id) - tieRank(a.id) || a.id.localeCompare(b.id)).map((x) => x.id);
  const cosRank = (qv, c) => sortIds(c.vecs.map(([id, r]) => ({ id, s: cosine(qv, r.embedding) })));
  const hybRank = (qv, c, qt) => sortIds(c.vecs.map(([id, r]) => ({ id, s: alpha * cosine(qv, r.embedding) + (1 - alpha) * kwCoverage(qt, bags.get(id) || new Set()) })));

  const mk = () => ({ p1: 0, p3: 0, mrr: 0 });
  const M = { keyword: mk(), "sem-4b": mk(), "hybrid-4b": mk() };
  if (c06) { M["sem-0.6b"] = mk(); M["hybrid-0.6b"] = mk(); }
  const acc = (m, ranked, expect) => { const r = rankOf(ranked, expect); if (r === 1) m.p1++; if (r <= 3) m.p3++; m.mrr += r === Infinity ? 0 : 1 / r; return r; };
  const lat = { "4b": 0, "0.6b": 0 };
  const rows = [];

  for (const { q, expect } of LABELS) {
    const qt = [...tokens(q)].filter((t) => !STOPWORDS.has(t));
    const kwRanked = sortIds(idx.skills.map((s) => { const b = bags.get(s.id); let c = 0; for (const t of qt) if (b.has(t)) c++; return { id: s.id, s: c }; }));

    let t0 = performance.now();
    const qv4b = await embedText(q, c4b.model);
    lat["4b"] += performance.now() - t0;

    let qv06 = null;
    if (c06) { t0 = performance.now(); qv06 = await embedText(q, c06.model); lat["0.6b"] += performance.now() - t0; }

    const row = { q, expect };
    row.kw = acc(M.keyword, kwRanked, expect);
    row.s4b = acc(M["sem-4b"], cosRank(qv4b, c4b), expect);
    row.h4b = acc(M["hybrid-4b"], hybRank(qv4b, c4b, qt), expect);
    if (c06) {
      row.s06 = acc(M["sem-0.6b"], cosRank(qv06, c06), expect);
      row.h06 = acc(M["hybrid-0.6b"], hybRank(qv06, c06, qt), expect);
    }
    rows.push(row);
  }

  const n = LABELS.length, pct = (x) => ((x / n) * 100).toFixed(0) + "%", f3 = (x) => (x / n).toFixed(3);
  const ms4b = (lat["4b"] / n).toFixed(1), ms06 = c06 ? (lat["0.6b"] / n).toFixed(1) : null;
  const summary = {
    queries: n, alpha,
    models: { "4b": `${c4b.model} (${c4b.dim}-d)`, "0.6b": c06 ? `${c06.model} (${c06.dim}-d)` : null },
    meanQueryMs: { "4b": Number(ms4b), "0.6b": ms06 ? Number(ms06) : null },
    methods: Object.fromEntries(Object.entries(M).map(([k, m]) => [k, { "P@1": pct(m.p1), "P@3": pct(m.p3), MRR: f3(m.mrr) }])),
  };
  if (asJson) { process.stdout.write(JSON.stringify({ summary, rows }, null, 2) + "\n"); return; }

  const latOf = { keyword: "~0 ms (no embedding)", "sem-4b": `${ms4b} ms`, "hybrid-4b": `${ms4b} ms (+kw)`, "sem-0.6b": `${ms06} ms`, "hybrid-0.6b": `${ms06} ms (+kw)` };
  console.log(`\nRetrieval A/B — ${n} labeled queries · alpha ${alpha}`);
  console.log(`  primary 4b = ${c4b.model} (${c4b.dim}-d)   comparison 0.6b = ${c06 ? `${c06.model} (${c06.dim}-d)` : "(absent — set SKILLS_06_CORPUS)"}\n`);
  console.log("Method        P@1    P@3    MRR     mean-query-latency");
  for (const [k, m] of Object.entries(M)) console.log(`${k.padEnd(13)} ${pct(m.p1).padEnd(6)} ${pct(m.p3).padEnd(6)} ${f3(m.mrr).padEnd(7)} ${latOf[k]}`);

  const r = (x) => (x === undefined ? " " : x === Infinity ? "—" : x);
  const cols = c06 ? ["kw", "s4b", "h4b", "s06", "h06"] : ["kw", "s4b", "h4b"];
  console.log(`\nPer-query first-hit rank (lower = better, — = miss):\n`);
  console.log("  " + cols.map((c) => c.padStart(3)).join(" ") + "  query");
  for (const x of rows) {
    const cells = c06 ? [x.kw, x.s4b, x.h4b, x.s06, x.h06] : [x.kw, x.s4b, x.h4b];
    console.log("  " + cells.map((c) => String(r(c)).padStart(3)).join(" ") + "  " + x.q.slice(0, 54));
  }
  console.log("");
}
main().catch((e) => { console.error(e); process.exit(2); });
