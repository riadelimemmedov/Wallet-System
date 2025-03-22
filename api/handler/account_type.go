package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/riad/banksystemendtoend/api"
	"github.com/riad/banksystemendtoend/api/dto"
	handler_interface "github.com/riad/banksystemendtoend/api/interface/handler"
	interface_service "github.com/riad/banksystemendtoend/api/interface/service"
	"github.com/riad/banksystemendtoend/api/utils"
	"github.com/riad/banksystemendtoend/util/config"
)

type accountTypeHandler struct {
	service interface_service.AccountTypeService
}

func NewAccountTypeHandler(service interface_service.AccountTypeService) handler_interface.AccountTypeHandler {
	return &accountTypeHandler{service: service}
}

func (h *accountTypeHandler) CreateAccountType(ctx *gin.Context) {
	var req dto.CreateAccountTypeRequest
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	inputAccountType := strings.ToUpper(req.AccountType)
	if !config.IsValidAccountType(inputAccountType) {
		err = fmt.Errorf(utils.GetValidAccountTypesMessage())
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	createdAccountType, err := h.service.CreateAccountType(ctx, inputAccountType, req.Description)
	if err != nil {
		if utils.IsDuplicateError(err) {
			ctx.JSON(http.StatusConflict, api.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	rsp := dto.AccountTypeResponse{
		AccountType: createdAccountType.AccountType,
		Description: createdAccountType.Description,
		IsActive:    createdAccountType.IsActive,
		CreatedAt:   createdAccountType.CreatedAt,
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": rsp})
}

func (h *accountTypeHandler) GetAccountType(ctx *gin.Context) {
	accountType := ctx.Param("accountType")
	if accountType == "" {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(fmt.Errorf("account type is required")))
		return
	}

	accountTypeData, err := h.service.GetAccountType(ctx, accountType)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, api.ErrorResponse(fmt.Errorf("account type not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}
	rsp := dto.AccountResponse{
		AccountType: accountTypeData.AccountType,
		IsActive:    accountTypeData.IsActive,
		CreatedAt:   accountTypeData.CreatedAt,
	}
	ctx.JSON(http.StatusOK, gin.H{"data": rsp})
}

func (h *accountTypeHandler) ListAccountTypes(ctx *gin.Context) {
	accountTypes, err := h.service.ListAccountTypes(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	var rsp []dto.AccountTypeResponse
	for _, accountType := range accountTypes {
		rsp = append(rsp, dto.AccountTypeResponse{
			AccountType: accountType.AccountType,
			Description: accountType.Description,
			IsActive:    accountType.IsActive,
			CreatedAt:   accountType.CreatedAt,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"data": rsp})
}

// ! Implement later
func (h *accountTypeHandler) UpdateAccountType(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "UpdateAccountType"})
}

func (h *accountTypeHandler) DeleteAccountType(ctx *gin.Context) {
	accountType := ctx.Param("accountType")
	if accountType == "" {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(fmt.Errorf("account type is required")))
		return
	}

	err := h.service.DeleteAccountType(ctx, accountType)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, api.ErrorResponse(fmt.Errorf("account type not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	message := fmt.Sprintf("Account type %s deleted successfully", accountType)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}
