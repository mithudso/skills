<!-- hub-reference-banner -->
> **Reference file — part of the `aws-cloud` hub.** Created 2026-06-30 to consolidate AWS
> **secretless / workload-identity** patterns that were thin in `aws-core` and `aws-serverless`.
> Sibling topics are reference files under `aws-cloud` — **not** standalone skills. IAM policy/role
> mechanics live in `references/aws-core.md`; detective/preventive security services in
> `references/aws-security-services.md`; CI/CD pipeline design in `devops-containers-cicd`. For
> Okta/IdP internals use `okta-expert` (`security-review`).

---

---
name: aws-secretless-auth
description: >
  AWS secretless workload identity — proving "who a workload is" to AWS without long-lived access
  keys. TRIGGER: the AWS credential-provider chain & where workloads get creds; EC2 instance profiles
  & IMDSv2 (hop limit, session-token requirement, why IMDSv1 is an SSRF risk); ECS task roles &
  container credential endpoint; **EKS IRSA** (IAM Roles for Service Accounts, OIDC provider,
  projected service-account token, sts:AssumeRoleWithWebIdentity); **EKS Pod Identity** (the newer
  agent-based alternative, association API, no per-cluster OIDC trust, IRSA-vs-Pod-Identity choice);
  Lambda execution roles; IAM Roles Anywhere (X.509 for on-prem/hybrid); **GitHub Actions / GitLab
  OIDC federation** to AWS (web-identity, no stored keys, trust-policy `sub`/`aud` conditions);
  cross-account role assumption & external IDs; SAML/OIDC identity federation & IAM Identity Center;
  eliminating static IAM users. SKIP: IAM policy/role/SCP authoring → references/aws-core.md; KMS/
  GuardDuty/Config/WAF → references/aws-security-services.md; secrets *storage* (Secrets Manager/SSM)
  → references/aws-core.md; CI/CD pipeline design → devops-containers-cicd; Okta/IdP product config →
  okta-expert.
version: "1.0.0"
updated: "2026-06-30"
category: developer
tags:
  - aws
  - security
  - iam
  - workload-identity
  - oidc
  - eks
  - secretless
keywords:
  - credential provider chain
  - imdsv2
  - instance profile
  - ecs task role
  - eks irsa
  - iam roles for service accounts
  - eks pod identity
  - assume role with web identity
  - projected service account token
  - lambda execution role
  - iam roles anywhere
  - github actions oidc
  - gitlab oidc
  - web identity federation
  - cross-account role
  - external id
  - iam identity center
  - sts assumerole
  - oidc provider
  - static access keys
---

# AWS Secretless Workload Authentication

> The goal: **no long-lived AWS access keys anywhere** — every workload obtains short-lived STS credentials by proving its identity (instance, container, pod, CI job, on-prem cert) to AWS. This reference covers the *workload-identity* mechanisms; IAM policy/role authoring is in `aws-core.md`. Verified-as-of 2026-06-30.

## The credential-provider chain (where creds come from)

Every AWS SDK/CLI resolves credentials through an ordered chain; the secretless goal is to land on an **environment-provided role**, never on static keys:

1. Explicit code params → 2. Environment variables (`AWS_ACCESS_KEY_ID` … — avoid in prod) → 3. Shared config/`credentials` file & SSO → 4. **Container credentials** (`AWS_CONTAINER_CREDENTIALS_*`, used by ECS/EKS Pod Identity) → 5. **Web-identity token** (`AWS_WEB_IDENTITY_TOKEN_FILE`, used by IRSA/CI OIDC) → 6. **IMDS** (EC2 instance profile).

The presence of #4/#5/#6 is what makes a workload secretless — STS hands out temporary creds that auto-refresh.

## EC2 — instance profiles & IMDSv2

An EC2 instance assumes a role via an **instance profile**; the SDK fetches temp creds from the **Instance Metadata Service (IMDS)** at `169.254.169.254`.

- **Always enforce IMDSv2** (session-oriented, token-required `PUT` then `GET`). **IMDSv1** (simple `GET`) is exploitable via **SSRF** — a vulnerable app tricked into fetching the metadata URL leaks the role's credentials. Enforce: `HttpTokens=required`, and set **`HttpPutResponseHopLimit=1`** so containers on the host can't reach IMDS through an extra network hop. New AMIs/launch templates can default IMDSv2-only.
- Scope the instance role tightly; instances inherit *one* profile, so co-locating differently-privileged workloads on one instance over-grants.

## ECS — task roles

