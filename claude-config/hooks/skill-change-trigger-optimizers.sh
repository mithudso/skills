#!/usr/bin/env bash
# PostToolUse hook: when an Edit/Write touches a SKILL.md file under any
# skills/<name>/ directory, inject a system reminder telling Claude to
# run the skill-optimizer iteratively until no medium-or-higher findings
# remain. skill-optimizer is the single owning loop; embedded prompt
# blocks are audited via a bounded prompt-deep-optimizer pass-bundle
# dispatch per the contract's Composed artifacts rule, never a second
# convergence loop.
#
# Loop caps and exit conditions are NOT defined here. They live in
# ~/.claude/skill-consolidation/convergence-and-severity.md (canonical);
# this hook cites that file instead of inlining numbers, per its rule
# that orchestrating agents and hooks must not re-derive exit conditions
# or caps.
#
# Drift tripwire: when the edited file is prompt-deep-optimizer's or
# skill-optimizer's own SKILL.md, grep the loop-consumer files for
# pass-count/cap mentions and surface them in the injected reminder so
# stale copies get flagged against the new SKILL.md.
#
# Triggered from ~/.claude/settings.json hooks.PostToolUse with matcher "Edit|Write".

set -u

# Read the hook payload from stdin once.
payload="$(cat 2>/dev/null || true)"
if [ -z "$payload" ]; then
  exit 0
fi

# Extract the affected file path.
file_path="$(printf '%s' "$payload" | jq -r '.tool_input.file_path // .tool_response.filePath // empty' 2>/dev/null || true)"
if [ -z "$file_path" ]; then
  exit 0
fi

# Only act on SKILL.md (or *-skill.md) under a skills/<name>/ directory.
case "$file_path" in
  */skills/*/SKILL.md|*/skills/*/skill.md|*/skills/*-skill.md)
    ;;
  *)
    exit 0
    ;;
esac

contract="${HOME}/.claude/skill-consolidation/convergence-and-severity.md"

# Drift tripwire: on pdo/sko SKILL.md edits, surface pass-count/cap mentions
# in the loop-consumer files so mismatches with the edited SKILL.md get flagged.
tripwire=""
case "$file_path" in
  */skills/prompt-deep-optimizer/SKILL.md|*/skills/skill-optimizer/SKILL.md)
    matches=""
    for f in \
      "${HOME}/.claude/agents/prompt-optimizer-loop.md" \
      "${HOME}/.claude/agents/convergence-loop-runner.md" \
      "${HOME}/.claude/skills/prompt-helper-optimizer/references/prompt-optimization-algorithms.md"
    do
      [ -f "$f" ] || continue
      m="$(grep -inE '[0-9]+[- ]([a-z]+[- ])?pass|pass[- ]count|iteration cap|hard cap|--max-iter|cap (of |is |at )?[0-9]' "$f" 2>/dev/null || true)"
      if [ -n "$m" ]; then
        matches="${matches}--- ${f}
${m}
"
      fi
    done
    if [ -n "$matches" ]; then
      tripwire="

DRIFT TRIPWIRE: this edit touched an optimizer's own SKILL.md. The following pass-count/cap mentions in loop-consumer files may now be stale — verify each against the edited SKILL.md and flag any mismatch in your response:
${matches}"
    fi
    ;;
esac

# Emit the additionalContext payload back to Claude.
jq -n --arg f "$file_path" --arg c "$contract" --arg t "$tripwire" '{
  hookSpecificOutput: {
    hookEventName: "PostToolUse",
    additionalContext: (
      "Skill file changed: " + $f + "\n\n" +
      "Per the standing instruction, run the skill-optimizer iteratively against this skill until no medium-or-higher findings remain. Iteration caps and exit conditions are defined canonically in " + $c + " — cite that file; do not re-derive caps or instability rules here.\n\n" +
      "Procedure:\n" +
      "1. Invoke the `skill-optimizer` skill (via the Skill tool) and apply its recommendations.\n" +
      "2. If the skill embeds runnable prompt blocks, audit them via prompt-deep-optimizer'\''s pass-bundle dispatch per the Composed artifacts rule in the cited contract.\n" +
      "3. Re-run the skill-optimizer pass. Stop when it reports no medium-or-higher findings, or on any canonical exit condition (iteration cap, no progress, instability) per " + $c + ".\n" +
      "4. If neither skill is installed locally, fall back to `mcp__tam_mcp__tam_search_skills` for any registered optimizer skill, then to `mcp__tam_mcp__tam_optimize_prompt` against the skill text as a final fallback.\n" +
      "5. Report the per-iteration finding-severity distribution in the response." +
      $t
    )
  }
}' 2>/dev/null || true
