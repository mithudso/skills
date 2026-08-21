---
name: dns-caching
description: >-
  Highly optimized DNS Caching skill compiled from 35 master concepts.
  Sourced via Wikipedia mapping, IETF RFC standards, and CNCF architectures.
  Distilled via offline pipeline and Ollama semantic indexing.
version: "1.1.0"
updated: "2026-08-19"
category: infrastructure
keywords:
  - dns caching
  - resolver
  - unbound
  - coredns
  - dga detection
  - ambient mesh dns
tags:
  - dns
  - infrastructure
  - caching
  - networking
  - ietf
  - cncf
---
# DNS Caching (v2. Anchor-Expanded)

Comprehensive operational architecture for DNS caching. Concepts grounded in Wikipedia taxonomies, verified against IETF RFC standards, and aligned with CNCF cloud-native patterns. Distilled and semantically deduped.

## 1. Core Architecture & Resolvers
- **Enterprise / ISP (RFC 1034/1035)**: Unbound, BIND9, PowerDNS Recursor, Knot Resolver. Heavy focus on root hints and authoritative zone transfers (AXFR/IXFR).
- **Edge / Local**: dnsmasq, systemd-resolved. Critical for OS/Browser cache interception.
- **Cloud Native (CNCF)**: CoreDNS, NodeLocal DNSCache, Envoy Filter. Designed for dynamic, ephemeral microservice endpoints.

## 2. Operations & Telemetry
- **Tools**: dig (`+trace`), dnsperf, tcpdump (port 53).
- **Logging**: dnstap (binary wire-format logging, superior to text query logs). Prometheus metrics exported via node agents.
- **Administration**: Infrastructure as Code (Terraform), runtime Cache Flushing.

## 3. Resilience, Security & Standards
- **Protocol Standards**: DNSSEC (RFC 4033), QNAME Minimization (RFC 7816 - privacy), Serve Stale (RFC 8767 - resilience against upstream outages).
- **Threat Mitigation**: Cache Poisoning (Kaminsky defense), DNS Amplification blocklists. 
- **Deep Frontier (ML)**: LLM-based DGA Detection via temporal-behavioral shifts, bypassing traditional lexical analysis.
- **Encrypted DNS**: DoH (RFC 8484), DoT (RFC 7858), ODoH (RFC 9230), DDR (Discovery of Designated Resolvers - RFC 9462).

## 4. Advanced Mesh & Routing
- **Global / Edge**: GSLB, Multi-CDN Steering via Real-User Monitoring (RUM) ingestion.
- **Ambient Mesh DNS**: Sidecarless Interception (ztunnel, Istio Ambient).
- **Identity-Aware Resolution**: SPIFFE/SPIRE binding, where the mesh refuses IP resolution without a valid SVID attestation.
