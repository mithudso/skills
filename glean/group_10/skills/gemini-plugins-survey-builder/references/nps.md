# Net Promoter Score (NPS) Survey Template

## Survey Goal
Measure customer loyalty and likelihood to recommend the product or service.

---

## Questions

### Q1 — Core NPS Question (Required)
**Type:** Numeric Scale (0–10)

> On a scale of 0 to 10, how likely are you to recommend [Company/Product Name] to a friend or colleague?

- 0 = Not at all likely
- 10 = Extremely likely

---

### Q2 — Follow-Up: Reason for Score (Required)
**Type:** Open Text

> What is the primary reason for your score?

---

### Q3 — Improvement Opportunity (Optional)
**Type:** Open Text

> What could we do to improve your experience?

---

### Q4 — Overall Experience (Optional)
**Type:** Multiple Choice (Single Select)

> How would you describe your overall experience with [Company/Product Name]?

- Excellent
- Good
- Neutral
- Poor
- Very Poor

---

## Scoring Logic

| Score Range | Segment     | Description                                      |
|-------------|-------------|--------------------------------------------------|
| 9–10        | Promoters   | Loyal enthusiasts who will refer others          |
| 7–8         | Passives    | Satisfied but unenthusiastic, vulnerable to churn|
| 0–6         | Detractors  | Unhappy customers who may damage brand reputation|

**NPS Formula:**
```
NPS = % Promoters − % Detractors
```

Score ranges from **−100** (all detractors) to **+100** (all promoters).

---

## Display Logic (Recommended)

- Show Q3 only to **Detractors** (score 0–6) and **Passives** (score 7–8).
- Show a thank-you message to **Promoters** (score 9–10) after Q2.

---

## End of Survey Message

> Thank you for your feedback! Your response helps us improve [Company/Product Name].
