# MMS — PR Review Loop Context

## GitHub Repository
- Owner/Repo: `10gen/mms`

## Build Variants
- `code_health` — Linting, static analysis
- `unit_java` — Java unit tests
- `int` — Integration tests
- `js` — JavaScript tests
- `local_openapi_required` — OpenAPI validation
- `e2e_local_required` — End-to-end tests
- `bazel_linux_x86_64` — Bazel build verification

## Local Test Commands
```bash
# Java unit tests (scope to affected packages when possible)
bazel test //server/src/unit/...

# Frontend unit tests
bazel test //client/packages/...:js_unit_test

# Code quality checks
bazel build //server:mms
```

## Ticket Prefix
Format: `CLOUDP-######` (six digits)

## AI Code Review Bot
- Bot name: `augmentcode[bot]`
- Trigger phrase: `augment review`
- Enabled: Yes

## Commit Message Format
```
CLOUDP-XXXXXX: Short description

- List of changes
```

## Evergreen Project
- Project: `mms`
- Patches URL: https://evergreen.mongodb.com/patches/mine (filter by `mms`)
