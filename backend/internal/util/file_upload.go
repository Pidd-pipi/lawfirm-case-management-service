package util

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var allowedExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
}

// SaveUploadedFile 保存上传文件到 uploadDir，返回可访问的相对路径。
func SaveUploadedFile(ctx context.Context, uploadDir string, maxMB int64, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
	if file.Size > maxMB*1024*1024 {
		return "", fmt.Errorf("file too large: %d bytes", file.Size)
	}
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open upload file: %w", err)
	}
	defer src.Close()
	return writeUpload(ctx, uploadDir, ext, src)
}

// writeUpload 将 src 写入 uploadDir 下的新文件。
func writeUpload(ctx context.Context, uploadDir, ext string, src io.Reader) (string, error) {
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, name)
	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("create destination file: %w", err)
	}
	defer out.Close()
	if err := copyWithContext(out, src, ctx); err != nil {
		return "", err
	}
	return "/uploads/" + name, nil
}

// copyWithContext 复制 src 到 dst。
func copyWithContext(dst io.Writer, src io.Reader, ctx context.Context) error {
	_, err := io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("write upload file: %w", err)
	}
	return nil
}
