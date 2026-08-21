<!-- hub-reference-banner -->
> **Reference file — part of the `devops-linux-admin` hub.** Formerly the standalone `linux-networking-stack` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

name: linux-networking-stack
title: Linux Networking Stack — nftables, Network Namespaces, tc, XDP
description: >
  The Linux kernel networking dataplane and how to engineer it. Covers nftables (the modern netfilter
  framework — tables/chains/base-chain hooks/priorities/policies, rules, named sets, and verdict maps that
  replace iptables), network namespaces (netns + veth pairs + bridges to build isolated network stacks, the
  basis of container networking), tc traffic control (qdiscs classless and classful — fq_codel, HTB, TBF,
  netem — plus classes and filters for shaping/policing/QoS), and XDP (the eBPF-based express data path —
  XDP_DROP/PASS/TX/REDIRECT actions, generic/native/offload modes, AF_XDP zero-copy sockets). Use to build a
  firewall, isolate network stacks, shape/emulate traffic, or process packets at line rate.
---

# Linux Networking Stack — nftables, Network Namespaces, tc, XDP

## Overview

This reference covers the four layers you actually reach for to *engineer* (not just diagnose) Linux
networking: the firewall/NAT framework (**nftables**), the isolation primitive plus its plumbing
(**network namespaces** + veth + bridge), the queueing/shaping engine (**tc** traffic control), and the
high-performance early hook (**XDP**). It sits beside `references/linux-sysadmin.md` (the `ip`/`ss`/`tcpdump`
diagnostic cheat-sheet) and `references/ebpf-observability.md` (XDP/tc *as places to attach eBPF programs* —
Cilium, AF_XDP). The packet's journey ties them together: NIC → **XDP** (earliest, pre-skb) → tc ingress →
netfilter **nftables** (prerouting/input/forward/output/postrouting) → routing → tc egress → wire, with every
stage living inside whatever **network namespace** owns the device.

## Core Concepts

### 1. Network namespaces, veth, and bridges

A **network namespace** (CLONE_NEWNET) is an independent copy of the entire network stack: its own
interfaces, IP addresses, routing tables, neighbor tables, ports, and netfilter/nftables ruleset. This is the
foundation of container networking (Docker/Podman/Kubernetes each give a container its own netns). See
`references/linux-cgroups-namespaces.md` for the namespace primitive itself.

- **`ip netns add ns1`** / `ip netns list` / `ip netns exec ns1 <cmd>` (run a command inside a namespace).
- A **veth pair** is a virtual cable: two interfaces, packets in one end emerge from the other.
  `ip link add veth0 type veth peer name veth1`, then move one end into a namespace:
  `ip link set veth1 netns ns1`.
- A **bridge** (`ip link add br0 type bridge`) is a software L2 switch in the root namespace; connect each
  namespace's veth host-end to it (`ip link set veth0 master br0`) to let many namespaces talk — exactly how
  container runtimes build a default network. Add NAT to the outside world with an nftables masquerade rule.

Minimal two-namespace connectivity: create veth pair, push each end into a netns, assign IPs in the same
subnet (`ip -n ns1 addr add 10.0.0.1/24 dev veth1`), bring up (`ip -n ns1 link set veth1 up` + `lo`), ping.

### 2. nftables — tables, chains, rules, sets, maps

nftables is the modern netfilter framework that replaces iptables/ip6tables/arptables/ebtables with one tool
(`nft`) and one ruleset language. The hierarchy is **table → chain → rule**.

- **Tables** are the top-level container, scoped to a **family**: `ip`, `ip6`, `inet` (dual v4+v6 — preferred),
  `arp`, `bridge`, `netdev`. `nft add table inet filter`.
- **Chains** come in two kinds:
  - **Base chains** hook into the packet path and declare type/hook/priority/policy:
    `nft 'add chain inet filter input { type filter hook input priority 0; policy drop; }'`. Chain **types**:
    `filter`, `nat`, `route`. **Hooks**: `prerouting`, `input`, `forward`, `output`, `postrouting`.
    **Priority** (lower runs first; nat-dst at -100, filter at 0). **Policy**: `accept` (default) or `drop` —
    the verdict when no rule matches (a default-drop input chain is the secure baseline).
  - **Regular chains** have no hook; you `jump`/`goto` to them for organization.
- **Rules** = matches + statements: `nft add rule inet filter input tcp dport 22 accept`.
- **Verdict statements**: `accept`, `drop`, `reject` (sends ICMP), `continue`, `return` (leave chain),
  `jump <chain>` (call, returns), `goto <chain>` (tail-call, does not return), `queue` (to userspace).
- **Named sets** hold many values of one type (addresses, ports, ifnames) for O(1) matching and live updates:
  `nft add set inet filter blocklist { type ipv4_addr\; }` then `tcp ... ip saddr @blocklist drop`. Sets can be
  `flags interval` (ranges/CIDRs), `timeout` (auto-expiring entries — great for fail2ban-style dynamic blocks).
