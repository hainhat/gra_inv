package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RSVPRateLimit creates a rate limiter middleware for RSVP endpoint
// Limit: 3 requests per minute per IP
func RSVPRateLimit() gin.HandlerFunc {
	// Define rate: 3 requests per minute
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  3,
	}

	// Create in-memory store
	store := memory.NewStore()

	// Create limiter instance
	instance := limiter.New(store, rate)

	// Create Gin middleware
	middleware := mgin.NewMiddleware(instance,
		mgin.WithErrorHandler(func(c *gin.Context, err error) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Spam ít thôi",
			})
			c.Abort()
		}),
		mgin.WithLimitReachedHandler(func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Spam ít thôi",
			})
			c.Abort()
		}),
	)

	return middleware
}
