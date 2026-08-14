# MMS — Confluence Editing Context

## Wiki Base URL
`https://wiki.corp.mongodb.com`

## Common Spaces
- `MMS` — Atlas/MMS engineering docs
- `CLOUDP` — Atlas product and architecture docs
- `HELP` — support runbooks

## Frequently Edited Page Types
- Feature wikis: `wiki.corp.mongodb.com/spaces/MMS/pages/...`
- Runbooks: `wiki.corp.mongodb.com/spaces/HELP/pages/...`
- Architecture docs: `wiki.corp.mongodb.com/spaces/CLOUDP/pages/...`

## Notes
Most MMS wiki pages use the markdown macro format. Always confirm storage format
before editing — use `confluence_get_page` with `bodyType: storage` first.
