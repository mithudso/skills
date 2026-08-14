# wt-cli-tool

**Category:** Science, Biology & Medicine
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/storage-engines/wt-cli-tool/skills/wt-cli-tool

## Description
How to use the WiredTiger `wt` command-line utility to inspect, dump, load, verify, salvage, or otherwise operate on a WiredTiger database from the shell. Use this skill whenever the user asks about the `wt` tool, mentions inspecting WiredTiger files (`.wt` files, the history store, the metadata, the WAL), wants to recover or salvage a corrupted WiredTiger database, dump or load a WiredTiger table, run `printlog` against a WiredTiger journal, look at a MongoDB data directory at the storage-engine level, or asks any question of the form "how do I X with wt", "what does `wt <subcommand>` do", or "can wt show me Y". Use this skill even when the user does not explicitly say "wt" — phrasing like "show me what's in this WiredTiger directory", "extract data from these `.wt` files", "verify our storage files", or "run recovery on this database" should all pull this skill in.

---

# Using the `wt` command-line utility

`wt` is WiredTiger's command-line entry point. It opens a WiredTiger database the same way an embedding application would (`wiredtiger_open` + `WT_SESSION`) and then runs a single subcommand against it. Almost everything you can do via the C API has a `wt` equivalent, which makes it the right tool for inspecting, exporting, repairing, or sanity-checking a database from the shell.

This skill exists because, in practice, users describe what they want ("I have a corrupted MongoDB data dir, can I get the documents out?", "diff two checkpoints", "is this WAL replayable?") rather than naming a subcommand. Translating intent into the right `wt` invocation — including the right *global* flags — is most of the work.

## Safety: read-only is the default

**Always open the database read-only (`-r`) unless the user has explicitly authorized a write.** `wt` is a privileged tool — many of its subcommands rewrite or destroy data in place, sometimes irreversibly, and `-r` is the seatbelt that prevents any of them from doing so regardless of which subcommand you pair it with.

Concretely, when this skill runs:

- **Suggest, don't execute, anything that writes.** If the user's request resolves to a destructive command — `salvage`, `drop`, `truncate`, `write`, `alter`, `compact`, `create`, `load`/`loadtext`, `downgrade`, or any global open mode that writes (`-R`, `-S`, default open without `-r`) — print the proposed command and the reason for running it, name what it will rewrite, and **stop and ask the user to confirm** before running it. Do not run it autonomously. This rule applies even when the user's intent obviously requires the write (e.g. "salvage this collection"). The user must explicitly authorize the specific write step.
- **Default every diagnostic to `-r`.** `list`, `stat`, `dump`, `read`, `printlog`, `verify` are all read-only-safe and pair naturally with `-r`. Never drop `-r` "just in case" — read-only mode rejects destructive operations as a feature, which gives the user a second chance to notice if a command was wrong.
- **Salvage is in a class by itself.** It rewrites the file and *discards* anything it can't recover. It is a last resort, after diagnosis, after a backup, and only with explicit user authorization. See the `salvage` subcommand below for the full warning.
- **Never use `--no-verify`-style escape hatches** like the `-F` (force) flag on `salvage` unless the user explicitly asks for them. They suppress the safety checks that would otherwise stop a wrong command.

When in doubt, recommend the read-only path, capture diagnostic output, and let the user decide whether the write is worth running.

## How to think about `wt`

A `wt` invocation has two parts:

```
wt [GLOBAL_OPTIONS] <subcommand> [SUBCOMMAND_OPTIONS] [args...]
```

The **global options** control how the database is *opened*: home directory, recovery, read-only-ness, encryption, logging, etc. The **subcommand** is what to do once it's open. The two are independent — e.g. you can pair `-r` (read-only) with almost any subcommand that doesn't write.

Because `wt` calls `wiredtiger_open`, it must be run *against* a database directory. Either `cd` into it, or pass `-h <dir>`. The error `unable to open... no such file or directory` almost always means the user is in the wrong directory or forgot `-h`.

A few non-obvious things worth keeping in front of mind:

- **`wt` defaults to refusing to run recovery**. If the database wasn't cleanly shut down, you'll get an error and have to opt in with `-R` (run recovery) or `-S` (salvage). This is intentional — silent recovery on a damaged database is how you destroy evidence.
- **For diagnosing a possibly-corrupt database, prefer `-r` (read-only)**. It opens the connection so it cannot write back, which means whatever state the user has is preserved no matter what subcommand you run.
- **MongoDB users need a config string**. MongoDB writes WiredTiger files with a non-default journal path and snappy compression; opening such a directory with stock `wt` fails. Pass `-C "log=(enabled=true,path=journal,compressor=snappy)"`. The `wt -?` output includes this string as a hint.
- **If `wt` doesn't know about snappy** (e.g. a vanilla WT checkout without `--with-builtins=snappy`), you'll see `unknown compressor 'snappy': Invalid argument`.

