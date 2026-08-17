<!-- hub-reference-banner -->
> **Reference file — part of the `career-and-formal-writing` hub.** Formerly the standalone `policy-and-governance-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: policy-and-governance-writing
description: Authoring and reviewing prescriptive policy documents — corporate policies, information security policies, acceptable-use policies, governance frameworks, and standards. Covers RFC 2119/8174 normative language (MUST/SHOULD/MAY semantics with the RFC 8174 lowercase clarification), the policy-vs-standard-vs-procedure-vs-guideline hierarchy (NIST SP 800-12), required policy components (scope, definitions, roles and responsibilities, exception clauses, review schedule, version history), ISO/IEC 27001 Annex A policy patterns, Plain Writing Act 2010 accessibility requirements, change-management and version-control discipline for policies. TRIGGER: user asks to write, draft, review, or audit a policy, standard, governance document, or framework; user says "make this a policy", "write an acceptable-use policy", "infosec policy", "AUP", "ISMS policy", "policy template", "policy scope section", "exception process", "policy review cadence", "MUST vs SHOULD", "RFC 2119 keywords", "policy vs procedure", "ISO 27001 policy", "NIST 800-12", "Plain Writing Act"; user needs to convert a procedure or runbook into a prescriptive policy. SKIP: writing an RFC, design doc, ADR, or proposal (use rfc-and-design-docs — those documents PROPOSE, this skill PRESCRIBES); writing a runbook or step-by-step procedure (use runbook-craft); converting policy prose into plain-language consumer copy (use plain-language); pure prose editing or audience tone calibration with no governance structure (use writing-expert); rhetorical-argument structure for persuasion (use rhetorical-frameworks-deep); executive memos or comms (use executive-comms).
---

# Policy and Governance Writing

## Overview

Policies prescribe. They do not propose, persuade, or describe — they bind. A reader of a policy needs three things in the first 60 seconds: who is bound by it, what they must do, and what happens if they don't. Everything else (rationale, history, examples) is supporting material.

Policy writing is distinct from RFCs (which propose changes for community review) and from runbooks (which describe a procedure step-by-step). A policy lives at the top of a documentation hierarchy: policies set mandatory direction, standards specify required configurations, procedures specify steps, and guidelines offer recommendations. Confusing the four is the most common policy authoring failure.

This skill covers the RFC 2119 / RFC 8174 (BCP 14) normative-keyword vocabulary, the NIST SP 800-12 policy hierarchy, the ISO/IEC 27001 policy structure pattern, the Plain Writing Act 2010 accessibility constraints, exception-handling clauses, and the change-management discipline that keeps policies enforceable over time.

## Core Concepts

### 1. RFC 2119 / RFC 8174 normative keywords (BCP 14)

Together, RFC 2119 and RFC 8174 form Best Current Practice 14. They define the interpretation of ten specific keywords when written in ALL CAPS:

| Keyword | Meaning |
|---|---|
| MUST / REQUIRED / SHALL | Absolute requirement |
| MUST NOT / SHALL NOT | Absolute prohibition |
| SHOULD / RECOMMENDED | Strong default — deviation requires documented reason and weighing of implications |
| SHOULD NOT / NOT RECOMMENDED | Strong default against — deviation requires documented reason |
| MAY / OPTIONAL | Truly optional; implementers must interoperate with both choices |

**RFC 8174 clarification:** The normative meaning attaches only when the keyword is in ALL CAPITALS. Lowercase "must", "should", and "may" carry their normal English meaning and have no normative weight. This eliminates ambiguity when authors write "the system should respond quickly" (descriptive prose) vs. "the system SHOULD respond within 200ms" (normative rule with documented exception path).

**Required incantation** for any policy that uses these keywords:

> The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in BCP 14 [RFC2119] [RFC8174] when, and only when, they appear in all capitals, as shown here.

Without that statement, the keywords are merely English words and the policy loses its normative bite.

### 2. The policy / standard / procedure / guideline hierarchy

NIST SP 800-12 Rev. 1 places policies at the top of the governance hierarchy. Each layer answers a different question:

- **Policy** — *what* must be true, and *why*. Mandatory. Strategic. Approved at the highest organizational level. Changes rarely (annually at most).
- **Standard** — *which* specific implementation satisfies the policy (TLS 1.3, AES-256, FIPS 140-3). Mandatory. Tactical. Changes when technology moves.
- **Procedure** — *how* to perform a specific task that satisfies the standard (step-by-step). Mandatory for the role performing the task. Operational. Changes whenever the system changes.
- **Guideline** — *suggested* approach when no mandatory standard applies. Advisory. Helps interpret policy where rigid rules would over-constrain.

NIST SP 800-12 also distinguishes three policy *types*: Program Policy (organization-wide security program charter), Issue-Specific Policy (one topic — e.g., acceptable use, remote work, encryption), and System-Specific Policy (one system — e.g., the production payment platform).

**The hierarchy test:** if you can answer "yes" to "would this need to change when we upgrade the firewall?", it is not a policy. It is a standard or procedure.

### 3. Required components of an enforceable policy

Drawing on ISO/IEC 27001 Annex A documentation requirements and NIST SP 800-12, an enforceable policy needs all of these sections. Missing any one creates an enforcement gap.

1. **Title and identifier** — unique policy ID, version number, effective date, supersedes prior version
2. **Purpose** — one-paragraph statement of why the policy exists (the business or regulatory driver)
3. **Scope** — who, what, where the policy applies; explicit in-scope and out-of-scope lists
4. **Definitions** — every term of art used in the policy, listed alphabetically
5. **Policy statements** — the normative rules, using BCP 14 keywords in ALL CAPS
6. **Roles and responsibilities** — named role titles (not individuals) mapped to obligations under this policy
7. **Exceptions** — how to request and approve a documented deviation
8. **Enforcement** — consequences of violation (referencing HR, contract, or regulatory mechanisms)
9. **Related documents** — pointers to standards, procedures, and laws this policy ties into
10. **Review schedule** — review cadence, owner, next review date, change-history table

### 4. The scope section as the contract boundary

The scope section is the load-bearing wall of the policy. It defines who is bound. A policy with vague scope is unenforceable because the accused can always argue they were not in scope.

A defensible scope section answers three questions explicitly:

- **People** — which employees, contractors, vendors, partners, customers? Include role titles, not individuals.
- **Assets** — which systems, data classifications, networks, locations, devices?
- **Activities** — which operations, transactions, or behaviors?

State each as an inclusion list AND an exclusion list. The exclusion list is what stops scope creep during enforcement. Example:

> **In scope:** all employees and contractors of Acme Corp; all systems classified as Confidential or Restricted; all activities performed on Acme-managed devices or networks.
>
> **Out of scope:** customers of Acme's public products; Public-classified marketing material; activities performed on personal devices outside Acme networks.

### 5. The exception clause

Every policy needs a documented exception path. Without one, every deviation is a violation, and operators will route around the policy quietly instead of asking permission. A well-built exception clause has four elements:

1. **Who can grant** — the approval authority (named role, not person)
2. **What must be documented** — business justification, compensating controls, scope of exception, duration
3. **How long it lasts** — maximum duration before re-review (typically 90 or 180 days)
4. **How it is tracked** — where the exception register lives, who audits it

Exceptions are not bypasses. They are time-bounded, risk-accepted, documented deviations. The exception register is itself an enforcement tool — auditors read it to find systemic risk.

### 6. Review schedule and change management

ISO/IEC 27001 clause 7.5.3 (Control of documented information) requires that documents be reviewed and updated as necessary, and that changes and current revision status be identified. Concretely, every policy needs:

- **Review cadence** — at least annual; sooner if regulation or technology changes
- **Document owner** — the role accountable for triggering review (typically the issuing authority or a designated steward)
- **Approval authority** — who signs the next version (often the same authority that signed the original)
- **Change-history table** — a row per version showing version, date, author, change summary, approver

A policy that has not been reviewed in three years is presumptively stale. Auditors flag it; operators ignore it.

### 7. Plain Writing Act 2010 and policy accessibility

The Plain Writing Act of 2010 requires US federal agencies to use "clear Government communication that the public can understand and use." Even outside the federal context, the same constraints apply to any policy that binds a non-specialist audience (employees, contractors, customers). Plain-writing requirements for policies:

- Short sentences (average 15–20 words; max ~25)
- Active voice ("The user MUST encrypt", not "Data MUST be encrypted by the user")
- Common words over jargon (or define jargon in the Definitions section)
- One idea per paragraph
- Lists for parallel items
- Headings that match what readers search for

Plain language does **not** mean simplistic. A policy can require complex behavior in simple sentences. The constraint is on the prose, not the rule.

### 8. ISO/IEC 27001 Annex A policy patterns

ISO/IEC 27001:2022 requires an information security policy (clause 5.2) approved by top management. Annex A (now 93 controls organized into four themes — Organizational, People, Physical, Technological) lists topic-specific policies that mature ISMS programs author separately:

- Access control policy (A.5.15)
- Asset management policy (A.5.9)
- Cryptography policy (A.8.24)
- Information classification policy (A.5.12)
- Acceptable use policy (A.5.10)
- Supplier security policy (A.5.19)
- Incident management policy (A.5.24)

Each follows the same skeleton (purpose, scope, definitions, statements, roles, exceptions, enforcement, review). The discipline is the consistency: an auditor opening any policy in the suite should see the same section names in the same order.

### 9. Policy vs RFC, ADR, and design doc

A common authoring error is writing a "policy" that is actually a proposal. The four documents have different jobs:

| Document | Voice | Audience | Outcome |
|---|---|---|---|
| RFC / Design doc | Proposes | Peers | "Let's discuss whether to do this" |
| ADR (Architecture Decision Record) | Records | Future engineers | "Here is what we decided and why" |
| Policy | Prescribes | Bound parties | "You MUST do this" |
| Runbook | Instructs | Operator on shift | "Step 1: log in to..." |

If the document contains the words "we propose", "we suggest exploring", "open questions", or "alternatives considered with no decision", it is a proposal, not a policy. Write the RFC first, get the decision, then write the policy that codifies the outcome.

### 10. Version-controlled policies as code

Modern policy programs treat policies as version-controlled artifacts in Git, not Word documents in a SharePoint folder. Benefits:

- **Diff-able changes** — every edit shows what changed and who approved it
- **Pull-request review** — same control gate as code
- **Tagged releases** — version 1.0, 1.1, 2.0 are immutable refs
- **Linkable** — policy IDs resolve to permalinks
- **Searchable** — grep across the policy corpus

The downside: Markdown lacks signature blocks. Most programs publish a rendered PDF for the signed copy of record and keep the source-of-truth Markdown in Git. Tag the commit hash that produced the PDF in the change-history table.

## Templates and Examples

### Template: Issue-Specific Policy skeleton (Markdown)

```markdown
# [Policy Title]

