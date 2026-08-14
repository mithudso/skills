# Trust, Rapport & Psychological Safety — Detail Reference

Deep detail behind `SKILL.md`. Three literatures: interpersonal trust, rapport, psychological safety. 27 sources; primary papers from AMR, AMJ, ASQ, JPSP, JAP, Risk Analysis, Psychological Inquiry.

---

## 1. Interpersonal Trust Models

### Mayer, Davis & Schoorman (1995) — Integrative Model (ABI)
*Academy of Management Review* 20(3), 709-734. ~50k citations.

Trust = "the willingness of a party to be vulnerable to the actions of another party" based on positive expectations. Three antecedents of **trustworthiness**:
- **Ability** — domain-specific skills, competence, expertise. (Note: domain-specific — someone trusted to architect a cluster is not automatically trusted on contracts.)
- **Benevolence** — the trustee wants to do good for the trustor beyond self-interest; genuine care.
- **Integrity** — the trustee adheres to a set of principles the trustor finds acceptable (honesty, consistency, fairness).

Separate trait: **propensity to trust** — the trustor's general disposition to trust; matters most *before* relationship-specific information exists.

Key sequencing insight: **integrity perceptions form fast** (early, from reputation and small signals); **benevolence perceptions grow slowly** with accumulated experience. Ability can be judged quickly within a domain.

Operator application: On a new account, you cannot manufacture benevolence on day one — but you can establish integrity (keep small promises, be transparent about limits) and ability (demonstrate competence fast) immediately.

### McAllister (1995) — Cognitive vs Affective Trust
*Academy of Management Journal* 38(1), 24-59. Field study, 194 managers.
- **Cognition-based trust** — rests on evidence: reliability, competence, track record. "I trust them because the evidence says they're capable and dependable."
- **Affect-based trust** — rests on emotional bonds, reciprocated care, genuine concern. "I trust them because we have a relationship."

The two are distinct but linked; a baseline of cognitive trust is typically a prerequisite before affective trust develops. Maps onto Mayer: cognitive ≈ ability/integrity evidence; affective ≈ benevolence felt over time.

### The Trust Equation — Maister, Green & Galford, *The Trusted Advisor* (2000)
**Trustworthiness = (Credibility + Reliability + Intimacy) / Self-Orientation**
- **Credibility** — words; can they be believed? (expertise, accuracy, honesty about what you know)
- **Reliability** — actions; do they keep promises, consistently? (the compounding variable — every kept micro-commitment adds up)
- **Intimacy** — emotional safety; is it safe to confide in them? (the differentiator — makes the customer share the *real* problem)
- **Self-Orientation** (denominator) — focus on self (quota, ego, agenda) vs the client. The single biggest trust-killer and the **highest-leverage variable**: halving self-orientation does more than doubling any numerator term.

Crosswalk to academic models: Credibility + Reliability ≈ ability / cognitive trust; Intimacy ≈ affective trust; low Self-Orientation ≈ benevolence.

---

## 2. Swift Trust in Temporary Teams

### Meyerson, Weick & Kramer (1996) — "Swift Trust and Temporary Groups"
In Kramer & Tyler (eds), *Trust in Organizations*, Sage, 166-195.

Temporary teams (incident bridges, new vendor relationships, consulting engagements, cross-functional task forces, film crews) must act interdependently from day one without time to build trust through repeated interaction. **Swift trust** resolves this: members import trust **presumptively** from **category/role-based cues** — professional roles, credentials, institutional reputation, shared professional standards — rather than personal history. The group "acts as if trust were present," then verifies and adjusts.

Critical property: swift trust is **fragile** — maintained by action and continued task engagement, with no relational ballast, so a single confirming or disconfirming event swings it sharply.

Operator application: A first TAM engagement, a new escalation bridge with an unfamiliar SRE, or a kickoff with a new champion all run on swift trust. Borrow credibility from your role, references, and MongoDB's standards; then protect it because one early miss does outsized damage.

