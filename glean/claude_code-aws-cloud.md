# aws-cloud

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude Code
**Original Path:** claude-code/aws-cloud

## Description
AWS hub — infra, serverless, AI/ML, private connectivity, databases. TRIGGER: IAM (policies, roles, SCPs, federation), AWS CLI/SSO, boto3/aws-sdk, EC2/EBS/S3, VPC design, Well-Architected, CloudWatch/CloudTrail (aws-core); Lambda, EventBridge, Step Functions, API Gateway, ECS/EKS Fargate, CDK/SAM/CloudFormation, serverless patterns (aws-serverless); Bedrock & SageMaker — foundation models, agents, RAG, guardrails, training/endpoints (aws-ai-ml); PrivateLink & VPC endpoints — interface/gateway/GWLB endpoints, endpoint services, private DNS/split-horizon, endpoint policies, PrivateLink vs peering vs Transit Gateway, endpoint troubleshooting (aws-privatelink-vpc-endpoints); AWS databases — RDS/Aurora/DynamoDB/DocumentDB/Neptune, DocumentDB-vs-Atlas, CockroachDB, IndexedDB (databases-aws-cockroach-indexeddb). SKIP: non-AWS cloud; MongoDB Atlas platform (incl. Atlas private endpoints) → mongodb-atlas-expert; Cloudflare R2 / R2-vs-S3 egress economics → cloudflare-platform.

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
| AWS databases | RDS/Aurora/DynamoDB/DocumentDB/Neptune, DocumentDB-vs-Atlas, CockroachDB, IndexedDB | `references/databases-aws-cockroach-indexeddb.md` |

<!-- cross-hub-map -->
## Cross-hub map — where every aws-cloud topic lives

Family split across hubs. Deep material not in Sub-skill routing table → sibling hub below — **activate hub or `Read` its `references/<name>.md`**. All former standalone skills now references under one hub (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `aws-cloud` | aws-cloud | `references/aws-core.md`, `references/aws-serverless.md`, `references/aws-ai-ml.md`, `references/aws-privatelink-vpc-endpoints.md`, `references/databases-aws-cockroach-indexeddb.md` |
| `mongodb-atlas-expert` | Atlas-consumer private connectivity | `references/mongodb-aws-networking.md`, `references/mongodb-atlas-azure.md`, `references/mongodb-atlas-gcp.md` |
| `networking` | DNS protocol internals, split-horizon | `references/dns-deep-dive.md` |