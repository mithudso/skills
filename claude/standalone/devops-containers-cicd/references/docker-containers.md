<!-- hub-reference-banner -->
> **Reference file — part of the `devops-containers-cicd` hub.** Formerly the standalone `docker-containers` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: docker-containers
title: "Docker and Containers Expert"
description: >
  Docker and container engineering expert covering Dockerfile authoring, image optimization,
  container security, Docker Compose v2, container networking, registries, BuildKit, Podman,
  OCI standards, and container anti-patterns.
  TRIGGER: user asks about Dockerfile, multi-stage build, container image size, distroless or
  Chainguard images, BuildKit cache mounts or secrets, Docker Compose health checks or secrets,
  container networking bridge/overlay/macvlan, image scanning with Trivy or Grype, image signing
  with Cosign, rootless containers, Podman vs Docker, Docker Swarm vs Kubernetes, or container
  security hardening (non-root user, read-only filesystem, capability dropping, seccomp).
  SKIP: Kubernetes cluster management or Helm (use kubernetes-networking); CI/CD pipeline
  construction beyond container build/scan/push steps (use cicd-pipelines); AWS ECS/Fargate
  or cloud-specific container orchestration (use aws-serverless or aws-core).
category: developer
version: "1.1.0"
updated: "2026-05-29"
triggers:
  - "Docker"
  - "Dockerfile"
  - "container"
  - "docker compose"
  - "docker-compose"
  - "multi-stage build"
  - "container image"
  - "container security"
  - "BuildKit"
  - "distroless"
  - "Podman"
  - "OCI"
  - "container registry"
  - "Docker Hub"
  - "GHCR"
  - "Harbor"
  - "container networking"
  - "docker volume"
  - "container orchestration"
  - "docker build"
  - "image optimization"
  - "container scanning"
  - "Trivy"
  - "Chainguard"
  - "Alpine image"
  - "scratch image"
  - "Docker Swarm"
  - "rootless container"
  - "docker secret"
  - "Quadlet"
whenToUse:
  - Writing or reviewing a Dockerfile
  - Reducing container image size with multi-stage builds or distroless base images
  - Securing container workloads (non-root user, capability dropping, read-only filesystem)
  - Designing a Docker Compose stack with health checks, secrets, and network isolation
  - Troubleshooting container networking (bridge, overlay, macvlan)
  - Choosing between Docker and Podman for a project
  - Integrating container build, scan, and push into CI/CD pipelines
  - Configuring BuildKit cache mounts, build secrets, or multi-platform builds
  - Scanning container images for vulnerabilities with Trivy, Grype, or Docker Scout
  - Implementing rootless containers or Podman Quadlet systemd integration
whenNotToUse:
  - Kubernetes cluster management, Helm charts, or CNI networking (use kubernetes-networking)
  - CI/CD pipeline structure beyond container build/scan/push (use cicd-pipelines)
  - AWS ECS, Fargate, or App Runner container orchestration (use aws-serverless)
related_skills:
  - cicd-pipelines
  - kubernetes-networking
  - code-packaging
  - aws-serverless
  - linux-sysadmin
---

# Docker and Containers Expert Context

Expert reference for Docker, OCI containers, Podman, and the container ecosystem. Covers Dockerfile authoring, image optimization, security hardening, Compose orchestration, networking, registries, BuildKit, and production deployment patterns. A response from this skill is correct when it applies the right pattern for the stated constraint (security, size, platform, tooling), flags dangerous defaults, and stays within the stated scope boundaries.

> **Staleness note:** Version-specific claims (Docker Engine 23.0+, BuildKit defaults, Compose Spec v5.0.0, DHI open-source Dec 2025, Falco v0.43.0+ Jan 2026) were current as of May 2026. Verify against current Docker release notes before relying on version-gated behavior.

**Navigation by task:**
- Dockerfile authoring and layer optimization → §1 Dockerfile Best Practices
- Multi-stage builds (Go/scratch, Python/distroless, Java/jlink, Rust/cargo-chef) → §2
- BuildKit cache mounts, secrets, SSH forwarding, heredoc, multi-platform → §3
- Security hardening (scanning, runtime, seccomp, signing, secrets) → §4
- Docker Compose v2 (profiles, health checks, override files, watch, Bake) → §5
- Container networking (bridge, host, overlay, macvlan) → §6
- Volumes and storage → §7
- Registry comparison and lifecycle management → §8
- Podman compatibility, rootless, Quadlet → §9
- OCI standards and artifacts → §10
- Docker Swarm vs Kubernetes decision → §11
- Anti-patterns → §12
- Production checklists (Dockerfile, security, Compose) → §13
- Quick reference commands → `references/docker-quick-reference.md`
- Reference links → `references/docker-references.md`

---

## 1. Dockerfile Best Practices

### 1.1 Base Image Selection

Choose the smallest base image that satisfies your runtime requirements.

| Base Image | Size | Shell | Pkg Manager | Use Case |
|---|---|---|---|---|
| `scratch` | 0 B | No | No | Static binaries (Go, Rust) |
| `distroless/static` | ~2 MB | No | No | Static binaries with CA certs |
| `distroless/base` | ~20 MB | No | No | Dynamic binaries (C/C++) |
| `distroless/java` | ~80 MB | No | No | JVM applications |
| `alpine:3.21` | ~5 MB | Yes | apk | General-purpose minimal |
| `node:22-alpine` | ~50 MB | Yes | apk+npm | Node.js applications |
| `chainguard/node` | ~45 MB | No | No | Hardened Node.js |
| `ubuntu:24.04` | ~78 MB | Yes | apt | Full Linux userland |

