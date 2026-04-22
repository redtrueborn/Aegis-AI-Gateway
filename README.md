# Iron Queue — GitLab Backlog, Architecture Study Plan, and Reference Stack

## Goal

Build a production-style job orchestration service in Go with Postgres, Terraform, and Ansible.

## Product Definition

Iron Queue is a control-plane service that:

* accepts jobs over HTTP
* persists state in Postgres
* dispatches work to worker agents
* supports retries, cancellation, timeouts, and leases
* exposes logs, metrics, and health endpoints
* is provisioned with Terraform and configured/deployed with Ansible

## Lean GitLab Operating Model

Use the lightest structure that still gives control.

### Hierarchy

* **Epic** = phase or major system area
* **Issue** = concrete deliverable that can be completed in a few focused sessions
* **Task** = execution step inside the issue

Do **not** create separate issues for every tiny implementation action unless multiple people need to own them independently.

### Daily Operating Rule

* plan at the **epic** and **issue** level
* execute at the **task** level
* review progress from the **board**, not from a giant issue list

### Suggested Labels

Keep labels minimal and functional.

**Workflow**

* status::backlog
* status::ready
* status::in-progress
* status::review
* status::done

**Component**

* component::api
* component::db
* component::scheduler
* component::worker
* component::deploy
* component::observability
* component::docs

**Priority**

* priority::p0
* priority::p1
* priority::p2

### Suggested Milestones

1. M1 Foundations
2. M2 Persistence and API
3. M3 Scheduler and Worker
4. M4 Hardening and Operations
5. M5 Infrastructure and Deployment
6. M6 Cleanup and Review

### Suggested Board Columns

* Backlog
* Ready
* In Progress
* Review
* Done

### Board Filter Strategy

Use one main board filtered to this project and these workflow labels. Only move issues across columns. Tasks stay inside the issue.

### Issue Template

Use this for every implementation issue.

```md
## Outcome
What will exist when this issue is done?

## Scope
- in scope
- in scope
- in scope

## Out of scope
- not included
- not included

## Tasks
- [ ] task 1
- [ ] task 2
- [ ] task 3

## Definition of done
- [ ] code implemented
- [ ] tests added or updated
- [ ] docs updated if needed
- [ ] manually verified

## Notes
Key design constraints, edge cases, or references.
```

### Epic Template

Use this for each phase epic.

```md
## Goal
What major outcome does this epic produce?

## Included issues
- issue 1
- issue 2
- issue 3

## Exit criteria
- measurable result
- measurable result
- measurable result

## Risks
- risk 1
- risk 2
```

### Working Rule for This Project

* **Epics** should stay few and stable
* **Issues** should be the main planning unit
* **Tasks** should contain implementation steps
* avoid label sprawl
* avoid duplicate tracking in external docs once work is in GitLab

---

## Lean GitLab Setup for Iron Queue

## Epics

### Epic: M1 Foundations

**Goal:** Service skeleton, local dev setup, lifecycle control, and base docs.

**Issues**

1. Bootstrap repository and project layout
2. Add configuration system
3. Add logging and graceful shutdown
4. Build base HTTP server and foundational endpoints
5. Stand up local development stack

### Epic: M2 Persistence and API

**Goal:** Durable schema, migrations, repositories, and core job API.

**Issues**

1. Design core data model
2. Add migration system
3. Implement repository layer
4. Build create/read/list job API
5. Build cancel and retry API

### Epic: M3 Scheduler and Worker

**Goal:** Job state machine, lease/claim loop, worker execution, heartbeats, and recovery.

**Issues**

1. Implement scheduler state machine
2. Implement lease acquisition and dispatch loop
3. Implement retry and backoff policy
4. Build worker agent
5. Add heartbeats and timeout supervision
6. Add cancellation propagation

### Epic: M4 Hardening and Operations

**Goal:** Error model, metrics, tests, resilience, and security baseline.

**Issues**

1. Define error model and API semantics
2. Add observability and metrics
3. Implement testing strategy
4. Run load and chaos drills
5. Add security baseline

