package repository

import (
	"nexa/backend/internal/model"

	"gorm.io/gorm"
)

type PasienRepository interface {
	Create(pasien *model.Pasien) error
	FindByID(id uint) (*model.Pasien, error)
	FindByNIK(nik string) (*model.Pasien, error)
	FindByNoBPJS(noBpjs string) (*model.Pasien, error)
	Search(query string, page, limit int) ([]model.Pasien, int64, error)
	Update(pasien *model.Pasien) error
	SoftDelete(id uint) error
}

type pasienRepository struct {
	db *gorm.DB
}

func NewPasienRepository(db *gorm.DB) PasienRepository {
	return &pasienRepository{db: db}
}

func (r *pasienRepository) Create(pasien *model.Pasien) error {
	return r.db.Create(pasien).Error
}

func (r *pasienRepository) FindByID(id uint) (*model.Pasien, error) {
	var pasien model.Pasien
	err := r.db.Where("id_pasien = ? AND deleted_at IS NULL", id).First(&pasien).Error
	if err != nil {
		return nil, err
	}
	return &pasien, nil
}

func (r *pasienRepository) FindByNIK(nik string) (*model.Pasien, error) {
	var pasien model.Pasien
	err := r.db.Where("nik = ? AND deleted_at IS NULL", nik).First(&pasien).Error
	if err != nil {
		return nil, err
	}
	return &pasien, nil
}

func (r *pasienRepository) FindByNoBPJS(noBpjs string) (*model.Pasien, error) {
	var pasien model.Pasien
	err := r.db.Where("no_bpjs = ? AND deleted_at IS NULL", noBpjs).First(&pasien).Error
	if err != nil {
		return nil, err
	}
	return &pasien, nil
}

func (r *pasienRepository) Search(query string, page, limit int) ([]model.Pasien, int64, error) {
	var pasiens []model.Pasien
	var total int64

	offset := (page - 1) * limit

	db := r.db.Where("deleted_at IS NULL")
	if query != "" {
		db = db.Where("nik ILIKE ? OR no_bpjs ILIKE ? OR nama ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	if err := db.Model(&model.Pasien{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(limit).Find(&pasiens).Error; err != nil {
		return nil, 0, err
	}

	return pasiens, total, nil
}

func (r *pasienRepository) Update(pasien *model.Pasien) error {
	return r.db.Save(pasien).Error
}

func (r *pasienRepository) SoftDelete(id uint) error {
	return r.db.Model(&model.Pasien{}).Where("id_pasien = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}