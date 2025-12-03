package controllers

import (
	"context"
	"net/http"
	"os"

	"graduation_invitation/backend/config"
	"graduation_invitation/backend/models"
	"graduation_invitation/backend/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

// GoogleCredential represents the JWT credential from Google Identity Services
type GoogleCredential struct {
	Credential string `json:"credential" binding:"required"`
}

// POST /api/auth/google/verify - Verify Google Identity Services JWT token
func VerifyGoogleToken(c *gin.Context) {
	var req GoogleCredential
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	// Verify the JWT token with Google
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	payload, err := idtoken.Validate(context.Background(), req.Credential, clientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid Google token",
		})
		return
	}

	// Extract user info from payload
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	googleID, _ := payload.Claims["sub"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)

	// Check if email is verified
	if !emailVerified {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Email not verified",
		})
		return
	}

	// Find or create user
	var user models.User
	result := config.DB.Where("email = ? OR google_id = ?", email, googleID).First(&user)

	if result.Error != nil {
		// User doesn't exist, create new user
		user = models.User{
			Email:        email,
			FullName:     name,
			Avatar:       picture,
			GoogleID:     googleID,
			AuthProvider: "google",
			Role:         "user",
			Password:     "", // No password for Google users
		}

		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create user",
			})
			return
		}
	} else {
		// User exists, update Google ID if not set
		if user.GoogleID == "" {
			user.GoogleID = googleID
			user.AuthProvider = "google"
			config.DB.Save(&user)
		}
	}

	// Generate JWT tokens
	accessToken, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate access token",
		})
		return
	}

	refreshToken, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate refresh token",
		})
		return
	}

	// Save refresh token
	user.RefreshToken = refreshToken
	config.DB.Save(&user)

	// Return tokens
	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"full_name": user.FullName,
			"avatar":    user.Avatar,
			"role":      user.Role,
		},
	})
}
