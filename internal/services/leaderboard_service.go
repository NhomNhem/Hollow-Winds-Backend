package services

import (
	"context"
	"fmt"
	"time"

	"github.com/NhomNhem/GameFeel-Backend/internal/database"
	"github.com/NhomNhem/GameFeel-Backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

// GetHollowWildsLeaderboard returns the specified leaderboard for Hollow Wilds
func (s *LeaderboardService) GetHollowWildsLeaderboard(ctx context.Context, lbType, scope, character string, limit, offset int) (*models.HollowWildsLeaderboardResponse, error) {
	if database.Pool == nil {
		// Mock leaderboard for development
		return &models.HollowWildsLeaderboardResponse{
			Type:      lbType,
			Scope:     scope,
			Character: character,
			Total:     1,
			Entries: []models.HollowWildsLeaderboardEntry{
				{
					Rank:        1,
					PlayerID:    "MOCK_PLAYER",
					DisplayName: "Mock Player",
					Value:       42,
					Character:   character,
					UpdatedAt:   time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	query := `
		SELECT 
			p.playfab_id,
			COALESCE(p.display_name, 'Anonymous'),
			le.value,
			le.character,
			le.world_seed,
			le.combat_build,
			le.updated_at,
			le.run_metadata,
			RANK() OVER (ORDER BY le.value DESC) as rank
		FROM leaderboard_entries le
		JOIN players p ON le.player_id = p.id
		WHERE le.type = $1
		  AND ($2 = 'global' OR le.character = $3)
		ORDER BY le.value DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := database.Pool.Query(ctx, query, lbType, scope, character, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query leaderboard: %w", err)
	}
	defer rows.Close()

	var entries []models.HollowWildsLeaderboardEntry
	for rows.Next() {
		var entry models.HollowWildsLeaderboardEntry
		var updatedAt time.Time
		err := rows.Scan(
			&entry.PlayerID,
			&entry.DisplayName,
			&entry.Value,
			&entry.Character,
			&entry.WorldSeed,
			&entry.CombatBuild,
			&updatedAt,
			&entry.RunMetadata,
			&entry.Rank,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}
		entry.UpdatedAt = updatedAt.Format("2006-01-02T15:04:05Z")
		entries = append(entries, entry)
	}

	// Get total
	var total int
	countQuery := `
		SELECT COUNT(*) FROM leaderboard_entries 
		WHERE type = $1 AND ($2 = 'global' OR character = $3)
	`
	err = database.Pool.QueryRow(ctx, countQuery, lbType, scope, character).Scan(&total)
	if err != nil {
		total = 0
	}

	return &models.HollowWildsLeaderboardResponse{
		Type:      lbType,
		Scope:     scope,
		Character: character,
		Total:     total,
		Entries:   entries,
	}, nil
}

// SubmitHollowWildsEntry submits a new run result to the leaderboard
func (s *LeaderboardService) SubmitHollowWildsEntry(ctx context.Context, playerID uuid.UUID, req models.LeaderboardSubmitRequest) (*models.LeaderboardSubmitResponse, error) {
	if database.Pool == nil {
		// Mock submission for development
		return &models.LeaderboardSubmitResponse{
			Success:        true,
			GlobalRank:     10,
			CharacterRank:  5,
			IsPersonalBest: true,
		}, nil
	}

	// Get current personal best
	var currentBest int64
	err := database.Pool.QueryRow(ctx, `
		SELECT value FROM leaderboard_entries 
		WHERE player_id = $1 AND type = $2 AND character = $3
	`, playerID, req.Type, req.Character).Scan(&currentBest)

	isPB := false
	if err == pgx.ErrNoRows {
		isPB = true
	} else if err != nil {
		return nil, fmt.Errorf("failed to get personal best: %w", err)
	} else if req.Value > currentBest {
		isPB = true
	}

	if !isPB {
		return nil, fmt.Errorf("value_too_low: Submitted value does not beat personal best")
	}

	// Submit entry (upsert)
	_, err = database.Pool.Exec(ctx, `
		INSERT INTO leaderboard_entries (player_id, type, value, character, world_seed, combat_build, run_metadata, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (player_id, type, character) DO UPDATE
		SET value = EXCLUDED.value,
		    world_seed = EXCLUDED.world_seed,
		    combat_build = EXCLUDED.combat_build,
		    run_metadata = EXCLUDED.run_metadata,
		    updated_at = NOW()
	`, playerID, req.Type, req.Value, req.Character, req.WorldSeed, req.CombatBuild, req.RunMetadata)

	if err != nil {
		return nil, fmt.Errorf("failed to submit entry: %w", err)
	}

	// Calculate ranks (simplified for now, ideally use a more efficient query or cache)
	var globalRank int
	err = database.Pool.QueryRow(ctx, `
		SELECT COUNT(*) + 1 FROM leaderboard_entries 
		WHERE type = $1 AND value > $2
	`, req.Type, req.Value).Scan(&globalRank)

	var characterRank int
	err = database.Pool.QueryRow(ctx, `
		SELECT COUNT(*) + 1 FROM leaderboard_entries 
		WHERE type = $1 AND character = $2 AND value > $3
	`, req.Type, req.Character, req.Value).Scan(&characterRank)

	return &models.LeaderboardSubmitResponse{
		Success:        true,
		GlobalRank:     globalRank,
		CharacterRank:  characterRank,
		IsPersonalBest: true,
	}, nil
}

// GetPlayerHollowWildsStats returns rankings for a player across all types
func (s *LeaderboardService) GetPlayerHollowWildsStats(ctx context.Context, playerID uuid.UUID) (*models.PlayerLeaderboardResponse, error) {
	if database.Pool == nil {
		// Mock stats for development
		return &models.PlayerLeaderboardResponse{
			Entries: []models.PlayerLeaderboardEntry{
				{
					Type:          "longest_run_days",
					GlobalRank:    10,
					CharacterRank: 5,
					Character:     "RIMBA",
					Value:         15,
					PersonalBest:  true,
				},
			},
		}, nil
	}

	query := `
		WITH ranked_global AS (
			SELECT type, player_id, value, RANK() OVER (PARTITION BY type ORDER BY value DESC) as rank
			FROM leaderboard_entries
		),
		ranked_character AS (
			SELECT type, player_id, character, value, RANK() OVER (PARTITION BY type, character ORDER BY value DESC) as rank
			FROM leaderboard_entries
		)
		SELECT 
			le.type,
			rg.rank as global_rank,
			rc.rank as character_rank,
			le.character,
			le.value
		FROM leaderboard_entries le
		JOIN ranked_global rg ON le.player_id = rg.player_id AND le.type = rg.type
		JOIN ranked_character rc ON le.player_id = rc.player_id AND le.type = rc.type AND le.character = rc.character
		WHERE le.player_id = $1
	`

	rows, err := database.Pool.Query(ctx, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query player stats: %w", err)
	}
	defer rows.Close()

	var entries []models.PlayerLeaderboardEntry
	for rows.Next() {
		var entry models.PlayerLeaderboardEntry
		err := rows.Scan(
			&entry.Type,
			&entry.GlobalRank,
			&entry.CharacterRank,
			&entry.Character,
			&entry.Value,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}
		entry.PersonalBest = true
		entries = append(entries, entry)
	}

	return &models.PlayerLeaderboardResponse{
		Entries: entries,
	}, nil
}
