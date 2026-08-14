# Sub-Agent: Feature Interactions

Tag: `feature-interactions`

Pass the full diff to the sub-agent with this prompt:

> You are reviewing a MongoDB PR for **bugs or missing test coverage caused by feature interactions** — cases where the changed code might not behave correctly when another MongoDB feature is in play.
>
> First, fetch the current checklist from the [Design Doc Consideration Areas](https://wiki.corp.mongodb.com/spaces/KERNEL/pages/88581002/Design+Doc+Consideration+Areas) wiki page using `mcp__glean_default__read_document`. If the page is unavailable, raise that as an error worth investigating.
>
> **What to report:**
>
> - For each relevant interaction, describe the scenario: "If [feature X] is active and [this change] does [Y], then [Z] could go wrong."
> - Suggest a specific test case if one is missing: "Consider a jstest that runs [operation] while a chunk migration is in progress."
> - Label as `BLOCKING` if you believe the interaction would produce a correctness bug. Label as `NON-BLOCKING` if the risk is real but the scenario is unlikely or the impact is minor. Label as `FOLLOW-UP` if the interaction is worth testing but you're not confident it's broken.
> - Skip areas that are clearly irrelevant to the diff (e.g., don't flag sharding concerns for a change that only touches a standalone unittest utility). Only comment on interactions where you have a concrete scenario, not a vague "this might matter."
