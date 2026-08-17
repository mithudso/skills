<!-- hub-reference-banner -->
> **Reference file — part of the `devops-linux-admin` hub.** Created 2026-07-01 as the cross-cutting
> **regex & text-processing deep-dive**, extracting and deepening the grep/sed/awk material referenced from
> `references/shell-scripting.md` (which remains the Bash/Zsh scripting guide). For Python's `re` module,
> use the `lang-python` hub. Sibling topics are reference files under the devops hubs — **not** standalone
> skills. Ignore any "use the X skill" pointer that names a bare sibling; load that topic's
> `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: regex-text-processing
title: Regex & Text Processing (grep / sed / awk / PCRE)
description: >
  Regular-expression flavors and the classic Unix text-processing toolchain. TRIGGER: the difference
  between POSIX BRE, POSIX ERE, and PCRE (Perl-compatible) and which tool uses which; grep/egrep/grep -P;
  sed substitution, addressing, and hold space; awk field processing, patterns/actions, and built-in
  variables; anchors, character classes, quantifiers, alternation, grouping, and backreferences; greedy vs
  lazy matching; capture groups and named groups; lookahead/lookbehind (PCRE only); common escaping and
  quoting pitfalls in the shell; catastrophic backtracking and ReDoS; choosing grep vs sed vs awk vs a real
  language.
  SKIP: writing/debugging the surrounding Bash or Zsh script (use shell-scripting.md); Python's re module
  and regex library (use lang-python); JSON/YAML structured querying with jq/yq (use shell-scripting.md);
  full parsers for recursive/nested grammars (regex is the wrong tool — use a real parser).
triggers:
  - regex
  - regular expression
  - BRE ERE PCRE
  - grep sed awk
  - grep -P
  - sed substitution
  - awk field processing
  - capture group
  - backreference
  - lookahead lookbehind
  - greedy vs lazy
  - catastrophic backtracking
  - ReDoS
  - character class
version: "1.0"
updated: "2026-07-01"
category: developer
tags:
  - regex
  - text-processing
  - grep
  - sed
  - awk
  - pcre
  - posix
  - redos
whenToUse:
  - Deciding whether a task needs grep, sed, awk, or a real programming language
  - Debugging why a regex matches in one tool (grep -P) but not another (plain grep)
  - Writing a sed or awk one-liner for log/CSV extraction or in-place edits
  - Diagnosing a regex that hangs on certain inputs (catastrophic backtracking / ReDoS)
---

# Regex & Text Processing — flavors, tools, and traps

Two things trip people up constantly: **which regex flavor a tool speaks**, and **which tool to reach for**.
Get those right and most text-processing tasks become one-liners.

## Regex flavors — and who speaks them

| Flavor | Used by | Notable traits |
| --- | --- | --- |
| **POSIX BRE** (Basic) | `grep`, `sed` (default), `ed` | `+ ? { } ( ) \|` are **literal** unless backslash-escaped: `\{`, `\(`, `\|`. No lazy quantifiers. |
| **POSIX ERE** (Extended) | `grep -E`/`egrep`, `sed -E`, `awk` | `+ ? { } ( ) \|` are **metacharacters** directly. Still no lookaround, no lazy, no `\d`. |
| **PCRE** (Perl-compatible) | `grep -P`, Perl, Python `re`, PCRE2 libs | Lookahead/lookbehind, lazy `*?`/`+?`, `\d \w \s`, named groups `(?<name>...)`, backreferences. |

Rules of thumb:

- POSIX ERE and BRE have **no** `\d`, `\w`, `\s`, no lookaround, no lazy quantifiers. Use POSIX classes
  instead: `[[:digit:]]`, `[[:alpha:]]`, `[[:space:]]`.
- `\d`, `\b`, lookahead/lookbehind, and lazy matching mean you are in **PCRE** — that's `grep -P`, Python,
  or Perl, **not** plain grep/sed/awk.
- Anchors `^ $`, classes `[...]`, quantifiers `* + ? {n,m}`, alternation `|`, and grouping `(...)` exist in
  all flavors, but escaping differs (BRE needs `\(`, `\{`, `\|`).

## The toolchain — grep, sed, awk

- **grep** — *find lines that match*. `grep -E` (ERE), `grep -P` (PCRE), `-o` (print only match), `-i`,
  `-v` (invert), `-r` (recurse), `-c` (count). Reach for it first for pure matching/filtering.
- **sed** — *stream editing*, mostly substitution: `sed -E 's/pat/repl/g'`. Use a non-`/` delimiter
  (`s|...|...|`) for paths. Backreferences in the replacement are `\1..\9`; `&` is the whole match.
  In-place edits: GNU `sed -i`, BSD/macOS `sed -i ''` (a required empty backup suffix — a top portability
  gotcha).
- **awk** — *field-oriented records*. `awk -F, '$3 > 100 {print $1, $3}'`. Splits each line into `$1..$NF`;
  built-ins `NR` (record #), `NF` (field count), `FS`/`OFS`. Best when the data is columnar or you need
  arithmetic, aggregation, or per-field logic.

Decision guide: **match → grep; simple substitute → sed; columns/arithmetic/state → awk; nested or
recursive structure (JSON, HTML, code) → a real parser**, not regex.
