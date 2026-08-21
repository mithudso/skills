<!-- hub-reference-banner -->
> **Reference file — part of the `devops-containers-cicd` hub.** Formerly the standalone `cicd-pipelines` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: cicd-pipelines
title: "CI/CD Pipelines Expert"
description: >-
  CI/CD pipeline architecture and operations expert. Covers GitHub Actions (reusable
  workflows, composite actions, matrix strategies, caching, self-hosted runners), CI
  testing strategies (test pyramid, unit/integration/E2E ratios), deployment strategies
  (blue-green, canary, rolling, A/B testing), progressive delivery with feature flags
  (LaunchDarkly, Unleash, Flagsmith, flag lifecycle management), GitOps (ArgoCD, Flux,
  Flagger, drift detection), artifact management (container image building, SBOM, SLSA
  provenance, signing), environment promotion (dev/staging/prod, ephemeral PR
  environments), secrets management (OIDC, Vault, short-lived credentials), monorepo CI
  patterns (Turborepo, Nx, Bazel, affected builds, remote caching), and supply chain
  security (dependency locking, egress firewalls, attestation).
  TRIGGER: designing, building, optimizing, or troubleshooting CI/CD pipelines; choosing
  deployment strategies; implementing feature flags or progressive delivery; setting up
  GitOps with ArgoCD or Flux; configuring GitHub Actions caching or matrix strategies;
  hardening pipeline security; monorepo CI patterns; managing artifacts or supply chain
  integrity.
  SKIP: Kubernetes cluster configuration beyond deployment (use kubernetes-networking);
  Terraform infrastructure provisioning (use terraform-kafka-infra); application-level
  testing patterns (use testing-and-vitest-expert).
version: "1.2.0"
updated: "2026-05-29"
category: developer
tags:
  - cicd
  - devops
  - github-actions
  - deployment
  - gitops
  - feature-flags
  - testing
  - pipelines
  - containers
  - security
keywords:
  - CI/CD
  - continuous integration
  - continuous deployment
  - continuous delivery
  - GitHub Actions
  - reusable workflows
  - composite actions
  - matrix strategy
  - blue-green deployment
  - canary deployment
  - rolling deployment
  - feature flags
  - LaunchDarkly
  - Unleash
  - progressive delivery
  - GitOps
  - ArgoCD
  - Flux
  - Flagger
  - artifact management
  - SBOM
  - SLSA
  - environment promotion
  - secrets management
  - OIDC
  - monorepo
  - Turborepo
  - Nx
  - Bazel
  - supply chain security
  - test pyramid
  - Docker caching
  - self-hosted runners
whenToUse:
  - "Designing, building, or optimizing CI/CD pipelines"
  - "Choosing deployment strategies (blue-green, canary, rolling)"
  - "Implementing feature flags or progressive delivery"
  - "Setting up GitOps workflows with ArgoCD or Flux"
  - "Configuring GitHub Actions workflows, caching, or matrix strategies"
  - "Hardening pipeline security or managing secrets with OIDC"
  - "Implementing monorepo CI patterns with Turborepo or Nx"
  - "Managing artifacts, container images, or supply chain integrity"
  - "Configuring environment promotion and ephemeral PR environments"
  - "Troubleshooting CI/CD failures, flaky tests, or slow pipelines"
whenNotToUse:
  - "Kubernetes cluster configuration beyond deployment — use kubernetes-networking"
  - "Terraform infrastructure provisioning — use terraform-kafka-infra"
  - "Application-level testing patterns — use testing-and-vitest-expert"
related_skills:
  - kubernetes-networking
  - terraform-kafka-infra
  - git-workflows
  - security-reviewer
  - testing-and-vitest-expert
  - code-packaging
---

# CI/CD Pipelines Expert Context

CI/CD (Continuous Integration / Continuous Delivery / Continuous Deployment) is the backbone of modern software delivery. CI automates the build-and-test cycle on every commit; CD automates the promotion of validated artifacts through environments to production. A well-designed pipeline catches defects early, enforces quality gates, and delivers value to users with minimal manual intervention.

