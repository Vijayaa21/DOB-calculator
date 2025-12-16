package routes

import (
	"github.com/dob-calculator/internal/handler"
	"github.com/dob-calculator/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// Setup configures all routes for the application
func Setup(app *fiber.App, userHandler *handler.UserHandler) {
	// Apply global middleware
	app.Use(middleware.CORS())
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// User routes
	users := api.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.ListUsers)
	users.Get("/:id", userHandler.GetUserByID)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)

	// Also register routes without /api/v1 prefix for backward compatibility
	app.Post("/users", userHandler.CreateUser)
	app.Get("/users", userHandler.ListUsers)
	app.Get("/users/:id", userHandler.GetUserByID)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)
}
