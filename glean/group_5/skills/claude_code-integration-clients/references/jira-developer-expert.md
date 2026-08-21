<!-- hub-reference-banner -->
> **Reference file — part of the `integration-clients` hub.** Formerly the standalone `jira-developer-expert` skill.
> Sibling topics in this family are now reference files under the hubs (`integration-clients`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: jira-developer-expert
description: >
  Jira Cloud developer expert: REST API v3 (issues, search, workflows), Jira Software agile API
  (boards, sprints, epics), Jira Service Management API (requests, SLAs, organizations), Forge app
  development (modules, triggers, web triggers, storage), OAuth 2.0 3LO auth flows, webhooks,
  and Connect-to-Forge migration guidance.
  TRIGGER: user is building a Jira integration, Forge app, or webhook handler; asking about Jira
  REST endpoints, JQL, Forge modules, OAuth scopes, pagination, or JSM request types; migrating
  a Connect app to Forge.
  SKIP: general project management usage (not developer/API questions); non-Jira Atlassian products
  (Confluence, Bitbucket) unless the question is cross-product auth — use atlassian skills instead.
version: "1.1.0"
updated: "2026-05-29"
category: developer
origin: local
tags:
  - jira
  - atlassian
  - forge
  - rest-api
  - oauth
  - webhooks
  - jira-software
  - jira-service-management
keywords:
  - jira
  - jira cloud
  - REST API v3
  - JQL
  - Forge
  - Connect
  - OAuth 2.0
  - 3LO
  - webhooks
  - boards
  - sprints
  - epics
  - JSM
  - service management
  - SLA
  - request types
  - pagination
  - rate limiting
whenToUse:
  - Building a Jira Cloud integration using REST API v3
  - Developing a Forge app with modules, triggers, or web triggers
  - Implementing OAuth 2.0 3LO auth flow for Jira
  - Setting up Jira webhooks and making handlers idempotent
  - Using the Jira Software Agile API for boards, sprints, and epics
  - Working with Jira Service Management requests, SLAs, or organizations
  - Migrating a Connect app to Forge
  - Handling 401/403/404/429 errors correctly in a Jira integration
whenNotToUse:
  - General Jira usage questions with no API/developer angle
  - Confluence, Bitbucket, or Trello API questions — use appropriate Atlassian skill
  - Jira Data Center or Server (on-prem) — this skill covers Cloud only
related_skills:
  - atlassian
  - jira-cli
  - jira-ticket-creation
---

# Jira Developer Expert

Jira Cloud developer reference: REST API v3, Jira Software, Jira Service Management, Forge apps,
OAuth 2.0 3LO, webhooks, and Connect-to-Forge migration.

**Forge is the default for all new Jira Cloud apps. Treat Connect as legacy-only.**

## Quick API Reference

| Operation | Method + Path |
|---|---|
| Create issue | `POST /rest/api/3/issue` |
| Edit issue | `PUT /rest/api/3/issue/{issueIdOrKey}` |
| Transition issue | `POST /rest/api/3/issue/{issueIdOrKey}/transitions` |
| Add comment | `POST /rest/api/3/issue/{issueIdOrKey}/comment` |
| Upload attachment | `POST /rest/api/3/issue/{issueIdOrKey}/attachments` |
| Search with JQL | `GET` or `POST /rest/api/3/search` |
| List boards | `GET /rest/agile/1.0/board` |
| List board sprints | `GET /rest/agile/1.0/board/{boardId}/sprint` |
| Move issues to sprint | `POST /rest/agile/1.0/sprint/{sprintId}/issue` |
| List customers | `GET /rest/servicedeskapi/customer` |
| List organizations | `GET /rest/servicedeskapi/organization` |
| Get request SLAs | `GET /rest/servicedeskapi/request/{issueIdOrKey}/sla` |
| List request types | `GET /rest/servicedeskapi/servicedesk/{serviceDeskId}/requesttype` |

## Auth and Scopes

**OAuth 2.0 (3LO) endpoints:**
- `https://auth.atlassian.com/authorize`
- `https://auth.atlassian.com/oauth/token`
- `https://auth.atlassian.com/oauth/revoke`

**Common scopes** (request least-privilege only):

| Scope | Use for |
|---|---|
| `read:jira-work` | Read issues, projects, workflows |
| `write:jira-work` | Create and edit issues |
| `read:jira-user` | Read user profiles |
| `manage:jira-project` | Manage projects |
| `manage:jira-configuration` | Manage instance configuration |

## Webhook Rules

- Prefer webhooks/triggers over polling.
- Apply JQL filters aggressively to limit event volume.
- Assume retries and duplicate delivery — make handlers idempotent.
- Keep handlers fast; defer heavy work to async queues.
- Use Forge triggers/web triggers for new cloud app work.

## Coding Standards

1. **Forge-first** for all new Jira Cloud apps.
2. Use `async/await`; never mix callbacks.
3. **Paginate every list endpoint** — never assume a single page covers all results.
4. Handle `429` with exponential backoff and respect `Retry-After`.
5. Never hardcode secrets — use Forge Storage API or environment variables.
6. Do not key logic off custom field display names when stable field IDs exist.
7. Treat Connect as legacy-only — new features go in Forge.

## Primary Sources

| Topic | URL |
|---|---|
| Jira REST API v3 | https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/ |
| Jira webhooks | https://developer.atlassian.com/cloud/jira/platform/webhooks/ |
| OAuth 3LO | https://developer.atlassian.com/cloud/jira/platform/oauth-2-3lo-apps/ |
| Jira scopes | https://developer.atlassian.com/cloud/jira/platform/scopes/ |
| Jira Software REST | https://developer.atlassian.com/cloud/jira/software/rest/ |
| Jira Service Management REST | https://developer.atlassian.com/cloud/jira/service-desk/rest/intro/ |
| Forge home | https://developer.atlassian.com/platform/forge/ |
| Forge manifest reference | https://developer.atlassian.com/platform/forge/manifest-reference/ |
| Forge security best practices | https://developer.atlassian.com/platform/forge/security-best-practices/ |
| Connect modules for Jira | https://developer.atlassian.com/cloud/jira/platform/about-connect-modules-for-jira/ |
