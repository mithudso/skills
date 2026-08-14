<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `firedrill-integration-tester` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: firedrill-integration-tester
description: >
  Firedrill integration testing patterns — scenario execution, safety mechanisms, abort conditions,
  scoring rubrics, multi-step agent orchestration, and post-drill reporting using the
  mdb_case_assistant MCP firedrill tools.
  TRIGGER: building or running automated firedrill/game-day validation systems, designing
  scenario-based integration tests with multi-step agent workflows, implementing safety mechanisms
  for drill execution, creating scoring rubrics for incident response evaluation, generating
  post-drill scorecards and reports, orchestrating agent-driven test execution with human checkpoints.
  SKIP: real production chaos engineering with live failure injection (use chaos-engineering tools);
  unit testing pure functions or isolated modules (use software-engineering-patterns,
  references/testing-and-vitest-expert.md); load/performance testing (use a dedicated load-testing
  tool such as k6 or Locust); general incident management that is not a simulation (use
  tam-operations, references/incident-response/SKILL.md).
version: 1.2.2
updated: "2026-06-01"
category: developer
tags:
  - firedrill
  - testing
  - integration
  - chaos-engineering
  - scenarios
  - agent
  - game-day
  - incident-simulation
  - scoring-rubric
  - safety
  - mcp
triggers:
  - firedrill
  - game day
  - incident simulation
  - drill scenario
  - firedrill scoring
  - abort drill
  - drill report
  - mdb_case_firedrill
  - agent-driven test orchestration
  - human checkpoint
  - reveal conditions
  - preflight go no-go
related_skills:
  - incident-response
  - testing-and-vitest-expert
  - software-engineering-patterns
  - ai-agent-engineering
---

# Firedrill Integration Tester

Firedrills (also called game days or incident simulations) are controlled exercises that validate incident response procedures, system resilience, and team readiness without real customer impact. This skill covers agent-driven patterns for end-to-end firedrill orchestration — from scenario selection through execution, scoring, and post-drill reporting.

**Key differentiator from general chaos engineering:** firedrills test **human and process** response, not just system behavior. The agent orchestrates the simulation, plays customer/system personas, evaluates response quality, and produces actionable scorecards.

## When to use this skill

- Building automated firedrill execution systems
- Designing scenario-based integration tests with multi-step workflows
- Implementing safety mechanisms (abort, blast radius, timeouts)
- Creating scoring rubrics for incident response evaluation
- Orchestrating agent-driven test execution with human checkpoints
- Generating post-drill reports and remediation tracking

## When NOT to use this skill

- Real production chaos engineering with live failure injection → use Litmus, AWS FIS, or Chaos Monkey
- Unit testing pure functions or isolated modules → use `software-engineering-patterns` (references/testing-and-vitest-expert.md)
- Load/performance testing → use a dedicated load-testing tool (k6, Locust, Gatling)
- General incident management that is not a simulation → use `tam-operations` (references/incident-response/SKILL.md)

---

## Core concepts

### Firedrill vs chaos engineering vs load testing

| Dimension | Firedrill | Chaos Engineering | Load Testing |
|-----------|-----------|-------------------|--------------|
| Primary target | People + process | System resilience | Capacity limits |
| Blast radius | Zero (simulated) | Controlled real | Controlled real |
| Artifacts created | Simulated cases/tickets | Real failures | Real traffic |
| Success metric | Response quality score | System recovery time | Throughput/latency |
| Abort trigger | Real incident detected | Threshold breach | SLA violation |

### Simulation boundary

Firedrills operate within a strict boundary — no real customer-facing artifacts are created:
- Case numbers use `DRILL-` prefix (never real case IDs)
- Linked tickets use `DRILL-HELP-` prefix
- No real Salesforce, Jira, or PagerDuty artifacts are opened
- All state lives in the drill engine, not production systems

### Scenario lifecycle

```
PREFLIGHT --> RUNNING --> CLOSING --> CLOSED
    |             |
    v             v
 ABORTED      ABORTED
```

---

## Quick start

