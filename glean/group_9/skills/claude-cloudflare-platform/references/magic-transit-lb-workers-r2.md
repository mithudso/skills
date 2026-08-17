# Magic Transit/WAN, Load Balancing & Argo, Workers & R2 (infrastructure angle)

Provenance: /dr deep-research run 2026-06-12 (skill: cloudflare-platform). Volatile values
(pricing, limits, vendor performance claims) `verified-as-of: 2026-06-12`.

## Contents

- [Magic Transit](#magic-transit)
- [Tunnel health checks and the MTU/MSS pitfall](#tunnel-health-checks-and-the-mtumss-pitfall)
- [Magic WAN](#magic-wan)
- [Load Balancing](#load-balancing)
- [Argo Smart Routing](#argo-smart-routing)
- [Workers from the infrastructure angle](#workers-from-the-infrastructure-angle)
- [R2 object storage](#r2-object-storage)
- [References](#references)

## Magic Transit

Network-layer (L3) DDoS protection + acceleration for whole IP prefixes: the customer brings
their own public IP space (BYOIP) and **Cloudflare BGP-announces those prefixes from its
anycast network**, so all Internet traffic for the customer's network ingests at the nearest
PoP, where DDoS mitigation and Magic Firewall rules apply before clean delivery
[^1][^4][^5][^6]. Magic Transit exists because the orange-cloud only covers L7/HTTP — L3/4
protection for non-HTTP services needs prefix-level interception [^5].

- **Delivery off-ramps:** anycast GRE tunnels, IPsec tunnels, or CNI (Cloudflare Network
  Interconnect — physical/virtual cross-connect) [^1][^6].
- **Anycast GRE:** GRE is stateless and endpoints bind to anycast IPs rather than devices, so
  one configured tunnel effectively terminates at *every* Cloudflare DC — any PoP can
  encapsulate and forward independently ("talking to us from over 500 colos" — practitioner) [^1][^7].
- **Deployment shapes:** **ingress-only** (return traffic egresses directly via Direct Server
  Return) or **ingress + egress** (policy-based routing/default routes back into the tunnels) [^1][^6].
- **Minimum prefix /24** — longer prefixes aren't globally routable; customers without a /24
  can lease Cloudflare IPs [^1].
- Onboarding distinguishes IP prefixes (LOA/IRR/RPKI ROA alignment) from BGP prefixes (what's
  announced and billed); Cloudflare prepends AS13335 ahead of the customer ASN; 0–3 prepends
  configurable (API or communities 13335:50101–50103) [^1].
- **On Demand:** advertise/withdraw prefixes dynamically (API/dashboard, BGP sessions to three
  redundant route reflectors, or traffic-threshold automation) — scrubbing only under attack.
  Use one control plane per prefix; follow the documented multi-step withdrawal to avoid
  "BGP zombie" blackholes [^1].
- Pricing is opaque/enterprise-only (tentative — practitioner reports) [^7].

## Tunnel health checks and the MTU/MSS pitfall

Every Cloudflare DC sends ICMP probes encapsulated in the tunnel protocol (rate tunable).
Probes are **unidirectional** (reply returns outside the tunnel) or **bidirectional** (reply
re-enters the tunnel, using the /31 or /30 interface addresses) [^1]. Tunnel states: healthy /
degraded (≥0.1% probe failures in 5 min) / down (all probes in the last second failing).
Instead of route withdrawal, **priority penalties** (+500,000 / +1,000,000) preserve
last-resort paths; recovery is deliberately slow (hysteresis, up to 30 min) to prevent
flapping [^1]. Separate endpoint health checks observe beyond the tunnel router but don't
steer [^1].

**The classic post-onboarding breakage is MTU/MSS:** GRE adds 24 bytes → tunnel MTU **1476**,
so IPv4 TCP MSS must be clamped to **1436**. Ingress-only/DSR customers must clamp on
edge-router transit ports. Unclamped traffic yields dropped packets and PMTUD-blackhole-style
"hard-to-debug connectivity issues" — Cloudflare's own troubleshooting path is largely MSS
forensics (including asking Cloudflare to clear the DF bit) [^1][^8][^9].

## Magic WAN

Same plumbing as WAN-as-a-service replacing MPLS: branches, DCs, VPCs, and roaming users
connect via anycast GRE/IPsec/CNI, the **Magic WAN Connector** appliance (zero-touch on-ramp,
>1 Gbps via SFP+), WARP device clients, or Cloudflare Tunnel; **Magic Firewall** enforces
network-layer rules at every PoP [^1][^2][^3]. Rebrand note: docs now say "**Cloudflare WAN
(formerly Magic WAN)**" (`verified-as-of: 2026-06-12`) [^1].

## Load Balancing

Hierarchy: load balancer (hostname) → **pools** → **endpoints/origins**, with **monitors**
attached to pools [^1].

- **Monitors:** interval, timeout, retries, expected status codes, expected body substring
  (must appear in the first 10 KB), follow-redirects, `consecutive_up`/`consecutive_down`
  anti-flapping, configurable probe **regions**, and `probe_zone` (simulate zone) so probes
  inherit zone features like Authenticated Origin Pulls or Argo [^1].
- **Pool health:** healthy/degraded/critical vs a health threshold; mandatory fallback pool;
  transitions fire notifications and appear in health-check event logs [^1].
- **Global steering:** `off` (failover by pool order), `random` (weighted), `geo`
  (region/country/PoP maps), `dynamic_latency` (RTT from health probes), `proximity`
  (GeoIP/lat-long), `least_outstanding_requests`, `least_connections` [^1].
- **Origin steering** within a pool: random, hash (CF-Connecting-IP), least_outstanding_requests,
  least_connections, per-origin weights, load shedding [^1].
- **Session affinity:** cookie (`cookie`/`ip_cookie`, default TTL 23 h) or header-based
  (default 1800 s); overrides steering/weights; **requires proxied mode** [^1].
- **Proxied (L7) vs DNS-only:** proxied hides origins, fast failover, integrates with
  cache/WAF/Workers; DNS-only just returns healthy IPs, is hostage to resolver TTL caching,
  has no session affinity, and exposes origins. L4 balancing only via Spectrum [^1].
- **Health-check false positives** (Cloudflare's own FAQ): origin firewall blocking probe IPs,
  rate limiting, redirects without follow-redirects, body string past 10 KB, missing Host
  header, missing probe_zone with mTLS/Argo. Real bug: monitors not sending the configured
  Host header to cloudflared-fronted origins (TCP-check workaround) [^1][^17].

## Argo Smart Routing

Routes dynamic/uncached (cache-miss) traffic edge↔origin over the fastest measured path across
Cloudflare's network — built from real-time latency/loss telemetry of Cloudflare's own
traffic; runs **on top of BGP**, not as a replacement [^10][^12]. The headline "~30% faster
TTFB" is a **vendor claim** [^12]. **Argo for Packets** extends to L3/L4 for Magic
Transit/WAN, reusing tunnel health probes as a one-way latency mesh [^1][^10][^11].

- **Pricing** (`verified-as-of: 2026-06-12`): $5/month per domain + $0.10/GB after the first
  GB — billed on traffic between Cloudflare and visitors in both directions, i.e. **including
  cached traffic Argo never accelerated** (long-standing practitioner critique; qualified)
  [^13][^16][^19].
- **Disconfirming:** near-zero benefit when origins sit close to users or on premium transit;
  can occasionally *add* latency; gains concentrate on long-haul/poorly-routed geographies.
  Run your own A/B tests — dashboard analytics often show "not enough requests" or 0%
  optimized [^13][^14][^15].
- **Naming archaeology:** "Argo Tiered Cache" was split off and made free (~Sept 2021) as
  plain **Tiered Cache** (qualified) [^1][^13]; "**Argo Tunnel**" was renamed **Cloudflare
  Tunnel** in April 2021 — only Smart Routing kept the Argo name [^18].

## Workers from the infrastructure angle

(The JS runtime — APIs, compat flags, wrangler — is owned by
`javascript-runtimes-deno-bun-edge`; this section covers where Workers sit in the network.)

Workers run on every machine in every PoP; the runtime is **V8 isolates, not containers** —
thousands per process, ~100× faster start than container cold-start, an order of magnitude
less memory; isolates can be evicted, so global state is non-durable [^1].

**Traffic binding:** (`verified-as-of: 2026-06-12` for limits)

- `workers.dev` subdomain — auto-provisioned, explicitly non-production.
- **Zone routes** — URL patterns where the Worker fronts an origin; cannot be the target of a
  same-zone `fetch()`. 1,000 routes/zone; 1,000 routed zones per Worker.
- **Custom domains** — the Worker *is* the origin; same-zone fetch works. 100/zone [^1].

On a matched route the Worker executes in the request path **ahead of normal cache lookup**
and reaches origins via `fetch()` subrequests whose responses flow back through caching logic
(explicit control via the Cache API) — qualified, assembled from Workers docs [^1].
**Subrequests:** 50/invocation Free, **10,000/invocation Paid (raisable)** — higher than the
older widely-cited 1,000 figure; max 6 simultaneous connections awaiting response headers; CPU
10 ms Free / default 30 s (up to 5 min) Paid; 128 MB per isolate [^1].

**Smart Placement** (opt-in, `placement.mode = "smart"`): moves execution near the *backend*
when the Worker averages >1 subrequest to a backend — compares measured duration across
candidate PoPs, ~15 min analysis, keeps 1% of requests un-placed as a control group; Jan 2026
added **placement hints** (`placement.region` with cloud region IDs, or host/hostname probes —
unsuitable for anycast backends) [^1][^20]. **Cron Triggers:** 5/account Free, 250 Paid,
15-min wall-clock cap — scheduled invocations decoupled from any eyeball PoP [^1].

**Durable Objects** (one line): a single-instance, globally-addressable coordination primitive
— strongly consistent storage + serialized execution for state that must live in exactly one
place [^1].

## R2 object storage

S3-API-compatible object storage whose architectural pitch is **zero egress fees** — egress
via Workers API, S3 API, and r2.dev is free, vs S3's $0.05–0.09/GB egress ladder. That is what
makes R2-as-CDN-origin economics work [^1][^21].

- **Pricing** (`verified-as-of: 2026-06-12`): Standard $0.015/GB-mo; Infrequent Access
  $0.01/GB-mo (+$0.01/GB retrieval, 30-day minimum); Class A ops $4.50/M ($9/M IA), Class B
  $0.36/M ($0.90/M IA); free tier 10 GB + 1M A + 10M B [^1].
- **"Zero egress" ≠ free serving:** every GET is a billed Class B op unless absorbed by cache
  or free tier [^22]. Scenario analysis: R2 ~99% cheaper for hot image-serving, but **S3 ~73%
  cheaper for cold archival** (storage-class granularity/Intelligent-Tiering) [^21]. S3 keeps
  a deeper feature surface [^21][^23].
- **Placement:** automatic (near creator); best-effort **location hints**
  (wnam/enam/weur/eeur/apac/oc, honored only at first creation of a bucket name); hard
  **jurisdictional restrictions** (`eu`, `fedramp`) — immutable, with jurisdiction-specific S3
  endpoints (`<account>.eu.r2.cloudflarestorage.com`) [^1].
- **Public exposure:** `r2.dev` is dev-only (rate-limited, throttled, no cache/WAF). Binding a
  **custom domain** puts the bucket behind the zone's orange-cloud pipeline — Cache (Smart
  Tiered Cache recommended; "Cache Everything" rule needed since only some extensions cache by
  default), WAF, Access, bot management. This is the canonical static-site/CDN-offload pattern [^1].
- **Limits:** 5 TiB max object (effectively 4.995 TiB); single upload/part ≤ ~4.995 GiB;
  10,000 multipart parts; 1 write/sec per key; 1M buckets/account [^1].
- **Event notifications:** object-create/delete (incl. lifecycle) → Cloudflare Queues, with
  prefix/suffix filters; ≤100 rules/bucket [^1].

## References

1. developers.cloudflare.com — magic-transit (about, tunnel-health-checks, mtu-mss,
   gre-ipsec-tunnels, traffic-steering, advertise-prefixes), magic-wan / cloudflare-wan,
   load-balancing (proxy-modes, health-details, FAQ, origin-level-steering),
   argo-smart-routing, workers (how-workers-works, limits, routing, placement, changelogs),
   r2 (pricing, limits, public-buckets, data-location, event-notifications) (tier: docs; one source)
2. https://blog.cloudflare.com/magic-wan-firewall/ — Magic WAN/Firewall architecture (tier: blog)
3. https://blog.cloudflare.com/magic-wan-connector/ — Connector appliance (tier: blog)
4. https://blog.cloudflare.com/flow-based-monitoring-for-magic-transit/ — BGP-attraction model (tier: blog)
5. https://blog.kingsmill.io/2022/07/setting-up-cloudflare-magic-transit/ — practitioner GRE/DSR setup (tier: blog)
6. https://www.securityscientist.net/blog/12-questions-and-answers-about-cloudflare-magic-transit/ (tier: blog)
7. https://www.reddit.com/r/networking/comments/1b0r6iw/ — anycast-GRE observations; pricing opacity (tier: forum)
8. https://www.travteks.com/blog/mtu-calculator-launch/ — stacked-tunnel MTU/MSS failures (tier: blog)
9. https://networkengineering.stackexchange.com/questions/80787/ — GRE 1476 MTU problem class (tier: forum)
10. https://blog.cloudflare.com/argo-v2/ — Argo for Packets mechanism (tier: blog)
11. https://blog.cloudflare.com/turbo-charge-gaming-and-streaming-with-argo-for-udp/ (tier: blog)
12. cloudflare.com Argo product page + 2017 press release — 30%/35% vendor claims (tier: vendor marketing)
13. https://community.cloudflare.com/t/is-argo-worth-it-for-u-s-only-traffic/246708 + related threads (tier: forum)
14. https://lowendtalk.com/discussion/194316/ — negative Argo cost/benefit reports (tier: forum)
15. https://wpjohnny.com/is-cloudflare-argo-routing-service-worth-the-cost/ (tier: blog)
16. https://news.ycombinator.com/item?id=14367378 — Argo billed-on-all-traffic critique (tier: forum)
17. https://github.com/cloudflare/cloudflared/issues/1261 — LB monitor Host-header bug (tier: forum)
18. https://blog.cloudflare.com/tunnel-for-everyone/ — Argo Tunnel → Cloudflare Tunnel rename (tier: blog)
19. https://www.fdaytalk.com/how-cloudflare-argo-pricing-works/ + cloudways.com — pricing confirmation (tier: blog)
20. https://blog.cloudflare.com/announcing-workers-smart-placement/ — placement algorithm, 1% control (tier: blog)
21. https://www.vantage.sh/blog/cloudflare-r2-aws-s3-comparison — R2-vs-S3 scenario costs (tier: analyst blog)
22. https://www.reddit.com/r/CloudFlare/comments/1ic51x1/ — Class-B cost reality (tier: forum)
23. https://www.digitalapplied.com/blog/cloudflare-r2-vs-aws-s3-comparison — feature-gap framing (tier: blog)
