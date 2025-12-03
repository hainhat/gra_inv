package middleware

import (
	"net/http"
	"strings"

	"graduation_invitation/backend/config"
	"graduation_invitation/backend/models"
	"graduation_invitation/backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthJWT kiểm tra token và thêm user vào context
func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy header Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing Authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid Authorization format"})
			return
		}

		// Xác minh token
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired token"})
			return
		}

		// Kiểm tra token type (nếu có claim này)
		// Access token mới sẽ có claim "token_type": "access"
		// Token cũ (24h) sẽ không có claim này -> vẫn cho phép để tương thích ngược (hoặc chặn nếu muốn strict)
		if tokenType, ok := claims["token_type"].(string); ok {
			if tokenType != "access" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid token type. Access token required."})
				return
			}
		}

		// Lấy user ID từ claims
		userIDFloat, ok := claims["id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid token claims"})
			return
		}
		userID := uint(userIDFloat)

		// Tìm user trong database (đảm bảo vẫn tồn tại)
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not found"})
			return
		}

		// Lưu thông tin user vào context
		c.Set("user", user)
		c.Next()
	}
}
