## Why

Hollow Wilds requires a backend API to support player authentication, save/load functionality, and analytics tracking for Early Access launch. The current GameFeel backend provides a foundation (Fiber, Supabase, Redis) but needs new endpoints and data models specific to Hollow Wilds' survival gameplay, deterministic world generation, and progression systems.

## What Changes

- New auth endpoints with refresh token support (extends existing JWT auth)
- Player save/load endpoints with version conflict detection
- Save backup system for data recovery
- Extended analytics events for survival gameplay tracking
- New database tables: `players`, `player_saves`, `player_save_backups`, `analytics_events`
- Redis caching for save data and session management

## Capabilities

### New Capabilities

- `auth-refresh`: JWT refresh token flow for extended sessions without re-authenticating via PlayFab
- `player-save`: Save and load game state with version control and conflict detection
- `save-backup`: Manual and automatic backup management for save data recovery
- `analytics-extended`: Extended event tracking for survival gameplay (death, crafting, exploration)

### Modified Capabilities

- `auth-login`: **MODIFIED** - Change from PlayFab session token to `playfab_session_ticket` parameter format

## Impact

- **Database**: New tables (`players`, `player_saves`, `player_save_backups`, extended `analytics_events`)
- **API**: New endpoints under `/api/v1/auth/refresh`, `/api/v1/player/save`, `/api/v1/analytics/events`
- **Models**: New data models for save data, backups, and extended analytics
- **Services**: New save service with version control and conflict resolution
- **Cache**: Redis caching for save data (5-minute TTL) and session blacklisting
- **Existing**: Auth login endpoint parameter format change (backward incompatible)
