package repository

import (
	"nexa/backend/internal/model"

	"gorm.io/gorm"
)

type RekamMedisRepository interface {
	Create(rekam *model.RekamMedis) error
	FindByID(id uint) (*model.RekamMedis, error)
	FindByPasienID(pasienID uint, page, limit int) ([]model.RekamMedis, int64, error)
	Update(rekam *model.RekamMedis) error
	SoftDelete(id uint) error
	AddDiagnosis(diagnosis *model.RekamDiagnosis) error
	RemoveDiagnosis(id uint) error
	AddTindakan(tindakan *model.DetailTindakan) error
	RemoveTindakan(id uint) error
	GetDiagnosis(rekamID uint) ([]model.RekamDiagnosis, error)
	GetTindakan(rekamID uint) ([]model.DetailTindakan, error)
}

type rekamMedisRepository struct {
	db *gorm.DB
}

func NewRekamMedisRepository(db *gorm.DB) RekamMedisRepository {
	return &rekamMedisRepository{db: db}
}

func (r *rekamMedisRepository) Create(rekam *model.RekamMedis) error {
	return r.db.Create(rekam).Error
}

func (r *rekamMedisRepository) FindByID(id uint) (*model.RekamMedis, error) {
	var rekam model.RekamMedis
	err := r.db.Preload("Pasien").Preload("Dokter").
		Where("id_rekam = ? AND deleted_at IS NULL", id).
		First(&rekam).Error
	if err != nil {
		return nil, err
	}
	return &rekam, nil
}

func (r *rekamMedisRepository) FindByPasienID(pasienID uint, page, limit int) ([]model.RekamMedis, int64, error) {
	var rekams []model.RekamMedis
	var total int64

	offset := (page - 1) * limit

	db := r.db.Where("id_pasien = ? AND deleted_at IS NULL", pasienID)

	if err := db.Model(&model.RekamMedis{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Preload("Pasien").Preload("Dokter").
		Order("tanggal_kunjungan DESC").
		Offset(offset).Limit(limit).
		Find(&rekams).Error; err != nil {
		return nil, 0, err
	}

	return rekams, total, nil
}

func (r *rekamMedisRepository) Update(rekam *model.RekamMedis) error {
	return r.db.Save(rekam).Error
}

func (r *rekamMedisRepository) SoftDelete(id uint) error {
	return r.db.Model(&model.RekamMedis{}).Where("id_rekam = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *rekamMedisRepository) AddDiagnosis(diagnosis *model.RekamDiagnosis) error {
	return r.db.Create(diagnosis).Error
}

func (r *rekamMedisRepository) RemoveDiagnosis(id uint) error {
	return r.db.Delete(&model.RekamDiagnosis{}, id).Error
}

func (r *rekamMedisRepository) AddTindakan(tindakan *model.DetailTindakan) error {
	return r.db.Create(tindakan).Error
}

func (r *rekamMedisRepository) RemoveTindakan(id uint) error {
	return r.db.Delete(&model.DetailTindakan{}, id).Error
}

func (r *rekamMedisRepository) GetDiagnosis(rekamID uint) ([]model.RekamDiagnosis, error) {
	var diagnoses []model.RekamDiagnosis
	err := r.db.Preload("Diagnosis").
		Where("id_rekam = ?", rekamID).
		Find(&diagnoses).Error
	return diagnoses, err
}

func (r *rekamMedisRepository) GetTindakan(rekamID uint) ([]model.DetailTindakan, error) {
	var tindakans []model.DetailTindakan
	err := r.db.Preload("Tindakan").
		Where("id_rekam = ?", rekamID).
		Find(&tindakans).Error
	return tindakans, err
}
