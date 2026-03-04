package api

import (
	"log"
	"strconv"

	"github.com/NhomNhem/GameFeel-Backend/internal/models"
	"github.com/NhomNhem/GameFeel-Backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// LeaderboardHandler handles leaderboard endpoints
type LeaderboardHandler struct {
	leaderboardService *services.LeaderboardService
}

// NewLeaderboardHandler creates a new leaderboard handler
func NewLeaderboardHandler() *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboardService: services.NewLeaderboardService(),
	}
}

// GetGlobalLeaderboard handles global leaderboard request
// @Summary Get global leaderboard
// @Description Get top players ranked by total stars collected
// @Tags Leaderboard
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param perPage query int false "Results per page (max 100)" default(100)
// @Success 200 {object} models.APIResponse{data=models.GlobalLeaderboardResponse} "Global leaderboard"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /leaderboard/global [get]
func (h *LeaderboardHandler) GetGlobalLeaderboard(c *fiber.Ctx) error {
	// Parse query parameters
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(c.Query("perPage", "100"))
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 100
	}

	// Get leaderboard
	leaderboard, err := h.leaderboardService.GetGlobalLeaderboard(c.Context(), page, perPage)
	if err != nil {
		log.Printf("Failed to get global leaderboard: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to retrieve global leaderboard",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    leaderboard,
	})
}

// GetLevelLeaderboard handles per-level leaderboard request
// @Summary Get level leaderboard
// @Description Get best times for a specific level
// @Tags Leaderboard
// @Produce json
// @Param levelId path string true "Level ID (e.g., 1-1)"
// @Param mapId query string false "Map ID filter (optional)"
// @Param limit query int false "Result limit (max 100)" default(100)
// @Success 200 {object} models.APIResponse{data=models.LevelLeaderboardResponse} "Level leaderboard"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /leaderboard/level/{levelId} [get]
func (h *LeaderboardHandler) GetLevelLeaderboard(c *fiber.Ctx) error {
	// Parse path and query parameters
	levelID := c.Params("levelId")
	if levelID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "levelId is required",
			},
		})
	}

	mapID := c.Query("mapId", "")
	
	limit, err := strconv.Atoi(c.Query("limit", "100"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 100
	}

	// Get leaderboard
	leaderboard, err := h.leaderboardService.GetLevelLeaderboard(c.Context(), levelID, mapID, limit)
	if err != nil {
		log.Printf("Failed to get level leaderboard for %s: %v", levelID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to retrieve level leaderboard",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    leaderboard,
	})
}

// GetPlayerStats handles player stats request
// @Summary Get player leaderboard stats
// @Description Get authenticated player's global rank and statistics
// @Tags Leaderboard
// @Produce json
// @Param Authorization header string true "Bearer JWT token" default(Bearer )
// @Success 200 {object} models.APIResponse{data=models.PlayerStatsResponse} "Player stats"
// @Failure 401 {object} models.APIResponse{error=models.APIError} "Unauthorized"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /leaderboard/player/me [get]
// @Security BearerAuth
func (h *LeaderboardHandler) GetPlayerStats(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeUnauthorized,
				Message: "User not authenticated",
			},
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid user ID",
			},
		})
	}

	// Get player stats
	stats, err := h.leaderboardService.GetPlayerStats(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get player stats for %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to retrieve player stats",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    stats,
	})
}

// GetLevelStats handles level analytics request
// @Summary Get level statistics
// @Description Get aggregated analytics for a specific level
// @Tags Analytics
// @Produce json
// @Param levelId path string true "Level ID (e.g., 1-1)"
// @Param mapId query string false "Map ID filter (optional)"
// @Success 200 {object} models.APIResponse{data=models.LevelStatsResponse} "Level statistics"
// @Failure 500 {object} models.APIResponse{error=models.APIError} "Internal server error"
// @Router /analytics/level-stats/{levelId} [get]
func (h *LeaderboardHandler) GetLevelStats(c *fiber.Ctx) error {
	// Parse path and query parameters
	levelID := c.Params("levelId")
	if levelID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "levelId is required",
			},
		})
	}

	mapID := c.Query("mapId", "")

	// Get level stats
	stats, err := h.leaderboardService.GetLevelStats(c.Context(), levelID, mapID)
	if err != nil {
		log.Printf("Failed to get level stats for %s: %v", levelID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to retrieve level statistics",
			},
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    stats,
	})
}
