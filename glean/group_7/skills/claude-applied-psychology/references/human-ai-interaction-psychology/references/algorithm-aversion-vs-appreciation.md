# Algorithm Aversion vs. Algorithm Appreciation

> Provenance: authored for the `human-ai-interaction-psychology` skill via /dr.
> Primary sources: Dietvorst, Simmons & Massey (2015, JEP:General; 2018,
> Management Science); Logg, Minson & Moore (2019, OBHDP); and the integrative
> reviews (MIS Quarterly, 2024; Reich/Kaju/Maglio, 2022) reconciling the two.

These are two *opposite* documented effects. Both are real. The applied skill is
predicting **which one** a given person, task, and moment will produce, and
designing for it.

## Algorithm aversion (Dietvorst, Simmons & Massey, 2015)

**Claim:** people *erroneously avoid algorithms after seeing them err* — they
abandon an algorithm faster than they abandon a human forecaster making the same
mistakes, even when the algorithm objectively outperforms the human.

**Design (the canonical study):** participants forecast outcomes (e.g., students'
academic performance, states' airline-passenger numbers). They saw the
statistical model's forecasts, their own, or both, sometimes *witnessing the
model's errors*. Then they chose which to bet on for a bonus.

**Core findings:**
- After seeing the model err, participants chose it far less often — despite the
  model having lower average error than their own judgment. They were "betting
  on the worse forecaster."
- **Error visibility is the trigger**, not the error *rate*. People tolerate
  human error as expected; an algorithm's error is read as evidence the
  algorithm is *broken* and won't improve. They hold algorithms to a higher,
  less forgiving standard.
- The aversion is partly affective/expectational: violated expectations of
  near-perfection, and a belief that algorithms can't learn or handle
  particulars the way a human can.

**The fix — "Overcoming Algorithm Aversion" (Dietvorst et al., 2018):** giving
people even a *small amount of control* — the ability to **adjust the
algorithm's output**, even by a capped/tiny amount — substantially increased
their willingness to use it and their satisfaction, with little cost to (and
sometimes a gain in) accuracy. **Perceived agency over the algorithm restores
reliance.**

## Algorithm appreciation (Logg, Minson & Moore, 2019)

**Claim:** in many settings people *prefer algorithmic to human judgment* — they
weight identical advice **more** heavily when told it came from an algorithm.

**Design:** six experiments; advice attributed to "an algorithm" vs. "other
people." Reliance measured by **Weight-On-Advice (WOA)** — how far participants
moved their estimate toward the advice after seeing it (0 = ignore, 1 = fully
adopt).

**Core findings:**
- Across numeric estimates (e.g., a person's weight from a photo) and forecasts
  (song popularity, romantic attraction), people put **higher WOA on
  algorithmic advice** than on human advice — "algorithm appreciation."
- **Boundary conditions (where appreciation shrinks or flips to aversion):**
  - **Choosing between the algorithm and one's OWN judgment** (vs. an external
    human advisor) — people cling to their own view (self-bias).
  - **Domain expertise** — experts/forecasters discounted the algorithm and
    relied more on themselves.
  - Appreciation was measured **before** seeing the algorithm err (pre-feedback),
    which is precisely where Dietvorst's aversion has not yet been triggered.

## Reconciliation: which effect fires when?

The two literatures are not contradictory — they describe **different points on
the same moderating axes.** The empirical landscape is admittedly fragmented
(MISQ 2024 review), but these moderators are the reliable predictors:

| Moderator | → Appreciation (over-weight AI) | → Aversion (discount AI) |
| --- | --- | --- |
| **Error visibility** | No error seen yet (pre-feedback) | A *visible* algorithmic error |
| **Task type** | Objective / numeric / quantifiable | Subjective, moral, "human," taste-based |
| **User expertise** | Non-expert / lay user | Domain expert |
| **Reference advisor** | "Algorithm vs. another person" | "Algorithm vs. **my own** judgment" |
| **Control over output** | (n/a) | Eased when user can adjust the output |
| **Recovery** | — | Aversion softens if the algorithm visibly *improves/learns* |
| **Stakes / accountability** | Lower personal accountability | High accountability + fear of blame for an algorithmic error |

Mechanistic read: **appreciation** is often a heuristic ("data/algorithms are
objective and precise"); **aversion** is triggered by *salient error* plus
*domain ego* plus the *fear of being blamed for deferring to a machine*.

## Operator playbooks

**If your failure mode is AVERSION (under-reliance) — common with expert
customers, or right after a visible miss:**
1. **Add adjustability** — let users edit/tune the suggestion (the strongest
   evidence-based lever).
2. **Set expectations before the first error** — "it will be wrong roughly X% of
   the time, here's the kind of case it misses" — so a visible miss confirms
   rather than violates expectations.
3. **Reframe as advisor, not replacement** — preserve the expert's agency and
   final say.
4. **Show the track record against their own baseline** — make the net-better
   comparison concrete (win-rate vs. status quo), not abstract.
5. **Surface learning/improvement** — show that the system updates; aversion is
   partly the belief that algorithms can't get better.

**If your failure mode is APPRECIATION-as-OVER-RELIANCE — common with
non-experts, objective-looking tasks, confident output:**
1. Treat fluent/precise output as a **risk cue**, not reassurance — it inflates
   WOA regardless of correctness.
2. Apply the reliance interventions in
   `reliance-interventions-and-design.md` (commit-first, friction on
   likely-error cases, honest uncertainty).
3. Instrument **over-reliance rate** (followed-when-wrong), not just agreement.

**Diagnosing a customer in the wild — ask:**
- Are they experts in this task? (expert → aversion-leaning)
- Have they *seen* the tool fail visibly yet? (yes → aversion spike)
- Is the task objective/numeric or subjective/"human"? (subjective → aversion)
- Are they comparing the AI to *their own* judgment or to a human colleague?
  (own → self-bias/aversion)
- Can they adjust the output? (no → aversion risk; yes → reliance rises)
