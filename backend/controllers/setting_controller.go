package controllers

import (
	"graduation_invitation/backend/config"
	"graduation_invitation/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PUBLIC API
// GET /api/settings/:key
func GetSettingByKey(c *gin.Context) {
	key := c.Param("key")
	var setting models.Setting
	if err := config.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Setting not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    setting,
	})
}

// ADMIN APIs
// GET /api/admin/settings
func AdminGetSettings(c *gin.Context) {
	var settings []models.Setting
	if err := config.DB.Order("key asc").Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch settings",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}

// PUT /api/admin/settings/:key
func AdminUpdateSetting(c *gin.Context) {
	key := c.Param("key")

	var req struct {
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
		})
		return
	}

	var setting models.Setting
	if err := config.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Setting not found",
		})
		return
	}

	setting.Value = req.Value

	if err := config.DB.Save(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update setting",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Setting updated successfully",
		"data":    setting,
	})
}
