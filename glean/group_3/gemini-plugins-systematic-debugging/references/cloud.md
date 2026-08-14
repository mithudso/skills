# MongoDB Atlas / Cloud Debugging Context

When debugging in Atlas or Cloud Manager repos:

- Check Atlas operator logs: `kubectl logs -n mongodb <pod-name>`
- Evergreen project: `mongodb-mongo-master-atlas`
- Atlas CLI debugging: `atlas logs download <cluster> --type mongodb`
- Common failure patterns: Atlas operator reconciliation loops, MDB resource status errors
