# Metrics & Their Critiques — CSAT, NPS, CES

> Provenance: reference under `applied-psychology ▸ customer-satisfaction-loyalty-science`. Created
> 2026-06-18 via `/dr`. Confidence tags: **[F]** well-replicated/primary · **[Q]** qualified (2
> sources / one authoritative) · **[C]** contested or vendor-origin (disagreement preserved).
> Volatile vendor claims dated; verified-as-of 2026-06-18.

## Contents
- 1. CSAT — Customer Satisfaction Score
- 2. NPS — Net Promoter Score (and the "one-number" debate)
- 3. CES — Customer Effort Score (and "the effort effect")
- 4. Cross-cutting synthesis: which metric to trust
- Sources

---

## 1. CSAT — Customer Satisfaction Score

**What it measures.** CSAT captures **transactional / attribute satisfaction** — a customer's
here-and-now reaction to a specific, recent, concrete experience (a support conversation, an order, an
onboarding step), not the overall relationship. Vendors distinguish *transactional CSAT* (pinned to a
touchpoint) from *relationship CSAT* (broader brand satisfaction), but the metric's center of gravity is
transactional. **[F]**

**How computed.** A single item — "How satisfied were you with [X]?" — on a 1–5 (dominant), 1–7, or
1–10 scale. The standard report is the **top-box / top-two-box percentage**:
`CSAT% = (responses in the "satisfied" zone ÷ total) × 100`. On 1–5, "satisfied" = 4s and 5s
(top-two-box); on 1–10, usually 9–10. Multiple vendors cite top-two-box as the most accurate retention
predictor among cut-points (attributed to Qualtrics). It can also be reported as a mean ("composite"),
which is easy to misread. **[F]**

**Lineage.** The CSAT-as-percentage convention traces to the **American Customer Satisfaction Index
(ACSI)**, founded 1994 at the University of Michigan by economist **Claes Fornell**, modeled on the
Swedish Customer Satisfaction Barometer; the methodology established the 0–100 reporting scale.
**[Q]** — *the 1996 Fornell methodology paper was not independently fetched; treat the precise citation
as needing primary verification.*

**Strengths.** Simple, intuitive, flexible across any touchpoint; good for operational management and
pinpointing where a specific interaction fell short; easy to segment by queue/region/plan tier and
trend over time. **[F]**

**Documented critiques / limits.**
- **Weak predictor of future behavior/loyalty.** "Satisfied" customers defect routinely. Gallup's work
  argues satisfaction is "a relatively poor indicator of future customer behavior" and that only
  *emotionally* (not merely *rationally*) satisfied customers behave differently from dissatisfied ones;
  "rationally satisfied customers behave no differently than customers who are dissatisfied." **[F]**
- **Recency bias** — measures the moment, not the relationship. **[F]**
- **Self-selection / non-response (volunteer) bias** and **survival bias** — only the very happy or very
  angry respond; churned customers are never captured. **[F]**
- **Negative skew / ceiling effect.** Satisfaction distributions cluster at the top because dissatisfied
  customers leave, so there's "very little opportunity at the bottom end of the scale" — which is *why*
  practitioners treat mid-scale (5–6 on a 10-pt) as effectively dissatisfied. **[Q]**
- **Doesn't explain "why"** or capture effort — pair with Kano (feature prioritization), root-cause
  analysis, and CES (effort). **[F]**

## 2. NPS — Net Promoter Score

### The original claim (Reichheld, HBR 2003)

Frederick F. Reichheld, "The One Number You Need to Grow," *Harvard Business Review*, Dec 2003 (Reprint
R0312C). NPS was developed by Reichheld with **Satmetrix** and **Bain & Company** (joint trademark
holders). **[F]**

**Item & calculation.** One question — "How likely is it that you would recommend [company/product] to a
friend or colleague?" — on a **0–10** scale.
- **Promoters** = 9–10 · **Passives** = 7–8 · **Detractors** = 0–6
- `NPS = %Promoters − %Detractors` (−100 to +100; passives are in the base but not added). **[F]**

**The argument.** Based on ~2 years of research (with Satmetrix) linking survey answers to actual
purchasing/referral behavior and growth, Reichheld claimed "would you recommend" was the single best
predictor of top-line growth — because recommending puts the customer's *own reputation* on the line
("sacrifice"). The 2003 article reported that the recommend question was the best or second-best
predictor of growth in 11 of 14 case-study industries; the often-quoted larger benchmarking figures
(hundreds of companies, ~28 industries, a ~16% median, exemplars such as eBay/Amazon/USAA at 75–80%)
come from the broader Satmetrix dataset and Reichheld's later 2006 book *The Ultimate Question*, not
the 2003 article itself. **[Q — the 2003-vs-later-benchmark provenance is easy to conflate; figures
attributed to their correct source.]**

