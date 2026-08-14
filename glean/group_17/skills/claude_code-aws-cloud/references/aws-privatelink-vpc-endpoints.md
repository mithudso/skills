<!-- hub-reference-banner -->
> **Reference file — part of the `aws-cloud` hub.** Created by /dr deep-research (2026-06-12) on
> AWS PrivateLink and VPC endpoints — the general mechanism. Sibling topics in this family are
> reference files under the hubs (`aws-cloud`) — **not** standalone skills. For the MongoDB-Atlas
> CONSUMER side of private connectivity, do not use this file alone: see `mongodb-aws-networking`,
> `mongodb-atlas-azure`, and `mongodb-atlas-gcp` under the `mongodb-atlas-expert` hub.

---

---
name: aws-privatelink-vpc-endpoints
description: >
  AWS PrivateLink and VPC endpoints — the general mechanism, both sides: interface vs gateway vs
  Gateway Load Balancer (GWLB) vs resource endpoints; provider-side endpoint services (NLB/GWLB,
  acceptance workflow, allowed principals, cross-account and cross-region, proxy protocol v2);
  private DNS integration and split-horizon resolution of endpoint hostnames; endpoint policies
  and data perimeters; security groups and NACLs on endpoint ENIs; quotas, bandwidth scaling, and
  pricing; PrivateLink vs VPC peering vs Transit Gateway decision framework; troubleshooting
  connectivity through endpoints.
  TRIGGER: choosing an endpoint type; exposing a service over PrivateLink; pendingAcceptance;
  endpoint hostname resolves to public IPs; on-prem resolution of vpce names; endpoint policy or
  data-perimeter design; SG/NACL rules for endpoint ENIs; 350s idle timeouts through endpoints;
  centralized endpoints hub-and-spoke; PrivateLink vs peering vs TGW.
  SKIP: Atlas private endpoints (consumer side) — mongodb-aws-networking / mongodb-atlas-azure /
  mongodb-atlas-gcp; DNS protocol internals — networking hub; general VPC design, NAT, SG-vs-NACL
  basics — aws-core.
version: "1.0.0"
updated: "2026-06-12"
category: developer
whenToUse: >
  Use when designing, operating, securing, pricing, or debugging AWS private connectivity built on
  VPC endpoints or PrivateLink endpoint services — on either the consumer or the provider side —
  or when choosing between PrivateLink, VPC peering, and Transit Gateway.
keywords:
  - aws privatelink
  - vpc endpoint
  - interface endpoint
  - gateway endpoint
  - gateway load balancer endpoint
  - endpoint service
  - private dns
  - split-horizon dns
  - endpoint policy
  - vpce
  - transit gateway vs privatelink
  - vpc peering vs privatelink
  - cross-region privatelink
  - proxy protocol v2
tags:
  - aws
  - networking
  - privatelink
  - vpc
  - dns
  - security
---

# AWS PrivateLink & VPC Endpoints — the General Mechanism

## Contents

