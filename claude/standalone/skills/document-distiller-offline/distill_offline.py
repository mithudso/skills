#!/usr/bin/env python3
"""
distill_offline.py — the offline engine for /distill-offline.

This is the deterministic, zero-LLM-token half of the document-distiller
pipeline. The LLM does exactly ONE thing the LLM can't be replaced for:
read the doc and emit a compact JSON array of classified, deduped,
salience-scored knowledge units. Everything mechanical around that — rendering
markdown from the JSON, writing the central + source-adjacent files, computing
the diff/add scope, and an honest degraded heuristic extract — lives here and
costs no tokens.

Why this saves tokens vs the original /distill (see docs in the skill):
  - The original has the model emit BOTH json and markdown. The markdown is
    fully derivable from the json, so ~half the output tokens are spent twice.
    Here the model emits json only; `render` produces the markdown.
  - All file I/O, section ordering, salience sorting, dedup counting, the
    `## Removed` section, and source-adjacent copies are done here, not by the
    model.
  - `diff` scopes --diff/--add so the model never re-reads unchanged text.

Subcommands:
  render   compact units JSON  -> central .md/.json + source-adjacent copies
  diff     old + new files     -> added/removed hunks JSON (scopes the LLM pass)
  merge    existing + new units-> deduped-against-existing units JSON (re-emit)
  fetch    URL                 -> plain text (offline ingest; no MCP round-trip)
  extract  plain text          -> HEURISTIC candidate units (degraded, honest)

All modes are stdlib-only. Read the skill for how the LLM pass hands off here.
"""

from __future__ import annotations

import argparse
import datetime as _dt
import difflib
import html
import json
import math
import os
import re
import subprocess
import sys
import urllib.error
import urllib.request
from html.parser import HTMLParser
from pathlib import Path

DIST_DIR = Path.home() / ".claude" / "distillations"
INDEX_DIR = DIST_DIR / ".index"

# Ollama is used ONLY as an offline semantic layer (dedup + novelty filtering).
# It never reads or generates prose — embeddings only — so it adds no LLM tokens.
OLLAMA_URL = os.environ.get("OLLAMA_HOST", "http://127.0.0.1:11434").rstrip("/")
DEFAULT_EMBED_MODEL = os.environ.get("DISTILL_EMBED_MODEL", "mxbai-embed-large")

# Taxonomy order drives section order in the markdown. Keep in sync with SKILL.md.
TYPE_ORDER = [
    "concept",
    "fact",
    "actionable",
    "question",
    "problem",
    "statement",
    "quote",
    "idea",
]
TYPE_HEADINGS = {
    "concept": "Concepts",
    "fact": "Facts",
    "actionable": "Actionables",
    "question": "Questions",
    "problem": "Problems",
    "statement": "Statements",
    "quote": "Quotes",
    "idea": "Ideas",
}
SALIENCE_RANK = {"high": 0, "medium": 1, "low": 2}


def _today() -> str:
    return _dt.date.today().isoformat()


def _slugify(text: str) -> str:
    text = re.sub(r"[^\w\s-]", "", (text or "").lower()).strip()
    text = re.sub(r"[\s_-]+", "-", text)
    return text[:80] or "distillation"


# --------------------------------------------------------------------------- #
# render — the emit phase (the main token win: model never writes markdown)
# --------------------------------------------------------------------------- #

def _unit_salience(u: dict) -> int:
    return SALIENCE_RANK.get((u.get("salience") or "medium").lower(), 1)


