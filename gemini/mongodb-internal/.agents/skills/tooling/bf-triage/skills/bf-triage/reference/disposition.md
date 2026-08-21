# BF disposition vocabulary

Two orthogonal closed-set enums consumed by the report's
`Recommended next steps` section and the Run-summary block.

The two axes are independent — many real cases need both fields to
fully describe the recommendation. Use `unknown` / `mixed` /
prose escape hatches when the agent cannot commit to a value within
the time budget.

## `fix_location` — where does the code change go?

| Value | Means | Typical paths |
| --- | --- | --- |
| `mongo_test` | Test / scaffolding-only change | `jstests/`, `buildscripts/`, `*_test.cpp` |
| `mongo_server` | Server-side C++/Rust change | `src/` production code |
| `dsi_config` | DSI YAML / variant / mongodb-setup change | `configurations/`, `evergreen/system_perf/master/variants.yml`, `configurations/mongodb_setup/`, `configurations/infrastructure_provisioning/` |
| `dsi_workload` | DSI locust workload source change | `workloads/<team>/<name>/src/` |
| `evergreen_config` | CI config change | `etc/evergreen.yml`, `evergreen/buildvariants/` |
| `infra_image` | Build host / image / AMI change | `kudzu/` image package list, AMI registry |
| `cross_team` | Primary fix lives in another team's project | Pair with the receiving team in `Re-route` step; `disposition` is usually `keep_open_pending_fix` or `close_duplicate` |
| `none` | No code change is needed | Pair with a `disposition` value below |
| `mixed` | More than one of the above | Spell out the split in prose under Recommended next steps |
| `unknown` | Agent could not determine within budget | Pair with `disposition = keep_open_investigating` |

## `disposition` — what happens to the BF ticket itself?

| Value | Means |
| --- | --- |
| `keep_open_pending_fix` | BF stays open; linked fix in some project must land first |
| `keep_open_accept_as_known` | BF stays open as known-bad tracking; no fix planned |
| `keep_open_investigating` | Needs more investigation before a firm call |
| `close_duplicate` | Duplicate of another BF — link as `duplicates` |
| `close_gone_away` | Caused-by and fixed-by SHAs cancel each other (common for perf-change BFs the detector auto-fixed) |
| `close_wont_fix` | Explicitly decided not to fix |

## Common pairings

| Real case | `fix_location` | `disposition` |
| --- | --- | --- |
| Test-only flake widened to fit observed variance | `mongo_test` | `keep_open_pending_fix` |
| Real server bug with a SERVER-N follow-up filed | `mongo_server` | `keep_open_pending_fix` |
| Backport-presence check found a missing cherry-pick | `mongo_server` (or `mongo_test`) | `keep_open_pending_fix` |
| Variant-incompatibility / DSI config drift | `dsi_config` | `keep_open_pending_fix` |
| Workload source change in `workloads/<team>/<name>/src/` | `dsi_workload` | `keep_open_pending_fix` |
| Build host / AMI flake (one-off) | `infra_image` or `none` | `close_wont_fix` |
| Perf-change BF, detector named `fixed by <sha>` that cancels caused-by | `none` | `close_gone_away` |
| Perf change real but expected (team accepts the regression) | `none` | `keep_open_accept_as_known` |
| Duplicate of an open BF | `none` | `close_duplicate` |
| Owning team is not the current one (re-route) | `cross_team` | `keep_open_pending_fix` |
| Out-of-budget / inconclusive | `unknown` | `keep_open_investigating` |

## Reporting rules

- Both fields are mandatory in the Run summary block. Use
  `unknown` rather than blank if the agent cannot commit.
- The `Recommended next steps` section's prose must be **consistent
  with** the two fields — e.g. if `disposition = close_gone_away`,
  the prose must justify closing (cite the caused-by and fixed-by
  SHAs), not recommend Reproduce.
- For `cross_team` re-routes, name the receiving team and cite the
  routing-rule source from `team_knowledge.md` /
  `workflow_overview.md`.
- For perf-change BFs, the
  [`log_patterns.md`](log_patterns.md) § "Performance-change BF
  interpretation" table maps each detector signal to the expected
  `(fix_location, disposition)` pair — use it as the default.
