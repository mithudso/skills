<!-- hub-reference-banner -->
> **Reference file — part of the `devops-observability` hub.** Formerly the standalone `ebpf-observability` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

name: ebpf-observability
title: eBPF for Linux Observability, Networking & Security
description: >
  eBPF — programmable in-kernel execution for Linux observability, networking, and security. Covers the
  shared foundation (the in-kernel VM and bytecode, the verifier, JIT, maps, helpers, program/attach types,
  CO-RE/BTF/libbpf portability) and the three application domains: observability/tracing (bpftrace, bcc,
  kprobe/uprobe/tracepoint/fentry, ring buffers, continuous profiling with Parca/Pixie), networking (XDP, tc,
  Cilium CNI, kube-proxy replacement, Hubble), and runtime security (BPF LSM, Tetragon, Falco, KubeArmor,
  seccomp). Includes development frameworks (libbpf, cilium/ebpf for Go, aya for Rust) and verifier
  troubleshooting.
  TRIGGER: writing or debugging an eBPF/bpftrace/bcc program; CO-RE/BTF/vmlinux.h portability; choosing
  XDP vs tc vs socket hooks; Cilium kube-proxy replacement or network policy; Hubble flow visibility; eBPF
  runtime security (BPF LSM, Tetragon, Falco); continuous profiling (Parca, Pixie); verifier rejections
  (instruction limit, stack 512, unbounded loop, invalid mem access).
  SKIP: general Linux host diagnostics without eBPF (ip/ss/tcpdump/systemd/perf) → linux-sysadmin;
  Kubernetes CNI/Service/ingress objects without eBPF specifics → kubernetes-networking; app-level
  OpenTelemetry instrumentation → nodejs-observability.
version: "1.0"
category: developer
updated: "2026-05-31"
tags:
  - ebpf
  - bpftrace
  - cilium
  - observability
  - linux
  - security
  - networking
  - tracing
---

# eBPF for Linux Observability, Networking & Security

## Overview

eBPF (extended Berkeley Packet Filter) lets you run sandboxed programs inside the Linux kernel — at
function entries, network hooks, security checkpoints — **without** changing kernel source or loading a
kernel module. It is the unifying technology behind a generation of cloud-native tooling: Cilium (the most
common CNI in production per the CNCF 2025 survey), Falco, Tetragon, Pixie, Parca, and the bcc/bpftrace
tracing tools. The pitch is "kernel superpowers without the risk": programmable visibility and control at
the kernel level, with the kernel statically guaranteeing each program is safe to run.

The three domains in this reference all sit on **one shared foundation** — the eBPF VM, verifier, maps, and
CO-RE portability. Read the foundation section first; the observability/networking/security sections inherit
it rather than re-derive it.

---

## Core Concepts (the shared foundation)

### 1. The in-kernel VM and bytecode

eBPF programs compile (via Clang/LLVM, target `bpf`) to a compact RISC-like bytecode with 11 64-bit
registers and a fixed instruction set. Programs are **event-driven**: they attach to a hook and run when that
hook fires (a syscall, a packet arriving, a function entry). Programs cannot block, cannot loop unbounded,
and have a **512-byte stack limit**.

### 2. The verifier (the safety contract)

Before any program runs, the kernel verifier performs static analysis over **every possible execution path**
to guarantee:
- **Memory safety** — no out-of-bounds reads/writes; every pointer is bounds-checked.
- **Termination** — no infinite loops; loops must have a provable upper bound (bounded loops allowed since
  kernel 5.3; before that, manual unrolling).
- **Type/argument safety** — helper-function arguments must match expected types; no uninitialized
  registers or stack.

The verifier's path-exploration budget is ~1 million instructions (the *analysis* limit), distinct from the
program's own instruction count. Unprivileged programs historically capped at `BPF_MAXINSNS` = 4096
instructions; privileged programs are bounded by the verifier's complexity budget instead.

### 3. JIT compilation

After verification, the JIT translates bytecode into native machine instructions for the host CPU so the
program runs at near-native speed. The kernel may retain the original bytecode, the verifier's analyzed
form, and the JIT image.

### 4. Maps (state + kernel↔user-space channel)

Maps are typed key/value data structures, created via the `bpf()` syscall, that persist state between
program invocations and bridge kernel programs to user space. Key types:
- **Hash / LRU hash** — general key/value; LRU auto-evicts.
- **Array / per-CPU array** — index-keyed; per-CPU variants avoid lock contention by giving each CPU its own
  slot (aggregate in user space).
- **Ring buffer** (`BPF_MAP_TYPE_RINGBUF`, 5.8+) — the modern way to stream events to user space; replaces
  the older per-CPU perf buffer with a single shared, ordered buffer and lower overhead.
