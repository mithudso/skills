# Sub-Agent: Observability & Future-Proofing

Tag: `observability`

Pass the full diff to the sub-agent with this prompt:

> Review this MongoDB C++ PR for observability gaps and future-proofing concerns. Focus on:
>
> **Observability & Diagnosability**
>
> - Logging: key entry/exit points and error branches logged at the right level?
>   Logs should include actionable context (paths, names, error codes) not just "something failed."
>   Refer to docs/logging.md in the 10gen/mongo repo if the PR is targeting this codebase.
> - Metrics / counters: should this be tracked in serverStatus, $currentOp, slow query logs, or $queryStats?
> - Server parameters / knobs: is there a case for a serverParameter to tune behavior without a code change?
>   Flag if a change is risky enough that a customer might need a knob to work around a problem.
> - FTDC impact: does a new metric bloat FTDC output? Any redundant fields?
> - Metric naming consistency: same word order and singular/plural convention as adjacent metrics
> - Feature flag: is the code change protected by a feature flag? Should it be?
>
> **Future-Proofing & Modularity (non-blocking or optional)**
>
> - Logic duplicated across 3+ sites worth extracting
> - Giant functions with repeated blocks — could a helper make the algorithm read like English?
> - Could this be a static method / free function rather than a method reaching into a subclass?
> - Is there a shared helper file this utility belongs in?
> - Temporal comments ("// TODO: remove after X") have a ticket number
>
> Pattern: "could keep some of these in here with LOGV2_DEBUG at higher verbosity" /
> "might be useful to include the path we were looking for in this error"
