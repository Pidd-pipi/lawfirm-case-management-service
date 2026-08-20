package middleware

import (
	"net/http"

	"cylawcase/internal/config"

	"github.com/gin-gonic/gin"
)

// OriginGuard 校验请求 Origin 是否在允许来源列表内，未配置时使用默认来源。
func OriginGuard(cfg *config.Config) gin.HandlerFunc {
	allowed := cfg.CORSOrigins
	if len(allowed) == 0 {
		allowed = []string{config.DefaultCORSOrigin}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && origin != allowed[0] {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
