package usecase

import (
	"testing"

	"nexa/backend/internal/model"
	"nexa/backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreatePasienHappyPath(t *testing.T) {
	repo := testutil.NewMockPasien()
	repo.CreateFn = func(*model.Pasien) error { return nil }

	pasien := testutil.NewValidPasien()

	require.NoError(t, NewPasienUsecase(repo).CreatePasien(pasien))

	assert.Equal(t, 1, repo.CreateCalls, "Create must be called exactly once")
	assert.Equal(t, 1, repo.FindByNIKCalls)
	assert.Equal(t, 0, repo.FindByNoBPJSCalls, "NoBPJS is nil, so the BPJS lookup must be skipped")
	assert.False(t, pasien.CreatedAt.IsZero())
	assert.False(t, pasien.UpdatedAt.IsZero())
}

func TestCreatePasienDuplicateNIK(t *testing.T) {
	repo := testutil.NewMockPasien()
	repo.FindByNIKFn = func(string) (*model.Pasien, error) {
		return &model.Pasien{IDPasien: 5, NIK: "2290134500000000"}, nil
	}

	err := NewPasienUsecase(repo).CreatePasien(testutil.NewValidPasien())

	assert.ErrorIs(t, err, ErrPasienDuplicateNIK)
	assert.Equal(t, 0, repo.CreateCalls, "Create must not run when the NIK is already registered")
}

func TestCreatePasienDuplicateBPJS(t *testing.T) {
	repo := testutil.NewMockPasien()
	noBPJS := "3117600000000"

	pasien := testutil.NewValidPasien()
	pasien.NoBPJS = &noBPJS

	repo.FindByNoBPJSFn = func(string) (*model.Pasien, error) {
		return &model.Pasien{IDPasien: 77, NIK: "1111111111111111", NoBPJS: &noBPJS}, nil
	}

	err := NewPasienUsecase(repo).CreatePasien(pasien)

	assert.ErrorIs(t, err, ErrPasienDuplicateBPJS)
	assert.Equal(t, 1, repo.FindByNoBPJSCalls)
	assert.Equal(t, 0, repo.CreateCalls)
}

func TestCreatePasienEmptyNoBPJSSkipsLookup(t *testing.T) {
	repo := testutil.NewMockPasien()
	repo.CreateFn = func(*model.Pasien) error { return nil }

	empty := ""
	pasien := testutil.NewValidPasien()
	pasien.NoBPJS = &empty

	require.NoError(t, NewPasienUsecase(repo).CreatePasien(pasien))

	assert.Equal(t, 0, repo.FindByNoBPJSCalls,
		"NoBPJS is an empty string, so the BPJS lookup must be skipped")
	assert.Equal(t, 1, repo.CreateCalls)
}

func TestGetPasienByIDFound(t *testing.T) {
	repo := testutil.NewMockPasien()
	want := testutil.NewValidPasien()
	want.IDPasien = 42
	repo.FindByIDFn = func(uint) (*model.Pasien, error) { return want, nil }

	got, err := NewPasienUsecase(repo).GetPasienByID(42)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, uint(42), got.IDPasien)
	assert.Equal(t, want.NIK, got.NIK)
}

func TestGetPasienByIDNotFound(t *testing.T) {
	repo := testutil.NewMockPasien()

	got, err := NewPasienUsecase(repo).GetPasienByID(999)

	assert.ErrorIs(t, err, ErrPasienNotFound)
	assert.Nil(t, got)
}

func TestGetPasienByNIKNotFound(t *testing.T) {
	repo := testutil.NewMockPasien()

	got, err := NewPasienUsecase(repo).GetPasienByNIK("0000000000000000")

	assert.ErrorIs(t, err, ErrPasienNotFound)
	assert.Nil(t, got)
}

func TestUpdatePasienNotFound(t *testing.T) {
	repo := testutil.NewMockPasien()

	err := NewPasienUsecase(repo).UpdatePasien(testutil.NewValidPasien())

	assert.ErrorIs(t, err, ErrPasienNotFound)
	assert.Equal(t, 0, repo.UpdateCalls)
}

func TestUpdatePasienPreservesCreatedAt(t *testing.T) {
	repo := testutil.NewMockPasien()

	existing := testutil.NewValidPasien()
	existing.IDPasien = 42
	existing.CreatedAt = mustParseTime(t, "1995-05-15")
	existing.UpdatedAt = mustParseTime(t, "2020-01-01")

	repo.FindByIDFn = func(uint) (*model.Pasien, error) { return existing, nil }
	repo.UpdateFn = func(*model.Pasien) error { return nil }

	before := existing.CreatedAt
	update := testutil.NewValidPasien()
	update.IDPasien = 42

	require.NoError(t, NewPasienUsecase(repo).UpdatePasien(update))

	assert.Equal(t, 1, repo.UpdateCalls)
	assert.Equal(t, before, existing.CreatedAt, "the usecase must preserve the stored CreatedAt")
	assert.True(t, existing.UpdatedAt.After(before),
		"UpdatedAt should advance beyond the fixture value")
}

func TestSearchPasienPaginationClamping(t *testing.T) {
	cases := []struct {
		name     string
		page     int
		limit    int
		wantPage int
		wantLmt  int
	}{
		{name: "zero values default to page 1 limit 10", page: 0, limit: 0, wantPage: 1, wantLmt: 10},
		{name: "negative page and oversized limit clamp", page: -1, limit: 500, wantPage: 1, wantLmt: 10},
		{name: "valid values pass through untouched", page: 4, limit: 30, wantPage: 4, wantLmt: 30},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := testutil.NewMockPasien()
			repo.SearchFn = func(string, int, int) ([]model.Pasien, int64, error) {
				return nil, 0, nil
			}

			_, _, err := NewPasienUsecase(repo).SearchPasien("tuti", tc.page, tc.limit)
			require.NoError(t, err)
			require.Equal(t, 1, repo.SearchCalls)
			assert.Equal(t, tc.wantPage, repo.LastSearchPage(), "page clamping")
			assert.Equal(t, tc.wantLmt, repo.LastSearchLimit(), "limit clamping")
		})
	}
}

func TestDeletePasienNotFound(t *testing.T) {
	repo := testutil.NewMockPasien()

	err := NewPasienUsecase(repo).DeletePasien(888)

	assert.ErrorIs(t, err, ErrPasienNotFound)
	assert.Equal(t, 0, repo.SoftDeleteCalls)
}

func TestDeletePasienSoftDeletesWithID(t *testing.T) {
	repo := testutil.NewMockPasien()
	repo.FindByIDFn = func(uint) (*model.Pasien, error) { return testutil.NewValidPasien(), nil }

	var deletedID uint
	repo.SoftDeleteFn = func(id uint) error {
		deletedID = id
		return nil
	}

	require.NoError(t, NewPasienUsecase(repo).DeletePasien(321))
	assert.Equal(t, 1, repo.SoftDeleteCalls)
	assert.Equal(t, uint(321), deletedID)
}

func TestCreatePasienPropagatesRepositoryError(t *testing.T) {
	repo := testutil.NewMockPasien()
	repo.FindByNIKFn = func(string) (*model.Pasien, error) { return nil, gorm.ErrInvalidDB }

	err := NewPasienUsecase(repo).CreatePasien(testutil.NewValidPasien())

	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrPasienDuplicateNIK,
		"a database failure must not be reported as a duplicate NIK")
}
