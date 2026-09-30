package http

import (
	"time"

	"nexa/backend/internal/models"

	"gorm.io/gorm"
)

type UserResponse struct {
	IDUser       uint      `json:"id_user"`
	Username     string    `json:"username"`
	Nama         string    `json:"nama"`
	Profesi      string    `json:"profesi"`
	Spesialisasi *string   `json:"spesialisasi"`
	NoSTR        *string   `json:"no_str"`
	IDRole       uint      `json:"id_role"`
	Role         string    `json:"role"`
	IDUnit       uint      `json:"id_unit"`
	UnitNama     string    `json:"unit_nama"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    *string   `json:"deleted_at,omitempty"`
}

func toUserResponse(db *gorm.DB, u *models.User) UserResponse {
	var deletedAt *string
	if u.DeletedAt != nil {
		s := u.DeletedAt.Format("2006-01-02 15:04:05")
		deletedAt = &s
	}

	roleName := ""
	if u.Role != nil {
		roleName = u.Role.NamaRole
	} else if u.IDRole > 0 {
		var role models.Role
		if err := db.First(&role, u.IDRole).Error; err == nil {
			roleName = role.NamaRole
		}
	}

	unitNama := ""
	if u.Unit != nil {
		unitNama = u.Unit.NamaUnit
	} else if u.IDUnit > 0 {
		var unit models.Unit
		if err := db.First(&unit, u.IDUnit).Error; err == nil {
			unitNama = unit.NamaUnit
		}
	}

	return UserResponse{
		IDUser:       u.IDUser,
		Username:     u.Username,
		Nama:         u.Nama,
		Profesi:      u.Profesi,
		Spesialisasi: u.Spesialisasi,
		NoSTR:        u.NoSTR,
		IDRole:       u.IDRole,
		Role:         roleName,
		IDUnit:       u.IDUnit,
		UnitNama:     unitNama,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}
