package server

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apimgr/wordList/src/config"
	"github.com/apimgr/wordList/src/words"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

//go:embed openapi.json
var openAPISpec []byte

var templates *template.Template

// New creates a new HTTP server
func New(cfg *config.Config, version string) *http.Server {
	// Parse templates
	var err error
	templates, err = template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		panic(fmt.Sprintf("Failed to parse templates: %v", err))
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(securityHeaders)
	r.Use(requestBodyLimit(1 << 20))
	if cfg.Server.RateLimit.Enabled {
		r.Use(newRateLimitMiddleware(cfg.Server.RateLimit.Requests, time.Duration(cfg.Server.RateLimit.Window)*time.Second))
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.WebSecurity.CORS},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Web routes
	r.Get("/", homeHandler(cfg))
	r.Get("/generate", generatePageHandler(cfg))
	r.Post("/generate", generatePageHandler(cfg))
	r.Get("/search", searchPageHandler(cfg))
	r.Get("/server/about", infoPageHandler(cfg, "About", "WordList serves the EFF Large Wordlist for secure, memorable passphrases."))
	r.Get("/server/help", infoPageHandler(cfg, "Help", "Generate a passphrase in the browser or use the public REST API and wordList-cli."))
	r.Get("/server/privacy", infoPageHandler(cfg, "Privacy", "Generated passphrases are returned only to the requesting client and are never stored by this service. Access logs omit request and response bodies."))
	r.Get("/server/terms", infoPageHandler(cfg, "Terms", "The service provides the EFF Large Wordlist and passphrase generation for public use. Respect the EFF source attribution shown in this repository."))
	r.Get("/server/healthz", healthPageHandler(cfg, version))
	r.Get("/server/docs/swagger", apiDocsHandler(cfg))
	r.Get("/server/docs/graphql", graphqlDocsHandler(cfg))

	// Special files
	r.Get("/robots.txt", robotsHandler(cfg))
	r.Get("/security.txt", securityHandler(cfg))
	r.Get("/.well-known/security.txt", securityHandler(cfg))
	r.Get("/manifest.json", manifestHandler(cfg))
	r.Get("/sw.js", serviceWorkerHandler)
	r.Get("/api/swagger", swaggerJSONHandler)
	r.Get("/api/graphql", graphqlHandler)
	r.Post("/api/graphql", graphqlHandler)
	r.Get("/api/healthz", healthHandler(version))
	r.Get("/graphql", graphqlHandler)
	r.Post("/graphql", graphqlHandler)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", apiRootHandler)
		r.Get("/server/healthz", healthHandler(version))
		r.Get("/server/about", apiInfoHandler("about", "WordList serves the EFF Large Wordlist for secure, memorable passphrases."))
		r.Get("/server/help", apiInfoHandler("help", "Generate a passphrase in the browser or use the public REST API and wordList-cli."))
		r.Get("/server/privacy", apiInfoHandler("privacy", "Generated passphrases are returned only to the requesting client and are never stored by this service."))
		r.Get("/server/terms", apiInfoHandler("terms", "The service provides the EFF Large Wordlist and passphrase generation for public use."))
		r.Get("/server/swagger", swaggerJSONHandler)
		r.Get("/server/graphql", graphqlHandler)
		r.Post("/server/graphql", graphqlHandler)
		r.Get("/graphql", graphqlHandler)
		r.Post("/graphql", graphqlHandler)

		// Random words
		r.Get("/random", apiRandomHandler)
		r.Get("/random/{count}", apiRandomHandler)
		r.Get("/random.txt", apiRandomTextHandler)
		r.Get("/random/{count}.txt", apiRandomTextHandler)

		// Passphrase generation
		r.Get("/passphrase", apiPassphraseHandler)
		r.Get("/passphrase/{count}", apiPassphraseHandler)
		r.Get("/passphrase.txt", apiPassphraseTextHandler)
		r.Get("/passphrase/{count}.txt", apiPassphraseTextHandler)

		// Search
		r.Get("/search/{query}", apiSearchHandler)
		r.Get("/search/{query}.txt", apiSearchTextHandler)

		// Browse
		r.Get("/words", apiAllWordsHandler)
		r.Get("/words.txt", apiAllWordsTextHandler)
		r.Get("/dice/{code}", apiByDiceCodeHandler)
		r.Get("/dice/{code}.txt", apiByDiceCodeTextHandler)
		r.Get("/letters/{letter}", apiByLetterHandler)
		r.Get("/letters/{letter}.txt", apiByLetterTextHandler)
		r.Get("/lengths/{length}", apiByLengthHandler)
		r.Get("/lengths/{length}.txt", apiByLengthTextHandler)

		// Stats
		r.Get("/stats", apiStatsHandler)
		r.Get("/letters", apiLettersHandler)
		r.Get("/count", apiCountHandler)
		r.Get("/count.txt", apiCountTextHandler)
	})

	return &http.Server{
		Addr:              net.JoinHostPort(strings.Trim(cfg.Server.Address, "[]"), cfg.Server.Port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 16,
	}
}

