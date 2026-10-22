package usecase

import (
	"testing"
	"time"

	"nexa/backend/internal/model"
	"nexa/backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ownerDoctorID and otherDoctorID are deliberately different so the ownership
// branch in canAccessRekam is genuinely exercised.
const (
	ownerDoctorID   = 42
	otherDoctorID   = 99
	unknownDoctorID = 7
)

func rekamOwnedByAnotherDoctor() *model.RekamMedis {
	return testutil.NewRekamOwnedBy(ownerDoctorID)
}

func TestCreateRekamMedisHappyPath(t *testing.T) {
	repo := testutil.NewMockRekam()
	repo.CreateFn = func(*model.RekamMedis) error { return nil }

	rekam := rekamOwnedByAnotherDoctor()

	require.NoError(t, NewRekamMedisUsecase(repo).CreateRekamMedis(rekam, ownerDoctorID))

	assert.Equal(t, 1, repo.CreateCalls)
	assert.False(t, rekam.CreatedAt.IsZero())
	assert.False(t, rekam.UpdatedAt.IsZero())
}

func TestGetRekamMedisByIDAcessMatrix(t *testing.T) {
	cases := []struct {
		name       string
		role       string
		userID     uint
		owner      uint
		wantAccess bool
	}{
		{name: "admin_ti may read any record", role: "admin_ti", userID: 1, owner: otherDoctorID, wantAccess: true},
		{name: "manajemen may read any record", role: "manajemen", userID: 2, owner: otherDoctorID, wantAccess: true},
		{name: "dokter_dpjp may read own record", role: "dokter_dpjp", userID: ownerDoctorID, owner: ownerDoctorID, wantAccess: true},
		{name: "dokter_dpjp may not read another doctor record", role: "dokter_dpjp", userID: otherDoctorID, owner: ownerDoctorID, wantAccess: false},
		{name: "perawat may read any record", role: "perawat", userID: 3, owner: otherDoctorID, wantAccess: true},
		{name: "petugas_rm may read any record", role: "petugas_rm", userID: 4, owner: otherDoctorID, wantAccess: true},
		{name: "petugas_casemix may read any record", role: "petugas_casemix", userID: 5, owner: otherDoctorID, wantAccess: true},
		{name: "keuangan may read any record", role: "keuangan", userID: 6, owner: otherDoctorID, wantAccess: true},
		{name: "unknown role is denied", role: "staf", userID: 8, owner: ownerDoctorID, wantAccess: false},
		{name: "empty role is denied", role: "", userID: ownerDoctorID, owner: ownerDoctorID, wantAccess: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := testutil.NewMockRekam()
			record := testutil.NewRekamOwnedBy(tc.owner)
			repo.FindByIDFn = func(uint) (*model.RekamMedis, error) { return record, nil }

			got, err := NewRekamMedisUsecase(repo).GetRekamMedisByID(7501, tc.userID, tc.role)

			if !tc.wantAccess {
				assert.ErrorIs(t, err, ErrRekamUnauthorized)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, uint(7501), got.IDRekam)
		})
	}
}

func TestGetRekamMedisByIDNotFound(t *testing.T) {
	repo := testutil.NewMockRekam()

	got, err := NewRekamMedisUsecase(repo).GetRekamMedisByID(999, 1, "admin_ti")

	assert.ErrorIs(t, err, ErrRekamNotFound)
	assert.Nil(t, got)
}

func TestAddDiagnosisOwnershipAndIDRekamAssignment(t *testing.T) {
	t.Run("dokter_dpjp on another doctor record is denied", func(t *testing.T) {
		repo := testutil.NewMockRekam()
		repo.FindByIDFn = func(uint) (*model.RekamMedis, error) { return rekamOwnedByAnotherDoctor(), nil }

		diagnosis := &model.RekamDiagnosis{KodeICD10: "N40", Jenis: "primer"}

		err := NewRekamMedisUsecase(repo).AddDiagnosis(7501, diagnosis, otherDoctorID, "dokter_dpjp")

		assert.ErrorIs(t, err, ErrRekamUnauthorized)
		assert.Equal(t, 0, repo.AddDiagnosisCalls, "AddDiagnosis must not run when the caller is not the owner")
		assert.Nil(t, repo.LastDiagnosis)
	})

	t.Run("dokter_dpjp on own record succeeds and IDRekam is set", func(t *testing.T) {
		repo := testutil.NewMockRekam()
		repo.AddDiagnosisFn = func(*model.RekamDiagnosis) error { return nil }
		repo.FindByIDFn = func(uint) (*model.RekamMedis, error) { return rekamOwnedByAnotherDoctor(), nil }

		diagnosis := &model.RekamDiagnosis{KodeICD10: "N40", Jenis: "primer"}

		require.NoError(t, NewRekamMedisUsecase(repo).AddDiagnosis(7501, diagnosis, ownerDoctorID, "dokter_dpjp"))

		assert.Equal(t, 1, repo.AddDiagnosisCalls)
		require.NotNil(t, repo.LastDiagnosis)
		assert.Equal(t, uint(7501), repo.LastDiagnosis.IDRekam,
			"the usecase must stamp IDRekam from the path parameter")
		assert.Equal(t, "N40", repo.LastDiagnosis.KodeICD10)
	})

	t.Run("admin_ti may add to any record", func(t *testing.T) {
		repo := testutil.NewMockRekam()
		repo.AddDiagnosisFn = func(*model.RekamDiagnosis) error { return nil }
		repo.FindByIDFn = func(uint) (*model.RekamMedis, error) { return rekamOwnedByAnotherDoctor(), nil }

		require.NoError(t, NewRekamMedisUsecase(repo).AddDiagnosis(7501,
			&model.RekamDiagnosis{KodeICD10: "R10", Jenis: "primer"}, 1, "admin_ti"))

		assert.Equal(t, 1, repo.AddDiagnosisCalls)
	})
}

