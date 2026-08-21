# Packaging Python Apps as Standalone Distributables

`lang-python` hub reference. Shipping a Python *application* to people who don't
have your environment — zipapps, interpreter-bundlers (PyInstaller), and
compilers (Nuitka).

**Scope:** this is **application distribution**. The *development* toolchain
(envs, lockfiles, `uv build`/`uv publish`) and **library/wheel** packaging are in
`uv-python-toolchain`; **supply-chain** signing/SBOM/audit is in
`python-supply-chain-security`. Here the question is "how does a user run my app
without `pip install`?"

---

## 1. The distribution spectrum

Ordered by how much they assume about the target machine:

| Approach | Needs Python on target? | Output | Best for |
| --- | --- | --- | --- |
| `zipapp` (stdlib) | **yes** | single `.pyz` | tiny pure-Python tools |
| **shiv** / **pex** | **yes** (+shebang) | self-contained zipapp w/ deps | internal/server deploys where Python exists |
| **PyInstaller** | no | folder or one-file bundle (interpreter + deps) | desktop apps, quick standalone EXE |
| **cx_Freeze** | no | frozen bundle | cross-platform freezing alternative |
| **Nuitka** | no | compiled C → native binary | performance + hardened distribution |
| **PyOxidizer** | no | single Rust-embedded binary | single-file embed (project less active) |
| **Briefcase** (BeeWare) | no | native app/installer (incl. mobile) | GUI/desktop/mobile app stores |

---

## 2. zipapp / shiv / pex — "still needs Python"

- **`zipapp`** (stdlib): `python -m zipapp myapp -m "pkg.cli:main"` → one `.pyz`
  the user runs with `python app.pyz`. Pure-Python only; no bundled interpreter.
- **shiv** (LinkedIn): builds a zipapp that **includes your dependencies**;
  unpacks them to a cache on first run. Requires the target to have a compatible
  Python and shebang support; great for controlled fleets, not for users who
  may not have Python.
- **pex** (Pants): conceptually the same (executable zip with deps), with more
  build/resolution control. Both are ideal when "Python is guaranteed on the
  box" (CI, servers, internal tooling) and you just want one shippable file.

---

## 3. PyInstaller — bundle the interpreter

The mainstream choice (multi-million monthly downloads). Analyzes imports,
gathers dependencies, and packages them **with a CPython interpreter** into a
runnable bundle.

- `--onedir` (default): a folder (faster start, easier to debug, larger
  footprint). `--onefile`: a single executable that self-extracts to a temp dir on
  each launch, convenient but with slower cold start.
- **Hidden imports**: dynamic/`importlib` imports aren't seen by static analysis;
  add `--hidden-import` or a hook. **Data files**: `--add-data`. Many libraries
  ship PyInstaller hooks via `pyinstaller-hooks-contrib`.
- Build is fast (it copies, doesn't compile). You must **build on each target
  OS**; no cross-compiling.

---

## 4. Nuitka — a real Python compiler

Translates your Python to optimized **C** and builds a native binary (it does
*not* just bundle CPython). Typical 2–4× speedups on compute-heavy code, and the
source isn't trivially recoverable.

- `--standalone` (self-contained folder) / `--onefile` (single binary).
- Slower builds and a C toolchain requirement, but the best long-term choice for
  professional desktop software or performance-sensitive distribution.
- Same per-OS build rule; same data-file/plugin considerations (`--include-data-*`,
  `--enable-plugin=` for e.g. PySide/Tkinter).

PyInstaller vs Nuitka: **PyInstaller** for speed of build and quick standalones;
**Nuitka** when you want performance and a compiled, harder-to-reverse artifact.

---

## 5. GUI/mobile and single-binary niches

- **Briefcase** (BeeWare): packages into native installers/app-bundles across
  macOS/Windows/Linux **and iOS/Android** — the path when you're shipping a GUI
  app to app stores.
- **cx_Freeze**: a long-standing freezer; reasonable PyInstaller alternative.
- **PyOxidizer**: embeds Python in a Rust binary for a true single file and can
  import modules from memory; powerful but development has slowed, so evaluate
  project activity before adopting.

---

## 6. Cross-platform realities & pitfalls

- **No cross-compiling.** A Windows `.exe` must be built on Windows, a macOS app
  on macOS, etc. Use CI matrices (GitHub Actions runners per-OS) to produce all
  artifacts.
- **Antivirus false positives.** One-file PyInstaller binaries are frequently
  flagged (self-extraction looks like malware). Code-sign and, on macOS,
  **notarize**; prefer `--onedir` for enterprise to reduce flags.
- **Hidden imports & data files** are the #1 "works locally, breaks frozen" bug —
  test the *built* artifact, declare dynamic imports and bundled resources
  explicitly, and access data via `importlib.resources`, not `__file__` paths.
- **Size.** Bundled interpreter + scientific stack → large binaries; exclude
  unused modules, and don't bundle if you control the runtime (use shiv/pex).
- **Startup cost.** `--onefile`/onefile-Nuitka extract on launch; for
  latency-sensitive CLIs prefer onedir.
- **Don't confuse with packaging for PyPI.** If your users *are* Python devs,
  shipping a wheel + a `[project.scripts]` entry point (pipx/uv install) is
  simpler than freezing — see `uv-python-toolchain`.

---

## Sources

- [PyInstaller documentation](https://pyinstaller.org/en/stable/) (onefile/onedir, hidden imports, hooks)
- [Nuitka — the Python compiler](https://nuitka.net/) (standalone/onefile, performance, plugins)
- [PyInstaller vs Nuitka vs cx_Freeze comparison (Sparx Engineering)](https://sparxeng.com/blog/software/python-standalone-executable-generators-pyinstaller-nuitka-cx-freeze)
- [Comparing PEX, PyOxidizer, and PyInstaller (oriolrius.cat)](https://oriolrius.cat/2024/10/25/comparing-python-executable-packaging-tools-pex-pyoxidizer-and-pyinstaller/)
- [shiv documentation](https://shiv.readthedocs.io/) / [zipapp — Python stdlib](https://docs.python.org/3/library/zipapp.html)
- [Briefcase (BeeWare) documentation](https://briefcase.readthedocs.io/)
