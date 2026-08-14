# mozjs-upgrade-integrate

**Category:** Frontend & Web Development
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/query-integration/mozjs-upgrade-integrate/skills/mozjs-upgrade-integrate

## Description
Use when integrating upgraded MozJS/SpiderMonkey sources into the mongo repository. Covers updating import.sh, running extract.sh, generating platform-specific files with gen-config.sh, building, and bumping version metadata. Use when asked to "integrate mozjs", "generate platform files for mozjs", "run gen-config for mozjs", or "finish mozjs integration".

---

# MozJS Upgrade -- Integration (Per-Platform)

## How to Use This File

This file can be used two ways:

- **Via the skill system** (marketplace): invoke the `mozjs-upgrade-integrate` skill.
- **Directly** (no setup needed): tell Claude
  `Read .agents/skills/mozjs-upgrade-integrate/SKILL.md and perform the integration for <platform>. Branch is <branch>.`
  Claude will read this file and follow the steps. For platform-specific details, also read the
  relevant reference file under `.agents/skills/mozjs-upgrade-integrate/references/`.

When used directly, Claude should confirm which platform it is running on before starting, check
that the branch is checked out and up to date (`git pull`), and proceed step by step — reporting any
hiccup before doing anything destructive.

## Overview

Integration brings the prepared SpiderMonkey fork into the mongo repo. This must be done on each of
7 platforms. The gen-config.sh step is platform-specific and must run natively (or under Rosetta for
macOS x86_64).

## Platform Matrix

| Platform       | Architecture | gen-config args  | Build command                                                   | Notes                           |
| -------------- | ------------ | ---------------- | --------------------------------------------------------------- | ------------------------------- |
| Linux x86_64   | x86_64       | `x86_64 linux`   | `bazel build install-devcore`                                   | Primary development platform    |
| Linux arm64    | aarch64      | `aarch64 linux`  | `bazel build install-devcore`                                   |                                 |
| macOS arm64    | aarch64      | `aarch64 macOS`  | `bazel build install-devcore`                                   | See platform-macos.md           |
| macOS x86_64   | x86_64       | `x86_64 macOS`   | `bazel build install-devcore`                                   | Use Rosetta: `arch -x86_64 zsh` |
| Windows x86_64 | x86_64       | `x86_64 windows` | `bazel build install-devcore`                                   | See platform-windows.md         |
| Linux ppc64le  | ppc64le      | `ppc64le linux`  | `python bazel/bazelisk.py build --config=local install-devcore` | See platform-ppc-s390x.md       |
| Linux s390x    | s390x        | `s390x linux`    | `python bazel/bazelisk.py build --config=local install-devcore` | See platform-ppc-s390x.md       |

## Step 1: Update scripts/import.sh

Edit `src/third_party/mozjs/scripts/import.sh`:

```bash
VERSION="{MAJOR}.{MINOR}.0esr"
LIB_GIT_BRANCH=spidermonkey-esr{MAJOR}.{MINOR}-cpp-only
LIB_GIT_REVISION=<commit_hash_from_fork>
```

## Step 2: Pull and Extract Sources

`import.sh` only needs git — no python venv or LLVM_OBJDUMP needed yet.

```bash
cd src/third_party/mozjs
bash scripts/import.sh      # clones the fork; also removes any existing mozilla-release/
```

Then set the required env vars before running `extract.sh`:

```bash
export LLVM_OBJDUMP=/opt/mongodbtoolchain/v5/bin/llvm-objdump   # Linux; do NOT add full v5 bin to PATH
source python3-venv/bin/activate                                  # for buildscripts used by gen-config
bash scripts/extract.sh
```

## Step 3: Generate Platform-Specific Files

Before spawning integration hosts, push the mongo branch to origin so remote hosts can fetch it:

```bash
git push origin <your-branch>
```

Then run gen-config for the current platform:

```bash
cd src/third_party/mozjs
bash scripts/gen-config.sh <arch> <os>
```

Each platform produces files in `platform/<arch>/<os>/`.

After gen-config, drop the spurious diff in `selfhosted.out.h` (3-byte change in the bytecode blob;
not needed by the build):

```bash
git checkout -- src/third_party/mozjs/extract/js/src/selfhosted.out.h
```

For platform-specific modifications required before running gen-config.sh, see the reference files.
See `references/spawn-integration-hosts.md` for how to provision Linux x86_64 and Windows hosts.

## Step 4: Verify the Build

Moving `mozilla-release/` out of the tree is **non-optional** — leaving it in place causes
duplicate-source errors and breaks the build.

```bash
mv src/third_party/mozjs/mozilla-release/ ../
bazel build install-devcore
# On ppc/s390x: python bazel/bazelisk.py build --config=local install-devcore
```

Expected output (cold build with EngFlow remote cache): ~9700 actions, ~54s wall time. Without
remote cache (local only): significantly longer.

## Step 5: Commit Per Platform

Commit only the platform files and push:

```bash
git add src/third_party/mozjs/platform/<arch>/<os>/
git commit -m "SERVER-XXXXX MozJS ESR {MAJOR}.{MINOR}.{PATCH} platform config: <Platform>"
git push origin <your-branch>
```


## Step 6: Bump Version Metadata

Done once (not per-platform), on the primary workstation. Update the version string in:

1. `README.third_party.md`
2. `sbom.json`

Search for the old version to find all occurrences: `git grep <old_version>`

Then run: `bazel run //:format`

## Repo-Specific Context

Before starting, check for a repo-specific context file in `references/`. Use
[references/_template.md](references/_template.md) as a template if none exists for your repo.

## Platform-Specific References

- macOS: `references/platform-macos.md`
- Windows: `references/platform-windows.md`
- PPC/s390x: `references/platform-ppc-s390x.md`
- Host setup: `references/host-setup-quickstart.md`
- Spawning hosts: [references/spawn-integration-hosts.md](references/spawn-integration-hosts.md) (script: [references/spawn-integration-hosts.sh](references/spawn-integration-hosts.sh))