| Policy ID | Version | Effective | Supersedes | Owner | Approver |
|---|---|---|---|---|---|
| POL-NNN | 1.0 | YYYY-MM-DD | n/a | [Role] | [Role] |

## 1. Purpose
[One paragraph: why this policy exists.]

## 2. Scope
**In scope:** [people, assets, activities]
**Out of scope:** [explicit exclusions]

## 3. Definitions
- **Term** — definition.

## 4. Normative Language
The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and
"OPTIONAL" in this document are to be interpreted as described in
BCP 14 [RFC2119] [RFC8174] when, and only when, they appear in all
capitals, as shown here.

## 5. Policy Statements
5.1 [Role] MUST [behavior].
5.2 [Role] MUST NOT [prohibition].
5.3 [Role] SHOULD [strong default]. Deviation requires documented exception per Section 7.

## 6. Roles and Responsibilities
- **[Role A]** — [obligations]
- **[Role B]** — [obligations]

## 7. Exceptions
Exceptions to this policy MUST be approved in writing by [Approval Role].
Each exception MUST include: business justification, compensating controls,
scope, duration (maximum 180 days), and approver signature. Exceptions are
tracked in the [Exception Register location].

## 8. Enforcement
Violations of this policy MAY result in [HR / contractual / regulatory consequence].

