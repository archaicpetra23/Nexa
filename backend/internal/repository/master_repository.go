package repository

import (
	"nexa/backend/internal/model"

	"gorm.io/gorm"
)

type MasterRepository interface {
	SearchDiagnosis(query string, limit int) ([]model.Diagnosis, error)
	SearchTindakan(query string, limit int) ([]model.Tindakan, error)
	GetAllDiagnosis() ([]model.Diagnosis, error)
	GetAllTindakan() ([]model.Tindakan, error)
	GetCBGSByCode(code string) (*model.TarifCBGs, error)
	GetAllCBGS() ([]model.TarifCBGs, error)
}

type masterRepository struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) MasterRepository {
	return &masterRepository{db: db}
}

func (r *masterRepository) SearchDiagnosis(query string, limit int) ([]model.Diagnosis, error) {
	var diagnoses []model.Diagnosis
	err := r.db.Where("nama_diagnosis ILIKE ?", "%"+query+"%").
		Limit(limit).
		Find(&diagnoses).Error
	return diagnoses, err
}

func (r *masterRepository) SearchTindakan(query string, limit int) ([]model.Tindakan, error) {
	var tindakans []model.Tindakan
	err := r.db.Where("nama_tindakan ILIKE ?", "%"+query+"%").
		Limit(limit).
		Find(&tindakans).Error
	return tindakans, err
}

func (r *masterRepository) GetAllDiagnosis() ([]model.Diagnosis, error) {
	var diagnoses []model.Diagnosis
	err := r.db.Find(&diagnoses).Error
	return diagnoses, err
}

func (r *masterRepository) GetAllTindakan() ([]model.Tindakan, error) {
	var tindakans []model.Tindakan
	err := r.db.Find(&tindakans).Error
	return tindakans, err
}

func (r *masterRepository) GetCBGSByCode(code string) (*model.TarifCBGs, error) {
	var cbgs model.TarifCBGs
	err := r.db.Where("kode_cbgs = ?", code).First(&cbgs).Error
	if err != nil {
		return nil, err
	}
	return &cbgs, nil
}

func (r *masterRepository) GetAllCBGS() ([]model.TarifCBGs, error) {
	var cbgsList []model.TarifCBGs
	err := r.db.Find(&cbgsList).Error
	return cbgsList, err
}