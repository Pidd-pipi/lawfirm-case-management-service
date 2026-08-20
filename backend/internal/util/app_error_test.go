package util

import (
	"errors"
	"testing"
)

func TestWrapPreservesChain(t *testing.T) {
	orig := NewAppError(40400, "not found")
	wrapped := Wrap(orig, "Client[id=%d] get failed", 1)
	var appErr *AppError
	if !errors.As(wrapped, &appErr) {
		t.Fatalf("Wrap 丢失了错误链，errors.As 找不到 AppError: %v", wrapped)
	}
	if appErr.Code != 40400 {
		t.Fatalf("code=%d, want 40400", appErr.Code)
	}
}

func TestIsAppError(t *testing.T) {
	if !IsAppError(NewAppError(42200, "bad")) {
		t.Fatal("IsAppError 对直接 AppError 应返回 true")
	}
	if !IsAppError(Wrap(NewAppError(42200, "bad"), "x")) {
		t.Fatal("IsAppError 对包装后的 AppError 应返回 true")
	}
}
