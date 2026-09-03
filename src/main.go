package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/apimgr/wordList/src/config"
	"github.com/apimgr/wordList/src/mode"
	"github.com/apimgr/wordList/src/paths"
	"github.com/apimgr/wordList/src/server"
	"github.com/apimgr/wordList/src/words"
)

var (
	Version   = "1.0.0"
	BuildTime = "unknown"
)

func main() {
	// CLI flags
	showHelp := flag.Bool("help", false, "Show help")
	showVersion := flag.Bool("version", false, "Show version")
	showStatus := flag.Bool("status", false, "Show service status")
	configDir := flag.String("config", "", "Configuration directory")
	dataDir := flag.String("data", "", "Data directory")
	logsDir := flag.String("logs", "", "Logs directory")
	address := flag.String("address", "", "Listen address")
	port := flag.String("port", "", "Listen port")
	modeFlag := flag.String("mode", "", "Application mode (dev/development, prod/production)")
	updateCmd := flag.String("update", "", "Update command (stable, beta, nightly)")
	serviceCmd := flag.String("service", "", "Service command (install, uninstall, start, stop, restart, status)")
	maintenanceMode := flag.String("maintenance", "", "Maintenance mode (on/off)")

	flag.Parse()

	// Unused vars to satisfy compiler
	_ = logsDir

	// Handle help
	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	// Handle version
	if *showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	// Handle update command
	if *updateCmd != "" {
		handleUpdateCommand(*updateCmd)
		os.Exit(0)
	}

	// Initialize mode
	if err := mode.Initialize(*modeFlag); err != nil {
		log.Printf("Warning: invalid mode: %v", err)
	}

	// Handle status check
	if *showStatus {
		checkStatus()
		os.Exit(0)
	}

	// Handle service commands
	if *serviceCmd != "" {
		handleServiceCommand(*serviceCmd)
		os.Exit(0)
	}

	// Handle maintenance mode
	if *maintenanceMode != "" {
		handleMaintenanceMode(*maintenanceMode)
		os.Exit(0)
	}

	// Initialize paths
	paths.Init(*configDir, *dataDir)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override config with CLI flags
	if *address != "" {
		cfg.Server.Address = *address
	}
	if *port != "" {
		cfg.Server.Port = *port
	}

	// Load word data
	if err := words.Load(); err != nil {
		log.Fatalf("Failed to load words: %v", err)
	}

	// Create server
	srv := server.New(cfg, Version)

	// Channel to listen for errors
	errChan := make(chan error, 1)

	// Start server in goroutine
	go func() {
		printStartup(cfg)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for interrupt signal or error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		fmt.Println("\n🛑 Shutting down gracefully...")
	case err := <-errChan:
		log.Printf("Server error: %v", err)
	}

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Server stopped")
}

func printHelp() {
	fmt.Printf(`wordList v%s - EFF Large Wordlist API

Usage: wordList [options]

Options:
  --help             Show this help message
  --version          Show version information
  --status           Check service status
  --config DIR       Set configuration directory
  --data DIR         Set data directory
  --logs DIR         Set logs directory
  --address ADDR     Set listen address (default: 0.0.0.0)
  --port PORT        Set listen port (default: auto)
  --mode MODE        Application mode (dev, prod)
  --update BRANCH    Update from branch (stable, beta, nightly)

Service Management:
  --service install    Install as system service
  --service uninstall  Uninstall system service
  --service start      Start the service
  --service stop       Stop the service
  --service restart    Restart the service
  --service status     Show service status

Maintenance:
  --maintenance on     Enable maintenance mode
  --maintenance off    Disable maintenance mode

Examples:
  wordList                          Start with default settings
  wordList --port 8080              Start on port 8080
  wordList --mode dev --port 8080   Start in development mode
  wordList --update stable          Update to stable branch
  wordList --config /etc/wordlist   Use custom config directory

Documentation: https://wordlist.apimgr.us
`, Version)
}

func printStartup(cfg *config.Config) {
	fmt.Println()
	fmt.Printf("✅ wordList v%s started successfully\n", Version)
	fmt.Printf("📡 Listening on http://%s:%s\n", getDisplayAddress(cfg), cfg.Server.Port)
	fmt.Printf("📊 Swagger UI: http://%s:%s/swagger\n", getDisplayAddress(cfg), cfg.Server.Port)
	fmt.Printf("📚 API Docs: http://%s:%s/api\n", getDisplayAddress(cfg), cfg.Server.Port)
	fmt.Printf("📝 Words loaded: %d\n", words.Count())
	fmt.Println()
}

func getDisplayAddress(cfg *config.Config) string {
	if cfg.Server.FQDN != "" {
		return cfg.Server.FQDN
	}
	if cfg.Server.Address == "0.0.0.0" || cfg.Server.Address == "" {
		hostname, err := os.Hostname()
		if err == nil && hostname != "" {
			return hostname
		}
		return "localhost"
	}
	return cfg.Server.Address
}

func checkStatus() {
	// Try to connect to the server
	cfg, _ := config.Load()
	addr := fmt.Sprintf("http://localhost:%s/healthz", cfg.Server.Port)

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(addr)
	if err != nil {
		fmt.Println("❌ Service is not running")
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Service is running")
		os.Exit(0)
	}

	fmt.Println("⚠️ Service returned unexpected status")
	os.Exit(2)
}

func handleUpdateCommand(branch string) {
	validBranches := map[string]bool{
		"stable":  true,
		"beta":    true,
		"nightly": true,
	}

	if !validBranches[branch] {
		fmt.Printf("Error: invalid update branch %q (valid: stable, beta, nightly)\n", branch)
		os.Exit(1)
	}

	fmt.Printf("Updating wordList from %s branch...\n", branch)

	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("Error: git is not installed")
		os.Exit(1)
	}

	cmd := exec.Command("git", "pull", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Update failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Update complete. Please rebuild the application.")
}

func handleServiceCommand(cmd string) {
	switch cmd {
	case "install":
		fmt.Println("Service installation not yet implemented")
		fmt.Println("Use systemd/launchd/rc.d to manage the service")
	case "uninstall":
		fmt.Println("Service uninstallation not yet implemented")
	case "start":
		fmt.Println("Use 'systemctl start wordlist' or run the binary directly")
	case "stop":
		fmt.Println("Use 'systemctl stop wordlist' or send SIGTERM to the process")
	case "restart":
		fmt.Println("Use 'systemctl restart wordlist'")
	case "status":
		fmt.Println("Use 'systemctl status wordlist' or --status flag")
	default:
		fmt.Printf("Unknown service command: %s\n", cmd)
		fmt.Println("Available commands: install, uninstall, start, stop, restart, status")
	}
}

func handleMaintenanceMode(m string) {
	switch m {
	case "on":
		fmt.Println("Maintenance mode: ON")
		fmt.Println("Note: Maintenance mode is handled at runtime, not persisted")
	case "off":
		fmt.Println("Maintenance mode: OFF")
	default:
		fmt.Printf("Invalid maintenance mode: %s (use 'on' or 'off')\n", m)
	}
}
