# Cherry-Pick Troubleshooting

This reference covers the cherry-pick strategy used during MozJS ESR upgrades, the base commits, the
resolver script, common conflict patterns, and when to escalate.

## Why Cherry-Pick Instead of Rebase?

A `git rebase --onto` does not work reliably for this upgrade because the ESR branches diverge
significantly between versions. The commit history is not shared -- Mozilla's
`mozilla-firefox/firefox` and our fork have different commit graphs, and the base code changes
substantially between ESR releases.

A full-range `git cherry-pick <first>..<last>` also fails because the first commit deletes thousands
of files from the Firefox tree (leaving only SpiderMonkey), which creates an unmanageable number of
merge conflicts when applied to a different ESR base.

The recommended approach is to cherry-pick commits **one by one**, using the resolver script for the
first commit (file deletion) and resolving any remaining conflicts manually.

## The 4 Base Commits

These foundational commits transform a full Firefox checkout into a SpiderMonkey-only repository.
They are applied in order to every new ESR branch.

### Commit 1: "Removed all Non-SpiderMonkey Files"

Deletes everything outside the SpiderMonkey directory tree. This is by far the largest commit and
always produces merge conflicts because the new ESR has files that did not exist in the previous
ESR. The resolver script handles this automatically.

### Commit 2: "Removed Rust from SpiderMonkey Repo"

Removes Rust source files and build system references from the SpiderMonkey tree. MongoDB's
embedding does not use the Rust components.

### Commit 3: "Replaced Rust dependencies in SpiderMonkey with C++ implementation"

Adds C++ replacements for functionality that SpiderMonkey originally implemented in Rust (notably
`encoding_rs` for character encoding). This is a MongoDB-specific addition.

### Commit 4: "Add TestLatin1.cpp"

Adds a test file for the C++ Latin-1 encoding implementation that replaced the Rust `encoding_rs`
library.

## The Resolver Script (Commit 1)

The script
`.agents/skills/mozjs-upgrade-fork/references/resolve_delete_non_spidermonkey_files_cherry_pick.sh`
automates conflict resolution for the first cherry-pick.

### What It Does

1. Parses the cherry-pick output for lines matching `Version HEAD of <path> left in tree`.
2. Deletes each of those files (they exist in the new ESR but were deleted in our commit).
3. Stages all deletions with `git add -u`.
4. Continues the cherry-pick with `GIT_EDITOR=true git cherry-pick --continue`.

### How to Use It

```bash
# Run the cherry-pick, capturing output (it will "fail" with conflicts)
git cherry-pick <commit_one_hash> > my_output_file.txt 2>&1 || true

# Run the resolver
bash .agents/skills/mozjs-upgrade-fork/references/resolve_delete_non_spidermonkey_files_cherry_pick.sh my_output_file.txt deleted_files.txt

# Clean up
rm -f my_output_file.txt deleted_files.txt
```

The script must be run from the root of the spidermonkey repository. The two arguments are:

- `input_file` -- the captured output of the failed cherry-pick
- `output_file` -- a scratch file where extracted filenames are written (can be discarded after)

## Common Conflict Patterns

### Files that exist in the new ESR but were deleted in our fork

This is the most common conflict pattern, especially for commit 1. Resolution: **delete them**. Our
fork intentionally removes non-SpiderMonkey files. If the resolver script missed some files, delete
them manually and stage the changes:

```bash
rm <conflicting_file>
git add <conflicting_file>
```

### Files renamed between ESR versions

Mozilla occasionally renames or moves files between ESR releases. If a cherry-pick fails because the
target file has been renamed:

1. Identify the new path of the file in the target ESR.
2. Apply the same modification to the file at its new path.
3. Stage the resolution and continue.

### MongoDB modifications that conflict with upstream changes

When our `// MONGODB MODIFICATION` changes conflict with upstream Mozilla changes:

1. Examine what changed upstream -- is it a refactor, a bug fix, or a new feature?
2. Keep our modification but adapt it to the new API or code structure.
3. Update the `// MONGODB MODIFICATION` comment if the justification changed.
4. If the upstream change makes our modification unnecessary, remove it (with a note in the commit
   message).

### New files added by upstream that conflict with our additions

If Mozilla added a file with the same name as one we added, compare the contents. Usually our file
should take precedence, but check whether upstream added functionality we need.

## Verifying a Cherry-Pick

After each cherry-pick, compare the diff against the original branch:

```bash
# View the diff of the just-applied commit
git diff HEAD~1 HEAD

# Compare against the same commit on the old branch
git log --oneline <old_branch> | head -20
git show <same_commit_on_old_branch> --stat
```

The diffs should be structurally equivalent -- same files modified, same logical changes.
Differences are expected in line numbers and surrounding context due to upstream changes, but the
**intent** of each modification should be preserved.

## The Annotation Convention

All modifications to SpiderMonkey files **must** be annotated:

```cpp
// MONGODB MODIFICATION: Justification for change.
```

This convention serves two purposes:

1. Makes MongoDB-specific changes easily identifiable during future upgrades.
2. Provides context for why the change was made, which helps when resolving conflicts.

When resolving conflicts, ensure the annotation is preserved or updated. If you add a new
modification during conflict resolution, annotate it.

## When to Stop and Seek Help

Escalate to a team member (or stop the automated process) if you encounter:

- **Semantic conflicts** -- the upstream code changed the behavior or API of something we modified,
  and the correct resolution is not obvious.
- **API changes** -- SpiderMonkey APIs used by MongoDB embedding changed signatures, return types,
  or semantics. These require understanding of both the SpiderMonkey change and our usage.
- **Build system changes** -- Mozilla restructured their build system in a way that breaks our
  configure flags or build process.
- **New Rust dependencies** -- upstream added new Rust code that our C++ replacement layer does not
  cover. This may require writing new C++ replacement code.
- **More than 5 files with non-trivial conflicts in a single cherry-pick** -- this suggests a larger
  structural change that needs careful review.

In all these cases, document what you have found (the commit, the conflict, and your analysis) and
hand off to someone with SpiderMonkey domain expertise.