```javascript
// 1. List available scenarios
const scenarios = await mdb_case_firedrill_list_scenarios();

// 2. Start a drill
const drill = await mdb_case_firedrill_start({ scenario_id: 'scenario-1' });

// 3. Confirm preflight and go
await mdb_case_firedrill_confirm_preflight({
  drill_id: drill.drill_id,
  go_no_go: { ir_available: true, moderator_available: true,
               no_active_incidents: true, bridge_ready: true },
});

// 4. Post a comment as the IR
await mdb_case_firedrill_post_comment({
  drill_id: drill.drill_id,
  comment: { body: 'Investigating the reported replication lag.', author_role: 'TAM/IC' },
});

// 5. Check status and scorecard
const status = await mdb_case_firedrill_status();

// 6. Close and get report
const report = await mdb_case_firedrill_close({
  drill_id: drill.drill_id,
  retro_inputs: { what_worked: 'Fast triage', what_slowed_us: 'Tool gaps',
                  action_items: ['Create runbook'] },
});
```

---

## Phase 1: Preflight (safety gate)

All checks must pass before the drill transitions to running. ANY false check blocks the drill.

```javascript
const goNoGo = {
  ir_available: true,          // participant confirmed
  moderator_available: true,   // moderator confirmed
  no_active_incidents: true,   // no conflicting real S1/S2 incidents
  bridge_ready: true,          // communication channel ready
  abort_criteria_confirmed: true,
};
const go = Object.values(goNoGo).every(v => v === true);
```

**MCP tool:** `mdb_case_firedrill_confirm_preflight` — all `go_no_go` fields must be `true`.

**Pre-drill safety checklist:**
- [ ] No active S1/S2 incidents in the account
- [ ] All participants confirmed availability
- [ ] Bridge/communication channel operational
- [ ] Abort criteria defined and understood
- [ ] Scenario reviewed by moderator
- [ ] Real-incident auto-abort detection is active

---

## Phase 2: Running (execution loop)

```javascript
async function runDrillLoop(drillId, scenario) {
  for (const beat of scenario.beats) {
    await presentStimulus(drillId, beat);

    const response = await waitForResponse(drillId, {
      timeout: beat.maxResponseTime || 300_000,
    });

    const score = evaluateResponse(response, beat.rubric);
    const reveals = checkRevealConditions(response, beat.reveals);
    if (reveals.length > 0) await deliverReveals(drillId, reveals);

    if (shouldAbort(drillId)) {
      await mdb_case_firedrill_abort({ drill_id: drillId, reason: 'Safety condition triggered' });
      return { status: 'aborted' };
    }
  }
  return { status: 'complete' };
}
```

**MCP tools during running:**
- `mdb_case_firedrill_post_comment` — inject IR/customer/moderator messages
- `mdb_case_firedrill_create_linked_ticket` — simulate escalation artifacts
- `mdb_case_firedrill_status` — poll current state and scorecard

### Timeout configuration

```javascript
const TIMEOUT_CONFIG = {
  beatResponseTimeout:  5 * 60 * 1000,  // 5 min per beat
  maxDrillDuration:    45 * 60 * 1000,  // 45 min total
  preflightExpiry:     10 * 60 * 1000,  // 10 min for go/no-go
  abortGracePeriod:     2 * 60 * 1000,  // 2 min cleanup after abort
};
```

### Polling with exponential backoff

```javascript
async function pollUntilResponse(drillId, timeout = 300_000) {
  const startTime = Date.now();
  while (Date.now() - startTime < timeout) {
    const status = await mdb_case_firedrill_status();
    if (status.lastActivity > startTime) return status;
    const elapsed = Date.now() - startTime;
    const delay = Math.min(5000 * Math.pow(2, Math.floor(elapsed / 30000)), 30000);
    await sleep(delay);
  }
  return { timeout: true };
}
```

---

## Safety mechanisms

### Abort hierarchy

```
Level 0: Auto-abort (system-detected real incident)
Level 1: Moderator abort (manual kill switch)
Level 2: IR abort (participant requests stop)
Level 3: Timeout abort (max duration exceeded)
```

### Blast radius containment

1. **Namespace isolation:** all drill artifacts use `DRILL-` prefixed identifiers
2. **No production writes:** drill engine never creates real cases, tickets, or alerts
3. **State isolation:** drill state in ephemeral storage, not production databases
4. **Role sandboxing:** customer persona only reveals pre-scripted facts

---

## Scoring and grading

### Rubric structure

