# Customer Satisfaction (CSAT) Survey Template

## Survey Goal
Measure customer satisfaction with a specific interaction, product, or service experience.

---

## Questions

### Q1 — Core CSAT Question (Required)
**Type:** Likert Scale (1–5)

> How satisfied are you with [interaction/product/service]?

| Score | Label               |
|-------|---------------------|
| 1     | Very Dissatisfied   |
| 2     | Dissatisfied        |
| 3     | Neutral             |
| 4     | Satisfied           |
| 5     | Very Satisfied      |

---

### Q2 — Ease of Experience (Optional)
**Type:** Likert Scale (1–5)

> How easy was it to [complete your task / get help / use the product]?

| Score | Label        |
|-------|--------------|
| 1     | Very Difficult |
| 2     | Difficult      |
| 3     | Neutral        |
| 4     | Easy           |
| 5     | Very Easy      |

---

### Q3 — What Went Well (Optional)
**Type:** Open Text

> What did we do well?

---

### Q4 — What Could Be Improved (Optional)
**Type:** Open Text

> What could we have done better?

---

### Q5 — Contact Follow-Up Consent (Optional)
**Type:** Multiple Choice (Single Select)

> Would you like a team member to follow up with you about your experience?

- Yes, please contact me
- No, thank you

---

## Scoring Logic

**CSAT Score Formula:**
```
CSAT (%) = (Number of satisfied responses [4–5] ÷ Total responses) × 100
```

| CSAT %    | Interpretation          |
|-----------|-------------------------|
| 80–100%   | Excellent               |
| 60–79%    | Good                    |
| 40–59%    | Needs Improvement       |
| Below 40% | Critical — Take Action  |

---

## Display Logic (Recommended)

- Show Q3 and Q4 only after Q1 is answered.
- Show Q5 (follow-up consent) only for scores of 1 or 2 (dissatisfied responses).

---

## End of Survey Message

> Thank you for taking the time to share your feedback. We review every response and use it to improve your experience.
