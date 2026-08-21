<!-- hub-reference-banner -->
> **Reference file — part of the `devops-containers-cicd` hub.** Formerly the standalone `terraform-kafka-infra` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: terraform-kafka-infra
title: "Terraform & Apache Kafka Infrastructure"
description: >-
  Infrastructure-as-code with Terraform/OpenTofu and event-streaming with Apache Kafka.
  Covers HCL fundamentals, variable management, modules, remote state backends (S3/GCS/
  HCP Terraform), workspaces, state manipulation (import, moved, check blocks), Terragrunt,
  testing pyramid (terraform validate, tflint, checkov, native terraform test, Terratest),
  Kafka architecture (brokers, topics, partitions, consumer groups, KRaft mode), producer/
  consumer configuration, exactly-once semantics, tiered storage, Kafka Connect (Debezium
  CDC, SMTs, Schema Registry), Kafka Streams vs Flink, ksqlDB, Kafka security (SASL,
  TLS, ACLs), and managed services (Confluent Cloud, Amazon MSK, Azure Event Hubs).
  TRIGGER: writing, reviewing, or debugging Terraform HCL; designing or troubleshooting
  Kafka topics, partitions, or consumer groups; choosing between managed Kafka offerings;
  migrating existing infrastructure into Terraform state; setting up Kafka Connect
  pipelines or stream-processing jobs; securing Kafka clusters; deciding between
  Terraform and OpenTofu.
  SKIP: Kubernetes cluster configuration (use kubernetes-networking); CI/CD pipeline
  authoring (use cicd-pipelines); MongoDB Atlas infrastructure (use mongodb-atlas-iac
  or mongodb-atlas-terraform).
version: "1.2.0"
updated: "2026-05-29"
category: developer
tags:
  - terraform
  - kafka
  - infrastructure
  - iac
  - hcl
  - kafka-connect
  - kafka-streams
  - confluent
  - msk
  - opentofu
keywords:
  - Terraform
  - OpenTofu
  - HCL
  - terraform module
  - remote state
  - S3 backend
  - DynamoDB lock
  - terraform import
  - moved block
  - check block
  - Terragrunt
  - tflint
  - checkov
  - Terratest
  - terraform test
  - Apache Kafka
  - KRaft mode
  - Kafka broker
  - Kafka topic
  - consumer group
  - exactly-once semantics
  - Kafka Connect
  - Debezium
  - Schema Registry
  - Kafka Streams
  - Apache Flink
  - ksqlDB
  - SASL
  - mTLS
  - Confluent Cloud
  - Amazon MSK
  - Azure Event Hubs
  - tiered storage
whenToUse:
  - "Writing, reviewing, or debugging Terraform HCL configurations"
  - "Designing Kafka topics, partitions, and consumer group topology"
  - "Troubleshooting Kafka consumer lag or rebalance storms"
  - "Choosing between Confluent Cloud, Amazon MSK, or Azure Event Hubs"
  - "Migrating existing infrastructure into Terraform state"
  - "Setting up Kafka Connect pipelines (Debezium CDC, JDBC, S3 sink)"
  - "Securing Kafka clusters with SASL/SCRAM, TLS, or ACLs"
  - "Deciding between Terraform and OpenTofu"
  - "Writing terraform test or Terratest integration tests"
  - "Implementing exactly-once semantics in Kafka producers/consumers"
whenNotToUse:
  - "Kubernetes cluster configuration — use kubernetes-networking"
  - "CI/CD pipeline authoring — use cicd-pipelines"
  - "MongoDB Atlas infrastructure — use mongodb-atlas-iac or mongodb-atlas-terraform"
related_skills:
  - kubernetes-networking
  - cicd-pipelines
  - mongodb-atlas-iac
  - mongodb-atlas-terraform
  - mongodb-kafka-connector
---

# Terraform & Apache Kafka Infrastructure Skill

## Overview

