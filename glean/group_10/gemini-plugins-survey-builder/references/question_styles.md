# Survey Design Best Practices
## Exhaustive Guide for Psychologists & Quant UX Researchers

This guide outlines the gold standards for survey methodology, focusing on construct validity, bias reduction, and advanced choice architecture.

---

## 1. Construct Validity & Question Wording

### Avoid Double-Barreled Questions
A double-barreled question touches upon more than one issue yet allows for only one answer. This creates "noisy" data because you cannot tell which part of the question the respondent was answering.
* **❌ Bad:** "Was the onboarding process easy and enjoyable?" (What if it was easy but boring? Or enjoyable but difficult?)
* **✅ Good:** Use one construct per question.
    * Question 1: "How easy was it to complete the onboarding?" (Construct: Effort/Ease)
    * Question 2: "How enjoyable was the onboarding experience?" (Construct: Satisfaction/Delight)

### Avoid Leading & Loaded Questions
Leading questions use biased language that nudges a respondent toward a specific answer, while loaded questions contain an implicit assumption.
* **❌ Bad:** "How much do you like the new streamlined navigation?" (Assumes the navigation *is* streamlined).
* **✅ Good:** "Please rate your experience using the new navigation menu."

### Use Consistent Scale Anchors
Switching scale directions or labels (e.g., using "1 = Poor" in one section and "1 = Strongly Agree" in another) increases cognitive load and leads to "mistake" entries.
* **Best Practice:** Keep anchors consistent across the survey. If you start with a 1–5 scale where 5 is the most positive, stick to that convention throughout.
* **Pro-tip:** Place negative labels on the left and positive on the right (Western reading order).

---

## 2. Advanced Choice Architecture

### Forced Choice vs. "Select All That Apply"
"Select all that apply" often leads to "satisficing," where respondents pick the first few items and skip the rest.
* **Best Practice:** Use **"Pick Top 3"** or **Ranking** to force the participant to make trade-offs.
* **✅ Professional:** "From the list below, rank the top 3 features that are most essential to your daily work."

### Prioritization: MaxDiff & Best-Worst Scaling
Avoid Likert scale matrices (rating 10 items from 1-5) for prioritization, as respondents often rate everything as "Highly Important."
* **Best Practice:** Use **MaxDiff**. Show the user a subset of items and ask them to pick only the **Most Important** and the **Least Important**. This creates a definitive forced-choice hierarchy.

---

## 3. Options & Logical Constraints

### Mutually Exclusive & Exhaustive Ranges
Numerical ranges must not overlap. If they do, the respondent won't know which box to check, and your data will be skewed.
* **❌ Bad (Overlapping):**
    * [ ] 1 - 2 years
    * [ ] 2 - 3 years
* **✅ Good (Mutually Exclusive):**
    * [ ] 1 year to less than 2 years
    * [ ] 2 years to less than 3 years
    * *OR*
    * [ ] 1 - 2
    * [ ] 3 - 4

### "None of the Above" & "Other" Logic
* **Exclusive Logic:** The "None of the above" option must automatically deselect all other selected answers to prevent contradictory data.
* **The "Other" Catch-all:** Always include an "Other (Please specify)" option with a text box for multiple-answer questions.

### Strategic Randomization
* **Randomize:** Multiple-choice lists and question order (within a block) to prevent Primacy/Recency bias.
* **Do NOT Randomize:** Sequential scales (e.g., Likert scales, frequency scales, or numerical ranges) must stay in their logical order.

---

## 4. High-Quality Qualitative Data

### Specific Open-Ended Prompts
Vague questions result in vague data. Force the respondent to ground their answer in a specific memory or behavior.
* **❌ Bad:** "Why did you answer the previous question that way?"
* **✅ Good:** "Please describe a specific example of when the [Feature Name] did not work as you expected. What was the impact on your task?"

---

## 5. Summary Checklist for Researchers

| Rule | Best Practice |
| :--- | :--- |
| **Constructs** | One idea per question. Avoid "and/or" in stems. |
| **Anchors** | Use the same scale labels (e.g., 1-5) across the entire survey. |
| **Prioritization** | Use MaxDiff or Ranking instead of Likert Matrices. |
| **Ranges** | Ensure 0% overlap in numerical categories. |
| **Bias** | Randomize lists; use neutral, non-leading wording. |
| **Logic** | "None of the above" must be a mutually exclusive toggle. |
