## ADDED Requirements

### Requirement: Create Manual Backup
The system SHALL create a backup of the current save data when requested by the player.

#### Scenario: Successful backup creation
- **WHEN** authenticated player requests backup
- **THEN** system creates backup with current save_version and returns backup_id

#### Scenario: Backup limit reached
- **WHEN** player already has 10 backups
- **THEN** system deletes oldest backup before creating new one

#### Scenario: No save data to backup
- **WHEN** player requests backup but has no save data
- **THEN** system returns 404 error

### Requirement: List Player Backups
The system SHALL return a list of all backups for the authenticated player.

#### Scenario: Backups exist
- **WHEN** player requests backup list
- **THEN** system returns array of backups with backup_id, save_version, and created_at

#### Scenario: No backups exist
- **WHEN** player requests backup list but has none
- **THEN** system returns empty array

### Requirement: Restore from Backup
The system SHALL restore save data from a selected backup.

#### Scenario: Successful restore
- **WHEN** player requests restore with valid backup_id
- **THEN** system copies backup data to player_saves and increments save_version

#### Scenario: Invalid backup_id
- **WHEN** player requests restore with non-existent backup_id
- **THEN** system returns 404 error

#### Scenario: Backup belongs to different player
- **WHEN** player requests restore with another player's backup_id
- **THEN** system returns 404 error (do not reveal backup exists)

### Requirement: Automatic Backup on Major Update
The system SHALL create automatic backup when save_version crosses major version threshold (10, 20, 30, etc.).

#### Scenario: Major version threshold crossed
- **WHEN** save_version changes from 9 to 10
- **THEN** system creates automatic backup before updating save

#### Scenario: Non-major version update
- **WHEN** save_version changes from 10 to 11
- **THEN** system does NOT create automatic backup
