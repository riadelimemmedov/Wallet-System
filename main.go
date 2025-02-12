// main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/riad/banksystemendtoend/util/config"
	setup "github.com/riad/banksystemendtoend/util/db"
)

// ! Run contains the core application logic and handles graceful shutdown
func run() error {
	if err := setup.InitializeEnvironment(config.DevEnvironment); err != nil {
		return err
	}
	log.Println("✅ Database connection established")

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("🔄 Shutting down server...")

	return nil
}

// ! Run core service
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if err := run(); err != nil {
		log.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}
