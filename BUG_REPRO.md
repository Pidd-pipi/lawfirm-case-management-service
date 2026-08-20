# BUG_REPRO：账单状态机允许 pending 直接开票

## Bug 是什么

`backend/internal/service/billing_service.go` 的 `billingStatusCanFlow` 转换表漏判：pending 账单可绕过 paid 直接进入 invoiced；void 可再转回其他状态；`backend/internal/constants/billing.go` 的状态枚举漏掉 invoiced。

## 如何触发

对 pending 账单调用开票接口，或跑定向测试：

```bash
cd backend
go test ./internal/constants -run '^TestIsValidBillingStatusInvoiced$' -count=1
go test ./internal/constants -run '^TestBillingStatusTransitions$' -count=1
go test ./internal/service -run '^TestBillingStatusCanFlowInvoiced$' -count=1
go test ./internal/service -run '^TestBillingStatusCanFlowPaid$' -count=1
```

## 错误信息

pending 账单被直接标记为已开票，跳过了已支付状态。
