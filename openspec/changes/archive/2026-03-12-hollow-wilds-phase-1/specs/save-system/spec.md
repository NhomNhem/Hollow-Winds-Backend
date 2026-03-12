## ADDED Requirements

### Requirement: Save State Persistence
The system SHALL store and retrieve complete game states as JSON for each player.

#### Scenario: Successful save
- **WHEN** client sends a valid JSON save body to the save endpoint
- **THEN** system persists the state to the `player_saves` table and returns a new `save_version`

#### Scenario: Version conflict
- **WHEN** client attempts to save with an outdated `version` number
- **THEN** system returns a `version_conflict` (409) error

### Requirement: Save Data Caching
The system SHALL use Redis to cache player save data for frequent retrieval.

#### Scenario: Save cache hit
- **WHEN** client requests save data and it's present in Redis
- **THEN** system returns the cached state without querying the database

#### Scenario: Cache invalidation
- **WHEN** a player successfully updates their save data
- **THEN** system clears the associated Redis cache entry

### Requirement: Save Backups
The system SHALL provide a mechanism to create and manage snapshots of player save data.

#### Scenario: Manual backup
- **WHEN** client requests a manual backup
- **THEN** system copies the current state to the `player_save_backups` table

#### Scenario: Automatic backup
- **WHEN** a player's save version increments by 10
- **THEN** system automatically creates a snapshot in the backups table
