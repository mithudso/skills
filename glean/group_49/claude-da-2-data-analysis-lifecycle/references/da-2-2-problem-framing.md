<!-- hub-reference-banner -->
> **Reference file — part of the `da-2-data-analysis-lifecycle` hub.** Formerly the standalone `da-2-2-problem-framing` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-2-2-problem-framing
description: |
  Expert knowledge on problem framing as the opening phase of the data analysis lifecycle — the practice of translating a vague business challenge into a precise, scoped, and measurable analytical question before any data work begins.

  TRIGGER: Use when a user asks how to start a data analysis or data science project, how to define or scope an analysis problem, how to write a problem statement for analytics, how to align business goals with data questions, how to translate a business need into an analytical objective, or when they ask about the discovery phase of the analytics lifecycle. Also trigger when someone asks "what question should my analysis answer?" or "how do I frame this for data analysis?".

  SKIP: Skip when the question is about a specific sub-technique that has its own node in the taxonomy — such as writing a formal business/research question (da-2-2-1), hypothesis formulation (da-2-2-2), or specific process frameworks like CRISP-DM (da-2-1-1), KDD (da-2-1-2), SEMMA (da-2-1-3), OSEMN (da-2-1-4), or TDSP (da-2-1-5), which each handle problem framing as one step inside a full end-to-end process. Also skip for problem framing as a general product/UX design activity — this skill covers it specifically within the Data Analysis Lifecycle context.
---

# Problem Framing in the Data Analysis Lifecycle

## What it is

Problem framing is the opening phase of any data analysis project: the structured activity of converting an abstract, often vague business challenge into a specific, measurable, and actionable analytical question. Nothing else in the lifecycle — data collection, modeling, visualization, deployment — can succeed if this step is skipped or done carelessly.

The core translation: **business language → analytical language**.

> "One of the most overlooked skills in data science isn't coding or model tuning, it's the ability to translate a business question into a data question." — Medium/@Ann-MarieA [1]

> "Over 80% of data science initiatives fail, largely because teams skip right to analyzing the data before agreeing on the problem to be solved." — MIT Sloan Management Review [2]

---

## Why it comes first

In the data analysis lifecycle, problem framing precedes data collection, exploration, and modeling for several reasons:

1. **Prevents wasted effort.** Without a clear target question, analysts often collect irrelevant data or run "fishing expedition" analyses — trolling data for unexpected correlations with no guiding hypothesis [2].
2. **Prevents scope creep.** A written, stakeholder-agreed problem statement constrains the project to relevant issues. Scope creep (unexpected delivery requirements causing delays and cost increases) is one of the most common analytics project failures [3].
3. **Enables success evaluation.** You cannot judge whether an analysis succeeded without knowing what question it was supposed to answer and what "good enough" looks like.
4. **Aligns multiple stakeholders.** Marketing, sales, operations, and technical teams each have different assumptions about "the problem." Framing surfaces and resolves those differences before they derail the project mid-stream [2].

---

## The translation process: business problem → analytical problem

Business stakeholders use language like "improve revenue," "reduce churn," or "understand our customers." These phrases are not analytical questions — they are starting points that require decomposition.

The translation follows this direction:

```
Business objective
  → Root cause investigation (5 Whys, fishbone)
    → Specific question with measurable outcome
      → Constraints and scope boundary
        → Success criteria
          → Documented problem statement
```

**Example (bad framing):** "Build a report on all customers."  
**Example (good framing):** "Which customer profiles show the highest probability of churning after month one, and what early behaviors predict that?" [4]

The difference: the second version names a population, a measurable outcome (churn probability), and an implied decision (early intervention). It is answerable, bounded, and action-oriented.

---

## Key components of a well-framed problem

### 1. Business objective
What decision or action will this analysis support? Stakeholders need actionable decisions, not findings. Frame the deliverable around who needs to decide what, by when, and under what constraints [1].

### 2. Root cause identification
Use the **5 Whys** technique: ask "why does this problem exist?" repeatedly until you reach an underlying cause rather than a symptom. This prevents solving surface effects while leaving causes intact [3][4].

