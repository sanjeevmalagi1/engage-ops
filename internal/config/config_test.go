package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadMergesMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "providers.yml", `
providers:
  customerio:
    app_api_key: "secret"
`)
	writeFile(t, dir, "resources.yml", `
resources:
  - provider: customerio
    type: segment
    name: active_users
    attributes:
      name: "Active Users"
`)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if _, ok := cfg.Providers["customerio"]; !ok {
		t.Error("expected customerio provider config to be present")
	}
	if len(cfg.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(cfg.Resources))
	}
	if got := cfg.Resources[0].Address().String(); got != "customerio_segment.active_users" {
		t.Errorf("unexpected resource address: %s", got)
	}
}

func TestLoadRejectsDuplicateAddresses(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "resources.yml", `
resources:
  - provider: customerio
    type: segment
    name: active_users
    attributes: {}
  - provider: customerio
    type: segment
    name: active_users
    attributes: {}
`)

	if _, err := Load(dir); err == nil {
		t.Fatal("expected error for duplicate resource address, got nil")
	}
}

func TestLoadErrorsOnEmptyDir(t *testing.T) {
	dir := t.TempDir()
	if _, err := Load(dir); err == nil {
		t.Fatal("expected error when no YAML files are present, got nil")
	}
}
