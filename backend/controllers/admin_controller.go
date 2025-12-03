package controllers

import (
	"net/http"
	"strconv"
	"time"

	"graduation_invitation/backend/config"
	"graduation_invitation/backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ==================== USER MANAGEMENT ====================

// GET /api/admin/users - Lấy danh sách users với phân trang
func AdminGetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	offset := (page - 1) * limit

	var users []models.User
	var total int64

	query := config.DB.Model(&models.User{})

	// Tìm kiếm theo email hoặc tên
	if search != "" {
		query = query.Where("email ILIKE ? OR full_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể lấy danh sách users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GET /api/admin/users/:id - Lấy chi tiết user
func AdminGetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User không tồn tại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// POST /api/admin/users - Tạo user mới
func AdminCreateUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name" binding:"required"`
		Phone    string `json:"phone"`
		Role     string `json:"role" binding:"required,oneof=admin user"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dữ liệu không hợp lệ",
			"error":   err.Error(),
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể tạo user",
		})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone:    req.Phone,
		Role:     req.Role,
		Avatar:   req.Avatar,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể lưu user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tạo user thành công",
		"data":    user,
	})
}

// PUT /api/admin/users/:id - Cập nhật user
func AdminUpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User không tồn tại",
		})
		return
	}

	var req struct {
		Email    string `json:"email" binding:"omitempty,email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		Phone    string `json:"phone"`
		Role     string `json:"role" binding:"omitempty,oneof=admin user"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dữ liệu không hợp lệ",
			"error":   err.Error(),
		})
		return
	}

	// Cập nhật các trường
	if req.Email != "" && req.Email != user.Email {
		// Kiểm tra email mới có trùng không
		var existing models.User
		if err := config.DB.First(&existing, "email = ? AND id != ?", req.Email, id).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Email đã được sử dụng",
			})
			return
		}
		user.Email = req.Email
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	// Cập nhật mật khẩu nếu có
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Không thể cập nhật mật khẩu",
			})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể cập nhật user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cập nhật user thành công",
		"data":    user,
	})
}

// DELETE /api/admin/users/:id - Xóa user
func AdminDeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Không cho phép xóa chính mình
	currentUser, _ := c.Get("user")
	if currentUser.(models.User).ID == parseUint(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Không thể xóa chính mình",
		})
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User không tồn tại",
		})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể xóa user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Xóa user thành công",
	})
}

// ==================== RSVP MANAGEMENT ====================

// GET /api/admin/rsvps - Lấy danh sách RSVPs với phân trang
func AdminGetRSVPs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	search := c.Query("search")
	offset := (page - 1) * limit

	var rsvps []models.RSVP
	var total int64

	query := config.DB.Model(&models.RSVP{}).Preload("User")

	// Filter theo status
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Tìm kiếm theo tên guest hoặc message
	if search != "" {
		query = query.Where("guest_name ILIKE ? OR message ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&rsvps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể lấy danh sách RSVPs",
		})
		return
	}

	// ✅ Transform data để thêm thông tin nhận biết user
	type RSVPResponse struct {
		ID           uint         `json:"id"`
		UserID       *uint        `json:"user_id"`
		User         *models.User `json:"user,omitempty"`
		GuestName    string       `json:"guest_name"`
		GuestEmail   string       `json:"guest_email"`
		GuestPhone   string       `json:"guest_phone"`
		Status       string       `json:"status"`
		GuestCount   int          `json:"guest_count"`
		Message      string       `json:"message"`
		IsLoggedIn   bool         `json:"is_logged_in"`
		DisplayName  string       `json:"display_name"`
		DisplayEmail string       `json:"display_email"`
		DisplayPhone string       `json:"display_phone"`
		CreatedAt    time.Time    `json:"created_at"`
		UpdatedAt    time.Time    `json:"updated_at"`
	}

	response := make([]RSVPResponse, 0, len(rsvps))
	for _, rsvp := range rsvps {
		item := RSVPResponse{
			ID:         rsvp.ID,
			UserID:     rsvp.UserID,
			GuestName:  rsvp.GuestName,
			GuestEmail: rsvp.GuestEmail,
			GuestPhone: rsvp.GuestPhone,
			Status:     rsvp.Status,
			GuestCount: rsvp.GuestCount,
			Message:    rsvp.Message,
			CreatedAt:  rsvp.CreatedAt,
			UpdatedAt:  rsvp.UpdatedAt,
		}

		// ✅ Kiểm tra user đã đăng nhập hay chưa
		if rsvp.UserID != nil && rsvp.User.ID != 0 {
			item.IsLoggedIn = true
			item.User = &rsvp.User
			item.DisplayName = rsvp.User.FullName
			item.DisplayEmail = rsvp.User.Email
			item.DisplayPhone = rsvp.User.Phone
		} else {
			item.IsLoggedIn = false
			item.DisplayName = rsvp.GuestName
			item.DisplayEmail = rsvp.GuestEmail
			item.DisplayPhone = rsvp.GuestPhone
		}

		response = append(response, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GET /api/admin/rsvps/:id - Lấy chi tiết RSVP
func AdminGetRSVP(c *gin.Context) {
	id := c.Param("id")
	var rsvp models.RSVP

	if err := config.DB.Preload("User").First(&rsvp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "RSVP không tồn tại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rsvp,
	})
}

// PUT /api/admin/rsvps/:id - Cập nhật RSVP
func AdminUpdateRSVP(c *gin.Context) {
	id := c.Param("id")
	var rsvp models.RSVP

	if err := config.DB.First(&rsvp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "RSVP không tồn tại",
		})
		return
	}

	var req struct {
		Status     string `json:"status" binding:"omitempty,oneof=yes no maybe"`
		GuestCount int    `json:"guest_count"`
		Message    string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dữ liệu không hợp lệ",
			"error":   err.Error(),
		})
		return
	}

	// Cập nhật các trường
	if req.Status != "" {
		rsvp.Status = req.Status
	}

	if req.GuestCount > 0 {
		rsvp.GuestCount = req.GuestCount
	}

	if req.Message != "" {
		rsvp.Message = req.Message
	}

	if err := config.DB.Save(&rsvp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể cập nhật RSVP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cập nhật RSVP thành công",
		"data":    rsvp,
	})
}

// DELETE /api/admin/rsvps/:id - Xóa RSVP
func AdminDeleteRSVP(c *gin.Context) {
	id := c.Param("id")
	var rsvp models.RSVP

	if err := config.DB.First(&rsvp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "RSVP không tồn tại",
		})
		return
	}

	if err := config.DB.Delete(&rsvp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể xóa RSVP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Xóa RSVP thành công",
	})
}

// ==================== DASHBOARD STATS ====================

// GET /api/admin/dashboard - Lấy thống kê tổng quan
func AdminGetDashboard(c *gin.Context) {
	var totalUsers int64
	var totalRSVPs int64
	var yesRSVPs int64
	var noRSVPs int64
	var maybeRSVPs int64

	config.DB.Model(&models.User{}).Count(&totalUsers)
	config.DB.Model(&models.RSVP{}).Count(&totalRSVPs)
	config.DB.Model(&models.RSVP{}).Where("status = ?", "yes").Count(&yesRSVPs)
	config.DB.Model(&models.RSVP{}).Where("status = ?", "no").Count(&noRSVPs)
	config.DB.Model(&models.RSVP{}).Where("status = ?", "maybe").Count(&maybeRSVPs)

	// Lấy RSVPs gần đây
	var recentRSVPs []models.RSVP
	config.DB.Preload("User").Order("created_at desc").Limit(5).Find(&recentRSVPs)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"totalUsers": totalUsers,
			"totalRSVPs": totalRSVPs,
			"stats": gin.H{
				"yes":   yesRSVPs,
				"no":    noRSVPs,
				"maybe": maybeRSVPs,
			},
			"recentRSVPs": recentRSVPs,
		},
	})
}

// Helper function
func parseUint(s string) uint {
	val, _ := strconv.ParseUint(s, 10, 32)
	return uint(val)
}
