package util

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestParseTokenEmptyErrorIs(t *testing.T) {
	_, err := ParseToken("secret", "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
	if !errors.Is(err, ErrTokenInvalid) {
		t.Fatalf("expected ErrTokenInvalid, got %v", err)
	}
}

func TestParseTokenExpiredErrorIs(t *testing.T) {
	secret := "test-secret-123456"
	token, err := GenerateToken(secret, -time.Hour, 7, "lawyer", "lawyer")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	_, err = ParseToken(secret, token)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
	if !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

func TestParseTokenMalformedErrorIs(t *testing.T) {
	_, err := ParseToken("secret", "not-a-token")
	if err == nil {
		t.Fatal("expected error for malformed token")
	}
	if !errors.Is(err, ErrTokenInvalid) {
		t.Fatalf("expected ErrTokenInvalid, got %v", err)
	}
}

func TestParseTokenSignatureErrorIs(t *testing.T) {
	claims := jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	raw, err := tok.SignedString([]byte("wrong-secret"))
	if err != nil {
		t.Fatalf("SignedString error: %v", err)
	}
	_, err = ParseToken("right-secret", raw)
	if err == nil {
		t.Fatal("expected error for signature mismatch")
	}
	if !errors.Is(err, ErrTokenInvalid) {
		t.Fatalf("expected ErrTokenInvalid, got %v", err)
	}
}
