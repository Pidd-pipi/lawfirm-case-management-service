package util

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteFilesReturnsFirstError(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.txt")
	err := DeleteFiles([]string{missing})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDeleteFilesWrapsError(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.txt")
	err := DeleteFiles([]string{missing})
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, got %v", err)
	}
}

func TestDeleteFilesContinuesAfterError(t *testing.T) {
	dir := t.TempDir()
	real := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(real, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(dir, "nope.txt")
	err := DeleteFiles([]string{missing, real})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if _, statErr := os.Stat(real); !os.IsNotExist(statErr) {
		t.Fatalf("real file should still be removed, stat err=%v", statErr)
	}
}