- **Stack trace** — capture kernel/user stacks for profiling and flamegraphs.
- **Maps-of-maps, sockmap, devmap, cpumap** — used by networking programs for redirects and nested config.

### 5. Helpers and kfuncs

Programs cannot call arbitrary kernel functions. They call a curated set of **BPF helpers**
(`bpf_map_lookup_elem`, `bpf_probe_read_kernel`, `bpf_ktime_get_ns`, `bpf_perf_event_output`, etc.). Newer
kernels also expose **kfuncs** — kernel functions explicitly allow-listed for BPF, a more flexible successor
to fixed helpers.

### 6. Program & attach types

The program type determines what the program can do and where it attaches. Major families:
- **Tracing/observability:** `kprobe`/`kretprobe` (dynamic kernel function entry/return — *unstable* API),
  `tracepoint` and `raw_tracepoint` (static, *stable* API — prefer these), `fentry`/`fexit` (BTF-based,
  lower overhead than kprobe, 5.5+), `uprobe`/`uretprobe` (user-space functions), `perf_event` (sampling/
  profiling).
- **Networking:** `XDP` (driver-level, pre-`sk_buff`), `tc`/`sched_cls` (after `sk_buff`, ingress+egress),
  `cgroup/sock` and socket hooks (connect/sendmsg load balancing), `sk_msg`, `sockops`.
- **Security:** `BPF_PROG_TYPE_LSM` (attach to LSM hooks; can **deny** an operation by returning non-zero),
  `seccomp`/cgroup device.

### 7. CO-RE, BTF & libbpf (portability — the second half of the foundation)

The historic pain of eBPF was portability: programs read kernel structs whose **memory layout changes
between kernel versions** (field offsets shift, fields get renamed/removed), so a program built against one
kernel broke on another. The old fix (BCC) shipped Clang on every host to recompile at runtime — heavy and
fragile.

**CO-RE (Compile Once – Run Everywhere)** fixes this with four cooperating pieces:
- **BTF (BPF Type Format)** — compact type/debug metadata describing kernel and program types. The running
  kernel exposes its own types at `/sys/kernel/btf/vmlinux`.
- **Clang** records *relocations* — "this access is field X of struct Y" — instead of hard-coded offsets.
- **`vmlinux.h`** — a single generated header (`bpftrace -e` / `bpftool btf dump file /sys/kernel/btf/vmlinux
  format c > vmlinux.h`) that replaces the need for matching kernel headers.
- **libbpf** — at load time, reads the program's BTF + relocations and the *target* kernel's BTF, then
  rewrites field offsets, handles field existence/size differences, and patches Kconfig-dependent values so
  the same compiled object runs across kernels.

Result: a single small `.o` (or a Go/Rust binary embedding it) runs across kernel versions without Clang on
the target. This is the basis of modern BTF/CO-RE perf tools (Brendan Gregg, Andrii Nakryiko).

---

## Tools / Frameworks

### Observability & tracing
- **bpftrace** — high-level tracing language (awk-like). Pattern: `probe /filter/ { action }`. Best for
  one-liners and short scripts; compiles to BPF under the hood. Probe types: `kprobe`/`kretprobe`,
  `tracepoint`, `uprobe`/`uretprobe`, `profile`/`interval`, `software`/`hardware`, `BEGIN`/`END`.
- **bcc (BPF Compiler Collection)** — Python/C++ framework; ships the canonical `tools/` (execsnoop,
  opensnoop, biolatency, tcpconnect, runqlat, profile). Still widely used; heavier than CO-RE because it
  embeds Clang.
- **libbpf + libbpf-bootstrap** — the C reference path for production CO-RE tools; skeletons auto-generate
  load/attach code and auto-attach by annotation.
- **cilium/ebpf** (Go, pure-Go, no cgo) — load/compile/attach from Go; `bpf2go` embeds the compiled object
  into the Go binary. Powers Cilium, Falco-adjacent tooling.
- **aya** (Rust, no libbpf dependency) — write both kernel and user side in Rust.
- **eunomia-bpf** — package/distribute eBPF tools as JSON or WASM via libbpf.

### Networking
- **Cilium** — CNCF *graduated* eBPF CNI; most-used CNI in production (CNCF 2025). Attaches at XDP (fastest),
  tc ingress/egress (policy), and cgroup/socket hooks (load balancing).
- **Hubble** — Cilium's flow-visibility layer: network/L7 observability, service maps, flow logs.

### Security
- **Tetragon** — Cilium's Kubernetes-aware security observability + runtime enforcement; uses kprobes,
  tracepoints, and LSM hooks; can synchronously kill/deny. K8s-identity-aware (pod/namespace/labels).
