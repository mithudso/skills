# Satisfaction Theory — Expectancy-Disconfirmation & the Kano Model

> Provenance: reference under `applied-psychology ▸ customer-satisfaction-loyalty-science`. Created
> 2026-06-18 via `/dr`. Confidence: **[F]** well-replicated/primary · **[Q]** qualified · **[C]**
> contested.

## Contents
- 1. Expectancy-Disconfirmation Theory (Oliver)
- 2. The Kano Model
- 3. Using them together (satisfiers vs delighters)
- Sources

---

## 1. Expectancy-Disconfirmation Theory (EDT / the disconfirmation paradigm)

**Origin & standing.** Richard L. Oliver (1977, 1980; book *Satisfaction: A Behavioral Perspective on
the Consumer*, 1997/2010). EDT is the **dominant paradigm** for how satisfaction forms. **[F]**

**The mechanism.** Satisfaction is a function of two things:
1. **Pre-purchase expectations** (a baseline/standard), and
2. **Disconfirmation** — the gap between expectations and **perceived performance**:
   - **Positive disconfirmation** (performance > expectations) → satisfaction
   - **Negative disconfirmation** (performance < expectations) → dissatisfaction
   - **Confirmation** (performance ≈ expectations) → neutral / "met" (not a strong positive). **[F]**

So satisfaction is **relative**, not absolute: the *same* objective performance produces satisfaction or
dissatisfaction depending on what the customer expected. **[F]**

**Components to know.**
- **Expectations** — the comparison standard (predictive, ideal, "should", desired, or minimum-tolerable
  — the polysemy is itself a measurement problem; see §critiques).
- **Perceived performance** — the customer's read of what was delivered.
- **Disconfirmation** — the subjective comparison (often measured directly: "better/worse than
  expected"), which is a *stronger* satisfaction driver than the arithmetic gap in many studies. **[Q]**

**Adjacent constructs.**
- **Assimilation-contrast effects** — small gaps get *assimilated* (perception pulled toward
  expectation); large gaps get *contrasted* (exaggerated). Explains why minor shortfalls are forgiven and
  large ones over-punished. **[Q]**
- **Zone of tolerance** — a band between "desired" and "adequate" service within which performance
  variation goes largely unnoticed; this concept bridges EDT to SERVQUAL (`service-quality-servqual.md`).
  **[Q]**

**Why it matters for a TAM.** The lever isn't only *delivery* — it's the **expectations you set**.
Over-promising (sales, marketing, your own commitments) raises the bar and manufactures negative
disconfirmation even when delivery is unchanged. Deliberate expectation-setting and "under-promise,
over-deliver" are direct applications. This is the same mechanism as SERVQUAL's Gap 4 (communication).

**Critiques.** Expectations are hard to measure and **polysemic** (respondents read "expect" as
prediction vs. ideal vs. norm), which destabilizes results; some studies find **perceived performance
alone** predicts satisfaction about as well as the disconfirmation gap (i.e., expectations add little in
some contexts — echoing the SERVPERF-vs-SERVQUAL debate). **[Q/C]**

## 2. The Kano Model

**Origin.** Noriaki Kano et al. (1984), "Attractive Quality and Must-Be Quality," *Hinshitsu / Journal of
the Japanese Society for Quality Control* 14(2):147–156 (originally Japanese). **[F]** It rejects the
single-axis assumption that "more quality = more satisfaction," replacing it with a **two-axis** map:
horizontal = degree of attribute *functionality/fulfillment*; vertical = resulting *satisfaction*.

**The five categories** (an attribute is classified by how satisfaction responds to its presence/absence):

| Category | Present | Absent | Curve | Strategic read |
| --- | --- | --- | --- | --- |
| **Must-be** (basic, dissatisfier) | Merely expected (no credit) | Strong dissatisfaction | Asymptotic on the downside | **Table stakes** — can't win on them, only lose |
| **One-dimensional** (performance, satisfier) | Satisfaction ↑ | Dissatisfaction ↑ | ~Linear | **Compete head-to-head**; more is better |
| **Attractive** (delighter, exciter) | Delight (disproportionate ↑) | No penalty | Asymptotic on the upside | **Differentiation**; unexpected value |
| **Indifferent** | No effect | No effect | Flat | **Stop investing** |
| **Reverse** | *Dissatisfaction* for some | Satisfaction | Inverted | **Watch segments** — some customers want the opposite |

**[F]** for the must-be / one-dimensional / attractive / indifferent / reverse taxonomy.

**The Kano questionnaire.** Each attribute is probed with a **functional/dysfunctional question pair**
("How do you feel if the product *has* feature X?" / "…if it *does not have* X?"), each answered on a
5-point like/expect/neutral/tolerate/dislike scale; the answer pair is cross-tabulated against a Kano
evaluation table to classify the attribute (with an "questionable" cell for inconsistent answers).
**[F]** A common refinement adds a **self-stated-importance** question to break ties / prioritize within
a category. **[Q]**

**The dynamic — "the natural decay of delight."** Categories **migrate over time**: today's *attractive*
delighters become tomorrow's *one-dimensional* expectations and eventually *must-be* basics as the market
adopts them (the classic example: a feature that delighted early then became table stakes). Implication:
delighters are a **depreciating asset** — you must keep finding new ones. **[Q]** (well-attested in the
practitioner/QM literature; precise longitudinal evidence is thinner, so flagged qualified.)

**Critiques.** Classification can be unstable across samples and over time; the functional/dysfunctional
question format is cognitively demanding and sensitive to wording; segment heterogeneity means one
attribute can land in different categories for different customers (the *reverse* and *indifferent* cells
expose this); and the model describes *response shape*, not *how much to invest*. **[Q]**

## 3. Using EDT + Kano together — satisfiers vs delighters

The two theories are complementary and resolve a common mistake (the reflex to "delight customers"):

- **EDT** tells you satisfaction is *relative to expectations* — so a delighter only delights while it's
  *unexpected*; once expected (Kano decay), it becomes a must-be whose absence now *dissatisfies*.
- **Kano** tells you **where to invest**: secure the **must-be's first** (their failure destroys
  satisfaction asymmetrically — see the negativity-bias finding in `loyalty-profit-chain.md` §asymmetry),
  compete on **one-dimensional** attributes, and use **attractive** delighters for differentiation —
  knowing they depreciate.
