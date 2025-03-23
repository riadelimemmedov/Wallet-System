package dependency

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riad/banksystemendtoend/api/handler"
	handler_interface "github.com/riad/banksystemendtoend/api/interface/handler"
	"github.com/riad/banksystemendtoend/api/repository"
	"github.com/riad/banksystemendtoend/api/service"
	db "github.com/riad/banksystemendtoend/db/sqlc"
)

type DependencyContainer struct {
	handlers map[string][]RouteHandler

	AccountTypeHandler handler_interface.AccountTypeHandler
}

type RouteHandler struct {
	Method      string
	Path        string
	HandlerFunc gin.HandlerFunc
}

func NewDependencyContainer(store db.Store) *DependencyContainer {
	container := &DependencyContainer{
		handlers: make(map[string][]RouteHandler),
	}

	container.registerAccountTypeHandlers(store)
	return container
}

func (c *DependencyContainer) registerAccountTypeHandlers(store db.Store) {
	accountTypeRepo := repository.NewAccountTypeRepository(store)
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
