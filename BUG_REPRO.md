# BUG_REPRO：批量删除文档循环 defer 堆积吞错误

## Bug 是什么

`backend/internal/util/batchio.go` 的 `DeleteFiles` 在循环内 `defer` 删除并吞掉错误：循环 defer 堆积到函数尾才执行，且返回的是本地 `err` 而不是命名返回值，删除失败时错误被静默丢弃。`backend/internal/service/document_service.go` 的 `filesForDeletion` 未剥离 `/uploads/` 前缀。

## 如何触发

批量删除包含缺失文件的路径列表，或跑定向测试：

```bash
cd backend
go test ./internal/util -run '^TestDeleteFilesReturnsFirstError$' -count=1
go test ./internal/util -run '^TestDeleteFilesWrapsError$' -count=1
go test ./internal/util -run '^TestDeleteFilesContinuesAfterError$' -count=1
go test ./internal/service -run '^TestFilesForDeletionStripsPrefix$' -count=1
```

## 错误信息

批量删除部分失败时返回 nil 错误，文件句柄/清理动作堆积到函数结束时才执行。
