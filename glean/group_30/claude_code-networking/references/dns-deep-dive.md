<!-- hub-reference-banner -->
> **Reference file — part of the `networking` hub.** Created directly as a hub reference by /dr deep
> research (2026-06-12) — not a folded standalone skill. Sibling topics in this family live as
> reference files under the `networking` hub — **not** standalone skills. Ignore any "use the X skill"
> pointers that name a bare sibling; load that topic's `references/<name>.md` from the owning hub
> (see the hub's routing table).

---

---
name: dns-deep-dive
title: DNS Deep Dive — Resolution Path, Records, TTLs, DNSSEC, Encrypted DNS, Resolvers & Debugging
description: >-
  Protocol-level DNS expert — resolution path & delegation, record types
  (A/AAAA/CNAME/SRV/TXT/SOA/CAA/MX/SVCB-HTTPS), TTLs & negative caching, DNSSEC, encrypted DNS
  (DoT/DoH/DoQ, ECH), split-horizon DNS, EDNS(0), resolver stacks, dig/delv/drill/kdig. TRIGGER:
  name resolves differently here vs there; NXDOMAIN/NODATA/SERVFAIL/REFUSED triage;
  delegation/glue/lame-delegation problems; DNSSEC bogus or expired RRSIG; CNAME at apex; negative
  caching bit me during a migration; split DNS leaking; DoH bypassing internal DNS;
  EDNS/truncation/fragmentation (Flag Day 1232); systemd-resolved, macOS mDNSResponder,
  glibc-vs-musl, Kubernetes ndots:5 quirks; dig +trace caveats. SKIP: MongoDB/Atlas SRV seedlist
  specifics → mongodb-expert (references/mongodb-connection-string.md) / mongodb-atlas-expert
  (references/mongodb-aws-networking.md); host-level net diagnostics (ip/ss/tcpdump) →
  devops-linux-admin; Cloudflare product behavior → cloudflare-platform.
version: 1.0.0
updated: 2026-06-12
category: networking
whenToUse:
  - "Why does this name resolve differently on my laptop vs the server vs inside the pod?"
  - "Triage NXDOMAIN vs NODATA vs SERVFAIL vs REFUSED, or a DNSSEC bogus / expired-RRSIG SERVFAIL"
  - "A migration went sideways — negative caching, resolver TTL clamps, or stale delegation/glue"
  - "Design split-horizon/private DNS (.internal, home.arpa, BIND views) or an encrypted-DNS (DoH/DoT/DoQ) rollout"
  - "Interpret dig/delv/drill/kdig output, dig +trace caveats, or Kubernetes ndots:5 query amplification"
keywords:
  - DNS resolution
  - recursive resolver
  - authoritative nameserver
  - delegation and glue records
  - DNS record types
  - SOA serial
  - CAA record
  - SVCB HTTPS record
  - negative caching
  - DNSSEC validation
  - NSEC3
  - DNS over HTTPS
  - DNS over TLS
  - DNS over QUIC
  - Encrypted ClientHello
  - split-horizon DNS
  - EDNS0
  - systemd-resolved
  - mDNSResponder
  - ndots
  - dig delv drill
tags: [networking, dns, dnssec, encrypted-dns, resolvers, debugging, dr-generated]
---

# DNS Deep Dive — Resolution, Records, TTLs, DNSSEC, Encrypted DNS, Resolvers & Debugging

Volatile claims in this reference (deployment status, defaults, adoption numbers, standards status)
are stamped at header level: **verified-as-of: 2026-06-12** unless individually noted.

## Contents

1. [Resolution path & delegation](#1-resolution-path--delegation)
2. [Record types](#2-record-types)
3. [TTLs, caching & negative caching](#3-ttls-caching--negative-caching)
4. [DNSSEC](#4-dnssec)
5. [Encrypted DNS — DoT/DoH/DoQ, discovery, ODoH, ECH](#5-encrypted-dns)
6. [Split-horizon & private DNS](#6-split-horizon--private-dns)
7. [EDNS(0)](#7-edns0)
8. [Resolver behavior differences — systemd-resolved, macOS, glibc/musl](#8-resolver-behavior-differences)
9. [DNS debugging — dig, delv, drill, kdig](#9-dns-debugging--dig-delv-drill-kdig)
10. [Cross-references](#cross-references)
11. [References](#references)

## 1. Resolution path & delegation

**Roles** (RFC 9499 terminology[^1]): a **stub resolver** is the minimal client library that sends a
query with RD=1 and depends on a recursive resolver; a **recursive (full-service) resolver** acts in
recursive mode with a cache, performing the iterative walk itself (sending RD=0 upstream); an
**authoritative server** answers for zones it holds (AA bit set) and returns **referrals** for
delegated children. A **delegation** exists when the parent zone carries an NS RRset for the child
origin; the boundary is the **zone cut**, and the child's top node is the **apex** (owner of SOA + NS).

**The walk**: stub → recursive (cache check) → root (referral to TLD) → TLD (referral to zone's NS,
with glue if needed) → authoritative (AA answer) → cached and returned. Caching short-circuits any
prefix of this path — most queries never reach the root.

- **Root servers**: 13 named identities (a–m.root-servers.net), each an anycast IPv4+IPv6 address,
  operated by 12 organizations, with hundreds-to-thousands of anycast instances behind them[^2]. The
  "13" comes from fitting all root NS + addresses in a classic 512-byte unfragmented UDP response —
  anycast grew around the limit rather than replacing it. Any specific instance count is a snapshot;
  don't quote one as stable.
- **Priming (RFC 8109)**[^3]: at startup a resolver sends QNAME=".", QTYPE=NS to a root-hints address
  and replaces the possibly stale hints file with the live root NS RRset + addresses. Don't
  hand-maintain root hints.
- **Glue records**: address RRs the *parent* serves for nameservers whose names live *inside* the
  delegated child zone ("in-domain" — RFC 9499 deprecates "in-bailiwick/out-of-bailiwick" as
  historic[^1]). Glue breaks the circular dependency (ns1.example.com serving example.com). Changing
  a glued nameserver's IP must happen **at the registrar/parent**, not just in the zone file. For
  out-of-zone nameservers (ns1.dnsprovider.net), the parent does not serve glue; resolvers rightly
  distrust out-of-zone additional data (cache-poisoning vector).
- **QNAME minimization (RFC 9156, obsoletes 7816)**[^4]: resolvers send only the labels needed to
  each ancestor server (QTYPE A/AAAA recommended, not NS). Costs up to ~26% more lookups; breaks
  against servers that mis-answer NXDOMAIN for **empty non-terminals** (some CDNs) — BIND's default
  `relaxed` mode falls back to the full QNAME on failure, `strict` does not. Default-on in Unbound
  (since 1.7.2), Knot Resolver, and BIND (relaxed).
- **Anycast (RFC 4786)**: the same service IP advertised from many sites; gives DNS lower latency,
  referral-size-free horizontal scale, and DDoS localization.

**Failure modes**: **lame delegation** — parent lists an NS that isn't authoritative (test each
delegated NS directly with `+norecurse` and check AA); **parent/child NS mismatch** — RFC 1034
expects the sets consistent; large-scale studies (CAIDA) found millions of inconsistent domains
causing skewed load and dead-server queries; **missing glue** for in-domain NS → hard deadlock.
Authoritative servers should not offer recursion (open-resolver abuse).

## 2. Record types

| Type | Core semantics | Sharp edges |
|---|---|---|
| **A / AAAA** | 32-bit IPv4 / 128-bit IPv6 address (RFC 1035 / RFC 3596) | Parallel A+AAAA lookups are the norm (getaddrinfo) — broken middleboxes mishandling that pattern cause hangs (§8) |
| **NS** | Names a zone's authoritative servers; parent-side NS RRset *is* the delegation | Keep parent and child NS RRsets identical |
| **CNAME** | Owner is an alias for the target **for all types**; no other data may coexist at the node (RFC 1034) | **Never legal at the apex** — apex must hold SOA+NS, which can't coexist with CNAME[^5]. CNAME+MX/TXT coexistence violations break mail |
| **SRV** (RFC 2782)[^6] | `_service._proto.name TTL IN SRV priority weight port target` — lowest priority first; weight load-splits within a priority; target "." = service explicitly absent | Target MUST NOT be an alias (widely violated, formally illegal); target must have address records, not be an IP literal |
| **TXT** | Free-form; carries SPF, DKIM, domain-verification tokens | 255-byte limit is per **character-string**, not per record — split long values into multiple quoted strings in ONE RR; consumers concatenate without spaces[^7]. Exactly one SPF record per name (RFC 7208) |
| **SOA** | 7 fields: MNAME (primary), RNAME (mailbox, first dot = @), SERIAL (zone version), REFRESH, RETRY, EXPIRE, MINIMUM | Keep RETRY < REFRESH < EXPIRE; forgetting a SERIAL bump strands secondaries. MINIMUM now means **negative-cache TTL** (RFC 2308, §3) |
| **CAA** (RFC 8659)[^8] | `flags tag value`; tags `issue`, `issuewild`, `iodef`; CA climbs from FQDN toward root, **first CAA RRset found wins and search stops** | `issue ";"` = no issuance. Absent `issuewild` → `issue` governs wildcards too. Critical flag (bit 0): CA MUST NOT issue if it can't understand the tag. CAs are required (CA/Browser Forum) to check CAA at issuance — list every CA you actually use |
| **PTR** | Reverse mapping under in-addr.arpa (octets reversed) / ip6.arpa (nibble-reversed) | Driven mostly by outbound-mail reputation (FCrDNS); delegated via RIRs |
| **MX** | Exchange + 16-bit preference (lower = preferred) | **Null MX** (RFC 7505): single MX, preference 0, exchange "." = "no mail accepted" |
| **SVCB / HTTPS** (RFC 9460)[^9] | **AliasMode** (SvcPriority=0): standards-based apex aliasing without touching other types. **ServiceMode** (>0): endpoint + SvcParams (`alpn`, `port`, `ipv4hint`/`ipv6hint`, `ech`, `mandatory`) | The standards answer to CNAME-at-apex; client support still uneven. HTTPS RR also carries the ECH config (§5) |

**Apex aliasing decision tree**: plain A/AAAA (static) → provider ALIAS/ANAME/CNAME-flattening
(**never standardized** — draft-ietf-dnsop-aname expired[^10]; TTL ownership and geo-routing degrade
because the *authoritative server's* resolver does target selection, not the client) → HTTPS-RR
AliasMode (standard, support growing). Never a literal CNAME at the apex.

## 3. TTLs, caching & negative caching

- TTL is the **maximum** cache duration; caches decrement it, so the population-average effective
  TTL ≈ half the authoritative value. **Effective TTL = clamp(record TTL, resolver floor, resolver
  cap)** — your published TTL is a request, not a command. Known clamps (verified-as-of 2026-06-12):
  Unbound `cache-max-ttl` 86400s default; BIND `max-cache-ttl` 1 week, `max-ncache-ttl` 3h; Google
  Public DNS caps positive TTLs at 6h[^11]. Many resolvers floor TTLs at 60–300s — sub-minute
  failover via TTL alone is not dependable.
- **Negative caching (RFC 2308)**[^12]: negative TTL = **min(SOA MINIMUM, SOA record's own TTL)**,
  delivered via the SOA in the authority section. **NXDOMAIN** is cached per `<QNAME, QCLASS>` — it
  suppresses *every* type at that name; **NODATA** (NOERROR, zero answers) is cached per
  `<QNAME, QTYPE, QCLASS>`. Recommended negative TTL 1–3h.
- **Serve-stale (RFC 8767)**: resolvers MAY answer with expired records when refresh fails;
  suggested client-response timer ~1.8s, stale answers served with TTL=30, retention 1–3 days.
  Off by default in Unbound (`serve-expired: no`); BIND requires `stale-answer-enable`. ISC itself
  criticizes serve-stale (masks real outages; can serve a stale A after the name became a CNAME) —
  treat it as a niche resilience tool, not a default[^13].
- **Cache poisoning context**: pre-2008, 16-bit TXID + fixed source port made off-path spoofing
  feasible (Kaminsky-class: race random-subdomain queries + in-bailiwick extra data). Source-port
  randomization added ~16 bits of entropy as mitigation; DNSSEC is the protocol-level fix.

**Migration patterns**: lower TTLs **at least one full old-TTL period before** a cutover; lower the
SOA negative TTL before adding brand-new names; never query a name before it's published (you'll
cache an NXDOMAIN that blocks all types — the classic "works everywhere except here"). Ultra-low
TTLs (≤60s) as a default posture destroy cache efficiency for negligible agility (APNIC
measurement: ~half the internet runs ≤1-min TTLs)[^14].

## 4. DNSSEC

DNSSEC provides **origin authentication and integrity** — not confidentiality, not DoS protection
(RFC 4033)[^15].

- **Records** (RFC 4034)[^16]: **DNSKEY** (flags 256 = ZSK, 257 = KSK/SEP); **RRSIG** (signature
  with absolute inception/**expiration** timestamps — the root of the most common outage class);
  **DS** (digest of a DNSKEY, lives in the *parent* — the cross-zone trust link); **NSEC/NSEC3**
  (authenticated denial; NSEC permits trivial zone-walking, NSEC3 hashes names but remains
  offline-dictionary-attackable). **RFC 9276**: NSEC3 iterations **MUST be 0**, salt SHOULD be empty
  (`1 0 0 -`) — extra iterations burn CPU for nothing and validators may SERVFAIL high-iteration
  zones[^17]. **CDS/CDNSKEY** (RFC 8078) automate DS updates parent-side; extended by RFC 9615
  (authenticated bootstrap) and RFC 9859; ccTLDs scanning CDS include .ch, .cz, .se, .za
  (verified-as-of 2026-06-12).
- **Chain of trust**: root trust anchor → DNSKEY → DS → DNSKEY → … → RRset. Four validation states
  (RFC 4033): **Secure**, **Insecure** (proven unsigned delegation), **Bogus** (should validate,
  doesn't), **Indeterminate**. Bogus → **SERVFAIL** to clients — security failure is
  indistinguishable from server failure at the RCODE level. **AD bit** = resolver asserts validated;
  **CD bit** = "don't validate upstream, give me the data."
- **Root KSK**: RFC 5011 automated rollover (30-day add hold-down; REVOKE flag bit 8). The first
  root rollover (KSK-2010→KSK-2017) executed 2018-10-11 after a year's postponement when telemetry
  showed stale trust anchors[^18].
- **Algorithms** (verified-as-of 2026-06-12): use **13 (ECDSA P-256/SHA-256)** or **15 (Ed25519)**
  for new deployments; 8 (RSA/SHA-256) remains fine. **RFC 9904 (2025) obsoletes RFC 8624** and
  moves guidance into IANA registry columns; SHA-1 DS digests MUST NOT be used[^19].
- **Deployment reality** (verified-as-of 2026-06-12; figures conflict — preserved): ~30–36% of users
  sit behind validating resolvers (APNIC); ~92% of TLDs are signed but only ~4–5% of .com
  delegations; query-weighted end-to-end use estimates range ~0.47%–1% depending on methodology.
  Adoption is stalling: the downside of misconfiguration (total SERVFAIL outage) is asymmetric to
  the upside[^20].

**Failure catalog**: (1) **expired RRSIGs** from a stalled signer — most common class; (2)
**DS/DNSKEY mismatch after operator/registrar migration** — moving DNS hosts while signed without
re-coordinating DS; (3) **Slack 2021**: Route 53 emitted wrong NSEC type bitmaps on wildcard
responses; resolvers cached them as proof-of-NODATA; rollback throttled by a 24h DS TTL[^21]; (4)
**.de outage 2026-05-05**: DENIC published non-validatable DNSSEC signatures during a scheduled
key rollover starting ~19:30 UTC; validating resolvers SERVFAILed huge swaths of .de; Cloudflare
ended impact for 1.1.1.1 at 22:17 UTC (~2h47m) by treating .de as insecure — a Negative-Trust-Anchor
intervention (RFC 7646)[^22].

**Patterns**: automate re-signing with monitoring on RRSIG expiry headroom; short DS TTLs around
changes; discriminate DNSSEC-SERVFAIL from infra-SERVFAIL with `+cd` (§9); use `delv` to find the
exact broken link.

## 5. Encrypted DNS

| Transport | Spec | Port / framing | Properties |
|---|---|---|---|
| **DoT** | RFC 7858[^23] | TCP 853, TLS, 2-octet length framing; ALPN `dot` mandated by later specs (RFC 9539/9463) | Trivially identifiable/blockable; TCP head-of-line blocking |
| **DoH** | RFC 8484[^24] | HTTPS 443; `application/dns-message`; GET (`?dns=` base64url) and POST both MUST | Indistinguishable from web traffic — anti-censorship strength, enterprise-control complaint; HTTP/2 multiplexing (transport HoL remains) |
| **DoQ** | RFC 9250[^25] | QUIC UDP 853, ALPN `doq`; one stream per query; Message ID MUST be 0 | No transport HoL; 1-RTT (+0-RTT for replayable QUERY/NOTIFY **only**); measured ~33% faster than DoT/DoH single-query (PAM 2022); resolver support still narrow (AdGuard, NextDNS) |

- **Discovery**: **DDR (RFC 9462)** — client queries SVCB for `_dns.resolver.arpa` at its known
  Do53 resolver; *verified* discovery requires the designated resolver's TLS cert to contain the
  original resolver IP. **DNR (RFC 9463)** — network provisioning via DHCPv6 option 144 / DHCPv4
  option 162 / RA option 144; takes precedence over DDR.
- **Oblivious DoH (RFC 9230, Experimental)**[^26]: client → proxy → target with HPKE; proxy sees
  who-but-not-what, target sees what-but-not-who; collapses if they collude. Niche but alive
  (Cloudflare target + partner proxies); the standards-track generalization is **Oblivious HTTP
  (RFC 9458)** (verified-as-of 2026-06-12).
- **ECH** (verified-as-of 2026-06-12): **published March 2026 as RFC 9849**, with **RFC 9848**
  defining the `ech` SvcParam (ECHConfigList in HTTPS/SVCB records)[^27]. Encrypts the entire
  ClientHello (SNI + ALPN) toward the client-facing server; outer SNI shows only the fronting
  provider. **ECH is only meaningful over encrypted DNS** — RFC 9848 states the same name otherwise
  leaks in the DNS query. Deployment is concentrated (Cloudflare + Meta per SIGCOMM 2025); Chrome
  ~117+/Firefox 118+ enable it gated on DoH; Russia blocks Cloudflare ECH (Nov 2024); China blocks
  the encrypted-DNS prerequisite.
- **The gap**: encryption is overwhelmingly **stub→recursive**. Recursive→authoritative (ADoT) has
  only **RFC 9539** (Experimental, 2024): unilateral opportunistic probing of port 853, passive-only
  threat model.
- **Policy plumbing**: Firefox checks the **canary domain `use-application-dns.net`** via the system
  resolver — NXDOMAIN disables *auto*-DoH (never user-chosen DoH)[^28]. Chrome auto-upgrades
  same-provider only (never switches providers) and disables Secure DNS on managed devices.

**Enterprise pattern**: run your own DoH/DoT endpoint, advertise via DNR+DDR, set browser policy,
serve the canary — blocking port 853 alone does nothing about DoH on 443. **Anti-patterns**:
deploying ECH without encrypted DNS; ODoH with same-operator proxy+target; non-replayable ops in
DoQ 0-RTT.

## 6. Split-horizon & private DNS

**Split-horizon** (= split-view/split-brain) DNS returns different answers for the same name,
usually keyed on query source address — internal clients get RFC 1918 addresses, external clients
the public view[^29].

- **BIND views**[^30]: `match-clients`/`match-destinations` ACLs, evaluated **in order, first match
  wins** (a `match-clients { any; }` view listed first shadows everything); once any view exists,
  all zones must live inside views. Each view = independent zone copies **and its own cache**
  (shareable via `attach-cache`; `in-view` shares a zone instance). Primary/secondary pairs with
  views need **TSIG keys in the match ACLs** so transfers land in the right view. ISC acknowledges
  views' complexity reputation.
- **Other servers**: Unbound `access-control-tag`/`view:` (local-data only — per-view *forwarding*
  is a long-open feature request); PowerDNS Auth gained BIND-style views in **5.0.0 (Aug 2025)** —
  narrowest-netmask wins, ordering irrelevant; dnsmasq `server=/example.corp/10.0.0.1` conditional
  forwarding, `address=/domain/IP` local authority. BIND conditional forwarding: per-zone
  `type forward; forward only|first`.
- **Namespace choice** (verified-as-of 2026-06-12): **never squat unregistered TLDs**. `.local` is
  reserved for mDNS (RFC 6762) — site-wide unicast `.local` (the old AD `corp.local` pattern)
  collides with Bonjour/Avahi. `.home/.corp/.mail` are ICANN high-risk collision strings, deferred
  indefinitely. **`home.arpa`** (RFC 8375) is for homenets. **`.internal`** was provisionally selected by IANA
  (Jan 2024) and reserved for private use by ICANN Board resolution (2024-07-29 — cite the ICANN
  Board materials for the resolution itself; the IANA news page covers only the provisional
  determination)[^31]; the companion IETF draft (draft-davies-internal-tld) is still an I-D, not an
  RFC. Caveats: DNSSEC-validating resolvers will
  *prove .internal's nonexistence* (root NSEC) unless the local resolver serves it beneath the
  validation point, and **no public CA can issue for .internal** — you need a private CA. Best
  practice remains a subdomain of a registered domain you own.
- **RFC 1918 reverse**: resolvers should serve empty RFC 1918 in-addr.arpa zones locally
  (RFC 6303/BCP 163); otherwise private PTR queries leak to the **AS112** anycast sink
  (RFC 7534)[^32]. Internal views with real PTR data must explicitly override the built-in empty
  zones (BIND `empty-zones-enable`).
- **Certificates**: public CAs stopped issuing for internal names in 2015–2016 (CA/B Forum).
  Split-horizon on a **public domain you own** is what keeps publicly-trusted certs (incl. ACME
  DNS-01) working for internally resolved names[^33].

**Failure modes**: cache mixing across horizons (roaming clients, shared forwarders); **DoH-bypassing
clients** that both break internal resolution and exfiltrate the internal namespace to a public
resolver; DNSSEC vs two-truths (sign per-view consistently or don't sign the internal view);
monitoring blind spots — **probe both views explicitly** in change control; view-ordering bugs.
Split DNS hides names as hygiene, not access control.

## 7. EDNS(0)

RFC 6891[^34]: the **OPT pseudo-RR** (TYPE 41) rides in the Additional section — exactly one per
message, never cached, hop-by-hop. Field overloading: CLASS = requestor's max UDP payload;
TTL field = extended-RCODE byte | version | **DO bit** | Z. Version mismatch → BADVERS.

- **Why**: classic DNS capped UDP at 512 bytes; DNSSEC answers don't fit.
- **Fragmentation history**: 4096-byte buffers invited IP fragmentation — unreliable (middleboxes
  drop fragments) and a poisoning vector. **DNS Flag Day 2019** removed the retry-without-EDNS
  workaround (timeouts now mean server-down). **DNS Flag Day 2020** set the de-facto default buffer
  to **1232 bytes** (IPv6 min MTU 1280 − 48); BIND ≥9.16.8 and Unbound default to it; oversized
  answers truncate (TC=1) and complete over **TCP** (RFC 7766 makes TCP support mandatory)[^35].
  Note the live contradiction: RFC 6891's own 4096 guidance is superseded by Flag Day 2020 practice.
- **ECS (RFC 7871)**: option 8 attaches a truncated client prefix (/24 v4, /56 v6 recommended) for
  CDN geo-steering; SOURCE PREFIX-LENGTH=0 is the client opt-out. The RFC is Informational and
  self-deprecating on privacy. Cloudflare 1.1.1.1 deliberately sends no ECS (relies on anycast
  density), measurably degrading some CDNs' steering; Google Public DNS sends it — a genuine
  privacy/performance fork[^36].
- **Extended DNS Errors (RFC 8914)**: option 15; 16-bit INFO-CODE + free text. Turns opaque
  SERVFAILs into "6 = DNSSEC Bogus", "9 = DNSKEY Missing", "15 = Blocked". Diagnostic only,
  unauthenticated.

**Patterns**: leave buffers at 1232; permit TCP/53 through firewalls (blocking it breaks
DNSSEC-sized answers outright); when comparing dig results across hosts pin `+bufsize=` (dig ≥9.18
defaults to 1232, older to 4096).

## 8. Resolver behavior differences

Scope note: this section covers Linux (systemd-resolved, glibc, musl), macOS, and Kubernetes.
Windows DNS Client (NRPT, `ipconfig /displaydns`, GPO DoH policy) is not yet covered — research
queued as a follow-up.

### systemd-resolved (Linux)

- Stub listener **127.0.0.53:53** (caching, LLMNR/mDNS integration, optional DNSSEC);
  **127.0.0.54** = pass-through proxy stub. Four `/etc/resolv.conf` modes: stub symlink
  (recommended), static, uplink, foreign — check with `ls -l /etc/resolv.conf`[^37].
- **Routing**: per-link DNS + `Domains=`; `~example.com` = route-only domain; **most-labels match
  wins**; `~.` makes a link the default DNS route. `.local` is never sent to unicast DNS unless
  configured as a search/routing domain (`Domains=~local` is the corporate-unicast-.local
  workaround). Synthesized names (`localhost`, `_gateway`, /etc/hosts) never hit the wire.
- **DNSSEC=** upstream default is `allow-downgrade` (documented as downgrade-attackable), but
  **mainstream distros ship DNSSEC=no** (Fedora explicitly; Ubuntu disabled it after bugs) — assume
  validation is OFF unless you enabled it (verified-as-of 2026-06-12)[^38].
- Operate with `resolvectl status|query|flush-caches`. `resolvectl query` exercises the real OS
  path (routing domains + cache) — dig does not.

### macOS — mDNSResponder & scoped resolvers

- **mDNSResponder** is the system resolver (Apple's discoveryd replacement experiment lasted
  10.10–10.10.4, 2015, then reverted). Apps resolve via `getaddrinfo` → mDNSResponder using the
  full **scoped-resolver** configuration; **`/etc/resolv.conf` is auto-generated and misleading** —
  only BSD-resolver clients like dig read it, so dig output can disagree with what apps do[^39].
- **`/etc/resolver/<domain>`** files create per-domain resolvers (keys: `nameserver`, `port`,
  `search_order`, `timeout`, `options ndots:n`); the "Super" DNS client routes each query to the
  client with the most matching labels — the sanctioned mechanism for unicast `.local` coexistence
  and VPN split-DNS.
- Inspect with `scutil --dns`; query the real path with `dscacheutil -q host -a name <host>`; flush:
  `sudo dscacheutil -flushcache; sudo killall -HUP mDNSResponder`.

### glibc stub resolver

- `/etc/resolv.conf` options[^40]: `ndots:` default **1** (cap 15); `timeout:` 5s; `attempts:` 2;
  **MAXNS=3** nameservers tried in order; `rotate`, `use-vc` (force TCP), `no-aaaa`, `trust-ad`.
  `getaddrinfo` sends **A and AAAA in parallel on one socket**; `single-request`/
  `single-request-reopen` serialize it for broken middleboxes. Search list unlimited since
  glibc 2.26; `LOCALDOMAIN`/`RES_OPTIONS` env overrides.
- `/etc/nsswitch.conf` `hosts:` ordering (`files dns` etc.) is why /etc/hosts wins; action modifiers
  `[NOTFOUND=return]` change fallthrough; glibc ≥2.33 auto-reloads it.

### musl (Alpine) differences

- Queries **all nameservers in parallel**, first answer wins — "primary then fallback" semantics
  silently disappear; inconsistent upstreams become nondeterministic[^41].
- **No TCP fallback until musl 1.2.4** — truncated responses simply failed (the classic Alpine
  large-response bug). `search`/`domain` unsupported before 1.1.13; names with ≥ndots dots are
  tried **only** absolute (no search fallback, unlike glibc).

### Kubernetes ndots:5

Pod resolv.conf ships `search <ns>.svc.cluster.local svc.cluster.local cluster.local` + `ndots:5`:
external names with <5 dots walk the whole search list first (3+ NXDOMAIN round trips × parallel
A+AAAA) — millions of junk queries/hour against CoreDNS at scale[^42]. Mitigations: trailing-dot
FQDNs, per-pod `dnsConfig` ndots:1–2, NodeLocal DNSCache. Counterpoint: blanket-lowering ndots
breaks short internal names, differently on musl vs glibc pods — tune per workload.

## 9. DNS debugging — dig, delv, drill, kdig

### dig essentials

`dig [@server] name type [+opts]`; `.digrc` applies silently — **use `-r` in scripts**. Exit code 0
includes NXDOMAIN — **parse the `status:` line, never the exit code**[^43].

| Flag | Use |
|---|---|
| `+short +identify` | terse answer + which server answered |
| `+noall +answer` | stable tab-separated scripting output |
| `+norecurse` | ask a cache "what do you already have" / clean authoritative queries |
| `+dnssec` | set DO bit, display RRSIGs — **dig never validates**; AD in a response is the upstream resolver's claim |
| `+cd` | checking-disabled: retrieve data a validator would suppress |
| `+trace` | client-side iteration from the root — see caveats below |
| `+nssearch` | SOA serial from every authoritative NS — one-shot lagging-secondary check |
| `+tcp`, `+ignore`, `+bufsize=N` | truncation/fragmentation probing |
| `+subnet=A/L`, `+subnet=0` | send/suppress ECS when debugging geo answers |
| `+https[=path]`, `+tls` | DoH/DoT endpoint testing (BIND ≥9.18-era dig; no +quic — use kdig) |

**Header reading**: flags `qr aa tc rd ra ad cd`; status **NOERROR** (may still be NODATA),
**NXDOMAIN**, **SERVFAIL** (infra *or* DNSSEC — indistinguishable alone), **REFUSED** (policy:
wrong server for the zone, or recursion denied).

**`dig +trace` caveats (ISC)**[^44]: it simulates a **cold cache** — it shows the delegation path,
*not* what your resolver will answer; `@server` only scopes the initial root-NS lookup; glueless
referrals are chased through **your /etc/resolv.conf servers**, so split-DNS environments can
silently contaminate the trace. Prefer `delv +ns` for faithful cold-cache resolver emulation.

### delv — DNSSEC triage

delv uses **named's own resolver+validator code** and fully validates against a built-in root trust
anchor, labeling answers "fully validated"/"unsigned" and explaining *why* validation failed
(vs dig's display-only `+dnssec`)[^45]. Key moves: `+vtrace` (validator decisions); `+cd` (pull
bogus data and self-validate — the canonical "why is this SERVFAIL"); `delv +ns` (full recursive
cold-cache emulation, the better +trace). Limits: static trust anchor (no RFC 5011 tracking — an
old binary breaks after a root KSK roll), class IN only.

### drill / kdig

- **drill** (ldns): `-T` trace from root (@server ignored), `-S` chase signatures upward, `-D` DO
  bit; **exit code is DNSSEC-aware** (non-zero = bogus/error) unlike dig's transport-only codes.
- **kdig** (Knot): dig-compatible plus the broadest encrypted transports — `+tls`, `+https`,
  **`+quic`**, `+json`, EDNS bufsize default 1232.

### Recipes

1. **Cache vs authoritative discrepancy**: `dig name @recursive` vs
   `dig name @ns1.zone +norecurse` (check AA); add `dig +nssearch zone` to compare serials across
   all secondaries.
2. **Delegation check**: `dig +noall +answer +additional NS zone @parent-tld-server` (NS without
   additional A/AAAA = missing glue) vs `dig NS zone @child-ns +norecurse` — sets must match.
3. **SERVFAIL root-cause**: re-run with `+cd`. Works with +cd but not without → **DNSSEC validation
   failure** → `delv @resolver name +vtrace` for the broken link. Fails both ways → infrastructure.
   Ask for RFC 8914 EDE codes where supported.
4. **Split-DNS verification**: loop `dig +short +identify name @resolver` across
   internal/public resolvers — public NXDOMAIN for internal-only names is the *expected* state.
5. **OS-path vs wire**: `resolvectl query` (Linux) / `dscacheutil -q host` (macOS) exercise the OS
   resolver path that dig deliberately bypasses (§8).

**Anti-patterns**: trusting +trace as "what my resolver sees"; scripting on dig exit codes; reading
AD as proof of local validation; interpreting REFUSED as an outage; forgetting `.digrc`.

## Cross-references

- **MongoDB SRV seedlist / Atlas DNS**: `mongodb-expert` → references/mongodb-connection-string.md
  (SRV/TXT seedlist format), `mongodb-atlas-expert` → references/mongodb-aws-networking.md
  (PrivateLink/VPC DNS) — protocol mechanics live here; product behavior lives there.
- **Host-level diagnostics** (ip/ss/tcpdump/nmap and general sysadmin DNS triage):
  `devops-linux-admin` → references/linux-sysadmin.md.
- **Cloudflare 1.1.1.1, proxied DNS, CNAME flattening as a product**: `cloudflare-platform`.
- **Sibling networking topics** (own runs/spokes): TCP & congestion control, QUIC, TLS internals,
  load balancing, VPN/ZTNA, IP addressing — see the `networking` hub routing table as they land.

## References

[^1]: RFC 9499 — DNS Terminology. https://www.rfc-editor.org/rfc/rfc9499.html
[^2]: Netnod — DNS root server FAQ. https://www.netnod.se/dns/dns-root-server-faq
[^3]: RFC 8109 — Initializing a DNS Resolver with Priming Queries. https://www.rfc-editor.org/rfc/rfc8109.html
[^4]: RFC 9156 — DNS Query Name Minimisation to Improve Privacy. https://www.rfc-editor.org/rfc/rfc9156.html
[^5]: ISC KB — CNAME at the apex of a zone. https://kb.isc.org/docs/aa-01640
[^6]: RFC 2782 — A DNS RR for specifying the location of services (SRV). https://www.rfc-editor.org/rfc/rfc2782.html
[^7]: ISC KB — TXT/SPF records longer than 255 characters. https://kb.isc.org/docs/aa-00356
[^8]: RFC 8659 — DNS Certification Authority Authorization (CAA). https://datatracker.ietf.org/doc/html/rfc8659
[^9]: RFC 9460 — Service Binding via SVCB and HTTPS RRs. https://datatracker.ietf.org/doc/html/rfc9460
[^10]: draft-ietf-dnsop-aname (expired — ANAME never standardized). https://datatracker.ietf.org/doc/draft-ietf-dnsop-aname/
[^11]: Google Public DNS FAQ (6h TTL cap). https://developers.google.com/speed/public-dns/faq ; Unbound defaults: https://nlnetlabs.nl/documentation/unbound/unbound.conf/ ; BIND ARM: https://bind9.readthedocs.io/en/latest/reference.html
[^12]: RFC 2308 — Negative Caching of DNS Queries. https://www.rfc-editor.org/rfc/rfc2308
[^13]: RFC 8767 — Serving Stale Data to Improve DNS Resiliency. https://www.rfc-editor.org/rfc/rfc8767 ; ISC criticism: https://www.isc.org/blogs/2020-serve-stale/
[^14]: APNIC — Stop using ridiculously low DNS TTLs. https://blog.apnic.net/2019/11/12/stop-using-ridiculously-low-dns-ttls/
[^15]: RFC 4033 — DNS Security Introduction and Requirements. https://www.rfc-editor.org/rfc/rfc4033
[^16]: RFC 4034 — Resource Records for DNSSEC. https://www.rfc-editor.org/rfc/rfc4034 ; NSEC3: RFC 5155 https://www.rfc-editor.org/rfc/rfc5155
[^17]: RFC 9276 — Guidance for NSEC3 Parameter Settings. https://www.rfc-editor.org/rfc/rfc9276
[^18]: RFC 5011 — Automated Updates of DNSSEC Trust Anchors. https://www.rfc-editor.org/rfc/rfc5011 ; KSK roll analysis: https://blog.apnic.net/2018/10/31/analyzing-the-ksk-roll/
[^19]: RFC 9904 — DNSSEC Algorithm Recommendations (obsoletes 8624). https://www.rfc-editor.org/rfc/rfc9904 ; IANA registry: https://www.iana.org/assignments/dns-sec-alg-numbers/dns-sec-alg-numbers.xhtml
[^20]: Huston — Measuring the use of DNSSEC. https://circleid.com/posts/20230910-measuring-the-use-of-dnssec ; live map: https://stats.labs.apnic.net/dnssec
[^21]: Slack engineering — What happened during Slack's DNSSEC rollout. https://slack.engineering/what-happened-during-slacks-dnssec-rollout/
[^22]: Cloudflare — the .de DNSSEC outage. https://blog.cloudflare.com/de-tld-outage-dnssec/
[^23]: RFC 7858 — DNS over TLS. https://www.rfc-editor.org/rfc/rfc7858
[^24]: RFC 8484 — DNS Queries over HTTPS. https://www.rfc-editor.org/rfc/rfc8484
[^25]: RFC 9250 — DNS over Dedicated QUIC Connections. https://www.rfc-editor.org/rfc/rfc9250 ; DoQ measurement (PAM 2022): https://vaibhavbajpai.com/documents/papers/proceedings/doq-pam-2022.pdf
[^26]: RFC 9230 — Oblivious DNS over HTTPS. https://www.rfc-editor.org/rfc/rfc9230 ; RFC 9458 — Oblivious HTTP. https://www.ietf.org/rfc/rfc9458.html
[^27]: RFC 9849 — TLS Encrypted Client Hello. https://datatracker.ietf.org/doc/rfc9849/ ; RFC 9848 — Bootstrapping ECH with DNS Service Bindings. https://www.rfc-editor.org/rfc/rfc9848 ; DDR: RFC 9462 https://www.rfc-editor.org/rfc/rfc9462 ; DNR: RFC 9463 https://www.rfc-editor.org/rfc/rfc9463 ; ADoT probing: RFC 9539 https://www.rfc-editor.org/rfc/rfc9539
[^28]: Mozilla — canary domain use-application-dns.net. https://support.mozilla.org/en-US/kb/canary-domain-use-application-dnsnet
[^29]: Wikipedia — Split-horizon DNS. https://en.wikipedia.org/wiki/Split-horizon_DNS
[^30]: ISC KB — Understanding views in BIND 9. https://kb.isc.org/docs/aa-00851 ; PowerDNS 5.0 views: https://doc.powerdns.com/authoritative/views.html
[^31]: IANA — proposed private-use TLD (.internal). https://www.iana.org/news/2024/proposed-private-use-tld ; draft-davies-internal-tld: https://datatracker.ietf.org/doc/draft-davies-internal-tld/ ; home.arpa: RFC 8375 https://www.rfc-editor.org/rfc/rfc8375.html ; mDNS .local: RFC 6762 https://www.rfc-editor.org/rfc/rfc6762.html
[^32]: RFC 6303 — Locally Served DNS Zones. https://www.rfc-editor.org/rfc/rfc6303.html ; RFC 7534 — AS112 Nameserver Operations. https://www.rfc-editor.org/rfc/rfc7534.html
[^33]: CA/Browser Forum — internal names guidance. https://cabforum.org/working-groups/server/internal-names/
[^34]: RFC 6891 — Extension Mechanisms for DNS (EDNS(0)). https://www.rfc-editor.org/rfc/rfc6891
[^35]: APNIC — DNS Flag Day 2020. https://blog.apnic.net/2020/09/17/dns-flag-day-2020-what-you-need-to-know/ ; ISC dig bufsize behavior: https://kb.isc.org/docs/behavior-dig-versions-edns-bufsize ; RFC 7766 — DNS over TCP: https://datatracker.ietf.org/doc/html/rfc7766
[^36]: RFC 7871 — EDNS Client Subnet. https://www.rfc-editor.org/rfc/rfc7871 ; APNIC Labs — Privacy and ECS: https://labs.apnic.net/index.php/2024/07/23/privacy-and-dns-client-subnet/ ; RFC 8914 — Extended DNS Errors: https://www.rfc-editor.org/rfc/rfc8914
[^37]: systemd-resolved.service(8). https://man7.org/linux/man-pages/man8/systemd-resolved.service.8.html ; resolved.conf(5): https://man7.org/linux/man-pages/man5/resolved.conf.5.html
[^38]: Fedora — systemd-resolved change (DNSSEC=no). https://fedoraproject.org/wiki/Changes/systemd-resolved
[^39]: macOS resolver(5). https://manp.gs/mac/5/resolver ; Understanding DNS on macOS: https://mikebian.co/understanding-dns-requests-on-macos/
[^40]: resolv.conf(5). https://man7.org/linux/man-pages/man5/resolv.conf.5.html ; nsswitch.conf(5): https://manpages.debian.org/bookworm/manpages/nsswitch.conf.5.en.html
[^41]: musl — functional differences from glibc. https://wiki.musl-libc.org/functional-differences-from-glibc.html
[^42]: Pracucci — Kubernetes DNS resolution and ndots. https://pracucci.com/kubernetes-dns-resolution-ndots-options-and-why-it-may-affect-application-performances.html
[^43]: dig(1), BIND 9.20. https://manpages.debian.org/unstable/bind9-dnsutils/dig.1.en.html
[^44]: ISC KB — dig and the +trace option. https://kb.isc.org/docs/aa-00208
[^45]: delv(1), BIND 9.20. https://manpages.debian.org/unstable/bind9-dnsutils/delv.1.en.html ; drill(1): https://manpages.debian.org/unstable/ldnsutils/drill.1.en.html ; kdig: https://www.knot-dns.cz/docs/latest/html/man_kdig.html