This skill covers infrastructure-as-code with Terraform/HCL (and its open-source fork OpenTofu) and event-streaming architecture with Apache Kafka. It addresses provisioning cloud resources declaratively, managing infrastructure state, building reusable modules, wiring Kafka producers/consumers, and operating managed Kafka services at scale.

---

## Terraform Core

### HCL Fundamentals

```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  required_version = ">= 1.6"
}

provider "aws" {
  region = var.aws_region
}

resource "aws_s3_bucket" "data" {
  bucket = "${var.env}-${var.project}-data"
  tags   = local.common_tags
}

data "aws_vpc" "main" {
  filter {
    name   = "tag:Name"
    values = ["${var.env}-vpc"]
  }
}

locals {
  common_tags = {
    Environment = var.env
    Project     = var.project
    ManagedBy   = "terraform"
  }
}
```

### Variable Management

```hcl
variable "env" {
  type        = string
  description = "Deployment environment (dev|staging|prod)"
  validation {
    condition     = contains(["dev", "staging", "prod"], var.env)
    error_message = "env must be dev, staging, or prod."
  }
}

variable "db_password" {
  type      = string
  sensitive = true  # suppressed from logs and plan output
}
```

Variable precedence (highest to lowest):
1. `-var` flags and `-var-file` CLI arguments
2. `terraform.tfvars` / `*.auto.tfvars`
3. Environment variables (`TF_VAR_<name>`)
4. Default values in `variable` blocks

### Expressions and Functions

```hcl
# Dynamic blocks
resource "aws_security_group" "app" {
  name = "${var.env}-app"
  dynamic "ingress" {
    for_each = var.ingress_rules
    content {
      from_port   = ingress.value.port
      to_port     = ingress.value.port
      protocol    = "tcp"
      cidr_blocks = ingress.value.cidr_blocks
    }
  }
}

# for_each over a set
resource "aws_iam_user" "devs" {
  for_each = toset(var.developer_names)
  name     = each.value
}

# count with conditional
resource "aws_cloudwatch_alarm" "high_cpu" {
  count = var.enable_alerts ? 1 : 0
}
```

### Key Provider Ecosystems

| Provider | Registry Path | Primary Use |
|---|---|---|
| AWS | `hashicorp/aws` | EC2, S3, RDS, VPC, IAM, Lambda, EKS |
| GCP | `hashicorp/google` | GCE, GCS, GKE, Cloud SQL, BigQuery |
| Azure | `hashicorp/azurerm` | VMs, AKS, Blob Storage, Azure SQL |
| Kubernetes | `hashicorp/kubernetes` | Deployments, Services, ConfigMaps |
| Helm | `hashicorp/helm` | Helm chart releases on K8s |
| MongoDB Atlas | `mongodb/mongodbatlas` | Atlas clusters, databases, users |
| Confluent | `confluentinc/confluent` | Kafka topics, connectors, schemas |

---

## Terraform Modules & State

### Standard Module Structure

```
modules/
  eks-cluster/
    main.tf
    variables.tf
    outputs.tf
    versions.tf
    README.md
    examples/basic/main.tf
    tests/cluster_test.go

environments/
  dev/
    main.tf
    terraform.tfvars
    backend.tf
  staging/
  prod/
```

### Module Best Practices

- **One function per module**: networking, compute, IAM, database — not combined
- **Shallow nesting**: maximum 2-3 levels
- **Semantic versioning**: pin module sources to versions (`~> 2.0`)
- **No hardcoded values**: all configuration through variables with validation
- **Explicit outputs**: expose ARNs, IDs, endpoints other modules need
- **Avoid over-modularization**: wrapping a single resource adds friction without benefit

### Remote State Backends

**S3 + DynamoDB (AWS):**
```hcl
terraform {
  backend "s3" {
    bucket         = "my-company-tfstate"
    key            = "prod/eks/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    kms_key_id     = "arn:aws:kms:us-east-1:123456789:key/abc123"
    dynamodb_table = "terraform-state-lock"
  }
}
```

DynamoDB table requires `LockID` as partition key (String). State locking prevents concurrent `apply` operations from corrupting state.

