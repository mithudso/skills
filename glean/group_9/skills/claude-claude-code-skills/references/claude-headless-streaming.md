# Claude Code headless streaming: the `-p` / stream-json client contract

> Provenance: researched 2026-08-04 for a Python TUI that shells out to `claude`.
> Flag claims come from real `claude --help` on **2.1.221** plus empirical
> acceptance probes. Event shapes are derived from the installed
> `claude-agent-sdk` **0.2.87** (the authoritative consumer of this wire
> protocol) plus official docs. **No live full stream was captured** — the probe
> environment blocked org verification — so the example stream below is
> reconstructed, not recorded. Two wire lines are labeled as genuinely observed.
> `verified-as-of: 2026-08-04`

## The contract in one paragraph

`claude -p --output-format stream-json --verbose` emits **NDJSON, one object per
line**, with a top-level `type`. The set is `system` (subtype-discriminated),
`assistant`, `user`, `result`, `stream_event`, `rate_limit_event`,
`control_response` / `control_request` / `control_cancel_request`, and
`prompt_suggestion`. `result` is the terminal event for a turn — but it is
**neither guaranteed nor necessarily last**.

## Five rules that decide whether your client works

**1. Close the child's stdin.** With stdin left open, `claude -p` waits for piped
input and prints, verbatim (observed on stderr):

```
Warning: no stdin data received in 3s, proceeding without it. If piping from a
slow command, redirect stdin explicitly: < /dev/null to skip, or wait longer.
```

That is a **fixed 3-second penalty on every run**, and it is often most of a
perceived slow start. Pass `stdin=DEVNULL` (or `< /dev/null`) unless you are
using `--input-format stream-json`, in which case stdin is your input channel.

**2. Buffer speculatively; never assume one line is one object.** The SDK
accumulates until `json.loads` succeeds, caps the buffer at 1 MiB
(`_DEFAULT_MAX_BUFFER_SIZE = 1024 * 1024`), and skips non-JSON lines when *not*
mid-parse — its own comment: *"Skip non-JSON lines (e.g. `[SandboxDebug]`) when
not mid-parse — they corrupt the buffer otherwise."* Copy both behaviors.

**3. Ignore unknown `type` values.** The SDK's parser ends in a
forward-compatible default branch that skips unrecognized types *"so newer CLI
versions don't crash older SDK versions."* Drift here is **additive** — new
fields, versioned (`api_error_status` since v2.1.110, `capabilities` v2.1.205+) —
not renames. A negation sweep found **no** source claiming the schema is unstable
or that a breaking change landed.

**4. There are four terminal states, not two.**

| State | Detection |
|---|---|
| Success | `result` with `subtype == "success"` → render `result` |
| Run-level error | `result` with an `error_*` subtype → **the `result` field is absent**; use `subtype`, `errors`, `api_error_status` |
| API-level error | `result` with `subtype == "success"` **but `is_error == true`** → `api_error_status` carries the HTTP code (429/500/529) |
| **No `result` at all** | Process exits having emitted no `result`. Docs claim failures print as the result on stdout; an observed startup auth failure instead gave **exit 1, zero stdout lines, message on stderr only.** |

Treat "exited without `result`" as first-class. Also **do not break on `result`** —
docs: *"A small number of trailing system events, such as `prompt_suggestion`,
can arrive after it, so iterate the stream to completion."*

**5. Cancellation forces an architecture change, not a flag.** In bare-argv
prompt mode there is no clean interrupt. Three independent confirmations: the SDK
raises `Exception("Control requests require streaming mode")`; its transport
always adds `--input-format stream-json` and writes prompts to stdin; docs state
single-message input *"does not support … Real-time interruption."* To cancel and
keep the session, send the prompt on stdin:

```json
{"type":"user","message":{"role":"user","content":"…"},"parent_tool_use_id":null,"session_id":"default"}
```

then interrupt in-band:

```json
{"type":"control_request","request_id":"req_1_a3f19c04","request":{"subtype":"interrupt"}}
```

The reply arrives as `{"type":"control_response","response":{"request_id":"…","subtype":"success"|"error"}}`.
Other shipping subtypes: `set_permission_mode`, `set_model`, `mcp_status`,
`get_context_usage`, `stop_task`, `rewind_files`, `initialize`. **Feature-detect**
via `capabilities` on `system/init` (`interrupt_receipt_v1`,
`interrupt_cancel_queued_v1`), not version strings.

## "Still thinking" vs "hung"

A wall-clock startup timer cannot tell these apart. With
`--include-partial-messages` you get token-level liveness: `stream_event` wraps
the **raw** Claude API SSE events (`message_start`, `content_block_delta`, …),
text at `event.delta.text` for `text_delta`, and `ttft_ms` on `message_start`.
Tool input arrives as `input_json_delta` / `partial_json` — **you accumulate it
yourself**. Stream events are main-session only; subagent deltas are never
forwarded. Without that flag there is **no documented heartbeat**.

Liveness signals: `stream_event` deltas · `system/api_retry` (carries `attempt`,
`max_retries`, `retry_delay_ms`, `error_status`) = *retrying, not hung* ·
`system/task_progress` · `rate_limit_event` · one `assistant` message per turn.
Failure signals: **process exit plus stderr**, not elapsed time. Prefer bounded
runs — `--max-turns` and `--max-budget-usd` yield deterministic
`error_max_turns` / `error_max_budget_usd` results instead of open-ended hangs.
`claude auth status` is a cheap JSON preflight (exit 0, no full startup).