```javascript
const SCORING_RUBRIC = {
  dimensions: [
    {
      id: 'time_to_first_response',
      weight: 0.15,
      thresholds: {
        pass: { maxSeconds: 120 },
        partial: { maxSeconds: 300 },
        fail: { maxSeconds: Infinity },
      },
    },
    {
      id: 'correct_diagnosis',
      weight: 0.30,
      evaluation: 'llm_judge',
      criteria: 'IR correctly identifies the root cause category',
      partialCredit: 'IR identifies symptoms but not root cause',
    },
    {
      id: 'appropriate_escalation',
      weight: 0.20,
      evaluation: 'deterministic',
      checkFn: (timeline) => timeline.some(e => e.type === 'linked_ticket_created'),
    },
    {
      id: 'customer_communication',
      weight: 0.20,
      evaluation: 'llm_judge',
      criteria: 'IR provides clear, empathetic, technically accurate updates',
    },
    {
      id: 'documentation_quality',
      weight: 0.15,
      evaluation: 'llm_judge',
      criteria: 'IR documents findings and next steps clearly',
    },
  ],
  grading: { pass: 0.80, partial: 0.60, fail: 0.0 },
};
```

### Hybrid evaluation: deterministic + LLM-as-judge

- **Deterministic:** binary check against concrete criteria (e.g., linked ticket created)
- **Threshold-based:** time measurements against pass/partial/fail thresholds
- **LLM-as-judge:** quality evaluation for communication, diagnosis, documentation
- **Fuzzy matches weight 0.5x** to prevent noise from dominating exact matches

### Grade tiers

| Grade | Score range | Meaning |
|-------|-------------|---------|
| PASS | 0.80–1.00 | Met or exceeded expectations |
| PARTIAL | 0.60–0.79 | Showed understanding; missed key elements |
| FAIL | 0.00–0.59 | Did not meet minimum expectations |

---

## Reveal-condition state machine

Customer persona reveals information progressively based on IR actions:

```javascript
const REVEAL_CONDITIONS = [
  {
    id: 'reveal_error_logs',
    trigger: { type: 'keyword_in_comment', keywords: ['logs', 'error', 'stacktrace'] },
    content: 'Here are the error logs from the last 30 minutes: [simulated log data]',
    unlocks: ['reveal_connection_string'],
  },
  {
    id: 'reveal_cluster_info',
    trigger: { type: 'action_taken', action: 'linked_ticket_created' },
    content: 'The cluster name is Production-East-1, running MongoDB 7.0',
  },
];
```

Prerequisites chain: `reveal_connection_string` requires `reveal_error_logs` to have been unlocked first.

---

## Human checkpoint pattern

Critical decisions require human approval before the agent continues:

```javascript
async function humanCheckpoint({ type, message, timeout = 300_000 }) {
  console.log(`[CHECKPOINT:${type}] ${message}`);
  const decision = await waitForHumanInput(timeout);
  if (!decision) return { approved: false, reason: 'timeout' };
  return { approved: decision.go, checks: decision.checks };
}
```

**Checkpoint types:**
1. **Go/No-Go** — before preflight → running transition
2. **Abort Confirmation** — before aborting (unless auto-abort triggered)
3. **Score Override** — when LLM judge confidence is low
4. **Retro Input** — collecting what-worked/what-slowed-us

---

## Post-drill reporting

The `mdb_case_firedrill_close` tool returns a frozen scorecard and retrospective markdown. Track per-person improvement across drills using `irScoreHistory` and identify team-level patterns via `commonFailures` dimensions.

Full scorecard schema, trend tracking structure, report markdown template, agent orchestration loop, complete scenario definition schema, and customer persona prompt template: [references/scenario-patterns.md](./references/scenario-patterns.md)

---

## MCP tool reference

| Shorthand | Full MCP tool name | Phase |
|-----------|-------------------|-------|
| `firedrill_list_scenarios` | `mdb_case_firedrill_list_scenarios` | Pre-start |
| `firedrill_start` | `mdb_case_firedrill_start` | Preflight |
| `firedrill_confirm_preflight` | `mdb_case_firedrill_confirm_preflight` | Preflight |
| `firedrill_status` | `mdb_case_firedrill_status` | Any (read-only) |
| `firedrill_post_comment` | `mdb_case_firedrill_post_comment` | Running |
| `firedrill_create_linked_ticket` | `mdb_case_firedrill_create_linked_ticket` | Running |
| `firedrill_abort` | `mdb_case_firedrill_abort` | Any |
| `firedrill_close` | `mdb_case_firedrill_close` | Closing |

