# PowerPC and Z-Series Platform Guide

## Host Access

File a ticket at https://jira.mongodb.org/browse/DEVPROD-15375 to get machine access. Choose the
highest available RHEL version. This takes time -- file early.

## SSH

Use SSH agent forwarding from your workstation to forward SSH keys:

```bash
ssh -o ForwardAgent=yes user@host
# or equivalently:
ssh -A user@host
```

## extract.sh Modification

You must remove `--disable-tests` from the configure flags in `extract.sh`:

```diff
     --enable-optimize \
     --disable-js-shell \
-    --disable-tests \
-    --disable-new-pass-manager
+
```

## Dependencies

cbindgen is required. Install Rust first, then:

```bash
cargo install cbindgen
```

## Build Command

Use the bazelisk wrapper instead of `bazel` directly (standard `bazel build` does not work on these
architectures):

```bash
python bazel/bazelisk.py build --config=local install-devcore
```

On ppc64le, you may also need to add `--disable_warnings_as_errors`.

## Workflow Summary

1. SSH into the host with agent forwarding
2. Clone the repo and check out your branch
3. Set git credentials (local, not global):
   ```bash
   git config user.name "Mariano Shaar"
   git config user.email "mariano.shaar@mongodb.com"
   ```
4. Run `scripts/import.sh` and `scripts/extract.sh` (with the modification above)
5. Run `scripts/gen-config.sh <arch> linux` (e.g., `ppc64le linux` or `s390x linux`)
6. Move `mozilla-release/` out of the way
7. Build with `python bazel/bazelisk.py build --config=local install-devcore`
8. Commit the platform-specific files
