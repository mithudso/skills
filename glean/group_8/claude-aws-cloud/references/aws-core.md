<!-- hub-reference-banner -->
> **Reference file — part of the `aws-cloud` hub.** Formerly the standalone `aws-core` skill.
> Sibling topics in this family are now reference files under the hubs (`aws-cloud`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: aws-core
description: >
  AWS core infrastructure expert: IAM (policies, roles, trust relationships, SCPs, OIDC federation,
  permissions boundaries, Access Analyzer), AWS CLI v2/SSO, boto3, @aws-sdk (JS/TS), EC2 instance
  selection (Graviton, spot/on-demand, Auto Scaling), EBS volume types, S3 storage classes and
  lifecycle rules, S3 encryption and presigned URLs, VPC design (subnets, NAT, Transit Gateway,
  PrivateLink, security groups vs NACLs), Well-Architected Framework reviews (all 6 pillars),
  CloudWatch alarms/composite alarms/anomaly detection, and AWS troubleshooting workflows
  (CloudTrail, Logs Insights, X-Ray, AWS Health).
  TRIGGER: user asks about IAM policies, roles, SCPs, permissions boundaries, AWS CLI/SSO profiles,
  boto3, aws-sdk, EC2 instance types, EBS volumes, S3 lifecycle or encryption, VPC design, NAT
  gateway, Transit Gateway, security groups, NACLs, Well-Architected review, CloudWatch alarms,
  composite alarms, CloudTrail, X-Ray, or general AWS infrastructure troubleshooting.
  SKIP: Lambda, EventBridge, Step Functions, ECS, EKS, CDK, SAM, or API Gateway questions — use
  aws-serverless. Bedrock, SageMaker, or foundation model questions — use aws-ai-ml.
version: "1.1.0"
updated: "2026-05-29"
category: developer
tags:
  - aws
  - iam
  - ec2
  - s3
  - vpc
  - cloudwatch
  - well-architected
  - devops
  - cloud
keywords:
  - aws
  - iam
  - ec2
  - s3
  - vpc
  - cloudwatch
  - boto3
  - aws-sdk
  - aws-cli
  - well-architected
  - cloudtrail
  - eventbridge
  - transit-gateway
  - spot-instances
  - graviton
  - presigned-url
  - credential-chain
  - access-analyzer
  - permissions-boundary
  - scp
  - imdsv2
  - sse-kms
  - composite-alarm
whenToUse:
  - Designing or auditing IAM policies, roles, trust relationships, SCPs, or OIDC federation
  - Configuring AWS CLI profiles, SSO, credential chains, or SDK retry/error handling
  - Writing boto3 (Python) or @aws-sdk/client-* (JavaScript/TypeScript) integration code
  - Choosing EC2 instance types, EBS volume types, spot/on-demand mix, or placement groups
  - Designing S3 bucket policies, lifecycle rules, encryption, or presigned URL flows
  - Architecting VPCs: subnets, NAT gateways, Transit Gateway, PrivateLink, security groups vs NACLs
  - Running or preparing for a Well-Architected Framework review across any of the 6 pillars
  - Setting up CloudWatch alarms, composite alarms, anomaly detection, SNS, or EventBridge rules
  - Troubleshooting AWS issues using CloudTrail, CloudWatch Logs Insights, X-Ray, or AWS Health
whenNotToUse:
  - Lambda, EventBridge, Step Functions, ECS, EKS, CDK, SAM, or API Gateway — use aws-serverless
  - Bedrock, SageMaker, or any AWS AI/ML service — use aws-ai-ml
  - Redshift, Athena, Glue, or data warehouse services with no IAM/networking question
related_skills:
  - aws-serverless
  - aws-ai-ml
  - software-engineering-patterns
---

# AWS Core Services

Practical reference for AWS engineers — IAM, CLI/SDK, EC2, S3, networking, Well-Architected,
troubleshooting, and alerting patterns.

**Scope boundary:** This skill covers AWS core infrastructure services. For serverless (Lambda, API Gateway, SAM/CDK), use `aws-serverless`. For AI/ML (Bedrock, SageMaker), use `aws-ai-ml`.

**Output guidance:** When activated, produce the most actionable artifact for the task: a policy JSON snippet, a CLI command, a decision table, an architecture description, or a step-by-step troubleshooting sequence. Always flag security risks (e.g., overly permissive policies) before presenting a solution.

**Safety guardrail:** Never generate IAM policies with `"Action": "*"` or `"Resource": "*"` without explicitly calling out the security risk and recommending a scoped alternative. Always ask the user for their AWS region and account ID before generating resource ARNs that contain region or account values.

---

## Quick Reference Tables

### EC2 Instance Family Selection

| Workload                            | Family    | Example Type   | Notes                                        |
|-------------------------------------|-----------|----------------|----------------------------------------------|
| General web / microservices         | M8g       | m8g.xlarge     | Graviton4; best price-perf for most workloads |
| General web (prev-gen)              | M7g       | m7g.xlarge     | Graviton3; solid choice if M8g unavailable   |
| Compute-intensive (CI, HPC, render) | C8g / C7g | c8g.2xlarge    | Graviton4/3; 60% better energy efficiency    |
| In-memory / large datasets          | R8g / R7g | r8g.4xlarge    | Graviton4/3; up to 768 GiB RAM               |
| ML inference / GPU                  | G5 / G4dn | g5.xlarge      | NVIDIA A10G (G5) / T4 (G4dn)                 |
| Burst workloads (dev/test)          | T4g       | t4g.medium     | Graviton2; cheapest per vCPU                 |
| Storage-optimized (OLAP, Kafka)     | I8g / I4i | i8g.xlarge     | NVMe SSD; up to 30 TB local                  |
| High-throughput networking          | C7gn      | c7gn.16xlarge  | 200 Gbps network                             |

**Graviton generation map:** Graviton2 → T4g, M6g, C6g | Graviton3 → M7g, C7g, R7g | Graviton4 → M8g, C8g, R8g, I8g

**Decision rule:** Default to Graviton (Arm64) unless you have an x86-only binary dependency. Most modern runtimes (Node 18+, Python 3.9+, Java 17+, .NET 6+, Go 1.21+) run on ARM natively. Verify AMI availability in the target region before committing.

---

### EBS Volume Type Selection

| Type              | Use Case                          | Max IOPS  | Max Throughput | Notes                                  |
|-------------------|-----------------------------------|-----------|----------------|----------------------------------------|
| gp3               | Default for most workloads        | 16,000    | 1,000 MB/s     | IOPS & throughput configurable, independent of size |
| io2 Block Express | High-performance DB (Oracle, SAP) | 256,000   | 4,000 MB/s     | Multi-attach supported; prefer over io1 |
| st1               | Sequential big data, Kafka        | 500       | 500 MB/s       | Throughput-optimized HDD               |
| sc1               | Cold archives, infrequent access  | 250       | 250 MB/s       | Cheapest magnetic; lowest cost option  |
| io1               | Legacy high-IOPS workloads        | 64,000    | 1,000 MB/s     | Prefer io2 for new deployments         |

---

### S3 Storage Class Selection

| Class                          | Min Duration | Retrieval  | Use Case                               |
|--------------------------------|--------------|------------|----------------------------------------|
| S3 Standard                    | None         | Instant    | Hot data, frequent access              |
| S3 Intelligent-Tiering         | None         | Instant    | Unknown/variable access pattern; auto-tiers with no retrieval fee |
| S3 Standard-IA                 | 30 days      | Instant    | < once/month access; disaster recovery |
| S3 One Zone-IA                 | 30 days      | Instant    | Reproducible data, single-AZ OK        |
| S3 Glacier Instant Retrieval   | 90 days      | Instant    | Archives accessed ~once/quarter        |
| S3 Glacier Flexible Retrieval  | 90 days      | 1–12 hrs   | Compliance archives                    |
| S3 Glacier Deep Archive        | 180 days     | 12–48 hrs  | 7–10 year regulatory retention         |

**Lifecycle transition order (only forward — cannot reverse):**
Standard → Standard-IA → One Zone-IA → Glacier Instant → Glacier Flexible → Deep Archive

Note: S3 Intelligent-Tiering is configured as a storage class at object creation or via a separate lifecycle rule — it is not a step in the sequential transition chain above.

---

### VPC Design Checklist

| Layer            | Guideline                                                                        |
|------------------|----------------------------------------------------------------------------------|
| CIDR sizing      | /16 per VPC; /26 public subnets (~59 IPs); /21 private subnets (~2,043 IPs)      |
| AZ coverage      | Minimum 2 AZs; 3 for production                                                  |
| Internet access  | IGW + public subnets for load balancers; one NAT Gateway per AZ for private egress |
| S3/DynamoDB      | Deploy Gateway Endpoints (free); avoids NAT Gateway data processing charges      |
| Service mesh     | Interface Endpoints for ECR, SSM, KMS, Secrets Manager — keeps traffic off internet |
| Cross-VPC        | VPC Peering for ≤ 10 VPCs (non-transitive); Transit Gateway for hub-and-spoke    |
| Segmentation     | Security Groups for instance-level stateful control; NACLs only for subnet-wide blocklisting |
| DNS              | Enable DNS hostnames + resolution; use Route 53 Private Hosted Zones             |
| Bastion hosts    | Replace with Systems Manager Session Manager (no inbound port 22 required)       |

---

### IAM Policy Evaluation Order

When an IAM principal makes a request, AWS evaluates policies in this order. An explicit `Deny` at any layer wins immediately.

1. **Explicit DENY** (any policy type) → request denied immediately
2. **SCP** (AWS Organizations) — if SCP does not allow the action → implicit DENY, even if IAM policy allows it
3. **Resource Control Policy (RCP)** — org-level guardrail on resource access (newer than SCP; applies to resource side)
4. **Permissions Boundary** — limits the *maximum* permissions an identity-based policy can grant; does not grant permissions by itself
5. **Identity-based policy** (inline or managed attached to user/role) — must have explicit ALLOW
6. **Resource-based policy** (S3 bucket policy, SQS policy, etc.) — same-account: UNION with identity-based; cross-account: *both* must allow
7. **Session policy** (passed at `sts:AssumeRole`) — further narrows, cannot expand permissions

**Effective permission:** `(identity-based ∩ permissions-boundary)` UNION `resource-based` — after passing SCP/RCP layer.
**Cross-account note:** For cross-account access, both the resource-based policy in the target account AND the identity-based policy in the source account must explicitly allow the action.

---

### AWS Credential Provider Chain (CLI + boto3 + @aws-sdk)

When no explicit credentials are provided, SDKs check in this order (first match wins):

1. **Environment variables** — `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` / `AWS_SESSION_TOKEN`
2. **Shared credentials file** — `~/.aws/credentials` (profile selected by `AWS_PROFILE` or `--profile`)
3. **AWS SSO** — IAM Identity Center token cache (`~/.aws/sso/cache/`)
4. **AWS config file** — `~/.aws/config` credential_process or assume-role chains
5. **Container credentials** — `AWS_CONTAINER_CREDENTIALS_RELATIVE_URI` (ECS task role)
6. **Instance metadata** — IMDSv2 (EC2 instance role; requires `PUT` token fetch, not IMDSv1 `GET`)

**Debugging tip:** `aws sts get-caller-identity --profile <profile>` confirms which identity is resolved. `AWS_SDK_LOAD_CONFIG=1` forces the JS SDK to honor `~/.aws/config`.

---

### CloudWatch Alarm States

| State             | Meaning                                              |
|-------------------|------------------------------------------------------|
| OK                | Metric within threshold for all evaluation periods   |
| ALARM             | Metric breached threshold for N consecutive periods  |
| INSUFFICIENT_DATA | Alarm just created, metric gap, or data too sparse   |

Use **composite alarms** to reduce alert noise: combine child metric alarms with AND/OR logic into a parent composite alarm → single SNS notification on genuine degradation.

---

## Key Patterns

### S3 Encryption Options

| Method        | Key management          | Use case                                              | Notes                                        |
|---------------|-------------------------|-------------------------------------------------------|----------------------------------------------|
| SSE-S3        | AWS-managed (AES-256)   | Default; no compliance requirement for key control    | Enabled by default on all new buckets (2023+)|
| SSE-KMS       | AWS KMS CMK             | Audit trail per-object; cross-account access control  | KMS API calls add cost + latency             |
| DSSE-KMS      | AWS KMS CMK (dual-layer)| Regulated industries requiring dual encryption layer  | Two independent AES-256 encryption layers    |
| SSE-C         | Customer-provided key   | Customer retains full key material outside AWS        | Disabled by default for new buckets (2026+); key must be sent with every request |
| CSE           | Client-side (pre-upload)| Maximum control; data encrypted before leaving client | Requires client-side SDK; AWS sees only ciphertext |

**Decision rule:** Start with SSE-S3 (free, zero config). Upgrade to SSE-KMS when you need per-object CloudTrail key-usage logs or cross-account key policies. Use DSSE-KMS only when a compliance mandate explicitly requires dual-layer encryption.

---

### EC2 Auto Scaling: Mixed On-Demand + Spot

```json
{
  "MixedInstancesPolicy": {
    "InstancesDistribution": {
      "OnDemandBaseCapacity": 2,
      "OnDemandPercentageAboveBaseCapacity": 30,
      "SpotAllocationStrategy": "price-capacity-optimized"
    },
    "LaunchTemplate": {
      "LaunchTemplateSpecification": {
        "LaunchTemplateId": "lt-0abcdef1234567890",
        "Version": "$Latest"
      },
      "Overrides": [
        {"InstanceType": "c8g.xlarge"},
        {"InstanceType": "c7g.xlarge"},
        {"InstanceType": "c6g.xlarge"}
      ]
    }
  }
}
```

**Strategy guide:**
- `price-capacity-optimized` — recommended default; picks pools least likely to be interrupted at lowest price
- `capacity-optimized` — prioritises availability over price; best for stateful/hard-to-interrupt jobs
- `lowest-price` — cheapest but highest interruption risk; suitable only for fault-tolerant batch
- `OnDemandBaseCapacity` — minimum number of On-Demand instances always running (set ≥ 1 for critical workloads)
- Specify 3–6 instance type overrides per family so the ASG can fall back when a Spot pool is unavailable

---

### IAM: Least-Privilege Service-to-Service Policy

**Don't do this** (wildcard resource — overly permissive):
```json
{
  "Effect": "Allow",
  "Action": "s3:*",
  "Resource": "*"
}
```

**Do this** (scoped action + scoped resource + condition):
```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Sid": "ReadTargetBucketOnly",
    "Effect": "Allow",
    "Action": ["s3:GetObject", "s3:ListBucket"],
    "Resource": [
      "arn:aws:s3:::my-bucket",
      "arn:aws:s3:::my-bucket/*"
    ],
    "Condition": {
      "StringEquals": { "aws:RequestedRegion": "us-east-1" }
    }
  }]
}
```

**Trust policy for EC2 role:**
```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": { "Service": "ec2.amazonaws.com" },
    "Action": "sts:AssumeRole"
  }]
}
```
Never use `"Principal": "*"` in a trust policy — it allows any identity in any account to assume the role.

---

### AWS CLI v2: SSO Setup (IAM Identity Center)

```bash
# One-time setup — creates [sso-session] block in ~/.aws/config (v2.22+)
aws configure sso

# Authenticate (opens browser)
aws sso login --sso-session my-org

# Use a named profile
aws s3 ls --profile prod-readonly
AWS_PROFILE=prod-readonly aws ec2 describe-instances

# ~/.aws/config (modern sso-session format):
# [sso-session my-org]
# sso_start_url  = https://my-org.awsapps.com/start
# sso_region     = us-east-1
# sso_registration_scopes = sso:account:access
#
# [profile prod-readonly]
# sso_session    = my-org
# sso_account_id = 123456789012
# sso_role_name  = ReadOnly
# region         = us-east-1
```

---

### boto3: Retry Configuration and Error Handling

```python
import boto3
from botocore.config import Config
from botocore.exceptions import ClientError

config = Config(
    retries={
        "mode": "adaptive",       # adaptive | standard | legacy
        "total_max_attempts": 5,  # includes the initial attempt
    },
    connect_timeout=5,
    read_timeout=30,
)

s3 = boto3.client("s3", config=config)

try:
    response = s3.get_object(Bucket="my-bucket", Key="path/to/obj")
    data = response["Body"].read()
except ClientError as e:
    code = e.response["Error"]["Code"]
    if code == "NoSuchKey":
        pass  # handle missing object — return None, raise custom exception, etc.
    elif code in ("AccessDenied", "403"):
        pass  # handle permission error — check IAM policy, bucket policy
    else:
        raise  # re-raise unexpected errors; don't swallow silently
```

Use `"mode": "adaptive"` for services with heavy throttling (DynamoDB, Lambda invocations, API Gateway). `"mode": "standard"` suits most other cases.

---

### @aws-sdk/client-* (JavaScript/TypeScript)

```typescript
import { S3Client, GetObjectCommand, NoSuchKey } from "@aws-sdk/client-s3";
import { fromSSO } from "@aws-sdk/credential-providers";

const client = new S3Client({
  region: "us-east-1",
  credentials: fromSSO({ profile: "prod-readonly" }), // or fromIni(), fromEnv(), fromInstanceMetadata()
  maxAttempts: 5,   // built-in retry; default is 3
});

try {
  const response = await client.send(new GetObjectCommand({
    Bucket: "my-bucket",
    Key: "path/to/obj",
  }));
  const body = await response.Body?.transformToString();
} catch (err) {
  if (err instanceof NoSuchKey) {
    // handle missing object
  } else {
    throw err;
  }
}
```

**Credential providers (JS SDK v3):** `fromEnv()` → `fromIni()` → `fromSSO()` → `fromInstanceMetadata()` — matches the CLI chain. Import from `@aws-sdk/credential-providers`.

---

### S3: Presigned URL Generation

```python
import boto3

s3 = boto3.client("s3", region_name="us-east-1")

# Generate a time-limited GET URL (no AWS credentials needed by the caller)
url = s3.generate_presigned_url(
    "get_object",
    Params={"Bucket": "my-bucket", "Key": "reports/2026-q1.pdf"},
    ExpiresIn=3600,   # seconds; max 7 days for IAM role credentials
)

# PUT URL (allow direct upload without exposing bucket credentials)
put_url = s3.generate_presigned_url(
    "put_object",
    Params={"Bucket": "my-bucket", "Key": "uploads/user-file.csv", "ContentType": "text/csv"},
    ExpiresIn=900,
)
```

**Notes:** Presigned URLs inherit the signer's permissions at generation time, not call time. If the signing role is revoked before the URL expires, the URL still works until expiry (unless you use presigned URL with STS temporary credentials, which expire with the session).

---

### CloudWatch Composite Alarm (CloudFormation)

```yaml
CPUAlarm:
  Type: AWS::CloudWatch::Alarm
  Properties:
    MetricName: CPUUtilization
    Namespace: AWS/EC2
    Dimensions:
      - Name: AutoScalingGroupName
        Value: !Ref MyASG
    Threshold: 80
    ComparisonOperator: GreaterThanThreshold
    EvaluationPeriods: 3
    Period: 60
    Statistic: Average
    TreatMissingData: notBreaching

MemoryAlarm:
  Type: AWS::CloudWatch::Alarm
  Properties:
    MetricName: mem_used_percent
    Namespace: CWAgent
    Threshold: 85
    ComparisonOperator: GreaterThanThreshold
    EvaluationPeriods: 3
    Period: 60
    Statistic: Average
    TreatMissingData: notBreaching

AppDegradedComposite:
  Type: AWS::CloudWatch::CompositeAlarm
  Properties:
    AlarmRule: !Sub "ALARM(${CPUAlarm}) OR ALARM(${MemoryAlarm})"
    AlarmActions:
      - !Ref AlertSNSTopic
    TreatMissingData: notBreaching
```

**Best practice:** Set `TreatMissingData: notBreaching` on metric alarms feeding a composite alarm, so a brief metrics gap doesn't trigger a false ALARM state.

---

## Well-Architected Pillar Summary

| Pillar                 | Core Question                                | Key AWS Tools                                  | Common Finding                         |
|------------------------|----------------------------------------------|------------------------------------------------|----------------------------------------|
| Operational Excellence | Can we run + evolve operations efficiently?  | CloudFormation, Systems Manager, Ops Center    | Manual runbooks not automated          |
| Security               | Are we protecting data + systems?            | IAM, KMS, GuardDuty, Security Hub, Macie       | Overly permissive IAM; missing MFA     |
| Reliability            | Can we recover from failures automatically?  | Multi-AZ, Route 53 health checks, AWS Backup   | Single-AZ deployments; no tested DR    |
| Performance Efficiency | Are we using compute efficiently?            | Graviton, CloudFront, ElastiCache, Auto Scaling | Oversized on-demand instances          |
| Cost Optimization      | Are we spending only what's needed?          | Cost Explorer, Savings Plans, Spot, Compute Optimizer | Idle resources; no Savings Plans  |
| Sustainability         | Are we minimizing environmental impact?      | Graviton, managed services, right-sizing, Spot | x86 instances where Graviton fits      |

**Review cadence:** Run a Well-Architected Review every 3–6 months, or before any major architectural change. Use the AWS Well-Architected Tool in the console to track findings and improvement plans.

---

## Operational Patterns

### IAM Access Analyzer

Use IAM Access Analyzer to detect external exposure and generate least-privilege policies from CloudTrail activity:

```bash
# Find externally accessible resources in your account
aws accessanalyzer list-findings \
  --analyzer-arn arn:aws:access-analyzer:us-east-1:123456789012:analyzer/MyAnalyzer

# Start least-privilege policy generation (uses last 90 days of CloudTrail)
# Step 1: kick off generation job
aws accessanalyzer start-policy-generation \
  --policy-generation-details '{"principalArn":"arn:aws:iam::123456789012:role/MyRole"}'

# Step 2: retrieve the generated policy (use jobId from step 1 output)
aws accessanalyzer get-generated-policy --job-id <jobId>

# Validate a policy document before attaching it
aws accessanalyzer validate-policy \
  --policy-document file://policy.json \
  --policy-type IDENTITY_POLICY
```

**Key use cases:**
- Detect S3 buckets, IAM roles, KMS keys, SQS queues, and Lambda functions exposed to the internet or other accounts
- Generate scoped IAM policies from real CloudTrail usage rather than hand-authoring from scratch
- Validate policies before deployment to catch overly permissive statements early

**Note:** All examples use `us-east-1` and account `123456789012` as placeholders — substitute your actual region and account ID.

---

## Troubleshooting Workflow

When an AWS action fails or behaves unexpectedly, follow this sequence. Start at step 1 for all failures; jump to step 4 for latency or distributed-tracing issues; jump to step 5 to rule out a regional AWS outage first.

1. **CloudTrail** — Find the API call: `aws cloudtrail lookup-events --lookup-attributes AttributeKey=EventName,AttributeValue=<Action>`. Check `errorCode` and `errorMessage` in the event.
2. **IAM Policy Simulator** — Test whether a specific principal can perform an action: console at `https://policysim.aws.amazon.com` or `aws iam simulate-principal-policy`.
3. **CloudWatch Logs Insights** — Query application logs:
   ```
   fields @timestamp, @message
   | filter @message like /ERROR/
   | sort @timestamp desc
   | limit 50
   ```
4. **X-Ray ServiceLens** — For distributed tracing: view service map, identify latency outliers, correlate trace IDs with CloudWatch log entries using `filter @message like /<trace-id>`.
5. **AWS Health** — Check for service disruptions: `aws health describe-events --filter '{"eventStatusCodes":["open"]}'`.
6. **Systems Manager Session Manager** — Access instances without SSH: `aws ssm start-session --target <instance-id>`. No inbound security group rules required.
