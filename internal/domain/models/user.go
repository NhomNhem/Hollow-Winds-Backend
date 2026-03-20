package models

import (
	"time"

	"github.com/google/uuid"
)

// SystemRole represents a user's global platform-level role
type SystemRole string

const (
	RoleSuperAdmin  SystemRole = "super_admin"
	RoleAdmin       SystemRole = "admin"
	RoleGameManager SystemRole = "game_manager"
	RoleSupport     SystemRole = "support"
	RoleUser        SystemRole = "user"
	RoleBetaTester  SystemRole = "beta_tester"
)

// User represents a user in the system
type User struct {
	ID          uuid.UUID  `json:"id"`
	PlayFabID   string     `json:"playfabId"`
	DisplayName *string    `json:"displayName,omitempty"`
	SystemRole  SystemRole `json:"systemRole"`

	// Currency (Legacy - TODO: Move to game-specific storage)
	Gold     int `json:"gold"`
	Diamonds int `json:"diamonds"`

	// Progression (Legacy - TODO: Move to game-specific storage)
	MaxMapUnlocked      int `json:"maxMapUnlocked"`
	TotalStarsCollected int `json:"totalStarsCollected"`

	// Metadata
	CreatedAt            time.Time  `json:"createdAt"`
	LastLoginAt          time.Time  `json:"lastLoginAt"`
	LastPlayedAt         *time.Time `json:"lastPlayedAt,omitempty"`
	TotalPlayTimeSeconds *int       `json:"totalPlayTimeSeconds,omitempty"`

	// Social
	FacebookID *string `json:"facebookId,omitempty"`
	GoogleID   *string `json:"googleId,omitempty"`

	// Flags
	IsBanned  bool       `json:"isBanned"`
	BanReason *string    `json:"banReason,omitempty"`
	BannedAt  *time.Time `json:"bannedAt,omitempty"`

	// GDPR
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

// AuthRequest represents the login request
type AuthRequest struct {
	PlayFabID   string  `json:"playfabId" validate:"required"`
	DisplayName *string `json:"displayName,omitempty"`
}

// AuthResponse represents the login response
type AuthResponse struct {
	JWT          string `json:"jwt"`
	RefreshToken string `json:"refreshToken,omitempty"`
	User         User   `json:"user"`
	ExpiresIn    int    `json:"expiresIn"` // seconds
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID    string     `json:"userId"`
	PlayFabID string     `json:"playfabId"`
	Role      SystemRole `json:"role"`
	IssuedAt  int64      `json:"iat"`
	ExpiresAt int64      `json:"exp"`
}
