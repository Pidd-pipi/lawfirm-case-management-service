# BUG_REPRO：JWT 解析错误链断裂

## Bug 是什么

`backend/internal/util/jwt.go` 的 `ParseToken` 用 `%v` 包装解析错误，丢失错误链：

- 过期、畸形、签名不匹配等错误都无法被 `errors.Is` 识别为过期/非法哨兵。
- `backend/internal/middleware/auth.go` 的 `AuthRequired` 对无法识别的解析错误误判为服务器内部错误，返回 500 而不是 401。

## 如何触发

携带已过期或篡改的 JWT 访问任意受保护接口，或直接跑定向测试：

```bash
cd backend
go test ./internal/util -run '^TestParseTokenExpiredErrorIs$' -count=1
go test ./internal/util -run '^TestParseTokenMalformedErrorIs$' -count=1
go test ./internal/middleware -run '^TestAuthRequiredRejectsInvalidToken$' -count=1
```

## 错误信息

过期 token 请求返回 HTTP 500，而正确行为应为 HTTP 401。
