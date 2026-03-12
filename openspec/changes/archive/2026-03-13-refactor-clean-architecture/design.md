## Context

The current project uses a flat service layer that directly manages database and Redis connections. Handlers are tightly coupled to these services. This refactor introduces Clean Architecture to provide a more modular and testable structure.

## Goals / Non-Goals

**Goals:**
- Separate project into Domain, Usecase, Infrastructure, and Delivery layers.
- Implement Dependency Injection for all components.
- Define interfaces for all data access (Repositories) and business logic (Usecases).
- Enable mock-based unit testing for usecases.

**Non-Goals:**
- Change the existing database schema or API behavior (refactor only).
- Implement a DI framework (manual constructor injection will be used for simplicity).

## Decisions

### 1. Layered Directory Structure
**Decision:** Organize `internal/` into the following sub-packages:
- `domain/`: Entities and Interface definitions (ports).
- `usecase/`: Business logic implementations.
- `infrastructure/`: Data access implementations (PostgreSQL, Redis, PlayFab).
- `delivery/`: API handlers (Fiber).

### 2. Dependency Injection Strategy
**Decision:** Use manual constructor injection in `cmd/server/main.go`.
**Rationale:** The project size doesn't yet justify the complexity of a DI container like Wire or Dig. Manual injection is explicit and easy to follow.

### 3. Interface Definitions
**Decision:** Interfaces will be defined in the `domain/` layer.
**Rationale:** This ensures the Dependency Rule—high-level logic only knows about abstractions, and implementations depend on the domain.

## Risks / Trade-offs

- **[Risk] Refactor Complexity** → **[Mitigation]** Refactor one module at a time (e.g., Auth first) and verify with existing integration tests.
- **[Trade-off] Boilerplate** → Clean Architecture requires more files (interfaces, wrappers). This is accepted for better long-term maintainability.