### Lewicki & Bunker (1996) — Stages of Trust Development
Same Sage volume. Trust deepens through three sequential levels:
1. **Calculus-based** — cost/benefit, fear of sanction. Fragile, transactional. Breaks easily.
2. **Knowledge-based** — predictability from accumulated experience; you can forecast the other's behavior.
3. **Identification-based** — shared values, mutual understanding; parties can act for one another. Most resilient; absorbs shocks.

Operator application: Most early customer relationships sit at calculus/early-knowledge — treat them as breakable. Long-tenured strategic accounts can reach identification-based trust, which is what lets them weather an outage without churning.

---

## 3. Trust Violation & Repair

### Kim, Ferrin, Cooper & Dirks (2004) — "Removing the Shadow of Suspicion"
*Journal of Applied Psychology* 89(1), 104-118. The central repair finding:
- **Competence-based violation** (a mistake, missed SLA, bad recommendation): **apology repairs trust better than denial** — people believe competence can improve and a single lapse isn't diagnostic of permanent inability.
- **Integrity-based violation** (a lie, broken commitment, withholding): **denial repairs trust better than apology — when you are actually innocent** — because integrity violations are seen as more diagnostic and stable; apologizing (admitting guilt) for an integrity breach is especially damaging.

Follow-up nuance: apology beats denial when guilt is later revealed; denial beats apology when innocence is later revealed. So: **deny only if true**, because a denial later exposed as false is catastrophic.

### Reticence (Ferrin, Kim et al.)
Staying silent / non-committal is generally **worse than either a clear apology or a clear denial** — it leaves suspicion unresolved. Do not "go quiet" after a trust hit.

### Slovic (1993) — The Asymmetry Principle
"Perceived Risk, Trust, and Democracy," *Risk Analysis* 13(6). Trust is "created slowly but destroyed in an instant by a single mishap." A manifestation of negativity bias: trust-destroying events are more visible, weighted as more credible, and a single bad event can erase many good ones.

Operator application: Asymmetry justifies over-investing in **prevention** and **rapid recovery**. The math of trust is not linear — you cannot average your way back. Caught early + owned fast is the only cheap path.

#### Repair decision aid
1. Classify the breach: competence or integrity?
2. Competence → **apologize, take responsibility, show the concrete fix + prevention.**
3. Integrity, you are innocent → **deny clearly with evidence**; do not apologize away guilt you don't own.
4. Integrity, you are at fault → apology + visible remediation; expect it to be slow (integrity violations are stickiest).
5. Never go silent (reticence is worst).
6. Honor the asymmetry: act faster and more visibly than feels proportionate.

---

## 4. Rapport & Relationship Dynamics

### Tickle-Degnen & Rosenthal (1990) — Three Components of Rapport
*Psychological Inquiry* 1(4), 285-293. Rapport =
- **Mutual attentiveness** — focused, engaged mutual attention.
- **Positivity** — warmth, friendliness, mutual liking.
- **Coordination** — synchrony, balance, "being in sync."

Weighting shifts over time: **early** = positivity + attentiveness dominate; **later** = coordination + attentiveness matter most (positivity is assumed). So a first meeting needs visible warmth; a long relationship needs demonstrated sync.

