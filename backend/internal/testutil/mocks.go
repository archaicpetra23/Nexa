package testutil

import (
	"time"

	"nexa/backend/internal/model"
	"nexa/backend/internal/repository"
)

// This package holds the repository stubs and entity fixtures shared by the
// usecase unit tests. The mocks are ordinary Go files (not _test.go) so the
// stubs are importable from other packages' test files; they are never
// referenced from production code.

// timeZero is a fixed, non-zero timestamp for entity fixtures.
func timeZero() time.Time {
	return time.Date(1995, 5, 15, 0, 0, 0, 0, time.UTC)
}

// NewMockPasien returns a stub whose lookups all report gorm.ErrRecordNotFound.
func NewMockPasien() *mockPasienRepository { return newMockPasien() }

// NewMockKlaim returns a stub whose lookups all report gorm.ErrRecordNotFound.
func NewMockKlaim() *mockKlaimRepository { return newMockKlaim() }

// NewMockRekam returns a stub whose lookups all report gorm.ErrRecordNotFound.
func NewMockRekam() *mockRekamMedisRepository { return newMockRekam() }

// NewValidKlaim builds a minimally valid draft Klaim for create-path tests.
func NewValidKlaim() *model.Klaim { return newValidKlaim() }

// NewValidPasien builds a minimally valid Pasien for create-path tests.
func NewValidPasien() *model.Pasien { return newValidPasien() }

// NewRekamOwnedBy builds a rekam medis owned by doctorID.
func NewRekamOwnedBy(doctorID uint) *model.RekamMedis { return newRekamOwnedBy(doctorID) }

// LastSearchPage and LastSearchLimit expose the arguments the stub's Search
// received, so usecase tests can assert pagination clamping.
func (m *mockKlaimRepository) LastSearchPage() int  { return m.LastSearchArgs.Page }
func (m *mockKlaimRepository) LastSearchLimit() int { return m.LastSearchArgs.Limit }

func (m *mockPasienRepository) LastSearchPage() int  { return m.LastSearchArgs.Page }
func (m *mockPasienRepository) LastSearchLimit() int { return m.LastSearchArgs.Limit }

// SatisfiesRepositoryInterfaces guards against signature drift: if a method is
// added to an interface, this package stops compiling instead of the usecase
// tests silently passing against a partial stub.
var (
	_ repository.PasienRepository     = (*mockPasienRepository)(nil)
	_ repository.KlaimRepository      = (*mockKlaimRepository)(nil)
	_ repository.RekamMedisRepository = (*mockRekamMedisRepository)(nil)
)
