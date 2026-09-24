package handlers

import (
	"net/http"
	"nexa/backend/internal/models"
	"nexa/backend/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB        *gorm.DB
	JWTSecret string
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	var user models.User
	if err := h.DB.Preload("Role").Where("username = ? AND deleted_at IS NULL", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
		return
	}

	if !pkg.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
		return
	}

	token, err := pkg.GenerateToken(user.IDUser, user.Role.NamaRole, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Token generation failed"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt", token, 3600*8, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data": gin.H{
			"user": gin.H{
				"id_user":  user.IDUser,
				"nama":     user.Nama,
				"username": user.Username,
				"profesi":  user.Profesi,
				"role":     user.Role.NamaRole,
			},
		},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.DB.Preload("Role").Where("id_user = ? AND deleted_at IS NULL", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id_user":  user.IDUser,
			"nama":     user.Nama,
			"username": user.Username,
			"profesi":  user.Profesi,
			"role":     user.Role.NamaRole,
		},
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "No refresh token"})
		return
	}

	claims, err := pkg.ValidateToken(cookie, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid refresh token"})
		return
	}

	token, err := pkg.GenerateToken(claims.UserID, claims.Role, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Token generation failed"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt", token, 3600*8, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token refreshed",
		"data": gin.H{
			"user": gin.H{
				"user_id": claims.UserID,
				"role":    claims.Role,
			},
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logout successful"})
}
