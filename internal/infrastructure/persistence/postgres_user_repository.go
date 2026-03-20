package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/NhomNhem/NhemDangFugBixs-Core/internal/domain/models"
	"github.com/NhomNhem/NhemDangFugBixs-Core/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresUserRepository struct {
	db *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if r.db == nil {
		return &models.User{
			ID:         id,
			PlayFabID:  "MOCK_PLAYFAB_ID",
			SystemRole: models.RoleUser,
			CreatedAt:  time.Now(),
		}, nil
	}

	var user models.User
	err := r.db.QueryRow(ctx, `
		SELECT id, playfab_id, display_name, system_role, gold, diamonds, 
		       max_map_unlocked, total_stars_collected, created_at, last_login_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID, &user.PlayFabID, &user.DisplayName, &user.SystemRole, &user.Gold, &user.Diamonds,
		&user.MaxMapUnlocked, &user.TotalStarsCollected, &user.CreatedAt, &user.LastLoginAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *postgresUserRepository) GetByPlayFabID(ctx context.Context, playfabID string) (*models.User, error) {
	if r.db == nil {
		return nil, nil
	}

	var user models.User
	err := r.db.QueryRow(ctx, `
		SELECT id, playfab_id, display_name, system_role, gold, diamonds, 
		       max_map_unlocked, total_stars_collected, created_at, last_login_at
		FROM users
		WHERE playfab_id = $1
	`, playfabID).Scan(
		&user.ID, &user.PlayFabID, &user.DisplayName, &user.SystemRole, &user.Gold, &user.Diamonds,
		&user.MaxMapUnlocked, &user.TotalStarsCollected, &user.CreatedAt, &user.LastLoginAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by playfab id: %w", err)
	}

	return &user, nil
}

func (r *postgresUserRepository) Create(ctx context.Context, user *models.User) error {
	if r.db == nil {
		return nil
	}
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (playfab_id, display_name, system_role, created_at, last_login_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, user.PlayFabID, user.DisplayName, user.SystemRole, user.CreatedAt, user.LastLoginAt).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *postgresUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	if r.db == nil {
		return nil
	}
	_, err := r.db.Exec(ctx, `UPDATE users SET last_login_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *postgresUserRepository) UpdateSystemRole(ctx context.Context, id uuid.UUID, role models.SystemRole) error {
	if r.db == nil {
		return nil
	}
	_, err := r.db.Exec(ctx, `UPDATE users SET system_role = $1 WHERE id = $2`, role, id)
	return err
}
