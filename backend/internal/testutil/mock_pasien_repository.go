package testutil

import (
	"testing"

	"nexa/backend/internal/model"
	"nexa/backend/internal/repository"

	"gorm.io/gorm"
)

// mockPasienRepository is a hand-written stub satisfying PasienRepository.
// Each method delegates to a function field and bumps a call counter, so a
// test can prove the usecase really invoked the dependency (or really did not)
// rather than passing on a silent no-op.
//
// Panics when the Fn field is nil: a nil stub method means the test forgot to
// arrange it, which is a bug in the test, not a condition to swallow.
type mockPasienRepository struct {
	CreateFn       func(*model.Pasien) error
	FindByIDFn     func(uint) (*model.Pasien, error)
	FindByNIKFn    func(string) (*model.Pasien, error)
	FindByNoBPJSFn func(string) (*model.Pasien, error)
	SearchFn       func(string, int, int) ([]model.Pasien, int64, error)
	UpdateFn       func(*model.Pasien) error
	SoftDeleteFn   func(uint) error

	CreateCalls       int
	FindByIDCalls     int
	FindByNIKCalls    int
	FindByNoBPJSCalls int
	SearchCalls       int
	UpdateCalls       int
	SoftDeleteCalls   int

	// LastSearchArgs records the arguments Search received, so pagination
	// clamping can be asserted rather than assumed.
	LastSearchArgs struct {
		Query string
		Page  int
		Limit int
	}
}

func (m *mockPasienRepository) Create(pasien *model.Pasien) error {
	m.CreateCalls++
	if m.CreateFn == nil {
		panic("mockPasienRepository.Create called with no CreateFn set")
	}
	return m.CreateFn(pasien)
}

func (m *mockPasienRepository) FindByID(id uint) (*model.Pasien, error) {
	m.FindByIDCalls++
	if m.FindByIDFn == nil {
		panic("mockPasienRepository.FindByID called with no FindByIDFn set")
	}
	return m.FindByIDFn(id)
}

func (m *mockPasienRepository) FindByNIK(nik string) (*model.Pasien, error) {
	m.FindByNIKCalls++
	if m.FindByNIKFn == nil {
		panic("mockPasienRepository.FindByNIK called with no FindByNIKFn set")
	}
	return m.FindByNIKFn(nik)
}

func (m *mockPasienRepository) FindByNoBPJS(noBpjs string) (*model.Pasien, error) {
	m.FindByNoBPJSCalls++
	if m.FindByNoBPJSFn == nil {
		panic("mockPasienRepository.FindByNoBPJS called with no FindByNoBPJSFn set")
	}
	return m.FindByNoBPJSFn(noBpjs)
}

func (m *mockPasienRepository) Search(query string, page, limit int) ([]model.Pasien, int64, error) {
	m.SearchCalls++
	m.LastSearchArgs = struct {
		Query string
		Page  int
		Limit int
	}{Query: query, Page: page, Limit: limit}
	if m.SearchFn == nil {
		panic("mockPasienRepository.Search called with no SearchFn set")
	}
	return m.SearchFn(query, page, limit)
}

func (m *mockPasienRepository) Update(pasien *model.Pasien) error {
	m.UpdateCalls++
	if m.UpdateFn == nil {
		panic("mockPasienRepository.Update called with no UpdateFn set")
	}
	return m.UpdateFn(pasien)
}

func (m *mockPasienRepository) SoftDelete(id uint) error {
	m.SoftDeleteCalls++
	if m.SoftDeleteFn == nil {
		panic("mockPasienRepository.SoftDelete called with no SoftDeleteFn set")
	}
	return m.SoftDeleteFn(id)
}

// Compile-time proof that the stub satisfies the interface it stands in for.
var _ repository.PasienRepository = (*mockPasienRepository)(nil)

// newMockPasien returns a stub whose lookups all report "not found", which is
// the baseline for the happy-path and duplicate-path tests alike.
func newMockPasien() *mockPasienRepository {
	notFound := func() (*model.Pasien, error) { return nil, gorm.ErrRecordNotFound }
	return &mockPasienRepository{
		FindByIDFn:     func(uint) (*model.Pasien, error) { return notFound() },
		FindByNIKFn:    func(string) (*model.Pasien, error) { return notFound() },
		FindByNoBPJSFn: func(string) (*model.Pasien, error) { return notFound() },
	}
}

// newValidPasien builds a minimally valid Pasien for the create-path tests.
func newValidPasien() *model.Pasien {
	return &model.Pasien{
		IDPasien:     0,
		NIK:          "2290134500000000",
		Nama:         "Tuti Saputra",
		TanggalLahir: timeZero(),
		JenisKelamin: "P",
		Alamat:       "Jl. Melati No. 12, Tangerang",
	}
}

// TestMockPasienRepositorySatisfiesInterface guards against a signature drift
// between the interface and the stub: if a method is added to the interface,
// this file stops compiling instead of silently passing.
func TestMockPasienRepositorySatisfiesInterface(t *testing.T) {
	var repo repository.PasienRepository = newMockPasien()
	if repo == nil {
		t.Fatal("nil PasienRepository")
	}
}
