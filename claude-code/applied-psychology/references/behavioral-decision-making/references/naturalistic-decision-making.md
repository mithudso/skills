# Naturalistic Decision Making (NDM) & Recognition-Primed Decision (RPD)

The **expert-intuition counterpoint** to the heuristics-and-biases program in the main skill.
Heuristics-and-biases studies how naive subjects err on artificial lab tasks; NDM studies how
*experienced* practitioners succeed in messy, high-stakes, real-world settings. Both are right —
the truth is conditional on the **environment** and the **decider's experience in it**. This file
gives you the model, the field, the famous Kahneman–Klein truce on when intuition is trustworthy,
and the operator rule for trusting (or distrusting) a fast expert read.

> **Read alongside §4 of SKILL.md** (bounded rationality / ecological rationality). NDM is the field
> evidence for Gigerenzer's claim that simple, recognition-based strategies are often *adaptive*, not
> defective. NDM and fast-and-frugal heuristics are cousins; both push back on the deficit framing.

---

## 1. The Recognition-Primed Decision model (Gary Klein)

RPD describes how experts actually decide under time pressure, uncertainty, high stakes, and shifting
conditions — settings where the classical "list all options, score each, pick the max" procedure is
impossible. Klein and colleagues built it from field studies of **fireground commanders**, then
replicated the pattern in critical-care nurses, military commanders, pilots, nuclear-plant operators,
and chess masters (Klein, Calderwood & Clinton-Cirocco, 1986/1993; Klein, *Sources of Power*, 1998).

**The core loop: cue → pattern → action-script.**

1. **Cue recognition.** The expert reads the situation and recognizes it as *typical* (or atypical) —
   a pattern match against thousands of stored cases. Recognition delivers four things at once:
   plausible **goals**, relevant **cues** (what to attend to), **expectancies** (what should happen
   next — violated expectancies are the tripwire that says "this case is not typical"), and a set of
   **typical actions**.
2. **Action-script retrieval.** Recognition of the pattern *primes* the first workable course of
   action — not a menu, just the action that has worked for this pattern before.
3. **Mental simulation.** The expert runs the candidate action forward in imagination ("if I do this,
   then…"). If the simulation reveals no fatal flaw, they act. If it does, they tweak the action or
   reject it and simulate the next one.

**Three variants (increasing difficulty):**
- **Simple match** — recognize the situation, execute the first adequate action. No comparison at all.
- **Diagnose the situation** — the case is not immediately typical; the expert does feature-matching
  or story-building (sensemaking) until it resolves into a recognizable pattern, then acts.
- **Evaluate a single course of action** — the action is clear but its outcome is uncertain, so the
  expert mentally simulates it before committing.

**The load-bearing contrast: serial vs concurrent option evaluation.** Classical decision theory
assumes **concurrent comparison** — you generate the full option set and evaluate them side by side
against common criteria (this is the normative model; see `da-33-prescriptive-analytics`). RPD is
**serial / satisficing** — the expert generates ONE option (the first the pattern primes),
mentally simulates it, and either takes it or moves to the next. Klein's field finding: skilled
decision-makers rarely compared options at all; the first option they considered was usually
workable, and generating it via recognition was the whole point. This is satisficing (Simon, §4 of
SKILL.md) realized in expert cognition, and it is *faster and often as good* as comparison when the
environment is learnable and time is short.

---

## 2. The NDM field, macrocognition, and sensemaking

**Naturalistic Decision Making** is the research movement (origin: a 1989 Dayton, Ohio conference;
Klein, Orasanu, Calderwood & Zsambok, *Decision Making in Action*, 1993). Its program: study how
proficient people make decisions in *real* conditions — time pressure, ambiguity, high stakes,
shifting goals, team coordination — rather than on context-free lab gambles. RPD is its flagship model.

**Macrocognition** — the cognitive functions experts actually use in the wild, as opposed to the
"microcognition" of controlled lab paradigms. The recognized macrocognitive functions: **decision
making, sensemaking, problem detection, planning, adapting/replanning, and coordinating** (Klein,
Ross, Moon, Klein, Hoffman & Hollnagel, 2003; Klein & Crandall).

**Sensemaking** — the process of fitting data into a frame (an explanatory structure) and fitting the
frame to the data, iterating as new cues arrive. Klein's **data/frame theory** (Klein, Phillips, Rall
& Peluso, 2007): a frame shapes which data you notice; surprising data forces you to question,
elaborate, or replace the frame. This is the front half of the RPD loop — getting to a recognizable
pattern when the situation is *not* immediately typical. Operator analogue: the moment a customer
signal doesn't fit your mental model of the account, that *violated expectancy* is the cue to re-frame,
not to force the old story.

