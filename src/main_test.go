package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/apimgr/wordList/src/config"
	"github.com/apimgr/wordList/src/words"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = orig

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintHelp(t *testing.T) {
	out := captureStdout(t, printHelp)
	if !strings.Contains(out, "Usage: wordList") {
		t.Errorf("printHelp() output missing usage: %q", out)
	}
	if !strings.Contains(out, Version) {
		t.Errorf("printHelp() output missing version: %q", out)
	}
}

func TestPrintStartup(t *testing.T) {
	if err := words.Load(); err != nil {
		t.Fatalf("words.Load() error = %v", err)
	}
	cfg := &config.Config{Server: config.ServerConfig{Port: "8080", Address: "0.0.0.0"}}
	out := captureStdout(t, func() { printStartup(cfg) })
	if !strings.Contains(out, "started successfully") {
		t.Errorf("printStartup() output missing banner: %q", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("printStartup() output missing port: %q", out)
	}
}

func TestGetDisplayAddress(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{FQDN: "example.com"}}
	if got := getDisplayAddress(cfg); got != "example.com" {
		t.Errorf("getDisplayAddress() = %q, want example.com", got)
	}

	cfg = &config.Config{Server: config.ServerConfig{Address: "203.0.113.5"}}
	if got := getDisplayAddress(cfg); got != "203.0.113.5" {
		t.Errorf("getDisplayAddress() = %q, want 203.0.113.5", got)
	}

	cfg = &config.Config{Server: config.ServerConfig{Address: "0.0.0.0"}}
	if got := getDisplayAddress(cfg); got == "" {
		t.Error("getDisplayAddress() returned empty string for wildcard address")
	}
}

func TestHandleServiceCommand(t *testing.T) {
	for _, cmd := range []string{"install", "uninstall", "start", "stop", "restart", "status", "unknown"} {
		out := captureStdout(t, func() { handleServiceCommand(cmd) })
		if out == "" {
			t.Errorf("handleServiceCommand(%q) produced no output", cmd)
		}
	}
}

func TestHandleMaintenanceMode(t *testing.T) {
	tests := map[string]string{
		"on":      "ON",
		"off":     "OFF",
		"invalid": "Invalid maintenance mode",
	}
	for in, want := range tests {
		out := captureStdout(t, func() { handleMaintenanceMode(in) })
		if !strings.Contains(out, want) {
			t.Errorf("handleMaintenanceMode(%q) output = %q, want to contain %q", in, out, want)
		}
	}
}
