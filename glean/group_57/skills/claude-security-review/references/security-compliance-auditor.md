<!-- hub-reference-banner -->
> **Reference file — part of the `security-review` hub.** Formerly the standalone `security-compliance-auditor` skill.
> Sibling topics in this family are now reference files under the hubs (`security-review`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: security-compliance-auditor
title: "Security Compliance Auditor"
description: >-
  Audits arbitrary codebases for security and compliance gaps using mongodb-ai-policy-audit-kit
  scanning patterns plus MongoDB security and compliance policy mapping. Covers secrets and
  credential exposure, auth/authz risks, insecure storage or transport handling, dependency
  and supply-chain concerns, logging/privacy/PII/Atlas-data handling, insecure defaults and
  configuration drift, AI-tooling policy compliance, and repository governance and
  release-gating controls.
  TRIGGER: auditing any repository for security or compliance gaps; checking for hardcoded
  secrets or credentials; reviewing auth/authz risks; evaluating dependency supply-chain
  posture; assessing AI coding policy compliance; checking repository governance controls
  (CODEOWNERS, branch protection, SSDLC gates).
  SKIP: configuring HTTP security headers on a running server (use http-security-headers);
  designing auth flows for a new application (use web-auth-patterns or api-design-patterns);
  Kubernetes RBAC configuration (use kubernetes-networking).
version: "1.1.0"
updated: "2026-05-29"
category: developer
tags:
  - security
  - compliance
  - audit
  - ssdlc
  - secrets
  - pii
  - supply-chain
  - mongodb-policy
keywords:
  - SSDLC
  - SAST
  - SCA
  - SBOM
  - security audit
  - credential exposure
  - hardcoded secrets
  - AI coding policy
  - CODEOWNERS
  - branch protection
  - THIRD_PARTY_NOTICES
  - PII redaction
  - Atlas data
  - Atlas logs
  - encryption in transit
  - signed commits
  - net-new Critical
  - net-new High
  - security artifact bundle
  - DAST
  - secret scanning
  - CodeQL
whenToUse:
  - "Auditing a repository for security or compliance gaps"
  - "Checking for hardcoded secrets, credentials, or API keys in code"
  - "Reviewing auth/authz risks across middleware and transport configuration"
  - "Evaluating dependency supply-chain posture (license, SCA, SBOM)"
  - "Assessing AI coding policy compliance (approved tooling, human review)"
  - "Checking repository governance controls (CODEOWNERS, branch protection, SSDLC)"
  - "Reviewing logging code for PII or Atlas data exposure"
  - "Producing a structured audit report with severity, confidence, and remediation"
whenNotToUse:
  - "Configuring HTTP security headers on a running server — use http-security-headers"
  - "Designing auth flows for a new application — use web-auth-patterns or api-design-patterns"
  - "Kubernetes RBAC configuration — use kubernetes-networking"
related_skills:
  - http-security-headers
  - software-engineering-patterns
  - security-reviewer
  - devops-infra
---

# Security Compliance Auditor

Audits an arbitrary codebase for security and compliance gaps using the mongodb-ai-policy-audit-kit scanning patterns plus MongoDB security and compliance policy mapping.

## When to use this skill

Use this skill when you need to audit a repository for:

- secrets and credential exposure
- auth/authz risks
- insecure storage or transport handling
- dependency and supply-chain concerns
- logging, privacy, PII, or Atlas-data handling concerns
- insecure defaults and configuration drift
- AI-tooling policy compliance
- repository governance and release-gating controls

This skill is designed for **arbitrary codebases**, not just the repository that contains the policy context file.

## Skill guidance

1. Treat the policy file at `docs/mdb_security_and_compliance_policy_context.md` (relative to the mdb-case-assistant repo root) as the policy source of truth for control mapping and severity language.
2. Reuse the upstream audit-kit patterns instead of inventing new scanners:
   - literal policy/tool string scanning from `audit_policy_scan.py`
   - governance/config/dependency heuristics from `audit_governance.py`
   - rule-based policy findings from `audit_rules.py`
3. Distinguish finding confidence explicitly:
   - **Confirmed** = direct evidence from files, config, scan output, or reproducible code paths
   - **Heuristic** = suspicious pattern inferred from static analysis but not fully proven
   - **Manual confirmation required** = process, org approval, branch protection, Jira, or release-governance evidence not provable from the repo alone
4. Always separate:
   - **repo evidence**
   - **process evidence**
   - **organizational approval evidence**
5. Do not claim compliance based only on absence of evidence. If a control is not directly inspectable, mark it as manual confirmation required.
6. When the upstream audit-kit is available locally, run the bundled wrapper script in this skill to collect reproducible evidence artifacts before writing conclusions.
7. Keep heuristic extensions clearly labeled as **Skill heuristic extension**, and only add them when supported by either the policy context file or an analogous audit-kit detection family.
8. The final report must include severity, policy section, control label, evidence, impacted files, remediation guidance, and confidence classification for each finding.

## Recommended workflow

### Phase 1 — Collect evidence

Read the target repository and gather:

- dependency manifests (`package.json`, `requirements.txt`, `go.mod`, `Cargo.toml`, etc.)
- CI/workflow configs
- repository governance files (`CODEOWNERS`, PR templates, contributing/security docs)
- secret-bearing config files (`.env*`, app config, deployment manifests)
- logging and telemetry code
- auth/authz middleware and transport configuration
- any AI-tooling or developer-policy docs

If the audit-kit is available locally, run the bundled wrapper script from its own repo:

```bash
/path/to/mongodb-ai-policy-audit-kit/scripts/run_audit_kit_suite.sh \
  --target <TARGET_REPO_PATH> \
  --audit-kit /path/to/mongodb-ai-policy-audit-kit \
  --output <TARGET_REPO_PATH>/.audit-artifacts
```

Add `--server-repo` when auditing a server/backend repository.

### Phase 2 — Interpret the upstream scan outputs

The upstream audit-kit contributes three evidence streams:

1. **Policy scan** (`audit_policy_scan.py`)
   - banned/disallowed AI tools
   - direct LLM vendor endpoint references
   - line-level file, line, check id, severity, and text evidence
2. **Governance scan** (`audit_governance.py`)
   - dependency allow/deny heuristics
   - telemetry config issues
   - hardcoded credentials
   - Python AI-call data-flow concerns
   - vector-search quantization heuristic
3. **Rule scan** (`audit_rules.py`)
   - banned tools
   - evaluated but not broadly licensed tools
   - direct AI vendor APIs
   - AI authorship markers
   - server-repo traceability hints

Use the upstream results as evidence, but do not stop there. Correlate them with the target repository's actual implementation and the policy controls below.

### Phase 3 — Map findings to policy controls

Use these exact policy sections and preserve the terminology:

1. **Secure Software Development Lifecycle (SSDLC) & Code Management**
   - security gates
   - vulnerability remediation
   - repository controls
   - code ownership
   - secrets management
   - environment separation
2. **Artificial Intelligence (AI) Coding Policy**
   - approved tooling
   - non-sensitive inputs
   - authorship and accountability
   - meaningful human input
   - two human review phases
3. **Open Source Software (OSS) & Third-Party Dependencies**
   - license allow/deny posture
   - attribution and tracking
   - third-party code segregation
   - `THIRD_PARTY_NOTICES`
4. **Data Protection, Privacy, and Logging**
   - PII redaction
   - Atlas data / Atlas logs handling
   - encryption in transit and at rest

### Phase 4 — Produce the audit report

For each finding include:

- finding id or title
- severity: **Critical / High / Moderate / Low**
- confidence: **Confirmed / Heuristic / Manual confirmation required**
- policy section heading
- control label
- requirement statement from policy
- policy citation(s)
- evidence observed
- impacted files
- remediation guidance
- residual uncertainty or follow-up evidence needed

Use `templates/audit-report-format.md` in this skill as the output structure.

## Direct checks vs heuristic checks

### Directly checkable from a repository

- presence or absence of `CODEOWNERS`
- secret-like strings in code or config
- direct LLM vendor API calls
- hardcoded MongoDB connection strings
- logging/redaction code paths
- dependency manifests and obvious banned/unapproved frameworks
- telemetry configuration in `.vscode/settings.json`
- presence of `THIRD_PARTY_NOTICES`
- CI workflows for SAST, SCA, SBOM, secret scanning, CodeQL, dependency review
- encryption and transport configuration visible in code or deployment manifests

### Heuristic or manual-only controls

- whether prod data is reused in non-prod
- whether release decisions allowed net-new Critical/High issues
- whether security artifact bundles are complete
- whether reviewers performed meaningful human review
- whether AI prompts ever included sensitive data
- whether signed commits and branch protection are actually enforced at the hosted-repo settings layer
- whether legal/security exceptions were approved outside the repo

## Bundled context

### Upstream audit-kit mapping

Source repo: `10gen/mongodb-ai-policy-audit-kit`
Reference commit: `8dee1acff4d0f585e45420da4571e3a2e8ae3d6f`

The skill reuses these upstream implementation patterns:

- `scripts/audit_policy_scan.py` — literal multi-extension line scanning, JSON output
- `scripts/audit_governance.py` — dependency/config/code heuristics for governance and data-flow signals
- `scripts/audit_rules.py` — rule-engine pattern with JSON export and server-repo traceability mode
- `.github/workflows/audit.yml` — CI-gating pattern and optional fail-on-warning posture
- `examples/security-dashboard.example.html` — structured evidence-oriented reporting model

The skill does **not** generalize MongoDB-internal Slack/team aliases, TODO org placeholders, or repo-specific skip paths from the upstream kit.

### Policy mapping guidance

Policy source of truth: `docs/mdb_security_and_compliance_policy_context.md` in the mdb-case-assistant repo.

Preserve these exact terms in findings where relevant:

- SSDLC
- SAST
- SCA
- SBOM
- security artifact bundle
- DAST/API testing
- secret scanning
- CODEOWNERS
- signed commits
- Types 1-4 / Type 5
- net-new Critical/High
- meaningful human input
- THIRD_PARTY_NOTICES
- customer Atlas data
- Atlas logs
- PII
- sensitive
- encryption in transit and at rest

## Practical usage examples

### Audit an application repository

```bash
~/.claude/skills/security-compliance-auditor/scripts/run_audit_kit_suite.sh \
  --target <TARGET_REPO_PATH> \
  --audit-kit ~/src/mongodb-ai-policy-audit-kit \
  --output <TARGET_REPO_PATH>/.audit-artifacts
```

Then read the repository plus the generated artifacts and write a report using:
`~/.claude/skills/security-compliance-auditor/templates/audit-report-format.md`

### Audit a server repository

```bash
~/.claude/skills/security-compliance-auditor/scripts/run_audit_kit_suite.sh \
  --target <TARGET_REPO_PATH> \
  --audit-kit ~/src/mongodb-ai-policy-audit-kit \
  --output <TARGET_REPO_PATH>/.audit-artifacts \
  --server-repo
```

Use the `--server-repo` results to look for missing AI traceability signals alongside the normal security/compliance controls.

## Limitations and extension points

- The upstream audit-kit is static-analysis-heavy and regex/heuristic based; it is not a replacement for a security review.
- `audit_governance.py` only performs shallow Python AST data-flow checks; non-Python data flow still needs manual review.
- The wrapper in this skill expects a local clone of `mongodb-ai-policy-audit-kit` and does not vendor the upstream scripts.
- Good extensions for later:
  - richer license and SBOM evidence parsing
  - repo-hosted settings ingestion for branch protection and signed commits
  - SARIF/CodeQL/Snyk/Black Duck artifact ingestion
  - stronger auth/authz and encryption checks by language/framework
