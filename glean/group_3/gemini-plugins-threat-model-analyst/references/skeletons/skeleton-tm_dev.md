# Skeleton: 4-developer-brief.md

> **⛔ Copy the template content below VERBATIM (excluding the outer code fence). Replace `[FILL]` placeholders.**
> **⛔ DO NOT include CVSS scores, CWE IDs, OWASP categories, or security taxonomy jargon. Plain English only.**
> **⛔ ONLY include design-confident items: architecture and design decisions the team can address NOW, before code is written.**
> **⛔ EXCLUDE: input sanitization/validation, library version issues, certificate version issues, SQL injection, XSS, buffer overflow — these are implementation details for later. Include only findings that impact the design of the features listed in the the original document provided**
> **⛔ Every item MUST describe the risk in plain English AND its specific impact on THIS system (not a generic security risk statement).**
> **⛔ DO NOT use bold inline headers. Use `### ` markdown headings for each issue.**

---

```markdown
# Developer Security Brief

This brief summarizes design-level security considerations for **[FILL: system name]**. It covers only items we are confident about based on the current architecture — issues that should inform design decisions before implementation begins.

This is not a comprehensive security audit. Detailed findings, including implementation-level guidance and full remediation plans, are in [3-findings.md](3-findings.md).

---

## Summary

| # | Issue | Affected Component | Priority |
|---|-------|--------------------|----------|
[REPEAT: one row per included issue, numbered sequentially]
| [FILL: N] | [FILL: short issue name] | [FILL: component name] | [FILL: High / Medium / Low] |
[END-REPEAT]

---

## Issues

[REPEAT: one section per included issue]

### [FILL: N]. [FILL: Short, plain-English issue name]

**Affected component:** [FILL: component name]

**Risk:** [FILL: 1-3 plain-English sentences describing the risk. No CVSS, CWE, or OWASP references. Explain what an attacker could do or what could go wrong.]

**Impact on this system:** [FILL: 1-3 sentences describing the specific impact on THIS system. Be concrete — name the data, service, or user group affected. E.g., "An attacker could forge requests to the internal metadata service, potentially retrieving cloud credentials used by the [ServiceName] component."]

**Design consideration:** [FILL: 1-3 sentences describing what should be considered at design time. Frame as a design decision, not a code fix. E.g., "Consider explicitly defining which external URLs [ServiceName] is permitted to contact, and whether an allowlist can be enforced at the network or service level before routing is implemented."]

[END-REPEAT]
```

**Fixed rules baked into this skeleton:**
- No CVSS scores, CWE IDs, OWASP tags, STRIDE categories, or tier labels
- Each issue has exactly 3 fields: Risk, Impact on this system, Design consideration
- Priority uses only: `High`, `Medium`, `Low` (no SDL Bugbar Severity levels)
- Issues are numbered sequentially (not FIND-XX)
- Plain English throughout — no security taxonomy
