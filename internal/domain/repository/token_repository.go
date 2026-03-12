package repository

import (
	"context"
	"time"
)

// TokenRepository defines the interface for token management (refresh, blacklist)
type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, token string, userID string, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, token string) (string, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	BlacklistJWT(ctx context.Context, jti string, ttl time.Duration) error
	IsJWTBlacklisted(ctx context.Context, jti string) (bool, error)
}
