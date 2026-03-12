## 1. Project Reorganization

- [x] 1.1 Create new directory structure: `internal/{domain,usecase,infrastructure,delivery}`
- [x] 1.2 Move existing models to `internal/domain/models` (or similar)

## 2. Infrastructure Layer (Adapters)

- [x] 2.1 Implement `PostgresSaveRepository` in `internal/infrastructure/persistence`
- [x] 2.2 Implement `RedisCacheRepository` in `internal/infrastructure/cache`
- [x] 2.3 Implement `PlayFabIdentityRepository` in `internal/infrastructure/identity`

## 3. Usecase Layer (Business Logic)

- [x] 3.1 Define `AuthUsecase` interface and implement it in `internal/usecase/auth`
- [x] 3.2 Define `PlayerUsecase` interface and implement it in `internal/usecase/player`
- [x] 3.3 Define `LeaderboardUsecase` interface and implement it in `internal/usecase/leaderboard`

## 4. Delivery Layer (Handlers)

- [x] 4.1 Refactor `HollowWildsHandler` to receive usecases via constructor
- [x] 4.2 Refactor `LeaderboardHandler` to receive usecases via constructor

## 5. Main & DI

- [x] 5.1 Update `cmd/server/main.go` to initialize repositories and inject them into usecases
- [x] 5.2 Inject usecases into handlers during registration
- [x] 5.3 Verify all existing integration tests pass: `go run scripts/test_hollow_wilds.go`
