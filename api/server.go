package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/riad/banksystemendtoend/api/dependency"
	"github.com/riad/banksystemendtoend/api/middleware"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	setup "github.com/riad/banksystemendtoend/util/db"
	"go.uber.org/zap"
)

type Server struct {
	router       *gin.Engine
	store        db.Store
	dependencies *dependency.DependencyContainer
}

// NewServer creates and configures a new server instance
func NewServer() (*Server, error) {
	store, err := db.GetSQLStore(setup.GetStore())
	if err != nil {
		return nil, err
	}

	dependencies := dependency.NewDependencyContainer(store)

	server := &Server{
		store:        store,
		dependencies: dependencies,
	}

	server.setupRouter()
	return server, nil
}

// setupRouter configures all the API routes and middleware
func (s *Server) setupRouter() {
	router := gin.Default()

	apiKey, err := middleware.NewAPIKey()

	if err != nil {
		zap.L().Fatal("Failed to create API key", zap.Error(err))
	}

	router.Use(middleware.Cors())
	router.Use(middleware.TimeOut(3 * time.Minute))
	router.Use(apiKey.ValidateAPIKey())

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		//Account Routes
		accounts := v1.Group("/accounts")
		{
			accounts.GET("", s.listAccounts)
		}

		// Account Type Routes - dynamically register from dependency container
		accountTypes := v1.Group("/account-types")
		for _, route := range s.dependencies.GetRouteHandlers("account-types") {
			accountTypes.Handle(route.Method, route.Path, route.HandlerFunc)
		}
	}

	//Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})
	s.router = router
}

// Start launches the HTTP server on the specified address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