def _render_unit_line(u: dict) -> str:
    t = u["type"]
    text = u.get("text", "").strip()
    anchor = u.get("source_anchor", "")
    sal = (u.get("salience") or "").lower()
    anchor_md = f"_{anchor}_" if anchor else ""

    if t == "quote":
        attr = u.get("attribution")
        q = f'> "{text}"'
        if attr:
            q += f" — {attr}"
        if anchor:
            q += f", {anchor_md}"
        if sal == "high":
            q += " · salience: high"
        return q

    if t == "concept":
        # Concepts: bold the leading term if the text is "Term — gloss".
        head, sep, rest = text.partition(" — ")
        body = f"**{head}** — {rest}" if sep else f"**{text}**"
        line = f"- {body}"
    elif t == "actionable":
        line = f"- [ ] {text}"
    else:
        line = f"- {text}"

    if anchor:
        line = f"{line} — {anchor_md}"
    # Salience shown for concepts (matches the original's convention) and any
    # explicitly-high unit, to keep the list scannable; omitted otherwise.
    if sal and (t == "concept" or sal == "high"):
        line += f" · salience: {sal}"
    return line


def render_markdown(data: dict) -> str:
    src = data.get("source", {})
    units = data.get("units", [])
    generated = data.get("generated") or _today()

    active = [u for u in units if u.get("status") != "removed"]
    removed = [u for u in units if u.get("status") == "removed"]
    after_dedup = len(active)
    extracted = after_dedup + sum(len(u.get("duplicates", []) or []) for u in units)

    title = src.get("title") or "Untitled document"
    ref = src.get("ref") or src.get("url") or src.get("path") or "(in-context document)"

    lines = [f"# Distilled: {title}", ""]
    lines.append(f"- Source: `{ref}`")
    lines.append(f"- Distilled: {generated}")
    lines.append(f"- Units: {extracted} extracted ({after_dedup} after dedup)")
    if data.get("anchors_note"):
        lines.append(f"- Anchors: {data['anchors_note']}")
    lines.append("")

    by_type: dict[str, list[dict]] = {t: [] for t in TYPE_ORDER}
    for u in active:
        t = u.get("type", "statement")
        if t not in by_type:
            t = "statement"  # unknown/typo'd type still renders, never dropped
        by_type[t].append(u)

    for t in TYPE_ORDER:
        group = by_type.get(t) or []
        if not group:
            continue
        group.sort(key=_unit_salience)
        lines.append(f"## {TYPE_HEADINGS[t]}")
        lines.append("")
        for u in group:
            lines.append(_render_unit_line(u))
        lines.append("")

    if removed:
        lines.append("## Removed")
        lines.append("")
        for u in removed:
            anchor = u.get("source_anchor", "")
            when = u.get("removed_on") or generated
            anchor_md = f" — _{anchor}_" if anchor else ""
            lines.append(f"- {u.get('text','').strip()}{anchor_md} · removed {when}")
        lines.append("")

    followups = data.get("followups") or []
    if followups:
        lines.append("## Suggested follow-ups (/dr)")
        lines.append("")
        for f in followups:
            concept = f.get("concept", "").strip()
            why = f.get("why", "").strip()
            lines.append(f"- `/dr {concept}` — {why}" if why else f"- `/dr {concept}`")
        lines.append("")

    return "\n".join(lines).rstrip() + "\n"


def _write_pair(stem_dir: Path, slug: str, generated: str, data: dict,
                json_only: bool) -> list[str]:
    stem_dir.mkdir(parents=True, exist_ok=True)
    written = []
    json_path = stem_dir / f"{slug}-{generated}.json"
    json_path.write_text(json.dumps(data, indent=2, ensure_ascii=False) + "\n")
    written.append(str(json_path))
    if not json_only:
        md_path = stem_dir / f"{slug}-{generated}.md"
        md_path.write_text(render_markdown(data))
        written.append(str(md_path))
    return written


def cmd_render(args) -> int:
    data = _load_json_input(args.in_file)
    src = data.setdefault("source", {})
    generated = data.get("generated") or _today()
    data["generated"] = generated

    slug = args.slug or _slugify(src.get("title") or _basename_of(src))
    written = _write_pair(DIST_DIR, slug, generated, data, args.json_only)

    # Source-adjacent copy: only for file inputs, best-effort.
    adjacent = []
    if src.get("type") == "file" and src.get("ref"):
        srcp = Path(src["ref"]).expanduser()
        adj_slug = f"{srcp.stem}-distilled"
        try:
            adjacent = _write_pair(srcp.parent, adj_slug, generated, data,
                                    args.json_only)
        except OSError as e:
            print(f"WARN source-adjacent copy failed ({e}); central copy kept.",
                  file=sys.stderr)

    result = {
        "written": written,
        "source_adjacent": adjacent,
        "units_after_dedup": len([u for u in data.get("units", [])
                                  if u.get("status") != "removed"]),
    }
    print(json.dumps(result, indent=2))
    return 0


