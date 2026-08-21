---
name: devops-containers-cicd
description: >-
  Containers, orchestration, CI/CD, IaC & delivery sub-hub (devops family). TRIGGER: Docker/Dockerfile, image optimization, container runtime; Kubernetes workloads/networking/ingress; CI/CD pipeline design (GitHub Actions, reusable workflows, caching, release automation); Terraform/OpenTofu & Apache Kafka infra; git branching/merge/rebase, conventional commits, branch protection, monorepo CI, semantic-release; library packaging & distribution (npm/PyPI/crates, ESM/CJS, monorepo, semver, provenance); self-healing systems & autonomic computing (MAPE-K, Kubernetes reconciliation, AIOps remediation, runbook automation, agentic-SRE); chaos engineering & resilience testing (fault injection, GameDays, blast-radius control, chaos-in-CI). SKIP: Linux kernel/admin internals → devops-linux-internals / devops-linux-admin; logging/tracing/metrics → devops-observability.
origin: local
---

# devops-containers-cicd

Containers, orchestration, CI/CD, IaC & delivery sub-hub (devops family).

This hub routes to on-demand reference files under `references/`. See each spoke for depth.

<!-- cross-hub-map -->
## Cross-hub map — where every devops topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `devops-linux-internals` | devops-linux-internals | `references/linux-kernel-architecture.md`, `references/linux-boot-init.md`, `references/linux-memory-numa.md`, `references/linux-storage-filesystems.md`, … |
| `devops-linux-admin` | devops-linux-admin | `references/linux-sysadmin.md`, `references/systemd.md`, `references/linux-package-management.md`, `references/shell-scripting.md`, … |
| `devops-containers-cicd` | devops-containers-cicd | `references/docker-containers.md`, `references/kubernetes-networking.md`, `references/cicd-pipelines.md`, `references/terraform-kafka-infra.md`, … |
| `devops-observability` | devops-observability | `references/nodejs-observability.md`, `references/pino-structured-logging.md`, `references/sentry-monitoring.md`, `references/ebpf-observability.md`, … |