**GCS (GCP):**
```hcl
terraform {
  backend "gcs" {
    bucket = "my-company-tfstate"
    prefix = "prod/gke"
  }
}
```

GCS provides built-in state locking. Enable Object Versioning on the bucket for state recovery.

**Terraform Cloud / HCP Terraform:**
```hcl
terraform {
  cloud {
    organization = "my-org"
    workspaces {
      name = "prod-eks"
    }
  }
}
```

### State Manipulation

```bash
# Import existing resource into state
terraform import aws_s3_bucket.data my-existing-bucket

# Import block (Terraform 1.5+ — preferred)
import {
  to = aws_s3_bucket.data
  id = "my-existing-bucket"
}

# Move resource within state (rename without destroy/recreate)
terraform state mv aws_instance.old aws_instance.new

# Remove resource from state (unmanage without destroying)
terraform state rm aws_instance.legacy

# Generate config for import (Terraform 1.5+)
terraform plan -generate-config-out=generated.tf
```

### moved Block (Terraform 1.1+)

The declarative `moved` block is the preferred way to rename or reorganize resources without destroying them. Unlike `terraform state mv`, it is tracked in version control.

```hcl
moved {
  from = aws_instance.old_name
  to   = aws_instance.new_name
}

moved {
  from = aws_s3_bucket.data
  to   = module.storage.aws_s3_bucket.data
}
```

`moved` blocks are safe to leave in place — they are no-ops once the state migration has been applied.

### check Block (Terraform 1.5+)

`check` blocks add post-condition assertions that run after every `plan` and `apply` without blocking the operation.

```hcl
check "rds_publicly_accessible" {
  assert {
    condition     = !aws_db_instance.main.publicly_accessible
    error_message = "RDS instance must not be publicly accessible."
  }
}
```

### Terragrunt

Terragrunt adds DRY backend configuration, module orchestration via `dependency` blocks, and `run-all` commands across modules.

```hcl
remote_state {
  backend = "s3"
  config = {
    bucket = "my-tfstate"
    key    = "${path_relative_to_include()}/terraform.tfstate"
    region = "us-east-1"
  }
}
```

---

## Terraform Testing

### Testing Pyramid

| Layer | Tool | When to Run | Cost |
|---|---|---|---|
| Syntax | `terraform validate` | Every commit | Free (no API) |
| Linting | `tflint` | Every commit | Free (no API) |
| Security | `checkov` / `tfsec` | Every commit | Free (no API) |
| Unit | `terraform test` (native) | Per PR | Free (mocked) |
| Integration | Terratest | Per PR / nightly | Real cloud cost |
| Compliance | Sentinel / OPA | Plan stage | Terraform Cloud |

### terraform validate and plan

```bash
terraform init -backend=false   # For CI validation without backend
terraform validate
terraform plan -out=tfplan
terraform show -json tfplan      # Parse plan output in CI
```

### TFLint

```bash
tflint --init
tflint --recursive
```

Catches: invalid instance types, deprecated resource arguments, naming convention violations.

### Checkov

```bash
checkov -d ./environments/prod
checkov -d . --output sarif --output-file results.sarif
```

Supports 750+ built-in checks covering CIS benchmarks, PCI-DSS, HIPAA, SOC2 controls.

### Native terraform test (Terraform 1.6+)

```hcl
run "bucket_created_with_correct_name" {
  command = plan

  assert {
    condition     = aws_s3_bucket.data.bucket == "prod-myproject-data"
    error_message = "Bucket name did not match expected pattern"
  }
}
```

### Terraform vs Pulumi vs CDK vs CloudFormation

| Dimension | Terraform/OpenTofu | Pulumi | AWS CDK | CloudFormation |
|---|---|---|---|---|
| Language | HCL (declarative DSL) | TypeScript, Python, Go, Java, .NET | TypeScript, Python, Java, .NET | YAML/JSON |
| Cloud support | All major clouds + 3000+ providers | All major clouds | AWS only | AWS only |
| Learning curve | Low | Medium | Medium | High (verbose YAML) |
| Multi-cloud | Excellent | Excellent | No | No |

