package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/riad/banksystemendtoend/api"
	logger "github.com/riad/banksystemendtoend/pkg/log"
	environment_config "github.com/riad/banksystemendtoend/util/config"
	setup "github.com/riad/banksystemendtoend/util/db"
	"go.uber.org/zap"
)

// ! Application represents the main application configuration
type Application struct {
	server *api.Server
}

// ! NewApplication initializes a new application instance
func NewApplication() (*Application, error) {

	logger.GetLogger().Info("Starting wallet system application")
	server, err := api.NewServer()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize server: %w", err)
	}

	return &Application{
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

	logger.GetLogger().Info("✅ Database connection established")

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		if err := app.server.Start(fmt.Sprintf(":%s", port)); err != nil {
			logger.GetLogger().Error("Failed to start server", zap.Error(err))
		}
	}()

	return app.waitForShutdown()
}

// ! waitForShutdown handles graceful shutdown logic
func (app *Application) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.GetLogger().Info("🔄 Shutting down server...")

	return nil
}

// !main is the entry point of the application
func main() {
	if err := run(); err != nil {
		fmt.Printf("❌ Error occurred: %v\n", err)
		os.Exit(1)
	}
}
