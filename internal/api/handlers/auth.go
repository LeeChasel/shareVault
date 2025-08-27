package handlers

import (
	"net/http"

	"github.com/LeeChasel/shareVault/internal/api/dto"
	"github.com/LeeChasel/shareVault/internal/models"
	"github.com/LeeChasel/shareVault/internal/service"
	"github.com/LeeChasel/shareVault/internal/utils"
	"github.com/gin-gonic/gin"
)

func Register(services *service.ApplicationServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "參數錯誤"})
			return
		}

		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密碼加密失敗"})
			return
		}
		user := models.User{
			Email:    req.Email,
			Username: req.Username,
			Password: hash,
		}

		isUserNameExist, err := services.UserService.ExistsByUsername(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if isUserNameExist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "此用戶名稱已被使用，請更換"})
			return
		}

		isEmailExist, err := services.UserService.ExistsByEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if isEmailExist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "此 email 已被使用，請更換"})
			return
		}

		if err := services.UserService.Create(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "註冊失敗"})
			return
		}

		resp := dto.RegisterResponse{
			ID:       user.ID.String(),
			Email:    user.Email,
			Username: user.Username,
		}
		c.JSON(http.StatusOK, resp)
	}
}

func Login(services *service.ApplicationServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "參數錯誤"})
			return
		}

		user, err := services.UserService.FindByEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "伺服器錯誤"})
			return
		}
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "此 email 尚未註冊，請先註冊"})
			return
		}
		if !utils.CheckPassword(user.Password, req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "密碼錯誤，請再試一次"})
			return
		}

		token, err := utils.GenerateJWT(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT 產生失敗"})
			return
		}

		resp := dto.LoginResponse{Token: token}
		c.JSON(http.StatusOK, resp)
	}
}
