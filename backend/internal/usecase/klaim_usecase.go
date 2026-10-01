package usecase

import (
	"errors"
	"nexa/backend/internal/model"
	"nexa/backend/internal/repository"
	"time"
)

var (
	ErrKlaimNotFound       = errors.New("klaim tidak ditemukan")
	ErrKlaimDuplicateRekam = errors.New("klaim untuk rekam medis ini sudah ada")
	ErrKlaimAlasanRequired = errors.New("alasan pending/tolak wajib diisi")
	ErrKlaimInvalidStatus  = errors.New("status klaim tidak valid")
)

type KlaimUsecase interface {
	CreateKlaim(klaim *model.Klaim, userID uint) error
	GetKlaimByID(id uint) (*model.Klaim, error)
	UpdateStatusKlaim(id uint, status string, alasan *string, userID uint) error
	SearchKlaim(status string, page, limit int) ([]model.Klaim, int64, error)
	DeleteKlaim(id uint) error
}

type klaimUsecase struct {
	repo repository.KlaimRepository
}

func NewKlaimUsecase(repo repository.KlaimRepository) KlaimUsecase {
	return &klaimUsecase{repo: repo}
}

func (uc *klaimUsecase) CreateKlaim(klaim *model.Klaim, userID uint) error {
	existing, err := uc.repo.FindByRekamID(klaim.IDRekam)
	if err == nil && existing != nil {
		return ErrKlaimDuplicateRekam
	}

	if err := validateKlaimStatus(klaim.StatusKlaim, klaim.AlasanPendingTolak); err != nil {
		return err
	}

	klaim.CreatedAt = time.Now()
	klaim.UpdatedAt = time.Now()

	return uc.repo.Create(klaim)
}

func (uc *klaimUsecase) GetKlaimByID(id uint) (*model.Klaim, error) {
	klaim, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, ErrKlaimNotFound
	}
	return klaim, nil
}

func (uc *klaimUsecase) UpdateStatusKlaim(id uint, status string, alasan *string, userID uint) error {
	klaim, err := uc.repo.FindByID(id)
	if err != nil {
		return ErrKlaimNotFound
	}

	if err := validateKlaimStatus(status, alasan); err != nil {
		return err
	}

	klaim.StatusKlaim = status
	klaim.AlasanPendingTolak = alasan
	klaim.UpdatedAt = time.Now()

	return uc.repo.Update(klaim)
}

func (uc *klaimUsecase) SearchKlaim(status string, page, limit int) ([]model.Klaim, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.repo.Search(status, page, limit)
}

func (uc *klaimUsecase) DeleteKlaim(id uint) error {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return ErrKlaimNotFound
	}
	return uc.repo.SoftDelete(id)
}

func validateKlaimStatus(status string, alasan *string) error {
	validStatuses := map[string]bool{
		"draft":     true,
		"pending":   true,
		"disetujui": true,
		"ditolak":   true,
	}

	if !validStatuses[status] {
		return ErrKlaimInvalidStatus
	}

	if (status == "pending" || status == "ditolak") && (alasan == nil || *alasan == "") {
		return ErrKlaimAlasanRequired
	}

	return nil
}
