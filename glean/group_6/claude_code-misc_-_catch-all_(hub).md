# Misc / Catch-All (hub)

**Category:** Writing, Communication & Marketing
**Platform:** Claude Code
**Original Path:** claude-code/misc-catch-all

## Description
The task asks me to return the fixed file content directly:

---

The task asks me to return the fixed file content directly:

---
name: misc-catch-all
title: "Misc / Catch-All (hub)"
description: "Catch-all hub for standalone skills that do not belong to any coherent >=8 family — reached by elimination, not topic match. TRIGGER: only when no other hub or skill fits and you have scanned the routing table below for the specific standalone skill. SKIP: anything matching a real domain hub (mongodb, frontend, security-review, data-analytics, deep-optimizer, writing, devops, etc.) -> that hub. NOTE: this is a deliberate junk drawer; prefer the per-skill routing table in the body over the hub description."
category: custom
origin: local
version: "1.0.0"
updated: "2026-06-21"
model: Codex-opus-4-8
effort: high
tags:
  - misc
  - catch-all
  - unfiled
  - standalone
whenToUse:
  - "a standalone skill that has no coherent family and no strong hub fit"
  - "you have checked every domain hub and none matches — scan this routing table by skill name"
whenNotToUse:
  - "the topic matches any real domain hub → use that hub, not this drawer"
related_skills: []
metadata:
  changelog:
    - "2026-06-21 created junk-drawer hub v1.0.0: folded 36 unfiled standalone skills via skill-tree-architect at user request (no coherent family; reached by elimination — body routing table carries per-skill vocab since the description cannot enumerate all members under the Glean cap)"
---
# Misc / Catch-All (hub)

**Deliberate junk drawer**: standalone skills, no >=8 family, no domain hub fit. Can't enumerate all triggers → **route by elimination**. Match request to row, **Read listed `references/` file**.

> Skill gains >=7 siblings → promote via `skill-tree-architect`. Holding pen, not home.

## Routing table

| Skill | Use when | Reference file |
| --- | --- | --- |
| `10gen` | 10gen GitHub repo intelligence for case-tracker and troubleshooting Chrome extension work | `references/10gen.md` |
| `case-mcp-server-guide` | Start, use, troubleshoot MDB Case Assistant local MCP server and mdb_case_* too… | `references/case-mcp-server-guide.md` |
| `dash-docsets` | Dash (Kapeli) docset ecosystem — offline API-doc format and tooling | `references/dash-docsets.md` |
| `diagrams-as-code` | (see reference) | `references/diagrams-as-code.md` |
| `find-skills` | Discover/install agent skills when user asks "how do I do X", "find skill for X", "is… | `references/find-skills.md` |
| `graph-network-3d-visualization` | Draw/lay out graphs/networks and 3D node-link browser scenes — rendering/layout layer, not analy… | `references/graph-network-3d-visualization.md` |
| `grill-me` | Stress-test plan/design/idea/decision. Assesses user topic kno… | `references/grill-me.md` |
| `mindmaps-and-concept-maps` | Hierarchical knowledge diagrams: Buzan radial mind maps vs. Nov… | `references/mindmaps-and-concept-maps.md` |
| `mongodb-kb` | MongoDB KB index for troubleshooting, support, escalation | `references/mongodb-kb.md` |
| `offer-design-and-value-proposition` | Design offers/value props — strategy layer beneath marketing copy. Osterwalder Val… | `references/offer-design-and-value-proposition.md` |
| `programmatic-charting-libraries` | (see reference) | `references/programmatic-charting-libraries.md` |
| `repo-file-analyzer` | Analyze repo file-by-file, generate per-file summaries (purpose, exports, deps), upload e… | `references/repo-file-analyzer.md` |
| `skill-folder-sync` | Compare two .Codex/skills folders and one-way sync, hub-and-spoke aware | `references/skill-folder-sync.md` |
| `skill-lookup` | Search/retrieve/install Agent Skills from prompts.chat registry via MCP. User asks to… | `references/skill-lookup.md` |
| `social-marketing-and-cbsm` | Apply social marketing and CBSM to public-good behavior change. Andreasen's | `references/social-marketing-and-cbsm.md` |
| `webcrypto-vault-reviewer` | Web Crypto API security review and implementation patterns — SubtleCrypto, PBKDF2 key derivation | `references/webcrypto-vault-reviewer.md` |
| `wordcloud-generation` | (see reference) | `references/wordcloud-generation.md` |

<!-- cross-hub-map -->
## Cross-hub map — where every misc-catch-all topic lives

Family split across hubs. Task not in routing table → reference under sibling hub — **activate hub or `Read` its `references/<name>.md`**. All former standalone skills now references under hubs (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `misc-catch-all` | Misc / Catch-All (unfiled standalone skills) | `references/10gen.md`, `references/api-design-patterns.md`, `references/backprop.md`, `references/build.md`, … |