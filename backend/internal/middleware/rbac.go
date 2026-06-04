package middleware

import (
	"cowork/internal/dto/response"
	"cowork/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// RequireRole returns a middleware that checks if the authenticated user's role
// (set by AuthMiddleware) is among the allowed roles. If not, it returns a
// forbidden error and aborts the request.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			response.Error(c, errcode.ErrForbidden, "permission denied")
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			response.Error(c, errcode.ErrForbidden, "permission denied")
			c.Abort()
			return
		}

		for _, allowed := range roles {
			if role == allowed {
				c.Next()
				return
			}
		}

		response.Error(c, errcode.ErrForbidden, "permission denied")
		c.Abort()
	}
}
