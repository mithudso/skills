# MMS — Graphite PR Stack Context

## Build Verification
Run before submitting stack to ensure CI won't immediately fail:
```bash
bazel build //server:mms
```

## Pre-commit Hooks
MMS pre-commit hooks run on `gt submit`. Expected output includes Java compile check
and clang-format. Hook failures block submission — fix before re-running `gt submit`.

## PR Title Format
Prefix with Jira ticket: `CLOUDP-123456: description of change`