- **Falco** — CNCF runtime-threat detection; largest rule library; eBPF syscall instrumentation; alerts and
  can kill from user space.
- **KubeArmor** — CNCF Sandbox; LSM-based enforcement (BPF-LSM, AppArmor, SELinux) for process/file/network
  policy.

### Profiling
- **Parca** — continuous profiling DaemonSet; samples stacks ~19×/s/CPU via eBPF → flamegraphs.
- **Pixie** — auto-instrumented K8s observability; captures HTTP/gRPC/DNS from kernel buffers (no app
  changes); auto distributed tracing.

---

## Methodology

### Picking the right tracing hook
1. **Prefer stable over dynamic.** Use `tracepoint`/`raw_tracepoint` (stable ABI) before `kprobe`
   (function names/args can change between kernels). Use `fentry`/`fexit` over kprobe when available — lower
   overhead, BTF-typed args.
2. **kprobe/uprobe for what tracepoints don't cover.** Dynamic instrumentation of any kernel/user function;
   accept the version-fragility tradeoff.
3. **perf_event/profile for sampling.** CPU/off-CPU flamegraphs, continuous profiling.
4. **Stream events via ring buffer**, not per-CPU perf buffer, on 5.8+.

### Picking the right network hook
1. **XDP** for the highest-throughput, earliest path (DDoS drop, L4 load balancing) — runs in the driver
   before `sk_buff` allocation. Modes: native (driver), offloaded (NIC), generic (slow, fallback).
2. **tc (sched_cls)** for full L3/L4 policy with `sk_buff` context, both ingress and egress.
3. **Socket / cgroup hooks** for connection-time (connect/sendmsg) service load balancing — this is how
   Cilium replaces kube-proxy without per-packet NAT.

> For the non-eBPF side of these hooks — XDP modes/actions and AF_XDP, the `tc` qdisc/class/filter model
> (HTB/fq_codel/netem), nftables, and netns/veth/bridge topologies as operational dataplane tools — load
> `references/linux-networking-stack.md`.

### Cilium kube-proxy replacement
Cilium implements the K8s Service abstraction in eBPF hash maps (O(1) lookup regardless of Service/endpoint
count), replacing iptables. East-west: rewrites the destination at the socket (`connect()`) level, avoiding
per-packet DNAT. North-south: XDP-accelerated L4 LB with DSR (Direct Server Return) and Maglev consistent
hashing. Set `kubeProxyReplacement: true`.

### Runtime security flow
Detection: attach to syscalls/kprobes/tracepoints, emit events (Falco/Tetragon) → alert. Enforcement: attach
to **BPF LSM** hooks (file open, process exec, socket) and return a non-zero value to **deny** the operation
in-kernel, synchronously — or have Tetragon/Falco send a kill signal. BPF LSM requires kernel ≥5.7 with
`CONFIG_BPF_LSM=y` and `lsm=...,bpf` on the boot cmdline.

### CO-RE build workflow
`clang -target bpf -g -O2 -c prog.bpf.c` → `bpftool gen skeleton` (or `bpf2go`) → load with libbpf/cilium-ebpf,
which relocates against the target's `/sys/kernel/btf/vmlinux`. Generate `vmlinux.h` once with `bpftool btf
dump`. No kernel headers or on-host Clang required at runtime.

---

## Practical Patterns

- **bpftrace one-liners** (the fastest way to answer a kernel question):
  - Count syscalls by process: `bpftrace -e 'tracepoint:raw_syscalls:sys_enter { @[comm] = count(); }'`
  - Trace `open()` with filename: `bpftrace -e 'tracepoint:syscalls:sys_enter_openat { printf("%s %s\n", comm, str(args.filename)); }'`
  - read() latency histogram: `bpftrace -e 'kprobe:vfs_read { @start[tid]=nsecs; } kretprobe:vfs_read /@start[tid]/ { @ns=hist(nsecs-@start[tid]); delete(@start[tid]); }'`
- **Use bcc canned tools first** (`execsnoop`, `opensnoop`, `biolatency`, `tcpconnect`, `runqlat`) before
  writing custom programs — they solve most production questions.
- **Per-CPU maps + user-space aggregation** for hot counters to avoid lock contention.
- **Ring buffer for events**, hash/array maps for state; size ring buffers as powers of two.
- **Cilium + Hubble** for "who talks to whom" without sidecars; Tetragon for exec/file/network enforcement.
- **Parca as a DaemonSet** for always-on flamegraphs (Meta cut fleet CPU ~20% with eBPF profiling; Datadog
  ~35% by switching agents).

## Anti-Patterns

- **Relying on kprobe for stable tooling** — kernel function signatures change; prefer tracepoints/fentry,
  or accept and version-gate the fragility.
- **BCC-on-every-host in production** — shipping Clang/LLVM to each node is heavy and breaks on header
  mismatch; use CO-RE/libbpf instead.
- **Large monolithic programs** — they hit verifier complexity/instruction limits. Design each program to do
  one task; split and chain via tail calls or BPF-to-BPF function calls.
- **Unbounded loops / uninitialized stack** — guaranteed verifier rejection.
- **Per-packet iptables at scale** — the whole point of Cilium's eBPF datapath is to escape iptables
  O(n) rule chains.
- **Treating XDP as a general policy hook** — XDP runs before `sk_buff`, so it lacks the context for full
  L3/L4 policy; use tc for that.

## Troubleshooting (verifier & load failures)

- **`R0 invalid mem access` / `invalid stack`** — out-of-bounds pointer or stack access. Bounds-check before
  dereferencing; ensure pointers from `bpf_probe_read*` are validated.
- **`back-edge` / "infinite loop detected" / "unbounded loop"** — give the loop a compile-time-bounded
  counter or use `#pragma unroll`.
