package http

import (
	"net/http"
	"strconv"

	"nexa/backend/internal/model"
	"nexa/backend/internal/models"
	"nexa/backend/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MasterHandler struct {
	db   *gorm.DB
	repo repository.MasterRepository
}

func NewMasterHandler(db *gorm.DB) *MasterHandler {
	return &MasterHandler{db: db, repo: repository.NewMasterRepository(db)}
}

type CreateUnitRequest struct {
	NamaUnit string `json:"nama_unit" binding:"required"`
}

type UpdateUnitRequest struct {
	NamaUnit string `json:"nama_unit" binding:"required"`
}

type CreateDiagnosisRequest struct {
	KodeICD10     string  `json:"kode_icd10" binding:"required"`
	NamaDiagnosis string  `json:"nama_diagnosis" binding:"required"`
	Kategori      *string `json:"kategori"`
}

type UpdateDiagnosisRequest struct {
	NamaDiagnosis string  `json:"nama_diagnosis" binding:"required"`
	Kategori      *string `json:"kategori"`
}

type CreateTindakanRequest struct {
	KodeTindakan string  `json:"kode_tindakan" binding:"required"`
	NamaTindakan string  `json:"nama_tindakan" binding:"required"`
	TarifStandar float64 `json:"tarif_standar" binding:"gte=0"`
}

type UpdateTindakanRequest struct {
	NamaTindakan string  `json:"nama_tindakan" binding:"required"`
	TarifStandar float64 `json:"tarif_standar" binding:"gte=0"`
}

type CreateCBGSRequest struct {
	KodeCBGS  string  `json:"kode_cbgs" binding:"required"`
	Deskripsi string  `json:"deskripsi" binding:"required"`
	Tarif     float64 `json:"tarif" binding:"gte=0"`
}

type UpdateCBGSRequest struct {
	Deskripsi string  `json:"deskripsi" binding:"required"`
	Tarif     float64 `json:"tarif" binding:"gte=0"`
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

// ICD-10 Handlers

// ListICD10 handles GET /api/v1/master/icd10
func (h *MasterHandler) ListICD10(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	diagnoses, total, err := h.repo.ListDiagnosis(search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memuat daftar ICD-10",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar ICD-10 berhasil dimuat",
		"data": gin.H{
			"items": diagnoses,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// SearchICD10 handles GET /api/v1/master/icd10/search
func (h *MasterHandler) SearchICD10(c *gin.Context) {
	q := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	diagnoses, err := h.repo.SearchDiagnosisByCodeOrName(q, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mencari ICD-10",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pencarian ICD-10 berhasil",
		"data":    diagnoses,
	})
}

// GetICD10 handles GET /api/v1/master/icd10/:kode
func (h *MasterHandler) GetICD10(c *gin.Context) {
	kode := c.Param("kode")
	diagnosis, err := h.repo.GetDiagnosisByCode(kode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Kode ICD-10 tidak ditemukan",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail ICD-10 berhasil dimuat",
		"data":    diagnosis,
	})
}

// CreateICD10 handles POST /api/v1/master/icd10
func (h *MasterHandler) CreateICD10(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateDiagnosisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data ICD-10 tidak valid",
			"data":    nil,
		})
		return
	}

	diagnosis := &model.Diagnosis{
		KodeICD10:     req.KodeICD10,
		NamaDiagnosis: req.NamaDiagnosis,
		Kategori:      req.Kategori,
	}

	if err := h.repo.CreateDiagnosis(diagnosis); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Kode sudah terdaftar",
			"data":    nil,
		})
		return
	}

	h.logActivity(userID, "diagnosis", "INSERT", nil)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "ICD-10 berhasil dibuat",
		"data":    diagnosis,
	})
}

// UpdateICD10 handles PUT /api/v1/master/icd10/:kode
func (h *MasterHandler) UpdateICD10(c *gin.Context) {
	userID := c.GetUint("user_id")
	kode := c.Param("kode")

	var req UpdateDiagnosisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data ICD-10 tidak valid",
			"data":    nil,
		})
		return
	}

	existing, err := h.repo.GetDiagnosisByCode(kode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Kode ICD-10 tidak ditemukan",
			"data":    nil,
		})
		return
	}

	existing.NamaDiagnosis = req.NamaDiagnosis
	existing.Kategori = req.Kategori

	if err := h.repo.UpdateDiagnosis(existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui ICD-10",
			"data":    nil,
		})
		return
	}

	h.logActivity(userID, "diagnosis", "UPDATE", nil)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "ICD-10 berhasil diperbarui",
		"data":    existing,
	})
}

// DeleteICD10 handles DELETE /api/v1/master/icd10/:kode
func (h *MasterHandler) DeleteICD10(c *gin.Context) {
	userID := c.GetUint("user_id")
	kode := c.Param("kode")

	_, err := h.repo.GetDiagnosisByCode(kode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Kode ICD-10 tidak ditemukan",
			"data":    nil,
		})
		return
	}

	if err := h.repo.DeleteDiagnosis(kode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus ICD-10",
			"data":    nil,
		})
		return
	}

	h.logActivity(userID, "diagnosis", "DELETE", nil)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "ICD-10 berhasil dihapus",
		"data":    nil,
	})
}

// ICD-9 CM (Tindakan) Handlers

