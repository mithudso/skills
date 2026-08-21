# Build and Test Reference

This reference covers building SpiderMonkey from the fork and running the test suites to verify that
the cherry-picked commits produce a working build.

## Install Dependencies

### Mozilla Bootstrap (optional but recommended)

Mozilla provides a bootstrap script that installs system-level build dependencies:

```bash
curl -L https://hg.mozilla.org/mozilla-central/raw-file/default/python/mozboot/bin/bootstrap.py -O
python3 bootstrap.py --application-choice=js --vcs=git
```

This is especially useful on fresh machines or VMs where build tools may be missing.

### Rust and cbindgen

The SpiderMonkey configure script requires `cbindgen` (generates C bindings from Rust code). Even
though we remove Rust source files, the build system still checks for `cbindgen`:

```bash
# Install Rust via rustup if not already present
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source "$HOME/.cargo/env"

# Install cbindgen
cargo install cbindgen
```

### autoconf2.13

> **Note (ESR 140.3+):** The `autoconf2.13` step is no longer required. The repository ships with a
> `configure` script you can invoke directly. If you run `autoconf` and get
> `configure.in: No such file or directory`, skip this step -- just use `../configure` directly.

## Configure

Run configure from a build directory inside `js/src`:

```bash
cd js/src
rm -rf _build && mkdir _build && cd _build

../configure \
    --disable-jemalloc \
    --disable-jit \
    --disable-wasm-moz-intgemm \
    --with-system-zlib \
    --without-intl-api \
    --enable-optimize \
    --enable-tests \
    --disable-bootstrap
    # Note: --with-system-icu is omitted. ESR 140.11+ requires system ICU >= 76.1;
    # most Linux distros ship ICU 70-74. With --without-intl-api also set,
    # SpiderMonkey builds cleanly without system ICU. Add --with-system-icu back
    # only if your system has ICU >= 76.1 (check: pkg-config --modversion icu-i18n).
```

### Configure Flags Explained

| Flag                         | Purpose                                                                                                                                                                                                                                            |
| ---------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `--disable-jemalloc`         | MongoDB provides its own memory management (tcmalloc/system allocator). Using jemalloc would conflict with the server's allocator.                                                                                                                 |
| `--disable-jit`              | JIT compilation is disabled for server embedding. MongoDB uses SpiderMonkey only as a JavaScript interpreter, not for performance-critical JS execution.                                                                                           |
| `--with-system-zlib`         | Use the system zlib instead of Mozilla's bundled copy. MongoDB already links against system zlib.                                                                                                                                                  |
| `--without-intl-api`         | Disable Mozilla's Intl API. MongoDB uses ICU directly for internationalization rather than through SpiderMonkey's Intl support.                                                                                                                    |
| `--enable-optimize`          | Build with optimizations enabled. Required for meaningful test results and to match production build behavior.                                                                                                                                     |
| `--enable-tests`             | Build the test suites. Required to run verification after cherry-picks.                                                                                                                                                                            |
| `--disable-wasm-moz-intgemm` | Disable the WASM integer GEMM library. Not needed for MongoDB's use case.                                                                                                                                                                          |
| `--with-system-icu`          | Use the system ICU library. ESR 140.11+ requires ICU ≥ 76.1. Most distros ship ICU 70-74, so **omit this flag** unless your system meets the requirement (`pkg-config --modversion icu-i18n`). With `--without-intl-api` set, omitting it is safe. |
| `--disable-bootstrap`        | Skip Mozilla's automatic dependency bootstrapping during configure. We manage dependencies ourselves.                                                                                                                                              |

## Build

```bash
make -j$(nproc)
```

A full build typically takes 5-15 minutes depending on the machine. If the build fails, check:

- Missing system dependencies (run `bootstrap.py` if needed)
- `cbindgen` not found (ensure `~/.cargo/bin` is in `PATH`)
- Cherry-pick conflicts that were not fully resolved

## Test Suites

Run the test suites in order. The build must succeed before any tests can run.

### 1. `make check` -- Main Test Suite

```bash
make check
```

**Important notes:**

- `check_vanilla_allocations.py` appears to hang but actually takes approximately 10 minutes. Do not
  kill it prematurely.
- This test checks for raw allocation calls (`operator new`, `malloc`, etc.) in compiled objects. It
  will report `TEST-UNEXPECTED-FAIL` for any object file containing vanilla allocations.
