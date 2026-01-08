package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Printf("M3UFromXtream Version %s\n\n", version)

	var err error
	config, err = LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger with configured level
	InitLogger(config.Logging.Level, os.Stdout)
	LogInfo("M3UFromXtream starting - Version %s", version)
	LogDebug("Log level set to: %s", config.Logging.Level)

	// Detect mode: CLI vs Web
	if isCLIMode() {
		LogInfo("Running in CLI mode")
		runCLI()
	} else {
		LogInfo("Running in web server mode")
		runWebServer()
	}
}

// isCLIMode checks if the application is being run in CLI mode
func isCLIMode() bool {
	// Check if command-line arguments indicate CLI mode (3-4 args: url, username, password, [output-file])
	if len(os.Args) >= 4 && len(os.Args) <= 5 {
		// Make sure the first arg doesn't look like a flag
		return !strings.HasPrefix(os.Args[1], "-")
	}
	return false
}

// runCLI executes the application in CLI mode
func runCLI() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: M3UFromXtream <url> <username> <password> [output-file]")
		fmt.Println("Example: M3UFromXtream http://example.com:8080 user pass output.m3u")
		os.Exit(1)
	}

	baseURL := os.Args[1]
	username := os.Args[2]
	password := os.Args[3]
	outputFile := "playlist.m3u"
	if len(os.Args) > 4 {
		outputFile = os.Args[4]
	}

	LogInfo("CLI Mode - Output file: %s", outputFile)
	LogDebug("Base URL: %s, Username: %s", baseURL, username)

	// Create XtreamConfig from CLI arguments
	xtreamCfg := XtreamConfig{
		BaseURL:           baseURL,
		Username:          username,
		Password:          password,
		RequestTimeoutSec: config.Xtream.RequestTimeoutSec,
	}

	content, err := generateM3UContent(&xtreamCfg)
	if err != nil {
		LogError("Failed to generate M3U content: %v", err)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
		LogError("Failed to write output file: %v", err)
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	LogInfo("M3U playlist successfully written to: %s", outputFile)
	fmt.Printf("M3U playlist successfully created: %s\n", outputFile)
}

// runWebServer executes the application in web server mode
func runWebServer() {
	if err := config.Validate(); err != nil {
		LogError("Configuration validation failed: %v", err)
		fmt.Printf("Configuration error: %v\n", err)
		fmt.Println("\nPlease set the following environment variables:")
		fmt.Println("  M3U_XTREAM_BASE_URL    - Xtream API URL")
		fmt.Println("  M3U_XTREAM_USERNAME    - Xtream username")
		fmt.Println("  M3U_XTREAM_PASSWORD    - Xtream password")
		fmt.Println("\nOr create a config.json file with the required settings.")
		os.Exit(1)
	}

	LogInfo("Configuration validated successfully")
	LogDebug("Server will listen on %s:%d", config.Server.Host, config.Server.Port)
	LogDebug("Xtream API URL: %s", config.Xtream.BaseURL)

	mux := http.NewServeMux()

	// Register handlers with logging middleware
	mux.HandleFunc("/m3u", loggingMiddleware(handlerM3U))
	mux.HandleFunc("/health", loggingMiddleware(handlerHealth))
	mux.HandleFunc("/config", loggingMiddleware(handlerConfig))

	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  time.Duration(config.Server.ReadTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeoutSec) * time.Second,
	}

	LogInfo("Starting HTTP server on %s", addr)
	LogInfo("Endpoints registered: /m3u, /health, /config")

	fmt.Printf("Starting web server on %s\n", addr)
	fmt.Printf("Endpoints:\n")
	fmt.Printf("  - M3U Playlist:  http://%s/m3u\n", getDisplayAddr(config.Server.Host, config.Server.Port))
	fmt.Printf("  - Health Check:  http://%s/health\n", getDisplayAddr(config.Server.Host, config.Server.Port))
	fmt.Printf("  - Config Info:   http://%s/config\n", getDisplayAddr(config.Server.Host, config.Server.Port))
	fmt.Println()

	if err := server.ListenAndServe(); err != nil {
		LogError("Server failed to start: %v", err)
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

// getDisplayAddr returns a user-friendly address for display purposes
func getDisplayAddr(host string, port int) string {
	if host == "0.0.0.0" || host == "" {
		return "localhost:" + strconv.Itoa(port)
	}
	return host + ":" + strconv.Itoa(port)
}
