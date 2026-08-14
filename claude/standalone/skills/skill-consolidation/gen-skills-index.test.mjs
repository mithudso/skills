// gen-skills-index.test.mjs — proves the semantic layer works.
//   node --test gen-skills-index.test.mjs
// Offline unit tests always run. Live integration tests (the embedding server) self-skip when
// the server is down, so the suite is safe in restricted CI. Live tests that mutate a corpus use
// a throwaway SKILLS_ROOT_DIR + SKILLS_OUT_DIR, so they never touch the real artifacts or race
// the WatchPaths agent.
import { test } from "node:test";
import assert from "node:assert/strict";
import { execFile } from "node:child_process";
import { promisify } from "node:util";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import {
  cosine, embedSource, srcHash, splitDesc, build, embedText, tokens, STOPWORDS, kwCoverage, OLLAMA_HOST,
} from "./gen-skills-index.mjs";

const HERE = path.dirname(fileURLToPath(import.meta.url));
const SCRIPT = path.join(HERE, "gen-skills-index.mjs");
const run = promisify(execFile);
const cli = (args, env) =>
  run(process.execPath, [SCRIPT, ...args], { cwd: HERE, env: env || process.env })
    .then((r) => ({ code: 0, ...r }))
    .catch((e) => ({ code: e.code ?? 1, stdout: e.stdout ?? "", stderr: e.stderr ?? "" }));

const serverUp = await fetch(`${OLLAMA_HOST}/api/version`, { signal: AbortSignal.timeout(2000) })
  .then((r) => r.ok).catch(() => false);
if (!serverUp) console.error(`[skip] embedding server down at ${OLLAMA_HOST} — live tests skipped`);

// ── offline unit tests (no network) ─────────────────────────────────────────
test("cosine: identical=1, orthogonal=0, opposite=-1, degenerate=0", () => {
  assert.equal(cosine([1, 2, 3], [1, 2, 3]), 1);
  assert.equal(cosine([1, 0], [0, 1]), 0);
  assert.equal(cosine([1, 0], [-1, 0]), -1);
  assert.equal(cosine([], []), 0); // never NaN
});

test("embedSource joins name+summary+triggers and drops empties", () => {
  assert.equal(embedSource({ name: "n", summary: "s", triggers: ["a", "b"] }), "n\ns\na; b");
  assert.equal(embedSource({ name: "n", summary: "", triggers: [] }), "n");
});

test("srcHash is 16-hex, stable, and content-sensitive", () => {
  const a = srcHash({ name: "x", summary: "y", triggers: ["z"] });
  assert.match(a, /^[0-9a-f]{16}$/);
  assert.equal(a, srcHash({ name: "x", summary: "y", triggers: ["z"] }));
  assert.notEqual(a, srcHash({ name: "x", summary: "y2", triggers: ["z"] }));
});

test("splitDesc parses summary/TRIGGER/SKIP and harvests skipTo peers", () => {
  const r = splitDesc("Does X. TRIGGER: foo; bar baz SKIP: when Y → other-skill");
  assert.equal(r.summary, "Does X");
  assert.deepEqual(r.triggers, ["foo", "bar baz"]);
  assert.ok(r.skipTo.includes("other-skill"));
});

test("tokens: lowercases, keeps alnum runs ≥3, dedupes", () => {
  const t = tokens("Atlas CPU on a Read-heavy WORKLOAD a an");
  assert.ok(t.has("atlas") && t.has("cpu") && t.has("read") && t.has("heavy") && t.has("workload"));
  assert.ok(!t.has("on") && !t.has("a") && !t.has("an")); // <3 chars dropped
  assert.equal(tokens("Cpu cpu CPU").size, 1); // dedup, case-insensitive
  assert.equal(tokens("").size, 0);
});

test("kwCoverage: fraction of query tokens present in the bag; 0 on empty query", () => {
  const bag = tokens("rust borrow checker lifetimes ownership");
  assert.equal(kwCoverage(["borrow", "checker"], bag), 1);
  assert.equal(kwCoverage(["borrow", "missing"], bag), 0.5);
  assert.equal(kwCoverage([], bag), 0); // never NaN
  assert.ok(STOPWORDS.has("the") && !STOPWORDS.has("borrow"));
});

test("build() returns a well-formed, populated index (offline, deterministic)", async () => {
  const idx = await build();
  assert.equal(idx.schemaVersion, 1);
  assert.ok(idx.counts.total >= 100, `expected many skills, got ${idx.counts.total}`);
  for (const s of idx.skills.slice(0, 5)) {
    assert.ok(typeof s.id === "string" && s.id.length > 0);
    assert.ok(Array.isArray(s.triggers));
  }
});

