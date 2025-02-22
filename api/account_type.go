package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/riad/banksystemendtoend/api/dto"
	"github.com/riad/banksystemendtoend/api/utils"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/util/config"
)

// ! CreateAccountType creates a new account type
func (server *Server) createAccountType(ctx *gin.Context) {
	var req dto.CreateAccountTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	// Check if the account type is valid
	inputAccountType := strings.ToUpper(req.AccountType)
	if !config.IsValidAccountType(inputAccountType) {
		err := fmt.Errorf("invalid account type: must be one of: SAVINGS, CHECKING, FIXED_DEPOSIT, MONEY_MARKET")
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	// Create account type schema
	arg := db.CreateAccountTypeParams{
		AccountType: inputAccountType,
		Description: req.Description,
	}

	// Check if the account type already exists IsDuplicateError
	createdAccountType, err := server.store.CreateAccountType(ctx, arg)
	if err != nil {
		if utils.IsDuplicateError(err) {
			ctx.JSON(http.StatusConflict, ErrorResponse(fmt.Errorf("account type already exists")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	rsp := dto.AccountTypeResponse{
		AccountType: createdAccountType.AccountType,
		Description: createdAccountType.Description,
		IsActive:    createdAccountType.IsActive,
		CreatedAt:   createdAccountType.CreatedAt,
	}
	ctx.JSON(http.StatusCreated, rsp)
}
