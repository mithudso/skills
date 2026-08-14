<!-- hub-reference-banner -->
> **Reference file — part of the `da-2-data-analysis-lifecycle` hub.** Formerly the standalone `da-2-2-1-business-research-question-definition` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-2-2-1-business-research-question-definition
description: |
  Expert knowledge for defining and formulating business and research questions at the start of a data analysis project, within the Problem Framing stage of the Data Analysis Lifecycle. Covers how to translate a business objective into a precise, answerable analytical question — the critical first act before selecting methods, collecting data, or building models.

  TRIGGER: Use this skill when the user is working at the problem-scoping stage of a data analysis project and needs help with any of: translating a business goal into an analytical question, writing or critiquing a problem statement, applying SMART criteria to a data question, aligning stakeholders on what question to answer, or distinguishing the business question from the data question. Also trigger when phrases like "define the problem", "frame the question", "what should we be trying to answer", "business objective → analytical goal", or "scope the analysis" appear in a data analysis context.

  SKIP: Do not use for hypothesis formulation once the question is already defined (that belongs to da-2-2-2-hypothesis-formulation), for overall process framework selection (CRISP-DM, KDD, SEMMA — see da-2-1-x skills), for choosing statistical or modeling methods, or for question formulation in academic/clinical research contexts (e.g., PICO for systematic reviews) unless the user is explicitly mapping those frameworks onto a business analytics workflow.
---

# Business/Research Question Definition

## Position in the taxonomy

**Data Analysis > Data Analysis Lifecycle (Process) > Problem framing > Business/research question definition**

This node covers the act of producing a well-formed question from a raw business situation. It is the first concrete deliverable of problem framing. Adjacent nodes handle what comes after: hypothesis formulation, success-metric definition, and stakeholder alignment. Those are separate steps; this skill focuses on getting the question right.

---

## What is a business/research question in data analysis?

A **business question** is a precise statement of what decision or outcome the organization wants to improve, expressed in non-technical language tied to business metrics (revenue, churn, cost, lead time, satisfaction scores).

A **research question** (or **analytical question**) is the translation of that business question into a form that can be answered with data — specifying what to measure, in what population, over what time window, and under what conditions.

These two formulations are not the same thing. Moving between them is the central skill here.

> "A business goal states objectives in business terminology, while a data mining goal states project objectives in technical terms." — CRISP-DM specification, via KDnuggets [1]

The gap between the two is where most analytics projects fail. Teams that skip this translation step optimize analytical accuracy while solving the wrong problem.

---

## Why this step matters

- **Scoping**: a sharp question tells you what data to collect and what to ignore. Without it, analysts pull everything and drown.
- **Alignment**: stakeholders often disagree about what the real problem is. Writing a question forces that disagreement to surface early.
- **Evaluation criteria**: you cannot know when you are done unless you know what the question was.
- **Preventing solution-first thinking**: stakeholders frequently arrive with a pre-chosen solution ("build a dashboard") rather than a problem. A structured question-definition process redirects to outcomes.

A Harvard study cited by ProductLeadership.com found that 85% of CXOs and senior executives reported their organizations struggle with problem diagnosis — the upstream failure that drives wasted analytical work [2].

---

## The translation process: from business objective to analytical question

### Step 1 — Capture the raw business objective

Start with what the business *wants*, not what it *thinks the data should do*. Typical raw objectives:

- "Increase sales"
- "Reduce customer churn"
- "Improve operational efficiency in the warehouse"

These are too vague to analyze directly. They need decomposition.

### Step 2 — Identify the decision that analytics must support

Analytics does not improve metrics directly — it improves *decisions*, which then move metrics. Map the raw objective to a repeatable decision someone actually makes.

> "Analytics cannot improve metrics directly...unless you change your behavior – unless you do something different – your churn rate is not going to improve." — KDnuggets, *Bringing Business Clarity to CRISP-DM* [1]

Questions to ask here:
- Who makes the decision this analysis will inform?
- How often do they make it?
- What information do they currently lack?
- What would they do differently with a better answer?

### Step 3 — Apply the SMART criteria

Refine the candidate question against five dimensions [2][3]:

| Dimension | Question to ask |
|---|---|
| **Specific** | Does it name the exact phenomenon, population, and metric? |
| **Measurable** | Is there a quantifiable outcome we can compare against a baseline? |
| **Achievable** | Is there plausibly data available to answer it? |
| **Relevant** | Does it connect to a decision that moves a business metric? |
| **Time-bound** | Does it specify a time window for the analysis? |

**Before SMART**: "Should we do something about customer churn?"

**After SMART**: "Which customer profiles show the highest 90-day churn probability based on product usage in their first 30 days, and what is the earliest behavioral signal?"

### Step 4 — Separate the business question from the analytical question

Write both down explicitly. The CRISP-DM Business Understanding phase formalizes this as two required outputs [1][4]:

- **Business objective**: "Reduce monthly subscription cancellations by 10% within two quarters."
- **Data mining / analytical goal**: "Build a classifier that identifies subscribers in their first 60 days with >70% probability of cancelling before day 90, using product engagement events as features."

If you cannot write the analytical goal without the business objective visible, the translation is not complete.

### Step 5 — Test against the decision-action link

Ask: *"If we had a perfect answer to this question, what would we do differently tomorrow?"*

If the answer is "nothing" or "we'd make a report," the question is not yet decision-linked. Rework until the answer names a concrete action.

---

## The C.L.A.R.I.F.Y. framework

