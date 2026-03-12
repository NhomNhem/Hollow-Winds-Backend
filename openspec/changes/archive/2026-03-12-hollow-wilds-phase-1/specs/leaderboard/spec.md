## ADDED Requirements

### Requirement: Leaderboard Type Management
The system SHALL support and track three metrics: `longest_run_days`, `sebilah_soul_level`, and `bosses_killed`.

#### Scenario: Valid leaderboard query
- **WHEN** client requests a leaderboard with a supported `type`
- **THEN** system returns ranked entries matching that metric

#### Scenario: Unsupported metric
- **WHEN** client requests a leaderboard for an unknown type
- **THEN** system returns a `validation_error` (422)

### Requirement: Personal Best Submission
The system SHALL only update a player's leaderboard entry if the submitted value exceeds their current personal best.

#### Scenario: New personal best
- **WHEN** client submits a value higher than the stored best for a given metric
- **THEN** system updates the database and returns the new `global_rank`

#### Scenario: Value below personal best
- **WHEN** client submits a value lower than or equal to the stored best
- **THEN** system returns a `value_too_low` error (400)

### Requirement: Character-Based Scoping
The system SHALL allow players to filter leaderboards globally or by character (RIMBA, DARA, BAYU, SARI).

#### Scenario: Per-character filter
- **WHEN** client requests the `longest_run_days` for the character `RIMBA`
- **THEN** system returns ranked entries containing only `RIMBA` player runs