// ListICD9 handles GET /api/v1/master/icd9
func (h *MasterHandler) ListICD9(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	tindakans, total, err := h.repo.ListTindakan(search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memuat daftar ICD-9 CM",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar ICD-9 CM berhasil dimuat",
		"data": gin.H{
			"items": tindakans,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// SearchICD9 handles GET /api/v1/master/icd9/search
func (h *MasterHandler) SearchICD9(c *gin.Context) {
	q := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	tindakans, err := h.repo.SearchTindakanByCodeOrName(q, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mencari ICD-9 CM",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pencarian ICD-9 CM berhasil",
		"data":    tindakans,
	})
}

// GetICD9 handles GET /api/v1/master/icd9/:kode
func (h *MasterHandler) GetICD9(c *gin.Context) {
	kode := c.Param("kode")
	tindakan, err := h.repo.GetTindakanByCode(kode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Kode tindakan tidak ditemukan",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail tindakan berhasil dimuat",
		"data":    tindakan,
	})
}

// CreateICD9 handles POST /api/v1/master/icd9
func (h *MasterHandler) CreateICD9(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateTindakanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tindakan tidak valid",
			"data":    nil,
		})
		return
	}

	tindakan := &model.Tindakan{
		KodeTindakan: req.KodeTindakan,
		NamaTindakan: req.NamaTindakan,
		TarifStandar: req.TarifStandar,
	}

	if err := h.repo.CreateTindakan(tindakan); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Kode sudah terdaftar",
			"data":    nil,
		})
		return
	}

	h.logActivity(userID, "tindakan", "INSERT", nil)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Tindakan berhasil dibuat",
		"data":    tindakan,
	})
}

// UpdateICD9 handles PUT /api/v1/master/icd9/:kode
func (h *MasterHandler) UpdateICD9(c *gin.Context) {
	userID := c.GetUint("user_id")
	kode := c.Param("kode")

	var req UpdateTindakanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tindakan tidak valid",
			"data":    nil,
		})
		return
	}

	existing, err := h.repo.GetTindakanByCode(kode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Kode tindakan tidak ditemukan",
			"data":    nil,
		})
		return
	}

	existing.NamaTindakan = req.NamaTindakan
	existing.TarifStandar = req.TarifStandar

	if err := h.repo.UpdateTindakan(existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui tindakan",
			"data":    nil,
		})
		return
	}

	h.logActivity(userID, "tindakan", "UPDATE", nil)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tindakan berhasil diperbarui",
		"data":    existing,
	})
}

// DeleteICD9 handles DELETE /api/v1/master/icd9/:kode
func (h *MasterHandler) DeleteICD9(c *gin.Context) {
	userID := c.GetUint("user_id")
	kode := c.Param("kode")

	_, err := h.repo.GetTindakanByCode(kode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Kode tindakan tidak ditemukan",
			"data":    nil,
		})
		return
	}

	if err := h.repo.DeleteTindakan(kode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus tindakan",
			"data":    nil,
		})
		return
	}

	h.logActivity(userID, "tindakan", "DELETE", nil)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tindakan berhasil dihapus",
		"data":    nil,
	})
}

// CBGs Handlers

// ListCBGS handles GET /api/v1/master/cbgs
func (h *MasterHandler) ListCBGS(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	cbgsList, total, err := h.repo.ListCBGS(search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memuat daftar CBGs",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar CBGs berhasil dimuat",
		"data": gin.H{
			"items": cbgsList,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetCBGS handles GET /api/v1/master/cbgs/:kode
func (h *MasterHandler) GetCBGS(c *gin.Context) {
	kode := c.Param("kode")
	cbgs, err := h.repo.GetCBGSByCode(kode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Kode CBGs tidak ditemukan",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail CBGs berhasil dimuat",
		"data":    cbgs,
	})
}

// CreateCBGS handles POST /api/v1/master/cbgs
func (h *MasterHandler) CreateCBGS(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateCBGSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data CBGs tidak valid",
			"data":    nil,
		})
		return
	}

	cbgs := &model.TarifCBGs{
		KodeCBGS:  req.KodeCBGS,
		Deskripsi: req.Deskripsi,
		Tarif:     req.Tarif,
	}

	if err := h.repo.CreateCBGS(cbgs); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Kode sudah terdaftar",
			"data":    nil,
		})
		return
	}

	h.logActivity(userID, "tarif_cbgs", "INSERT", nil)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "CBGs berhasil dibuat",
		"data":    cbgs,
	})
}

// UpdateCBGS handles PUT /api/v1/master/cbgs/:kode
func (h *MasterHandler) UpdateCBGS(c *gin.Context) {
	userID := c.GetUint("user_id")
	kode := c.Param("kode")

	var req UpdateCBGSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data CBGs tidak valid",
			"data":    nil,
		})
		return
	}

	existing, err := h.repo.GetCBGSByCode(kode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Kode CBGs tidak ditemukan",
			"data":    nil,
		})
		return
	}

	existing.Deskripsi = req.Deskripsi
	existing.Tarif = req.Tarif

	if err := h.repo.UpdateCBGS(existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui CBGs",
			"data":    nil,
		})
		return
	}

	h.logActivity(userID, "tarif_cbgs", "UPDATE", nil)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "CBGs berhasil diperbarui",
		"data":    existing,
	})
}

// DeleteCBGS handles DELETE /api/v1/master/cbgs/:kode
func (h *MasterHandler) DeleteCBGS(c *gin.Context) {
	userID := c.GetUint("user_id")
	kode := c.Param("kode")

	_, err := h.repo.GetCBGSByCode(kode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Kode CBGs tidak ditemukan",
			"data":    nil,
		})
		return
	}

	if err := h.repo.DeleteCBGS(kode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus CBGs",
			"data":    nil,
		})
		return
	}

	h.logActivity(userID, "tarif_cbgs", "DELETE", nil)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "CBGs berhasil dihapus",
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