### The academic "one-number" critique (genuinely CONTESTED)

- **Keiningham, Cooil, Andreassen & Aksoy (2007)**, *Journal of Marketing* 71(3):39–51. Used longitudinal
  data (21 firms, 15,500+ interviews, Norwegian Customer Satisfaction Barometer, 2–4-year periods) to
  *replicate Reichheld's/Satmetrix's own methodology* and compare NPS against the **ACSI** — deliberately
  using the very industries Reichheld cited as NPS exemplars. **Finding:** "fails to replicate his
  assertions regarding the 'clear superiority' of Net Promoter"; flatly, **"Net Promoter is not superior
  to the ACSI for the data under investigation"**; R² values for NPS vs. ACSI "comparable (and
  statistically indistinguishable)." They also flagged that a supportive earlier study (Morgan & Rego
  2006) misunderstood the data fields so that "Net Promoter was not actually examined." **[F]**
- **Kristensen & Eskildsen (2011/2014)**, *TQM Journal* — controlled experiment, Danish insurance
  (~2,000 obs.); NPS "a very poor predictor of both customer loyalty and customer satisfaction," highly
  sensitive to the underlying distribution, lower precision than ACSI/EPSI. **[F]**
- **Fisher & Kordupleski (2019)**, *Applied Stochastic Models in Business & Industry* 35(1):138–151 — NPS
  "has not lived up to its claimed benefits," ignores **Relative Value** (the AT&T tradition), focuses on
  the user (often not the decision-maker), is uncalibrated by competitive info, has no agreed aggregation
  standard; "a false god." (Preprint arXiv:1806.10452.) **[F]**

**The recommend-intention vs. actual-behavior gap.** Christina Stahlkopf (C Space), "Where Net Promoter
Score Goes Wrong," *HBR* (Oct 2019), compared 2,000 consumers' NPS ratings against what they had
*actually done*: ~50% classified promoters but **~69% had actually recommended**; ~16% classified
detractors but **only ~4% had actually discouraged** anyone; detractors were "seven times more likely to
have either recommended a brand or said nothing at all than to have disparaged it." **[Q — figures from
secondary summaries; the HBR article is paywalled and the exact figures were not verified against the
original.]**

**Disconfirming the disconfirmers (preserve the nuance).** Jeff Sauro / MeasuringU notes the Keiningham
datasets, like Reichheld's own, have **insufficient statistical power** (few firms per industry) to
*definitively* dismiss NPS; across the literature NPS "isn't always clearly the best predictor… but even
when it's not, it's a close contender," and it correlates highly with CSAT. Honest summary: **NPS and
CSAT have roughly comparable predictive ability; neither is reliably superior; the NPS→growth causation
is unestablished and could run the other way.** **[F]**

**The Bain/Reichheld defense (the other side).** Bain argues critics attack *absolute bottom-up* scores,
whereas growth prediction relies on **relative / "top-down" / competitive-benchmark NPS**:
- *Bottom-up* (individual classification) predicts individual LTV (promoter-vs-detractor LTV gaps cited
  at 3–8×). **[C — Bain claim]**
- *Relative/top-down* (third-party benchmarking vs. direct competitors) is what Bain says explains growth:
  "differences in relative competitive NPS explain anywhere from **10% to 70%** of the variation in
  subsequent revenue growth among direct competitors," NPS leader grows "more than two times faster."
  **[C — Bain proprietary research, not independently replicated, as of 2026.]**
- Reichheld ("The Dubious Argument Against NPS," 2019) concedes NPS is widely *misused* — gamed, treated
  as an end in itself, tied to frontline comp — and argues that's the real problem.

**Benchmarking caveats (both sides partly agree).** Comparing **absolute** NPS across regions/cultures is
unreliable (e.g., Asian respondents may give fewer 9s/10s — cultural rating-scale bias); there is **no
agreed standard for aggregating NPS to a company level**; small companies often lack a valid benchmark
(some practitioners say ignore benchmarks, track your own trend). **[F]**

**Net assessment for the TAM:** NPS is operationally useful as a *relationship pulse* and a trigger for
closed-loop follow-up, but the academic literature does **not** support the "one number / singular
superior predictor of growth" framing. Treat it as a contested business claim, not an established fact.

## 3. CES — Customer Effort Score

**Origin.** Developed in **2010** by the **Corporate Executive Board (CEB; now Gartner)**, introduced in
Dixon, Freeman & Toman, "Stop Trying to Delight Your Customers," *HBR* Jul–Aug 2010 (Reprint R1007L).
**[F]**