## Global options reference

| Flag | Purpose |
|---|---|
| `-h <dir>` | Database home directory. Default is `.` (cwd). |
| `-C <config>` | Extra `wiredtiger_open` config (e.g. MongoDB's `log=(...,compressor=snappy)`). |
| `-E <key>` | Encryption secret key. |
| `-r` | Open read-only. **Default. Use this unless the user has explicitly authorized a write.** |
| `-R` | Run recovery. **Writes** to the DB (replays the WAL into the on-disk state). Required if the DB wasn't cleanly shut down — but only if the user wants to repair, not just inspect. Confirm with the user before using. |
| `-S` | Run salvage recovery. **Writes**, lossy. Last resort. Confirm with the user before using, and only on a copy. |
| `-L` | Force logging off. Useful for `printlog`-only inspection in some contexts. |
| `-l <path>` | Open a database that's mid-`live_restore`, with the source path. |
| `-m` | Verify metadata while opening. |
| `-p` | Disable cursor pre-fetching. Use when dumping/verifying possibly-corrupt data. |
| `-B` | Maintain release 3.3 log file compatibility. |
| `-V` | Print library version and exit. |
| `-v` | Verbose. |
| `-?` | Print help. Works with subcommands too: `wt dump -?`. |

`-L`, `-R`, `-S` are mutually exclusive. `-r` cannot combine with `-R` or `-S`.

## Subcommands

Each subcommand below shows its synopsis, then the things actually worth knowing about it. Where details get long (common workflows), follow the pointer to the reference file.

### Inspecting / reading

#### `list` — show what's in the database

```
wt list [-cv] [-f output] [uri]
```

The most common starting point. Without args, prints every URI in the metadata. With a URI, prints details for just that one. `-c` adds checkpoint info, `-v` adds the full schema config string. Read-only-safe.

#### `stat` — runtime statistics

```
wt stat [-f] [uri]
```

Prints engine-wide stats, or table-level stats if a URI is given. `-f` limits output to "fast" stats (cheap to collect — same as the `statistics=(fast)` cursor config). Stat names match the `stat.h` definitions in the wiredtiger source tree, so they're searchable.

#### `dump` — export a table

```
wt dump [-ejnprx] [-c checkpoint] [-f output] [-k key]
        [-l lower-bound] [-t timestamp] [-u upper-bound] [-w window] uri
```

The workhorse. Defaults to a text format that `wt load` can re-import. Important options:

- `-j` — JSON output. Easier for ad-hoc inspection and for piping to `jq`.
- `-p` — pretty-print (printable characters left as-is). Combine with `-x` to keep raw byte arrays in hex while pretty-printing the rest. **`-p` output is not loadable.**
- `-x` — hex-encode everything.
- `-c <checkpoint>` — dump as of a named checkpoint instead of "now". Pair with `-r` to inspect prior states without disturbing the live one.
- `-t <timestamp>` — dump as of an MVCC timestamp.
- `-k <key>` / `-l <lower>` / `-u <upper>` / `-w <window>` — restrict the range. `-w n` dumps `n` records on either side of `-k`. `-n` means "if `-k` isn't found, return whatever `search_near` lands on" (useful for "is anything near this key?").
- `-r` *(subcommand flag, not the global `-r`)* — reverse order.


#### `read` — fetch specific keys

```
wt read uri key ...
```

Prints the value for each key. Exits non-zero if any key is missing. Only works on tables with string or recno keys *and* string values. For anything more complex, use `dump -k`.

#### `printlog` — display the WAL

```
wt printlog [-mux] [-f output]
            [-l start-file,start-offset]
            [-l start-file,start-offset,end-file,end-offset]
```

Decodes the write-ahead log to text. Output is JSON-like. Notable:

- By default user data is **redacted** — keys and values appear as placeholders. Pass `-u` to show real bytes (only do this when the user explicitly wants application data).
- `-m` filters to message-type records only (the things logged by `WT_SESSION.log_printf`). Good for finding application annotations in the log without wading through commits.
- `-x` adds hex alongside string-formatted items.
- `-l` bounds the log range. The argument is comma-separated LSN components.

`printlog` doesn't run recovery, so it's safe on an unclean database. (`util_main.c` forces `log=(enabled=false)` for this command.)

### Modifying — destructive, requires explicit user authorization

> **Every subcommand in this group writes to the database.** Some are also lossy (`salvage`, `truncate`, `drop`). Never run any of these autonomously on a database the user did not explicitly tell you to modify. The default response when one of them looks like the answer is: explain the proposed command, what it will rewrite, and **ask the user to confirm before running it**. For investigations, use `dump`/`printlog`/`verify`/`stat`/`list` under `-r` instead.

#### `create` — make a table

```
wt create [-c configuration] uri
```

Equivalent to `WT_SESSION.create(uri, configuration)`. **Writes** (adds metadata + new file). Pass a config like `"key_format=S,value_format=S"` via `-c`.

#### `alter` — change table config

```
wt alter uri configuration ...
```

Pairs of `uri configuration` **modify** config values that were originally passed to `create`. The URI may be a *prefix* as long as it matches exactly one object.

#### `drop` — delete a table

```
wt drop uri
```

**Destroys the table and its file.** Equivalent to `WT_SESSION.drop` with `force`. Irreversible — confirm a backup is in hand first.

#### `truncate` — empty a table

```
wt truncate uri
```

**Removes every row** from the table. Whole-table only — there is no range form here (use the C API for that). Irreversible — confirm a backup is in hand first.

#### `write` — insert/remove records

```
wt write [-aor] uri key value ...
wt write -a uri value ...           # column store: append
wt write -r uri key                 # remove key
```

**Writes** key/value pairs (refuses to overwrite without `-o`). `-o` allows overwriting, `-a` appends to a column store (assigns recno keys), `-r` removes a key. Be especially careful with `-o` (silent overwrite) and `-r` (delete).

#### `compact` — reclaim space

```
wt compact [-c configuration] uri
```

**Rewrites the table in place**, returning unused space to the OS. `-c` forwards a config to `WT_SESSION.compact` (e.g. `"timeout=0"`). Generally safe (compact is crash-safe), but it does take a write lock on the file and changes its layout.

#### `salvage` — recover a corrupt table

```
wt salvage [-F] uri
```

> **STOP before suggesting `salvage`.** It **rewrites the file in place** and **discards any pages it cannot recover** — silent data loss is the *expected* outcome. No dry-run, no undo.
>
> Last resort, only after read-only diagnosis (`verify`, `list`, `dump`) has identified specific damage *and* the user has explicitly chosen salvage over alternatives (restore from backup, rebuild an index, accept the loss). Never run autonomously even if the prompt says "the database is corrupted" — diagnose read-only first and present findings.
>
> When authorized: (1) `cp -a` the data directory aside, (2) run salvage against the **copy**, never the original, (3) verify the copy, (4) decide whether to swap it in. **Do not pass `-F`** (force) without an additional explicit confirmation — it bypasses the file-format safety check. Pair with `-S` at the global level if recovery itself is needed to open the DB.

#### `load` / `loadtext` — import data

```
wt load     [-ajn] [-f input] [-r name] [uri configuration ...]
wt loadtext [-f input] uri
```

**Writes** records into a (possibly existing) table. `load` reads `dump` output (text or `-j` JSON). `loadtext` reads bare lines — for row stores, line pairs are key then value; for column stores, every line is appended as a value. `-r` renames the destination URI; `-n` makes overwrite an error (use it to be safer); `-a` ignores incoming recno keys and assigns new ones. Confirm with the user before running, and prefer `-n` to avoid silent overwrite.

#### `backup` — copy the database

```
wt backup [-t uri] directory
```

Hot backup of the whole database into the target directory (which becomes a fresh, openable WiredTiger home). `-t uri` restricts to specific data sources. Backup is read-only with respect to the source, but it does write to `directory` — make sure the user has approved the destination path.

### Lifecycle / verification

#### `verify` — structural integrity check

```
wt verify [-acSstu] [-d <dump-spec>] [uri]
```

Walks the table's pages and confirms invariants. Without a URI, verifies every table. Important options:

- `-a` — abort on the first damaged table (default is to continue and report all).
- `-c` — continue past per-page errors within a single table.
- `-S` — strict mode (treat any oddity as an error).
- `-s` — verify against the *stable* timestamp (only valid after rollback-to-stable).
- `-d <spec>` — emit verification dumps. Specs include `dump_address`, `dump_blocks`, `dump_layout`, `dump_tree_shape`, `dump_offsets=<offsets>`, `dump_pages`. `-k` and `-u` further restrict what the page/block dump shows. Use `-p` (global) to disable prefetching when looking at corrupted data.

`verify` is the right diagnostic when `dump` reports specific corruption — it can pinpoint where in the tree the damage is.

#### `downgrade` — set compatibility version

```
wt downgrade -V <release>
```

Marks the database as compatible with an older release so an older binary can open it. `-V` is required.

#### `copyright`

```
wt copyright
```

Prints the copyright. No-op for diagnostic purposes; included for completeness.

## URIs

`wt` arguments expecting a "uri" use WiredTiger's URI scheme:

| URI | Means |
|---|---|
| `table:<name>` | A user table (most common). |
| `file:<name>.wt` | A raw file (the underlying btree). |
| `colgroup:<table>:<group>` | A column group of a table. |
| `index:<table>:<index>` | An index of a table. |
| `lsm:<name>` | An LSM tree (legacy). |

For convenience, several commands let you pass a bare name (`mytable`) and they default to `table:mytable`. Where a URI is required, prefer the explicit form — it removes ambiguity for `index:` and `colgroup:` cases.

A few special URIs are *not* allowed as `wt` arguments and will be rejected: `backup:`, `config:`, `statistics:` (these only make sense via the C API).

## Choosing the right invocation

When the user describes a goal, walk through these decisions:

1. **Where is the database?** If the user gives you a path, that's `-h <path>`; otherwise confirm the cwd is right.
2. **Add `-r` (read-only) by default.** Drop it only after the user explicitly authorizes a write — and even then, propose the command and wait for confirmation before running anything destructive.
3. **Is it a MongoDB data directory?** If yes, add `-C "log=(enabled=true,path=journal,compressor=snappy)"`.
4. **Did the user name a MongoDB namespace?** (`animals.cats`, `local.oplog.rs`, `mydb.orders`, `system.users`, an index on a collection, etc.) `wt` knows nothing about MongoDB namespaces — it only sees `collection-N-<id>.wt` and `index-N-<id>.wt` files. Before doing anything else, dump `_mdb_catalog.wt` to map the namespace to its on-disk ident, then operate on that file. The recipe is in `references/mongodb-bson.md` ("namespace → ident → documents"). This applies even for read-only diagnostic commands like `verify` and `stat` — point them at the resolved file, not the namespace name.
5. **Does the task read or write?** Read-only tasks pair naturally with `-r`. Write tasks (recovery, salvage, drop, truncate, write, alter, compact, create, load, downgrade) require explicit user authorization first — see "Safety: read-only is the default" above.
6. **Pick the subcommand.** Use `list` to discover URIs first if the user doesn't know them.
7. **Pick subcommand options based on what they need to see** — JSON for piping, hex for binary data, ranges for big tables, checkpoint/timestamp for historical state.

## Common workflows

For end-to-end recipes — inspecting a MongoDB data directory, salvaging a corrupted DB, exporting a table to JSON for analysis, comparing two checkpoints, reading the WAL to find a specific operation — see `references/common-workflows.md`.

## MongoDB BSON values

`wt dump`, `wt verify -d dump_pages`, and `wt printlog` emit MongoDB values as raw BSON bytes. Decode them with this skill's bundled `scripts/wt_to_mdb_bson.py` (a copy of `src/third_party/wiredtiger/tools/wt_to_mdb_bson.py`). Locate it with `find ~/.claude -name "wt_to_mdb_bson.py" 2>/dev/null | head -1`. See `references/mongodb-bson.md` for the catalog-mapping recipe and the snappy-not-builtin gotcha whenever the user is working against a MongoDB dbpath.

**Rule of thumb when a MongoDB namespace shows up:** the first `wt` invocation is *always* a catalog dump to find the file, not the work the user asked for. Users name things like `local.oplog.rs` or `myapp.users`; `wt` sees `collection-18-8170640145721574010.wt`. Dump `_mdb_catalog.wt`, decode with `scripts/wt_to_mdb_bson.py -m dump -j`, filter to the namespace, and read off `.value.ident` (and `.value.idxIdent[<index name>]` for indexes). Only then run the verify / dump / stat / printlog command, against `file:<ident>.wt`. The oplog is the most common case — `local.oplog.rs` lives in some `collection-N-<id>.wt` file you'll never find by name.

## Building `wt`

If the user doesn't have `wt` installed, see `references/building-wt.md` for build paths from a MongoDB checkout (preferred — Bazel target includes snappy as a builtin) or a standalone WiredTiger checkout (CMake/ninja).

## Companion tools

`wt_to_mdb_bson.py` is bundled (see above). Other focused inspection tools live alongside `wt` in `src/third_party/wiredtiger/tools/` in a MongoDB checkout (not bundled): `wt_binary_decode.py` (raw page format), `wt_ckpt_decode.py` (block-checkpoint metadata).

## When to suggest a different tool

`wt` is the right answer for almost any "I want to look at / extract / verify" task on a WiredTiger DB. Suggest something else when:

- Application-level views like query plans or index selection — that's MongoDB-level, not WiredTiger-level. (For just decoding BSON out of `wt dump`/`verify`/`printlog`, use `wt_to_mdb_bson.py` — see `references/mongodb-bson.md`.)
- Live-database changes on a running MongoDB — `wt` can't open a database another process has open. Stop the server first, or use `mongo`/`mongosh`.
- Performance investigation rather than data inspection — point at `wtperf` (`bench/wtperf/`) or `WT_SESSION` statistics through the running app.