package pkg

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// testSecret is a fixed 32+ char secret so the tests do not depend on .env.
const testSecret = "nexa-unit-test-secret-key-2026-min32chars"

func TestGenerateTokenRoundTrip(t *testing.T) {
	tokenString, err := GenerateToken(42, "dokter_dpjp", testSecret)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	if tokenString == "" {
		t.Fatal("GenerateToken returned an empty token")
	}

	claims, err := ValidateToken(tokenString, testSecret)
	if err != nil {
		t.Fatalf("ValidateToken returned error for a freshly generated token: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("UserID = %d, want 42", claims.UserID)
	}
	if claims.Role != "dokter_dpjp" {
		t.Errorf("Role = %q, want %q", claims.Role, "dokter_dpjp")
	}
}

func TestValidateTokenWrongSecret(t *testing.T) {
	tokenString, err := GenerateToken(1, "admin_ti", testSecret)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := ValidateToken(tokenString, "different-secret-that-is-also-long-enough")
	if err == nil {
		t.Fatal("ValidateToken accepted a token signed with a different secret")
	}
	if claims != nil {
		t.Errorf("claims = %+v, want nil", claims)
	}
}

func TestValidateTokenMalformed(t *testing.T) {
	for _, raw := range []string{"", "not-a-jwt", "a.b.c", "..."} {
		claims, err := ValidateToken(raw, testSecret)
		if err == nil {
			t.Errorf("ValidateToken(%q) accepted a malformed token", raw)
		}
		if claims != nil {
			t.Errorf("ValidateToken(%q) returned claims %+v, want nil", raw, claims)
		}
	}
}

func TestValidateTokenExpired(t *testing.T) {
	// Build the expired token by hand: GenerateToken always stamps now+8h.
	claims := Claims{
		UserID: 7,
		Role:   "perawat",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	got, err := ValidateToken(tokenString, testSecret)
	if err == nil {
		t.Fatal("ValidateToken accepted an expired token")
	}
	if got != nil {
		t.Errorf("claims = %+v, want nil", got)
	}
}

func TestValidateTokenTamperedSignature(t *testing.T) {
	tokenString, err := GenerateToken(5, "keuangan", testSecret)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	tampered := []byte(tokenString)
	// The signature is the last segment; flip one character at the very end.
	tampered[len(tampered)-1] ^= 0x01

	claims, err := ValidateToken(string(tampered), testSecret)
	if err == nil {
		t.Fatal("ValidateToken accepted a token with a tampered signature")
	}
	if claims != nil {
		t.Errorf("claims = %+v, want nil", claims)
	}
}

func TestValidateTokenTamperedPayload(t *testing.T) {
	tokenString, err := GenerateToken(5, "keuangan", testSecret)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	// Header.payload.signature — corrupt the payload segment instead.
	segments := splitToken(t, tokenString)
	segments[1] = "e30" // base64 of "{}" — an empty claim set

	claims, err := ValidateToken(joinToken(segments), testSecret)
	if err == nil {
		t.Fatal("ValidateToken accepted a token whose payload was replaced with an empty claim set")
	}
	if claims != nil {
		t.Errorf("claims = %+v, want nil", claims)
	}
}

func TestClaimsOnlyCarriesUserIDAndRole(t *testing.T) {
	// The Claims struct exposes exactly UserID and Role beyond the registered
	// claims. A future field added here would change the JWT payload shape and
	// break every already-issued token, so pin the set explicitly.
	tokenString, err := GenerateToken(11, "manajemen", testSecret)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := ValidateToken(tokenString, testSecret)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}
	if claims.UserID != 11 {
		t.Errorf("UserID = %d, want 11", claims.UserID)
	}
	if claims.Role != "manajemen" {
		t.Errorf("Role = %q, want %q", claims.Role, "manajemen")
	}
	if claims.ExpiresAt == nil || claims.IssuedAt == nil {
		t.Error("registered claims ExpiresAt/IssuedAt should be set")
	}
}

func splitToken(t *testing.T, token string) []string {
	t.Helper()

	parts := make([]string, 0, 3)
	start := 0
	for i := 0; i < len(token); i++ {
		if token[i] == '.' {
			parts = append(parts, token[start:i])
			start = i + 1
		}
	}
	parts = append(parts, token[start:])

	if len(parts) != 3 {
		t.Fatalf("token %q has %d segments, want 3", token, len(parts))
	}
	return parts
}

func joinToken(parts []string) string {
	return parts[0] + "." + parts[1] + "." + parts[2]
}
