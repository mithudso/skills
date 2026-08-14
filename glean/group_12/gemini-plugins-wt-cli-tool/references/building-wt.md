# Building the `wt` tool

If the user doesn't have `wt` installed, the right build path depends on which tree they're in.

## From a MongoDB checkout (preferred for MongoDB users)

The default path on a mongo checkout is the full server build, which produces `wt` alongside `mongod` and the rest of the test/install binaries:

```bash
cd /path/to/mongo
bazel build install-dist-test
# wt lands at bazel-bin/install/bin/wt
./bazel-bin/install/bin/wt -V
```

That's also what most MongoDB devs already have built locally, so on a working dev box `wt` is usually already there — just point at `bazel-bin/install/bin/wt`. The binary already includes snappy (and zstd / zlib / lz4) as builtins, so the MongoDB connection config — `-C "log=(enabled=true,path=journal,compressor=snappy)"` — works without an `extensions=[...]` workaround.

If the user wants only `wt` (skipping the full server build), the narrow target is `//src/third_party/wiredtiger:wt`:

```bash
bazel build //src/third_party/wiredtiger:wt
# binary lands at bazel-bin/src/third_party/wiredtiger/wt
```

This is the same binary as the one ending up in `install/bin/wt`, just without the `mongod` build along for the ride. With a warm remote cache it completes in under a minute.

## From a standalone WiredTiger checkout

If the user is in a standalone WT tree (not a mongo checkout), `wt` is built alongside the rest of WiredTiger via CMake:

```bash
cd <wiredtiger>/build && ninja wt
# binary lands at build/wt
```

A vanilla WT build won't include snappy, which MongoDB needs — see the snappy-extension fallback in `mongodb-bson.md`. (Or use the Bazel target above if a mongo checkout is available.)

## Verifying the build

```bash
<path-to>/wt -V        # prints the WiredTiger library version
<path-to>/wt -?        # prints the help, including the MongoDB config-string hint
```

If `-V` runs cleanly, the binary is good. If MongoDB-side commands later fail with `unknown compressor 'snappy': Invalid argument`, you've got a vanilla WT build without snappy as a builtin — switch to the Bazel target above, or follow the extension-loading workaround in `mongodb-bson.md`.
