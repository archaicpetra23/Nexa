package usecase

import (
	"errors"
	"nexa/backend/internal/model"
	"nexa/backend/internal/repository"
	"time"
)

var (
	ErrRekamNotFound      = errors.New("rekam medis tidak ditemukan")
	ErrRekamUnauthorized  = errors.New("tidak memiliki akses ke rekam medis ini")
	ErrDiagnosisDuplicate = errors.New("diagnosis sudah ada untuk rekam medis ini")
	ErrTindakanDuplicate  = errors.New("tindakan sudah ada untuk rekam medis ini")
)

type RekamMedisUsecase interface {
	CreateRekamMedis(rekam *model.RekamMedis, userID uint) error
	GetRekamMedisByID(id uint, userID uint, role string) (*model.RekamMedis, error)
	GetRekamMedisByPasien(pasienID uint, page, limit int) ([]model.RekamMedis, int64, error)
	UpdateRekamMedis(rekam *model.RekamMedis, userID uint, role string) error
	DeleteRekamMedis(id uint, userID uint, role string) error
	AddDiagnosis(rekamID uint, diagnosis *model.RekamDiagnosis, userID uint, role string) error
	RemoveDiagnosis(rekamID, diagnosisID uint, userID uint, role string) error
	AddTindakan(rekamID uint, tindakan *model.DetailTindakan, userID uint, role string) error
	RemoveTindakan(rekamID, tindakanID uint, userID uint, role string) error
}

type rekamMedisUsecase struct {
	repo repository.RekamMedisRepository
}

func NewRekamMedisUsecase(repo repository.RekamMedisRepository) RekamMedisUsecase {
	return &rekamMedisUsecase{repo: repo}
}

func (uc *rekamMedisUsecase) CreateRekamMedis(rekam *model.RekamMedis, userID uint) error {
	rekam.CreatedAt = time.Now()
	rekam.UpdatedAt = time.Now()
	return uc.repo.Create(rekam)
}

func (uc *rekamMedisUsecase) GetRekamMedisByID(id uint, userID uint, role string) (*model.RekamMedis, error) {
	rekam, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, ErrRekamNotFound
	}

	if !canAccessRekam(rekam, userID, role) {
		return nil, ErrRekamUnauthorized
	}

	return rekam, nil
}

func (uc *rekamMedisUsecase) GetRekamMedisByPasien(pasienID uint, page, limit int) ([]model.RekamMedis, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.repo.FindByPasienID(pasienID, page, limit)
}

func (uc *rekamMedisUsecase) UpdateRekamMedis(rekam *model.RekamMedis, userID uint, role string) error {
	existing, err := uc.repo.FindByID(rekam.IDRekam)
	if err != nil {
		return ErrRekamNotFound
	}

	if !canAccessRekam(existing, userID, role) {
		return ErrRekamUnauthorized
	}

	rekam.CreatedAt = existing.CreatedAt
	rekam.UpdatedAt = time.Now()

	return uc.repo.Update(rekam)
}

func (uc *rekamMedisUsecase) DeleteRekamMedis(id uint, userID uint, role string) error {
	rekam, err := uc.repo.FindByID(id)
	if err != nil {
		return ErrRekamNotFound
	}

	if !canAccessRekam(rekam, userID, role) {
		return ErrRekamUnauthorized
	}

	return uc.repo.SoftDelete(id)
}

func (uc *rekamMedisUsecase) AddDiagnosis(rekamID uint, diagnosis *model.RekamDiagnosis, userID uint, role string) error {
	rekam, err := uc.repo.FindByID(rekamID)
	if err != nil {
		return ErrRekamNotFound
	}

	if !canAccessRekam(rekam, userID, role) {
		return ErrRekamUnauthorized
	}

	diagnosis.IDRekam = rekamID
	return uc.repo.AddDiagnosis(diagnosis)
}

func (uc *rekamMedisUsecase) RemoveDiagnosis(rekamID, diagnosisID uint, userID uint, role string) error {
	rekam, err := uc.repo.FindByID(rekamID)
	if err != nil {
		return ErrRekamNotFound
	}

	if !canAccessRekam(rekam, userID, role) {
		return ErrRekamUnauthorized
	}

	return uc.repo.RemoveDiagnosis(diagnosisID)
}

func (uc *rekamMedisUsecase) AddTindakan(rekamID uint, tindakan *model.DetailTindakan, userID uint, role string) error {
	rekam, err := uc.repo.FindByID(rekamID)
	if err != nil {
		return ErrRekamNotFound
	}

	if !canAccessRekam(rekam, userID, role) {
		return ErrRekamUnauthorized
	}

	tindakan.IDRekam = rekamID
	return uc.repo.AddTindakan(tindakan)
}

func (uc *rekamMedisUsecase) RemoveTindakan(rekamID, tindakanID uint, userID uint, role string) error {
	rekam, err := uc.repo.FindByID(rekamID)
	if err != nil {
		return ErrRekamNotFound
	}

	if !canAccessRekam(rekam, userID, role) {
		return ErrRekamUnauthorized
	}

	return uc.repo.RemoveTindakan(tindakanID)
}

func canAccessRekam(rekam *model.RekamMedis, userID uint, role string) bool {
	if role == "admin_ti" || role == "manajemen" {
		return true
	}

	if role == "dokter_dpjp" && rekam.IDDokter == userID {
		return true
	}

	if role == "petugas_rm" || role == "perawat" || role == "petugas_casemix" || role == "keuangan" {
		return true
	}

	return false
}
