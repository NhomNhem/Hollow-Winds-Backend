package usecase

import (
	"context"

	"github.com/NhomNhem/NhemDangFugBixs-Core/internal/domain/models"
	"github.com/google/uuid"
)

// LevelUsecase defines the business logic for level progression
type LevelUsecase interface {
	CompleteLevel(ctx context.Context, userID uuid.UUID, req *models.LevelCompletionRequest) (*models.LevelCompletionResponse, error)
}
