---
name: customer-satisfaction-loyalty-science
description: >-
  The SCIENCE of measuring and moving B2B customer satisfaction, loyalty, and advocacy — the constructs
  and theory beneath any single CSAT/NPS number, for a TAM. Replication-honest about contested findings.
  TRIGGER: what CSAT/NPS/CES each measures + their critiques; the NPS "one-number" debate / does NPS
  predict growth; Customer Effort Score; expectancy-disconfirmation (Oliver); the Kano model (satisfiers
  vs delighters); SERVQUAL / RATER / Gaps Model; the satisfaction→loyalty→revenue chain & Service-Profit
  Chain; "satisfied customers defect" / apostle model; satisfaction vs loyalty vs advocacy; why B2B
  satisfaction differs (multi-stakeholder, false loyalty); which metric to trust.
  SKIP: churn/CLV/survival analytics → da-23-customer-lifetime-value; survey-item or feedback-reply
  WRITING → career-and-formal-writing / content-and-marketing-writing; outcome-based CS / proving value
  → value-realization-outcome-cs; pre-purchase buying-committee psychology → enterprise-b2b-buyer-psychology.
version: 1.0.0
updated: 2026-06-18
category: applied-psychology
whenToUse: >-
  Load when the question is about the THEORY and MEASUREMENT of customer satisfaction, loyalty, or
  advocacy — what a metric (CSAT, NPS, CES) actually measures and its documented critiques, the
  satisfaction/quality models (expectancy-disconfirmation, Kano, SERVQUAL/Gaps), the
  satisfaction→loyalty→revenue chain (Service-Profit Chain, Jones & Sasser apostle model), the
  satisfaction-vs-loyalty-vs-advocacy distinction (Oliver's four loyalty phases), or why B2B
  account satisfaction behaves differently from B2C. Not for survey-writing, churn analytics, or
  value-realization mechanics (cross-referenced).
keywords:
  - customer satisfaction
  - customer loyalty
  - NPS
  - Net Promoter Score
  - CSAT
  - Customer Effort Score
  - CES
  - expectancy disconfirmation
  - Kano model
  - SERVQUAL
  - RATER
  - gaps model
  - service-profit chain
  - apostle model
  - zone of indifference
  - satisfaction loyalty advocacy
  - Oliver loyalty phases
  - Keiningham one number
  - Effortless Experience
  - B2B customer satisfaction
  - voice of customer programs
  - customer advocacy
tags:
  - applied-psychology
  - customer-satisfaction
  - customer-loyalty
  - measurement
  - tam
  - cx
  - voice-of-customer
  - service-quality
---

# Customer Satisfaction & Loyalty Science

The measurement-and-theory layer beneath every customer feedback number. This is **not** a survey-writing
guide and **not** a churn-analytics toolkit — it is the science of *what satisfaction and loyalty are,
how the field measures them, what each metric actually captures, and which findings are solid vs.
contested*. Audience: a B2B SaaS / MongoDB TAM who has to read an account's satisfaction signals
honestly and act to move them.

**Replication honesty is the point.** Much of this domain is vendor-marketed as settled when the
peer-reviewed literature is divided. The single most important meta-finding: **no single metric
("one number") reliably predicts growth or retention across contexts** — the best predictor is
industry-dependent, and combining metrics beats any one of them.[^dehaan][^keiningham] State that
honesty whenever you cite a metric.

## How to use this skill

1. Pick the question type and read the matching reference (table below).
2. When you cite a finding, carry its **confidence tag**: **well-replicated** (state as fact),
   **qualified** (2 sources / one authoritative study), or **contested/vendor-origin** (preserve the
   disagreement — never launder a vendor claim into a fact).
3. For a TAM, always end at the **B2B reading** (§7 / `b2b-and-application.md`): an account is a
   *buying center*, not one respondent, and a single account-level score hides stakeholder variance
   and false loyalty.

## Reference routing table

| Question | Read |
| --- | --- |
| What does CSAT / NPS / CES each measure, how is it computed, and what are its documented critiques? The NPS "one-number" debate. Which metric to trust. | `references/metrics-and-critiques.md` |
| Satisfaction theory: expectancy-disconfirmation (Oliver), the Kano model (must-be / performance / delighter), prioritizing satisfiers vs. delighters | `references/satisfaction-theory.md` |
| Service-quality measurement: SERVQUAL, the five RATER dimensions, the Gaps Model, SERVQUAL-vs-SERVPERF critique | `references/service-quality-servqual.md` |
| The satisfaction→loyalty→retention→revenue chain; Service-Profit Chain (Heskett); "satisfied customers defect" / zone of indifference / apostle model (Jones & Sasser); causal asymmetry & diminishing returns | `references/loyalty-profit-chain.md` |
| Satisfaction vs loyalty vs advocacy (Oliver's four loyalty phases; attitudinal vs behavioral; false loyalty); B2B vs B2C differences; how a TAM reads an account score | `references/b2b-and-application.md` |

The five references hold the full detail with inline citations. The rest of this file is the
**executive synthesis** — enough to answer most questions directly, with pointers to the depth.

---

## 1. The metrics, in one screen

Three metrics dominate. They measure **different things on different time axes** and are complements,
not substitutes.[^silabs][^qualtrics]

| Metric | Question | Time axis | Measures | Best at | Blind to |
| --- | --- | --- | --- | --- | --- |
| **CSAT** (Customer Satisfaction Score) | "How satisfied were you with [X]?" (1–5 typical; report **top-box %**) | Past, **transactional** (a touchpoint) | Attribute/episode satisfaction | Pinpointing where one interaction fell short | Relationship, future behavior, "why" |
| **NPS** (Net Promoter Score) | "How likely to recommend?" (0–10; **%promoters − %detractors**, −100…+100) | Relationship, **intention/future** | Stated recommend-intention (advocacy proxy) | Relationship pulse, closed-loop trigger | Effort, the "why", actual behavior gap |
| **CES** (Customer Effort Score) | "[Company] made it easy to handle my issue" (1–7 agree; **top-3-box %**) | Past, **process-specific** | Effort/ease of a service interaction | Predicting *dis*loyalty in support contexts | Satisfaction, delight, relationship, "why" |

Key construct facts (well-replicated):
- **Promoters 9–10, Passives 7–8, Detractors 0–6** for NPS; passives are counted in the base but not added.[^reichheld2003]
- **CES evolved** from a 1–5 inverted "how much effort" item (CES 1.0, 2010) to a 1–7 agreement
  "made it easy" item (CES 2.0, ~2013) because the inverted scale produced false positives and
  "effort" translated poorly; "easy" is more universal.[^greatbrook][^etouchpoint]
- All three share survey pathologies: **self-selection / non-response bias, survival bias** (churned
  customers never respond), **recency bias**, and **gaming** when tied to incentives.[^supportlogic]

Full detail, formulas, lineage (ACSI/Fornell), and citations → `references/metrics-and-critiques.md`.

## 2. The "one-number" debate (genuinely contested — do not pick a side as fact)

Reichheld's *"The One Number You Need to Grow"* (HBR 2003) claimed "would you recommend" was the
single best predictor of growth.[^reichheld2003] The peer-reviewed literature does **not** support
the superiority claim:

- **Keiningham, Cooil, Andreassen & Aksoy (2007, *Journal of Marketing*)** replicated Reichheld's own
  method on 21 firms / 15,500+ interviews (Norwegian barometer) and concluded **"Net Promoter is not
  superior to the ACSI"** — comparable, statistically indistinguishable predictive power.[^keiningham]
- **Kristensen & Eskildsen (2011)** and **Fisher & Kordupleski (2019)** independently found NPS a poor
  / no-better predictor, distribution-sensitive, and uncalibrated by competitive info.[^kristensen][^fisher]
- **de Haan, Verhoef & Wiesel (2015, *IJRM*)** compared CSAT, NPS, CES across 93 firms / 18 industries
  controlling for self-selection: **top-2-box CSAT predicted retention best overall, the best metric
  varies by industry, and combining metrics beats any single one.**[^dehaan]

**The honest, replicated summary:** NPS and CSAT have *roughly comparable* predictive ability; neither
is reliably superior; the NPS→growth *causation* is unestablished (could run growth→NPS).[^measuringu]

**The other side (preserve it):** Bain/Reichheld argue critics attack the wrong target — *absolute*
bottom-up scores — whereas growth prediction relies on **relative / competitive-benchmark NPS** (you
vs. direct competitors). Bain claims relative NPS explains 10–70% of revenue-growth variation and that
the NPS leader grows ~2× faster — but this is **Bain proprietary research, not independently
replicated**; treat it as a contested business claim.[^bain]

**TAM takeaway:** NPS is useful as a relationship pulse and closed-loop trigger. Never present it (or
any single number) as *the* predictor of account growth. The CES "better than CSAT/NPS" claim is
similarly **CEB/Gartner's own research, largely un-replicated, and scoped to service/support contexts**
— and the one independent cross-metric test (de Haan) actually ranks top-2-box CSAT ahead of CES for
retention.[^dixon2010][^dehaan]

## 3. Why satisfaction happens — the theory

**Expectancy-Disconfirmation Theory (EDT) — Oliver (1980).** The dominant satisfaction paradigm.
Satisfaction is a function of **prior expectations** and **disconfirmation** (the gap between expected
and perceived performance): positive disconfirmation (performance > expectations) → satisfaction;
negative → dissatisfaction; confirmation → neutral/met.[^oliver1980] This is *why* a merely-"met
expectations" experience scores neutral, and why over-promising (raising expectations) manufactures
dissatisfaction even when delivery is unchanged. Critiques: "expectation" is hard to measure and
polysemic; perceived performance alone sometimes drives satisfaction. Detail → `references/satisfaction-theory.md`.

**The Kano Model — Kano (1984).** Maps product/service attributes to satisfaction response, on two
axes (functionality vs. satisfaction). Five categories:
- **Must-be** (basic / dissatisfiers): absence enrages, presence is merely expected — asymptotic on the
  downside. *Table stakes; you can't win on them, only lose.*
- **One-dimensional** (performance / satisfiers): satisfaction rises ~linearly with performance. *Where
  you compete head-to-head.*
- **Attractive** (delighters / exciters): presence delights, absence isn't penalized. *Differentiation.*
- **Indifferent**: customer doesn't care. *Stop investing.*
- **Reverse**: more of it makes some customers *less* satisfied. *Watch for segments.*

Used to **prioritize satisfiers vs. delighters**, and it has a *dynamic*: today's delighters decay into
tomorrow's expected must-be's (the "natural decay of delight").[^kano1984] This is the antidote to the
naïve "delight customers" reflex — must-be failures destroy satisfaction faster than delighters build
it (cf. the asymmetry in §5). Detail → `references/satisfaction-theory.md`.

## 4. Measuring service quality — SERVQUAL / RATER / the Gaps Model

**SERVQUAL — Parasuraman, Zeithaml & Berry (1985 model; 1988 instrument).** Service quality = the gap
between customer **expectations (E)** and **perceptions (P)** of performance; a 22-item instrument
scoring **Q = P − E**.[^pzb1988] The original 10 determinants collapsed (1988) to **five RATER
dimensions**:

| | Dimension | Definition |
| --- | --- | --- |
| **R** | Reliability | Perform the promised service dependably and accurately. *PZB's repeated finding: most important to customers.* |
| **A** | Assurance | Knowledge/courtesy of employees; ability to inspire trust and confidence. |
| **T** | Tangibles | Physical (and digital) evidence — facilities, equipment, materials. *Consistently least important.* |
| **E** | Empathy | Caring, individualized attention. |
| **R** | Responsiveness | Willingness to help and provide prompt service. |

**The Gaps Model (5 gaps)** explains *why* quality falls short, with **Gap 5 (the customer gap,
expected vs. perceived) = f(Gaps 1–4)**:
1. **Gap 1 — Listening/Knowledge**: what customers expect vs. what management *thinks* they expect.
2. **Gap 2 — Standards/Design**: management understanding not translated into clear service standards.
3. **Gap 3 — Delivery/Performance**: failing to deliver to your own standards.
4. **Gap 4 — Communication**: external promises (marketing/sales) ≠ what's delivered (over-promising).
5. **Gap 5 — Customer**: the outcome gap; close Gaps 1–4 and it shrinks.[^pzb1985]

**Contested:** SERVQUAL's **difference (P−E) scores** are psychometrically weak; **Cronin & Taylor
(1992) SERVPERF** (performance-only, no expectations battery) is more parsimonious and often
psychometrically superior. The **Carrillat et al. (2007) meta-analysis** found the two **about equally
valid** at predicting overall quality — but SERVQUAL retains better *diagnostic* power (it tells you
*which* gap to close).[^cronin1992][^carrillat] Detail → `references/service-quality-servqual.md`.

## 5. The satisfaction → loyalty → revenue chain (and why it's not linear)

**Service-Profit Chain — Heskett et al. (1994).** A causal chain: *internal service quality → employee
satisfaction → employee retention/productivity → external service value → customer satisfaction →
customer loyalty → revenue growth & profitability*, with the **"satisfaction mirror"** (employee and
customer satisfaction reflect each other).[^spc1994] The loyalty→profit tail rests on Reichheld &
Sasser's **"Zero Defections" (1990): a 5-point cut in defection can raise profits ~25–85%.**[^zerodefections]
A meta-analysis (Hogreve et al. 2017) found all links significant — **but** a contrarian stream
(Silvestro & Cross 2000; Silvestro 2002) failed to support the *satisfaction-mirror* link, even finding
employee *dis*satisfaction correlated with profit in some retail settings, and most evidence is
**correlational** so the causal ordering isn't firmly established.[^hogreve][^silvestro]

**The satisfaction-loyalty curve is NON-LINEAR — Jones & Sasser (1995, "Why Satisfied Customers
Defect").** Xerox found **"totally satisfied" customers were 6× more likely to repurchase than merely
"satisfied" ones.** The curve's *steepness depends on competitive intensity*: where switching is easy,
any drop from total satisfaction collapses loyalty; in low-competition/lock-in markets the curve is
flat.[^jonessasser] The **Apostle Model** (satisfaction × loyalty 2×2):

| | Will leave | Stays loyal |
| --- | --- | --- |
| **Satisfied** | **Mercenaries** (satisfied, not loyal; chase price) | **Loyalists / Apostles** (satisfied + loyal + advocate) |
| **Dissatisfied** | **Defectors / Terrorists** (spread negative word-of-mouth) | **Hostages** (dissatisfied but stuck — lock-in / high switching costs) |

**Causal asymmetry & diminishing returns — Anderson & Mittal (2000).** The chain links are
**asymmetric** (negative attribute performance hurts satisfaction *more* than equal positive
performance helps — negativity bias) and **non-linear** (diminishing returns on intentions at high
satisfaction).[^andersonmittal] This is the academic grounding for: **lifting a detractor out of
dissatisfaction matters more than delighting an already-happy account**, and **you cannot average
across the curve** — an average CSAT/NPS misleads. Detail → `references/loyalty-profit-chain.md`.

## 6. Satisfaction vs loyalty vs advocacy — three different things

- **Satisfaction** = attitudinal, post-consumption *evaluation* (did it meet the standard).
- **Loyalty** (Oliver 1999) = *"a deeply held commitment to rebuy… despite situational influences and
  marketing efforts having the potential to cause switching."* Has **attitudinal** (commitment) and
  **behavioral** (repeat purchase) dimensions. Repeat buying from **inertia, habit, or lock-in is NOT
  true loyalty** ("spurious"/"false" loyalty; the "hostage").[^oliver1999][^dickbasu]
- **Advocacy** = active recommendation/referral (positive word-of-mouth) — a *behavioral* step beyond
  private loyalty; what NPS's "would you recommend" tries to capture.

**Oliver's four phases of loyalty** — sequential, each with its own vulnerability: **cognitive**
(loyalty to information/price — weakest, "phantom") → **affective** (liking from satisfying
experiences) → **conative** (commitment/intention to rebuy) → **action** (intention + inertia +
willingness to overcome obstacles — deepest).[^oliver1999] **Satisfaction is necessary but becomes a
weaker driver of loyalty as commitment forms through fortitude and social bonding** — and Reichheld's
"satisfaction trap" reports **65–85% of "satisfied" customers still defect.**[^oliver1999][^jonessasser]

## 7. The B2B reading — what changes for a TAM

This is where a TAM must land. (Detail → `references/b2b-and-application.md`.)

1. **An account is a buying center, not one respondent.** Multiple roles (economic sponsor, technical
   buyer, end users, procurement, finance, champion) evaluate you on *different criteria* and reach
   *different satisfaction levels*; B2B buying committees run ~7–12 stakeholders (a practitioner figure
   [Q]).[^sis] Academically [F], **Austen, Herbst & Bertels (2012, *Industrial Marketing
   Management*)** show account-level satisfaction is a **non-linear** function of individual members'
   satisfaction and that **negative individual judgments carry more weight** — a single dissatisfied
   stakeholder can sink the account, and a simple average hides it.[^austen]
2. **A single account-level CSAT/NPS hides stakeholder variance.** High executive satisfaction can mask
   daily user frustration (or vice versa). Read **role-divergence within an account as a risk signal in
   itself**, and weight a detractor *decision-maker* far above a detractor *end-user*. The highest churn
   risk often sits in the roles least likely to be surveyed.[^sis]
3. **Dissatisfaction ≠ near-term churn (the dangerous lag).** High switching costs and multi-year
   contracts mean a customer can be a genuine detractor with no intent to leave — Jones & Sasser's
   **"hostage"** / false loyalty in B2B form. Lock-in retention is not loyalty; it can flip to
   active churn (and "terrorist" word-of-mouth) the moment a credible alternative appears.[^jonessasser]
4. **Oliver explicitly scoped his loyalty model to B2C**, noting B2B has "additional variables, such as
   power dependencies."[^oliver1999] Apply the four-phase model to accounts with care: the account team
   itself is a loyalty mechanism (relationship/social bonding) that can create attitudinal loyalty
   independent of product satisfaction — and can mask product dissatisfaction.

**Net operating stance:** measure multi-stakeholder, segment by role, watch asymmetry (fix detractors
before delighting promoters), distrust the single number, and separate true loyalty from lock-in.

## Output guidance — what a good answer looks like for a TAM

- **Citing a metric honestly:** name what it measures, its time axis, and its confidence tag — e.g.
  *"Their relationship NPS is 40; that's a useful pulse and a closed-loop trigger, but NPS does not
  reliably predict growth (Keiningham 2007), so I wouldn't read it as a renewal forecast."* Never present
  one number as the verdict.
- **Reading a split account score (worked example):** account CSAT looks healthy at 4.4/5, but the
  economic sponsor sits at 2/5 while end users are at 4.8. Don't average it. Flag the role-divergence as a
  risk, weight the dissatisfied decision-maker (Austen et al. 2012 — negatives weigh more), and ask
  whether the renewal is true loyalty or a *hostage* held by switching costs (Jones & Sasser). The
  decision-maker's 2/5 is the headline, not the 4.4 average.
- **Prioritizing an action:** fix a must-be / high-effort failure (Kano + CES) before adding a delighter;
  negativity asymmetry means the basic failure is costing more satisfaction than any delighter would add.

---

## Anti-patterns

- **Presenting one number as "the" predictor of growth/retention.** The replicated evidence says no
  single metric is reliably superior; the best one is context-dependent; combine them.[^dehaan][^keiningham]
- **Laundering a vendor claim into a fact.** Bain's relative-NPS growth figures and CEB's CES
  superiority claim are proprietary and largely un-replicated — cite them *as* contested.[^bain][^dixon2010]
- **Averaging across the satisfaction-loyalty curve.** Asymmetry + non-linearity mean a mean CSAT/NPS
  misleads; segment and look at the tails (top-box, detractors).[^andersonmittal][^jonessasser]
- **Chasing delight before fixing must-be's.** Kano + negativity bias: basic-attribute failures destroy
  satisfaction faster than delighters build it.[^kano1984][^andersonmittal]
- **Reading a healthy account-level B2B score as safety.** It can hide a dissatisfied decision-maker or a
  hostage held only by switching costs.[^austen][^jonessasser]
- **Treating "satisfied" as "loyal" as "advocating."** They are distinct constructs; satisfaction is a
  weak, non-universal precursor to loyalty.[^oliver1999]

## Cross-references (do NOT duplicate these here)

- **Churn / retention / CLV analytics & survival models** → `da-23-customer-lifetime-value`,
  `da-24-survival-analysis` (this skill is the satisfaction/loyalty *constructs*; those own the math).
- **Writing NPS/CSAT survey items, scale design, response ops** → `career-and-formal-writing`
  (survey-question-writing) and **writing feedback replies at scale** → `content-and-marketing-writing`
  (nps-response-writing).
- **Proving/realizing the business value the customer set out to achieve, outcome-based CS, TTV** →
  `value-realization-outcome-cs`.
- **Driving adoption via motivation/habit/behavior change** → `behavior-change-psychology` (this hub).
- **Pre-purchase enterprise buying-committee psychology** (who decides, "no decision", champion/mobilizer,
  switching-cost *as a buying fear*) → `enterprise-b2b-buyer-psychology` (this hub). This skill is the
  **post-purchase** satisfaction/loyalty counterpart.
- **De-escalating an angry customer / reading the room** → `emotion-and-affect-psychology` (this hub);
  **building/repairing trust** → `trust-and-psychological-safety` (this hub).

## References (sources)

Confidence: **[F]** well-replicated / primary; **[Q]** qualified (2 sources or single authoritative);
**[C]** contested or vendor-origin (disagreement preserved). Volatile vendor claims are dated.

[^reichheld2003]: **[F]** Reichheld, F. F. (2003). "The One Number You Need to Grow." *Harvard Business Review* 81(12). https://hbr.org/2003/12/the-one-number-you-need-to-grow — NPS definition, 0–10 scale, promoter/passive/detractor cuts.
[^keiningham]: **[F]** Keiningham, Cooil, Andreassen & Aksoy (2007). "A Longitudinal Examination of Net Promoter and Firm Revenue Growth." *Journal of Marketing* 71(3):39–51. https://journals.sagepub.com/doi/10.1509/jmkg.71.3.039 — fails to replicate NPS superiority; "not superior to the ACSI."
[^dehaan]: **[F]** de Haan, Verhoef & Wiesel (2015). "The predictive ability of different customer feedback metrics for retention." *Int'l J. of Research in Marketing* 32(2):195–206. https://www.sciencedirect.com/science/article/abs/pii/S0167811615000324 — top-2-box CSAT best overall; best metric is industry-dependent; combining beats any single metric.
[^kristensen]: **[Q]** Kristensen & Eskildsen (2011/2014). "Is the Net Promoter Score a trustworthy performance measure?" *TQM Journal*. https://www.emerald.com/insight/content/doi/10.1108/TQM-03-2011-0021/full/html — NPS a poor predictor, distribution-sensitive.
[^fisher]: **[Q]** Fisher & Kordupleski (2019). "Good and bad market research: A critical review of Net Promoter Score." *Applied Stochastic Models in Business & Industry* 35(1):138–151. https://ideas.repec.org/a/wly/apsmbi/v35y2019i1p138-151.html (preprint arXiv:1806.10452).
[^measuringu]: **[F]** Sauro, J. — MeasuringU summaries of the NPS-growth literature. https://measuringu.com/nps-growth/ — NPS and CSAT roughly comparable; causation unestablished; underpowered datasets on both sides.
[^bain]: **[C]** Bain & Company, "The Numbers behind the Net Promoter System" (as of 2026). https://www.netpromotersystem.com/about/numbers-behind-the-net-promoter-system/ — relative/competitive NPS, 10–70% growth-variation and ~2×-faster claims; Bain proprietary, not independently replicated.
[^dixon2010]: **[C]** Dixon, Freeman & Toman (2010). "Stop Trying to Delight Your Customers." *Harvard Business Review* 88(7/8). https://hbr.org/2010/07/stop-trying-to-delight-your-customers — origin of CES / "the effort effect"; CEB research, claim that CES beats CSAT/NPS in service contexts is largely un-replicated.
[^greatbrook]: **[F]** Great Brook Consulting (Van Bennekom), quoting *The Effortless Experience* on CES 1.0 vs CES 2.0 wording/scales. https://greatbrook.com/effortless-experience-questionnaire-design-survey-administration-issues/
[^etouchpoint]: **[Q]** eTouchPoint, "The Customer Effort Score: A Primer" (CES 2.0, Net Easy Score, formulas). https://www.etouchpoint.com/customer-effort-score-primer-recommendations-cx-practitioners/
[^supportlogic]: **[F]** SupportLogic, "Limitations of CSAT and NPS." https://www.supportlogic.com/resources/blog/your-customer-satisfaction-metrics-dont-show-the-full-picture/ — self-selection, survival, recency bias shared across metrics.
[^silabs]: **[Q]** SI Labs, "Customer Satisfaction Score (CSAT)" incl. CSAT/NPS/CES comparison. https://www.si-labs.com/en/articles/customer-satisfaction-score/
[^qualtrics]: **[Q]** Qualtrics, "What is CSAT?" https://www.qualtrics.com/en-gb/articles/customer-experience/what-is-csat/ — top-box/top-two-box reporting; transactional vs relationship CSAT.
[^oliver1980]: **[F]** Oliver, R. L. (1980). "A Cognitive Model of the Antecedents and Consequences of Satisfaction Decisions." *Journal of Marketing Research* 17(4):460–469 — expectancy-disconfirmation paradigm.
[^kano1984]: **[F]** Kano, N., Seraku, N., Takahashi, F. & Tsuji, S. (1984). "Attractive Quality and Must-Be Quality." *Hinshitsu / J. Japanese Society for Quality Control* 14(2):147–156 — the five quality categories.
[^pzb1985]: **[F]** Parasuraman, Zeithaml & Berry (1985). "A Conceptual Model of Service Quality and Its Implications for Future Research." *Journal of Marketing* 49(4):41–50. https://eli.johogo.com/Class/parsus.pdf — the Gaps Model; Gap5 = f(Gap1–4); 10 determinants.
[^pzb1988]: **[F]** Parasuraman, Zeithaml & Berry (1988). "SERVQUAL: A Multiple-Item Scale…" *Journal of Retailing* 64(1):12–40 — 22-item instrument, 5 dimensions (P−E).
[^cronin1992]: **[F]** Cronin & Taylor (1992). "Measuring Service Quality: A Reexamination and Extension." *Journal of Marketing* 56(3):55–68 — SERVPERF (performance-only) challenge. https://bishtref.com/articles/10.1177/002224299205600304
[^carrillat]: **[Q]** Carrillat, Jaramillo & Mulki (2007). "The validity of the SERVQUAL and SERVPERF scales: a meta-analytic view…" *Int'l J. Service Industry Management* 18(5):472–490. https://www.emerald.com/insight/content/doi/10.1108/09564230710826250/full/html — about equally valid; SERVQUAL more diagnostic.
[^spc1994]: **[F]** Heskett, Jones, Loveman, Sasser & Schlesinger (1994). "Putting the Service-Profit Chain to Work." *Harvard Business Review* 72(2):164–174 (HBR Classic R0807L). https://hbr.org/2008/07/putting-the-service-profit-chain-to-work
[^zerodefections]: **[F]** Reichheld & Sasser (1990). "Zero Defections: Quality Comes to Services." *Harvard Business Review* 68(5):105–111. https://hbr.org/1990/09/zero-defections-quality-comes-to-services — 5% defection cut → ~25–85% profit lift.
[^hogreve]: **[F]** Hogreve, Iseke, Derfuss & Eller (2017). "The Service–Profit Chain: A Meta-Analytic Test…" *Journal of Marketing* 81(3):41–61 — links significant but effect sizes vary; original model fit poor.
[^silvestro]: **[F]** Silvestro & Cross (2000), *Int'l J. Service Industry Management* 11(3):244–268; Silvestro (2002), *Int'l J. Operations & Production Management* 22(1):30–49 — challenge the satisfaction mirror; inverse employee-sat↔profit in a UK grocer.
[^jonessasser]: **[F]** Jones, T. O. & Sasser, W. E. Jr. (1995). "Why Satisfied Customers Defect." *Harvard Business Review* 73(6):88–99 — Xerox 6×; non-linear curve by competitive intensity; apostle model (apostles/mercenaries/hostages/terrorists); false loyalty.
[^andersonmittal]: **[F]** Anderson, E. W. & Mittal, V. (2000). "Strengthening the Satisfaction-Profit Chain." *Journal of Service Research* 3(2):107–120 — asymmetry and non-linearity in the chain links; builds on Mittal, Ross & Baldasare (1998), *Journal of Marketing* 62(1):33–47.
[^oliver1999]: **[F]** Oliver, R. L. (1999). "Whence Consumer Loyalty?" *Journal of Marketing* 63 (Special Issue):33–44. https://foster.uw.edu/wp-content/uploads/2016/07/12_Oliver_1999.pdf — loyalty definition, four phases (cognitive→affective→conative→action), satisfaction as weak/necessary precursor, explicit B2C scope caveat.
[^dickbasu]: **[Q]** Dick, A. S. & Basu, K. (1994). "Customer Loyalty: Toward an Integrated Conceptual Framework." *Journal of the Academy of Marketing Science* 22(2):99–113 — attitudinal × behavioral loyalty 2×2 (true/latent/spurious/no loyalty).
[^austen]: **[F]** Austen, V., Herbst, U. & Bertels, V. (2012). "When 3 + 3 does not equal 5 + 1: New insights into the measurement of industrial customer satisfaction." *Industrial Marketing Management* 41(6):973–983. https://www.sciencedirect.com/science/article/abs/pii/S0019850111002483 — buying-center satisfaction is non-linear in individual members; negative judgments weigh more.
[^sis]: **[Q]** SIS International / practitioner consensus (as of 2026) on B2B CX: 7–12-stakeholder buying committees, role-divergence as risk, churn risk in un-surveyed roles. https://www.sisinternational.com/ — practitioner-sourced operational guidance (underlying constructs are academically grounded).
