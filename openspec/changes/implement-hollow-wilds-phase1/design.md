## Context

Hollow Wilds is a survival game requiring persistent player state, authentication, and analytics. The existing GameFeel backend provides the tech stack (Fiber, Supabase, Redis) but needs extension for Hollow Wilds' specific requirements:

- Deterministic world generation (server doesn't store chunk data, only player state)
- Save version control with conflict detection for multiplayer readiness
- Extended analytics for survival gameplay metrics
- Refresh token flow for extended play sessions

**Stakeholders:**
- Players: Need reliable save/load, seamless auth
- Developers: Need analytics for balancing and bug tracking
- Operations: Need monitoring and data recovery capabilities

**Constraints:**
- Must work with Unity client using PlayFab SDK
- Save data is JSONB (flexible schema for game state evolution)
- Redis caching for performance (5-minute TTL for saves)
- Phase 1 is single-player; Phase 3 adds multiplayer

## Goals / Non-Goals

**Goals:**
- Implement auth with refresh tokens (JWT + PlayFab)
- Implement save/load with version control
- Implement save backup system
- Extend analytics for survival events
- Achieve sub-200ms response times for save operations
- Support 10,000+ concurrent players

**Non-Goals:**
- Multiplayer synchronization (Phase 3)
- Leaderboard implementation (separate change)
- Shop/economy endpoints (Phase 2)
- WebSocket real-time features (Phase 3)
- Client-side code changes (Unity integration handled separately)

## Decisions

### Save Data Storage
**Decision:** Store full game state as JSONB in single `player_saves` table

**Rationale:**
- Flexible schema accommodates game state changes without migrations
- Atomic save/load (single row per player)
- Supabase PostgreSQL handles JSONB efficiently

**Alternatives considered:**
- Normalized tables (inventory, player stats, etc.): Too rigid for evolving game design
- Document database (MongoDB): Adds complexity, PostgreSQL JSONB sufficient

### Version Conflict Strategy
**Decision:** Optimistic locking with `save_version` field

**Rationale:**
- Supports future multiplayer without major refactor
- Client handles conflict resolution (fetch latest, merge, resubmit)
- Simple to implement and test

**Alternatives considered:**
- Last-write-wins: Loses data in concurrent scenarios
- Pessimistic locking: Poor UX, requires session management

### Refresh Token Storage
**Decision:** Store refresh tokens in Redis with TTL

**Rationale:**
- Fast validation (sub-millisecond)
- Automatic expiration via Redis TTL
- Easy revocation on logout

**Alternatives considered:**
- Database storage: Slower, adds load
- Stateless JWT refresh: No revocation capability

### Backup Strategy
**Decision:** Manual backup trigger + automatic backup on major version bumps

**Rationale:**
- Manual: Player-controlled, on-demand recovery
- Automatic: Safety net for major updates
- Limited to 10 backups per player to control storage

**Alternatives considered:**
- Continuous backups: Too storage-intensive
- No backups: Poor player experience on data loss

## Risks / Trade-offs

**Risk:** JSONB schema drift as game evolves
→ **Mitigation:** Version save_data schema; add migration functions

**Risk:** Refresh token theft enables session hijacking
→ **Mitigation:** Bind tokens to device fingerprint; implement rotation

**Risk:** Save conflicts frustrate players in future multiplayer
→ **Mitigation:** Clear UI messaging; auto-retry with merge

**Trade-off:** Storing full state as JSONB makes partial updates harder
→ **Acceptable:** Full save/load is simpler; partial updates come in Phase 3

**Trade-off:** Manual backups require player action
→ **Acceptable:** Automatic backups on major milestones complements manual

**Trade-off:** Redis caching means 5-minute stale data window
→ **Acceptable:** Save data changes infrequently; cache invalidation on write

## Migration Plan

1. **Database migration**: Run SQL migrations for new tables
2. **Deploy new endpoints**: Backward-compatible (except auth parameter change)
3. **Update Unity client**: Switch to new auth parameter format
4. **Monitor**: Watch for save conflicts, auth failures
5. **Rollback**: Revert to previous auth format if critical issues

## Open Questions

- Should backups have expiration (e.g., 30 days)?
- What's the maximum save data size before rejection?
- Should analytics be batched or real-time for Phase 1?
