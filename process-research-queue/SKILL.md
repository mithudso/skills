---
name: process-research-queue
description: >-
  Automates the research of pending concepts in the global research queue and injects them into the 3D concept tree.
version: 1.0.0
updated: 2026-08-18
category: workflow
model: claude-sonnet-5
effort: medium
triggers:
  - process the research queue
  - /process-queue
  - research pending concepts
whenToUse:
  - "I want to research the items in my queue"
  - "process the concept queue"
---

# Process Research Queue Workflow

This skill automates the process of reading the global research queue, conducting deep research on pending items, and integrating the findings into the global concept tree.

## Step 1: Read the Queue
1. Use `view_file` to read `~/.global-ai-hub/concept-tree/RESEARCH_QUEUE.md`.
2. Identify all incomplete items denoted by `- [ ] Concept: \`[Topic]\` | Parent: \`[Parent Concept]\``.
3. If there are no incomplete items, inform the user and exit.

## Step 2: Conduct Deep Research
For each incomplete item:
1. Use `invoke_subagent` to spawn a `research` subagent to investigate the `[Topic]`.
2. The prompt should instruct the subagent to conduct deep web research on the concept and return a concise summary of the key findings.

## Step 3: Update the Concept Tree
1. Read `~/.global-ai-hub/concept-tree/tree.json` using `view_file`.
2. For each researched concept, add a new node to the JSON array:
   ```json
   {
     "concept": "[Topic]",
     "parentConcept": "[Parent Concept]",
     "researchedAt": "YYYY-MM-DD",
     "childConcepts": []
   }
   ```
   *Note: Ensure the `parentConcept` matches an existing concept in the tree. If it doesn't, attach it to a logical parent or root (`null`). Also, make sure to add the new topic to the `childConcepts` array of the parent node.*
3. Use `write_to_file` to save the updated `tree.json`.

## Step 4: Mark Queue as Complete
1. Use `multi_replace_file_content` or `replace_file_content` to modify `~/.global-ai-hub/concept-tree/RESEARCH_QUEUE.md`.
2. Change the `- [ ]` checkboxes to `- [x]` for all the items you successfully processed.
3. Provide a summary of the concepts researched to the user.
