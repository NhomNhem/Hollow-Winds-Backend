## MODIFIED Requirements

### Requirement: Leaderboard Type Management
The system SHALL support and track leaderboard metrics via a `LeaderboardRepository`.

#### Scenario: Valid leaderboard query
- **WHEN** client requests a leaderboard with a supported `type`
- **THEN** system returns ranked entries from the repository matching that metric

### Requirement: Personal Best Submission
The system SHALL update a player's leaderboard entry via the `LeaderboardRepository` if the value exceeds their personal best.

#### Scenario: New personal best
- **WHEN** client submits a value higher than the stored best
- **THEN** system updates the repository and returns the new ranks
