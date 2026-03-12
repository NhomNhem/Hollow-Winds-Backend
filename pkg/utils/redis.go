package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// InitRedis initializes the Redis client
func InitRedis() error {
	host := os.Getenv("UPSTASH_REDIS_HOST")
	port := os.Getenv("UPSTASH_REDIS_PORT")
	password := os.Getenv("UPSTASH_REDIS_PASSWORD")
	if password == "" {
		password = os.Getenv("REDIS_PASSWORD")
	}

	db := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		fmt.Sscanf(dbStr, "%d", &db)
	}

	// If not configured, skip initialization (development mode)
	if host == "" {
		return nil
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// CacheSaveData caches player save data with TTL
func CacheSaveData(ctx context.Context, playerID string, saveData string, ttl time.Duration) error {
	if RedisClient == nil {
		return nil // Skip if Redis not configured
	}

	key := fmt.Sprintf("player:save:%s", playerID)
	return RedisClient.Set(ctx, key, saveData, ttl).Err()
}

// GetCachedSaveData retrieves cached player save data
func GetCachedSaveData(ctx context.Context, playerID string) (string, error) {
	if RedisClient == nil {
		return "", nil // Skip if Redis not configured
	}

	key := fmt.Sprintf("player:save:%s", playerID)
	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Cache miss
	}
	return val, err
}

// InvalidateSaveCache invalidates cached save data for a player
func InvalidateSaveCache(ctx context.Context, playerID string) error {
	if RedisClient == nil {
		return nil // Skip if Redis not configured
	}

	key := fmt.Sprintf("player:save:%s", playerID)
	return RedisClient.Del(ctx, key).Err()
}

// StoreRefreshToken stores a refresh token in Redis with TTL
func StoreRefreshToken(ctx context.Context, token string, playerID string, ttl time.Duration) error {
	if RedisClient == nil {
		return nil // Skip if Redis not configured
	}

	key := fmt.Sprintf("session:%s", token)
	return RedisClient.Set(ctx, key, playerID, ttl).Err()
}

// GetRefreshToken retrieves the player ID for a refresh token
func GetRefreshToken(ctx context.Context, token string) (string, error) {
	if RedisClient == nil {
		return "", nil // Skip if Redis not configured
	}

	key := fmt.Sprintf("session:%s", token)
	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Token not found
	}
	return val, err
}

// DeleteRefreshToken removes a refresh token from Redis
func DeleteRefreshToken(ctx context.Context, token string) error {
	if RedisClient == nil {
		return nil // Skip if Redis not configured
	}

	key := fmt.Sprintf("session:%s", token)
	return RedisClient.Del(ctx, key).Err()
}

// BlacklistJWT adds a JWT JTI to the blacklist with TTL
func BlacklistJWT(ctx context.Context, jti string, ttl time.Duration) error {
	if RedisClient == nil {
		return nil // Skip if Redis not configured
	}

	key := fmt.Sprintf("session:%s:blacklist", jti)
	return RedisClient.Set(ctx, key, "1", ttl).Err()
}

// IsJWTBlacklisted checks if a JWT JTI is blacklisted
func IsJWTBlacklisted(ctx context.Context, jti string) (bool, error) {
	if RedisClient == nil {
		return false, nil // Skip if Redis not configured
	}

	key := fmt.Sprintf("session:%s:blacklist", jti)
	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil // Not blacklisted
	}
	return val == "1", err
}
