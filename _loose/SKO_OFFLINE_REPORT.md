# Offline skill-optimizer — fixes applied (141 skills, ~/.claude/skills)

Zero-token mechanical validator (`skill_optimizer_offline.py`). Read-only audit → targeted fixes.

## Before → after
- critical: 6 → 0
- high: 4 → 1 (the 1 remaining is a false positive; see below)
- medium: 16 → 15
- low: 219 → 221 (+2 from new reference-file pointers)

## Real defects fixed (7)
| Skill | Defect | Fix |
|---|---|---|
| misc-catch-all | stray line before `---` → frontmatter unparsable, skill never loaded | removed stray line |
| web-text-mirror | unquoted `description` w/ `: ` → invalid YAML, skill never loaded | folded block scalar |
| eval-driven-development | description 2221 chars (index cap 1536) → truncated every session | trimmed to 1520 |
| smart-contract-security | description 1623 chars | trimmed to 1532 |
| telemetry-pipeline | body ~15.7k tok (10k hard ceiling) | moved 2 sections → `references/`, dropped TOC |
| customer-context-architect | body ~13.4k tok | moved output-format specs → `references/output-and-file-formats.md` |
| bartending-career-track | body ~10.1k tok | dropped redundant TOC (→9.9k) |

No content deleted from the 3 oversized skills — relocated to `references/` (progressive disclosure).

## NOT fixed — false positives / heuristic noise (deliberate)
- **technical-writing-craft "2 broken links" (high):** the paths (`../connection-strings.md`, `../retry-policy.md`) are inside a blockquote **teaching example** for cross-reference vs forward-link. Fixing = corrupting the example.
- **3 "malformed table" (medium/low):** validator mis-parses valid tables (and a YAML list `- PPL`). Tables are correct.
- **3 "duplicate heading" (medium/low):** subheadings under different parent sections — legitimate structure.
- **~221 low + ~15 medium prose (readability FK-grade, empty-quantifier "voice", ai-ism counts):** heuristic-only (the tool itself says "verify"), subjective, no line locations. Bulk-rewriting live global skills to chase these would degrade quality without a real `sko` judgment pass. Left for per-skill `sko` when next touched.

## Savings

**One-time (this audit):** zero-token validator instead of the LLM `sko` loop on 141 skills ≈ **$308 avoided** (~$2.18/skill × 141 at Opus rates).

**Recurring (the fixes, per year):**
- Always-loaded router layer: desc trims save ~200 tok/session, but reviving 2 dead skills re-adds ~276 tok/session → **net ≈ neutral**. Correctness fixes, not token savings.
- Per-invocation body savings (3 oversized skills): ~9.9k tokens relocated out of always-in-body. At ~104 invocations/skill/yr, section skipped ~60% of the time ≈ **~0.6M input tokens/yr ≈ $4–10/yr**.

**Bottom line:** fixes are primarily *correctness* — 2 skills that silently failed to load now work, 2 descriptions no longer truncate, 3 skills under the hard ceiling. Dollar savings from the fixes are minor (<$10/yr). The big number ($308) is the one-time offline-vs-LLM audit savings.
