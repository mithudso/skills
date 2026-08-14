<!-- hub-reference-banner -->
> **Reference file — part of the `da-applied-and-communication` hub.** Formerly the standalone `da-11-ethics-and-privacy` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-11-ethics-and-privacy
title: Data Ethics and Privacy
version: "1.0.0"
updated: "2026-05-30"
category: data-analysis
origin: local
description: >
  Data ethics and privacy for the working analyst — bias and fairness (including
  the impossibility theorems), GDPR / CCPA / HIPAA / EU AI Act (2026), differential
  privacy, k-anonymity / l-diversity / t-closeness, synthetic data, IRB review,
  model cards / datasheets for datasets / data statements, and AI alignment
  considerations for analysis pipelines.
  TRIGGER: user is auditing a model or dataset for bias, designing a privacy-
  preserving release, mapping a project to GDPR/HIPAA/CCPA, applying differential
  privacy, considering re-identification risk, writing a model card or
  datasheet, preparing for IRB review, or asks about EU AI Act compliance.
  SKIP: cluster-level encryption at rest (mongodb-encryption); MongoDB-specific
  compliance certification questions (mongodb-compliance); broader org security
  audit (security-compliance-auditor); cryptography internals.
triggers:
  - data ethics
  - bias and fairness
  - GDPR analysis
  - HIPAA analysis
  - CCPA
  - EU AI Act
  - differential privacy
  - k-anonymity
  - synthetic data privacy
  - model card
  - datasheet for datasets
  - IRB review
  - reidentification risk
keywords:
  - bias
  - fairness
  - demographic-parity
  - equalized-odds
  - calibration
  - GDPR
  - CCPA
  - HIPAA
  - EU-AI-Act
  - differential-privacy
  - epsilon-delta
  - k-anonymity
  - l-diversity
  - t-closeness
  - synthetic-data
  - IRB
  - model-card
  - datasheet
  - Fairlearn
  - AIF360
when_to_use:
  - Auditing a model or dataset for bias / fairness
  - Designing a privacy-preserving data release
  - Mapping a project to GDPR / HIPAA / CCPA / EU AI Act compliance
  - Applying differential privacy to outputs or training
  - Assessing re-identification risk on a "de-identified" dataset
  - Writing a model card or datasheet for datasets
  - Preparing for IRB / ethics-board review
when_not_to_use:
  - Cluster-level encryption at rest — use mongodb-encryption
  - MongoDB-specific compliance certifications — use mongodb-compliance
  - General security review of an application — use security-compliance-auditor
  - Cryptography algorithm choice — use webcrypto-vault-reviewer or domain-specific skills
related_skills:
  - mongodb-compliance
  - security-compliance-auditor
  - mongodb-encryption
  - da-1-6-3-reproducibility-replicability
  - da-12-ab-testing-causal-inference
---

# Data Ethics and Privacy

The discipline of analyzing data without harming the people the data describes. This skill covers the working analyst's view: fairness audits, privacy techniques, the regulatory landscape, and the documentation artifacts that make the work auditable.

## When to use this skill

Activate when the user:
- is auditing a model or dataset for bias / fairness
- is designing a privacy-preserving data release
- is mapping a project to GDPR / HIPAA / CCPA / EU AI Act
- is applying differential privacy
- is assessing re-identification risk on "de-identified" data
- needs a model card or datasheet for datasets
- is preparing for IRB / ethics-board review

## When NOT to use this skill

- Cluster encryption at rest → `mongodb-encryption`
- MongoDB compliance certifications → `mongodb-compliance`
- App-level security review → `security-compliance-auditor`
- Cryptography internals → cryptography-specific skills

---

## Bias and fairness

### Types of bias

| Type | Definition | Example |
|---|---|---|
| **Selection bias** | Sample is non-random | Survey only respondents who opt in |
| **Measurement bias** | Instrument is wrong | Scale that under-weighs above 200 lb |
| **Label bias** | Ground-truth labels reflect historical bias | "Hired" labels reflect biased past hiring |
| **Sampling bias** | Some groups under-represented | Face-recognition training set 80% white |
| **Aggregation bias** | One model fits all subgroups poorly | One model for all ages misses elderly patterns |
| **Deployment bias** | Model used differently than designed | Trained for screening, used for sentencing |
| **Evaluation bias** | Test set unrepresentative | Benchmark covers only common cases |

