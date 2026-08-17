# MozJS Upgrade Validate — Repo-Specific Context

Copy this file to your repo as `.agents/skills/mozjs-upgrade-validate/references/<your-repo>.md`
and fill in the values below.

## Evergreen Projects

| Platform       | Project                          |
| -------------- | -------------------------------- |
| linux/x86_64   | `mongodb-mongo-master`           |
| linux/arm64    | `mongodb-mongo-master`           |
| windows/x86_64 | `mongodb-mongo-master`           |
| macOS/x86_64   | `mongodb-mongo-master-nightly`   |
| macOS/arm64    | `mongodb-mongo-master-nightly`   |
| linux/ppc64le  | `mongodb-mongo-master-nightly`   |
| linux/s390x    | `mongodb-mongo-master-nightly`   |

## SpiderMonkey Fork Owner

Ping **Santiago Roche** and **Chris Wolff** to update the default branch of
`mongodb-forks/spidermonkey` after the PR lands.

## Large Patch Threshold

MozJS upgrade patches can exceed the Evergreen size limit. Try without `--large` first;
add it if Evergreen rejects the patch. If `--large` still fails, use the git hack branch
approach (see `references/git-hack-branch.md`).
