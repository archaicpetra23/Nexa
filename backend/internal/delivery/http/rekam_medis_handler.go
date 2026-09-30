package http

import (
	"net/http"
	"strconv"
	"time"

	"nexa/backend/internal/model"
	"nexa/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type RekamMedisHandler struct {
	usecase usecase.RekamMedisUsecase
}

func NewRekamMedisHandler(uc usecase.RekamMedisUsecase) *RekamMedisHandler {
	return &RekamMedisHandler{usecase: uc}
}

type CreateRekamRequest struct {
	IDPasien        uint    `json:"id_pasien" binding:"required"`
	IDDokter        uint    `json:"id_dokter" binding:"required"`
	JenisPerawatan  string  `json:"jenis_perawatan" binding:"required,oneof=rawat_jalan rawat_inap"`
	Keluhan         string  `json:"keluhan" binding:"required"`
	Catatan         *string `json:"catatan"`
	TanggalPulang   *string `json:"tanggal_pulang"`
}

type AddDiagnosisRequest struct {
	KodeICD10 string `json:"kode_icd10" binding:"required"`
	Jenis     string `json:"jenis" binding:"required,oneof=primer sekunder"`
}

type AddTindakanRequest struct {
	KodeTindakan string `json:"kode_tindakan" binding:"required"`
	Jumlah       int    `json:"jumlah" binding:"required,gt=0"`
}

// ListRekam handles GET /api/v1/rekam-medis?pasien_id=
func (h *RekamMedisHandler) ListRekam(c *gin.Context) {
	pasienID, _ := strconv.ParseUint(c.Query("pasien_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	rekams, total, err := h.usecase.GetRekamMedisByPasien(uint(pasienID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memuat daftar rekam medis",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar rekam medis berhasil dimuat",
		"data": gin.H{
			"items": rekams,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetRekam handles GET /api/v1/rekam-medis/:id
func (h *RekamMedisHandler) GetRekam(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	rekam, err := h.usecase.GetRekamMedisByID(uint(id), userID, role)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Rekam medis tidak ditemukan atau tidak dapat diakses",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail rekam medis berhasil dimuat",
		"data":    rekam,
	})
}

// CreateRekam handles POST /api/v1/rekam-medis
func (h *RekamMedisHandler) CreateRekam(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateRekamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data rekam medis tidak valid",
			"data":    nil,
		})
		return
	}

	rekam := &model.RekamMedis{
		IDPasien:       req.IDPasien,
		IDDokter:       req.IDDokter,
		JenisPerawatan: req.JenisPerawatan,
		Keluhan:        req.Keluhan,
		Catatan:        req.Catatan,
	}

	if req.TanggalPulang != nil {
		var t time.Time
		if err := parseDate(*req.TanggalPulang, &t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format tanggal_pulang harus YYYY-MM-DD",
				"data":    nil,
			})
			return
		}
		rekam.TanggalPulang = &t
	}

	if err := h.usecase.CreateRekamMedis(rekam, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menyimpan rekam medis",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Rekam medis berhasil dibuat",
		"data":    rekam,
	})
}

// AddDiagnosis handles POST /api/v1/rekam-medis/:id/diagnosis
func (h *RekamMedisHandler) AddDiagnosis(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	rekamID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req AddDiagnosisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data diagnosis tidak valid",
			"data":    nil,
		})
		return
	}

	diagnosis := &model.RekamDiagnosis{
		KodeICD10: req.KodeICD10,
		Jenis:     req.Jenis,
	}

	if err := h.usecase.AddDiagnosis(uint(rekamID), diagnosis, userID, role); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Tidak dapat menambahkan diagnosis ke rekam medis ini",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Diagnosis berhasil ditambahkan",
		"data":    diagnosis,
	})
}

// AddTindakan handles POST /api/v1/rekam-medis/:id/tindakan
func (h *RekamMedisHandler) AddTindakan(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	rekamID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req AddTindakanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tindakan tidak valid. Jumlah harus lebih dari 0",
			"data":    nil,
		})
		return
	}

	tindakan := &model.DetailTindakan{
		KodeTindakan: req.KodeTindakan,
		Jumlah:       req.Jumlah,
	}

	if err := h.usecase.AddTindakan(uint(rekamID), tindakan, userID, role); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Tidak dapat menambahkan tindakan ke rekam medis ini",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Tindakan berhasil ditambahkan",
		"data":    tindakan,
	})
}