### Fairness metrics

For a model with prediction `Ŷ`, true label `Y`, and protected attribute `A` (e.g., race, gender):

- **Demographic parity** — `P(Ŷ=1 | A=0) = P(Ŷ=1 | A=1)`. Same positive-prediction rate across groups. Reasonable when base rates *should* be equal.
- **Equalized odds** — `P(Ŷ=1 | Y=y, A=0) = P(Ŷ=1 | Y=y, A=1)` for `y ∈ {0,1}`. Same TPR and FPR across groups.
- **Equal opportunity** — Equalized odds restricted to `Y=1`. Same TPR.
- **Calibration within groups** — `P(Y=1 | Ŷ=s, A=a)` is the same for each group `a` at each predicted score `s`.

### The impossibility theorems

Chouldechova (2017) and Kleinberg/Mullainathan/Raghavan (2017) showed: **calibration and equalized odds cannot both hold unless base rates are equal across groups, or the predictor is perfect.** You must choose which definition matters for your context. The political-philosophy debate is unavoidable; the math will not pick for you.

### Fairness toolkits

| Toolkit | Language | Strength |
|---|---|---|
| **Fairlearn** (Microsoft) | Python | Mitigation algorithms + dashboard |
| **AIF360** (IBM) | Python / R | 70+ metrics and 10+ mitigation algorithms |
| **What-If Tool** (Google) | Browser-based | Interactive what-if exploration |
| **`fairness` package** | R | Statistical fairness tests |

Pattern: pick 2-3 metrics relevant to the use case, compute them, document the choice and tradeoffs in the model card.

---

## The regulatory landscape (2026)

### GDPR (EU, in force since 2018)

Key principles for analysts:
- **Purpose limitation** — data collected for purpose X can't be used for unrelated purpose Y without new consent
- **Data minimization** — collect only what's necessary
- **Right to erasure (Art. 17)** — users can demand deletion. Hard if your training data is fixed.
- **Lawful basis** — consent, contract, legitimate interest, vital interest, public task, legal obligation (one of these must apply)
- **DPIA** (Data Protection Impact Assessment) — required for high-risk processing
- **Cross-border transfers** — Schrems II killed Privacy Shield; SCCs + transfer impact assessments needed
- Maximum fine: 4% of global annual revenue or €20M, whichever is higher

### CCPA / CPRA (California, in force since 2020 / 2023)

