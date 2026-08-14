#!/usr/bin/env python3
"""
md_to_pdf.py — convert a bf-triage Markdown report to a styled PDF.

Adapted from the ftdc-analysis skill
(10gen/employees: home/luke.pearson/llm-ftdc-analysis/SKILL.md).
Same approach: render Markdown → HTML with a couple of common extensions,
wrap in an inline-CSS template, then use WeasyPrint to write a PDF.

Usage:
    md_to_pdf.py [--auto-install] <input.md> <output.pdf>

Dependencies (one-time):
    pip3 install markdown weasyprint
    (or pass --auto-install and let the script handle it)

WeasyPrint pulls Pango/Cairo/gdk-pixbuf via the system package manager on
fresh hosts. On Ubuntu these are usually pre-installed; on macOS run
    brew install pango cairo gdk-pixbuf libffi
before the pip install. The --auto-install path can run pip but cannot
install brew / apt packages — those have to be done by the user.

--auto-install behaviour:
    When the Python dependencies are missing AND --auto-install is set,
    the script runs `python3 -m pip install ...` for `markdown` and
    `weasyprint`. The install target is:
      - the active venv when $VIRTUAL_ENV is set (no --user flag);
      - the user's site-packages otherwise (--user flag).
    The script never installs into the system site-packages and never
    invokes sudo. It prints the exact pip command before running it.

Exit codes (consumed by the skill's opt-in dispatch):
    0   PDF generated successfully
    1   bad CLI usage
    2   missing Python deps (markdown / weasyprint) — caller should
        print the install hint and continue with the MD-only report
        (or re-invoke with --auto-install)
    3   I/O error reading the input MD file
    4   WeasyPrint rendering / system-lib error (Pango / Cairo missing)
    5   --auto-install was requested but pip itself failed
"""

from __future__ import annotations

import os
import sys
from pathlib import Path

CSS = """
@page { size: letter; margin: 0.75in; }

/* Body & typography */
body {
    font-family: -apple-system, "Segoe UI", Helvetica, Arial, sans-serif;
    font-size: 10pt;
    line-height: 1.5;
    color: #1a1a1a;
}

/* Headings — h1 starts a new page (handy for Mode B combined outputs) */
h1 {
    font-size: 18pt;
    border-bottom: 2px solid #2563eb;
    padding-bottom: 6px;
    margin-top: 24pt;
    color: #1e3a5f;
    page-break-before: always;
}
h1:first-of-type {
    /* don't force a page break before the very first heading */
    page-break-before: avoid;
}
h2 {
    font-size: 14pt;
    border-bottom: 1px solid #d1d5db;
    padding-bottom: 4px;
    margin-top: 20pt;
    color: #1e40af;
    page-break-after: avoid;
}
h3 {
    font-size: 12pt;
    margin-top: 16pt;
    color: #1e40af;
    page-break-after: avoid;
}

/* Tables */
table {
    border-collapse: collapse;
    width: 100%;
    margin: 10pt 0;
    font-size: 8.5pt;
    page-break-inside: avoid;
}
th {
    background-color: #1e40af;
    color: white;
    padding: 6px 8px;
    text-align: left;
    font-weight: 600;
}
td {
    padding: 5px 8px;
    border-bottom: 1px solid #e5e7eb;
    vertical-align: top;
}
tr:nth-child(even) { background-color: #f8fafc; }

/* Inline + block code */
code {
    background-color: #f1f5f9;
    padding: 1px 4px;
    border-radius: 3px;
    font-family: "SFMono-Regular", Consolas, Menlo, monospace;
    font-size: 8.5pt;
}
pre {
    background-color: #1e293b;
    color: #e2e8f0;
    padding: 12px;
    border-radius: 6px;
    font-size: 8pt;
    line-height: 1.4;
    page-break-inside: avoid;
}
pre code { background-color: transparent; padding: 0; color: #e2e8f0; }

/* Blockquotes — bf-triage uses these for "Limited evidence" banners */
blockquote {
    border-left: 4px solid #f59e0b;
    background-color: #fffbeb;
    padding: 8px 12px;
    margin: 10pt 0;
    color: #92400e;
}

strong { color: #1e3a5f; }

/* Avoid breaking lists awkwardly */
ul, ol { page-break-inside: auto; }
li { page-break-inside: avoid; }
"""


