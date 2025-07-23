package handlers

import (
	"net/http"

	"github.com/LeeChasel/shareVault/backend/dto"
	"github.com/LeeChasel/shareVault/backend/models"
	"github.com/LeeChasel/shareVault/backend/repository"
	"github.com/LeeChasel/shareVault/backend/services"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "參數錯誤"})
		return
	}

	hash, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密碼加密失敗"})
		return
	}
	user := models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: hash,
	}

	userRepo := c.MustGet("repos").(*repository.Repositories).User
	// 檢查是否與現有使用者衝突
	if userRepo.ExistsByUsername(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "此用戶名稱已被使用，請更換"})
		return
	}
	if userRepo.ExistsByEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "此 email 已被使用，請更換"})
		return
	}

	if err := userRepo.Create(&user); err != nil {
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

func Login(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "參數錯誤"})
		return
	}
	// Default username and password
	if req.Username != "admin" || req.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤"})
		return
	}
	token, err := services.GenerateJWT(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT 產生失敗"})
		return
	}
	// Set JWT cookie
	c.SetCookie("jwt", token, 86400, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "login success"})
}

func Logout(c *gin.Context) {
	// Clear JWT cookie
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}