def _basename_of(src: dict) -> str:
    ref = src.get("ref") or src.get("url") or src.get("path") or ""
    if not ref:
        return "distillation"
    ref = ref.rstrip("/")
    return Path(ref).name or ref.split("/")[-1] or "distillation"


# --------------------------------------------------------------------------- #
# diff — scope --diff so the LLM never re-reads unchanged text
# --------------------------------------------------------------------------- #

def _read_text(path_or_url: str) -> str:
    p = Path(path_or_url).expanduser()
    if p.exists():
        return p.read_text(errors="replace")
    raise FileNotFoundError(path_or_url)


def _unified_hunks(old: str, new: str) -> dict:
    """Return added and removed line groups using difflib (git-independent)."""
    old_lines = old.splitlines()
    new_lines = new.splitlines()
    sm = difflib.SequenceMatcher(a=old_lines, b=new_lines, autojunk=False)
    added, removed = [], []
    for tag, i1, i2, j1, j2 in sm.get_opcodes():
        if tag in ("replace", "insert"):
            for j in range(j1, j2):
                added.append({"new_line": j + 1, "text": new_lines[j]})
        if tag in ("replace", "delete"):
            for i in range(i1, i2):
                removed.append({"old_line": i + 1, "text": old_lines[i]})
    return {"added": added, "removed": removed}


def cmd_diff(args) -> int:
    old = _read_text(args.old)
    new = _read_text(args.new)
    hunks = _unified_hunks(old, new)
    # Drop blank-only added/removed lines — no signal to classify.
    hunks["added"] = [h for h in hunks["added"] if h["text"].strip()]
    hunks["removed"] = [h for h in hunks["removed"] if h["text"].strip()]
    print(json.dumps(hunks, indent=2, ensure_ascii=False))
    return 0


# --------------------------------------------------------------------------- #
# merge — dedup NEW units against an EXISTING set only (for --add / --diff)
# --------------------------------------------------------------------------- #

def _norm(s: str) -> str:
    return re.sub(r"\s+", " ", (s or "").lower()).strip()


def _next_id(units: list[dict]) -> int:
    mx = 0
    for u in units:
        m = re.match(r"u(\d+)", u.get("id", ""))
        if m:
            mx = max(mx, int(m.group(1)))
    return mx + 1


def cmd_merge(args) -> int:
    existing = _load_json_input(args.existing)
    new_units = _load_json_input(args.new_units)
    if isinstance(new_units, dict):
        new_units = new_units.get("units", [])

    base = existing.get("units", [])
    # Existing units NOT marked removed this run are the dedup targets.
    targets = [u for u in base if u.get("status") != "removed"]
    target_norms = [(u, _norm(u.get("text", ""))) for u in targets]

    # Semantic dedup (opt-in): catches reworded restatements the lexical ratio
    # misses. Falls back to lexical if ollama is unavailable.
    sem_targets = sem_cands = None
    if args.semantic and targets:
        try:
            sem_targets = _ollama_embed([u.get("text", "") for u in targets], args.model)
            sem_cands = _ollama_embed([nu.get("text", "") for nu in new_units], args.model)
        except OllamaUnavailable as e:
            print(f"WARN {e}; semantic dedup off, using lexical.", file=sys.stderr)
            sem_targets = sem_cands = None

    nid = _next_id(base)
    folded = 0
    kept = 0
    today = _today()
    for ci, nu in enumerate(new_units):
        ntext = _norm(nu.get("text", ""))
        best, best_ratio = None, 0.0
        for ti, (u, tn) in enumerate(target_norms):
            if sem_targets is not None:
                r = _cosine(sem_cands[ci], sem_targets[ti])
            else:
                r = difflib.SequenceMatcher(a=ntext, b=tn, autojunk=False).ratio()
            if r > best_ratio:
                best, best_ratio = u, r
        if best is not None and best_ratio >= args.threshold:
            best.setdefault("duplicates", []).append(nu.get("id") or f"u{nid:03d}")
            folded += 1
        else:
            uid = f"u{nid:03d}"
            nid += 1
            nu["id"] = uid
            nu.setdefault("canonical", uid)
            nu.setdefault("added_in", today)
            base.append(nu)
            kept += 1

    existing["units"] = base
    existing["generated"] = today
    print(json.dumps(existing, indent=2, ensure_ascii=False))
    print(f"# merged: {kept} new, {folded} folded as duplicates "
          f"(threshold {args.threshold})", file=sys.stderr)
    return 0


