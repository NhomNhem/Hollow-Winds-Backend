package usecase

import (
	"context"

	"github.com/NhomNhem/NhemDangFugBixs-Core/internal/domain/models"
	"github.com/google/uuid"
)

// AnalyticsUsecase defines the interface for analytics-related business logic
type AnalyticsUsecase interface {
	TrackEvents(ctx context.Context, playerID *uuid.UUID, events []models.AnalyticsEvent) (int, int, error)
}
