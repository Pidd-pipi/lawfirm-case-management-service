package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterConcurrentSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rl := NewRateLimiter(100000)
	engine := gin.New()
	engine.Use(rl.Limit())
	engine.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })
	srv := httptest.NewServer(engine)
	defer srv.Close()

	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				resp, err := http.Get(srv.URL + "/ping")
				if err == nil {
					_ = resp.Body.Close()
				}
			}
		}()
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				_ = rl.Snapshot()
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestRateLimiterBucketsIsolated(t *testing.T) {
	rl := NewRateLimiter(10)
	got := rl.Buckets()
	got["evil"] = &bucket{tokens: 42, lastFill: time.Now()}
	if _, ok := rl.Buckets()["evil"]; ok {
		t.Fatalf("Buckets 返回内部引用，修改泄漏回限流器")
	}
}

func TestRateLimitMetricsConcurrent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rl := NewRateLimiter(100000)
	engine := gin.New()
	engine.Use(rl.Limit())
	engine.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })
	engine.GET("/metrics", RateLimitMetrics(rl))
	srv := httptest.NewServer(engine)
	defer srv.Close()

	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				resp, err := http.Get(srv.URL + "/ping")
				if err == nil {
					_ = resp.Body.Close()
				}
			}
		}()
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 200; j++ {
				resp, err := http.Get(srv.URL + "/metrics")
				if err == nil {
					_ = resp.Body.Close()
				}
			}
		}()
	}
	close(start)
	wg.Wait()
}
