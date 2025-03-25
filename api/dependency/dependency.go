package dependency

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/riad/banksystemendtoend/api/handler"
	handler_interface "github.com/riad/banksystemendtoend/api/interface/handler"
	"github.com/riad/banksystemendtoend/api/repository"
	"github.com/riad/banksystemendtoend/api/service"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/pkg/cache"
	"github.com/riad/banksystemendtoend/pkg/redis"
)

type DependencyContainer struct {
	handlers     map[string][]RouteHandler
	redisClient  *redis.Client
	cacheService *cache.Service

	AccountTypeHandler handler_interface.AccountTypeHandler
}

type RouteHandler struct {
	Method      string
	Path        string
	HandlerFunc gin.HandlerFunc
}

func NewDependencyContainer(store db.Store, redisConfig redis.Config) (*DependencyContainer, error) {
	// Initialize Redis client
	redisClient, err := redis.NewClient(redisConfig)
	if err != nil {
		return nil, err
	}

	cacheService := cache.NewService(redisClient, "wallet_app", 60*time.Minute)

	container := &DependencyContainer{
		handlers:     make(map[string][]RouteHandler),
		redisClient:  redisClient,
		cacheService: cacheService,
	}

	container.registerAccountTypeHandlers(store, cacheService)
	return container, nil
}

func (c *DependencyContainer) registerAccountTypeHandlers(store db.Store, cacheService *cache.Service) {
	accountTypeRepo := repository.NewAccountTypeRepository(store, cacheService)
	accountTypeService := service.NewAccountTypeService(accountTypeRepo)
	accountTypeHandler := handler.NewAccountTypeHandler(accountTypeService)

	c.AccountTypeHandler = accountTypeHandler

	c.handlers["account-types"] = []RouteHandler{
		{
			Method:      http.MethodPost,
			Path:        "",
			HandlerFunc: accountTypeHandler.CreateAccountType,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:account_type",
			HandlerFunc: accountTypeHandler.GetAccountType,
		},
		{
			Method:      http.MethodGet,
			Path:        "",
			HandlerFunc: accountTypeHandler.ListAccountTypes,
		},
		{
			Method:      http.MethodPut,
			Path:        "/:account_type",
			HandlerFunc: accountTypeHandler.UpdateAccountType,
		},
		{
			Method:      http.MethodPatch,
			Path:        "/:account_type",
			HandlerFunc: accountTypeHandler.UpdateAccountType,
		},
		{
			Method:      http.MethodDelete,
			Path:        "/:account_type",
			HandlerFunc: accountTypeHandler.DeleteAccountType,
		},
	}
}

func (c *DependencyContainer) GetRouteHandlers(groupPrefix string) []RouteHandler {
	return c.handlers[groupPrefix]
}

func (c *DependencyContainer) GetCacheService() *cache.Service {
	return c.cacheService
}

func (c *DependencyContainer) Close() error {
	if c.redisClient != nil {
		return c.redisClient.Close()
	}
	return nil
}
