<!-- hub-reference-banner -->
> **Reference file — part of the `devops-containers-cicd` hub.** Formerly the standalone `kubernetes-networking` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: kubernetes-networking
title: "Kubernetes & Networking"
description: >-
  Kubernetes workloads, networking, and platform engineering expert. Covers core
  primitives (pods, deployments, StatefulSets, DaemonSets, Jobs, HPA/VPA/KEDA),
  networking (CNI plugins: Cilium/Calico/Flannel, service mesh: Istio/Linkerd/Cilium,
  ingress, Gateway API, NetworkPolicy, CoreDNS), operators, Helm/Kustomize, GitOps
  (ArgoCD/Flux), observability (Prometheus/Grafana/OTel), and security (RBAC, PSS,
  admission controllers, secrets management). Covers production patterns,
  anti-patterns, and troubleshooting for 2025-2026 Kubernetes releases.
  TRIGGER: designing or debugging Kubernetes workloads; choosing or configuring a CNI
  plugin; setting up or troubleshooting service mesh; configuring ingress or Gateway
  API; writing NetworkPolicy; building or evaluating Kubernetes Operators; packaging
  with Helm or Kustomize; setting up GitOps with ArgoCD or Flux; instrumenting
  Prometheus/Grafana/OTel; hardening RBAC, PSS, or secrets management;
  troubleshooting cluster issues.
  SKIP: CI/CD pipeline authoring (use cicd-pipelines); Terraform/IaC provisioning
  (use terraform-kafka-infra); application-level code patterns (use backend-patterns); self-healing architecture, MAPE-K control loops, autonomic computing theory, closed-loop/AIOps remediation design (use self-healing-systems-autonomic-computing).
version: "1.2.0"
updated: "2026-05-29"
category: developer
tags:
  - kubernetes
  - k8s
  - networking
  - CNI
  - cilium
  - calico
  - service-mesh
  - istio
  - linkerd
  - ingress
  - gateway-api
  - helm
  - kustomize
  - argocd
  - flux
  - gitops
  - prometheus
  - grafana
  - rbac
  - security
  - operators
  - statefulset
  - devops
  - platform-engineering
keywords:
  - kubernetes
  - kubectl
  - k8s
  - helm chart
  - kustomize
  - argocd
  - flux
  - cilium
  - calico
  - istio
  - linkerd
  - ingress controller
  - gateway api
  - networkpolicy
  - statefulset
  - kubernetes operator
  - pod security
  - CoreDNS
  - kube-proxy
  - kubeconfig
  - HPA
  - VPA
  - KEDA
  - pod disruption budget
  - topology spread constraints
  - sealed secrets
  - external secrets operator
whenToUse:
  - "Designing or debugging Kubernetes workloads (pods, deployments, StatefulSets)"
  - "Choosing or configuring a CNI plugin (Cilium vs Calico vs Flannel)"
  - "Setting up or troubleshooting service mesh (Istio, Linkerd, Cilium)"
  - "Configuring ingress, Gateway API, or load balancing"
  - "Writing NetworkPolicy or debugging connectivity between pods"
  - "Building or evaluating a Kubernetes Operator"
  - "Packaging manifests with Helm or Kustomize"
  - "Setting up GitOps with ArgoCD or Flux"
  - "Instrumenting monitoring with Prometheus, Grafana, or OpenTelemetry"
  - "Hardening cluster security (RBAC, PSS, admission controllers, secrets)"
  - "Troubleshooting pod startup, DNS resolution, or node health"
whenNotToUse:
  - "CI/CD pipeline authoring — use cicd-pipelines"
  - "Terraform or IaC provisioning — use terraform-kafka-infra"
  - "Application-level code patterns — use backend-patterns"
related_skills:
  - cicd-pipelines
  - terraform-kafka-infra
  - security-compliance-auditor
  - microservices-patterns
---

# Kubernetes & Networking

## When To Use This Skill

