# Debiasing Procedures & Worked Operator Scenarios

You cannot will a bias away — awareness alone barely moves the needle. What works is
*structured procedures* that force a deliberate (System-2) pass at the moment of decision.
Below: the evidence-backed debiasers, then worked scenarios.

---

## The four debiasers that earn their keep

### Consider-the-opposite
Before committing, explicitly generate reasons the favored conclusion could be **wrong**, and
arguments for the alternative. The best-evidenced general-purpose debiaser; particularly
effective against **anchoring**, **overconfidence**, and **confirmation bias**.
- *Mechanics:* "What would have to be true for the opposite call to be right? What evidence am
  I discounting?" Assign someone the explicit role of arguing the other side.
- *Caveat:* asymmetric — strong against high anchors and overconfidence; weaker against some
  other biases. Use it as the default first move, not a universal cure.

### Premortem (Gary Klein)
*Before* acting, assume the decision has already failed catastrophically, then have everyone
write down why. "Prospective hindsight" makes people ~30% better at identifying failure causes
and legitimizes dissent that groupthink would otherwise suppress.
- *When:* any consequential, hard-to-reverse decision — migration plans, architecture choices,
  big renewal/expansion bets, go-live decisions.
- *Output:* a ranked list of failure modes → mitigations folded into the plan.

### Reference-class forecasting (Flyvbjerg; Kahneman's "outside view")
Forecast from the **distribution of comparable past cases**, not from this case's inside
details. The cure for the **planning fallacy** and optimism bias.
- *Steps:* (1) define a reference class of similar past projects; (2) get the distribution of
  their actual outcomes (cost, duration, success rate); (3) position this case in that
  distribution; (4) adjust only with strong, specific justification.
- *Effect:* replaces hopeful assumptions with base rates — the single biggest lever on chronic
  estimate optimism.

### Checklists
Force a comprehensive, disciplined pass so System 1 can't skip a step. Effective for recurring,
high-stakes, error-prone decisions (Gawande's surgical/aviation evidence).
- *Caveat:* only work with **consistent** adherence — partial use creates false confidence.
  Keep them short and killable; a checklist nobody runs is worse than none.

> Why not just "be aware"? Debiasing-by-education has a weak track record. The reliable wins
> come from changing the *procedure* or the *environment* (and from boosts — teaching a better
> decision rule — per `references/choice-architecture.md`), not from telling people to try harder.

---

## Worked scenarios

### Scenario A — Negotiating price / scope with a customer
- **Effects in play:** anchoring (their opening number), reference dependence (their status quo),
  framing (gain vs loss).
- **Moves:** open first with a defensible anchor when you can; if they anchor low, name it and
  re-anchor against your own reference (total cost of the current state, value delivered) before
  you start trading. Frame the renewal as *protecting* current capability (loss frame relative to
  their reference point) rather than as *net-new* spend — but only if it's honestly true.
- **Your own guard:** run consider-the-opposite on your walk-away number so their anchor doesn't drag it.

### Scenario B — A migration/upgrade keeps slipping
- **Effects in play:** present bias / hyperbolic discounting (upfront cost, delayed benefit →
  perpetually deferred), planning fallacy (each new date is optimistic), status-quo bias.
- **Moves:** reference-class forecast the timeline from comparable migrations (not this team's
  hope); premortem the go-live; reduce the upfront cost the customer feels (stage it, do the
  hard first step *with* them = EAST "Easy"/"Timely"); set a default date with a commitment device.

### Scenario C — Customer won't move off a failing design "because of how much they've invested"
- **Effects in play:** sunk-cost fallacy / escalation of commitment.
- **Moves:** explicitly separate past (unrecoverable) spend from the forward decision — "from
  today, ignoring what's already spent, which path has the better expected outcome?" Provide a
  face-saving reframe (the prior investment *taught* something) so abandoning isn't an admission
  of failure.

### Scenario D — Your own capacity forecast / account-risk call
- **Effects in play:** overconfidence, confirmation bias, availability (the last loud incident).
- **Moves:** widen intervals via reference class; consider-the-opposite on your "this account is
  healthy" prior; check whether one vivid recent event (availability) is distorting your
  risk ranking; use base rates of churn for similar accounts, not vibes.

### Scenario E — Post-incident review
- **Effects in play:** hindsight bias ("obviously going to fail"), confirmation bias in
  root-cause.
- **Moves:** reconstruct the timeline with only information **knowable at the time**; resist the
  single-cause narrative; see `postmortem-writing` for blameless-framing and hindsight controls.

---

## One-line field guide

> Name the mechanism → pick the matching procedure (consider-the-opposite / premortem /
> reference class / checklist) or environment change (default / framing / EAST) → for anything
> long-term, prefer a **boost** (teach the customer the better rule) over a nudge you control.
> Never tell the customer they're being irrational.
