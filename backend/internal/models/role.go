package models

import "time"

type Role struct {
	IDRole   uint   `gorm:"primaryKey;column:id_role" json:"id_role"`
	NamaRole string `gorm:"unique;not null;column:nama_role" json:"nama_role"`
}

func (Role) TableName() string {
	return "roles"
}

type User struct {
	IDUser       uint       `gorm:"primaryKey;column:id_user" json:"id_user"`
	Nama         string     `gorm:"not null" json:"nama"`
	Profesi      string     `gorm:"not null" json:"profesi"`
	Spesialisasi *string    `json:"spesialisasi"`
	NoSTR        *string    `gorm:"unique;column:no_str" json:"no_str"`
	Username     string     `gorm:"unique;not null" json:"username"`
	PasswordHash string     `gorm:"not null;column:password_hash" json:"-"`
	IDRole       uint       `gorm:"not null;column:id_role" json:"id_role"`
	IDUnit       uint       `gorm:"not null;column:id_unit" json:"id_unit"`
	CreatedAt    time.Time  `gorm:"default:now();column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"default:now();column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	Role         *Role      `gorm:"foreignKey:IDRole" json:"role,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type Unit struct {
	IDUnit   uint   `gorm:"primaryKey;column:id_unit" json:"id_unit"`
	NamaUnit string `gorm:"unique;not null;column:nama_unit" json:"nama_unit"`
}

func (Unit) TableName() string {
	return "units"
}
