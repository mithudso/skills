# aws-cloud

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude
**Original Path:** claude/standalone/aws-cloud

## Description
AWS hub — infra, serverless, AI/ML, private connectivity, databases, security, identity. TRIGGER: IAM, CLI/SSO, boto3/aws-sdk, EC2/EBS/S3, VPC, Well-Architected, CloudWatch/CloudTrail (aws-core); Lambda, EventBridge, Step Functions, API Gateway, ECS/EKS Fargate, CDK/SAM/CloudFormation (aws-serverless); Bedrock & SageMaker — models, agents, RAG, guardrails (aws-ai-ml); PrivateLink & VPC endpoints — interface/gateway/GWLB, endpoint services, private DNS/split-horizon, endpoint policies, vs peering/TGW, troubleshooting (aws-privatelink-vpc-endpoints); AWS databases — RDS/Aurora/DynamoDB/DocumentDB/Neptune, DocumentDB-vs-Atlas, CockroachDB, IndexedDB (databases-aws-cockroach-indexeddb); Aurora engine deep-dive — storage-compute separation, 6-way quorum, Serverless v2, Global Database, Aurora DSQL (aws-aurora); security services — KMS, GuardDuty, Security Hub, Config, WAF/Shield, Macie, Inspector, Access Analyzer (aws-security-services); secretless workload identity — IMDSv2, ECS task roles, EKS IRSA & Pod Identity, GitHub/GitLab OIDC, Roles Anywhere (aws-secretless-auth). SKIP: non-AWS cloud; MongoDB Atlas platform → mongodb-atlas-expert; the PostgreSQL engine itself (MVCC, vacuum, planner, WAL) → postgresql-expert; Cloudflare R2/egress → cloudflare-platform; app/web code-level security review → security-review.

---

# aws-cloud

AWS hub — core infra, serverless, AI/ML, database ecosystem.

Routes to `references/` on demand. See each spoke for depth.

## Sub-skill routing table

| Sub-topic | When to load | Reference |
| --- | --- | --- |
| AWS core infrastructure | IAM, CLI/SSO, SDKs, EC2/EBS/S3, VPC design basics, Well-Architected, CloudWatch/CloudTrail | `references/aws-core.md` |
| Serverless & containers | Lambda, EventBridge, Step Functions, API Gateway, ECS/EKS Fargate, CDK/SAM/CloudFormation | `references/aws-serverless.md` |
| AI/ML services | Bedrock, SageMaker, foundation models, agents, RAG, guardrails | `references/aws-ai-ml.md` |
| PrivateLink & VPC endpoints | interface/gateway/GWLB/resource endpoints, endpoint services (provider side, acceptance, cross-account/region), private DNS & split-horizon, endpoint policies, SGs on endpoint ENIs, quotas/pricing, PrivateLink vs peering vs TGW, endpoint connectivity troubleshooting | `references/aws-privatelink-vpc-endpoints.md` |
| AWS databases | RDS/Aurora/DynamoDB/DocumentDB/Neptune, DocumentDB-vs-Atlas, CockroachDB, IndexedDB — the cross-engine "which AWS database?" decision guide | `references/databases-aws-cockroach-indexeddb.md` |
| Aurora engine deep-dive | Aurora storage-compute separation, distributed log-structured storage, 6-way quorum across 3 AZs, "the log is the database", reader endpoints/replicas, fast failover, Serverless v2 (ACUs, scale-to-zero), Global Database & write forwarding, Blue/Green, I/O-Optimized, Backtrack, Babelfish, and Aurora DSQL (active-active distributed PostgreSQL-compatible) | `references/aws-aurora.md` |
| Security services | KMS (CMK, key policy vs grants, envelope encryption, multi-Region keys), GuardDuty, Security Hub (ASFF, standards), AWS Config (rules, conformance packs, remediation), CloudTrail-as-security, WAF & Shield, Macie, Inspector, Detective, IAM Access Analyzer, Network Firewall, ACM, multi-account delegated admin | `references/aws-security-services.md` |
| Secretless workload identity | credential-provider chain, EC2 instance profiles & IMDSv2 (hop limit, SSRF), ECS task roles, EKS IRSA vs Pod Identity, Lambda execution roles, IAM Roles Anywhere (X.509), GitHub/GitLab OIDC federation, cross-account assumption & external IDs, IAM Identity Center, eliminating static keys | `references/aws-secretless-auth.md` |

<!-- cross-hub-map -->
## Cross-hub map — where every aws-cloud topic lives

Family split across hubs. Deep material not in Sub-skill routing table → sibling hub below — **activate hub or `Read` its `references/<name>.md`**. All former standalone skills now references under one hub (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `aws-cloud` | aws-cloud | `references/aws-core.md`, `references/aws-serverless.md`, `references/aws-ai-ml.md`, `references/aws-privatelink-vpc-endpoints.md`, `references/databases-aws-cockroach-indexeddb.md`, `references/aws-aurora.md`, `references/aws-security-services.md`, `references/aws-secretless-auth.md` |
| `postgresql-expert` | PostgreSQL engine & data plane (MVCC, vacuum, planner, indexes, WAL/replication, partitioning, types/pgvector, isolation/RLS) | `references/mvcc-vacuum-and-bloat.md`, `references/query-planner-and-statistics.md`, `references/wal-replication-and-recovery.md` |
| `mongodb-atlas-expert` | Atlas-consumer private connectivity | `references/mongodb-aws-networking.md`, `references/mongodb-atlas-azure.md`, `references/mongodb-atlas-gcp.md` |
| `networking` | DNS protocol internals, split-horizon | `references/dns-deep-dive.md` |