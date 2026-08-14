# MongoDB Server — Graphite PR Stack Context

## Build Verification
Run a targeted test before submitting stack:
```bash
python buildscripts/resmoke.py run --suite core jstests/core/<your_test>.js
```

## Pre-commit Hooks
Server pre-commit hooks run clang-format and eslint. Fix all hook failures before
running `gt submit`. To run manually: `pre-commit run --all-files`

## PR Title Format
Prefix with Jira ticket: `SERVER-XXXXX: description of change`

## Draft Mode
Always open Server PRs as draft. Graphite does this automatically when using
`gt submit --draft`.