**Decision guidance:**
- Default to Terraform — largest ecosystem, talent pool, lowest risk
- Choose Pulumi when HCL's limited logic is a real pain point (complex conditionals, loops)
- Choose CDK when team already uses it and AWS-only is acceptable

### OpenTofu

OpenTofu is the Linux Foundation-maintained open-source fork of Terraform (created 2023 after HashiCorp changed license to BSL 1.1). Drop-in replacement for Terraform 1.x.

```bash
tofu init
tofu plan
tofu apply
```

Choose OpenTofu when: organization requires a fully OSS toolchain, or avoiding HashiCorp/IBM licensing restrictions.

---

## Kafka Architecture

### Core Components

```
Producers → [Topic: orders] → Consumers
               ├── Partition 0 → Replica (Leader: broker-1, Follower: broker-2)
               ├── Partition 1 → Replica (Leader: broker-2, Follower: broker-3)
               └── Partition 2 → Replica (Leader: broker-3, Follower: broker-1)
```

- **Broker**: stores partitions, serves client requests, coordinates replication. Cluster of 3+ brokers provides HA.
- **Topic**: named stream of records divided into partitions.
- **Partition**: ordered, immutable log. Unit of parallelism and replication.
- **Offset**: sequential integer identifying each message within a partition.
- **Consumer Group**: one or more consumers sharing work. Each partition is assigned to exactly one consumer in a group.
- **Replication Factor**: `replication.factor=3` with `min.insync.replicas=2` is a common production setting.

### KRaft Mode (Kafka 4.0+)

ZooKeeper is fully removed in Kafka 4.0. KRaft uses a built-in Raft quorum of controller nodes:

```properties
# Combined broker+controller (small clusters / dev)
process.roles=broker,controller
node.id=1
controller.quorum.voters=1@kafka1:9093,2@kafka2:9093,3@kafka3:9093
```

For clusters with 6+ brokers, use dedicated controller-only nodes (3 or 5) to isolate metadata operations from the data path.

### Topic Configuration

```bash
kafka-topics.sh --create \
  --bootstrap-server kafka:9092 \
  --topic orders \
  --partitions 12 \
  --replication-factor 3 \
  --config retention.ms=604800000 \
  --config min.insync.replicas=2
```

**Partition sizing guidance:**
- Start with 2-4 partitions per broker for moderate workloads
- Each consumer in a group maps to one partition; more partitions = more parallelism
- Partitions can be increased but not decreased — plan conservatively

### Producer Configuration

```properties
bootstrap.servers=kafka1:9092,kafka2:9092,kafka3:9092
acks=all
enable.idempotence=true
retries=2147483647
max.in.flight.requests.per.connection=5
compression.type=lz4
linger.ms=5
batch.size=65536
```

### Consumer Configuration

```properties
bootstrap.servers=kafka1:9092,kafka2:9092
group.id=order-processing-service
auto.offset.reset=earliest
enable.auto.commit=false     # Manual commit for at-least-once guarantee
max.poll.records=500
session.timeout.ms=30000
heartbeat.interval.ms=10000
```

### Exactly-Once Semantics (EOS)

Three layers required:
1. **Idempotent Producer**: `enable.idempotence=true` — deduplicates retries at broker
2. **Transactions**: `transactional.id=my-producer-1` — atomic write across multiple partitions/topics
3. **Idempotent Consumer**: application-level deduplication using message keys or business IDs

```java
producer.initTransactions();
try {
    producer.beginTransaction();
    producer.send(new ProducerRecord<>("topic-a", key, value));
    producer.send(new ProducerRecord<>("topic-b", key, value));
    producer.commitTransaction();
} catch (Exception e) {
    producer.abortTransaction();
}
```

**Cost**: Expect 10-30% throughput reduction from transactional overhead.

### Tiered Storage

