package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "GameFeel Backend v1.0.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    code,
					"message": err.Error(),
				},
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: getEnv("ALLOWED_ORIGINS", "*"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "GameFeel Backend is running",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// Root endpoint with API info
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "GameFeel API v1",
			"endpoints": []string{
				"GET  /health",
				"GET  /api/v1/",
				"POST /api/v1/auth/login",
				"POST /api/v1/levels/complete",
				"POST /api/v1/talents/upgrade",
				"POST /api/v1/payments/create-session",
				"POST /api/v1/analytics/events",
			},
		})
	})

	// TODO: Register route handlers
	// auth := api.Group("/auth")
	// levels := api.Group("/levels")
	// talents := api.Group("/talents")
	// payments := api.Group("/payments")
	// analytics := api.Group("/analytics")

	// Get port from env or default to 8080
	port := getEnv("PORT", "8080")

	// Start server
	log.Printf("🚀 Server starting on port %s...", port)
	log.Printf("📝 Environment: %s", getEnv("ENV", "development"))
	log.Printf("🔗 Health check: http://localhost:%s/health", port)
	log.Printf("🔗 API docs: http://localhost:%s/api/v1/", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
