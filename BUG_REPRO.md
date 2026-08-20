# BUG_REPRO：空 accept_date 返回零值时间指针

## Bug 是什么

`backend/internal/dto/dto_case.go` 的 `ParseAcceptDate` 对空串返回指向零值 `time.Time` 的非 nil 指针，非法日期也吞掉错误返回零值；`ValidateAcceptDate` 对零值放行；`backend/internal/dto/dto_common.go` 的 `PageQuery.Normalize` 对 nil 接收者无保护会 panic。

## 如何触发

创建案件时 `accept_date` 留空，或跑定向测试：

```bash
cd backend
go test ./internal/dto -run '^TestParseAcceptDateEmptyNil$' -count=1
go test ./internal/dto -run '^TestParseAcceptDateInvalidError$' -count=1
go test ./internal/dto -run '^TestValidateAcceptDateZero$' -count=1
go test ./internal/dto -run '^TestPageQueryNormalizeNilReceiver$' -count=1
```

## 错误信息

留空的 accept_date 被写入为零值时间，分页参数为 nil 时触发 panic。
