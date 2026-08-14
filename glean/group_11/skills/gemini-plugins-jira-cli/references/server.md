# Jira CLI — MongoDB Server Context

## Projects
- `SERVER` — MongoDB Server issues
- `WT` — WiredTiger issues
- `TOOLS` — MongoDB Tools

## Common Queries

```bash
# List open SERVER issues assigned to me
jira issue list -p SERVER --assignee $(jira me) --status "In Progress"

# Find issues linked to a commit
jira issue list -p SERVER --jql "text ~ \"<commit-hash>\""
```

## Workflow States
SERVER uses: Open → In Progress → Waiting for Review → Closed
