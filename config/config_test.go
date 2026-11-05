package config

import "testing"

func TestLoadConfig_AllPresent(t *testing.T) {
	t.Setenv("MONGO_URI", "mongodb://localhost:27017")
	t.Setenv("MONGO_DB_NAME", "quicktask")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("PORT", "9090")

	cfg := LoadConfig()

	if cfg.MongoURI != "mongodb://localhost:27017" {
		t.Errorf("MongoURI = %q", cfg.MongoURI)
	}
	if cfg.MongoDB != "quicktask" {
		t.Errorf("MongoDB = %q", cfg.MongoDB)
	}
	if cfg.JWTSecret != "secret" {
		t.Errorf("JWTSecret = %q", cfg.JWTSecret)
	}
	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want 9090", cfg.Port)
	}
}

func TestLoadConfig_PortFallback(t *testing.T) {
	t.Setenv("MONGO_URI", "mongodb://localhost:27017")
	t.Setenv("MONGO_DB_NAME", "quicktask")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("PORT", "") // unset -> fallback

	cfg := LoadConfig()

	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want default 8080", cfg.Port)
	}
}

func TestGetenv_Fallback(t *testing.T) {
	t.Setenv("SOME_KEY", "")
	if got := getenv("SOME_KEY", "fallback"); got != "fallback" {
		t.Errorf("getenv fallback = %q, want 'fallback'", got)
	}
	t.Setenv("SOME_KEY", "value")
	if got := getenv("SOME_KEY", "fallback"); got != "value" {
		t.Errorf("getenv = %q, want 'value'", got)
	}
}
