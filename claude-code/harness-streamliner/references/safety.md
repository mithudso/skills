# Safety protocol — backup, validate, restore

Harness config is high-blast-radius: a malformed `settings.json` can break every
future session. The auto-apply mode is only safe because of this loop. Follow it
exactly for **every file you write**.

## The loop (per file)

```bash
F="$HOME/.claude/settings.json"          # the file being edited
TS="$(date -u +%Y%m%dT%H%M%SZ)"
BAK="$F.hso-bak.$TS"

# 1. BACK UP (never edit without one)
cp "$F" "$BAK"

# 2. EDIT — prefer jq into a temp file, then move into place.
#    jq can't produce structurally-broken JSON the way a text Edit can.
jq '<transform>' "$F" > "$F.hso-tmp" && mv "$F.hso-tmp" "$F"

# 3. VALIDATE — must be valid JSON. If not, RESTORE and abort this fix.
if ! jq -e . "$F" >/dev/null 2>&1; then
  cp "$BAK" "$F"
  echo "RESTORED $F from $BAK — fix aborted (invalid JSON after edit)"
  # report this prominently; do not continue editing this file
fi
```

## Rules

1. **No backup, no edit.** If `cp` to the backup fails, stop — don't proceed.
2. **One file at a time.** Finish the backup→edit→validate cycle on one file
   before starting the next. Don't batch edits across files then validate at the
   end — you lose the clean restore point.
3. **Prefer `jq` over text Edit.** Structural transforms via `jq` (dedupe a
   bucket, drop a key, filter an array) can't leave dangling commas or broken
   braces. Use a surgical `Edit` only when `jq` is awkward, and validate the same
   way.
4. **Validate after every write**, not just at the end. Catch corruption at the
   step that caused it. Validity is necessary but not sufficient: a `jq` edit can
   produce *valid JSON with the wrong content* (e.g. an empty `allow` list). For
   any list edit, also assert the element count changed by exactly the intended
   amount before moving on — restore if it didn't.
5. **Restore on any failure** — invalid JSON, or the re-scan showing a new
   problem you didn't intend. Restoring is always correct when in doubt.
6. **Keep backups.** Leave the `*.hso-bak.<ts>` files; list their paths in the
   report so the user can roll back manually if they dislike a change.
7. **Never edit flag-only findings** (secrets, mode conflicts, posture changes).
   No backup makes those safe — they're judgment calls, not corruption risks.

## Confirm-before-delete cases

Some removals are mechanically safe (valid JSON after) but semantically risky —
the JSON is fine, but you may have deleted something the user wanted. Confirm
with the user before removing:

- An MCP server (even a duplicate — confirm *which* duplicate).
- A hook whose script is missing (it may be a temporarily-moved script).
- An `env` var you can't positively classify as dead.

Permission dedupe/subsumption and `autoMode`↔`permissions` overlap removal are
safe to apply without per-item confirmation, because they provably don't change
effective permissions — the broader/identical rule still grants the same access.

## After applying

Re-run the scanner to confirm convergence:

```bash
~/.claude/skills/harness-streamliner/scripts/scan.sh [PROJECT_DIR|none]
```

A clean second pass (empty "Applied") means done. If Medium+ items remain,
either an edit didn't take (check for a RESTORE line) or removing one rule
revealed another — investigate, don't loop blindly.
