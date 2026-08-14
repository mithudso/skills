# MMS — Confluence Reading Context

## Common Spaces
- `MMS` — Atlas/MMS engineering docs
- `CLOUDP` — Atlas product and architecture
- `HELP` — support runbooks and on-call guides

## Frequently Referenced Pages
- MMS local dev setup: search `run mms local` in MMS space
- Atlas architecture overview: search `architecture` in CLOUDP space
- On-call runbooks: HELP space → search by alert name

## CQL Examples
```
space = "MMS" AND title ~ "local development"
space = "CLOUDP" AND ancestor = <parent-page-id>
space = "HELP" AND label = "runbook"
```