# --------------------------------------------------------------------------- #
# fetch — offline URL -> text (avoids an MCP/WebFetch round-trip when possible)
# --------------------------------------------------------------------------- #

class _TextExtractor(HTMLParser):
    # Tags whose *contents* are markup/style/script, never page text — dropped whole.
    _SKIP = {"script", "style", "noscript", "head", "nav", "footer", "svg",
             "template", "iframe", "object", "embed", "canvas"}
    # Tags that imply a line/paragraph break in the extracted text.
    _BREAK = {"p", "br", "div", "li", "h1", "h2", "h3", "h4", "h5", "h6",
              "tr", "section", "article", "header", "blockquote", "pre"}

    def __init__(self):
        super().__init__(convert_charrefs=True)
        self.parts: list[str] = []
        self._skip_depth = 0

    def handle_starttag(self, tag, attrs):
        if tag in self._SKIP:
            self._skip_depth += 1
        elif tag in self._BREAK:
            self.parts.append("\n")

    def handle_startendtag(self, tag, attrs):
        if tag in self._BREAK:
            self.parts.append("\n")

    def handle_endtag(self, tag):
        if tag in self._SKIP and self._skip_depth:
            self._skip_depth -= 1

    def handle_data(self, data):
        if self._skip_depth == 0 and data.strip():
            self.parts.append(data)

    def text(self) -> str:
        raw = "".join(self.parts)
        raw = html.unescape(raw)
        raw = re.sub(r"[ \t]+", " ", raw)
        raw = re.sub(r"\n[ \t]+", "\n", raw)
        raw = re.sub(r"\n{3,}", "\n\n", raw)
        return raw.strip()


def _looks_like_html(text: str, path: str | None = None) -> bool:
    """Detect HTML by extension or by content — used to auto-strip before distill."""
    if path and path.lower().rsplit(".", 1)[-1] in ("html", "htm", "xhtml"):
        return True
    head = text[:4096].lower()
    return ("<!doctype html" in head or "<html" in head or "<body" in head
            or ("<div" in head and "</" in head))


def _html_to_text(body: str) -> str:
    """Strip all tags/CSS/JS and return only the page text."""
    parser = _TextExtractor()
    parser.feed(body)
    parser.close()
    return parser.text()


def cmd_fetch(args) -> int:
    try:
        out = subprocess.run(
            ["curl", "-fsSL", "--max-time", str(args.timeout), args.url],
            capture_output=True, text=True, check=True,
        )
    except (subprocess.CalledProcessError, FileNotFoundError) as e:
        print(f"ERROR fetch failed: {e}", file=sys.stderr)
        return 2
    body = out.stdout
    if _looks_like_html(body):
        body = _html_to_text(body)
    print(body)
    return 0


def cmd_clean(args) -> int:
    """Read a local file (or stdin) and, if it is HTML, discard every tag, CSS
    block, and script, printing only the page text. Non-HTML passes through
    unchanged so it is safe to run on any local file before distilling."""
    if args.in_file in (None, "-"):
        body = sys.stdin.read()
        path = None
    else:
        body = _read_text(args.in_file)
        path = args.in_file
    if args.force_html or _looks_like_html(body, path):
        body = _html_to_text(body)
    print(body)
    return 0


