## Why

The current backend implementation has tightly coupled layers where services directly access global variables (like `database.Pool` and `utils.RedisClient`) and handlers are directly instantiated. This makes testing difficult, violates SOLID principles, and hinders long-term maintainability as the project scales.

## What Changes

- **Layer Separation**: Restructure the project into `domain` (entities/interfaces), `usecase` (logic), `infrastructure` (data access), and `delivery` (API/Fiber).
- **Dependency Injection**: Implement constructor-based DI to pass dependencies (repositories, services) into their consumers.
- **Interface-Driven Development**: Define explicit interfaces for repositories and usecases to decouple business logic from implementation details.
- **SOLID Compliance**: Ensure Single Responsibility for each component and Dependency Inversion where high-level logic depends on abstractions.

## Capabilities

### New Capabilities
- `clean-architecture-framework`: The structural foundation for the new architecture.
- `dependency-injection`: The mechanism for managing component lifecycles and providing dependencies.

### Modified Capabilities
- `auth`: Refactor logic to use repositories and usecases.
- `save-system`: Refactor logic to use repositories and usecases.
- `leaderboard`: Refactor logic to use repositories and usecases.
- `analytics`: Refactor logic to use repositories and usecases.

## Impact

- **Codebase**: Major reorganization of `internal/` directory.
- **Testing**: Enable unit testing with mocks for all business logic.
- **Maintainability**: Clearer boundaries and easier navigation of the system.
