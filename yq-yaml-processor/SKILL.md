---
name: yq-yaml-processor
description: >-
  Uses yq for structurally parsing, querying, and updating YAML, XML, and CSV files safely without breaking formatting.
---

# yq Instructions

When the user asks you to modify or query complex YAML files (such as Kubernetes manifests, GitHub Actions workflows, or docker-compose files), you MUST use `yq` instead of `sed`, `awk`, or manual string replacement tools. This ensures you do not break the YAML structure.

## Querying YAML
```bash
# Read a specific property
yq '.metadata.name' config.yaml

# Read an array element
yq '.services.web.ports[0]' docker-compose.yaml
```

## Modifying YAML (In-Place)
Always use the `-i` flag for in-place modifications.
```bash
# Update an image tag in a Kubernetes deployment
yq -i '.spec.template.spec.containers[0].image = "nginx:1.21"' deployment.yaml

# Add a new environment variable to a service
yq -i '.services.api.environment += {"NEW_VAR": "value"}' docker-compose.yml
```

## Format Conversion
`yq` can also seamlessly convert between JSON, YAML, and XML.
```bash
# Convert YAML to JSON
yq -o=json '.' file.yaml

# Convert JSON to YAML
yq -p=json -o=yaml '.' file.json
```
