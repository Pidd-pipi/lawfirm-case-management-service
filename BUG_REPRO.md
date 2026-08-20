# BUG_REPRO：空 CORS 配置触发 nil 切片越界

## Bug 是什么

`backend/internal/config/config.go` 的 `parseCSV` 对空输入返回 nil 切片，`Load()` 对空的 `APP_CORS_ORIGINS` 未提供默认来源；`backend/internal/middleware/cors.go` 与 `origin_guard.go` 直接按下标取值，导致空配置下启动即 panic。

## 如何触发

把 `APP_CORS_ORIGINS` 设为空后启动服务并接收跨域请求，或跑定向测试：

```bash
cd backend
go test ./internal/config -run '^TestParseCSVEmptyReturnsNonNil$' -count=1
go test ./internal/middleware -run '^TestCORSDefaultsWhenEmpty$' -count=1
go test ./internal/middleware -run '^TestOriginGuardAllowsWhenEmpty$' -count=1
```

## 错误信息

```
panic: runtime error: index out of range [0] with length 0

goroutine 1 [running]:
cylawcase/internal/middleware.CORS(...)
	.../internal/middleware/cors.go:13
```