ECS gives **each task** its own role (the **task role**), distinct from the EC2/host **instance role** and the **task execution role** (which pulls images/secrets for the agent). The SDK fetches creds from a container endpoint (`AWS_CONTAINER_CREDENTIALS_RELATIVE_URI`). Per-task roles mean two containers on the same host get different, least-privilege identities — never share one host role across tasks.

## EKS — IRSA vs Pod Identity

Two ways to give a **Kubernetes pod** an AWS IAM role. Know both and when to pick each:

### IRSA (IAM Roles for Service Accounts) — the original
- Each cluster has an **OIDC identity provider** registered in IAM. A Kubernetes **service account** is annotated with an IAM role ARN; the pod gets a **projected, short-lived service-account JWT**; the SDK calls **`sts:AssumeRoleWithWebIdentity`** with it.
- The IAM role's **trust policy** conditions on the OIDC provider + the specific `system:serviceaccount:<ns>:<sa>` (`sub`) and `sts.amazonaws.com` (`aud`).
- **Strengths:** fine-grained, mature, well-understood, works anywhere. **Friction:** you must create/manage an **OIDC provider per cluster**, and the role trust policy is cluster-specific — so a role isn't reusable across many clusters without editing trust policies, and at fleet scale this is heavy.

### EKS Pod Identity — the newer model (recommended default for new clusters)
- An **EKS Pod Identity Agent** add-on runs on the cluster; you create a **Pod Identity Association** (cluster + namespace + service account → IAM role) via the EKS API — **no per-cluster OIDC provider** to register.
- The IAM role's trust policy targets the **EKS service principal** (`pods.eks.amazonaws.com`) with `sts:AssumeRole` + `sts:TagSession` — **the same role can be reused across many clusters** because trust isn't tied to a cluster-specific OIDC issuer. Associations are managed as EKS resources, not IAM annotations.
- **Choose Pod Identity** for: managing many clusters, simpler/reusable roles, less IAM sprawl. **Stay on / choose IRSA** for: features Pod Identity doesn't yet cover (verify the current gap list), non-EKS Kubernetes, or existing fleets already standardized on IRSA. Both issue short-lived creds; both are secretless.

## Lambda & on-prem

- **Lambda** functions assume an **execution role** automatically; creds arrive via environment/SDK with no configuration. Scope per-function; never bake keys into env vars.
- **IAM Roles Anywhere** extends roles to workloads **outside AWS** (on-prem servers, other clouds) using **X.509 certificates** from a trusted private CA — the server presents a cert, a trust anchor + profile map it to a role, STS issues temp creds. The hybrid alternative to copying access keys to a data-center host.

## CI/CD & cross-account federation

- **GitHub Actions / GitLab CI → AWS via OIDC.** Register the CI provider's OIDC issuer in IAM and create a role whose trust policy conditions on the issuer (`token.actions.githubusercontent.com`), the **`aud`**, and a tightly-scoped **`sub`** (e.g. a specific repo + branch/environment: `repo:org/repo:ref:refs/heads/main`). The pipeline calls `AssumeRoleWithWebIdentity` with the job's OIDC token — **zero stored AWS keys in CI secrets.** This is the modern standard; retire any `AWS_ACCESS_KEY_ID` repo secrets. Lock the `sub` condition down — a wildcard `sub` lets any repo/branch assume the role.
- **Cross-account assumption** — a principal in account A assumes a role in account B whose trust policy names A. Use an **`ExternalId`** for third-party/confused-deputy scenarios (the vendor must present a shared, non-guessable ID), and condition on source account/org.
- **Human access** → **IAM Identity Center** (SSO) federated to your IdP (Okta/Entra/Google) issues short-lived console/CLI sessions per permission set — eliminating IAM *users* for people. (IdP product config → `okta-expert`.)

## The elimination checklist

1. Delete static IAM **users**/access keys for workloads; replace with roles (instance profile / task role / IRSA or Pod Identity / execution role / Roles Anywhere).
2. Replace CI access-key secrets with **OIDC federation**.
3. Enforce **IMDSv2-only** with hop limit 1 across all instances.
4. Use **IAM Identity Center** for humans, not IAM users.
5. Audit residual long-lived keys with **IAM Access Analyzer unused-access** + credential reports (`aws-security-services.md`).

> **Mental model:** secretless = "the environment proves identity, STS hands out a short-lived token." EC2→IMDSv2, ECS→task role, EKS→IRSA *or* Pod Identity, Lambda→execution role, CI→OIDC, on-prem→Roles Anywhere, humans→Identity Center. Any place you'd otherwise paste an access key is a place one of these belongs instead.
