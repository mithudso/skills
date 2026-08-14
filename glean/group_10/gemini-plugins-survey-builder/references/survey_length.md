# Best Practices for Survey Length: Developers & Dev Tools

When designing surveys for developers, engineers, and users of technical tools, efficiency is paramount. This demographic is highly sensitive to "context switching" and time-sink activities. Adhering to strict length and structure guidelines ensures higher response rates and better data quality.

## 1. Completion Time Benchmarks
Time is the most critical metric for developer surveys. Always state the estimated time upfront.

* **Standard Target:** Aim for a **5-minute completion time**. This is the "sweet spot" for maintaining focus without causing drop-off.
* **Complex/Forced Choice:** For deep-dive research or surveys using "forced choice" (MaxDiff) methodologies, you may extend to **10 minutes**, provided the value proposition to the developer is clear.
* **The Hard Ceiling:** Never exceed **15 minutes**. Beyond this point, data quality degrades significantly due to survey fatigue.

## 2. In-Product & Intercept Surveys
Surveys delivered via "toast" notifications, banners, or within a CLI/IDE environment must be minimally intrusive.

* **The 2-Question Rule:** In-product intercepts should be **2 questions maximum**.
    * *Example:* One quantitative (e.g., NPS or Rating) and one qualitative (e.g., "Why?").
* **Zero-Disruption Goal:** Ensure the survey does not block the developer's primary workflow or command execution.

## 3. Question Count Guidelines
* **General Surveys:** Aim for **8–10 questions**.
* **Exhaustive Research:** Cap at **15 questions** total. If you need more data, consider splitting the survey into two different cohorts (A/B testing the questions).
* **Screeners:** Limit screening questions to **2–3**. Do not make users work just to find out they aren't eligible.

## 4. Response Scales & Logic
Standardization helps developers move through the survey quickly by reducing cognitive load.

* **Likert Scales:** * Prefer **5-point scales** for simplicity and speed.
    * Use **7-point (or higher) scales** only when measuring very subtle nuances in sentiment where high granularity is required for statistical significance.
* **Multi-Select Options:** Limit multiple-choice lists to a **maximum of 10 options**.
    * This count must include "Other (please specify)" and "None of the above."
    * If a list exceeds 10, use a searchable dropdown or categorize the options.
* **Binary Logic:** Use Yes/No questions to branch (skip logic) only when necessary to keep the path short.

## 5. Instructions & Copywriting
Developers value precision and hate marketing fluff.

* **Directness:** Instructions should be **clear, technical, and direct**. 
    * *Bad:* "We would be so grateful if you could take a moment to share your thoughts on our new API."
    * *Good:* "Rate the latency of the v2.0 Auth endpoint."
* **No Ambiguity:** Define technical terms if there is any chance of regional or stack-based variation in meaning.
* **Progress Indicators:** Always include a progress bar so the user knows exactly how close they are to the end.

## 6. Optimization for Technical UX
* **Mobile vs. Desktop:** While most developers take surveys on desktops (often while coding), ensure the survey UI handles code snippets or technical diagrams gracefully on all screens.
* **Keyboard Navigation:** Ensure the survey can be completed using only the keyboard (Tab/Enter). Power users will appreciate the ability to fly through the form.
* **The "Save and Continue" Fallacy:** For developer surveys, do not rely on "save and continue." If they leave the tab, they likely won't come back. Keep it short enough to finish in one sitting.
