package handlers

import (
	"net/http"

	"github.com/LeeChasel/sharevault/backend/models"
	"github.com/LeeChasel/sharevault/backend/services"
	"github.com/gin-gonic/gin"
)

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