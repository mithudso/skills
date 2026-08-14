# Heuristics & Biases — Catalog with Operator Application

The heuristics-and-biases program (Tversky & Kahneman, 1974, *Science* 185:1124) holds
that judgment under uncertainty relies on a few mental shortcuts that are economical and
usually effective but produce **systematic, predictable** errors. Below: each effect, the
canonical demonstration, and one TAM/customer application note.

> Caution before applying: Gigerenzer's ecological-rationality program shows many of these
> "shortcuts" are *adaptive* and often outperform complex strategies under uncertainty.
> Diagnose whether you're seeing a genuine error or a fit-to-environment heuristic before
> "fixing" it. See SKILL.md §4.

---

## The three foundational heuristics (1974)

### Availability
Judging frequency/probability by how easily instances come to mind. Vivid, recent, or
emotionally charged events feel more probable than they are.
- *Demonstration:* people judge words starting with "r" as more common than words with "r"
  in third position (the former are easier to retrieve); the reverse is true.
- *TAM note:* a single dramatic outage makes a customer over-weight that failure mode in
  planning while ignoring silent, higher-probability risks. Counter with base rates.

### Representativeness
Judging probability by similarity to a stereotype/prototype, neglecting base rates and
sample size.
- *Demonstration:* the **conjunction fallacy** ("Linda the bank teller" — people rate
  "bank teller AND feminist" as more probable than "bank teller," which is logically
  impossible). Also base-rate neglect in the cab/medical-test problems.
- *TAM note:* "this account *looks like* our churned accounts" can be useful pattern-matching
  but ignores the base rate of churn; check the actual rate before acting.

### Anchoring-and-adjustment
An initial value (even arbitrary) anchors the estimate; subsequent adjustment is
insufficient, biasing the final answer toward the anchor.
- *Demonstration:* spinning a wheel of fortune to a random number then asking for the % of
  African nations in the UN — higher wheel values produced higher estimates.
- *TAM note:* **the most useful negotiation effect.** The first number (price, scope, timeline)
  anchors the entire discussion. Open first when you have a defensible anchor; when anchored
  against, explicitly re-anchor before adjusting. *Consider-the-opposite* is the best
  debiaser for anchoring.

---

## High-frequency operator biases

### Confirmation bias
Seeking, interpreting, and weighting evidence to favor a pre-existing belief.
- *TAM note:* the engine behind premature root-cause calls and "the customer is fine" right
  before a churn. Counter: actively seek disconfirming evidence; assign a devil's advocate.

### Hindsight bias ("knew-it-all-along")
After an outcome, overestimating how predictable it was.
- *TAM note:* corrupts postmortems — people retrofit the failure as "obvious." Reconstruct
  what was *knowable at the time*, in timeline order. See `postmortem-writing`.

### Overconfidence (and the planning fallacy)
Confidence intervals are too narrow; predictions too optimistic. The **planning fallacy** is
the special case for task duration/cost — plans systematically underestimate.
- *Demonstration:* "90% confident" ranges contain the truth far less than 90% of the time.
- *TAM note:* migration/upgrade timelines and capacity forecasts are chronically optimistic.
  Counter with **reference-class forecasting** + **premortem**; widen the interval.

### Status-quo bias & the default effect
Disproportionate preference for the current state and for the pre-selected option.
- *Demonstration:* organ-donation and 401(k) enrollment rates swing dramatically with
  opt-in vs opt-out defaults (same choice, different default).
- *TAM note:* "let's keep what we have" is often status-quo bias, not a reasoned no. Make the
  better option the default, cut switching friction, or impose a decision deadline so inaction
  isn't the free win.

### Sunk-cost fallacy (escalation of commitment)
Continuing because of unrecoverable past investment rather than future expected value.
- *Demonstration:* the Concorde fallacy; finishing a bad movie "because I paid."
- *TAM note:* "we've already invested 18 months in this cluster design / this vendor" — reframe
  to the **marginal** decision from today; make past spend explicitly irrelevant to the forward call.

---

## Other named effects you'll meet

- **Framing effect** — equivalent options framed as gains vs losses flip the choice (see
  prospect theory in SKILL.md §3). "95% success" vs "5% failure."
- **Endowment effect** — valuing something more once you own it (often attributed to loss
  aversion; the attribution is contested — see `references/replication-status.md`).
- **Base-rate neglect** — ignoring prior probabilities; a representativeness symptom.
- **Recency / primacy** — order effects in sequential information.
- **Affect heuristic** — substituting "how do I feel about it" for a probability/risk judgment.
- **Halo effect** — one salient positive trait colors the whole evaluation (a polished demo
  inflating the perceived reliability of the product).
- **Gambler's fallacy / hot-hand** — misreading randomness in streaks.

For each, the operator pattern is the same: **name the mechanism, then design the decision
environment around it** — never tell the customer they're being irrational.
