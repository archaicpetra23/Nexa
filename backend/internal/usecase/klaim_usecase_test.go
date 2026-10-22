package usecase

import (
	"testing"

	"nexa/backend/internal/model"
	"nexa/backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateKlaimHappyPath(t *testing.T) {
	repo := testutil.NewMockKlaim()
	repo.CreateFn = func(*model.Klaim) error { return nil }

	klaim := testutil.NewValidKlaim()
	klaim.StatusKlaim = "draft"

	require.NoError(t, NewKlaimUsecase(repo).CreateKlaim(klaim, 22))

	assert.Equal(t, 1, repo.CreateCalls, "Create must be called exactly once")
	assert.False(t, klaim.CreatedAt.IsZero(), "CreatedAt should be stamped by the usecase")
	assert.False(t, klaim.UpdatedAt.IsZero(), "UpdatedAt should be stamped by the usecase")
}

func TestCreateKlaimDuplicateRekam(t *testing.T) {
	repo := testutil.NewMockKlaim()
	repo.FindByRekamIDFn = func(uint) (*model.Klaim, error) {
		return &model.Klaim{IDKlaim: 99, IDRekam: 48648}, nil
	}

	klaim := testutil.NewValidKlaim()

	err := NewKlaimUsecase(repo).CreateKlaim(klaim, 22)

	assert.ErrorIs(t, err, ErrKlaimDuplicateRekam)
	assert.Equal(t, 0, repo.CreateCalls, "Create must not run when the rekam already has a claim")
}

func TestCreateKlaimPendingWithoutAlasan(t *testing.T) {
	repo := testutil.NewMockKlaim()
	repo.CreateFn = func(*model.Klaim) error { return nil }

	klaim := testutil.NewValidKlaim()
	klaim.StatusKlaim = "pending"
	klaim.AlasanPendingTolak = nil

	err := NewKlaimUsecase(repo).CreateKlaim(klaim, 22)

	assert.ErrorIs(t, err, ErrKlaimAlasanRequired)
	assert.Equal(t, 0, repo.CreateCalls, "Create must not run when the reason is missing")
}

func TestCreateKlaimDitolakWithAlasan(t *testing.T) {
	repo := testutil.NewMockKlaim()
	repo.CreateFn = func(*model.Klaim) error { return nil }

	alasan := "Koding tidak sesuai"
	klaim := testutil.NewValidKlaim()
	klaim.StatusKlaim = "ditolak"
	klaim.AlasanPendingTolak = &alasan

	require.NoError(t, NewKlaimUsecase(repo).CreateKlaim(klaim, 22))
	assert.Equal(t, 1, repo.CreateCalls)
}

func TestUpdateStatusKlaimNotFound(t *testing.T) {
	repo := testutil.NewMockKlaim()
	// FindByIDFn already returns gorm.ErrRecordNotFound.

	err := NewKlaimUsecase(repo).UpdateStatusKlaim(404, "disetujui", nil, 22)

	assert.ErrorIs(t, err, ErrKlaimNotFound)
	assert.Equal(t, 0, repo.UpdateCalls, "Update must not run when the claim does not exist")
}

func TestUpdateStatusKlaimDitolakWithoutAlasan(t *testing.T) {
	repo := testutil.NewMockKlaim()
	repo.FindByIDFn = func(uint) (*model.Klaim, error) {
		return testutil.NewValidKlaim(), nil
	}

	err := NewKlaimUsecase(repo).UpdateStatusKlaim(12770, "ditolak", nil, 22)

	assert.ErrorIs(t, err, ErrKlaimAlasanRequired)
	assert.Equal(t, 0, repo.UpdateCalls, "Update must not run when the reason is missing")
}

func TestUpdateStatusKlaimDitolakWithAlasanMutatesRecord(t *testing.T) {
	repo := testutil.NewMockKlaim()
	repo.UpdateFn = func(*model.Klaim) error { return nil }

	var stored *model.Klaim
	repo.FindByIDFn = func(uint) (*model.Klaim, error) {
		stored = testutil.NewValidKlaim()
		return stored, nil
	}

	alasan := "Berkas penunjang tidak lengkap"
	require.NoError(t, NewKlaimUsecase(repo).UpdateStatusKlaim(12770, "ditolak", &alasan, 22))

	assert.Equal(t, 1, repo.UpdateCalls)
	require.NotNil(t, stored)
	assert.Equal(t, "ditolak", stored.StatusKlaim)
	require.NotNil(t, stored.AlasanPendingTolak)
	assert.Equal(t, "Berkas penunjang tidak lengkap", *stored.AlasanPendingTolak)
	assert.False(t, stored.UpdatedAt.IsZero())
}

