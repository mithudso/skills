# Operations: recipes for shepherding a bump

Repo-specific mechanics for 10gen/mongo. Commands and IDs here were correct as
of 2026-07; verify against the current tree when something looks off (the bot
script is the source of truth for its own behavior).

## Bot facts

- Script: `buildscripts/modules/atlas/bump_sls_pin.py` (tests in
  `buildscripts/modules/atlas/tests/test_bump_sls_pin.py`).
- Runs as task `bump_sls_pin` on private variant `sls-pin-bump`, defined in
  `src/mongo/db/modules/atlas/atlas_dev.yml`, cron `0 14 * * 2` (Tuesdays
  14:00 UTC), project `mongodb-mongo-master`.
- Pin lives at `pinned_sls_commit` in `buildscripts/modules/atlas/manifest.json`;
  the import also refreshes two `flags-state.json` copies
  (`buildscripts/modules/atlas/` and
  `src/mongo/db/modules/atlas/jstests/disagg_storage/libs/`) plus SLS protos.
- One open bump PR at a time, dependabot-style: the bot force-pushes
  `server-disagg-storage/bump-sls-pin` and edits the open PR if one exists,
  else opens fresh. Fresh SERVER ticket per bump (revert-by-ticket design).
- Pushes/PRs via the `mongo-pr-bot` GitHub App (credentials are Evergreen
  project expansions).

## The DSC verification patch

The bot submits it via the Evergreen REST API (PUT `/patches/` with
`finalize: false`, then POST `action: finalize`), selecting:

- variant: `enterprise-amazon-linux2023-arm64-all-feature-flags-extra-system-deps`
- tasks: the task TAG selector `.assigned_to_jira_team_server_disagg`
- alias: empty (that is why it is not the `required` set)

Reproduce or extend one by hand with the CLI:

```
evergreen patch -p mongodb-mongo-master --reuse-patch <patch_id> -f
```

(`--reuse-patch` replays the resolved variant/task selection; the CLI's `-t`
takes literal names, not `.tag` selectors, so reuse is the faithful form.)

### Reading patch results

- Spruce UI: `https://spruce.corp.mongodb.com/version/<patch_id>`.
- REST: base `https://evergreen.corp.mongodb.com/api/rest/v2`, bearer token
  from `evergreen client get-oauth-token`. Gotchas: the CLI can hang when the
  cached token is stale; a fresh interactive `evergreen login` run fixes it.
  Never hand-rotate the dex refresh token.
- Raw task logs: append `&text=true` to the `task_log_raw` URL (otherwise you
  get a redirect to the Spruce SPA). Per-test logs: use the test-result's
  TestLogs URL.
- Before trusting a red/green verdict on a hand-submitted patch, verify via
  REST that the patch actually contains the files you think it does
  (`evergreen patch` has silently dropped uncommitted/untracked content).

## Variant A gate: old-pin-safety test

Before landing a fix standalone (Variant A), prove the OLD SLS binary
tolerates it:

1. On a branch containing ONLY the fix (no pin bump), submit a DSC-only patch
   (reuse-patch form above) against master, which is still on the old pin.
2. Green disagg tasks: old-pin-safe; land the fix standalone, then the bot PR
   merges untouched.
3. Red with rejection signatures (e.g. the old binary hard-rejecting a config
   field with `UnknownField(...)`): NOT old-pin-safe; go Variant B.

Old binaries differ per service: historically the old CMS tolerated an unknown
`flag_config` field while the old page/log binaries hard-rejected it. Test,
do not reason by analogy.

## Variant B: bundle assembly

1. Fresh personal branch off current master (do not rebase a stale branch).
2. Import the bot's bump commit: `git cherry-pick -n <bot-commit>` or
   `git checkout <bot-commit> -- <allowlisted paths>`. Then verify the staged
   set is EXACTLY the import allowlist (manifest + flags-state x2 + protos);
   drop anything else.
3. Add the fix as its own commit (squash-merge collapses them; PR title and
   description become the landed commit message).
4. Run format/lint on YOUR files only; keep the imported protos/json
   bot-verbatim.
5. Validate: DSC-only patch first (fast signal), then a normal required patch
   before requesting review.
6. Open a DRAFT PR under a fresh ticket. Close the bot PR with a signed
   comment explaining the supersession and linking the successor.
7. Never push to `server-disagg-storage/bump-sls-pin`.

