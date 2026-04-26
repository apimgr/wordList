package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apimgr/wordList/src/paths"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	WebUI       WebUIConfig       `yaml:"web_ui"`
	WebRobots   WebRobotsConfig   `yaml:"web_robots"`
	WebSecurity WebSecurityConfig `yaml:"web_security"`
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	FQDN         string        `yaml:"fqdn"`
	Address      string        `yaml:"address"`
	Mode         string        `yaml:"mode"`
	UpdateBranch string        `yaml:"update_branch"`
	Logging      LoggingConfig `yaml:"logging"`
	Admin        AdminConfig   `yaml:"admin"`
	Session      SessionConfig `yaml:"session"`
}

type AdminConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	APIToken string `yaml:"api_token"`
}

type SessionConfig struct {
	Timeout int `yaml:"timeout"`
}

type LoggingConfig struct {
	AccessFormat string `yaml:"access_format"`
	Level        string `yaml:"level"`
}

type WebUIConfig struct {
	Theme   string `yaml:"theme"`
	Logo    string `yaml:"logo"`
	Favicon string `yaml:"favicon"`
}

type WebRobotsConfig struct {
	Allow []string `yaml:"allow"`
	Deny  []string `yaml:"deny"`
}

type WebSecurityConfig struct {
	Admin string `yaml:"admin"`
	CORS  string `yaml:"cors"`
}

// Default configuration
func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         "64851",
			FQDN:         "",
			Address:      "0.0.0.0",
			Mode:         "production",
			UpdateBranch: "stable",
			Logging: LoggingConfig{
				AccessFormat: "apache",
				Level:        "info",
			},
			Admin: AdminConfig{
				Username: "admin",
				Password: "",
				APIToken: "",
			},
			Session: SessionConfig{
				Timeout: 3600,
			},
		},
		WebUI: WebUIConfig{
			Theme:   "dark",
			Logo:    "",
			Favicon: "",
		},
		WebRobots: WebRobotsConfig{
			Allow: []string{"/", "/api"},
			Deny:  []string{"/admin"},
		},
		WebSecurity: WebSecurityConfig{
			Admin: "security@wordlist.apimgr.us",
			CORS:  "*",
		},
	}
}

// migrateYamlToYml migrates .yaml files to .yml extension
func migrateYamlToYml(path string) {
	// Only process if path ends with .yml
	if !strings.HasSuffix(path, ".yml") {
		return
	}

	// Check if .yaml version exists
	yamlPath := strings.TrimSuffix(path, ".yml") + ".yaml"
	if _, err := os.Stat(yamlPath); err == nil {
		// Old .yaml file exists, check if new .yml file doesn't exist
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Rename .yaml to .yml
			if err := os.Rename(yamlPath, path); err == nil {
				fmt.Printf("Migrated configuration: %s -> %s\n", yamlPath, path)
			}
		}
	}
}

// Load loads configuration from file or creates default
func Load() (*Config, error) {
	cfg := defaultConfig()

	configFile := filepath.Join(paths.ConfigDir(), "server.yml")

	// Migrate old .yaml file to .yml if needed
	migrateYamlToYml(configFile)

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create default config file
		if err := Save(cfg); err != nil {
			return cfg, err
		}
		return cfg, nil
	}

	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return cfg, err
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Save saves configuration to file
func Save(cfg *Config) error {
	configDir := paths.ConfigDir()

	// Create config directory if needed
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "server.yml")

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	// Add header comment
	content := "# wordList Configuration\n# https://wordlist.apimgr.us\n\n" + string(data)

	return os.WriteFile(configFile, []byte(content), 0644)
}
