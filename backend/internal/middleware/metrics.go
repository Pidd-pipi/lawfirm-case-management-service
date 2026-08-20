package middleware

import "github.com/gin-gonic/gin"

// RateLimitMetrics 返回限流器各 IP 剩余令牌总数，供运维观察。
func RateLimitMetrics(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		total := 0
		for _, b := range rl.buckets {
			total += b.tokens
		}
		c.JSON(200, gin.H{"total_tokens": total})
	}
}
