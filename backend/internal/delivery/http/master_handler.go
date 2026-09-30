package http

import (
	"net/http"
	"strconv"

	"nexa/backend/internal/model"
	"nexa/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MasterHandler struct {
	db *gorm.DB
}

func NewMasterHandler(db *gorm.DB) *MasterHandler {
	return &MasterHandler{db: db}
}

type CreateUnitRequest struct {
	NamaUnit string `json:"nama_unit" binding:"required"`
}

type UpdateUnitRequest struct {
	NamaUnit string `json:"nama_unit" binding:"required"`
}

// GetRoles handles GET /api/v1/master/roles
func (h *MasterHandler) GetRoles(c *gin.Context) {
	var roles []models.Role
	if err := h.db.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memuat data role",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar role berhasil dimuat",
		"data":    roles,
	})
}

// GetUnits handles GET /api/v1/master/units
func (h *MasterHandler) GetUnits(c *gin.Context) {
	var units []models.Unit
	if err := h.db.Order("nama_unit ASC").Find(&units).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memuat data unit",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar unit berhasil dimuat",
		"data":    units,
	})
}

// CreateUnit handles POST /api/v1/master/units
func (h *MasterHandler) CreateUnit(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Nama unit wajib diisi",
			"data":    nil,
		})
		return
	}

	// Check unique
	var existing models.Unit
	if err := h.db.Where("nama_unit = ?", req.NamaUnit).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Nama unit sudah digunakan",
			"data":    nil,
		})
		return
	}

	unit := &models.Unit{
		NamaUnit: req.NamaUnit,
	}

	if err := h.db.Create(unit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat unit",
			"data":    nil,
		})
		return
	}

	// Log activity
	h.logActivity(userID, "units", "INSERT", &unit.IDUnit)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Unit berhasil dibuat",
		"data":    unit,
	})
}

// UpdateUnit handles PUT /api/v1/master/units/:id
func (h *MasterHandler) UpdateUnit(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID unit tidak valid",
			"data":    nil,
		})
		return
	}

	var req UpdateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Nama unit wajib diisi",
			"data":    nil,
		})
		return
	}

	var unit models.Unit
	if err := h.db.First(&unit, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Unit tidak ditemukan",
			"data":    nil,
		})
		return
	}

	// Check unique
	var existing models.Unit
	if err := h.db.Where("nama_unit = ? AND id_unit != ?", req.NamaUnit, id).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Nama unit sudah digunakan",
			"data":    nil,
		})
		return
	}

	unit.NamaUnit = req.NamaUnit
	if err := h.db.Save(&unit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui unit",
			"data":    nil,
		})
		return
	}

	// Log activity
	h.logActivity(userID, "units", "UPDATE", &unit.IDUnit)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Unit berhasil diperbarui",
		"data":    unit,
	})
}

// DeleteUnit handles DELETE /api/v1/master/units/:id
func (h *MasterHandler) DeleteUnit(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID unit tidak valid",
			"data":    nil,
		})
		return
	}

	var unit models.Unit
	if err := h.db.First(&unit, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Unit tidak ditemukan",
			"data":    nil,
		})
		return
	}

	// Hard delete for master data (units)
	if err := h.db.Delete(&unit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus unit",
			"data":    nil,
		})
		return
	}

	// Log activity
	h.logActivity(userID, "units", "DELETE", &unit.IDUnit)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Unit berhasil dihapus",
		"data":    nil,
	})
}

func (h *MasterHandler) logActivity(userID uint, tabel string, aktivitas string, dataID *uint) {
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