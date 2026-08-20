package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cylawcase/internal/config"

	"github.com/gin-gonic/gin"
)

func TestCORSDefaultsWhenEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{CORSOrigins: nil}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CORS panicked on empty config: %v", r)
		}
	}()
	_ = CORS(cfg)
}

func TestOriginGuardAllowsWhenEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{CORSOrigins: nil}
	engine := gin.New()
	engine.Use(OriginGuard(cfg))
	engine.GET("/x", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	srv := httptest.NewServer(engine)
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/x", nil)
	req.Header.Set("Origin", "http://localhost:28031")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d, want 200", resp.StatusCode)
	}
}
