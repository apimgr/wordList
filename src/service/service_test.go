package service

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDetectServiceManager(t *testing.T) {
	got := DetectServiceManager()
	switch runtime.GOOS {
	case "linux":
		if got != ServiceSystemd && got != ServiceRunit && got != ServiceUnknown {
			t.Errorf("unexpected service type on linux: %v", got)
		}
	case "darwin":
		if got != ServiceLaunchd {
			t.Errorf("expected ServiceLaunchd on darwin, got %v", got)
		}
	case "windows":
		if got != ServiceWindows {
			t.Errorf("expected ServiceWindows on windows, got %v", got)
		}
	}
}

func TestGetBinaryPath(t *testing.T) {
	path := GetBinaryPath()
	if path == "" {
		t.Fatal("GetBinaryPath() returned empty string")
	}
	if runtime.GOOS != "windows" && path != "/usr/local/bin/wordList" {
		t.Errorf("GetBinaryPath() = %q", path)
	}
}

func TestCopyBinary(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src-bin")
	dst := filepath.Join(dir, "nested", "dst-bin")

	if err := os.WriteFile(src, []byte("binary-content"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := copyBinary(src, dst); err != nil {
		t.Fatalf("copyBinary() error = %v", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read copied binary: %v", err)
	}
	if string(data) != "binary-content" {
		t.Errorf("copied content = %q", data)
	}

	if err := copyBinary(filepath.Join(dir, "does-not-exist"), dst); err == nil {
		t.Error("expected error copying missing source file")
	}
}
