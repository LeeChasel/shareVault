package middleware

import (
	"net/http"
	"strings"

	"github.com/LeeChasel/shareVault/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const AUTH_HEADER = "Authorization"

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AUTH_HEADER)
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未經授權，請先登入"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的 JWT"})
			c.Abort()
			return
		}

		userId, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "登入的 userId 格式錯誤"})
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
