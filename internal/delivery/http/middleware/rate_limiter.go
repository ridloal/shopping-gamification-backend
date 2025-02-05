package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(rdb *redis.Client, rateLimit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := "rate_limiter:" + clientIP

		// Increment the request count for the client IP
		count, err := rdb.Incr(context.Background(), key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// Set the expiration for the key if it's the first request
		if count == 1 {
			rdb.Expire(context.Background(), key, duration)
		}

		// Check if the request count exceeds the rate limit
		if count > int64(rateLimit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