## 1. GitHub Actions — Workflows, Reusable Workflows, and Composite Actions

### Core Workflow Anatomy

A GitHub Actions workflow is a YAML file in `.github/workflows/` triggered by events (push, pull_request, workflow_dispatch, schedule, repository_dispatch). Each workflow contains jobs that run on runners, and each job contains ordered steps.

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: 'npm'
      - run: npm ci
      - run: npm test
```

### Reusable Workflows (`workflow_call`)

Reusable workflows define a full multi-job pipeline once and call it from multiple repositories. As of 2026, GitHub supports up to **10 nested reusable workflows** and **50 total workflow invocations** per run.

```yaml
# .github/workflows/shared-ci.yml (callee)
on:
  workflow_call:
    inputs:
      node-version:
        type: string
        default: '22'
    secrets:
      NPM_TOKEN:
        required: true

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ inputs.node-version }}
      - run: npm ci
        env:
          NPM_TOKEN: ${{ secrets.NPM_TOKEN }}
      - run: npm test
```

**Key rules:** Secrets are NOT automatically forwarded — pass them explicitly via the `secrets:` map. Matrices must be declared inside the reusable workflow, not passed as inputs from the caller.

### Composite Actions

Composite actions encapsulate reusable step logic without requiring a separate job. Use them for smaller repeated tasks (setup, caching, formatting, linting).

```yaml
# .github/actions/setup-project/action.yml
name: 'Setup Project'
description: 'Install deps and restore cache'
inputs:
  node-version:
    default: '22'
runs:
  using: composite
  steps:
    - uses: actions/setup-node@v4
      with:
        node-version: ${{ inputs.node-version }}
        cache: 'npm'
    - run: npm ci
      shell: bash
```

**When to use which:**
- Reusable workflows: multi-job pipelines (build + test matrix, deploy)
- Composite actions: single-purpose step bundles (setup, formatting, packaging)

### Matrix Strategies

Matrices run a job across multiple configurations in parallel.

```yaml
strategy:
  fail-fast: false
  matrix:
    os: [ubuntu-latest, windows-latest, macos-latest]
    node: [20, 22]
    exclude:
      - os: windows-latest
        node: 20
```

Dynamic matrices from upstream jobs enable monorepo-aware parallelization:

```yaml
strategy:
  matrix:
    package: ${{ fromJson(needs.detect.outputs.packages) }}
```

### Caching Strategies

Four-layer caching hierarchy for optimal build times:

1. **Package manager cache**: `actions/setup-node@v4` with `cache: 'pnpm'`
2. **Build cache**: Turborepo remote cache via `TURBO_TOKEN` and `TURBO_TEAM`
3. **Heavy dependency cache**: `actions/cache@v4` for Playwright browsers, native modules
4. **Docker layer cache**: `type=gha` backend for Docker builds (can cut build times by 80%+)

```yaml
- uses: docker/build-push-action@v6
  with:
    push: true
    tags: myapp:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

GitHub Actions cache is limited to **10 GB per repository**. Use `hashFiles()` on lockfiles for cache keys. Warm the main branch cache on a schedule so new PRs get cache hits.

### Self-Hosted Runners and Cost Optimization

Runner pricing: Ubuntu $0.008/min, Windows $0.016/min (2x), macOS $0.08/min (10x). Use **Actions Runner Controller (ARC)** on Kubernetes for auto-scaling self-hosted runners.

Cost reduction patterns:
- `concurrency` groups with `cancel-in-progress: true` to avoid duplicate runs
- Timeout enforcement (`timeout-minutes:`) to prevent runaway jobs
- Path-based triggers to skip irrelevant workflows

### 2026 Security Roadmap

GitHub's 2026 security roadmap introduces:

- **Workflow dependency locking**: A `dependencies:` section locks direct and transitive dependencies with commit SHAs for deterministic execution
- **Policy-driven execution controls**: Ruleset-based actor and event rules
- **Scoped secrets**: Secrets bound to specific repos, branches, environments, workflow identities, or paths
- **Native egress firewall**: Layer 7 firewall for GitHub-hosted runners with domain/IP/HTTP method allowlists

