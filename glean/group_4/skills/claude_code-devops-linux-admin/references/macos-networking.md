<!-- hub-reference-banner -->
> **Reference file — part of the `devops-linux-admin` hub.** Created by `/dr` deep research 2026-06-12 (not a folded standalone skill).
> Sibling topics in this family are reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Linux-side host diagnostics live in `references/linux-sysadmin.md` and `references/linux-networking-stack.md`; this file is the macOS counterpart.

---
name: macos-networking
title: macOS Networking Stack & Diagnostics — SystemConfiguration/scutil, mDNSResponder, PF, Network Extension, Wi-Fi tooling
description: >-
  macOS networking stack & diagnostics — SystemConfiguration/configd + scutil (--dns/--nwi/--proxy),
  networksetup/locations/ipconfig, mDNSResponder DNS path (/etc/resolver, scoped resolvers, cache flushing,
  why dig bypasses it), PF/pfctl & Application Firewall (socketfilterfw), Network Extension VPNs (utun,
  split DNS, primary-service rank), Wi-Fi tooling (wdutil, sysdiagnose), macOS↔Linux CLI map, failure modes
  (captive portal, Private Relay, AWDL). TRIGGER: diagnosing DNS/routing/VPN/firewall/Wi-Fi on macOS;
  "dig works but the browser fails" on a Mac; VPN split-DNS or DNS-not-restored; finding the macOS
  equivalent of a Linux networking command. SKIP: Linux host networking (ip/ss/nftables) →
  references/linux-sysadmin.md / references/linux-networking-stack.md; DNS protocol internals (DNSSEC,
  DoH/DoT, delegation) → networking (references/dns-deep-dive.md); generic packet-capture methodology;
  macOS privilege/TCC internals → devops-linux-internals (references/linux-mac-privilege.md).
version: "1.0.0"
category: developer
updated: "2026-06-12"
whenToUse: >
  Use when diagnosing or scripting macOS network behavior: reading scutil --dns/--nwi output, fixing VPN
  split-DNS, flushing the resolver cache, configuring PF or the application firewall, understanding why a
  VPN's utun became the primary interface, replacing the airport CLI, mapping Linux commands to Darwin, or
  working a captive-portal / Private Relay / AWDL failure mode.
keywords:
  - macos networking
  - scutil
  - SystemConfiguration
  - configd
  - mDNSResponder
  - networksetup
  - ipconfig getpacket
  - /etc/resolver
  - dscacheutil
  - pfctl
  - socketfilterfw
  - Network Extension
  - NEPacketTunnelProvider
  - utun
  - split DNS
  - wdutil
  - sysdiagnose
  - captive portal
  - iCloud Private Relay
  - AWDL
tags:
  - macos
  - networking
  - dns
  - vpn
  - firewall
  - wifi
  - diagnostics
  - sysadmin
---

# macOS Networking Stack & Diagnostics

`verified-as-of: 2026-06-12` (volatile claims in this file — version timelines, Sequoia/Tahoe behavior changes, Private Relay mechanics — were verified against fetched sources on this date.)

## Contents

