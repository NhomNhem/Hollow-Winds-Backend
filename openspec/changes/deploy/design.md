## Context

The backend has been upgraded with Hollow Wilds Phase 1 features and verified in local mock mode. We now need to transition to a live production environment using Fly.io for application hosting and Supabase for the PostgreSQL database.

## Goals / Non-Goals

**Goals:**
- Deploy the updated containerized backend to Fly.io.
- Ensure production Supabase matches the new Phase 1 schema.
- Securely configure all production secrets.

**Non-Goals:**
- Automated CI/CD pipeline setup (focus is on the immediate manual deployment).
- Load testing or performance tuning.

## Decisions

### 1. Hosting Platform
**Decision:** Continue using Fly.io.
**Rationale:** The project already contains a `fly.toml` and Dockerfile. Fly.io provides a straightforward path for containerized Go applications with easy scaling and secret management.

### 2. Secret Management
**Decision:** Use `fly secrets set` for all sensitive environment variables.
**Rationale:** Avoids committing secrets to version control and ensures they are injected securely into the production runtime.

### 3. Database Migration Strategy
**Decision:** Execute the `004_hollow_wilds_phase1.sql` migration manually via the Supabase SQL Editor or the migration tool pointing to the production URL.
**Rationale:** Given the Phase 1 nature of the pivot, a controlled manual execution allows for immediate verification and rollback if schema conflicts occur with existing "GameFeel" tables.

## Risks / Trade-offs

- **[Risk] Production Data Conflict** → **[Mitigation]** The new tables (`players`, `player_saves`) are additive and do not modify existing generic `users` tables, reducing risk. A full database backup will be taken via Supabase before execution.
- **[Risk] Connectivity Issues** → **[Mitigation]** Use the `/health` endpoint to verify database and Redis connectivity immediately after deployment.
