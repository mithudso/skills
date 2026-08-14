<!-- hub-reference-banner -->
> **Reference file — part of the `devops-containers-cicd` hub.** New `/dr`-researched reference (concept-family-explorer run, 2026-06-15).
> Sibling topics in this family are reference files under the hubs — **not** standalone skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling skill; load that topic's `references/<name>.md` from the owning hub.

---

---
name: self-healing-systems-autonomic-computing
title: Self-Healing Systems & Autonomic Computing
description: >
  Systems-level self-healing and autonomic computing: IBM autonomic computing (Kephart & Chess four self-* properties) and the MAPE-K control loop (Monitor-Analyze-Plan-Execute over shared Knowledge); self-healing infrastructure (Kubernetes level-triggered reconciliation, controllers/operators, liveness/readiness/startup probes, replica/node/PV recovery); closed-loop AIOps remediation (anomaly->RCA->playbook->action->verify), the L0-L5 remediation maturity ladder and runbook automation (StackStorm/Rundeck/Ansible); resilience self-healing patterns (circuit breaker, bulkhead, timeouts, retries+backoff+jitter); and safety-gated autonomous remediation (dry-run, approval, checkpoint/rollback, Faithful Undo, Transactional No-Regression; agentic-SRE frontier STRATUS/OpsAgent/AIOpsLab). TRIGGER: "autonomic computing / four self-* properties", "MAPE-K loop", "how does Kubernetes self-healing / reconciliation work", "closed-loop auto-remediation", "runbook automation", "remediation maturity ladder", "safety gate / rollback for automated remediation", "agentic SRE / STRATUS". SKIP: incident PROCESS/comms, SLO/error-budget DEFINITION, on-call, postmortems -> tam-operations / incident-comms / postmortem-writing; telemetry/OTel/eBPF MONITORING -> devops-observability; application-level resilience pattern code (retry/circuit-breaker/fallback) -> autoremediation; Kubernetes workload/CNI/mesh authoring -> kubernetes-networking; code-level patch synthesis (APR, GenProg) -> automated-program-repair; agent/compound-system self-improvement -> ai-agents-orchestration.
category: developer
---

# Self-Healing Systems & Autonomic Computing

## Overview

This skill is the **systems-level** reference for software and infrastructure that detects its own faults and acts to restore a healthy state without a human in the loop — the closed-loop *Monitor -> decide -> act -> verify* discipline. It runs from the founding theory (IBM autonomic computing, the MAPE-K control loop) through concrete infrastructure mechanisms (Kubernetes reconciliation, probes, operators), AIOps closed-loop remediation, resilience patterns, and the 2025-26 safety-gated agentic-SRE frontier.

It is deliberately scoped *away* from three neighbors: code-level patch synthesis (APR) lives in `automated-program-repair`; the *implementation* of application resilience patterns (retry/circuit-breaker/fallback code, thresholds, idempotency keys) and telemetry instrumentation live in `devops-observability`; and incident *process* (comms, SLO definition, on-call, postmortems) lives in `tam-operations` / `incident-comms` / `postmortem-writing`. This skill owns the **architecture and mechanism** of automated healing.

## When to Use This Skill

| Task | Examples |
|---|---|
| Apply autonomic theory | four self-* properties, autonomic element = managed element + manager |
| Reason about a control loop | MAPE-K phases over shared Knowledge, async messaging bus |
| Explain infra self-healing | Kubernetes reconcile loop, probes, replica/node/PV recovery |
| Design closed-loop remediation | anomaly -> RCA -> playbook -> action -> verify -> audit |
| Place a remediation maturity level | L0 alerting through L5 autonomous + verified rollback |
| Add safety gating | dry-run, approval, checkpoint/rollback, Faithful Undo, TNR |
| Evaluate the agentic-SRE frontier | STRATUS, OpsAgent, MicroRemed, AIOpsLab, ITBench |

Not this skill: the *consensus algorithm* a self-healing replicated system relies on — Raft/Paxos leader election, quorum intersection, the CAP/FLP theory behind why a reconcile loop can lose liveness under partition — is `distributed-systems-consensus` (control loops here; agreement protocols there).

