package pkg

import (
	"strings"
	"testing"
)

func TestHashPasswordProducesBcryptCost12Hash(t *testing.T) {
	const password = "password123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned an empty hash")
	}
	if hash == password {
		t.Error("HashPassword returned the plaintext password")
	}
	if len(hash) != 60 {
		t.Errorf("hash length = %d bytes, want 60 (bcrypt format)", len(hash))
	}
	if !strings.HasPrefix(hash, "$2a$12$") {
		t.Errorf("hash prefix = %q, want %q (bcrypt, cost 12)", hash[:7], "$2a$12$")
	}
}

func TestCheckPasswordHashCorrectPassword(t *testing.T) {
	const password = "rahasia-nexa-2026"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash rejected the correct password")
	}
}

func TestCheckPasswordHashWrongPassword(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	for _, wrong := range []string{"password124", "Password123", "password123 ", ""} {
		if CheckPasswordHash(wrong, hash) {
			t.Errorf("CheckPasswordHash accepted wrong password %q", wrong)
		}
	}
}

func TestCheckPasswordHashRejectsMalformedHash(t *testing.T) {
	for _, bad := range []string{"", "not-a-hash", "$2a$12$short"} {
		if CheckPasswordHash("password123", bad) {
			t.Errorf("CheckPasswordHash accepted malformed hash %q", bad)
		}
	}
}

func TestHashPasswordUsesRandomSalt(t *testing.T) {
	const password = "sama-sekali"

	first, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	second, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if first == second {
		t.Error("two hashes of the same password are identical — bcrypt salt is not being randomised")
	}
	if !CheckPasswordHash(password, first) || !CheckPasswordHash(password, second) {
		t.Error("at least one of the two hashes fails verification")
	}
}
