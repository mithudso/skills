<!-- hub-reference-banner -->
> **Reference file — part of the `devops-containers-cicd` hub.** Created 2026-07-01 as the Git **internals
> deep-dive** — the object model and plumbing beneath everyday Git. For usage-level branching, merge/rebase
> policy, conventional commits, hooks, and release automation, use `references/git-workflows.md` (the
> porcelain/workflow guide). Sibling topics are reference files under the devops hubs — **not** standalone
> skills. Ignore any "use the X skill" pointer that names a bare sibling; load that topic's
> `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: git-internals
title: Git Internals (Object Model & Plumbing)
description: >
  How Git actually stores history — the content-addressable object database, the commit DAG, refs, and the
  plumbing commands beneath the porcelain. TRIGGER: the object model (blob/tree/commit/annotated-tag) and
  SHA-1/SHA-256 content addressing; the commit DAG and how history is a graph of parent pointers; refs,
  symbolic refs, HEAD, packed-refs, and the reflog; refspecs and how fetch/push map remote↔local refs; the
  index/staging area and tree objects; plumbing vs porcelain (cat-file, hash-object, rev-parse, update-ref,
  ls-tree, rev-list); packfiles, delta compression, gc, and repacking; loose vs packed objects; worktrees
  and the .git directory layout; recovering lost commits via reflog/fsck; history rewriting mechanics
  (rebase, filter-repo) and why rewriting changes SHAs.
  SKIP: branching strategy, conventional commits, merge/squash policy, hooks, semantic-release (use
  git-workflows.md); CI/CD pipeline design (use cicd-pipelines.md); GitHub/GitLab platform features and
  branch-protection rulesets (use git-workflows.md).
triggers:
  - git internals
  - git object model
  - git plumbing
  - blob tree commit object
  - content addressable
  - commit DAG
  - git refs and reflog
  - refspec
  - packfile
  - git gc repack
  - cat-file hash-object
  - update-ref
  - git worktree
  - recover lost commit
  - .git directory
version: "1.0"
updated: "2026-07-01"
category: developer
tags:
  - git
  - version-control
  - internals
  - plumbing
  - object-model
  - dag
  - packfile
  - reflog
whenToUse:
  - Explaining how Git stores a commit's snapshot vs a diff (it stores snapshots via trees)
  - Debugging a corrupt or bloated repo (loose objects, packfiles, gc)
  - Recovering a commit that was orphaned by a bad reset/rebase (reflog, fsck)
  - Understanding why rebase/filter-repo rewrites SHAs and detaches downstream clones
  - Reasoning about refspecs when fetch/push behaves unexpectedly
---

# Git Internals — the object database beneath the porcelain

Git is a **content-addressable filesystem** with a version-control UI layered on top. Everything in
`.git/objects` is keyed by the hash of its content. Understanding the four object types and how refs point
into them explains almost every "why did Git do that?" question.

## The four object types

| Object | Stores | Points to |
| --- | --- | --- |
| **blob** | raw file bytes (no name, no mode) | nothing |
| **tree** | a directory listing: `mode  type  sha  name` rows | blobs and subtrees |
| **commit** | a snapshot pointer + metadata (author, committer, message) | exactly one root **tree**, and 0..n **parent** commits |
| **tag** (annotated) | a named, signed pointer with its own message | usually a commit |

Key consequences:

- A commit stores a **full snapshot** (via its root tree), not a diff. Diffs are computed on demand by
  comparing two trees. Delta compression happens later in packfiles, not in the logical model.
- Identical file content anywhere in history is stored **once** — the blob SHA is the same, so trees just
  reference it. This is why renaming a huge unchanged file costs almost nothing.
- Object IDs are the SHA of `"<type> <size>\0<content>"` (zlib-compressed on disk). Default is SHA-1;
  SHA-256 repos exist but are not yet interoperable with SHA-1 repos.

## The commit DAG

History is a **directed acyclic graph** of commits linked by parent pointers:

- A normal commit has **one** parent. A **merge** commit has **two or more**. The very first commit has
  **zero** (a "root" commit).
- Branches and tags are just **refs** — movable (branch) or fixed (tag) labels pointing at a commit SHA.
  There is nothing "inside" a branch; deleting a branch only removes the label.
- `git log` walks parent pointers backwards from a starting commit. Order is topological, not chronological
  — commit timestamps can lie (rebases, clock skew).

## Refs, HEAD, and the reflog

- **Refs** live under `.git/refs/` (or compacted into `.git/packed-refs`): `refs/heads/main`,
  `refs/tags/v1.0`, `refs/remotes/origin/main`. Each is a file containing one 40-hex SHA.
- **HEAD** is a symbolic ref, normally `ref: refs/heads/main`. "Detached HEAD" means HEAD holds a raw SHA
  instead of pointing at a branch.
- The **reflog** (`git reflog`, `.git/logs/`) records every local movement of HEAD and branches. It is your
  undo history: a commit orphaned by a bad `reset --hard` or `rebase` is still reachable via its reflog
  entry until it is garbage-collected. `git reflog` + `git reset`/`git branch <name> <sha>` recovers it.

## Refspecs — how fetch/push map refs

A refspec is `[+]<src>:<dst>`. The default fetch refspec
`+refs/heads/*:refs/remotes/origin/*` copies every remote branch into your `origin/*` tracking refs; the
leading `+` allows non-fast-forward updates to tracking refs. Push refspecs work in reverse
(`refs/heads/main:refs/heads/main`). Most "my push/fetch mapped the wrong branch" surprises are a custom or
misread refspec.

## The index (staging area)

The **index** (`.git/index`) is a binary file listing the paths, modes, and blob SHAs that will form the
*next* commit's tree. `git add` writes blobs and updates the index; `git commit` snapshots the index into a
tree and creates a commit pointing at it. "Staged vs working tree vs HEAD" is just three trees being
compared: `git diff` (working↔index), `git diff --cached` (index↔HEAD).

## Plumbing vs porcelain

**Porcelain** = the human commands (`add`, `commit`, `merge`, `log`). **Plumbing** = the low-level,
script-stable commands that operate directly on objects and refs. The useful plumbing set:

| Command | Does |
| --- | --- |
| `git cat-file -t <sha>` / `-p <sha>` | show an object's **type** / **pretty content** |
| `git hash-object [-w] <file>` | compute (and optionally store) a blob SHA |
| `git ls-tree <tree>` | list a tree's entries |
| `git rev-parse <rev>` | resolve a ref/expression to a full SHA |
| `git rev-list <rev>` | walk the DAG, list reachable commit SHAs |
| `git update-ref <ref> <sha>` | move a ref safely (records reflog) |
| `git symbolic-ref HEAD` | read/set what HEAD points to |
| `git fsck` | verify object connectivity; report dangling/unreachable objects |

`git cat-file -p HEAD^{tree}` then drilling into subtrees is the fastest way to *see* the object model
directly.

## Packfiles, loose objects, and gc

- New objects start **loose** — one zlib-compressed file per object under `.git/objects/ab/cdef…`.
- `git gc` (auto-triggered by many commands) repacks loose objects into a **packfile**
  (`.git/objects/pack/*.pack` + `.idx`), applying **delta compression** between similar objects and
  removing unreachable ones past the grace period. `git repack -ad` forces a full repack.
- A bloated `.git` usually means large blobs still reachable from history (not just working-tree size).
  `git count-objects -vH` and `git rev-list --objects --all | git cat-file --batch-check` locate them.

## Worktrees and the .git layout

- `git worktree add ../feat feat` creates a **second working directory** sharing the *same* object database
  and refs, on a different branch — no reclone, no stash dance. Linked worktrees store their own HEAD/index
  under `.git/worktrees/<name>/`; the main `.git` remains the object store.
- Bare repos (`git init --bare`) have no working tree — just the object DB and refs — which is why servers
  use them.

## History rewriting — why SHAs change

A commit's SHA is derived from its content **including its parent SHAs and tree**. So any edit to a commit —
`rebase`, `commit --amend`, `filter-repo` — produces a **new** commit object with a new SHA, and every
descendant is likewise re-created (their parent pointer changed). Consequences:

- Rewriting **published** history forces collaborators to re-sync (their old commits are now orphaned),
  which is why force-pushing shared branches is dangerous. Prefer `--force-with-lease` over `--force`.
- `git filter-repo` (the modern replacement for the slow, foot-gun `filter-branch`) is the tool for purging
  a secret or a huge file from *all* history — it rewrites every affected commit.
- After a rewrite, the old commits linger in the reflog and as unreachable objects until gc; `git reflog`
  is still your recovery path if a rewrite went wrong.
