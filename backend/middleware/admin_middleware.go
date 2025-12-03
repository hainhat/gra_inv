package middleware

import (
	"net/http"

	"graduation_invitation/backend/models"

	"github.com/gin-gonic/gin"
)

// RequireAdmin kiểm tra user phải có role admin
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			return
		}

		user := userCtx.(models.User)
		if user.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Chỉ admin mới có quyền truy cập",
			})
			return
		}

		c.Next()
	}
}
