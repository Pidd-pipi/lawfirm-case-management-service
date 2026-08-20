package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"cylawcase/internal/constants"
	"cylawcase/internal/util"

	"github.com/gin-gonic/gin"
)

func TestErrorHandlerClassifiesAppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	engine := gin.New()
	engine.Use(ErrorHandler(logger))
	engine.GET("/x", func(c *gin.Context) {
		_ = c.Error(util.NewAppError(constants.CodeNotFound, constants.MsgNotFound))
	})
	srv := httptest.NewServer(engine)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/x")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status=%d, want 404", resp.StatusCode)
	}
}

func TestAppErrorStatusNotFound(t *testing.T) {
	if got := appErrorStatus(constants.CodeNotFound); got != http.StatusNotFound {
		t.Fatalf("appErrorStatus(CodeNotFound)=%d, want 404", got)
	}
}