---

## 3. RPD vs classical decision theory (the explicit contrast)

| | Classical / normative (SEU) | Recognition-Primed Decision |
| --- | --- | --- |
| Where it lives | `da-33-prescriptive-analytics`, lab gambles | Real fire grounds, ORs, command posts |
| Option generation | Full set, exhaustive | One at a time, recognition-primed |
| Evaluation | Concurrent comparison vs common criteria | Serial mental simulation of one option |
| Decision rule | Maximize expected utility | Satisfice — take first option that works |
| What makes it good | Logical coherence | Pattern validity + deep experience |
| Failure mode | Intractable under time pressure; ignores expertise | Fails in low-validity / unfamiliar domains |

NDM does **not** claim experts are bias-free, nor that comparison is never worth it. It claims that
in **high-validity, well-practiced** domains the recognition route is fast, frugal, and accurate — and
that lab studies on novices systematically *understate* expert competence. The heuristics-and-biases
program and NDM are complementary lenses, not rivals; which one applies depends on the environment.

---

## 4. When is intuition trustworthy? Kahneman–Klein (2009)

The definitive reconciliation. Daniel Kahneman (skeptic of expert intuition; heuristics-and-biases)
and Gary Klein (champion of expert intuition; NDM) ran a ~6-year **adversarial collaboration** and
published "Conditions for Intuitive Expertise: **A Failure to Disagree**" (*American Psychologist*,
2009, 64(6), 515–526). They converged on **two conditions** that must *both* hold before an intuitive
judgment deserves trust:

1. **High-validity environment** — the domain contains stable, learnable regularities that genuinely
   predict outcomes. (Firefighting, anesthesiology, chess, reading a familiar customer = high validity.
   Long-range political/economic/stock-pick forecasting = **low/zero validity** — no stable cues to
   learn, so "expert intuition" there is illusory regardless of confidence or experience.)
2. **Adequate opportunity to learn** — the expert practiced in that environment with **rapid, clear
   feedback** over enough repetitions to internalize the regularities. Practice without prompt, valid
   feedback breeds confident-but-wrong intuition (the radiologist who never learns the biopsy result;
   the forecaster whose calls are never scored).

**The trap they jointly flagged:** *subjective confidence is not a valid cue to accuracy.* People feel
equally certain in high- and low-validity domains. Fluency and confidence are produced by System 1
whether or not the underlying environment supports skill. So you cannot use "I'm sure" to decide
whether to trust your gut — you must audit the *environment* and the *feedback history* instead.

This is the bridge to the main skill: a heuristic/intuition is trustworthy exactly when it is
**ecologically rational** (§4) — matched to a high-validity, well-learned environment. Outside those
conditions, fall back on the debiasing procedures in `references/debiasing-and-application.md`.

---

## 5. The premortem (Klein) — where NDM meets debiasing

