<!-- hub-reference-banner -->
> **Reference file — part of the `devops-linux-admin` hub.** Formerly the standalone `systemd` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

name: systemd
title: systemd Deep Dive — Units, journald, Resource Control, systemd-oomd, System Extensions
description: >
  systemd internals beyond day-to-day systemctl. The unit/object model and dependency+ordering semantics
  (Wants/Requires/Requisite/BindsTo/PartOf, After/Before, the .service/.socket/.target/.timer/.mount/.slice/
  .scope/.path/.device type taxonomy, service Type=, socket activation, generators, drop-ins, presets);
  systemd-journald (structured binary journal, Storage= volatile/persistent/auto/none, rotation/vacuum,
  rate-limiting, ForwardToSyslog, Forward Secure Sealing, journal namespaces, remote journals); cgroup-v2
  resource control (the slice/scope/service hierarchy, CPUWeight/CPUQuota/AllowedCPUs, the
  MemoryMin/Low/High/Max protection+throttle ladder, IOWeight, TasksMax, Delegate=); systemd-oomd
  (PSI-driven userspace OOM, ManagedOOMSwap / ManagedOOMMemoryPressure, oomd.conf); and image-based
  extensions (systemd-sysext for /usr+/opt, systemd-confext for /etc, DDIs/dm-verity, portable services
  via portablectl, and the execution-sandbox surface ProtectSystem/PrivateTmp/DynamicUser audited by
  systemd-analyze security).
category: developer
tags:
  - developer
  - devops
  - linux
  - systemd
  - cgroups
  - observability
---

# systemd Deep Dive

This reference goes past the `systemctl start/enable/status` and `journalctl -u` basics
that live in `references/linux-sysadmin.md` (§3.3–3.5) and into how systemd actually
models the system: the unit graph, the journal subsystem, the cgroup-v2 resource-control
tree, PSI-driven OOM management, and image-based system/configuration extensions plus the
execution sandbox. For routine service ops and OOM-killer triage, start in `linux-sysadmin`;
come here for *why* and for the advanced knobs.

`systemd` is the init system (PID 1) and service/session manager on most modern Linux
distros (RHEL 7+, Ubuntu 15.04+, Debian 8+, Fedora, SUSE, Arch). PID 1 is the manager;
`systemctl` is the client that talks to it over D-Bus. Authoritative source of truth is
the freedesktop.org man pages — defer to `systemd.unit(5)`, `systemd.exec(5)`,
`systemd.resource-control(5)`, `journald.conf(5)`, `oomd.conf(5)`, and `systemd-sysext(8)`
for exact flags and version-gated behavior.

---

## 1. The unit & object model (the shared foundation)

Everything systemd manages is a **unit** — a declarative config file describing an object
and its relationships. The unit graph (dependencies + ordering) is the foundation every
other subsystem here builds on, so understand it first.

### Unit types (suffix = type)

| Suffix | Object | Notes |
| --- | --- | --- |
| `.service` | A process/daemon | The workhorse; `Type=` governs readiness semantics |
| `.socket` | An IPC/network socket | Enables **socket activation** (lazy start, boot parallelism) |
| `.target` | A grouping/sync point | Replaces SysV runlevels (`multi-user.target`, `graphical.target`) |
| `.timer` | A cron-like trigger | Activates a matching `.service` on `OnCalendar=`/`OnBootSec=` |
| `.mount` / `.automount` | A mount point | Auto-generated from `/etc/fstab` by a generator |
| `.path` | Path-based activation | Starts a service when a file/dir appears or changes |
| `.slice` | A cgroup branch | Resource-control grouping (see §3) |
| `.scope` | Externally-forked processes | Like a service but for processes systemd did **not** fork (e.g. user sessions, `machinectl`) |
| `.device` | A udev device | Lets units depend on hardware appearing |
| `.swap` | A swap area | Generated from fstab |

### Service `Type=` (readiness protocol — the most common footgun)

- `simple` (default if `ExecStart=` set, no `Type=`): considered started the instant the
  binary is `fork`/`exec`'d. systemd does **not** know when the service is actually ready.
- `exec`: like `simple` but "started" once the `execve()` succeeds (catches bad binaries).
- `forking`: the classic double-fork daemon; use with `PIDFile=`.
- `oneshot`: runs to completion then exits (with `RemainAfterExit=yes` to stay "active");
  ideal for setup scripts and ordering anchors.
