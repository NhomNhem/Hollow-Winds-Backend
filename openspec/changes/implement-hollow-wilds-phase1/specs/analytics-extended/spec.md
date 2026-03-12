## ADDED Requirements

### Requirement: Submit Analytics Events
The system SHALL accept batch submission of analytics events for survival gameplay tracking.

#### Scenario: Successful event submission
- **WHEN** client submits valid analytics events
- **THEN** system stores events and returns count of accepted events

#### Scenario: Mixed valid and invalid events
- **WHEN** client submits batch with some invalid events
- **THEN** system accepts valid events, rejects invalid ones, returns counts

#### Scenario: Invalid event schema
- **WHEN** client submits event missing required fields (event_name, timestamp)
- **THEN** system rejects that specific event with validation error

### Requirement: Standard Event Types
The system SHALL support the following standard event types: session_start, session_end, player_death, player_spawn, item_crafted, item_used, enemy_killed, boss_killed, biome_entered, poi_discovered, sebilah_evolved, build_changed, base_object_placed, quest_flag_set.

#### Scenario: player_death event
- **WHEN** player dies in game
- **THEN** client submits event with cause, day_count, biome, character, build in payload

#### Scenario: item_crafted event
- **WHEN** player crafts an item
- **THEN** client submits event with item_id and materials_used array in payload

#### Scenario: boss_killed event
- **WHEN** player defeats a boss
- **THEN** client submits event with boss_id, character, build, day_count in payload

### Requirement: Event Payload Flexibility
The system SHALL accept arbitrary JSON payload for each event type to support evolving analytics needs.

#### Scenario: New payload fields
- **WHEN** client submits event with new payload fields
- **THEN** system accepts and stores the new fields without schema validation

#### Scenario: Empty payload
- **WHEN** client submits event with empty payload object
- **THEN** system accepts event with null/empty payload

### Requirement: Session Tracking
The system SHALL associate events with session_id for player session analysis.

#### Scenario: Events with session_id
- **WHEN** client submits events with session_id
- **THEN** system links all events to that session for querying

#### Scenario: Events without session_id
- **WHEN** client submits events without session_id
- **THEN** system accepts events but marks session as unknown

### Requirement: Rate Limiting on Analytics
The system SHALL rate limit analytics submissions to prevent abuse.

#### Scenario: Normal submission rate
- **WHEN** client submits 10 events per minute
- **THEN** system accepts all events

#### Scenario: Excessive submission rate
- **WHEN** client submits 100+ events per minute
- **THEN** system returns 429 rate limit error