def _try_imports() -> tuple[bool, str]:
    """Returns (ok, missing_module_name)."""
    try:
        import markdown  # noqa: F401
    except ImportError:
        return False, "markdown"
    try:
        from weasyprint import HTML  # noqa: F401
    except ImportError:
        return False, "weasyprint"
    return True, ""


def _auto_install() -> int:
    """Run `pip install` for markdown + weasyprint into venv or --user.

    Returns 0 on success, 5 on pip failure.
    """
    import subprocess

    pkgs = ["markdown", "weasyprint"]
    in_venv = bool(os.environ.get("VIRTUAL_ENV"))
    cmd = [sys.executable, "-m", "pip", "install", "--quiet"]
    if not in_venv:
        cmd.append("--user")
    cmd.extend(pkgs)

    print(
        f"[md_to_pdf] auto-installing dependencies ({'venv' if in_venv else '--user'} target):",
        file=sys.stderr,
    )
    print(f"[md_to_pdf]   {' '.join(cmd)}", file=sys.stderr)
    try:
        subprocess.check_call(cmd)
    except subprocess.CalledProcessError as exc:
        print(
            f"[md_to_pdf] pip install failed (rc={exc.returncode}). "
            "Run it manually and re-try.",
            file=sys.stderr,
        )
        return 5
    except FileNotFoundError:
        print(
            f"[md_to_pdf] {sys.executable} -m pip not available. "
            "Install pip or run pip install manually.",
            file=sys.stderr,
        )
        return 5
    return 0


def render(md_path: Path, pdf_path: Path, auto_install: bool = False) -> int:
    ok, missing = _try_imports()
    if not ok:
        if not auto_install:
            print(
                f"[md_to_pdf] missing dependency: {missing}. Install with:\n"
                "    pip3 install --user markdown weasyprint\n"
                "  (on macOS also: brew install pango cairo gdk-pixbuf libffi)\n"
                "  Or re-invoke this script with --auto-install to let it pip-install for you.",
                file=sys.stderr,
            )
            return 2
        rc = _auto_install()
        if rc != 0:
            return rc
        import importlib
        importlib.invalidate_caches()
        ok, missing = _try_imports()
        if not ok:
            print(
                f"[md_to_pdf] post-install import still failing for {missing}. "
                "Check that the installed version is on this interpreter's path.",
                file=sys.stderr,
            )
            return 2

    import markdown
    from weasyprint import HTML

    try:
        md_text = md_path.read_text(encoding="utf-8")
    except OSError as exc:
        print(f"[md_to_pdf] cannot read {md_path}: {exc}", file=sys.stderr)
        return 3

    html_body = markdown.markdown(
        md_text,
        extensions=["tables", "fenced_code", "codehilite", "sane_lists"],
    )

    full_html = (
        "<!DOCTYPE html>\n"
        '<html><head><meta charset="utf-8">'
        f"<style>{CSS}</style>"
        f"</head><body>{html_body}</body></html>"
    )

    try:
        HTML(string=full_html, base_url=str(md_path.parent)).write_pdf(str(pdf_path))
    except Exception as exc:  # noqa: BLE001 — surface WeasyPrint-internal errors verbatim
        print(
            f"[md_to_pdf] WeasyPrint failed: {exc}\n"
            "  Common cause: missing Pango/Cairo system libraries. On Ubuntu:\n"
            "    sudo apt-get install -y libpango-1.0-0 libpangoft2-1.0-0 \\\n"
            "        libcairo2 libgdk-pixbuf2.0-0 libffi-dev shared-mime-info",
            file=sys.stderr,
        )
        return 4

    print(f"[md_to_pdf] wrote {pdf_path} ({pdf_path.stat().st_size} bytes)")
    return 0


def main(argv: list[str]) -> int:
    args = list(argv[1:])
    auto_install = False
    if "--auto-install" in args:
        auto_install = True
        args = [a for a in args if a != "--auto-install"]
    if len(args) != 2:
        print(
            "Usage: md_to_pdf.py [--auto-install] <input.md> <output.pdf>",
            file=sys.stderr,
        )
        return 1
    md_path = Path(args[0]).expanduser().resolve()
    pdf_path = Path(args[1]).expanduser().resolve()
    if not md_path.exists():
        print(f"[md_to_pdf] input file not found: {md_path}", file=sys.stderr)
        return 3
    pdf_path.parent.mkdir(parents=True, exist_ok=True)
    return render(md_path, pdf_path, auto_install=auto_install)


if __name__ == "__main__":
    sys.exit(main(sys.argv))
