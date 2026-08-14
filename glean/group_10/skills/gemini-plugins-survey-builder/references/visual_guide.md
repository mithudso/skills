# Visual & Structural Guidelines for Qualtrics Surveys

Creating a professional, accessible, and user-friendly survey is essential for maximizing response rates and ensuring data quality. A visually cluttered or inconsistent survey leads to respondent fatigue and "straight-lining." 

Follow these guidelines to create a polished and effective Qualtrics instrument.

---

## 1. Typography & Readability
Consistency in text is the hallmark of a professional survey.

* **Font Selection:** Use clean, sans-serif fonts (e.g., Arial, Helvetica, or Open Sans) for better readability on digital screens.
* **Consistent Sizing:** * **Question Text:** 14pt – 16pt.
    * **Answer Choices:** 12pt – 14pt.
    * **Instructional Text:** Use italics or a slightly smaller font (11pt) to differentiate from the main question.
* **Minimal Formatting:** Avoid excessive use of bolding, underlining, or colors within the question text. Use bolding only for emphasis on critical words (e.g., "Select **three** options").

## 2. Structural Organization (Blocks & Themes)
Proper organization helps the respondent understand the flow of the survey and allows for better backend logic management.

* **Thematic Grouping:** Group related questions into **Blocks**. For example, all demographic questions should be in a "Demographics" block.
* **Block Naming:** Every block should have a descriptive title (e.g., "Customer Satisfaction," "Product Usage"). While you can hide these from the respondent, they are vital for your own organizational clarity in the Survey Flow.
* **Transition Text:** Use a "Descriptive Text" question at the start of a new block to provide context or instructions for the upcoming section.

## 3. Pacing & Page Breaks
Long scrolling pages overwhelm respondents. Strategic breaks keep them engaged.

* **One Topic per Page:** Ideally, put page breaks between different topics or question types.
* **Avoid "The Wall of Questions":** Limit each page to 1–3 short questions. For complex matrix questions, the question should likely sit on its own page.
* **Force Page Breaks:** Use the "Add Page Break" feature between questions within a block to control the rhythm.

## 4. Visual Layout & Theme
Qualtrics provides various themes; choose one that minimizes distraction.

* **Clean Themes:** Use the "Minimal" or "Static" themes. Avoid backgrounds with busy images or high-contrast patterns.
* **Progress Bars:** Always include a progress bar (placed at the top or bottom). This reduces drop-off by showing the respondent how much effort remains.
* **Logos:** If using a company or institutional logo, ensure it is high resolution and centered or left-aligned at the top. Do not let it take up more than 15% of the screen height.

## 5. Question Design Consistency
* **Scale Direction:** Keep scales consistent throughout the survey. If "1" is "Strongly Disagree" in the first block, do not make it "Strongly Agree" in the next.
* **Mobile Optimization:** Always enable "Mobile Friendly" for Matrix questions. This converts them into an accordion style on smartphones, preventing horizontal scrolling.
* **Alignment:** Standardize the alignment of text and buttons. Left-aligned text is generally easier to read for Western audiences.

## 6. Buttons & Navigation
* **Standardized Buttons:** Use "Next" and "Back" (if allowed). Ensure the button text is clear.
* **The "Back" Button:** Enable the back button unless your survey logic (like randomization) strictly forbids it. Respondents appreciate being able to correct a mistake.

## 7. Accessibility (WCAG Compliance)
* **High Contrast:** Ensure text has a high contrast ratio against the background (e.g., black text on a white or very light grey background).
* **Alt-Text:** If you use images, provide descriptive Alt-Text for screen readers.
* **Check Accessibility:** Use the "Check Survey Accessibility" tool under the *Tools* menu in Qualtrics before launching.

---