// Template data
type PageData struct {
	SiteTitle string
	BaseURL   string
	Theme     string
	Title     string
	Content   string
	Passphrase string
	WordCount int
	Entropy float64
	Words []string
	Query string
	Results []string
	Error string
}

func newPageData(cfg *config.Config) PageData {
	return PageData{
		SiteTitle: "WordList",
		BaseURL:   func() string { if cfg.Server.FQDN != "" { return "https://" + cfg.Server.FQDN }; return "http://localhost:" + cfg.Server.Port }(),
		Theme:     cfg.WebUI.Theme,
	}
}

// Web handlers
func homeHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newPageData(cfg)
		_ = templates.ExecuteTemplate(w, "home.html", data)
	}
}

func generatePageHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newPageData(cfg)
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "invalid form submission", http.StatusBadRequest)
				return
			}
			count, err := strconv.Atoi(r.FormValue("count"))
			if err != nil || count < 1 || count > 20 {
				data.Error = "Choose between 1 and 20 words."
			} else {
				separator := r.FormValue("separator")
				if !r.PostForm.Has("separator") {
					separator = "-"
				}
				if len(separator) > 8 {
					data.Error = "Separator must be at most 8 characters."
				} else {
					data.Passphrase, err = words.PassphraseSecure(count, separator, r.FormValue("capitalize") != "false")
					if err != nil {
						data.Error = "Secure random generation is temporarily unavailable. Please retry."
					} else {
						data.WordCount = count
						data.Entropy = words.EntropyBits(count)
						data.Words = strings.Split(data.Passphrase, separator)
					}
				}
			}
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		_ = templates.ExecuteTemplate(w, "generate.html", data)
	}
}

func searchPageHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newPageData(cfg)
		data.Query = strings.TrimSpace(r.URL.Query().Get("q"))
		if len(data.Query) > 100 {
			data.Error = "Search query must be 100 characters or fewer."
		} else if data.Query != "" {
			data.Results = words.Search(data.Query)
			if len(data.Results) > 100 {
				data.Results = data.Results[:100]
			}
		}
		_ = templates.ExecuteTemplate(w, "search.html", data)
	}
}

func apiDocsHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newPageData(cfg)
		_ = templates.ExecuteTemplate(w, "openapi.html", data)
	}
}

func infoPageHandler(cfg *config.Config, title, content string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newPageData(cfg)
		data.Title = title
		data.Content = content
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = templates.ExecuteTemplate(w, "info.html", data)
	}
}

func healthPageHandler(cfg *config.Config, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newPageData(cfg)
		data.Title = "Service health"
		data.Content = fmt.Sprintf("Healthy · %d EFF words · version %s", words.Count(), version)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = templates.ExecuteTemplate(w, "info.html", data)
	}
}

func healthHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, map[string]interface{}{
			"status": "healthy", "timestamp": time.Now().UTC().Format(time.RFC3339),
			"word_count": words.Count(), "service": "wordList", "version": version,
		})
	}
}

func graphqlDocsHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newPageData(cfg)
		data.Title = "GraphQL"
		data.Content = "GraphQL endpoint documentation is served locally. Submit a JSON request containing a query to /api/v1/server/graphql."
		_ = templates.ExecuteTemplate(w, "info.html", data)
	}
}

func apiRootHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, map[string]interface{}{"service": "wordList", "version": "v1", "word_count": words.Count(), "endpoints": []string{"/api/v1/passphrase", "/api/v1/words", "/api/v1/dice/{code}", "/api/v1/server/healthz", "/api/v1/server/swagger", "/api/v1/server/graphql"}})
}

func robotsHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "User-agent: *")
		for _, path := range cfg.WebRobots.Allow {
			fmt.Fprintf(w, "Allow: %s\n", path)
		}
		for _, path := range cfg.WebRobots.Deny {
			fmt.Fprintf(w, "Disallow: %s\n", path)
		}
	}
}

func securityHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Contact: mailto:%s\n", cfg.WebSecurity.Admin)
		fmt.Fprintln(w, "Preferred-Languages: en")
	}
}

func manifestHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/manifest+json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":             "WordList",
			"short_name":       "WordList",
			"description":      "EFF Large Wordlist API for Passphrases",
			"start_url":        "/",
			"display":          "standalone",
			"background_color": "#1e1e2e",
			"theme_color":      "#6366f1",
			"icons": []map[string]string{
				{"src": "/static/icon.svg", "sizes": "any", "type": "image/svg+xml", "purpose": "any maskable"},
			},
		})
	}
}

// API response helper
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func textResponse(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(text))
}

// API handlers

func requestCount(w http.ResponseWriter, r *http.Request, key string, fallback, minimum, maximum int) (int, bool) {
	value := chi.URLParam(r, key)
	if value == "" {
		value = r.URL.Query().Get(key)
	}
	if value == "" {
		return fallback, true
	}
	count, err := strconv.Atoi(value)
	if err != nil || count < minimum || count > maximum {
		writeAPIError(w, http.StatusBadRequest, fmt.Sprintf("count must be between %d and %d", minimum, maximum))
		return 0, false
	}
	return count, true
}

func writeAPIError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"type": "about:blank", "title": http.StatusText(status), "status": status, "detail": message})
}

func writeTextError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message + "\n"))
}

func apiRandomHandler(w http.ResponseWriter, r *http.Request) {
	count, ok := requestCount(w, r, "count", 6, 1, 100)
	if !ok {
		return
	}
	randomWords, err := words.RandomSecure(count)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "secure random source unavailable")
		return
	}
	jsonResponse(w, map[string]interface{}{
		"count": len(randomWords),
		"words": randomWords,
	})
}

func apiRandomTextHandler(w http.ResponseWriter, r *http.Request) {
	count, ok := requestCount(w, r, "count", 6, 1, 100)
	if !ok {
		return
	}
	randomWords, err := words.RandomSecure(count)
	if err != nil {
		writeTextError(w, http.StatusInternalServerError, "secure random source unavailable")
		return
	}
	textResponse(w, strings.Join(randomWords, "\n"))
}

func apiPassphraseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	count, ok := requestCount(w, r, "count", 6, 1, 20)
	if !ok {
		return
	}

	separator := "-"
	if r.URL.Query().Has("separator") {
		separator = r.URL.Query().Get("separator")
	}
	if len(separator) > 8 {
		writeAPIError(w, http.StatusBadRequest, "separator must be at most 8 characters")
		return
	}

	capitalize := r.URL.Query().Get("capitalize") != "false"

	passphrase, err := words.PassphraseSecure(count, separator, capitalize)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "secure random source unavailable")
		return
	}
	passphraseWords := strings.Split(passphrase, separator)

	jsonResponse(w, map[string]interface{}{
		"passphrase":  passphrase,
		"entropy_bits": words.EntropyBits(count),
		"words":       passphraseWords,
		"word_count":  count,
		"separator":   separator,
		"capitalized": capitalize,
	})
}

func apiPassphraseTextHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	count, ok := requestCount(w, r, "count", 6, 1, 20)
	if !ok {
		return
	}

	separator := "-"
	if r.URL.Query().Has("separator") {
		separator = r.URL.Query().Get("separator")
	}
	if len(separator) > 8 {
		writeAPIError(w, http.StatusBadRequest, "separator must be at most 8 characters")
		return
	}

	capitalize := r.URL.Query().Get("capitalize") != "false"

	passphrase, err := words.PassphraseSecure(count, separator, capitalize)
	if err != nil {
		writeTextError(w, http.StatusInternalServerError, "secure random source unavailable")
		return
	}
	textResponse(w, passphrase)
}

func apiSearchHandler(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	if query == "" || len(query) > 100 {
		writeAPIError(w, http.StatusBadRequest, "query must contain 1 to 100 characters")
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "prefix"
	}
	if mode != "prefix" && mode != "contains" {
		writeAPIError(w, http.StatusBadRequest, "mode must be prefix or contains")
		return
	}
	limit, err := searchLimit(r)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	var results []string

	if mode == "contains" {
		results = words.SearchContains(query)
	} else {
		results = words.Search(query)
	}

	if len(results) > limit {
		results = results[:limit]
	}

	jsonResponse(w, map[string]interface{}{
		"query":   query,
		"mode":    mode,
		"count":   len(results),
		"words":   results,
		"limited": len(results) == limit,
	})
}

