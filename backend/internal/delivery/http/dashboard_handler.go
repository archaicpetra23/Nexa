package http

import (
	"net/http"

	"nexa/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// topDiagnosisRow is one entry of the dashboard "top_diagnosis" list.
type topDiagnosisRow struct {
	KodeICD10     string `json:"kode_icd10"`
	NamaDiagnosis string `json:"nama_diagnosis"`
	Jumlah        int64  `json:"jumlah"`
}

// dashboardStats is the /dashboard/stats payload.
type dashboardStats struct {
	TotalPasien           int64             `json:"total_pasien"`
	TotalKlaim            int64             `json:"total_klaim"`
	KlaimPending          int64             `json:"klaim_pending"`
	KlaimDisetujui        int64             `json:"klaim_disetujui"`
	KlaimDitolak          int64             `json:"klaim_ditolak"`
	KlaimDraft            int64             `json:"klaim_draft"`
	TotalNominalDisetujui float64           `json:"total_nominal_disetujui"`
	TopDiagnosis          []topDiagnosisRow `json:"top_diagnosis"`
}

// GetStats handles GET /api/v1/dashboard/stats
func (h *DashboardHandler) GetStats(c *gin.Context) {
	// Initialise the slice so an empty result serialises as [] and never null.
	// The frontend iterates this field directly.
	stats := dashboardStats{TopDiagnosis: []topDiagnosisRow{}}

	// Count total pasien (soft-deleted)
	if err := h.db.Model(&model.Pasien{}).
		Where("deleted_at IS NULL").
		Count(&stats.TotalPasien).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghitung total pasien",
			"data":    nil,
		})
		return
	}

	// Count total klaim (soft-deleted) and status breakdown
	if err := h.db.Model(&model.Klaim{}).
		Where("deleted_at IS NULL").
		Count(&stats.TotalKlaim).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghitung total klaim",
			"data":    nil,
		})
		return
	}

	// Count by status
	statuses := []struct {
		Status string
		Field  *int64
	}{
		{"pending", &stats.KlaimPending},
		{"disetujui", &stats.KlaimDisetujui},
		{"ditolak", &stats.KlaimDitolak},
		{"draft", &stats.KlaimDraft},
	}

	for _, s := range statuses {
		if err := h.db.Model(&model.Klaim{}).
			Where("deleted_at IS NULL AND status_klaim = ?", s.Status).
			Count(s.Field).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal menghitung klaim status " + s.Status,
				"data":    nil,
			})
			return
		}
	}

	// Sum nominal_klaim for disetujui
	if err := h.db.Model(&model.Klaim{}).
		Where("deleted_at IS NULL AND status_klaim = 'disetujui'").
		Select("COALESCE(SUM(nominal_klaim), 0)").
		Scan(&stats.TotalNominalDisetujui).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghitung total nominal disetujui",
			"data":    nil,
		})
		return
	}

	// Top 10 diagnosis by frequency in rekam_diagnosis
	if err := h.db.Model(&model.RekamDiagnosis{}).
		Select("d.kode_icd10, d.nama_diagnosis, COUNT(*) as jumlah").
		Joins("JOIN diagnosis d ON d.kode_icd10 = rekam_diagnosis.kode_icd10").
		Group("d.kode_icd10, d.nama_diagnosis").
		Order("jumlah DESC").
		Limit(10).
		Scan(&stats.TopDiagnosis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghitung top diagnosis",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Statistik dashboard berhasil dimuat",
		"data":    stats,
	})
}
