package http

import (
	"net/http"
	"strconv"
	"time"

	"nexa/backend/internal/model"
	"nexa/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

type KlaimAdminResponse struct {
	IDKlaim      uint    `json:"id_klaim"`
	IDRekam      uint    `json:"id_rekam"`
	KodeCBGS     string  `json:"kode_cbgs"`
	StatusKlaim  string  `json:"status_klaim"`
	NominalKlaim float64 `json:"nominal_klaim"`
	PasienNama   string  `json:"pasien_nama"`
	PetugasNama  string  `json:"petugas_nama"`
	CreatedAt    string  `json:"created_at"`
}

type AuditTrailResponse struct {
	IDLog          uint      `json:"id_log"`
	IDUser         *uint     `json:"id_user"`
	Username       string    `json:"username"`
	TabelTerdampak string    `json:"tabel_terdampak"`
	Aktivitas      string    `json:"aktivitas"`
	DataID         *int      `json:"data_id"`
	Waktu          time.Time `json:"waktu"`
}

// GetUsers handles GET /api/v1/admin/users
func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	roleFilter := c.Query("role")

	query := h.db.Preload("Role").Preload("Unit").Where("deleted_at IS NULL")

	if search != "" {
		query = query.Where("nama ILIKE ? OR username ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if roleFilter != "" {
		query = query.Joins("JOIN roles ON users.id_role = roles.id_role").
			Where("roles.nama_role = ?", roleFilter)
	}

	var total int64
	query.Model(&models.User{}).Count(&total)

	var users []models.User
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memuat data pengguna",
			"data":    nil,
		})
		return
	}

	response := make([]UserResponse, len(users))
	for i, u := range users {
		response[i] = toUserResponse(h.db, &u)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar pengguna berhasil dimuat",
		"data": gin.H{
			"items": response,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetSubmissions handles GET /api/v1/admin/submissions
func (h *AdminHandler) GetSubmissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	dari := c.Query("dari")
	sampai := c.Query("sampai")
	search := c.Query("search")

	query := h.db.Preload("RekamMedis.Pasien").
		Preload("PetugasCasemix").
		Where("deleted_at IS NULL")

	if status != "" {
		query = query.Where("status_klaim = ?", status)
	}
	if dari != "" {
		query = query.Where("tanggal_klaim >= ?", dari)
	}
	if sampai != "" {
		query = query.Where("tanggal_klaim <= ?", sampai)
	}
	if search != "" {
		query = query.Joins("JOIN rekam_medis ON klaim.id_rekam = rekam_medis.id_rekam").
			Joins("JOIN pasien ON rekam_medis.id_pasien = pasien.id_pasien").
			Where("pasien.nama ILIKE ?", "%"+search+"%")
	}

	var total int64
	query.Model(&model.Klaim{}).Count(&total)

	var klaims []model.Klaim
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&klaims).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memuat data klaim",
			"data":    nil,
		})
		return
	}

	response := make([]KlaimAdminResponse, len(klaims))
	for i, k := range klaims {
		pasienNama := ""
		if k.RekamMedis != nil && k.RekamMedis.Pasien != nil {
			pasienNama = k.RekamMedis.Pasien.Nama
		}
		petugasNama := ""
		if k.PetugasCasemix != nil {
			petugasNama = k.PetugasCasemix.Nama
		}
		response[i] = KlaimAdminResponse{
			IDKlaim:      k.IDKlaim,
			IDRekam:      k.IDRekam,
			KodeCBGS:     k.KodeCBGS,
			StatusKlaim:  k.StatusKlaim,
			NominalKlaim: k.NominalKlaim,
			PasienNama:   pasienNama,
			PetugasNama:  petugasNama,
			CreatedAt:    k.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar klaim berhasil dimuat",
		"data": gin.H{
			"items": response,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetAuditTrail handles GET /api/v1/admin/audit-trail
func (h *AdminHandler) GetAuditTrail(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	userID := c.Query("user_id")
	aktivitas := c.Query("aktivitas")
	dari := c.Query("dari")
	sampai := c.Query("sampai")

	query := h.db.Preload("User")

	if userID != "" {
		query = query.Where("id_user = ?", userID)
	}
	if aktivitas != "" {
		query = query.Where("aktivitas = ?", aktivitas)
	}
	if dari != "" {
		query = query.Where("waktu >= ?", dari)
	}
	if sampai != "" {
		endDate, err := time.Parse("2006-01-02", sampai)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format tanggal sampai tidak valid",
				"data":    nil,
			})
			return
		}
		query = query.Where("waktu < ?", endDate.AddDate(0, 0, 1))
	}

	var total int64
	query.Model(&model.LogAktivitas{}).Count(&total)

	var logs []model.LogAktivitas
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("waktu DESC").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memuat audit trail",
			"data":    nil,
		})
		return
	}

	response := make([]AuditTrailResponse, len(logs))
	for i, l := range logs {
		username := ""
		if l.User != nil {
			username = l.User.Nama
		}
		response[i] = AuditTrailResponse{
			IDLog:          l.IDLog,
			IDUser:         l.IDUser,
			Username:       username,
			TabelTerdampak: l.TabelTerdampak,
			Aktivitas:      l.Aktivitas,
			DataID:         l.DataID,
			Waktu:          l.Waktu,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Audit trail berhasil dimuat",
		"data": gin.H{
			"items": response,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