- `notify` / `notify-reload`: the service calls `sd_notify(READY=1)` — the **correct** way
  to signal readiness so ordered successors don't start too early.
- `dbus`: ready once it acquires a `BusName=`.

### Dependencies vs ordering — keep them separate

These are **orthogonal axes**. A dependency says "pull this unit in"; an ordering directive
says "sequence relative to it." Neither implies the other.

**Requirement (dependency) directives:**
- `Wants=` — soft/optional: pull the target in, but its failure does **not** fail us. The
  preferred default (declared via the `.wants/` symlink dir created by `enable`).
- `Requires=` — hard: if the dependency fails to start or is stopped, we're stopped too.
- `Requisite=` — like Requires but does **not** start it; fails immediately if it isn't
  already up.
- `BindsTo=` — stronger than Requires: if the bound unit stops *for any reason* (incl.
  device unplug), we stop too. Common with `.device` units.
- `PartOf=` — one-directional propagation of stop/restart (not start) — restart the parent,
  children restart; useful for a `.target` controlling a group.
- `Conflicts=` — negative dependency; starting one stops the other.

**Ordering directives:**
- `After=` / `Before=` — pure sequencing. `After=foo` means "do not start until foo's job
  finished." **Wants/Requires do NOT imply After** — without explicit ordering, units start
  in **parallel**. The idiomatic pair is `Wants=foo` + `After=foo`.

### `[Install]`, drop-ins, presets, generators

- `[Install]` (`WantedBy=`, `RequiredBy=`, `Also=`, `Alias=`) is only read by
  `systemctl enable/disable`; it creates the `.wants/`/`.requires/` symlinks. It has no
  effect at runtime — an enabled unit is just a symlink in a target's `.wants/` dir.
- **Drop-ins**: `/etc/systemd/system/foo.service.d/override.conf` overlays/extends a unit
  without editing the vendor file. `systemctl edit foo` creates one. Last-wins per setting;
  list settings reset with an empty assignment (`ExecStart=` then `ExecStart=...`).
- **Unit search precedence** (highest first): `/etc/systemd/system` → `/run/systemd/system`
  → `/usr/lib/systemd/system`. `systemctl cat foo` shows the effective merged unit.
- **Presets** (`systemctl preset`): vendor policy for whether units are enabled on install.
- **Generators**: programs run early at boot that synthesize units dynamically
  (`systemd-fstab-generator`, `systemd-gpt-auto-generator`, etc.).

---

## 2. systemd-journald — the structured journal

`systemd-journald` collects logs from the kernel, stdout/stderr of services, syslog, and
native structured `sd_journal` calls into an indexed **binary** journal keyed by trusted
fields (`_PID`, `_UID`, `_SYSTEMD_UNIT`, `_BOOT_ID`, …) that the daemon stamps and the
sender cannot forge.

### Storage (`Storage=` in `journald.conf`)

- `volatile` — memory only, under `/run/log/journal` (lost on reboot).
- `persistent` — on disk under `/var/log/journal` (auto-created); falls back to `/run` early
  in boot.
- `auto` (default on many distros) — persistent **iff** `/var/log/journal/` exists, else
  volatile. So "make logs survive reboot" = `mkdir -p /var/log/journal && systemctl restart
  systemd-journald` (or set `Storage=persistent`).
- `none` — drop all storage; forwarding (syslog/console/kmsg) still works.

### Sizing, rotation, retention

`SystemMaxUse=`, `SystemKeepFree=`, `SystemMaxFileSize=`, `MaxFileSec=`, `MaxRetentionSec=`
bound disk use (defaults: ~10% of fs, capped at 4 GiB). Manual reclaim:
`journalctl --vacuum-size=500M`, `--vacuum-time=7d`, `--disk-usage`.

### Rate limiting & forwarding

- `RateLimitIntervalSec=` / `RateLimitBurst=` drop floods per-service (a frequent cause of
  "missing" log lines under load).
- `ForwardToSyslog=`, `ForwardToKMsg=`, `ForwardToConsole=`, `ForwardToWall=` bridge to other
  sinks (set `ForwardToSyslog=yes` when an rsyslog/journald hybrid is in play).

### Forward Secure Sealing (FSS) & integrity

FSS (`journalctl --setup-keys`, then `Seal=yes`) cryptographically seals persistent journals
using Seekable Sequential Key Generators so past entries can't be silently altered;
`journalctl --verify` checks integrity. Tamper-evident, not tamper-proof.

