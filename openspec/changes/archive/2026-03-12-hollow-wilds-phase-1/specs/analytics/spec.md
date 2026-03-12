## ADDED Requirements

### Requirement: Batch Event Submission
The system SHALL support the ingestion of multiple analytics events in a single HTTP request.

#### Scenario: Successful batch
- **WHEN** client submits a valid array of events
- **THEN** system returns the count of `accepted` and `rejected` events

#### Scenario: Validation of event names
- **WHEN** client submits a batch with unknown event types
- **THEN** system rejects only those invalid events and accepts the rest

### Requirement: Event Context Logging
The system SHALL capture and store the `player_id`, `session_id`, and `timestamp` for each survival event.

#### Scenario: Player death log
- **WHEN** a `player_death` event is tracked
- **THEN** system records the cause, day_count, and character into the `analytics_events` table
