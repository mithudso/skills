# Cloudflare DNS (authoritative + 1.1.1.1) and debugging through the proxy

Provenance: /dr deep-research run 2026-06-12 (skill: cloudflare-platform). Volatile values
(timeouts, plan gates, TTLs) `verified-as-of: 2026-06-12`.

## Contents

- [Authoritative DNS](#authoritative-dns)
- [The 1.1.1.1 public resolver](#the-1111-public-resolver)
- [No ECS: the geo-routing consequence](#no-ecs-the-geo-routing-consequence)
- [July 14, 2025 resolver outage](#july-14-2025-resolver-outage)
- [The 52x error family](#the-52x-error-family)
- [1xxx body-level codes](#1xxx-body-level-codes)
- [Debugging workflow through the proxy](#debugging-workflow-through-the-proxy)
- [References](#references)

## Authoritative DNS

Cloudflare authoritative DNS runs on the global anycast network — the same NS IPs announced
from every PoP [^1][^4]. Two onboarding modes:

- **Full (primary) setup** — repoint registrar NS records; Cloudflare becomes authoritative.
- **Partial (CNAME) setup** — keep your existing provider; point individual hostnames at
  `{hostname}.cdn.cloudflare.net`. Business/Enterprise only; no DNS-infrastructure DDoS
  protection; apex proxying only if the external provider supports CNAME flattening/ALIAS
  (RFC 1912 forbids apex CNAMEs) [^1][^7].

**Proxy status is the core switch.** Only A/AAAA/CNAME can be proxied; MX/TXT/SRV etc. are
always DNS-only. Proxied records answer with anycast edge IPs; DNS-only returns the real
origin IP [^1][^7][^21]. **Why mail must stay DNS-only:** an MX target resolving to edge IPs
sends SMTP to a proxy that only speaks HTTP/S — mail breaks. The trade-off is real: the
grey-clouded MX A-record leaks the origin IP (the dashboard warns about exactly this) [^1][^7].
Subtleties: if any A/AAAA on a name is proxied, *all* records on that name are treated as
proxied; a proxied hop anywhere in a CNAME chain proxies the whole name [^1]. Proxied records
get a forced **Auto TTL = 300 s** (not editable) [^1][^7].

**CNAME flattening** resolves the chain server-side and returns final A/AAAA — mandatory at
the apex, default for proxied CNAMEs. Flatten-all can break third-party CNAME domain
verification; CNAME into another Cloudflare account triggers **error 1014** [^1][^7].
**Wildcards** are proxyable on all plans (since May 2022); Universal SSL covers only one label
deep [^1][^4]. **DNSSEC** is one-click zone signing (add DS at registrar), with multi-signer
DNSSEC (RFC 8901) and secondary-mode options [^1][^7]. **Foundation DNS** adds advanced
nameservers on `foundationdns.com/.net` using **two separate anycast groups**; custom (vanity)
nameservers brand NS on your own domain [^1][^3].

## The 1.1.1.1 public resolver

Anycast resolver on 1.1.1.1/1.0.0.1 (+ 2606:4700:4700::1111/1001), launched with APNIC.
Encrypted transports (`verified-as-of: 2026-06-12`): **DoT** (port 853, `one.one.one.one`),
**DoH** (`https://cloudflare-dns.com/dns-query`, wireformat + JSON, HTTP/2 and HTTP/3), and
**ODoH** — **not RFC 9250 DoQ** (a 2026 vendor article claiming Cloudflare DoQ GA is
contradicted by Cloudflare's own docs; DoH-over-HTTP/3 ≠ DoQ; treat as contested/likely wrong)
[^1][^7][^27]. Filtering variants: **1.1.1.2/1.0.0.2** (malware; DoH
`security.cloudflare-dns.com/dns-query`) and **1.1.1.3/1.0.0.3** (malware+adult;
`family.cloudflare-dns.com`) [^1][^5].

Privacy: source IPs truncated and deleted within 25 h, shared only with APNIC, independently
audited (KPMG) [^1][^6][^9]. The consumer WARP app tunnels device DNS through 1.1.1.1
(`warp=on` in `/cdn-cgi/trace`) [^1][^20].

## No ECS: the geo-routing consequence

1.1.1.1 **does not send EDNS Client Subnet upstream** [^1][^6]. Consequence: ECS-dependent
CDNs and traffic managers (Akamai, Azure Traffic Manager, MaxMind-based geosteering) locate
1.1.1.1 users by *resolver-PoP* location, not client subnet. Measured effects: 2× slower
Akamai fetches; Australia→US misrouting; ISP-embedded caches (Netflix/Akamai in-ISP) become
unselectable [^8][^9][^10][^11][^12]. Counterpoint: anycast density keeps continent/country
accuracy mostly fine — the loss is at city/carrier/ISP-embedded-cache granularity. Both sides
are partially right [^8][^9][^10]. Debugging: compare `dig @1.1.1.1` vs `dig @8.8.8.8` vs the
zone's authoritative servers.

## July 14, 2025 resolver outage

62-minute global outage (21:52–22:54 UTC). Root cause: a June 6 config error attached resolver
prefixes to a non-production Data Localization Suite topology; a July 14 change refreshed
global config and **withdrew the resolver's BGP prefixes from all production PoPs**. The
"BGP hijack" narrative was debunked — Tata AS4755's 1.1.1.0/24 announcement was a side-effect
exposed by Cloudflare's own withdrawal (RPKI ROV limited its spread), not the cause.
Remediation: retire legacy static-IP topology systems, staged health-gated rollouts
[^2][^13][^14][^15][^16][^17][^18]. Operational lesson: configure resolver diversity
(e.g. 1.1.1.1 + a second provider), not two addresses of the same provider [^17].

## The 52x error family

All 52x are **Cloudflare-generated** and describe the edge→origin leg
(`verified-as-of: 2026-06-12` for timeout values):

| Code | Meaning | First checks |
|---|---|---|
| 520 | Catch-all: origin returned unexpected/empty/protocol-violating response | crashed app; security plugin blocking CF IPs; malformed/oversized headers [^1][^22][^23] |
| 521 | Origin actively refused TCP | web server down; port blocked; origin firewall rejecting CF ranges [^1][^22] |
| 522 | TCP timeout: no SYN+ACK within **19 s**, or no ACK of the request within **90 s** post-connect | #1 cause: CF ranges blocked/rate-limited in iptables/.htaccess; overloaded origin; stale origin IP [^1][^7][^24][^25] |
| 523 | Origin unreachable (routing failure / wrong origin IP) | [^1][^22] |
| 524 | TCP connect OK but no HTTP response within the **Proxy Read Timeout — default 120 s** (older sources say 100 s; docs now say 120, with ~1 s skew from the Pingora migration). Enterprise: raisable to 6,000 s (Cache Rule "Proxy Read Timeout" / zone API). Free/Pro/Business: cannot raise — grey-cloud the slow endpoint or convert to status-polling. Separate non-adjustable 30 s Proxy Write Timeout | [^1][^7][^19][^21] |
| 525 | TLS handshake to origin failed (Full/Full strict) | see proxy-tls-caching.md [^1][^22] |
| 526 | Origin cert invalid under Full (strict) | see proxy-tls-caching.md [^1][^22] |
| 530 | Cloudflare can't resolve the origin hostname; HTML body carries the specific 1xxx code (commonly 1016) | [^1] |

**52x is not always the origin's fault** — documented 522s with healthy origins (peering/edge
issues, partial-region failures). Check cloudflarestatus.com and test per-colo before blaming
the origin [^7][^21][^25].

## 1xxx body-level codes

Families [^1][^7]: config/DNS (1000 DNS points to prohibited IP, 1001, 1004, 1014 cross-account
CNAME, 1016 origin DNS error), access enforcement (1003 direct-IP access, 1005 ASN ban,
1006/1007/1008/1106 IP ban, 1009 country ban, 1010, 1012, 1020 access rules), 1013 SNI/Host
mismatch, **1015 rate limited**, rewrite errors (1035–1041), 1033 (Tunnel error),
**Workers: 1101 = Worker threw an exception, 1102 = Worker exceeded CPU/resource limits**,
1200 cache connection limit.

## Debugging workflow through the proxy

1. **Did the response traverse Cloudflare?** `server: cloudflare` and `cf-cache-status`
   headers; Cloudflare-branded error pages show the error code with the **Ray ID at the
   bottom** — a plain/host-branded error page means the origin generated it [^19][^22][^23].
2. **Which PoP?** `Cf-Ray: 230b030023ae2822-SJC` = ray ID + **IATA code of the serving colo**.
   The ray ID is also forwarded to the origin as a request header — log it (nginx
   `$http_cf_ray`) to correlate edge↔origin. With Argo/Tiered Cache the colo reflects the DC
   contacting the origin, not the ingress PoP. The ray ID is what Cloudflare support asks for
   [^1][^19][^20].
3. **Client-side ground truth:** every proxied zone exposes **`/cdn-cgi/trace`** (fields
   `fl,h,ip,ts,visit_scheme,uag,colo,http,loc,tls,sni,warp,gateway,kex`) — `colo=` identifies
   your PoP; unavailable on grey-clouded hosts [^1][^20][^28].
4. **Isolate the origin:** `curl -I --resolve example.com:443:ORIGIN_IP https://example.com/`
   (or hosts-file edit) bypasses Cloudflare. Direct-OK + proxied-52x ⇒ the CF→origin path
   (usually the origin firewall) is at fault [^22][^23][^25][^26].
5. **Coarse switch:** "Pause Cloudflare on Site" (Overview → Advanced Actions) keeps DNS
   authoritative but routes traffic straight to origin zone-wide; grey-clouding does it per
   record [^1][^23].
6. **Platform-side:** cloudflarestatus.com (per-PoP status, re-routed colos) and Cloudflare
   Radar [^22][^23].
7. **Origin-side evidence:** origins log Cloudflare IPs as peers — restore real client IPs
   (CF-Connecting-IP; see proxy-tls-caching.md) before reading origin access logs [^1][^25].
8. **499 in Cloudflare analytics** = client closed the connection before Cloudflare responded
   (nginx-lineage code) — a *client* disconnect, not an origin failure; an impatient client
   logs 499 where a patient one would have hit 522/524. Since Jan 2026, HTTP/3 cancellations
   log 499 immediately; Cloudflare advises excluding them from error-rate SLOs [^1].

**Origin secrecy caveat:** any DNS-only record pointing at the same origin (MX target, direct
subdomain) leaks the IP and re-enables direct attack/bypass — `--resolve` against leaked IPs
is the standard origin-exposure test [^1][^7][^25].

## References

1. developers.cloudflare.com — dns/proxy-status, cname-flattening, partial-setup,
   foundation-dns, dnssec; 1.1.1.1 setup/encryption/privacy; support 5xx + per-error 520–530,
   1xxx index, error-499; fundamentals http-headers, cdn-cgi-endpoint, pause-cloudflare;
   workers errors (tier: docs; one source)
2. https://blog.cloudflare.com/cloudflare-1-1-1-1-incident-on-july-14-2025/ — outage postmortem (tier: blog)
3. https://blog.cloudflare.com/foundation-dns-launch/ — two anycast groups (tier: blog)
4. https://blog.cloudflare.com/wildcard-proxy-for-everyone/ — wildcard proxying, network scale (tier: blog)
5. https://blog.cloudflare.com/introducing-1-1-1-1-for-families/ — 1.1.1.2/1.1.1.3 (tier: blog)
6. blog.cloudflare.com privacy-examination posts + KPMG audit PDF — 25 h deletion, no-ECS commitments (tier: blog/paper)
7. community.cloudflare.com #17424 (ECS), #228623 (MX), #316628 (524), #393404 (522), #516358 (DoQ), #44273 (526), #640838 (1102) (tier: forum)
8. https://sajalkayan.com/post/cloudflare-1dot1dot1dot1.html — measured Akamai 2× slowdown (tier: blog)
9. https://blog.apnic.net/2024/07/23/privacy-and-dns-client-subnet/ — ECS privacy vs accuracy (tier: blog)
10. dnswiz.app/blog/ecs-lies — ECS impact granularity breakdown (tier: blog)
11. akamai.com — "A Look at the ECS Behavior of DNS Resolvers" (tier: paper)
12. medium.com/nextdns — ISP-embedded cache dependence on ECS (tier: blog)
13. https://www.kentik.com/blog/cloudflares-dns-downtime-why-bgp-hijacks-were-never-to-blame/ (tier: blog)
14. https://www.thousandeyes.com/blog/cloudflare-outage-analysis-july-14-2025 (tier: blog)
15. bleepingcomputer.com — outage not attack/hijack (tier: news)
16. theregister.com — outage timeline (tier: news)
17. pulse.internetsociety.org — resolver-diversity lessons (tier: blog)
18. https://anuragbhatia.com/post/2025/07/cloudflare-dns-outage/ — AS4755 suppressed-route evidence (tier: blog)
19. http.dev/cf-ray and http.dev/524 — header anatomy, 120 s default (tier: reference)
20. bornoe.org + blog.tsinbei.com — /cdn-cgi/trace field walkthroughs (tier: blog)
21. reddit.com/r/CloudFlare — proxied vs DNS-only, 522/524 practitioner counterpoints (tier: forum)
22. my.ultrawebhosting.com KB — error-page branding heuristic, per-code fixes (tier: kb)
23. contabo.com blog — curl debugging, pause-proxy workflow (tier: blog)
24. cyberoptik.net glossary — 19 s SYN timing corroboration (tier: blog)
25. bigiron.cc "Cloudflare 522 but the Origin Is Healthy" + herish.me — --resolve diagnosis, origin-exposure testing (tier: blog)
26. everything.curl.dev + support.cpanel.net — --resolve/Host-header mechanics (tier: docs)
27. expressvpn.com/blog/dns-over-quic — DoQ GA claim (contested) (tier: blog)
28. dev.to/richardkazuomiller — colo observation, plan-based PoP routing anecdote (tier: blog)
