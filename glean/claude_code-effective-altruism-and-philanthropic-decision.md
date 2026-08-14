# effective-altruism-and-philanthropic-decision

**Category:** Frontend & Web Development
**Platform:** Claude Code
**Original Path:** claude-code/effective-altruism-and-philanthropic-decision

## Description
Guide evidence-based, normative decisions about where to give and how to maximize philanthropic impact. TRIGGER: effective altruism (EA), Singer's drowning-child/shallow-pond argument, MacAskill, Toby Ord; GiveWell charity evaluation, cost-per-outcome, cost per QALY/DALY; scale-neglectedness-tractability (SNT/ITN) framework; earning to give; longtermism critique; analytic vs emotional giving tension; measurability bias; EA critiques / post-FTX; moral circle expansion; philanthropic portfolio strategy; value pluralism in giving; "where should I donate to do the most good"; "is EA effective"; "how to evaluate a charity rigorously". SKIP: donor psychology / warm-glow / compassion fade → psychology-of-charitable-giving; moral foundations → moral-psychology; cognitive biases / decision heuristics → behavioral-decision-making; donor campaign / marketing → venture-cause-nonprofit-marketing; 501(c)(3) formation → venture-nc-nonprofit-formation; fundraising ops → venture-nonprofit-fundraising-ops.

---

# Effective Altruism & Philanthropic Decision-Making

Spoke of the `applied-psychology` hub. This file covers the **normative/decision** layer of
philanthropic giving: the intellectual framework for deciding *where* and *how much* to give
to maximize impact. For the psychological mechanisms that drive *why* people give emotionally,
see `psychology-of-charitable-giving`.

---

## 1. Foundations — Singer's Argument & the EA Movement

### The Drowning-Child Thought Experiment

Peter Singer's 1972 essay *Famine, Affluence, and Morality* introduced the argument that
reframes charitable giving as a moral obligation, not a supererogatory virtue:

> *If it is in our power to prevent something bad from happening, without thereby sacrificing
> anything of comparable moral importance, we ought, morally, to do it.*

The "shallow pond" case: you pass a drowning child. You can save them at trivial cost (muddy
clothes). Most people agree you must act. Singer's move: geographic distance is morally
irrelevant. A child dying of preventable disease in a low-income country triggers the same
obligation — and modern aid infrastructure removes the practical barrier. This remains "the
most famous argument in modern philosophy" in some philosophical accounts.

**EA's addition to Singer (circa 2009, Ord & MacAskill):** Singer's 1972 essay established
*that* we must give; EA added *where*. Toby Ord investigated cost-effectiveness variance
across charities and found orders-of-magnitude differences. Giving What We Can (2009) and
80,000 Hours formalized the two-part commitment: give significantly (≥10% of income) *and*
give to the highest-impact opportunities.

**Key figures:**
- Peter Singer — utilitarian philosopher; *The Most Good You Can Do* (2015)
- Will MacAskill — *Doing Good Better* (2015), *What We Owe the Future* (2022)
- Toby Ord — *The Precipice* (2020); founded Giving What We Can

---

## 2. Cost-Effectiveness Analysis

### Core Metrics

| Metric | Definition | EA use |
|--------|-----------|--------|
| DALY (Disability-Adjusted Life Year) | One year of healthy life lost | Compare disease-burden interventions |
| QALY (Quality-Adjusted Life Year) | One year at full health = 1.0 | Healthcare/policy comparisons |
| Cost-per-outcome | $ to achieve a measurable result (e.g., life saved) | GiveWell's primary unit |
| Moral weight | Relative value assigned to different outcomes | Converts heterogeneous benefits into a single comparable unit |

### GiveWell's Approach

GiveWell (founded 2006) operationalizes cost-effectiveness at scale:

