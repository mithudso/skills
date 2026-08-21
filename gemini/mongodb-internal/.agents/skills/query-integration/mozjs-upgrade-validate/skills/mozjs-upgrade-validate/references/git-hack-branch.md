# Git Hack Branch — Large Patch Workaround

This is needed when the MozJS patch is too large even for `evergreen patch --large` (server
timeout).

The idea: create a separate branch that contains ONLY the Evergreen YAML changes. This small branch
is submitted as the Evergreen patch. The YAML changes include a "git hack" that checks out the
actual MozJS branch at runtime.

## Automated Approach

Use the helper script:

```bash
cd mozjsupgrade
./gen-hack-branch.sh <your_mozjs_branch_name>
```

## Manual Steps

If the script doesn't work or needs adjustment:

### 1. Create the git hack branch script

Create `evergreen/functions/git_hack_branch.sh`:

```bash
cd src
set -o errexit
set -o verbose
git checkout "$branch_name"
```

### 2. Add the hack branch definition to Evergreen YAML

In `etc/evergreen_yml_components/definitions.yml`, add:

```yaml
"git hack branch": &git_hack_branch
  command: subprocess.exec
  display_name: "get the branch we actually want to test"
  params:
    binary: bash
    args:
      - "src/evergreen/functions/git_hack_branch.sh"
    env:
      branch_name: <your_branch_name>
```

### 3. Insert the hack after every git_get_project call

Insert `*git_hack_branch` after every `*git_get_project` call in the following functions in
`definitions.yml`:

1. `git get project and add git tag`
2. `do non-compile setup`
3. `do benchmark setup`
4. `generate powercycle tasks`
5. `run powercycle sentinel`
6. The compile setup function (~line 2527)

The insertion looks like:

```yaml
    - *git_get_project
    - *git_hack_branch    # <-- ADD THIS LINE
```

### 4. Create the Evergreen-only branch

Create a new branch from the SAME base commit as your MozJS branch (not from the MozJS branch
itself). This branch should contain ONLY:

- The `definitions.yml` changes
- The `evergreen/functions/git_hack_branch.sh` file

```bash
# From the base commit (before MozJS changes)
git checkout -b mozjs-evergreen-hack <base_commit>
# Apply only the Evergreen changes
git add etc/evergreen_yml_components/definitions.yml evergreen/functions/git_hack_branch.sh
git commit -m "Temporary: git hack branch for MozJS upgrade testing"
```

### 5. Submit the small patch

```bash
evergreen patch -p mongodb-mongo-master -a required -f -y -u
```

This small patch will pass the size check. At runtime, the hack branch step will checkout your
actual MozJS branch with the full changes.

### 6. Cleanup

After testing completes successfully, delete the hack branch and remove the `git_hack_branch.sh`
file before merging. The `definitions.yml` changes should NOT be part of the final MozJS PR.
