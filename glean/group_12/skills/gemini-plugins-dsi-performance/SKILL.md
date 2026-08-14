---
name: dsi-performance
description: Use when querying Evergreen CI/CD data, finding failed tasks, checking waterfall status, fetching task logs, or analyzing performance test results via the Evergreen GraphQL or REST API.
source: 10gen/dsi
license: Internal
mongodb:
  team: release-quality
  owner: jawwad.asghar@mongodb.com
  internal: true
---

# DSI Performance Testing — Evergreen API

Query Evergreen CI/CD data using the **GraphQL API** (preferred) or REST API v2.

**Use GraphQL** for most queries - it allows fetching nested data in a single request and requesting only needed fields.

**Use REST API** for downloading task logs (GraphQL doesn't support this).

## Authentication

Both APIs use the same credentials from the Evergreen CLI:

```bash
evergreen client get-user    # Username
evergreen client get-api-key # API key
```

| API     | Endpoint                                      |
| ------- | --------------------------------------------- |
| GraphQL | `https://evergreen.mongodb.com/graphql/query` |
| REST v2 | `https://evergreen.mongodb.com/rest/v2`       |

## Key Queries

| Query                        | Description                              |
| ---------------------------- | ---------------------------------------- |
| `mainlineCommits(options)`   | List versions for a project              |
| `version(versionId)`         | Get version with buildVariants and tasks |
| `task(taskId, execution?)`   | Get task details with tests              |
| `patch(id)`                  | Get patch details                        |
| `project(projectIdentifier)` | Get project configuration                |

**Entity hierarchy:** project → versions → builds/buildVariants → tasks → tests

## Common Patterns

**Get recent mainline commits for sys-perf:**

```graphql
{
  mainlineCommits(options: { projectIdentifier: "sys-perf", limit: 5 }) {
    versions {
      version { id revision author message createTime }
    }
  }
}
```

**Filter for failed tasks:**

```graphql
{
  mainlineCommits(
    options: { projectIdentifier: "sys-perf", limit: 5 }
    buildVariantOptions: { statuses: ["failed", "system-failed"] }
  ) {
    versions {
      version {
        id
        buildVariants(options: { statuses: ["failed", "system-failed"] }) {
          displayName
          tasks { displayName status }
        }
      }
    }
  }
}
```

**Download task logs (REST):**

```bash
curl -H "Api-User: $(evergreen client get-user)" \
     -H "Api-Key: $(evergreen client get-api-key)" \
     "https://evergreen.mongodb.com/rest/v2/tasks/TASK_ID/build/TaskLogs?type=task_log"
```

## Task Statuses

| Status           | Meaning                 |
| ---------------- | ----------------------- |
| `success`        | Passed                  |
| `failed`         | Failed                  |
| `system-failed`  | Infrastructure failure  |
| `task-timed-out` | Exceeded timeout        |
| `will-run`       | Scheduled, waiting      |
| `blocked`        | Waiting on dependencies |

## Requesters

| Value                 | Description         |
| --------------------- | ------------------- |
| `gitter_request`      | Mainline commit     |
| `patch_request`       | CLI patch           |
| `github_pull_request` | GitHub PR           |

See `references/evergreen.md` for performance patch submission and monitoring.
