---
name: mozjs-upgrade-validate
description: >-
  Use when testing the MozJS/SpiderMonkey upgrade in Evergreen CI/CD or finalizing the upgrade.
  Covers submitting patches across all platforms, the SpiderMonkey-Debug variant, the git hack
  branch workaround for oversized patches, and updating the default branch. Use when asked to
  "test mozjs in evergreen", "submit mozjs patch", "mozjs patch too large", or "finalize mozjs upgrade".
source: 10gen/mongo
license: Internal
mongodb:
  team: query-integration
  owner: mariano.shaar@mongodb.com
  internal: true
---

## Repo-Specific Context

Before starting, check for a repo-specific context file in `references/`. Use
[references/_template.md](references/_template.md) as a template if none exists for your repo.

### Required Platforms

Trigger an Evergreen patch targeting ALL supported platform combinations:

| Platform       | Project                                                  |
| -------------- | -------------------------------------------------------- |
| linux/x86_64   | `mongodb-mongo-master`                                   |
| linux/arm64    | `mongodb-mongo-master`                                   |
| windows/x86_64 | `mongodb-mongo-master`                                   |
| macOS/x86_64   | `mongodb-mongo-master` or `mongodb-mongo-master-nightly` |
| macOS/arm64    | `mongodb-mongo-master` or `mongodb-mongo-master-nightly` |
| linux/ppc64le  | `mongodb-mongo-master-nightly`                           |
| linux/s390x    | `mongodb-mongo-master-nightly`                           |

Note: PPC, s390x, and some macOS variants are only on `mongodb-mongo-master-nightly`.

### Submitting Patches

```bash
# Standard platforms on mongodb-mongo-master
evergreen patch -p mongodb-mongo-master -a required -f -y -u

# Nightly platforms (PPC, s390x)
evergreen patch -p mongodb-mongo-master-nightly \
  -v enterprise-rhel-81-ppc64le -v enterprise-rhel-81-ppc64le-dynamic \
  -v enterprise-rhel-83-s390x -v enterprise-rhel-83-s390x-dynamic \
  -v enterprise-rhel-9-ppc64le -v enterprise-rhel-9-ppc64le-dynamic \
  -v enterprise-rhel-9-s390x -v enterprise-rhel-9-s390x-dynamic \
  -t all -f -y -u
```

If Evergreen rejects the patch for being too large, add `--large` to the command.

### SpiderMonkey-Debug Variant

The `Shared Library Enterprise RHEL 8 SpiderMonkey Debug` variant enables `MOZ_ASSERT` calls. This
MUST be fully passing before merging. Submit it explicitly:

```bash
evergreen patch -p mongodb-mongo-master -v "enterprise-rhel8-arm64-dynamic-spider-monkey-dbg" -t all -f -y -u --large
```

### Large Patch Workaround

If `--large` is needed but still fails with a server timeout:

1. First try `evergreen patch --large`
2. If that fails, use the **git hack branch** approach. See
   `.agents/skills/mozjs-upgrade-validate/references/git-hack-branch.md` for the full procedure, or
   run:
   ```bash
   # Run from the mongo repo root:
   bash .agents/skills/mozjs-upgrade-validate/references/gen-hack-branch.sh <your_mozjs_branch_name>
   ```

### Diagnosing Failures

Use the `evergreen-cicd` skill for failure diagnosis. Key priorities for MozJS upgrades:

1. Compilation failures — often platform-specific (especially Windows MSVC)
2. SpiderMonkey-Debug variant failures — `MOZ_ASSERT` assertions
3. JS test failures — check if they're in the exclude list

### Final Steps (Step 8)

Once the PR lands:

1. Ping **Santiago Roche** and **Chris Wolff** to update the default branch of
   `mongodb-forks/spidermonkey` to the new `spidermonkey-esr{MAJOR}.{MINOR}-cpp-only` branch
2. Verify the default branch was updated on GitHub