// ── live integration: read-only against the real corpus ──────────────────────
test("embedText returns a vector matching the corpus dim", { skip: !serverUp }, async () => {
  const v = await embedText("hello world");
  assert.ok(Array.isArray(v) && v.length > 0);
  const corpus = JSON.parse(fs.readFileSync(path.join(HERE, "SKILLS-EMBEDDINGS.json"), "utf8"));
  assert.equal(v.length, corpus.dim);
});

test("--search ranks a MongoDB query onto a mongodb/atlas skill (exit 0)", { skip: !serverUp }, async () => {
  const { code, stdout } = await cli(["--search", "atlas cluster slow query high CPU", "--top=5", "--json"]);
  assert.equal(code, 0);
  const out = JSON.parse(stdout);
  assert.ok(out.results.length > 0);
  assert.match(out.results[0].id, /mongodb|atlas|mql/);
});

test("--search empty query exits 2", { skip: !serverUp }, async () => {
  assert.equal((await cli(["--search", ""])).code, 2);
});

test("--search is hybrid by default (mode+alpha+per-hit cosine/coverage); --no-hybrid is pure cosine", { skip: !serverUp }, async () => {
  const def = JSON.parse((await cli(["--search", "the borrow checker rejects my lifetimes", "--top=3", "--json"])).stdout);
  assert.equal(def.mode, "hybrid");
  assert.equal(typeof def.alpha, "number");
  assert.ok("cosine" in def.results[0] && "coverage" in def.results[0]);

  const pure = JSON.parse((await cli(["--search", "the borrow checker rejects my lifetimes", "--top=3", "--no-hybrid", "--json"])).stdout);
  assert.equal(pure.mode, "cosine");
  assert.ok(!("alpha" in pure));
  assert.ok(!("coverage" in pure.results[0]));
});

test("--search impossible threshold exits 1 with a fall-back message", { skip: !serverUp }, async () => {
  const { code, stderr } = await cli(["--search", "atlas slow query", "--threshold=0.999"]);
  assert.equal(code, 1);
  assert.match(stderr, /regex routing|threshold/i);
});

// ── degraded path: a forced-unreachable server must fail OPEN, not error. Runs in any state
// (incl. restricted CI) because it points OLLAMA_HOST at a dead port; the keyword fallback
// reads only the offline SKILLS-INDEX.json, so it returns results with exit 0. ──
test("--search fails open to keyword-only when the embedding server is unreachable (exit 0)", async () => {
  const env = { ...process.env, OLLAMA_HOST: "http://127.0.0.1:1" };
  const { code, stdout, stderr } = await cli(["--search", "atlas slow query high cpu", "--top=5", "--json"], env);
  assert.equal(code, 0);
  assert.match(stderr, /degraded to keyword-only/i);
  const out = JSON.parse(stdout);
  assert.equal(out.mode, "keyword");
  assert.equal(out.degraded, true);
  assert.ok(out.results.length > 0, "keyword fallback should return at least one hit");
  assert.ok("coverage" in out.results[0]);
});

// ── live integration: isolated corpus (a new skill embeds, reuse, no-churn) ──
test("incremental --embed: new skill adds a vector, then reuse, then no-churn", { skip: !serverUp }, async () => {
  const tmpRoot = fs.mkdtempSync(path.join(os.tmpdir(), "skroot-"));
  const tmpOut = fs.mkdtempSync(path.join(os.tmpdir(), "skout-"));
  const env = { ...process.env, SKILLS_ROOT_DIR: tmpRoot, SKILLS_OUT_DIR: tmpOut };
  const mk = (id, desc) => {
    fs.mkdirSync(path.join(tmpRoot, id));
    fs.writeFileSync(path.join(tmpRoot, id, "SKILL.md"), `---\nname: ${id}\ndescription: ${desc}\n---\n`);
  };
  try {
    mk("alpha-skill", "Alpha does indexing. TRIGGER: alpha index test");
    let r = await cli(["--embed"], env);
    assert.equal(r.code, 0);
    assert.match(r.stderr, /embedded 1, reused 0/);
    const corpus = JSON.parse(fs.readFileSync(path.join(tmpOut, "SKILLS-EMBEDDINGS.json"), "utf8"));
    assert.ok(corpus.vectors["alpha-skill"], "new skill should have a vector");
    assert.equal(corpus.vectors["alpha-skill"].embedding.length, corpus.dim);

    mk("beta-skill", "Beta does search. TRIGGER: beta search test");
    r = await cli(["--embed"], env);
    assert.match(r.stderr, /embedded 1, reused 1/, "adding one skill re-embeds only the new one");

    r = await cli(["--embed"], env);
    assert.match(r.stderr, /already up to date/, "no change ⇒ no-churn (no rewrite)");
  } finally {
    fs.rmSync(tmpRoot, { recursive: true, force: true });
    fs.rmSync(tmpOut, { recursive: true, force: true });
  }
});
