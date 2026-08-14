# cloudflare-platform

**Category:** Science, Biology & Medicine
**Platform:** Claude
**Original Path:** claude/standalone/cloudflare-platform

## Description
Cloudflare platform from the network/infrastructure angle — orange-cloud reverse proxy effects on DNS answers, source IPs, TLS; authoritative DNS + 1.1.1.1; Tunnel (cloudflared); WARP; Access/ZTNA; Magic Transit/WAN; load balancing & Argo; Workers/R2 as infrastructure; caching & Cache Rules; 52x, cf-ray/trace debugging. TRIGGER: site behind Cloudflare misbehaving (520-526, cf-cache-status, cf-ray, /cdn-cgi/trace); origin sees only Cloudflare IPs (CF-Connecting-IP); SSL modes (Flexible vs Full strict), 525/526; proxied vs DNS-only, CNAME flattening, MX grey-cloud; 1.1.1.1/DoH no-ECS geo effects; cloudflared Tunnel + Access; WARP split tunnels & posture; Magic Transit BGP/GRE/MTU; LB steering & health checks; Argo; R2 zero-egress CDN origin; purge, tiered cache. SKIP: Workers JS runtime, compat flags, wrangler → lang-js-ts; DNS/BGP protocol internals → networking; host-level net tooling → devops-linux-admin; ZTNA concepts/NIST 800-207, Okta/IdP internals → security-review / okta-expert.

---

# Cloudflare Platform (network/infrastructure angle)

Researched 2026-06-12 from Cloudflare docs/blog + independent practitioner sources (3+
independent sources per concept; volatile values stamped `verified-as-of: 2026-06-12`).
Scope: network/infrastructure platform. Workers **JS runtime** (APIs, compat flags,
wrangler) owned by `lang-js-ts` (references/javascript-runtimes-deno-bun-edge.md) —
cross-reference, no duplicate.

## The one mental model

DNS record **proxied** (orange cloud) swaps origin IP for Cloudflare **anycast
edge IPs** in every DNS answer. From there: clients connect (TCP+TLS) to nearest
Cloudflare PoP, Cloudflare opens *second, independent* connection to origin, origin
only sees connections from Cloudflare published IP ranges. Every other behavior
in skill — client-IP restoration headers, two independent TLS sessions + encryption
modes, edge caching, 52x errors on edge→origin leg, traceroute terminating at
PoP — falls from that one fact. DNS-only (grey cloud) records bypass all, leak
origin IP.

## Quick reference (highest-frequency facts)

| Question | Answer | Detail |
|---|---|---|
| Origin sees wrong client IP | Read `CF-Connecting-IP` (or XFF); trust ONLY from Cloudflare ranges; mod_remoteip / ngx realip | references/proxy-tls-caching.md |
| Which SSL mode | Full (strict) only end-to-end-verified mode; Flexible = plaintext to origin + redirect loops | references/proxy-tls-caching.md |
| HTML not caching | By design — caching by file extension; HTML needs Cache Rule ("Eligible for cache" = Cache Everything) | references/proxy-tls-caching.md |
| `cf-cache-status: DYNAMIC` | "Not eligible, no rule makes eligible" — normal for HTML, not error | references/proxy-tls-caching.md |
| 521 vs 522 vs 524 | refused / TCP timeout (19 s SYN) / connected but no HTTP response in 120 s | references/dns-and-debugging.md |
| 525 vs 526 | TLS handshake to origin failed / cert invalid under Full (strict) | references/proxy-tls-caching.md, dns-and-debugging.md |
| Locate serving PoP | `cf-ray: <id>-IATA` header; `/cdn-cgi/trace` → `colo=` | references/dns-and-debugging.md |
| Bypass proxy to test origin | `curl -I --resolve host:443:ORIGIN_IP https://host/` or "Pause Cloudflare on Site" | references/dns-and-debugging.md |
| Why mail breaks when proxied | Proxied records answer with edge IPs speaking only HTTP/S; MX targets must stay DNS-only (leak IP) | references/dns-and-debugging.md |
| Expose internal app safely | cloudflared Tunnel (outbound-only, port 7844) + Access policy; validate `Cf-Access-Jwt-Assertion` at origin | references/zero-trust-tunnel-warp-access.md |
| WARP can't reach private net | Carve corporate CIDR out of default RFC1918 split-tunnel **exclude** list | references/zero-trust-tunnel-warp-access.md |
| Magic Transit breakage after onboarding | Almost always MTU/MSS: GRE MTU 1476, clamp IPv4 MSS to 1436 | references/magic-transit-lb-workers-r2.md |
| R2 "free egress" | Egress free; every GET still billed Class B op unless cached | references/magic-transit-lb-workers-r2.md |

## Routing table

| Topic | Reference |
|---|---|
| Orange-cloud proxy fundamentals (DNS answers, source IPs, client-IP restoration, ports, origin exposure), TLS modes (Universal SSL, Origin CA, Authenticated Origin Pulls, 525/526), caching defaults, cf-cache-status, Cache Rules, purge, Tiered Cache, Cache Reserve | `references/proxy-tls-caching.md` |
| Authoritative DNS (full/partial setup, proxied vs DNS-only, CNAME flattening, DNSSEC, Foundation DNS), 1.1.1.1 resolver (DoH/DoT, filtering variants, no-ECS consequences, July 2025 outage), 52x/1xxx error catalog, cf-ray anatomy, /cdn-cgi/trace, debugging workflow | `references/dns-and-debugging.md` |
| Cloudflare Tunnel/cloudflared (architecture, tokens, ingress rules, replicas, private networks, failure modes), WARP client (modes, split tunnels, posture, CGNAT addressing), Access/ZTNA (apps, policies, IdPs, service tokens, JWT validation, SSH) | `references/zero-trust-tunnel-warp-access.md` |
| Magic Transit & Magic WAN (BGP prefix advertisement, anycast GRE/IPsec/CNI, health checks, MTU/MSS), Load Balancing (pools/monitors/steering/affinity), Argo Smart Routing & Tiered Cache history, Workers infra angle (routes, isolates, Smart Placement), R2 (zero egress, public buckets, limits) | `references/magic-transit-lb-workers-r2.md` |

## Cross-references

- **Workers as JS runtime** (V8 isolates API surface, nodejs_compat, compatibility dates,
  wrangler workflow): `lang-js-ts` (references/javascript-runtimes-deno-bun-edge.md) —
  that hub owns runtime; this skill owns where Workers sit in request path.
- **Python Workers** (Pyodide at edge): `lang-python` (references/python-in-browser-wasm.md).
- **AI Gateway / LLM proxying on Cloudflare**: `ai-llm-model-layer` (llm-ai-gateways concept).
- **ZTNA conceptual model & IdP side** (NIST 800-207, Okta integration): `security-review` /
  `okta-expert`.
- **DNS protocol internals (DNSSEC, resolution path, DoH/DoT mechanics)**: `networking` hub
  (references/dns-deep-dive.md; BGP/transport siblings landing as hub grows).
- **Host-level network tool recipes (ip/ss/dig/tcpdump/iptables)**: `devops-linux-admin`.

## Sources

Per-concept numbered source lists at bottom of each reference file. Primary corpus:
developers.cloudflare.com (product truth, one source), blog.cloudflare.com engineering posts,
independent practitioner/forum sources (community.cloudflare.com, GitHub issues, Stack
Exchange, third-party engineering blogs) — 60+ distinct sources across four references.