# psychology

**Category:** Science, Biology & Medicine
**Platform:** Claude
**Original Path:** claude/standalone/psychology

## Description
Psychology domain hub — clinical, developmental, social & positive psychology as a deep knowledge domain (vs applied-psychology = TAM/operator lens). SPOKES: psychology-positive; psychology-clinical-personality; psychology-neurodevelopmental; psychology-stress-trauma; psychology-social; psychology-confidence-identity; psychology-influence-depth. TRIGGER: personality disorders (narcissism, BPD, psychopathy, dark triad); neurodevelopmental (ADHD, autism, masking); positive psychology / PERMA / wellbeing; stress science / allostatic load / PTSD / polyvagal theory; social hierarchy / status / dominance; self-efficacy / impostor syndrome; influence theory (Cialdini, compliance, cult dynamics). SKIP: operator psychology (behavior change, nudges, decision-making, human-AI) → applied-psychology; persuasion in practice → content-and-marketing-writing; moral psychology → applied-psychology.

---

# Psychology (router hub)

Deep psychology domain knowledge — clinical, developmental, social, and positive. For operator/TAM behavioral lens → `applied-psychology`. Routes to 7 spokes on demand.

## Routing table

| Topic | Load when | Skill |
|---|---|---|
| Positive psychology / wellbeing | "What does research say about happiness / flourishing / meaning?"; PERMA, flow, eudaimonia, hedonic adaptation | `psychology-positive` |
| Personality disorders | narcissism, psychopathy, BPD, dark triad, cluster B, antisocial PD, defense mechanisms, attachment styles | `psychology-clinical-personality` |
| Neurodevelopmental | ADHD (exec function, dopamine, RSD), autism spectrum (masking, sensory, social cognition, stimming) | `psychology-neurodevelopmental` |
| Stress & trauma | allostatic load, HPA axis, cortisol, chronic stress, PTSD, burnout, polyvagal theory | `psychology-stress-trauma` |
| Social psychology | status, dominance vs prestige, social comparison, conformity, obedience, in-group/out-group, social identity | `psychology-social` |
| Confidence & identity | self-efficacy (mechanistic), self-esteem, impostor syndrome, identity formation | `psychology-confidence-identity` |
| Influence & persuasion theory | Cialdini deep-dive, compliance ladders, social proof mechanisms, cult dynamics, inoculation, persuasion resistance | `psychology-influence-depth` |

## Boundary: psychology vs applied-psychology

| Question type | Route |
|---|---|
| "How does narcissism work / what causes it / what does research show?" | **psychology** (clinical depth) |
| "How do I deal with a narcissistic stakeholder?" | **applied-psychology** (operator lens) |
| "What is Cialdini's reciprocity principle and its mechanisms?" | **psychology-influence-depth** (theory) |
| "How do I use reciprocity to drive product adoption?" | **applied-psychology** → applied-psychology (references/persuasion-and-influence-psychology.md) spoke |
| "What is polyvagal theory?" | **psychology-stress-trauma** (mechanistic) |
| "Why won't my customer trust me?" | **applied-psychology** → applied-psychology (references/trust-and-psychological-safety.md) spoke |

Rule: "how does it work / what does research say?" → psychology hub. "how do I use it?" → applied-psychology.