- [Overview](#overview)
- [Core Concepts](#core-concepts)
  - [1. The configuration model: configd, the dynamic store, IPMonitor](#1-the-configuration-model-configd-the-dynamic-store-ipmonitor)
  - [2. scutil — reading and watching live state](#2-scutil--reading-and-watching-live-state)
  - [3. networksetup, locations, and ipconfig](#3-networksetup-locations-and-ipconfig)
  - [4. The DNS resolution path and mDNSResponder](#4-the-dns-resolution-path-and-mdnsresponder)
  - [5. .local, mDNS/Bonjour, and dns-sd](#5-local-mdnsbonjour-and-dns-sd)
  - [6. Resolver configuration and cache flushing](#6-resolver-configuration-and-cache-flushing)
  - [7. PF packet filter and the Application Firewall](#7-pf-packet-filter-and-the-application-firewall)
  - [8. Network Extension: VPNs, filters, proxies, and interface reordering](#8-network-extension-vpns-filters-proxies-and-interface-reordering)
  - [9. Wi-Fi diagnostics: wdutil, sysdiagnose, Wireless Diagnostics](#9-wi-fi-diagnostics-wdutil-sysdiagnose-wireless-diagnostics)
  - [10. CLI mapping: macOS vs Linux](#10-cli-mapping-macos-vs-linux)
- [Methodology: a macOS network diagnosis workflow](#methodology-a-macos-network-diagnosis-workflow)
- [Practical Patterns](#practical-patterns)
- [Anti-Patterns](#anti-patterns)
- [Troubleshooting: macOS-specific failure modes](#troubleshooting-macos-specific-failure-modes)
- [References](#references)

## Overview

macOS networking is *not* "BSD with a GUI." Configuration flows through a dedicated daemon (`configd`) and an
in-memory key-value store; DNS flows through a system resolver daemon (`mDNSResponder`) that most Unix DNS
tools bypass; VPNs are Network Extension processes that win primary-interface election rather than scripts
that edit resolv.conf; and the classic BSD tools (`ifconfig`, `netstat`, `route`) coexist with
Apple-specific ones (`scutil`, `networksetup`, `ipconfig`, `wdutil`, `dns-sd`, `networkQuality`). The single
most common diagnostic mistake on macOS is testing DNS with `dig` and believing the result describes what
apps see; it doesn't (see §4). The second most common is assuming Linux mental models (one resolv.conf, one
routing truth, `ss`/`ip` available), which don't hold.

## Core Concepts

### 1. The configuration model: configd, the dynamic store, IPMonitor

**configd** maintains the desired and current network state and notifies clients of changes. It loads
configuration-agent bundles from `/System/Library/SystemConfiguration/`: notably `IPConfiguration` (DHCP/
static/link-local address management), `IPMonitor` (primary-service election), `PreferencesMonitor`,
`KernelEventMonitor`, `InterfaceNamer`, `PPPController`.[^1][^2]

Two stores matter:

- **Dynamic store (`SCDynamicStore`)** — in-memory, owned by configd, recreated at every boot. Holds the
  active configuration and live state. `scutil` is a thin client over it.[^1][^2]
- **Persistent store** — `/Library/Preferences/SystemConfiguration/preferences.plist` (since 10.3). Top-level
  keys: `Sets` (one per network *location*), `CurrentSet`, `NetworkServices`, `System`. Apple marks the
  format private — change it via `networksetup`/System Settings, not by editing the plist.[^2]

**Key namespaces.** `Setup:/...` keys are desired configuration (mapped from preferences by
PreferencesMonitor); `State:/...` keys are live state published by agents (e.g., what DHCP actually
assigned).[^2][^3] Propagation: preferences edit → `Setup:` keys → per-protocol agents configure interfaces
and publish `State:/Network/Service/<ServiceID>/...` → **IPMonitor** elects the primary service.

**Primary service election.** IPMonitor picks the highest service in the current set's **`ServiceOrder`**
array that is currently usable, publishes `State:/Network/Global/IPv4` (`PrimaryService`,
`PrimaryInterface`, `Router`), installs that service's default route, and assembles the global DNS and proxy
configuration from the primary service, preferring `Setup:` (manually configured) values over `State:`
(DHCP-supplied) ones.[^2] The primary service therefore determines the default route, DNS, and proxies all at
once. VPN-type services take priority over non-VPN services and cannot be reordered in the UI.[^4] A
service can be demoted from election by setting `PrimaryRank Second` on its `State:` service key, a
practitioner technique (single-source; treat as tentative) used to stop a VPN's DNS from going global.[^5]

Caveat (multiple independent reports, no Apple acknowledgment): service order is occasionally **not**
honored when two services share a gateway: the default route stays on the first-*connected* interface until
you toggle the other one.[^6]

### 2. scutil — reading and watching live state

```bash
scutil --dns          # the resolver configuration apps actually use
scutil --nwi          # network-state info: interface priority order per protocol
scutil --proxy        # State:/Network/Global/Proxies
scutil --get HostName             # also: LocalHostName, ComputerName; --set to write (root)
scutil -r host.example.com        # reachability (SCNetworkReachability flags); -W to watch
sudo scutil           # interactive store browser
  > list State:/Network/.*        # pattern-list keys
  > show State:/Network/Global/IPv4
  > n.add State:/Network/Global/IPv4 ; n.watch    # watch for changes
```

**`scutil --dns` anatomy.**[^7][^8] Two sections:

- **"DNS configuration"** — resolver #1 is the default (the primary service's DNS). *Supplemental*
  resolvers follow: entries with a `domain` and flag `Supplemental` are **match-domain (split-DNS)**
  resolvers sourced from `/etc/resolver/*` files or a VPN's `SupplementalMatchDomains`; queries matching the
  domain go to that resolver instead of the default. Entries for `local`, `254.169.in-addr.arpa`, and the
  `*.ip6.arpa` zones with `options: mdns` and `order` ≈ 300000+ are the multicast-DNS zones.
- **"DNS configuration (for scoped queries)"** — per-interface resolvers (flag `Scoped`, with
  `if_index : 6 (en0)`), used when a query is explicitly bound to an interface, the normal arrangement for
  VPN tunnels.

Field semantics: `search domain[n]` = suffixes appended to unqualified names; `order` = resolver(5)
`search_order` (lower wins among clients for the same domain); `flags: Request A records, Request AAAA
records` = which record types are queried (AAAA is dropped when no routable IPv6 exists); `reach` = the
server's SCNetworkReachability flags. Remember: `scutil --dns` shows *configuration*, not guaranteed runtime
*behavior*: mDNSResponder has been observed preferring servers that support encrypted transports over the
listed order (practitioner reports).[^9]

**`scutil --nwi`** lists interfaces in priority order per protocol: first entry is effectively primary, so a
connected full-tunnel VPN shows its `utun` first; per-interface `flags` (e.g. `0x5 (IPv4,DNS)`) tell you
which interface supplies DNS. This is the standard way to find the physical interface underneath a VPN.[^10]

**Writes don't stick.** `State:` writes via interactive scutil are transient: agents republish their keys
and the whole store is rebuilt at boot.[^3] Persistent change belongs to `networksetup`/System Settings.

### 3. networksetup, locations, and ipconfig

**Three-tool mental model** (each operates on a different layer):

| Tool | Layer | Persistence |
| --- | --- | --- |
| `networksetup` | **Setup** (preferences.plist via SC) | survives reboot; agents apply it |
| `scutil` | **State** (dynamic store) | live truth; writes are transient |
| `ipconfig` | per-interface **IPConfiguration agent** | direct interrogation/control; `set` is temporary |

**networksetup essentials:**[^11]

```bash
networksetup -listallnetworkservices        # '*' prefix = disabled service
networksetup -listallhardwareports          # "Wi-Fi" ↔ en0 ↔ MAC mapping
networksetup -listnetworkserviceorder       # priority order; -ordernetworkservices to change
networksetup -getinfo "Wi-Fi"
networksetup -setdnsservers "Wi-Fi" 1.1.1.1 9.9.9.9    # literal 'empty' to clear → revert to DHCP DNS
networksetup -getdnsservers "Wi-Fi"         # shows ONLY manually set servers — DHCP DNS prints
                                            # "There aren't any DNS Servers set..."
networksetup -setsearchdomains "Wi-Fi" corp.example.com
networksetup -setwebproxy "Wi-Fi" proxy.corp.example 8080    # also -setsecurewebproxy, -setautoproxyurl (PAC)
networksetup -setairportnetwork en0 SSID password   # airport verbs take the DEVICE, not service name
networksetup -listlocations          # also -getcurrentlocation, -switchtolocation X, -createlocation X populate
```

Gotchas: since Big Sur, standard users can only read settings, toggle Wi-Fi power, and switch Wi-Fi
networks — mutations need admin (`** Error: Command requires admin privileges.`).[^12] networksetup is
**blind to most VPN/utun services** (they're not network services in the Setup model), so
`-setdnsservers` against a tunnel fails — VPN scripts use scutil writes instead.[^13]

**Locations are "Sets".** A location is one dictionary under `Sets` in preferences.plist holding its own
`ServiceOrder` and service references; switching (`networksetup -switchtolocation`, or the older
`scselect`) repoints `CurrentSet` and reconfigures immediately.[^2][^11] A clean "Automatic-DNS" location is
a useful tool for captive-portal debugging (§Troubleshooting).

**ipconfig (Darwin's, not Windows'):** talks directly to the IPConfiguration agent.[^14]

```bash
ipconfig getifaddr en0           # bare IPv4 — the scripting one-liner
ipconfig getpacket en0           # the accepted DHCP ACK: yiaddr, server_identifier, lease_time (hex
                                 # seconds), subnet_mask, router, domain_name_server. Empty if no DHCP.
ipconfig getsummary en0          # full agent state: link, SSID/BSSID (redacted on new macOS), DHCP data
ipconfig getoption en0 router    # single DHCP option
sudo ipconfig set en0 DHCP       # classic lease renew — temporary service, not persistent
```

`getpacket`/`getsummary` are the canonical way to see *DHCP-supplied* DNS (which `-getdnsservers` hides).
SSID/BSSID fields in `getsummary` became `<redacted>` in macOS 15.6 (see §9).[^15]

### 4. The DNS resolution path and mDNSResponder

**mDNSResponder is THE system resolver.** Since Mac OS X 10.6 it serves as the system-wide unicast DNS stub
resolver *and* the mDNS/Bonjour responder (UDP 5353), caching replies for client processes. It is
open-source (apple-oss-distributions).[^16] Apps resolve via `getaddrinfo()` → libinfo's mdns module →
`DNSServiceQueryRecord` IPC into mDNSResponder; `/etc/hosts` is honored on this path and re-read on
SIGHUP.[^17]

**The discoveryd saga** (the cautionary tale behind "restart your resolver"): OS X 10.10 Yosemite (Oct 2014)
replaced mDNSResponder with **discoveryd**, a closed-source rewrite. It brought duplicate-hostname renames
("MacBook (2)", "(3)", …, an interplay with Bonjour Sleep Proxy), Wi-Fi drops, CPU spins, and resolution
failures. Apple removed it in 10.10.4 beta 4 (May 26, 2015); 10.10.4 final (June 30, 2015) shipped with
mDNSResponder restored.[^18]

**Who bypasses the system resolver (the load-bearing diagnostic fact).** `dig`, `nslookup`, and `host`
embed their own resolver code: they read `/etc/resolv.conf` (or your `@server`) directly and ignore
`/etc/hosts`, `/etc/resolver/`, scoped resolvers, and match domains. `ping`, `curl`, browsers, and GUI apps
use the system path. Apple's own dig man page warns it "does not use the host name and address resolution or
the DNS query routing mechanisms used by other processes."[^19] `/etc/resolv.conf` itself is an
auto-generated **compatibility stub** carrying only the default unscoped resolvers: split-DNS configuration
*cannot be represented in it by design*.[^20]

Test the path apps actually use:

```bash
dscacheutil -q host -a name internal.corp.example   # system path
dns-sd -G v4v6 internal.corp.example                # full resolver path incl. search domains
ping -c1 internal.corp.example                      # quick proxy
```

**Observing mDNSResponder:**

```bash
sudo log stream --info --predicate 'process == "mDNSResponder"'
sudo log config --subsystem com.apple.mDNSResponder --mode "level:debug"   # then stream with --debug
```

The classic `sudo killall -INFO mDNSResponder` state dump is deprecated; the daemon points to `dns-sd -O`,
which then fails with "State dump is only enabled in internal builds". **There is currently no supported way
to view the DNS cache contents on a stock Mac**; you can only flush it.[^21] DNS names in unified logs are
redacted as `<private>` unless a profile enables private data.

### 5. .local, mDNS/Bonjour, and dns-sd

macOS special-cases `.local` to multicast DNS per RFC 6762 — `scutil --dns` always shows the
`domain: local, options: mdns` resolver. Corporate **unicast** `.local` zones therefore break name
resolution and AD binding; Apple's guidance is to use a registered domain.[^22] The documented escape hatch:
resolver(5) supports two clients for the same domain ordered by `search_order` — an `/etc/resolver/local`
file with a `nameserver` and a low `search_order` makes unicast `.local` resolve while mDNS keeps
working.[^23]

`dns-sd` is the Bonjour/system-path test tool: `dns-sd -B _ssh._tcp` (browse), `dns-sd -R` (register),
`dns-sd -q name rrtype` (query through the system path), `dns-sd -G v4v6 name` (address lookup).[^16]

Stale Bonjour records (device vanished without an mDNS "Goodbye" packet) self-purge ~15 s after a failed
resolve; interface bounce or sleep also clears cached records (Apple QA1310).[^24]

### 6. Resolver configuration and cache flushing

**/etc/resolver/ files.** One file per domain (filename = domain), keywords from resolver(5): `nameserver`
(up to 3), `port`, `search_order`, `timeout`, `options`. They affect **only the system resolver path**:
`dig` will never see them, which is the #1 false bug report. Verify with `scutil --dns` (the file appears as
a supplemental resolver) and test with `dscacheutil -q host`.[^23][^19] Failures are silent: the entry can
look healthy in `scutil --dns` yet not be consulted. A March-2026 report (single source — tentative, watch
item) says macOS 26 broke /etc/resolver for custom TLDs not in the IANA root zone (including `.test`):
mDNSResponder intercepts them as mDNS and serves a long-TTL negative answer while `scutil --dns` still shows
the resolver.[^25]

**Flushing the cache** — current canonical incantation:

```bash
sudo dscacheutil -flushcache; sudo killall -HUP mDNSResponder
```

The HUP does the real DNS work: the handler drops the in-memory cache and re-reads config including
`/etc/hosts`; `dscacheutil -flushcache` clears the separate Directory Services cache and on 10.7+
contributes nothing to DNS specifically (harmless, kept by convention).[^26] Version archaeology:
`lookupd -flushcache` (≤10.4), `dscacheutil -flushcache` alone (10.5–10.6), `killall -HUP mDNSResponder`
(10.7+), and `discoveryutil mdnsflushcache` only on 10.10.0–10.10.3.[^27] Killing mDNSResponder outright is
safe (launchd restarts it; momentary stall); a wedged daemon ignoring HUP can be bounced with
`sudo launchctl stop com.apple.mDNSResponder; sudo launchctl start com.apple.mDNSResponder` (practitioner
fallback).

"Flushed but still stale" usually means a different cache layer: upstream/router resolver, negative caching
(NXDOMAIN held per the zone's SOA minimum), an app-internal cache (Chrome/Electron keep their own), or a
DNS-proxying agent (WARP, Tailscale MagicDNS) in front of mDNSResponder.[^26]

### 7. PF packet filter and the Application Firewall

**PF** (ported from OpenBSD; syntax frozen around the OpenBSD 4.5 era) replaced `ipfw`, which was deprecated
when PF arrived in 10.7 Lion and removed in 10.10.[^28] PF is **disabled by default** and is *not* enabled at
boot: `com.apple.pfctl.plist` loads `/etc/pf.conf` without `-e`.[^29]

```bash
sudo pfctl -e                 # enable (use -d to disable; separate commands, never piped)
sudo pfctl -E                 # enable + take a reference (Apple extension)
sudo pfctl -X <token>         # release that reference
sudo pfctl -s References      # who has PF enabled (pid/name/token)
sudo pfctl -nf /etc/pf.conf   # validate; -f to load
sudo pfctl -sr | -ss          # rules / states
sudo pfctl -sr -a 'com.apple/*'   # rules inside Apple's anchors
```

The `-E`/`-X`/`-s References` reference-count contract is Apple-specific: each system component that needs PF
takes a reference; PF turns off only when the last is released.[^30] The stock `/etc/pf.conf` is a chain of
`com.apple` anchor statements in strict order (scrub → nat → rdr → dummynet → anchor → load anchor);
`/etc/pf.anchors/com.apple` defines sub-anchors like `200.AirDrop` and `250.ApplicationFirewall`.[^29]

Gotchas: `/etc/pf.conf` is **overwritten by macOS updates** (even point releases); put custom rules in your
own anchor file plus a LaunchDaemon running `pfctl -E -f <ruleset>`; manually loaded rules don't reliably
survive network changes (practitioners add `WatchPaths` on SystemConfiguration plists to re-load).[^31]

**Application Firewall (ALF)** — per-application **inbound** firewall at the socket layer, distinct from and
layered above PF; enabling ALF takes a PF reference (`pfctl -s References` shows it) and ALF's stealth/
block-all features materialize as rules in the `com.apple/250.ApplicationFirewall` anchor.[^29]

```bash
/usr/libexec/ApplicationFirewall/socketfilterfw --getglobalstate   # --setglobalstate on|off to change
# other verbs: --listapps, --add <path>, --blockapp/--unblockapp <path>,
#   --getallowsigned/--setallowsigned on|off, --getstealthmode/--setstealthmode on|off
```

**macOS 15 Sequoia change:** ALF settings are no longer stored in `/Library/Preferences/com.apple.alf.plist`;
anything parsing that plist broke (osquery's `alf` table went empty); `socketfilterfw` is the supported
surface and Apple DTS states the new storage "is not considered API". `--listapps` also stopped printing app
paths (names + state only).[^32]

### 8. Network Extension: VPNs, filters, proxies, and interface reordering

**Provider types:** `NEPacketTunnelProvider` (packet-tunnel VPN), `NEAppProxyProvider` /
`NETransparentProxyProvider` (flow-level proxying), `NEFilterDataProvider` + `NEFilterControlProvider`
(content filter — the data provider sees content but runs in an exfiltration-blocking sandbox; the control
provider has network access but never sees content), `NEFilterPacketProvider`, `NEDNSProxyProvider`,
`NEDNSSettingsManager` (system-wide encrypted-DNS), `NEVPNManager` (built-in IKEv2/IPsec personal VPN).[^33]
Legacy Network Kernel Extensions were deprecated in Catalina and refused by default in Big Sur; modern
clients ship **system extensions** (`systemextensionsctl list` to inspect).[^34]

**How a VPN captures traffic.** The provider creates a **`utun`** interface and hands the system
`NEPacketTunnelNetworkSettings` whose `includedRoutes`/`excludedRoutes` install routes: a default route in
`includedRoutes` = full tunnel; a subset = split tunnel. More-specific entries in the system routing table
still win, and apps that scope connections to an interface bypass the tunnel, unless `enforceRoutes` is set
(split-tunnel only; Apple DTS: it does **not** apply to default routes) or `includeAllNetworks` is set (full
capture minus DHCP, captive-portal probes, and some Apple services).[^35] `includeAllNetworks` has documented
sharp edges: Mullvad still refuses to use it (their post concerns the iOS app; the NE flag is shared across
Apple platforms) (breakage up to "no internet until reboot"), profiles with iAN
conflict with each other, and it blocks the provider's own out-of-tunnel bootstrap traffic.[^36]

**Interface reordering, concretely:** the connected tunnel's service outranks physical services in
IPMonitor's election, so `State:/Network/Global/IPv4 PrimaryInterface` becomes `utunN`, `scutil --nwi` lists
the utun first, and `netstat -rn` shows the default (or scoped) routes on it. Each reconnect creates a *new*
utunN; a stale utun lingers (with leftover routes) as long as the owning process keeps its control socket
open; route deletion won't remove the interface.[^37]

**Split DNS on the NE side.** `NEDNSSettings.matchDomains` lists the domains routed to the tunnel's
resolvers (surfacing in `scutil --dns` as `Supplemental` match-domain resolvers, order ≈ 103000, plus scoped
entries on the utun). An **empty string** in matchDomains makes the tunnel resolver the default. By default
every matchDomains entry is *also appended to the global search list*; set `matchDomainsNoSearch = true`
for pure routing domains. One resolver pool per NEDNSSettings: you cannot route different domains to
different servers from one settings object.[^38] Known regression class: IKEv2 mobileconfig profiles on
Ventura ≥13.3 degrading Supplemental resolvers to merely Scoped (split DNS silently broken).[^39]
**DNS-not-restored-after-disconnect** is a recurring vendor bug class: teardown treated as a no-op when the
VPN pushed the same DNS IP as the underlay; macOS pinning a restored resolver to the wrong `if_index` after
sleep; clients wiping interface DNS on disconnect. First-line fix is the flush incantation; durable fix is
client-side.[^40]

**The Big Sur exclusion-list saga** (worth knowing as precedent): macOS 11.0 shipped a
`ContentFilterExclusionList` exempting ~53 Apple apps from NE content filters, so Little Snitch/LuLu-class
firewalls couldn't see App Store/Maps/iCloud traffic, and malware could piggyback on excluded processes
(demonstrated by Patrick Wardle). Apple removed the list in 11.2.[^41]

**Transparent proxies:** `NETransparentProxyProvider` differs from its superclass: returning `false` from
`handleNewFlow` lets the flow go *direct* (not refused); it ignores NEDNSSettings/NEProxySettings and does
not support `includeAllNetworks`/`enforceRoutes`. DTS guidance: content filter to *observe*, transparent
proxy to *modify*.[^42]

### 9. Wi-Fi diagnostics: wdutil, sysdiagnose, Wireless Diagnostics

**airport is dead.** The private-framework `airport` CLI was functionally killed in macOS 14.4 (March 2024):
it prints only a deprecation notice pointing at Wireless Diagnostics and `wdutil`; the binary is entirely
absent on current macOS.[^43] Programmatic replacement is CoreWLAN, which now effectively requires Location
authorization (and in practice a signed app bundle).

**wdutil** (every subcommand requires sudo):[^44]

```bash
sudo wdutil info        # NETWORK (primary service, DNS), WIFI (SSID/BSSID/RSSI/noise/Tx rate/PHY/
                        # MCS/NSS/channel), BLUETOOTH, AWDL sections
sudo wdutil diagnose -q # live tests + ~50-file bundle → /var/tmp/WirelessDiagnostics_*.tar.gz
```

**Redaction:** SSID/BSSID/MAC are location-identifying, so `wdutil info` prints `<redacted>` even under sudo
unless the calling context holds Location Services permission. TCC grants Location to *signed app bundles*,
not bare CLI binaries; workarounds are a small signed GUI wrapper (e.g. wifi-unredactor) or
`system_profiler SPAirPortDataType`, which still shows the current SSID (not BSSID) without
permission.[^45] The workaround landscape is unstable across point releases: `networksetup
-getairportnetwork` died in 15.0, the `ipconfig getsummary` SSID trick died in 15.6.[^15]

**sysdiagnose:** `sudo sysdiagnose` (or Control-Option-Command-Shift-Period) → `/var/tmp/
sysdiagnose_<date>.tar.gz`, containing netstat/routing dumps, `ifconfig.txt`, a `network-info` folder
(SystemConfiguration plists), Wi-Fi state, and `system_logs.logarchive` (open in Console.app — the most
important artifact per Apple DTS). The archive can contain names, IPs, and paths; scrub before
sharing.[^46]

**Wireless Diagnostics.app** (Option-click the Wi-Fi menu icon): the Assistant writes a report to
`/var/tmp`; the **Window menu** hides the real tools: **Scan** (the surviving `airport -s` replacement:
SSIDs, channels, security, signal) and **Sniffer** (raw 802.11 captures to `.wcap` in `/var/tmp`, readable in
Wireshark).[^47]

### 10. CLI mapping: macOS vs Linux

| Task | macOS (Darwin) | Linux |
| --- | --- | --- |
| Routing table | `netstat -rn` | `ip route` |
| Single route lookup | `route -n get <dest>` | `ip route get <dest>` |
| Listening sockets + process | `lsof -nP -iTCP -sTCP:LISTEN`; `netstat -anv` (has pid column) | `ss -tlnp` |
| Interface config | `ifconfig` (no `ip` binary exists) | `ip addr` |
| IPv4 neighbors | `arp -a` | `ip neigh` |
| IPv6 neighbors | `ndp -an` | `ip -6 neigh` |
| Per-process traffic | `nettop` | `iftop`/`nethogs` (3rd-party) |
| Speed/bufferbloat test | `networkQuality` (built-in) | none built-in |
| Kernel net state | `sysctl net.*` — **no /proc on macOS** | `/proc/net/*`, `sysctl` |
| Which process owns a port | `lsof -i tcp:<port>` | `ss -tlnp sport = :<port>` / `fuser` |

There is **no `ss` and no `ip`** on macOS.[^48] `networkQuality` (macOS 12+) measures capacity plus
responsiveness in RPM (round trips per minute under load); useful flags: `-v`, `-c` (JSON), `-I en0` (bind
interface), `-p` (test through Private Relay), `-s` (sequential).[^49]

## Methodology: a macOS network diagnosis workflow

1. **Establish which interface is primary and who supplies DNS:** `scutil --nwi`, then
   `scutil --dns | head -40`. A `utun` first = a VPN owns default route/DNS.
2. **Confirm the route the kernel will pick:** `route -n get <dest>` (compare against `netstat -rn`).
3. **Test name resolution on BOTH paths:** `dig <name>` (direct path) vs
   `dscacheutil -q host -a name <name>` / `dns-sd -G v4v6 <name>` (system path). Divergence localizes the
   fault: system-path-only failure → resolver config (match domains, scoped resolvers, /etc/resolver);
   dig-only failure → resolv.conf stub or upstream server, usually a red herring.
4. **Inspect resolver provenance:** in `scutil --dns`, which entry *should* match the failing domain? Is it
   Supplemental (match domain) or only Scoped? What `order` does it carry? Then check
   `State:/Network/Service/<id>/DNS` in interactive scutil for the publishing service.
5. **Check interception layers:** `systemextensionsctl list` (content filters, transparent proxies, DNS
   proxies), `pfctl -s References && pfctl -sr -a 'com.apple/*'`, Private Relay status, "Limit IP Address
   Tracking".
6. **Wi-Fi layer:** `sudo wdutil info` (RSSI/noise/CCA/channel; AWDL state), Wireless Diagnostics → Scan.
7. **Escalate to evidence capture:** `sudo wdutil diagnose`, `sudo sysdiagnose`, and
   `log stream --predicate 'process == "mDNSResponder"'` while reproducing.
8. **Reset state cheaply before chasing ghosts:** flush DNS, toggle the interface, switch to a clean network
   Location, reboot; network transitions and sleep legitimately clear resolver/mDNS state.

## Practical Patterns

- **Find the physical interface under a VPN:** `scutil --nwi` (first non-utun entry); or walk
  `State:/Network/Global/IPv4` → `PrimaryService` → service keys in interactive scutil.
- **See DHCP-supplied DNS** (which `networksetup -getdnsservers` hides): `ipconfig getpacket en0` →
  `domain_name_server`, or `ipconfig getsummary en0`.
- **Make unicast `.local` work:** `/etc/resolver/local` with `nameserver <ip>` and `search_order 100`.
- **Dev-TLD resolver:** `/etc/resolver/test` → `nameserver 127.0.0.1` (+ `port 53535` if needed); verify in
  `scutil --dns`, test with `dscacheutil`, never with dig. (Watch item: macOS 26 custom-TLD regression.[^25])
- **Persistent PF rules:** custom anchor file in `/etc/pf.anchors/`, referenced from a copy of the stock
  ruleset, loaded by your own LaunchDaemon with `pfctl -E -f`; never edit `/etc/pf.conf` inline (updates
  overwrite it).
- **Script-readable SSID on modern macOS:** `system_profiler SPAirPortDataType` (SSID only), or a signed
  Location-authorized helper; budget for the approach breaking at point releases.
- **Watch config changes live while reproducing:** interactive `scutil` → `n.add State:/Network/Global/IPv4`
  → `n.watch`; or `log stream --predicate 'subsystem == "com.apple.SystemConfiguration"'`.
- **JSON capacity/latency evidence for a ticket:** `networkQuality -v -c`.

## Anti-Patterns

- **Testing system DNS with dig/nslookup** — they bypass everything that makes macOS DNS macOS.
- **Editing `/etc/resolv.conf`** — regenerated on every network event; split DNS can't live there.
- **Hand-editing `preferences.plist` or `com.apple.alf.plist`** — private formats; Sequoia removed the ALF
  plist entirely.
- **Inline edits to `/etc/pf.conf`** or to `com.apple.pfctl.plist` (the latter needs SIP off and updates
  revert it) — use your own anchor + LaunchDaemon.
- **Scripting against `airport`** — gone. Migrate to wdutil/CoreWLAN/system_profiler.
- **Assuming `scutil --dns` order == runtime behavior** — encrypted-transport preference and server
  penalization can reorder reality; confirm with mDNSResponder logs when it matters.
- **`kill -9` on mDNSResponder as routine hygiene** — HUP suffices; `-9` is for a wedged daemon only.
- **Relying on `PrimaryInterface` being physical** — full-tunnel VPNs legitimately make it `utunN`; naive
  tooling breaks.

## Troubleshooting: macOS-specific failure modes

**VPN split DNS fails.** Symptom matrix: *dig resolves internal names but browsers fail* → scoped/match
config broken while resolv.conf points at the VPN DNS; *browsers work but dig fails* → split DNS is correct
and resolv.conf only carries the LAN resolver (expected; not a bug). Root causes, in observed frequency:
client set search domains but no match domains (suffix completion ≠ query routing); supplemental resolver
degraded to scoped-only (Ventura 13.3 mobileconfig regression class); resolver pinned to the wrong if_index
after sleep/reconnect; Private Relay inserting itself (single report). Workflow: §Methodology steps 3–5;
remedies: fix client matchDomains, flush DNS, `/etc/resolver/<domain>` as a stopgap.[^38][^39][^40]

**Captive portal never appears.** Detection: macOS probes `http://captive.apple.com/hotspot-detect.html`
(UA `CaptiveNetworkSupport-… wispr`); any non-`Success` answer pops the Captive Network Assistant
(`/System/Library/CoreServices/Captive Network Assistant.app`); timeout = "no internet, no popup".[^50]
CNA-suppression causes: hard-coded DNS on the interface defeating the portal's DNS hijack (fix: clean
"Automatic" Location); the Catalina-era autostart bug (fix: open CNA manually or browse to
`http://captive.apple.com`); content-filter agents intercepting the probe; **on the network side**,
whitelisting `captive.apple.com` in the walled garden, which makes devices believe they're online.[^51]
Since Big Sur, macOS also supports the IETF Captive Portal API (DHCP/RA-advertised `captive+json`
endpoint).[^52]

**iCloud Private Relay interference.** PR (iCloud+, macOS 12+) routes Safari + insecure-HTTP + DNS through a
two-hop MASQUE/QUIC proxy (`mask.icloud.com` / `mask-h2.icloud.com`); destination sees a shared egress IP
(daily CSV at `mask-api.icloud.com/egress-ip-ranges.csv`); PR DNS ignores local/VPN resolver settings for
covered traffic.[^53] An on-device NE VPN takes precedence (PR traffic inside the tunnel isn't eligible),
**but** Mullvad measured PR QUIC heartbeats bypassing PF rules and the routing table outside the tunnel:
"VPN disables Private Relay" is only mostly true.[^54] Sanctioned network-side block: return NXDOMAIN (never
a timeout/drop) for the two mask hostnames; users then get a per-network disable prompt. Per-network user
disable: toggle off "Limit IP Address Tracking" (which also, without any iCloud+ subscription, stops
known-tracker proxying for Mail/Safari, a separate setting that has itself caused captive-portal and
connectivity oddities).[^55]

**mDNSResponder cache/state.** Stale answers after VPN changes or network transitions → flush incantation
(§6); remember negative caching and app-internal caches before blaming the daemon; sleep/wake and interface
bounce clear mDNS records by design (QA1310).[^24][^26]

**Per-app Local Network permission blocks LAN/mDNS (macOS 15+).** Sequoia brought iOS-style per-app
Local Network privacy gating to macOS: an app denied the permission silently fails to reach LAN hosts or
browse Bonjour, so `ping` to a printer or `dns-sd -B` can fail from one terminal app and work from another.
Check System Settings → Privacy & Security → Local Network for the calling app before debugging the
network itself (Apple documents the control in the Mac User Guide[^57]; behavior reports are
practitioner-corroborated but the gating UX has shifted across 15.x point releases — treat specifics as
qualified).

**AWDL (awdl0) latency spikes.** AirDrop/AirPlay/Handoff hop the single Wi-Fi radio to social channels
(6 / 44 / 149), stalling AP traffic: periodic 50–200 ms (reported up to ~1 s) ping spikes correlated with
`AWDL_PEER_PRESENCE (awdl0)` log events. `sudo ifconfig awdl0 down` works but is reverted within seconds on
Ventura+ (practitioner consensus); durable mitigations: put the AP on AWDL's preferred channel (149 where
legal, 44 EU), disable AirDrop/Bluetooth, accept it.[^56]

**SSID `<redacted>` everywhere.** Not a bug: TCC/Location gating (§9). Use a signed Location-authorized
helper or `system_profiler SPAirPortDataType`.

## References

[^1]: configd(8) — https://real-world-systems.com/docs/configd.1.html
[^2]: Apple, *System Configuration Programming Guidelines* (archived ADC PDF) — https://leopard-adc.pepas.com/documentation/Networking/Conceptual/SystemConfigFrameworks/SystemConfigFrameworks.pdf
[^3]: Frank Denis, *Programmatically changing network configuration on OS X* — https://00f.net/2011/08/14/programmatically-changing-network-configuration-on-osx/
[^4]: Apple, *Change the order of network services on Mac* — https://support.apple.com/guide/mac-help/mchlp2711/mac
[^5]: floyd.ch — PrimaryRank demotion of VPN services — https://www.floyd.ch/?p=1342
[^6]: apple.stackexchange — service order ignored (Mojave/Big Sur reports) — https://apple.stackexchange.com/questions/349903/
[^7]: scutil(8) — https://ss64.com/mac/scutil.html ; https://keith.github.io/xcode-man-pages/scutil.8.html
[^8]: Rakhesh Sasidharan, *macOS VPN doesn't use the VPN DNS* (scutil --dns anatomy) — https://rakhesh.com/infrastructure/macos-vpn-doesnt-use-the-vpn-dns/
[^9]: Mike Bianco, *Understanding DNS requests on macOS* — https://mikebian.co/understanding-dns-requests-on-macos/
[^10]: thewayeye, *Finding the underlying physical network interface of a VPN* — https://thewayeye.net/posts/underlying-physical-network-interface-vpn/
[^11]: networksetup(8) — https://ss64.com/mac/networksetup.html
[^12]: Kandji, *Managing network settings on macOS Big Sur* — https://blog.kandji.io/managing-network-settings-on-macos-big-sur-and-mac-address-randomization-in-ios-14
[^13]: openconnect vpnc-scripts issue #48 (networksetup blind to utun) — https://gitlab.com/openconnect/vpnc-scripts/-/issues/48
[^14]: ipconfig(8) — https://ss64.com/mac/ipconfig.html ; https://keith.github.io/xcode-man-pages/ipconfig.8.html
[^15]: StackOverflow, *Get current Wi-Fi SSID on macOS Sequoia+* (15.0/15.6 breakage timeline) — https://stackoverflow.com/questions/78994709/
[^16]: mDNSResponder(8) + Apple open source README — https://manpagez.com/man/8/mDNSResponder/osx-10.11.6.php ; https://github.com/apple-oss-distributions/mDNSResponder
[^17]: Julia Evans, *getaddrinfo is kind of weird* — https://jvns.ca/blog/2022/02/23/getaddrinfo-is-kind-of-weird/ ; A. Jesse Jiryu Davis — https://emptysqua.re/blog/getaddrinfo-cpython-mac-and-bsd/
[^18]: Ars Technica discoveryd coverage (Jan/May/Jun 2015) — https://arstechnica.com/gadgets/2015/01/why-dns-in-os-x-10-10-is-broken-and-what-you-can-do-to-fix-it/ ; https://arstechnica.com/gadgets/2015/05/new-os-x-beta-dumps-discoveryd-restores-mdnsresponder-to-fix-dns-bugs/ ; Craig Hockenberry — https://furbo.org/2015/05/05/discoveryd-clusterfuck/
[^19]: Apple macnetworkprog list (dig warning; dns-sd -G full path) — https://lists.apple.com/archives/macnetworkprog/2012/Jul/msg00010.html ; Gordon Davisson — https://stackoverflow.com/questions/50914268/
[^20]: Tailscale issue #4413 (resolv.conf cannot express split DNS) — https://github.com/tailscale/tailscale/issues/4413
[^21]: kunall.is, *mDNSResponder state dump deprecation* — https://kunall.is/posts/mdnsresponder/ ; corroborated — https://rbf.dev/blog/2024/01/how-apple-accidentally-broke-my-spotify/
[^22]: Apple, *If you use a network with a .local top-level domain* — https://support.apple.com/en-us/101903 ; RFC 6762 — https://datatracker.ietf.org/doc/html/rfc6762
[^23]: resolver(5) — https://www.manpagez.com/man/5/resolver/
[^24]: Apple Technical Q&A QA1310 (stale mDNS records, sleep clears cache) — https://developer.apple.com/library/archive/qa/qa1310/_index.html
[^25]: adamamyl gist, *macOS 26 /etc/resolver custom-TLD breakage* (single source) — https://gist.github.com/adamamyl/81b78eced40feae50eae7c4f3bec1f5a
[^26]: apple.stackexchange, *Flush DNS cache on Sierra/High Sierra* (dscacheutil does nothing for DNS) — https://apple.stackexchange.com/questions/303110/ ; animatedcreativity layered-cache analysis — https://animatedcreativity.com/tutorials/macos-dns-caching-why-etc-hosts-edits-dont-take-effect-immediately-and-the-dscacheutil-mdn/
[^27]: Apple HT202516, *Reset the DNS cache* (per-version commands) — https://support.apple.com/en-us/HT202516
[^28]: PF-on-macOS history — https://manjusri.ucsc.edu/2015/03/10/PF-on-Mac-OS-X/ ; PaperCut KB (ipfw removed 10.10) — https://www.papercut.com/kb/Main/MacPortForwarding
[^29]: stock pf.conf incl. -E/-X contract comment — https://github.com/essandess/macOS-Fortress/blob/master/pf.conf ; Apple `container` project anchor ordering — https://github.com/apple/container/blob/129c2dc9/Sources/Services/ContainerAPIService/Client/PacketFilter.swift
[^30]: pfctl(8) — https://keith.github.io/xcode-man-pages/pfctl.8.html
[^31]: iyanmv, *Setting up PF correctly on macOS* (pf.conf overwritten by updates) — https://iyanmv.medium.com/setting-up-correctly-packet-filter-pf-firewall-on-any-macos-from-sierra-to-big-sur-47e70e062a0e ; WatchPaths reload pattern — https://stackoverflow.com/questions/45475242/
[^32]: Apple, macOS Sequoia 15 enterprise release notes (ALF plist removal) — https://support.apple.com/en-us/121011 ; Apple DTS thread — https://developer.apple.com/forums/thread/760135 ; Fleet/osquery fallout — https://github.com/fleetdm/fleet/issues/21802
[^33]: Apple, NetworkExtension framework docs — https://developer.apple.com/documentation/networkextension ; content-filter providers — https://developer.apple.com/documentation/networkextension/content-filter-providers
[^34]: Apple, *Deprecated kernel extensions and system extension alternatives* — https://developer.apple.com/support/kernel-extensions/
[^35]: Apple, *Routing your VPN network traffic* — https://developer.apple.com/documentation/networkextension/routing-your-vpn-network-traffic ; includeAllNetworks — https://developer.apple.com/documentation/NetworkExtension/NEVPNProtocol/includeAllNetworks ; DTS on enforceRoutes — https://developer.apple.com/forums/thread/832079
[^36]: Mullvad, *Why we still don't use includeAllNetworks* (Mar 2025) — https://mullvad.net/uk/blog/why-we-still-dont-use-includeallnetworks ; profile-conflict thread — https://developer.apple.com/forums/thread/669086
[^37]: utun lifecycle / stale interfaces — https://developer.apple.com/forums/thread/719060 ; https://apple.stackexchange.com/questions/484365/ ; PrimaryInterface=utun — https://github.com/exelban/stats/issues/2143
[^38]: NEDNSSettings.matchDomains — https://developer.apple.com/documentation/networkextension/nednssettings/matchdomains ; matchDomainsNoSearch semantics (WireGuard list) — https://lists.zx2c4.com/pipermail/wireguard/2021-November/007288.html
[^39]: apple.stackexchange, *Split DNS broken with mobileconfig VPN profile* (Ventura ≥13.3) — https://apple.stackexchange.com/questions/450523/
[^40]: Versa KB (same-DNS-IP teardown no-op) — https://support.versa-networks.com/support/solutions/articles/23000029468-dns-resolution-failure-after-vpn-disconnect-macos- ; Tunnelblick #736 (if_index pinning) — https://github.com/Tunnelblick/Tunnelblick/issues/736
[^41]: Jamf, *The network is back in Big Sur* — https://www.jamf.com/blog/the-network-is-back-in-big-sur/ ; ZDNet — https://www.zdnet.com/article/apple-removes-feature-that-allowed-its-apps-to-bypass-macos-firewalls-and-vpns/ ; Patrick Wardle — https://x.com/patrickwardle/status/1349488392732491776
[^42]: NETransparentProxyProvider — https://developer.apple.com/documentation/NetworkExtension/NETransparentProxyProvider ; DTS guidance — https://developer.apple.com/forums/thread/761134
[^43]: Intuitibits, *Goodbye airport* — https://www.intuitibits.com/2024/03/14/goodbye-airport/ ; Apple forums — https://developer.apple.com/forums/thread/748161
[^44]: NetBeez, *wdutil for troubleshooting Wi-Fi on macOS* — https://netbeez.net/blog/wdutil-for-troubleshooting-wi-fi-issues-on-macos/
[^45]: wifi-unredactor (TCC/Location mechanics) — https://github.com/noperator/wifi-unredactor ; Go/TCC signing analysis — https://dev.to/jaisonerick/reading-wi-fi-data-from-go-on-macos-after-apple-removed-airport-19g
[^46]: sysdiagnose(1) — https://keith.github.io/xcode-man-pages/sysdiagnose.1.html ; Howard Oakley — https://eclecticlight.co/2016/02/08/more-useful-information-gleaned-from-sysdiagnose/ ; Apple DTS logarchive guidance — https://developer.apple.com/forums/thread/739560
[^47]: Apple, *Use Wireless Diagnostics on Mac* — https://support.apple.com/guide/mac-help/mchlf4de377f/mac ; Howard Oakley — https://eclecticlight.co/2021/12/23/coreservices-apps-wireless-diagnostics/
[^48]: macOS port/process lookups — https://stackoverflow.com/questions/3855127/ ; https://superuser.com/questions/627391/
[^49]: networkQuality(8) — https://manp.gs/mac/8/networkQuality ; TidBITS RPM explainer — https://tidbits.com/2022/04/22/use-apples-networkquality-tool-to-test-internet-responsiveness/
[^50]: Captive portal behavior (WBA) — https://captivebehavior.wballiance.com/ ; Purple guide — https://www.purple.ai/en-gb/guides/why-captive-portal-isnt-loading-on-iphone
[^51]: CNA failure threads — https://apple.stackexchange.com/questions/376401/ ; https://discussions.apple.com/thread/250907026 ; Graham Pugh (CNA internals/suppression) — https://grahamrpugh.com/2014/10/29/undocumented-change-to-captive-network-assistant-settings-in-yosemite.html
[^52]: Apple, *How to modernize your captive network* (Captive Portal API) — https://developer.apple.com/news/?id=q78sq5rv
[^53]: Apple, *Prepare your network for iCloud Private Relay* — https://developer.apple.com/icloud/prepare-your-network-for-icloud-private-relay/ ; Apple PR whitepaper (Dec 2021) — https://www.apple.com/privacy/docs/iCloud_Private_Relay_Overview_Dec2021.PDF ; measurement study — https://export.arxiv.org/pdf/2207.02112v3.pdf
[^54]: Mullvad, *Apple's Private Relay can cause the system to ignore firewall rules* — https://mullvad.net/en/blog/apples-private-relay-can-cause-the-system-to-ignore-firewall-rules ; DTS on VPN precedence — https://developer.apple.com/forums/thread/682274
[^55]: Apple, *About iCloud Private Relay* (per-network disable) — https://support.apple.com/en-us/102022 ; TidBITS interlocking-privacy-settings guide — https://tidbits.com/2022/06/20/solving-connectivity-problems-caused-by-interlocking-apple-privacy-settings/
[^56]: networkweather, *AWDL explained* — https://www.networkweather.com/learn/awdl/ ; Open Wireless Link research (TU Darmstadt) — https://tuprints.ulb.tu-darmstadt.de/bitstreams/ff8e8047-5d3c-4ae2-9f67-1781b467e723/download ; Ventura workaround breakage — https://apple.stackexchange.com/questions/451646/
[^57]: Apple, *Control access to your local network on Mac* — https://support.apple.com/guide/mac-help/control-access-to-your-local-network-mchla4f60997/mac
