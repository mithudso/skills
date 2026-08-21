# Context Files for LLMs: What They Are, Why They Work, and How to Build One

*Last updated: 2026-07-21*

## 1. What Is a Context File?

A context file is a persistent, version-controlled document that tells an LLM
what it needs to know *before* a conversation starts: who it's talking to,
what rules apply, what conventions to follow, what's been decided already. It
is loaded once per session (or once per turn, depending on the harness) rather
than re-explained by a human every time.

Familiar examples:

- `CLAUDE.md` — Claude Code's project-level context file (coding conventions,
  build commands, architecture notes, "don't do X" rules)
- `AGENTS.md` — the emerging cross-tool standard for the same purpose
- `.cursorrules` — Cursor's equivalent
- A system prompt block containing company voice guidelines, product facts,
  or a support agent's escalation policy
- A "customer context" doc a support or sales team maintains per account

The common thread: **static, reusable, high-signal facts that would otherwise
be retyped into every prompt.** A context file turns tribal knowledge and
repeated instructions into a single artifact the model reads once and applies
consistently.

## 2. Why Context Files Are Useful

**Consistency.** Without a context file, every session starts from a blank
slate; the model infers conventions from whatever code or messages it can
see, and infers wrong at least some of the time. A context file removes the
guesswork: the same rules apply every session, for every user, without
re-explaining.

**Time savings.** Any fact a human would otherwise type into the first
message of *every* session belongs in a context file instead. Onboarding a new
teammate or a new AI session becomes "read this file" instead of a 20-minute
verbal briefing repeated indefinitely.

**Correctness under scale.** As a project or account grows, the number of "oh
by the way" facts grows with it. A context file is the only mechanism that
scales: a human can't remember to mention 40 caveats every time, but a file
loaded automatically never forgets one.

**Cheaper inference.** This is the concrete financial case, covered in
Section 8 below: a stable context file is cacheable. A context re-explained
in prose every turn is not.

**Auditable behavior.** When the model does something wrong, the context file
is the first place to look, and the first place to fix it. Compare this to
prompting the model differently every session and hoping the fix "sticks,"
which it won't.

## 3. How to Use a Context File With an LLM

Mechanically, a context file is just text injected into the prompt before the
user's actual request. Three common injection points:

1. **System prompt** — the file's content becomes (or is appended to) the
   system prompt. This is the standard mechanism for CLAUDE.md/AGENTS.md-style
   files: the harness reads the file from disk and prepends it automatically.
2. **First user turn** — some tools inject the context file as if it were the
   first message in the conversation, before the real user question.
3. **Tool-retrieved reference** — for very large context files, don't inject
   the whole thing; give the model a tool to `Read` specific sections on
   demand (this is how large skill/reference systems like this one work: the
   top-level file is small and routes to deeper reference files only when
   relevant).

**Ordering matters.** Put the most stable content first (role, rules,
conventions) and the most volatile content last (today's date, the current
user message). This ordering is what makes prompt caching work (see Section
8), and it also matches how models attend to context: content at the start
and end of the window is used more reliably than content buried in the middle
("lost-in-the-middle").

**Don't just dump everything.** A context file is not a junk drawer. Every
line should be something the model would otherwise get wrong or have to ask
about. If a fact never affects model behavior, it doesn't belong in context.
It belongs in documentation a human reads, not a context file a model reads.

## 4. How to Create a Context File for a Specific Topic or Customer

Building one well is a distillation exercise, not a transcription exercise.
Steps:

**Step 1 — Collect the raw material.** Gather every source of truth: prior
conversations, tickets, docs, meeting notes, Slack threads, the actual
codebase or account records. Don't start from memory; start from source.

**Step 2 — Extract only what changes model behavior.** For each candidate
fact, ask: "If the model didn't know this, would it act differently or say
something wrong?" If yes, keep it. If it's background color with no behavior
implication, cut it. This is the single highest-impact step: in practice,
most first drafts run several times longer than they need to because this
step gets skipped.

**Step 3 — Structure by stability, not by topic.** Group content into layers:
things that almost never change (role, hard rules) at the top; things that
change per-session (current status, active priorities) lower down; anything
truly ephemeral (a secret valid only for this one request) never persisted to
the file at all.

**Step 4 — Write facts, not prose.** A context file is closer to structured
notes than an essay. Bullet points, tables, and short declarative sentences
beat paragraphs. The model doesn't need to be persuaded of anything; it needs
unambiguous facts.

