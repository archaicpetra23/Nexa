package http

import (
	"errors"
	"net/http"
	"strconv"

	"nexa/backend/internal/model"
	"nexa/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type KlaimHandler struct {
	usecase usecase.KlaimUsecase
}

func NewKlaimHandler(uc usecase.KlaimUsecase) *KlaimHandler {
	return &KlaimHandler{usecase: uc}
}

type CreateKlaimRequest struct {
	IDRekam            uint    `json:"id_rekam" binding:"required"`
	KodeCBGS           string  `json:"kode_cbgs" binding:"required"`
	NominalKlaim       float64 `json:"nominal_klaim" binding:"gte=0"`
	AlasanPendingTolak *string `json:"alasan_pending_tolak"`
}

type UpdateStatusRequest struct {
	StatusKlaim        string  `json:"status_klaim" binding:"required,oneof=draft pending disetujui ditolak"`
	AlasanPendingTolak *string `json:"alasan_pending_tolak"`
}

// CreateKlaim handles POST /api/v1/klaim
func (h *KlaimHandler) CreateKlaim(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateKlaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data klaim tidak valid",
			"data":    nil,
		})
		return
	}

	klaim := &model.Klaim{
		IDRekam:            req.IDRekam,
		KodeCBGS:           req.KodeCBGS,
		IDPetugasCasemix:   userID,
		StatusKlaim:        "draft",
		NominalKlaim:       req.NominalKlaim,
		AlasanPendingTolak: req.AlasanPendingTolak,
	}

	if err := h.usecase.CreateKlaim(klaim, userID); err != nil {
		handleKlaimError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Klaim berhasil dibuat",
		"data":    klaim,
	})
}

// GetKlaim handles GET /api/v1/klaim/:id
func (h *KlaimHandler) GetKlaim(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID klaim tidak valid",
			"data":    nil,
		})
		return
	}

	klaim, err := h.usecase.GetKlaimByID(uint(id))
	if err != nil {
		handleKlaimError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail klaim berhasil dimuat",
		"data":    klaim,
	})
}

// SearchKlaim handles GET /api/v1/klaim?status=&page=1&limit=10
func (h *KlaimHandler) SearchKlaim(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	klaims, total, err := h.usecase.SearchKlaim(status, page, limit)
	if err != nil {
		handleKlaimError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar klaim berhasil dimuat",
		"data": gin.H{
			"items": klaims,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// UpdateStatus handles PATCH /api/v1/klaim/:id/status
func (h *KlaimHandler) UpdateStatus(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID klaim tidak valid",
			"data":    nil,
		})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status klaim tidak valid. Gunakan: draft, pending, disetujui, ditolak",
			"data":    nil,
		})
		return
	}

	if err := h.usecase.UpdateStatusKlaim(uint(id), req.StatusKlaim, req.AlasanPendingTolak, userID); err != nil {
		handleKlaimError(c, err)
		return
	}

	klaim, _ := h.usecase.GetKlaimByID(uint(id))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status klaim berhasil diperbarui",
		"data":    klaim,
	})
}

// DeleteKlaim handles DELETE /api/v1/klaim/:id (soft delete)
func (h *KlaimHandler) DeleteKlaim(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID klaim tidak valid",
			"data":    nil,
		})
		return
	}

	if err := h.usecase.DeleteKlaim(uint(id)); err != nil {
		handleKlaimError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Klaim berhasil dihapus",
		"data":    nil,
	})
}

func handleKlaimError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrKlaimNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Data klaim tidak ditemukan",
			"data":    nil,
		})
	case errors.Is(err, usecase.ErrKlaimDuplicateRekam):
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Klaim untuk rekam medis ini sudah ada",
			"data":    nil,
		})
	case errors.Is(err, usecase.ErrKlaimAlasanRequired):
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": "Alasan wajib diisi ketika status klaim adalah pending atau ditolak",
			"data":    nil,
		})
	case errors.Is(err, usecase.ErrKlaimInvalidStatus):
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status klaim tidak valid",
			"data":    nil,
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan saat memproses klaim",
			"data":    nil,
		})
	}
}
