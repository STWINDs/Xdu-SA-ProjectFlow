package middleware

import (
	"strings"

	"cowork/internal/config"
	"cowork/internal/dto/response"
	"cowork/pkg/errcode"
	"cowork/pkg/jwt"

	"github.com/gin-gonic/gin"
)

var jwtConfig *config.JWTConfig

// SetJWTConfig sets the JWT configuration for the auth middleware.
func SetJWTConfig(cfg *config.JWTConfig) {
	jwtConfig = cfg
}

// AuthMiddleware validates the JWT access token in the Authorization header.
// On success, it sets "userID", "username", and "role" in the Gin context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, errcode.ErrUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			response.Error(c, errcode.ErrUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		if jwtConfig == nil {
			response.Error(c, errcode.ErrInternal, "JWT configuration not initialized")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(tokenString, jwtConfig.AccessSecret)
		if err != nil {
			response.Error(c, errcode.ErrUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}
