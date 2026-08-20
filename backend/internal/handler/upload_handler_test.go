package handler

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"cylawcase/internal/config"

	"github.com/gin-gonic/gin"
)

func TestUploadHandlerUsesRequestContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fw, err := mw.CreateFormFile("file", "a.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	if err := mw.Close(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := httptest.NewRequest(http.MethodPost, "/upload", &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := NewUploadHandler(&config.Config{UploadDir: t.TempDir(), UploadMaxMB: 10}, logger)
	h.UploadFile(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("code=%d, want 400 (handler should propagate request ctx)", w.Code)
	}
}
