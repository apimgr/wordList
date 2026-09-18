package mode

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestParseMode(t *testing.T) {
	tests := []struct {
		in      string
		want    Mode
		wantErr bool
	}{
		{"dev", Development, false},
		{"development", Development, false},
		{"prod", Production, false},
		{"production", Production, false},
		{"bogus", "", true},
	}

	for _, tt := range tests {
		got, err := ParseMode(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseMode(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
		}
		if got != tt.want {
			t.Errorf("ParseMode(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestSetAndGet(t *testing.T) {
	if err := Set("development"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if Get() != Development {
		t.Errorf("Get() = %v, want %v", Get(), Development)
	}
	if !IsDevelopment() {
		t.Error("IsDevelopment() = false, want true")
	}
	if IsProduction() {
		t.Error("IsProduction() = true, want false")
	}

	if err := Set("production"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if !IsProduction() {
		t.Error("IsProduction() = false, want true")
	}

	if err := Set("invalid"); err == nil {
		t.Error("Set(\"invalid\") expected error, got nil")
	}
}

func TestInitialize(t *testing.T) {
	// CLI flag takes priority
	if err := Initialize("development"); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	if Get() != Development {
		t.Errorf("Get() = %v, want %v", Get(), Development)
	}

	// Env var used when CLI flag empty
	os.Setenv("MODE", "production")
	defer os.Unsetenv("MODE")
	if err := Initialize(""); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	if Get() != Production {
		t.Errorf("Get() = %v, want %v", Get(), Production)
	}

	// Default when nothing set
	os.Unsetenv("MODE")
	if err := Initialize(""); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	if Get() != Production {
		t.Errorf("Get() = %v, want %v", Get(), Production)
	}

	// Invalid env var propagates error
	os.Setenv("MODE", "nonsense")
	if err := Initialize(""); err == nil {
		t.Error("Initialize() expected error for invalid MODE env var")
	}
}

func TestGetErrorDetail(t *testing.T) {
	if got := GetErrorDetail(nil); got != "" {
		t.Errorf("GetErrorDetail(nil) = %q, want empty", got)
	}

	Set("production")
	err := errors.New("boom")
	got := GetErrorDetail(err)
	if got != "An internal error occurred. Please try again later." {
		t.Errorf("GetErrorDetail() in production = %q", got)
	}

	Set("development")
	got = GetErrorDetail(err)
	if !strings.Contains(got, "boom") || !strings.Contains(got, "Stack trace") {
		t.Errorf("GetErrorDetail() in development missing detail: %q", got)
	}
}

func TestShouldShowDebugEndpoints(t *testing.T) {
	Set("development")
	if !ShouldShowDebugEndpoints() {
		t.Error("expected true in development")
	}
	Set("production")
	if ShouldShowDebugEndpoints() {
		t.Error("expected false in production")
	}
}

func TestGetCacheHeaders(t *testing.T) {
	Set("development")
	h := GetCacheHeaders()
	if h["Cache-Control"] != "no-cache, no-store, must-revalidate" {
		t.Errorf("dev Cache-Control = %q", h["Cache-Control"])
	}

	Set("production")
	h = GetCacheHeaders()
	if h["Cache-Control"] != "public, max-age=31536000, immutable" {
		t.Errorf("prod Cache-Control = %q", h["Cache-Control"])
	}
}

func TestGetLogLevel(t *testing.T) {
	Set("development")
	if GetLogLevel() != "debug" {
		t.Errorf("GetLogLevel() = %q, want debug", GetLogLevel())
	}
	Set("production")
	if GetLogLevel() != "info" {
		t.Errorf("GetLogLevel() = %q, want info", GetLogLevel())
	}
}

func TestShouldCacheTemplatesAndStaticFiles(t *testing.T) {
	Set("production")
	if !ShouldCacheTemplates() || !ShouldCacheStaticFiles() {
		t.Error("expected caching enabled in production")
	}
	Set("development")
	if ShouldCacheTemplates() || ShouldCacheStaticFiles() {
		t.Error("expected caching disabled in development")
	}
}

func TestShouldEnableAutoReloadAndProfiling(t *testing.T) {
	Set("development")
	if !ShouldEnableAutoReload() || !ShouldEnableProfiling() {
		t.Error("expected auto-reload and profiling enabled in development")
	}
	Set("production")
	if ShouldEnableAutoReload() || ShouldEnableProfiling() {
		t.Error("expected auto-reload and profiling disabled in production")
	}
}

func TestGetPanicRecoveryDetail(t *testing.T) {
	Set("production")
	if got := GetPanicRecoveryDetail("boom"); got != "Internal Server Error" {
		t.Errorf("GetPanicRecoveryDetail() production = %q", got)
	}

	Set("development")
	got := GetPanicRecoveryDetail("boom")
	if !strings.Contains(got, "boom") || !strings.Contains(got, "Stack trace") {
		t.Errorf("GetPanicRecoveryDetail() development missing detail: %q", got)
	}
}

func TestModeString(t *testing.T) {
	if Development.String() != "development" {
		t.Errorf("String() = %q", Development.String())
	}
}
