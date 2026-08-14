# Zero Trust connectivity: Tunnel (cloudflared), WARP, and Access

Provenance: /dr deep-research run 2026-06-12 (skill: cloudflare-platform). Volatile values
(ports, replica caps, plan gates, mode names) `verified-as-of: 2026-06-12`.

## Contents

- [Cloudflare Tunnel (cloudflared)](#cloudflare-tunnel-cloudflared)
- [Replicas and high availability](#replicas-and-high-availability)
- [Private network mode](#private-network-mode)
- [Tunnel operations and failure modes](#tunnel-operations-and-failure-modes)
- [WARP client](#warp-client)
- [Cloudflare Access / ZTNA](#cloudflare-access--ztna)
- [JWT validation at the origin](#jwt-validation-at-the-origin)
- [The Tunnel + Access pattern, and its limits](#the-tunnel--access-pattern-and-its-limits)
- [References](#references)

## Cloudflare Tunnel (cloudflared)

Tunnel inverts the exposure model: a Go daemon (`cloudflared`) on the origin host dials
**outbound-only** to the edge — no inbound firewall ports, no public A record, no origin IP
exposure; works behind CGNAT and double NAT [^1][^10][^19][^23][^24]. The daemon connects to
**port 7844, preferring QUIC (UDP 7844) with fallback to HTTP/2 (TCP 7844)** — allow egress
TCP+UDP 7844 [^1][^2][^13]. One tunnel = **four outbound connections to four distinct edge
servers across at least two data centers** [^1][^3][^10][^12].

**Types and credentials:**

- **Quick tunnels** (`cloudflared tunnel --url`) — no account, ephemeral `trycloudflare.com`
  URL; dev/demo only [^1][^28].
- **Named tunnels** persist with a UUID. Auth: (a) **credentials JSON** (UUID.json + cert.pem
  from `cloudflared tunnel login`) for **locally-managed** tunnels — YAML config with
  **ingress rules** (ordered hostname/path matchers → local services; a final catch-all,
  commonly `http_status:404`, is required); or (b) a **tunnel token** (single base64 string)
  for **remotely-managed** tunnels — config stored in Cloudflare's control plane, edited via
  dashboard/API [^1][^4][^10][^14]. Public hostname routes create CNAMEs to
  `<UUID>.cfargotunnel.com` [^1][^12]. Runs as a systemd/Windows service; one instance per
  host when service-installed [^1].

## Replicas and high availability

The same tunnel can run from multiple `cloudflared` processes ("replicas"), each adding 4 edge
connections; cap **25 replicas / 100 connections per tunnel** (`verified-as-of: 2026-06-12`)
[^1][^3]. Replicas are failover, **not traffic steering**: Cloudflare forwards to the
geographically closest replica and retries others on failure. For deterministic distribution
or session affinity, use Cloudflare Load Balancer with separate tunnel UUIDs per pool endpoint
[^1][^3][^12]. Guidance: keep replicas of one tunnel in the same physical location; split IP
space across multiple tunnels for control-plane redundancy and to avoid port exhaustion on
large private routes [^1].

## Private network mode

Beyond public hostnames, a tunnel can advertise **IP/CIDR routes** (RFC1918 or public) mapped
to it; WARP-enrolled clients then reach those subnets, with virtual networks disambiguating
overlapping CIDRs. **Top misconfiguration:** WARP's split-tunnel config **excludes RFC1918 by
default**, so you must carve your corporate CIDR out of the exclude list (or switch to include
mode) before private routing works [^1][^6]. The tunnel is asymmetric: server-initiated
(reverse-path) connections don't traverse it [^1].

## Tunnel operations and failure modes

`cloudflared` exposes a **Prometheus metrics endpoint** (default `127.0.0.1:<first free port
20241–20245>/metrics`; `0.0.0.0` in containers; `--metrics` to set) with HA-connection count,
request/error counters, QUIC RTT/loss [^1]. `cloudflared tunnel info` lists connector IDs and
colos [^1].

Documented failure modes (qualified — GitHub issues + community):

- UDP 7844 blocked → QUIC fails; fallback behavior is inconsistent (reports of both
  failure-to-fall-back and permanent fallback to HTTP/2, which breaks UDP/private-DNS
  features) [^13][^14].
- Throughput consistently below direct port-forwarding (one benchmark: 97/30 Mbps tunneled vs
  110/70 direct); long-haul non-US transfers measured 3–4× slower [^14].
- `cloudflared` memory grows unboundedly when request volume exceeds origin upload bandwidth
  (issue #1205 — no built-in backpressure) [^13].
- Tunnels can fail silently without external monitoring [^24]. Practitioners avoid Tunnel for
  latency-sensitive/game/bulk-media traffic [^23].

## WARP client

WARP is the device agent (now branded **"Cloudflare One Client"**;
`verified-as-of: 2026-06-12`). Three channels: the L3 tunnel (**WireGuard or MASQUE, both
UDP**), DoH for DNS policy (inside the tunnel), and HTTPS device orchestration [^1][^11].
**MASQUE (proxying over HTTP/3/QUIC on 443) is now the default protocol**; WireGuard remains
available [^1][^5].

**Modes** (legacy → current names, `verified-as-of: 2026-06-12`) [^1]:

| Legacy name | Current name | Behavior |
|---|---|---|
| Gateway with WARP | Traffic and DNS mode (`WarpWithDnsOverHttps`) | full tunnel: DNS + network + HTTP filtering, posture, DLP |
| Gateway with DoH | DNS only mode (`DnsOverHttps`) | DNS filtering only; no traffic routing |
| SWG without DNS filtering | Traffic only mode (`TunnelOnly`) | routes traffic; OS keeps its DNS |
| Proxy mode | Local proxy mode (`WarpProxy`) | localhost SOCKS5/HTTP proxy, default port 40000, MASQUE-only, 10 s request timeout |
| Device Information Only | Posture only mode (`PostureOnly`) | posture signals only |

**Split tunnels:** two mutually exclusive per-profile modes — **Exclude** (default; everything
through WARP except listed IPs/domains; RFC1918 + CGNAT excluded out of the box) and
**Include** (only listed routes traverse WARP) [^1][^11][^14]. Getting the exclude-list
carve-out wrong is the most common "private network unreachable through WARP" cause [^1][^6][^14].

**Enrollment, posture, addressing:** devices enroll into a Zero Trust org (`warp-cli
registration new <team>`, MDM, or GUI), gated by device-enrollment permission policies + IdP
login; system clock skew >20 s invalidates the enrollment JWT [^1][^11]. Posture checks are
client-reported (disk encryption, OS version, firewall, file/cert presence) or pulled from
third-party providers (CrowdStrike, Intune, SentinelOne), then referenced in Access/Gateway
policies [^1][^11]. Addressing: **172.16.0.2 hardcoded under WireGuard**; **MASQUE devices get
a unique IP from CGNAT 100.96.0.0/12** (IPv6 `2606:4700:0cf1:1000::/64`), enabling
WARP-to-WARP and return-path routing; custom RFC1918 pools ≥/24 possible [^1]. Interface MTU
~1280; double-tunnel scenarios need MSS clamping (qualified) [^1].

**Diagnostics:** `warp-cli` (status/settings/registration/debug posture) and `warp-diag` (zip
with `daemon.log` — the key artifact — settings, posture, MASQUE qlogs, optional pcaps);
remote collection via DEX [^1][^11]. The consumer **1.1.1.1/WARP app is a separate product**
— no Gateway policy or org enrollment [^1][^14]. Known conflicts: third-party VPN route/DNS
contention, `systemd-resolved` on Linux, "Happy Eyeballs checks failed" connect loops,
MITM-cert rotation cache issues (qualified) [^1][^14][^15].

## Cloudflare Access / ZTNA

Access is an identity-aware proxy: every request to a protected hostname is
authenticated/authorized **at the edge**, replacing network-location trust. Unauthenticated
browsers are redirected to `/cdn-cgi/access/login/<host>` for an IdP flow; Access then mints a
JWT and sets the `CF_Authorization` cookie [^1][^16]. Apps are **deny by default** [^1].

- **App types:** self-hosted (public hostname, typically behind Tunnel), SaaS (Access brokers
  SAML/OIDC in front of the SaaS app), private/non-HTTP (private IPs/hostnames,
  infrastructure SSH) [^1][^7]. Multiple apps can share a hostname; most-specific path wins [^1][^20].
- **Policies:** one action — Allow / Block / **Bypass** / **Service Auth** — with rule logic
  **Include = OR, Require = AND, Exclude = NOT**. Selectors: email/domain, IdP groups, SAML
  attributes/OIDC claims, IP ranges, country, device posture, WARP/Gateway presence, mTLS
  cert, service token, External Evaluation, user risk score (Enterprise) [^1]. Evaluation
  order: Service Auth & Bypass first (top-down), then Block/Allow top-down; first match wins
  (`verified-as-of: 2026-06-12`) [^1].
- **IdPs:** any SAML or OIDC provider, multiple IdPs simultaneously — including multiple
  instances of the same provider (e.g. several Okta orgs); one-time-PIN email is the zero-IdP
  fallback [^1][^6].
- **Service tokens** (Client ID + one-time-displayed Secret) authenticate machine-to-machine
  via `CF-Access-Client-Id`/`CF-Access-Client-Secret` headers against **Service Auth**
  policies; successful exchange returns a `CF_Authorization` JWT [^1].

## JWT validation at the origin

Access forwards the signed token as **`Cf-Access-Jwt-Assertion`** (header preferred over the
cookie). Origins should validate: RS256 signature against the **JWKS at
`https://<team>.cloudflareaccess.com/cdn-cgi/access/certs`** (match `kid`; don't cache
`public_cert`), `iss` = `https://<team>.cloudflareaccess.com`, and the per-application **AUD
tag** [^1][^16][^17][^18].

**Access alone is insufficient unless the app is reachable only via Tunnel** — anyone who
finds the origin IP bypasses the edge. Mitigations: Tunnel, origin JWT validation
(cloudflared can enforce it via `originRequest.access` with `audTag`), Aegis/CNI, or mTLS
[^1][^8][^16]. Caveat: cloudflared issue #784 — top-level `originRequest.access` was silently
ignored (only per-ingress worked), leaving "global" JWT enforcement unprotected (qualified) [^13].

## The Tunnel + Access pattern, and its limits

The dominant self-hosting pattern: **Tunnel for the data plane + Access for authn at the
edge** — free up to 50 users (`verified-as-of: 2026-06-12`, qualified) [^19][^20][^22][^23].
Conveniences: **App Launcher** portal; **short-lived certificates** replacing static SSH keys
(Access JWT exchanged for an ephemeral cert signed by a Cloudflare CA trusted by sshd);
**browser-rendered SSH/VNC/RDP** (self-hosted public apps only; username must match email
prefix; Bypass/Service Auth unsupported there) [^1][^7].

Vs VPN: no client for browser apps, per-app instead of network-level access, no lateral
movement, central revocation. Costs: ~20–50 ms added latency; **browser-cookie auth breaks
native mobile/desktop apps** that can't inject service-token headers (Home Assistant et al.) —
practitioners run hybrid Access + VPN/Tailscale stacks [^19][^21][^22][^25].

**Dangers:**

- **Bypass policies disable all Access controls and are unlogged.** Prefer Service Auth;
  webhook endpoints commonly force per-path Bypass carve-outs, and ordering mistakes
  (catch-all before bypass paths) silently break integrations [^1][^20].
- Migration anti-patterns: whole-subnet tunnel routes instead of per-app, legacy VPN kept
  "temporarily" forever, posture checks omitted (qualified) [^25].
- The **2025-11-18 Cloudflare outage** failed Access authentication closed for ~3 h (existing
  sessions survived; new logins failed), locking admins out of their own infra — keep
  break-glass procedures [^9][^26].
- Access is ZTNA, **not PAM** — no command/query-level session governance (qualified) [^27].

## References

1. developers.cloudflare.com — cloudflare-one: tunnel-with-firewall, tunnel-availability,
   deploy-replicas, tunnel configuration/metrics/local-management, private-net + add-routes,
   connectivity-options, client-architecture, client modes, troubleshooting/diagnostic-logs,
   device-ips, validating-json, application-token, access policies, service-tokens,
   browser-rendering, identity-providers, learning paths (tier: docs; one source)
2. https://blog.cloudflare.com/getting-cloudflare-tunnels-to-connect-to-the-cloudflare-network-with-quic/ (tier: blog)
3. https://blog.cloudflare.com/highly-available-and-highly-scalable-cloudflare-tunnels/ (tier: blog)
4. https://blog.cloudflare.com/ridiculously-easy-to-use-tunnels/ (tier: blog)
5. https://blog.cloudflare.com/masque-building-a-new-protocol-into-cloudflare-warp/ (tier: blog)
6. https://blog.cloudflare.com/cloudflare-access-for-saas/ — multiple simultaneous IdPs (tier: blog)
7. https://blog.cloudflare.com/browser-ssh-terminal-with-auditing/ — short-lived certs (tier: blog)
8. https://blog.cloudflare.com/access-aegis-cni/ — origin-bypass problem and fixes (tier: blog)
9. https://blog.cloudflare.com/18-november-2025-outage/ — Access fail-closed (tier: blog)
10. cloudsecop.net — tunnel deep-dive: connector internals, keepalive, replica practice (tier: blog)
11. cloudsecop.net — WARP device enrollment, posture providers (tier: blog)
12. zenn.dev/oymk/articles/87316d61b3530a — replica behavior + LB pairing, observed layout (tier: blog)
13. github.com/cloudflare/cloudflared issues #749, #721, #758, #784, #1205, #1309 — QUIC fallback bugs, silent Access-config ignore, unbounded memory (tier: forum)
14. community.cloudflare.com — QUIC→HTTP2 flapping, tunnel speed benchmarks, 1.1.1.1 vs One app (tier: forum)
15. reddit.com r/selfhosted & r/CloudFlare — port 7844 blocked, split-tunnel issues (tier: forum)
16. stackharbor.com — Access flow, JWKS validation, aud pitfalls (tier: blog)
17. reintech.io — Node jwks-rsa origin validation pattern (tier: blog)
18. github.com/ymyzk/cla-jwt-verifier — nginx auth_request JWT sidecar (tier: practitioner code)
19. marios.istos.dev — Tunnel+Access pattern, 50-user free tier, 20–50 ms latency (tier: blog)
20. josemanuelortega.me — Bypass-for-webhooks ordering trap (tier: blog)
21. seoullayer.com — mobile-app incompatibility with browser auth; hybrid stacks (tier: blog)
22. sumguy.com — ingress rules, Tunnel vs Tailscale vs VPN comparison (tier: blog)
23. ebourgess.dev — outbound-only rationale, what not to tunnel (tier: blog)
24. xda-developers.com — CGNAT benefits, silent-failure caution (tier: media)
25. decryptiondigest.com — VPN-migration gap patterns (tier: blog)
26. identityfusion.com + isaaclins.com — Nov-2025 lockout/break-glass lessons (tier: blog)
27. nhimg.org — Access vs PAM gap analysis (tier: blog)
28. blog.thms.uk — quick tunnels (trycloudflare) usage (tier: blog)
