package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/riad/banksystemendtoend/api/dependency"
	"github.com/riad/banksystemendtoend/api/middleware"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	logger "github.com/riad/banksystemendtoend/pkg/log"
	cache_setup "github.com/riad/banksystemendtoend/util/cache"
	db_setup "github.com/riad/banksystemendtoend/util/db"
	"go.uber.org/zap"
)

type Server struct {
	router       *gin.Engine
	store        db.Store
	dependencies *dependency.DependencyContainer
}

// NewServer creates and configures a new server instance
func NewServer() (*Server, error) {
	// Get the initialized store and Redis client
	store, err := db.GetSQLStore(db_setup.GetStore())
	if err != nil {
		logger.GetLogger().Error("Failed to get SQL store", zap.Error(err))
		return nil, err
	}

	// Get the Redis client that was initialized in InitializeEnvironment
	redisClient := cache_setup.GetRedisClient()

	// Create dependencies container with DbStore and Redis client
	dependencies, err := dependency.NewDependencyContainer(store, redisClient)
	if err != nil {
		logger.GetLogger().Error("Failed to create dependency container", zap.Error(err))
		return nil, err
	}

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
		cacheService := s.dependencies.GetCacheService()
		store, _ := db.GetSQLStore(db_setup.GetStore())

		redisStatus := "up"
		dbStatus := "up"

		if cacheService != nil {
			is_connect := cacheService.CheckRedisConnection()
			if !is_connect {
				redisStatus = "down"
				logger.GetLogger().Error("Failed to connect Redis", zap.Error(err))
			}
		} else {
			redisStatus = "not_configured"
			logger.GetLogger().Warn("Redis service not configured")
		}

		// Database connection test
		err = db_setup.CheckDBHealth(c.Request.Context(), store)
		if err != nil {
			dbStatus = "down"
			logger.GetLogger().Error("Failed to connect DB", zap.Error(err))
		}

		c.JSON(http.StatusOK, gin.H{"api": "up", "redis": redisStatus, "database": dbStatus})
	})
	s.router = router
}

// Start launches the HTTP server on the specified address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
