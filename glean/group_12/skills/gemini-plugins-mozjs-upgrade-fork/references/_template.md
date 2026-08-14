# MozJS Upgrade Fork — Repo-Specific Context

Copy this file to your repo as `.agents/skills/mozjs-upgrade-fork/references/<your-repo>.md`
and fill in the values below.

## Current Version

```
VERSION="<major>.<minor>.0esr"
LIB_GIT_BRANCH=spidermonkey-esr<major>.<minor>-cpp-only
LIB_GIT_REVISION=<commit_hash>
```

Located in: `src/third_party/mozjs/scripts/import.sh`

## Fork Branch

The spidermonkey fork lives at `github.com/mongodb-forks/spidermonkey`.
Working branch convention: `spidermonkey-esr<major>.<minor>-cpp-only`

## Known Excluded Tests

JS tests excluded from `make check-jstests` are listed in `references/exclude.txt`.
Update this list when new upstream failures are confirmed unrelated to MongoDB changes.

## Known Permanent Failures

- `./dist/bin/jsapi-tests`: `test_DeflateStringToUTF8Buffer` always fails (SERVER-99489)
- `make check`: `check_vanilla_allocations.py` may fail in non-MongoDB-modified code (upstream issue)
