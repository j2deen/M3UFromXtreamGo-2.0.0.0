package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// LoadConfig loads configuration with priority: env vars > config file > defaults
func LoadConfig() (*Config, error) {
	cfg := getDefaults()

	// Load from config file if it exists
	configPath := getConfigFilePath()
	if _, err := os.Stat(configPath); err == nil {
		if err := loadConfigFile(configPath, cfg); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	// Override with environment variables
	applyEnvVars(cfg)

	return cfg, nil
}

// getDefaults returns default configuration values
func getDefaults() *Config {
	return &Config{
		Mode: "web",
		Server: ServerConfig{
			Port:            8080,
			Host:            "0.0.0.0",
			ReadTimeoutSec:  30,
			WriteTimeoutSec: 30,
		},
		Xtream: XtreamConfig{
			BaseURL:           "",
			Username:          "",
			Password:          "",
			RequestTimeoutSec: 30,
		},
		Output: OutputConfig{
			DefaultFilename: "playlist.m3u",
		},
		Logging: LoggingConfig{
			Level: "INFO",
		},
	}
}

// getConfigFilePath returns the configuration file path
func getConfigFilePath() string {
	if path := os.Getenv("M3U_CONFIG_FILE"); path != "" {
		return path
	}
	return "./config.json"
}

// loadConfigFile loads configuration from a JSON file
func loadConfigFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, cfg)
}

// applyEnvVars overrides configuration with environment variables
func applyEnvVars(cfg *Config) {
	if mode := os.Getenv("M3U_MODE"); mode != "" {
		cfg.Mode = mode
	}
	if port := os.Getenv("M3U_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}
	if host := os.Getenv("M3U_SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if readTimeout := os.Getenv("M3U_SERVER_READ_TIMEOUT"); readTimeout != "" {
		if t, err := strconv.Atoi(readTimeout); err == nil {
			cfg.Server.ReadTimeoutSec = t
		}
	}
	if writeTimeout := os.Getenv("M3U_SERVER_WRITE_TIMEOUT"); writeTimeout != "" {
		if t, err := strconv.Atoi(writeTimeout); err == nil {
			cfg.Server.WriteTimeoutSec = t
		}
	}
	if baseURL := os.Getenv("M3U_XTREAM_BASE_URL"); baseURL != "" {
		cfg.Xtream.BaseURL = baseURL
	}
	if username := os.Getenv("M3U_XTREAM_USERNAME"); username != "" {
		cfg.Xtream.Username = username
	}
	if password := os.Getenv("M3U_XTREAM_PASSWORD"); password != "" {
		cfg.Xtream.Password = password
	}
	if requestTimeout := os.Getenv("M3U_XTREAM_REQUEST_TIMEOUT"); requestTimeout != "" {
		if t, err := strconv.Atoi(requestTimeout); err == nil {
			cfg.Xtream.RequestTimeoutSec = t
		}
	}
	if filename := os.Getenv("M3U_OUTPUT_FILENAME"); filename != "" {
		cfg.Output.DefaultFilename = filename
	}
	if logLevel := os.Getenv("M3U_LOG_LEVEL"); logLevel != "" {
		cfg.Logging.Level = logLevel
	}
}

// Validate ensures required configuration fields are present
func (c *Config) Validate() error {
	if c.Xtream.BaseURL == "" {
		return fmt.Errorf("xtream base_url is required")
	}
	if c.Xtream.Username == "" {
		return fmt.Errorf("xtream username is required")
	}
	if c.Xtream.Password == "" {
		return fmt.Errorf("xtream password is required")
	}
	if c.Mode == "web" && c.Server.Port == 0 {
		return fmt.Errorf("server port is required for web mode")
	}
	return nil
}
