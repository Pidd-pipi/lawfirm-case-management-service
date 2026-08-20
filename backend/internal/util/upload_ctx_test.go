package util

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeUploadFile(t *testing.T) *multipart.FileHeader {
	t.Helper()
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	part, err := mw.CreateFormFile("file", "report.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("hello pdf")); err != nil {
		t.Fatal(err)
	}
	if err := mw.Close(); err != nil {
		t.Fatal(err)
	}
	mr := multipart.NewReader(&b, mw.Boundary())
	form, err := mr.ReadForm(1 << 20)
	if err != nil {
		t.Fatal(err)
	}
	return form.File["file"][0]
}

func TestSaveUploadedFileCanceledNoDir(t *testing.T) {
	fh := makeUploadFile(t)
	dir := filepath.Join(t.TempDir(), "uploads")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := SaveUploadedFile(ctx, dir, 10, fh)
	if err == nil {
		t.Fatal("expected cancel error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if _, statErr := os.Stat(dir); !os.IsNotExist(statErr) {
		t.Fatalf("upload dir should not be created on cancel, stat err=%v", statErr)
	}
}

type slowReader struct {
	remaining int
	sleep     time.Duration
}

func (r *slowReader) Read(p []byte) (int, error) {
	if r.remaining <= 0 {
		return 0, io.EOF
	}
	r.remaining--
	time.Sleep(r.sleep)
	copy(p, []byte("chunkdata"))
	return 8, nil
}

func TestCopyWithContextStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	r := &slowReader{remaining: 500, sleep: 2 * time.Millisecond}
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()
	var buf bytes.Buffer
	err := copyWithContext(&buf, r, ctx)
	if err == nil {
		t.Fatal("expected cancel error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

type failingReader struct {
	data []byte
	at   int
}

func (r *failingReader) Read(p []byte) (int, error) {
	if r.at >= len(r.data) {
		return 0, errors.New("boom")
	}
	n := copy(p, r.data[r.at:])
	r.at += n
	return n, nil
}

func TestWriteUploadRemovesPartialOnError(t *testing.T) {
	dir := t.TempDir()
	_, err := writeUpload(context.Background(), dir, ".pdf", &failingReader{data: bytes.Repeat([]byte("x"), 1024)})
	if err == nil {
		t.Fatal("expected read error")
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Fatalf("partial file not removed, entries=%v", entries)
	}
}