# --------------------------------------------------------------------------- #
# extract — HEURISTIC candidate units (DEGRADED; honest first-pass only)
# --------------------------------------------------------------------------- #

_ACTION_RE = re.compile(r"^\s*(?:-\s*\[\s?\]|\d+\.|[-*])?\s*"
                        r"(should|must|need to|todo|do |run |add |use |ensure|"
                        r"make sure|remember to|don'?t forget)", re.I)
_SENT_SPLIT = re.compile(r"(?<=[.!?])\s+(?=[A-Z0-9\"'])")


def _heuristic_type(sent: str) -> str:
    s = sent.strip()
    if s.startswith(('"', "“", ">")) or (s.count('"') >= 2):
        return "quote"
    if s.endswith("?"):
        return "question"
    if _ACTION_RE.match(s):
        return "actionable"
    if re.search(r"\b(risk|fails?|cannot|broken|limitation|bug|issue|problem)\b", s, re.I):
        return "problem"
    if re.search(r"\d", s) and re.search(r"\b(is|are|was|were|has|have|costs?|takes?)\b", s, re.I):
        return "fact"
    return "statement"


def cmd_extract(args) -> int:
    if args.in_file in (None, "-"):
        text = sys.stdin.read()
        path = None
    else:
        text = _read_text(args.in_file)
        path = args.in_file
    if _looks_like_html(text, path):
        text = _html_to_text(text)  # distill only page text, never markup
    units = []
    nid = 1
    for lineno, line in enumerate(text.splitlines(), start=1):
        line = line.strip()
        if not line or len(line) < 20:
            continue
        for sent in _SENT_SPLIT.split(line):
            sent = sent.strip()
            if len(sent) < 20:
                continue
            units.append({
                "id": f"u{nid:03d}",
                "type": _heuristic_type(sent),
                "text": sent,
                "source_anchor": f"L{lineno}",
                "salience": "medium",
                "canonical": f"u{nid:03d}",
            })
            nid += 1
    result = {
        "_warning": "HEURISTIC DEGRADED EXTRACT — not a real distillation. "
                    "No semantic dedup, no cross-doc reconciliation, coarse "
                    "types. Feed to the LLM classify pass or treat as a rough "
                    "first cut only.",
        "units": units,
    }
    print(json.dumps(result, indent=2, ensure_ascii=False))
    return 0


# --------------------------------------------------------------------------- #
# semantic layer (ollama embeddings) — offline, token-free
#
# Why this is here and what it does NOT do: distilling ONE fresh file still
# needs the whole doc in the LLM's context (every part is a candidate unit), so
# no index shrinks that first pass. The index pays off across runs: --add,
# --diff, and multi-file corpus distillation. Embed each distilled unit once;
# on later runs the `novelty` filter drops candidate material already covered by
# a prior distillation BEFORE it reaches the LLM, so the classify pass only sees
# genuinely new content. It also lets dedup catch reworded restatements that the
# lexical difflib path misses. All of this is embeddings-only: zero LLM tokens.
# --------------------------------------------------------------------------- #

class OllamaUnavailable(RuntimeError):
    pass


def _ollama_embed(texts: list[str], model: str) -> list[list[float]]:
    """Batch-embed via ollama /api/embed. Raises OllamaUnavailable on any error
    so callers can fall back to the lexical path instead of dying."""
    if not texts:
        return []
    payload = json.dumps({"model": model, "input": texts}).encode()
    req = urllib.request.Request(
        f"{OLLAMA_URL}/api/embed",
        data=payload,
        headers={"Content-Type": "application/json"},
    )
    try:
        with urllib.request.urlopen(req, timeout=120) as resp:
            data = json.loads(resp.read())
    except (urllib.error.URLError, TimeoutError, OSError, ValueError) as e:
        raise OllamaUnavailable(f"ollama embed failed ({e})") from e
    embs = data.get("embeddings")
    if not embs or len(embs) != len(texts):
        raise OllamaUnavailable("ollama returned no/mismatched embeddings")
    return embs


