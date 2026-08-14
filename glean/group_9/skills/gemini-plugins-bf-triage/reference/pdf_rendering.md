# Reference: Step 9b — Optional PDF rendering

> Read this file only when the agent has decided to generate a PDF
> alongside the Markdown report. The main `SKILL.md` Step 9b section
> points here.

Run only when **either** of the following is true:

- `BF_TRIAGE_GENERATE_PDF=1` is set in the environment, **or**
- the user's prompt explicitly asks for a PDF (case-insensitive match
  on phrases like "and a pdf", "pdf too", "as pdf", "generate pdf").

Output path: replace the `.md` extension on the **versioned** MD path
computed in Step 9 with `.pdf` (same directory, same basename, same
`-vN` suffix). Examples:

- `./bf-reports/BF-43272-triage.md` → `./bf-reports/BF-43272-triage.pdf`
- `./bf-reports/BF-43272-triage-v3.md` → `./bf-reports/BF-43272-triage-v3.pdf`
- `./bf-reports/team-<slug>-<date>/BF-43272-triage.md` →
  `./bf-reports/team-<slug>-<date>/BF-43272-triage.pdf` (Mode B)

The PDF inherits the MD's suffix. The agent does NOT run the Step 9
`next_free` helper a second time for the PDF — pairing the MD and PDF
filenames is the goal. If the PDF path happens to already exist (it
shouldn't, given the MD just got a fresh `-vN`), overwrite it: the MD
is the source of truth.

## Invocation rules

The PDF step has two callers with different consent boundaries:

1. **Coordinator (parent agent in Mode A, or Mode B coordinator before
   fan-out, or Mode C coordinator after aggregation)** — invoke with
   `--auto-install`:

   ```bash
   python3 "${BF_TRIAGE_SKILL_DIR}/scripts/md_to_pdf.py" \
       --auto-install <report>.md <report>.pdf
   ```

   The user enabled the feature explicitly (env var or prompt), which
   carries consent to install the two pip packages. The script will
   install into the active venv (if `$VIRTUAL_ENV` is set) or to the
   user's site-packages (`--user`) — never system-wide, never `sudo`.

2. **Subagents (Mode B parallel-fan-out triagers, Mode C held-in
   triagers / graders)** — invoke **without** `--auto-install`:

   ```bash
   python3 "${BF_TRIAGE_SKILL_DIR}/scripts/md_to_pdf.py" \
       <report>.md <report>.pdf
   ```

   Subagents have no live user channel to confirm an install. If the
   deps are missing they degrade silently (exit 2) and the parent agent
   handles the install on the next coordinator-level invocation, or
   the user installs manually.

## Exit-code interpretation (graceful degradation)

The PDF step MUST NOT fail the triage run. Interpret the script's exit
code as follows:

| Exit code | Meaning | Action |
| --------- | ------- | ------ |
| `0` | PDF written | Log the path next to the MD path in the final user message. |
| `2` | Missing `markdown` / `weasyprint` Python deps | If invoked from a subagent: log a one-line note and continue MD-only. If invoked from the coordinator without `--auto-install` (e.g. user-disabled auto-install via prompt): log the manual `pip3 install` hint. |
| `3` | Cannot read input MD | Should not happen (we just wrote it). Log and continue. |
| `4` | WeasyPrint failed (usually Pango/Cairo missing on macOS) | Log the brew install hint the script prints. Continue MD-only. |
| `5` | `--auto-install` was requested but `pip install` itself failed | Log the failure verbatim. Suggest the user run `pip3 install --user markdown weasyprint` manually. Continue MD-only. |
| any other | Unknown | Log stderr verbatim. Continue MD-only. |

Do NOT prompt the user mid-run to install deps — the `--auto-install`
flag already encodes that consent. If the user prefers an explicit
prompt instead, they can unset `BF_TRIAGE_GENERATE_PDF` and ask for
PDF case-by-case in the prompt.

In Mode B, run Step 9b **per BF**, after each subagent's MD write.
The team `index.md` is intentionally NOT converted to PDF (it's a
short summary index, not a reading artifact).

## Stage 2 — Attach the PDF to a Jira issue (opt-in, Step 11b)

This is the **one sanctioned non-gateway write** in the skill. It runs
ONLY when all of these hold (otherwise skip silently and just post the
comment):

1. A comment is being posted this run (`POST_COMMENTS=auto`, or the
   `ask` gate got a "yes").
2. PDF attachment was opted into: the prompt asks for it
   (case-insensitive "attach the pdf", "post the pdf", "pdf to the
   ticket", "with the pdf") OR `BF_TRIAGE_ATTACH_PDF=1`.
3. Prerequisites resolve (see below). Any failure → skip the upload,
   keep the comment, note the skip in the final summary. NEVER fail
   the triage over an attachment.

Attachment is **coordinator-only** (like comment posting) — subagents
never upload.

### Prerequisites

- The PDF exists: run Stage 1 first (the normal Step 9b render). If
  `md_to_pdf.py` did not exit `0`, do not attempt upload.
- `python3` (stdlib only — no `curl`, no pip deps).
- A Jira Personal Access Token (PAT) is resolvable. The upload script
  searches, **first hit wins, regardless of which runtime is active**:
  1. `$JIRA_PERSONAL_TOKEN` in the environment.
  2. The first `JIRA_PERSONAL_TOKEN` found scanning `~/.cursor/mcp.json`,
     then `~/.claude.json`, then `~/.claude/settings.json` (the key may
     be nested anywhere, e.g. under `mcpServers.<server>.env`).
  So a Claude Code run will use a token sitting in the Cursor config,
  and vice-versa. If none is found the script exits `2` — skip the
  upload (the gateway has no attach tool, so there is no fallback).

### Procedure

Run the shipped script — it resolves the token in-process (never on a
command line, never echoed), builds the multipart upload, and POSTs to
the Jira REST attachments endpoint:

```bash
python3 "${BF_TRIAGE_SKILL_DIR}/scripts/attach_pdf_to_jira.py" \
    <BF-KEY> <report>.pdf
```

Interpret the exit code (never fail the triage over an attachment):

| Exit | Meaning | Action |
| ---- | ------- | ------ |
| `0` | Uploaded (HTTP 200) | Note `PDF attached to <BF-KEY>` + the printed attachment id in the final summary. |
| `2` | No PAT found (env or any MCP config) | Skip — keep the comment, note `PDF attach skipped (no Jira token)`. |
| `3` | PDF not readable | Stage 1 render must have failed; note and continue MD/comment-only. |
| `4` | Upload failed (non-200 / network) | Log the script's stderr verbatim, keep the comment, retry at most once. |
| `1` | Bad usage | Should not happen; log and continue. |

Confirm a success if desired via the read-only gateway tool
`jira_list_attachments(issue_key="<BF-KEY>")`.

### Security / visibility caveats (surface these to the user)

- **Attachments are NOT visibility-restricted the way comments are.**
  A `developers_only` comment is hidden from non-developers, but a PDF
  attachment is visible to **anyone who can view the issue**. This is
  the main reason attachment requires its own explicit opt-in (it is
  not implied by `POST_COMMENTS`): the triage PDF embeds log snippets,
  so the user must consciously accept the wider exposure.
- The upload **bypasses the MCP gateway** and its audit path — it is a
  direct authenticated write to Jira. Only do it under the opt-in.
- The PAT is plaintext in the MCP config; never echo it, never inline
  it on the command line, never write it into the report or any
  spooled file.
