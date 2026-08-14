# Decoding MongoDB BSON values from `wt` output

When the user is poking at a MongoDB data directory, the values inside every `wt dump`, `wt verify -d dump_pages`, and `wt printlog` are **BSON-encoded** — opaque bytes from `wt`'s point of view. This skill bundles a small Python tool that decodes those bytes into something a human (or `jq`) can read:

```
scripts/wt_to_mdb_bson.py        (relative to this skill's root directory)
```

The install path varies by plugin configuration. The recipes below assume `$BSON` is set to the absolute path — locate it first:

```bash
BSON=$(find ~/.claude -name "wt_to_mdb_bson.py" 2>/dev/null | head -1)
```

(The bundled copy is a snapshot of `src/third_party/wiredtiger/tools/wt_to_mdb_bson.py` from the MongoDB tree. If you need the latest upstream version — e.g. when MongoDB adds new BSON types — copy from a current checkout into `scripts/`.)

This reference covers when and how to use it, plus the catalog-mapping recipe that turns "I want to see `db.coll`" into "dump this `collection-N-<id>.wt` file."

## When to reach for it

If the user asks any of these, the answer involves `scripts/wt_to_mdb_bson.py`:

- "What documents are in this collection on disk?"
- "Get me the BSON of every record in `<ns>` from this dbpath."
- "What's the on-disk layout of the catalog / what's the ident for `<ns>`?"
- "Decode the WAL so I can see the actual writes that happened."
- "When I run `verify -d dump_pages` against a corrupt MongoDB collection, the V {...} lines are unreadable — render them."

If the user is on a non-MongoDB WiredTiger DB (e.g. a `wtperf` corpus, a standalone application's WT files), this script does not apply — values aren't BSON.

## The tool

Three modes — one per `wt` subcommand whose output it knows how to parse:

| Mode | Decodes output of | Notes |
|---|---|---|
| `-m dump`     | `wt dump -x <uri>` | Hex format is required (`-x`); the script `codecs.decode(value, 'hex')`s each value. |
| `-m verify`   | `wt verify -d dump_pages <uri>` | Looks for `V {...}` value lines and adds a decoded line beneath each. |
| `-m printlog` | `wt printlog -u -x` | Looks for `"value-hex": "..."` fields and decodes them. `-u` un-redacts user data; `-x` adds hex. Skip `-x` and there's nothing to decode. |

Two invocation styles:

```bash
# Pipe — preferred. You see exactly what wt is doing and can swap flags freely.
./wt -r -h /data -C "..." dump -x file:foo.wt | "$BSON" -m dump

# Self-invoking — convenient one-liner. Script invokes wt for you.
# Only useful for printlog mode, where the wt config is hardcoded inside.
# In dump/verify mode you still need to be in the right cwd, so the pipe form is usually clearer.
"$BSON" -m dump     -f /path/to/wt file:foo.wt
"$BSON" -m verify   -f /path/to/wt file:foo.wt
"$BSON" -m printlog -f /path/to/wt
```

Useful flag on either style:

- `-j` / `--json` — emit Canonical Extended JSON (one document per line for `dump` mode). Pipeable to `jq`. Without `-j`, output is a friendly Python-pretty form with `Key:` / `Value:` sections.

The script imports `bson` from PyMongo. On the MongoDB toolchain, `python3` already has it; otherwise `pip install pymongo`.

## Snappy: stock `wt` won't open MongoDB files

MongoDB writes WiredTiger files with `block_compressor=snappy` and a snappy-compressed journal. A `wt` binary that wasn't compiled with `--with-builtins=snappy` (e.g. a vanilla checkout's default build) will refuse the connection:

```
unknown compressor 'snappy': Invalid argument
```

Two fixes:

1. **Load the snappy extension at runtime** — every WT build also produces `libwiredtiger_snappy.so` under `<wt>/build/ext/compressors/snappy/`. Reference it in the `-C` config:

   ```bash
   EXT=$(find /path/to/wt-source/build/ext -name 'libwiredtiger_snappy.so' | head -1)
   wt -r -h /data \
      -C "extensions=[$EXT],log=(enabled=true,path=journal,compressor=snappy)" \
      list
   ```

   This is the path of least resistance when the user has a checkout but not a MongoDB build.

2. **Use a `wt` built from MongoDB's WT source** — `src/third_party/wiredtiger/` produces a `wt` with snappy already builtin. That's what `mongo/build/install/bin/wt` is.

Once snappy is loaded the rest of the workflow below works unchanged.

## Recipe: namespace → ident → documents

