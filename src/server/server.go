package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apimgr/wordList/src/admin"
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
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.WebSecurity.CORS},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Initialize admin handler
	adminHandler := admin.NewHandler(
		cfg.Server.Admin.Username,
		cfg.Server.Admin.Password,
		cfg.Server.Admin.APIToken,
		cfg.Server.Session.Timeout,
		false, // SSL enabled
		version,
	)

	// Register admin routes
	adminHandler.RegisterRoutes(r)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Web routes
	r.Get("/", homeHandler(cfg))
	r.Get("/generate", generatePageHandler(cfg))
	r.Get("/search", searchPageHandler(cfg))
	r.Get("/api", apiDocsHandler(cfg))
	r.Get("/swagger", swaggerHandler(cfg))

	// Health check
	r.Get("/healthz", healthHandler)

	// Special files
	r.Get("/robots.txt", robotsHandler(cfg))
	r.Get("/security.txt", securityHandler(cfg))
	r.Get("/.well-known/security.txt", securityHandler(cfg))
	r.Get("/manifest.json", manifestHandler(cfg))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
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
		r.Get("/letter/{letter}", apiByLetterHandler)
		r.Get("/letter/{letter}.txt", apiByLetterTextHandler)
		r.Get("/length/{length}", apiByLengthHandler)
		r.Get("/length/{length}.txt", apiByLengthTextHandler)

		// Stats
		r.Get("/stats", apiStatsHandler)
		r.Get("/letters", apiLettersHandler)
		r.Get("/count", apiCountHandler)
		r.Get("/count.txt", apiCountTextHandler)
	})

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Address, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

// Template data
type PageData struct {
	SiteTitle string
	BaseURL   string
	Theme     string
}

func newPageData(cfg *config.Config) PageData {
	return PageData{
		SiteTitle: "WordList",
		BaseURL:   fmt.Sprintf("http://localhost:%s", cfg.Server.Port),
		Theme:     cfg.WebUI.Theme,
	}
}

// Web handlers
func homeHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newPageData(cfg)
		templates.ExecuteTemplate(w, "home.html", data)
	}
}

func generatePageHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newPageData(cfg)
		templates.ExecuteTemplate(w, "generate.html", data)
	}
}

func searchPageHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newPageData(cfg)
		templates.ExecuteTemplate(w, "search.html", data)
	}
}

func apiDocsHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newPageData(cfg)
		templates.ExecuteTemplate(w, "openapi.html", data)
	}
}

func swaggerHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api", http.StatusFound)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "healthy",
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"word_count":  words.Count(),
		"service":     "wordList",
		"version":     "1.0.0",
	})
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
				{"src": "/static/icon-192.png", "sizes": "192x192", "type": "image/png"},
				{"src": "/static/icon-512.png", "sizes": "512x512", "type": "image/png"},
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
func apiRandomHandler(w http.ResponseWriter, r *http.Request) {
	count := 6
	if c := chi.URLParam(r, "count"); c != "" {
		if n, err := strconv.Atoi(c); err == nil && n > 0 && n <= 100 {
			count = n
		}
	}

	randomWords := words.Random(count)
	jsonResponse(w, map[string]interface{}{
		"count": len(randomWords),
		"words": randomWords,
	})
}

func apiRandomTextHandler(w http.ResponseWriter, r *http.Request) {
	count := 6
	if c := chi.URLParam(r, "count"); c != "" {
		if n, err := strconv.Atoi(c); err == nil && n > 0 && n <= 100 {
			count = n
		}
	}

	randomWords := words.Random(count)
	textResponse(w, strings.Join(randomWords, "\n"))
}

func apiPassphraseHandler(w http.ResponseWriter, r *http.Request) {
	count := 6
	if c := chi.URLParam(r, "count"); c != "" {
		if n, err := strconv.Atoi(c); err == nil && n >= 3 && n <= 20 {
			count = n
		}
	}

	separator := r.URL.Query().Get("separator")
	if separator == "" {
		separator = "-"
	}

	capitalize := r.URL.Query().Get("capitalize") != "false"

	passphrase := words.Passphrase(count, separator, capitalize)
	passphraseWords := strings.Split(passphrase, separator)

	jsonResponse(w, map[string]interface{}{
		"passphrase":  passphrase,
		"words":       passphraseWords,
		"word_count":  count,
		"separator":   separator,
		"capitalized": capitalize,
	})
}

func apiPassphraseTextHandler(w http.ResponseWriter, r *http.Request) {
	count := 6
	if c := chi.URLParam(r, "count"); c != "" {
		if n, err := strconv.Atoi(c); err == nil && n >= 3 && n <= 20 {
			count = n
		}
	}

	separator := r.URL.Query().Get("separator")
	if separator == "" {
		separator = "-"
	}

	capitalize := r.URL.Query().Get("capitalize") != "false"

	passphrase := words.Passphrase(count, separator, capitalize)
	textResponse(w, passphrase)
}

func apiSearchHandler(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	if query == "" {
		jsonResponse(w, map[string]interface{}{
			"error": "Query parameter is required",
		})
		return
	}

	mode := r.URL.Query().Get("mode")
	var results []string

	if mode == "contains" {
		results = words.SearchContains(query)
	} else {
		results = words.Search(query)
	}

	// Limit results
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 1000 {
			limit = n
		}
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
	if query == "" {
		textResponse(w, "Query parameter is required")
		return
	}

	mode := r.URL.Query().Get("mode")
	var results []string

	if mode == "contains" {
		results = words.SearchContains(query)
	} else {
		results = words.Search(query)
	}

	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 1000 {
			limit = n
		}
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
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}
	if pp := r.URL.Query().Get("per_page"); pp != "" {
		if n, err := strconv.Atoi(pp); err == nil && n > 0 && n <= 1000 {
			perPage = n
		}
	}

	start := (page - 1) * perPage
	end := start + perPage

	if start >= len(allWords) {
		start = len(allWords)
	}
	if end > len(allWords) {
		end = len(allWords)
	}

	pageWords := allWords[start:end]
	totalPages := (len(allWords) + perPage - 1) / perPage

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

func apiByLetterHandler(w http.ResponseWriter, r *http.Request) {
	letter := chi.URLParam(r, "letter")
	result := words.ByLetter(letter)

	jsonResponse(w, map[string]interface{}{
		"letter": letter,
		"count":  len(result),
		"words":  result,
	})
}

func apiByLetterTextHandler(w http.ResponseWriter, r *http.Request) {
	letter := chi.URLParam(r, "letter")
	result := words.ByLetter(letter)
	textResponse(w, strings.Join(result, "\n"))
}

func apiByLengthHandler(w http.ResponseWriter, r *http.Request) {
	length, err := strconv.Atoi(chi.URLParam(r, "length"))
	if err != nil || length < 1 {
		jsonResponse(w, map[string]interface{}{
			"error": "Invalid length parameter",
		})
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
	if err != nil || length < 1 {
		textResponse(w, "Invalid length parameter")
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
