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
	ListDiagnosis(search string, page, limit int) ([]model.Diagnosis, int64, error)
	ListTindakan(search string, page, limit int) ([]model.Tindakan, int64, error)
	ListCBGS(search string, page, limit int) ([]model.TarifCBGs, int64, error)
	SearchDiagnosisByCodeOrName(query string, limit int) ([]model.Diagnosis, error)
	SearchTindakanByCodeOrName(query string, limit int) ([]model.Tindakan, error)
	GetDiagnosisByCode(code string) (*model.Diagnosis, error)
	GetTindakanByCode(code string) (*model.Tindakan, error)
	CreateDiagnosis(d *model.Diagnosis) error
	UpdateDiagnosis(d *model.Diagnosis) error
	DeleteDiagnosis(code string) error
	CreateTindakan(t *model.Tindakan) error
	UpdateTindakan(t *model.Tindakan) error
	DeleteTindakan(code string) error
	CreateCBGS(c *model.TarifCBGs) error
	UpdateCBGS(c *model.TarifCBGs) error
	DeleteCBGS(code string) error
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

func (r *masterRepository) ListDiagnosis(search string, page, limit int) ([]model.Diagnosis, int64, error) {
	var diagnoses []model.Diagnosis
	var total int64
	q := r.db.Model(&model.Diagnosis{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("kode_icd10 ILIKE ? OR nama_diagnosis ILIKE ?", like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Session(&gorm.Session{}).Order("kode_icd10 ASC").
		Offset((page - 1) * limit).Limit(limit).Find(&diagnoses).Error
	return diagnoses, total, err
}

func (r *masterRepository) ListTindakan(search string, page, limit int) ([]model.Tindakan, int64, error) {
	var tindakans []model.Tindakan
	var total int64
	q := r.db.Model(&model.Tindakan{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("kode_tindakan ILIKE ? OR nama_tindakan ILIKE ?", like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Session(&gorm.Session{}).Order("kode_tindakan ASC").
		Offset((page - 1) * limit).Limit(limit).Find(&tindakans).Error
	return tindakans, total, err
}

func (r *masterRepository) ListCBGS(search string, page, limit int) ([]model.TarifCBGs, int64, error) {
	var cbgsList []model.TarifCBGs
	var total int64
	q := r.db.Model(&model.TarifCBGs{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("kode_cbgs ILIKE ? OR deskripsi ILIKE ?", like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Session(&gorm.Session{}).Order("kode_cbgs ASC").
		Offset((page - 1) * limit).Limit(limit).Find(&cbgsList).Error
	return cbgsList, total, err
}

func (r *masterRepository) SearchDiagnosisByCodeOrName(query string, limit int) ([]model.Diagnosis, error) {
	var diagnoses []model.Diagnosis
	err := r.db.Where("kode_icd10 ILIKE ? OR nama_diagnosis ILIKE ?", query+"%", "%"+query+"%").
		Order("kode_icd10 ASC").Limit(limit).Find(&diagnoses).Error
	return diagnoses, err
}

func (r *masterRepository) SearchTindakanByCodeOrName(query string, limit int) ([]model.Tindakan, error) {
	var tindakans []model.Tindakan
	err := r.db.Where("kode_tindakan ILIKE ? OR nama_tindakan ILIKE ?", query+"%", "%"+query+"%").
		Order("kode_tindakan ASC").Limit(limit).Find(&tindakans).Error
	return tindakans, err
}

func (r *masterRepository) GetDiagnosisByCode(code string) (*model.Diagnosis, error) {
	var diagnosis model.Diagnosis
	err := r.db.Where("kode_icd10 = ?", code).First(&diagnosis).Error
	if err != nil {
		return nil, err
	}
	return &diagnosis, nil
}

func (r *masterRepository) GetTindakanByCode(code string) (*model.Tindakan, error) {
	var tindakan model.Tindakan
	err := r.db.Where("kode_tindakan = ?", code).First(&tindakan).Error
	if err != nil {
		return nil, err
	}
	return &tindakan, nil
}

func (r *masterRepository) CreateDiagnosis(d *model.Diagnosis) error {
	return r.db.Create(d).Error
}

func (r *masterRepository) UpdateDiagnosis(d *model.Diagnosis) error {
	return r.db.Save(d).Error
}

func (r *masterRepository) DeleteDiagnosis(code string) error {
	return r.db.Where("kode_icd10 = ?", code).Delete(&model.Diagnosis{}).Error
}

func (r *masterRepository) CreateTindakan(t *model.Tindakan) error {
	return r.db.Create(t).Error
}

func (r *masterRepository) UpdateTindakan(t *model.Tindakan) error {
	return r.db.Save(t).Error
}

func (r *masterRepository) DeleteTindakan(code string) error {
	return r.db.Where("kode_tindakan = ?", code).Delete(&model.Tindakan{}).Error
}

func (r *masterRepository) CreateCBGS(c *model.TarifCBGs) error {
	return r.db.Create(c).Error
}

func (r *masterRepository) UpdateCBGS(c *model.TarifCBGs) error {
	return r.db.Save(c).Error
}

func (r *masterRepository) DeleteCBGS(code string) error {
	return r.db.Where("kode_cbgs = ?", code).Delete(&model.TarifCBGs{}).Error
}
