# MMS — Evergreen CI/CD Context

## Project Flag
```bash
evergreen patch -p mms -d "description" -y
```

## Key Variants
- `unit_java` — Java unit tests
- `int` — Integration tests
- `code_health` — Linting, static analysis

## Common Bazel Tasks
```bash
# Run unit tests locally before patching
bazel test //server/...

# Check build health
bazel build //server:mms
```

## Patch Workflow
```bash
evergreen patch -p mms -d "fix: my change" -y --variants unit_java,int
```

## Monitoring
Check patch status at: https://evergreen.mongodb.com/patches/mine
Filter by project: `mms`
