package http

import (
	"errors"
	"net/http"
	"strconv"

	"nexa/backend/internal/model"
	"nexa/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type PasienHandler struct {
	usecase usecase.PasienUsecase
}

func NewPasienHandler(uc usecase.PasienUsecase) *PasienHandler {
	return &PasienHandler{usecase: uc}
}

type CreatePasienRequest struct {
	NIK          string  `json:"nik" binding:"required,len=16,numeric"`
	NoBPJS       *string `json:"no_bpjs" binding:"omitempty,len=13,numeric"`
	Nama         string  `json:"nama" binding:"required"`
	TanggalLahir string  `json:"tanggal_lahir" binding:"required"`
	JenisKelamin string  `json:"jenis_kelamin" binding:"required,oneof=L P"`
	Alamat       string  `json:"alamat" binding:"required"`
	NoHP         *string `json:"no_hp"`
}

// CreatePasien handles POST /api/v1/pasien
func (h *PasienHandler) CreatePasien(c *gin.Context) {
	var req CreatePasienRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data pasien tidak valid. Pastikan NIK 16 digit, tanggal_lahir format YYYY-MM-DD, jenis_kelamin L atau P.",
			"data":    nil,
		})
		return
	}

	pasien := &model.Pasien{
		NIK:          req.NIK,
		NoBPJS:       req.NoBPJS,
		Nama:         req.Nama,
		JenisKelamin: req.JenisKelamin,
		Alamat:       req.Alamat,
		NoHP:         req.NoHP,
	}

	if err := parseDate(req.TanggalLahir, &pasien.TanggalLahir); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format tanggal_lahir harus YYYY-MM-DD.",
			"data":    nil,
		})
		return
	}

	if err := h.usecase.CreatePasien(pasien); err != nil {
		handlePasienError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pasien berhasil didaftarkan",
		"data":    pasien,
	})
}

// GetPasien handles GET /api/v1/pasien/:id
func (h *PasienHandler) GetPasien(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID pasien tidak valid",
			"data":    nil,
		})
		return
	}

	pasien, err := h.usecase.GetPasienByID(uint(id))
	if err != nil {
		handlePasienError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pasien ditemukan",
		"data":    pasien,
	})
}

// SearchPasien handles GET /api/v1/pasien?search=&page=1&limit=10
func (h *PasienHandler) SearchPasien(c *gin.Context) {
	query := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	pasiens, total, err := h.usecase.SearchPasien(query, page, limit)
	if err != nil {
		handlePasienError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar pasien berhasil dimuat",
		"data": gin.H{
			"items": pasiens,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// UpdatePasien handles PUT /api/v1/pasien/:id
func (h *PasienHandler) UpdatePasien(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID pasien tidak valid",
			"data":    nil,
		})
		return
	}

	existing, err := h.usecase.GetPasienByID(uint(id))
	if err != nil {
		handlePasienError(c, err)
		return
	}

	var req CreatePasienRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data pasien tidak valid",
			"data":    nil,
		})
		return
	}

	existing.NIK = req.NIK
	existing.NoBPJS = req.NoBPJS
	existing.Nama = req.Nama
	existing.JenisKelamin = req.JenisKelamin
	existing.Alamat = req.Alamat
	existing.NoHP = req.NoHP

	if err := parseDate(req.TanggalLahir, &existing.TanggalLahir); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format tanggal_lahir harus YYYY-MM-DD",
			"data":    nil,
		})
		return
	}

	if err := h.usecase.UpdatePasien(existing); err != nil {
		handlePasienError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pasien berhasil diperbarui",
		"data":    existing,
	})
}

// DeletePasien handles DELETE /api/v1/pasien/:id (soft delete)
func (h *PasienHandler) DeletePasien(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID pasien tidak valid",
			"data":    nil,
		})
		return
	}

	if err := h.usecase.DeletePasien(uint(id)); err != nil {
		handlePasienError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pasien berhasil dihapus",
		"data":    nil,
	})
}

func handlePasienError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrPasienNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Data pasien tidak ditemukan",
			"data":    nil,
		})
	case errors.Is(err, usecase.ErrPasienDuplicateNIK):
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "NIK sudah terdaftar. Gunakan NIK lain atau buka data yang sudah ada.",
			"data":    nil,
		})
	case errors.Is(err, usecase.ErrPasienDuplicateBPJS):
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Nomor BPJS sudah terdaftar",
			"data":    nil,
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan saat memproses data pasien",
			"data":    nil,
		})
	}
}
