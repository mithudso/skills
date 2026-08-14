# Evergreen Patch Creation — Detailed Reference

## CLI Flag Reference

| Flag                 | Short           | Description                                                                  |
| -------------------- | --------------- | ---------------------------------------------------------------------------- |
| `--project`          | `-p`            | Evergreen project identifier (e.g., `<PROJECT>`)                  |
| `--variants`         | `-v`            | Variant names, comma-separated. Use `"all"` for every variant                |
| `--tasks`            | `-t`            | Task names, comma-separated. Use `"all"` for every task in selected variants |
| `--alias`            | `-a`            | Use a predefined patch alias (configured by project admins)                  |
| `--finalize`         | `-f`            | Schedule tasks immediately without manual finalization                       |
| `--skip_confirm`     | `-y`            | Skip interactive confirmation prompts                                        |
| `--uncommitted`      | `-u`            | Include uncommitted local changes in the patch diff                          |
| `--description`      | `-d`            | Patch description text                                                       |
| `--auto-description` | `--ad`          | Use the last git commit message as description                               |
| `--browse`           |                 | Open the patch URL in the default browser                                    |
| `--ref`              |                 | Commit reference to test on top of merge-base (default: HEAD)                |
| `--large`            | `-l`            | Allow patches larger than 16MB                                               |
| `--json`             | `-j`            | Output patch as JSON; suppresses warnings/confirmations                      |
| `--verbose`          |                 | Show patch summary after submission                                          |
| `--param`            |                 | Set a parameter as `KEY=VALUE` (repeatable)                                  |
| `--path`             |                 | Path to Evergreen project config file                                        |
| `--regex_variants`   | `--rv`          | Select variants by regex pattern                                             |
| `--regex_tasks`      | `--rt`          | Select tasks by regex pattern                                                |
| `--include-modules`  |                 | Include module diffs from defined module paths                               |
| `--trigger-alias`    |                 | Trigger alias for tasks from other projects                                  |
| `--repeat`           | `--reuse`       | Reuse the same tasks/variants from the latest patch                          |
| `--repeat-failed`    | `--rf`          | Reuse only the failed tasks/variants from the latest patch                   |
| `--repeat-patch`     | `--reuse-patch` | Reuse tasks/variants from a specific patch ID                                |

## Project Detection

The project identifier may vary by branch or repo. Check your repo-specific references file (e.g., `references/server.md`, `references/mms.md`) for project naming conventions.

## Variant -> Task Mapping

Variant and task names vary by repository. Check your repo-specific references file (e.g., `references/server.md`, `references/mms.md`) for available variants, tasks, and common patterns.

## Patch Aliases

Aliases are predefined combinations of variants and tasks configured by project admins. Use with `-a <alias>`.

To discover all available aliases:

```bash
evergreen list --patch-aliases -p <PROJECT>
```

Check your repo-specific references file (e.g., `references/server.md`, `references/mms.md`) for commonly used aliases in your project.

## Patch Management Commands

### List existing patches

```bash
evergreen list-patches                    # Show 5 most recent
evergreen list-patches -n 10              # Show 10 most recent
evergreen list-patches -n 0               # Show all patches
evergreen list-patches -j                 # JSON output
evergreen list-patches -s                 # Show diff summary
evergreen list-patches -i <patch_id>      # Show specific patch
```

### Finalize an unfinalized patch

```bash
evergreen finalize-patch -i <patch_id>
```

### Cancel a patch

```bash
evergreen cancel-patch -i <patch_id>
```

### Validate project config

```bash
evergreen validate -f etc/evergreen.yml
evergreen validate -f etc/evergreen.yml -p <PROJECT>    # With project-specific validation
```

### Evaluate expanded config

```bash
evergreen evaluate -f etc/evergreen.yml --variants    # Show variant definitions
evergreen evaluate -f etc/evergreen.yml --tasks       # Show task definitions
```

### List project configuration

```bash
evergreen list --variants -p <PROJECT>          # All variant names
evergreen list --tasks -p <PROJECT>             # All task names
evergreen list --patch-aliases -p <PROJECT>     # All patch aliases
evergreen list --parameters -p <PROJECT>        # All project parameters
```