1. **Benchmark:** Cost of "doubling consumption" for a person in poverty for one year = 1 unit.
2. **Funding bar:** Grants must clear ~10× the benchmark to be recommended. The bar adjusts with available funding: when GiveWell expects less total funding to distribute, the bar rises (only the most cost-effective grants qualify); when more funding is available, the bar lowers (more opportunities can be funded).
3. **Model building:** GiveWell builds its own cost-effectiveness models rather than relying on published estimates (which can contain large errors — one published figure was off by ~100×).
4. **Moral weights:** Explicitly acknowledged as subjective, grounded in surveys and global health metrics. Averting the death of a child ≈ 100 units; averting one year of disability ≈ 2.3 units.

**Critical honesty:** GiveWell's own documentation states: *"We believe that cost-effectiveness
estimates should not be taken literally due to the significant uncertainty around them."*
Models are comparative tools, not precise measurements.

### Expected-Value Reasoning

EA relies heavily on expected-value (EV) reasoning: multiply outcome value × probability of
achieving it. This enables comparison of speculative high-upside interventions (pandemic
preparedness) against near-certain lower-upside ones (bednet distribution).

**Tension:** "Shut up and multiply" can produce conclusions that violate common-sense moral
intuitions. Critics argue EV maximization without side-constraints is a root cause of
intellectual and cultural failures in EA (including the FTX episode — see §5).

---

## 3. Cause Prioritization: Scale, Neglectedness, Tractability

The SNT (also called ITN: Importance, Tractability, Neglectedness) framework is the primary
tool for comparing cause areas when full cost-effectiveness modeling is impractical.

### The Framework

**Scale (Importance):** How big is the problem if unsolved? Measured by total welfare impact
— lives affected, severity, and breadth.

**Neglectedness:** How few resources are currently being directed here relative to the need?
Neglectedness signals opportunity *only* when caused by correctable structural reasons (market
failure, cultural blind spots, lack of awareness) — not when the problem is neglected simply
because it is genuinely intractable.

**Tractability (Solvability):** If resources doubled, what fraction of the problem would be
solved? Focuses on marginal returns to additional effort.

| Cause tier | Scale | Neglectedness | Tractability | EA verdict |
|------------|-------|---------------|-------------|------------|
| Global health (malaria, NTDs) | Enormous | Moderate | High (RCT-proven) | Top-tier |
| Factory farming / animal welfare | Vast (sheer numbers) | High | Moderate | Priority cause |
| Biosecurity / pandemic prevention | Catastrophic potential | High | Moderate | Growing priority |
| AI safety | Potentially existential | High (2024: less so) | Contested | Contested |
| Climate change (mitigation) | Enormous | Low (many funders) | Moderate | Less EA-prioritized |

**Known framework limitation:** The three factors correlate — neglected problems often remain
neglected because they are intractable. The framework can double-count or produce inflated
rankings. Use as a screening/intuition pump, not a precise algorithm. (EA Forum, 2024)

### Earning to Give

A career strategy derived from SNT reasoning: if your highest-talent path is finance, law,
or tech rather than direct nonprofit work, take the high-earning path and donate a large
fraction (commonly 20–50%) to high-impact organizations. The logic: financial capital is
fungible; labor bottlenecks are often not in giving roles.

**80,000 Hours guidance (2026):** Earning to give is most attractive when:
- You are an excellent fit for a high-earning but morally neutral career
- Target organizations are funding-constrained (can convert money into impact)
- You want career capital while maintaining social-impact engagement
- You are highly uncertain which cause is most pressing (money is flexible)

**Current nuance:** 80,000 Hours has shifted emphasis toward direct-work roles in EA-priority
areas (AI safety, biosecurity, policy) where talent — not funding — is the bottleneck. Earning
to give remains appropriate for many but is not the default recommendation.

### Moral Circle Expansion

Singer's argument generalizes beyond humans: if capacity to suffer is what grounds moral
concern, then non-human animals and future generations deserve moral consideration.

- **Animal welfare:** Factory farming causes suffering at a scale that rivals or exceeds
  global poverty by some welfare metrics. Organizations like Animal Charity Evaluators apply
  GiveWell-style rigor to this cause.
- **Future generations:** Longtermism (see §4) is the most contested extension.

---

