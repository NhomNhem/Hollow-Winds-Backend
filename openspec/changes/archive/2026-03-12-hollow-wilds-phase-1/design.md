## Context

The backend is transitioning from a generic "GameFeel" structure to a specific "Hollow Wilds" survival game architecture. Current systems like `users` and `level_completions` are being deprecated in favor of a specialized `players` model and a robust JSON-based save system.

## Goals / Non-Goals

**Goals:**
- Implement a secure, PlayFab-integrated authentication system with JWT.
- Provide a high-performance save/load system using Redis for caching.
- Enable competitive play via a multi-metric leaderboard system.
- Support detailed survival analytics through batch event logging.

**Non-Goals:**
- Real-time multiplayer synchronization (slated for Phase 3).
- In-game shop and economy implementation (slated for Phase 2).
- Migration of legacy "GameFeel" user data to the new system.

## Decisions

### 1. Data Modeling: `players` vs. `users`
**Decision:** Create a separate `players` table instead of modifying the existing `users` table.
**Rationale:** "Hollow Wilds" has significantly different data requirements (survival stats vs. star-based level progress). Maintaining a clean break avoids schema pollution and allows both systems to coexist during transition if necessary.

### 2. Save Persistence: JSONB in PostgreSQL
**Decision:** Use `JSONB` for the `save_data` column in `player_saves`.
**Rationale:** The game state is complex and highly nested. Storing it as a blob or across dozens of tables would be inflexible and slow to develop. `JSONB` allows for fast reads/writes while still permitting indexing on specific fields if needed later.

### 3. Caching Strategy: Write-Through with Invalidation
**Decision:** Use Redis to cache save data on every read (if missing) and invalidate it on every write.
**Rationale:** Players load their save once at the start but may save frequently. Caching the full JSON blob in Redis minimizes expensive PostgreSQL hits for a deterministic single-player state.

### 4. Versioning: Optimistic Locking
**Decision:** Include a `save_version` integer in the `player_saves` table.
**Rationale:** Prevents data loss if multiple clients (or accidental concurrent calls) attempt to save state. The server will reject any update where the provided `version` doesn't match the current server `version`.

## Risks / Trade-offs

- **[Risk] Redis Downtime** → **[Mitigation]** The system is designed to fallback to the database if Redis is unavailable, albeit with higher latency.
- **[Trade-off] High Database Size** → By storing full snapshots in `player_save_backups`, storage usage will grow. **[Mitigation]** Limit the number of backups to 10 per player.
