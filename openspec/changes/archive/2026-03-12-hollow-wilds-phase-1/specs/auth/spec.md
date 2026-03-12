## ADDED Requirements

### Requirement: PlayFab Session Ticket Validation
The system SHALL validate a PlayFab session ticket provided by the client against the PlayFab API.

#### Scenario: Valid ticket
- **WHEN** client provides a valid PlayFab session ticket and PlayFab ID
- **THEN** system successfully validates the ticket and retrieves player account info

#### Scenario: Invalid ticket
- **WHEN** client provides an invalid or expired PlayFab session ticket
- **THEN** system returns an `unauthorized` (401) error

### Requirement: JWT Issuance
The system SHALL issue a JWT access token upon successful PlayFab validation, containing the internal player UUID and PlayFab ID.

#### Scenario: Token generation
- **WHEN** player is authenticated
- **THEN** system returns a JWT token signed with `JWT_SECRET` that expires in 1 hour

### Requirement: Refresh Token Management
The system SHALL issue and manage refresh tokens to allow players to obtain new access tokens without re-authenticating with PlayFab.

#### Scenario: Refresh token usage
- **WHEN** client provides a valid refresh token to the refresh endpoint
- **THEN** system issues a new JWT access token

#### Scenario: Logout
- **WHEN** client calls the logout endpoint with a refresh token
- **THEN** system revokes the refresh token and blacklists the current JWT
