# Firedrill — Scenario Patterns and Agent Orchestration

## Full agent orchestration loop

```javascript
/**
 * Agent-driven firedrill orchestration.
 * The agent acts as both moderator and customer persona,
 * with human checkpoints at critical decision points.
 */
async function firedrillAgentLoop(config) {
  const { scenarioId, ir, moderator } = config;

  // Phase 1: List and select scenario
  const scenarios = await mcp.firedrill_list_scenarios();

  // Phase 2: Start drill (enters preflight)
  const drill = await mcp.firedrill_start({
    scenario_id: scenarioId,
    ir: { name: ir.name, email: ir.email },
    moderator: { name: moderator.name, email: moderator.email },
    bridge_link: config.bridgeLink,
    target_cluster: config.targetCluster,
  });

  // Phase 3: Human checkpoint — confirm preflight
  const preflightResult = await humanCheckpoint({
    type: 'go_no_go',
    message: `Drill ${drill.drill_id} ready. Confirm preflight checks:`,
    checks: scenarios.find(s => s.id === scenarioId).preflightChecks,
  });

  if (!preflightResult.approved) {
    await mcp.firedrill_abort({ drill_id: drill.drill_id, reason: 'Preflight not confirmed' });
    return { status: 'cancelled', reason: 'preflight_rejected' };
  }

  await mcp.firedrill_confirm_preflight({
    drill_id: drill.drill_id,
    go_no_go: preflightResult.checks,
  });

  // Phase 4: Execute scenario beats
  const scenario = scenarios.find(s => s.id === scenarioId);
  for (const beat of scenario.beats) {
    await mcp.firedrill_post_comment({
      drill_id: drill.drill_id,
      comment: {
        body: beat.customerMessage,
        author_role: 'IR',
        author_name: 'Customer Persona',
        is_internal: false,
      },
    });

    const status = await pollUntilResponse(drill.drill_id, beat.timeout);

    // Check reveal conditions
    if (meetsRevealCondition(status.lastComment, beat.revealTrigger)) {
      await mcp.firedrill_post_comment({
        drill_id: drill.drill_id,
        comment: {
          body: beat.revealContent,
          author_role: 'IR',
          author_name: 'Customer Persona',
          is_internal: false,
        },
      });
    }
  }

  // Phase 5: Close and report
  const report = await mcp.firedrill_close({
    drill_id: drill.drill_id,
    retro_inputs: await collectRetroInputs(drill.drill_id),
  });

  return { status: 'complete', report };
}
```

## Complete scenario definition schema

```javascript
const SCENARIO_SCHEMA = {
  id: 'scenario-1',
  title: 'Oplog Lag Causing Replication Delay',
  severity: 'S2',
  category: 'replication',

  initialContext: {
    customerMessage: 'Our writes are succeeding but reads from secondaries are stale. Started 20 minutes ago.',
    caseMetadata: {
      product: 'MongoDB Atlas',
      tier: 'M30',
      region: 'us-east-1',
    },
  },

  beats: [
    {
      id: 'beat-1',
      type: 'customer_message',
      content: 'We are seeing 5-second lag on secondary reads.',
      maxResponseTime: 120_000,
      rubric: {
        criteria: 'IR acknowledges the issue and asks clarifying questions',
        passExample: 'Asking about write volume, recent deployments, or cluster tier',
      },
      reveals: ['reveal_error_logs'],
    },
    {
      id: 'beat-2',
      type: 'escalation_needed',
      triggerCondition: 'After IR identifies oplog as potential cause',
      rubric: {
        criteria: 'IR creates a linked help ticket for cluster team review',
      },
      reveals: ['reveal_cluster_info'],
    },
  ],

  goldenPath: [
    'Acknowledge customer concern',
    'Ask about write patterns and recent changes',
    'Check oplog window and replication lag metrics',
    'Identify oplog overflow or heavy write burst',
    'Escalate to cluster team if needed',
    'Communicate timeline and next steps to customer',
  ],

  preflightChecks: ['ir_available', 'moderator_available', 'no_active_incidents', 'bridge_ready'],
};
```

## Customer persona prompt template

```text
You are playing a customer persona in a firedrill simulation.

SCENARIO: {scenario_title}
SEVERITY: {severity}
YOUR ROLE: {customer_role} at {customer_company}

PERSONALITY:
- Technical level: {technical_level} (junior/mid/senior)
- Frustration level: {frustration_level} (calm/concerned/frustrated/angry)
- Communication style: {style} (terse/verbose/technical/non-technical)

FACTS YOU KNOW (reveal ONLY when asked relevant questions):
{reveal_conditions_formatted}

RULES:
1. Never reveal facts unprompted — wait for the IR to ask
2. Stay in character throughout the drill
3. If the IR asks something not covered by your facts, say "I am not sure, let me check"
4. Escalate frustration if response time exceeds {escalation_threshold} minutes
5. Acknowledge good communication from the IR naturally

INITIAL MESSAGE:
{initial_customer_message}
```

## Report markdown template

```markdown
# Firedrill Report: {scenario_name}

**Drill ID**: {drill_id}
**Date**: {date}
**Duration**: {duration} minutes
**Grade**: {overall_grade} ({weighted_score_pct}%)

## Participants
- **IR**: {ir_name}
- **Moderator**: {moderator_name}

## Scorecard

| Dimension | Grade | Score | Detail |
|-----------|-------|-------|--------|
| Time to First Response | {grade} | {score} | {detail} |
| Correct Diagnosis | {grade} | {score} | {detail} |
| Appropriate Escalation | {grade} | {score} | {detail} |
| Customer Communication | {grade} | {score} | {detail} |
| Documentation Quality | {grade} | {score} | {detail} |

## Timeline
{timeline_entries}

## Retrospective
### What Worked
{what_worked}

### What Slowed Us
{what_slowed_us}

## Action Items
{action_items_checklist}

## Recommendations
{ai_generated_recommendations}
```
