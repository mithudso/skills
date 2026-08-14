---
name: devops-linux-admin
description: >-
  Linux/macOS sysadmin & host-operations sub-hub (devops family). TRIGGER: sysadmin troubleshooting (high CPU, OOM, disk full, DNS, port conflicts), systemd units/journald, package management, Bash/Zsh production shell scripting, host network diagnostics (ip, ss, dig, tcpdump, nmap, iptables/nftables); macOS networking stack & diagnostics — scutil/SystemConfiguration, networksetup/locations, mDNSResponder & DNS resolution order, /etc/resolver, pfctl/socketfilterfw, Network Extension/VPN utun & split DNS, wdutil/sysdiagnose Wi-Fi tooling, macOS↔Linux CLI mapping, captive portal/Private Relay/AWDL failure modes. SKIP: kernel/memory/storage/namespaces internals → devops-linux-internals; Docker/Kubernetes/CI-CD/IaC → devops-containers-cicd; logging/tracing/metrics/eBPF/perf → devops-observability; DNS protocol internals (DNSSEC, DoH/DoT mechanics, delegation) → networking.
version: "1.1.0"
updated: "2026-06-12"
origin: local
---

# devops-linux-admin

Linux/macOS sysadmin & host-operations sub-hub (devops family).

This hub routes to on-demand reference files under `references/`. See each spoke for depth.

## Sub-skill routing table

| Spoke | Reference file | Covers |
| --- | --- | --- |
| linux-sysadmin | `references/linux-sysadmin.md` | Linux-side network diagnostics (ip, ss, dig, tcpdump) + cross-platform sysadmin playbooks: monitoring, processes, logs, disks, permissions. macOS network topics → macos-networking |
| systemd | `references/systemd.md` | systemd internals: unit model, dependencies/ordering, journald |
| linux-package-management | `references/linux-package-management.md` | apt/dnf/pacman ecosystems: install, query, repair, pin |
| shell-scripting | `references/shell-scripting.md` | Production Bash/Zsh: error handling, parameter expansion, debugging |
| linux-networking-stack | `references/linux-networking-stack.md` | Linux dataplane engineering: nftables, network namespaces, tc, XDP |
| macos-networking | `references/macos-networking.md` | macOS networking stack & diagnostics: SystemConfiguration/scutil, networksetup, mDNSResponder & DNS resolution order, /etc/resolver, PF/pfctl & application firewall, Network Extension VPNs/utun/split DNS, wdutil/sysdiagnose, macOS↔Linux CLI map, captive portal/Private Relay/AWDL failure modes |

<!-- cross-hub-map -->
## Cross-hub map — where every devops topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `devops-linux-internals` | devops-linux-internals | `references/linux-kernel-architecture.md`, `references/linux-boot-init.md`, `references/linux-memory-numa.md`, `references/linux-storage-filesystems.md`, … |
| `devops-linux-admin` | devops-linux-admin | `references/linux-sysadmin.md`, `references/systemd.md`, `references/linux-package-management.md`, `references/shell-scripting.md`, `references/macos-networking.md`, … |
| `devops-containers-cicd` | devops-containers-cicd | `references/docker-containers.md`, `references/kubernetes-networking.md`, `references/cicd-pipelines.md`, `references/terraform-kafka-infra.md`, … |
| `devops-observability` | devops-observability | `references/nodejs-observability.md`, `references/pino-structured-logging.md`, `references/sentry-monitoring.md`, `references/ebpf-observability.md`, … |
