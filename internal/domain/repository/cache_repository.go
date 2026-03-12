package repository

import (
	"context"
	"time"
)

// CacheRepository defines the interface for distributed caching
type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
