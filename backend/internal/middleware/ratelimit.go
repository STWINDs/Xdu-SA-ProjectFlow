package middleware

import (
	"context"
	"fmt"
	"time"

	"cowork/internal/dto/response"
	"cowork/pkg/errcode"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitMiddleware returns a Gin handler that enforces IP-based rate limiting
// using Redis INCR + EXPIRE. If the limit is exceeded within the window,
// it responds with an error and aborts the request.
func RateLimitMiddleware(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rdb == nil {
			// Redis unavailable; skip rate limiting.
			c.Next()
			return
		}

		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s", ip)

		count, err := rdb.Incr(context.Background(), key).Result()
		if err != nil {
			// If Redis fails, allow the request through.
			c.Next()
			return
		}

		// Set expiration on first request in the window.
		if count == 1 {
			rdb.Expire(context.Background(), key, window)
		}

		if count > int64(limit) {
			response.Error(c, errcode.ErrLoginLocked, "too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}
