<!-- hub-reference-banner -->
> **Reference file — part of the `aws-cloud` hub.** Formerly the standalone `aws-ai-ml` skill.
> Sibling topics in this family are now reference files under the hubs (`aws-cloud`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: aws-ai-ml
title: AWS AI/ML Expert (Bedrock + SageMaker)
description: |
  Expert guidance on AWS AI/ML services: Amazon Bedrock (foundation models, agents, RAG, guardrails, AgentCore) and Amazon SageMaker (training, endpoints, MLOps pipelines).
  TRIGGER: designing or reviewing Bedrock or SageMaker applications; Converse API or InvokeModel usage; Bedrock Agents or Knowledge Bases (RAG); Bedrock Guardrails, fine-tuning, or model evaluation; choosing between on-demand/provisioned/batch inference; prompt caching, cross-region inference, intelligent routing; SageMaker training jobs, endpoints, Pipelines, Feature Store, Ground Truth, Canvas, HyperPod; integrating Bedrock with Lambda/EventBridge/Step Functions; troubleshooting IAM AccessDenied, throttling, cold starts, or model-not-enabled errors; any question containing "Bedrock", "SageMaker", "foundation model on AWS", "LLM inference on AWS", "model fine-tuning on AWS".
  SKIP: generic Python ML/data science (scikit-learn, pandas, PyTorch local training) with no AWS deployment; OpenAI SDK, Azure OpenAI, Google Vertex AI, or Hugging Face Inference Endpoints when no Bedrock/SageMaker integration is present; vector database questions for Pinecone, Weaviate, ChromaDB, or MongoDB Atlas Vector Search when no Bedrock/SageMaker integration is mentioned.
category: developer
version: "1.1.1"
updated: "2026-05-31"
keywords:
  - AWS
  - Amazon Bedrock
  - SageMaker
  - foundation models
  - RAG
  - Bedrock Agents
  - AgentCore
  - Knowledge Bases
  - Guardrails
  - Converse API
  - InvokeModel
  - provisioned throughput
  - prompt caching
  - cross-region inference
  - SageMaker Pipelines
  - MLOps
  - HyperPod
  - Claude on AWS
  - LLM inference
whenToUse:
  - Designing or reviewing applications that use Amazon Bedrock (Claude, Titan, Llama, Mistral, Cohere, Nova)
  - Using the Bedrock Converse API, InvokeModel API, or streaming inference
  - Building Bedrock Agents, AgentCore, or Knowledge Bases (RAG pipelines)
  - Configuring Bedrock Guardrails, model evaluation, or fine-tuning
  - Choosing between on-demand, provisioned throughput, and batch inference modes
  - Setting up prompt caching, cross-region inference, or intelligent prompt routing
  - Integrating Bedrock with Lambda, EventBridge, Step Functions, or API Gateway
  - Training or deploying ML models on SageMaker (training jobs, endpoints, HyperPod)
  - MLOps with SageMaker Pipelines, Model Registry, Feature Store, Ground Truth, Canvas
  - Troubleshooting IAM AccessDenied, throttling, cold starts, or model-not-enabled errors on Bedrock/SageMaker
whenNotToUse:
  - Generic Python ML/data science (scikit-learn, pandas, PyTorch local training) with no AWS deployment
  - OpenAI SDK, Azure OpenAI, Google Vertex AI, or Hugging Face Inference Endpoints (non-AWS)
  - Vector database questions (Pinecone, Weaviate, ChromaDB, MongoDB Atlas Vector Search) when no Bedrock/SageMaker integration is mentioned
  - General embeddings questions not tied to a Bedrock or SageMaker deployment
related_skills:
  - aws-core
  - aws-serverless
  - rag-architecture
  - mongodb-atlas-expert
---

# AWS AI/ML Expert (Bedrock + SageMaker)

Expert guidance on AWS AI/ML services: Amazon Bedrock (foundation models, agents, RAG, guardrails, AgentCore) and Amazon SageMaker (training, endpoints, MLOps pipelines).

## Reference context

A detailed context file lives at `references/aws-ai-ml-context.md` relative to this skill's directory. Read it before answering any Bedrock or SageMaker question. If unavailable, proceed from the inline knowledge below and note the gap.

The context file covers:
- Decision tables: Bedrock vs SageMaker, model selection, inference mode selection
- Bedrock Converse API patterns and SDK code examples
- Bedrock Agents, AgentCore, and Knowledge Bases architecture (RAG workflow)
- Guardrails configuration and pricing
- Provisioned throughput vs on-demand vs batch cost analysis
- Prompt caching, cross-region inference, intelligent routing, model distillation
- SageMaker endpoint types (real-time, serverless, async) with tradeoff table
- SageMaker Pipelines, Model Registry, Feature Store, Ground Truth, Canvas, HyperPod
- Integration patterns: Bedrock + Lambda, Bedrock + EventBridge, SageMaker + API Gateway
- Security: VPC endpoints, IAM least-privilege, KMS encryption, model access logging
- Common failure modes: throttling, model not enabled in region, IAM AccessDenied, endpoint cold starts

## How to use this skill

1. Read `references/aws-ai-ml-context.md` before answering any Bedrock or SageMaker question. If the file is missing, use the inline decision shortcuts below and ask 1–2 targeted questions to narrow scope before answering depth questions (e.g., "What runtime? What traffic volume?").
2. Apply the Bedrock-vs-SageMaker decision table when the user needs to pick a service. If no table row matches, ask 2–3 targeted context questions (model type, traffic pattern, data sovereignty requirements) before recommending.
3. For design questions ("how do I build X?"), produce numbered steps + key decision points rather than prose. For model selection, use the model comparison table — Claude is the default recommendation for enterprise text tasks.
4. For cost questions, check the inference mode table; note batch has a 24-hour max turnaround and provisioned throughput requires a 1-month minimum commitment.
5. Recommend prompt caching whenever the application has long, repeated system prompts (up to 85% latency/cost reduction).
6. For RAG, lead with Bedrock Knowledge Bases (fully managed) before recommending custom vector DB solutions.
7. For security, always mention VPC endpoints + IAM least-privilege + KMS at-rest encryption as a baseline.
8. For troubleshooting, apply the failure mode table:
   - Model not enabled in region → request access in Bedrock console (Models → Model access)
   - IAM AccessDenied → add `bedrock:InvokeModel` (and `bedrock:InvokeModelWithResponseStream` for streaming) to the role
   - Throttling → switch to provisioned throughput or reduce request concurrency
   - SageMaker endpoint cold start → use Provisioned Concurrency or scheduled keep-warm invocations

## Key decision shortcuts

### Bedrock vs SageMaker

| Need | Service |
|---|---|
| Hosted foundation model (Claude, Llama, etc.) with zero ML infra | Bedrock |
| Train a custom model, fine-tune with own data, or run proprietary weights | SageMaker |
| Production RAG pipeline with managed ingestion + retrieval | Bedrock Knowledge Bases |
| Complex ML pipelines with feature engineering, data labeling, model monitoring | SageMaker |
| No-code AutoML on tabular data | SageMaker Canvas |
| Large-scale distributed LLM pre-training or fine-tuning (100B+ params) | SageMaker HyperPod |

### Inference mode

| Use case | Mode | Notes |
|---|---|---|
| Interactive chat, real-time API responses | On-demand (Converse API) | Standard pricing |
| Document summarization, offline enrichment, batch eval | Batch inference | 50% cost savings; 24-hour max turnaround |
| High-volume, steady predictable load (>10M tokens/day) | Provisioned throughput | 15–40% savings; 1-month minimum |
| Long or repeated system prompts | Prompt caching | Up to 85% latency/cost reduction |

### Model selection (Bedrock)

| Task | Recommended model |
|---|---|
| Enterprise text, reasoning, code, safety-critical | Claude (Anthropic) |
| Cost-sensitive, simple classification, low-latency | Amazon Nova Micro/Lite |
| Multimodal (image + text input, video understanding) | Amazon Nova Pro or Claude 3.x (vision) |
| Embeddings, semantic search, RAG retrieval | Cohere Embed or Amazon Titan Embeddings |
| Open-weights fine-tuning | Meta Llama 4 |
| EU data sovereignty required | Mistral Large |

### SageMaker endpoint type

| Requirement | Endpoint type |
|---|---|
| Sub-100ms latency, persistent traffic | Real-time endpoint |
| Bursty/unpredictable traffic, cold starts acceptable | Serverless endpoint |
| Large payloads (up to 1 GB), long processing (up to 1 hour) | Async endpoint |
| Cost-sensitive batch scoring | Batch Transform job |

> Self-hosting an open-weights LLM on your own GPUs (EC2/EKS) rather than a managed endpoint — choosing/tuning a serving engine (vLLM, SGLang, TensorRT-LLM, TGI), PagedAttention/KV-cache sizing, continuous batching, chunked prefill, speculative decoding, TTFT/TPOT/goodput SLOs, or GPU autoscaling/cost-per-token — is a serving-runtime concern: see the `llm-inference-serving` skill (`ai-agent-engineering` hub). This skill covers the *managed* AWS path (Bedrock + SageMaker endpoints).