The **Fishbone (Ishikawa) diagram** is an alternative structured technique for mapping multiple contributing causes across categories (people, process, data, technology, environment).

### 3. Specific, measurable analytical question
The question must be:
- **Specific**: names the population, variable, and timeframe
- **Measurable**: has a quantifiable answer
- **Actionable**: the answer must feed a real decision
- **Relevant**: connected to the stated business objective
- **Time-bound**: has a defined scope in time or data coverage

This is the SMART criteria applied to analytical questions [3][5].

### 4. Constraints
Document all constraints upfront:
- **Data availability**: what data exists, in what form, covering what time period
- **Timeline**: when does the analysis need to be complete
- **Budget**: what compute, tooling, and personnel is available
- **Privacy and legal**: PII restrictions, regulatory requirements (GDPR, HIPAA, etc.)
- **Deployment**: where will the result live; latency requirements if a model will run in production [6]

### 5. Success criteria
Define what "done" and "good enough" look like before starting. Success criteria combine:
- **Business KPIs**: the metric the business cares about (revenue uplift, churn reduction rate, cost savings)
- **Analytical metrics**: the model or statistical evaluation metric that predicts the business KPI (accuracy, precision-recall, confidence interval width)

These two are distinct. A model can be 95% accurate and still fail the business criterion if it optimizes the wrong proxy [6].

### 6. Stakeholder map
Identify who:
- Owns the problem (decision authority)
- Will use the output (consumers of findings)
- Has domain knowledge needed for framing (subject matter experts)
- Must sign off before the project proceeds (gatekeepers)

Use a **RACI matrix** (Responsible, Accountable, Consulted, Informed) to document this [5].

---

## The C.L.A.R.I.F.Y. framework

A structured mnemonic for working through problem framing [4]:

| Letter | Step | What to do |
|--------|------|------------|
| **C** | Context | Map the business landscape: strategic priorities, affected stakeholders, whether this is a cost, revenue, risk, or experience problem |
| **L** | Look for symptoms | Identify visible issues (declining sales, high churn, low engagement) without conflating them with causes |
| **A** | Ask why | Apply 5 Whys to each symptom to reach underlying causes |
| **R** | Reframe into data questions | Convert business issues into specific, measurable analytical questions tied to real decisions |
| **I** | Identify stakeholders | Engage decision-makers early to clarify success metrics and how insights will influence choices |
| **F** | Filter the noise | Focus on the key metrics, trusted data sources, and relevant timeframes; avoid overanalyzing |
| **Y** | Yield to action | Verify the data question supports a specific decision with clear possible actions |

---

## Problem complexity categories

Not all analytics problems are framed the same way. A useful classification [5]:

| Type | Nature | Approach |
|------|--------|----------|
| **Simple** | Known complexities with clear discovery paths | Apply established best practices directly |
| **Complicated** | Known or unknown complexities but an identifiable path exists | Engage domain experts; decompose into sub-problems |
| **Complex** | Both problem nature and discovery path are uncertain | Iterative experimentation; accept that framing itself may evolve across multiple cycles |

Complex problems require the analyst to treat problem framing as iterative — the framing document is a living artifact revised as the data reveals surprises, not a one-time deliverable [1].

---

## Outputs of problem framing

A completed problem framing phase produces one or more of these artifacts:

**Problem Definition Document** containing:
- Business objective statement
- Analytical question(s)
- Success criteria (business KPIs + analytical metrics)
- Constraints (data, time, budget, legal)
- Stakeholder map / RACI
- Scope boundary (what is explicitly out of scope)
- Initial hypotheses to test

**Scope statement** — a brief, agreed description of what the project will and will not do.

**Question hierarchy** — decomposing the main question into sub-questions that can be answered independently and aggregated.

---

## Common pitfalls

### Starting with the solution
The most common failure: a stakeholder says "I need a dashboard showing X" or a data scientist says "let's try gradient boosting." Both skip problem framing. The solution type should be derived from the question, not chosen first [1][5].

