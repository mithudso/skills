# MongoDB Server — Jira Ticket Creation Context

## Project Key
`SERVER` — all MongoDB Server work

## Epic Link Field
Custom field ID: `customfield_10857`

```bash
# Create with epic link
jira issue create -p SERVER -t Bug -s "WiredTiger cache pressure on replica set" \
  --priority High --custom customfield_10857=SERVER-99999
```

```python
# MCP approach with epic link
jira_create_issue(
    project_key="SERVER",
    summary="WiredTiger cache pressure on replica set",
    issue_type="Bug",
    description="...",
    additional_fields={
        "priority": {"name": "Critical - P2"},
        "customfield_10857": "SERVER-99999",  # Epic link
    }
)
```

## Required Fields
- **Assignee**: Assign immediately to avoid triage queue
- **Component**: e.g. `Storage`, `Replication`, `Sharding`, `Query`
- **Fix Version**: Set to `7.x` or specific version if known

## Common Issue Types
- `Bug` — server defect (use for errors, failures, broken functionality)
- `Improvement` — server enhancement to existing features
- `New Feature` — new functionality
- `Task` — internal/tooling work

## Priority Guidance
- `Blocker - P1` — system completely unusable
- `Critical - P2` — production down, data loss, security
- `Major - P3` — default for most work
- `Minor - P4` — low impact
- `Trivial - P5` — cosmetic

## Example
```python
jira_create_issue(
    project_key="SERVER",
    summary="Aggregation pipeline returns incorrect results with $lookup",
    issue_type="Bug",
    description="""h2. Problem
The {{$lookup}} stage returns duplicate results when the foreign collection has a compound index.

h2. Steps to Reproduce
{code}
db.orders.aggregate([
  { $lookup: { from: "inventory", localField: "item", foreignField: "sku", as: "inventory_docs" } }
])
{code}

h2. Impact
* Incorrect query results returned to application
* Affects all deployments using $lookup with compound indexes""",
    additional_fields={
        "priority": {"name": "Critical - P2"},
        "customfield_10857": "SERVER-90000",
    }
)
```

## Critical Rules
- Always use JIRA wiki markup (h2., {code}, {{}})
- When uncertain: Bug + Major - P3
