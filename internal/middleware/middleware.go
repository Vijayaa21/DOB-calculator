package middleware

import (
	"time"

	"github.com/dob-calculator/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestID middleware adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request ID exists in header
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set the request ID in locals for access in handlers
		c.Locals("requestId", requestID)

		// Set the request ID in response header
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

// RequestLogger middleware logs request details and duration
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Get request ID
		requestID, _ := c.Locals("requestId").(string)

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request details
		logger.Info("HTTP Request",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		)

		return err
	}
}

// ErrorHandler is a custom error handler for Fiber
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Get request ID
	requestID, _ := c.Locals("requestId").(string)

	// Default status code
	code := fiber.StatusInternalServerError

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Log the error
	logger.Error("Request error",
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.Int("status", code),
		zap.Error(err),
	)

	// Return JSON error response
	return c.Status(code).JSON(fiber.Map{
		"error":      "internal_error",
		"message":    err.Error(),
		"request_id": requestID,
	})
}

// CORS middleware handles Cross-Origin Resource Sharing
func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