**Core finding ("the effort effect").** A study of **75,000+** people interacting with contact-center
reps or self-service found "over-the-top" service/delight makes little difference to loyalty in
**service/support** contexts; customers want a **quick, low-effort resolution**. The article reported CES
"is a better predictor of loyalty than CSAT or NPS" *in customer-service interactions specifically*.
Article data (as reported by CEB): of low-effort customers, **94% intended to repurchase** and **88%
intended to increase spend**, while on the high-effort side **81% intended to spread negative
word-of-mouth**. (The two sides are reported on different measures — repurchase/spend for low-effort,
negative-WOM intent for high-effort — so read them as directional, not a single matched pair.)
**[C — CEB's own research, not independently replicated in the article; figures as summarized.]**

**The book.** Dixon, Toman & DeLisi, *The Effortless Experience: Conquering the New Battleground for
Customer Loyalty* (Portfolio/Penguin, Sept 2013). Headline stat: "**96% of customers with a high-effort
service interaction become more disloyal, compared to just 9% who have a low-effort experience.**" Five
disloyalty drivers (per the HBR article): repeat contacts, channel switching, transfers, repeating
information, and perceived/generic additional effort. **[F]** (existence & content of the book and its
claims) / **[C]** (the predictive-power figures are CEB-proprietary).

**CES 1.0 vs CES 2.0 (well-corroborated).**

| | CES 1.0 (2010) | CES 2.0 (~2013, current standard) |
| --- | --- | --- |
| Form | Direct question | Agreement statement |
| Wording | "How much effort did you personally have to put forth to handle your request?" | "[The company] made it easy for me to handle my issue." |
| Scale | **1–5, inverted** (1 = very low effort; *lower is better*) | **1–7 agreement** (1 = strongly disagree … 7 = strongly agree; *higher is better*) |
| Scoring | Mean of responses | **Top-three-box %** (5+6+7 ÷ total); sometimes 0–100 |

**[F]** — corroborated independently (Great Brook quoting the book pp.157–158, eTouchPoint, RateNow,
Heart of the Customer, MeasuringU). **Why the change:** the inverted scale caused false
positives/negatives; some read it as "how hard did *I* try" rather than "how hard did the company make
it"; "effort" translates poorly internationally — "easy" is more universal; CES 2.0 reduces
social-desirability/self-assessment bias. CEB guidance: move customers to *at least a 5* on the 7-pt
scale (1→5 boosts loyalty ~22%, 5→7 only ~2% — diminishing returns). **[C — CEB figures.]**

**Variant — Net Easy Score (British Telecom):** an NPS-style calc — %easy (top 2 of 7) minus %difficult
(bottom 3); BT found "easy"-rating customers ~40% less likely to churn. **[Q]**

**What CES does NOT measure / limits.**
- Measures **effort/ease only** — not satisfaction, delight, relationship loyalty, or "why." A low-effort
  interaction can still leave a customer unsatisfied with the *outcome*. **[F]**
- It's **transactional** and breaks as a relationship/annual metric. **[Q]**
- Scope is **service/support-centric**; the effort→loyalty finding doesn't claim effort dominates loyalty
  in purchase/product contexts. **[F]**
- **Comparability is murky** — mixed scales (5-pt vs 7-pt), mixed wording, three competing calc methods
  (mean / top-3-box / Net Easy); few clean public benchmarks. **[F]**

**Independent cross-metric check (critical).** de Haan, Verhoef & Wiesel (2015), *IJRM* 32(2):195–206 —
compared CSAT, NPS, CES across 93 firms / 18 industries with multi-level probit controlling for
self-selection. **Top-2-box CSAT predicted retention best overall; the best metric varies by industry
and unit of analysis; combining metrics predicts better than any single one.** Direct independent
evidence against *any* "one-number" claim — and it ranks CSAT ahead of NPS and CES for retention. **[F]**

## 4. Cross-cutting synthesis — which metric to trust

1. **Each measures a different thing on a different time axis** — CSAT (past, transactional), NPS
   (relationship, intention/future), CES (past, process effort). Complements, not substitutes. **[F]**
2. **Every "single best metric" claim is contested by the peer-reviewed literature.** Keiningham (2007)
   kills NPS's superiority claim; de Haan (2015) shows the best metric is industry-dependent and that
   top-2-box CSAT often wins for retention; CES's superiority claim is service-scoped and CEB-proprietary.
   The replicated consensus: **combine metrics; no single number reliably predicts growth/retention
   across contexts.** **[F]**
3. **All three share survey pathologies** — self-selection/non-response, survival bias (churned customers
   vanish), recency bias, and gaming when tied to incentives. **[F]**
4. **For a TAM:** use NPS as a relationship pulse and closed-loop trigger; use CSAT (top-box) and CES on
   specific touchpoints; never present one number as the predictor of account health — segment, look at
   the tails, and (§b2b) read the account as a multi-stakeholder buying center.

## Sources

1. Qualtrics, "What is CSAT?" — https://www.qualtrics.com/en-gb/articles/customer-experience/what-is-csat/ — [industry/docs]
2. SurveyMonkey, "CSAT Score Formula, Scale and Benefits" — https://www.surveymonkey.com/learn/customer-feedback/csat-score/ — [industry]
3. Gladly, "Customer satisfaction score (CSAT)" — https://www.gladly.ai/glossary/customer-satisfaction-score/ — [industry]
4. Zendesk, "What is customer satisfaction score (CSAT)" — https://www.zendesk.com/blog/customer-experience/loyalty/customer-loyalty/customer-satisfaction-score/ — [industry]
5. SI Labs, "Customer Satisfaction Score (CSAT)" (CSAT/NPS/CES comparison) — https://www.si-labs.com/en/articles/customer-satisfaction-score/ — [blog]
6. Gallup, "Customer Satisfaction: A Flawed Measure" (2007) — https://news.gallup.com/businessjournal/28564/Customer-Satisfaction-Flawed-Measure.aspx — [industry/research]
7. MIT Sloan Management Review, "The High Price of Customer Satisfaction" (Keiningham et al.) — https://sloanreview.mit.edu/article/the-high-price-of-customer-satisfaction/ — [paper/industry]
8. SupportLogic, "Limitations of CSAT and NPS" — https://www.supportlogic.com/resources/blog/your-customer-satisfaction-metrics-dont-show-the-full-picture/ — [blog/industry]
9. Reichheld, "The One Number You Need to Grow," HBR Dec 2003 — https://hbr.org/2003/12/the-one-number-you-need-to-grow — [paper]
10. Keiningham, Cooil, Andreassen & Aksoy (2007), *Journal of Marketing* 71(3):39–51 — https://journals.sagepub.com/doi/10.1509/jmkg.71.3.039 — [paper]
11. MeasuringU (Sauro), "Does the Net Promoter Score Predict Company Growth?" — https://measuringu.com/nps-growth/ — [blog/research]
12. Kristensen & Eskildsen, "Is the NPS a trustworthy performance measure?" *TQM Journal* — https://www.emerald.com/insight/content/doi/10.1108/TQM-03-2011-0021/full/html — [paper]
13. Fisher & Kordupleski (2019), *ASMBI* 35(1):138–151 — https://ideas.repec.org/a/wly/apsmbi/v35y2019i1p138-151.html (preprint arXiv:1806.10452 — https://arxiv.org/pdf/1806.10452) — [paper]
14. Stahlkopf, "Where Net Promoter Score Goes Wrong," HBR Oct 2019 (paywalled) — https://hbr.org/2019/10/where-net-promoter-score-goes-wrong — [paper]
15. Bain, "The Numbers behind the Net Promoter System" — https://www.netpromotersystem.com/about/numbers-behind-the-net-promoter-system/ — [industry]
16. Bain, "The Benefits of a Competitive Benchmark NPS" — https://www.bain.com/insights/the-benefits-of-a-competitive-benchmark-net-promoter-score/ — [industry]
17. Dixon, Freeman & Toman, "Stop Trying to Delight Your Customers," HBR Jul–Aug 2010 — https://hbr.org/2010/07/stop-trying-to-delight-your-customers — [paper]
18. *The Effortless Experience* (Dixon, Toman, DeLisi, 2013), Penguin Random House — https://www.penguinrandomhouse.com/books/312730/the-effortless-experience-by-matthew-dixon/ — [book]
19. de Haan, Verhoef & Wiesel (2015), *IJRM* 32(2):195–206 — https://www.sciencedirect.com/science/article/abs/pii/S0167811615000324 (open PDF: https://pure.rug.nl/ws/files/17052354/In_press_version_IJRM.pdf) — [paper]
20. Great Brook Consulting, CES 1.0/2.0 (quoting the book) — https://greatbrook.com/effortless-experience-questionnaire-design-survey-administration-issues/ — [blog]
21. eTouchPoint, "The Customer Effort Score: A Primer" — https://www.etouchpoint.com/customer-effort-score-primer-recommendations-cx-practitioners/ — [blog/industry]
22. Delighted, "What is Customer Effort Score (CES)?" — https://delighted.com/what-is-customer-effort-score — [industry]
23. MeasuringU (Sauro), "10 Things to Know about the Customer Effort Score" — https://measuringu.com/customer-effort-score/ — [blog/research]
