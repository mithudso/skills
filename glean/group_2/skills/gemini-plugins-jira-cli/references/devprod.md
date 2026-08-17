# Jira CLI — DevProd / Build & Virtualization Context

## Projects
- `DEVPROD` — Developer Productivity team issues
- `BUILD` — Build & Virtualization team issues
- `BF` — Build Failures (auto-generated from Evergreen)

## Common Queries

```bash
# List open DEVPROD issues assigned to me
jira issue list -p DEVPROD --assignee $(jira me) --status "In Progress"

# View a BUILD issue
jira issue view BUILD-1234

# Create a DEVPROD ticket
jira issue create -p DEVPROD -t Story -s "Title" --priority Medium
```

## Workflow States
DEVPROD uses: Open → In Progress → In Review → Done
