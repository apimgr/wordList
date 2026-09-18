package ssl

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestNewManagerAndDisabledTLSConfig(t *testing.T) {
	m := NewManager(Config{Enabled: false})
	cfg, err := m.GetTLSConfig([]string{"example.com"})
	if err != nil {
		t.Fatalf("GetTLSConfig() error = %v", err)
	}
	if cfg != nil {
		t.Error("expected nil TLS config when SSL disabled")
	}
}

func TestGetTLSConfigNoCertsAvailable(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(Config{Enabled: true, CertPath: dir})
	_, err := m.GetTLSConfig([]string{"example.com"})
	if err == nil {
		t.Error("expected error when no certificates and Let's Encrypt disabled")
	}
}

func TestFindManualCerts(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(Config{CertPath: dir})

	cert, key := m.findManualCerts([]string{"example.com"})
	if cert != "" || key != "" {
		t.Fatalf("expected no certs found, got %q %q", cert, key)
	}

	crt := filepath.Join(dir, "example.com.crt")
	k := filepath.Join(dir, "example.com.key")
	if err := os.WriteFile(crt, []byte("cert"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(k, []byte("key"), 0644); err != nil {
		t.Fatal(err)
	}

	cert, key = m.findManualCerts([]string{"example.com"})
	if cert != crt || key != k {
		t.Errorf("findManualCerts() = %q, %q", cert, key)
	}

	m2 := NewManager(Config{})
	cert, key = m2.findManualCerts([]string{"example.com"})
	if cert != "" || key != "" {
		t.Error("expected empty result when CertPath unset")
	}
}

func TestFindExistingCerts(t *testing.T) {
	m := NewManager(Config{})
	cert, key := m.findExistingCerts([]string{"nonexistent.invalid"})
	if cert != "" || key != "" {
		t.Errorf("expected no existing certs, got %q %q", cert, key)
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "exists.txt")
	if err := os.WriteFile(f, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if !fileExists(f) {
		t.Error("expected fileExists() true for existing file")
	}
	if fileExists(filepath.Join(dir, "missing.txt")) {
		t.Error("expected fileExists() false for missing file")
	}
}

func TestGetHTTPHandlerFallback(t *testing.T) {
	m := NewManager(Config{})
	called := false
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	h := m.GetHTTPHandler(fallback)
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if !called {
		t.Error("expected fallback handler to be called when certManager is nil")
	}
}

func TestParseChallenge(t *testing.T) {
	tests := map[string]string{
		"http-01":     "http-01",
		"HTTP01":      "http-01",
		"http":        "http-01",
		"tls-alpn-01": "tls-alpn-01",
		"tlsalpn01":   "tls-alpn-01",
		"tls":         "tls-alpn-01",
		"dns-01":      "dns-01",
		"dns01":       "dns-01",
		"dns":         "dns-01",
		"  DNS  ":     "dns-01",
		"unknown":     "http-01",
	}
	for in, want := range tests {
		if got := ParseChallenge(in); got != want {
			t.Errorf("ParseChallenge(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestChallengeServer(t *testing.T) {
	cs := NewChallengeServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/some/other/path", nil)
	if cs.ServeHTTP(rec, req) {
		t.Error("expected ServeHTTP to return false for non-challenge path")
	}

	cs.SetToken("tok1", "auth1")

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/tok1", nil)
	if !cs.ServeHTTP(rec, req) {
		t.Fatal("expected ServeHTTP to return true for challenge path")
	}
	if rec.Body.String() != "auth1" {
		t.Errorf("challenge body = %q, want auth1", rec.Body.String())
	}

	cs.ClearToken("tok1")
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/tok1", nil)
	if !cs.ServeHTTP(rec, req) {
		t.Fatal("expected ServeHTTP to return true even for unknown token")
	}
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404 for cleared token", rec.Code)
	}
}