**Key parameter constraints:**
- `scenario_id`: must match `^scenario-\d+$`
- `drill_id`: must match `^drill-`
- `author_role` in `post_comment`: enum `'IR' | 'TAM/IC' | 'Moderator'`
- `go_no_go` in `confirm_preflight`: all values must be `true` to proceed

---

## Error handling

```javascript
async function safeDrillOperation(operation, drillId) {
  try {
    return await operation();
  } catch (error) {
    if (error.message?.includes('no active drill')) {
      return { status: 'no_drill', error };
    }
    if (error.message?.includes('invalid phase')) {
      const status = await mdb_case_firedrill_status();
      return { status: 'phase_mismatch', phase: status.phase, error };
    }
    if (error.message?.includes('timeout') || error.code === 'ETIMEDOUT') {
      return await operation(); // retry once
    }
    // Unknown error: abort as safety measure
    if (drillId) {
      await mdb_case_firedrill_abort({ drill_id: drillId, reason: `Unexpected error: ${error.message}` });
    }
    throw error;
  }
}
```

Always check current phase before state-changing operations:

```javascript
async function ensurePhase(expectedPhase) {
  const status = await mdb_case_firedrill_status();
  if (status.phase !== expectedPhase) {
    throw new Error(`Expected '${expectedPhase}' but found '${status.phase}'`);
  }
  return status;
}
```

---

## Anti-patterns

| Anti-pattern | Problem | Fix |
|-------------|---------|-----|
| Skipping preflight | Missed active S1 causes confusion between real and simulated | Always gate on preflight |
| Drill artifacts in production | Real Jira/Salesforce cases pollute production | All artifacts use `DRILL-` namespace |
| Scoring without rubric | Subjective pass/fail after the drill | Define rubric criteria before the drill starts |
| No abort path | Drill continues during a real incident | Auto-abort on real S1/S2; moderator always has manual kill switch |
| Over-scripting the IR | IR knows the scenario before it starts | IR discovers the problem organically |
| Binary pass/fail only | Discourages learning | Use three-tier grading (pass/partial/fail) |
| No retrospective | Closes without collecting lessons | Always collect what-worked, what-slowed-us, action items |
| Testing only happy paths | Validation theater | Include scenarios targeting known weak areas |

---

## Testing the firedrill engine itself

```javascript
describe('Firedrill Engine Integration', () => {
  it('transitions from preflight to running on confirm', async () => {
    const { drill_id } = await mcp.firedrill_start({ scenario_id: 'scenario-1', ir: {...} });
    await mcp.firedrill_confirm_preflight({
      drill_id,
      go_no_go: { ir_available: true, moderator_available: true,
                  no_active_incidents: true, bridge_ready: true },
    });
    const status = await mcp.firedrill_status();
    expect(status.phase).toBe('running');
  });

  it('aborts immediately when triggered', async () => {
    // ... setup running drill ...
    await mcp.firedrill_abort({ drill_id, reason: 'Real S1 incident detected' });
    const status = await mcp.firedrill_status();
    expect(status.phase).toBe('aborted');
  });

  it('scores partial credit for incomplete diagnosis', async () => {
    // ... setup running drill ...
    await mcp.firedrill_post_comment({
      drill_id,
      comment: { body: 'I see replication lag but not sure of root cause.', author_role: 'IR' },
    });
    const status = await mcp.firedrill_status();
    expect(status.scorecard?.scores?.correct_diagnosis?.grade).toBe('partial');
  });
});
```

---

## References

- [Gremlin: Introduction to GameDays](https://www.gremlin.com/community/tutorials/introduction-to-gamedays)
- [Shopify: Four Steps to Creating Effective Game Day Tests](https://shopify.engineering/four-steps-creating-effective-game-day-tests)
- [AWS Well-Architected: Run Simulations (SEC10-BP07)](https://docs.aws.amazon.com/wellarchitected/latest/security-pillar/sec_incident_response_run_game_days.html)
- [Scale Labs: Agentic Rubrics for Code Verification](https://labs.scale.com/blog/agentic-rubrics)
- [Google Cloud: A Methodical Approach to Agent Evaluation](https://cloud.google.com/blog/topics/developers-practitioners/a-methodical-approach-to-agent-evaluation)