Klein originated the **premortem** (developed ~1989–1991; popularized in "Performing a Project
Premortem," *Harvard Business Review*, Sept 2007). Before committing to a plan, the team is told to
*assume it has already failed catastrophically* and then write down every reason why. This inverts the
usual "will it work?" framing into "it didn't work — explain how," which licenses dissent and surfaces
risks that optimism and groupthink suppress.

- **Mechanism: prospective hindsight.** Imagining an outcome as already-certain makes people generate
  concrete causes far more readily. The technique builds on Mitchell, Russo & Pennington (1989), "Back
  to the Future," which reported prospective hindsight increased the number of correctly identified
  reasons for a future outcome by ~30%. (Treat the "30%" as the originating finding, not a guaranteed
  universal effect size.)
- **How it complements debiasing.** The premortem is the team-scale, emotionally-licensed version of
  *consider-the-opposite* — the best-evidenced general debiaser (see SKILL.md §7 and
  `references/debiasing-and-application.md`). It directly attacks overconfidence, the planning fallacy,
  and the optimism that RPD's fast pattern-match can carry when the pattern is *wrong*.
- **The deeper point:** the same Gary Klein who showed experts *should* trust fast recognition in
  high-validity settings also built the premortem to *slow experts down* when the stakes are high and
  the situation may be atypical. NDM is not "always trust your gut" — it is "know when your gut earned
  the right to be trusted, and install a System-2 check when it hasn't."

---

## 6. Operator note — fast expert read vs slow down and debias

Use this as a gate before you act on a gut call (renewal risk, escalation severity, architecture
recommendation, "this account is fine / in trouble"):

**Trust the fast expert read when ALL hold:**
- The domain is **high-validity** — there are real, repeating patterns (you've seen this account shape,
  this failure signature, this buying behavior many times).
- You personally have **logged the reps with feedback** — you've seen how cases like this resolved,
  not just observed them once.
- **Expectancies are being met** — the situation is unfolding the way the recognized pattern predicts;
  no surprising cue has fired.
- It's **reversible or low-stakes**, or time pressure genuinely forecloses deliberation.

**Slow down and run a System-2 check (premortem / consider-the-opposite / reference-class) when ANY
hold:**
- **Low-validity domain** — long-horizon forecast, novel market, one-off bet with no learnable base
  rate. Confidence here is a feeling, not evidence.
- **Novel or surprising situation** — a violated expectancy fired; the pattern match may be a
  *false* match (recognition mis-fires on superficial similarity — the failure mode of
  representativeness, §2 of SKILL.md).
- **High stakes and irreversible** — large renewal, public escalation, one-way-door architecture call.
- **You can't point to the feedback that trained the intuition** — then it isn't expertise, it's a
  guess wearing a confident face.

The honest stance: do not romanticize intuition *or* reflexively distrust it. Audit the environment and
your own feedback history first; let that — not how certain you feel — decide which system gets the call.

---

## Sources

- Klein, G. A. (1998). *Sources of Power: How People Make Decisions.* MIT Press. (RPD model; fireground
  and expert field studies)
- Klein, G., Calderwood, R., & Clinton-Cirocco, A. (1986/1993). "Rapid Decision Making on the Fire
  Ground." (the originating firefighter study; reprinted *J. Cognitive Engineering and Decision Making*)
- Klein, G., Orasanu, J., Calderwood, R., & Zsambok, C. E. (Eds.) (1993). *Decision Making in Action:
  Models and Methods.* (founding NDM volume)
- Kahneman, D., & Klein, G. (2009). "Conditions for Intuitive Expertise: A Failure to Disagree."
  *American Psychologist* 64(6), 515–526. (the adversarial collaboration; high-validity + learnability)
- Klein, G. (2007). "Performing a Project Premortem." *Harvard Business Review*, September 2007.
- Mitchell, D. J., Russo, J. E., & Pennington, N. (1989). "Back to the Future: Temporal Perspective in
  the Explanation of Events." *Journal of Behavioral Decision Making.* (prospective hindsight; premortem
  basis)
- Klein, G., Ross, K. G., Moon, B. M., Klein, D. E., Hoffman, R. R., & Hollnagel, E. (2003).
  "Macrocognition." *IEEE Intelligent Systems.*
- Klein, G., Phillips, J. K., Rall, E. L., & Peluso, D. A. (2007). "A Data/Frame Theory of
  Sensemaking." (in *Expertise Out of Context*, ed. Hoffman)
