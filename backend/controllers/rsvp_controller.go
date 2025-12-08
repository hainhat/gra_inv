package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"graduation_invitation/backend/config"
	"graduation_invitation/backend/models"
	"graduation_invitation/backend/utils"

	"github.com/gin-gonic/gin"
)

// POST /api/rsvp
func SubmitRSVP(c *gin.Context) {
	var req struct {
		GuestName      string `json:"guest_name"`
		GuestEmail     string `json:"guest_email"`
		GuestPhone     string `json:"guest_phone"`
		Status         string `json:"status"`
		GuestCount     int    `json:"guest_count"`
		Message        string `json:"message"`
		RecaptchaToken string `json:"recaptcha_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dữ liệu gửi lên không hợp lệ.",
		})
		return
	}

	// ✅ Verify reCAPTCHA token
	valid, err := utils.VerifyRecaptcha(req.RecaptchaToken)
	if err != nil {
		log.Printf("❌ reCAPTCHA verification error: %v", err)
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Xác minh reCAPTCHA thất bại. Vui lòng thử lại.",
		})
		return
	}

	if req.Status == "" {
		req.Status = "yes"
	}

	rsvp := models.RSVP{
		GuestName:  req.GuestName,
		GuestEmail: req.GuestEmail,
		GuestPhone: req.GuestPhone,
		Status:     req.Status,
		GuestCount: req.GuestCount,
		Message:    req.Message,
	}

	// ✅ Giải mã token và gắn user_id nếu có
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(token)
		if err == nil && claims != nil {
			if id, ok := claims["id"].(float64); ok {
				userID := uint(id)
				rsvp.UserID = &userID
				rsvp.GuestName = ""
				rsvp.GuestEmail = ""
				rsvp.GuestPhone = ""
			}
		}
	}

	// ✅ Lưu vào DB
	if err := config.DB.Create(&rsvp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể lưu RSVP. Vui lòng thử lại sau.",
		})
		return
	}

	// ✅ Gửi email xác nhận (bất đồng bộ)
	if req.GuestEmail != "" {
		go func() {
			err := utils.SendRSVPConfirmation(req.GuestEmail, req.GuestName)
			if err != nil {
				log.Printf("❌ Failed to send email to %s: %v", req.GuestEmail, err)
			} else {
				log.Printf("✅ Email sent successfully to %s", req.GuestEmail)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		//"message": "Cảm ơn bạn đã phản hồi!",
	})
}

// GET /api/rsvp/stats - Public endpoint for RSVP statistics
func GetStats(c *gin.Context) {
	var total int64
	var yes int64
	var no int64
	var maybe int64

	// Count total
	config.DB.Model(&models.RSVP{}).Count(&total)

	// Count by status
	config.DB.Model(&models.RSVP{}).Where("status = ?", "yes").Count(&yes)
	config.DB.Model(&models.RSVP{}).Where("status = ?", "no").Count(&no)
	config.DB.Model(&models.RSVP{}).Where("status = ?", "maybe").Count(&maybe)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total": total,
			"yes":   yes,
			"no":    no,
			"maybe": maybe,
		},
	})
}

// GET /api/rsvp/messages - Public endpoint to get RSVP messages with pagination
func GetRSVPMessages(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit

	var rsvps []models.RSVP
	var total int64

	// Count total messages
	config.DB.Model(&models.RSVP{}).Where("message != ?", "").Count(&total)

	// Get RSVPs with messages, ordered by newest first
	if err := config.DB.Preload("User").
		Where("message != ?", "").
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&rsvps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể lấy danh sách messages.",
		})
		return
	}

	// Transform data for frontend
	type MessageResponse struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Avatar    string `json:"avatar"`
		Message   string `json:"message"`
		CreatedAt string `json:"created_at"`
	}

	messages := make([]MessageResponse, 0, len(rsvps))
	for _, rsvp := range rsvps {
		msg := MessageResponse{
			ID:        rsvp.ID,
			Message:   rsvp.Message,
			CreatedAt: rsvp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Use user info if available, otherwise use guest info
		if rsvp.UserID != nil && rsvp.User.ID != 0 {
			msg.Name = rsvp.User.FullName
			msg.Avatar = rsvp.User.Avatar
		} else {
			msg.Name = rsvp.GuestName
			// Default avatar for guests (using UI Avatars)
			msg.Avatar = ""
		}

		messages = append(messages, msg)
	}

	// Calculate total pages
	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    messages,
		"total":   total,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}
