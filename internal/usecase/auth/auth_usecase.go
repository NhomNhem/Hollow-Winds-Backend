package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/NhomNhem/NhemDangFugBixs-Core/internal/domain/models"
	"github.com/NhomNhem/NhemDangFugBixs-Core/internal/domain/repository"
	"github.com/NhomNhem/NhemDangFugBixs-Core/internal/domain/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authUsecase struct {
	userRepo     repository.UserRepository
	playerRepo   repository.PlayerRepository
	identityRepo repository.IdentityRepository
	tokenRepo    repository.TokenRepository
}

// NewAuthUsecase creates a new authentication usecase
func NewAuthUsecase(
	userRepo repository.UserRepository,
	playerRepo repository.PlayerRepository,
	identityRepo repository.IdentityRepository,
	tokenRepo repository.TokenRepository,
) usecase.AuthUsecase {
	return &authUsecase{
		userRepo:     userRepo,
		playerRepo:   playerRepo,
		identityRepo: identityRepo,
		tokenRepo:    tokenRepo,
	}
}

func (u *authUsecase) Login(ctx context.Context, sessionTicket string, overridePlayFabID string) (*models.AuthResponse, error) {
	// 1. Validate PlayFab ticket
	playfabID, err := u.identityRepo.ValidateTicket(ctx, sessionTicket)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// 2. Allow override in dev mode
	if playfabID == "MOCK_PLAYFAB_ID" && overridePlayFabID != "" {
		playfabID = overridePlayFabID
	}

	// 3. Get or create platform user
	user, err := u.userRepo.GetByPlayFabID(ctx, playfabID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user = &models.User{
			ID:          uuid.New(),
			PlayFabID:   playfabID,
			SystemRole:  models.RoleUser,
			CreatedAt:   time.Now(),
			LastLoginAt: time.Now(),
		}
		if err := u.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	} else {
		// Update last login
		_ = u.userRepo.UpdateLastLogin(ctx, user.ID)
	}

	// 4. Ensure game-specific player profile exists (for Hollow Wilds)
	// In the future, this should be game-agnostic or moved to a dynamic handler
	player, _ := u.playerRepo.GetByPlayFabID(ctx, playfabID)
	if player == nil {
		player = &models.Player{
			ID:         user.ID, // Link by ID for consistency
			PlayFabID:  playfabID,
			CreatedAt:  time.Now(),
			LastSeenAt: time.Now(),
		}
		_ = u.playerRepo.Create(ctx, player)
	} else {
		_ = u.playerRepo.UpdateLastSeen(ctx, player.ID)
	}

	// 5. Generate JWT with proper role
	token, expiresIn, err := u.generateJWT(user.ID.String(), user.PlayFabID, user.SystemRole)
	if err != nil {
		return nil, err
	}

	// 6. Generate Refresh Token (30 days TTL per Master Plan)
	refreshToken := uuid.New().String()
	_ = u.tokenRepo.StoreRefreshToken(ctx, refreshToken, user.ID.String(), 30*24*time.Hour)

	return &models.AuthResponse{
		JWT:          token,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User:         *user,
	}, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokenResponse, error) {
	userIDStr, err := u.tokenRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil || userIDStr == "" {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in token")
	}

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	token, expiresIn, err := u.generateJWT(user.ID.String(), user.PlayFabID, user.SystemRole)
	if err != nil {
		return nil, err
	}

	return &models.RefreshTokenResponse{
		Token:     token,
		ExpiresIn: expiresIn,
	}, nil
}

func (u *authUsecase) Logout(ctx context.Context, refreshToken string, jti string) error {
	if refreshToken != "" {
		_ = u.tokenRepo.DeleteRefreshToken(ctx, refreshToken)
	}
	if jti != "" {
		_ = u.tokenRepo.BlacklistJWT(ctx, jti, 24*time.Hour)
	}
	return nil
}

func (u *authUsecase) LegacyLogin(ctx context.Context, playfabID, displayName, sessionToken string) (*models.AuthResponse, error) {
	// Re-use core Login logic but with display name hint
	user, err := u.userRepo.GetByPlayFabID(ctx, playfabID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user = &models.User{
			ID:          uuid.New(),
			PlayFabID:   playfabID,
			DisplayName: &displayName,
			SystemRole:  models.RoleUser,
			CreatedAt:   time.Now(),
			LastLoginAt: time.Now(),
		}
		if err := u.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	token, expiresIn, err := u.generateJWT(user.ID.String(), user.PlayFabID, user.SystemRole)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		JWT:       token,
		ExpiresIn: expiresIn,
		User:      *user,
	}, nil
}

func (u *authUsecase) generateJWT(userID, playfabID string, role models.SystemRole) (string, int, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-123"
	}

	// 60 minutes TTL per Master Plan
	expiresIn := 3600
	now := time.Now()
	expiresAt := now.Add(time.Duration(expiresIn) * time.Second)

	claims := models.JWTClaims{
		UserID:    userID,
		PlayFabID: playfabID,
		Role:      role,
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    claims.UserID,
		"sub":       claims.UserID,
		"playfabId": claims.PlayFabID,
		"role":      string(claims.Role),
		"iat":       claims.IssuedAt,
		"exp":       claims.ExpiresAt,
		"jti":       uuid.New().String(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresIn, nil
}
