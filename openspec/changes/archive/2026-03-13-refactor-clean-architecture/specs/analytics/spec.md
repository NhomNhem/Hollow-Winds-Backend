## MODIFIED Requirements

### Requirement: Batch Event Submission
The system SHALL support the ingestion of multiple analytics events via an `AnalyticsRepository`.

#### Scenario: Successful batch
- **WHEN** client submits a valid array of events
- **THEN** system persists the events via the repository and returns the count of accepted/rejected events

### Requirement: Event Context Logging
The system SHALL capture and store survival events in the `AnalyticsRepository`.

#### Scenario: Player death log
- **WHEN** a `player_death` event is tracked
- **THEN** system records the event context into the repository
