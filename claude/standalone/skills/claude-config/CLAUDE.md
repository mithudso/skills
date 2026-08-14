# Global instructions

## Be concise

Default to concise responses. Lead with the answer or outcome; cut preamble,
restated context, and option surveys. Give a recommendation, not an exhaustive
list. Expand only when the task genuinely needs detail (multi-step plans,
trade-off decisions, or when asked).

## Task visibility (always keep progress visible)

For any task with **3 or more distinct steps**, maintain a live task list so it's
always clear what's being worked on and how much is left:

- **Create it up front.** Call `TaskCreate` to add one item per step at the start
  of the work — use the task tools (`TaskCreate` / `TaskUpdate` / `TaskList`), not
  just prose narration.
- **One in-progress at a time.** Keep exactly one task `in_progress`; mark it
  `completed` the moment it's done and set the next one `in_progress`.
- **Update every turn.** Refresh task status at the end of each turn so the live
  checklist — and the status line — always reflect what's done and how many steps
  remain.
- **Granular over vague.** Prefer small, verifiable steps to a few broad ones; it
  makes the "steps left" count meaningful.

This applies to multi-step work in any project. Skip it only for genuinely
single-step requests.
