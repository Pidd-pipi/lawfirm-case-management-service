# BUG_REPRO：协办律师去重污染共享底层数组

## Bug 是什么

`backend/internal/util/slices.go` 的 `Unique`/`Filter` 使用 `s[:0]` 原地复用底层数组，返回的切片与调用方共享底层数组；`backend/internal/service/case_service.go` 的 `jsonCoLawyers` 未做去重与零值过滤，导致协办律师名单重复保存后被污染串场。

## 如何触发

重复保存案件的协办律师列表，或跑定向测试：

```bash
cd backend
go test ./internal/util -run '^TestUniquePreservesInput$' -count=1
go test ./internal/util -run '^TestFilterPreservesInput$' -count=1
go test ./internal/service -run '^TestJSONCoLawyersDedupes$' -count=1
go test ./internal/service -run '^TestJSONCoLawyersDropsZero$' -count=1
```

## 错误信息

协办律师名单在多次保存后内容串场、出现重复或丢失。