def _cosine(a: list[float], b: list[float]) -> float:
    dot = sum(x * y for x, y in zip(a, b))
    na = math.sqrt(sum(x * x for x in a))
    nb = math.sqrt(sum(y * y for y in b))
    return dot / (na * nb) if na and nb else 0.0


def _index_path(name: str) -> Path:
    return INDEX_DIR / f"{_slugify(name)}.json"


def cmd_index(args) -> int:
    """Embed the units of a distillation JSON and persist a vector index so
    future runs can dedup / novelty-filter against it. Idempotent per name."""
    data = _load_json_input(args.in_file)
    units = data.get("units", []) if isinstance(data, dict) else data
    units = [u for u in units if u.get("status") != "removed" and u.get("text")]
    src = data.get("source", {}) if isinstance(data, dict) else {}
    name = args.name or _slugify(src.get("title") or _basename_of(src))

    texts = [u["text"] for u in units]
    try:
        vectors = _ollama_embed(texts, args.model)
    except OllamaUnavailable as e:
        print(f"ERROR {e}; index needs ollama running.", file=sys.stderr)
        return 2

    entries = [
        {"id": u.get("id"), "text": u["text"], "type": u.get("type"),
         "source": src.get("ref"), "vector": v}
        for u, v in zip(units, vectors)
    ]
    INDEX_DIR.mkdir(parents=True, exist_ok=True)
    path = _index_path(name)
    path.write_text(json.dumps(
        {"model": args.model, "generated": _today(), "entries": entries},
        ensure_ascii=False))
    print(json.dumps({"index": str(path), "vectors": len(entries),
                      "model": args.model}, indent=2))
    return 0


def _load_indexes(names: list[str] | None) -> list[dict]:
    if names:
        paths = [_index_path(n) for n in names]
    else:
        paths = sorted(INDEX_DIR.glob("*.json")) if INDEX_DIR.exists() else []
    out = []
    for p in paths:
        if p.exists():
            out.extend(json.loads(p.read_text()).get("entries", []))
    return out


def cmd_novelty(args) -> int:
    """Given candidate units/lines and one or more prior indexes, emit only the
    candidates NOT already covered (cosine below --threshold to every indexed
    unit). This is the incremental/corpus token-saver: the LLM classify pass
    only ever sees novel material. Falls back to lexical-only if ollama is down."""
    cand = _load_json_input(args.in_file)
    if isinstance(cand, dict):
        cand = cand.get("units") or cand.get("added") or []
    cand_texts = [c.get("text", "") for c in cand]

    index = _load_indexes(args.against)
    if not index:
        # Nothing to compare against — everything is novel.
        print(json.dumps({"novel": cand, "matched": [],
                          "note": "no prior index; all candidates novel"},
                         ensure_ascii=False, indent=2))
        return 0

    index_texts = [e["text"] for e in index]
    semantic = True
    try:
        cand_vecs = _ollama_embed(cand_texts, args.model)
        idx_vecs = [e["vector"] for e in index]
    except OllamaUnavailable as e:
        print(f"WARN {e}; falling back to lexical novelty.", file=sys.stderr)
        semantic = False

    novel, matched = [], []
    for i, c in enumerate(cand):
        best_sim, best_j = 0.0, -1
        if semantic:
            for j, iv in enumerate(idx_vecs):
                s = _cosine(cand_vecs[i], iv)
                if s > best_sim:
                    best_sim, best_j = s, j
        else:
            for j, it in enumerate(index_texts):
                s = difflib.SequenceMatcher(a=_norm(cand_texts[i]), b=_norm(it),
                                            autojunk=False).ratio()
                if s > best_sim:
                    best_sim, best_j = s, j
        if best_sim >= args.threshold:
            matched.append({"candidate": c.get("text"), "similarity": round(best_sim, 3),
                            "matched_id": index[best_j].get("id"),
                            "matched_text": index[best_j].get("text")})
        else:
            novel.append(c)

    print(json.dumps({
        "mode": "semantic" if semantic else "lexical",
        "threshold": args.threshold,
        "novel_count": len(novel), "matched_count": len(matched),
        "novel": novel, "matched": matched,
    }, ensure_ascii=False, indent=2))
    return 0