## Bot-ops (maintainer knowledge)

### Script-version check

Confirm the bump task ran the current script: compare the task's base revision
against the latest master commit touching `bump_sls_pin.py`. A bot run from a
revision predating a script fix can reintroduce fixed failure modes (an early
bump PR was closed for exactly this).

### Orphan tickets (ticket exists, no PR)

The bot cuts the Jira ticket before pushing; if the final `git push` fails, the
task log ends with a "this ticket has no associated PR and needs manual
cleanup" warning. Handling:

- Close the orphan ticket (or have the operator do it).
- `remote: Repository not found` on the push (exit 128) is GitHub masking a
  403: the App's installation token lacks `contents: write` on 10gen/mongo
  even though reads still work. That is a GitHub App/permissions fix
  (DevProd), not a bot-code fix. It recurs every cron until fixed.

### Manual re-run

The variant is patchable. To re-run the bot outside the cron:

```
evergreen patch -p mongodb-mongo-master -v sls-pin-bump -t bump_sls_pin -f
```

A successful re-run cuts its OWN fresh ticket + PR (it does not adopt an
orphan). Close leftover orphans separately.

### Finding the bot's task log

The cron ACTIVATES a per-commit build (requester is `gitter_request`, not
`ad_hoc`). Scan recent gitter_request versions' `/builds` for an `sls-pin-bump`
build with `activated: true`; the task id is
`mongodb_mongo_master_sls_pin_bump_bump_sls_pin_<rev>_<version-create-time as
YY_MM_DD_HH_MM_SS>` (the timestamp is the commit's create time, not 14:00).
Task-type raw log + `&text=true` has the python output.

### ECR images for a pin

SLS images live in ECR account `664315256653`, region us-east-1, repos
`disagg-storage/{log,page,pagematerializer,scheduler,cellmetadata,...}`, image
tag = full SLS githash (= the pin). Retention is count-based (~3 weeks hard
guarantee; see the taxonomy's image-rot entry).

- Existence check:
  `aws --profile mms-scratch ecr describe-images --registry-id 664315256653
  --region us-east-1 --repository-name disagg-storage/log
  --image-ids imageTag=<pin>`
- Docker login for local pulls (expires ~12h):
  `aws --profile mms-scratch ecr get-login-password --region us-east-1 |
  docker login --username AWS --password-stdin
  664315256653.dkr.ecr.us-east-1.amazonaws.com`
  (SSO first if needed: `aws sso login --profile mms-scratch` -- interactive,
  hand to the operator.)

## Local repro (fast iteration only)

Evergreen is the real gate (live-image drift means local green is not
authoritative). For iterating on a specific failing test:

```
python3-venv/bin/python buildscripts/resmoke.py run \
  --suites=buildscripts/modules/atlas/suites/disagg_storage.yml \
  --installDir=bazel-bin/install-dist-test/bin <test path>
```

Gotchas:

- Requires the ECR docker login above; mongod must be built
  (`bazel build install-dist-test`), but a pin change alone needs no rebuild
  (it only changes docker images).
- Write resmoke output to a file; do not pipe to tail (buffering).
- On FAILURE the SLS fixture teardown can hang for a very long time. Do not
  wait: read the relevant SLS container logs for the verdict, then clean up:
  `docker ps -aq --filter name=dcp_ | xargs -r docker rm -f` and remove the
  `dcp_` networks.

## Checkpoint template

Path: `~/Documents/Tickets/<bump-ticket>/checkpoint.md`. Keep it terse; it is
resume state, not a report.

```markdown
# <bump-ticket> shepherding checkpoint
- Updated: <UTC timestamp>
- Bump PR: #<n> (state) | pin: <old> -> <new>
- Phase: orient | audit | triage | green | red-A | red-B | escalated | post-merge
- Competing PRs: none | #<n> (disposition)
- Commit audit: clean | contaminated (<paths>)
- Verification patch: <patch_id> (<status>)
- Failure classification:
  - <task/test>: <taxonomy category> -- <evidence link>
- Strategy: <merge as-is | Variant A (<fix ticket/PR>) | Variant B (<PR>) | SLS escalation>
- Patches submitted by me: <ids + purpose>
- Blocked on: <auth/human action/SLS answer/nothing>
- Next action: <single concrete step>
```

When authorized, post the current checkpoint as a Jira comment on the bump
ticket (signed) so any operator or future session can resume.