## 4. Longtermism — a Contested Frontier

**The claim:** Future people matter as much as present people; because there could be
vastly many future people, improving long-run outcomes (e.g., reducing extinction risk)
has astronomical expected value, potentially dwarfing near-term interventions.

**Key texts:** MacAskill, *What We Owe the Future* (2022); Ord, *The Precipice* (2020).
Strong longtermism motivated much EA funding toward AI safety, biosecurity, and
existential-risk research.

**Contested on multiple fronts — honest assessment:**

| Objection | Source | Strength |
|-----------|--------|----------|
| Meta-epistemic: billion-year consequences are unknowable; EV calculations are unjustified | Schwitzgebel (2024), Oxford Global Priorities Institute | Strong |
| Neglects urgent, tractable present suffering | Singer AMA (2024); Mounk (2024) | Moderate |
| Brackets historical injustice; sidelines reparations | Schliesser (2024, Crooked Timber) | Moderate |
| Creates problematic scale-creep away from near-term causes | Internal EA Forum debates | Moderate |
| Existence value arguments philosophically contested (Parfit, person-affecting views) | Standard population ethics | Strong |

**Singer's own position (2024):** "Placing too much emphasis on longtermism hinders the
prospects of EA becoming a broad and popular movement... some ways of doing good are hundreds
of times more effective than others [in near-term causes]."

**Working position for practitioners:** Near-term global health interventions (malaria,
NTDs, direct cash transfers) have strong RCT-backed evidence. Longtermism-motivated work
carries high uncertainty but potentially very high upside. A diversified philanthropic
portfolio (see §7 matrix) may hedge this uncertainty rather than concentrate on one time
horizon.

---

## 5. Honest Critiques of EA

### Measurability Bias

EA systematically advantages interventions that are:
- Quantifiable (lives saved, DALYs averted)
- Testable via RCTs
- Tractable within a funding cycle

This creates a selection effect: arts, social cohesion, democracy, legal reform, and
structural/systemic change are systematically under-funded even when they may have large
diffuse impacts. GiveWell's own analysts acknowledge that moral weights are "ultimately
subjective" despite the quantitative scaffolding.

### The Post-FTX Reputation Crisis (2022–2024)

Sam Bankman-Fried (SBF), EA's wealthiest funder via the FTX Foundation, was arrested on
fraud charges in 2022 (convicted 2023) after allegedly misusing customer funds at scale —
publicly framing his actions in EA expected-value terms. Movement leaders explicitly
distanced themselves: "A clear-thinking EA should strongly oppose 'ends justify the means'
reasoning" (MacAskill). The structural critique that survived: EA's individual-trust model
over institutional checks, deference to high-status insiders, and "counter-cultural"
self-image made it susceptible to rationalized misconduct. (Matthews, Vox 2022; EA Forum
post-mortem 2024.)

### Value Pluralism Objections

EA's utilitarian/welfarist foundations do not capture all defensible moral values:
- **Rights and side-constraints:** Some outcomes wrong regardless of aggregate welfare gain
- **Relational obligations:** Special duties to family, community, existing relationships
- **Non-welfarist values:** Justice, dignity, cultural continuity, autonomy
- **Process values:** How good is achieved matters, not just the outcome magnitude

GiveWell acknowledges this: moral weights are judgment calls, not derivations.

### Overhead and Measurement Limits

"Overhead ratios" (admin/fundraising costs as % of total) are widely criticized as
misleading charity metrics. EA is correctly skeptical of overhead-based ranking. But
GiveWell's own models acknowledge that cost-effectiveness estimates contain "significant
uncertainty" and should not be "taken literally." The precision of the framework can
project false confidence.

---

## 6. The Head-vs-Heart Bridge (Cross-Reference Point)

The central tension in philanthropic decision-making:

**Analytic/effective giving:** Reason-driven, comparative, statistically framed. Favors
"statistical lives" (future beneficiaries identified by probability). Optimizes across cause
areas. May feel emotionally cold.