MongoDB stores collections in files named `collection-<n>-<id>.wt` and indexes in `index-<n>-<id>.wt`. The mapping from a logical namespace (`animals.cats`, `local.oplog.rs`, `myapp.users`, etc.) to those files lives in `_mdb_catalog.wt`. **Whenever the user names a collection or index by namespace, the first `wt` step is always this catalog dump — `wt` itself only knows the ident.** Three steps:

```bash
HOME=/data
EXT=/path/to/libwiredtiger_snappy.so
CFG="extensions=[$EXT],log=(enabled=true,path=journal,compressor=snappy)"
BSON=$(find ~/.claude -name "wt_to_mdb_bson.py" 2>/dev/null | head -1)

# 1. Dump and decode the catalog. Each record maps a namespace to its idents.
wt -r -h "$HOME" -C "$CFG" dump -x file:_mdb_catalog.wt \
  | "$BSON" -m dump -j \
  | jq 'select(.value.md.ns == "animals.cats")'
# .value.ident       → the collection ident, e.g. "collection-7-884086260529450731"
# .value.idxIdent    → object: index name → index ident

# 2. Dump that collection's file. Values come out as BSON documents.
wt -r -h "$HOME" -C "$CFG" dump -x file:collection-7-884086260529450731.wt \
  | "$BSON" -m dump -j
```

For "list every namespace and its file":

```bash
wt -r -h "$HOME" -C "$CFG" dump -x file:_mdb_catalog.wt \
  | "$BSON" -m dump -j \
  | jq -r '.value.md.ns + "\t" + .value.ident'
```

The catalog also tells you the UUID, capped settings, indexes (with their full spec), and any TTL config — useful for "what was this collection's shape before it crashed?"

## Recipe: BSON-decoded page-level verify

When `verify` flags a collection, the page dump shows raw bytes. Pipe it through to read what's actually there:

```bash
wt -r -p -h "$HOME" -C "$CFG" verify -d dump_pages file:collection-7-...wt \
  | "$BSON" -m verify
```

The script preserves the verify output verbatim and inserts a decoded line under each `V {...}` it can decode. Use `-p` (the global-level prefetch-disable flag) to avoid speculative reads on suspect data.

## Recipe: BSON in the WAL

Journal records carry their values as `"value-hex": "<hex>"` fields once `printlog` is given `-u -x`:

```bash
wt -r -h "$HOME" -C "$CFG" printlog -u -x \
  | "$BSON" -m printlog
```

Note: keys are left as hex — for `_mdb_catalog` and the oplog they're either small integers or special record-id encodings, not BSON.

## What about index files?

`index-N-<id>.wt` records a server-encoded sort key plus a record-id pointing back into the collection — they're not BSON documents. `$BSON` (`wt_to_mdb_bson.py`) won't usefully decode them. For ad-hoc inspection use `wt dump -px` (printable + hex) and read the bytes; for a programmatic decode you need MongoDB's KeyString library (or `mongod --queryableBackupMode` against a copy of the dbpath).

## Companion tools

`wt_to_mdb_bson.py` is bundled in this skill at `scripts/wt_to_mdb_bson.py`. The other tools listed below live under `src/third_party/wiredtiger/tools/` in a MongoDB checkout — they are *not* bundled here, so they require a checkout. When the user's question matches one, point them at it instead of writing your own decoder:

| Tool | What it does |
|---|---|
| `scripts/wt_to_mdb_bson.py` (bundled) | Decode BSON in `dump` / `verify` / `printlog` output (this file). |
| `wt_binary_decode.py` | Decode raw WT page-format bytes (lengths, cell types, MVCC tuples). |
| `wt_ckpt_decode.py` | Inspect `WT_BLOCK` checkpoint metadata (file offsets, root addr, etc.). |
| `wt_cmp_dir` | Compare two WT home dirs at the file level (size, checksum, mtime). |
| `wt_cmp_uri.py` | Compare a single URI's contents between two homes. |
| `wt_disagg_addr_decode.py` | Decode disaggregated-storage page addresses. |
| `wt_timestamps` | Pretty-print MVCC start/stop/durable timestamps. |
| `wt_turtle_config_parse.py` | Parse the `WiredTiger.turtle` metadata file. |
| `wt_verify/` | Helpers for richer corruption analysis on top of `verify`. |
| `backup_analysis.py` | Analyze WT hot-backup outputs. |

These tools are coupled to the WiredTiger source they ship alongside — when there's any version skew suspected, run them out of the same checkout as the `wt` binary you used to produce the input.
