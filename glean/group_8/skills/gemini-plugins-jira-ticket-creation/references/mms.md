# MMS — Jira Ticket Creation Context

## Project Key
`CLOUDP` — all Atlas/MMS work

## Required Fields
- **Component**: Select the Atlas component (e.g. `Atlas UI`, `Atlas API`, `NDS`)
- **Team**: Set to your team label (customfield_12751)
- **Fix Version**: Leave blank unless targeting a specific release
- **customfield_10257**: `{"value": "Not Needed"}` — required for all tickets

## Common Issue Types
- `Story` — new feature or user-facing work (default)
- `Bug` — defect with customer impact
- `Task` — internal/tech debt work
- `Investigation` — research/spike work

## Priority Guidance
- `Critical - P2` — production impact, data loss, security
- `Major - P3` — default for most work
- `Minor - P4` — low priority enhancements

## Team Assignment (customfield_12751)

**CRITICAL FORMAT:** Must be an array of objects with "value" key:

```python
"customfield_12751": [{"value": "Team Name"}]
```

**Common Teams:**
- Atlas Clusters Durability Availability Performance Triage
- Atlas Clusters Fleet Rollout Management
- Atlas Clusters Platform I
- Atlas Clusters Platform II
- Atlas Clusters Platform Triage
- Atlas Clusters Resilience & Recovery
- Atlas Clusters Security Cluster Connectivity
- Atlas Clusters Security Data Plane Security
- Atlas Clusters Security Encryption & Compliance
- Atlas Clusters Security Triage
- Atlas Clusters Workload Management
- APIx DevOps Integrations
- APIx DevTools
- APIx Platform
- Automation
- Automation I
- Automation II
- Backup - Atlas
- Backup - DS
- Backup - Private Cloud
- Billing Platform
- Cloud Native Enablement
- Enterprise Advanced
- Frontend Platform
- IAM Authorization
- IAM Workload Identity
- Monarch
- Payments
- Pricing
- Search
- SRE

**ASK USER FOR TEAM if not obvious from context**

## Epic Link
Custom field ID: `customfield_10857`

```bash
# Create with epic link
jira issue create -p CLOUDP -t Bug -s "Atlas UI: login fails on Safari" \
  --priority High --custom customfield_10857=CLOUDP-275436
```

## Example
```python
jira_create_issue(
    project_key="CLOUDP",
    summary="Atlas cluster fails to restart after maintenance",
    issue_type="Bug",
    description="""h2. Problem
Cluster {{abc123}} failed to restart after maintenance window.

h2. Error
{code}
Error: timeout waiting for cluster to become healthy
{code}

h2. Impact
* 8 minutes of downtime
* Customer-facing service unavailable""",
    components="Atlas",
    additional_fields={
        "priority": {"name": "Critical - P2"},
        "customfield_12751": [{"value": "Automation"}],
        "customfield_10857": "CLOUDP-275436",
        "customfield_10257": {"value": "Not Needed"}
    }
)
```

## Critical Rules
- NEVER omit priority, customfield_12751, or customfield_10257 - ticket will fail
- Team name must match exactly (case-sensitive)
- When uncertain: Story + Major - P3 + Atlas + ask for team