## Sessions

`--continue`/`-c` resumes the most recent conversation **in the cwd**;
`--resume`/`-r <id|name>` a specific one — *"passing a session ID searches only
the current project directory and its git worktrees."* `--fork-session` copies
history under a new id; `--session-id <uuid>` pre-assigns; `--no-session-persistence`
writes nothing (print mode only). Capture the id from `result.session_id`
(present on **every** result, success or error) or `system/init.session_id`.
Transcripts live at `~/.claude/projects/<encoded-cwd>/<session-id>.jsonl`, where
`<encoded-cwd>` is the absolute cwd with every non-alphanumeric character
replaced by `-`; honors `CLAUDE_CONFIG_DIR`. **A mismatched cwd silently yields a
fresh session instead of the resumed one** — the most common resume failure.

## Signals

- **SIGTERM → exit 143** (documented): aborts the turn, terminates the process
  tree of any running Bash command, runs `SessionEnd` hooks, exits 143.
- **SIGINT mid-stream: undocumented.** Two third-party observations report exit
  130 (just POSIX 128+2, not a contract) and "kills the process entirely in
  headless mode". Whether the session stays resumable is unknown.
- Use SIGTERM for a graceful kill; use the in-band `control_request` interrupt
  when you need the session to survive.

## Shutdown timing to budget for

Background Bash shells are killed ~5 s after the final result and stdin close.
Background subagents are waited on, capped at 10 min
(`CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS`). Slow consumers get a drain wait scaled
to queue depth, **capped at 30 s** (was ~2 s before v2.1.214).

## Gotchas

1. **`CLAUDECODE` / `CLAUDE_CODE_*` inherited by a nested `claude -p` can hang
   it** (anthropics/claude-code#26190). Scrub them from the child environment
   when spawning from inside a Claude session.
2. **stderr carries ANSI escapes even when redirected to a file** — observed:
   `^[[0m^[[31m…Unable to verify organization…`. Strip ANSI before rendering.
3. **Pipe buffering is unresolved and decides your transport.** Issue #25670
   (2.1.41, macOS): stdout *"block-buffered instead of line-buffered… JSON lines
   accumulate in the pipe buffer (~4-8KB) and don't appear until the buffer fills
   or the process exits"*; bot-closed as duplicate, not fixed. Counter-evidence:
   the official SDK reads stdout via `PIPE` and works. **Discriminating test:**
   spawn with stdout as a pipe and check whether `system/init` arrives *before*
   the run completes. Early → pipes are fine. Only at exit → wrap in a pty.
4. **Hang after `result`** — #25629 (2.1.38, Linux): result emitted, stdout stays
   open, process hangs 5+ min. Community workaround: 30 s timer post-`result`,
   then SIGTERM → SIGKILL.
5. **Print mode ≠ TUI.** `-p` runs are classified `sdk-cli` and can diverge from
   the TUI request scaffold; HTTP MCP servers and claude.ai connectors reportedly
   don't load in `-p` (#34131, #26364, #59105).
6. **Piped stdin is capped at 10 MB** (since v2.1.128).
7. **Permission blocks surface on three channels**: `result.permission_denials`
   (element shape undocumented), a denied tool appearing as a `user` message with
   an `is_error` tool_result, and `--permission-prompt-tool` if you own the
   decision. `--permission-mode dontAsk` is the doc-recommended locked-down
   headless posture. `AskUserQuestion`, org-`ask` connector tools, and MCP tools
   marked `requiresUserInteraction` are denied **even under `bypassPermissions`**.
8. **Flags come in three tiers.** (a) in `--help` *and* docs — safe; (b)
   documented but hidden from `--help`: `--max-turns`, `--permission-prompt-tool`,
   `--system-prompt-file`; (c) in neither yet accepted by 2.1.221:
   `--max-thinking-tokens`, `--session-mirror`, `--task-budget`, `--thinking`,
   `--thinking-display`. Tier (c) has no stability contract, and `--max-turns`
   sitting in (b) is your only bounded-run guarantee.
9. **The schema is de-facto stable, not contractually documented.** Issue #24612
   asked for a schema reference and was closed "not planned": developers *"must
   either trial-and-error their way through the output or read the Agent SDK
   TypeScript type definitions."*
10. **SDK types are a superset of the wire format.** `types.py` marks
    `mirror_error` as *"SDK-synthesized … never emitted by the CLI subprocess."*
    Strong hint, not a contract.

## If your host is Python

`claude-agent-sdk` already implements the buffered NDJSON reader, the control
protocol, `interrupt()`, stderr callbacks, and exit-code mapping. Docs frame the
CLI as *"using the Agent SDK via the CLI (`claude -p`)"* and point to the
Python/TS packages *"for full programmatic control."* Shelling out is the
documented path for **non**-Python/TS hosts; from Python, weigh the SDK first.

## Open gaps

Pipe-vs-PTY buffering on current versions (test above) · no exit-code table
exists anywhere beyond 0/nonzero and 143 · `permission_denials` element shape ·
SIGINT semantics and whether the session survives · whether `--verbose` is
required at runtime for stream-json (the parser accepts its absence; pass it
anyway — the SDK does unconditionally) · whether `task_started`/`task_progress`/
`task_notification` are emitted unconditionally.
