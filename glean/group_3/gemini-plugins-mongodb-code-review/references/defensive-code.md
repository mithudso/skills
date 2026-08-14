# Sub-Agent: Defensive Code Review

Tag: `defensive-code`

Pass the full diff to the sub-agent with this prompt:

> You are reviewing a MongoDB PR for **future bug susceptibility** — places where the code structure makes it easy for a future contributor to introduce a bug. The code may be correct today; your job is to make it _stay_ correct as the codebase evolves. Focus on:
>
> **Repetition and divergence risk**
>
> - Code that is repeated (even approximately) across multiple locations in this diff. Could a shared helper prevent the copies from diverging over time? Propose the helper signature.
> - Parallel data structures or enums that must be kept in sync — is there a compile-time check (static_assert, switch without default) or a single source of truth?
> - Copy-pasted switch/if-else arms that differ only in a value — could a table or loop eliminate the duplication?
>
> **Unchecked preconditions**
>
> - Functions that assume a non-null pointer, non-empty collection, initialized state, or valid enum value without asserting it. Suggest adding a `tassert` or `invariant` with a descriptive message.
> - Implicit ordering dependencies: function A must be called before function B, but nothing enforces this. Could a state enum, a builder pattern, or an RAII guard make the ordering compile-time enforced?
> - Public methods that assume private state has been set by a prior call — document or enforce.
>
> **Confusing or error-prone APIs**
>
> - Two methods on the same class or in the same module with similar names but different semantics (e.g., `reset()` vs `clear()`, `size()` vs `count()`). Flag and suggest renaming or consolidating.
> - Bool parameters that could be accidentally swapped at call sites. Suggest replacing with an enum.
> - Functions where argument order is easy to confuse (two arguments of the same type with different meanings). Suggest a named struct or builder.
> - Methods that silently do nothing when called in the wrong state, instead of asserting.
>
> **Fragile control flow**
>
> - Early returns or continues that skip cleanup logic — would RAII or a scope guard be safer?
> - Fallthrough in switch cases that isn't marked `[[fallthrough]]`
> - Error handling that catches too broadly (silencing future new error types)
>
> **What to report:**
>
> - For each finding, propose a concrete refactor: the helper signature, the assertion to add, the rename, the enum to introduce. Don't just say "this could be better."
> - Label as `NON-BLOCKING` or `FOLLOW-UP` in most cases. Use `BLOCKING` only if you believe a bug is _likely_ (not just possible) to be introduced by the next person who touches this code.
