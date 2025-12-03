package controllers

import (
	"graduation_invitation/backend/config"
	"graduation_invitation/backend/models"
	"graduation_invitation/backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest là dữ liệu gửi từ frontend
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// POST /api/login
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid input"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "email = ?", req.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not found"})
		return
	}

	// Kiểm tra mật khẩu
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Wrong password"})
		return
	}

	// Sinh token JWT (giữ nguyên để tương thích)
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create token"})
		return
	}

	// Sinh access token và refresh token mới
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create refresh token"})
		return
	}

	// Lưu refresh token vào database
	user.RefreshToken = refreshToken
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"token":         token,        // Giữ nguyên để tương thích với code cũ
		"access_token":  accessToken,  // Token mới với thời gian ngắn
		"refresh_token": refreshToken, // Token để làm mới
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})
}

// GET /api/me
func Me(c *gin.Context) {
	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Not logged in"})
		return
	}

	user := userCtx.(models.User)

	// Check if user has RSVP'd
	var count int64
	config.DB.Model(&models.RSVP{}).Where("user_id = ?", user.ID).Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"full_name": user.FullName,
			"phone":     user.Phone,
			"avatar":    user.Avatar,
			"role":      user.Role,
			"has_rsvp":  count > 0,
		},
	})
}

// RegisterRequest chứa dữ liệu khi đăng ký
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone"`
}

// POST /api/register
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dữ liệu không hợp lệ hoặc thiếu trường bắt buộc",
		})
		return
	}

	// Kiểm tra email trùng
	var existing models.User
	if err := config.DB.First(&existing, "email = ?", req.Email).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Email đã được sử dụng",
		})
		return
	}

	// Hash mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Không thể tạo tài khoản"})
		return
	}

	// Tạo user mới
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone:    req.Phone,
		Role:     "user", // mặc định là user thường
	}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Không thể lưu người dùng"})
		return
	}

	// Sinh access token và refresh token mới
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create refresh token"})
		return
	}

	// Lưu refresh token vào database
	user.RefreshToken = refreshToken
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save refresh token"})
		return
	}

	// Trả về user và token
	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "Đăng ký thành công",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})
}

// GET /api/check-email?email=...
func CheckEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(400, gin.H{"exists": false})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "email = ?", email).Error; err == nil {
		c.JSON(200, gin.H{"exists": true})
	} else {
		c.JSON(200, gin.H{"exists": false})
	}
}

// POST /api/refresh - Làm mới access token
func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Refresh token is required",
		})
		return
	}

	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid or expired refresh token",
		})
		return
	}

	// Lấy user ID từ claims
	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid token claims",
		})
		return
	}
	userID := uint(userIDFloat)

	// Kiểm tra refresh token có khớp với database không
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Kiểm tra refresh token có khớp với token đã lưu
	if user.RefreshToken != req.RefreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Refresh token không hợp lệ hoặc đã bị thu hồi",
		})
		return
	}

	// Tạo access token mới
	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"access_token": newAccessToken,
	})
}

// POST /api/logout - Thu hồi refresh token
func Logout(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not logged in",
		})
		return
	}

	currentUser := user.(models.User)

	// Xóa refresh token khỏi database
	if err := config.DB.Model(&models.User{}).Where("id = ?", currentUser.ID).Update("refresh_token", "").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Đăng xuất thành công",
	})
}
