# BUG_REPRO：上传未响应取消信号残留半截文件

## Bug 是什么

`backend/internal/util/file_upload.go` 的 `SaveUploadedFile` 用 `io.Copy` 直接写盘，不接收也不检查 `context`；`backend/internal/handler/upload_handler.go` 传入 `context.Background()` 而不是请求 ctx。客户端取消后服务端仍继续写完整文件，且失败时不清除已创建的半截文件。

## 如何触发

上传较大文件后立即取消请求，或跑定向测试：

```bash
cd backend
go test ./internal/util -run '^TestSaveUploadedFileCanceledNoDir$' -count=1
go test ./internal/util -run '^TestCopyWithContextStopsOnCancel$' -count=1
go test ./internal/util -run '^TestWriteUploadRemovesPartialOnError$' -count=1
go test ./internal/handler -run '^TestUploadHandlerUsesRequestContext$' -count=1
```

## 错误信息

取消后上传目录中仍残留未完成的半截文件。
