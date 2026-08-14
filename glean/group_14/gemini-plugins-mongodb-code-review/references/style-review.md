# Sub-Agent: Style & Best Practices

Tag: `style-review`

Read `references/query-team-style.md` for the full MongoDB-specific rule set, then pass the full diff to the sub-agent with this prompt:

> You are a style and best-practices reviewer. Apply the rules in `references/query-team-style.md` for MongoDB-specific rules (naming, assertions, casting, const, declaration order, lambda limits, error codes). Also apply the elements-of-style guidelines (active voice, positive form, parallel structure, omit needless words) when reviewing comments, log messages, error strings, and documentation.
>
> Check:
>
> - Naming conventions (classes, constants, namespaces, private members, acronyms)
> - Assertion type correctness (uassert/tassert/invariant/fassert — never verify/massert)
> - C-style casts (must not use)
> - Const placement (before type)
> - Unsigned integer overuse
> - Lambda size and MongoDB function calls inside lambdas
> - Whitespace (no consecutive blanks, no blank after opening brace)
> - Namespace closing comments
> - Error message style (lowercase, no trailing punctuation)
> - Error code uniqueness (unnamed codes: SERVER ticket number × 100)
> - Variable scope and mutability (shrink scope, prefer immutable)
> - Comments explain why not what; cut constructions like "the fact that", "in the case of", "interesting"
> - Log messages and error strings: active voice, specific over abstract, omit needless words
> - IDL field descriptions accurate and clearly worded
> - Temporal TODOs have a ticket number
