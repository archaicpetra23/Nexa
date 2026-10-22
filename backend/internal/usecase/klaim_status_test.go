package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ptr returns a pointer to a copy of s, so each test row owns its own string
// and the table below stays readable.
func ptr(s string) *string {
	v := s
	return &v
}

func TestValidateKlaimStatus(t *testing.T) {
	cases := []struct {
		name    string
		status  string
		alasan  *string
		wantErr error
	}{
		{name: "draft without reason", status: "draft", alasan: nil, wantErr: nil},
		{name: "pending without reason", status: "pending", alasan: nil, wantErr: ErrKlaimAlasanRequired},
		{name: "pending with empty reason", status: "pending", alasan: ptr(""), wantErr: ErrKlaimAlasanRequired},
		{name: "pending with reason", status: "pending", alasan: ptr("Menunggu berkas penunjang"), wantErr: nil},
		{name: "ditolak without reason", status: "ditolak", alasan: nil, wantErr: ErrKlaimAlasanRequired},
		{name: "ditolak with empty reason", status: "ditolak", alasan: ptr(""), wantErr: ErrKlaimAlasanRequired},
		{name: "ditolak with reason", status: "ditolak", alasan: ptr("Koding tidak sesuai"), wantErr: nil},
		{name: "disetujui without reason", status: "disetujui", alasan: nil, wantErr: nil},
		{name: "disetujui with reason", status: "disetujui", alasan: ptr("alasan"), wantErr: nil},
		{name: "status with leading space", status: " Draft", alasan: nil, wantErr: ErrKlaimInvalidStatus},
		{name: "unknown status selesai", status: "selesai", alasan: nil, wantErr: ErrKlaimInvalidStatus},
		{name: "empty status", status: "", alasan: nil, wantErr: ErrKlaimInvalidStatus},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateKlaimStatus(tc.status, tc.alasan)

			if tc.wantErr == nil {
				require.NoError(t, err)
				return
			}
			// Assert by sentinel identity, not by message text.
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestValidateKlaimStatusReasonIsWhitespaceSensitive(t *testing.T) {
	// The backend checks `alasan == nil || *alasan == ""` — a whitespace-only
	// reason passes that check. Pin the current behaviour so a future tightening
	// is a deliberate change, not an accidental one.
	require.NoError(t, validateKlaimStatus("pending", ptr("   ")))
}
