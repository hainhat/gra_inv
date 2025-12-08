package routes

import (
	"graduation_invitation/backend/controllers"
	"graduation_invitation/backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Public routes
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register)
		api.GET("/check-email", controllers.CheckEmail)
		api.POST("/rsvp", middleware.RSVPRateLimit(), controllers.SubmitRSVP)
		api.GET("/rsvp/stats", controllers.GetStats)
		api.GET("/rsvp/messages", controllers.GetRSVPMessages)
		api.POST("/refresh", controllers.RefreshToken)

		// Google Identity Services route
		api.POST("/auth/google/verify", controllers.VerifyGoogleToken)

		// Protected routes (require JWT)
		auth := api.Group("/")
		auth.Use(middleware.AuthJWT())
		{
			auth.GET("/me", controllers.Me)
			auth.POST("/logout", controllers.Logout)
		}

		// Admin routes (require JWT + admin role)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthJWT())
		admin.Use(middleware.RequireAdmin())
		{
			// Dashboard
			admin.GET("/dashboard", controllers.AdminGetDashboard)

			// User management
			admin.GET("/users", controllers.AdminGetUsers)
			admin.GET("/users/:id", controllers.AdminGetUser)
			admin.POST("/users", controllers.AdminCreateUser)
			admin.PUT("/users/:id", controllers.AdminUpdateUser)
			admin.DELETE("/users/:id", controllers.AdminDeleteUser)

			// RSVP management
			admin.GET("/rsvps", controllers.AdminGetRSVPs)
			admin.GET("/rsvps/:id", controllers.AdminGetRSVP)
			admin.PUT("/rsvps/:id", controllers.AdminUpdateRSVP)
			admin.DELETE("/rsvps/:id", controllers.AdminDeleteRSVP)

			// Public route - lấy setting
			api.GET("/settings/:key", controllers.GetSettingByKey)

			// Admin routes - quản lý settings
			admin.GET("/settings", controllers.AdminGetSettings)
			admin.PUT("/settings/:key", controllers.AdminUpdateSetting)
		}
	}
}
