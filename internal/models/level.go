package models

import (
	"time"

	"github.com/google/uuid"
)

// LevelCompletionRequest represents the request to complete a level
type LevelCompletionRequest struct {
	LevelID      string  `json:"levelId" validate:"required"`
	MapID        string  `json:"mapId" validate:"required"`
	TimeSeconds  float64 `json:"timeSeconds" validate:"required,gt=0"`
	FinalHP      float64 `json:"finalHp" validate:"required,gte=0"`
	DashCount    int     `json:"dashCount" validate:"gte=0"`
	CounterCount int     `json:"counterCount" validate:"gte=0"`
	VulnerableKills int  `json:"vulnerableKills" validate:"gte=0"`
}

// LevelCompletionResponse represents the response after completing a level
type LevelCompletionResponse struct {
	Success          bool     `json:"success"`
	StarsEarned      int      `json:"starsEarned"`
	GoldEarned       int      `json:"goldEarned"`
	NewTotalGold     int      `json:"newTotalGold"`
	NewTotalStars    int      `json:"newTotalStars"`
	NextLevelUnlocked *string `json:"nextLevelUnlocked,omitempty"`
	MapUnlocked      *string  `json:"mapUnlocked,omitempty"`
	IsFirstCompletion bool    `json:"isFirstCompletion"`
	NewBestTime      bool     `json:"newBestTime"`
}

// LevelCompletion represents a level completion record
type LevelCompletion struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"userId"`
	LevelID     string     `json:"levelId"`
	MapID       string     `json:"mapId"`
	
	// Performance stats
	StarsEarned       int     `json:"starsEarned"`
	BestTimeSeconds   float64 `json:"bestTimeSeconds"`
	PlayCount         int     `json:"playCount"`
	
	// Latest run stats
	LastFinalHP       *float64 `json:"lastFinalHp,omitempty"`
	LastDashCount     *int     `json:"lastDashCount,omitempty"`
	LastCounterCount  *int     `json:"lastCounterCount,omitempty"`
	LastVulnerableKills *int   `json:"lastVulnerableKills,omitempty"`
	
	// Timestamps
	FirstCompletedAt  time.Time  `json:"firstCompletedAt"`
	LastPlayedAt      time.Time  `json:"lastPlayedAt"`
}

// LevelObjective represents a single objective for a level
type LevelObjective struct {
	Type      string  `json:"type"`      // "health", "time", "dash_count", etc.
	Threshold float64 `json:"threshold"` // Required value
	Operator  string  `json:"operator"`  // "gt", "gte", "lt", "lte", "eq"
}

// LevelConfig represents level configuration for validation
type LevelConfig struct {
	LevelID      string           `json:"levelId"`
	MapID        string           `json:"mapId"`
	MinTimeSeconds float64        `json:"minTimeSeconds"` // Minimum possible completion time (anti-cheat)
	BaseGold     int              `json:"baseGold"`       // Base gold reward
	Objectives   []LevelObjective `json:"objectives"`     // 3 objectives for 3 stars
}

// Anti-cheat detection levels
const (
	AntiCheatNone     = 0 // No issues
	AntiCheatSuspicious = 1 // Suspicious but allow
	AntiCheatModerate = 2 // Flag for review
	AntiCheatSevere   = 3 // Reject and potentially ban
)

// AntiCheatResult represents anti-cheat validation result
type AntiCheatResult struct {
	Level   int      `json:"level"`
	Reasons []string `json:"reasons,omitempty"`
}
