## ADDED Requirements

### Requirement: Load Player Save
The system SHALL return the complete game state for the authenticated player.

#### Scenario: Save data exists
- **WHEN** authenticated player requests save data
- **THEN** system returns full save data including world, player, inventory, sebilah, base, discovered_pois, and quest_flags

#### Scenario: No save data found
- **WHEN** player requests save data but has never saved
- **THEN** system returns 404 error with code "save_not_found"

#### Scenario: Cached save available
- **WHEN** player requests save data within 5 minutes of last save
- **THEN** system returns cached data from Redis (sub-100ms response)

### Requirement: Save Player Game State
The system SHALL persist the complete game state with version control.

#### Scenario: First save
- **WHEN** player saves for the first time
- **THEN** system creates new save record with save_version=1

#### Scenario: Subsequent save
- **WHEN** player saves with correct save_version
- **THEN** system updates save data and increments save_version

#### Scenario: Version conflict
- **WHEN** player saves with outdated save_version
- **THEN** system returns 409 error with server's current save_version

#### Scenario: Invalid save data schema
- **WHEN** player submits save data missing required fields
- **THEN** system returns 422 validation error

### Requirement: Save Data Structure
The system SHALL store game state with the following structure: world (seed, play_time_seconds, day_count), player (character, position, health, hunger, sanity, warmth), inventory (slots array, equipped_weapon), sebilah (weapon_id, soul_level, infusion_points), base (placed_objects array), discovered_pois (array), quest_flags (object).

#### Scenario: Complete save structure
- **WHEN** player submits complete save data
- **THEN** system validates all required fields are present

#### Scenario: Character validation
- **WHEN** player save has invalid character value
- **THEN** system rejects save (valid characters: RIMBA, DARA, BAYU, SARI)

### Requirement: Update Last Seen Timestamp
The system SHALL update the player's last_seen_at timestamp on save operations.

#### Scenario: Successful save
- **WHEN** player successfully saves
- **THEN** system updates player's last_seen_at to current timestamp
