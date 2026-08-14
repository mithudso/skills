<!-- hub-reference-banner -->
> **Reference file — part of the `da-3-data-acquisition-sampling` hub.** Formerly the standalone `da-3-2-collection-methods` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-3-2-collection-methods
version: "1.1.0"
updated: "2026-05-30"
description: |
  Categorical overview and comparison of data collection methods across the full taxonomy:
  surveys/questionnaires, interviews/focus groups, controlled experiments, direct
  observation, instrumentation and telemetry, web scraping, APIs, public/secondary
  datasets, sensors/IoT, and transactional logs. Covers comparison dimensions (cost,
  researcher control, internal/external validity, scale, bias risk, ethical obligations)
  and a decision framework for method selection.

  TRIGGER: User asks "which data collection method should I use", "how do I collect data
  for X", "compare surveys vs experiments", "what are the types of data collection",
  "tradeoffs between scraping and APIs", "when to use observational vs experimental data",
  "primary vs secondary data collection overview", "how do interviews compare to surveys",
  "when should I run a focus group", or any question about choosing among collection
  approaches at the categorical level.

  SKIP: Detailed execution of a single method — defer to sibling skills (planned):
  da-3-2-1-surveys-questionnaires (survey design, question writing, Likert scales,
    response rate tactics, survey platforms),
  da-3-2-2-experiments (RCT design, power analysis, A/B test mechanics, pre-registration),
  da-3-2-3-observation (structured vs. ethnographic observation, field notes,
    inter-rater reliability),
  da-3-2-4-instrumentation-telemetry (event schemas, log pipelines, APM tools),
  da-3-2-5-apis-web-scraping (REST pagination, authentication, HTML parsing,
    rate limiting, legal compliance),
  da-3-2-6-public-datasets (census/data.gov/Kaggle sourcing, license review,
    provenance documentation),
  da-3-2-7-sensors-iot (device protocols, MQTT, time-series ingestion, calibration).
  Also skip if the question is about sampling strategy (da-3-data-acquisition-sampling),
  data cleaning after collection (da-4-data-cleaning-preparation), or deep ethics
  analysis (da-11-ethics-and-privacy).
related_skills:
  - da-3-data-collection-acquisition
  - da-3-1-1-primary-vs-secondary
  - da-3-data-acquisition-sampling
  - da-4-data-cleaning-preparation
  - da-11-ethics-and-privacy
---

# Data Collection Methods: Categorical Overview

## What This Covers

Data collection is the systematic process of gathering raw information to answer a research question, train a model, or support an operational decision. The field divides into ten major method families. Choosing among them before touching any tool determines data quality, cost, and what analytical claims are defensible downstream.

This skill covers the categorical taxonomy and selection logic. For deep execution of any single method, see the SKIP list in the frontmatter.

---

## The Ten Method Families

> **Glossary of terms used throughout:** *Internal validity* — the degree to which a study establishes a cause-and-effect relationship (high in experiments, lower in observational studies). *External validity* — the degree to which findings generalize beyond the study context (high when collection occurs in natural settings). *Ecological validity* — a form of external validity; whether study conditions reflect real-world conditions. *Hawthorne effect* — participants alter their behavior because they know they are being observed.

### 1. Surveys and Questionnaires

