# MozJS Upgrade Integrate — Repo-Specific Context

Copy this file to your repo as `.agents/skills/mozjs-upgrade-integrate/references/<your-repo>.md`
and fill in the values below.

## Repository Layout

- MozJS sources: `src/third_party/mozjs/`
- import.sh / extract.sh / gen-config.sh: `src/third_party/mozjs/scripts/`
- Platform output files: `src/third_party/mozjs/platform/<arch>/<os>/`

## Build Command

```bash
# Standard platforms (Linux x86_64, arm64, macOS, Windows):
bazel build install-devcore

# PPC and s390x only:
python bazel/bazelisk.py build --config=local install-devcore
```

## LLVM_OBJDUMP Path

```bash
export LLVM_OBJDUMP=/opt/mongodbtoolchain/v5/bin/llvm-objdump
```

## Remote Cache

EngFlow remote cache is available on standard hosts. Expected cold build: ~9700 actions, ~54s wall time.

## Version Metadata Files

Files to update when bumping the version string:
- `README.third_party.md`
- `sbom.json`
