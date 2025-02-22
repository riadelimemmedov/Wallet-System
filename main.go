package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/riad/banksystemendtoend/api"
	"github.com/riad/banksystemendtoend/api/middleware"
	"github.com/riad/banksystemendtoend/api/middleware/logger/config"
	environment_config "github.com/riad/banksystemendtoend/util/config"
	setup "github.com/riad/banksystemendtoend/util/db"
	"go.uber.org/zap"
)

// ! Application represents the main application configuration
type Application struct {
	logger *middleware.Logger
	server *api.Server
}

// ! NewApplication initializes a new application instance
func NewApplication() (*Application, error) {
	logger, err := middleware.NewLogger(config.LoggerConfig{
		Filename:   "logs/app.log",
		TimeFormat: "2006-01-02 15:04:05",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	server, err := api.NewServer()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize server: %w", err)
	}

	return &Application{
		logger: logger,
		server: server,
	}, nil
}

// ! run contains the core application logic and handles graceful shutdown
func run() error {
	if err := setup.InitializeEnvironment(environment_config.DevEnvironment); err != nil {
		return err
	}

	app, err := NewApplication()
	if err != nil {
		return err
	}

	app.logger.ZapLogger.Info("✅ Database connection established")

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		if err := app.server.Start(":8080"); err != nil {
			app.logger.ZapLogger.Error("Failed to start server", zap.Error(err))
		}
	}()

	return app.waitForShutdown()
}

// ! waitForShutdown handles graceful shutdown logic
func (app *Application) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	app.logger.ZapLogger.Info("🔄 Shutting down server...")

	return nil
}

// !main is the entry point of the application
func main() {
	if err := run(); err != nil {
		fmt.Printf("❌ Error occurred: %v\n", err)
		os.Exit(1)
	}
}
