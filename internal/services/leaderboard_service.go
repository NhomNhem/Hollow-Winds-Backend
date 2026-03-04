package services

import (
	"context"
	"fmt"

	"github.com/NhomNhem/GameFeel-Backend/internal/database"
	"github.com/NhomNhem/GameFeel-Backend/internal/models"
	"github.com/google/uuid"
)

// LeaderboardService handles leaderboard business logic
type LeaderboardService struct{}

// NewLeaderboardService creates a new leaderboard service
func NewLeaderboardService() *LeaderboardService {
	return &LeaderboardService{}
}

// GetGlobalLeaderboard returns top players by total stars
func (s *LeaderboardService) GetGlobalLeaderboard(ctx context.Context, page, perPage int) (*models.GlobalLeaderboardResponse, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	offset := (page - 1) * perPage

	query := `
		WITH ranked_players AS (
			SELECT 
				u.playfab_id,
				u.display_name,
				u.total_stars_collected,
				u.max_map_unlocked,
				COUNT(DISTINCT lc.level_id) as levels_completed,
				RANK() OVER (ORDER BY u.total_stars_collected DESC) as rank
			FROM users u
			LEFT JOIN level_completions lc ON u.id = lc.user_id
			WHERE u.is_banned = false AND u.deleted_at IS NULL
			GROUP BY u.id
		)
		SELECT * FROM ranked_players
		ORDER BY rank ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := database.Pool.Query(ctx, query, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query global leaderboard: %w", err)
	}
	defer rows.Close()

	var leaderboard []models.LeaderboardEntry
	for rows.Next() {
		var entry models.LeaderboardEntry
		err := rows.Scan(
			&entry.PlayerID,
			&entry.DisplayName,
			&entry.TotalStars,
			&entry.MaxMapUnlocked,
			&entry.LevelsCompleted,
			&entry.Rank,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan leaderboard entry: %w", err)
		}
		leaderboard = append(leaderboard, entry)
	}

	// Get total count
	var total int
	err = database.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM users 
		WHERE is_banned = false AND deleted_at IS NULL
	`).Scan(&total)
	if err != nil {
		total = 0
	}

	return &models.GlobalLeaderboardResponse{
		Leaderboard: leaderboard,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
	}, nil
}

// GetLevelLeaderboard returns best times for a specific level
func (s *LeaderboardService) GetLevelLeaderboard(ctx context.Context, levelID, mapID string, limit int) (*models.LevelLeaderboardResponse, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	query := `
		SELECT 
			u.playfab_id,
			u.display_name,
			lc.best_time_seconds,
			lc.stars_earned,
			lc.play_count,
			lc.first_completed_at,
			RANK() OVER (ORDER BY lc.best_time_seconds ASC) as rank
		FROM level_completions lc
		JOIN users u ON lc.user_id = u.id
		WHERE lc.level_id = $1 
		  AND ($2 = '' OR lc.map_id = $2)
		  AND u.is_banned = false
		ORDER BY lc.best_time_seconds ASC
		LIMIT $3
	`

	rows, err := database.Pool.Query(ctx, query, levelID, mapID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query level leaderboard: %w", err)
	}
	defer rows.Close()

	var leaderboard []models.LeaderboardEntry
	for rows.Next() {
		var entry models.LeaderboardEntry
		err := rows.Scan(
			&entry.PlayerID,
			&entry.DisplayName,
			&entry.BestTime,
			&entry.Stars,
			&entry.PlayCount,
			&entry.FirstCompleted,
			&entry.Rank,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan level leaderboard entry: %w", err)
		}
		leaderboard = append(leaderboard, entry)
	}

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(DISTINCT lc.user_id)
		FROM level_completions lc
		JOIN users u ON lc.user_id = u.id
		WHERE lc.level_id = $1 
		  AND ($2 = '' OR lc.map_id = $2)
		  AND u.is_banned = false
	`
	err = database.Pool.QueryRow(ctx, countQuery, levelID, mapID).Scan(&total)
	if err != nil {
		total = 0
	}

	return &models.LevelLeaderboardResponse{
		LevelID:     levelID,
		MapID:       mapID,
		Leaderboard: leaderboard,
		Total:       total,
	}, nil
}

// GetPlayerStats returns authenticated player's position and stats
func (s *LeaderboardService) GetPlayerStats(ctx context.Context, userID uuid.UUID) (*models.PlayerStatsResponse, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	query := `
		WITH ranked_players AS (
			SELECT 
				u.id,
				u.playfab_id,
				u.display_name,
				u.total_stars_collected,
				u.max_map_unlocked,
				COUNT(DISTINCT lc.level_id) as levels_completed,
				COALESCE(AVG(lc.stars_earned), 0) as avg_stars,
				RANK() OVER (ORDER BY u.total_stars_collected DESC) as global_rank
			FROM users u
			LEFT JOIN level_completions lc ON u.id = lc.user_id
			WHERE u.is_banned = false AND u.deleted_at IS NULL
			GROUP BY u.id
		)
		SELECT 
			playfab_id,
			display_name,
			global_rank,
			total_stars_collected,
			max_map_unlocked,
			levels_completed,
			avg_stars
		FROM ranked_players 
		WHERE id = $1
	`

	var stats models.PlayerStatsResponse
	err := database.Pool.QueryRow(ctx, query, userID).Scan(
		&stats.PlayerID,
		&stats.DisplayName,
		&stats.GlobalRank,
		&stats.TotalStars,
		&stats.MaxMapUnlocked,
		&stats.LevelsCompleted,
		&stats.AverageStars,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get player stats: %w", err)
	}

	return &stats, nil
}

// GetLevelStats returns aggregated analytics for a level
func (s *LeaderboardService) GetLevelStats(ctx context.Context, levelID, mapID string) (*models.LevelStatsResponse, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	query := `
		SELECT 
			lc.level_id,
			COALESCE(lc.map_id, ''),
			COUNT(DISTINCT lc.user_id) as unique_players,
			COALESCE(AVG(lc.best_time_seconds), 0) as avg_time,
			COALESCE(MIN(lc.best_time_seconds), 0) as best_time,
			COALESCE(AVG(lc.stars_earned), 0) as avg_stars,
			COALESCE(SUM(lc.play_count), 0) as total_plays
		FROM level_completions lc
		WHERE lc.level_id = $1 AND ($2 = '' OR lc.map_id = $2)
		GROUP BY lc.level_id, lc.map_id
	`

	var stats models.LevelStatsResponse
	err := database.Pool.QueryRow(ctx, query, levelID, mapID).Scan(
		&stats.LevelID,
		&stats.MapID,
		&stats.UniquePlayers,
		&stats.AverageTime,
		&stats.BestTime,
		&stats.AverageStars,
		&stats.TotalPlays,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get level stats: %w", err)
	}

	return &stats, nil
}