1. [Overview](#overview)
2. [Endpoint types](#endpoint-types)
3. [Security groups and NACLs on endpoint ENIs](#security-groups-and-nacls-on-endpoint-enis)
4. [Endpoint services — the provider side](#endpoint-services--the-provider-side)
5. [Private DNS and split-horizon resolution](#private-dns-and-split-horizon-resolution)
6. [Endpoint policies](#endpoint-policies)
7. [Quotas, scaling, pricing](#quotas-scaling-pricing)
8. [Decision framework: PrivateLink vs VPC peering vs Transit Gateway](#decision-framework-privatelink-vs-vpc-peering-vs-transit-gateway)
9. [Troubleshooting connectivity through endpoints](#troubleshooting-connectivity-through-endpoints)
10. [Azure Private Link / GCP Private Service Connect equivalence](#azure-private-link--gcp-private-service-connect-equivalence)
11. [Anti-patterns](#anti-patterns)
12. [Cross-references](#cross-references)
13. [References](#references)

## Overview

AWS PrivateLink lets a consumer VPC reach a service (an AWS service, a SaaS product, or another
team's NLB-fronted application) through an **elastic network interface with a private IP inside the
consumer's own subnets**, instead of over the internet, a peering connection, or a Transit Gateway.
The consumer never learns the provider's network topology; the provider never gains a route into
the consumer. Connections are **unidirectional (consumer → service)** and the two VPCs' CIDRs are
never compared — overlapping address space is a non-issue.[^1][^16]

"VPC endpoint" is the umbrella term. AWS distinguishes **five endpoint types** — interface,
gateway, Gateway Load Balancer, resource, and service-network — and all of them are powered by
PrivateLink **except gateway endpoints**, which predate PrivateLink and are pure route-table
constructs.[^1][^2]

## Endpoint types

### Interface endpoints

- One **requester-managed ENI per chosen subnet** (max one subnet per AZ), with a private IP from
  the subnet CIDR. You can view but not manage the ENI; its IP does not change for the life of the
  endpoint.[^2][^3]
- Front most AWS services (SQS, STS, Secrets Manager, ECR, CloudWatch, KMS …), AWS Marketplace /
  SaaS services, and customer endpoint services.
- Traffic: **TCP and UDP**. UDP support shipped 2024-10-31 (with dual-stack NLBs); before that
  PrivateLink was TCP-only, and many secondary sources still say "TCP only" — they are stale. No
  ICMP traverses PrivateLink: you cannot ping an endpoint; probe with TCP instead.[^4]
- **MTU 8500 bytes**; larger packets are **dropped**, PMTUD is **not supported** (no ICMP
  "fragmentation needed" is generated), and MSS clamping is enforced on all packets.[^5]
- Idle timeout: established flows through interface endpoints idle out at **350 seconds** (same
  fixed timer family as NAT Gateway and NLB). The side that times out gets a TCP RST only when it
  next sends; the other side sees a half-open connection. Use TCP keepalives < 350 s for
  long-lived connections.[^6]
- AZ behavior: regional DNS name round-robins across healthy per-AZ ENIs; zonal names pin traffic
  in-AZ. Cross-AZ data transfer **to an interface endpoint is free since 2022-04-01** (does not
  apply to GWLB endpoints); per-GB data *processing* still applies.[^7]

### Gateway endpoints (S3 and DynamoDB only)

- **No ENI, no private IP, no PrivateLink.** Creating one adds a route to selected route tables
  with destination = an AWS-managed **prefix list** (e.g. `pl-63a5400a` for S3 us-east-1) and
  target = the endpoint. The service's public DNS names keep resolving to **public IPs**; the
  prefix-list route intercepts traffic to those IPs. Longest-prefix match applies.[^8]
- **Free** — no hourly or per-GB charge. This is why the default S3/DynamoDB strategy is "gateway
  endpoint unless you need access from outside the VPC".[^8][^9]
- **VPC-local only.** AWS: "Resources on the other side of a VPN connection, VPC peering
  connection, transit gateway, or Direct Connect connection in your VPC cannot use a gateway
  endpoint." Also same-region only (prefix lists are regional). For hybrid/cross-VPC access to S3,
  use an S3 **interface** endpoint.[^9]
- Access control = endpoint policy + route-table association (+ the client's SG/NACL). Client SG
  egress rules **can reference the prefix-list ID**; NACLs cannot — enumerate the service CIDRs.[^8]

### Gateway Load Balancer endpoints (GWLBE)

- A PrivateLink-powered ENI used as a **route-table target** that hands all IP packets to a
  **Gateway Load Balancer** (L3, all ports), which **GENEVE-encapsulates** them (UDP 6081) toward
  a fleet of inspection appliances (firewalls, IDS/IPS).[^10][^11]
- One AZ/subnet per GWLBE, fixed at creation; deploy one GWLBE per AZ and keep traffic zonal.
  Routing patterns use **edge association** (ingress route table on the IGW/VGW) plus subnet route
  tables; appliance MTU must accommodate GENEVE overhead (ELB docs: 68 bytes → appliance MTU ≥
  8568).[^10][^11]
- GWLB flow stickiness defaults to 5-tuple (configurable 3-/2-tuple); TCP flow idle timeout was
  fixed at 350 s, configurable **60–6000 s** since 2024-09 (`tcp.idle_timeout.seconds`); non-TCP
  flows fixed 120 s.[^11][^12]
- **No security groups and no endpoint policies** attach to GWLBEs.[^10]
- "Appliance mode" is a **Transit Gateway** setting, not a GWLB one: it pins both directions of a
  flow to one AZ so stateful appliances see symmetric traffic in multi-AZ TGW designs.[^13]

### Resource and service-network endpoints (2024-12+)

- **VPC resources over PrivateLink** (GA 2024-12-01): share an individual resource (database,
  cluster, even on-prem IP) without a load balancer, via a **Resource Gateway** (ingress point in
  the owner VPC) + **Resource Configuration** shared through AWS RAM; consumers create **resource
  endpoints**. TCP only; consumer→resource initiation only. **Service-network endpoints** attach a
  VPC to a VPC Lattice service network. Built on shared Lattice+PrivateLink infrastructure.[^14]

### S3: gateway vs interface — when each

| Need | Use |
| --- | --- |
| In-VPC access, lowest cost | Gateway endpoint (free) |
| Access from on-prem (DX/VPN), peered VPC, or TGW spoke | Interface endpoint |
| SG-level control on the path, or private DNS for S3 names | Interface endpoint |
| Both: free in-VPC + hybrid | Keep both; S3's "private DNS only for inbound resolver endpoint" option routes on-prem via the interface endpoint while VPC traffic rides the free gateway endpoint[^9][^15] |

## Security groups and NACLs on endpoint ENIs

- **Interface endpoints take security groups** (changeable later via `modify-vpc-endpoint`).
  Required inbound: the service port (usually TCP 443) from the client SGs or VPC CIDR. SGs are
  stateful — return traffic is automatic.[^3][^17]
- **If you don't pick an SG at creation, the VPC default SG is used** — which only allows inbound
  from members of itself. Clients in other SGs are then silently blocked: the single most common
  cause of "connection timed out" to a new endpoint.[^17][^18]
- The **client's** SG must allow egress to the endpoint ENI IPs/port. Pitfall: restricting client
  egress to *only* the endpoint's SG can break DNS/DHCP; allow egress to the VPC CIDR or add the
  resolver explicitly.[^17]
- Tightest pattern: endpoint SG inbound 443 **from the consumers' SG IDs**; pragmatic pattern:
  from the VPC CIDR.
- **Gateway endpoints have no ENI → no SG.** GWLBEs support neither SGs nor policies.[^8][^10]
- **NACLs** apply at the subnet boundary of the endpoint ENIs and are stateless: endpoint subnet
  needs inbound 443 from clients AND outbound **ephemeral 1024–65535** back; mirror rules on the
  client subnet. NACLs are not evaluated within a single subnet (documented exception: app→GWLBE
  in the same subnet IS evaluated).[^17][^10]

## Endpoint services — the provider side

### Creating and gating the service

- An endpoint service fronts a **Network Load Balancer** (TCP/UDP services) or **Gateway Load
  Balancer** (inspection). ALB-only backends need an NLB chained in front. Service name:
  `com.amazonaws.vpce.<region>.vpce-svc-xxxxxxxxxxxx`.[^19]
- Two independent controls:[^20]
  1. **Allowed principals** — ARNs (`arn:aws:iam::<account>:root`, a role/user, or `*`). AWS
     warning: principal `*` + auto-accept = your NLB is effectively public despite having no
     public IP. Removing a principal does NOT sever existing connections.
  2. **Acceptance** — `acceptance-required`: connection requests sit in `pendingAcceptance` until
     accepted (`accept-vpc-endpoint-connections`) or rejected. A consumer endpoint stuck in
     `pendingAcceptance`/`Rejected` **times out at connect time** — a top cross-account failure
     mode.
- **SNS connection notifications** per service: `Connect`, `Accept`, `Reject`, `Delete` events;
  the topic policy must allow `vpce.amazonaws.com` to publish.[^20]
- There is no connection draining at the PrivateLink layer; rejected endpoints should be deleted
  (consumer billing stops only on deletion).[^20]

### Client IP and NLB specifics

- **Client IP is not preserved through PrivateLink**: targets always see the NLB's private IPs.
  Remedy: enable **proxy protocol v2** on the target group — the PPv2 header carries the client
  IP *and the source VPC endpoint ID*. Both ends must speak PPv2 or connections break on the
  first byte (nginx `proxy_protocol`, Envoy/Istio listener filter, etc.). Consequence: you cannot
  SG-filter by true client IP on targets.[^21][^22]
- AZs: the service is consumable only in AZs the NLB has enabled — **match by AZ ID
  (`use1-az4`), not AZ name** (names are randomized per account). Run targets in ≥2 AZs; enable
  cross-zone load balancing if you can't put targets in every enabled AZ (costs inter-AZ transfer,
  loses zone isolation).[^19][^23]
- NLB TCP idle timeout defaults to **350 s**, configurable 60–6000 s since 2024 (TLS listeners
  fixed at 350 s; UDP flows 120 s). Keepalive guidance: interval below the idle timeout.[^24]
- One NLB belongs to at most one endpoint service; a service can have several NLBs — each consumer
  endpoint ENI is pinned at first connect to one same-AZ NLB, so keep listener/target config
  identical across NLBs.[^19]

### Cross-account and cross-region

- Cross-account is the default consumption model: the consumer creates an interface endpoint
  against your service name from their own account; no routes, no CIDR coordination.[^19]
- **Cross-region PrivateLink** (verified-as-of: 2026-06-12): GA 2024-11-26 for customer endpoint
  services; extended to AWS-managed services 2025-11-19. Provider adds "supported Regions" to the
  service (no remote infrastructure needed); consumer creates an endpoint in their region naming
  the remote **service region**. IAM-gated via `vpce:AllowMultiRegion` + condition keys
  (`ec2:VpceSupportedRegion`, `ec2:VpceServiceRegion`). Not supported on NLBs with a custom TCP
  idle timeout. Cross-region endpoints get **regional DNS names only** (no zonal); PrivateLink
  manages AZ failover but not cross-region failover.[^25][^26]

### Custom private DNS name for a service

- A provider can attach **one** private DNS name (wildcards allowed, e.g. `*.example.com`) so
  consumers keep calling `api.example.com`. Requires **domain-ownership verification via a TXT
  record** (AWS issues name/value, state `pendingVerification` → `verified`). If verification
  lapses, new connections are denied but existing ones survive. Not supported for GWLB
  services.[^27]
- Consumer side: enabling private DNS on the endpoint makes AWS create a private hosted zone in
  the **consumer's** VPC CNAME-ing the provider's name to the endpoint's regional DNS name.[^27][^16]

### VPC Lattice vs PrivateLink (one paragraph)

PrivateLink is L4, needs a provider-run NLB/GWLB, gates by allow-list + acceptance, and is the
cheaper, simpler primitive for exposing one TCP/UDP service to known consumers. VPC Lattice is
L7-capable (HTTP/gRPC routing, weighted targets), needs no load balancer, uses IAM auth policies,
also tolerates overlapping CIDRs, and fits multi-account microservice meshes. Since 2024-12 the two
share infrastructure (VPC resources / Resource Gateway). Prefer PrivateLink for non-HTTP protocols,
SaaS-to-known-customers, and cost; Lattice for service-mesh-shaped problems.[^28]

## Private DNS and split-horizon resolution

### The three name layers

1. **Endpoint-specific names** (always created): regional
   `{endpoint-id}.{service-id}.{region}.vpce.amazonaws.com` and zonal
   `{endpoint-id}-{az}.{service-id}.{region}.vpce.amazonaws.com`. These records live in **public
   DNS but resolve to the endpoint's PRIVATE IPs** — even from the internet. That is deliberate:
   on-prem clients can use `vpce-…` names with zero DNS plumbing (the IPs are useless without
   network reach).[^2][^29]
2. **The private DNS overlay** (the split-horizon mechanism, optional but default-on for AWS
   services): enabling "private DNS" creates a hidden AWS-managed **Route 53 private hosted zone**
   that makes the service's normal public name (e.g. `sqs.us-east-1.amazonaws.com`) resolve to the
   endpoint's private IPs **inside that VPC only**. Outside the VPC the same name still resolves
   publicly — textbook split-horizon, no client changes needed.[^2][^16][^30]
3. **Provider custom private DNS names** (see provider section above) — same trick, consumer-side
   PHZ created per endpoint.[^27]

### Prerequisites and resolution path

- Private DNS requires VPC attributes **`enableDnsSupport` AND `enableDnsHostnames` = true**, and
  the service must support it (GWLBEs don't; DynamoDB exposes no private-DNS toggle and AWS warns
  against DIY PHZ overrides of DynamoDB names).[^31][^3]
- Queries must reach the **Route 53 Resolver** at VPC CIDR+2 (or 169.254.169.253). It only answers
  from inside the VPC. Custom in-VPC DNS servers must forward to it. Precedence trap: a Resolver
  **forwarding rule for the same/more-specific domain beats the PHZ** — a classic accidental
  bypass of the overlay.[^31][^32]

### Hybrid (on-prem) resolution

- On-prem resolvers cannot see the overlay. Standard fix: **Route 53 Resolver inbound endpoint**
  in the VPC + conditional forwarder on the on-prem DNS for the service domains; queries arriving
  through the inbound endpoint get the private answers. Requires separate L3 reach (DX/VPN,
  optionally via TGW).[^33]
- Alternatives: point on-prem clients at the endpoint-specific `vpce-…` names (publicly
  resolvable → private IPs), or maintain manual CNAMEs on the on-prem DNS.[^29][^33]

### Peered VPCs / centralized endpoints

- The AWS-managed PHZ binds only to the endpoint's own VPC. Hub-and-spoke pattern: create the
  central endpoint with private DNS **disabled**, build your **own PHZ** for the service name with
  an alias to the endpoint's regional DNS name, associate it with every spoke VPC (cross-account
  association supported) — or attach PHZs to a **Route 53 Profile** (2024) shared via RAM (FIS
  replaced ~6000 per-VPC associations with one Profile per region for 13,000 endpoints).[^33][^34][^35]
- Failure mode: PHZ associated with a spoke that has no route to the endpoint ENIs — resolution
  succeeds, connection hangs.[^16]
- AWS-managed private DNS doesn't span multiple endpoints for one service; scale-out = private DNS
  off + custom PHZ with weighted alias records.[^36]

### Gateway endpoints: no DNS magic

Gateway endpoints change **routing, not DNS**: `dig s3.us-east-1.amazonaws.com` keeps returning
public IPs and that is correct — the prefix-list route intercepts the traffic. Seeing public IPs in
DNS output is NOT evidence the gateway endpoint is broken.[^8][^9]

### S3 private-DNS nuance

S3 interface endpoints gained private DNS in 2023-03, including the S3-only option **"enable
private DNS only for inbound resolver endpoint"** (`PrivateDnsOnlyForInboundResolverEndpoint`):
on-prem queries via the inbound resolver get interface-endpoint private IPs while in-VPC queries
keep resolving to public IPs and ride the free gateway endpoint. Requires an existing S3 gateway
endpoint in the VPC (and blocks its deletion while set).[^15][^9]

## Endpoint policies

- An endpoint policy is an IAM-style **resource policy on the endpoint** (gateway and most
  interface endpoints; not GWLBE; endpoints to non-AWS endpoint services always get full access).
  Default policy = `Allow * * *`.[^37]
- **Policies filter, never grant**: effective access = identity policy ∧ endpoint policy ∧
  resource policy. AWS data-perimeter framing: "you define the maximum allowed access through the
  network."[^37][^38]
- Size limit **20,480 characters, not adjustable**. Changes take effect within minutes. For
  **gateway** endpoints the `Principal` element must be `"*"` — scope with `aws:PrincipalArn`
  conditions instead (`"AWS": "<account-id>"` matches only the root user: documented foot-gun).[^37][^9]
- Not every service supports endpoint policies (then the endpoint silently allows full access).
  Verify live, per region:
  `aws ec2 describe-vpc-endpoint-services --service-name com.amazonaws.<region>.<svc> --query 'ServiceDetails[*].VpcEndpointPolicySupported'`.[^37]
- **Data-perimeter pattern**: endpoint policy allows only trusted identities from your network —
  `Allow` conditioned on `aws:PrincipalOrgID` (+ `aws:PrincipalIsAWSService` so service principals
  keep working); add `aws:ResourceOrgID` so requests through your endpoint can only target
  org-owned resources (anti-exfiltration: foreign credentials writing to a foreign bucket through
  your endpoint). The mirror-image network controls (`aws:SourceVpce`, `aws:SourceVpc`,
  `aws:VpceOrgID`, `aws:VpceAccount`) live in bucket policies/SCPs/RCPs, not the endpoint
  policy.[^38][^39]
- **Carve-out pitfalls** (the classic breakages):
  - Locked-down S3 endpoint policies break **ECR/ECS image pulls** — layers come from the
    AWS-owned bucket `arn:aws:s3:::prod-<region>-starport-layer-bucket/*`; allow `s3:GetObject`
    on it (403 `CannotPullContainerError` otherwise). Org-ID conditions also fail for these
    presigned fetches.[^40]
  - **SSM** needs AWS-managed public buckets (agent packages, `<region>-birdwatcher-prod`, patch
    baselines).[^41]
  - Start from **aws-samples/data-perimeter-policy-examples**, which ships per-service endpoint
    policies with the carve-outs pre-baked.[^39]

## Quotas, scaling, pricing

verified-as-of: 2026-06-12 (us-east-1 prices; check the pricing page for other regions)

| Item | Value |
| --- | --- |
| Interface endpoint | $0.01 per endpoint per AZ per hour + data processing $0.01/GB (first 1 PB/mo/region; $0.006 next 4 PB; $0.004 beyond — tiers pool regionally)[^42] |
| Gateway endpoint (S3/DynamoDB) | Free[^8] |
| GWLB endpoint | $0.01/hr + $0.0035/GB processed (AWS ELB pricing examples)[^43] |
| Resource endpoint | $0.01/AZ/hr + tiered per-GB + $0.02 per resource per hour[^42] |
| Cross-region | consumer additionally pays inter-region data transfer both directions (~$0.02/GB each way in AWS's example); provider pays ~$0.05/hr per active remote region[^42][^25] |
| Bandwidth | 10 Gbps per AZ baseline, auto-scales to 100 Gbps per AZ; beyond that, multiple endpoints + weighted DNS[^5][^36] |
| Interface + GWLB endpoints per VPC | 50 combined (adjustable, L-29B6F2EB) — the oft-quoted "50 interface endpoints" is actually the combined quota[^5] |
| Gateway endpoints | 20 per region default (adjustable), hard cap 255 per VPC[^5] |
| Endpoint policy size | 20,480 chars, not adjustable[^5] |
| MTU | 8500 bytes, larger dropped, no PMTUD, MSS clamping enforced[^5] |
| Idle timeouts | interface endpoint/NLB 350 s (NLB configurable 60–6000 s since 2024; not compatible with cross-region); GWLB TCP 350 s default configurable; UDP/non-TCP 120 s[^6][^24][^12] |

Cost framing vs NAT Gateway ($0.045/hr + $0.045/GB): gateway endpoints make S3/DynamoDB traffic
free; interface endpoints cut per-GB ~78% for chatty services (ECR, CloudWatch, STS). Counterweight:
endpoint-hours multiply as endpoints × AZs × VPCs (50 VPCs × 18 endpoints × 3 AZs ≈ $19.7k/mo before
data) — which is exactly why the centralized-endpoint pattern below exists.[^44][^35]

## Decision framework: PrivateLink vs VPC peering vs Transit Gateway

| Dimension | PrivateLink | VPC peering | Transit Gateway |
| --- | --- | --- | --- |
| Direction | Unidirectional, consumer→service, specific ports | Bidirectional, full network | Bidirectional, transitive hub |
| CIDR overlap | **Tolerated** (no routing exchanged) | Must not overlap | Must not overlap |
| Scale | Thousands of consumers per service | Non-transitive; ~125 peerings/VPC; full mesh = n(n−1)/2 | Up to 5,000 attachments, 50 Gbps/attachment |
| Exposure | One service, narrowest blast radius | Whole CIDR (subject to routes/SGs) | Whole CIDR via hub |
| Cost shape (us-east-1, verified-as-of: 2026-06-12) | $0.01/AZ/hr + $0.01/GB | No hourly; same-AZ data free (since 2021-05), cross-AZ $0.01/GB each side | ~$0.05/attachment/hr + ~$0.02/GB processed (regional variance, e.g. $0.07 ap-southeast-2) |
| Latency | + endpoint/NLB hop | Lowest (direct) | + TGW hop |
| Protocols | TCP + UDP (L4) | Any IP | Any IP |

Heuristic: need full network reachability? **No → PrivateLink.** Yes → exactly 2 VPCs: **peering**;
3+ VPCs / hybrid / transitive routing / centralized inspection: **TGW**. Exposing one service to
many parties, especially cross-org with CIDR collisions: **PrivateLink** — it is the standard SaaS
private-connectivity mechanism (per-consumer allow-list + acceptance; AWS Prescriptive Guidance
rates it lowest-TCO for SaaS network access).[^45][^46][^47]

**When NOT to use PrivateLink**: general-purpose bidirectional reachability; many services to
expose (one endpoint service + NLB each); very high throughput where per-GB processing beats
peering/TGW economics; protocols beyond TCP/UDP.[^46][^48]

### Combining: centralized interface endpoints (hub-and-spoke)

Place interface endpoints once in a shared-services VPC; spokes reach them via TGW. Required DNS
wiring: endpoints created with private DNS **off**, a custom PHZ per service name aliased to the
endpoint, associated with every spoke (or one Route 53 Profile). Cost trade: 1× endpoint-hours vs
TGW $0.02/GB on all endpoint traffic — wins with many VPCs and modest data volume, loses for
high-throughput VPCs. On-prem reaches centralized **interface** endpoints via DX/VPN→TGW (gateway
endpoints can't be centralized at all). The shared-services VPC becomes critical infrastructure —
a bad PHZ change ripples everywhere.[^34][^35][^33]

Adjacent options: **VPC sharing (RAM)** — multiple accounts launch into one shared VPC's subnets,
no inter-VPC connectivity needed at all, coarser isolation; **VPC Lattice** — L7 service network
across accounts, also CIDR-overlap tolerant.[^28]

## Troubleshooting connectivity through endpoints

### Ordered checklist

1. **DNS**: `dig <hostname>` from the client. Private-path IP (in the VPC CIDR) expected for
   interface endpoints with private DNS. Public IP ⇒ private DNS off, VPC DNS attributes off,
   wrong resolver, or (on-prem) no inbound resolver endpoint. For gateway endpoints a public IP
   is CORRECT — check the route table instead.[^18][^31]
2. **Endpoint SG**: inbound 443 (service port) from client SG/CIDR? Default-SG trap (see SG
   section). Client egress allows the endpoint?[^17][^18]
3. **NACLs**: both subnets, both directions, ephemeral 1024–65535 return ports.[^17]
4. **Endpoint policy**: restrictive policy silently blocks specific API calls with the network
   path fine (S3→ECR starport bucket being the classic).[^37][^40]
5. **Route table** (gateway/GWLBE only): prefix-list route present in the subnet's table?[^8]
6. **Endpoint state**: must be `Available`; `pendingAcceptance`/`Rejected` = connect timeouts.[^49]
7. **AZ mismatch**: verbatim error *"The VPC endpoint service … does not support the availability
   zone of the subnet: subnet-…"* — provider has no NLB node in that AZ **and** AZ names differ
   per account; compare **AZ IDs** (`describe-vpc-endpoint-services` lists supported AZs).[^23]
8. **From outside the VPC**: interface endpoints are reachable over peering/TGW/DX given a route +
   resolvable DNS; **gateway endpoints never are**.[^9]
9. **Provider side**: NLB target health (endpoint can look green while all targets are unhealthy),
   listener port match, NLB SG.[^49]

### Tools

- **VPC Reachability Analyzer** — static config analysis (no packets); supports endpoints,
  endpoint services, peering, TGW as path components; returns the blocking component with an
  explanation code. Same-region only; endpoints can't be a connection *initiator* (consistent with
  PrivateLink unidirectionality).[^50]
- **VPC Flow Logs on the endpoint ENIs** — confirm packets arrive/return. GWLB↔appliance legs are
  GENEVE-encapsulated; inner flows aren't visible there.[^51]
- **CloudWatch PrivateLink metrics** (1-min, free, since 2022-01): namespace
  `AWS/PrivateLinkEndpoints` — `ActiveConnections`, `NewConnections`, `BytesProcessed`,
  `PacketsDropped`, `RstPacketsReceived`; namespace `AWS/PrivateLinkServices` (provider) — same
  plus `RstPacketsSent`, `EndpointsCount`. Not published for gateway endpoints.[^52]
- **CloudTrail** — endpoint create/modify/delete are EC2 API actions; answers "who changed the
  policy/SGs/subnets".

### Failure-mode catalog

| Symptom | Likely cause / fix |
| --- | --- |
| Long-lived connection dies silently; RST on next write; rising `RstPacketsReceived` | 350 s idle timeout (endpoint or NLB). TCP keepalives < 350 s; or raise NLB `tcp.idle_timeout.seconds`[^6][^24] |
| Provider app sees wrong client IP / SG filtering by client IP fails | PrivateLink never preserves client IP — enable PPv2 on the target group AND parse it in the app[^21][^22] |
| Service hostname resolves to public IPs in-VPC | Private DNS disabled; `enableDnsHostnames`/`enableDnsSupport` off; resolver bypass; forwarding rule shadowing the PHZ[^18][^31][^32] |
| On-prem resolves public IPs | No Route 53 inbound resolver endpoint + conditional forwarding; or use `vpce-…` names directly[^33] |
| New endpoint times out for everyone | Default SG selected at creation (self-referencing inbound only)[^17][^18] |
| `CreateVpcEndpoint` AZ error | AZ-ID mismatch vs provider's NLB AZs[^23] |
| Container pulls fail 403 through locked-down S3 endpoint | Allow `s3:GetObject` on `prod-<region>-starport-layer-bucket/*`[^40] |
| Intermittent drops through GWLB inspection | Asymmetric AZ routing — enable TGW appliance mode / per-AZ route design; appliance NAT breaking 5-tuple match[^13][^53] |
| Large packets vanish | >8500 MTU dropped, no PMTUD; clamp MSS / cap MTU[^5] |
| `RequestLimitExceeded` on endpoint APIs | EC2-API throttling at account AND organization level[^5] |

## Azure Private Link / GCP Private Service Connect equivalence

Short comparative map only — deep multi-cloud is out of scope (Atlas-specific Azure/GCP private
connectivity lives in `mongodb-atlas-azure` / `mongodb-atlas-gcp`).

| AWS | Azure | GCP |
| --- | --- | --- |
| Interface VPC endpoint (ENI) | Private Endpoint (NIC in VNet subnet) | PSC endpoint (forwarding rule + internal IP) |
| Endpoint service ← NLB/GWLB | Private Link service ← Standard Load Balancer | Service attachment ← internal LB (+ NAT subnet) |
| Allowed principals + accept/reject | Auto-approval list by subscription, else manual | Consumer accept/reject lists by project |
| Private-DNS flag / custom PHZ | `privatelink.*` private DNS zones | Cloud DNS private zones (auto-records if producer sets a DNS domain) |

All three share the model: unidirectional consumer→service, provider-backbone transport, explicit
per-consumer authorization, overlapping-CIDR tolerance.[^54][^55][^56]

## Anti-patterns

- **Allowing principal `*` with auto-accept** on an endpoint service — your "private" NLB is now
  reachable by any AWS account.[^20]
- **Treating "DNS shows public IPs" as a gateway-endpoint failure** — that's how gateway endpoints
  work; check the route table.[^8]
- **Forgetting the endpoint SG entirely** (default SG) or restricting client egress so hard DNS
  breaks.[^17][^18]
- **Org-scoped endpoint policies without AWS-owned-bucket carve-outs** (ECR starport, SSM
  birdwatcher) — breaks image pulls and patching in subtle 403s.[^40][^41]
- **Expecting bidirectional traffic over PrivateLink** — it is consumer→service only; for
  east-west reachability use peering/TGW.[^45]
- **Centralizing gateway endpoints** — impossible; per-VPC only.[^9]
- **Relying on AZ names across accounts** — always compare AZ IDs for endpoint/AZ planning.[^23]
- **Pinning zonal DNS names without your own failover** — only the regional name health-checks
  across AZs.[^29]
- **Long-lived idle connections with no keepalive** through any PrivateLink/NLB path.[^6]

## Cross-references

- **Atlas consumer side**: `mongodb-aws-networking` (Atlas + AWS PrivateLink/peering),
  `mongodb-atlas-azure` (Azure Private Link for Atlas), `mongodb-atlas-gcp` (PSC for Atlas) —
  under the `mongodb-atlas-expert` hub. TCP-keepalive-vs-PrivateLink-timeout cases:
  MongoDB KB 000023037.
- **This hub**: `references/aws-core.md` (VPC design, SG vs NACL basics, NAT, Well-Architected),
  `references/aws-serverless.md` (Lambda-in-VPC patterns that consume endpoints).
- **networking hub**: DNS protocol internals (resolution path, split-horizon/BIND views, TTLs).
- **devops-linux-admin**: host-level diagnostics (`dig`, `ss`, `tcpdump`) used in the checklist.

## References

[^1]: AWS PrivateLink FAQs — https://aws.amazon.com/privatelink/faqs/
[^2]: Access AWS services through AWS PrivateLink — https://docs.aws.amazon.com/vpc/latest/privatelink/privatelink-access-aws-services.html
[^3]: Configure an interface endpoint / create — https://docs.aws.amazon.com/vpc/latest/privatelink/interface-endpoints.html and https://docs.aws.amazon.com/vpc/latest/privatelink/create-interface-endpoint.html
[^4]: UDP support for AWS PrivateLink (2024-10-31) — https://aws.amazon.com/about-aws/whats-new/2024/10/aws-udp-privatelink-dual-stack-network-load-balancers/
[^5]: AWS PrivateLink quotas — https://docs.aws.amazon.com/vpc/latest/privatelink/vpc-limits-endpoints.html
[^6]: Implementing long-running TCP connections within VPC networking — https://aws.amazon.com/blogs/networking-and-content-delivery/implementing-long-running-tcp-connections-within-vpc-networking/
[^7]: Data transfer price reduction for PrivateLink, TGW, Client VPN (2022-04) — https://aws.amazon.com/about-aws/whats-new/2022/04/aws-data-transfer-price-reduction-privatelink-transit-gateway-client-vpn-services/
[^8]: Gateway endpoints — https://docs.aws.amazon.com/vpc/latest/privatelink/gateway-endpoints.html
[^9]: Gateway endpoints for Amazon S3 — https://docs.aws.amazon.com/vpc/latest/privatelink/vpc-endpoints-s3.html
[^10]: Gateway Load Balancer endpoints — https://docs.aws.amazon.com/vpc/latest/privatelink/gateway-load-balancer-endpoints.html and https://docs.aws.amazon.com/vpc/latest/privatelink/vpce-gateway-load-balancer.html
[^11]: Gateway Load Balancers (ELB docs) — https://docs.aws.amazon.com/elasticloadbalancing/latest/gateway/gateway-load-balancers.html
[^12]: Configurable TCP idle timeout for GWLB — https://aws.amazon.com/blogs/networking-and-content-delivery/introducing-configurable-tcp-idle-timeout-for-gateway-load-balancer/
[^13]: Best practices for deploying Gateway Load Balancer — https://aws.amazon.com/blogs/networking-and-content-delivery/best-practices-for-deploying-gateway-load-balancer/
[^14]: Access VPC resources through AWS PrivateLink — https://docs.aws.amazon.com/vpc/latest/privatelink/privatelink-access-resources.html and https://aws.amazon.com/about-aws/whats-new/2024/12/access-vpc-resources-aws-privatelink/
[^15]: Private DNS support for Amazon S3 with AWS PrivateLink — https://aws.amazon.com/blogs/storage/introducing-private-dns-support-for-amazon-s3-with-aws-privatelink/
[^16]: How AWS PrivateLink and Private DNS actually work (Kevin Keller) — https://kevinkeller.org/posts/how-privatelink-and-dns-work/
[^17]: Security group and NACL settings for VPC endpoints (re:Post KC) — https://repost.aws/knowledge-center/security-network-acl-vpc-endpoint
[^18]: Fix gateway or interface endpoint connectivity (re:Post KC) — https://repost.aws/knowledge-center/vpc-fix-gateway-or-interface-endpoint
[^19]: Create a service powered by AWS PrivateLink — https://docs.aws.amazon.com/vpc/latest/privatelink/create-endpoint-service.html
[^20]: Configure an endpoint service — https://docs.aws.amazon.com/vpc/latest/privatelink/configure-endpoint-service.html
[^21]: Edit target group attributes — client IP preservation / PPv2 — https://docs.aws.amazon.com/elasticloadbalancing/latest/network/edit-target-group-attributes.html
[^22]: Preserving client IP in PrivateLink with proxy protocol v2 (zoph.me) — https://zoph.me/posts/2024-08-18-proxy-protocol-privatelink/ and https://aws.amazon.com/blogs/networking-and-content-delivery/preserving-client-ip-address-with-proxy-protocol-v2-and-network-load-balancer/
[^23]: Interface endpoint AZ errors (re:Post KC) — https://repost.aws/knowledge-center/interface-endpoint-availability-zone
[^24]: Network Load Balancers — idle timeout — https://docs.aws.amazon.com/elasticloadbalancing/latest/network/network-load-balancers.html
[^25]: Introducing cross-region connectivity for AWS PrivateLink — https://aws.amazon.com/blogs/networking-and-content-delivery/introducing-cross-region-connectivity-for-aws-privatelink/ and https://aws.amazon.com/about-aws/whats-new/2024/11/aws-privatelink-across-region-connectivity/
[^26]: PrivateLink cross-region for AWS services (2025-11) — https://aws.amazon.com/about-aws/whats-new/2025/11/aws-privatelink-cross-region-connectivity-aws-services/ and https://docs.aws.amazon.com/vpc/latest/privatelink/privatelink-share-your-services.html
[^27]: Manage DNS names for VPC endpoint services — https://docs.aws.amazon.com/vpc/latest/privatelink/manage-dns-names.html
[^28]: Exploring Amazon VPC Lattice (One Cloud Please) — https://onecloudplease.com/blog/exploring-amazon-vpc-lattice ; AWS VPC connectivity comparison — https://cloudrps.com/blog/aws-vpc-connectivity-peering-transit-gateway-privatelink/
[^29]: AWS PrivateLink for Amazon S3 — https://docs.aws.amazon.com/AmazonS3/latest/userguide/privatelink-interface-endpoints.html ; HA endpoint services whitepaper — https://docs.aws.amazon.com/whitepapers/latest/aws-privatelink/creating-highly-available-endpoint-services.html
[^30]: AWS PrivateLink and VPC endpoints complete guide (hidekazu-konishi) — https://hidekazu-konishi.com/entry/aws_privatelink_vpc_endpoints_complete_guide.html
[^31]: Amazon DNS concepts (enableDnsSupport/Hostnames, VPC+2 resolver) — https://docs.aws.amazon.com/vpc/latest/userguide/AmazonDNS-concepts.html ; configure private DNS (re:Post KC) — https://repost.aws/knowledge-center/vpc-interface-configure-dns
[^32]: Considerations for private hosted zones — https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/hosted-zone-private-considerations.html
[^33]: Integrating Transit Gateway with PrivateLink and Route 53 Resolver — https://aws.amazon.com/blogs/networking-and-content-delivery/integrating-aws-transit-gateway-with-aws-privatelink-and-amazon-route-53-resolver/
[^34]: Centralized access to VPC private endpoints (multi-VPC whitepaper) — https://docs.aws.amazon.com/whitepapers/latest/building-scalable-secure-multi-vpc-network-infrastructure/centralized-access-to-vpc-private-endpoints.html ; Centralize access using VPC interface endpoints — https://aws.amazon.com/blogs/networking-and-content-delivery/centralize-access-using-vpc-interface-endpoints/
[^35]: How FIS centralized 13,000 VPC endpoints (Route 53 Profiles) — https://aws.amazon.com/blogs/networking-and-content-delivery/how-fis-centralized-13000-vpc-endpoints-to-strengthen-security-and-simplify-operations/ ; Streamline DNS for PrivateLink with Route 53 Profiles — https://aws.amazon.com/blogs/networking-and-content-delivery/streamline-dns-management-for-aws-privatelink-deployment-with-amazon-route-53-profiles/
[^36]: Scale traffic using multiple interface endpoints — https://aws.amazon.com/blogs/networking-and-content-delivery/scale-traffic-using-multiple-interface-endpoints/
[^37]: Control access to VPC endpoints using endpoint policies — https://docs.aws.amazon.com/vpc/latest/privatelink/vpc-endpoints-access.html
[^38]: Establishing a data perimeter on AWS: allow only trusted identities — https://aws.amazon.com/blogs/security/establishing-a-data-perimeter-on-aws-allow-only-trusted-identities-to-access-company-data/ ; IAM data perimeters — https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_data-perimeters.html
[^39]: aws-samples/data-perimeter-policy-examples — https://github.com/aws-samples/data-perimeter-policy-examples
[^40]: ECR image pull errors through restricted S3 endpoints (re:Post KC) — https://repost.aws/knowledge-center/ecs-ecr-docker-image-error ; ECR interface VPC endpoints — https://docs.aws.amazon.com/AmazonECR/latest/userguide/vpc-endpoints.html
[^41]: SSM Agent technical details (AWS-managed buckets) — https://docs.aws.amazon.com/systems-manager/latest/userguide/ssm-agent-technical-details.html
[^42]: AWS PrivateLink pricing — https://aws.amazon.com/privatelink/pricing/
[^43]: Elastic Load Balancing pricing (GWLBE examples) — https://aws.amazon.com/elasticloadbalancing/pricing/
[^44]: VPC endpoint cost analysis — https://awsnegotiations.com/blog/vpc-endpoint-cost-analysis ; NAT Gateway vs VPC endpoints (Michal Drozd) — https://www.michal-drozd.com/en/blog/aws-nat-gateway-vs-vpc-endpoints/
[^45]: VPC peering pricing change (same-AZ free, 2021-05) — https://aws.amazon.com/about-aws/whats-new/2021/05/amazon-vpc-announces-pricing-change-for-vpc-peering/ ; peering vs PrivateLink (StackHarbor) — https://stackharbor.com/en/knowledge-base/awsvpc-peering-vs-privatelink/
[^46]: AWS PrivateLink vs VPC peering (ngrok) — https://ngrok.com/blog/aws-privatelink-vs-vpc-peering
[^47]: SaaS network access options (AWS Prescriptive Guidance) — https://docs.aws.amazon.com/prescriptive-guidance/latest/saas-network-access-options/evaluating.html ; PrivateLink use cases whitepaper — https://docs.aws.amazon.com/whitepapers/latest/aws-privatelink/use-case-examples.html
[^48]: Transport-layer tenant routing for SaaS using AWS PrivateLink — https://aws.amazon.com/blogs/networking-and-content-delivery/transport-layer-tenant-routing-for-saas-using-aws-privatelink/
[^49]: Troubleshoot endpoint-to-endpoint-service connectivity (re:Post KC) — https://repost.aws/knowledge-center/connect-endpoint-service-vpc
[^50]: How Reachability Analyzer works — https://docs.aws.amazon.com/vpc/latest/reachability/how-reachability-analyzer-works.html
[^51]: Troubleshoot GWLB connectivity (re:Post KC) — https://repost.aws/knowledge-center/elb-gwlb-connectivity-issues
[^52]: CloudWatch metrics for AWS PrivateLink — https://docs.aws.amazon.com/vpc/latest/privatelink/privatelink-cloudwatch-metrics.html ; usage insights blog — https://aws.amazon.com/blogs/networking-and-content-delivery/gain-usage-insights-with-amazon-cloudwatch-metrics-for-aws-privatelink/
[^53]: VPC routing enhancements and GWLB deployment patterns — https://aws.amazon.com/blogs/networking-and-content-delivery/vpc-routing-enhancements-and-gwlb-deployment-patterns/
[^54]: Azure for AWS professionals — networking — https://learn.microsoft.com/en-us/azure/architecture/aws-professional/networking
[^55]: Azure Private Link service overview — https://learn.microsoft.com/en-us/azure/private-link/private-link-service-overview
[^56]: GCP Private Service Connect — https://docs.cloud.google.com/vpc/docs/private-service-connect
[^57]: Choosing your VPC endpoint strategy for Amazon S3 — https://aws.amazon.com/blogs/architecture/choosing-your-vpc-endpoint-strategy-for-amazon-s3/
[^58]: PrivateLink network performance troubleshooting (re:Post KC) — https://repost.aws/knowledge-center/vpc-troubleshoot-network-performance-privatelink