Structured instruments delivered to a sample, typically mixing closed-ended (scaled, categorical) and open-ended questions ([Excelsior College Data Science](https://express.excelsior.edu/datascience/chapter/chapter-2-6-data-collection-methods-surveys-experiments-and-observations/)). They are the primary route to internal states — motivations, opinions, self-reported behavior — that cannot be directly observed.

**When to use:** You need breadth across a population; you are measuring attitudes or preferences; budget is moderate; real-time observation is impractical.

**Bias risks:** Non-response bias (vocal or available respondents over-represent themselves), social desirability bias (respondents give acceptable rather than true answers), recall error on past behavior, and wording effects from question framing ([Sopact](https://www.sopact.com/use-case/data-collection-methods)).

**Ethical obligation:** Informed consent, data anonymization, and clear statement of purpose at collection time.

---

### 2. Controlled Experiments and A/B Tests

Systematic manipulation of one or more independent variables under conditions designed to hold confounders constant, then measurement of outcomes in treatment vs. control groups ([Excelsior College Data Science](https://express.excelsior.edu/datascience/chapter/chapter-2-6-data-collection-methods-surveys-experiments-and-observations/)). Random assignment is the mechanism that makes causal inference defensible.

**When to use:** The research question is causal ("does X cause Y?"); an intervention is testable; resources support controlled conditions; the setting is a product, clinical, or policy context.

**Bias risks:** Artificial conditions reduce ecological validity; demand characteristics (participants behave differently when observed); Hawthorne effect; ethical limits on withholding potentially beneficial treatments ([Sopact](https://www.sopact.com/use-case/data-collection-methods)).

**Ethical obligation:** Informed consent for human subjects; institutional review board (IRB) oversight for medical or sensitive contexts; pre-registration of hypotheses to prevent p-hacking.

---

### 3. Direct Observation

Systematic recording of behaviors or events as they occur in natural or structured settings, without researcher manipulation ([Sopact](https://www.sopact.com/use-case/data-collection-methods)). Observation captures what people actually do versus what they say they do.

**When to use:** Behavioral authenticity matters; participants cannot articulate the behavior verbally; longitudinal tracking of natural change is needed.

**Bias risks:** Observer effect (participants alter behavior when watched); observer bias (researcher interprets ambiguous acts through their own frame); selection of which moments to record ([Excelsior College Data Science](https://express.excelsior.edu/datascience/chapter/chapter-2-6-data-collection-methods-surveys-experiments-and-observations/)).

**Ethical obligation:** Consent varies by setting — covert observation raises significant ethical concerns; digital behavioral data requires disclosure in privacy policies.

---

### 4. Instrumentation and Telemetry

Software-embedded event tracking that records user or system actions automatically as they occur — page views, button clicks, errors, latency, resource utilization ([Pingax](https://pingax.com/data-analytics/stage-2-data-collection/data-sources-and-techniques/data-collection-methods-apis-web-scraping-and-more/)). This is the high-scale, low-marginal-cost end of the observational spectrum.

**When to use:** You operate the software or device generating events; you need large-N behavioral data; real-time monitoring or alerting is required.

**Bias risks:** Instrumentation only captures what was instrumented — missing events create silent gaps; survivorship bias if inactive users churn before data is collected; logging errors corrupt the record.

**Ethical obligation:** Users must be notified of telemetry in privacy policy; PII must not appear in event payloads without explicit consent; data minimization principles apply.

---

### 5. APIs

Standardized programmatic interfaces that expose structured data from a service or platform to authorized callers ([Pingax](https://pingax.com/data-analytics/stage-2-data-collection/data-sources-and-techniques/data-collection-methods-apis-web-scraping-and-more/)). Output is typically JSON or XML. Examples include social media APIs, financial data feeds, weather services, and internal microservice endpoints.

**When to use:** The data owner provides a sanctioned access mechanism; data is structured and updated continuously; rate limits are acceptable for your volume needs.

**Bias risks:** The platform controls what data is surfaced — historical access, filtering, and field availability change at the provider's discretion; rate limits create temporal gaps.

**Ethical obligation:** Comply with terms of service; respect rate limits; store credentials securely; review data license for redistribution restrictions.

---

### 6. Web Scraping

Programmatic extraction of information from web pages by parsing HTML structure ([Pingax](https://pingax.com/data-analytics/stage-2-data-collection/data-sources-and-techniques/data-collection-methods-apis-web-scraping-and-more/)). Tools include BeautifulSoup, Playwright, and Selenium for JavaScript-rendered content.

**When to use:** No API exists for the target data; the data is publicly displayed; you have legal clarity on the site's terms of service.

**Bias risks:** Selection bias — only data the site chooses to display is captured; layout changes break parsers silently; sites present different content to different users (A/B tests, geo-targeting).

**Ethical obligation:** Review `robots.txt` and terms of service before collecting; do not overwhelm servers (rate limiting is required); legal status varies by jurisdiction and depends on whether the data is public vs. access-controlled ([Pingax](https://pingax.com/data-analytics/stage-2-data-collection/data-sources-and-techniques/data-collection-methods-apis-web-scraping-and-more/)).

---

### 7. Public and Secondary Datasets

Data already collected by another party — government statistical agencies (Census, BLS, CDC), academic repositories, open data portals, licensed commercial datasets, or previously published research ([Sopact](https://www.sopact.com/use-case/data-collection-methods)).

**When to use:** Speed is critical; population-level baselines are needed; budget is low; the research question aligns with an existing collection's original purpose.

**Bias risks:** Researcher disconnect from collection context; secondary data is typically published 2–3 years after collection, reducing recency; the original sampling frame may not match your target population ([Sopact](https://www.sopact.com/use-case/data-collection-methods)).

**Ethical obligation:** Cite the original data source; comply with redistribution license (many government datasets are public domain; many commercial datasets are not); respect any de-identification or aggregation constraints.

---

### 8. Sensors and IoT

Physical devices — temperature probes, accelerometers, GPS trackers, medical monitors, industrial gauges — that convert physical phenomena into digital signals at high frequency ([Pingax](https://pingax.com/data-analytics/stage-2-data-collection/data-sources-and-techniques/data-collection-methods-apis-web-scraping-and-more/)). Infrastructure typically involves message brokers (Apache Kafka, MQTT) and cloud ingestion pipelines.

**When to use:** The phenomenon is physical and continuous; real-time monitoring or alerting is required; environmental, industrial, or health measurement is the goal.

**Bias risks:** Sensor calibration drift introduces systematic error over time; placement of sensors determines what is measured and what is not; device failure creates gaps that may not be flagged.

**Ethical obligation:** In consumer or medical IoT, explicit consent for continuous passive collection is required; data minimization and secure transmission (TLS) are baseline requirements.

---

### 9. Transactional and Operational Logs

Records generated automatically by business systems as side-effects of operations — purchases, logins, support tickets, database writes, application error logs ([Pingax](https://pingax.com/data-analytics/stage-2-data-collection/data-sources-and-techniques/data-collection-methods-apis-web-scraping-and-more/)). Unlike purpose-built telemetry, these exist primarily to support the system's function, not data collection.

**When to use:** Behavioral ground truth is needed without burdening users; operational or compliance auditing; longitudinal patterns in business activity.

**Bias risks:** Logs reflect system behavior, not user intent; events that never complete (abandoned carts, crashed sessions) may not be recorded; schema changes silently alter what is captured.

**Ethical obligation:** Log retention policies must comply with GDPR, CCPA, and equivalent regulations; PII in logs (IP addresses, email) requires masking or deletion schedules.

---

### 10. Interviews and Focus Groups

One-on-one conversations (interviews) or facilitated small-group discussions (focus groups, typically 6–10 participants) used to explore reasoning, lived experience, and group dynamics ([Sopact](https://www.sopact.com/use-case/data-collection-methods)). Three interview structures exist: structured (fixed questions), semi-structured (guided but flexible), and unstructured (open exploration).

**When to use:** You need to understand *why* behind behavior, not just *what*; you are scoping a problem before quantitative work; you need edge cases or outlier explanations that surveys cannot surface.

**Bias risks:** Interviewer effect (respondent answers differently based on perceived interviewer expectations); groupthink in focus groups (dominant voices suppress minority views); social desirability bias; small N limits generalizability.

**Ethical obligation:** Informed consent; audio/video recording requires explicit permission; confidentiality of participant identity must be maintained in reporting.

---

## Comparison Matrix

| Method | Relative Cost | Researcher Control | Internal Validity | External Validity | Scale | Primary Bias Risks |
|---|---|---|---|---|---|---|
| Surveys | Low–Moderate | High (question design) | Moderate | Moderate | High | Response bias, social desirability, recall error |
| Experiments | High | Highest | Highest | May be limited (artificial conditions) | Low–Moderate | Demand characteristics, Hawthorne effect |
| Observation | Low–Moderate | Low (no manipulation) | Lower | High (natural context) | Moderate (digital: High) | Observer bias, observer effect |
| Instrumentation/Telemetry | Low (ongoing) | Moderate (what gets instrumented) | High (behavioral ground truth) | High | Very High | Instrumentation gaps, survivorship bias |
| APIs | Low–Moderate | Low (provider controls schema) | High (structured) | Source-dependent | High (rate-limited) | Provider filtering, terms changes |
| Web Scraping | Moderate | Low | Variable (layout fragility) | Moderate | Moderate–High | Selection bias, fragile parsers |
| Public/Secondary Datasets | Very Low | None (pre-existing) | Varies | Moderate (lag 2–3 years) | High | Population mismatch, recency |
| Sensors/IoT | High (infrastructure) | High (placement/config) | Very High | High | Very High (continuous) | Calibration drift, placement bias |
| Transactional Logs | Very Low (already exist) | Low | High (operational ground truth) | Moderate | High | Incomplete events, schema drift |
| Interviews / Focus Groups | Moderate–High (time) | High (question design) | Moderate | Moderate (small N) | Low | Interviewer effect, groupthink, social desirability |

Sources: ([Excelsior College Data Science](https://express.excelsior.edu/datascience/chapter/chapter-2-6-data-collection-methods-surveys-experiments-and-observations/); [Pingax](https://pingax.com/data-analytics/stage-2-data-collection/data-sources-and-techniques/data-collection-methods-apis-web-scraping-and-more/); [Sopact](https://www.sopact.com/use-case/data-collection-methods))

---

## Decision Framework: Work Backward from the Decision

Rather than picking a method first, name the decision that requires data, then ask what answer would change the decision, then choose the method that produces that answer most efficiently ([Sopact](https://www.sopact.com/use-case/data-collection-methods)).

| Decision Need | Lead Method | Complement |
|---|---|---|
| How widespread is a problem or attitude? | Survey | Secondary data for baseline; interviews for outlier explanation |
| Did an intervention cause an outcome? | Controlled experiment / A/B test | Survey of treatment/control + interviews for mechanism |
| What do users actually do (not say they do)? | Instrumentation/telemetry or observation | Survey to compare stated vs. actual behavior |
| What is happening in a physical environment? | Sensors/IoT | Transactional logs for operational context |
| What is the population-level baseline? | Public/secondary datasets | Brief survey to validate applicability |
| What is available without new collection cost? | Transactional logs or secondary data | Observation or survey to fill gaps |
| What data exists on a third-party platform? | APIs (sanctioned) or web scraping (if no API) | Verify terms of service before proceeding |

The Excelsior / Spotify integration model makes the sequencing concrete: surveys identify needs → experiments test solutions → telemetry monitors real-world effects at scale ([Excelsior College Data Science](https://express.excelsior.edu/datascience/chapter/chapter-2-6-data-collection-methods-surveys-experiments-and-observations/)).

---

## Cross-Cutting Principles

**Mixed methods reduce single-method blind spots.** Surveys capture stated intent; telemetry captures actual behavior. Neither alone is sufficient for most product or policy decisions ([Sopact](https://www.sopact.com/use-case/data-collection-methods)).

**Persistent identifiers across methods.** When combining methods (e.g., survey + log + experiment), a shared participant or user ID makes integration tractable. Without it, joining datasets requires probabilistic matching, which degrades validity ([Sopact](https://www.sopact.com/use-case/data-collection-methods)).

**Collection method determines what inference is possible.** Non-experimental methods — surveys, interviews, telemetry, scraping, logs, and observation — support associative and descriptive claims. Only controlled experiments with random assignment support causal claims. This cannot be corrected at the analysis stage ([Excelsior College Data Science](https://express.excelsior.edu/datascience/chapter/chapter-2-6-data-collection-methods-surveys-experiments-and-observations/)).

**Ethics are not optional post-processing.** Consent, disclosure, and data minimization must be built into collection design, not retrofitted. GDPR and CCPA compliance is a collection-phase obligation, not a storage-phase one.

**Bias compounds across the pipeline.** A convenience sample in collection produces a biased population estimate regardless of how sophisticated the downstream analysis is. Fix bias at the source.

---

## Common Pitfalls

- Choosing a method before articulating the decision it must inform — leads to data that is interesting but not actionable.
- Using surveys to measure behavior when instrumentation data exists and is cheaper and more accurate.
- Scraping without reviewing `robots.txt` or terms of service — creates legal and reputational risk.
- Treating secondary data as current when the lag is 2–3 years and the environment has shifted.
- Assuming that more data volume compensates for a flawed collection design.
- Conflating correlation in observational data with causation — requires an experimental design with random assignment.

---

## Scope Boundary

This skill covers the categorical choice among methods. For execution details on any individual method, see the SKIP clause in the frontmatter above — each sibling skill is listed there with its precise coverage scope.

---

## Sources

1. [Data Collection Methods: Surveys, Experiments, and Observations — Excelsior College Introduction to Data Science](https://express.excelsior.edu/datascience/chapter/chapter-2-6-data-collection-methods-surveys-experiments-and-observations/) — Textbook chapter covering the three classical primary methods with comparison tables and the Spotify integration model.

2. [Data Collection Methods: APIs, Web Scraping, and More — Pingax](https://pingax.com/data-analytics/stage-2-data-collection/data-sources-and-techniques/data-collection-methods-apis-web-scraping-and-more/) — Practitioner guide covering APIs, scraping, sensors, surveys, and transactional logs with cost/bias/ethics comparison.

3. [Data Collection Methods: The Complete Guide — Sopact](https://www.sopact.com/use-case/data-collection-methods) — Comprehensive overview of seven primary methods, comparison dimensions, decision framework, and AI-native collection principles including persistent identifiers and mixed-method integration.
