## 1. Local Build & Dependencies

- [x] 1.1 Run `go mod tidy` to clean up dependencies
- [x] 1.2 Compile the project: `go build -o server ./cmd/server`
- [x] 1.3 Verify the binary starts without immediate panics

## 2. Docker Verification

- [x] 2.1 Build Docker image: `docker build -t gamefeel-backend .`
- [x] 2.2 Verify image is listed in local registry
- [x] 2.3 (Optional) Run container locally to ensure runtime stability

## 3. API Integration Testing

- [x] 3.1 Start the local server
- [x] 3.2 Run the integration test suite: `go run scripts/test_hollow_wilds.go`
- [x] 3.3 Confirm all tests pass against the new endpoints
