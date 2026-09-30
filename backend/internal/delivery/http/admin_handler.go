package http

import (
	"net/http"

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

type UserResponse struct {
	IDUser      uint   `json:"id_user"`
	Username    string `json:"username"`
	Nama        string `json:"nama"`
	Profesi     string `json:"profesi"`
	Role        string `json:"role"`
	IDUnit      uint   `json:"id_unit"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
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

// GetUsers handles GET /api/v1/admin/users
func (h *AdminHandler) GetUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Preload("Role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memuat data pengguna",
			"data":    nil,
		})
		return
	}

	response := make([]UserResponse, len(users))
	for i, u := range users {
		var deletedAt *string
		if u.DeletedAt != nil {
			s := u.DeletedAt.Format("2006-01-02 15:04:05")
			deletedAt = &s
		}
		// BUG-004 FIX: guard nil Role pointer with fallback
		roleName := ""
		if u.Role != nil {
			roleName = u.Role.NamaRole
		} else if u.IDRole > 0 {
			var role models.Role
			if err := h.db.First(&role, u.IDRole).Error; err == nil {
				roleName = role.NamaRole
			}
		}
		response[i] = UserResponse{
			IDUser:    u.IDUser,
			Username:  u.Username,
			Nama:      u.Nama,
			Profesi:   u.Profesi,
			Role:      roleName,
			IDUnit:    u.IDUnit,
			DeletedAt: deletedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar pengguna berhasil dimuat",
		"data":    response,
	})
}

// GetSubmissions handles GET /api/v1/admin/submissions
func (h *AdminHandler) GetSubmissions(c *gin.Context) {
	var klaims []model.Klaim
	if err := h.db.Preload("RekamMedis.Pasien").
		Preload("PetugasCasemix").
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&klaims).Error; err != nil {
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
		"data":    response,
	})
}