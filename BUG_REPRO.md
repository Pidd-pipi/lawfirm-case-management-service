# BUG_REPRO：金额千分位分组共享底层数组污染负号

## Bug 是什么

`backend/internal/util/amount_formatter.go` 的 `groupDigits` 使用包级 `scratchGroups` 并在调用间 `[:0]` 复用，多次格式化会互相污染；负数符号位在分组反转时被交换到错误位置。`backend/internal/util/formatters.go` 的 `FormatMoneyList` 同样复用包级 `moneyBuf`。

## 如何触发

连续格式化多笔金额（尤其负数），或跑定向测试：

```bash
cd backend
go test ./internal/util -run '^TestGroupDigitsIsolated$' -count=1
go test ./internal/util -run '^TestGroupDigitsNegative$' -count=1
go test ./internal/util -run '^TestFormatMoneyListIsolated$' -count=1
```

## 错误信息

负数金额格式化后负号丢失或位置错乱，连续两次格式化结果互相污染。
