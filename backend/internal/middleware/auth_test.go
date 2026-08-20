package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cylawcase/internal/config"
	"cylawcase/internal/util"

	"github.com/gin-gonic/gin"
)

func TestAuthRequiredRejectsInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{JWTSecret: "test-secret-123456"}
	engine := gin.New()
	engine.Use(AuthRequired(cfg))
	engine.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })
	srv := httptest.NewServer(engine)
	defer srv.Close()

	expired, err := util.GenerateToken("test-secret-123456", -time.Hour, 7, "lawyer", "lawyer")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	for _, tok := range []string{"not-a-token", expired} {
		req, _ := http.NewRequest(http.MethodGet, srv.URL+"/x", nil)
		req.Header.Set("Authorization", "Bearer "+tok)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request error: %v", err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("token=%q status=%d, want 401", tok, resp.StatusCode)
		}
	}
}