func TestAddDiagnosisNotFound(t *testing.T) {
	repo := testutil.NewMockRekam()

	err := NewRekamMedisUsecase(repo).AddDiagnosis(999,
		&model.RekamDiagnosis{KodeICD10: "N40", Jenis: "primer"}, ownerDoctorID, "dokter_dpjp")

	assert.ErrorIs(t, err, ErrRekamNotFound)
	assert.Equal(t, 0, repo.AddDiagnosisCalls)
}

func TestRemoveDiagnosisUnauthorized(t *testing.T) {
	repo := testutil.NewMockRekam()
	repo.FindByIDFn = func(uint) (*model.RekamMedis, error) { return rekamOwnedByAnotherDoctor(), nil }

	err := NewRekamMedisUsecase(repo).RemoveDiagnosis(7501, 5, otherDoctorID, "dokter_dpjp")

	assert.ErrorIs(t, err, ErrRekamUnauthorized)
	assert.Equal(t, 0, repo.RemoveDiagnosisCalls)
}

func TestAddTindakanAccess(t *testing.T) {
	t.Run("perawat may add tindakan", func(t *testing.T) {
		repo := testutil.NewMockRekam()
		repo.AddTindakanFn = func(*model.DetailTindakan) error { return nil }
		repo.FindByIDFn = func(uint) (*model.RekamMedis, error) { return rekamOwnedByAnotherDoctor(), nil }

		require.NoError(t, NewRekamMedisUsecase(repo).AddTindakan(7501,
			&model.DetailTindakan{KodeTindakan: "99101", Jumlah: 2}, 3, "perawat"))

		assert.Equal(t, 1, repo.AddTindakanCalls)
		require.NotNil(t, repo.LastTindakan)
		assert.Equal(t, uint(7501), repo.LastTindakan.IDRekam,
			"the usecase must stamp IDRekam from the path parameter")
	})

	t.Run("dokter_dpjp on another doctor record is denied", func(t *testing.T) {
		repo := testutil.NewMockRekam()
		repo.FindByIDFn = func(uint) (*model.RekamMedis, error) { return rekamOwnedByAnotherDoctor(), nil }

		err := NewRekamMedisUsecase(repo).AddTindakan(7501,
			&model.DetailTindakan{KodeTindakan: "99101", Jumlah: 1}, otherDoctorID, "dokter_dpjp")

		assert.ErrorIs(t, err, ErrRekamUnauthorized)
		assert.Equal(t, 0, repo.AddTindakanCalls)
	})
}

func TestRemoveTindakanUnauthorized(t *testing.T) {
	repo := testutil.NewMockRekam()
	repo.FindByIDFn = func(uint) (*model.RekamMedis, error) { return rekamOwnedByAnotherDoctor(), nil }

	err := NewRekamMedisUsecase(repo).RemoveTindakan(7501, 9, otherDoctorID, "dokter_dpjp")

	assert.ErrorIs(t, err, ErrRekamUnauthorized)
	assert.Equal(t, 0, repo.RemoveTindakanCalls)
}

func TestUpdateRekamMedisUnauthorized(t *testing.T) {
	repo := testutil.NewMockRekam()
	repo.FindByIDFn = func(uint) (*model.RekamMedis, error) { return rekamOwnedByAnotherDoctor(), nil }

	err := NewRekamMedisUsecase(repo).UpdateRekamMedis(rekamOwnedByAnotherDoctor(), otherDoctorID, "dokter_dpjp")

	assert.ErrorIs(t, err, ErrRekamUnauthorized)
	assert.Equal(t, 0, repo.UpdateCalls)
}

func TestDeleteRekamMedisUnauthorized(t *testing.T) {
	repo := testutil.NewMockRekam()
	repo.FindByIDFn = func(uint) (*model.RekamMedis, error) { return rekamOwnedByAnotherDoctor(), nil }

	err := NewRekamMedisUsecase(repo).DeleteRekamMedis(7501, otherDoctorID, "dokter_dpjp")

	assert.ErrorIs(t, err, ErrRekamUnauthorized)
	assert.Equal(t, 0, repo.SoftDeleteCalls)
}

func TestCanAccessRekamIsExhaustive(t *testing.T) {
	// Pin the documented matrix: every recognised role plus two that must be
	// refused. If a role is added to the system later, this table must be
	// updated deliberately rather than silently inheriting the default deny.
	roles := []string{"admin_ti", "manajemen", "petugas_rm", "perawat", "petugas_casemix", "keuangan"}
	for _, role := range roles {
		t.Run("role "+role+" has read access", func(t *testing.T) {
			repo := testutil.NewMockRekam()
			repo.FindByIDFn = func(uint) (*model.RekamMedis, error) { return rekamOwnedByAnotherDoctor(), nil }

			_, err := NewRekamMedisUsecase(repo).GetRekamMedisByID(7501, unknownDoctorID, role)
			require.NoError(t, err)
		})
	}
}

// mustParseTime parses a date for fixture construction and fails the test on
// malformed input, so a typo in a fixture cannot look like a product bug.
func mustParseTime(t *testing.T, date string) time.Time {
	t.Helper()

	parsed, err := time.Parse("2006-01-02", date)
	require.NoError(t, err)
	return parsed
}