### Byrne (1971) — Similarity-Attraction (The Attraction Paradigm)
More shared attitudes/values → more attraction, via a reinforcement-affect model (similarity validates one's worldview = rewarding). **Proportion** of similar attitudes, not absolute count, drives liking.

Operator application: Find and surface *genuine* common ground (technical philosophy, prior tooling, problem framing) — not manufactured agreement.

### Zajonc (1968) — Mere-Exposure Effect
Repeated exposure increases liking, even without conscious recognition. Curve is positive but **decelerating** — the first few (positive) exposures matter most.

Operator application: Frequency and consistency of positive contact build affinity. A short, reliable weekly touchpoint can beat a rare big meeting for relationship warmth.

### Chartrand & Bargh (1999) — The Chameleon Effect
*JPSP* 76(6), 893-910. People nonconsciously mimic others' postures, gestures, mannerisms; this mimicry **increases liking and interaction smoothness**. Synchrony is rapport's behavioral signature.

**Caution:** the natural effect is unconscious. **Deliberate, exaggerated mimicry backfires if detected as inauthentic.** Let synchrony emerge; do not perform it. Matching communication pace/medium/formality is the safe, deliberate version.

### Gottman — The 5:1 "Magic Ratio"
Stable/satisfied couples maintain ~5 positive interactions per 1 negative during conflict (~20:1 outside conflict); below 5:1 predicts dissolution. Gottman predicted divorce with ~90%+ accuracy from short observations. Extended to teams: highest performers run positive-heavy (~5-6:1).

**Caveat (knowledge gap):** the specific *team* threshold rests partly on the Losada line of work, whose mathematical model was critiqued (Brown, Sokal & Friedman 2013). The **directional** claim (high performers run positive-heavy) is robust; treat any exact team number as a **heuristic, not a law**.

### Gable & Reis (2004) — Capitalization / Active-Constructive Responding (ACR)
*JPSP* 87(2), 228-245. How you respond to someone's **good news** predicts relationship quality more than how you respond to bad news. Four response styles to a shared win:
- **Active-Constructive** — enthusiastic, engaged, asks them to elaborate. **The only style that builds the bond.**
- Passive-Constructive — understated ("nice").
- Active-Destructive — points out downsides ("but won't that mean more on-call?").
- Passive-Destructive — ignores / changes the subject.

ACR links to higher intimacy, satisfaction, commitment, trust. Highest-ROI rapport micro-skill: when a customer or teammate shares a win, celebrate it and ask them to say more.

---

## 5. Psychological Safety (Edmondson)

### Edmondson (1999) — Foundational Construct
"Psychological Safety and Learning Behavior in Work Teams," *Administrative Science Quarterly* 44(4), 350-383. (AOM Outstanding Publication, 2000.)

**Team psychological safety** = "a shared belief held by members of a team that the team is safe for interpersonal risk taking." Study of 51 teams: psychological safety → **learning behavior** (asking questions, seeking feedback, discussing errors, experimenting) → performance. Learning behavior **mediates** the safety→performance link.

### The 7-item scale (Edmondson 1999)
7-point Likert. Items (3 reverse-scored, marked R):
1. If you make a mistake on this team, it is often held against you. (R)
2. Members of this team are able to bring up problems and tough issues.
3. People on this team sometimes reject others for being different. (R)
4. It is safe to take a risk on this team.
5. It is difficult to ask other members of this team for help. (R)
6. No one on this team would deliberately act in a way that undermines my efforts.
7. Working with members of this team, my unique skills and talents are valued and utilized.

**Caveat (knowledge gap):** exact verbatim varies across secondary reproductions; this is the most commonly reproduced version. **Check against the 1999 ASQ original before any formal measurement use.**

### Psychological safety ≠ trust
Edmondson's distinction:
- **Level** — psych safety is a *group/team climate*; trust is *dyadic* (one party → a specific other).
- **Focus** — trust is about whether *the other* will act well; psych safety is about whether *I* can take a risk here ("is it safe to be me?").
- **Time frame** — psych safety concerns *near-term* interpersonal consequences (embarrassment/punishment now); trust extends to longer-term vulnerability.

A team can have high psych safety (safe to challenge ideas in a meeting) without deep personal trust between any two members, and vice versa. Build both; do not assume one delivers the other.

### Leader behaviors (Edmondson, *The Fearless Organization*, 2019)
Build safety by:
1. **Framing work as a learning problem** — acknowledge uncertainty + interdependence.
2. **Modeling fallibility/curiosity** — "I may miss things; I need to hear from you."
3. **Inviting participation** with genuine questions.
4. **Responding productively** to bad news — appreciation, not punishment; destigmatize failure.

Erode it by: blaming the messenger, punishing honest mistakes, dismissing input, projecting know-it-all certainty.

### Timothy R. Clark (2020) — Four Stages of Psychological Safety
*The 4 Stages of Psychological Safety*. A progressive ladder on two dimensions (respect + permission):
1. **Inclusion Safety** — safe to belong / be yourself.
2. **Learner Safety** — safe to learn, ask questions, make mistakes.
3. **Contributor Safety** — safe to contribute and make a difference using your skills.
4. **Challenger Safety** — safe to challenge the status quo, dissent, propose change. The highest stage; where innovation lives.

Diagnostic value: most teams stall **before Challenger safety**, which kills dissent and innovation. Use the stages to locate *where* safety breaks.

### Google Project Aristotle (2012-2015)
180+ teams, 200+ interviews, 250+ attributes. **Psychological safety was the #1 of five dynamics** of effective teams (others: dependability, structure & clarity, meaning, impact). How the team interacted mattered more than individual talent.

### The two misconceptions to flag
- Psychological safety is **NOT** niceness, comfort, or lowered standards. High-performing cultures **pair high safety with high accountability/standards** (Edmondson's 2x2: high safety + high standards = the learning zone; high safety + low standards = the comfort zone). Never sell safety as "be nice."
- Trust is **NOT** blindness. Mature trust (knowledge/identification-based) is **calibrated by evidence**, not granted unconditionally.

---

## Applied playbook — TAM customer relationships AND internal teams

- **Trust = (Credibility + Reliability + Intimacy) / Self-Orientation.** Fastest advisor lever = **lower self-orientation** (visibly put the customer's outcome ahead of your quota/agenda). Reliability compounds (keep every small promise); intimacy differentiates (make it safe to share the real problem).
- **Sequence on a new account/champion:** lead with **integrity + ability** signals; benevolence/affective trust comes with time. Treat the early relationship as **swift trust** — borrow credibility from role, references, standards — and protect it because it's fragile.
- **Match the repair to the breach:** blown SLA / wrong recommendation (competence) → **apologize + show the fix**; perceived broken commitment / honesty question (integrity) → if innocent, **explain/deny with evidence**, never go silent. Honor the asymmetry — one mishap can erase months; over-invest in prevention and rapid recovery.
- **Rapport is built, not faked:** genuine similarity, more positive contact frequency, let synchrony happen, run a positive-heavy interaction ratio. Highest-ROI micro-skill: **active-constructive responding** to a customer's or teammate's *good* news.
- **Internal teams: build psychological safety, not just trust.** Group-level, strongest predictor of team performance. Leaders: frame work as learning, model fallibility, invite real input, reward truth-telling. Use the 7-item scale to measure, Clark's 4 stages to diagnose where it breaks, and pair safety with accountability.

## Knowledge gaps (carry as caveats)
1. **Edmondson 7-item exact wording** — secondary reproductions vary; verify against the 1999 ASQ original for formal measurement.
2. **Losada/team positive:negative ratio** — the marital 5:1 is well replicated; the team threshold rests partly on Losada's math, since critiqued (Brown, Sokal & Friedman 2013). Direction holds; exact team number is heuristic.

## Sources (27)
Mayer/Davis/Schoorman 1995 (AMR 20(3), JSTOR 258792); Schoorman/Mayer/Davis 2007 retrospective; McAllister 1995 (AMJ 38(1)); Maister/Green/Galford 2000 *The Trusted Advisor* (Trusted Advisor Associates; ModelThinkers); Meyerson/Weick/Kramer 1996 swift trust (Wikipedia; Sage Knowledge); Kim/Ferrin/Cooper/Dirks 2004 (JAP 89(1); PubMed 14769123; SSRN 398221); reticence (Semantic Scholar); Slovic 1993 (Risk Analysis 13(6), Wiley); Lewicki & Bunker 1996 (NUS PDF); Tickle-Degnen & Rosenthal 1990 (Psychological Inquiry 1(4), T&F); Byrne 1997 review (SAGE); Zajonc 1968 mere-exposure (Wikipedia; SimplyPsychology); Chartrand & Bargh 1999 (JPSP 76(6)); Gottman 5:1 (Gottman Institute; Psychology Today); Gable & Reis 2004 (JPSP 87(2)); Edmondson 1999 (ASQ 44(4), SAGE; MIT PDF); 7-item scale (Atlantis Press; NovoPsych TPS-7); psych safety ≠ trust (Redefining Comms; ResearchGate); *The Fearless Organization* 2019 (Mindtools; Root Inc); Clark 4 stages (LeaderFactor); Project Aristotle (Google re:Work; Psych Safety).
