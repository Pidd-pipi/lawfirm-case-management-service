package config

import "testing"

func TestParseCSVEmptyReturnsNonNil(t *testing.T) {
	got := parseCSV("  ,  ,")
	if got == nil {
		t.Fatal("parseCSV 对空输入返回 nil 切片")
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 items, got %d", len(got))
	}
}
