package http

import (
	"net/http"
	"strconv"
	"time"

	"nexa/backend/internal/model"
	"nexa/backend/internal/models"
	"nexa/backend/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminUserHandler struct {
	db *gorm.DB
}

func NewAdminUserHandler(db *gorm.DB) *AdminUserHandler {
	return &AdminUserHandler{db: db}
}

type CreateUserRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required,min=8"`
	Nama         string `json:"nama" binding:"required"`
	Profesi      string `json:"profesi" binding:"required"`
	Spesialisasi string `json:"spesialisasi"`
	NoSTR        string `json:"no_str"`
	IDRole       uint   `json:"id_role" binding:"required"`
	IDUnit       uint   `json:"id_unit" binding:"required"`
}

type UpdateUserRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Nama         string `json:"nama"`
	Profesi      string `json:"profesi"`
	Spesialisasi string `json:"spesialisasi"`
	NoSTR        string `json:"no_str"`
	IDRole       uint   `json:"id_role"`
	IDUnit       uint   `json:"id_unit"`
}

func (h *AdminUserHandler) CreateUser(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data pengguna tidak valid. Username, password (min 8), nama, profesi, id_role, id_unit wajib diisi.",
			"data":    nil,
		})
		return
	}

	// Check username unique
	var existing models.User
	if err := h.db.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Username sudah digunakan",
			"data":    nil,
		})
		return
	}

	// Check role exists
	var role models.Role
	if err := h.db.First(&role, req.IDRole).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Role tidak ditemukan",
			"data":    nil,
		})
		return
	}

	// Check unit exists
	var unit models.Unit
	if err := h.db.First(&unit, req.IDUnit).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unit tidak ditemukan",
			"data":    nil,
		})
		return
	}

	// Hash password
	hashedPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memproses password",
			"data":    nil,
		})
		return
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Nama:         req.Nama,
		Profesi:      req.Profesi,
		Spesialisasi: &req.Spesialisasi,
		NoSTR:        &req.NoSTR,
		IDRole:       req.IDRole,
		IDUnit:       req.IDUnit,
	}

	if err := h.db.Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat pengguna",
			"data":    nil,
		})
		return
	}

	// Log activity
	h.logActivity(userID, "users", "INSERT", &user.IDUser)

	// Preload relations for response
	h.db.Preload("Role").Preload("Unit").First(user, user.IDUser)

	response := toUserResponse(h.db, user)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pengguna berhasil dibuat",
		"data":    response,
	})
}

func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID pengguna tidak valid",
			"data":    nil,
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"data":    nil,
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Pengguna tidak ditemukan",
			"data":    nil,
		})
		return
	}

	// Check username unique if changed
	if req.Username != "" && req.Username != user.Username {
		var existing models.User
		if err := h.db.Where("username = ? AND id_user != ?", req.Username, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Username sudah digunakan",
				"data":    nil,
			})
			return
		}
		user.Username = req.Username
	}

	if req.Nama != "" {
		user.Nama = req.Nama
	}
	if req.Profesi != "" {
		user.Profesi = req.Profesi
	}
	if req.Spesialisasi != "" {
		user.Spesialisasi = &req.Spesialisasi
	}
	if req.NoSTR != "" {
		user.NoSTR = &req.NoSTR
	}
	if req.IDRole > 0 {
		var role models.Role
		if err := h.db.First(&role, req.IDRole).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Role tidak ditemukan",
				"data":    nil,
			})
			return
		}
		user.IDRole = req.IDRole
	}
	if req.IDUnit > 0 {
		var unit models.Unit
		if err := h.db.First(&unit, req.IDUnit).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Unit tidak ditemukan",
				"data":    nil,
			})
			return
		}
		user.IDUnit = req.IDUnit
	}
	if req.Password != "" {
		hashedPassword, err := pkg.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal memproses password",
				"data":    nil,
			})
			return
		}
		user.PasswordHash = hashedPassword
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui pengguna",
			"data":    nil,
		})
		return
	}

	// Log activity
	h.logActivity(userID, "users", "UPDATE", &user.IDUser)

	h.db.Preload("Role").Preload("Unit").First(&user, user.IDUser)
	response := toUserResponse(h.db, &user)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pengguna berhasil diperbarui",
		"data":    response,
	})
}

func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID pengguna tidak valid",
			"data":    nil,
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Pengguna tidak ditemukan",
			"data":    nil,
		})
		return
	}

	// Soft delete
	now := time.Now()
	user.DeletedAt = &now
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus pengguna",
			"data":    nil,
		})
		return
	}

	// Log activity
	h.logActivity(userID, "users", "DELETE", &user.IDUser)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pengguna berhasil dihapus",
		"data":    nil,
	})
}

func (h *AdminUserHandler) logActivity(userID uint, tabel string, aktivitas string, dataID *uint) {
	var idInt *int
	if dataID != nil {
		v := int(*dataID)
		idInt = &v
	}
	h.db.Create(&model.LogAktivitas{
		IDUser:         &userID,
		TabelTerdampak: tabel,
		Aktivitas:      aktivitas,
		DataID:         idInt,
	})
}
