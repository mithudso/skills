<!-- hub-reference-banner -->
> **Reference file — part of the `aws-cloud` hub.** Created 2026-06-30 to consolidate the AWS
> security-*services* layer (detective + preventive controls) that was previously scattered across
> `aws-core`. Sibling topics in this family are reference files under the `aws-cloud` hub — **not**
> standalone skills. IAM (identity) lives in `references/aws-core.md`; secretless workload auth in
> `references/aws-secretless-auth.md`; secrets storage in `references/aws-core.md`. For app/web
> security *review* (code-level vulns) use `security-review`; for bank AI governance use
> `bank-genai-model-risk-governance`.

---

---
name: aws-security-services
description: >
  AWS security-services layer — the detective, preventive, and data-protection controls beyond IAM
  identity. TRIGGER: KMS (CMK types, key policies vs grants, envelope encryption, rotation, multi-
  Region keys, ABAC on keys); AWS Certificate Manager (ACM); GuardDuty threat detection; Security
  Hub (CSPM, security standards, central findings, ASFF); AWS Config (rules, conformance packs,
  remediation, drift); CloudTrail-as-security (org trail, data events, CloudTrail Lake); WAF & Shield
  (web ACLs, managed rules, rate limits, DDoS, Shield Advanced); Macie (S3 PII discovery); Inspector
  (CVE/vuln scanning of EC2/ECR/Lambda); Detective (investigation graphs); IAM Access Analyzer
  (external-access + unused-access); Network Firewall; Secrets-vs-config (cross-ref); the
  detective-vs-preventive control model; multi-account security with Organizations & delegated
  admin. SKIP: IAM policies/roles/SCPs/permission boundaries → references/aws-core.md; secretless
  workload identity (IRSA, Pod Identity, OIDC federation, IMDSv2) → references/aws-secretless-auth.md;
  app/web code-level security review → security-review; bank/FSI AI governance → bank-genai-model-
  risk-governance; PrivateLink/network isolation mechanics → references/aws-privatelink-vpc-endpoints.md.
version: "1.0.0"
updated: "2026-06-30"
category: developer
tags:
  - aws
  - security
  - kms
  - guardduty
  - security-hub
  - aws-config
  - waf
  - encryption
keywords:
  - aws kms
  - customer managed key
  - key policy
  - kms grant
  - envelope encryption
  - multi-region key
  - guardduty
  - security hub
  - aws config
  - conformance pack
  - cloudtrail lake
  - aws waf
  - web acl
  - aws shield advanced
  - macie
  - amazon inspector
  - amazon detective
  - iam access analyzer
  - unused access
  - network firewall
  - asff
  - cspm
  - delegated administrator
---

# AWS Security Services

> Detective, preventive, and data-protection controls **beyond IAM identity**. IAM answers *who can do what* (`aws-core.md`); this reference covers *how data is protected, how threats are detected, and how posture is governed across accounts*. Verified-as-of 2026-06-30; service feature sets move quarterly — re-verify before quoting limits or "newest" features.

## The control taxonomy

Map every service to the control it provides — this is how you reason about coverage gaps:

| Control type | Question | AWS services |
|---|---|---|
| **Preventive** | Stop bad things happening | IAM/SCPs, KMS, WAF, Shield, Network Firewall, security groups |
| **Detective** | Notice bad things | GuardDuty, Config, CloudTrail, Security Hub, Macie, Inspector, Access Analyzer |
| **Responsive** | Investigate & remediate | Detective, Config remediation, Security Hub automations, EventBridge → Lambda/SSM |
| **Data protection** | Keep data confidential/intact | KMS, ACM, Macie, S3 Block Public Access |

A mature posture has all four. A common gap: heavy preventive IAM work but no detective layer (no GuardDuty/Config/Security Hub) — so misconfigurations go unnoticed.

## KMS — the cryptographic root

**Key Management Service** is the backbone of AWS encryption. Almost every "encrypt at rest" feature (S3 SSE-KMS, EBS, RDS, Secrets Manager) calls KMS.

