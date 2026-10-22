package testutil

import (
	"testing"
	"time"

	"nexa/backend/internal/model"
	"nexa/backend/internal/repository"

	"gorm.io/gorm"
)

// mockRekamMedisRepository stubs RekamMedisRepository. Function fields, call
// counters, panic on an unset Fn — same contract as the other two stubs.
type mockRekamMedisRepository struct {
	CreateFn          func(*model.RekamMedis) error
	FindByIDFn        func(uint) (*model.RekamMedis, error)
	FindByPasienIDFn  func(uint, int, int) ([]model.RekamMedis, int64, error)
	FindAllFn         func(int, int) ([]model.RekamMedis, int64, error)
	UpdateFn          func(*model.RekamMedis) error
	SoftDeleteFn      func(uint) error
	AddDiagnosisFn    func(*model.RekamDiagnosis) error
	RemoveDiagnosisFn func(uint) error
	AddTindakanFn     func(*model.DetailTindakan) error
	RemoveTindakanFn  func(uint) error
	GetDiagnosisFn    func(uint) ([]model.RekamDiagnosis, error)
	GetTindakanFn     func(uint) ([]model.DetailTindakan, error)

	CreateCalls          int
	FindByIDCalls        int
	FindByPasienIDCalls  int
	FindAllCalls         int
	UpdateCalls          int
	SoftDeleteCalls      int
	AddDiagnosisCalls    int
	RemoveDiagnosisCalls int
	AddTindakanCalls     int
	RemoveTindakanCalls  int
	GetDiagnosisCalls    int
	GetTindakanCalls     int

	// LastDiagnosis and LastTindakan capture what the usecase handed to the
	// repository, so the "IDRekam must be set by the usecase" rule is asserted
	// on the real argument instead of inferred.
	LastDiagnosis *model.RekamDiagnosis
	LastTindakan  *model.DetailTindakan
}

func (m *mockRekamMedisRepository) Create(rekam *model.RekamMedis) error {
	m.CreateCalls++
	if m.CreateFn == nil {
		panic("mockRekamMedisRepository.Create called with no CreateFn set")
	}
	return m.CreateFn(rekam)
}

func (m *mockRekamMedisRepository) FindByID(id uint) (*model.RekamMedis, error) {
	m.FindByIDCalls++
	if m.FindByIDFn == nil {
		panic("mockRekamMedisRepository.FindByID called with no FindByIDFn set")
	}
	return m.FindByIDFn(id)
}

func (m *mockRekamMedisRepository) FindByPasienID(pasienID uint, page, limit int) ([]model.RekamMedis, int64, error) {
	m.FindByPasienIDCalls++
	if m.FindByPasienIDFn == nil {
		panic("mockRekamMedisRepository.FindByPasienID called with no FindByPasienIDFn set")
	}
	return m.FindByPasienIDFn(pasienID, page, limit)
}

func (m *mockRekamMedisRepository) FindAll(page, limit int) ([]model.RekamMedis, int64, error) {
	m.FindAllCalls++
	if m.FindAllFn == nil {
		panic("mockRekamMedisRepository.FindAll called with no FindAllFn set")
	}
	return m.FindAllFn(page, limit)
}

func (m *mockRekamMedisRepository) Update(rekam *model.RekamMedis) error {
	m.UpdateCalls++
	if m.UpdateFn == nil {
		panic("mockRekamMedisRepository.Update called with no UpdateFn set")
	}
	return m.UpdateFn(rekam)
}

func (m *mockRekamMedisRepository) SoftDelete(id uint) error {
	m.SoftDeleteCalls++
	if m.SoftDeleteFn == nil {
		panic("mockRekamMedisRepository.SoftDelete called with no SoftDeleteFn set")
	}
	return m.SoftDeleteFn(id)
}

func (m *mockRekamMedisRepository) AddDiagnosis(diagnosis *model.RekamDiagnosis) error {
	m.AddDiagnosisCalls++
	m.LastDiagnosis = diagnosis
	if m.AddDiagnosisFn == nil {
		panic("mockRekamMedisRepository.AddDiagnosis called with no AddDiagnosisFn set")
	}
	return m.AddDiagnosisFn(diagnosis)
}

func (m *mockRekamMedisRepository) RemoveDiagnosis(id uint) error {
	m.RemoveDiagnosisCalls++
	if m.RemoveDiagnosisFn == nil {
		panic("mockRekamMedisRepository.RemoveDiagnosis called with no RemoveDiagnosisFn set")
	}
	return m.RemoveDiagnosisFn(id)
}

func (m *mockRekamMedisRepository) AddTindakan(tindakan *model.DetailTindakan) error {
	m.AddTindakanCalls++
	m.LastTindakan = tindakan
	if m.AddTindakanFn == nil {
		panic("mockRekamMedisRepository.AddTindakan called with no AddTindakanFn set")
	}
	return m.AddTindakanFn(tindakan)
}

func (m *mockRekamMedisRepository) RemoveTindakan(id uint) error {
	m.RemoveTindakanCalls++
	if m.RemoveTindakanFn == nil {
		panic("mockRekamMedisRepository.RemoveTindakan called with no RemoveTindakanFn set")
	}
	return m.RemoveTindakanFn(id)
}

func (m *mockRekamMedisRepository) GetDiagnosis(rekamID uint) ([]model.RekamDiagnosis, error) {
	m.GetDiagnosisCalls++
	if m.GetDiagnosisFn == nil {
		panic("mockRekamMedisRepository.GetDiagnosis called with no GetDiagnosisFn set")
	}
	return m.GetDiagnosisFn(rekamID)
}

func (m *mockRekamMedisRepository) GetTindakan(rekamID uint) ([]model.DetailTindakan, error) {
	m.GetTindakanCalls++
	if m.GetTindakanFn == nil {
		panic("mockRekamMedisRepository.GetTindakan called with no GetTindakanFn set")
	}
	return m.GetTindakanFn(rekamID)
}

var _ repository.RekamMedisRepository = (*mockRekamMedisRepository)(nil)

// newMockRekam returns a stub whose lookups all report "not found".
func newMockRekam() *mockRekamMedisRepository {
	notFound := func() (*model.RekamMedis, error) { return nil, gorm.ErrRecordNotFound }
	return &mockRekamMedisRepository{
		FindByIDFn: func(uint) (*model.RekamMedis, error) { return notFound() },
	}
}

// newRekamOwnedBy builds a rekam medis owned by doctorID. Tests that need to
// exercise the ownership rule pass a different userID/role than the owner.
func newRekamOwnedBy(doctorID uint) *model.RekamMedis {
	return &model.RekamMedis{
		IDRekam:          7501,
		IDPasien:         75833,
		IDDokter:         doctorID,
		TanggalKunjungan: time.Date(2026, 10, 17, 8, 58, 16, 0, time.UTC),
		JenisPerawatan:   "rawat_inap",
		Keluhan:          "Demam dan nyeri perut",
	}
}

func TestMockRekamMedisRepositorySatisfiesInterface(t *testing.T) {
	var repo repository.RekamMedisRepository = newMockRekam()
	if repo == nil {
		t.Fatal("nil RekamMedisRepository")
	}
}
