package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfig_Success(t *testing.T) {

	tempDir := t.TempDir()
	configJSON := `{
		"APP_ENV": "test",
		"HTTP_ADDR": ":8080",
		"DATABASE_URL": "postgres://localhost/test",
		"REQUEST_TIMEOUT": "30s",
		"LOGLEVEL": "debug"
	}`

	configPath := filepath.Join(tempDir, "config.json")
	if err := os.WriteFile(configPath, []byte(configJSON), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	withWorkingDir(t, tempDir)

	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected config, got nil")
	}

	if cfg.App_Env != "test" {
		t.Errorf("App_Env = %q, want %q", cfg.App_Env, "test")
	}
	if cfg.Http_Addr != ":8080" {
		t.Errorf("Http_Addr = %q, want %q", cfg.Http_Addr, ":8080")
	}
	if cfg.Database_URL != "postgres://localhost/test" {
		t.Errorf("Database_URL = %q, want %q", cfg.Database_URL, "postgres://localhost/test")
	}
	if cfg.Request_Timeout != "30s" {
		t.Errorf("Request_Timeout = %q, want %q", cfg.Request_Timeout, "30s")
	}
	if cfg.Log_Level != "debug" {
		t.Errorf("Log_Level = %q, want %q", cfg.Log_Level, "debug")
	}
}

func TestGetConfig_MissingFile(t *testing.T) {

	tempDir := t.TempDir()
	withWorkingDir(t, tempDir)

	cfg, err := GetConfig()
	if err == nil {
		t.Fatal("expected error for missing config file, got nil")
	}
	if cfg != nil {
		t.Fatalf("expected nil config when file missing, got %#v", cfg)
	}
}

func TestGetConfig_InvalidJSON(t *testing.T) {

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	if err := os.WriteFile(configPath, []byte(`{"APP_ENV":`), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	withWorkingDir(t, tempDir)

	cfg, err := GetConfig()
	if err == nil {
		t.Fatal("expected error for invalid config JSON, got nil")
	}
	if cfg != nil {
		t.Fatalf("expected nil config for invalid JSON, got %#v", cfg)
	}
}

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	})
}
