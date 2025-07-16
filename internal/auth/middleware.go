package auth

import (
	"net/http"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"strings"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authRepo *repository.AuthRepository
}

func NewAuthMiddleware(authRepo *repository.AuthRepository) *AuthMiddleware {
	return &AuthMiddleware{authRepo: authRepo}
}

func (m *AuthMiddleware) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// EXTRACT TOKEN FROM HEADER
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}
		tokenString := parts[1]

		// VALIDATE TOKEN
		claims, err := ValidateToken(tokenString, false)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// CHECK BLACKLIST TOKEN
		isBlacklisted, err := m.authRepo.IsTokenBlacklisted(c.Request.Context(), claims.TokenID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not verify token session"})
			return
		}
		if isBlacklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked (logged out)"})
			return
		}

		c.Set("user_npk", claims.NPK)
		c.Set("user_position_id", claims.PositionID)

		c.Next()
	}
}