**Docker Hardened Images (DHI):** In December 2025, Docker open-sourced 1,000+ hardened container images under Apache 2.0. Each ships with cryptographically signed metadata, SBOMs, and continuous CVE patching. Prefer DHI images over vanilla official images when available.

**Chainguard Images:** Minimal, FIPS-compliant base images rebuilt daily with automatic SBOM generation. Ideal for regulated environments.

### 1.2 Image Pinning Strategy

```dockerfile
# Tag-based (flexible but mutable)
FROM node:22-alpine

# Digest-pinned (immutable, reproducible)
FROM node:22-alpine@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c

# Version-ranged (compromise)
FROM python:3.12-slim
```

Always pin digests in production Dockerfiles. Use Docker Scout or Renovate to automate digest update PRs when upstream publishes new versions.

### 1.3 Layer Ordering and Cache Optimization

Docker evaluates instructions sequentially and caches results. Order instructions from least-changing to most-changing.

```dockerfile
# GOOD: Dependencies before source code
FROM node:22-alpine AS build
WORKDIR /app

# Layer 1: rarely changes
COPY package.json package-lock.json ./

# Layer 2: cached unless lock file changes
RUN npm ci --only=production

# Layer 3: changes frequently
COPY src/ ./src/
RUN npm run build
```

```dockerfile
# BAD: Invalidates dependency cache on every code change
COPY . .
RUN npm ci && npm run build
```

### 1.4 RUN Instruction Patterns

Combine related commands to reduce layers. Clean up in the same layer.

```dockerfile
# Package installation (Debian/Ubuntu)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    gnupg \
    && rm -rf /var/lib/apt/lists/*

# Package installation (Alpine)
RUN apk add --no-cache \
    curl \
    jq \
    openssl

# Pipe safety
RUN set -o pipefail && curl -fsSL https://example.com/install.sh | sh
```

### 1.5 COPY vs ADD

Use `COPY` for local file transfers. Reserve `ADD` for remote URL downloads and automatic tar extraction.

```dockerfile
# Preferred: explicit and predictable
COPY requirements.txt /app/

# Only when you need auto-extraction or remote fetch
ADD https://example.com/archive.tar.gz /tmp/
```

**Bind mount alternative** for files needed only during build:

```dockerfile
RUN --mount=type=bind,source=requirements.txt,target=/tmp/requirements.txt \
    pip install --requirement /tmp/requirements.txt
```

### 1.6 CMD vs ENTRYPOINT

```dockerfile
# ENTRYPOINT: the executable (rarely overridden)
# CMD: default arguments (easily overridden)
ENTRYPOINT ["node"]
CMD ["server.js"]

# Combined with helper script for init logic
COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["postgres"]
```

In entrypoint scripts, use `exec "$@"` so the application becomes PID 1 and receives Unix signals (SIGTERM, SIGINT) properly for graceful shutdown.

### 1.7 ARG vs ENV

```dockerfile
# ARG: build-time only, not persisted in image
ARG NODE_ENV=production
ARG APP_VERSION

# ENV: persisted in image, available at runtime
ENV NODE_ENV=${NODE_ENV}

# Prevent sensitive ARG leakage
RUN --mount=type=secret,id=api_key \
    export API_KEY=$(cat /run/secrets/api_key) && \
    ./configure --key=$API_KEY
```

### 1.8 .dockerignore

Always include a `.dockerignore` to reduce build context size and prevent secret leakage.

```gitignore
# Version control
.git
.gitignore

# Dependencies (rebuilt in container)
node_modules/
vendor/
__pycache__/
*.pyc

# Build artifacts
dist/
build/
*.egg-info/

# IDE and OS
.vscode/
.idea/
.DS_Store
Thumbs.db

# Secrets and config
.env
.env.*
*.pem
*.key
credentials.json

# Documentation
*.md
LICENSE
docs/

# Tests (unless needed for build)
tests/
__tests__/
*.test.*
*.spec.*

# Docker files (prevent recursive builds)
Dockerfile*
docker-compose*
.dockerignore
```

### 1.9 Health Checks

```dockerfile
# HTTP health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:3000/health || exit 1

# TCP health check
HEALTHCHECK --interval=15s --timeout=3s --retries=5 \
    CMD nc -z localhost 5432 || exit 1

# File-based health check
HEALTHCHECK --interval=30s CMD test -f /tmp/healthy || exit 1
```

### 1.10 Non-Root User

```dockerfile
# Debian/Ubuntu
RUN groupadd --system --gid 1001 appgroup && \
    useradd --system --uid 1001 --gid appgroup --no-log-init appuser

# Alpine
RUN addgroup -g 1001 -S appgroup && \
    adduser -S -u 1001 -G appgroup appuser

# Set ownership and switch
COPY --chown=appuser:appgroup . /app
USER appuser
```

Always assign explicit UID/GID values for reproducibility across hosts.

---

## 2. Multi-Stage Builds

### 2.1 Core Pattern

Separate build dependencies from runtime to produce minimal final images.

```dockerfile
# Stage 1: Build
FROM node:22-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: Production
FROM node:22-alpine AS production
RUN addgroup -g 1001 -S app && adduser -S -u 1001 -G app app
WORKDIR /app
COPY --from=builder --chown=app:app /app/dist ./dist
COPY --from=builder --chown=app:app /app/node_modules ./node_modules
COPY --from=builder --chown=app:app /app/package.json ./
USER app
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=5s CMD wget -qO- http://localhost:3000/health || exit 1
CMD ["node", "dist/server.js"]
```

