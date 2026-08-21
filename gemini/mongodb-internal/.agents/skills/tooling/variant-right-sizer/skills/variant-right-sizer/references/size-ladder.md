# Instance Size Ladder Reference

## Size Progression (smallest → largest)

```
small < medium < large < xlarge < 2xlarge < 4xlarge < 8xlarge < 12xlarge < 16xlarge < 24xlarge < 48xlarge
```

Not all sizes exist for every image family. Always verify with:

```bash
evergreen list --distros | grep "<image_prefix>"
```

## Common Image Families

### ARM64 (Graviton m8g)

```bash
evergreen list --distros | grep "m8g"
```

Known sizes:
- `*-m8g-xlarge` — 4 vCPU, 16 GB
- `*-m8g-2xlarge` — 8 vCPU, 32 GB
- `*-m8g-4xlarge` — 16 vCPU, 64 GB
- `*-m8g-8xlarge` — 32 vCPU, 128 GB
- `*-m8g-12xlarge` — 48 vCPU, 192 GB
- `*-m8g-16xlarge` — 64 vCPU, 256 GB
- `*-m8g-24xlarge` — 96 vCPU, 384 GB
- `*-m8g-48xlarge` — 192 vCPU, 768 GB

### ARM64 (Graviton r8g — memory-optimized)

r8g instances have 2x the memory per vCPU compared to m8g:

```bash
evergreen list --distros | grep "r8g"
```

- `*-r8g-xlarge` — 4 vCPU, 32 GB
- `*-r8g-2xlarge` — 8 vCPU, 64 GB
- `*-r8g-4xlarge` — 16 vCPU, 128 GB
- `*-r8g-8xlarge` — 32 vCPU, 256 GB

Use r8g when OOM is the failure mode but you don't want more CPUs (which increases parallelism and can make things worse).

### x86 Generic Sizes

```bash
evergreen list --distros | grep "rhel8.8-"
```

Known sizes:
- `rhel8.8-small` — ~2 vCPU, 8 GB
- `rhel8.8-large` — ~4 vCPU, 16 GB
- `rhel8.8-xlarge` — ~8 vCPU, 32 GB
- `rhel8.8-xxlarge` — ~16 vCPU, 64 GB

### Amazon Linux 2023

```bash
evergreen list --distros | grep "amazon2023"
```

- `amazon2023-cloud-small`
- `amazon2023-cloud-medium`
- `amazon2023-cloud-large`

## How to Extract the Image Prefix

Given a distro like `rhel8.8-arm64-m8g-8xlarge`:
1. Strip the size suffix: `rhel8.8-arm64-m8g`
2. Search: `evergreen list --distros | grep "rhel8.8-arm64-m8g"`
3. The results show all available sizes for that image family

For generic distros like `rhel8.8-large`:
1. Strip the size: `rhel8.8`
2. Search: `evergreen list --distros | grep "^rhel8.8-"`
3. Filter to the same naming pattern

## Cost Considerations

- Stepping down from 8xlarge to 4xlarge roughly halves compute cost
- r-series instances cost ~10-15% more than m-series at the same size, but if they let you use a smaller size, the net savings can be significant
- The goal is the smallest size where tests pass reliably — not the smallest possible size
