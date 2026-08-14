<!-- hub-reference-banner -->
> **Reference file — part of the `devops-observability` hub.** Formerly the standalone `linux-perf-tracing` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

name: linux-perf-tracing
title: Linux Performance Analysis & Tracing — perf, ftrace, USE Method, PSI
description: >
  Kernel-level Linux performance methodology and the tracing tools that implement it. Covers the USE method
  (utilization / saturation / errors checklist for every resource), perf (hardware/software counters,
  sampling profiles, call graphs, flame graphs, perf sched/trace/probe), ftrace (the in-kernel function
  tracer via tracefs, function/function_graph tracers, tracepoints, kprobe/uprobe dynamic events,
  histograms, trace-cmd front-end), and PSI (Pressure Stall Information — /proc/pressure and cgroup
  cpu/memory/io.pressure, the some-vs-full distinction). Use to find a bottleneck systematically, profile
  CPU/off-CPU time, trace kernel functions and tracepoints, or read pressure metrics.
---

# Linux Performance Analysis & Tracing — perf, ftrace, USE Method, PSI

## Overview

This reference is about *finding and characterizing* Linux performance problems at the kernel level. It pairs
a methodology (the USE method — what to measure, in what order) with the three native instruments that feed
it: **perf** (the swiss-army profiler built on perf_events / PMCs), **ftrace** (the in-kernel function and
event tracer exposed through tracefs), and **PSI** (pressure stall information — the canonical "how much
time was lost waiting for this resource" signal). It deliberately sits *above* the everyday cheat-sheet
tools (`vmstat`, `iostat`, `sar`, `top`) covered in `references/linux-sysadmin.md`, and *beside* the modern
programmable-tracing layer (`bpftrace`/`bcc`/eBPF) in `references/ebpf-observability.md`. Reach for eBPF when
you need custom in-kernel aggregation; reach for ftrace/perf-probe when a built-in tracer or a quick dynamic
probe is enough.

## Core Concepts

### 1. The USE Method (Brendan Gregg)

A methodology for a *complete* health check of every resource. For each resource, check three metric types:

- **Utilization** — the percent of time the resource was busy servicing work (or, for capacity resources like
  memory/disk space, the proportion consumed).
- **Saturation** — the degree of *extra* queued work that could not be serviced immediately (run-queue length,
  swap activity, I/O queue depth). Saturation is the early-warning signal that utilization alone hides.
- **Errors** — count of error events (NIC drops, disk read errors, ECC). **Check errors first** — they are
  fastest to interpret and often the actual cause.

Build a resource × metric matrix and walk it. It "solves ~80% of server issues with ~5% of the effort."
Resources to enumerate: CPU, memory capacity, network interfaces, storage I/O, storage capacity, and (often
forgotten) interconnects/buses, I/O controllers, and memory bandwidth.

The Linux checklist (the canonical mapping — keep this table; it is the heart of the method):

| Resource | Utilization | Saturation | Errors |
|---|---|---|---|
| CPU | `vmstat 1` (us+sy+st); `mpstat -P ALL 1`; `sar -u`; per-proc `pidstat 1` | `vmstat 1` (`r` > nCPU); `sar -q` (runq-sz); `perf sched latency`; `/proc/PID/schedstat` | `perf` processor-specific error events (PMCs) |
| Memory cap | `free -m`; `vmstat 1` (free); `sar -r`; per-proc `top` RES, `slabtop -s c` | `vmstat 1` (`si`/`so`); `sar -B` (pgscank/pgscand); `sar -W`; OOM via `dmesg \| grep -i kill` | `dmesg` (HW); dynamic-trace malloc failures |
| Network iface | `sar -n DEV 1`; `ip -s link`; `/proc/net/dev`; `nicstat` | `ifconfig`/`ip -s link` (overruns, dropped); `netstat -s` (retransmits); `sar -n EDEV`; `nicstat` Sat | `ip -s link`; `netstat -i`; `sar -n EDEV` |
| Storage I/O | `iostat -xz 1` (`%util`); `sar -d`; `iotop`; `pidstat -d` | `iostat -xz 1` (`aqu-sz` > 1, high `await`); `sar -d` | `/sys/.../ioerr_cnt`; `smartctl`; block-layer trace |
| Storage cap | `df -h`; `swapon -s`; `/proc/meminfo` | (n/a once full → ENOSPC) | `strace`/dynamic trace for ENOSPC; `/var/log/messages` |

Modern note: `aqu-sz` is the current `iostat` column name (older docs say `avgqu-sz`); `dstat` is deprecated
upstream (use `dool`); prefer PSI (below) over `vmstat`'s `r`/`si`/`so` heuristics where the kernel supports it.

### 2. perf — the perf_events profiler

`perf` is the front-end to the kernel's perf_events subsystem. It samples on **events**: hardware PMU counters
(cycles, instructions, cache-misses, branch-misses, LLC, stalled-cycles), software events (cpu-clock,
context-switches, page-faults, migrations), tracepoints, and dynamic kprobe/uprobe events.

- **`perf list`** — enumerate available events on this CPU/kernel.
- **`perf stat <cmd>`** — counting mode. Prints a summary of events over the run. First step in any
  investigation: tells you if the workload is CPU-bound, memory-bound, or branch-bound. `perf stat -d` adds
  cache/IPC detail; IPC (instructions-per-cycle) < ~1.0 usually means it is stalled on memory.
- **`perf record -g <cmd>`** / `perf record -g -p <pid> -- sleep 30` — sampling mode with call graphs.
  Use `--call-graph dwarf` when the target lacks frame pointers (most distro binaries) or `--call-graph lbr`
  on modern Intel; default `fp` needs `-fno-omit-frame-pointer` builds. `-F 99` sets 99 Hz (a prime, to avoid
  lockstep). `-a` profiles system-wide.
- **`perf report`** — interactive TUI of the recorded `perf.data`. `perf script` dumps raw samples.
- **`perf top`** — live `top`-like sampling profile.
- **`perf sched record` / `perf sched latency`** — scheduler latency / run-queue wait analysis (CPU saturation).
- **`perf trace`** — strace-like syscall tracer built on perf_events (lower overhead than `strace`).
- **`perf probe`** — define dynamic kprobes/uprobes (`perf probe --add 'tcp_sendmsg'`) then sample them.

### 3. Flame graphs

A flame graph is a visualization of sampled stack traces: x-axis = population (alphabetical, **not** time),
width = how often that stack was sampled, y-axis = stack depth. Wide frames = hot paths. Two paths to one:

- Native: `perf record -g … && perf script report flamegraph` (newer perf builds emit an interactive SVG/HTML).
- Classic (Brendan Gregg's FlameGraph scripts): `perf script | stackcollapse-perf.pl | flamegraph.pl > out.svg`.

Variants: **CPU** flame graphs (on-CPU, from sampled stacks), **off-CPU** flame graphs (time *blocked*, from
sched tracepoints/eBPF — complements CPU graphs for I/O- or lock-bound workloads), and **differential** flame
graphs (red/blue delta between two profiles).

### 4. ftrace — the in-kernel function tracer

ftrace is the kernel's built-in tracing infrastructure (by Steven Rostedt). It powers function tracing, trace
events, dynamic kprobe/uprobe events, and histograms. It is controlled entirely through the **tracefs**
filesystem, mounted at `/sys/kernel/tracing` (legacy: `/sys/kernel/debug/tracing`). No userspace agent needed.

Key tracers (write to `current_tracer`):

- **`function`** — traces every kernel function entry (cheap, broad).
- **`function_graph`** — call-graph with entry/exit and per-function durations (the most useful for "why is
  this slow"). Read the `DURATION` column; lines marked `+`/`!` flag latency outliers.
- **`nop`** — disables function tracing but lets trace events / dynamic events still fire.
- **`irqsoff` / `preemptoff` / `wakeup` / `wakeup_rt`** — latency tracers (max-latency-so-far recorders).

The raw tracefs workflow:

```sh
cd /sys/kernel/tracing
echo function_graph > current_tracer
echo 'tcp_*' > set_ftrace_filter      # limit scope (essential — full trace is firehose)
echo 1 > tracing_on ; sleep 1 ; echo 0 > tracing_on
cat trace                              # or trace_pipe to stream + drain
echo nop > current_tracer ; echo > trace   # reset
```

**Trace events** (static tracepoints) live under `events/` (e.g. `events/sched/sched_switch/enable`).
**Dynamic events**: `kprobe_events`/`uprobe_events` add probes at arbitrary addresses. **Histograms** (the
`hist` trigger) aggregate in-kernel, e.g. latency distributions, without dumping every event to userspace.

### 5. trace-cmd — the ftrace front-end

`trace-cmd` is the userspace CLI wrapper around ftrace (avoids hand-editing tracefs). Core verbs:

- `trace-cmd record -p function_graph -l 'tcp_*' sleep 1` → writes `trace.dat`.
- `trace-cmd record -e sched -e net sleep 1` → record by event subsystem.
- `trace-cmd report` → human-readable dump; `trace-cmd list` enumerates tracers/events/functions.
- `trace-cmd start/stop/show/reset/clear` → control without recording to file.
- KernelShark is the GUI viewer for `trace.dat`.

### 6. PSI — Pressure Stall Information

PSI (kernel ≥ 4.20, `CONFIG_PSI=y`) quantifies *time lost to resource contention*. System-wide it appears as
three files under `/proc/pressure/`: **cpu**, **memory**, **io**. Per-cgroup (cgroup v2) it appears as
`cpu.pressure`, `memory.pressure`, `io.pressure` in each cgroup directory (same format).

Format and semantics:

```
some avg10=0.00 avg60=0.00 avg300=0.00 total=12345
full avg10=0.00 avg60=0.00 avg300=0.00 total=6789
```

- **`some`** — share of wall-clock time in which *at least one* task was stalled waiting on the resource.
  Early indicator of contention. (CPU only exposes `some`.)
- **`full`** — share of time in which *all non-idle* tasks were stalled (nobody made progress). This is the
  productivity-loss signal — high `full` memory/io pressure is a strong sign you are thrashing or OOM-bound.
- `avg10/60/300` are running-average percentages over 10s/1m/5m; `total` is cumulative stall microseconds.

PSI is what powers **systemd-oomd** and `oomd` (act on memory pressure *before* the kernel OOM killer fires).
You can also register **PSI triggers** via `poll()` on the pressure files to get woken when pressure crosses a
threshold (e.g. "memory some > 150ms over 1s"). See `references/systemd.md` and
`references/linux-cgroups-namespaces.md` for the cgroup-controller/oomd angles.

## Tools / Frameworks

- **perf** (`linux-tools-common` / `perf` pkg) — sampling profiler, counters, sched/trace/probe.
- **ftrace + trace-cmd + KernelShark** — built-in function/event tracer and its CLI/GUI.
- **FlameGraph** (github.com/brendangregg/FlameGraph) — stackcollapse + flamegraph.pl.
- **PSI** (`/proc/pressure/*`, cgroup `*.pressure`) — pressure metrics; **systemd-oomd**/**oomd** consume it.
- Adjacent: **bpftrace/bcc** (`references/ebpf-observability.md`) for programmable tracing; `pidstat`,
  `sar`/`sysstat`, `vmstat`, `iostat` (`references/linux-sysadmin.md`) feed the USE matrix.

## Methodology

1. **Start with USE, errors first.** Walk the resource × metric matrix; check error counters before
   utilization. Most issues surface in the easy metrics (CPU saturation, memory capacity, NIC util, disk util).
2. **Confirm the bottleneck class with `perf stat`.** CPU-bound (high util, IPC ≥ 1) vs memory-bound (low IPC,
   high cache-miss) vs off-CPU (low util but slow — go to off-CPU analysis).
3. **For on-CPU: `perf record -g -F 99` → flame graph.** Find the widest tower; that is the hot path.
4. **For off-CPU / latency: PSI + `perf sched latency` + off-CPU flame graph** (sched tracepoints or eBPF).
5. **To explain *why* a specific kernel path is slow: `function_graph` (scoped with `set_ftrace_filter`) or a
   targeted `perf probe`/kprobe.** Scope first — never trace the whole kernel unfiltered.
6. **Watch PSI for trend/contention** rather than instantaneous utilization; alert on `full` memory/io pressure.

## Practical Patterns

- 30-second CPU profile of one process: `perf record -F 99 -g -p $PID -- sleep 30 && perf script | \
  stackcollapse-perf.pl | flamegraph.pl > cpu.svg`.
- Cheap "what is the kernel doing right now": `perf top` (or `perf top -e cache-misses`).
- Syscall survey with low overhead: `perf trace -p $PID` (preferred over `strace` for hot processes).
- Function-graph why-is-this-slow: scope `set_ftrace_filter` to the subsystem, use `function_graph`, read
  `DURATION`.
- Pressure-based alerting: poll `/proc/pressure/memory` for `full avg10` rising; let systemd-oomd act first.

## Anti-Patterns

- **Trusting utilization alone.** A CPU at 90% util may be fine; one at 50% util with a deep run queue
  (saturation) is the real fire. Always pair util with saturation (PSI `some`/`full`, `r`).
- **Unfiltered ftrace `function` tracing on a busy box.** The trace buffer fills instantly and overhead spikes.
  Always set `set_ftrace_filter`/`set_ftrace_pid` or use `function_graph` on a narrow scope.
- **`perf record` with default frame-pointer call graphs on distro binaries** → broken/short stacks. Use
  `--call-graph dwarf` or `lbr`.
- **Sampling at round frequencies (100/1000 Hz).** Use a prime like 99/997 to avoid aliasing with periodic
  kernel activity.
- **Reading PSI `some` as catastrophic.** `some` rising is normal under any load; it is `full` (especially
  memory/io) that signals lost productivity.
- **Reaching for eBPF when ftrace/perf-probe suffices** (and vice-versa — using ftrace for custom aggregation
  that a bpftrace histogram does in one line). See `references/ebpf-observability.md`.

## Troubleshooting

- `perf` shows `<not supported>` for HW events → running in a VM/container without PMU passthrough, or
  `kernel.perf_event_paranoid` too high (`sysctl -w kernel.perf_event_paranoid=1`, or `-1` for full access).
- Stacks are `[unknown]` / truncated → missing frame pointers; switch to `--call-graph dwarf`, install debug
  symbols, or add `--no-children`.
- ftrace `trace` is empty → check `tracing_on`, that `current_tracer` ≠ `nop`, and that `set_ftrace_filter`
  didn't filter everything out; confirm tracefs is mounted.
- `/proc/pressure/` missing → kernel built without `CONFIG_PSI` or booted with `psi=0`; boot with `psi=1`.
- Per-cgroup `*.pressure` missing → not on cgroup v2 unified hierarchy (see `references/linux-cgroups-namespaces.md`).

## References

- USE Method (overview + ACM Queue article): https://www.brendangregg.com/usemethod.html ; Linux checklist: https://www.brendangregg.com/USEmethod/use-linux.html
- CPU Flame Graphs: https://www.brendangregg.com/FlameGraphs/cpuflamegraphs.html ; FlameGraph repo: https://github.com/brendangregg/FlameGraph
- perf tutorial (perfwiki): https://perfwiki.github.io/main/tutorial/
- RHEL "Getting started with flamegraphs": https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/8/html/monitoring_and_managing_system_status_and_performance/getting-started-with-flamegraphs_monitoring-and-managing-system-status-and-performance
- ftrace — Function Tracer (kernel.org): https://docs.kernel.org/trace/ftrace.html
- trace-cmd man page (man7): https://man7.org/linux/man-pages/man1/trace-cmd.1.html ; intro: https://devkernel.io/posts/ftrace_trace_cmd/
- PSI — Pressure Stall Information (kernel.org): https://docs.kernel.org/accounting/psi.html ; Facebook PSI microsite: https://facebookmicrosites.github.io/psi/docs/overview ; Kubernetes PSI metrics: https://kubernetes.io/docs/reference/instrumentation/understand-psi-metrics/