- **Key types:** AWS-managed keys (`aws/service`, free, auto-rotated, no policy control) · **customer-managed keys (CMK)** (you control policy/rotation/grants — use these for anything sensitive or cross-account) · AWS-owned keys (invisible). Symmetric (default) vs asymmetric (RSA/ECC for sign/verify) vs HMAC.
- **Envelope encryption:** KMS encrypts a per-object **data key**, not your bulk data. The service generates a data key (`GenerateDataKey`), encrypts data locally with the plaintext key, stores the KMS-encrypted key alongside, and discards the plaintext. This is why KMS scales to large objects.
- **Key policy vs IAM vs grants:** the **key policy** is the root of trust for a CMK (a key with no policy granting access is unusable even by admins). Grants give temporary, scoped programmatic delegation (used by AWS services on your behalf). The effective access = key policy ∩ IAM ∩ grants ∩ SCPs.
- **Rotation:** enable automatic annual rotation for CMKs (PG: on by default for new keys in console); old key material is retained to decrypt old ciphertext. **Multi-Region keys** replicate key material across Regions for cross-Region DR/global tables.
- **ABAC on keys** via condition keys (`kms:EncryptionContext`) ties decryption to attributes (tenant, environment) — a key isolation pattern for multi-tenant data.

## Detective stack

- **GuardDuty** — continuous threat detection from CloudTrail, VPC Flow Logs, DNS logs, S3 data events, EKS audit logs, RDS login activity, and Lambda network activity, plus Malware Protection. Findings are severity-scored; route to Security Hub/EventBridge. No agents, no log plumbing — enable per-Region (or org-wide via delegated admin).
- **AWS Config** — records resource configuration *over time* and evaluates **rules** (managed + custom). **Conformance packs** bundle rules to a framework (CIS, PCI, NIST). Supports auto-**remediation** (SSM documents) and detects drift. This is your "is anything non-compliant *right now*, and when did it change?" service.
- **Security Hub** — the aggregation plane: ingests findings from GuardDuty, Inspector, Macie, Config, Access Analyzer, and partners into the normalized **ASFF** (AWS Security Finding Format), runs **security standards** (CIS, AWS FSBP, PCI, NIST 800-53), and centralizes across accounts/Regions. Start here for "one pane of glass."
- **Inspector** — automated **vulnerability** scanning of EC2, ECR container images, and Lambda for CVEs and (EC2) unintended network exposure; continuous, not point-in-time.
- **Macie** — ML-driven **sensitive-data discovery** in S3 (PII, credentials, financial data) and bucket-level security/posture (public, unencrypted, shared).
- **Detective** — builds **investigation graphs** from GuardDuty/CloudTrail/VPC logs to trace the blast radius of a finding (which roles, IPs, resources were involved over time). The "now investigate it" tool after GuardDuty fires.
- **IAM Access Analyzer** — external-access findings (resources reachable from outside the zone of trust) **and** unused-access analysis (unused roles/permissions for least-privilege cleanup) + policy generation/validation (see `aws-core.md` for CLI).

## Edge & network protection

- **WAF** — layer-7 **web ACLs** for CloudFront/ALB/API Gateway/AppSync: AWS Managed Rule groups (OWASP-style, bot control, account-takeover), custom rules, **rate-based rules**, and CAPTCHA/challenge. Log to S3/CloudWatch/Firehose.
- **Shield** — **Standard** (free, automatic L3/L4 DDoS for all AWS) and **Advanced** (paid: enhanced detection, 24×7 Shield Response Team, cost-protection for scaling during attacks, WAF integration).
- **Network Firewall** — managed stateful (Suricata-compatible) firewall for VPC traffic: domain allow/deny lists, IPS rules, egress filtering — the VPC-level complement to security groups/NACLs.
- **ACM** — provisions/renews TLS certificates for ELB/CloudFront/API Gateway (free public certs, auto-renewal); **Private CA** for internal PKI.

## Multi-account governance

Security at scale runs through **AWS Organizations**: delegate a **security/audit account** as **delegated administrator** for GuardDuty, Security Hub, Config, Macie, Inspector, and Access Analyzer so findings aggregate org-wide; enforce guardrails with **SCPs** (preventive) and **declarative policies**; centralize logs (org CloudTrail, Config aggregator) in a dedicated log-archive account. This is the Landing Zone / Control Tower pattern.

> **Mental model:** IAM is the lock on the door (`aws-core.md`); this layer is the alarm system, the cameras, the safe (KMS), and the security-guard station (Security Hub). Enable the detective trio (GuardDuty + Config + Security Hub) early — preventive controls without detection means you never learn when they fail.
