package repository

import (
	"context"

	"github.com/NhomNhem/NhemDangFugBixs-Core/internal/domain/models"
	"github.com/google/uuid"
)

// UserRepository defines the interface for global user data access
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByPlayFabID(ctx context.Context, playfabID string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	UpdateSystemRole(ctx context.Context, id uuid.UUID, role models.SystemRole) error
}
