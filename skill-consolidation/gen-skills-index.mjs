#!/usr/bin/env node
// gen-skills-index.mjs — consolidate the WHOLE skill library into one agent-ingestible index.
//
// Why this exists: the per-family `*-manifest.json` files only store the pre-TRIGGER
// `routingLine`, are split per family, and don't carry the structured TRIGGER/SKIP routing,
// peers, version, or usage-rank that an agent needs to pick ONE skill without reading 117
// SKILL.md files. This builds the single cross-family index that was missing. It does NOT
// replace any existing tooling — it READS the existing sources of truth and joins them:
//   • ~/.claude/skills/<id>/SKILL.md frontmatter (name, description, related_skills, version…)
//   • hub-registry.mjs              → spoke↔hub↔family (the canonical manifest parser)
//   • <hub>/references/<spoke>.md    → folded (cold-tier) spokes that have no top-level dir
//   • tiering/access-log.jsonl       → usage signal for the selection rank
//
// Outputs (both in this dir, so every workspace under $HOME can read them):
//   SKILLS-INDEX.json  — full machine index   |   SKILLS-INDEX.md — compact agent table
//
// Usage:  node gen-skills-index.mjs                 # write both outputs
//         node gen-skills-index.mjs --check          # exit 1 if SKILLS-INDEX.json is stale (CI)
//         node gen-skills-index.mjs --stdout         # print JSON, don't write
//         node gen-skills-index.mjs --quiet
//         node gen-skills-index.mjs --embed          # (re)build SKILLS-EMBEDDINGS.json (needs a local embedding server)
//         node gen-skills-index.mjs --search "<q>"   # hybrid (cosine + keyword) kNN over the embeddings
//                                                     #   [--alpha=0.92 | --no-hybrid (pure cosine) | --top=N --threshold=T --json]
//
// Zero npm dependencies. The OPTIONAL semantic layer (--embed/--search) talks to a local
// embedding server (Ollama, OpenAI-compatible /api/embed) over the built-in fetch; the
// default run and --check stay fully offline and byte-deterministic.

import fs from "node:fs";
import fsp from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import crypto from "node:crypto";
import { fileURLToPath } from "node:url";
import { loadHubRegistry } from "./hub-registry.mjs";

const HERE = path.dirname(fileURLToPath(import.meta.url));
const HOME = os.homedir();
// SKILLS_ROOT and the output dir accept env overrides so tests/benchmarks can run the whole
// pipeline against a throwaway library + corpus, fully isolated from the live artifacts.
const SKILLS_ROOT = process.env.SKILLS_ROOT_DIR || path.join(HOME, ".claude", "skills");
const OUT_DIR = process.env.SKILLS_OUT_DIR || HERE;
const OUT_JSON = path.join(OUT_DIR, "SKILLS-INDEX.json");
const OUT_MD = path.join(OUT_DIR, "SKILLS-INDEX.md");
const ACCESS_LOG = path.join(HERE, "tiering", "access-log.jsonl");
const RANK_WINDOW_DAYS = 30;
const OUT_EMB = path.join(OUT_DIR, "SKILLS-EMBEDDINGS.json");

// ── Semantic layer config (opt-in). Pinned model+dimension: cosine similarity is only
// valid within ONE model+dim, so --search enforces query model/dim == corpus model/dim. ──
const OLLAMA_HOST = (process.env.OLLAMA_HOST || "http://localhost:11434").replace(/\/+$/, "");
const EMBED_MODEL = process.env.SKILLS_EMBED_MODEL || "qwen3-embedding:4b";
const EMBED_KEY = process.env.OLLAMA_API_KEY || process.env.OPENAI_API_KEY || ""; // optional; never logged
const EMBED_CONCURRENCY = Number(process.env.SKILLS_EMBED_CONCURRENCY) || 8;
const SEARCH_TOP_N = 10;
const SEARCH_THRESHOLD = 0.30;
// Hybrid (default) blends cosine with keyword overlap: score = alpha*cosine + (1-alpha)*coverage,
// where coverage = (# content query tokens present in the skill's name/summary/triggers) / (# query tokens).
// alpha=0.92 was the A/B optimum (bench-search.mjs): hybrid-4b P@1 70% / P@3 87% / MRR 0.775,
// beating pure cosine-4b (67% / 80% / 0.744) on every metric. --no-hybrid falls back to pure cosine.
const SEARCH_ALPHA = 0.92;

