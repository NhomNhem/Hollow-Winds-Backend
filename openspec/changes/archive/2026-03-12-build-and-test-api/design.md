## Context

The backend has just received a major update (Hollow Wilds Phase 1). We need to validate that these changes haven't introduced build-time or runtime regressions.

## Goals / Non-Goals

**Goals:**
- Successfully compile the Go server.
- Successfully build the Docker container.
- Run the integration test suite and pass all cases.

**Non-Goals:**
- Performance/Load testing.
- Unit testing of individual functions (covered by integration tests).

## Decisions

### 1. Build Environment
**Decision:** Use Go 1.25+ local environment for initial verification, then Docker for isolation.
**Rationale:** Local builds are faster for iteration, while Docker ensures parity with production (Fly.io).

### 2. Test Execution
**Decision:** Run tests against a local instance of the server (using `localhost:8080`).
**Rationale:** Simplest setup for verification without needing complex CI orchestration at this stage.

## Risks / Trade-offs

- **[Risk] Environment Variables** → **[Mitigation]** Ensure `configs/.env` is present or mock variables are provided during build/test.
- **[Risk] External Dependencies (Supabase/Redis)** → **[Mitigation]** Tests should handle disconnected states or use local mocks if configured.
