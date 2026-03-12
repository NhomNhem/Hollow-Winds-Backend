## 1. Database Schema

- [ ] 1.1 Create `leaderboard_entries` table migration (id, user_id, level_id, best_time_seconds, completed_at, created_at, updated_at)
- [ ] 1.2 Create indexes on leaderboard_entries (user_id, level_id, best_time_seconds composite index)
- [ ] 1.3 Create `leaderboard_cache` table for Redis sync tracking (optional)
- [ ] 1.4 Run migrations and verify schema in Supabase

## 2. Models

- [ ] 2.1 Create `LeaderboardEntry` model in internal/models/leaderboard.go
- [ ] 2.2 Create `LeaderboardService` interface definition
- [ ] 2.3 Create request/response DTOs for leaderboard API endpoints

## 3. Service Layer

- [ ] 3.1 Implement `GetGlobalLeaderboard` function with pagination support
- [ ] 3.2 Implement `GetPlayerRank` function returning player rank and surrounding players
- [ ] 3.3 Implement `UpdateLeaderboardEntry` function for new completions
- [ ] 3.4 Implement `GetFriendsLeaderboard` function with social graph lookup
- [ ] 3.5 Add Redis caching layer for leaderboard reads (30-second TTL)
- [ ] 3.6 Add cache invalidation logic on leaderboard updates

## 4. API Handlers

- [ ] 4.1 Create `GET /api/v1/leaderboards/{levelId}` endpoint for global rankings
- [ ] 4.2 Create `GET /api/v1/leaderboards/{levelId}/me` endpoint for player's rank
- [ ] 4.3 Create `GET /api/v1/leaderboards/{levelId}/friends` endpoint for friends rankings
- [ ] 4.4 Add query parameter support for time period filtering (daily, weekly, all-time)
- [ ] 4.5 Add pagination query parameters (page, perPage)
- [ ] 4.6 Register new routes in internal/api/routes.go

## 5. Integration with Level Completion

- [ ] 5.1 Modify `POST /api/v1/levels/complete` handler to call leaderboard service
- [ ] 5.2 Add leaderboard update logic after successful level validation
- [ ] 5.3 Add error handling for leaderboard failures (should not block level completion)

## 6. Admin Endpoints

- [ ] 6.1 Create `DELETE /api/v1/admin/leaderboards/{levelId}` endpoint for resetting leaderboards
- [ ] 6.2 Create `GET /api/v1/admin/leaderboards/stats` endpoint for leaderboard statistics
- [ ] 6.3 Add admin middleware check to admin endpoints

## 7. Testing

- [ ] 7.1 Write unit tests for leaderboard service functions
- [ ] 7.2 Write integration tests for leaderboard API endpoints
- [ ] 7.3 Write tests for cache invalidation logic
- [ ] 7.4 Add load testing for leaderboard read operations
- [ ] 7.5 Test edge cases (empty leaderboards, ties, pagination boundaries)

## 8. Documentation

- [ ] 8.1 Add Swagger annotations for new leaderboard endpoints
- [ ] 8.2 Update docs/api.md with leaderboard API documentation
- [ ] 8.3 Run `swag init` to regenerate swagger docs
- [ ] 8.4 Add example requests/responses to documentation
