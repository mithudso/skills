# Sub-Agent: Architecture & Design

Tag: `architecture`

Pass the full diff to the sub-agent with this prompt:

> Review the architecture, overall design approach, and correctness of this MongoDB C++ PR. Apply the `coding-best-practices@xdg-claude` guidelines (deep modules, information hiding, pull complexity downward, general-purpose over special-case, design-it-twice) as a lens on every design decision below. Focus on:
>
> **Scope & Approach (Step 1)**
>
> - Right approach? Is there a fundamentally better way to solve this?
> - Comprehensibility: will the code be clear to a future reader without ticket context?
> - Should this be split into smaller PRs or sub-tasks?
> - Missing pieces: migration path, rollback story, feature flag, documentation?
> - Does the implementation match the stated design / ticket?
>
> **Correctness & Safety (blocking)**
>
> - Race conditions / concurrency story documented or enforced
> - Negative, zero, and overflow edge cases handled
> - Error paths don't silently corrupt state
> - Lifetimes / ownership semantics clear (raw ptr vs intrusive_ptr vs unique_ptr)
> - State that should be reset _is_ reset at the right lifecycle point
> - Related call sites: does this pattern appear elsewhere with the same bug?
>
> **API & Interface Shape (blocking or follow-up)**
>
> - Preconditions documented. Ideally encoded via an assertion _and_ public comments where possible.
>   If the assertion is publicly visible then a comment is not needed.
> - Optional parameters have clear semantics (nullopt meaning documented)
> - Bool parameters should be enum or named struct field
> - Virtual dispatch: is the indirection actually needed?
> - Struct vs class with getters — prefer struct when there's no invariant to protect
> - Return type consistency across related methods
> - Guard/RAII API — is the lifetime obvious at the call site?
>
> Pattern: "Could we do something virtual-y to handle this?" / "Should this be a StatusWith?"
> Look for refactors that make the code impossible to misuse.
