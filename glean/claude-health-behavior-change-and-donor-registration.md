# health-behavior-change-and-donor-registration

**Category:** Cloud, DevOps & Infrastructure
**Platform:** Claude
**Original Path:** claude/standalone/health-behavior-change-and-donor-registration

## Description
Apply health-behavior models and organ-donor registration intervention evidence. HBM (susceptibility, severity, benefits, barriers, cues, self-efficacy); TPB/TRA (attitude, norm, PBC, intention→behavior); EPPM (threat x efficacy, danger- vs fear-control); COM-B (capability, opportunity, motivation); the willingness-registration GAP; mandated/prompted/active choice, opt-in vs opt-out/presumed consent (behavioral-defaults evidence), reciprocity priming, loss-framing, social norms, UK BIT RCT 1M+. TRIGGER: HBM/TPB/EPPM for a health campaign; fear-appeal design; intention-action gap; donor registration campaign design; opt-in/opt-out/mandated/prompted-choice evidence (not policy advocacy); reciprocity or loss-framing for donation messaging; COM-B barrier mapping; social-norm messaging for registration; why supporters don't register. SKIP: SDT/Fogg/habit loops/stages-of-change → applied-psychology (references/behavior-change-psychology.md); OPO/OPTN/system/consent law → venture-organ-donation-system; presumed-consent policy advocacy → venture-organ-donation-frontier; Cialdini/norm theory/ELM/reactance → applied-psychology (references/persuasion-and-influence-psychology.md); nudge/defaults/choice-architecture → applied-psychology (references/behavioral-decision-making.md); EPPM for charitable giving/prosocial framing → psychology-of-charitable-giving.

---

<!-- SPOKE of the applied-psychology hub. Lives at applied-psychology/references/health-behavior-change-and-donor-registration/ when folded; currently installed as a standalone skill until folding completes. -->
# Health Behavior Change Models & Organ Donor Registration Interventions

> **Spoke of the `applied-psychology` hub.** For adoption/motivation models (SDT, Fogg,
> stages-of-change, habit loops, goal-setting) see
> `applied-psychology/references/behavior-change-psychology/SKILL.md`. For persuasion theory
> (Cialdini, ELM, reactance) see
> `applied-psychology/references/persuasion-and-influence-psychology/SKILL.md`. For EPPM
> applied to charitable giving or prosocial donation-appeal design, see
> `psychology-of-charitable-giving`.

---

## The core challenge: most people support donation. Very few register.

Surveys: 80–90% US/UK public support organ donation. 2023: ~58% Americans registered; pre-law-change UK ~40%. **Willingness–registration gap** = intention–behavior gap as public-health emergency. Need both formal behavior-change models and empirical intervention evidence.

---

## Decision table: pick the model for the diagnostic question

