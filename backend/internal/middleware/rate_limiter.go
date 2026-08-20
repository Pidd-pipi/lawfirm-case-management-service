package middleware

import (
	"net/http"
	"sync"
	"time"

	"cylawcase/internal/constants"

	"github.com/gin-gonic/gin"
)

type bucket struct {
	mu       sync.Mutex
	tokens   int
	lastFill time.Time
}

// RateLimiter 基于 IP 的简易令牌桶限流。
type RateLimiter struct {
	mu       sync.Mutex
	perMin   int
	buckets  map[string]*bucket
	capacity int
}

// NewRateLimiter 构造限流器。
func NewRateLimiter(perMin int) *RateLimiter {
	if perMin <= 0 {
		perMin = 120
	}
	return &RateLimiter{perMin: perMin, buckets: make(map[string]*bucket), capacity: perMin}
}

// Limit 返回限流中间件。
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		b := rl.getBucket(ip)
		if !rl.allow(b, now) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"code": constants.CodeTooManyRequests, "message": constants.MsgTooManyRequests, "data": nil})
			return
		}
		c.Next()
	}
}

// getBucket 取出（或创建）IP 对应的令牌桶，桶自身的状态由桶内互斥锁保护，
// 因此不同 IP 之间互不阻塞。
func (rl *RateLimiter) getBucket(ip string) *bucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if b, ok := rl.buckets[ip]; ok {
		return b
	}
	b := &bucket{tokens: rl.capacity, lastFill: time.Now()}
	rl.buckets[ip] = b
	return b
}

// allow 补充令牌并判断是否放行，对单个桶的状态读写均在桶锁保护下完成。
func (rl *RateLimiter) allow(b *bucket, now time.Time) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	elapsed := now.Sub(b.lastFill)
	if elapsed > 0 {
		b.tokens += int(elapsed.Minutes()) * rl.perMin
		if b.tokens > rl.capacity {
			b.tokens = rl.capacity
		}
		b.lastFill = now
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

// Snapshot 返回各 IP 当前令牌数的快照。
func (rl *RateLimiter) Snapshot() map[string]int {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	out := make(map[string]int, len(rl.buckets))
	for ip, b := range rl.buckets {
		b.mu.Lock()
		out[ip] = b.tokens
		b.mu.Unlock()
	}
	return out
}

// Buckets 返回限流桶供调用方查看。
func (rl *RateLimiter) Buckets() map[string]*bucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	out := make(map[string]*bucket, len(rl.buckets))
	for ip, b := range rl.buckets {
		out[ip] = &bucket{tokens: b.tokens, lastFill: b.lastFill}
	}
	return out
}
