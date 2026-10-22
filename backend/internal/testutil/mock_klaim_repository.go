package testutil

import (
	"testing"

	"nexa/backend/internal/model"
	"nexa/backend/internal/repository"

	"gorm.io/gorm"
)

// mockKlaimRepository stubs KlaimRepository. Same contract as the pasien stub:
// function fields, call counters, panic on an unset Fn.
type mockKlaimRepository struct {
	CreateFn        func(*model.Klaim) error
	FindByIDFn      func(uint) (*model.Klaim, error)
	FindByRekamIDFn func(uint) (*model.Klaim, error)
	UpdateFn        func(*model.Klaim) error
	SearchFn        func(string, int, int) ([]model.Klaim, int64, error)
	SoftDeleteFn    func(uint) error

	CreateCalls        int
	FindByIDCalls      int
	FindByRekamIDCalls int
	UpdateCalls        int
	SearchCalls        int
	SoftDeleteCalls    int

	LastSearchArgs struct {
		Status string
		Page   int
		Limit  int
	}
}

func (m *mockKlaimRepository) Create(klaim *model.Klaim) error {
	m.CreateCalls++
	if m.CreateFn == nil {
		panic("mockKlaimRepository.Create called with no CreateFn set")
	}
	return m.CreateFn(klaim)
}

func (m *mockKlaimRepository) FindByID(id uint) (*model.Klaim, error) {
	m.FindByIDCalls++
	if m.FindByIDFn == nil {
		panic("mockKlaimRepository.FindByID called with no FindByIDFn set")
	}
	return m.FindByIDFn(id)
}

func (m *mockKlaimRepository) FindByRekamID(rekamID uint) (*model.Klaim, error) {
	m.FindByRekamIDCalls++
	if m.FindByRekamIDFn == nil {
		panic("mockKlaimRepository.FindByRekamID called with no FindByRekamIDFn set")
	}
	return m.FindByRekamIDFn(rekamID)
}

func (m *mockKlaimRepository) Update(klaim *model.Klaim) error {
	m.UpdateCalls++
	if m.UpdateFn == nil {
		panic("mockKlaimRepository.Update called with no UpdateFn set")
	}
	return m.UpdateFn(klaim)
}

func (m *mockKlaimRepository) Search(status string, page, limit int) ([]model.Klaim, int64, error) {
	m.SearchCalls++
	m.LastSearchArgs = struct {
		Status string
		Page   int
		Limit  int
	}{Status: status, Page: page, Limit: limit}
	if m.SearchFn == nil {
		panic("mockKlaimRepository.Search called with no SearchFn set")
	}
	return m.SearchFn(status, page, limit)
}

func (m *mockKlaimRepository) SoftDelete(id uint) error {
	m.SoftDeleteCalls++
	if m.SoftDeleteFn == nil {
		panic("mockKlaimRepository.SoftDelete called with no SoftDeleteFn set")
	}
	return m.SoftDeleteFn(id)
}

var _ repository.KlaimRepository = (*mockKlaimRepository)(nil)

// newMockKlaim returns a stub whose lookups all report "not found".
func newMockKlaim() *mockKlaimRepository {
	notFound := func() (*model.Klaim, error) { return nil, gorm.ErrRecordNotFound }
	return &mockKlaimRepository{
		FindByIDFn:      func(uint) (*model.Klaim, error) { return notFound() },
		FindByRekamIDFn: func(uint) (*model.Klaim, error) { return notFound() },
	}
}

// newValidKlaim builds a minimally valid draft Klaim.
func newValidKlaim() *model.Klaim {
	return &model.Klaim{
		IDRekam:          48648,
		KodeCBGS:         "F-4-10-I",
		IDPetugasCasemix: 22,
		StatusKlaim:      "draft",
		NominalKlaim:     4500000,
	}
}

func TestMockKlaimRepositorySatisfiesInterface(t *testing.T) {
	var repo repository.KlaimRepository = newMockKlaim()
	if repo == nil {
		t.Fatal("nil KlaimRepository")
	}
}
