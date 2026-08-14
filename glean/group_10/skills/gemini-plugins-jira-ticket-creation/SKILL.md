---
name: jira-ticket-creation
description: >-
  Use when creating Jira tickets or bug reports. Covers Wiki Markup syntax,
  description templates, and field selection for any MongoDB Jira project.
source: 10gen/mongo
license: Internal
mongodb:
  team: devprod-bv
  owner: srdjan.pajic@mongodb.com
  internal: true
---

# JIRA Ticket Creation

> Replace `<PROJECT>` with your Jira project key. Check `references/` for your repo's project key, epic link field, and required fields.

**PREFER the `jira` CLI over Jira MCP tools.** See the `jira-cli` skill for CLI-based ticket creation. Fall back to the MCP approach below only if the CLI is unavailable or for operations the CLI cannot perform.

## Required Fields - ALL MUST BE PROVIDED

Every ticket creation MUST include these exact fields or it will fail:

```python
jira_create_issue(
    project_key="<PROJECT>",  # see references/ for your project key
    summary="<title>",      # Clear, specific, <255 chars
    issue_type="<type>",    # Bug|Improvement|New Feature|Task|Story
    description="<details>", # Include context, steps, impact
    additional_fields={
        "priority": {"name": "<priority>"},        # Blocker - P1 | Critical - P2 | Major - P3 | Minor - P4 | Trivial - P5
    }
)
```

## JIRA Description Formatting

**CRITICAL: JIRA uses Wiki Markup syntax, NOT Markdown!**

| Element     | JIRA Wiki Markup      | Example                        |
| ----------- | --------------------- | ------------------------------ |
| Heading 2   | `h2. Title`           | `h2. Problem`                  |
| Heading 3   | `h3. Title`           | `h3. Details`                  |
| Code Block  | `{code:cpp}...{code}` | `{code:cpp}if (x) {...}{code}` |
| Inline Code | `{{text}}`            | `{{methodName()}}`             |
| Bullet List | `* item`              | `* First item`                 |
| Bold        | `*text*`               | `*important*`                  |

## Field Selection Rules

**Issue Type:**

- Bug: Errors, failures, broken functionality, incorrect behavior
- Improvement: Enhancements to existing features
- New Feature: New functionality (default for feature work)
- Task: General work items, maintenance, process improvements
- Story: User-facing features or work items

**Priority:**

- Blocker - P1: System completely unusable, blocks all progress
- Critical - P2: Production down, data loss risk, security
- Major - P3: Everything else (default)
- Minor - P4: Low impact, nice-to-have
- Trivial - P5: Cosmetic, minor inconvenience

## Description Formatting

**CRITICAL: JIRA uses Wiki Markup, NOT Markdown**

Always format descriptions using JIRA wiki markup syntax:

**Headers:**

```
h1. Header 1
h2. Header 2
h3. Header 3
```

**Code/Monospace:**

```
{{code}} or {{methodName()}} or {{src/path/file.cpp}}
```

**Lists:**

```
* First level bullet
** Second level bullet
*** Third level bullet
# Numbered list
## Numbered sublist
```

**Text Formatting:**

```
*bold*
_italic_
{{monospace}}
```

**Panels for Callouts:**

```
{panel:title=Note}
Important information here
{panel}

{panel:title=Warning|borderStyle=solid|borderColor=#ccc|titleBGColor=#F7D6C1}
Warning message here
{panel}
```

**Tables:**

```
|| Heading 1 || Heading 2 || Heading 3 ||
| Cell 1-1 | Cell 1-2 | Cell 1-3 |
| Cell 2-1 | Cell 2-2 | Cell 2-3 |
```

**Code Blocks with Language:**

```
{code:java}
public void method() {
    // Java code
}
{code}

{code:python}
def function():
    # Python code
{code}

{code:bash}
# Bash commands
ls -la
{code}
```

**Common Pattern for Structured Descriptions:**

```
h2. Overview
Brief summary of what needs to be done.

h2. Background
Context and motivation for this work.

h2. Scope of Work

h3. 1. Implementation Changes
*Files to modify:*

* {{src/path/to/file1.cpp}}
** Add {{validateMethod()}} method
** Update {{doGetNext()}} to handle deduplication

* {{src/path/to/file2.h}}
** Add new member variable for tracking seen keys

h3. 2. Test Coverage
*File to modify:*
* {{src/path/to/test_file.cpp}}
** Add test methods covering:
*** Happy path scenarios
*** Error conditions
*** Edge cases

h2. Acceptance Criteria
* Criterion 1 completed
* Criterion 2 validated
* All tests pass

h2. Technical Notes
* Important implementation detail 1
* Important implementation detail 2

h2. Related Files
* Implementation: {{src/path/to/file1.cpp}}
* Tests: {{src/path/to/test_file.cpp}}
```

## Example

```python
jira_create_issue(
    project_key="<PROJECT>",
    summary="Fix aggregation pipeline incorrect results",
    issue_type="Bug",
    description="""h2. Problem
The aggregation pipeline returns duplicate results when using compound indexes.

h2. Steps to Reproduce
{code}
db.orders.aggregate([
  { $lookup: { from: "inventory", localField: "item", foreignField: "sku", as: "inventory_docs" } }
])
{code}

h2. Impact
* Incorrect query results returned to application
* Affects all deployments using this pattern""",
    additional_fields={
        "priority": {"name": "Critical - P2"},
    }
)
```

## Critical Rules

- Field values are case-sensitive - use exactly as shown
- **ALWAYS use JIRA wiki markup in descriptions (h2., {code}, {{}}), NEVER Markdown**
- When uncertain: Bug + Major - P3

## Common Errors and Troubleshooting

### Error: "Error calling tool 'create_issue'"

**Most Common Causes:**

1. **Missing required fields**

   - Must include: project_key, summary, issue_type, description
   - Must include in additional_fields: priority

2. **Invalid priority format**
   - WRONG: `"priority": "Major - P3"`
   - CORRECT: `"priority": {"name": "Major - P3"}`

### Testing Ticket Creation

If unsure about field format, you can:

1. Get an existing ticket with `jira_get_issue_MCP_DOCKER` using `fields="*all"`
2. Examine the exact format of priority, etc.
3. Use that exact format in your ticket creation

## Repo-Specific Context

Before creating a ticket, check if a `references/` directory exists alongside this skill.
Look for a file matching your repo name (e.g., `references/mms.md`, `references/server.md`).
Load it to get the correct project key, epic link field, and required fields for your repo.
Personal overrides in your local `references/` take precedence over central ones.