Use this skill when working on any of these topics:
- Designing or debugging Kubernetes workloads (pods, deployments, StatefulSets, DaemonSets, Jobs)
- Choosing or configuring a CNI plugin (Cilium, Calico, Flannel)
- Setting up or troubleshooting service mesh (Istio, Linkerd, Cilium Service Mesh)
- Configuring ingress, Gateway API, or load balancing
- Writing NetworkPolicy or debugging connectivity
- Building or evaluating Kubernetes Operators
- Packaging manifests with Helm or Kustomize
- Setting up GitOps with ArgoCD or Flux
- Instrumenting monitoring with Prometheus, Grafana, or OpenTelemetry
- Hardening cluster security (RBAC, PSS, admission controllers, secrets)
- Troubleshooting common cluster issues

**Quick decision trees are embedded in each section.** When unsure which tool or pattern to use, look for the "Choose X when..." guidance blocks.

---

## Overview

Kubernetes is a container orchestration platform built around a declarative control-plane model: you describe desired state via YAML manifests, and controllers continuously reconcile actual state toward it. The control plane consists of the API server, etcd (cluster state store), kube-scheduler, kube-controller-manager, and optionally cloud-controller-manager. Worker nodes run kubelet, kube-proxy (service networking via iptables/IPVS), and a container runtime (containerd, CRI-O).

As of 2025-2026, average CPU utilization in Kubernetes clusters is ~10% and memory ~23% — most clusters vastly overprovision. Production readiness requires proper resource requests/limits, autoscalers (HPA, VPA, KEDA), and cost visibility tooling. Kubernetes 1.30+ solidified CEL-based validation, sidecar container support as GA, and the Gateway API graduating as the preferred replacement for the Ingress API.

The ecosystem has stabilized around a core stack: Cilium or Calico for CNI, Helm + Kustomize for manifests, ArgoCD or Flux for GitOps, Prometheus + Grafana for observability, and either Vault or External Secrets Operator for secrets management.

---

## Core Concepts

### Pods
- The smallest schedulable unit; one or more containers sharing a network namespace and storage volumes.
- Never create bare Pods in production — they do not reschedule on node failure. Use Deployments, StatefulSets, DaemonSets, or Jobs instead.
- Each Pod gets a unique cluster-internal IP; containers within a Pod communicate via `localhost`.
- Key fields: `resources.requests/limits`, `livenessProbe`, `readinessProbe`, `startupProbe`, `securityContext`, `affinity/antiAffinity`, `topologySpreadConstraints`.

### Deployments
- Manages ReplicaSets to maintain N replicas; supports rolling update and recreate strategies.
- Rolling update parameters: `maxUnavailable` (default 25%) and `maxSurge` (default 25%) control rollout speed vs. availability.
- Always define `readinessProbe` — traffic only routes to pods that pass it.
- Create the backing Service *before* the Deployment so DNS resolves at startup.

