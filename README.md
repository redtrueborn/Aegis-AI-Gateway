# Aegis AI Gateway

**Aegis AI Gateway** is a Go-based AI infrastructure gateway for controlling, observing, and operating LLM requests in production-style backend systems.

The project is built to explore how modern AI systems should be wrapped with real backend/platform engineering: provider abstraction, request logging, timeout handling, retries, durable async jobs, cost tracking, structured output validation, and OpenTelemetry instrumentation.

This is not a chatbot.

This is infrastructure.

---

## Project Status

Aegis is under active solo development.

The repository is public for visibility, learning, and portfolio review. It is not currently accepting external contributions while the core architecture is being established.

Current focus:

- Go backend architecture
- LLM gateway design
- durable background jobs
- request observability
- AI infrastructure fundamentals
- production-oriented service patterns

---

## Why This Exists

Calling an LLM provider directly is easy.

Operating AI workloads reliably is harder.

Production AI systems need answers to questions like:

- Which provider handled this request?
- How long did the model call take?
- How much did the request cost?
- Did the request timeout, retry, or fail?
- Was the output valid JSON?
- Which prompt and model version produced the result?
- Can slow AI work run asynchronously?
- Can failures be traced and debugged?
- Can model calls be routed, limited, and observed?

Aegis is a focused project for building those backend controls in Go.

---

## Core Goals

Aegis is designed around five engineering goals:

1. **Control** — route AI requests through one internal gateway instead of scattering provider calls across services.
2. **Reliability** — handle timeouts, retries, queue-backed workloads, and failure states.
3. **Observability** — trace requests, provider calls, database operations, and worker jobs.
4. **Cost Awareness** — track model usage, latency, tokens, and estimated cost.
5. **Extensibility** — make providers, queues, RAG, and evals modular enough to evolve.

---

## Planned Architecture

```text
Client / Application
        |
        v
+---------------------+
|  Aegis API Gateway  |
+---------------------+
        |
        |-- Auth / API Keys
        |-- Rate Limiting
        |-- Request IDs
        |-- Provider Routing
        |-- Structured Output Validation
        |-- Usage Logging
        |
        v
+---------------------+
|   Provider Layer    |
+---------------------+
        |
        |-- Mock Provider
        |-- OpenAI-compatible Provider
        |-- Future: Anthropic / Gemini / Bedrock / vLLM
        |
        v
+---------------------+
|  LLM Provider APIs  |
+---------------------+
```

Async AI workloads:

```
Client / Application
        |
        v
+---------------------+
|  Aegis API Gateway  |
+---------------------+
        |
        v
+---------------------+
| Durable Job Queue   |
+---------------------+
        |
        v
+---------------------+
|  Worker Process     |
+---------------------+
        |
        |-- Batch LLM Calls
        |-- Document Ingestion
        |-- Embedding Jobs
        |-- Eval Runs
        |-- Retry / Backoff
        |-- Dead-letter Handling
        |
        v
+---------------------+
| Usage + Telemetry   |
+---------------------+

```

### Initial Feature Roadmap

#### v0.1 — Gateway Core
---

* [] Go HTTP API
* [] POST /v1/chat
* [] mock LLM provider
* [] provider interface
* [] request IDs
* [] context timeout handling
* [] Postgres request logging
* [] Docker Compose setup
* [] basic README and architecture docs
#### v0.2 — Durable Queue

* [] Postgres-backed job table
* [] worker process
* [] enqueue async AI job
* [] retry with backoff
* [] max attempts
* [] dead-letter state
* [] job execution logs
* [] graceful worker shutdown
#### v0.3 — Observability
---

* OpenTelemetry tracing
* HTTP request spans
* provider call spans
 database operation spans
 queue job spans
 latency tracking
 error classification
 basic metrics endpoint

#### v0.4 — RAG Hook
---

 document ingestion endpoint
 text chunking
 embedding provider interface
 pgvector-backed chunk storage
 retrieval endpoint
 source metadata
 retrieval logging

#### v0.5 — Eval Harness
---

 eval case model
 prompt version tracking
 model version tracking
 expected output traits
 eval run history
 latency/cost comparison
 regression testing notes


#### Tech Stack
---

