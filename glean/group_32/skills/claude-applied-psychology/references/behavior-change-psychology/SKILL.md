<!-- FOLDED SPOKE of the applied-psychology hub. This is the full former standalone skill; any "use the <X>" / SKIP pointer below now refers to a sibling reference under applied-psychology/references/<X>/. -->

---
name: behavior-change-psychology
description: >-
  Motivation and behavior-change psychology applied to technical account
  management, customer adoption, and enablement. Five evidence-based models
  in one coherent toolkit — Self-Determination Theory (autonomy/competence/
  relatedness, intrinsic vs extrinsic, the overjustification effect), the Fogg
  Behavior Model (B=MAP) and Tiny Habits, the Transtheoretical / Stages-of-Change
  model, habit formation (cue-routine-reward, implementation intentions, habit
  stacking), and goal-setting theory (Locke & Latham, SMART, the goal-gradient
  effect) — each mapped to driving durable behavior change in B2B accounts.
  TRIGGER: "how do I get this customer to adopt X", "drive product adoption /
  feature adoption", "they signed but aren't using it", "champion isn't engaging",
  "build an adoption / enablement plan", "make this usage a habit", "design an
  onboarding nudge", "why did our rollout stall", "incentive / reward is
  backfiring", "set adoption goals", "stages of change for a customer", "motivate
  a stakeholder", "behavior change model", "Fogg B=MAP", "self-determination
  theory", "tiny habits", "habit loop", "implementation intentions", "goal-gradient",
  "overjustification effect", "intrinsic vs extrinsic motivation". SKIP: TAM
  account deliverables, health scoring, QBR/EBR structure, churn-risk frameworks
  with no behavior-change-psychology angle (use tam-expertise); account health
  scoring algorithms (use account-health-scorer); cohort/retention math,
  activation funnels, or feature-adoption metrics as measurement (use
  da-34-cohort-retention-analytics and da-21-product-analytics); motivational-
  interviewing OARS technique and conversational question design (use
  interview-and-conversational); persuasion/influence for executive decisions
  (use executive-comms); the cognitive biases, prospect-theory framing, and
  choice architecture behind a one-time buyer/customer DECISION — anchoring,
  loss aversion, defaults, sunk-cost, debiasing your own judgment (use
  behavioral-decision-making).
related_skills:
  - tam-expertise
  - account-health-scorer
  - da-applied-and-communication
  - writing-expert
---

> **Folded spoke of the `applied-psychology` hub.** This file was the standalone `behavior-change-psychology` skill; it now lives at `applied-psychology/references/behavior-change-psychology/SKILL.md`.
> Any "use the <X> skill" / "(use ...)" pointer below to another psychology skill now refers to a **sibling reference** under `applied-psychology/references/<X>/` — not a top-level skill.

# Behavior-Change Psychology for TAM & Customer Adoption

Five evidence-based models of motivation and behavior change, combined into one
applied toolkit for technical account managers, customer success, and enablement
teams. The job-to-be-done: **a customer bought the product, but people aren't
changing what they actually do.** These models explain why, and what to do about it.

The throughline: adoption is a *behavior-change* problem, not a *feature-awareness*
problem. Telling users about value rarely changes behavior. Engineering motivation,
ability, prompts, habits, and goals does.

---

## Decision table — pick the model for the symptom

| Customer symptom | Primary model | What it tells you to do |
| --- | --- | --- |
| "They love it in demos but never log in." | Fogg (B=MAP) | A prompt is missing, or ability is too low. Add a trigger; shrink the first action. |
| "Mandated by an exec, team resents it." | Self-Determination Theory | Motivation is purely external. Build autonomy + competence to internalize it. |
| "We gave them an incentive and usage *dropped* after it ended." | Overjustification effect (SDT) | An expected, controlling reward crowded out intrinsic interest. Switch to informational feedback. |
| "Champion is enthusiastic; the wider team is indifferent." | Transtheoretical (Stages of Change) | Different stakeholders are in different stages. Stage-match your asks. |
| "They adopt during onboarding, then drift back to the old tool." | Habit formation | No stable cue-routine-reward loop formed. Anchor the new behavior to an existing routine. |
| "Usage is vague and stalled; no momentum." | Goal-setting (Locke & Latham) | No specific, challenging, committed goal with feedback. Set one; show progress. |
| "Engagement spikes near a milestone, then flatlines." | Goal-gradient effect | Effort rises near a visible finish line. Add intermediate milestones and progress markers. |