| Diagnostic question | Primary model | What it tells you |
| --- | --- | --- |
| Why does someone not feel personally at risk? | Health Belief Model | Low perceived susceptibility or severity; address cues to action |
| Why does someone who supports donation never register? | TPB intention-action gap | PBC barriers (don't know how, family conflict anxiety) dominate |
| Why did a fear-appeal campaign backfire? | EPPM | High threat + low efficacy → fear-control (denial, avoidance) |
| Why did the social norms message reduce registrations? | EPPM + social proof failure | Ambiguous norm cue; UK BIT trial found one variant backfired |
| What's the fastest single-message lift at a registration touchpoint? | BIT RCT evidence | Reciprocity frame ("If you needed one, would you take one?") |
| Should we switch to opt-out? | Defaults evidence (contested) | Opt-out lifts registration counts but causation on actual donation is debated |
| How to map barriers for a campaign intervention? | COM-B / Behaviour Change Wheel | Diagnose capability / opportunity / motivation deficits first |

**Response structure:** (1) Use decision table to select models. (2) Multi-model questions: lead with COM-B as diagnostic scaffold, layer HBM/TPB/EPPM for mechanism detail. (3) Always pair model recommendation with evidence caveat from Limitations section.

---

## 1. Health Belief Model (HBM)

**Provenance:** Developed 1950s–70s, Hochbaum, Rosenstock, colleagues at US Public Health Service. Self-efficacy added by Rosenstock, Strecher, & Becker (1988).

### Constructs

| Construct | Definition | Organ-donation example |
| --- | --- | --- |
| Perceived susceptibility | Belief one is at risk of a condition | "Someone in my family might need a transplant" |
| Perceived severity | Belief the condition is serious | "Without a transplant, patients die" |
| Perceived benefits | Belief the action reduces risk / produces benefit | "My organs could save up to eight lives" |
| Perceived barriers | Belief the action is costly, risky, or effortful | "I don't want doctors to give up on me", "my family will be upset" |
| Cues to action | Internal or external stimuli that trigger action | Seeing a "register now" prompt while renewing a license |
| Self-efficacy | Confidence one can perform the behavior | "I know how to register and it takes 2 minutes" |

**HBM for registration campaigns:**
1. Survey for dominant barrier (fear of family conflict, distrust of medical system, body-integrity concerns) — barriers predict non-registration most consistently.
2. Address barriers directly with brief factual corrections before registration ask.
3. Pair cue to action immediately after barrier addressed.

**Replication note:** HBM has strong face validity, decades of use. Meta-analyses show moderate predictive validity (barriers and benefits strongest; susceptibility and severity weaker). Predicts intention reasonably well, less reliable for actual behavioral follow-through without PBC and implementation intentions. Effect sizes vary substantially across behaviors and populations.

---

## 2. Theory of Reasoned Action (TRA) & Theory of Planned Behavior (TPB)

**Provenance:** TRA — Fishbein & Ajzen (1975); TPB — Ajzen (1985, 1991), adding PBC to TRA.

### Constructs

**TRA (two-component intention model):**
- **Attitude toward the behavior**: Personal evaluation of donation (good/bad, beneficial/harmful)
- **Subjective norm**: Perceived social pressure ("My family and friends would want me to register")
- → Combined attitudes + subjective norms → **Behavioral intention** → **Behavior**

**TPB extends TRA by adding:**
- **Perceived behavioral control (PBC)**: Belief in ability to perform the behavior ("I know how and I'm capable of registering"; functionally similar to self-efficacy)
- PBC predicts both intention AND behavior directly (especially when actual control is imperfect)

### Application to organ donation

2021 systematic review of 17 TPB-based intervention studies (Asgarimojarad et al., *Saudi J Kidney Dis Transpl*): TPB-based interventions improve organ donation behavior. **PBC strongest moderating predictor across cultures**: structural barriers (knowing how to register, family discussion anxiety, mistrust) dominate when attitudes already positive. Removing knowledge barriers + providing immediate registration opportunity addresses PBC directly — most consistent finding across community-based intervention studies.

**Cultural note:** Cross-cultural TPB studies: attitude → intention stronger in individualist cultures (Americans); subjective norm → intention stronger in collectivist cultures (Koreans). Adjust intervention design accordingly.

**On the intention–action gap:** TPB explains *intention* well (R² often 40–60% in donation studies). Explains *behavior* less well. Gap between strong intention and actual registration driven primarily by PBC (structural barriers) and absence of immediate registration opportunity — not weak attitudes or low motivation. Single most consistent finding across every systematic review: **provide an immediate registration opportunity at the point of contact**.

---

## 3. Extended Parallel Process Model (EPPM)

**Provenance:** Witte (1992), *Communication Monographs*. Integrates Leventhal's danger-control / fear-control framework with Protection Motivation Theory (Rogers) and prior fear-appeal work.

### Core mechanism

Fear appeal received → recipient assesses two things in sequence:
1. **Threat appraisal**: How susceptible am I? How severe is it?
2. **Efficacy appraisal**: How effective is the recommended action? Can I do it?

Perceived efficacy/threat ratio determines response:

| Threat | Efficacy | Response | Outcome |
| --- | --- | --- | --- |
| Low | Any | No processing | Message ignored |
| High | High | **Danger control** | Protective behavior (desired outcome) |
| High | Low | **Fear control** | Denial, reactance, avoidance (backfire) |

**Why fear-only donation messages backfire:** "Three people die every day on the waiting list" = high-threat message. If recipient simultaneously feels registration is complicated, or doctors will give up on them if registered (low response efficacy), fear-control kicks in — avoidance instead of registration.

**The EPPM fix:** Pair high-threat messaging with high-efficacy information (registration takes 2 minutes; doctors prioritize treatment regardless of donor status; conditions can be specified). BIT RCT found loss-frame "Three people die every day" performed well *because* simultaneous context (simple registration page) provided high-efficacy cue.

**Fear appeal meta-analysis honest note:** Witte & Allen (2000): fear appeals generally effective when efficacy is high, ineffective or counterproductive when efficacy is low. Effect sizes moderate (r ≈ .20) and heterogeneous. EPPM well-supported in health communication literature but most tests use self-report outcomes; behavioral follow-through with objectively measured registration is rarer.

> **Scope note:** This section covers EPPM applied to registration behavior. For EPPM applied to charitable giving or prosocial donation-appeal design, see `psychology-of-charitable-giving`.

---

## 4. COM-B / Behaviour Change Wheel

**Provenance:** Michie, van Stralen, & West (2011), *Implementation Science*. Synthesizes 19 behavior-change frameworks into unified system.

### COM-B model

Behavior (B) requires three conditions — all must be sufficient:
- **Capability (C)**: Physical and psychological ability (knowledge of how to register; knowing where to go; understanding what donation means)
- **Opportunity (O)**: Physical and social environment enabling behavior (registration prompt at DMV or ID renewal; family acceptance; time available)
- **Motivation (M)**: Reflective (conscious decision, goals) and automatic (habits, emotions, impulses) processes

### Organ-donation diagnostic map

| COM-B deficit | What it looks like | Intervention type |
| --- | --- | --- |
| Capability — knowledge | "I didn't know I could register at the DMV" | Education/information |
| Capability — skills | "I don't know how to have the family conversation" | Training/rehearsal |
| Opportunity — physical | No registration prompt in the checkout flow | Environmental restructuring (prompted choice) |
| Opportunity — social | Family disapproval is expected | Social norm messaging; community peer-leader programs |
| Motivation — reflective | "I mean to, I just haven't done it yet" | Implementation intentions; cues to action |
| Motivation — automatic | Mortality salience causes avoidance | Fear-managed messaging (EPPM); framing |

**COM-B vs TPB:** COM-B maps intervention types to deficits; TPB identifies specific attitudinal and normative mechanisms driving intention. Complementary, not redundant.

**Behaviour Change Wheel — nine intervention functions** (Michie et al. 2011):

| Function | Definition | Organ-donation example |
| --- | --- | --- |
| Education | Increase knowledge/understanding | "3 people die daily" statistic |
| Persuasion | Use communication to induce change | Reciprocity frame; loss frame |
| Incentivisation | Rewards/punishments for behavior | Donor-recipient priority queue (Israel, Singapore) |
| Coercion | Penalties for non-behavior | (Not used in registration contexts) |
| Training | Skills practice | Family-conversation rehearsal guides |
| Restriction | Environmental limits | N/A for donation |
| Environmental restructuring | Change physical/social context | Opt-out default; prompted-choice at DMV |
| Modelling | Provide exemplar to observe | Celebrity/peer registrant testimonials |
| Enablement | Remove barriers; support capability | Online registration simplification |

---

## 5. The Willingness–Registration Gap

Gap empirically documented and theoretically important. Key features:

- **Scale**: US: ~85% approve of donation (Gallup); ~58% registered (2023 HRSA data). UK pre-2020: 80%+ support, ~40% registered.
- **Mechanism:** Not primarily attitudinal — most non-registrants have positive attitudes. Dominant drivers:
  1. **PBC barriers**: Perceived complexity, uncertainty about family impact, distrust of medical system, body-integrity concerns
  2. **Mortality salience**: Registration confronts personal death; avoidance is automatic motivation response
  3. **Absence of immediate registration opportunity**: Intentions without proximate action opportunity decay rapidly (implementation-intentions research)
  4. **Social-norm uncertainty**: People underestimate how many peers are registered (pluralistic ignorance reducing registration)

- **Bridging it:**
  - Single most replicated finding: **provide immediate, frictionless registration opportunity at point of contact** (Golding et al. 2017 systematic review; Cochrane 2021: RR 1.30 for registration-behavior outcomes)
  - Brief myth-busting reduces attitudinal ambivalence, increases PBC and subjective norm (Feeley narrative-with-information studies)
  - Anticipated regret ("How would you feel if a family member needed an organ and you hadn't registered?") increases initial commitment (Lowe et al. 2024, *PLOS ONE*, US national survey)

---

## 6. Registration Intervention Evidence Base

> **Scope:** Behavioral evidence for registration-choice architecture (defaults, framing, prompted choice). Policy advocacy for/against presumed-consent legislation out of scope; see `venture-organ-donation-frontier` for policy layer.

### Active choice / Mandated choice / Prompted choice

| Frame | Definition | Evidence |
| --- | --- | --- |
| Opt-in (US status quo) | Must actively consent to be registered | Baseline; lowest registration rates among the three |
| Opt-out / Presumed consent | Registered by default unless explicit objection | Cross-country data: ~25–30% higher donation rates (Abadie & Gay 2006) — but contested (see below) |
| Mandated choice | Must make a yes/no decision; no default | Odds of registration ~2× opt-in (Kessler & Roth 2014, *AEJ: Economic Policy*); Virginia's experience shows ~24% remain undecided |
| Prompted choice | Embedded in routine transaction (ID renewal, DMV) | Italy quasi-experiment (Fieles-Ahmad & Schulze Spuentrup 2026): significant increase in consent registrations; consistent with BIT findings |

**Opt-out / presumed consent causal evidence:**

Johnson & Goldstein (2003, *Science*; 2004, *Transplantation*): dramatic cross-country differences — opt-out countries ~60 percentage points higher consent rates. Online experiment (n=161): opt-out roughly doubles self-reported donation intent vs opt-in.

Cross-country comparison conflates consent law with infrastructure, health system capacity, GDP, traffic mortality rates, Catholic religious culture, and whether families routinely consulted before procurement. Abadie & Gay (2006) used panel regression — found ~25–30% higher donation rates in presumed-consent countries — but critics note identifying variation still observational. Shepherd et al. (2014): presumed consent only increases donation rates when family consent also routinely sought and combined registry maintained; otherwise effect negligible.

**Practical takeaway:** Opt-out/presumed-consent associated with higher registration counts and somewhat higher donation rates, but causal mechanism not purely defaults — also reflects accompanying infrastructure investments and cultural conditions. US is opt-in and unlikely to change near-term; actionable intervention evidence is in prompted-choice and message framing at point of contact.

### UK Behavioural Insights Team RCT (2013)

Largest RCT on organ donor registration: **n = 1,085,322** UK road-tax/driving-licence applicants on GOV.UK, randomized to one of eight message variants or control. Six of eight tested variants shown below; remaining two didn't differ significantly from control.

**Results (selected variants):**

| Message variant | Odds ratio vs control | Notes |
| --- | --- | --- |
| Reciprocity: "If you needed an organ transplant, would you have one? If so, please help others." | **OR 1.38** (p < 0.001) | Single best performer |
| Loss frame: "Three people die every day because there aren't enough organ donors." | OR 1.33 (p < 0.001) | Not significantly different from reciprocity |
| Gain frame: "You could save or transform up to 9 lives as an organ donor." | OR 1.25 (p < 0.001) | Effective |
| Cognitive dissonance: "If you support donation in principle, show it." | OR 1.23 (p < 0.001) | Effective |
| Social norms alone: "Every day thousands of people who see this page decide to register." | **OR 0.94** (p < 0.05) | **Significantly reduced registrations** |
| Social norms + image of people | OR ~0.96 (not significantly different from social-norms-alone) | Also backfired |

**Why did the social-norms message backfire?** Message described norm as *thousands* registering daily — may have highlighted vastness of decision or paradoxically communicated "most people still need to decide" rather than "most people have done it." Lacked any efficacy component. Sharp warning: social-norm messaging doesn't always increase prosocial behavior in donation contexts.

**BIT trial implications:**
1. Reciprocity priming = single most evidence-supported message type at registration scale
2. Loss-framing comparably effective (not significantly better)
3. Social-norm messages need precise calibration — generic norm message can backfire
4. Projected impact of best message: ~96,000 additional registrations per year if applied consistently across GOV.UK

### Community-based interventions

2021 Cochrane systematic review (46 RCTs/cluster-RCTs/quasi-RCTs; n up to 1.29M):
- Registration behavior: **RR 1.30 (95% CI 1.19–1.43)** — low certainty, high heterogeneity (I² = 84%)
- Community interventions targeting specific ethnic groups: **RR 2.14** — but I² = 85%, high heterogeneity
- Classroom interventions with transplant community speakers: RR 1.33, I² = 0% — most consistent
- Brief interventions addressing myths/concerns + immediate registration opportunity: most consistent finding across all reviews

**Honest caveat:** Most community-intervention studies have high risk of bias and high heterogeneity. Cochrane reviewers rate overall certainty as *low*. Treat effect sizes as upper-bound estimates.

---

## 7. Barriers-to-Registration Quick Reference

Most common empirically documented barriers:

| Barrier | Prevalence | Intervention |
| --- | --- | --- |
| Doctors won't try as hard to save me | Very common (esp. low-trust populations) | Factual correction (organs retrieved only after death is confirmed) |
| Body-integrity / disfigurement concern | Common | Factual: body treated with respect; open-casket funerals proceed normally |
| Religious objection | Moderate; varies by faith | Faith-concordant messaging; most major religions support donation |
| Family will be upset | Moderate | Facilitate family conversation; provide conversation scripts |
| Don't know how to register | Moderate | Immediate registration opportunity; reduce friction |
| Haven't thought about it | Very common | Salience + cue to action at point of contact |
| "I'm too young / healthy to think about this" | Common in 18–30 | Reframe: registration is for others' benefit, not risk acknowledgment |

---

## Limitations and Replication Honesty

- **HBM:** Moderate predictive validity; barriers/benefits outperform susceptibility/severity. Doesn't specify how to combine constructs; less predictive of behavior than intention.
- **TPB:** Intention prediction strong; behavior prediction weaker and intention–behavior correlation decays over time. PBC conflates self-efficacy with actual controllability.
- **EPPM:** Well-supported in communication research; most tests use self-report; rare behavioral registration outcomes. Fear-control vs danger-control distinction theoretically clean but empirically harder to measure.
- **COM-B / BCW:** Widely used in health policy; framework taxonomic rather than predictive; doesn't specify which intervention functions produce largest effects.
- **Presumed-consent evidence:** Cross-country comparisons observational. Causal claims about opt-out *causing* higher donation rates contested.
- **BIT RCT (2013):** Strong internal validity (large n, true randomization). External validity: UK context, single online touchpoint; reciprocity effect may not generalize identically across cultures or settings.

---

## Sources (19 cited)

1. Rosenstock, IM, Strecher, VJ, & Becker, MH (1988). Social learning theory and the Health Belief Model. *Health Education Quarterly*, 15(2), 175–183.
2. Fishbein, M, & Ajzen, I (1975). *Belief, Attitude, Intention, and Behavior.* Addison-Wesley.
3. Ajzen, I (1991). The theory of planned behavior. *Organizational Behavior and Human Decision Processes*, 50(2), 179–211.
4. Witte, K (1992). Putting the fear back into fear appeals: The Extended Parallel Process Model. *Communication Monographs*, 59(4), 329–349.
5. Witte, K, & Allen, M (2000). A meta-analysis of fear appeals: Implications for effective public health campaigns. *Health Education & Behavior*, 27(5), 591–615.
6. Michie, S, van Stralen, MM, & West, R (2011). The behaviour change wheel: A new method for characterising and designing behaviour change interventions. *Implementation Science*, 6, 42.
7. Johnson, EJ, & Goldstein, DG (2003). Do defaults save lives? *Science*, 302(5649), 1338–1339.
8. Johnson, EJ, & Goldstein, DG (2004). Defaults and donation decisions. *Transplantation*, 78(12), 1713–1716.
9. Abadie, A, & Gay, S (2006). The impact of presumed consent legislation on cadaveric organ donation: A cross-country study. *Journal of Health Economics*, 25(4), 599–620.
10. Shepherd, L, O'Carroll, RE, & Ferguson, E (2014). The impact of presumed consent laws and institutions on deceased organ donation. *European Journal of Health Economics*, 15(8), 853–870.
11. Behavioural Insights Team / NHS Blood and Transplant (2013). *Applying Behavioural Insights to Organ Donation: Preliminary Results from a Randomised Controlled Trial.* Cabinet Office, UK.
12. Gill, C, et al. (2018). Effect of persuasive messages on NHS organ donor registrations: A pragmatic, quasi-randomised trial. *Trials*, 19(1), 514.
13. Golding, BF, Gilles, K, Ahvenainen, S, et al. (2017). A systematic narrative review of effects of community-based intervention on rates of organ donor registration. *Progress in Transplantation*, 27(3), 273–281.
14. Morgan, SE, et al. (2021). Interventions for increasing solid organ donor registration. *Cochrane Database of Systematic Reviews*, Issue 4. CD010829.
15. Hyde, MK, & White, KM (2009). To be a donor or not to be? Applying an extended TPB to predict posthumous organ donation intentions. *Journal of Applied Social Psychology*, 39(7), 1654–1674.
16. Feeley, TH, & Anker, AE, et al. (2020). Systematic literature review of interventions for promoting postmortem organ donation from a social marketing perspective. *Progress in Transplantation*, 30(2), 145–158.
17. Asgarimojarad, A, et al. (2021). Factors affecting organ donation registration: A systematic review using the Health Belief Model. *Health Psychology Open*, 8(2).
18. Lowe, SR, et al. (2024). Anticipated regret and organ donor registration: A national survey study. *PLOS ONE*, 19(3).
19. Kessler, JB, & Roth, AE (2014). Organ allocation policy and the decision to donate. *American Economic Journal: Economic Policy*, 6(4), 261–293.

---

## Cross-hub routing summary

| Topic | Route |
| --- | --- |
| SDT, Fogg B=MAP, habit loops, stages-of-change, goal-setting | `applied-psychology/references/behavior-change-psychology/` |
| Donation SYSTEM: OPO operations, UNOS allocation, consent law, NC registry | `venture-organ-donation-system` |
| Opt-out / presumed consent POLICY advocacy and tracking | `venture-organ-donation-frontier` |
| Social proof / descriptive + injunctive norms theory (Cialdini) | `applied-psychology/references/persuasion-and-influence-psychology/` |
| Nudge, status-quo bias, defaults as general decision framing | `applied-psychology/references/behavioral-decision-making/` |
| Fear, emotional regulation, mortality salience as affect construct | `applied-psychology/references/emotion-and-affect-psychology/` |
| EPPM for charitable giving / prosocial donation appeal design | `psychology-of-charitable-giving` |