### Epic: M5 Infrastructure and Deployment

**Goal:** Provision host with Terraform and deploy/configure app with Ansible.

**Issues**

1. Create Terraform project layout
2. Provision base infrastructure
3. Create Ansible inventory and roles
4. Automate application deployment
5. Add deployment verification and rollback notes

### Epic: M6 Cleanup and Review

**Goal:** Architecture clarity, operational readiness, and maintainability.

**Issues**

1. Write ADRs
2. Review package boundaries
3. Write operational runbook

---

## Issue Breakdown

### Issue: Bootstrap repository and project layout

**Labels:** `status::ready`, `priority::p0`, `component::docs`

**Tasks**

* [ ] choose module name
* [ ] initialize Go module
* [ ] create `cmd/`, `int

### Issue: Repository bootstrap

**Outcome:** Monorepo or single service repo initialized with clean layout.

**Tasks**

* choose module name
* initialize Go module
* create top-level layout: `cmd/`, `internal/`, `migrations/`, `deploy/terraform/`, `deploy/ansible/`, `docs/`
* add Makefile or task runner
* add `.editorconfig`, `.gitignore`, `.env.example`
* add README with local dev commands
* add pre-commit or lint script

**Definition of done**

* repo builds locally
* README gets a new developer to first successful run

### Issue: Configuration system

**Outcome:** Service loads config from environment deterministically.

**Tasks**

* define config struct
* implement env parsing and validation
* define defaults vs required settings
* add startup config validation
* document all env vars

### Issue: Logging and process lifecycle

**Outcome:** Service starts cleanly and shuts down cleanly.

**Tasks**

* add structured logging
* wire startup/shutdown logs
* implement signal handling
* implement graceful HTTP shutdown
* verify background goroutines stop on shutdown

### Issue: Basic HTTP server

**Outcome:** HTTP server with foundational endpoints.

**Tasks**

* create router setup
* add `/healthz`
* add `/readyz`
* add `/version`
* add request logging middleware
* add request ID middleware
* add timeout middleware

### Issue: Local development stack

**Outcome:** Local service and Postgres start with one command.

**Tasks**

* add Docker Compose for Postgres
* add local run instructions
* add DB bootstrap script
* verify service can connect locally

---

## Epic 2 — Persistence and API

### Issue: Core data model design

**Outcome:** Durable schema for jobs, attempts, workers, and audit events.

**Tasks**

* design ERD
* define job states: queued, leased, running, succeeded, failed, canceled, dead
* define attempt model
* define worker lease model
* define event log model
* define indexes by lookup path
* document invariants and state transitions

### Issue: Migration system

**Outcome:** Schema is reproducible and versioned.

**Tasks**

* choose migration tool or build simple SQL migration runner
* create initial schema migration
* add rollback strategy
* add migration apply command
* add migration status command

### Issue: Repository layer

**Outcome:** Database access isolated behind interfaces.

**Tasks**

* define repository interfaces
* implement job repository
* implement attempt repository
* implement worker/lease repository
* implement event repository
* use `context.Context` on all methods
* add transaction boundaries where needed

### Issue: Job create/read/list API

**Outcome:** Users can create and inspect jobs.

**Tasks**

* design request/response types
* implement POST `/jobs`
* implement GET `/jobs/{id}`
* implement GET `/jobs`
* add validation and error mapping
* add pagination for listing
* add tests for API handlers

### Issue: Cancel and retry API

**Outcome:** External operators can alter job execution state safely.

**Tasks**

* implement POST `/jobs/{id}/cancel`
* implement POST `/jobs/{id}/retry`
* enforce state transition rules
* append event log entries
* test illegal transitions

---

## Epic 3 — Scheduler and Worker Runtime

### Issue: Scheduler state machine

**Outcome:** Single source of truth for valid state transitions.

**Tasks**

* define transition matrix
* implement domain methods for transition checks
* document terminal vs retriable failure
* test all transitions

### Issue: Lease acquisition and dispatch loop

**Outcome:** Scheduler can safely claim queued jobs.

**Tasks**

* implement polling loop
* implement DB claim logic in transaction
* prevent duplicate lease acquisition
* set lease expiry timestamps
* record attempt creation
* test concurrent claim scenarios

### Issue: Retry and backoff policy

**Outcome:** Failed jobs retry predictably.

**Tasks**

* define retry classification
* define max attempts
* define exponential or bounded backoff
* schedule next attempt time
* send exhausted jobs to dead state
* test policy behavior

### Issue: Worker agent

**Outcome:** Independent process can request and execute work.

**Tasks**

* define worker registration identity
* implement poll for work
* implement command/task execution abstraction
* report success/failure
* report start and finish timestamps
* persist output or output summary

### Issue: Heartbeats and timeout supervision

**Outcome:** Stalled jobs are detected and recovered.

**Tasks**

* add worker heartbeat endpoint or DB update path
* implement heartbeat loop
* implement supervisor scan for expired leases
* transition stale jobs back to queued or failed
* add tests for orphaned work recovery

### Issue: Cancellation propagation

**Outcome:** Cancel requests actually stop work.

**Tasks**

* propagate context cancellation to worker execution
* define best-effort vs guaranteed cancellation semantics
* update final state on cancel
* test race between completion and cancellation

---

## Epic 4 — Hardening and Operations

### Issue: Error model and API semantics

**Outcome:** Consistent operational behavior under failure.

**Tasks**

* define internal error categories
* define API error response schema
* map DB, validation, and domain errors to HTTP status codes
* add correlation IDs to error logs

### Issue: Observability

**Outcome:** Operators can see system health and behavior.

**Tasks**

* define counters, gauges, and histograms
* add metrics endpoint
* add queue depth metric
* add job latency metric
* add retry count metric
* add worker heartbeat freshness metric

### Issue: Testing strategy

**Outcome:** Confidence in correctness under concurrency and failure.

**Tasks**

* add unit tests for domain/state machine
* add repository integration tests against Postgres
* add handler tests with `httptest`
* add scheduler race/concurrency tests
* run `go test -race`
* add failure injection cases

### Issue: Load and chaos drills

**Outcome:** Basic confidence under stress.

**Tasks**

* create seed workload generator
* run burst enqueue test
* simulate worker crash during execution
* simulate DB connection interruption
* measure recovery behavior
* document findings and fixes

### Issue: Security baseline

**Outcome:** Service is not careless by default.

**Tasks**

* minimize sensitive data in logs
* validate and bound request payloads
* define auth placeholder or simple token auth
* set DB least-privilege role if possible
* review shell execution safety if commands are supported

---

## Epic 5 — Infrastructure and Deployment

### Issue: Terraform project layout

**Outcome:** Infrastructure is defined as code with sane boundaries.

**Tasks**

* create root module
* create modules for compute, network, firewall/security rules, DNS if needed
* define variables and outputs
* configure remote state backend
* document state locking approach

### Issue: Base infrastructure provisioning

**Outcome:** VM and network are provisioned reproducibly.

**Tasks**

* provision VM
* provision networking/firewall
* provision static IP if needed
* provision DNS record if needed
* verify SSH access path

### Issue: Ansible inventory and role structure

**Outcome:** Host configuration is modular and repeatable.

**Tasks**

* create inventory structure
* create base role
* create app deploy role
* create systemd role or task set
* create reverse proxy role if needed
* add templates for env file and service unit

### Issue: Application deployment automation

**Outcome:** A fresh host can be configured and deployed end to end.

**Tasks**

* install required packages
* create service user
* create directories and permissions
* upload binary or build artifact
* render config
* install systemd service
* enable/start service
* restart on config change with handlers

### Issue: Deployment verification

**Outcome:** Deployment has objective verification.

**Tasks**

* verify health endpoint after deploy
* verify DB connectivity
* verify logs accessible on host
* verify service restart behavior
* document rollback steps

---

## Epic 6 — Architecture Review and Cleanup

### Issue: Architecture decision records

**Outcome:** Major decisions are written down.

**Tasks**

* ADR: why `net/http` over framework
* ADR: why `database/sql` over ORM
* ADR: polling vs LISTEN/NOTIFY wakeups
* ADR: job lease model and stale recovery
* ADR: retry policy
* ADR: single binary vs split services

### Issue: Package and boundary review

**Outcome:** Codebase structure reflects domain boundaries instead of accidental coupling.

**Tasks**

* review dependency direction
* remove package cycles
* tighten interfaces
* isolate transport from domain and persistence
* review config ownership

### Issue: Operational runbook

**Outcome:** Human operator can run system under incident conditions.

**Tasks**

* write startup/shutdown steps
* write deployment steps
* write common failure scenarios
* write queue stuck troubleshooting
* write DB issue troubleshooting
* write worker stale lease troubleshooting

---

## Architecture Materials — What to Study and Why

### 1. Web / Queue / Worker architecture

Study this first. It maps directly to the shape of the system: API layer, durable queue/state, worker execution.

Questions to answer while studying:

* where does the request end and async work begin?
* what must be durable before acknowledging job creation?
* what is the retry boundary?
* what components scale independently?

### 2. Scheduler-Agent-Supervisor pattern

This is the real architectural core for long-running or distributed task execution.

Questions to answer:

* who owns durable state?
* who owns progress tracking?
* how are stuck workers detected?
* how are partial failures recovered?

### 3. Queue-based load leveling

This helps you reason about load smoothing, decoupling, and why async execution exists.

Questions to answer:

* what load spikes should be absorbed by the queue?
* how does queue depth affect operator decisions?
* how do workers scale with backlog?

### 4. Software design for asynchronous jobs

This frames reliability, utilization, and scheduling efficiency.

Questions to answer:

* how are job classes separated?
* how do deadlines and retries affect capacity?
* where are expensive resources kept busy vs idle?

### 5. Data architecture for orchestration systems

Study your own schema like an architecture artifact.

Questions to answer:

* what are the write-hot tables?
* what indexes are mandatory?
* what transitions require transactions?
* what data is authoritative versus derived?

---

## Architecture Artifacts to Produce

1. Context diagram
2. Container diagram
3. Sequence diagram for job submission
4. Sequence diagram for lease acquisition and execution
5. State machine diagram for job lifecycle
6. ERD for Postgres schema
7. Deployment diagram for VM, DB access, and service
8. Runbook for failure recovery
9. ADR set for major technical decisions

---

## Reference Stack

### Go

* A Tour of Go
* Effective Go
* Go documentation index
* `context` package docs
* `database/sql` package docs
* Go database access guides
* `Querying for data` guide

### PostgreSQL

* Official PostgreSQL documentation
* transactions and locking
* indexes and query planning
* EXPLAIN
* advisory locks
* LISTEN/NOTIFY

### Terraform

* Terraform documentation
* configuration language overview
* state
* backends and state locking
* modules overview

### Ansible

* Ansible playbooks guide
* Ansible concepts
* roles
* handlers
* strategies
* sample setup

### GitLab

* epics
* issues
* tasks
* work items
* milestones
* labels
* issue boards
* roadmap
* description templates

### Architecture references

* Web-Queue-Worker architecture style
* Scheduler-Agent-Supervisor pattern
* Queue-based load leveling pattern
* AWS cloud design patterns
* AWS Well-Architected guidance for asynchronous and scheduled jobs

---

## Suggested Build Sequence

1. Repo bootstrap
2. Config, logging, graceful shutdown
3. Basic HTTP server
4. Local Postgres + migrations
5. Data model and repositories
6. Job create/read/list API
7. State machine
8. Lease claim loop
9. Worker agent
10. Heartbeats and timeouts
11. Retry/cancel semantics
12. Metrics and tests
13. Terraform provisioning
14. Ansible deployment
15. Architecture review and cleanup

---

## Definition of Overall Success

* a job can be created and durably stored
* scheduler can claim it exactly once per lease window
* worker can execute and report status
* retries and cancellation behave predictably
* stale work is recovered
* deployment is repeatable on a fresh host
* architecture decisions are documented
* another engineer can understand and operate the system from repo docs alone

