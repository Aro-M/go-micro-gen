<div align="center">

```
             _   _   __  __  ____    ____   _____  _   _ 
   ____ ___ | | | | |  \/  |/ ___|  / ___| | ____|| \ | |
  / __// _ \| |_| | | |\/| |\___ \ | |  _  |  _|  |  \| |
 | (_ | (_) |  _  | | |  | | ___) || |_| | | |___ | |\  |
  \__\ \___/|_| |_| |_|  |_||____/  \____| |_____||_| \_|
```

# go-micro-gen — The Ultimate Microservice Scaffolder

**One command to scaffold a production-ready Go microservice.**

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

</div>

---

## What is go-micro-gen?

`go-micro-gen` is a CLI tool that scaffolds a complete, production-ready Go microservice **in seconds**. 

Unlike a simple "Hello World" generator, `go-micro-gen` sets up an **entire ecosystem** ensuring you follow industry best practices right from the start.

### ✨ Key Features
- ✅ **Clean / Hexagonal Architecture** layout
- ✅ **OpenTelemetry** (traces + metrics) pre-configured
- ✅ **Cleanenv** configuration manager
- ✅ **PostgreSQL or MongoDB** implementations with the Repository Pattern
- ✅ **Message Broker** support (Kafka, RabbitMQ, NATS)
- ✅ **Cloud Provider Support** boilerplate (AWS & GCP)
- ✅ **Database Migrations** out of the box (with `golang-migrate`)
- ✅ **Docker & Docker-Compose** (App + DB + Redis + Prometheus + Grafana)
- ✅ **GitHub Actions / GitLab CI Pipeline** configurations
- ✅ **Standard Makefile** tools and `.golangci.yml` linter setup
- ✅ **Graceful shutdown**, unit tests, and health/readiness endpoints

---

## Architecture Overview

`go-micro-gen` enforces a Clean Architecture approach where dependencies flow inward towards the `Domain` layer. The service layer is completely agnostic to external databases and transports.

```mermaid
graph TD
    subgraph External
        HTTP[HTTP Transport]
        CLI[CLI Commands]
    end

    subgraph Adapters
        APIHandler[HTTP Handlers]
        Repo[Postgres/Mongo Repositories]
    end

    subgraph Core
        Service[Business Service]
        Domain[Domain Entities]
    end
    
    HTTP --> APIHandler
    CLI --> Service
    
    APIHandler -->|Uses| Service
    Service -->|Implements| Domain
    
    Repo -->|Implements Interface| Service
    Repo -->|Maps to| Domain

    classDef core fill:#d4edda,stroke:#28a745,stroke-width:2px;
    classDef adapters fill:#cce5ff,stroke:#007bff,stroke-width:2px;
    classDef external fill:#e2e3e5,stroke:#6c757d,stroke-width:2px;
    
    class Domain,Service core;
    class APIHandler,Repo adapters;
    class HTTP,CLI external;
```

---

## Installation

### Method 1: `go install` (Easiest)

If you have Go installed, you can simply run:

```bash
go install github.com/Aro-M/go-micro-gen@latest
```
*Make sure `$(go env GOPATH)/bin` is in your system `$PATH`.*

### Method 2: Build from source

```bash
git clone https://github.com/Aro-M/go-micro-gen
cd go-micro-gen
go build -o go-micro-gen ./cmd/...
sudo mv go-micro-gen /usr/local/bin/
```

### Verify Installation

```bash
go-micro-gen --version
go-micro-gen --help
```

---

## Quick Start

### Interactive Mode

Simply run the `generate` command and answer the prompts:

```bash
go-micro-gen generate
```

You'll see an interactive interface like this:

```text
? Service name:          order-service
? Go module path:        github.com/acme/order-service
? Architecture pattern:  [clean / hexagonal]
? Database:              [postgres / mongo / none]
? Include Redis?         (y/N)
? Include Docker setup?  (Y/n)
? Include Kubernetes?    (y/N)
? Include Helm charts?   (y/N)
? Cloud provider:        [aws / gcp / none]
? CI/CD provider:        [github / gitlab / none]
? Output directory:      ./order-service

  Service:  order-service
  Module:   github.com/acme/order-service
  Arch:     clean
  DB:       postgres
  Redis:    false
  Docker:   true
  K8s:      false
  Helm:     false
  Cloud:    aws
  CI:       github
  Output:   ./order-service

? Generate service with these settings? (Y/n) Y

🚀 Generating order-service ...
```

### Scripted / Non-Interactive Mode

You can supply arguments via flags for CI pipelines or scripts:

```bash
go-micro-gen generate \
  --name payment-service \
  --module github.com/acme/payment-service \
  --db postgres \
  --broker kafka \
  --transport grpc \
  --arch ddd \
  --cloud aws \
  --ci github \
  --redis=false \
  --docker=false \
  --k8s=false \
  --helm=false \
  --output ./payment-service \
  --yes
```

### Init Mode

If you already have an existing project (with a `go.mod`), you can inject the microservice structure directly into the current directory in various architectural options:

```bash
go-micro-gen init \
  --arch vertical \
  --db mongo \
  --broker rabbitmq \
  --transport both \
  --yes
```
This infers the module path and service name automatically and scaffolds the files within `.`.

---

## What Happens Under the Hood?

`go-micro-gen` embeds a complete set of `.tmpl` files and stitches together a service based on your selections. 

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant CLI as go-micro-gen
    participant Tmpl as Embedded Templates
    participant Disk as File System
    
    Dev->>CLI: go-micro-gen generate
    CLI->>Dev: Prompt for options (DB, Cloud, CI)
    Dev-->>CLI: Provide inputs
    
    CLI->>Tmpl: fs.WalkDir(templates)
    loop Every Template
        Tmpl-->>CLI: Parse file
        CLI->>CLI: Filter out unused tools (e.g., skip Mongo if Postgres selected)
        CLI->>CLI: Render variables (Service Name, Package, etc.)
        CLI->>Disk: Write output file securely
    end
    
    CLI-->>Dev: Print Success & Next Steps
```

---

## Structure of the Generated Service

A generated service (`order-service`) typically looks like this:

```
order-service/
│
├── cmd/
│   └── main.go                    # Entrypoint — DI & App wiring happens here
│
├── internal/                      # Internal packages
│   ├── config/                    # cleanenv mapped configs 
│   ├── domain/                    # Business objects (Entities)
│   ├── repository/                # DB interfaces & implementations
│   ├── service/                   # Pure business logic layer
│   └── transport/                 # HTTP routers and handlers using chi
│
├── pkg/                           # Public utilities
│   ├── health/                    # Liveness/Readiness probes
│   ├── logger/                    # Structured logging (log/slog)
│   └── telemetry/                 # OpenTelemetry traces & metrics
│
├── db/migrations/                 # Database schema migrations (.sql)
├── docker/                        # Docker & Docker Compose setup
├── .github/                       # CI workflows
├── .golangci.yml                  # Linter settings
├── .env.example                   # Environment variable templates
├── Makefile                       # Automation commands (build, run, test)
└── go.mod                         # Go module
```

---

## Roadmap

- [x] `go-micro-gen init` — inject structure into an existing project
- [x] `go-micro-gen add handler` — endpoint generation capability 
- [x] gRPC transport template
- [x] Kubernetes manifests Generation (Deployment + Service + ConfigMap)
- [x] Helm chart generation

---

## License

MIT © go-micro-gen contributors
