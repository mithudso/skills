# jira-cli

**Category:** Databases, Data Engineering & Analytics
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/jira-cli/skills/jira-cli

## Description
Use when creating, updating, querying, or transitioning Jira issues in any MongoDB repo.

---

# Jira CLI

> Replace `<PROJECT>` with your Jira project key (e.g., `CLOUDP`, `SERVER`, `DEVPROD`). Check `references/` for repo-specific defaults.

## Fetch Ticket Details

```bash
# View a ticket (summary, description, status, assignee, etc.)
jira issue view <PROJECT>-12345  # see references/ for your project key

# Plain text output (better for parsing)
jira issue view <PROJECT>-12345 --plain

# Include comments
jira issue view <PROJECT>-12345 --comments 5

# Raw API response (full JSON)
jira issue view <PROJECT>-12345 --raw
```

## List Issues

```bash
# List issues in the project
jira issue list -p <PROJECT>  # see references/ for your project key

# Filter by status
jira issue list -p <PROJECT> -s "In Progress"

# Filter by assignee (current user)
jira issue list -p <PROJECT> -a$(jira me)

# Custom JQL query
jira issue list -q "project = <PROJECT> AND status = 'In Progress' AND assignee = currentUser()"
```

## Create, Update, and Transition Issues

```bash
# Create an issue
jira issue create -p <PROJECT> -t Story -s "Title" --priority Medium  # see references/ for your project key

# Assign to yourself
jira issue assign <PROJECT>-12345 $(jira me)

# Transition to a new state
jira issue move <PROJECT>-12345 "In Progress"

# Add a comment
jira issue comment add <PROJECT>-12345 "Fixed in $(git rev-parse --short HEAD)"

# Edit summary or description
jira issue edit <PROJECT>-12345 -s "Updated title"
```

## Other Useful Commands

```bash
# View epics
jira epic list -p <PROJECT>  # see references/ for your project key

# View sprint board
jira sprint list -p <PROJECT> --board <board-id>

# Open ticket in browser
jira open <PROJECT>-12345

# Check current configured user
jira me
```

## Repo-Specific Context

Before running commands, check if a `references/` directory exists alongside this skill.
Look for a file matching your repo name (e.g., `references/mms.md`, `references/server.md`, `references/dsi.md`).
Load it to get the correct project key, workflow states, and common queries for your repo.
Personal overrides in your local `references/` take precedence over central ones.