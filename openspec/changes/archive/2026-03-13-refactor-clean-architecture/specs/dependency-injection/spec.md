## ADDED Requirements

### Requirement: Constructor Injection
All components SHALL receive their dependencies via constructors rather than accessing global variables.

#### Scenario: Verify handler instantiation
- **WHEN** a new API handler is created in `main.go`
- **THEN** it receives its corresponding usecase as an interface via its constructor

### Requirement: Dependency Inversion
Components SHALL depend on interfaces defined in the Domain layer rather than concrete implementations.

#### Scenario: Verify repository usage
- **WHEN** a usecase needs data access
- **THEN** it interacts with a repository interface from the Domain layer