## 9. Related Documents
- [Linked standard, procedure, regulation]

## 10. Review and Change History
**Review cadence:** Annual, or sooner upon material change to scope or law.
**Next review:** YYYY-MM-DD

| Version | Date | Author | Change | Approver |
|---|---|---|---|---|
| 1.0 | YYYY-MM-DD | [Name] | Initial issue | [Name] |
```

### Example: A well-formed policy statement

> **5.3** Employees with access to Restricted data MUST encrypt that data at rest using AES-256 or stronger when stored on portable devices. Employees SHOULD use the company-issued full-disk-encryption tool. Use of any other encryption tool requires a documented exception per Section 7.

This statement names the bound party (employees with access to Restricted data), uses ALL-CAPS keywords, references a standard (AES-256) without inlining it, gives a strong default, and points to the exception path.

### Example: A poorly-formed policy statement (and the fix)

Bad:

> Sensitive data should be protected appropriately by all relevant parties.

Why it fails: passive voice, vague subject ("all relevant parties"), lowercase "should" (non-normative per RFC 8174), undefined terms ("sensitive", "appropriately"), no standard, no exception path.

Fixed:

> Employees and contractors handling data classified as Restricted MUST encrypt that data at rest using a control listed in the Approved Encryption Standards register (STD-014). Exceptions require approval per Section 7.

## Anti-Patterns

1. **The descriptive policy** — Reads like a textbook chapter. Lots of background, no rules. Fix: every section that does not contain a MUST/MUST NOT/SHOULD belongs in a separate background or commentary document.

2. **Lowercase normatives** — Using "should" and "must" in lowercase while believing they bind readers. RFC 8174 says they don't. Fix: ALL CAPS or rewrite.

3. **The grab-bag scope** — "This policy applies to everyone and everything." Unenforceable because it is non-discriminating. Fix: explicit in-scope and out-of-scope lists.

4. **Personal-name responsibilities** — "Jane Smith MUST approve all exceptions." Jane leaves; policy breaks. Fix: name roles, not people. "The Chief Information Security Officer MUST approve all exceptions."

5. **No exception path** — Forces operators to violate or route around. Fix: every SHOULD or MUST that has plausible legitimate deviation needs a Section 7 exception path.

6. **Procedure leakage** — Step-by-step instructions inside a policy ("first open the console, then click Settings, then..."). Fix: pull steps into a referenced procedure document; the policy says *what*, the procedure says *how*.

7. **Frozen-in-amber policies** — Approved in 2019, references defunct systems and obsolete regulations. Fix: mandatory annual review with documented owner; missed review triggers escalation.

8. **The standard masquerading as policy** — "All servers MUST run TLS 1.3." When TLS 1.4 ships, the policy is wrong. Fix: policy says "All servers MUST use the cryptographic protocol versions listed in the Approved Cryptography Standard (STD-007)." The standard moves; the policy holds.

9. **Mixing prescriptive and proposal voice** — "We recommend that teams consider adopting..." That is an RFC sentence, not a policy sentence. Fix: decide first, then prescribe. If you cannot decide, you are not ready to write the policy.

10. **No definitions block** — Uses "sensitive data", "appropriate controls", "timely manner" with no definition. Every undefined term is a future enforcement dispute. Fix: alphabetized Definitions section with every term of art.

## Decision Heuristics

**Is this document a policy, standard, procedure, or guideline?**

- Does it answer *what must be true and why*? → Policy
- Does it answer *which specific implementation*? → Standard
- Does it answer *how to perform the task step-by-step*? → Procedure
- Does it answer *what is suggested when no rule applies*? → Guideline

**Do I need RFC 2119 keywords?**

- If the document binds anyone to do or not do something → Yes, include the BCP 14 incantation and use ALL CAPS
- If the document is descriptive or proposes a change → No, use normal English

**Should this be MUST or SHOULD?**

- Is there *any* legitimate scenario where deviation is acceptable? → SHOULD (with exception path)
- Is deviation a violation in all cases? → MUST
- If you write SHOULD without an exception path, change to MUST

**Should this be one policy or several?**

- One topic, one audience, one set of obligations → one policy
- Multiple topics that need different review owners or approval authorities → multiple policies
- Policies citing each other is healthy; policies copying each other is duplication risk

**Should this section live in the policy or in a referenced document?**

- Will it change more than once per year? → Referenced document (standard or procedure)
- Will it never change without re-approving the whole policy? → In the policy
- Is it implementation-specific? → Referenced standard or procedure

**Is this scope enforceable?**

- Can an auditor read the scope and tell whether a specific person, system, or activity is bound? → Yes
- If "depends" or "context-dependent", rewrite

## References

1. **RFC 2119** — Bradner, S. "Key words for use in RFCs to Indicate Requirement Levels." IETF, March 1997. <https://datatracker.ietf.org/doc/html/rfc2119>
2. **RFC 8174** — Leiba, B. "Ambiguity of Uppercase vs Lowercase in RFC 2119 Key Words." IETF, May 2017. <https://www.rfc-editor.org/rfc/rfc8174.html>
3. **NIST SP 800-12 Rev. 1** — "An Introduction to Information Security." NIST, June 2017. <https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-12r1.pdf>
4. **Plain Writing Act of 2010** — Public Law 111-274. Federal plain-language guidelines: <https://www.plainlanguage.gov/guidelines/>
5. **ISO/IEC 27001:2022** — "Information security, cybersecurity and privacy protection — Information security management systems — Requirements." International Organization for Standardization.
6. **Plain Language at Digital.gov** — <https://digital.gov/resources/plain-writing-act>

## Cross-references

- **rfc-and-design-docs** — for documents that *propose* (this skill is for documents that *prescribe*)
- **runbook-craft** — for step-by-step procedures (the *how* below the policy)
- **plain-language** — for translating policy prose into consumer-facing copy
- **writing-expert** — for general prose craft once the policy structure is set
- **security-compliance-auditor** — for auditing an existing policy against ISO 27001 / SOC 2 controls