### Journal namespaces & remote

- **Namespaces** (`LogNamespace=` on a unit + `journald@<ns>` instance) give a service an
  isolated journal store — useful for multi-tenant or high-volume separation.
- **Remote**: `systemd-journal-remote` / `-upload` / `-gatewayd` ship or serve journals over
  HTTP(S) for centralization without converting to text syslog.

Querying lives in `linux-sysadmin` §3.5 (`-u`, `-f`, `-b`, `-p`, `-k`, `-o json`, field
matches like `_PID=`/`_SYSTEMD_UNIT=`).

---

## 3. Resource control — cgroup v2

systemd is a **cgroup manager**: every service/scope gets its own cgroup, organized into a
tree of **slices**. On the unified (cgroup-v2) hierarchy, control is *hierarchical* — limits
on a parent slice bound the sum of its children.

### The slice tree

- `-.slice` (root) → `system.slice` (system services), `user.slice` (→ `user-<UID>.slice` →
  `session/app/...`), `machine.slice` (VMs/containers via machined).
- A unit joins a slice with `Slice=` (default `system.slice`). Put related services under a
  custom `Slice=my.slice` to cap them collectively.
- Inspect live with `systemd-cgls` (tree) and `systemd-cgtop` (live resource use).

### CPU

- `CPUWeight=` (1–10000, default 100) — *relative* share under contention. 200 ≈ twice the
  CPU of a 100-weight sibling **only when both are runnable and competing**.
- `CPUQuota=` — *absolute* cap, e.g. `CPUQuota=50%` = half a core; `200%` = two cores.
- `AllowedCPUs=` / `AllowedMemoryNodes=` — pin to specific CPUs/NUMA nodes (cpuset).

### Memory — the protection→throttle→kill ladder

- `MemoryMin=` — **hard** protection: memory below this is never reclaimed (guarantees
  working set survives system pressure).
- `MemoryLow=` — **soft** protection: reclaimed only if unprotected memory elsewhere is
  exhausted.
- `MemoryHigh=` — **throttle** (recommended primary knob): above it the cgroup is
  aggressively reclaimed and processes throttled, but **not** killed. Back-pressure.
- `MemoryMax=` — **hard limit / last line of defense**: exceed it and the cgroup OOM-killer
  fires. Pattern: set `MemoryHigh` as the working ceiling and `MemoryMax` a bit higher as the
  backstop.
- `MemorySwapMax=` — cap swap for the cgroup.

### IO (requires bfq/io controller on the backing block device)

- `IOWeight=` (1–10000, default 100) — relative bandwidth share.
- `IOReadBandwidthMax=` / `IOWriteBandwidthMax=` / `IOReadIOPSMax=` / `IOWriteIOPSMax=` —
  absolute per-device caps: `IOReadBandwidthMax=/dev/sda 50M`.

### Tasks & delegation

- `TasksMax=` — cap the number of tasks (pids) — a fork-bomb guard (default from
  `DefaultTasksMax`, often `15%` of `kernel.pid_max`).
- `Delegate=yes` — hand a subtree of the cgroup hierarchy to the unit's own manager so it can
  create sub-cgroups (this is how a rootless container runtime or a nested systemd manages
  its own children). Without delegation, only systemd may write the cgroup tree.

Most of these are runtime-settable: `systemctl set-property foo.service MemoryMax=1G`
(persists as a drop-in; add `--runtime` for transient).

---

## 4. systemd-oomd — PSI-driven userspace OOM killer

The kernel OOM killer fires only when memory is *already* exhausted and picks victims by a
crude heuristic, often after the box has thrashed into unresponsiveness. **systemd-oomd**
acts *earlier* and *smarter* from userspace using **PSI (Pressure Stall Information)** —
the kernel metric (Linux ≥ 4.20) that reports the % of wall-time tasks stalled waiting on
memory or IO.

### How it works

`systemd-oomd` polls PSI and cgroup memory/swap stats for monitored cgroups. When a
configured pressure/swap threshold is sustained over a window, it kills the *cgroup* (the
whole offending unit, not a lone process) that's the best candidate — so a runaway service
dies cleanly instead of taking a random victim down with it.

### Enabling monitoring (opt-in per cgroup)

- `ManagedOOMSwap=kill` — act when **system-wide swap** runs low. Because it's system-wide,
  it makes the most sense set on `-.slice` (root) with descendants eligible as candidates.
