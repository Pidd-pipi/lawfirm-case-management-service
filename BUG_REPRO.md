# BUG_REPRO：限流器并发竞态

## Bug 是什么

`backend/internal/middleware/rate_limiter.go` 的令牌桶限流存在并发竞态：

- `Limit()` 在读取/刷新 `bucket.tokens`、`bucket.lastFill` 之前就提前 `rl.mu.Unlock()`，并发请求共享同一个 `*bucket` 时在锁外并发读写。
- `Snapshot()` 与 `RateLimitMetrics()` 不加锁遍历内部 `rl.buckets` 并读取 `b.tokens`，与写路径并发产生 data race。
- `Buckets()` 直接返回内部 map 引用，外部修改会泄漏回限流器。

## 如何触发

并发访问登录接口（`POST /api/v1/auth/login` 使用该限流器）时，同时读取限流快照：

```bash
cd backend
go test -race ./internal/middleware -run '^TestRateLimiterConcurrentSnapshot$' -count=1
go test -race ./internal/middleware -run '^TestRateLimitMetricsConcurrent$' -count=1
go test -race ./internal/middleware -run '^TestRateLimiterBucketsIsolated$' -count=1
```

## 错误信息

```
WARNING: DATA RACE
Write at 0x00c000255500 by goroutine 50:
  .../internal/middleware/rate_limiter.go:43
Read at 0x00c00030a0c8 by goroutine 50:
  .../internal/middleware/rate_limiter.go:46
Previous write at 0x00c00030a0c8 by goroutine 45:
  .../internal/middleware/rate_limiter.go:51
```
