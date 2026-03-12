## 1. Database Schema

- [x] 1.1 Create `004_hollow_wilds_phase1.sql` migration with `players`, `player_saves`, `player_save_backups`, and `leaderboard_entries` tables
- [x] 1.2 Run migrations on target environments

## 2. Core Service Implementation

- [x] 2.1 Update `internal/models` with Hollow Wilds DTOs and player states
- [x] 2.2 Implement `HollowWildsService` for auth validation and save management
- [x] 2.3 Add Redis caching logic to `pkg/utils/redis.go`
- [x] 2.4 Extend `LeaderboardService` with Hollow Wilds multi-metric support

## 3. API Handlers & Routes

- [x] 3.1 Create `hollow_wilds_handler.go` with Auth, Save, and Analytics endpoints
- [x] 3.2 Update `leaderboard_handler.go` with Hollow Wilds ranking logic
- [x] 3.3 Register all Phase 1 routes in `cmd/server/main.go`
- [x] 3.4 Initialize Redis connection in `main.go`

## 4. Testing & Validation

- [x] 4.1 Write integration tests for PlayFab login flow
- [x] 4.2 Verify save/load cycle with versioning and backups
- [x] 4.3 Validate leaderboard submission and ranking accuracy
- [x] 4.4 Verify analytics batch processing
