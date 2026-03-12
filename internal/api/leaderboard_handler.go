package api

import (
	"log"
	"strconv"
	"strings"

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

// GetHollowWildsLeaderboard handles the new Hollow Wilds leaderboard request
// @Summary Get HW Leaderboard
// @Description Get ranked entries for a specific metric and scope
// @Tags Hollow Wilds
// @Produce json
// @Param type query string false "Metric type (longest_run_days, sebilah_soul_level, bosses_killed)" default(longest_run_days)
// @Param scope query string false "Scope (global, per_character)" default(global)
// @Param character query string false "Character filter (required if scope=per_character)"
// @Param limit query int false "Limit" default(100)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.HollowWildsLeaderboardResponse "Leaderboard data"
// @Router /leaderboard [get]
func (h *LeaderboardHandler) GetHollowWildsLeaderboard(c *fiber.Ctx) error {
	lbType := c.Query("type", "longest_run_days")
	scope := c.Query("scope", "global")
	character := c.Query("character", "")
	limit := c.QueryInt("limit", 100)
	offset := c.QueryInt("offset", 0)

	if scope == "per_character" && character == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "character is required for per_character scope",
			},
		})
	}

	leaderboard, err := h.leaderboardService.GetHollowWildsLeaderboard(c.Context(), lbType, scope, character, limit, offset)
	if err != nil {
		log.Printf("Failed to get Hollow Wilds leaderboard: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to retrieve leaderboard",
			},
		})
	}

	return c.JSON(leaderboard)
}

// SubmitHollowWildsEntry handles leaderboard submission
// @Summary Submit HW Run
// @Description Submit a result after a run to update personal best and ranks
// @Tags Hollow Wilds
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.LeaderboardSubmitRequest true "Run result"
// @Success 200 {object} models.LeaderboardSubmitResponse "Submission result"
// @Failure 400 {object} models.APIResponse{error=models.APIError} "Value too low or invalid request"
// @Router /leaderboard/submit [post]
func (h *LeaderboardHandler) SubmitHollowWildsEntry(c *fiber.Ctx) error {
	playerIDStr, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeUnauthorized,
				Message: "User not authenticated",
			},
		})
	}
	playerID, _ := uuid.Parse(playerIDStr)

	var req models.LeaderboardSubmitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInvalidRequest,
				Message: "Invalid request body",
			},
		})
	}

	result, err := h.leaderboardService.SubmitHollowWildsEntry(c.Context(), playerID, req)
	if err != nil {
		log.Printf("Failed to submit leaderboard entry: %v", err)
		if strings.Contains(err.Error(), "value_too_low") {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "value_too_low",
					Message: "Submitted value does not beat personal best",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to submit entry",
			},
		})
	}

	return c.JSON(result)
}

// GetPlayerHollowWildsStats handles request for player's own ranks
// @Summary Get HW Player Ranks
// @Description Get current ranks for the authenticated player across all types
// @Tags Hollow Wilds
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.PlayerLeaderboardResponse "Player rankings"
// @Router /leaderboard/player [get]
func (h *LeaderboardHandler) GetPlayerHollowWildsStats(c *fiber.Ctx) error {
	playerIDStr, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeUnauthorized,
				Message: "User not authenticated",
			},
		})
	}
	playerID, _ := uuid.Parse(playerIDStr)

	stats, err := h.leaderboardService.GetPlayerHollowWildsStats(c.Context(), playerID)
	if err != nil {
		log.Printf("Failed to get player Hollow Wilds stats: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    models.ErrCodeInternalError,
				Message: "Failed to retrieve player stats",
			},
		})
	}

	return c.JSON(stats)
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