- **TAM application:** prioritize fixing basic-attribute failures (reliability, "it just works,"
  low-effort support — note this connects to CES and SERVQUAL Reliability) **before** investing in
  delight; and manage expectations so you're not manufacturing negative disconfirmation. "Stop trying to
  delight" (the CES research, `metrics-and-critiques.md` §3) is the service-context corollary of Kano +
  EDT: in support, *removing effort* (a must-be) beats *adding delight*.

## Sources

1. Oliver, R. L. (1980). "A Cognitive Model of the Antecedents and Consequences of Satisfaction Decisions." *Journal of Marketing Research* 17(4):460–469 — [paper] (expectancy-disconfirmation).
2. Oliver, R. L. (1997/2010). *Satisfaction: A Behavioral Perspective on the Consumer* (2nd ed., Routledge) — [book].
3. Tse, D. K. & Wilton, P. C. (1988). "Models of Consumer Satisfaction Formation." *Journal of Marketing Research* 25(2):204–212 — [paper] (disconfirmation operationalization).
4. Kano, N., Seraku, N., Takahashi, F. & Tsuji, S. (1984). "Attractive Quality and Must-Be Quality." *Hinshitsu / J. Japanese Society for Quality Control* 14(2):147–156 — [paper] (the five categories).
5. ASQ (American Society for Quality), "Kano Model Analysis" — https://asq.org/quality-resources — [docs/industry] (questionnaire & categories).
6. Zeithaml, Berry & Parasuraman (1993). "The Nature and Determinants of Customer Expectations of Service." *Journal of the Academy of Marketing Science* 21(1):1–12 — [paper] (zone of tolerance; desired vs adequate).
