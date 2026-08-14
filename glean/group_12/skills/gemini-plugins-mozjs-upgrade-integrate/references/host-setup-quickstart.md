# Quick-Start Script for New Hosts

## Generic Host Provisioning

Spawn hosts use the `amazon2023-latest-large` distro (Linux x86_64) or equivalent. See
`spawn-integration-hosts.md` for how to spawn hosts and get SSH commands.

**Always work from `~/mongo` — bazel requires it and Claude sessions should start there.**

### Prerequisites (run once after workstation setup)

The with-Rust path requires Rust, cargo, and cbindgen — none are installed by the workstation setup:

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --default-toolchain 1.93.1 --profile minimal
source "$HOME/.cargo/env"
cargo install cbindgen
```

`llvm-objdump` is also needed by `extract.sh` but is not on `PATH` by default. Do **not** add the
full toolchain v5 bin directory to `PATH` — it switches the compiler to v5 g++ which statically
links libstdc++, and moz.configure explicitly rejects that. Instead, set just the one variable:

```bash
export LLVM_OBJDUMP=/opt/mongodbtoolchain/v5/bin/llvm-objdump
```

Add this to `~/.bashrc` so it persists across sessions.

### Full Setup Sequence

```bash
git clone git@github.com:10gen/mongo.git
cd mongo
git config user.name "<your name>"
git config user.email "<your email>"

timeout 600 bash etc/set_up_workstation.sh || true   # EngFlow auth blocks on headless hosts; venv finishes before that
source ~/.bashrc
source python3-venv/bin/activate

# Install Rust toolchain (not included in workstation setup)
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --default-toolchain 1.93.1 --profile minimal
source "$HOME/.cargo/env"
cargo install cbindgen

# Make llvm-objdump available without pulling in toolchain v5 compiler
echo 'export LLVM_OBJDUMP=/opt/mongodbtoolchain/v5/bin/llvm-objdump' >> ~/.bashrc
export LLVM_OBJDUMP=/opt/mongodbtoolchain/v5/bin/llvm-objdump

git fetch origin <your_branch>
git checkout <your_branch>

cd src/third_party/mozjs
rm -rf mozilla-release
bash scripts/import.sh
bash scripts/extract.sh
bash scripts/gen-config.sh <ARCH> <OS>

# Verify build
cd ../../..
mv src/third_party/mozjs/mozilla-release/ ../
bazel build install-devcore
# On ppc/s390x: python bazel/bazelisk.py build --config=local install-devcore
```

## Tips

- **SSH agent forwarding**: Use `ssh -A user@host` to forward your SSH keys when connecting to
  remote hosts. This avoids needing to copy keys to each machine.
- **Start Claude in `~/mongo`**: When running a Claude session for the integration step,
  `cd ~/mongo` first so bazel commands work without path confusion.
- **Windows hosts**: See the Windows-specific reference (`platform-windows.md`) for MozillaBuild
  shell setup and other Windows-specific requirements.
- **Commit per-platform**: After verifying the build succeeds on a given platform, commit the
  platform-specific files for a clean git history.
