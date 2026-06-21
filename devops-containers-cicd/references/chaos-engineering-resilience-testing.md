<!-- hub-reference-banner -->
> **Reference file — part of the `devops-containers-cicd` hub.** New `/dr`-researched reference (concept-family-explorer frontier run, 2026-06-15).
> Sibling topics in this family are reference files under the hubs — **not** standalone skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling skill; load that topic's `references/<name>.md` from the owning hub.

---

---
name: chaos-engineering-resilience-testing
title: Chaos Engineering & Resilience Testing
description: >
  The testing/validation discipline of deliberately injecting failure to build confidence a distributed system withstands turbulent conditions: the five Principles of Chaos (steady-state hypothesis, vary real-world events, run in production, automate continuously, minimize blast radius); the testing-vs-experimentation distinction; the Netflix lineage (Simian Army -> Chaos Monkey -> FIT/ChAP/Chaos Kong); the hypothesis-driven control/experiment-group loop and GameDays (MTTD/MTTR); fault-injection mechanisms (instance/pod kill, latency/error injection, resource exhaustion, network partition, DNS/IO/time faults); blast-radius control + automated abort/stop conditions; the tooling landscape (Chaos Monkey, Gremlin, Litmus, Chaos Mesh, AWS FIS, Azure Chaos Studio, Steadybit) and Kubernetes-native chaos CRDs; continuous verification / chaos-in-CI with resiliency-score gating; and SRE error-budget governance of how aggressively to experiment. TRIGGER: "design a chaos experiment / steady-state hypothesis", "run a GameDay", "inject latency/pod-kill", "blast radius / stop condition", "Chaos Mesh / LitmusChaos / AWS FIS / Gremlin", "chaos in CI / resiliency-score gate", "test our failover/graceful degradation". SKIP: the remediation MECHANISM (MAPE-K, Kubernetes self-healing, AIOps closed-loop, agentic-SRE) -> self-healing-systems-autonomic-computing; retry/circuit-breaker/fallback CODE -> autoremediation; incident comms/postmortems/SLO process -> incident-comms / postmortem-writing / tam-operations; k8s workload/CNI/mesh authoring -> kubernetes-networking; telemetry/OTel -> devops-observability.
category: developer
---

# Chaos Engineering & Resilience Testing

## Overview