Kafka 3.6+ offloads older log segments to S3/GCS/Azure Blob while maintaining local tier for recent data. Enables cheaper long-term retention, decoupled storage/broker scaling, and faster cluster recovery.

---

## Kafka Connect & Streams

### Kafka Connect Architecture

```
External System → Source Connector → Kafka Topic → Sink Connector → External System
                        ↑                                 ↑
                   Transforms (SMTs)              Transforms (SMTs)
                        ↑                                 ↑
                 Schema Registry               Schema Registry
```

### Common Source Connector Patterns

| Pattern | Connector | Use Case |
|---|---|---|
| CDC (Change Data Capture) | Debezium (MySQL, PG, MongoDB, SQL Server) | Real-time DB replication |
| JDBC Polling | Confluent JDBC Source | Batch DB ingestion |
| File Spooling | SpoolDir Source | Log/CSV file ingestion |
| REST Polling | HTTP Source | API data ingestion |

### Common Sink Connector Patterns

| Pattern | Connector | Use Case |
|---|---|---|
| Data Warehouse | BigQuery Sink, Snowflake Sink | Analytics loading |
| Object Storage | S3 Sink, GCS Sink | Data lake (Parquet/Avro) |
| Search | Elasticsearch Sink | Full-text search indexing |
| Cache Invalidation | Redis Sink | Cache warming/invalidation |

### Single Message Transforms (SMTs)

```json
{
  "transforms": "mask,rename",
  "transforms.mask.type": "org.apache.kafka.connect.transforms.MaskField$Value",
  "transforms.mask.fields": "credit_card,ssn",
  "transforms.rename.type": "org.apache.kafka.connect.transforms.ReplaceField$Value",
  "transforms.rename.renames": "customerId:customer_id"
}
```

**Warning**: SMTs run synchronously in the connector's data path. Avoid expensive operations (external API calls, complex computation) — use Kafka Streams or Flink instead.

### Schema Registry

Confluent Schema Registry provides a centralized schema store supporting Avro, Protobuf, and JSON Schema.

**Compatibility modes:**
- `BACKWARD` (default): new schema can read data written by old schema (add optional fields with defaults)
- `FORWARD`: old schema can read data written by new schema
- `FULL`: both backward and forward
- `NONE`: no compatibility checking (dangerous for production)

**Schema evolution rules:**
- Add fields with default values (backward-compatible)
- Never remove required fields or change field types
- Mark fields deprecated before removing them

### Kafka Streams vs Apache Flink

| Dimension | Kafka Streams | Apache Flink |
|---|---|---|
| Deployment | Library embedded in app | Separate cluster |
| Language | Java/Scala (JVM only) | Java, Scala, Python (PyFlink) |
| Operations overhead | Low | High |
| Exactly-once | Yes (within Kafka ecosystem) | Yes (cross-system via two-phase commit) |
| SQL support | ksqlDB (separate service) | Flink SQL / Table API (built-in) |
| Best for | Kafka-native microservices, simple stateful ops | Complex stateful processing, multi-source joins |

---

## Kafka Security

### Authentication Mechanisms

| Mechanism | Description | Production Recommendation |
|---|---|---|
| SASL/PLAIN | Username/password in plaintext | Only with TLS; dev/staging only |
| SASL/SCRAM-SHA-512 | Challenge-response, hashed credentials | Good for most teams |
| SASL/GSSAPI (Kerberos) | Enterprise SSO | Enterprise with existing Kerberos |
| SASL/OAUTHBEARER | OAuth 2.0 / OIDC tokens | Cloud-native, short-lived tokens |
| mTLS | Mutual TLS certificate auth | High-security environments |

**Key distinction**: `SASL/PLAIN` (credential type) vs `SASL_SSL` (security protocol with TLS). Always use `SASL_SSL` in production — `SASL_PLAINTEXT` sends credentials unencrypted.

### ACLs

