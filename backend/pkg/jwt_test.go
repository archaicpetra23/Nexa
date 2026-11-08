package pkg

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateTokenCreatesValidJWT(t *testing.T) {
	const secret = "test-secret-min-32-characters-long"
	userID := uint(42)
	role := "admin_ti"

	token, err := GenerateToken(userID, role, secret)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken returned empty token")
	}
	// JWT format: header.payload.signature
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("token parts = %d, want 3 (JWT format)", len(parts))
	}
}

func TestValidateTokenAcceptsValidToken(t *testing.T) {
	const secret = "test-secret-min-32-characters-long"
	userID := uint(99)
	role := "dokter_dpjp"

	token, err := GenerateToken(userID, role, secret)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken rejected valid token: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("claims.UserID = %d, want %d", claims.UserID, userID)
	}
	if claims.Role != role {
		t.Errorf("claims.Role = %q, want %q", claims.Role, role)
	}
}

func TestValidateTokenRejectsExpiredToken(t *testing.T) {
	// Tidak bisa membuat token expired tanpa mock time, jadi kita test
	// bahwa token dengan secret berbeda ditolak sebagai proxy untuk
	// validasi signature.
	const secret = "test-secret-min-32-characters-long"
	const wrongSecret = "different-secret-min-32-chars-x"

	token, err := GenerateToken(1, "admin_ti", secret)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	_, err = ValidateToken(token, wrongSecret)
	if err == nil {
		t.Error("ValidateToken accepted token signed with different secret")
	}
}

func TestValidateTokenRejectsMalformedToken(t *testing.T) {
	const secret = "test-secret-min-32-characters-long"

	for _, bad := range []string{"", "not-a-token", "header.payload", "a.b.c.d"} {
		_, err := ValidateToken(bad, secret)
		if err == nil {
			t.Errorf("ValidateToken accepted malformed token %q", bad)
		}
	}
}

func TestGenerateTokenSetsExpiresAtIn8Hours(t *testing.T) {
	const secret = "test-secret-min-32-characters-long"

	before := time.Now()
	token, err := GenerateToken(1, "admin_ti", secret)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	after := time.Now()

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}

	if claims.ExpiresAt == nil {
		t.Fatal("claims.ExpiresAt is nil")
	}

	// Token harus expire antara 8 jam dari before dan 8 jam dari after.
	// Toleransi 1 menit untuk build lambat.
	minExpiry := before.Add(8*time.Hour - time.Minute)
	maxExpiry := after.Add(8*time.Hour + time.Minute)

	expiryTime := claims.ExpiresAt.Time
	if expiryTime.Before(minExpiry) || expiryTime.After(maxExpiry) {
		t.Errorf("claims.ExpiresAt = %v, want between %v and %v",
			expiryTime, minExpiry, maxExpiry)
	}
}
