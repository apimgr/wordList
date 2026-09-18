package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apimgr/wordList/src/config"
	"github.com/apimgr/wordList/src/words"
)

func testServer(t *testing.T) http.Handler {
	t.Helper()
	if err := words.Load(); err != nil {
		t.Fatalf("words.Load() error = %v", err)
	}
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:    "64851",
			Address: "0.0.0.0",
		},
		WebUI: config.WebUIConfig{Theme: "dark"},
		WebRobots: config.WebRobotsConfig{
			Allow: []string{"/", "/api"},
			Deny:  []string{"/admin"},
		},
		WebSecurity: config.WebSecurityConfig{
			Admin: "security@example.com",
			CORS:  "*",
		},
	}
	return New(cfg, "test").Handler
}

func doRequest(t *testing.T, h http.Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestHealthHandler(t *testing.T) {
	h := testServer(t)
	rec := doRequest(t, h, http.MethodGet, "/healthz")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body["status"] != "healthy" {
		t.Errorf("status field = %v", body["status"])
	}
}

func TestRobotsHandler(t *testing.T) {
	h := testServer(t)
	rec := doRequest(t, h, http.MethodGet, "/robots.txt")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Allow: /") || !strings.Contains(body, "Disallow: /admin") {
		t.Errorf("unexpected robots.txt body: %q", body)
	}
}

func TestSecurityHandler(t *testing.T) {
	h := testServer(t)
	for _, path := range []string{"/security.txt", "/.well-known/security.txt"} {
		rec := doRequest(t, h, http.MethodGet, path)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d", path, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "security@example.com") {
			t.Errorf("%s missing contact email: %q", path, rec.Body.String())
		}
	}
}

func TestManifestHandler(t *testing.T) {
	h := testServer(t)
	rec := doRequest(t, h, http.MethodGet, "/manifest.json")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/manifest+json" {
		t.Errorf("Content-Type = %q", ct)
	}
}

func TestSwaggerRedirect(t *testing.T) {
	h := testServer(t)
	rec := doRequest(t, h, http.MethodGet, "/swagger")
	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/api" {
		t.Errorf("Location = %q", loc)
	}
}

func TestAPIRandomHandler(t *testing.T) {
	h := testServer(t)

	rec := doRequest(t, h, http.MethodGet, "/api/v1/random")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var body struct {
		Count int      `json:"count"`
		Words []string `json:"words"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body.Count != 6 {
		t.Errorf("default count = %d, want 6", body.Count)
	}

	rec = doRequest(t, h, http.MethodGet, "/api/v1/random/3")
	json.Unmarshal(rec.Body.Bytes(), &body)
	if body.Count != 3 {
		t.Errorf("count = %d, want 3", body.Count)
	}

	rec = doRequest(t, h, http.MethodGet, "/api/v1/random/3.txt")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("Content-Type = %q", ct)
	}
}

func TestAPIPassphraseHandler(t *testing.T) {
	h := testServer(t)

	rec := doRequest(t, h, http.MethodGet, "/api/v1/passphrase/4?separator=_&capitalize=false")
	var body struct {
		Passphrase string `json:"passphrase"`
		WordCount  int    `json:"word_count"`
	}
	json.Unmarshal(rec.Body.Bytes(), &body)
	if body.WordCount != 4 {
		t.Errorf("word_count = %d, want 4", body.WordCount)
	}
	if !strings.Contains(body.Passphrase, "_") {
		t.Errorf("passphrase does not use custom separator: %q", body.Passphrase)
	}

	rec = doRequest(t, h, http.MethodGet, "/api/v1/passphrase.txt")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestAPISearchHandler(t *testing.T) {
	h := testServer(t)

	all := words.All()
	if len(all) == 0 {
		t.Fatal("no words loaded")
	}
	prefix := all[0][:1]

	rec := doRequest(t, h, http.MethodGet, "/api/v1/search/"+prefix)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var body struct {
		Count int `json:"count"`
	}
	json.Unmarshal(rec.Body.Bytes(), &body)
	if body.Count == 0 {
		t.Error("expected search results")
	}

	rec = doRequest(t, h, http.MethodGet, "/api/v1/search/"+prefix+"?mode=contains&limit=1")
	json.Unmarshal(rec.Body.Bytes(), &body)
	if body.Count > 1 {
		t.Errorf("expected limit to cap results, got %d", body.Count)
	}

	rec = doRequest(t, h, http.MethodGet, "/api/v1/search/"+prefix+".txt")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestAPIAllWordsHandler(t *testing.T) {
	h := testServer(t)

	rec := doRequest(t, h, http.MethodGet, "/api/v1/words?page=1&per_page=10")
	var body struct {
		Words      []string `json:"words"`
		TotalPages int      `json:"total_pages"`
	}
	json.Unmarshal(rec.Body.Bytes(), &body)
	if len(body.Words) != 10 {
		t.Errorf("expected 10 words, got %d", len(body.Words))
	}

	rec = doRequest(t, h, http.MethodGet, "/api/v1/words.txt")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestAPIByLetterHandler(t *testing.T) {
	h := testServer(t)

	rec := doRequest(t, h, http.MethodGet, "/api/v1/letter/a")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	rec = doRequest(t, h, http.MethodGet, "/api/v1/letter/a.txt")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestAPIByLengthHandler(t *testing.T) {
	h := testServer(t)

	rec := doRequest(t, h, http.MethodGet, "/api/v1/length/5")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	rec = doRequest(t, h, http.MethodGet, "/api/v1/length/notanumber")
	var body struct {
		Error string `json:"error"`
	}
	json.Unmarshal(rec.Body.Bytes(), &body)
	if body.Error == "" {
		t.Error("expected error for invalid length")
	}

	rec = doRequest(t, h, http.MethodGet, "/api/v1/length/notanumber.txt")
	if !strings.Contains(rec.Body.String(), "Invalid length") {
		t.Errorf("expected invalid length text, got %q", rec.Body.String())
	}
}

func TestAPIStatsLettersCountHandlers(t *testing.T) {
	h := testServer(t)

	for _, path := range []string{"/api/v1/stats", "/api/v1/letters", "/api/v1/count", "/api/v1/count.txt"} {
		rec := doRequest(t, h, http.MethodGet, path)
		if rec.Code != http.StatusOK {
			t.Errorf("%s status = %d", path, rec.Code)
		}
	}
}

func TestAPISearchMissingQuery(t *testing.T) {
	h := testServer(t)
	// chi requires a non-empty path segment, so this exercises the empty-query branch indirectly
	rec := doRequest(t, h, http.MethodGet, "/api/v1/search/%20")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestWebPageHandlers(t *testing.T) {
	h := testServer(t)
	for _, path := range []string{"/", "/generate", "/search", "/api"} {
		rec := doRequest(t, h, http.MethodGet, path)
		if rec.Code != http.StatusOK {
			t.Errorf("%s status = %d", path, rec.Code)
		}
	}
}
