package dto

import (
	"testing"
	"time"
)

func TestParseAcceptDateEmptyNil(t *testing.T) {
	got, err := ParseAcceptDate("")
	if err != nil {
		t.Fatalf("ParseAcceptDate empty error: %v", err)
	}
	if got != nil {
		t.Fatalf("ParseAcceptDate empty = %v, want nil", *got)
	}
}

func TestParseAcceptDateInvalidError(t *testing.T) {
	got, err := ParseAcceptDate("2026-13-99")
	if err == nil {
		t.Fatalf("ParseAcceptDate invalid should error, got %v", got)
	}
}

func TestValidateAcceptDateZero(t *testing.T) {
	z := time.Time{}
	if err := ValidateAcceptDate(&z); err == nil {
		t.Fatal("ValidateAcceptDate zero should error")
	}
}
