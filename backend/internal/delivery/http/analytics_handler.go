package http

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"nexa/backend/internal/model"
	"nexa/backend/internal/repository"
	"nexa/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalyticsHandler struct {
	db           *gorm.DB
	klaimRepo    repository.KlaimRepository
	klaimUsecase usecase.KlaimUsecase
}

func NewAnalyticsHandler(db *gorm.DB) *AnalyticsHandler {
	klaimRepo := repository.NewKlaimRepository(db)
	klaimUC := usecase.NewKlaimUsecase(klaimRepo)
	return &AnalyticsHandler{
		db:           db,
		klaimRepo:    klaimRepo,
		klaimUsecase: klaimUC,
	}
}

func (h *AnalyticsHandler) ExportKlaim(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	status := c.Query("status")
	from := c.Query("from")
	to := c.Query("to")
	unitID := c.Query("unit")

	query := h.db.Preload("RekamMedis.Pasien").
		Preload("PetugasCasemix").
		Where("deleted_at IS NULL")

	if status != "" {
		query = query.Where("status_klaim = ?", status)
	}
	if from != "" {
		query = query.Where("tanggal_klaim >= ?", from)
	}
	if to != "" {
		query = query.Where("tanggal_klaim <= ?", to)
	}
	if unitID != "" {
		query = query.Joins("JOIN rekam_medis ON klaim.id_rekam = rekam_medis.id_rekam").
			Where("rekam_medis.id_unit = ?", unitID)
	}

	var klaims []model.Klaim
	if err := query.Order("created_at DESC").Find(&klaims).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memuat data klaim",
			"data":    nil,
		})
		return
	}

	if format == "json" {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    klaims,
		})
		return
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{
		"ID Klaim",
		"Status",
		"Nominal",
		"Nama Pasien",
		"Tanggal Klaim",
		"Kode CBGS",
		"Petugas",
	}

	if err := writer.Write(header); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat CSV",
			"data":    nil,
		})
		return
	}

	for _, k := range klaims {
		pasienNama := ""
		if k.RekamMedis != nil && k.RekamMedis.Pasien != nil {
			pasienNama = k.RekamMedis.Pasien.Nama
		}
		petugasNama := ""
		if k.PetugasCasemix != nil {
			petugasNama = k.PetugasCasemix.Nama
		}
		tanggalKlaim := k.TanggalKlaim.Format("2006-01-02")

		record := []string{
			strconv.FormatUint(uint64(k.IDKlaim), 10),
			k.StatusKlaim,
			fmt.Sprintf("%.2f", k.NominalKlaim),
			pasienNama,
			tanggalKlaim,
			k.KodeCBGS,
			petugasNama,
		}

		if err := writer.Write(record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal membuat CSV",
				"data":    nil,
			})
			return
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat CSV",
			"data":    nil,
		})
		return
	}

	filename := fmt.Sprintf("klaim-export-%s.csv", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv", buf.Bytes())
}