- **Maps & verdict maps (vmaps)** map a key to a value/verdict, collapsing many rules into one lookup:
  `tcp dport vmap { 22 : accept, 80 : accept, 443 : accept, 23 : drop }`. Named maps:
  `ip daddr map @nat_targets`.

Idiomatic usage is a single declarative ruleset file loaded with `nft -f /etc/nftables.conf` (atomic
replace), not incremental `add` commands. NAT example (masquerade for a bridge subnet):
`nft 'add chain inet nat postrouting { type nat hook postrouting priority 100; }'` then
`nft add rule inet nat postrouting ip saddr 10.0.0.0/24 oifname "eth0" masquerade`.

### 3. tc — traffic control (qdiscs, classes, filters)

`tc` (iproute2) manages the kernel's packet queueing on egress (and, via ingress qdisc, some ingress). The
model: a **qdisc** (queueing discipline) attached to a device, optionally **classes** forming a tree, and
**filters** that classify packets into classes.

- **Classless qdiscs** (just shape/schedule, no sub-classes):
  - **fq_codel** — fair-queueing + CoDel AQM; the modern sane default, fights bufferbloat.
  - **TBF** (Token Bucket Filter) — simple rate limit to a single ceiling.
  - **netem** — *network emulator*: inject delay, jitter, loss, duplication, reordering, corruption — the go-to
    for testing under bad-network conditions. `tc qdisc add dev eth0 root netem delay 100ms 20ms loss 1%`.
- **Classful qdiscs** (build a class tree for link-sharing/QoS):
  - **HTB** (Hierarchical Token Bucket) — the workhorse for guaranteed-rate + ceiling link sharing.
    `tc qdisc add dev eth0 root handle 1: htb default 30`, then
    `tc class add dev eth0 parent 1: classid 1:10 htb rate 100mbit ceil 200mbit`. Each leaf class can carry its
    own classless qdisc (commonly fq_codel) for intra-class fairness.
- **Filters** classify into classes by matching (`u32` on header fields, `flower` on flow keys, `bpf`):
  `tc filter add dev eth0 parent 1: protocol ip prio 1 u32 match ip dport 443 0xffff flowid 1:10`.
- Handles use `major:minor` notation; `root` is egress, `ingress`/`clsact` qdiscs handle ingress and host
  eBPF/tc attach points (the latter is where `references/ebpf-observability.md` hooks tc-BPF / Cilium).

Shaping vs policing: shaping (HTB/TBF) *buffers and delays* to a rate; policing (`tc ... action police`)
*drops* over-rate packets without buffering.

### 4. XDP — eXpress Data Path

XDP runs an eBPF program at the **earliest** point in the RX path — in the driver, before an `sk_buff` is even
allocated — so it can drop/redirect/transmit at line rate (millions of pps) while bypassing most of the stack.
It is the packet-processing counterpart to tc-BPF (which runs later, with an skb). Program logic and the eBPF
toolchain live in `references/ebpf-observability.md`; this section is the operational/decision view.

- **Return actions**: `XDP_DROP` (drop here — DDoS scrubbing, the canonical XDP use), `XDP_PASS` (continue to
  the normal stack), `XDP_TX` (bounce back out the same NIC, possibly modified — e.g. simple load balancer),
  `XDP_REDIRECT` (send to another NIC, a CPU, or an AF_XDP socket), `XDP_ABORTED` (error/drop + tracepoint).
- **Modes**: **native** (driver implements XDP hook in its RX path — full speed, needs driver support),
  **offloaded** (program runs on the SmartNIC, zero host CPU — needs NIC support, e.g. Netronome),
  **generic/SKB** (kernel emulates the hook in `netif_receive_skb()` after skb alloc — works on any driver but
  loses most of the speed benefit; fine for dev/testing).
- **AF_XDP** is a socket family (`XDP_REDIRECT` target) giving userspace zero-copy, kernel-bypass access to raw
  frames via UMEM rings — the basis for high-speed userspace networking (DPDK-like) without leaving the kernel
  driver model.
- Load/manage with `ip link set dev eth0 xdp obj prog.o sec xdp` (or `xdpgeneric`/`xdpoffload`), or via
  `bpftool`/libxdp. Picking the hook: XDP for earliest drop/redirect at line rate; **tc (clsact) BPF** when you
  need the skb, egress processing, or full packet metadata.

## Tools / Frameworks

- **`nft`** (nftables) — the single firewall/NAT tool; `/etc/nftables.conf` declarative ruleset.
- **`ip` / iproute2** — `ip netns`, `ip link ... type veth|bridge`, `ip -n <ns>`; the namespace/veth/bridge tooling.
- **`tc`** (iproute2) — qdisc/class/filter management; `tc -s qdisc show` for stats/drops.
- **XDP/libxdp/bpftool/`ip link ... xdp`** — load and manage XDP programs (program authoring →
  `references/ebpf-observability.md`).
- Diagnostics: `ss`, `tcpdump`, `nft monitor`, `bridge` (FDB) — see `references/linux-sysadmin.md`.

## Methodology

