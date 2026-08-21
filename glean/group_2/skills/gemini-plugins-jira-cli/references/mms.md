# MMS / Atlas — Jira CLI Context

## Default Project
- Primary: `CLOUDP` (Atlas/MMS features, bugs, tasks)
- Secondary: `HELP` (customer-facing support tickets)

## Workflow States
Open → In Progress → In Review → Done

## Common Queries

```bash
# My open CLOUDP tickets
jira issue list -p CLOUDP --assignee $(jira me) --status "In Progress"

# View a CLOUDP ticket
jira issue view CLOUDP-123456

# Create a CLOUDP ticket
jira issue create -p CLOUDP -t Story -s "Title" --priority Medium

# List HELP tickets for triage
jira issue list -p HELP --status "Open" --order-by created
```

## Ticket ID Format
`CLOUDP-######` — six digits