Chaos engineering is "the discipline of experimenting on a distributed system in order to build confidence in the system's capability to withstand turbulent conditions in production" — the canonical definition from the Netflix team that coined the term [Rosenthal et al., O'Reilly / arXiv:1702.05843]. This skill is the **TESTING / validation** half of resilience work: you deliberately inject failure to discover whether a system *does* survive turbulence. The **remediation mechanism** it validates — MAPE-K loops, Kubernetes self-healing, retry/circuit-breaker/fallback code, AIOps closed-loop remediation, agentic-SRE — is owned by `self-healing-systems-autonomic-computing` and `autoremediation`; chaos engineering *exercises* those mechanisms but does not implement them.

The defining mental shift is **testing vs experimentation**: a test makes a binary true/false assertion about a *known* property and generates no new knowledge; chaos engineering is experimentation that surfaces *new* knowledge about emergent behavior in complex systems — "how close is our system to the edge of chaos?" [arXiv:1702.05843, citing Dekker]. Use it to drive a hypothesis-driven loop (define steady state → hypothesize → vary real-world events → try to disprove → minimize blast radius → automate), and to choose and safely operate fault-injection tooling.

## When to Use

- Designing a chaos experiment or writing a falsifiable steady-state hypothesis
- Planning or facilitating a GameDay (human detection/diagnosis/remediation, MTTD/MTTR)
- Selecting a fault-injection mechanism (pod kill, latency, resource exhaustion, network partition, DNS/IO/time faults) for a real-world event
- Choosing a chaos tool (Chaos Monkey, Gremlin, LitmusChaos, Chaos Mesh, AWS FIS, Azure Chaos Studio, Steadybit, Chaos Toolkit, Pumba, Toxiproxy)
- Bounding blast radius and wiring automated abort/stop conditions
- Adding continuous verification / chaos-in-CI with a resiliency-score release gate
- Deciding how aggressively to experiment given the SRE error budget

## When NOT to Use (deferrals)

- Implementing the recovery itself (MAPE-K, K8s reconciliation, AIOps remediation, agentic-SRE) → `self-healing-systems-autonomic-computing`
- Retry/circuit-breaker/fallback/idempotency CODE → `autoremediation`
- Incident comms, postmortems, SLO *process* → `incident-comms` / `postmortem-writing` / `tam-operations`
- Kubernetes workload/CNI/service-mesh authoring → `kubernetes-networking`
- Telemetry/OTel/eBPF instrumentation → `devops-observability`

## Core Concepts

1. **Steady-state hypothesis.** Define normal behavior as a measurable *output* of the system (throughput, error rate, latency percentiles) — not internal attributes. Chaos verifies the system *does* work, not *how* it works. Every experiment and GameDay starts with a written, falsifiable hypothesis; "let's see what happens when we kill a pod" is not one [principlesofchaos.org; Hassan].

2. **Control-group vs experiment-group design.** Hypothesize that steady state continues in both a control group and an experimental group, inject variables only into the experimental group, then try to **disprove** the hypothesis by detecting a steady-state difference. Netflix's ChAP siphons a fraction of real production traffic into paired control/experiment clusters and statistically compares them, which both implements this design and bounds blast radius [principlesofchaos.org; arXiv:1702.05849].

3. **The five Principles of Chaos.** The manifesto lists five equal "Advanced Principles": (1) build a hypothesis around steady-state behavior; (2) vary real-world events (hardware/software failures, traffic spikes, scaling events) prioritized by impact or frequency; (3) run experiments in production; (4) automate experiments to run continuously; (5) **minimize blast radius**. (Blast-radius was a later 2017+ addition to the original four, but the current manifesto presents all five as peers) [principlesofchaos.org].

4. **Testing vs experimentation.** Tests assert a known property (binary, no new knowledge); chaos engineering is experimentation that generates new knowledge about emergent, complex-system behavior. This is why it is distinct from integration testing [arXiv:1702.05843].

5. **The Netflix lineage.** Chaos Monkey randomly terminates production instances during business hours so weaknesses surface at a controlled time (open-sourced 2012; 65k+ instances killed in year one; 2.0 integrated with Spinnaker, terminate-only). The Simian Army adds Latency Monkey (deprecated after it caused cascading failures) and the region-scale Chaos Gorilla; **Chaos Kong** (full-region failover) is a later, separately-documented Netflix exercise. **FIT** (Failure Injection Testing) is the controlled "failure-as-a-service" successor with bounded scope, defined injection points, and a global "Halt all Failures" control [netflixtechblog; arXiv:1702.05843].

6. **ChAP & continuous verification.** ChAP (Chaos Automation Platform) is the orchestration layer on FIT+Hystrix that auto-detects whether a service is resilient instead of relying on a human watching dashboards; its Monocle dashboard surfaces fallbacks/timeouts/retries and a per-service criticality score. ChAP is the seed of **Continuous Verification (CV)** — practitioners typically implement CV as a proactive pipeline stage after CI/CD that compares live behavior to a baseline (often gating or rolling back a deploy). Because confidence in past results decays as the system evolves, experiments are automated and re-run continuously [arXiv:1702.05849; verica.io].

7. **Blast-radius control + automated abort.** Scope to a subset first (one pod, 1% of traffic, one AZ) via selection modes (AWS FIS `selectionMode` COUNT/PERCENT/ALL; Chaos Mesh `mode: fixed | fixed-percent | random-max-percent`); codify the abort path as stop conditions wired to alarms teams already trust (AWS FIS stops on a CloudWatch alarm; a stopped experiment cannot resume), halt buttons, and auto-revert [docs.aws.amazon.com/fis; chaos-mesh.org; gremlin.com].

8. **GameDays.** A structured, time-boxed, multi-team exercise that tests *human* response (detection, diagnosis, remediation) as well as system resilience — distinct from automated experiments that verify *automation*. Separate observer and responder channels; observers hold a kill switch and inject failure at an unannounced time. Metrics: time-to-first-alert, time-to-detection, time-to-correct-diagnosis, remediation correctness, time-to-resolution. Stage the blast radius (single pod → latency into a critical dependency → concurrent multi-failure) [Hassan].

9. **Fault-injection taxonomy.** Map real-world events to mechanisms: instance/pod termination (Chaos Monkey, FIS terminate, `kubectl delete pod`), latency injection (`tc`, NetworkChaos delay) for timeouts/circuit breakers, resource exhaustion (`stress-ng`, StressChaos, `dd`) for OOM/GC/autoscaling, network partition (`iptables`, NetworkChaos partition) for multi-AZ/split-brain, plus DNS/IO/time faults [Hassan; chaos-mesh.org].

10. **Kubernetes-native chaos via CRDs.** Chaos Mesh declares experiments as CRDs (PodChaos, NetworkChaos, StressChaos, DNSChaos, HTTPChaos, IOChaos, TimeChaos, KernelChaos), injected by a privileged Chaos Daemon DaemonSet, with selectors + mode for native blast-radius targeting and RBAC gating. LitmusChaos (CNCF) splits intent into ChaosExperiment (installable fault template), ChaosEngine (binds app→experiment), and ChaosHub (shareable catalog), and validates the steady-state hypothesis with **Probes** (SOT/EOT/Edge/Continuous/OnChaos) that determine the verdict [chaos-mesh.org; litmuschaos.io].

11. **Chaos-in-CI / resiliency-score gating.** Run experiments on every deploy, on a schedule, or as a pipeline gate (CircleCI + Chaos Toolkit; Flux CD GitOps → Chaos Mesh). Harness HCE runs a chaos step in CD and auto-**rolls back** the deploy if the resiliency score is below threshold — turning chaos into a hard quality gate. Resilience decays without continuous re-verification, so this catches a fault-tolerance regression the week it ships [verica.io].

12. **SRE error-budget governance.** Error budget = 1 − SLO is the control loop balancing reliability against velocity: a healthy budget licenses more aggressive chaos; a near-drained budget self-polices toward stability and halts non-critical change. Chaos experiments must themselves run *within* the budget and be bounded by stop conditions. Chaos is the proactive complement to SLO/error-budget governance [sre.google, *Embracing Risk*]; Google's **DiRT** (Disaster Recovery Testing) is a separate, parallel company-wide disaster-exercise program.

## References Outline

- **Principles & definition** — `principlesofchaos.org` (the five Advanced Principles, steady state, experiment design) and Rosenthal et al., *Chaos Engineering* [O'Reilly / arXiv:1702.05843] (canonical definition, testing-vs-experimentation, complex-systems view).
- **Netflix platform lineage** — *The Netflix Simian Army* [netflixtechblog] (Chaos Monkey origin, Simian Army, 2.0/Spinnaker) and *A Platform for Automating Chaos Experiments — ChAP* [arXiv:1702.05849] (traffic-siphoning control/experiment clusters, Monocle, criticality score).
- **Managed & cloud-native tooling** — AWS Fault Injection Service experiment-template docs (Targets/Actions/Stop conditions/roleArn), `chaos-mesh.org` (CRD grammar + Chaos Daemon), `litmuschaos.io` (ChaosExperiment/ChaosEngine/Probes), Gremlin tool-comparison (open-source vs commercial landscape).
- **Practice, CV & governance** — Hassan, *Chaos Engineering in Practice* (GameDay runbook, hypothesis discipline, fault taxonomy, MTTD/MTTR), Verica, *Continuous Verification* (CV as a release gate, resiliency-score gating, decay), Google SRE *Embracing Risk* (`sre.google`) (error budgets).
- **Boundary cross-refs** — defer remediation mechanism to `self-healing-systems-autonomic-computing` / `autoremediation`; incident process to `incident-comms` / `postmortem-writing` / `tam-operations`; k8s authoring to `kubernetes-networking`; telemetry to `devops-observability`.

## Frontier

- **Litmus MCP Server** — exposing chaos experiments to AI/agent (MCP) clients to discover/run/monitor experiments from LLM workflows.
- **Resiliency-score-as-release-gate** (Harness HCE) and broader CV platforms (Verica) extending chaos beyond availability to security/performance/cost.
- **AI/agentic-SRE intersection** — chaos validating autonomous closed-loop remediation (MAPE-K) and AIOps self-healing rather than human-only runbooks (the loop itself stays in `self-healing-systems-autonomic-computing`).
- **Real-outage validation** — large cloud dependency-cascade events reinforcing run-in-production dependency-failure experiments.
- **Game-Day-Continuous** (fault injection on every deploy) and pre-announced "Disasterpiece Theater"-style exercises as organizational practice.