1. **Isolate first with netns** when building or reproducing a topology — veth pairs into namespaces, a bridge
   in root, NAT out via nftables masquerade. This is the container-networking pattern in miniature.
2. **Firewall with nftables, default-drop**: `inet` family, base `input` chain `policy drop`, allow loopback +
   established/related + the few services you expose; use named sets (with `timeout`) for dynamic block/allow lists.
3. **Shape/QoS with tc**: HTB tree for guaranteed-rate link sharing, fq_codel leaves for fairness; netem to
   *test* resilience to latency/loss before production.
4. **Reach for XDP only for line-rate drop/redirect** (DDoS, L4 LB); otherwise tc-BPF or nftables is simpler.
   Choose the hook by where in the path you need to act (XDP earliest pre-skb → tc → netfilter).

## Practical Patterns

- Two-namespace lab: `ip netns add a; ip netns add b; ip link add va type veth peer name vb;
  ip link set va netns a; ip link set vb netns b; ip -n a addr add 10.0.0.1/24 dev va; ...; ping`.
- Simulate a flaky WAN for testing: `tc qdisc add dev eth0 root netem delay 150ms 30ms distribution normal loss 2% reorder 5%`.
- Rate-limit one tenant: HTB class at `rate 50mbit ceil 100mbit` with a `flower`/`u32` filter on its IP.
- Drop a spoofed source range at the NIC: XDP program returning `XDP_DROP` on a match (DDoS scrubbing) — far
  cheaper than an nftables drop because it never builds an skb.
- Atomic firewall reload: edit `/etc/nftables.conf`, `nft -c -f` to syntax-check, then `nft -f` to swap in.

## Anti-Patterns

- **Mixing iptables and nftables on the same box** without understanding `iptables-nft` translation → confusing
  duplicate/ordered rules. Pick one; on modern distros prefer native nftables `inet` rulesets.
- **`policy accept` base input chain** "to be safe" — that is a wide-open firewall. Default-drop and allowlist.
- **Incremental `nft add rule` in scripts** instead of a declarative `-f` ruleset → non-atomic, drift-prone.
- **Leaving netns/veth/bridges orphaned** after a failed lab (`ip netns del`, delete the bridge) — they leak and
  confuse later runs.
- **HTB without a `default` class** → unclassified traffic falls into class 0 (direct, unshaped), silently
  bypassing your QoS.
- **Reaching for XDP when tc/nftables suffices.** XDP has no skb, limited helpers, and driver/mode constraints;
  use it only when line-rate early drop/redirect is the actual requirement (see `references/ebpf-observability.md`).
- **Generic/SKB XDP in production expecting native speed** — it runs after skb alloc and gives little benefit.

## Troubleshooting

- nftables rule "not matching" → check chain **priority/hook** ordering (another base chain at lower priority
  may have already accepted/dropped), and family (`ip` vs `ip6` vs `inet`). `nft monitor` shows live trace;
  `nft --debug=netlink` aids syntax issues. Add `meta nftrace set 1` + `nft monitor trace` to trace a packet.
- netns "no connectivity" → verify both veth ends are `up`, `lo` is up inside the ns, IPs share a subnet, and a
  route/default gateway exists; bridge needs `ip link set br0 up` and member ports enslaved + up.
- tc shaping "not working" → confirm the qdisc is on the **egress** path (shaping is egress; for ingress shaping
  redirect to an IFB device), check `tc -s qdisc show` for drops/overlimits, and that filters actually match
  (flowid points at a real class).
- XDP "won't load" → driver lacks **native** support (fall back to `xdpgeneric` for testing), verifier rejected
  the program, or another XDP prog is already attached (`ip link set dev eth0 xdp off` first). See verifier
  troubleshooting in `references/ebpf-observability.md`.

## References

- nftables wiki — Quick reference (10 min): https://wiki.nftables.org/wiki-nftables/index.php/Quick_reference-nftables_in_10_minutes ; Sets: https://wiki.nftables.org/wiki-nftables/index.php/Sets ; Verdict Maps: https://wiki.nftables.org/wiki-nftables/index.php/Verdict_Maps_(vmaps)
- ArchWiki nftables: https://wiki.archlinux.org/title/Nftables
- Linux Network Namespaces (Scott Lowe): https://blog.scottlowe.org/2013/09/04/introducing-linux-network-namespaces/ ; netns + bridge + veth lab: https://medium.com/@masud.educations/setting-up-linux-network-namespace-and-bridge-for-network-isolation-9a9bba6e75de
- RHEL 9 Linux traffic control: https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/9/html/configuring_and_managing_networking/linux-traffic-control_configuring-and-managing-networking ; ArchWiki Advanced traffic control: https://wiki.archlinux.org/title/Advanced_traffic_control ; tc(8): https://man7.org/linux/man-pages/man8/tc.8.html
- XDP (Wikipedia overview): https://en.wikipedia.org/wiki/Express_Data_Path ; Datadog "gentle introduction to XDP": https://www.datadoghq.com/blog/xdp-intro/ ; Red Hat "Get started with XDP": https://developers.redhat.com/blog/2021/04/01/get-started-with-xdp
