package middleware

import (
	"github.com/NhomNhem/NhemDangFugBixs-Core/internal/domain/models"
	"github.com/gofiber/fiber/v2"
)

// roleHierarchy defines the order of roles from lowest to highest
var roleHierarchy = map[models.SystemRole]int{
	models.RoleUser:       0,
	models.RoleAdmin:      1,
	models.RoleSuperAdmin: 2,
}

// RoleMiddleware checks if the user has the required minimum role
func RoleMiddleware(minRole models.SystemRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleStr, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    models.ErrCodeUnauthorized,
					Message: "Role not found in context",
				},
			})
		}

		userRole := models.SystemRole(userRoleStr)

		// Check hierarchy
		if roleHierarchy[userRole] < roleHierarchy[minRole] {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "FORBIDDEN",
					Message: "Insufficient permissions",
				},
			})
		}

		return c.Next()
	}
}
