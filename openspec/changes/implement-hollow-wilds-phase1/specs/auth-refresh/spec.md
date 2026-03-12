## ADDED Requirements

### Requirement: Refresh JWT Token
The system SHALL issue a new JWT access token when provided with a valid refresh token.

#### Scenario: Valid refresh token
- **WHEN** client submits a valid, non-expired refresh token
- **THEN** system returns a new JWT access token with 3600-second expiry

#### Scenario: Expired refresh token
- **WHEN** client submits an expired refresh token
- **THEN** system returns 401 error with code "invalid_refresh_token"

#### Scenario: Already revoked refresh token
- **WHEN** client submits a refresh token that was revoked via logout
- **THEN** system returns 401 error with code "invalid_refresh_token"

### Requirement: Issue Refresh Token on Login
The system SHALL issue a refresh token alongside the JWT access token during authentication.

#### Scenario: Successful PlayFab authentication
- **WHEN** client provides valid PlayFab session ticket
- **THEN** system returns JWT access token AND refresh token

#### Scenario: Refresh token storage
- **WHEN** refresh token is issued
- **THEN** system stores token in Redis with 7-day TTL

### Requirement: Logout and Token Revocation
The system SHALL invalidate the refresh token when client requests logout.

#### Scenario: Successful logout
- **WHEN** authenticated client requests logout
- **THEN** system invalidates refresh token and returns success

#### Scenario: Logout with invalid token
- **WHEN** client requests logout with already-invalid token
- **THEN** system returns success (idempotent operation)
