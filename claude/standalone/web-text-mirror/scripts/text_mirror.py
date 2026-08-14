#!/usr/bin/env python3
"""Recursive text-only website mirror via trafilatura.

BFS-crawls one or more seed URLs (same-host scoping per seed), extracts each page
to Markdown, and appends it to a single output file preceded by a `URL:` header.
Public pages only, robots.txt-respecting, rate-limited. Auto-installs trafilatura
and lxml if missing. Uncapped by default; writes incrementally so partial output
is usable and the run is resumable-by-hand.

Usage:
  python text_mirror.py URL [URL ...] [--out FILE] [--delay SECONDS] [--max-pages N]
Examples:
  python text_mirror.py https://wiki.example.org/
  python text_mirror.py https://a.example/ https://b.example/x/ --out mirror.md --max-pages 500
"""
import argparse
import importlib.util
import os
import subprocess
import sys
import time
from collections import deque
from urllib.parse import urljoin, urldefrag, urlparse
from urllib.robotparser import RobotFileParser


def ensure_deps():
    """Install trafilatura/lxml into the current interpreter if not importable.

    Tries progressively more permissive pip strategies so it works in a venv, a
    normal user install, and a PEP 668 externally-managed (e.g. Homebrew) Python.
    """
    def missing():
        return [pkg for mod, pkg in (("trafilatura", "trafilatura"), ("lxml", "lxml"))
                if importlib.util.find_spec(mod) is None]

    need = missing()
    if not need:
        return
    print(f"[bootstrap] installing missing packages: {', '.join(need)}", flush=True)
    strategies = (
        ["--quiet"],
        ["--quiet", "--user"],
        ["--quiet", "--user", "--break-system-packages"],
        ["--quiet", "--break-system-packages"],
    )
    for extra in strategies:
        try:
            subprocess.check_call([sys.executable, "-m", "pip", "install", *extra, *need])
        except subprocess.CalledProcessError:
            continue
        importlib.invalidate_caches()
        if not missing():
            return
    sys.exit("[bootstrap] ERROR: could not install "
             f"{', '.join(need)}. Install manually (e.g. in a venv:\n"
             "  python3 -m venv .venv && . .venv/bin/activate && pip install trafilatura)")


ensure_deps()
import trafilatura  # noqa: E402
from trafilatura.settings import use_config  # noqa: E402
from lxml import html as LH  # noqa: E402

UA = "trafilatura-text-mirror/1.0 (+public-archive)"
CFG = use_config()
CFG.set("DEFAULT", "USER_AGENTS", UA)
CFG.set("DEFAULT", "DOWNLOAD_TIMEOUT", "30")

SKIP_EXT = (".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp", ".ico", ".css",
            ".js", ".pdf", ".zip", ".gz", ".tar", ".mp4", ".mp3", ".avi",
            ".woff", ".woff2", ".ttf", ".rss", ".atom")
# auth/action/non-content URLs skipped to stay public + avoid crawl explosions
SKIP_SUBSTR = ("/login", "/register", "/account", "/members/", "/lost-password",
               "/whats-new/", "/misc/", "/goto/", "/find-new/", "/attachments/",
               "/logout", "oauth", "action=", "oldid=", "diff=", "special:",
               "printable=", "veaction=", "redlink", "feed=", "&do=", "/help/")

LOG_FILE = None


def log(msg):
    line = f"[{time.strftime('%H:%M:%S')}] {msg}"
    print(line, flush=True)
    if LOG_FILE:
        with open(LOG_FILE, "a") as f:
            f.write(line + "\n")


def norm(url):
    return urldefrag(url)[0].strip()


def wanted(url, host):
    try:
        p = urlparse(url)
    except Exception:
        return False
    if p.scheme not in ("http", "https") or p.netloc != host:
        return False
    low = url.lower()
    if low.endswith(SKIP_EXT):
        return False
    if any(s in low for s in SKIP_SUBSTR):
        return False
    return True