Area	Technology
Language	Go
API	Go HTTP server
Database	Postgres
Vector Search	pgvector planned
Queue	Postgres-backed durable queue
Observability	OpenTelemetry planned
Local Dev	Docker Compose
CI	GitHub Actions planned

#### Repository Structure
---

aegis-ai-gateway/
  cmd/
    api/                 # HTTP API entrypoint
    worker/              # background worker entrypoint

  internal/
    auth/                # API key/auth controls
    config/              # configuration loading
    gateway/             # request handling and gateway orchestration
    providers/           # LLM provider interfaces and implementations
    queue/               # durable async job engine
    storage/             # Postgres access layer
    telemetry/           # OpenTelemetry setup
    usage/               # request usage, tokens, cost, latency
    validation/          # structured output validation

  migrations/            # database migrations

  docs/
    architecture.md       # system architecture notes
    ai-gateway-design.md  # gateway design decisions
    queue-design.md       # durable queue design
    observability.md      # tracing/metrics design

  docker-compose.yml
  Makefile
  README.md

#### Local Development
---

Local setup will evolve as the project stabilizes.

Expected local dependencies:

Go 1.22+
Docker
Docker Compose
Postgres

Start local services:

docker compose up -d

Run database migrations:

make migrate

Run the API:

make api

Run the worker:

make worker

Run tests:

make test

#### Example API Shape
---

Chat Request
```
POST /v1/chat
Content-Type: application/json
Authorization: Bearer <api_key>
{
  "model": "mock-fast",
  "messages": [
    {
      "role": "user",
      "content": "Summarize this support ticket."
    }
  ],
  "response_format": "json",
  "metadata": {
    "tenant_id": "demo",
    "workflow": "support_triage"
  }
}
Chat Response
{
  "id": "req_123",
  "provider": "mock",
  "model": "mock-fast",
  "status": "succeeded",
  "latency_ms": 42,
  "output": {
    "summary": "Mock response generated successfully."
  }
}
```

### Data Model Draft
---

ai_requests
id
tenant_id
provider
model
prompt_version
request_hash
status
latency_ms
input_tokens
output_tokens
estimated_cost
error_type
created_at
updated_at
jobs
id
type
payload
status
attempts
max_attempts
run_after
locked_by
locked_until
idempotency_key
error_message
created_at
updated_at

#### Planned job states:
---

queued
running
succeeded
failed
dead_letter
cancelled

#### Design Principles
---

1. Keep the gateway boring

The gateway should be predictable infrastructure, not a clever framework.

2. Make failures visible

Every failed request or job should have enough metadata to debug what happened.

3. Prefer explicit state

Jobs, requests, provider calls, retries, and eval runs should have clear state transitions.

4. Use Go for production patterns

The project is intentionally written in Go to sharpen backend/platform skills:

interfaces
context cancellation
HTTP middleware
Postgres transactions
worker loops
graceful shutdown
structured errors
observability
integration testing
5. Build AI infrastructure, not AI hype

The goal is not to build another chatbot. The goal is to build the backend control layer around AI workloads.

#### What This Project Demonstrates

Aegis is intended to demonstrate practical backend/platform engineering around AI systems:

Go service architecture
LLM provider abstraction
production-style HTTP APIs
request lifecycle tracking
timeout and retry handling
durable background processing
Postgres-backed state
observability design
RAG infrastructure foundations
eval infrastructure foundations
Roadmap Notes

Near-term focus is intentionally narrow:

1. Make one LLM request controllable.
2. Make the request observable.
3. Make slow AI work asynchronous.
4. Make failures traceable.
5. Add retrieval.
6. Add evals.

The project will not start with complex model serving, GPU infrastructure, fine-tuning, agents, or Kubernetes inference. Those belong later.

Contribution Policy

Aegis is currently a solo portfolio/research project.

External contributions are not being accepted yet while the architecture is still changing.

If this changes, the repository will add:

CONTRIBUTING.md
issue templates
pull request templates
code of conduct
public roadmap labels
License

License decision pending.

During early development, this project is public for visibility and review. A formal license will be added before any stable release.

Author

Built by Alex Trew as part of a focused transition into production AI infrastructure and platform engineering.

Primary focus:

Go backend systems
AI infrastructure
observability
durable async workloads
production service design

Use this now. Do not over-polish it. The README only matters if the repo starts getting real commits behind it.
::contentReference[oaicite:1]{index=1}
