# macOS Platform Guide (arm64 and x86_64)

## Before You Start

**Verify `mozilla-release/` is at the correct ESR version** before running gen-config. A stale
directory from a previous upgrade cycle will silently produce wrong output (the Rust crate's
`Cargo.toml` will be missing, configure will detect the wrong version, etc.).

```bash
cat mozilla-release/config/milestone.txt   # should match target ESR, e.g. 140.10.2
```

If it doesn't match, run `import.sh` first to re-fetch the correct sources.

**Python venv path varies by engineer.** SKILL.md refers to `python3-venv/bin/activate` but there is
no standard path on macOS laptops — find your own:

```bash
find ~/mongo -maxdepth 3 -name activate 2>/dev/null | head -5
```

Use whatever path `activate` resolves to in all commands below.

## Dependencies

Install the required packages:

```bash
brew install zlib lld
```

## Required mozilla-release/ Modifications Before gen-config.sh

The following modifications must be applied to files inside `mozilla-release/` before running
`gen-config.sh` on macOS.

### 1. moz.configure (~line 926): Replace pkg_check_modules with hardcoded zlib config

Replace the `pkg_check_modules` call for zlib with hardcoded paths:

```python
# pkg_check_modules("MOZ_ZLIB", "zlib >= 1.2.3", when="--with-system-zlib")
set_config("MOZ_SYSTEM_ZLIB", True, when="--with-system-zlib")
set_config("MOZ_ZLIB_CFLAGS", ["-I/opt/homebrew/opt/zlib/include"], when="--with-system-zlib")
set_config("MOZ_ZLIB_LIBS", ["z"], when="--with-system-zlib")
```

### 2. build/moz.configure/flags.configure: Remove -flegacy-pass-manager lines

Remove the following lines (inside the `clang >= 13.0.0` block):

```python
# Remove these lines:
if compiler.type == "clang":
    return namespace(flags=["-flegacy-pass-manager"], enabled=False)
if compiler.type == "clang-cl":
    return namespace(flags=["-Xclang", "-flegacy-pass-manager"], enabled=False)
```

### 3. build/moz.configure/toolchain.configure (~line 2003)

Update the version check flags (inside `# Check the kind of linker`):

```python
version_check = ["-Wl", "--version"]
```

## Required Environment Variables for gen-config.sh

Set these before running `gen-config.sh` (in addition to the standard ones in SKILL.md):

```bash
export CXXFLAGS="-std=c++17"
export MACOSX_DEPLOYMENT_TARGET=11.0
```

These prevent configure from picking up an excessively new deployment target (e.g. `26.5` on macOS
Tahoe) and ensure C++17 is used, avoiding `std::is_array_v` errors from `mozilla/Assertions.h`.

## macOS x86_64 via Rosetta

Run gen-config under Rosetta with the same env vars as arm64:

```bash
arch -x86_64 /bin/bash -c '
export PATH="$(brew --prefix llvm)/bin:$PATH"
export LLVM_OBJDUMP="$(brew --prefix llvm)/bin/llvm-objdump"
export CXXFLAGS="-std=c++17"
export MACOSX_DEPLOYMENT_TARGET=11.0
source /path/to/mongo/venv/bin/activate
bash scripts/gen-config.sh x86_64 macOS
'
```

Note: the MongoDB toolchain (`/opt/mongodbtoolchain`) does not exist on macOS.

**Reusing `mozilla-release/` between arm64 and x86_64 runs.** The modifications applied in the
section above persist in the directory. After the arm64 run, move `mozilla-release/` out of the tree
for the build, then move it back in for the x86_64 run — no need to re-apply the modifications or
re-run `import.sh`.

## Building on macOS

Use `--config=local` to bypass the EngFlow remote cache, which may not be reachable from a laptop
and will cause the build to hang indefinitely:

```bash
bazel build install-devcore --config=local
```
