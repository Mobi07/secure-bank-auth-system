package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/mobi07/secure_bank_auth/internal/constants"
	"github.com/mobi07/secure_bank_auth/internal/domain"
)

func GetUserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get(constants.UserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)
	if !ok {
		return 0, false
	}

	return userID, true
}

func GetUserRole(c *gin.Context) (domain.Role, bool) {
	value, exists := c.Get(constants.RoleKey)
	if !exists {
		return "", false
	}

	role, ok := value.(domain.Role)
	if !ok {
		return "", false
	}

	return role, true
}