A practitioner-level mnemonic for structuring question definition sessions (source: Ziad Diab, Medium [5]):

| Letter | Step |
|---|---|
| **C** — Context | Understand the broader business situation and strategic priorities before touching data |
| **L** — Look for Symptoms | Identify visible issues; do not confuse symptoms with root causes |
| **A** — Ask Why | Use 5-Whys or similar to uncover underlying causes of symptoms |
| **R** — Reframe into Data Questions | Transform business issues into specific, measurable analytical questions |
| **I** — Identify Stakeholders | Engage decision-makers early; understand their success criteria |
| **F** — Filter the Noise | Focus on key metrics and trusted sources; avoid over-scoping |
| **Y** — Yield to Action | Verify the question supports a concrete decision or business action |

---

## Worked example

**Starting point (stakeholder statement):** "We're struggling to hit our daily order fulfillment goal."

**Step 1 — Objective extracted:** Improve daily fulfillment rate.

**Step 2 — Decision identified:** Warehouse floor managers decide labor allocation by shift. They currently rely on yesterday's order count. The decision is made twice daily.

**Step 3 — SMART applied:**
- *Specific*: Predict order volume by shift for each of our three warehouse zones.
- *Measurable*: Error < 30% mean absolute percentage error vs. actual volume.
- *Achievable*: Historical order data by zone exists for 18 months.
- *Relevant*: A 30% MAPE improvement is projected to reduce labor overstaffing costs by 20%.
- *Time-bound*: Predict the next shift (4-hour horizon).

**Business question:** "Can we predict short-term order volume per warehouse zone accurately enough to improve labor scheduling?"

**Analytical question:** "Can a time-series model trained on 18 months of zone-level order history achieve MAPE < 30% on 4-hour-ahead volume forecasts for each of three zones?"

*(Adapted from Towards Data Science [3])*

---

## Common pitfalls

### Accepting a solution disguised as a question
Stakeholders often say "We need a dashboard that shows X." That is not a question — it is a pre-selected deliverable. Push back to: "What decision does that dashboard need to support, and how often?"

### Confusing output with outcome
"Build a model" or "make a report" are outputs. Outcomes are decisions made and metrics moved. Always close the loop: output → decision → metric.

### Over-specifying technical details too early
A good business question should be readable by a non-technical stakeholder. If it contains model names, SQL, or algorithm references, you have skipped to method selection.

### Under-specifying the population or time window
"Why do customers churn?" is not answerable as stated. "What behavioral signals in the first 30 days of subscription predict cancellation within 90 days, among B2B customers acquired via paid search?" is.

### Assuming shared understanding
The same words mean different things to different stakeholders. "Customer" might mean account, contact, or seat depending on who you ask. Define every term in the question explicitly.

### Treating problem definition as a one-time activity
Business conditions evolve. Questions that were correct at project start may need revision at midpoint. Build in a review checkpoint before the modeling phase, not just at kickoff [4].

### Confirmation bias in scoping
Anchoring on a hypothesis before the question is fully defined narrows which data gets collected and which alternatives get considered. Question definition should precede hypothesis formulation.

---

## Relationship to adjacent concepts

| Concept | Relationship |
|---|---|
| **Problem framing (da-2-2)** | Parent concept; question definition is its first concrete output |
| **Hypothesis formulation (da-2-2-2)** | Comes after the question is defined; turns the question into a testable proposition |
| **CRISP-DM Business Understanding (da-2-1-1)** | CRISP-DM is a process framework; its Business Understanding phase is where question definition happens within that framework |
| **Success criteria / metric definition** | Downstream; requires an established question before metrics can be specified |
| **SMART questions (Google Data Analytics)** | A widely taught application of SMART to data question formulation [6] |

---

## Quick reference checklist

Before moving from question definition to hypothesis formulation or data collection, verify:

- [ ] The business question is written in non-technical language a stakeholder can confirm.
- [ ] The analytical question specifies: population, metric, time window, and comparison point.
- [ ] Both questions are written down separately and reviewed by at least one stakeholder.
- [ ] The decision this question informs is named, including who makes it and how often.
- [ ] The question passes all five SMART criteria.
- [ ] You can answer "what would we do differently if we had a perfect answer?"
- [ ] All key terms in the question are defined (no assumed shared meaning).

---

## Sources

1. Groth, R. (2017). "Bringing Business Clarity to CRISP-DM." *KDnuggets*. https://www.kdnuggets.com/2017/01/business-clarity-crisp-dm.html

2. ProductLeadership.com. "Business Problem Framing in the Age of Data Analytics." https://www.productleadership.com/blog/business-problem-framing-in-the-age-of-data-analytics/

3. Towards Data Science. "Essential Questions to Ask Before Starting a Data Science Project." https://towardsdatascience.com/essential-questions-to-ask-before-starting-a-data-science-project-cd633dcd9d55/

4. Splendor, E. U. (Medium). "Defining Business Problems for Analytic Thinking & Data-Driven Decision Making." https://medium.com/@Splendor001/defining-business-problems-for-analytic-thinking-data-driven-decision-making-efb1b166d8cc

5. Diab, Z. (Medium). "The Art of Data Problem Framing." https://medium.com/@ziadkdiab/the-art-of-data-problem-framing-7d6881e34bb4

6. TrustEd Institute. "SMART Questions Methodology — Google Data Analytics Certificate." https://trustedinstitute.com/concept/google-data-analytics/ask-questions-data-driven-decisions/smart-questions-methodology/