func apiSearchTextHandler(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	if query == "" || len(query) > 100 {
		writeTextError(w, http.StatusBadRequest, "query must contain 1 to 100 characters")
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "prefix"
	}
	if mode != "prefix" && mode != "contains" {
		writeTextError(w, http.StatusBadRequest, "mode must be prefix or contains")
		return
	}
	limit, err := searchLimit(r)
	if err != nil {
		writeTextError(w, http.StatusBadRequest, err.Error())
		return
	}
	var results []string

	if mode == "contains" {
		results = words.SearchContains(query)
	} else {
		results = words.Search(query)
	}

	if len(results) > limit {
		results = results[:limit]
	}

	textResponse(w, strings.Join(results, "\n"))
}

func apiAllWordsHandler(w http.ResponseWriter, r *http.Request) {
	allWords := words.All()

	// Pagination
	page := 1
	perPage := 100

	if p := r.URL.Query().Get("page"); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 {
			writeAPIError(w, http.StatusBadRequest, "page must be a positive integer")
			return
		}
		page = n
	}
	if pp := r.URL.Query().Get("per_page"); pp != "" {
		n, err := strconv.Atoi(pp)
		if err != nil || n < 1 || n > 1000 {
			writeAPIError(w, http.StatusBadRequest, "per_page must be between 1 and 1000")
			return
		}
		perPage = n
	}

	totalPages := (len(allWords) + perPage - 1) / perPage
	if page > totalPages {
		page = totalPages + 1
	}
	start := len(allWords)
	end := len(allWords)
	if page <= totalPages {
		start = (page - 1) * perPage
		end = min(start+perPage, len(allWords))
	}

	pageWords := allWords[start:end]

	jsonResponse(w, map[string]interface{}{
		"words":       pageWords,
		"count":       len(pageWords),
		"total":       len(allWords),
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
		"has_next":    page < totalPages,
		"has_prev":    page > 1,
	})
}

func apiAllWordsTextHandler(w http.ResponseWriter, r *http.Request) {
	allWords := words.All()
	textResponse(w, strings.Join(allWords, "\n"))
}

func apiByDiceCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	word, ok := words.ByDiceCode(code)
	if !ok {
		writeAPIError(w, http.StatusNotFound, "no word is mapped to that dice code")
		return
	}
	jsonResponse(w, map[string]string{"dice": code, "word": word})
}

func apiByDiceCodeTextHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	word, ok := words.ByDiceCode(strings.TrimSuffix(code, ".txt"))
	if !ok {
		writeTextError(w, http.StatusNotFound, "no word is mapped to that dice code")
		return
	}
	textResponse(w, word)
}

func apiByLetterHandler(w http.ResponseWriter, r *http.Request) {
	letter := chi.URLParam(r, "letter")
	if len(letter) != 1 || (letter[0] < 'a' || letter[0] > 'z') && (letter[0] < 'A' || letter[0] > 'Z') {
		writeAPIError(w, http.StatusBadRequest, "letter must be one ASCII letter")
		return
	}
	result := words.ByLetter(letter)

	jsonResponse(w, map[string]interface{}{
		"letter": letter,
		"count":  len(result),
		"words":  result,
	})
}

func apiByLetterTextHandler(w http.ResponseWriter, r *http.Request) {
	letter := chi.URLParam(r, "letter")
	if len(letter) != 1 || (letter[0] < 'a' || letter[0] > 'z') && (letter[0] < 'A' || letter[0] > 'Z') {
		writeTextError(w, http.StatusBadRequest, "letter must be one ASCII letter")
		return
	}
	result := words.ByLetter(letter)
	textResponse(w, strings.Join(result, "\n"))
}

func apiByLengthHandler(w http.ResponseWriter, r *http.Request) {
	length, err := strconv.Atoi(chi.URLParam(r, "length"))
	if err != nil || length < 1 || length > 100 {
		writeAPIError(w, http.StatusBadRequest, "length must be between 1 and 100")
		return
	}

	result := words.ByLength(length)

	jsonResponse(w, map[string]interface{}{
		"length": length,
		"count":  len(result),
		"words":  result,
	})
}