### 2.2 Go Binary with Scratch

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app ./cmd/server

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app
USER 65534:65534
ENTRYPOINT ["/app"]
```

**Impact:** Typical Go image drops from 800+ MB to 10-15 MB.

### 2.3 Python with Distroless

```dockerfile
FROM python:3.12-slim AS builder
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir --prefix=/install -r requirements.txt
COPY . .

FROM gcr.io/distroless/python3-debian12
WORKDIR /app
COPY --from=builder /install /usr/local
COPY --from=builder /app .
USER nonroot:nonroot
CMD ["main.py"]
```

### 2.4 Java with JLink

```dockerfile
FROM eclipse-temurin:21-jdk-alpine AS builder
WORKDIR /app
COPY . .
RUN ./gradlew bootJar
RUN jlink --add-modules $(jdeps --print-module-deps build/libs/app.jar) \
    --strip-debug --compress=2 --no-header-files --no-man-pages \
    --output /custom-jre

FROM alpine:3.21
COPY --from=builder /custom-jre /opt/java
COPY --from=builder /app/build/libs/app.jar /app/app.jar
ENV PATH="/opt/java/bin:${PATH}"
USER 1001:1001
ENTRYPOINT ["java", "-jar", "/app/app.jar"]
```

### 2.5 Rust with Cargo Chef

```dockerfile
FROM rust:1.82-alpine AS chef
RUN apk add --no-cache musl-dev && cargo install cargo-chef
WORKDIR /app

FROM chef AS planner
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

FROM chef AS builder
COPY --from=planner /app/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json
COPY . .
RUN cargo build --release

FROM scratch
COPY --from=builder /app/target/release/myapp /myapp
USER 65534:65534
ENTRYPOINT ["/myapp"]
```

### 2.6 Frontend Static Assets with Nginx

```dockerfile
FROM node:22-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:1.27-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
HEALTHCHECK --interval=30s CMD wget -qO- http://localhost/ || exit 1
```

### 2.7 Parallel Build Stages

BuildKit executes independent stages concurrently:

```dockerfile
FROM golang:1.23 AS build-api
WORKDIR /api
COPY api/ .
RUN go build -o /bin/api

FROM node:22-alpine AS build-frontend
WORKDIR /web
COPY web/ .
RUN npm ci && npm run build

FROM alpine:3.21
COPY --from=build-api /bin/api /usr/local/bin/
COPY --from=build-frontend /web/dist /var/www/html
```

Both `build-api` and `build-frontend` stages run simultaneously.

---

## 3. BuildKit Features

### 3.1 Enabling BuildKit

```bash
# Environment variable
export DOCKER_BUILDKIT=1
docker build .

# Docker daemon config (/etc/docker/daemon.json)
{ "features": { "buildkit": true } }

# BuildKit is the default builder since Docker Engine 23.0+
# Verify:
docker buildx version
```

### 3.2 Cache Mounts

Persist package manager caches across builds without including them in image layers.

```dockerfile
# Node.js
RUN --mount=type=cache,target=/root/.npm \
    npm ci --only=production

# Python
RUN --mount=type=cache,target=/root/.cache/pip \
    pip install -r requirements.txt

# Go
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app ./cmd/server

# Rust
RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=/app/target \
    cargo build --release

# apt (Debian/Ubuntu)
RUN --mount=type=cache,target=/var/cache/apt \
    --mount=type=cache,target=/var/lib/apt/lists \
    apt-get update && apt-get install -y curl
```

Cache mounts can reduce build times by 10x or more for dependency installation.

### 3.3 Build Secrets

Secrets are mounted in memory for a single RUN instruction and never written to image layers or build cache.

```dockerfile
# syntax=docker/dockerfile:1

# Mount a secret file
RUN --mount=type=secret,id=npmrc,target=/root/.npmrc \
    npm ci --only=production

# Read secret into env var
RUN --mount=type=secret,id=api_key \
    export API_KEY=$(cat /run/secrets/api_key) && \
    ./configure --api-key=$API_KEY
```

```bash
# Build with secrets
docker build --secret id=npmrc,src=$HOME/.npmrc \
             --secret id=api_key,src=./api-key.txt \
             -t myapp:v1 .
```

### 3.4 SSH Agent Forwarding

```dockerfile
RUN --mount=type=ssh \
    git clone git@github.com:org/private-repo.git

RUN --mount=type=ssh \
    pip install git+ssh://git@github.com/org/private-lib.git
```

```bash
docker build --ssh default=$SSH_AUTH_SOCK -t myapp .
```

### 3.5 Heredoc Syntax

Write multi-line scripts without shell escaping gymnastics.

```dockerfile
# syntax=docker/dockerfile:1

