package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/apimgr/wordList/src/paths"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()
	if cfg.Server.Port != "64851" {
		t.Errorf("default port = %q", cfg.Server.Port)
	}
	if cfg.Server.Mode != "production" {
		t.Errorf("default mode = %q", cfg.Server.Mode)
	}
	if cfg.WebUI.Theme != "dark" {
		t.Errorf("default theme = %q", cfg.WebUI.Theme)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	paths.Init(dir, dir)

	cfg := defaultConfig()
	cfg.Server.Port = "12345"

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.Server.Port != "12345" {
		t.Errorf("loaded port = %q, want 12345", loaded.Server.Port)
	}
}

func TestLoadCreatesDefaultWhenMissing(t *testing.T) {
	dir := t.TempDir()
	paths.Init(dir, dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Server.Port != "64851" {
		t.Errorf("expected default port, got %q", cfg.Server.Port)
	}

	if _, err := os.Stat(filepath.Join(dir, "server.yml")); err != nil {
		t.Errorf("expected server.yml to be created: %v", err)
	}
}

func TestMigrateYamlToYml(t *testing.T) {
	dir := t.TempDir()
	ymlPath := filepath.Join(dir, "server.yml")
	yamlPath := filepath.Join(dir, "server.yaml")

	if err := os.WriteFile(yamlPath, []byte("server:\n  port: \"9999\"\n"), 0644); err != nil {
		t.Fatal(err)
	}

	migrateYamlToYml(ymlPath)

	if _, err := os.Stat(yamlPath); !os.IsNotExist(err) {
		t.Error("expected .yaml file to be renamed away")
	}
	if _, err := os.Stat(ymlPath); err != nil {
		t.Errorf("expected .yml file to exist: %v", err)
	}

	// non-.yml paths are ignored
	migrateYamlToYml(filepath.Join(dir, "other.txt"))
}
