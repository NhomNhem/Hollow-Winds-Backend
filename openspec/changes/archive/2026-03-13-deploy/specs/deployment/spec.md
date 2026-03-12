## ADDED Requirements

### Requirement: Fly.io Deployment
The system SHALL be deployable to Fly.io using the `fly deploy` command.

#### Scenario: Successful production push
- **WHEN** the developer executes `fly deploy` from the project root
- **THEN** Fly.io builds the image, updates the application, and starts the new instance successfully

### Requirement: Environment Configuration
The production environment SHALL have all required secrets and variables configured via Fly.io.

#### Scenario: Verify secrets
- **WHEN** the application starts in production
- **THEN** it successfully retrieves `DATABASE_URL`, `REDIS_URL`, `JWT_SECRET`, and `PLAYFAB_TITLE_ID` from the environment
