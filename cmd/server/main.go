package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dob-calculator/config"
	"github.com/dob-calculator/db"
	"github.com/dob-calculator/internal/handler"
	"github.com/dob-calculator/internal/logger"
	"github.com/dob-calculator/internal/middleware"
	"github.com/dob-calculator/internal/repository"
	"github.com/dob-calculator/internal/routes"
	"github.com/dob-calculator/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Init(); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting DOB Calculator API")


	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	database, err := db.Connect(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	// Initialize layers
	userRepo := repository.NewUserRepository(database.Queries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		AppName:      "DOB Calculator API v1.0.0",
	})

	// Setup routes
	routes.Setup(app, userHandler)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			logger.Error("Server shutdown error", zap.Error(err))
		}
	}()

	addr := ":" + cfg.Server.Port
	logger.Info("Server starting", zap.String("address", addr))
	if err := app.Listen(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
