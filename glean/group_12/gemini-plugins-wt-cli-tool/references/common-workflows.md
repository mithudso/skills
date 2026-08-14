# Common `wt` workflows

End-to-end recipes for tasks people actually ask about. Each is written so you can adapt it to the user's specifics — paths, table names, encryption keys.

## Inspect a MongoDB data directory

MongoDB writes its WiredTiger files with a different journal layout and snappy compression, so stock `wt` will fail unless you tell it. Use `-r` to keep the directory pristine:

```bash
wt -r \
   -h /path/to/mongo/data \
   -C "log=(enabled=true,path=journal,compressor=snappy)" \
   list -v
```

This prints the metadata for every collection and index. The `wt` help output (`wt -?`) includes the same `log=...` config as a hint for exactly this case.

To dump a specific MongoDB collection (collection files are named `collection-<n>--<id>.wt`):

```bash
wt -r -h /path/to/data -C "log=(...,compressor=snappy)" \
   list                          # find the right URI
wt -r -h /path/to/data -C "log=(...,compressor=snappy)" \
   dump -j file:collection-7--1234567890.wt > coll7.json
```

The values are BSON-encoded — `wt` will give you the raw bytes. Pipe through this skill's bundled decoder (`scripts/wt_to_mdb_bson.py -m dump`; locate it with `find ~/.claude -name "wt_to_mdb_bson.py" 2>/dev/null | head -1`) to render them as JSON. If you also need the namespace → file ident mapping (since collection files are named `collection-N-<id>.wt`, not `db.coll`), see `references/mongodb-bson.md` for the full recipe.

## Diagnose a possibly-corrupt database

Always start read-only. The goal is to gather evidence without changing anything.

```bash
wt -r -h /broken list -v          # what's there?
wt -r -h /broken stat             # any obviously suspicious counters?
wt -r -p -h /broken verify        # full structural check, prefetch off
```

`-p` (global) disables cursor pre-fetching, which is recommended when working with possibly-corrupt data so the engine doesn't speculatively read pages it can't trust.

If `verify` flags a specific table, dig in:

```bash
wt -r -p -h /broken verify -d dump_layout file:bad.wt
wt -r -p -h /broken verify -d dump_pages -k file:bad.wt
```

## Read the WAL for a specific event

`printlog` redacts user data by default. Use `-u` only when the user explicitly wants to see real keys/values.

```bash
wt -h /path printlog                          # redacted, full log
wt -h /path printlog -m                       # message records only
wt -h /path printlog -u                       # show real keys/values
wt -h /path printlog -l 5,0                   # from log file 5 offset 0
wt -h /path printlog -l 5,0,7,4096            # bounded range
```

Output is JSON-ish, so `wt printlog ... | jq` works for filtering.

## Common pitfalls

- **"unable to open... no such file or directory"** — almost always wrong cwd or missing `-h`.
- **"WT_RUN_RECOVERY"** — the database needs recovery. `-R` runs it; `-r` works around it for read-only inspection if recovery isn't strictly necessary; `-S` for salvage.
- **MongoDB DB silently fails to open** — forgot the `-C "log=(...,compressor=snappy)"` config.
- **`read` returns nothing on a binary table** — `read` only handles string/recno keys with string values. Use `dump -k <key>` instead.
