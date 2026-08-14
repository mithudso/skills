# networking

**Category:** Science, Biology & Medicine
**Platform:** Claude
**Original Path:** claude/standalone/networking

## Description
Networking hub — protocol internals & network diagnostics (cloud-connectivity spokes queued). TRIGGER: protocol-level DNS — resolution path (stub/recursive/authoritative, delegation, glue, QNAME minimization), record types (A/AAAA/CNAME/SRV/TXT/NS/SOA/CAA/SVCB-HTTPS), NXDOMAIN/SERVFAIL/REFUSED triage, TTLs & negative caching, DNSSEC validation failures, encrypted DNS (DoT/DoH/DoQ, ECH, Oblivious DoH), split-horizon/private DNS (.internal, home.arpa, BIND views), EDNS(0)/Flag Days/ECS, resolver stack differences (systemd-resolved, mDNSResponder, glibc/musl, ndots), dig/delv/drill/kdig debugging (dns-deep-dive). SKIP: host-level net admin & macOS network config/tooling (ip/ss/tcpdump, scutil/networksetup) → devops-linux-admin; Cloudflare product behavior → cloudflare-platform; AWS VPC/PrivateLink → aws-cloud / mongodb-atlas-expert; MongoDB SRV seedlist → mongodb-expert (references/mongodb-connection-string.md); Node dns/getaddrinfo API usage → lang-js-ts.

---

# networking

Networking hub — protocol internals and network diagnostics. This hub routes to on-demand
reference files under `references/`. Created 2026-06-12 as the anchor for the networking concept
family (first spoke: protocol-level DNS); sibling spokes from the family queue — TCP/transport
internals, routing/BGP, QUIC, TLS internals, IP addressing, network error taxonomy, packet-capture
tooling, plus the cloud-connectivity topics (load balancing, VPN/ZTNA) the description's
"spokes queued" parenthetical refers to — are added as their research runs complete.

## Sub-skill routing table

| Reference | Load when the task involves |
| --- | --- |
| `references/dns-deep-dive.md` | DNS resolution path (stub/recursive/authoritative, delegation, glue, QNAME minimization, priming); record types A/AAAA/CNAME/SRV/TXT/NS/SOA/CAA/PTR/MX/SVCB-HTTPS and apex aliasing; TTLs, negative caching (RFC 2308), serve-stale, resolver TTL clamps; DNSSEC (RRSIG/DNSKEY/DS/NSEC3, chain of trust, SERVFAIL triage, KSK rollover, algorithm guidance); encrypted DNS (DoT/DoH/DoQ, DDR/DNR, Oblivious DoH, ECH); split-horizon & private DNS (BIND views, .internal/home.arpa/.local, RFC 1918 reverse, AS112); EDNS(0) (OPT RR, Flag Days, 1232 bufsize, ECS, Extended DNS Errors); resolver stacks (systemd-resolved, macOS mDNSResponder/scoped resolvers, glibc vs musl, Kubernetes ndots:5); debugging with dig/delv/drill/kdig (+trace caveats, +cd SERVFAIL discrimination, delegation checks) |

## How to use this hub

1. Match the task to a routing-table row and `Read` that reference file for depth.
2. If a networking topic has no row yet, it is either queued for research (see family list above) or
   owned by a neighboring hub — check the SKIP pointers in the description.

## Cross-hub boundaries

| Neighboring owner | Owns |
| --- | --- |
| `devops-linux-admin` | host-level diagnostics & sysadmin recipes (ip, ss, dig-as-cheatsheet, tcpdump, nmap, nftables admin); macOS network config/tooling (scutil, networksetup); Linux networking dataplane (references/linux-networking-stack.md — nftables, netns, tc, XDP) |
| `cloudflare-platform` | Cloudflare products: proxied DNS, 1.1.1.1 as a product, Tunnel, WARP, Magic Transit |
| `aws-cloud` | VPC design, PrivateLink (AWS side) |
| `mongodb-expert` / `mongodb-atlas-expert` | MongoDB SRV seedlist connection strings; Atlas private networking DNS |
| `lang-js-ts` (nodejs-http-networking) | Node.js HTTP/TLS/dgram APIs |