**Step 5 — Name it precisely.** For a topic-specific file: name the domain and
scope in the first line ("MongoDB Atlas performance conventions for Project
X"). For a customer file: lead with identity and current state (who they are,
what tier/plan, what's open right now), then history, in reverse-recency
order so the freshest facts are easiest to find.

**Step 6 — Version it and date it.** Context files go stale. Put a last-
updated date at the top and treat staleness detection as a first-class
concern (see Section 10): a wrong fact silently followed is worse than no
fact at all.

**Step 7 — Test by deletion.** After a first draft, delete a line and check
whether the model's behavior changes without it. If it doesn't, the line
wasn't earning its place. This is the cheapest quality check available, and
it is almost never actually run. It should be.

## 5. Worked Examples

### Example A — Project context file (CLAUDE.md-style)

```markdown
# Project: billing-service

## Role
You are assisting with the billing-service Go microservice.

## Hard rules
- Never commit directly to `main` — open a PR.
- All monetary values are stored as integer cents, never floats.
- Migrations require a rollback script in the same PR.

## Conventions
- Error wrapping: use `fmt.Errorf("doing X: %w", err)`, always with context.
- Tests: table-driven, one `t.Run` per case, no shared mutable fixtures.

## Current state (updated 2026-07-15)
- Migrating from Stripe API v1 to v2; v1 code is deprecated but not yet removed.
- Do not add new callers of `legacy/stripe_v1.go`.
```

Every line here fails or passes the Step 2 test: each one changes what the
model would otherwise do wrong (default to floats for money, commit to main,
call the deprecated Stripe client).

### Example B — Customer/account context file (TAM- or support-style)

```markdown
# Account: Acme Corp (Tier: Enterprise, Renewal: 2026-11-01)

## Current state (updated 2026-07-20)
- 2 open P2 cases: connection pool exhaustion (case #4821), slow
  aggregation on `orders` collection (case #4830).
- Primary technical contact: Priya Shah (DBA lead). CC all replies to her.
- Contract does NOT include Atlas Search — do not recommend it without
  flagging the upsell.

## History that still matters
- 2026-03: migrated off self-managed MongoDB to Atlas. Some internal
  runbooks still reference the old topology — treat those as stale.
- Known sensitivity: a prior TAM over-promised a timeline on case #3990 that
  slipped 3 weeks. Be conservative on ETAs with this account.
```

Notice this isn't a transcript of every interaction: it's the subset that
would cause an assistant to say or do the wrong thing if missing (recommend
an unpurchased feature, ignore the DBA lead, over-promise a date).

### Example C — Topic-specific reference file (this repository's own pattern)

This very skill system uses context files as a routing layer: a small
top-level file states scope and triggers, and defers deep content to
`references/*.md` files loaded only when relevant. That's context engineering
applied recursively: the context file about how to use context files is
itself built to avoid dumping everything into one window at once.

## 6. Use Cases

- **Coding agents** (Claude Code, Cursor, Copilot): project conventions, build
  commands, architecture constraints, "don't touch this" warnings.
- **Customer support / account management**: per-account facts, open issues,
  sensitivities, contract boundaries, exactly the kind of thing that
  currently lives in a human's head and evaporates when they change roles.
- **Domain expert assistants**: a stable body of domain rules (compliance
  requirements, style guides, API conventions) that should apply identically
  across every session regardless of which model or which day it is.
- **Multi-agent systems**: a shared context file lets independent agents stay
  aligned on facts (schema, terminology, current sprint) without each one
  re-deriving them from scratch.
- **Personal AI assistants**: standing preferences (tone, formatting,
  recurring commitments) that shouldn't need restating each session.

## 7. Why It's Valuable (Beyond Convenience)

Section 2's "correctness under scale" point holds for a single team; this is
what changes when the scale gets larger still. A single user prompting a
model carefully can get good output without any of this. The value compounds
further when:

- Sessions are short-lived and stateless: a context file is the substitute
  for continuity a human collaborator would otherwise provide from memory.
- The cost of a wrong answer is high (compliance, money, customer trust): a
  context file converts "the model might remember this" into "the model
  always has this."

## 8. Why It's Cheaper (Prompt Caching)

This is the part that's easy to miss: a well-structured context file isn't
just organizationally useful, it directly reduces API cost via prompt
caching.

Anthropic's caching mechanism lets you mark a stable prefix of your prompt
(e.g., `cache_control: {"type": "ephemeral"}`) so that identical leading
content is *read from cache* on subsequent calls instead of reprocessed at
full price. A context file, by construction static and reused across many
calls, is exactly the kind of content that benefits:

```json
{
  "role": "user",
  "content": [
    {
      "type": "text",
      "text": "<the static context file content>",
      "cache_control": {"type": "ephemeral"}
    }
  ]
}
```

Rules of thumb for capturing this savings:
- Put the context file content first, before anything session-specific.
- Never put a cache breakpoint inside content that changes per turn, since it
  defeats the cache and adds write cost with no benefit.
- A per-session or per-account block (dynamic) should sit *after* the cached
  static file, with a shorter TTL if it's cached at all.

Practically, this means the same context file re-explained in prose inside
every prompt (no caching possible, because the text differs slightly each
time, or because there's no file structure to cache against) costs
meaningfully more than the same facts committed to a stable, unchanging file
loaded the same way every time. The savings compound: a support team running
thousands of sessions a day against the same account context file is doing
the token-cost equivalent of paying full price vs. bulk price for the same
information, every single call.

## 9. When NOT to Use a Context File

- **One-off, single-session tasks.** If there's no repetition, there's
  nothing to amortize, so just put the facts in the prompt directly.
- **Facts that change every call.** Highly volatile data (today's stock
  price, the current timestamp, a one-time token) doesn't belong in a cached
  static file: it breaks the cache and adds staleness risk for no benefit.
  Inject it dynamically instead.
- **Secrets or credentials valid for a single request.** Never persist these
  to a file that outlives the request. That's a security liability, not a
  convenience.
- **When it becomes a junk drawer.** If a context file balloons because
  everything got added and nothing was ever removed, it stops being
  high-signal and starts actively degrading output (see "context rot" below).
  At that point the fix is aggressive pruning, not more file.
- **Cross-account bleed.** Never let one account's or user's context file
  leak into another's session: a shared cache layer must not contain PII or
  account-specific facts that could surface for the wrong user.

## 10. Known Limitations

**Context rot.** Every frontier model's output quality degrades as context
grows, well before the context window physically overflows (per 2025 "context
rot" research, a 200K-token model can show meaningful degradation by around
50K tokens). A bigger context file is not free; it competes with everything
else in the window for the model's effective attention.

**Lost-in-the-middle.** Models attend to the start and end of context far
more reliably than the middle. A fact buried in the middle of a long context
file is less reliably followed than the same fact at the top or bottom.
Structure the file accordingly, and don't just append forever.

**Staleness.** A context file is only as good as its last update. Nothing
forces someone to update it when the underlying facts change, and a
confidently-stated stale fact is more dangerous than no fact: the model will
act on it without hedging. Put a last-updated date on every file and treat
"is this still true" as an explicit, recurring check, not an assumption.

**Context poisoning.** If a hallucination or wrong fact makes it into the
file, the model will keep referencing it turn after turn, compounding the
error instead of self-correcting. A context file needs the same review
discipline as any other piece of source-of-truth documentation, arguably
more, since a model won't push back on it the way a skeptical human might.

**It's not a substitute for retrieval.** A context file holds what's true
*most of the time, for most sessions*. Anything genuinely per-query
(searching a large corpus, looking up one specific record) belongs in
retrieval (RAG) or a tool call, not stuffed into the static file "just in
case." Conflating the two is the most common authoring mistake: it produces
a context file that's both bloated and still incomplete.

**No enforcement mechanism.** A context file is advisory, not a hard
constraint. The model can still ignore or contradict it, especially under
long conversations or adversarial input. For anything safety-critical,
back the context file's rules with actual guardrails (validation, tool
permissions, human review) rather than relying on the model reading and
obeying a markdown file.

## 11. Summary

A context file is the cheapest, most durable lever available for making LLM
behavior consistent, correct, and inexpensive at scale. Build it by
distilling, not transcribing, only the facts that change model behavior;
structure it stable-to-volatile; keep it current; and treat its size as a
cost, not a feature. Used well, it turns "explain this every time" into
"state this once, cheaply, forever." Used poorly, it becomes bloat that
actively degrades the same output it was meant to improve.
