package middleware

import (
	"github.com/LeeChasel/shareVault/backend/repository"
	"github.com/gin-gonic/gin"
)

// 將 repositories 注入 Gin context
func InjectRepos(repos *repository.Repositories) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("repos", repos)
		c.Next()
	}
}
