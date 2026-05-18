package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadConfig_Valid(t *testing.T) {
	path := writeConfig(t, `
core:
  address: "https://api.k9.io"
  runas: "nobody"
  connection_timeout: 10
proxy:
  http_listen: ":8080"
  http_mode: "release"
  http_tls: false
  cache_dir: "/tmp/cache"
`)
	cfg := LoadConfig(path)
	if cfg.Core.Address != "https://api.k9.io" {
		t.Fatalf("unexpected address: %q", cfg.Core.Address)
	}
	if cfg.Core.Connection_Timeout != 10 {
		t.Fatalf("expected timeout 10, got %d", cfg.Core.Connection_Timeout)
	}
}

func TestLoadConfig_DefaultTimeout(t *testing.T) {
	path := writeConfig(t, `
core:
  address: "https://api.k9.io"
  runas: "nobody"
proxy:
  http_listen: ":8080"
  http_mode: "release"
  http_tls: false
  cache_dir: "/tmp/cache"
`)
	cfg := LoadConfig(path)
	if cfg.Core.Connection_Timeout != 5 {
		t.Fatalf("expected default timeout 5, got %d", cfg.Core.Connection_Timeout)
	}
}