# --------------------------------------------------------------------------- #
# plumbing
# --------------------------------------------------------------------------- #

def _load_json_input(path):
    if path in (None, "-"):
        return json.load(sys.stdin)
    return json.loads(Path(path).expanduser().read_text())


def main(argv=None) -> int:
    p = argparse.ArgumentParser(description="Offline engine for /distill-offline.")
    sub = p.add_subparsers(dest="cmd", required=True)

    r = sub.add_parser("render", help="compact units JSON -> md+json files")
    r.add_argument("--in", dest="in_file", default="-",
                   help="units JSON file, or - for stdin")
    r.add_argument("--slug", help="override output slug")
    r.add_argument("--json-only", action="store_true", help="skip markdown")
    r.set_defaults(func=cmd_render)

    d = sub.add_parser("diff", help="old+new files -> added/removed hunks JSON")
    d.add_argument("--old", required=True)
    d.add_argument("--new", required=True)
    d.set_defaults(func=cmd_diff)

    m = sub.add_parser("merge", help="dedup new units against existing set")
    m.add_argument("--existing", required=True, help="existing distillation JSON")
    m.add_argument("--new-units", dest="new_units", default="-",
                   help="new units JSON, or - for stdin")
    m.add_argument("--threshold", type=float, default=0.82,
                   help="near-duplicate ratio 0-1 (lexical ~0.82; "
                        "with --semantic try ~0.90)")
    m.add_argument("--semantic", action="store_true",
                   help="use ollama embeddings for dedup (catches rewordings)")
    m.add_argument("--model", default=DEFAULT_EMBED_MODEL,
                   help="ollama embedding model")
    m.set_defaults(func=cmd_merge)

    ix = sub.add_parser("index", help="embed a distillation's units -> vector index")
    ix.add_argument("--in", dest="in_file", default="-",
                    help="distillation JSON, or - for stdin")
    ix.add_argument("--name", help="index name (defaults to source slug)")
    ix.add_argument("--model", default=DEFAULT_EMBED_MODEL,
                    help="ollama embedding model")
    ix.set_defaults(func=cmd_index)

    nv = sub.add_parser("novelty",
                         help="drop candidates already covered by prior indexes")
    nv.add_argument("--in", dest="in_file", default="-",
                    help="candidate units/hunks JSON, or - for stdin")
    nv.add_argument("--against", nargs="*",
                    help="index names to compare against (default: all)")
    nv.add_argument("--threshold", type=float, default=0.90,
                    help="cosine/ratio above which a candidate is NOT novel")
    nv.add_argument("--model", default=DEFAULT_EMBED_MODEL,
                    help="ollama embedding model")
    nv.set_defaults(func=cmd_novelty)

    f = sub.add_parser("fetch", help="URL -> plain text (offline)")
    f.add_argument("--url", required=True)
    f.add_argument("--timeout", type=int, default=30)
    f.set_defaults(func=cmd_fetch)

    cl = sub.add_parser("clean",
                        help="local HTML file -> page text only (drop tags/CSS/JS)")
    cl.add_argument("--in", dest="in_file", default="-",
                    help="file to strip, or - for stdin")
    cl.add_argument("--force-html", action="store_true",
                    help="strip as HTML even if not auto-detected")
    cl.set_defaults(func=cmd_clean)

    e = sub.add_parser("extract", help="text -> HEURISTIC candidate units (degraded)")
    e.add_argument("--in", dest="in_file", default="-",
                   help="text file, or - for stdin")
    e.set_defaults(func=cmd_extract)

    args = p.parse_args(argv)
    try:
        return args.func(args)
    except (FileNotFoundError, json.JSONDecodeError) as e:
        print(f"ERROR {e}", file=sys.stderr)
        return 2


if __name__ == "__main__":
    raise SystemExit(main())
