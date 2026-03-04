package models

import "time"

// LeaderboardEntry represents a single leaderboard entry
type LeaderboardEntry struct {
	Rank            int     `json:"rank"`
	PlayerID        string  `json:"playerId"`
	DisplayName     *string `json:"displayName,omitempty"`
	TotalStars      int     `json:"totalStars,omitempty"`
	BestTime        float64 `json:"bestTime,omitempty"`
	Stars           int     `json:"stars,omitempty"`
	PlayCount       int     `json:"playCount,omitempty"`
	LevelsCompleted int     `json:"levelsCompleted,omitempty"`
	MaxMapUnlocked  int     `json:"maxMapUnlocked,omitempty"`
	FirstCompleted  *time.Time `json:"firstCompleted,omitempty"`
}

// GlobalLeaderboardResponse for global rankings
type GlobalLeaderboardResponse struct {
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
	Total       int                `json:"total"`
	Page        int                `json:"page"`
	PerPage     int                `json:"perPage"`
}

// LevelLeaderboardResponse for per-level rankings
type LevelLeaderboardResponse struct {
	LevelID     string             `json:"levelId"`
	MapID       string             `json:"mapId,omitempty"`
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
	Total       int                `json:"total"`
}

// PlayerStatsResponse for player position & stats
type PlayerStatsResponse struct {
	PlayerID        string  `json:"playerId"`
	DisplayName     *string `json:"displayName,omitempty"`
	GlobalRank      int     `json:"globalRank"`
	TotalStars      int     `json:"totalStars"`
	MaxMapUnlocked  int     `json:"maxMapUnlocked"`
	LevelsCompleted int     `json:"levelsCompleted"`
	AverageStars    float64 `json:"averageStars"`
}

// LevelStatsResponse for level analytics
type LevelStatsResponse struct {
	LevelID         string  `json:"levelId"`
	MapID           string  `json:"mapId,omitempty"`
	UniquePlayers   int     `json:"uniquePlayers"`
	AverageTime     float64 `json:"averageTime"`
	BestTime        float64 `json:"bestTime"`
	AverageStars    float64 `json:"averageStars"`
	TotalPlays      int     `json:"totalPlays"`
	CompletionRate  float64 `json:"completionRate,omitempty"`
}
