package repository

import (
	"nexa/backend/internal/model"

	"gorm.io/gorm"
)

type KlaimRepository interface {
	Create(klaim *model.Klaim) error
	FindByID(id uint) (*model.Klaim, error)
	FindByRekamID(rekamID uint) (*model.Klaim, error)
	Update(klaim *model.Klaim) error
	Search(status string, page, limit int) ([]model.Klaim, int64, error)
	SoftDelete(id uint) error
}

type klaimRepository struct {
	db *gorm.DB
}

func NewKlaimRepository(db *gorm.DB) KlaimRepository {
	return &klaimRepository{db: db}
}

func (r *klaimRepository) Create(klaim *model.Klaim) error {
	return r.db.Create(klaim).Error
}

func (r *klaimRepository) FindByID(id uint) (*model.Klaim, error) {
	var klaim model.Klaim
	err := r.db.Preload("RekamMedis").Preload("KodeCBGSNavigation").Preload("PetugasCasemix").
		Where("id_klaim = ? AND deleted_at IS NULL", id).
		First(&klaim).Error
	if err != nil {
		return nil, err
	}
	return &klaim, nil
}

func (r *klaimRepository) FindByRekamID(rekamID uint) (*model.Klaim, error) {
	var klaim model.Klaim
	err := r.db.Where("id_rekam = ? AND deleted_at IS NULL", rekamID).First(&klaim).Error
	if err != nil {
		return nil, err
	}
	return &klaim, nil
}

func (r *klaimRepository) Update(klaim *model.Klaim) error {
	return r.db.Save(klaim).Error
}

func (r *klaimRepository) Search(status string, page, limit int) ([]model.Klaim, int64, error) {
	var klaims []model.Klaim
	var total int64

	offset := (page - 1) * limit

	db := r.db.Where("deleted_at IS NULL")
	if status != "" {
		db = db.Where("status_klaim = ?", status)
	}

	if err := db.Model(&model.Klaim{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Preload("RekamMedis").Preload("KodeCBGSNavigation").Preload("PetugasCasemix").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&klaims).Error; err != nil {
		return nil, 0, err
	}

	return klaims, total, nil
}

func (r *klaimRepository) SoftDelete(id uint) error {
	return r.db.Model(&model.Klaim{}).Where("id_klaim = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}