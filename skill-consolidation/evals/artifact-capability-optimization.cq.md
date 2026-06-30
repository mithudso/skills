# Competency questions — AI-Assisted Artifact & Capability Optimization

Family seeded by concept-family-explorer (subject: claude skill optimization, llm skills,
prompt optimization, document optimization, code optimization, claude plugins, agent
optimization). Each bullet is a question someone working this family should be able to ask
and reach a skill. `→ GAP` marks a question that does NOT yet resolve to a HAVE/researched
node (these are the above-threshold gaps from the scored table).

- How do I audit and improve a Claude skill's trigger accuracy, token budget, and peer-deferral edges? → skill-optimizer (HAVE)
- Which prompt-optimization algorithm (APE / OPRO / DSPy-MIPROv2 / GEPA / TextGrad) fits my training data, and how do I run a multi-pass prompt audit to convergence? → prompt-deep-optimizer (HAVE)
- How do I iteratively critique and fix a document until no Medium-or-higher findings remain? → document-critique / ddo (HAVE)
- How do I review and fix a whole source file or repo to convergence with build/lint/test gating? → code-deep-optimizer (HAVE)
- How do I structure, package, and distribute a Claude Code plugin (hooks, commands, skills, agents)? → claude-code-plugins (HAVE)
- How do I choose a multi-agent orchestration topology and build the agent harness (action space, observation format, guardrails)? → ai-agents-orchestration / agent-harness-construction (HAVE)
- How do I evaluate whether an LLM-app change actually improved output quality — LLM-as-judge calibration, golden datasets, CI regression suites? → ai-agent-engineering ▸ Eval-Driven Development (HAVE)
- How can I automatically design or search over an agent's architecture / multi-step prompt-program and optimize the whole compound system from a metric (ADAS, AFlow, GEPA-for-systems, Darwin-Gödel Machine, Trace)? → GAP
- How do I compress a long prompt or context to cut tokens without losing task quality (LLMLingua, gist tokens, soft-prompt distillation)? → GAP
- How do I detect and repair "prompt rot" and re-optimize my prompts and skills when the underlying model version changes? → GAP