const ARGV = process.argv.slice(2);
const FLAGS = new Set(ARGV);
// "--name value" OR "--name=value"; undefined if absent, "" if present without a value.
function argVal(name) {
  for (let i = 0; i < ARGV.length; i++) {
    if (ARGV[i] === name) return ARGV[i + 1] ?? "";
    if (ARGV[i].startsWith(name + "=")) return ARGV[i].slice(name.length + 1);
  }
  return undefined;
}
const QUIET = FLAGS.has("--quiet");
const log = (...a) => { if (!QUIET) console.error(...a); };
const tilde = (p) => p.startsWith(HOME) ? "~" + p.slice(HOME.length) : p;

// ── Frontmatter parsing (mirrors build.mjs so the two stay byte-compatible) ──
function frontmatter(md) {
  const mm = md.match(/^---\n([\s\S]*?)\n---/);
  return mm ? mm[1] : "";
}
function field(fm, name) {
  const re = new RegExp(`^${name}:\\s*(?:[|>][-+]?)?\\s*(.*)$`, "m");
  const mm = fm.match(re);
  if (!mm) return "";
  const val = mm[1].trim();
  if (val) return val.replace(/^["']|["']$/g, "");
  const lines = fm.split("\n");
  const idx = lines.findIndex((l) => new RegExp(`^${name}:`).test(l));
  const buf = [];
  for (let i = idx + 1; i < lines.length; i++) {
    if (/^\s+\S/.test(lines[i])) buf.push(lines[i].trim());
    else break;
  }
  return buf.join(" ").replace(/^["']|["']$/g, "");
}
// YAML block list:  name:\n  - a\n  - b
function listField(fm, name) {
  const lines = fm.split("\n");
  const idx = lines.findIndex((l) => new RegExp(`^${name}:\\s*$`).test(l));
  if (idx === -1) return [];
  const out = [];
  for (let i = idx + 1; i < lines.length; i++) {
    const m = lines[i].match(/^\s+-\s+(.*)$/);
    if (!m) break;
    out.push(m[1].trim().replace(/^["']|["']$/g, ""));
  }
  return out;
}

// ── Description splitter: summary / TRIGGER list / SKIP routing edges ─────────
function splitDesc(desc) {
  const d = (desc || "").replace(/\s+/g, " ").trim();
  const tIdx = d.search(/\bTRIGGER:/);
  const sIdx = d.search(/\bSKIP:/);
  const cut = [tIdx, sIdx].filter((x) => x >= 0);
  const summary = (cut.length ? d.slice(0, Math.min(...cut)) : d).trim().replace(/[.\s]+$/, "");
  const triggerBlock = tIdx >= 0 ? d.slice(tIdx + 8, sIdx > tIdx ? sIdx : undefined) : "";
  const skipBlock = sIdx >= 0 ? d.slice(sIdx + 5) : "";
  const splitClauses = (s) => s.split(/;|·/).map((x) => x.trim()).filter(Boolean);
  const triggers = splitClauses(triggerBlock);
  const skip = splitClauses(skipBlock);
  // SKIP clauses often end "... → target-skill" / "-> target" / "use target": harvest peer ids
  const skipTo = [];
  for (const clause of skip) {
    const m = clause.match(/(?:→|->|use)\s+([a-z0-9][a-z0-9-]+[a-z0-9])\b/i);
    if (m) skipTo.push(m[1].toLowerCase());
  }
  return { summary, triggers, skip, skipTo };
}

// ── Usage rank from the tiering access-log (recent accesses = selection weight) ──
function loadRank() {
  const rank = new Map();
  let raw;
  try { raw = fs.readFileSync(ACCESS_LOG, "utf8"); } catch { return rank; }
  const cutoff = Date.now() - RANK_WINDOW_DAYS * 86400000;
  for (const line of raw.split("\n")) {
    if (!line.trim()) continue;
    let e; try { e = JSON.parse(line); } catch { continue; }
    const t = Date.parse(e.ts);
    if (!e.spoke || !(t >= cutoff)) continue;
    rank.set(e.spoke, (rank.get(e.spoke) || 0) + 1);
  }
  return rank;
}

// ── Resolve every skill id → its SKILL.md source path ────────────────────────
// Top-level skills live at <root>/<id>/SKILL.md (symlinks into .agents/skills are
// transparently followed). Folded spokes with no top-level dir are read from the
// verbatim copy at <root>/<hub>/references/<spoke>.md that build.mjs maintains.
function resolveSources(reg) {
  const src = new Map(); // id -> { path, kind }
  let dirents = [];
  try { dirents = fs.readdirSync(SKILLS_ROOT, { withFileTypes: true }); } catch {}
  for (const d of dirents) {
    if (d.name.startsWith(".")) continue;
    const p = path.join(SKILLS_ROOT, d.name, "SKILL.md");
    if (fs.existsSync(p)) src.set(d.name, { path: p, kind: "top-level" });
  }
  for (const spoke of reg.spokes) {
    if (src.has(spoke)) continue;
    const hub = reg.spokeHub.get(spoke);
    const p = path.join(SKILLS_ROOT, hub, "references", `${spoke}.md`);
    if (fs.existsSync(p)) src.set(spoke, { path: p, kind: "folded-spoke" });
  }
  return src;
}

// Parallel read pool — this IS the cache-warm: one bounded fan-out read of every
// SKILL.md in the main process, so a high-latency / cold-cache FS pays once.
async function readAll(paths, concurrency = 32) {
  const out = new Map();
  let i = 0;
  await Promise.all(Array.from({ length: concurrency }, async () => {
    while (i < paths.length) {
      const p = paths[i++];
      try { out.set(p, await fsp.readFile(p, "utf8")); } catch { out.set(p, null); }
    }
  }));
  return out;
}

// ── Semantic layer (opt-in): embedding client, source text, hashing, cosine. The
// default run and --check never reach this code — they stay offline + deterministic. ──
function embedSource(s) {
  return [s.name, s.summary, (s.triggers || []).join("; ")].filter(Boolean).join("\n");
}
function srcHash(s) {
  return crypto.createHash("sha256").update(EMBED_MODEL + "\n" + embedSource(s)).digest("hex").slice(0, 16);
}
function cosine(a, b) {
  let dot = 0, na = 0, nb = 0;
  const n = Math.min(a.length, b.length);
  for (let i = 0; i < n; i++) { dot += a[i] * b[i]; na += a[i] * a[i]; nb += b[i] * b[i]; }
  return na && nb ? dot / (Math.sqrt(na) * Math.sqrt(nb)) : 0;
}
// Keyword-overlap helpers (used by --hybrid and the benchmark). Mirror the regex router's notion
// of a literal token match: lowercase alnum runs of length ≥3, minus a small stop list.
const STOPWORDS = new Set("the and for with from this that into your you are how can what when use using does keeps come back over out a an of to in on it my our i".split(" "));
const tokens = (s) => new Set((s || "").toLowerCase().match(/[a-z0-9]{3,}/g) || []);
function kwCoverage(queryToks, bag) {
  if (!queryToks.length) return 0;
  let m = 0; for (const t of queryToks) if (bag.has(t)) m++;
  return m / queryToks.length;
}
async function embedText(text, model = EMBED_MODEL) {
  const headers = { "content-type": "application/json" };
  if (EMBED_KEY) headers.authorization = `Bearer ${EMBED_KEY}`; // secret: never logged/printed
  const res = await fetch(`${OLLAMA_HOST}/api/embed`, {
    method: "POST", headers,
    body: JSON.stringify({ model, input: text }),
  });
  if (!res.ok) throw new Error(`embed HTTP ${res.status}`);
  const j = await res.json();
  const v = Array.isArray(j.embeddings) ? j.embeddings[0] : j.embedding;
  if (!Array.isArray(v) || !v.length) throw new Error("no embedding vector in response");
  return v;
}

async function build() {
  const reg = loadHubRegistry();
  const rank = loadRank();
  const src = resolveSources(reg);
  const contents = await readAll([...src.values()].map((s) => s.path));

  const skills = [];
  for (const [id, info] of src) {
    const md = contents.get(info.path);
    if (md == null) continue;
    const fm = frontmatter(md);
    const { summary, triggers, skip, skipTo } = splitDesc(field(fm, "description"));
    const related = listField(fm, "related_skills");
    const hub = reg.spokeHub.get(id) || null;     // the hub this skill is a SPOKE of
    const isHub = reg.hubs.has(id);               // this skill is itself a router/hub
    // Family: a spoke inherits its hub's family; a hub belongs to its own family;
    // a true standalone has none. This groups each hub WITH its spokes for selection.
    const family = hub ? famOf(reg, hub) : (isHub ? famOf(reg, id) : null);
    skills.push({
      id,
      name: field(fm, "name") || id,
      kind: info.kind,
      isHub,
      hub,
      family,
      version: field(fm, "version") || null,
      updated: field(fm, "updated") || null,
      origin: field(fm, "origin") || null,
      summary,
      triggers,
      skip,
      // Peer set for discovery: declared related_skills ∪ SKIP-routed targets ∪ owning hub.
      peers: [...new Set([...related, ...skipTo, ...(hub ? [hub] : [])])].filter((x) => x !== id),
      rankAccess30d: rank.get(id) || 0,
      srcBytes: Buffer.byteLength(md),
      path: tilde(info.path),
    });
  }
  // Selection order: most-used first, then alphabetical for stable diffs.
  skills.sort((a, b) => b.rankAccess30d - a.rankAccess30d || a.id.localeCompare(b.id));

  const byFamily = {};
  for (const s of skills) (byFamily[s.family || "(standalone)"] ||= []).push(s.id);

  return {
    schemaVersion: 1,
    skillsRoot: tilde(SKILLS_ROOT),
    counts: {
      total: skills.length,
      topLevel: skills.filter((s) => s.kind === "top-level").length,
      foldedSpokes: skills.filter((s) => s.kind === "folded-spoke").length,
      families: Object.keys(byFamily).length,
    },
    families: Object.fromEntries(Object.entries(byFamily).sort()),
    skills,
  };
}

// Family for a hub: hub-registry exposes families as a set, not a hub→family map,
// so re-derive by scanning the manifests once and caching.
let _famCache = null;
function famOf(reg, hub) {
  if (!_famCache) {
    _famCache = new Map();
    try {
      for (const f of fs.readdirSync(HERE).filter((x) => x.endsWith("-manifest.json") && x !== "consolidation-manifest.json")) {
        const m = JSON.parse(fs.readFileSync(path.join(HERE, f), "utf8"));
        if (m && m.hubs) for (const h of Object.keys(m.hubs)) _famCache.set(h, m.family || null);
      }
    } catch {}
  }
  return _famCache.get(hub) || null;
}

function renderMd(idx) {
  const out = [];
  out.push("<!-- AUTO-GENERATED by gen-skills-index.mjs — do not edit by hand. -->");
  out.push("# Skill Library Index");
  out.push("");
  out.push(`${idx.counts.total} skills (${idx.counts.topLevel} top-level + ${idx.counts.foldedSpokes} folded spokes) across ${idx.counts.families} families. Skills root: \`${idx.skillsRoot}\`.`);
  out.push("");
  out.push("**Agent usage:** scan the table to pick the single best-matching skill by its triggers, then READ that skill's `SKILL.md` (or its `references/` file) before answering — this table is for routing, not depth. Rows are ordered by recent-use rank within each family.");
  out.push("");
  const fams = Object.keys(idx.families).sort();
  const byId = new Map(idx.skills.map((s) => [s.id, s]));
  for (const fam of fams) {
    out.push(`## ${fam}`);
    out.push("");
    out.push("| skill | v | summary | top triggers | peers |");
    out.push("| --- | --- | --- | --- | --- |");
    for (const id of idx.families[fam]) {
      const s = byId.get(id);
      const cell = (x) => (x || "").replace(/\|/g, "\\|").replace(/\n/g, " ");
      const trg = cell(s.triggers.slice(0, 4).join("; ")) || "—";
      const peers = s.peers.length ? s.peers.slice(0, 6).map((p) => `\`${p}\``).join(", ") : "—";
      const tag = s.isHub ? " *(hub)*" : s.kind === "folded-spoke" ? " *(folded)*" : "";
      out.push(`| \`${s.id}\`${tag} | ${s.version || "—"} | ${cell(s.summary).slice(0, 140) || "—"} | ${trg} | ${peers} |`);
    }
    out.push("");
  }
  return out.join("\n");
}

// ── --embed : (re)generate SKILLS-EMBEDDINGS.json incrementally, keyed by source hash ──
async function embed(idx) {
  try { await embedText("preflight"); }
  catch (e) {
    console.error(`embedding server unreachable at ${OLLAMA_HOST} (${e.message}). ` +
      `Start it (e.g. 'ollama serve') or set OLLAMA_HOST. Leaving ${path.basename(OUT_EMB)} untouched.`);
    process.exit(2);
  }
  let prev = {};
  try { prev = JSON.parse(fs.readFileSync(OUT_EMB, "utf8")); } catch {}
  const reusable = prev.model === EMBED_MODEL ? (prev.vectors || {}) : {}; // model change ⇒ re-embed all
  const items = idx.skills.map((s) => ({ id: s.id, text: embedSource(s), hash: srcHash(s) }));
  const vectors = {};
  let reused = 0, embedded = 0, failed = 0, dim = null, i = 0;
  await Promise.all(Array.from({ length: EMBED_CONCURRENCY }, async () => {
    while (i < items.length) {
      const it = items[i++];
      const old = reusable[it.id];
      if (old && old.hash === it.hash && Array.isArray(old.embedding)) {
        vectors[it.id] = old; reused++; dim ||= old.embedding.length; continue;
      }
      try {
        const e = await embedText(it.text);
        vectors[it.id] = { hash: it.hash, embedding: e }; embedded++; dim ||= e.length;
      } catch (err) {
        failed++; log(`embed FAIL ${it.id}: ${err.message}`);
        if (old && Array.isArray(old.embedding)) vectors[it.id] = old; // keep prior, don't lose corpus
      }
    }
  }));
  const ids = Object.keys(vectors);
  // No-churn idempotency: if nothing was (re)embedded, the model matches, and the id set is
  // unchanged, the corpus is byte-identical except for generatedAt — skip the write so the
  // periodic/WatchPaths automation has zero side effects when the library hasn't changed.
  const sameSet = prev.model === EMBED_MODEL && Object.keys(reusable).length === ids.length && ids.every((id) => reusable[id]);
  if (embedded === 0 && failed === 0 && sameSet && fs.existsSync(OUT_EMB)) {
    log(`${tilde(OUT_EMB)} already up to date — ${ids.length} vectors (model ${EMBED_MODEL}, dim ${prev.dim}); embedded 0, reused ${reused}.`);
    return;
  }
  const payload = {
    schemaVersion: 1, model: EMBED_MODEL, dim: dim || (prev.dim ?? null),
    generatedAt: new Date().toISOString(), vectors,
  };
  const tmp = OUT_EMB + ".tmp";
  fs.writeFileSync(tmp, JSON.stringify(payload) + "\n");
  fs.renameSync(tmp, OUT_EMB); // atomic: a concurrent WatchPaths embed + manual --embed never corrupt
  log(`Wrote ${tilde(OUT_EMB)} — ${ids.length} vectors ` +
    `(model ${EMBED_MODEL}, dim ${payload.dim}); embedded ${embedded}, reused ${reused}, failed ${failed}.`);
  if (failed && !embedded && !reused) process.exit(2);
}

// ── --search "<query>" : hybrid (cosine + keyword) kNN over the corpus. Fail-OPEN: if the
// embedding server is unreachable or the corpus is missing/empty, it degrades to keyword-only
// ranking over the always-present SKILLS-INDEX.json (exit 0) rather than erroring. The only
// non-zero exits are an empty query / dim mismatch (2) and a genuine no-match (1). ──
async function search(query) {
  const q = (query || "").trim();
  if (!q) { console.error('usage: node gen-skills-index.mjs --search "<query>" [--alpha=A] [--no-hybrid] [--top=N] [--threshold=T] [--json]'); process.exit(2); }
  const top = Number(argVal("--top")) || SEARCH_TOP_N;
  const tRaw = argVal("--threshold");
  const threshold = tRaw === undefined || tRaw === "" ? SEARCH_THRESHOLD : Number(tRaw);
  const hybrid = !FLAGS.has("--no-hybrid"); // hybrid is the default; --no-hybrid forces pure cosine
  const aRaw = argVal("--alpha");
  const alpha = aRaw === undefined || aRaw === "" ? SEARCH_ALPHA : Number(aRaw);
  const wantJson = FLAGS.has("--json");
  // Keyword bags from the always-present offline index. Built once, before any network call,
  // so they serve both the hybrid scorer below and the no-server keyword fallback.
  const qToks = [...tokens(q)].filter((t) => !STOPWORDS.has(t));
  const meta = new Map(); // id → { rank: rankAccess30d, summary, bag } (no extra embedding)
  try { for (const s of JSON.parse(fs.readFileSync(OUT_JSON, "utf8")).skills) meta.set(s.id, { rank: s.rankAccess30d || 0, summary: s.summary || "", bag: tokens([s.name, s.summary, (s.triggers || []).join(" ")].join(" ")) }); } catch {}

  // Degraded ranking when the semantic layer is unavailable: pure keyword coverage, no server
  // needed. Prints hits (exit 0) or, if nothing shares a keyword, points at the regex router
  // (exit 1). Never exits 2 — server-down is an expected, recoverable condition.
  const keywordFallback = (reason) => {
    console.error(`${reason} — degraded to keyword-only ranking (no semantic similarity).`);
    const scored = [...meta.entries()]
      .map(([id, m]) => ({ id, cov: kwCoverage(qToks, m.bag || new Set()), rank: m.rank || 0, summary: m.summary || "" }))
      .filter((x) => x.cov > 0)
      .sort((a, b) => b.cov - a.cov || b.rank - a.rank || a.id.localeCompare(b.id));
    const hits = scored.slice(0, top);
    if (!hits.length) {
      console.error(`no skill shares a keyword with "${q}". Fall back to the TRIGGER/SKIP routing in SKILLS-INDEX.md.`);
      process.exit(1);
    }
    if (wantJson) {
      process.stdout.write(JSON.stringify({
        query: q, mode: "keyword", degraded: true, top,
        results: hits.map(({ id, cov, rank }) => ({ id, score: Number(cov.toFixed(4)), coverage: Number(cov.toFixed(4)), rankAccess30d: rank })),
      }, null, 2) + "\n");
    } else {
      for (const h of hits) process.stdout.write(`${h.cov.toFixed(3)}  ${h.id}  —  ${(h.summary || "").slice(0, 90)}\n`);
    }
  };

  let corpus;
  try { corpus = JSON.parse(fs.readFileSync(OUT_EMB, "utf8")); }
  catch { return keywordFallback(`No ${path.basename(OUT_EMB)} (run: node gen-skills-index.mjs --embed to enable semantic search)`); }
  const entries = Object.entries(corpus.vectors || {}).filter(([, r]) => Array.isArray(r.embedding));
  if (!entries.length) return keywordFallback(`${path.basename(OUT_EMB)} has no vectors`);
  let qv;
  try { qv = await embedText(q); }
  catch (e) { return keywordFallback(`embedding server unreachable at ${OLLAMA_HOST} (${e.message})`); }
  if (corpus.dim && qv.length !== corpus.dim) {
    console.error(`dim mismatch: query ${qv.length} vs corpus ${corpus.dim} (model ${corpus.model}). Re-embed with the same model.`);
    process.exit(2);
  }
  const scored = entries.map(([id, r]) => {
    const m = meta.get(id) || {};
    const cos = cosine(qv, r.embedding);
    const cov = hybrid ? kwCoverage(qToks, m.bag || new Set()) : 0;
    return { id, cos, cov, score: hybrid ? alpha * cos + (1 - alpha) * cov : cos, rank: m.rank || 0, summary: m.summary || "" };
  }).sort((a, b) => b.score - a.score || b.rank - a.rank || a.id.localeCompare(b.id));
  const hits = scored.filter((x) => x.score >= threshold).slice(0, top);
  if (!hits.length) {
    console.error(`no skill scored ≥ ${threshold} for "${q}" (best ${scored[0] ? scored[0].score.toFixed(3) : "n/a"}). ` +
      `Lower --threshold or fall back to regex routing.`);
    process.exit(1);
  }
  if (wantJson) {
    process.stdout.write(JSON.stringify({
      query: q, model: corpus.model, dim: corpus.dim, threshold, top, mode: hybrid ? "hybrid" : "cosine", ...(hybrid ? { alpha } : {}),
      results: hits.map(({ id, score, cos, cov, rank }) => ({ id, score: Number(score.toFixed(4)), ...(hybrid ? { cosine: Number(cos.toFixed(4)), coverage: Number(cov.toFixed(4)) } : {}), rankAccess30d: rank })),
    }, null, 2) + "\n");
  } else {
    for (const h of hits) process.stdout.write(`${h.score.toFixed(3)}  ${h.id}  —  ${(h.summary || "").slice(0, 90)}\n`);
  }
}

async function main() {
  const searchQuery = argVal("--search");
  if (searchQuery !== undefined) { await search(searchQuery); return; }
  const idx = await build();
  if (FLAGS.has("--embed")) { await embed(idx); return; }
  if (FLAGS.has("--stdout")) {
    process.stdout.write(JSON.stringify(idx, null, 2) + "\n");
    return;
  }
  if (FLAGS.has("--check")) {
    let prev = null;
    try { prev = JSON.parse(fs.readFileSync(OUT_JSON, "utf8")); } catch {}
    const norm = (o) => o && { ...o, generatedAt: undefined };
    const stale = JSON.stringify(norm(prev)) !== JSON.stringify(norm({ ...idx, generatedAt: undefined }));
    if (stale) { console.error("SKILLS-INDEX.json is STALE — run: node gen-skills-index.mjs"); process.exit(1); }
    log("SKILLS-INDEX.json is up to date.");
    return;
  }
  const stamped = { generatedAt: new Date().toISOString(), ...idx };
  fs.writeFileSync(OUT_JSON, JSON.stringify(stamped, null, 2) + "\n");
  fs.writeFileSync(OUT_MD, renderMd(idx) + "\n");
  log(`Wrote ${tilde(OUT_JSON)} and ${tilde(OUT_MD)} — ${idx.counts.total} skills, ${idx.counts.families} families.`);
}

// Auto-run only as the process entrypoint; importing this module (tests, benchmarks) is
// side-effect-free and exposes the routines below.
const IS_ENTRY = !!process.argv[1] && path.resolve(process.argv[1]) === fileURLToPath(import.meta.url);
if (IS_ENTRY) main().catch((e) => { console.error(e); process.exit(2); });

export { build, embed, search, embedText, embedSource, srcHash, cosine, splitDesc, renderMd, tokens, STOPWORDS, kwCoverage, EMBED_MODEL, OLLAMA_HOST };
