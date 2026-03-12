## ADDED Requirements

### Requirement: Database Schema Integrity
The production Supabase database SHALL match the required schema for Hollow Wilds Phase 1.

#### Scenario: Verify tables
- **WHEN** the production database is inspected
- **THEN** tables `players`, `player_saves`, `player_save_backups`, and `leaderboard_entries` MUST exist with correct column definitions

### Requirement: Endpoint Accessibility
The production API SHALL be accessible via the public URL with HTTPS.

#### Scenario: Health check
- **WHEN** a GET request is made to `https://api.hollowwilds.com/health`
- **THEN** the system returns a 200 OK status with "ok" state
