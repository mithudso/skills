# Windows Platform Guide (x86_64)

## Overview

Windows integration uses the **MozillaBuild shell** for mozjs script steps and **PowerShell** for
bazel/git steps. The repo must live on the **Z: drive** (`Z:\mongo`) — the MozillaBuild shell maps
it as `/z/mongo`, and bazel performs much better on Z: than on C:.

> **RDP vs SSH**: RDP is recommended. SSH through Cygwin sshd works but has known limitations
> (Python/shell prompts don't display correctly, path mismatches between Cygwin and Windows). These
> instructions document both paths where they differ.

---

## 1. Spawn the Host

Spawn a `windows-vsCurrent-xxlarge` host from [Evergreen](https://evergreen.mongodb.com/spawn). For
MongoDB 4.6+, `windows-2022-large` also works for SSH-only access.

---

## 2. Initial Machine Setup (from your local machine)

```bash
REMOTE_HOST=<YOUR HOST DNS NAME>

# Set the Administrator password (needed for RDP)
ssh Administrator@$REMOTE_HOST "net user Administrator <PASSWORD>"

# Copy SSH key so git clone works on the host
ssh Administrator@$REMOTE_HOST mkdir /cygdrive/c/Users/Administrator/.ssh
scp $HOME/.ssh/id_rsa Administrator@$REMOTE_HOST:/cygdrive/c/Users/Administrator/.ssh/id_rsa

# Copy .gitconfig so commits have the right name/email
scp $HOME/.gitconfig Administrator@$REMOTE_HOST:/cygdrive/c/Users/Administrator/.gitconfig

# Copy .evergreen.yml for patch builds
scp $HOME/.evergreen.yml Administrator@$REMOTE_HOST:/cygdrive/c/Users/Administrator/.evergreen.yml
```

Add `-o StrictHostKeyChecking=no` to each `ssh`/`scp` command to skip the host key prompt.

---

## 3. Clone the Repo to Z:

### Via RDP (PowerShell or cmd.exe)

```powershell
Z:
git clone git@github.com:10gen/mongo.git
cd mongo
```

### Via SSH (Cygwin bash)

Git on Windows (MinGW) cannot use Cygwin's forwarded SSH agent. The spawn setup above copies
`id_rsa` directly to the host, which git uses automatically. If you skipped that step, use the
`cygpath -m` workaround instead:

```bash
# Only needed if id_rsa was NOT copied in step 2
git config --global core.sshCommand "$(cygpath -m "$(which ssh)")"
```

Then clone to Z:

```bash
cd /cygdrive/z
git clone git@github.com:10gen/mongo.git
cd mongo
git fetch origin <branch> && git checkout <branch>
```

> **64-bit git is required** — `scripts/import.sh` makes large git requests that 32-bit Cygwin git
> cannot handle. Ensure the chocolatey version is on PATH:
>
> ```bash
> choco install git -y
> export PATH="/c/Program Files/Git/cmd:$PATH"
> ```

---

## 4. Configure Bazel (Critical — C: will fill up)

The C: drive fills up with a single build. Redirect everything to Z: **before** running bazel for
the first time.

Create `C:\Users\Administrator\.bazelrc` (or `.bazelrc.local` in the repo root) and add at the top:

```
startup --output_user_root=Z:/b
startup --output_base=Z:/b/b
common --action_env=TMP=Z:/b
common --action_env=TEMP=Z:/b
```

Then create the directory:

```powershell
mkdir Z:\b
```

---

## 5. Install Prerequisites

From **Cygwin bash**:

```bash
choco install make -y
choco install git -y
choco install llvm -y   # needed for clang-cl during MozJS configure
```

**Rust**: `windows-vsCurrent-xxlarge` hosts have Rust pre-installed at
`C:\Users\Administrator\.cargo\`. Verify: `ls /cygdrive/c/Users/Administrator/.cargo/bin/cargo`. If
missing, install:

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --default-toolchain 1.93.1 --profile minimal
```

> `~/.cargo/env` does **not** exist in Cygwin home on Windows hosts — Rust is installed to the
> Windows user profile. Always use the explicit path:
>
> ```bash
> export PATH="/cygdrive/c/Users/Administrator/.cargo/bin:$PATH"
> ```

Install cbindgen:

```bash
export PATH="/cygdrive/c/Users/Administrator/.cargo/bin:$PATH"
cargo install cbindgen
```

---

## 6. Install MozillaBuild

MozillaBuild must be downloaded to a Windows-accessible path and run via `cmd /c` — executing `.exe`
files directly from Cygwin `/tmp` fails with "Permission denied":

```bash
# In Cygwin bash
curl -L https://ftp.mozilla.org/pub/mozilla/libraries/win32/MozillaBuildSetup-Latest.exe \
  -o /cygdrive/c/Users/Administrator/mozillabuild.exe
chmod +x /cygdrive/c/Users/Administrator/mozillabuild.exe
cmd /c "C:\\Users\\Administrator\\mozillabuild.exe /S"
sleep 30
# Default install path: C:\mozilla-build
```

---

## 7. Install Node.js and Claude Code (SSH workflow)

The pre-installed Node on Evergreen Windows hosts is v8, which is too old. Install LTS via choco:

```bash
choco install nodejs-lts -y
export PATH="/cygdrive/c/Program Files/nodejs:$PATH"
```

npm on Windows has an `${APPDATA}` prefix bug in Cygwin. Set the prefix explicitly before installing
Claude:

```bash
npm config set prefix "C:/Users/Administrator/npm-global"
npm install -g @anthropic-ai/claude-code
```

Add all new tools to PATH and persist in `~/.bashrc`:

```bash
echo 'export PATH="/cygdrive/c/Users/Administrator/npm-global:$PATH"' >> ~/.bashrc
echo 'export PATH="/cygdrive/c/Users/Administrator/.cargo/bin:$PATH"' >> ~/.bashrc
echo 'export PATH="/cygdrive/c/Program Files/nodejs:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

Verify: `claude --version`

> If claude is not found after install, locate it with:
>
> ```bash
> find /cygdrive/c /cygdrive/z -name "claude" 2>/dev/null
> ```
>
> Then add that directory to PATH.

---

## 8. Running Bazel from Cygwin (Critical Setup)

Bazel has several Cygwin-specific pitfalls. Set these environment variables in `~/.bashrc`
**before** running any bazel command:

```bash
export USERPROFILE='C:\Users\Administrator'
export HOMEDRIVE='C:'
export HOMEPATH='\Users\Administrator'
export BAZELISK_HOME='C:\Users\Administrator\AppData\Local\bazelisk'
```

Without `USERPROFILE`/`HOMEDRIVE`/`HOMEPATH`, bazel-managed Python tools call
`Path("~").expanduser()` and crash with `RuntimeError: Could not determine home directory`.

Without `BAZELISK_HOME`, `$LOCALAPPDATA` (which Cygwin sets to a posix path like `/data/tmp`) causes
bazelisk to cache downloads in a broken location.

**Bazelisk picks the wrong launcher on Windows.** `bazel/bazelisk.py` finds `tools/bazel` (a bash
script), sees it's executable under Cygwin, and tries to run it — but Windows can't `CreateProcess`
a bash script (`WinError 193`). The fix is to call the downloaded `bazel.exe` directly:

```bash
BAZEL_EXE=$(find "$BAZELISK_HOME/downloads/bazelbuild" -name "bazel*.exe" 2>/dev/null | head -1)
"$BAZEL_EXE" build install-devcore
```

Or from PowerShell (avoids all Cygwin path issues entirely — recommended for bazel steps):

```powershell
$env:Path += ";C:\Users\Administrator\.local\bin;"
bazel build install-devcore
```

**Python version**: Cygwin PATH may put Python 2.7 (`/cygdrive/c/Python27`) ahead of 3.14. Always
use `python3` or a full path when invoking Python scripts manually.

**MSVC version**: `.bazelrc` pins `BAZEL_VC_FULL_VERSION=14.44.35207`. Spawn hosts created before
**2026-01-16** ship MSVC 14.31 and fail with
`Auto-Configuration Error: Microsoft Visual C++ build tools 14.44.35207 could not be found`. Fix:
spawn a fresh host (base image has 14.44), or install via VS Installer:
`vs_installer.exe modify --add Microsoft.VisualStudio.Component.VC.14.44.17.14.x86.x64`

---

## 10. Build Zlib

The moz.configure zlib fix (step 12a) needs the bazel-built zlib library. Build it first.

From **PowerShell** in `Z:\mongo`:

```powershell
# Install bazel if needed
C:\Python\Python310\python.exe buildscripts/install_bazel.py
$env:Path += ";C:\Users\Administrator\.local\bin;"

# Build (this triggers EngFlow auth — copy the URL and open it on your local machine)
bazel build install-devcore
```

The zlib output will be at:

```
Z:\mongo\bazel-out\x64_windows-fastbuild\bin\src\third_party\zlib\zlib
```

---

## 11. Pull MozJS Sources

Launch the **MozillaBuild shell** and keep it open for steps 10–12:

```powershell
C:\mozilla-build\start-shell.bat
```

```bash
# In MozillaBuild shell
cd /z/mongo/src/third_party/mozjs
git checkout <your-mozjs-upgrade-branch>
rm -rf extract include mozilla-release

# scripts/import.sh uses /tmp by default; on Windows change it to /z to avoid slow cross-drive copies
# Edit scripts/import.sh and change:
#   LIB_GIT_DIR=$(mktemp -d /tmp/import-spidermonkey.XXXXXX)
# to:
#   LIB_GIT_DIR=$(mktemp -d /z/import-spidermonkey.XXXXXX)
# (The file already has the Windows line commented out for reference)

bash scripts/import.sh
```

> **OOM on import.sh?** Make sure PATH uses the 64-bit chocolatey git, not the Cygwin git.

---

## 12. Required Modifications Before extract.sh

Apply all of the following to files inside `mozilla-release/` **before** running `extract.sh`. Do
**not** commit these changes.

### 10a. moz.configure (~line 918): Hardcode zlib paths

`pkg_check_modules` does not work on Windows. Comment it out and hardcode the paths:

```diff
-pkg_check_modules("MOZ_ZLIB", "zlib >= 1.2.3", when="--with-system-zlib")
+# pkg_check_modules("MOZ_ZLIB", "zlib >= 1.2.3", when="--with-system-zlib")
 set_config("MOZ_SYSTEM_ZLIB", True, when="--with-system-zlib")
-
+set_config("MOZ_ZLIB_CFLAGS", ["-IZ:/mongo/src/third_party/zlib"], when="--with-system-zlib")
+set_config("MOZ_ZLIB_LIBS", ["Z:/mongo/bazel-out/x64_windows-fastbuild/bin/src/third_party/zlib/zlib"], when="--with-system-zlib")
```

### 10b. build/moz.configure/pkg.configure (~line 15): Remove WINNT exclusion

```diff
-return compile_environment and target.os not in ("WINNT", "OSX", "Android")
+return compile_environment and target.os not in ("OSX", "Android")
```

### 10c. ICU double-conversion fix

Without system ICU on Windows, three files in `intl/icu/source/i18n/` fail to compile. Apply this
patch from inside `mozilla-release/`:

```diff
# intl/icu/source/i18n/number_decimalquantity.cpp  (two occurrences)
-#ifdef JS_HAS_INTL_API
+#if JS_HAS_INTL_API || !MOZ_SYSTEM_ICU

# intl/icu/source/i18n/number_utils.cpp  (two occurrences)
-#ifdef JS_HAS_INTL_API
+#if JS_HAS_INTL_API || !MOZ_SYSTEM_ICU

# intl/icu/source/i18n/units_converter.cpp  (two occurrences)
-#ifdef JS_HAS_INTL_API
+#if JS_HAS_INTL_API || !MOZ_SYSTEM_ICU
```

Save as `icu_fix.patch` and apply:

```bash
cd mozilla-release
patch -p1 < /z/icu_fix.patch
cd ..
```

---

## 13. Run extract.sh and gen-config.sh

```bash
# In MozillaBuild shell, from src/third_party/mozjs/
bash scripts/extract.sh > output.txt
# Pipe to output.txt so only warnings/errors appear in the shell.
```

**Common failures:**

| Symptom                    | Fix                                                 |
| -------------------------- | --------------------------------------------------- |
| `msys2 make` or wrong make | `export PATH="/c/ProgramData/chocolatey/bin:$PATH"` |
| C compiler not found       | `export PATH="$PATH:/c/Program Files/LLVM/bin/"`    |
| `extract/` folder is empty | `extract.sh` failed — check `output.txt` for errors |
| Python script errors       | See Python fix below                                |

**Python version fix** (if configure scripts fail with Python errors):

In `mozilla-release/js/src/configure`, change:

```bash
PYTHON3="${PYTHON3:-python3}"
# to:
PYTHON3="/c/Python/Python310/python.exe"
```

Then run gen-config:

```bash
bash scripts/gen-config.sh x86_64 windows
```

---

## 14. Post gen-config Cleanup

From Cygwin bash (not MozillaBuild):

```bash
cd /cygdrive/z/mongo
# Revert the spurious 3-byte change in selfhosted.out.h
git checkout -- src/third_party/mozjs/extract/js/src/selfhosted.out.h
```

---

## 15. Verify the Build

Move `mozilla-release/` out of the tree first (leaving it causes duplicate-source errors):

```bash
mv src/third_party/mozjs/mozilla-release/ ../
```

Then build (PowerShell or MozillaBuild shell):

```bash
bazel build install-devcore
```

If bazel is not found: `export PATH="$PATH:/c/Users/Administrator/.local/bin"` (MozillaBuild) or
`$env:Path += ";C:\Users\Administrator\.local\bin;"` (PowerShell).

---

## 16. Commit

```bash
git add src/third_party/mozjs/platform/x86_64/windows/
git commit -m "SKUNK-11 MozJS ESR {MAJOR}.{MINOR}.{PATCH} platform config: Windows x86_64"
git push origin <branch>
```

If files in `extract/` or `include/` are also modified, the pulled MozJS version differs from the
base commit — investigate before committing.

---

## Troubleshooting

- **All script steps must run in the MozillaBuild shell** — not PowerShell, not Cygwin bash.
- **C: drive fills up**: Set `.bazelrc` `output_user_root` to Z: before the first build (step 4). To
  recover: delete `C:\Users\Administrator\_bazel_Administrator\<cache>`, run
  `bazel clean --expunge`, configure `.bazelrc`, retry.
- **Bazel not found**: Re-add to PATH: `$env:Path += ";C:\Users\Administrator\.local\bin;"`
- **git not 64-bit**: `choco install git -y` then ensure it's on PATH before Cygwin's git.
- **EngFlow auth**: When bazel first runs it prints a URL. Open it on your **local** machine (not
  the host). If the shell looks stuck, press Enter.
- **SSH git clone fails**: Ensure `id_rsa` was copied to the host in step 2. Alternatively use
  `git config --global core.sshCommand "$(cygpath -m "$(which ssh)")"`.
- **Claude not found after install**: Run `find /cygdrive/c /cygdrive/z -name "claude" 2>/dev/null`
  to locate it, then add that directory to PATH in `~/.bashrc`.

---

## Historical: MSVC Compilation Errors (ESR 115.7)

These were fixed in patches already committed to the repo. Listed here in case they resurface:

- Atomic ops assembly (MSVC vs GCC/Clang) → Windows-specific intrinsics
- JIT flag not propagated → explicit Windows disable
- `JSProtoKey` trailing comma → `/Zc:preprocessor` flag
- Variadic macro naming extension → MSVC-compatible rewrite
- `__atomic_store_n` → `std::atomic_store_explicit`
- `RecGroup` zero-length `TypeDef` array → empty initializer list
- `__builtin_frame_address` → `_AddressOfReturnAddress`
- `WeakHeapPtr<BaseShape*>` missing GCPolicy header → add include
- `SetMallocMaxDirtyPageModifier` undefined → delete unused function
- `PRMJ_now()` declared `static` → remove static
- `emitInstanceCall0` zero-length array init → use `nullptr`
- SSE intrinsics not supported → disable relevant optimizations
