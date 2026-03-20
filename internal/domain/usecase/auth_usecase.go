package usecase

import (
	"context"

	"github.com/NhomNhem/NhemDangFugBixs-Core/internal/domain/models"
)

// AuthUsecase defines the interface for authentication business logic
type AuthUsecase interface {
	// Login performs authentication via PlayFab and returns platform tokens
	Login(ctx context.Context, sessionTicket string, overridePlayFabID string) (*models.AuthResponse, error)

	// RefreshToken updates the access token using a valid refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokenResponse, error)

	// Logout invalidates the refresh token and blacklists the current JWT
	Logout(ctx context.Context, refreshToken string, jti string) error

	// LegacyLogin provides backward compatibility for Phase 1 games
	LegacyLogin(ctx context.Context, playfabID, displayName, sessionToken string) (*models.AuthResponse, error)
}