- **Verification:** failures are acceptable **only** if they are in code that was NOT modified by
  MongoDB. Check the failing object file names -- if they are in upstream SpiderMonkey code (not
  annotated with `// MONGODB MODIFICATION`), the failures are expected and can be ignored.

To filter test output for unexpected failures:

```bash
make check 2>&1 | sed '/^test-[^(pass|known)]/!d' | sed '/.*check_vanilla_allocations.*/d'
```

If the filtered output is empty, the suite passed.

### 2. `make check-jstests` -- JavaScript Test Suite

```bash
make check-jstests JSTESTS_EXTRA_ARGS=-x=/path/to/.agents/skills/mozjs-upgrade-fork/references/exclude.txt
```

The exclude file is located at `.agents/skills/mozjs-upgrade-fork/references/exclude.txt` in the
mongo repository. It lists tests that are expected to fail because we build without the Intl API
(`--without-intl-api`), which affects Unicode property escapes, case folding, and locale-dependent
behavior.

**The exclude list must be updated for each ESR version.** When upgrading to a new ESR, some tests
may be added or removed. Run `check-jstests` without exclusions first to identify the full set of
failures, then update `exclude.txt` accordingly.

To filter test output for unexpected failures:

```bash
make check-jstests JSTESTS_EXTRA_ARGS=-x=/path/to/exclude.txt 2>&1 \
    | sed '/^test-[^(pass|known)]/!d'
```

If the filtered output is empty, the suite passed.

### 3. `./dist/bin/jsapi-tests` -- JSAPI Test Suite

```bash
./dist/bin/jsapi-tests
```

**Expected failure:** `test_DeflateStringToUTF8Buffer` is known to fail (tracked in SERVER-99489).
This is a pre-existing issue unrelated to the upgrade.

To filter:

```bash
./dist/bin/jsapi-tests 2>&1 \
    | sed '/^test-[^(pass|known)]/!d' \
    | sed '/.*DeflateStringToUTF8Buffer.*/d'
```

If the filtered output is empty, the suite passed.

### 4. `make check-jit-test` -- JIT Test Suite

```bash
make check-jit-test
```

**Many failures are expected** because JIT is disabled (`--disable-jit`). This suite is included for
completeness but is not a gate for the upgrade. Review the output for any failures that seem
unrelated to JIT being disabled.

## Interpreting Results

The general pattern for verifying test results:

1. Run the test suite, capturing output.
2. Filter with the regex `^test-[^(pass|known)]` to find unexpected failures.
3. Remove known exceptions (`check_vanilla_allocations`, `DeflateStringToUTF8Buffer`, exclude list).
4. If the resulting failure list is **empty**, the suite passed.

The `upgrade.sh` script automates this by writing filtered failures to log files:

- `make_check_failures.txt`
- `jstests_failures.txt`
- `jsapi_tests_failures.txt`

Empty log files indicate success.

## Troubleshooting Build Failures

### `cbindgen: command not found`

Ensure Rust and cbindgen are installed and in your PATH:

```bash
cargo install cbindgen
export PATH="$HOME/.cargo/bin:$PATH"
```

### ICU version too old or not found

ESR 140.11+ requires system ICU ≥ 76.1 when `--with-system-icu` is used. Most Linux distros ship ICU
70-74. Check your version:

```bash
pkg-config --modversion icu-i18n
```

**If version < 76.1:** Remove `--with-system-icu` from configure. With `--without-intl-api` also
set, SpiderMonkey builds cleanly without system ICU (this is the recommended approach).

**If ICU is missing entirely:** Install it, then re-run configure without `--with-system-icu`:

```bash
# Ubuntu/Debian
sudo apt-get install libicu-dev

# RHEL/CentOS
sudo yum install libicu-devel
```

### `configure: error: zlib not found`

Install the zlib development package:

```bash
# Ubuntu/Debian
sudo apt-get install zlib1g-dev

# RHEL/CentOS
sudo yum install zlib-devel
```

### Compilation errors in cherry-picked code

If the build fails after cherry-picking, the most likely cause is an unresolved or incorrectly
resolved merge conflict. Check:

1. `git diff` for any remaining conflict markers (`<<<<<<<`, `=======`, `>>>>>>>`).
2. Compare the failing file against the same file on the previous branch to identify what changed.
3. Check the SpiderMonkey migration guide for API changes between ESR versions:
   https://github.com/mozilla-spidermonkey/spidermonkey-embedding-examples/blob/next/docs/Migration%20Guide.md
