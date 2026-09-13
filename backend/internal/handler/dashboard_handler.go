package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mobi07/secure_bank_auth/internal/constants"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID, exists := c.Get(constants.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	role, exists := c.Get(constants.RoleKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "welcome to your dashboard",
		"user_id": userID,
		"role":    role,
	})

}