def get_robots(host):
    rp = RobotFileParser()
    rp.set_url(f"https://{host}/robots.txt")
    try:
        rp.read()
    except Exception as e:
        log(f"robots-read-error {host}: {e}")
    return rp


def links_from(htmltext, base, host):
    try:
        doc = LH.fromstring(htmltext)
    except Exception:
        return []
    out = []
    for a in doc.xpath("//a[@href]"):
        u = norm(urljoin(base, a.get("href")))
        if wanted(u, host):
            out.append(u)
    return out


def crawl(seed, host, fout, delay, max_pages):
    rp = get_robots(host)
    q = deque([norm(seed)])
    seen = set(q)
    n = 0
    while q:
        if max_pages and n >= max_pages:
            log(f"{host}: hit max_pages={max_pages}"); break
        url = q.popleft()
        if not rp.can_fetch(UA, url):
            log(f"robots-skip {url}"); continue
        downloaded = None
        try:
            downloaded = trafilatura.fetch_url(url, config=CFG)
        except Exception as e:
            log(f"fetch-error {url}: {e}")
        time.sleep(delay)
        if not downloaded:
            log(f"no-content {url}"); continue
        n += 1
        try:
            text = trafilatura.extract(downloaded, url=url, output_format="markdown",
                                       include_comments=True, include_tables=True,
                                       favor_recall=True, config=CFG)
        except Exception as e:
            log(f"extract-error {url}: {e}"); text = None
        fout.write("\n\n" + "=" * 90 + "\n")
        fout.write(f"URL: {url}\n")
        fout.write("=" * 90 + "\n\n")
        fout.write((text or "[no extractable text]") + "\n")
        fout.flush()
        if n % 25 == 0:
            log(f"{host}: {n} pages saved, queue={len(q)}, seen={len(seen)}")
        for u in links_from(downloaded, url, host):
            if u not in seen:
                seen.add(u); q.append(u)
    log(f"{host}: DONE {n} pages, {len(seen)} urls discovered")
    return n


def main():
    global LOG_FILE
    ap = argparse.ArgumentParser(description="Recursive text-only site mirror via trafilatura.")
    ap.add_argument("seeds", nargs="+", help="one or more seed URLs to crawl")
    ap.add_argument("--out", default=os.environ.get("OUT_FILE", "text-mirror/mirror.md"),
                    help="output Markdown file (default: text-mirror/mirror.md)")
    ap.add_argument("--delay", type=float, default=float(os.environ.get("CRAWL_DELAY", "1.0")),
                    help="seconds between requests (default: 1.0)")
    ap.add_argument("--max-pages", type=int, default=int(os.environ.get("MAX_PAGES", "0")),
                    help="per-site page cap; 0 = unlimited (default: 0)")
    args = ap.parse_args()

    out_file = os.path.abspath(args.out)
    out_dir = os.path.dirname(out_file) or "."
    os.makedirs(out_dir, exist_ok=True)
    LOG_FILE = os.path.join(out_dir, "crawl.log")
    open(LOG_FILE, "w").close()

    with open(out_file, "w") as fout:
        fout.write("# Text-only mirror (trafilatura)\n")
        fout.write(f"# Generated {time.strftime('%Y-%m-%d %H:%M:%S')}\n")
        fout.write(f"# Delay={args.delay}s  MaxPages/site={args.max_pages or 'unlimited'}\n")
        total = 0
        for seed in args.seeds:
            host = urlparse(seed).netloc
            if not host:
                log(f"skip invalid seed (no host): {seed}"); continue
            log(f"=== crawling {host} from {seed} ===")
            fout.write(f"\n\n{'#' * 90}\n# SITE: {host}  (seed: {seed})\n{'#' * 90}\n")
            total += crawl(seed, host, fout, args.delay, args.max_pages)
        log(f"ALL DONE: {total} pages total -> {out_file}")


if __name__ == "__main__":
    main()
