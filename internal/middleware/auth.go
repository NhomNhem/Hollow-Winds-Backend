package middleware

import (
	"strings"

	"github.com/NhomNhem/GameFeel-Backend/internal/models"
	"github.com/NhomNhem/GameFeel-Backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware() fiber.Handler {
	authService := services.NewAuthService()

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeUnauthorized,
					Message: "Authorization header required",
				},
			})
		}

		// Extract token from "Bearer <token>"
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			// No "Bearer " prefix
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeUnauthorized,
					Message: "Invalid authorization format. Use: Bearer <token>",
				},
			})
		}

		// Verify JWT
		claims, err := authService.VerifyJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeInvalidToken,
					Message: "JWT token is invalid or expired",
				},
			})
		}

		// Store user info in context for downstream handlers
		c.Locals("userId", claims.UserID)
		c.Locals("playfabId", claims.PlayFabID)

		return c.Next()
	}
}
