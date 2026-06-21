---
name: misc-catch-all
title: "Misc / Catch-All (hub)"
description: "Catch-all hub for standalone skills that do not belong to any coherent >=8 family — reached by elimination, not topic match. TRIGGER: only when no other hub or skill fits and you have scanned the routing table below for the specific standalone skill. SKIP: anything matching a real domain hub (mongodb, frontend, security-review, data-analytics, deep-optimizer, writing, devops, etc.) -> that hub. NOTE: this is a deliberate junk drawer; prefer the per-skill routing table in the body over the hub description."
category: custom
origin: local
version: "1.0.0"
updated: "2026-06-21"
model: claude-opus-4-8
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

A **deliberate junk drawer**: standalone skills with no coherent >=8 family and no strong fit into a domain hub. Unlike a real hub, the description here cannot enumerate every member's trigger vocabulary, so **routing into this hub happens by elimination**. Once here, match the request to a row below and **Read the listed `references/` file**.

> If a skill in this drawer later gains >=7 siblings (a real family emerges), promote it out into its own hub via `skill-tree-architect`. This drawer is a holding pen, not a home.

## Routing table

| Skill | Use when | Reference file |
| --- | --- | --- |
| `10gen` | 10gen GitHub repo intelligence for case-tracker and troubleshooting Chrome extension work | `references/10gen.md` |
| `case-mcp-server-guide` | Guidance for starting, using, and troubleshooting the MDB Case Assistant local case MCP server and its mdb_case_* too… | `references/case-mcp-server-guide.md` |
| `dash-docsets` | Dash (Kapeli) docset ecosystem — the offline API-documentation format and its tooling | `references/dash-docsets.md` |
| `diagrams-as-code` | (see reference) | `references/diagrams-as-code.md` |
| `find-skills` | Helps users discover and install agent skills when they ask questions like "how do I do X", "find a skill for X", "is… | `references/find-skills.md` |
| `graph-network-3d-visualization` | Drawing and laying out graphs/networks and 3D node-link scenes in the browser — the rendering/layout layer, not analy… | `references/graph-network-3d-visualization.md` |
| `grill-me` | Calibrated grilling session for stress-testing a plan, design, idea, or decision. First assesses the user's topic kno… | `references/grill-me.md` |
| `mindmaps-and-concept-maps` | Structural and programmatic reference for hierarchical knowledge diagrams: Buzan radial single-root mind maps vs. Nov… | `references/mindmaps-and-concept-maps.md` |
| `mongodb-kb` | MongoDB Knowledge Base article index for troubleshooting, customer support, and escalation | `references/mongodb-kb.md` |
| `offer-design-and-value-proposition` | Design compelling offers and value propositions — the strategy layer beneath marketing copy. Covers Osterwalder's Val… | `references/offer-design-and-value-proposition.md` |
| `programmatic-charting-libraries` | (see reference) | `references/programmatic-charting-libraries.md` |
| `repo-file-analyzer` | Recursively analyze a repo file-by-file, generating per-file summaries (purpose, exports, dependencies), and upload e… | `references/repo-file-analyzer.md` |
| `skill-folder-sync` | Compare two .claude/skills folders and one-way sync them, hub-and-spoke aware | `references/skill-folder-sync.md` |
| `skill-lookup` | Search, retrieve, and install Agent Skills from the prompts.chat registry using MCP tools. Use when the user asks to… | `references/skill-lookup.md` |
| `social-marketing-and-cbsm` | Apply social marketing and CBSM to public-good behavior change. Andreasen's | `references/social-marketing-and-cbsm.md` |
| `webcrypto-vault-reviewer` | Web Crypto API security review and implementation patterns — SubtleCrypto, PBKDF2 key derivation | `references/webcrypto-vault-reviewer.md` |
| `wordcloud-generation` | (see reference) | `references/wordcloud-generation.md` |

<!-- cross-hub-map -->
## Cross-hub map — where every misc-catch-all topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `misc-catch-all` | Misc / Catch-All (unfiled standalone skills) | `references/10gen.md`, `references/api-design-patterns.md`, `references/backprop.md`, `references/build.md`, … |
