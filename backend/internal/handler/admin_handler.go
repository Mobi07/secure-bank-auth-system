package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "admin user management endpoint",
	})
}
