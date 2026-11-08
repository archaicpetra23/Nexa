package usecase

import (
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"nexa/backend/internal/model"
)

// MockKlaimRepository for testing
type MockKlaimRepository struct {
	CreateFunc      func(klaim *model.Klaim) error
	FindByIDFunc    func(id uint) (*model.Klaim, error)
	FindByRekamIDFunc func(rekamID uint) (*model.Klaim, error)
	UpdateFunc      func(klaim *model.Klaim) error
	SearchFunc      func(status string, page, limit int) ([]model.Klaim, int64, error)
	SoftDeleteFunc  func(id uint) error
}

func (m *MockKlaimRepository) Create(klaim *model.Klaim) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(klaim)
	}
	return nil
}

func (m *MockKlaimRepository) FindByID(id uint) (*model.Klaim, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockKlaimRepository) FindByRekamID(rekamID uint) (*model.Klaim, error) {
	if m.FindByRekamIDFunc != nil {
		return m.FindByRekamIDFunc(rekamID)
	}
	return nil, nil
}

func (m *MockKlaimRepository) Update(klaim *model.Klaim) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(klaim)
	}
	return nil
}

func (m *MockKlaimRepository) Search(status string, page, limit int) ([]model.Klaim, int64, error) {
	return m.SearchFunc(status, page, limit)
}

func (m *MockKlaimRepository) SoftDelete(id uint) error {
	return m.SoftDeleteFunc(id)
}

func TestCreateKlaimSuccess(t *testing.T) {
	repo := &MockKlaimRepository{
		FindByRekamIDFunc: func(rekamID uint) (*model.Klaim, error) {
			return nil, gorm.ErrRecordNotFound // gorm.ErrRecordNotFound simulation
		},
		CreateFunc: func(klaim *model.Klaim) error {
			return nil
		},
	}

	uc := NewKlaimUsecase(repo)

	klaim := &model.Klaim{
		IDRekam:      1,
		KodeCBGS:     "A1",
		IDPetugasCasemix: 2,
		StatusKlaim:  "draft",
		NominalKlaim: 1000000,
	}

	err := uc.CreateKlaim(klaim, 1)
	if err != nil {
		t.Errorf("CreateKlaim returned error: %v", err)
	}
}

func TestCreateKlaimDuplicateRekam(t *testing.T) {
	existing := &model.Klaim{IDKlaim: 1}
	repo := &MockKlaimRepository{
		FindByRekamIDFunc: func(rekamID uint) (*model.Klaim, error) {
			return existing, nil
		},
	}

	uc := NewKlaimUsecase(repo)

	klaim := &model.Klaim{IDRekam: 1}
	err := uc.CreateKlaim(klaim, 1)
	if !errors.Is(err, ErrKlaimDuplicateRekam) {
		t.Errorf("CreateKlaim expected ErrKlaimDuplicateRekam, got %v", err)
	}
}

func TestCreateKlaimInvalidStatus(t *testing.T) {
	repo := &MockKlaimRepository{
		FindByRekamIDFunc: func(rekamID uint) (*model.Klaim, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	uc := NewKlaimUsecase(repo)

	// invalid status "cancelled"
	klaim := &model.Klaim{
		IDRekam:      1,
		KodeCBGS:     "A1",
		IDPetugasCasemix: 2,
		StatusKlaim:  "cancelled",
	}

	err := uc.CreateKlaim(klaim, 1)
	if !errors.Is(err, ErrKlaimInvalidStatus) {
		t.Errorf("CreateKlaim expected ErrKlaimInvalidStatus, got %v", err)
	}
}

func TestCreateKlaimPendingWithoutAlasan(t *testing.T) {
	repo := &MockKlaimRepository{
		FindByIDFunc: func(id uint) (*model.Klaim, error) {
			return nil, ErrKlaimAlasanRequired // ini yang dibutuhkan
		},
	}

	uc := NewKlaimUsecase(repo)

	klaim := &model.Klaim{
		IDRekam:      1,
		KodeCBGS:     "A1",
		IDPetugasCasemix: 2,
		StatusKlaim:  "pending",
	}

	err := uc.CreateKlaim(klaim, 1)
	if !errors.Is(err, ErrKlaimAlasanRequired) {
		t.Errorf("CreateKlaim expected ErrKlaimAlasanRequired, got %v", err)
	}
}

func TestUpdateStatusKlaimSuccess(t *testing.T) {
	now := time.Now()
	existing := &model.Klaim{
		IDKlaim:       1,
		StatusKlaim:   "draft",
		AlasanPendingTolak: nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	repo := &MockKlaimRepository{
		FindByIDFunc: func(id uint) (*model.Klaim, error) {
			return existing, nil
		},
		UpdateFunc: func(klaim *model.Klaim) error {
			return nil
		},
	}

	uc := NewKlaimUsecase(repo)

	alasan := "berkas lengkap"
	err := uc.UpdateStatusKlaim(1, "pending", &alasan, 1)
	if err != nil {
		t.Errorf("UpdateStatusKlaim returned error: %v", err)
	}
	if existing.StatusKlaim != "pending" {
		t.Errorf("Status not updated, got %q", existing.StatusKlaim)
	}
	if existing.AlasanPendingTolak == nil || *existing.AlasanPendingTolak != alasan {
		t.Errorf("Alasan not set correctly")
	}
}

func TestUpdateStatusKlaimInvalidTransition(t *testing.T) {
	existing := &model.Klaim{
		IDKlaim:     1,
		StatusKlaim: "disetujui",
	}
	repo := &MockKlaimRepository{
		FindByIDFunc: func(id uint) (*model.Klaim, error) {
			return existing, nil
		},
	}

	uc := NewKlaimUsecase(repo)

	// Cannot go back from disetujui to draft
	err := uc.UpdateStatusKlaim(1, "draft", nil, 1)
	if !errors.Is(err, ErrKlaimInvalidStatus) {
		t.Errorf("Expected ErrKlaimInvalidStatus for invalid transition, got %v", err)
	}
}

func TestSearchKlaimPaginated(t *testing.T) {
	mockKlaims := []model.Klaim{
		{IDKlaim: 1, StatusKlaim: "draft"},
		{IDKlaim: 2, StatusKlaim: "pending"},
	}
	repo := &MockKlaimRepository{
		SearchFunc: func(status string, page, limit int) ([]model.Klaim, int64, error) {
			return mockKlaims, 2, nil
		},
	}

	uc := NewKlaimUsecase(repo)

	items, total, err := uc.SearchKlaim("pending", 1, 10)
	if err != nil {
		t.Errorf("SearchKlaim returned error: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("items count = %d, want 2", len(items))
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
}
