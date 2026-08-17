# Orange-cloud proxy, TLS, and caching

Provenance: /dr deep-research run 2026-06-12 (skill: cloudflare-platform). Volatile values
(port lists, default extension list, TTL tables, pricing, rate limits) `verified-as-of: 2026-06-12`.

## Contents

- [Orange-cloud reverse proxy fundamentals](#orange-cloud-reverse-proxy-fundamentals)
- [Restoring the real client IP](#restoring-the-real-client-ip)
- [Proxied ports](#proxied-ports)
- [TLS through Cloudflare](#tls-through-cloudflare)
- [Caching defaults and cf-cache-status](#caching-defaults-and-cf-cache-status)
- [Cache Rules, TTLs, purge, tiered cache](#cache-rules-ttls-purge-tiered-cache)
- [References](#references)

## Orange-cloud reverse proxy fundamentals

Setting a DNS record to **proxied** changes what authoritative DNS answers: resolvers receive
**Cloudflare anycast IPs** (same ranges announced via BGP from every PoP) instead of the origin
IP [^1][^7]. Only A, AAAA, and CNAME records carrying HTTP/S traffic can be proxied; proxied
CNAMEs are flattened to anycast A/AAAA answers [^1]. **DNS-only (grey cloud)** records return
the real origin value and bypass all Cloudflare processing — every grey-clouded record on the
same origin is a potential origin-IP leak [^1][^2][^3].

Once proxied, the origin only sees TCP connections **from Cloudflare's published ranges**
(cloudflare.com/ips — 15 IPv4 prefixes, e.g. 173.245.48.0/20, 104.16.0.0/13, 172.64.0.0/13,
plus 7 IPv6 prefixes, e.g. 2606:4700::/32) [^1][^3][^4]. Firewall consequences:

1. **Allowlist the Cloudflare ranges** or risk auto-blocking what looks like a handful of IPs
   sending huge volume [^1].
2. **Deny everything else** — attackers find origin IPs via DNS history (SecurityTrails),
   Shodan/Censys certificate-SAN scans, mail bounces, and outbound connections, then hit the
   origin directly with a forged Host header, skipping the WAF [^2][^3][^6].
3. Allow-CF-ranges is necessary but **not sufficient**: Cloudflare egress IPs are shared by
   millions of zones, so any customer's Worker can reach your origin from a "legitimate"
   Cloudflare IP. Cloudflare's own docs grade the allowlist approach only "moderately secure".
   Stronger: Authenticated Origin Pulls, Cloudflare Tunnel (no public IP at all), or Enterprise
   dedicated egress IPs [^1][^2].

**Ping/traceroute to a proxied hostname measures the path to the nearest Cloudflare PoP, not
the origin.** Confirm which colo serves you via `/cdn-cgi/trace`. Anycast latency is
peering-dependent — free-plan zones have been observed routed to distant PoPs (e.g. DTAG users
via New York), so proxying can *worsen* TTFB in specific geographies [^1][^7]. Port scans
against anycast IPs show many ports "open" because the edge serves all customers on shared IPs [^1].

## Restoring the real client IP

The origin's `REMOTE_ADDR` becomes a Cloudflare IP. The real client arrives in:

- `CF-Connecting-IP` — single IP, always set/overwritten at the edge [^1][^5]
- `True-Client-IP` — Enterprise synonym (Akamai compatibility), same value [^1]
- `X-Forwarded-For` — appended chain; equals CF-Connecting-IP when no prior XFF [^1][^5]

Restore with Apache **mod_remoteip** (`RemoteIPHeader CF-Connecting-IP` + the Cloudflare ranges
as `RemoteIPTrustedProxy`; mod_cloudflare is deprecated) or NGINX **ngx_http_realip_module**
(`set_real_ip_from <each CF range>; real_ip_header CF-Connecting-IP;`) [^1][^4].

**Trust boundary (critical):** these headers are trivially spoofable by anyone who can reach
the origin directly. Trust them **only** when the peer IP is verified to be Cloudflare —
ideally combined with Authenticated Origin Pulls mTLS [^5][^6][^9]. A proxy stack
(e.g. Traefik) must have the Cloudflare ranges in its trusted-IPs list or it will honor forged
headers [^9].

**Pseudo-IPv4** maps IPv6 clients to hashed Class-E (240.0.0.0/4) IPv4 addresses for
IPv4-only origin software — either as an added `Cf-Pseudo-IPv4` header or by overwriting
`CF-Connecting-IP`/XFF (real IPv6 preserved in `CF-Connecting-IPv6`) [^1].

## Proxied ports

`verified-as-of: 2026-06-12` — HTTP: 80, 8080, 8880, 2052, 2082, 2086, 2095. HTTPS: 443,
2053, 2083, 2087, 2096, 8443. **Caching is disabled on the ten nonstandard ports** (2052,
2053, 2082, 2083, 2086, 2087, 2095, 2096, 8880, 8443) — i.e. caching stays active on 80, 443,
and 8080; Enterprise can re-enable it per port via Cache Rule. Anything else requires
grey-clouding or Spectrum (arbitrary TCP/UDP, Enterprise) [^1][^8].

## TLS through Cloudflare

A proxied request involves **two independent TLS sessions** — client↔edge and edge↔origin —
which negotiate versions and ciphers independently: a visitor can speak TLS 1.3 to the edge
while Cloudflare↔origin runs TLS 1.2. The zone's "Minimum TLS Version" governs only the client
side [^1][^16].

**Edge certificates:** every active zone gets a free auto-renewed **Universal SSL** DV cert
covering the apex and first-level subdomains only — `dev.www.example.com` is NOT covered (fix:
Advanced Certificate Manager, Total TLS, or custom certs) [^1]. Universal SSL **requires SNI**
clients; Free-plan certs additionally require ECDSA support [^1][^10]. **Total TLS** (needs
ACM + full DNS setup) auto-issues per-hostname 90-day certs for proxied hostnames Universal
doesn't cover; skips hostnames used with Load Balancing, Tunnel, or Spectrum [^1].

**Encryption modes** control only the edge→origin hop (API values `off`, `flexible`, `full`,
`strict`, `origin_pull`):

| Mode | Edge→origin | Risk |
|---|---|---|
| Off | cleartext both hops | everything |
| Flexible | plain HTTP to origin | padlock lies: CF→origin leg is interceptable plaintext (documented ISP content injection via this gap [^13]); classic `ERR_TOO_MANY_REDIRECTS` loop when origin redirects HTTP→HTTPS [^12][^14]; only works on port 443 |
| Full | matches the visitor's protocol (HTTPS visitor → HTTPS to origin, **no cert validation** — accepts self-signed/expired; HTTP visitor → HTTP to origin) | MITM-able on the back hop [^1][^11][^12] |
| Full (strict) | HTTPS validated against a publicly trusted CA or Cloudflare Origin CA | the destination state [^1][^11] |
| Strict (SSL-Only Origin Pull) | always HTTPS to origin regardless of visitor scheme | — |

Contested nuance: Troy Hunt defends Flexible as incremental progress for TLS-incapable origins
[^11]; joepie91's critique documents real interception through it [^13]. Technically undisputed:
Flexible sends plaintext to the origin.

**Origin CA certificates** — free Cloudflare-issued certs (up to 200 SANs, single-level
wildcards, no IP SANs) trusted **only by the Cloudflare edge**: pause Cloudflare or grey-cloud
the record and browsers throw `NET::ERR_CERT_AUTHORITY_INVALID`. Cloudflare currently sends no
expiry notifications for them [^1].

**Authenticated Origin Pulls (AOP)** — the edge presents a client certificate the origin
verifies (mTLS), at three levels: global (shared Cloudflare cert — proves only "came from
Cloudflare's network, possibly another customer's zone"), zone-level, per-hostname (own
uploaded certs). Requires Full/Full (strict); does not function in Off/Flexible. AOP is what
makes `CF-Connecting-IP` actually trustworthy [^1][^9]. Client mTLS cannot pass *through* the
edge — the edge terminates the client handshake, so client cryptographic material never
reaches the origin [^15].

**Errors:** 525 = edge→origin TLS handshake failure under Full/Full (strict) (no cert, 443
closed, no SNI support on origin, cipher mismatch). 526 = certificate validation failure under
Full (strict) (expired/self-signed/name-mismatch/incomplete chain) — dropping to Full "fixes"
526 only by disabling validation [^1][^12].

**Certificate Transparency Monitoring** (opt-in, public beta as of 2026-06) emails on any CA
logging a cert for your domain; expect routine alerts from Cloudflare partner-CA renewals [^1][^17].

## Caching defaults and cf-cache-status

Cloudflare caches by **file extension only** (not MIME type) from a fixed list of ~56
extensions (CSS, JS, JPG, PNG, WOFF2, PDF, ZIP, MP4, AVIF, ZST, APK…; `verified-as-of:
2026-06-12`) plus robots.txt — **HTML and JSON are not cached by default** [^1][^18][^19].
Cacheability additionally requires GET, no `Set-Cookie`, and no
`Cache-Control: private/no-store/no-cache/max-age=0` [^1]. Origin sends `s-maxage` on HTML?
Still ignored until a Cache Rule marks HTML eligible [^19]. Without origin cache headers,
default Edge TTLs by status: 200/206/301 → 120 min; 302/303 → 20 min; 404/410 → 3 min [^1].
Request collapsing (cache lock) dedupes concurrent misses per colo [^1].

`cf-cache-status` values [^1][^18]:

| Value | Meaning |
|---|---|
| HIT | served from cache |
| MISS | not in cache; fetched; cacheable |
| DYNAMIC | not eligible and no rule makes it eligible — the **normal** status for HTML; #1 "Cloudflare isn't caching my site" confusion |
| BYPASS | origin said no-cache/private/max-age=0, or sent Set-Cookie/Authorization |
| EXPIRED | found stale; refetched synchronously |
| STALE | served stale (origin unreachable) |
| UPDATING | served stale while async revalidation refreshes |
| REVALIDATED | synchronous conditional revalidation (INM/IMS) |
| NONE/UNKNOWN | never touched cache: Worker-generated, WAF block, redirects |

`Age` appears only on cache-served responses [^1].

## Cache Rules, TTLs, purge, tiered cache

**Cache Rules** supersede Page Rules. Page Rules: first-match-wins; modern Rules: stackable,
**last-match-wins**, and take precedence over Page Rules [^1][^20]. Migration gotcha: selecting
"Eligible for cache" enables **Cache Everything by default** — the opposite of Page Rules. So
caching HTML = one Cache Rule marking it eligible, plus bypass rules for `/login`, `/cart`,
etc. (cache-everything + Edge TTL override can strip `Set-Cookie` and break sessions) [^1][^18].

**TTL precedence** (`verified-as-of: 2026-06-12`): Edge Cache TTL rules can respect, override,
or bypass origin `Cache-Control`. `s-maxage` targets the edge, `max-age` the browser; Edge TTL
overrides `s-maxage`; Browser Cache TTL overrides `max-age` when higher or absent. Origin
Cache Control strict-adherence is always on for Free/Pro/Business, toggleable Enterprise; with
it on, `max-age=0` revalidates rather than bypasses. Minimum Edge TTL is plan-gated (Free 2 h →
Enterprise 1 s) [^1][^21].

**Purge:** by single URL (instant, preferred), by tag/hostname/prefix, or everything — all
plans, differentiated by rate limits (Free 5 req/min … Enterprise 50 req/s) [^1].

**Tiered Cache:** lower tiers face visitors; only **upper tiers** contact the origin — higher
hit ratio, fewer origin connections. Topologies: Smart (latency-picked single upper tier per
origin), Generic Global, Regional, Custom (Enterprise). Confirm via `CacheTieredFill` [^1].

**Cache Reserve:** R2-backed persistent "ultimate upper tier" (paid: $0.015/GB-mo, $4.50/M
class-A writes, $0.36/M class-B reads; `verified-as-of: 2026-06-12`). Eligibility: cacheable +
TTL ≥ 10 h + Content-Length. 30-day sliding retention. Purge-by-URL purges it instantly;
tag/host/prefix purges only force revalidation. (Smart/Regional Tiered Cache + Cache Reserve
are being re-marketed as "Smart Shield" as of 2026.) [^1]

**Cache keys:** default = scheme + host + URI-with-query (+ Origin header for CORS). Custom
keys (headers, cookies, geo/device/lang) are Enterprise; ignore/sort-query-string and
device-type are all-plan. The scheme is part of the default key — changing SSL mode
(Flexible→Full) busts the cache [^1]. **Vary is ignored** for caching decisions except
`Vary: Accept-Encoding` and opt-in Vary for Images (Accept-based WebP/AVIF) — Accept-based
content negotiation behind Cloudflare can serve the wrong variant unless you build a custom
cache key on Accept (production postmortem: WebP served to non-WebP clients) [^1][^22][^23].

## References

1. developers.cloudflare.com — proxy-status, how-cloudflare-works, cloudflare-ip-addresses,
   network-ports, pseudo-ipv4, restoring-original-visitor-ips, http-headers,
   protect-your-origin-server, ssl/ssl-modes/*, universal-ssl, origin-ca,
   authenticated-origin-pull, total-tls, ct-monitoring, cache/default-cache-behavior,
   cache-control, edge-browser-cache-ttl, cache-rules/*, purge-cache, tiered-cache,
   cache-reserve, cache-keys, vary-for-images (tier: docs; one source)
2. https://pentesting.se/en/blog/cloudflare-origin-bypass — allow-CF-IPs bypassability, origin-IP discovery (tier: blog)
3. https://kbeezie.com — blocking direct-to-origin, realip config (tier: blog)
4. https://laamanen.net/ubuntu-vps-cloudflare-essential-security — UFW CF-only + realip (tier: blog)
5. https://adam-p.ca/blog/2022/03/x-forwarded-for — XFF parsing, CF-Connecting-IP semantics tested (tier: blog)
6. devsec-blog.com — X-Forwarded-For spoofing when proxy is bypassable (tier: blog)
7. community.cloudflare.com #586399, #662386, #608356 + lowendtalk #192516 — ping/traceroute hits PoP; DTAG-via-NYC latency (tier: forum)
8. https://blog.cloudflare.com/cloudflare-now-supporting-more-ports — port list rationale (tier: blog)
9. josemanuelortega.me — CF-Connecting-IP only trustworthy with AOP; Traefik trustedIPs pitfall (tier: blog)
10. https://blog.cloudflare.com/introducing-universal-ssl + community #98078, #110502 — SNI/ECDSA requirements (tier: blog/forum)
11. https://www.troyhunt.com/cloudflare-ssl-and-unhealthy-security-absolutism — mode ladder + Flexible-as-progress argument (tier: blog)
12. panelica.com + stackharbor.com SSL-mode guides — redirect loops, 525 causes (tier: blog)
13. http://cryto.net/~joepie91/blog/2016/07/14/cloudflare-we-have-a-problem — TPB/Airtel interception via Flexible (tier: blog)
14. community.cloudflare.com/t/63531 — community deprecation of Flexible (tier: forum)
15. hetneo.link — why client mTLS can't traverse a terminating edge (tier: blog)
16. community.cloudflare.com/t/325172 — TLS 1.3 client-side with TLS 1.2 origin-side, empirical (tier: forum)
17. https://blog.cloudflare.com/introducing-certificate-transparency-monitoring — CT crawler (tier: blog)
18. community.cloudflare.com #477213, #593111 — DYNAMIC semantics, last-match-wins (tier: forum)
19. agentcookbooks.com — measured: s-maxage ignored for HTML until Cache Rule opt-in (tier: blog)
20. https://blog.cloudflare.com/future-of-page-rules + /introducing-cache-rules + /cache-rules-go-ga — Page Rules EOL, precedence (tier: blog)
21. https://www.jonoalderson.com/performance/http-caching — independent TTL-precedence analysis (tier: blog)
22. https://simonwillison.net/2023/Nov/20/ — Vary-ignored content-negotiation hazard (tier: blog)
23. github.com/umbraco/Umbraco.Cloud.Issues #815 — WebP-to-non-WebP-clients postmortem (tier: forum)
