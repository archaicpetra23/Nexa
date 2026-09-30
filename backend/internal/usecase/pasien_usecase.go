package usecase

import (
	"errors"
	"nexa/backend/internal/model"
	"nexa/backend/internal/repository"
	"time"
)

var (
	ErrPasienNotFound      = errors.New("pasien tidak ditemukan")
	ErrPasienDuplicateNIK  = errors.New("NIK sudah terdaftar")
	ErrPasienDuplicateBPJS = errors.New("No BPJS sudah terdaftar")
)

type PasienUsecase interface {
	CreatePasien(pasien *model.Pasien) error
	GetPasienByID(id uint) (*model.Pasien, error)
	GetPasienByNIK(nik string) (*model.Pasien, error)
	SearchPasien(query string, page, limit int) ([]model.Pasien, int64, error)
	UpdatePasien(pasien *model.Pasien) error
	DeletePasien(id uint) error
}

type pasienUsecase struct {
	repo repository.PasienRepository
}

func NewPasienUsecase(repo repository.PasienRepository) PasienUsecase {
	return &pasienUsecase{repo: repo}
}

func (uc *pasienUsecase) CreatePasien(pasien *model.Pasien) error {
	existing, err := uc.repo.FindByNIK(pasien.NIK)
	if err == nil && existing != nil {
		return ErrPasienDuplicateNIK
	}

	if pasien.NoBPJS != nil && *pasien.NoBPJS != "" {
		existingBPJS, err := uc.repo.FindByNoBPJS(*pasien.NoBPJS)
		if err == nil && existingBPJS != nil {
			return ErrPasienDuplicateBPJS
		}
	}

	pasien.CreatedAt = time.Now()
	pasien.UpdatedAt = time.Now()

	return uc.repo.Create(pasien)
}

func (uc *pasienUsecase) GetPasienByID(id uint) (*model.Pasien, error) {
	pasien, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, ErrPasienNotFound
	}
	return pasien, nil
}

func (uc *pasienUsecase) GetPasienByNIK(nik string) (*model.Pasien, error) {
	pasien, err := uc.repo.FindByNIK(nik)
	if err != nil {
		return nil, ErrPasienNotFound
	}
	return pasien, nil
}

func (uc *pasienUsecase) SearchPasien(query string, page, limit int) ([]model.Pasien, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.repo.Search(query, page, limit)
}

func (uc *pasienUsecase) UpdatePasien(pasien *model.Pasien) error {
	existing, err := uc.repo.FindByID(pasien.IDPasien)
	if err != nil {
		return ErrPasienNotFound
	}

	pasien.CreatedAt = existing.CreatedAt
	pasien.UpdatedAt = time.Now()

	return uc.repo.Update(pasien)
}

func (uc *pasienUsecase) DeletePasien(id uint) error {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return ErrPasienNotFound
	}
	return uc.repo.SoftDelete(id)
}