## MODIFIED Requirements

### Requirement: PlayFab Session Ticket Validation
The system SHALL validate a PlayFab session ticket provided by the client against the PlayFab API via a dedicated `IdentityRepository`.

#### Scenario: Valid ticket
- **WHEN** client provides a valid PlayFab session ticket and PlayFab ID
- **THEN** system successfully validates the ticket via the repository and retrieves player account info

#### Scenario: Invalid ticket
- **WHEN** client provides an invalid or expired PlayFab session ticket
- **THEN** system returns an `unauthorized` (401) error

### Requirement: Refresh Token Management
The system SHALL issue and manage refresh tokens via a `TokenRepository` to allow players to obtain new access tokens.

#### Scenario: Refresh token usage
- **WHEN** client provides a valid refresh token to the refresh endpoint
- **THEN** system issues a new JWT access token after verifying with the repository

#### Scenario: Logout
- **WHEN** client calls the logout endpoint with a refresh token
- **THEN** system revokes the refresh token in the repository and blacklists the current JWT