**Sources**: [GitHub Actions 2026 Security Roadmap](https://github.blog/news-insights/product-news/whats-coming-to-our-github-actions-2026-security-roadmap/), [GitHub Actions Reusable Workflows Docs](https://docs.github.com/en/actions/how-tos/reuse-automations/reuse-workflows)

## 2. CI Testing Strategies

### The Testing Pyramid (2026 Ratios)

| Scenario | Unit | Integration | E2E |
|---|---|---|---|
| Classic | 70% | 20% | 10% |
| AI-assisted dev | 60% | 25% | 15% |
| High-stakes / regulated | 50% | 30% | 20% |

### Test Layer Characteristics

- **Unit tests**: Run in milliseconds, pinpoint failures precisely, test functions/modules in isolation. Run on every commit.
- **Integration tests**: Catch communication issues between services/components. Run early after unit tests pass. Use test containers or service stubs.
- **E2E tests**: Validate critical user paths through the full stack. Resource-intensive, slower, more brittle. Run later in the pipeline or on merge to main.

### Pipeline Test Organization

```yaml
jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - run: npm test -- --coverage
    # Fast: ~30 seconds

  integration-tests:
    needs: unit-tests
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: test
    steps:
      - run: npm run test:integration
    # Medium: ~2-5 minutes

  e2e-tests:
    needs: integration-tests
    runs-on: ubuntu-latest
    steps:
      - run: npx playwright test
    # Slow: ~5-15 minutes
```

### Test Optimization in CI

- **Parallelization**: Use matrix strategies to shard tests across runners
- **Affected-only testing**: In monorepos, run tests only for changed packages and their dependents
- **Flaky test quarantine**: Isolate flaky tests into a separate non-blocking job; track and fix systematically
- **Test result caching**: Skip tests when inputs have not changed (Nx, Turborepo, Bazel)
- **Fail-fast pipelines**: Run linting and type-checking before expensive test suites

## 3. Deployment Strategies

### Rolling Deployment

Sequentially replaces instances with zero downtime. Kubernetes implementation:

```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 1          # One extra pod during update
    maxUnavailable: 0    # Always maintain full capacity
```

**Pros:** Minimal infrastructure overhead, zero downtime, simple.
**Cons:** Mixed versions serve traffic simultaneously, slow rollback, difficult to test new version in isolation.

### Blue-Green Deployment

Two identical production environments. Traffic switches entirely after validation.

```
[Load Balancer] ---> [Blue: v1.0 (live)]
                     [Green: v1.1 (staged, smoke-tested)]
                          |
                    switch traffic
                          |
[Load Balancer] ---> [Green: v1.1 (live)]
                     [Blue: v1.0 (standby for rollback)]
```

**Implementation pattern (AWS ECS):** Build image → deploy to green service → run smoke tests against green endpoint → switch ALB listener target group → keep blue as instant rollback target.

**Pros:** Instant rollback, full isolation for testing.
**Cons:** Doubles infrastructure cost, database migrations must be backward-compatible.

### Canary Deployment

Routes a small traffic percentage to the new version, progressively increasing if metrics hold.

```yaml
# Istio VirtualService for canary traffic splitting
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
spec:
  http:
    - route:
        - destination:
            host: myapp
            subset: stable
          weight: 95
        - destination:
            host: myapp
            subset: canary
          weight: 5
```

**Progression:** 5% → 25% → 50% → 100%, with metric gates (error rate, latency p99) at each stage. Automated rollback if thresholds are breached.

**Pros:** Minimal blast radius, real-user validation, data-driven promotion.
**Cons:** Complex routing infrastructure, requires strong observability.

### Auto-Rollback Pattern

```bash
for i in $(seq 1 10); do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$HEALTH_ENDPOINT")
  if [ "$STATUS" != "200" ]; then
    kubectl rollout undo deployment/myapp
    exit 1
  fi
  sleep 5
done
```

## 4. Feature Flags and Progressive Delivery

### The 5 Flag Types

| Type | Lifetime | Default | Purpose |
|---|---|---|---|
| Release | 30 days post-GA | Off | Hide incomplete features |
| Experiment | Test window | Off | A/B testing |
| Ops / Kill Switch | Permanent | On | Emergency circuit breaker |
| Permission | Long-lived | Off | Subscription/role gating |
| Configuration | Avoid | — | Use env vars instead |

### Naming Convention

Pattern: `{type}-{team}-{feature}-{context}`

Examples: `release-onboarding-wizard-v2`, `permission-pro-advanced-reports`, `ops-kill-payment-gateway`

### Progressive Rollout Rings

| Ring | Audience | Duration | Metric Gates |
|---|---|---|---|
| Internal | Team only | 1-3 days | Functional correctness |
| Canary | 1-5% users | 2-3 days | Error rate, latency p50/p99 |
| Beta | 10-25% users | 3-7 days | Conversion, revenue |
| Expansion | 50% users | 3-5 days | All metrics stable |
| GA | 100% | Stabilize 7-14 days | Confirm, then remove flag |

### Flag Evaluation Performance

- Target: < 1ms for cached SDK evaluation, < 10ms including network
- Use local SDK caching (30-60 second TTL)
- Fetch flags at startup, update via streaming or polling
- Edge evaluation (Unleash Proxy, LaunchDarkly Edge) for latency-critical paths

### Tools Comparison

| Feature | LaunchDarkly | Unleash | Flagsmith |
|---|---|---|---|
| Hosting | SaaS | Self-hosted / SaaS | Self-hosted / SaaS |
| Pricing | Enterprise | Open-source core | Open-source core |
| SDK support | 25+ languages | 15+ languages | 15+ languages |
| A/B testing | Built-in metrics | Via Unleash Proxy | Built-in |
| Best for | Enterprise scale | K8s self-hosting | Flexibility |

### Implementation Best Practices

1. **Wrap flag evaluations** in dedicated functions — enables removal by changing one file
2. **Test both states** — every gated code path needs enabled AND disabled tests
3. **Enforce flag count limits** — 20-30 active flags per service maximum
4. **Quarterly audits** — 2-4 hours per quarter reviewing flag necessity and removal barriers
5. **Assign owner + expiration** at creation — flags without owners become permanent bugs
6. **Separate config from flags** — timeouts, batch sizes, rate limits belong in env vars

## 5. GitOps — ArgoCD, Flux, and Progressive Delivery

### Core GitOps Principles

1. **Declarative**: All infrastructure and application config is declared in code
2. **Versioned**: Git is the single source of truth — every change is an auditable commit
3. **Automated**: Reconciliation controllers pull desired state and apply it continuously
4. **Self-healing**: Drift from declared state is detected and corrected automatically

### ArgoCD vs Flux

| Aspect | ArgoCD | Flux |
|---|---|---|
| Architecture | Centralized controller + web UI | Modular controllers |
| UI | Built-in web dashboard | CLI-first (optional Weave GitOps UI) |
| Helm support | Native | Via Helm Controller |
| Kustomize | Native | Via Kustomize Controller |
| Multi-cluster | App-of-Apps pattern | Kustomization per cluster |
| Image automation | External (Argo Image Updater) | Native (Image Automation Controller) |
| Best for | Teams wanting visual management | Teams wanting CLI-native, modular GitOps |

### Repository Structure Patterns

**Mono-repo** (small-medium teams):
```
infrastructure/
  base/
    deployment.yaml
    service.yaml
  overlays/
    dev/
      kustomization.yaml
    staging/
      kustomization.yaml
    production/
      kustomization.yaml
```

**Multi-repo** (large orgs with strict team boundaries):
- App repo: application code + Dockerfile
- Config repo: Kubernetes manifests per environment
- Infrastructure repo: cluster-level resources (Terraform, Crossplane)

### Environment Promotion with Kustomize

```yaml
# overlays/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
patchesStrategicMerge:
  - deployment-patch.yaml  # replicas: 5, resource limits
images:
  - name: myapp
    newTag: v1.2.3
```

Promotion = updating the image tag in the production overlay and committing to Git.

### Secrets in GitOps

| Approach | How it works | Git safety |
|---|---|---|
| Sealed Secrets | Encrypt with public key, decrypt in-cluster | Encrypted values committed |
| External Secrets Operator | References paths in Vault/AWS SM/GCP SM | No secrets in Git |
| SOPS | Encrypts YAML fields with KMS keys (Flux native) | Encrypted fields committed |

## 6. Artifact Management and Container Images

### Build Once, Deploy Everywhere

The CI pipeline produces a versioned, immutable artifact once. The CD pipeline promotes that exact artifact through environments. Never rebuild for each environment.

### Container Image Best Practices

- **Multi-stage builds** to minimize image size and attack surface
- **Layer ordering**: Stack from least to most frequently mutated for cache reuse
- **Deterministic tags**: Use commit SHA or semantic version, never `latest` in production
- **Scanning**: Trivy or Grype scans images for CVEs before promotion
- **Signing**: Cosign signs images with SLSA provenance attestation

```yaml
- uses: docker/build-push-action@v6
  with:
    push: true
    tags: |
      ghcr.io/org/app:${{ github.sha }}
      ghcr.io/org/app:${{ steps.version.outputs.tag }}
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

### SBOM and Supply Chain Integrity

- **SBOM** (Software Bill of Materials): Generate with Syft or Trivy for every release artifact
- **SLSA** (Supply-chain Levels for Software Artifacts): Build provenance attestation, artifact signing, build environment isolation
- **Dependency locking**: Pin all CI dependencies by commit SHA; GitHub's 2026 `dependencies:` section automates this

## 7. Environment Promotion and Secrets Management

### Environment Hierarchy

```
PR Preview (ephemeral) -> Dev (auto-deploy) -> Staging (QA, prod-identical) -> Production (approval gates)
```

### Ephemeral PR Environments

```yaml
helm upgrade --install pr-${{ github.event.pull_request.number }} ./chart \
  --set image.tag=${{ github.sha }} \
  --set resources.memory=256Mi \
  --namespace previews
```

### Secrets Management Patterns

| Pattern | Tool | Mechanism |
|---|---|---|
| OIDC federation | GitHub Actions OIDC | Short-lived tokens exchanged for cloud credentials |
| Centralized vault | HashiCorp Vault | Runtime injection, short-lived leases |
| Platform-native | AWS Secrets Manager, GCP Secret Manager | Environment-scoped, auto-rotation |
| Pipeline-native | GitHub Secrets | Encrypted at rest, masked in logs |

**OIDC is the preferred pattern** — eliminates long-lived secrets entirely:

```yaml
- uses: aws-actions/configure-aws-credentials@v4
  with:
    role-to-assume: arn:aws:iam::123456789:role/github-actions
    aws-region: us-east-1
    # No access keys needed — GitHub generates OIDC tokens
```

### Secrets Anti-patterns

- Storing secrets in code, environment files, or CI artifacts
- Using long-lived credentials when OIDC is available
- Printing secrets to logs (even masked, they can leak in error traces)
- Sharing secrets across environments without scoping
- Not rotating secrets on team member departure

## 8. Monorepo CI Patterns

### Change Detection

Only build and test what changed:

```yaml
- uses: dorny/paths-filter@v3
  id: filter
  with:
    filters: |
      frontend:
        - 'packages/frontend/**'
        - 'packages/shared/**'
      backend:
        - 'packages/backend/**'
        - 'packages/shared/**'
```

### Build Tools Comparison

| Tool | Language support | Remote caching | Affected builds | Best for |
|---|---|---|---|---|
| Turborepo | JavaScript/TypeScript | Vercel / custom | Via `--filter` | JS/TS monorepos |
| Nx | Multi-language | Nx Cloud / custom | `nx affected` | Polyglot monorepos |
| Bazel | Any language | Remote Execution API | Native | Enterprise scale |
| Rush | JavaScript/TypeScript | Via build cache | Incremental | Large JS/TS monorepos |

### Turborepo CI Pattern

```yaml
- run: pnpm turbo build --filter='...[origin/main]'
  env:
    TURBO_TOKEN: ${{ secrets.TURBO_TOKEN }}
    TURBO_TEAM: ${{ secrets.TURBO_TEAM }}
```

Requires `fetch-depth: 0` to access full Git history for comparison.

### Nx Affected Pattern

```bash
npx nx affected --target=build --base=origin/main
npx nx affected --target=test --base=origin/main
```

If a teammate already built packages with the same inputs, remote caching provides a cache hit even on a fresh CI machine.

## 9. Release Management

### Semantic Versioning

`MAJOR.MINOR.PATCH` — major for breaking API changes, minor for backward-compatible features, patch for bug fixes.

### Conventional Commits

Format: `type(scope): description`

Types: `feat`, `fix`, `docs`, `chore`, `perf`, `refactor`, `test`, `ci`, `build`

Enables automated changelog generation and version bumping via tools like `semantic-release`, `release-please`, or `changesets`.

### Release Train Model

Bi-weekly releases with feature freeze during QA week. Hotfixes deploy immediately from separate branches. Feature flags decouple deployment from release — deploy the code, enable the feature when ready.

## 10. Anti-Patterns and Troubleshooting

### Common Anti-Patterns

| Anti-Pattern | Problem | Fix |
|---|---|---|
| `latest` tag in production | Non-deterministic deployments | Use commit SHA or semver tags |
| Manual deployment steps | Error-prone, inconsistent | Automate through pipeline |
| Shared mutable runners | State leakage between jobs | Use ephemeral runners or containers |
| Monolithic pipeline | Slow feedback, all-or-nothing | Split into focused workflows |
| No rollback plan | Extended outages on failure | Define rollback in every deploy job |
| Secret sprawl | 29M secrets exposed on GitHub in 2025 | Centralized vault + OIDC |
| Stale feature flags | Code complexity, testing burden | Enforce expiration + quarterly audits |
| Skipping staging | Production surprises | Environment parity is non-negotiable |

### Troubleshooting Checklist

1. **Build failures**: Check dependency resolution, cache invalidation, runner environment drift
2. **Flaky tests**: Isolate to quarantine job, check for timing dependencies, shared state, network calls
3. **Slow pipelines**: Profile step durations, add caching layers, parallelize, use affected-only builds
4. **Deployment failures**: Check health endpoints, resource limits, image pull policies, RBAC
5. **Secret errors**: Verify secret names match, check environment scoping, OIDC role trust policy
6. **Cache misses**: Verify key patterns, check cache eviction (10 GB limit), ensure lockfile stability
7. **GitOps sync failures**: Check webhook validation, CRD version compatibility, resource quotas

## References

1. [GitHub Actions 2026 Security Roadmap](https://github.blog/news-insights/product-news/whats-coming-to-our-github-actions-2026-security-roadmap/)
2. [GitHub Actions Reusable Workflows Docs](https://docs.github.com/en/actions/how-tos/reuse-automations/reuse-workflows)
3. [Testing Pyramid 2026 Guide](https://testomat.io/blog/testing-pyramid-role-in-modern-software-testing-strategies/)
4. [GitOps Guide 2026 (AskAnTech)](https://www.askantech.com/gitops-infrastructure-management-continuous-deployment-argocd-flux/)
5. [Feature Flags 12 Best Practices](https://designrevision.com/blog/feature-flags-best-practices)
6. [Secrets Management in CI/CD (Infisical)](https://infisical.com/blog/secrets-management-cicd)
7. [Docker Layer Caching in GitHub Actions](https://www.blacksmith.sh/blog/cache-is-king-a-guide-for-docker-layer-caching-in-github-actions)
8. [Deployment Strategies Explained (Koyeb)](https://www.koyeb.com/blog/blue-green-rolling-and-canary-continuous-deployments-explained)

### Cross-References

- **git-workflows** — branching strategies that feed CI triggers
- **kubernetes-networking** — service mesh, deployment manifests
- **testing-and-vitest-expert** — test runner configuration for CI integration
- **security-reviewer** — pipeline security audit patterns
- **terraform-kafka-infra** — infrastructure provisioning alongside pipelines