## Core Concepts

1. **Autonomic computing & the four self-* properties.** Kephart & Chess (*The Vision of Autonomic Computing*, IEEE Computer 2003) defined self-managing systems given high-level admin objectives, exhibiting **self-configuration, self-optimization, self-healing, self-protection**. The biological / autonomic-nervous-system framing (the paper's analogy is heart rate and body temperature regulating without conscious control) — system-level behavior emerging from interacting elements — is the conceptual root of all later self-healing work.

2. **The autonomic element.** The architectural unit (from IBM's *Architectural Blueprint for Autonomic Computing*, not the Kephart & Chess vision paper): one or more *managed elements* controlled by a single *autonomic manager* that observes and acts via **sensors and effectors**. System-level behavior emerges from many interacting autonomic elements.

3. **The MAPE-K control loop.** The canonical autonomic-manager internals (IBM *Architectural Blueprint* / Redbook *Problem Determination Using Self-Managing Autonomic Technology*): **Monitor** (collect/aggregate/filter metrics & topology), **Analyze** (correlate/model/predict), **Plan** (policy-guided action structuring), **Execute** (controlled plan execution), over a shared **Knowledge** base. The four parts share an *asynchronous messaging bus*, not a strict control flow — Plan can ask Monitor for more data; Monitor can trigger replanning. Instrumentation-agnostic (SNMP, JMX, DMTF/CIM, web-services sensors/effectors).

4. **Kubernetes reconciliation as MAPE-K Execute.** Self-healing in Kubernetes is controllers continuously reconciling *actual* state toward the desired `.spec`. Controllers are **level-triggered, idempotent reconcile loops**: watch -> workqueue -> `Reconcile()` -> close the gap -> requeue with exponential backoff (Kubernetes docs, *Controllers*, kubernetes.io/docs/concepts/architecture/controller/). This is the systems-level analog of the MAPE-K loop.

5. **Kubernetes self-healing primitives.** Concrete heal mechanisms: container restarts (`restartPolicy` + **liveness probe** restarts deadlocked containers), **readiness probe** gates traffic, **startup probe** for slow boots (Kubernetes docs, *Pod Lifecycle — Container probes*); replica replacement via ReplicaSet/StatefulSet/DaemonSet; node-failure reschedule; PersistentVolume reattach; Service endpoint removal of unhealthy pods (Kubernetes docs, *Self-Healing*). Operators extend reconciliation to custom resources. *Workload/CNI/mesh authoring depth routes to `kubernetes-networking`.*

6. **Closed-loop AIOps remediation.** The detect->decide->act->verify production loop: anomaly detection -> correlation/root-cause analysis -> playbook selection -> automated action -> verification -> audit. The action half is what distinguishes self-healing from monitoring. *Telemetry/anomaly-detection instrumentation (OTel, eBPF, Sentry) routes to `devops-observability`.*

7. **Remediation maturity ladder.** L0 alerting -> L1 manual runbooks -> L2 orchestration (StackStorm/Rundeck/Ansible) -> L3 AI-assisted -> L4/L5 fully agentic autonomous with *verified rollback*. The ladder frames how much human gating remains at each level and is the planning lens for adopting auto-remediation safely.

8. **Runbook automation & error-budget-driven triggers.** Codified runbooks (StackStorm rules/workflows, Rundeck jobs, Ansible playbooks) turn recurring alerts into executable, validated actions. Best practice: validate a runbook via chaos test before letting it fire automatically. Error-budget burn-rate can act as the trigger to escalate from alert to automated action. *SLO/error-budget DEFINITION and on-call process route to `tam-operations`/`incident-comms`.*

9. **Resilience / application-level self-healing patterns.** The fault-tolerance mechanisms that complement infra reconciliation: **circuit breaker** (a self-resetting Closed/Open/Half-Open finite-state machine that trips on a failure threshold to prevent cascading failure — Fowler/Nygard *Release It!*, popularized by Hystrix), **bulkheads** (isolate blast radius), **fail-fast**, **timeouts**, and **retries with exponential backoff + full jitter**. Standard chain: breaker -> retry -> timeout -> call, with a per-dependency bulkhead. *Concept reference here; implementation code routes to `devops-observability` / application frameworks.*

10. **Safety-gated autonomous remediation.** The central concern as systems gain autonomy: confidence thresholds, dry-run, approval workflows, checkpoint/rollback, and formal safety specs. STRATUS (NeurIPS 2025, *A Multi-agent System for Autonomous Reliability Engineering of Modern Clouds*) formalizes **Transactional No-Regression (TNR)** (the paper also writes "Non-Regression") and **Faithful Undo** — every automated action carries an undo operator — so a remediation can be safely reversed. This is the gating that makes L4/L5 autonomy acceptable in regulated/financial systems.

11. **The 2025-26 agentic-SRE frontier.** Autonomic self-healing is converging with LLM multi-agent systems at the *infrastructure* level: **STRATUS** (specialized detection/diagnosis/mitigation agents in a safety state machine; >=1.5x mitigation success on AIOpsLab/ITBench), **OpsAgent** (self-evolving, production-deployed), and end-to-end microservice remediation benchmarks (**MicroRemed/E2E-MR**: diagnosis report -> executable Ansible playbook -> auto-execute -> verify; **AIOpsLab**, **ITBench**, **OPENRCA**). *Agent/compound-system self-improvement design routes to `ai-agents-orchestration`.*

12. **Predictive (proactive) self-healing.** The frontier beyond reactive loops: correlate multi-source telemetry (metrics, logs, kernel dmesg, NVMe SMART, thermal) to stage remediation playbooks *before* failures materialize — moving the loop left of the incident.

## References Outline

A full skill would carry these sub-files under `references/`:

- `autonomic-computing.md` — Kephart & Chess, the four self-* properties, the autonomic element (Architectural Blueprint), sensors/effectors.
- `mape-k-loop.md` — Monitor/Analyze/Plan/Execute over Knowledge, async bus, instrumentation-agnostic sensors/effectors, MAPE-K mapped onto modern controllers.
- `kubernetes-self-healing.md` — level-triggered idempotent reconciliation (Controllers doc), probes (Pod Lifecycle doc), replica/node/PV recovery, operators; explicit hand-off to `kubernetes-networking` for authoring depth.
- `aiops-closed-loop.md` — anomaly->RCA->playbook->action->verify->audit, the L0-L5 maturity ladder, runbook automation (StackStorm/Rundeck/Ansible), error-budget triggers.
- `resilience-patterns.md` — circuit breaker FSM, bulkhead, fail-fast, timeouts, retries+backoff+jitter as self-healing concepts; pointer to `devops-observability` for implementation.
- `safety-gated-remediation.md` — dry-run/approval/checkpoint/rollback, Faithful Undo, Transactional No-Regression; STRATUS/OpsAgent/MicroRemed and the AIOpsLab/ITBench benchmarks; predictive self-healing.

## Key Landmarks & Citations

- Kephart & Chess, *The Vision of Autonomic Computing*, IEEE Computer 2003 — the four self-* properties; the autonomic-nervous-system analogy.
- IBM, *An Architectural Blueprint for Autonomic Computing* (and Redbook SG24-6665, *Problem Determination Using Self-Managing Autonomic Technology*) — the autonomic element (managed element + manager, sensors/effectors) and the canonical MAPE-K loop definition.
- Kubernetes docs, *Controllers* (kubernetes.io/docs/concepts/architecture/controller/) — level-triggered reconciliation; *Pod Lifecycle — Container probes* — liveness/readiness/startup; *Self-Healing* — replica/node/PV recovery.
- Fowler, *CircuitBreaker* (martinfowler.com) + Nygard, *Release It!* — Closed/Open/Half-Open FSM, bulkheads, timeouts, retries+jitter.
- STRATUS (NeurIPS 2025) — *A Multi-agent System for Autonomous Reliability Engineering of Modern Clouds*; Transactional No-Regression + Faithful Undo; AIOpsLab/ITBench evaluation.
- Companion references: `automated-program-repair` (code-level patching), `kubernetes-networking`, `devops-observability` (resilience-pattern implementation + telemetry), `tam-operations` (incident process/SLO).
