# Severity Calibration by Document Type

Referenced from SKILL.md Step 3. Apply these overrides to findings from
`document-critique` before recording them in the findings table. Write the
calibrated severity in the table's "Calibrated?" column. All tiers stay within
the canonical Blocking/Major/Medium/Minor/Nit ladder of
`~/.claude/skill-consolidation/convergence-and-severity.md` (see its
"Appendix: anchored severity examples" for the family-wide anchors).

Default severity comes from `document-critique`. Override only when a column
entry differs from "Default."

| Finding type | Default | Runbook | Incident post-mortem | Weekly update | RFC / spec | Customer-facing deliverable (case analysis / account review / JIMP / RCA letter) | KB article |
|---|---|---|---|---|---|---|---|
| Missing prerequisites section | Major | **Blocking** | Major | Minor | Major | Major | Major |
| Passive voice throughout | Medium | **Major** | Medium | Medium | Minor | Medium | Medium |
| Missing rollback procedure | Major | **Blocking** | Major | N/A | N/A | N/A | Major |
| Unverified version number | Medium | **Major** | Medium | Minor | **Blocking** | **Blocking** | **Major** |
| Missing owner / responsible party | Medium | **Major** | **Major** | Medium | Medium | **Major** | Medium |
| Section ordering suboptimal | Medium | Minor | Medium | Medium | Medium | Medium | Medium |
| Vague time references ("recently") | Medium | **Major** | **Major** | **Major** | Medium | **Major** | Medium |
| Load-bearing customer-visible claim still unverified after Pass 10.5 and unhedged¹ | Medium | N/A | N/A | N/A | N/A | **Blocking** | **Major** |
| Internal-only artifact leaked (Internal KB link, HELP-/Jira prefix, unshared roadmap date)² | Medium | Medium | Medium | Medium | Medium | **Blocking** | **Blocking** (Public / customer-shareable articles only; N/A for Internal-flagged articles) |
| Missing AS-OF date | Minor | Minor | Minor | Minor | Minor | **Major** | **Medium** (KB last-reviewed conventions) |

¹ "Unhedged" = no "as of \<date\>" or "(verify)" qualifier. A hedged or
non-load-bearing unverified claim is **Major** instead, so a closable fix path
(verify, hedge, or cut) always exists and convergence stays reachable. Per
Pass 10.5, a claim that cannot be verified because of auth walls is recorded
as BLOCKED/unverifiable per the canonical contract rather than silently fixed —
the Blocking anchor means "do not ship to a customer unverified."

² Keyed on audience for KB articles (internal references are Internal-flagged
articles' purpose; matches mongodb-kb's Public-only sharing rule). This row
upgrades the engine's former Pass 13 voice-only note, which fired only on the
final iteration.

The three anchor rows above (unverified customer-visible claim, internal-only
artifact leak, missing AS-OF date) fire during normal severity scoring on every
iteration — not as a final-iteration voice check.

If the document type doesn't match any column, use the default severity.
Calibration applies before findings are written to the table so the table
already reflects document-type-appropriate severity.