**Affective/emotional giving:** Driven by sympathy, personal connection, narrative. Responds
strongly to identifiable individual victims. Generates warm-glow and in-group preference.

**The critical empirical finding (Small, Loewenstein & Slovic, 2007):** Priming deliberative
thinking *reduces* donations to identifiable victims *without increasing* donations to
statistical victims — a net reduction in generosity. This is the central empirical problem for
EA's prescriptive project: effective-giving advocacy may suppress emotional giving faster than
it activates analytical giving, producing worse aggregate outcomes in the short run.

**This skill's boundary:** The *normative* question (how *should* one give, what frameworks
exist) belongs here. The *descriptive* question (why *do* people give or fail to give, what
psychological mechanisms drive donations) belongs to `psychology-of-charitable-giving`.

**Practitioner implication:** Efforts to move donors toward more effective giving must grapple
with this tension. Purely analytic appeals may depress total giving. Hybrid strategies that
preserve emotional engagement while adding cost-effectiveness filters may outperform either
pure approach.

---

## 7. Quick Reference: Decision Framework for Philanthropic Giving

### Cause Prioritization Checklist (SNT)

```
1. SCALE — How many people/beings affected, how severely?
2. NEGLECTEDNESS — Is this underfunded for bad reasons (not mere intractability)?
3. TRACTABILITY — Would marginal resources move the needle meaningfully?
4. EVIDENCE BASE — RCT-proven vs. plausible vs. speculative?
5. PERSONAL FIT — Where can you uniquely contribute (labor or dollars)?
```

### Giving Strategy Matrix

| Donor profile | Recommended approach |
|---------------|----------------------|
| Values evidence, near-term impact | GiveWell top charities (malaria, NTDs, cash transfers) |
| Values animal welfare | Animal Charity Evaluators top picks |
| Wants portfolio diversification across time horizons | Near-term (50–70%) + longtermism (30–50%) |
| Career: high-earning, uncertain direct-work fit | Earning to give (10–50% of income) |
| Career: direct-work fit in priority area | Pursue talent role; give what is feasible |
| Values local community, relational ties | Accept value pluralism; blend effective + local |
| Already gives to specific charities; wants to evaluate | Apply Red Flags checklist; check GiveWell / ACE listing; run SNT screen on cause area |

### Charity Evaluation Red Flags

- Overhead ratio used as primary metric (misleading)
- No external evaluation or evidence of effectiveness
- Vague outcome claims without counterfactual reasoning
- Single dramatic narrative without population-level impact data
- Cost-effectiveness estimates not published or auditable

---

## 8. Sources & Further Reading

Full bibliography (17 sources) in `references/sources.md`. Key anchors used in this file:

- Singer (1972) — the drowning-child/shallow-pond foundational argument
- GiveWell (2025) — cost-effectiveness methodology and moral weights
- Small, Loewenstein & Slovic (2007) — deliberative thought reduces charitable giving
- 80,000 Hours (2016/2026) — SNT framework and earning-to-give guidance
- Schwitzgebel (2024), Mounk (2024), Schliesser (2024) — longtermism and EA critiques
- Matthews / Vox (2022), EA Forum (2024) — post-FTX structural analysis

---

## Cross-Skill Routing

| Topic | Where it lives |
|-------|---------------|
| Why people give emotionally (warm-glow, identifiable victim, compassion fade) | `psychology-of-charitable-giving` |
| Donor retention, ask psychology, stewardship, lapsed-donor win-back | `fundraising-and-donor-psychology` |
| Moral foundations, fairness perception, values conflict | `moral-psychology` |
| Cognitive biases, expected value, dual-process reasoning | `behavioral-decision-making` |
| Nonprofit campaign design, donor acquisition, marketing | `venture-cause-nonprofit-marketing` |
| Fundraising operations, donor cultivation, major gifts | `venture-nonprofit-fundraising-ops` |
| Nonprofit 501(c)(3) formation and legal structure | `venture-nc-nonprofit-formation` |
| Donation psychology in health-behavior context | `health-behavior-change-and-donor-registration` |