### Vague problem statements
"Improve customer retention" is not a problem statement. Without a specific population, a measurable target, and a timeframe, there is no way to know whether the analysis answered the question [3].

### Stakeholder mismatch
When multiple stakeholders have different (unstated) definitions of the problem, the analysis will satisfy none of them. Run explicit alignment sessions before framing is considered complete [2].

### Confusing symptoms with causes
Treating a data quality problem as the problem to solve (e.g., cleaning bad records) rather than asking why the data is bad — and whether it will recur — is a common framing error [2].

### Overfitting the frame to available data
Analysts sometimes unconsciously narrow the question to what their current dataset can answer rather than what the business actually needs to know. This produces technically correct but strategically irrelevant analysis [6].

### Letting stakeholders prescribe solutions
Stakeholders should define the problem; analysts define the solution approach. When stakeholders specify "build a regression model" rather than stating the decision they need to make, the frame is wrong [4].

### Proxy label confusion
When the direct outcome cannot be measured (e.g., "customer satisfaction"), analysts use proxy signals (e.g., repeat purchases, NPS scores). Framing must explicitly document which proxy is being used and acknowledge its limitations [6].

---

## Relationship to the broader data analysis lifecycle

Problem framing is the first phase; it gates everything downstream:

```
[Problem framing]  ←— this skill
    ↓
Business/research question definition  (da-2-2-1)
Hypothesis formulation  (da-2-2-2)
    ↓
Data collection and understanding
    ↓
Exploratory data analysis
    ↓
Modeling / analysis
    ↓
Evaluation against success criteria
    ↓
Communication and deployment
```

The sub-skills `da-2-2-1` (business/research question definition) and `da-2-2-2` (hypothesis formulation) cover specific output artifacts produced during problem framing. Process frameworks such as CRISP-DM (`da-2-1-1`), KDD (`da-2-1-2`), SEMMA (`da-2-1-3`), OSEMN (`da-2-1-4`), and TDSP (`da-2-1-5`) each embed problem framing as one named phase inside a full end-to-end process — they do not replace this standalone treatment of the framing activity.

---

## Worked example

**Scenario**: A subscription SaaS company is losing customers.

| Stage | Content |
|-------|---------|
| **Business objective** | Reduce monthly churn to below 3% within two quarters |
| **Symptoms identified** | Month-over-month churn rate has increased from 2.1% to 4.7% over six months |
| **5 Whys result** | Customers who onboard without completing a key setup step churn at 3× the rate of those who complete it |
| **Analytical question** | "What is the probability that a given new customer will churn within 90 days, based on their first-14-day product activity?" |
| **Success criteria** | Business: reduce 90-day churn by 30%. Analytical: model achieves AUC > 0.75 on held-out data |
| **Constraints** | 18 months of event log data available; no PII in model features; must run in real-time scoring pipeline with <200ms latency |
| **Out of scope** | Pricing changes, product feature changes, marketing campaign design |

The resulting problem definition document gives the data team a clear, evaluable target before they touch a single row of data.

---

## Sources

1. Ann-Marie A., "From Business Problem to Data Problem," Medium. https://medium.com/@Ann-MarieA/from-business-problem-to-data-problem-99fee59c512e

2. Ransbotham, S. et al., "Framing Data Science Problems the Right Way From the Start," MIT Sloan Management Review. https://sloanreview.mit.edu/article/framing-data-science-problems-the-right-way-from-the-start/

3. Render Analytics, "Why It is Important to Properly Frame and Scope a Data Analytics Project." https://www.renderanalytics.net/post/why-it-is-important-to-properly-frame-and-scope-a-data-analytics-project

4. Diab, Z., "The Art of Data Problem Framing," Medium. https://medium.com/@ziadkdiab/the-art-of-data-problem-framing-7d6881e34bb4

5. Product Leadership, "Business Problem Framing in the Age of Data Analytics." https://www.productleadership.com/blog/business-problem-framing-in-the-age-of-data-analytics/

6. Google for Developers, "Framing an ML Problem," Machine Learning Crash Course. https://developers.google.com/machine-learning/problem-framing/ml-framing