- `ManagedOOMMemoryPressure=kill` — act when a cgroup's **memory PSI** exceeds the limit.
- `ManagedOOMMemoryPressureLimit=` — per-unit override of the global pressure threshold.

### Global config — `oomd.conf`

Loaded from the first of `/etc/systemd/`, `/run/systemd/`, `/usr/local/lib/systemd/`,
`/usr/lib/systemd/`. Keys: `SwapUsedLimit=` (default ~90%), `DefaultMemoryPressureLimit=`
(default ~60%), `DefaultMemoryPressureDurationSec=` (how long pressure must persist before
acting). Enabled by default on Fedora 34+ and Ubuntu 22.04+. **Requirements**: cgroup-v2 +
kernel PSI support; it complements, doesn't replace, the kernel OOM killer and `MemoryMax`.

---

## 5. System & configuration extensions, portable services, sandboxing

### systemd-sysext / systemd-confext (image-based extensions)

For immutable / image-based OSes (and increasingly elsewhere), you can layer additional
files onto a read-only base without rebuilding the image:

- **sysext** extends `/usr/` (OS vendor tree) and `/opt/` (third-party). Images live in
  `/var/lib/extensions/`, `/etc/extensions/`, `/run/extensions/`.
- **confext** extends `/etc/` only. Images live in `/var/lib/confexts/`, etc.
- Mechanism: each image is **overlayed via overlayfs** onto the target hierarchy at merge
  time; the merged dirs become read-only. `unmerge` removes the overlay; the base reappears
  untouched.
- Each image (a dir, a `.raw` **DDI** = Discoverable Disk Image, optionally dm-verity-signed)
  must carry `/usr/lib/extension-release.d/extension-release.<NAME>` (sysext) or
  `/etc/extension-release.d/extension-release.<NAME>` (confext). This os-release-style file
  must match the host's `ID=` (or `_any`) and `VERSION_ID=`/`SYSEXT_LEVEL=`/`CONFEXT_LEVEL=`,
  plus `ARCHITECTURE=`. Mismatch ⇒ refused (prevents loading an extension built for a
  different OS version).
- Commands: `systemd-sysext merge | unmerge | refresh | list | status` (and the
  `systemd-sysext.service` that auto-merges at boot). `systemd-sysupdate` handles A/B image
  updates of the base + extensions.

### Portable services (`portablectl`)

Since v239, a "portable service" bundles a service, its binaries, and dependencies into an
image (dir, raw, or squashfs) that runs **from** the image. `portablectl attach <image>`
(talking to `systemd-portabled`) drops generated units + drop-ins onto the host and applies
a **security profile** (`default`, `nonetwork`, `strict`, `trusted`) that pre-sets
sandboxing. `detach` reverses it. It's a lighter-weight middle ground between a raw service
and a full container.

### Execution sandbox (systemd.exec hardening)

Any service can be locked down with kernel-backed sandboxing directives — these are the same
levers portable-service profiles set:

- **Filesystem**: `ProtectSystem=strict` (whole FS read-only except explicit paths),
  `ProtectHome=`, `ReadWritePaths=`/`ReadOnlyPaths=`/`InaccessiblePaths=`, `PrivateTmp=yes`
  (private `/tmp`), `PrivateDevices=`, `ProtectKernelTunables=`, `ProtectKernelModules=`,
  `ProtectControlGroups=`.
- **Identity**: `DynamicUser=yes` (transient UID/GID per run — no static user needed),
  `User=`/`Group=`, `SupplementaryGroups=`.
- **Capabilities & syscalls**: `CapabilityBoundingSet=`, `AmbientCapabilities=`
  (e.g. grant only `CAP_NET_BIND_SERVICE` instead of running as root for port < 1024),
  `NoNewPrivileges=yes`, `SystemCallFilter=@system-service`, `RestrictAddressFamilies=`,
  `MemoryDenyWriteExecute=yes`.
- **Namespaces**: `PrivateNetwork=`, `PrivateUsers=`, `RestrictNamespaces=`.

**Audit it**: `systemd-analyze security <unit>` (added v240) scores each service 0–10 (lower
= more locked down) and itemizes which directives are unset — the fastest way to harden a
unit incrementally.

---

## Practical patterns

- **Readiness, not race**: long-init daemons should use `Type=notify` + `sd_notify(READY=1)`
  so ordered successors (`After=`) actually wait. Don't paper over races with `sleep` in
  ExecStartPre.
