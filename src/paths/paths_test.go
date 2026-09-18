package paths

import (
	"os"
	"testing"
)

func resetOverrides() {
	configDir = ""
	dataDir = ""
}

func TestInitOverrides(t *testing.T) {
	defer resetOverrides()

	Init("/tmp/cfg", "/tmp/data")
	if ConfigDir() != "/tmp/cfg" {
		t.Errorf("ConfigDir() = %q, want /tmp/cfg", ConfigDir())
	}
	if DataDir() != "/tmp/data" {
		t.Errorf("DataDir() = %q, want /tmp/data", DataDir())
	}

	// Empty strings should not overwrite existing overrides
	Init("", "")
	if ConfigDir() != "/tmp/cfg" {
		t.Errorf("ConfigDir() after empty Init = %q, want /tmp/cfg", ConfigDir())
	}
	if DataDir() != "/tmp/data" {
		t.Errorf("DataDir() after empty Init = %q, want /tmp/data", DataDir())
	}
}

func TestConfigDirDefault(t *testing.T) {
	defer resetOverrides()
	resetOverrides()

	dir := ConfigDir()
	if dir == "" {
		t.Fatal("ConfigDir() returned empty string")
	}

	if os.Geteuid() == 0 {
		if dir != "/etc/apimgr/wordList" {
			t.Errorf("ConfigDir() as root = %q", dir)
		}
	} else {
		if dir == "/etc/apimgr/wordList" {
			t.Errorf("ConfigDir() as non-root should not be /etc path")
		}
	}
}

func TestDataDirDefault(t *testing.T) {
	defer resetOverrides()
	resetOverrides()

	dir := DataDir()
	if dir == "" {
		t.Fatal("DataDir() returned empty string")
	}

	if os.Geteuid() == 0 {
		if dir != "/var/lib/apimgr/wordList" {
			t.Errorf("DataDir() as root = %q", dir)
		}
	}
}

func TestLogDir(t *testing.T) {
	defer resetOverrides()
	resetOverrides()

	dir := LogDir()
	if dir == "" {
		t.Fatal("LogDir() returned empty string")
	}

	if os.Geteuid() == 0 {
		if dir != "/var/log/apimgr/wordList" {
			t.Errorf("LogDir() as root = %q", dir)
		}
	}
}