```bash
# Grant producer permission
kafka-acls.sh --bootstrap-server kafka:9092 \
  --add \
  --allow-principal User:order-service \
  --operation Write \
  --operation Describe \
  --topic orders

# Grant consumer group permission
kafka-acls.sh --bootstrap-server kafka:9092 \
  --add \
  --allow-principal User:analytics-service \
  --operation Read \
  --topic orders \
  --group analytics-consumer-group
```

### Security Best Practices

1. Use TLS encryption for all inter-broker and client-broker communication
2. Prefer SASL/SCRAM-SHA-512 for team deployments
3. Use SASL/OAUTHBEARER with short-lived tokens for cloud-native / service mesh environments
4. Apply principle of least privilege with ACLs — producers get Write only, consumers get Read only
5. Rotate credentials regularly; use secrets managers (Vault, AWS Secrets Manager) for credential injection
6. Enable audit logging (`authorizer.class.name=kafka.security.authorizer.AclAuthorizer`)
7. Separate admin credentials from application credentials

---

## Managed Kafka Services

### Comparison Matrix

| Feature | Confluent Cloud | Amazon MSK | Azure Event Hubs |
|---|---|---|---|
| Full Kafka compatibility | Yes | Yes | Partial |
| Multi-cloud | Yes | No (AWS only) | No (Azure only) |
| Schema Registry | Managed | Self-managed / AWS Glue | No |
| Kafka Connect | Managed (200+ connectors) | MSK Connect (self-managed workers) | Limited |
| ksqlDB | Managed | No | No |
| Auth | API Keys, SASL, OAuth | IAM, SASL/SCRAM, mTLS | SAS tokens, Azure AD |
| Best cloud fit | Any / Multi-cloud | AWS | Azure |

**Decision guidance:**
- Confluent Cloud: multi-cloud, teams wanting fully managed Schema Registry + ksqlDB + connectors
- Amazon MSK: AWS-native teams, IAM auth, lower cost for basic Kafka usage
- Azure Event Hubs: Azure-native teams with existing Event Hubs investment; not for complex Kafka-native use cases

### Amazon MSK Terraform

```hcl
resource "aws_msk_cluster" "main" {
  cluster_name           = "${var.env}-kafka"
  kafka_version          = var.kafka_version
  number_of_broker_nodes = 3

  broker_node_group_info {
    instance_type  = "kafka.m5.large"
    client_subnets = var.private_subnet_ids
    storage_info {
      ebs_storage_info { volume_size = 1000 }
    }
  }

  encryption_info {
    encryption_in_transit {
      client_broker = "TLS"
      in_cluster    = true
    }
  }

  client_authentication {
    sasl { iam = true }
  }
}
```

---

## Common Patterns

### Terraform Environment Promotion

```
environments/
  base/          # Shared modules called by all envs
  dev/           # tfvars: small instances, no HA
  staging/       # tfvars: prod-like but smaller
  prod/          # tfvars: full HA, backups enabled
```

### Kafka Outbox Pattern (transactional messaging)

```
App → INSERT into orders table + outbox table (same DB transaction)
Debezium CDC → Capture outbox table changes → Kafka topic
Consumer → Process event from Kafka
```

Guarantees message is published if and only if the DB transaction commits. Eliminates dual-write inconsistency.

### Kafka Dead Letter Queue (DLQ)

```json
{
  "errors.tolerance": "all",
  "errors.deadletterqueue.topic.name": "orders.dlq",
  "errors.deadletterqueue.context.headers.enable": true,
  "errors.log.enable": true
}
```

---

## Anti-Patterns

### Terraform Anti-Patterns

| Anti-Pattern | Problem | Fix |
|---|---|---|
| Storing secrets in tfvars | Credentials in version control | Use AWS Secrets Manager / Vault / env vars |
| Monolithic root modules | One giant main.tf — hard to plan, slow to apply | Split into focused modules |
| Manual state edits | State corruption | Use `terraform state` commands |
| No state locking | Concurrent apply corrupts state | DynamoDB/GCS lock |
| Unpinned providers | Breaking changes pulled silently | Use `~>` version constraints |
| `terraform apply` without plan review | Unexpected changes in production | Always review plan first |
| Over-using `count` for logical branching | Resource destruction on list reordering | Prefer `for_each` with maps |

