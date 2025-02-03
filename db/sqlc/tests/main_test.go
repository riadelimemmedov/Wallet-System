// File: db/db_test.go
package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/riad/banksystemendtoend/util/config"
	setup "github.com/riad/banksystemendtoend/util/db"
)

var ctx = context.Background()

// Main test run before all tests
func TestMain(m *testing.M) {
	if err := setup.InitializeEnvironment(config.TestEnvironment); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	code := m.Run()
	os.Exit(code)
}

// CleanupDB provides a helper function for cleaning up the test database between tests.
func CleanupDB(t *testing.T) {
	t.Helper()
	if err := setup.DropAllData(ctx, setup.GetStore()); err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}
}
