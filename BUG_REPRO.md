# BUG_REPRO：业务错误包装断链后统一返回 500

## Bug 是什么

`backend/internal/util/app_error.go` 的 `Wrap` 用 `%v` 包装错误，丢失 `AppError` 错误链；`backend/internal/middleware/error_handler.go` 因此无法用 `errors.As` 识别业务错误，把本应 404/422 的业务错误统一映射成 500。

## 如何触发

任意经过 `Wrap` 包装的业务错误，或跑定向测试：

```bash
cd backend
go test ./internal/util -run '^TestWrapPreservesChain$' -count=1
go test ./internal/util -run '^TestIsAppError$' -count=1
go test ./internal/middleware -run '^TestErrorHandlerClassifiesAppError$' -count=1
go test ./internal/middleware -run '^TestAppErrorStatusNotFound$' -count=1
```

## 错误信息

本应返回 404 的“资源不存在”被包装后返回 HTTP 500。
