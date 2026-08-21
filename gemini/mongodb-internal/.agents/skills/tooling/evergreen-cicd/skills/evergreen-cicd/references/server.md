# MongoDB Server — Evergreen CI/CD Context

## Project Flag
```bash
evergreen patch -p mongodb-mongo-master -d "description" -y
```

## Key Aliases
- `required` — required suite (runs on commit queue)
- `commit_queue` — commit queue validation

## Common Task Patterns
- `jscore` / `noPassthrough` — core JS tests
- `replica_sets` — replication tests
- `sharding` — sharding tests
- `concurrency` — concurrency tests
- `auth` — authentication tests

## Local Test Verification
```bash
# Run a specific suite before patching
python buildscripts/resmoke.py run --suite core jstests/core/my_test.js
```

## Patch Workflow
```bash
evergreen patch -p mongodb-mongo-master -d "fix: my change" -y --alias required
```
