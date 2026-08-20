package config

import "testing"

func TestLoadCORSOriginsDefault(t *testing.T) {
	t.Setenv("APP_CORS_ORIGINS", "")
	cfg := Load()
	if cfg.CORSOrigins == nil || len(cfg.CORSOrigins) == 0 {
		t.Fatalf("Load 未为空的 APP_CORS_ORIGINS 提供默认来源，got %#v", cfg.CORSOrigins)
	}
}