### Services
- Provides a stable virtual IP (ClusterIP) and DNS name for a set of pods selected by label.
- Types: `ClusterIP` (internal only), `NodePort` (exposes on every node's port), `LoadBalancer` (provisions cloud LB), `ExternalName` (DNS CNAME alias).
- Headless Services (`clusterIP: None`) return individual pod IPs — required by StatefulSets for stable DNS per pod.

### Namespaces
- Logical partitioning of cluster resources; not a security boundary by themselves.
- Use for environment separation (dev/staging/prod) or team isolation with ResourceQuotas and LimitRanges.
- Apply Pod Security Admission labels at the namespace level (`pod-security.kubernetes.io/enforce: restricted`).

### ConfigMaps & Secrets
- ConfigMaps: non-sensitive config data; mounted as files or injected as env vars.
- Secrets: base64-encoded (NOT encrypted) by default in etcd; use encryption-at-rest config or an external secrets solution.
- Never store Secrets in Git without encryption. See Secrets Management section.

### RBAC
- Four objects: `Role`/`ClusterRole` (what), `RoleBinding`/`ClusterRoleBinding` (who gets what).
- Principle of least privilege: grant only the verbs and resources actually needed.
- Avoid `cluster-admin` bindings for workloads. Audit with `kubectl auth can-i --list`.
- Set `automountServiceAccountToken: false` on pods that don't call the API server.
- Prefer `Role` + `RoleBinding` (namespace-scoped) over `ClusterRole` bindings wherever possible.

### DaemonSets
- Runs exactly one pod per node (or per matching node subset via `nodeSelector`/`affinity`).
- Use for: CNI agents (Cilium, Calico), log collectors (Fluentd, Fluent Bit), node monitoring (node-exporter), security agents.
- Rolling update strategy: `maxUnavailable` controls how many nodes are updated simultaneously.
- Tolerate `NoSchedule` taints on control-plane nodes if the DaemonSet must run everywhere.

### PersistentVolumes & StorageClasses
- `volumeBindingMode: WaitForFirstConsumer` — delays PV provisioning until a pod is scheduled; prevents cross-zone binding mismatches.
- `reclaimPolicy: Retain` — PV persists after PVC deletion; must be manually reclaimed. Use for databases in production.
- CSI (Container Storage Interface) is the standard driver model; avoid deprecated in-tree volume plugins.
- Access modes: `ReadWriteOnce` (single node), `ReadWriteMany` (NFS/CephFS), `ReadWriteOncePod` (single pod — K8s 1.22+).

### Jobs & CronJobs
- **Job**: runs one or more pods to completion; retries on failure up to `backoffLimit`.
- **CronJob**: schedules Jobs on a cron schedule (`spec.schedule: "0 2 * * *"`).
- `completionMode: Indexed` (K8s 1.24+): each pod gets a unique index via `JOB_COMPLETION_INDEX` env var — useful for parallel batch processing with sharding.
- Set `activeDeadlineSeconds` to prevent runaway Jobs.
- Always set `ttlSecondsAfterFinished` to auto-clean completed Job pods (default: never cleaned up).

### ResourceQuota & LimitRange

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: team-quota
  namespace: team-a
spec:
  hard:
    requests.cpu: "10"
    requests.memory: 20Gi
    limits.cpu: "20"
    limits.memory: 40Gi
    pods: "50"
```

```yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: default-limits
  namespace: team-a
spec:
  limits:
  - type: Container
    default:
      cpu: 500m
      memory: 256Mi
    defaultRequest:
      cpu: 100m
      memory: 128Mi
    max:
      cpu: "4"
      memory: 4Gi
```

### HPA / VPA / KEDA
- **HPA**: scales replica count based on CPU, memory, or custom metrics.
- **VPA**: adjusts requests/limits per pod; cannot run simultaneously with HPA on the same metric.
- **KEDA**: scales to zero based on event sources (Kafka lag, queue depth, etc.); most flexible for async workloads.

---

## Networking

### Kubernetes Networking Model
- Every Pod gets a unique routable IP; pods on different nodes communicate without NAT.
- Services abstract pod IPs via kube-proxy (iptables/IPVS) or eBPF-based CNI bypass.
  For the eBPF datapath itself — Cilium's XDP/tc/socket hooks, kube-proxy replacement
  internals, Hubble flow visibility, and Tetragon runtime security — load
  `references/ebpf-observability.md` in this hub.
- DNS: CoreDNS resolves `<service>.<namespace>.svc.cluster.local`.

### CNI Plugins

**Cilium** (recommended for new clusters)
- Built on Linux eBPF — bypasses iptables entirely, replacing kube-proxy.
- P99 latency ~0.8ms (vs. Calico ~1.4ms in 2026 benchmarks); ~40-50% faster than iptables-based CNIs.
- Native NetworkPolicy enforcement, L7 policy (HTTP/gRPC-aware), Hubble for flow observability.
- Sidecarless mTLS via Cilium Service Mesh (kernel-level, avoids 10-15% CPU tax of sidecar proxies).

**Calico**
- Supports BGP routing (no overlay) for bare-metal and on-prem.
- eBPF dataplane in Calico v3.28+ narrows performance gap with Cilium to <5% in most scenarios.
- Better when: BGP integration is required, Windows nodes are needed, or operational familiarity with Calico exists.

**Flannel**
- Simplest CNI — VXLAN overlay, easy to operate.
- Does NOT support Kubernetes NetworkPolicy natively.
- Suitable only for small dev/test clusters; not recommended for production.

**Choose:**
- Cilium for new clusters, security-centric environments, or eBPF-forward teams
- Calico for BGP integration, large-scale enterprise, or Windows nodes
- Cloud-managed (AWS VPC CNI, GKE Dataplane V2, Azure CNI) when running EKS/GKE/AKS

### Service Mesh

| | Istio | Linkerd | Cilium Service Mesh |
|---|---|---|---|
| Architecture | Envoy sidecar + Istiod | Rust micro-proxy sidecar | eBPF kernel (no sidecar) |
| Latency overhead | High | Low | Lowest |
| Operational complexity | High | Low | Low (if Cilium already in use) |
| Memory overhead | High (~25-50GB more than Linkerd at 500 services) | Low | Minimal |

- **Choose Istio** when: advanced traffic management (fault injection, circuit breaking, sophisticated retries), dedicated platform team, or Envoy-specific filters.
- **Choose Linkerd** when: minimal operational burden, small team, or mesh that installs in minutes.
- **Choose Cilium Service Mesh** when: already using Cilium as CNI, want zero sidecar overhead.

### NetworkPolicy

Default behavior: all ingress and egress is allowed unless a NetworkPolicy selects the pod.

**Production default-deny pattern:**
```yaml
# Step 1: deny all ingress and egress
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
---
# Step 2: allow egress to CoreDNS (required for any DNS resolution)
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-dns-egress
spec:
  podSelector: {}
  policyTypes:
  - Egress
  egress:
  - ports:
    - port: 53
      protocol: UDP
    - port: 53
      protocol: TCP
```

Without the DNS egress allow, pods in a default-deny namespace cannot resolve any names.

### Ingress & Gateway API

**Ingress (legacy — frozen API)**
- Routes HTTP/HTTPS traffic to services based on host/path rules.
- Kubernetes project recommends migrating to Gateway API; Ingress API is feature-frozen.

**Gateway API (preferred)**
- Role-oriented: GatewayClass (infra team), Gateway (cluster ops), HTTPRoute/TCPRoute (app team).
- Supports traffic splitting, header-based routing, request mirroring natively.
- **ReferenceGrant gotcha**: cross-namespace routing requires a `ReferenceGrant` in the target namespace. Without it, the route silently fails.

```yaml
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferenceGrant
metadata:
  name: allow-frontend-route
  namespace: backend
spec:
  from:
  - group: gateway.networking.k8s.io
    kind: HTTPRoute
    namespace: frontend
  to:
  - group: ""
    kind: Service
```

**Ingress Controllers:**

| Controller | Best for | Notes |
|---|---|---|
| NGINX Ingress | Stability, large community | Most widely deployed |
| Traefik | Auto-discovery, Let's Encrypt | Great for dynamic environments |
| Envoy Gateway | Gateway API-first | Advanced L7, distributed tracing |
| HAProxy Ingress | High performance, TCP | Strong raw throughput |

### CoreDNS
- Resolves `<svc>.<ns>.svc.cluster.local` and `<pod-ip>.<ns>.pod.cluster.local`.
- Key plugins: `kubernetes`, `forward` (upstream DNS), `cache` (TTL caching), `health`, `ready`, `prometheus`.
- `ndots:5` default causes up to 5 DNS search-domain lookups for external names. Lower to `ndots:2` for most workloads.
- Scale CoreDNS replicas based on cluster size; use `DNSAutoscaler` or HPA on DNS query rate.

---

## Operators & StatefulSets

### StatefulSets
- Use for workloads requiring stable network identity, ordered scaling, and persistent storage (databases, Kafka, ZooKeeper, Elasticsearch).
- Each Pod gets a predictable DNS name: `<name>-<ordinal>.<headless-svc>.<ns>.svc.cluster.local`.
- VolumeClaimTemplates provision a unique PVC per Pod; PVCs are NOT deleted when the StatefulSet is scaled down.
- Always use `podDisruptionBudget` to protect against simultaneous disruptions during node drains.

### Operators
- Operators encode operational knowledge as Kubernetes controllers managing CRDs.
- In 2025, Operators are non-negotiable for production databases: they handle failover, backup, restore, scaling, and upgrades.
- Built with controller-runtime (Go) or Kopf (Python); use Server-Side Apply (SSA) for managing child objects.

**Operator reconcile pattern:**
1. Compute desired state from CR spec.
2. List current owned objects (StatefulSets, Services, PDBs, Secrets).
3. Apply desired objects via SSA with a consistent `fieldManager`.
4. Update CR `.status` to reflect observed state.
5. Return `ctrl.Result{RequeueAfter: ...}` for time-based reconciliation.

**Finalizer pattern (critical for cleanup):**
```go
if !cr.DeletionTimestamp.IsZero() {
    if err := r.cleanupExternalResources(ctx, cr); err != nil {
        return ctrl.Result{}, err
    }
    controllerutil.RemoveFinalizer(cr, "myoperator.example.com/cleanup")
    return ctrl.Result{}, r.Update(ctx, cr)
}
```

Without finalizers, the CR is deleted before the operator can clean up external resources.

---

## GitOps & Deployment

### Helm
- CNCF graduated; ~75% adoption among Kubernetes users (CNCF 2025 Survey). Helm 4 released November 2025.
- Best for: third-party software distribution, app packaging shared with other teams.
- Pin chart versions and image tags; avoid `latest`. Use `helm diff` plugin before upgrades.

### Kustomize
- Built into `kubectl` (`kubectl apply -k`); patch-based, no templating language.
- Best for: environment-specific configuration of internal apps (dev/staging/prod).

### Hybrid Pattern (recommended)
- **Helm for packaging** third-party charts and distributing internal apps.
- **Kustomize for environment overlays**, using `helmCharts` entries in `kustomization.yaml` to render Helm charts then apply patches.

### ArgoCD vs Flux CD

| | ArgoCD | Flux CD |
|---|---|---|
| UI | Built-in web dashboard | No built-in UI |
| Architecture | Centralized API server + controllers | Modular native Kubernetes controllers |
| RBAC | Built-in team RBAC | Kubernetes-native RBAC |
| Multi-cluster | Native app-of-apps pattern | Requires manual bootstrap per cluster |
| Onboarding | Faster (UI-driven) | Steeper (YAML-first) |

- **Choose ArgoCD** when: developer experience matters, managing many applications across clusters, or need team-based RBAC.
- **Choose Flux** when: automation-first, lightweight GitOps as Kubernetes-native controllers.

**ArgoCD ApplicationSet — multi-cluster GitOps:**
```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: guestbook
spec:
  generators:
  - clusters: {}
  template:
    spec:
      project: default
      source:
        repoURL: https://github.com/org/gitops-repo
        path: apps/guestbook/{{name}}
      destination:
        server: "{{server}}"
        namespace: guestbook
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
```

---

## Monitoring & Observability

### The Observability Stack (2025)
- **Metrics**: Prometheus + Grafana (or Grafana Mimir for long-term storage).
- **Logs**: Loki + Grafana or EFK.
- **Traces**: Tempo (Grafana stack) or Jaeger with OpenTelemetry SDK instrumentation.

### Prometheus PrometheusRule CRD

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: app-alerts
  labels:
    prometheus: kube-prometheus
spec:
  groups:
  - name: app.rules
    rules:
    - alert: PodCrashLooping
      expr: rate(kube_pod_container_status_restarts_total[15m]) > 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Pod {{ $labels.pod }} is crash-looping"
    - alert: HighCPUThrottling
      expr: rate(container_cpu_throttled_seconds_total[5m]) > 0.5
      for: 10m
      labels:
        severity: warning
```

---

## Security

### Pod Security Standards (PSS)
Three levels enforced by the built-in Pod Security Admission controller:
- **Privileged**: no restrictions (system components only).
- **Baseline**: prevents known privilege escalations; suitable for general workloads.
- **Restricted**: requires non-root, no host path mounts, drop all capabilities.

```yaml
metadata:
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/warn: restricted
    pod-security.kubernetes.io/audit: restricted
```

Minimal compliant `securityContext` for restricted namespaces:
```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 2000
  seccompProfile:
    type: RuntimeDefault
containers:
- securityContext:
    allowPrivilegeEscalation: false
    readOnlyRootFilesystem: true
    capabilities:
      drop: ["ALL"]
```

### Secrets Management

| Option | How it works | Git safety |
|---|---|---|
| Sealed Secrets (Bitnami) | Encrypt with cluster public key; decrypt in-cluster | Encrypted values committed |
| External Secrets Operator | Syncs from AWS SM / GCP SM / Vault / Azure KV | No secrets in Git |
| HashiCorp Vault | Dynamic secrets with TTL; Vault Agent Injector | No secrets in Git |

**Production rule:** secrets must be: (1) sourced from an external vault, (2) delivered via sync or mount, (3) encrypted at rest in etcd.

### Admission Controllers
- **OPA Gatekeeper**: policy-as-code using Rego; enforces custom constraints.
- **Kyverno**: Kubernetes-native policy engine; policies written as YAML; supports generate, mutate, validate.
- **ValidatingAdmissionPolicy**: native CEL-based validation without external controller (GA in 1.30).

```yaml
# Kyverno policy: require signed images
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-image-signature
spec:
  validationFailureAction: Enforce
  rules:
  - name: check-image-signature
    match:
      resources:
        kinds: [Pod]
    verifyImages:
    - imageReferences: ["registry.example.com/*"]
      attestors:
      - entries:
        - keyless:
            subject: "https://github.com/my-org/*"
            issuer: "https://token.actions.githubusercontent.com"
```

---

## Common Patterns

### Health Probes (always define all three)
```yaml
livenessProbe:
  httpGet: { path: /healthz, port: 8080 }
  initialDelaySeconds: 10
  periodSeconds: 10
readinessProbe:
  httpGet: { path: /ready, port: 8080 }
  initialDelaySeconds: 5
startupProbe:
  httpGet: { path: /healthz, port: 8080 }
  failureThreshold: 30
  periodSeconds: 10
```

### Graceful Shutdown
- Set `terminationGracePeriodSeconds` appropriately (default 30s is often too short for databases).
- Handle `SIGTERM` in app code; drain in-flight requests before exiting.
- Use `preStop` hook to allow kube-proxy to drain connections: `lifecycle.preStop.exec.command: ["sleep", "5"]`.

### Pod Disruption Budgets
```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
spec:
  minAvailable: 2
  selector:
    matchLabels: { app: my-app }
```

### Topology Spread Constraints
```yaml
topologySpreadConstraints:
- maxSkew: 1
  topologyKey: topology.kubernetes.io/zone
  whenUnsatisfiable: DoNotSchedule
  labelSelector:
    matchLabels: { app: my-app }
```

### Resource Right-Sizing
- Start with VPA in recommendation mode to observe actual usage before setting requests/limits.
- CPU: set requests based on average, limits at 2-4x (or omit CPU limits to avoid throttling spikes).
- Memory: set requests = limits for predictable Guaranteed QoS class on critical workloads.

---

## Anti-Patterns

| Anti-Pattern | Consequence | Fix |
|---|---|---|
| Bare Pods | Don't reschedule on node failure | Use Deployments/StatefulSets |
| No resource requests/limits | Noisy-neighbor issues, OOM kills | Set via LimitRange defaults |
| `latest` image tag | Non-reproducible deployments | Use immutable SHA or semver tags |
| Secrets in ConfigMaps or Git (unencrypted) | Credential exposure | Sealed Secrets or ESO |
| Overuse of `cluster-admin` | Broadest possible permissions | Least-privilege Roles |
| Single replica for stateful services | No HA | StatefulSet with ≥3 replicas |
| No PodDisruptionBudget on databases | Node drain evicts all replicas | Always set PDB |
| Skipping readinessProbes | Traffic routes to unready pods | Always define readinessProbe |
| Flannel in production without NetworkPolicy | No traffic isolation between tenants | Use Calico or Cilium |
| High `ndots` DNS setting | 5 extra lookup chains per external name | Lower to `ndots:2` |
| No seccompProfile | Unrestricted syscall surface | Set `RuntimeDefault` |
| Missing finalizers on Operators | Orphaned cloud resources on deletion | Always add finalizer pattern |
| Deployments for stateful workloads | Shared PVC access, unstable pod names | Use StatefulSets |

---

## Troubleshooting

### Pod Won't Start
```bash
kubectl describe pod <name> -n <ns>   # Events section shows scheduler/image/probe failures
kubectl logs <pod> --previous          # Logs from crashed container
kubectl get events --sort-by=.metadata.creationTimestamp -n <ns>
```
Common causes: `ImagePullBackoff` (bad image/registry auth), `CrashLoopBackOff` (app crash), `Pending` (insufficient resources), `OOMKilled` (memory limit too low).

### Service Not Reachable
```bash
kubectl get endpoints <svc-name> -n <ns>   # Must show pod IPs; empty = label selector mismatch
kubectl run debug --image=nicolaka/netshoot -it --rm -- curl http://<svc>.<ns>.svc.cluster.local
```

### DNS Resolution Failures
```bash
kubectl run debug --image=busybox -it --rm -- nslookup <svc>.<ns>.svc.cluster.local
kubectl logs -n kube-system -l k8s-app=kube-dns
```

### High CPU/Memory
```bash
kubectl top pods -n <ns> --sort-by=cpu
kubectl top nodes
kubectl describe node <node>   # Check Allocated resources vs Capacity
```

### NetworkPolicy Debugging
- With Cilium: `hubble observe --pod <name> --verdict DROPPED` for real-time flow visibility.
- With Calico: `calicoctl get policy` and Flow Logs.
- Use `netshoot` pod (`nicolaka/netshoot`) as a debug container.

### Node Not Ready
```bash
kubectl get nodes
kubectl describe node <node>
journalctl -u kubelet -n 100
```
Common causes: kubelet not running, disk pressure, memory pressure, CNI plugin crashed, certificate expiry (`kubeadm certs check-expiration`).

### Helm Release Stuck
```bash
helm list -n <ns> -a
helm rollback <release> -n <ns>
```

---

## References

- [Kubernetes Services, Load Balancing, and Networking](https://kubernetes.io/docs/concepts/services-networking/)
- [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
- [RBAC Good Practices](https://kubernetes.io/docs/concepts/security/rbac-good-practices/)
- [Cilium vs Calico vs Flannel CNI Comparison 2026](https://sanj.dev/post/cilium-calico-flannel-cni-performance-comparison/)
- [Service Meshes Decoded: Istio vs Linkerd vs Cilium (LiveWyer)](https://livewyer.io/blog/service-meshes-decoded-istio-vs-linkerd-vs-cilium/)
- [ArgoCD vs Flux CD GitOps 2025](https://www.zignuts.com/blog/argo-cd-vs-flux-cd--comparison)
- [Kubernetes Secrets Management 2025 (Infisical)](https://infisical.com/blog/kubernetes-secrets-management-2025)
- [Grafana Mimir 3.0 KubeCon 2025](https://www.morningstar.com/news/business-wire/20251105571677/)