### Kafka Anti-Patterns

| Anti-Pattern | Problem | Fix |
|---|---|---|
| Too many small topics | Thousands of 1-partition topics kill broker performance | Consolidate with routing keys |
| Too few partitions | Cannot scale consumers beyond partition count | Plan partitions upfront |
| Using Kafka as a database | No random-access queries | Sink to a proper store |
| Ignoring consumer lag | Consumers falling behind without alert | Alert on `records-lag-max` |
| Auto-commit before processing | Message lost if processing fails | `enable.auto.commit=false` |
| SASL_PLAINTEXT in production | Credentials in plaintext | Always use SASL_SSL |
| No Schema Registry | Deserialization errors on producer changes | Use Schema Registry |
| Expensive SMT logic | Blocks connector throughput | Move complex transforms to Kafka Streams |
| Not setting `min.insync.replicas` | `acks=all` can succeed with only 1 replica | Always set to 2 for 3-replica clusters |

---

## Troubleshooting

### Terraform Troubleshooting

```bash
TF_LOG=DEBUG terraform plan 2>&1 | tee tf-debug.log
terraform apply -refresh-only   # Sync state to match real infra
terraform plan -detailed-exitcode
# Exit 0: no changes; Exit 1: error; Exit 2: changes needed
terraform force-unlock <lock-id>  # Lock stuck from killed apply
```

**Common errors:**
- `Error acquiring the state lock`: Another apply running, or a previous one died. Check DynamoDB for stale lock.
- `Resource already exists`: Import the resource with `terraform import` before applying.
- `Cycle error`: Circular dependency. Use `depends_on` carefully or restructure.

### Kafka Troubleshooting

```bash
# Check consumer group lag
kafka-consumer-groups.sh --bootstrap-server kafka:9092 \
  --describe --group my-consumer-group

# Reset offsets (reprocess from beginning)
kafka-consumer-groups.sh --bootstrap-server kafka:9092 \
  --group my-group --topic orders --reset-offsets --to-earliest --execute

# Check topic end offsets
kafka-run-class.sh kafka.tools.GetOffsetShell \
  --bootstrap-server kafka:9092 --topic orders --time -1
```

**Key metrics to monitor:**

| Metric | Warning Threshold | Description |
|---|---|---|
| `records-lag-max` | > 1000 | Consumer group falling behind |
| `BytesInPerSec` | Near capacity | Ingress saturation |
| `UnderReplicatedPartitions` | > 0 | Broker failure or network partition |
| `ActiveControllerCount` | != 1 | Controller election issue |

**Common issues:**
- `LEADER_NOT_AVAILABLE`: Normal briefly after topic creation; persistent = broker down
- `NOT_ENOUGH_REPLICAS`: `min.insync.replicas` not met — check broker health
- Consumer rebalance storms: Increase `session.timeout.ms`, check `max.poll.interval.ms` vs processing time

---

## References

**Terraform:**
- [HashiCorp Terraform Documentation](https://developer.hashicorp.com/terraform/docs)
- [Standard Module Structure](https://developer.hashicorp.com/terraform/language/modules/develop/structure)
- [Terraform Best Practices - AWS Prescriptive Guidance](https://docs.aws.amazon.com/prescriptive-guidance/latest/terraform-aws-provider-best-practices/structure.html)
- [How to Test Terraform Code (Spacelift)](https://spacelift.io/blog/terraform-test)

**Kafka:**
- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [Confluent Developer Documentation](https://developer.confluent.io/)
- [Exactly-Once Semantics (Confluent)](https://www.confluent.io/blog/exactly-once-semantics-are-possible-heres-how-apache-kafka-does-it/)
- [Kafka Streams vs Flink Comparison (Conduktor)](https://www.conduktor.io/glossary/kafka-streams-vs-apache-flink)
- [Confluent Cloud vs Amazon MSK](https://www.confluent.io/confluent-cloud-vs-amazon-msk/)