func TestUpdateStatusKlaimDisetujuiWithoutAlasan(t *testing.T) {
	repo := testutil.NewMockKlaim()
	repo.UpdateFn = func(*model.Klaim) error { return nil }
	repo.FindByIDFn = func(uint) (*model.Klaim, error) { return testutil.NewValidKlaim(), nil }

	require.NoError(t, NewKlaimUsecase(repo).UpdateStatusKlaim(12770, "disetujui", nil, 22))
	assert.Equal(t, 1, repo.UpdateCalls)
}

func TestUpdateStatusKlaimUnknownStatus(t *testing.T) {
	repo := testutil.NewMockKlaim()
	repo.FindByIDFn = func(uint) (*model.Klaim, error) { return testutil.NewValidKlaim(), nil }

	err := NewKlaimUsecase(repo).UpdateStatusKlaim(12770, "arsip", nil, 22)

	assert.ErrorIs(t, err, ErrKlaimInvalidStatus)
	assert.Equal(t, 0, repo.UpdateCalls)
}

func TestSearchKlaimPaginationClamping(t *testing.T) {
	cases := []struct {
		name     string
		page     int
		limit    int
		wantPage int
		wantLmt  int
	}{
		{name: "zero values default to page 1 limit 10", page: 0, limit: 0, wantPage: 1, wantLmt: 10},
		{name: "negative page and oversized limit clamp", page: -5, limit: 999, wantPage: 1, wantLmt: 10},
		{name: "valid values pass through untouched", page: 3, limit: 25, wantPage: 3, wantLmt: 25},
		{name: "limit above the cap clamps to 10", page: 1, limit: 500, wantPage: 1, wantLmt: 10},
		{name: "negative limit defaults to 10", page: 2, limit: -1, wantPage: 2, wantLmt: 10},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := testutil.NewMockKlaim()
			repo.SearchFn = func(string, int, int) ([]model.Klaim, int64, error) {
				return nil, 0, nil
			}

			_, _, err := NewKlaimUsecase(repo).SearchKlaim("pending", tc.page, tc.limit)
			require.NoError(t, err)
			require.Equal(t, 1, repo.SearchCalls)
			assert.Equal(t, tc.wantPage, repo.LastSearchPage(), "page clamping")
			assert.Equal(t, tc.wantLmt, repo.LastSearchLimit(), "limit clamping")
		})
	}
}

func TestDeleteKlaimNotFound(t *testing.T) {
	repo := testutil.NewMockKlaim()

	err := NewKlaimUsecase(repo).DeleteKlaim(999)

	assert.ErrorIs(t, err, ErrKlaimNotFound)
	assert.Equal(t, 0, repo.SoftDeleteCalls)
}

func TestDeleteKlaimSoftDeletesWithID(t *testing.T) {
	repo := testutil.NewMockKlaim()
	repo.FindByIDFn = func(uint) (*model.Klaim, error) { return testutil.NewValidKlaim(), nil }

	var deletedID uint
	repo.SoftDeleteFn = func(id uint) error {
		deletedID = id
		return nil
	}

	require.NoError(t, NewKlaimUsecase(repo).DeleteKlaim(555))
	assert.Equal(t, 1, repo.SoftDeleteCalls)
	assert.Equal(t, uint(555), deletedID, "SoftDelete must receive the claim id")
}

func TestCreateKlaimPropagatesRepositoryError(t *testing.T) {
	repo := testutil.NewMockKlaim()
	// A non-ErrRecordNotFound error from the duplicate lookup must not be
	// mistaken for "no existing claim".
	repo.FindByRekamIDFn = func(uint) (*model.Klaim, error) {
		return nil, gorm.ErrInvalidDB
	}

	err := NewKlaimUsecase(repo).CreateKlaim(testutil.NewValidKlaim(), 22)

	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrKlaimDuplicateRekam)
}