- Right to know what's collected, right to delete, right to opt out of sale
- **Sensitive personal information** category (CPRA, 2023): geolocation, biometric, racial/ethnic, religious, sexual orientation, immigration status, health
- Applies to businesses meeting thresholds (revenue, volume of CA residents' data, or % of revenue from selling data)

### HIPAA (US healthcare)

- **PHI** = Protected Health Information (any identifiable health data held by a covered entity or business associate)
- **Safe Harbor de-identification** — remove 18 specific identifiers (name, geo smaller than state, dates more granular than year, phone, fax, email, SSN, MRN, plan beneficiary number, account, certificate, license, vehicle ID, device ID, URLs, IP, biometric, photos, other unique identifying numbers/codes)
- **Expert determination** — alternative to Safe Harbor; a qualified statistician certifies low re-identification risk
- **Business Associate Agreement (BAA)** required for any vendor that touches PHI

### EU AI Act (2024-2026 phased)

Risk tiers, with phased application:
- **Unacceptable risk** — banned (social scoring, manipulation, real-time biometric in public spaces with exceptions)
- **High risk** — strict obligations (CE marking, risk management system, data governance, transparency, human oversight, accuracy/robustness/cybersecurity, post-market monitoring). Examples: AI in critical infrastructure, education, employment, essential services, law enforcement
- **Limited risk** — transparency obligations (e.g., chatbots must disclose they're AI)
- **Minimal risk** — most current AI (spam filters, video games); no specific obligations
- **General Purpose AI (GPAI)** models above a compute threshold (10^25 FLOPS) have additional obligations

Maximum fine: 7% of global annual turnover or €35M for prohibited AI; lower tiers for other violations.

### Comparison

| Region | Law | Trigger | Max penalty |
|---|---|---|---|
| EU | GDPR | Personal data of EU residents | 4% revenue or €20M |
| California | CCPA/CPRA | Business thresholds + CA residents | $7,500 per intentional violation |
| US healthcare | HIPAA | PHI handled by covered entity / BA | $1.5M/year/violation category |
| EU | AI Act | AI system put on the EU market | 7% revenue or €35M |
| Brazil | LGPD | Similar to GDPR | 2% revenue, capped at R$50M |

---

## Privacy-preserving techniques

### Differential privacy (DP)

A mathematical guarantee: an algorithm `A` is `(ε, δ)`-differentially private if for any two datasets `D` and `D'` differing in one record, and any output `S`:

```
P(A(D) ∈ S) ≤ exp(ε) · P(A(D') ∈ S) + δ
```

Intuition: a single person's data has bounded influence on the output. Smaller `ε` = stronger privacy = less utility.

- **ε ≤ 1**: strong
- **1 < ε ≤ 10**: moderate (most production deployments)
- **ε > 10**: weak; mostly a marketing claim

Mechanisms:
- **Laplace mechanism** — add Laplace noise to count / sum queries
- **Gaussian mechanism** — Gaussian noise; pairs with `(ε, δ)` framing (δ > 0)
- **Exponential mechanism** — for non-numeric outputs

Deployments:
- **Apple** uses local DP (ε ≈ 4-8 per query) for QuickType, emoji, Health
- **Google RAPPOR** for browser telemetry
- **US Census 2020** added DP noise to aggregate releases

DP is the only privacy technique with a provable mathematical guarantee. The cost is utility loss; the gain is composability across queries (the privacy budget).

### k-anonymity, l-diversity, t-closeness

Older syntactic models. A dataset is **k-anonymous** if every record's quasi-identifiers (attributes that could re-identify when combined) match at least k-1 others.

Failure modes:
- **Homogeneity attack** — all k records share a sensitive value (l-diversity addresses this by requiring l distinct sensitive values per group)
- **Background knowledge attack** — attacker knows the target's age and zip, narrows to a group, learns sensitive attribute
- **Skewness attack** — group's distribution of the sensitive attribute differs from the population's (t-closeness addresses this)

Famous breaches that motivated DP:
- **Netflix Prize (2006)** — Narayanan and Shmatikov re-identified users by linking the anonymized ratings to public IMDB profiles
- **AOL search log release (2006)** — user 4417749 identified by NYT from their queries within days
- **Sweeney (1997)** — re-identified the Massachusetts governor in a "de-identified" hospital release using just ZIP + birthdate + gender

The takeaway: **"de-identified" via field removal is not privacy.** Use differential privacy if you have a privacy guarantee to meet.

### Synthetic data

Generate fake records that statistically resemble real data. Tools: SDV (Synthetic Data Vault), Mostly AI, Hazy, Gretel.

**The trap**: synthetic data does not automatically inherit privacy. A GAN trained on a real dataset can memorize and reproduce records. To get a privacy guarantee, train the generator with differential privacy. Utility-vs-privacy tradeoff is real and dataset-specific.

---

## IRB and the Common Rule

Required when human-subjects research is involved (typically academic / publicly-funded). The Belmont principles:

1. **Respect for persons** — informed consent, special protections for vulnerable populations
2. **Beneficence** — minimize harm, maximize benefit
3. **Justice** — fair distribution of research benefits and burdens

Industry data science usually doesn't require IRB review, but the same principles are a reasonable ethical floor. Some companies (Microsoft Research, Google DeepMind) maintain their own internal ethics boards.

---

## Algorithmic accountability documentation

### Model cards (Mitchell et al, 2019)

A short structured document for each model release covering: model details, intended use, factors (e.g., demographic groups), metrics (including fairness), evaluation data, training data, quantitative analyses, ethical considerations, caveats. Hugging Face hub adopted them; pull requests for missing model cards are now routine.

### Datasheets for datasets (Gebru et al, 2018)

The analog for datasets: motivation, composition, collection process, preprocessing, uses, distribution, maintenance. Particularly useful before a public release of a benchmark or training corpus.

### Data statements (Bender & Friedman, 2018)

A linguistically-focused analog for NLP datasets: speaker demographics, language variety, annotator demographics, speech situation.

A reasonable practice: every model and dataset shipped to production includes a model card / datasheet committed alongside the artifact in version control.

---

## AI alignment considerations for analysis pipelines

Two emerging risks specific to LLM-in-the-loop analysis:

1. **Prompt injection in data pipelines.** If the analysis pipeline pipes user-controlled text into a prompt (e.g., "summarize these customer comments"), an adversarial comment can change the pipeline's behavior. Mitigations: separate untrusted input from system instructions, use the structured prompt format, escape / mark with `<untrusted_...>` tags as the mdb-tam server already does for LLM calls (`server/src/live/recommender.js`).
2. **Sycophancy and confirmation bias from LLMs.** An LLM analyst will tend to agree with the framing of the question. If you ask "is this a strong effect?" you get a different answer than "is this a weak effect?". Counter with structured prompts that ask for both directions.

---

## Practical checklists

### Pre-release privacy checklist

- [ ] What's the lawful basis for processing this data?
- [ ] Is the use compatible with the original collection purpose?
- [ ] What are the quasi-identifiers in this dataset?
- [ ] Have we tested for re-identification with a realistic attacker model?
- [ ] What's the privacy guarantee? (DP ε, k-anonymity k, ...)
- [ ] Is a DPIA required? If so, has it been done?
- [ ] Are cross-border transfers in scope? SCCs / adequacy decision in place?
- [ ] Is there a deletion / data-subject-request workflow?

### Pre-release fairness checklist

- [ ] What protected attributes are in scope? (Even if not in the model, are they recoverable from features?)
- [ ] Which fairness metric does this use case require, and why?
- [ ] What's the metric value across groups? Disparate impact ratio < 0.8?
- [ ] Have we tested for intersectional subgroups (e.g., Black women, not just "Black" and "women" separately)?
- [ ] Who is harmed by a false positive vs a false negative? Are those costs distributed evenly?
- [ ] Is the model card committed?

---

## References

1. Barocas, S., Hardt, M., & Narayanan, A. (2023). *Fairness and Machine Learning*. fairmlbook.org
2. Dwork, C., & Roth, A. (2014). *The Algorithmic Foundations of Differential Privacy*.
3. Mitchell, M. et al. (2019). "Model Cards for Model Reporting." FAT*.
4. Gebru, T. et al. (2018). "Datasheets for Datasets." arXiv:1803.09010
5. Bender, E. M. & Friedman, B. (2018). "Data Statements for Natural Language Processing." TACL.
6. Chouldechova, A. (2017). "Fair prediction with disparate impact." *Big Data*.
7. Kleinberg, J., Mullainathan, S., Raghavan, M. (2017). "Inherent Trade-Offs in the Fair Determination of Risk Scores." ITCS.
8. GDPR text — https://gdpr-info.eu/
9. CCPA/CPRA — https://oag.ca.gov/privacy/ccpa
10. HIPAA Privacy Rule — https://www.hhs.gov/hipaa/
11. EU AI Act — https://artificialintelligenceact.eu/
12. Sweeney, L. (2002). "k-Anonymity: A Model for Protecting Privacy."
13. Narayanan, A. & Shmatikov, V. (2008). "Robust De-anonymization of Large Sparse Datasets."
14. Microsoft Fairlearn — https://fairlearn.org/
15. IBM AIF360 — https://aif360.res.ibm.com/
