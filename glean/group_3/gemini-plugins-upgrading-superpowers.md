# upgrading-superpowers

**Category:** AI, Agents & Prompt Engineering
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/upgrading-superpowers/skills/upgrading-superpowers

## Description
Use when upgrading, syncing, or diffing the vendored obra/superpowers skills in 10gen/agent-skills against a new upstream superpowers release.

---

# Upgrading Vendored Superpowers Skills

This repo vendors a subset of skills from [obra/superpowers](https://github.com/obra/superpowers) (MIT) with deliberate MongoDB-specific modifications. Upstream keeps evolving; a naive copy of a new release silently deletes those modifications, and a naive keep of our copy misses upstream fixes.

**Read [references/vendored-skills.md](references/vendored-skills.md) before touching anything.** It records which skills are vendored, the upstream tag they were last synced to, the standing customizations that must survive every sync, and the skills that were deliberately pruned and must not be re-added.

## Process

### 1. Fetch upstream at a release tag

```bash
git clone https://github.com/obra/superpowers /tmp/superpowers
git -C /tmp/superpowers tag --sort=-v:refname | head -1   # newest release, e.g. v6.1.1
```

Always sync against a tag, never `main` HEAD — the recorded tag is what makes the next upgrade's 3-way merge possible. The local Claude plugin cache (`~/.claude/plugins/cache/claude-plugins-official/superpowers/<version>/`) matches the tag of the same version and is fine for quick diffing, but use the clone for merging: you need `git show <tag>:...` for base versions.

### 2. Three-way merge each vendored SKILL.md

For each skill in the inventory, with `BASE` = the last-synced tag recorded in the inventory and `NEW` = the target tag:

```bash
git -C /tmp/superpowers show "BASE:skills/<name>/SKILL.md" > /tmp/base.md
git -C /tmp/superpowers show "NEW:skills/<name>/SKILL.md"  > /tmp/theirs.md
git merge-file -p <repo-skill>/SKILL.md /tmp/base.md /tmp/theirs.md > /tmp/merged.md
```

- A clean merge keeps our customizations and picks up upstream changes automatically.
- A conflict means upstream changed a line we customized. Apply the upstream improvement, then re-apply the customization per the inventory — never resolve by just taking one side.
- Frontmatter: ours always wins structurally (upstream has only `name` + `description`; we add `source`, `license`, `mongodb:`). If upstream changed its `description`, adopt the new wording inside our frontmatter.

### 3. Adaptation pass on every merged file

Upstream text reintroduces patterns we deliberately remove. After merging, re-scrub:

- Replace `superpowers:<skill>` cross-references with the bare skill name — these skills are installed standalone here, not under the superpowers plugin namespace.
- Redirect references to pruned skills using the redirect map in the inventory.
- Move upstream supporting `.md` files into `references/` and rewrite relative links to match (upstream keeps them flat next to SKILL.md).
- Keep the inventory's excluded upstream files (test fixtures, CREATION-LOG.md, visual-companion assets) out of the repo.
- Preserve the MongoDB-added content listed per skill in the inventory (e.g. `references/server.md` / `references/cloud.md` index entries, stack-neutral examples).

### 4. Sync supporting files

Diff each skill's non-SKILL.md upstream files against our `references/` (and `scripts/` if present). Vendor new prompt templates and load-bearing scripts; delete our copies of files upstream removed — but only when nothing in our SKILL.md still links to them.

### 5. Validate

First, diff each merged SKILL.md body against upstream's and confirm every remaining difference maps to a documented customization in the inventory — this catches both merge corruption and undocumented drift. A difference you can't explain means either a lost upstream change (spurious "ours deleted" from a wrong base tag) or a customization nobody recorded; resolve it before moving on.

```bash
skill-validator check .agents/skills/<category>/<name>/skills/<name> --strict --allow-extra-frontmatter
grep -rn "superpowers:" .agents/skills/<category>/<name>/   # must be empty
```

Also grep merged files for the pruned skill names, confirm every relative link resolves, and check code fences stay balanced (`grep -c '^ *\`\`\`'` per file must be even). Registration files (marketplace.json, category plugin.json, CATALOG.md) only change if a skill is added or removed — an in-place upgrade touches none of them.

### 6. Update the inventory

Bump the last-synced tag in [references/vendored-skills.md](references/vendored-skills.md) and record any new customization decisions made during the merge. This step is what keeps the next upgrade cheap — skipping it forces the next person to reverse-engineer the base version from diffs.

## Common Mistakes

- **Copying upstream wholesale because "it's mostly the same"** — the customizations are small and scattered; wholesale copy deletes them invisibly. Always 3-way merge.
- **Syncing from `main` or the plugin cache without recording a tag** — leaves the next upgrade with no merge base.
- **Resolving a merge conflict by picking a side** — conflicts are exactly the places where both an upstream fix and a customization apply; you almost always want both.
- **Re-adding pruned skills** because new upstream text references them — follow the redirect map instead.
- **Skipping the inventory update** — the upgrade works but the provenance is lost again.