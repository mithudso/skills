# variant-right-sizer

**Category:** Cloud, DevOps & Infrastructure
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/variant-right-sizer/skills/variant-right-sizer

## Description
Use when optimizing Evergreen variant instance sizes for cost. Iteratively scales down instance sizes, tests via patches, and finds the smallest viable size for each distro used by a variant.

---

# Variant Right-Sizer

Find the smallest viable instance size for an Evergreen build variant by iteratively stepping down sizes and testing via patches.

## Required Inputs

Gather from the user before starting:
- **Variant name** (e.g., `enterprise-rhel8-arm64-debug-tsan`)
- **File where the variant is defined** (optional — will search Evergreen YAML files if not provided)
- **Evergreen project ID** (optional — infer from repo or ask)

## Step 1: Find the Variant

Search the Evergreen YAML configuration for the variant definition:

```bash
# Search for the variant name in all YAML files
grep -rn "name: <variant_name>" etc/evergreen_yml_components/ .evergreen.yml evergreen/
```

Read the full variant block including its `run_on`, `expansions`, and `tasks` sections.

## Step 2: Inventory All Instance Sizes

Identify every instance size reference in the variant. Record each as `(context, current_distro)`:

| Context | Where to Look |
|---------|--------------|
| `run_on` | The variant's default distro |
| `large_distro_name` | In the variant's `expansions` block or inherited anchor |
| `xlarge_distro_name` | In the variant's `expansions` block or inherited anchor |
| `core_analyzer_distro_name` | In the variant's `expansions` block or inherited anchor |
| Task-level `distros` | Individual tasks or task groups that override the default distro |

Also check expansion anchors (e.g., `<<: *some-template`) — follow the anchor to find distro references defined there.

Build a list like:

```
1. run_on: rhel8.8-arm64-m8g-8xlarge
2. large_distro_name: rhel8.8-arm64-m8g-8xlarge
3. task compile_TG distro: rhel8.8-arm64-m8g-8xlarge
4. task jsCore distro: rhel8.8-arm64-m8g-8xlarge
```

Deduplicate — if multiple contexts use the same distro, they can be changed together.

## Step 3: Look Up Available Sizes

For each distro image family, look up what sizes exist in Evergreen:

```bash
evergreen list --distros | grep "<image_prefix>"
```

For example, if the variant uses `rhel8.8-arm64-m8g-8xlarge`:

```bash
evergreen list --distros | grep "rhel8.8-arm64-m8g"
```

This returns all available sizes for that image family. Build a **size ladder** — an ordered list from smallest to largest. See `references/size-ladder.md` for known progressions.

**CRITICAL RULE**: Only change the size suffix. Never change the image or OS. For example:
- `rhel8.8-arm64-m8g-8xlarge` → `rhel8.8-arm64-m8g-4xlarge` (correct)
- `rhel8.8-arm64-m8g-8xlarge` → `amazon2023-arm64-m8g-4xlarge` (WRONG — changed image)

## Step 4: Step Down Loop

For each unique distro in the inventory, starting with the largest:

### 4a. Step Down One Size

Move one step down the size ladder (e.g., `8xlarge` → `4xlarge`).

### 4b. Edit the YAML

Change ALL references to this distro in the variant (run_on, expansions, task-level overrides). Use find-and-replace within the variant block — don't change other variants.

### 4c. Validate

```bash
evergreen evaluate etc/evergreen.yml   # or .evergreen.yml depending on the repo
```

### 4d. Commit and Push

```bash
git add <changed_files>
git commit -m "<variant>: Step down <context> from <old_size> to <new_size>"
git push origin <branch>
```

### 4e. Submit Patch

```bash
evergreen patch -p <project-id> -v <variant-name> -t all -d "<variant>: test <new_size>" -f -y
```

### 4f. Monitor with pr-review-loop

Use the `pr-review-loop` skill workflow to monitor the patch. Specifically:

1. Use `/loop 10m` to check patch status
2. Wait for completion (green or failed)

### 4g. Analyze Result

**If GREEN**: This size works. Record it as the current smallest viable size. Go back to 4a and try stepping down again.

**If FAILED — check for OOM**:

Look for OOM indicators in the task logs (use `get_task_log_detailed`):
- "killed by the OOM killer"
- "signal 9" / "signal: killed"
- Core dumps being uploaded (`dump_mongod.*.core`)
- `fixture_setup` failures across multiple jobs (port/memory contention)
- Widespread cascading failures from a single resource issue

If OOM is confirmed:
1. Check if memory-optimized instances are available:
   ```bash
   evergreen list --distros | grep "r8g"   # or r7g, r6g
   ```
2. If r-series instances exist for this image family: try the r-series at the same or one step smaller vCPU count (r-series has 2x memory per vCPU vs m-series)
3. If r-series instances don't exist: ask the user if they want to request them from the infrastructure team for optimal cost
4. If the user declines: revert to the last working size

**If FAILED — not OOM** (test failures, timeouts that aren't resource-related):

The smaller instance likely introduced the failure. Revert to the last working size. This is fine — it means that instance size was already right-sized.

### 4h. Move to Next Distro

Once the current distro has been right-sized (either stepped down successfully or confirmed at its current size), move to the next unique distro in the inventory.

## Step 5: Final Report

After all instance sizes have been tested, print a summary:

```markdown
## Right-Sizing Report for <variant_name>

| Context | Original Distro | Final Distro | Change |
|---------|----------------|--------------|--------|
| run_on | rhel8.8-arm64-m8g-8xlarge | rhel8.8-arm64-m8g-4xlarge | Downsized |
| large_distro_name | rhel8.8-arm64-m8g-8xlarge | rhel8.8-arm64-m8g-4xlarge | Downsized |
| task: compile_TG | rhel8.8-arm64-m8g-8xlarge | rhel8.8-arm64-m8g-8xlarge | Already right-sized |

### Patches Submitted
- <patch_url_1> — 8xlarge → 4xlarge: GREEN
- <patch_url_2> — 4xlarge → 2xlarge: FAILED (OOM)
- <patch_url_3> — r8g-4xlarge: GREEN (memory-optimized alternative)
```

## When to Stop

- All instance sizes in the inventory have been tested at their smallest viable size
- User asks to stop
- 5 consecutive failed attempts on the same instance size (likely flaky, not size-related)

## Key References

- `references/size-ladder.md` — Known size progressions and how to look them up
- `pr-review-loop` skill — For patch monitoring and failure diagnosis
- `evergreen-cicd` skill — For patch submission patterns
- OOM indicators: "killed", "signal 9", core dumps, widespread fixture_setup failures