- **Reliable restarts**: `Restart=on-failure`, `RestartSec=`, and `StartLimitIntervalSec=` /
  `StartLimitBurst=` to avoid crash-loop hammering; `systemctl reset-failed` to clear.
- **Cap a noisy service**: `systemctl set-property foo.service MemoryHigh=2G MemoryMax=3G
  CPUQuota=150%` — immediate, persisted as a drop-in.
- **Group + cap a workload**: define `my.slice` with limits, set `Slice=my.slice` on each
  member; cgroup-v2 enforces the aggregate.
- **Survive reboots for logs**: `Storage=persistent` (or create `/var/log/journal`).
- **Override vendor units the right way**: `systemctl edit foo` (drop-in), never edit
  `/usr/lib/systemd/system/...`.

## Anti-patterns

- Assuming `Requires=`/`Wants=` orders units — it does **not**; add `After=`.
- Using `Type=simple` for a daemon that needs warm-up, then chaining dependents that start
  too early.
- Hand-editing vendor unit files in `/usr/lib` (lost on package upgrade) instead of drop-ins.
- Relying solely on the kernel OOM killer for a memory-hungry service when `MemoryMax=` +
  `MemoryHigh=` + systemd-oomd would contain it before the host thrashes.
- Setting `MemoryMax=` with no `MemoryHigh=` — you get a hard kill with no graceful
  back-pressure warning.
- Editing files inside a **merged** sysext overlay expecting persistence — the overlay is
  read-only and ephemeral; change the image and `refresh`.
- Running a service as root "to be safe" — it's the opposite; `DynamicUser=` +
  `CapabilityBoundingSet=` + `systemd-analyze security` is safer.

## Troubleshooting

- `systemctl status foo` / `journalctl -u foo -b` — first stop; read the `Active:`/`Result:`
  lines and recent journal.
- `systemctl list-dependencies foo` and `systemd-analyze critical-chain foo` — why a unit
  started late or didn't start; `systemd-analyze blame` / `plot` for boot-time regressions.
- `systemctl cat foo` — see the *effective* merged unit (vendor + drop-ins).
- `systemd-cgls` / `systemd-cgtop` — which cgroup is eating CPU/mem; confirms `Slice=`/limits
  took effect.
- `systemctl show foo -p MemoryMax -p CPUQuotaPerSecUSec` — read back the live, resolved
  resource values.
- OOM forensics: `journalctl -k | grep -i oom` (kernel killer) vs
  `journalctl -u systemd-oomd` (userspace oomd actions) — they're different killers.
- `systemd-analyze verify foo.service` — lint a unit file for errors before deploying.
- Extension won't merge: check `extension-release.<NAME>` `ID=`/`VERSION_ID=` matches the host.

## References

- systemd.unit(5), systemd.service(5), systemd.target(5), systemd.socket(5) — freedesktop.org
  man pages (latest). https://www.freedesktop.org/software/systemd/man/latest/systemd.unit.html
- "systemd: Unit dependencies and order" — Fedora Magazine.
  https://fedoramagazine.org/systemd-unit-dependencies-and-order/
- journald.conf(5) — freedesktop.org.
  https://www.freedesktop.org/software/systemd/man/latest/journald.conf.html
- systemd.resource-control(5) — freedesktop.org / man7.
  https://man7.org/linux/man-pages/man5/systemd.resource-control.5.html
- RHEL 8/9 "Configuring resource management by using cgroups-v2 and systemd" — Red Hat docs.
  https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/8/html/managing_monitoring_and_updating_the_kernel/assembly_configuring-resource-management-using-systemd_managing-monitoring-and-updating-the-kernel
- systemd-oomd.service(8) + oomd.conf(5) — man pages; Fedora "EnableSystemdOomd" change.
  https://man7.org/linux/man-pages/man8/systemd-oomd.service.8.html ,
  https://fedoraproject.org/wiki/Changes/EnableSystemdOomd
- systemd-sysext(8) / systemd-confext(8) — Arch/Debian man pages; UAPI Group "Extension
  Images" spec. https://uapi-group.org/specifications/specs/extension_image/
- PORTABLE_SERVICES.md — systemd repo.
  https://github.com/systemd/systemd/blob/main/docs/PORTABLE_SERVICES.md
- systemd/Sandboxing — ArchWiki; systemd-analyze security.
  https://wiki.archlinux.org/title/Systemd/Sandboxing
