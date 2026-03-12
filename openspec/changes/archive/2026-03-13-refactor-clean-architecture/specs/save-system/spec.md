## MODIFIED Requirements

### Requirement: Save State Persistence
The system SHALL store and retrieve complete game states as JSON via a `SaveRepository`.

#### Scenario: Successful save
- **WHEN** client sends a valid JSON save body to the save endpoint
- **THEN** system persists the state via the repository and returns a new `save_version`

### Requirement: Save Data Caching
The system SHALL use a `CacheRepository` to cache player save data for frequent retrieval.

#### Scenario: Save cache hit
- **WHEN** client requests save data and it's present in the cache repository
- **THEN** system returns the cached state without querying the primary data store

### Requirement: Save Backups
The system SHALL provide a mechanism to create and manage snapshots of player save data via the `SaveRepository`.

#### Scenario: Manual backup
- **WHEN** client requests a manual backup
- **THEN** system creates a snapshot via the repository