func apiByLengthTextHandler(w http.ResponseWriter, r *http.Request) {
	length, err := strconv.Atoi(chi.URLParam(r, "length"))
	if err != nil || length < 1 || length > 100 {
		writeTextError(w, http.StatusBadRequest, "length must be between 1 and 100")
		return
	}

	result := words.ByLength(length)
	textResponse(w, strings.Join(result, "\n"))
}

func apiStatsHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, words.Stats())
}

func apiLettersHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, map[string]interface{}{
		"letters": words.Letters(),
		"count":   len(words.Letters()),
	})
}

func apiCountHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, map[string]interface{}{
		"count": words.Count(),
	})
}

func apiCountTextHandler(w http.ResponseWriter, r *http.Request) {
	textResponse(w, strconv.Itoa(words.Count()))
}


func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'; connect-src 'self'; object-src 'none'; base-uri 'self'; frame-ancestors 'none'")
		next.ServeHTTP(w, r)
	})
}

func requestBodyLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}

type rateWindow struct {
	start time.Time
	count int
}

func newRateLimitMiddleware(limit int, window time.Duration) func(http.Handler) http.Handler {
	var mu sync.Mutex
	clients := make(map[netip.Addr]rateWindow)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				host = r.RemoteAddr
			}
			ip, err := netip.ParseAddr(strings.Trim(host, "[]"))
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ip = ip.Unmap()
			now := time.Now()
			mu.Lock()
			if len(clients) >= 10000 {
				for key, value := range clients {
					if now.Sub(value.start) >= window {
						delete(clients, key)
					}
				}
			}
			entry, exists := clients[ip]
			if !exists && len(clients) >= 10000 {
				mu.Unlock()
				w.Header().Set("Retry-After", strconv.Itoa(max(1, int(window.Seconds()))))
				writeAPIError(w, http.StatusTooManyRequests, "request capacity temporarily unavailable")
				return
			}
			if entry.start.IsZero() || now.Sub(entry.start) >= window {
				entry = rateWindow{start: now}
			}
			entry.count++
			clients[ip] = entry
			limited := entry.count > limit
			retry := max(1, int(window.Seconds()-now.Sub(entry.start).Seconds()))
			mu.Unlock()
			if limited {
				w.Header().Set("Retry-After", strconv.Itoa(retry))
				writeAPIError(w, http.StatusTooManyRequests, "too many requests; retry later")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func swaggerJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.oai.openapi+json;version=3.0")
	w.Header().Set("Cache-Control", "public, max-age=300")
	_, _ = w.Write(openAPISpec)
}

type graphQLRequest struct {
	Query string `json:"query"`
}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, `<!doctype html><html lang="en"><meta charset="utf-8"><title>GraphQL</title><h1>WordList GraphQL</h1><p>POST a JSON object with a query field to this endpoint.</p></html>`)
		return
	}
	var request graphQLRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&request); err != nil {
		writeAPIError(w, http.StatusBadRequest, "request must be valid JSON containing a query")
		return
	}
	query := strings.ToLower(request.Query)
	if len(request.Query) > 4096 || strings.TrimSpace(query) == "" {
		writeAPIError(w, http.StatusBadRequest, "query must contain 1 to 4096 characters")
		return
	}
	data := make(map[string]interface{})
	if strings.Contains(query, "wordcount") || strings.Contains(query, "word_count") || strings.Contains(query, "count") {
		data["wordCount"] = words.Count()
	}
	if strings.Contains(query, "stats") {
		data["stats"] = words.Stats()
	}
	if strings.Contains(query, "words") {
		data["words"] = words.All()
	}
	if strings.Contains(query, "passphrase") {
		writeAPIError(w, http.StatusBadRequest, "passphrases are available through the REST API only")
		return
	}
	if len(data) == 0 {
		writeAPIError(w, http.StatusBadRequest, "query requested no supported fields")
		return
	}
	jsonResponse(w, map[string]interface{}{"data": data})
}

func serviceWorkerHandler(w http.ResponseWriter, r *http.Request) {
	file, err := staticFS.Open("static/sw.js")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = io.Copy(w, file)
}


func apiInfoHandler(page, content string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, map[string]string{"page": page, "content": content})
	}
}


func searchLimit(r *http.Request) (int, error) {
	limit := 100
	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 1000 {
			return 0, errors.New("limit must be between 1 and 1000")
		}
		limit = parsed
	}
	return limit, nil
}