Most real accounts need **two or three models layered**. A typical adoption play:
set a clear goal (goal-setting) → make the first action trivially easy and prompted
(Fogg) → repeat it against a reliable cue until automatic (habit) → support
autonomy/competence so motivation internalizes (SDT) → stage-match the asks across
the buying group (TTM). See *Combined adoption playbook* below.

---

## Model 1 — Self-Determination Theory (Deci & Ryan)

**Core claim.** Durable, high-quality motivation comes from satisfying three innate
psychological needs. Frustrate them and you get compliance at best, resistance at
worst. [Ryan & Deci 2000](https://selfdeterminationtheory.org/SDT/documents/2000_RyanDeci_SDT.pdf)

The three basic psychological needs:
- **Autonomy** — feeling the behavior is self-chosen, not coerced. (Not independence; it's volition.)
- **Competence** — feeling effective, making visible progress, experiencing mastery.
- **Relatedness** — feeling connected to and valued by others who matter.

**The motivation continuum (Organismic Integration Theory).** Extrinsic motivation
is not one thing. It runs along a continuum of *internalization*:

`amotivation → external regulation → introjected → identified → integrated → intrinsic`

- **External regulation** — "I use it because my manager checks." (rewards/punishments)
- **Introjected** — "I'd feel guilty if I didn't." (internalized pressure, ego)
- **Identified** — "This genuinely helps me hit my goals." (personally valued)
- **Integrated** — "This is how my team works now." (aligned with identity)
- **Intrinsic** — "I find this genuinely satisfying to use."

The TAM goal is to move stakeholders *rightward* — from external to identified/integrated.
You rarely get true intrinsic motivation for enterprise software, and you don't need
it; **identified/integrated regulation is the durable target.** [Ryan & Deci 2020 (CEP)](https://selfdeterminationtheory.org/wp-content/uploads/2020/04/2020_RyanDeci_CEP_PrePrint.pdf), [Organismic Integration Theory overview](https://psychologyfanatic.com/organismic-integration-theory/)

**The overjustification effect.** Giving an *expected, tangible, controlling* reward
for an activity someone already finds interesting can *reduce* their intrinsic
motivation — the perceived reason for acting shifts from internal interest to the
external reward. In Deci's classic 1971 study, paid puzzle-solvers played less in
free time than unpaid ones. The mechanism: a shift in *perceived locus of causality*
from internal to external. Crucially, **informational** rewards (that signal
competence) are far less corrosive than **controlling** ones. [Overjustification effect (Wikipedia)](https://en.wikipedia.org/wiki/Overjustification_effect), [The Decision Lab](https://thedecisionlab.com/biases/overjustification-effect), [Deci, Koestner & Ryan 2001 meta-analysis](https://journals.sagepub.com/doi/10.3102/00346543071001001)

**TAM application:**
- **Autonomy:** Offer paths, not mandates. "Here are three ways teams like yours
  roll this out — which fits you?" Co-design the adoption plan with the champion
  rather than handing them one.
- **Competence:** Sequence enablement from quick wins to mastery. Celebrate visible
  progress (dashboards, certifications, "you're now in the top 10% of power users").
- **Relatedness:** Connect users to a community/peer cohort; make the CSM a trusted
  partner, not an auditor. User groups and customer advisory boards satisfy relatedness.
- **Avoid the overjustification trap:** Don't bolt cash/swag bounties onto a feature
  users already like — when the bounty ends, usage can fall *below* baseline. Prefer
  *informational* recognition ("your team's query latency dropped 40%") over
  *controlling* incentives ("log in 5× this week for a gift card").

> SDT is *the* lens for the "exec mandated it and the team resents it" problem and
> for diagnosing when an incentive program will backfire.

---

## Model 2 — Fogg Behavior Model (B = MAP) & Tiny Habits

**Core formula.** `B = MAP`: a **B**ehavior happens only when **M**otivation,
**A**bility, and a **P**rompt converge *at the same moment*. Remove any one and the
behavior doesn't occur. Developed by BJ Fogg (Stanford Behavior Design Lab), expanded
in *Tiny Habits* (2019). [BJ Fogg — behaviormodel.org](https://www.behaviormodel.org/), [Behavioral Scientist](https://www.thebehavioralscientist.com/articles/fogg-behavior-model)

**The action line (compensatory relationship).** Motivation and Ability trade off
along a curved threshold. If a behavior is *easy enough*, even low motivation crosses
the line. If motivation is *sky-high*, people will tolerate a hard behavior. Because
**motivation fluctuates (the "motivation wave") and is unreliable, the durable lever
is almost always Ability — make the behavior easier.** [The Decision Lab — Fogg](https://thedecisionlab.com/reference-guide/psychology/fogg-behavior-model)

**Six ability factors** (what "simplicity" means — the scarcest one is your bottleneck):
time, money, physical effort, mental effort (cognitive load), social deviance, and
non-routine (unfamiliarity). [Triple Whale guide](https://www.triplewhale.com/blog/fogg-behavior-model)

**Three prompt types** (match the prompt to where the user is):
- **Spark** — raises motivation when motivation is low (highlight a benefit/risk).
- **Facilitator** — raises ability when motivation is high but the task feels hard (templates, defaults, "do it for me").
- **Signal** — a simple reminder when motivation *and* ability are already sufficient.

**Tiny Habits recipe.** `After I [ANCHOR MOMENT], I will [TINY BEHAVIOR]. Then I
celebrate.` The tiny behavior must take <30 seconds and need almost no motivation;
the celebration wires in positive emotion to accelerate habit formation. [EasyHabits — Tiny Habits](https://www.easyhabits.io/blog/tiny-habits-bj-fogg), [BeWay — B=MAP & Tiny Habits](https://blog.beway.com/en/bmap-big-changes-with-tiny-habits/)

**TAM application:**
- Diagnose stalled adoption with three questions: *Is there a prompt? Is the first
  action easy enough? Is there enough motivation for how hard it is?* Fix the missing
  element — usually start with prompt and ability before trying to pump motivation.
- **Shrink the first action.** Don't ask a user to "migrate their workflow." Ask them
  to "run one saved query." Time-to-first-value should be minutes, not weeks.
- **Engineer prompts into existing systems:** in-app nudges, Slack/email triggers tied
  to a real moment, calendar holds — not a generic "please use the product" email.
- **Pre-load ability:** ship templates, sample dashboards, sensible defaults, and
  "white-glove" first-config so the user's effort approaches zero.

> Fogg is the first tool to reach for on *any* "they're not using it" symptom because
> it's the fastest diagnostic: prompt? ability? motivation?

---

## Model 3 — Transtheoretical Model (Stages of Change)

**Core claim.** Behavior change is a *process through stages*, not a single event.
People at different stages need different interventions. Developed by Prochaska &
DiClemente (1980s). [Transtheoretical model (Wikipedia)](https://en.wikipedia.org/wiki/Transtheoretical_model), [Simply Psychology](https://www.simplypsychology.org/transtheoretical-model.html)

**The stages:**
1. **Precontemplation** — not even considering change; unaware of the problem/value. (No intent within ~6 months.)
2. **Contemplation** — aware of the upside, but ambivalent; weighing pros vs cons. (Intends to act within ~6 months.)
3. **Preparation** — committed and planning; has an action plan, often a small first step taken.
4. **Action** — actively doing the new behavior (recent, <6 months).
5. **Maintenance** — sustaining the behavior, guarding against relapse (6+ months).
   (A sixth stage, **termination**, where relapse risk is gone, is often omitted in practice.)

**Three supporting constructs:**
- **Decisional balance** — the perceived pros vs cons of changing. Progress happens
  as pros come to outweigh cons.
- **Self-efficacy** — confidence the user can sustain the behavior under temptation.
- **Processes of change** — the *how*: ten techniques (consciousness-raising, etc.)
  that move people between stages. The key practice is **stage-matching** — using the
  right technique for the stage. Stage-matched interventions roughly *double* the odds
  of progressing to action. [Stages of Change Theory — StatPearls/NIH](https://www.ncbi.nlm.nih.gov/books/NBK556005/), [ScienceDirect — TTM overview](https://www.sciencedirect.com/topics/nursing-and-health-professions/transtheoretical-model)

**TAM application — stage-matched plays across the buying group:**
- **Precontemplation:** Don't pitch features. Raise awareness of the cost of the
  status quo. Share peer benchmarks and a relevant "what others are losing by not
  doing this" narrative.
- **Contemplation:** Resolve ambivalence. Build a decisional-balance sheet *with* them;
  address the cons explicitly (migration cost, retraining fear). This is where
  motivational-interviewing change-talk helps (→ `writing-expert` (references/interview-and-conversational.md)).
- **Preparation:** Co-build a concrete 30-60-90 adoption plan and a tiny first step.
- **Action:** Remove friction, provide hands-on enablement, celebrate early wins.
- **Maintenance:** Guard against relapse to the old tool — QBRs, health monitoring,
  expansion to new use cases, reinforce habits.
- Map each stakeholder (champion, end users, economic buyer) to a stage; they're
  rarely all in the same one. The most common mistake is pitching Action-stage tactics
  to a Precontemplation-stage audience.

**Known limitation (state it honestly).** TTM is widely critiqued: the stage time-
boundaries (6 months, 30 days) are arbitrary, people skip/regress through stages
non-sequentially, and systematic reviews question whether stage-based interventions
beat non-staged ones. Use it as a *diagnostic vocabulary for readiness*, not a rigid
ladder. For a factors-based alternative, see **COM-B / Behavior Change Wheel** (Michie
et al.) in *Adjacent models*. [Wikipedia — TTM criticism](https://en.wikipedia.org/wiki/Transtheoretical_model), [The Behavioral Scientist — TTM](https://www.thebehavioralscientist.com/glossary/transtheoretical-model)

---

## Model 4 — Habit Formation

**Core claim.** Durable adoption = behavior that runs automatically, with low
conscious effort, triggered by a context cue. The brain shifts control from the
prefrontal cortex (effortful) to the basal ganglia (automatic) through repeated
cue→behavior→reward cycles. [Simply Psychology — science of habit formation](https://www.simplypsychology.com/articles/science-of-habit-formation)

**The habit loop (Duhigg):** **Cue → Routine → Reward.** A cue triggers a routine;
the routine yields a reward the brain wants to repeat. James Clear's *Atomic Habits*
inserts a fourth element, **craving**, between cue and routine — useful when designing
a habit deliberately (cue → craving → response → reward). [Habi — habit loop](https://habi.app/insights/habit-loop/), [James Clear — habit triggers](https://jamesclear.com/habit-triggers)

**How long it takes.** Lally et al. (UCL, 96 participants) found automaticity forms
in **18-254 days, median ~66** — *not* the mythical 21. Simple behaviors form faster
than complex ones. Plan adoption support over months, not weeks. [Mentalzon — 66 days / Lally](https://mentalzon.com/en/post/7770/66-days-to-build-a-new-habit-why-it%E2%80%99s-not-a-myth-but-real-habit-psychology)

**Two high-impact techniques:**
- **Implementation intentions (Gollwitzer):** `I will [behavior] at [time] in
  [location/situation].` A meta-analysis of 94 studies (8,000+ people) found a
  medium-to-large effect (d ≈ 0.65) on goal attainment. Pre-deciding the *when/where*
  bridges the intention-action gap. [Habit-formation guide (Gollwitzer d=0.65)](https://goalsandprogress.com/habit-formation-complete-guide/)
- **Habit stacking:** `After I [existing habit], I will [new habit].` Anchors the new
  behavior to an existing, well-cued routine so you borrow its trigger and neural
  pathway. (This is the Tiny Habits anchor in habit-loop language.)

**TAM application:**
- **Anchor product usage to an existing work ritual:** "After the daily standup, the
  on-call engineer checks the cluster health dashboard." The standup is the cue; you
  don't have to manufacture a new one.
- **Design the reward to be immediate and visible.** The basal-ganglia loop needs a
  reward *now*: a fast result, a satisfying confirmation, a saved-time indicator.
  Delayed ROI doesn't reinforce a habit; instant feedback does.
- **Use implementation intentions in enablement:** have champions commit, in writing,
  to *when and where* the team will use the feature ("every Monday at 9am, in the
  release-planning meeting, we open the X report"), not just *that* they will.
- **Target a "usage threshold."** Identify the repetition count that predicts
  stickiness (e.g., Slack's ~2,000-messages signal) and drive the account to it
  before declaring onboarding done. [Slack PAI threshold — SaaS onboarding science](https://www.saasfactor.co/blogs/the-science-of-saas-onboarding-a-comprehensive-framework-for-reducing-friction-improving-activation-and-preventing-churn)

---

## Model 5 — Goal-Setting Theory & the Goal-Gradient Effect

**Core claim (Locke & Latham).** Specific, challenging goals produce higher
performance than vague ("do your best") or easy goals — *provided* certain conditions
hold. The most-validated theory in organizational psychology. [Positive Psychology — goal-setting theory](https://positivepsychology.com/goal-setting-theory/), [Mindtools — Locke's theory](https://www.mindtools.com/azazlu3/lockes-goal-setting-theory/)

**Five principles:**
1. **Clarity** — specific, measurable goals beat vague ones.
2. **Challenge** — harder (but attainable) goals drive more effort than easy ones.
3. **Commitment** — the person must own the goal; self-efficacy fuels commitment.
4. **Feedback** — progress feedback is required to adjust effort/strategy.
5. **Task complexity** — for complex tasks, break into sub-goals and provide support/time.

**SMART** (Specific, Measurable, Achievable, Realistic, Time-bound) operationalizes
clarity + challenge + a deadline. Useful as a *checklist*, but it's a packaging of the
principles, not a replacement for them — SMART alone omits commitment and feedback. [Notion — goal-setting theory](https://www.notion.com/blog/goal-setting-theory)

**The goal-gradient effect (Hull; revived by Kivetz, Urminsky & Zheng 2006).** Effort
*accelerates* as people get closer to a goal. In the famous coffee-card study,
customers bought faster the nearer they were to a free drink; raters rated more songs
near the reward. Critically, the **endowed progress effect** shows that giving people
*artificial* initial progress (a punch card pre-stamped 2 of 12 vs a blank 10-stamp
card — same 10 to go) speeds completion, because the finish line feels closer. [Kivetz et al. 2006 — Goal-Gradient Hypothesis Resurrected (PDF)](https://home.uchicago.edu/ourminsky/Goal-Gradient_Illusionary_Goal_Progress.pdf), [Learning Loop — goal gradient](https://learningloop.io/plays/psychology/goal-gradient-effect)

**TAM application:**
- **Set adoption goals as SMART + committed:** not "drive adoption" but "the data team
  runs ≥20 production queries/week on the new cluster by end of Q3," co-owned with the
  champion (commitment), with a weekly usage dashboard (feedback).
- **Decompose big rollouts into sub-goals** (task complexity): per-team, per-use-case,
  per-milestone — each a visible finish line.
- **Exploit goal-gradient in onboarding UX and success plans:** progress bars,
  completion checklists, "3 of 5 onboarding steps done," tiered milestones. Show the
  near finish line. [Lampstellar — goal gradient for customer success](https://www.lampstellar.com/blog/customer-success/power-goal-gradient-effect-elevate-customer-success)
- **Use endowed progress:** start the customer's onboarding checklist *already partly
  complete* ("Account created ✓, SSO configured ✓ — 2 of 6 done") to pull them toward
  completion.

**The dark side (use challenging goals responsibly).** "Goals Gone Wild" (Ordóñez,
Schweitzer, Galinsky, Bazerman 2009) documents that over-aggressive goals can narrow
focus, crowd out unstated priorities, increase risk-taking, and motivate *unethical
shortcuts* — especially when someone is just short of the goal, and when ability is
lacking. Latham & Itzchakov (2024) frame these as boundary-condition failures: the
side effects appear when moderators (notably **ability** and **commitment**) are
ignored. For customers: don't set adoption targets so aggressive they drive vanity
metrics or gaming (e.g., logins with no real work). Pair stretch goals with capability
and genuine commitment. [HBS — Goals Gone Wild (PDF)](https://www.hbs.edu/ris/Publication%20Files/09-083.pdf), [Goal setting (Wikipedia) — dark side & moderators](https://en.wikipedia.org/wiki/Goal_setting)

---

## Combined adoption playbook (how the models compose)

A durable adoption motion layers all five. Run roughly in this order:

1. **Frame the goal** (Goal-Setting): co-set a specific, challenging, committed
   adoption goal with feedback. Decompose into sub-goals with visible finish lines.
2. **Locate readiness** (TTM): map each stakeholder to a stage; stage-match the ask.
   Don't sell Action tactics to a Precontemplation audience.
3. **Engineer the first behavior** (Fogg B=MAP): add a prompt, crush ability barriers,
   shrink the first action to <30 seconds of effort. Reach time-to-first-value fast.
4. **Repeat into a habit** (Habit loop): anchor usage to an existing ritual (habit
   stacking), use written implementation intentions, ensure an immediate reward, and
   drive to the usage threshold over ~weeks-to-months.
5. **Internalize the motivation** (SDT): support autonomy (choices not mandates),
   competence (quick wins → mastery, informational feedback), and relatedness (peer
   community). Move stakeholders from external → identified/integrated regulation.
6. **Reward without backfiring** (overjustification guardrail): prefer informational
   recognition over expected, controlling incentives — especially for behaviors users
   already value. Watch for usage cliffs when an incentive ends.

**Worked example — stalled feature rollout:**
A bank's platform team bought a feature; six weeks in, usage is near zero.
- *Fogg diagnosis:* no in-app prompt, and first use requires a 9-step config (ability
  too low). → TAM ships a pre-built config and a Slack trigger tied to deploys.
- *TTM diagnosis:* champion is in Action, end-users in Precontemplation. → TAM runs a
  status-quo-cost session for end-users, not a feature demo.
- *Goal-setting:* co-set "≥15 deploys/week use the feature by end of next month,"
  weekly dashboard. Onboarding checklist starts 2/6 complete (endowed progress).
- *Habit:* "after each deploy, check the X panel" — anchored to the existing deploy ritual.
- *SDT:* offer three rollout paths (autonomy); spotlight the team's latency win
  (informational competence feedback), *not* a gift-card-per-login scheme.

---

## Anti-patterns

- **Treating adoption as an awareness problem.** More demos and feature emails rarely
  change behavior. Engineer prompts, ability, habits, and goals instead.
- **Pumping motivation while ignoring ability.** Motivation is a wave; it recedes. The
  cheaper, more durable lever is making the behavior easier (Fogg).
- **Bolting controlling incentives onto valued behaviors.** Expected, tangible rewards
  can crowd out intrinsic interest (overjustification); usage often falls *below*
  baseline when they stop. Prefer informational feedback.
- **One message for the whole buying group.** Stakeholders sit in different stages
  (TTM); a single pitch mis-serves most of them.
- **Vague goals ("drive adoption").** No clarity, no challenge, no feedback → no effect.
  Make goals SMART *and* committed.
- **Over-aggressive targets without capability.** Stretch goals minus ability/commitment
  → vanity metrics, gaming, narrowed focus (Goals Gone Wild).
- **Declaring "onboarded" before a habit forms.** Adoption that hasn't crossed a usage
  threshold and become a cued routine relapses to the old tool.
- **Expecting habits in 21 days.** Real automaticity is ~66 days median (18-254); plan
  the success motion accordingly.
- **Treating TTM stages as a rigid, linear ladder.** People skip and regress; stages
  are a readiness vocabulary, not a fixed sequence.

---

## Adjacent models (deliberately out of scope — pointers only)

- **COM-B / Behavior Change Wheel** (Michie, van Stralen & West 2011) — behavior needs
  **C**apability, **O**pportunity, **M**otivation. A factors-based ("why") complement to
  TTM's process ("how") view; the most-cited modern alternative. [COM-B model](https://www.habitweekly.com/models-frameworks/the-com-b-model)
- **Motivational interviewing (OARS)** — the conversational technique for resolving
  ambivalence in Contemplation. → `writing-expert` (references/interview-and-conversational.md).
- **Hook Model** (Eyal) — trigger/action/variable-reward/investment; a consumer-product
  habit-engineering variant of the loop above.
- **Nudge / choice architecture** (Thaler & Sunstein) — defaults and friction design;
  overlaps Fogg's ability lever.
- **Adoption/activation metrics** (aha moment, time-to-value, feature velocity, cohort
  retention) — the *measurement* layer. → `da-applied-and-communication` (references/da-21-product-analytics.md), `da-34-cohort-retention-analytics`.
- **The learning SCIENCE of making the training itself stick** — cognitive load, retrieval
  practice, spacing, interleaving, deliberate practice, the illusion of fluency. This skill
  gets a customer *motivated* to learn and keep practicing; designing the enablement so the
  knowledge actually *retains* is → `learning-and-expertise-psychology`.

---

## References

Self-Determination Theory:
- [Ryan & Deci 2000 — SDT and the facilitation of intrinsic motivation (PDF)](https://selfdeterminationtheory.org/SDT/documents/2000_RyanDeci_SDT.pdf)
- [Ryan & Deci 2020 — Intrinsic & extrinsic motivation from an SDT perspective (PDF)](https://selfdeterminationtheory.org/wp-content/uploads/2020/04/2020_RyanDeci_CEP_PrePrint.pdf)
- [Organismic Integration Theory overview](https://psychologyfanatic.com/organismic-integration-theory/)
- [Overjustification effect — Wikipedia](https://en.wikipedia.org/wiki/Overjustification_effect)
- [Overjustification effect — The Decision Lab](https://thedecisionlab.com/biases/overjustification-effect)
- [Deci, Koestner & Ryan 2001 — reward meta-analysis](https://journals.sagepub.com/doi/10.3102/00346543071001001)

Fogg Behavior Model & Tiny Habits:
- [BJ Fogg — behaviormodel.org](https://www.behaviormodel.org/)
- [The Behavioral Scientist — Fogg Behavior Model](https://www.thebehavioralscientist.com/articles/fogg-behavior-model)
- [The Decision Lab — Fogg Behavior Model](https://thedecisionlab.com/reference-guide/psychology/fogg-behavior-model)
- [Triple Whale — complete guide to M, A, P](https://www.triplewhale.com/blog/fogg-behavior-model)
- [EasyHabits — Tiny Habits explained](https://www.easyhabits.io/blog/tiny-habits-bj-fogg)

Transtheoretical Model:
- [Transtheoretical model — Wikipedia (incl. criticism)](https://en.wikipedia.org/wiki/Transtheoretical_model)
- [Simply Psychology — TTM](https://www.simplypsychology.org/transtheoretical-model.html)
- [Stages of Change Theory — StatPearls / NIH](https://www.ncbi.nlm.nih.gov/books/NBK556005/)
- [ScienceDirect — TTM overview](https://www.sciencedirect.com/topics/nursing-and-health-professions/transtheoretical-model)
- [The Behavioral Scientist — TTM glossary & critique](https://www.thebehavioralscientist.com/glossary/transtheoretical-model)

Habit formation:
- [Simply Psychology — science of habit formation](https://www.simplypsychology.com/articles/science-of-habit-formation)
- [Habi — the habit loop (cue/craving/response/reward)](https://habi.app/insights/habit-loop/)
- [James Clear — habit triggers & stacking](https://jamesclear.com/habit-triggers)
- [Goals & Progress — habit formation (Gollwitzer d=0.65; Lally)](https://goalsandprogress.com/habit-formation-complete-guide/)
- [Mentalzon — 66 days / Lally et al.](https://mentalzon.com/en/post/7770/66-days-to-build-a-new-habit-why-it%E2%80%99s-not-a-myth-but-real-habit-psychology)

Goal-setting & goal-gradient:
- [Positive Psychology — Locke & Latham goal-setting theory](https://positivepsychology.com/goal-setting-theory/)
- [Mindtools — Locke's goal-setting theory](https://www.mindtools.com/azazlu3/lockes-goal-setting-theory/)
- [Notion — the 5 principles & SMART](https://www.notion.com/blog/goal-setting-theory)
- [Kivetz, Urminsky & Zheng 2006 — Goal-Gradient Hypothesis Resurrected (PDF)](https://home.uchicago.edu/ourminsky/Goal-Gradient_Illusionary_Goal_Progress.pdf)
- [Learning Loop — goal-gradient effect](https://learningloop.io/plays/psychology/goal-gradient-effect)
- [Lampstellar — goal-gradient for customer success](https://www.lampstellar.com/blog/customer-success/power-goal-gradient-effect-elevate-customer-success)
- [HBS — Goals Gone Wild (PDF)](https://www.hbs.edu/ris/Publication%20Files/09-083.pdf)
- [Goal setting — Wikipedia (dark side, moderators)](https://en.wikipedia.org/wiki/Goal_setting)

Applied customer adoption / behavior change:
- [SaaSFactor — the science of SaaS onboarding (activation, usage thresholds)](https://www.saasfactor.co/blogs/the-science-of-saas-onboarding-a-comprehensive-framework-for-reducing-friction-improving-activation-and-preventing-churn)

Adjacent:
- [COM-B model & Behavior Change Wheel (Michie et al.)](https://www.habitweekly.com/models-frameworks/the-com-b-model)
