## Why

The project is pivoting from "GameFeel" to "Hollow Wilds," a survival game with Malaysian folklore themes. This change implements the Phase 1 backend requirements to support the new game's core loop, including authentication, state persistence, and competitive features.

## What Changes

- **Auth System Update**: Transition from the old GameFeel auth to Hollow Wilds auth using PlayFab session tickets and JWT with refresh token support.
- **Player State Management**: New save/load system for complex game states (world, inventory, sebilah, base) with version control and backups.
- **Competitive Leaderboards**: Implementation of three specific leaderboard types (longest run, soul level, bosses killed) with global and per-character scopes.
- **Enhanced Analytics**: Batch event tracking for survival-specific actions (player death, item crafting, etc.).
- **Infrastructure**: Redis integration for caching save data and managing session blacklisting/rate limiting.

## Capabilities

### New Capabilities
- `auth`: PlayFab session ticket validation, JWT issuance, and refresh token management.
- `save-system`: Versioned player state persistence with automatic and manual backups.
- `leaderboard`: Multi-metric rankings with character-based filtering and personal best tracking.
- `analytics`: Batch event tracking for in-game survival metrics.

### Modified Capabilities
- None (Initial Phase 1 implementation for Hollow Wilds).

## Impact

- **Database**: New tables `players`, `player_saves`, `player_save_backups`, `leaderboard_entries`.
- **API**: New endpoints at `/api/v1/auth/hw`, `/api/v1/player`, `/api/v1/leaderboard`, and `/api/v1/analytics`.
- **Dependencies**: Redis (Upstash) for caching and session management.
- **Legacy**: Moves away from `users`, `level_completions`, and `user_talents` tables/logic from the GameFeel phase.
