# mini-marketplace (Go)

A small backend service written in Go to demonstrate production-oriented backend engineering practices.

The project intentionally focuses on architecture, testability, and correctness rather than feature breadth. It models a simplified marketplace with a real database, clean layering, and comprehensive automated tests.

---

## What This Project Demonstrates
- Idiomatic Go service structure
- Interface-driven design and dependency injection
- Explicit error handling and context propagation
- Real persistence using SQLite (database/sql)
- Multiple layers of automated testing (unit, handler, integration)
- Production-style repository and service boundaries

## Architechture
```
cmd/server          Application entrypoint and wiring
internal/http
  └── api           HTTP handlers (transport layer)
  └── middleware    Prometheus middleware
internal/metrics    Prometheus metrics
internal/service    Business logic
internal/repository
  └── sqlite        SQLite implementation
internal/model      Domain models
```
### Design principles
- Handlers depend on interfaces, not concrete services
- Services are transport-agnostic
- Persistence is isolated behind repositories
- ```context.Context``` flows from HTTP → service → database

### API Design

The API follows REST principles:
- Resource-oriented URLs (`/products`, `/products/{id}`)
- Standard HTTP methods (GET, POST, PUT, PATCH, DELETE)
- Stateless requests
- Consistent JSON responses
- Clear status codes

## Implemented Features
- JSON API with proper status codes
- SQLite-backed repository
- Context-aware database queries with timeouts
- Dependency injection via interfaces
- CI pipeline (Github Actions)

### Transactions
All write operations are executed within database transactions to ensure atomicity and consistency, even for multi-step operations such as update and delete

### Concurrency control
The database layer uses a semaphore to limit concurrent access, preventing connection exhaustion and reducing errors under heavy load. The concurrency limit can be configured via the SEM_MAX environment variable.

### Observability
The service exposes Prometheus-compatible metrics at ```/metrics```, including request counts, latency histograms, in-flight requests, semaphore usage and Go runtime metrics.

## Testing
- Unit tests (table-driven)
- HTTP handler tests (```httptest```)
- Integration tests with in-memory SQLite
- Error-path and cancellation coverage
- Load testing with containers

## Running the project
```bash
make run
```
You can open a demo UI in your browser:
```
http://localhost:8080/
```

## Running Tests
```bash
go test ./...
```
With race detector
```bash
go test -race ./...
```

You can also load up a container environment with limited resources and Locust for load testing
```bash
make up
```
You can then open ```http://localhost:8089``` for Locust interface and ```http://localhost:3000``` for Grafana dashboards. Login with ```admin``` ```admin```

## Why This Project
This repository exists to show how a small Go service can be structured in a realistic, production-ready way, even at a limited scope.

The emphasis is on clear boundaries, correctness, and testability, which scale well as complexity grows.

### Notes
- SQLite is used for simplicity; the repository abstraction allows replacing it without changes to handlers or services.
- Integration tests use an in-memory database for speed and determinism.

## What I Learned
- Designing testable Go code using interfaces
- Propagating context across service boundaries
- Structuring repositories around real database behavior
- Writing meaningful tests beyond the happy path
