package repository

import (
	"context"

	"github.com/NhomNhem/GameFeel-Backend/internal/domain/models"
)

// AnalyticsRepository defines the interface for analytics data access
type AnalyticsRepository interface {
	RecordEvents(ctx context.Context, events []models.AnalyticsEvent) (int, int, error)
}