- **`unreachable insn`** — dead code; the verifier rejects any unreachable instruction.
- **"too many instructions" / complexity limit** — program (or analyzed path count) exceeds the budget;
  simplify, split into tail-called programs, or reduce branching.
- **Exceeds 512-byte stack** — move large buffers into a per-CPU array map instead of the stack.
- **Use of uninitialized register/stack** — initialize all variables/struct fields before use; the verifier
  tracks definedness.
- **CO-RE relocation failure / field not found** — target kernel lacks BTF (`/sys/kernel/btf/vmlinux`
  missing) or the field genuinely doesn't exist on that kernel; use `bpf_core_field_exists()` guards and the
  `___suffix` ambiguity-resolution flavors.
- **`bpf_probe_read` failures** — on 5.5+ use the address-space-explicit `bpf_probe_read_kernel` /
  `bpf_probe_read_user` instead of the legacy `bpf_probe_read`.
- **Permission denied loading** — needs `CAP_BPF` (+`CAP_PERFMON`/`CAP_NET_ADMIN` depending on type), or
  root; check `kernel.unprivileged_bpf_disabled` sysctl. Inspect with `bpftool prog`/`bpftool map`; raise
  `ulimit -l` (locked memory) on older kernels.

## References

- eBPF.io — core infrastructure & applications landscape: https://ebpf.io/infrastructure/ , https://ebpf.io/applications/
- eBPF Docs — CO-RE concept: https://docs.ebpf.io/concepts/core/
- Andrii Nakryiko — "BPF CO-RE (Compile Once – Run Everywhere)": https://nakryiko.com/posts/bpf-portability-and-co-re/
- Brendan Gregg — "BPF binaries: BTF, CO-RE, and the future of BPF perf tools": https://www.brendangregg.com/blog/2020-11-04/bpf-co-re-btf-libbpf.html ; Linux eBPF tracing tools: https://www.brendangregg.com/ebpf.html
- bpftrace — One-Liner Tutorial: https://bpftrace.org/tutorial-one-liners ; intro: https://www.brendangregg.com/blog/2019-08-19/bpftrace.html
- Linux kernel docs — eBPF verifier: https://docs.kernel.org/bpf/verifier.html ; BPF Design Q&A: https://docs.kernel.org/bpf/bpf_design_QA.html
- Cilium — kube-proxy-free / docs: https://docs.cilium.io/en/stable/network/kubernetes/kubeproxy-free/ ; repo: https://github.com/cilium/cilium
- Tetragon — runtime enforcement FAQ: https://tetragon.io/docs/installation/faq/
- AccuKnox — BPF-LSM runtime security: https://accuknox.com/blog/runtime-security-ebpf-bpf-lsm
- eunomia — eBPF Ecosystem Progress 2024–2025: https://eunomia.dev/blog/2025/02/12/ebpf-ecosystem-progress-in-20242025-a-technical-deep-dive/
- Parca/Pixie continuous profiling: https://fatihkoc.net/posts/ebpf-parca-observability/ ; Metoro top eBPF tools: https://metoro.io/blog/top-ebpf-observability-tools
- Groundcover — eBPF verifier errors & debugging: https://www.groundcover.com/ebpf/ebpf-verifier
- O'Reilly — *Learning eBPF* (Rice), ch.5 CO-RE/BTF/libbpf: https://www.oreilly.com/library/view/learning-ebpf/9781098135119/ch05.html
