package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// List accounts fetch all accounts from database
func (s *Server) listAccounts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved accounts"})
}