# Multi-line script
RUN <<EOF
set -ex
apt-get update
apt-get install -y --no-install-recommends curl jq
rm -rf /var/lib/apt/lists/*
EOF

# Inline file creation
COPY <<EOF /etc/nginx/conf.d/default.conf
server {
    listen 80;
    server_name _;
    location / {
        proxy_pass http://app:3000;
    }
}
EOF

# Multi-file copy
COPY <<-app.py <<-config.yaml /app/
#!/usr/bin/env python3
import os
print("Hello")
app.py
database:
  host: localhost
  port: 5432
config.yaml
```

### 3.6 Multi-Platform Builds

```bash
# Create a builder with multi-platform support
docker buildx create --name multiarch --use

# Build for multiple architectures
docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --tag myapp:v1.0 \
    --push .

# Inspect manifest
docker buildx imagetools inspect myapp:v1.0
```

```dockerfile
# Platform-aware base image selection
FROM --platform=$BUILDPLATFORM golang:1.23 AS builder
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /app

FROM --platform=$TARGETPLATFORM alpine:3.21
COPY --from=builder /app /usr/local/bin/app
```

### 3.7 Build Arguments and Metadata

```dockerfile
# syntax=docker/dockerfile:1

# OCI image annotations
FROM alpine:3.21
LABEL org.opencontainers.image.source="https://github.com/org/repo"
LABEL org.opencontainers.image.description="My application"
LABEL org.opencontainers.image.licenses="MIT"

ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
LABEL org.opencontainers.image.created=$BUILD_DATE
LABEL org.opencontainers.image.revision=$VCS_REF
LABEL org.opencontainers.image.version=$VERSION
```

---

## 4. Container Security

### 4.1 Vulnerability Scanning

```bash
# Trivy -- most popular open-source scanner
trivy image --severity CRITICAL,HIGH myapp:v1.2.3

# Fail CI on critical vulnerabilities
trivy image --exit-code 1 --severity CRITICAL myapp:v1.2.3

# Generate SBOM
trivy image --format spdx-json --output sbom.json myapp:v1.2.3

# Scan filesystem (before build)
trivy fs --severity HIGH,CRITICAL .

# Grype -- alternative scanner by Anchore
grype myapp:v1.2.3 --fail-on high

# Docker Scout -- integrated into Docker CLI
docker scout quickview myapp:v1.2.3
docker scout cves myapp:v1.2.3 --only-severity critical,high
docker scout recommendations myapp:v1.2.3
```

**CI/CD integration rule:** Fail the build on HIGH and CRITICAL CVEs. A "zero CVE" scan does not mean the image is secure -- scanners miss application-level misconfigurations, hardcoded secrets, and malicious custom binaries.

### 4.2 Runtime Hardening

> **Kernel primitives behind these flags:** `--cap-drop`/`--cap-add`, `--security-opt seccomp`,
> rootless (`userns-remap`), and `--memory`/`--cpus` are all thin wrappers over Linux **capabilities,
> seccomp, user namespaces, and cgroup v2 controllers**. For how the isolation actually works — the
> eight namespace types, userns UID/GID mapping (subuid/subgid, newuidmap), cgroup v2 controllers/PSI,
> and the container-escape surface (CVE-2022-0492 `release_agent`, why `--privileged` removes isolation) —
> load `references/linux-cgroups-namespaces.md` in this hub.

```yaml
# docker-compose.yml security configuration
services:
  app:
    image: myapp:v1.2.3
    read_only: true
    tmpfs:
      - /tmp:rw,noexec,nosuid,size=64m
    security_opt:
      - no-new-privileges:true
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE    # Only if binding ports < 1024
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 128M
    ulimits:
      nproc: 200
      nofile:
        soft: 65536
        hard: 65536
```

### 4.3 Seccomp Profiles

Restrict system calls to only what the application needs.

```json
{
  "defaultAction": "SCMP_ACT_ERRNO",
  "architectures": ["SCMP_ARCH_X86_64"],
  "syscalls": [
    {
      "names": [
        "read", "write", "open", "close", "stat", "fstat",
        "mmap", "mprotect", "munmap", "brk", "socket",
        "connect", "accept", "sendto", "recvfrom",
        "epoll_create1", "epoll_ctl", "epoll_wait",
        "futex", "nanosleep", "clock_gettime", "exit_group"
      ],
      "action": "SCMP_ACT_ALLOW"
    }
  ]
}
```

```bash
docker run --security-opt seccomp=custom-profile.json myapp:v1
```

### 4.4 Image Signing and Verification

```bash
# Cosign -- sign and verify OCI images
cosign sign --key cosign.key myregistry.io/myapp:v1.2.3

# Verify before deployment
cosign verify --key cosign.pub myregistry.io/myapp:v1.2.3

# Keyless signing with Sigstore (OIDC identity)
cosign sign myregistry.io/myapp:v1.2.3
cosign verify --certificate-identity=user@example.com \
              --certificate-oidc-issuer=https://accounts.google.com \
              myregistry.io/myapp:v1.2.3

# Docker Content Trust (Notary v2)
export DOCKER_CONTENT_TRUST=1
docker push myregistry.io/myapp:v1.2.3
```

### 4.5 Secrets Management

**Build-time secrets (BuildKit):**

```dockerfile
RUN --mount=type=secret,id=npm_token,target=/root/.npmrc \
    npm ci --only=production

RUN --mount=type=secret,id=db_password \
    export DB_PASS=$(cat /run/secrets/db_password) && \
    ./migrate --password=$DB_PASS
```

**Runtime secrets (Docker Swarm / Compose):**

```yaml
services:
  app:
    secrets:
      - db_password
      - api_key
    environment:
      DB_PASSWORD_FILE: /run/secrets/db_password

secrets:
  db_password:
    file: ./secrets/db_password.txt
  api_key:
    external: true
```

**Application code to read file-based secrets:**

```javascript
const { readFileSync } = require('fs');

function getSecret(name) {
  const filePath = process.env[`${name}_FILE`];
  if (filePath) {
    return readFileSync(filePath, 'utf-8').trim();
  }
  return process.env[name] || '';
}

const dbPassword = getSecret('DB_PASSWORD');
```

### 4.6 Runtime Monitoring

**Falco** (v0.43.0+, January 2026) detects anomalous behavior at runtime:

```yaml
# Custom Falco rule
- rule: Shell Spawned in Container
  desc: Alert on shell processes in containers
  condition: >
    spawned_process and container and
    proc.name in (bash, sh, zsh, dash, ash)
  output: >
    Shell spawned in container
    (container=%container.name image=%container.image.repository
     user=%user.name proc=%proc.name cmdline=%proc.cmdline)
  priority: WARNING
```

### 4.7 Docker Socket Protection

The Docker daemon socket (`/var/run/docker.sock`) grants root-equivalent access. Never mount it into application containers.

```yaml
# BAD: gives container full host control
volumes:
  - /var/run/docker.sock:/var/run/docker.sock

# If CI/CD runners need Docker access, use:
# - Docker-in-Docker (dind) with TLS
# - Kaniko for rootless image builds
# - Podman (no daemon socket to mount)
```

---

## 5. Docker Compose v2

### 5.1 File Structure

Compose v2 prefers `compose.yaml` as the canonical filename (over `docker-compose.yml`). Both are still supported, but `compose.yaml` takes precedence when both exist.

```yaml
# compose.yaml
name: my-application

services:
  api:
    build:
      context: ./api
      dockerfile: Dockerfile
      args:
        NODE_ENV: production
    ports:
      - "3000:3000"
    environment:
      DATABASE_URL: postgres://db:5432/mydb
    depends_on:
      db:
        condition: service_healthy
    networks:
      - frontend
      - backend
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:3000/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 15s

  db:
    image: postgres:16-alpine
    volumes:
      - pgdata:/var/lib/postgresql/data
    environment:
      POSTGRES_DB: mydb
      POSTGRES_USER: app
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
    secrets:
      - db_password
    networks:
      - backend
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app -d mydb"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
    volumes:
      - redis-data:/data
    networks:
      - backend
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s

volumes:
  pgdata:
    driver: local
  redis-data:
    driver: local

networks:
  frontend:
    driver: bridge
  backend:
    driver: bridge
    internal: true   # No external access

secrets:
  db_password:
    file: ./secrets/db_password.txt
```

### 5.2 Service Profiles

Group services for different environments without separate files.

```yaml
services:
  api:
    build: ./api
    profiles: []   # Always started (no profile = always active)

  db:
    image: postgres:16-alpine
    profiles: []

  worker:
    build: ./worker
    profiles: ["full", "worker"]

  debug-tools:
    image: nicolaka/netshoot
    profiles: ["debug"]

  prometheus:
    image: prom/prometheus
    profiles: ["monitoring"]
```

```bash
# Start core services only
docker compose up

# Start with workers
docker compose --profile full up

# Start with monitoring
docker compose --profile monitoring up

# Start everything
docker compose --profile full --profile monitoring --profile debug up
```

### 5.3 Depends-On with Health Checks

```yaml
services:
  api:
    depends_on:
      db:
        condition: service_healthy
        restart: true              # Restart api if db restarts
      redis:
        condition: service_healthy
      migrations:
        condition: service_completed_successfully
```

Without `condition: service_healthy`, `depends_on` only waits for the container to start, not for the service inside it to be ready. This is one of the most common Docker Compose mistakes.

### 5.4 Compose Override Files

```bash
# compose.yaml           -- base configuration
# compose.override.yaml  -- local development overrides (auto-merged)
# compose.prod.yaml      -- production overrides

# Development (auto-merges compose.override.yaml)
docker compose up

# Production
docker compose -f compose.yaml -f compose.prod.yaml up -d

# Validate merged config
docker compose -f compose.yaml -f compose.prod.yaml config
```

### 5.5 Resource Limits

```yaml
services:
  api:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 256M
      replicas: 3
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
        window: 120s
```

### 5.6 Logging Configuration

```yaml
services:
  api:
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"
        tag: "{{.Name}}/{{.ID}}"

  # Alternative: send to Loki
  worker:
    logging:
      driver: loki
      options:
        loki-url: "http://loki:3100/loki/api/v1/push"
        loki-retries: "5"
        loki-batch-size: "400"
```

### 5.7 Docker Compose Watch (Hot Reload)

```yaml
services:
  api:
    build: ./api
    develop:
      watch:
        - action: sync
          path: ./api/src
          target: /app/src
        - action: rebuild
          path: ./api/package.json
```

```bash
docker compose watch
```

### 5.8 Docker Bake (Compose Spec v5.0.0)

As of December 2025, Compose Spec v5.0.0 delegates builds to Docker Bake for faster builds, better caching, and multi-platform support.

```hcl
# docker-bake.hcl
group "default" {
  targets = ["api", "worker"]
}

target "api" {
  context    = "./api"
  dockerfile = "Dockerfile"
  tags       = ["myregistry.io/api:latest"]
  platforms  = ["linux/amd64", "linux/arm64"]
  cache-from = ["type=registry,ref=myregistry.io/api:cache"]
  cache-to   = ["type=registry,ref=myregistry.io/api:cache,mode=max"]
}

target "worker" {
  context    = "./worker"
  tags       = ["myregistry.io/worker:latest"]
}
```

```bash
docker buildx bake --push
```

---

## 6. Container Networking

### 6.1 Network Drivers

| Driver | Scope | Use Case |
|---|---|---|
| `bridge` | Single host | Default; container-to-container on same host |
| `host` | Single host | Maximum network performance (no NAT) |
| `overlay` | Multi-host | Swarm/cross-node communication (VXLAN) |
| `macvlan` | Single host | Container appears as physical device on LAN |
| `ipvlan` | Single host | Like macvlan but shares host MAC address |
| `none` | Single host | Complete network isolation |

### 6.2 User-Defined Bridge Networks

```bash
# Create a custom bridge network
docker network create --driver bridge \
    --subnet 172.20.0.0/16 \
    --ip-range 172.20.240.0/20 \
    my-network

# Containers on the same user-defined bridge resolve each other by name
docker run --network my-network --name api myapp:v1
docker run --network my-network --name db postgres:16
# api can reach db at hostname "db"
```

User-defined bridges provide DNS-based service discovery. The default `bridge` network does not -- only `--link` (deprecated) works there.

### 6.3 Network Isolation in Compose

```yaml
networks:
  public:
    driver: bridge
  private:
    driver: bridge
    internal: true       # No internet access

services:
  nginx:
    networks: [public, private]
  api:
    networks: [private]
  db:
    networks: [private]
```

The `internal: true` flag prevents any external network access, ensuring the database is only reachable by services explicitly connected to that network.

### 6.4 Overlay Networks (Multi-Host)

```bash
# Initialize Swarm
docker swarm init

# Create overlay network
docker network create --driver overlay --attachable my-overlay

# Services across nodes communicate transparently
docker service create --network my-overlay --name api myapp:v1
docker service create --network my-overlay --name db postgres:16
```

Overlay networks use VXLAN encapsulation under the hood to create a virtual L2 network spanning multiple hosts.

### 6.5 Macvlan Networks

```bash
# Containers get IPs on the host network (no NAT)
docker network create -d macvlan \
    --subnet=192.168.1.0/24 \
    --gateway=192.168.1.1 \
    -o parent=eth0 \
    my-macvlan

docker run --network my-macvlan --ip 192.168.1.100 myapp:v1
```

Four modes: bridge (virtual bridge on host), VEPA (hairpin through external switch), private (containers cannot see each other), passthru (exclusive NIC access).

---

## 7. Volumes and Storage

### 7.1 Volume Types

```yaml
services:
  db:
    volumes:
      # Named volume (Docker-managed, preferred for persistence)
      - pgdata:/var/lib/postgresql/data

      # Bind mount (host path, good for development)
      - ./init-scripts:/docker-entrypoint-initdb.d:ro

      # tmpfs (in-memory, good for sensitive temp data)
      - type: tmpfs
        target: /tmp
        tmpfs:
          size: 64m
          mode: 1777

volumes:
  pgdata:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/postgres
```

### 7.2 Volume Best Practices

- Use **named volumes** for database data, not bind mounts
- Mark host-mounted config as `:ro` (read-only) when the container should not write to it
- Use `tmpfs` for ephemeral scratch space and sensitive temp files
- Back up named volumes with `docker run --rm -v pgdata:/data -v $(pwd):/backup alpine tar czf /backup/pgdata.tar.gz /data`
- Never store application state in the container's writable layer

### 7.3 Volume Drivers

```yaml
volumes:
  shared-data:
    driver: local
    driver_opts:
      type: nfs
      o: addr=nfs-server.local,rw
      device: ":/exports/data"

  encrypted-data:
    driver: local
    driver_opts:
      type: tmpfs
      o: "size=100m"
```

---

## 8. Container Registries

### 8.1 Registry Comparison

| Registry | Type | Scanning | Free Tier | OCI v2 |
|---|---|---|---|---|
| Docker Hub | Public SaaS | Scout | 1 private repo | Yes |
| GitHub (GHCR) | Public SaaS | Dependabot | Unlimited | Yes |
| AWS ECR | Cloud | Inspector | 500 MB/mo | Yes |
| Google AR | Cloud | On-Demand | 500 MB/mo | Yes |
| Azure ACR | Cloud | Defender | None | Yes |
| Harbor | Self-hosted | Trivy/Clair | Unlimited | Yes |
| Quay | Hybrid | Clair | Unlimited (public) | Yes |

### 8.2 Registry Operations

```bash
# Docker Hub
docker login
docker tag myapp:v1 username/myapp:v1
docker push username/myapp:v1

# GitHub Container Registry (GHCR)
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
docker tag myapp:v1 ghcr.io/org/myapp:v1
docker push ghcr.io/org/myapp:v1

# AWS ECR
aws ecr get-login-password | docker login --username AWS --password-stdin 123456789.dkr.ecr.us-east-1.amazonaws.com
docker tag myapp:v1 123456789.dkr.ecr.us-east-1.amazonaws.com/myapp:v1
docker push 123456789.dkr.ecr.us-east-1.amazonaws.com/myapp:v1

# Harbor
docker login harbor.example.com
docker tag myapp:v1 harbor.example.com/project/myapp:v1
docker push harbor.example.com/project/myapp:v1
```

### 8.3 Image Lifecycle Management

```bash
# List tags
docker images myapp --format "{{.Tag}} {{.Size}} {{.CreatedAt}}"

# Remove old images
docker image prune -a --filter "until=720h"   # Remove images older than 30 days

# Registry garbage collection (Harbor)
# Harbor admin UI: Configuration > Garbage Collection > Schedule
```

Tag strategy: `<app>:<semver>-<sha>` (e.g., `myapp:1.2.3-abc1234`). Always push both the specific tag and a mutable convenience tag (`myapp:1.2`, `myapp:latest`).

---

## 9. Podman Compatibility

### 9.1 Architecture Differences

| Feature | Docker | Podman |
|---|---|---|
| Daemon | dockerd (privileged) | Daemonless |
| Rootless | Opt-in | Default |
| CLI | docker | podman (drop-in) |
| Compose | docker compose | podman compose (via podman-compose or docker-compose) |
| Socket | /var/run/docker.sock | Per-user socket |
| Init system | N/A | Quadlet (systemd) |
| Kubernetes | N/A | `podman generate kube` / `podman play kube` |
| Pods | N/A | Native pod support |

### 9.2 CLI Compatibility

```bash
# Drop-in alias
alias docker=podman

# Most commands work identically
podman build -t myapp .
podman run -d --name api -p 3000:3000 myapp
podman ps
podman logs api
podman exec -it api sh
podman stop api
podman rm api
```

**Compatibility gaps:**
- Docker Buildx multi-platform builds (limited in Podman)
- Docker Swarm (not supported; use Kubernetes instead)
- Docker Compose compatibility is ~95% via podman-compose

### 9.3 Rootless Containers

```bash
# Podman runs rootless by default
podman run --rm alpine id
# uid=0(root) gid=0(root) -- root INSIDE the container
# but mapped to unprivileged user OUTSIDE via user namespaces

# Check user namespace mapping
podman unshare cat /proc/self/uid_map
```

Rootless adds ~25-30% overhead to startup (irrelevant for long-running services, measurable in CI/CD with many short-lived containers).

### 9.4 Quadlet (Systemd Integration)

Generate systemd service units from container definitions.

```ini
# ~/.config/containers/systemd/myapp.container
[Container]
Image=myapp:v1.2.3
PublishPort=3000:3000
Volume=myapp-data:/data
Environment=NODE_ENV=production
HealthCmd=wget -qO- http://localhost:3000/health
HealthInterval=30s

[Service]
Restart=always

[Install]
WantedBy=default.target
```

```bash
systemctl --user daemon-reload
systemctl --user start myapp
systemctl --user enable myapp
```

### 9.5 Podman Pods

```bash
# Create a pod (shared network namespace, like Kubernetes)
podman pod create --name my-pod -p 8080:80

podman run -d --pod my-pod --name nginx nginx:1.27-alpine
podman run -d --pod my-pod --name api myapp:v1

# All containers in the pod share localhost
# Generate Kubernetes YAML from running pod
podman generate kube my-pod > deployment.yaml
```

### 9.6 Memory Comparison

Tests with 50 concurrent nginx containers show Podman consuming 15-20% less total memory than Docker due to elimination of the daemon process and connection overhead.

---

## 10. OCI Standards

### 10.1 OCI Specifications

The Open Container Initiative (OCI) defines three specifications:

| Spec | Purpose | Status |
|---|---|---|
| Image Spec | How to package container images (layers, config, manifest) | v1.1 (stable) |
| Runtime Spec | How to run containers (lifecycle, namespaces, cgroups) | v1.2 (stable) |
| Distribution Spec | How to push/pull images from registries | v1.1 (stable) |

Both Docker and Podman produce identical OCI images and can run each other's containers. Any OCI-compliant runtime (runc, crun, youki, gVisor, Kata) can execute OCI images.

### 10.2 OCI Artifacts

OCI registries can store non-image artifacts using the same distribution APIs:

- Helm charts
- WASM modules
- SBOMs
- Signatures (Cosign)
- Policy bundles (OPA)

```bash
# Push a Helm chart to an OCI registry
helm push mychart-0.1.0.tgz oci://ghcr.io/org/charts

# Push an arbitrary artifact with ORAS
oras push ghcr.io/org/artifacts/config:v1 config.yaml:application/yaml
```

---

## 11. Container Orchestration Overview

### 11.1 Docker Swarm vs Kubernetes

| Feature | Docker Swarm | Kubernetes |
|---|---|---|
| Complexity | Low (built into Docker) | High (separate control plane) |
| Setup time | Minutes | Hours to days |
| Scaling | Manual or simple rules | HPA, VPA, KEDA, custom metrics |
| Networking | Overlay + routing mesh | CNI plugins (Calico, Cilium, Flannel) |
| Storage | Volume drivers | CSI drivers, PV/PVC |
| Service discovery | Built-in DNS | CoreDNS + Services/Ingress |
| Secret management | Docker secrets | Kubernetes Secrets + external (Vault) |
| Community | Declining | Dominant |
| Best for | Small teams, simple stacks | Large-scale, complex microservices |

### 11.2 When to Use What

- **Docker Compose:** Single-host development and small production deployments
- **Docker Swarm:** Multi-node orchestration for teams that find Kubernetes too complex; can scale to hundreds of nodes
- **Kubernetes:** Large-scale production, complex microservices, multi-cloud, when you need the ecosystem (Istio, ArgoCD, Prometheus, etc.)

### 11.3 Migration Path

```text
Docker Compose (single host)
    |
    v  (need multi-host?)
Docker Swarm (simple multi-node)
    |
    v  (need advanced scheduling, service mesh, custom controllers?)
Kubernetes (full orchestration)
```

Podman can generate Kubernetes YAML from running containers/pods, easing the transition.

---

## 12. Anti-Patterns and Common Mistakes

### 12.1 Image Anti-Patterns

| Anti-Pattern | Problem | Fix |
|---|---|---|
| Using `:latest` in production | Non-deterministic deployments | Pin specific version + digest |
| Fat base images | Bloated images, larger attack surface | Use Alpine, distroless, or scratch |
| Single-stage builds | Build tools in production image | Use multi-stage builds |
| `COPY . .` before `RUN npm ci` | Cache invalidation on every code change | Copy lock files first, then source |
| No `.dockerignore` | Inflated build context, secret leakage | Always maintain `.dockerignore` |
| Storing secrets in ENV | Visible in `docker inspect`, logs | Use BuildKit secrets or runtime mounts |
| Running as root | Privilege escalation risk | `USER nonroot` |

### 12.2 Runtime Anti-Patterns

| Anti-Pattern | Problem | Fix |
|---|---|---|
| No resource limits | Container can starve host | Set `--memory` and `--cpus` |
| Mounting Docker socket | Root-equivalent access to host | Use Kaniko, Podman, or dind with TLS |
| No health checks | Orchestrator cannot detect failures | Add `HEALTHCHECK` in Dockerfile |
| `--privileged` flag | Full host access | Use `--cap-add` for specific capabilities |
| Data in container layer | Data loss on container restart | Use named volumes |
| `depends_on` without health | Service starts before dependency is ready | Use `condition: service_healthy` |
| No logging config | Disk fills up with json-file logs | Set `max-size` and `max-file` |
| Ignoring signal handling | Ungraceful shutdown (SIGKILL after timeout) | Handle SIGTERM, use `exec` in entrypoint |

### 12.3 Build Anti-Patterns

| Anti-Pattern | Problem | Fix |
|---|---|---|
| `apt-get update` alone | Stale packages from cached layer | Combine with `apt-get install` in one RUN |
| Not cleaning apt cache | Wasted layer space | `rm -rf /var/lib/apt/lists/*` |
| Multiple `RUN pip install` | Excessive layers, cache misses | Single `RUN pip install -r requirements.txt` |
| `ADD` for local files | Unpredictable behavior (auto-extraction) | Use `COPY` for local files |
| Hard-coded config | Cannot reuse across environments | Use ENV, ARG, or config mounts |

---

## 13. Production Checklist

### 13.1 Dockerfile Checklist

- [ ] Minimal base image (Alpine, distroless, or Chainguard)
- [ ] Multi-stage build separating build from runtime
- [ ] Digest-pinned base images for production
- [ ] Non-root USER with explicit UID/GID
- [ ] HEALTHCHECK instruction defined
- [ ] `.dockerignore` excludes secrets, tests, docs, VCS
- [ ] No secrets in ENV or build args (use BuildKit secrets)
- [ ] Packages sorted alphabetically in install commands
- [ ] Cache-friendly layer ordering (dependencies before source)
- [ ] OCI labels for metadata

### 13.2 Security Checklist

- [ ] CI/CD vulnerability scanning (fail on HIGH/CRITICAL)
- [ ] SBOM generated and stored with image
- [ ] Image signed (Cosign or Docker Content Trust)
- [ ] Read-only root filesystem (`read_only: true`)
- [ ] All capabilities dropped, only needed ones added back
- [ ] `no-new-privileges` security option set
- [ ] Docker socket never mounted into application containers
- [ ] Resource limits (CPU, memory) set for all containers
- [ ] Network isolation with internal networks where appropriate
- [ ] Runtime monitoring (Falco or equivalent)

### 13.3 Compose / Deployment Checklist

- [ ] Health checks on all services
- [ ] `depends_on` with `condition: service_healthy`
- [ ] Named volumes for all persistent data
- [ ] Restart policies (`unless-stopped` or `always`)
- [ ] Logging configured with size rotation
- [ ] Resource limits in deploy section
- [ ] Secrets via Docker secrets or external manager
- [ ] `docker compose config` validates before deploy
- [ ] Environment-specific overrides via compose files

---

## 14. Quick Reference Commands

Full command cheatsheet (build, run, inspect, cleanup, compose, scanning, registry): see `references/docker-quick-reference.md`.

---

## 15. Cross-Skill References

### CI/CD Pipelines (cicd-pipelines)

Container images are a primary CI/CD artifact. Key integration points:
- Build and scan images in CI (GitHub Actions, GitLab CI)
- Push to registry with provenance attestation (SLSA)
- Use multi-stage builds to keep CI fast and images small
- Implement image promotion across environments (dev -> staging -> prod)
- Cache Docker layers in CI using registry-backed cache (`--cache-from`, `--cache-to`)
- Sign images in the pipeline with Cosign for supply chain security

```yaml
# GitHub Actions: Build, scan, push
- uses: docker/build-push-action@v6
  with:
    push: true
    tags: ghcr.io/org/myapp:${{ github.sha }}
    cache-from: type=registry,ref=ghcr.io/org/myapp:cache
    cache-to: type=registry,ref=ghcr.io/org/myapp:cache,mode=max

- name: Scan image
  run: trivy image --exit-code 1 --severity CRITICAL ghcr.io/org/myapp:${{ github.sha }}
```

### Code Packaging (code-packaging)

Container images are the universal packaging format for deployable artifacts:
- Use multi-stage builds to separate the library/app build from the runtime packaging
- Package monorepo services as individual container images using shared base stages
- Leverage BuildKit cache mounts to speed up `npm ci`, `pip install`, `go mod download`
- Publish both library artifacts (npm, PyPI) and container images from the same CI pipeline
- Use OCI artifacts to store SBOMs, signatures, and Helm charts alongside images

---

## 16. References

Official docs, security/scanning guides, image optimization, BuildKit, Podman, Compose, and anti-pattern resources: see `references/docker-references.md`.
